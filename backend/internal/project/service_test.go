package project

import (
	"context"
	"errors"
	"syntaxia/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testProjectPathProject = "/test/project"

// Mock implementations for testing
type MockProjectLogger struct {
	mock.Mock
}

func (m *MockProjectLogger) Debug(msg string) {
	m.Called(msg)
}

func (m *MockProjectLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockProjectLogger) Warning(msg string) {
	m.Called(msg)
}

func (m *MockProjectLogger) Error(msg string) {
	m.Called(msg)
}

func (m *MockProjectLogger) Fatal(msg string) {
	m.Called(msg)
}

// Mock EventBus for testing
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Emit(eventName string, data ...interface{}) {
	m.Called(eventName, data)
}

// Mock TreeBuilder for testing
type MockTreeBuilder struct {
	mock.Mock
}

func (m *MockTreeBuilder) BuildTree(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*domain.FileNode, error) {
	args := m.Called(dirPath, useGitignore, useCustomIgnore)
	return args.Get(0).([]*domain.FileNode), args.Error(1)
}

func (m *MockTreeBuilder) InvalidateCache() {
	m.Called()
}

// Mock GitRepository for testing
type MockGitRepository struct {
	mock.Mock
}

func (m *MockGitRepository) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	args := m.Called(projectRoot)
	return args.Get(0).([]domain.FileStatus), args.Error(1)
}

func (m *MockGitRepository) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	args := m.Called(projectRoot, branchName, limit)
	return args.Get(0).([]domain.CommitWithFiles), args.Error(1)
}

func (m *MockGitRepository) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	args := m.Called(projectRoot, filePath, commitHash)
	return args.String(0), args.Error(1)
}

func (m *MockGitRepository) GetGitignoreContent(projectRoot string) (string, error) {
	args := m.Called(projectRoot)
	return args.String(0), args.Error(1)
}

