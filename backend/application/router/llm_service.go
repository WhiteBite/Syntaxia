package router

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"strings"
	"time"
)

// LLMService предоставляет интеграцию LLM с Router/Planner
type LLMService struct {
	log        domain.Logger
	llmClient  domain.LLMClient
	fileReader domain.FileReader
	enabled    bool
}

// LLMConfig конфигурация для Router LLM сервиса
type LLMConfig struct {
	Enabled             bool                   `json:"enabled"`
	LLMConfig           domain.LLMConfig       `json:"llm_config"`
	FallbackToHeuristic bool                   `json:"fallback_to_heuristic"`
	MaxRetries          int                    `json:"max_retries"`
	Timeout             time.Duration          `json:"timeout"`
	Options             map[string]interface{} `json:"options"`
}

// LLMPipelineRequest запрос на создание пайплайна через LLM
type LLMPipelineRequest struct {
	TaskID       string                 `json:"task_id"`
	TaskName     string                 `json:"task_name"`
	Description  string                 `json:"description"`
	StepFile     string                 `json:"step_file"`
	Budgets      domain.TaskBudgets     `json:"budgets"`
	Dependencies []string               `json:"dependencies"`
	Context      map[string]interface{} `json:"context"`
	Options      map[string]interface{} `json:"options"`
}

// LLMPipelineResponse ответ от LLM с пайплайном
type LLMPipelineResponse struct {
	Policy       *PipelinePolicy `json:"policy"`
	Confidence   float64         `json:"confidence"`
	Reasoning    string          `json:"reasoning"`
	FallbackUsed bool            `json:"fallback_used"`
	Error        string          `json:"error,omitempty"`
}

// NewLLMService создает новый сервис Router LLM
func NewLLMService(config LLMConfig, log domain.Logger) *LLMService {
	// This constructor is kept for backward compatibility but should not be used
	// The container should use NewLLMServiceWithClient instead
	return &LLMService{
		log:     log,
		enabled: config.Enabled,
	}
}

// NewLLMServiceWithClient создает новый сервис Router LLM с внедренным клиентом
func NewLLMServiceWithClient(config LLMConfig, log domain.Logger, llmClient domain.LLMClient, fileReader domain.FileReader) *LLMService {
	return &LLMService{
		log:        log,
		llmClient:  llmClient,
		fileReader: fileReader,
		enabled:    config.Enabled,
	}
}

// IsEnabled проверяет, включен ли LLM роутер
func (r *LLMService) IsEnabled() bool {
	return r.enabled && r.llmClient != nil
}

// CreatePipelineWithLLM создает пайплайн с помощью LLM
func (r *LLMService) CreatePipelineWithLLM(ctx context.Context, task domain.Task, contextPack map[string]interface{}) (*LLMPipelineResponse, error) {
	if !r.IsEnabled() {
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        "LLM router is disabled",
		}, nil
	}

	// Проверяем доступность LLM сервиса
	if err := r.llmClient.HealthCheck(ctx); err != nil {
		r.log.Warning(fmt.Sprintf("LLM server not available, falling back to heuristic: %v", err))
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        fmt.Sprintf("LLM server not available: %v", err),
		}, nil
	}

	// Создаем запрос для LLM
	request := LLMPipelineRequest{
		TaskID:       task.ID,
		TaskName:     task.Name,
		Description:  "Create a task execution pipeline",
		StepFile:     task.StepFile,
		Budgets:      task.Budgets,
		Dependencies: task.DependsOn,
		Context:      contextPack,
		Options: map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
			"max_tokens":  2048,
		},
	}

	// Генерируем промпт для LLM
	prompt := r.generatePipelinePrompt(request)

	// Получаем ответ от LLM с GBNF грамматикой
	response, err := r.llmClient.GenerateWithGBNF(ctx, prompt, "", request.Options)
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to generate pipeline with LLM: %v", err))
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        fmt.Sprintf("LLM generation failed: %v", err),
		}, nil
	}

	// Парсим ответ LLM
	policy, confidence, reasoning, err := r.parseLLMResponse([]byte(response.Content))
	if err != nil {
		r.log.Error(fmt.Sprintf("Failed to parse LLM response: %v", err))
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        fmt.Sprintf("Failed to parse LLM response: %v", err),
		}, nil
	}

	return &LLMPipelineResponse{
		Policy:       policy,
		Confidence:   confidence,
		Reasoning:    reasoning,
		FallbackUsed: false,
	}, nil
}

