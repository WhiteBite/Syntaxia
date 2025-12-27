package ai

import (
	"strings"
	"testing"
	"time"

	"syntaxia/domain"
)

// promptCacheTestLogger implements domain.Logger for prompt cache testing
type promptCacheTestLogger struct{}

func (m *promptCacheTestLogger) Debug(message string)   {}
func (m *promptCacheTestLogger) Info(message string)    {}
func (m *promptCacheTestLogger) Warning(message string) {}
func (m *promptCacheTestLogger) Error(message string)   {}
func (m *promptCacheTestLogger) Fatal(message string)   {}

func newTestService() *PromptCacheService {
	return NewPromptCacheService(&promptCacheTestLogger{}, nil)
}

func newTestServiceWithConfig(config *PromptCacheConfig) *PromptCacheService {
	return NewPromptCacheService(&promptCacheTestLogger{}, config)
}

func TestNewPromptCacheService(t *testing.T) {
	tests := []struct {
		name           string
		config         *PromptCacheConfig
		wantMinTokens  int
		wantMaxAge     time.Duration
		wantRefreshInt time.Duration
	}{
		{
			name:           "nil config uses defaults",
			config:         nil,
			wantMinTokens:  DefaultMinCacheableTokens,
			wantMaxAge:     DefaultMaxCacheAge,
			wantRefreshInt: DefaultRepoMapRefreshMinutes * time.Minute,
		},
		{
			name: "custom config overrides defaults",
			config: &PromptCacheConfig{
				MinCacheableTokens:     512,
				MaxCacheAge:            10 * time.Minute,
				RepoMapRefreshInterval: 20 * time.Minute,
			},
			wantMinTokens:  512,
			wantMaxAge:     10 * time.Minute,
			wantRefreshInt: 20 * time.Minute,
		},
		{
			name: "zero values use defaults",
			config: &PromptCacheConfig{
				MinCacheableTokens: 0,
				MaxCacheAge:        0,
			},
			wantMinTokens:  DefaultMinCacheableTokens,
			wantMaxAge:     DefaultMaxCacheAge,
			wantRefreshInt: DefaultRepoMapRefreshMinutes * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestServiceWithConfig(tt.config)

			if svc.config.MinCacheableTokens != tt.wantMinTokens {
				t.Errorf("MinCacheableTokens = %d, want %d",
					svc.config.MinCacheableTokens, tt.wantMinTokens)
			}
			if svc.config.MaxCacheAge != tt.wantMaxAge {
				t.Errorf("MaxCacheAge = %v, want %v",
					svc.config.MaxCacheAge, tt.wantMaxAge)
			}
			if svc.config.RepoMapRefreshInterval != tt.wantRefreshInt {
				t.Errorf("RepoMapRefreshInterval = %v, want %v",
					svc.config.RepoMapRefreshInterval, tt.wantRefreshInt)
			}
		})
	}
}

func TestBuildStructuredPrompt(t *testing.T) {
	// Create large content to exceed min cacheable tokens
	largeContent := strings.Repeat("x", DefaultMinCacheableTokens*CharsPerToken+100)

	tests := []struct {
		name              string
		input             PromptInput
		wantCachedCount   int
		wantDynamicCount  int
		wantCachedTokens  bool // true if should have cached tokens
		wantDynamicTokens bool // true if should have dynamic tokens
	}{
		{
			name: "empty input",
			input: PromptInput{
				SystemPrompt:  "",
				RepoMap:       "",
				ProjectRules:  "",
				SelectedFiles: nil,
				UserMessage:   "",
			},
			wantCachedCount:   0,
			wantDynamicCount:  0,
			wantCachedTokens:  false,
			wantDynamicTokens: false,
		},
		{
			name: "only system prompt (large)",
			input: PromptInput{
				SystemPrompt: largeContent,
			},
			wantCachedCount:   1,
			wantDynamicCount:  0,
			wantCachedTokens:  true,
			wantDynamicTokens: false,
		},
		{
			name: "all cached blocks (large)",
			input: PromptInput{
				SystemPrompt: largeContent,
				RepoMap:      largeContent,
				ProjectRules: largeContent,
			},
			wantCachedCount:   3,
			wantDynamicCount:  0,
			wantCachedTokens:  true,
			wantDynamicTokens: false,
		},
		{
			name: "only dynamic blocks",
			input: PromptInput{
				SelectedFiles: []SelectedFile{
					{Path: "main.go", Content: "package main"},
				},
				UserMessage: "Fix the bug",
			},
			wantCachedCount:   0,
			wantDynamicCount:  2,
			wantCachedTokens:  false,
			wantDynamicTokens: true,
		},
		{
			name: "mixed cached and dynamic",
			input: PromptInput{
				SystemPrompt: largeContent,
				RepoMap:      largeContent,
				SelectedFiles: []SelectedFile{
					{Path: "main.go", Content: "package main"},
					{Path: "util.go", Content: "package util"},
				},
				UserMessage: "Refactor this code",
			},
			wantCachedCount:   2,
			wantDynamicCount:  3,
			wantCachedTokens:  true,
			wantDynamicTokens: true,
		},
		{
			name: "small system prompt not cached",
			input: PromptInput{
				SystemPrompt: "You are a helpful assistant.",
				UserMessage:  "Hello",
			},
			wantCachedCount:   1, // Block is created but not cacheable
			wantDynamicCount:  1,
			wantCachedTokens:  true,  // Still has tokens
			wantDynamicTokens: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			result := svc.BuildStructuredPrompt(tt.input)

			if len(result.CachedBlocks) != tt.wantCachedCount {
				t.Errorf("CachedBlocks count = %d, want %d",
					len(result.CachedBlocks), tt.wantCachedCount)
			}
			if len(result.DynamicBlocks) != tt.wantDynamicCount {
				t.Errorf("DynamicBlocks count = %d, want %d",
					len(result.DynamicBlocks), tt.wantDynamicCount)
			}
			if tt.wantCachedTokens && result.CachedTokens == 0 {
				t.Error("Expected CachedTokens > 0")
			}
			if tt.wantDynamicTokens && result.DynamicTokens == 0 {
				t.Error("Expected DynamicTokens > 0")
			}
			if result.TotalTokens != result.CachedTokens+result.DynamicTokens {
				t.Errorf("TotalTokens = %d, want %d",
					result.TotalTokens, result.CachedTokens+result.DynamicTokens)
			}
		})
	}
}

