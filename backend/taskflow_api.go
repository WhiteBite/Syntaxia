package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"syntaxia/application/taskflow"
	"syntaxia/domain"
	"time"
)

// LoadTasks loads tasks from plan.yaml
func (a *App) LoadTasks() ([]domain.Task, error) {
	return a.taskflowService.LoadTasks()
}

// GetTaskStatus returns task status
func (a *App) GetTaskStatus(taskID string) (*domain.TaskStatus, error) {
	return a.taskflowService.GetTaskStatus(taskID)
}

// UpdateTaskStatus updates task status
func (a *App) UpdateTaskStatus(taskID string, state domain.TaskState, message string) error {
	return a.taskflowService.UpdateTaskStatus(taskID, state, message)
}

// ExecuteTask executes a task
func (a *App) ExecuteTask(taskID string) error {
	return a.taskflowService.ExecuteTask(taskID)
}

// ExecuteTaskflow executes the entire taskflow
func (a *App) ExecuteTaskflow() error {
	return a.taskflowService.ExecuteTaskflow()
}

// GetReadyTasks returns tasks ready for execution
func (a *App) GetReadyTasks() ([]domain.Task, error) {
	return a.taskflowService.GetReadyTasks()
}

// GetTaskDependencies returns task dependencies
func (a *App) GetTaskDependencies(taskID string) ([]domain.Task, error) {
	return a.taskflowService.GetTaskDependencies(taskID)
}

// ValidateTaskflow validates taskflow correctness
func (a *App) ValidateTaskflow() error {
	return a.taskflowService.ValidateTaskflow()
}

// GetTaskflowProgress returns execution progress
func (a *App) GetTaskflowProgress() (float64, error) {
	return a.taskflowService.GetTaskflowProgress()
}

// ResetTaskflow resets taskflow
func (a *App) ResetTaskflow() error {
	return a.taskflowService.ResetTaskflow()
}

// StartAutonomousTask starts an autonomous task
func (a *App) StartAutonomousTask(requestJson string) (string, error) {
	var request domain.AutonomousTaskRequest
	if err := json.Unmarshal([]byte(requestJson), &request); err != nil {
		validationErr := domain.NewValidationError("Invalid JSON request format", map[string]interface{}{
			"originalError": err.Error(),
			"requestJson":   requestJson,
		})
		return "", a.transformDomainError(validationErr)
	}

	result, err := a.taskflowService.StartAutonomousTask(a.ctx, request)
	if err != nil {
		if domainErr, ok := err.(*domain.DomainError); ok {
			return "", a.transformDomainError(domainErr)
		}
		return "", domain.NewInternalError("Unexpected error occurred", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		configErr := domain.NewConfigurationError("Failed to serialize response", err)
		return "", a.transformDomainError(configErr)
	}

	return string(resultJson), nil
}

// CancelAutonomousTask cancels an autonomous task
func (a *App) CancelAutonomousTask(taskId string) error {
	return a.taskflowService.CancelAutonomousTask(a.ctx, taskId)
}

// GetAutonomousTaskStatus gets autonomous task status
func (a *App) GetAutonomousTaskStatus(taskId string) (string, error) {
	status, err := a.taskflowService.GetAutonomousTaskStatus(a.ctx, taskId)
	if err != nil {
		return "", a.transformError(err)
	}

	statusJson, err := json.Marshal(status)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal status", err)
		return "", a.transformError(marshalErr)
	}

	return string(statusJson), nil
}

// ListAutonomousTasks lists all autonomous tasks
func (a *App) ListAutonomousTasks(projectPath string) (string, error) {
	tasks, err := a.taskflowService.ListAutonomousTasks(a.ctx, projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to list autonomous tasks: %w", err)
	}

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tasks: %w", err)
	}

	return string(tasksJson), nil
}

