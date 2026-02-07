package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"testing"
)

// BenchmarkGitCache_Set benchmarks cache set operations
func BenchmarkGitCache_Set(b *testing.B) {
	cache := NewGitCache()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, i, 1)
	}
}

// BenchmarkGitCache_Get benchmarks cache get operations
func BenchmarkGitCache_Get(b *testing.B) {
	cache := NewGitCache()

	// Pre-populate
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Set(key, i, 1)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		cache.Get(key)
	}
}

// BenchmarkGitCache_ConcurrentReadWrite benchmarks concurrent operations
func BenchmarkGitCache_ConcurrentReadWrite(b *testing.B) {
	cache := NewGitCache()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key-%d", i%100)
			if i%2 == 0 {
				cache.Set(key, i, 1)
			} else {
				cache.Get(key)
			}
			i++
		}
	})
}

// BenchmarkRepository_GetBranches_NoCache benchmarks without cache
func BenchmarkRepository_GetBranches_NoCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := New(&domain.NoopLogger{}).(*Repository)
		_, _ = repo.GetBranches(tmpDir)
	}
}

// BenchmarkRepository_GetBranches_WithCache benchmarks with cache
func BenchmarkRepository_GetBranches_WithCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Warm up cache
	_, _ = repo.GetBranches(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetBranches(tmpDir)
	}
}

// BenchmarkRepository_GetFileAtRef_NoCache benchmarks file retrieval without cache
func BenchmarkRepository_GetFileAtRef_NoCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := New(&domain.NoopLogger{}).(*Repository)
		_, _ = repo.GetFileAtRef(tmpDir, "test.txt", "HEAD")
	}
}

// BenchmarkRepository_GetFileAtRef_WithCache benchmarks file retrieval with cache
func BenchmarkRepository_GetFileAtRef_WithCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Warm up cache
	_, _ = repo.GetFileAtRef(tmpDir, "test.txt", "HEAD")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetFileAtRef(tmpDir, "test.txt", "HEAD")
	}
}

// BenchmarkRepository_ListFilesAtRef_NoCache benchmarks file listing without cache
func BenchmarkRepository_ListFilesAtRef_NoCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := New(&domain.NoopLogger{}).(*Repository)
		_, _ = repo.ListFilesAtRef(tmpDir, "HEAD")
	}
}

// BenchmarkRepository_ListFilesAtRef_WithCache benchmarks file listing with cache
func BenchmarkRepository_ListFilesAtRef_WithCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Warm up cache
	_, _ = repo.ListFilesAtRef(tmpDir, "HEAD")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.ListFilesAtRef(tmpDir, "HEAD")
	}
}

// BenchmarkRepository_GetCommitHistory_NoCache benchmarks commit history without cache
func BenchmarkRepository_GetCommitHistory_NoCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := New(&domain.NoopLogger{}).(*Repository)
		_, _ = repo.GetCommitHistory(tmpDir, 10)
	}
}

// BenchmarkRepository_GetCommitHistory_WithCache benchmarks commit history with cache
func BenchmarkRepository_GetCommitHistory_WithCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepo(b)
	defer os.RemoveAll(tmpDir)

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Warm up cache
	_, _ = repo.GetCommitHistory(tmpDir, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetCommitHistory(tmpDir, 10)
	}
}

// BenchmarkRepository_MultipleFiles_NoCache simulates building context from 50 files without cache
func BenchmarkRepository_MultipleFiles_NoCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepoWithMultipleFiles(b, 50)
	defer os.RemoveAll(tmpDir)

	files, _ := exec.Command("git", "-C", tmpDir, "ls-files").Output()
	fileList := []string{}
	for _, line := range splitLines(string(files)) {
		if line != "" {
			fileList = append(fileList, line)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repo := New(&domain.NoopLogger{}).(*Repository)
		for _, file := range fileList {
			_, _ = repo.GetFileAtRef(tmpDir, file, "HEAD")
		}
	}
}

// BenchmarkRepository_MultipleFiles_WithCache simulates building context from 50 files with cache
func BenchmarkRepository_MultipleFiles_WithCache(b *testing.B) {
	if !isGitAvailable() {
		b.Skip("Git not available")
	}

	tmpDir := createBenchmarkGitRepoWithMultipleFiles(b, 50)
	defer os.RemoveAll(tmpDir)

	files, _ := exec.Command("git", "-C", tmpDir, "ls-files").Output()
	fileList := []string{}
	for _, line := range splitLines(string(files)) {
		if line != "" {
			fileList = append(fileList, line)
		}
	}

	repo := New(&domain.NoopLogger{}).(*Repository)

	// Warm up cache
	for _, file := range fileList {
		_, _ = repo.GetFileAtRef(tmpDir, file, "HEAD")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, file := range fileList {
			_, _ = repo.GetFileAtRef(tmpDir, file, "HEAD")
		}
	}
}

// Helper functions for benchmarks

func createBenchmarkGitRepo(b *testing.B) string {
	tmpDir, err := os.MkdirTemp("", "git-bench-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "bench@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Bench User").Run()

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("benchmark test content"), 0644); err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	return tmpDir
}

func createBenchmarkGitRepoWithMultipleFiles(b *testing.B, fileCount int) string {
	tmpDir, err := os.MkdirTemp("", "git-bench-multi-*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		b.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git
	exec.Command("git", "-C", tmpDir, "config", "user.email", "bench@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Bench User").Run()

	// Create multiple test files
	for i := 0; i < fileCount; i++ {
		testFile := filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i))
		content := fmt.Sprintf("This is test file number %d with some content", i)
		if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
			os.RemoveAll(tmpDir)
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Add and commit
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit with multiple files").Run()

	return tmpDir
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
