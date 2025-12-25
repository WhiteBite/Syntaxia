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

// ESLintAnalyzer реализует StaticAnalyzer для TypeScript/JavaScript с использованием ESLint
type ESLintAnalyzer struct {
	log domain.Logger
}

// NewESLintAnalyzer создает новый анализатор для TypeScript/JavaScript
func NewESLintAnalyzer(log domain.Logger) *ESLintAnalyzer {
	return &ESLintAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ TypeScript/JavaScript кода
func (a *ESLintAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running ESLint analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что ESLint установлен
	if err := a.checkESLintInstalled(); err != nil {
		return nil, fmt.Errorf("ESLint not installed: %w", err)
	}

	// Строим команду для ESLint
	args := []string{"--format", "json"}

	// Добавляем конфигурационный файл если указан
	if config.ConfigFile != "" {
		args = append(args, "--config", config.ConfigFile)
	}

	// Добавляем правила если указаны
	if len(config.Rules) > 0 {
		args = append(args, "--rule", strings.Join(config.Rules, ","))
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		args = append(args, "--no-eslintrc")
		args = append(args, "--rule", strings.Join(config.ExcludeRules, ","))
	}

	// Добавляем путь к проекту
	args = append(args, ".")

	// Создаем команду
	cmd := exec.CommandContext(ctx, "npx", append([]string{"eslint"}, args...)...) //nolint:gosec // External tool command
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
		// ESLint может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("ESLint found issues: %v", err))
	}

	// Парсим JSON вывод ESLint
	issues, err := a.parseESLintOutput(output)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse ESLint output: %v", err))
		// Возвращаем результат с ошибкой парсинга
		result.Error = fmt.Sprintf("Failed to parse output: %v", err)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	a.log.Info(fmt.Sprintf("ESLint analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *ESLintAnalyzer) GetSupportedLanguages() []string {
	return []string{"typescript", "ts", "javascript", "js"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *ESLintAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeESLint
}

// ValidateConfig проверяет корректность конфигурации
func (a *ESLintAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	supportedLanguages := []string{"typescript", "ts", "javascript", "js"}
	isSupported := false
	for _, lang := range supportedLanguages {
		if config.Language == lang {
			isSupported = true
			break
		}
	}

	if !isSupported {
		return fmt.Errorf("ESLint analyzer only supports TypeScript/JavaScript languages")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем, что проект содержит TypeScript/JavaScript файлы
	if !a.hasTypeScriptFilePaths(config.ProjectPath) {
		return fmt.Errorf("no TypeScript/JavaScript files found in project path")
	}

	return nil
}

// checkESLintInstalled проверяет, установлен ли ESLint
func (a *ESLintAnalyzer) checkESLintInstalled() error {
	cmd := exec.Command("npx", "eslint", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ESLint not found: %w", err)
	}
	return nil
}

// hasTypeScriptFilePaths проверяет, есть ли TypeScript/JavaScript файлы в проекте
func (a *ESLintAnalyzer) hasTypeScriptFilePaths(projectPath string) bool {
	patterns := []string{
		filepath.Join(projectPath, "**/*.ts"),
		filepath.Join(projectPath, "**/*.tsx"),
		filepath.Join(projectPath, "**/*.js"),
		filepath.Join(projectPath, "**/*.jsx"),
	}

	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err == nil && len(matches) > 0 {
			return true
		}
	}
	return false
}

// parseESLintOutput парсит JSON вывод ESLint
func (a *ESLintAnalyzer) parseESLintOutput(output []byte) ([]*domain.StaticIssue, error) {
	var eslintResults []struct {
		FilePathPath string `json:"filePath"`
		Messages     []struct {
			RuleID         string `json:"ruleId"`
			Severity       int    `json:"severity"`
			Message        string `json:"message"`
			LineNumber     int    `json:"line"`
			ColumnStart    int    `json:"column"`
			NodeType       string `json:"nodeType,omitempty"`
			MessageID      string `json:"messageId,omitempty"`
			EndLineNumber  int    `json:"endLineNumber,omitempty"`
			EndColumnStart int    `json:"endColumnStart,omitempty"`
			Fix            struct {
				Range []int  `json:"range"`
				Text  string `json:"text"`
			} `json:"fix,omitempty"`
		} `json:"messages"`
		ErrorCount          int `json:"errorCount"`
		WarningCount        int `json:"warningCount"`
		FixableErrorCount   int `json:"fixableErrorCount"`
		FixableWarningCount int `json:"fixableWarningCount"`
	}

	if err := json.Unmarshal(output, &eslintResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ESLint output: %w", err)
	}

	var issues []*domain.StaticIssue

	for _, result := range eslintResults {
		for _, message := range result.Messages {
			// Конвертируем severity
			severity := a.convertSeverity(message.Severity)

			issue := &domain.StaticIssue{
				File:     result.FilePathPath,
				Line:     message.LineNumber,
				Column:   message.ColumnStart,
				Severity: severity,
				Message:  message.Message,
				Code:     message.RuleID,
				Category: a.getCategory(message.RuleID),
			}

			// Добавляем предложения по исправлению если есть
			if message.Fix.Text != "" {
				issue.Suggestions = []string{message.Fix.Text}
			}

			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// convertSeverity конвертирует severity ESLint в общий формат
func (a *ESLintAnalyzer) convertSeverity(eslintSeverity int) string {
	switch eslintSeverity {
	case 0:
		return severityInfo
	case 1:
		return severityWarning
	case 2:
		return severityError
	default:
		return severityWarning
	}
}

// getCategory определяет категорию проблемы по коду правила
func (a *ESLintAnalyzer) getCategory(ruleID string) string {
	if strings.HasPrefix(ruleID, "@typescript-eslint/") {
		return "typescript"
	} else if strings.HasPrefix(ruleID, "react/") {
		return "react"
	} else if strings.HasPrefix(ruleID, "import/") {
		return "imports"
	} else if strings.HasPrefix(ruleID, "prefer-") {
		return "style"
	} else if strings.HasPrefix(ruleID, "no-") {
		return "best-practices"
	}
	return severityOther
}
