package testengine

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

// Test type constants
const (
	testTypeUnit = "unit"
)

/*
Targeted Tests Strategy Documentation

Стратегия выбора целевых тестов основана на анализе графа зависимостей:

1. IMPACTED GRAPH ANALYSIS:
   - Строим symbol graph для анализа зависимостей между файлами
   - Находим прямые зависимости измененных файлов
   - Рекурсивно находим косвенные зависимости (до 3 уровней глубины)
   - Избегаем циклических зависимостей через visited map

2. TARGETED TEST SELECTION:
   - Для каждого затронутого файла находим связанные тесты
   - Приоритизируем smoke тесты для быстрой проверки
   - Исключаем дублирующиеся тесты через testSet
   - Fallback к полному набору тестов при отсутствии targeted

3. SMOKE TEST INTEGRATION:
   - Smoke тесты добавляются в начало списка для приоритетного выполнения
   - Автоматическое обнаружение smoke тестов через анализаторы
   - Поддержка комбинированного режима affected+smoke

4. PERFORMANCE OPTIMIZATION:
   - Ограничение глубины поиска зависимостей (maxDepth = 3)
   - Кэширование результатов анализа
   - Параллельное выполнение независимых тестов

5. FALLBACK STRATEGY:
   - При ошибках анализа зависимостей -> полный набор тестов
   - При отсутствии targeted тестов -> все тесты
   - При ошибках smoke тестов -> только affected тесты
*/

// TestEngineImpl реализует TestEngine
type TestEngineImpl struct {
	log           domain.Logger
	testRunners   map[string]domain.TestRunner
	testAnalyzers map[string]domain.TestAnalyzer
	symbolGraph   domain.SymbolGraphBuilder
}

// NewTestEngine создает новый движок тестирования
func NewTestEngine(log domain.Logger, symbolGraph domain.SymbolGraphBuilder) *TestEngineImpl {
	return &TestEngineImpl{
		log:           log,
		testRunners:   make(map[string]domain.TestRunner),
		testAnalyzers: make(map[string]domain.TestAnalyzer),
		symbolGraph:   symbolGraph,
	}
}

// RegisterTestRunner регистрирует runner для языка
func (e *TestEngineImpl) RegisterTestRunner(language string, runner domain.TestRunner) {
	e.testRunners[language] = runner
	e.log.Info(fmt.Sprintf("Registered test runner for language: %s", language))
}

// RegisterTestAnalyzer регистрирует анализатор для языка
func (e *TestEngineImpl) RegisterTestAnalyzer(language string, analyzer domain.TestAnalyzer) {
	e.testAnalyzers[language] = analyzer
	e.log.Info(fmt.Sprintf("Registered test analyzer for language: %s", language))
}

// RunTests выполняет тесты согласно конфигурации
func (e *TestEngineImpl) RunTests(ctx context.Context, config *domain.TestConfig) ([]*domain.TestResult, error) {
	e.log.Info(fmt.Sprintf("Running tests for language: %s, scope: %s", config.Language, config.Scope))

	runner, exists := e.testRunners[config.Language]
	if !exists {
		return nil, fmt.Errorf("no test runner registered for language: %s", config.Language)
	}

	// Обнаруживаем тесты
	suite, err := e.DiscoverTests(ctx, config.ProjectPath, config.Language)
	if err != nil {
		return nil, fmt.Errorf("failed to discover tests: %w", err)
	}

	// Фильтруем тесты согласно scope
	filteredTests := e.filterTestsByScope(suite.Tests, config.Scope)

	// Создаем новый suite с отфильтрованными тестами
	filteredSuite := &domain.TestSuite{
		Name:        suite.Name,
		Language:    suite.Language,
		ProjectPath: suite.ProjectPath,
		Tests:       filteredTests,
		Config:      config,
	}

	// Запускаем тесты
	results, err := runner.RunTestSuite(ctx, filteredSuite)
	if err != nil {
		return nil, fmt.Errorf("failed to run test suite: %w", err)
	}

	e.log.Info(fmt.Sprintf("Completed %d tests", len(results)))
	return results, nil
}

