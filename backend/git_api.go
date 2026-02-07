package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"syntaxia/domain"
	"syntaxia/infrastructure/git"
)

// Git input validation constants
const (
	minCommitHashLen = 7
	maxCommitHashLen = 40
)

// Regex patterns for git validation
var (
	commitHashRegex = regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`)
	branchNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._/-]+$`)
)

// validateProjectRoot validates that projectRoot is not empty and exists
func validateProjectRoot(projectRoot string) error {
	if projectRoot == "" {
		return fmt.Errorf("projectRoot is required")
	}
	if _, err := os.Stat(projectRoot); os.IsNotExist(err) {
		return fmt.Errorf("projectRoot does not exist: %s", projectRoot)
	}
	return nil
}

// validateGitInput validates projectRoot and optional file paths
func validateGitInput(projectRoot string, paths ...string) error {
	if err := validateProjectRoot(projectRoot); err != nil {
		return err
	}

	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return fmt.Errorf("failed to resolve project root: %w", err)
	}
	absProjectRoot = filepath.Clean(absProjectRoot)

	for _, path := range paths {
		if path == "" {
			continue
		}
		// Check for path traversal
		if strings.Contains(path, "..") {
			fullPath := filepath.Join(projectRoot, path)
			absFullPath, err := filepath.Abs(fullPath)
			if err != nil {
				return fmt.Errorf("failed to resolve path %s: %w", path, err)
			}
			absFullPath = filepath.Clean(absFullPath)

			if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
				return fmt.Errorf("path traversal not allowed: %s", path)
			}
		}
	}
	return nil
}

// validateBranchName validates git branch name format
func validateBranchName(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name is required")
	}
	// Git branch names cannot contain: space, ~, ^, :, ?, *, [, \, control chars
	// They also cannot start with - or end with .lock
	if strings.HasPrefix(branch, "-") {
		return fmt.Errorf("branch name cannot start with '-': %s", branch)
	}
	if strings.HasSuffix(branch, ".lock") {
		return fmt.Errorf("branch name cannot end with '.lock': %s", branch)
	}
	if strings.Contains(branch, "..") {
		return fmt.Errorf("branch name cannot contain '..': %s", branch)
	}
	if !branchNameRegex.MatchString(branch) {
		return fmt.Errorf("invalid branch name format: %s", branch)
	}
	return nil
}

// validateCommitHash validates git commit hash format
func validateCommitHash(hash string) error {
	if hash == "" {
		return fmt.Errorf("commit hash is required")
	}
	if len(hash) < minCommitHashLen || len(hash) > maxCommitHashLen {
		return fmt.Errorf("commit hash must be %d-%d hex characters: %s", minCommitHashLen, maxCommitHashLen, hash)
	}
	if !commitHashRegex.MatchString(hash) {
		return fmt.Errorf("invalid commit hash format (must be hex): %s", hash)
	}
	return nil
}

// validateGitRef validates a git reference (branch name or commit hash)
func validateGitRef(ref string) error {
	if ref == "" {
		return fmt.Errorf("git ref is required")
	}
	// Try as commit hash first
	if commitHashRegex.MatchString(ref) {
		return nil
	}
	// Try as branch name
	return validateBranchName(ref)
}

// IsGitAvailable checks if git is available on the system
func (a *App) IsGitAvailable() bool {
	return a.projectHandler.IsGitAvailable()
}

// IsGitRepository checks if the given path is a git repository
func (a *App) IsGitRepository(projectPath string) bool {
	if projectPath == "" {
		return false
	}
	return a.gitRepo.IsGitRepository(projectPath)
}

// GetUncommittedFiles returns list of uncommitted files in a git repository
func (a *App) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	if err := validateProjectRoot(projectRoot); err != nil {
		return nil, err
	}
	return a.projectHandler.GetUncommittedFiles(projectRoot)
}

// GetRichCommitHistory returns commit history with file changes
func (a *App) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	if err := validateProjectRoot(projectRoot); err != nil {
		return nil, err
	}
	if branchName != "" {
		if err := validateBranchName(branchName); err != nil {
			return nil, err
		}
	}
	if limit <= 0 {
		limit = 50 // default limit
	}
	return a.projectHandler.GetRichCommitHistory(projectRoot, branchName, limit)
}

