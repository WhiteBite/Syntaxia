package build

import (
	"context"
	"fmt"

	"syntaxia/domain"
)

// TestService предоставляет высокоуровневый API для работы с тестами
type TestService struct {
	log        domain.Logger
	testEngine domain.TestEngine
}

// NewTestService создает новый сервис тестирования
func NewTestService(log domain.Logger, testEngine domain.TestEngine) *TestService {
	return &TestService{
		log:        log,
		testEngine: testEngine,
	}
}

// RunTests выполняет тесты согласно конфигурации
func (s *TestService) RunTests(ctx context.Context, config *domain.TestConfig) ([]*domain.TestResult, error) {
	s.log.Info(fmt.Sprintf("Running tests with scope: %s", config.Scope))

	if config.Scope == domain.TestScopeAffected || config.Scope == domain.TestScopeAffectedSmoke {
		changedFiles := []string{}
		affectedGraph, err := s.testEngine.BuildAffectedGraph(ctx, changedFiles, config.ProjectPath)
		if err != nil {
			s.log.Warning(fmt.Sprintf("Failed to build affected graph: %v", err))
			return s.testEngine.RunTests(ctx, config)
		}
		return s.testEngine.RunTargetedTests(ctx, config, affectedGraph)
	}
	return s.testEngine.RunTests(ctx, config)
}

// RunTargetedTests выполняет целевые тесты для затронутых файлов
func (s *TestService) RunTargetedTests(ctx context.Context, config *domain.TestConfig, changedFiles []string) ([]*domain.TestResult, error) {
	s.log.Info(fmt.Sprintf("Running targeted tests for %d changed files", len(changedFiles)))

	affectedGraph, err := s.testEngine.BuildAffectedGraph(ctx, changedFiles, config.ProjectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build affected graph: %w", err)
	}
	return s.testEngine.RunTargetedTests(ctx, config, affectedGraph)
}

// DiscoverTests обнаруживает тесты в проекте
func (s *TestService) DiscoverTests(ctx context.Context, projectPath, language string) (*domain.TestSuite, error) {
	return s.testEngine.DiscoverTests(ctx, projectPath, language)
}

// BuildAffectedGraph строит граф затронутых файлов
func (s *TestService) BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*domain.AffectedGraph, error) {
	return s.testEngine.BuildAffectedGraph(ctx, changedFiles, projectPath)
}

// GetTestCoverage получает покрытие тестами
func (s *TestService) GetTestCoverage(ctx context.Context, testPath string) (*domain.TestCoverage, error) {
	return s.testEngine.GetTestCoverage(ctx, testPath)
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (s *TestService) GetSupportedLanguages() []string {
	return s.testEngine.GetSupportedLanguages()
}

// RunSmokeTests выполняет только smoke тесты
func (s *TestService) RunSmokeTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	config := &domain.TestConfig{
		Language:    language,
		ProjectPath: projectPath,
		Scope:       domain.TestScopeSmoke,
		Timeout:     60,
		Verbose:     true,
	}
	return s.RunTests(ctx, config)
}

// RunUnitTests выполняет только unit тесты
func (s *TestService) RunUnitTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	config := &domain.TestConfig{
		Language:    language,
		ProjectPath: projectPath,
		Scope:       domain.TestScopeUnit,
		Timeout:     30,
		Coverage:    true,
	}
	return s.RunTests(ctx, config)
}

// RunIntegrationTests выполняет только integration тесты
func (s *TestService) RunIntegrationTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	config := &domain.TestConfig{
		Language:    language,
		ProjectPath: projectPath,
		Scope:       domain.TestScopeIntegration,
		Timeout:     300,
		Verbose:     true,
	}
	return s.RunTests(ctx, config)
}

// ValidateTestResults валидирует результаты тестов
func (s *TestService) ValidateTestResults(results []*domain.TestResult) *domain.TestValidationResult {
	validation := &domain.TestValidationResult{
		TotalTests:      len(results),
		FailedTestPaths: []string{},
	}

	for _, result := range results {
		validation.TotalDuration += result.Duration
		if result.Success {
			validation.PassedTests++
		} else {
			validation.FailedTests++
			validation.FailedTestPaths = append(validation.FailedTestPaths, result.TestPath)
		}
	}

	validation.Success = validation.FailedTests == 0
	if validation.TotalTests > 0 {
		validation.SuccessRate = float64(validation.PassedTests) / float64(validation.TotalTests) * 100
	}
	return validation
}
