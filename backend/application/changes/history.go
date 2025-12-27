package changes

import "fmt"

// Undo undoes the last operation
func (m *Manager) Undo() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.historyIndex < 0 || len(m.history) == 0 {
		return fmt.Errorf("nothing to undo")
	}

	entry := m.history[m.historyIndex]
	if entry.Undone {
		return fmt.Errorf("already undone")
	}

	if entry.Type == "apply" {
		for i := len(entry.Changes) - 1; i >= 0; i-- {
			if err := m.rollbackChange(entry.Changes[i]); err != nil {
				return fmt.Errorf("failed to undo change: %w", err)
			}
			entry.Changes[i].Applied = false
		}
		if group, exists := m.groups[entry.GroupID]; exists {
			group.Applied = false
		}
	} else {
		for _, change := range entry.Changes {
			if err := m.applyChangeToSandbox(change); err != nil {
				return fmt.Errorf("failed to undo rollback: %w", err)
			}
			change.Applied = true
		}
		if group, exists := m.groups[entry.GroupID]; exists {
			group.Applied = true
		}
	}

	entry.Undone = true
	m.historyIndex--

	m.log.Info(fmt.Sprintf("Undone: %s", entry.Description))

	return nil
}

// Redo redoes the last undone operation
func (m *Manager) Redo() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	nextIndex := m.historyIndex + 1
	if nextIndex >= len(m.history) {
		return fmt.Errorf("nothing to redo")
	}

	entry := m.history[nextIndex]
	if !entry.Undone {
		return fmt.Errorf("not undone, cannot redo")
	}

	if entry.Type == "apply" {
		for _, change := range entry.Changes {
			if err := m.applyChangeToSandbox(change); err != nil {
				return fmt.Errorf("failed to redo apply: %w", err)
			}
			change.Applied = true
		}
		if group, exists := m.groups[entry.GroupID]; exists {
			group.Applied = true
		}
	} else {
		for i := len(entry.Changes) - 1; i >= 0; i-- {
			if err := m.rollbackChange(entry.Changes[i]); err != nil {
				return fmt.Errorf("failed to redo rollback: %w", err)
			}
			entry.Changes[i].Applied = false
		}
		if group, exists := m.groups[entry.GroupID]; exists {
			group.Applied = false
		}
	}

	entry.Undone = false
	m.historyIndex = nextIndex

	m.log.Info(fmt.Sprintf("Redone: %s", entry.Description))

	return nil
}

// CanUndo returns true if undo is available
func (m *Manager) CanUndo() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.historyIndex < 0 || len(m.history) == 0 {
		return false
	}
	return !m.history[m.historyIndex].Undone
}

// CanRedo returns true if redo is available
func (m *Manager) CanRedo() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	nextIndex := m.historyIndex + 1
	if nextIndex >= len(m.history) {
		return false
	}
	return m.history[nextIndex].Undone
}

// GetHistory returns the change history
func (m *Manager) GetHistory() []*HistoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*HistoryEntry, len(m.history))
	copy(result, m.history)
	return result
}

// ClearHistory clears the undo/redo history
func (m *Manager) ClearHistory() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.history = make([]*HistoryEntry, 0)
	m.historyIndex = -1
	m.log.Info("Cleared change history")
}