// generatePipelinePrompt генерирует промпт для LLM
func (r *LLMService) generatePipelinePrompt(request LLMPipelineRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are an AI assistant that creates task execution pipelines. ")
	prompt.WriteString("Generate a JSON response with a pipeline configuration for the following task:\n\n")

	prompt.WriteString(fmt.Sprintf("Task ID: %s\n", request.TaskID))
	prompt.WriteString(fmt.Sprintf("Task Name: %s\n", request.TaskName))
	prompt.WriteString(fmt.Sprintf("Description: %s\n", request.Description))
	prompt.WriteString(fmt.Sprintf("Step File: %s\n", request.StepFile))
	prompt.WriteString(fmt.Sprintf("Dependencies: %v\n", request.Dependencies))

	if request.Budgets.MaxFiles > 0 {
		prompt.WriteString(fmt.Sprintf("Max Files: %d\n", request.Budgets.MaxFiles))
	}
	if request.Budgets.MaxChangedLines > 0 {
		prompt.WriteString(fmt.Sprintf("Max Changed Lines: %d\n", request.Budgets.MaxChangedLines))
	}

	prompt.WriteString("\nAvailable pipeline steps:\n")
	prompt.WriteString("- retrieve: Extract context and dependencies\n")
	prompt.WriteString("- ast_synth: Synthesize AST and symbol graph\n")
	prompt.WriteString("- compile: Compile and build the project\n")
	prompt.WriteString("- test: Run tests and validation\n")
	prompt.WriteString("- static: Perform static analysis\n")
	prompt.WriteString("- format: Format code\n")
	prompt.WriteString("- validate: Validate results\n")
	prompt.WriteString("- repair: Fix issues and retry\n")

	prompt.WriteString("\nPipeline policy options:\n")
	prompt.WriteString("- EnableRetrieve: Enable context retrieval\n")
	prompt.WriteString("- EnableASTSynth: Enable AST synthesis\n")
	prompt.WriteString("- EnableCompile: Enable compilation\n")
	prompt.WriteString("- EnableTest: Enable testing\n")
	prompt.WriteString("- EnableStatic: Enable static analysis\n")
	prompt.WriteString("- EnableFormat: Enable formatting\n")
	prompt.WriteString("- EnableValidate: Enable validation\n")
	prompt.WriteString("- EnableRepair: Enable repair\n")
	prompt.WriteString("- FailFast: Stop on first error\n")
	prompt.WriteString("- RetryFailed: Retry failed steps\n")
	prompt.WriteString("- ParallelSteps: Execute steps in parallel\n")

	prompt.WriteString("\nGenerate a JSON response with the following structure:\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"schemaVersion\": \"1.0\",\n")
	prompt.WriteString("  \"edits\": [\n")
	prompt.WriteString("    {\n")
	prompt.WriteString("      \"kind\": \"pipelineConfig\",\n")
	prompt.WriteString("      \"path\": \"task-pipeline\",\n")
	prompt.WriteString("      \"language\": \"pipeline\",\n")
	prompt.WriteString("      \"operation\": {\n")
	prompt.WriteString("        \"engine\": \"llm-router\",\n")
	prompt.WriteString("        \"action\": \"createPipeline\",\n")
	prompt.WriteString("        \"params\": {\n")
	prompt.WriteString("          \"policy\": {\n")
	prompt.WriteString("            \"enableRetrieve\": true,\n")
	prompt.WriteString("            \"enableASTSynth\": true,\n")
	prompt.WriteString("            \"enableCompile\": false,\n")
	prompt.WriteString("            \"enableTest\": false,\n")
	prompt.WriteString("            \"enableStatic\": false,\n")
	prompt.WriteString("            \"enableFormat\": false,\n")
	prompt.WriteString("            \"enableValidate\": true,\n")
	prompt.WriteString("            \"enableRepair\": false,\n")
	prompt.WriteString("            \"failFast\": true,\n")
	prompt.WriteString("            \"retryFailed\": true,\n")
	prompt.WriteString("            \"parallelSteps\": false,\n")
	prompt.WriteString("            \"maxRetries\": 3,\n")
	prompt.WriteString("            \"timeout\": 1800000000000\n")
	prompt.WriteString("          },\n")
	prompt.WriteString("          \"confidence\": 0.95,\n")
	prompt.WriteString("          \"reasoning\": \"Explanation of pipeline configuration\"\n")
	prompt.WriteString("        }\n")
	prompt.WriteString("      },\n")
	prompt.WriteString("      \"post\": {\n")
	prompt.WriteString("        \"formatters\": [\"pipeline-validator\"]\n")
	prompt.WriteString("      }\n")
	prompt.WriteString("    }\n")
	prompt.WriteString("  ]\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

// llmPipelineResponse represents the structure of LLM pipeline response
type llmPipelineResponse struct {
	SchemaVersion string `json:"schemaVersion"`
	Edits         []struct {
		Kind      string `json:"kind"`
		Operation struct {
			Params struct {
				Policy     map[string]interface{} `json:"policy"`
				Confidence float64                `json:"confidence"`
				Reasoning  string                 `json:"reasoning"`
			} `json:"params"`
		} `json:"operation"`
	} `json:"edits"`
}

// defaultPipelinePolicy returns a default pipeline policy
func defaultPipelinePolicy() *PipelinePolicy {
	return &PipelinePolicy{
		EnableRetrieve: true, EnableASTSynth: true, EnableCompile: false, EnableTest: false,
		EnableStatic: false, EnableRepair: false, EnableFormat: false, EnableValidate: true,
		FailFast: true, RetryFailed: true, MaxRetries: 3, ParallelSteps: false, Timeout: 30 * time.Minute,
	}
}

// applyPolicyFromMap applies policy values from a map to a PipelinePolicy
func applyPolicyFromMap(policy *PipelinePolicy, m map[string]interface{}) {
	if m == nil {
		return
	}
	if v, ok := m["enableRetrieve"].(bool); ok {
		policy.EnableRetrieve = v
	}
	if v, ok := m["enableASTSynth"].(bool); ok {
		policy.EnableASTSynth = v
	}
	if v, ok := m["enableCompile"].(bool); ok {
		policy.EnableCompile = v
	}
	if v, ok := m["enableTest"].(bool); ok {
		policy.EnableTest = v
	}
	if v, ok := m["enableStatic"].(bool); ok {
		policy.EnableStatic = v
	}
	if v, ok := m["enableFormat"].(bool); ok {
		policy.EnableFormat = v
	}
	if v, ok := m["enableValidate"].(bool); ok {
		policy.EnableValidate = v
	}
	if v, ok := m["enableRepair"].(bool); ok {
		policy.EnableRepair = v
	}
	if v, ok := m["failFast"].(bool); ok {
		policy.FailFast = v
	}
	if v, ok := m["retryFailed"].(bool); ok {
		policy.RetryFailed = v
	}
	if v, ok := m["parallelSteps"].(bool); ok {
		policy.ParallelSteps = v
	}
	if v, ok := m["maxRetries"].(float64); ok {
		policy.MaxRetries = int(v)
	}
	if v, ok := m["timeout"].(float64); ok {
		policy.Timeout = time.Duration(v)
	}
}

// parseLLMResponse парсит ответ от LLM
func (r *LLMService) parseLLMResponse(response []byte) (*PipelinePolicy, float64, string, error) {
	var llmResponse llmPipelineResponse
	if err := json.Unmarshal(response, &llmResponse); err != nil {
		return nil, 0, "", fmt.Errorf("failed to unmarshal LLM response: %w", err)
	}

	if len(llmResponse.Edits) == 0 {
		return nil, 0, "", fmt.Errorf("no pipeline configuration in LLM response")
	}

	edit := llmResponse.Edits[0]
	if edit.Kind != "pipelineConfig" {
		return nil, 0, "", fmt.Errorf("unexpected edit kind: %s", edit.Kind)
	}

	policy := defaultPipelinePolicy()
	applyPolicyFromMap(policy, edit.Operation.Params.Policy)

	return policy, edit.Operation.Params.Confidence, edit.Operation.Params.Reasoning, nil
}

// ValidatePipelineWithLLM валидирует пайплайн с помощью LLM
func (r *LLMService) ValidatePipelineWithLLM(ctx context.Context, pipeline *TaskPipeline) (*LLMPipelineResponse, error) {
	if !r.IsEnabled() {
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        "LLM router is disabled",
		}, nil
	}

	// Создаем промпт для валидации
	prompt := r.generateValidationPrompt(pipeline)

	// Получаем ответ от LLM
	response, err := r.llmClient.GenerateEditsJSON(ctx, prompt, map[string]interface{}{
		"temperature": 0.1,
		"top_p":       0.9,
		"max_tokens":  1024,
	})
	if err != nil {
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        fmt.Sprintf("LLM validation failed: %v", err),
		}, nil
	}

	// Парсим ответ валидации
	validation, confidence, reasoning, err := r.parseValidationResponse(response)
	if err != nil {
		return &LLMPipelineResponse{
			FallbackUsed: true,
			Error:        fmt.Sprintf("Failed to parse validation response: %v", err),
		}, nil
	}

	return &LLMPipelineResponse{
		Policy:       validation,
		Confidence:   confidence,
		Reasoning:    reasoning,
		FallbackUsed: false,
	}, nil
}

