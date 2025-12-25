package project

import (
	"context"
	"syntaxia/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	benchProjServicePath = "/test/project"
	branchMain           = "main"
)

// Mock implementations for benchmarking
type mockProjectLogger struct{}

func (m *mockProjectLogger) Debug(msg string)   {}
func (m *mockProjectLogger) Info(msg string)    {}
func (m *mockProjectLogger) Warning(msg string) {}
func (m *mockProjectLogger) Error(msg string)   {}
func (m *mockProjectLogger) Fatal(msg string)   {}

// Mock EventBus for benchmarking
type mockEventBus struct{}

func (m *mockEventBus) Emit(eventName string, data ...interface{}) {}

// Mock TreeBuilder for benchmarking
type mockTreeBuilder struct {
	delayMs int
}

func (m *mockTreeBuilder) BuildTree(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*domain.FileNode, error) {
	// Create a mock file tree structure
	nodes := []*domain.FileNode{
		{
			Name:    "file1.go",
			Path:    benchProjServicePath + "/file1.go",
			RelPath: "file1.go",
			IsDir:   false,
			Size:    1024,
		},
		{
			Name:    "dir1",
			Path:    benchProjServicePath + "/dir1",
			RelPath: "dir1",
			IsDir:   true,
			Children: []*domain.FileNode{
				{
					Name:    "file2.js",
					Path:    benchProjServicePath + "/dir1/file2.js",
					RelPath: "dir1/file2.js",
					IsDir:   false,
					Size:    2048,
				},
			},
		},
	}
	return nodes, nil
}

func (m *mockTreeBuilder) InvalidateCache() {
	// No-op for benchmark mock
}

// Mock GitRepository for benchmarking
type mockGitRepository struct {
	delayMs int
}

func (m *mockGitRepository) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	return []domain.FileStatus{
		{
			Path:   "file1.go",
			Status: "modified",
		},
		{
			Path:   "file2.js",
			Status: "added",
		},
	}, nil
}

func (m *mockGitRepository) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	return []domain.CommitWithFiles{
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
	}, nil
}

func (m *mockGitRepository) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	return "file content", nil
}

func (m *mockGitRepository) GetGitignoreContent(projectRoot string) (string, error) {
	return "node_modules/\n*.log\n", nil
}

func (m *mockGitRepository) IsGitAvailable() bool {
	return true
}

func (m *mockGitRepository) GetBranches(projectRoot string) ([]string, error) {
	return []string{branchMain, "develop", "feature/test"}, nil
}

func (m *mockGitRepository) GetCurrentBranch(projectRoot string) (string, error) {
	return branchMain, nil
}

func (m *mockGitRepository) GetAllFiles(projectPath string) ([]string, error) {
	return []string{"file1.go", "file2.js", "README.md"}, nil
}

func (m *mockGitRepository) GenerateDiff(projectPath string) (string, error) {
	return "diff --git a/file1.go b/file1.go...", nil
}

func (m *mockGitRepository) IsGitRepository(projectPath string) bool {
	return true
}

func (m *mockGitRepository) CloneRepository(url, targetPath string, depth int) error {
	return nil
}

func (m *mockGitRepository) CheckoutBranch(projectPath, branch string) error {
	return nil
}

func (m *mockGitRepository) CheckoutCommit(projectPath, commitHash string) error {
	return nil
}

func (m *mockGitRepository) GetCommitHistory(projectPath string, limit int) ([]domain.CommitInfo, error) {
	return []domain.CommitInfo{
		{Hash: "abc123", Subject: "Initial commit", Author: "test", Date: "2024-01-01"},
	}, nil
}

func (m *mockGitRepository) FetchRemoteBranches(projectPath string) ([]string, error) {
	return []string{"origin/main", "origin/develop"}, nil
}

func (m *mockGitRepository) ListFilesAtRef(projectPath, ref string) ([]string, error) {
	return []string{"file1.go", "file2.js"}, nil
}

func (m *mockGitRepository) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
	return "file content at ref", nil
}

// Mock ContextService for benchmarking
type mockContextService struct {
	delayMs int
}

func (m *mockContextService) GenerateContextAsync(ctx context.Context, rootDir string, includedPaths []string) {
	// Mock async context generation
}

func BenchmarkProjectService_ListFiles(b *testing.B) {
	// Setup
	mockLogger := &mockProjectLogger{}
	mockBus := &mockEventBus{}
	mockTreeBuilder := &mockTreeBuilder{delayMs: 5}
	mockGitRepo := &mockGitRepository{}
	mockContextSvc := &mockContextService{}

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	b.ResetTimer()
	b.ReportAllocs()

	dirPath := benchProjServicePath
	useGitignore := true
	useCustomIgnore := true

	for i := 0; i < b.N; i++ {
		_, err := service.ListFiles(dirPath, useGitignore, useCustomIgnore)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProjectService_GetUncommittedFiles(b *testing.B) {
	// Setup
	mockLogger := &mockProjectLogger{}
	mockBus := &mockEventBus{}
	mockTreeBuilder := &mockTreeBuilder{}
	mockGitRepo := &mockGitRepository{delayMs: 10}
	mockContextSvc := &mockContextService{}

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	b.ResetTimer()
	b.ReportAllocs()

	projectRoot := benchProjServicePath

	for i := 0; i < b.N; i++ {
		_, err := service.GetUncommittedFiles(projectRoot)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProjectService_GetRichCommitHistory(b *testing.B) {
	// Setup
	mockLogger := &mockProjectLogger{}
	mockBus := &mockEventBus{}
	mockTreeBuilder := &mockTreeBuilder{}
	mockGitRepo := &mockGitRepository{delayMs: 15}
	mockContextSvc := &mockContextService{}

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	b.ResetTimer()
	b.ReportAllocs()

	projectRoot := benchProjServicePath
	branchName := branchMain
	limit := 10

	for i := 0; i < b.N; i++ {
		_, err := service.GetRichCommitHistory(projectRoot, branchName, limit)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProjectService_IsGitAvailable(b *testing.B) {
	// Setup
	mockLogger := &mockProjectLogger{}
	mockBus := &mockEventBus{}
	mockTreeBuilder := &mockTreeBuilder{}
	mockGitRepo := &mockGitRepository{}
	mockContextSvc := &mockContextService{}

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		result := service.IsGitAvailable()
		assert.True(b, result)
	}
}

func BenchmarkProjectService_GenerateContext(b *testing.B) {
	// Setup
	mockLogger := &mockProjectLogger{}
	mockBus := &mockEventBus{}
	mockTreeBuilder := &mockTreeBuilder{}
	mockGitRepo := &mockGitRepository{}
	mockContextSvc := &mockContextService{delayMs: 20}

	service := NewService(mockLogger, mockBus, mockTreeBuilder, mockGitRepo, mockContextSvc)

	b.ResetTimer()
	b.ReportAllocs()

	ctx := context.Background()
	rootDir := benchProjServicePath
	includedPaths := []string{"file1.go", "dir1/file2.js"}

	for i := 0; i < b.N; i++ {
		service.GenerateContext(ctx, rootDir, includedPaths)
	}
}
