// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// RepoMap represents a compact map of the repository structure
type RepoMap struct {
	ProjectRoot string         `json:"projectRoot"`
	Entries     []RepoMapEntry `json:"entries"`
	TotalTokens int            `json:"totalTokens"`
	GeneratedAt time.Time      `json:"generatedAt"`
}

// RepoMapEntry represents a single file entry in the repo map
type RepoMapEntry struct {
	FilePath string            `json:"filePath"`
	Symbols  []SymbolSignature `json:"symbols"`
	Rank     float64           `json:"rank"` // PageRank score
}

// SymbolSignature represents a symbol with its signature
type SymbolSignature struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`      // class, function, method, interface, etc.
	Signature string `json:"signature"` // e.g. "login(email: string): Promise<User>"
	Line      int    `json:"line"`
}

// RepoMapBuilder builds compact repository maps for AI context
type RepoMapBuilder interface {
	// BuildMap generates a repository map within token budget
	BuildMap(ctx context.Context, projectRoot string, maxTokens int) (*RepoMap, error)

	// InvalidateCache clears cached map for a project
	InvalidateCache(projectRoot string)
}

// RepoMapBuilderImpl implements RepoMapBuilder
type RepoMapBuilderImpl struct {
	analyzerRegistry analysis.AnalyzerRegistry
	fileRanker       FileRanker
	log              domain.Logger

	mu    sync.RWMutex
	cache map[string]*repoMapCacheEntry
}

type repoMapCacheEntry struct {
	repoMap   *RepoMap
	expiresAt time.Time
}

const (
	cacheExpiration     = 5 * time.Minute
	defaultMaxTokens    = 8000
	tokensPerChar       = 4 // ~4 chars per token
	maxFilesInMap       = 200
	maxSymbolsPerFile   = 50
	signatureMaxLen     = 120
	minRankThreshold    = 0.001
)

// NewRepoMapBuilder creates a new RepoMapBuilder
func NewRepoMapBuilder(
	analyzerRegistry analysis.AnalyzerRegistry,
	fileRanker FileRanker,
	log domain.Logger,
) *RepoMapBuilderImpl {
	return &RepoMapBuilderImpl{
		analyzerRegistry: analyzerRegistry,
		fileRanker:       fileRanker,
		log:              log,
		cache:            make(map[string]*repoMapCacheEntry),
	}
}

// BuildMap generates a repository map within the token budget
func (b *RepoMapBuilderImpl) BuildMap(ctx context.Context, projectRoot string, maxTokens int) (*RepoMap, error) {
	if maxTokens <= 0 {
		maxTokens = defaultMaxTokens
	}

	// Check cache
	if cached := b.getCached(projectRoot); cached != nil && cached.TotalTokens <= maxTokens {
		return cached, nil
	}

	b.log.Info(fmt.Sprintf("Building repo map for %s with max %d tokens", projectRoot, maxTokens))

	// Get ranked files
	rankedFiles, err := b.fileRanker.RankFiles(ctx, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to rank files: %w", err)
	}

	// Build entries for top-ranked files
	entries := make([]RepoMapEntry, 0, len(rankedFiles))
	totalTokens := 0
	filesProcessed := 0

	for _, rf := range rankedFiles {
		if filesProcessed >= maxFilesInMap {
			break
		}

		if rf.Rank < minRankThreshold {
			continue
		}

		entry, entryTokens, err := b.buildEntry(ctx, projectRoot, rf)
		if err != nil {
			b.log.Warning(fmt.Sprintf("Failed to build entry for %s: %v", rf.FilePath, err))
			continue
		}

		// Check if adding this entry exceeds budget
		if totalTokens+entryTokens > maxTokens {
			// Try to fit with fewer symbols
			entry, entryTokens = b.truncateEntry(entry, maxTokens-totalTokens)
			if entryTokens == 0 {
				continue
			}
		}

		entries = append(entries, entry)
		totalTokens += entryTokens
		filesProcessed++

		if totalTokens >= maxTokens {
			break
		}
	}

	repoMap := &RepoMap{
		ProjectRoot: projectRoot,
		Entries:     entries,
		TotalTokens: totalTokens,
		GeneratedAt: time.Now(),
	}

	// Cache the result
	b.setCache(projectRoot, repoMap)

	b.log.Info(fmt.Sprintf("Built repo map: %d files, %d tokens", len(entries), totalTokens))

	return repoMap, nil
}

// InvalidateCache clears the cache for a project
func (b *RepoMapBuilderImpl) InvalidateCache(projectRoot string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.cache, projectRoot)
}

// getCached returns cached repo map if valid
func (b *RepoMapBuilderImpl) getCached(projectRoot string) *RepoMap {
	b.mu.RLock()
	defer b.mu.RUnlock()

	entry, ok := b.cache[projectRoot]
	if !ok || time.Now().After(entry.expiresAt) {
		return nil
	}
	return entry.repoMap
}

// setCache stores repo map in cache
func (b *RepoMapBuilderImpl) setCache(projectRoot string, repoMap *RepoMap) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.cache[projectRoot] = &repoMapCacheEntry{
		repoMap:   repoMap,
		expiresAt: time.Now().Add(cacheExpiration),
	}
}

// buildEntry builds a RepoMapEntry for a file
func (b *RepoMapBuilderImpl) buildEntry(ctx context.Context, projectRoot string, rf RankedFile) (RepoMapEntry, int, error) {
	fullPath := filepath.Join(projectRoot, rf.FilePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return RepoMapEntry{}, 0, fmt.Errorf("failed to read file: %w", err)
	}

	analyzer := b.analyzerRegistry.GetAnalyzer(rf.FilePath)
	if analyzer == nil {
		// No analyzer available, create minimal entry
		return RepoMapEntry{
			FilePath: rf.FilePath,
			Symbols:  []SymbolSignature{},
			Rank:     rf.Rank,
		}, b.estimateTokens(rf.FilePath + ":\n"), nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, rf.FilePath, content)
	if err != nil {
		b.log.Warning(fmt.Sprintf("Failed to extract symbols from %s: %v", rf.FilePath, err))
		return RepoMapEntry{
			FilePath: rf.FilePath,
			Symbols:  []SymbolSignature{},
			Rank:     rf.Rank,
		}, b.estimateTokens(rf.FilePath + ":\n"), nil
	}

	// Convert to signatures and filter
	signatures := b.symbolsToSignatures(symbols)

	// Limit symbols per file
	if len(signatures) > maxSymbolsPerFile {
		signatures = signatures[:maxSymbolsPerFile]
	}

	entry := RepoMapEntry{
		FilePath: rf.FilePath,
		Symbols:  signatures,
		Rank:     rf.Rank,
	}

	tokens := b.estimateEntryTokens(entry)

	return entry, tokens, nil
}

// symbolsToSignatures converts analysis symbols to signatures
func (b *RepoMapBuilderImpl) symbolsToSignatures(symbols []analysis.Symbol) []SymbolSignature {
	signatures := make([]SymbolSignature, 0, len(symbols))

	for _, sym := range symbols {
		// Skip non-important symbols
		if !b.isImportantSymbol(sym) {
			continue
		}

		sig := SymbolSignature{
			Name: sym.Name,
			Kind: string(sym.Kind),
			Line: sym.StartLine,
		}

		// Use existing signature or build one
		if sym.Signature != "" {
			sig.Signature = truncateString(sym.Signature, signatureMaxLen)
		} else {
			sig.Signature = b.buildSignature(sym)
		}

		signatures = append(signatures, sig)
	}

	// Sort by line number
	sort.Slice(signatures, func(i, j int) bool {
		return signatures[i].Line < signatures[j].Line
	})

	return signatures
}

// isImportantSymbol checks if a symbol should be included in the map
func (b *RepoMapBuilderImpl) isImportantSymbol(sym analysis.Symbol) bool {
	switch sym.Kind {
	case analysis.KindFunction, analysis.KindMethod, analysis.KindClass,
		analysis.KindInterface, analysis.KindStruct, analysis.KindType,
		analysis.KindComponent, analysis.KindComposable:
		return true
	case analysis.KindConstant, analysis.KindVariable:
		// Include only exported constants/variables
		return len(sym.Name) > 0 && sym.Name[0] >= 'A' && sym.Name[0] <= 'Z'
	default:
		return false
	}
}

// buildSignature builds a signature string for a symbol
func (b *RepoMapBuilderImpl) buildSignature(sym analysis.Symbol) string {
	switch sym.Kind {
	case analysis.KindClass:
		return fmt.Sprintf("class %s", sym.Name)
	case analysis.KindInterface:
		return fmt.Sprintf("interface %s", sym.Name)
	case analysis.KindStruct:
		return fmt.Sprintf("struct %s", sym.Name)
	case analysis.KindFunction:
		return fmt.Sprintf("func %s(...)", sym.Name)
	case analysis.KindMethod:
		if sym.Parent != "" {
			return fmt.Sprintf("func (%s) %s(...)", sym.Parent, sym.Name)
		}
		return fmt.Sprintf("func %s(...)", sym.Name)
	default:
		return sym.Name
	}
}

// truncateEntry reduces entry size to fit within token budget
func (b *RepoMapBuilderImpl) truncateEntry(entry RepoMapEntry, maxTokens int) (RepoMapEntry, int) {
	if maxTokens <= 0 {
		return RepoMapEntry{}, 0
	}

	// Start with file path only
	baseTokens := b.estimateTokens(entry.FilePath + ":\n")
	if baseTokens > maxTokens {
		return RepoMapEntry{}, 0
	}

	// Add symbols until we hit the limit
	truncated := RepoMapEntry{
		FilePath: entry.FilePath,
		Symbols:  make([]SymbolSignature, 0),
		Rank:     entry.Rank,
	}

	currentTokens := baseTokens
	for _, sym := range entry.Symbols {
		symTokens := b.estimateSymbolTokens(sym)
		if currentTokens+symTokens > maxTokens {
			break
		}
		truncated.Symbols = append(truncated.Symbols, sym)
		currentTokens += symTokens
	}

	return truncated, currentTokens
}

// estimateEntryTokens estimates tokens for an entry
func (b *RepoMapBuilderImpl) estimateEntryTokens(entry RepoMapEntry) int {
	tokens := b.estimateTokens(entry.FilePath + ":\n")
	for _, sym := range entry.Symbols {
		tokens += b.estimateSymbolTokens(sym)
	}
	return tokens
}

// estimateSymbolTokens estimates tokens for a symbol
func (b *RepoMapBuilderImpl) estimateSymbolTokens(sym SymbolSignature) int {
	// Format: "  kind name signature\n"
	line := fmt.Sprintf("  %s %s", sym.Kind, sym.Signature)
	return b.estimateTokens(line)
}

// estimateTokens estimates token count for text
func (b *RepoMapBuilderImpl) estimateTokens(text string) int {
	return (len(text) + tokensPerChar - 1) / tokensPerChar
}

// Format formats the repo map as a string for AI context
func (rm *RepoMap) Format() string {
	var sb strings.Builder

	sb.WriteString("# Repository Map\n\n")

	for _, entry := range rm.Entries {
		sb.WriteString(entry.FilePath)
		sb.WriteString(":\n")

		for _, sym := range entry.Symbols {
			sb.WriteString("  ")
			sb.WriteString(sym.Kind)
			sb.WriteString(" ")
			sb.WriteString(sym.Signature)
			sb.WriteString("\n")
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatCompact formats the repo map in a more compact form
func (rm *RepoMap) FormatCompact() string {
	var sb strings.Builder

	for _, entry := range rm.Entries {
		sb.WriteString(entry.FilePath)
		sb.WriteString(": ")

		names := make([]string, 0, len(entry.Symbols))
		for _, sym := range entry.Symbols {
			names = append(names, sym.Name)
		}
		sb.WriteString(strings.Join(names, ", "))
		sb.WriteString("\n")
	}

	return sb.String()
}

// truncateString truncates a string to maxLen
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