// generateValidationPrompt генерирует промпт для валидации пайплайна
func (r *LLMService) generateValidationPrompt(pipeline *TaskPipeline) string {
	var prompt strings.Builder

	prompt.WriteString("You are an AI assistant that validates task execution pipelines. ")
	prompt.WriteString("Analyze the following pipeline and provide validation feedback:\n\n")

	prompt.WriteString(fmt.Sprintf("Task ID: %s\n", pipeline.TaskID))
	prompt.WriteString(fmt.Sprintf("Number of Steps: %d\n", len(pipeline.Steps)))
	prompt.WriteString(fmt.Sprintf("Policy: %+v\n", pipeline.Policy))

	prompt.WriteString("\nPipeline Steps:\n")
	for i, step := range pipeline.Steps {
		prompt.WriteString(fmt.Sprintf("%d. %s (%s) - Priority: %d\n",
			i+1, step.Name, step.Type, step.Priority))
	}

	prompt.WriteString("\nGenerate a JSON response with validation results:\n")
	prompt.WriteString("{\n")
	prompt.WriteString("  \"schemaVersion\": \"1.0\",\n")
	prompt.WriteString("  \"edits\": [\n")
	prompt.WriteString("    {\n")
	prompt.WriteString("      \"kind\": \"pipelineValidation\",\n")
	prompt.WriteString("      \"path\": \"pipeline-validation\",\n")
	prompt.WriteString("      \"language\": \"validation\",\n")
	prompt.WriteString("      \"operation\": {\n")
	prompt.WriteString("        \"engine\": \"llm-validator\",\n")
	prompt.WriteString("        \"action\": \"validatePipeline\",\n")
	prompt.WriteString("        \"params\": {\n")
	prompt.WriteString("          \"isValid\": true,\n")
	prompt.WriteString("          \"confidence\": 0.95,\n")
	prompt.WriteString("          \"reasoning\": \"Validation explanation\",\n")
	prompt.WriteString("          \"suggestions\": [\"suggestion1\", \"suggestion2\"],\n")
	prompt.WriteString("          \"warnings\": [\"warning1\", \"warning2\"]\n")
	prompt.WriteString("        }\n")
	prompt.WriteString("      },\n")
	prompt.WriteString("      \"post\": {\n")
	prompt.WriteString("        \"formatters\": [\"validation-formatter\"]\n")
	prompt.WriteString("      }\n")
	prompt.WriteString("    }\n")
	prompt.WriteString("  ]\n")
	prompt.WriteString("}\n")

	return prompt.String()
}

