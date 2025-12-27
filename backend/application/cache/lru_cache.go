// Package cache provides caching utilities for analysis results.
package cache

import (
	"container/list"
	"sync"
	"time"
)

// Entry represents a cache entry with metadata.
type Entry[V any] struct {
	Key       string
	Value     V
	Size      int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

// LRUCache is a thread-safe generic LRU cache with TTL support.
type LRUCache[V any] struct {
	mu       sync.RWMutex
	capacity int
	maxSize  int64
	ttl      time.Duration

	items    map[string]*list.Element
	eviction *list.List

	// Metrics
	currentSize int64
	hits        int64
	misses      int64
	evictions   int64
}

// Config holds LRU cache configuration.
type Config struct {
	Capacity int           // Maximum number of items
	MaxSize  int64         // Maximum total size in bytes (0 = unlimited)
	TTL      time.Duration // Time-to-live for entries (0 = no expiration)
}

// DefaultConfig returns default cache configuration.
func DefaultConfig() Config {
	return Config{
		Capacity: 1000,
		MaxSize:  100 * 1024 * 1024, // 100MB
		TTL:      30 * time.Minute,
	}
}

// NewLRUCache creates a new LRU cache with the given configuration.
func NewLRUCache[V any](cfg Config) *LRUCache[V] {
	if cfg.Capacity <= 0 {
		cfg.Capacity = 1000
	}
	return &LRUCache[V]{
		capacity: cfg.Capacity,
		maxSize:  cfg.MaxSize,
		ttl:      cfg.TTL,
		items:    make(map[string]*list.Element),
		eviction: list.New(),
	}
}

// Get retrieves a value from the cache.
// Returns the value and true if found and not expired, otherwise zero value and false.
func (c *LRUCache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	elem, exists := c.items[key]
	if !exists {
		c.misses++
		return zero, false
	}

	entry := elem.Value.(*Entry[V])

	// Check TTL expiration
	if c.ttl > 0 && time.Now().After(entry.ExpiresAt) {
		c.removeLocked(key)
		c.misses++
		return zero, false
	}

	// Move to front (most recently used)
	c.eviction.MoveToFront(elem)
	c.hits++

	return entry.Value, true
}

// Set adds or updates a value in the cache.
func (c *LRUCache[V]) Set(key string, value V, size int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.setLocked(key, value, size)
}

func (c *LRUCache[V]) setLocked(key string, value V, size int64) {
	now := time.Now()
	expiresAt := time.Time{}
	if c.ttl > 0 {
		expiresAt = now.Add(c.ttl)
	}

	// Update existing entry
	if elem, exists := c.items[key]; exists {
		oldEntry := elem.Value.(*Entry[V])
		c.currentSize -= oldEntry.Size

		oldEntry.Value = value
		oldEntry.Size = size
		oldEntry.CreatedAt = now
		oldEntry.ExpiresAt = expiresAt

		c.currentSize += size
		c.eviction.MoveToFront(elem)
		return
	}

	// Evict if necessary before adding new entry
	c.evictIfNeededLocked(size)

	entry := &Entry[V]{
		Key:       key,
		Value:     value,
		Size:      size,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	elem := c.eviction.PushFront(entry)
	c.items[key] = elem
	c.currentSize += size
}

// Delete removes a key from the cache.
func (c *LRUCache[V]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.removeLocked(key)
}

func (c *LRUCache[V]) removeLocked(key string) {
	elem, exists := c.items[key]
	if !exists {
		return
	}

	entry := elem.Value.(*Entry[V])
	c.currentSize -= entry.Size
	c.eviction.Remove(elem)
	delete(c.items, key)
}

// evictIfNeededLocked evicts entries until there's room for newSize.
func (c *LRUCache[V]) evictIfNeededLocked(newSize int64) {
	// Evict by capacity
	for c.eviction.Len() >= c.capacity {
		c.evictOldestLocked()
	}

	// Evict by size
	if c.maxSize > 0 {
		for c.currentSize+newSize > c.maxSize && c.eviction.Len() > 0 {
			c.evictOldestLocked()
		}
	}
}

func (c *LRUCache[V]) evictOldestLocked() {
	elem := c.eviction.Back()
	if elem == nil {
		return
	}

	entry := elem.Value.(*Entry[V])
	c.currentSize -= entry.Size
	c.eviction.Remove(elem)
	delete(c.items, entry.Key)
	c.evictions++
}

// Clear removes all entries from the cache.
func (c *LRUCache[V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.eviction = list.New()
	c.currentSize = 0
}

// Len returns the number of items in the cache.
func (c *LRUCache[V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.eviction.Len()
}

// Size returns the current total size of cached items.
func (c *LRUCache[V]) Size() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.currentSize
}

// Stats returns cache statistics.
func (c *LRUCache[V]) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hitRate := float64(0)
	total := c.hits + c.misses
	if total > 0 {
		hitRate = float64(c.hits) / float64(total)
	}

	return CacheStats{
		Items:       c.eviction.Len(),
		Capacity:    c.capacity,
		CurrentSize: c.currentSize,
		MaxSize:     c.maxSize,
		Hits:        c.hits,
		Misses:      c.misses,
		Evictions:   c.evictions,
		HitRate:     hitRate,
	}
}

// CacheStats holds cache statistics.
type CacheStats struct {
	Items       int     `json:"items"`
	Capacity    int     `json:"capacity"`
	CurrentSize int64   `json:"currentSize"`
	MaxSize     int64   `json:"maxSize"`
	Hits        int64   `json:"hits"`
	Misses      int64   `json:"misses"`
	Evictions   int64   `json:"evictions"`
	HitRate     float64 `json:"hitRate"`
}

// Keys returns all keys in the cache (for debugging/testing).
func (c *LRUCache[V]) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

// Contains checks if a key exists in the cache without updating access time.
func (c *LRUCache[V]) Contains(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, exists := c.items[key]
	if !exists {
		return false
	}

	// Check TTL expiration
	if c.ttl > 0 {
		entry := elem.Value.(*Entry[V])
		if time.Now().After(entry.ExpiresAt) {
			return false
		}
	}

	return true
}

// GetOrSet retrieves a value from cache or computes and stores it if not present.
// The compute function is called only if the key is not in cache.
// Returns the value and true if it was from cache, false if computed.
func (c *LRUCache[V]) GetOrSet(key string, compute func() (V, int64, error)) (V, bool, error) {
	// Try to get from cache first
	if value, ok := c.Get(key); ok {
		return value, true, nil
	}

	// Compute the value
	value, size, err := compute()
	if err != nil {
		var zero V
		return zero, false, err
	}

	// Store in cache
	c.Set(key, value, size)
	return value, false, nil
}

// Cleanup removes expired entries from the cache.
// Call this periodically if TTL is enabled.
func (c *LRUCache[V]) Cleanup() int {
	if c.ttl == 0 {
		return 0
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	removed := 0

	for key, elem := range c.items {
		entry := elem.Value.(*Entry[V])
		if now.After(entry.ExpiresAt) {
			c.currentSize -= entry.Size
			c.eviction.Remove(elem)
			delete(c.items, key)
			removed++
		}
	}

	return removed
}
