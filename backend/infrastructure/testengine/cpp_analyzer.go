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

// C++ analyzer constants
const (
	cppIncludePattern       = `^\s*#include\s*[<"]([^>"]+)[>"]`
	cppTestFuncGTest        = `TEST(?:_F|_P)?\s*\(\s*(\w+)\s*,\s*(\w+)\s*\)`
	cppTestFuncCatch2       = `TEST_CASE\s*\(\s*"([^"]+)"`
	cppTestFuncBoost        = `BOOST_AUTO_TEST_CASE\s*\(\s*(\w+)\s*\)`
	cppTestSuiteGTest       = `TEST_F\s*\(\s*(\w+)\s*,`
	cppTestSuiteCatch2      = `TEST_CASE\s*\([^,]+,\s*"([^"]+)"`
)

// CppTestAnalyzer implements domain.TestAnalyzer for C++
type CppTestAnalyzer struct {
	log domain.Logger
}

// NewCppTestAnalyzer creates a new C++ test analyzer
func NewCppTestAnalyzer(log domain.Logger) *CppTestAnalyzer {
	return &CppTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *CppTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for C++ test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractIncludes(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *CppTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding C++ tests for file: %s", filePath))

	var testFiles []string

	// Get file name without extension
	ext := filepath.Ext(filePath)
	fileName := strings.TrimSuffix(filepath.Base(filePath), ext)

	// Strategy 1: Look for <filename>_test.cpp in test/ or tests/ directory
	testDirs := []string{
		filepath.Join(projectPath, "test"),
		filepath.Join(projectPath, "tests"),
	}

	testExtensions := []string{".cpp", ".cc", ".cxx"}

	for _, testDir := range testDirs {
		for _, testExt := range testExtensions {
			testFile := filepath.Join(testDir, fileName+"_test"+testExt)
			if a.fileExists(testFile) {
				relPath, _ := filepath.Rel(projectPath, testFile)
				testFiles = append(testFiles, relPath)
			}

			// Also check test_<filename>
			testFile2 := filepath.Join(testDir, "test_"+fileName+testExt)
			if a.fileExists(testFile2) {
				relPath, _ := filepath.Rel(projectPath, testFile2)
				testFiles = append(testFiles, relPath)
			}
		}
	}

	// Strategy 2: Search for tests that include this file
	includeTests, err := a.findTestsIncludingFile(ctx, filePath, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests including file: %v", err))
	} else {
		testFiles = append(testFiles, includeTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *CppTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"// smoke",
		"/* smoke",
		"smoke_test",
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

// extractIncludes extracts #include statements from C++ code
func (a *CppTestAnalyzer) extractIncludes(content string) []string {
	var includes []string
	includeSet := make(map[string]bool)

	lines := strings.Split(content, "\n")
	includeRe := regexp.MustCompile(cppIncludePattern)

	for _, line := range lines {
		// Check for include pattern
		if matches := includeRe.FindStringSubmatch(line); len(matches) > 1 {
			includePath := matches[1]
			if !includeSet[includePath] {
				includes = append(includes, includePath)
				includeSet[includePath] = true
			}
		}
	}

	return includes
}

// findTestsIncludingFile finds tests that include a given file
func (a *CppTestAnalyzer) findTestsIncludingFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	var testFiles []string

	fileName := filepath.Base(filePath)

	testDirs := []string{
		filepath.Join(projectPath, "test"),
		filepath.Join(projectPath, "tests"),
	}

	for _, testDir := range testDirs {
		if _, err := os.Stat(testDir); os.IsNotExist(err) {
			continue
		}

		err := filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				if a.shouldSkipDirectory(info.Name()) {
					return filepath.SkipDir
				}
				return nil
			}

			// Only check C++ test files
			if !a.isTestFile(info.Name()) {
				return nil
			}

			// Check if test includes the file
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if a.includesFile(string(content), fileName) {
				relPath, err := filepath.Rel(projectPath, path)
				if err != nil {
					return nil
				}
				testFiles = append(testFiles, relPath)
			}

			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed to search for tests: %w", err)
		}
	}

	return testFiles, nil
}

// includesFile checks if content includes a specific file
func (a *CppTestAnalyzer) includesFile(content, fileName string) bool {
	// Check for direct include
	patterns := []string{
		fmt.Sprintf(`#include\s*[<"]%s[>"]`, regexp.QuoteMeta(fileName)),
		fmt.Sprintf(`#include\s*[<"][^>"]*/%s[>"]`, regexp.QuoteMeta(fileName)),
	}

	for _, pattern := range patterns {
		if matched, _ := regexp.MatchString(pattern, content); matched {
			return true
		}
	}

	return false
}

// fileExists checks if a file exists
func (a *CppTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// shouldSkipDirectory checks if directory should be skipped
func (a *CppTestAnalyzer) shouldSkipDirectory(name string) bool {
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
func (a *CppTestAnalyzer) isTestFile(name string) bool {
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

// removeDuplicates removes duplicate strings from slice
func (a *CppTestAnalyzer) removeDuplicates(items []string) []string {
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

// ExtractTestFunctions extracts test function names from a file
func (a *CppTestAnalyzer) ExtractTestFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var functions []string
	contentStr := string(content)

	// Google Test pattern: TEST(Suite, Name) or TEST_F(Suite, Name)
	gtestRe := regexp.MustCompile(cppTestFuncGTest)
	gtestMatches := gtestRe.FindAllStringSubmatch(contentStr, -1)
	for _, match := range gtestMatches {
		if len(match) > 2 {
			functions = append(functions, fmt.Sprintf("%s.%s", match[1], match[2]))
		}
	}

	// Catch2 pattern: TEST_CASE("name")
	catch2Re := regexp.MustCompile(cppTestFuncCatch2)
	catch2Matches := catch2Re.FindAllStringSubmatch(contentStr, -1)
	for _, match := range catch2Matches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	// Boost.Test pattern: BOOST_AUTO_TEST_CASE(name)
	boostRe := regexp.MustCompile(cppTestFuncBoost)
	boostMatches := boostRe.FindAllStringSubmatch(contentStr, -1)
	for _, match := range boostMatches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	return functions, nil
}

// ExtractTestSuites extracts test suite names from a file
func (a *CppTestAnalyzer) ExtractTestSuites(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var suites []string
	suiteSet := make(map[string]bool)
	contentStr := string(content)

	// Google Test fixture pattern
	gtestRe := regexp.MustCompile(cppTestSuiteGTest)
	gtestMatches := gtestRe.FindAllStringSubmatch(contentStr, -1)
	for _, match := range gtestMatches {
		if len(match) > 1 && !suiteSet[match[1]] {
			suites = append(suites, match[1])
			suiteSet[match[1]] = true
		}
	}

	// Catch2 tag pattern
	catch2Re := regexp.MustCompile(cppTestSuiteCatch2)
	catch2Matches := catch2Re.FindAllStringSubmatch(contentStr, -1)
	for _, match := range catch2Matches {
		if len(match) > 1 && !suiteSet[match[1]] {
			suites = append(suites, match[1])
			suiteSet[match[1]] = true
		}
	}

	return suites, nil
}

// DetectTestFramework detects which test framework is used in a file
func (a *CppTestAnalyzer) DetectTestFramework(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unknown"
	}

	contentStr := string(content)

	// Check for Google Test
	if strings.Contains(contentStr, "gtest/gtest.h") ||
		strings.Contains(contentStr, "TEST(") ||
		strings.Contains(contentStr, "TEST_F(") {
		return "googletest"
	}

	// Check for Catch2
	if strings.Contains(contentStr, "catch2/catch") ||
		strings.Contains(contentStr, "catch.hpp") ||
		strings.Contains(contentStr, "TEST_CASE(") {
		return "catch2"
	}

	// Check for Boost.Test
	if strings.Contains(contentStr, "boost/test") ||
		strings.Contains(contentStr, "BOOST_AUTO_TEST_CASE") {
		return "boost"
	}

	return "unknown"
}
