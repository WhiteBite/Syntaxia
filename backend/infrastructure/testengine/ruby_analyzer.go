package testengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// RubyTestAnalyzer analyzes Ruby test output
type RubyTestAnalyzer struct {
	log domain.Logger
}

// NewRubyTestAnalyzer creates a new Ruby test analyzer
func NewRubyTestAnalyzer(log domain.Logger) *RubyTestAnalyzer {
	return &RubyTestAnalyzer{log: log}
}

// GetLanguage returns the language identifier
func (a *RubyTestAnalyzer) GetLanguage() string {
	return "ruby"
}

// AnalyzeOutput analyzes Ruby test output and extracts test results
func (a *RubyTestAnalyzer) AnalyzeOutput(output string) ([]*domain.TestResult, error) {
	// Try RSpec JSON format first
	if results := a.parseRSpecJSON(output); len(results) > 0 {
		return results, nil
	}

	// Try RSpec text format
	if results := a.parseRSpecText(output); len(results) > 0 {
		return results, nil
	}

	// Try Minitest format
	if results := a.parseMinitestOutput(output); len(results) > 0 {
		return results, nil
	}

	// Fallback: create single result from output
	return a.createFallbackResult(output), nil
}

// parseRSpecJSON parses RSpec JSON output
func (a *RubyTestAnalyzer) parseRSpecJSON(output string) []*domain.TestResult {
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
		Summary struct {
			Duration     float64 `json:"duration"`
			ExampleCount int     `json:"example_count"`
			FailureCount int     `json:"failure_count"`
			PendingCount int     `json:"pending_count"`
		} `json:"summary"`
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

// parseRSpecText parses RSpec text output
func (a *RubyTestAnalyzer) parseRSpecText(output string) []*domain.TestResult {
	var results []*domain.TestResult

	// Pattern for RSpec failure: "rspec ./spec/file_spec.rb:10 # description"
	failurePattern := regexp.MustCompile(`rspec\s+(\./)?([^\s:]+):(\d+)\s+#\s+(.+)`)

	// Pattern for summary: "10 examples, 2 failures"
	summaryPattern := regexp.MustCompile(`(\d+)\s+examples?,\s+(\d+)\s+failures?`)

	lines := strings.Split(output, "\n")

	// Extract failures
	for _, line := range lines {
		if matches := failurePattern.FindStringSubmatch(line); matches != nil {
			result := &domain.TestResult{
				TestPath: matches[2],
				TestName: matches[4],
				Language: "ruby",
				Success:  false,
				Error:    "Test failed",
			}
			results = append(results, result)
		}
	}

	// If we found failures, return them
	if len(results) > 0 {
		return results
	}

	// Check summary for success
	if matches := summaryPattern.FindStringSubmatch(output); matches != nil {
		if matches[2] == "0" {
			// All tests passed - create a summary result
			result := &domain.TestResult{
				TestPath: "spec",
				TestName: "RSpec Suite",
				Language: "ruby",
				Success:  true,
				Output:   output,
			}
			return []*domain.TestResult{result}
		}
	}

	return nil
}

// parseMinitestOutput parses Minitest output
func (a *RubyTestAnalyzer) parseMinitestOutput(output string) []*domain.TestResult {
	var results []*domain.TestResult

	// Pattern for Minitest failure: "Failure: TestClass#test_method [file.rb:10]:"
	failurePattern := regexp.MustCompile(`(Failure|Error):\s+(\w+)#(\w+)\s+\[([^\]]+):(\d+)\]:`)

	// Pattern for summary: "10 runs, 20 assertions, 1 failures, 0 errors"
	summaryPattern := regexp.MustCompile(`(\d+)\s+runs?,\s+(\d+)\s+assertions?,\s+(\d+)\s+failures?,\s+(\d+)\s+errors?`)

	lines := strings.Split(output, "\n")

	// Extract failures
	for i, line := range lines {
		if matches := failurePattern.FindStringSubmatch(line); matches != nil {
			errorMsg := ""
			// Get error message from next lines
			if i+1 < len(lines) {
				errorMsg = strings.TrimSpace(lines[i+1])
			}

			result := &domain.TestResult{
				TestPath: matches[4],
				TestName: matches[2] + "#" + matches[3],
				Language: "ruby",
				Success:  false,
				Error:    errorMsg,
			}
			results = append(results, result)
		}
	}

	// If we found failures, return them
	if len(results) > 0 {
		return results
	}

	// Check summary for success
	if matches := summaryPattern.FindStringSubmatch(output); matches != nil {
		failures := matches[3]
		errors := matches[4]
		if failures == "0" && errors == "0" {
			// All tests passed
			result := &domain.TestResult{
				TestPath: "test",
				TestName: "Minitest Suite",
				Language: "ruby",
				Success:  true,
				Output:   output,
			}
			return []*domain.TestResult{result}
		}
	}

	return nil
}

// createFallbackResult creates a fallback result when parsing fails
func (a *RubyTestAnalyzer) createFallbackResult(output string) []*domain.TestResult {
	result := &domain.TestResult{
		TestPath: "unknown",
		TestName: "Ruby Tests",
		Language: "ruby",
		Output:   output,
	}

	// Try to determine success from output
	outputLower := strings.ToLower(output)
	if strings.Contains(outputLower, "0 failures") &&
		strings.Contains(outputLower, "0 errors") {
		result.Success = true
	} else if strings.Contains(outputLower, "failure") ||
		strings.Contains(outputLower, "error") ||
		strings.Contains(outputLower, "failed") {
		result.Success = false
		result.Error = "Test execution failed"
	} else {
		result.Success = true
	}

	return []*domain.TestResult{result}
}

// RubyTestFailure represents a Ruby test failure
type RubyTestFailure struct {
	TestName string
	Location string
	Message  string
}

// ExtractFailures extracts failure details from test output
func (a *RubyTestAnalyzer) ExtractFailures(output string) []*RubyTestFailure {
	var failures []*RubyTestFailure

	// RSpec failure pattern
	rspecFailure := regexp.MustCompile(`(?m)^\s*\d+\)\s+(.+)\n\s+Failure/Error:\s+(.+)\n`)

	// Minitest failure pattern
	minitestFailure := regexp.MustCompile(`(?m)(Failure|Error):\s+(\w+#\w+)\s+\[([^\]]+)\]:\n(.+)`)

	// Extract RSpec failures
	if matches := rspecFailure.FindAllStringSubmatch(output, -1); matches != nil {
		for _, match := range matches {
			failure := &RubyTestFailure{
				TestName: match[1],
				Message:  match[2],
			}
			failures = append(failures, failure)
		}
	}

	// Extract Minitest failures
	if matches := minitestFailure.FindAllStringSubmatch(output, -1); matches != nil {
		for _, match := range matches {
			failure := &RubyTestFailure{
				TestName: match[2],
				Location: match[3],
				Message:  match[4],
			}
			failures = append(failures, failure)
		}
	}

	return failures
}

// GetTestCoverage extracts coverage information from output
func (a *RubyTestAnalyzer) GetTestCoverage(output string) *domain.TestCoverage {
	coverage := &domain.TestCoverage{
		Percentage: 0,
	}

	// SimpleCov coverage pattern: "Coverage report generated... 85.5% covered"
	coveragePattern := regexp.MustCompile(`(\d+\.?\d*)%\s+covered`)

	if matches := coveragePattern.FindStringSubmatch(output); matches != nil {
		var percent float64
		if _, err := parseFloatSafe(matches[1], &percent); err == nil {
			coverage.Percentage = percent
		}
	}

	// Alternative pattern: "Line Coverage: 85.5%"
	lineCoveragePattern := regexp.MustCompile(`Line Coverage:\s*(\d+\.?\d*)%`)
	if matches := lineCoveragePattern.FindStringSubmatch(output); matches != nil {
		var percent float64
		if _, err := parseFloatSafe(matches[1], &percent); err == nil {
			coverage.Percentage = percent
		}
	}

	return coverage
}

// parseFloatSafe safely parses a float string
func parseFloatSafe(s string, result *float64) (bool, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return false, nil
	}

	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return false, err
	}

	*result = f
	return true, nil
}

