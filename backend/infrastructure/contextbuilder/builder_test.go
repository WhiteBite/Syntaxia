package contextbuilder

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestParseContext(t *testing.T) {
	ctx := `--- File: main.go ---
package main

func main() {}

--- File: utils.go ---
package main

func helper() {}
`

	entries := parseContext(ctx)
	if len(entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(entries))
	}

	// Should be sorted by path
	if entries[0].Path != "main.go" {
		t.Errorf("expected first entry 'main.go', got %q", entries[0].Path)
	}
	if entries[1].Path != "utils.go" {
		t.Errorf("expected second entry 'utils.go', got %q", entries[1].Path)
	}
}

func TestParseContext_Empty(t *testing.T) {
	entries := parseContext("")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for empty input, got %d", len(entries))
	}
}

func TestParseContext_NoHeaders(t *testing.T) {
	entries := parseContext("just some text without headers")
	if len(entries) != 0 {
		t.Errorf("expected 0 entries for input without headers, got %d", len(entries))
	}
}

func TestStripComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name:     "line comment //",
			input:    "code // comment\nmore code",
			contains: "code",
			excludes: "comment",
		},
		{
			name:     "line comment #",
			input:    "code # comment\nmore code",
			contains: "code",
			excludes: "comment",
		},
		{
			name:     "block comment",
			input:    "code /* block comment */ more",
			contains: "code",
			excludes: "block comment",
		},
		{
			name:     "multiline block",
			input:    "code /* multi\nline\ncomment */ end",
			contains: "code",
			excludes: "multi",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripComments(tt.input)
			if !strings.Contains(result, tt.contains) {
				t.Errorf("result should contain %q, got %q", tt.contains, result)
			}
		})
	}
}

func TestBuildTree(t *testing.T) {
	paths := []string{
		"src/main.go",
		"src/utils/helper.go",
		"README.md",
	}

	tree := buildTree(paths)

	if tree == "" {
		t.Error("expected non-empty tree")
	}
	if !strings.Contains(tree, "src") {
		t.Error("tree should contain 'src'")
	}
	if !strings.Contains(tree, "main.go") {
		t.Error("tree should contain 'main.go'")
	}
	if !strings.Contains(tree, "README.md") {
		t.Error("tree should contain 'README.md'")
	}
}

func TestBuildTree_Empty(t *testing.T) {
	tree := buildTree([]string{})
	if tree != "" {
		t.Errorf("expected empty tree for empty paths, got %q", tree)
	}
}

func TestBuildTree_Deterministic(t *testing.T) {
	paths := []string{"b.go", "a.go", "c.go"}

	tree1 := buildTree(paths)
	tree2 := buildTree(paths)

	if tree1 != tree2 {
		t.Error("buildTree should be deterministic")
	}
}

func TestBuildFromContext_Plain(t *testing.T) {
	ctx := `--- File: main.go ---
package main

func main() {}
`

	result, err := BuildFromContext("plain", ctx, BuildOptions{})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	if !strings.Contains(result, "--- File: main.go ---") {
		t.Error("result should contain file header")
	}
	if !strings.Contains(result, "package main") {
		t.Error("result should contain file content")
	}
}

func TestBuildFromContext_PlainWithStripComments(t *testing.T) {
	ctx := `--- File: main.go ---
package main
// this is a comment
func main() {}
`

	result, err := BuildFromContext("plain", ctx, BuildOptions{StripComments: true})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	if strings.Contains(result, "this is a comment") {
		t.Error("comments should be stripped")
	}
}

func TestBuildFromContext_Manifest(t *testing.T) {
	ctx := `--- File: src/main.go ---
package main

--- File: src/utils.go ---
package main
`

	result, err := BuildFromContext("manifest", ctx, BuildOptions{})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	if !strings.Contains(result, "Manifest:") {
		t.Error("result should contain 'Manifest:'")
	}
	if !strings.Contains(result, "main.go") {
		t.Error("result should contain file names")
	}
}

func TestBuildFromContext_JSON(t *testing.T) {
	ctx := `--- File: main.go ---
package main

func main() {}
`

	result, err := BuildFromContext("json", ctx, BuildOptions{})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	var entries []entry
	if err := json.Unmarshal([]byte(result), &entries); err != nil {
		t.Fatalf("result should be valid JSON: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Path != "main.go" {
		t.Errorf("expected path 'main.go', got %q", entries[0].Path)
	}
}

func TestBuildFromContext_DefaultFormat(t *testing.T) {
	ctx := `--- File: main.go ---
package main
`

	result, err := BuildFromContext("unknown_format", ctx, BuildOptions{})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	// Should default to manifest format
	if !strings.Contains(result, "Manifest:") {
		t.Error("unknown format should default to manifest")
	}
}

func TestBuildFromContext_EmptyContext(t *testing.T) {
	result, err := BuildFromContext("plain", "", BuildOptions{})
	if err != nil {
		t.Fatalf("BuildFromContext failed: %v", err)
	}

	if result != "" {
		t.Errorf("expected empty result for empty context, got %q", result)
	}
}

func TestBuildTree_WindowsPaths(t *testing.T) {
	paths := []string{
		"src\\main.go",
		"src\\utils\\helper.go",
	}

	tree := buildTree(paths)

	if !strings.Contains(tree, "src") {
		t.Error("tree should handle Windows paths")
	}
	if !strings.Contains(tree, "main.go") {
		t.Error("tree should contain main.go")
	}
}

func TestParseContext_Sorted(t *testing.T) {
	ctx := `--- File: z.go ---
content z

--- File: a.go ---
content a

--- File: m.go ---
content m
`

	entries := parseContext(ctx)
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	// Should be sorted alphabetically
	if entries[0].Path != "a.go" {
		t.Errorf("first entry should be 'a.go', got %q", entries[0].Path)
	}
	if entries[1].Path != "m.go" {
		t.Errorf("second entry should be 'm.go', got %q", entries[1].Path)
	}
	if entries[2].Path != "z.go" {
		t.Errorf("third entry should be 'z.go', got %q", entries[2].Path)
	}
}