// parseValidationResponse парсит ответ валидации от LLM
func (r *LLMService) parseValidationResponse(response []byte) (*PipelinePolicy, float64, string, error) {
	var validationResponse llmPipelineResponse
	if err := json.Unmarshal(response, &validationResponse); err != nil {
		return nil, 0, "", fmt.Errorf("failed to unmarshal validation response: %w", err)
	}

	if len(validationResponse.Edits) == 0 {
		return nil, 0, "", fmt.Errorf("no validation result in LLM response")
	}

	validation := validationResponse.Edits[0]
	policy := &PipelinePolicy{}
	applyPolicyFromMap(policy, validation.Operation.Params.Policy)

	return policy, validation.Operation.Params.Confidence, validation.Operation.Params.Reasoning, nil
}

// LoadGBNFGrammar загружает GBNF грамматику из файла
func (r *LLMService) LoadGBNFGrammar(grammarPath string) (string, error) {
	if grammarPath == "" {
		grammarPath = "docs/schemas/edits.gbnf"
	}

	data, err := r.fileReader.ReadFile(grammarPath)
	if err != nil {
		return "", fmt.Errorf("failed to read GBNF grammar file: %w", err)
	}

	return string(data), nil
}

// GetLLMStatus возвращает статус LLM сервиса
func (r *LLMService) GetLLMStatus(ctx context.Context) map[string]interface{} {
	status := map[string]interface{}{
		"enabled":    r.enabled,
		"llm_client": r.llmClient != nil,
	}

	if r.llmClient != nil {
		// Проверяем доступность LLM сервиса
		if err := r.llmClient.HealthCheck(ctx); err != nil {
			status["llm_server_available"] = false
			status["llm_server_error"] = err.Error()
		} else {
			status["llm_server_available"] = true

			// Получаем информацию о модели
			if modelInfo, err := r.llmClient.GetModelInfo(ctx); err == nil {
				status["model_info"] = modelInfo
			}
		}
	}

	return status
}
