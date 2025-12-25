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

// SolveCommand представляет команду решения задач
type SolveCommand struct {
	container *CLIContainer
}

// NewSolveCommand создает новую команду решения задач
func NewSolveCommand(container *CLIContainer) *SolveCommand {
	return &SolveCommand{
		container: container,
	}
}

// Execute выполняет команду решения задач
func (c *SolveCommand) Execute(ctx context.Context, args []string) error {
	// Создаем флаги для команды
	fs := flag.NewFlagSet("solve", flag.ExitOnError)
	var (
		task        = fs.String("task", "", "Task description to solve")
		projectPath = fs.String("project", ".", "Project path")
		output      = fs.String("output", "", "Output file for solution (JSON)")
		provider    = fs.String("provider", "openai", "AI provider (openai, gemini, localai)")
		model       = fs.String("model", "", "AI model to use")
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

	// Проверяем обязательные параметры
	if *task == "" {
		return fmt.Errorf("task description is required (use -task flag)")
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
		fmt.Printf("Solving task: %s\n", *task)
		fmt.Printf("Project path: %s\n", absPath)
		fmt.Printf("AI provider: %s\n", *provider)
		if *model != "" {
			fmt.Printf("AI model: %s\n", *model)
		}
	}

	// Получаем настройки (пока не используем, но могут понадобиться в будущем)
	_, err = c.container.SettingsService.GetSettingsDTO()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}

	// Создаем системный промпт
	systemPrompt := c.createSystemPrompt(absPath, *provider, *model)

	if *verbose {
		fmt.Printf("System prompt length: %d characters\n", len(systemPrompt))
	}

	// Генерируем код
	generatedCode, err := c.container.AIService.GenerateCode(ctx, systemPrompt, *task)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	if *verbose {
		fmt.Printf("Generated code length: %d characters\n", len(generatedCode))
	}

	// Создаем результат решения
	solveResult := &SolveResult{
		Task:          *task,
		ProjectPath:   absPath,
		Provider:      *provider,
		Model:         *model,
		GeneratedCode: generatedCode,
		Timestamp:     time.Now(),
	}

	// Выводим результат
	if *output != "" {
		// Сохраняем в файл
		data, err := json.MarshalIndent(solveResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal solve result: %w", err)
		}

		if err := os.WriteFile(*output, data, 0o644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Solution saved to: %s\n", *output)
	} else {
		// Выводим в stdout
		data, err := json.MarshalIndent(solveResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal solve result: %w", err)
		}

		fmt.Println(string(data))
	}

	return nil
}

// createSystemPrompt создает системный промпт для AI
func (c *SolveCommand) createSystemPrompt(projectPath, provider, model string) string {
	prompt := fmt.Sprintf(`You are an expert software developer working on a project at: %s

You have access to the following AI provider: %s`, projectPath, provider)

	if model != "" {
		prompt += fmt.Sprintf("\nUsing model: %s", model)
	}

	prompt += `

Your task is to:
1. Analyze the codebase structure
2. Understand the existing patterns and conventions
3. Generate high-quality, production-ready code
4. Follow best practices for the language and framework
5. Ensure code is well-documented and maintainable

Please provide:
- Clear, readable code
- Appropriate comments and documentation
- Error handling where necessary
- Tests if applicable
- Any additional context or explanations

Focus on writing code that integrates well with the existing codebase.`

	return prompt
}

// printHelp выводит справку по команде
func (c *SolveCommand) printHelp() {
	fmt.Printf(`ark solve - Solve coding tasks using AI

Usage: ark solve [options]

Options:
  -task string
        Task description to solve (required)
  -project string
        Project path (default ".")
  -output string
        Output file for solution (JSON)
  -provider string
        AI provider: openai, gemini, localai (default "openai")
  -model string
        AI model to use (uses default if not specified)
  -verbose
        Verbose output
  -help
        Show this help message

Examples:
  ark solve --task "add error handling to login function"
  ark solve --task "implement user authentication" --project ./my-app
  ark solve --task "add unit tests" --provider gemini --output solution.json
  ark solve --task "refactor database queries" --verbose
`)
}

// SolveResult представляет результат решения задачи
type SolveResult struct {
	Task          string    `json:"task"`
	ProjectPath   string    `json:"project_path"`
	Provider      string    `json:"provider"`
	Model         string    `json:"model"`
	GeneratedCode string    `json:"generated_code"`
	Timestamp     time.Time `json:"timestamp"`
}
