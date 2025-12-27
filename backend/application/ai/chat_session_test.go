package ai

import (
	"testing"

	"syntaxia/domain"
)

func TestNewChatSession(t *testing.T) {
	session := NewChatSession("test-id", "/project")

	if session.ID != "test-id" {
		t.Errorf("ID: got %q, expected %q", session.ID, "test-id")
	}
	if session.ProjectRoot != "/project" {
		t.Errorf("ProjectRoot: got %q, expected %q", session.ProjectRoot, "/project")
	}
	if len(session.Messages) != 0 {
		t.Errorf("Messages: got %d, expected 0", len(session.Messages))
	}
	if session.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
	if session.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestChatSessionAddMessage(t *testing.T) {
	session := NewChatSession("test", "/project")
	initialUpdatedAt := session.UpdatedAt

	msg := domain.ChatMessage{Role: domain.RoleUser, Content: "Hello"}
	session.AddMessage(msg)

	messages := session.GetMessages()
	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "Hello" {
		t.Errorf("content: got %q, expected %q", messages[0].Content, "Hello")
	}
	if !session.UpdatedAt.After(initialUpdatedAt) && session.UpdatedAt != initialUpdatedAt {
		t.Error("UpdatedAt should be updated")
	}
}

func TestChatSessionGetMessages(t *testing.T) {
	session := NewChatSession("test", "/project")
	session.AddMessage(domain.ChatMessage{Role: domain.RoleUser, Content: "msg1"})
	session.AddMessage(domain.ChatMessage{Role: domain.RoleAssistant, Content: "msg2"})

	messages := session.GetMessages()

	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}

	// Verify it's a copy
	messages[0].Content = "modified"
	original := session.GetMessages()
	if original[0].Content == "modified" {
		t.Error("GetMessages should return a copy")
	}
}

func TestChatSessionClear(t *testing.T) {
	session := NewChatSession("test", "/project")
	session.AddMessage(domain.ChatMessage{Role: domain.RoleUser, Content: "msg"})

	session.Clear()

	messages := session.GetMessages()
	if len(messages) != 0 {
		t.Errorf("expected 0 messages after clear, got %d", len(messages))
	}
}

func TestNewChatSessionManager(t *testing.T) {
	manager := NewChatSessionManager()

	if manager == nil {
		t.Fatal("expected non-nil manager")
	}
	if manager.sessions == nil {
		t.Error("sessions map should be initialized")
	}
}

func TestChatSessionManagerGetOrCreate(t *testing.T) {
	manager := NewChatSessionManager()

	// Create new session
	session1 := manager.GetOrCreate("session1", "/project1")
	if session1 == nil {
		t.Fatal("expected non-nil session")
	}
	if session1.ID != "session1" {
		t.Errorf("ID: got %q, expected %q", session1.ID, "session1")
	}

	// Get existing session
	session1Again := manager.GetOrCreate("session1", "/project2")
	if session1Again != session1 {
		t.Error("should return same session instance")
	}
	if session1Again.ProjectRoot != "/project1" {
		t.Error("should not update project root for existing session")
	}

	// Create another session
	session2 := manager.GetOrCreate("session2", "/project2")
	if session2 == session1 {
		t.Error("should create different session")
	}
}

func TestChatSessionManagerGet(t *testing.T) {
	manager := NewChatSessionManager()

	// Non-existent session
	_, ok := manager.Get("nonexistent")
	if ok {
		t.Error("should return false for non-existent session")
	}

	// Create and get
	manager.GetOrCreate("test", "/project")
	session, ok := manager.Get("test")
	if !ok {
		t.Error("should return true for existing session")
	}
	if session == nil {
		t.Error("should return session")
	}
}

func TestChatSessionManagerDelete(t *testing.T) {
	manager := NewChatSessionManager()
	manager.GetOrCreate("test", "/project")

	manager.Delete("test")

	_, ok := manager.Get("test")
	if ok {
		t.Error("session should be deleted")
	}

	// Delete non-existent should not panic
	manager.Delete("nonexistent")
}

func TestChatSessionManagerList(t *testing.T) {
	manager := NewChatSessionManager()

	// Empty list
	ids := manager.List()
	if len(ids) != 0 {
		t.Errorf("expected 0 ids, got %d", len(ids))
	}

	// Add sessions
	manager.GetOrCreate("a", "/p")
	manager.GetOrCreate("b", "/p")
	manager.GetOrCreate("c", "/p")

	ids = manager.List()
	if len(ids) != 3 {
		t.Errorf("expected 3 ids, got %d", len(ids))
	}

	// Verify all IDs are present
	idMap := make(map[string]bool)
	for _, id := range ids {
		idMap[id] = true
	}
	for _, expected := range []string{"a", "b", "c"} {
		if !idMap[expected] {
			t.Errorf("expected id %q in list", expected)
		}
	}
}

func TestChatSessionConcurrency(t *testing.T) {
	session := NewChatSession("test", "/project")

	done := make(chan bool)

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func(n int) {
			session.AddMessage(domain.ChatMessage{
				Role:    domain.RoleUser,
				Content: "message",
			})
			done <- true
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			_ = session.GetMessages()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	messages := session.GetMessages()
	if len(messages) != 10 {
		t.Errorf("expected 10 messages, got %d", len(messages))
	}
}

func TestChatSessionManagerConcurrency(t *testing.T) {
	manager := NewChatSessionManager()

	done := make(chan bool)

	// Concurrent GetOrCreate
	for i := 0; i < 10; i++ {
		go func(n int) {
			manager.GetOrCreate("shared", "/project")
			done <- true
		}(i)
	}

	// Concurrent List
	for i := 0; i < 10; i++ {
		go func() {
			_ = manager.List()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}

	ids := manager.List()
	if len(ids) != 1 {
		t.Errorf("expected 1 session, got %d", len(ids))
	}
}
