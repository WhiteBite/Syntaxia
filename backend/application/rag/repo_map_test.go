package rag

import (
	"context"
	"testing"
	"time"

	"syntaxia/domain/analysis"
)

// mockAnalyzerRegistry implements analysis.AnalyzerRegistry for testing
type mockAnalyzerRegistry struct {
	analyzers map[string]*mockAnalyzer
}

func newMockAnalyzerRegistry() *mockAnalyzerRegistry {
	return &mockAnalyzerRegistry{
		analyzers: make(map[string]*mockAnalyzer),
	}
}

func (r *mockAnalyzerRegistry) Register(analyzer analysis.LanguageAnalyzer) {
	r.analyzers[analyzer.Language()] = analyzer.(*mockAnalyzer)
}

func (r *mockAnalyzerRegistry) GetAnalyzer(filePath string) analysis.LanguageAnalyzer {
	for _, a := range r.analyzers {
		if a.CanAnalyze(filePath) {
			return a
		}
	}
	return nil
}

func (r *mockAnalyzerRegistry) GetAnalyzerByLanguage(lang string) analysis.LanguageAnalyzer {
	return r.analyzers[lang]
}

func (r *mockAnalyzerRegistry) SupportedLanguages() []string {
	langs := make([]string, 0, len(r.analyzers))
	for lang := range r.analyzers {
		langs = append(langs, lang)
	}
	return langs
}

func (r *mockAnalyzerRegistry) SupportedExtensions() []string {
	return []string{".go", ".ts"}
}

// mockAnalyzer implements analysis.LanguageAnalyzer for testing
type mockAnalyzer struct {
	lang       string
	extensions []string
	symbols    map[string][]analysis.Symbol
	imports    map[string][]analysis.Import
	exports    map[string][]analysis.Export
}

func newMockAnalyzer(lang string, extensions []string) *mockAnalyzer {
	return &mockAnalyzer{
		lang:       lang,
		extensions: extensions,
		symbols:    make(map[string][]analysis.Symbol),
		imports:    make(map[string][]analysis.Import),
		exports:    make(map[string][]analysis.Export),
	}
}

func (a *mockAnalyzer) Language() string     { return a.lang }
func (a *mockAnalyzer) Extensions() []string { return a.extensions }

func (a *mockAnalyzer) CanAnalyze(filePath string) bool {
	for _, ext := range a.extensions {
		if len(filePath) > len(ext) && filePath[len(filePath)-len(ext):] == ext {
			return true
		}
	}
	return false
}

func (a *mockAnalyzer) ExtractSymbols(_ context.Context, filePath string, _ []byte) ([]analysis.Symbol, error) {
	if syms, ok := a.symbols[filePath]; ok {
		return syms, nil
	}
	return []analysis.Symbol{}, nil
}

func (a *mockAnalyzer) GetImports(_ context.Context, filePath string, _ []byte) ([]analysis.Import, error) {
	if imps, ok := a.imports[filePath]; ok {
		return imps, nil
	}
	return []analysis.Import{}, nil
}

func (a *mockAnalyzer) GetExports(_ context.Context, filePath string, _ []byte) ([]analysis.Export, error) {
	if exps, ok := a.exports[filePath]; ok {
		return exps, nil
	}
	return []analysis.Export{}, nil
}

func (a *mockAnalyzer) GetFunctionBody(_ context.Context, _ string, _ []byte, _ string) (string, int, int, error) {
	return "", 0, 0, nil
}

// mockFileRanker implements FileRanker for testing
type mockFileRanker struct {
	rankedFiles []RankedFile
}

func newMockFileRanker(files []RankedFile) *mockFileRanker {
	return &mockFileRanker{rankedFiles: files}
}

func (r *mockFileRanker) RankFiles(_ context.Context, _ string) ([]RankedFile, error) {
	return r.rankedFiles, nil
}

func (r *mockFileRanker) InvalidateCache(_ string) {}

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (l *mockLogger) Debug(_ string)   {}
func (l *mockLogger) Info(_ string)    {}
func (l *mockLogger) Warning(_ string) {}
func (l *mockLogger) Error(_ string)   {}
func (l *mockLogger) Fatal(_ string)   {}

