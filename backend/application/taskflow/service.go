// Package taskflow provides task workflow management services.
package taskflow

import (
	"context"
	"fmt"
	"os"
	"syntaxia/application/router"
	"syntaxia/domain"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// RouterPlanner interface to avoid circular imports
type RouterPlanner interface {
	CreatePipeline(ctx context.Context, task domain.Task, policy *router.PipelinePolicy) (*router.TaskPipeline, error)
	ExecutePipeline(ctx context.Context, pipeline *router.TaskPipeline) error
	GetPipelineStatus(pipeline *router.TaskPipeline) map[string]any
}

// RouterLLMService interface to avoid circular imports
type RouterLLMService interface {
	CreatePipelineWithLLM(ctx context.Context, task domain.Task, contextPack map[string]any) (*router.LLMPipelineResponse, error)
}

// LLMPipelineResponse is an alias for router.LLMPipelineResponse
type LLMPipelineResponse = router.LLMPipelineResponse

// Service implements TaskflowService and TaskTypeProvider
type Service struct {
	log              domain.Logger
	config           domain.TaskflowConfig
	tasks            map[string]domain.Task
	statuses         map[string]*domain.TaskStatus
	mu               sync.RWMutex
	planPath         string
	statusPath       string
	planner          RouterPlanner
	routerLlmService RouterLLMService
	guardrails       domain.GuardrailService
	repo             domain.TaskflowRepository
	gitRepo          domain.GitRepository
}

// NewService creates a new taskflow service
func NewService(log domain.Logger, planner RouterPlanner, routerLlmService RouterLLMService, guardrails domain.GuardrailService, repo domain.TaskflowRepository, gitRepo domain.GitRepository) domain.TaskflowService {
	service := &Service{
		log:              log,
		tasks:            make(map[string]domain.Task),
		statuses:         make(map[string]*domain.TaskStatus),
		planPath:         "tasks/plan.yaml",
		statusPath:       "tasks/status.json",
		planner:          planner,
		routerLlmService: routerLlmService,
		guardrails:       guardrails,
		repo:             repo,
		gitRepo:          gitRepo,
		config: domain.TaskflowConfig{
			AutoStart:     true,
			MaxConcurrent: 3,
			RetryAttempts: 3,
			RetryDelay:    5 * time.Second,
			Timeout:       30 * time.Minute,
			EnableLogging: true,
			EnableMetrics: true,
		},
	}

	if _, err := service.LoadTasks(); err != nil {
		log.Warning(fmt.Sprintf("Failed to load tasks: %v", err))
	}

	return service
}

// GetTaskType returns task type by ID (TaskTypeProvider implementation)
func (s *Service) GetTaskType(taskID string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, exists := s.tasks[taskID]
	if !exists {
		return "", fmt.Errorf("task not found: %s", taskID)
	}

	return domain.TaskTypeFromID(taskID).String(), nil
}

// LoadTasks loads tasks from plan.yaml
func (s *Service) LoadTasks() ([]domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.planPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read plan file: %w", err)
	}

	var plan struct {
		Version int `yaml:"version"`
		Tasks   []struct {
			ID        string   `yaml:"id"`
			Name      string   `yaml:"name"`
			DependsOn []string `yaml:"dependsOn"`
			StepFile  string   `yaml:"stepFile"`
			Budgets   struct {
				MaxFiles        int `yaml:"maxFiles"`
				MaxChangedLines int `yaml:"maxChangedLines"`
			} `yaml:"budgets"`
			Status string `yaml:"status"`
		} `yaml:"tasks"`
	}

	if err := yaml.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("failed to parse plan file: %w", err)
	}

	statuses, err := s.loadStatuses()
	if err != nil {
		s.log.Warning(fmt.Sprintf("Failed to load statuses: %v", err))
	}

	tasks := make([]domain.Task, 0, len(plan.Tasks))
	for _, taskData := range plan.Tasks {
		state := domain.TaskStateTodo
		if status, exists := statuses[taskData.ID]; exists {
			state = status
		}

		task := domain.Task{
			ID:        taskData.ID,
			Name:      taskData.Name,
			State:     state,
			DependsOn: taskData.DependsOn,
			StepFile:  taskData.StepFile,
			Budgets: domain.TaskBudgets{
				MaxFiles:        taskData.Budgets.MaxFiles,
				MaxChangedLines: taskData.Budgets.MaxChangedLines,
			},
			Status:    taskData.Status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Metadata:  make(map[string]interface{}),
		}

		s.tasks[task.ID] = task
		tasks = append(tasks, task)
	}

	s.log.Info(fmt.Sprintf("Loaded %d tasks from plan", len(tasks)))
	return tasks, nil
}

