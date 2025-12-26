package tools

import (
	"strings"
	"testing"

	"syntaxia/domain"
)

// MockToolHandler implements ToolHandler for testing
type MockToolHandler struct {
	tools       []domain.Tool
	canHandle   map[string]bool
	executeFunc func(toolName string, args map[string]any, projectRoot string) (string, error)
}

func (m *MockToolHandler) GetTools() []domain.Tool {
	return m.tools
}

func (m *MockToolHandler) CanHandle(toolName string) bool {
	if m.canHandle == nil {
		return false
	}
	return m.canHandle[toolName]
}

func (m *MockToolHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	if m.executeFunc != nil {
		return m.executeFunc(toolName, args, projectRoot)
	}
	return "executed: " + toolName, nil
}

func TestNewHandlerRegistry(t *testing.T) {
	registry := NewHandlerRegistry(nil)

	if registry == nil {
		t.Fatal("expected non-nil registry")
	}
	if len(registry.handlers) != 0 {
		t.Errorf("expected empty handlers, got: %d", len(registry.handlers))
	}
}

func TestHandlerRegistry_Register(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler := &MockToolHandler{
		tools: []domain.Tool{{Name: "test_tool"}},
	}

	registry.Register(handler)

	if len(registry.handlers) != 1 {
		t.Errorf("expected 1 handler, got: %d", len(registry.handlers))
	}
}

func TestHandlerRegistry_GetAllTools(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler1 := &MockToolHandler{
		tools: []domain.Tool{{Name: "tool1"}, {Name: "tool2"}},
	}
	handler2 := &MockToolHandler{
		tools: []domain.Tool{{Name: "tool3"}},
	}

	registry.Register(handler1)
	registry.Register(handler2)

	tools := registry.GetAllTools()

	if len(tools) != 3 {
		t.Errorf("expected 3 tools, got: %d", len(tools))
	}
}

func TestHandlerRegistry_GetAllTools_Empty(t *testing.T) {
	registry := NewHandlerRegistry(nil)

	tools := registry.GetAllTools()

	if len(tools) != 0 {
		t.Errorf("expected 0 tools, got: %d", len(tools))
	}
}

func TestHandlerRegistry_Execute_Success(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler := &MockToolHandler{
		canHandle: map[string]bool{"my_tool": true},
		executeFunc: func(toolName string, args map[string]any, projectRoot string) (string, error) {
			return "result from " + toolName, nil
		},
	}

	registry.Register(handler)

	result, err := registry.Execute("my_tool", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "my_tool") {
		t.Errorf("expected tool name in result, got: %s", result)
	}
}

func TestHandlerRegistry_Execute_NoHandler(t *testing.T) {
	registry := NewHandlerRegistry(nil)

	_, err := registry.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "no handler found") {
		t.Errorf("expected 'no handler found' error, got: %v", err)
	}
}

func TestHandlerRegistry_CanHandle_True(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler := &MockToolHandler{
		canHandle: map[string]bool{"supported_tool": true},
	}

	registry.Register(handler)

	if !registry.CanHandle("supported_tool") {
		t.Error("expected CanHandle to return true")
	}
}

func TestHandlerRegistry_CanHandle_False(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler := &MockToolHandler{
		canHandle: map[string]bool{"other_tool": true},
	}

	registry.Register(handler)

	if registry.CanHandle("unsupported_tool") {
		t.Error("expected CanHandle to return false")
	}
}

func TestHandlerRegistry_CanHandle_Empty(t *testing.T) {
	registry := NewHandlerRegistry(nil)

	if registry.CanHandle("any_tool") {
		t.Error("expected CanHandle to return false for empty registry")
	}
}

func TestHandlerRegistry_MultipleHandlers(t *testing.T) {
	registry := NewHandlerRegistry(nil)
	handler1 := &MockToolHandler{
		canHandle: map[string]bool{"tool1": true},
		executeFunc: func(toolName string, args map[string]any, projectRoot string) (string, error) {
			return "handler1", nil
		},
	}
	handler2 := &MockToolHandler{
		canHandle: map[string]bool{"tool2": true},
		executeFunc: func(toolName string, args map[string]any, projectRoot string) (string, error) {
			return "handler2", nil
		},
	}

	registry.Register(handler1)
	registry.Register(handler2)

	result1, _ := registry.Execute("tool1", map[string]any{}, "/project")
	result2, _ := registry.Execute("tool2", map[string]any{}, "/project")

	if result1 != "handler1" {
		t.Errorf("expected handler1 result, got: %s", result1)
	}
	if result2 != "handler2" {
		t.Errorf("expected handler2 result, got: %s", result2)
	}
}


func TestNewBaseHandler(t *testing.T) {
	logger := &MockLogger{}
	handler := NewBaseHandler(logger)

	if handler.Logger != logger {
		t.Error("expected logger to be set")
	}
}

func TestNewBaseHandler_NilLogger(t *testing.T) {
	handler := NewBaseHandler(nil)

	if handler.Logger != nil {
		t.Error("expected nil logger")
	}
}
