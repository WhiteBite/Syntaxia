package router

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"sort"
	"strings"
	"time"
)

// PlannerService предоставляет эвристический планировщик TPL (DAG) без LLM
type PlannerService struct {
	log            domain.Logger
	buildService   domain.IBuildService
	testService    domain.ITestService
	staticAnalyzer domain.IStaticAnalyzerService
	repairService  domain.RepairService
}

// NewPlannerService создает новый сервис планировщика
func NewPlannerService(
	log domain.Logger,
	buildService domain.IBuildService,
	testService domain.ITestService,
	staticAnalyzer domain.IStaticAnalyzerService,
	repairService domain.RepairService,
) *PlannerService {
	return &PlannerService{
		log:            log,
		buildService:   buildService,
		testService:    testService,
		staticAnalyzer: staticAnalyzer,
		repairService:  repairService,
	}
}

// TaskPipelineStepType определяет тип шага пайплайна
type TaskPipelineStepType string

const (
	StepTypeRetrieve TaskPipelineStepType = "retrieve"
	StepTypeASTSynth TaskPipelineStepType = "ast_synth"
	StepTypeCompile  TaskPipelineStepType = "compile"
	StepTypeTest     TaskPipelineStepType = "test"
	StepTypeStatic   TaskPipelineStepType = "static"
	StepTypeRepair   TaskPipelineStepType = "repair"
	StepTypeFormat   TaskPipelineStepType = "format"
	StepTypeValidate TaskPipelineStepType = "validate"
)

// TaskPipelineStepStatus определяет статус шага пайплайна
type TaskPipelineStepStatus string

const (
	StepStatusPending   TaskPipelineStepStatus = "pending"
	StepStatusRunning   TaskPipelineStepStatus = "running"
	StepStatusCompleted TaskPipelineStepStatus = "completed"
	StepStatusFailed    TaskPipelineStepStatus = "failed"
)

// TaskPipelineStepResult содержит результат выполнения шага
type TaskPipelineStepResult struct {
	Success     bool               `json:"success"`
	Message     string             `json:"message"`
	Data        map[string]any     `json:"data,omitempty"`
	Metrics     map[string]float64 `json:"metrics,omitempty"`
	Artifacts   []string           `json:"artifacts,omitempty"`
	Warnings    []string           `json:"warnings,omitempty"`
	Suggestions []string           `json:"suggestions,omitempty"`
}

// TaskPipelineStep представляет шаг в пайплайне задачи
type TaskPipelineStep struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	Type        TaskPipelineStepType    `json:"type"`
	Status      TaskPipelineStepStatus  `json:"status"`
	Priority    int                     `json:"priority"`
	DependsOn   []string                `json:"depends_on"`
	Config      map[string]any          `json:"config"`
	StartedAt   *time.Time              `json:"started_at,omitempty"`
	CompletedAt *time.Time              `json:"completed_at,omitempty"`
	Duration    time.Duration           `json:"duration,omitempty"`
	Error       string                  `json:"error,omitempty"`
	Result      *TaskPipelineStepResult `json:"result,omitempty"`
}

// TaskPipelineStatus определяет статус пайплайна
type TaskPipelineStatus string

const (
	PipelineStatusPending   TaskPipelineStatus = "pending"
	PipelineStatusRunning   TaskPipelineStatus = "running"
	PipelineStatusCompleted TaskPipelineStatus = "completed"
	PipelineStatusFailed    TaskPipelineStatus = "failed"
)

// PipelinePolicy определяет политику выполнения пайплайна
type PipelinePolicy struct {
	EnableRetrieve bool          `json:"enable_retrieve"`
	EnableASTSynth bool          `json:"enable_ast_synth"`
	EnableCompile  bool          `json:"enable_compile"`
	EnableTest     bool          `json:"enable_test"`
	EnableStatic   bool          `json:"enable_static"`
	EnableRepair   bool          `json:"enable_repair"`
	EnableFormat   bool          `json:"enable_format"`
	EnableValidate bool          `json:"enable_validate"`
	FailFast       bool          `json:"fail_fast"`
	RetryFailed    bool          `json:"retry_failed"`
	MaxRetries     int           `json:"max_retries"`
	ParallelSteps  bool          `json:"parallel_steps"`
	Timeout        time.Duration `json:"timeout"`
}