// GetTaskStatus returns task status
func (s *Service) GetTaskStatus(taskID string) (*domain.TaskStatus, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, exists := s.statuses[taskID]
	if !exists {
		return nil, fmt.Errorf("task status not found: %s", taskID)
	}

	return status, nil
}

// UpdateTaskStatus updates task status
func (s *Service) UpdateTaskStatus(taskID string, state domain.TaskState, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	status, exists := s.statuses[taskID]
	if !exists {
		status = &domain.TaskStatus{TaskID: taskID}
		s.statuses[taskID] = status
	}

	status.State = state
	status.Message = message

	switch state {
	case domain.TaskStateDone, domain.TaskStateFailed, domain.TaskStateBlocked:
		if status.CompletedAt == nil {
			now := time.Now()
			status.CompletedAt = &now
			if status.StartedAt != nil {
				status.Duration = now.Sub(*status.StartedAt)
			}
		}
	}

	if err := s.saveStatuses(); err != nil {
		return fmt.Errorf("failed to save statuses: %w", err)
	}

	s.log.Info(fmt.Sprintf("Updated task %s status to %s: %s", taskID, state, message))
	return nil
}

// ExecuteTask executes a task
func (s *Service) ExecuteTask(taskID string) error {
	s.mu.Lock()
	task, exists := s.tasks[taskID]
	if !exists {
		s.mu.Unlock()
		return fmt.Errorf("task not found: %s", taskID)
	}
	s.mu.Unlock()

	if err := s.checkDependencies(taskID); err != nil {
		return err
	}
	if err := s.validateWithGuardrails(taskID); err != nil {
		return err
	}

	now := time.Now()
	status := &domain.TaskStatus{
		TaskID:    taskID,
		State:     domain.TaskStateTodo,
		Progress:  0.0,
		Message:   "Starting task execution",
		StartedAt: &now,
	}

	s.mu.Lock()
	s.statuses[taskID] = status
	s.mu.Unlock()

	s.log.Info(fmt.Sprintf("Creating pipeline for task: %s", taskID))

	pipeline, err := s.planner.CreatePipeline(context.Background(), task, nil)
	if err != nil {
		return fmt.Errorf("failed to create pipeline: %w", err)
	}

	s.log.Info(fmt.Sprintf("Executing pipeline for task: %s", taskID))
	if err := s.planner.ExecutePipeline(context.Background(), pipeline); err != nil {
		return fmt.Errorf("failed to execute pipeline: %w", err)
	}

	pipelineStatus := s.planner.GetPipelineStatus(pipeline)
	progress := pipelineStatus["progress"].(float64)
	status.Progress = progress
	status.Message = fmt.Sprintf("Pipeline completed with %d steps", len(pipeline.Steps))

	if pipeline.Status == PipelineStatusCompleted {
		status.State = domain.TaskStateDone
		status.Message = "Task completed successfully via pipeline"
	} else {
		status.State = domain.TaskStateFailed
		status.Message = fmt.Sprintf("Pipeline failed: %s", pipeline.Error)
	}

	if s.guardrails != nil && (strings.Contains(taskID, "scaffold") || strings.Contains(taskID, "deps_fix")) {
		s.guardrails.DisableEphemeralMode()
	}

	return s.UpdateTaskStatus(taskID, status.State, status.Message)
}

// ExecuteTaskflow executes the entire taskflow
func (s *Service) ExecuteTaskflow() error {
	s.log.Info("Starting taskflow execution")

	for {
		readyTasks, err := s.GetReadyTasks()
		if err != nil {
			return fmt.Errorf("failed to get ready tasks: %w", err)
		}

		if len(readyTasks) == 0 {
			s.log.Info("No more tasks to execute")
			break
		}

		for _, task := range readyTasks {
			if err := s.ExecuteTask(task.ID); err != nil {
				s.log.Error(fmt.Sprintf("Failed to execute task %s: %v", task.ID, err))
				if err := s.UpdateTaskStatus(task.ID, domain.TaskStateFailed, err.Error()); err != nil {
					s.log.Error(fmt.Sprintf("Failed to update task status: %v", err))
				}
			}
		}
	}

	s.log.Info("Taskflow execution completed")
	return nil
}

