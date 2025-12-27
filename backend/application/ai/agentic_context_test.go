package ai

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"syntaxia/domain"
)

// MockAIProvider implements domain.AIProvider for testing
type MockAIProvider struct {
	responses []string
	callCount int
}

func (m *MockAIProvider) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	if m.callCount >= len(m.responses) {
		return domain.AIResponse{Content: `{"tool_calls": [{"name": "mark_context_complete", "arguments": {"reasoning": "No more responses"}}]}`}, nil
	}
	response := m.responses[m.callCount]
	m.callCount++
	return domain.AIResponse{Content: response, TokensUsed: 100}, nil
}

func (m *MockAIProvider) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	return fmt.Errorf("not implemented")
}

func (m *MockAIProvider) ListModels(ctx context.Context) ([]string, error) {
	return []string{"test-model"}, nil
}

func (m *MockAIProvider) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{Name: "mock"}
}

func (m *MockAIProvider) ValidateRequest(req domain.AIRequest) error {
	return nil
}

func (m *MockAIProvider) EstimateTokens(req domain.AIRequest) (int, error) {
	return 100, nil
}

func (m *MockAIProvider) GetPricing(model string) domain.PricingInfo {
	return domain.PricingInfo{}
}

// MockFileTools implements FileToolsProvider for testing
type MockFileTools struct {
	files   map[string]string
	matches []ContentMatch
}

func NewMockFileTools() *MockFileTools {
	return &MockFileTools{
		files: map[string]string{
			"auth/AuthService.ts": `export class AuthService {
	login(username: string, password: string): Promise<User> {}
	logout(): void {}
	refreshToken(): Promise<string> {}
}`,
			"config/oauth.ts": `export const oauthConfig = {
	googleClientId: "xxx",
	googleClientSecret: "yyy",
};`,
			"api/auth.api.ts": `import { AuthService } from '../auth/AuthService';
export function setupAuthRoutes(app: Express) {}`,
		},
		matches: []ContentMatch{},
	}
}

func (m *MockFileTools) SearchFiles(pattern, directory, projectRoot string) ([]string, error) {
	var results []string
	patternLower := strings.ToLower(pattern)
	for path := range m.files {
		if strings.Contains(strings.ToLower(path), patternLower) {
			results = append(results, path)
		}
	}
	return results, nil
}

func (m *MockFileTools) SearchContent(pattern, filePattern, projectRoot string, maxResults int) ([]ContentMatch, error) {
	var results []ContentMatch
	patternLower := strings.ToLower(pattern)

	for path, content := range m.files {
		lines := strings.Split(content, "\n")
		for i, line := range lines {
			if strings.Contains(strings.ToLower(line), patternLower) {
				results = append(results, ContentMatch{
					FilePath: path,
					Line:     i + 1,
					Content:  strings.TrimSpace(line),
				})
				if len(results) >= maxResults {
					return results, nil
				}
			}
		}
	}
	return results, nil
}

func (m *MockFileTools) ReadFile(path, projectRoot string) (string, error) {
	if content, ok := m.files[path]; ok {
		return content, nil
	}
	return "", fmt.Errorf("file not found: %s", path)
}

// MockSymbolTools implements SymbolToolsProvider for testing
type MockSymbolTools struct {
	symbols map[string][]SymbolInfo
}

func NewMockSymbolTools() *MockSymbolTools {
	return &MockSymbolTools{
		symbols: map[string][]SymbolInfo{
			"auth/AuthService.ts": {
				{Name: "AuthService", Kind: "class", FilePath: "auth/AuthService.ts", Line: 1},
				{Name: "login", Kind: "method", FilePath: "auth/AuthService.ts", Line: 2},
				{Name: "logout", Kind: "method", FilePath: "auth/AuthService.ts", Line: 3},
				{Name: "refreshToken", Kind: "method", FilePath: "auth/AuthService.ts", Line: 4},
			},
			"config/oauth.ts": {
				{Name: "oauthConfig", Kind: "constant", FilePath: "config/oauth.ts", Line: 1},
			},
		},
	}
}

