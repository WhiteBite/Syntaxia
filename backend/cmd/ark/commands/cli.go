package commands

import (
	"context"
)

// CLI представляет интерфейс командной строки
type CLI struct {
	container *CLIContainer
}

// NewCLI создает новый экземпляр CLI
func NewCLI(container *CLIContainer) *CLI {
	return &CLI{
		container: container,
	}
}

// Index выполняет команду индексации проекта
func (c *CLI) Index(ctx context.Context, args []string) error {
	indexCmd := NewIndexCommand(c.container)
	return indexCmd.Execute(ctx, args)
}

// Solve выполняет команду решения задач
func (c *CLI) Solve(ctx context.Context, args []string) error {
	solveCmd := NewSolveCommand(c.container)
	return solveCmd.Execute(ctx, args)
}

// Result выполняет команду показа результатов
func (c *CLI) Result(ctx context.Context, args []string) error {
	resultCmd := NewResultCommand(c.container)
	return resultCmd.Execute(ctx, args)
}

// Verify выполняет команду верификации проекта
func (c *CLI) Verify(ctx context.Context, args []string) error {
	verifyCmd := NewVerifyCommand(c.container)
	return verifyCmd.Execute(ctx, args)
}
