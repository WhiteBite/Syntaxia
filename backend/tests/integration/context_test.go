package integration

import (
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"syntaxia/domain"
	"syntaxia/infrastructure/memory"
)

// TestContextMemory_SaveAndRetrieve tests saving and retrieving context from SQLite
func TestContextMemory_SaveAndRetrieve(t *testing.T) {
	// Setup temp directory
	tempDir := t.TempDir()

	// Create context memory
	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Test data
	testFiles := []string{
		"/project/src/main.go",
		"/project/src/utils.go",
		"/project/README.md",
	}

	ctx := &domain.ConversationContext{
		ID:           "test_ctx_001",
		ProjectRoot:  "/project",
		Topic:        "Test Context",
		Summary:      "Testing context save and retrieve",
		Files:        testFiles,
		Symbols:      []string{"main", "utils"},
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		MessageCount: 5,
	}

	// Save context
	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save context: %v", err)
	}

	// Retrieve context
	retrieved, err := cm.GetContext("test_ctx_001")
	if err != nil {
		t.Fatalf("Failed to retrieve context: %v", err)
	}

	// Verify fields
	if retrieved.ID != ctx.ID {
		t.Errorf("ID mismatch: got %s, want %s", retrieved.ID, ctx.ID)
	}
	if retrieved.ProjectRoot != ctx.ProjectRoot {
		t.Errorf("ProjectRoot mismatch: got %s, want %s", retrieved.ProjectRoot, ctx.ProjectRoot)
	}
	if retrieved.Topic != ctx.Topic {
		t.Errorf("Topic mismatch: got %s, want %s", retrieved.Topic, ctx.Topic)
	}
	if retrieved.Summary != ctx.Summary {
		t.Errorf("Summary mismatch: got %s, want %s", retrieved.Summary, ctx.Summary)
	}
	if len(retrieved.Files) != len(testFiles) {
		t.Errorf("Files count mismatch: got %d, want %d", len(retrieved.Files), len(testFiles))
	}
	for i, f := range retrieved.Files {
		if f != testFiles[i] {
			t.Errorf("File[%d] mismatch: got %s, want %s", i, f, testFiles[i])
		}
	}
}

// TestContextMemory_GetRecentContexts tests retrieving recent contexts
func TestContextMemory_GetRecentContexts(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	projectRoot := "/test/project"

	// Save multiple contexts
	for i := 0; i < 5; i++ {
		ctx := &domain.ConversationContext{
			ID:           "ctx_" + string(rune('a'+i)),
			ProjectRoot:  projectRoot,
			Topic:        "Topic " + string(rune('A'+i)),
			Files:        []string{"/file" + string(rune('1'+i)) + ".go"},
			CreatedAt:    time.Now().Add(time.Duration(i) * time.Minute),
			LastAccessed: time.Now().Add(time.Duration(i) * time.Minute),
		}
		if err := cm.SaveContext(ctx); err != nil {
			t.Fatalf("Failed to save context %d: %v", i, err)
		}
	}

	// Get recent contexts
	recent, err := cm.GetRecentContexts(projectRoot, 3)
	if err != nil {
		t.Fatalf("Failed to get recent contexts: %v", err)
	}

	if len(recent) != 3 {
		t.Errorf("Expected 3 recent contexts, got %d", len(recent))
	}

	// Verify order (most recent first)
	if len(recent) >= 2 {
		if recent[0].LastAccessed.Before(recent[1].LastAccessed) {
			t.Error("Contexts not ordered by last_accessed DESC")
		}
	}
}

// TestContextMemory_FindByTopic tests finding contexts by topic
func TestContextMemory_FindByTopic(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	projectRoot := "/search/project"

	// Save contexts with different topics
	contexts := []struct {
		id    string
		topic string
	}{
		{"ctx_1", "Backend API refactoring"},
		{"ctx_2", "Frontend Vue components"},
		{"ctx_3", "Backend database migration"},
		{"ctx_4", "Testing infrastructure"},
	}

	for _, c := range contexts {
		ctx := &domain.ConversationContext{
			ID:           c.id,
			ProjectRoot:  projectRoot,
			Topic:        c.topic,
			Files:        []string{"/file.go"},
			CreatedAt:    time.Now(),
			LastAccessed: time.Now(),
		}
		if err := cm.SaveContext(ctx); err != nil {
			t.Fatalf("Failed to save context: %v", err)
		}
	}

	// Search for "Backend"
	found, err := cm.FindContextByTopic(projectRoot, "Backend")
	if err != nil {
		t.Fatalf("Failed to find contexts: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 contexts with 'Backend', got %d", len(found))
	}
}

