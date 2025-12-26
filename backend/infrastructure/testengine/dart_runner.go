package testengine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// Dart test runner constants
const (
	dartDefaultTimeout   = 300
	dartTestFilePattern  = "_test.dart"
	dartTestDir          = "test"
	dartSkipBuild        = "build"
	dartSkipDartTool     = ".dart_tool"
	dartSkipPackages     = ".packages"
	dartPubspecFile      = "pubspec.yaml"
	flutterTestCommand   = "flutter"
	dartTestCommand      = "dart"
)

// DartTestRunner implements domain.TestRunner for Dart/Flutter projects
type DartTestRunner struct {
	log domain.Logger
}

// NewDartTestRunner creates a new Dart test runner
func NewDartTestRunner(log domain.Logger) *DartTestRunner {
	return &DartTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *DartTestRunner) GetLanguage() string {
	return "dart"
}

// RunTest executes a single Dart test
func (r *DartTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Dart test: %s", testPath))

	startTime := time.Now()

	// Detect if Flutter or pure Dart project
	isFlutter := r.isFlutterProject(config.ProjectPath)

	// Build command
	var cmd *exec.Cmd
	if isFlutter {
		cmd = exec.CommandContext(ctx, flutterTestCommand, "test", testPath)
	} else {
		cmd = exec.CommandContext(ctx, dartTestCommand, "test", testPath)
	}
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
		Language: "dart",
		Duration: duration,
		Output:   string(output),
		Metadata: map[string]interface{}{
			"isFlutter": isFlutter,
		},
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Dart test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Dart test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite executes a test suite
func (r *DartTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Dart test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	// Try batch execution first
	batchResult, err := r.runBatchTests(ctx, suite)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Batch test execution failed, falling back to individual tests: %v", err))
	} else {
		return batchResult, nil
	}

	// Fallback: run tests individually
	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed Dart test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests finds all Dart tests in a project
func (r *DartTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Dart tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Check for test directory
	testDir := filepath.Join(projectPath, dartTestDir)
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		r.log.Info("No test/ directory found")
		return tests, nil
	}

	err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if r.shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file is a Dart test file
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
				"framework": "dart_test",
			},
		}

		tests = append(tests, testInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover Dart tests: %w", err)
	}

	r.log.Info(fmt.Sprintf("Discovered %d Dart test files", len(tests)))
	return tests, nil
}

// runBatchTests runs all tests in batch
func (r *DartTestRunner) runBatchTests(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	startTime := time.Now()

	isFlutter := r.isFlutterProject(suite.ProjectPath)

	var cmd *exec.Cmd
	if isFlutter {
		cmd = exec.CommandContext(ctx, flutterTestCommand, "test", "--reporter=expanded")
	} else {
		cmd = exec.CommandContext(ctx, dartTestCommand, "test", "--reporter=expanded")
	}
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Parse output and create results
	results := r.parseTestOutput(string(output), suite.Tests, duration)

	if err != nil && len(results) == 0 {
		return nil, fmt.Errorf("dart test batch execution failed: %w", err)
	}

	return results, nil
}

// parseTestOutput parses dart test output and creates test results
func (r *DartTestRunner) parseTestOutput(output string, tests []*domain.TestInfo, totalDuration float64) []*domain.TestResult {
	var results []*domain.TestResult

	avgDuration := totalDuration / float64(len(tests))

	for _, test := range tests {
		result := &domain.TestResult{
			TestPath: test.Path,
			TestName: test.Name,
			Language: "dart",
			Duration: avgDuration,
			Output:   output,
		}

		// Check if test passed or failed based on output
		testName := strings.TrimSuffix(test.Name, "_test.dart")
		if strings.Contains(output, "All tests passed") {
			result.Success = true
		} else if strings.Contains(output, testName) && strings.Contains(output, "FAILED") {
			result.Success = false
			result.Error = "Test failed"
		} else if strings.Contains(output, "Some tests failed") {
			result.Success = false
			result.Error = "Some tests failed"
		} else {
			result.Success = !strings.Contains(output, "FAILED")
		}

		results = append(results, result)
	}

	return results
}

// isFlutterProject checks if the project is a Flutter project
func (r *DartTestRunner) isFlutterProject(projectPath string) bool {
	pubspecPath := filepath.Join(projectPath, dartPubspecFile)
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return false
	}

	// Check for flutter dependency
	return strings.Contains(string(content), "flutter:")
}

// shouldSkipDirectory checks if directory should be skipped during discovery
func (r *DartTestRunner) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		dartSkipBuild,
		dartSkipDartTool,
		dartSkipPackages,
		".git",
		"ios",
		"android",
		"web",
		"linux",
		"macos",
		"windows",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// isTestFile checks if file is a Dart test file
func (r *DartTestRunner) isTestFile(name string) bool {
	return strings.HasSuffix(name, dartTestFilePattern)
}

// analyzeTestFile analyzes test file to determine its type
func (r *DartTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit"
	}

	contentStr := strings.ToLower(string(content))

	// Check for widget tests
	if strings.Contains(contentStr, "widgettester") ||
		strings.Contains(contentStr, "testwidgets") ||
		strings.Contains(contentStr, "pumpwidget") {
		return "widget"
	}

	// Check for integration tests
	if strings.Contains(contentStr, "integration") ||
		strings.Contains(contentStr, "integrationtestwidgetsbinding") {
		return "integration"
	}

	// Check for smoke tests
	if strings.Contains(contentStr, "smoke") {
		return "smoke"
	}

	return "unit"
}
