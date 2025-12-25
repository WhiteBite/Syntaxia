package ai

import (
	"fmt"
	"syntaxia/domain"
)

// KeyExtractorFunc defines a function type for extracting API keys
type KeyExtractorFunc func(settings domain.SettingsDTO) string

// KeyResolverImpl implements the KeyResolver interface using a registry pattern
type KeyResolverImpl struct {
	extractors map[string]KeyExtractorFunc
}

// NewKeyResolverImpl creates a new KeyResolverImpl with default extractors
func NewKeyResolverImpl() *KeyResolverImpl {
	return &KeyResolverImpl{
		extractors: map[string]KeyExtractorFunc{
			"openai":     func(s domain.SettingsDTO) string { return s.OpenAIAPIKey },
			"gemini":     func(s domain.SettingsDTO) string { return s.GeminiAPIKey },
			"openrouter": func(s domain.SettingsDTO) string { return s.OpenRouterAPIKey },
			"localai":    func(s domain.SettingsDTO) string { return s.LocalAIAPIKey },
		},
	}
}

// GetKey retrieves the API key for a given provider type from settings
func (r *KeyResolverImpl) GetKey(providerType string, settings domain.SettingsDTO) (string, error) {
	extractor, exists := r.extractors[providerType]
	if !exists {
		return "", fmt.Errorf("no key extractor registered for provider: %s", providerType)
	}

	key := extractor(settings)
	return key, nil
}

// RegisterExtractor registers a new key extractor for a provider type
func (r *KeyResolverImpl) RegisterExtractor(providerType string, extractor KeyExtractorFunc) {
	r.extractors[providerType] = extractor
}

// GetSupportedProviders returns a list of supported provider types
func (r *KeyResolverImpl) GetSupportedProviders() []string {
	providers := make([]string, 0, len(r.extractors))
	for provider := range r.extractors {
		providers = append(providers, provider)
	}
	return providers
}
