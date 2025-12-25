package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"syntaxia/domain"
	projectservice "syntaxia/internal/project"
)

// ProjectHandler handles all project-related operations
// Delegates to internal/project.Service for core functionality
type ProjectHandler struct {
	log            domain.Logger
	bus            domain.EventBus
	projectService *projectservice.Service
	fileWatcher    domain.FileSystemWatcher
	fileReader     domain.FileContentReader
	gitRepo        domain.GitRepository
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(
	log domain.Logger,
	bus domain.EventBus,
	projectService *projectservice.Service,
	fileWatcher domain.FileSystemWatcher,
	fileReader domain.FileContentReader,
	gitRepo domain.GitRepository,
) *ProjectHandler {
	return &ProjectHandler{
		log:            log,
		bus:            bus,
		projectService: projectService,
		fileWatcher:    fileWatcher,
		fileReader:     fileReader,
		gitRepo:        gitRepo,
	}
}

// ListFiles delegates to projectService
func (h *ProjectHandler) ListFiles(dirPath string, useGitignore, useCustomIgnore bool) ([]*domain.FileNode, error) {
	return h.projectService.ListFiles(dirPath, useGitignore, useCustomIgnore)
}

// GetCurrentDirectory returns the current working directory
func (h *ProjectHandler) GetCurrentDirectory() (string, error) {
	return os.Getwd()
}

// StartFileWatcher starts watching a directory for changes
func (h *ProjectHandler) StartFileWatcher(rootDirPath string) error {
	return h.fileWatcher.Start(rootDirPath)
}

// StopFileWatcher stops the file watcher
func (h *ProjectHandler) StopFileWatcher() {
	h.fileWatcher.Stop()
}

// ReadFileContent reads the content of a file
func (h *ProjectHandler) ReadFileContent(ctx context.Context, rootDir, relPath string) (string, error) {
	contents, err := h.fileReader.ReadContents(ctx, []string{relPath}, rootDir, nil)
	if err != nil {
		return "", err
	}
	if content, ok := contents[relPath]; ok {
		return content, nil
	}
	return "", fmt.Errorf("file not found: %s", relPath)
}

// GetFileStats returns file statistics
func (h *ProjectHandler) GetFileStats(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file stats: %w", err)
	}

	stats := map[string]interface{}{
		"name":    fileInfo.Name(),
		"size":    fileInfo.Size(),
		"modTime": fileInfo.ModTime().Unix(),
		"isDir":   fileInfo.IsDir(),
		"mode":    fileInfo.Mode().String(),
	}

	statsJSON, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("failed to marshal file stats: %w", err)
	}

	return string(statsJSON), nil
}

// GenerateContext delegates context generation to project service
func (h *ProjectHandler) GenerateContext(ctx context.Context, rootDir string, includedPaths []string) {
	h.projectService.GenerateContext(ctx, rootDir, includedPaths)
}

// === Git Operations - delegate to projectService ===

// IsGitAvailable checks if git is available
func (h *ProjectHandler) IsGitAvailable() bool {
	return h.projectService.IsGitAvailable()
}

// GetUncommittedFiles returns all uncommitted files
func (h *ProjectHandler) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	return h.projectService.GetUncommittedFiles(projectRoot)
}

// GetRichCommitHistory returns detailed commit history
func (h *ProjectHandler) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	return h.projectService.GetRichCommitHistory(projectRoot, branchName, limit)
}

// GetFileContentAtCommit returns file content at a specific commit
func (h *ProjectHandler) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	return h.gitRepo.GetFileContentAtCommit(projectRoot, filePath, commitHash)
}

// GetGitignoreContent returns .gitignore content
func (h *ProjectHandler) GetGitignoreContent(projectRoot string) (string, error) {
	return h.gitRepo.GetGitignoreContent(projectRoot)
}

// GetBranches returns all git branches
func (h *ProjectHandler) GetBranches(projectRoot string) ([]string, error) {
	return h.gitRepo.GetBranches(projectRoot)
}

// GetCurrentBranch returns the current git branch
func (h *ProjectHandler) GetCurrentBranch(projectRoot string) (string, error) {
	return h.gitRepo.GetCurrentBranch(projectRoot)
}

// ClearCache clears the file tree cache
func (h *ProjectHandler) ClearCache() {
	h.projectService.InvalidateCache()
}
