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

// RuffAnalyzer реализует StaticAnalyzer для Python с использованием Ruff
type RuffAnalyzer struct {
	log domain.Logger
}

// NewRuffAnalyzer создает новый анализатор для Python
func NewRuffAnalyzer(log domain.Logger) *RuffAnalyzer {
	return &RuffAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ Python кода
func (a *RuffAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running Ruff analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что Ruff установлен
	if err := a.checkRuffInstalled(); err != nil {
		a.log.Warning(fmt.Sprintf("Ruff not available: %v", err))
		// Возвращаем успешный результат с предупреждением вместо ошибки
		return &domain.StaticAnalysisResult{
			Success:     true,
			Language:    config.Language,
			ProjectPath: config.ProjectPath,
			Analyzer:    config.Analyzer,
			Duration:    time.Since(startTime).Seconds(),
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
			Error:       fmt.Sprintf("Ruff not available: %v", err),
		}, nil
	}

	// Строим команду для Ruff
	args := []string{"check", "--output-format", "json"}

	// Добавляем конфигурационный файл если указан
	if config.ConfigFile != "" {
		args = append(args, "--config", config.ConfigFile)
	}

	// Добавляем правила если указаны
	if len(config.Rules) > 0 {
		args = append(args, "--select", strings.Join(config.Rules, ","))
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		args = append(args, "--ignore", strings.Join(config.ExcludeRules, ","))
	}

	// Добавляем путь к проекту
	args = append(args, config.ProjectPath)

	// Создаем команду
	cmd := exec.CommandContext(ctx, "ruff", args...)
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
		// Ruff может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("Ruff found issues: %v", err))
	}

	// Парсим JSON вывод Ruff
	issues, err := a.parseRuffOutput(output)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse Ruff output: %v", err))
		// Возвращаем результат с ошибкой парсинга
		result.Error = fmt.Sprintf("Failed to parse output: %v", err)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("Ruff analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *RuffAnalyzer) GetSupportedLanguages() []string {
	return []string{"python", "py"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *RuffAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeRuff
}

// ValidateConfig проверяет корректность конфигурации
func (a *RuffAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	supportedLanguages := []string{"python", "py"}
	isSupported := false
	for _, lang := range supportedLanguages {
		if config.Language == lang {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("Ruff analyzer only supports Python language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем, что проект содержит Python файлы
	if !a.hasPythonFilePaths(config.ProjectPath) {
		return fmt.Errorf("no Python files found in project path")
	}

	return nil
}

// checkRuffInstalled проверяет, установлен ли Ruff
func (a *RuffAnalyzer) checkRuffInstalled() error {
	cmd := exec.Command("ruff", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Ruff not found: %w", err)
	}
	return nil
}

// hasPythonFilePaths проверяет, есть ли Python файлы в проекте
func (a *RuffAnalyzer) hasPythonFilePaths(projectPath string) bool {
	patterns := []string{
		filepath.Join(projectPath, "**/*.py"),
		filepath.Join(projectPath, "src/**/*.py"),
		filepath.Join(projectPath, "main/**/*.py"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return true
		}
	}
	return false
}

// parseRuffOutput парсит JSON вывод Ruff
func (a *RuffAnalyzer) parseRuffOutput(output []byte) ([]*domain.StaticIssue, error) {
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines))

	for _, line := range lines {
		if issue := a.parseRuffLine(strings.TrimSpace(line)); issue != nil {
			issues = append(issues, issue)
		}
	}
	return issues, nil
}

// parseRuffLine parses a single Ruff output line
func (a *RuffAnalyzer) parseRuffLine(line string) *domain.StaticIssue {
	if line == "" || !strings.Contains(line, "\"code\"") {
		return nil
	}

	code := extractJSONString(line, "\"code\": \"")
	message := extractJSONString(line, "\"message\": \"")
	filename := extractJSONString(line, "\"filename\": \"")
	row := extractJSONInt(line, "\"row\": ")
	column := extractJSONInt(line, "\"column\": ")

	if code == "" || filename == "" {
		return nil
	}

	severity := "warning"
	if strings.HasPrefix(code, "E") {
		severity = "error"
	}

	return &domain.StaticIssue{
		File: filename, Line: row, Column: column,
		Severity: severity, Message: message, Code: code, Category: a.getCategory(code),
	}
}

// extractJSONString extracts a string value from JSON-like line
func extractJSONString(line, prefix string) string {
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

// extractJSONInt extracts an int value from JSON-like line
func extractJSONInt(line, prefix string) int {
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

// getCategory определяет категорию проблемы по коду
func (a *RuffAnalyzer) getCategory(code string) string {
	if strings.HasPrefix(code, "E") {
		return "error"
	} else if strings.HasPrefix(code, "W") {
		return "warning"
	} else if strings.HasPrefix(code, "F") {
		return "flake8"
	} else if strings.HasPrefix(code, "I") {
		return "imports"
	} else if strings.HasPrefix(code, "N") {
		return "naming"
	} else if strings.HasPrefix(code, "UP") {
		return "pyupgrade"
	}
	return severityOther
}
