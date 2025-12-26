package context

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"syntaxia/domain"
)

// MockLogger implements domain.Logger for testing
type MockLogger struct {
	DebugMessages   []string
	InfoMessages    []string
	WarningMessages []string
	ErrorMessages   []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		DebugMessages:   make([]string, 0),
		InfoMessages:    make([]string, 0),
		WarningMessages: make([]string, 0),
		ErrorMessages:   make([]string, 0),
	}
}

func (m *MockLogger) Debug(message string)   { m.DebugMessages = append(m.DebugMessages, message) }
func (m *MockLogger) Info(message string)    { m.InfoMessages = append(m.InfoMessages, message) }
func (m *MockLogger) Warning(message string) { m.WarningMessages = append(m.WarningMessages, message) }
func (m *MockLogger) Error(message string)   { m.ErrorMessages = append(m.ErrorMessages, message) }
func (m *MockLogger) Fatal(message string)   { m.ErrorMessages = append(m.ErrorMessages, message) }

// MockFileContentReader implements domain.FileContentReader for testing
type MockFileContentReader struct {
	Contents map[string]string
	Err      error
}

func NewMockFileContentReader() *MockFileContentReader {
	return &MockFileContentReader{
		Contents: make(map[string]string),
	}
}

func (m *MockFileContentReader) ReadContents(
	_ context.Context,
	filePaths []string,
	_ string,
	_ func(current, total int64),
) (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	result := make(map[string]string)
	for _, path := range filePaths {
		if content, ok := m.Contents[path]; ok {
			result[path] = content
		}
	}
	return result, nil
}

// MockTreeBuilder implements domain.TreeBuilder for testing
type MockTreeBuilder struct {
	Tree []*domain.FileNode
	Err  error
}

func NewMockTreeBuilder() *MockTreeBuilder {
	return &MockTreeBuilder{
		Tree: make([]*domain.FileNode, 0),
	}
}

func (m *MockTreeBuilder) BuildTree(_ string, _, _ bool) ([]*domain.FileNode, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Tree, nil
}

func (m *MockTreeBuilder) InvalidateCache() {}

// setupTestProject creates a temporary project structure for testing
func setupTestProject(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()

	files := map[string]string{
		"main.go":            "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n",
		"auth/login.go":      "package auth\n\nfunc Login() error {\n\treturn nil\n}\n",
		"auth/logout.go":     "package auth\n\nfunc Logout() error {\n\treturn nil\n}\n",
		"api/handler.go":     "package api\n\nfunc Handler() {}\n",
		"config/settings.go": "package config\n\nvar Settings = map[string]string{}\n",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file %s: %v", fullPath, err)
		}
	}

	return tempDir
}

func TestNewSmartContextCollector(t *testing.T) {
	logger := NewMockLogger()
	fileReader := NewMockFileContentReader()
	treeBuilder := NewMockTreeBuilder()

	collector := NewSmartContextCollector(logger, fileReader, treeBuilder)

	if collector == nil {
		t.Fatal("expected non-nil collector")
	}
	if collector.log != logger {
		t.Error("logger not set correctly")
	}
	if collector.fileReader != fileReader {
		t.Error("fileReader not set correctly")
	}
	if collector.treeBuilder != treeBuilder {
		t.Error("treeBuilder not set correctly")
	}
}

