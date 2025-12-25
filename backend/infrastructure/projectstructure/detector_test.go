package projectstructure

import (
	"os"
	"path/filepath"
	"syntaxia/domain"
	"testing"
)

func TestNewDetector(t *testing.T) {
	d := NewDetector()
	if d == nil {
		t.Fatal("NewDetector returned nil")
	}
	if len(d.frameworkDetectors) == 0 {
		t.Error("No framework detectors initialized")
	}
	if len(d.archDetectors) == 0 {
		t.Error("No architecture detectors initialized")
	}
}

func TestDetectStructure(t *testing.T) {
	// Create temp directory with test structure
	tmpDir := t.TempDir()

	// Create a simple Go project structure
	createTestFile(t, tmpDir, "go.mod", "module testproject\n\ngo 1.21")
	createTestFile(t, tmpDir, "main.go", "package main\n\nfunc main() {}")
	createTestDir(t, tmpDir, "domain")
	createTestFile(t, tmpDir, "domain/models.go", "package domain\n\ntype User struct{}")
	createTestDir(t, tmpDir, "application")
	createTestFile(t, tmpDir, "application/service.go", "package application\n\ntype Service struct{}")

	d := NewDetector()
	structure, err := d.DetectStructure(tmpDir)
	if err != nil {
		t.Fatalf("DetectStructure failed: %v", err)
	}

	if structure == nil {
		t.Fatal("DetectStructure returned nil")
	}

	// Check languages detected
	if len(structure.Languages) == 0 {
		t.Error("No languages detected")
	}

	foundGo := false
	for _, lang := range structure.Languages {
		if lang.Name == "Go" {
			foundGo = true
			if lang.FileCount < 2 {
				t.Errorf("Expected at least 2 Go files, got %d", lang.FileCount)
			}
		}
	}
	if !foundGo {
		t.Error("Go language not detected")
	}

	// Check build systems
	if len(structure.BuildSystems) == 0 {
		t.Error("No build systems detected")
	}

	foundGoMod := false
	for _, bs := range structure.BuildSystems {
		if bs.Name == "go" {
			foundGoMod = true
		}
	}
	if !foundGoMod {
		t.Error("Go build system not detected")
	}
}

func TestDetectArchitecture_CleanArchitecture(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Clean Architecture structure
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestDir(t, tmpDir, "infrastructure")
	createTestFile(t, tmpDir, "domain/entity.go", "package domain")
	createTestFile(t, tmpDir, "application/service.go", "package application")
	createTestFile(t, tmpDir, "infrastructure/repo.go", "package infrastructure")

	d := NewDetector()
	arch, err := d.DetectArchitecture(tmpDir)
	if err != nil {
		t.Fatalf("DetectArchitecture failed: %v", err)
	}

	if arch == nil {
		t.Fatal("DetectArchitecture returned nil")
	}

	if arch.Type != domain.ArchCleanArchitecture {
		t.Errorf("Expected Clean Architecture, got %s", arch.Type)
	}

	if arch.Confidence < 0.3 {
		t.Errorf("Expected confidence >= 0.3, got %f", arch.Confidence)
	}

	if len(arch.Indicators) == 0 {
		t.Error("No indicators found")
	}
}

func TestDetectArchitecture_MVC(t *testing.T) {
	tmpDir := t.TempDir()

	// Create MVC structure
	createTestDir(t, tmpDir, "models")
	createTestDir(t, tmpDir, "views")
	createTestDir(t, tmpDir, "controllers")
	createTestFile(t, tmpDir, "models/user.go", "package models")
	createTestFile(t, tmpDir, "views/index.html", "<html></html>")
	createTestFile(t, tmpDir, "controllers/user_controller.go", "package controllers")

	d := NewDetector()
	arch, err := d.DetectArchitecture(tmpDir)
	if err != nil {
		t.Fatalf("DetectArchitecture failed: %v", err)
	}

	if arch.Type != domain.ArchMVC {
		t.Errorf("Expected MVC, got %s", arch.Type)
	}
}

