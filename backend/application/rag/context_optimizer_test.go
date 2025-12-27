package rag

import (
	"context"
	"strings"
	"testing"

	"syntaxia/domain/analysis"
)

// optimizerMockRegistry implements analysis.AnalyzerRegistry for testing
type optimizerMockRegistry struct {
	analyzer analysis.LanguageAnalyzer
}

func (m *optimizerMockRegistry) Register(_ analysis.LanguageAnalyzer)                     {}
func (m *optimizerMockRegistry) GetAnalyzer(_ string) analysis.LanguageAnalyzer           { return m.analyzer }
func (m *optimizerMockRegistry) GetAnalyzerByLanguage(_ string) analysis.LanguageAnalyzer { return m.analyzer }
func (m *optimizerMockRegistry) SupportedLanguages() []string                             { return []string{"go"} }
func (m *optimizerMockRegistry) SupportedExtensions() []string                            { return []string{".go"} }

// optimizerMockAnalyzer implements analysis.LanguageAnalyzer for testing
type optimizerMockAnalyzer struct {
	symbols []analysis.Symbol
	err     error
}

func (m *optimizerMockAnalyzer) Language() string         { return "go" }
func (m *optimizerMockAnalyzer) Extensions() []string     { return []string{".go"} }
func (m *optimizerMockAnalyzer) CanAnalyze(_ string) bool { return true }

func (m *optimizerMockAnalyzer) ExtractSymbols(_ context.Context, _ string, _ []byte) ([]analysis.Symbol, error) {
	return m.symbols, m.err
}

func (m *optimizerMockAnalyzer) GetImports(_ context.Context, _ string, _ []byte) ([]analysis.Import, error) {
	return nil, nil
}

func (m *optimizerMockAnalyzer) GetExports(_ context.Context, _ string, _ []byte) ([]analysis.Export, error) {
	return nil, nil
}

func (m *optimizerMockAnalyzer) GetFunctionBody(_ context.Context, _ string, _ []byte, _ string) (string, int, int, error) {
	return "", 0, 0, nil
}

// optimizerMockLogger implements domain.Logger for testing
type optimizerMockLogger struct{}

func (m *optimizerMockLogger) Debug(_ string)   {}
func (m *optimizerMockLogger) Info(_ string)    {}
func (m *optimizerMockLogger) Warning(_ string) {}
func (m *optimizerMockLogger) Error(_ string)   {}
func (m *optimizerMockLogger) Fatal(_ string)   {}

func TestContextOptimizer_Optimize(t *testing.T) {
	tests := []struct {
		name           string
		files          []FileContent
		opts           OptimizeOptions
		wantFileCount  int
		wantFirstPath  string
		checkContent   func([]FileContent) bool
	}{
		{
			name: "basic optimization",
			files: []FileContent{
				{Path: "a.go", Content: "package main\n\n\n\nfunc main() {}", Tokens: 10, Priority: 1},
				{Path: "b.go", Content: "package util\n\nfunc helper() {}", Tokens: 8, Priority: 2},
			},
			opts:          OptimizeOptions{RemoveEmptyLines: true, PositionalSort: true},
			wantFileCount: 2,
			wantFirstPath: "b.go",
			checkContent: func(files []FileContent) bool {
				return !strings.Contains(files[0].Content, "\n\n\n")
			},
		},
		{
			name: "token budget trimming",
			files: []FileContent{
				{Path: "a.go", Content: "package main", Tokens: 100, Priority: 1},
				{Path: "b.go", Content: "package util", Tokens: 100, Priority: 2},
			},
			opts:          OptimizeOptions{MaxTokens: 150, PositionalSort: true},
			wantFileCount: 2, // First file fits, second gets truncated
			wantFirstPath: "b.go",
		},
		{
			name: "remove comments",
			files: []FileContent{
				{Path: "test.go", Content: "// comment\npackage main\n/* block */\nfunc main() {}", Tokens: 20},
			},
			opts:          OptimizeOptions{RemoveComments: true},
			wantFileCount: 1,
			checkContent: func(files []FileContent) bool {
				return !strings.Contains(files[0].Content, "// comment") &&
					!strings.Contains(files[0].Content, "/* block */")
			},
		},
		{
			name:          "empty input",
			files:         []FileContent{},
			opts:          OptimizeOptions{},
			wantFileCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})
			result, err := optimizer.Optimize(context.Background(), tt.files, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != tt.wantFileCount {
				t.Errorf("got %d files, want %d", len(result), tt.wantFileCount)
			}

			if tt.wantFirstPath != "" && len(result) > 0 && result[0].Path != tt.wantFirstPath {
				t.Errorf("first file = %s, want %s", result[0].Path, tt.wantFirstPath)
			}

			if tt.checkContent != nil && !tt.checkContent(result) {
				t.Error("content check failed")
			}
		})
	}
}

