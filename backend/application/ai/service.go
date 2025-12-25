// Package ai provides AI service functionality for code generation.
package ai

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"syntaxia/domain"
	"sync"
	"sync/atomic"
	"time"
)

// Service orchestrates AI-related operations with caching.
type Service struct {
	settingsService    SettingsProvider
	log                domain.Logger
	providerRegistry   map[string]domain.AIProviderFactory
	intelligentService *IntelligentService

	providerCache   map[string]domain.AIProvider
	providerCacheMu sync.RWMutex

	responseCache   map[string]*cachedAIResponse
	responseCacheMu sync.RWMutex

	totalRequests   int64
	cacheHits       int64
	cacheMisses     int64
	totalTokensUsed int64

	stopCh   chan struct{}
	stopOnce sync.Once
	wg       sync.WaitGroup
}

// SettingsProvider interface for settings access
type SettingsProvider interface {
	GetSettingsDTO() (domain.SettingsDTO, error)
}

type cachedAIResponse struct {
	content   string
	timestamp time.Time
	tokens    int
}

const (
	maxResponseCacheSize       = 100
	responseCacheTTL           = 30 * time.Minute
	cacheCleanupInterval       = 5 * time.Minute
	DefaultTemperature         = 0.7
	DefaultMaxTokens           = 4000
	DefaultTopP                = 1.0
	DefaultTimeout             = 60 * time.Second
	DefaultStreamTimeout       = 120 * time.Second
	deterministicTempThreshold = 0.1
)

// NewService creates a new AI Service.
func NewService(
	settingsService SettingsProvider,
	log domain.Logger,
	providerRegistry map[string]domain.AIProviderFactory,
	intelligentService *IntelligentService,
) *Service {
	service := &Service{
		settingsService:    settingsService,
		log:                log,
		providerRegistry:   providerRegistry,
		intelligentService: intelligentService,
		providerCache:      make(map[string]domain.AIProvider),
		responseCache:      make(map[string]*cachedAIResponse),
		stopCh:             make(chan struct{}),
	}
	service.wg.Add(1)
	go service.cleanupResponseCache()
	return service
}

// Shutdown gracefully stops the AI service
func (s *Service) Shutdown(ctx context.Context) error {
	s.stopOnce.Do(func() { close(s.stopCh) })

	done := make(chan struct{})
	go func() { s.wg.Wait(); close(done) }()

	select {
	case <-done:
		s.log.Info("AIService shutdown complete")
	case <-ctx.Done():
		s.log.Warning("AIService shutdown timed out")
		return ctx.Err()
	}

	s.responseCacheMu.Lock()
	s.responseCache = make(map[string]*cachedAIResponse)
	s.responseCacheMu.Unlock()

	s.providerCacheMu.Lock()
	s.providerCache = make(map[string]domain.AIProvider)
	s.providerCacheMu.Unlock()

	return nil
}

func (s *Service) cleanupResponseCache() {
	defer s.wg.Done()
	ticker := time.NewTicker(cacheCleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.performCacheCleanup()
		}
	}
}

func (s *Service) performCacheCleanup() {
	s.responseCacheMu.Lock()
	defer s.responseCacheMu.Unlock()

	now := time.Now()
	for key, entry := range s.responseCache {
		if now.Sub(entry.timestamp) > responseCacheTTL {
			delete(s.responseCache, key)
		}
	}

	for len(s.responseCache) > maxResponseCacheSize {
		oldest := time.Now()
		oldestKey := ""
		for key, entry := range s.responseCache {
			if entry.timestamp.Before(oldest) {
				oldest = entry.timestamp
				oldestKey = key
			}
		}
		if oldestKey == "" {
			break
		}
		delete(s.responseCache, oldestKey)
	}
}

func (s *Service) getCacheKey(systemPrompt, userPrompt, model string, temperature float64, maxTokens int, topP float64) string {
	h := sha256.New()
	h.Write([]byte(systemPrompt))
	h.Write([]byte(userPrompt))
	h.Write([]byte(model))
	fmt.Fprintf(h, "%.2f:%d:%.2f", temperature, maxTokens, topP)
	return hex.EncodeToString(h.Sum(nil))[:16]
}

// GetMetrics returns AI service metrics
func (s *Service) GetMetrics() map[string]any {
	s.responseCacheMu.RLock()
	cacheSize := len(s.responseCache)
	s.responseCacheMu.RUnlock()

	s.providerCacheMu.RLock()
	providerCount := len(s.providerCache)
	s.providerCacheMu.RUnlock()

	return map[string]any{
		"total_requests":      atomic.LoadInt64(&s.totalRequests),
		"cache_hits":          atomic.LoadInt64(&s.cacheHits),
		"cache_misses":        atomic.LoadInt64(&s.cacheMisses),
		"total_tokens_used":   atomic.LoadInt64(&s.totalTokensUsed),
		"response_cache_size": cacheSize,
		"cached_providers":    providerCount,
	}
}

// GetIntelligentService returns the intelligent AI service
func (s *Service) GetIntelligentService() *IntelligentService {
	return s.intelligentService
}
