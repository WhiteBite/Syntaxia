// Package changes provides change management with history, undo/redo, and conflict detection.
package changes

import (
	"time"
)

// ChangeOperation represents the type of change operation
type ChangeOperation int

const (
	// OpCreate indicates a new file creation
	OpCreate ChangeOperation = iota
	// OpModify indicates a file modification
	OpModify
	// OpDelete indicates a file deletion
	OpDelete
)

// String returns string representation of ChangeOperation
func (op ChangeOperation) String() string {
	switch op {
	case OpCreate:
		return "create"
	case OpModify:
		return "modify"
	case OpDelete:
		return "delete"
	default:
		return "unknown"
	}
}

// Change represents a single file change with metadata
type Change struct {
	ID              string          `json:"id"`
	Path            string          `json:"path"`
	Operation       ChangeOperation `json:"operation"`
	Content         []byte          `json:"content,omitempty"`
	OriginalContent []byte          `json:"originalContent,omitempty"`
	OriginalHash    string          `json:"originalHash,omitempty"`
	Timestamp       time.Time       `json:"timestamp"`
	SessionID       string          `json:"sessionId"`
	Applied         bool            `json:"applied"`
}

// ChangeGroup represents a group of related changes (e.g., from one AI operation)
type ChangeGroup struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"sessionId"`
	Description string    `json:"description"`
	Changes     []*Change `json:"changes"`
	CreatedAt   time.Time `json:"createdAt"`
	Applied     bool      `json:"applied"`
}

// ConflictInfo contains information about a detected conflict
type ConflictInfo struct {
	Path            string    `json:"path"`
	ExpectedHash    string    `json:"expectedHash"`
	CurrentHash     string    `json:"currentHash"`
	ChangeID        string    `json:"changeId"`
	DetectedAt      time.Time `json:"detectedAt"`
	OriginalContent []byte    `json:"originalContent,omitempty"`
	CurrentContent  []byte    `json:"currentContent,omitempty"`
}

// HistoryEntry represents an entry in the undo/redo history
type HistoryEntry struct {
	ID          string    `json:"id"`
	GroupID     string    `json:"groupId"`
	Type        string    `json:"type"` // "apply" or "rollback"
	Changes     []*Change `json:"changes"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Undone      bool      `json:"undone"`
}

// ManagerConfig contains configuration for ChangeManager
type ManagerConfig struct {
	MaxHistorySize     int           `json:"maxHistorySize"`
	ConflictCheckDelay time.Duration `json:"conflictCheckDelay"`
	EnableAutoSave     bool          `json:"enableAutoSave"`
}

// DefaultManagerConfig returns default configuration
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		MaxHistorySize:     100,
		ConflictCheckDelay: 500 * time.Millisecond,
		EnableAutoSave:     false,
	}
}
