package staticanalyzer

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

// Clippy analyzer constants
const (
	clippyDefaultTimeout = 300 // 5 minutes
)

// ClippyAnalyzer реализует StaticAnalyzer для Rust с использованием clippy
type ClippyAnalyzer struct {
	log domain.Logger
}

// NewClippyAnalyzer создает новый анализатор для Rust
func NewClippyAnalyzer(log domain.Logger) *ClippyAnalyzer {
	return &ClippyAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ Rust кода
func (a *ClippyAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running clippy analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что clippy установлен
	if err := a.checkClippyInstalled(); err != nil {
		return nil, fmt.Errorf("clippy not installed: %w", err)
	}

	// Строим команду для clippy
	args := []string{"clippy", "--message-format=json"}

	// Добавляем флаги для более строгого анализа
	args = append(args, "--", "-D", "warnings")

	// Добавляем дополнительные правила если указаны
	if len(config.Rules) > 0 {
		for _, rule := range config.Rules {
			args = append(args, "-W", rule)
		}
	}

	// Добавляем исключения если указаны
	if len(config.ExcludeRules) > 0 {
		for _, rule := range config.ExcludeRules {
			args = append(args, "-A", rule)
		}
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "cargo", args...)
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
		// clippy может вернуть ошибку, если найдет проблемы
		// Это не всегда означает, что анализ не удался
		a.log.Warning(fmt.Sprintf("Clippy found issues: %v", err))
	}

	// Парсим JSON вывод clippy
	issues, parseErr := a.parseClippyOutput(output)
	if parseErr != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse clippy output: %v", parseErr))
		result.Error = fmt.Sprintf("Failed to parse output: %v", parseErr)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	// Определяем успешность на основе наличия ошибок
	result.Success = result.Summary.ErrorCount == 0

	a.log.Info(fmt.Sprintf("Clippy analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *ClippyAnalyzer) GetSupportedLanguages() []string {
	return []string{"rust", "rs"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *ClippyAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeClippy
}

// ValidateConfig проверяет корректность конфигурации
func (a *ClippyAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "rust" && config.Language != "rs" {
		return fmt.Errorf("clippy analyzer only supports Rust language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем наличие Cargo.toml
	cargoPath := filepath.Join(config.ProjectPath, "Cargo.toml")
	if _, err := os.Stat(cargoPath); os.IsNotExist(err) {
		return fmt.Errorf("Cargo.toml not found in project path")
	}

	return nil
}

// checkClippyInstalled проверяет, установлен ли clippy
func (a *ClippyAnalyzer) checkClippyInstalled() error {
	cmd := exec.Command("cargo", "clippy", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("clippy not found, install with: rustup component add clippy")
	}
	return nil
}

// clippyMessage представляет сообщение от clippy в JSON формате
type clippyMessage struct {
	Reason  string `json:"reason"`
	Message *struct {
		Code *struct {
			Code        string `json:"code"`
			Explanation string `json:"explanation,omitempty"`
		} `json:"code"`
		Level   string `json:"level"`
		Message string `json:"message"`
		Spans   []struct {
			FileName    string `json:"file_name"`
			LineStart   int    `json:"line_start"`
			LineEnd     int    `json:"line_end"`
			ColumnStart int    `json:"column_start"`
			ColumnEnd   int    `json:"column_end"`
			IsPrimary   bool   `json:"is_primary"`
			Label       string `json:"label,omitempty"`
		} `json:"spans"`
		Children []struct {
			Level   string `json:"level"`
			Message string `json:"message"`
		} `json:"children,omitempty"`
	} `json:"message"`
}

// parseClippyOutput парсит JSON вывод clippy
func (a *ClippyAnalyzer) parseClippyOutput(output []byte) ([]*domain.StaticIssue, error) {
	lines := strings.Split(string(output), "\n")
	issues := make([]*domain.StaticIssue, 0, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var msg clippyMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			// Пропускаем строки, которые не являются JSON
			continue
		}

		// Обрабатываем только compiler-message
		if msg.Reason != "compiler-message" || msg.Message == nil {
			continue
		}

		// Пропускаем сообщения без spans
		if len(msg.Message.Spans) == 0 {
			continue
		}

		// Находим primary span
		var primarySpan *struct {
			FileName    string `json:"file_name"`
			LineStart   int    `json:"line_start"`
			LineEnd     int    `json:"line_end"`
			ColumnStart int    `json:"column_start"`
			ColumnEnd   int    `json:"column_end"`
			IsPrimary   bool   `json:"is_primary"`
			Label       string `json:"label,omitempty"`
		}

		for i := range msg.Message.Spans {
			if msg.Message.Spans[i].IsPrimary {
				primarySpan = &msg.Message.Spans[i]
				break
			}
		}

		if primarySpan == nil {
			primarySpan = &msg.Message.Spans[0]
		}

		// Конвертируем severity
		severity := a.convertSeverity(msg.Message.Level)

		// Получаем код ошибки
		code := ""
		if msg.Message.Code != nil {
			code = msg.Message.Code.Code
		}

		issue := &domain.StaticIssue{
			File:     primarySpan.FileName,
			Line:     primarySpan.LineStart,
			Column:   primarySpan.ColumnStart,
			Severity: severity,
			Message:  msg.Message.Message,
			Code:     code,
			Category: a.getCategory(code),
		}

		// Добавляем suggestions из children
		for _, child := range msg.Message.Children {
			if child.Level == "help" || child.Level == "note" {
				issue.Suggestions = append(issue.Suggestions, child.Message)
			}
		}

		issues = append(issues, issue)
	}

	return issues, nil
}

// convertSeverity конвертирует severity clippy в общий формат
func (a *ClippyAnalyzer) convertSeverity(clippySeverity string) string {
	switch strings.ToLower(clippySeverity) {
	case "error":
		return severityError
	case "warning":
		return severityWarning
	case "note":
		return severityInfo
	case "help":
		return severityHint
	default:
		return severityWarning
	}
}

// getCategory определяет категорию проблемы по коду
func (a *ClippyAnalyzer) getCategory(code string) string {
	if code == "" {
		return severityOther
	}

	// Clippy lint categories
	switch {
	case strings.HasPrefix(code, "clippy::"):
		lintName := strings.TrimPrefix(code, "clippy::")
		return a.categorizeLint(lintName)
	case strings.HasPrefix(code, "E"):
		return "compiler-error"
	case strings.HasPrefix(code, "W"):
		return "compiler-warning"
	default:
		return severityOther
	}
}

// categorizeLint категоризирует clippy lint
func (a *ClippyAnalyzer) categorizeLint(lintName string) string {
	// Основные категории clippy lints
	complexityLints := []string{
		"cognitive_complexity", "too_many_arguments", "too_many_lines",
		"type_complexity", "excessive_precision",
	}

	styleLints := []string{
		"needless_return", "redundant_closure", "single_match",
		"match_bool", "if_same_then_else", "collapsible_if",
	}

	perfLints := []string{
		"needless_collect", "unnecessary_to_owned", "clone_on_copy",
		"useless_vec", "box_collection",
	}

	correctnessLints := []string{
		"eq_op", "erasing_op", "almost_swapped", "suspicious_arithmetic_impl",
		"misrefactored_assign_op",
	}

	for _, lint := range complexityLints {
		if strings.Contains(lintName, lint) {
			return "complexity"
		}
	}

	for _, lint := range styleLints {
		if strings.Contains(lintName, lint) {
			return "style"
		}
	}

	for _, lint := range perfLints {
		if strings.Contains(lintName, lint) {
			return "performance"
		}
	}

	for _, lint := range correctnessLints {
		if strings.Contains(lintName, lint) {
			return "correctness"
		}
	}

	return "clippy"
}
