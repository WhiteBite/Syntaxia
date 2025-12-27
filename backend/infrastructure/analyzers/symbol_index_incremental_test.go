package analyzers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSymbolIndex_IncrementalUpdate(t *testing.T) {
	tests := []struct {
		name           string
		initialFiles   map[string]string
		modifiedFiles  map[string]string
		deletedFiles   []string
		expectedChange int
		wantErr        bool
	}{
		{
			name: "no changes",
			initialFiles: map[string]string{
				"main.go": "package main\nfunc Hello() {}",
			},
			modifiedFiles:  nil,
			deletedFiles:   nil,
			expectedChange: 0,
			wantErr:        false,
		},
		{
			name: "single file modified",
			initialFiles: map[string]string{
				"main.go": "package main\nfunc Hello() {}",
			},
			modifiedFiles: map[string]string{
				"main.go": "package main\nfunc Hello() {}\nfunc World() {}",
			},
			expectedChange: 1,
			wantErr:        false,
		},
		{
			name: "new file added",
			initialFiles: map[string]string{
				"main.go": "package main\nfunc Hello() {}",
			},
			modifiedFiles: map[string]string{
				"utils.go": "package main\nfunc Utils() {}",
			},
			expectedChange: 1,
			wantErr:        false,
		},
		{
			name: "file deleted",
			initialFiles: map[string]string{
				"main.go":  "package main\nfunc Hello() {}",
				"utils.go": "package main\nfunc Utils() {}",
			},
			deletedFiles:   []string{"utils.go"},
			expectedChange: 1,
			wantErr:        false,
		},
		{
			name: "multiple changes",
			initialFiles: map[string]string{
				"main.go":   "package main\nfunc Hello() {}",
				"utils.go":  "package main\nfunc Utils() {}",
				"delete.go": "package main\nfunc Delete() {}",
			},
			modifiedFiles: map[string]string{
				"main.go": "package main\nfunc Hello() {}\nfunc Modified() {}",
				"new.go":  "package main\nfunc New() {}",
			},
			deletedFiles:   []string{"delete.go"},
			expectedChange: 3,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			registry := NewAnalyzerRegistry()
			idx := NewSymbolIndex(registry)
			ctx := context.Background()

			// Create initial files
			for name, content := range tt.initialFiles {
				path := filepath.Join(tmpDir, name)
				require.NoError(t, os.WriteFile(path, []byte(content), 0644))
			}

			// Initial index
			err := idx.IndexProject(ctx, tmpDir)
			require.NoError(t, err)

			// Record mod times
			for name := range tt.initialFiles {
				relPath := name
				info, _ := os.Stat(filepath.Join(tmpDir, name))
				idx.mu.Lock()
				idx.fileModTimes[relPath] = info.ModTime()
				idx.mu.Unlock()
			}

			// Wait to ensure different mod times
			time.Sleep(10 * time.Millisecond)

			// Apply modifications
			for name, content := range tt.modifiedFiles {
				path := filepath.Join(tmpDir, name)
				require.NoError(t, os.WriteFile(path, []byte(content), 0644))
			}

			// Delete files
			for _, name := range tt.deletedFiles {
				path := filepath.Join(tmpDir, name)
				_ = os.Remove(path)
			}

			// Run incremental update
			changed, err := idx.IncrementalUpdate(ctx, tmpDir)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedChange, changed)
		})
	}
}

