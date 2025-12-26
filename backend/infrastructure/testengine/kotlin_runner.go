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

// Kotlin test runner constants
const (
	kotlinTestTimeout     = 300 // 5 minutes default timeout
	kotlinTestFileSuffix  = "Test.kt"
	kotlinTestDirName     = "test"
)

// KotlinTestRunner реализует TestRunner для Kotlin
type KotlinTestRunner struct {
	log domain.Logger
}

// NewKotlinTestRunner создает новый runner для Kotlin тестов
func NewKotlinTestRunner(log domain.Logger) *KotlinTestRunner {
	return &KotlinTestRunner{
		log: log,
	}
}

// GetLanguage возвращает язык
func (r *KotlinTestRunner) GetLanguage() string {
	return "kotlin"
}

// RunTest выполняет один Kotlin тест
func (r *KotlinTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Kotlin test: %s", testPath))

	startTime := time.Now()

	// Получаем имя тестового класса из пути
	testClassName := r.extractTestClassName(testPath)

	// Определяем команду gradle
	gradleCmd := r.getGradleCommand(config.ProjectPath)
	if gradleCmd == "" {
		return &domain.TestResult{
			TestPath: testPath,
			TestName: testClassName,
			Language: "kotlin",
			Success:  false,
			Error:    "gradle not found in PATH",
		}, nil
	}

	// Строим команду для запуска теста
	args := []string{"test", fmt.Sprintf("--tests=%s", testClassName)}
	if config.Verbose {
		args = append(args, "--info")
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, gradleCmd, args...)
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
		TestName: testClassName,
		Language: "kotlin",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Kotlin test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Kotlin test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор Kotlin тестов
func (r *KotlinTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Kotlin test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed Kotlin test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает Kotlin тесты в проекте
func (r *KotlinTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Kotlin tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Стандартные директории для тестов
	testDirs := []string{
		"src/test/kotlin",
		"src/test/java", // Kotlin tests can also be in java directory
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
				return err
			}

			if info.IsDir() {
				return nil
			}

			// Проверяем, является ли файл Kotlin тестом
			if !r.isKotlinTestFile(info.Name()) {
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
					"className": r.extractTestClassName(relPath),
				},
			}

			tests = append(tests, testInfo)
			return nil
		})

		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to walk test directory %s: %v", testDir, err))
		}
	}

	r.log.Info(fmt.Sprintf("Discovered %d Kotlin test files", len(tests)))
	return tests, nil
}

// getGradleCommand возвращает команду для запуска Gradle
func (r *KotlinTestRunner) getGradleCommand(projectPath string) string {
	// Проверяем наличие gradlew wrapper
	if r.hasGradleWrapper(projectPath) {
		if os.PathSeparator == '\\' {
			return filepath.Join(projectPath, "gradlew.bat")
		}
		return filepath.Join(projectPath, "gradlew")
	}

	// Используем системный gradle
	if _, err := exec.LookPath("gradle"); err == nil {
		return "gradle"
	}

	return ""
}

// hasGradleWrapper проверяет наличие Gradle wrapper
func (r *KotlinTestRunner) hasGradleWrapper(projectPath string) bool {
	wrapperFiles := []string{"gradlew", "gradlew.bat"}
	for _, wrapper := range wrapperFiles {
		if _, err := os.Stat(filepath.Join(projectPath, wrapper)); err == nil {
			return true
		}
	}
	return false
}

// extractTestClassName извлекает имя класса из пути к файлу
func (r *KotlinTestRunner) extractTestClassName(testPath string) string {
	// Убираем расширение .kt
	className := strings.TrimSuffix(filepath.Base(testPath), ".kt")

	// Пытаемся получить полное имя класса с пакетом
	testPath = strings.ReplaceAll(testPath, "\\", "/")
	parts := strings.Split(testPath, "/")

	startIdx := -1
	for i, part := range parts {
		if part == "kotlin" && i > 0 && parts[i-1] == "test" {
			startIdx = i + 1
			break
		}
		if part == "java" && i > 0 && parts[i-1] == "test" {
			startIdx = i + 1
			break
		}
	}

	if startIdx > 0 && startIdx < len(parts) {
		packageParts := parts[startIdx:]
		if len(packageParts) > 0 {
			packageParts[len(packageParts)-1] = strings.TrimSuffix(packageParts[len(packageParts)-1], ".kt")
			return strings.Join(packageParts, ".")
		}
	}

	return className
}

// isKotlinTestFile проверяет, является ли файл Kotlin тестом
func (r *KotlinTestRunner) isKotlinTestFile(fileName string) bool {
	if !strings.HasSuffix(fileName, ".kt") {
		return false
	}

	// Стандартные паттерны именования тестов
	testPatterns := []string{
		"Test.kt",
		"Tests.kt",
		"Spec.kt",
		"IT.kt", // Integration Test
	}

	for _, pattern := range testPatterns {
		if strings.HasSuffix(fileName, pattern) {
			return true
		}
	}

	// Также проверяем префикс Test
	if strings.HasPrefix(fileName, "Test") {
		return true
	}

	return false
}

// analyzeTestFile анализирует файл теста для определения его типа
func (r *KotlinTestRunner) analyzeTestFile(filePath string) string {
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
		"@springboottest",
		"@datajpatest",
		"@webmvctest",
		"@integrationtest",
		"integration",
		"it.kt",
		"@testcontainers",
		"@container",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentStr, indicator) || strings.Contains(fileName, indicator) {
			return "integration"
		}
	}

	return "unit"
}

// ParseTestOutput парсит вывод gradle test для извлечения результатов
func (r *KotlinTestRunner) ParseTestOutput(output string) (passed, failed, skipped int) {
	// Формат: "5 tests completed, 2 failed, 1 skipped"
	re := regexp.MustCompile(`(\d+)\s+tests?\s+completed.*?(\d+)\s+failed.*?(\d+)\s+skipped`)
	matches := re.FindStringSubmatch(output)

	if len(matches) >= 4 {
		fmt.Sscanf(matches[1], "%d", &passed)
		fmt.Sscanf(matches[2], "%d", &failed)
		fmt.Sscanf(matches[3], "%d", &skipped)
		passed = passed - failed // completed includes failed
	}

	return
}
