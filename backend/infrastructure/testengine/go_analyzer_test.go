package testengine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// TestNewGoTestAnalyzer tests analyzer creation
func TestNewGoTestAnalyzer(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.log)
}

// TestGoTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestGoTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-deps-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests := []struct {
		name            string
		content         string
		expectedImports []string
	}{
		{
			name: "multi_line_imports",
			content: `package main

import (
	"testing"
	"os"
	"path/filepath"
)

func TestSomething(t *testing.T) {}
`,
			expectedImports: []string{"testing", "os", "path/filepath"},
		},
		{
			name: "single_import",
			content: `package main

import "testing"

func TestSingle(t *testing.T) {}
`,
			expectedImports: []string{"testing"},
		},
		{
			name: "imports_with_aliases",
			content: `package main

import (
	"testing"
	myalias "some/package"
	. "dot/import"
	_ "blank/import"
)

func TestAliases(t *testing.T) {}
`,
			expectedImports: []string{"testing"},
		},
		{
			name: "imports_with_comments",
			content: `package main

import (
	"testing"
	"os"
)

func TestComments(t *testing.T) {}
`,
			expectedImports: []string{"testing", "os"},
		},
		{
			name: "no_imports",
			content: `package main

func TestNoImports(t *testing.T) {}
`,
			expectedImports: []string{},
		},
		{
			name: "mixed_single_and_multi",
			content: `package main

import "fmt"

import (
	"testing"
	"os"
)

func TestMixed(t *testing.T) {}
`,
			expectedImports: []string{"fmt", "testing", "os"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+"_test.go")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
			require.NoError(t, err)

			for _, expected := range tt.expectedImports {
				assert.Contains(t, deps, expected, "Should contain import: %s", expected)
			}
		})
	}
}

// TestGoTestAnalyzer_AnalyzeTestDependencies_NonExistent tests with non-existent file
func TestGoTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/test_file.go")
	assert.Error(t, err)
}

// TestGoTestAnalyzer_FindTestsForFile tests finding tests for source file
func TestGoTestAnalyzer_FindTestsForFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.go": `package main

func Helper() {}
`,
		"utils_test.go": `package main

import "testing"

func TestHelper(t *testing.T) {}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.go", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
	assert.Contains(t, tests, "utils_test.go")
}

// TestGoTestAnalyzer_FindTestsForFile_NoTestFile tests when no test file exists
func TestGoTestAnalyzer_FindTestsForFile_NoTestFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-no-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.go": `package main

func Helper() {}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.go", tempDir)
	require.NoError(t, err)

	// Should return empty but no error
	assert.Empty(t, tests)
}

// TestGoTestAnalyzer_FindTestsForFile_InSubdirectory tests finding tests in subdirectory
func TestGoTestAnalyzer_FindTestsForFile_InSubdirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-subdir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create subdirectory
	subDir := filepath.Join(tempDir, "pkg")
	err = os.MkdirAll(subDir, 0o755)
	require.NoError(t, err)

	files := map[string]string{
		"pkg/utils.go": `package pkg

func Helper() {}
`,
		"pkg/utils_test.go": `package pkg

import "testing"

func TestHelper(t *testing.T) {}
`,
		"pkg/other_test.go": `package pkg

import "testing"

func TestOther(t *testing.T) {}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "pkg/utils.go", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
}

// TestGoTestAnalyzer_IsSmokeTest tests smoke test detection
func TestGoTestAnalyzer_IsSmokeTest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-smoke-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "smoke_in_comment",
			fileName: "api_test.go",
			content: `package main

// smoke test for API
func TestAPI(t *testing.T) {}
`,
			expected: true,
		},
		{
			name:     "smoke_in_block_comment",
			fileName: "basic_test.go",
			content: `package main

/* smoke test */
func TestBasic(t *testing.T) {}
`,
			expected: true,
		},
		{
			name:     "smoke_in_filename",
			fileName: "smoke_test.go",
			content: `package main

func TestSmoke(t *testing.T) {}
`,
			expected: true,
		},
		{
			name:     "smoketest_in_content",
			fileName: "health_test.go",
			content: `package main

func TestSmokeTest(t *testing.T) {}
`,
			expected: true,
		},
		{
			name:     "smoke_colon_tag",
			fileName: "tagged_test.go",
			content: `package main

// smoke: basic health check
func TestHealth(t *testing.T) {}
`,
			expected: true,
		},
		{
			name:     "not_smoke_test",
			fileName: "unit_test.go",
			content: `package main

import "testing"

func TestAdd(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}
`,
			expected: false,
		},
		{
			name:     "smoke_underscore",
			fileName: "check_test.go",
			content: `package main

func TestSmoke_Check(t *testing.T) {}
`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			isSmoke, err := analyzer.IsSmokeTest(context.Background(), testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isSmoke)
		})
	}
}

// TestGoTestAnalyzer_IsSmokeTest_NonExistent tests smoke detection for non-existent file
func TestGoTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/test_file.go")
	assert.Error(t, err)
}

// TestGoTestAnalyzer_getPackageName tests package name extraction
func TestGoTestAnalyzer_getPackageName(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-package-name-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "main_package",
			content: `package main

func main() {}
`,
			expected: "main",
		},
		{
			name: "custom_package",
			content: `package testengine

func Test() {}
`,
			expected: "testengine",
		},
		{
			name: "package_with_comment",
			content: `// Package utils provides utility functions
package utils

func Helper() {}
`,
			expected: "utils",
		},
		{
			name: "test_package",
			content: `package mypackage_test

import "testing"

func TestSomething(t *testing.T) {}
`,
			expected: "mypackage_test",
		},
		{
			name: "package_with_block_comment",
			content: `/*
Package docs provides documentation utilities.
*/
package docs

func Generate() {}
`,
			expected: "docs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+".go")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := analyzer.getPackageName(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGoTestAnalyzer_getPackageName_NonExistent tests package name for non-existent file
func TestGoTestAnalyzer_getPackageName_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	result := analyzer.getPackageName("/nonexistent/file.go")
	assert.Equal(t, "", result)
}

// TestGoTestAnalyzer_findTestsInPackage tests finding all tests in a package
func TestGoTestAnalyzer_findTestsInPackage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-pkg-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.go": `package mypackage

func Helper() {}
`,
		"utils_test.go": `package mypackage

import "testing"

func TestHelper(t *testing.T) {}
`,
		"other_test.go": `package mypackage

import "testing"

func TestOther(t *testing.T) {}
`,
		"external_test.go": `package mypackage_test

import "testing"

func TestExternal(t *testing.T) {}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.findTestsInPackage(context.Background(), ".", "mypackage", tempDir)
	require.NoError(t, err)

	// Should find utils_test.go and other_test.go but not external_test.go
	assert.Len(t, tests, 2)
}

