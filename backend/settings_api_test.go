package main

import (
	"testing"
)

func TestMatchesIgnorePattern(t *testing.T) {
	app := &App{}

	tests := []struct {
		name     string
		path     string
		rules    string
		expected bool
	}{
		// Empty rules
		{"empty rules", "file.go", "", false},
		
		// Exact filename match
		{"exact filename", "go.mod", "go.mod", true},
		{"exact filename in subdir", "subdir/go.mod", "go.mod", true},
		{"exact filename no match", "go.sum", "go.mod", false},
		
		// Directory patterns (with trailing slash)
		{"dir pattern root", ".git/config", ".git/", true},
		{"dir pattern nested", "src/.git/config", ".git/", true},
		{"dir pattern file inside", ".monitoring/logs/app.log", ".monitoring/", true},
		{"dir pattern no match", "git/file.txt", ".git/", false},
		
		// Wildcard patterns
		{"wildcard extension", "app.log", "*.log", true},
		{"wildcard extension nested", "logs/app.log", "*.log", true},
		{"wildcard no match", "app.txt", "*.log", false},
		
		// Path prefix
		{"path prefix", "vendor/pkg/file.go", "vendor", true},
		{"path prefix nested", "node_modules/lodash/index.js", "node_modules", true},
		
		// Comments should be ignored
		{"comment line", "file.go", "# comment\nfile.go", true},
		{"only comment", "file.go", "# comment", false},
		
		// Multiple rules
		{"multiple rules first match", "go.mod", "go.mod\ngo.sum", true},
		{"multiple rules second match", "go.sum", "go.mod\ngo.sum", true},
		{"multiple rules no match", "main.go", "go.mod\ngo.sum", false},
		
		// Windows path normalization
		{"windows path", "src\\file.go", "src/", true},
		
		// Deploy directory
		{"deploy dir", ".deploy/config.yaml", ".deploy/", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.matchesIgnorePattern(tt.path, tt.rules)
			if result != tt.expected {
				t.Errorf("matchesIgnorePattern(%q, %q) = %v, want %v", 
					tt.path, tt.rules, result, tt.expected)
			}
		})
	}
}

func TestMatchPattern(t *testing.T) {
	app := &App{}

	tests := []struct {
		name     string
		path     string
		pattern  string
		isDir    bool
		expected bool
	}{
		// Exact filename
		{"exact match", "go.mod", "go.mod", false, true},
		{"exact in subdir", "sub/go.mod", "go.mod", false, true},
		
		// Directory patterns
		{"dir at root", ".git/config", ".git", true, true},
		{"dir nested", "a/.git/b", ".git", true, true},
		{"dir no match", "git/file", ".git", true, false},
		
		// Wildcards
		{"wildcard star", "test.log", "*.log", false, true},
		
		// Path prefix
		{"prefix match", "vendor/lib/x.go", "vendor", false, true},
		{"prefix no match", "myvendor/x.go", "vendor", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := app.matchPattern(tt.path, tt.pattern, tt.isDir)
			if result != tt.expected {
				t.Errorf("matchPattern(%q, %q, %v) = %v, want %v",
					tt.path, tt.pattern, tt.isDir, result, tt.expected)
			}
		})
	}
}
