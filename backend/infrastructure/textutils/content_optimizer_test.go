package textutils

import (
	"context"
	"syntaxia/domain/analysis"
	"strings"
	"testing"
)

// mockCommentStripper implements CommentStripperInterface for testing
type mockCommentStripper struct{}

func (m *mockCommentStripper) Strip(content string, filePath string) string {
	// Simple mock: remove lines starting with //
	lines := strings.Split(content, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, "//") {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

func TestContentOptimizer_Optimize_CollapseEmptyLines(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	input := "line1\n\n\n\n\nline2"
	result := opt.Optimize(ctx, input, "file.go", OptimizeOptions{
		CollapseEmptyLines: true,
	})

	if strings.Contains(result, "\n\n\n\n") {
		t.Error("should collapse multiple empty lines")
	}
}

func TestContentOptimizer_Optimize_StripLicense(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	input := `/*
 * Copyright 2024 Example Corp
 * MIT License
 */

package main`

	result := opt.Optimize(ctx, input, "main.go", OptimizeOptions{
		StripLicense: true,
	})

	if strings.Contains(result, "Copyright") {
		t.Error("should strip license header")
	}
	if !strings.Contains(result, "package main") {
		t.Error("should preserve code")
	}
}

func TestContentOptimizer_Optimize_StripComments(t *testing.T) {
	opt := NewContentOptimizerSimple(&mockCommentStripper{})
	ctx := context.Background()

	input := `package main
// This is a comment
func main() {}`

	result := opt.Optimize(ctx, input, "main.go", OptimizeOptions{
		StripComments: true,
	})

	if strings.Contains(result, "This is a comment") {
		t.Error("should strip comments")
	}
	if !strings.Contains(result, "func main()") {
		t.Error("should preserve code")
	}
}

func TestContentOptimizer_Optimize_CompactJSON(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	input := `{
  "name": "test",
  "value": 123
}`

	result := opt.Optimize(ctx, input, "config.json", OptimizeOptions{
		CompactDataFiles: true,
	})

	if strings.Contains(result, "\n") {
		t.Error("should compact JSON to single line")
	}
	if !strings.Contains(result, `"name":"test"`) {
		t.Error("should preserve JSON content")
	}
}

func TestContentOptimizer_Optimize_TrimWhitespace(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	input := "line1   \nline2  \t\nline3"
	result := opt.Optimize(ctx, input, "file.go", OptimizeOptions{
		TrimWhitespace: true,
	})

	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.HasSuffix(line, " ") || strings.HasSuffix(line, "\t") {
			t.Error("should trim trailing whitespace")
		}
	}
}

func TestContentOptimizer_Optimize_Combined(t *testing.T) {
	opt := NewContentOptimizerSimple(&mockCommentStripper{})
	ctx := context.Background()

	input := `/*
 * Copyright 2024 Example
 * MIT License
 */

package main

// Comment to remove

func main() {
    // Another comment
}
`

	result := opt.Optimize(ctx, input, "main.go", OptimizeOptions{
		CollapseEmptyLines: true,
		StripLicense:       true,
		StripComments:      true,
		TrimWhitespace:     true,
	})

	if strings.Contains(result, "Copyright") {
		t.Error("should strip license")
	}
	if strings.Contains(result, "Comment to remove") {
		t.Error("should strip comments")
	}
	// Check that result is shorter than input (optimizations applied)
	if len(result) >= len(input) {
		t.Error("optimized result should be shorter than input")
	}
}

func TestContentOptimizer_Optimize_EmptyContent(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	result := opt.Optimize(ctx, "", "file.go", AggressiveOptimizeOptions())

	if result != "" {
		t.Error("should return empty string for empty input")
	}
}

func TestContentOptimizer_OptimizeWithDefaults(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	input := "line1\n\n\n\nline2   "
	result := opt.OptimizeWithDefaults(ctx, input, "file.go")

	// Default options include CollapseEmptyLines and TrimWhitespace
	if strings.Contains(result, "\n\n\n\n") {
		t.Error("default should collapse empty lines")
	}
}

func TestContentOptimizer_OptimizeAggressive(t *testing.T) {
	opt := NewContentOptimizerSimple(&mockCommentStripper{})
	ctx := context.Background()

	input := `/*
 * Copyright 2024
 */
// comment
code`

	result := opt.OptimizeAggressive(ctx, input, "file.go")

	if strings.Contains(result, "Copyright") {
		t.Error("aggressive should strip license")
	}
}

func TestContentOptimizer_SkeletonMode(t *testing.T) {
	// Create mock analyzer and registry
	analyzer := &mockAnalyzer{
		lang: "go",
		exts: []string{".go"},
		symbols: []analysis.Symbol{
			{Name: "main", Kind: analysis.KindPackage},
			{Name: "User", Kind: analysis.KindStruct},
		},
	}
	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".go": analyzer},
	}

	opt := NewContentOptimizer(registry, nil)
	ctx := context.Background()

	input := `package main

type User struct {
    ID   int
    Name string
    // lots of fields...
}

func (u *User) GetID() int {
    return u.ID
}
`

	result := opt.Optimize(ctx, input, "main.go", OptimizeOptions{
		SkeletonMode: true,
	})

	// Skeleton should be much shorter
	if len(result) >= len(input) {
		t.Error("skeleton should be shorter than original")
	}
	if !strings.Contains(result, "User") {
		t.Error("skeleton should contain type name")
	}
}