func TestBlockTypes(t *testing.T) {
	largeContent := strings.Repeat("x", DefaultMinCacheableTokens*CharsPerToken+100)

	tests := []struct {
		name      string
		input     PromptInput
		wantTypes map[BlockType]int // expected count per type
	}{
		{
			name: "system prompt type",
			input: PromptInput{
				SystemPrompt: largeContent,
			},
			wantTypes: map[BlockType]int{
				BlockTypeSystem: 1,
			},
		},
		{
			name: "repo map type",
			input: PromptInput{
				RepoMap: largeContent,
			},
			wantTypes: map[BlockType]int{
				BlockTypeRepoMap: 1,
			},
		},
		{
			name: "project rules type",
			input: PromptInput{
				ProjectRules: largeContent,
			},
			wantTypes: map[BlockType]int{
				BlockTypeProjectRules: 1,
			},
		},
		{
			name: "selected files type",
			input: PromptInput{
				SelectedFiles: []SelectedFile{
					{Path: "a.go", Content: "package a"},
					{Path: "b.go", Content: "package b"},
				},
			},
			wantTypes: map[BlockType]int{
				BlockTypeSelectedFile: 2,
			},
		},
		{
			name: "user message type",
			input: PromptInput{
				UserMessage: "Hello",
			},
			wantTypes: map[BlockType]int{
				BlockTypeUserMessage: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			result := svc.BuildStructuredPrompt(tt.input)

			// Count types in all blocks
			typeCounts := make(map[BlockType]int)
			for _, block := range result.CachedBlocks {
				typeCounts[block.Type]++
			}
			for _, block := range result.DynamicBlocks {
				typeCounts[block.Type]++
			}

			for wantType, wantCount := range tt.wantTypes {
				if typeCounts[wantType] != wantCount {
					t.Errorf("BlockType %s count = %d, want %d",
						wantType, typeCounts[wantType], wantCount)
				}
			}
		})
	}
}

func TestGenerateCacheHeaders(t *testing.T) {
	largeContent := strings.Repeat("x", DefaultMinCacheableTokens*CharsPerToken+100)

	tests := []struct {
		name                   string
		input                  PromptInput
		provider               string
		wantAnthropicBlocks    int
		wantOpenAIStaticPrefix bool
	}{
		{
			name: "anthropic with cached blocks",
			input: PromptInput{
				SystemPrompt: largeContent,
				RepoMap:      largeContent,
			},
			provider:               "anthropic",
			wantAnthropicBlocks:    2,
			wantOpenAIStaticPrefix: false,
		},
		{
			name: "claude alias for anthropic",
			input: PromptInput{
				SystemPrompt: largeContent,
			},
			provider:               "claude",
			wantAnthropicBlocks:    1,
			wantOpenAIStaticPrefix: false,
		},
		{
			name: "openai with cached blocks",
			input: PromptInput{
				SystemPrompt: largeContent,
				RepoMap:      largeContent,
			},
			provider:               "openai",
			wantAnthropicBlocks:    0,
			wantOpenAIStaticPrefix: true,
		},
		{
			name: "unknown provider",
			input: PromptInput{
				SystemPrompt: largeContent,
			},
			provider:               "unknown",
			wantAnthropicBlocks:    0,
			wantOpenAIStaticPrefix: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			prompt := svc.BuildStructuredPrompt(tt.input)
			headers := svc.GenerateCacheHeaders(prompt, tt.provider)

			if len(headers.AnthropicCacheControl) != tt.wantAnthropicBlocks {
				t.Errorf("AnthropicCacheControl blocks = %d, want %d",
					len(headers.AnthropicCacheControl), tt.wantAnthropicBlocks)
			}

			hasOpenAIPrefix := headers.OpenAICacheHints != nil &&
				headers.OpenAICacheHints.StaticPrefix > 0
			if hasOpenAIPrefix != tt.wantOpenAIStaticPrefix {
				t.Errorf("OpenAI StaticPrefix present = %v, want %v",
					hasOpenAIPrefix, tt.wantOpenAIStaticPrefix)
			}
		})
	}
}