func TestCollectContext(t *testing.T) {
	tests := []struct {
		name             string
		request          domain.SmartContextRequest
		setupProject     bool
		expectedStrategy string
		wantErr          bool
		checkFunc        func(t *testing.T, result *domain.SmartContextResult)
	}{
		{
			name: "with selected files uses selected_with_imports strategy",
			request: domain.SmartContextRequest{
				Task:          "fix authentication bug",
				SelectedFiles: []string{"auth/login.go"},
				MaxTokens:     10000,
			},
			setupProject:     true,
			expectedStrategy: "selected_with_imports",
			wantErr:          false,
			checkFunc: func(t *testing.T, result *domain.SmartContextResult) {
				if result.Strategy != "selected_with_imports" {
					t.Errorf("expected strategy 'selected_with_imports', got '%s'", result.Strategy)
				}
			},
		},
		{
			name: "without selected files uses task_analysis strategy",
			request: domain.SmartContextRequest{
				Task:      "fix authentication bug",
				MaxTokens: 10000,
			},
			setupProject:     true,
			expectedStrategy: "task_analysis",
			wantErr:          false,
			checkFunc: func(t *testing.T, result *domain.SmartContextResult) {
				if result.Strategy != "task_analysis" {
					t.Errorf("expected strategy 'task_analysis', got '%s'", result.Strategy)
				}
			},
		},
		{
			name: "default max tokens when not specified",
			request: domain.SmartContextRequest{
				Task:          "test task",
				SelectedFiles: []string{"main.go"},
				MaxTokens:     0,
			},
			setupProject: true,
			wantErr:      false,
			checkFunc: func(t *testing.T, result *domain.SmartContextResult) {
				if result == nil {
					t.Error("expected non-nil result")
				}
			},
		},
		{
			name: "project structure is included",
			request: domain.SmartContextRequest{
				Task:          "analyze project",
				SelectedFiles: []string{"main.go"},
				MaxTokens:     50000,
			},
			setupProject: true,
			wantErr:      false,
			checkFunc: func(t *testing.T, result *domain.SmartContextResult) {
				if result.ProjectStructure == "" {
					t.Error("expected non-empty project structure")
				}
			},
		},
		{
			name: "total tokens calculated correctly",
			request: domain.SmartContextRequest{
				Task:          "test",
				SelectedFiles: []string{"main.go"},
				MaxTokens:     50000,
			},
			setupProject: true,
			wantErr:      false,
			checkFunc: func(t *testing.T, result *domain.SmartContextResult) {
				if result.TotalTokens <= 0 {
					t.Error("expected positive total tokens")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)

			var projectRoot string
			if tt.setupProject {
				projectRoot = setupTestProject(t)
			} else {
				projectRoot = t.TempDir()
			}
			tt.request.ProjectRoot = projectRoot

			ctx := context.Background()
			result, err := collector.CollectContext(ctx, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("CollectContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && result == nil {
				t.Error("expected non-nil result when no error")
				return
			}

			if tt.checkFunc != nil && result != nil {
				tt.checkFunc(t, result)
			}
		})
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		name     string
		task     string
		expected []string
	}{
		{
			name:     "auth keywords",
			task:     "fix auth login bug",
			expected: []string{"auth", "login"},
		},
		{
			name:     "user keywords",
			task:     "update user profile page",
			expected: []string{"user", "profile"},
		},
		{
			name:     "api keywords",
			task:     "add new api handler for service",
			expected: []string{"api", "handler", "service"},
		},
		{
			name:     "test keywords",
			task:     "write unit test for spec",
			expected: []string{"test", "spec"},
		},
		{
			name:     "config keywords",
			task:     "update config settings",
			expected: []string{"config", "settings"},
		},
		{
			name:     "empty task",
			task:     "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
			keywords := collector.extractKeywords(tt.task)

			for _, expected := range tt.expected {
				found := false
				for _, kw := range keywords {
					if kw == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected keyword '%s' not found in %v", expected, keywords)
				}
			}
		})
	}
}

func TestExpandWithImports(t *testing.T) {
	tests := []struct {
		name        string
		files       []string
		maxExpected int
	}{
		{
			name:        "single file",
			files:       []string{"main.go"},
			maxExpected: 1,
		},
		{
			name:        "multiple files",
			files:       []string{"main.go", "auth/login.go", "api/handler.go"},
			maxExpected: 3,
		},
		{
			name:        "duplicate files deduplicated",
			files:       []string{"main.go", "main.go", "auth/login.go"},
			maxExpected: 2,
		},
		{
			name:        "empty files",
			files:       []string{},
			maxExpected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
			result := collector.expandWithImports(tt.files, "")

			if len(result) > tt.maxExpected {
				t.Errorf("expected at most %d files, got %d", tt.maxExpected, len(result))
			}

			seen := make(map[string]bool)
			for _, f := range result {
				if seen[f] {
					t.Errorf("duplicate file found: %s", f)
				}
				seen[f] = true
			}
		})
	}
}

func TestFindFilesByTask(t *testing.T) {
	tests := []struct {
		name          string
		task          string
		expectedFiles []string
	}{
		{
			name:          "find auth files",
			task:          "fix authentication login",
			expectedFiles: []string{"login.go"},
		},
		{
			name:          "find config files",
			task:          "update settings configuration",
			expectedFiles: []string{"settings.go"},
		},
		{
			name:          "find handler files",
			task:          "add api handler",
			expectedFiles: []string{"handler.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
			projectRoot := setupTestProject(t)

			files := collector.findFilesByTask(tt.task, projectRoot)

			for _, expected := range tt.expectedFiles {
				found := false
				for _, f := range files {
					if filepath.Base(f) == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected file '%s' not found in %v", expected, files)
				}
			}
		})
	}
}

