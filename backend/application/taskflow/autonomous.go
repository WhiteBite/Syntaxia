package taskflow

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
	"time"
)

// StartAutonomousTask starts an autonomous task
func (s *Service) StartAutonomousTask(ctx context.Context, request domain.AutonomousTaskRequest) (*domain.AutonomousTaskResponse, error) {
	s.log.Info(fmt.Sprintf("Starting autonomous task: %s", request.Task))

	if err := s.validateAutonomousTaskRequest(request); err != nil {
		return nil, domain.NewValidationError("Invalid autonomous task request", map[string]interface{}{
			"task":        request.Task,
			"projectPath": request.ProjectPath,
			"slaPolicy":   request.SlaPolicy,
		})
	}

	taskID := fmt.Sprintf("autonomous_%d", time.Now().Unix())

	if s.hasRunningTasks() {
		return nil, domain.NewInvalidTaskStateError(taskID, "new", "no_running_tasks")
	}

	status := &domain.AutonomousTaskStatus{
		TaskId:                 taskID,
		Status:                 "pending",
		CurrentStep:            "initializing",
		Progress:               0.0,
		EstimatedTimeRemaining: DefaultEstimatedTimeSeconds,
		StartedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	if err := s.createTaskStatus(taskID, request); err != nil {
		return nil, domain.NewInternalError("Failed to create task status", err)
	}

	go s.safeExecuteAutonomousTask(ctx, request, status)

	return &domain.AutonomousTaskResponse{
		TaskId:  taskID,
		Status:  "accepted",
		Message: "Task accepted for processing",
	}, nil
}

// CancelAutonomousTask cancels an autonomous task
func (s *Service) CancelAutonomousTask(ctx context.Context, taskID string) error {
	s.log.Info(fmt.Sprintf("Cancelling autonomous task: %s", taskID))

	s.mu.Lock()
	defer s.mu.Unlock()

	status, exists := s.statuses[taskID]
	if !exists {
		return domain.NewTaskNotFoundError(taskID)
	}

	if status.State == domain.TaskStateDone {
		return domain.NewInvalidTaskStateError(taskID, string(status.State), "cancellable")
	}

	status.State = domain.TaskStateFailed
	status.Message = "Task cancelled by user"

	if err := s.saveStatuses(); err != nil {
		return domain.NewInternalError("Failed to save task status after cancellation", err)
	}

	s.log.Info(fmt.Sprintf("Task %s cancelled successfully", taskID))
	return nil
}

// GetAutonomousTaskStatus gets autonomous task status
func (s *Service) GetAutonomousTaskStatus(ctx context.Context, taskID string) (*domain.AutonomousTaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.statuses[taskID]
	if !exists {
		return nil, domain.NewTaskNotFoundError(taskID)
	}

	autonomousStatus := &domain.AutonomousTaskStatus{
		TaskId:                 taskID,
		Status:                 string(status.State),
		CurrentStep:            status.Message,
		Progress:               status.Progress * 100,
		EstimatedTimeRemaining: s.calculateEstimatedTimeRemaining(status),
		StartedAt:              time.Now(),
		UpdatedAt:              status.UpdatedAt,
		Error:                  status.Error,
	}

	if status.StartedAt != nil {
		autonomousStatus.StartedAt = *status.StartedAt
	}

	return autonomousStatus, nil
}

// ListAutonomousTasks returns list of autonomous tasks
func (s *Service) ListAutonomousTasks(ctx context.Context, projectPath string) ([]domain.AutonomousTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var autonomousTasks []domain.AutonomousTask

	for taskID, status := range s.statuses {
		if strings.HasPrefix(taskID, "autonomous_") {
			task := domain.AutonomousTask{
				ID:          taskID,
				Name:        "Autonomous Task",
				Description: status.Message,
				Status:      string(status.State),
				ProjectPath: projectPath,
				Progress:    status.Progress,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Error:       status.Error,
			}

			if status.StartedAt != nil {
				task.CreatedAt = *status.StartedAt
			}
			if status.CompletedAt != nil {
				task.CompletedAt = status.CompletedAt
			}

			autonomousTasks = append(autonomousTasks, task)
		}
	}

	s.log.Info(fmt.Sprintf("Found %d autonomous tasks for project %s", len(autonomousTasks), projectPath))
	return autonomousTasks, nil
}

// GetTaskLogs returns logs for a task
func (s *Service) GetTaskLogs(ctx context.Context, taskID string) ([]domain.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.statuses[taskID]
	if !exists {
		return nil, domain.NewTaskNotFoundError(taskID)
	}

	var logs []domain.LogEntry

	if status.StartedAt != nil {
		logs = append(logs, domain.LogEntry{
			ID:        fmt.Sprintf("%s-created", taskID),
			TaskID:    taskID,
			Level:     "INFO",
			Message:   "Task created",
			Timestamp: *status.StartedAt,
			Metadata:  map[string]interface{}{"event": "task_created"},
		})
	}

	logs = append(logs, domain.LogEntry{
		ID:        fmt.Sprintf("%s-status", taskID),
		TaskID:    taskID,
		Level:     "INFO",
		Message:   fmt.Sprintf("Task status: %s - %s", status.State, status.Message),
		Timestamp: status.UpdatedAt,
		Metadata:  map[string]interface{}{"state": string(status.State), "progress": status.Progress},
	})

	if status.Error != "" {
		logs = append(logs, domain.LogEntry{
			ID:        fmt.Sprintf("%s-error", taskID),
			TaskID:    taskID,
			Level:     "ERROR",
			Message:   status.Error,
			Timestamp: status.UpdatedAt,
			Metadata:  map[string]interface{}{"event": "task_error"},
		})
	}

	if status.CompletedAt != nil {
		logs = append(logs, domain.LogEntry{
			ID:        fmt.Sprintf("%s-completed", taskID),
			TaskID:    taskID,
			Level:     "INFO",
			Message:   "Task completed",
			Timestamp: *status.CompletedAt,
			Metadata:  map[string]interface{}{"event": "task_completed", "duration": status.Duration.String()},
		})
	}

	s.log.Debug(fmt.Sprintf("Retrieved %d log entries for task %s", len(logs), taskID))
	return logs, nil
}

// PauseTask pauses task execution
func (s *Service) PauseTask(ctx context.Context, taskID string) error {
	return s.changeTaskState(taskID, domain.TaskStateTodo, domain.TaskStateBlocked, "Task paused by user", "paused")
}

// ResumeTask resumes paused task execution
func (s *Service) ResumeTask(ctx context.Context, taskID string) error {
	return s.changeTaskState(taskID, domain.TaskStateBlocked, domain.TaskStateTodo, "Task resumed by user", "resumed")
}
