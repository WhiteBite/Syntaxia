package router

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
	"time"
)

// ExecutePipeline выполняет пайплайн
func (r *PlannerService) ExecutePipeline(ctx context.Context, pipeline *TaskPipeline) error {
	r.log.Info(fmt.Sprintf("Executing pipeline for task: %s", pipeline.TaskID))

	now := time.Now()
	pipeline.StartedAt = &now
	pipeline.Status = PipelineStatusRunning

	// Выполняем шаги последовательно или параллельно
	if pipeline.Policy.ParallelSteps {
		return r.executePipelineParallel(ctx, pipeline)
	}
	return r.executePipelineSequential(ctx, pipeline)
}

// executePipelineSequential выполняет пайплайн последовательно
func (r *PlannerService) executePipelineSequential(ctx context.Context, pipeline *TaskPipeline) error {
	for _, step := range pipeline.Steps {
		if err := r.executeStep(ctx, step); err != nil {
			if pipeline.Policy.FailFast {
				pipeline.Status = PipelineStatusFailed
				pipeline.Error = err.Error()
				now := time.Now()
				pipeline.CompletedAt = &now
				pipeline.Duration = now.Sub(*pipeline.StartedAt)
				return fmt.Errorf("step %s failed: %w", step.ID, err)
			}
			// Продолжаем выполнение других шагов
		}
	}

	pipeline.Status = PipelineStatusCompleted
	now := time.Now()
	pipeline.CompletedAt = &now
	pipeline.Duration = now.Sub(*pipeline.StartedAt)

	r.log.Info(fmt.Sprintf("Pipeline completed for task: %s", pipeline.TaskID))
	return nil
}

// executePipelineParallel выполняет пайплайн параллельно
// В реальной реализации нужно учитывать зависимости между шагами
func (r *PlannerService) executePipelineParallel(ctx context.Context, pipeline *TaskPipeline) error {
	return r.executePipelineSequential(ctx, pipeline)
}

// executeStep выполняет один шаг пайплайна
func (r *PlannerService) executeStep(ctx context.Context, step *TaskPipelineStep) error {
	r.log.Info(fmt.Sprintf("Executing step: %s (%s)", step.Name, step.ID))

	now := time.Now()
	step.StartedAt = &now
	step.Status = StepStatusRunning

	var err error
	switch step.Type {
	case StepTypeRetrieve:
		err = r.executeRetrieveStep(ctx, step)
	case StepTypeASTSynth:
		err = r.executeASTSynthStep(ctx, step)
	case StepTypeCompile:
		err = r.executeCompileStep(ctx, step)
	case StepTypeTest:
		err = r.executeTestStep(ctx, step)
	case StepTypeStatic:
		err = r.executeStaticStep(ctx, step)
	case StepTypeFormat:
		err = r.executeFormatStep(ctx, step)
	case StepTypeValidate:
		err = r.executeValidateStep(ctx, step)
	case StepTypeRepair:
		err = r.executeRepairStep(ctx, step)
	default:
		err = fmt.Errorf("unknown step type: %s", step.Type)
	}

	// Обновляем статус шага
	completedAt := time.Now()
	step.CompletedAt = &completedAt
	step.Duration = completedAt.Sub(*step.StartedAt)

	if err != nil {
		step.Status = StepStatusFailed
		step.Error = err.Error()
		r.log.Error(fmt.Sprintf("Step %s failed: %v", step.ID, err))
	} else {
		step.Status = StepStatusCompleted
		step.Result = &TaskPipelineStepResult{
			Success: true,
			Message: fmt.Sprintf("Step %s completed successfully", step.Name),
		}
		r.log.Info(fmt.Sprintf("Step %s completed successfully", step.ID))
	}

	return err
}

// executeRetrieveStep выполняет шаг извлечения контекста
func (r *PlannerService) executeRetrieveStep(_ context.Context, _ *TaskPipelineStep) error {
	// Симуляция выполнения шага retrieve
	time.Sleep(100 * time.Millisecond)
	return nil
}

// executeASTSynthStep выполняет шаг синтеза AST
func (r *PlannerService) executeASTSynthStep(_ context.Context, _ *TaskPipelineStep) error {
	// Симуляция выполнения шага ast_synth
	time.Sleep(200 * time.Millisecond)
	return nil
}

// executeCompileStep выполняет шаг компиляции
func (r *PlannerService) executeCompileStep(ctx context.Context, step *TaskPipelineStep) error {
	projectPath, ok := step.Config["project_path"].(string)
	if !ok {
		return fmt.Errorf("project_path not found in step config")
	}

	language, ok := step.Config["language"].(string)
	if !ok {
		language = "go"
	}

	result, err := r.buildService.Build(ctx, projectPath, language)
	if err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	step.Result = &TaskPipelineStepResult{
		Success: result.Success,
		Message: result.Output,
		Data: map[string]any{
			"duration": result.Duration,
		},
		Artifacts: result.Artifacts,
	}

	if !result.Success {
		return fmt.Errorf("compilation failed: %s", result.Output)
	}

	return nil
}