func (m *MockGitRepository) IsGitAvailable() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockGitRepository) GetBranches(projectRoot string) ([]string, error) {
	args := m.Called(projectRoot)
	if branches := args.Get(0); branches != nil {
		return branches.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGitRepository) GetCurrentBranch(projectRoot string) (string, error) {
	args := m.Called(projectRoot)
	return args.String(0), args.Error(1)
}

func (m *MockGitRepository) GetAllFiles(projectPath string) ([]string, error) {
	args := m.Called(projectPath)
	if files := args.Get(0); files != nil {
		return files.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGitRepository) GenerateDiff(projectPath string) (string, error) {
	args := m.Called(projectPath)
	return args.String(0), args.Error(1)
}

func (m *MockGitRepository) IsGitRepository(projectPath string) bool {
	args := m.Called(projectPath)
	return args.Bool(0)
}

func (m *MockGitRepository) CloneRepository(url, targetPath string, depth int) error {
	args := m.Called(url, targetPath, depth)
	return args.Error(0)
}

func (m *MockGitRepository) CheckoutBranch(projectPath, branch string) error {
	args := m.Called(projectPath, branch)
	return args.Error(0)
}

func (m *MockGitRepository) CheckoutCommit(projectPath, commitHash string) error {
	args := m.Called(projectPath, commitHash)
	return args.Error(0)
}

func (m *MockGitRepository) GetCommitHistory(projectPath string, limit int) ([]domain.CommitInfo, error) {
	args := m.Called(projectPath, limit)
	if commits := args.Get(0); commits != nil {
		return commits.([]domain.CommitInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGitRepository) FetchRemoteBranches(projectPath string) ([]string, error) {
	args := m.Called(projectPath)
	if branches := args.Get(0); branches != nil {
		return branches.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGitRepository) ListFilesAtRef(projectPath, ref string) ([]string, error) {
	args := m.Called(projectPath, ref)
	if files := args.Get(0); files != nil {
		return files.([]string), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockGitRepository) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
	args := m.Called(projectPath, filePath, ref)
	return args.String(0), args.Error(1)
}

// Mock ContextService for testing
type MockContextService struct {
	mock.Mock
}

func (m *MockContextService) GenerateContextAsync(ctx context.Context, rootDir string, includedPaths []string) {
	m.Called(ctx, rootDir, includedPaths)
}

func TestProjectService_NewService(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	// Execute
	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, mockLogger, service.log)
	assert.Equal(t, mockBus, service.bus)
	assert.Equal(t, mockTreeBuilder, service.treeBuilder)
	assert.Equal(t, mockGitRepo, service.gitRepo)
	assert.Equal(t, mockContextSvc, service.contextSvc)
}

func TestProjectService_ListFiles_Success(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	dirPath := testProjectPathProject
	useGitignore := true
	useCustomIgnore := true

	fileNodes := []*domain.FileNode{
		{
			Name:  "file1.go",
			Path:  "/test/project/file1.go",
			IsDir: false,
		},
		{
			Name:  "dir1",
			Path:  "/test/project/dir1",
			IsDir: true,
		},
	}

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockTreeBuilder.On("BuildTree", dirPath, useGitignore, useCustomIgnore).Return(fileNodes, nil)

	// Execute
	nodes, err := service.ListFiles(dirPath, useGitignore, useCustomIgnore)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	assert.Equal(t, fileNodes, nodes)
	assert.Equal(t, 2, len(nodes))

	mockLogger.AssertExpectations(t)
	mockTreeBuilder.AssertExpectations(t)
}

func TestProjectService_ListFiles_Error(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	dirPath := testProjectPathProject
	useGitignore := true
	useCustomIgnore := true

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockLogger.On("Error", mock.AnythingOfType("string")).Return()
	mockTreeBuilder.On("BuildTree", dirPath, useGitignore, useCustomIgnore).Return(([]*domain.FileNode)(nil), errors.New("build tree failed"))

	// Execute
	nodes, err := service.ListFiles(dirPath, useGitignore, useCustomIgnore)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, nodes)
	assert.Contains(t, err.Error(), "build tree failed")

	mockLogger.AssertExpectations(t)
	mockTreeBuilder.AssertExpectations(t)
}

func TestProjectService_GetUncommittedFiles_Success(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	projectRoot := testProjectPathProject

	fileStatuses := []domain.FileStatus{
		{
			Path:   "file1.go",
			Status: "modified",
		},
		{
			Path:   "file2.js",
			Status: "added",
		},
	}

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockGitRepo.On("GetUncommittedFiles", projectRoot).Return(fileStatuses, nil)

	// Execute
	files, err := service.GetUncommittedFiles(projectRoot)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, files)
	assert.Equal(t, fileStatuses, files)
	assert.Equal(t, 2, len(files))

	mockLogger.AssertExpectations(t)
	mockGitRepo.AssertExpectations(t)
}

func TestProjectService_GetUncommittedFiles_Error(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	projectRoot := testProjectPathProject

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockLogger.On("Error", mock.AnythingOfType("string")).Return()
	mockGitRepo.On("GetUncommittedFiles", projectRoot).Return(([]domain.FileStatus)(nil), errors.New("git operation failed"))

	// Execute
	files, err := service.GetUncommittedFiles(projectRoot)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, files)
	assert.Contains(t, err.Error(), "git operation failed")

	mockLogger.AssertExpectations(t)
	mockGitRepo.AssertExpectations(t)
}

func TestProjectService_GetRichCommitHistory_Success(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	projectRoot := testProjectPathProject
	branchName := "main"
	limit := 10

	commits := []domain.CommitWithFiles{
		{
			Hash:    "abc123",
			Subject: "Initial commit",
			Files:   []string{"README.md"},
		},
		{
			Hash:    "def456",
			Subject: "Add feature",
			Files:   []string{"feature.go", "test.go"},
		},
	}

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockGitRepo.On("GetRichCommitHistory", projectRoot, branchName, limit).Return(commits, nil)

	// Execute
	history, err := service.GetRichCommitHistory(projectRoot, branchName, limit)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, history)
	assert.Equal(t, commits, history)
	assert.Equal(t, 2, len(history))

	mockLogger.AssertExpectations(t)
	mockGitRepo.AssertExpectations(t)
}

func TestProjectService_GetRichCommitHistory_Error(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	projectRoot := testProjectPathProject
	branchName := "main"
	limit := 10

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockLogger.On("Error", mock.AnythingOfType("string")).Return()
	mockGitRepo.On("GetRichCommitHistory", projectRoot, branchName, limit).Return(([]domain.CommitWithFiles)(nil), errors.New("git operation failed"))

	// Execute
	history, err := service.GetRichCommitHistory(projectRoot, branchName, limit)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, history)
	assert.Contains(t, err.Error(), "git operation failed")

	mockLogger.AssertExpectations(t)
	mockGitRepo.AssertExpectations(t)
}

func TestProjectService_IsGitAvailable(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test cases
	testCases := []struct {
		name      string
		available bool
		expected  bool
	}{
		{"Git available", true, true},
		{"Git not available", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mocks
			mockGitRepo.On("IsGitAvailable").Return(tc.available).Once()

			// Execute
			result := service.IsGitAvailable()

			// Assert
			assert.Equal(t, tc.expected, result)

			mockGitRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_GenerateContext(t *testing.T) {
	// Setup
	mockLogger := new(MockProjectLogger)
	mockBus := new(MockEventBus)
	mockTreeBuilder := new(MockTreeBuilder)
	mockGitRepo := new(MockGitRepository)
	mockContextSvc := new(MockContextService)

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	// Test data
	ctx := context.Background()
	rootDir := testProjectPathProject
	includedPaths := []string{"file1.go", "dir1/file2.js"}

	// Setup mocks
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockContextSvc.On("GenerateContextAsync", ctx, rootDir, includedPaths).Return()

	// Execute
	service.GenerateContext(ctx, rootDir, includedPaths)

	// Assert
	mockLogger.AssertExpectations(t)
	mockContextSvc.AssertExpectations(t)
}