func TestCacheOperations(t *testing.T) {
	t.Run("register and check cached block", func(t *testing.T) {
		svc := newTestService()

		block := PromptBlock{
			Type:       BlockTypeSystem,
			Content:    "test content",
			Cacheable:  true,
			TokenCount: 100,
			Hash:       "testhash123",
		}

		svc.RegisterCachedBlock(block)

		// Initially not "hit" status
		if svc.IsCached(block.Hash) {
			t.Error("Block should not be cached (status is pending)")
		}

		// Update status to hit
		svc.UpdateCacheStatus(block.Hash, "hit")

		if !svc.IsCached(block.Hash) {
			t.Error("Block should be cached after status update")
		}
	})

	t.Run("cache stats", func(t *testing.T) {
		svc := newTestService()

		// Register multiple blocks
		for i := 0; i < 3; i++ {
			block := PromptBlock{
				Type:       BlockTypeSystem,
				Content:    strings.Repeat("x", i*100),
				Cacheable:  true,
				TokenCount: 100 + i*50,
				Hash:       "hash" + string(rune('a'+i)),
			}
			svc.RegisterCachedBlock(block)
		}

		// Update some statuses
		svc.UpdateCacheStatus("hasha", "hit")
		svc.UpdateCacheStatus("hashb", "hit")
		svc.UpdateCacheStatus("hashc", "miss")

		stats := svc.GetCacheStats()

		if stats.TotalBlocks != 3 {
			t.Errorf("TotalBlocks = %d, want 3", stats.TotalBlocks)
		}
		if stats.HitCount != 2 {
			t.Errorf("HitCount = %d, want 2", stats.HitCount)
		}
		if stats.MissCount != 1 {
			t.Errorf("MissCount = %d, want 1", stats.MissCount)
		}
	})

	t.Run("invalidate cache", func(t *testing.T) {
		svc := newTestService()

		block := PromptBlock{
			Type:       BlockTypeSystem,
			Content:    "test",
			Cacheable:  true,
			TokenCount: 100,
			Hash:       "testhash",
		}
		svc.RegisterCachedBlock(block)

		svc.InvalidateCache()

		stats := svc.GetCacheStats()
		if stats.TotalBlocks != 0 {
			t.Errorf("TotalBlocks after invalidate = %d, want 0", stats.TotalBlocks)
		}
	})

	t.Run("invalidate by type", func(t *testing.T) {
		svc := newTestService()

		// Register blocks of different types
		svc.RegisterCachedBlock(PromptBlock{
			Type: BlockTypeSystem, Hash: "sys1", TokenCount: 100,
		})
		svc.RegisterCachedBlock(PromptBlock{
			Type: BlockTypeRepoMap, Hash: "repo1", TokenCount: 100,
		})
		svc.RegisterCachedBlock(PromptBlock{
			Type: BlockTypeSystem, Hash: "sys2", TokenCount: 100,
		})

		svc.InvalidateByType(BlockTypeSystem)

		stats := svc.GetCacheStats()
		if stats.TotalBlocks != 1 {
			t.Errorf("TotalBlocks after type invalidate = %d, want 1", stats.TotalBlocks)
		}
	})
}