// TaskPipeline представляет полный пайплайн для задачи
type TaskPipeline struct {
	TaskID      string              `json:"task_id"`
	Steps       []*TaskPipelineStep `json:"steps"`
	Status      TaskPipelineStatus  `json:"status"`
	CreatedAt   time.Time           `json:"created_at"`
	StartedAt   *time.Time          `json:"started_at,omitempty"`
	CompletedAt *time.Time          `json:"completed_at,omitempty"`
	Duration    time.Duration       `json:"duration,omitempty"`
	Error       string              `json:"error,omitempty"`
	Policy      *PipelinePolicy     `json:"policy"`
}

// CreatePipeline создает пайплайн для задачи
func (r *PlannerService) CreatePipeline(_ context.Context, task domain.Task, policy *PipelinePolicy) (*TaskPipeline, error) {
	r.log.Info(fmt.Sprintf("Creating pipeline for task: %s", task.ID))

	if policy == nil {
		policy = r.determinePolicy(task)
	}

	steps := r.createPipelineSteps(task, policy)
	r.sortPipelineSteps(steps)

	pipeline := &TaskPipeline{
		TaskID:    task.ID,
		Steps:     steps,
		Status:    PipelineStatusPending,
		CreatedAt: time.Now(),
		Policy:    policy,
	}

	r.log.Info(fmt.Sprintf("Created pipeline with %d steps for task: %s", len(steps), task.ID))
	return pipeline, nil
}

// determinePolicy определяет политику выполнения на основе задачи
func (r *PlannerService) determinePolicy(task domain.Task) *PipelinePolicy {
	policy := &PipelinePolicy{
		EnableRetrieve: true,
		EnableASTSynth: true,
		EnableCompile:  true,
		EnableTest:     true,
		EnableStatic:   true,
		EnableRepair:   true,
		EnableFormat:   true,
		EnableValidate: true,
		FailFast:       true,
		RetryFailed:    true,
		MaxRetries:     3,
		ParallelSteps:  false,
		Timeout:        30 * time.Minute,
	}

	// Настройка политики на основе типа задачи
	r.applyTaskSpecificPolicy(task, policy)

	return policy
}

// applyTaskSpecificPolicy применяет специфичные настройки политики
func (r *PlannerService) applyTaskSpecificPolicy(task domain.Task, policy *PipelinePolicy) {
	switch {
	case strings.Contains(task.ID, "ark-100"),
		strings.Contains(task.ID, "ark-110"),
		strings.Contains(task.ID, "ark-120"),
		strings.Contains(task.ID, "ark-140"):
		policy.EnableCompile = false
		policy.EnableTest = false
		policy.EnableStatic = false
		policy.EnableRepair = false
	case strings.Contains(task.ID, "ark-130"),
		strings.Contains(task.ID, "ark-150"):
		policy.EnableRepair = true
	case strings.Contains(task.ID, "ark-160"):
		policy.EnableASTSynth = false
		policy.EnableCompile = false
		policy.EnableStatic = false
		policy.EnableRepair = false
	case strings.Contains(task.ID, "ark-170"):
		policy.EnableASTSynth = false
		policy.EnableCompile = false
		policy.EnableTest = false
		policy.EnableRepair = false
	case strings.Contains(task.ID, "ark-180"):
		policy.EnableCompile = false
		policy.EnableTest = false
		policy.EnableStatic = false
		policy.EnableRepair = false
		policy.EnableFormat = false
	case strings.Contains(task.ID, "ark-999"):
		policy.EnableRepair = true
		policy.FailFast = true
	}
}

// GetPipelineStatus возвращает статус пайплайна
func (r *PlannerService) GetPipelineStatus(pipeline *TaskPipeline) map[string]any {
	completed, failed, pending, running := 0, 0, 0, 0

	for _, step := range pipeline.Steps {
		switch step.Status {
		case StepStatusCompleted:
			completed++
		case StepStatusFailed:
			failed++
		case StepStatusPending:
			pending++
		case StepStatusRunning:
			running++
		}
	}

	return map[string]any{
		"task_id":     pipeline.TaskID,
		"status":      pipeline.Status,
		"total_steps": len(pipeline.Steps),
		"completed":   completed,
		"failed":      failed,
		"pending":     pending,
		"running":     running,
		"progress":    float64(completed) / float64(len(pipeline.Steps)),
		"duration":    pipeline.Duration,
		"error":       pipeline.Error,
	}
}

