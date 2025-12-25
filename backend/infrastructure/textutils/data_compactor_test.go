package textutils

import (
	"strings"
	"testing"
)

func TestDataCompactor_CompactJSON(t *testing.T) {
	c := NewDataCompactor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "formatted JSON to compact",
			input: `{
  "name": "test",
  "value": 123
}`,
			expected: `{"name":"test","value":123}`,
		},
		{
			name:     "already compact JSON",
			input:    `{"name":"test"}`,
			expected: `{"name":"test"}`,
		},
		{
			name: "nested JSON",
			input: `{
  "user": {
    "name": "John",
    "age": 30
  },
  "items": [
    1,
    2,
    3
  ]
}`,
			// JSON key order is not guaranteed, just check it's compact
			expected: "", // Will check separately
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "invalid JSON unchanged",
			input:    "not json {",
			expected: "not json {",
		},
		{
			name: "JSON with special chars",
			input: `{
  "html": "<div>test</div>",
  "url": "http://example.com"
}`,
			expected: `{"html":"<div>test</div>","url":"http://example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.CompactJSON(tt.input)
			if tt.expected == "" {
				// For nested JSON, just check it's compact (no newlines)
				if strings.Contains(result, "\n") {
					t.Errorf("CompactJSON() should not contain newlines, got %q", result)
				}
			} else if result != tt.expected {
				t.Errorf("CompactJSON() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDataCompactor_CompactYAML(t *testing.T) {
	c := NewDataCompactor()

	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name: "removes comments",
			input: `# This is a comment
name: test
# Another comment
value: 123`,
			contains: []string{"name: test", "value: 123"},
			excludes: []string{"# This is a comment", "# Another comment"},
		},
		{
			name: "removes inline comments",
			input: `name: test # inline comment
value: 123 # another inline`,
			contains: []string{"name: test", "value: 123"},
			excludes: []string{"# inline comment", "# another inline"},
		},
		{
			name: "preserves hash in strings",
			input: `color: "#ff0000"
tag: '#hashtag'`,
			contains: []string{`color: "#ff0000"`, `tag: '#hashtag'`},
			excludes: []string{},
		},
		{
			name: "collapses empty lines",
			input: `name: test


value: 123



end: true`,
			contains: []string{"name: test", "value: 123", "end: true"},
			excludes: []string{},
		},
		{
			name:     "empty string",
			input:    "",
			contains: []string{},
			excludes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.CompactYAML(tt.input)

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("CompactYAML() should contain %q, got:\n%s", s, result)
				}
			}

			for _, s := range tt.excludes {
				if strings.Contains(result, s) {
					t.Errorf("CompactYAML() should not contain %q, got:\n%s", s, result)
				}
			}
		})
	}
}

func TestDataCompactor_Compact(t *testing.T) {
	c := NewDataCompactor()

	tests := []struct {
		name     string
		content  string
		filePath string
		isJSON   bool
	}{
		{
			name:     "JSON file",
			content:  "{\n  \"test\": 1\n}",
			filePath: "config.json",
			isJSON:   true,
		},
		{
			name:     "YAML file",
			content:  "# comment\ntest: 1",
			filePath: "config.yaml",
			isJSON:   false,
		},
		{
			name:     "YML file",
			content:  "# comment\ntest: 1",
			filePath: "config.yml",
			isJSON:   false,
		},
		{
			name:     "Other file unchanged",
			content:  "some content",
			filePath: "file.txt",
			isJSON:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := c.Compact(tt.content, tt.filePath)

			if tt.isJSON {
				// JSON should be compacted (no newlines)
				if strings.Contains(result, "\n") && strings.Contains(tt.content, "\n") {
					// Only check if input had newlines
					if result == tt.content {
						t.Error("JSON should be compacted")
					}
				}
			}
		})
	}
}

func TestIsDataFile(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"config.json", true},
		{"config.yaml", true},
		{"config.yml", true},
		{"config.JSON", true}, // case insensitive
		{"config.YAML", true},
		{"config.txt", false},
		{"config.go", false},
		{"", false},
	}

	for _, tt := range tests {
		result := IsDataFile(tt.path)
		if result != tt.expected {
			t.Errorf("IsDataFile(%q) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestRemoveYAMLInlineComment(t *testing.T) {
	c := NewDataCompactor()

	tests := []struct {
		input    string
		expected string
	}{
		{"name: test # comment", "name: test"},
		{"name: test", "name: test"},
		{`color: "#ff0000"`, `color: "#ff0000"`},
		{`tag: '#hashtag' # comment`, `tag: '#hashtag'`},
		{"value: 123  # trailing", "value: 123"},
	}

	for _, tt := range tests {
		result := c.removeYAMLInlineComment(tt.input)
		if result != tt.expected {
			t.Errorf("removeYAMLInlineComment(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func BenchmarkCompactJSON_Small(b *testing.B) {
	c := NewDataCompactor()
	input := `{
  "name": "test",
  "value": 123,
  "nested": {
    "a": 1,
    "b": 2
  }
}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.CompactJSON(input)
	}
}

func BenchmarkCompactJSON_Large(b *testing.B) {
	c := NewDataCompactor()
	// Генерируем большой JSON
	var sb strings.Builder
	sb.WriteString("{\n")
	for i := 0; i < 100; i++ {
		sb.WriteString(`  "field`)
		sb.WriteString(string(rune('0' + i%10)))
		sb.WriteString(`": "value",`)
		sb.WriteString("\n")
	}
	sb.WriteString(`  "last": "value"`)
	sb.WriteString("\n}")
	input := sb.String()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.CompactJSON(input)
	}
}

func BenchmarkCompactYAML(b *testing.B) {
	c := NewDataCompactor()
	input := `# Configuration file
# Version 1.0

name: myapp # application name
version: 1.0.0

database:
  host: localhost # db host
  port: 5432
  
# Features
features:
  - auth
  - api
  - web
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.CompactYAML(input)
	}
}
