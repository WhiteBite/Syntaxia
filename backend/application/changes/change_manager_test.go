package changes

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// mockFS implements domain.FileSystemProvider for testing
type mockFS struct {
	files map[string][]byte
	mu    sync.RWMutex
}

func newMockFS() *mockFS {
	return &mockFS{files: make(map[string][]byte)}
}

func (m *mockFS) ReadFile(filename string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	content, ok := m.files[filename]
	if !ok {
		return nil, os.ErrNotExist
	}
	return content, nil
}

func (m *mockFS) WriteFile(filename string, data []byte, perm int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[filename] = data
	return nil
}

func (m *mockFS) MkdirAll(path string, perm int) error {
	return nil
}

func (m *mockFS) setFile(path string, content []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.files[path] = content
}

// mockSandbox implements domain.SandboxFS for testing
type mockSandbox struct {
	changes     map[string]*domain.SandboxFileChange
	projectRoot string
	mu          sync.RWMutex
}

func newMockSandbox() *mockSandbox {
	return &mockSandbox{
		changes:     make(map[string]*domain.SandboxFileChange),
		projectRoot: "/project",
	}
}

func (s *mockSandbox) ReadFile(filename string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if change, ok := s.changes[filename]; ok {
		if change.Operation == domain.SandboxOpDelete {
			return nil, os.ErrNotExist
		}
		return change.Content, nil
	}
	return nil, os.ErrNotExist
}

func (s *mockSandbox) WriteFile(filename string, data []byte, perm int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.changes[filename] = &domain.SandboxFileChange{
		Path:      filename,
		Content:   data,
		Operation: domain.SandboxOpModify,
		Timestamp: time.Now(),
	}
	return nil
}

func (s *mockSandbox) MkdirAll(path string, perm int) error {
	return nil
}

func (s *mockSandbox) DeleteFile(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.changes[filename] = &domain.SandboxFileChange{
		Path:      filename,
		Operation: domain.SandboxOpDelete,
		Timestamp: time.Now(),
	}
	return nil
}

func (s *mockSandbox) GetChanges() []*domain.SandboxFileChange {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*domain.SandboxFileChange, 0, len(s.changes))
	for _, c := range s.changes {
		result = append(result, c)
	}
	return result
}

func (s *mockSandbox) GetDiff(filename string) (string, error)   { return "", nil }
func (s *mockSandbox) GetAllDiffs() (string, error)              { return "", nil }
func (s *mockSandbox) Apply() error                              { return nil }
func (s *mockSandbox) Discard()                                  { s.changes = make(map[string]*domain.SandboxFileChange) }
func (s *mockSandbox) DiscardFile(filename string)               { delete(s.changes, filename) }
func (s *mockSandbox) HasChanges() bool                          { return len(s.changes) > 0 }
func (s *mockSandbox) GetChangeCount() int                       { return len(s.changes) }
func (s *mockSandbox) GetProjectRoot() string                    { return s.projectRoot }
func (s *mockSandbox) SetProjectRoot(root string)                { s.projectRoot = root }

func TestChangeOperation_String(t *testing.T) {
	tests := []struct {
		name string
		op   ChangeOperation
		want string
	}{
		{"create", OpCreate, "create"},
		{"modify", OpModify, "modify"},
		{"delete", OpDelete, "delete"},
		{"unknown", ChangeOperation(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewManager(t *testing.T) {
	tests := []struct {
		name   string
		config *ManagerConfig
	}{
		{"with default config", nil},
		{"with custom config", &ManagerConfig{MaxHistorySize: 50}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sandbox := newMockSandbox()
			realFS := newMockFS()
			log := &domain.NoopLogger{}

			manager := NewManager(sandbox, realFS, log, tt.config)

			require.NotNil(t, manager)
			assert.NotEmpty(t, manager.GetSessionID())
		})
	}
}

func TestManager_StartAndEndGroup(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"simple group", "Test changes", false},
		{"empty description", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

			groupID := manager.StartGroup(tt.description)
			assert.NotEmpty(t, groupID)

			group, err := manager.EndGroup()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, groupID, group.ID)
				assert.Equal(t, tt.description, group.Description)
			}
		})
	}
}

func TestManager_EndGroup_NoActiveGroup(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	_, err := manager.EndGroup()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active change group")
}

