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

// RustTestAnalyzer implements TestAnalyzer for Rust
type RustTestAnalyzer struct {
log domain.Logger
}

// NewRustTestAnalyzer creates a new analyzer for Rust tests
func NewRustTestAnalyzer(log domain.Logger) *RustTestAnalyzer {
return &RustTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *RustTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
a.log.Info(fmt.Sprintf("Analyzing dependencies for Rust test: %s", testPath))

content, err := os.ReadFile(testPath)
if err != nil {
return nil, fmt.Errorf("failed to read test file: %w", err)
}

var dependencies []string
contentStr := string(content)

// Find use statements
useRegex := regexp.MustCompile("use\\s+([^;]+);")
matches := useRegex.FindAllStringSubmatch(contentStr, -1)
for _, match := range matches {
if len(match) > 1 {
dep := strings.TrimSpace(match[1])
if !strings.HasPrefix(dep, "std::") && !strings.HasPrefix(dep, "core::") {
dependencies = append(dependencies, dep)
}
}
}

// Find extern crate
externRegex := regexp.MustCompile("extern\\s+crate\\s+(\\w+);")
externMatches := externRegex.FindAllStringSubmatch(contentStr, -1)
for _, match := range externMatches {
if len(match) > 1 {
dependencies = append(dependencies, match[1])
}
}

a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
return dependencies, nil
}

// FindTestsForFile finds tests for a file
func (a *RustTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
a.log.Info(fmt.Sprintf("Finding tests for Rust file: %s", filePath))

var testFiles []string
moduleName := strings.TrimSuffix(filepath.Base(filePath), ".rs")

// Check if file has inline tests
fullPath := filepath.Join(projectPath, filePath)
if hasTests, _ := a.fileHasTests(fullPath); hasTests {
testFiles = append(testFiles, filePath)
}

// Look for test files in tests/
testsDir := filepath.Join(projectPath, "tests")
if _, err := os.Stat(testsDir); err == nil {
testPatterns := []string{
moduleName + ".rs",
moduleName + "_test.rs",
"test_" + moduleName + ".rs",
}
for _, pattern := range testPatterns {
testFile := filepath.Join("tests", pattern)
fullTestPath := filepath.Join(projectPath, testFile)
if _, err := os.Stat(fullTestPath); err == nil {
testFiles = append(testFiles, testFile)
}
}
}

// Remove duplicates
uniqueTests := make(map[string]bool)
var uniqueTestFiles []string
for _, test := range testFiles {
if !uniqueTests[test] {
uniqueTests[test] = true
uniqueTestFiles = append(uniqueTestFiles, test)
}
}

a.log.Info(fmt.Sprintf("Found %d test files for %s", len(uniqueTestFiles), filePath))
return uniqueTestFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *RustTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
content, err := os.ReadFile(testPath)
if err != nil {
return false, fmt.Errorf("failed to read test file: %w", err)
}

contentStr := strings.ToLower(string(content))
smokeIndicators := []string{"smoke", "#[ignore]", "// smoke", "/* smoke"}

for _, indicator := range smokeIndicators {
if strings.Contains(contentStr, indicator) {
return true, nil
}
}

fileName := strings.ToLower(filepath.Base(testPath))
return strings.Contains(fileName, "smoke"), nil
}

// fileHasTests checks if file contains tests
func (a *RustTestAnalyzer) fileHasTests(filePath string) (bool, error) {
content, err := os.ReadFile(filePath)
if err != nil {
return false, err
}

contentStr := string(content)
testPatterns := []string{"#[test]", "#[cfg(test)]", "mod tests"}

for _, pattern := range testPatterns {
if strings.Contains(contentStr, pattern) {
return true, nil
}
}

return false, nil
}