func TestSymbolIndex_UpdateFile(t *testing.T) {
	tests := []struct {
		name            string
		initialContent  string
		updatedContent  string
		expectedSymbols []string
	}{
		{
			name:            "update adds new symbol",
			initialContent:  "package main\nfunc Hello() {}",
			updatedContent:  "package main\nfunc Hello() {}\nfunc World() {}",
			expectedSymbols: []string{"main", "Hello", "World"},
		},
		{
			name:            "update removes symbol",
			initialContent:  "package main\nfunc Hello() {}\nfunc World() {}",
			updatedContent:  "package main\nfunc Hello() {}",
			expectedSymbols: []string{"main", "Hello"},
		},
		{
			name:            "update renames symbol",
			initialContent:  "package main\nfunc OldName() {}",
			updatedContent:  "package main\nfunc NewName() {}",
			expectedSymbols: []string{"main", "NewName"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewAnalyzerRegistry()
			idx := NewSymbolIndex(registry)
			ctx := context.Background()

			// Index initial content
			err := idx.IndexFile(ctx, "test.go", []byte(tt.initialContent))
			require.NoError(t, err)

			// Update file
			err = idx.UpdateFile(ctx, "test.go", []byte(tt.updatedContent))
			require.NoError(t, err)

			// Verify symbols
			symbols := idx.GetSymbolsInFile("test.go")
			symbolNames := make([]string, len(symbols))
			for i, s := range symbols {
				symbolNames[i] = s.Name
			}

			assert.ElementsMatch(t, tt.expectedSymbols, symbolNames)
		})
	}
}

func TestSymbolIndex_RemoveFile(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)
	ctx := context.Background()

	// Index a file
	content := "package main\nfunc ToRemove() {}"
	err := idx.IndexFile(ctx, "remove.go", []byte(content))
	require.NoError(t, err)

	// Verify symbol exists
	symbols := idx.SearchByName("ToRemove")
	assert.Len(t, symbols, 1)

	// Remove file
	idx.RemoveFile("remove.go")

	// Verify symbol is gone
	symbols = idx.SearchByName("ToRemove")
	assert.Len(t, symbols, 0)

	// Verify file is not tracked
	idx.mu.RLock()
	_, exists := idx.fileModTimes["remove.go"]
	idx.mu.RUnlock()
	assert.False(t, exists)
}

func TestSymbolIndex_GetChangedFiles(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	// Create initial file
	file1 := filepath.Join(tmpDir, "file1.go")
	require.NoError(t, os.WriteFile(file1, []byte("package main\nfunc F1() {}"), 0644))

	// Index and record mod time
	ctx := context.Background()
	err := idx.IndexProject(ctx, tmpDir)
	require.NoError(t, err)

	info, _ := os.Stat(file1)
	idx.mu.Lock()
	idx.fileModTimes["file1.go"] = info.ModTime()
	idx.mu.Unlock()

	// No changes yet
	changed, err := idx.GetChangedFiles(tmpDir)
	require.NoError(t, err)
	assert.Len(t, changed, 0)

	// Wait and modify
	time.Sleep(10 * time.Millisecond)
	require.NoError(t, os.WriteFile(file1, []byte("package main\nfunc F1() {}\nfunc F2() {}"), 0644))

	// Should detect change
	changed, err = idx.GetChangedFiles(tmpDir)
	require.NoError(t, err)
	assert.Contains(t, changed, "file1.go")
}