func TestManager_RecordChange(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		content   []byte
		operation ChangeOperation
		existing  []byte
	}{
		{"create new file", "new.txt", []byte("new content"), OpCreate, nil},
		{"modify existing", "existing.txt", []byte("modified"), OpModify, []byte("original")},
		{"delete file", "delete.txt", nil, OpDelete, []byte("to delete")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			realFS := newMockFS()
			if tt.existing != nil {
				realFS.setFile(tt.path, tt.existing)
			}

			manager := NewManager(newMockSandbox(), realFS, &domain.NoopLogger{}, nil)
			manager.StartGroup("test")

			change, err := manager.RecordChange(tt.path, tt.content, tt.operation)

			require.NoError(t, err)
			assert.NotEmpty(t, change.ID)
			assert.Equal(t, tt.path, change.Path)
			assert.Equal(t, tt.operation, change.Operation)
			assert.Equal(t, tt.content, change.Content)

			if tt.existing != nil && tt.operation != OpCreate {
				assert.Equal(t, tt.existing, change.OriginalContent)
				assert.NotEmpty(t, change.OriginalHash)
			}
		})
	}
}

func TestManager_ApplyGroup(t *testing.T) {
	tests := []struct {
		name    string
		changes []struct {
			path    string
			content []byte
			op      ChangeOperation
		}
		wantErr bool
	}{
		{
			name: "apply single change",
			changes: []struct {
				path    string
				content []byte
				op      ChangeOperation
			}{
				{"file.txt", []byte("content"), OpCreate},
			},
			wantErr: false,
		},
		{
			name: "apply multiple changes",
			changes: []struct {
				path    string
				content []byte
				op      ChangeOperation
			}{
				{"file1.txt", []byte("content1"), OpCreate},
				{"file2.txt", []byte("content2"), OpCreate},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sandbox := newMockSandbox()
			manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

			groupID := manager.StartGroup("test")
			for _, c := range tt.changes {
				_, err := manager.RecordChange(c.path, c.content, c.op)
				require.NoError(t, err)
			}
			_, err := manager.EndGroup()
			require.NoError(t, err)

			err = manager.ApplyGroup(groupID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, len(tt.changes), sandbox.GetChangeCount())
			}
		})
	}
}

