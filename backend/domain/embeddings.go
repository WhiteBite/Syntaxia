package domain

import (
	"context"
	"time"
)

// EmbeddingVector represents a vector embedding
type EmbeddingVector []float32

// EmbeddingModel represents available embedding models
type EmbeddingModel string

const (
	EmbeddingModelOpenAI   EmbeddingModel = "text-embedding-ada-002"
	EmbeddingModelOpenAI3S EmbeddingModel = "text-embedding-3-small"
	EmbeddingModelOpenAI3L EmbeddingModel = "text-embedding-3-large"
	EmbeddingModelLocal    EmbeddingModel = "all-MiniLM-L6-v2"
	EmbeddingModelCodeBERT EmbeddingModel = "codebert-base"
)

// EmbeddingDimensions returns the dimension size for each model
func (m EmbeddingModel) Dimensions() int {
	switch m {
	case EmbeddingModelOpenAI:
		return 1536
	case EmbeddingModelOpenAI3S:
		return 1536
	case EmbeddingModelOpenAI3L:
		return 3072
	case EmbeddingModelLocal:
		return 384
	case EmbeddingModelCodeBERT:
		return 768
	default:
		return 1536
	}
}

// EmbeddingRequest represents a request to generate embeddings
type EmbeddingRequest struct {
	Texts []string       `json:"texts"`
	Model EmbeddingModel `json:"model"`
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	Embeddings []EmbeddingVector `json:"embeddings"`
	Model      EmbeddingModel    `json:"model"`
	TokensUsed int               `json:"tokensUsed"`
}

// CodeChunk represents a chunk of code for embedding
type CodeChunk struct {
	ID         string    `json:"id"`
	FilePath   string    `json:"filePath"`
	Content    string    `json:"content"`
	StartLine  int       `json:"startLine"`
	EndLine    int       `json:"endLine"`
	ChunkType  ChunkType `json:"chunkType"`
	SymbolName string    `json:"symbolName,omitempty"`
	SymbolKind string    `json:"symbolKind,omitempty"`
	Language   string    `json:"language"`
	TokenCount int       `json:"tokenCount"`
	Hash       string    `json:"hash"` // for change detection
}

// ChunkType represents the type of code chunk
type ChunkType string

const (
	ChunkTypeFile     ChunkType = "file"
	ChunkTypeFunction ChunkType = "function"
	ChunkTypeClass    ChunkType = "class"
	ChunkTypeMethod   ChunkType = "method"
	ChunkTypeBlock    ChunkType = "block"
)

