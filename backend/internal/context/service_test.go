package context

import (
	"context"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const testProjectPathService = "/test/project"

// MockFileContentReader is a mock implementation of domain.FileContentReader
type MockFileContentReader struct {
	mock.Mock
}

func (m *MockFileContentReader) ReadContents(ctx context.Context, filePaths []string, rootDir string, progress func(current, total int64)) (map[string]string, error) {
	args := m.Called(ctx, filePaths, rootDir, progress)
	return args.Get(0).(map[string]string), args.Error(1)
}

// MockTokenCounter is a mock implementation of TokenCounter
type MockTokenCounter struct {
	mock.Mock
}

func (m *MockTokenCounter) CountTokens(text string) int {
	args := m.Called(text)
	return args.Int(0)
}

// MockLogger is a mock implementation of domain.Logger
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

// MockEventBus is a mock implementation of domain.EventBus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Emit(event string, data ...interface{}) {
	m.Called(event, data)
}

func TestService_BuildContext(t *testing.T) {
	// Setup
	mockFileReader := new(MockFileContentReader)
	mockTokenCounter := new(MockTokenCounter)
	mockLogger := new(MockLogger)
	mockBus := new(MockEventBus)

	// Create temporary directory for context storage
	tempDir, err := os.MkdirTemp("", "context_service_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := &Service{
		fileReader:   mockFileReader,
		tokenCounter: mockTokenCounter,
		eventBus:     mockBus,
		logger:       mockLogger,
		contextDir:   tempDir,
		streams:      make(map[string]*Stream),
	}

	// Test data
	projectPath := testProjectPathService
	includedPaths := []string{"src/main.go", "src/util.go"}
	fileContents := map[string]string{
		"src/main.go": "package main\n\nfunc main() {\n\tprintln(\"Hello World\")\n}",
		"src/util.go": "package main\n\nfunc util() string {\n\treturn \"utility\"\n}",
	}

	// Setup mocks
	mockFileReader.On("ReadContents", mock.Anything, includedPaths, projectPath, mock.AnythingOfType("func(int64, int64)")).Return(fileContents, nil)
	mockTokenCounter.On("CountTokens", mock.AnythingOfType("string")).Return(150)
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockLogger.On("Warning", mock.AnythingOfType("string")).Return().Maybe()
	mockBus.On("Emit", mock.AnythingOfType("string"), mock.Anything).Return()

	// Execute
	ctx := context.Background()
	options := &BuildOptions{
		IncludeManifest: true,
		StripComments:   false,
		MaxTokens:       1000,
	}

	result, err := service.BuildContext(ctx, projectPath, includedPaths, options)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, projectPath, result.ProjectPath)
	assert.Equal(t, includedPaths, result.Files)
	assert.Equal(t, 300, result.TokenCount)                  // Token count for 2 files, 150 tokens each
	assert.Contains(t, result.Content, "STREAMING_CONTEXT:") // Now expects streaming reference

	// Verify stream was created (streaming context file)
	contextFile := filepath.Join(tempDir, result.ID+".ctx")
	assert.FileExists(t, contextFile)

	mockFileReader.AssertExpectations(t)
	mockTokenCounter.AssertExpectations(t)
	mockLogger.AssertExpectations(t)
	// Skipping event bus assertions for now
}

func TestService_BuildContext_TokenLimitExceeded(t *testing.T) {
	// Setup
	mockFileReader := new(MockFileContentReader)
	mockTokenCounter := new(MockTokenCounter)
	mockLogger := new(MockLogger)
	mockBus := new(MockEventBus)

	tempDir, err := os.MkdirTemp("", "context_service_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := &Service{
		fileReader:   mockFileReader,
		tokenCounter: mockTokenCounter,
		eventBus:     mockBus,
		logger:       mockLogger,
		contextDir:   tempDir,
		streams:      make(map[string]*Stream),
	}

	// Test data
	projectPath := testProjectPathService
	includedPaths := []string{"src/large.go"}
	fileContents := map[string]string{
		"src/large.go": strings.Repeat("// Large file content\n", 1000),
	}

	// Setup mocks - return token count that exceeds limit
	mockFileReader.On("ReadContents", mock.Anything, includedPaths, projectPath, mock.AnythingOfType("func(int64, int64)")).Return(fileContents, nil)
	mockTokenCounter.On("CountTokens", mock.AnythingOfType("string")).Return(2000)
	mockLogger.On("Info", mock.AnythingOfType("string")).Return()
	mockLogger.On("Warning", mock.AnythingOfType("string")).Return().Maybe()
	mockBus.On("Emit", mock.AnythingOfType("string"), mock.Anything).Return()

	// Execute
	ctx := context.Background()
	options := &BuildOptions{
		MaxTokens: 1000, // Set limit lower than returned count
	}

	result, err := service.BuildContext(ctx, projectPath, includedPaths, options)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "context would exceed token limit")

	mockFileReader.AssertExpectations(t)
	mockTokenCounter.AssertExpectations(t)
	// Skipping event bus assertions for now
}

func TestService_GetContext(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_service_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := &Service{
		contextDir: tempDir,
		streams:    make(map[string]*Stream),
	}

	// Create a test context
	testContext := &domain.Context{
		ID:          "test-context-id",
		Name:        "Test Context",
		Description: "Test context description",
		Content:     "Test content",
		Files:       []string{"test.go"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectPath: testProjectPathService,
		TokenCount:  100,
	}

	// Save the context manually
	err = service.saveContext(testContext)
	assert.NoError(t, err)

	// Execute
	ctx := context.Background()
	result, err := service.GetContext(ctx, "test-context-id")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testContext.ID, result.ID)
	assert.Equal(t, testContext.Name, result.Name)
	assert.Equal(t, testContext.Content, result.Content)
}

func TestService_GetContext_NotFound(t *testing.T) {
	// Setup
	tempDir, err := os.MkdirTemp("", "context_service_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	service := &Service{
		contextDir: tempDir,
		streams:    make(map[string]*Stream),
	}

	// Execute
	ctx := context.Background()
	result, err := service.GetContext(ctx, "non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")
}