// GetTaskLogs returns logs for a specific task
func (a *App) GetTaskLogs(taskId string) (string, error) {
	logs, err := a.taskflowService.GetTaskLogs(a.ctx, taskId)
	if err != nil {
		return "", fmt.Errorf("failed to get task logs: %w", err)
	}

	logsJson, err := json.Marshal(logs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal logs: %w", err)
	}

	return string(logsJson), nil
}

// PauseTask pauses an autonomous task
func (a *App) PauseTask(taskId string) error {
	return a.taskflowService.PauseTask(a.ctx, taskId)
}

// ResumeTask resumes a paused autonomous task
func (a *App) ResumeTask(taskId string) error {
	return a.taskflowService.ResumeTask(a.ctx, taskId)
}

// ============ TASK PROTOCOL API ENDPOINTS ============

// ExecuteTaskProtocol executes the full Task Protocol verification for a project
func (a *App) ExecuteTaskProtocol(configJson string) (string, error) {
	var config domain.TaskProtocolConfig
	if err := json.Unmarshal([]byte(configJson), &config); err != nil {
		return "", fmt.Errorf("failed to parse protocol config JSON: %w", err)
	}

	result, err := a.taskProtocolService.ExecuteProtocol(a.ctx, &config)
	if err != nil {
		return "", fmt.Errorf("protocol execution failed: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protocol result: %w", err)
	}

	return string(resultJson), nil
}

// ValidateTaskProtocolStage executes a single Task Protocol verification stage
func (a *App) ValidateTaskProtocolStage(stage string, configJson string) (string, error) {
	var config domain.TaskProtocolConfig
	if err := json.Unmarshal([]byte(configJson), &config); err != nil {
		return "", fmt.Errorf("failed to parse protocol config JSON: %w", err)
	}

	var protocolStage domain.ProtocolStage
	switch stage {
	case "linting":
		protocolStage = domain.StageLinting
	case "building":
		protocolStage = domain.StageBuilding
	case "testing":
		protocolStage = domain.StageTesting
	case "guardrails":
		protocolStage = domain.StageGuardrails
	default:
		return "", fmt.Errorf("unsupported protocol stage: %s", stage)
	}

	result, err := a.taskProtocolService.ValidateStage(a.ctx, protocolStage, &config)
	if err != nil {
		return "", fmt.Errorf("stage validation failed: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal stage result: %w", err)
	}

	return string(resultJson), nil
}

// ValidateAIGeneratedCode validates AI-generated code using Task Protocol
func (a *App) ValidateAIGeneratedCode(requestJson string) (string, error) {
	var request struct {
		ProjectPath    string   `json:"projectPath"`
		Context        string   `json:"context"`
		Languages      []string `json:"languages"`
		GeneratedFiles []struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		} `json:"generatedFiles"`
	}

	if err := json.Unmarshal([]byte(requestJson), &request); err != nil {
		return "", fmt.Errorf("failed to parse validation request JSON: %w", err)
	}

	validationRequest := &taskflow.AICodeValidationRequest{
		ProjectPath: request.ProjectPath,
		Context:     request.Context,
		Languages:   request.Languages,
	}

	result, err := a.taskflowProtocolIntegration.ValidateAIGeneratedCode(a.ctx, validationRequest)
	if err != nil {
		return "", fmt.Errorf("AI code validation failed: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal validation result: %w", err)
	}

	return string(resultJson), nil
}

// GetTaskProtocolConfiguration loads the current Task Protocol configuration
func (a *App) GetTaskProtocolConfiguration(projectPath string, languages []string) (string, error) {
	config := &domain.TaskProtocolConfig{
		ProjectPath: projectPath,
		Languages:   languages,
		EnabledStages: []domain.ProtocolStage{
			domain.StageLinting,
			domain.StageBuilding,
			domain.StageTesting,
			domain.StageGuardrails,
		},
		MaxRetries: 3,
		FailFast:   false,
		SelfCorrection: domain.SelfCorrectionConfig{
			Enabled:      true,
			MaxAttempts:  5,
			AIAssistance: true,
		},
		Timeouts: map[string]time.Duration{
			"linting":    5 * time.Minute,
			"building":   10 * time.Minute,
			"testing":    15 * time.Minute,
			"guardrails": 2 * time.Minute,
		},
	}

	configPath := filepath.Join(projectPath, "task_protocol.yaml")
	if loadedConfig, err := a.taskProtocolConfigService.LoadConfiguration(configPath); err == nil {
		config = loadedConfig
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protocol configuration: %w", err)
	}

	return string(configJson), nil
}

