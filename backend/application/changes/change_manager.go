package changes

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"syntaxia/domain"
)

// SessionStats contains statistics about a session
type SessionStats struct {
	FilesChanged   int  `json:"filesChanged"`
	LinesAdded     int  `json:"linesAdded"`
	LinesRemoved   int  `json:"linesRemoved"`
	CanUndo        bool `json:"canUndo"`
	CanRedo        bool `json:"canRedo"`
	HasConflicts   bool `json:"hasConflicts"`
	PendingChanges int  `json:"pendingChanges"`
	AppliedChanges int  `json:"appliedChanges"`
	GroupCount     int  `json:"groupCount"`
}

// Manager manages file changes with history, undo/redo, and conflict detection
type Manager struct {
	mu               sync.RWMutex
	sandbox          domain.SandboxFS
	realFS           domain.FileSystemProvider
	log              domain.Logger
	config           *ManagerConfig
	currentGroup     *ChangeGroup
	groups           map[string]*ChangeGroup
	history          []*HistoryEntry
	historyIndex     int
	sessionID        string
	fileHashes       map[string]string // path -> hash at time of read
	changeCounter    int
	groupCounter     int
	conflictDetector *ConflictDetector
}

// NewManager creates a new change manager
func NewManager(
	sandbox domain.SandboxFS,
	realFS domain.FileSystemProvider,
	log domain.Logger,
	config *ManagerConfig,
) *Manager {
	if config == nil {
		config = DefaultManagerConfig()
	}
	return &Manager{
		sandbox:          sandbox,
		realFS:           realFS,
		log:              log,
		config:           config,
		groups:           make(map[string]*ChangeGroup),
		history:          make([]*HistoryEntry, 0),
		historyIndex:     -1,
		sessionID:        generateID("session"),
		fileHashes:       make(map[string]string),
		conflictDetector: NewConflictDetector(realFS, log),
	}
}

// StartGroup begins a new change group for related changes
func (m *Manager) StartGroup(description string) string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.groupCounter++
	groupID := fmt.Sprintf("group-%d-%d", time.Now().UnixNano(), m.groupCounter)

	m.currentGroup = &ChangeGroup{
		ID:          groupID,
		SessionID:   m.sessionID,
		Description: description,
		Changes:     make([]*Change, 0),
		CreatedAt:   time.Now(),
		Applied:     false,
	}

	m.groups[groupID] = m.currentGroup
	m.log.Info(fmt.Sprintf("Started change group: %s (%s)", groupID, description))

	return groupID
}

// EndGroup finalizes the current change group
func (m *Manager) EndGroup() (*ChangeGroup, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.currentGroup == nil {
		return nil, fmt.Errorf("no active change group")
	}

	group := m.currentGroup
	m.currentGroup = nil

	m.log.Info(fmt.Sprintf("Ended change group: %s with %d changes", group.ID, len(group.Changes)))

	return group, nil
}

// RecordChange records a file change to the current group
func (m *Manager) RecordChange(path string, content []byte, operation ChangeOperation) (*Change, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var originalContent []byte
	var originalHash string

	if operation != OpCreate {
		data, err := m.realFS.ReadFile(path)
		if err == nil {
			originalContent = data
			originalHash = computeHash(data)
			m.fileHashes[path] = originalHash
		}
	}

	m.changeCounter++
	change := &Change{
		ID:              fmt.Sprintf("change-%d-%d", time.Now().UnixNano(), m.changeCounter),
		Path:            path,
		Operation:       operation,
		Content:         content,
		OriginalContent: originalContent,
		OriginalHash:    originalHash,
		Timestamp:       time.Now(),
		SessionID:       m.sessionID,
		Applied:         false,
	}

	if m.currentGroup != nil {
		m.currentGroup.Changes = append(m.currentGroup.Changes, change)
	}

	m.log.Info(fmt.Sprintf("Recorded change: %s %s", operation.String(), path))

	return change, nil
}

// CheckConflict checks if a file has been modified since it was read
func (m *Manager) CheckConflict(path string) (*ConflictInfo, error) {
	m.mu.RLock()
	expectedHash, hasHash := m.fileHashes[path]
	m.mu.RUnlock()

	if !hasHash {
		return nil, nil
	}

	currentContent, err := m.realFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for conflict check: %w", err)
	}

	currentHash := computeHash(currentContent)

	if currentHash != expectedHash {
		m.log.Warning(fmt.Sprintf("Conflict detected for %s: expected %s, got %s",
			path, expectedHash[:8], currentHash[:8]))

		return &ConflictInfo{
			Path:           path,
			ExpectedHash:   expectedHash,
			CurrentHash:    currentHash,
			DetectedAt:     time.Now(),
			CurrentContent: currentContent,
		}, nil
	}

	return nil, nil
}

