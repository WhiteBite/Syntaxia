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

// Ktlint analyzer constants
const (
	ktlintDefaultTimeout = 300 // 5 minutes
)

// KtlintAnalyzer реализует StaticAnalyzer для Kotlin с использованием ktlint
type KtlintAnalyzer struct {
	log domain.Logger
}

// NewKtlintAnalyzer создает новый анализатор для Kotlin
func NewKtlintAnalyzer(log domain.Logger) *KtlintAnalyzer {
	return &KtlintAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ Kotlin кода
func (a *KtlintAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running ktlint analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Пробуем сначала gradle ktlintCheck, затем standalone ktlint
	var result *domain.StaticAnalysisResult
	var err error

	if a.hasGradleKtlint(config.ProjectPath) {
		result, err = a.runGradleKtlint(ctx, config)
	} else if a.checkKtlintInstalled() {
		result, err = a.runStandaloneKtlint(ctx, config)
	} else {
		return nil, fmt.Errorf("ktlint not available: neither gradle ktlintCheck nor standalone ktlint found")
	}

	if err != nil {
		return nil, err
	}

	result.Duration = time.Since(startTime).Seconds()
	a.log.Info(fmt.Sprintf("Ktlint analysis completed in %.2fs, found %d issues", result.Duration, len(result.Issues)))

	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *KtlintAnalyzer) GetSupportedLanguages() []string {
	return []string{"kotlin", "kt"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *KtlintAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeKtlint
}

// ValidateConfig проверяет корректность конфигурации
func (a *KtlintAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "kotlin" && config.Language != "kt" {
		return fmt.Errorf("ktlint analyzer only supports Kotlin language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	return nil
}

// hasGradleKtlint проверяет наличие ktlint plugin в Gradle проекте
func (a *KtlintAnalyzer) hasGradleKtlint(projectPath string) bool {
	// Проверяем build.gradle.kts
	ktsPath := filepath.Join(projectPath, "build.gradle.kts")
	if content, err := os.ReadFile(ktsPath); err == nil {
		if strings.Contains(string(content), "ktlint") {
			return true
		}
	}

	// Проверяем build.gradle
	gradlePath := filepath.Join(projectPath, "build.gradle")
	if content, err := os.ReadFile(gradlePath); err == nil {
		if strings.Contains(string(content), "ktlint") {
			return true
		}
	}

	return false
}

// checkKtlintInstalled проверяет, установлен ли standalone ktlint
func (a *KtlintAnalyzer) checkKtlintInstalled() bool {
	cmd := exec.Command("ktlint", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// runGradleKtlint запускает ktlint через Gradle
func (a *KtlintAnalyzer) runGradleKtlint(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	result := &domain.StaticAnalysisResult{
		Success:     true,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    config.Analyzer,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	// Определяем команду gradle
	gradleCmd := a.getGradleCommand(config.ProjectPath)
	if gradleCmd == "" {
		return nil, fmt.Errorf("gradle not found")
	}

	// Запускаем gradle ktlintCheck
	cmd := exec.CommandContext(ctx, gradleCmd, "ktlintCheck", "--continue")
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Устанавливаем переменные окружения
	cmd.Env = os.Environ()
	if config.EnvVars != nil {
		for key, value := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		// ktlint может вернуть ошибку, если найдет проблемы
		a.log.Warning(fmt.Sprintf("Ktlint found issues: %v", err))
	}

	// Парсим вывод
	issues := a.parseGradleKtlintOutput(string(output))
	result.Issues = issues
	result.Summary = generateSummary(issues)
	result.Success = result.Summary.ErrorCount == 0

	return result, nil
}

// runStandaloneKtlint запускает standalone ktlint
func (a *KtlintAnalyzer) runStandaloneKtlint(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	result := &domain.StaticAnalysisResult{
		Success:     true,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    config.Analyzer,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	// Строим команду для ktlint
	args := []string{"--reporter=json"}

	// Добавляем путь к файлам
	args = append(args, "**/*.kt")

	// Создаем команду
	cmd := exec.CommandContext(ctx, "ktlint", args...)
	executil.HideWindow(cmd)
	cmd.Dir = config.ProjectPath

	// Устанавливаем переменные окружения
	cmd.Env = os.Environ()
	if config.EnvVars != nil {
		for key, value := range config.EnvVars {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		// ktlint может вернуть ошибку, если найдет проблемы
		a.log.Warning(fmt.Sprintf("Ktlint found issues: %v", err))
	}

	// Парсим JSON вывод
	issues, parseErr := a.parseKtlintJSONOutput(output)
	if parseErr != nil {
		// Fallback to text parsing
		a.log.Warning(fmt.Sprintf("Failed to parse JSON output, trying text: %v", parseErr))
		issues = a.parseKtlintTextOutput(string(output))
	}

	result.Issues = issues
	result.Summary = generateSummary(issues)
	result.Success = result.Summary.ErrorCount == 0

	return result, nil
}

// getGradleCommand возвращает команду для запуска Gradle
func (a *KtlintAnalyzer) getGradleCommand(projectPath string) string {
	// Проверяем наличие gradlew wrapper
	if a.hasGradleWrapper(projectPath) {
		if os.PathSeparator == '\\' {
			return filepath.Join(projectPath, "gradlew.bat")
		}
		return filepath.Join(projectPath, "gradlew")
	}

	// Используем системный gradle
	if _, err := exec.LookPath("gradle"); err == nil {
		return "gradle"
	}

	return ""
}

// hasGradleWrapper проверяет наличие Gradle wrapper
func (a *KtlintAnalyzer) hasGradleWrapper(projectPath string) bool {
	wrapperFiles := []string{"gradlew", "gradlew.bat"}
	for _, wrapper := range wrapperFiles {
		if _, err := os.Stat(filepath.Join(projectPath, wrapper)); err == nil {
			return true
		}
	}
	return false
}

// ktlintJSONIssue представляет issue в JSON формате ktlint
type ktlintJSONIssue struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Message string `json:"message"`
	Rule    string `json:"rule"`
}

// parseKtlintJSONOutput парсит JSON вывод ktlint
func (a *KtlintAnalyzer) parseKtlintJSONOutput(output []byte) ([]*domain.StaticIssue, error) {
	var jsonIssues []ktlintJSONIssue
	if err := json.Unmarshal(output, &jsonIssues); err != nil {
		return nil, err
	}

	issues := make([]*domain.StaticIssue, 0, len(jsonIssues))
	for _, ji := range jsonIssues {
		issue := &domain.StaticIssue{
			File:     ji.File,
			Line:     ji.Line,
			Column:   ji.Column,
			Severity: severityWarning,
			Message:  ji.Message,
			Code:     ji.Rule,
			Category: a.categorizeRule(ji.Rule),
		}
		issues = append(issues, issue)
	}

	return issues, nil
}

// parseKtlintTextOutput парсит текстовый вывод ktlint
func (a *KtlintAnalyzer) parseKtlintTextOutput(output string) []*domain.StaticIssue {
	var issues []*domain.StaticIssue

	// ktlint text format: "file.kt:10:5: Message (rule-name)"
	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, ".kt:") {
			continue
		}

		issue := a.parseKtlintLine(line)
		if issue != nil {
			issues = append(issues, issue)
		}
	}

	return issues
}

// parseKtlintLine парсит одну строку вывода ktlint
func (a *KtlintAnalyzer) parseKtlintLine(line string) *domain.StaticIssue {
	// Format: "file.kt:10:5: Message (rule-name)"
	parts := strings.SplitN(line, ":", 4)
	if len(parts) < 4 {
		return nil
	}

	file := parts[0]
	lineNum := a.parseInt(parts[1])
	column := a.parseInt(parts[2])
	messageAndRule := strings.TrimSpace(parts[3])

	// Extract rule from message
	message := messageAndRule
	rule := ""
	if idx := strings.LastIndex(messageAndRule, "("); idx != -1 {
		if endIdx := strings.LastIndex(messageAndRule, ")"); endIdx > idx {
			rule = messageAndRule[idx+1 : endIdx]
			message = strings.TrimSpace(messageAndRule[:idx])
		}
	}

	return &domain.StaticIssue{
		File:     file,
		Line:     lineNum,
		Column:   column,
		Severity: severityWarning,
		Message:  message,
		Code:     rule,
		Category: a.categorizeRule(rule),
	}
}

// parseGradleKtlintOutput парсит вывод gradle ktlintCheck
func (a *KtlintAnalyzer) parseGradleKtlintOutput(output string) []*domain.StaticIssue {
	var issues []*domain.StaticIssue

	lines := strings.Split(output, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Gradle ktlint output format: "file.kt:10:5: Message (rule-name)"
		if strings.Contains(line, ".kt:") && !strings.HasPrefix(line, ">") {
			issue := a.parseKtlintLine(line)
			if issue != nil {
				issues = append(issues, issue)
			}
		}
	}

	return issues
}

// parseInt парсит строку в int
func (a *KtlintAnalyzer) parseInt(s string) int {
	var i int
	_, _ = fmt.Sscanf(s, "%d", &i)
	return i
}

// categorizeRule категоризирует правило ktlint
func (a *KtlintAnalyzer) categorizeRule(rule string) string {
	if rule == "" {
		return severityOther
	}

	// ktlint rule categories
	switch {
	case strings.Contains(rule, "indent"):
		return "formatting"
	case strings.Contains(rule, "spacing"):
		return "formatting"
	case strings.Contains(rule, "import"):
		return "imports"
	case strings.Contains(rule, "naming"):
		return "naming"
	case strings.Contains(rule, "comment"):
		return "documentation"
	case strings.Contains(rule, "max-line-length"):
		return "formatting"
	case strings.Contains(rule, "no-wildcard-imports"):
		return "imports"
	case strings.Contains(rule, "no-unused"):
		return "unused"
	default:
		return "style"
	}
}
