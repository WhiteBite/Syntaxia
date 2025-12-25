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

// ErrorProneAnalyzer реализует StaticAnalyzer для Java с использованием ErrorProne
type ErrorProneAnalyzer struct {
	log domain.Logger
}

// NewErrorProneAnalyzer создает новый анализатор для Java
func NewErrorProneAnalyzer(log domain.Logger) *ErrorProneAnalyzer {
	return &ErrorProneAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ Java кода
func (a *ErrorProneAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running ErrorProne analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что Java и ErrorProne доступны
	if err := a.checkJavaEnvironment(); err != nil {
		a.log.Warning(fmt.Sprintf("Java environment not available: %v", err))
		// Возвращаем успешный результат с предупреждением вместо ошибки
		return &domain.StaticAnalysisResult{
			Success:     true,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Duration:    time.Since(startTime).Seconds(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
			Error:       fmt.Sprintf("Java environment not available: %v", err),
		}, nil
	}

	// Строим команду для ErrorProne
	args := []string{"-Xplugin:ErrorProne"}

	// Добавляем конфигурационный файл если указан
	if config.ConfigFile != "" {
		args = append(args, "-XepPatchChecks:"+config.ConfigFile)
	}

	// Добавляем правила если указаны
	if len(config.Rules) > 0 {
		args = append(args, "-XepDisableAllChecks", "-Xep:"+strings.Join(config.Rules, ","))
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		args = append(args, "-XepDisable:"+strings.Join(config.ExcludeRules, ","))
	}

	// Добавляем путь к проекту
	args = append(args, config.ProjectPath)

	// Создаем команду
	cmd := exec.CommandContext(ctx, "javac", args...)
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
		// ErrorProne может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("ErrorProne found issues: %v", err))
	}

	// Парсим вывод ErrorProne
	issues, err := a.parseErrorProneOutput(output)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse ErrorProne output: %v", err))
		// Возвращаем результат с ошибкой парсинга
		result.Error = fmt.Sprintf("Failed to parse output: %v", err)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("ErrorProne analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *ErrorProneAnalyzer) GetSupportedLanguages() []string {
	return []string{"java"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *ErrorProneAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeErrorProne
}

// ValidateConfig проверяет корректность конфигурации
func (a *ErrorProneAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "java" {
		return fmt.Errorf("ErrorProne analyzer only supports Java language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем, что проект содержит Java файлы
	if !a.hasJavaFilePaths(config.ProjectPath) {
		return fmt.Errorf("no Java files found in project path")
	}

	return nil
}

// checkJavaEnvironment проверяет доступность Java окружения
func (a *ErrorProneAnalyzer) checkJavaEnvironment() error {
	// Проверяем Java
	cmd := exec.Command("java", "-version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Java not found: %w", err)
	}

	// Проверяем javac
	cmd = exec.Command("javac", "-version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("javac not found: %w", err)
	}

	return nil
}

// hasJavaFilePaths проверяет, есть ли Java файлы в проекте
func (a *ErrorProneAnalyzer) hasJavaFilePaths(projectPath string) bool {
	sep := string(filepath.Separator)
	patterns := []string{
		projectPath + sep + "**" + sep + "*.java",
		projectPath + sep + "src" + sep + "**" + sep + "*.java",
		projectPath + sep + "main" + sep + "**" + sep + "*.java",
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return true
		}
	}
	return false
}

// parseErrorProneOutput парсит вывод ErrorProne
func (a *ErrorProneAnalyzer) parseErrorProneOutput(output []byte) ([]*domain.StaticIssue, error) {
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines)/2)

	// ErrorProne выводит ошибки в формате:
	// filename:line:column: [error] message
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ":") {
			continue
		}

		// Парсим строку ошибки
		parts := strings.SplitN(line, ":", 4)
		if len(parts) < 4 {
			continue
		}

		file := parts[0]
		lineNum := 0
		columnNum := 0

		// Парсим номер строки
		if _, err := fmt.Sscanf(parts[1], "%d", &lineNum); err != nil {
			continue
		}

		// Парсим номер колонки
		if _, err := fmt.Sscanf(parts[2], "%d", &columnNum); err != nil {
			continue
		}

		// Извлекаем сообщение
		message := parts[3]
		if strings.HasPrefix(message, " [error] ") {
			message = strings.TrimPrefix(message, " [error] ")
		} else if strings.HasPrefix(message, " [warning] ") {
			message = strings.TrimPrefix(message, " [warning] ")
		}

		// Определяем severity
		severity := severityWarning
		if strings.Contains(line, "[error]") {
			severity = severityError
		}

		issue := &domain.StaticIssue{
			File:     file,
			Line:     lineNum,
			Column:   columnNum,
			Severity: severity,
			Message:  message,
			Code:     "errorprone",
			Category: a.getCategory(message),
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// getCategory определяет категорию проблемы по сообщению
func (a *ErrorProneAnalyzer) getCategory(message string) string {
	message = strings.ToLower(message)

	if strings.Contains(message, "null") {
		return "null-safety"
	} else if strings.Contains(message, "unused") {
		return "unused-code"
	} else if strings.Contains(message, "deprecated") {
		return "deprecation"
	} else if strings.Contains(message, "concurrent") {
		return "concurrency"
	} else if strings.Contains(message, "resource") {
		return "resource-management"
	}
	return severityOther
}
