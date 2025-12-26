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

// Dart analyzer constants
const (
	dartImportPattern     = `^import\s+['"]([^'"]+)['"]`
	dartTestFuncPattern   = `\btest\s*\(\s*['"]([^'"]+)['"]`
	dartTestGroupPattern  = `\bgroup\s*\(\s*['"]([^'"]+)['"]`
	dartTestWidgetPattern = `\btestWidgets\s*\(\s*['"]([^'"]+)['"]`
)

// DartTestAnalyzer implements domain.TestAnalyzer for Dart
type DartTestAnalyzer struct {
	log domain.Logger
}

// NewDartTestAnalyzer creates a new Dart test analyzer
func NewDartTestAnalyzer(log domain.Logger) *DartTestAnalyzer {
	return &DartTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *DartTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for Dart test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractImports(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *DartTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding Dart tests for file: %s", filePath))

	var testFiles []string

	// Get file name without extension
	fileName := strings.TrimSuffix(filepath.Base(filePath), ".dart")
	dir := filepath.Dir(filePath)

	// Strategy 1: Look for <filename>_test.dart in test/ directory
	testDir := filepath.Join(projectPath, "test")
	testFile1 := filepath.Join(testDir, fileName+"_test.dart")
	if a.fileExists(testFile1) {
		relPath, _ := filepath.Rel(projectPath, testFile1)
		testFiles = append(testFiles, relPath)
	}

	// Strategy 2: Mirror directory structure in test/
	relDir, err := filepath.Rel(filepath.Join(projectPath, "lib"), dir)
	if err == nil && relDir != "." {
		testFile2 := filepath.Join(testDir, relDir, fileName+"_test.dart")
		if a.fileExists(testFile2) {
			relPath, _ := filepath.Rel(projectPath, testFile2)
			testFiles = append(testFiles, relPath)
		}
	}

	// Strategy 3: Search for tests that import this file
	importTests, err := a.findTestsImportingFile(ctx, filePath, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests importing file: %v", err))
	} else {
		testFiles = append(testFiles, importTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *DartTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"@smoke",
		"// smoke",
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

// extractImports extracts import statements from Dart code
func (a *DartTestAnalyzer) extractImports(content string) []string {
	var imports []string
	importSet := make(map[string]bool)

	lines := strings.Split(content, "\n")
	importRe := regexp.MustCompile(dartImportPattern)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Check for import pattern
		if matches := importRe.FindStringSubmatch(line); len(matches) > 1 {
			importPath := matches[1]
			if !importSet[importPath] {
				imports = append(imports, importPath)
				importSet[importPath] = true
			}
		}
	}

	return imports
}

// findTestsImportingFile finds tests that import a given file
func (a *DartTestAnalyzer) findTestsImportingFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	var testFiles []string

	// Convert file path to package import path
	importPath := a.filePathToImportPath(filePath, projectPath)
	if importPath == "" {
		return nil, nil
	}

	testDir := filepath.Join(projectPath, "test")
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		return nil, nil
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

		// Only check test files
		if !strings.HasSuffix(info.Name(), "_test.dart") {
			return nil
		}

		// Check if test imports the file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if a.importsFile(string(content), importPath, filePath) {
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

	return testFiles, nil
}

// importsFile checks if content imports a specific file
func (a *DartTestAnalyzer) importsFile(content, importPath, filePath string) bool {
	// Check for package import
	if importPath != "" && strings.Contains(content, importPath) {
		return true
	}

	// Check for relative import
	fileName := filepath.Base(filePath)
	if strings.Contains(content, fileName) {
		return true
	}

	return false
}

// filePathToImportPath converts file path to Dart package import path
func (a *DartTestAnalyzer) filePathToImportPath(filePath, projectPath string) string {
	// Read pubspec.yaml to get package name
	pubspecPath := filepath.Join(projectPath, "pubspec.yaml")
	content, err := os.ReadFile(pubspecPath)
	if err != nil {
		return ""
	}

	// Extract package name
	packageName := a.extractPackageName(string(content))
	if packageName == "" {
		return ""
	}

	// Convert file path to import path
	libDir := filepath.Join(projectPath, "lib")
	relPath, err := filepath.Rel(libDir, filePath)
	if err != nil {
		return ""
	}

	// Convert to package import format
	relPath = strings.ReplaceAll(relPath, string(filepath.Separator), "/")
	return fmt.Sprintf("package:%s/%s", packageName, relPath)
}

// extractPackageName extracts package name from pubspec.yaml content
func (a *DartTestAnalyzer) extractPackageName(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "name:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

// fileExists checks if a file exists
func (a *DartTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// shouldSkipDirectory checks if directory should be skipped
func (a *DartTestAnalyzer) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		"build",
		".dart_tool",
		".packages",
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

// removeDuplicates removes duplicate strings from slice
func (a *DartTestAnalyzer) removeDuplicates(items []string) []string {
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
func (a *DartTestAnalyzer) ExtractTestFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var functions []string
	testRe := regexp.MustCompile(dartTestFuncPattern)
	widgetRe := regexp.MustCompile(dartTestWidgetPattern)

	// Find test() calls
	matches := testRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	// Find testWidgets() calls
	widgetMatches := widgetRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range widgetMatches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	return functions, nil
}

// ExtractTestGroups extracts test group names from a file
func (a *DartTestAnalyzer) ExtractTestGroups(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var groups []string
	groupRe := regexp.MustCompile(dartTestGroupPattern)

	matches := groupRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			groups = append(groups, match[1])
		}
	}

	return groups, nil
}
