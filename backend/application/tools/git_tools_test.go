package tools

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	tmpDir := t.TempDir()

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skip("git not available")
	}

	// Configure git user
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	return tmpDir
}

func TestGitStatus_ReturnsChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create a file
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	// Should show untracked file
	if !containsStr(result, "test.txt") {
		t.Errorf("expected to see test.txt in status, got: %s", result)
	}
}

func TestGitStatus_CleanRepo(t *testing.T) {
	tmpDir := setupGitRepo(t)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "clean") {
		t.Errorf("expected clean status, got: %s", result)
	}
}

func TestGitDiff_NoChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "No differences") {
		t.Errorf("expected no differences, got: %s", result)
	}
}

func TestGitLog_EmptyRepo(t *testing.T) {
	tmpDir := setupGitRepo(t)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{"limit": float64(5)}, tmpDir)

	// Empty repo might return error or "No commits"
	if err == nil && !containsStr(result, "No commits") && result != "" {
		// If there's a result, it should be valid
		t.Logf("git log result: %s", result)
	}
}

func TestGitLog_WithCommits(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create and commit a file
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{"limit": float64(5)}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsStr(result, "Initial commit") {
		t.Errorf("expected to see commit message, got: %s", result)
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || containsStrHelper(s, substr))
}

func containsStrHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
