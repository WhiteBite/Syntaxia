package textutils

import (
	"context"
	"syntaxia/domain/analysis"
	"strings"
	"testing"
)

// mockAnalyzer implements analysis.LanguageAnalyzer for testing
type mockAnalyzer struct {
	lang    string
	exts    []string
	symbols []analysis.Symbol
	imports []analysis.Import
}

func (m *mockAnalyzer) Language() string     { return m.lang }
func (m *mockAnalyzer) Extensions() []string { return m.exts }
func (m *mockAnalyzer) CanAnalyze(filePath string) bool {
	for _, ext := range m.exts {
		if strings.HasSuffix(filePath, ext) {
			return true
		}
	}
	return false
}

func (m *mockAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	return m.symbols, nil
}

func (m *mockAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	return m.imports, nil
}

func (m *mockAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	return nil, nil
}

func (m *mockAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	return "", 0, 0, nil
}

// mockRegistry implements AnalyzerRegistry for testing
type mockRegistry struct {
	analyzers map[string]analysis.LanguageAnalyzer
}

func (r *mockRegistry) GetAnalyzer(filePath string) analysis.LanguageAnalyzer {
	for ext, analyzer := range r.analyzers {
		if strings.HasSuffix(filePath, ext) {
			return analyzer
		}
	}
	return nil
}

func TestSkeletonGenerator_Generate_Go(t *testing.T) {
	analyzer := &mockAnalyzer{
		lang: "go",
		exts: []string{".go"},
		symbols: []analysis.Symbol{
			{Name: "main", Kind: analysis.KindPackage},
			{Name: "User", Kind: analysis.KindStruct},
			{Name: "UserService", Kind: analysis.KindInterface},
			{Name: "NewUser", Kind: analysis.KindFunction, Signature: "func NewUser(id int) *User"},
			{Name: "GetID", Kind: analysis.KindMethod, Parent: "User", Signature: "func (u *User) GetID() int"},
			{Name: "MaxSize", Kind: analysis.KindConstant},
		},
		imports: []analysis.Import{
			{Path: "fmt"},
			{Path: "context"},
		},
	}

	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".go": analyzer},
	}

	gen := NewSkeletonGenerator(registry)
	ctx := context.Background()

	result := gen.Generate(ctx, "package main\n...", "main.go")

	// Check package
	if !strings.Contains(result, "package main") {
		t.Error("should contain package declaration")
	}

	// Check imports
	if !strings.Contains(result, "import") {
		t.Error("should contain imports")
	}

	// Check struct
	if !strings.Contains(result, "type User struct") {
		t.Error("should contain User struct")
	}

	// Check interface
	if !strings.Contains(result, "type UserService interface") {
		t.Error("should contain UserService interface")
	}

	// Check function
	if !strings.Contains(result, "func NewUser") {
		t.Error("should contain NewUser function")
	}

	// Check method reference
	if !strings.Contains(result, "GetID") {
		t.Error("should reference GetID method")
	}
}

func TestSkeletonGenerator_Generate_TypeScript(t *testing.T) {
	analyzer := &mockAnalyzer{
		lang: "typescript",
		exts: []string{".ts"},
		symbols: []analysis.Symbol{
			{Name: "User", Kind: analysis.KindInterface},
			{Name: "UserService", Kind: analysis.KindClass},
			{Name: "createUser", Kind: analysis.KindFunction, Signature: "function createUser(name: string): User"},
		},
		imports: []analysis.Import{
			{Path: "./types", Names: []string{"User", "Config"}},
		},
	}

	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".ts": analyzer},
	}

	gen := NewSkeletonGenerator(registry)
	ctx := context.Background()

	result := gen.Generate(ctx, "...", "service.ts")

	if !strings.Contains(result, "import") {
		t.Error("should contain imports")
	}

	if !strings.Contains(result, "interface User") {
		t.Error("should contain User interface")
	}

	if !strings.Contains(result, "class UserService") {
		t.Error("should contain UserService class")
	}

	if !strings.Contains(result, "function createUser") {
		t.Error("should contain createUser function")
	}
}

