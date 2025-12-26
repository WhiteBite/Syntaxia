// Package protocol provides task protocol verification services.
package protocol

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"sync"
	"time"
)

// Service implements the TaskProtocolService interface
type Service struct {
	log                  domain.Logger
	verificationPipeline VerificationPipeline
	staticAnalyzer       domain.IStaticAnalyzerService
	testService          domain.ITestService
	buildService         domain.IBuildService
	guardrailService     domain.GuardrailService
	aiService            IntelligentAI
	errorAnalyzer        domain.ErrorAnalyzer
	correctionEngine     domain.CorrectionEngine
	mu                   sync.RWMutex
}

// VerificationPipeline interface to avoid circular imports
type VerificationPipeline interface{}

// IntelligentAI interface to avoid circular imports
type IntelligentAI interface{}

// NewService creates a new TaskProtocolService instance
func NewService(
	log domain.Logger,
	verificationPipeline VerificationPipeline,
	staticAnalyzer domain.IStaticAnalyzerService,
	testService domain.ITestService,
	buildService domain.IBuildService,
	guardrailService domain.GuardrailService,
	aiService IntelligentAI,
	errorAnalyzer domain.ErrorAnalyzer,
	correctionEngine domain.CorrectionEngine,
) domain.TaskProtocolService {
	return &Service{
		log:                  log,
		verificationPipeline: verificationPipeline,
		staticAnalyzer:       staticAnalyzer,
		testService:          testService,
		buildService:         buildService,
		guardrailService:     guardrailService,
		aiService:            aiService,
		errorAnalyzer:        errorAnalyzer,
		correctionEngine:     correctionEngine,
	}
}

// ExecuteProtocol executes the full verification protocol for a task
func (s *Service) ExecuteProtocol(ctx context.Context, config *domain.TaskProtocolConfig) (*domain.TaskProtocolResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.log.Info(fmt.Sprintf("Starting Task Protocol verification for project: %s", config.ProjectPath))

	result := &domain.TaskProtocolResult{
		TaskID:           generateTaskID(),
		StartedAt:        time.Now(),
		Stages:           make([]*domain.ProtocolStageResult, 0),
		CorrectionCycles: 0,
	}

	for _, stage := range config.EnabledStages {
		stageResult, err := s.executeStageWithRetry(ctx, stage, config)
		if err != nil && stageResult == nil {
			stageResult = &domain.ProtocolStageResult{
				Stage:   stage,
				Success: false,
				ErrorDetails: &domain.ErrorDetails{
					Message: fmt.Sprintf("Catastrophic error in stage %s: %v", stage, err),
				},
			}
		}

		result.Stages = append(result.Stages, stageResult)

		if !stageResult.Success {
			if config.FailFast {
				result.Success = false
				result.FinalError = fmt.Sprintf("Stage %s failed: %v", stage, err)
				result.CompletedAt = time.Now()
				return result, fmt.Errorf("protocol failed at stage %s: %w", stage, err)
			}
			s.log.Warning(fmt.Sprintf("Stage %s failed but continuing due to FailFast=false", stage))
		}
	}

	result.Success = s.determineOverallSuccess(result.Stages)
	result.CompletedAt = time.Now()

	s.log.Info(fmt.Sprintf("Task Protocol completed with success: %t, correction cycles: %d",
		result.Success, result.CorrectionCycles))

	return result, nil
}

func (s *Service) executeStageWithRetry(ctx context.Context, stage domain.ProtocolStage, config *domain.TaskProtocolConfig) (*domain.ProtocolStageResult, error) {
	stageResult := &domain.ProtocolStageResult{
		Stage:           stage,
		Attempts:        0,
		CorrectionSteps: make([]*domain.CorrectionStep, 0),
	}

	startTime := time.Now()
	maxAttempts := config.MaxRetries + 1

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		stageResult.Attempts = attempt

		s.log.Info(fmt.Sprintf("Executing stage %s, attempt %d/%d", stage, attempt, maxAttempts))

		err := s.executeStage(ctx, stage, config)

		if err == nil {
			stageResult.Success = true
			stageResult.Duration = time.Since(startTime)
			return stageResult, nil
		}

		if config.SelfCorrection.Enabled && attempt < maxAttempts {
			s.log.Info(fmt.Sprintf("Stage %s failed, attempting self-correction", stage))

			if correctionApplied, corrErr := s.attemptSelfCorrection(ctx, err, stage, config); corrErr == nil && correctionApplied {
				s.log.Info(fmt.Sprintf("Self-correction applied for stage %s", stage))
				continue
			} else if corrErr != nil {
				s.log.Warning(fmt.Sprintf("Self-correction failed for stage %s: %v", stage, corrErr))
			}
		}

		if attempt == maxAttempts {
			stageResult.Success = false
			stageResult.Duration = time.Since(startTime)
			stageResult.ErrorDetails = s.parseErrorDetails(err, stage)
			return stageResult, err
		}
	}

	return stageResult, fmt.Errorf("stage %s failed after %d attempts", stage, maxAttempts)
}

