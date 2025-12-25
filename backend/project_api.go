package main

import (
	"os"
	"syntaxia/domain"
)

// GetStartupPath returns the path passed via command line (from context menu)
func (a *App) GetStartupPath() (string, error) {
	return a.startupPath, nil
}

// ClearStartupPath clears the startup path after it's been used
func (a *App) ClearStartupPath() error {
	a.startupPath = ""
	return nil
}

// SelectDirectory opens a directory selection dialog
func (a *App) SelectDirectory() (string, error) {
	return a.bridge.OpenDirectoryDialog()
}

// GetCurrentDirectory returns the current working directory
func (a *App) GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// PathExists checks if a path (file or directory) exists
func (a *App) PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ListFiles lists files in a directory with optional gitignore and custom ignore support
func (a *App) ListFiles(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*domain.FileNode, error) {
	return a.projectHandler.ListFiles(dirPath, useGitignore, useCustomIgnore)
}

// ClearFileTreeCache clears the file tree cache (call after changing ignore rules)
func (a *App) ClearFileTreeCache() {
	a.projectHandler.ClearCache()
}

// ReadFileContent reads file content from a project
func (a *App) ReadFileContent(rootDir, relPath string) (string, error) {
	return a.projectHandler.ReadFileContent(a.ctx, rootDir, relPath)
}

// StartFileWatcher starts watching a directory for file changes
func (a *App) StartFileWatcher(rootDirPath string) error {
	return a.projectHandler.StartFileWatcher(rootDirPath)
}

// StopFileWatcher stops the file watcher
func (a *App) StopFileWatcher() {
	a.projectHandler.StopFileWatcher()
}
