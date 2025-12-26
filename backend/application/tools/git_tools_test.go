package tools

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"syntaxia/domain"
)

// MockGitContextBuilder implements domain.GitContextBuilder for testing
type MockGitContextBuilder struct {
	recentChanges  []domain.RecentChange
	coChangedFiles []string
	suggestFiles   []string
	err            error
}

func (m *MockGitContextBuilder) GetRecentChanges(since, pathFilter string) ([]domain.RecentChange, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.recentChanges, nil
}

func (m *MockGitContextBuilder) GetCoChangedFiles(filePath string, limit int) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.coChangedFiles, nil
}

func (m *MockGitContextBuilder) SuggestContextFiles(task string, currentFiles []string, limit int) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.suggestFiles, nil
}

func (m *MockGitContextBuilder) GetRelatedByAuthor(filePath string, limit int) ([]string, error) {
	return nil, nil
}

func setupGitRepo(t *testing.T) string {
	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skip("git not available")
	}

	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()

	return tmpDir
}

func TestGitStatus_ReturnsChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !strings.Contains(result, "test.txt") {
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
	if !strings.Contains(result, "clean") {
		t.Errorf("expected clean status, got: %s", result)
	}
}

func TestGitStatus_ModifiedFile(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("modified"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "modified") {
		t.Errorf("expected modified status, got: %s", result)
	}
}

func TestGitDiff_NoChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No differences") {
		t.Errorf("expected no differences, got: %s", result)
	}
}

func TestGitDiff_WithChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("modified content"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "No differences") {
		t.Errorf("expected differences, got: %s", result)
	}
}

func TestGitDiff_Staged(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("staged content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{"staged": true}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(result, "No differences") {
		t.Errorf("expected staged differences, got: %s", result)
	}
}

func TestGitDiff_WithPath(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("initial1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("initial2"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()
	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("modified1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("modified2"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{"path": "file1.txt"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "file1") {
		t.Errorf("expected diff for file1, got: %s", result)
	}
}

func TestGitLog_EmptyRepo(t *testing.T) {
	tmpDir := setupGitRepo(t)

	handler := NewGitToolsHandler(nil, nil)
	_, err := handler.Execute("git_log", map[string]any{"limit": float64(5)}, tmpDir)

	// Empty repo returns error or "No commits" - both are acceptable
	if err != nil {
		t.Logf("git log on empty repo returned error (expected): %v", err)
	}
}

func TestGitLog_WithCommits(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{"limit": float64(5)}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Initial commit") {
		t.Errorf("expected to see commit message, got: %s", result)
	}
}

func TestGitLog_WithLimit(t *testing.T) {
	tmpDir := setupGitRepo(t)

	for i := 1; i <= 5; i++ {
		os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"+string(rune('0'+i))), 0644)
		exec.Command("git", "-C", tmpDir, "add", ".").Run()
		exec.Command("git", "-C", tmpDir, "commit", "-m", "Commit "+string(rune('0'+i))).Run()
	}

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{"limit": float64(2)}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(result), "\n")
	// Header line + 2 commits
	if len(lines) < 2 {
		t.Errorf("expected at least 2 lines, got: %d", len(lines))
	}
}

func TestGitLog_WithPath(t *testing.T) {
	tmpDir := setupGitRepo(t)

	os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("content1"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Add file1").Run()

	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("content2"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Add file2").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{
		"limit": float64(10),
		"path":  "file1.txt",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "file1") {
		t.Errorf("expected commits for file1, got: %s", result)
	}
}

func TestGitChangedFiles_Success(t *testing.T) {
	mock := &MockGitContextBuilder{
		recentChanges: []domain.RecentChange{
			{FilePath: "src/main.go", ChangeCount: 5, LastChanged: time.Now()},
			{FilePath: "src/utils.go", ChangeCount: 3, LastChanged: time.Now()},
		},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_changed_files", map[string]any{
		"since": "1 week ago",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "main.go") || !strings.Contains(result, "utils.go") {
		t.Errorf("expected changed files, got: %s", result)
	}
}

func TestGitChangedFiles_NoChanges(t *testing.T) {
	mock := &MockGitContextBuilder{
		recentChanges: []domain.RecentChange{},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_changed_files", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No recent changes") {
		t.Errorf("expected 'No recent changes', got: %s", result)
	}
}

func TestGitChangedFiles_NotInitialized(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)
	_, err := handler.Execute("git_changed_files", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error when git context not initialized")
	}
}

