package tools

import (
	"errors"
	"strings"
	"testing"

	"syntaxia/domain"
)

// MockProjectStructureService implements ProjectStructureService for testing
type MockProjectStructureService struct {
	architecture    *domain.ArchitectureInfo
	frameworks      []domain.FrameworkInfo
	conventions     *domain.ConventionInfo
	summary         string
	relatedLayers   []domain.LayerInfo
	suggestedFiles  []string
	err             error
}

func (m *MockProjectStructureService) DetectArchitecture(projectRoot string) (*domain.ArchitectureInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.architecture, nil
}

func (m *MockProjectStructureService) DetectFrameworks(projectRoot string) ([]domain.FrameworkInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.frameworks, nil
}

func (m *MockProjectStructureService) DetectConventions(projectRoot string) (*domain.ConventionInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.conventions, nil
}

func (m *MockProjectStructureService) GetArchitectureSummary(projectRoot string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.summary, nil
}

func (m *MockProjectStructureService) GetRelatedLayers(projectRoot, filePath string) ([]domain.LayerInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.relatedLayers, nil
}

func (m *MockProjectStructureService) SuggestRelatedFiles(projectRoot, filePath string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.suggestedFiles, nil
}

func TestNewProjectStructureToolsHandler(t *testing.T) {
	service := &MockProjectStructureService{}
	handler := NewProjectStructureToolsHandler(nil, service)

	if handler == nil {
		t.Fatal("expected non-nil handler")
	}
	if handler.service != service {
		t.Error("expected service to be set")
	}
}

func TestProjectStructureToolsHandler_CanHandle(t *testing.T) {
	handler := NewProjectStructureToolsHandler(nil, nil)

	tests := []struct {
		toolName string
		expected bool
	}{
		{"detect_architecture", true},
		{"detect_frameworks", true},
		{"detect_conventions", true},
		{"get_project_structure", true},
		{"get_related_layers", true},
		{"suggest_related_files", true},
		{"unknown_tool", false},
	}

	for _, tt := range tests {
		t.Run(tt.toolName, func(t *testing.T) {
			if got := handler.CanHandle(tt.toolName); got != tt.expected {
				t.Errorf("CanHandle(%s) = %v, want %v", tt.toolName, got, tt.expected)
			}
		})
	}
}

func TestProjectStructureToolsHandler_GetTools(t *testing.T) {
	handler := NewProjectStructureToolsHandler(nil, nil)
	tools := handler.GetTools()

	if len(tools) != 6 {
		t.Errorf("expected 6 tools, got: %d", len(tools))
	}
}