// collectTargetedTests collects tests for affected files
func (e *TestEngineImpl) collectTargetedTests(ctx context.Context, config *domain.TestConfig, affectedGraph *domain.AffectedGraph, analyzer domain.TestAnalyzer) ([]string, map[string]bool) {
	var targetTests []string
	testSet := make(map[string]bool)

	for _, file := range affectedGraph.AffectedFiles {
		tests, err := analyzer.FindTestsForFile(ctx, file, config.ProjectPath)
		if err != nil {
			e.log.Warning(fmt.Sprintf("Failed to find tests for file %s: %v", file, err))
			continue
		}
		for _, test := range tests {
			if !testSet[test] {
				targetTests = append(targetTests, test)
				testSet[test] = true
			}
		}
	}
	return targetTests, testSet
}

// addSmokeTests adds smoke tests to the beginning of test list
func (e *TestEngineImpl) addSmokeTests(ctx context.Context, config *domain.TestConfig, targetTests []string, testSet map[string]bool) []string {
	if config.Scope != domain.TestScopeAffectedSmoke && config.Scope != domain.TestScopeSmoke {
		return targetTests
	}

	smokeTests, err := e.findSmokeTests(ctx, config.ProjectPath, config.Language)
	if err != nil {
		e.log.Warning(fmt.Sprintf("Failed to find smoke tests: %v", err))
		return targetTests
	}

	e.log.Info(fmt.Sprintf("Found %d smoke tests", len(smokeTests)))
	for _, test := range smokeTests {
		if !testSet[test] {
			targetTests = append([]string{test}, targetTests...)
			testSet[test] = true
		}
	}
	return targetTests
}

// RunTargetedTests выполняет целевые тесты для затронутых файлов
func (e *TestEngineImpl) RunTargetedTests(ctx context.Context, config *domain.TestConfig, affectedGraph *domain.AffectedGraph) ([]*domain.TestResult, error) {
	e.log.Info(fmt.Sprintf("Running targeted tests for %d affected files", len(affectedGraph.AffectedFiles)))

	runner, exists := e.testRunners[config.Language]
	if !exists {
		return nil, fmt.Errorf("no test runner registered for language: %s", config.Language)
	}

	analyzer, exists := e.testAnalyzers[config.Language]
	if !exists {
		e.log.Warning(fmt.Sprintf("No test analyzer for language: %s, falling back to all tests", config.Language))
		return e.RunTests(ctx, config)
	}

	targetTests, testSet := e.collectTargetedTests(ctx, config, affectedGraph, analyzer)
	targetTests = e.addSmokeTests(ctx, config, targetTests, testSet)

	if len(targetTests) == 0 {
		e.log.Warning("No targeted tests found, running all tests")
		return e.RunTests(ctx, config)
	}

	e.log.Info(fmt.Sprintf("Found %d targeted tests", len(targetTests)))

	results := make([]*domain.TestResult, 0, len(targetTests))
	for _, testPath := range targetTests {
		result, err := runner.RunTest(ctx, testPath, config)
		if err != nil {
			e.log.Warning(fmt.Sprintf("Failed to run test %s: %v", testPath, err))
			result = &domain.TestResult{Success: false, TestPath: testPath, Error: err.Error()}
		}
		results = append(results, result)
	}

	e.log.Info(fmt.Sprintf("Completed %d targeted tests", len(results)))
	return results, nil
}

// DiscoverTests обнаруживает тесты в проекте
func (e *TestEngineImpl) DiscoverTests(ctx context.Context, projectPath, language string) (*domain.TestSuite, error) {
	runner, exists := e.testRunners[language]
	if !exists {
		return nil, fmt.Errorf("no test runner registered for language: %s", language)
	}

	tests, err := runner.DiscoverTests(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to discover tests: %w", err)
	}

	suite := &domain.TestSuite{
		Name:        filepath.Base(projectPath),
		Language:    language,
		ProjectPath: projectPath,
		Tests:       tests,
	}

	e.log.Info(fmt.Sprintf("Discovered %d tests for language: %s", len(tests), language))
	return suite, nil
}

