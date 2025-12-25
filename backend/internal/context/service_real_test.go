package context

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"syntaxia/infrastructure/filereader"
)

func TestContextService_RealFileReader_GoApp(t *testing.T) {
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

	// Use REAL file reader
	realFileReader := filereader.NewSecureFileReader(logger)

	service, err := NewService(realFileReader, tokenCounter, eventBus, logger)
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Shutdown(context.Background())

	t.Run("RealFileReader_SingleFile", func(t *testing.T) {
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

		if len(ctx.Files) == 0 {
			t.Errorf("Expected files in context, got none")
		}

		if ctx.TokenCount == 0 {
			t.Errorf("Expected non-zero token count")
		}
	})

	t.Run("RealFileReader_MultipleFiles", func(t *testing.T) {
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

	t.Run("RealFileReader_SubdirectoryFile", func(t *testing.T) {
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

	t.Run("RealFileReader_AllFiles", func(t *testing.T) {
		logger.warnings = nil
		logger.infos = nil

		// This is what the frontend sends - all files
		includedPaths := []string{"main.go", "go.mod", "services/user_service.go"}

		ctx, err := service.BuildContext(context.Background(), projectPath, includedPaths, nil)
		if err != nil {
			t.Fatalf("BuildContext failed: %v", err)
		}

		t.Logf("Context Files: %v", ctx.Files)
		t.Logf("Context TokenCount: %d", ctx.TokenCount)
		t.Logf("Logger warnings: %v", logger.warnings)
		t.Logf("Logger infos: %v", logger.infos)

		if len(ctx.Files) != 3 {
			t.Errorf("Expected 3 files, got %d: %v", len(ctx.Files), ctx.Files)
		}

		if ctx.TokenCount == 0 {
			t.Errorf("Expected non-zero token count")
		}
	})
}
