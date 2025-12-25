package semantic

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// ServiceImpl implements SemanticSearchService
type ServiceImpl struct {
	embeddingProvider domain.EmbeddingProvider
	vectorStore       domain.VectorStore
	symbolIndex       analysis.SymbolIndex
	log               domain.Logger
	chunker           domain.CodeChunker

	// Indexing state
	indexingMu    sync.RWMutex
	indexingState map[string]*IndexingState
}

// IndexingState tracks the state of project indexing
type IndexingState struct {
	ProjectID    string
	InProgress   bool
	Progress     float64
	TotalFiles   int
	IndexedFiles int
	StartedAt    time.Time
	Error        error
}

// SymbolInfoForChunking represents symbol info for chunking
type SymbolInfoForChunking struct {
	Name      string
	Kind      string
	StartLine int
	EndLine   int
}

// NewService creates a new semantic search service
func NewService(
	embeddingProvider domain.EmbeddingProvider,
	vectorStore domain.VectorStore,
	symbolIndex analysis.SymbolIndex,
	log domain.Logger,
	chunker domain.CodeChunker,
) *ServiceImpl {
	return &ServiceImpl{
		embeddingProvider: embeddingProvider,
		vectorStore:       vectorStore,
		symbolIndex:       symbolIndex,
		log:               log,
		chunker:           chunker,
		indexingState:     make(map[string]*IndexingState),
	}
}

// startIndexingState initializes indexing state
func (s *ServiceImpl) startIndexingState(projectID string) (*IndexingState, error) {
	s.indexingMu.Lock()
	defer s.indexingMu.Unlock()

	if state, exists := s.indexingState[projectID]; exists && state.InProgress {
		return nil, fmt.Errorf("indexing already in progress for project")
	}

	state := &IndexingState{ProjectID: projectID, InProgress: true, StartedAt: time.Now()}
	s.indexingState[projectID] = state
	return state, nil
}

// Search performs semantic search
func (s *ServiceImpl) Search(ctx context.Context, req domain.SemanticSearchRequest) (*domain.SemanticSearchResponse, error) {
	startTime := time.Now()
	projectID := generateProjectID(req.ProjectRoot)

	if req.TopK == 0 {
		req.TopK = 10
	}
	if req.MinScore == 0 {
		req.MinScore = 0.5
	}

	s.log.Info(fmt.Sprintf("Semantic search: query='%s', topK=%d", domain.TruncateString(req.Query, 50), req.TopK))

	switch req.SearchType {
	case domain.SearchTypeKeyword:
		return s.keywordSearch(ctx, req, startTime)
	case domain.SearchTypeHybrid:
		return s.hybridSearch(ctx, req, startTime)
	default:
		return s.semanticSearch(ctx, projectID, req, startTime)
	}
}