// TestContextMemory_UpdateExisting tests updating an existing context
func TestContextMemory_UpdateExisting(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Save initial context
	ctx := &domain.ConversationContext{
		ID:           "update_test",
		ProjectRoot:  "/project",
		Topic:        "Initial Topic",
		Files:        []string{"/file1.go"},
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
		MessageCount: 1,
	}
	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save initial context: %v", err)
	}

	// Update context
	ctx.Topic = "Updated Topic"
	ctx.Files = []string{"/file1.go", "/file2.go", "/file3.go"}
	ctx.MessageCount = 10

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to update context: %v", err)
	}

	// Retrieve and verify
	retrieved, err := cm.GetContext("update_test")
	if err != nil {
		t.Fatalf("Failed to retrieve updated context: %v", err)
	}

	if retrieved.Topic != "Updated Topic" {
		t.Errorf("Topic not updated: got %s", retrieved.Topic)
	}
	if len(retrieved.Files) != 3 {
		t.Errorf("Files not updated: got %d files", len(retrieved.Files))
	}
	if retrieved.MessageCount != 10 {
		t.Errorf("MessageCount not updated: got %d", retrieved.MessageCount)
	}
}

// TestContextMemory_EmptyFiles tests context with empty files array
func TestContextMemory_EmptyFiles(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	ctx := &domain.ConversationContext{
		ID:           "empty_files_test",
		ProjectRoot:  "/project",
		Topic:        "Empty Files Context",
		Files:        []string{}, // Empty
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save context with empty files: %v", err)
	}

	retrieved, err := cm.GetContext("empty_files_test")
	if err != nil {
		t.Fatalf("Failed to retrieve context: %v", err)
	}

	if retrieved.Files == nil {
		t.Error("Files should not be nil")
	}
}

// TestContextMemory_Preferences tests preference storage
func TestContextMemory_Preferences(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Set preferences
	prefs := map[string]string{
		"theme":    "dark",
		"language": "ru",
		"fontSize": "14",
	}

	for k, v := range prefs {
		if err := cm.SetPreference(k, v); err != nil {
			t.Fatalf("Failed to set preference %s: %v", k, err)
		}
	}

	// Get individual preference
	theme, err := cm.GetPreference("theme")
	if err != nil {
		t.Fatalf("Failed to get preference: %v", err)
	}
	if theme != "dark" {
		t.Errorf("Theme mismatch: got %s, want dark", theme)
	}

	// Get all preferences
	allPrefs, err := cm.GetAllPreferences()
	if err != nil {
		t.Fatalf("Failed to get all preferences: %v", err)
	}

	if len(allPrefs) != 3 {
		t.Errorf("Expected 3 preferences, got %d", len(allPrefs))
	}
}

// TestContextMemory_DifferentProjects tests contexts from different projects
func TestContextMemory_DifferentProjects(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Save contexts for different projects
	projects := []string{"/project/alpha", "/project/beta", "/project/gamma"}

	for i, proj := range projects {
		for j := 0; j < 3; j++ {
			ctx := &domain.ConversationContext{
				ID:           proj + "_ctx_" + string(rune('0'+j)),
				ProjectRoot:  proj,
				Topic:        "Topic " + string(rune('A'+j)),
				Files:        []string{"/file.go"},
				CreatedAt:    time.Now(),
				LastAccessed: time.Now().Add(time.Duration(i*3+j) * time.Minute),
			}
			if err := cm.SaveContext(ctx); err != nil {
				t.Fatalf("Failed to save context: %v", err)
			}
		}
	}

	// Get contexts for specific project
	alphaContexts, err := cm.GetRecentContexts("/project/alpha", 10)
	if err != nil {
		t.Fatalf("Failed to get alpha contexts: %v", err)
	}

	if len(alphaContexts) != 3 {
		t.Errorf("Expected 3 contexts for alpha, got %d", len(alphaContexts))
	}

	// Verify all belong to alpha
	for _, ctx := range alphaContexts {
		if ctx.ProjectRoot != "/project/alpha" {
			t.Errorf("Context from wrong project: %s", ctx.ProjectRoot)
		}
	}
}


