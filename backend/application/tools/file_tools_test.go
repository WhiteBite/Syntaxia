package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchFiles_MatchesPattern(t *testing.T) {
	// Create temp directory with test files
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "utils.go"), []byte("package utils"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "readme.md"), []byte("# README"), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "main.go") || !contains(result, "utils.go") {
		t.Errorf("expected to find .go files, got: %s", result)
	}
}

func TestSearchFiles_NoMatches(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("search_files", map[string]any{"pattern": "*.xyz"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "No files found") {
		t.Errorf("expected 'No files found' message, got: %s", result)
	}
}

func TestReadFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	content := "line1\nline2\nline3\nline4\nline5"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("read_file", map[string]any{"path": "test.txt"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "line1") || !contains(result, "line5") {
		t.Errorf("expected file content, got: %s", result)
	}
}

func TestReadFile_NotFound(t *testing.T) {
	tmpDir := t.TempDir()

	handler := NewFileToolsHandler(nil, nil)
	_, err := handler.Execute("read_file", map[string]any{"path": "nonexistent.txt"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestReadFile_WithLineRange(t *testing.T) {
	tmpDir := t.TempDir()
	content := "line1\nline2\nline3\nline4\nline5"
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(content), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("read_file", map[string]any{
		"path":       "test.txt",
		"start_line": float64(2),
		"end_line":   float64(4),
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "line2") || !contains(result, "line4") {
		t.Errorf("expected lines 2-4, got: %s", result)
	}
}

func TestListDirectory_Flat(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("list_directory", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "file1.txt") || !contains(result, "subdir") {
		t.Errorf("expected directory contents, got: %s", result)
	}
}

func TestGetFileInfo_Success(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("get_file_info", map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "test.go") || !contains(result, ".go") {
		t.Errorf("expected file info, got: %s", result)
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

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("list_functions", map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "main") || !contains(result, "helper") {
		t.Errorf("expected function names, got: %s", result)
	}
}

func TestSearchContent_MatchesPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("func TestFunction() {\n\treturn nil\n}"), 0644)

	handler := NewFileToolsHandler(nil, nil)
	result, err := handler.Execute("search_content", map[string]any{"pattern": "TestFunction"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(result, "TestFunction") {
		t.Errorf("expected to find pattern, got: %s", result)
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr)))
}
