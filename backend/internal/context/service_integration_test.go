package context

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// IntegrationMockLogger for testing
type IntegrationMockLogger struct {
	warnings []string
	infos    []string
	errors   []string
}

func (m *IntegrationMockLogger) Info(msg string)    { m.infos = append(m.infos, msg) }
func (m *IntegrationMockLogger) Warning(msg string) { m.warnings = append(m.warnings, msg) }
func (m *IntegrationMockLogger) Error(msg string)   { m.errors = append(m.errors, msg) }
func (m *IntegrationMockLogger) Debug(msg string)   {}
func (m *IntegrationMockLogger) Fatal(msg string)   { panic(msg) }

// IntegrationMockTokenCounter for testing
type IntegrationMockTokenCounter struct{}

func (m *IntegrationMockTokenCounter) CountTokens(text string) int {
	// Simple approximation: 1 token per 4 chars
	return len(text) / 4
}

// IntegrationMockEventBus for testing
type IntegrationMockEventBus struct {
	events []string
}

func (m *IntegrationMockEventBus) Emit(event string, data ...interface{}) {
	m.events = append(m.events, event)
}

func (m *IntegrationMockEventBus) On(event string, handler func(data interface{}))  {}
func (m *IntegrationMockEventBus) Off(event string, handler func(data interface{})) {}

// IntegrationMockFileReader that simulates reading files
type IntegrationMockFileReader struct {
	log *IntegrationMockLogger
}

func (m *IntegrationMockFileReader) ReadContents(
	ctx context.Context,
	filePaths []string,
	rootDir string,
	progress func(current, total int64),
) (map[string]string, error) {
	result := make(map[string]string)

	for _, path := range filePaths {
		// Try to read actual file
		fullPath := filepath.Join(rootDir, path)
		data, err := os.ReadFile(fullPath)
		if err != nil {
			m.log.Warning("Cannot read file: " + path + " - " + err.Error())
			continue
		}
		result[path] = string(data)
		m.log.Info("Read file: " + path + " (" + string(rune(len(data))) + " bytes)")
	}

	return result, nil
}

func TestContextService_BuildContext_GoApp(t *testing.T) {
	// Get path to test_folder/go-app
	workspaceRoot, err := filepath.Abs("../../..")
	if err != nil {
		t.Fatalf("Failed to get workspace root: %v", err)
	}

	projectPath := filepath.Join(workspaceRoot, "test_folder", "go-app")

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Skipf("Test folder does not exist: %s", projectPath)
	}

	t.Logf("Testing with project path: %s", projectPath)

	logger := &IntegrationMockLogger{}
	tokenCounter := &IntegrationMockTokenCounter{}
	eventBus := &IntegrationMockEventBus{}
	fileReader := &IntegrationMockFileReader{log: logger}

	service, err := NewService(fileReader, tokenCounter, eventBus, logger)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Shutdown(context.Background())

	t.Run("BuildContext_SingleFile", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		includedPaths := []string{"main.go"}

		ctx, err := service.BuildContext(context.Background(), projectPath, includedPaths, nil)
		if err != nil {
			t.Fatalf("BuildContext failed: %v", err)
		}

		t.Logf("Context ID: %s", ctx.ID)
		t.Logf("Context Files: %v", ctx.Files)
		t.Logf("Context TokenCount: %d", ctx.TokenCount)
		t.Logf("Logger warnings: %v", logger.warnings)
		t.Logf("Logger infos: %v", logger.infos)

		if len(ctx.Files) == 0 {
			t.Errorf("Expected files in context, got none")
		}

		if ctx.TokenCount == 0 {
			t.Errorf("Expected non-zero token count")
		}
	})

	t.Run("BuildContext_MultipleFiles", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		includedPaths := []string{"main.go", "go.mod"}

		ctx, err := service.BuildContext(context.Background(), projectPath, includedPaths, nil)
		if err != nil {
			t.Fatalf("BuildContext failed: %v", err)
		}

		t.Logf("Context Files: %v", ctx.Files)
		t.Logf("Context TokenCount: %d", ctx.TokenCount)

		if len(ctx.Files) != 2 {
			t.Errorf("Expected 2 files, got %d: %v", len(ctx.Files), ctx.Files)
		}
	})

	t.Run("BuildContext_SubdirectoryFile", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		includedPaths := []string{"services/user_service.go"}

		ctx, err := service.BuildContext(context.Background(), projectPath, includedPaths, nil)
		if err != nil {
			t.Fatalf("BuildContext failed: %v", err)
		}

		t.Logf("Context Files: %v", ctx.Files)
		t.Logf("Logger warnings: %v", logger.warnings)

		if len(ctx.Files) == 0 {
			t.Errorf("Expected files in context, got none. Warnings: %v", logger.warnings)
		}
	})
}
