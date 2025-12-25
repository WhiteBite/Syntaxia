package exec

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/internal/executil"
)

// CommandRunnerImpl реализует интерфейс CommandRunner для выполнения команд
type CommandRunnerImpl struct {
	log domain.Logger
}

// NewCommandRunnerImpl создает новый экземпляр CommandRunnerImpl
func NewCommandRunnerImpl(log domain.Logger) *CommandRunnerImpl {
	return &CommandRunnerImpl{
		log: log,
	}
}

// RunCommand выполняет команду с заданным контекстом и аргументами
func (c *CommandRunnerImpl) RunCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	c.log.Debug(fmt.Sprintf("Executing command: %s %v", name, args))

	cmd := exec.CommandContext(ctx, name, args...)
	executil.HideWindow(cmd)
	output, err := cmd.CombinedOutput()

	if err != nil {
		c.log.Warning(fmt.Sprintf("Command failed: %s %v - %v", name, args, err))
		return output, fmt.Errorf("command %s failed: %w", name, err)
	}

	c.log.Debug(fmt.Sprintf("Command succeeded: %s %v", name, args))
	return output, nil
}

// RunCommandInDir выполняет команду в указанной директории
func (c *CommandRunnerImpl) RunCommandInDir(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	c.log.Debug(fmt.Sprintf("Executing command in directory %s: %s %v", dir, name, args))

	// Проверяем, что директория существует
	if !filepath.IsAbs(dir) {
		return nil, fmt.Errorf("directory path must be absolute: %s", dir)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	executil.HideWindow(cmd)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if err != nil {
		c.log.Warning(fmt.Sprintf("Command failed in directory %s: %s %v - %v", dir, name, args, err))
		return output, fmt.Errorf("command %s failed in directory %s: %w", name, dir, err)
	}

	c.log.Debug(fmt.Sprintf("Command succeeded in directory %s: %s %v", dir, name, args))
	return output, nil
}