func TestDetectArchitecture_Success(t *testing.T) {
	service := &MockProjectStructureService{
		architecture: &domain.ArchitectureInfo{
			Type:        "Clean Architecture",
			Confidence:  0.85,
			Description: "Clean Architecture pattern detected",
			Indicators:  []string{"domain layer", "application layer"},
			Layers: []domain.LayerInfo{
				{Name: "domain", Path: "backend/domain", Description: "Business entities"},
			},
		},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("detect_architecture", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Clean Architecture") {
		t.Errorf("expected architecture type, got: %s", result)
	}
	if !strings.Contains(result, "85%") {
		t.Errorf("expected confidence, got: %s", result)
	}
}

func TestDetectArchitecture_Error(t *testing.T) {
	service := &MockProjectStructureService{
		err: errors.New("detection failed"),
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	_, err := handler.Execute("detect_architecture", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDetectFrameworks_Success(t *testing.T) {
	service := &MockProjectStructureService{
		frameworks: []domain.FrameworkInfo{
			{
				Name:          "Vue.js",
				Version:       "3.0",
				Category:      "Frontend",
				Language:      "TypeScript",
				ConfigFiles:   []string{"vite.config.ts"},
				BestPractices: []string{"Use composition API", "Use TypeScript"},
			},
		},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("detect_frameworks", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Vue.js") {
		t.Errorf("expected framework name, got: %s", result)
	}
}

func TestDetectFrameworks_NoFrameworks(t *testing.T) {
	service := &MockProjectStructureService{
		frameworks: []domain.FrameworkInfo{},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("detect_frameworks", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No frameworks") {
		t.Errorf("expected no frameworks message, got: %s", result)
	}
}

func TestDetectConventions_Success(t *testing.T) {
	service := &MockProjectStructureService{
		conventions: &domain.ConventionInfo{
			NamingStyle:     "camelCase",
			FolderStructure: "by-feature",
			FileNaming: domain.FileNamingStyle{
				Style:    "snake_case",
				Suffixes: []string{"_test.go", "_mock.go"},
			},
			TestConventions: domain.TestConventions{
				Location:   "same directory",
				FileSuffix: "_test.go",
			},
			ImportStyle: domain.ImportStyle{
				AbsoluteImports: true,
				ImportOrder:     []string{"stdlib", "external", "internal"},
			},
		},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("detect_conventions", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "camelCase") {
		t.Errorf("expected naming style, got: %s", result)
	}
}

func TestGetRelatedLayers_Success(t *testing.T) {
	service := &MockProjectStructureService{
		relatedLayers: []domain.LayerInfo{
			{Name: "domain", Path: "backend/domain", Description: "Business logic"},
			{Name: "application", Path: "backend/application", Dependencies: []string{"domain"}},
		},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("get_related_layers", map[string]any{
		"path": "backend/application/service.go",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "domain") {
		t.Errorf("expected layer name, got: %s", result)
	}
}

func TestGetRelatedLayers_MissingPath(t *testing.T) {
	handler := NewProjectStructureToolsHandler(nil, &MockProjectStructureService{})

	_, err := handler.Execute("get_related_layers", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestGetRelatedLayers_NoLayers(t *testing.T) {
	service := &MockProjectStructureService{
		relatedLayers: []domain.LayerInfo{},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("get_related_layers", map[string]any{
		"path": "random.go",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No architectural layers") {
		t.Errorf("expected no layers message, got: %s", result)
	}
}

func TestSuggestRelatedFiles_Success(t *testing.T) {
	service := &MockProjectStructureService{
		suggestedFiles: []string{"service_test.go", "repository.go"},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("suggest_related_files", map[string]any{
		"path": "service.go",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "service_test.go") {
		t.Errorf("expected suggested file, got: %s", result)
	}
}

func TestSuggestRelatedFiles_MissingPath(t *testing.T) {
	handler := NewProjectStructureToolsHandler(nil, &MockProjectStructureService{})

	_, err := handler.Execute("suggest_related_files", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing path")
	}
}

func TestSuggestRelatedFiles_NoSuggestions(t *testing.T) {
	service := &MockProjectStructureService{
		suggestedFiles: []string{},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("suggest_related_files", map[string]any{
		"path": "isolated.go",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "No related files") {
		t.Errorf("expected no suggestions message, got: %s", result)
	}
}

func TestProjectStructureToolsHandler_UnknownTool(t *testing.T) {
	handler := NewProjectStructureToolsHandler(nil, &MockProjectStructureService{})

	_, err := handler.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown project structure tool") {
		t.Errorf("expected 'unknown project structure tool' error, got: %v", err)
	}
}

func TestDetectFrameworks_ManyBestPractices(t *testing.T) {
	service := &MockProjectStructureService{
		frameworks: []domain.FrameworkInfo{
			{
				Name:     "React",
				Category: "Frontend",
				Language: "TypeScript",
				BestPractices: []string{
					"Practice 1",
					"Practice 2",
					"Practice 3",
					"Practice 4",
					"Practice 5",
				},
			},
		},
	}
	handler := NewProjectStructureToolsHandler(nil, service)

	result, err := handler.Execute("detect_frameworks", map[string]any{}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should truncate to 3 and show "and X more"
	if !strings.Contains(result, "and 2 more") {
		t.Errorf("expected truncation message, got: %s", result)
	}
}