// TestContextSummaryJSON tests JSON serialization of context summaries
func TestContextSummaryJSON(t *testing.T) {
	type contextSummaryJSON struct {
		ID          string                  `json:"id"`
		Name        string                  `json:"name,omitempty"`
		ProjectPath string                  `json:"projectPath"`
		FileCount   int                     `json:"fileCount"`
		TotalSize   int64                   `json:"totalSize"`
		TokenCount  int                     `json:"tokenCount"`
		LineCount   int                     `json:"lineCount"`
		CreatedAt   string                  `json:"createdAt"`
		Metadata    *domain.ContextMetadata `json:"metadata,omitempty"`
	}

	// Test with files in metadata
	summary := contextSummaryJSON{
		ID:          "test_123",
		Name:        "Test Context",
		ProjectPath: "/project",
		FileCount:   5,
		TotalSize:   1024,
		TokenCount:  500,
		LineCount:   100,
		CreatedAt:   time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Metadata: &domain.ContextMetadata{
			SelectedFiles: []string{
				"/project/main.go",
				"/project/utils.go",
			},
		},
	}

	// Serialize
	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Deserialize
	var parsed contextSummaryJSON
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify
	if parsed.ID != summary.ID {
		t.Errorf("ID mismatch")
	}
	if parsed.Name != summary.Name {
		t.Errorf("Name mismatch")
	}
	if parsed.Metadata == nil {
		t.Fatal("Metadata is nil")
	}
	if len(parsed.Metadata.SelectedFiles) != 2 {
		t.Errorf("SelectedFiles count mismatch: got %d", len(parsed.Metadata.SelectedFiles))
	}
}

// TestContextFilesRestoration tests that files are properly restored from context
func TestContextFilesRestoration(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Original files to save
	originalFiles := []string{
		"src/main.go",
		"src/handlers/api.go",
		"src/models/user.go",
		"config/settings.yaml",
		"README.md",
	}

	// Save context
	ctx := &domain.ConversationContext{
		ID:           "restore_test",
		ProjectRoot:  "/my/project",
		Topic:        "Restoration Test",
		Files:        originalFiles,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	// Retrieve
	retrieved, err := cm.GetContext("restore_test")
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	// Verify all files are restored
	if len(retrieved.Files) != len(originalFiles) {
		t.Fatalf("Files count mismatch: got %d, want %d", len(retrieved.Files), len(originalFiles))
	}

	// Create a map for easier comparison
	originalMap := make(map[string]bool)
	for _, f := range originalFiles {
		originalMap[f] = true
	}

	for _, f := range retrieved.Files {
		if !originalMap[f] {
			t.Errorf("Unexpected file in restored context: %s", f)
		}
	}
}

// TestContextMemory_LargeFileList tests context with many files
func TestContextMemory_LargeFileList(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Generate large file list
	var files []string
	for i := 0; i < 500; i++ {
		files = append(files, filepath.Join("src", "module"+string(rune('a'+i%26)), "file"+string(rune('0'+i%10))+".go"))
	}

	ctx := &domain.ConversationContext{
		ID:           "large_files_test",
		ProjectRoot:  "/large/project",
		Topic:        "Large File List",
		Files:        files,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save large context: %v", err)
	}

	retrieved, err := cm.GetContext("large_files_test")
	if err != nil {
		t.Fatalf("Failed to retrieve large context: %v", err)
	}

	if len(retrieved.Files) != 500 {
		t.Errorf("Files count mismatch: got %d, want 500", len(retrieved.Files))
	}
}

// TestContextMemory_SpecialCharacters tests context with special characters in topic/files
func TestContextMemory_SpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	ctx := &domain.ConversationContext{
		ID:           "special_chars_test",
		ProjectRoot:  "/проект/тест",
		Topic:        "Тема с кириллицей и 'кавычками' и \"двойными\"",
		Summary:      "Summary with <html> & special chars: \n\t",
		Files: []string{
			"/путь/к/файлу.go",
			"/path with spaces/file.go",
			"/path/file'with'quotes.go",
		},
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	retrieved, err := cm.GetContext("special_chars_test")
	if err != nil {
		t.Fatalf("Failed to retrieve: %v", err)
	}

	if retrieved.Topic != ctx.Topic {
		t.Errorf("Topic with special chars not preserved")
	}
	if retrieved.ProjectRoot != ctx.ProjectRoot {
		t.Errorf("Cyrillic project root not preserved")
	}
	if len(retrieved.Files) != 3 {
		t.Errorf("Files with special chars not preserved")
	}
}

// TestContextMemory_ConcurrentAccess tests concurrent read/write
func TestContextMemory_ConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Run concurrent operations
	done := make(chan bool, 10)

	// Writers
	for i := 0; i < 5; i++ {
		go func(id int) {
			ctx := &domain.ConversationContext{
				ID:           "concurrent_" + string(rune('a'+id)),
				ProjectRoot:  "/project",
				Topic:        "Concurrent " + string(rune('A'+id)),
				Files:        []string{"/file.go"},
				CreatedAt:    time.Now(),
				LastAccessed: time.Now(),
			}
			_ = cm.SaveContext(ctx)
			done <- true
		}(i)
	}

	// Readers
	for i := 0; i < 5; i++ {
		go func() {
			_, _ = cm.GetRecentContexts("/project", 10)
			done <- true
		}()
	}

	// Wait for all
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify data integrity
	contexts, err := cm.GetRecentContexts("/project", 10)
	if err != nil {
		t.Fatalf("Failed to get contexts after concurrent access: %v", err)
	}

	if len(contexts) != 5 {
		t.Errorf("Expected 5 contexts, got %d", len(contexts))
	}
}

