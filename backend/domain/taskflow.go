package domain

import (
	"context"
	"time"
)

// TaskState состояние задачи
type TaskState string

const (
	TaskStateTodo    TaskState = "todo"
	TaskStateRunning TaskState = "running"
	TaskStateDone    TaskState = "done"
	TaskStateBlocked TaskState = "blocked"
	TaskStateFailed  TaskState = "failed"
)

// Task представляет задачу в taskflow
type Task struct {
	ID          string
	Name        string
	Description string
	State       TaskState
	DependsOn   []string
	StepFile    string
	Budgets     TaskBudgets
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	StartedAt   *time.Time
	CompletedAt *time.Time
	Error       string
	Metadata    map[string]interface{}
}

// TaskBudgets бюджетные ограничения задачи
type TaskBudgets struct {
	MaxFiles        int `json:"maxFiles"`
	MaxChangedLines int `json:"maxChangedLines"`
}

// TaskStatus статус выполнения задачи
type TaskStatus struct {
	TaskID      string
	State       TaskState
	Progress    float64 // 0.0 - 1.0
	Message     string
	Error       string
	StartedAt   *time.Time
	CompletedAt *time.Time
	UpdatedAt   time.Time
	Duration    time.Duration
}

// TaskflowConfig конфигурация taskflow
type TaskflowConfig struct {
	AutoStart     bool
	MaxConcurrent int
	RetryAttempts int
	RetryDelay    time.Duration
	Timeout       time.Duration
	EnableLogging bool
	EnableMetrics bool
}

// TaskflowRepository интерфейс для работы с taskflow statuses
type TaskflowRepository interface {
	// LoadStatuses загружает статусы задач из хранилища
	LoadStatuses() (map[string]TaskState, error)

	// SaveStatuses сохраняет статусы задач в хранилище
	SaveStatuses(statuses map[string]TaskState) error
}

// TaskflowService интерфейс для сервиса taskflow
type TaskflowService interface {
	// LoadTasks загружает задачи из plan.yaml
	LoadTasks() ([]Task, error)

	// GetTaskStatus возвращает статус задачи
	GetTaskStatus(taskID string) (*TaskStatus, error)

	// UpdateTaskStatus обновляет статус задачи
	UpdateTaskStatus(taskID string, state TaskState, message string) error

	// ExecuteTask выполняет задачу
	ExecuteTask(taskID string) error

	// ExecuteTaskflow выполняет весь taskflow
	ExecuteTaskflow() error

	// GetReadyTasks возвращает готовые к выполнению задачи
	GetReadyTasks() ([]Task, error)

	// GetTaskDependencies возвращает зависимости задачи
	GetTaskDependencies(taskID string) ([]Task, error)

	// ValidateTaskflow проверяет корректность taskflow
	ValidateTaskflow() error

	// GetTaskflowProgress возвращает прогресс выполнения
	GetTaskflowProgress() (float64, error)

	// ResetTaskflow сбрасывает taskflow
	ResetTaskflow() error

	// StartAutonomousTask запускает автономную задачу
	StartAutonomousTask(ctx context.Context, request AutonomousTaskRequest) (*AutonomousTaskResponse, error)

	// CancelAutonomousTask отменяет автономную задачу
	CancelAutonomousTask(ctx context.Context, taskID string) error

	// GetAutonomousTaskStatus получает статус автономной задачи
	GetAutonomousTaskStatus(ctx context.Context, taskID string) (*AutonomousTaskStatus, error)

	// ListAutonomousTasks возвращает список автономных задач
	ListAutonomousTasks(ctx context.Context, projectPath string) ([]AutonomousTask, error)

	// GetTaskLogs возвращает логи задачи
	GetTaskLogs(ctx context.Context, taskID string) ([]LogEntry, error)

	// PauseTask приостанавливает задачу
	PauseTask(ctx context.Context, taskID string) error

	// ResumeTask возобновляет задачу
	ResumeTask(ctx context.Context, taskID string) error
}
