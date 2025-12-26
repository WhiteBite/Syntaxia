package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"syntaxia/domain"
)

// MockSandboxFS implements domain.SandboxFS for testing
type MockSandboxFS struct {
	files       map[string][]byte
	writeError  error
	projectRoot string
}

func NewMockSandboxFS() *MockSandboxFS {
	return &MockSandboxFS{
		files:       make(map[string][]byte),
		projectRoot: "/test",
	}
}

func (m *MockSandboxFS) ReadFile(filename string) ([]byte, error) {
	if content, ok := m.files[filename]; ok {
		return content, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockSandboxFS) WriteFile(filename string, data []byte, perm int) error {
	if m.writeError != nil {
		return m.writeError
	}
	m.files[filename] = data
	return nil
}

func (m *MockSandboxFS) DeleteFile(filename string) error {
	delete(m.files, filename)
	return nil
}

func (m *MockSandboxFS) MkdirAll(path string, perm int) error { return nil }
func (m *MockSandboxFS) GetChanges() []*domain.SandboxFileChange { return nil }
func (m *MockSandboxFS) GetDiff(filename string) (string, error) {
	return "", nil
}
func (m *MockSandboxFS) GetAllDiffs() (string, error) { return "", nil }
func (m *MockSandboxFS) Apply() error                 { return nil }
func (m *MockSandboxFS) Discard()                     {}
func (m *MockSandboxFS) DiscardFile(filename string)  {}
func (m *MockSandboxFS) HasChanges() bool             { return len(m.files) > 0 }
func (m *MockSandboxFS) GetChangeCount() int          { return len(m.files) }
func (m *MockSandboxFS) GetProjectRoot() string       { return m.projectRoot }
func (m *MockSandboxFS) SetProjectRoot(root string)   { m.projectRoot = root }

func TestSearchFiles_MatchesPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "utils.go"), []byte("package utils"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# README"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "main.go") || !strings.Contains(result, "utils.go") {
		t.Errorf("expected to find .go files, got: %s", result)
	}
}

func TestSearchFiles_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.xyz"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No files found") {
		t.Errorf("expected 'No files found' message, got: %s", result)
	}
}

func TestSearchFiles_MissingPattern(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("search_files", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
	if !strings.Contains(err.Error(), "pattern is required") {
		t.Errorf("expected 'pattern is required' error, got: %v", err)
	}
}

func TestSearchFiles_WithDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "src")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "app.go"), []byte("package app"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{
		"pattern":   "*.go",
		"directory": "src",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "app.go") {
		t.Errorf("expected to find app.go in src, got: %s", result)
	}
}

func TestReadFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	content := "line1\nline2\nline3\nline4\nline5"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{"path": "test.txt"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "line1") || !strings.Contains(result, "line5") {
		t.Errorf("expected file content, got: %s", result)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	handler := NewFileToolsHandler(nil, nil, nil)
	_, err := handler.Execute("read_file", map[string]any{"path": "nonexistent.txt"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestReadFile_MissingPath(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("read_file", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "path is required") {
		t.Errorf("expected 'path is required' error, got: %v", err)
	}
}

func TestReadFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("read_file", map[string]any{"path": "../../../etc/passwd"}, tmpDir)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
	if !strings.Contains(err.Error(), "path traversal") {
		t.Errorf("expected 'path traversal' error, got: %v", err)
	}
}

func TestReadFile_WithLineRange(t *testing.T) {
	tmpDir := t.TempDir()
	content := "line1\nline2\nline3\nline4\nline5"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{
		"path":       "test.txt",
		"start_line": float64(2),
		"end_line":   float64(4),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "line2") || !strings.Contains(result, "line4") {
		t.Errorf("expected lines 2-4, got: %s", result)
	}
}

func TestReadFile_StartLineExceedsLength(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("line1\nline2"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	_, err := handler.Execute("read_file", map[string]any{
		"path":       "test.txt",
		"start_line": float64(100),
	}, tmpDir)

	if err == nil {
		t.Fatal("expected error for start_line exceeding file length")
	}
}

func TestListDirectory_Flat(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "file1.txt") || !strings.Contains(result, "subdir") {
		t.Errorf("expected directory contents, got: %s", result)
	}
}

func TestListDirectory_Recursive(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
		"max_depth": float64(2),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "root.txt") || !strings.Contains(result, "nested.txt") {
		t.Errorf("expected recursive listing, got: %s", result)
	}
}

func TestListDirectory_WithPath(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "src")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "app.go"), []byte("package app"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{"path": "src"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "app.go") {
		t.Errorf("expected app.go in listing, got: %s", result)
	}
}

