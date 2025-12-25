package filereader

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// MockLogger for testing
type MockLogger struct {
	warnings []string
	infos    []string
	errors   []string
}

func (m *MockLogger) Info(msg string)    { m.infos = append(m.infos, msg) }
func (m *MockLogger) Warning(msg string) { m.warnings = append(m.warnings, msg) }
func (m *MockLogger) Error(msg string)   { m.errors = append(m.errors, msg) }
func (m *MockLogger) Debug(msg string)   {}
func (m *MockLogger) Fatal(msg string)   { panic(msg) }

func TestReadContents_GoAppProject(t *testing.T) {
	// Get absolute path to test_folder/go-app
	workspaceRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get workspace root: %v", err)
	}

	projectPath := filepath.Join(workspaceRoot, "test_folder", "go-app")

	// Check if test folder exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Skipf("Test folder does not exist: %s", projectPath)
	}

	t.Logf("Testing with project path: %s", projectPath)

	logger := &MockLogger{}
	reader := NewSecureFileReader(logger)

	// Test 1: Read single file with relative path
	t.Run("ReadSingleFile_RelativePath", func(t *testing.T) {
		filePaths := []string{"main.go"}

		contents, err := reader.ReadContents(context.Background(), filePaths, projectPath, nil)
		if err != nil {
			t.Fatalf("ReadContents failed: %v", err)
		}

		t.Logf("Contents keys: %v", getKeys(contents))
		t.Logf("Logger warnings: %v", logger.warnings)
		t.Logf("Logger infos: %v", logger.infos)

		if len(contents) == 0 {
			t.Errorf("Expected content for main.go, got empty map")
		}

		content, exists := contents["main.go"]
		if !exists {
			t.Errorf("Expected key 'main.go' in contents, got keys: %v", getKeys(contents))
		} else {
			t.Logf("main.go content length: %d bytes", len(content))
			if len(content) == 0 {
				t.Errorf("main.go content is empty")
			}
		}
	})

	// Test 2: Read multiple files
	t.Run("ReadMultipleFiles", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		filePaths := []string{"main.go", "go.mod"}

		contents, err := reader.ReadContents(context.Background(), filePaths, projectPath, nil)
		if err != nil {
			t.Fatalf("ReadContents failed: %v", err)
		}

		t.Logf("Contents keys: %v", getKeys(contents))

		if len(contents) != 2 {
			t.Errorf("Expected 2 files, got %d. Keys: %v", len(contents), getKeys(contents))
		}

		for _, path := range filePaths {
			if _, exists := contents[path]; !exists {
				t.Errorf("Missing content for %s", path)
			}
		}
	})

	// Test 3: Read file in subdirectory
	t.Run("ReadFileInSubdirectory", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		filePaths := []string{"services/user_service.go"}

		contents, err := reader.ReadContents(context.Background(), filePaths, projectPath, nil)
		if err != nil {
			t.Fatalf("ReadContents failed: %v", err)
		}

		t.Logf("Contents keys: %v", getKeys(contents))
		t.Logf("Logger warnings: %v", logger.warnings)

		if len(contents) == 0 {
			t.Errorf("Expected content for services/user_service.go, got empty map")
		}

		// Check both possible key formats
		_, exists1 := contents["services/user_service.go"]
		_, exists2 := contents["services\\user_service.go"]

		if !exists1 && !exists2 {
			t.Errorf("Missing content for services/user_service.go. Keys: %v", getKeys(contents))
		}
	})

	// Test 4: Expand directory
	t.Run("ExpandDirectory", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		filePaths := []string{"services"}

		contents, err := reader.ReadContents(context.Background(), filePaths, projectPath, nil)
		if err != nil {
			t.Fatalf("ReadContents failed: %v", err)
		}

		t.Logf("Contents keys after directory expansion: %v", getKeys(contents))
		t.Logf("Logger infos: %v", logger.infos)
		t.Logf("Logger warnings: %v", logger.warnings)

		if len(contents) == 0 {
			t.Errorf("Expected files from services directory, got empty map")
		}
	})

	// Test 5: Mixed - files and directory
	t.Run("MixedFilesAndDirectory", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		filePaths := []string{"main.go", "services"}

		contents, err := reader.ReadContents(context.Background(), filePaths, projectPath, nil)
		if err != nil {
			t.Fatalf("ReadContents failed: %v", err)
		}

		t.Logf("Contents keys: %v", getKeys(contents))
		t.Logf("Number of files: %d", len(contents))

		if len(contents) < 2 {
			t.Errorf("Expected at least 2 files (main.go + services/*), got %d", len(contents))
		}

		if _, exists := contents["main.go"]; !exists {
			t.Errorf("Missing main.go in contents")
		}
	})
}

func getKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
