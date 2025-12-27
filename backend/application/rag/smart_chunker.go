// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// ChunkOptions configures how files are chunked
type ChunkOptions struct {
	MaxTokensPerChunk int  // Maximum tokens per chunk (default: 500)
	IncludeContext    int  // Lines of context before/after (default: 3)
	PreserveStructure bool // Keep class/function boundaries (default: true)
}

// Chunk represents a piece of code extracted from a file
type Chunk struct {
	Content   string  `json:"content"`
	StartLine int     `json:"startLine"`
	EndLine   int     `json:"endLine"`
	Type      string  `json:"type"`      // "function", "class", "method", "block", "file"
	Name      string  `json:"name"`      // function/class name if applicable
	Tokens    int     `json:"tokens"`    // estimated token count
	Relevance float64 `json:"relevance"` // 0-1, if query provided
}

// SmartChunker splits files into meaningful chunks using AST analysis
type SmartChunker interface {
	// ChunkFile splits a file into chunks based on AST structure
	ChunkFile(ctx context.Context, filePath string, opts ChunkOptions) ([]Chunk, error)

	// ChunkByRelevance returns only relevant chunks for a query
	ChunkByRelevance(ctx context.Context, filePath string, query string, maxTokens int) ([]Chunk, error)

	// ChunkContent chunks content directly without reading from file
	ChunkContent(ctx context.Context, content []byte, filePath string, opts ChunkOptions) ([]Chunk, error)
}

// SmartChunkerImpl implements SmartChunker using language analyzers
type SmartChunkerImpl struct {
	analyzerRegistry analysis.AnalyzerRegistry
	log              domain.Logger
}

// Chunking constants
const (
	defaultMaxTokensPerChunk = 500
	defaultContextLines      = 3
	charsPerToken            = 4
	minChunkTokens           = 50
	summaryMaxLines          = 20
)

// NewSmartChunker creates a new SmartChunker
func NewSmartChunker(
	analyzerRegistry analysis.AnalyzerRegistry,
	log domain.Logger,
) *SmartChunkerImpl {
	return &SmartChunkerImpl{
		analyzerRegistry: analyzerRegistry,
		log:              log,
	}
}

// ChunkFile splits a file into chunks based on AST structure
func (c *SmartChunkerImpl) ChunkFile(ctx context.Context, filePath string, opts ChunkOptions) ([]Chunk, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return c.ChunkContent(ctx, content, filePath, opts)
}

// ChunkContent chunks content directly without reading from file
func (c *SmartChunkerImpl) ChunkContent(ctx context.Context, content []byte, filePath string, opts ChunkOptions) ([]Chunk, error) {
	c.setDefaultOptions(&opts)

	analyzer := c.analyzerRegistry.GetAnalyzer(filePath)
	if analyzer == nil {
		return c.chunkByLines(content, filePath, opts), nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		c.log.Warning(fmt.Sprintf("AST parsing failed for %s, falling back to line-based: %v", filePath, err))
		return c.chunkByLines(content, filePath, opts), nil
	}

	if len(symbols) == 0 {
		return c.chunkByLines(content, filePath, opts), nil
	}

	return c.chunkBySymbols(content, filePath, symbols, opts), nil
}

// ChunkByRelevance returns only relevant chunks for a query
func (c *SmartChunkerImpl) ChunkByRelevance(ctx context.Context, filePath string, query string, maxTokens int) ([]Chunk, error) {
	opts := ChunkOptions{
		MaxTokensPerChunk: defaultMaxTokensPerChunk,
		IncludeContext:    defaultContextLines,
		PreserveStructure: true,
	}

	chunks, err := c.ChunkFile(ctx, filePath, opts)
	if err != nil {
		return nil, err
	}

	queryTerms := c.extractQueryTerms(query)
	for i := range chunks {
		chunks[i].Relevance = c.calculateRelevance(chunks[i], queryTerms)
	}

	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].Relevance > chunks[j].Relevance
	})

	return c.selectChunksWithinBudget(chunks, maxTokens), nil
}

func (c *SmartChunkerImpl) setDefaultOptions(opts *ChunkOptions) {
	if opts.MaxTokensPerChunk == 0 {
		opts.MaxTokensPerChunk = defaultMaxTokensPerChunk
	}
	if opts.IncludeContext == 0 {
		opts.IncludeContext = defaultContextLines
	}
}

func (c *SmartChunkerImpl) chunkBySymbols(content []byte, filePath string, symbols []analysis.Symbol, opts ChunkOptions) []Chunk {
	lines := strings.Split(string(content), "\n")
	chunks := make([]Chunk, 0, len(symbols))

	topLevelSymbols := c.filterTopLevelSymbols(symbols)

	for _, sym := range topLevelSymbols {
		chunk := c.createChunkFromSymbol(lines, sym, opts)
		if chunk.Tokens >= minChunkTokens {
			chunks = append(chunks, chunk)
		}
	}

	if len(chunks) == 0 {
		return c.chunkByLines(content, filePath, opts)
	}

	chunks = c.fillGaps(lines, chunks, opts)
	chunks = c.splitLargeChunks(chunks, opts)

	return chunks
}