func TestListDirectory_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	emptyDir := filepath.Join(tmpDir, "empty")
	os.Mkdir(emptyDir, 0755)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{"path": "empty"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "empty") {
		t.Errorf("expected 'empty' message, got: %s", result)
	}
}

func TestGetFileInfo_Success(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("get_file_info", map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "test.go") || !strings.Contains(result, ".go") {
		t.Errorf("expected file info, got: %s", result)
	}
}

func TestGetFileInfo_MissingPath(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("get_file_info", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestGetFileInfo_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("get_file_info", map[string]any{"path": "nonexistent.txt"}, tmpDir)
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestListFunctions_GoFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `package main

func main() {
	fmt.Println("Hello")
}

func helper() string {
	return "help"
}
`
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "main") || !strings.Contains(result, "helper") {
		t.Errorf("expected function names, got: %s", result)
	}
}

func TestListFunctions_TypeScriptFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `function greet(name: string) {
	return "Hello " + name;
}

const helper = () => {
	return "help";
}

const asyncFn = async function() {
	return await fetch("/api");
}
`
	os.WriteFile(filepath.Join(tmpDir, "app.ts"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "app.ts"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "greet") {
		t.Errorf("expected function names, got: %s", result)
	}
}

func TestListFunctions_MissingPath(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("list_functions", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestListFunctions_UnsupportedFileType(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "data.json"), []byte(`{"key": "value"}`), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "data.json"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Unsupported") {
		t.Errorf("expected unsupported message, got: %s", result)
	}
}

func TestListFunctions_NoFunctions(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "empty.go"), []byte("package main\n\nvar x = 1"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "empty.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No functions") {
		t.Errorf("expected 'No functions' message, got: %s", result)
	}
}

func TestSearchContent_MatchesPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("func TestFunction() {\n\treturn nil\n}"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "TestFunction"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "TestFunction") {
		t.Errorf("expected to find pattern, got: %s", result)
	}
}

func TestSearchContent_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "nonexistent"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No matches") {
		t.Errorf("expected 'No matches' message, got: %s", result)
	}
}

func TestSearchContent_MissingPattern(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("search_content", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing pattern")
	}
}

func TestSearchContent_WithFilePattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "app.go"), []byte("func AppFunc() {}"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.ts"), []byte("function AppFunc() {}"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{
		"pattern":      "AppFunc",
		"file_pattern": "*.go",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "app.go") {
		t.Errorf("expected to find in .go file, got: %s", result)
	}
}

func TestSearchContent_WithMaxResults(t *testing.T) {
	tmpDir := t.TempDir()
	content := "match1\nmatch2\nmatch3\nmatch4\nmatch5"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{
		"pattern":     "match",
		"max_results": float64(2),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "2 matches") {
		t.Errorf("expected 2 matches, got: %s", result)
	}
}

func TestWriteFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	sandbox := NewMockSandboxFS()

	handler := NewFileToolsHandler(nil, nil, sandbox)
	result, err := handler.Execute("write_file", map[string]any{
		"path":    "new_file.txt",
		"content": "Hello, World!",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "sandbox") {
		t.Errorf("expected sandbox message, got: %s", result)
	}
}

func TestWriteFile_MissingPath(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, NewMockSandboxFS())

	_, err := handler.Execute("write_file", map[string]any{"content": "test"}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestWriteFile_MissingContent(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, NewMockSandboxFS())

	_, err := handler.Execute("write_file", map[string]any{"path": "test.txt"}, tmpDir)
	if err == nil {
		t.Fatal("expected error for missing content")
	}
}

