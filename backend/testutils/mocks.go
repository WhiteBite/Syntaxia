package testutils

import (
	"context"
	"syntaxia/domain"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockLogger is a shared mock implementation of domain.Logger for tests
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Warning(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Fatal(msg string) {
	m.Called(msg)
}

// NewMockLogger creates a new instance of MockLogger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// MockStaticAnalyzerService is a shared mock implementation of domain.IStaticAnalyzerService for tests
type MockStaticAnalyzerService struct {
	mock.Mock
}

func (m *MockStaticAnalyzerService) AnalyzeProject(ctx context.Context, projectPath string, languages []string) (*domain.StaticAnalysisReport, error) {
	args := m.Called(ctx, projectPath, languages)
	return args.Get(0).(*domain.StaticAnalysisReport), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeFile(ctx context.Context, filePath, language string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, filePath, language)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) GetSupportedAnalyzers() []domain.StaticAnalyzerType {
	args := m.Called()
	return args.Get(0).([]domain.StaticAnalyzerType)
}

func (m *MockStaticAnalyzerService) GetAnalyzerForLanguage(language string) (domain.StaticAnalyzer, error) {
	args := m.Called(language)
	return args.Get(0).(domain.StaticAnalyzer), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeGoProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeTypeScriptProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeJavaScriptProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeJavaProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzePythonProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) AnalyzeCppProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).(*domain.StaticAnalysisResult), args.Error(1)
}

func (m *MockStaticAnalyzerService) ValidateAnalysisResults(results map[string]*domain.StaticAnalysisResult) *domain.StaticAnalysisValidationResult {
	args := m.Called(results)
	return args.Get(0).(*domain.StaticAnalysisValidationResult)
}

// MockTestService is a shared mock implementation of domain.ITestService for tests
type MockTestService struct {
	mock.Mock
}

func (m *MockTestService) RunTests(ctx context.Context, config *domain.TestConfig) ([]*domain.TestResult, error) {
	args := m.Called(ctx, config)
	return args.Get(0).([]*domain.TestResult), args.Error(1)
}

func (m *MockTestService) RunTargetedTests(ctx context.Context, config *domain.TestConfig, changedFiles []string) ([]*domain.TestResult, error) {
	args := m.Called(ctx, config, changedFiles)
	return args.Get(0).([]*domain.TestResult), args.Error(1)
}

func (m *MockTestService) DiscoverTests(ctx context.Context, projectPath, language string) (*domain.TestSuite, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).(*domain.TestSuite), args.Error(1)
}

func (m *MockTestService) BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*domain.AffectedGraph, error) {
	args := m.Called(ctx, changedFiles, projectPath)
	return args.Get(0).(*domain.AffectedGraph), args.Error(1)
}

func (m *MockTestService) GetTestCoverage(ctx context.Context, testPath string) (*domain.TestCoverage, error) {
	args := m.Called(ctx, testPath)
	return args.Get(0).(*domain.TestCoverage), args.Error(1)
}

func (m *MockTestService) GetSupportedLanguages() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockTestService) RunSmokeTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).([]*domain.TestResult), args.Error(1)
}

func (m *MockTestService) RunUnitTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).([]*domain.TestResult), args.Error(1)
}

func (m *MockTestService) RunIntegrationTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).([]*domain.TestResult), args.Error(1)
}

func (m *MockTestService) ValidateTestResults(results []*domain.TestResult) *domain.TestValidationResult {
	args := m.Called(results)
	return args.Get(0).(*domain.TestValidationResult)
}

// MockBuildService is a shared mock implementation of domain.IBuildService for tests
type MockBuildService struct {
	mock.Mock
}

func (m *MockBuildService) Build(ctx context.Context, projectPath, language string) (*domain.BuildResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).(*domain.BuildResult), args.Error(1)
}

func (m *MockBuildService) TypeCheck(ctx context.Context, projectPath, language string) (*domain.TypeCheckResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).(*domain.TypeCheckResult), args.Error(1)
}

func (m *MockBuildService) BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*domain.BuildResult, *domain.TypeCheckResult, error) {
	args := m.Called(ctx, projectPath, language)
	return args.Get(0).(*domain.BuildResult), args.Get(1).(*domain.TypeCheckResult), args.Error(2)
}

func (m *MockBuildService) BuildMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*domain.BuildResult, error) {
	args := m.Called(ctx, projectPath, languages)
	return args.Get(0).(map[string]*domain.BuildResult), args.Error(1)
}

