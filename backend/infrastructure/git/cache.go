package git

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// CacheEntry represents a single cache entry with expiration
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// CacheStats tracks cache performance metrics
type CacheStats struct {
	Hits   int64
	Misses int64
	Size   int64
}

// GitCache provides thread-safe caching for Git operations with TTL support
type GitCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheEntry
	stats CacheStats
}

// NewGitCache creates a new Git cache instance
func NewGitCache() *GitCache {
	return &GitCache{
		cache: make(map[string]*CacheEntry),
	}
}

// Get retrieves a value from cache if it exists and hasn't expired
func (c *GitCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		c.stats.Misses++
		return nil, false
	}

	// Check expiration
	if time.Now().After(entry.ExpiresAt) {
		c.stats.Misses++
		return nil, false
	}

	c.stats.Hits++
	return entry.Data, true
}

// Set stores a value in cache with the specified TTL
func (c *GitCache) Set(key string, data interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.stats.Size = int64(len(c.cache))
}

// Invalidate removes all cache entries matching the pattern
func (c *GitCache) Invalidate(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key := range c.cache {
		if strings.Contains(key, pattern) {
			delete(c.cache, key)
		}
	}
	c.stats.Size = int64(len(c.cache))
}

// InvalidateProject removes all cache entries for a specific project
func (c *GitCache) InvalidateProject(projectPath string) {
	c.Invalidate(projectPath)
}

// GetStats returns current cache statistics
func (c *GitCache) GetStats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := c.stats
	stats.Size = int64(len(c.cache))
	return stats
}

// Clear removes all entries from the cache
func (c *GitCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*CacheEntry)
	c.stats = CacheStats{}
}

// CleanExpired removes all expired entries from the cache
func (c *GitCache) CleanExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removed := 0

	for key, entry := range c.cache {
		if now.After(entry.ExpiresAt) {
			delete(c.cache, key)
			removed++
		}
	}

	c.stats.Size = int64(len(c.cache))
	return removed
}

// GetSize returns the current number of entries in the cache
func (c *GitCache) GetSize() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// Has checks if a key exists in cache and hasn't expired
func (c *GitCache) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return false
	}

	return time.Now().Before(entry.ExpiresAt)
}

// buildCacheKey creates a cache key for branches
func buildBranchesCacheKey(projectPath string) string {
	return fmt.Sprintf("branches:%s", projectPath)
}

// buildCurrentBranchCacheKey creates a cache key for current branch
func buildCurrentBranchCacheKey(projectPath string) string {
	return fmt.Sprintf("current-branch:%s", projectPath)
}

// buildFilesAtRefCacheKey creates a cache key for files at ref
func buildFilesAtRefCacheKey(projectPath, ref string) string {
	return fmt.Sprintf("files:%s:%s", projectPath, ref)
}

// buildFileAtRefCacheKey creates a cache key for file content at ref
func buildFileAtRefCacheKey(projectPath, filePath, ref string) string {
	return fmt.Sprintf("file:%s:%s:%s", projectPath, ref, filePath)
}

// buildCommitHistoryCacheKey creates a cache key for commit history
func buildCommitHistoryCacheKey(projectPath string, limit int) string {
	return fmt.Sprintf("commits:%s:%d", projectPath, limit)
}

// buildStatusCacheKey creates a cache key for git status
func buildStatusCacheKey(projectPath string) string {
	return fmt.Sprintf("status:%s", projectPath)
}
