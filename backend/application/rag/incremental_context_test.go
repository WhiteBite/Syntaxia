package rag

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// incCtxMockLogger implements domain.Logger for testing
type incCtxMockLogger struct{}

func (m *incCtxMockLogger) Debug(message string)   {}
func (m *incCtxMockLogger) Info(message string)    {}
func (m *incCtxMockLogger) Warning(message string) {}
func (m *incCtxMockLogger) Error(message string)   {}
func (m *incCtxMockLogger) Fatal(message string)   {}

// incCtxMockFileReader implements domain.FileContentReader for testing
type incCtxMockFileReader struct {
	contents map[string]string
}

func (m *incCtxMockFileReader) ReadContents(
	ctx context.Context,
	filePaths []string,
	rootDir string,
	progress func(current, total int64),
) (map[string]string, error) {
	result := make(map[string]string)
	for _, path := range filePaths {
		if content, ok := m.contents[path]; ok {
			result[path] = content
		}
	}
	return result, nil
}

// incCtxMockRepoMapBuilder implements RepoMapBuilder for testing
type incCtxMockRepoMapBuilder struct {
	repoMap *RepoMap
	err     error
}

func (m *incCtxMockRepoMapBuilder) BuildMap(ctx context.Context, projectRoot string, maxTokens int) (*RepoMap, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.repoMap != nil {
		return m.repoMap, nil
	}
	return &RepoMap{
		ProjectRoot: projectRoot,
		Entries: []RepoMapEntry{
			{FilePath: "main.go", Symbols: []SymbolSignature{{Name: "main", Kind: "function"}}},
		},
		TotalTokens: 100,
		GeneratedAt: time.Now(),
	}, nil
}

func (m *incCtxMockRepoMapBuilder) InvalidateCache(projectRoot string) {}

func TestStartSession(t *testing.T) {
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		&incCtxMockFileReader{contents: make(map[string]string)},
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	session, err := svc.StartSession(context.Background(), "/tmp/test-project")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if session.ID == "" {
		t.Error("expected non-empty session ID")
	}

	if session.ProjectRoot != "/tmp/test-project" {
		t.Errorf("expected project root /tmp/test-project, got %s", session.ProjectRoot)
	}

	if session.RepoMap == "" {
		t.Error("expected non-empty repo map")
	}

	if len(session.LoadedFiles) != 0 {
		t.Errorf("expected 0 loaded files, got %d", len(session.LoadedFiles))
	}

	if session.TotalTokens <= 0 {
		t.Error("expected positive token count from repo map")
	}
}

func TestRequestFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := map[string]string{
		"main.go":  "package main\n\nfunc main() {\n\tprintln(\"hello\")\n}",
		"utils.go": "package main\n\nfunc helper() string {\n\treturn \"help\"\n}",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileReader := &incCtxMockFileReader{contents: testFiles}
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Request files
	update, err := svc.RequestFiles(context.Background(), session.ID, []string{"main.go", "utils.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	if len(update.NewFiles) != 2 {
		t.Errorf("expected 2 new files, got %d", len(update.NewFiles))
	}

	if update.TokensAdded <= 0 {
		t.Error("expected positive tokens added")
	}

	// Request same files again - should not add duplicates
	update2, err := svc.RequestFiles(context.Background(), session.ID, []string{"main.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	if len(update2.NewFiles) != 0 {
		t.Errorf("expected 0 new files (already loaded), got %d", len(update2.NewFiles))
	}
}

func TestRequestFiles_TokenBudget(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-budget-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a large file
	largeContent := make([]byte, 100000) // ~25000 tokens
	for i := range largeContent {
		largeContent[i] = 'a'
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "large.go"), largeContent, 0644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}

	fileReader := &incCtxMockFileReader{contents: map[string]string{
		"large.go": string(largeContent),
	}}

	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Set a small token budget
	if err := svc.SetMaxTokens(session.ID, 1000); err != nil {
		t.Fatalf("failed to set max tokens: %v", err)
	}

	// Request large file - should be truncated or skipped
	update, err := svc.RequestFiles(context.Background(), session.ID, []string{"large.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Either file is truncated or skipped
	if len(update.NewFiles) > 0 {
		// File was truncated
		if update.TotalTokens > 1000 {
			t.Errorf("total tokens %d exceeds budget 1000", update.TotalTokens)
		}
	} else if len(update.Skipped) == 0 {
		t.Error("expected file to be either loaded (truncated) or skipped")
	}
}

func TestGetDiff(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-diff-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{
		"main.go":  "package main",
		"utils.go": "package main",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileReader := &incCtxMockFileReader{contents: testFiles}
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Request first file
	_, err = svc.RequestFiles(context.Background(), session.ID, []string{"main.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Get diff - should show main.go as added
	diff, err := svc.GetDiff(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get diff: %v", err)
	}

	if len(diff.Added) != 1 || diff.Added[0] != "main.go" {
		t.Errorf("expected main.go in added, got %v", diff.Added)
	}

	// Request second file
	_, err = svc.RequestFiles(context.Background(), session.ID, []string{"utils.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Get diff again - should show utils.go as added
	diff2, err := svc.GetDiff(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get diff: %v", err)
	}

	if len(diff2.Added) != 1 || diff2.Added[0] != "utils.go" {
		t.Errorf("expected utils.go in added, got %v", diff2.Added)
	}
}

func TestEndSession(t *testing.T) {
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		&incCtxMockFileReader{contents: make(map[string]string)},
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), "/tmp/test")
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// End session
	err = svc.EndSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to end session: %v", err)
	}

	// Try to get session - should fail
	_, err = svc.GetSession(context.Background(), session.ID)
	if err == nil {
		t.Error("expected error when getting ended session")
	}
}

func TestGetSession(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-get-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{"main.go": "package main"}
	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileReader := &incCtxMockFileReader{contents: testFiles}
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Load a file
	_, err = svc.RequestFiles(context.Background(), session.ID, []string{"main.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Get session
	retrieved, err := svc.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if retrieved.ID != session.ID {
		t.Errorf("expected session ID %s, got %s", session.ID, retrieved.ID)
	}

	if len(retrieved.LoadedFiles) != 1 {
		t.Errorf("expected 1 loaded file, got %d", len(retrieved.LoadedFiles))
	}
}

func TestGetSessionContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-format-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{"main.go": "package main\n\nfunc main() {}"}
	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileReader := &incCtxMockFileReader{contents: testFiles}
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Load a file
	_, err = svc.RequestFiles(context.Background(), session.ID, []string{"main.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Get formatted context
	contextStr, err := svc.GetSessionContext(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get session context: %v", err)
	}

	// Check that context contains expected parts
	if !incCtxContains(contextStr, "main.go") {
		t.Error("expected main.go in context")
	}

	if !incCtxContains(contextStr, "package main") {
		t.Error("expected package main in context")
	}
}

func TestUnloadFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "incremental-context-unload-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFiles := map[string]string{
		"main.go":  "package main",
		"utils.go": "package main",
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write file: %v", err)
		}
	}

	fileReader := &incCtxMockFileReader{contents: testFiles}
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		fileReader,
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	// Start session
	session, err := svc.StartSession(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("failed to start session: %v", err)
	}

	// Load files
	_, err = svc.RequestFiles(context.Background(), session.ID, []string{"main.go", "utils.go"})
	if err != nil {
		t.Fatalf("failed to request files: %v", err)
	}

	// Get session to check loaded files
	retrieved, err := svc.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	initialTokens := retrieved.TotalTokens

	// Unload one file
	err = svc.UnloadFiles(context.Background(), session.ID, []string{"main.go"})
	if err != nil {
		t.Fatalf("failed to unload files: %v", err)
	}

	// Check that file was unloaded
	retrieved, err = svc.GetSession(context.Background(), session.ID)
	if err != nil {
		t.Fatalf("failed to get session: %v", err)
	}

	if len(retrieved.LoadedFiles) != 1 {
		t.Errorf("expected 1 loaded file after unload, got %d", len(retrieved.LoadedFiles))
	}

	if _, exists := retrieved.LoadedFiles["main.go"]; exists {
		t.Error("main.go should have been unloaded")
	}

	if retrieved.TotalTokens >= initialTokens {
		t.Error("token count should have decreased after unload")
	}
}

func TestGetActiveSessions(t *testing.T) {
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		&incCtxMockFileReader{contents: make(map[string]string)},
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	if svc.GetActiveSessions() != 0 {
		t.Error("expected 0 active sessions initially")
	}

	// Start sessions
	session1, _ := svc.StartSession(context.Background(), "/tmp/test1")
	session2, _ := svc.StartSession(context.Background(), "/tmp/test2")

	if svc.GetActiveSessions() != 2 {
		t.Errorf("expected 2 active sessions, got %d", svc.GetActiveSessions())
	}

	// End one session
	_ = svc.EndSession(context.Background(), session1.ID)

	if svc.GetActiveSessions() != 1 {
		t.Errorf("expected 1 active session, got %d", svc.GetActiveSessions())
	}

	// End remaining session
	_ = svc.EndSession(context.Background(), session2.ID)

	if svc.GetActiveSessions() != 0 {
		t.Errorf("expected 0 active sessions, got %d", svc.GetActiveSessions())
	}
}

func TestSessionNotFound(t *testing.T) {
	svc := NewIncrementalContextService(
		&incCtxMockLogger{},
		&incCtxMockFileReader{contents: make(map[string]string)},
		&incCtxMockRepoMapBuilder{},
	)
	defer svc.Stop()

	nonExistentID := "non-existent-session"

	// All operations should fail with session not found
	_, err := svc.RequestFiles(context.Background(), nonExistentID, []string{"file.go"})
	if err == nil {
		t.Error("expected error for non-existent session")
	}

	_, err = svc.GetDiff(context.Background(), nonExistentID)
	if err == nil {
		t.Error("expected error for non-existent session")
	}

	_, err = svc.GetSession(context.Background(), nonExistentID)
	if err == nil {
		t.Error("expected error for non-existent session")
	}

	_, err = svc.GetSessionContext(context.Background(), nonExistentID)
	if err == nil {
		t.Error("expected error for non-existent session")
	}

	err = svc.EndSession(context.Background(), nonExistentID)
	if err == nil {
		t.Error("expected error for non-existent session")
	}
}

func TestFormatSessionContext(t *testing.T) {
	session := &ContextSession{
		ID:          "test-session",
		ProjectRoot: "/tmp/test",
		RepoMap:     "# Repository Map\n\nmain.go:\n  function main\n",
		LoadedFiles: map[string]string{
			"main.go":  "package main\n\nfunc main() {}",
			"utils.go": "package main\n\nfunc helper() {}",
		},
		FileTokens: map[string]int{
			"main.go":  10,
			"utils.go": 10,
		},
		TotalTokens: 120,
	}

	formatted := FormatSessionContext(session)

	// Check that repo map is included
	if !incCtxContains(formatted, "Repository Map") {
		t.Error("expected Repository Map in formatted context")
	}

	// Check that loaded files section is included
	if !incCtxContains(formatted, "Loaded Files") {
		t.Error("expected Loaded Files section")
	}

	// Check that files are included
	if !incCtxContains(formatted, "main.go") {
		t.Error("expected main.go in formatted context")
	}

	if !incCtxContains(formatted, "utils.go") {
		t.Error("expected utils.go in formatted context")
	}
}

func TestFormatContextUpdate(t *testing.T) {
	tests := []struct {
		name     string
		update   *ContextUpdate
		contains []string
	}{
		{
			name: "with new files",
			update: &ContextUpdate{
				NewFiles: map[string]string{
					"main.go": "package main",
				},
				TokensAdded: 3,
				TotalTokens: 103,
			},
			contains: []string{"Newly Loaded Files", "main.go"},
		},
		{
			name: "with skipped files",
			update: &ContextUpdate{
				NewFiles:    map[string]string{},
				Skipped:     []string{"large.go"},
				TokensAdded: 0,
				TotalTokens: 100,
			},
			contains: []string{"No new files loaded"},
		},
		{
			name: "empty update",
			update: &ContextUpdate{
				NewFiles:    map[string]string{},
				TokensAdded: 0,
				TotalTokens: 100,
			},
			contains: []string{"No new files loaded"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatted := FormatContextUpdate(tt.update)

			for _, expected := range tt.contains {
				if !incCtxContains(formatted, expected) {
					t.Errorf("expected %q in formatted output", expected)
				}
			}
		})
	}
}

func TestCopyMap(t *testing.T) {
	original := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}

	copied := copyMap(original)

	// Check that values are copied
	if copied["key1"] != "value1" || copied["key2"] != "value2" {
		t.Error("copied map has incorrect values")
	}

	// Modify original and check that copy is not affected
	original["key1"] = "modified"

	if copied["key1"] != "value1" {
		t.Error("copy was affected by modification to original")
	}
}

func TestCopyIntMap(t *testing.T) {
	original := map[string]int{
		"key1": 100,
		"key2": 200,
	}

	copied := copyIntMap(original)

	// Check that values are copied
	if copied["key1"] != 100 || copied["key2"] != 200 {
		t.Error("copied map has incorrect values")
	}

	// Modify original and check that copy is not affected
	original["key1"] = 999

	if copied["key1"] != 100 {
		t.Error("copy was affected by modification to original")
	}
}


// incCtxContains is a helper function to check if a string contains a substring
func incCtxContains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && incCtxContainsHelper(s, substr))
}

func incCtxContainsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