func (m *MockSymbolTools) ListSymbols(path, projectRoot string) ([]SymbolInfo, error) {
	if symbols, ok := m.symbols[path]; ok {
		return symbols, nil
	}
	return []SymbolInfo{}, nil
}

func (m *MockSymbolTools) SearchSymbols(query, kind, projectRoot string) ([]SymbolInfo, error) {
	var results []SymbolInfo
	queryLower := strings.ToLower(query)

	for _, symbols := range m.symbols {
		for _, s := range symbols {
			if strings.Contains(strings.ToLower(s.Name), queryLower) {
				if kind == "" || strings.EqualFold(s.Kind, kind) {
					results = append(results, s)
				}
			}
		}
	}
	return results, nil
}

// MockLogger implements domain.Logger for testing
type MockLogger struct{}

func (m *MockLogger) Debug(msg string)                          {}
func (m *MockLogger) Info(msg string)                           {}
func (m *MockLogger) Warning(msg string)                        {}
func (m *MockLogger) Error(msg string)                          {}
func (m *MockLogger) Fatal(msg string)                          {}
func (m *MockLogger) WithField(key string, value any) domain.Logger { return m }
func (m *MockLogger) WithFields(fields map[string]any) domain.Logger { return m }

func TestAgenticContextService_GatherContext_BasicFlow(t *testing.T) {
	// Setup mock AI that simulates context gathering flow
	mockAI := &MockAIProvider{
		responses: []string{
			// Iteration 1: Search for auth files
			`{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "auth"}}]}`,
			// Iteration 2: Get symbols and request file
			`{"tool_calls": [
				{"name": "get_file_symbols", "arguments": {"path": "auth/AuthService.ts"}},
				{"name": "request_file", "arguments": {"path": "auth/AuthService.ts", "reason": "main auth service"}}
			]}`,
			// Iteration 3: Search for oauth and request config
			`{"tool_calls": [
				{"name": "search_content_for_context", "arguments": {"query": "oauth"}},
				{"name": "request_file", "arguments": {"path": "config/oauth.ts", "reason": "oauth configuration"}}
			]}`,
			// Iteration 4: Mark complete
			`{"tool_calls": [{"name": "mark_context_complete", "arguments": {"reasoning": "Found auth service and oauth config"}}]}`,
		},
	}

	service := NewAgenticContextService(
		mockAI,
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	ctx := context.Background()
	result, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:          "Add Google OAuth authentication",
		ProjectRoot:   "/test/project",
		MaxIterations: 10,
		MaxTokens:     100000,
	})

	if err != nil {
		t.Fatalf("GatherContext failed: %v", err)
	}

	// Verify results
	if len(result.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result.Files))
	}

	// Check that expected files are in context
	fileMap := make(map[string]bool)
	for _, f := range result.Files {
		fileMap[f.Path] = true
	}

	if !fileMap["auth/AuthService.ts"] {
		t.Error("Expected auth/AuthService.ts in context")
	}
	if !fileMap["config/oauth.ts"] {
		t.Error("Expected config/oauth.ts in context")
	}

	// Verify iterations
	if result.Iterations != 4 {
		t.Errorf("Expected 4 iterations, got %d", result.Iterations)
	}

	// Verify tool calls were logged
	if len(result.ToolCalls) < 4 {
		t.Errorf("Expected at least 4 tool calls, got %d", len(result.ToolCalls))
	}
}

func TestAgenticContextService_GatherContext_MaxIterations(t *testing.T) {
	// Mock AI that never completes
	mockAI := &MockAIProvider{
		responses: []string{
			`{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "test"}}]}`,
			`{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "test2"}}]}`,
			`{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "test3"}}]}`,
		},
	}

	service := NewAgenticContextService(
		mockAI,
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	ctx := context.Background()
	result, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:          "Test task",
		ProjectRoot:   "/test/project",
		MaxIterations: 3,
	})

	if err != nil {
		t.Fatalf("GatherContext failed: %v", err)
	}

	if result.Iterations != 3 {
		t.Errorf("Expected 3 iterations (max), got %d", result.Iterations)
	}
}

