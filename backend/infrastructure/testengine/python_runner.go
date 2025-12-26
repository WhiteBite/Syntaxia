package testengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// Python test runner constants
const (
	// Default timeout for Python tests in seconds
	pythonDefaultTimeout = 300

	// Test file patterns
	pythonTestFilePrefix = "test_"
	pythonTestFileSuffix = "_test.py"

	// Test function/class patterns
	pythonTestFuncPrefix  = "test_"
	pythonTestClassPrefix = "Test"

	// Directories to skip during discovery
	pythonSkipPycache = "__pycache__"
	pythonSkipVenv    = ".venv"
	pythonSkipVenv2   = "venv"
	pythonSkipTox     = ".tox"
	pythonSkipEggs    = ".eggs"
	pythonSkipBuild   = "build"
	pythonSkipDist    = "dist"
)

// PythonTestRunner implements domain.TestRunner for Python projects
type PythonTestRunner struct {
	log domain.Logger
}

// NewPythonTestRunner creates a new Python test runner
func NewPythonTestRunner(log domain.Logger) *PythonTestRunner {
	return &PythonTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *PythonTestRunner) GetLanguage() string {
	return "python"
}

// RunTest executes a single Python test
func (r *PythonTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Python test: %s", testPath))

	startTime := time.Now()

	// Detect test framework and build command
	args, err := r.buildTestCommand(testPath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build test command: %w", err)
	}

	// Create command
	cmd := exec.CommandContext(ctx, "python", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Set environment variables
	cmd.Env = os.Environ()
	if config.EnvVars != nil {
		for key, value := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Run command
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	result := &domain.TestResult{
		TestPath: testPath,
		TestName: filepath.Base(testPath),
		Language: "python",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Python test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Python test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite executes a test suite
func (r *PythonTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Python test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	// Check if pytest is available for batch execution
	if r.isPytestAvailable(suite.Config.ProjectPath) {
		// Run all tests at once with pytest for better performance
		batchResult, err := r.runPytestBatch(ctx, suite)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Batch pytest execution failed, falling back to individual tests: %v", err))
		} else {
			return batchResult, nil
		}
	}

	// Fallback: run tests individually
	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed Python test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests finds all Python tests in a project
func (r *PythonTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Python tests in: %s", projectPath))

	var tests []*domain.TestInfo

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip excluded directories
			if r.shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file is a Python test file
		if !r.isTestFile(info.Name()) {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to get relative path for %s: %v", path, err))
			return nil
		}

		// Analyze test file for type
		testType := r.analyzeTestFile(path)

		testInfo := &domain.TestInfo{
			Path: relPath,
			Name: filepath.Base(path),
			Type: testType,
			Metadata: map[string]string{
				"framework": r.detectFramework(path),
			},
		}

		tests = append(tests, testInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover Python tests: %w", err)
	}

	r.log.Info(fmt.Sprintf("Discovered %d Python test files", len(tests)))
	return tests, nil
}

// buildTestCommand builds the command arguments for running a test
func (r *PythonTestRunner) buildTestCommand(testPath string, config *domain.TestConfig) ([]string, error) {
	args := []string{"-m"}

	// Detect framework
	framework := r.detectFrameworkFromConfig(config, testPath)

	switch framework {
	case "pytest":
		args = append(args, "pytest")
		args = append(args, testPath)

		if config.Verbose {
			args = append(args, "-v")
		}

		args = append(args, "--tb=short", "-q")

		if config.Coverage {
			args = append(args, "--cov", "--cov-report=term-missing")
		}

		if config.Timeout > 0 {
			args = append(args, fmt.Sprintf("--timeout=%d", config.Timeout))
		}

	case "unittest":
		args = append(args, "unittest")
		// Convert file path to module path for unittest
		modulePath := r.filePathToModulePath(testPath)
		args = append(args, modulePath)

		if config.Verbose {
			args = append(args, "-v")
		}

	default:
		// Default to pytest
		args = append(args, "pytest", testPath, "-v", "--tb=short")
	}

	return args, nil
}

// runPytestBatch runs all tests in batch using pytest
func (r *PythonTestRunner) runPytestBatch(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	startTime := time.Now()

	args := []string{"-m", "pytest"}

	// Add all test paths
	for _, test := range suite.Tests {
		args = append(args, test.Path)
	}

	// Add common options
	args = append(args, "-v", "--tb=short")

	// Add JSON report for parsing
	args = append(args, "--json-report", "--json-report-file=-")

	if suite.Config != nil {
		if suite.Config.Coverage {
			args = append(args, "--cov", "--cov-report=term-missing")
		}
		if suite.Config.Timeout > 0 {
			args = append(args, fmt.Sprintf("--timeout=%d", suite.Config.Timeout))
		}
	}

	cmd := exec.CommandContext(ctx, "python", args...)
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Try to parse JSON report
	results := r.parsePytestOutput(string(output), suite.Tests, duration)

	if err != nil && len(results) == 0 {
		// If parsing failed and command failed, return error
		return nil, fmt.Errorf("pytest batch execution failed: %w", err)
	}

	return results, nil
}

// parsePytestOutput parses pytest output and creates test results
func (r *PythonTestRunner) parsePytestOutput(output string, tests []*domain.TestInfo, totalDuration float64) []*domain.TestResult {
	var results []*domain.TestResult

	// Try to parse JSON report first
	jsonResults := r.parseJSONReport(output)
	if len(jsonResults) > 0 {
		return jsonResults
	}

	// Fallback: parse text output
	avgDuration := totalDuration / float64(len(tests))

	for _, test := range tests {
		result := &domain.TestResult{
			TestPath: test.Path,
			TestName: test.Name,
			Language: "python",
			Duration: avgDuration,
			Output:   output,
		}

		// Check if test passed or failed based on output
		if strings.Contains(output, "PASSED") || strings.Contains(output, "passed") {
			result.Success = true
		} else if strings.Contains(output, "FAILED") || strings.Contains(output, "failed") {
			result.Success = false
			result.Error = "Test failed"
		} else {
			// Unknown status, assume passed if no error
			result.Success = true
		}

		results = append(results, result)
	}

	return results
}

// parseJSONReport parses pytest JSON report
func (r *PythonTestRunner) parseJSONReport(output string) []*domain.TestResult {
	// Find JSON in output
	jsonStart := strings.Index(output, "{\"created\":")
	if jsonStart == -1 {
		return nil
	}

	jsonEnd := strings.LastIndex(output, "}")
	if jsonEnd == -1 || jsonEnd <= jsonStart {
		return nil
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var report struct {
		Tests []struct {
			NodeID   string  `json:"nodeid"`
			Outcome  string  `json:"outcome"`
			Duration float64 `json:"duration"`
		} `json:"tests"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &report); err != nil {
		return nil
	}

	var results []*domain.TestResult
	for _, test := range report.Tests {
		result := &domain.TestResult{
			TestPath: test.NodeID,
			TestName: filepath.Base(test.NodeID),
			Language: "python",
			Duration: test.Duration,
			Success:  test.Outcome == "passed",
		}

		if test.Outcome != "passed" {
			result.Error = fmt.Sprintf("Test %s: %s", test.NodeID, test.Outcome)
		}

		results = append(results, result)
	}

	return results
}

// shouldSkipDirectory checks if directory should be skipped during discovery
func (r *PythonTestRunner) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		pythonSkipPycache,
		pythonSkipVenv,
		pythonSkipVenv2,
		pythonSkipTox,
		pythonSkipEggs,
		pythonSkipBuild,
		pythonSkipDist,
		"node_modules",
		".git",
		".pytest_cache",
		".mypy_cache",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// isTestFile checks if file is a Python test file
func (r *PythonTestRunner) isTestFile(name string) bool {
	if !strings.HasSuffix(name, ".py") {
		return false
	}

	// Check for test_*.py pattern
	if strings.HasPrefix(name, pythonTestFilePrefix) {
		return true
	}

	// Check for *_test.py pattern
	if strings.HasSuffix(name, pythonTestFileSuffix) {
		return true
	}

	// Check for conftest.py (pytest configuration)
	if name == "conftest.py" {
		return false // Not a test file, but pytest config
	}

	return false
}

// analyzeTestFile analyzes test file to determine its type
func (r *PythonTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit" // Default to unit test
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke tests
	if strings.Contains(contentStr, "smoke") ||
		strings.Contains(contentStr, "@pytest.mark.smoke") {
		return "smoke"
	}

	// Check for integration tests
	if strings.Contains(contentStr, "integration") ||
		strings.Contains(contentStr, "@pytest.mark.integration") ||
		strings.Contains(contentStr, "database") ||
		strings.Contains(contentStr, "http") ||
		strings.Contains(contentStr, "requests") {
		return "integration"
	}

	// Default to unit test
	return "unit"
}

// detectFramework detects the test framework used in a file
func (r *PythonTestRunner) detectFramework(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "pytest" // Default to pytest
	}

	contentStr := string(content)

	// Check for pytest imports/decorators
	if strings.Contains(contentStr, "import pytest") ||
		strings.Contains(contentStr, "from pytest") ||
		strings.Contains(contentStr, "@pytest") {
		return "pytest"
	}

	// Check for unittest imports
	if strings.Contains(contentStr, "import unittest") ||
		strings.Contains(contentStr, "from unittest") ||
		strings.Contains(contentStr, "unittest.TestCase") {
		return "unittest"
	}

	// Default to pytest (more common)
	return "pytest"
}

// detectFrameworkFromConfig detects framework from config or file
func (r *PythonTestRunner) detectFrameworkFromConfig(config *domain.TestConfig, testPath string) string {
	// Check if pytest is available in project
	if r.isPytestAvailable(config.ProjectPath) {
		return "pytest"
	}

	// Fallback to file analysis
	fullPath := filepath.Join(config.ProjectPath, testPath)
	return r.detectFramework(fullPath)
}

// isPytestAvailable checks if pytest is available in the project
func (r *PythonTestRunner) isPytestAvailable(projectPath string) bool {
	cmd := exec.Command("python", "-m", "pytest", "--version")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	err := cmd.Run()
	return err == nil
}

// filePathToModulePath converts file path to Python module path
func (r *PythonTestRunner) filePathToModulePath(filePath string) string {
	// Remove .py extension
	modulePath := strings.TrimSuffix(filePath, ".py")

	// Replace path separators with dots
	modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")
	modulePath = strings.ReplaceAll(modulePath, "/", ".")

	return modulePath
}
