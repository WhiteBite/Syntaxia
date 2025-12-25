package domain

import "context"

// StreamChunk represents a chunk of streamed response
type StreamChunk struct {
	Content      string `json:"content"`
	Done         bool   `json:"done"`
	Error        string `json:"error,omitempty"`
	TokensUsed   int    `json:"tokensUsed,omitempty"`
	FinishReason string `json:"finishReason,omitempty"`
}

// AIProvider определяет контракт для любого провайдера LLM.
// Это позволяет легко заменять OpenAI на Gemini, Claude или локальные модели.
type AIProvider interface {
	// Generate выполняет основной запрос к LLM для генерации контента.
	Generate(ctx context.Context, req AIRequest) (AIResponse, error)

	// GenerateStream выполняет запрос с потоковой передачей ответа
	GenerateStream(ctx context.Context, req AIRequest, onChunk func(chunk StreamChunk)) error

	// ListModels возвращает список доступных моделей
	ListModels(ctx context.Context) ([]string, error)

	// GetProviderInfo возвращает информацию о провайдере
	GetProviderInfo() ProviderInfo

	// ValidateRequest проверяет корректность запроса
	ValidateRequest(req AIRequest) error

	// EstimateTokens оценивает количество токенов в запросе
	EstimateTokens(req AIRequest) (int, error)

	// GetPricing возвращает информацию о стоимости
	GetPricing(model string) PricingInfo
}

// IntelligentAIProvider расширенный интерфейс для интеллектуальных возможностей
type IntelligentAIProvider interface {
	AIProvider

	// GenerateIntelligent выполняет интеллектуальный запрос с автоматической оптимизацией
	GenerateIntelligent(ctx context.Context, req IntelligentRequest) (AIResponse, error)

	// OptimizePrompt автоматически оптимизирует промпт
	OptimizePrompt(ctx context.Context, originalPrompt string, context []string) (string, error)

	// SelectOptimalModel выбирает оптимальную модель на основе запроса
	SelectOptimalModel(ctx context.Context, req AIRequest, strategy ModelSelectionStrategy) (string, error)

	// CompressContext сжимает контекст для экономии токенов
	CompressContext(ctx context.Context, context []string, maxTokens int) ([]string, error)

	// AnalyzeResponse анализирует ответ и предоставляет рекомендации
	AnalyzeResponse(ctx context.Context, response AIResponse, originalRequest AIRequest) (ResponseAnalysis, error)
}

// ProviderInfo информация о провайдере
type ProviderInfo struct {
	Name            string
	Version         string
	Capabilities    []string
	Limitations     []string
	SupportedModels []string
}

// PricingInfo информация о стоимости
type PricingInfo struct {
	InputTokensPer1K  float64
	OutputTokensPer1K float64
	Currency          string
	Model             string
}

// ResponseAnalysis анализ ответа
type ResponseAnalysis struct {
	QualityScore      float64
	RelevanceScore    float64
	CompletenessScore float64
	Suggestions       []string
	Improvements      []string
	Confidence        float64
}
