package domain

import "context"

// StaticAnalyzerType определяет тип статического анализатора
type StaticAnalyzerType string

const (
	StaticAnalyzerTypeStaticcheck StaticAnalyzerType = "staticcheck" // Go
	StaticAnalyzerTypeESLint      StaticAnalyzerType = "eslint"      // TypeScript/JavaScript
	StaticAnalyzerTypeErrorProne  StaticAnalyzerType = "errorprone"  // Java
	StaticAnalyzerTypeRuff        StaticAnalyzerType = "ruff"        // Python
	StaticAnalyzerTypeClangTidy   StaticAnalyzerType = "clang-tidy"  // C/C++
)

// StaticIssue представляет проблему, найденную статическим анализатором
type StaticIssue struct {
	File        string   `json:"file"`
	Line        int      `json:"line"`
	Column      int      `json:"column"`
	Severity    string   `json:"severity"` // "error", "warning", "info", "hint"
	Message     string   `json:"message"`
	Code        string   `json:"code,omitempty"`
	Category    string   `json:"category,omitempty"`
	Confidence  string   `json:"confidence,omitempty"` // "high", "medium", "low"
	Suggestions []string `json:"suggestions,omitempty"`
}

// StaticAnalysisResult представляет результат статического анализа
type StaticAnalysisResult struct {
	Success     bool                   `json:"success"`
	Language    string                 `json:"language"`
	ProjectPath string                 `json:"projectPath"`
	Analyzer    StaticAnalyzerType     `json:"analyzer"`
	Issues      []*StaticIssue         `json:"issues"`
	Summary     *StaticAnalysisSummary `json:"summary"`
	Duration    float64                `json:"duration"`
	Error       string                 `json:"error,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// StaticAnalysisSummary представляет сводку статического анализа
type StaticAnalysisSummary struct {
	TotalIssues       int            `json:"totalIssues"`
	ErrorCount        int            `json:"errorCount"`
	WarningCount      int            `json:"warningCount"`
	InfoCount         int            `json:"infoCount"`
	HintCount         int            `json:"hintCount"`
	SeverityBreakdown map[string]int `json:"severityBreakdown"`
	CategoryBreakdown map[string]int `json:"categoryBreakdown"`
	FilesAnalyzed     int            `json:"filesAnalyzed"`
	FilesWithIssues   int            `json:"filesWithIssues"`
}

// StaticAnalyzerConfig определяет конфигурацию статического анализатора
type StaticAnalyzerConfig struct {
	Language     string             `json:"language"`
	ProjectPath  string             `json:"projectPath"`
	Analyzer     StaticAnalyzerType `json:"analyzer"`
	ConfigFile   string             `json:"configFile,omitempty"`
	Rules        []string           `json:"rules,omitempty"`
	ExcludeRules []string           `json:"excludeRules,omitempty"`
	Severity     string             `json:"severity"` // "error", "warning", "info", "hint"
	Timeout      int                `json:"timeout"`  // в секундах
	Parallel     bool               `json:"parallel"`
	OutputFormat string             `json:"outputFormat"` // "json", "text", "sarif"
	EnvVars      map[string]string  `json:"envVars,omitempty"`
}

// StaticAnalyzer определяет интерфейс для статического анализатора
type StaticAnalyzer interface {
	// Analyze выполняет статический анализ
	Analyze(ctx context.Context, config *StaticAnalyzerConfig) (*StaticAnalysisResult, error)

	// GetSupportedLanguages возвращает поддерживаемые языки
	GetSupportedLanguages() []string

	// GetAnalyzerType возвращает тип анализатора
	GetAnalyzerType() StaticAnalyzerType

	// ValidateConfig проверяет корректность конфигурации
	ValidateConfig(config *StaticAnalyzerConfig) error
}

// StaticAnalyzerEngine определяет интерфейс для движка статического анализа
type StaticAnalyzerEngine interface {
	// RegisterAnalyzer регистрирует анализатор
	RegisterAnalyzer(analyzer StaticAnalyzer)

	// AnalyzeProject выполняет анализ проекта
	AnalyzeProject(ctx context.Context, projectPath string, languages []string) (map[string]*StaticAnalysisResult, error)

	// AnalyzeFile выполняет анализ одного файла
	AnalyzeFile(ctx context.Context, filePath string, config *StaticAnalyzerConfig) (*StaticAnalysisResult, error)

	// GetSupportedAnalyzers возвращает поддерживаемые анализаторы
	GetSupportedAnalyzers() []StaticAnalyzerType

	// GetAnalyzerForLanguage возвращает анализатор для языка
	GetAnalyzerForLanguage(language string) (StaticAnalyzer, error)

	// GenerateReport генерирует отчет о статическом анализе
	GenerateReport(results map[string]*StaticAnalysisResult, projectPath string) *StaticAnalysisReport
}

// StaticAnalysisReport представляет отчет о статическом анализе
type StaticAnalysisReport struct {
	ProjectPath     string                           `json:"projectPath"`
	Timestamp       string                           `json:"timestamp"`
	TotalDuration   float64                          `json:"totalDuration"`
	Results         map[string]*StaticAnalysisResult `json:"results"`
	Summary         *StaticAnalysisReportSummary     `json:"summary"`
	Recommendations []string                         `json:"recommendations,omitempty"`
}

// StaticAnalysisReportSummary представляет сводку отчета
type StaticAnalysisReportSummary struct {
	TotalIssues       int            `json:"totalIssues"`
	TotalErrors       int            `json:"totalErrors"`
	TotalWarnings     int            `json:"totalWarnings"`
	LanguagesAnalyzed []string       `json:"languagesAnalyzed"`
	AnalyzersUsed     []string       `json:"analyzersUsed"`
	CriticalIssues    []*StaticIssue `json:"criticalIssues,omitempty"`
	Success           bool           `json:"success"`
}

// LanguageAnalysisValidation represents validation result for a specific language
type LanguageAnalysisValidation struct {
	Language   string `json:"language"`
	Success    bool   `json:"success"`
	IssueCount int    `json:"issueCount"`
	Error      string `json:"error,omitempty"`
}

// StaticAnalysisValidationResult представляет результат валидации статического анализа
type StaticAnalysisValidationResult struct {
	Success         bool                                   `json:"success"`
	TotalLanguages  int                                    `json:"totalLanguages"`
	SuccessCount    int                                    `json:"successCount"`
	FailureCount    int                                    `json:"failureCount"`
	SuccessRate     float64                                `json:"successRate"`
	TotalIssues     int                                    `json:"totalIssues"`
	TotalErrors     int                                    `json:"totalErrors"`
	TotalWarnings   int                                    `json:"totalWarnings"`
	FailedLanguages []string                               `json:"failedLanguages"`
	CriticalIssues  []*StaticIssue                         `json:"criticalIssues"`
	Languages       map[string]*LanguageAnalysisValidation `json:"languages"`
}
