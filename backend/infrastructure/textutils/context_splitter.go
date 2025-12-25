package textutils

import (
	"fmt"
	"regexp"
	"syntaxia/domain"
	"strings"
	"unicode/utf8"
)

// This is a simplified token estimation. A more accurate one would use a proper tokenizer.
// For now, assume 1 token = 4 characters.
func approxTokens(s string) int {
	return utf8.RuneCountInString(s) / 4
}

// ContextSplitterImpl provides functionality to split a large text context.
type ContextSplitterImpl struct {
	log domain.Logger
}

// NewContextSplitter creates a new instance of ContextSplitterImpl.
func NewContextSplitter(log domain.Logger) domain.ContextSplitter {
	return &ContextSplitterImpl{log: log}
}

// SplitContext splits the given context text into chunks based on provided settings.
func (s *ContextSplitterImpl) SplitContext(ctxText string, settings domain.SplitSettings) ([]string, error) {
	if settings.MaxTokensPerChunk <= 0 {
		return nil, fmt.Errorf("MaxTokensPerChunk must be greater than 0")
	}

	totalTokens := approxTokens(ctxText)
	if totalTokens <= settings.MaxTokensPerChunk {
		s.log.Info(fmt.Sprintf("Context fits within single chunk (%d tokens). No splitting needed.", totalTokens))
		return []string{ctxText}, nil
	}

	s.log.Info(fmt.Sprintf("Splitting context (total %d tokens) with max %d tokens per chunk, strategy '%s', overlap %d.",
		totalTokens, settings.MaxTokensPerChunk, settings.SplitStrategy, settings.OverlapTokens))

	switch settings.SplitStrategy {
	case "file":
		return s.splitByFileHeaders(ctxText, settings.MaxTokensPerChunk)
	case "token":
		return s.splitByTokenCount(ctxText, settings.MaxTokensPerChunk, settings.OverlapTokens)
	case "smart": // Smart strategy tries file headers first, then falls back to token if a single file is too big
		chunks, err := s.splitByFileHeaders(ctxText, settings.MaxTokensPerChunk)
		if err != nil {
			return nil, err
		}
		// Check if any chunk is still too large after file-based splitting
		needsFurtherSplitting := false
		for _, chunk := range chunks {
			if approxTokens(chunk) > settings.MaxTokensPerChunk {
				needsFurtherSplitting = true
				break
			}
		}
		if needsFurtherSplitting {
			s.log.Info("Smart split: File-based splitting still resulted in large chunks. Falling back to token-based splitting.")
			return s.splitByTokenCount(ctxText, settings.MaxTokensPerChunk, settings.OverlapTokens)
		}
		return chunks, nil
	default:
		return nil, fmt.Errorf("unknown split strategy: %s", settings.SplitStrategy)
	}
}

// splitByFileHeaders attempts to split the context by "--- File: " headers, keeping files whole.
// If a single file content exceeds the token limit, it will still be in its own chunk.
func (s *ContextSplitterImpl) splitByFileHeaders(text string, tokenLimit int) ([]string, error) {
	// Regular expression to find "--- File: " headers
	re := regexp.MustCompile(`(?m)^--- File: .*? ---\s*`)
	idxs := re.FindAllStringIndex(text, -1)

	if len(idxs) == 0 {
		s.log.Warning("No file headers found for file-based splitting. Returning text as single chunk.")
		return []string{text}, nil
	}

	var chunks []string
	currentChunkBuilder := strings.Builder{}
	currentChunkTokens := 0

	// Handle optional manifest at the beginning
	firstFileHeaderIdx := idxs[0][0]
	if firstFileHeaderIdx > 0 {
		manifestPart := text[:firstFileHeaderIdx]
		currentChunkBuilder.WriteString(manifestPart)
		currentChunkTokens += approxTokens(manifestPart)
	}

	for i := 0; i < len(idxs); i++ {
		start := idxs[i][0]
		end := len(text)
		if i+1 < len(idxs) {
			end = idxs[i+1][0]
		}
		fileSegment := text[start:end]
		fileSegmentTokens := approxTokens(fileSegment)

		// If adding this file would exceed the limit, start a new chunk
		// unless it's the very first file in an empty currentChunk, or the file itself is larger than the limit
		// In the latter case, the file will become its own chunk.
		if currentChunkTokens > 0 && currentChunkTokens+fileSegmentTokens > tokenLimit {
			chunks = append(chunks, strings.TrimSpace(currentChunkBuilder.String()))
			currentChunkBuilder.Reset()
			currentChunkTokens = 0
		}
		currentChunkBuilder.WriteString(fileSegment)
		currentChunkTokens += fileSegmentTokens
	}

	if currentChunkBuilder.Len() > 0 {
		chunks = append(chunks, strings.TrimSpace(currentChunkBuilder.String()))
	}

	s.log.Info(fmt.Sprintf("File-based splitting resulted in %d chunks.", len(chunks)))
	return chunks, nil
}

// splitByTokenCount splits the text purely by token count, with optional overlap.
func (s *ContextSplitterImpl) splitByTokenCount(text string, tokenLimit, overlapTokens int) ([]string, error) {
	if tokenLimit <= overlapTokens {
		return nil, fmt.Errorf("tokenLimit (%d) must be greater than overlapTokens (%d)", tokenLimit, overlapTokens)
	}

	var chunks []string
	runes := []rune(text)
	textLength := len(runes)

	// Convert token limit to character limit (approximate)
	charLimit := tokenLimit * 4
	overlapChars := overlapTokens * 4

	currentPos := 0
	for currentPos < textLength {
		endPos := currentPos + charLimit
		if endPos > textLength {
			endPos = textLength
		}

		chunk := string(runes[currentPos:endPos])
		chunks = append(chunks, chunk)

		// If we've reached the end, break
		if endPos >= textLength {
			break
		}

		// Move start position for next chunk, considering overlap
		// Ensure we always make forward progress
		step := charLimit - overlapChars
		if step <= 0 {
			step = 1 // Minimum step to avoid infinite loop
		}
		currentPos += step
	}

	s.log.Info(fmt.Sprintf("Token-based splitting resulted in %d chunks.", len(chunks)))
	return chunks, nil
}