func TestAgenticContextService_GatherContext_TokenBudget(t *testing.T) {
	mockFileTools := NewMockFileTools()
	// Add a large file
	mockFileTools.files["large/file.ts"] = strings.Repeat("x", 10000) // ~2500 tokens

	mockAI := &MockAIProvider{
		responses: []string{
			`{"tool_calls": [{"name": "request_file", "arguments": {"path": "large/file.ts", "reason": "test"}}]}`,
			`{"tool_calls": [{"name": "request_file", "arguments": {"path": "auth/AuthService.ts", "reason": "test"}}]}`,
		},
	}

	service := NewAgenticContextService(
		mockAI,
		mockFileTools,
		NewMockSymbolTools(),
		&MockLogger{},
	)

	ctx := context.Background()
	result, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:          "Test task",
		ProjectRoot:   "/test/project",
		MaxTokens:     3000, // Budget that will be exceeded after first file
		MaxIterations: 10,
	})

	if err != nil {
		t.Fatalf("GatherContext failed: %v", err)
	}

	// Should stop due to token budget
	if result.TotalTokens < 2000 {
		t.Errorf("Expected tokens > 2000, got %d", result.TotalTokens)
	}
}

func TestAgenticContextService_GatherContext_DuplicateFiles(t *testing.T) {
	mockAI := &MockAIProvider{
		responses: []string{
			`{"tool_calls": [{"name": "request_file", "arguments": {"path": "auth/AuthService.ts", "reason": "first request"}}]}`,
			`{"tool_calls": [{"name": "request_file", "arguments": {"path": "auth/AuthService.ts", "reason": "duplicate request"}}]}`,
			`{"tool_calls": [{"name": "mark_context_complete", "arguments": {"reasoning": "done"}}]}`,
		},
	}

	service := NewAgenticContextService(
		mockAI,
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	ctx := context.Background()
	result, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:        "Test task",
		ProjectRoot: "/test/project",
	})

	if err != nil {
		t.Fatalf("GatherContext failed: %v", err)
	}

	// Should only have one file despite two requests
	if len(result.Files) != 1 {
		t.Errorf("Expected 1 file (no duplicates), got %d", len(result.Files))
	}
}

func TestAgenticContextService_GatherContext_NoToolCalls(t *testing.T) {
	mockAI := &MockAIProvider{
		responses: []string{
			"I don't need any tools, the task is clear enough.",
		},
	}

	service := NewAgenticContextService(
		mockAI,
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	ctx := context.Background()
	result, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:        "Simple task",
		ProjectRoot: "/test/project",
	})

	if err != nil {
		t.Fatalf("GatherContext failed: %v", err)
	}

	if result.Iterations != 1 {
		t.Errorf("Expected 1 iteration, got %d", result.Iterations)
	}

	if len(result.Files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(result.Files))
	}
}

func TestAgenticContextService_GatherContext_ContextTimeout(t *testing.T) {
	// MockAIProvider that respects context cancellation
	mockAI := &MockAIProviderWithContext{
		responses: []string{
			`{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "test"}}]}`,
		},
	}

	service := NewAgenticContextService(
		mockAI,
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	// Create already cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := service.GatherContext(ctx, AgenticContextRequest{
		Task:        "Test task",
		ProjectRoot: "/test/project",
	})

	if err == nil {
		t.Error("Expected error due to cancelled context")
	}
}

// MockAIProviderWithContext respects context cancellation
type MockAIProviderWithContext struct {
	responses []string
	callCount int
}

func (m *MockAIProviderWithContext) Generate(ctx context.Context, req domain.AIRequest) (domain.AIResponse, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return domain.AIResponse{}, ctx.Err()
	default:
	}

	if m.callCount >= len(m.responses) {
		return domain.AIResponse{Content: `{"tool_calls": [{"name": "mark_context_complete", "arguments": {"reasoning": "No more responses"}}]}`}, nil
	}
	response := m.responses[m.callCount]
	m.callCount++
	return domain.AIResponse{Content: response, TokensUsed: 100}, nil
}

