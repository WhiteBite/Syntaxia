package sandbox

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

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
	return &mockFS{
		files: make(map[string][]byte),
	}
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

func (m *mockFS) getFile(path string) ([]byte, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	content, ok := m.files[path]
	return content, ok
}

func TestSandboxFS_WriteAndRead(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Write to sandbox
	content := []byte("hello world")
	err := sandbox.WriteFile("test.txt", content, 0o644)
	require.NoError(t, err)

	// Read from sandbox
	read, err := sandbox.ReadFile("test.txt")
	require.NoError(t, err)
	assert.Equal(t, content, read)

	// Verify real FS was not modified
	_, ok := mockFS.getFile("test.txt")
	assert.False(t, ok, "real FS should not be modified")
}

func TestSandboxFS_ReadFromRealFS(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/existing.txt", []byte("original content"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Read existing file from real FS
	content, err := sandbox.ReadFile("/project/existing.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("original content"), content)
}

func TestSandboxFS_ModifyExistingFile(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/file.txt", []byte("original"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Modify file
	err := sandbox.WriteFile("/project/file.txt", []byte("modified"), 0o644)
	require.NoError(t, err)

	// Read modified content from sandbox
	content, err := sandbox.ReadFile("/project/file.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("modified"), content)

	// Verify original is preserved in real FS
	original, ok := mockFS.getFile("/project/file.txt")
	assert.True(t, ok)
	assert.Equal(t, []byte("original"), original)

	// Check change operation
	changes := sandbox.GetChanges()
	require.Len(t, changes, 1)
	assert.Equal(t, domain.SandboxOpModify, changes[0].Operation)
	assert.Equal(t, []byte("original"), changes[0].OriginalContent)
}

func TestSandboxFS_DeleteFile(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/to-delete.txt", []byte("delete me"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Delete file
	err := sandbox.DeleteFile("/project/to-delete.txt")
	require.NoError(t, err)

	// Reading deleted file should return ErrNotExist
	_, err = sandbox.ReadFile("/project/to-delete.txt")
	assert.ErrorIs(t, err, os.ErrNotExist)

	// Verify real FS still has the file
	content, ok := mockFS.getFile("/project/to-delete.txt")
	assert.True(t, ok)
	assert.Equal(t, []byte("delete me"), content)

	// Check change operation
	changes := sandbox.GetChanges()
	require.Len(t, changes, 1)
	assert.Equal(t, domain.SandboxOpDelete, changes[0].Operation)
}

func TestSandboxFS_DeleteNonExistentFile(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Delete non-existent file should fail
	err := sandbox.DeleteFile("nonexistent.txt")
	assert.Error(t, err)
}

func TestSandboxFS_GetDiff(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/file.txt", []byte("line1\nline2\nline3"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Modify file
	err := sandbox.WriteFile("/project/file.txt", []byte("line1\nmodified\nline3"), 0o644)
	require.NoError(t, err)

	// Get diff
	diff, err := sandbox.GetDiff("/project/file.txt")
	require.NoError(t, err)

	// Verify diff contains expected markers
	assert.Contains(t, diff, "diff --git")
	assert.Contains(t, diff, "---")
	assert.Contains(t, diff, "+++")
	assert.Contains(t, diff, "-line1")
	assert.Contains(t, diff, "+line1")
}

func TestSandboxFS_GetDiff_NewFile(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Create new file
	err := sandbox.WriteFile("new-file.txt", []byte("new content"), 0o644)
	require.NoError(t, err)

	// Get diff
	diff, err := sandbox.GetDiff("new-file.txt")
	require.NoError(t, err)

	// Verify diff shows new file
	assert.Contains(t, diff, "new file mode")
	assert.Contains(t, diff, "/dev/null")
	assert.Contains(t, diff, "+new content")
}

func TestSandboxFS_GetDiff_DeletedFile(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/to-delete.txt", []byte("old content"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Delete file
	err := sandbox.DeleteFile("/project/to-delete.txt")
	require.NoError(t, err)

	// Get diff
	diff, err := sandbox.GetDiff("/project/to-delete.txt")
	require.NoError(t, err)

	// Verify diff shows deleted file
	assert.Contains(t, diff, "deleted file mode")
	assert.Contains(t, diff, "-old content")
}

func TestSandboxFS_GetAllDiffs(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/file1.txt", []byte("original1"))
	mockFS.setFile("/project/file2.txt", []byte("original2"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Modify multiple files
	err := sandbox.WriteFile("/project/file1.txt", []byte("modified1"), 0o644)
	require.NoError(t, err)

	err = sandbox.WriteFile("/project/file2.txt", []byte("modified2"), 0o644)
	require.NoError(t, err)

	// Get all diffs
	allDiffs, err := sandbox.GetAllDiffs()
	require.NoError(t, err)

	// Verify both files are in diff
	assert.Contains(t, allDiffs, "file1.txt")
	assert.Contains(t, allDiffs, "file2.txt")
}

func TestSandboxFS_Apply(t *testing.T) {
	// Create temp directory for real file operations
	tempDir, err := os.MkdirTemp("", "sandbox-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a real file
	existingFile := filepath.Join(tempDir, "existing.txt")
	err = os.WriteFile(existingFile, []byte("original"), 0o644)
	require.NoError(t, err)

	// Use real FS for apply test
	realFS := &realFileSystem{}
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS(tempDir, realFS, log)

	// Create new file in sandbox
	err = sandbox.WriteFile("new.txt", []byte("new content"), 0o644)
	require.NoError(t, err)

	// Modify existing file in sandbox
	err = sandbox.WriteFile(existingFile, []byte("modified"), 0o644)
	require.NoError(t, err)

	// Apply changes
	err = sandbox.Apply()
	require.NoError(t, err)

	// Verify new file was created
	newContent, err := os.ReadFile(filepath.Join(tempDir, "new.txt"))
	require.NoError(t, err)
	assert.Equal(t, []byte("new content"), newContent)

	// Verify existing file was modified
	modifiedContent, err := os.ReadFile(existingFile)
	require.NoError(t, err)
	assert.Equal(t, []byte("modified"), modifiedContent)

	// Verify sandbox is cleared
	assert.False(t, sandbox.HasChanges())
}

func TestSandboxFS_Discard(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Add some changes
	err := sandbox.WriteFile("file1.txt", []byte("content1"), 0o644)
	require.NoError(t, err)
	err = sandbox.WriteFile("file2.txt", []byte("content2"), 0o644)
	require.NoError(t, err)

	assert.Equal(t, 2, sandbox.GetChangeCount())

	// Discard all
	sandbox.Discard()

	assert.False(t, sandbox.HasChanges())
	assert.Equal(t, 0, sandbox.GetChangeCount())
}

func TestSandboxFS_DiscardFile(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Add some changes
	err := sandbox.WriteFile("file1.txt", []byte("content1"), 0o644)
	require.NoError(t, err)
	err = sandbox.WriteFile("file2.txt", []byte("content2"), 0o644)
	require.NoError(t, err)

	assert.Equal(t, 2, sandbox.GetChangeCount())

	// Discard single file
	sandbox.DiscardFile("file1.txt")

	assert.Equal(t, 1, sandbox.GetChangeCount())

	// Verify file1 is gone but file2 remains
	_, err = sandbox.ReadFile("file1.txt")
	assert.Error(t, err)

	content, err := sandbox.ReadFile("file2.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("content2"), content)
}

func TestSandboxFS_HasChanges(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	assert.False(t, sandbox.HasChanges())

	err := sandbox.WriteFile("test.txt", []byte("content"), 0o644)
	require.NoError(t, err)

	assert.True(t, sandbox.HasChanges())
}

func TestSandboxFS_GetChangeCount(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	assert.Equal(t, 0, sandbox.GetChangeCount())

	_ = sandbox.WriteFile("file1.txt", []byte("1"), 0o644)
	assert.Equal(t, 1, sandbox.GetChangeCount())

	_ = sandbox.WriteFile("file2.txt", []byte("2"), 0o644)
	assert.Equal(t, 2, sandbox.GetChangeCount())

	// Overwriting same file shouldn't increase count
	_ = sandbox.WriteFile("file1.txt", []byte("1-updated"), 0o644)
	assert.Equal(t, 2, sandbox.GetChangeCount())
}

func TestSandboxFS_ConcurrentAccess(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			filename := "file" + string(rune('0'+idx%10)) + ".txt"
			_ = sandbox.WriteFile(filename, []byte("content"), 0o644)
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			filename := "file" + string(rune('0'+idx%10)) + ".txt"
			_, _ = sandbox.ReadFile(filename)
		}(i)
	}

	// Concurrent status checks
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = sandbox.HasChanges()
			_ = sandbox.GetChangeCount()
			_ = sandbox.GetChanges()
		}()
	}

	wg.Wait()

	// Should complete without race conditions
	assert.True(t, sandbox.HasChanges())
}

func TestSandboxFS_PathNormalization(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Write with different path formats
	err := sandbox.WriteFile("src/file.txt", []byte("content"), 0o644)
	require.NoError(t, err)

	// Read with same path
	content, err := sandbox.ReadFile("src/file.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("content"), content)

	// Should be same file
	assert.Equal(t, 1, sandbox.GetChangeCount())
}

func TestSandboxFS_GetChanges_Sorted(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Add files in random order
	_ = sandbox.WriteFile("z.txt", []byte("z"), 0o644)
	_ = sandbox.WriteFile("a.txt", []byte("a"), 0o644)
	_ = sandbox.WriteFile("m.txt", []byte("m"), 0o644)

	changes := sandbox.GetChanges()
	require.Len(t, changes, 3)

	// Should be sorted by path
	assert.Equal(t, "a.txt", changes[0].Path)
	assert.Equal(t, "m.txt", changes[1].Path)
	assert.Equal(t, "z.txt", changes[2].Path)
}

func TestSandboxChangeOperation_String(t *testing.T) {
	assert.Equal(t, "create", opString(domain.SandboxOpCreate))
	assert.Equal(t, "modify", opString(domain.SandboxOpModify))
	assert.Equal(t, "delete", opString(domain.SandboxOpDelete))
	assert.Equal(t, "unknown", opString(domain.SandboxChangeOperation(99)))
}

func TestSandboxFS_EmptyDiff(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Get diff for non-existent change
	diff, err := sandbox.GetDiff("nonexistent.txt")
	require.NoError(t, err)
	assert.Empty(t, diff)

	// Get all diffs when empty
	allDiffs, err := sandbox.GetAllDiffs()
	require.NoError(t, err)
	assert.Empty(t, allDiffs)
}

func TestSandboxFS_ApplyEmpty(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Apply with no changes should succeed
	err := sandbox.Apply()
	require.NoError(t, err)
}

func TestSandboxFS_MultipleModifications(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/file.txt", []byte("original"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Multiple modifications to same file
	_ = sandbox.WriteFile("/project/file.txt", []byte("v1"), 0o644)
	_ = sandbox.WriteFile("/project/file.txt", []byte("v2"), 0o644)
	_ = sandbox.WriteFile("/project/file.txt", []byte("v3"), 0o644)

	// Should only have one change entry
	assert.Equal(t, 1, sandbox.GetChangeCount())

	// Content should be latest
	content, err := sandbox.ReadFile("/project/file.txt")
	require.NoError(t, err)
	assert.Equal(t, []byte("v3"), content)

	// Original should still be preserved
	changes := sandbox.GetChanges()
	require.Len(t, changes, 1)
	assert.Equal(t, []byte("original"), changes[0].OriginalContent)
}

func TestSandboxFS_SetProjectRoot(t *testing.T) {
	mockFS := newMockFS()
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/old-project", mockFS, log)

	assert.Equal(t, "/old-project", sandbox.GetProjectRoot())

	sandbox.SetProjectRoot("/new-project")
	assert.Equal(t, "/new-project", sandbox.GetProjectRoot())
}

// realFileSystem implements domain.FileSystemProvider using actual filesystem
type realFileSystem struct{}

func (r *realFileSystem) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (r *realFileSystem) WriteFile(filename string, data []byte, perm int) error {
	return os.WriteFile(filename, data, os.FileMode(perm))
}

func (r *realFileSystem) MkdirAll(path string, perm int) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

func TestSandboxFS_ApplyWithDelete(t *testing.T) {
	// Create temp directory for real file operations
	tempDir, err := os.MkdirTemp("", "sandbox-delete-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a file to delete
	fileToDelete := filepath.Join(tempDir, "to-delete.txt")
	err = os.WriteFile(fileToDelete, []byte("delete me"), 0o644)
	require.NoError(t, err)

	realFS := &realFileSystem{}
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS(tempDir, realFS, log)

	// Delete file in sandbox
	err = sandbox.DeleteFile(fileToDelete)
	require.NoError(t, err)

	// File should still exist on disk
	_, err = os.Stat(fileToDelete)
	require.NoError(t, err)

	// Apply changes
	err = sandbox.Apply()
	require.NoError(t, err)

	// File should be deleted
	_, err = os.Stat(fileToDelete)
	assert.True(t, os.IsNotExist(err))
}

func TestSandboxFS_ApplyCreatesDirectories(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "sandbox-mkdir-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	realFS := &realFileSystem{}
	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS(tempDir, realFS, log)

	// Create file in nested directory
	err = sandbox.WriteFile("deep/nested/dir/file.txt", []byte("content"), 0o644)
	require.NoError(t, err)

	// Apply changes
	err = sandbox.Apply()
	require.NoError(t, err)

	// Verify file was created with directories
	content, err := os.ReadFile(filepath.Join(tempDir, "deep/nested/dir/file.txt"))
	require.NoError(t, err)
	assert.Equal(t, []byte("content"), content)
}

func TestSandboxFS_DiffFormat(t *testing.T) {
	mockFS := newMockFS()
	mockFS.setFile("/project/test.go", []byte("package main\n\nfunc main() {\n\tprintln(\"hello\")\n}\n"))

	log := &domain.NoopLogger{}
	sandbox := NewSandboxFS("/project", mockFS, log)

	// Modify file
	newContent := "package main\n\nfunc main() {\n\tprintln(\"hello world\")\n}\n"
	err := sandbox.WriteFile("/project/test.go", []byte(newContent), 0o644)
	require.NoError(t, err)

	diff, err := sandbox.GetDiff("/project/test.go")
	require.NoError(t, err)

	// Verify diff format
	assert.True(t, strings.HasPrefix(diff, "diff --git"))
	assert.Contains(t, diff, "--- a/")
	assert.Contains(t, diff, "+++ b/")
	assert.Contains(t, diff, "@@")
}
