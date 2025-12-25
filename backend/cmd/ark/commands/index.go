package commands

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"time"
)

// IndexCommand представляет команду индексации
type IndexCommand struct {
	container *CLIContainer
}

// NewIndexCommand создает новую команду индексации
func NewIndexCommand(container *CLIContainer) *IndexCommand {
	return &IndexCommand{
		container: container,
	}
}

// Execute выполняет команду индексации
func (c *IndexCommand) Execute(ctx context.Context, args []string) error {
	// Создаем флаги для команды
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	var (
		projectPath = fs.String("project", ".", "Project path to index")
		output      = fs.String("output", "", "Output file for index data (JSON)")
		language    = fs.String("language", "go", "Primary language to index")
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
		fmt.Printf("Indexing project: %s\n", absPath)
		fmt.Printf("Language: %s\n", *language)
	}

	// Индексируем файлы проекта
	files, err := c.container.ProjectService.ListFiles(absPath, true, true)
	if err != nil {
		return fmt.Errorf("failed to list project files: %w", err)
	}

	if *verbose {
		fmt.Printf("Found %d files\n", len(files))
	}

	// Строим граф символов
	symbolGraph, err := c.container.SymbolGraph.BuildSymbolGraph(ctx, absPath, *language)
	if err != nil {
		return fmt.Errorf("failed to build symbol graph: %w", err)
	}

	if *verbose {
		fmt.Printf("Built symbol graph with %d nodes\n", len(symbolGraph.Nodes))
	}

	// Создаем результат индексации
	indexResult := &IndexResult{
		ProjectPath: absPath,
		Language:    *language,
		Files:       files,
		SymbolGraph: symbolGraph,
		Timestamp:   time.Now(),
	}

	// Выводим результат
	if *output != "" {
		// Сохраняем в файл
		data, err := json.MarshalIndent(indexResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal index result: %w", err)
		}

		if err := os.WriteFile(*output, data, 0o644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}

		fmt.Printf("Index data saved to: %s\n", *output)
	} else {
		// Выводим в stdout
		data, err := json.MarshalIndent(indexResult, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal index result: %w", err)
		}

		fmt.Println(string(data))
	}

	return nil
}

// printHelp выводит справку по команде
func (c *IndexCommand) printHelp() {
	fmt.Printf(`ark index - Index project files and build symbol graph

Usage: ark index [options]

Options:
  -project string
        Project path to index (default ".")
  -output string
        Output file for index data (JSON)
  -language string
        Primary language to index (default "go")
  -verbose
        Verbose output
  -help
        Show this help message

Examples:
  ark index --project ./my-project
  ark index --project ./my-project --output index.json --language typescript
  ark index --project ./my-project --verbose
`)
}

// IndexResult представляет результат индексации
type IndexResult struct {
	ProjectPath string              `json:"project_path"`
	Language    string              `json:"language"`
	Files       []*domain.FileNode  `json:"files"`
	SymbolGraph *domain.SymbolGraph `json:"symbol_graph"`
	Timestamp   time.Time           `json:"timestamp"`
}
