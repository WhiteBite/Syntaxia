package domain

import "time"

// AIRequest представляет унифицированный запрос к AI провайдеру.
type AIRequest struct {
	Model        string
	SystemPrompt string
	UserPrompt   string
	// Дополнительные параметры для интеллектуальной системы
	Temperature      float64
	MaxTokens        int
	TopP             float64
	FrequencyPenalty float64
	PresencePenalty  float64
	// Контекст для лучшего понимания
	ContextFiles []string
	ProjectType  string
	// Метаданные запроса
	RequestID  string
	Priority   RequestPriority
	RetryCount int
	Timeout    time.Duration
	// Грамматика для структурированного вывода
	Grammar string
}

// AIResponse представляет унифицированный ответ от AI провайдера.
type AIResponse struct {
	Content string
	// Дополнительные метаданные
	TokensUsed     int
	ModelUsed      string
	ProcessingTime time.Duration
	FinishReason   string
	// Метаданные для интеллектуальной системы
	Confidence  float64
	Suggestions []string
	NextActions []string
	// Ошибки и предупреждения
	Warnings []string
	Errors   []string
}

// RequestPriority определяет приоритет запроса
type RequestPriority int

const (
	PriorityLow RequestPriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// AIProviderConfig содержит конфигурацию для провайдера
type AIProviderConfig struct {
	ProviderType string
	APIKey       string
	BaseURL      string
	Timeout      time.Duration
	MaxRetries   int
	RateLimit    RateLimitConfig
}

// RateLimitConfig настройки ограничения скорости
type RateLimitConfig struct {
	RequestsPerMinute int
	TokensPerMinute   int
	BurstSize         int
}

// IntelligentRequest представляет интеллектуальный запрос с автоматической оптимизацией
type IntelligentRequest struct {
	BaseRequest    AIRequest
	Optimization   OptimizationConfig
	FallbackConfig FallbackConfig
	Monitoring     MonitoringConfig
}

// OptimizationConfig настройки оптимизации запроса
type OptimizationConfig struct {
	AutoOptimizePrompt bool
	ContextCompression bool
	TokenOptimization  bool
	ModelSelection     ModelSelectionStrategy
}

// ModelSelectionStrategy стратегия выбора модели
type ModelSelectionStrategy string

const (
	StrategyFastest  ModelSelectionStrategy = "fastest"
	StrategyCheapest ModelSelectionStrategy = "cheapest"
	StrategyBest     ModelSelectionStrategy = "best"
	StrategyBalanced ModelSelectionStrategy = "balanced"
)

// FallbackConfig настройки fallback стратегии
type FallbackConfig struct {
	EnableFallback      bool
	FallbackModels      []string
	FallbackProviders   []string
	MaxFallbackAttempts int
}

// MonitoringConfig настройки мониторинга
type MonitoringConfig struct {
	EnableMetrics        bool
	EnableLogging        bool
	EnableTracing        bool
	PerformanceThreshold time.Duration
}
