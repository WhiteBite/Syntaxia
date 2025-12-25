package ai

import (
	"context"
	"syntaxia/domain"
)

// ProviderConfig contains configuration for creating AI providers and model fetchers
type ProviderConfig struct {
	FactoryFunc  func(apiKey, host string, log domain.Logger) (domain.AIProvider, error)
	ModelFetcher func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error)
}

// GetProviderRegistry returns a unified registry of all AI provider configurations
// This eliminates duplication between createModelFetchers and createProviderFactory
func GetProviderRegistry(openRouterHost string) map[string]ProviderConfig {
	return map[string]ProviderConfig{
		"gemini": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				return NewGemini(apiKey, host, log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewGemini(apiKey, host, log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
		"openai": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				return NewOpenAI(apiKey, host, log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewOpenAI(apiKey, host, log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
		"openrouter": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				effectiveHost := host
				if effectiveHost == "" {
					effectiveHost = openRouterHost
				}
				return NewOpenAI(apiKey, effectiveHost, log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewOpenAI(apiKey, openRouterHost, log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
		"localai": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				return NewLocalAI(apiKey, host, log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewLocalAI(apiKey, host, log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
		"qwen": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				return NewQwen(apiKey, host, log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewQwen(apiKey, host, log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
		"qwen-cli": {
			FactoryFunc: func(apiKey, host string, log domain.Logger) (domain.AIProvider, error) {
				return NewQwenCLI(log)
			},
			ModelFetcher: func(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
				p, err := NewQwenCLI(log)
				if err != nil {
					return nil, err
				}
				return p.ListModels(ctx)
			},
		},
	}
}

// HostResolver resolves provider-specific host configuration at runtime.
type HostResolver func(providerType string) (string, error)

// NewAIProviderFactoryRegistry builds a registry mapping provider type to domain.AIProviderFactory implementations.
// This enables Open/Closed principle compliance for selecting providers in application services.
func NewAIProviderFactoryRegistry(
	log domain.Logger,
	openRouterHost string,
	resolveHost HostResolver,
) map[string]domain.AIProviderFactory {
	configs := GetProviderRegistry(openRouterHost)
	factories := make(map[string]domain.AIProviderFactory, len(configs))

	for providerType, cfg := range configs {
		cfg := cfg
		factories[providerType] = func(pt, apiKey string) (domain.AIProvider, error) {
			host := ""
			if resolveHost != nil {
				resolvedHost, err := resolveHost(pt)
				if err != nil {
					return nil, err
				}
				host = resolvedHost
			}
			return cfg.FactoryFunc(apiKey, host, log)
		}
	}

	return factories
}
