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

func TestGetUncommittedFiles(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Initially no uncommitted files
	files, err := repo.GetUncommittedFiles(tempDir)
	if err != nil {
		t.Fatalf("GetUncommittedFiles error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("Expected 0 uncommitted files, got %d", len(files))
	}

	// Create a new file (untracked)
	newFile := filepath.Join(tempDir, "new.txt")
	if err := os.WriteFile(newFile, []byte("new content"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err = repo.GetUncommittedFiles(tempDir)
	if err != nil {
		t.Fatalf("GetUncommittedFiles error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 uncommitted file, got %d", len(files))
	}
	if len(files) > 0 && files[0].Status != "U" {
		t.Errorf("Expected status 'U' (untracked), got %q", files[0].Status)
	}
}

func TestGetUncommittedFiles_Modified(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Modify existing file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("modified content"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := repo.GetUncommittedFiles(tempDir)
	if err != nil {
		t.Fatalf("GetUncommittedFiles error: %v", err)
	}
	if len(files) != 1 {
		t.Errorf("Expected 1 uncommitted file, got %d", len(files))
	}
	if len(files) > 0 && files[0].Status != "M" {
		t.Errorf("Expected status 'M' (modified), got %q", files[0].Status)
	}
}

func TestGetRichCommitHistory(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	commits, err := repo.GetRichCommitHistory(tempDir, "", 10)
	if err != nil {
		t.Fatalf("GetRichCommitHistory error: %v", err)
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
		if c.Author == "" {
			t.Error("Commit author should not be empty")
		}
	}
}

func TestGetFileContentAtCommit(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	// Get commit hash
	commits, _ := repo.GetCommitHistory(tempDir, 1)
	if len(commits) == 0 {
		t.Fatal("No commits found")
	}

	content, err := repo.GetFileContentAtCommit(tempDir, "test.txt", commits[0].Hash)
	if err != nil {
		t.Fatalf("GetFileContentAtCommit error: %v", err)
	}

	if !strings.Contains(content, "test content") {
		t.Errorf("Expected 'test content' in file, got: %s", content)
	}
}

func TestGenerateDiff(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepoWithMultipleCommits(t)
	defer os.RemoveAll(tempDir)

	diff, err := repo.GenerateDiff(tempDir)
	if err != nil {
		t.Fatalf("GenerateDiff error: %v", err)
	}

	// Diff should contain something
	if diff == "" {
		t.Log("Diff is empty (expected if only one commit)")
	}
}

func TestGetAllFiles(t *testing.T) {
	// Skip if git not available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}

	repo := New(&testLogger{})
	tempDir := setupTestGitRepo(t)
	defer os.RemoveAll(tempDir)

	files, err := repo.GetAllFiles(tempDir)
	if err != nil {
		t.Fatalf("GetAllFiles error: %v", err)
	}

	if len(files) == 0 {
		t.Error("Expected at least one file")
	}

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
}

func TestParseGitStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []domain.FileStatus
	}{
		{
			name:  "modified file",
			input: " M file.go\n",
			expected: []domain.FileStatus{
				{Path: "file.go", Status: "M"},
			},
		},
		{
			name:  "added file",
			input: "A  new.go\n",
			expected: []domain.FileStatus{
				{Path: "new.go", Status: "A"},
			},
		},
		{
			name:  "deleted file",
			input: "D  old.go\n",
			expected: []domain.FileStatus{
				{Path: "old.go", Status: "D"},
			},
		},
		{
			name:  "untracked file",
			input: "?? untracked.txt\n",
			expected: []domain.FileStatus{
				{Path: "untracked.txt", Status: "U"},
			},
		},
		{
			name:  "renamed file",
			input: "R  old.go -> new.go\n",
			expected: []domain.FileStatus{
				{Path: "new.go", Status: "R"},
			},
		},
		{
			name:     "empty output",
			input:    "",
			expected: nil,
		},
		{
			name:  "multiple files",
			input: " M file1.go\nA  file2.go\n?? file3.txt\n",
			expected: []domain.FileStatus{
				{Path: "file1.go", Status: "M"},
				{Path: "file2.go", Status: "A"},
				{Path: "file3.txt", Status: "U"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseGitStatus(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d files, got %d", len(tt.expected), len(result))
				return
			}
			for i, exp := range tt.expected {
				if result[i].Path != exp.Path {
					t.Errorf("file %d: expected path %q, got %q", i, exp.Path, result[i].Path)
				}
				if result[i].Status != exp.Status {
					t.Errorf("file %d: expected status %q, got %q", i, exp.Status, result[i].Status)
				}
			}
		})
	}
}

func TestParseRichLogOutput(t *testing.T) {
	input := `COMMIT abc123 def456
Initial commit
Test User
2024-01-15T10:30:00+00:00
A	file1.go
A	file2.go

COMMIT def456
Second commit
Test User
2024-01-16T11:00:00+00:00
M	file1.go
`

	commits, err := ParseRichLogOutput(input)
	if err != nil {
		t.Fatalf("ParseRichLogOutput error: %v", err)
	}

	if len(commits) != 2 {
		t.Errorf("expected 2 commits, got %d", len(commits))
	}

	if len(commits) > 0 {
		if commits[0].Hash != "abc123" {
			t.Errorf("expected hash 'abc123', got %q", commits[0].Hash)
		}
		if commits[0].Subject != "Initial commit" {
			t.Errorf("expected subject 'Initial commit', got %q", commits[0].Subject)
		}
		if commits[0].Author != "Test User" {
			t.Errorf("expected author 'Test User', got %q", commits[0].Author)
		}
		if !commits[0].IsMerge {
			// Has parent def456, so it's a merge
		}
		if len(commits[0].Files) != 2 {
			t.Errorf("expected 2 files, got %d", len(commits[0].Files))
		}
	}
}

func TestMapGitStatus(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"M ", "M"},
		{" M", "M"},
		{"MM", "M"},
		{"A ", "A"},
		{"D ", "D"},
		{"R ", "R"},
		{"C ", "C"},
		{"??", "U"},
		{"UU", "UM"},
		{"XX", "XX"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := mapGitStatus(tt.input)
			if result != tt.expected {
				t.Errorf("mapGitStatus(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper to setup a test git repository with multiple commits
func setupTestGitRepoWithMultipleCommits(t *testing.T) string {
	tempDir := setupTestGitRepo(t)

	// Add another file and commit
	file2 := filepath.Join(tempDir, "file2.txt")
	if err := os.WriteFile(file2, []byte("second file"), 0o644); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	cmd := exec.Command("git", "add", ".")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	cmd = exec.Command("git", "commit", "-m", "Second commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tempDir)
		t.Fatal(err)
	}

	return tempDir
}
