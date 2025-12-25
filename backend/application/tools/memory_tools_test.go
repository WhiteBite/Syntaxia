package tools

import (
	"syntaxia/domain"
	"testing"
	"time"
)

// MockContextMemory implements domain.ContextMemory for testing
type MockContextMemory struct {
	contexts    []*domain.ConversationContext
	preferences map[string]string
}

func (m *MockContextMemory) SaveContext(ctx *domain.ConversationContext) error {
	ctx.LastAccessed = time.Now()
	m.contexts = append(m.contexts, ctx)
	return nil
}

func (m *MockContextMemory) GetContext(id string) (*domain.ConversationContext, error) {
	for _, c := range m.contexts {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, nil
}

func (m *MockContextMemory) FindContextByTopic(projectRoot, topic string) ([]*domain.ConversationContext, error) {
	var result []*domain.ConversationContext
	for _, c := range m.contexts {
		if c.ProjectRoot == projectRoot && containsStr(c.Topic, topic) {
			result = append(result, c)
		}
	}
	return result, nil
}

func (m *MockContextMemory) GetRecentContexts(projectRoot string, limit int) ([]*domain.ConversationContext, error) {
	var result []*domain.ConversationContext
	for _, c := range m.contexts {
		if c.ProjectRoot == projectRoot {
			result = append(result, c)
			if len(result) >= limit {
				break
			}
		}
	}
	return result, nil
}

func (m *MockContextMemory) SetPreference(key, value string) error {
	if m.preferences == nil {
		m.preferences = make(map[string]string)
	}
	m.preferences[key] = value
	return nil
}

func (m *MockContextMemory) GetPreference(key string) (string, error) {
	if m.preferences == nil {
		return "", nil
	}
	return m.preferences[key], nil
}

func (m *MockContextMemory) GetAllPreferences() (map[string]string, error) {
	if m.preferences == nil {
		return make(map[string]string), nil
	}
	return m.preferences, nil
}

func (m *MockContextMemory) Close() error {
	return nil
}

func TestSaveContext_Success(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewMemoryToolsHandler(nil, mock)

	result, err := handler.Execute("save_context", map[string]any{
		"topic":   "auth-feature",
		"summary": "Authentication implementation",
		"files":   []any{"auth.go", "login.go"},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "saved") {
		t.Errorf("expected success message, got: %s", result)
	}
	if len(mock.contexts) != 1 {
		t.Errorf("expected 1 context saved, got: %d", len(mock.contexts))
	}
}

func TestFindContext_Found(t *testing.T) {
	mock := &MockContextMemory{
		contexts: []*domain.ConversationContext{
			{ProjectRoot: "/project", Topic: "auth-feature", Summary: "Auth stuff", Files: []string{"auth.go"}},
		},
	}
	handler := NewMemoryToolsHandler(nil, mock)

	result, err := handler.Execute("find_context", map[string]any{"topic": "auth"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "auth-feature") {
		t.Errorf("expected to find context, got: %s", result)
	}
}

func TestFindContext_NotFound(t *testing.T) {
	mock := &MockContextMemory{}
	handler := NewMemoryToolsHandler(nil, mock)

	result, err := handler.Execute("find_context", map[string]any{"topic": "nonexistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "No contexts found") {
		t.Errorf("expected 'not found' message, got: %s", result)
	}
}

func TestGetRecentContexts_Success(t *testing.T) {
	mock := &MockContextMemory{
		contexts: []*domain.ConversationContext{
			{ProjectRoot: "/project", Topic: "topic1", Summary: "Summary 1", LastAccessed: time.Now()},
			{ProjectRoot: "/project", Topic: "topic2", Summary: "Summary 2", LastAccessed: time.Now()},
		},
	}
	handler := NewMemoryToolsHandler(nil, mock)

	result, err := handler.Execute("get_recent_contexts", map[string]any{"limit": float64(10)}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "topic1") || !containsStr(result, "topic2") {
		t.Errorf("expected recent contexts, got: %s", result)
	}
}

func TestMemoryTools_NotInitialized(t *testing.T) {
	handler := NewMemoryToolsHandler(nil, nil)

	_, err := handler.Execute("save_context", map[string]any{"topic": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error when context memory not initialized")
	}
}
