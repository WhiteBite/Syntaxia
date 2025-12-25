package testengine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"strings"
	"time"
)

// GoTestRunner реализует TestRunner для Go
type GoTestRunner struct {
	log domain.Logger
}

// NewGoTestRunner создает новый runner для Go тестов
func NewGoTestRunner(log domain.Logger) *GoTestRunner {
	return &GoTestRunner{
		log: log,
	}
}

// RunTest выполняет один Go тест
func (r *GoTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Go test: %s", testPath))

	startTime := time.Now()

	// Строим команду для запуска теста
	args := []string{"test"}

	if config.Verbose {
		args = append(args, "-v")
	}

	if config.Coverage {
		args = append(args, "-cover")
	}

	// Добавляем таймаут если указан
	if config.Timeout > 0 {
		args = append(args, "-timeout", fmt.Sprintf("%ds", config.Timeout))
	}

	// Добавляем путь к тесту
	args = append(args, testPath)

	// Создаем команду
	cmd := exec.CommandContext(ctx, "go", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Устанавливаем переменные окружения
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
		TestName: filepath.Base(testPath),
		Language: "go",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Go test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Go test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор Go тестов
func (r *GoTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Go test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	// Если параллельное выполнение включено, запускаем тесты параллельно
	if suite.Config != nil && suite.Config.Parallel {
		// Простая реализация параллельного выполнения
		// В будущем можно добавить более сложную логику с горутинами
		for _, test := range suite.Tests {
			result, err := r.RunTest(ctx, test.Path, suite.Config)
			if err != nil {
				r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
			}
			results = append(results, result)
		}
	} else {
		// Последовательное выполнение
		for _, test := range suite.Tests {
			result, err := r.RunTest(ctx, test.Path, suite.Config)
			if err != nil {
				r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
			}
			results = append(results, result)
		}
	}

	r.log.Info(fmt.Sprintf("Completed Go test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает Go тесты в проекте
func (r *GoTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Go tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Ищем все *_test.go файлы
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем vendor и node_modules
			if info.Name() == "vendor" || info.Name() == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, "_test.go") {
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
			Name: filepath.Base(path),
			Type: testType,
			Metadata: map[string]string{
				"package": r.getPackageName(path),
			},
		}

		tests = append(tests, testInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover Go tests: %w", err)
	}

	r.log.Info(fmt.Sprintf("Discovered %d Go test files", len(tests)))
	return tests, nil
}

// GetLanguage возвращает язык
func (r *GoTestRunner) GetLanguage() string {
	return "go"
}

// analyzeTestFile анализирует файл теста для определения его типа
func (r *GoTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit" // По умолчанию считаем unit тестом
	}

	contentStr := string(content)

	// Проверяем на smoke тесты (обычно содержат "smoke" в названии или комментариях)
	if strings.Contains(strings.ToLower(contentStr), "smoke") {
		return "smoke"
	}

	// Проверяем на integration тесты (обычно используют внешние зависимости)
	if strings.Contains(contentStr, "integration") ||
		strings.Contains(contentStr, "TestMain") ||
		strings.Contains(contentStr, "database") ||
		strings.Contains(contentStr, "http") {
		return "integration"
	}

	// По умолчанию считаем unit тестом
	return "unit"
}

// getPackageName получает имя пакета из Go файла
func (r *GoTestRunner) getPackageName(filePath string) string {
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
