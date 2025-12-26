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

// TestGoTestRunner_GetLanguage tests language identifier
func TestGoTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	assert.Equal(t, "go", runner.GetLanguage())
}

// TestNewGoTestRunner tests runner creation
func TestNewGoTestRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	assert.NotNil(t, runner)
	assert.NotNil(t, runner.log)
}

// TestGoTestRunner_analyzeTestFile tests test file type analysis
func TestGoTestRunner_analyzeTestFile(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "go-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "smoke_test_comment",
			content: `package main
// smoke test for basic functionality
func TestSmoke(t *testing.T) {}
`,
			expected: "smoke",
		},
		{
			name: "smoke_in_function_name",
			content: `package main
func TestSmokeBasic(t *testing.T) {}
`,
			expected: "smoke",
		},
		{
			name: "integration_test_with_testmain",
			content: `package main
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
func TestIntegration(t *testing.T) {}
`,
			expected: "integration",
		},
		{
			name: "integration_test_with_database",
			content: `package main
import "database/sql"
func TestDatabase(t *testing.T) {}
`,
			expected: "integration",
		},
		{
			name: "integration_test_with_http",
			content: `package main
import "net/http"
func TestHTTP(t *testing.T) {}
`,
			expected: "integration",
		},
		{
			name: "integration_keyword",
			content: `package main
// integration test
func TestAPI(t *testing.T) {}
`,
			expected: "integration",
		},
		{
			name: "unit_test_default",
			content: `package main
func TestAdd(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}
`,
			expected: "unit",
		},
		{
			name: "unit_test_simple",
			content: `package main
import "testing"
func TestSimple(t *testing.T) {}
`,
			expected: "unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+"_test.go")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.analyzeTestFile(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGoTestRunner_analyzeTestFile_NonExistent tests analyzing non-existent file
func TestGoTestRunner_analyzeTestFile_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	result := runner.analyzeTestFile("/nonexistent/file_test.go")
	assert.Equal(t, "unit", result) // Default to unit
}

// TestGoTestRunner_getPackageName tests package name extraction
func TestGoTestRunner_getPackageName(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "go-package-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

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
			name: "package_test",
			content: `package mypackage_test

import "testing"

func TestSomething(t *testing.T) {}
`,
			expected: "mypackage_test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".go")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.getPackageName(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGoTestRunner_getPackageName_NonExistent tests package name for non-existent file
func TestGoTestRunner_getPackageName_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	result := runner.getPackageName("/nonexistent/file.go")
	assert.Equal(t, "", result)
}

// TestGoTestRunner_DiscoverTests tests test discovery
func TestGoTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with Go project structure
	tempDir, err := os.MkdirTemp("", "go-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"main.go": `package main
func main() {}
`,
		"main_test.go": `package main
import "testing"
func TestMain(t *testing.T) {}
`,
		"utils/helper.go": `package utils
func Helper() {}
`,
		"utils/helper_test.go": `package utils
import "testing"
func TestHelper(t *testing.T) {}
`,
		"vendor/lib/lib_test.go": `package lib
func TestLib(t *testing.T) {}
`,
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find main_test.go and helper_test.go
	// Should NOT find vendor/lib/lib_test.go
	assert.Len(t, tests, 2)

	// Verify test files are found
	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["main_test.go"], "main_test.go should be found")
	assert.True(t, foundFiles["helper_test.go"], "helper_test.go should be found")
	assert.False(t, foundFiles["lib_test.go"], "vendor files should be skipped")
}

// TestGoTestRunner_DiscoverTests_SkipsNodeModules tests that node_modules is skipped
func TestGoTestRunner_DiscoverTests_SkipsNodeModules(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-skip-node-modules-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files including node_modules
	testFiles := map[string]string{
		"main_test.go": `package main
func TestMain(t *testing.T) {}
`,
		"node_modules/pkg/pkg_test.go": `package pkg
func TestPkg(t *testing.T) {}
`,
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should only find main_test.go
	assert.Len(t, tests, 1)
	assert.Equal(t, "main_test.go", tests[0].Name)
}

// TestGoTestRunner_DiscoverTests_EmptyProject tests discovery in empty project
func TestGoTestRunner_DiscoverTests_EmptyProject(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-empty-project-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	assert.Empty(t, tests)
}

// TestGoTestRunner_DiscoverTests_NonExistentPath tests discovery with non-existent path
func TestGoTestRunner_DiscoverTests_NonExistentPath(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	_, err := runner.DiscoverTests(context.Background(), "/nonexistent/path")
	assert.Error(t, err)
}

// TestGoTestRunner_RunTest tests single test execution
func TestGoTestRunner_RunTest(t *testing.T) {
	// Create temp directory with a simple Go test
	tempDir, err := os.MkdirTemp("", "go-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create go.mod
	goMod := `module testproject
go 1.21
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644)
	require.NoError(t, err)

	// Create a simple test file
	testContent := `package main

import "testing"

func TestSimple(t *testing.T) {
	if 1+1 != 2 {
		t.Error("math is broken")
	}
}
`
	err = os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
		Coverage:    true,
		Timeout:     30,
	}

	result, err := runner.RunTest(context.Background(), "./...", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "./...", result.TestPath)
	assert.Equal(t, "go", result.Language)
	// Test should pass
	assert.True(t, result.Success)
}