func (c *SmartChunkerImpl) filterTopLevelSymbols(symbols []analysis.Symbol) []analysis.Symbol {
	topLevel := make([]analysis.Symbol, 0, len(symbols))

	for _, sym := range symbols {
		if sym.Parent == "" && c.isChunkableSymbol(sym.Kind) {
			topLevel = append(topLevel, sym)
		}
	}

	sort.Slice(topLevel, func(i, j int) bool {
		return topLevel[i].StartLine < topLevel[j].StartLine
	})

	return topLevel
}

func (c *SmartChunkerImpl) isChunkableSymbol(kind analysis.SymbolKind) bool {
	switch kind {
	case analysis.KindFunction, analysis.KindMethod, analysis.KindClass,
		analysis.KindStruct, analysis.KindInterface, analysis.KindType:
		return true
	default:
		return false
	}
}

func (c *SmartChunkerImpl) createChunkFromSymbol(lines []string, sym analysis.Symbol, opts ChunkOptions) Chunk {
	startLine := max(1, sym.StartLine-opts.IncludeContext)
	endLine := min(len(lines), sym.EndLine+opts.IncludeContext)

	var contentBuilder strings.Builder
	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		contentBuilder.WriteString(lines[i])
		if i < endLine-1 {
			contentBuilder.WriteString("\n")
		}
	}

	content := contentBuilder.String()

	return Chunk{
		Content:   content,
		StartLine: startLine,
		EndLine:   endLine,
		Type:      string(sym.Kind),
		Name:      sym.Name,
		Tokens:    c.estimateTokens(content),
	}
}

func (c *SmartChunkerImpl) fillGaps(lines []string, chunks []Chunk, opts ChunkOptions) []Chunk {
	if len(chunks) == 0 {
		return chunks
	}

	result := make([]Chunk, 0, len(chunks)*2)
	prevEnd := 0

	for _, chunk := range chunks {
		if chunk.StartLine > prevEnd+1 {
			gapChunk := c.createGapChunk(lines, prevEnd+1, chunk.StartLine-1, opts)
			if gapChunk.Tokens >= minChunkTokens {
				result = append(result, gapChunk)
			}
		}
		result = append(result, chunk)
		prevEnd = chunk.EndLine
	}

	if prevEnd < len(lines) {
		gapChunk := c.createGapChunk(lines, prevEnd+1, len(lines), opts)
		if gapChunk.Tokens >= minChunkTokens {
			result = append(result, gapChunk)
		}
	}

	return result
}

func (c *SmartChunkerImpl) createGapChunk(lines []string, startLine, endLine int, _ ChunkOptions) Chunk {
	var contentBuilder strings.Builder
	for i := startLine - 1; i < endLine && i < len(lines); i++ {
		contentBuilder.WriteString(lines[i])
		if i < endLine-1 {
			contentBuilder.WriteString("\n")
		}
	}

	content := contentBuilder.String()

	return Chunk{
		Content:   content,
		StartLine: startLine,
		EndLine:   endLine,
		Type:      "block",
		Name:      "",
		Tokens:    c.estimateTokens(content),
	}
}

func (c *SmartChunkerImpl) splitLargeChunks(chunks []Chunk, opts ChunkOptions) []Chunk {
	result := make([]Chunk, 0, len(chunks))

	for _, chunk := range chunks {
		if chunk.Tokens <= opts.MaxTokensPerChunk {
			result = append(result, chunk)
			continue
		}

		subChunks := c.splitChunk(chunk, opts)
		result = append(result, subChunks...)
	}

	return result
}

func (c *SmartChunkerImpl) splitChunk(chunk Chunk, opts ChunkOptions) []Chunk {
	lines := strings.Split(chunk.Content, "\n")
	maxLinesPerChunk := (opts.MaxTokensPerChunk * charsPerToken) / 80 // assume ~80 chars per line

	if maxLinesPerChunk < 10 {
		maxLinesPerChunk = 10
	}

	result := make([]Chunk, 0)
	currentStart := 0

	for currentStart < len(lines) {
		currentEnd := min(currentStart+maxLinesPerChunk, len(lines))

		var contentBuilder strings.Builder
		for i := currentStart; i < currentEnd; i++ {
			contentBuilder.WriteString(lines[i])
			if i < currentEnd-1 {
				contentBuilder.WriteString("\n")
			}
		}

		content := contentBuilder.String()
		subChunk := Chunk{
			Content:   content,
			StartLine: chunk.StartLine + currentStart,
			EndLine:   chunk.StartLine + currentEnd - 1,
			Type:      chunk.Type,
			Name:      chunk.Name,
			Tokens:    c.estimateTokens(content),
		}

		result = append(result, subChunk)
		currentStart = currentEnd
	}

	return result
}

