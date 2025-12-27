package ai

import (
	"sync"
	"time"

	"syntaxia/domain"
)

// ChatSession manages a chat session with message history
type ChatSession struct {
	ID          string               `json:"id"`
	Messages    []domain.ChatMessage `json:"messages"`
	ProjectRoot string               `json:"projectRoot"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	mu          sync.RWMutex
}

// NewChatSession creates a new chat session
func NewChatSession(id, projectRoot string) *ChatSession {
	now := time.Now()
	return &ChatSession{
		ID:          id,
		Messages:    make([]domain.ChatMessage, 0),
		ProjectRoot: projectRoot,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddMessage adds a message to the session
func (s *ChatSession) AddMessage(msg domain.ChatMessage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = append(s.Messages, msg)
	s.UpdatedAt = time.Now()
}

// GetMessages returns all messages
func (s *ChatSession) GetMessages() []domain.ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()
	msgs := make([]domain.ChatMessage, len(s.Messages))
	copy(msgs, s.Messages)
	return msgs
}

// Clear clears the session messages
func (s *ChatSession) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Messages = make([]domain.ChatMessage, 0)
	s.UpdatedAt = time.Now()
}

// ChatSessionManager manages multiple chat sessions
type ChatSessionManager struct {
	sessions map[string]*ChatSession
	mu       sync.RWMutex
}

// NewChatSessionManager creates a new session manager
func NewChatSessionManager() *ChatSessionManager {
	return &ChatSessionManager{
		sessions: make(map[string]*ChatSession),
	}
}

// GetOrCreate gets or creates a session
func (m *ChatSessionManager) GetOrCreate(id, projectRoot string) *ChatSession {
	m.mu.Lock()
	defer m.mu.Unlock()

	if session, ok := m.sessions[id]; ok {
		return session
	}

	session := NewChatSession(id, projectRoot)
	m.sessions[id] = session
	return session
}

// Get gets a session by ID
func (m *ChatSessionManager) Get(id string) (*ChatSession, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[id]
	return session, ok
}

// Delete deletes a session
func (m *ChatSessionManager) Delete(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, id)
}

// List returns all session IDs
func (m *ChatSessionManager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.sessions))
	for id := range m.sessions {
		ids = append(ids, id)
	}
	return ids
}
