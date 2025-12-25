package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"testing"
)

type testLogger struct{}

func (l *testLogger) Debug(msg string)   {}
func (l *testLogger) Info(msg string)    {}
func (l *testLogger) Warning(msg string) {}
func (l *testLogger) Error(msg string)   {}
func (l *testLogger) Fatal(msg string)   {}

func TestIsGitAvailable(t *testing.T) {
	repo := New(&testLogger{})

	// Git should be available on most dev machines
	available := repo.IsGitAvailable()

	// Just check it doesn't panic - result depends on environment
	t.Logf("Git available: %v", available)
}

func TestIsGitRepository(t *testing.T) {
	repo := New(&testLogger{})

	// Create temp dir without git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Should return false for non-git directory
	if repo.IsGitRepository(tempDir) {
		t.Error("Expected false for non-git directory")
	}
}

func TestIsGitRepository_WithGit(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	// Create temp dir and init git
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Init git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Should return true for git directory
	if !repo.IsGitRepository(tempDir) {
		t.Error("Expected true for git directory")
	}
}

func TestGetBranches(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	// Create temp git repo with a commit
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	branches, err := repo.GetBranches(tempDir)
	if err != nil {
		t.Fatalf("GetBranches error: %v", err)
	}

	// Should have at least one branch (main or master)
	if len(branches) == 0 {
		t.Error("Expected at least one branch")
	}

	t.Logf("Branches: %v", branches)
}

func TestGetCurrentBranch(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	branch, err := repo.GetCurrentBranch(tempDir)
	if err != nil {
		t.Fatalf("GetCurrentBranch error: %v", err)
	}

	if branch == "" {
		t.Error("Expected non-empty branch name")
	}

	t.Logf("Current branch: %s", branch)
}

func TestGetCommitHistory(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	commits, err := repo.GetCommitHistory(tempDir, 10)
	if err != nil {
		t.Fatalf("GetCommitHistory error: %v", err)
	}

	if len(commits) == 0 {
		t.Error("Expected at least one commit")
	}

	// Check commit structure
	for _, c := range commits {
		if c.Hash == "" {
			t.Error("Commit hash should not be empty")
		}
		if c.Subject == "" {
			t.Error("Commit subject should not be empty")
		}
	}

	t.Logf("Commits: %d", len(commits))
}

func TestListFilesAtRef(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Get current branch
	branch, _ := repo.GetCurrentBranch(tempDir)

	files, err := repo.ListFilesAtRef(tempDir, branch)
	if err != nil {
		t.Fatalf("ListFilesAtRef error: %v", err)
	}

	// Should have test.txt from setup
	found := false
	for _, f := range files {
		if f == "test.txt" {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected test.txt in files list")
	}

	t.Logf("Files at %s: %v", branch, files)
}

func TestGetFileAtRef(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	branch, _ := repo.GetCurrentBranch(tempDir)

	content, err := repo.GetFileAtRef(tempDir, "test.txt", branch)
	if err != nil {
		t.Fatalf("GetFileAtRef error: %v", err)
	}

	if !strings.Contains(content, "test content") {
		t.Errorf("Expected 'test content' in file, got: %s", content)
	}
}

func TestCheckoutBranch(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})

	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Create a new branch
	cmd := exec.Command("git", "branch", "test-branch")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	// Checkout the new branch
	err := repo.CheckoutBranch(tempDir, "test-branch")
	if err != nil {
		t.Fatalf("CheckoutBranch error: %v", err)
	}

	// Verify we're on the new branch
	branch, _ := repo.GetCurrentBranch(tempDir)
	if branch != "test-branch" {
		t.Errorf("Expected 'test-branch', got '%s'", branch)
	}
}

// Helper to setup a test git repository
func setupTestGitRepo(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatal(err)
	}

	// Init git
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	// Configure git user for commits
	cmd = exec.Command("git", "config", "user.email", "test@test.com")
	cmd.Dir = tempDir
	cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0o644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	// Add and commit
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	return tempDir
}

// Verify Repository implements GitRepository interface
var _ domain.GitRepository = (*Repository)(nil)