// semanticSearch performs pure semantic search
func (s *ServiceImpl) semanticSearch(ctx context.Context, projectID string, req domain.SemanticSearchRequest, startTime time.Time) (*domain.SemanticSearchResponse, error) {
	resp, err := s.generateEmbeddingsWithRetry(ctx, []string{req.Query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	results, err := s.vectorStore.Search(ctx, projectID, resp.Embeddings[0], req.TopK*2, req.MinScore)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	results = s.applyFilters(results, req.Filters)
	if len(results) > req.TopK {
		results = results[:req.TopK]
	}

	return &domain.SemanticSearchResponse{
		Results:      results,
		TotalResults: len(results),
		QueryTime:    time.Since(startTime),
		SearchType:   domain.SearchTypeSemantic,
	}, nil
}

// keywordSearch performs keyword-based search
func (s *ServiceImpl) keywordSearch(_ context.Context, req domain.SemanticSearchRequest, startTime time.Time) (*domain.SemanticSearchResponse, error) {
	if s.symbolIndex == nil {
		return &domain.SemanticSearchResponse{
			Results:    []domain.SemanticSearchResult{},
			QueryTime:  time.Since(startTime),
			SearchType: domain.SearchTypeKeyword,
		}, nil
	}

	symbols := s.symbolIndex.SearchByName(req.Query)
	results := make([]domain.SemanticSearchResult, 0, len(symbols))

	for _, sym := range symbols {
		results = append(results, domain.SemanticSearchResult{
			Chunk: domain.CodeChunk{
				FilePath:   sym.FilePath,
				StartLine:  sym.StartLine,
				EndLine:    sym.EndLine,
				SymbolName: sym.Name,
				ChunkType:  domain.ChunkTypeFunction,
			},
			Score: 0.8,
		})
	}

	results = s.applyFilters(results, req.Filters)
	if len(results) > req.TopK {
		results = results[:req.TopK]
	}

	return &domain.SemanticSearchResponse{
		Results:      results,
		TotalResults: len(results),
		QueryTime:    time.Since(startTime),
		SearchType:   domain.SearchTypeKeyword,
	}, nil
}

// hybridSearch combines semantic and keyword search
func (s *ServiceImpl) hybridSearch(ctx context.Context, req domain.SemanticSearchRequest, startTime time.Time) (*domain.SemanticSearchResponse, error) {
	projectID := generateProjectID(req.ProjectRoot)

	// Get semantic results
	semanticResp, err := s.generateEmbeddingsWithRetry(ctx, []string{req.Query})
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	semanticResults, err := s.vectorStore.Search(ctx, projectID, semanticResp.Embeddings[0], req.TopK*2, req.MinScore*0.8)
	if err != nil {
		return nil, fmt.Errorf("failed to search vector store: %w", err)
	}

	// Get keyword results
	var keywordResults []domain.SemanticSearchResult
	if s.symbolIndex != nil {
		symbols := s.symbolIndex.SearchByName(req.Query)
		for _, sym := range symbols {
			keywordResults = append(keywordResults, domain.SemanticSearchResult{
				Chunk: domain.CodeChunk{
					FilePath:   sym.FilePath,
					StartLine:  sym.StartLine,
					EndLine:    sym.EndLine,
					SymbolName: sym.Name,
					ChunkType:  domain.ChunkTypeFunction,
				},
				Score: 0.7,
			})
		}
	}

	// Merge and deduplicate results
	resultMap := make(map[string]domain.SemanticSearchResult)
	for _, r := range semanticResults {
		key := fmt.Sprintf("%s:%d", r.Chunk.FilePath, r.Chunk.StartLine)
		resultMap[key] = r
	}
	for _, r := range keywordResults {
		key := fmt.Sprintf("%s:%d", r.Chunk.FilePath, r.Chunk.StartLine)
		if existing, ok := resultMap[key]; ok {
			existing.Score = (existing.Score + r.Score) / 2
			resultMap[key] = existing
		} else {
			resultMap[key] = r
		}
	}

	// Convert to slice and sort
	results := make([]domain.SemanticSearchResult, 0, len(resultMap))
	for _, r := range resultMap {
		results = append(results, r)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	results = s.applyFilters(results, req.Filters)
	if len(results) > req.TopK {
		results = results[:req.TopK]
	}

	return &domain.SemanticSearchResponse{
		Results:      results,
		TotalResults: len(results),
		QueryTime:    time.Since(startTime),
		SearchType:   domain.SearchTypeHybrid,
	}, nil
}

// FindSimilar finds similar code
func (s *ServiceImpl) FindSimilar(_ context.Context, req domain.SimilarCodeRequest) (*domain.SemanticSearchResponse, error) {
	startTime := time.Now()

	topK := req.TopK
	if topK == 0 {
		topK = 10
	}

	// This is a simplified implementation - full implementation would read file content
	return &domain.SemanticSearchResponse{
		Results:      []domain.SemanticSearchResult{},
		TotalResults: 0,
		QueryTime:    time.Since(startTime),
		SearchType:   domain.SearchTypeSemantic,
	}, nil
}

// GetClusters returns code clusters
func (s *ServiceImpl) GetClusters(_ context.Context, _ string, _ int) ([]domain.ClusterInfo, error) {
	return []domain.ClusterInfo{}, nil
}

// GetStats returns indexing statistics
func (s *ServiceImpl) GetStats(ctx context.Context, projectRoot string) (*domain.VectorStoreStats, error) {
	projectID := generateProjectID(projectRoot)
	return s.vectorStore.GetStats(ctx, projectID)
}

// IsIndexed checks if a project is indexed
func (s *ServiceImpl) IsIndexed(ctx context.Context, projectRoot string) bool {
	projectID := generateProjectID(projectRoot)
	stats, err := s.vectorStore.GetStats(ctx, projectID)
	if err != nil || stats == nil {
		return false
	}
	return stats.TotalChunks > 0
}

// InvalidateFile marks a file for re-indexing
func (s *ServiceImpl) InvalidateFile(ctx context.Context, projectRoot string, filePath string) error {
	projectID := generateProjectID(projectRoot)
	return s.vectorStore.Delete(ctx, projectID, filePath)
}

// GetIndexingState returns the current indexing state
func (s *ServiceImpl) GetIndexingState(projectRoot string) *IndexingState {
	projectID := generateProjectID(projectRoot)
	s.indexingMu.RLock()
	defer s.indexingMu.RUnlock()
	return s.indexingState[projectID]
}
