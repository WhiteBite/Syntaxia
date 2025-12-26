package staticanalyzer

import (
	"os"
	"testing"
)

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"0", 0},
		{"1", 1},
		{"10", 10},
		{"123", 123},
		{"999999", 999999},
		{"-1", -1},
		{"-100", -100},
		{"", 0},
		{"abc", 0},
		{"12abc", 12},
		{"abc12", 0},
		{"  10", 10},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	// Test with existing file
	tempFile, err := os.CreateTemp("", "test_file_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	defer os.Remove(tempPath)

	if !fileExists(tempPath) {
		t.Errorf("fileExists(%q) = false, want true", tempPath)
	}

	// Test with non-existing file
	nonExistentPath := "/non/existent/path/to/file.txt"
	if fileExists(nonExistentPath) {
		t.Errorf("fileExists(%q) = true, want false", nonExistentPath)
	}

	// Test with empty path
	if fileExists("") {
		t.Error("fileExists(\"\") = true, want false")
	}
}

func TestFileExists_Directory(t *testing.T) {
	// Test with existing directory
	tempDir, err := os.MkdirTemp("", "test_dir_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// fileExists returns true for directories too (it just checks if path exists)
	if !fileExists(tempDir) {
		t.Errorf("fileExists(%q) = false, want true for directory", tempDir)
	}
}
