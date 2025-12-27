package rag

import (
	"context"
	"math"
	"testing"
)

func TestFileRanker_RankFiles_EmptyProject(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Empty project should return empty slice
	ranked, err := ranker.RankFiles(context.Background(), "/nonexistent")
	if err != nil {
		t.Fatalf("RankFiles() error = %v", err)
	}

	if len(ranked) != 0 {
		t.Errorf("RankFiles() for empty project = %d files, want 0", len(ranked))
	}
}

func TestFileRanker_Cache(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Manually set cache to test cache operations
	testFiles := []RankedFile{{FilePath: "test.go", Rank: 1.0}}
	ranker.setCache("/test", testFiles, nil)

	// Check cache exists
	cached := ranker.getCached("/test")
	if cached == nil {
		t.Error("Expected cache to be populated after setCache()")
	}

	// Invalidate
	ranker.InvalidateCache("/test")

	// Check cache cleared
	cached = ranker.getCached("/test")
	if cached != nil {
		t.Error("Expected cache to be cleared after InvalidateCache()")
	}
}

func TestPageRank_SimpleGraph(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Create a simple graph: A -> B -> C
	graph := &importGraph{
		nodes: map[string]int{
			"a.go": 0,
			"b.go": 1,
			"c.go": 2,
		},
		files: []string{"a.go", "b.go", "c.go"},
		adjacency: [][]int{
			{1},    // a imports b
			{2},    // b imports c
			{},     // c imports nothing
		},
		reverseAdj: [][]int{
			{},     // nothing imports a
			{0},    // a imports b
			{1},    // b imports c
		},
		exportCounts: []int{1, 2, 3},
		fileSizes:    []int64{100, 200, 300},
	}

	ranks := ranker.pageRank(graph)

	if len(ranks) != 3 {
		t.Fatalf("pageRank() returned %d ranks, want 3", len(ranks))
	}

	// C should have highest rank (most imported transitively)
	// Ranks are normalized, so we check relative ordering
	if ranks[2] < ranks[0] {
		t.Error("Expected c.go to have higher rank than a.go")
	}
}

func TestPageRank_CyclicGraph(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Create a cyclic graph: A -> B -> C -> A
	graph := &importGraph{
		nodes: map[string]int{
			"a.go": 0,
			"b.go": 1,
			"c.go": 2,
		},
		files: []string{"a.go", "b.go", "c.go"},
		adjacency: [][]int{
			{1},    // a imports b
			{2},    // b imports c
			{0},    // c imports a (cycle)
		},
		reverseAdj: [][]int{
			{2},    // c imports a
			{0},    // a imports b
			{1},    // b imports c
		},
		exportCounts: []int{1, 1, 1},
		fileSizes:    []int64{100, 100, 100},
	}

	ranks := ranker.pageRank(graph)

	if len(ranks) != 3 {
		t.Fatalf("pageRank() returned %d ranks, want 3", len(ranks))
	}

	// All nodes should have similar ranks in a cycle
	avgRank := (ranks[0] + ranks[1] + ranks[2]) / 3
	for i, rank := range ranks {
		diff := math.Abs(rank - avgRank)
		if diff > 0.3 { // Allow some variance due to normalization
			t.Errorf("Node %d rank %.4f differs too much from average %.4f", i, rank, avgRank)
		}
	}
}

func TestPageRank_DanglingNodes(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Graph with dangling node (no outgoing edges)
	graph := &importGraph{
		nodes: map[string]int{
			"a.go": 0,
			"b.go": 1,
		},
		files: []string{"a.go", "b.go"},
		adjacency: [][]int{
			{1},    // a imports b
			{},     // b imports nothing (dangling)
		},
		reverseAdj: [][]int{
			{},     // nothing imports a
			{0},    // a imports b
		},
		exportCounts: []int{1, 1},
		fileSizes:    []int64{100, 100},
	}

	ranks := ranker.pageRank(graph)

	if len(ranks) != 2 {
		t.Fatalf("pageRank() returned %d ranks, want 2", len(ranks))
	}

	// Both should have valid ranks
	for i, rank := range ranks {
		if rank <= 0 {
			t.Errorf("Node %d has invalid rank %.4f", i, rank)
		}
	}
}

