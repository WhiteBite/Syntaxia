package testengine

import (
	"context"
	"encoding/json"
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

// Ruby test runner constants
const (
	rubyDefaultTimeout  = 300
	rubyTestFileSuffix  = "_spec.rb"
	rubyMinitestSuffix  = "_test.rb"
	rspecConfigFile     = ".rspec"
	gemfileFile         = "Gemfile"
)

// RubyTestRunner implements domain.TestRunner for Ruby projects
type RubyTestRunner struct {
	log domain.Logger
}

// NewRubyTestRunner creates a new Ruby test runner
func NewRubyTestRunner(log domain.Logger) *RubyTestRunner {
	return &RubyTestRunner{log: log}
}

// GetLanguage returns the language identifier
func (r *RubyTestRunner) GetLanguage() string {
	return "ruby"
}

// RunTest executes a single Ruby test
func (r *RubyTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Ruby test: %s", testPath))

	startTime := time.Now()

	// Determine test framework
	framework := r.detectTestFramework(config.ProjectPath, testPath)

	// Build command based on framework
	cmdName, args := r.buildTestCommand(testPath, config, framework)

	// Create command
	cmd := exec.CommandContext(ctx, cmdName, args...)
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
		Language: "ruby",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Ruby test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Ruby test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite executes a test suite
func (r *RubyTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Ruby test suite with %d tests", len(suite.Tests)))

	// Determine framework
	framework := r.detectProjectFramework(suite.ProjectPath)

	// Try batch execution
	if framework == "rspec" && r.isRSpecAvailable(suite.ProjectPath) {
		batchResult, err := r.runRSpecBatch(ctx, suite)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Batch RSpec execution failed: %v", err))
		} else {
			return batchResult, nil
		}
	}

	if framework == "minitest" && r.isMinitestAvailable(suite.ProjectPath) {
		batchResult, err := r.runMinitestBatch(ctx, suite)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Batch Minitest execution failed: %v", err))
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

	r.log.Info(fmt.Sprintf("Completed Ruby test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests finds all Ruby tests in a project
func (r *RubyTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Ruby tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Standard test directories
	testDirs := []string{
		"spec",
		"test",
		"tests",
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
				if r.shouldSkipDirectory(info.Name()) {
					return filepath.SkipDir
				}
				return nil
			}

			// Check if file is a Ruby test file
			if !r.isTestFile(info.Name()) {
				return nil
			}

			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				r.log.Warning(fmt.Sprintf("Failed to get relative path for %s: %v", path, err))
				return nil
			}

			testType := r.analyzeTestFile(path)
			framework := r.detectTestFramework(projectPath, relPath)

			testInfo := &domain.TestInfo{
				Path: relPath,
				Name: filepath.Base(path),
				Type: testType,
				Metadata: map[string]string{
					"framework": framework,
				},
			}

			tests = append(tests, testInfo)
			return nil
		})

		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to walk test directory %s: %v", testDir, err))
		}
	}

	r.log.Info(fmt.Sprintf("Discovered %d Ruby test files", len(tests)))
	return tests, nil
}

// detectTestFramework determines which test framework to use
func (r *RubyTestRunner) detectTestFramework(projectPath, testPath string) string {
	// Check file suffix
	if strings.HasSuffix(testPath, rubyTestFileSuffix) {
		return "rspec"
	}
	if strings.HasSuffix(testPath, rubyMinitestSuffix) {
		return "minitest"
	}

	// Check directory
	if strings.Contains(testPath, "spec/") || strings.HasPrefix(testPath, "spec") {
		return "rspec"
	}

	// Check Gemfile for framework
	return r.detectProjectFramework(projectPath)
}

// detectProjectFramework detects the test framework from project configuration
func (r *RubyTestRunner) detectProjectFramework(projectPath string) string {
	gemfilePath := filepath.Join(projectPath, gemfileFile)
	content, err := os.ReadFile(gemfilePath)
	if err != nil {
		// Default to RSpec if no Gemfile
		if _, err := os.Stat(filepath.Join(projectPath, "spec")); err == nil {
			return "rspec"
		}
		return "minitest"
	}

	contentStr := string(content)
	if strings.Contains(contentStr, "rspec") {
		return "rspec"
	}
	if strings.Contains(contentStr, "minitest") {
		return "minitest"
	}

	// Check for spec directory
	if _, err := os.Stat(filepath.Join(projectPath, "spec")); err == nil {
		return "rspec"
	}

	return "minitest"
}

