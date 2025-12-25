package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"syntaxia/cmd/ark/commands"
)

const (
	appName = "ark"
	version = "1.0.0"
)

func main() {
	// Парсим флаги
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	// Показываем версию если запрошено
	if showVersion {
		fmt.Printf("%s version %s\n", appName, version)
		os.Exit(0)
	}

	// Получаем аргументы командной строки
	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// Создаем контекст
	ctx := context.Background()

	// Создаем CLI контейнер
	container, err := commands.NewCLIContainer(ctx, "", "", false)
	if err != nil {
		log.Fatalf("Failed to create CLI container: %v", err)
	}

	// Создаем CLI команды
	cli := commands.NewCLI(container)

	// Выполняем команду
	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "index":
		if err := cli.Index(ctx, commandArgs); err != nil {
			log.Fatalf("Index command failed: %v", err)
		}
	case "solve":
		if err := cli.Solve(ctx, commandArgs); err != nil {
			log.Fatalf("Solve command failed: %v", err)
		}
	case "result":
		if err := cli.Result(ctx, commandArgs); err != nil {
			log.Fatalf("Result command failed: %v", err)
		}
	case "verify":
		if err := cli.Verify(ctx, commandArgs); err != nil {
			log.Fatalf("Verify command failed: %v", err)
		}
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf(`%s - ARK/Syntaxia CLI

Usage: %s <command> [options]

Commands:
  index   - Index project files and build symbol graph
  solve   - Solve coding tasks using AI
  result  - Show results and reports
  verify  - Verify project quality and health
  help    - Show this help message

Examples:
  %s index --project ./my-project
  %s solve --task "add error handling"
  %s result --format json

Use '%s <command> --help' for more information about a command.
`, appName, appName, appName, appName, appName, appName)
}
