package changes

import (
	"fmt"
	"sync"
	"time"

	"syntaxia/domain"
)

// ConflictResolution represents how to resolve a conflict
type ConflictResolution string

const (
	// ResolutionKeepOurs keeps our (sandbox) version
	ResolutionKeepOurs ConflictResolution = "keep_ours"
	// ResolutionKeepTheirs keeps their (filesystem) version
	ResolutionKeepTheirs ConflictResolution = "keep_theirs"
	// ResolutionMerge attempts to merge changes (not implemented yet)
	ResolutionMerge ConflictResolution = "merge"
)

// Conflict represents a detected file conflict
type Conflict struct {
	FilePath       string    `json:"filePath"`
	ExpectedHash   string    `json:"expectedHash"`
	ActualHash     string    `json:"actualHash"`
	Message        string    `json:"message"`
	OurContent     []byte    `json:"ourContent,omitempty"`
	TheirContent   []byte    `json:"theirContent,omitempty"`
	OriginalContent []byte   `json:"originalContent,omitempty"`
	DetectedAt     time.Time `json:"detectedAt"`
}

// ConflictDetector detects and manages file conflicts
type ConflictDetector struct {
	mu         sync.RWMutex
	fileHashes map[string]string // path -> hash at snapshot time
	realFS     domain.FileSystemProvider
	log        domain.Logger
}

// NewConflictDetector creates a new conflict detector
func NewConflictDetector(realFS domain.FileSystemProvider, log domain.Logger) *ConflictDetector {
	return &ConflictDetector{
		fileHashes: make(map[string]string),
		realFS:     realFS,
		log:        log,
	}
}

// SnapshotFile captures the current hash of a file for later conflict detection
func (d *ConflictDetector) SnapshotFile(path string) error {
	content, err := d.realFS.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file for snapshot: %w", err)
	}

	hash := computeHash(content)

	d.mu.Lock()
	d.fileHashes[path] = hash
	d.mu.Unlock()

	d.log.Debug(fmt.Sprintf("Snapshotted file %s with hash %s", path, hash[:8]))

	return nil
}

// SnapshotFiles captures hashes for multiple files
func (d *ConflictDetector) SnapshotFiles(paths []string) error {
	var errs []error

	for _, path := range paths {
		if err := d.SnapshotFile(path); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", path, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to snapshot %d files", len(errs))
	}

	return nil
}

// CheckConflict checks if a single file has been modified since snapshot
func (d *ConflictDetector) CheckConflict(path string) (*Conflict, error) {
	d.mu.RLock()
	expectedHash, hasSnapshot := d.fileHashes[path]
	d.mu.RUnlock()

	if !hasSnapshot {
		return nil, nil // No snapshot, no conflict to detect
	}

	currentContent, err := d.realFS.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for conflict check: %w", err)
	}

	actualHash := computeHash(currentContent)

	if actualHash != expectedHash {
		return &Conflict{
			FilePath:     path,
			ExpectedHash: expectedHash,
			ActualHash:   actualHash,
			Message:      fmt.Sprintf("file %s was modified externally since it was read", path),
			TheirContent: currentContent,
			DetectedAt:   time.Now(),
		}, nil
	}

	return nil, nil
}

// CheckAllConflicts checks all snapshotted files for conflicts
func (d *ConflictDetector) CheckAllConflicts() []Conflict {
	d.mu.RLock()
	paths := make([]string, 0, len(d.fileHashes))
	for path := range d.fileHashes {
		paths = append(paths, path)
	}
	d.mu.RUnlock()

	var conflicts []Conflict

	for _, path := range paths {
		conflict, err := d.CheckConflict(path)
		if err != nil {
			d.log.Warning(fmt.Sprintf("Error checking conflict for %s: %v", path, err))
			continue
		}
		if conflict != nil {
			conflicts = append(conflicts, *conflict)
		}
	}

	return conflicts
}

// ResolveConflict resolves a conflict using the specified resolution strategy
func (d *ConflictDetector) ResolveConflict(path string, resolution ConflictResolution) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	switch resolution {
	case ResolutionKeepOurs:
		// Update snapshot to current filesystem state (we'll overwrite it anyway)
		content, err := d.realFS.ReadFile(path)
		if err != nil {
			// File might not exist yet if it's a create operation
			delete(d.fileHashes, path)
			return nil
		}
		d.fileHashes[path] = computeHash(content)
		d.log.Info(fmt.Sprintf("Resolved conflict for %s: keeping ours", path))

	case ResolutionKeepTheirs:
		// Update snapshot to current filesystem state (accept their changes)
		content, err := d.realFS.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file for resolution: %w", err)
		}
		d.fileHashes[path] = computeHash(content)
		d.log.Info(fmt.Sprintf("Resolved conflict for %s: keeping theirs", path))

	case ResolutionMerge:
		return fmt.Errorf("merge resolution not implemented yet")

	default:
		return fmt.Errorf("unknown resolution strategy: %s", resolution)
	}

	return nil
}

// ClearSnapshot removes the snapshot for a file
func (d *ConflictDetector) ClearSnapshot(path string) {
	d.mu.Lock()
	delete(d.fileHashes, path)
	d.mu.Unlock()
}

// ClearAllSnapshots removes all file snapshots
func (d *ConflictDetector) ClearAllSnapshots() {
	d.mu.Lock()
	d.fileHashes = make(map[string]string)
	d.mu.Unlock()
}

// GetSnapshotCount returns the number of snapshotted files
func (d *ConflictDetector) GetSnapshotCount() int {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return len(d.fileHashes)
}

// HasSnapshot checks if a file has been snapshotted
func (d *ConflictDetector) HasSnapshot(path string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.fileHashes[path]
	return exists
}
