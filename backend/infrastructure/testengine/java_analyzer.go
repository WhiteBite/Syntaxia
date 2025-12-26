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

// Java analyzer constants
const (
	javaTestAnnotation     = "@Test"
	javaTestClassSuffix    = "Test.java"
	javaTestClassSuffixAlt = "Tests.java"
)

// JavaTestAnalyzer implements domain.TestAnalyzer for Java
type JavaTestAnalyzer struct {
	log domain.Logger
}

// NewJavaTestAnalyzer creates a new Java test analyzer
func NewJavaTestAnalyzer(log domain.Logger) *JavaTestAnalyzer {
	return &JavaTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *JavaTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for Java test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractImports(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *JavaTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding Java tests for file: %s", filePath))

	var testFiles []string

	// Get class name from file
	className := strings.TrimSuffix(filepath.Base(filePath), ".java")

	// Standard test directories
	testDirs := []string{
		"src/test/java",
		"src/test",
		"test",
		"tests",
	}

	// Look for test files in standard locations
	for _, testDir := range testDirs {
		// Strategy 1: Same package structure in test directory
		relDir := filepath.Dir(filePath)
		if strings.HasPrefix(relDir, "src/main/java/") {
			relDir = strings.TrimPrefix(relDir, "src/main/java/")
		}

		testFile := filepath.Join(testDir, relDir, className+"Test.java")
		if a.fileExists(filepath.Join(projectPath, testFile)) {
			testFiles = append(testFiles, testFile)
		}

		testFile2 := filepath.Join(testDir, relDir, className+"Tests.java")
		if a.fileExists(filepath.Join(projectPath, testFile2)) {
			testFiles = append(testFiles, testFile2)
		}
	}

	// Strategy 2: Search for tests that import this class
	importTests, err := a.findTestsImportingClass(ctx, className, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests importing class: %v", err))
	} else {
		testFiles = append(testFiles, importTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *JavaTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"@tag(\"smoke\")",
		"@category(smoke",
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

// extractImports extracts import statements from Java code
func (a *JavaTestAnalyzer) extractImports(content string) []string {
	var imports []string
	importSet := make(map[string]bool)

	// Match import statements
	importRe := regexp.MustCompile(`import\s+(static\s+)?([a-zA-Z_][\w.]*);`)
	matches := importRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			importPath := match[2]
			if !importSet[importPath] {
				imports = append(imports, importPath)
				importSet[importPath] = true
			}
		}
	}

	return imports
}

// findTestsImportingClass finds tests that import a given class
func (a *JavaTestAnalyzer) findTestsImportingClass(ctx context.Context, className, projectPath string) ([]string, error) {
	var testFiles []string

	testDirs := []string{
		"src/test/java",
		"src/test",
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
				return nil
			}

			if info.IsDir() || !strings.HasSuffix(info.Name(), ".java") {
				return nil
			}

			// Check if file imports the class
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}

			if a.importsClass(string(content), className) {
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

// importsClass checks if content imports a specific class
func (a *JavaTestAnalyzer) importsClass(content, className string) bool {
	// Check for direct import
	pattern := fmt.Sprintf(`import\s+[\w.]*\.%s;`, className)
	matched, _ := regexp.MatchString(pattern, content)
	if matched {
		return true
	}

	// Check for wildcard import (less precise)
	if strings.Contains(content, className) {
		return true
	}

	return false
}

// fileExists checks if a file exists
func (a *JavaTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// removeDuplicates removes duplicate strings from slice
func (a *JavaTestAnalyzer) removeDuplicates(items []string) []string {
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