func TestPageRank_EntryPointBoost(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Graph with entry point
	graph := &importGraph{
		nodes: map[string]int{
			"main.go":  0, // Entry point
			"utils.go": 1,
		},
		files: []string{"main.go", "utils.go"},
		adjacency: [][]int{
			{1},    // main imports utils
			{},     // utils imports nothing
		},
		reverseAdj: [][]int{
			{},     // nothing imports main
			{0},    // main imports utils
		},
		exportCounts: []int{1, 5},
		fileSizes:    []int64{100, 100},
	}

	ranks := ranker.pageRank(graph)

	// main.go should get entry point boost
	// Even though utils.go is imported, main.go should have competitive rank
	if ranks[0] < 0.1 {
		t.Errorf("Entry point main.go has too low rank: %.4f", ranks[0])
	}
}

func TestPageRank_TestFileDiscount(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Graph with test file
	graph := &importGraph{
		nodes: map[string]int{
			"service.go":      0,
			"service_test.go": 1, // Test file
		},
		files: []string{"service.go", "service_test.go"},
		adjacency: [][]int{
			{},     // service imports nothing
			{0},    // test imports service
		},
		reverseAdj: [][]int{
			{1},    // test imports service
			{},     // nothing imports test
		},
		exportCounts: []int{5, 1},
		fileSizes:    []int64{500, 300},
	}

	ranks := ranker.pageRank(graph)

	// Test file should have lower rank due to discount
	if ranks[1] >= ranks[0] {
		t.Errorf("Test file rank %.4f should be lower than service rank %.4f", ranks[1], ranks[0])
	}
}