func TestReadAndTruncate(t *testing.T) {
	tests := []struct {
		name           string
		files          []string
		maxTokens      int
		expectedCount  int
		checkTruncated bool
	}{
		{
			name:          "read single file",
			files:         []string{"main.go"},
			maxTokens:     10000,
			expectedCount: 1,
		},
		{
			name:          "read multiple files",
			files:         []string{"main.go", "auth/login.go"},
			maxTokens:     10000,
			expectedCount: 2,
		},
		{
			name:           "truncate when exceeding limit",
			files:          []string{"main.go"},
			maxTokens:      5,
			expectedCount:  1,
			checkTruncated: true,
		},
		{
			name:          "empty files list",
			files:         []string{},
			maxTokens:     10000,
			expectedCount: 0,
		},
		{
			name:          "non-existent file skipped",
			files:         []string{"nonexistent.go"},
			maxTokens:     10000,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
			projectRoot := setupTestProject(t)

			result := collector.readAndTruncate(tt.files, projectRoot, tt.maxTokens)

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d files, got %d", tt.expectedCount, len(result))
			}

			for _, f := range result {
				if f.Path == "" {
					t.Error("file path should not be empty")
				}
				if f.Reason == "" {
					t.Error("file reason should not be empty")
				}
			}
		})
	}
}

func TestBuildCompactStructure(t *testing.T) {
	tests := []struct {
		name          string
		checkContains []string
	}{
		{
			name:          "includes project structure header",
			checkContains: []string{"PROJECT STRUCTURE:"},
		},
		{
			name:          "includes directories",
			checkContains: []string{"auth/", "api/", "config/"},
		},
		{
			name:          "includes files",
			checkContains: []string{"main.go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewMockLogger()
			fileReader := NewMockFileContentReader()
			treeBuilder := NewMockTreeBuilder()

			collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
			projectRoot := setupTestProject(t)

			structure := collector.buildCompactStructure(projectRoot)

			for _, expected := range tt.checkContains {
				if !containsSubstr(structure, expected) {
					t.Errorf("expected structure to contain '%s'", expected)
				}
			}
		})
	}
}

func TestBuildCompactStructure_SkipsIgnoredDirs(t *testing.T) {
	tempDir := t.TempDir()

	ignoredDirs := []string{"node_modules", ".git", "vendor", "dist", "build", "__pycache__"}
	for _, dir := range ignoredDirs {
		dirPath := filepath.Join(tempDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
		filePath := filepath.Join(dirPath, "test.txt")
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	regularDir := filepath.Join(tempDir, "src")
	if err := os.MkdirAll(regularDir, 0755); err != nil {
		t.Fatalf("failed to create src dir: %v", err)
	}

	logger := NewMockLogger()
	fileReader := NewMockFileContentReader()
	treeBuilder := NewMockTreeBuilder()

	collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
	structure := collector.buildCompactStructure(tempDir)

	for _, dir := range ignoredDirs {
		if containsSubstr(structure, dir+"/") {
			t.Errorf("expected structure to NOT contain ignored dir '%s/'", dir)
		}
	}

	if !containsSubstr(structure, "src/") {
		t.Errorf("expected structure to contain 'src/'")
	}
}

func TestCollectContext_LogsInfo(t *testing.T) {
	logger := NewMockLogger()
	fileReader := NewMockFileContentReader()
	treeBuilder := NewMockTreeBuilder()

	collector := NewSmartContextCollector(logger, fileReader, treeBuilder)
	projectRoot := setupTestProject(t)

	req := domain.SmartContextRequest{
		Task:          "test task",
		ProjectRoot:   projectRoot,
		SelectedFiles: []string{"main.go"},
		MaxTokens:     10000,
	}

	ctx := context.Background()
	_, err := collector.CollectContext(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(logger.InfoMessages) == 0 {
		t.Error("expected at least one info message to be logged")
	}

	found := false
	for _, msg := range logger.InfoMessages {
		if containsSubstr(msg, "Smart context collected") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected log message about context collection, got: %v", logger.InfoMessages)
	}
}

func containsSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
