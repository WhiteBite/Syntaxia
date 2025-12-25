package symbolgraph

import (
	"context"

	"syntaxia/application/rag"
	"syntaxia/domain"
)

// CallStackAnalyzerAdapter adapts CallStackAnalyzer to application.CallStackAnalyzerInterface
type CallStackAnalyzerAdapter struct {
	analyzer *CallStackAnalyzer
}

// NewCallStackAnalyzerAdapter creates a new adapter
func NewCallStackAnalyzerAdapter(log domain.Logger) *CallStackAnalyzerAdapter {
	return &CallStackAnalyzerAdapter{
		analyzer: NewCallStackAnalyzer(log),
	}
}

// AnalyzeCallStack implements rag.CallStackAnalyzerInterface
func (a *CallStackAnalyzerAdapter) AnalyzeCallStack(
	ctx context.Context,
	projectRoot, filePath, symbolName string,
	maxDepth int,
) (*rag.CallStackResult, error) {
	result, err := a.analyzer.AnalyzeCallStack(ctx, projectRoot, filePath, symbolName, maxDepth)
	if err != nil {
		return nil, err
	}

	// Convert to application type
	return &rag.CallStackResult{
		RootSymbol:   result.RootSymbol,
		Callers:      result.Callers,
		Callees:      result.Callees,
		Dependencies: result.Dependencies,
		RelatedFiles: result.RelatedFiles,
		TotalSymbols: result.TotalSymbols,
	}, nil
}

// GetTransitiveDependencies implements application.CallStackAnalyzerInterface
func (a *CallStackAnalyzerAdapter) GetTransitiveDependencies(
	ctx context.Context,
	projectRoot, filePath, symbolName string,
	maxDepth int,
) ([]*domain.SymbolNode, error) {
	return a.analyzer.GetTransitiveDependencies(ctx, projectRoot, filePath, symbolName, maxDepth)
}