// GetFileContentAtCommit returns file content at a specific commit
func (a *App) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	if err := validateGitInput(projectRoot, filePath); err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("filePath is required")
	}
	if err := validateCommitHash(commitHash); err != nil {
		return "", err
	}
	return a.projectHandler.GetFileContentAtCommit(projectRoot, filePath, commitHash)
}

// GetGitignoreContent returns the content of .gitignore file
func (a *App) GetGitignoreContent(projectRoot string) (string, error) {
	if err := validateProjectRoot(projectRoot); err != nil {
		return "", err
	}
	return a.projectHandler.GetGitignoreContent(projectRoot)
}

// GetBranches returns all git branches
func (a *App) GetBranches(projectRoot string) (string, error) {
	if err := validateProjectRoot(projectRoot); err != nil {
		return "", err
	}
	branches, err := a.gitRepo.GetBranches(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %w", err)
	}

	branchesJson, err := json.Marshal(branches)
	if err != nil {
		return "", fmt.Errorf("failed to marshal branches: %w", err)
	}

	return string(branchesJson), nil
}

// GetCurrentBranch returns the current git branch
func (a *App) GetCurrentBranch(projectRoot string) (string, error) {
	if err := validateProjectRoot(projectRoot); err != nil {
		return "", err
	}
	branch, err := a.gitRepo.GetCurrentBranch(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return branch, nil
}

// CloneRepository clones a remote git repository
func (a *App) CloneRepository(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("repository URL is required")
	}
	tempDir, err := os.MkdirTemp("", "Syntaxia-git-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	if err := a.gitRepo.CloneRepository(url, tempDir, 1); err != nil {
		os.RemoveAll(tempDir)
		return "", err
	}

	return tempDir, nil
}

// CheckoutBranch switches to a specific branch in a git repository
func (a *App) CheckoutBranch(projectPath, branch string) error {
	if err := validateProjectRoot(projectPath); err != nil {
		return err
	}
	if err := validateBranchName(branch); err != nil {
		return err
	}
	return a.gitRepo.CheckoutBranch(projectPath, branch)
}

// CheckoutCommit switches to a specific commit in a git repository
func (a *App) CheckoutCommit(projectPath, commitHash string) error {
	if err := validateProjectRoot(projectPath); err != nil {
		return err
	}
	if err := validateCommitHash(commitHash); err != nil {
		return err
	}
	return a.gitRepo.CheckoutCommit(projectPath, commitHash)
}

// GetCommitHistory returns recent commits for selection
func (a *App) GetCommitHistory(projectPath string, limit int) (string, error) {
	if err := validateProjectRoot(projectPath); err != nil {
		return "", err
	}
	if limit <= 0 {
		limit = 50 // default limit
	}
	commits, err := a.gitRepo.GetCommitHistory(projectPath, limit)
	if err != nil {
		return "", err
	}

	commitsJson, err := json.Marshal(commits)
	if err != nil {
		return "", fmt.Errorf("failed to marshal commits: %w", err)
	}

	return string(commitsJson), nil
}

// GetRemoteBranches returns all remote branches
func (a *App) GetRemoteBranches(projectPath string) (string, error) {
	if err := validateProjectRoot(projectPath); err != nil {
		return "", err
	}
	branches, err := a.gitRepo.FetchRemoteBranches(projectPath)
	if err != nil {
		return "", err
	}

	branchesJson, err := json.Marshal(branches)
	if err != nil {
		return "", fmt.Errorf("failed to marshal branches: %w", err)
	}

	return string(branchesJson), nil
}

// CleanupTempRepository removes a temporary cloned repository
func (a *App) CleanupTempRepository(path string) error {
	tempDir := os.TempDir()
	if !strings.HasPrefix(path, tempDir) && !strings.Contains(path, "Syntaxia-git-") {
		return fmt.Errorf("refusing to remove non-temp path: %s", path)
	}
	return os.RemoveAll(path)
}

// ListFilesAtRef returns list of files at a specific branch/commit without checkout
func (a *App) ListFilesAtRef(projectPath, ref string) (string, error) {
	if err := validateProjectRoot(projectPath); err != nil {
		return "", err
	}
	if err := validateGitRef(ref); err != nil {
		return "", err
	}
	files, err := a.gitRepo.ListFilesAtRef(projectPath, ref)
	if err != nil {
		return "", err
	}
	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// GetFileAtRef returns file content at a specific branch/commit without checkout
func (a *App) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
	if err := validateGitInput(projectPath, filePath); err != nil {
		return "", err
	}
	if filePath == "" {
		return "", fmt.Errorf("filePath is required")
	}
	if err := validateGitRef(ref); err != nil {
		return "", err
	}
	return a.gitRepo.GetFileAtRef(projectPath, filePath, ref)
}

