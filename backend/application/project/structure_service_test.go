package project

import (
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"testing"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string)                                       {}
func (m *mockLogger) Info(msg string)                                        {}
func (m *mockLogger) Warning(msg string)                                     {}
func (m *mockLogger) Error(msg string)                                       {}
func (m *mockLogger) Fatal(msg string)                                       {}
func (m *mockLogger) Debugf(format string, args ...interface{})              {}
func (m *mockLogger) Infof(format string, args ...interface{})               {}
func (m *mockLogger) Warningf(format string, args ...interface{})            {}
func (m *mockLogger) Errorf(format string, args ...interface{})              {}
func (m *mockLogger) Fatalf(format string, args ...interface{})              {}
func (m *mockLogger) WithField(key string, value interface{}) domain.Logger  { return m }
func (m *mockLogger) WithFields(fields map[string]interface{}) domain.Logger { return m }

func TestNewStructureServiceLazy(t *testing.T) {
	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	if service == nil {
		t.Fatal("NewStructureService returned nil")
	}
	if service.detector == nil {
		t.Error("Detector not initialized")
	}
	if service.logger == nil {
		t.Error("Logger not set")
	}
}

func TestStructureService_DetectStructure(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go project
	createTestFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21")
	createTestFile(t, tmpDir, "main.go", "package main\n\nfunc main() {}")
	createTestDir(t, tmpDir, "domain")
	createTestFile(t, tmpDir, "domain/user.go", "package domain")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	structure, err := service.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("DetectStructure failed: %v", err)
	}

	if structure == nil {
		t.Fatal("DetectStructure returned nil")
	}

	// Check that languages are detected
	if len(structure.Languages) == 0 {
		t.Error("No languages detected")
	}

	// Check that build systems are detected
	if len(structure.BuildSystems) == 0 {
		t.Error("No build systems detected")
	}
}

func TestStructureService_DetectArchitecture(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Clean Architecture structure
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestDir(t, tmpDir, "infrastructure")
	createTestFile(t, tmpDir, "domain/entity.go", "package domain")
	createTestFile(t, tmpDir, "application/service.go", "package application")
	createTestFile(t, tmpDir, "infrastructure/repo.go", "package infrastructure")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	arch, err := service.DetectArchitecture(tmpDir)
	if err != nil {
		t.Fatalf("DetectArchitecture failed: %v", err)
	}

	if arch == nil {
		t.Fatal("DetectArchitecture returned nil")
	}

	if arch.Type != domain.ArchCleanArchitecture {
		t.Errorf("Expected Clean Architecture, got %s", arch.Type)
	}
}

func TestStructureService_DetectConventions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with snake_case naming
	createTestFile(t, tmpDir, "user_service.go", "package main")
	createTestFile(t, tmpDir, "user_repository.go", "package main")
	createTestFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	conventions, err := service.DetectConventions(tmpDir)
	if err != nil {
		t.Fatalf("DetectConventions failed: %v", err)
	}

	if conventions == nil {
		t.Fatal("DetectConventions returned nil")
	}

	// Check test conventions
	if conventions.TestConventions.Framework != "go test" {
		t.Errorf("Expected 'go test' framework, got '%s'", conventions.TestConventions.Framework)
	}
}

func TestStructureService_DetectFrameworks(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Vue project
	createTestFile(t, tmpDir, "package.json", `{
		"name": "test",
		"dependencies": {
			"vue": "^3.4.0"
		}
	}`)

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	frameworks, err := service.DetectFrameworks(tmpDir)
	if err != nil {
		t.Fatalf("DetectFrameworks failed: %v", err)
	}

	found := false
	for _, fw := range frameworks {
		if fw.Name == "Vue.js" {
			found = true
		}
	}
	if !found {
		t.Error("Vue.js not detected")
	}
}

func TestStructureService_GetRelatedLayers(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Clean Architecture structure
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestDir(t, tmpDir, "infrastructure")
	createTestFile(t, tmpDir, "domain/user.go", "package domain")
	createTestFile(t, tmpDir, "application/user_service.go", "package application")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	layers, err := service.GetRelatedLayers(tmpDir, filepath.Join(tmpDir, "application", "user_service.go"))
	if err != nil {
		t.Fatalf("GetRelatedLayers failed: %v", err)
	}

	if len(layers) == 0 {
		t.Log("No related layers found (may be expected)")
	}
}