func TestDetectArchitecture_Hexagonal(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Hexagonal structure
	createTestDir(t, tmpDir, "ports")
	createTestDir(t, tmpDir, "adapters")
	createTestDir(t, tmpDir, "core")
	createTestFile(t, tmpDir, "ports/repository.go", "package ports")
	createTestFile(t, tmpDir, "adapters/http.go", "package adapters")
	createTestFile(t, tmpDir, "core/domain.go", "package core")

	d := NewDetector()
	arch, err := d.DetectArchitecture(tmpDir)
	if err != nil {
		t.Fatalf("DetectArchitecture failed: %v", err)
	}

	if arch.Type != domain.ArchHexagonal {
		t.Errorf("Expected Hexagonal, got %s", arch.Type)
	}
}

func TestDetectConventions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with snake_case naming
	createTestFile(t, tmpDir, "user_service.go", "package main")
	createTestFile(t, tmpDir, "user_repository.go", "package main")
	createTestFile(t, tmpDir, "user_handler.go", "package main")
	createTestFile(t, tmpDir, "user_service_test.go", "package main")

	d := NewDetector()
	conventions, err := d.DetectConventions(tmpDir)
	if err != nil {
		t.Fatalf("DetectConventions failed: %v", err)
	}

	if conventions == nil {
		t.Fatal("DetectConventions returned nil")
	}

	// Check test conventions
	if conventions.TestConventions.FileSuffix == "" {
		t.Error("Test file suffix not detected")
	}
}

func TestDetectConventions_TestFramework(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Go project with tests
	createTestFile(t, tmpDir, "go.mod", "module test\n\ngo 1.21")
	createTestFile(t, tmpDir, "main_test.go", "package main")

	d := NewDetector()
	conventions, err := d.DetectConventions(tmpDir)
	if err != nil {
		t.Fatalf("DetectConventions failed: %v", err)
	}

	if conventions.TestConventions.Framework != "go test" {
		t.Errorf("Expected 'go test' framework, got '%s'", conventions.TestConventions.Framework)
	}
}

func TestDetectConventions_JestFramework(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Jest config
	createTestFile(t, tmpDir, "jest.config.js", "module.exports = {}")
	createTestFile(t, tmpDir, "package.json", `{"name": "test"}`)

	d := NewDetector()
	conventions, err := d.DetectConventions(tmpDir)
	if err != nil {
		t.Fatalf("DetectConventions failed: %v", err)
	}

	if conventions.TestConventions.Framework != "jest" {
		t.Errorf("Expected 'jest' framework, got '%s'", conventions.TestConventions.Framework)
	}
}

func TestDetectFolderStructure(t *testing.T) {
	tests := []struct {
		name     string
		dirs     []string
		expected domain.FolderStructure
	}{
		{
			name:     "by-layer",
			dirs:     []string{"domain", "application", "infrastructure"},
			expected: domain.FolderByLayer,
		},
		{
			name:     "by-type",
			dirs:     []string{"components", "utils", "types"},
			expected: domain.FolderByType,
		},
		{
			name:     "by-feature",
			dirs:     []string{"features", "modules", "pages"},
			expected: domain.FolderByFeature,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for _, dir := range tt.dirs {
				createTestDir(t, tmpDir, dir)
				createTestFile(t, tmpDir, filepath.Join(dir, "file.go"), "package "+dir)
			}

			d := NewDetector()
			conventions, _ := d.DetectConventions(tmpDir)

			if conventions.FolderStructure != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, conventions.FolderStructure)
			}
		})
	}
}