// BuildAffectedGraph строит граф затронутых файлов
// Стратегия: используем symbol graph для анализа зависимостей и находим все файлы,
// которые могут быть затронуты изменениями, включая прямые и косвенные зависимости
func (e *TestEngineImpl) BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*domain.AffectedGraph, error) {
	e.log.Info(fmt.Sprintf("Building affected graph for %d changed files", len(changedFiles)))

	// Строим граф символов для анализа зависимостей
	symbolGraph, err := e.symbolGraph.BuildGraph(ctx, projectPath)
	if err != nil {
		e.log.Warning(fmt.Sprintf("Failed to build symbol graph: %v", err))
		// Возвращаем простой граф только с измененными файлами
		return &domain.AffectedGraph{
			ChangedFiles:  changedFiles,
			AffectedFiles: changedFiles,
			Dependencies:  make(map[string][]string),
			TestMapping:   make(map[string][]string),
		}, nil
	}

	// Находим затронутые файлы через зависимости
	affectedFiles := make(map[string]bool)
	dependencies := make(map[string][]string)
	testMapping := make(map[string][]string)

	// Добавляем измененные файлы
	for _, file := range changedFiles {
		affectedFiles[file] = true
	}

	// Находим зависимости для каждого измененного файла с улучшенным алгоритмом
	for _, file := range changedFiles {
		fileDeps, err := e.findFileDependencies(file, symbolGraph)
		if err != nil {
			e.log.Warning(fmt.Sprintf("Failed to find dependencies for %s: %v", file, err))
			continue
		}

		dependencies[file] = fileDeps

		// Добавляем прямые зависимости
		for _, dep := range fileDeps {
			affectedFiles[dep] = true
		}

		// Находим косвенные зависимости (зависимости зависимостей)
		indirectDeps := e.findIndirectDependencies(fileDeps, symbolGraph)
		for _, dep := range indirectDeps {
			affectedFiles[dep] = true
		}
	}

	// Конвертируем map в slice
	affectedFilesList := make([]string, 0, len(affectedFiles))
	for file := range affectedFiles {
		affectedFilesList = append(affectedFilesList, file)
	}

	// Находим тесты для затронутых файлов
	for _, file := range affectedFilesList {
		tests, err := e.findTestsForFile(ctx, file, projectPath)
		if err != nil {
			e.log.Warning(fmt.Sprintf("Failed to find tests for %s: %v", file, err))
			continue
		}
		testMapping[file] = tests
	}

	graph := &domain.AffectedGraph{
		ChangedFiles:  changedFiles,
		AffectedFiles: affectedFilesList,
		Dependencies:  dependencies,
		TestMapping:   testMapping,
	}

	e.log.Info(fmt.Sprintf("Built affected graph with %d affected files", len(affectedFilesList)))
	return graph, nil
}

// GetTestCoverage получает покрытие тестами
func (e *TestEngineImpl) GetTestCoverage(ctx context.Context, testPath string) (*domain.TestCoverage, error) {
	// Базовая реализация - в будущем можно добавить интеграцию с инструментами покрытия
	return &domain.TestCoverage{
		Percentage: 0.0,
		Lines:      0,
		Functions:  0,
		Branches:   0,
		Files:      make(map[string]float64),
	}, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (e *TestEngineImpl) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(e.testRunners))
	for lang := range e.testRunners {
		languages = append(languages, lang)
	}
	return languages
}

// filterTestsByScope фильтрует тесты согласно scope
func (e *TestEngineImpl) filterTestsByScope(tests []*domain.TestInfo, scope domain.TestScope) []*domain.TestInfo {
	var filtered []*domain.TestInfo

	for _, test := range tests {
		switch scope {
		case domain.TestScopeAll:
			filtered = append(filtered, test)
		case domain.TestScopeSmoke:
			if test.Type == "smoke" {
				filtered = append(filtered, test)
			}
		case domain.TestScopeUnit:
			if test.Type == testTypeUnit {
				filtered = append(filtered, test)
			}
		case domain.TestScopeIntegration:
			if test.Type == "integration" {
				filtered = append(filtered, test)
			}
		case domain.TestScopeAffected, domain.TestScopeAffectedSmoke:
			// Для affected scope фильтрация происходит в RunTargetedTests
			filtered = append(filtered, test)
		}
	}

	return filtered
}

