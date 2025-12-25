package main

import (
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"time"
)

// ApplyEdits applies edits from Edits JSON
func (a *App) ApplyEdits(edits *domain.EditsJSON) ([]*domain.ApplyResult, error) {
	return a.applyService.ApplyEdits(a.ctx, edits)
}

// ApplySingleEdit applies a single edit
func (a *App) ApplySingleEdit(edit *domain.Edit) (*domain.ApplyResult, error) {
	return a.applyService.ApplySingleEdit(a.ctx, edit)
}

// ValidateEdits validates edits correctness
func (a *App) ValidateEdits(edits *domain.EditsJSON) error {
	return a.applyService.ValidateEdits(a.ctx, edits)
}

// RollbackEdits rolls back edits
func (a *App) RollbackEdits(results []*domain.ApplyResult) error {
	return a.applyService.RollbackEdits(a.ctx, results)
}

// GenerateDiff generates diff between two states
func (a *App) GenerateDiff(beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	return a.diffService.GenerateDiff(a.ctx, beforePath, afterPath, format)
}

// GenerateDiffFromResults generates diff from apply results
func (a *App) GenerateDiffFromResults(results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
	return a.diffService.GenerateDiffFromResults(a.ctx, results, format)
}

// GenerateDiffFromEdits generates diff from Edits JSON
func (a *App) GenerateDiffFromEdits(edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
	return a.diffService.GenerateDiffFromEdits(a.ctx, edits, format)
}

// PublishDiff publishes diff
func (a *App) PublishDiff(diff *domain.DiffResult) error {
	return a.diffService.PublishDiff(a.ctx, diff)
}

// GenerateAndPublishDiff generates and publishes diff
func (a *App) GenerateAndPublishDiff(beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	return a.diffService.GenerateAndPublishDiff(a.ctx, beforePath, afterPath, format)
}

// TestBackend is a simple test for backend functionality
func (a *App) TestBackend(allFilesJson string, rootDir string) (string, error) {
	var allFiles []*domain.FileNode
	if err := json.Unmarshal([]byte(allFilesJson), &allFiles); err != nil {
		return "", fmt.Errorf("failed to parse files JSON: %w", err)
	}

	testResult := map[string]interface{}{
		"status":     "success",
		"filesCount": len(allFiles),
		"rootDir":    rootDir,
		"timestamp":  time.Now().Unix(),
		"message":    "Backend работает корректно",
	}

	if len(allFiles) > 0 {
		testResult["sampleFile"] = allFiles[0].RelPath
	}

	resultJson, err := json.Marshal(testResult)
	if err != nil {
		return "", fmt.Errorf("failed to marshal test result: %w", err)
	}

	return string(resultJson), nil
}