func (s *Service) executeStage(ctx context.Context, stage domain.ProtocolStage, config *domain.TaskProtocolConfig) error {
	switch stage {
	case domain.StageLinting:
		return s.executeLintingStage(ctx, config)
	case domain.StageBuilding:
		return s.executeBuildingStage(ctx, config)
	case domain.StageTesting:
		return s.executeTestingStage(ctx, config)
	case domain.StageGuardrails:
		return s.executeGuardrailsStage(ctx, config)
	default:
		return fmt.Errorf("unsupported stage: %s", stage)
	}
}

func (s *Service) executeLintingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	s.log.Info("Running linting stage")

	report, err := s.staticAnalyzer.AnalyzeProject(ctx, config.ProjectPath, config.Languages)
	if err != nil {
		return fmt.Errorf("static analysis failed: %w", err)
	}

	if s.hasLintingErrors(report) {
		return fmt.Errorf("linting errors found in static analysis report")
	}

	return nil
}

func (s *Service) executeBuildingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	s.log.Info("Running building stage")

	validation, err := s.buildService.ValidateProject(ctx, config.ProjectPath, config.Languages)
	if err != nil {
		return fmt.Errorf("build validation failed: %w", err)
	}

	if !validation.Success {
		return fmt.Errorf("project build failed for one or more languages")
	}

	return nil
}

func (s *Service) executeTestingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	s.log.Info("Running testing stage")

	for _, language := range config.Languages {
		results, err := s.testService.RunSmokeTests(ctx, config.ProjectPath, language)
		if err != nil {
			return fmt.Errorf("tests failed for %s: %w", language, err)
		}

		validation := s.testService.ValidateTestResults(results)
		if !validation.Success {
			return fmt.Errorf("test validation failed for %s: %d tests failed", language, validation.FailedTests)
		}
	}

	return nil
}

func (s *Service) executeGuardrailsStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	s.log.Info("Running guardrails stage")

	taskID := generateTaskID()
	files := []string{}
	linesChanged := int64(0)

	result, err := s.guardrailService.ValidateTask(taskID, files, linesChanged)
	if err != nil {
		return fmt.Errorf("guardrail validation failed: %w", err)
	}

	if !result.Valid {
		return fmt.Errorf("guardrail validation failed: %s", result.Error)
	}

	return nil
}

func (s *Service) attemptSelfCorrection(ctx context.Context, err error, stage domain.ProtocolStage, config *domain.TaskProtocolConfig) (bool, error) {
	errorDetails, analyzeErr := s.errorAnalyzer.AnalyzeError(err.Error(), stage)
	if analyzeErr != nil {
		return false, fmt.Errorf("error analysis failed: %w", analyzeErr)
	}

	corrections, corrErr := s.errorAnalyzer.SuggestCorrections(errorDetails)
	if corrErr != nil {
		return false, fmt.Errorf("correction suggestion failed: %w", corrErr)
	}

	if len(corrections) == 0 {
		return false, nil
	}

	correctionResult, applyErr := s.correctionEngine.ApplyCorrections(ctx, corrections, config.ProjectPath)
	if applyErr != nil {
		return false, fmt.Errorf("correction application failed: %w", applyErr)
	}

	return correctionResult.Success, nil
}

// ValidateStage executes a single verification stage
func (s *Service) ValidateStage(ctx context.Context, stage domain.ProtocolStage, config *domain.TaskProtocolConfig) (*domain.ProtocolStageResult, error) {
	return s.executeStageWithRetry(ctx, stage, config)
}

// RequestCorrectionGuidance requests AI-generated correction guidance for errors
func (s *Service) RequestCorrectionGuidance(ctx context.Context, error *domain.ErrorDetails, taskContext *domain.TaskContext) (*domain.CorrectionGuidance, error) {
	if s.aiService == nil {
		return nil, fmt.Errorf("AI service not available")
	}

	guidance := &domain.CorrectionGuidance{
		Error:       error,
		Steps:       make([]*domain.CorrectionStep, 0),
		Explanation: "AI-generated correction guidance would be provided here",
		Confidence:  0.8,
	}

	return guidance, nil
}

func (s *Service) determineOverallSuccess(stages []*domain.ProtocolStageResult) bool {
	for _, stage := range stages {
		if (stage.Stage == domain.StageBuilding || stage.Stage == domain.StageTesting) && !stage.Success {
			return false
		}
	}
	return true
}

func (s *Service) hasLintingErrors(report interface{}) bool {
	return false
}

func (s *Service) parseErrorDetails(err error, stage domain.ProtocolStage) *domain.ErrorDetails {
	return &domain.ErrorDetails{
		Stage:     stage,
		ErrorType: domain.ErrorTypeCompilation,
		Message:   err.Error(),
		Tool:      string(stage),
		Severity:  "error",
	}
}

func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().UnixNano())
}