func TestContextOptimizer_OptimizeContent(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		filePath    string
		opts        OptimizeOptions
		wantContain string
		wantExclude string
	}{
		{
			name:        "remove single line comments",
			content:     "// comment\ncode line\n// another",
			filePath:    "test.go",
			opts:        OptimizeOptions{RemoveComments: true},
			wantContain: "code line",
			wantExclude: "// comment",
		},
		{
			name:        "remove multi-line comments",
			content:     "/* block\ncomment */\ncode",
			filePath:    "test.go",
			opts:        OptimizeOptions{RemoveComments: true},
			wantContain: "code",
			wantExclude: "block",
		},
		{
			name:        "collapse empty lines",
			content:     "line1\n\n\n\nline2",
			filePath:    "test.go",
			opts:        OptimizeOptions{RemoveEmptyLines: true},
			wantContain: "line1\n\nline2",
			wantExclude: "\n\n\n",
		},
		{
			name:        "python hash comments",
			content:     "# comment\ncode = 1\n# another",
			filePath:    "test.py",
			opts:        OptimizeOptions{RemoveComments: true},
			wantContain: "code = 1",
			wantExclude: "# comment",
		},
		{
			name:        "preserve code without options",
			content:     "// comment\nfunc main() {}",
			filePath:    "test.go",
			opts:        OptimizeOptions{},
			wantContain: "// comment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})
			result := optimizer.OptimizeContent(context.Background(), tt.content, tt.filePath, tt.opts)

			if tt.wantContain != "" && !strings.Contains(result, tt.wantContain) {
				t.Errorf("result should contain %q, got: %s", tt.wantContain, result)
			}

			if tt.wantExclude != "" && strings.Contains(result, tt.wantExclude) {
				t.Errorf("result should not contain %q, got: %s", tt.wantExclude, result)
			}
		})
	}
}

func TestContextOptimizer_CalculateSNR(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantMin float64
		wantMax float64
	}{
		{
			name:    "high signal code",
			content: "func main() {\n\treturn nil\n}\nfunc helper() {\n\treturn err\n}",
			wantMin: 0.3,
			wantMax: 1.0,
		},
		{
			name:    "mostly comments",
			content: "// comment 1\n// comment 2\n// comment 3\ncode",
			wantMin: 0.1,
			wantMax: 0.5,
		},
		{
			name:    "empty content",
			content: "",
			wantMin: 0.0,
			wantMax: 0.0,
		},
		{
			name:    "mixed content",
			content: "func test() {\n// comment\n\treturn nil\n}",
			wantMin: 0.1,
			wantMax: 0.8,
		},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snr := optimizer.CalculateSNR(tt.content)

			if snr < tt.wantMin || snr > tt.wantMax {
				t.Errorf("SNR = %f, want between %f and %f", snr, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestContextOptimizer_CollapseFunctions(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		symbols     []analysis.Symbol
		wantContain string
	}{
		{
			name: "collapse large function",
			content: `func LargeFunc() {
	line1
	line2
	line3
	line4
	line5
	line6
}`,
			symbols: []analysis.Symbol{
				{Name: "LargeFunc", Kind: analysis.KindFunction, StartLine: 1, EndLine: 8},
			},
			wantContain: collapsedBodyToken,
		},
		{
			name: "keep small function",
			content: `func Small() {
	return
}`,
			symbols: []analysis.Symbol{
				{Name: "Small", Kind: analysis.KindFunction, StartLine: 1, EndLine: 3},
			},
			wantContain: "return",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &optimizerMockRegistry{
				analyzer: &optimizerMockAnalyzer{symbols: tt.symbols},
			}

			optimizer := NewContextOptimizer(registry, &optimizerMockLogger{})
			result := optimizer.OptimizeContent(context.Background(), tt.content, "test.go", OptimizeOptions{
				CollapseFunctions: true,
			})

			if !strings.Contains(result, tt.wantContain) {
				t.Errorf("result should contain %q, got: %s", tt.wantContain, result)
			}
		})
	}
}

func TestContextOptimizer_SortByPosition(t *testing.T) {
	tests := []struct {
		name      string
		files     []FileContent
		wantOrder []string
	}{
		{
			name: "sort by priority",
			files: []FileContent{
				{Path: "low.go", Priority: 1, Relevance: 0.5},
				{Path: "high.go", Priority: 3, Relevance: 0.5},
				{Path: "mid.go", Priority: 2, Relevance: 0.5},
			},
			wantOrder: []string{"high.go", "mid.go", "low.go"},
		},
		{
			name: "sort by relevance when same priority",
			files: []FileContent{
				{Path: "a.go", Priority: 1, Relevance: 0.3},
				{Path: "b.go", Priority: 1, Relevance: 0.9},
				{Path: "c.go", Priority: 1, Relevance: 0.6},
			},
			wantOrder: []string{"b.go", "c.go", "a.go"},
		},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := optimizer.sortByPosition(tt.files)

			for i, wantPath := range tt.wantOrder {
				if result[i].Path != wantPath {
					t.Errorf("position %d: got %s, want %s", i, result[i].Path, wantPath)
				}
			}
		})
	}
}