// BuildContextAtRef builds context from files at a specific git ref without checkout
func (a *App) BuildContextAtRef(projectPath string, files []string, ref string, optionsJson string) (string, error) {
	if err := validateProjectRoot(projectPath); err != nil {
		return "", err
	}
	if err := validateGitRef(ref); err != nil {
		return "", err
	}
	// Validate all file paths
	for _, file := range files {
		if err := validateGitInput(projectPath, file); err != nil {
			return "", err
		}
	}

	var contents []string

	for _, file := range files {
		content, err := a.gitRepo.GetFileAtRef(projectPath, file, ref)
		if err != nil {
			a.log.Info(fmt.Sprintf("Warning: Failed to read file %s at ref %s: %v", file, ref, err))
			continue
		}
		contents = append(contents, fmt.Sprintf("// File: %s (ref: %s)\n%s", file, ref, content))
	}

	result := strings.Join(contents, "\n\n---\n\n")
	return result, nil
}

// GetGitignoreContentForProject returns .gitignore content for a project
func (a *App) GetGitignoreContentForProject(projectPath string) (string, error) {
	if err := validateProjectRoot(projectPath); err != nil {
		return "", err
	}
	content, err := a.gitRepo.GetGitignoreContent(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to get .gitignore content: %w", err)
	}
	return content, nil
}

// AddToGitignore adds a pattern to .gitignore file
func (a *App) AddToGitignore(projectPath string, pattern string) error {
	if err := validateProjectRoot(projectPath); err != nil {
		return err
	}
	if pattern == "" {
		return fmt.Errorf("pattern is required")
	}
	gitignorePath := filepath.Join(projectPath, ".gitignore")

	content, err := os.ReadFile(gitignorePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to read .gitignore: %w", err)
	}

	newContent := string(content)
	if !strings.HasSuffix(newContent, "\n") && newContent != "" {
		newContent += "\n"
	}
	newContent += pattern + "\n"

	if err := os.WriteFile(gitignorePath, []byte(newContent), 0o600); err != nil {
		return fmt.Errorf("failed to write .gitignore: %w", err)
	}

	return nil
}

// ============ GitHub API ENDPOINTS ============

// IsGitHubURL checks if URL is a GitHub repository
func (a *App) IsGitHubURL(url string) bool {
	return git.IsGitHubURL(url)
}

// GitHubGetBranches returns branches for a GitHub repository via API
func (a *App) GitHubGetBranches(repoURL string) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	branches, err := api.GetBranches(repo.Owner, repo.Name)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(branches)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GitHubGetCommits returns commits for a GitHub repository branch via API
func (a *App) GitHubGetCommits(repoURL, branch string, limit int) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	commits, err := api.GetCommits(repo.Owner, repo.Name, branch, limit)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(commits)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GitHubListFiles returns file list for a GitHub repository at specific ref via API
func (a *App) GitHubListFiles(repoURL, ref string) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	files, err := api.ListFiles(repo.Owner, repo.Name, ref)
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// GitHubGetFileContent returns file content from GitHub repository via API
func (a *App) GitHubGetFileContent(repoURL, filePath, ref string) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	content, err := api.GetRawFileContent(repo.Owner, repo.Name, filePath, ref)
	if err != nil {
		return "", err
	}

	return content, nil
}

// GitHubBuildContext builds context from GitHub files via API (no clone)
func (a *App) GitHubBuildContext(repoURL string, files []string, ref string) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	var contents []string
	for _, file := range files {
		content, err := api.GetRawFileContent(repo.Owner, repo.Name, file, ref)
		if err != nil {
			a.log.Info(fmt.Sprintf("Warning: Failed to read GitHub file %s: %v", file, err))
			continue
		}
		contents = append(contents, fmt.Sprintf("// File: %s (GitHub: %s/%s@%s)\n%s", file, repo.Owner, repo.Name, ref, content))
	}

	result := strings.Join(contents, "\n\n---\n\n")
	return result, nil
}

// GitHubGetDefaultBranch returns the default branch for a GitHub repository
func (a *App) GitHubGetDefaultBranch(repoURL string) (string, error) {
	api := git.NewGitHubAPI()

	repo, err := git.ParseGitHubURL(repoURL)
	if err != nil {
		return "", err
	}

	return api.GetDefaultBranch(repo.Owner, repo.Name)
}

