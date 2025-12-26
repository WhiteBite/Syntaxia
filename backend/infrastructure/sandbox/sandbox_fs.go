package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syntaxia/domain"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// SandboxFS implements domain.SandboxFS interface
// All changes are stored in memory until explicitly applied
type SandboxFS struct {
	mu          sync.RWMutex
	changes     map[string]*domain.SandboxFileChange
	projectRoot string
	realFS      domain.FileSystemProvider
	log         domain.Logger
}

// Compile-time check that SandboxFS implements domain.SandboxFS
var _ domain.SandboxFS = (*SandboxFS)(nil)

// NewSandboxFS creates a new sandbox filesystem
func NewSandboxFS(projectRoot string, realFS domain.FileSystemProvider, log domain.Logger) *SandboxFS {
	return &SandboxFS{
		changes:     make(map[string]*domain.SandboxFileChange),
		projectRoot: projectRoot,
		realFS:      realFS,
		log:         log,
	}
}

// ReadFile reads from sandbox if change exists, otherwise from real FS
func (s *SandboxFS) ReadFile(filename string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedPath := s.normalizePath(filename)

	if change, ok := s.changes[normalizedPath]; ok {
		if change.Operation == domain.SandboxOpDelete {
			return nil, os.ErrNotExist
		}
		return change.Content, nil
	}

	return s.realFS.ReadFile(filename)
}

// WriteFile saves changes to sandbox, not to real filesystem
func (s *SandboxFS) WriteFile(filename string, data []byte, perm int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedPath := s.normalizePath(filename)

	// Read original content for diff (ignore errors - file may not exist)
	original, err := s.realFS.ReadFile(filename)

	op := domain.SandboxOpModify
	if err != nil || original == nil {
		op = domain.SandboxOpCreate
		original = nil
	}

	s.changes[normalizedPath] = &domain.SandboxFileChange{
		Path:            normalizedPath,
		Content:         data,
		OriginalContent: original,
		Operation:       op,
		Timestamp:       time.Now(),
	}

	s.log.Info(fmt.Sprintf("Sandbox: %s file %s", opString(op), normalizedPath))
	return nil
}

// MkdirAll is a no-op in sandbox (directories are virtual)
func (s *SandboxFS) MkdirAll(path string, perm int) error {
	// In sandbox mode, directories are created implicitly when files are written
	return nil
}

// DeleteFile marks a file as deleted in the sandbox
func (s *SandboxFS) DeleteFile(filename string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedPath := s.normalizePath(filename)

	// Read original content for potential restore
	original, err := s.realFS.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("file does not exist: %s", filename)
	}

	s.changes[normalizedPath] = &domain.SandboxFileChange{
		Path:            normalizedPath,
		OriginalContent: original,
		Operation:       domain.SandboxOpDelete,
		Timestamp:       time.Now(),
	}

	s.log.Info(fmt.Sprintf("Sandbox: delete file %s", normalizedPath))
	return nil
}

// GetChanges returns all pending changes sorted by path
func (s *SandboxFS) GetChanges() []*domain.SandboxFileChange {
	s.mu.RLock()
	defer s.mu.RUnlock()

	changes := make([]*domain.SandboxFileChange, 0, len(s.changes))
	for _, change := range s.changes {
		changes = append(changes, change)
	}

	// Sort by path for consistent ordering
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Path < changes[j].Path
	})

	return changes
}

// GetDiff returns unified diff for a single file
func (s *SandboxFS) GetDiff(filename string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	normalizedPath := s.normalizePath(filename)
	change, ok := s.changes[normalizedPath]
	if !ok {
		return "", nil
	}

	return s.generateUnifiedDiff(change), nil
}

// GetAllDiffs returns unified diff for all changes
func (s *SandboxFS) GetAllDiffs() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.changes) == 0 {
		return "", nil
	}

	// Get sorted paths for consistent output
	paths := make([]string, 0, len(s.changes))
	for path := range s.changes {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	var builder strings.Builder
	for i, path := range paths {
		change := s.changes[path]
		diff := s.generateUnifiedDiff(change)
		if diff != "" {
			if i > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(diff)
		}
	}

	return builder.String(), nil
}