// GetReadyTasks returns tasks ready for execution
func (s *Service) GetReadyTasks() ([]domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var readyTasks []domain.Task

	for _, task := range s.tasks {
		if task.State != domain.TaskStateTodo {
			continue
		}

		ready := true
		for _, depID := range task.DependsOn {
			dep, exists := s.tasks[depID]
			if !exists || dep.State != domain.TaskStateDone {
				ready = false
				break
			}
		}

		if ready {
			readyTasks = append(readyTasks, task)
		}
	}

	return readyTasks, nil
}

// GetTaskDependencies returns task dependencies
func (s *Service) GetTaskDependencies(taskID string) ([]domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found: %s", taskID)
	}

	var deps []domain.Task
	for _, depID := range task.DependsOn {
		if dep, exists := s.tasks[depID]; exists {
			deps = append(deps, dep)
		}
	}

	return deps, nil
}

// ValidateTaskflow validates taskflow correctness
func (s *Service) ValidateTaskflow() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for taskID := range s.tasks {
		if !visited[taskID] {
			if s.hasCycle(taskID, visited, recStack) {
				return fmt.Errorf("circular dependency detected")
			}
		}
	}

	for _, task := range s.tasks {
		if task.StepFile != "" {
			if _, err := os.Stat(task.StepFile); os.IsNotExist(err) {
				return fmt.Errorf("step file not found: %s", task.StepFile)
			}
		}
	}

	return nil
}

// GetTaskflowProgress returns execution progress
func (s *Service) GetTaskflowProgress() (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.tasks) == 0 {
		return 0.0, nil
	}

	completed := 0
	for _, task := range s.tasks {
		if task.State == domain.TaskStateDone {
			completed++
		}
	}

	return float64(completed) / float64(len(s.tasks)), nil
}

// ResetTaskflow resets taskflow
func (s *Service) ResetTaskflow() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for taskID := range s.tasks {
		s.statuses[taskID] = &domain.TaskStatus{
			TaskID: taskID,
			State:  domain.TaskStateTodo,
		}
	}

	if err := s.saveStatuses(); err != nil {
		return fmt.Errorf("failed to save statuses: %w", err)
	}

	s.log.Info("Taskflow reset completed")
	return nil
}

func (s *Service) loadStatuses() (map[string]domain.TaskState, error) {
	if s.repo == nil {
		return make(map[string]domain.TaskState), nil
	}
	return s.repo.LoadStatuses()
}

func (s *Service) saveStatuses() error {
	if s.repo == nil {
		return nil
	}

	statuses := make(map[string]domain.TaskState)
	for taskID, status := range s.statuses {
		statuses[taskID] = status.State
	}

	return s.repo.SaveStatuses(statuses)
}

func (s *Service) hasCycle(taskID string, visited, recStack map[string]bool) bool {
	visited[taskID] = true
	recStack[taskID] = true

	task, exists := s.tasks[taskID]
	if !exists {
		return false
	}

	for _, depID := range task.DependsOn {
		if !visited[depID] {
			if s.hasCycle(depID, visited, recStack) {
				return true
			}
		} else if recStack[depID] {
			return true
		}
	}

	recStack[taskID] = false
	return false
}

func (s *Service) checkDependencies(taskID string) error {
	deps, err := s.GetTaskDependencies(taskID)
	if err != nil {
		return fmt.Errorf("failed to get dependencies: %w", err)
	}
	for _, dep := range deps {
		if dep.State != domain.TaskStateDone {
			return fmt.Errorf("dependency %s is not completed (state: %s)", dep.ID, dep.State)
		}
	}
	return nil
}

func (s *Service) validateWithGuardrails(taskID string) error {
	if s.guardrails == nil {
		return nil
	}

	taskType := domain.TaskTypeFromID(taskID).String()
	if taskType == "scaffold" || taskType == "deps_fix" {
		if err := s.guardrails.EnableEphemeralMode(taskID, taskType, 5*time.Minute); err != nil {
			s.log.Warning(fmt.Sprintf("Failed to enable ephemeral mode for task %s: %v", taskID, err))
		}
	}

	validationResult, err := s.guardrails.ValidateTask(taskID, []string{}, 0)
	if err != nil {
		s.log.Error(fmt.Sprintf("Guardrail validation failed for task %s: %v", taskID, err))
		return fmt.Errorf("guardrail validation failed: %w", err)
	}
	if !validationResult.Valid {
		s.log.Error(fmt.Sprintf("Task %s failed guardrail validation: %s", taskID, validationResult.Error))
		return fmt.Errorf("task validation failed: %s", validationResult.Error)
	}
	return nil
}
