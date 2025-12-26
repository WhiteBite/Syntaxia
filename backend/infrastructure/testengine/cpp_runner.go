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
	"syntaxia/internal/executil"
	"time"
)

// C++ test runner constants
const (
	cppDefaultTimeout     = 300
	cppTestFilePattern    = "_test.cpp"
	cppTestDir            = "test"
	cppTestsDir           = "tests"
	cppBuildDir           = "build"
	cppCMakeFile          = "CMakeLists.txt"
	cppMakefile           = "Makefile"
	cppGoogleTestPattern  = "TEST\\s*\\("
	cppCatch2Pattern      = "TEST_CASE\\s*\\("
	cppBoostTestPattern   = "BOOST_AUTO_TEST_CASE\\s*\\("
)

// CppTestRunner implements domain.TestRunner for C++ projects
type CppTestRunner struct {
	log domain.Logger
}

// NewCppTestRunner creates a new C++ test runner
func NewCppTestRunner(log domain.Logger) *CppTestRunner {
	return &CppTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *CppTestRunner) GetLanguage() string {
	return "cpp"
}

// RunTest executes a single C++ test
func (r *CppTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running C++ test: %s", testPath))

	startTime := time.Now()

	// Detect build system
	buildSystem := r.detectBuildSystem(config.ProjectPath)

	result := &domain.TestResult{
		TestPath: testPath,
		TestName: filepath.Base(testPath),
		Language: "cpp",
		Metadata: map[string]interface{}{
			"buildSystem": buildSystem,
		},
	}

	// Build and run test based on build system
	var output string
	var err error

	switch buildSystem {
	case "cmake":
		output, err = r.runCMakeTest(ctx, testPath, config)
	case "make":
		output, err = r.runMakeTest(ctx, testPath, config)
	default:
		output, err = r.runDirectTest(ctx, testPath, config)
	}

	result.Duration = time.Since(startTime).Seconds()
	result.Output = output

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("C++ test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("C++ test passed for %s in %.2fs", testPath, result.Duration))
	}

	return result, nil
}

// RunTestSuite executes a test suite
func (r *CppTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running C++ test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	// Try batch execution first for CMake projects
	buildSystem := r.detectBuildSystem(suite.ProjectPath)
	if buildSystem == "cmake" {
		batchResult, err := r.runCMakeBatchTests(ctx, suite)
		if err == nil {
			return batchResult, nil
		}
		r.log.Warning(fmt.Sprintf("Batch test execution failed, falling back to individual tests: %v", err))
	}

	// Fallback: run tests individually
	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed C++ test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests finds all C++ tests in a project
func (r *CppTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering C++ tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Check for test directories
	testDirs := []string{
		filepath.Join(projectPath, cppTestDir),
		filepath.Join(projectPath, cppTestsDir),
	}

	for _, testDir := range testDirs {
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			continue
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

			// Check if file is a C++ test file
			if !r.isTestFile(info.Name()) {
				return nil
			}

			// Get relative path
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				r.log.Warning(fmt.Sprintf("Failed to get relative path for %s: %v", path, err))
				return nil
			}

			// Analyze test file for framework
			framework := r.detectTestFramework(path)

			testInfo := &domain.TestInfo{
				Path: relPath,
				Name: filepath.Base(path),
				Type: "unit",
				Metadata: map[string]string{
					"framework": framework,
				},
			}

			tests = append(tests, testInfo)
			return nil
		})

		if err != nil {
			r.log.Warning(fmt.Sprintf("Error walking test directory %s: %v", testDir, err))
		}
	}

	r.log.Info(fmt.Sprintf("Discovered %d C++ test files", len(tests)))
	return tests, nil
}

