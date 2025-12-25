package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/application/protocol"
	"syntaxia/application/taskflow"
	"syntaxia/domain"
	"sync"
	"sync/atomic"
)

// TaskflowHandler handles taskflow and autonomous task operations
type TaskflowHandler struct {
	log                         domain.Logger
	taskflowService             domain.TaskflowService
	guardrailService            domain.GuardrailService
	repairService               domain.RepairService
	taskProtocolService         domain.TaskProtocolService
	taskProtocolConfigService   *protocol.ConfigService
	taskflowProtocolIntegration *taskflow.ProtocolIntegration
	buildService                domain.IBuildService

	// Active task tracking for cancellation
	activeTasks   map[string]context.CancelFunc
	activeTasksMu sync.Mutex

	// Metrics
	totalTasks      int64
	activeTaskCount int64
	failedTasks     int64
}

// NewTaskflowHandler creates a new taskflow handler
func NewTaskflowHandler(
	log domain.Logger,
	taskflowService domain.TaskflowService,
	guardrailService domain.GuardrailService,
	repairService domain.RepairService,
	taskProtocolService domain.TaskProtocolService,
	taskProtocolConfigService *protocol.ConfigService,
	taskflowProtocolIntegration *taskflow.ProtocolIntegration,
	buildService domain.IBuildService,
) *TaskflowHandler {
	return &TaskflowHandler{
		log:                         log,
		taskflowService:             taskflowService,
		guardrailService:            guardrailService,
		repairService:               repairService,
		taskProtocolService:         taskProtocolService,
		taskProtocolConfigService:   taskProtocolConfigService,
		taskflowProtocolIntegration: taskflowProtocolIntegration,
		buildService:                buildService,
		activeTasks:                 make(map[string]context.CancelFunc),
	}
}

// === Taskflow Operations ===

// LoadTasks loads tasks from plan.yaml
func (h *TaskflowHandler) LoadTasks() ([]domain.Task, error) {
	return h.taskflowService.LoadTasks()
}

// GetTaskStatus returns task status
func (h *TaskflowHandler) GetTaskStatus(taskID string) (*domain.TaskStatus, error) {
	return h.taskflowService.GetTaskStatus(taskID)
}

// UpdateTaskStatus updates task status
func (h *TaskflowHandler) UpdateTaskStatus(taskID string, state domain.TaskState, message string) error {
	return h.taskflowService.UpdateTaskStatus(taskID, state, message)
}

// ExecuteTask executes a task
func (h *TaskflowHandler) ExecuteTask(taskID string) error {
	atomic.AddInt64(&h.activeTaskCount, 1)
	defer atomic.AddInt64(&h.activeTaskCount, -1)
	atomic.AddInt64(&h.totalTasks, 1)

	err := h.taskflowService.ExecuteTask(taskID)
	if err != nil {
		atomic.AddInt64(&h.failedTasks, 1)
	}
	return err
}

// ExecuteTaskflow executes entire taskflow
func (h *TaskflowHandler) ExecuteTaskflow() error {
	return h.taskflowService.ExecuteTaskflow()
}

// GetReadyTasks returns ready tasks
func (h *TaskflowHandler) GetReadyTasks() ([]domain.Task, error) {
	return h.taskflowService.GetReadyTasks()
}

// GetTaskDependencies returns task dependencies
func (h *TaskflowHandler) GetTaskDependencies(taskID string) ([]domain.Task, error) {
	return h.taskflowService.GetTaskDependencies(taskID)
}

// ValidateTaskflow validates taskflow
func (h *TaskflowHandler) ValidateTaskflow() error {
	return h.taskflowService.ValidateTaskflow()
}

// GetTaskflowProgress returns taskflow progress
func (h *TaskflowHandler) GetTaskflowProgress() (float64, error) {
	return h.taskflowService.GetTaskflowProgress()
}

// ResetTaskflow resets taskflow
func (h *TaskflowHandler) ResetTaskflow() error {
	return h.taskflowService.ResetTaskflow()
}

