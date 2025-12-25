package tools

import (
	"testing"
)

func TestPreferencesHandler_CanHandle(t *testing.T) {
	handler := NewPreferencesToolsHandler(nil, nil)

	tests := []struct {
		toolName string
		expected bool
	}{
		{"set_preference", true},
		{"get_preferences", true},
		{"unknown_tool", false},
		{"save_context", false},
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

func TestPreferencesHandler_GetTools(t *testing.T) {
	handler := NewPreferencesToolsHandler(nil, nil)
	tools := handler.GetTools()

	if len(tools) != 2 {
		t.Errorf("expected 2 tools, got %d", len(tools))
	}

	expectedTools := map[string]bool{
		"set_preference":  false,
		"get_preferences": false,
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

func TestSetPreference_Success(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewPreferencesToolsHandler(nil, mock)

	result, err := handler.Execute("set_preference", map[string]any{
		"key":   "theme",
		"value": "dark",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "theme") || !containsStr(result, "dark") {
		t.Errorf("expected success message with key/value, got: %s", result)
	}
	if mock.preferences["theme"] != "dark" {
		t.Errorf("preference not saved correctly")
	}
}

func TestSetPreference_MissingKey(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewPreferencesToolsHandler(nil, mock)

	_, err := handler.Execute("set_preference", map[string]any{
		"value": "dark",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when key missing")
	}
}

func TestSetPreference_NotInitialized(t *testing.T) {
	handler := NewPreferencesToolsHandler(nil, nil)

	_, err := handler.Execute("set_preference", map[string]any{
		"key":   "theme",
		"value": "dark",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
}

func TestGetPreferences_SingleKey(t *testing.T) {
	mock := &MockContextMemory{
		preferences: map[string]string{"theme": "dark", "lang": "en"},
	}
	handler := NewPreferencesToolsHandler(nil, mock)

	result, err := handler.Execute("get_preferences", map[string]any{
		"key": "theme",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "theme") || !containsStr(result, "dark") {
		t.Errorf("expected preference value, got: %s", result)
	}
}

func TestGetPreferences_AllPreferences(t *testing.T) {
	mock := &MockContextMemory{
		preferences: map[string]string{"theme": "dark", "lang": "en"},
	}
	handler := NewPreferencesToolsHandler(nil, mock)

	result, err := handler.Execute("get_preferences", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "theme") || !containsStr(result, "lang") {
		t.Errorf("expected all preferences, got: %s", result)
	}
}

func TestGetPreferences_NotFound(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewPreferencesToolsHandler(nil, mock)

	result, err := handler.Execute("get_preferences", map[string]any{
		"key": "nonexistent",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "not found") {
		t.Errorf("expected 'not found' message, got: %s", result)
	}
}

func TestGetPreferences_Empty(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewPreferencesToolsHandler(nil, mock)

	result, err := handler.Execute("get_preferences", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "No preferences") {
		t.Errorf("expected 'no preferences' message, got: %s", result)
	}
}

func TestGetPreferences_NotInitialized(t *testing.T) {
	handler := NewPreferencesToolsHandler(nil, nil)

	_, err := handler.Execute("get_preferences", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
}

func TestPreferencesHandler_UnknownTool(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewPreferencesToolsHandler(nil, mock)

	_, err := handler.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
}
