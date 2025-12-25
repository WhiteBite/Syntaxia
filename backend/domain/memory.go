package domain

import "time"

// ContextMemory interface for storing conversation context
type ContextMemory interface {
	SaveContext(ctx *ConversationContext) error
	GetContext(id string) (*ConversationContext, error)
	FindContextByTopic(projectRoot, topic string) ([]*ConversationContext, error)
	GetRecentContexts(projectRoot string, limit int) ([]*ConversationContext, error)
	SetPreference(key, value string) error
	GetPreference(key string) (string, error)
	GetAllPreferences() (map[string]string, error)
	Close() error
}

// ConversationContext represents stored context for a conversation
type ConversationContext struct {
	ID           string    `json:"id"`
	ProjectRoot  string    `json:"projectRoot"`
	Topic        string    `json:"topic"`
	Files        []string  `json:"files"`
	Symbols      []string  `json:"symbols"`
	Summary      string    `json:"summary"`
	LastAccessed time.Time `json:"lastAccessed"`
	CreatedAt    time.Time `json:"createdAt"`
	MessageCount int       `json:"messageCount"`
}

// UserPreference stores user preferences
type UserPreference struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	UpdatedAt time.Time `json:"updatedAt"`
}
