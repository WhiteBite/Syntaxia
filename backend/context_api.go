package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"syntaxia/domain"
)

// RequestsyntaxiaContextGeneration generates context for selected files
func (a *App) RequestsyntaxiaContextGeneration(rootDir string, includedPaths []string) {
	a.projectHandler.GenerateContext(a.ctx, rootDir, includedPaths)
}

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

// BuildContext builds context and returns ContextSummary (OOM-safe)
func (a *App) BuildContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	var options domain.ContextBuildOptions
	if strings.TrimSpace(optionsJson) != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			validationErr := domain.NewValidationError("failed to parse options JSON", map[string]interface{}{
				"originalError": err.Error(),
				"optionsJson":   optionsJson,
			})
			return "", a.transformError(validationErr)
		}
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContext called with empty includedPaths - this may include all project files")
	}

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, &options)
	if err != nil {
		return "", a.transformError(err)
	}

	contextJson, err := json.Marshal(summary)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summary", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextJson), nil
}

// GetContextContent returns paginated context content for memory-safe viewing
func (a *App) GetContextContent(contextID string, startLine int, lineCount int) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if lineCount <= 0 {
		lineCount = 1000
	}

	chunk, err := a.contextService.ReadContextChunk(a.ctx, contextID, startLine, lineCount)
	if err != nil {
		return "", a.transformError(err)
	}

	chunkJson, err := json.Marshal(chunk)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context chunk", err)
		return "", a.transformError(marshalErr)
	}

	return string(chunkJson), nil
}

// GetFullContextContent returns the full context content as a string
func (a *App) GetFullContextContent(contextID string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	content, err := a.contextService.ReadContextContent(a.ctx, contextID)
	if err != nil {
		return "", a.transformError(err)
	}

	return content, nil
}

// BuildContextLegacy is deprecated - use BuildContext instead
func (a *App) BuildContextLegacy() (string, error) {
	return "", a.transformError(domain.NewConfigurationError("legacy context building is no longer supported", nil))
}

// GetContext retrieves context metadata by ID
func (a *App) GetContext(contextID string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	summary, err := a.contextService.GetContextSummary(a.ctx, contextID)
	if err != nil {
		return "", a.transformError(err)
	}

	contextJson, err := json.Marshal(summary)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summary", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextJson), nil
}

// contextSummaryJSON is a JSON-friendly version of ContextSummary with name field
type contextSummaryJSON struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	ProjectPath string                 `json:"projectPath"`
	FileCount   int                    `json:"fileCount"`
	TotalSize   int64                  `json:"totalSize"`
	TokenCount  int                    `json:"tokenCount"`
	LineCount   int                    `json:"lineCount"`
	CreatedAt   string                 `json:"createdAt"`
	Metadata    *domain.ContextMetadata `json:"metadata,omitempty"`
}

// GetProjectContexts lists all stored context summaries for a project path
// Combines contexts from both JSON files (built contexts) and SQLite (saved contexts)
func (a *App) GetProjectContexts(projectPath string) (string, error) {
	var allSummaries []contextSummaryJSON

	// 1. Get contexts from JSON files (built contexts)
	if a.contextService != nil {
		summaries, err := a.contextService.GetProjectContextSummaries(a.ctx, projectPath)
		if err == nil && summaries != nil {
			for _, s := range summaries {
				allSummaries = append(allSummaries, contextSummaryJSON{
					ID:          s.ID,
					ProjectPath: s.ProjectPath,
					FileCount:   s.FileCount,
					TotalSize:   s.TotalSize,
					TokenCount:  s.TokenCount,
					LineCount:   s.LineCount,
					CreatedAt:   s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Metadata:    &s.Metadata,
				})
			}
		}
	}

	// 2. Get contexts from SQLite (saved contexts via SaveContextMemory)
	if a.analysisContainer != nil {
		contextMemory := a.analysisContainer.GetContextMemory()
		if contextMemory != nil {
			savedContexts, err := contextMemory.GetRecentContexts(projectPath, 100)
			if err == nil && savedContexts != nil {
				// Build set of existing IDs to avoid duplicates
				existingIDs := make(map[string]bool)
				for _, s := range allSummaries {
					existingIDs[s.ID] = true
				}

				for _, ctx := range savedContexts {
					// Skip if already exists in JSON summaries
					if existingIDs[ctx.ID] {
						continue
					}

					allSummaries = append(allSummaries, contextSummaryJSON{
						ID:          ctx.ID,
						Name:        ctx.Topic, // Topic becomes Name for frontend
						ProjectPath: ctx.ProjectRoot,
						FileCount:   len(ctx.Files),
						TotalSize:   0,
						LineCount:   0,
						CreatedAt:   ctx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
						Metadata: &domain.ContextMetadata{
							SelectedFiles: ctx.Files,
						},
					})
				}
			}
		}
	}

	contextsJson, err := json.Marshal(allSummaries)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summaries", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextsJson), nil
}

// DeleteContext removes context metadata and associated content from disk
func (a *App) DeleteContext(contextID string) error {
	if a.contextService == nil {
		return a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if err := a.contextService.DeleteContext(a.ctx, contextID); err != nil {
		return a.transformError(err)
	}

	return nil
}

// BuildContextFromRequest builds context using provided options and returns ContextSummary
func (a *App) BuildContextFromRequest(projectPath string, includedPaths []string, options *domain.ContextBuildOptions) (*domain.ContextSummary, error) {
	if a.contextService == nil {
		return nil, a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContextFromRequest called with empty includedPaths")
	}

	if options == nil {
		options = &domain.ContextBuildOptions{}
	}

	a.log.Info(fmt.Sprintf("[BuildContextFromRequest] OutputFormat received: '%s', StripComments: %v, MaxTokens: %d",
		options.OutputFormat, options.StripComments, options.MaxTokens))

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, options)
	if err != nil {
		return nil, a.transformError(err)
	}

	return summary, nil
}

// GetContextLines returns a chunk of context content between startLine and endLine inclusive
func (a *App) GetContextLines(contextID string, startLine, endLine int64) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if endLine < startLine {
		return "", a.transformError(domain.NewValidationError("invalid line range", map[string]interface{}{
			"startLine": startLine,
			"endLine":   endLine,
		}))
	}

	lineCount := int(endLine-startLine) + 1
	chunk, err := a.contextService.ReadContextChunk(a.ctx, contextID, int(startLine), lineCount)
	if err != nil {
		return "", a.transformError(err)
	}

	chunkJson, err := json.Marshal(chunk)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context chunk", err)
		return "", a.transformError(marshalErr)
	}

	return string(chunkJson), nil
}

// CreateStreamingContext delegates to BuildContext to create a disk-backed context summary
func (a *App) CreateStreamingContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	return a.BuildContext(projectPath, includedPaths, optionsJson)
}

// GetStreamingContext returns context summary metadata for streaming compatibility
func (a *App) GetStreamingContext(contextID string) (string, error) {
	return a.GetContext(contextID)
}

// CloseStreamingContext removes a streaming context and associated resources
func (a *App) CloseStreamingContext(contextID string) error {
	return a.DeleteContext(contextID)
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