// findSmokeTests находит smoke тесты
func (e *TestEngineImpl) findSmokeTests(ctx context.Context, projectPath, language string) ([]string, error) {
	analyzer, exists := e.testAnalyzers[language]
	if !exists {
		return nil, fmt.Errorf("no test analyzer for language: %s", language)
	}

	// Обнаруживаем все тесты
	runner, exists := e.testRunners[language]
	if !exists {
		return nil, fmt.Errorf("no test runner for language: %s", language)
	}

	tests, err := runner.DiscoverTests(ctx, projectPath)
	if err != nil {
		return nil, err
	}

	var smokeTests []string
	for _, test := range tests {
		isSmoke, err := analyzer.IsSmokeTest(ctx, test.Path)
		if err != nil {
			continue
		}
		if isSmoke {
			smokeTests = append(smokeTests, test.Path)
		}
	}

	return smokeTests, nil
}

// buildNodeIndex creates a map of node ID to node for fast lookup
func buildNodeIndex(graph *domain.SymbolGraph) map[string]*domain.SymbolNode {
	index := make(map[string]*domain.SymbolNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		index[node.ID] = node
	}
	return index
}

// findFileDependencies находит зависимости файла через граф символов
func (e *TestEngineImpl) findFileDependencies(filePath string, graph *domain.SymbolGraph) ([]string, error) {
	dependencySet := make(map[string]bool)
	nodeIndex := buildNodeIndex(graph)

	// Collect node IDs for this file
	fileNodeIDs := make(map[string]bool)
	for _, node := range graph.Nodes {
		if node.Path == filePath {
			fileNodeIDs[node.ID] = true
		}
	}

	// Find dependencies via edges
	for _, edge := range graph.Edges {
		if !fileNodeIDs[edge.From] {
			continue
		}
		if targetNode, ok := nodeIndex[edge.To]; ok && targetNode.Path != filePath {
			dependencySet[targetNode.Path] = true
		}
	}

	dependencies := make([]string, 0, len(dependencySet))
	for dep := range dependencySet {
		dependencies = append(dependencies, dep)
	}
	return dependencies, nil
}

// findIndirectDependencies находит косвенные зависимости (зависимости зависимостей)
// Стратегия: рекурсивно проходим по графу зависимостей на глубину до 3 уровней
// для избежания бесконечных циклов и оптимизации производительности
func (e *TestEngineImpl) findIndirectDependencies(directDeps []string, symbolGraph *domain.SymbolGraph) []string {
	if symbolGraph == nil {
		return []string{}
	}

	indirectDeps := make(map[string]bool)
	visited := make(map[string]bool)
	maxDepth := 3

	// Рекурсивная функция для поиска зависимостей
	var findDepsRecursive func(files []string, depth int)
	findDepsRecursive = func(files []string, depth int) {
		if depth >= maxDepth {
			return
		}

		for _, file := range files {
			if visited[file] {
				continue
			}
			visited[file] = true

			// Находим зависимости для текущего файла
			deps, err := e.findFileDependencies(file, symbolGraph)
			if err != nil {
				continue
			}
			for _, dep := range deps {
				if !visited[dep] {
					indirectDeps[dep] = true
					findDepsRecursive([]string{dep}, depth+1)
				}
			}
		}
	}

	// Запускаем рекурсивный поиск
	findDepsRecursive(directDeps, 1)

	// Конвертируем в слайс
	result := make([]string, 0, len(indirectDeps))
	for dep := range indirectDeps {
		result = append(result, dep)
	}

	return result
}

// findTestsForFile находит тесты для файла
func (e *TestEngineImpl) findTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	// Определяем язык по расширению файла
	ext := strings.ToLower(filepath.Ext(filePath))
	var language string

	switch ext {
	case ".go":
		language = "go"
	case ".ts", ".tsx":
		language = "typescript"
	case ".js", ".jsx":
		language = "javascript"
	case ".java":
		language = "java"
	case ".py":
		language = "python"
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	analyzer, exists := e.testAnalyzers[language]
	if !exists {
		return nil, fmt.Errorf("no test analyzer for language: %s", language)
	}

	return analyzer.FindTestsForFile(ctx, filePath, projectPath)
}
