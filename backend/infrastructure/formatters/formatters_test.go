package formatters

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"syntaxia/domain"
)

type testLogger struct{}

func (l *testLogger) Debug(msg string)   {}
func (l *testLogger) Info(msg string)    {}
func (l *testLogger) Warning(msg string) {}
func (l *testLogger) Error(msg string)   {}
func (l *testLogger) Fatal(msg string)   {}

// =============================================================================
// Go Formatter Tests
// =============================================================================

func TestNewGoFormatter(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)
	if formatter == nil {
		t.Fatal("NewGoFormatter returned nil")
	}
}

func TestGoFormatter_GetSupportedLanguages(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)

	langs := formatter.GetSupportedLanguages()
	if len(langs) != 1 {
		t.Errorf("expected 1 language, got %d", len(langs))
	}
	if langs[0] != "go" {
		t.Errorf("expected 'go', got %q", langs[0])
	}
}

func TestGoFormatter_FormatFile_NotGoFile(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)

	err := formatter.FormatFile(context.Background(), "test.py")
	if err == nil {
		t.Error("expected error for non-Go file")
	}
	if !strings.Contains(err.Error(), "not a Go file") {
		t.Errorf("expected 'not a Go file' error, got: %v", err)
	}
}

func TestGoFormatter_FormatContent_UnsupportedLanguage(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)

	_, err := formatter.FormatContent(context.Background(), "code", "python")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
	if !strings.Contains(err.Error(), "unsupported language") {
		t.Errorf("expected 'unsupported language' error, got: %v", err)
	}
}

func TestGoFormatter_FormatFile(t *testing.T) {
	// Skip if gofmt not available
	if _, err := exec.LookPath("gofmt"); err != nil {
		t.Skip("gofmt not available")
	}

	log := &testLogger{}
	formatter := NewGoFormatter(log)

	// Create temp file with unformatted Go code
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.go")
	unformatted := `package main
func main(){println("hello")}`

	if err := os.WriteFile(tmpFile, []byte(unformatted), 0o644); err != nil {
		t.Fatal(err)
	}

	err := formatter.FormatFile(context.Background(), tmpFile)
	if err != nil {
		t.Fatalf("FormatFile failed: %v", err)
	}

	// Read formatted content
	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatal(err)
	}

	// Should be formatted (gofmt adds spaces)
	if !strings.Contains(string(content), "func main() {") {
		t.Error("file should be formatted")
	}
}

func TestGoFormatter_FormatContent(t *testing.T) {
	// Skip if gofmt not available
	if _, err := exec.LookPath("gofmt"); err != nil {
		t.Skip("gofmt not available")
	}

	log := &testLogger{}
	formatter := NewGoFormatter(log)

	unformatted := `package main
func main(){println("hello")}`

	result, err := formatter.FormatContent(context.Background(), unformatted, "go")
	if err != nil {
		t.Fatalf("FormatContent failed: %v", err)
	}

	if !strings.Contains(result, "func main() {") {
		t.Error("content should be formatted")
	}
}

func TestGoFormatter_FixImports_NotGoFile(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)

	err := formatter.FixImports(context.Background(), "test.py")
	if err == nil {
		t.Error("expected error for non-Go file")
	}
}

func TestGoFormatter_FixImportsInContent_UnsupportedLanguage(t *testing.T) {
	log := &testLogger{}
	formatter := NewGoFormatter(log)

	_, err := formatter.FixImportsInContent(context.Background(), "code", "python")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
}

// =============================================================================
// TypeScript Formatter Tests
// =============================================================================

func TestNewTypeScriptFormatter(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)
	if formatter == nil {
		t.Fatal("NewTypeScriptFormatter returned nil")
	}
}

func TestTypeScriptFormatter_GetSupportedLanguages(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	langs := formatter.GetSupportedLanguages()
	if len(langs) != 2 {
		t.Errorf("expected 2 languages, got %d", len(langs))
	}

	hasTS := false
	hasTypescript := false
	for _, l := range langs {
		if l == "ts" {
			hasTS = true
		}
		if l == "typescript" {
			hasTypescript = true
		}
	}

	if !hasTS || !hasTypescript {
		t.Error("expected 'ts' and 'typescript' in supported languages")
	}
}

func TestTypeScriptFormatter_FormatFile_NotTSFile(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	err := formatter.FormatFile(context.Background(), "test.go")
	if err == nil {
		t.Error("expected error for non-TypeScript file")
	}
	if !strings.Contains(err.Error(), "not a TypeScript file") {
		t.Errorf("expected 'not a TypeScript file' error, got: %v", err)
	}
}

func TestTypeScriptFormatter_FormatContent_UnsupportedLanguage(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	_, err := formatter.FormatContent(context.Background(), "code", "go")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
	if !strings.Contains(err.Error(), "unsupported language") {
		t.Errorf("expected 'unsupported language' error, got: %v", err)
	}
}

func TestTypeScriptFormatter_FixImports_NotTSFile(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	err := formatter.FixImports(context.Background(), "test.go")
	if err == nil {
		t.Error("expected error for non-TypeScript file")
	}
}

func TestTypeScriptFormatter_FixImportsInContent_UnsupportedLanguage(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	_, err := formatter.FixImportsInContent(context.Background(), "code", "go")
	if err == nil {
		t.Error("expected error for unsupported language")
	}
}

func TestTypeScriptFormatter_FindPackageJsonDir(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	// Create temp directory structure
	tmpDir := t.TempDir()
	srcDir := filepath.Join(tmpDir, "src")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Create package.json in root
	pkgJSON := filepath.Join(tmpDir, "package.json")
	if err := os.WriteFile(pkgJSON, []byte(`{"name": "test"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	// Test finding package.json from nested file
	testFile := filepath.Join(srcDir, "test.ts")
	result := formatter.findPackageJsonDir(testFile)

	if result != tmpDir {
		t.Errorf("expected %q, got %q", tmpDir, result)
	}
}

func TestTypeScriptFormatter_FindPackageJsonDir_NotFound(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	// Create temp directory without package.json
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.ts")

	result := formatter.findPackageJsonDir(testFile)
	if result != "" {
		t.Errorf("expected empty string when package.json not found, got %q", result)
	}
}

func TestTypeScriptFormatter_FormatFile_TSX(t *testing.T) {
	log := &testLogger{}
	formatter := NewTypeScriptFormatter(log)

	// Should accept .tsx files
	err := formatter.FormatFile(context.Background(), "/nonexistent/test.tsx")
	// Will fail because file doesn't exist, but shouldn't fail on extension check
	if err != nil && strings.Contains(err.Error(), "not a TypeScript file") {
		t.Error("should accept .tsx files")
	}
}

// =============================================================================
// Interface Compliance Tests
// =============================================================================

// Verify that formatters can be used with domain.Logger
func TestFormatters_AcceptDomainLogger(t *testing.T) {
	var log domain.Logger = &domain.NoopLogger{}

	goFormatter := NewGoFormatter(log)
	if goFormatter == nil {
		t.Error("GoFormatter should accept domain.Logger")
	}

	tsFormatter := NewTypeScriptFormatter(log)
	if tsFormatter == nil {
		t.Error("TypeScriptFormatter should accept domain.Logger")
	}
}