func TestContextOptimizer_TrimToTokenBudget(t *testing.T) {
	tests := []struct {
		name          string
		files         []FileContent
		maxTokens     int
		wantFileCount int
		wantPaths     []string
	}{
		{
			name: "all files fit",
			files: []FileContent{
				{Path: "a.go", Tokens: 100},
				{Path: "b.go", Tokens: 100},
			},
			maxTokens:     300,
			wantFileCount: 2,
			wantPaths:     []string{"a.go", "b.go"},
		},
		{
			name: "trim to budget",
			files: []FileContent{
				{Path: "a.go", Tokens: 100},
				{Path: "b.go", Tokens: 100},
				{Path: "c.go", Tokens: 100},
			},
			maxTokens:     250,
			wantFileCount: 2,
			wantPaths:     []string{"a.go", "b.go"},
		},
		{
			name: "truncate last file",
			files: []FileContent{
				{Path: "a.go", Content: strings.Repeat("x", 400), Tokens: 100},
				{Path: "b.go", Content: strings.Repeat("y", 800), Tokens: 200},
			},
			maxTokens:     250,
			wantFileCount: 2,
		},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := optimizer.trimToTokenBudget(tt.files, tt.maxTokens)

			if len(result) != tt.wantFileCount {
				t.Errorf("got %d files, want %d", len(result), tt.wantFileCount)
			}

			if tt.wantPaths != nil {
				for i, wantPath := range tt.wantPaths {
					if i < len(result) && result[i].Path != wantPath {
						t.Errorf("file %d: got %s, want %s", i, result[i].Path, wantPath)
					}
				}
			}
		})
	}
}

func TestContextOptimizer_CalculateOptimalStrategy(t *testing.T) {
	tests := []struct {
		name                  string
		files                 []FileContent
		targetTokens          int
		wantRemoveComments    bool
		wantRemoveEmptyLines  bool
		wantCollapseFunctions bool
	}{
		{
			name: "no optimization needed",
			files: []FileContent{
				{Tokens: 100},
			},
			targetTokens:          200,
			wantRemoveComments:    false,
			wantRemoveEmptyLines:  false,
			wantCollapseFunctions: false,
		},
		{
			name: "light optimization",
			files: []FileContent{
				{Tokens: 120},
			},
			targetTokens:          100,
			wantRemoveComments:    false,
			wantRemoveEmptyLines:  true,
			wantCollapseFunctions: false,
		},
		{
			name: "aggressive optimization",
			files: []FileContent{
				{Tokens: 300},
			},
			targetTokens:          100,
			wantRemoveComments:    true,
			wantRemoveEmptyLines:  true,
			wantCollapseFunctions: true,
		},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := optimizer.CalculateOptimalStrategy(tt.files, tt.targetTokens)

			if opts.RemoveComments != tt.wantRemoveComments {
				t.Errorf("RemoveComments = %v, want %v", opts.RemoveComments, tt.wantRemoveComments)
			}

			if opts.RemoveEmptyLines != tt.wantRemoveEmptyLines {
				t.Errorf("RemoveEmptyLines = %v, want %v", opts.RemoveEmptyLines, tt.wantRemoveEmptyLines)
			}

			if opts.CollapseFunctions != tt.wantCollapseFunctions {
				t.Errorf("CollapseFunctions = %v, want %v", opts.CollapseFunctions, tt.wantCollapseFunctions)
			}
		})
	}
}

func TestContextOptimizer_IsSignalLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"func main() {", true},
		{"return nil", true},
		{"if err != nil {", true},
		{"x := 5", true},
		{"// comment", false},
		{"", false},
		{"class MyClass:", true},
		{"def method(self):", true},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := optimizer.isSignalLine(tt.line); got != tt.want {
				t.Errorf("isSignalLine(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestContextOptimizer_IsNoiseLine(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"// comment", true},
		{"# python comment", true},
		{"/* block start", true},
		{"* continuation", true},
		{"", true},
		{"func main() {", false},
		{"{", false},
		{"}", false},
	}

	optimizer := NewContextOptimizer(&optimizerMockRegistry{}, &optimizerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			if got := optimizer.isNoiseLine(tt.line); got != tt.want {
				t.Errorf("isNoiseLine(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"file.go", "go"},
		{"path/to/file.ts", "ts"},
		{"noext", ""},
		{".hidden", "hidden"},
		{"file.test.js", "js"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := getFileExtension(tt.path); got != tt.want {
				t.Errorf("getFileExtension(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
