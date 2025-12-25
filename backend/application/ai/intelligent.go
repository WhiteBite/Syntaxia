package ai

import (
	"context"
	"fmt"
	"math/rand"
	"syntaxia/domain"
	"strings"
	"time"
)

// IntelligentService provides intelligent AI capabilities
type IntelligentService struct {
	settingsService SettingsProvider
	log             domain.Logger
	providerGetter  domain.AIProviderGetter
	rateLimiter     *RateLimiter
	metrics         *MetricsCollector
}

// NewIntelligentService creates a new intelligent AI service
func NewIntelligentService(
	settingsService SettingsProvider,
	log domain.Logger,
	rateLimiter *RateLimiter,
	metrics *MetricsCollector,
) *IntelligentService {
	return &IntelligentService{
		settingsService: settingsService,
		log:             log,
		rateLimiter:     rateLimiter,
		metrics:         metrics,
	}
}

// SetProviderGetter sets the provider source for AI access
func (s *IntelligentService) SetProviderGetter(getter domain.AIProviderGetter) {
	s.providerGetter = getter
}

// IntelligentGenerationOptions options for intelligent generation
type IntelligentGenerationOptions struct {
	Temperature            float64
	MaxTokens              int
	TopP                   float64
	Priority               domain.RequestPriority
	Timeout                time.Duration
	MaxRetries             int
	AutoOptimizePrompt     bool
	ContextCompression     bool
	TokenOptimization      bool
	ModelSelectionStrategy domain.ModelSelectionStrategy
	EnableFallback         bool
	FallbackModels         []string
	FallbackProviders      []string
	MaxFallbackAttempts    int
	PerformanceThreshold   time.Duration
	ProjectType            string
	CodeStyle              string
}

// IntelligentGenerationResult result of intelligent generation
type IntelligentGenerationResult struct {
	Content        string
	ModelUsed      string
	TokensUsed     int
	ProcessingTime time.Duration
	QualityScore   float64
	Suggestions    []string
	Warnings       []string
	RequestID      string
	Provider       string
}

