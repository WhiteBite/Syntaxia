package domain

import "context"

// AIProviderGetter provides access to AI providers
// This interface breaks the circular dependency between AIService and IntelligentAIService
type AIProviderGetter interface {
	// GetProvider returns the current AI provider and selected model
	GetProvider(ctx context.Context) (AIProvider, string, error)
}
