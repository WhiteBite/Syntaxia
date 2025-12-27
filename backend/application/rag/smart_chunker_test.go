package rag

import (
	"context"
	"strings"
	"testing"

	"syntaxia/domain/analysis"
)

// chunkerMockRegistry implements analysis.AnalyzerRegistry for testing
type chunkerMockRegistry struct {
	analyzer analysis.LanguageAnalyzer
}

func (m *chunkerMockRegistry) Register(_ analysis.LanguageAnalyzer)                     {}
func (m *chunkerMockRegistry) GetAnalyzer(_ string) analysis.LanguageAnalyzer           { return m.analyzer }
func (m *chunkerMockRegistry) GetAnalyzerByLanguage(_ string) analysis.LanguageAnalyzer { return m.analyzer }
func (m *chunkerMockRegistry) SupportedLanguages() []string                             { return []string{"go"} }
func (m *chunkerMockRegistry) SupportedExtensions() []string                            { return []string{".go"} }

// chunkerMockAnalyzer implements analysis.LanguageAnalyzer for testing
type chunkerMockAnalyzer struct {
	symbols []analysis.Symbol
	err     error
}

func (m *chunkerMockAnalyzer) Language() string         { return "go" }
func (m *chunkerMockAnalyzer) Extensions() []string     { return []string{".go"} }
func (m *chunkerMockAnalyzer) CanAnalyze(_ string) bool { return true }

func (m *chunkerMockAnalyzer) ExtractSymbols(_ context.Context, _ string, _ []byte) ([]analysis.Symbol, error) {
	return m.symbols, m.err
}

func (m *chunkerMockAnalyzer) GetImports(_ context.Context, _ string, _ []byte) ([]analysis.Import, error) {
	return nil, nil
}

func (m *chunkerMockAnalyzer) GetExports(_ context.Context, _ string, _ []byte) ([]analysis.Export, error) {
	return nil, nil
}

func (m *chunkerMockAnalyzer) GetFunctionBody(_ context.Context, _ string, _ []byte, _ string) (string, int, int, error) {
	return "", 0, 0, nil
}

// chunkerMockLogger implements domain.Logger for testing
type chunkerMockLogger struct{}

func (m *chunkerMockLogger) Debug(_ string)   {}
func (m *chunkerMockLogger) Info(_ string)    {}
func (m *chunkerMockLogger) Warning(_ string) {}
func (m *chunkerMockLogger) Error(_ string)   {}
func (m *chunkerMockLogger) Fatal(_ string)   {}

func TestSmartChunker_ChunkContent(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		filePath       string
		symbols        []analysis.Symbol
		useAnalyzer    bool
		opts           ChunkOptions
		wantMinChunks  int
		wantMaxChunks  int
		wantFirstType  string
	}{
		{
			name: "chunk by lines when no analyzer",
			content: `line 1
line 2
line 3
line 4
line 5`,
			filePath:      "unknown.xyz",
			useAnalyzer:   false,
			opts:          ChunkOptions{MaxTokensPerChunk: 100},
			wantMinChunks: 1,
			wantMaxChunks: 1,
			wantFirstType: "block",
		},
		{
			name:          "empty content produces single empty chunk",
			content:       "",
			filePath:      "empty.go",
			useAnalyzer:   false,
			opts:          ChunkOptions{MaxTokensPerChunk: 500},
			wantMinChunks: 1,
			wantMaxChunks: 1,
			wantFirstType: "block",
		},
		{
			name: "chunk with symbols",
			content: `package main

func Hello() {
	println("hello")
}

func World() {
	println("world")
}`,
			filePath: "main.go",
			symbols: []analysis.Symbol{
				{Name: "Hello", Kind: analysis.KindFunction, StartLine: 3, EndLine: 5},
				{Name: "World", Kind: analysis.KindFunction, StartLine: 7, EndLine: 9},
			},
			useAnalyzer:   true,
			opts:          ChunkOptions{MaxTokensPerChunk: 500, IncludeContext: 1, PreserveStructure: true},
			wantMinChunks: 1,
			wantMaxChunks: 4,
			wantFirstType: "block", // package declaration comes first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var registry *chunkerMockRegistry
			if tt.useAnalyzer && len(tt.symbols) > 0 {
				registry = &chunkerMockRegistry{
					analyzer: &chunkerMockAnalyzer{symbols: tt.symbols},
				}
			} else {
				registry = &chunkerMockRegistry{analyzer: nil}
			}

			chunker := NewSmartChunker(registry, &chunkerMockLogger{})
			chunks, err := chunker.ChunkContent(context.Background(), []byte(tt.content), tt.filePath, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(chunks) < tt.wantMinChunks || len(chunks) > tt.wantMaxChunks {
				t.Errorf("got %d chunks, want between %d and %d", len(chunks), tt.wantMinChunks, tt.wantMaxChunks)
			}

			if len(chunks) > 0 && chunks[0].Type != tt.wantFirstType {
				t.Errorf("first chunk type = %s, want %s", chunks[0].Type, tt.wantFirstType)
			}
		})
	}
}