func (m *MockBuildService) TypeCheckMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*domain.TypeCheckResult, error) {
	args := m.Called(ctx, projectPath, languages)
	return args.Get(0).(map[string]*domain.TypeCheckResult), args.Error(1)
}

func (m *MockBuildService) ValidateProject(ctx context.Context, projectPath string, languages []string) (*domain.ProjectValidationResult, error) {
	args := m.Called(ctx, projectPath, languages)
	return args.Get(0).(*domain.ProjectValidationResult), args.Error(1)
}

func (m *MockBuildService) GetSupportedLanguages() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func (m *MockBuildService) DetectLanguages(ctx context.Context, projectPath string) ([]string, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).([]string), args.Error(1)
}

// MockGuardrailService is a shared mock implementation of domain.GuardrailService for tests
type MockGuardrailService struct {
	mock.Mock
}

func (m *MockGuardrailService) ValidateTask(taskID string, files []string, linesChanged int64) (*domain.TaskValidationResult, error) {
	args := m.Called(taskID, files, linesChanged)
	return args.Get(0).(*domain.TaskValidationResult), args.Error(1)
}

func (m *MockGuardrailService) ValidatePath(path string) ([]domain.GuardrailViolation, error) {
	args := m.Called(path)
	return args.Get(0).([]domain.GuardrailViolation), args.Error(1)
}

func (m *MockGuardrailService) ValidateBudget(budgetType domain.BudgetType, current int64) ([]domain.BudgetViolation, error) {
	args := m.Called(budgetType, current)
	return args.Get(0).([]domain.BudgetViolation), args.Error(1)
}

func (m *MockGuardrailService) EnableEphemeralMode(taskID, taskType string, duration time.Duration) error {
	args := m.Called(taskID, taskType, duration)
	return args.Error(0)
}

func (m *MockGuardrailService) DisableEphemeralMode() {
	m.Called()
}

func (m *MockGuardrailService) GetPolicies() ([]domain.GuardrailPolicy, error) {
	args := m.Called()
	return args.Get(0).([]domain.GuardrailPolicy), args.Error(1)
}

func (m *MockGuardrailService) GetBudgetPolicies() ([]domain.BudgetPolicy, error) {
	args := m.Called()
	return args.Get(0).([]domain.BudgetPolicy), args.Error(1)
}

func (m *MockGuardrailService) AddPolicy(policy domain.GuardrailPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockGuardrailService) RemovePolicy(policyID string) error {
	args := m.Called(policyID)
	return args.Error(0)
}

func (m *MockGuardrailService) UpdatePolicy(policy domain.GuardrailPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockGuardrailService) AddBudgetPolicy(policy domain.BudgetPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockGuardrailService) RemoveBudgetPolicy(policyID string) error {
	args := m.Called(policyID)
	return args.Error(0)
}

func (m *MockGuardrailService) UpdateBudgetPolicy(policy domain.BudgetPolicy) error {
	args := m.Called(policy)
	return args.Error(0)
}

func (m *MockGuardrailService) GetConfig() domain.GuardrailConfig {
	args := m.Called()
	return args.Get(0).(domain.GuardrailConfig)
}

func (m *MockGuardrailService) UpdateConfig(config domain.GuardrailConfig) error {
	args := m.Called(config)
	return args.Error(0)
}

func (m *MockGuardrailService) SetTaskTypeProvider(taskTypeProvider domain.TaskTypeProvider) {
	m.Called(taskTypeProvider)
}

// MockTaskflowService is a shared mock implementation of domain.TaskflowService for tests
type MockTaskflowService struct {
	mock.Mock
}

func (m *MockTaskflowService) LoadTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskflowService) GetTaskStatus(taskID string) (*domain.TaskStatus, error) {
	args := m.Called(taskID)
	return args.Get(0).(*domain.TaskStatus), args.Error(1)
}

func (m *MockTaskflowService) UpdateTaskStatus(taskID string, state domain.TaskState, message string) error {
	args := m.Called(taskID, state, message)
	return args.Error(0)
}

func (m *MockTaskflowService) ExecuteTask(taskID string) error {
	args := m.Called(taskID)
	return args.Error(0)
}

func (m *MockTaskflowService) ExecuteTaskflow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTaskflowService) GetReadyTasks() ([]domain.Task, error) {
	args := m.Called()
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskflowService) GetTaskDependencies(taskID string) ([]domain.Task, error) {
	args := m.Called(taskID)
	return args.Get(0).([]domain.Task), args.Error(1)
}

func (m *MockTaskflowService) ValidateTaskflow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTaskflowService) GetTaskflowProgress() (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTaskflowService) ResetTaskflow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTaskflowService) StartAutonomousTask(ctx context.Context, request domain.AutonomousTaskRequest) (*domain.AutonomousTaskResponse, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(*domain.AutonomousTaskResponse), args.Error(1)
}

