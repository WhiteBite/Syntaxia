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

// Build tool constants
const (
	buildToolMaven  = "maven"
	buildToolGradle = "gradle"
)

// JavaTestRunner реализует TestRunner для Java
type JavaTestRunner struct {
	log domain.Logger
}

// NewJavaTestRunner создает новый runner для Java тестов
func NewJavaTestRunner(log domain.Logger) *JavaTestRunner {
	return &JavaTestRunner{
		log: log,
	}
}

// GetLanguage возвращает язык
func (r *JavaTestRunner) GetLanguage() string {
	return "java"
}

// RunTest выполняет один Java тест
func (r *JavaTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Java test: %s", testPath))

	startTime := time.Now()

	// Определяем инструмент сборки
	buildTool := r.detectBuildTool(config.ProjectPath)
	r.log.Info(fmt.Sprintf("Detected build tool: %s", buildTool))

	// Получаем имя тестового класса из пути
	testClassName := r.extractTestClassName(testPath)

	// Строим команду для запуска теста
	cmdName, args := r.buildTestCommand(buildTool, testClassName, config)

	// Создаем команду
	cmd := exec.CommandContext(ctx, cmdName, args...)
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
		Language: "java",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("Java test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("Java test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор Java тестов
func (r *JavaTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running Java test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed Java test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает Java тесты в проекте
func (r *JavaTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering Java tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Стандартные директории для тестов
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
				return err
			}

			if info.IsDir() {
				return nil
			}

			// Проверяем, является ли файл Java тестом
			if !r.isJavaTestFile(info.Name()) {
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
					"buildTool": r.detectBuildTool(projectPath),
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

	r.log.Info(fmt.Sprintf("Discovered %d Java test files", len(tests)))
	return tests, nil
}

// detectBuildTool определяет инструмент сборки по файлам проекта
func (r *JavaTestRunner) detectBuildTool(projectPath string) string {
	// Проверяем наличие pom.xml (Maven)
	if _, err := os.Stat(filepath.Join(projectPath, "pom.xml")); err == nil {
		return buildToolMaven
	}

	// Проверяем наличие build.gradle или build.gradle.kts (Gradle)
	gradleFiles := []string{"build.gradle", "build.gradle.kts"}
	for _, gradleFile := range gradleFiles {
		if _, err := os.Stat(filepath.Join(projectPath, gradleFile)); err == nil {
			return buildToolGradle
		}
	}

	// По умолчанию используем Maven
	return buildToolMaven
}

// buildTestCommand строит команду для запуска теста
func (r *JavaTestRunner) buildTestCommand(buildTool, testClassName string, config *domain.TestConfig) (string, []string) {
	switch buildTool {
	case buildToolMaven:
		args := []string{"test", fmt.Sprintf("-Dtest=%s", testClassName)}
		if config.Verbose {
			args = append(args, "-X")
		}
		return "mvn", args

	case buildToolGradle:
		args := []string{"test", fmt.Sprintf("--tests=%s", testClassName)}
		if config.Verbose {
			args = append(args, "--info")
		}
		// Определяем, использовать ли gradlew
		if r.hasGradleWrapper(config.ProjectPath) {
			return r.getGradleWrapperCommand(), args
		}
		return "gradle", args

	default:
		return "mvn", []string{"test", fmt.Sprintf("-Dtest=%s", testClassName)}
	}
}

// hasGradleWrapper проверяет наличие Gradle wrapper
func (r *JavaTestRunner) hasGradleWrapper(projectPath string) bool {
	wrapperFiles := []string{"gradlew", "gradlew.bat"}
	for _, wrapper := range wrapperFiles {
		if _, err := os.Stat(filepath.Join(projectPath, wrapper)); err == nil {
			return true
		}
	}
	return false
}

// getGradleWrapperCommand возвращает команду для Gradle wrapper в зависимости от ОС
func (r *JavaTestRunner) getGradleWrapperCommand() string {
	if os.PathSeparator == '\\' {
		return ".\\gradlew.bat"
	}
	return "./gradlew"
}

// extractTestClassName извлекает имя класса из пути к файлу
func (r *JavaTestRunner) extractTestClassName(testPath string) string {
	// Убираем расширение .java
	className := strings.TrimSuffix(filepath.Base(testPath), ".java")

	// Пытаемся получить полное имя класса с пакетом
	// Ищем src/test/java в пути
	testPath = strings.ReplaceAll(testPath, "\\", "/")
	parts := strings.Split(testPath, "/")

	startIdx := -1
	for i, part := range parts {
		if part == "java" && i > 0 && parts[i-1] == "test" {
			startIdx = i + 1
			break
		}
	}

	if startIdx > 0 && startIdx < len(parts) {
		// Собираем полное имя класса
		packageParts := parts[startIdx:]
		if len(packageParts) > 0 {
			// Последний элемент - имя файла, убираем расширение
			packageParts[len(packageParts)-1] = strings.TrimSuffix(packageParts[len(packageParts)-1], ".java")
			return strings.Join(packageParts, ".")
		}
	}

	return className
}

// isJavaTestFile проверяет, является ли файл Java тестом
func (r *JavaTestRunner) isJavaTestFile(fileName string) bool {
	if !strings.HasSuffix(fileName, ".java") {
		return false
	}

	// Стандартные паттерны именования тестов
	testPatterns := []string{
		"Test.java",
		"Tests.java",
		"TestCase.java",
		"IT.java", // Integration Test
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
func (r *JavaTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit" // По умолчанию считаем unit тестом
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
		"it.java",
		"@testcontainers",
		"@container",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentStr, indicator) || strings.Contains(fileName, indicator) {
			return "integration"
		}
	}

	// По умолчанию считаем unit тестом
	return "unit"
}

// parseTestAnnotations парсит аннотации тестов из Java файла
func (r *JavaTestRunner) parseTestAnnotations(content string) []string {
	var annotations []string

	// Ищем аннотации @Test, @ParameterizedTest, @RepeatedTest и т.д.
	annotationRegex := regexp.MustCompile(`@(Test|ParameterizedTest|RepeatedTest|DisplayName|Tag)\b`)
	matches := annotationRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 1 {
			annotations = append(annotations, match[1])
		}
	}

	return annotations
}
