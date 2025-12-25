package domain

import "context"

// TestScope определяет область тестирования
type TestScope string

const (
	TestScopeAll           TestScope = "all"            // Все тесты
	TestScopeAffected      TestScope = "affected"       // Только затронутые файлы
	TestScopeSmoke         TestScope = "smoke"          // Только smoke тесты
	TestScopeUnit          TestScope = "unit"           // Только unit тесты
	TestScopeIntegration   TestScope = "integration"    // Только integration тесты
	TestScopeAffectedSmoke TestScope = "affected+smoke" // Затронутые + smoke
)

// TestResult представляет результат выполнения теста
type TestResult struct {
	Success  bool                   `json:"success"`
	TestPath string                 `json:"testPath"`
	TestName string                 `json:"testName"`
	Language string                 `json:"language"`
	Duration float64                `json:"duration"`
	Output   string                 `json:"output"`
	Error    string                 `json:"error,omitempty"`
	Coverage *TestCoverage          `json:"coverage,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// TestCoverage представляет покрытие тестами
type TestCoverage struct {
	Percentage float64            `json:"percentage"`
	Lines      int                `json:"lines"`
	Functions  int                `json:"functions"`
	Branches   int                `json:"branches"`
	Files      map[string]float64 `json:"files,omitempty"`
}

// TestSuite представляет набор тестов
type TestSuite struct {
	Name        string      `json:"name"`
	Language    string      `json:"language"`
	ProjectPath string      `json:"projectPath"`
	Tests       []*TestInfo `json:"tests"`
	Config      *TestConfig `json:"config"`
}

// TestInfo представляет информацию о тесте
type TestInfo struct {
	Path        string            `json:"path"`
	Name        string            `json:"name"`
	Type        string            `json:"type"` // "unit", "integration", "smoke"
	TargetFiles []string          `json:"targetFiles,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// TestConfig определяет конфигурацию тестирования
type TestConfig struct {
	Language        string            `json:"language"`
	ProjectPath     string            `json:"projectPath"`
	Scope           TestScope         `json:"scope"`
	Parallel        bool              `json:"parallel"`
	Timeout         int               `json:"timeout"` // в секундах
	Coverage        bool              `json:"coverage"`
	Verbose         bool              `json:"verbose"`
	EnvVars         map[string]string `json:"envVars,omitempty"`
	TestPatterns    []string          `json:"testPatterns,omitempty"`
	ExcludePatterns []string          `json:"excludePatterns,omitempty"`
}

// AffectedGraph представляет граф затронутых файлов
type AffectedGraph struct {
	ChangedFiles  []string            `json:"changedFiles"`
	AffectedFiles []string            `json:"affectedFiles"`
	Dependencies  map[string][]string `json:"dependencies"`
	TestMapping   map[string][]string `json:"testMapping"` // file -> tests
}

// TestEngine определяет интерфейс для выполнения тестов
type TestEngine interface {
	// RunTests выполняет тесты согласно конфигурации
	RunTests(ctx context.Context, config *TestConfig) ([]*TestResult, error)

	// RunTargetedTests выполняет целевые тесты для затронутых файлов
	RunTargetedTests(ctx context.Context, config *TestConfig, affectedGraph *AffectedGraph) ([]*TestResult, error)

	// DiscoverTests обнаруживает тесты в проекте
	DiscoverTests(ctx context.Context, projectPath, language string) (*TestSuite, error)

	// BuildAffectedGraph строит граф затронутых файлов
	BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*AffectedGraph, error)

	// GetTestCoverage получает покрытие тестами
	GetTestCoverage(ctx context.Context, testPath string) (*TestCoverage, error)

	// GetSupportedLanguages возвращает поддерживаемые языки
	GetSupportedLanguages() []string
}

// TestRunner определяет интерфейс для запуска тестов конкретного языка
type TestRunner interface {
	// RunTest выполняет один тест
	RunTest(ctx context.Context, testPath string, config *TestConfig) (*TestResult, error)

	// RunTestSuite выполняет набор тестов
	RunTestSuite(ctx context.Context, suite *TestSuite) ([]*TestResult, error)

	// DiscoverTests обнаруживает тесты
	DiscoverTests(ctx context.Context, projectPath string) ([]*TestInfo, error)

	// GetLanguage возвращает язык
	GetLanguage() string
}

// TestAnalyzer определяет интерфейс для анализа тестов
type TestAnalyzer interface {
	// AnalyzeTestDependencies анализирует зависимости тестов
	AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error)

	// FindTestsForFile находит тесты для файла
	FindTestsForFile(ctx context.Context, filePath string, projectPath string) ([]string, error)

	// IsSmokeTest определяет, является ли тест smoke тестом
	IsSmokeTest(ctx context.Context, testPath string) (bool, error)
}

// TestValidationResult представляет результат валидации тестов
type TestValidationResult struct {
	Success         bool     `json:"success"`
	TotalTests      int      `json:"totalTests"`
	PassedTests     int      `json:"passedTests"`
	FailedTests     int      `json:"failedTests"`
	SkippedTests    int      `json:"skippedTests"`
	SuccessRate     float64  `json:"successRate"`
	TotalDuration   float64  `json:"totalDuration"`
	FailedTestPaths []string `json:"failedTestPaths"`
}
