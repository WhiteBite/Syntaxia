// Package rag provides backward compatibility aliases for semantic search.
// The actual implementation is in syntaxia/application/semantic.
package rag

import (
	"syntaxia/application/semantic"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// SemanticSearchService is an alias for semantic.ServiceImpl
type SemanticSearchService = semantic.ServiceImpl

// IndexingState is an alias for semantic.IndexingState
type IndexingState = semantic.IndexingState

// SymbolInfoForChunking is an alias for semantic.SymbolInfoForChunking
type SymbolInfoForChunking = semantic.SymbolInfoForChunking

// NewSemanticSearchService creates a new semantic search service.
func NewSemanticSearchService(
	embeddingProvider domain.EmbeddingProvider,
	vectorStore domain.VectorStore,
	symbolIndex analysis.SymbolIndex,
	log domain.Logger,
	chunker domain.CodeChunker,
) *SemanticSearchService {
	return semantic.NewService(embeddingProvider, vectorStore, symbolIndex, log, chunker)
}