// buildTestCommand builds the command for running a test
func (r *RubyTestRunner) buildTestCommand(testPath string, config *domain.TestConfig, framework string) (string, []string) {
	useBundler := r.hasBundler(config.ProjectPath)

	if framework == "rspec" {
		return r.buildRSpecCommand(testPath, config, useBundler)
	}
	return r.buildMinitestCommand(testPath, config, useBundler)
}

// buildRSpecCommand builds RSpec command
func (r *RubyTestRunner) buildRSpecCommand(testPath string, config *domain.TestConfig, useBundler bool) (string, []string) {
	args := []string{}

	if useBundler {
		args = append(args, "exec", "rspec")
	}

	// Add format for JSON output
	args = append(args, "--format", "json")
	args = append(args, "--format", "progress")

	// Add test path
	args = append(args, testPath)

	if useBundler {
		return "bundle", args
	}
	// Without bundler, args already contains just the rspec arguments
	return "rspec", args
}

// buildMinitestCommand builds Minitest command
func (r *RubyTestRunner) buildMinitestCommand(testPath string, config *domain.TestConfig, useBundler bool) (string, []string) {
	args := []string{}

	if useBundler {
		args = append(args, "exec", "ruby")
	}

	args = append(args, testPath)

	if useBundler {
		return "bundle", args
	}
	return "ruby", args
}

