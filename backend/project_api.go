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

// GetFileDependencies returns dependencies for a single file
func (a *App) GetFileDependencies(projectRoot, filePath string) ([]domain.FileDependency, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if err := validateRelativePath(projectRoot, filePath); err != nil {
		return nil, err
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}
	return a.container.DependencyAnalyzer.AnalyzeFile(projectRoot, filePath)
}

// GetProjectDependencyGraph returns the complete dependency graph for a project
func (a *App) GetProjectDependencyGraph(projectRoot string, ignorePatterns []string) (*domain.DependencyGraph, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}
	return a.container.DependencyAnalyzer.BuildGraph(projectRoot, ignorePatterns)
}

// ClearDependencyCache clears the cached dependency graph
func (a *App) ClearDependencyCache() error {
	if a.container.DependencyAnalyzer == nil {
		return fmt.Errorf("dependency analyzer not initialized")
	}
	a.container.DependencyAnalyzer.ClearCache()
	return nil
}

// SearchFiles searches for files in a project using the backend search index
func (a *App) SearchFiles(projectRoot, query string, options domain.SearchOptions) ([]domain.FileSearchResult, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if query == "" {
		return []domain.FileSearchResult{}, nil
	}
	if a.container.FileSearcher == nil {
		return nil, fmt.Errorf("file searcher not initialized")
	}

	// Set default max results if not specified
	if options.MaxResults == 0 {
		options.MaxResults = 100
	}

	return a.container.FileSearcher.SearchFiles(projectRoot, query, options)
}

// FilterFilesByExtension filters files by extension using pre-computed metadata
func (a *App) FilterFilesByExtension(projectRoot string, includeExts, excludeExts []string) ([]*domain.FileNode, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}

	// Get cached tree
	tree, err := a.projectHandler.ListFiles(projectRoot, true, true)
	if err != nil {
		return nil, err
	}

	// Filter tree by extensions
	return filterTreeByExtensions(tree, includeExts, excludeExts), nil
}

// FilterFilesByWeight filters files by token weight using pre-computed TotalSize
func (a *App) FilterFilesByWeight(projectRoot string, minTokens int) ([]*domain.FileNode, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}

	// Get cached tree
	tree, err := a.projectHandler.ListFiles(projectRoot, true, true)
	if err != nil {
		return nil, err
	}

	// Filter tree by weight
	return filterTreeByWeight(tree, minTokens), nil
}

// filterTreeByExtensions filters file tree by extensions
func filterTreeByExtensions(nodes []*domain.FileNode, includeExts, excludeExts []string) []*domain.FileNode {
	if len(includeExts) == 0 && len(excludeExts) == 0 {
		return nodes
	}

	result := make([]*domain.FileNode, 0, len(nodes))

	for _, node := range nodes {
		if node.IsDir {
			// For directories, recursively filter children
			filteredChildren := filterTreeByExtensions(node.Children, includeExts, excludeExts)
			if len(filteredChildren) > 0 {
				nodeCopy := *node
				nodeCopy.Children = filteredChildren
				result = append(result, &nodeCopy)
			}
		} else {
			// For files, check extension
			ext := filepath.Ext(node.Name)
			
			// Check exclude list first
			if len(excludeExts) > 0 && containsString(excludeExts, ext) {
				continue
			}
			
			// Check include list
			if len(includeExts) > 0 {
				if containsString(includeExts, ext) {
					result = append(result, node)
				}
			} else {
				// No include list, include all (except excluded)
				result = append(result, node)
			}
		}
	}

	return result
}

// filterTreeByWeight filters file tree by token weight
func filterTreeByWeight(nodes []*domain.FileNode, minTokens int) []*domain.FileNode {
	if minTokens <= 0 {
		return nodes
	}

	const bytesPerToken = 4 // Approximate: 1 token ≈ 4 bytes

	result := make([]*domain.FileNode, 0, len(nodes))

	for _, node := range nodes {
		if node.IsDir {
			// For directories, check if TotalSize meets threshold
			estimatedTokens := int(node.TotalSize / bytesPerToken)
			if estimatedTokens >= minTokens {
				// Recursively filter children
				filteredChildren := filterTreeByWeight(node.Children, minTokens)
				if len(filteredChildren) > 0 {
					nodeCopy := *node
					nodeCopy.Children = filteredChildren
					result = append(result, &nodeCopy)
				}
			}
		} else {
			// For files, check if Size meets threshold
			estimatedTokens := int(node.Size / bytesPerToken)
			if estimatedTokens >= minTokens {
				result = append(result, node)
			}
		}
	}

	return result
}

// containsString checks if a string slice contains a value
func containsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}


// GetFileDependenciesBatch returns dependencies for multiple files (batch query)
func (a *App) GetFileDependenciesBatch(projectRoot string, filePaths []string) (map[string][]domain.FileDependency, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}

	// Validate all paths
	for _, filePath := range filePaths {
		if err := validateRelativePath(projectRoot, filePath); err != nil {
			return nil, fmt.Errorf("invalid path %s: %w", filePath, err)
		}
	}

	return a.container.DependencyAnalyzer.GetFileDependenciesBatch(projectRoot, filePaths)
}

// GetIncomingDependencies returns files that import this file
func (a *App) GetIncomingDependencies(projectRoot, filePath string) ([]string, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if err := validateRelativePath(projectRoot, filePath); err != nil {
		return nil, err
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}

	return a.container.DependencyAnalyzer.GetIncomingDependencies(projectRoot, filePath)
}

// GetOutgoingDependencies returns files that this file imports
func (a *App) GetOutgoingDependencies(projectRoot, filePath string) ([]string, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if err := validateRelativePath(projectRoot, filePath); err != nil {
		return nil, err
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}

	return a.container.DependencyAnalyzer.GetOutgoingDependencies(projectRoot, filePath)
}

// IsDependencyGraphCached checks if graph is cached
func (a *App) IsDependencyGraphCached(projectRoot string) bool {
	if projectRoot == "" || a.container.DependencyAnalyzer == nil {
		return false
	}
	return a.container.DependencyAnalyzer.IsDependencyGraphCached(projectRoot)
}

// GetDependencyGraphStats returns cache statistics
func (a *App) GetDependencyGraphStats(projectRoot string) (*domain.DependencyGraphStats, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if a.container.DependencyAnalyzer == nil {
		return nil, fmt.Errorf("dependency analyzer not initialized")
	}

	return a.container.DependencyAnalyzer.GetDependencyGraphStats(projectRoot)
}

// UpdateDependencyFile updates single file in cached graph (incremental update)
func (a *App) UpdateDependencyFile(projectRoot, filePath string) error {
	if projectRoot == "" {
		return fmt.Errorf("projectRoot is required")
	}
	if err := validateRelativePath(projectRoot, filePath); err != nil {
		return err
	}
	if a.container.DependencyAnalyzer == nil {
		return fmt.Errorf("dependency analyzer not initialized")
	}

	return a.container.DependencyAnalyzer.UpdateFile(projectRoot, filePath)
}

// RemoveDependencyFile removes file from cached graph
func (a *App) RemoveDependencyFile(projectRoot, filePath string) error {
	if projectRoot == "" {
		return fmt.Errorf("projectRoot is required")
	}
	if err := validateRelativePath(projectRoot, filePath); err != nil {
		return err
	}
	if a.container.DependencyAnalyzer == nil {
		return fmt.Errorf("dependency analyzer not initialized")
	}

	return a.container.DependencyAnalyzer.RemoveFile(projectRoot, filePath)
}
