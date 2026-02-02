package main

import (
	"encoding/json"
	"syntaxia/domain"
)

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

// GetStreamingContext returns context summary metadata for streaming compatibility
func (a *App) GetStreamingContext(contextID string) (string, error) {
	return a.GetContext(contextID)
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

// contextSummaryJSON is a JSON-friendly version of ContextSummary with name field
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

// CloseStreamingContext removes a streaming context and associated resources
func (a *App) CloseStreamingContext(contextID string) error {
	return a.DeleteContext(contextID)
}

// GetContextStats returns statistics about a context
func (a *App) GetContextStats(contextID string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	summary, err := a.contextService.GetContextSummary(a.ctx, contextID)
	if err != nil {
		return "", a.transformError(err)
	}

	stats := map[string]interface{}{
		"id":         summary.ID,
		"fileCount":  summary.FileCount,
		"totalSize":  summary.TotalSize,
		"tokenCount": summary.TokenCount,
		"lineCount":  summary.LineCount,
		"createdAt":  summary.CreatedAt,
	}

	statsJson, err := json.Marshal(stats)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context stats", err)
		return "", a.transformError(marshalErr)
	}

	return string(statsJson), nil
}