// runRSpecBatch runs all tests using RSpec
func (r *RubyTestRunner) runRSpecBatch(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	startTime := time.Now()

	useBundler := r.hasBundler(suite.ProjectPath)
	cmdName := "rspec"
	args := []string{}

	if useBundler {
		cmdName = "bundle"
		args = append(args, "exec", "rspec")
	}

	// Add JSON format
	args = append(args, "--format", "json")

	// Add all test paths
	for _, test := range suite.Tests {
		args = append(args, test.Path)
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Parse JSON output
	results := r.parseRSpecJSON(string(output), duration)

	if len(results) == 0 && err != nil {
		return nil, fmt.Errorf("RSpec batch execution failed: %w", err)
	}

	return results, nil
}

// runMinitestBatch runs all tests using Minitest
func (r *RubyTestRunner) runMinitestBatch(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	startTime := time.Now()

	useBundler := r.hasBundler(suite.ProjectPath)
	cmdName := "ruby"
	args := []string{}

	if useBundler {
		cmdName = "bundle"
		args = append(args, "exec", "ruby")
	}

	// Run rake test if available
	if r.hasRakeTask(suite.ProjectPath, "test") {
		if useBundler {
			args = []string{"exec", "rake", "test"}
		} else {
			cmdName = "rake"
			args = []string{"test"}
		}
	} else {
		// Run individual test files
		args = append(args, "-e", "")
		for _, test := range suite.Tests {
			args[len(args)-1] += fmt.Sprintf("require './%s'; ", test.Path)
		}
	}

	cmd := exec.CommandContext(ctx, cmdName, args...)
	executil.HideWindow(cmd)
	cmd.Dir = suite.ProjectPath
	cmd.Env = os.Environ()

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	// Parse Minitest output
	results := r.parseMinitestOutput(string(output), suite.Tests, duration)

	if len(results) == 0 && err != nil {
		return nil, fmt.Errorf("Minitest batch execution failed: %w", err)
	}

	return results, nil
}

// parseRSpecJSON parses RSpec JSON output
func (r *RubyTestRunner) parseRSpecJSON(output string, totalDuration float64) []*domain.TestResult {
	// Find JSON in output
	jsonStart := strings.Index(output, "{\"version\":")
	if jsonStart == -1 {
		jsonStart = strings.Index(output, "{\"examples\":")
	}
	if jsonStart == -1 {
		return nil
	}

	jsonStr := output[jsonStart:]

	var report struct {
		Examples []struct {
			ID          string  `json:"id"`
			Description string  `json:"description"`
			FullDesc    string  `json:"full_description"`
			Status      string  `json:"status"`
			FilePath    string  `json:"file_path"`
			LineNumber  int     `json:"line_number"`
			RunTime     float64 `json:"run_time"`
			Exception   *struct {
				Class     string   `json:"class"`
				Message   string   `json:"message"`
				Backtrace []string `json:"backtrace"`
			} `json:"exception,omitempty"`
		} `json:"examples"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &report); err != nil {
		return nil
	}

	var results []*domain.TestResult
	for _, example := range report.Examples {
		result := &domain.TestResult{
			TestPath: example.FilePath,
			TestName: example.FullDesc,
			Language: "ruby",
			Duration: example.RunTime,
			Success:  example.Status == "passed",
		}

		if example.Exception != nil {
			result.Error = example.Exception.Class + ": " + example.Exception.Message
			if len(example.Exception.Backtrace) > 0 {
				result.Output = strings.Join(example.Exception.Backtrace, "\n")
			}
		}

		results = append(results, result)
	}

	return results
}

// parseMinitestOutput parses Minitest text output
func (r *RubyTestRunner) parseMinitestOutput(output string, tests []*domain.TestInfo, totalDuration float64) []*domain.TestResult {
	var results []*domain.TestResult

	// Pattern for Minitest summary: "10 runs, 20 assertions, 1 failures, 0 errors"
	summaryPattern := regexp.MustCompile(`(\d+)\s+runs?,\s+(\d+)\s+assertions?,\s+(\d+)\s+failures?,\s+(\d+)\s+errors?`)

	// Check if all tests passed
	allPassed := false
	if matches := summaryPattern.FindStringSubmatch(output); matches != nil {
		failures := matches[3]
		errors := matches[4]
		allPassed = failures == "0" && errors == "0"
	}

	avgDuration := totalDuration / float64(len(tests))

	for _, test := range tests {
		result := &domain.TestResult{
			TestPath: test.Path,
			TestName: test.Name,
			Language: "ruby",
			Duration: avgDuration,
			Output:   output,
			Success:  allPassed,
		}

		if !allPassed {
			result.Error = "Test failed"
		}

		results = append(results, result)
	}

	return results
}

// Helper methods

func (r *RubyTestRunner) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		"vendor",
		"node_modules",
		".git",
		".bundle",
		"tmp",
		"log",
		"coverage",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}
	return false
}

func (r *RubyTestRunner) isTestFile(name string) bool {
	if !strings.HasSuffix(name, ".rb") {
		return false
	}
	return strings.HasSuffix(name, rubyTestFileSuffix) ||
		strings.HasSuffix(name, rubyMinitestSuffix)
}

func (r *RubyTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit"
	}

	contentStr := strings.ToLower(string(content))
	fileName := strings.ToLower(filepath.Base(filePath))

	// Check for integration tests
	if strings.Contains(contentStr, "integration") || strings.Contains(fileName, "integration") {
		return "integration"
	}

	// Check for feature/acceptance tests
	if strings.Contains(contentStr, "feature") || strings.Contains(fileName, "feature") ||
		strings.Contains(contentStr, "capybara") {
		return "feature"
	}

	// Check for request/controller tests
	if strings.Contains(fileName, "request") || strings.Contains(fileName, "controller") {
		return "request"
	}

	return "unit"
}

func (r *RubyTestRunner) hasBundler(projectPath string) bool {
	gemfilePath := filepath.Join(projectPath, gemfileFile)
	_, err := os.Stat(gemfilePath)
	return err == nil
}

func (r *RubyTestRunner) isRSpecAvailable(projectPath string) bool {
	useBundler := r.hasBundler(projectPath)

	var cmd *exec.Cmd
	if useBundler {
		cmd = exec.Command("bundle", "exec", "rspec", "--version")
	} else {
		cmd = exec.Command("rspec", "--version")
	}
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	return cmd.Run() == nil
}

func (r *RubyTestRunner) isMinitestAvailable(projectPath string) bool {
	useBundler := r.hasBundler(projectPath)

	var cmd *exec.Cmd
	if useBundler {
		cmd = exec.Command("bundle", "exec", "ruby", "-e", "require 'minitest'")
	} else {
		cmd = exec.Command("ruby", "-e", "require 'minitest'")
	}
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	return cmd.Run() == nil
}

func (r *RubyTestRunner) hasRakeTask(projectPath, taskName string) bool {
	useBundler := r.hasBundler(projectPath)

	var cmd *exec.Cmd
	if useBundler {
		cmd = exec.Command("bundle", "exec", "rake", "-T", taskName)
	} else {
		cmd = exec.Command("rake", "-T", taskName)
	}
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.Contains(string(output), taskName)
}