func (m *MockAIProviderWithContext) GenerateStream(ctx context.Context, req domain.AIRequest, onChunk func(chunk domain.StreamChunk)) error {
	return fmt.Errorf("not implemented")
}

func (m *MockAIProviderWithContext) ListModels(ctx context.Context) ([]string, error) {
	return []string{"test-model"}, nil
}

func (m *MockAIProviderWithContext) GetProviderInfo() domain.ProviderInfo {
	return domain.ProviderInfo{Name: "mock"}
}

func (m *MockAIProviderWithContext) ValidateRequest(req domain.AIRequest) error {
	return nil
}

func (m *MockAIProviderWithContext) EstimateTokens(req domain.AIRequest) (int, error) {
	return 100, nil
}

func (m *MockAIProviderWithContext) GetPricing(model string) domain.PricingInfo {
	return domain.PricingInfo{}
}

func TestAgenticContextService_ExecuteTools(t *testing.T) {
	service := NewAgenticContextService(
		&MockAIProvider{},
		NewMockFileTools(),
		NewMockSymbolTools(),
		&MockLogger{},
	)

	tests := []struct {
		name        string
		toolCall    parsedToolCall
		wantErr     bool
		checkResult func(t *testing.T, result string)
	}{
		{
			name: "search_files_for_context",
			toolCall: parsedToolCall{
				Name:      "search_files_for_context",
				Arguments: map[string]any{"pattern": "auth"},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result string) {
				if !strings.Contains(result, "AuthService") {
					t.Error("Expected result to contain AuthService")
				}
			},
		},
		{
			name: "search_content_for_context",
			toolCall: parsedToolCall{
				Name:      "search_content_for_context",
				Arguments: map[string]any{"query": "login"},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result string) {
				if !strings.Contains(result, "login") {
					t.Error("Expected result to contain login")
				}
			},
		},
		{
			name: "get_file_symbols",
			toolCall: parsedToolCall{
				Name:      "get_file_symbols",
				Arguments: map[string]any{"path": "auth/AuthService.ts"},
			},
			wantErr: false,
			checkResult: func(t *testing.T, result string) {
				if !strings.Contains(result, "AuthService") {
					t.Error("Expected result to contain AuthService symbol")
				}
			},
		},
		{
			name: "unknown_tool",
			toolCall: parsedToolCall{
				Name:      "unknown_tool",
				Arguments: map[string]any{},
			},
			wantErr: true,
		},
		{
			name: "search_files_missing_pattern",
			toolCall: parsedToolCall{
				Name:      "search_files_for_context",
				Arguments: map[string]any{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestedFiles := make(map[string]bool)
			result := &AgenticContextResult{Files: []ContextFileResult{}}

			toolResult, err := service.executeTool(tt.toolCall, "/test/project", requestedFiles, result)

			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.checkResult != nil {
				tt.checkResult(t, toolResult)
			}
		})
	}
}

func TestAgenticContextService_ParseToolCalls(t *testing.T) {
	service := &AgenticContextService{}

	tests := []struct {
		name     string
		response string
		expected int
	}{
		{
			name:     "single tool call",
			response: `{"tool_calls": [{"name": "search_files_for_context", "arguments": {"pattern": "test"}}]}`,
			expected: 1,
		},
		{
			name:     "multiple tool calls",
			response: `{"tool_calls": [{"name": "tool1", "arguments": {}}, {"name": "tool2", "arguments": {}}]}`,
			expected: 2,
		},
		{
			name:     "no tool calls",
			response: "Just a regular response without any JSON",
			expected: 0,
		},
		{
			name:     "invalid JSON",
			response: `{"tool_calls": [invalid]}`,
			expected: 0,
		},
		{
			name:     "tool call with text around",
			response: `Let me search for files. {"tool_calls": [{"name": "search", "arguments": {}}]} That should help.`,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := service.parseToolCalls(tt.response)
			if len(calls) != tt.expected {
				t.Errorf("Expected %d tool calls, got %d", tt.expected, len(calls))
			}
		})
	}
}

