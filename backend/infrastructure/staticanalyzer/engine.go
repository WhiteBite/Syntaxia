package staticanalyzer

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
	"time"
)

// Severity level constants
const (
	severityError   = "error"
	severityWarning = "warning"
	severityInfo    = "info"
	severityHint    = "hint"
	severityOther   = "other"
)

// generateSummary создает сводку анализа из списка issues (общая функция для всех анализаторов)
func generateSummary(issues []*domain.StaticIssue) *domain.StaticAnalysisSummary {
	summary := &domain.StaticAnalysisSummary{
		TotalIssues:       len(issues),
		SeverityBreakdown: make(map[string]int),
		CategoryBreakdown: make(map[string]int),
		FilesAnalyzed:     0,
		FilesWithIssues:   0,
	}

	filesWithIssues := make(map[string]bool)

	for _, issue := range issues {
		summary.SeverityBreakdown[issue.Severity]++
		summary.CategoryBreakdown[issue.Category]++
		filesWithIssues[issue.File] = true

		switch issue.Severity {
		case severityError:
			summary.ErrorCount++
		case severityWarning:
			summary.WarningCount++
		case severityInfo:
			summary.InfoCount++
		case severityHint:
			summary.HintCount++
		}
	}

	summary.FilesWithIssues = len(filesWithIssues)
	return summary
}

// StaticAnalyzerEngineImpl реализует StaticAnalyzerEngine
type StaticAnalyzerEngineImpl struct {
	log         domain.Logger
	analyzers   map[string]domain.StaticAnalyzer
	languageMap map[string]domain.StaticAnalyzerType
}

// NewStaticAnalyzerEngine создает новый движок статического анализа
func NewStaticAnalyzerEngine(log domain.Logger) *StaticAnalyzerEngineImpl {
	engine := &StaticAnalyzerEngineImpl{
		log:         log,
		analyzers:   make(map[string]domain.StaticAnalyzer),
		languageMap: make(map[string]domain.StaticAnalyzerType),
	}

	// Устанавливаем соответствие языков и анализаторов
	engine.languageMap["go"] = domain.StaticAnalyzerTypeStaticcheck
	engine.languageMap["typescript"] = domain.StaticAnalyzerTypeESLint
	engine.languageMap["ts"] = domain.StaticAnalyzerTypeESLint
	engine.languageMap["javascript"] = domain.StaticAnalyzerTypeESLint
	engine.languageMap["js"] = domain.StaticAnalyzerTypeESLint
	engine.languageMap["java"] = domain.StaticAnalyzerTypeErrorProne
	engine.languageMap["python"] = domain.StaticAnalyzerTypeRuff
	engine.languageMap["py"] = domain.StaticAnalyzerTypeRuff
	engine.languageMap["c"] = domain.StaticAnalyzerTypeClangTidy
	engine.languageMap["cpp"] = domain.StaticAnalyzerTypeClangTidy
	engine.languageMap["cc"] = domain.StaticAnalyzerTypeClangTidy

	return engine
}

// RegisterAnalyzer регистрирует анализатор
func (e *StaticAnalyzerEngineImpl) RegisterAnalyzer(analyzer domain.StaticAnalyzer) {
	analyzerType := analyzer.GetAnalyzerType()
	e.analyzers[string(analyzerType)] = analyzer
}

// AnalyzeProject выполняет анализ проекта
func (e *StaticAnalyzerEngineImpl) AnalyzeProject(ctx context.Context, projectPath string, languages []string) (map[string]*domain.StaticAnalysisResult, error) {
	e.log.Info(fmt.Sprintf("Analyzing project: %s for languages: %v", projectPath, languages))

	results := make(map[string]*domain.StaticAnalysisResult)

	for _, language := range languages {
		analyzer, err := e.GetAnalyzerForLanguage(language)
		if err != nil {
			e.log.Warning(fmt.Sprintf("No analyzer for language %s: %v", language, err))
			continue
		}

		config := &domain.StaticAnalyzerConfig{
			Language:     language,
			ProjectPath:  projectPath,
			Analyzer:     analyzer.GetAnalyzerType(),
			Timeout:      300, // 5 минут таймаут
			OutputFormat: "json",
		}

		result, err := analyzer.Analyze(ctx, config)
		if err != nil {
			e.log.Warning(fmt.Sprintf("Analysis failed for language %s: %v", language, err))
			result = &domain.StaticAnalysisResult{
				Success:     false,
				Language:    language,
				ProjectPath: projectPath,
				Analyzer:    analyzer.GetAnalyzerType(),
				Error:       err.Error(),
			}
		}

		results[language] = result
	}

	e.log.Info(fmt.Sprintf("Completed analysis for %d languages", len(results)))
	return results, nil
}

