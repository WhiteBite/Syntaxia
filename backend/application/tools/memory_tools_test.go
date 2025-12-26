package tools

import (
	"errors"
	"strings"
	"testing"
	"time"

	"syntaxia/domain"
)

// mockMemoryForMemoryTools implements domain.ContextMemory for memory_tools testing
type mockMemoryForMemoryTools struct {
	contexts       map[string]*domain.ConversationContext
	preferences    map[string]string
	saveError      error
	findError      error
	getRecentError error
	savedContexts  []*domain.ConversationContext
}

func newMockMemoryForMemoryTools() *mockMemoryForMemoryTools {
	return &mockMemoryForMemoryTools{
		contexts:      make(map[string]*domain.ConversationContext),
		preferences:   make(map[string]string),
		savedContexts: make([]*domain.ConversationContext, 0),
	}
}

func (m *mockMemoryForMemoryTools) SaveContext(ctx *domain.ConversationContext) error {
	if m.saveError != nil {
		return m.saveError
	}
	m.savedContexts = append(m.savedContexts, ctx)
	m.contexts[ctx.Topic] = ctx
	return nil
}

func (m *mockMemoryForMemoryTools) GetContext(id string) (*domain.ConversationContext, error) {
	if ctx, ok := m.contexts[id]; ok {
		return ctx, nil
	}
	return nil, errors.New("context not found")
}

func (m *mockMemoryForMemoryTools) FindContextByTopic(projectRoot, topic string) ([]*domain.ConversationContext, error) {
	if m.findError != nil {
		return nil, m.findError
	}
	var results []*domain.ConversationContext
	for _, ctx := range m.contexts {
		if ctx.ProjectRoot == projectRoot && strings.Contains(ctx.Topic, topic) {
			results = append(results, ctx)
		}
	}
	return results, nil
}

func (m *mockMemoryForMemoryTools) GetRecentContexts(projectRoot string, limit int) ([]*domain.ConversationContext, error) {
	if m.getRecentError != nil {
		return nil, m.getRecentError
	}
	var results []*domain.ConversationContext
	for _, ctx := range m.contexts {
		if ctx.ProjectRoot == projectRoot {
			results = append(results, ctx)
			if len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (m *mockMemoryForMemoryTools) SetPreference(key, value string) error {
	if m.preferences == nil {
		m.preferences = make(map[string]string)
	}
	m.preferences[key] = value
	return nil
}

func (m *mockMemoryForMemoryTools) GetPreference(key string) (string, error) {
	if m.preferences == nil {
		return "", errors.New("not found")
	}
	if v, ok := m.preferences[key]; ok {
		return v, nil
	}
	return "", errors.New("not found")
}

func (m *mockMemoryForMemoryTools) GetAllPreferences() (map[string]string, error) {
	if m.preferences == nil {
		return make(map[string]string), nil
	}
	return m.preferences, nil
}

func (m *mockMemoryForMemoryTools) Close() error { return nil }

// containsStr helper for checking substrings
func containsStr(result, substr string) bool {
	return strings.Contains(result, substr)
}

// === Constructor Tests ===

func TestNewMemoryToolsHandler(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.ContextMemory != mockMemory {
		t.Error("expected ContextMemory to be set")
	}
}

// === CanHandle Tests ===

func TestMemoryToolsHandler_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     bool
	}{
		{"save_context", "save_context", true},
		{"find_context", "find_context", true},
		{"get_recent_contexts", "get_recent_contexts", true},
		{"unknown_tool", "unknown_tool", false},
		{"file_tool", "read_file", false},
		{"git_tool", "git_status", false},
		{"empty_string", "", false},
	}

	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.CanHandle(tt.toolName); got != tt.want {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

// === GetTools Tests ===

func TestMemoryToolsHandler_GetTools(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())
	tools := handler.GetTools()

	if len(tools) == 0 {
		t.Fatal("expected non-empty tools list")
	}

	expectedTools := []string{"save_context", "find_context", "get_recent_contexts"}
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q not found", expected)
		}
	}
}

func TestMemoryToolsHandler_GetTools_HasRequiredParameters(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())
	tools := handler.GetTools()

	for _, tool := range tools {
		if tool.Name == "save_context" {
			if len(tool.Parameters.Required) == 0 {
				t.Error("save_context should have required parameters")
			}
			found := false
			for _, req := range tool.Parameters.Required {
				if req == "topic" {
					found = true
					break
				}
			}
			if !found {
				t.Error("save_context should require 'topic' parameter")
			}
		}
		if tool.Name == "find_context" {
			found := false
			for _, req := range tool.Parameters.Required {
				if req == "topic" {
					found = true
					break
				}
			}
			if !found {
				t.Error("find_context should require 'topic' parameter")
			}
		}
	}
}

