package testengine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// Test framework constants
const (
	frameworkVitest = "vitest"
	frameworkJest   = "jest"
	frameworkMocha  = "mocha"
)

// TypeScriptTestRunner реализует TestRunner для TypeScript/JavaScript
type TypeScriptTestRunner struct {
	log domain.Logger
}

// NewTypeScriptTestRunner создает новый runner для TypeScript тестов
func NewTypeScriptTestRunner(log domain.Logger) *TypeScriptTestRunner {
	return &TypeScriptTestRunner{
		log: log,
	}
}

// GetLanguage возвращает язык
func (r *TypeScriptTestRunner) GetLanguage() string {
	return "typescript"
}

// RunTest выполняет один TypeScript/JavaScript тест
func (r *TypeScriptTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running TypeScript test: %s", testPath))

	startTime := time.Now()

	// Определяем фреймворк тестирования
	framework := r.detectTestFramework(config.ProjectPath)
	r.log.Info(fmt.Sprintf("Detected test framework: %s", framework))

	// Строим команду для запуска теста
	args, err := r.buildTestCommand(framework, testPath, config)
	if err != nil {
		return nil, fmt.Errorf("failed to build test command: %w", err)
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "npx", args...)
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
		Language: "typescript",
		Duration: duration,
		Output:   string(output),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		r.log.Warning(fmt.Sprintf("TypeScript test failed for %s: %v", testPath, err))
	} else {
		result.Success = true
		r.log.Info(fmt.Sprintf("TypeScript test passed for %s in %.2fs", testPath, duration))
	}

	return result, nil
}

// RunTestSuite выполняет набор TypeScript тестов
func (r *TypeScriptTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	r.log.Info(fmt.Sprintf("Running TypeScript test suite with %d tests", len(suite.Tests)))

	var results []*domain.TestResult

	for _, test := range suite.Tests {
		result, err := r.RunTest(ctx, test.Path, suite.Config)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to run test %s: %v", test.Path, err))
		}
		results = append(results, result)
	}

	r.log.Info(fmt.Sprintf("Completed TypeScript test suite with %d results", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает TypeScript/JavaScript тесты в проекте
func (r *TypeScriptTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	r.log.Info(fmt.Sprintf("Discovering TypeScript tests in: %s", projectPath))

	var tests []*domain.TestInfo

	// Паттерны для поиска тестовых файлов
	testPatterns := []string{
		"*.test.ts", "*.spec.ts",
		"*.test.tsx", "*.spec.tsx",
		"*.test.js", "*.spec.js",
		"*.test.jsx", "*.spec.jsx",
	}

	// Ищем тестовые файлы
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем node_modules, dist, build и скрытые директории
			if r.shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		// Проверяем, соответствует ли файл паттернам тестов
		fileName := info.Name()
		isTestFile := false
		for _, pattern := range testPatterns {
			if matched, _ := filepath.Match(pattern, fileName); matched {
				isTestFile = true
				break
			}
		}

		// Также проверяем директорию __tests__
		if !isTestFile && strings.Contains(path, "__tests__") {
			ext := filepath.Ext(fileName)
			if ext == ".ts" || ext == ".tsx" || ext == ".js" || ext == ".jsx" {
				isTestFile = true
			}
		}

		if !isTestFile {
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
			Name: fileName,
			Type: testType,
			Metadata: map[string]string{
				"framework": r.detectTestFramework(projectPath),
			},
		}

		tests = append(tests, testInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to discover TypeScript tests: %w", err)
	}

	r.log.Info(fmt.Sprintf("Discovered %d TypeScript test files", len(tests)))
	return tests, nil
}

// detectTestFramework определяет фреймворк тестирования по package.json
func (r *TypeScriptTestRunner) detectTestFramework(projectPath string) string {
	packageJSONPath := filepath.Join(projectPath, "package.json")
	content, err := os.ReadFile(packageJSONPath)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Failed to read package.json: %v", err))
		return frameworkJest // По умолчанию используем Jest
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(content, &packageJSON); err != nil {
		r.log.Warning(fmt.Sprintf("Failed to parse package.json: %v", err))
		return frameworkJest
	}

	// Проверяем devDependencies и dependencies
	for _, depsKey := range []string{"devDependencies", "dependencies"} {
		if deps, ok := packageJSON[depsKey].(map[string]interface{}); ok {
			if _, hasVitest := deps["vitest"]; hasVitest {
				return frameworkVitest
			}
			if _, hasJest := deps["jest"]; hasJest {
				return frameworkJest
			}
			if _, hasMocha := deps["mocha"]; hasMocha {
				return frameworkMocha
			}
		}
	}

	// Проверяем наличие конфигурационных файлов
	configFiles := map[string]string{
		"vitest.config.ts":  frameworkVitest,
		"vitest.config.js":  frameworkVitest,
		"vitest.config.mts": frameworkVitest,
		"jest.config.ts":    frameworkJest,
		"jest.config.js":    frameworkJest,
		"jest.config.mjs":   frameworkJest,
		".mocharc.json":     frameworkMocha,
		".mocharc.js":       frameworkMocha,
		".mocharc.yml":      frameworkMocha,
	}

	for configFile, framework := range configFiles {
		if _, err := os.Stat(filepath.Join(projectPath, configFile)); err == nil {
			return framework
		}
	}

	return frameworkJest // По умолчанию используем Jest
}

// buildTestCommand строит команду для запуска теста
func (r *TypeScriptTestRunner) buildTestCommand(framework, testPath string, config *domain.TestConfig) ([]string, error) {
	var args []string

	switch framework {
	case frameworkVitest:
		args = []string{"vitest", "run", testPath}
		if config.Coverage {
			args = append(args, "--coverage")
		}
		args = append(args, "--reporter=json")

	case frameworkJest:
		args = []string{"jest", testPath, "--json"}
		if config.Coverage {
			args = append(args, "--coverage")
		}
		if config.Verbose {
			args = append(args, "--verbose")
		}

	case frameworkMocha:
		args = []string{"mocha", testPath, "--reporter", "json"}
		if config.Timeout > 0 {
			args = append(args, "--timeout", fmt.Sprintf("%d", config.Timeout*1000))
		}

	default:
		return nil, fmt.Errorf("unsupported test framework: %s", framework)
	}

	return args, nil
}

// analyzeTestFile анализирует файл теста для определения его типа
func (r *TypeScriptTestRunner) analyzeTestFile(filePath string) string {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "unit" // По умолчанию считаем unit тестом
	}

	contentStr := strings.ToLower(string(content))

	// Проверяем на smoke тесты
	if strings.Contains(contentStr, "smoke") {
		return "smoke"
	}

	// Проверяем на integration тесты
	integrationIndicators := []string{
		"integration",
		"e2e",
		"end-to-end",
		"api test",
		"database",
		"http",
		"fetch(",
		"axios",
		"supertest",
	}

	for _, indicator := range integrationIndicators {
		if strings.Contains(contentStr, indicator) {
			return "integration"
		}
	}

	// По умолчанию считаем unit тестом
	return "unit"
}

// shouldSkipDir проверяет, нужно ли пропустить директорию
func (r *TypeScriptTestRunner) shouldSkipDir(name string) bool {
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