func (c *SmartChunkerImpl) chunkByLines(content []byte, _ string, opts ChunkOptions) []Chunk {
	lines := strings.Split(string(content), "\n")
	maxLinesPerChunk := (opts.MaxTokensPerChunk * charsPerToken) / 80

	if maxLinesPerChunk < 20 {
		maxLinesPerChunk = 20
	}

	chunks := make([]Chunk, 0)
	currentStart := 0

	for currentStart < len(lines) {
		currentEnd := min(currentStart+maxLinesPerChunk, len(lines))

		var contentBuilder strings.Builder
		for i := currentStart; i < currentEnd; i++ {
			contentBuilder.WriteString(lines[i])
			if i < currentEnd-1 {
				contentBuilder.WriteString("\n")
			}
		}

		chunkContent := contentBuilder.String()
		chunk := Chunk{
			Content:   chunkContent,
			StartLine: currentStart + 1,
			EndLine:   currentEnd,
			Type:      "block",
			Tokens:    c.estimateTokens(chunkContent),
		}

		chunks = append(chunks, chunk)
		currentStart = currentEnd
	}

	return chunks
}

func (c *SmartChunkerImpl) extractQueryTerms(query string) []string {
	words := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_')
	})

	terms := make([]string, 0, len(words))
	for _, word := range words {
		if len(word) >= 2 && !isStopWord(word) {
			terms = append(terms, word)
		}
	}

	return terms
}

func (c *SmartChunkerImpl) calculateRelevance(chunk Chunk, queryTerms []string) float64 {
	if len(queryTerms) == 0 {
		return 0.5
	}

	contentLower := strings.ToLower(chunk.Content)
	nameLower := strings.ToLower(chunk.Name)

	var score float64
	matchedTerms := 0

	for _, term := range queryTerms {
		if strings.Contains(nameLower, term) {
			score += 0.4
			matchedTerms++
		} else if strings.Contains(contentLower, term) {
			score += 0.2
			matchedTerms++
		}
	}

	if matchedTerms > 0 {
		score += float64(matchedTerms) / float64(len(queryTerms)) * 0.4
	}

	if chunk.Type == "function" || chunk.Type == "method" {
		score += 0.1
	}

	if score > 1.0 {
		score = 1.0
	}

	return score
}

func (c *SmartChunkerImpl) selectChunksWithinBudget(chunks []Chunk, maxTokens int) []Chunk {
	result := make([]Chunk, 0)
	totalTokens := 0

	for _, chunk := range chunks {
		if totalTokens+chunk.Tokens > maxTokens {
			continue
		}
		result = append(result, chunk)
		totalTokens += chunk.Tokens
	}

	return result
}

func (c *SmartChunkerImpl) estimateTokens(text string) int {
	return len(text) / charsPerToken
}

// GenerateFileSummary creates a summary for large files
func (c *SmartChunkerImpl) GenerateFileSummary(ctx context.Context, filePath string, content []byte) string {
	analyzer := c.analyzerRegistry.GetAnalyzer(filePath)
	if analyzer == nil {
		return c.generateLineSummary(content)
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return c.generateLineSummary(content)
	}

	return c.generateSymbolSummary(filePath, symbols)
}

func (c *SmartChunkerImpl) generateLineSummary(content []byte) string {
	lines := strings.Split(string(content), "\n")
	if len(lines) <= summaryMaxLines {
		return string(content)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("// File summary (%d lines total)\n", len(lines)))

	for i := 0; i < summaryMaxLines/2 && i < len(lines); i++ {
		sb.WriteString(lines[i])
		sb.WriteString("\n")
	}

	sb.WriteString("// ... truncated ...\n")

	startIdx := len(lines) - summaryMaxLines/2
	for i := startIdx; i < len(lines); i++ {
		sb.WriteString(lines[i])
		if i < len(lines)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func (c *SmartChunkerImpl) generateSymbolSummary(filePath string, symbols []analysis.Symbol) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("// File: %s\n", filepath.Base(filePath)))
	sb.WriteString("// Symbols:\n")

	for _, sym := range symbols {
		if sym.Parent != "" {
			continue
		}

		if sym.Signature != "" {
			sb.WriteString(fmt.Sprintf("//   %s\n", sym.Signature))
		} else {
			sb.WriteString(fmt.Sprintf("//   %s %s\n", sym.Kind, sym.Name))
		}
	}

	return sb.String()
}

func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"may": true, "might": true, "must": true, "shall": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "at": true, "by": true, "from": true, "as": true,
		"and": true, "or": true, "but": true, "if": true, "then": true,
		"this": true, "that": true, "these": true, "those": true,
		"it": true, "its": true, "i": true, "me": true, "my": true,
	}
	return stopWords[word]
}
