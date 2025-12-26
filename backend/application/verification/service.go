package verification

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"time"
)

// FormatterService interface to avoid circular imports
type FormatterService interface {
	FormatProject(ctx context.Context, projectPath, language string) (*domain.FormatResult, error)
}

// Service предоставляет высокоуровневый API для verification pipeline
type Service struct {
	log              domain.Logger
	buildService     domain.IBuildService
	testService      domain.ITestService
	staticAnalyzer   domain.IStaticAnalyzerService
	formatterService FormatterService
	reportWriter     domain.FileSystemWriter
	taskProtocol     domain.TaskProtocolService
}

// NewService создает новый сервис verification pipeline
func NewService(
	log domain.Logger,
	buildService domain.IBuildService,
	testService domain.ITestService,
	staticAnalyzer domain.IStaticAnalyzerService,
	formatterService FormatterService,
	reportWriter domain.FileSystemWriter,
	taskProtocol domain.TaskProtocolService,
) *Service {
	return &Service{
		log:              log,
		buildService:     buildService,
		testService:      testService,
		staticAnalyzer:   staticAnalyzer,
		formatterService: formatterService,
		reportWriter:     reportWriter,
		taskProtocol:     taskProtocol,
	}
}

// RunVerificationPipeline выполняет полный verification pipeline
func (s *Service) RunVerificationPipeline(ctx context.Context, config *domain.VerificationConfig) (*domain.VerificationResult, error) {
	s.log.Info(fmt.Sprintf("Starting verification pipeline for project: %s", config.ProjectPath))

	result := &domain.VerificationResult{
		ProjectPath: config.ProjectPath,
		Languages:   config.Languages,
		StartedAt:   time.Now().UTC().Format(time.RFC3339),
		Steps:       make([]*domain.VerificationStep, 0),
	}

	// Шаг 1: Форматирование
	formatStep := s.runStep(ctx, "format", config, s.runFormatStep)
	result.Steps = append(result.Steps, formatStep)

	// Шаг 2: Build и Type Check
	buildStep := s.runStep(ctx, "build-typecheck", config, s.runBuildTypeCheckStep)
	result.Steps = append(result.Steps, buildStep)
	if buildStep.Error != nil {
		result.Success = false
		result.CompletedAt = time.Now().UTC().Format(time.RFC3339)
		return result, fmt.Errorf("build/typecheck failed: %w", buildStep.Error)
	}

	// Шаг 3: Smoke Tests
	testStep := s.runStep(ctx, "smoke-tests", config, s.runSmokeTestsStep)
	result.Steps = append(result.Steps, testStep)
	if testStep.Error != nil {
		result.Success = false
		result.CompletedAt = time.Now().UTC().Format(time.RFC3339)
		return result, fmt.Errorf("smoke tests failed: %w", testStep.Error)
	}

	// Шаг 4: Static Analysis
	staticStep := s.runStep(ctx, "static-analysis", config, s.runStaticAnalysisStep)
	result.Steps = append(result.Steps, staticStep)

	// Определяем общий успех
	result.Success = true
	for _, step := range result.Steps {
		if !step.Success && step.Name != "format" && step.Name != "static-analysis" {
			result.Success = false
			break
		}
	}

	result.CompletedAt = time.Now().UTC().Format(time.RFC3339)

	if err := s.saveVerificationReport(ctx, result); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to save verification report: %v", err))
	}

	s.log.Info(fmt.Sprintf("Verification pipeline completed with success: %t", result.Success))
	return result, nil
}

type stepFunc func(ctx context.Context, config *domain.VerificationConfig) (interface{}, error)

func (s *Service) runStep(ctx context.Context, name string, config *domain.VerificationConfig, fn stepFunc) *domain.VerificationStep {
	step := &domain.VerificationStep{
		Name:      name,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}

	result, err := fn(ctx, config)
	step.Result = result
	step.Success = err == nil
	step.Error = err
	step.CompletedAt = time.Now().UTC().Format(time.RFC3339)

	if err != nil {
		s.log.Warning(fmt.Sprintf("%s step failed: %v", name, err))
	}

	return step
}