func TestSkeletonGenerator_Generate_Python(t *testing.T) {
	analyzer := &mockAnalyzer{
		lang: "python",
		exts: []string{".py"},
		symbols: []analysis.Symbol{
			{Name: "UserService", Kind: analysis.KindClass},
			{Name: "create_user", Kind: analysis.KindFunction, Signature: "def create_user(name: str) -> User"},
		},
		imports: []analysis.Import{
			{Path: "typing", Names: []string{"Optional", "List"}},
		},
	}

	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".py": analyzer},
	}

	gen := NewSkeletonGenerator(registry)
	ctx := context.Background()

	result := gen.Generate(ctx, "...", "service.py")

	if !strings.Contains(result, "from typing import") {
		t.Error("should contain imports")
	}

	if !strings.Contains(result, "class UserService") {
		t.Error("should contain UserService class")
	}

	if !strings.Contains(result, "def create_user") {
		t.Error("should contain create_user function")
	}
}

func TestSkeletonGenerator_NoAnalyzer(t *testing.T) {
	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{},
	}

	gen := NewSkeletonGenerator(registry)
	ctx := context.Background()

	result := gen.Generate(ctx, "content", "file.unknown")

	if result != "" {
		t.Error("should return empty string for unsupported file")
	}
}

func TestSkeletonGenerator_NilRegistry(t *testing.T) {
	gen := NewSkeletonGenerator(nil)
	ctx := context.Background()

	result := gen.Generate(ctx, "content", "main.go")

	if result != "" {
		t.Error("should return empty string with nil registry")
	}
}

func TestSkeletonGenerator_CanGenerateSkeleton(t *testing.T) {
	analyzer := &mockAnalyzer{lang: "go", exts: []string{".go"}}
	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".go": analyzer},
	}

	gen := NewSkeletonGenerator(registry)

	if !gen.CanGenerateSkeleton("main.go") {
		t.Error("should be able to generate skeleton for .go files")
	}

	if gen.CanGenerateSkeleton("file.unknown") {
		t.Error("should not be able to generate skeleton for unknown files")
	}
}

func TestIsSkeletonSupported(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"app.ts", true},
		{"app.tsx", true},
		{"app.js", true},
		{"app.jsx", true},
		{"main.py", true},
		{"Main.java", true},
		{"Main.kt", true},
		{"main.rs", true},
		{"Main.cs", true},
		{"App.vue", true},
		{"main.dart", true},
		{"file.txt", false},
		{"file.md", false},
		{"file.json", false},
	}

	for _, tt := range tests {
		result := IsSkeletonSupported(tt.path)
		if result != tt.expected {
			t.Errorf("IsSkeletonSupported(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestSupportedSkeletonExtensions(t *testing.T) {
	exts := SupportedSkeletonExtensions()

	if len(exts) == 0 {
		t.Error("should return supported extensions")
	}

	// Check some expected extensions
	expected := []string{".go", ".ts", ".py", ".java"}
	for _, exp := range expected {
		found := false
		for _, ext := range exts {
			if ext == exp {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %q in supported extensions", exp)
		}
	}
}

func BenchmarkSkeletonGenerator_Generate(b *testing.B) {
	analyzer := &mockAnalyzer{
		lang: "go",
		exts: []string{".go"},
		symbols: []analysis.Symbol{
			{Name: "main", Kind: analysis.KindPackage},
			{Name: "User", Kind: analysis.KindStruct},
			{Name: "Config", Kind: analysis.KindStruct},
			{Name: "Service", Kind: analysis.KindInterface},
			{Name: "NewUser", Kind: analysis.KindFunction},
			{Name: "NewConfig", Kind: analysis.KindFunction},
			{Name: "GetID", Kind: analysis.KindMethod, Parent: "User"},
			{Name: "GetName", Kind: analysis.KindMethod, Parent: "User"},
			{Name: "MaxSize", Kind: analysis.KindConstant},
			{Name: "DefaultTimeout", Kind: analysis.KindConstant},
		},
		imports: []analysis.Import{
			{Path: "fmt"},
			{Path: "context"},
			{Path: "time"},
		},
	}

	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".go": analyzer},
	}

	gen := NewSkeletonGenerator(registry)
	ctx := context.Background()
	content := "package main\n..."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gen.Generate(ctx, content, "main.go")
	}
}