func TestManager_ApplyGroup_NotFound(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.ApplyGroup("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "group not found")
}

func TestManager_ApplyGroup_AlreadyApplied(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
	_, _ = manager.EndGroup()

	err := manager.ApplyGroup(groupID)
	require.NoError(t, err)

	err = manager.ApplyGroup(groupID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already applied")
}

func TestManager_RollbackGroup(t *testing.T) {
	tests := []struct {
		name    string
		changes []struct {
			path     string
			content  []byte
			original []byte
			op       ChangeOperation
		}
	}{
		{
			name: "rollback create",
			changes: []struct {
				path     string
				content  []byte
				original []byte
				op       ChangeOperation
			}{
				{"new.txt", []byte("new"), nil, OpCreate},
			},
		},
		{
			name: "rollback modify",
			changes: []struct {
				path     string
				content  []byte
				original []byte
				op       ChangeOperation
			}{
				{"existing.txt", []byte("modified"), []byte("original"), OpModify},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			realFS := newMockFS()
			sandbox := newMockSandbox()

			for _, c := range tt.changes {
				if c.original != nil {
					realFS.setFile(c.path, c.original)
				}
			}

			manager := NewManager(sandbox, realFS, &domain.NoopLogger{}, nil)

			groupID := manager.StartGroup("test")
			for _, c := range tt.changes {
				_, err := manager.RecordChange(c.path, c.content, c.op)
				require.NoError(t, err)
			}
			_, _ = manager.EndGroup()

			err := manager.ApplyGroup(groupID)
			require.NoError(t, err)

			err = manager.RollbackGroup(groupID)
			require.NoError(t, err)

			group, _ := manager.GetGroup(groupID)
			assert.False(t, group.Applied)
		})
	}
}

func TestManager_RollbackGroup_NotApplied(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
	_, _ = manager.EndGroup()

	err := manager.RollbackGroup(groupID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not applied")
}

func TestManager_RollbackSession(t *testing.T) {
	sandbox := newMockSandbox()
	manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

	// Create and apply multiple groups
	for i := 0; i < 3; i++ {
		groupID := manager.StartGroup("test")
		_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
		_, _ = manager.EndGroup()
		_ = manager.ApplyGroup(groupID)
	}

	err := manager.RollbackSession()
	require.NoError(t, err)

	// All groups should be unapplied
	groups := manager.GetSessionGroups()
	for _, g := range groups {
		assert.False(t, g.Applied)
	}
}

func TestManager_UndoRedo(t *testing.T) {
	sandbox := newMockSandbox()
	manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

	// Initially no undo/redo available
	assert.False(t, manager.CanUndo())
	assert.False(t, manager.CanRedo())

	// Create and apply a group
	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	// Now undo should be available
	assert.True(t, manager.CanUndo())
	assert.False(t, manager.CanRedo())

	// Undo
	err := manager.Undo()
	require.NoError(t, err)
	assert.False(t, manager.CanUndo())
	assert.True(t, manager.CanRedo())

	// Redo
	err = manager.Redo()
	require.NoError(t, err)
	assert.True(t, manager.CanUndo())
	assert.False(t, manager.CanRedo())
}

func TestManager_Undo_NothingToUndo(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.Undo()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to undo")
}

func TestManager_Redo_NothingToRedo(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.Redo()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nothing to redo")
}

func TestManager_ConflictDetection(t *testing.T) {
	tests := []struct {
		name           string
		originalContent []byte
		modifiedContent []byte
		expectConflict bool
	}{
		{
			name:           "no conflict - same content",
			originalContent: []byte("original"),
			modifiedContent: []byte("original"),
			expectConflict: false,
		},
		{
			name:           "conflict - content changed",
			originalContent: []byte("original"),
			modifiedContent: []byte("modified externally"),
			expectConflict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			realFS := newMockFS()
			realFS.setFile("file.txt", tt.originalContent)

			manager := NewManager(newMockSandbox(), realFS, &domain.NoopLogger{}, nil)

			// Record a change (this stores the hash)
			manager.StartGroup("test")
			_, _ = manager.RecordChange("file.txt", []byte("new content"), OpModify)
			_, _ = manager.EndGroup()

			// Simulate external modification
			realFS.setFile("file.txt", tt.modifiedContent)

			// Check for conflict
			conflict, err := manager.CheckConflict("file.txt")
			require.NoError(t, err)

			if tt.expectConflict {
				require.NotNil(t, conflict)
				assert.Equal(t, "file.txt", conflict.Path)
			} else {
				assert.Nil(t, conflict)
			}
		})
	}
}

func TestManager_CheckGroupConflicts(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file1.txt", []byte("original1"))
	realFS.setFile("file2.txt", []byte("original2"))

	manager := NewManager(newMockSandbox(), realFS, &domain.NoopLogger{}, nil)

	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file1.txt", []byte("modified1"), OpModify)
	_, _ = manager.RecordChange("file2.txt", []byte("modified2"), OpModify)
	_, _ = manager.EndGroup()

	// Modify one file externally
	realFS.setFile("file1.txt", []byte("external change"))

	conflicts, err := manager.CheckGroupConflicts(groupID)
	require.NoError(t, err)
	assert.Len(t, conflicts, 1)
	assert.Equal(t, "file1.txt", conflicts[0].Path)
}

func TestManager_CheckGroupConflicts_NotFound(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	_, err := manager.CheckGroupConflicts("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "group not found")
}

func TestManager_GetHistory(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	// Create and apply multiple groups
	for i := 0; i < 3; i++ {
		groupID := manager.StartGroup("test")
		_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
		_, _ = manager.EndGroup()
		_ = manager.ApplyGroup(groupID)
	}

	history := manager.GetHistory()
	assert.Len(t, history, 3)
}

func TestManager_ClearHistory(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	assert.Len(t, manager.GetHistory(), 1)

	manager.ClearHistory()

	assert.Len(t, manager.GetHistory(), 0)
	assert.False(t, manager.CanUndo())
}

func TestManager_NewSession(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	oldSessionID := manager.GetSessionID()

	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	newSessionID := manager.NewSession()

	assert.NotEqual(t, oldSessionID, newSessionID)
	assert.Len(t, manager.GetSessionGroups(), 0)
	assert.Len(t, manager.GetHistory(), 0)
}

func TestManager_GetPendingAndAppliedChanges(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	// Create group with changes but don't apply
	manager.StartGroup("pending")
	_, _ = manager.RecordChange("pending.txt", []byte("pending"), OpCreate)
	_, _ = manager.EndGroup()

	// Create and apply another group
	groupID := manager.StartGroup("applied")
	_, _ = manager.RecordChange("applied.txt", []byte("applied"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	pending := manager.GetPendingChanges()
	applied := manager.GetAppliedChanges()

	assert.Len(t, pending, 1)
	assert.Equal(t, "pending.txt", pending[0].Path)

	assert.Len(t, applied, 1)
	assert.Equal(t, "applied.txt", applied[0].Path)
}

func TestManager_HistorySizeLimit(t *testing.T) {
	config := &ManagerConfig{MaxHistorySize: 3}
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, config)

	// Create more entries than the limit
	for i := 0; i < 5; i++ {
		groupID := manager.StartGroup("test")
		_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
		_, _ = manager.EndGroup()
		_ = manager.ApplyGroup(groupID)
	}

	history := manager.GetHistory()
	assert.Len(t, history, 3)
}

func TestManager_ConcurrentAccess(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	var wg sync.WaitGroup
	numGoroutines := 50

	// Concurrent group operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			groupID := manager.StartGroup("concurrent")
			_, _ = manager.RecordChange("file.txt", []byte("content"), OpCreate)
			_, _ = manager.EndGroup()
			_ = manager.ApplyGroup(groupID)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = manager.GetSessionGroups()
			_ = manager.GetHistory()
			_ = manager.CanUndo()
			_ = manager.CanRedo()
		}()
	}

	wg.Wait()
}

func TestManager_UpdateFileHash(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("content"))

	manager := NewManager(newMockSandbox(), realFS, &domain.NoopLogger{}, nil)

	err := manager.UpdateFileHash("file.txt")
	require.NoError(t, err)

	// Now modify the file
	realFS.setFile("file.txt", []byte("modified"))

	conflict, err := manager.CheckConflict("file.txt")
	require.NoError(t, err)
	require.NotNil(t, conflict)
}

func TestManager_UpdateFileHash_NotFound(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.UpdateFileHash("nonexistent.txt")
	assert.Error(t, err)
}

func TestComputeHash(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
	}{
		{"empty", []byte{}},
		{"simple", []byte("hello")},
		{"unicode", []byte("привет мир")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := computeHash(tt.content)
			assert.Len(t, hash, 64) // SHA256 hex = 64 chars

			// Same content should produce same hash
			hash2 := computeHash(tt.content)
			assert.Equal(t, hash, hash2)
		})
	}
}

