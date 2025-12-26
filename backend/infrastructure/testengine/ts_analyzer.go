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

// TypeScriptTestAnalyzer реализует TestAnalyzer для TypeScript/JavaScript
type TypeScriptTestAnalyzer struct {
	log domain.Logger
}

// NewTypeScriptTestAnalyzer создает новый анализатор для TypeScript тестов
func NewTypeScriptTestAnalyzer(log domain.Logger) *TypeScriptTestAnalyzer {
	return &TypeScriptTestAnalyzer{
		log: log,
	}
}

// AnalyzeTestDependencies анализирует зависимости тестов
func (a *TypeScriptTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for TypeScript test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	var dependencies []string
	contentStr := string(content)

	// Ищем ES6 импорты: import ... from '...'
	es6ImportRegex := regexp.MustCompile(`import\s+(?:[\w\s{},*]+\s+from\s+)?['"]([^'"]+)['"]`)
	es6Matches := es6ImportRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range es6Matches {
		if len(match) > 1 {
			dependencies = append(dependencies, match[1])
		}
	}

	// Ищем CommonJS require: require('...')
	requireRegex := regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	requireMatches := requireRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range requireMatches {
		if len(match) > 1 {
			dependencies = append(dependencies, match[1])
		}
	}

	// Ищем динамические импорты: import('...')
	dynamicImportRegex := regexp.MustCompile(`import\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	dynamicMatches := dynamicImportRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range dynamicMatches {
		if len(match) > 1 {
			dependencies = append(dependencies, match[1])
		}
	}

	// Убираем дубликаты
	uniqueDeps := make(map[string]bool)
	var result []string
	for _, dep := range dependencies {
		if !uniqueDeps[dep] {
			uniqueDeps[dep] = true
			result = append(result, dep)
		}
	}

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(result), testPath))
	return result, nil
}

// FindTestsForFile находит тесты для файла
func (a *TypeScriptTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding tests for TypeScript file: %s", filePath))

	var testFiles []string

	// Получаем имя файла без расширения
	ext := filepath.Ext(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), ext)
	dir := filepath.Dir(filePath)

	// Возможные расширения тестовых файлов
	testExtensions := []string{".test.ts", ".spec.ts", ".test.tsx", ".spec.tsx", ".test.js", ".spec.js"}

	// Ищем соответствующие тестовые файлы в той же директории
	for _, testExt := range testExtensions {
		testFile := filepath.Join(dir, baseName+testExt)
		fullTestPath := filepath.Join(projectPath, testFile)

		if _, err := os.Stat(fullTestPath); err == nil {
			testFiles = append(testFiles, testFile)
		}
	}

	// Ищем тесты в директории __tests__
	testsDir := filepath.Join(dir, "__tests__")
	fullTestsDir := filepath.Join(projectPath, testsDir)
	if info, err := os.Stat(fullTestsDir); err == nil && info.IsDir() {
		for _, testExt := range testExtensions {
			testFile := filepath.Join(testsDir, baseName+testExt)
			fullTestPath := filepath.Join(projectPath, testFile)

			if _, err := os.Stat(fullTestPath); err == nil {
				testFiles = append(testFiles, testFile)
			}
		}
	}

	// Ищем тесты, которые импортируют данный файл
	importTests, err := a.findTestsImportingFile(ctx, filePath, projectPath)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to find tests importing file: %v", err))
	} else {
		testFiles = append(testFiles, importTests...)
	}

	// Убираем дубликаты
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

// IsSmokeTest определяет, является ли тест smoke тестом
func (a *TypeScriptTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	content, err := os.ReadFile(testPath)
	if err != nil {
		return false, fmt.Errorf("failed to read test file: %w", err)
	}

	contentStr := strings.ToLower(string(content))

	// Проверяем различные признаки smoke тестов
	smokeIndicators := []string{
		"smoke",
		"smoke_test",
		"smoketest",
		"// smoke",
		"/* smoke",
		"@smoke",
		"describe.smoke",
		"it.smoke",
	}

	for _, indicator := range smokeIndicators {
		if strings.Contains(contentStr, indicator) {
			return true, nil
		}
	}

	// Проверяем имя файла
	fileName := strings.ToLower(filepath.Base(testPath))
	if strings.Contains(fileName, "smoke") {
		return true, nil
	}

	return false, nil
}

// findTestsImportingFile находит тесты, которые импортируют указанный файл
func (a *TypeScriptTestAnalyzer) findTestsImportingFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	var tests []string

	// Получаем имя модуля для поиска в импортах
	ext := filepath.Ext(filePath)
	baseName := strings.TrimSuffix(filepath.Base(filePath), ext)
	dir := filepath.Dir(filePath)

	// Паттерны для поиска тестовых файлов
	testPatterns := []string{
		"*.test.ts", "*.spec.ts",
		"*.test.tsx", "*.spec.tsx",
		"*.test.js", "*.spec.js",
	}

	// Ищем тестовые файлы в проекте
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем node_modules и другие служебные директории
			if a.shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Проверяем, является ли файл тестовым
		fileName := info.Name()
		isTestFile := false
		for _, pattern := range testPatterns {
			if matched, _ := filepath.Match(pattern, fileName); matched {
				isTestFile = true
				break
			}
		}

		if !isTestFile {
			return nil
		}

		// Читаем содержимое тестового файла
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Проверяем, импортирует ли тест наш файл
		contentStr := string(content)
		relDir := strings.ReplaceAll(dir, "\\", "/")

		// Паттерны для поиска импорта
		importPatterns := []string{
			fmt.Sprintf(`from\s+['"].*%s['"]`, baseName),
			fmt.Sprintf(`from\s+['"].*/%s['"]`, baseName),
			fmt.Sprintf(`require\s*\(\s*['"].*%s['"]\s*\)`, baseName),
			fmt.Sprintf(`from\s+['"].*%s/%s['"]`, relDir, baseName),
		}

		for _, pattern := range importPatterns {
			re := regexp.MustCompile(pattern)
			if re.MatchString(contentStr) {
				relPath, err := filepath.Rel(projectPath, path)
				if err != nil {
					continue
				}
				tests = append(tests, relPath)
				break
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to search for tests: %w", err)
	}

	return tests, nil
}

// shouldSkipDir проверяет, нужно ли пропустить директорию
func (a *TypeScriptTestAnalyzer) shouldSkipDir(name string) bool {
	skipDirs := []string{
		"node_modules",
		"dist",
		"build",
		"coverage",
		".git",
		".next",
		".nuxt",
		"out",
	}

	for _, skip := range skipDirs {
		if name == skip {
			return true
		}
	}

	// Пропускаем скрытые директории
	if strings.HasPrefix(name, ".") && name != "." && name != ".." {
		return true
	}

	return false
}