// === Autonomous Task Operations ===

// StartAutonomousTask starts an autonomous task
func (h *TaskflowHandler) StartAutonomousTask(ctx context.Context, requestJSON string) (string, error) {
	var request domain.AutonomousTaskRequest
	if err := json.Unmarshal([]byte(requestJSON), &request); err != nil {
		return "", fmt.Errorf("invalid JSON request format: %w", err)
	}

	// Create cancellable context
	taskCtx, cancel := context.WithCancel(ctx)

	result, err := h.taskflowService.StartAutonomousTask(taskCtx, request)
	if err != nil {
		cancel()
		return "", err
	}

	// Store cancel function for later cancellation
	h.activeTasksMu.Lock()
	h.activeTasks[result.TaskId] = cancel
	h.activeTasksMu.Unlock()

	atomic.AddInt64(&h.activeTaskCount, 1)
	atomic.AddInt64(&h.totalTasks, 1)

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to serialize response: %w", err)
	}

	return string(resultJSON), nil
}

// CancelAutonomousTask cancels an autonomous task
func (h *TaskflowHandler) CancelAutonomousTask(ctx context.Context, taskID string) error {
	h.activeTasksMu.Lock()
	cancel, exists := h.activeTasks[taskID]
	if exists {
		cancel()
		delete(h.activeTasks, taskID)
	}
	h.activeTasksMu.Unlock()

	atomic.AddInt64(&h.activeTaskCount, -1)
	return h.taskflowService.CancelAutonomousTask(ctx, taskID)
}

// GetAutonomousTaskStatus returns autonomous task status
func (h *TaskflowHandler) GetAutonomousTaskStatus(ctx context.Context, taskID string) (string, error) {
	status, err := h.taskflowService.GetAutonomousTaskStatus(ctx, taskID)
	if err != nil {
		return "", err
	}

	statusJSON, err := json.Marshal(status)
	if err != nil {
		return "", fmt.Errorf("failed to marshal status: %w", err)
	}

	return string(statusJSON), nil
}

// ListAutonomousTasks lists autonomous tasks
func (h *TaskflowHandler) ListAutonomousTasks(ctx context.Context, projectPath string) (string, error) {
	tasks, err := h.taskflowService.ListAutonomousTasks(ctx, projectPath)
	if err != nil {
		return "", err
	}

	tasksJSON, err := json.Marshal(tasks)
	if err != nil {
		return "", fmt.Errorf("failed to marshal tasks: %w", err)
	}

	return string(tasksJSON), nil
}

// GetTaskLogs returns task logs
func (h *TaskflowHandler) GetTaskLogs(ctx context.Context, taskID string) (string, error) {
	logs, err := h.taskflowService.GetTaskLogs(ctx, taskID)
	if err != nil {
		return "", err
	}

	logsJSON, err := json.Marshal(logs)
	if err != nil {
		return "", fmt.Errorf("failed to marshal logs: %w", err)
	}

	return string(logsJSON), nil
}

// PauseTask pauses a task
func (h *TaskflowHandler) PauseTask(ctx context.Context, taskID string) error {
	return h.taskflowService.PauseTask(ctx, taskID)
}

// ResumeTask resumes a task
func (h *TaskflowHandler) ResumeTask(ctx context.Context, taskID string) error {
	return h.taskflowService.ResumeTask(ctx, taskID)
}

// === Guardrail Operations ===

// ValidatePath validates path against policies
func (h *TaskflowHandler) ValidatePath(path string) ([]domain.GuardrailViolation, error) {
	return h.guardrailService.ValidatePath(path)
}

// ValidateBudget validates budget constraints
func (h *TaskflowHandler) ValidateBudget(budgetType domain.BudgetType, current int64) ([]domain.BudgetViolation, error) {
	return h.guardrailService.ValidateBudget(budgetType, current)
}