func TestGenerateID(t *testing.T) {
	id1 := generateID("test")
	id2 := generateID("test")

	assert.Contains(t, id1, "test-")
	assert.Contains(t, id2, "test-")
	assert.NotEqual(t, id1, id2)
}


// =============================================================================
// ConflictDetector Tests
// =============================================================================

func TestNewConflictDetector(t *testing.T) {
	realFS := newMockFS()
	log := &domain.NoopLogger{}

	detector := NewConflictDetector(realFS, log)

	require.NotNil(t, detector)
	assert.Equal(t, 0, detector.GetSnapshotCount())
}

func TestConflictDetector_SnapshotFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		content []byte
		wantErr bool
	}{
		{"snapshot existing file", "file.txt", []byte("content"), false},
		{"snapshot nonexistent file", "missing.txt", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			realFS := newMockFS()
			if tt.content != nil {
				realFS.setFile(tt.path, tt.content)
			}

			detector := NewConflictDetector(realFS, &domain.NoopLogger{})

			err := detector.SnapshotFile(tt.path)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.True(t, detector.HasSnapshot(tt.path))
			}
		})
	}
}

func TestConflictDetector_SnapshotFiles(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file1.txt", []byte("content1"))
	realFS.setFile("file2.txt", []byte("content2"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})

	err := detector.SnapshotFiles([]string{"file1.txt", "file2.txt"})

	require.NoError(t, err)
	assert.Equal(t, 2, detector.GetSnapshotCount())
	assert.True(t, detector.HasSnapshot("file1.txt"))
	assert.True(t, detector.HasSnapshot("file2.txt"))
}

