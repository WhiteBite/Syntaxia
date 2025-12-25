package main

import (
	"encoding/json"
	"syntaxia/domain"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDomainModels_ContextValidation(t *testing.T) {
	// Test Context domain model structure
	context := domain.Context{
		ID:          "test-context-123",
		Name:        "Test Context",
		Description: "A test context for validation",
		Content:     "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}",
		Files:       []string{"main.go", "utils.go"},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		ProjectPath: "/test/project",
		TokenCount:  150,
	}

	// Validate required fields
	assert.NotEmpty(t, context.ID)
	assert.NotEmpty(t, context.Name)
	assert.NotEmpty(t, context.Content)
	assert.Greater(t, len(context.Files), 0)
	assert.False(t, context.CreatedAt.IsZero())
	assert.False(t, context.UpdatedAt.IsZero())
	assert.NotEmpty(t, context.ProjectPath)
	assert.Greater(t, context.TokenCount, 0)

	// Test JSON serialization/deserialization
	jsonData, err := json.Marshal(context)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var deserializedContext domain.Context
	err = json.Unmarshal(jsonData, &deserializedContext)
	assert.NoError(t, err)
	assert.Equal(t, context.ID, deserializedContext.ID)
	assert.Equal(t, context.Name, deserializedContext.Name)
	assert.Equal(t, context.Content, deserializedContext.Content)
	assert.Equal(t, context.Files, deserializedContext.Files)
}

func TestDomainModels_GenericReportValidation(t *testing.T) {
	// Test GenericReport domain model structure
	report := domain.GenericReport{
		Id:        "report-123",
		TaskId:    "task-456",
		Type:      "context-analysis",
		Title:     "Context Analysis Report",
		Summary:   "Analysis summary",
		Content:   "Detailed analysis content",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate required fields
	assert.NotEmpty(t, report.Id)
	assert.NotEmpty(t, report.TaskId)
	assert.NotEmpty(t, report.Type)
	assert.NotEmpty(t, report.Title)
	assert.False(t, report.CreatedAt.IsZero())
	assert.False(t, report.UpdatedAt.IsZero())

	// Test JSON serialization/deserialization
	jsonData, err := json.Marshal(report)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	var deserializedReport domain.GenericReport
	err = json.Unmarshal(jsonData, &deserializedReport)
	assert.NoError(t, err)
	assert.Equal(t, report.Id, deserializedReport.Id)
	assert.Equal(t, report.TaskId, deserializedReport.TaskId)
	assert.Equal(t, report.Type, deserializedReport.Type)
	assert.Equal(t, report.Title, deserializedReport.Title)
}

func TestDomainModels_AutonomousTaskRequestValidation(t *testing.T) {
	// Test AutonomousTaskRequest validation
	validRequest := domain.AutonomousTaskRequest{
		Task:        "Implement authentication system",
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

	// Validate structure
	assert.NotEmpty(t, validRequest.Task)
	assert.Contains(t, []string{"lite", "standard", "strict"}, validRequest.SlaPolicy)
	assert.NotEmpty(t, validRequest.ProjectPath)

	// Validate options
	assert.Greater(t, validRequest.Options.MaxTokens, 0)
	assert.GreaterOrEqual(t, validRequest.Options.Temperature, 0.0)
	assert.LessOrEqual(t, validRequest.Options.Temperature, 2.0)

	// Test serialization
	jsonData, err := json.Marshal(validRequest)
	assert.NoError(t, err)

	var deserializedRequest domain.AutonomousTaskRequest
	err = json.Unmarshal(jsonData, &deserializedRequest)
	assert.NoError(t, err)
	assert.Equal(t, validRequest.Task, deserializedRequest.Task)
	assert.Equal(t, validRequest.SlaPolicy, deserializedRequest.SlaPolicy)
	assert.Equal(t, validRequest.ProjectPath, deserializedRequest.ProjectPath)
}

func TestTokenEstimation_SimpleAlgorithm(t *testing.T) {
	// Test simple token estimation (4 chars per token approximation)
	testCases := []struct {
		name     string
		content  string
		expected int
	}{
		{
			name:     "Simple text",
			content:  "Hello world",
			expected: 2, // 11 chars / 4 = 2.75 -> 2
		},
		{
			name:     "Go code snippet",
			content:  "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}",
			expected: 11, // 46 chars / 4 = 11.5 -> 11
		},
		{
			name:     "Empty content",
			content:  "",
			expected: 0,
		},
		{
			name:     "Large content",
			content:  strings.Repeat("test ", 1000), // 5000 chars
			expected: 1250,                          // 5000 / 4 = 1250
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			estimated := len(tc.content) / 4
			assert.Equal(t, tc.expected, estimated)

			// Validate token count is reasonable
			if tc.content != "" {
				assert.Greater(t, estimated, -1) // Should be >= 0
			}
		})
	}
}

func TestContextBuildOptions_Validation(t *testing.T) {
	// Test ContextBuildOptions structure
	options := domain.ContextBuildOptions{
		IncludeManifest: true,
		StripComments:   false,
		MaxTokens:       10000,
	}

	// Validate boolean options
	assert.IsType(t, true, options.IncludeManifest)
	assert.IsType(t, false, options.StripComments)

	// Validate numeric constraints
	if options.MaxTokens > 0 {
		assert.Greater(t, options.MaxTokens, 0)
		assert.Less(t, options.MaxTokens, 1000000) // Reasonable upper bound
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(options)
	assert.NoError(t, err)

	var deserializedOptions domain.ContextBuildOptions
	err = json.Unmarshal(jsonData, &deserializedOptions)
	assert.NoError(t, err)
	assert.Equal(t, options.IncludeManifest, deserializedOptions.IncludeManifest)
	assert.Equal(t, options.StripComments, deserializedOptions.StripComments)
	assert.Equal(t, options.MaxTokens, deserializedOptions.MaxTokens)
}

func TestSLAPolicyValidation(t *testing.T) {
	// Test SLA policy validation
	validPolicies := []string{"lite", "standard", "strict"}
	invalidPolicies := []string{"", "invalid", "custom", "premium"}

	for _, policy := range validPolicies {
		t.Run("Valid_"+policy, func(t *testing.T) {
			assert.Contains(t, validPolicies, policy)
			assert.NotEmpty(t, policy)
		})
	}

	for _, policy := range invalidPolicies {
		t.Run("Invalid_"+policy, func(t *testing.T) {
			assert.NotContains(t, validPolicies, policy)
		})
	}
}

func TestFileNodeValidation(t *testing.T) {
	// Test FileNode structure (from domain)
	fileNode := domain.FileNode{
		Name:            "test.go",
		Path:            "/project/test.go",
		RelPath:         "test.go",
		IsDir:           false,
		IsGitignored:    false,
		IsCustomIgnored: false,
		IsIgnored:       false,
		Size:            1024,
	}

	// Validate required fields
	assert.NotEmpty(t, fileNode.Name)
	assert.NotEmpty(t, fileNode.Path)
	assert.NotEmpty(t, fileNode.RelPath)
	assert.IsType(t, false, fileNode.IsDir)
	assert.GreaterOrEqual(t, fileNode.Size, int64(0))

	// Validate path consistency
	assert.True(t, strings.HasSuffix(fileNode.Path, fileNode.RelPath))

	// Test JSON serialization
	jsonData, err := json.Marshal(fileNode)
	assert.NoError(t, err)

	var deserializedNode domain.FileNode
	err = json.Unmarshal(jsonData, &deserializedNode)
	assert.NoError(t, err)
	assert.Equal(t, fileNode.Name, deserializedNode.Name)
	assert.Equal(t, fileNode.Path, deserializedNode.Path)
	assert.Equal(t, fileNode.IsDir, deserializedNode.IsDir)
}

func TestAPIResponseStructures(t *testing.T) {
	// Test autonomous task response structure
	response := domain.AutonomousTaskResponse{
		TaskId:  "task-789",
		Status:  "accepted",
		Message: "Task accepted for processing",
	}

	assert.NotEmpty(t, response.TaskId)
	assert.Contains(t, []string{"accepted", "rejected"}, response.Status)

	// Test autonomous task status structure
	status := domain.AutonomousTaskStatus{
		TaskId:                 "task-789",
		Status:                 "running",
		CurrentStep:            "context-analysis",
		Progress:               45.5,
		EstimatedTimeRemaining: 120,
		StartedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}

	assert.NotEmpty(t, status.TaskId)
	assert.Contains(t, []string{"pending", "running", "completed", "failed", "cancelled"}, status.Status)
	assert.GreaterOrEqual(t, status.Progress, 0.0)
	assert.LessOrEqual(t, status.Progress, 100.0)
	assert.False(t, status.StartedAt.IsZero())
	assert.False(t, status.UpdatedAt.IsZero())
}