func TestShouldSkipDirectory(t *testing.T) {
	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{"node_modules", "node_modules", true},
		{"vendor", "vendor", true},
		{"build", "build", true},
		{"dist", "dist", true},
		{".git", ".git", true},
		{".idea", ".idea", true},
		{"hidden dir", ".hidden", true},
		{"src", "src", false},
		{"components", "components", false},
		{"utils", "utils", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipDirectory(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSymbolIndex_UpdateIncremental(t *testing.T) {
	tests := []struct {
		name          string
		initialFiles  map[string]string
		update        IncrementalUpdate
		expectSymbols []string
	}{
		{
			name: "add single file",
			initialFiles: map[string]string{
				"main.go": "package main\nfunc Main() {}",
			},
			update: IncrementalUpdate{
				AddedFiles: []string{"utils.go"},
			},
			expectSymbols: []string{"main", "Main"},
		},
		{
			name: "modify file",
			initialFiles: map[string]string{
				"main.go": "package main\nfunc Original() {}",
			},
			update: IncrementalUpdate{
				ModifiedFiles: []string{"main.go"},
			},
			expectSymbols: []string{"main"},
		},
		{
			name: "delete file",
			initialFiles: map[string]string{
				"main.go":  "package main\nfunc Main() {}",
				"utils.go": "package main\nfunc Utils() {}",
			},
			update: IncrementalUpdate{
				DeletedFiles: []string{"utils.go"},
			},
			expectSymbols: []string{"main", "Main"},
		},
		{
			name: "mixed operations",
			initialFiles: map[string]string{
				"a.go": "package main\nfunc A() {}",
				"b.go": "package main\nfunc B() {}",
			},
			update: IncrementalUpdate{
				AddedFiles:    []string{"c.go"},
				ModifiedFiles: []string{"a.go"},
				DeletedFiles:  []string{"b.go"},
			},
			expectSymbols: []string{"main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			registry := NewAnalyzerRegistry()
			idx := NewSymbolIndex(registry)
			ctx := context.Background()

			// Create initial files
			for name, content := range tt.initialFiles {
				path := filepath.Join(tmpDir, name)
				require.NoError(t, os.WriteFile(path, []byte(content), 0644))
			}

			// Initial index
			err := idx.IndexProject(ctx, tmpDir)
			require.NoError(t, err)

			// Create added files
			for _, name := range tt.update.AddedFiles {
				path := filepath.Join(tmpDir, name)
				content := "package main\nfunc Added() {}"
				require.NoError(t, os.WriteFile(path, []byte(content), 0644))
			}

			// Modify files
			for _, name := range tt.update.ModifiedFiles {
				path := filepath.Join(tmpDir, name)
				content := "package main\nfunc Modified() {}"
				require.NoError(t, os.WriteFile(path, []byte(content), 0644))
			}

			// Delete files
			for _, name := range tt.update.DeletedFiles {
				path := filepath.Join(tmpDir, name)
				_ = os.Remove(path)
			}

			// Apply incremental update
			err = idx.UpdateIncremental(ctx, tt.update)
			assert.NoError(t, err)
		})
	}
}

func TestSymbolIndex_GetFileHash(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)
	ctx := context.Background()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.go")
	content := "package main\nfunc Test() {}"
	require.NoError(t, os.WriteFile(testFile, []byte(content), 0644))

	// Index project
	err := idx.IndexProject(ctx, tmpDir)
	require.NoError(t, err)

	// Record mod time and set project root
	info, _ := os.Stat(testFile)
	idx.mu.Lock()
	idx.fileModTimes["test.go"] = info.ModTime()
	idx.projectRoot = tmpDir
	idx.mu.Unlock()

	// Get hash
	hash := idx.GetFileHash("test.go")
	assert.NotEmpty(t, hash)

	// Hash for non-existent file should be empty
	hash = idx.GetFileHash("nonexistent.go")
	assert.Empty(t, hash)
}

func TestSymbolIndex_SetFileHash(t *testing.T) {
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)

	// SetFileHash is a no-op for SymbolIndexImpl
	// Just verify it doesn't panic
	idx.SetFileHash("test.go", "somehash")
}

func TestIncrementalUpdate_EmptyUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	registry := NewAnalyzerRegistry()
	idx := NewSymbolIndex(registry)
	ctx := context.Background()

	// Create initial file
	testFile := filepath.Join(tmpDir, "main.go")
	require.NoError(t, os.WriteFile(testFile, []byte("package main\nfunc Main() {}"), 0644))

	// Index project
	err := idx.IndexProject(ctx, tmpDir)
	require.NoError(t, err)

	initialCount := len(idx.symbols)

	// Apply empty update
	err = idx.UpdateIncremental(ctx, IncrementalUpdate{})
	assert.NoError(t, err)

	// Symbol count should remain the same
	assert.Equal(t, initialCount, len(idx.symbols))
}