// AnalyzeTestDependencies analyzes test dependencies for a Ruby test file
func (a *RubyTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for Ruby test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	return a.extractRequires(string(content)), nil
}

// FindTestsForFile finds tests for a given source file
func (a *RubyTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding Ruby tests for file: %s", filePath))

	var testFiles []string

	// Get base name without extension
	baseName := strings.TrimSuffix(filepath.Base(filePath), ".rb")

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
			if err != nil || info.IsDir() {
				return nil
			}

			fileName := info.Name()

			// Check for matching test files
			if strings.HasSuffix(fileName, "_spec.rb") || strings.HasSuffix(fileName, "_test.rb") {
				testBaseName := strings.TrimSuffix(strings.TrimSuffix(fileName, "_spec.rb"), "_test.rb")
				if testBaseName == baseName {
					relPath, _ := filepath.Rel(projectPath, path)
					testFiles = append(testFiles, relPath)
				}
			}

			return nil
		})

		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to walk test directory %s: %v", testDir, err))
		}
	}

	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *RubyTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))
	fileName := strings.ToLower(filepath.Base(testPath))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		":smoke",
		"tag: :smoke",
		"tags: [:smoke",
	}

	for _, indicator := range smokeIndicators {
		if strings.Contains(contentStr, indicator) || strings.Contains(fileName, indicator) {
			return true, nil
		}
	}

	return false, nil
}

// extractRequires extracts require statements from Ruby code
func (a *RubyTestAnalyzer) extractRequires(content string) []string {
	var requires []string
	seen := make(map[string]bool)

	// Match require and require_relative statements
	requireRe := regexp.MustCompile(`require(?:_relative)?\s+['"]([^'"]+)['"]`)
	matches := requireRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			req := match[1]
			if !seen[req] {
				requires = append(requires, req)
				seen[req] = true
			}
		}
	}

	return requires
}