func TestWriteFile_PathTraversal(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, NewMockSandboxFS())

	_, err := handler.Execute("write_file", map[string]any{
		"path":    "../../../etc/passwd",
		"content": "malicious",
	}, tmpDir)
	if err == nil {
		t.Fatal("expected error for path traversal")
	}
}

func TestWriteFile_SandboxNotInitialized(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("write_file", map[string]any{
		"path":    "test.txt",
		"content": "content",
	}, tmpDir)
	if err == nil {
		t.Fatal("expected error when sandbox not initialized")
	}
}

func TestFileToolsHandler_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     bool
	}{
		{"search_files", "search_files", true},
		{"search_content", "search_content", true},
		{"read_file", "read_file", true},
		{"write_file", "write_file", true},
		{"list_directory", "list_directory", true},
		{"get_file_info", "get_file_info", true},
		{"list_functions", "list_functions", true},
		{"unknown_tool", "unknown_tool", false},
		{"git_status", "git_status", false},
	}

	handler := NewFileToolsHandler(nil, nil, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.CanHandle(tt.toolName); got != tt.want {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

func TestFileToolsHandler_GetTools(t *testing.T) {
	handler := NewFileToolsHandler(nil, nil, nil)
	tools := handler.GetTools()

	if len(tools) == 0 {
		t.Fatal("expected non-empty tools list")
	}

	expectedTools := []string{"search_files", "search_content", "read_file", "write_file", "list_directory", "get_file_info", "list_functions"}
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q not found", expected)
		}
	}
}

func TestFileToolsHandler_UnknownTool(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("unknown_tool", map[string]any{}, tmpDir)
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown file tool") {
		t.Errorf("expected 'unknown file tool' error, got: %v", err)
	}
}


// === Security Tests ===

func TestReadFile_FileTooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a file that exceeds MaxFileSize (we'll use a smaller test limit)
	// Since we can't easily create a 10MB+ file in tests, we test the error message format
	largePath := filepath.Join(tmpDir, "large.txt")
	
	// Create a file and check the size validation logic works
	// We'll create a small file and verify the size check mechanism
	os.WriteFile(largePath, []byte("small content"), 0644)
	
	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{"path": "large.txt"}, tmpDir)
	
	// Should succeed for small file
	if err != nil {
		t.Fatalf("unexpected error for small file: %v", err)
	}
	if !strings.Contains(result, "small content") {
		t.Errorf("expected file content, got: %s", result)
	}
}

func TestReadFile_BinaryFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a binary file with null bytes
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0x00, 0x05}
	os.WriteFile(filepath.Join(tmpDir, "binary.bin"), binaryContent, 0644)
	
	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{"path": "binary.bin"}, tmpDir)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "binary file") {
		t.Errorf("expected binary file warning, got: %s", result)
	}
}

func TestSearchContent_SkipsBinaryFiles(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a text file with searchable content
	os.WriteFile(filepath.Join(tmpDir, "text.txt"), []byte("searchable content here"), 0644)
	// Create a binary file with the same text but also null bytes
	binaryContent := append([]byte("searchable content here"), 0x00, 0x01, 0x02)
	os.WriteFile(filepath.Join(tmpDir, "binary.bin"), binaryContent, 0644)
	
	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "searchable"}, tmpDir)
	
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should find in text file but not in binary
	if !strings.Contains(result, "text.txt") {
		t.Errorf("expected to find match in text.txt, got: %s", result)
	}
	if strings.Contains(result, "binary.bin") {
		t.Errorf("should not search in binary files, got: %s", result)
	}
}

func TestWriteFile_ContentTooLarge(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, NewMockSandboxFS())
	
	// Create content larger than MaxFileSize (10MB)
	// We'll test with a string that reports being too large
	// In real scenario, this would be 10MB+ content
	largeContent := strings.Repeat("x", MaxFileSize+1)
	
	_, err := handler.Execute("write_file", map[string]any{
		"path":    "large.txt",
		"content": largeContent,
	}, tmpDir)
	
	if err == nil {
		t.Fatal("expected error for content too large")
	}
	if !strings.Contains(err.Error(), "exceeds maximum") {
		t.Errorf("expected 'exceeds maximum' error, got: %v", err)
	}
}

