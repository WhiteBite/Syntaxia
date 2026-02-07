package git

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestGitCache_GetSet(t *testing.T) {
	cache := NewGitCache()

	// Set value
	cache.Set("test-key", "test-value", 1*time.Minute)

	// Get value
	val, ok := cache.Get("test-key")
	if !ok {
		t.Fatal("Expected cache hit")
	}
	if val.(string) != "test-value" {
		t.Errorf("Expected 'test-value', got %v", val)
	}
}

func TestGitCache_GetMiss(t *testing.T) {
	cache := NewGitCache()

	// Get non-existent key
	_, ok := cache.Get("non-existent")
	if ok {
		t.Error("Expected cache miss for non-existent key")
	}

	// Check stats
	stats := cache.GetStats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}
}

func TestGitCache_Expiration(t *testing.T) {
	cache := NewGitCache()

	// Set with short TTL
	cache.Set("test-key", "test-value", 50*time.Millisecond)

	// Should be available immediately
	_, ok := cache.Get("test-key")
	if !ok {
		t.Error("Expected cache hit before expiration")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	_, ok = cache.Get("test-key")
	if ok {
		t.Error("Expected cache miss after expiration")
	}
}

func TestGitCache_Invalidate(t *testing.T) {
	cache := NewGitCache()

	// Set multiple entries
	cache.Set("branches:project1", []string{"main"}, 1*time.Minute)
	cache.Set("branches:project2", []string{"dev"}, 1*time.Minute)
	cache.Set("file:project1:main:file.go", "content", 1*time.Minute)
	cache.Set("commits:project1:50", []string{"abc123"}, 1*time.Minute)

	// Invalidate all branches
	cache.Invalidate("branches:")

	// Check branches are invalidated
	_, ok1 := cache.Get("branches:project1")
	_, ok2 := cache.Get("branches:project2")
	if ok1 || ok2 {
		t.Error("Expected branches to be invalidated")
	}

	// Check file cache remains
	_, ok3 := cache.Get("file:project1:main:file.go")
	if !ok3 {
		t.Error("Expected file cache to remain")
	}

	// Check commits remain
	_, ok4 := cache.Get("commits:project1:50")
	if !ok4 {
		t.Error("Expected commits cache to remain")
	}
}

func TestGitCache_InvalidateProject(t *testing.T) {
	cache := NewGitCache()

	// Set entries for multiple projects
	cache.Set("branches:project1", []string{"main"}, 1*time.Minute)
	cache.Set("file:project1:main:file.go", "content1", 1*time.Minute)
	cache.Set("branches:project2", []string{"dev"}, 1*time.Minute)
	cache.Set("file:project2:dev:file.go", "content2", 1*time.Minute)

	// Invalidate project1
	cache.InvalidateProject("project1")

	// Check project1 entries are gone
	_, ok1 := cache.Get("branches:project1")
	_, ok2 := cache.Get("file:project1:main:file.go")
	if ok1 || ok2 {
		t.Error("Expected project1 entries to be invalidated")
	}

	// Check project2 entries remain
	_, ok3 := cache.Get("branches:project2")
	_, ok4 := cache.Get("file:project2:dev:file.go")
	if !ok3 || !ok4 {
		t.Error("Expected project2 entries to remain")
	}
}

func TestGitCache_Stats(t *testing.T) {
	cache := NewGitCache()

	cache.Set("key1", "val1", 1*time.Minute)
	cache.Set("key2", "val2", 1*time.Minute)

	// Hit
	cache.Get("key1")
	cache.Get("key1")

	// Miss
	cache.Get("key3")
	cache.Get("key4")

	stats := cache.GetStats()
	if stats.Hits != 2 {
		t.Errorf("Expected 2 hits, got %d", stats.Hits)
	}
	if stats.Misses != 2 {
		t.Errorf("Expected 2 misses, got %d", stats.Misses)
	}
	if stats.Size != 2 {
		t.Errorf("Expected size 2, got %d", stats.Size)
	}
}

func TestGitCache_Concurrent(t *testing.T) {
	cache := NewGitCache()

	var wg sync.WaitGroup
	iterations := 100

	// Concurrent writes
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			cache.Set(key, n, 1*time.Minute)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", n)
			cache.Get(key)
		}(i)
	}

	wg.Wait()

	// Verify no race conditions occurred
	stats := cache.GetStats()
	if stats.Size > int64(iterations) {
		t.Errorf("Expected size <= %d, got %d", iterations, stats.Size)
	}
}

func TestGitCache_Clear(t *testing.T) {
	cache := NewGitCache()

	// Add entries
	cache.Set("key1", "val1", 1*time.Minute)
	cache.Set("key2", "val2", 1*time.Minute)
	cache.Set("key3", "val3", 1*time.Minute)

	// Clear
	cache.Clear()

	// Verify empty
	if cache.GetSize() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", cache.GetSize())
	}

	stats := cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 || stats.Size != 0 {
		t.Error("Expected all stats to be reset after clear")
	}
}

