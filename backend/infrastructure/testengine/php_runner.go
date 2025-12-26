package testengine

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// PHP test runner constants
const (
	// Default timeout for PHP tests in seconds
	phpDefaultTimeout = 300

	// Test file patterns
	phpTestFileSuffix = "Test.php"

	// PHPUnit configuration files
	phpunitXML     = "phpunit.xml"
	phpunitXMLDist = "phpunit.xml.dist"
)

// PHPTestRunner implements domain.TestRunner for PHP projects
type PHPTestRunner struct {
	log domain.Logger
}

// NewPHPTestRunner creates a new PHP test runner
func NewPHPTestRunner(log domain.Logger) *PHPTestRunner {
	return &PHPTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *PHPTestRunner) GetLanguage() string {
	return "php"
}

// RunTest executes a single PHP test
func (r *PHPTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running PHP test: %s", testPath))

	startTime := time.Now()

	// Build PHPUnit command
	args, err := r.buildTestCommand(testPath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build test command: %w", err)
	}

	// Get PHPUnit executable path
	phpunitPath := r.findPHPUnitPath(config.ProjectPath)

	// Create command
	cmd := exec.CommandContext(ctx, phpunitPath, args...)
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
		Language: "php",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("PHP test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("PHP test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite executes a test suite
func (r *PHPTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running PHP test suite with %d tests", len(suite.Tests)))

	// Check if PHPUnit is available for batch execution
	if r.isPHPUnitAvailable(suite.Config.ProjectPath) {
		// Run all tests at once with PHPUnit for better performance
		batchResult, err := r.runPHPUnitBatch(ctx, suite)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Batch PHPUnit execution failed, falling back to individual tests: %v", err))
		} else {
			return batchResult, nil
		}
	}

	// Fallback: run tests individually
	var results []*domain.TestResult
	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed PHP test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests finds all PHP tests in a project
func (r *PHPTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering PHP tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Standard test directories
	testDirs := []string{
		"tests",
		"test",
		"Tests",
		"Test",
		"src/Test",
		"src/Tests",
	}

	for _, testDir := range testDirs {
		fullTestDir := filepath.Join(projectPath, testDir)
		if _, err := os.Stat(fullTestDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(fullTestDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// Skip vendor and cache directories
				if r.shouldSkipDirectory(info.Name()) {
					return filepath.SkipDir
				}
				return nil
			}

			// Check if file is a PHP test file
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
					"framework": "phpunit",
				},
			}

			tests = append(tests, testInfo)
			return nil
		})

		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to walk test directory %s: %v", testDir, err))
		}
	}

	r.log.Info(fmt.Sprintf("Discovered %d PHP test files", len(tests)))
	return tests, nil
}

// buildTestCommand builds the command arguments for running a test
func (r *PHPTestRunner) buildTestCommand(testPath string, config *domain.TestConfig) ([]string, error) {
	args := []string{}

	// Add configuration file if exists
	configFile := r.findPHPUnitConfig(config.ProjectPath)
	if configFile != "" {
		args = append(args, "--configuration", configFile)
	}

	// Add test path
	args = append(args, testPath)

	// Add verbose flag
	if config.Verbose {
		args = append(args, "--verbose")
	}

	// Add coverage if requested
	if config.Coverage {
		args = append(args, "--coverage-text")
	}

	// Add testdox for better output
	args = append(args, "--testdox")

	return args, nil
}

// runPHPUnitBatch runs all tests in batch using PHPUnit
func (r *PHPTestRunner) runPHPUnitBatch(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	startTime := time.Now()

	phpunitPath := r.findPHPUnitPath(suite.ProjectPath)
	args := []string{}

	// Add configuration file if exists
	configFile := r.findPHPUnitConfig(suite.ProjectPath)
	if configFile != "" {
		args = append(args, "--configuration", configFile)
	}

	// Add JSON log for parsing
	logFile := filepath.Join(os.TempDir(), fmt.Sprintf("phpunit-log-%d.json", time.Now().UnixNano()))
	args = append(args, "--log-junit", logFile)

	// Add all test paths
	for _, test := range suite.Tests {
		args = append(args, test.Path)
	}

	if suite.Config != nil && suite.Config.Verbose {
		args = append(args, "--verbose")
	}

	cmd := exec.CommandContext(ctx, phpunitPath, args...)
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Try to parse JUnit XML log
	results := r.parseJUnitXML(logFile, suite.Tests, duration)

	// Clean up log file
	_ = os.Remove(logFile)

	if err != nil && len(results) == 0 {
		// If parsing failed and command failed, parse text output
		results = r.parseTextOutput(string(output), suite.Tests, duration)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("PHPUnit batch execution failed: %w", err)
	}

	return results, nil
}