// === Execute Tests ===

func TestMemoryToolsHandler_Execute_UnknownTool(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())

	_, err := handler.Execute("unknown_tool", map[string]any{}, "/test")
	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown memory tool") {
		t.Errorf("expected 'unknown memory tool' error, got: %v", err)
	}
}

// === SaveContext Tests ===

func TestSaveContext_Success(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("save_context", map[string]any{
		"topic":   "authentication",
		"summary": "Auth module implementation",
		"files":   []any{"auth.go", "login.go"},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Context saved") {
		t.Errorf("expected success message, got: %s", result)
	}
	if !strings.Contains(result, "authentication") {
		t.Errorf("expected topic in result, got: %s", result)
	}
	if len(mockMemory.savedContexts) != 1 {
		t.Errorf("expected 1 saved context, got: %d", len(mockMemory.savedContexts))
	}
}

func TestSaveContext_WithoutOptionalFields(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("save_context", map[string]any{
		"topic": "minimal",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Context saved") {
		t.Errorf("expected success message, got: %s", result)
	}
}

func TestSaveContext_ContextMemoryNotInitialized(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, nil)

	_, err := handler.Execute("save_context", map[string]any{
		"topic": "test",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected 'not initialized' error, got: %v", err)
	}
}

func TestSaveContext_SaveError(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.saveError = errors.New("database error")
	handler := NewMemoryToolsHandler(nil, mockMemory)

	_, err := handler.Execute("save_context", map[string]any{
		"topic": "test",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when save fails")
	}
	if !strings.Contains(err.Error(), "failed to save context") {
		t.Errorf("expected 'failed to save context' error, got: %v", err)
	}
}

func TestSaveContext_FilesConversion(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	_, err := handler.Execute("save_context", map[string]any{
		"topic": "files-test",
		"files": []any{"file1.go", "file2.go", "file3.go"},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mockMemory.savedContexts) != 1 {
		t.Fatal("expected 1 saved context")
	}
	if len(mockMemory.savedContexts[0].Files) != 3 {
		t.Errorf("expected 3 files, got: %d", len(mockMemory.savedContexts[0].Files))
	}
}

func TestSaveContext_MixedFilesArray(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Test with mixed types in files array (only strings should be added)
	_, err := handler.Execute("save_context", map[string]any{
		"topic": "mixed-test",
		"files": []any{"file1.go", 123, "file2.go", nil},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mockMemory.savedContexts[0].Files) != 2 {
		t.Errorf("expected 2 string files, got: %d", len(mockMemory.savedContexts[0].Files))
	}
}

// === FindContext Tests ===

func TestFindContext_Success(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["auth-topic"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "auth-topic",
		Summary:     "Authentication implementation",
		Files:       []string{"auth.go", "login.go"},
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("find_context", map[string]any{
		"topic": "auth",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "auth-topic") {
		t.Errorf("expected to find context, got: %s", result)
	}
	if !strings.Contains(result, "Authentication implementation") {
		t.Errorf("expected summary in result, got: %s", result)
	}
}

func TestFindContext_NoResults(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("find_context", map[string]any{
		"topic": "nonexistent",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No contexts found") {
		t.Errorf("expected 'No contexts found' message, got: %s", result)
	}
}

func TestFindContext_MissingTopic(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())

	_, err := handler.Execute("find_context", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing topic")
	}
	if !strings.Contains(err.Error(), "topic is required") {
		t.Errorf("expected 'topic is required' error, got: %v", err)
	}
}

func TestFindContext_EmptyTopic(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, newMockMemoryForMemoryTools())

	_, err := handler.Execute("find_context", map[string]any{
		"topic": "",
	}, "/project")

	if err == nil {
		t.Fatal("expected error for empty topic")
	}
	if !strings.Contains(err.Error(), "topic is required") {
		t.Errorf("expected 'topic is required' error, got: %v", err)
	}
}

func TestFindContext_ContextMemoryNotInitialized(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, nil)

	_, err := handler.Execute("find_context", map[string]any{
		"topic": "test",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected 'not initialized' error, got: %v", err)
	}
}

func TestFindContext_FindError(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.findError = errors.New("database error")
	handler := NewMemoryToolsHandler(nil, mockMemory)

	_, err := handler.Execute("find_context", map[string]any{
		"topic": "test",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when find fails")
	}
	if !strings.Contains(err.Error(), "failed to find context") {
		t.Errorf("expected 'failed to find context' error, got: %v", err)
	}
}

func TestFindContext_MultipleResults(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["auth-login"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "auth-login",
		Summary:     "Login flow",
		Files:       []string{"login.go"},
	}
	mockMemory.contexts["auth-logout"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "auth-logout",
		Summary:     "Logout flow",
		Files:       []string{"logout.go", "session.go"},
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("find_context", map[string]any{
		"topic": "auth",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "auth-login") || !strings.Contains(result, "auth-logout") {
		t.Errorf("expected both contexts in result, got: %s", result)
	}
}

// === GetRecentContexts Tests ===

func TestGetRecentContexts_Success(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["recent1"] = &domain.ConversationContext{
		ProjectRoot:  "/project",
		Topic:        "recent1",
		Summary:      "Recent context 1",
		LastAccessed: time.Now(),
	}
	mockMemory.contexts["recent2"] = &domain.ConversationContext{
		ProjectRoot:  "/project",
		Topic:        "recent2",
		Summary:      "Recent context 2",
		LastAccessed: time.Now().Add(-time.Hour),
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Recent contexts") {
		t.Errorf("expected 'Recent contexts' header, got: %s", result)
	}
}

func TestGetRecentContexts_WithLimit(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	for i := 0; i < 5; i++ {
		mockMemory.contexts[string(rune('a'+i))] = &domain.ConversationContext{
			ProjectRoot:  "/project",
			Topic:        string(rune('a' + i)),
			Summary:      "Context",
			LastAccessed: time.Now(),
		}
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("get_recent_contexts", map[string]any{
		"limit": float64(2),
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Recent contexts") {
		t.Errorf("expected recent contexts, got: %s", result)
	}
}

func TestGetRecentContexts_DefaultLimit(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Should use default limit of 10 when not specified
	_, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetRecentContexts_NoResults(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No recent contexts") {
		t.Errorf("expected 'No recent contexts' message, got: %s", result)
	}
}

func TestGetRecentContexts_ContextMemoryNotInitialized(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, nil)

	_, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected 'not initialized' error, got: %v", err)
	}
}

func TestGetRecentContexts_GetRecentError(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.getRecentError = errors.New("database error")
	handler := NewMemoryToolsHandler(nil, mockMemory)

	_, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when get recent fails")
	}
	if !strings.Contains(err.Error(), "failed to get recent contexts") {
		t.Errorf("expected 'failed to get recent contexts' error, got: %v", err)
	}
}

func TestGetRecentContexts_DifferentProjectRoot(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["ctx1"] = &domain.ConversationContext{
		ProjectRoot:  "/project1",
		Topic:        "ctx1",
		Summary:      "Project 1 context",
		LastAccessed: time.Now(),
	}
	mockMemory.contexts["ctx2"] = &domain.ConversationContext{
		ProjectRoot:  "/project2",
		Topic:        "ctx2",
		Summary:      "Project 2 context",
		LastAccessed: time.Now(),
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "ctx2") {
		t.Errorf("should not include contexts from different project, got: %s", result)
	}
}

// === Edge Cases ===

func TestSaveContext_ProjectRootPreserved(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	_, err := handler.Execute("save_context", map[string]any{
		"topic": "test",
	}, "/my/project/path")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mockMemory.savedContexts) != 1 {
		t.Fatal("expected 1 saved context")
	}
	if mockMemory.savedContexts[0].ProjectRoot != "/my/project/path" {
		t.Errorf("expected project root to be preserved, got: %s", mockMemory.savedContexts[0].ProjectRoot)
	}
}

func TestFindContext_OutputFormat(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["test-topic"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "test-topic",
		Summary:     "Test summary",
		Files:       []string{"a.go", "b.go", "c.go"},
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("find_context", map[string]any{
		"topic": "test",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Check output contains file count
	if !strings.Contains(result, "3 files") {
		t.Errorf("expected file count in output, got: %s", result)
	}
}

func TestGetRecentContexts_DateFormat(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	testDate := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	mockMemory.contexts["dated"] = &domain.ConversationContext{
		ProjectRoot:  "/project",
		Topic:        "dated",
		Summary:      "Dated context",
		LastAccessed: testDate,
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("get_recent_contexts", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Check date format (YYYY-MM-DD)
	if !strings.Contains(result, "2024-06-15") {
		t.Errorf("expected formatted date in output, got: %s", result)
	}
}


// === Additional Edge Case Tests ===

func TestSaveContext_EmptyTopic(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Empty topic should still work (topic is extracted as empty string)
	result, err := handler.Execute("save_context", map[string]any{
		"topic": "",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Context saved") {
		t.Errorf("expected success message, got: %s", result)
	}
}

func TestSaveContext_AllFields(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("save_context", map[string]any{
		"topic":   "full-context",
		"summary": "Complete context with all fields",
		"files":   []any{"file1.go", "file2.go", "file3.go"},
	}, "/my/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "full-context") {
		t.Errorf("expected topic in result, got: %s", result)
	}

	// Verify saved context
	if len(mockMemory.savedContexts) != 1 {
		t.Fatal("expected 1 saved context")
	}
	ctx := mockMemory.savedContexts[0]
	if ctx.Topic != "full-context" {
		t.Errorf("expected topic 'full-context', got: %s", ctx.Topic)
	}
	if ctx.Summary != "Complete context with all fields" {
		t.Errorf("expected summary, got: %s", ctx.Summary)
	}
	if len(ctx.Files) != 3 {
		t.Errorf("expected 3 files, got: %d", len(ctx.Files))
	}
	if ctx.ProjectRoot != "/my/project" {
		t.Errorf("expected project root '/my/project', got: %s", ctx.ProjectRoot)
	}
}

func TestFindContext_PartialMatch(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["authentication-login"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "authentication-login",
		Summary:     "Login implementation",
		Files:       []string{"login.go"},
	}
	mockMemory.contexts["authentication-logout"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "authentication-logout",
		Summary:     "Logout implementation",
		Files:       []string{"logout.go"},
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("find_context", map[string]any{
		"topic": "authentication",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should find both contexts with partial match
	if !strings.Contains(result, "authentication-login") || !strings.Contains(result, "authentication-logout") {
		t.Errorf("expected both contexts, got: %s", result)
	}
}

func TestGetRecentContexts_LimitAsInt(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["ctx1"] = &domain.ConversationContext{
		ProjectRoot:  "/project",
		Topic:        "ctx1",
		Summary:      "Context 1",
		LastAccessed: time.Now(),
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Test with integer limit (not float64)
	result, err := handler.Execute("get_recent_contexts", map[string]any{
		"limit": 5, // int instead of float64
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "ctx1") {
		t.Errorf("expected context in result, got: %s", result)
	}
}

func TestSaveContext_EmptyFilesArray(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("save_context", map[string]any{
		"topic": "no-files",
		"files": []any{},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Context saved") {
		t.Errorf("expected success message, got: %s", result)
	}
	if len(mockMemory.savedContexts[0].Files) != 0 {
		t.Errorf("expected empty files array")
	}
}

func TestSaveContext_NilFilesValue(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	handler := NewMemoryToolsHandler(nil, mockMemory)

	result, err := handler.Execute("save_context", map[string]any{
		"topic": "nil-files",
		"files": nil,
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Context saved") {
		t.Errorf("expected success message, got: %s", result)
	}
}

func TestFindContext_CaseSensitive(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["AUTH"] = &domain.ConversationContext{
		ProjectRoot: "/project",
		Topic:       "AUTH",
		Summary:     "Auth context",
		Files:       []string{},
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Search with lowercase - depends on mock implementation
	result, err := handler.Execute("find_context", map[string]any{
		"topic": "AUTH",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "AUTH") {
		t.Errorf("expected to find AUTH context, got: %s", result)
	}
}

func TestGetRecentContexts_ZeroLimit(t *testing.T) {
	mockMemory := newMockMemoryForMemoryTools()
	mockMemory.contexts["ctx"] = &domain.ConversationContext{
		ProjectRoot:  "/project",
		Topic:        "ctx",
		Summary:      "Context",
		LastAccessed: time.Now(),
	}
	handler := NewMemoryToolsHandler(nil, mockMemory)

	// Zero limit should use default (10)
	result, err := handler.Execute("get_recent_contexts", map[string]any{
		"limit": float64(0),
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should still return results with default limit
	if strings.Contains(result, "No recent contexts") {
		t.Errorf("expected contexts with default limit, got: %s", result)
	}
}


func TestNewMemoryToolsHandler_NilMemory(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, nil)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.ContextMemory != nil {
		t.Error("expected nil context memory")
	}
}
