package sandbox

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"strings"
	"time"
)

const defaultEngine = "docker"

// RunnerImpl реализует SandboxRunner
type RunnerImpl struct {
	log domain.Logger
}

// NewSandboxRunner создает новый sandbox runner
func NewSandboxRunner(log domain.Logger) *RunnerImpl {
	return &RunnerImpl{
		log: log,
	}
}

// Run выполняет команду в песочнице
func (r *RunnerImpl) Run(ctx context.Context, config domain.SandboxConfig, command []string) (*domain.SandboxResult, error) {
	startTime := time.Now()
	result := &domain.SandboxResult{}

	// Проверяем доступность движка
	if !r.IsAvailable(ctx) {
		result.Success = false
		result.Error = fmt.Sprintf("Sandbox engine %s not available", config.Engine)
		return result, nil
	}

	// Создаем контейнер
	containerID, err := r.createContainer(ctx, config, command)
	if err != nil {
		result.Success = false
		result.Error = fmt.Sprintf("Failed to create container: %v", err)
		return result, nil
	}
	result.ContainerID = containerID

	// Запускаем контейнер с таймаутом
	runCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Timeout)*time.Second)
	defer cancel()

	output, err := r.runContainer(runCtx, containerID)
	result.Output = output

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.ExitCode = 1
	} else {
		result.Success = true
		result.ExitCode = 0
	}

	// Получаем логи
	logs, err := r.getContainerLogs(ctx, containerID)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Failed to get container logs: %v", err))
	} else {
		result.Logs = logs
	}

	// Очищаем ресурсы
	if err := r.Cleanup(ctx, containerID); err != nil {
		r.log.Warning(fmt.Sprintf("Failed to cleanup container %s: %v", containerID, err))
	}

	result.Duration = time.Since(startTime).Seconds()
	return result, nil
}

// Cleanup очищает ресурсы песочницы
func (r *RunnerImpl) Cleanup(ctx context.Context, containerID string) error {
	if containerID == "" {
		return nil
	}

	// Останавливаем контейнер
	stopCmd := exec.CommandContext(ctx, "docker", "stop", containerID)
	executil.HideWindow(stopCmd)
	if err := stopCmd.Run(); err != nil {
		r.log.Warning(fmt.Sprintf("Failed to stop container %s: %v", containerID, err))
	}

	// Удаляем контейнер
	rmCmd := exec.CommandContext(ctx, "docker", "rm", containerID)
	executil.HideWindow(rmCmd)
	if err := rmCmd.Run(); err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}

	return nil
}

// IsAvailable проверяет доступность движка песочницы
func (r *RunnerImpl) IsAvailable(ctx context.Context) bool {
	// Проверяем наличие Docker
	if _, err := exec.LookPath("docker"); err == nil {
		// Проверяем, что Docker daemon работает
		cmd := exec.CommandContext(ctx, "docker", "info")
		executil.HideWindow(cmd)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	// Проверяем наличие Podman
	if _, err := exec.LookPath("podman"); err == nil {
		// Проверяем, что Podman работает
		cmd := exec.CommandContext(ctx, "podman", "info")
		executil.HideWindow(cmd)
		if err := cmd.Run(); err == nil {
			return true
		}
	}

	return false
}

// GetInfo возвращает информацию о движке
func (r *RunnerImpl) GetInfo(ctx context.Context) (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// Проверяем Docker
	if _, err := exec.LookPath("docker"); err == nil {
		cmd := exec.CommandContext(ctx, "docker", "version", "--format", "json")
		executil.HideWindow(cmd)
		output, err := cmd.Output()
		if err == nil {
			var dockerInfo map[string]interface{}
			if err := json.Unmarshal(output, &dockerInfo); err == nil {
				info["docker"] = dockerInfo
			}
		}
	}

	// Проверяем Podman
	if _, err := exec.LookPath("podman"); err == nil {
		cmd := exec.CommandContext(ctx, "podman", "version", "--format", "json")
		executil.HideWindow(cmd)
		output, err := cmd.Output()
		if err == nil {
			var podmanInfo map[string]interface{}
			if err := json.Unmarshal(output, &podmanInfo); err == nil {
				info["podman"] = podmanInfo
			}
		}
	}

	return info, nil
}

// createContainer создает контейнер
func (r *RunnerImpl) createContainer(ctx context.Context, config domain.SandboxConfig, command []string) (string, error) {
	engine := defaultEngine
	if config.Engine != "" {
		engine = config.Engine
	}

	args := []string{"create"}

	// Добавляем ограничения ресурсов
	if config.MemoryLimit != "" {
		args = append(args, "--memory", config.MemoryLimit)
	}
	if config.CPULimit != "" {
		args = append(args, "--cpus", config.CPULimit)
	}

	// Добавляем режим сети
	if config.NetworkMode != "" {
		args = append(args, "--network", config.NetworkMode)
	} else {
		args = append(args, "--network", "none")
	}

	// Добавляем точки монтирования
	for _, mount := range config.Mounts {
		mountArg := fmt.Sprintf("%s:%s", mount.Source, mount.Target)
		if mount.ReadOnly {
			mountArg += ":ro"
		}
		args = append(args, "--mount", mountArg)
	}

	// Добавляем переменные окружения
	for key, value := range config.EnvVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Добавляем пользователя
	if config.User != "" {
		args = append(args, "--user", config.User)
	}

	// Добавляем рабочую директорию
	if config.WorkingDir != "" {
		args = append(args, "--workdir", config.WorkingDir)
	}

	// Добавляем образ и команду
	args = append(args, config.Image)
	args = append(args, command...)

	cmd := exec.CommandContext(ctx, engine, args...)
	executil.HideWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	containerID := strings.TrimSpace(string(output))
	return containerID, nil
}

// runContainer запускает контейнер
func (r *RunnerImpl) runContainer(ctx context.Context, containerID string) (string, error) {
	engine := defaultEngine

	cmd := exec.CommandContext(ctx, engine, "start", "-a", containerID)
	executil.HideWindow(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("container execution failed: %w", err)
	}

	return string(output), nil
}

// getContainerLogs получает логи контейнера
func (r *RunnerImpl) getContainerLogs(ctx context.Context, containerID string) (string, error) {
	engine := defaultEngine

	cmd := exec.CommandContext(ctx, engine, "logs", containerID)
	executil.HideWindow(cmd)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}

	return string(output), nil
}