// ============ GitLab API ENDPOINTS ============

// IsGitLabURL checks if URL is a GitLab repository
func (a *App) IsGitLabURL(url string) bool {
	return git.IsGitLabURL(url)
}

// GitLabGetBranches returns branches for a GitLab repository
func (a *App) GitLabGetBranches(repoURL string) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	branches, err := api.GetBranches(repo.Host, repo.Namespace, repo.Name)
	if err != nil {
		return "", fmt.Errorf("failed to get branches: %w", err)
	}

	branchesJson, err := json.Marshal(branches)
	if err != nil {
		return "", fmt.Errorf("failed to marshal branches: %w", err)
	}

	return string(branchesJson), nil
}

// GitLabGetDefaultBranch returns the default branch for a GitLab repository
func (a *App) GitLabGetDefaultBranch(repoURL string) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	return api.GetDefaultBranch(repo.Host, repo.Namespace, repo.Name)
}

// GitLabGetCommits returns commits for a GitLab repository branch
func (a *App) GitLabGetCommits(repoURL, branch string, limit int) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	commits, err := api.GetCommits(repo.Host, repo.Namespace, repo.Name, branch, limit)
	if err != nil {
		return "", fmt.Errorf("failed to get commits: %w", err)
	}

	commitsJson, err := json.Marshal(commits)
	if err != nil {
		return "", fmt.Errorf("failed to marshal commits: %w", err)
	}

	return string(commitsJson), nil
}

// GitLabListFiles returns list of files in a GitLab repository at a specific ref
func (a *App) GitLabListFiles(repoURL, ref string) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	files, err := api.ListFiles(repo.Host, repo.Namespace, repo.Name, ref)
	if err != nil {
		return "", fmt.Errorf("failed to list files: %w", err)
	}

	filesJson, err := json.Marshal(files)
	if err != nil {
		return "", fmt.Errorf("failed to marshal files: %w", err)
	}

	return string(filesJson), nil
}

// GitLabGetFileContent returns content of a file from GitLab repository
func (a *App) GitLabGetFileContent(repoURL, filePath, ref string) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	return api.GetFileContent(repo.Host, repo.Namespace, repo.Name, filePath, ref)
}

// GitLabBuildContext builds context from GitLab repository files
func (a *App) GitLabBuildContext(repoURL string, files []string, ref string) (string, error) {
	repo, err := git.ParseGitLabURL(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse GitLab URL: %w", err)
	}

	api := git.NewGitLabAPI()
	var contextBuilder strings.Builder

	for _, filePath := range files {
		content, err := api.GetFileContent(repo.Host, repo.Namespace, repo.Name, filePath, ref)
		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to get file %s: %v", filePath, err))
			continue
		}

		contextBuilder.WriteString(fmt.Sprintf("// File: %s\n", filePath))
		contextBuilder.WriteString(content)
		contextBuilder.WriteString("\n\n")
	}

	return contextBuilder.String(), nil
}

// ============ GIT CACHE MANAGEMENT ============

// ClearGitCache clears all cached git data
func (a *App) ClearGitCache() error {
	if repo, ok := a.gitRepo.(interface{ ClearCache() }); ok {
		repo.ClearCache()
		return nil
	}
	return fmt.Errorf("git repository does not support cache clearing")
}

// InvalidateGitCacheForProject invalidates all cache entries for a specific project
func (a *App) InvalidateGitCacheForProject(projectPath string) error {
	if err := validateProjectRoot(projectPath); err != nil {
		return err
	}

	if repo, ok := a.gitRepo.(interface{ InvalidateProjectCache(string) }); ok {
		repo.InvalidateProjectCache(projectPath)
		return nil
	}
	return fmt.Errorf("git repository does not support cache invalidation")
}

// GetGitCacheStats returns cache statistics
func (a *App) GetGitCacheStats() (string, error) {
	if repo, ok := a.gitRepo.(interface{ GetCacheStats() map[string]interface{} }); ok {
		stats := repo.GetCacheStats()
		result, err := json.Marshal(stats)
		if err != nil {
			return "", fmt.Errorf("failed to marshal cache stats: %w", err)
		}
		return string(result), nil
	}
	return "", fmt.Errorf("git repository does not support cache statistics")
}
