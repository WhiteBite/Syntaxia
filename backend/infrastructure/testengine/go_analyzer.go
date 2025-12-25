package testengine

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"strings"
)

// GoTestAnalyzer реализует TestAnalyzer для Go
type GoTestAnalyzer struct {
	log domain.Logger
}

// NewGoTestAnalyzer создает новый анализатор для Go тестов
func NewGoTestAnalyzer(log domain.Logger) *GoTestAnalyzer {
	return &GoTestAnalyzer{
		log: log,
	}
}

// AnalyzeTestDependencies анализирует зависимости тестов
func (a *GoTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Analyzing dependencies for test: %s", testPath))

	content, err := os.ReadFile(testPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read test file: %w", err)
	}

	var dependencies []string
	contentStr := string(content)

	// Ищем импорты
	importRegex := regexp.MustCompile(`import\s+\(([\s\S]*?)\)`)
	matches := importRegex.FindStringSubmatch(contentStr)

	if len(matches) > 1 {
		imports := matches[1]
		// Разбираем импорты
		lines := strings.Split(imports, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}

			// Убираем кавычки и комментарии
			line = strings.Trim(line, `"`)
			if idx := strings.Index(line, "//"); idx != -1 {
				line = strings.TrimSpace(line[:idx])
			}

			if line != "" {
				dependencies = append(dependencies, line)
			}
		}
	}

	// Ищем также однострочные импорты
	singleImportRegex := regexp.MustCompile(`import\s+"([^"]+)"`)
	singleMatches := singleImportRegex.FindAllStringSubmatch(contentStr, -1)
	for _, match := range singleMatches {
		if len(match) > 1 {
			dependencies = append(dependencies, match[1])
		}
	}

	a.log.Info(fmt.Sprintf("Found %d dependencies for test %s", len(dependencies), testPath))
	return dependencies, nil
}

// FindTestsForFile находит тесты для файла
func (a *GoTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	a.log.Info(fmt.Sprintf("Finding tests for file: %s", filePath))

	var testFiles []string

	// Получаем имя файла без расширения
	fileName := strings.TrimSuffix(filepath.Base(filePath), ".go")
	dir := filepath.Dir(filePath)

	// Ищем соответствующий тестовый файл
	testFile := filepath.Join(dir, fileName+"_test.go")
	fullTestPath := filepath.Join(projectPath, testFile)

	if _, err := os.Stat(fullTestPath); err == nil {
		testFiles = append(testFiles, testFile)
	}

	// Также ищем тесты в том же пакете
	packageName := a.getPackageName(filepath.Join(projectPath, filePath))
	if packageName != "" {
		packageTests, err := a.findTestsInPackage(ctx, dir, packageName, projectPath)
		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to find tests in package: %v", err))
		} else {
			testFiles = append(testFiles, packageTests...)
		}
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
func (a *GoTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
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
		"smoke:",
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

// findTestsInPackage находит все тесты в пакете
func (a *GoTestAnalyzer) findTestsInPackage(ctx context.Context, dir, packageName, projectPath string) ([]string, error) {
	var tests []string

	// Ищем все *_test.go файлы в директории
	testPattern := filepath.Join(projectPath, dir, "*_test.go")
	matches, err := filepath.Glob(testPattern)
	if err != nil {
		return nil, fmt.Errorf("failed to glob test files: %w", err)
	}

	for _, match := range matches {
		// Проверяем, что тест принадлежит тому же пакету
		if a.getPackageName(match) == packageName {
			relPath, err := filepath.Rel(projectPath, match)
			if err != nil {
				a.log.Warning(fmt.Sprintf("Failed to get relative path for %s: %v", match, err))
				continue
			}
			tests = append(tests, relPath)
		}
	}

	return tests, nil
}

// getPackageName получает имя пакета из Go файла
func (a *GoTestAnalyzer) getPackageName(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return ""
	}

	// Простой regex для поиска package
	re := regexp.MustCompile(`package\s+(\w+)`)
	matches := re.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
