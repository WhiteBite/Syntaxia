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

// TestNewPythonTestAnalyzer tests analyzer creation
func TestNewPythonTestAnalyzer(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.log)
}

// TestPythonTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestPythonTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-deps-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name            string
		content         string
		expectedImports []string
	}{
		{
			name: "simple_imports",
			content: `import os
import sys
import json

def test_something():
    pass
`,
			expectedImports: []string{"os", "sys", "json"},
		},
		{
			name: "from_imports",
			content: `from pathlib import Path
from typing import List, Dict
from collections import defaultdict

def test_something():
    pass
`,
			expectedImports: []string{"pathlib", "typing", "collections"},
		},
		{
			name: "relative_imports",
			content: `from .utils import helper
from ..models import User

def test_something():
    pass
`,
			expectedImports: []string{".utils", ".models"},
		},
		{
			name: "mixed_imports",
			content: `import pytest
from unittest import TestCase
from .helpers import setup

def test_something():
    pass
`,
			expectedImports: []string{"pytest", "unittest", ".helpers"},
		},
		{
			name:            "no_imports",
			content:         `def test_simple():\n    assert True`,
			expectedImports: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, "test_"+tt.name+".py")
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

// TestPythonTestAnalyzer_AnalyzeTestDependencies_NonExistent tests with non-existent file
func TestPythonTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/test_file.py")
	assert.Error(t, err)
}

// TestPythonTestAnalyzer_FindTestsForFile tests finding tests for source file
func TestPythonTestAnalyzer_FindTestsForFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.py":      `def helper(): pass`,
		"test_utils.py": `from utils import helper\ndef test_helper(): pass`,
		"utils_test.py": `from utils import helper\ndef test_helper2(): pass`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.py", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
}

// TestPythonTestAnalyzer_FindTestsForFile_InTestsDir tests finding tests in tests/ directory
func TestPythonTestAnalyzer_FindTestsForFile_InTestsDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-tests-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testsDir := filepath.Join(tempDir, "tests")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	files := map[string]string{
		"utils.py":            `def helper(): pass`,
		"tests/test_utils.py": `from utils import helper\ndef test_helper(): pass`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.py", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
}

// TestPythonTestAnalyzer_IsSmokeTest tests smoke test detection
func TestPythonTestAnalyzer_IsSmokeTest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-smoke-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "pytest_smoke_marker",
			fileName: "test_api.py",
			content:  `import pytest\n@pytest.mark.smoke\ndef test_api(): pass`,
			expected: true,
		},
		{
			name:     "smoke_comment",
			fileName: "test_basic.py",
			content:  `# smoke test\ndef test_basic(): pass`,
			expected: true,
		},
		{
			name:     "smoke_in_filename",
			fileName: "test_smoke.py",
			content:  `def test_smoke(): pass`,
			expected: true,
		},
		{
			name:     "smoke_in_content",
			fileName: "test_health.py",
			content:  `def test_smoke_check(): pass`,
			expected: true,
		},
		{
			name:     "not_smoke_test",
			fileName: "test_unit.py",
			content:  `def test_add():\n    assert 1 + 1 == 2`,
			expected: false,
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

// TestPythonTestAnalyzer_IsSmokeTest_NonExistent tests smoke detection for non-existent file
func TestPythonTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/test_file.py")
	assert.Error(t, err)
}

