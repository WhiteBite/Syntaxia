package testengine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// CSharpTestRunner реализует TestRunner для C#
type CSharpTestRunner struct {
	log domain.Logger
}

// NewCSharpTestRunner создает новый runner для C# тестов
func NewCSharpTestRunner(log domain.Logger) *CSharpTestRunner {
	return &CSharpTestRunner{
		log: log,
	}
}

// GetLanguage возвращает язык
func (r *CSharpTestRunner) GetLanguage() string {
	return "csharp"
}

// RunTest выполняет один C# тест
func (r *CSharpTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running C# test: %s", testPath))

	startTime := time.Now()

	// Получаем имя тестового класса/метода из пути
	testFilter := r.extractTestFilter(testPath)

	// Строим команду для запуска теста
	args := []string{"test", "--nologo"}
	if testFilter != "" {
		args = append(args, "--filter", testFilter)
	}
	if config.Verbose {
		args = append(args, "-v", "detailed")
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "dotnet", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Устанавливаем переменные окружения
	cmd.Env = os.Environ()
	if config.EnvVars != nil {
		for key, value := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	// Запускаем команду
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	result := &domain.TestResult{
		TestPath: testPath,
		TestName: testFilter,
		Language: "csharp",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("C# test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("C# test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор C# тестов
func (r *CSharpTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running C# test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed C# test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает C# тесты в проекте
func (r *CSharpTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering C# tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Ищем тестовые проекты (обычно содержат .Tests или Test в имени)
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			// Пропускаем bin, obj, .git
			if info.Name() == "bin" || info.Name() == "obj" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Проверяем, является ли файл C# тестом
		if !r.isCSharpTestFile(info.Name(), path) {
			return nil
		}

		// Получаем относительный путь
		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to get relative path for %s: %v", path, err))
			return nil
		}

		// Анализируем тест файл для определения типа тестов
		testType := r.analyzeTestFile(path)

		testInfo := &domain.TestInfo{
			Path: relPath,
			Name: info.Name(),
			Type: testType,
			Metadata: map[string]string{
				"framework": r.detectTestFramework(path),
				"className": r.extractClassName(path),
			},
		}

		tests = append(tests, testInfo)
		return nil
	})

	if err != nil {
		r.log.Warning(fmt.Sprintf("Failed to walk project directory: %v", err))
	}

	r.log.Info(fmt.Sprintf("Discovered %d C# test files", len(tests)))
	return tests, nil
}

// extractTestFilter извлекает фильтр теста из пути к файлу
func (r *CSharpTestRunner) extractTestFilter(testPath string) string {
	// Убираем расширение .cs
	className := strings.TrimSuffix(filepath.Base(testPath), ".cs")
	return fmt.Sprintf("FullyQualifiedName~%s", className)
}

// isCSharpTestFile проверяет, является ли файл C# тестом
func (r *CSharpTestRunner) isCSharpTestFile(fileName string, filePath string) bool {
	if !strings.HasSuffix(fileName, ".cs") {
		return false
	}

	// Стандартные паттерны именования тестов
	testPatterns := []string{
		"Test.cs",
		"Tests.cs",
		"TestCase.cs",
		"Spec.cs",
		"Specs.cs",
	}

	for _, pattern := range testPatterns {
		if strings.HasSuffix(fileName, pattern) {
			return true
		}
	}

	// Проверяем содержимое файла на наличие тестовых атрибутов
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

// analyzeTestFile анализирует файл теста для определения его типа
func (r *CSharpTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit"
	}

	contentStr := strings.ToLower(string(content))
	fileName := strings.ToLower(filepath.Base(filePath))

	// Проверяем на smoke тесты
	if strings.Contains(contentStr, "smoke") || strings.Contains(fileName, "smoke") {
		return "smoke"
	}

	// Проверяем на integration тесты
	integrationIndicators := []string{
		"integration",
		"webapplicationfactory",
		"testserver",
		"httpclient",
		"[collection(",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentStr, indicator) || strings.Contains(fileName, indicator) {
			return "integration"
		}
	}

	return "unit"
}

// detectTestFramework определяет используемый тестовый фреймворк
func (r *CSharpTestRunner) detectTestFramework(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unknown"
	}

	contentStr := string(content)

	// xUnit
	if strings.Contains(contentStr, "[Fact]") || strings.Contains(contentStr, "[Theory]") {
		return "xunit"
	}

	// NUnit
	if strings.Contains(contentStr, "[Test]") || strings.Contains(contentStr, "[TestCase]") {
		return "nunit"
	}

	// MSTest
	if strings.Contains(contentStr, "[TestMethod]") || strings.Contains(contentStr, "[TestClass]") {
		return "mstest"
	}

	return "unknown"
}

// extractClassName извлекает имя класса из файла
func (r *CSharpTestRunner) extractClassName(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return strings.TrimSuffix(filepath.Base(filePath), ".cs")
	}

	// Ищем объявление класса
	classRe := regexp.MustCompile(`class\s+(\w+)`)
	matches := classRe.FindStringSubmatch(string(content))
	if len(matches) > 1 {
		return matches[1]
	}

	return strings.TrimSuffix(filepath.Base(filePath), ".cs")
}

// ParseTestOutput парсит вывод dotnet test
func (r *CSharpTestRunner) ParseTestOutput(output string) *domain.TestValidationResult {
	result := &domain.TestValidationResult{
		Success:         true,
		FailedTestPaths: []string{},
	}

	lines := strings.Split(output, "\n")

	// Regex для парсинга результатов
	passedRe := regexp.MustCompile(`Passed:\s*(\d+)`)
	failedRe := regexp.MustCompile(`Failed:\s*(\d+)`)
	skippedRe := regexp.MustCompile(`Skipped:\s*(\d+)`)
	totalRe := regexp.MustCompile(`Total:\s*(\d+)`)

	for _, line := range lines {
		if matches := passedRe.FindStringSubmatch(line); len(matches) > 1 {
			result.PassedTests, _ = strconv.Atoi(matches[1])
		}
		if matches := failedRe.FindStringSubmatch(line); len(matches) > 1 {
			result.FailedTests, _ = strconv.Atoi(matches[1])
		}
		if matches := skippedRe.FindStringSubmatch(line); len(matches) > 1 {
			result.SkippedTests, _ = strconv.Atoi(matches[1])
		}
		if matches := totalRe.FindStringSubmatch(line); len(matches) > 1 {
			result.TotalTests, _ = strconv.Atoi(matches[1])
		}

		// Ищем failed тесты
		if strings.Contains(line, "Failed") && strings.Contains(line, ".") {
			// Извлекаем имя теста
			failedTestRe := regexp.MustCompile(`Failed\s+([\w.]+)`)
			if matches := failedTestRe.FindStringSubmatch(line); len(matches) > 1 {
				result.FailedTestPaths = append(result.FailedTestPaths, matches[1])
			}
		}
	}

	result.Success = result.FailedTests == 0
	if result.TotalTests > 0 {
		result.SuccessRate = float64(result.PassedTests) / float64(result.TotalTests) * 100
	}

	return result
}

// parseInt parses string to int, returns 0 on error
func parseInt(s string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}