// EmbeddedChunk represents a code chunk with its embedding
type EmbeddedChunk struct {
	Chunk     CodeChunk       `json:"chunk"`
	Embedding EmbeddingVector `json:"embedding"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

// SemanticSearchRequest represents a semantic search query
type SemanticSearchRequest struct {
	Query       string         `json:"query"`
	ProjectRoot string         `json:"projectRoot"`
	TopK        int            `json:"topK"`
	MinScore    float32        `json:"minScore"`
	Filters     *SearchFilters `json:"filters,omitempty"`
	SearchType  SearchType     `json:"searchType"`
}

// SearchFilters for filtering search results
type SearchFilters struct {
	Languages   []string    `json:"languages,omitempty"`
	ChunkTypes  []ChunkType `json:"chunkTypes,omitempty"`
	FilePaths   []string    `json:"filePaths,omitempty"`
	ExcludeDirs []string    `json:"excludeDirs,omitempty"`
}

// SearchType represents the type of search
type SearchType string

const (
	SearchTypeSemantic SearchType = "semantic"
	SearchTypeKeyword  SearchType = "keyword"
	SearchTypeHybrid   SearchType = "hybrid"
)

// SemanticSearchResult represents a single search result
type SemanticSearchResult struct {
	Chunk      CodeChunk `json:"chunk"`
	Score      float32   `json:"score"`
	Highlights []string  `json:"highlights,omitempty"`
	Reason     string    `json:"reason,omitempty"`
}

// SemanticSearchResponse represents the search response
type SemanticSearchResponse struct {
	Results      []SemanticSearchResult `json:"results"`
	TotalResults int                    `json:"totalResults"`
	QueryTime    time.Duration          `json:"queryTime"`
	SearchType   SearchType             `json:"searchType"`
}

// SimilarCodeRequest represents a request to find similar code
type SimilarCodeRequest struct {
	FilePath    string  `json:"filePath"`
	StartLine   int     `json:"startLine"`
	EndLine     int     `json:"endLine"`
	TopK        int     `json:"topK"`
	MinScore    float32 `json:"minScore"`
	ExcludeSelf bool    `json:"excludeSelf"`
}

// ClusterInfo represents a cluster of similar code
type ClusterInfo struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Chunks      []CodeChunk     `json:"chunks"`
	Centroid    EmbeddingVector `json:"centroid"`
	Size        int             `json:"size"`
}

// EmbeddingProvider generates embeddings for text
type EmbeddingProvider interface {
	// GenerateEmbeddings generates embeddings for the given texts
	GenerateEmbeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error)

	// GetModelInfo returns information about the embedding model
	GetModelInfo() EmbeddingModelInfo

	// ValidateRequest validates the embedding request
	ValidateRequest(req EmbeddingRequest) error
}

// EmbeddingModelInfo contains information about an embedding model
type EmbeddingModelInfo struct {
	Model      EmbeddingModel `json:"model"`
	Dimensions int            `json:"dimensions"`
	MaxTokens  int            `json:"maxTokens"`
	Provider   string         `json:"provider"`
}

// VectorStore stores and retrieves embeddings
type VectorStore interface {
	// Store stores an embedded chunk
	Store(ctx context.Context, projectID string, chunk EmbeddedChunk) error

	// StoreBatch stores multiple embedded chunks
	StoreBatch(ctx context.Context, projectID string, chunks []EmbeddedChunk) error

	// Search performs vector similarity search
	Search(ctx context.Context, projectID string, query EmbeddingVector, topK int, minScore float32) ([]SemanticSearchResult, error)

	// Delete removes embeddings for a file
	Delete(ctx context.Context, projectID string, filePath string) error

	// DeleteProject removes all embeddings for a project
	DeleteProject(ctx context.Context, projectID string) error

	// GetStats returns statistics about stored embeddings
	GetStats(ctx context.Context, projectID string) (*VectorStoreStats, error)

	// GetChunkByID retrieves a specific chunk
	GetChunkByID(ctx context.Context, projectID string, chunkID string) (*EmbeddedChunk, error)

	// ListChunks lists all chunks for a file
	ListChunks(ctx context.Context, projectID string, filePath string) ([]EmbeddedChunk, error)
}

// VectorStoreStats contains statistics about the vector store
type VectorStoreStats struct {
	TotalChunks int       `json:"totalChunks"`
	TotalFiles  int       `json:"totalFiles"`
	TotalTokens int       `json:"totalTokens"`
	LastUpdated time.Time `json:"lastUpdated"`
	IndexSize   int64     `json:"indexSize"`
	Dimensions  int       `json:"dimensions"`
}

// SemanticSearchService provides semantic search capabilities
type SemanticSearchService interface {
	// IndexProject indexes all files in a project
	IndexProject(ctx context.Context, projectRoot string) error

	// IndexFile indexes a single file
	IndexFile(ctx context.Context, projectRoot string, filePath string) error

	// Search performs semantic search
	Search(ctx context.Context, req SemanticSearchRequest) (*SemanticSearchResponse, error)

	// FindSimilar finds similar code
	FindSimilar(ctx context.Context, req SimilarCodeRequest) (*SemanticSearchResponse, error)

	// GetClusters returns code clusters
	GetClusters(ctx context.Context, projectRoot string, numClusters int) ([]ClusterInfo, error)

	// GetStats returns indexing statistics
	GetStats(ctx context.Context, projectRoot string) (*VectorStoreStats, error)

	// IsIndexed checks if a project is indexed
	IsIndexed(ctx context.Context, projectRoot string) bool

	// InvalidateFile marks a file for re-indexing
	InvalidateFile(ctx context.Context, projectRoot string, filePath string) error
}

// RAGService provides Retrieval Augmented Generation
type RAGService interface {
	// RetrieveContext retrieves relevant context for a query
	RetrieveContext(ctx context.Context, query string, projectRoot string, maxTokens int) ([]CodeChunk, error)

	// RetrieveAndRank retrieves and re-ranks results
	RetrieveAndRank(ctx context.Context, query string, projectRoot string, topK int) ([]SemanticSearchResult, error)

	// HybridSearch performs hybrid keyword + semantic search
	HybridSearch(ctx context.Context, req SemanticSearchRequest) (*SemanticSearchResponse, error)
}