func (s *Service) runFormatStep(ctx context.Context, config *domain.VerificationConfig) (interface{}, error) {
	s.log.Info("Running format step")
	results := make([]*domain.FormatResult, 0, len(config.Languages))

	for _, language := range config.Languages {
		result, err := s.formatterService.FormatProject(ctx, config.ProjectPath, language)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to format %s: %v", language, err))
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

func (s *Service) runBuildTypeCheckStep(ctx context.Context, config *domain.VerificationConfig) (interface{}, error) {
	s.log.Info("Running build/typecheck step")

	validation, err := s.buildService.ValidateProject(ctx, config.ProjectPath, config.Languages)
	if err != nil {
		return nil, fmt.Errorf("project validation failed: %w", err)
	}

	if !validation.Success {
		return nil, fmt.Errorf("project validation failed for some languages")
	}

	return validation, nil
}

func (s *Service) runSmokeTestsStep(ctx context.Context, config *domain.VerificationConfig) (interface{}, error) {
	s.log.Info("Running smoke tests step")

	var allResults []*domain.TestResult

	for _, language := range config.Languages {
		results, err := s.testService.RunSmokeTests(ctx, config.ProjectPath, language)
		if err != nil {
			return nil, fmt.Errorf("smoke tests failed for %s: %w", language, err)
		}

		validation := s.testService.ValidateTestResults(results)
		if !validation.Success {
			return nil, fmt.Errorf("smoke tests failed for %s: %d tests failed", language, validation.FailedTests)
		}

		allResults = append(allResults, results...)
	}

	return allResults, nil
}

func (s *Service) runStaticAnalysisStep(ctx context.Context, config *domain.VerificationConfig) (interface{}, error) {
	s.log.Info("Running static analysis step")

	report, err := s.staticAnalyzer.AnalyzeProject(ctx, config.ProjectPath, config.Languages)
	if err != nil {
		return nil, fmt.Errorf("static analysis failed: %w", err)
	}

	return report, nil
}

func (s *Service) saveVerificationReport(ctx context.Context, result *domain.VerificationResult) error {
	reportsDir := "tasks/reports"
	if err := s.reportWriter.MkdirAll(reportsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102_150405")
	filename := fmt.Sprintf("verification_%s.json", timestamp)
	fpath := filepath.Join(reportsDir, filename)

	jsonData, err := s.serializeVerificationResult(result)
	if err != nil {
		return fmt.Errorf("failed to serialize verification result: %w", err)
	}

	if err := s.reportWriter.WriteFile(fpath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write verification report: %w", err)
	}

	s.log.Info(fmt.Sprintf("Saved verification report to %s", fpath))
	return nil
}

func (s *Service) serializeVerificationResult(result *domain.VerificationResult) ([]byte, error) {
	content := fmt.Sprintf(`{
  "projectPath": "%s",
  "languages": %v,
  "success": %t,
  "startedAt": "%s",
  "completedAt": "%s",
  "steps": %d
}`, result.ProjectPath, result.Languages, result.Success, result.StartedAt, result.CompletedAt, len(result.Steps))

	return []byte(content), nil
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (s *Service) GetSupportedLanguages() []string {
	return s.buildService.GetSupportedLanguages()
}

// DetectLanguages определяет языки в проекте
func (s *Service) DetectLanguages(ctx context.Context, projectPath string) ([]string, error) {
	return s.buildService.DetectLanguages(ctx, projectPath)
}

// RunTaskProtocolVerification executes the Task Protocol verification for a task
func (s *Service) RunTaskProtocolVerification(ctx context.Context, config *domain.TaskProtocolConfig) (*domain.TaskProtocolResult, error) {
	s.log.Info(fmt.Sprintf("Starting Task Protocol verification for project: %s", config.ProjectPath))

	if s.taskProtocol == nil {
		return nil, fmt.Errorf("task protocol service not available")
	}

	result, err := s.taskProtocol.ExecuteProtocol(ctx, config)
	if err != nil {
		s.log.Error(fmt.Sprintf("Task Protocol verification failed: %v", err))
		return result, err
	}

	if err := s.saveTaskProtocolReport(ctx, result); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to save task protocol report: %v", err))
	}

	s.log.Info(fmt.Sprintf("Task Protocol verification completed with success: %t", result.Success))
	return result, nil
}

// RunTaskProtocolStage executes a single Task Protocol verification stage
func (s *Service) RunTaskProtocolStage(ctx context.Context, stage domain.ProtocolStage, config *domain.TaskProtocolConfig) (*domain.ProtocolStageResult, error) {
	s.log.Info(fmt.Sprintf("Running Task Protocol stage: %s", stage))

	if s.taskProtocol == nil {
		return nil, fmt.Errorf("task protocol service not available")
	}

	return s.taskProtocol.ValidateStage(ctx, stage, config)
}

// CreateTaskProtocolConfig creates a TaskProtocolConfig from VerificationConfig
func (s *Service) CreateTaskProtocolConfig(verifyConfig *domain.VerificationConfig) *domain.TaskProtocolConfig {
	return &domain.TaskProtocolConfig{
		ProjectPath: verifyConfig.ProjectPath,
		Languages:   verifyConfig.Languages,
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
}

func (s *Service) saveTaskProtocolReport(ctx context.Context, result *domain.TaskProtocolResult) error {
	reportsDir := "tasks/reports/protocols"
	if err := s.reportWriter.MkdirAll(reportsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create protocol reports directory: %w", err)
	}

	timestamp := time.Now().UTC().Format("20060102_150405")
	filename := fmt.Sprintf("task_protocol_%s_%s.json", result.TaskID, timestamp)
	fpath := filepath.Join(reportsDir, filename)

	jsonData, err := s.serializeTaskProtocolResult(result)
	if err != nil {
		return fmt.Errorf("failed to serialize task protocol result: %w", err)
	}

	if err := s.reportWriter.WriteFile(fpath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write task protocol report: %w", err)
	}

	s.log.Info(fmt.Sprintf("Saved task protocol report to %s", fpath))
	return nil
}

func (s *Service) serializeTaskProtocolResult(result *domain.TaskProtocolResult) ([]byte, error) {
	content := fmt.Sprintf(`{
  "taskId": "%s",
  "success": %t,
  "startedAt": "%s",
  "completedAt": "%s",
  "correctionCycles": %d,
  "stages": %d,
  "finalError": "%s"
}`,
		result.TaskID,
		result.Success,
		result.StartedAt.Format(time.RFC3339),
		result.CompletedAt.Format(time.RFC3339),
		result.CorrectionCycles,
		len(result.Stages),
		result.FinalError)

	return []byte(content), nil
}
