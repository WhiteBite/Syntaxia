package analyzers

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestSymbolIndex_EnsureIndexed_OnlyOnce(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	// Create temp project
	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func Hello() {}
func World() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// First call should index
	err := idx.EnsureIndexed(ctx, tmpDir)
	if err != nil {
		t.Fatalf("EnsureIndexed failed: %v", err)
	}

	stats1 := idx.Stats()
	if stats1["total_symbols"] == 0 {
		t.Error("Expected symbols after indexing")
	}

	// Second call should return cached result (not re-index)
	err = idx.EnsureIndexed(ctx, tmpDir)
	if err != nil {
		t.Fatalf("Second EnsureIndexed failed: %v", err)
	}

	stats2 := idx.Stats()
	if stats1["total_symbols"] != stats2["total_symbols"] {
		t.Errorf("Stats changed between calls: %d vs %d", stats1["total_symbols"], stats2["total_symbols"])
	}
}

func TestSymbolIndex_EnsureIndexed_Concurrent(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func Concurrent() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	// Run 10 concurrent EnsureIndexed calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := idx.EnsureIndexed(ctx, tmpDir); err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Concurrent EnsureIndexed failed: %v", err)
	}

	// Should have indexed exactly once
	if !idx.IsIndexed() {
		t.Error("Expected index to be built")
	}
}

func TestSymbolIndex_Invalidate(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func ToInvalidate() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Index
	if err := idx.EnsureIndexed(ctx, tmpDir); err != nil {
		t.Fatal(err)
	}

	if !idx.IsIndexed() {
		t.Error("Expected indexed after EnsureIndexed")
	}

	// Invalidate
	idx.Invalidate()

	if idx.IsIndexed() {
		t.Error("Expected not indexed after Invalidate")
	}

	stats := idx.Stats()
	if stats["total_symbols"] != 0 {
		t.Errorf("Expected 0 symbols after invalidate, got %d", stats["total_symbols"])
	}

	// Re-index should work
	if err := idx.EnsureIndexed(ctx, tmpDir); err != nil {
		t.Fatal(err)
	}

	if !idx.IsIndexed() {
		t.Error("Expected indexed after re-indexing")
	}
}

func TestSymbolIndex_InvalidateFile(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	tmpDir := t.TempDir()

	// Create two files
	file1 := filepath.Join(tmpDir, "file1.go")
	file2 := filepath.Join(tmpDir, "file2.go")

	content1 := `package main
func FromFile1() {}
`
	content2 := `package main
func FromFile2() {}
`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := idx.EnsureIndexed(ctx, tmpDir); err != nil {
		t.Fatal(err)
	}

	statsBefore := idx.Stats()
	symbolsBefore := statsBefore["total_symbols"]

	// Invalidate only file1
	idx.InvalidateFile("file1.go")

	statsAfter := idx.Stats()
	symbolsAfter := statsAfter["total_symbols"]

	if symbolsAfter >= symbolsBefore {
		t.Errorf("Expected fewer symbols after InvalidateFile, got %d vs %d", symbolsAfter, symbolsBefore)
	}

	// file2 symbols should still exist
	file2Symbols := idx.GetSymbolsInFile("file2.go")
	if len(file2Symbols) == 0 {
		t.Error("Expected file2 symbols to remain after InvalidateFile(file1)")
	}
}

func TestSymbolIndex_ProjectChange_ReIndexes(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	// Create two temp projects
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	file1 := filepath.Join(tmpDir1, "main.go")
	file2 := filepath.Join(tmpDir2, "main.go")

	content1 := `package main
func Project1Func() {}
`
	content2 := `package main
func Project2Func() {}
func AnotherFunc() {}
`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// Index project 1
	if err := idx.EnsureIndexed(ctx, tmpDir1); err != nil {
		t.Fatal(err)
	}

	stats1 := idx.Stats()

	// Index project 2 - should auto-invalidate and re-index
	if err := idx.EnsureIndexed(ctx, tmpDir2); err != nil {
		t.Fatal(err)
	}

	stats2 := idx.Stats()

	// Project 2 has more symbols
	if stats2["total_symbols"] <= stats1["total_symbols"] {
		t.Errorf("Expected more symbols in project2, got %d vs %d", stats2["total_symbols"], stats1["total_symbols"])
	}
}