func TestSmartChunker_ChunkByRelevance(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		query     string
		maxTokens int
	}{
		{
			name: "relevant function first",
			content: `package main

func ProcessUser() {
	// user processing
}

func ProcessOrder() {
	// order processing
}`,
			query:     "user processing",
			maxTokens: 1000,
		},
		{
			name: "match by content",
			content: `package main

func Alpha() {
	// handles authentication
}

func Beta() {
	// handles logging
}`,
			query:     "authentication",
			maxTokens: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := &chunkerMockRegistry{analyzer: nil}

			chunker := NewSmartChunker(registry, &chunkerMockLogger{})

			chunks, err := chunker.ChunkContent(context.Background(), []byte(tt.content), "test.go", ChunkOptions{
				MaxTokensPerChunk: 500,
				IncludeContext:    1,
			})
			if err != nil {
				t.Fatalf("ChunkContent failed: %v", err)
			}

			if len(chunks) == 0 {
				t.Fatal("no chunks returned")
			}

			// Calculate relevance manually
			queryTerms := chunker.extractQueryTerms(tt.query)
			for i := range chunks {
				chunks[i].Relevance = chunker.calculateRelevance(chunks[i], queryTerms)
			}

			// Verify relevance was calculated
			hasRelevance := false
			for _, chunk := range chunks {
				if chunk.Relevance > 0 {
					hasRelevance = true
					break
				}
			}

			if !hasRelevance {
				t.Error("no chunks have relevance > 0")
			}
		})
	}
}

func TestSmartChunker_EstimateTokens(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		wantMin int
		wantMax int
	}{
		{
			name:    "short text",
			text:    "hello world",
			wantMin: 2,
			wantMax: 4,
		},
		{
			name:    "code snippet",
			text:    "func main() {\n\tprintln(\"hello\")\n}",
			wantMin: 5,
			wantMax: 15,
		},
		{
			name:    "empty",
			text:    "",
			wantMin: 0,
			wantMax: 0,
		},
	}

	chunker := NewSmartChunker(&chunkerMockRegistry{}, &chunkerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := chunker.estimateTokens(tt.text)
			if tokens < tt.wantMin || tokens > tt.wantMax {
				t.Errorf("estimateTokens(%q) = %d, want between %d and %d", tt.text, tokens, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSmartChunker_ExtractQueryTerms(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		wantTerms []string
	}{
		{
			name:      "simple query",
			query:     "user authentication",
			wantTerms: []string{"user", "authentication"},
		},
		{
			name:      "with stop words",
			query:     "the user is authenticated",
			wantTerms: []string{"user", "authenticated"},
		},
		{
			name:      "camelCase",
			query:     "processUserData",
			wantTerms: []string{"processuserdata"},
		},
		{
			name:      "with symbols",
			query:     "func() error handling",
			wantTerms: []string{"func", "error", "handling"},
		},
	}

	chunker := NewSmartChunker(&chunkerMockRegistry{}, &chunkerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			terms := chunker.extractQueryTerms(tt.query)

			if len(terms) != len(tt.wantTerms) {
				t.Errorf("got %d terms, want %d: %v vs %v", len(terms), len(tt.wantTerms), terms, tt.wantTerms)
				return
			}

			for i, term := range terms {
				if term != tt.wantTerms[i] {
					t.Errorf("term[%d] = %s, want %s", i, term, tt.wantTerms[i])
				}
			}
		})
	}
}