// GetGuardrailPolicies returns all policies
func (h *TaskflowHandler) GetGuardrailPolicies() ([]domain.GuardrailPolicy, error) {
	return h.guardrailService.GetPolicies()
}

// GetBudgetPolicies returns budget policies
func (h *TaskflowHandler) GetBudgetPolicies() ([]domain.BudgetPolicy, error) {
	return h.guardrailService.GetBudgetPolicies()
}

// AddGuardrailPolicy adds a policy
func (h *TaskflowHandler) AddGuardrailPolicy(policy domain.GuardrailPolicy) error {
	return h.guardrailService.AddPolicy(policy)
}

// RemoveGuardrailPolicy removes a policy
func (h *TaskflowHandler) RemoveGuardrailPolicy(policyID string) error {
	return h.guardrailService.RemovePolicy(policyID)
}

// UpdateGuardrailPolicy updates a policy
func (h *TaskflowHandler) UpdateGuardrailPolicy(policy domain.GuardrailPolicy) error {
	return h.guardrailService.UpdatePolicy(policy)
}

// AddBudgetPolicy adds a budget policy
func (h *TaskflowHandler) AddBudgetPolicy(policy domain.BudgetPolicy) error {
	return h.guardrailService.AddBudgetPolicy(policy)
}

// RemoveBudgetPolicy removes a budget policy
func (h *TaskflowHandler) RemoveBudgetPolicy(policyID string) error {
	return h.guardrailService.RemoveBudgetPolicy(policyID)
}

// UpdateBudgetPolicy updates a budget policy
func (h *TaskflowHandler) UpdateBudgetPolicy(policy domain.BudgetPolicy) error {
	return h.guardrailService.UpdateBudgetPolicy(policy)
}

// === Repair Operations ===

// ExecuteRepair executes repair cycle
func (h *TaskflowHandler) ExecuteRepair(ctx context.Context, projectPath, errorOutput, language string, maxAttempts int) (*domain.RepairResult, error) {
	req := domain.RepairRequest{
		ProjectPath: projectPath,
		ErrorOutput: errorOutput,
		Language:    language,
		MaxAttempts: maxAttempts,
	}
	return h.repairService.ExecuteRepair(ctx, req)
}

// GetAvailableRepairRules returns available repair rules
func (h *TaskflowHandler) GetAvailableRepairRules(language string) ([]domain.RepairRule, error) {
	return h.repairService.GetAvailableRules(language)
}

// AddRepairRule adds a repair rule
func (h *TaskflowHandler) AddRepairRule(rule domain.RepairRule) error {
	return h.repairService.AddRule(rule)
}

// RemoveRepairRule removes a repair rule
func (h *TaskflowHandler) RemoveRepairRule(ruleID string) error {
	return h.repairService.RemoveRule(ruleID)
}

// ValidateRepairRule validates a repair rule
func (h *TaskflowHandler) ValidateRepairRule(rule domain.RepairRule) error {
	return h.repairService.ValidateRule(rule)
}

// GetMetrics returns handler metrics
func (h *TaskflowHandler) GetMetrics() map[string]interface{} {
	h.activeTasksMu.Lock()
	activeCount := len(h.activeTasks)
	h.activeTasksMu.Unlock()

	return map[string]interface{}{
		"total_tasks":         atomic.LoadInt64(&h.totalTasks),
		"active_tasks":        atomic.LoadInt64(&h.activeTaskCount),
		"failed_tasks":        atomic.LoadInt64(&h.failedTasks),
		"tracked_cancellable": activeCount,
	}
}

// Cleanup cleans up resources
func (h *TaskflowHandler) Cleanup() {
	h.activeTasksMu.Lock()
	defer h.activeTasksMu.Unlock()

	for taskID, cancel := range h.activeTasks {
		h.log.Info(fmt.Sprintf("Cancelling task %s during cleanup", taskID))
		cancel()
	}
	h.activeTasks = make(map[string]context.CancelFunc)
}