func TestConflictDetector_CheckConflict_NoConflict(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("original"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})
	_ = detector.SnapshotFile("file.txt")

	// File unchanged
	conflict, err := detector.CheckConflict("file.txt")

	require.NoError(t, err)
	assert.Nil(t, conflict)
}

func TestConflictDetector_CheckConflict_DetectsConflict(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("original"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})
	_ = detector.SnapshotFile("file.txt")

	// Modify file externally
	realFS.setFile("file.txt", []byte("modified externally"))

	conflict, err := detector.CheckConflict("file.txt")

	require.NoError(t, err)
	require.NotNil(t, conflict)
	assert.Equal(t, "file.txt", conflict.FilePath)
	assert.NotEqual(t, conflict.ExpectedHash, conflict.ActualHash)
	assert.Contains(t, conflict.Message, "modified externally")
}

func TestConflictDetector_CheckConflict_NoSnapshot(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("content"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})

	// No snapshot taken
	conflict, err := detector.CheckConflict("file.txt")

	require.NoError(t, err)
	assert.Nil(t, conflict) // No conflict if no snapshot
}

func TestConflictDetector_CheckAllConflicts(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file1.txt", []byte("original1"))
	realFS.setFile("file2.txt", []byte("original2"))
	realFS.setFile("file3.txt", []byte("original3"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})
	_ = detector.SnapshotFiles([]string{"file1.txt", "file2.txt", "file3.txt"})

	// Modify some files externally
	realFS.setFile("file1.txt", []byte("modified1"))
	realFS.setFile("file3.txt", []byte("modified3"))

	conflicts := detector.CheckAllConflicts()

	assert.Len(t, conflicts, 2)

	paths := make(map[string]bool)
	for _, c := range conflicts {
		paths[c.FilePath] = true
	}
	assert.True(t, paths["file1.txt"])
	assert.True(t, paths["file3.txt"])
	assert.False(t, paths["file2.txt"])
}