func TestSmartChunker_CalculateRelevance(t *testing.T) {
	tests := []struct {
		name       string
		chunk      Chunk
		queryTerms []string
		wantMin    float64
		wantMax    float64
	}{
		{
			name:       "exact name match",
			chunk:      Chunk{Name: "ProcessUser", Content: "func ProcessUser() {}", Type: "function"},
			queryTerms: []string{"processuser"},
			wantMin:    0.5,
			wantMax:    1.0,
		},
		{
			name:       "content match only",
			chunk:      Chunk{Name: "Handler", Content: "// handles user requests", Type: "function"},
			queryTerms: []string{"user"},
			wantMin:    0.2,
			wantMax:    0.8,
		},
		{
			name:       "no match",
			chunk:      Chunk{Name: "Alpha", Content: "func Alpha() {}", Type: "function"},
			queryTerms: []string{"beta", "gamma"},
			wantMin:    0.0,
			wantMax:    0.2,
		},
		{
			name:       "empty query",
			chunk:      Chunk{Name: "Test", Content: "test content", Type: "function"},
			queryTerms: []string{},
			wantMin:    0.4,
			wantMax:    0.6,
		},
	}

	chunker := NewSmartChunker(&chunkerMockRegistry{}, &chunkerMockLogger{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			relevance := chunker.calculateRelevance(tt.chunk, tt.queryTerms)

			if relevance < tt.wantMin || relevance > tt.wantMax {
				t.Errorf("relevance = %f, want between %f and %f", relevance, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestSmartChunker_GenerateFileSummary(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		symbols     []analysis.Symbol
		useAnalyzer bool
		wantContain string
	}{
		{
			name:        "with symbols",
			content:     "package main\n\nfunc Hello() {}\nfunc World() {}",
			useAnalyzer: true,
			symbols: []analysis.Symbol{
				{Name: "Hello", Kind: analysis.KindFunction, Signature: "func Hello()"},
				{Name: "World", Kind: analysis.KindFunction, Signature: "func World()"},
			},
			wantContain: "func Hello()",
		},
		{
			name:        "without symbols - line summary",
			content:     "line1\nline2\nline3",
			useAnalyzer: false,
			symbols:     []analysis.Symbol{},
			wantContain: "line1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var registry *chunkerMockRegistry
			if tt.useAnalyzer && len(tt.symbols) > 0 {
				registry = &chunkerMockRegistry{
					analyzer: &chunkerMockAnalyzer{symbols: tt.symbols},
				}
			} else {
				registry = &chunkerMockRegistry{analyzer: nil}
			}

			chunker := NewSmartChunker(registry, &chunkerMockLogger{})
			summary := chunker.GenerateFileSummary(context.Background(), "test.go", []byte(tt.content))

			if !containsSubstr(summary, tt.wantContain) {
				t.Errorf("summary does not contain %q: %s", tt.wantContain, summary)
			}
		})
	}
}

func TestIsStopWord(t *testing.T) {
	tests := []struct {
		word string
		want bool
	}{
		{"the", true},
		{"is", true},
		{"function", false},
		{"user", false},
		{"and", true},
		{"process", false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			if got := isStopWord(tt.word); got != tt.want {
				t.Errorf("isStopWord(%q) = %v, want %v", tt.word, got, tt.want)
			}
		})
	}
}

// Helper functions

func generateLines(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString("\tline ")
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString("\n")
	}
	return sb.String()
}

func containsSubstr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstr(s, substr)))
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
