package context

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddLineNumbers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line",
			input:    "hello world",
			expected: "1 | hello world",
		},
		{
			name:     "multiple lines",
			input:    "line1\nline2\nline3",
			expected: "1 | line1\n2 | line2\n3 | line3",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "1 | ",
		},
		{
			name:     "ten lines with padding",
			input:    "1\n2\n3\n4\n5\n6\n7\n8\n9\n10",
			expected: " 1 | 1\n 2 | 2\n 3 | 3\n 4 | 4\n 5 | 5\n 6 | 6\n 7 | 7\n 8 | 8\n 9 | 9\n10 | 10",
		},
		{
			name:     "preserves indentation",
			input:    "func main() {\n    fmt.Println(\"hello\")\n}",
			expected: "1 | func main() {\n2 |     fmt.Println(\"hello\")\n3 | }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addLineNumbers(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestEscapeForFormat(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		format   OutputFormat
		expected string
	}{
		{
			name:     "XML escapes ampersand",
			content:  "a & b",
			format:   FormatXML,
			expected: "a &amp; b",
		},
		{
			name:     "XML escapes less than",
			content:  "a < b",
			format:   FormatXML,
			expected: "a &lt; b",
		},
		{
			name:     "XML escapes greater than",
			content:  "a > b",
			format:   FormatXML,
			expected: "a &gt; b",
		},
		{
			name:     "XML escapes all special chars",
			content:  "<div>&test</div>",
			format:   FormatXML,
			expected: "&lt;div&gt;&amp;test&lt;/div&gt;",
		},
		{
			name:     "Markdown no escape",
			content:  "<div>&test</div>",
			format:   FormatMarkdown,
			expected: "<div>&test</div>",
		},
		{
			name:     "Plain no escape",
			content:  "<div>&test</div>",
			format:   FormatPlain,
			expected: "<div>&test</div>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := escapeForFormat(tt.content, tt.format)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatFileHeader(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		format   OutputFormat
		contains string
	}{
		{
			name:     "XML format",
			filePath: "src/main.go",
			format:   FormatXML,
			contains: "<file path=\"src/main.go\">",
		},
		{
			name:     "Markdown format",
			filePath: "src/main.go",
			format:   FormatMarkdown,
			contains: "## File: src/main.go",
		},
		{
			name:     "Plain format",
			filePath: "src/main.go",
			format:   FormatPlain,
			contains: "--- File: src/main.go ---",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileHeader(tt.filePath, tt.format)
			assert.True(t, strings.Contains(result, tt.contains), "Expected %q to contain %q", result, tt.contains)
		})
	}
}

func TestFormatFileFooter(t *testing.T) {
	tests := []struct {
		name     string
		format   OutputFormat
		contains string
	}{
		{
			name:     "XML format",
			format:   FormatXML,
			contains: "</file>",
		},
		{
			name:     "Markdown format",
			format:   FormatMarkdown,
			contains: "```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFileFooter(tt.format)
			assert.True(t, strings.Contains(result, tt.contains), "Expected %q to contain %q", result, tt.contains)
		})
	}
}