// UpdateTaskProtocolConfiguration saves an updated Task Protocol configuration
func (a *App) UpdateTaskProtocolConfiguration(configJson string) error {
	var config domain.TaskProtocolConfig
	if err := json.Unmarshal([]byte(configJson), &config); err != nil {
		return fmt.Errorf("failed to parse protocol config JSON: %w", err)
	}

	if err := a.taskProtocolConfigService.ValidateConfiguration(&config); err != nil {
		return fmt.Errorf("configuration validation failed: %w", err)
	}

	configPath := filepath.Join(config.ProjectPath, "task_protocol.yaml")
	if err := a.taskProtocolConfigService.SaveConfiguration(&config, configPath); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	a.log.Info(fmt.Sprintf("Task Protocol configuration updated for project: %s", config.ProjectPath))
	return nil
}

// RequestTaskProtocolCorrectionGuidance requests AI guidance for correcting errors
func (a *App) RequestTaskProtocolCorrectionGuidance(errorDetailsJson string, contextJson string) (string, error) {
	var errorDetails domain.ErrorDetails
	if err := json.Unmarshal([]byte(errorDetailsJson), &errorDetails); err != nil {
		return "", fmt.Errorf("failed to parse error details JSON: %w", err)
	}

	var taskContext domain.TaskContext
	if err := json.Unmarshal([]byte(contextJson), &taskContext); err != nil {
		return "", fmt.Errorf("failed to parse task context JSON: %w", err)
	}

	guidance, err := a.taskProtocolService.RequestCorrectionGuidance(a.ctx, &errorDetails, &taskContext)
	if err != nil {
		return "", fmt.Errorf("correction guidance request failed: %w", err)
	}

	guidanceJson, err := json.Marshal(guidance)
	if err != nil {
		return "", fmt.Errorf("failed to marshal correction guidance: %w", err)
	}

	return string(guidanceJson), nil
}

// CreateTaskProtocolForProject creates a default Task Protocol configuration for a project
func (a *App) CreateTaskProtocolForProject(projectPath string, languages []string) (string, error) {
	if len(languages) == 0 {
		detectedLanguages, err := a.buildService.DetectLanguages(a.ctx, projectPath)
		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to detect languages for %s: %v", projectPath, err))
			languages = []string{"go"}
		} else {
			languages = detectedLanguages
		}
	}

	config := &domain.TaskProtocolConfig{
		ProjectPath: projectPath,
		Languages:   languages,
		EnabledStages: []domain.ProtocolStage{
			domain.StageLinting,
			domain.StageBuilding,
			domain.StageTesting,
			domain.StageGuardrails,
		},
		MaxRetries: 3,
		FailFast:   false,
		SelfCorrection: domain.SelfCorrectionConfig{
			Enabled:      true,
			MaxAttempts:  5,
			AIAssistance: true,
		},
		Timeouts: map[string]time.Duration{
			"linting":    5 * time.Minute,
			"building":   10 * time.Minute,
			"testing":    15 * time.Minute,
			"guardrails": 2 * time.Minute,
		},
	}

	configPath := filepath.Join(projectPath, "task_protocol.yaml")
	if err := a.taskProtocolConfigService.SaveConfiguration(config, configPath); err != nil {
		return "", fmt.Errorf("failed to save protocol configuration: %w", err)
	}

	configJson, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal protocol configuration: %w", err)
	}

	a.log.Info(fmt.Sprintf("Task Protocol configuration created for project: %s", projectPath))
	return string(configJson), nil
}
