package embeddings

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"unicode/utf8"
)

// ChunkerConfig configuration for code chunking
type ChunkerConfig struct {
	MaxChunkTokens int  `json:"maxChunkTokens"`
	MinChunkTokens int  `json:"minChunkTokens"`
	OverlapTokens  int  `json:"overlapTokens"`
	PreferSymbols  bool `json:"preferSymbols"`  // prefer function/class boundaries
	IncludeContext bool `json:"includeContext"` // include surrounding context
}

// DefaultChunkerConfig returns default chunking configuration
func DefaultChunkerConfig() ChunkerConfig {
	return ChunkerConfig{
		MaxChunkTokens: 512,
		MinChunkTokens: 50,
		OverlapTokens:  50,
		PreferSymbols:  true,
		IncludeContext: true,
	}
}

// CodeChunker splits code into chunks for embedding
type CodeChunker struct {
	config ChunkerConfig
}

// NewCodeChunker creates a new code chunker
func NewCodeChunker(config ChunkerConfig) *CodeChunker {
	return &CodeChunker{config: config}
}

// ChunkFile splits a file into chunks
func (c *CodeChunker) ChunkFile(filePath string, content []byte, symbols []SymbolInfo) []domain.CodeChunk {
	language := detectLanguage(filePath)
	lines := strings.Split(string(content), "\n")

	// If we have symbols and prefer symbol-based chunking
	if c.config.PreferSymbols && len(symbols) > 0 {
		return c.chunkBySymbols(filePath, lines, symbols, language)
	}

	// Fall back to fixed-size chunking
	return c.chunkBySize(filePath, lines, language)
}

// SymbolInfo represents a symbol for chunking
type SymbolInfo struct {
	Name      string
	Kind      string
	StartLine int
	EndLine   int
}

// chunkBySymbols creates chunks based on symbol boundaries
func (c *CodeChunker) chunkBySymbols(filePath string, lines []string, symbols []SymbolInfo, language string) []domain.CodeChunk {
	var chunks []domain.CodeChunk
	usedLines := make(map[int]bool)

	// First, create chunks for each symbol
	for _, sym := range symbols {
		if sym.StartLine < 1 || sym.EndLine > len(lines) {
			continue
		}

		// Get symbol content
		symbolLines := lines[sym.StartLine-1 : sym.EndLine]
		content := strings.Join(symbolLines, "\n")
		tokenCount := estimateTokens(content)

		// Skip if too small
		if tokenCount < c.config.MinChunkTokens {
			continue
		}

		// If too large, split into sub-chunks
		if tokenCount > c.config.MaxChunkTokens {
			subChunks := c.splitLargeSymbol(filePath, symbolLines, sym, language)
			chunks = append(chunks, subChunks...)
		} else {
			chunk := domain.CodeChunk{
				ID:         generateChunkID(filePath, sym.StartLine, sym.EndLine),
				FilePath:   filePath,
				Content:    content,
				StartLine:  sym.StartLine,
				EndLine:    sym.EndLine,
				ChunkType:  mapSymbolKindToChunkType(sym.Kind),
				SymbolName: sym.Name,
				SymbolKind: sym.Kind,
				Language:   language,
				TokenCount: tokenCount,
				Hash:       hashContent(content),
			}
			chunks = append(chunks, chunk)
		}

		// Mark lines as used
		for i := sym.StartLine; i <= sym.EndLine; i++ {
			usedLines[i] = true
		}
	}

	// Create chunks for remaining code (imports, constants, etc.)
	chunks = append(chunks, c.chunkRemainingLines(filePath, lines, usedLines, language)...)

	return chunks
}