// TestExtractTopicFromMessage tests topic extraction
func TestExtractTopicFromMessage(t *testing.T) {
	tests := []struct {
		message  string
		expected string
	}{
		{"работа над backend api", "backend api"},
		{"fix the login bug", "the login bug"},
		{"implement user authentication", "user authentication"},
		{"добавь новую функцию", "новую функцию"},
		{"refactor the database layer", "the database layer"},
		{"simple message", "simple message"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := memory.ExtractTopicFromMessage(tt.message)
			if result != tt.expected {
				t.Errorf("ExtractTopicFromMessage(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}


// TestContextMemory_WindowsPathNormalization tests that Windows paths with backslashes are handled correctly
func TestContextMemory_WindowsPathNormalization(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Simulate Windows-style paths (with backslashes)
	windowsFiles := []string{
		"C:\\Users\\dev\\project\\src\\main.go",
		"C:\\Users\\dev\\project\\src\\utils\\helper.go",
		"C:\\Users\\dev\\project\\README.md",
	}

	ctx := &domain.ConversationContext{
		ID:           "windows_path_test",
		ProjectRoot:  "C:\\Users\\dev\\project",
		Topic:        "Windows Path Test",
		Files:        windowsFiles,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	// Save context
	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save context with Windows paths: %v", err)
	}

	// Retrieve context
	retrieved, err := cm.GetContext("windows_path_test")
	if err != nil {
		t.Fatalf("Failed to retrieve context: %v", err)
	}

	// Verify files are preserved exactly as saved
	if len(retrieved.Files) != len(windowsFiles) {
		t.Fatalf("Files count mismatch: got %d, want %d", len(retrieved.Files), len(windowsFiles))
	}

	for i, f := range retrieved.Files {
		if f != windowsFiles[i] {
			t.Errorf("File[%d] mismatch: got %s, want %s", i, f, windowsFiles[i])
		}
	}
}

// TestContextMemory_MixedPathStyles tests contexts with mixed path styles
func TestContextMemory_MixedPathStyles(t *testing.T) {
	tempDir := t.TempDir()

	cm, err := memory.NewContextMemory(tempDir)
	if err != nil {
		t.Fatalf("Failed to create context memory: %v", err)
	}
	defer cm.Close()

	// Mix of Unix and Windows style paths (can happen in cross-platform scenarios)
	mixedFiles := []string{
		"src/main.go",                    // Unix relative
		"src\\utils\\helper.go",          // Windows relative
		"/project/src/api.go",            // Unix absolute
		"C:\\project\\config\\app.yaml",  // Windows absolute
	}

	ctx := &domain.ConversationContext{
		ID:           "mixed_path_test",
		ProjectRoot:  "/project",
		Topic:        "Mixed Path Styles",
		Files:        mixedFiles,
		CreatedAt:    time.Now(),
		LastAccessed: time.Now(),
	}

	if err := cm.SaveContext(ctx); err != nil {
		t.Fatalf("Failed to save context with mixed paths: %v", err)
	}

	retrieved, err := cm.GetContext("mixed_path_test")
	if err != nil {
		t.Fatalf("Failed to retrieve context: %v", err)
	}

	// All paths should be preserved
	if len(retrieved.Files) != len(mixedFiles) {
		t.Fatalf("Files count mismatch: got %d, want %d", len(retrieved.Files), len(mixedFiles))
	}

	for i, f := range retrieved.Files {
		if f != mixedFiles[i] {
			t.Errorf("File[%d] mismatch: got %s, want %s", i, f, mixedFiles[i])
		}
	}
}
