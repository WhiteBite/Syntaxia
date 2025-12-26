package tools

import (
	"fmt"
	"path/filepath"
	"strings"

	"syntaxia/domain"
)

// ProjectStructureService interface for project structure analysis
type ProjectStructureService interface {
	DetectArchitecture(projectRoot string) (*domain.ArchitectureInfo, error)
	DetectFrameworks(projectRoot string) ([]domain.FrameworkInfo, error)
	DetectConventions(projectRoot string) (*domain.ConventionInfo, error)
	GetArchitectureSummary(projectRoot string) (string, error)
	GetRelatedLayers(projectRoot, filePath string) ([]domain.LayerInfo, error)
	SuggestRelatedFiles(projectRoot, filePath string) ([]string, error)
}

// ProjectStructureToolsHandler handles project structure analysis tools
type ProjectStructureToolsHandler struct {
	BaseHandler
	service ProjectStructureService
}

// NewProjectStructureToolsHandler creates a new project structure tools handler
func NewProjectStructureToolsHandler(logger domain.Logger, service ProjectStructureService) *ProjectStructureToolsHandler {
	return &ProjectStructureToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		service:     service,
	}
}

var projectStructureToolNames = map[string]bool{
	"detect_architecture":   true,
	"detect_frameworks":     true,
	"detect_conventions":    true,
	"get_project_structure": true,
	"get_related_layers":    true,
	"suggest_related_files": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *ProjectStructureToolsHandler) CanHandle(toolName string) bool {
	return projectStructureToolNames[toolName]
}

// GetTools returns the list of project structure tools
func (h *ProjectStructureToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "detect_architecture",
			Description: "Detect the architecture pattern of the project (Clean Architecture, Hexagonal, MVC, MVVM, Layered, DDD).",
			Parameters:  domain.ToolParameters{Type: "object", Properties: map[string]domain.ToolProperty{}},
		},
		{
			Name:        "detect_frameworks",
			Description: "Detect frameworks and libraries used in the project (Vue, React, Gin, Express, Spring, etc.).",
			Parameters:  domain.ToolParameters{Type: "object", Properties: map[string]domain.ToolProperty{}},
		},
		{
			Name:        "detect_conventions",
			Description: "Detect coding conventions in the project (naming style, folder structure, test conventions).",
			Parameters:  domain.ToolParameters{Type: "object", Properties: map[string]domain.ToolProperty{}},
		},
		{
			Name:        "get_project_structure",
			Description: "Get complete project structure analysis including architecture, frameworks, conventions.",
			Parameters:  domain.ToolParameters{Type: "object", Properties: map[string]domain.ToolProperty{}},
		},
		{
			Name:        "get_related_layers",
			Description: "Get architectural layers related to a specific file.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the file (relative to project root)"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "suggest_related_files",
			Description: "Suggest related files based on architecture patterns.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the file (relative to project root)"},
				},
				Required: []string{"path"},
			},
		},
	}
}

