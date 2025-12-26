package ai

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter implements token bucket rate limiting per provider
type RateLimiter struct {
	buckets map[string]*tokenBucket
	mu      sync.RWMutex
	config  map[string]rateLimitConfig
}

type tokenBucket struct {
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

type rateLimitConfig struct {
	tokensPerSecond float64
	maxTokens       float64
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string]*tokenBucket),
		config: map[string]rateLimitConfig{
			"openai":     {tokensPerSecond: 10, maxTokens: 60},
			"gemini":     {tokensPerSecond: 10, maxTokens: 60},
			"openrouter": {tokensPerSecond: 5, maxTokens: 30},
			"localai":    {tokensPerSecond: 100, maxTokens: 100},
		},
	}
}

// CheckLimit checks if request is within rate limit
func (r *RateLimiter) CheckLimit(provider string) error {
	r.mu.RLock()
	bucket, exists := r.buckets[provider]
	r.mu.RUnlock()

	if !exists {
		r.mu.Lock()
		if bucket, exists = r.buckets[provider]; !exists {
			config := r.getConfig(provider)
			bucket = &tokenBucket{tokens: config.maxTokens, lastRefill: time.Now()}
			r.buckets[provider] = bucket
		}
		r.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	config := r.getConfig(provider)
	elapsed := time.Since(bucket.lastRefill).Seconds()
	bucket.tokens += elapsed * config.tokensPerSecond
	if bucket.tokens > config.maxTokens {
		bucket.tokens = config.maxTokens
	}
	bucket.lastRefill = time.Now()

	if bucket.tokens < 1 {
		return fmt.Errorf("rate limit exceeded for provider %s, please wait", provider)
	}

	bucket.tokens--
	return nil
}

func (r *RateLimiter) getConfig(provider string) rateLimitConfig {
	if config, ok := r.config[provider]; ok {
		return config
	}
	return rateLimitConfig{tokensPerSecond: 10, maxTokens: 60}
}
