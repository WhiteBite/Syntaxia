package projectstructure

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestCachedDetector_DetectStructure_Caches(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	// First call - should detect
	result1, err := cd.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("DetectStructure failed: %v", err)
	}
	if result1 == nil {
		t.Fatal("Expected non-nil result")
	}

	// Second call - should return cached
	result2, err := cd.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("Second DetectStructure failed: %v", err)
	}

	// Should be same pointer (cached)
	if result1 != result2 {
		t.Error("Expected cached result to be same pointer")
	}

	stats := cd.Stats()
	if stats["structure_entries"] != 1 {
		t.Errorf("Expected 1 structure entry, got %d", stats["structure_entries"])
	}
}

func TestCachedDetector_TTL_Expiration(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	// Very short TTL for testing
	cd := NewCachedDetectorWithTTL(50 * time.Millisecond)

	// First call
	result1, _ := cd.DetectStructure(tmpDir)

	// Wait for TTL to expire
	time.Sleep(100 * time.Millisecond)

	// Second call - should re-detect (new pointer)
	result2, _ := cd.DetectStructure(tmpDir)

	if result1 == result2 {
		t.Error("Expected new result after TTL expiration")
	}
}

func TestCachedDetector_Invalidate(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	// Populate cache
	_, _ = cd.DetectStructure(tmpDir)
	_, _ = cd.DetectArchitecture(tmpDir)
	_, _ = cd.DetectFrameworks(tmpDir)
	_, _ = cd.DetectConventions(tmpDir)

	stats := cd.Stats()
	if stats["structure_entries"] != 1 {
		t.Error("Expected cache to be populated")
	}

	// Invalidate
	cd.Invalidate(tmpDir)

	stats = cd.Stats()
	if stats["structure_entries"] != 0 {
		t.Error("Expected cache to be cleared after Invalidate")
	}
}

func TestCachedDetector_InvalidateAll(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()
	createTestProject(t, tmpDir1)
	createTestProject(t, tmpDir2)

	cd := NewCachedDetector()

	// Populate cache for both projects
	_, _ = cd.DetectStructure(tmpDir1)
	_, _ = cd.DetectStructure(tmpDir2)

	stats := cd.Stats()
	if stats["structure_entries"] != 2 {
		t.Errorf("Expected 2 entries, got %d", stats["structure_entries"])
	}

	// Invalidate all
	cd.InvalidateAll()

	stats = cd.Stats()
	if stats["structure_entries"] != 0 {
		t.Error("Expected all caches to be cleared")
	}
}

func TestCachedDetector_Concurrent(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	var wg sync.WaitGroup
	errors := make(chan error, 20)

	// Run concurrent detections
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cd.DetectStructure(tmpDir)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent detection failed: %v", err)
	}

	// Should have exactly 1 cached entry
	stats := cd.Stats()
	if stats["structure_entries"] != 1 {
		t.Errorf("Expected 1 entry after concurrent access, got %d", stats["structure_entries"])
	}
}

func TestCachedDetector_SetTTL(t *testing.T) {
	cd := NewCachedDetector()

	if cd.GetTTL() != DefaultCacheTTL {
		t.Errorf("Expected default TTL %v, got %v", DefaultCacheTTL, cd.GetTTL())
	}

	newTTL := 10 * time.Minute
	cd.SetTTL(newTTL)

	if cd.GetTTL() != newTTL {
		t.Errorf("Expected TTL %v, got %v", newTTL, cd.GetTTL())
	}
}

func TestCachedDetector_DetectArchitecture_Caches(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	result1, _ := cd.DetectArchitecture(tmpDir)
	result2, _ := cd.DetectArchitecture(tmpDir)

	if result1 != result2 {
		t.Error("Expected cached architecture result")
	}
}

func TestCachedDetector_DetectFrameworks_Caches(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	result1, _ := cd.DetectFrameworks(tmpDir)
	result2, _ := cd.DetectFrameworks(tmpDir)

	// For slices, check length (can't compare pointers directly)
	if len(result1) != len(result2) {
		t.Error("Expected same frameworks result")
	}
}

func TestCachedDetector_DetectConventions_Caches(t *testing.T) {
	tmpDir := t.TempDir()
	createTestProject(t, tmpDir)

	cd := NewCachedDetector()

	result1, _ := cd.DetectConventions(tmpDir)
	result2, _ := cd.DetectConventions(tmpDir)

	if result1 != result2 {
		t.Error("Expected cached conventions result")
	}
}

// createTestProject creates a minimal test project structure
func createTestProject(t *testing.T, dir string) {
	t.Helper()

	// Create go.mod
	goMod := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create main.go
	mainGo := filepath.Join(dir, "main.go")
	if err := os.WriteFile(mainGo, []byte("package main\nfunc main() {}\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create domain directory
	domainDir := filepath.Join(dir, "domain")
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		t.Fatal(err)
	}
	modelGo := filepath.Join(domainDir, "model.go")
	if err := os.WriteFile(modelGo, []byte("package domain\ntype User struct{}\n"), 0644); err != nil {
		t.Fatal(err)
	}
}
