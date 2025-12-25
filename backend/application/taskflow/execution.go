package taskflow

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"time"
)

// safeExecuteAutonomousTask executes autonomous task with comprehensive error recovery
func (s *Service) safeExecuteAutonomousTask(ctx context.Context, request domain.AutonomousTaskRequest, status *domain.AutonomousTaskStatus) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error(fmt.Sprintf("PANIC in autonomous task execution: %v", r))
			s.updateAutonomousTaskStatus(status.TaskId, "failed",
				fmt.Sprintf("Task execution panicked: %v", r), 100.0)
			s.notifyTaskFailure(status.TaskId, fmt.Sprintf("Internal error: %v", r))
		}
	}()

	if err := s.executeAutonomousTask(ctx, request, status); err != nil {
		s.log.Error(fmt.Sprintf("Autonomous task execution failed: %v", err))
		s.updateAutonomousTaskStatus(status.TaskId, "failed", err.Error(), 100.0)
		s.notifyTaskFailure(status.TaskId, err.Error())
	}
}

// executeAutonomousTask executes autonomous task with self-correction loop
func (s *Service) executeAutonomousTask(ctx context.Context, request domain.AutonomousTaskRequest, status *domain.AutonomousTaskStatus) error {
	basePipeline, planningTask, err := s.planAutonomousTask(ctx, request, status)
	if err != nil {
		return err
	}

	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		s.log.Info(fmt.Sprintf("[Task %s] Starting pipeline execution, attempt %d/%d.", status.TaskId, i+1, maxRetries))
		currentPipeline := *basePipeline

		if err := s.planner.ExecutePipeline(ctx, &currentPipeline); err == nil && currentPipeline.Status == PipelineStatusCompleted {
			s.finishAutonomousTask(request, status)
			return nil
		}

		s.log.Error(fmt.Sprintf("[Task %s] Pipeline execution failed", status.TaskId))
		if err := s.attemptRepair(ctx, planningTask, &currentPipeline, status, i); err != nil {
			return err
		}
	}
	return fmt.Errorf("task failed after %d repair attempts", maxRetries)
}

// planAutonomousTask creates execution plan for autonomous task
func (s *Service) planAutonomousTask(ctx context.Context, request domain.AutonomousTaskRequest, status *domain.AutonomousTaskStatus) (*TaskPipeline, domain.Task, error) {
	s.updateAutonomousTaskStatus(status.TaskId, "running", "Planning task...", 10.0)
	s.log.Info(fmt.Sprintf("[Task %s] Generating execution plan for: %s", status.TaskId, request.Task))

	contextPack, err := s.buildContextForTask(ctx, request)
	if err != nil {
		s.log.Error(fmt.Sprintf("[Task %s] Failed to build context: %v", status.TaskId, err))
		return nil, domain.Task{}, fmt.Errorf("failed to build context for task: %w", err)
	}

	planningTask := domain.Task{
		ID: status.TaskId, Name: "Autonomous Planning Task",
		Metadata: map[string]interface{}{"original_request": request.Task, "sla_policy": request.SlaPolicy, "project_path": request.ProjectPath},
	}

	var policy *PipelinePolicy
	if llmResponse, err := s.routerLlmService.CreatePipelineWithLLM(ctx, planningTask, contextPack); err == nil && llmResponse != nil && !llmResponse.FallbackUsed {
		policy = llmResponse.Policy
		s.log.Info(fmt.Sprintf("[Task %s] Using LLM-defined policy.", status.TaskId))
	} else {
		s.log.Info(fmt.Sprintf("[Task %s] Using heuristic policy.", status.TaskId))
	}

	basePipeline, err := s.planner.CreatePipeline(ctx, planningTask, policy)
	if err != nil {
		s.log.Error(fmt.Sprintf("[Task %s] Failed to create execution plan: %v", status.TaskId, err))
		return nil, domain.Task{}, fmt.Errorf("failed to create execution plan: %w", err)
	}

	s.log.Info(fmt.Sprintf("[Task %s] Execution plan generated with %d steps.", status.TaskId, len(basePipeline.Steps)))
	s.updateAutonomousTaskStatus(status.TaskId, "running", "Execution plan created. Starting execution...", 20.0)
	return basePipeline, planningTask, nil
}

// finishAutonomousTask completes the task and generates report
func (s *Service) finishAutonomousTask(request domain.AutonomousTaskRequest, status *domain.AutonomousTaskStatus) {
	s.updateAutonomousTaskStatus(status.TaskId, "running", "Generating final report...", 95.0)
	if diff, err := s.gitRepo.GenerateDiff(request.ProjectPath); err != nil {
		s.log.Error(fmt.Sprintf("[Task %s] Failed to generate git diff: %v", status.TaskId, err))
	} else {
		s.log.Info(fmt.Sprintf("[Task %s] Git Diff:\n%s", status.TaskId, diff))
	}
	s.updateAutonomousTaskStatus(status.TaskId, "completed", "Task completed successfully", 100.0)
	s.log.Info(fmt.Sprintf("[Task %s] Autonomous task finished.", status.TaskId))
}

