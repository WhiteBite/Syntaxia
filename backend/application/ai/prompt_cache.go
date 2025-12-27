// Package ai provides AI service functionality including prompt caching.
package ai

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"syntaxia/domain"
)

// PromptCacheService manages prompt caching for AI providers (Anthropic/OpenAI style).
// It separates prompts into cached (static) and dynamic parts to optimize API usage.
type PromptCacheService struct {
	log domain.Logger

	// Cache for static content hashes and their cache status
	staticCache   map[string]*CachedPromptBlock
	staticCacheMu sync.RWMutex

	// Configuration
	config PromptCacheConfig
}

// PromptCacheConfig configuration for prompt caching
type PromptCacheConfig struct {
	// MinCacheableTokens minimum tokens for a block to be cacheable
	MinCacheableTokens int
	// MaxCacheAge maximum age for cached content
	MaxCacheAge time.Duration
	// RepoMapRefreshInterval how often to refresh repo map cache
	RepoMapRefreshInterval time.Duration
}

// CachedPromptBlock represents a cached prompt block
type CachedPromptBlock struct {
	Hash        string    `json:"hash"`
	Content     string    `json:"content"`
	TokenCount  int       `json:"tokenCount"`
	CachedAt    time.Time `json:"cachedAt"`
	LastUsed    time.Time `json:"lastUsed"`
	UseCount    int       `json:"useCount"`
	BlockType   BlockType `json:"blockType"`
	CacheStatus string    `json:"cacheStatus"` // "hit", "miss", "pending"
}

// BlockType identifies the type of prompt block
type BlockType string

const (
	BlockTypeSystem       BlockType = "system"
	BlockTypeRepoMap      BlockType = "repo_map"
	BlockTypeProjectRules BlockType = "project_rules"
	BlockTypeSelectedFile BlockType = "selected_file"
	BlockTypeUserMessage  BlockType = "user_message"
)

// PromptBlock represents a single block in the structured prompt
type PromptBlock struct {
	Type       BlockType `json:"type"`
	Content    string    `json:"content"`
	Cacheable  bool      `json:"cacheable"`
	TokenCount int       `json:"tokenCount"`
	Hash       string    `json:"hash,omitempty"`
}

// StructuredPrompt represents a prompt with cached and dynamic parts
type StructuredPrompt struct {
	// CachedBlocks are static parts that can be cached by the provider
	CachedBlocks []PromptBlock `json:"cachedBlocks"`
	// DynamicBlocks are parts that change per request
	DynamicBlocks []PromptBlock `json:"dynamicBlocks"`
	// TotalTokens estimated total tokens
	TotalTokens int `json:"totalTokens"`
	// CachedTokens tokens in cached blocks
	CachedTokens int `json:"cachedTokens"`
	// DynamicTokens tokens in dynamic blocks
	DynamicTokens int `json:"dynamicTokens"`
}

// CacheHeaders represents cache control headers for AI providers
type CacheHeaders struct {
	// Anthropic-style cache control
	AnthropicCacheControl []AnthropicCacheBlock `json:"anthropicCacheControl,omitempty"`
	// OpenAI-style cache hints
	OpenAICacheHints *OpenAICacheHints `json:"openaiCacheHints,omitempty"`
}

// AnthropicCacheBlock represents Anthropic's cache_control block
type AnthropicCacheBlock struct {
	Type         string `json:"type"` // "ephemeral"
	BlockIndex   int    `json:"blockIndex"`
	CacheBreaker string `json:"cacheBreaker,omitempty"`
}

// OpenAICacheHints represents OpenAI's caching hints
type OpenAICacheHints struct {
	// PredictedOutputTokens for response prediction caching
	PredictedOutputTokens int `json:"predictedOutputTokens,omitempty"`
	// StaticPrefix indicates static content prefix length
	StaticPrefix int `json:"staticPrefix,omitempty"`
}

// Default configuration values
const (
	DefaultMinCacheableTokens    = 1024
	DefaultMaxCacheAge           = 5 * time.Minute
	DefaultRepoMapRefreshMinutes = 10
	CharsPerToken                = 4
)

// NewPromptCacheService creates a new PromptCacheService
func NewPromptCacheService(log domain.Logger, config *PromptCacheConfig) *PromptCacheService {
	cfg := PromptCacheConfig{
		MinCacheableTokens:     DefaultMinCacheableTokens,
		MaxCacheAge:            DefaultMaxCacheAge,
		RepoMapRefreshInterval: DefaultRepoMapRefreshMinutes * time.Minute,
	}
	if config != nil {
		if config.MinCacheableTokens > 0 {
			cfg.MinCacheableTokens = config.MinCacheableTokens
		}
		if config.MaxCacheAge > 0 {
			cfg.MaxCacheAge = config.MaxCacheAge
		}
		if config.RepoMapRefreshInterval > 0 {
			cfg.RepoMapRefreshInterval = config.RepoMapRefreshInterval
		}
	}

	return &PromptCacheService{
		log:         log,
		staticCache: make(map[string]*CachedPromptBlock),
		config:      cfg,
	}
}

