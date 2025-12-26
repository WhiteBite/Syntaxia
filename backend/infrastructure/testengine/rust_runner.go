package testengine

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// Rust test runner constants
const (
	rustTestTimeout     = 300 // 5 minutes default timeout
	rustTestFilePattern = "_test.rs"
	rustTestDirName     = "tests"
)

// RustTestRunner реализует TestRunner для Rust
type RustTestRunner struct {
	log domain.Logger
}

// NewRustTestRunner создает новый runner для Rust тестов
func NewRustTestRunner(log domain.Logger) *RustTestRunner {
	return &RustTestRunner{
		log: log,
	}
}

// RunTest выполняет один Rust тест
func (r *RustTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Rust test: %s", testPath))

	startTime := time.Now()

	// Строим команду для запуска теста
	args := []string{"test"}

	// Добавляем имя теста если указан конкретный тест
	if testPath != "" && testPath != "./..." {
		// Извлекаем имя теста из пути
		testName := r.extractTestName(testPath)
		if testName != "" {
			args = append(args, testName)
		}
	}

	// Добавляем флаги
	if config.Verbose {
		args = append(args, "--", "--nocapture")
	}

	// Добавляем таймаут если указан
	timeout := config.Timeout
	if timeout == 0 {
		timeout = rustTestTimeout
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "cargo", args...)
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
		TestName: filepath.Base(testPath),
		Language: "rust",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Rust test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Rust test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор Rust тестов
func (r *RustTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Rust test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	// Для Rust эффективнее запустить все тесты одной командой
	if len(suite.Tests) > 1 {
		result, err := r.runAllTests(ctx, suite.ProjectPath, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run all tests: %v", err))
		}
		results = append(results, result)
	} else {
		// Запускаем тесты по одному
		for _, test := range suite.Tests {
			result, err := r.RunTest(ctx, test.Path, suite.Config)
			if err != nil {
				r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
			}
			results = append(results, result)
		}
	}

	r.log.Info(fmt.Sprintf("Completed Rust test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает Rust тесты в проекте
func (r *RustTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Rust tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Проверяем наличие Cargo.toml
	cargoPath := filepath.Join(projectPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Cargo.toml not found in %s", projectPath)
	}

	// Ищем тесты в src/ (unit tests) и tests/ (integration tests)
	srcTests, err := r.discoverUnitTests(projectPath)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Failed to discover unit tests: %v", err))
	} else {
		tests = append(tests, srcTests...)
	}

	integrationTests, err := r.discoverIntegrationTests(projectPath)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Failed to discover integration tests: %v", err))
	} else {
		tests = append(tests, integrationTests...)
	}

	r.log.Info(fmt.Sprintf("Discovered %d Rust test files", len(tests)))
	return tests, nil
}

// GetLanguage возвращает язык
func (r *RustTestRunner) GetLanguage() string {
	return "rust"
}

// runAllTests запускает все тесты проекта
func (r *RustTestRunner) runAllTests(ctx context.Context, projectPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	startTime := time.Now()

	args := []string{"test"}
	if config != nil && config.Verbose {
		args = append(args, "--", "--nocapture")
	}

	cmd := exec.CommandContext(ctx, "cargo", args...)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime).Seconds()

	result := &domain.TestResult{
		TestPath: projectPath,
		TestName: "all",
		Language: "rust",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
	} else {
		result.Success = true
	}

	return result, nil
}

// discoverUnitTests находит unit тесты в src/
func (r *RustTestRunner) discoverUnitTests(projectPath string) ([]*domain.TestInfo, error) {
	var tests []*domain.TestInfo

	srcPath := filepath.Join(projectPath, "src")
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return tests, nil
	}

	err := filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".rs") {
			return nil
		}

		// Проверяем наличие #[test] или mod tests в файле
		hasTests, err := r.fileHasTests(path)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to check tests in %s: %v", path, err))
			return nil
		}

		if hasTests {
			relPath, err := filepath.Rel(projectPath, path)
			if err != nil {
				return nil
			}

			testInfo := &domain.TestInfo{
				Path: relPath,
				Name: filepath.Base(path),
				Type: "unit",
				Metadata: map[string]string{
					"module": r.getModuleName(path),
				},
			}
			tests = append(tests, testInfo)
		}

		return nil
	})

	return tests, err
}

// discoverIntegrationTests находит integration тесты в tests/
func (r *RustTestRunner) discoverIntegrationTests(projectPath string) ([]*domain.TestInfo, error) {
	var tests []*domain.TestInfo

	testsPath := filepath.Join(projectPath, rustTestDirName)
	if _, err := os.Stat(testsPath); os.IsNotExist(err) {
		return tests, nil
	}

	err := filepath.Walk(testsPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".rs") {
			return nil
		}

		relPath, err := filepath.Rel(projectPath, path)
		if err != nil {
			return nil
		}

		testType := "integration"
		// Проверяем на smoke тесты
		if strings.Contains(strings.ToLower(filepath.Base(path)), "smoke") {
			testType = "smoke"
		}

		testInfo := &domain.TestInfo{
			Path: relPath,
			Name: filepath.Base(path),
			Type: testType,
			Metadata: map[string]string{
				"module": r.getModuleName(path),
			},
		}
		tests = append(tests, testInfo)

		return nil
	})

	return tests, err
}

// fileHasTests проверяет наличие тестов в файле
func (r *RustTestRunner) fileHasTests(filePath string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	contentStr := string(content)

	// Ищем #[test] атрибут или mod tests блок
	testPatterns := []string{
		"#[test]",
		"#[cfg(test)]",
		"mod tests",
	}

	for _, pattern := range testPatterns {
		if strings.Contains(contentStr, pattern) {
			return true, nil
		}
	}

	return false, nil
}

// getModuleName извлекает имя модуля из пути файла
func (r *RustTestRunner) getModuleName(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".rs")
}

// extractTestName извлекает имя теста из пути
func (r *RustTestRunner) extractTestName(testPath string) string {
	// Для Rust тестов имя теста - это имя модуля или функции
	base := filepath.Base(testPath)
	return strings.TrimSuffix(base, ".rs")
}

// parseTestOutput парсит вывод cargo test для извлечения результатов
func (r *RustTestRunner) parseTestOutput(output string) (passed, failed, ignored int) {
	// Формат: "test result: ok. 5 passed; 0 failed; 1 ignored"
	re := regexp.MustCompile(`test result:.*?(\d+) passed.*?(\d+) failed.*?(\d+) ignored`)
	matches := re.FindStringSubmatch(output)

	if len(matches) >= 4 {
		fmt.Sscanf(matches[1], "%d", &passed)
		fmt.Sscanf(matches[2], "%d", &failed)
		fmt.Sscanf(matches[3], "%d", &ignored)
	}

	return
}