// Apply applies all sandbox changes to the real filesystem
func (s *SandboxFS) Apply() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.changes) == 0 {
		return nil
	}

	// Collect all changes to apply
	toApply := make([]*domain.SandboxFileChange, 0, len(s.changes))
	for _, change := range s.changes {
		toApply = append(toApply, change)
	}

	// Sort by path for consistent application order
	sort.Slice(toApply, func(i, j int) bool {
		return toApply[i].Path < toApply[j].Path
	})

	// Apply changes atomically (best effort)
	appliedPaths := make([]string, 0, len(toApply))
	for _, change := range toApply {
		fullPath := s.getFullPath(change.Path)

		switch change.Operation {
		case domain.SandboxOpCreate, domain.SandboxOpModify:
			// Ensure directory exists
			dir := filepath.Dir(fullPath)
			if err := s.realFS.MkdirAll(dir, 0o755); err != nil {
				// Rollback applied changes
				s.rollbackApplied(appliedPaths)
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}

			if err := s.realFS.WriteFile(fullPath, change.Content, 0o644); err != nil {
				// Rollback applied changes
				s.rollbackApplied(appliedPaths)
				return fmt.Errorf("failed to write file %s: %w", change.Path, err)
			}
			appliedPaths = append(appliedPaths, change.Path)
			s.log.Info(fmt.Sprintf("Applied: %s %s", opString(change.Operation), change.Path))

		case domain.SandboxOpDelete:
			if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
				// Rollback applied changes
				s.rollbackApplied(appliedPaths)
				return fmt.Errorf("failed to delete file %s: %w", change.Path, err)
			}
			appliedPaths = append(appliedPaths, change.Path)
			s.log.Info(fmt.Sprintf("Applied: delete %s", change.Path))
		}
	}

	// Clear sandbox after successful apply
	s.changes = make(map[string]*domain.SandboxFileChange)
	s.log.Info(fmt.Sprintf("Sandbox: applied %d changes", len(appliedPaths)))

	return nil
}

// Discard clears all pending changes
func (s *SandboxFS) Discard() {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := len(s.changes)
	s.changes = make(map[string]*domain.SandboxFileChange)
	s.log.Info(fmt.Sprintf("Sandbox: discarded %d changes", count))
}

// DiscardFile removes a single file from pending changes
func (s *SandboxFS) DiscardFile(filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	normalizedPath := s.normalizePath(filename)
	if _, ok := s.changes[normalizedPath]; ok {
		delete(s.changes, normalizedPath)
		s.log.Info(fmt.Sprintf("Sandbox: discarded changes for %s", normalizedPath))
	}
}

// HasChanges returns true if there are pending changes
func (s *SandboxFS) HasChanges() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.changes) > 0
}

// GetChangeCount returns the number of changed files
func (s *SandboxFS) GetChangeCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.changes)
}

// GetProjectRoot returns the project root path
func (s *SandboxFS) GetProjectRoot() string {
	return s.projectRoot
}

// SetProjectRoot updates the project root path
func (s *SandboxFS) SetProjectRoot(root string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.projectRoot = root
}

// normalizePath normalizes file path for consistent storage
func (s *SandboxFS) normalizePath(path string) string {
	// Convert to forward slashes and clean
	cleaned := filepath.ToSlash(filepath.Clean(path))

	// If path is absolute and within project root, make it relative
	if filepath.IsAbs(path) && s.projectRoot != "" {
		rel, err := filepath.Rel(s.projectRoot, path)
		if err == nil && !strings.HasPrefix(rel, "..") {
			return filepath.ToSlash(rel)
		}
	}

	return cleaned
}

// getFullPath returns the full filesystem path for a normalized path
func (s *SandboxFS) getFullPath(normalizedPath string) string {
	if filepath.IsAbs(normalizedPath) {
		return normalizedPath
	}
	return filepath.Join(s.projectRoot, normalizedPath)
}

// generateUnifiedDiff generates unified diff for a file change
func (s *SandboxFS) generateUnifiedDiff(change *domain.SandboxFileChange) string {
	var oldContent, newContent string

	switch change.Operation {
	case domain.SandboxOpCreate:
		oldContent = ""
		newContent = string(change.Content)
	case domain.SandboxOpModify:
		oldContent = string(change.OriginalContent)
		newContent = string(change.Content)
	case domain.SandboxOpDelete:
		oldContent = string(change.OriginalContent)
		newContent = ""
	}

	// Use go-diff for generating diff
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, true)

	// Generate unified diff format
	return s.formatUnifiedDiff(change.Path, change.Operation, oldContent, newContent, diffs)
}