// BuildStructuredPrompt builds a structured prompt with cached and dynamic parts
func (s *PromptCacheService) BuildStructuredPrompt(input PromptInput) *StructuredPrompt {
	result := &StructuredPrompt{
		CachedBlocks:  make([]PromptBlock, 0),
		DynamicBlocks: make([]PromptBlock, 0),
	}

	// 1. System prompt → CACHED
	if input.SystemPrompt != "" {
		block := s.createBlock(BlockTypeSystem, input.SystemPrompt, true)
		result.CachedBlocks = append(result.CachedBlocks, block)
		result.CachedTokens += block.TokenCount
	}

	// 2. Repo map → CACHED (update rarely)
	if input.RepoMap != "" {
		block := s.createBlock(BlockTypeRepoMap, input.RepoMap, true)
		result.CachedBlocks = append(result.CachedBlocks, block)
		result.CachedTokens += block.TokenCount
	}

	// 3. Project rules → CACHED
	if input.ProjectRules != "" {
		block := s.createBlock(BlockTypeProjectRules, input.ProjectRules, true)
		result.CachedBlocks = append(result.CachedBlocks, block)
		result.CachedTokens += block.TokenCount
	}

	// 4. Selected files → DYNAMIC
	for _, file := range input.SelectedFiles {
		block := s.createBlock(BlockTypeSelectedFile, file.Content, false)
		block.Hash = file.Path // Use path as identifier
		result.DynamicBlocks = append(result.DynamicBlocks, block)
		result.DynamicTokens += block.TokenCount
	}

	// 5. User message → DYNAMIC
	if input.UserMessage != "" {
		block := s.createBlock(BlockTypeUserMessage, input.UserMessage, false)
		result.DynamicBlocks = append(result.DynamicBlocks, block)
		result.DynamicTokens += block.TokenCount
	}

	result.TotalTokens = result.CachedTokens + result.DynamicTokens
	return result
}

// PromptInput represents input for building a structured prompt
type PromptInput struct {
	SystemPrompt  string         `json:"systemPrompt"`
	RepoMap       string         `json:"repoMap"`
	ProjectRules  string         `json:"projectRules"`
	SelectedFiles []SelectedFile `json:"selectedFiles"`
	UserMessage   string         `json:"userMessage"`
}

// SelectedFile represents a file selected for context
type SelectedFile struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// createBlock creates a prompt block with token estimation
func (s *PromptCacheService) createBlock(blockType BlockType, content string, cacheable bool) PromptBlock {
	tokenCount := s.estimateTokens(content)
	hash := s.hashContent(content)

	// Only cache if meets minimum token threshold
	if cacheable && tokenCount < s.config.MinCacheableTokens {
		cacheable = false
	}

	return PromptBlock{
		Type:       blockType,
		Content:    content,
		Cacheable:  cacheable,
		TokenCount: tokenCount,
		Hash:       hash,
	}
}

// GenerateCacheHeaders generates provider-specific cache headers
func (s *PromptCacheService) GenerateCacheHeaders(
	prompt *StructuredPrompt,
	provider string,
) *CacheHeaders {
	headers := &CacheHeaders{}

	switch provider {
	case "anthropic", "claude":
		headers.AnthropicCacheControl = s.generateAnthropicHeaders(prompt)
	case "openai":
		headers.OpenAICacheHints = s.generateOpenAIHeaders(prompt)
	}

	return headers
}

// generateAnthropicHeaders generates Anthropic-style cache control
func (s *PromptCacheService) generateAnthropicHeaders(prompt *StructuredPrompt) []AnthropicCacheBlock {
	blocks := make([]AnthropicCacheBlock, 0)

	for i, block := range prompt.CachedBlocks {
		if block.Cacheable {
			blocks = append(blocks, AnthropicCacheBlock{
				Type:       "ephemeral",
				BlockIndex: i,
			})
		}
	}

	return blocks
}

// generateOpenAIHeaders generates OpenAI-style cache hints
func (s *PromptCacheService) generateOpenAIHeaders(prompt *StructuredPrompt) *OpenAICacheHints {
	return &OpenAICacheHints{
		StaticPrefix: prompt.CachedTokens,
	}
}

// UpdateCacheStatus updates the cache status for a block
func (s *PromptCacheService) UpdateCacheStatus(hash string, status string) {
	s.staticCacheMu.Lock()
	defer s.staticCacheMu.Unlock()

	if cached, ok := s.staticCache[hash]; ok {
		cached.CacheStatus = status
		cached.LastUsed = time.Now()
		cached.UseCount++
	}
}

