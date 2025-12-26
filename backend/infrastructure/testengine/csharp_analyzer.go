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

// CSharpTestAnalyzer implements domain.TestAnalyzer for C#
type CSharpTestAnalyzer struct {
	log domain.Logger
}

// NewCSharpTestAnalyzer creates a new C# test analyzer
func NewCSharpTestAnalyzer(log domain.Logger) *CSharpTestAnalyzer {
	return &CSharpTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *CSharpTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for C# test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractUsings(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *CSharpTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding C# tests for file: %s", filePath))

	var testFiles []string

	// Get class name from file
	className := strings.TrimSuffix(filepath.Base(filePath), ".cs")

	// Look for test files with standard naming conventions
	testPatterns := []string{
		className + "Test.cs",
		className + "Tests.cs",
		className + "TestCase.cs",
		className + "Spec.cs",
		className + "Specs.cs",
	}

	// Walk through project to find test files
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Skip bin, obj, .git directories
			if info.Name() == "bin" || info.Name() == "obj" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file matches test patterns
		for _, pattern := range testPatterns {
			if info.Name() == pattern {
				relPath, err := filepath.Rel(projectPath, path)
				if err != nil {
					return nil
				}
				testFiles = append(testFiles, relPath)
				return nil
			}
		}

		return nil
	})

	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to walk project directory: %v", err))
	}

	// Also search for tests that reference this class
	importTests, err := a.findTestsReferencingClass(ctx, className, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests referencing class: %v", err))
	} else {
		testFiles = append(testFiles, importTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// FindTestsForSymbol finds tests for a specific symbol (class, method)
func (a *CSharpTestAnalyzer) FindTestsForSymbol(ctx context.Context, symbolName, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding C# tests for symbol: %s", symbolName))

	var testFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if info.Name() == "bin" || info.Name() == "obj" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".cs") {
			return nil
		}

		// Check if file is a test file and references the symbol
		if a.isTestFile(path) && a.referencesSymbol(path, symbolName) {
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				return nil
			}
			testFiles = append(testFiles, relPath)
		}

		return nil
	})

	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to walk project directory: %v", err))
	}

	a.log.Info(fmt.Sprintf("Found %d test files for symbol %s", len(testFiles), symbolName))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *CSharpTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"[category(\"smoke\")]",
		"[trait(\"category\", \"smoke\")]",
		"[testcategory(\"smoke\")]",
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

// AnalyzeTestFile analyzes a C# test file and returns test information
func (a *CSharpTestAnalyzer) AnalyzeTestFile(ctx context.Context, testPath string) (*domain.TestInfo, error) {
	a.log.Info(fmt.Sprintf("Analyzing C# test file: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := string(content)
	fileName := filepath.Base(testPath)

	testInfo := &domain.TestInfo{
		Path:     testPath,
		Name:     fileName,
		Type:     a.determineTestType(contentStr, fileName),
		Metadata: make(map[string]string),
	}

	// Extract test framework
	testInfo.Metadata["framework"] = a.detectTestFramework(contentStr)

	// Extract class name
	testInfo.Metadata["className"] = a.extractClassName(contentStr)

	// Count test methods
	testInfo.Metadata["testCount"] = fmt.Sprintf("%d", a.countTestMethods(contentStr))

	return testInfo, nil
}

// extractUsings extracts using statements from C# code
func (a *CSharpTestAnalyzer) extractUsings(content string) []string {
	var usings []string
	usingSet := make(map[string]bool)

	// Match using statements
	usingRe := regexp.MustCompile(`using\s+([\w.]+)\s*;`)
	matches := usingRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			usingPath := match[1]
			if !usingSet[usingPath] {
				usings = append(usings, usingPath)
				usingSet[usingPath] = true
			}
		}
	}

	return usings
}

// findTestsReferencingClass finds tests that reference a given class
func (a *CSharpTestAnalyzer) findTestsReferencingClass(ctx context.Context, className, projectPath string) ([]string, error) {
	var testFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if info.Name() == "bin" || info.Name() == "obj" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".cs") {
			return nil
		}

		// Check if file is a test file
		if !a.isTestFile(path) {
			return nil
		}

		// Check if file references the class
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if a.referencesClass(string(content), className) {
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				return nil
			}
			testFiles = append(testFiles, relPath)
		}

		return nil
	})

	return testFiles, err
}

// isTestFile checks if a file is a test file
func (a *CSharpTestAnalyzer) isTestFile(filePath string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	contentStr := string(content)
	testIndicators := []string{
		"[Test]",
		"[TestMethod]",
		"[Fact]",
		"[Theory]",
		"[TestCase]",
	}

	for _, indicator := range testIndicators {
		if strings.Contains(contentStr, indicator) {
			return true
		}
	}

	return false
}

// referencesClass checks if content references a specific class
func (a *CSharpTestAnalyzer) referencesClass(content, className string) bool {
	// Check for direct usage
	patterns := []string{
		fmt.Sprintf(`new\s+%s\s*\(`, className),
		fmt.Sprintf(`%s\.`, className),
		fmt.Sprintf(`<%s>`, className),
		fmt.Sprintf(`\(%s\s+`, className),
	}

	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, content)
		if matched {
			return true
		}
	}

	return false
}

// referencesSymbol checks if a file references a specific symbol
func (a *CSharpTestAnalyzer) referencesSymbol(filePath, symbolName string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), symbolName)
}

// determineTestType determines the type of test based on content and filename
func (a *CSharpTestAnalyzer) determineTestType(content, fileName string) string {
	contentLower := strings.ToLower(content)
	fileNameLower := strings.ToLower(fileName)

	// Check for smoke tests
	if strings.Contains(contentLower, "smoke") || strings.Contains(fileNameLower, "smoke") {
		return "smoke"
	}

	// Check for integration tests
	integrationIndicators := []string{
		"integration",
		"webapplicationfactory",
		"testserver",
		"httpclient",
		"[collection(",
		"testcontainers",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentLower, indicator) {
			return "integration"
		}
	}

	return "unit"
}

// detectTestFramework detects the test framework used
func (a *CSharpTestAnalyzer) detectTestFramework(content string) string {
	// xUnit
	if strings.Contains(content, "[Fact]") || strings.Contains(content, "[Theory]") {
		return "xunit"
	}

	// NUnit
	if strings.Contains(content, "[Test]") || strings.Contains(content, "[TestCase]") {
		return "nunit"
	}

	// MSTest
	if strings.Contains(content, "[TestMethod]") || strings.Contains(content, "[TestClass]") {
		return "mstest"
	}

	return "unknown"
}

// extractClassName extracts the class name from content
func (a *CSharpTestAnalyzer) extractClassName(content string) string {
	classRe := regexp.MustCompile(`class\s+(\w+)`)
	matches := classRe.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// countTestMethods counts the number of test methods in the content
func (a *CSharpTestAnalyzer) countTestMethods(content string) int {
	count := 0

	testAttributes := []string{
		`\[Test\]`,
		`\[TestMethod\]`,
		`\[Fact\]`,
		`\[Theory\]`,
		`\[TestCase\(`,
	}

	for _, attr := range testAttributes {
		re := regexp.MustCompile(attr)
		matches := re.FindAllString(content, -1)
		count += len(matches)
	}

	return count
}

// removeDuplicates removes duplicate strings from slice
func (a *CSharpTestAnalyzer) removeDuplicates(items []string) []string {
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
