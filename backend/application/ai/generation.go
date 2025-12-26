package ai

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"sync/atomic"
	"time"
)

// GenerationOptions options for code generation
type GenerationOptions struct {
	Model       string
	Temperature float64
	MaxTokens   int
	TopP        float64
	Priority    domain.RequestPriority
	Timeout     time.Duration
}

type generationParams struct {
	model       string
	temperature float64
	maxTokens   int
	topP        float64
	timeout     time.Duration
	priority    domain.RequestPriority
	useCache    bool
}

func applyOptions(params *generationParams, options *GenerationOptions) {
	if options == nil {
		return
	}
	if options.Model != "" {
		params.model = options.Model
	}
	if options.Temperature != 0 {
		params.temperature = options.Temperature
		if params.temperature > deterministicTempThreshold {
			params.useCache = false
		}
	}
	if options.MaxTokens != 0 {
		params.maxTokens = options.MaxTokens
	}
	if options.TopP != 0 {
		params.topP = options.TopP
	}
	if options.Timeout != 0 {
		params.timeout = options.Timeout
	}
	if options.Priority != domain.PriorityLow {
		params.priority = options.Priority
	}
}

func (s *Service) checkCache(cacheKey string, useCache bool) (string, bool) {
	if !useCache {
		return "", false
	}
	s.responseCacheMu.RLock()
	defer s.responseCacheMu.RUnlock()
	if cached, ok := s.responseCache[cacheKey]; ok && time.Since(cached.timestamp) < responseCacheTTL {
		atomic.AddInt64(&s.cacheHits, 1)
		return cached.content, true
	}
	return "", false
}

func (s *Service) generateCodeInternal(ctx context.Context, systemPrompt, userPrompt string, options *GenerationOptions) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	atomic.AddInt64(&s.totalRequests, 1)

	provider, model, err := s.getProvider(ctx)
	if err != nil {
		return "", err
	}

	params := &generationParams{
		model: model, temperature: DefaultTemperature, maxTokens: DefaultMaxTokens,
		topP: DefaultTopP, timeout: DefaultTimeout, priority: domain.PriorityNormal, useCache: true,
	}
	applyOptions(params, options)

	cacheKey := s.getCacheKey(systemPrompt, userPrompt, params.model, params.temperature, params.maxTokens, params.topP)
	if content, found := s.checkCache(cacheKey, params.useCache); found {
		return content, nil
	}
	atomic.AddInt64(&s.cacheMisses, 1)

	req := domain.AIRequest{
		Model: params.model, SystemPrompt: systemPrompt, UserPrompt: userPrompt,
		Temperature: params.temperature, MaxTokens: params.maxTokens, TopP: params.topP,
		RequestID: fmt.Sprintf("req_%d", time.Now().UnixNano()),
		Priority:  params.priority, Timeout: params.timeout,
	}

	tctx, cancel := context.WithTimeout(ctx, params.timeout)
	defer cancel()

	resp, err := provider.Generate(tctx, req)
	if err != nil {
		return "", fmt.Errorf("AI generation failed: %w", err)
	}

	if params.useCache && resp.Content != "" {
		s.responseCacheMu.Lock()
		s.responseCache[cacheKey] = &cachedAIResponse{content: resp.Content, timestamp: time.Now(), tokens: resp.TokensUsed}
		s.responseCacheMu.Unlock()
	}

	atomic.AddInt64(&s.totalTokensUsed, int64(resp.TokensUsed))
	return resp.Content, nil
}

// GenerateCode generates code using AI
func (s *Service) GenerateCode(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	return s.generateCodeInternal(ctx, systemPrompt, userPrompt, nil)
}

// GenerateCodeWithOptions generates code with additional options
func (s *Service) GenerateCodeWithOptions(ctx context.Context, systemPrompt, userPrompt string, options GenerationOptions) (string, error) {
	return s.generateCodeInternal(ctx, systemPrompt, userPrompt, &options)
}

// GenerateCodeStream generates code with streaming response
func (s *Service) GenerateCodeStream(ctx context.Context, systemPrompt, userPrompt string, onChunk func(chunk domain.StreamChunk)) error {
	if ctx == nil {
		ctx = context.Background()
	}
	atomic.AddInt64(&s.totalRequests, 1)

	provider, model, err := s.getProvider(ctx)
	if err != nil {
		onChunk(domain.StreamChunk{Done: true, Error: err.Error()})
		return err
	}

	req := domain.AIRequest{
		Model: model, SystemPrompt: systemPrompt, UserPrompt: userPrompt,
		Temperature: DefaultTemperature, MaxTokens: DefaultMaxTokens, TopP: DefaultTopP,
		RequestID: fmt.Sprintf("stream_%d", time.Now().UnixNano()),
		Priority:  domain.PriorityNormal, Timeout: DefaultStreamTimeout,
	}

	tctx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	return provider.GenerateStream(tctx, req, onChunk)
}

// GenerateIntelligentCode uses intelligent system for code generation
func (s *Service) GenerateIntelligentCode(ctx context.Context, task, codeContext string, options IntelligentGenerationOptions) (*IntelligentGenerationResult, error) {
	return s.intelligentService.GenerateIntelligentCode(ctx, task, codeContext, options)
}

// GetProviderInfo returns information about current provider
func (s *Service) GetProviderInfo(ctx context.Context) (*domain.ProviderInfo, error) {
	provider, _, err := s.getProvider(ctx)
	if err != nil {
		return nil, err
	}
	info := provider.GetProviderInfo()
	return &info, nil
}

// ListAvailableModels returns list of available models
func (s *Service) ListAvailableModels(ctx context.Context) ([]string, error) {
	provider, _, err := s.getProvider(ctx)
	if err != nil {
		return nil, err
	}
	return provider.ListModels(ctx)
}
