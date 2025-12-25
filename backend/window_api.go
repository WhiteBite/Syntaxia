package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/infrastructure/wailsbridge"
)

// --- Window State Management ---

// GetWindowState returns current window state (position, size, maximized, fullscreen)
func (a *App) GetWindowState() map[string]interface{} {
	state := a.bridge.GetWindowState()
	return map[string]interface{}{
		"x":          state.X,
		"y":          state.Y,
		"width":      state.Width,
		"height":     state.Height,
		"maximized":  state.Maximized,
		"fullscreen": state.Fullscreen,
	}
}

// SaveWindowState saves current window state to config file
func (a *App) SaveWindowState() error {
	state := a.bridge.GetWindowState()

	// Skip saving if window is already closed (size is 0)
	if state.Width == 0 || state.Height == 0 {
		return nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".syntaxia")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, "window-state.json"), data, 0o644)
}

// LoadWindowState loads and applies saved window state
func (a *App) LoadWindowState() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filepath.Join(homeDir, ".syntaxia", "window-state.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No saved state, use defaults
		}
		return err
	}

	var state wailsbridge.WindowState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	a.bridge.SetWindowState(state)
	return nil
}

// ResetWindowState resets window to default size and position
func (a *App) ResetWindowState() {
	state := wailsbridge.WindowState{
		Width:      1600,
		Height:     900,
		X:          100,
		Y:          100,
		Maximized:  false,
		Fullscreen: false,
	}
	a.bridge.SetWindowState(state)
}

// WindowMaximize maximizes the window
func (a *App) WindowMaximize() {
	a.bridge.WindowMaximize()
}

// WindowToggleMaximize toggles maximized state
func (a *App) WindowToggleMaximize() {
	a.bridge.WindowToggleMaximize()
}
