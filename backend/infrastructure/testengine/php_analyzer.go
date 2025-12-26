package testengine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// PHP analyzer constants
const (
	phpTestAnnotation  = "@test"
	phpTestMethodRegex = `function\s+(test\w+)\s*\(`
)

// PHPTestAnalyzer implements domain.TestAnalyzer for PHP
type PHPTestAnalyzer struct {
	log domain.Logger
}

// NewPHPTestAnalyzer creates a new PHP test analyzer
func NewPHPTestAnalyzer(log domain.Logger) *PHPTestAnalyzer {
	return &PHPTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *PHPTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for PHP test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractUseStatements(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *PHPTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding PHP tests for file: %s", filePath))

	var testFiles []string

	// Get class name from file
	className := strings.TrimSuffix(filepath.Base(filePath), ".php")

	// Standard test directories
	testDirs := []string{
		"tests",
		"test",
		"Tests",
		"Test",
		"tests/Unit",
		"tests/Feature",
		"tests/Integration",
	}

	// Look for test files in standard locations
	for _, testDir := range testDirs {
		// Strategy 1: Same name with Test suffix
		testFile := filepath.Join(testDir, className+"Test.php")
		if a.fileExists(filepath.Join(projectPath, testFile)) {
			testFiles = append(testFiles, testFile)
		}

		// Strategy 2: Look in subdirectories
		fullTestDir := filepath.Join(projectPath, testDir)
		if _, err := os.Stat(fullTestDir); err == nil {
			_ = filepath.Walk(fullTestDir, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				if info.Name() == className+"Test.php" {
					relPath, _ := filepath.Rel(projectPath, path)
					testFiles = append(testFiles, relPath)
				}
				return nil
			})
		}
	}

	// Strategy 3: Search for tests that use this class
	useTests, err := a.findTestsUsingClass(ctx, className, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests using class: %v", err))
	} else {
		testFiles = append(testFiles, useTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *PHPTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"@group smoke",
		"smoketest",
	}

	for _, indicator := range smokeIndicators {
		if strings.Contains(contentStr, indicator) {
			return true, nil
		}
	}

	// Check file name
	fileName := strings.ToLower(filepath.Base(testPath))
	if strings.Contains(fileName, "smoke") {
		return true, nil
	}

	return false, nil
}

// extractUseStatements extracts use statements from PHP code
func (a *PHPTestAnalyzer) extractUseStatements(content string) []string {
	var uses []string
	useSet := make(map[string]bool)

	// Match use statements
	useRe := regexp.MustCompile(`use\s+([a-zA-Z_\\][\w\\]*);`)
	matches := useRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			usePath := match[1]
			if !useSet[usePath] {
				uses = append(uses, usePath)
				useSet[usePath] = true
			}
		}
	}

	return uses
}

// findTestsUsingClass finds tests that use a given class
func (a *PHPTestAnalyzer) findTestsUsingClass(ctx context.Context, className, projectPath string) ([]string, error) {
	var testFiles []string

	testDirs := []string{
		"tests",
		"test",
		"Tests",
		"Test",
	}

	for _, testDir := range testDirs {
		fullTestDir := filepath.Join(projectPath, testDir)
		if _, err := os.Stat(fullTestDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(fullTestDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				if info.Name() == "vendor" || info.Name() == ".git" {
					return filepath.SkipDir
				}
				return nil
			}

			if !strings.HasSuffix(info.Name(), "Test.php") {
				return nil
			}

			// Check if file uses the class
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if a.usesClass(string(content), className) {
				relPath, err := filepath.Rel(projectPath, path)
				if err != nil {
					return nil
				}
				testFiles = append(testFiles, relPath)
			}

			return nil
		})

		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to walk test directory %s: %v", testDir, err))
		}
	}

	return testFiles, nil
}

// usesClass checks if content uses a specific class
func (a *PHPTestAnalyzer) usesClass(content, className string) bool {
	// Check for use statement
	pattern := fmt.Sprintf(`use\s+[\w\\]*\\%s;`, className)
	matched, _ := regexp.MatchString(pattern, content)
	if matched {
		return true
	}

	// Check for direct usage
	if strings.Contains(content, className) {
		return true
	}

	return false
}

// ExtractTestMethods extracts test method names from a PHP test file
func (a *PHPTestAnalyzer) ExtractTestMethods(content string) []string {
	var methods []string

	// Match test methods (methods starting with "test")
	methodRe := regexp.MustCompile(phpTestMethodRegex)
	matches := methodRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			methods = append(methods, match[1])
		}
	}

	// Also match methods with @test annotation
	annotationRe := regexp.MustCompile(`/\*\*[^*]*\*\s*@test[^*]*\*/\s*(?:public\s+)?function\s+(\w+)\s*\(`)
	annotationMatches := annotationRe.FindAllStringSubmatch(content, -1)

	for _, match := range annotationMatches {
		if len(match) > 1 {
			methods = append(methods, match[1])
		}
	}

	return a.removeDuplicates(methods)
}

// GetTestGroups extracts PHPUnit groups from a test file
func (a *PHPTestAnalyzer) GetTestGroups(content string) []string {
	var groups []string
	groupSet := make(map[string]bool)

	// Match @group annotations
	groupRe := regexp.MustCompile(`@group\s+(\w+)`)
	matches := groupRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			group := match[1]
			if !groupSet[group] {
				groups = append(groups, group)
				groupSet[group] = true
			}
		}
	}

	return groups
}

// fileExists checks if a file exists
func (a *PHPTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// removeDuplicates removes duplicate strings from slice
func (a *PHPTestAnalyzer) removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