func TestFormatForProvider(t *testing.T) {
	largeContent := strings.Repeat("x", DefaultMinCacheableTokens*CharsPerToken+100)

	tests := []struct {
		name             string
		input            PromptInput
		provider         string
		wantSystemNonEmpty bool
		wantUserNonEmpty   bool
	}{
		{
			name: "format with all parts",
			input: PromptInput{
				SystemPrompt: largeContent,
				RepoMap:      "repo map content",
				SelectedFiles: []SelectedFile{
					{Path: "main.go", Content: "package main"},
				},
				UserMessage: "Fix bug",
			},
			provider:         "openai",
			wantSystemNonEmpty: true,
			wantUserNonEmpty:   true,
		},
		{
			name: "format with only system",
			input: PromptInput{
				SystemPrompt: largeContent,
			},
			provider:         "anthropic",
			wantSystemNonEmpty: true,
			wantUserNonEmpty:   false,
		},
		{
			name: "format with only user message",
			input: PromptInput{
				UserMessage: "Hello",
			},
			provider:         "openai",
			wantSystemNonEmpty: false,
			wantUserNonEmpty:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			prompt := svc.BuildStructuredPrompt(tt.input)
			formatted := svc.FormatForProvider(prompt, tt.provider)

			hasSystem := formatted.SystemPrompt != ""
			hasUser := formatted.UserPrompt != ""

			if hasSystem != tt.wantSystemNonEmpty {
				t.Errorf("SystemPrompt non-empty = %v, want %v", hasSystem, tt.wantSystemNonEmpty)
			}
			if hasUser != tt.wantUserNonEmpty {
				t.Errorf("UserPrompt non-empty = %v, want %v", hasUser, tt.wantUserNonEmpty)
			}
			if formatted.Provider != tt.provider {
				t.Errorf("Provider = %s, want %s", formatted.Provider, tt.provider)
			}
		})
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantTokens int
	}{
		{
			name:       "empty content",
			content:    "",
			wantTokens: 0,
		},
		{
			name:       "short content",
			content:    "hello",
			wantTokens: 1, // 5 chars / 4 = 1
		},
		{
			name:       "longer content",
			content:    strings.Repeat("x", 100),
			wantTokens: 25, // 100 / 4 = 25
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			got := svc.estimateTokens(tt.content)

			if got != tt.wantTokens {
				t.Errorf("estimateTokens() = %d, want %d", got, tt.wantTokens)
			}
		})
	}
}

func TestHashContent(t *testing.T) {
	svc := newTestService()

	t.Run("same content same hash", func(t *testing.T) {
		content := "test content"
		hash1 := svc.hashContent(content)
		hash2 := svc.hashContent(content)

		if hash1 != hash2 {
			t.Errorf("Same content produced different hashes: %s vs %s", hash1, hash2)
		}
	})

	t.Run("different content different hash", func(t *testing.T) {
		hash1 := svc.hashContent("content1")
		hash2 := svc.hashContent("content2")

		if hash1 == hash2 {
			t.Error("Different content produced same hash")
		}
	})

	t.Run("hash length is 16", func(t *testing.T) {
		hash := svc.hashContent("test")
		if len(hash) != 16 {
			t.Errorf("Hash length = %d, want 16", len(hash))
		}
	})
}

func TestCleanupExpired(t *testing.T) {
	config := &PromptCacheConfig{
		MinCacheableTokens: 10,
		MaxCacheAge:        100 * time.Millisecond,
	}
	svc := newTestServiceWithConfig(config)

	// Register a block
	block := PromptBlock{
		Type:       BlockTypeSystem,
		Content:    "test",
		Cacheable:  true,
		TokenCount: 100,
		Hash:       "testhash",
	}
	svc.RegisterCachedBlock(block)

	// Should not be cleaned up immediately
	removed := svc.CleanupExpired()
	if removed != 0 {
		t.Errorf("CleanupExpired removed %d blocks, want 0", removed)
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	// Should be cleaned up now
	removed = svc.CleanupExpired()
	if removed != 1 {
		t.Errorf("CleanupExpired removed %d blocks, want 1", removed)
	}

	stats := svc.GetCacheStats()
	if stats.TotalBlocks != 0 {
		t.Errorf("TotalBlocks after cleanup = %d, want 0", stats.TotalBlocks)
	}
}

func TestCacheableThreshold(t *testing.T) {
	config := &PromptCacheConfig{
		MinCacheableTokens: 100,
	}
	svc := newTestServiceWithConfig(config)

	tests := []struct {
		name         string
		tokenCount   int
		wantCacheable bool
	}{
		{
			name:         "below threshold",
			tokenCount:   50,
			wantCacheable: false,
		},
		{
			name:         "at threshold",
			tokenCount:   100,
			wantCacheable: true,
		},
		{
			name:         "above threshold",
			tokenCount:   200,
			wantCacheable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := strings.Repeat("x", tt.tokenCount*CharsPerToken)
			block := svc.createBlock(BlockTypeSystem, content, true)

			if block.Cacheable != tt.wantCacheable {
				t.Errorf("Cacheable = %v, want %v", block.Cacheable, tt.wantCacheable)
			}
		})
	}
}

// Verify interface compliance
var _ domain.Logger = (*promptCacheTestLogger)(nil)
