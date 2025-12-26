package semantic

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"syntaxia/domain"
)

// generateProjectID generates a unique project ID from path
func generateProjectID(projectRoot string) string {
	hash := sha256.Sum256([]byte(projectRoot))
	return hex.EncodeToString(hash[:8])
}

// shouldSkipDir checks if directory should be skipped during indexing
func shouldSkipDir(name string) bool {
	skipDirs := []string{
		".git", ".svn", ".hg",
		"node_modules", "vendor", "venv", ".venv",
		"build", "dist", "target", "out",
		".idea", ".vscode", ".vs",
		"__pycache__", ".pytest_cache",
		"coverage", ".nyc_output",
	}
	return slices.Contains(skipDirs, name)
}

// isCodeFile checks if file is a code file
func isCodeFile(path string) bool {
	codeExtensions := []string{
		".go", ".py", ".js", ".ts", ".jsx", ".tsx", ".vue",
		".java", ".kt", ".scala", ".rs", ".cpp", ".c", ".h",
		".cs", ".rb", ".php", ".swift", ".dart", ".lua",
		".sql", ".graphql", ".proto",
	}
	ext := strings.ToLower(filepath.Ext(path))
	return slices.Contains(codeExtensions, ext)
}

// createEmbeddingText creates text for embedding from chunk
func createEmbeddingText(chunk domain.CodeChunk) string {
	var sb strings.Builder
	if chunk.SymbolName != "" {
		sb.WriteString(chunk.SymbolName)
		sb.WriteString(": ")
	}
	sb.WriteString(chunk.Content)
	return sb.String()
}

// applyFilters applies search filters to results
func (s *ServiceImpl) applyFilters(results []domain.SemanticSearchResult, filters *domain.SearchFilters) []domain.SemanticSearchResult {
	if filters == nil {
		return results
	}

	filtered := make([]domain.SemanticSearchResult, 0, len(results))
	for _, r := range results {
		chunk := &r.Chunk
		if !matchesLanguageFilter(chunk, filters.Languages) ||
			!matchesChunkTypeFilter(chunk, filters.ChunkTypes) ||
			!matchesFilePathFilter(chunk, filters.FilePaths) ||
			isExcludedDir(chunk, filters.ExcludeDirs) {
			continue
		}
		filtered = append(filtered, r)
	}
	return filtered
}

func matchesLanguageFilter(chunk *domain.CodeChunk, languages []string) bool {
	if len(languages) == 0 {
		return true
	}
	return slices.Contains(languages, chunk.Language)
}

func matchesChunkTypeFilter(chunk *domain.CodeChunk, chunkTypes []domain.ChunkType) bool {
	if len(chunkTypes) == 0 {
		return true
	}
	return slices.Contains(chunkTypes, chunk.ChunkType)
}

func matchesFilePathFilter(chunk *domain.CodeChunk, filePaths []string) bool {
	if len(filePaths) == 0 {
		return true
	}
	for _, fp := range filePaths {
		if strings.HasPrefix(chunk.FilePath, fp) {
			return true
		}
	}
	return false
}

func isExcludedDir(chunk *domain.CodeChunk, excludeDirs []string) bool {
	for _, dir := range excludeDirs {
		if strings.Contains(chunk.FilePath, dir) {
			return true
		}
	}
	return false
}

// generateEmbeddingsWithRetry generates embeddings with retry logic
func (s *ServiceImpl) generateEmbeddingsWithRetry(ctx context.Context, texts []string) (*domain.EmbeddingResponse, error) {
	maxRetries := 3
	baseDelay := 1 * time.Second

	req := domain.EmbeddingRequest{Texts: texts}

	var lastErr error
	for attempt := range maxRetries {
		resp, err := s.embeddingProvider.GenerateEmbeddings(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if attempt < maxRetries-1 {
			delay := baseDelay * time.Duration(1<<attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}
