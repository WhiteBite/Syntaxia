package ai

import (
	"sync"
	"time"
)

// MetricsCollector collects AI generation metrics (thread-safe)
type MetricsCollector struct {
	generations []GenerationMetric
	mu          sync.RWMutex
	maxSize     int
}

// GenerationMetric represents a single generation metric
type GenerationMetric struct {
	Provider   string
	Model      string
	Duration   time.Duration
	TokensUsed int
	Timestamp  time.Time
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{generations: make([]GenerationMetric, 0, 1000), maxSize: 1000}
}

// RecordGeneration records a generation metric
func (m *MetricsCollector) RecordGeneration(provider, model string, duration time.Duration, tokens int) {
	metric := GenerationMetric{Provider: provider, Model: model, Duration: duration, TokensUsed: tokens, Timestamp: time.Now()}

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.generations) >= m.maxSize {
		m.generations = m.generations[m.maxSize/10:]
	}
	m.generations = append(m.generations, metric)
}

// GetMetrics returns aggregated metrics
func (m *MetricsCollector) GetMetrics() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	totalTokens := 0
	totalDuration := time.Duration(0)
	providerCounts := make(map[string]int)

	for _, gen := range m.generations {
		totalTokens += gen.TokensUsed
		totalDuration += gen.Duration
		providerCounts[gen.Provider]++
	}

	return map[string]any{
		"total_generations": len(m.generations),
		"total_tokens":      totalTokens,
		"total_duration_ms": totalDuration.Milliseconds(),
		"by_provider":       providerCounts,
	}
}