func TestResolveImport(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	nodes := map[string]int{
		"src/utils.ts":       0,
		"src/utils/index.ts": 1,
		"src/components/Button.tsx": 2,
	}

	tests := []struct {
		name       string
		fromFile   string
		importPath string
		want       string
	}{
		{
			name:       "relative import same dir",
			fromFile:   "src/main.ts",
			importPath: "./utils",
			want:       "src/utils.ts",
		},
		{
			name:       "relative import with extension",
			fromFile:   "src/main.ts",
			importPath: "./utils.ts",
			want:       "src/utils.ts",
		},
		{
			name:       "relative import index file",
			fromFile:   "src/main.ts",
			importPath: "./utils",
			want:       "src/utils.ts", // Prefers direct file over index
		},
		{
			name:       "empty import",
			fromFile:   "src/main.ts",
			importPath: "",
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ranker.resolveImport(tt.fromFile, tt.importPath, nodes)
			if got != tt.want {
				t.Errorf("resolveImport() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRankedFile_Sorting(t *testing.T) {
	files := []RankedFile{
		{FilePath: "low.go", Rank: 0.1},
		{FilePath: "high.go", Rank: 0.9},
		{FilePath: "mid.go", Rank: 0.5},
	}

	// Sort by rank descending (as done in RankFiles)
	for i := 0; i < len(files)-1; i++ {
		for j := i + 1; j < len(files); j++ {
			if files[j].Rank > files[i].Rank {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	if files[0].FilePath != "high.go" {
		t.Errorf("First file should be high.go, got %s", files[0].FilePath)
	}
	if files[1].FilePath != "mid.go" {
		t.Errorf("Second file should be mid.go, got %s", files[1].FilePath)
	}
	if files[2].FilePath != "low.go" {
		t.Errorf("Third file should be low.go, got %s", files[2].FilePath)
	}
}

func TestImportGraph_Structure(t *testing.T) {
	graph := &importGraph{
		nodes:        make(map[string]int),
		files:        make([]string, 0),
		adjacency:    make([][]int, 0),
		reverseAdj:   make([][]int, 0),
		exportCounts: make([]int, 0),
		fileSizes:    make([]int64, 0),
	}

	// Add nodes
	files := []string{"a.go", "b.go", "c.go"}
	for i, f := range files {
		graph.nodes[f] = i
		graph.files = append(graph.files, f)
		graph.adjacency = append(graph.adjacency, []int{})
		graph.reverseAdj = append(graph.reverseAdj, []int{})
		graph.exportCounts = append(graph.exportCounts, 0)
		graph.fileSizes = append(graph.fileSizes, 100)
	}

	// Add edge: a -> b
	graph.adjacency[0] = append(graph.adjacency[0], 1)
	graph.reverseAdj[1] = append(graph.reverseAdj[1], 0)

	// Verify structure
	if len(graph.files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(graph.files))
	}

	if len(graph.adjacency[0]) != 1 {
		t.Errorf("Expected 1 outgoing edge from a.go, got %d", len(graph.adjacency[0]))
	}

	if len(graph.reverseAdj[1]) != 1 {
		t.Errorf("Expected 1 incoming edge to b.go, got %d", len(graph.reverseAdj[1]))
	}
}

func TestPageRank_Convergence(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Create a larger graph to test convergence
	n := 10
	graph := &importGraph{
		nodes:        make(map[string]int),
		files:        make([]string, n),
		adjacency:    make([][]int, n),
		reverseAdj:   make([][]int, n),
		exportCounts: make([]int, n),
		fileSizes:    make([]int64, n),
	}

	for i := 0; i < n; i++ {
		name := string(rune('a'+i)) + ".go"
		graph.nodes[name] = i
		graph.files[i] = name
		graph.adjacency[i] = []int{}
		graph.reverseAdj[i] = []int{}
		graph.exportCounts[i] = 1
		graph.fileSizes[i] = 100
	}

	// Create chain: 0 -> 1 -> 2 -> ... -> n-1
	for i := 0; i < n-1; i++ {
		graph.adjacency[i] = append(graph.adjacency[i], i+1)
		graph.reverseAdj[i+1] = append(graph.reverseAdj[i+1], i)
	}

	ranks := ranker.pageRank(graph)

	// All ranks should be positive
	for i, rank := range ranks {
		if rank <= 0 {
			t.Errorf("Node %d has non-positive rank: %.6f", i, rank)
		}
	}

	// Sum should be close to n (normalized)
	sum := 0.0
	for _, rank := range ranks {
		sum += rank
	}

	// After normalization, max rank is 1.0
	maxRank := 0.0
	for _, rank := range ranks {
		if rank > maxRank {
			maxRank = rank
		}
	}

	if math.Abs(maxRank-1.0) > 0.01 {
		t.Errorf("Max rank should be ~1.0 after normalization, got %.4f", maxRank)
	}
}

func TestGetTopFiles(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Pre-populate cache with test data
	testFiles := []RankedFile{
		{FilePath: "a.go", Rank: 1.0},
		{FilePath: "b.go", Rank: 0.8},
		{FilePath: "c.go", Rank: 0.6},
		{FilePath: "d.go", Rank: 0.4},
		{FilePath: "e.go", Rank: 0.2},
	}

	ranker.setCache("/test", testFiles, nil)

	tests := []struct {
		name    string
		n       int
		wantLen int
	}{
		{
			name:    "top 3",
			n:       3,
			wantLen: 3,
		},
		{
			name:    "top 10 (more than available)",
			n:       10,
			wantLen: 5,
		},
		{
			name:    "top 0",
			n:       0,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ranker.GetTopFiles(context.Background(), "/test", tt.n)
			if err != nil {
				t.Fatalf("GetTopFiles() error = %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("GetTopFiles() len = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestGetFileRank(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Pre-populate cache
	testFiles := []RankedFile{
		{FilePath: "main.go", Rank: 1.0},
		{FilePath: "utils.go", Rank: 0.5},
	}

	ranker.setCache("/test", testFiles, nil)

	tests := []struct {
		name     string
		filePath string
		wantRank float64
	}{
		{
			name:     "existing file",
			filePath: "main.go",
			wantRank: 1.0,
		},
		{
			name:     "another existing file",
			filePath: "utils.go",
			wantRank: 0.5,
		},
		{
			name:     "non-existing file",
			filePath: "notfound.go",
			wantRank: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ranker.GetFileRank(context.Background(), "/test", tt.filePath)
			if err != nil {
				t.Fatalf("GetFileRank() error = %v", err)
			}

			if got != tt.wantRank {
				t.Errorf("GetFileRank() = %.4f, want %.4f", got, tt.wantRank)
			}
		})
	}
}


func TestFileRanker_GetImportGraph(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Set up cache with graph
	graph := &importGraph{
		nodes: map[string]int{"a.go": 0, "b.go": 1},
		files: []string{"a.go", "b.go"},
		adjacency: [][]int{
			{1}, // a imports b
			{},
		},
		reverseAdj:   [][]int{{}, {0}},
		exportCounts: []int{1, 2},
		fileSizes:    []int64{100, 200},
	}

	ranker.setCache("/test", []RankedFile{}, graph)

	importGraph, err := ranker.GetImportGraph(context.Background(), "/test")
	if err != nil {
		t.Fatalf("GetImportGraph() error = %v", err)
	}

	if len(importGraph) != 2 {
		t.Errorf("GetImportGraph() returned %d files, want 2", len(importGraph))
	}

	if len(importGraph["a.go"]) != 1 {
		t.Errorf("a.go should have 1 import, got %d", len(importGraph["a.go"]))
	}

	if importGraph["a.go"][0] != "b.go" {
		t.Errorf("a.go should import b.go, got %s", importGraph["a.go"][0])
	}
}

func TestFileRanker_GetImportGraph_NoCache(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// No cache set, should return error
	_, err := ranker.GetImportGraph(context.Background(), "/nonexistent")
	if err == nil {
		t.Error("GetImportGraph() should return error when no cache")
	}
}

func TestPageRank_SingleNode(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	graph := &importGraph{
		nodes:        map[string]int{"single.go": 0},
		files:        []string{"single.go"},
		adjacency:    [][]int{{}},
		reverseAdj:   [][]int{{}},
		exportCounts: []int{1},
		fileSizes:    []int64{100},
	}

	ranks := ranker.pageRank(graph)

	if len(ranks) != 1 {
		t.Fatalf("pageRank() returned %d ranks, want 1", len(ranks))
	}

	if ranks[0] != 1.0 {
		t.Errorf("Single node should have rank 1.0, got %.4f", ranks[0])
	}
}

func TestPageRank_EmptyGraph(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	graph := &importGraph{
		nodes:        map[string]int{},
		files:        []string{},
		adjacency:    [][]int{},
		reverseAdj:   [][]int{},
		exportCounts: []int{},
		fileSizes:    []int64{},
	}

	ranks := ranker.pageRank(graph)

	if len(ranks) != 0 {
		t.Errorf("pageRank() on empty graph should return empty slice, got %d", len(ranks))
	}
}

func TestPageRank_HighExportBoost(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	// Two isolated nodes, one with many exports
	graph := &importGraph{
		nodes: map[string]int{
			"few_exports.go":  0,
			"many_exports.go": 1,
		},
		files:      []string{"few_exports.go", "many_exports.go"},
		adjacency:  [][]int{{}, {}},
		reverseAdj: [][]int{{}, {}},
		exportCounts: []int{1, 20}, // many_exports has 20 exports
		fileSizes:    []int64{100, 100},
	}

	ranks := ranker.pageRank(graph)

	// File with more exports should have higher rank
	if ranks[1] <= ranks[0] {
		t.Errorf("File with more exports should have higher rank: few=%.4f, many=%.4f", ranks[0], ranks[1])
	}
}

func TestResolveImport_AliasImport(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	nodes := map[string]int{
		"src/components/Button.vue": 0,
	}

	// Test @/ alias resolution
	result := ranker.resolveImport("src/main.ts", "@/components/Button", nodes)

	// Should try to resolve but may not find exact match
	// The function returns the resolved path attempt
	if result == "" {
		t.Error("resolveImport() should return non-empty for alias import")
	}
}

func TestResolveImport_ParentDirectory(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	logger := &mockLogger{}

	ranker := NewFileRanker(registry, logger)

	nodes := map[string]int{
		"src/utils.ts": 0,
	}

	result := ranker.resolveImport("src/components/Button.tsx", "../utils", nodes)

	if result != "src/utils.ts" {
		t.Errorf("resolveImport() = %q, want 'src/utils.ts'", result)
	}
}