// executeTestStep выполняет шаг тестирования
func (r *PlannerService) executeTestStep(ctx context.Context, step *TaskPipelineStep) error {
	projectPath, ok := step.Config["project_path"].(string)
	if !ok {
		if pipeline, ok := ctx.Value("pipeline").(*TaskPipeline); ok {
			if task, ok := pipeline.Steps[0].Config["task"].(domain.Task); ok {
				projectPath = task.Metadata["project_path"].(string)
			}
		}
		if projectPath == "" {
			return fmt.Errorf("project_path not found for test step")
		}
	}

	testConfig := &domain.TestConfig{
		ProjectPath: projectPath,
		Language:    "go",
		Scope:       domain.TestScopeAll,
		Coverage:    true,
		Timeout:     int((5 * time.Minute).Seconds()),
	}

	results, err := r.testService.RunTests(ctx, testConfig)
	if err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	success := true
	messages := make([]string, 0, len(results))
	totalDuration := 0.0

	for _, res := range results {
		if !res.Success {
			success = false
		}
		messages = append(messages, fmt.Sprintf("Test: %s, Success: %v, Output: %s", res.TestName, res.Success, res.Output))
		totalDuration += res.Duration
	}

	step.Result = &TaskPipelineStepResult{
		Success: success,
		Message: strings.Join(messages, "\n"),
		Data: map[string]any{
			"duration":  totalDuration,
			"tests_run": len(results),
		},
	}

	if !success {
		return fmt.Errorf("one or more tests failed")
	}

	return nil
}

// executeStaticStep выполняет шаг статического анализа
func (r *PlannerService) executeStaticStep(ctx context.Context, step *TaskPipelineStep) error {
	projectPath, ok := step.Config["project_path"].(string)
	if !ok {
		if pipeline, ok := ctx.Value("pipeline").(*TaskPipeline); ok {
			if task, ok := pipeline.Steps[0].Config["task"].(domain.Task); ok {
				projectPath = task.Metadata["project_path"].(string)
			}
		}
		if projectPath == "" {
			return fmt.Errorf("project_path not found for static analysis step")
		}
	}

	languages, ok := step.Config["languages"].([]string)
	if !ok {
		languages = []string{"go"}
	}

	failOnError, ok := step.Config["fail_on_error"].(bool)
	if !ok {
		failOnError = false
	}

	results, err := r.staticAnalyzer.AnalyzeProject(ctx, projectPath, languages)
	if err != nil {
		return fmt.Errorf("static analysis failed: %w", err)
	}

	success := true
	totalIssues := 0
	totalDuration := 0.0

	for _, result := range results.Results {
		totalIssues += len(result.Issues)
		totalDuration += result.Duration
	}

	if failOnError && totalIssues > 0 {
		success = false
	}

	step.Result = &TaskPipelineStepResult{
		Success: success,
		Message: fmt.Sprintf("Static analysis completed with %d issues.", totalIssues),
		Data: map[string]any{
			"total_issues": totalIssues,
			"duration":     totalDuration,
		},
	}

	if !success {
		return fmt.Errorf("static analysis found critical issues")
	}

	return nil
}

// executeFormatStep выполняет шаг форматирования
func (r *PlannerService) executeFormatStep(_ context.Context, _ *TaskPipelineStep) error {
	// Симуляция выполнения шага format
	time.Sleep(150 * time.Millisecond)
	return nil
}

// executeValidateStep выполняет шаг валидации
func (r *PlannerService) executeValidateStep(_ context.Context, _ *TaskPipelineStep) error {
	// Симуляция выполнения шага validate
	time.Sleep(100 * time.Millisecond)
	return nil
}

// executeRepairStep выполняет шаг исправления
func (r *PlannerService) executeRepairStep(ctx context.Context, step *TaskPipelineStep) error {
	projectPath, ok := step.Config["project_path"].(string)
	if !ok {
		if pipeline, ok := ctx.Value("pipeline").(*TaskPipeline); ok {
			if task, ok := pipeline.Steps[0].Config["task"].(domain.Task); ok {
				projectPath = task.Metadata["project_path"].(string)
			}
		}
		if projectPath == "" {
			return fmt.Errorf("project_path not found for repair step")
		}
	}

	errorOutput, ok := step.Config["error_output"].(string)
	if !ok {
		step.Result = &TaskPipelineStepResult{
			Success: true,
			Message: "No errors to repair.",
		}
		return nil
	}

	language, ok := step.Config["language"].(string)
	if !ok {
		language = "go"
	}

	maxAttempts, ok := step.Config["max_attempts"].(int)
	if !ok {
		maxAttempts = 3
	}

	repairRequest := domain.RepairRequest{
		ProjectPath: projectPath,
		ErrorOutput: errorOutput,
		Language:    language,
		MaxAttempts: maxAttempts,
	}

	result, err := r.repairService.ExecuteRepair(ctx, repairRequest)
	if err != nil {
		return fmt.Errorf("repair execution failed: %w", err)
	}

	step.Result = &TaskPipelineStepResult{
		Success: result.Success,
		Message: fmt.Sprintf("Repair finished. Success: %v. Attempts: %d.", result.Success, result.Attempts),
		Data: map[string]any{
			"attempts": result.Attempts,
		},
	}

	if !result.Success {
		return fmt.Errorf("automated repair failed after %d attempts", result.Attempts)
	}

	return nil
}
