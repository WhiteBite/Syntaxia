package tools

import (
	"testing"

	"syntaxia/domain"
)

// mockCallGraphBuilder implements domain.CallGraphBuilder for testing
type mockCallGraphBuilder struct{}

func (m *mockCallGraphBuilder) Build(projectRoot string) (*domain.CallGraph, error) {
	return &domain.CallGraph{Nodes: make(map[string]*domain.CallGraphNode)}, nil
}

func (m *mockCallGraphBuilder) GetCallers(functionID string) []domain.CallGraphNode {
	return nil
}

func (m *mockCallGraphBuilder) GetCallees(functionID string) []domain.CallGraphNode {
	return nil
}

func (m *mockCallGraphBuilder) GetImpact(functionID string, maxDepth int) []domain.CallGraphNode {
	return nil
}

func (m *mockCallGraphBuilder) GetCallChain(startID, endID string, maxDepth int) [][]string {
	return nil
}

func TestCallGraphHandler_CanHandle(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	tests := []struct {
		toolName string
		expected bool
	}{
		{"get_callers", true},
		{"get_callees", true},
		{"get_impact", true},
		{"get_call_chain", true},
		{"unknown_tool", false},
		{"git_status", false},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			result := handler.CanHandle(tt.toolName)
			if result != tt.expected {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.toolName, result, tt.expected)
			}
		})
	}
}

func TestCallGraphHandler_GetTools(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)
	tools := handler.GetTools()

	if len(tools) != 4 {
		t.Errorf("expected 4 tools, got %d", len(tools))
	}

	expectedTools := map[string]bool{
		"get_callers":    false,
		"get_callees":    false,
		"get_impact":     false,
		"get_call_chain": false,
	}

	for _, tool := range tools {
		if _, ok := expectedTools[tool.Name]; ok {
			expectedTools[tool.Name] = true
		}
	}

	for name, found := range expectedTools {
		if !found {
			t.Errorf("expected tool %s not found", name)
		}
	}
}

func TestGetCallers_NotInitialized(t *testing.T) {
	handler := NewCallGraphToolsHandler(nil, nil)

	_, err := handler.Execute("get_callers", map[string]any{
		"function_id": "pkg.Function",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when call graph not initialized")
	}
}

func TestGetCallers_MissingFunctionID(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	_, err := handler.Execute("get_callers", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when function_id missing")
	}
}

func TestGetCallees_MissingFunctionID(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	_, err := handler.Execute("get_callees", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when function_id missing")
	}
}

func TestGetImpact_MissingFunctionID(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	_, err := handler.Execute("get_impact", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when function_id missing")
	}
}

func TestGetCallChain_MissingParams(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	tests := []struct {
		name string
		args map[string]any
	}{
		{"missing both", map[string]any{}},
		{"missing end", map[string]any{"start_function": "a"}},
		{"missing start", map[string]any{"end_function": "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Execute("get_call_chain", tt.args, "/project")
			if err == nil {
				t.Fatal("expected error for missing params")
			}
		})
	}
}

func TestCallGraphHandler_UnknownTool(t *testing.T) {
	callGraph := &mockCallGraphBuilder{}
	handler := NewCallGraphToolsHandler(nil, callGraph)

	_, err := handler.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}
