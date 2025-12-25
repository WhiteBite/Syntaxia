package context

import (
	"context"
	"syntaxia/domain"
)

// ContextBuilderAdapter adapts Service to implement domain.ContextBuilder interface
type ContextBuilderAdapter struct {
	service *Service
}

// NewContextBuilderAdapter creates a new adapter that wraps Service
func NewContextBuilderAdapter(service *Service) *ContextBuilderAdapter {
	return &ContextBuilderAdapter{service: service}
}

// BuildContext implements domain.ContextBuilder interface
func (a *ContextBuilderAdapter) BuildContext(ctx context.Context, projectPath string, includedPaths []string, options *domain.ContextBuildOptions) (*domain.ContextSummary, error) {
	return a.service.BuildContextSummary(ctx, projectPath, includedPaths, options)
}
