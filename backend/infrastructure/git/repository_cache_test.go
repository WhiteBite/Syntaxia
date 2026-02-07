package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"testing"
	"time"
)

// TestRepository_GetBranches_Cached tests that GetBranches uses cache
func TestRepository_GetBranches_Cached(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	// Create temp git repo
	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call - cache miss
	branches1, err := repo.GetBranches(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get branches: %v", err)
	}

	stats1 := repo.cache.GetStats()
	if stats1.Misses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", stats1.Misses)
	}

	// Second call - should hit cache
	branches2, err := repo.GetBranches(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get branches: %v", err)
	}

	stats2 := repo.cache.GetStats()
	if stats2.Hits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", stats2.Hits)
	}

	// Results should be identical
	if len(branches1) != len(branches2) {
		t.Errorf("Branch lists differ: %d vs %d", len(branches1), len(branches2))
	}
}

// TestRepository_GetCurrentBranch_Cached tests that GetCurrentBranch uses cache
func TestRepository_GetCurrentBranch_Cached(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call - cache miss
	branch1, err := repo.GetCurrentBranch(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	stats1 := repo.cache.GetStats()
	if stats1.Misses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", stats1.Misses)
	}

	// Second call - should hit cache
	branch2, err := repo.GetCurrentBranch(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get current branch: %v", err)
	}

	stats2 := repo.cache.GetStats()
	if stats2.Hits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", stats2.Hits)
	}

	// Results should be identical
	if branch1 != branch2 {
		t.Errorf("Branch names differ: %s vs %s", branch1, branch2)
	}
}

// TestRepository_ListFilesAtRef_Cached tests that ListFilesAtRef uses cache
func TestRepository_ListFilesAtRef_Cached(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call - cache miss
	files1, err := repo.ListFilesAtRef(tmpDir, "HEAD")
	if err != nil {
		t.Fatalf("Failed to list files: %v", err)
	}

	stats1 := repo.cache.GetStats()
	if stats1.Misses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", stats1.Misses)
	}

	// Second call - should hit cache
	files2, err := repo.ListFilesAtRef(tmpDir, "HEAD")
	if err != nil {
		t.Fatalf("Failed to list files: %v", err)
	}

	stats2 := repo.cache.GetStats()
	if stats2.Hits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", stats2.Hits)
	}

	// Results should be identical
	if len(files1) != len(files2) {
		t.Errorf("File lists differ: %d vs %d", len(files1), len(files2))
	}
}

// TestRepository_GetFileAtRef_Cached tests that GetFileAtRef uses cache (immutable)
func TestRepository_GetFileAtRef_Cached(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call - cache miss
	content1, err := repo.GetFileAtRef(tmpDir, "test.txt", "HEAD")
	if err != nil {
		t.Fatalf("Failed to get file: %v", err)
	}

	stats1 := repo.cache.GetStats()
	if stats1.Misses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", stats1.Misses)
	}

	// Second call - should hit cache
	content2, err := repo.GetFileAtRef(tmpDir, "test.txt", "HEAD")
	if err != nil {
		t.Fatalf("Failed to get file: %v", err)
	}

	stats2 := repo.cache.GetStats()
	if stats2.Hits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", stats2.Hits)
	}

	// Content should be identical
	if content1 != content2 {
		t.Error("File contents differ")
	}
}

// TestRepository_GetCommitHistory_Cached tests that GetCommitHistory uses cache
func TestRepository_GetCommitHistory_Cached(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call - cache miss
	commits1, err := repo.GetCommitHistory(tmpDir, 10)
	if err != nil {
		t.Fatalf("Failed to get commit history: %v", err)
	}

	stats1 := repo.cache.GetStats()
	if stats1.Misses != 1 {
		t.Errorf("Expected 1 cache miss, got %d", stats1.Misses)
	}

	// Second call - should hit cache
	commits2, err := repo.GetCommitHistory(tmpDir, 10)
	if err != nil {
		t.Fatalf("Failed to get commit history: %v", err)
	}

	stats2 := repo.cache.GetStats()
	if stats2.Hits != 1 {
		t.Errorf("Expected 1 cache hit, got %d", stats2.Hits)
	}

	// Results should be identical
	if len(commits1) != len(commits2) {
		t.Errorf("Commit lists differ: %d vs %d", len(commits1), len(commits2))
	}
}

