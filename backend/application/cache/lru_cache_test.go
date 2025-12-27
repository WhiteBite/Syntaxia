package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLRUCache_BasicOperations(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(c *LRUCache[string])
		key      string
		wantVal  string
		wantOK   bool
	}{
		{
			name: "get existing key",
			setup: func(c *LRUCache[string]) {
				c.Set("key1", "value1", 10)
			},
			key:     "key1",
			wantVal: "value1",
			wantOK:  true,
		},
		{
			name:    "get non-existing key",
			setup:   func(c *LRUCache[string]) {},
			key:     "missing",
			wantVal: "",
			wantOK:  false,
		},
		{
			name: "update existing key",
			setup: func(c *LRUCache[string]) {
				c.Set("key1", "value1", 10)
				c.Set("key1", "value2", 10)
			},
			key:     "key1",
			wantVal: "value2",
			wantOK:  true,
		},
		{
			name: "delete key",
			setup: func(c *LRUCache[string]) {
				c.Set("key1", "value1", 10)
				c.Delete("key1")
			},
			key:     "key1",
			wantVal: "",
			wantOK:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[string](Config{Capacity: 10})
			tt.setup(cache)

			val, ok := cache.Get(tt.key)
			assert.Equal(t, tt.wantOK, ok)
			assert.Equal(t, tt.wantVal, val)
		})
	}
}

func TestLRUCache_Eviction(t *testing.T) {
	tests := []struct {
		name          string
		capacity      int
		insertCount   int
		wantLen       int
		wantEvictions int64
	}{
		{
			name:          "no eviction when under capacity",
			capacity:      5,
			insertCount:   3,
			wantLen:       3,
			wantEvictions: 0,
		},
		{
			name:          "eviction when at capacity",
			capacity:      3,
			insertCount:   5,
			wantLen:       3,
			wantEvictions: 2,
		},
		{
			name:          "single item capacity",
			capacity:      1,
			insertCount:   3,
			wantLen:       1,
			wantEvictions: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[int](Config{Capacity: tt.capacity})

			for i := 0; i < tt.insertCount; i++ {
				cache.Set(string(rune('a'+i)), i, 1)
			}

			assert.Equal(t, tt.wantLen, cache.Len())
			assert.Equal(t, tt.wantEvictions, cache.Stats().Evictions)
		})
	}
}

func TestLRUCache_LRUOrder(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 3})

	// Insert 3 items
	cache.Set("a", "1", 1)
	cache.Set("b", "2", 1)
	cache.Set("c", "3", 1)

	// Access "a" to make it most recently used
	cache.Get("a")

	// Insert new item, should evict "b" (least recently used)
	cache.Set("d", "4", 1)

	// "a" should still exist (was accessed)
	val, ok := cache.Get("a")
	assert.True(t, ok)
	assert.Equal(t, "1", val)

	// "b" should be evicted
	_, ok = cache.Get("b")
	assert.False(t, ok)

	// "c" and "d" should exist
	_, ok = cache.Get("c")
	assert.True(t, ok)
	_, ok = cache.Get("d")
	assert.True(t, ok)
}

func TestLRUCache_SizeEviction(t *testing.T) {
	tests := []struct {
		name        string
		maxSize     int64
		items       []struct{ key string; size int64 }
		wantLen     int
		wantSize    int64
	}{
		{
			name:    "evict by size",
			maxSize: 100,
			items: []struct{ key string; size int64 }{
				{"a", 40},
				{"b", 40},
				{"c", 40}, // Should evict "a"
			},
			wantLen:  2,
			wantSize: 80,
		},
		{
			name:    "large item evicts multiple",
			maxSize: 100,
			items: []struct{ key string; size int64 }{
				{"a", 30},
				{"b", 30},
				{"c", 30},
				{"d", 80}, // Should evict a, b, c
			},
			wantLen:  1,
			wantSize: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewLRUCache[string](Config{
				Capacity: 100,
				MaxSize:  tt.maxSize,
			})

			for _, item := range tt.items {
				cache.Set(item.key, item.key, item.size)
			}

			assert.Equal(t, tt.wantLen, cache.Len())
			assert.Equal(t, tt.wantSize, cache.Size())
		})
	}
}

func TestLRUCache_TTL(t *testing.T) {
	cache := NewLRUCache[string](Config{
		Capacity: 10,
		TTL:      50 * time.Millisecond,
	})

	cache.Set("key1", "value1", 1)

	// Should exist immediately
	val, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", val)

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("key1")
	assert.False(t, ok)
}