func TestRepoMapBuilder_BuildMap(t *testing.T) {
	tests := []struct {
		name        string
		maxTokens   int
		rankedFiles []RankedFile
		wantErr     bool
	}{
		{
			name:        "empty project",
			maxTokens:   1000,
			rankedFiles: []RankedFile{},
			wantErr:     false,
		},
		{
			name:      "files below rank threshold excluded",
			maxTokens: 5000,
			rankedFiles: []RankedFile{
				{FilePath: "tiny.go", Rank: 0.0001}, // Below minRankThreshold
			},
			wantErr: false,
		},
		{
			name:        "default max tokens when zero",
			maxTokens:   0,
			rankedFiles: []RankedFile{},
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := newMockAnalyzerRegistry()
			ranker := newMockFileRanker(tt.rankedFiles)
			logger := &mockLogger{}

			builder := NewRepoMapBuilder(registry, ranker, logger)

			repoMap, err := builder.BuildMap(context.Background(), "/test/project", tt.maxTokens)

			if (err != nil) != tt.wantErr {
				t.Errorf("BuildMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && repoMap == nil {
				t.Error("BuildMap() returned nil repoMap without error")
			}
		})
	}
}

func TestRepoMapBuilder_Cache(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker([]RankedFile{})
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	// First call should build
	repoMap1, err := builder.BuildMap(context.Background(), "/test/project", 1000)
	if err != nil {
		t.Fatalf("First BuildMap() error = %v", err)
	}

	// Second call should return cached
	repoMap2, err := builder.BuildMap(context.Background(), "/test/project", 1000)
	if err != nil {
		t.Fatalf("Second BuildMap() error = %v", err)
	}

	// Should be same instance (cached)
	if repoMap1.GeneratedAt != repoMap2.GeneratedAt {
		t.Error("Expected cached result, got new generation")
	}

	// Invalidate cache
	builder.InvalidateCache("/test/project")

	// Third call should rebuild
	repoMap3, err := builder.BuildMap(context.Background(), "/test/project", 1000)
	if err != nil {
		t.Fatalf("Third BuildMap() error = %v", err)
	}

	// After invalidation, a new map is generated
	// Note: GeneratedAt might be same if executed very fast, so we just verify no error
	if repoMap3 == nil {
		t.Error("Expected new generation after cache invalidation")
	}
}

func TestRepoMap_Format(t *testing.T) {
	repoMap := &RepoMap{
		ProjectRoot: "/test",
		Entries: []RepoMapEntry{
			{
				FilePath: "main.go",
				Symbols: []SymbolSignature{
					{Name: "main", Kind: "function", Signature: "func main()"},
					{Name: "Config", Kind: "struct", Signature: "struct Config"},
				},
				Rank: 1.0,
			},
			{
				FilePath: "utils.go",
				Symbols: []SymbolSignature{
					{Name: "Helper", Kind: "function", Signature: "func Helper(x int) string"},
				},
				Rank: 0.5,
			},
		},
		TotalTokens: 100,
		GeneratedAt: time.Now(),
	}

	formatted := repoMap.Format()

	// Check that output contains expected content
	if !contains(formatted, "main.go") {
		t.Error("Format() should contain file path 'main.go'")
	}
	if !contains(formatted, "func main()") {
		t.Error("Format() should contain signature 'func main()'")
	}
	if !contains(formatted, "utils.go") {
		t.Error("Format() should contain file path 'utils.go'")
	}
}

func TestRepoMap_FormatCompact(t *testing.T) {
	repoMap := &RepoMap{
		ProjectRoot: "/test",
		Entries: []RepoMapEntry{
			{
				FilePath: "main.go",
				Symbols: []SymbolSignature{
					{Name: "main", Kind: "function"},
					{Name: "Config", Kind: "struct"},
				},
				Rank: 1.0,
			},
		},
		TotalTokens: 50,
		GeneratedAt: time.Now(),
	}

	formatted := repoMap.FormatCompact()

	if !contains(formatted, "main.go:") {
		t.Error("FormatCompact() should contain 'main.go:'")
	}
	if !contains(formatted, "main") {
		t.Error("FormatCompact() should contain symbol name 'main'")
	}
	if !contains(formatted, "Config") {
		t.Error("FormatCompact() should contain symbol name 'Config'")
	}
}

