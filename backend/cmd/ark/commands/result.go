package commands

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ResultCommand представляет команду показа результатов
type ResultCommand struct {
	container *CLIContainer
}

// NewResultCommand создает новую команду показа результатов
func NewResultCommand(container *CLIContainer) *ResultCommand {
	return &ResultCommand{
		container: container,
	}
}

// Execute выполняет команду показа результатов
func (c *ResultCommand) Execute(ctx context.Context, args []string) error {
	// Создаем флаги для команды
	fs := flag.NewFlagSet("result", flag.ExitOnError)
	var (
		projectPath = fs.String("project", ".", "Project path")
		format      = fs.String("format", "json", "Output format (json, text)")
		output      = fs.String("output", "", "Output file")
		reportType  = fs.String("type", "all", "Report type (all, ux, guardrails, tasks)")
		verbose     = fs.Bool("verbose", false, "Verbose output")
		help        = fs.Bool("help", false, "Show help")
	)

	// Парсим аргументы
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Показываем помощь если запрошено
	if *help {
		c.printHelp()
		return nil
	}

	// Проверяем существование проекта
	if _, err := os.Stat(*projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project path does not exist: %s", *projectPath)
	}

	// Получаем абсолютный путь
	absPath, err := filepath.Abs(*projectPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if *verbose {
		fmt.Printf("Generating results for project: %s\n", absPath)
		fmt.Printf("Report type: %s\n", *reportType)
		fmt.Printf("Output format: %s\n", *format)
	}

	// Создаем результат
	result := &ResultData{
		ProjectPath: absPath,
		ReportType:  *reportType,
		Timestamp:   time.Now(),
	}

	// Собираем данные в зависимости от типа отчета
	switch *reportType {
	case "all", "ux":
		if err := c.collectUXMetrics(result); err != nil {
			return fmt.Errorf("failed to collect UX metrics: %w", err)
		}
	}

	if *reportType == "all" || *reportType == "guardrails" {
		if err := c.collectGuardrailInfo(result); err != nil {
			return fmt.Errorf("failed to collect guardrail info: %w", err)
		}
	}

	if *reportType == "all" || *reportType == "tasks" {
		if err := c.collectTaskInfo(result); err != nil {
			return fmt.Errorf("failed to collect task info: %w", err)
		}
	}

	// Выводим результат
	if *output != "" {
		// Сохраняем в файл
		var data []byte
		var err error

		if *format == "json" {
			data, err = json.MarshalIndent(result, "", "  ")
		} else {
			data = []byte(c.formatAsText(result))
		}

		if err != nil {
			return fmt.Errorf("failed to marshal result: %w", err)
		}

		if err := os.WriteFile(*output, data, 0o644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Results saved to: %s\n", *output)
	} else {
		// Выводим в stdout
		if *format == "json" {
			data, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal result: %w", err)
			}
			fmt.Println(string(data))
		} else {
			fmt.Println(c.formatAsText(result))
		}
	}

	return nil
}

// collectUXMetrics собирает UX метрики
func (c *ResultCommand) collectUXMetrics(result *ResultData) error {
	// Здесь можно добавить сбор UX метрик
	// Пока возвращаем заглушку
	result.UXMetrics = &UXMetricsData{
		WhyViewReports:     []string{},
		TimeToGreenReports: []string{},
		DerivedDiffReports: []string{},
	}
	return nil
}

// collectGuardrailInfo собирает информацию о guardrails
func (c *ResultCommand) collectGuardrailInfo(result *ResultData) error {
	// Получаем конфигурацию guardrails
	config := c.container.GuardrailService.GetConfig()

	result.Guardrails = &GuardrailData{
		FailClosed:           config.FailClosed,
		EnableEphemeralMode:  config.EnableEphemeralMode,
		EnableTaskValidation: config.EnableTaskValidation,
		EnableBudgetTracking: config.EnableBudgetTracking,
		EnablePathValidation: config.EnablePathValidation,
	}
	return nil
}

// collectTaskInfo собирает информацию о задачах
func (c *ResultCommand) collectTaskInfo(result *ResultData) error {
	// Здесь можно добавить сбор информации о задачах
	// Пока возвращаем заглушку
	result.Tasks = &TaskData{
		TotalTasks: 0,
		Completed:  0,
		Failed:     0,
		Pending:    0,
	}
	return nil
}

// formatAsText форматирует результат как текст
func (c *ResultCommand) formatAsText(result *ResultData) string {
	text := "ARK Results Report\n"
	text += "==================\n\n"
	text += fmt.Sprintf("Project Path: %s\n", result.ProjectPath)
	text += fmt.Sprintf("Report Type: %s\n", result.ReportType)
	text += fmt.Sprintf("Generated: %s\n\n", result.Timestamp.Format(time.RFC3339))

	if result.Guardrails != nil {
		text += "Guardrails Configuration:\n"
		text += fmt.Sprintf("  Fail Closed: %t\n", result.Guardrails.FailClosed)
		text += fmt.Sprintf("  Ephemeral Mode: %t\n", result.Guardrails.EnableEphemeralMode)
		text += fmt.Sprintf("  Task Validation: %t\n", result.Guardrails.EnableTaskValidation)
		text += fmt.Sprintf("  Budget Tracking: %t\n", result.Guardrails.EnableBudgetTracking)
		text += fmt.Sprintf("  Path Validation: %t\n\n", result.Guardrails.EnablePathValidation)
	}

	if result.Tasks != nil {
		text += fmt.Sprintf("Task Status:\n")
		text += fmt.Sprintf("  Total: %d\n", result.Tasks.TotalTasks)
		text += fmt.Sprintf("  Completed: %d\n", result.Tasks.Completed)
		text += fmt.Sprintf("  Failed: %d\n", result.Tasks.Failed)
		text += fmt.Sprintf("  Pending: %d\n\n", result.Tasks.Pending)
	}

	return text
}

// printHelp выводит справку по команде
func (c *ResultCommand) printHelp() {
	fmt.Printf(`ark result - Show results and reports

Usage: ark result [options]

Options:
  -project string
        Project path (default ".")
  -format string
        Output format: json, text (default "json")
  -output string
        Output file (stdout if not specified)
  -type string
        Report type: all, ux, guardrails, tasks (default "all")
  -verbose
        Verbose output
  -help
        Show this help message

Examples:
  ark result --project ./my-project
  ark result --format text --type guardrails
  ark result --output report.json --type all
  ark result --verbose
`)
}

// ResultData представляет данные результатов
type ResultData struct {
	ProjectPath string         `json:"project_path"`
	ReportType  string         `json:"report_type"`
	Timestamp   time.Time      `json:"timestamp"`
	UXMetrics   *UXMetricsData `json:"ux_metrics,omitempty"`
	Guardrails  *GuardrailData `json:"guardrails,omitempty"`
	Tasks       *TaskData      `json:"tasks,omitempty"`
}

// UXMetricsData представляет UX метрики
type UXMetricsData struct {
	WhyViewReports     []string `json:"why_view_reports"`
	TimeToGreenReports []string `json:"time_to_green_reports"`
	DerivedDiffReports []string `json:"derived_diff_reports"`
}

// GuardrailData представляет данные guardrails
type GuardrailData struct {
	FailClosed           bool `json:"fail_closed"`
	EnableEphemeralMode  bool `json:"enable_ephemeral_mode"`
	EnableTaskValidation bool `json:"enable_task_validation"`
	EnableBudgetTracking bool `json:"enable_budget_tracking"`
	EnablePathValidation bool `json:"enable_path_validation"`
}

// TaskData представляет данные задач
type TaskData struct {
	TotalTasks int `json:"total_tasks"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Pending    int `json:"pending"`
}
