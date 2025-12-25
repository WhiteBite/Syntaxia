package analyzers

import (
	"os"
	"path/filepath"
	"syntaxia/domain/analysis"
	"sync"
	"testing"
)

func TestCallGraphBuilder_EnsureBuilt_OnlyOnce(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func Hello() { World() }
func World() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// First call should build
	graph1, err := builder.EnsureBuilt(tmpDir)
	if err != nil {
		t.Fatalf("EnsureBuilt failed: %v", err)
	}
	if graph1 == nil {
		t.Fatal("Expected non-nil graph")
	}

	nodeCount1 := len(graph1.Nodes)
	if nodeCount1 == 0 {
		t.Error("Expected nodes after building")
	}

	// Second call should return cached result
	graph2, err := builder.EnsureBuilt(tmpDir)
	if err != nil {
		t.Fatalf("Second EnsureBuilt failed: %v", err)
	}

	if graph1 != graph2 {
		t.Error("Expected same graph instance on second call")
	}
}

func TestCallGraphBuilder_EnsureBuilt_Concurrent(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func Concurrent() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, 10)
	graphs := make(chan *analysis.CallGraph, 10)

	// Run 10 concurrent EnsureBuilt calls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g, err := builder.EnsureBuilt(tmpDir)
			if err != nil {
				errors <- err
				return
			}
			graphs <- g
		}()
	}

	wg.Wait()
	close(errors)
	close(graphs)

	for err := range errors {
		t.Errorf("Concurrent EnsureBuilt failed: %v", err)
	}

	// All should return the same graph
	var firstGraph *analysis.CallGraph
	for g := range graphs {
		if firstGraph == nil {
			firstGraph = g
		} else if g != firstGraph {
			t.Error("Expected same graph instance from all concurrent calls")
		}
	}

	if !builder.IsBuilt() {
		t.Error("Expected IsBuilt() to return true")
	}
}

func TestCallGraphBuilder_Invalidate(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	content := `package main
func ToInvalidate() {}
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Build
	_, err := builder.EnsureBuilt(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if !builder.IsBuilt() {
		t.Error("Expected IsBuilt() after EnsureBuilt")
	}

	// Invalidate
	builder.Invalidate()

	if builder.IsBuilt() {
		t.Error("Expected not built after Invalidate")
	}

	if builder.GetProjectRoot() != "" {
		t.Error("Expected empty project root after Invalidate")
	}

	// Re-build should work
	graph, err := builder.EnsureBuilt(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	if !builder.IsBuilt() {
		t.Error("Expected IsBuilt() after re-building")
	}

	if len(graph.Nodes) == 0 {
		t.Error("Expected nodes after re-building")
	}
}

func TestCallGraphBuilder_ProjectChange_Rebuilds(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	// Create two temp projects
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	file1 := filepath.Join(tmpDir1, "main.go")
	file2 := filepath.Join(tmpDir2, "main.go")

	content1 := `package main
func Project1() {}
`
	content2 := `package main
func Project2A() {}
func Project2B() {}
func Project2C() {}
`
	if err := os.WriteFile(file1, []byte(content1), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte(content2), 0644); err != nil {
		t.Fatal(err)
	}

	// Build project 1
	graph1, err := builder.EnsureBuilt(tmpDir1)
	if err != nil {
		t.Fatal(err)
	}
	nodes1 := len(graph1.Nodes)

	if builder.GetProjectRoot() != tmpDir1 {
		t.Errorf("Expected project root %s, got %s", tmpDir1, builder.GetProjectRoot())
	}

	// Build project 2 - should auto-invalidate and rebuild
	graph2, err := builder.EnsureBuilt(tmpDir2)
	if err != nil {
		t.Fatal(err)
	}
	nodes2 := len(graph2.Nodes)

	if builder.GetProjectRoot() != tmpDir2 {
		t.Errorf("Expected project root %s, got %s", tmpDir2, builder.GetProjectRoot())
	}

	// Project 2 has more functions
	if nodes2 <= nodes1 {
		t.Errorf("Expected more nodes in project2, got %d vs %d", nodes2, nodes1)
	}
}

func TestCallGraphBuilder_GetProjectRoot(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	// Initially empty
	if builder.GetProjectRoot() != "" {
		t.Error("Expected empty project root initially")
	}

	tmpDir := t.TempDir()
	goFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(goFile, []byte("package main\nfunc X(){}"), 0644); err != nil {
		t.Fatal(err)
	}

	_, _ = builder.EnsureBuilt(tmpDir)

	if builder.GetProjectRoot() != tmpDir {
		t.Errorf("Expected project root %s, got %s", tmpDir, builder.GetProjectRoot())
	}
}