// AnalyzeFile выполняет анализ одного файла
func (e *StaticAnalyzerEngineImpl) AnalyzeFile(ctx context.Context, filePath string, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	e.log.Info(fmt.Sprintf("Analyzing file: %s", filePath))

	analyzer, err := e.GetAnalyzerForLanguage(config.Language)
	if err != nil {
		return nil, fmt.Errorf("no analyzer for language %s: %w", config.Language, err)
	}

	// Валидируем конфигурацию
	if err := analyzer.ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	startTime := time.Now()
	result, err := analyzer.Analyze(ctx, config)
	duration := time.Since(startTime).Seconds()

	if result != nil {
		result.Duration = duration
	}

	if err != nil {
		e.log.Warning(fmt.Sprintf("Analysis failed for file %s: %v", filePath, err))
		if result == nil {
			result = &domain.StaticAnalysisResult{
				Success:     false,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    config.Analyzer,
				Error:       err.Error(),
				Duration:    duration,
			}
		}
	}

	return result, nil
}

// GetSupportedAnalyzers возвращает поддерживаемые анализаторы
func (e *StaticAnalyzerEngineImpl) GetSupportedAnalyzers() []domain.StaticAnalyzerType {
	analyzers := make([]domain.StaticAnalyzerType, 0, len(e.analyzers))
	for analyzerType := range e.analyzers {
		analyzers = append(analyzers, domain.StaticAnalyzerType(analyzerType))
	}
	return analyzers
}

// GetAnalyzerForLanguage возвращает анализатор для языка
func (e *StaticAnalyzerEngineImpl) GetAnalyzerForLanguage(language string) (domain.StaticAnalyzer, error) {
	analyzerType, exists := e.languageMap[strings.ToLower(language)]
	if !exists {
		return nil, fmt.Errorf("no analyzer mapped for language: %s", language)
	}

	analyzer, exists := e.analyzers[string(analyzerType)]
	if !exists {
		return nil, fmt.Errorf("analyzer %s not registered", analyzerType)
	}

	return analyzer, nil
}

// GenerateReport генерирует отчет о статическом анализе
func (e *StaticAnalyzerEngineImpl) GenerateReport(results map[string]*domain.StaticAnalysisResult, projectPath string) *domain.StaticAnalysisReport {
	report := &domain.StaticAnalysisReport{
		ProjectPath:   projectPath,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		TotalDuration: 0.0,
		Results:       results,
		Summary:       &domain.StaticAnalysisReportSummary{},
	}

	var totalIssues, totalErrors, totalWarnings int
	var languagesAnalyzed []string
	var analyzersUsed []string
	var criticalIssues []*domain.StaticIssue

	for language, result := range results {
		if result.Success {
			languagesAnalyzed = append(languagesAnalyzed, language)
			analyzersUsed = append(analyzersUsed, string(result.Analyzer))
			report.TotalDuration += result.Duration

			if result.Summary != nil {
				totalIssues += result.Summary.TotalIssues
				totalErrors += result.Summary.ErrorCount
				totalWarnings += result.Summary.WarningCount

				// Добавляем критические проблемы (ошибки)
				for _, issue := range result.Issues {
					if issue.Severity == severityError {
						criticalIssues = append(criticalIssues, issue)
					}
				}
			}
		}
	}

	report.Summary.TotalIssues = totalIssues
	report.Summary.TotalErrors = totalErrors
	report.Summary.TotalWarnings = totalWarnings
	report.Summary.LanguagesAnalyzed = languagesAnalyzed
	report.Summary.AnalyzersUsed = analyzersUsed
	report.Summary.CriticalIssues = criticalIssues
	report.Summary.Success = totalErrors == 0

	// Генерируем рекомендации
	report.Recommendations = e.generateRecommendations(report)

	return report
}

// generateRecommendations генерирует рекомендации на основе результатов анализа
func (e *StaticAnalyzerEngineImpl) generateRecommendations(report *domain.StaticAnalysisReport) []string {
	var recommendations []string

	if report.Summary.TotalErrors > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Fix %d critical errors before proceeding", report.Summary.TotalErrors))
	}

	if report.Summary.TotalWarnings > 0 {
		recommendations = append(recommendations, fmt.Sprintf("Consider addressing %d warnings to improve code quality", report.Summary.TotalWarnings))
	}

	if len(report.Summary.CriticalIssues) > 0 {
		recommendations = append(recommendations, "Review and fix critical issues in the following files:")
		for _, issue := range report.Summary.CriticalIssues {
			recommendations = append(recommendations, fmt.Sprintf("  - %s:%d: %s", issue.File, issue.Line, issue.Message))
		}
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "No critical issues found. Code quality looks good!")
	}

	return recommendations
}