func TestGitCoChanged_Success(t *testing.T) {
	mock := &MockGitContextBuilder{
		coChangedFiles: []string{"related1.go", "related2.go"},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_co_changed", map[string]any{
		"file_path": "main.go",
		"limit":     float64(5),
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "related1.go") {
		t.Errorf("expected co-changed files, got: %s", result)
	}
}

func TestGitCoChanged_MissingFilePath(t *testing.T) {
	mock := &MockGitContextBuilder{}
	handler := NewGitToolsHandler(nil, mock)

	_, err := handler.Execute("git_co_changed", map[string]any{}, "/project")
	if err == nil {
		t.Fatal("expected error for missing file_path")
	}
}

func TestGitCoChanged_NoResults(t *testing.T) {
	mock := &MockGitContextBuilder{
		coChangedFiles: []string{},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_co_changed", map[string]any{
		"file_path": "isolated.go",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No co-changed") {
		t.Errorf("expected 'No co-changed' message, got: %s", result)
	}
}

func TestGitSuggestContext_Success(t *testing.T) {
	mock := &MockGitContextBuilder{
		suggestFiles: []string{"auth.go", "user.go", "session.go"},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_suggest_context", map[string]any{
		"task":          "implement login feature",
		"current_files": []any{"main.go"},
		"limit":         float64(5),
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "auth.go") {
		t.Errorf("expected suggested files, got: %s", result)
	}
}

func TestGitSuggestContext_NoSuggestions(t *testing.T) {
	mock := &MockGitContextBuilder{
		suggestFiles: []string{},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_suggest_context", map[string]any{
		"task": "unknown task",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No context suggestions") {
		t.Errorf("expected 'No context suggestions', got: %s", result)
	}
}

func TestGitToolsHandler_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     bool
	}{
		{"git_status", "git_status", true},
		{"git_diff", "git_diff", true},
		{"git_log", "git_log", true},
		{"git_changed_files", "git_changed_files", true},
		{"git_co_changed", "git_co_changed", true},
		{"git_suggest_context", "git_suggest_context", true},
		{"unknown_tool", "unknown_tool", false},
		{"read_file", "read_file", false},
	}

	handler := NewGitToolsHandler(nil, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.CanHandle(tt.toolName); got != tt.want {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

func TestGitToolsHandler_GetTools(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)
	tools := handler.GetTools()

	if len(tools) == 0 {
		t.Fatal("expected non-empty tools list")
	}

	expectedTools := []string{"git_status", "git_diff", "git_log", "git_changed_files", "git_co_changed", "git_suggest_context"}
	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q not found", expected)
		}
	}
}

func TestGitToolsHandler_UnknownTool(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)
	_, err := handler.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown git tool") {
		t.Errorf("expected 'unknown git tool' error, got: %v", err)
	}
}

// Tests for input validation

