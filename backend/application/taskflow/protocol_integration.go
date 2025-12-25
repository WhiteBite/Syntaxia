package taskflow

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
	"time"
)

// ProtocolIntegration integrates Task Protocol with Taskflow services
type ProtocolIntegration struct {
	log                 domain.Logger
	taskflowService     domain.TaskflowService
	taskProtocolService domain.TaskProtocolService
	configService       ProtocolConfigService
	aiService           IntelligentAI
}

// ProtocolConfigService interface for protocol configuration
type ProtocolConfigService interface {
	GetConfigurationForProject(ctx context.Context, projectPath string, languages []string) (*domain.TaskProtocolConfig, error)
}

// IntelligentAI interface for AI operations
type IntelligentAI interface{}

// NewProtocolIntegration creates a new integration service
func NewProtocolIntegration(
	log domain.Logger,
	taskflowService domain.TaskflowService,
	taskProtocolService domain.TaskProtocolService,
	configService ProtocolConfigService,
	aiService IntelligentAI,
) *ProtocolIntegration {
	return &ProtocolIntegration{
		log:                 log,
		taskflowService:     taskflowService,
		taskProtocolService: taskProtocolService,
		configService:       configService,
		aiService:           aiService,
	}
}

// ExecuteTaskWithProtocol executes a taskflow task with protocol verification
func (t *ProtocolIntegration) ExecuteTaskWithProtocol(ctx context.Context, taskID string) (*domain.TaskProtocolResult, error) {
	t.log.Info(fmt.Sprintf("Executing task %s with protocol verification", taskID))

	task, err := t.getTaskDetails(taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get task details: %w", err)
	}

	if err := t.taskflowService.UpdateTaskStatus(taskID, domain.TaskStateRunning, "Starting protocol verification"); err != nil {
		t.log.Warning(fmt.Sprintf("Failed to update task status: %v", err))
	}

	protocolConfig, err := t.createProtocolConfigForTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol config: %w", err)
	}

	result, err := t.taskProtocolService.ExecuteProtocol(ctx, protocolConfig)
	if err != nil {
		_ = t.taskflowService.UpdateTaskStatus(taskID, domain.TaskStateFailed, fmt.Sprintf("Protocol verification failed: %v", err))
		return result, fmt.Errorf("protocol execution failed: %w", err)
	}

	if result.Success {
		if err := t.taskflowService.UpdateTaskStatus(taskID, domain.TaskStateDone, "Protocol verification completed successfully"); err != nil {
			t.log.Warning(fmt.Sprintf("Failed to update task status: %v", err))
		}
	} else {
		if err := t.taskflowService.UpdateTaskStatus(taskID, domain.TaskStateFailed, fmt.Sprintf("Protocol verification failed: %s", result.FinalError)); err != nil {
			t.log.Warning(fmt.Sprintf("Failed to update task status: %v", err))
		}
	}

	t.storeProtocolResultInTask(taskID, result)

	t.log.Info(fmt.Sprintf("Task %s protocol execution completed with success: %t", taskID, result.Success))
	return result, nil
}

// ExecuteTaskflowWithProtocol executes an entire taskflow with protocol verification
func (t *ProtocolIntegration) ExecuteTaskflowWithProtocol(ctx context.Context, options *ProtocolOptions) (*ProtocolResult, error) {
	t.log.Info("Executing taskflow with protocol verification")

	result := &ProtocolResult{
		StartedAt:   time.Now(),
		TaskResults: make(map[string]*domain.TaskProtocolResult),
		Success:     true,
	}

	tasks, err := t.taskflowService.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	for _, task := range tasks {
		if t.shouldExecuteTask(task, options) {
			t.log.Info(fmt.Sprintf("Executing task: %s", task.ID))

			taskResult, err := t.ExecuteTaskWithProtocol(ctx, task.ID)
			if err != nil {
				t.log.Error(fmt.Sprintf("Task %s failed: %v", task.ID, err))
				result.Success = false
				result.FailedTasks = append(result.FailedTasks, task.ID)

				if options.FailFast {
					result.CompletedAt = time.Now()
					return result, fmt.Errorf("taskflow execution failed at task %s: %w", task.ID, err)
				}
			} else {
				result.TaskResults[task.ID] = taskResult
				if !taskResult.Success {
					result.Success = false
					result.FailedTasks = append(result.FailedTasks, task.ID)
				}
			}
		}
	}

	result.CompletedAt = time.Now()
	t.log.Info(fmt.Sprintf("Taskflow execution completed with success: %t", result.Success))
	return result, nil
}

