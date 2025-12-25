package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/infrastructure/git"
	"strings"
)

// IsGitAvailable checks if git is available on the system
func (a *App) IsGitAvailable() bool {
	return a.projectHandler.IsGitAvailable()
}

// IsGitRepository checks if the given path is a git repository
func (a *App) IsGitRepository(projectPath string) bool {
	return a.gitRepo.IsGitRepository(projectPath)
}

// GetUncommittedFiles returns list of uncommitted files in a git repository
func (a *App) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	return a.projectHandler.GetUncommittedFiles(projectRoot)
}

// GetRichCommitHistory returns commit history with file changes
func (a *App) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	return a.projectHandler.GetRichCommitHistory(projectRoot, branchName, limit)
}

// GetFileContentAtCommit returns file content at a specific commit
func (a *App) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	return a.projectHandler.GetFileContentAtCommit(projectRoot, filePath, commitHash)
}

// GetGitignoreContent returns the content of .gitignore file
func (a *App) GetGitignoreContent(projectRoot string) (string, error) {
	return a.projectHandler.GetGitignoreContent(projectRoot)
}

// GetBranches returns all git branches
func (a *App) GetBranches(projectRoot string) (string, error) {
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
	branch, err := a.gitRepo.GetCurrentBranch(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return branch, nil
}

// CloneRepository clones a remote git repository
func (a *App) CloneRepository(url string) (string, error) {
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
	return a.gitRepo.CheckoutBranch(projectPath, branch)
}

// CheckoutCommit switches to a specific commit in a git repository
func (a *App) CheckoutCommit(projectPath, commitHash string) error {
	return a.gitRepo.CheckoutCommit(projectPath, commitHash)
}

// GetCommitHistory returns recent commits for selection
func (a *App) GetCommitHistory(projectPath string, limit int) (string, error) {
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
	return a.gitRepo.GetFileAtRef(projectPath, filePath, ref)
}

// BuildContextAtRef builds context from files at a specific git ref without checkout
func (a *App) BuildContextAtRef(projectPath string, files []string, ref string, optionsJson string) (string, error) {
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
	content, err := a.gitRepo.GetGitignoreContent(projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to get .gitignore content: %w", err)
	}
	return content, nil
}

// AddToGitignore adds a pattern to .gitignore file
func (a *App) AddToGitignore(projectPath string, pattern string) error {
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
