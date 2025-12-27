package changes

import (
	"fmt"
	"sort"
)

// GetGroup returns a change group by ID
func (m *Manager) GetGroup(groupID string) (*ChangeGroup, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	group, exists := m.groups[groupID]
	if !exists {
		return nil, fmt.Errorf("group not found: %s", groupID)
	}
	return group, nil
}

// GetSessionGroups returns all groups for the current session
func (m *Manager) GetSessionGroups() []*ChangeGroup {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*ChangeGroup
	for _, group := range m.groups {
		if group.SessionID == m.sessionID {
			result = append(result, group)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	return result
}

// GetPendingChanges returns all pending (unapplied) changes
func (m *Manager) GetPendingChanges() []*Change {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Change
	for _, group := range m.groups {
		for _, change := range group.Changes {
			if !change.Applied {
				result = append(result, change)
			}
		}
	}
	return result
}

// GetAppliedChanges returns all applied changes
func (m *Manager) GetAppliedChanges() []*Change {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []*Change
	for _, group := range m.groups {
		for _, change := range group.Changes {
			if change.Applied {
				result = append(result, change)
			}
		}
	}
	return result
}

// NewSession starts a new session, clearing current state
func (m *Manager) NewSession() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.sessionID = generateID("session")
	m.currentGroup = nil
	m.groups = make(map[string]*ChangeGroup)
	m.history = make([]*HistoryEntry, 0)
	m.historyIndex = -1
	m.fileHashes = make(map[string]string)

	m.log.Info(fmt.Sprintf("Started new session: %s", m.sessionID))

	return m.sessionID
}

// GetSessionID returns the current session ID
func (m *Manager) GetSessionID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.sessionID
}

// UpdateFileHash updates the stored hash for a file (call after reading)
func (m *Manager) UpdateFileHash(path string) error {
	content, err := m.realFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	m.mu.Lock()
	m.fileHashes[path] = computeHash(content)
	m.mu.Unlock()

	return nil
}