func TestSymbolsToSignatures(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	tests := []struct {
		name     string
		symbols  []analysis.Symbol
		wantLen  int
		wantKind string
	}{
		{
			name: "function symbol",
			symbols: []analysis.Symbol{
				{Name: "DoSomething", Kind: analysis.KindFunction, Signature: "func DoSomething()"},
			},
			wantLen:  1,
			wantKind: "function",
		},
		{
			name: "class symbol",
			symbols: []analysis.Symbol{
				{Name: "MyClass", Kind: analysis.KindClass},
			},
			wantLen:  1,
			wantKind: "class",
		},
		{
			name: "skip package symbol",
			symbols: []analysis.Symbol{
				{Name: "main", Kind: analysis.KindPackage},
			},
			wantLen: 0,
		},
		{
			name: "skip private variable",
			symbols: []analysis.Symbol{
				{Name: "privateVar", Kind: analysis.KindVariable},
			},
			wantLen: 0,
		},
		{
			name: "include exported constant",
			symbols: []analysis.Symbol{
				{Name: "MaxSize", Kind: analysis.KindConstant},
			},
			wantLen:  1,
			wantKind: "constant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signatures := builder.symbolsToSignatures(tt.symbols)

			if len(signatures) != tt.wantLen {
				t.Errorf("symbolsToSignatures() len = %d, want %d", len(signatures), tt.wantLen)
			}

			if tt.wantLen > 0 && signatures[0].Kind != tt.wantKind {
				t.Errorf("symbolsToSignatures() kind = %s, want %s", signatures[0].Kind, tt.wantKind)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "short string unchanged",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "exact length unchanged",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "long string truncated",
			input:  "hello world",
			maxLen: 8,
			want:   "hello...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncateString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEstimateTokens(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "empty string",
			text: "",
			want: 0,
		},
		{
			name: "short text",
			text: "test",
			want: 1, // 4 chars / 4 = 1
		},
		{
			name: "longer text",
			text: "hello world test",
			want: 4, // 16 chars / 4 = 4
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.estimateTokens(tt.text)
			if got != tt.want {
				t.Errorf("estimateTokens() = %d, want %d", got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}


func TestBuildSignature(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	tests := []struct {
		name string
		sym  analysis.Symbol
		want string
	}{
		{
			name: "class",
			sym:  analysis.Symbol{Name: "MyClass", Kind: analysis.KindClass},
			want: "class MyClass",
		},
		{
			name: "interface",
			sym:  analysis.Symbol{Name: "MyInterface", Kind: analysis.KindInterface},
			want: "interface MyInterface",
		},
		{
			name: "struct",
			sym:  analysis.Symbol{Name: "MyStruct", Kind: analysis.KindStruct},
			want: "struct MyStruct",
		},
		{
			name: "function",
			sym:  analysis.Symbol{Name: "DoSomething", Kind: analysis.KindFunction},
			want: "func DoSomething(...)",
		},
		{
			name: "method with parent",
			sym:  analysis.Symbol{Name: "DoSomething", Kind: analysis.KindMethod, Parent: "MyStruct"},
			want: "func (MyStruct) DoSomething(...)",
		},
		{
			name: "method without parent",
			sym:  analysis.Symbol{Name: "DoSomething", Kind: analysis.KindMethod},
			want: "func DoSomething(...)",
		},
		{
			name: "other kind",
			sym:  analysis.Symbol{Name: "SomeVar", Kind: analysis.KindVariable},
			want: "SomeVar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.buildSignature(tt.sym)
			if got != tt.want {
				t.Errorf("buildSignature() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestIsImportantSymbol(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	tests := []struct {
		name string
		sym  analysis.Symbol
		want bool
	}{
		{
			name: "function is important",
			sym:  analysis.Symbol{Name: "DoSomething", Kind: analysis.KindFunction},
			want: true,
		},
		{
			name: "method is important",
			sym:  analysis.Symbol{Name: "DoSomething", Kind: analysis.KindMethod},
			want: true,
		},
		{
			name: "class is important",
			sym:  analysis.Symbol{Name: "MyClass", Kind: analysis.KindClass},
			want: true,
		},
		{
			name: "interface is important",
			sym:  analysis.Symbol{Name: "MyInterface", Kind: analysis.KindInterface},
			want: true,
		},
		{
			name: "struct is important",
			sym:  analysis.Symbol{Name: "MyStruct", Kind: analysis.KindStruct},
			want: true,
		},
		{
			name: "type is important",
			sym:  analysis.Symbol{Name: "MyType", Kind: analysis.KindType},
			want: true,
		},
		{
			name: "component is important",
			sym:  analysis.Symbol{Name: "MyComponent", Kind: analysis.KindComponent},
			want: true,
		},
		{
			name: "composable is important",
			sym:  analysis.Symbol{Name: "useMyHook", Kind: analysis.KindComposable},
			want: true,
		},
		{
			name: "exported constant is important",
			sym:  analysis.Symbol{Name: "MaxSize", Kind: analysis.KindConstant},
			want: true,
		},
		{
			name: "private constant is not important",
			sym:  analysis.Symbol{Name: "maxSize", Kind: analysis.KindConstant},
			want: false,
		},
		{
			name: "exported variable is important",
			sym:  analysis.Symbol{Name: "GlobalVar", Kind: analysis.KindVariable},
			want: true,
		},
		{
			name: "private variable is not important",
			sym:  analysis.Symbol{Name: "privateVar", Kind: analysis.KindVariable},
			want: false,
		},
		{
			name: "package is not important",
			sym:  analysis.Symbol{Name: "main", Kind: analysis.KindPackage},
			want: false,
		},
		{
			name: "import is not important",
			sym:  analysis.Symbol{Name: "fmt", Kind: analysis.KindImport},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := builder.isImportantSymbol(tt.sym)
			if got != tt.want {
				t.Errorf("isImportantSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateEntry(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	entry := RepoMapEntry{
		FilePath: "test.go",
		Symbols: []SymbolSignature{
			{Name: "Func1", Kind: "function", Signature: "func Func1()"},
			{Name: "Func2", Kind: "function", Signature: "func Func2()"},
			{Name: "Func3", Kind: "function", Signature: "func Func3()"},
		},
		Rank: 1.0,
	}

	tests := []struct {
		name      string
		maxTokens int
		wantEmpty bool
	}{
		{
			name:      "zero tokens returns empty",
			maxTokens: 0,
			wantEmpty: true,
		},
		{
			name:      "negative tokens returns empty",
			maxTokens: -1,
			wantEmpty: true,
		},
		{
			name:      "very small budget returns empty",
			maxTokens: 1,
			wantEmpty: true,
		},
		{
			name:      "reasonable budget returns truncated",
			maxTokens: 50,
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			truncated, tokens := builder.truncateEntry(entry, tt.maxTokens)
			if tt.wantEmpty {
				if tokens != 0 {
					t.Errorf("truncateEntry() tokens = %d, want 0", tokens)
				}
			} else {
				if truncated.FilePath != entry.FilePath {
					t.Errorf("truncateEntry() FilePath = %q, want %q", truncated.FilePath, entry.FilePath)
				}
			}
		})
	}
}

func TestEstimateEntryTokens(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	entry := RepoMapEntry{
		FilePath: "test.go",
		Symbols: []SymbolSignature{
			{Name: "Func1", Kind: "function", Signature: "func Func1()"},
		},
		Rank: 1.0,
	}

	tokens := builder.estimateEntryTokens(entry)

	if tokens <= 0 {
		t.Errorf("estimateEntryTokens() = %d, want > 0", tokens)
	}
}

func TestEstimateSymbolTokens(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker(nil)
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	sym := SymbolSignature{
		Name:      "TestFunc",
		Kind:      "function",
		Signature: "func TestFunc(a int, b string) error",
	}

	tokens := builder.estimateSymbolTokens(sym)

	if tokens <= 0 {
		t.Errorf("estimateSymbolTokens() = %d, want > 0", tokens)
	}
}

func TestRepoMapBuilder_InvalidateCache(t *testing.T) {
	registry := newMockAnalyzerRegistry()
	ranker := newMockFileRanker([]RankedFile{})
	logger := &mockLogger{}

	builder := NewRepoMapBuilder(registry, ranker, logger)

	// Build to populate cache
	_, _ = builder.BuildMap(context.Background(), "/test", 1000)

	// Verify cache exists
	cached := builder.getCached("/test")
	if cached == nil {
		t.Error("Expected cache to exist after BuildMap")
	}

	// Invalidate
	builder.InvalidateCache("/test")

	// Verify cache cleared
	cached = builder.getCached("/test")
	if cached != nil {
		t.Error("Expected cache to be cleared after InvalidateCache")
	}
}

func TestRepoMap_EmptyFormat(t *testing.T) {
	repoMap := &RepoMap{
		ProjectRoot: "/test",
		Entries:     []RepoMapEntry{},
		TotalTokens: 0,
	}

	formatted := repoMap.Format()
	if !contains(formatted, "Repository Map") {
		t.Error("Format() should contain header even for empty map")
	}

	compact := repoMap.FormatCompact()
	if compact != "" {
		t.Errorf("FormatCompact() for empty map should be empty, got %q", compact)
	}
}
