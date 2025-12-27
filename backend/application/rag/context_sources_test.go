package rag

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"syntaxia/domain"
)

// ctxSourceMockLogger implements domain.Logger for testing
type ctxSourceMockLogger struct{}

func (m *ctxSourceMockLogger) Debug(message string)   {}
func (m *ctxSourceMockLogger) Info(message string)    {}
func (m *ctxSourceMockLogger) Warning(message string) {}
func (m *ctxSourceMockLogger) Error(message string)   {}
func (m *ctxSourceMockLogger) Fatal(message string)   {}

// ctxSourceMockFileReader implements domain.FileContentReader for testing
type ctxSourceMockFileReader struct {
	contents map[string]string
}

func (m *ctxSourceMockFileReader) ReadContents(
	ctx context.Context,
	filePaths []string,
	rootDir string,
	progress func(current, total int64),
) (map[string]string, error) {
	result := make(map[string]string)
	for _, path := range filePaths {
		if content, ok := m.contents[path]; ok {
			result[path] = content
		}
	}
	return result, nil
}

// ctxSourceMockGitRepo implements domain.GitRepository for testing
type ctxSourceMockGitRepo struct {
	isGitRepo        bool
	allFiles         []string
	uncommittedFiles []domain.FileStatus
	diff             string
	commitHistory    []domain.CommitInfo
}

func (m *ctxSourceMockGitRepo) IsGitRepository(projectPath string) bool {
	return m.isGitRepo
}

func (m *ctxSourceMockGitRepo) GetAllFiles(projectPath string) ([]string, error) {
	return m.allFiles, nil
}

func (m *ctxSourceMockGitRepo) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	return m.uncommittedFiles, nil
}

func (m *ctxSourceMockGitRepo) GenerateDiff(projectPath string) (string, error) {
	return m.diff, nil
}

func (m *ctxSourceMockGitRepo) GetCommitHistory(projectPath string, limit int) ([]domain.CommitInfo, error) {
	if limit > len(m.commitHistory) {
		return m.commitHistory, nil
	}
	return m.commitHistory[:limit], nil
}

// Stub implementations for other GitRepository methods
func (m *ctxSourceMockGitRepo) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	return nil, nil
}
func (m *ctxSourceMockGitRepo) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	return "", nil
}
func (m *ctxSourceMockGitRepo) GetGitignoreContent(projectRoot string) (string, error) { return "", nil }
func (m *ctxSourceMockGitRepo) IsGitAvailable() bool                                   { return true }
func (m *ctxSourceMockGitRepo) GetBranches(projectRoot string) ([]string, error)       { return nil, nil }
func (m *ctxSourceMockGitRepo) GetCurrentBranch(projectRoot string) (string, error)    { return "main", nil }
func (m *ctxSourceMockGitRepo) CloneRepository(url, targetPath string, depth int) error {
	return nil
}
func (m *ctxSourceMockGitRepo) CheckoutBranch(projectPath, branch string) error     { return nil }
func (m *ctxSourceMockGitRepo) CheckoutCommit(projectPath, commitHash string) error { return nil }
func (m *ctxSourceMockGitRepo) ListFilesAtRef(projectPath, ref string) ([]string, error) {
	return nil, nil
}
func (m *ctxSourceMockGitRepo) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
	return "", nil
}
func (m *ctxSourceMockGitRepo) FetchRemoteBranches(projectPath string) ([]string, error) { return nil, nil }

func TestGetAvailableSources(t *testing.T) {
	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)

	sources := manager.GetAvailableSources()

	expectedSources := []SourceType{
		SourceCode,
		SourceGitDiff,
		SourceGitHistory,
		SourceTests,
		SourceDocs,
		SourceTerminal,
	}

	if len(sources) != len(expectedSources) {
		t.Errorf("expected %d sources, got %d", len(expectedSources), len(sources))
	}

	for i, expected := range expectedSources {
		if sources[i] != expected {
			t.Errorf("expected source %s at index %d, got %s", expected, i, sources[i])
		}
	}
}