// GenerateIntelligentCode performs intelligent code generation
func (s *IntelligentService) GenerateIntelligentCode(
	ctx context.Context, task, codeContext string, options IntelligentGenerationOptions,
) (*IntelligentGenerationResult, error) {
	startTime := time.Now()

	dto, err := s.settingsService.GetSettingsDTO()
	if err != nil {
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	intelligentReq := s.buildIntelligentRequest(task, dto, options)
	provider, model, err := s.selectOptimalProvider(ctx, intelligentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to select provider: %w", err)
	}

	optimizedPrompt := s.optimizePrompt(task, codeContext, options)

	req := domain.AIRequest{
		Model: model, SystemPrompt: s.buildSystemPrompt(options), UserPrompt: optimizedPrompt,
		Temperature: options.Temperature, MaxTokens: options.MaxTokens, TopP: options.TopP,
		RequestID: generateRequestID(), Priority: options.Priority, Timeout: options.Timeout,
	}

	if err := s.rateLimiter.CheckLimit(provider.GetProviderInfo().Name); err != nil {
		return nil, fmt.Errorf("rate limit exceeded: %w", err)
	}

	var response domain.AIResponse
	var lastErr error

	for attempt := 0; attempt <= options.MaxRetries; attempt++ {
		if attempt > 0 {
			s.log.Info(fmt.Sprintf("Retry attempt %d for request %s", attempt, req.RequestID))
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		response, lastErr = provider.Generate(ctx, req)
		if lastErr == nil {
			break
		}

		if attempt == options.MaxRetries && options.EnableFallback {
			response, lastErr = s.tryFallback(ctx, req)
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("all generation attempts failed: %w", lastErr)
	}

	analysis := s.analyzeResponse(response)
	s.metrics.RecordGeneration(provider.GetProviderInfo().Name, model, time.Since(startTime), response.TokensUsed)

	return &IntelligentGenerationResult{
		Content: response.Content, ModelUsed: response.ModelUsed, TokensUsed: response.TokensUsed,
		ProcessingTime: response.ProcessingTime, QualityScore: analysis.QualityScore,
		Suggestions: analysis.Suggestions, Warnings: response.Warnings,
		RequestID: req.RequestID, Provider: provider.GetProviderInfo().Name,
	}, nil
}

func (s *IntelligentService) buildIntelligentRequest(task string, _ domain.SettingsDTO, options IntelligentGenerationOptions) domain.IntelligentRequest {
	return domain.IntelligentRequest{
		BaseRequest: domain.AIRequest{UserPrompt: task, Priority: options.Priority, Timeout: options.Timeout},
		Optimization: domain.OptimizationConfig{
			AutoOptimizePrompt: options.AutoOptimizePrompt, ContextCompression: options.ContextCompression,
			TokenOptimization: options.TokenOptimization, ModelSelection: options.ModelSelectionStrategy,
		},
		FallbackConfig: domain.FallbackConfig{
			EnableFallback: options.EnableFallback, FallbackModels: options.FallbackModels,
			FallbackProviders: options.FallbackProviders, MaxFallbackAttempts: options.MaxFallbackAttempts,
		},
		Monitoring: domain.MonitoringConfig{EnableMetrics: true, EnableLogging: true, PerformanceThreshold: options.PerformanceThreshold},
	}
}

func (s *IntelligentService) selectOptimalProvider(ctx context.Context, req domain.IntelligentRequest) (domain.AIProvider, string, error) {
	if s.providerGetter == nil {
		return nil, "", fmt.Errorf("provider getter not initialized")
	}

	provider, model, err := s.providerGetter.GetProvider(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get provider: %w", err)
	}

	if model == "" {
		if intelligentProvider, ok := provider.(domain.IntelligentAIProvider); ok {
			model, err = intelligentProvider.SelectOptimalModel(ctx, req.BaseRequest, req.Optimization.ModelSelection)
			if err != nil {
				return nil, "", fmt.Errorf("failed to select optimal model: %w", err)
			}
		} else {
			return nil, "", fmt.Errorf("no model selected for provider")
		}
	}

	return provider, model, nil
}

func (s *IntelligentService) optimizePrompt(task, codeContext string, options IntelligentGenerationOptions) string {
	if !options.AutoOptimizePrompt {
		return task
	}

	optimized := strings.TrimSpace(task)
	if codeContext != "" {
		optimized = fmt.Sprintf("Context:\n%s\n\nTask:\n%s", codeContext, optimized)
	}
	if options.ProjectType != "" {
		optimized = fmt.Sprintf("%s\n\nPlease provide solution optimized for %s project.", optimized, options.ProjectType)
	}
	return optimized
}

func (s *IntelligentService) buildSystemPrompt(options IntelligentGenerationOptions) string {
	basePrompt := "You are an expert software developer. Your task is to implement the user's request by providing the necessary code changes in the form of a standard git diff."
	if options.ProjectType != "" {
		basePrompt += fmt.Sprintf(" Focus on %s best practices.", options.ProjectType)
	}
	if options.CodeStyle != "" {
		basePrompt += fmt.Sprintf(" Follow %s coding style.", options.CodeStyle)
	}
	return basePrompt
}

func (s *IntelligentService) tryFallback(_ context.Context, _ domain.AIRequest) (domain.AIResponse, error) {
	s.log.Warning("Fallback providers not supported in current implementation")
	return domain.AIResponse{}, fmt.Errorf("fallback not supported")
}

func (s *IntelligentService) analyzeResponse(response domain.AIResponse) domain.ResponseAnalysis {
	analysis := domain.ResponseAnalysis{QualityScore: 0.8, RelevanceScore: 0.8, CompletenessScore: 0.8, Confidence: 0.8}
	if strings.Contains(response.Content, "diff --git") {
		analysis.QualityScore += 0.1
		analysis.CompletenessScore += 0.1
	}
	if len(response.Content) > 100 {
		analysis.CompletenessScore += 0.1
	}
	analysis.Suggestions = []string{"Review the generated code before applying", "Test the changes in a development environment"}
	return analysis
}

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%d", time.Now().Unix(), rand.Intn(1000)) //nolint:gosec
}