// CheckGroupConflicts checks all files in a group for conflicts
func (m *Manager) CheckGroupConflicts(groupID string) ([]*ConflictInfo, error) {
	m.mu.RLock()
	group, exists := m.groups[groupID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}

	var conflicts []*ConflictInfo

	for _, change := range group.Changes {
		if change.Operation == OpCreate {
			continue
		}

		conflict, err := m.CheckConflict(change.Path)
		if err != nil {
			m.log.Warning(fmt.Sprintf("Error checking conflict for %s: %v", change.Path, err))
			continue
		}

		if conflict != nil {
			conflict.ChangeID = change.ID
			conflict.OriginalContent = change.OriginalContent
			conflicts = append(conflicts, conflict)
		}
	}

	return conflicts, nil
}

// ApplyGroup applies all changes in a group to the sandbox
func (m *Manager) ApplyGroup(groupID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	group, exists := m.groups[groupID]
	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	if group.Applied {
		return fmt.Errorf("group already applied: %s", groupID)
	}

	for _, change := range group.Changes {
		if err := m.applyChangeToSandbox(change); err != nil {
			return fmt.Errorf("failed to apply change %s: %w", change.ID, err)
		}
		change.Applied = true
	}

	group.Applied = true

	m.addHistoryEntry(&HistoryEntry{
		ID:          generateID("history"),
		GroupID:     groupID,
		Type:        "apply",
		Changes:     group.Changes,
		Timestamp:   time.Now(),
		Description: group.Description,
		Undone:      false,
	})

	m.log.Info(fmt.Sprintf("Applied group %s with %d changes", groupID, len(group.Changes)))

	return nil
}

// RollbackGroup rolls back all changes in a group
func (m *Manager) RollbackGroup(groupID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	group, exists := m.groups[groupID]
	if !exists {
		return fmt.Errorf("group not found: %s", groupID)
	}

	if !group.Applied {
		return fmt.Errorf("group not applied, cannot rollback: %s", groupID)
	}

	for i := len(group.Changes) - 1; i >= 0; i-- {
		change := group.Changes[i]
		if err := m.rollbackChange(change); err != nil {
			return fmt.Errorf("failed to rollback change %s: %w", change.ID, err)
		}
		change.Applied = false
	}

	group.Applied = false

	m.addHistoryEntry(&HistoryEntry{
		ID:          generateID("history"),
		GroupID:     groupID,
		Type:        "rollback",
		Changes:     group.Changes,
		Timestamp:   time.Now(),
		Description: fmt.Sprintf("Rollback: %s", group.Description),
		Undone:      false,
	})

	m.log.Info(fmt.Sprintf("Rolled back group %s with %d changes", groupID, len(group.Changes)))

	return nil
}

// RollbackSession rolls back all changes in the current session
func (m *Manager) RollbackSession() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var appliedGroups []*ChangeGroup
	for _, group := range m.groups {
		if group.SessionID == m.sessionID && group.Applied {
			appliedGroups = append(appliedGroups, group)
		}
	}

	sort.Slice(appliedGroups, func(i, j int) bool {
		return appliedGroups[i].CreatedAt.After(appliedGroups[j].CreatedAt)
	})

	for _, group := range appliedGroups {
		for i := len(group.Changes) - 1; i >= 0; i-- {
			change := group.Changes[i]
			if err := m.rollbackChange(change); err != nil {
				m.log.Warning(fmt.Sprintf("Failed to rollback change %s: %v", change.ID, err))
				continue
			}
			change.Applied = false
		}
		group.Applied = false
	}

	m.log.Info(fmt.Sprintf("Rolled back session %s with %d groups", m.sessionID, len(appliedGroups)))

	return nil
}

// applyChangeToSandbox applies a single change to the sandbox
func (m *Manager) applyChangeToSandbox(change *Change) error {
	switch change.Operation {
	case OpCreate, OpModify:
		return m.sandbox.WriteFile(change.Path, change.Content, 0o644)
	case OpDelete:
		return m.sandbox.DeleteFile(change.Path)
	default:
		return fmt.Errorf("unknown operation: %d", change.Operation)
	}
}

