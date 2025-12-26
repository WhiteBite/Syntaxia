package testengine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
	"time"
)

// SwiftTestRunner implements domain.TestRunner for Swift projects
type SwiftTestRunner struct {
	log domain.Logger
}

// NewSwiftTestRunner creates a new Swift test runner
func NewSwiftTestRunner(log domain.Logger) *SwiftTestRunner {
	return &SwiftTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *SwiftTestRunner) GetLanguage() string {
	return "swift"
}

// RunTest executes a single Swift test
func (r *SwiftTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	startTime := time.Now()
	projectPath := config.ProjectPath

	result := &domain.TestResult{
		TestPath: testPath,
		TestName: r.extractTestName(testPath),
		Language: "swift",
	}

	// Determine build system
	buildSystem := r.detectBuildSystem(projectPath)

	var cmd *exec.Cmd
	switch buildSystem {
	case "spm":
		// Swift Package Manager: swift test --filter <test>
		testFilter := r.getTestFilter(testPath)
		args := []string{"test"}
		if testFilter != "" {
			args = append(args, "--filter", testFilter)
		}
		if config.Verbose {
			args = append(args, "-v")
		}
		cmd = exec.CommandContext(ctx, "swift", args...)
	case "xcode":
		// Xcode project: xcodebuild test
		cmd = r.createXcodebuildTestCommand(ctx, projectPath, testPath, config)
	default:
		result.Success = false
		result.Error = "no Swift build system detected (Package.swift or .xcodeproj required)"
		return result, nil
	}

	cmd.Dir = projectPath
	r.applyEnvVars(cmd, config)

	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(startTime).Seconds()

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		// Store failure message in Error field
		if failMsg := r.extractFailureMessage(string(output)); failMsg != "" {
			result.Error = failMsg
		}
		return result, nil
	}

	result.Success = true
	return result, nil
}

// RunTestSuite executes a test suite
func (r *SwiftTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Swift test suite: %s", suite.Name))

	// If no specific tests, run all tests
	if len(suite.Tests) == 0 {
		result, err := r.runAllTests(ctx, suite.ProjectPath, suite.Config)
		if err != nil {
			return nil, err
		}
		return []*domain.TestResult{result}, nil
	}

	var results []*domain.TestResult
	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// DiscoverTests finds all Swift tests in a project
func (r *SwiftTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Swift tests in: %s", projectPath))

	buildSystem := r.detectBuildSystem(projectPath)
	if buildSystem == "" {
		return nil, fmt.Errorf("no Swift build system detected (Package.swift or .xcodeproj required)")
	}

	var tests []*domain.TestInfo

	// Discover tests based on build system
	switch buildSystem {
	case "spm":
		spmTests, err := r.discoverSPMTests(projectPath)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to discover SPM tests: %v", err))
		} else {
			tests = append(tests, spmTests...)
		}
	case "xcode":
		xcodeTests, err := r.discoverXcodeTests(projectPath)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to discover Xcode tests: %v", err))
		} else {
			tests = append(tests, xcodeTests...)
		}
	}

	r.log.Info(fmt.Sprintf("Discovered %d Swift tests", len(tests)))
	return tests, nil
}

// runAllTests runs all tests in the project
func (r *SwiftTestRunner) runAllTests(ctx context.Context, projectPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	startTime := time.Now()

	result := &domain.TestResult{
		TestPath: projectPath,
		TestName: "all",
		Language: "swift",
	}

	buildSystem := r.detectBuildSystem(projectPath)

	var cmd *exec.Cmd
	switch buildSystem {
	case "spm":
		args := []string{"test"}
		if config.Verbose {
			args = append(args, "-v")
		}
		cmd = exec.CommandContext(ctx, "swift", args...)
	case "xcode":
		cmd = r.createXcodebuildTestCommand(ctx, projectPath, "", config)
	default:
		result.Success = false
		result.Error = "no Swift build system detected"
		return result, nil
	}

	cmd.Dir = projectPath
	r.applyEnvVars(cmd, config)

	output, err := cmd.CombinedOutput()
	result.Output = string(output)
	result.Duration = time.Since(startTime).Seconds()

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		if failMsg := r.extractFailureMessage(string(output)); failMsg != "" {
			result.Error = failMsg
		}
		return result, nil
	}

	result.Success = true
	return result, nil
}

// detectBuildSystem detects the Swift build system used
func (r *SwiftTestRunner) detectBuildSystem(projectPath string) string {
	// Check for Swift Package Manager
	if _, err := os.Stat(filepath.Join(projectPath, "Package.swift")); err == nil {
		return "spm"
	}

	// Check for Xcode project
	entries, err := os.ReadDir(projectPath)
	if err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".xcodeproj") ||
				strings.HasSuffix(entry.Name(), ".xcworkspace") {
				return "xcode"
			}
		}
	}

	return ""
}

// discoverSPMTests discovers tests in Swift Package Manager projects
func (r *SwiftTestRunner) discoverSPMTests(projectPath string) ([]*domain.TestInfo, error) {
	var tests []*domain.TestInfo

	testsDir := filepath.Join(projectPath, "Tests")
	if _, err := os.Stat(testsDir); os.IsNotExist(err) {
		return tests, nil
	}

	err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".swift") {
			return nil
		}

		// Check if file contains XCTest
		hasTests, _ := r.fileHasTests(path)
		if !hasTests {
			return nil
		}

		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			relPath = path
		}

		testType := "unit"
		if r.isSmokeTest(path) {
			testType = "smoke"
		} else if strings.Contains(strings.ToLower(relPath), "integration") {
			testType = "integration"
		}

		tests = append(tests, &domain.TestInfo{
			Path: relPath,
			Name: info.Name(),
			Type: testType,
			Metadata: map[string]string{
				"language": "swift",
			},
		})

		return nil
	})

	return tests, err
}

