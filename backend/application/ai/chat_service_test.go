package ai

import (
"context"
"testing"
"time"

"syntaxia/domain"
)

type mockToolExecutor struct {
tools   []domain.Tool
results map[string]domain.ToolResult
}

func newMockToolExecutor() *mockToolExecutor {
return &mockToolExecutor{
tools: []domain.Tool{
{
Name:        "read_file",
Description: "Read file contents",
Parameters: domain.ToolParameters{
Type: "object",
Properties: map[string]domain.ToolProperty{
"path": {Type: "string", Description: "File path"},
},
Required: []string{"path"},
},
},
},
results: make(map[string]domain.ToolResult),
}
}

func (m *mockToolExecutor) GetAvailableTools() []domain.Tool {
return m.tools
}

func (m *mockToolExecutor) ExecuteTool(call domain.ToolCall, _ string) domain.ToolResult {
if result, ok := m.results[call.Name]; ok {
result.ToolCallID = call.ID
return result
}
return domain.ToolResult{ToolCallID: call.ID, Content: "mock result"}
}

func TestDefaultChatServiceConfig(t *testing.T) {
config := DefaultChatServiceConfig()
if config.MaxIterations != 15 {
t.Errorf("MaxIterations: got %d, expected 15", config.MaxIterations)
}
if config.MaxHistoryLength != 50 {
t.Errorf("MaxHistoryLength: got %d, expected 50", config.MaxHistoryLength)
}
}

func TestNewChatService(t *testing.T) {
logger := &domain.NoopLogger{}
executor := newMockToolExecutor()
service := NewChatService(logger, nil, executor)
if service == nil {
t.Fatal("expected non-nil service")
}
}

func TestParseToolCallsSingle(t *testing.T) {
response := "{\"tool_calls\": [{\"name\": \"read_file\", \"arguments\": {}}]}"
calls := parseToolCalls(response, 10)
if len(calls) != 1 {
t.Errorf("got %d, expected 1", len(calls))
}
}

func TestParseToolCallsNone(t *testing.T) {
calls := parseToolCalls("Plain text", 10)
if len(calls) != 0 {
t.Errorf("got %d, expected 0", len(calls))
}
}

func TestParseToolCallsEmpty(t *testing.T) {
calls := parseToolCalls("{\"tool_calls\": []}", 10)
if len(calls) != 0 {
t.Errorf("got %d, expected 0", len(calls))
}
}

func TestParseToolCallsLimit(t *testing.T) {
response := "{\"tool_calls\": [{\"name\": \"a\", \"arguments\": {}}, {\"name\": \"b\", \"arguments\": {}}]}"
calls := parseToolCalls(response, 1)
if len(calls) != 1 {
t.Errorf("got %d, expected 1", len(calls))
}
}

func TestCleanResponseOnlyJSON(t *testing.T) {
result := cleanResponseFromToolCalls("{\"tool_calls\": []}")
if result != "" {
t.Errorf("got %q, expected empty", result)
}
}

func TestCleanResponseTextBefore(t *testing.T) {
result := cleanResponseFromToolCalls("Text {\"tool_calls\": []}")
if result != "Text" {
t.Errorf("got %q, expected Text", result)
}
}

func TestCleanResponseNoJSON(t *testing.T) {
result := cleanResponseFromToolCalls("Plain")
if result != "Plain" {
t.Errorf("got %q, expected Plain", result)
}
}

func TestFormatSmartContextNil(t *testing.T) {
if result := formatSmartContext(nil); result != "" {
t.Errorf("nil context should return empty, got %q", result)
}
}

func TestFormatSmartContextWithStructure(t *testing.T) {
ctx := &domain.SmartContextResult{ProjectStructure: "src/"}
if result := formatSmartContext(ctx); result == "" {
t.Error("expected non-empty result")
}
}

func TestChatServicePrepareMessages(t *testing.T) {
service := NewChatService(&domain.NoopLogger{}, nil, newMockToolExecutor())
req := ChatRequest{Messages: []domain.ChatMessage{{Role: domain.RoleUser, Content: "Hi"}}}
messages := service.prepareMessages(req)
if len(messages) < 2 {
t.Errorf("expected at least 2 messages, got %d", len(messages))
}
}

func TestChatServiceExecuteSingleTool(t *testing.T) {
executor := newMockToolExecutor()
executor.results["read_file"] = domain.ToolResult{Content: "content"}
service := NewChatService(&domain.NoopLogger{}, nil, executor)
call := domain.ToolCall{ID: "test", Name: "read_file", Arguments: map[string]any{"path": "f.go"}}
result := service.executeSingleTool(context.Background(), call, "/p")
if result.Result != "content" {
t.Errorf("got %q", result.Result)
}
}

type slowToolExecutor struct{ delay time.Duration }

func (s *slowToolExecutor) GetAvailableTools() []domain.Tool { return nil }
func (s *slowToolExecutor) ExecuteTool(call domain.ToolCall, _ string) domain.ToolResult {
time.Sleep(s.delay)
return domain.ToolResult{ToolCallID: call.ID, Content: "done"}
}

func TestChatServiceExecuteToolTimeout(t *testing.T) {
config := DefaultChatServiceConfig()
config.ToolTimeout = 1 * time.Millisecond
service := NewChatServiceWithConfig(&domain.NoopLogger{}, nil, &slowToolExecutor{delay: 50 * time.Millisecond}, config)
call := domain.ToolCall{ID: "slow", Name: "slow"}
result := service.executeSingleTool(context.Background(), call, "/p")
if result.Error != "tool execution timed out" {
t.Errorf("expected timeout, got %q", result.Error)
}
}