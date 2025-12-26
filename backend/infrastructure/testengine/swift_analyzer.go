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

// Swift analyzer constants
const (
	swiftImportPattern    = `^import\s+(\w+)`
	swiftTestFuncPattern  = `func\s+(test\w+)\s*\(`
	swiftTestClassPattern = `class\s+(\w+Tests?)\s*:\s*XCTestCase`
)

// SwiftTestAnalyzer implements domain.TestAnalyzer for Swift
type SwiftTestAnalyzer struct {
	log domain.Logger
}

// NewSwiftTestAnalyzer creates a new Swift test analyzer
func NewSwiftTestAnalyzer(log domain.Logger) *SwiftTestAnalyzer {
	return &SwiftTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *SwiftTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for Swift test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("swift: failed to read test file: %w", err)
	}

	dependencies := a.extractImports(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *SwiftTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding Swift tests for file: %s", filePath))

	var testFiles []string

	// Get file name without extension
	fileName := strings.TrimSuffix(filepath.Base(filePath), ".swift")
	dir := filepath.Dir(filePath)

	// Strategy 1: Look for <filename>Tests.swift in Tests/ directory
	testDir := filepath.Join(projectPath, "Tests")
	testFile1 := filepath.Join(testDir, fileName+"Tests.swift")
	if a.fileExists(testFile1) {
		relPath, _ := filepath.Rel(projectPath, testFile1)
		testFiles = append(testFiles, relPath)
	}

	// Strategy 2: Mirror directory structure in Tests/
	relDir, err := filepath.Rel(filepath.Join(projectPath, "Sources"), dir)
	if err == nil && relDir != "." {
		// Try to find matching test target
		testTargets, _ := a.findTestTargets(projectPath)
		for _, target := range testTargets {
			testFile2 := filepath.Join(testDir, target, fileName+"Tests.swift")
			if a.fileExists(testFile2) {
				relPath, _ := filepath.Rel(projectPath, testFile2)
				testFiles = append(testFiles, relPath)
			}
		}
	}

	// Strategy 3: Search for tests that import this module
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
func (a *SwiftTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("swift: failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"@smoke",
		"// smoke",
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

// extractImports extracts import statements from Swift code
func (a *SwiftTestAnalyzer) extractImports(content string) []string {
	var imports []string
	importSet := make(map[string]bool)

	lines := strings.Split(content, "\n")
	importRe := regexp.MustCompile(swiftImportPattern)

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

// findTestsImportingFile finds tests that import a given module
func (a *SwiftTestAnalyzer) findTestsImportingFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	var testFiles []string

	// Get module name from Package.swift or project structure
	moduleName := a.detectModuleName(filePath, projectPath)
	if moduleName == "" {
		return nil, nil
	}

	testDir := filepath.Join(projectPath, "Tests")
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
		if !strings.HasSuffix(info.Name(), ".swift") {
			return nil
		}

		// Check if test imports the module
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		if a.importsModule(string(content), moduleName) {
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				return nil
			}
			testFiles = append(testFiles, relPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("swift: failed to search for tests: %w", err)
	}

	return testFiles, nil
}

// importsModule checks if content imports a specific module
func (a *SwiftTestAnalyzer) importsModule(content, moduleName string) bool {
	importPattern := fmt.Sprintf(`import\s+%s\b`, regexp.QuoteMeta(moduleName))
	re := regexp.MustCompile(importPattern)
	return re.MatchString(content)
}

// detectModuleName detects the module name for a source file
func (a *SwiftTestAnalyzer) detectModuleName(filePath, projectPath string) string {
	// Try to extract from Package.swift
	packagePath := filepath.Join(projectPath, "Package.swift")
	content, err := os.ReadFile(packagePath)
	if err != nil {
		return ""
	}

	// Extract target name from path
	relPath, err := filepath.Rel(filepath.Join(projectPath, "Sources"), filePath)
	if err != nil {
		return ""
	}

	parts := strings.Split(relPath, string(filepath.Separator))
	if len(parts) > 0 {
		targetName := parts[0]
		// Verify target exists in Package.swift
		if strings.Contains(string(content), targetName) {
			return targetName
		}
	}

	return ""
}

// findTestTargets finds test target directories
func (a *SwiftTestAnalyzer) findTestTargets(projectPath string) ([]string, error) {
	var targets []string

	testDir := filepath.Join(projectPath, "Tests")
	entries, err := os.ReadDir(testDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && !a.shouldSkipDirectory(entry.Name()) {
			targets = append(targets, entry.Name())
		}
	}

	return targets, nil
}

// fileExists checks if a file exists
func (a *SwiftTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// shouldSkipDirectory checks if directory should be skipped
func (a *SwiftTestAnalyzer) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		".build",
		".git",
		"DerivedData",
		"Pods",
		"Carthage",
		".swiftpm",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// removeDuplicates removes duplicate strings from slice
func (a *SwiftTestAnalyzer) removeDuplicates(items []string) []string {
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
func (a *SwiftTestAnalyzer) ExtractTestFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("swift: failed to read file: %w", err)
	}

	var functions []string
	testRe := regexp.MustCompile(swiftTestFuncPattern)

	matches := testRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	return functions, nil
}

// ExtractTestClasses extracts test class names from a file
func (a *SwiftTestAnalyzer) ExtractTestClasses(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("swift: failed to read file: %w", err)
	}

	var classes []string
	classRe := regexp.MustCompile(swiftTestClassPattern)

	matches := classRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			classes = append(classes, match[1])
		}
	}

	return classes, nil
}