// splitLargeSymbol splits a large symbol into smaller chunks
func (c *CodeChunker) splitLargeSymbol(filePath string, lines []string, sym SymbolInfo, language string) []domain.CodeChunk {
	chunks := make([]domain.CodeChunk, 0, len(lines)/50+1)

	currentStart := 0
	currentTokens := 0
	currentLines := make([]string, 0, 50)

	for i, line := range lines {
		lineTokens := estimateTokens(line)

		if currentTokens+lineTokens > c.config.MaxChunkTokens && len(currentLines) > 0 {
			// Create chunk
			content := strings.Join(currentLines, "\n")
			chunk := domain.CodeChunk{
				ID:         generateChunkID(filePath, sym.StartLine+currentStart, sym.StartLine+i-1),
				FilePath:   filePath,
				Content:    content,
				StartLine:  sym.StartLine + currentStart,
				EndLine:    sym.StartLine + i - 1,
				ChunkType:  domain.ChunkTypeBlock,
				SymbolName: sym.Name,
				SymbolKind: sym.Kind,
				Language:   language,
				TokenCount: currentTokens,
				Hash:       hashContent(content),
			}
			chunks = append(chunks, chunk)

			// Start new chunk with overlap
			overlapLines := getOverlapLines(currentLines, c.config.OverlapTokens)
			currentLines = overlapLines
			currentStart = i - len(overlapLines)
			currentTokens = estimateTokens(strings.Join(overlapLines, "\n"))
		}

		currentLines = append(currentLines, line)
		currentTokens += lineTokens
	}

	// Add remaining lines
	if len(currentLines) > 0 {
		content := strings.Join(currentLines, "\n")
		chunk := domain.CodeChunk{
			ID:         generateChunkID(filePath, sym.StartLine+currentStart, sym.EndLine),
			FilePath:   filePath,
			Content:    content,
			StartLine:  sym.StartLine + currentStart,
			EndLine:    sym.EndLine,
			ChunkType:  domain.ChunkTypeBlock,
			SymbolName: sym.Name,
			SymbolKind: sym.Kind,
			Language:   language,
			TokenCount: currentTokens,
			Hash:       hashContent(content),
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}

// chunkBySize creates fixed-size chunks
func (c *CodeChunker) chunkBySize(filePath string, lines []string, language string) []domain.CodeChunk {
	chunks := make([]domain.CodeChunk, 0, len(lines)/50+1)

	currentStart := 1
	currentTokens := 0
	currentLines := make([]string, 0, 50)

	for i, line := range lines {
		lineTokens := estimateTokens(line)

		if currentTokens+lineTokens > c.config.MaxChunkTokens && len(currentLines) > 0 {
			content := strings.Join(currentLines, "\n")
			// EndLine should be currentStart + len(currentLines) - 1 (1-based, inclusive)
			endLine := currentStart + len(currentLines) - 1
			chunk := domain.CodeChunk{
				ID:         generateChunkID(filePath, currentStart, endLine),
				FilePath:   filePath,
				Content:    content,
				StartLine:  currentStart,
				EndLine:    endLine,
				ChunkType:  domain.ChunkTypeBlock,
				Language:   language,
				TokenCount: currentTokens,
				Hash:       hashContent(content),
			}
			chunks = append(chunks, chunk)

			// Start new chunk with overlap
			overlapLines := getOverlapLines(currentLines, c.config.OverlapTokens)
			currentLines = overlapLines
			currentStart = i + 2 - len(overlapLines) // i+1 is next line (0-based), +1 for 1-based
			currentTokens = estimateTokens(strings.Join(overlapLines, "\n"))
		}

		currentLines = append(currentLines, line)
		currentTokens += lineTokens
	}

	// Add remaining lines
	if len(currentLines) > 0 && currentTokens >= c.config.MinChunkTokens {
		content := strings.Join(currentLines, "\n")
		endLine := currentStart + len(currentLines) - 1
		if endLine > len(lines) {
			endLine = len(lines)
		}
		chunk := domain.CodeChunk{
			ID:         generateChunkID(filePath, currentStart, endLine),
			FilePath:   filePath,
			Content:    content,
			StartLine:  currentStart,
			EndLine:    endLine,
			ChunkType:  domain.ChunkTypeBlock,
			Language:   language,
			TokenCount: currentTokens,
			Hash:       hashContent(content),
		}
		chunks = append(chunks, chunk)
	}

	return chunks
}

// chunkRemainingLines creates chunks for lines not covered by symbols
func (c *CodeChunker) chunkRemainingLines(filePath string, lines []string, usedLines map[int]bool, language string) []domain.CodeChunk {
	chunks := make([]domain.CodeChunk, 0, len(lines)/50+1)
	currentLines := make([]string, 0, 50)
	currentStart := 0

	for i, line := range lines {
		lineNum := i + 1
		if usedLines[lineNum] {
			// Flush current chunk if any
			if len(currentLines) > 0 {
				content := strings.Join(currentLines, "\n")
				tokenCount := estimateTokens(content)
				if tokenCount >= c.config.MinChunkTokens {
					chunk := domain.CodeChunk{
						ID:         generateChunkID(filePath, currentStart, i),
						FilePath:   filePath,
						Content:    content,
						StartLine:  currentStart,
						EndLine:    i,
						ChunkType:  domain.ChunkTypeBlock,
						Language:   language,
						TokenCount: tokenCount,
						Hash:       hashContent(content),
					}
					chunks = append(chunks, chunk)
				}
				currentLines = nil
			}
			continue
		}

		if len(currentLines) == 0 {
			currentStart = lineNum
		}
		currentLines = append(currentLines, line)
	}

	// Flush remaining
	if len(currentLines) > 0 {
		content := strings.Join(currentLines, "\n")
		tokenCount := estimateTokens(content)
		if tokenCount >= c.config.MinChunkTokens {
			chunk := domain.CodeChunk{
				ID:         generateChunkID(filePath, currentStart, len(lines)),
				FilePath:   filePath,
				Content:    content,
				StartLine:  currentStart,
				EndLine:    len(lines),
				ChunkType:  domain.ChunkTypeBlock,
				Language:   language,
				TokenCount: tokenCount,
				Hash:       hashContent(content),
			}
			chunks = append(chunks, chunk)
		}
	}

	return chunks
}

// Helper functions

func generateChunkID(filePath string, startLine, endLine int) string {
	data := fmt.Sprintf("%s:%d:%d", filePath, startLine, endLine)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

func hashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:16])
}

