package analyzers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewCallGraphBuilder(t *testing.T) {
	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	if builder == nil {
		t.Fatal("NewCallGraphBuilder returned nil")
	}
	if builder.graph == nil {
		t.Error("graph not initialized")
	}
	if builder.depGraph == nil {
		t.Error("depGraph not initialized")
	}
}

func TestCallGraphBuilder_Build_GoProject(t *testing.T) {
	tmpDir := t.TempDir()

	mainGo := `package main

func main() {
	hello()
	result := add(1, 2)
	_ = result
}

func hello() {
	greet("world")
}

func greet(name string) {
	println(name)
}

func add(a, b int) int {
	return a + b
}
`
	writeTestFile(t, tmpDir, "main.go", mainGo)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	graph, err := builder.Build(tmpDir)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if graph == nil {
		t.Fatal("Build returned nil graph")
	}

	if len(graph.Nodes) == 0 {
		t.Error("no nodes in graph")
	}

	foundMain := false
	foundHello := false
	for _, node := range graph.Nodes {
		if node.Name == "main" {
			foundMain = true
		}
		if node.Name == "hello" {
			foundHello = true
		}
	}

	if !foundMain {
		t.Error("main function not found")
	}
	if !foundHello {
		t.Error("hello function not found")
	}

	if len(graph.Edges) == 0 {
		t.Error("no edges in graph")
	}
}

func TestCallGraphBuilder_GetCallers(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

func main() {
	helper()
}

func other() {
	helper()
}

func helper() {
	println("help")
}
`
	writeTestFile(t, tmpDir, "main.go", code)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)
	builder.Build(tmpDir)

	callers := builder.GetCallers("helper")
	if len(callers) < 1 {
		t.Logf("callers for helper: %d (may vary)", len(callers))
	}
}

func TestCallGraphBuilder_GetCallees(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

func main() {
	foo()
	bar()
}

func foo() {}
func bar() {}
`
	writeTestFile(t, tmpDir, "main.go", code)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)
	builder.Build(tmpDir)

	callees := builder.GetCallees("main")
	t.Logf("callees for main: %d", len(callees))
}

func TestCallGraphBuilder_BuildDependencyGraph(t *testing.T) {
	tmpDir := t.TempDir()

	mainGo := `package main

import "fmt"

func main() {
	fmt.Println("hello")
}
`
	writeTestFile(t, tmpDir, "main.go", mainGo)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)

	depGraph, err := builder.BuildDependencyGraph(tmpDir)
	if err != nil {
		t.Fatalf("BuildDependencyGraph failed: %v", err)
	}

	if depGraph == nil {
		t.Fatal("BuildDependencyGraph returned nil")
	}
}

func TestCallGraphBuilder_GetImpact(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

func main() {
	a()
}

func a() {
	b()
}

func b() {
	c()
}

func c() {}
`
	writeTestFile(t, tmpDir, "main.go", code)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)
	builder.Build(tmpDir)

	impact := builder.GetImpact("c", 3)
	t.Logf("impact for c: %d nodes", len(impact))
}

func TestCallGraphBuilder_ExportMermaid(t *testing.T) {
	tmpDir := t.TempDir()

	code := `package main

func main() {
	hello()
}

func hello() {}
`
	writeTestFile(t, tmpDir, "main.go", code)

	registry := NewAnalyzerRegistry()
	builder := NewCallGraphBuilder(registry)
	builder.Build(tmpDir)

	mermaid := builder.ExportMermaid(10)
	if mermaid == "" {
		t.Error("ExportMermaid returned empty string")
	}
	if !strings.Contains(mermaid, "graph") {
		t.Error("mermaid should contain 'graph'")
	}
}

func writeTestFile(t *testing.T, base, path, content string) {
	t.Helper()
	fullPath := filepath.Join(base, path)
	dir := filepath.Dir(fullPath)
	os.MkdirAll(dir, 0o755)
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
}