func TestAgenticContextService_IsContextComplete(t *testing.T) {
	service := &AgenticContextService{}

	tests := []struct {
		name     string
		calls    []parsedToolCall
		expected bool
	}{
		{
			name:     "empty calls",
			calls:    []parsedToolCall{},
			expected: false,
		},
		{
			name: "no complete marker",
			calls: []parsedToolCall{
				{Name: "search_files_for_context", Arguments: map[string]any{}},
			},
			expected: false,
		},
		{
			name: "has complete marker",
			calls: []parsedToolCall{
				{Name: "mark_context_complete", Arguments: map[string]any{"reasoning": "done"}},
			},
			expected: true,
		},
		{
			name: "complete marker among others",
			calls: []parsedToolCall{
				{Name: "request_file", Arguments: map[string]any{}},
				{Name: "mark_context_complete", Arguments: map[string]any{}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.isContextComplete(tt.calls)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestContextFileTools_SearchFiles(t *testing.T) {
	// This test requires actual filesystem, skip in unit tests
	t.Skip("Requires filesystem setup")
}

func TestContextFileTools_ReadFile_PathTraversal(t *testing.T) {
	tools := NewContextFileTools(&MockLogger{})

	_, err := tools.ReadFile("../../../etc/passwd", "/test/project")
	if err == nil {
		t.Error("Expected error for path traversal attempt")
	}
	if !strings.Contains(err.Error(), "path traversal") {
		t.Errorf("Expected path traversal error, got: %v", err)
	}
}

func TestSimpleSymbolTools_ExtractGoSymbols(t *testing.T) {
	content := `package main

func main() {
	fmt.Println("Hello")
}

type User struct {
	Name string
}

type Reader interface {
	Read() error
}

const MaxSize = 100

var globalVar = "test"
`

	symbols := extractGoSymbols(content, "test.go")

	expectedSymbols := map[string]string{
		"main":      "function",
		"User":      "class",
		"Reader":    "interface",
		"MaxSize":   "constant",
		"globalVar": "variable",
	}

	for name, kind := range expectedSymbols {
		found := false
		for _, s := range symbols {
			if s.Name == name && s.Kind == kind {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected symbol %s of kind %s not found", name, kind)
		}
	}
}

func TestSimpleSymbolTools_ExtractTSSymbols(t *testing.T) {
	content := `export function greet(name: string): string {
	return "Hello " + name;
}

export class UserService {
	getUser(id: number): User {}
}

export interface User {
	id: number;
	name: string;
}

export type UserId = number;

export const API_URL = "https://api.example.com";
`

	symbols := extractTSSymbols(content, "test.ts")

	expectedSymbols := map[string]string{
		"greet":       "function",
		"UserService": "class",
		"User":        "interface",
		"UserId":      "type",
		"API_URL":     "constant",
	}

	for name, kind := range expectedSymbols {
		found := false
		for _, s := range symbols {
			if s.Name == name && s.Kind == kind {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected symbol %s of kind %s not found", name, kind)
		}
	}
}

func TestSimpleSymbolTools_ExtractPythonSymbols(t *testing.T) {
	content := `def greet(name):
    return f"Hello {name}"

async def fetch_data(url):
    pass

class UserService:
    def get_user(self, id):
        pass
`

	symbols := extractPythonSymbols(content, "test.py")

	expectedSymbols := map[string]string{
		"greet":       "function",
		"fetch_data":  "function",
		"UserService": "class",
	}

	for name, kind := range expectedSymbols {
		found := false
		for _, s := range symbols {
			if s.Name == name && s.Kind == kind {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected symbol %s of kind %s not found", name, kind)
		}
	}
}

func TestDefaultValues(t *testing.T) {
	if DefaultMaxIterations != 5 {
		t.Errorf("Expected DefaultMaxIterations to be 5, got %d", DefaultMaxIterations)
	}

	if DefaultIterationTimeout != 30*time.Second {
		t.Errorf("Expected DefaultIterationTimeout to be 30s, got %v", DefaultIterationTimeout)
	}

	if DefaultMaxContextTokens != 100000 {
		t.Errorf("Expected DefaultMaxContextTokens to be 100000, got %d", DefaultMaxContextTokens)
	}
}
