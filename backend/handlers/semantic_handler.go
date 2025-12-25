package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
)

// SemanticHandler handles semantic search API requests
type SemanticHandler struct {
	semanticSearch domain.SemanticSearchService
	ragService     domain.RAGService
	log            domain.Logger
}

// NewSemanticHandler creates a new semantic handler
func NewSemanticHandler(
	semanticSearch domain.SemanticSearchService,
	ragService domain.RAGService,
	log domain.Logger,
) *SemanticHandler {
	return &SemanticHandler{
		semanticSearch: semanticSearch,
		ragService:     ragService,
		log:            log,
	}
}

// SemanticSearchRequest represents a search request from frontend
type SemanticSearchRequest struct {
	Query       string   `json:"query"`
	ProjectRoot string   `json:"projectRoot"`
	TopK        int      `json:"topK"`
	MinScore    float32  `json:"minScore"`
	SearchType  string   `json:"searchType"`
	Languages   []string `json:"languages,omitempty"`
	ChunkTypes  []string `json:"chunkTypes,omitempty"`
}

// Search performs semantic search
func (h *SemanticHandler) Search(ctx context.Context, requestJSON string) (string, error) {
	var req SemanticSearchRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	// Set defaults
	if req.TopK == 0 {
		req.TopK = 10
	}
	if req.MinScore == 0 {
		req.MinScore = 0.5
	}

	searchType := domain.SearchTypeHybrid
	switch req.SearchType {
	case "semantic":
		searchType = domain.SearchTypeSemantic
	case "keyword":
		searchType = domain.SearchTypeKeyword
	}

	searchReq := domain.SemanticSearchRequest{
		Query:       req.Query,
		ProjectRoot: req.ProjectRoot,
		TopK:        req.TopK,
		MinScore:    req.MinScore,
		SearchType:  searchType,
	}

	// Add filters if provided
	if len(req.Languages) > 0 || len(req.ChunkTypes) > 0 {
		searchReq.Filters = &domain.SearchFilters{}
		if len(req.Languages) > 0 {
			searchReq.Filters.Languages = req.Languages
		}
		if len(req.ChunkTypes) > 0 {
			for _, ct := range req.ChunkTypes {
				searchReq.Filters.ChunkTypes = append(searchReq.Filters.ChunkTypes, domain.ChunkType(ct))
			}
		}
	}

	results, err := h.semanticSearch.Search(ctx, searchReq)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(resultJSON), nil
}

// FindSimilarRequest represents a find similar request
type FindSimilarRequest struct {
	FilePath    string  `json:"filePath"`
	StartLine   int     `json:"startLine"`
	EndLine     int     `json:"endLine"`
	TopK        int     `json:"topK"`
	MinScore    float32 `json:"minScore"`
	ExcludeSelf bool    `json:"excludeSelf"`
}

// FindSimilar finds similar code
func (h *SemanticHandler) FindSimilar(ctx context.Context, requestJSON string) (string, error) {
	var req FindSimilarRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	if req.TopK == 0 {
		req.TopK = 5
	}
	if req.MinScore == 0 {
		req.MinScore = 0.5
	}

	searchReq := domain.SimilarCodeRequest{
		FilePath:    req.FilePath,
		StartLine:   req.StartLine,
		EndLine:     req.EndLine,
		TopK:        req.TopK,
		MinScore:    req.MinScore,
		ExcludeSelf: req.ExcludeSelf,
	}

	results, err := h.semanticSearch.FindSimilar(ctx, searchReq)
	if err != nil {
		return "", fmt.Errorf("find similar failed: %w", err)
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(resultJSON), nil
}

// IndexProject indexes a project for semantic search
func (h *SemanticHandler) IndexProject(ctx context.Context, projectRoot string) error {
	h.log.Info(fmt.Sprintf("Starting semantic indexing for: %s", projectRoot))
	return h.semanticSearch.IndexProject(ctx, projectRoot)
}

// IndexFile indexes a single file
func (h *SemanticHandler) IndexFile(ctx context.Context, projectRoot, filePath string) error {
	return h.semanticSearch.IndexFile(ctx, projectRoot, filePath)
}

// GetStats returns indexing statistics
func (h *SemanticHandler) GetStats(ctx context.Context, projectRoot string) (string, error) {
	stats, err := h.semanticSearch.GetStats(ctx, projectRoot)
	if err != nil {
		return "", err
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("failed to marshal stats: %w", err)
	}

	return string(statsJSON), nil
}

// IsIndexed checks if a project is indexed
func (h *SemanticHandler) IsIndexed(ctx context.Context, projectRoot string) bool {
	return h.semanticSearch.IsIndexed(ctx, projectRoot)
}

// RetrieveContextRequest represents a RAG context request
type RetrieveContextRequest struct {
	Query       string `json:"query"`
	ProjectRoot string `json:"projectRoot"`
	MaxTokens   int    `json:"maxTokens"`
}

// RetrieveContext retrieves relevant context using RAG
func (h *SemanticHandler) RetrieveContext(ctx context.Context, requestJSON string) (string, error) {
	var req RetrieveContextRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	if req.MaxTokens == 0 {
		req.MaxTokens = 4000
	}

	chunks, err := h.ragService.RetrieveContext(ctx, req.Query, req.ProjectRoot, req.MaxTokens)
	if err != nil {
		return "", fmt.Errorf("retrieve context failed: %w", err)
	}

	chunksJSON, err := json.Marshal(chunks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal chunks: %w", err)
	}

	return string(chunksJSON), nil
}

// HybridSearch performs hybrid search using RAG
func (h *SemanticHandler) HybridSearch(ctx context.Context, requestJSON string) (string, error) {
	var req SemanticSearchRequest
	if err := json.Unmarshal([]byte(requestJSON), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	if req.TopK == 0 {
		req.TopK = 10
	}
	if req.MinScore == 0 {
		req.MinScore = 0.3
	}

	searchReq := domain.SemanticSearchRequest{
		Query:       req.Query,
		ProjectRoot: req.ProjectRoot,
		TopK:        req.TopK,
		MinScore:    req.MinScore,
		SearchType:  domain.SearchTypeHybrid,
	}

	results, err := h.ragService.HybridSearch(ctx, searchReq)
	if err != nil {
		return "", fmt.Errorf("hybrid search failed: %w", err)
	}

	resultJSON, err := json.Marshal(results)
	if err != nil {
		return "", fmt.Errorf("failed to marshal results: %w", err)
	}

	return string(resultJSON), nil
}