// TestGoTestAnalyzer_findTestsInPackage_EmptyDir tests with empty directory
func TestGoTestAnalyzer_findTestsInPackage_EmptyDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-empty-pkg-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.findTestsInPackage(context.Background(), ".", "mypackage", tempDir)
	require.NoError(t, err)

	assert.Empty(t, tests)
}

// TestGoTestAnalyzer_FindTestsForFile_DuplicateRemoval tests that duplicates are removed
func TestGoTestAnalyzer_FindTestsForFile_DuplicateRemoval(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-dedup-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.go": `package main

func Helper() {}
`,
		"utils_test.go": `package main

import "testing"

func TestHelper(t *testing.T) {}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.go", tempDir)
	require.NoError(t, err)

	// Check no duplicates
	seen := make(map[string]bool)
	for _, test := range tests {
		assert.False(t, seen[test], "Duplicate test found: %s", test)
		seen[test] = true
	}
}

// TestGoTestAnalyzer_AnalyzeTestDependencies_ComplexImports tests complex import scenarios
func TestGoTestAnalyzer_AnalyzeTestDependencies_ComplexImports(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-complex-imports-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

import "encoding/json"

func TestComplex(t *testing.T) {}
`

	testFile := filepath.Join(tempDir, "complex_test.go")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	expectedImports := []string{
		"context", "fmt", "io", "net/http", "os",
		"path/filepath", "strings", "testing", "time", "encoding/json",
	}

	for _, expected := range expectedImports {
		assert.Contains(t, deps, expected, "Should contain import: %s", expected)
	}
}

// TestGoTestAnalyzer_IsSmokeTest_CaseInsensitive tests case insensitivity
func TestGoTestAnalyzer_IsSmokeTest_CaseInsensitive(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-smoke-case-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "uppercase_SMOKE",
			content: `package main

// SMOKE test
func TestAPI(t *testing.T) {}
`,
			expected: true,
		},
		{
			name: "mixed_case_Smoke",
			content: `package main

// Smoke Test
func TestAPI(t *testing.T) {}
`,
			expected: true,
		},
		{
			name: "lowercase_smoke",
			content: `package main

// smoke test
func TestAPI(t *testing.T) {}
`,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+"_test.go")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			isSmoke, err := analyzer.IsSmokeTest(context.Background(), testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isSmoke)
		})
	}
}

// TestGoTestAnalyzer_FindTestsForFile_WithExternalTestPackage tests with _test package
func TestGoTestAnalyzer_FindTestsForFile_WithExternalTestPackage(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-external-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.go": `package utils

func Helper() string { return "help" }
`,
		"utils_test.go": `package utils

import "testing"

func TestHelperInternal(t *testing.T) {}
`,
		"utils_external_test.go": `package utils_test

import (
	"testing"
	"myproject/utils"
)

func TestHelperExternal(t *testing.T) {
	_ = utils.Helper()
}
`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.go", tempDir)
	require.NoError(t, err)

	// Should find at least the internal test file
	assert.NotEmpty(t, tests)
}

// TestGoTestAnalyzer_AnalyzeTestDependencies_EmptyImportBlock tests empty import block
func TestGoTestAnalyzer_AnalyzeTestDependencies_EmptyImportBlock(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "go-empty-import-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `package main

import (
)

func TestEmpty(t *testing.T) {}
`

	testFile := filepath.Join(tempDir, "empty_test.go")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	analyzer := NewGoTestAnalyzer(log)

	deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
	require.NoError(t, err)

	assert.Empty(t, deps)
}