func (m *MockTaskflowService) CancelAutonomousTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskflowService) GetAutonomousTaskStatus(ctx context.Context, taskID string) (*domain.AutonomousTaskStatus, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(*domain.AutonomousTaskStatus), args.Error(1)
}

func (m *MockTaskflowService) ListAutonomousTasks(ctx context.Context, projectPath string) ([]domain.AutonomousTask, error) {
	args := m.Called(ctx, projectPath)
	return args.Get(0).([]domain.AutonomousTask), args.Error(1)
}

func (m *MockTaskflowService) GetTaskLogs(ctx context.Context, taskID string) ([]domain.LogEntry, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).([]domain.LogEntry), args.Error(1)
}

func (m *MockTaskflowService) PauseTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

func (m *MockTaskflowService) ResumeTask(ctx context.Context, taskID string) error {
	args := m.Called(ctx, taskID)
	return args.Error(0)
}

// MockApplyEngine is a shared mock implementation of domain.ApplyEngine for tests
type MockApplyEngine struct {
	mock.Mock
}

func (m *MockApplyEngine) ApplyOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	args := m.Called(ctx, op)
	return args.Get(0).(*domain.ApplyResult), args.Error(1)
}

func (m *MockApplyEngine) ApplyOperations(ctx context.Context, ops []*domain.ApplyOperation) ([]*domain.ApplyResult, error) {
	args := m.Called(ctx, ops)
	return args.Get(0).([]*domain.ApplyResult), args.Error(1)
}

func (m *MockApplyEngine) ApplyEdit(ctx context.Context, edit domain.Edit) error {
	args := m.Called(ctx, edit)
	return args.Error(0)
}

func (m *MockApplyEngine) ValidateOperation(ctx context.Context, op *domain.ApplyOperation) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *MockApplyEngine) RollbackOperation(ctx context.Context, result *domain.ApplyResult) error {
	args := m.Called(ctx, result)
	return args.Error(0)
}

func (m *MockApplyEngine) RegisterFormatter(language string, formatter domain.Formatter) {
	m.Called(language, formatter)
}

func (m *MockApplyEngine) RegisterImportFixer(language string, fixer domain.ImportFixer) {
	m.Called(language, fixer)
}

// MockFormatter is a shared mock implementation of domain.Formatter for tests
type MockFormatter struct {
	mock.Mock
}

func (m *MockFormatter) FormatFile(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *MockFormatter) FormatContent(ctx context.Context, content, language string) (string, error) {
	args := m.Called(ctx, content, language)
	return args.Get(0).(string), args.Error(1)
}

func (m *MockFormatter) GetSupportedLanguages() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

// MockVerificationPipelineService is a shared mock implementation of application.VerificationPipelineService for tests
type MockVerificationPipelineService struct {
	mock.Mock
}

// RunTaskProtocolVerification implements the VerificationPipelineService interface
func (m *MockVerificationPipelineService) RunTaskProtocolVerification(ctx context.Context, config *domain.TaskProtocolConfig) (*domain.TaskProtocolResult, error) {
	args := m.Called(ctx, config)
	return args.Get(0).(*domain.TaskProtocolResult), args.Error(1)
}

// CreateTaskProtocolConfig implements the VerificationPipelineService interface
func (m *MockVerificationPipelineService) CreateTaskProtocolConfig(verifyConfig *domain.VerificationConfig) *domain.TaskProtocolConfig {
	args := m.Called(verifyConfig)
	return args.Get(0).(*domain.TaskProtocolConfig)
}

// MockIntelligentAIService is a shared mock implementation of application.IntelligentAIService for tests
type MockIntelligentAIService struct {
	mock.Mock
}

// MockCommandRunner is a shared mock implementation of domain.CommandRunner for tests
type MockCommandRunner struct {
	mock.Mock
}

func (m *MockCommandRunner) RunCommand(ctx context.Context, name string, args ...string) ([]byte, error) {
	callArgs := append([]interface{}{ctx, name}, args)
	argsCalled := m.Called(callArgs...)
	return argsCalled.Get(0).([]byte), argsCalled.Error(1)
}

func (m *MockCommandRunner) RunCommandInDir(ctx context.Context, dir, name string, args ...string) ([]byte, error) {
	callArgs := append([]interface{}{ctx, dir, name}, args)
	argsCalled := m.Called(callArgs...)
	return argsCalled.Get(0).([]byte), argsCalled.Error(1)
}
