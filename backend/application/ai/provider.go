package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"syntaxia/domain"
)

// GetProvider returns current provider and model
func (s *Service) GetProvider(ctx context.Context) (domain.AIProvider, string, error) {
	return s.getProvider(ctx)
}

func (s *Service) getProvider(_ context.Context) (domain.AIProvider, string, error) {
	dto, err := s.settingsService.GetSettingsDTO()
	if err != nil {
		return nil, "", fmt.Errorf("could not get settings: %w", err)
	}

	providerType := dto.SelectedProvider
	if providerType == "" {
		return nil, "", fmt.Errorf("no AI provider selected")
	}

	apiKey := s.getAPIKey(dto, providerType)
	if apiKey == "" && providerType != "localai" && providerType != "qwen-cli" {
		return nil, "", fmt.Errorf("API key for %s is not set", providerType)
	}

	cacheKey := fmt.Sprintf("%s:%s", providerType, s.hashKey(apiKey))

	s.providerCacheMu.RLock()
	if cachedProvider, ok := s.providerCache[cacheKey]; ok {
		s.providerCacheMu.RUnlock()
		return cachedProvider, s.getModelForProvider(dto, providerType), nil
	}
	s.providerCacheMu.RUnlock()

	factory, exists := s.providerRegistry[providerType]
	if !exists {
		return nil, "", fmt.Errorf("no factory registered for provider %s", providerType)
	}

	provider, err := factory(providerType, apiKey)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create provider %s: %w", providerType, err)
	}

	s.providerCacheMu.Lock()
	s.providerCache[cacheKey] = provider
	s.providerCacheMu.Unlock()

	return provider, s.getModelForProvider(dto, providerType), nil
}

func (s *Service) getAPIKey(dto domain.SettingsDTO, providerType string) string {
	switch providerType {
	case "openai":
		return dto.OpenAIAPIKey
	case "gemini":
		return dto.GeminiAPIKey
	case "openrouter":
		return dto.OpenRouterAPIKey
	case "localai":
		return dto.LocalAIAPIKey
	case "qwen":
		return dto.QwenAPIKey
	case "qwen-cli":
		return ""
	default:
		return ""
	}
}

func (s *Service) hashKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])[:8]
}

func (s *Service) getModelForProvider(dto domain.SettingsDTO, providerType string) string {
	model, ok := dto.SelectedModels[providerType]
	if !ok || model == "" {
		if models, ok := dto.AvailableModels[providerType]; ok && len(models) > 0 {
			model = models[0]
			s.log.Warning(fmt.Sprintf("No model selected for %s, falling back to: %s", providerType, model))
		}
	}
	return model
}

// InvalidateProviderCache clears the provider cache
func (s *Service) InvalidateProviderCache() {
	s.providerCacheMu.Lock()
	s.providerCache = make(map[string]domain.AIProvider)
	s.providerCacheMu.Unlock()
	s.log.Info("Provider cache invalidated")
}