// formatUnifiedDiff formats diffs as unified diff
func (s *SandboxFS) formatUnifiedDiff(path string, op domain.SandboxChangeOperation, oldContent, newContent string, diffs []diffmatchpatch.Diff) string {
	var builder strings.Builder

	// Header
	oldPath := "a/" + path
	newPath := "b/" + path

	switch op {
	case domain.SandboxOpCreate:
		oldPath = "/dev/null"
		builder.WriteString(fmt.Sprintf("diff --git %s %s\n", oldPath, newPath))
		builder.WriteString("new file mode 100644\n")
	case domain.SandboxOpDelete:
		newPath = "/dev/null"
		builder.WriteString(fmt.Sprintf("diff --git %s %s\n", oldPath, newPath))
		builder.WriteString("deleted file mode 100644\n")
	default:
		builder.WriteString(fmt.Sprintf("diff --git %s %s\n", oldPath, newPath))
	}

	builder.WriteString(fmt.Sprintf("--- %s\n", oldPath))
	builder.WriteString(fmt.Sprintf("+++ %s\n", newPath))

	// Generate hunks
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// Simple unified diff generation
	hunks := s.generateHunks(oldLines, newLines, diffs)
	builder.WriteString(hunks)

	return builder.String()
}

// generateHunks generates diff hunks from line-based comparison
func (s *SandboxFS) generateHunks(oldLines, newLines []string, diffs []diffmatchpatch.Diff) string {
	var builder strings.Builder

	// Convert character-based diffs to line-based
	oldLineNum := 1
	newLineNum := 1

	// Simple approach: show all changes in one hunk
	oldLen := len(oldLines)
	newLen := len(newLines)

	if oldLen == 0 && newLen == 0 {
		return ""
	}

	// Handle empty old content (new file)
	if oldLen == 1 && oldLines[0] == "" {
		oldLen = 0
	}
	// Handle empty new content (deleted file)
	if newLen == 1 && newLines[0] == "" {
		newLen = 0
	}

	builder.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", oldLineNum, oldLen, newLineNum, newLen))

	// Show removed lines
	for _, line := range oldLines {
		if line != "" || oldLen > 0 {
			builder.WriteString("-" + line + "\n")
		}
	}

	// Show added lines
	for _, line := range newLines {
		if line != "" || newLen > 0 {
			builder.WriteString("+" + line + "\n")
		}
	}

	return builder.String()
}

// rollbackApplied attempts to rollback already applied changes
func (s *SandboxFS) rollbackApplied(appliedPaths []string) {
	for _, path := range appliedPaths {
		change, ok := s.changes[path]
		if !ok {
			continue
		}

		fullPath := s.getFullPath(path)

		switch change.Operation {
		case domain.SandboxOpCreate:
			// Remove created file
			if err := os.Remove(fullPath); err != nil {
				s.log.Warning(fmt.Sprintf("Rollback: failed to remove %s: %v", path, err))
			}
		case domain.SandboxOpModify:
			// Restore original content
			if change.OriginalContent != nil {
				if err := s.realFS.WriteFile(fullPath, change.OriginalContent, 0o644); err != nil {
					s.log.Warning(fmt.Sprintf("Rollback: failed to restore %s: %v", path, err))
				}
			}
		case domain.SandboxOpDelete:
			// Restore deleted file
			if change.OriginalContent != nil {
				dir := filepath.Dir(fullPath)
				_ = s.realFS.MkdirAll(dir, 0o755)
				if err := s.realFS.WriteFile(fullPath, change.OriginalContent, 0o644); err != nil {
					s.log.Warning(fmt.Sprintf("Rollback: failed to restore deleted %s: %v", path, err))
				}
			}
		}
	}
}

// opString returns string representation of SandboxChangeOperation
func opString(op domain.SandboxChangeOperation) string {
	switch op {
	case domain.SandboxOpCreate:
		return "create"
	case domain.SandboxOpModify:
		return "modify"
	case domain.SandboxOpDelete:
		return "delete"
	default:
		return "unknown"
	}
}
