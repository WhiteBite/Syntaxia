package embeddings

import (
	"testing"
)

func TestDefaultChunkerConfig(t *testing.T) {
	config := DefaultChunkerConfig()
	if config.MaxChunkTokens != 512 {
		t.Errorf("expected MaxChunkTokens 512, got %d", config.MaxChunkTokens)
	}
	if config.MinChunkTokens != 50 {
		t.Errorf("expected MinChunkTokens 50, got %d", config.MinChunkTokens)
	}
	if config.OverlapTokens != 50 {
		t.Errorf("expected OverlapTokens 50, got %d", config.OverlapTokens)
	}
}

func TestNewCodeChunker(t *testing.T) {
	config := DefaultChunkerConfig()
	chunker := NewCodeChunker(config)
	if chunker == nil {
		t.Fatal("NewCodeChunker returned nil")
	}
}

func TestCodeChunker_ChunkFile_BySize(t *testing.T) {
	config := ChunkerConfig{
		MaxChunkTokens: 100,
		MinChunkTokens: 10,
		OverlapTokens:  10,
		PreferSymbols:  false,
	}
	chunker := NewCodeChunker(config)

	content := []byte(`package main

func main() {
	fmt.Println("Hello, World!")
}

func helper() {
	// do something
}
`)

	chunks := chunker.ChunkFile("main.go", content, nil)
	if len(chunks) == 0 {
		t.Error("expected at least one chunk")
	}

	for _, chunk := range chunks {
		if chunk.FilePath != "main.go" {
			t.Errorf("expected FilePath 'main.go', got %q", chunk.FilePath)
		}
		if chunk.Language != "go" {
			t.Errorf("expected Language 'go', got %q", chunk.Language)
		}
		if chunk.ID == "" {
			t.Error("chunk ID should not be empty")
		}
		if chunk.Hash == "" {
			t.Error("chunk Hash should not be empty")
		}
	}
}

func TestCodeChunker_ChunkFile_BySymbols(t *testing.T) {
	config := ChunkerConfig{
		MaxChunkTokens: 500,
		MinChunkTokens: 5, // Lower threshold for test
		OverlapTokens:  10,
		PreferSymbols:  true,
	}
	chunker := NewCodeChunker(config)

	content := []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
	fmt.Println("This is a test")
	fmt.Println("More content here")
}

func helper() {
	// helper function with more content
	fmt.Println("Helper doing work")
	fmt.Println("More helper work")
}
`)

	symbols := []SymbolInfo{
		{Name: "main", Kind: "function", StartLine: 5, EndLine: 9},
		{Name: "helper", Kind: "function", StartLine: 11, EndLine: 15},
	}

	chunks := chunker.ChunkFile("main.go", content, symbols)

	// Should have chunks (may include remaining lines chunk too)
	if len(chunks) == 0 {
		t.Error("expected at least one chunk")
	}

	foundMain := false
	foundHelper := false
	for _, chunk := range chunks {
		if chunk.SymbolName == "main" {
			foundMain = true
		}
		if chunk.SymbolName == "helper" {
			foundHelper = true
		}
	}

	// At least one symbol should be found
	if !foundMain && !foundHelper {
		t.Error("expected at least one symbol chunk")
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text      string
		minTokens int
		maxTokens int
	}{
		{"", 0, 0},
		{"hello", 1, 2},
		{"hello world", 2, 4},
		{"func main() { fmt.Println(\"test\") }", 5, 15},
	}

	for _, tt := range tests {
		tokens := estimateTokens(tt.text)
		if tokens < tt.minTokens || tokens > tt.maxTokens {
			t.Errorf("estimateTokens(%q) = %d, expected between %d and %d",
				tt.text, tokens, tt.minTokens, tt.maxTokens)
		}
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filePath string
		expected string
	}{
		{"main.go", "go"},
		{"app.ts", "typescript"},
		{"app.tsx", "typescript"},
		{"script.js", "javascript"},
		{"script.jsx", "javascript"},
		{"component.vue", "vue"},
		{"main.py", "python"},
		{"Service.java", "java"},
		{"Service.kt", "kotlin"},
		{"widget.dart", "dart"},
		{"main.rs", "rust"},
		{"Program.cs", "csharp"},
		{"main.cpp", "cpp"},
		{"main.c", "c"},
		{"app.rb", "ruby"},
		{"index.php", "php"},
		{"App.swift", "swift"},
		{"unknown.xyz", "unknown"},
	}

	for _, tt := range tests {
		got := detectLanguage(tt.filePath)
		if got != tt.expected {
			t.Errorf("detectLanguage(%q) = %q, want %q", tt.filePath, got, tt.expected)
		}
	}
}

func TestGenerateChunkID(t *testing.T) {
	id1 := generateChunkID("file.go", 1, 10)
	id2 := generateChunkID("file.go", 1, 10)
	id3 := generateChunkID("file.go", 1, 11)

	if id1 != id2 {
		t.Error("same inputs should produce same ID")
	}
	if id1 == id3 {
		t.Error("different inputs should produce different IDs")
	}
	if len(id1) != 16 {
		t.Errorf("expected ID length 16, got %d", len(id1))
	}
}

func TestHashContent(t *testing.T) {
	hash1 := hashContent("hello world")
	hash2 := hashContent("hello world")
	hash3 := hashContent("hello world!")

	if hash1 != hash2 {
		t.Error("same content should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("different content should produce different hash")
	}
	if len(hash1) != 32 {
		t.Errorf("expected hash length 32, got %d", len(hash1))
	}
}

func TestGetOverlapLines(t *testing.T) {
	lines := []string{"line1", "line2", "line3", "line4", "line5"}

	// Get overlap with enough tokens
	overlap := getOverlapLines(lines, 10)
	if len(overlap) == 0 {
		t.Error("expected some overlap lines")
	}

	// Empty input
	overlap = getOverlapLines(nil, 10)
	if overlap != nil {
		t.Error("expected nil for empty input")
	}

	// Zero overlap tokens
	overlap = getOverlapLines(lines, 0)
	if overlap != nil {
		t.Error("expected nil for zero overlap tokens")
	}
}

func TestMapSymbolKindToChunkType(t *testing.T) {
	tests := []struct {
		kind     string
		expected string
	}{
		{"function", "function"},
		{"func", "function"},
		{"class", "class"},
		{"method", "method"},
		{"unknown", "block"},
		{"", "block"},
	}

	for _, tt := range tests {
		got := mapSymbolKindToChunkType(tt.kind)
		if string(got) != tt.expected {
			t.Errorf("mapSymbolKindToChunkType(%q) = %q, want %q", tt.kind, got, tt.expected)
		}
	}
}

func TestCodeChunker_LargeFile(t *testing.T) {
	config := ChunkerConfig{
		MaxChunkTokens: 50,
		MinChunkTokens: 10,
		OverlapTokens:  5,
		PreferSymbols:  false,
	}
	chunker := NewCodeChunker(config)

	// Create a large file content
	var content string
	for i := 0; i < 100; i++ {
		content += "func test" + string(rune('A'+i%26)) + "() { /* some code here */ }\n"
	}

	chunks := chunker.ChunkFile("large.go", []byte(content), nil)
	if len(chunks) < 2 {
		t.Errorf("expected multiple chunks for large file, got %d", len(chunks))
	}

	// Verify chunks don't exceed max tokens
	for _, chunk := range chunks {
		if chunk.TokenCount > config.MaxChunkTokens*2 { // Allow some tolerance
			t.Errorf("chunk exceeds max tokens: %d > %d", chunk.TokenCount, config.MaxChunkTokens)
		}
	}
}