func TestGitToolsHandler_ValidateArgs(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name     string
		args     map[string]any
		required []string
		wantErr  bool
		errMsg   string
	}{
		{
			name:     "all required present",
			args:     map[string]any{"file_path": "test.go", "limit": 10},
			required: []string{"file_path"},
			wantErr:  false,
		},
		{
			name:     "missing required",
			args:     map[string]any{"limit": 10},
			required: []string{"file_path"},
			wantErr:  true,
			errMsg:   "missing required argument: file_path",
		},
		{
			name:     "nil value",
			args:     map[string]any{"file_path": nil},
			required: []string{"file_path"},
			wantErr:  true,
			errMsg:   "argument file_path cannot be nil",
		},
		{
			name:     "empty required list",
			args:     map[string]any{},
			required: []string{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateArgs(tt.args, tt.required)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGitToolsHandler_GetStringArg(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name       string
		args       map[string]any
		key        string
		defaultVal string
		want       string
	}{
		{"existing string", map[string]any{"path": "test.go"}, "path", "", "test.go"},
		{"missing key", map[string]any{}, "path", "default.go", "default.go"},
		{"wrong type", map[string]any{"path": 123}, "path", "default.go", "default.go"},
		{"nil value", map[string]any{"path": nil}, "path", "default.go", "default.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.getStringArg(tt.args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getStringArg() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGitToolsHandler_GetIntArg(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name       string
		args       map[string]any
		key        string
		defaultVal int
		want       int
	}{
		{"float64 value", map[string]any{"limit": float64(25)}, "limit", 10, 25},
		{"int value", map[string]any{"limit": 30}, "limit", 10, 30},
		{"missing key", map[string]any{}, "limit", 10, 10},
		{"wrong type", map[string]any{"limit": "abc"}, "limit", 10, 10},
		{"nil value", map[string]any{"limit": nil}, "limit", 10, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.getIntArg(tt.args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getIntArg() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGitToolsHandler_GetBoolArg(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name       string
		args       map[string]any
		key        string
		defaultVal bool
		want       bool
	}{
		{"true value", map[string]any{"staged": true}, "staged", false, true},
		{"false value", map[string]any{"staged": false}, "staged", true, false},
		{"missing key", map[string]any{}, "staged", true, true},
		{"wrong type", map[string]any{"staged": "yes"}, "staged", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handler.getBoolArg(tt.args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("getBoolArg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitToolsHandler_GetStringArrayArg(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name    string
		args    map[string]any
		key     string
		want    []string
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid array",
			args: map[string]any{"files": []any{"a.go", "b.go"}},
			key:  "files",
			want: []string{"a.go", "b.go"},
		},
		{
			name: "missing key",
			args: map[string]any{},
			key:  "files",
			want: nil,
		},
		{
			name: "nil value",
			args: map[string]any{"files": nil},
			key:  "files",
			want: nil,
		},
		{
			name:    "not an array",
			args:    map[string]any{"files": "not-array"},
			key:     "files",
			wantErr: true,
			errMsg:  "must be an array",
		},
		{
			name:    "array with non-string",
			args:    map[string]any{"files": []any{"a.go", 123}},
			key:     "files",
			wantErr: true,
			errMsg:  "must be a string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handler.getStringArrayArg(tt.args, tt.key)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("expected error containing %q, got: %v", tt.errMsg, err)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != len(tt.want) {
					t.Errorf("getStringArrayArg() len = %d, want %d", len(got), len(tt.want))
				}
				for i := range got {
					if got[i] != tt.want[i] {
						t.Errorf("getStringArrayArg()[%d] = %q, want %q", i, got[i], tt.want[i])
					}
				}
			}
		})
	}
}

func TestGitLog_LimitValidation(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create a commit
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	handler := NewGitToolsHandler(nil, nil)

	tests := []struct {
		name  string
		limit any
	}{
		{"negative limit clamped to 1", float64(-5)},
		{"zero limit clamped to 1", float64(0)},
		{"over 100 clamped to 100", float64(200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Execute("git_log", map[string]any{"limit": tt.limit}, tmpDir)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGitCoChanged_NilFilePath(t *testing.T) {
	mock := &MockGitContextBuilder{}
	handler := NewGitToolsHandler(nil, mock)

	_, err := handler.Execute("git_co_changed", map[string]any{"file_path": nil}, "/project")
	if err == nil {
		t.Fatal("expected error for nil file_path")
	}
	if !strings.Contains(err.Error(), "cannot be nil") {
		t.Errorf("expected 'cannot be nil' error, got: %v", err)
	}
}

func TestGitSuggestContext_InvalidCurrentFiles(t *testing.T) {
	mock := &MockGitContextBuilder{}
	handler := NewGitToolsHandler(nil, mock)

	_, err := handler.Execute("git_suggest_context", map[string]any{
		"current_files": "not-an-array",
	}, "/project")

	if err == nil {
		t.Fatal("expected error for invalid current_files type")
	}
	if !strings.Contains(err.Error(), "must be an array") {
		t.Errorf("expected 'must be an array' error, got: %v", err)
	}
}

func TestGitToolsHandler_NilArgs(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	// Should not panic with nil args
	_, err := handler.Execute("git_status", nil, t.TempDir())
	// git_status doesn't require args, so it should work (or fail for other reasons like no git repo)
	_ = err // We just want to ensure no panic
}


// === Additional Edge Case Tests ===

func TestGitDiff_Truncation(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create initial commit
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Create large change
	largeContent := strings.Repeat("modified line\n", 1000)
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(largeContent), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Large diffs should be truncated
	if len(result) > maxDiffSize+100 { // Allow some margin for truncation message
		t.Errorf("diff should be truncated, got length: %d", len(result))
	}
}

func TestGitStatus_AddedFile(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create initial commit
	os.WriteFile(filepath.Join(tmpDir, "initial.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Add new file and stage it
	os.WriteFile(filepath.Join(tmpDir, "new.txt"), []byte("new"), 0644)
	exec.Command("git", "-C", tmpDir, "add", "new.txt").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "added") || !strings.Contains(result, "new.txt") {
		t.Errorf("expected added status for new.txt, got: %s", result)
	}
}

func TestGitStatus_DeletedFile(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create and commit file
	os.WriteFile(filepath.Join(tmpDir, "todelete.txt"), []byte("content"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Delete and stage
	os.Remove(filepath.Join(tmpDir, "todelete.txt"))
	exec.Command("git", "-C", tmpDir, "add", "todelete.txt").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "deleted") {
		t.Errorf("expected deleted status, got: %s", result)
	}
}

func TestGitChangedFiles_WithPathFilter(t *testing.T) {
	mock := &MockGitContextBuilder{
		recentChanges: []domain.RecentChange{
			{FilePath: "src/main.go", ChangeCount: 5, LastChanged: time.Now()},
		},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_changed_files", map[string]any{
		"since":       "1 week ago",
		"path_filter": "src/",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "main.go") {
		t.Errorf("expected filtered results, got: %s", result)
	}
}

func TestGitChangedFiles_Error(t *testing.T) {
	mock := &MockGitContextBuilder{
		err: errors.New("git error"),
	}

	handler := NewGitToolsHandler(nil, mock)
	_, err := handler.Execute("git_changed_files", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error from git context")
	}
	if !strings.Contains(err.Error(), "failed to get recent changes") {
		t.Errorf("expected 'failed to get recent changes' error, got: %v", err)
	}
}

func TestGitCoChanged_Error(t *testing.T) {
	mock := &MockGitContextBuilder{
		err: errors.New("git error"),
	}

	handler := NewGitToolsHandler(nil, mock)
	_, err := handler.Execute("git_co_changed", map[string]any{
		"file_path": "test.go",
	}, "/project")

	if err == nil {
		t.Fatal("expected error from git context")
	}
	if !strings.Contains(err.Error(), "failed to get co-changed files") {
		t.Errorf("expected 'failed to get co-changed files' error, got: %v", err)
	}
}

func TestGitCoChanged_EmptyFilePath(t *testing.T) {
	mock := &MockGitContextBuilder{}
	handler := NewGitToolsHandler(nil, mock)

	_, err := handler.Execute("git_co_changed", map[string]any{
		"file_path": "",
	}, "/project")

	if err == nil {
		t.Fatal("expected error for empty file_path")
	}
	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("expected 'cannot be empty' error, got: %v", err)
	}
}

func TestGitCoChanged_LimitValidation(t *testing.T) {
	mock := &MockGitContextBuilder{
		coChangedFiles: []string{"file1.go"},
	}

	handler := NewGitToolsHandler(nil, mock)

	tests := []struct {
		name  string
		limit any
	}{
		{"negative limit", float64(-5)},
		{"zero limit", float64(0)},
		{"over max limit", float64(200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Execute("git_co_changed", map[string]any{
				"file_path": "test.go",
				"limit":     tt.limit,
			}, "/project")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGitSuggestContext_Error(t *testing.T) {
	mock := &MockGitContextBuilder{
		err: errors.New("suggestion error"),
	}

	handler := NewGitToolsHandler(nil, mock)
	_, err := handler.Execute("git_suggest_context", map[string]any{
		"task": "test task",
	}, "/project")

	if err == nil {
		t.Fatal("expected error from git context")
	}
	if !strings.Contains(err.Error(), "failed to suggest context files") {
		t.Errorf("expected 'failed to suggest context files' error, got: %v", err)
	}
}

func TestGitSuggestContext_NotInitialized(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	_, err := handler.Execute("git_suggest_context", map[string]any{
		"task": "test",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when git context not initialized")
	}
	if !strings.Contains(err.Error(), "git context not initialized") {
		t.Errorf("expected 'git context not initialized' error, got: %v", err)
	}
}

func TestGitCoChanged_NotInitialized(t *testing.T) {
	handler := NewGitToolsHandler(nil, nil)

	_, err := handler.Execute("git_co_changed", map[string]any{
		"file_path": "test.go",
	}, "/project")

	if err == nil {
		t.Fatal("expected error when git context not initialized")
	}
	if !strings.Contains(err.Error(), "git context not initialized") {
		t.Errorf("expected 'git context not initialized' error, got: %v", err)
	}
}

func TestGitSuggestContext_LimitValidation(t *testing.T) {
	mock := &MockGitContextBuilder{
		suggestFiles: []string{"file1.go"},
	}

	handler := NewGitToolsHandler(nil, mock)

	tests := []struct {
		name  string
		limit any
	}{
		{"negative limit", float64(-5)},
		{"zero limit", float64(0)},
		{"over max limit", float64(200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.Execute("git_suggest_context", map[string]any{
				"task":  "test",
				"limit": tt.limit,
			}, "/project")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestGitSuggestContext_WithCurrentFiles(t *testing.T) {
	mock := &MockGitContextBuilder{
		suggestFiles: []string{"related.go"},
	}

	handler := NewGitToolsHandler(nil, mock)
	result, err := handler.Execute("git_suggest_context", map[string]any{
		"task":          "implement feature",
		"current_files": []any{"main.go", "utils.go"},
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "related.go") {
		t.Errorf("expected suggestions, got: %s", result)
	}
}


func TestGitStatus_UntrackedFile(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create untracked file
	os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("content"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "untracked") {
		t.Errorf("expected untracked status, got: %s", result)
	}
}

func TestGitStatus_StagedAndUnstaged(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create and commit file
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	// Modify, stage, then modify again
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("staged"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("unstaged"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should show MM status (staged + unstaged)
	if !strings.Contains(result, "test.txt") {
		t.Errorf("expected file in status, got: %s", result)
	}
}

func TestGitLog_DefaultLimit(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create multiple commits
	for i := 0; i < 15; i++ {
		os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte(fmt.Sprintf("content%d", i)), 0644)
		exec.Command("git", "-C", tmpDir, "add", ".").Run()
		exec.Command("git", "-C", tmpDir, "commit", "-m", fmt.Sprintf("Commit %d", i)).Run()
	}

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_log", map[string]any{}, tmpDir) // No limit specified

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should use default limit of 10
	if !strings.Contains(result, "commits") {
		t.Errorf("expected commits in result, got: %s", result)
	}
}


func TestNewGitToolsHandler(t *testing.T) {
	mock := &MockGitContextBuilder{}
	handler := NewGitToolsHandler(nil, mock)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.GitContext != mock {
		t.Error("expected git context to be set")
	}
}

func TestGitStatus_ShortLine(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create a file to have some status
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("content"), 0644)

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_status", map[string]any{}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty result")
	}
}

func TestGitDiff_NoStagedChanges(t *testing.T) {
	tmpDir := setupGitRepo(t)

	// Create and commit a file
	os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("initial"), 0644)
	exec.Command("git", "-C", tmpDir, "add", ".").Run()
	exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

	handler := NewGitToolsHandler(nil, nil)
	result, err := handler.Execute("git_diff", map[string]any{"staged": true}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No differences") {
		t.Errorf("expected no differences for clean staged, got: %s", result)
	}
}
