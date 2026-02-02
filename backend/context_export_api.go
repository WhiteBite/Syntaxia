package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syntaxia/domain"
)

// ExportContext exports context with specified settings
func (a *App) ExportContext(settingsJson string) (domain.ExportResult, error) {
	var settings domain.ExportSettings
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		validationErr := domain.NewValidationError("failed to parse export settings", map[string]interface{}{
			"originalError": err.Error(),
			"settingsJson":  settingsJson,
		})
		return domain.ExportResult{}, a.transformError(validationErr)
	}

	result, err := a.exportService.Export(a.ctx, settings)
	if err != nil {
		return domain.ExportResult{}, a.transformError(err)
	}

	return result, nil
}

// CleanupTempFiles cleans up temporary export files
func (a *App) CleanupTempFiles(filePath string) error {
	if filePath == "" {
		return nil
	}

	if !strings.Contains(filePath, "syntaxia-export-") {
		return fmt.Errorf("not a temp export file")
	}

	tempDir := filepath.Dir(filePath)
	return os.RemoveAll(tempDir)
}

// ExportProject exports an entire project
func (a *App) ExportProject(projectPath string, format string, optionsJson string) (string, error) {
	var options map[string]interface{}
	if optionsJson != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			return "", fmt.Errorf("failed to parse options JSON: %w", err)
		}
	}

	exportSettings := domain.ExportSettings{
		ProjectPath: projectPath,
		Format:      format,
		Options:     options,
	}

	result, err := a.exportService.Export(a.ctx, exportSettings)
	if err != nil {
		return "", fmt.Errorf("failed to export project: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export result: %w", err)
	}

	return string(resultJson), nil
}

// GetExportHistory returns export history
func (a *App) GetExportHistory(projectPath string) (string, error) {
	history, err := a.exportService.GetExportHistory(a.ctx, projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to get export history: %w", err)
	}

	historyJson, err := json.Marshal(history)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export history: %w", err)
	}

	return string(historyJson), nil
}