func TestConflictDetector_ResolveConflict(t *testing.T) {
	tests := []struct {
		name       string
		resolution ConflictResolution
		wantErr    bool
	}{
		{"keep ours", ResolutionKeepOurs, false},
		{"keep theirs", ResolutionKeepTheirs, false},
		{"merge not implemented", ResolutionMerge, true},
		{"unknown resolution", ConflictResolution("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			realFS := newMockFS()
			realFS.setFile("file.txt", []byte("content"))

			detector := NewConflictDetector(realFS, &domain.NoopLogger{})
			_ = detector.SnapshotFile("file.txt")

			err := detector.ResolveConflict("file.txt", tt.resolution)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConflictDetector_ClearSnapshot(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("content"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})
	_ = detector.SnapshotFile("file.txt")

	assert.True(t, detector.HasSnapshot("file.txt"))

	detector.ClearSnapshot("file.txt")

	assert.False(t, detector.HasSnapshot("file.txt"))
}

func TestConflictDetector_ClearAllSnapshots(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file1.txt", []byte("content1"))
	realFS.setFile("file2.txt", []byte("content2"))

	detector := NewConflictDetector(realFS, &domain.NoopLogger{})
	_ = detector.SnapshotFiles([]string{"file1.txt", "file2.txt"})

	assert.Equal(t, 2, detector.GetSnapshotCount())

	detector.ClearAllSnapshots()

	assert.Equal(t, 0, detector.GetSnapshotCount())
}

// =============================================================================
// RollbackToChange Tests
// =============================================================================

func TestManager_RollbackToChange(t *testing.T) {
	sandbox := newMockSandbox()
	manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

	// Create and apply a group with multiple changes
	groupID := manager.StartGroup("test")
	change1, _ := manager.RecordChange("file1.txt", []byte("content1"), OpCreate)
	change2, _ := manager.RecordChange("file2.txt", []byte("content2"), OpCreate)
	change3, _ := manager.RecordChange("file3.txt", []byte("content3"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	// Rollback to change2 (should rollback change3 and change2)
	err := manager.RollbackToChange(change2.ID)

	require.NoError(t, err)
	assert.True(t, change1.Applied)
	assert.False(t, change2.Applied)
	assert.False(t, change3.Applied)
}

func TestManager_RollbackToChange_NotFound(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.RollbackToChange("nonexistent")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "change not found")
}

// =============================================================================
// RollbackLastN Tests
// =============================================================================

func TestManager_RollbackLastN(t *testing.T) {
	tests := []struct {
		name           string
		totalChanges   int
		rollbackN      int
		expectedApplied int
	}{
		{"rollback 1 of 3", 3, 1, 2},
		{"rollback 2 of 3", 3, 2, 1},
		{"rollback all", 3, 3, 0},
		{"rollback more than available", 3, 5, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sandbox := newMockSandbox()
			manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

			// Create and apply changes
			groupID := manager.StartGroup("test")
			for i := 0; i < tt.totalChanges; i++ {
				_, _ = manager.RecordChange(fmt.Sprintf("file%d.txt", i), []byte("content"), OpCreate)
			}
			_, _ = manager.EndGroup()
			_ = manager.ApplyGroup(groupID)

			err := manager.RollbackLastN(tt.rollbackN)

			require.NoError(t, err)
			applied := manager.GetAppliedChanges()
			assert.Len(t, applied, tt.expectedApplied)
		})
	}
}

func TestManager_RollbackLastN_InvalidN(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	err := manager.RollbackLastN(0)
	assert.Error(t, err)

	err = manager.RollbackLastN(-1)
	assert.Error(t, err)
}

// =============================================================================
// GetRollbackPreview Tests
// =============================================================================

func TestManager_GetRollbackPreview(t *testing.T) {
	sandbox := newMockSandbox()
	manager := NewManager(sandbox, newMockFS(), &domain.NoopLogger{}, nil)

	// Create and apply changes
	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("file1.txt", []byte("content1"), OpCreate)
	_, _ = manager.RecordChange("file2.txt", []byte("content2"), OpCreate)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	preview := manager.GetRollbackPreview()

	assert.Len(t, preview, 2)
	// Should be sorted newest first
	assert.True(t, preview[0].Timestamp.After(preview[1].Timestamp) || 
		preview[0].Timestamp.Equal(preview[1].Timestamp))
}

func TestManager_GetRollbackPreview_Empty(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	preview := manager.GetRollbackPreview()

	assert.Len(t, preview, 0)
}

// =============================================================================
// ApplyWithConflictCheck Tests
// =============================================================================

func TestManager_ApplyWithConflictCheck_NoConflicts(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("existing.txt", []byte("original"))

	sandbox := newMockSandbox()
	manager := NewManager(sandbox, realFS, &domain.NoopLogger{}, nil)

	// Snapshot the file
	_ = manager.GetConflictDetector().SnapshotFile("existing.txt")

	// Create changes
	manager.StartGroup("test")
	_, _ = manager.RecordChange("new.txt", []byte("new content"), OpCreate)
	_, _ = manager.EndGroup()

	conflicts, err := manager.ApplyWithConflictCheck()

	require.NoError(t, err)
	assert.Len(t, conflicts, 0)
}

func TestManager_ApplyWithConflictCheck_WithConflicts(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("file.txt", []byte("original"))

	sandbox := newMockSandbox()
	manager := NewManager(sandbox, realFS, &domain.NoopLogger{}, nil)

	// Snapshot the file
	_ = manager.GetConflictDetector().SnapshotFile("file.txt")

	// Create a modify change
	manager.StartGroup("test")
	_, _ = manager.RecordChange("file.txt", []byte("our changes"), OpModify)
	_, _ = manager.EndGroup()

	// Simulate external modification
	realFS.setFile("file.txt", []byte("external changes"))

	conflicts, err := manager.ApplyWithConflictCheck()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "conflicts")
	assert.Len(t, conflicts, 1)
	assert.Equal(t, "file.txt", conflicts[0].FilePath)
}

// =============================================================================
// GetSessionStats Tests
// =============================================================================

func TestManager_GetSessionStats(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("existing.txt", []byte("original\ncontent\nhere"))

	sandbox := newMockSandbox()
	manager := NewManager(sandbox, realFS, &domain.NoopLogger{}, nil)

	// Create and apply some changes
	groupID := manager.StartGroup("test")
	_, _ = manager.RecordChange("new.txt", []byte("line1\nline2"), OpCreate)
	_, _ = manager.RecordChange("existing.txt", []byte("modified"), OpModify)
	_, _ = manager.EndGroup()
	_ = manager.ApplyGroup(groupID)

	// Create pending changes
	manager.StartGroup("pending")
	_, _ = manager.RecordChange("pending.txt", []byte("pending"), OpCreate)
	_, _ = manager.EndGroup()

	stats := manager.GetSessionStats()

	assert.Equal(t, 3, stats.FilesChanged)
	assert.Equal(t, 2, stats.AppliedChanges)
	assert.Equal(t, 1, stats.PendingChanges)
	assert.Equal(t, 2, stats.GroupCount)
	assert.True(t, stats.CanUndo)
	assert.False(t, stats.CanRedo)
	assert.Greater(t, stats.LinesAdded, 0)
}

func TestManager_GetSessionStats_Empty(t *testing.T) {
	manager := NewManager(newMockSandbox(), newMockFS(), &domain.NoopLogger{}, nil)

	stats := manager.GetSessionStats()

	assert.Equal(t, 0, stats.FilesChanged)
	assert.Equal(t, 0, stats.AppliedChanges)
	assert.Equal(t, 0, stats.PendingChanges)
	assert.Equal(t, 0, stats.GroupCount)
	assert.False(t, stats.CanUndo)
	assert.False(t, stats.CanRedo)
	assert.False(t, stats.HasConflicts)
}

// =============================================================================
// CountLines Tests
// =============================================================================

func TestCountLines(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    int
	}{
		{"empty", []byte{}, 0},
		{"single line", []byte("hello"), 1},
		{"two lines", []byte("hello\nworld"), 2},
		{"three lines", []byte("a\nb\nc"), 3},
		{"trailing newline", []byte("hello\n"), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countLines(tt.content)
			assert.Equal(t, tt.want, got)
		})
	}
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestManager_FullWorkflow_WithConflictDetection(t *testing.T) {
	realFS := newMockFS()
	realFS.setFile("config.json", []byte(`{"version": 1}`))

	sandbox := newMockSandbox()
	manager := NewManager(sandbox, realFS, &domain.NoopLogger{}, nil)

	// 1. Snapshot files before editing
	err := manager.GetConflictDetector().SnapshotFile("config.json")
	require.NoError(t, err)

	// 2. Make changes
	groupID := manager.StartGroup("Update config")
	_, err = manager.RecordChange("config.json", []byte(`{"version": 2}`), OpModify)
	require.NoError(t, err)
	_, err = manager.EndGroup()
	require.NoError(t, err)

	// 3. Apply with conflict check (should succeed)
	conflicts, err := manager.ApplyWithConflictCheck()
	require.NoError(t, err)
	assert.Len(t, conflicts, 0)

	// 4. Check stats
	stats := manager.GetSessionStats()
	assert.Equal(t, 1, stats.FilesChanged)
	assert.Equal(t, 1, stats.AppliedChanges)
	assert.True(t, stats.CanUndo)

	// 5. Undo
	err = manager.Undo()
	require.NoError(t, err)

	// 6. Redo
	err = manager.Redo()
	require.NoError(t, err)

	// 7. Rollback
	err = manager.RollbackGroup(groupID)
	require.NoError(t, err)

	group, _ := manager.GetGroup(groupID)
	assert.False(t, group.Applied)
}

func TestManager_ConcurrentConflictDetection(t *testing.T) {
	realFS := newMockFS()
	for i := 0; i < 10; i++ {
		realFS.setFile(fmt.Sprintf("file%d.txt", i), []byte(fmt.Sprintf("content%d", i)))
	}

	manager := NewManager(newMockSandbox(), realFS, &domain.NoopLogger{}, nil)
	detector := manager.GetConflictDetector()

	var wg sync.WaitGroup

	// Concurrent snapshots
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_ = detector.SnapshotFile(fmt.Sprintf("file%d.txt", idx))
		}(i)
	}

	// Concurrent conflict checks
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, _ = detector.CheckConflict(fmt.Sprintf("file%d.txt", idx))
		}(i)
	}

	// Concurrent stats
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = manager.GetSessionStats()
		}()
	}

	wg.Wait()
}
