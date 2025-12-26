package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"syntaxia/domain"
)

// validateRelativePath validates that a relative path doesn't escape the project root
func validateRelativePath(projectRoot, relPath string) error {
	if relPath == "" {
		return fmt.Errorf("path is required")
	}

	// Check for path traversal attempts
	if strings.Contains(relPath, "..") {
		fullPath := filepath.Join(projectRoot, relPath)
		absProjectRoot, err := filepath.Abs(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to resolve project root: %w", err)
		}
		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
		absProjectRoot = filepath.Clean(absProjectRoot)
		absFullPath = filepath.Clean(absFullPath)

		if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
			return fmt.Errorf("path traversal not allowed")
		}
	}
	return nil
}

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
	if dirPath == "" {
		return nil, fmt.Errorf("dirPath is required")
	}
	// Verify the directory exists
	info, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory does not exist: %s", dirPath)
		}
		return nil, fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("path is not a directory: %s", dirPath)
	}
	return a.projectHandler.ListFiles(dirPath, useGitignore, useCustomIgnore)
}

// ClearFileTreeCache clears the file tree cache (call after changing ignore rules)
func (a *App) ClearFileTreeCache() {
	a.projectHandler.ClearCache()
}

// ReadFileContent reads file content from a project
func (a *App) ReadFileContent(rootDir, relPath string) (string, error) {
	if rootDir == "" {
		return "", fmt.Errorf("rootDir is required")
	}
	if err := validateRelativePath(rootDir, relPath); err != nil {
		return "", err
	}
	return a.projectHandler.ReadFileContent(a.ctx, rootDir, relPath)
}

// StartFileWatcher starts watching a directory for file changes
func (a *App) StartFileWatcher(rootDirPath string) error {
	// Update sandbox project root when project is opened
	if a.container.SandboxFS != nil {
		a.container.SandboxFS.SetProjectRoot(rootDirPath)
	}
	return a.projectHandler.StartFileWatcher(rootDirPath)
}

// StopFileWatcher stops the file watcher
func (a *App) StopFileWatcher() {
	a.projectHandler.StopFileWatcher()
}