// TestPythonTestAnalyzer_extractImports tests import extraction
func TestPythonTestAnalyzer_extractImports(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		content  string
		expected []string
	}{
		{
			name:     "import_statement",
			content:  "import os",
			expected: []string{"os"},
		},
		{
			name:     "from_import",
			content:  "from pathlib import Path",
			expected: []string{"pathlib"},
		},
		{
			name:     "multiple_imports",
			content:  "import os\nimport sys\nfrom typing import List",
			expected: []string{"os", "sys", "typing"},
		},
		{
			name:     "with_comments",
			content:  "# comment\nimport os\n# another comment",
			expected: []string{"os"},
		},
		{
			name:     "empty_content",
			content:  "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractImports(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_filePathToModuleName tests file path to module name conversion
func TestPythonTestAnalyzer_filePathToModuleName(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "simple_file",
			filePath: "utils.py",
			expected: "utils",
		},
		{
			name:     "nested_file",
			filePath: "src/utils.py",
			expected: "src.utils",
		},
		{
			name:     "deeply_nested",
			filePath: "src/pkg/module.py",
			expected: "src.pkg.module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.filePathToModuleName(tt.filePath)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_shouldSkipDirectory tests directory skip logic
func TestPythonTestAnalyzer_shouldSkipDirectory(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{"__pycache__", "__pycache__", true},
		{".venv", ".venv", true},
		{"venv", "venv", true},
		{".tox", ".tox", true},
		{".eggs", ".eggs", true},
		{"build", "build", true},
		{"dist", "dist", true},
		{"node_modules", "node_modules", true},
		{".git", ".git", true},
		{".pytest_cache", ".pytest_cache", true},
		{".mypy_cache", ".mypy_cache", true},
		{"src", "src", false},
		{"tests", "tests", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.shouldSkipDirectory(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_isTestFile tests test file detection
func TestPythonTestAnalyzer_isTestFile(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		fileName string
		expected bool
	}{
		{"test_prefix", "test_utils.py", true},
		{"test_suffix", "utils_test.py", true},
		{"regular_file", "utils.py", false},
		{"conftest", "conftest.py", false},
		{"non_python", "test_utils.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.isTestFile(tt.fileName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_removeDuplicates tests duplicate removal
func TestPythonTestAnalyzer_removeDuplicates(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no_duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with_duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty_slice",
			input:    []string{},
			expected: nil,
		},
		{
			name:     "all_same",
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.removeDuplicates(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_fileExists tests file existence check
func TestPythonTestAnalyzer_fileExists(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-exists-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	existingFile := filepath.Join(tempDir, "exists.py")
	err = os.WriteFile(existingFile, []byte(""), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	assert.True(t, analyzer.fileExists(existingFile))
	assert.False(t, analyzer.fileExists(filepath.Join(tempDir, "nonexistent.py")))
}

// TestPythonTestAnalyzer_importsModule tests module import detection
func TestPythonTestAnalyzer_importsModule(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	tests := []struct {
		name       string
		content    string
		moduleName string
		expected   bool
	}{
		{
			name:       "direct_import",
			content:    "import mymodule",
			moduleName: "mymodule",
			expected:   true,
		},
		{
			name:       "from_import",
			content:    "from mymodule import something",
			moduleName: "mymodule",
			expected:   true,
		},
		{
			name:       "no_import",
			content:    "def test(): pass",
			moduleName: "mymodule",
			expected:   false,
		},
		{
			name:       "partial_match",
			content:    "import mymodule.submodule",
			moduleName: "mymodule.submodule",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.importsModule(tt.content, tt.moduleName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestPythonTestAnalyzer_ExtractTestFunctions tests test function extraction
func TestPythonTestAnalyzer_ExtractTestFunctions(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-extract-funcs-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `def test_one():
    pass

def test_two():
    pass

def helper():
    pass

def test_three():
    pass
`
	testFile := filepath.Join(tempDir, "test_funcs.py")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	funcs, err := analyzer.ExtractTestFunctions(testFile)
	require.NoError(t, err)

	assert.Len(t, funcs, 3)
	assert.Contains(t, funcs, "test_one")
	assert.Contains(t, funcs, "test_two")
	assert.Contains(t, funcs, "test_three")
	assert.NotContains(t, funcs, "helper")
}

// TestPythonTestAnalyzer_ExtractTestFunctions_NonExistent tests with non-existent file
func TestPythonTestAnalyzer_ExtractTestFunctions_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	_, err := analyzer.ExtractTestFunctions("/nonexistent/test_file.py")
	assert.Error(t, err)
}

// TestPythonTestAnalyzer_ExtractTestClasses tests test class extraction
func TestPythonTestAnalyzer_ExtractTestClasses(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "py-extract-classes-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	content := `class TestOne:
    def test_method(self):
        pass

class TestTwo(unittest.TestCase):
    def test_method(self):
        pass

class Helper:
    pass

class TestThree:
    pass
`
	testFile := filepath.Join(tempDir, "test_classes.py")
	err = os.WriteFile(testFile, []byte(content), 0o644)
	require.NoError(t, err)

	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	classes, err := analyzer.ExtractTestClasses(testFile)
	require.NoError(t, err)

	assert.Len(t, classes, 3)
	assert.Contains(t, classes, "TestOne")
	assert.Contains(t, classes, "TestTwo")
	assert.Contains(t, classes, "TestThree")
	assert.NotContains(t, classes, "Helper")
}

// TestPythonTestAnalyzer_ExtractTestClasses_NonExistent tests with non-existent file
func TestPythonTestAnalyzer_ExtractTestClasses_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewPythonTestAnalyzer(log)

	_, err := analyzer.ExtractTestClasses("/nonexistent/test_file.py")
	assert.Error(t, err)
}