// Execute executes a project structure tool
func (h *ProjectStructureToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	switch toolName {
	case "detect_architecture":
		return h.detectArchitecture(projectRoot)
	case "detect_frameworks":
		return h.detectFrameworks(projectRoot)
	case "detect_conventions":
		return h.detectConventions(projectRoot)
	case "get_project_structure":
		return h.service.GetArchitectureSummary(projectRoot)
	case "get_related_layers":
		return h.getRelatedLayers(args, projectRoot)
	case "suggest_related_files":
		return h.suggestRelatedFiles(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown project structure tool: %s", toolName)
	}
}

func (h *ProjectStructureToolsHandler) detectArchitecture(projectRoot string) (string, error) {
	arch, err := h.service.DetectArchitecture(projectRoot)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Architecture: %s\n", arch.Type))
	result.WriteString(fmt.Sprintf("Confidence: %.0f%%\n", arch.Confidence*100))
	result.WriteString(fmt.Sprintf("Description: %s\n", arch.Description))

	if len(arch.Indicators) > 0 {
		result.WriteString("\nIndicators:\n")
		for _, ind := range arch.Indicators {
			result.WriteString(fmt.Sprintf("  - %s\n", ind))
		}
	}

	if len(arch.Layers) > 0 {
		result.WriteString("\nArchitectural Layers:\n")
		for _, layer := range arch.Layers {
			result.WriteString(fmt.Sprintf("  [%s] %s\n", layer.Name, layer.Path))
			if layer.Description != "" {
				result.WriteString(fmt.Sprintf("    Description: %s\n", layer.Description))
			}
			if len(layer.Dependencies) > 0 {
				result.WriteString(fmt.Sprintf("    Dependencies: %v\n", layer.Dependencies))
			}
		}
	}

	return result.String(), nil
}

func (h *ProjectStructureToolsHandler) detectFrameworks(projectRoot string) (string, error) {
	frameworks, err := h.service.DetectFrameworks(projectRoot)
	if err != nil {
		return "", err
	}

	if len(frameworks) == 0 {
		return "No frameworks detected in this project.", nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Detected %d frameworks:\n\n", len(frameworks)))

	for _, fw := range frameworks {
		version := ""
		if fw.Version != "" {
			version = " v" + fw.Version
		}
		result.WriteString(fmt.Sprintf("[%s] %s%s\n", fw.Category, fw.Name, version))
		result.WriteString(fmt.Sprintf("  Language: %s\n", fw.Language))

		if len(fw.ConfigFiles) > 0 {
			result.WriteString(fmt.Sprintf("  Config files: %v\n", fw.ConfigFiles))
		}

		if len(fw.BestPractices) > 0 {
			result.WriteString("  Best Practices:\n")
			for i, bp := range fw.BestPractices {
				if i >= 3 {
					result.WriteString(fmt.Sprintf("    ... and %d more\n", len(fw.BestPractices)-3))
					break
				}
				result.WriteString(fmt.Sprintf("    • %s\n", bp))
			}
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

func (h *ProjectStructureToolsHandler) detectConventions(projectRoot string) (string, error) {
	conventions, err := h.service.DetectConventions(projectRoot)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	result.WriteString("Project Conventions:\n\n")
	result.WriteString(fmt.Sprintf("Naming Style: %s\n", conventions.NamingStyle))
	result.WriteString(fmt.Sprintf("Folder Structure: %s\n", conventions.FolderStructure))

	result.WriteString("\nFile Naming:\n")
	result.WriteString(fmt.Sprintf("  Style: %s\n", conventions.FileNaming.Style))
	if len(conventions.FileNaming.Suffixes) > 0 {
		result.WriteString(fmt.Sprintf("  Common suffixes: %v\n", conventions.FileNaming.Suffixes))
	}

	result.WriteString("\nTest Conventions:\n")
	result.WriteString(fmt.Sprintf("  Location: %s\n", conventions.TestConventions.Location))
	result.WriteString(fmt.Sprintf("  File suffix: %s\n", conventions.TestConventions.FileSuffix))

	result.WriteString("\nImport Style:\n")
	result.WriteString(fmt.Sprintf("  Absolute imports: %v\n", conventions.ImportStyle.AbsoluteImports))
	result.WriteString(fmt.Sprintf("  Import order: %v\n", conventions.ImportStyle.ImportOrder))

	return result.String(), nil
}

func (h *ProjectStructureToolsHandler) getRelatedLayers(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	layers, err := h.service.GetRelatedLayers(projectRoot, filepath.Join(projectRoot, path))
	if err != nil {
		return "", err
	}

	if len(layers) == 0 {
		return fmt.Sprintf("No architectural layers found related to %s", path), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Layers related to %s:\n\n", path))

	for _, layer := range layers {
		result.WriteString(fmt.Sprintf("[%s] %s\n", layer.Name, layer.Path))
		if layer.Description != "" {
			result.WriteString(fmt.Sprintf("  Description: %s\n", layer.Description))
		}
		if len(layer.Dependencies) > 0 {
			result.WriteString(fmt.Sprintf("  Dependencies: %v\n", layer.Dependencies))
		}
		result.WriteString("\n")
	}

	return result.String(), nil
}

func (h *ProjectStructureToolsHandler) suggestRelatedFiles(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	suggestions, err := h.service.SuggestRelatedFiles(projectRoot, filepath.Join(projectRoot, path))
	if err != nil {
		return "", err
	}

	if len(suggestions) == 0 {
		return fmt.Sprintf("No related files found for %s", path), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Related files for %s:\n\n", path))
	for _, suggestion := range suggestions {
		result.WriteString(fmt.Sprintf("  - %s\n", suggestion))
	}

	return result.String(), nil
}
