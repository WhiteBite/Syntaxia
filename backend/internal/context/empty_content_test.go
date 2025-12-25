package context

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestService_BuildContext_EmptyContentDueToPathTraversal tests the scenario where
// backend returns fileCount > 0 but all files are filtered out due to path traversal,
// resulting in empty content (0 bytes, 0 lines)
func TestService_BuildContext_EmptyContentDueToPathTraversal(t *testing.T) {
	// Setup
	mockFileReader := new(MockFileContentReader)
	mockTokenCounter := new(MockTokenCounter)
	mockLogger := new(MockLogger)
	mockBus := new(MockEventBus)

	// Create temporary directory for context storage
	tempDir, err := os.MkdirTemp("", "context_empty_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create service
	svc, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	assert.NoError(t, err)
	svc.contextDir = tempDir

	// Mock logger to accept any calls
	mockLogger.On("Debug", mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Warning", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()

	// Mock event bus
	mockBus.On("Emit", mock.Anything, mock.Anything).Maybe()

	// Simulate path traversal scenario:
	// - Request 18 files
	// - All files are filtered out by path validation
	// - ReadContents returns empty map (no files)
	projectPath := filepath.Join(tempDir, "project")
	err = os.MkdirAll(projectPath, 0o755)
	assert.NoError(t, err)

	// Files that would trigger path traversal warnings
	requestedFiles := []string{
		"../outside/file1.kt",
		"../outside/file2.kt",
		"../outside/file3.kt",
		"../outside/file4.kt",
		"../outside/file5.kt",
		"../outside/file6.kt",
		"../outside/file7.kt",
		"../outside/file8.kt",
		"../outside/file9.kt",
		"../outside/file10.kt",
		"../outside/file11.kt",
		"../outside/file12.kt",
		"../outside/file13.kt",
		"../outside/file14.kt",
		"../outside/file15.kt",
		"../outside/file16.kt",
		"../outside/file17.kt",
		"../outside/file18.kt",
	}

	// Mock ReadContents to return empty map (all files filtered)
	mockFileReader.On("ReadContents", mock.Anything, mock.Anything, projectPath, mock.Anything).
		Return(map[string]string{}, nil)

	// Build context
	ctx := context.Background()
	result, err := svc.BuildContext(ctx, projectPath, requestedFiles, &BuildOptions{
		MaxTokens:       1000,
		StripComments:   false,
		IncludeManifest: false,
	})

	// Assertions
	assert.NoError(t, err, "BuildContext should not return error")
	assert.NotNil(t, result, "Result should not be nil")

	// Calculate metrics from result
	fileCount := len(result.Files)
	contentSize := len(result.Content)
	lineCount := len(strings.Split(result.Content, "\n"))
	if result.Content == "" {
		lineCount = 0
	}

	// The bug: backend may return Files list with requested files
	// But actual content is empty (0 bytes, 0 lines)
	t.Logf("Result: FileCount=%d, ContentSize=%d, LineCount=%d, TokenCount=%d (requested %d files)",
		fileCount, contentSize, lineCount, result.TokenCount, len(requestedFiles))

	// Current behavior (BUG):
	// The Files array might contain requested files even though content is empty
	// This is the bug we're testing for

	// Expected behavior (FIX):
	// If content is empty, Files should also be empty
	if contentSize == 0 {
		assert.Equal(t, 0, fileCount, "FileCount should be 0 when content is empty")
		assert.Equal(t, 0, lineCount, "LineCount should be 0 when content is empty")
		assert.Equal(t, 0, result.TokenCount, "TokenCount should be 0 when content is empty")
	}

	// Verify the context file was created
	contextPath := filepath.Join(tempDir, result.ID+".ctx")
	assert.FileExists(t, contextPath, "Context file should exist")

	// Verify context file is empty or minimal
	content, err := os.ReadFile(contextPath)
	assert.NoError(t, err)
	if contentSize == 0 {
		assert.LessOrEqual(t, len(content), 10, "Context file should be empty or nearly empty when no content")
	}
}

// TestService_BuildContext_PartialContentDueToPathTraversal tests the scenario where
// some files are valid and some are filtered out
func TestService_BuildContext_PartialContentDueToPathTraversal(t *testing.T) {
	// Setup
	mockFileReader := new(MockFileContentReader)
	mockTokenCounter := new(MockTokenCounter)
	mockLogger := new(MockLogger)
	mockBus := new(MockEventBus)

	tempDir, err := os.MkdirTemp("", "context_partial_test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	svc, err := NewService(mockFileReader, mockTokenCounter, mockBus, mockLogger)
	assert.NoError(t, err)
	svc.contextDir = tempDir

	mockLogger.On("Debug", mock.Anything).Maybe()
	mockLogger.On("Info", mock.Anything).Maybe()
	mockLogger.On("Warning", mock.Anything).Maybe()
	mockLogger.On("Error", mock.Anything).Maybe()
	mockBus.On("Emit", mock.Anything, mock.Anything).Maybe()

	projectPath := filepath.Join(tempDir, "project")
	err = os.MkdirAll(projectPath, 0o755)
	assert.NoError(t, err)

	// Request 10 files, but only 3 are valid
	requestedFiles := []string{
		"valid1.kt",
		"valid2.kt",
		"valid3.kt",
		"../outside/invalid1.kt",
		"../outside/invalid2.kt",
		"../outside/invalid3.kt",
		"../outside/invalid4.kt",
		"../outside/invalid5.kt",
		"../outside/invalid6.kt",
		"../outside/invalid7.kt",
	}

	// Mock ReadContents to return only valid files
	validContent := map[string]string{
		"valid1.kt": "package test\nfun main() {}\n",
		"valid2.kt": "package test\nclass Test {}\n",
		"valid3.kt": "package test\nval x = 1\n",
	}
	mockFileReader.On("ReadContents", mock.Anything, mock.Anything, projectPath, mock.Anything).
		Return(validContent, nil)

	mockTokenCounter.On("CountTokens", mock.Anything).Return(10)

	ctx := context.Background()
	result, err := svc.BuildContext(ctx, projectPath, requestedFiles, &BuildOptions{
		MaxTokens:       1000,
		StripComments:   false,
		IncludeManifest: false,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)

	fileCount := len(result.Files)
	contentSize := len(result.Content)
	lineCount := len(strings.Split(result.Content, "\n"))

	t.Logf("Result: FileCount=%d, ContentSize=%d, LineCount=%d (requested %d files)",
		fileCount, contentSize, lineCount, len(requestedFiles))

	// Expected: FileCount should be 3 (actual files), not 10 (requested)
	assert.Equal(t, 3, fileCount, "FileCount should be 3 (actual files included)")
	assert.Greater(t, contentSize, 0, "ContentSize should be > 0")
	assert.Greater(t, lineCount, 0, "LineCount should be > 0")
}
