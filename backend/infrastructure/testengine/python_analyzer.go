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

// Python analyzer constants
const (
	// Import patterns
	pythonImportPattern       = `^import\s+(\w+(?:\.\w+)*)`
	pythonFromImportPattern   = `^from\s+(\w+(?:\.\w+)*)\s+import`
	pythonRelativeImportPat   = `^from\s+\.+(\w*(?:\.\w+)*)\s+import`
	
	// Test patterns
	pythonTestFuncPattern     = `def\s+(test_\w+)\s*\(`
	pythonTestClassPattern    = `class\s+(Test\w+)\s*[:\(]`
	pythonTestMethodPattern   = `def\s+(test_\w+)\s*\(self`
)

// PythonTestAnalyzer implements domain.TestAnalyzer for Python
type PythonTestAnalyzer struct {
	log domain.Logger
}

// NewPythonTestAnalyzer creates a new Python test analyzer
func NewPythonTestAnalyzer(log domain.Logger) *PythonTestAnalyzer {
	return &PythonTestAnalyzer{log: log}
}

// AnalyzeTestDependencies analyzes test dependencies
func (a *PythonTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for Python test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	dependencies := a.extractImports(string(content))

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile finds tests for a given source file
func (a *PythonTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding Python tests for file: %s", filePath))

	var testFiles []string

	// Get file name without extension
	fileName := strings.TrimSuffix(filepath.Base(filePath), ".py")
	dir := filepath.Dir(filePath)

	// Strategy 1: Look for test_<filename>.py in same directory
	testFile1 := filepath.Join(dir, "test_"+fileName+".py")
	if a.fileExists(filepath.Join(projectPath, testFile1)) {
		testFiles = append(testFiles, testFile1)
	}

	// Strategy 2: Look for <filename>_test.py in same directory
	testFile2 := filepath.Join(dir, fileName+"_test.py")
	if a.fileExists(filepath.Join(projectPath, testFile2)) {
		testFiles = append(testFiles, testFile2)
	}

	// Strategy 3: Look in tests/ directory at same level
	testsDir := filepath.Join(dir, "tests")
	testFile3 := filepath.Join(testsDir, "test_"+fileName+".py")
	if a.fileExists(filepath.Join(projectPath, testFile3)) {
		testFiles = append(testFiles, testFile3)
	}

	// Strategy 4: Look in project root tests/ directory
	rootTestsDir := "tests"
	testFile4 := filepath.Join(rootTestsDir, "test_"+fileName+".py")
	if a.fileExists(filepath.Join(projectPath, testFile4)) {
		testFiles = append(testFiles, testFile4)
	}

	// Strategy 5: Search for tests that import this module
	importTests, err := a.findTestsImportingModule(ctx, filePath, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests importing module: %v", err))
	} else {
		testFiles = append(testFiles, importTests...)
	}

	// Remove duplicates
	testFiles = a.removeDuplicates(testFiles)

	a.log.Info(fmt.Sprintf("Found %d test files for %s", len(testFiles), filePath))
	return testFiles, nil
}

// IsSmokeTest determines if a test is a smoke test
func (a *PythonTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Check for smoke test indicators
	smokeIndicators := []string{
		"smoke",
		"@pytest.mark.smoke",
		"@smoke",
		"# smoke",
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

// extractImports extracts import statements from Python code
func (a *PythonTestAnalyzer) extractImports(content string) []string {
	var imports []string
	importSet := make(map[string]bool)

	lines := strings.Split(content, "\n")

	// Compile regex patterns
	importRe := regexp.MustCompile(pythonImportPattern)
	fromImportRe := regexp.MustCompile(pythonFromImportPattern)
	relativeImportRe := regexp.MustCompile(pythonRelativeImportPat)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for "import X" pattern
		if matches := importRe.FindStringSubmatch(line); len(matches) > 1 {
			module := matches[1]
			if !importSet[module] {
				imports = append(imports, module)
				importSet[module] = true
			}
		}

		// Check for "from X import Y" pattern
		if matches := fromImportRe.FindStringSubmatch(line); len(matches) > 1 {
			module := matches[1]
			if !importSet[module] {
				imports = append(imports, module)
				importSet[module] = true
			}
		}

		// Check for relative imports "from .X import Y"
		if matches := relativeImportRe.FindStringSubmatch(line); len(matches) > 1 {
			module := matches[1]
			if module != "" && !importSet[module] {
				imports = append(imports, "."+module)
				importSet[module] = true
			}
		}
	}

	return imports
}

// findTestsImportingModule finds tests that import a given module
func (a *PythonTestAnalyzer) findTestsImportingModule(ctx context.Context, filePath, projectPath string) ([]string, error) {
	var testFiles []string

	// Convert file path to module name
	moduleName := a.filePathToModuleName(filePath)
	if moduleName == "" {
		return nil, nil
	}

	// Walk through project looking for test files
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			// Skip excluded directories
			if a.shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Only check test files
		if !a.isTestFile(info.Name()) {
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
		return nil, fmt.Errorf("failed to search for tests: %w", err)
	}

	return testFiles, nil
}

// importsModule checks if content imports a specific module
func (a *PythonTestAnalyzer) importsModule(content, moduleName string) bool {
	// Check for direct import
	if strings.Contains(content, "import "+moduleName) {
		return true
	}

	// Check for from import
	if strings.Contains(content, "from "+moduleName) {
		return true
	}

	// Check for partial module name (e.g., "from mypackage.mymodule import")
	parts := strings.Split(moduleName, ".")
	if len(parts) > 0 {
		lastPart := parts[len(parts)-1]
		if strings.Contains(content, "import "+lastPart) ||
			strings.Contains(content, "from "+lastPart) {
			return true
		}
	}

	return false
}

// filePathToModuleName converts file path to Python module name
func (a *PythonTestAnalyzer) filePathToModuleName(filePath string) string {
	// Remove .py extension
	modulePath := strings.TrimSuffix(filePath, ".py")

	// Replace path separators with dots
	modulePath = strings.ReplaceAll(modulePath, string(filepath.Separator), ".")
	modulePath = strings.ReplaceAll(modulePath, "/", ".")

	// Remove leading dots
	modulePath = strings.TrimPrefix(modulePath, ".")

	return modulePath
}

// fileExists checks if a file exists
func (a *PythonTestAnalyzer) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// shouldSkipDirectory checks if directory should be skipped
func (a *PythonTestAnalyzer) shouldSkipDirectory(name string) bool {
	skipDirs := []string{
		"__pycache__",
		".venv",
		"venv",
		".tox",
		".eggs",
		"build",
		"dist",
		"node_modules",
		".git",
		".pytest_cache",
		".mypy_cache",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	return false
}

// isTestFile checks if file is a Python test file
func (a *PythonTestAnalyzer) isTestFile(name string) bool {
	if !strings.HasSuffix(name, ".py") {
		return false
	}

	// Check for test_*.py pattern
	if strings.HasPrefix(name, "test_") {
		return true
	}

	// Check for *_test.py pattern
	if strings.HasSuffix(name, "_test.py") {
		return true
	}

	return false
}

// removeDuplicates removes duplicate strings from slice
func (a *PythonTestAnalyzer) removeDuplicates(items []string) []string {
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
func (a *PythonTestAnalyzer) ExtractTestFunctions(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var functions []string
	funcRe := regexp.MustCompile(pythonTestFuncPattern)

	matches := funcRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			functions = append(functions, match[1])
		}
	}

	return functions, nil
}

// ExtractTestClasses extracts test class names from a file
func (a *PythonTestAnalyzer) ExtractTestClasses(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var classes []string
	classRe := regexp.MustCompile(pythonTestClassPattern)

	matches := classRe.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) > 1 {
			classes = append(classes, match[1])
		}
	}

	return classes, nil
}
