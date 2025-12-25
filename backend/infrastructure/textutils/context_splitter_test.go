package textutils

import (
	"strings"
	"testing"
)

// mockLogger for testing
type mockLogger struct{}

func (l *mockLogger) Debug(msg string)                                     {}
func (l *mockLogger) Info(msg string)                                      {}
func (l *mockLogger) Warning(msg string)                                   {}
func (l *mockLogger) Error(msg string)                                     {}
func (l *mockLogger) Fatal(msg string)                                     {}
func (l *mockLogger) Debugf(format string, args ...interface{})            {}
func (l *mockLogger) Infof(format string, args ...interface{})             {}
func (l *mockLogger) Warningf(format string, args ...interface{})          {}
func (l *mockLogger) Errorf(format string, args ...interface{})            {}
func (l *mockLogger) Fatalf(format string, args ...interface{})            {}
func (l *mockLogger) WithField(key string, value interface{}) interface{}  { return l }
func (l *mockLogger) WithFields(fields map[string]interface{}) interface{} { return l }

func TestApproxTokens(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"test", 1},
		{"hello world", 2}, // 11 chars / 4 = 2
		{"a", 0},           // 1 char / 4 = 0
		{"abcd", 1},        // 4 chars / 4 = 1
	}

	for _, tt := range tests {
		got := approxTokens(tt.input)
		if got != tt.expected {
			t.Errorf("approxTokens(%q) = %d, want %d", tt.input, got, tt.expected)
		}
	}
}

func TestApproxTokens_Unicode(t *testing.T) {
	// Unicode characters should be counted as runes, not bytes
	input := "привет мир" // Russian "hello world"
	tokens := approxTokens(input)
	if tokens < 2 {
		t.Errorf("approxTokens should handle unicode, got %d tokens", tokens)
	}
}

func TestNewContextSplitter(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	if splitter == nil {
		t.Fatal("NewContextSplitter returned nil")
	}
}

func TestSplitContext_FitsInSingleChunk(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	text := "short text"

	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 1000,
		OverlapTokens:     0,
		SplitStrategy:     "token",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0] != text {
		t.Errorf("chunk content mismatch")
	}
}

func TestSplitContext_InvalidMaxTokens(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})

	_, err := splitter.SplitContext("text", struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 0,
		OverlapTokens:     0,
		SplitStrategy:     "token",
	})

	if err == nil {
		t.Error("expected error for MaxTokensPerChunk = 0")
	}

	_, err = splitter.SplitContext("text", struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: -1,
		OverlapTokens:     0,
		SplitStrategy:     "token",
	})

	if err == nil {
		t.Error("expected error for MaxTokensPerChunk < 0")
	}
}

func TestSplitContext_UnknownStrategy(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	longText := strings.Repeat("word ", 1000)

	_, err := splitter.SplitContext(longText, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 100,
		OverlapTokens:     0,
		SplitStrategy:     "unknown_strategy",
	})

	if err == nil {
		t.Error("expected error for unknown strategy")
	}
}

func TestSplitContext_TokenStrategy(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	// Create text that's definitely larger than limit
	longText := strings.Repeat("word ", 100) // ~500 chars = ~125 tokens

	chunks, err := splitter.SplitContext(longText, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 50,
		OverlapTokens:     5, // Small overlap to avoid issues
		SplitStrategy:     "token",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) < 2 {
		t.Errorf("expected multiple chunks, got %d", len(chunks))
	}
}

func TestSplitContext_FileStrategy(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	text := `--- File: file1.go ---
package main

func main() {}

--- File: file2.go ---
package utils

func helper() {}
`
	// Make limit large enough to fit both files
	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 1000,
		OverlapTokens:     0,
		SplitStrategy:     "file",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk (both files fit), got %d", len(chunks))
	}
}

func TestSplitContext_FileStrategy_SplitFiles(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	// Create files that need to be split
	file1 := "--- File: file1.go ---\n" + strings.Repeat("code ", 100)
	file2 := "--- File: file2.go ---\n" + strings.Repeat("more ", 100)
	text := file1 + "\n" + file2

	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 50, // Small limit to force splitting
		OverlapTokens:     0,
		SplitStrategy:     "file",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}
}

func TestSplitContext_SmartStrategy(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	text := `--- File: file1.go ---
package main

--- File: file2.go ---
package utils
`
	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 1000,
		OverlapTokens:     0,
		SplitStrategy:     "smart",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) == 0 {
		t.Error("expected at least 1 chunk")
	}
}

func TestSplitContext_SmartStrategy_FallbackToToken(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	// Single large file that exceeds limit
	text := "--- File: large.go ---\n" + strings.Repeat("code ", 100)

	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 30,
		OverlapTokens:     3,
		SplitStrategy:     "smart",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	// Should fall back to token-based splitting
	if len(chunks) < 2 {
		t.Errorf("expected multiple chunks from fallback, got %d", len(chunks))
	}
}

func TestSplitContext_NoFileHeaders(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	text := strings.Repeat("plain text without headers ", 100)

	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 1000,
		OverlapTokens:     0,
		SplitStrategy:     "file",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	// Should return as single chunk when no headers found
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk for text without headers, got %d", len(chunks))
	}
}

func TestSplitByTokenCount_OverlapGreaterThanLimit(t *testing.T) {
	splitter := &ContextSplitterImpl{log: &mockLogger{}}
	text := strings.Repeat("word ", 100)

	_, err := splitter.splitByTokenCount(text, 10, 20) // overlap > limit
	if err == nil {
		t.Error("expected error when overlap > limit")
	}
}

func TestSplitByTokenCount_OverlapEqualsLimit(t *testing.T) {
	splitter := &ContextSplitterImpl{log: &mockLogger{}}
	text := strings.Repeat("word ", 100)

	_, err := splitter.splitByTokenCount(text, 10, 10) // overlap == limit
	if err == nil {
		t.Error("expected error when overlap == limit")
	}
}

func TestSplitContext_WithManifest(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})
	text := `Manifest:
├── src
│   └── main.go

--- File: src/main.go ---
package main
`
	chunks, err := splitter.SplitContext(text, struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 1000,
		OverlapTokens:     0,
		SplitStrategy:     "file",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	// Manifest should be included
	if !strings.Contains(chunks[0], "Manifest:") {
		t.Error("manifest should be preserved")
	}
}

func TestSplitContext_EmptyText(t *testing.T) {
	splitter := NewContextSplitter(&mockLogger{})

	chunks, err := splitter.SplitContext("", struct {
		MaxTokensPerChunk int
		OverlapTokens     int
		SplitStrategy     string
	}{
		MaxTokensPerChunk: 100,
		OverlapTokens:     0,
		SplitStrategy:     "token",
	})

	if err != nil {
		t.Fatalf("SplitContext failed: %v", err)
	}
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk for empty text, got %d", len(chunks))
	}
}