// CreateTaskProtocolForAIGeneration creates a protocol for AI-generated code validation
func (t *ProtocolIntegration) CreateTaskProtocolForAIGeneration(ctx context.Context, aiRequest *AICodeGenerationRequest) (*domain.TaskProtocolConfig, error) {
	t.log.Info("Creating task protocol for AI code generation")

	languages, err := t.detectLanguagesFromContext(aiRequest.Context)
	if err != nil {
		return nil, fmt.Errorf("failed to detect languages: %w", err)
	}

	config, err := t.configService.GetConfigurationForProject(ctx, aiRequest.ProjectPath, languages)
	if err != nil {
		return nil, fmt.Errorf("failed to get protocol configuration: %w", err)
	}

	config = t.customizeConfigForAIGeneration(config, aiRequest)

	return config, nil
}

// ValidateAIGeneratedCode validates AI-generated code using the task protocol
func (t *ProtocolIntegration) ValidateAIGeneratedCode(ctx context.Context, request *AICodeValidationRequest) (*domain.TaskProtocolResult, error) {
	t.log.Info("Validating AI-generated code with task protocol")

	protocolConfig, err := t.CreateTaskProtocolForAIGeneration(ctx, &AICodeGenerationRequest{
		ProjectPath: request.ProjectPath,
		Context:     request.Context,
		Languages:   request.Languages,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create protocol config: %w", err)
	}

	result, err := t.taskProtocolService.ExecuteProtocol(ctx, protocolConfig)
	if err != nil {
		return result, fmt.Errorf("validation protocol failed: %w", err)
	}

	if !result.Success && protocolConfig.SelfCorrection.Enabled && protocolConfig.SelfCorrection.AIAssistance {
		t.log.Info("Attempting AI-assisted correction")

		correctionResult, err := t.attemptAICorrection(ctx, result, request)
		if err != nil {
			t.log.Warning(fmt.Sprintf("AI correction failed: %v", err))
		} else if correctionResult != nil {
			result = correctionResult
		}
	}

	return result, nil
}

func (t *ProtocolIntegration) getTaskDetails(taskID string) (*domain.Task, error) {
	return &domain.Task{
		ID:        taskID,
		Name:      fmt.Sprintf("Task %s", taskID),
		State:     domain.TaskStateTodo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Metadata:  make(map[string]interface{}),
	}, nil
}

func (t *ProtocolIntegration) createProtocolConfigForTask(ctx context.Context, task *domain.Task) (*domain.TaskProtocolConfig, error) {
	projectPath := "/default/project/path"
	if path, exists := task.Metadata["projectPath"]; exists {
		if pathStr, ok := path.(string); ok {
			projectPath = pathStr
		}
	}

	languages := []string{"go"}
	if langs, exists := task.Metadata["languages"]; exists {
		if langSlice, ok := langs.([]string); ok {
			languages = langSlice
		}
	}

	config, err := t.configService.GetConfigurationForProject(ctx, projectPath, languages)
	if err != nil {
		return nil, err
	}

	if task.Budgets.MaxFiles > 0 {
		config.MaxRetries = 5
	}

	return config, nil
}

func (t *ProtocolIntegration) shouldExecuteTask(task domain.Task, options *ProtocolOptions) bool {
	if options.TaskFilter != "" && task.ID != options.TaskFilter {
		return false
	}

	if task.State == domain.TaskStateDone && !options.ForceRerun {
		return false
	}

	return true
}

func (t *ProtocolIntegration) storeProtocolResultInTask(taskID string, result *domain.TaskProtocolResult) {
	metadata := map[string]interface{}{
		"protocolResult": map[string]interface{}{
			"success":          result.Success,
			"correctionCycles": result.CorrectionCycles,
			"stagesExecuted":   len(result.Stages),
			"executionTime":    result.CompletedAt.Sub(result.StartedAt).String(),
		},
	}

	t.log.Debug(fmt.Sprintf("Storing protocol result metadata for task %s: %+v", taskID, metadata))
}

func (t *ProtocolIntegration) detectLanguagesFromContext(context string) ([]string, error) {
	languages := make([]string, 0)

	if containsGoKeywords(context) {
		languages = append(languages, "go")
	}
	if containsTypeScriptKeywords(context) {
		languages = append(languages, "typescript")
	}
	if containsJavaScriptKeywords(context) {
		languages = append(languages, "javascript")
	}

	if len(languages) == 0 {
		languages = append(languages, "go")
	}

	return languages, nil
}

func (t *ProtocolIntegration) customizeConfigForAIGeneration(config *domain.TaskProtocolConfig, request *AICodeGenerationRequest) *domain.TaskProtocolConfig {
	config.SelfCorrection.Enabled = true
	config.SelfCorrection.AIAssistance = true
	config.SelfCorrection.MaxAttempts = 5
	config.MaxRetries = 3

	config.EnabledStages = []domain.ProtocolStage{
		domain.StageLinting,
		domain.StageBuilding,
		domain.StageTesting,
		domain.StageGuardrails,
	}

	return config
}

func (t *ProtocolIntegration) attemptAICorrection(ctx context.Context, result *domain.TaskProtocolResult, request *AICodeValidationRequest) (*domain.TaskProtocolResult, error) {
	var failedStage *domain.ProtocolStageResult
	for _, stage := range result.Stages {
		if !stage.Success {
			failedStage = stage
			break
		}
	}

	if failedStage == nil || failedStage.ErrorDetails == nil {
		return nil, fmt.Errorf("no failed stage with error details found")
	}

	taskContext := &domain.TaskContext{
		ProjectPath: request.ProjectPath,
		Languages:   request.Languages,
		Files:       request.ChangedFiles,
	}

	guidance, err := t.taskProtocolService.RequestCorrectionGuidance(ctx, failedStage.ErrorDetails, taskContext)
	if err != nil {
		return nil, fmt.Errorf("failed to get AI correction guidance: %w", err)
	}

	t.log.Info(fmt.Sprintf("Received AI correction guidance: %s", guidance.Explanation))

	return result, nil
}

func containsGoKeywords(context string) bool {
	goKeywords := []string{"package", "func", "import", "go.mod", "go.sum", ".go"}
	for _, keyword := range goKeywords {
		if strings.Contains(context, keyword) {
			return true
		}
	}
	return false
}

func containsTypeScriptKeywords(context string) bool {
	tsKeywords := []string{"interface", "type", "typescript", ".ts", ".tsx", "tsconfig.json"}
	for _, keyword := range tsKeywords {
		if strings.Contains(context, keyword) {
			return true
		}
	}
	return false
}

func containsJavaScriptKeywords(context string) bool {
	jsKeywords := []string{"javascript", ".js", ".jsx", "package.json", "node_modules"}
	for _, keyword := range jsKeywords {
		if strings.Contains(context, keyword) {
			return true
		}
	}
	return false
}

// ProtocolOptions defines taskflow protocol execution options
type ProtocolOptions struct {
	TaskFilter string
	FailFast   bool
	ForceRerun bool
	Parallel   bool
}

// ProtocolResult represents taskflow protocol execution result
type ProtocolResult struct {
	StartedAt   time.Time                             `json:"startedAt"`
	CompletedAt time.Time                             `json:"completedAt"`
	Success     bool                                  `json:"success"`
	TaskResults map[string]*domain.TaskProtocolResult `json:"taskResults"`
	FailedTasks []string                              `json:"failedTasks"`
}

// AICodeGenerationRequest represents AI code generation request
type AICodeGenerationRequest struct {
	ProjectPath string   `json:"projectPath"`
	Context     string   `json:"context"`
	Languages   []string `json:"languages"`
	Task        string   `json:"task"`
}

// AICodeValidationRequest represents AI code validation request
type AICodeValidationRequest struct {
	ProjectPath  string   `json:"projectPath"`
	Context      string   `json:"context"`
	Languages    []string `json:"languages"`
	ChangedFiles []string `json:"changedFiles"`
}