// rollbackChange rolls back a single change
func (m *Manager) rollbackChange(change *Change) error {
	switch change.Operation {
	case OpCreate:
		m.sandbox.DiscardFile(change.Path)
		return nil
	case OpModify:
		if change.OriginalContent != nil {
			return m.sandbox.WriteFile(change.Path, change.OriginalContent, 0o644)
		}
		m.sandbox.DiscardFile(change.Path)
		return nil
	case OpDelete:
		if change.OriginalContent != nil {
			return m.sandbox.WriteFile(change.Path, change.OriginalContent, 0o644)
		}
		return nil
	default:
		return fmt.Errorf("unknown operation: %d", change.Operation)
	}
}

// addHistoryEntry adds an entry to history, managing size limits
func (m *Manager) addHistoryEntry(entry *HistoryEntry) {
	if m.historyIndex < len(m.history)-1 {
		m.history = m.history[:m.historyIndex+1]
	}

	m.history = append(m.history, entry)
	m.historyIndex = len(m.history) - 1

	if len(m.history) > m.config.MaxHistorySize {
		excess := len(m.history) - m.config.MaxHistorySize
		m.history = m.history[excess:]
		m.historyIndex -= excess
		if m.historyIndex < -1 {
			m.historyIndex = -1
		}
	}
}

// RollbackToChange rolls back all changes up to and including the specified change
func (m *Manager) RollbackToChange(changeID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Find the change and its group
	var targetGroup *ChangeGroup
	var targetChangeIndex int = -1

	for _, group := range m.groups {
		if group.SessionID != m.sessionID || !group.Applied {
			continue
		}
		for i, change := range group.Changes {
			if change.ID == changeID {
				targetGroup = group
				targetChangeIndex = i
				break
			}
		}
		if targetGroup != nil {
			break
		}
	}

	if targetGroup == nil {
		return fmt.Errorf("change not found: %s", changeID)
	}

	// Rollback all changes in the group from the end to the target change
	for i := len(targetGroup.Changes) - 1; i >= targetChangeIndex; i-- {
		change := targetGroup.Changes[i]
		if !change.Applied {
			continue
		}
		if err := m.rollbackChange(change); err != nil {
			return fmt.Errorf("failed to rollback change %s: %w", change.ID, err)
		}
		change.Applied = false
	}

	// Check if all changes in group are rolled back
	allRolledBack := true
	for _, change := range targetGroup.Changes {
		if change.Applied {
			allRolledBack = false
			break
		}
	}
	if allRolledBack {
		targetGroup.Applied = false
	}

	m.log.Info(fmt.Sprintf("Rolled back to change %s", changeID))

	return nil
}

// RollbackLastN rolls back the last N applied changes
func (m *Manager) RollbackLastN(n int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if n <= 0 {
		return fmt.Errorf("n must be positive, got %d", n)
	}

	// Collect all applied changes sorted by timestamp (newest first)
	var appliedChanges []*Change
	for _, group := range m.groups {
		if group.SessionID != m.sessionID {
			continue
		}
		for _, change := range group.Changes {
			if change.Applied {
				appliedChanges = append(appliedChanges, change)
			}
		}
	}

	// Sort by timestamp descending (newest first)
	sort.Slice(appliedChanges, func(i, j int) bool {
		return appliedChanges[i].Timestamp.After(appliedChanges[j].Timestamp)
	})

	// Rollback up to n changes
	rolledBack := 0
	for _, change := range appliedChanges {
		if rolledBack >= n {
			break
		}
		if err := m.rollbackChange(change); err != nil {
			m.log.Warning(fmt.Sprintf("Failed to rollback change %s: %v", change.ID, err))
			continue
		}
		change.Applied = false
		rolledBack++
	}

	// Update group applied status
	for _, group := range m.groups {
		if group.SessionID != m.sessionID {
			continue
		}
		allRolledBack := true
		for _, change := range group.Changes {
			if change.Applied {
				allRolledBack = false
				break
			}
		}
		if allRolledBack && group.Applied {
			group.Applied = false
		}
	}

	m.log.Info(fmt.Sprintf("Rolled back %d changes", rolledBack))

	return nil
}

