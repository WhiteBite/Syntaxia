package rag

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"syntaxia/domain"
)

// Service implements RAGService for Retrieval Augmented Generation
type Service struct {
	semanticSearch    domain.SemanticSearchService
	embeddingProvider domain.EmbeddingProvider
	log               domain.Logger
}

// NewService creates a new RAG service
func NewService(
	semanticSearch domain.SemanticSearchService,
	embeddingProvider domain.EmbeddingProvider,
	log domain.Logger,
) *Service {
	return &Service{
		semanticSearch:    semanticSearch,
		embeddingProvider: embeddingProvider,
		log:               log,
	}
}

// RetrieveContext retrieves relevant context for a query
func (r *Service) RetrieveContext(ctx context.Context, query string, projectRoot string, maxTokens int) ([]domain.CodeChunk, error) {
	r.log.Info(fmt.Sprintf("RAG: Retrieving context for query, maxTokens=%d", maxTokens))

	// Use hybrid search for best results
	searchReq := domain.SemanticSearchRequest{
		Query:       query,
		ProjectRoot: projectRoot,
		TopK:        50, // Get more results for filtering
		MinScore:    0.3,
		SearchType:  domain.SearchTypeHybrid,
	}

	results, err := r.semanticSearch.Search(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	// Select chunks that fit within token budget
	selectedChunks := make([]domain.CodeChunk, 0, len(results.Results))
	totalTokens := 0

	for _, result := range results.Results {
		chunkTokens := result.Chunk.TokenCount
		if chunkTokens == 0 {
			chunkTokens = len(result.Chunk.Content) / 4
		}

		if totalTokens+chunkTokens > maxTokens {
			// Try to fit smaller chunks
			continue
		}

		selectedChunks = append(selectedChunks, result.Chunk)
		totalTokens += chunkTokens
	}

	r.log.Info(fmt.Sprintf("RAG: Selected %d chunks with %d tokens", len(selectedChunks), totalTokens))

	// Sort by file path and line number for better readability
	sort.Slice(selectedChunks, func(i, j int) bool {
		if selectedChunks[i].FilePath != selectedChunks[j].FilePath {
			return selectedChunks[i].FilePath < selectedChunks[j].FilePath
		}
		return selectedChunks[i].StartLine < selectedChunks[j].StartLine
	})

	return selectedChunks, nil
}

// RetrieveAndRank retrieves and re-ranks results using cross-encoder style scoring
func (r *Service) RetrieveAndRank(ctx context.Context, query string, projectRoot string, topK int) ([]domain.SemanticSearchResult, error) {
	r.log.Info(fmt.Sprintf("RAG: Retrieve and rank, topK=%d", topK))

	// First stage: retrieve candidates with semantic search
	searchReq := domain.SemanticSearchRequest{
		Query:       query,
		ProjectRoot: projectRoot,
		TopK:        topK * 3, // Get more candidates for re-ranking
		MinScore:    0.2,
		SearchType:  domain.SearchTypeSemantic,
	}

	results, err := r.semanticSearch.Search(ctx, searchReq)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	if len(results.Results) == 0 {
		return []domain.SemanticSearchResult{}, nil
	}

	// Second stage: re-rank using additional signals
	rerankedResults := r.rerank(query, results.Results)

	// Return top K
	if len(rerankedResults) > topK {
		rerankedResults = rerankedResults[:topK]
	}

	return rerankedResults, nil
}

// calcSymbolBoost calculates boost for symbol name matches
func calcSymbolBoost(symbolName string, queryTerms []string) float32 {
	if symbolName == "" {
		return 0
	}
	symbolLower := strings.ToLower(symbolName)
	var boost float32
	for _, term := range queryTerms {
		if strings.Contains(symbolLower, term) {
			boost += 0.15
		}
	}
	return boost
}

// calcContentBoost calculates boost for content matches
func calcContentBoost(content string, queryTerms []string) float32 {
	if len(queryTerms) == 0 {
		return 0
	}
	contentLower := strings.ToLower(content)
	matchCount := 0
	for _, term := range queryTerms {
		if strings.Contains(contentLower, term) {
			matchCount++
		}
	}
	return float32(matchCount) / float32(len(queryTerms)) * 0.1
}

// calcChunkTypeBoost returns boost based on chunk type
func calcChunkTypeBoost(chunkType domain.ChunkType) float32 {
	if chunkType == domain.ChunkTypeFunction || chunkType == domain.ChunkTypeMethod {
		return 0.05
	}
	return 0
}

// rerank re-ranks results using multiple signals
func (r *Service) rerank(query string, results []domain.SemanticSearchResult) []domain.SemanticSearchResult {
	queryTerms := strings.Fields(strings.ToLower(query))

	for i := range results {
		chunk := &results[i].Chunk
		score := results[i].Score
		score += calcSymbolBoost(chunk.SymbolName, queryTerms)
		score += calcContentBoost(chunk.Content, queryTerms)
		score += calcChunkTypeBoost(chunk.ChunkType)
		if chunk.TokenCount > 400 {
			score -= 0.02
		}
		results[i].Score = score
	}

	sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
	return results
}

// HybridSearch performs hybrid keyword + semantic search with advanced merging
func (r *Service) HybridSearch(ctx context.Context, req domain.SemanticSearchRequest) (*domain.SemanticSearchResponse, error) {
	startTime := time.Now()

	r.log.Info(fmt.Sprintf("RAG: Hybrid search for '%s'", domain.TruncateString(req.Query, 50)))

	// Perform semantic search
	semanticReq := req
	semanticReq.SearchType = domain.SearchTypeSemantic
	semanticReq.TopK = req.TopK * 2

	semanticResults, err := r.semanticSearch.Search(ctx, semanticReq)
	if err != nil {
		return nil, fmt.Errorf("semantic search failed: %w", err)
	}

	// Perform keyword search
	keywordReq := req
	keywordReq.SearchType = domain.SearchTypeKeyword
	keywordReq.TopK = req.TopK

	keywordResults, err := r.semanticSearch.Search(ctx, keywordReq)
	if err != nil {
		// Keyword search failure is not critical
		r.log.Warning(fmt.Sprintf("Keyword search failed: %v", err))
		keywordResults = &domain.SemanticSearchResponse{Results: []domain.SemanticSearchResult{}}
	}

	// Reciprocal Rank Fusion (RRF) for merging results
	k := 60 // RRF constant
	scores := make(map[string]float32)
	chunks := make(map[string]domain.CodeChunk)

	// Add semantic results
	for rank, result := range semanticResults.Results {
		key := chunkKey(result.Chunk)
		scores[key] += 1.0 / float32(k+rank+1)
		chunks[key] = result.Chunk
	}

	// Add keyword results
	for rank, result := range keywordResults.Results {
		key := chunkKey(result.Chunk)
		scores[key] += 1.0 / float32(k+rank+1)
		if _, exists := chunks[key]; !exists {
			chunks[key] = result.Chunk
		}
	}

	// Convert to results
	mergedResults := make([]domain.SemanticSearchResult, 0, len(scores))
	for key, score := range scores {
		mergedResults = append(mergedResults, domain.SemanticSearchResult{
			Chunk: chunks[key],
			Score: score,
		})
	}

	// Sort by RRF score
	sort.Slice(mergedResults, func(i, j int) bool {
		return mergedResults[i].Score > mergedResults[j].Score
	})

	// Apply re-ranking
	mergedResults = r.rerank(req.Query, mergedResults)

	// Limit to topK
	if len(mergedResults) > req.TopK {
		mergedResults = mergedResults[:req.TopK]
	}

	return &domain.SemanticSearchResponse{
		Results:      mergedResults,
		TotalResults: len(mergedResults),
		QueryTime:    time.Since(startTime),
		SearchType:   domain.SearchTypeHybrid,
	}, nil
}

// BuildContextPrompt builds a context prompt from retrieved chunks
func (r *Service) BuildContextPrompt(chunks []domain.CodeChunk) string {
	if len(chunks) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## Relevant Code Context\n\n")

	currentFile := ""
	for _, chunk := range chunks {
		if chunk.FilePath != currentFile {
			if currentFile != "" {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("### File: %s\n", chunk.FilePath))
			currentFile = chunk.FilePath
		}

		if chunk.SymbolName != "" {
			sb.WriteString(fmt.Sprintf("\n#### %s %s (lines %d-%d)\n",
				chunk.SymbolKind, chunk.SymbolName, chunk.StartLine, chunk.EndLine))
		} else {
			sb.WriteString(fmt.Sprintf("\n#### Lines %d-%d\n", chunk.StartLine, chunk.EndLine))
		}

		sb.WriteString("```")
		sb.WriteString(chunk.Language)
		sb.WriteString("\n")
		sb.WriteString(chunk.Content)
		sb.WriteString("\n```\n")
	}

	return sb.String()
}

// Helper functions

func chunkKey(chunk domain.CodeChunk) string {
	return fmt.Sprintf("%s:%d:%d", chunk.FilePath, chunk.StartLine, chunk.EndLine)
}
