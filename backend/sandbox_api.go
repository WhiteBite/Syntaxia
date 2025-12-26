package main

import (
	"encoding/json"
	"fmt"
)

// SandboxChangeDTO represents a file change for frontend
type SandboxChangeDTO struct {
	Path      string `json:"path"`
	Operation string `json:"operation"`
	Diff      string `json:"diff"`
}

// SandboxChangesResponse contains all pending changes
type SandboxChangesResponse struct {
	Changes    []SandboxChangeDTO `json:"changes"`
	TotalCount int                `json:"totalCount"`
}

// GetSandboxChanges returns all pending changes in the sandbox
func (a *App) GetSandboxChanges() (string, error) {
	if a.container.SandboxFS == nil {
		return "{\"changes\":[], \"totalCount\":0}", nil
	}

	changes := a.container.SandboxFS.GetChanges()
	dtos := make([]SandboxChangeDTO, 0, len(changes))

	for _, change := range changes {
		op := "modify"
		switch change.Operation {
		case 0: // SandboxOpCreate
			op = "create"
		case 1: // SandboxOpModify
			op = "modify"
		case 2: // SandboxOpDelete
			op = "delete"
		}

		diff, _ := a.container.SandboxFS.GetDiff(change.Path)

		dtos = append(dtos, SandboxChangeDTO{
			Path:      change.Path,
			Operation: op,
			Diff:      diff,
		})
	}

	response := SandboxChangesResponse{
		Changes:    dtos,
		TotalCount: len(dtos),
	}

	result, err := json.Marshal(response)
	if err != nil {
		return "", fmt.Errorf("failed to marshal sandbox changes: %w", err)
	}

	return string(result), nil
}

// GetSandboxDiff returns the diff for a specific file
func (a *App) GetSandboxDiff(path string) (string, error) {
	if a.container.SandboxFS == nil {
		return "", nil
	}

	return a.container.SandboxFS.GetDiff(path)
}

// GetSandboxAllDiffs returns unified diff for all changes
func (a *App) GetSandboxAllDiffs() (string, error) {
	if a.container.SandboxFS == nil {
		return "", nil
	}

	return a.container.SandboxFS.GetAllDiffs()
}

// ApplySandboxChanges applies all pending changes to real filesystem
func (a *App) ApplySandboxChanges() error {
	if a.container.SandboxFS == nil {
		return fmt.Errorf("sandbox not initialized")
	}

	return a.container.SandboxFS.Apply()
}

// DiscardSandboxChanges discards all pending changes
func (a *App) DiscardSandboxChanges() error {
	if a.container.SandboxFS == nil {
		return nil
	}

	a.container.SandboxFS.Discard()
	return nil
}

// DiscardSandboxFile discards changes for a specific file
func (a *App) DiscardSandboxFile(path string) error {
	if a.container.SandboxFS == nil {
		return nil
	}

	a.container.SandboxFS.DiscardFile(path)
	return nil
}

// HasSandboxChanges returns true if there are pending changes
func (a *App) HasSandboxChanges() bool {
	if a.container.SandboxFS == nil {
		return false
	}

	return a.container.SandboxFS.HasChanges()
}

// GetSandboxChangeCount returns the number of pending changes
func (a *App) GetSandboxChangeCount() int {
	if a.container.SandboxFS == nil {
		return 0
	}

	return a.container.SandboxFS.GetChangeCount()
}
