package context

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_isTestFile tests the isTestFile method
func TestService_isTestFile(t *testing.T) {
	service := &Service{}

	tests := []struct {
		path     string
		expected bool
		desc     string
	}{
		// Go tests
		{"main_test.go", true, "Go test file"},
		{"service_test.go", true, "Go test with underscore"},
		{"main.go", false, "Go source file"},

		// JS/TS tests
		{"app.test.js", true, "JS test file"},
		{"app.spec.ts", true, "TS spec file"},
		{"app.test.tsx", true, "TSX test file"},
		{"app.js", false, "JS source file"},

		// Python tests
		{"test_main.py", true, "Python test prefix"},
		{"main_test.py", true, "Python test suffix"},
		{"main.py", false, "Python source file"},

		// Directory-based (need /dir/ pattern)
		{"src/tests/main.go", true, "File in tests dir"},
		{"src/__tests__/app.js", true, "File in __tests__ dir"},
		{"project/src/__tests__/app.js", true, "Nested __tests__ dir"},
		{"project/e2e/login.spec.ts", true, "File in e2e dir"},

		// Edge cases
		{"src/main.go", false, "Source file"},
		{"lib/utils.ts", false, "Lib file"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			result := service.isTestFile(tt.path)
			assert.Equal(t, tt.expected, result, "isTestFile(%q)", tt.path)
		})
	}
}

// TestService_filterTestFiles tests the filterTestFiles method
func TestService_filterTestFiles(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything).Return()

	service := &Service{logger: mockLogger}

	input := []string{
		"main.go",
		"main_test.go",
		"service.go",
		"service_test.go",
		"utils.go",
		"tests/helper.go",
	}

	result := service.filterTestFiles(input)

	expected := []string{
		"main.go",
		"service.go",
		"utils.go",
	}

	assert.Equal(t, len(expected), len(result))
	for i, path := range result {
		assert.Equal(t, expected[i], path)
	}
}

// TestService_stripLicenseHeader tests license header removal
func TestService_stripLicenseHeader(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name     string
		content  string
		filePath string
		contains string
		excludes string
	}{
		{
			name: "Go MIT license",
			content: `/*
 * Copyright 2024 Example Corp
 * MIT License
 */

package main

func main() {}`,
			filePath: "main.go",
			contains: "package main",
			excludes: "Copyright",
		},
		{
			name: "Python license",
			content: `# Copyright 2024 Example
# Licensed under MIT

def main():
    pass`,
			filePath: "main.py",
			contains: "def main():",
			excludes: "Copyright",
		},
		{
			name: "No license preserved",
			content: `package main

func main() {}`,
			filePath: "main.go",
			contains: "package main",
			excludes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.stripLicenseHeader(tt.content, tt.filePath)
			if tt.contains != "" {
				assert.Contains(t, result, tt.contains)
			}
			if tt.excludes != "" {
				assert.NotContains(t, result, tt.excludes)
			}
		})
	}
}

// TestService_compactDataFile tests JSON/YAML compaction
func TestService_compactDataFile(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name        string
		content     string
		filePath    string
		hasNewlines bool
	}{
		{
			name: "JSON compacted",
			content: `{
  "name": "test",
  "value": 123
}`,
			filePath:    "config.json",
			hasNewlines: false,
		},
		{
			name: "YAML comments removed",
			content: `# Comment
name: test
# Another comment
value: 123`,
			filePath:    "config.yaml",
			hasNewlines: true, // YAML keeps structure
		},
		{
			name:        "Non-data file unchanged",
			content:     "some content\nwith newlines",
			filePath:    "file.txt",
			hasNewlines: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.compactDataFile(tt.content, tt.filePath)
			if tt.hasNewlines {
				// For YAML and non-data files, just check it's not empty
				assert.NotEmpty(t, result)
			} else {
				// JSON should be compacted to single line
				assert.False(t, strings.Contains(result, "\n"), "should not contain newlines")
			}
		})
	}
}

// TestService_trimTrailingWhitespace tests whitespace trimming
func TestService_trimTrailingWhitespace(t *testing.T) {
	service := &Service{}

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
			name:     "no trailing whitespace unchanged",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.trimTrailingWhitespace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestService_collapseEmptyLines tests empty line collapsing
func TestService_collapseEmptyLines(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "multiple empty lines collapsed",
			input:    "line1\n\n\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
		{
			name:     "two empty lines preserved",
			input:    "line1\n\n\nline2",
			expected: "line1\n\n\nline2",
		},
		{
			name:     "no empty lines unchanged",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "fast path - no triple newlines",
			input:    "line1\n\nline2",
			expected: "line1\n\nline2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.collapseEmptyLines(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestService_applyContentOptimizations tests the full optimization pipeline
func TestService_applyContentOptimizations(t *testing.T) {
	service := &Service{}

	input := `/*
 * Copyright 2024 Example
 * MIT License
 */

package main

func main() {
    println("hello")
}


`

	// Test with all optimizations enabled
	options := &BuildOptions{
		StripLicense:       true,
		TrimWhitespace:     true,
		CollapseEmptyLines: true,
	}

	result := service.applyContentOptimizations(input, "main.go", options)

	// Should not contain license
	assert.NotContains(t, result, "Copyright")

	// Should contain code
	assert.Contains(t, result, "package main")
	assert.Contains(t, result, "func main()")

	// Should be shorter than input
	assert.Less(t, len(result), len(input))
}

// TestService_applyContentOptimizations_NilOptions tests nil options handling
func TestService_applyContentOptimizations_NilOptions(t *testing.T) {
	service := &Service{}

	input := "some content"
	result := service.applyContentOptimizations(input, "file.go", nil)

	// Should return unchanged content
	assert.Equal(t, input, result)
}

// TestService_compactJSON tests JSON compaction
func TestService_compactJSON(t *testing.T) {
	service := &Service{}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "formatted to compact",
			input: `{
  "name": "test",
  "value": 123
}`,
			expected: `{"name":"test","value":123}`,
		},
		{
			name:     "already compact",
			input:    `{"name":"test"}`,
			expected: `{"name":"test"}`,
		},
		{
			name:     "invalid JSON unchanged",
			input:    "not json {",
			expected: "not json {",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.compactJSON(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestService_compactYAML tests YAML compaction
func TestService_compactYAML(t *testing.T) {
	service := &Service{}

	input := `# Comment to remove
name: test
# Another comment
value: 123`

	result := service.compactYAML(input)

	// Should not contain comments
	assert.NotContains(t, result, "# Comment")
	assert.NotContains(t, result, "# Another")

	// Should contain data
	assert.Contains(t, result, "name: test")
	assert.Contains(t, result, "value: 123")
}