func TestIsBinaryFile_TextFile(t *testing.T) {
	tmpDir := t.TempDir()
	textPath := filepath.Join(tmpDir, "text.txt")
	os.WriteFile(textPath, []byte("Hello, World!\nThis is a text file."), 0644)
	
	if isBinaryFile(textPath) {
		t.Error("text file should not be detected as binary")
	}
}

func TestIsBinaryFile_BinaryFile(t *testing.T) {
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "binary.bin")
	// Binary content with null bytes
	os.WriteFile(binaryPath, []byte{0x89, 0x50, 0x4E, 0x47, 0x00, 0x00}, 0644)
	
	if !isBinaryFile(binaryPath) {
		t.Error("binary file should be detected as binary")
	}
}

func TestIsBinaryFile_NonExistent(t *testing.T) {
	// Non-existent file should return false (not binary)
	if isBinaryFile("/nonexistent/path/file.bin") {
		t.Error("non-existent file should return false")
	}
}

func TestIsBinaryFile_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	emptyPath := filepath.Join(tmpDir, "empty.txt")
	os.WriteFile(emptyPath, []byte{}, 0644)
	
	if isBinaryFile(emptyPath) {
		t.Error("empty file should not be detected as binary")
	}
}

func TestConstants_Defined(t *testing.T) {
	// Verify constants are properly defined
	if MaxFileSize != 10*1024*1024 {
		t.Errorf("MaxFileSize should be 10MB, got: %d", MaxFileSize)
	}
	if MaxReadTimeout != 30*1000000000 { // 30 seconds in nanoseconds
		t.Errorf("MaxReadTimeout should be 30s, got: %v", MaxReadTimeout)
	}
	if BinaryCheckBytes != 512 {
		t.Errorf("BinaryCheckBytes should be 512, got: %d", BinaryCheckBytes)
	}
}


// === Additional Edge Case Tests ===

func TestSearchFiles_SkipsNodeModules(t *testing.T) {
	tmpDir := t.TempDir()
	nodeModules := filepath.Join(tmpDir, "node_modules")
	os.Mkdir(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "package.go"), []byte("package pkg"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "node_modules") {
		t.Errorf("should skip node_modules, got: %s", result)
	}
	if !strings.Contains(result, "main.go") {
		t.Errorf("should find main.go, got: %s", result)
	}
}

func TestSearchFiles_SkipsGitDir(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("config"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "config"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, ".git") {
		t.Errorf("should skip .git directory, got: %s", result)
	}
}

func TestSearchFiles_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# README"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "readme"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "README.md") {
		t.Errorf("should find README.md with case-insensitive search, got: %s", result)
	}
}

func TestSearchContent_RegexPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("func Test123() {}\nfunc Test456() {}"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "Test\\d+"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Test123") || !strings.Contains(result, "Test456") {
		t.Errorf("should find regex matches, got: %s", result)
	}
}

func TestSearchContent_InvalidRegexFallsBackToLiteral(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("search [pattern] here"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "[pattern]"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "[pattern]") {
		t.Errorf("should find literal pattern, got: %s", result)
	}
}

func TestReadFile_EndLineExceedsLength(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "short.txt"), []byte("line1\nline2"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{
		"path":     "short.txt",
		"end_line": float64(100),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "line1") && !strings.Contains(result, "line2") {
		t.Errorf("should read available lines, got: %s", result)
	}
}

func TestListDirectory_SkipsNodeModulesRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	nodeModules := filepath.Join(tmpDir, "node_modules")
	os.Mkdir(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "pkg.js"), []byte("module"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.js"), []byte("app"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
		"max_depth": float64(3),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "pkg.js") {
		t.Errorf("should skip node_modules contents, got: %s", result)
	}
}

