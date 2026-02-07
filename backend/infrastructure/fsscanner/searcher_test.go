package fsscanner

import (
	"syntaxia/domain"
	"testing"
)

// mockTreeBuilder для тестирования
type mockTreeBuilder struct {
	tree []*domain.FileNode
	err  error
}

func (m *mockTreeBuilder) BuildTree(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*domain.FileNode, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tree, nil
}

func (m *mockTreeBuilder) InvalidateCache() {}

// mockLogger для тестирования
type mockLogger struct{}

func (l *mockLogger) Debug(message string)   {}
func (l *mockLogger) Info(message string)    {}
func (l *mockLogger) Warning(message string) {}
func (l *mockLogger) Error(message string)   {}
func (l *mockLogger) Fatal(message string)   {}

func createTestTree() []*domain.FileNode {
	return []*domain.FileNode{
		{
			Name:  "root",
			Path:  "/test/project",
			IsDir: true,
			Children: []*domain.FileNode{
				{
					Name:        "main.go",
					Path:        "/test/project/main.go",
					IsDir:       false,
					Size:        1024,
					ContentType: "text",
				},
				{
					Name:        "helper.go",
					Path:        "/test/project/helper.go",
					IsDir:       false,
					Size:        512,
					ContentType: "text",
				},
				{
					Name:  "src",
					Path:  "/test/project/src",
					IsDir: true,
					Children: []*domain.FileNode{
						{
							Name:        "component.ts",
							Path:        "/test/project/src/component.ts",
							IsDir:       false,
							Size:        2048,
							ContentType: "text",
						},
						{
							Name:        "utils.ts",
							Path:        "/test/project/src/utils.ts",
							IsDir:       false,
							Size:        768,
							ContentType: "text",
						},
					},
				},
			},
		},
	}
}

func TestBuildSearchIndex(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	err := searcher.BuildSearchIndex("/test/project")
	if err != nil {
		t.Fatalf("BuildSearchIndex failed: %v", err)
	}

	stats := searcher.GetSearchStats()
	if stats["cached_indices"].(int) != 1 {
		t.Errorf("Expected 1 cached index, got %d", stats["cached_indices"])
	}
}

func TestSearchFiles_ExactMatch(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    10,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   false,
	}

	results, err := searcher.SearchFiles("/test/project", "main", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "main.go" {
		t.Errorf("Expected main.go, got %s", results[0].Name)
	}
}

func TestSearchFiles_CaseInsensitive(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    10,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   false,
	}

	results, err := searcher.SearchFiles("/test/project", "MAIN", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result for case-insensitive search, got %d", len(results))
	}
}

func TestSearchFiles_WithPath(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    10,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   true,
	}

	results, err := searcher.SearchFiles("/test/project", "src", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	// Should find files in src directory
	if len(results) < 2 {
		t.Errorf("Expected at least 2 results when searching path, got %d", len(results))
	}
}

func TestSearchFiles_FileTypeFilter(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    10,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   true,
		FileTypes:     []string{".ts"},
	}

	results, err := searcher.SearchFiles("/test/project", "component", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 .ts file, got %d", len(results))
	}

	if len(results) > 0 && results[0].Name != "component.ts" {
		t.Errorf("Expected component.ts, got %s", results[0].Name)
	}
}

func TestSearchFiles_MaxResults(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    2,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   true,
	}

	// Search for something that matches multiple files
	results, err := searcher.SearchFiles("/test/project", ".", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	if len(results) > 2 {
		t.Errorf("Expected max 2 results, got %d", len(results))
	}
}

func TestInvalidateSearchIndex(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	// Build index
	err := searcher.BuildSearchIndex("/test/project")
	if err != nil {
		t.Fatalf("BuildSearchIndex failed: %v", err)
	}

	stats := searcher.GetSearchStats()
	if stats["cached_indices"].(int) != 1 {
		t.Errorf("Expected 1 cached index before invalidation")
	}

	// Invalidate
	searcher.InvalidateSearchIndex("/test/project")

	stats = searcher.GetSearchStats()
	if stats["cached_indices"].(int) != 0 {
		t.Errorf("Expected 0 cached indices after invalidation, got %d", stats["cached_indices"])
	}
}

func TestSearchFiles_Scoring(t *testing.T) {
	tree := createTestTree()
	mockBuilder := &mockTreeBuilder{tree: tree}
	logger := &mockLogger{}

	searcher := NewFileSearcher(mockBuilder, logger)

	options := domain.SearchOptions{
		MaxResults:    10,
		FuzzyMatch:    false,
		CaseSensitive: false,
		IncludePath:   false,
	}

	results, err := searcher.SearchFiles("/test/project", "main", options)
	if err != nil {
		t.Fatalf("SearchFiles failed: %v", err)
	}

	if len(results) > 0 {
		// Exact match at start should have high score
		if results[0].Score < 0.8 {
			t.Errorf("Expected high score for exact match, got %f", results[0].Score)
		}
	}
}