func TestLRUCache_Contains(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 10})

	cache.Set("key1", "value1", 1)

	assert.True(t, cache.Contains("key1"))
	assert.False(t, cache.Contains("missing"))
}

func TestLRUCache_Clear(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 10})

	cache.Set("key1", "value1", 10)
	cache.Set("key2", "value2", 10)

	assert.Equal(t, 2, cache.Len())
	assert.Equal(t, int64(20), cache.Size())

	cache.Clear()

	assert.Equal(t, 0, cache.Len())
	assert.Equal(t, int64(0), cache.Size())
}

func TestLRUCache_Stats(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 10, MaxSize: 1000})

	cache.Set("key1", "value1", 10)
	cache.Get("key1") // hit
	cache.Get("key2") // miss

	stats := cache.Stats()

	assert.Equal(t, 1, stats.Items)
	assert.Equal(t, 10, stats.Capacity)
	assert.Equal(t, int64(10), stats.CurrentSize)
	assert.Equal(t, int64(1000), stats.MaxSize)
	assert.Equal(t, int64(1), stats.Hits)
	assert.Equal(t, int64(1), stats.Misses)
	assert.Equal(t, 0.5, stats.HitRate)
}

func TestLRUCache_GetOrSet(t *testing.T) {
	cache := NewLRUCache[int](Config{Capacity: 10})

	computeCalls := 0
	compute := func() (int, int64, error) {
		computeCalls++
		return 42, 8, nil
	}

	// First call should compute
	val, fromCache, err := cache.GetOrSet("key1", compute)
	require.NoError(t, err)
	assert.Equal(t, 42, val)
	assert.False(t, fromCache)
	assert.Equal(t, 1, computeCalls)

	// Second call should use cache
	val, fromCache, err = cache.GetOrSet("key1", compute)
	require.NoError(t, err)
	assert.Equal(t, 42, val)
	assert.True(t, fromCache)
	assert.Equal(t, 1, computeCalls) // Not incremented
}

func TestLRUCache_Cleanup(t *testing.T) {
	cache := NewLRUCache[string](Config{
		Capacity: 10,
		TTL:      50 * time.Millisecond,
	})

	cache.Set("key1", "value1", 1)
	cache.Set("key2", "value2", 1)

	// Wait for expiration
	time.Sleep(60 * time.Millisecond)

	// Add non-expired item
	cache.Set("key3", "value3", 1)

	removed := cache.Cleanup()

	assert.Equal(t, 2, removed)
	assert.Equal(t, 1, cache.Len())
	assert.True(t, cache.Contains("key3"))
}

func TestLRUCache_Keys(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 10})

	cache.Set("a", "1", 1)
	cache.Set("b", "2", 1)
	cache.Set("c", "3", 1)

	keys := cache.Keys()

	assert.Len(t, keys, 3)
	assert.Contains(t, keys, "a")
	assert.Contains(t, keys, "b")
	assert.Contains(t, keys, "c")
}

func TestLRUCache_Concurrent(t *testing.T) {
	cache := NewLRUCache[int](Config{Capacity: 100})

	var wg sync.WaitGroup
	numGoroutines := 10
	numOps := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := string(rune('a' + (id*numOps+j)%26))
				cache.Set(key, id*numOps+j, 1)
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := string(rune('a' + j%26))
				cache.Get(key)
			}
		}()
	}

	wg.Wait()

	// Should not panic and maintain consistency
	assert.LessOrEqual(t, cache.Len(), 100)
}

func TestLRUCache_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.Equal(t, 1000, cfg.Capacity)
	assert.Equal(t, int64(100*1024*1024), cfg.MaxSize)
	assert.Equal(t, 30*time.Minute, cfg.TTL)
}

func TestLRUCache_ZeroCapacity(t *testing.T) {
	// Should use default capacity
	cache := NewLRUCache[string](Config{Capacity: 0})

	cache.Set("key1", "value1", 1)
	val, ok := cache.Get("key1")

	assert.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestLRUCache_UpdateSize(t *testing.T) {
	cache := NewLRUCache[string](Config{Capacity: 10, MaxSize: 100})

	cache.Set("key1", "value1", 30)
	assert.Equal(t, int64(30), cache.Size())

	// Update with different size
	cache.Set("key1", "value2", 50)
	assert.Equal(t, int64(50), cache.Size())
	assert.Equal(t, 1, cache.Len())
}