func TestStructureService_SuggestRelatedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with related names
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestFile(t, tmpDir, "domain/user.go", "package domain")
	createTestFile(t, tmpDir, "application/user_service.go", "package application")
	createTestFile(t, tmpDir, "user_test.go", "package main")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	suggestions, err := service.SuggestRelatedFiles(tmpDir, filepath.Join(tmpDir, "domain", "user.go"))
	if err != nil {
		t.Fatalf("SuggestRelatedFiles failed: %v", err)
	}

	// May or may not find suggestions depending on structure
	t.Logf("Found %d suggestions", len(suggestions))
}

func TestStructureService_GetStructureJSON(t *testing.T) {
	tmpDir := t.TempDir()

	createTestFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21")
	createTestFile(t, tmpDir, "main.go", "package main")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	jsonStr, err := service.GetStructureJSON(tmpDir)
	if err != nil {
		t.Fatalf("GetStructureJSON failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("GetStructureJSON returned empty string")
	}

	// Check that it's valid JSON-like structure
	if !strings.Contains(jsonStr, "projectType") {
		t.Error("JSON doesn't contain expected field 'projectType'")
	}
	if !strings.Contains(jsonStr, "languages") {
		t.Error("JSON doesn't contain expected field 'languages'")
	}
}

func TestStructureService_GetArchitectureSummary(t *testing.T) {
	tmpDir := t.TempDir()

	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21")
	createTestFile(t, tmpDir, "domain/entity.go", "package domain")
	createTestFile(t, tmpDir, "application/service.go", "package application")

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	summary, err := service.GetArchitectureSummary(tmpDir)
	if err != nil {
		t.Fatalf("GetArchitectureSummary failed: %v", err)
	}

	if summary == "" {
		t.Error("GetArchitectureSummary returned empty string")
	}

	// Check that summary contains expected sections
	if !strings.Contains(summary, "Project Structure Analysis") {
		t.Error("Summary doesn't contain header")
	}
	if !strings.Contains(summary, "Project Type") {
		t.Error("Summary doesn't contain Project Type")
	}
}

func TestStructureService_EmptyProject(t *testing.T) {
	tmpDir := t.TempDir()

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	structure, err := service.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("DetectStructure failed on empty project: %v", err)
	}

	if structure == nil {
		t.Fatal("DetectStructure returned nil for empty project")
	}

	// Should still return valid structure with defaults
	if structure.ProjectType == "" {
		t.Error("ProjectType should have a default value")
	}
}

func TestStructureService_NonExistentPath(t *testing.T) {
	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	_, err := service.DetectStructure("/non/existent/path/12345")
	// Should not panic, may return error or empty structure
	t.Logf("Error for non-existent path: %v", err)
}

func TestStructureService_ComplexProject(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a complex project structure
	dirs := []string{
		"backend/domain",
		"backend/application",
		"backend/infrastructure",
		"backend/handlers",
		"frontend/src/components",
		"frontend/src/features",
		"frontend/src/stores",
	}

	for _, dir := range dirs {
		createTestDir(t, tmpDir, dir)
	}

	files := map[string]string{
		"backend/go.mod":                      "module backend\n\ngo 1.21",
		"backend/domain/user.go":              "package domain",
		"backend/application/user_service.go": "package application",
		"backend/handlers/user_handler.go":    "package handlers",
		"frontend/package.json":               `{"dependencies": {"vue": "^3.4.0"}}`,
		"frontend/src/App.vue":                "<template></template>",
	}

	for path, content := range files {
		createTestFile(t, tmpDir, path, content)
	}

	logger := &mockLogger{}
	service := NewStructureServiceLazy(logger)

	structure, err := service.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("DetectStructure failed: %v", err)
	}

	// Should detect multiple languages
	if len(structure.Languages) < 2 {
		t.Errorf("Expected at least 2 languages, got %d", len(structure.Languages))
	}

	// Should detect frameworks
	if len(structure.Frameworks) == 0 {
		t.Log("No frameworks detected in complex project")
	}
}

// Helper functions
func createTestDir(t *testing.T, base, path string) {
	t.Helper()
	fullPath := filepath.Join(base, path)
	if err := os.MkdirAll(fullPath, 0o755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", fullPath, err)
	}
}

func createTestFile(t *testing.T, base, path, content string) {
	t.Helper()
	fullPath := filepath.Join(base, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("Failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create file %s: %v", fullPath, err)
	}
}