// TestGoTestRunner_RunTest_WithEnvVars tests test execution with environment variables
func TestGoTestRunner_RunTest_WithEnvVars(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-env-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create go.mod
	goMod := `module testproject
go 1.21
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644)
	require.NoError(t, err)

	// Create test file
	testContent := `package main

import "testing"

func TestEnv(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		EnvVars: map[string]string{
			"TEST_VAR": "test_value",
		},
	}

	result, err := runner.RunTest(context.Background(), "./...", config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestGoTestRunner_RunTest_FailingTest tests execution of failing test
func TestGoTestRunner_RunTest_FailingTest(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-fail-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create go.mod
	goMod := `module testproject
go 1.21
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644)
	require.NoError(t, err)

	// Create failing test file
	testContent := `package main

import "testing"

func TestFailing(t *testing.T) {
	t.Error("this test always fails")
}
`
	err = os.WriteFile(filepath.Join(tempDir, "main_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
	}

	result, err := runner.RunTest(context.Background(), "./...", config)

	assert.NoError(t, err) // No error from runner itself
	assert.NotNil(t, result)
	assert.False(t, result.Success) // Test should fail
	assert.NotEmpty(t, result.Error)
}

// TestGoTestRunner_RunTestSuite tests running a test suite
func TestGoTestRunner_RunTestSuite(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-suite-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create go.mod
	goMod := `module testproject
go 1.21
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644)
	require.NoError(t, err)

	// Create test files
	testContent := `package main

import "testing"

func TestOne(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "one_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	suite := &domain.TestSuite{
		Name:        "test_suite",
		Language:    "go",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "./...", Name: "all_tests", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: tempDir,
			Verbose:     false,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 1)
}

// TestGoTestRunner_RunTestSuite_Parallel tests parallel test execution
func TestGoTestRunner_RunTestSuite_Parallel(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-parallel-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create go.mod
	goMod := `module testproject
go 1.21
`
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goMod), 0o644)
	require.NoError(t, err)

	// Create test file
	testContent := `package main

import "testing"

func TestParallel(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "parallel_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	suite := &domain.TestSuite{
		Name:        "parallel_suite",
		Language:    "go",
		ProjectPath: tempDir,
		Tests: []*domain.TestInfo{
			{Path: "./...", Name: "test1", Type: "unit"},
		},
		Config: &domain.TestConfig{
			ProjectPath: tempDir,
			Parallel:    true,
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.NotNil(t, results)
}

// TestGoTestRunner_RunTestSuite_EmptySuite tests running empty test suite
func TestGoTestRunner_RunTestSuite_EmptySuite(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	suite := &domain.TestSuite{
		Name:        "empty_suite",
		Language:    "go",
		ProjectPath: "/tmp",
		Tests:       []*domain.TestInfo{},
		Config: &domain.TestConfig{
			ProjectPath: "/tmp",
		},
	}

	results, err := runner.RunTestSuite(context.Background(), suite)

	assert.NoError(t, err)
	assert.Empty(t, results)
}

// TestGoTestRunner_DiscoverTests_WithMetadata tests that discovered tests have metadata
func TestGoTestRunner_DiscoverTests_WithMetadata(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "go-metadata-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test file
	testContent := `package mypackage

import "testing"

func TestWithMetadata(t *testing.T) {}
`
	err = os.WriteFile(filepath.Join(tempDir, "meta_test.go"), []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewGoTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)
	require.Len(t, tests, 1)

	// Check metadata
	assert.Equal(t, "mypackage", tests[0].Metadata["package"])
	assert.Equal(t, "meta_test.go", tests[0].Name)
}
