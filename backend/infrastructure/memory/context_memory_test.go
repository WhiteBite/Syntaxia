package memory

import (
	"syntaxia/domain"
	"testing"
	"time"
)

const testValue2 = "value2"

func TestNewContextMemory(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	if cm == nil {
		t.Fatal("NewContextMemory returned nil")
	}
	if cm.db == nil {
		t.Error("db not initialized")
	}
}

func TestContextMemory_SaveAndGetContext(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	ctx := &domain.ConversationContext{
		ID:           "test-ctx-1",
		ProjectRoot:  "/test/project",
		Topic:        "auth implementation",
		Files:        []string{"auth.go", "auth_test.go"},
		Summary:      "Working on authentication",
		LastAccessed: time.Now(),
		CreatedAt:    time.Now(),
		MessageCount: 5,
	}

	err = cm.SaveContext(ctx)
	if err != nil {
		t.Fatalf("SaveContext failed: %v", err)
	}

	retrieved, err := cm.GetContext("test-ctx-1")
	if err != nil {
		t.Fatalf("GetContext failed: %v", err)
	}

	if retrieved.ID != ctx.ID {
		t.Errorf("ID mismatch: got %q, want %q", retrieved.ID, ctx.ID)
	}
	if retrieved.Topic != ctx.Topic {
		t.Errorf("Topic mismatch: got %q, want %q", retrieved.Topic, ctx.Topic)
	}
	if len(retrieved.Files) != len(ctx.Files) {
		t.Errorf("Files count mismatch: got %d, want %d", len(retrieved.Files), len(ctx.Files))
	}
}

func TestContextMemory_FindContextByTopic(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	// Save multiple contexts
	contexts := []*domain.ConversationContext{
		{ID: "ctx-1", ProjectRoot: "/project", Topic: "auth login", LastAccessed: time.Now(), CreatedAt: time.Now()},
		{ID: "ctx-2", ProjectRoot: "/project", Topic: "auth logout", LastAccessed: time.Now(), CreatedAt: time.Now()},
		{ID: "ctx-3", ProjectRoot: "/project", Topic: "user profile", LastAccessed: time.Now(), CreatedAt: time.Now()},
	}

	for _, ctx := range contexts {
		cm.SaveContext(ctx)
	}

	found, err := cm.FindContextByTopic("/project", "auth")
	if err != nil {
		t.Fatalf("FindContextByTopic failed: %v", err)
	}

	if len(found) < 2 {
		t.Errorf("expected at least 2 contexts with 'auth', got %d", len(found))
	}
}

func TestContextMemory_GetRecentContexts(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	// Save contexts
	for i := 0; i < 5; i++ {
		ctx := &domain.ConversationContext{
			ID:           "ctx-" + string(rune('a'+i)),
			ProjectRoot:  "/project",
			Topic:        "topic " + string(rune('a'+i)),
			LastAccessed: time.Now(),
			CreatedAt:    time.Now(),
		}
		cm.SaveContext(ctx)
	}

	recent, err := cm.GetRecentContexts("/project", 3)
	if err != nil {
		t.Fatalf("GetRecentContexts failed: %v", err)
	}

	if len(recent) > 3 {
		t.Errorf("expected max 3 contexts, got %d", len(recent))
	}
}

func TestContextMemory_Preferences(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	// Set preferences
	err = cm.SetPreference("exclude_tests", "true")
	if err != nil {
		t.Fatalf("SetPreference failed: %v", err)
	}

	err = cm.SetPreference("max_files", "50")
	if err != nil {
		t.Fatalf("SetPreference failed: %v", err)
	}

	// Get single preference
	val, err := cm.GetPreference("exclude_tests")
	if err != nil {
		t.Fatalf("GetPreference failed: %v", err)
	}
	if val != boolTrue {
		t.Errorf("expected 'true', got %q", val)
	}

	// Get all preferences
	prefs, err := cm.GetAllPreferences()
	if err != nil {
		t.Fatalf("GetAllPreferences failed: %v", err)
	}

	if len(prefs) < 2 {
		t.Errorf("expected at least 2 preferences, got %d", len(prefs))
	}
	if prefs["max_files"] != "50" {
		t.Errorf("max_files should be '50', got %q", prefs["max_files"])
	}
}

func TestContextMemory_UpdatePreference(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}
	defer cm.Close()

	cm.SetPreference("key", "value1")
	cm.SetPreference("key", testValue2)

	val, _ := cm.GetPreference("key")
	if val != testValue2 {
		t.Errorf("expected 'value2', got %q", val)
	}
}

func TestExtractTopicFromMessage(t *testing.T) {
	tests := []struct {
		message string
		expect  string
	}{
		{"work on auth module", "auth module"},
		{"fix the login bug", "login bug"},
		{"implement user registration", "user registration"},
		{"random message here", "random message here"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			topic := ExtractTopicFromMessage(tt.message)
			if topic == "" {
				t.Error("expected non-empty topic")
			}
		})
	}
}

func TestContextMemory_Close(t *testing.T) {
	tmpDir := t.TempDir()
	cm, err := NewContextMemory(tmpDir)
	if err != nil {
		t.Fatalf("NewContextMemory failed: %v", err)
	}

	err = cm.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
}
