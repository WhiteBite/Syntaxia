package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	t.Helper()
	tmpDir := t.TempDir()

	runGit(t, tmpDir, "init")
	runGit(t, tmpDir, "config", "user.email", "test@test.com")
	runGit(t, tmpDir, "config", "user.name", "Test User")

	writeFile(t, tmpDir, "README.md", "# Test")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Initial commit")

	return tmpDir
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	os.MkdirAll(filepath.Dir(path), 0o755)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", name, err)
	}
}

func TestNewContextBuilder(t *testing.T) {
	cb := NewContextBuilder("/tmp/test")
	if cb == nil {
		t.Fatal("NewContextBuilder returned nil")
	}
	if cb.projectRoot != "/tmp/test" {
		t.Errorf("expected projectRoot '/tmp/test', got %q", cb.projectRoot)
	}
}

func TestContextBuilder_GetRecentChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Add more commits
	writeFile(t, tmpDir, "file1.go", "package main")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Add file1")

	writeFile(t, tmpDir, "file1.go", "package main\n// updated")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Update file1")

	cb := NewContextBuilder(tmpDir)
	changes, err := cb.GetRecentChanges("1 year ago", "")
	if err != nil {
		t.Fatalf("GetRecentChanges failed: %v", err)
	}

	if len(changes) == 0 {
		t.Error("expected some changes")
	}

	// file1.go should have higher change count
	foundFile1 := false
	for _, c := range changes {
		if c.FilePath == "file1.go" {
			foundFile1 = true
			if c.ChangeCount < 2 {
				t.Errorf("expected file1.go to have at least 2 changes, got %d", c.ChangeCount)
			}
		}
	}
	if !foundFile1 {
		t.Error("file1.go not found in changes")
	}
}

func TestContextBuilder_GetCoChangedFiles(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create files that change together
	writeFile(t, tmpDir, "service.go", "package main")
	writeFile(t, tmpDir, "service_test.go", "package main")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Add service")

	writeFile(t, tmpDir, "service.go", "package main\n// v2")
	writeFile(t, tmpDir, "service_test.go", "package main\n// v2")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Update service")

	cb := NewContextBuilder(tmpDir)
	coChanged, err := cb.GetCoChangedFiles("service.go", 10)
	if err != nil {
		t.Fatalf("GetCoChangedFiles failed: %v", err)
	}

	found := false
	for _, f := range coChanged {
		if f == "service_test.go" {
			found = true
			break
		}
	}
	if !found {
		t.Error("service_test.go should be co-changed with service.go")
	}
}

func TestContextBuilder_GetRelatedByAuthor(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Same author changes multiple files
	writeFile(t, tmpDir, "auth.go", "package auth")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Add auth")

	writeFile(t, tmpDir, "user.go", "package user")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Add user")

	cb := NewContextBuilder(tmpDir)
	related, err := cb.GetRelatedByAuthor("auth.go", 10)
	if err != nil {
		t.Fatalf("GetRelatedByAuthor failed: %v", err)
	}

	// Should find other files by same author
	t.Logf("Related files: %v", related)
}

func TestContextBuilder_SuggestContextFiles(t *testing.T) {
	tmpDir := setupGitRepo(t)

	writeFile(t, tmpDir, "auth.go", "package auth")
	writeFile(t, tmpDir, "auth_test.go", "package auth")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "feat: add authentication")

	cb := NewContextBuilder(tmpDir)
	suggestions, err := cb.SuggestContextFiles("fix auth bug", []string{"auth.go"}, 10)
	if err != nil {
		t.Fatalf("SuggestContextFiles failed: %v", err)
	}

	t.Logf("Suggestions: %v", suggestions)
}

func TestParseUnixTime(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1234567890", 1234567890},
		{"0", 0},
		{"999", 999},
	}

	for _, tt := range tests {
		var ts int64
		parseUnixTime(tt.input, &ts)
		if ts != tt.expected {
			t.Errorf("parseUnixTime(%q) = %d, want %d", tt.input, ts, tt.expected)
		}
	}
}

func TestParseUnixTime_Reset(t *testing.T) {
	// Test that parseUnixTime resets the value
	var ts int64 = 999
	parseUnixTime("123", &ts)
	if ts != 123 {
		t.Errorf("expected 123, got %d (should reset previous value)", ts)
	}
}

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		input    string
		minWords int
	}{
		{"fix authentication bug in login", 2},
		{"add user registration feature", 2},
		{"the a an is are", 0}, // all stop words
	}

	for _, tt := range tests {
		keywords := extractKeywords(tt.input)
		if len(keywords) < tt.minWords {
			t.Errorf("extractKeywords(%q) = %v, expected at least %d words", tt.input, keywords, tt.minWords)
		}
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !containsString(slice, "a") {
		t.Error("expected to find 'a'")
	}
	if !containsString(slice, "c") {
		t.Error("expected to find 'c'")
	}
	if containsString(slice, "d") {
		t.Error("should not find 'd'")
	}
	if containsString(nil, "a") {
		t.Error("should not find in nil slice")
	}
}