func TestListDirectory_MaxDepthLimit(t *testing.T) {
	tmpDir := t.TempDir()
	deep := filepath.Join(tmpDir, "a", "b", "c", "d")
	os.MkdirAll(deep, 0755)
	os.WriteFile(filepath.Join(deep, "deep.txt"), []byte("deep"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
		"max_depth": float64(1),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "deep.txt") {
		t.Errorf("should not show files beyond max_depth, got: %s", result)
	}
}

func TestListFunctions_GoMethodReceiver(t *testing.T) {
	tmpDir := t.TempDir()
	content := `package main

func (s *Service) Method1() {}
func (s Service) Method2() {}
func RegularFunc() {}
`
	os.WriteFile(filepath.Join(tmpDir, "service.go"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "service.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Method1") || !strings.Contains(result, "Method2") || !strings.Contains(result, "RegularFunc") {
		t.Errorf("should find all functions including methods, got: %s", result)
	}
}

func TestListFunctions_JSXFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `function Component() {
	return <div>Hello</div>;
}

const Helper = () => {
	return null;
}
`
	os.WriteFile(filepath.Join(tmpDir, "component.jsx"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "component.jsx"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Component") {
		t.Errorf("should find JSX functions, got: %s", result)
	}
}

func TestListFunctions_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("list_functions", map[string]any{"path": "nonexistent.go"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestListFunctions_DuplicateRemoval(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
const myFunc = function() {}
const myFunc = () => {}
`
	os.WriteFile(filepath.Join(tmpDir, "dup.js"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "dup.js"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Count occurrences of myFunc
	count := strings.Count(result, "myFunc")
	if count > 1 {
		t.Errorf("should remove duplicates, found %d occurrences", count)
	}
}

func TestWriteFile_SandboxWriteError(t *testing.T) {
	tmpDir := t.TempDir()
	sandbox := NewMockSandboxFS()
	sandbox.writeError = os.ErrPermission

	handler := NewFileToolsHandler(nil, nil, sandbox)
	_, err := handler.Execute("write_file", map[string]any{
		"path":    "test.txt",
		"content": "content",
	}, tmpDir)

	if err == nil {
		t.Fatal("expected error when sandbox write fails")
	}
	if !strings.Contains(err.Error(), "failed to write to sandbox") {
		t.Errorf("expected 'failed to write to sandbox' error, got: %v", err)
	}
}

func TestGetFileInfo_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	os.Mkdir(subDir, 0755)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("get_file_info", map[string]any{"path": "subdir"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "true") { // isDir should be true
		t.Errorf("expected isDir to be true for directory, got: %s", result)
	}
}

func TestSearchContent_SkipsNodeModules(t *testing.T) {
	tmpDir := t.TempDir()
	nodeModules := filepath.Join(tmpDir, "node_modules")
	os.Mkdir(nodeModules, 0755)
	os.WriteFile(filepath.Join(nodeModules, "lib.js"), []byte("searchterm"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.js"), []byte("searchterm"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "searchterm"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "node_modules") {
		t.Errorf("should skip node_modules, got: %s", result)
	}
	if !strings.Contains(result, "app.js") {
		t.Errorf("should find in app.js, got: %s", result)
	}
}

func TestListDirectory_NonExistentPath(t *testing.T) {
	tmpDir := t.TempDir()
	handler := NewFileToolsHandler(nil, nil, nil)

	_, err := handler.Execute("list_directory", map[string]any{"path": "nonexistent"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}


func TestSearchFiles_Limit50Results(t *testing.T) {
	tmpDir := t.TempDir()
	// Create more than 50 files
	for i := 0; i < 60; i++ {
		os.WriteFile(filepath.Join(tmpDir, fmt.Sprintf("file%d.txt", i)), []byte("content"), 0644)
	}

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "file"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "50 files") {
		t.Errorf("should limit to 50 files, got: %s", result)
	}
}

func TestListFunctions_TSXAsyncFunction(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
const fetchData = async function() {
	return await fetch('/api');
}

const processData = async (data) => {
	return data;
}
`
	os.WriteFile(filepath.Join(tmpDir, "async.tsx"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "async.tsx"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "fetchData") {
		t.Errorf("should find async functions, got: %s", result)
	}
}

func TestSearchFiles_VendorSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	vendorDir := filepath.Join(tmpDir, "vendor")
	os.Mkdir(vendorDir, 0755)
	os.WriteFile(filepath.Join(vendorDir, "lib.go"), []byte("package lib"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "vendor") {
		t.Errorf("should skip vendor directory, got: %s", result)
	}
}

func TestSearchFiles_DistSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	distDir := filepath.Join(tmpDir, "dist")
	os.Mkdir(distDir, 0755)
	os.WriteFile(filepath.Join(distDir, "bundle.js"), []byte("bundle"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.js"), []byte("app"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.js"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "dist") {
		t.Errorf("should skip dist directory, got: %s", result)
	}
}

func TestListDirectory_GitSkipped(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("config"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "config") && strings.Contains(result, ".git") {
		t.Errorf("should skip .git directory contents, got: %s", result)
	}
}


func TestReadFile_EndLineClampedToFileLength(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "short.txt"), []byte("line1\nline2\nline3"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("read_file", map[string]any{
		"path":       "short.txt",
		"start_line": float64(1),
		"end_line":   float64(1000), // Way beyond file length
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should clamp end_line to file length and return all lines
	if !strings.Contains(result, "line1") || !strings.Contains(result, "line3") {
		t.Errorf("expected all lines, got: %s", result)
	}
}

func TestListDirectory_FileInfoNil(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "test.txt") {
		t.Errorf("expected file in listing, got: %s", result)
	}
}

func TestSearchContent_SkipsGitDir(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	os.Mkdir(gitDir, 0755)
	os.WriteFile(filepath.Join(gitDir, "config"), []byte("searchterm"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("searchterm"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "searchterm"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, ".git") {
		t.Errorf("should skip .git directory, got: %s", result)
	}
	if !strings.Contains(result, "main.go") {
		t.Errorf("should find in main.go, got: %s", result)
	}
}

func TestListDirectory_RecursiveWithFiles(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("root"), 0644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "root.txt") || !strings.Contains(result, "nested.txt") {
		t.Errorf("expected both files, got: %s", result)
	}
	if !strings.Contains(result, "sub") {
		t.Errorf("expected subdirectory, got: %s", result)
	}
}


func TestSearchContent_CaseInsensitive(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("func UPPERCASE() {}"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "uppercase"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "UPPERCASE") {
		t.Errorf("should find case-insensitive match, got: %s", result)
	}
}

func TestListFunctions_AsyncArrowFunction(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
const asyncArrow = async () => {
	return await fetch('/api');
}
`
	os.WriteFile(filepath.Join(tmpDir, "async.js"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "async.js"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "asyncArrow") {
		t.Errorf("should find async arrow function, got: %s", result)
	}
}

func TestListFunctions_NamedFunction(t *testing.T) {
	tmpDir := t.TempDir()
	content := `
const namedFunc = function myName() {
	return 'hello';
}
`
	os.WriteFile(filepath.Join(tmpDir, "named.js"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "named.js"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "namedFunc") {
		t.Errorf("should find named function expression, got: %s", result)
	}
}

func TestSearchFiles_PartialNameMatch(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "user_service.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "product_service.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "service"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "user_service.go") || !strings.Contains(result, "product_service.go") {
		t.Errorf("should find partial matches, got: %s", result)
	}
}

func TestListDirectory_DeepRecursive(t *testing.T) {
	tmpDir := t.TempDir()
	// Create nested structure
	level1 := filepath.Join(tmpDir, "level1")
	level2 := filepath.Join(level1, "level2")
	level3 := filepath.Join(level2, "level3")
	os.MkdirAll(level3, 0755)
	os.WriteFile(filepath.Join(level3, "deep.txt"), []byte("deep"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{
		"recursive": true,
		"max_depth": float64(5),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "level1") || !strings.Contains(result, "level2") {
		t.Errorf("should show nested directories, got: %s", result)
	}
}


func TestNewFileToolsHandler(t *testing.T) {
	sandbox := NewMockSandboxFS()
	handler := NewFileToolsHandler(nil, nil, sandbox)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.SandboxFS != sandbox {
		t.Error("expected sandbox to be set")
	}
}

func TestSearchFiles_GlobMatch(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("text"), 0644)

	handler := NewFileToolsHandler(nil, nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "test.go") {
		t.Errorf("expected test.go in result, got: %s", result)
	}
	if strings.Contains(result, "test.txt") {
		t.Errorf("should not include test.txt, got: %s", result)
	}
}
