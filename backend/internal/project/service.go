package project

import (
	"context"
	"syntaxia/domain"
)

// Service handles project-level operations including file listing, git operations, and context generation
type Service struct {
	log         domain.Logger
	bus         domain.EventBus
	treeBuilder domain.TreeBuilder
	gitRepo     domain.GitRepository
	contextSvc  ContextService // Interface to context service
}

// ContextService interface for context operations
type ContextService interface {
	GenerateContextAsync(ctx context.Context, rootDir string, includedPaths []string)
}

// NewService creates a new project service
func NewService(
	log domain.Logger,
	bus domain.EventBus,
	treeBuilder domain.TreeBuilder,
	gitRepo domain.GitRepository,
	contextSvc ContextService,
) *Service {
	return &Service{
		log:         log,
		bus:         bus,
		treeBuilder: treeBuilder,
		gitRepo:     gitRepo,
		contextSvc:  contextSvc,
	}
}

// ListFiles lists all files in a directory with optional ignore rules
func (s *Service) ListFiles(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*domain.FileNode, error) {
	s.log.Info("Listing files for directory: " + dirPath)
	nodes, err := s.treeBuilder.BuildTree(dirPath, useGitignore, useCustomIgnore)
	if err != nil {
		s.log.Error("Failed to build file tree: " + err.Error())
		return nil, err
	}
	s.log.Info("File tree built successfully")
	return nodes, nil
}

// InvalidateCache clears the file tree cache
func (s *Service) InvalidateCache() {
	s.log.Info("Invalidating file tree cache")
	s.treeBuilder.InvalidateCache()
}

// GetUncommittedFiles returns all uncommitted files in the git repository
func (s *Service) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	s.log.Info("Getting uncommitted files for: " + projectRoot)
	files, err := s.gitRepo.GetUncommittedFiles(projectRoot)
	if err != nil {
		s.log.Error("Failed to get uncommitted files: " + err.Error())
		return nil, err
	}
	return files, nil
}

// GetRichCommitHistory returns detailed commit history with file changes
func (s *Service) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	s.log.Info("Getting commit history for: " + projectRoot)
	commits, err := s.gitRepo.GetRichCommitHistory(projectRoot, branchName, limit)
	if err != nil {
		s.log.Error("Failed to get commit history: " + err.Error())
		return nil, err
	}
	return commits, nil
}

// IsGitAvailable checks if git is available in the system
func (s *Service) IsGitAvailable() bool {
	return s.gitRepo.IsGitAvailable()
}

// GenerateContext starts asynchronous context generation for the project
func (s *Service) GenerateContext(ctx context.Context, rootDir string, includedPaths []string) {
	s.log.Info("Starting context generation")
	// Use the context service for safe context generation
	s.contextSvc.GenerateContextAsync(ctx, rootDir, includedPaths)
}
