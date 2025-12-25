package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

func TestGenerateDiff_NoChanges(t *testing.T) {
	// Setup: Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "test-git-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Initialize repository
	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	// Create and commit a file
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0o644)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Make second commit with same content (no actual changes)
	_, err = worktree.Commit("second commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
		AllowEmptyCommits: true,
	})
	require.NoError(t, err)

	// Test - GenerateDiff compares HEAD with HEAD~1
	service := NewService(&domain.NoopLogger{})
	diff, err := service.GenerateDiff(tmpDir)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, diff, "Diff should be empty when commits have no changes")
}

func TestGenerateDiff_WithChanges(t *testing.T) {
	// Setup: Create a temporary git repository
	tmpDir, err := os.MkdirTemp("", "test-git-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Initialize repository
	repo, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)

	// Create and commit a file
	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0o644)
	require.NoError(t, err)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Modify the file and commit again
	err = os.WriteFile(testFile, []byte("modified content"), 0o644)
	require.NoError(t, err)

	_, err = worktree.Add("test.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("second commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
		},
	})
	require.NoError(t, err)

	// Test - GenerateDiff compares HEAD with HEAD~1
	service := NewService(&domain.NoopLogger{})
	diff, err := service.GenerateDiff(tmpDir)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, diff, "Diff should not be empty when commits have changes")
	assert.Contains(t, diff, "test.txt", "Diff should mention the modified file")
}