// GetRollbackPreview returns a preview of changes that would be rolled back
func (m *Manager) GetRollbackPreview() []*Change {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var appliedChanges []*Change
	for _, group := range m.groups {
		if group.SessionID != m.sessionID {
			continue
		}
		for _, change := range group.Changes {
			if change.Applied {
				appliedChanges = append(appliedChanges, change)
			}
		}
	}

	// Sort by timestamp descending (newest first - order they would be rolled back)
	sort.Slice(appliedChanges, func(i, j int) bool {
		return appliedChanges[i].Timestamp.After(appliedChanges[j].Timestamp)
	})

	return appliedChanges
}

// ApplyWithConflictCheck applies all pending changes with conflict detection
func (m *Manager) ApplyWithConflictCheck() ([]Conflict, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check for conflicts first
	var conflicts []Conflict
	for _, group := range m.groups {
		if group.SessionID != m.sessionID || group.Applied {
			continue
		}
		for _, change := range group.Changes {
			if change.Operation == OpCreate {
				continue // New files can't have conflicts
			}
			conflict, err := m.conflictDetector.CheckConflict(change.Path)
			if err != nil {
				m.log.Warning(fmt.Sprintf("Error checking conflict for %s: %v", change.Path, err))
				continue
			}
			if conflict != nil {
				conflicts = append(conflicts, Conflict{
					FilePath:        conflict.FilePath,
					ExpectedHash:    conflict.ExpectedHash,
					ActualHash:      conflict.ActualHash,
					Message:         conflict.Message,
					OurContent:      change.Content,
					TheirContent:    conflict.TheirContent,
					OriginalContent: change.OriginalContent,
					DetectedAt:      conflict.DetectedAt,
				})
			}
		}
	}

	// If there are conflicts, return them without applying
	if len(conflicts) > 0 {
		return conflicts, fmt.Errorf("detected %d conflicts, resolve before applying", len(conflicts))
	}

	// Apply all pending changes
	for _, group := range m.groups {
		if group.SessionID != m.sessionID || group.Applied {
			continue
		}
		for _, change := range group.Changes {
			if err := m.applyChangeToSandbox(change); err != nil {
				return nil, fmt.Errorf("failed to apply change %s: %w", change.ID, err)
			}
			change.Applied = true
		}
		group.Applied = true

		m.addHistoryEntry(&HistoryEntry{
			ID:          generateID("history"),
			GroupID:     group.ID,
			Type:        "apply",
			Changes:     group.Changes,
			Timestamp:   time.Now(),
			Description: group.Description,
			Undone:      false,
		})
	}

	return nil, nil
}

// GetSessionStats returns statistics about the current session
func (m *Manager) GetSessionStats() SessionStats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := SessionStats{
		CanUndo: m.canUndoLocked(),
		CanRedo: m.canRedoLocked(),
	}

	changedFiles := make(map[string]bool)
	for _, group := range m.groups {
		if group.SessionID != m.sessionID {
			continue
		}
		stats.GroupCount++
		for _, change := range group.Changes {
			changedFiles[change.Path] = true
			if change.Applied {
				stats.AppliedChanges++
			} else {
				stats.PendingChanges++
			}

			// Calculate lines added/removed
			if change.Content != nil {
				stats.LinesAdded += countLines(change.Content)
			}
			if change.OriginalContent != nil {
				stats.LinesRemoved += countLines(change.OriginalContent)
			}
		}
	}

	stats.FilesChanged = len(changedFiles)

	// Check for conflicts
	conflicts := m.conflictDetector.CheckAllConflicts()
	stats.HasConflicts = len(conflicts) > 0

	return stats
}

// GetConflictDetector returns the conflict detector for direct access
func (m *Manager) GetConflictDetector() *ConflictDetector {
	return m.conflictDetector
}

// canUndoLocked checks if undo is available (must be called with lock held)
func (m *Manager) canUndoLocked() bool {
	if m.historyIndex < 0 || len(m.history) == 0 {
		return false
	}
	return !m.history[m.historyIndex].Undone
}

// canRedoLocked checks if redo is available (must be called with lock held)
func (m *Manager) canRedoLocked() bool {
	nextIndex := m.historyIndex + 1
	if nextIndex >= len(m.history) {
		return false
	}
	return m.history[nextIndex].Undone
}

// countLines counts the number of lines in content
func countLines(content []byte) int {
	if len(content) == 0 {
		return 0
	}
	count := 1
	for _, b := range content {
		if b == '\n' {
			count++
		}
	}
	return count
}