func TestContentOptimizer_CanGenerateSkeleton(t *testing.T) {
	// Without registry
	opt := NewContentOptimizerSimple(nil)
	if opt.CanGenerateSkeleton("main.go") {
		t.Error("should return false without registry")
	}

	// With registry
	analyzer := &mockAnalyzer{lang: "go", exts: []string{".go"}}
	registry := &mockRegistry{
		analyzers: map[string]analysis.LanguageAnalyzer{".go": analyzer},
	}
	opt2 := NewContentOptimizer(registry, nil)
	if !opt2.CanGenerateSkeleton("main.go") {
		t.Error("should return true with registry")
	}
}

func TestDefaultOptimizeOptions(t *testing.T) {
	opts := DefaultOptimizeOptions()

	if !opts.CollapseEmptyLines {
		t.Error("default should enable CollapseEmptyLines")
	}
	if !opts.TrimWhitespace {
		t.Error("default should enable TrimWhitespace")
	}
	if opts.StripComments {
		t.Error("default should not enable StripComments")
	}
}

func TestAggressiveOptimizeOptions(t *testing.T) {
	opts := AggressiveOptimizeOptions()

	if !opts.CollapseEmptyLines {
		t.Error("aggressive should enable CollapseEmptyLines")
	}
	if !opts.StripLicense {
		t.Error("aggressive should enable StripLicense")
	}
	if !opts.StripComments {
		t.Error("aggressive should enable StripComments")
	}
	if !opts.CompactDataFiles {
		t.Error("aggressive should enable CompactDataFiles")
	}
}

func TestEstimateSavings(t *testing.T) {
	tests := []struct {
		name     string
		opts     OptimizeOptions
		minSaved int
	}{
		{"skeleton mode", SkeletonOptimizeOptions(), 50},
		{"aggressive", AggressiveOptimizeOptions(), 20},
		{"default", DefaultOptimizeOptions(), 5},
		{"empty", OptimizeOptions{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			savings := EstimateSavings(tt.opts)
			if savings < tt.minSaved {
				t.Errorf("EstimateSavings() = %d, want >= %d", savings, tt.minSaved)
			}
		})
	}
}

func TestOptimizeOptionsFromFlags(t *testing.T) {
	opts := OptimizeOptionsFromFlags(true, true, false, true, false, true)

	if !opts.CollapseEmptyLines {
		t.Error("should set CollapseEmptyLines")
	}
	if !opts.StripLicense {
		t.Error("should set StripLicense")
	}
	if opts.StripComments {
		t.Error("should not set StripComments")
	}
	if !opts.CompactDataFiles {
		t.Error("should set CompactDataFiles")
	}
	if opts.SkeletonMode {
		t.Error("should not set SkeletonMode")
	}
	if !opts.TrimWhitespace {
		t.Error("should set TrimWhitespace")
	}
}

func TestContentOptimizer_OptimizeBatch(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	files := map[string]string{
		"file1.go": "line1\n\n\n\nline2",
		"file2.go": "line3\n\n\n\nline4",
	}

	result := opt.OptimizeBatch(ctx, files, OptimizeOptions{
		CollapseEmptyLines: true,
	})

	if len(result) != 2 {
		t.Errorf("should return %d files, got %d", 2, len(result))
	}

	for _, content := range result {
		if strings.Contains(content, "\n\n\n\n") {
			t.Error("batch should optimize all files")
		}
	}
}

func TestContentOptimizer_OptimizeBatchWithStats(t *testing.T) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	files := map[string]string{
		"file1.go": "line1\n\n\n\n\n\n\nline2",
		"file2.go": "line3\n\n\n\n\n\n\nline4",
	}

	result, stats := opt.OptimizeBatchWithStats(ctx, files, OptimizeOptions{
		CollapseEmptyLines: true,
	})

	if len(result) != 2 {
		t.Error("should return all files")
	}

	if stats.FilesProcessed != 2 {
		t.Errorf("FilesProcessed = %d, want 2", stats.FilesProcessed)
	}

	if stats.SavedBytes <= 0 {
		t.Error("should save some bytes")
	}

	if stats.SavedPercent <= 0 {
		t.Error("should have positive save percent")
	}
}

func BenchmarkContentOptimizer_Optimize_Default(b *testing.B) {
	opt := NewContentOptimizerSimple(&mockCommentStripper{})
	ctx := context.Background()

	input := `package main

import "fmt"

// Main function
func main() {
    fmt.Println("Hello")
}


`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opt.OptimizeWithDefaults(ctx, input, "main.go")
	}
}

func BenchmarkContentOptimizer_Optimize_Aggressive(b *testing.B) {
	opt := NewContentOptimizerSimple(&mockCommentStripper{})
	ctx := context.Background()

	input := `/*
 * Copyright 2024 Example
 */
package main

import "fmt"

// Main function
func main() {
    // Print hello
    fmt.Println("Hello")
}


`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opt.OptimizeAggressive(ctx, input, "main.go")
	}
}

func BenchmarkContentOptimizer_OptimizeBatch(b *testing.B) {
	opt := NewContentOptimizerSimple(nil)
	ctx := context.Background()

	files := make(map[string]string)
	for i := 0; i < 100; i++ {
		files["file"+string(rune('0'+i%10))+".go"] = "line1\n\n\n\nline2\n\n\n\nline3"
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		opt.OptimizeBatch(ctx, files, DefaultOptimizeOptions())
	}
}