// parseJUnitXML parses PHPUnit JUnit XML output
func (r *PHPTestRunner) parseJUnitXML(logFile string, tests []*domain.TestInfo, totalDuration float64) []*domain.TestResult {
	content, err := os.ReadFile(logFile)
	if err != nil {
		return nil
	}

	var testSuites struct {
		XMLName    xml.Name `xml:"testsuites"`
		TestSuites []struct {
			Name      string `xml:"name,attr"`
			Tests     int    `xml:"tests,attr"`
			Failures  int    `xml:"failures,attr"`
			Errors    int    `xml:"errors,attr"`
			Time      string `xml:"time,attr"`
			TestCases []struct {
				Name      string `xml:"name,attr"`
				Class     string `xml:"class,attr"`
				File      string `xml:"file,attr"`
				Time      string `xml:"time,attr"`
				Failure   *struct {
					Message string `xml:"message,attr"`
					Type    string `xml:"type,attr"`
					Content string `xml:",chardata"`
				} `xml:"failure"`
				Error *struct {
					Message string `xml:"message,attr"`
					Type    string `xml:"type,attr"`
					Content string `xml:",chardata"`
				} `xml:"error"`
			} `xml:"testcase"`
		} `xml:"testsuite"`
	}

	if err := xml.Unmarshal(content, &testSuites); err != nil {
		return nil
	}

	var results []*domain.TestResult
	for _, suite := range testSuites.TestSuites {
		for _, tc := range suite.TestCases {
			var duration float64
			_, _ = fmt.Sscanf(tc.Time, "%f", &duration)

			result := &domain.TestResult{
				TestPath: tc.File,
				TestName: tc.Name,
				Language: "php",
				Duration: duration,
				Success:  tc.Failure == nil && tc.Error == nil,
			}

			if tc.Failure != nil {
				result.Error = tc.Failure.Message
				result.Output = tc.Failure.Content
			}
			if tc.Error != nil {
				result.Error = tc.Error.Message
				result.Output = tc.Error.Content
			}

			results = append(results, result)
		}
	}

	return results
}

// parseTextOutput parses PHPUnit text output
func (r *PHPTestRunner) parseTextOutput(output string, tests []*domain.TestInfo, totalDuration float64) []*domain.TestResult {
	var results []*domain.TestResult
	avgDuration := totalDuration / float64(len(tests))

	for _, test := range tests {
		result := &domain.TestResult{
			TestPath: test.Path,
			TestName: test.Name,
			Language: "php",
			Duration: avgDuration,
			Output:   output,
		}

		// Check if test passed or failed based on output
		if strings.Contains(output, "OK") && !strings.Contains(output, "FAILURES") {
			result.Success = true
		} else if strings.Contains(output, "FAILURES") || strings.Contains(output, "ERRORS") {
			result.Success = false
			result.Error = "Test failed"
		} else {
			result.Success = true
		}

		results = append(results, result)
	}

	return results
}

// shouldSkipDirectory checks if directory should be skipped during discovery
func (r *PHPTestRunner) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		"vendor",
		"node_modules",
		".git",
		"cache",
		".phpunit.cache",
		"var",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// isTestFile checks if file is a PHP test file
func (r *PHPTestRunner) isTestFile(name string) bool {
	if !strings.HasSuffix(name, ".php") {
		return false
	}

	// Check for *Test.php pattern
	if strings.HasSuffix(name, phpTestFileSuffix) {
		return true
	}

	return false
}

// analyzeTestFile analyzes test file to determine its type
func (r *PHPTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit"
	}

	contentStr := strings.ToLower(string(content))
	fileName := strings.ToLower(filepath.Base(filePath))

	// Check for smoke tests
	if strings.Contains(contentStr, "smoke") || strings.Contains(fileName, "smoke") {
		return "smoke"
	}

	// Check for integration tests
	integrationIndicators := []string{
		"integration",
		"database",
		"@group integration",
		"kerneltestcase",
		"webtestcase",
		"functionaltestcase",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentStr, indicator) || strings.Contains(fileName, indicator) {
			return "integration"
		}
	}

	// Check for feature tests (Laravel)
	if strings.Contains(fileName, "feature") || strings.Contains(contentStr, "featuretest") {
		return "feature"
	}

	return "unit"
}

// findPHPUnitPath finds the PHPUnit executable
func (r *PHPTestRunner) findPHPUnitPath(projectPath string) string {
	// Check for vendor/bin/phpunit (Composer)
	vendorPath := filepath.Join(projectPath, "vendor", "bin", "phpunit")
	if _, err := os.Stat(vendorPath); err == nil {
		return vendorPath
	}

	// Check for Windows variant
	vendorPathWin := filepath.Join(projectPath, "vendor", "bin", "phpunit.bat")
	if _, err := os.Stat(vendorPathWin); err == nil {
		return vendorPathWin
	}

	// Fallback to global phpunit
	return "phpunit"
}

// findPHPUnitConfig finds PHPUnit configuration file
func (r *PHPTestRunner) findPHPUnitConfig(projectPath string) string {
	configFiles := []string{phpunitXML, phpunitXMLDist}

	for _, configFile := range configFiles {
		configPath := filepath.Join(projectPath, configFile)
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	return ""
}

// isPHPUnitAvailable checks if PHPUnit is available in the project
func (r *PHPTestRunner) isPHPUnitAvailable(projectPath string) bool {
	phpunitPath := r.findPHPUnitPath(projectPath)

	cmd := exec.Command(phpunitPath, "--version")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	err := cmd.Run()
	return err == nil
}

// PHPUnitJSONResult represents PHPUnit JSON output structure
type PHPUnitJSONResult struct {
	Event   string  `json:"event"`
	Suite   string  `json:"suite,omitempty"`
	Test    string  `json:"test,omitempty"`
	Status  string  `json:"status,omitempty"`
	Time    float64 `json:"time,omitempty"`
	Message string  `json:"message,omitempty"`
	Trace   string  `json:"trace,omitempty"`
}

// parseJSONOutput parses PHPUnit JSON output
func (r *PHPTestRunner) parseJSONOutput(output string) []*domain.TestResult {
	var results []*domain.TestResult

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var jsonResult PHPUnitJSONResult
		if err := json.Unmarshal([]byte(line), &jsonResult); err != nil {
			continue
		}

		if jsonResult.Event == "test" && jsonResult.Test != "" {
			result := &domain.TestResult{
				TestPath: jsonResult.Test,
				TestName: jsonResult.Test,
				Language: "php",
				Duration: jsonResult.Time,
				Success:  jsonResult.Status == "pass",
			}

			if jsonResult.Status != "pass" {
				result.Error = jsonResult.Message
			}

			results = append(results, result)
		}
	}

	return results
}