func TestCollectFromSources_EmptySources(t *testing.T) {
	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)

	result, err := manager.CollectFromSources(context.Background(), SourceRequest{
		ProjectRoot: "/tmp/test",
		Sources:     []SourceType{},
		MaxTokens:   10000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 0 {
		t.Errorf("expected 0 items, got %d", len(result.Items))
	}

	if result.TotalTokens != 0 {
		t.Errorf("expected 0 tokens, got %d", result.TotalTokens)
	}
}

func TestCollectFromSources_TerminalOutput(t *testing.T) {
	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)

	terminalOutput := "Error: undefined variable 'foo' at line 42"

	result, err := manager.CollectFromSources(context.Background(), SourceRequest{
		ProjectRoot:    "/tmp/test",
		Sources:        []SourceType{SourceTerminal},
		TerminalOutput: terminalOutput,
		MaxTokens:      10000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	item := result.Items[0]
	if item.Source != SourceTerminal {
		t.Errorf("expected source %s, got %s", SourceTerminal, item.Source)
	}

	if item.Content != terminalOutput {
		t.Errorf("expected content %q, got %q", terminalOutput, item.Content)
	}

	if item.Relevance != 0.95 {
		t.Errorf("expected relevance 0.95, got %f", item.Relevance)
	}
}

func TestCollectFromSources_GitDiff(t *testing.T) {
	gitRepo := &ctxSourceMockGitRepo{
		isGitRepo: true,
		uncommittedFiles: []domain.FileStatus{
			{Path: "main.go", Status: "M"},
			{Path: "new_file.go", Status: "A"},
		},
		diff: "diff --git a/main.go b/main.go\n+added line",
	}

	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, gitRepo)

	result, err := manager.CollectFromSources(context.Background(), SourceRequest{
		ProjectRoot: "/tmp/test",
		Sources:     []SourceType{SourceGitDiff},
		MaxTokens:   10000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	item := result.Items[0]
	if item.Source != SourceGitDiff {
		t.Errorf("expected source %s, got %s", SourceGitDiff, item.Source)
	}

	if item.Metadata["changedFiles"] != "2" {
		t.Errorf("expected 2 changed files in metadata, got %s", item.Metadata["changedFiles"])
	}
}

func TestCollectFromSources_GitHistory(t *testing.T) {
	gitRepo := &ctxSourceMockGitRepo{
		isGitRepo: true,
		commitHistory: []domain.CommitInfo{
			{Hash: "abc123def456", Subject: "feat: add feature", Author: "Test User", Date: "2024-01-15"},
			{Hash: "def456abc789", Subject: "fix: bug fix", Author: "Test User", Date: "2024-01-14"},
		},
	}

	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, gitRepo)

	result, err := manager.CollectFromSources(context.Background(), SourceRequest{
		ProjectRoot: "/tmp/test",
		Sources:     []SourceType{SourceGitHistory},
		MaxTokens:   10000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}

	item := result.Items[0]
	if item.Source != SourceGitHistory {
		t.Errorf("expected source %s, got %s", SourceGitHistory, item.Source)
	}

	if item.Metadata["commits"] != "2" {
		t.Errorf("expected 2 commits in metadata, got %s", item.Metadata["commits"])
	}
}

func TestCalculateSourceBudgets(t *testing.T) {
	tests := []struct {
		name      string
		sources   []SourceType
		maxTokens int
		wantCode  bool
	}{
		{
			name:      "single source gets full budget",
			sources:   []SourceType{SourceCode},
			maxTokens: 10000,
			wantCode:  true,
		},
		{
			name:      "multiple sources split budget",
			sources:   []SourceType{SourceCode, SourceTests, SourceGitDiff},
			maxTokens: 10000,
			wantCode:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)
			budgets := manager.calculateSourceBudgets(tt.sources, tt.maxTokens)

			totalBudget := 0
			for _, budget := range budgets {
				totalBudget += budget
			}

			// Total budget should be close to maxTokens (may differ due to rounding)
			if totalBudget > tt.maxTokens {
				t.Errorf("total budget %d exceeds maxTokens %d", totalBudget, tt.maxTokens)
			}

			if tt.wantCode {
				if _, ok := budgets[SourceCode]; !ok {
					t.Error("expected code source to have budget")
				}
			}
		})
	}
}

func TestTerminalOutputStorage(t *testing.T) {
	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)

	sessionID := "test-session-123"
	output := "Error: compilation failed"

	// Set terminal output
	manager.SetTerminalOutput(sessionID, output)

	// Get terminal output
	retrieved := manager.GetTerminalOutput(sessionID)
	if retrieved != output {
		t.Errorf("expected %q, got %q", output, retrieved)
	}

	// Clear terminal output
	manager.ClearTerminalOutput(sessionID)

	// Should be empty now
	retrieved = manager.GetTerminalOutput(sessionID)
	if retrieved != "" {
		t.Errorf("expected empty string, got %q", retrieved)
	}
}