// RegisterCachedBlock registers a block as cached
func (s *PromptCacheService) RegisterCachedBlock(block PromptBlock) {
	s.staticCacheMu.Lock()
	defer s.staticCacheMu.Unlock()

	s.staticCache[block.Hash] = &CachedPromptBlock{
		Hash:        block.Hash,
		Content:     block.Content,
		TokenCount:  block.TokenCount,
		CachedAt:    time.Now(),
		LastUsed:    time.Now(),
		UseCount:    1,
		BlockType:   block.Type,
		CacheStatus: "pending",
	}
}

// IsCached checks if a block is already cached
func (s *PromptCacheService) IsCached(hash string) bool {
	s.staticCacheMu.RLock()
	defer s.staticCacheMu.RUnlock()

	cached, ok := s.staticCache[hash]
	if !ok {
		return false
	}

	// Check if cache is still valid
	if time.Since(cached.CachedAt) > s.config.MaxCacheAge {
		return false
	}

	return cached.CacheStatus == "hit"
}

// GetCacheStats returns cache statistics
func (s *PromptCacheService) GetCacheStats() CacheStats {
	s.staticCacheMu.RLock()
	defer s.staticCacheMu.RUnlock()

	stats := CacheStats{
		TotalBlocks: len(s.staticCache),
	}

	for _, block := range s.staticCache {
		stats.TotalTokens += block.TokenCount
		stats.TotalUseCount += block.UseCount

		switch block.CacheStatus {
		case "hit":
			stats.HitCount++
		case "miss":
			stats.MissCount++
		}
	}

	if stats.HitCount+stats.MissCount > 0 {
		stats.HitRate = float64(stats.HitCount) / float64(stats.HitCount+stats.MissCount)
	}

	return stats
}

// CacheStats represents cache statistics
type CacheStats struct {
	TotalBlocks   int     `json:"totalBlocks"`
	TotalTokens   int     `json:"totalTokens"`
	TotalUseCount int     `json:"totalUseCount"`
	HitCount      int     `json:"hitCount"`
	MissCount     int     `json:"missCount"`
	HitRate       float64 `json:"hitRate"`
}

// InvalidateCache invalidates all cached blocks
func (s *PromptCacheService) InvalidateCache() {
	s.staticCacheMu.Lock()
	defer s.staticCacheMu.Unlock()

	s.staticCache = make(map[string]*CachedPromptBlock)
	s.log.Info("Prompt cache invalidated")
}

// InvalidateByType invalidates cached blocks of a specific type
func (s *PromptCacheService) InvalidateByType(blockType BlockType) {
	s.staticCacheMu.Lock()
	defer s.staticCacheMu.Unlock()

	for hash, block := range s.staticCache {
		if block.BlockType == blockType {
			delete(s.staticCache, hash)
		}
	}

	s.log.Info("Prompt cache invalidated for type: " + string(blockType))
}

// CleanupExpired removes expired cache entries
func (s *PromptCacheService) CleanupExpired() int {
	s.staticCacheMu.Lock()
	defer s.staticCacheMu.Unlock()

	removed := 0
	now := time.Now()

	for hash, block := range s.staticCache {
		if now.Sub(block.CachedAt) > s.config.MaxCacheAge {
			delete(s.staticCache, hash)
			removed++
		}
	}

	if removed > 0 {
		s.log.Info("Cleaned up expired prompt cache entries")
	}

	return removed
}

// FormatForProvider formats the structured prompt for a specific provider
func (s *PromptCacheService) FormatForProvider(
	prompt *StructuredPrompt,
	provider string,
) FormattedPrompt {
	var systemParts []string
	var userParts []string

	// Combine cached blocks into system prompt
	for _, block := range prompt.CachedBlocks {
		systemParts = append(systemParts, block.Content)
	}

	// Combine dynamic blocks into user prompt
	for _, block := range prompt.DynamicBlocks {
		if block.Type == BlockTypeUserMessage {
			userParts = append(userParts, block.Content)
		} else {
			// Selected files go before user message
			userParts = append([]string{block.Content}, userParts...)
		}
	}

	formatted := FormattedPrompt{
		SystemPrompt: joinWithSeparator(systemParts, "\n\n"),
		UserPrompt:   joinWithSeparator(userParts, "\n\n"),
		Provider:     provider,
	}

	// Add cache headers
	formatted.CacheHeaders = s.GenerateCacheHeaders(prompt, provider)

	return formatted
}

// FormattedPrompt represents a prompt formatted for a specific provider
type FormattedPrompt struct {
	SystemPrompt string        `json:"systemPrompt"`
	UserPrompt   string        `json:"userPrompt"`
	Provider     string        `json:"provider"`
	CacheHeaders *CacheHeaders `json:"cacheHeaders,omitempty"`
}

// estimateTokens estimates token count for content
func (s *PromptCacheService) estimateTokens(content string) int {
	return len(content) / CharsPerToken
}

// hashContent generates a hash for content
func (s *PromptCacheService) hashContent(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])[:16]
}

// joinWithSeparator joins strings with a separator, skipping empty strings
func joinWithSeparator(parts []string, sep string) string {
	var result string
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i > 0 && result != "" {
			result += sep
		}
		result += part
	}
	return result
}
