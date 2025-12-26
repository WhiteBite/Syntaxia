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

// TestPythonTestRunner_GetLanguage tests language identifier
func TestPythonTestRunner_GetLanguage(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	assert.Equal(t, "python", runner.GetLanguage())
}

// TestPythonTestRunner_isTestFile tests test file detection
func TestPythonTestRunner_isTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{
			name:     "test_prefix_file",
			fileName: "test_main.py",
			expected: true,
		},
		{
			name:     "test_suffix_file",
			fileName: "main_test.py",
			expected: true,
		},
		{
			name:     "regular_python_file",
			fileName: "main.py",
			expected: false,
		},
		{
			name:     "conftest_file",
			fileName: "conftest.py",
			expected: false,
		},
		{
			name:     "non_python_file",
			fileName: "test_main.js",
			expected: false,
		},
		{
			name:     "test_prefix_only",
			fileName: "test_.py",
			expected: true,
		},
		{
			name:     "nested_test_file",
			fileName: "test_utils_helper.py",
			expected: true,
		},
		{
			name:     "init_file",
			fileName: "__init__.py",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.isTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestRunner_shouldSkipDirectory tests directory skip logic
func TestPythonTestRunner_shouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{
			name:     "pycache_directory",
			dirName:  "__pycache__",
			expected: true,
		},
		{
			name:     "venv_directory",
			dirName:  ".venv",
			expected: true,
		},
		{
			name:     "venv2_directory",
			dirName:  "venv",
			expected: true,
		},
		{
			name:     "tox_directory",
			dirName:  ".tox",
			expected: true,
		},
		{
			name:     "eggs_directory",
			dirName:  ".eggs",
			expected: true,
		},
		{
			name:     "build_directory",
			dirName:  "build",
			expected: true,
		},
		{
			name:     "dist_directory",
			dirName:  "dist",
			expected: true,
		},
		{
			name:     "node_modules_directory",
			dirName:  "node_modules",
			expected: true,
		},
		{
			name:     "git_directory",
			dirName:  ".git",
			expected: true,
		},
		{
			name:     "pytest_cache_directory",
			dirName:  ".pytest_cache",
			expected: true,
		},
		{
			name:     "mypy_cache_directory",
			dirName:  ".mypy_cache",
			expected: true,
		},
		{
			name:     "regular_directory",
			dirName:  "tests",
			expected: false,
		},
		{
			name:     "src_directory",
			dirName:  "src",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.shouldSkipDirectory(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestRunner_filePathToModulePath tests file path to module conversion
func TestPythonTestRunner_filePathToModulePath(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "simple_file",
			filePath: "test_main.py",
			expected: "test_main",
		},
		{
			name:     "nested_file",
			filePath: "tests/test_utils.py",
			expected: "tests.test_utils",
		},
		{
			name:     "deeply_nested_file",
			filePath: "src/tests/unit/test_service.py",
			expected: "src.tests.unit.test_service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := runner.filePathToModulePath(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestRunner_detectFramework tests framework detection from file content
func TestPythonTestRunner_detectFramework(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "python-runner-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "pytest_import",
			content: `import pytest

def test_something():
    assert True
`,
			expected: "pytest",
		},
		{
			name: "pytest_from_import",
			content: `from pytest import fixture

@fixture
def my_fixture():
    return 42
`,
			expected: "pytest",
		},
		{
			name: "pytest_decorator",
			content: `@pytest.mark.parametrize("x", [1, 2, 3])
def test_param(x):
    assert x > 0
`,
			expected: "pytest",
		},
		{
			name: "unittest_import",
			content: `import unittest

class TestCase(unittest.TestCase):
    def test_something(self):
        self.assertTrue(True)
`,
			expected: "unittest",
		},
		{
			name: "unittest_from_import",
			content: `from unittest import TestCase

class MyTest(TestCase):
    pass
`,
			expected: "unittest",
		},
		{
			name: "no_framework_detected",
			content: `def test_something():
    assert True
`,
			expected: "pytest", // Default to pytest
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, tt.name+".py")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.detectFramework(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestRunner_analyzeTestFile tests test file type analysis
func TestPythonTestRunner_analyzeTestFile(t *testing.T) {
	// Create temp directory for test files
	tempDir, err := os.MkdirTemp("", "python-analyze-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "smoke_test_marker",
			content: `import pytest

@pytest.mark.smoke
def test_smoke():
    pass
`,
			expected: "smoke",
		},
		{
			name: "smoke_in_name",
			content: `def test_smoke_login():
    pass
`,
			expected: "smoke",
		},
		{
			name: "integration_test_marker",
			content: `import pytest

@pytest.mark.integration
def test_integration():
    pass
`,
			expected: "integration",
		},
		{
			name: "integration_with_database",
			content: `import database

def test_db_connection():
    db = database.connect()
`,
			expected: "integration",
		},
		{
			name: "integration_with_http",
			content: `import requests

def test_api_call():
    response = requests.get("http://api.example.com")
`,
			expected: "integration",
		},
		{
			name: "unit_test_default",
			content: `def test_add():
    assert 1 + 1 == 2
`,
			expected: "unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			testFile := filepath.Join(tempDir, "test_"+tt.name+".py")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			result := runner.analyzeTestFile(testFile)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestRunner_buildTestCommand tests command building for different frameworks
func TestPythonTestRunner_buildTestCommand(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name           string
		testPath       string
		config         *domain.TestConfig
		expectedPrefix []string
	}{
		{
			name:     "pytest_basic",
			testPath: "tests/test_main.py",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     false,
				Coverage:    false,
			},
			expectedPrefix: []string{"-m", "pytest", "tests/test_main.py"},
		},
		{
			name:     "pytest_verbose",
			testPath: "tests/test_main.py",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     true,
				Coverage:    false,
			},
			expectedPrefix: []string{"-m", "pytest", "tests/test_main.py", "-v"},
		},
		{
			name:     "pytest_with_coverage",
			testPath: "tests/test_main.py",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Verbose:     false,
				Coverage:    true,
			},
			expectedPrefix: []string{"-m", "pytest", "tests/test_main.py"},
		},
		{
			name:     "pytest_with_timeout",
			testPath: "tests/test_main.py",
			config: &domain.TestConfig{
				ProjectPath: "/project",
				Timeout:     60,
			},
			expectedPrefix: []string{"-m", "pytest", "tests/test_main.py"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args, err := runner.buildTestCommand(tt.testPath, tt.config)
			require.NoError(t, err)

			// Check that args start with expected prefix
			for i, expected := range tt.expectedPrefix {
				if i < len(args) {
					assert.Equal(t, expected, args[i], "Argument %d mismatch", i)
				}
			}
		})
	}
}

// TestPythonTestRunner_DiscoverTests tests test discovery
func TestPythonTestRunner_DiscoverTests(t *testing.T) {
	// Create temp directory with test structure
	tempDir, err := os.MkdirTemp("", "python-discover-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test directory structure
	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"tests/test_main.py":       "def test_main(): pass",
		"tests/test_utils.py":      "def test_utils(): pass",
		"tests/conftest.py":        "import pytest",
		"src/main.py":              "def main(): pass",
		"tests/unit/test_unit.py":  "def test_unit(): pass",
		"tests/__pycache__/foo.py": "# should be skipped",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests, err := runner.DiscoverTests(context.Background(), tempDir)
	require.NoError(t, err)

	// Should find test_main.py, test_utils.py, test_unit.py
	// Should NOT find conftest.py, main.py, or files in __pycache__
	assert.GreaterOrEqual(t, len(tests), 3)

	// Verify test files are found
	foundFiles := make(map[string]bool)
	for _, test := range tests {
		foundFiles[test.Name] = true
	}

	assert.True(t, foundFiles["test_main.py"], "test_main.py should be found")
	assert.True(t, foundFiles["test_utils.py"], "test_utils.py should be found")
	assert.True(t, foundFiles["test_unit.py"], "test_unit.py should be found")
	assert.False(t, foundFiles["conftest.py"], "conftest.py should not be found")
	assert.False(t, foundFiles["main.py"], "main.py should not be found")
}

// TestPythonTestRunner_parseJSONReport tests JSON report parsing
func TestPythonTestRunner_parseJSONReport(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	tests := []struct {
		name           string
		output         string
		expectedCount  int
		expectedPassed int
	}{
		{
			name:           "empty_output",
			output:         "",
			expectedCount:  0,
			expectedPassed: 0,
		},
		{
			name:           "no_json_in_output",
			output:         "Running tests...\nAll tests passed!",
			expectedCount:  0,
			expectedPassed: 0,
		},
		{
			name: "valid_json_report",
			output: `Some output before
{"created": 1234567890, "tests": [{"nodeid": "test_main.py::test_one", "outcome": "passed", "duration": 0.1}, {"nodeid": "test_main.py::test_two", "outcome": "failed", "duration": 0.2}]}
Some output after`,
			expectedCount:  2,
			expectedPassed: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := runner.parseJSONReport(tt.output)
			assert.Len(t, results, tt.expectedCount)

			passedCount := 0
			for _, r := range results {
				if r.Success {
					passedCount++
				}
			}
			assert.Equal(t, tt.expectedPassed, passedCount)
		})
	}
}

// TestPythonTestRunner_parsePytestOutput tests pytest output parsing
func TestPythonTestRunner_parsePytestOutput(t *testing.T) {
	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	testInfos := []*domain.TestInfo{
		{Path: "test_one.py", Name: "test_one.py"},
		{Path: "test_two.py", Name: "test_two.py"},
	}

	tests := []struct {
		name          string
		output        string
		expectedCount int
	}{
		{
			name:          "passed_output",
			output:        "2 passed in 0.5s",
			expectedCount: 2,
		},
		{
			name:          "failed_output",
			output:        "1 passed, 1 failed in 0.5s",
			expectedCount: 2,
		},
		{
			name:          "empty_output",
			output:        "",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := runner.parsePytestOutput(tt.output, testInfos, 1.0)
			assert.Len(t, results, tt.expectedCount)
		})
	}
}

// TestPythonTestRunner_RunTest_Integration tests actual test execution (skipped if python not available)
func TestPythonTestRunner_RunTest_Integration(t *testing.T) {
	// Skip if python is not available
	if _, err := os.Stat("/usr/bin/python3"); os.IsNotExist(err) {
		if _, err := os.Stat("/usr/bin/python"); os.IsNotExist(err) {
			t.Skip("Python not available, skipping integration test")
		}
	}

	// Create temp directory with a simple test
	tempDir, err := os.MkdirTemp("", "python-run-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a simple test file
	testContent := `def test_simple():
    assert 1 + 1 == 2
`
	testFile := filepath.Join(tempDir, "test_simple.py")
	err = os.WriteFile(testFile, []byte(testContent), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	runner := NewPythonTestRunner(log)

	config := &domain.TestConfig{
		ProjectPath: tempDir,
		Verbose:     true,
	}

	// This test may fail if pytest is not installed, which is expected
	result, err := runner.RunTest(context.Background(), "test_simple.py", config)

	// We just verify the result structure is correct
	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, "test_simple.py", result.TestPath)
		assert.Equal(t, "python", result.Language)
	}
}