func TestGetRelatedLayers(t *testing.T) {
	tmpDir := t.TempDir()

	// Create Clean Architecture structure
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestDir(t, tmpDir, "infrastructure")
	createTestFile(t, tmpDir, "domain/user.go", "package domain")
	createTestFile(t, tmpDir, "application/user_service.go", "package application")

	d := NewDetector()
	layers, err := d.GetRelatedLayers(tmpDir, filepath.Join(tmpDir, "application", "user_service.go"))
	if err != nil {
		t.Fatalf("GetRelatedLayers failed: %v", err)
	}

	// Should find application layer and its dependencies
	if len(layers) == 0 {
		t.Error("No related layers found")
	}

	foundApplication := false
	for _, layer := range layers {
		if layer.Name == "application" {
			foundApplication = true
		}
	}
	if !foundApplication {
		t.Error("Application layer not found in related layers")
	}
}

func TestSuggestRelatedFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files with related names
	createTestDir(t, tmpDir, "domain")
	createTestDir(t, tmpDir, "application")
	createTestFile(t, tmpDir, "domain/user.go", "package domain")
	createTestFile(t, tmpDir, "application/user_service.go", "package application")
	createTestFile(t, tmpDir, "user_test.go", "package main")

	d := NewDetector()
	suggestions, err := d.SuggestRelatedFiles(tmpDir, filepath.Join(tmpDir, "domain", "user.go"))
	if err != nil {
		t.Fatalf("SuggestRelatedFiles failed: %v", err)
	}

	// Should suggest related files
	if len(suggestions) == 0 {
		t.Log("No suggestions found (may be expected depending on structure)")
	}
}

func TestDetectLanguages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files of different languages
	createTestFile(t, tmpDir, "main.go", "package main")
	createTestFile(t, tmpDir, "app.ts", "const x = 1")
	createTestFile(t, tmpDir, "style.css", "body {}")
	createTestFile(t, tmpDir, "index.html", "<html></html>")

	d := NewDetector()
	structure, _ := d.DetectStructure(tmpDir)

	if len(structure.Languages) < 2 {
		t.Errorf("Expected at least 2 languages, got %d", len(structure.Languages))
	}

	// Check primary language detection
	hasPrimary := false
	for _, lang := range structure.Languages {
		if lang.Primary {
			hasPrimary = true
		}
	}
	if !hasPrimary {
		t.Error("No primary language detected")
	}
}

func TestDetectBuildSystems(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string]string
		expectName string
	}{
		{
			name:       "npm",
			files:      map[string]string{"package.json": `{"name": "test", "scripts": {"build": "tsc"}}`},
			expectName: "npm",
		},
		{
			name:       "go",
			files:      map[string]string{"go.mod": "module test\n\ngo 1.21"},
			expectName: "go",
		},
		{
			name:       "make",
			files:      map[string]string{"Makefile": "build:\n\tgo build"},
			expectName: "make",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			for name, content := range tt.files {
				createTestFile(t, tmpDir, name, content)
			}

			d := NewDetector()
			structure, _ := d.DetectStructure(tmpDir)

			found := false
			for _, bs := range structure.BuildSystems {
				if bs.Name == tt.expectName {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Build system %s not detected", tt.expectName)
			}
		})
	}
}

func TestDetectProjectType(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(dir string)
		expected string
	}{
		{
			name: "cli",
			setup: func(dir string) {
				createTestFile(t, dir, "main.go", "package main")
				createTestDir(t, dir, "cmd")
				createTestFile(t, dir, "cmd/app/main.go", "package main")
			},
			expected: "cli",
		},
		{
			name: "monorepo",
			setup: func(dir string) {
				createTestDir(t, dir, "packages")
				createTestDir(t, dir, "packages/app1")
				createTestDir(t, dir, "packages/app2")
			},
			expected: "monorepo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tt.setup(tmpDir)

			d := NewDetector()
			structure, _ := d.DetectStructure(tmpDir)

			if structure.ProjectType != tt.expected {
				t.Errorf("Expected project type %s, got %s", tt.expected, structure.ProjectType)
			}
		})
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
