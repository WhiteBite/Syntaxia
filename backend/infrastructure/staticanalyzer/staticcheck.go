package staticanalyzer

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"strings"
	"time"
)

// StaticcheckAnalyzer реализует StaticAnalyzer для Go с использованием staticcheck
type StaticcheckAnalyzer struct {
	log domain.Logger
}

// NewStaticcheckAnalyzer создает новый анализатор для Go
func NewStaticcheckAnalyzer(log domain.Logger) *StaticcheckAnalyzer {
	return &StaticcheckAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ Go кода
func (a *StaticcheckAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running staticcheck analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что staticcheck установлен
	if err := a.checkStaticcheckInstalled(); err != nil {
		return nil, fmt.Errorf("staticcheck not installed: %w", err)
	}

	// Строим команду для staticcheck
	args := []string{"-f", "json"}

	// Добавляем конфигурационный файл если указан
	if config.ConfigFile != "" {
		args = append(args, "-conf", config.ConfigFile)
	}

	// Добавляем правила если указаны
	if len(config.Rules) > 0 {
		args = append(args, "-checks", strings.Join(config.Rules, ","))
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		args = append(args, "-exclude", strings.Join(config.ExcludeRules, ","))
	}

	// Добавляем путь к проекту
	args = append(args, "./...")

	// Создаем команду
	cmd := exec.CommandContext(ctx, "staticcheck", args...)
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
		// staticcheck может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("Staticcheck found issues: %v", err))
	}

	// Парсим JSON вывод staticcheck
	issues, err := a.parseStaticcheckOutput(output)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse staticcheck output: %v", err))
		// Возвращаем результат с ошибкой парсинга
		result.Error = fmt.Sprintf("Failed to parse output: %v", err)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("Staticcheck analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *StaticcheckAnalyzer) GetSupportedLanguages() []string {
	return []string{"go"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *StaticcheckAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeStaticcheck
}

// ValidateConfig проверяет корректность конфигурации
func (a *StaticcheckAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "go" {
		return fmt.Errorf("staticcheck analyzer only supports Go language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем, что проект содержит Go файлы
	if !a.hasGoFilePaths(config.ProjectPath) {
		return fmt.Errorf("no Go files found in project path")
	}

	return nil
}

// checkStaticcheckInstalled проверяет, установлен ли staticcheck
func (a *StaticcheckAnalyzer) checkStaticcheckInstalled() error {
	cmd := exec.Command("staticcheck", "-version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("staticcheck not found in PATH: %w", err)
	}
	return nil
}

// hasGoFilePaths проверяет, есть ли Go файлы в проекте
func (a *StaticcheckAnalyzer) hasGoFilePaths(projectPath string) bool {
	pattern := filepath.Join(projectPath, "**/*.go")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return false
	}
	return len(matches) > 0
}

// parseStaticcheckOutput парсит JSON вывод staticcheck
func (a *StaticcheckAnalyzer) parseStaticcheckOutput(output []byte) ([]*domain.StaticIssue, error) {
	// staticcheck может выводить несколько JSON объектов, разделенных новой строкой
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var staticcheckIssue struct {
			Rule     string `json:"code"`
			Severity string `json:"severity"`
			Location struct {
				FilePath    string `json:"file"`
				LineNumber  int    `json:"line"`
				ColumnStart int    `json:"column"`
			} `json:"location"`
			End struct {
				LineNumber  int `json:"line"`
				ColumnStart int `json:"column"`
			} `json:"end"`
			Message string `json:"message"`
		}

		if err := json.Unmarshal([]byte(line), &staticcheckIssue); err != nil {
			a.log.Warning(fmt.Sprintf("Failed to parse staticcheck issue: %v", err))
			continue
		}

		// Конвертируем severity
		severity := a.convertSeverity(staticcheckIssue.Severity)

		issue := &domain.StaticIssue{
			File:     staticcheckIssue.Location.FilePath,
			Line:     staticcheckIssue.Location.LineNumber,
			Column:   staticcheckIssue.Location.ColumnStart,
			Severity: severity,
			Message:  staticcheckIssue.Message,
			Code:     staticcheckIssue.Rule,
			Category: a.getCategory(staticcheckIssue.Rule),
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// convertSeverity конвертирует severity staticcheck в общий формат
func (a *StaticcheckAnalyzer) convertSeverity(staticcheckSeverity string) string {
	switch strings.ToLower(staticcheckSeverity) {
	case severityError:
		return severityError
	case severityWarning:
		return severityWarning
	case severityInfo:
		return severityInfo
	default:
		return severityWarning // По умолчанию считаем warning
	}
}

// getCategory определяет категорию проблемы по коду
func (a *StaticcheckAnalyzer) getCategory(code string) string {
	if strings.HasPrefix(code, "SA") {
		return "static-analysis"
	} else if strings.HasPrefix(code, "ST") {
		return "style"
	} else if strings.HasPrefix(code, "U") {
		return "unused"
	} else if strings.HasPrefix(code, "QF") {
		return "quickfix"
	}
	return severityOther
}