// attemptRepair attempts to repair a failed pipeline step
func (s *Service) attemptRepair(ctx context.Context, planningTask domain.Task, pipeline *TaskPipeline, status *domain.AutonomousTaskStatus, attempt int) error {
	s.updateAutonomousTaskStatus(status.TaskId, "running", "Execution failed. Attempting self-correction...", 80.0+float64(attempt)*5)

	failedStep := s.findFailedStep(pipeline)
	if failedStep == nil {
		return fmt.Errorf("pipeline failed but no failed step found. Final status: %s", pipeline.Status)
	}

	s.log.Info(fmt.Sprintf("[Task %s] Found failed step: %s. Attempting to repair.", status.TaskId, failedStep.Name))
	repairPipeline, err := s.createRepairPipeline(ctx, planningTask, failedStep)
	if err != nil {
		return fmt.Errorf("failed to create repair pipeline: %w", err)
	}

	if err := s.planner.ExecutePipeline(ctx, repairPipeline); err != nil {
		return fmt.Errorf("repair pipeline execution failed: %w", err)
	}
	if repairPipeline.Status != PipelineStatusCompleted {
		return fmt.Errorf("repair pipeline did not complete successfully. Final status: %s", repairPipeline.Status)
	}
	s.log.Info(fmt.Sprintf("[Task %s] Repair successful. Retrying main pipeline.", status.TaskId))
	return nil
}

func (s *Service) findFailedStep(pipeline *TaskPipeline) *TaskPipelineStep {
	for _, step := range pipeline.Steps {
		if step.Status == StepStatusFailed {
			return step
		}
	}
	return nil
}

func (s *Service) createRepairPipeline(_ context.Context, task domain.Task, failedStep *TaskPipelineStep) (*TaskPipeline, error) {
	repairStep := &TaskPipelineStep{
		ID:        fmt.Sprintf("%s-repair-%d", task.ID, time.Now().Unix()),
		Name:      fmt.Sprintf("Repair for: %s", failedStep.Name),
		Type:      StepTypeRepair,
		Status:    StepStatusPending,
		Priority:  1,
		DependsOn: []string{},
		Config: map[string]interface{}{
			"task_id":           task.ID,
			"project_path":      task.Metadata["project_path"],
			"error_output":      failedStep.Error,
			"repair_strategies": []string{"auto_fix"},
			"max_attempts":      1,
		},
	}

	return &TaskPipeline{
		TaskID:    fmt.Sprintf("%s-repair", task.ID),
		Steps:     []*TaskPipelineStep{repairStep},
		Status:    PipelineStatusPending,
		CreatedAt: time.Now(),
		Policy:    &PipelinePolicy{FailFast: true},
	}, nil
}

func (s *Service) changeTaskState(taskID string, fromState, toState domain.TaskState, message, action string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, exists := s.statuses[taskID]
	if !exists {
		return domain.NewTaskNotFoundError(taskID)
	}

	if status.State != fromState {
		return domain.NewInvalidTaskStateError(taskID, string(status.State), string(fromState))
	}

	status.State = toState
	status.Message = message
	status.UpdatedAt = time.Now()

	if err := s.saveStatuses(); err != nil {
		return domain.NewInternalError("Failed to save task status after "+action, err)
	}

	s.log.Info(fmt.Sprintf("Task %s %s successfully", taskID, action))
	return nil
}

func (s *Service) calculateEstimatedTimeRemaining(status *domain.TaskStatus) int64 {
	if status.State == domain.TaskStateDone || status.State == domain.TaskStateFailed {
		return 0
	}
	if status.StartedAt == nil {
		return DefaultEstimatedTimeSeconds
	}
	if status.Progress < 0.01 {
		return DefaultEstimatedTimeSeconds
	}
	if status.Progress >= 1.0 {
		return 0
	}

	elapsed := time.Since(*status.StartedAt)
	totalEstimated := elapsed.Seconds() / status.Progress
	remaining := totalEstimated * (1.0 - status.Progress)

	if remaining > MaxEstimatedTimeSeconds {
		remaining = MaxEstimatedTimeSeconds
	}

	return int64(remaining)
}

func (s *Service) validateAutonomousTaskRequest(request domain.AutonomousTaskRequest) error {
	if request.Task == "" {
		return fmt.Errorf("task description cannot be empty")
	}
	if request.ProjectPath == "" {
		return fmt.Errorf("project path cannot be empty")
	}
	if request.SlaPolicy == "" {
		return fmt.Errorf("SLA policy cannot be empty")
	}

	validSLAPolicies := []string{"lite", "standard", "strict"}
	for _, policy := range validSLAPolicies {
		if request.SlaPolicy == policy {
			return nil
		}
	}

	return fmt.Errorf("invalid SLA policy: %s, must be one of: %v", request.SlaPolicy, validSLAPolicies)
}

func (s *Service) hasRunningTasks() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, status := range s.statuses {
		if status.State == domain.TaskStateTodo {
			return true
		}
	}
	return false
}

func (s *Service) createTaskStatus(taskID string, _ domain.AutonomousTaskRequest) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.statuses[taskID] = &domain.TaskStatus{
		TaskID: taskID,
		State:  domain.TaskStateTodo,
	}

	return s.saveStatuses()
}

func (s *Service) notifyTaskFailure(taskID string, errorMsg string) {
	s.log.Error(fmt.Sprintf("Task %s failed: %s", taskID, errorMsg))
}

func (s *Service) updateAutonomousTaskStatus(taskID, status, message string, progress float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	taskStatus, exists := s.statuses[taskID]
	if !exists {
		taskStatus = &domain.TaskStatus{TaskID: taskID}
		s.statuses[taskID] = taskStatus
	}

	taskStatus.Message = message
	taskStatus.Progress = progress / 100.0

	switch status {
	case "running":
		taskStatus.State = domain.TaskStateRunning
	case "completed":
		taskStatus.State = domain.TaskStateDone
	case "failed":
		taskStatus.State = domain.TaskStateFailed
	}
}

func (s *Service) buildContextForTask(_ context.Context, request domain.AutonomousTaskRequest) (map[string]interface{}, error) {
	return map[string]interface{}{
		"task":         request.Task,
		"project_path": request.ProjectPath,
		"sla_policy":   request.SlaPolicy,
	}, nil
}
