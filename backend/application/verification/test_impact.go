package verification

import (
	"context"
	"path/filepath"
	"strings"
	"sync"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// TestImpactAnalyzer analyzes which tests are affected by file changes
type TestImpactAnalyzer struct {
	symbolIndex analysis.SymbolIndex
	testService domain.ITestService
	log         domain.Logger

	mu        sync.RWMutex
	testIndex map[string][]string // file -> test files that cover it
	indexed   bool
}

// NewTestImpactAnalyzer creates a new test impact analyzer
func NewTestImpactAnalyzer(
	symbolIndex analysis.SymbolIndex,
	testService domain.ITestService,
	log domain.Logger,
) *TestImpactAnalyzer {
	return &TestImpactAnalyzer{
		symbolIndex: symbolIndex,
		testService: testService,
		log:         log,
		testIndex:   make(map[string][]string),
	}
}

// BuildTestIndex builds the test index for a project
func (a *TestImpactAnalyzer) BuildTestIndex(ctx context.Context, projectRoot string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.testIndex = make(map[string][]string)

	languages := a.testService.GetSupportedLanguages()
	for _, lang := range languages {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		suite, err := a.testService.DiscoverTests(ctx, projectRoot, lang)
		if err != nil {
			a.log.Warning("Failed to discover tests for " + lang + ": " + err.Error())
			continue
		}

		if suite == nil || len(suite.Tests) == 0 {
			continue
		}

		for _, test := range suite.Tests {
			for _, targetFile := range test.TargetFiles {
				a.testIndex[targetFile] = append(a.testIndex[targetFile], test.Path)
			}

			// Also map test file to itself
			a.testIndex[test.Path] = append(a.testIndex[test.Path], test.Path)
		}
	}

	a.indexed = true
	return nil
}

// GetAffectedTests returns test files affected by changed files
func (a *TestImpactAnalyzer) GetAffectedTests(changedFiles []string) []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	testSet := make(map[string]bool)

	for _, file := range changedFiles {
		// Direct mapping from index
		if tests, ok := a.testIndex[file]; ok {
			for _, test := range tests {
				testSet[test] = true
			}
		}

		// Check if file is a test file itself
		if a.isTestFile(file) {
			testSet[file] = true
		}

		// Find tests by convention (e.g., main.go -> main_test.go)
		conventionTests := a.findTestsByConvention(file)
		for _, test := range conventionTests {
			testSet[test] = true
		}
	}

	tests := make([]string, 0, len(testSet))
	for test := range testSet {
		tests = append(tests, test)
	}
	return tests
}

// RunTargetedTests runs only tests affected by changed files
func (a *TestImpactAnalyzer) RunTargetedTests(
	ctx context.Context,
	projectRoot string,
	changedFiles []string,
) ([]*domain.TestResult, error) {
	affectedTests := a.GetAffectedTests(changedFiles)

	if len(affectedTests) == 0 {
		return []*domain.TestResult{}, nil
	}

	config := &domain.TestConfig{
		ProjectPath:  projectRoot,
		Scope:        domain.TestScopeAffected,
		TestPatterns: affectedTests,
		Parallel:     true,
	}

	return a.testService.RunTargetedTests(ctx, config, changedFiles)
}

// IsIndexed returns whether the test index has been built
func (a *TestImpactAnalyzer) IsIndexed() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.indexed
}

// Clear clears the test index
func (a *TestImpactAnalyzer) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.testIndex = make(map[string][]string)
	a.indexed = false
}

// GetTestIndex returns a copy of the test index
func (a *TestImpactAnalyzer) GetTestIndex() map[string][]string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	result := make(map[string][]string, len(a.testIndex))
	for k, v := range a.testIndex {
		result[k] = append([]string{}, v...)
	}
	return result
}

// isTestFile checks if a file is a test file
func (a *TestImpactAnalyzer) isTestFile(file string) bool {
	base := filepath.Base(file)

	// Go test files
	if strings.HasSuffix(base, "_test.go") {
		return true
	}

	// JavaScript/TypeScript test files
	if strings.Contains(base, ".test.") || strings.Contains(base, ".spec.") {
		return true
	}

	// Python test files
	if strings.HasPrefix(base, "test_") && strings.HasSuffix(base, ".py") {
		return true
	}

	// Check if in test directory (normalize path separators)
	normalizedPath := strings.ReplaceAll(file, "\\", "/")
	return strings.Contains(normalizedPath, "/test/") ||
		strings.Contains(normalizedPath, "/tests/") ||
		strings.Contains(normalizedPath, "/__tests__/") ||
		strings.HasPrefix(normalizedPath, "tests/") ||
		strings.HasPrefix(normalizedPath, "__tests__/")
}

// findTestsByConvention finds test files by naming convention
func (a *TestImpactAnalyzer) findTestsByConvention(file string) []string {
	var tests []string
	ext := filepath.Ext(file)
	base := strings.TrimSuffix(filepath.Base(file), ext)
	dir := filepath.Dir(file)

	switch ext {
	case ".go":
		// main.go -> main_test.go
		testFile := filepath.Join(dir, base+"_test.go")
		if _, ok := a.testIndex[testFile]; ok {
			tests = append(tests, testFile)
		}

	case ".ts", ".tsx", ".js", ".jsx":
		// component.ts -> component.test.ts, component.spec.ts
		for _, pattern := range []string{".test", ".spec"} {
			testFile := filepath.Join(dir, base+pattern+ext)
			if _, ok := a.testIndex[testFile]; ok {
				tests = append(tests, testFile)
			}
		}
		// Also check __tests__ directory
		testsDir := filepath.Join(dir, "__tests__")
		testFile := filepath.Join(testsDir, base+ext)
		if _, ok := a.testIndex[testFile]; ok {
			tests = append(tests, testFile)
		}

	case ".py":
		// module.py -> test_module.py
		testFile := filepath.Join(dir, "test_"+base+".py")
		if _, ok := a.testIndex[testFile]; ok {
			tests = append(tests, testFile)
		}
	}

	return tests
}