// TestRepository_ClearCache tests cache clearing
func TestRepository_ClearCache(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Populate cache
	_, _ = repo.GetBranches(tmpDir)
	_, _ = repo.GetCurrentBranch(tmpDir)

	// Verify cache has entries
	if repo.cache.GetSize() == 0 {
		t.Error("Expected cache to have entries")
	}

	// Clear cache
	repo.ClearCache()

	// Verify cache is empty
	if repo.cache.GetSize() != 0 {
		t.Errorf("Expected cache to be empty, got size %d", repo.cache.GetSize())
	}
}

// TestRepository_InvalidateProjectCache tests project-specific cache invalidation
func TestRepository_InvalidateProjectCache(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir1 := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir1)

	tmpDir2 := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir2)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Populate cache for both projects
	_, _ = repo.GetBranches(tmpDir1)
	_, _ = repo.GetBranches(tmpDir2)

	initialSize := repo.cache.GetSize()
	if initialSize != 2 {
		t.Errorf("Expected cache size 2, got %d", initialSize)
	}

	// Invalidate project1
	repo.InvalidateProjectCache(tmpDir1)

	// Verify project1 cache is gone but project2 remains
	if repo.cache.Has(buildBranchesCacheKey(tmpDir1)) {
		t.Error("Expected project1 cache to be invalidated")
	}

	if !repo.cache.Has(buildBranchesCacheKey(tmpDir2)) {
		t.Error("Expected project2 cache to remain")
	}
}

// TestRepository_GetCacheStats tests cache statistics
func TestRepository_GetCacheStats(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Generate some cache activity
	_, _ = repo.GetBranches(tmpDir)      // miss
	_, _ = repo.GetBranches(tmpDir)      // hit
	_, _ = repo.GetCurrentBranch(tmpDir) // miss
	_, _ = repo.GetCurrentBranch(tmpDir) // hit

	stats := repo.GetCacheStats()

	hits, ok := stats["hits"].(int64)
	if !ok || hits != 2 {
		t.Errorf("Expected 2 hits, got %v", stats["hits"])
	}

	misses, ok := stats["misses"].(int64)
	if !ok || misses != 2 {
		t.Errorf("Expected 2 misses, got %v", stats["misses"])
	}

	hitRate, ok := stats["hitRate"].(float64)
	if !ok || hitRate != 50.0 {
		t.Errorf("Expected 50%% hit rate, got %v", stats["hitRate"])
	}
}

// TestRepository_CacheExpiration tests that cache entries expire correctly
func TestRepository_CacheExpiration(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Manually set a short TTL entry
	cacheKey := buildBranchesCacheKey(tmpDir)
	repo.cache.Set(cacheKey, []string{"main"}, 50*time.Millisecond)

	// Should be available immediately
	if !repo.cache.Has(cacheKey) {
		t.Error("Expected cache entry to exist")
	}

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	if repo.cache.Has(cacheKey) {
		t.Error("Expected cache entry to be expired")
	}
}

// TestRepository_CachePerformance tests cache performance improvement
func TestRepository_CachePerformance(t *testing.T) {
	if !isGitAvailable() {
		t.Skip("Git not available")
	}

	tmpDir := createTempGitRepo(t)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// First call (cache miss) - measure time
	start1 := time.Now()
	_, err := repo.GetBranches(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get branches: %v", err)
	}
	duration1 := time.Since(start1)

	// Second call (cache hit) - measure time
	start2 := time.Now()
	_, err = repo.GetBranches(tmpDir)
	if err != nil {
		t.Fatalf("Failed to get branches: %v", err)
	}
	duration2 := time.Since(start2)

	// Cache hit should be significantly faster (at least 10x)
	if duration2 > duration1/10 {
		t.Logf("Warning: Cache hit (%v) not significantly faster than miss (%v)", duration2, duration1)
	} else {
		t.Logf("Cache performance: miss=%v, hit=%v (%.1fx faster)", duration1, duration2, float64(duration1)/float64(duration2))
	}
}

// Helper functions

func isGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func createTempGitRepo(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "git-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to add files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to commit: %v", err)
	}

	return tmpDir
}