// stepBuilder помогает создавать шаги пайплайна
type stepBuilder struct {
	taskID  string
	stepID  int
	steps   []*TaskPipelineStep
	prevDep string
}

func newStepBuilder(taskID string) *stepBuilder {
	return &stepBuilder{
		taskID: taskID,
		stepID: 1,
		steps:  make([]*TaskPipelineStep, 0),
	}
}

func (sb *stepBuilder) addStep(name string, stepType TaskPipelineStepType, priority int, config map[string]any) {
	dependsOn := []string{}
	if sb.prevDep != "" {
		dependsOn = []string{sb.prevDep}
	}

	step := &TaskPipelineStep{
		ID:        fmt.Sprintf("%s-step-%d", sb.taskID, sb.stepID),
		Name:      name,
		Type:      stepType,
		Status:    StepStatusPending,
		Priority:  priority,
		DependsOn: dependsOn,
		Config:    config,
	}

	sb.steps = append(sb.steps, step)
	sb.prevDep = step.ID
	sb.stepID++
}

// createPipelineSteps создает шаги пайплайна на основе политики
func (r *PlannerService) createPipelineSteps(task domain.Task, policy *PipelinePolicy) []*TaskPipelineStep {
	sb := newStepBuilder(task.ID)

	if policy.EnableRetrieve {
		sb.addStep("Retrieve Context", StepTypeRetrieve, 1, map[string]any{
			"task_id":      task.ID,
			"step_file":    task.StepFile,
			"budgets":      task.Budgets,
			"dependencies": task.DependsOn,
		})
	}

	if policy.EnableASTSynth {
		sb.addStep("AST Synthesis", StepTypeASTSynth, 2, map[string]any{
			"task_id":      task.ID,
			"symbol_graph": true,
			"context_pack": true,
		})
	}

	if policy.EnableCompile {
		sb.addStep("Compile", StepTypeCompile, 3, map[string]any{
			"task_id":      task.ID,
			"project_path": task.Metadata["project_path"],
			"language":     "go",
			"build_mode":   "debug",
		})
	}

	if policy.EnableTest {
		sb.addStep("Test", StepTypeTest, 4, map[string]any{
			"task_id":      task.ID,
			"project_path": task.Metadata["project_path"],
			"test_mode":    "targeted",
			"coverage":     true,
		})
	}

	if policy.EnableStatic {
		sb.addStep("Static Analysis", StepTypeStatic, 5, map[string]any{
			"task_id":       task.ID,
			"project_path":  task.Metadata["project_path"],
			"analyzers":     []string{"staticcheck", "go vet", "eslint"},
			"fail_on_error": false,
		})
	}

	if policy.EnableFormat {
		sb.addStep("Format", StepTypeFormat, 6, map[string]any{
			"task_id":    task.ID,
			"formatters": []string{"gofmt", "goimports", "prettier"},
		})
	}

	if policy.EnableValidate {
		sb.addStep("Validate", StepTypeValidate, 7, map[string]any{
			"task_id":          task.ID,
			"validation_rules": []string{"syntax", "semantics", "policy"},
		})
	}

	if policy.EnableRepair {
		sb.addStep("Repair", StepTypeRepair, 8, map[string]any{
			"task_id":           task.ID,
			"project_path":      task.Metadata["project_path"],
			"repair_strategies": []string{"auto_fix", "suggestions", "rollback"},
			"max_attempts":      policy.MaxRetries,
		})
	}

	return sb.steps
}

// sortPipelineSteps сортирует шаги по приоритету и зависимостям
func (r *PlannerService) sortPipelineSteps(steps []*TaskPipelineStep) {
	sort.Slice(steps, func(i, j int) bool {
		// Сначала по приоритету
		if steps[i].Priority != steps[j].Priority {
			return steps[i].Priority < steps[j].Priority
		}
		// Затем по количеству зависимостей
		return len(steps[i].DependsOn) < len(steps[j].DependsOn)
	})
}
