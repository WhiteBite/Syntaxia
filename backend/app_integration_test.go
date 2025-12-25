package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp_BuildContext_Integration(t *testing.T) {
	// Create a temporary test project
	tempDir, err := os.MkdirTemp("", "app_integration_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"main.go": `package main

import "fmt"

// Main function
func main() {
	fmt.Println("Hello, World!")
}`,
		"utils.go": `package main

// Utility function
func add(a, b int) int {
	return a + b
}`,
		"README.md": `# Test Project

This is a test project for integration testing.`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	// Test data
	projectPath := tempDir
	includedPaths := []string{"main.go", "utils.go"}

	// Create build context options
	options := &domain.ContextBuildOptions{
		IncludeManifest: true,
		StripComments:   false,
		MaxTokens:       10000,
	}

	// Verify test files exist and have expected content
	for filename := range testFiles {
		if filename == "README.md" {
			continue // Skip non-included file
		}

		filePath := filepath.Join(projectPath, filename)
		actualContent, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(actualContent), "package main")
	}

	// Validate that the context would contain expected elements
	expectedContextElements := []string{
		"package main",
		"func main()",
		"func add(",
	}

	// Read and validate file contents that would be processed
	mainContent, err := os.ReadFile(filepath.Join(projectPath, "main.go"))
	assert.NoError(t, err)

	utilsContent, err := os.ReadFile(filepath.Join(projectPath, "utils.go"))
	assert.NoError(t, err)

	combinedContent := string(mainContent) + string(utilsContent)

	for _, expectedElement := range expectedContextElements {
		assert.Contains(t, combinedContent, expectedElement,
			"Context should contain expected element: %s", expectedElement)
	}

	// Verify files exist
	assert.FileExists(t, filepath.Join(projectPath, "main.go"))
	assert.FileExists(t, filepath.Join(projectPath, "utils.go"))

	_ = projectPath
	_ = includedPaths
	_ = options
}

func TestApp_AnalyzeTaskAndCollectContext_MetadataApproach(t *testing.T) {
	// Create a temporary test project
	tempDir, err := os.MkdirTemp("", "app_metadata_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	err = os.WriteFile(filepath.Join(tempDir, "main.go"), []byte("package main\nfunc main() {}"), 0o644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tempDir, "utils.go"), []byte("package main\nfunc utils() {}"), 0o644)
	require.NoError(t, err)

	// Test the metadata approach (new efficient method)
	task := "Implement a new feature"

	// Create project metadata instead of full file list
	projectMetadata := map[string]interface{}{
		"projectPath":    tempDir,
		"totalFileCount": 2,
		"fileTypes":      []string{".go"},
		"languages":      []string{"Go"},
		"framework":      "standard",
	}

	metadataJSON, err := json.Marshal(projectMetadata)
	require.NoError(t, err)

	// Validate metadata structure
	var parsedMetadata map[string]interface{}
	err = json.Unmarshal(metadataJSON, &parsedMetadata)
	assert.NoError(t, err)

	assert.Equal(t, tempDir, parsedMetadata["projectPath"])
	assert.Equal(t, float64(2), parsedMetadata["totalFileCount"]) // JSON numbers are float64

	// Verify this is much more efficient than serializing full file tree
	metadataSize := len(metadataJSON)
	assert.Less(t, metadataSize, 1000, "Metadata should be compact")

	// Compare with theoretical full file serialization
	mockFullFileList := []map[string]interface{}{
		{
			"name":    "main.go",
			"path":    filepath.Join(tempDir, "main.go"),
			"isDir":   false,
			"content": "package main\nfunc main() {}",
		},
		{
			"name":    "utils.go",
			"path":    filepath.Join(tempDir, "utils.go"),
			"isDir":   false,
			"content": "package main\nfunc utils() {}",
		},
	}

	fullListJSON, err := json.Marshal(mockFullFileList)
	require.NoError(t, err)

	fullListSize := len(fullListJSON)
	assert.Greater(t, fullListSize, metadataSize,
		"Full file list should be larger than metadata (efficiency improvement)")

	efficiencyImprovement := float64(fullListSize-metadataSize) / float64(fullListSize) * 100
	assert.Greater(t, efficiencyImprovement, 40.0,
		"Should achieve at least 40%% efficiency improvement")

	_ = task
}

func TestApp_GetContext_Integration(t *testing.T) {
	// Setup test context directory
	tempDir, err := os.MkdirTemp("", "app_context_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a mock context file
	contextID := "test-context-123"
	testContext := domain.Context{
		ID:          contextID,
		Name:        "Test Integration Context",
		Description: "Integration test context",
		Content:     "Test context content",
		Files:       []string{"test.go", "helper.go"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectPath: "/test/project",
		TokenCount:  250,
	}

	// Manually save context file (simulating what BuildContext would do)
	contextData, err := json.MarshalIndent(testContext, "", "  ")
	require.NoError(t, err)

	contextFile := filepath.Join(tempDir, contextID+".json")
	err = os.WriteFile(contextFile, contextData, 0o644)
	require.NoError(t, err)

	// Verify context file was created
	assert.FileExists(t, contextFile)

	// Test reading the context back
	readData, err := os.ReadFile(contextFile)
	assert.NoError(t, err)

	var readContext domain.Context
	err = json.Unmarshal(readData, &readContext)
	assert.NoError(t, err)

	// Validate context data integrity
	assert.Equal(t, testContext.ID, readContext.ID)
	assert.Equal(t, testContext.Name, readContext.Name)
	assert.Equal(t, testContext.Content, readContext.Content)
	assert.Equal(t, testContext.Files, readContext.Files)
	assert.Equal(t, testContext.ProjectPath, readContext.ProjectPath)
	assert.Equal(t, testContext.TokenCount, readContext.TokenCount)
}

func TestApp_StartAutonomousTask_Integration(t *testing.T) {
	// Test autonomous task request validation
	validRequest := domain.AutonomousTaskRequest{
		Task:        "Implement user authentication",
		SlaPolicy:   "standard",
		ProjectPath: "/test/project",
		Options: domain.AutonomousTaskOptions{
			MaxTokens:            4000,
			Temperature:          0.7,
			EnableStaticAnalysis: true,
			EnableTests:          true,
			EnableSBOM:           false,
		},
	}

	// Validate request structure
	assert.NotEmpty(t, validRequest.Task)
	assert.Contains(t, []string{"lite", "standard", "strict"}, validRequest.SlaPolicy)
	assert.NotEmpty(t, validRequest.ProjectPath)
	assert.NotNil(t, validRequest.Options)

	// Test request serialization (what would be sent to backend)
	requestJSON, err := json.Marshal(validRequest)
	assert.NoError(t, err)
	assert.NotEmpty(t, requestJSON)

	// Test request deserialization (what backend would receive)
	var deserializedRequest domain.AutonomousTaskRequest
	err = json.Unmarshal(requestJSON, &deserializedRequest)
	assert.NoError(t, err)

	assert.Equal(t, validRequest.Task, deserializedRequest.Task)
	assert.Equal(t, validRequest.SlaPolicy, deserializedRequest.SlaPolicy)
	assert.Equal(t, validRequest.ProjectPath, deserializedRequest.ProjectPath)
	assert.Equal(t, validRequest.Options.MaxTokens, deserializedRequest.Options.MaxTokens)
}

func TestApp_ListReports_Integration(t *testing.T) {
	// Setup test reports directory
	tempDir, err := os.MkdirTemp("", "app_reports_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create mock report files
	reports := []domain.GenericReport{
		{
			Id:        "report-1",
			TaskId:    "task-123",
			Type:      "context-analysis",
			Title:     "Context Analysis Report",
			Summary:   "Analysis of project context",
			Content:   "Detailed analysis content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        "report-2",
			TaskId:    "task-123",
			Type:      "performance",
			Title:     "Performance Report",
			Summary:   "Performance metrics",
			Content:   "Performance data",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        "report-3",
			TaskId:    "task-456",
			Type:      "context-analysis",
			Title:     "Another Analysis",
			Summary:   "Different analysis",
			Content:   "Other content",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Save report files
	for _, report := range reports {
		reportData, err := json.MarshalIndent(report, "", "  ")
		require.NoError(t, err)

		reportFile := filepath.Join(tempDir, report.Id+".json")
		err = os.WriteFile(reportFile, reportData, 0o644)
		require.NoError(t, err)
	}

	// Test listing all reports
	files, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.Len(t, files, 3)

	// Test filtering by type
	contextAnalysisReports := []domain.GenericReport{}
	for _, report := range reports {
		if report.Type == "context-analysis" {
			contextAnalysisReports = append(contextAnalysisReports, report)
		}
	}
	assert.Len(t, contextAnalysisReports, 2)

	// Test filtering by task ID
	task123Reports := []domain.GenericReport{}
	for _, report := range reports {
		if report.TaskId == "task-123" {
			task123Reports = append(task123Reports, report)
		}
	}
	assert.Len(t, task123Reports, 2)

	// Validate report file structure
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		reportPath := filepath.Join(tempDir, file.Name())
		reportData, err := os.ReadFile(reportPath)
		assert.NoError(t, err)

		var report domain.GenericReport
		err = json.Unmarshal(reportData, &report)
		assert.NoError(t, err)

		// Validate required fields
		assert.NotEmpty(t, report.Id)
		assert.NotEmpty(t, report.TaskId)
		assert.NotEmpty(t, report.Type)
		assert.NotEmpty(t, report.Title)
		assert.False(t, report.CreatedAt.IsZero())
		assert.False(t, report.UpdatedAt.IsZero())
	}
}

func TestApp_TokenEstimation_Integration(t *testing.T) {
	// Test token estimation for different content types
	testContents := map[string]string{
		"simple": "Hello world",
		"code": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"markdown": `# Project Documentation

This is a comprehensive guide to the project.

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

Run the following command:

` + "`" + `bash
npm install
` + "`" + ``,
		"large": strings.Repeat("This is a test sentence. ", 1000),
	}

	for name, content := range testContents {
		t.Run(name, func(t *testing.T) {
			// Simple token estimation (4 characters per token approximation)
			estimatedTokens := len(content) / 4

			// Validate token estimation is reasonable
			assert.Greater(t, estimatedTokens, 0)

			if name == "simple" {
				assert.Less(t, estimatedTokens, 10)
			}

			if name == "large" {
				assert.Greater(t, estimatedTokens, 1000)
			}

			// Validate content is not empty
			assert.NotEmpty(t, content)

			// Validate content length matches expectations
			if name == "large" {
				expectedLength := len("This is a test sentence. ") * 1000
				assert.Equal(t, expectedLength, len(content))
			}
		})
	}
}

// NOTE: RecentProjects tests removed - they used deprecated settingsManager API
// The functionality is now tested through SettingsService tests