// detectBuildSystem detects the build system used in the project
func (r *CppTestRunner) detectBuildSystem(projectPath string) string {
	// Check for CMakeLists.txt
	if _, err := os.Stat(filepath.Join(projectPath, cppCMakeFile)); err == nil {
		return "cmake"
	}

	// Check for Makefile
	if _, err := os.Stat(filepath.Join(projectPath, cppMakefile)); err == nil {
		return "make"
	}

	return "direct"
}

// runCMakeTest runs a test using CMake/CTest
func (r *CppTestRunner) runCMakeTest(ctx context.Context, testPath string, config *domain.TestConfig) (string, error) {
	buildDir := filepath.Join(config.ProjectPath, cppBuildDir)

	// Ensure build directory exists and project is configured
	if err := r.ensureCMakeBuild(ctx, config.ProjectPath, buildDir); err != nil {
		return "", fmt.Errorf("failed to configure CMake: %w", err)
	}

	// Build the test target
	testName := strings.TrimSuffix(filepath.Base(testPath), ".cpp")
	cmd := exec.CommandContext(ctx, "cmake", "--build", buildDir, "--target", testName)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(buildOutput), fmt.Errorf("failed to build test: %w", err)
	}

	// Run the test executable
	testExe := filepath.Join(buildDir, testName)
	if _, err := os.Stat(testExe); os.IsNotExist(err) {
		// Try with .exe extension on Windows
		testExe = filepath.Join(buildDir, testName+".exe")
	}

	runCmd := exec.CommandContext(ctx, testExe)
	executil.HideWindow(runCmd)
	runCmd.Dir = buildDir

	output, err := runCmd.CombinedOutput()
	return string(buildOutput) + "\n" + string(output), err
}

// runMakeTest runs a test using Make
func (r *CppTestRunner) runMakeTest(ctx context.Context, testPath string, config *domain.TestConfig) (string, error) {
	testName := strings.TrimSuffix(filepath.Base(testPath), ".cpp")

	// Build the test
	cmd := exec.CommandContext(ctx, "make", testName)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(buildOutput), fmt.Errorf("failed to build test: %w", err)
	}

	// Run the test executable
	testExe := filepath.Join(config.ProjectPath, testName)
	runCmd := exec.CommandContext(ctx, testExe)
	executil.HideWindow(runCmd)
	runCmd.Dir = config.ProjectPath

	output, err := runCmd.CombinedOutput()
	return string(buildOutput) + "\n" + string(output), err
}

// runDirectTest compiles and runs a test directly with g++
func (r *CppTestRunner) runDirectTest(ctx context.Context, testPath string, config *domain.TestConfig) (string, error) {
	testName := strings.TrimSuffix(filepath.Base(testPath), ".cpp")
	outputExe := filepath.Join(config.ProjectPath, testName)

	// Detect test framework and add appropriate flags
	framework := r.detectTestFramework(filepath.Join(config.ProjectPath, testPath))
	args := []string{"-o", outputExe, testPath}

	switch framework {
	case "googletest":
		args = append(args, "-lgtest", "-lgtest_main", "-pthread")
	case "catch2":
		// Catch2 is header-only, no additional libs needed
	case "boost":
		args = append(args, "-lboost_unit_test_framework")
	}

	// Compile
	cmd := exec.CommandContext(ctx, "g++", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	buildOutput, err := cmd.CombinedOutput()
	if err != nil {
		return string(buildOutput), fmt.Errorf("failed to compile test: %w", err)
	}

	// Run
	runCmd := exec.CommandContext(ctx, outputExe)
	executil.HideWindow(runCmd)
	runCmd.Dir = config.ProjectPath

	output, err := runCmd.CombinedOutput()

	// Cleanup executable
	_ = os.Remove(outputExe)

	return string(buildOutput) + "\n" + string(output), err
}

// runCMakeBatchTests runs all tests using CTest
func (r *CppTestRunner) runCMakeBatchTests(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	buildDir := filepath.Join(suite.ProjectPath, cppBuildDir)

	// Ensure build
	if err := r.ensureCMakeBuild(ctx, suite.ProjectPath, buildDir); err != nil {
		return nil, fmt.Errorf("failed to configure CMake: %w", err)
	}

	// Build all tests
	cmd := exec.CommandContext(ctx, "cmake", "--build", buildDir)
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath

	if _, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("failed to build tests: %w", err)
	}

	// Run CTest
	startTime := time.Now()
	ctestCmd := exec.CommandContext(ctx, "ctest", "--output-on-failure")
	executil.HideWindow(ctestCmd)
	ctestCmd.Dir = buildDir

	output, err := ctestCmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Parse results
	return r.parseCTestOutput(string(output), suite.Tests, duration, err), nil
}