func TestGitCache_CleanExpired(t *testing.T) {
	cache := NewGitCache()

	// Add entries with different TTLs
	cache.Set("short1", "val1", 50*time.Millisecond)
	cache.Set("short2", "val2", 50*time.Millisecond)
	cache.Set("long1", "val3", 1*time.Minute)
	cache.Set("long2", "val4", 1*time.Minute)

	// Wait for short TTL entries to expire
	time.Sleep(100 * time.Millisecond)

	// Clean expired
	removed := cache.CleanExpired()

	if removed != 2 {
		t.Errorf("Expected 2 entries removed, got %d", removed)
	}

	if cache.GetSize() != 2 {
		t.Errorf("Expected 2 entries remaining, got %d", cache.GetSize())
	}

	// Verify long TTL entries still exist
	_, ok1 := cache.Get("long1")
	_, ok2 := cache.Get("long2")
	if !ok1 || !ok2 {
		t.Error("Expected long TTL entries to remain")
	}
}

func TestGitCache_Has(t *testing.T) {
	cache := NewGitCache()

	// Non-existent key
	if cache.Has("non-existent") {
		t.Error("Expected Has to return false for non-existent key")
	}

	// Add key
	cache.Set("test-key", "value", 1*time.Minute)

	// Should exist
	if !cache.Has("test-key") {
		t.Error("Expected Has to return true for existing key")
	}

	// Add expired key
	cache.Set("expired-key", "value", 10*time.Millisecond)
	time.Sleep(50 * time.Millisecond)

	// Should not exist (expired)
	if cache.Has("expired-key") {
		t.Error("Expected Has to return false for expired key")
	}
}

func TestGitCache_MultipleDataTypes(t *testing.T) {
	cache := NewGitCache()

	// Test different data types
	cache.Set("string", "test", 1*time.Minute)
	cache.Set("int", 42, 1*time.Minute)
	cache.Set("slice", []string{"a", "b", "c"}, 1*time.Minute)
	cache.Set("map", map[string]int{"a": 1, "b": 2}, 1*time.Minute)

	// Retrieve and verify
	if val, ok := cache.Get("string"); !ok || val.(string) != "test" {
		t.Error("String value mismatch")
	}

	if val, ok := cache.Get("int"); !ok || val.(int) != 42 {
		t.Error("Int value mismatch")
	}

	if val, ok := cache.Get("slice"); !ok || len(val.([]string)) != 3 {
		t.Error("Slice value mismatch")
	}

	if val, ok := cache.Get("map"); !ok || len(val.(map[string]int)) != 2 {
		t.Error("Map value mismatch")
	}
}

func TestGitCache_CacheKeyBuilders(t *testing.T) {
	tests := []struct {
		name     string
		builder  func() string
		expected string
	}{
		{
			name:     "branches key",
			builder:  func() string { return buildBranchesCacheKey("/path/to/project") },
			expected: "branches:/path/to/project",
		},
		{
			name:     "current branch key",
			builder:  func() string { return buildCurrentBranchCacheKey("/path/to/project") },
			expected: "current-branch:/path/to/project",
		},
		{
			name:     "files at ref key",
			builder:  func() string { return buildFilesAtRefCacheKey("/path/to/project", "main") },
			expected: "files:/path/to/project:main",
		},
		{
			name:     "file at ref key",
			builder:  func() string { return buildFileAtRefCacheKey("/path/to/project", "file.go", "abc123") },
			expected: "file:/path/to/project:abc123:file.go",
		},
		{
			name:     "commit history key",
			builder:  func() string { return buildCommitHistoryCacheKey("/path/to/project", 50) },
			expected: "commits:/path/to/project:50",
		},
		{
			name:     "status key",
			builder:  func() string { return buildStatusCacheKey("/path/to/project") },
			expected: "status:/path/to/project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.builder()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGitCache_OverwriteExisting(t *testing.T) {
	cache := NewGitCache()

	// Set initial value
	cache.Set("key", "value1", 1*time.Minute)

	// Overwrite with new value
	cache.Set("key", "value2", 1*time.Minute)

	// Should get new value
	val, ok := cache.Get("key")
	if !ok {
		t.Fatal("Expected cache hit")
	}
	if val.(string) != "value2" {
		t.Errorf("Expected 'value2', got %v", val)
	}

	// Size should still be 1
	if cache.GetSize() != 1 {
		t.Errorf("Expected size 1, got %d", cache.GetSize())
	}
}

func TestGitCache_LargeDataset(t *testing.T) {
	cache := NewGitCache()

	// Add many entries
	count := 1000
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, i, 1*time.Minute)
	}

	// Verify size
	if cache.GetSize() != count {
		t.Errorf("Expected size %d, got %d", count, cache.GetSize())
	}

	// Verify random access
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i*10)
		val, ok := cache.Get(key)
		if !ok {
			t.Errorf("Expected to find key %s", key)
		}
		if val.(int) != i*10 {
			t.Errorf("Expected value %d, got %v", i*10, val)
		}
	}
}