func estimateTokens(text string) int {
	// Rough estimation: ~4 characters per token for code
	return utf8.RuneCountInString(text) / 4
}

func getOverlapLines(lines []string, overlapTokens int) []string {
	if len(lines) == 0 {
		return nil
	}

	tokens := 0
	startIdx := len(lines)

	for i := len(lines) - 1; i >= 0; i-- {
		lineTokens := estimateTokens(lines[i])
		if tokens+lineTokens > overlapTokens {
			break
		}
		tokens += lineTokens
		startIdx = i
	}

	if startIdx >= len(lines) {
		return nil
	}

	result := make([]string, len(lines)-startIdx)
	copy(result, lines[startIdx:])
	return result
}

var extToLanguage = map[string]string{
	".go": "go", ".ts": "typescript", ".tsx": "typescript",
	".js": "javascript", ".jsx": "javascript", ".mjs": "javascript",
	".vue": "vue", ".py": "python", ".java": "java",
	".kt": "kotlin", ".kts": "kotlin", ".dart": "dart",
	".rs": "rust", ".cs": "csharp", ".cpp": "cpp", ".cc": "cpp", ".cxx": "cpp", ".hpp": "cpp",
	".c": "c", ".h": "c", ".rb": "ruby", ".php": "php", ".swift": "swift",
}

func detectLanguage(filePath string) string {
	if lang, ok := extToLanguage[strings.ToLower(filepath.Ext(filePath))]; ok {
		return lang
	}
	return "unknown"
}

func mapSymbolKindToChunkType(kind string) domain.ChunkType {
	switch strings.ToLower(kind) {
	case "function", "func":
		return domain.ChunkTypeFunction
	case "class":
		return domain.ChunkTypeClass
	case "method":
		return domain.ChunkTypeMethod
	default:
		return domain.ChunkTypeBlock
	}
}
