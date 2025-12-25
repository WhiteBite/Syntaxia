package staticanalyzer

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"strings"
	"time"
)

// ClangTidyAnalyzer реализует StaticAnalyzer для C/C++ с использованием ClangTidy
type ClangTidyAnalyzer struct {
	log domain.Logger
}

// NewClangTidyAnalyzer создает новый анализатор для C/C++
func NewClangTidyAnalyzer(log domain.Logger) *ClangTidyAnalyzer {
	return &ClangTidyAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ C/C++ кода
func (a *ClangTidyAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running ClangTidy analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что ClangTidy установлен
	if err := a.checkClangTidyInstalled(); err != nil {
		a.log.Warning(fmt.Sprintf("ClangTidy not available: %v", err))
		// Возвращаем успешный результат с предупреждением вместо ошибки
		return &domain.StaticAnalysisResult{
			Success:     true,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Duration:    time.Since(startTime).Seconds(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
			Error:       fmt.Sprintf("ClangTidy not available: %v", err),
		}, nil
	}

	// Строим команду для ClangTidy
	args := []string{"--format-style=json"}

	// Добавляем конфигурационный файл если указан
	if config.ConfigFile != "" {
		args = append(args, "--config-file="+config.ConfigFile)
	}

	// Добавляем правила если указаны
	if len(config.Rules) > 0 {
		args = append(args, "-checks="+strings.Join(config.Rules, ","))
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		args = append(args, "-checks=-"+strings.Join(config.ExcludeRules, ","))
	}

	// Добавляем путь к проекту
	args = append(args, config.ProjectPath)

	// Создаем команду
	cmd := exec.CommandContext(ctx, "clang-tidy", args...)
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

	result := &domain.StaticAnalysisResult{
		Success:     true,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    config.Analyzer,
		Duration:    duration,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	if err != nil {
		// ClangTidy может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("ClangTidy found issues: %v", err))
	}

	// Парсим JSON вывод ClangTidy
	issues, err := a.parseClangTidyOutput(output)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse ClangTidy output: %v", err))
		// Возвращаем результат с ошибкой парсинга
		result.Error = fmt.Sprintf("Failed to parse output: %v", err)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("ClangTidy analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *ClangTidyAnalyzer) GetSupportedLanguages() []string {
	return []string{"c", "cpp", "cc"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *ClangTidyAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeClangTidy
}

// ValidateConfig проверяет корректность конфигурации
func (a *ClangTidyAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	supportedLanguages := []string{"c", "cpp", "cc"}
	isSupported := false
	for _, lang := range supportedLanguages {
		if config.Language == lang {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("ClangTidy analyzer only supports C/C++ languages")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем, что проект содержит C/C++ файлы
	if !a.hasCppFilePaths(config.ProjectPath) {
		return fmt.Errorf("no C/C++ files found in project path")
	}

	return nil
}

// checkClangTidyInstalled проверяет, установлен ли ClangTidy
func (a *ClangTidyAnalyzer) checkClangTidyInstalled() error {
	cmd := exec.Command("clang-tidy", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ClangTidy not found: %w", err)
	}
	return nil
}

// hasCppFilePaths проверяет, есть ли C/C++ файлы в проекте
func (a *ClangTidyAnalyzer) hasCppFilePaths(projectPath string) bool {
	sep := string(filepath.Separator)
	patterns := []string{
		projectPath + sep + "**" + sep + "*.c",
		projectPath + sep + "**" + sep + "*.cpp",
		projectPath + sep + "**" + sep + "*.cc",
		projectPath + sep + "**" + sep + "*.h",
		projectPath + sep + "**" + sep + "*.hpp",
		projectPath + sep + "src" + sep + "**" + sep + "*.c",
		projectPath + sep + "src" + sep + "**" + sep + "*.cpp",
		projectPath + sep + "src" + sep + "**" + sep + "*.cc",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return true
		}
	}
	return false
}

// parseClangTidyOutput парсит JSON вывод ClangTidy
func (a *ClangTidyAnalyzer) parseClangTidyOutput(output []byte) ([]*domain.StaticIssue, error) {
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines))

	for _, line := range lines {
		if issue := a.parseClangTidyLine(strings.TrimSpace(line)); issue != nil {
			issues = append(issues, issue)
		}
	}
	return issues, nil
}

// parseClangTidyLine parses a single ClangTidy output line
func (a *ClangTidyAnalyzer) parseClangTidyLine(line string) *domain.StaticIssue {
	if line == "" || !strings.Contains(line, "\"DiagnosticName\"") {
		return nil
	}

	diagnosticName := a.extractClangString(line, "\"DiagnosticName\": \"")
	message := a.extractClangString(line, "\"Message\": \"")
	filePath := a.extractClangString(line, "\"FilePathPath\": \"")
	lineNum := a.extractClangInt(line, "\"FilePathLineNumber\": ")
	columnNum := a.extractClangInt(line, "\"FilePathColumnStart\": ")

	if diagnosticName == "" || filePath == "" {
		return nil
	}

	severity := severityWarning
	if strings.Contains(diagnosticName, severityError) {
		severity = severityError
	}

	return &domain.StaticIssue{
		File: filePath, Line: lineNum, Column: columnNum,
		Severity: severity, Message: message, Code: diagnosticName, Category: a.getCategory(diagnosticName),
	}
}

// extractClangString extracts a string value from JSON-like line
func (a *ClangTidyAnalyzer) extractClangString(line, prefix string) string {
	start := strings.Index(line, prefix)
	if start == -1 {
		return ""
	}
	start += len(prefix)
	end := strings.Index(line[start:], "\"")
	if end == -1 {
		return ""
	}
	return line[start : start+end]
}

// extractClangInt extracts an int value from JSON-like line
func (a *ClangTidyAnalyzer) extractClangInt(line, prefix string) int {
	start := strings.Index(line, prefix)
	if start == -1 {
		return 0
	}
	start += len(prefix)
	end := strings.Index(line[start:], ",")
	if end == -1 {
		end = strings.Index(line[start:], "}")
	}
	if end == -1 {
		return 0
	}
	var val int
	_, _ = fmt.Sscanf(line[start:start+end], "%d", &val)
	return val
}

// getCategory определяет категорию проблемы по имени диагностики
func (a *ClangTidyAnalyzer) getCategory(diagnosticName string) string {
	diagnosticName = strings.ToLower(diagnosticName)

	if strings.Contains(diagnosticName, "unused") {
		return "unused-code"
	} else if strings.Contains(diagnosticName, "null") {
		return "null-safety"
	} else if strings.Contains(diagnosticName, "memory") {
		return "memory-management"
	} else if strings.Contains(diagnosticName, "performance") {
		return "performance"
	} else if strings.Contains(diagnosticName, "modernize") {
		return "modernization"
	} else if strings.Contains(diagnosticName, "bugprone") {
		return "bug-prone"
	}
	return severityOther
}