// ensureCMakeBuild ensures CMake project is configured and ready to build
func (r *CppTestRunner) ensureCMakeBuild(ctx context.Context, projectPath, buildDir string) error {
	// Create build directory if needed
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return fmt.Errorf("failed to create build directory: %w", err)
	}

	// Check if already configured
	if _, err := os.Stat(filepath.Join(buildDir, "CMakeCache.txt")); err == nil {
		return nil
	}

	// Configure
	cmd := exec.CommandContext(ctx, "cmake", "-S", projectPath, "-B", buildDir)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("cmake configure failed: %s", string(output))
	}

	return nil
}

// parseCTestOutput parses CTest output and creates test results
func (r *CppTestRunner) parseCTestOutput(output string, tests []*domain.TestInfo, totalDuration float64, runErr error) []*domain.TestResult {
	var results []*domain.TestResult

	avgDuration := totalDuration / float64(len(tests))
	if len(tests) == 0 {
		avgDuration = totalDuration
	}

	for _, test := range tests {
		result := &domain.TestResult{
			TestPath: test.Path,
			TestName: test.Name,
			Language: "cpp",
			Duration: avgDuration,
			Output:   output,
		}

		testName := strings.TrimSuffix(test.Name, ".cpp")
		testName = strings.TrimSuffix(testName, "_test")

		// Check test result in output
		passedPattern := regexp.MustCompile(fmt.Sprintf(`%s.*Passed`, regexp.QuoteMeta(testName)))
		failedPattern := regexp.MustCompile(fmt.Sprintf(`%s.*Failed`, regexp.QuoteMeta(testName)))

		if passedPattern.MatchString(output) {
			result.Success = true
		} else if failedPattern.MatchString(output) || runErr != nil {
			result.Success = false
			result.Error = "Test failed"
		} else {
			// Default based on overall result
			result.Success = runErr == nil
		}

		results = append(results, result)
	}

	return results
}

// detectTestFramework detects which test framework is used
func (r *CppTestRunner) detectTestFramework(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unknown"
	}

	contentStr := string(content)

	// Check for Google Test
	if matched, _ := regexp.MatchString(cppGoogleTestPattern, contentStr); matched {
		return "googletest"
	}

	// Check for Catch2
	if matched, _ := regexp.MatchString(cppCatch2Pattern, contentStr); matched {
		return "catch2"
	}

	// Check for Boost.Test
	if matched, _ := regexp.MatchString(cppBoostTestPattern, contentStr); matched {
		return "boost"
	}

	return "unknown"
}

// shouldSkipDirectory checks if directory should be skipped during discovery
func (r *CppTestRunner) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		"build",
		".git",
		"cmake-build-debug",
		"cmake-build-release",
		"out",
		"bin",
		"lib",
		"third_party",
		"vendor",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// isTestFile checks if file is a C++ test file
func (r *CppTestRunner) isTestFile(name string) bool {
	// Check for common test file patterns
	testPatterns := []string{
		"_test.cpp",
		"_test.cc",
		"_test.cxx",
		"test_",
		"_tests.cpp",
		"_tests.cc",
	}

	nameLower := strings.ToLower(name)
	for _, pattern := range testPatterns {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	return false
}