func TestFilterCodeFiles(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected int
	}{
		{
			name:     "filters code files",
			files:    []string{"main.go", "utils.ts", "readme.md", "config.json"},
			expected: 2,
		},
		{
			name:     "excludes test files",
			files:    []string{"main.go", "main_test.go", "utils.spec.ts"},
			expected: 1,
		},
		{
			name:     "empty input",
			files:    []string{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterCodeFiles(tt.files)
			if len(result) != tt.expected {
				t.Errorf("expected %d files, got %d", tt.expected, len(result))
			}
		})
	}
}

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"main_test.go", true},
		{"utils.test.ts", true},
		{"component.spec.js", true},
		{"main.go", false},
		{"utils.ts", false},
		{"testing_utils.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isTestFile(tt.path)
			if result != tt.expected {
				t.Errorf("isTestFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsDocFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"README.md", true},
		{"docs/guide.md", true},
		{"CHANGELOG.md", true},
		{"CONTRIBUTING.md", true},
		{"notes.txt", true},
		{"main.go", false},
		{"config.json", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isDocFile(tt.path)
			if result != tt.expected {
				t.Errorf("isDocFile(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCtxSourceEstimateTokens(t *testing.T) {
	tests := []struct {
		text     string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"hello world", 2},  // 11 chars / 4 = 2
		{"a longer piece of text that should have more tokens", 12}, // 51 chars / 4 = 12
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			result := estimateTokens(tt.text)
			if result != tt.expected {
				t.Errorf("estimateTokens(%q) = %d, want %d", tt.text, result, tt.expected)
			}
		})
	}
}

func TestTruncateToTokens(t *testing.T) {
	tests := []struct {
		name      string
		text      string
		maxTokens int
	}{
		{
			name:      "no truncation needed",
			text:      "short",
			maxTokens: 100,
		},
		{
			name:      "truncation needed",
			text:      "this is a longer text that needs truncation because it exceeds the token limit",
			maxTokens: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateToTokens(tt.text, tt.maxTokens)
			// Just verify it doesn't panic and returns something
			if result == "" && tt.text != "" {
				t.Error("truncateToTokens returned empty string for non-empty input")
			}
		})
	}
}

func TestCalculateRelevance(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		query    string
		minScore float64
	}{
		{
			name:     "empty query",
			filePath: "main.go",
			query:    "",
			minScore: 0.5,
		},
		{
			name:     "matching query",
			filePath: "auth/login.go",
			query:    "auth login",
			minScore: 0.6,
		},
		{
			name:     "main file boost",
			filePath: "main.go",
			query:    "test",
			minScore: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateRelevance(tt.filePath, tt.query)
			if result < tt.minScore {
				t.Errorf("calculateRelevance(%q, %q) = %f, want >= %f",
					tt.filePath, tt.query, result, tt.minScore)
			}
		})
	}
}

func TestFormatSourceResult(t *testing.T) {
	result := &SourceResult{
		Items: []SourceItem{
			{
				Source:    SourceCode,
				Path:      "main.go",
				Content:   "package main",
				Tokens:    3,
				Relevance: 0.8,
			},
			{
				Source:    SourceTerminal,
				Path:      "terminal_output",
				Content:   "Error: test failed",
				Tokens:    4,
				Relevance: 0.95,
			},
		},
		TotalTokens: 7,
		Sources: map[SourceType]int{
			SourceCode:     3,
			SourceTerminal: 4,
		},
	}

	formatted := FormatSourceResult(result)

	// Check that it contains expected sections
	if !ctxSourceContains(formatted, "# Code Files") {
		t.Error("expected Code Files section")
	}

	if !ctxSourceContains(formatted, "# Terminal Output") {
		t.Error("expected Terminal Output section")
	}

	if !ctxSourceContains(formatted, "main.go") {
		t.Error("expected main.go in output")
	}
}

func TestCollectFromSources_WithRealFiles(t *testing.T) {
	// Create temp directory with test files
	tmpDir, err := os.MkdirTemp("", "context-sources-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := map[string]string{
		"main.go":       "package main\n\nfunc main() {}",
		"utils.go":      "package main\n\nfunc helper() {}",
		"main_test.go":  "package main\n\nfunc TestMain(t *testing.T) {}",
		"README.md":     "# Test Project\n\nThis is a test.",
		"docs/guide.md": "# Guide\n\nHow to use.",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	manager := NewContextSourceManager(&ctxSourceMockLogger{}, nil, nil)

	// Test collecting docs
	result, err := manager.CollectFromSources(context.Background(), SourceRequest{
		ProjectRoot: tmpDir,
		Sources:     []SourceType{SourceDocs},
		MaxTokens:   10000,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Items) < 1 {
		t.Error("expected at least 1 doc item")
	}

	// Check that README is included
	hasReadme := false
	for _, item := range result.Items {
		if ctxSourceContains(item.Path, "README") {
			hasReadme = true
			break
		}
	}

	if !hasReadme {
		t.Error("expected README.md to be included")
	}
}

func ctxSourceContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && ctxSourceContainsHelper(s, substr))
}

func ctxSourceContainsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