// discoverXcodeTests discovers tests in Xcode projects
func (r *SwiftTestRunner) discoverXcodeTests(projectPath string) ([]*domain.TestInfo, error) {
	var tests []*domain.TestInfo

	// Look for test targets in common locations
	testDirs := []string{
		filepath.Join(projectPath, "Tests"),
		filepath.Join(projectPath, "UnitTests"),
		filepath.Join(projectPath, "UITests"),
	}

	// Also search for *Tests directories
	entries, err := os.ReadDir(projectPath)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() && strings.HasSuffix(entry.Name(), "Tests") {
				testDirs = append(testDirs, filepath.Join(projectPath, entry.Name()))
			}
		}
	}

	for _, testDir := range testDirs {
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}

			if !strings.HasSuffix(info.Name(), ".swift") {
				return nil
			}

			hasTests, _ := r.fileHasTests(path)
			if !hasTests {
				return nil
			}

			relPath, _ := filepath.Rel(projectPath, path)

			testType := "unit"
			if strings.Contains(strings.ToLower(relPath), "uitest") {
				testType = "ui"
			} else if strings.Contains(strings.ToLower(relPath), "integration") {
				testType = "integration"
			} else if r.isSmokeTest(path) {
				testType = "smoke"
			}

			tests = append(tests, &domain.TestInfo{
				Path: relPath,
				Name: info.Name(),
				Type: testType,
				Metadata: map[string]string{
					"language": "swift",
				},
			})

			return nil
		})

		if err != nil {
			r.log.Warning(fmt.Sprintf("Error walking test directory %s: %v", testDir, err))
		}
	}

	return tests, nil
}

// fileHasTests checks if a Swift file contains XCTest tests
func (r *SwiftTestRunner) fileHasTests(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	contentStr := string(content)

	// Check for XCTest patterns
	patterns := []string{
		`import\s+XCTest`,
		`:\s*XCTestCase`,
		`func\s+test\w+\s*\(`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if re.MatchString(contentStr) {
			return true, nil
		}
	}

	return false, nil
}

// isSmokeTest checks if a test file is a smoke test
func (r *SwiftTestRunner) isSmokeTest(filePath string) bool {
	fileName := strings.ToLower(filepath.Base(filePath))
	return strings.Contains(fileName, "smoke")
}

// extractTestName extracts test name from path
func (r *SwiftTestRunner) extractTestName(testPath string) string {
	base := filepath.Base(testPath)
	return strings.TrimSuffix(base, ".swift")
}

// getTestFilter creates a test filter for swift test command
func (r *SwiftTestRunner) getTestFilter(testPath string) string {
	// Extract test class name from file
	base := filepath.Base(testPath)
	name := strings.TrimSuffix(base, ".swift")
	return name
}

// createXcodebuildTestCommand creates xcodebuild test command
func (r *SwiftTestRunner) createXcodebuildTestCommand(ctx context.Context, projectPath, testPath string, config *domain.TestConfig) *exec.Cmd {
	args := []string{"test"}

	// Find project or workspace
	entries, _ := os.ReadDir(projectPath)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".xcworkspace") {
			args = append(args, "-workspace", entry.Name())
			break
		} else if strings.HasSuffix(entry.Name(), ".xcodeproj") {
			args = append(args, "-project", entry.Name())
		}
	}

	// Add scheme if available
	args = append(args, "-scheme", r.detectScheme(projectPath))

	// Add destination
	args = append(args, "-destination", "platform=macOS")

	// Add specific test if provided
	if testPath != "" {
		testName := r.extractTestName(testPath)
		args = append(args, "-only-testing:"+testName)
	}

	return exec.CommandContext(ctx, "xcodebuild", args...)
}

// detectScheme detects the Xcode scheme to use
func (r *SwiftTestRunner) detectScheme(projectPath string) string {
	// Try to find scheme from project name
	entries, _ := os.ReadDir(projectPath)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".xcodeproj") {
			return strings.TrimSuffix(entry.Name(), ".xcodeproj")
		}
	}
	return "default"
}

// applyEnvVars applies environment variables to command
func (r *SwiftTestRunner) applyEnvVars(cmd *exec.Cmd, config *domain.TestConfig) {
	if config.EnvVars != nil {
		cmd.Env = os.Environ()
		for k, v := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}
}

// extractFailureMessage extracts failure message from test output
func (r *SwiftTestRunner) extractFailureMessage(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "error:") || strings.Contains(line, "failed") {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

// parseTestOutput parses swift test output for statistics
func (r *SwiftTestRunner) parseTestOutput(output string) (passed, failed, skipped int) {
	// Pattern: Test Suite 'All tests' passed at ...
	// Pattern: Executed X tests, with Y failures
	re := regexp.MustCompile(`Executed\s+(\d+)\s+tests?,\s+with\s+(\d+)\s+failures?`)
	matches := re.FindStringSubmatch(output)
	if len(matches) >= 3 {
		fmt.Sscanf(matches[1], "%d", &passed)
		fmt.Sscanf(matches[2], "%d", &failed)
		passed = passed - failed
	}
	return
}
