package staticanalyzer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"time"
)

// DotnetFormatAnalyzer реализует StaticAnalyzer для C# с использованием dotnet format
type DotnetFormatAnalyzer struct {
	log domain.Logger
}

// NewDotnetFormatAnalyzer создает новый анализатор для C#
func NewDotnetFormatAnalyzer(log domain.Logger) *DotnetFormatAnalyzer {
	return &DotnetFormatAnalyzer{
		log: log,
	}
}

// Analyze выполняет статический анализ C# кода
func (a *DotnetFormatAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	a.log.Info(fmt.Sprintf("Running dotnet format analysis for project: %s", config.ProjectPath))

	startTime := time.Now()

	// Проверяем, что dotnet установлен
	if err := a.checkDotnetInstalled(); err != nil {
		return nil, fmt.Errorf("dotnet not installed: %w", err)
	}

	// Проверяем наличие .csproj или .sln файла
	if !a.hasDotnetProject(config.ProjectPath) {
		return nil, fmt.Errorf("no .csproj or .sln file found in project path")
	}

	// Строим команду для dotnet format
	args := []string{"format", "--verify-no-changes", "--verbosity", "diagnostic"}

	// Добавляем severity если указан
	if config.Severity != "" {
		args = append(args, "--severity", config.Severity)
	}

	// Создаем команду
	cmd := exec.CommandContext(ctx, "dotnet", args...)
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
		Analyzer:    domain.StaticAnalyzerTypeDotnetFormat,
		Duration:    duration,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}

	if err != nil {
		// dotnet format возвращает ошибку, если найдены проблемы форматирования
		a.log.Warning(fmt.Sprintf("Dotnet format found issues: %v", err))
		result.Success = false
	}

	// Парсим вывод dotnet format
	issues, parseErr := a.parseDotnetFormatOutput(string(output))
	if parseErr != nil {
		a.log.Warning(fmt.Sprintf("Failed to parse dotnet format output: %v", parseErr))
		result.Error = fmt.Sprintf("Failed to parse output: %v", parseErr)
	} else {
		result.Issues = issues
	}

	// Генерируем сводку
	result.Summary = generateSummary(issues)

	// Определяем успешность на основе наличия ошибок
	result.Success = result.Summary.ErrorCount == 0 && result.Summary.WarningCount == 0

	a.log.Info(fmt.Sprintf("Dotnet format analysis completed in %.2fs, found %d issues", duration, len(issues)))
	return result, nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (a *DotnetFormatAnalyzer) GetSupportedLanguages() []string {
	return []string{"csharp", "cs"}
}

// GetAnalyzerType возвращает тип анализатора
func (a *DotnetFormatAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return domain.StaticAnalyzerTypeDotnetFormat
}

// ValidateConfig проверяет корректность конфигурации
func (a *DotnetFormatAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if config.Language != "csharp" && config.Language != "cs" {
		return fmt.Errorf("dotnet format analyzer only supports C# language")
	}

	if config.ProjectPath == "" {
		return fmt.Errorf("project path is required")
	}

	// Проверяем наличие .csproj или .sln файла
	if !a.hasDotnetProject(config.ProjectPath) {
		return fmt.Errorf("no .csproj or .sln file found in project path")
	}

	return nil
}

// checkDotnetInstalled проверяет, установлен ли dotnet
func (a *DotnetFormatAnalyzer) checkDotnetInstalled() error {
	cmd := exec.Command("dotnet", "--version")
	executil.HideWindow(cmd)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("dotnet not found, install .NET SDK from https://dotnet.microsoft.com/download")
	}
	return nil
}

// hasDotnetProject проверяет наличие .csproj или .sln файла
func (a *DotnetFormatAnalyzer) hasDotnetProject(projectPath string) bool {
	// Проверяем .sln файлы
	slnFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.sln"))
	if len(slnFiles) > 0 {
		return true
	}

	// Проверяем .csproj файлы
	csprojFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.csproj"))
	if len(csprojFiles) > 0 {
		return true
	}

	// Проверяем .csproj в поддиректориях
	found := false
	_ = filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".csproj") {
			found = true
			return filepath.SkipAll
		}
		return nil
	})

	return found
}

// parseDotnetFormatOutput парсит вывод dotnet format
func (a *DotnetFormatAnalyzer) parseDotnetFormatOutput(output string) ([]*domain.StaticIssue, error) {
	var issues []*domain.StaticIssue

	lines := strings.Split(output, "\n")

	// Regex patterns for different dotnet format output formats
	// Format: "path/File.cs(10,5): warning IDE0001: Message"
	issueRe := regexp.MustCompile(`([^(]+)\((\d+),(\d+)\):\s+(error|warning|info)\s+(\w+):\s+(.+)`)

	// Format for formatting issues: "Would format: path/File.cs"
	formatRe := regexp.MustCompile(`Would format:\s+(.+)`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try to match standard issue format
		if matches := issueRe.FindStringSubmatch(line); len(matches) >= 7 {
			issue := &domain.StaticIssue{
				File:     matches[1],
				Line:     a.parseInt(matches[2]),
				Column:   a.parseInt(matches[3]),
				Severity: matches[4],
				Code:     matches[5],
				Message:  matches[6],
				Category: a.categorizeIssue(matches[5]),
			}
			issues = append(issues, issue)
			continue
		}

		// Try to match "Would format" output
		if matches := formatRe.FindStringSubmatch(line); len(matches) >= 2 {
			issue := &domain.StaticIssue{
				File:     matches[1],
				Severity: severityWarning,
				Code:     "FORMAT001",
				Message:  "File needs formatting",
				Category: "formatting",
			}
			issues = append(issues, issue)
		}
	}

	return issues, nil
}

// categorizeIssue категоризирует проблему по коду
func (a *DotnetFormatAnalyzer) categorizeIssue(code string) string {
	if code == "" {
		return severityOther
	}

	// IDE codes - code style
	if strings.HasPrefix(code, "IDE") {
		return "style"
	}

	// CS codes - compiler warnings/errors
	if strings.HasPrefix(code, "CS") {
		return "compiler"
	}

	// CA codes - code analysis
	if strings.HasPrefix(code, "CA") {
		return "analysis"
	}

	// SA codes - StyleCop
	if strings.HasPrefix(code, "SA") {
		return "style"
	}

	// FORMAT codes - formatting
	if strings.HasPrefix(code, "FORMAT") {
		return "formatting"
	}

	return severityOther
}

// parseInt parses string to int
func (a *DotnetFormatAnalyzer) parseInt(s string) int {
	var val int
	_, _ = fmt.Sscanf(s, "%d", &val)
	return val
}
