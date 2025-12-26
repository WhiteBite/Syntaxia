package textutils

import (
	"strings"
	"testing"
)

func TestWhitespaceOptimizer_CollapseEmptyLines(t *testing.T) {
	w := NewWhitespaceOptimizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no empty lines",
			input:    "line1\nline2\nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "single empty line preserved",
			input:    "line1\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "two empty lines preserved",
			input:    "line1\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
		{
			name:     "three empty lines collapsed to two",
			input:    "line1\n\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
		{
			name:     "many empty lines collapsed",
			input:    "line1\n\n\n\n\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
		{
			name:     "multiple groups of empty lines",
			input:    "a\n\n\n\nb\n\n\n\n\nc",
			expected: "a\n\n\nb\n\n\nc",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace only lines not collapsed (performance optimization)",
			input:    "line1\n   \n   \n   \n   \nline2",
			expected: "line1\n   \n   \n   \n   \nline2", // whitespace-only lines are kept as-is
		},
		{
			name:     "truly empty lines collapsed",
			input:    "line1\n\n\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := w.CollapseEmptyLines(tt.input)
			if result != tt.expected {
				t.Errorf("CollapseEmptyLines() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestWhitespaceOptimizer_TrimTrailingWhitespace(t *testing.T) {
	w := NewWhitespaceOptimizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "trailing spaces removed",
			input:    "line1   \nline2  \nline3",
			expected: "line1\nline2\nline3",
		},
		{
			name:     "trailing tabs removed",
			input:    "line1\t\t\nline2\t",
			expected: "line1\nline2",
		},
		{
			name:     "mixed trailing whitespace",
			input:    "line1 \t \nline2",
			expected: "line1\nline2",
		},
		{
			name:     "no trailing whitespace unchanged",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := w.TrimTrailingWhitespace(tt.input)
			if result != tt.expected {
				t.Errorf("TrimTrailingWhitespace() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestIsWhitespaceOnly(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", true},
		{"\t", true},
		{"  \t  ", true},
		{"a", false},
		{" a ", false},
		{"\ta\t", false},
	}

	for _, tt := range tests {
		result := isWhitespaceOnly(tt.input)
		if result != tt.expected {
			t.Errorf("isWhitespaceOnly(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// Benchmark для проверки производительности
func BenchmarkCollapseEmptyLines_Small(b *testing.B) {
	w := NewWhitespaceOptimizer()
	input := "line1\n\n\n\nline2\n\n\n\nline3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.CollapseEmptyLines(input)
	}
}

func BenchmarkCollapseEmptyLines_Large(b *testing.B) {
	w := NewWhitespaceOptimizer()
	// Симулируем большой файл с множеством пустых строк
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("func example() {\n")
		sb.WriteString("    return nil\n")
		sb.WriteString("}\n")
		sb.WriteString("\n\n\n\n\n") // 5 пустых строк
	}
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.CollapseEmptyLines(input)
	}
}

func BenchmarkCollapseEmptyLines_NoChanges(b *testing.B) {
	w := NewWhitespaceOptimizer()
	// Файл без множественных пустых строк - должен быть быстрый путь
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("func example() {\n")
		sb.WriteString("    return nil\n")
		sb.WriteString("}\n\n")
	}
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.CollapseEmptyLines(input)
	}
}
