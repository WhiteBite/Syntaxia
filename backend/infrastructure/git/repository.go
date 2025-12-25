package git

import (
	"bufio"
	"fmt"
	"os/exec"
	"syntaxia/domain"
	"syntaxia/internal/executil"
	"strings"
	"time"
)

type Repository struct {
	log domain.Logger
}

func New(log domain.Logger) domain.GitRepository {
	return &Repository{
		log: log,
	}
}

func (r *Repository) IsGitAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

func (r *Repository) GetUncommittedFiles(projectRoot string) ([]domain.FileStatus, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}

	files := parseGitStatus(string(output))
	r.log.Info(fmt.Sprintf("Found %d uncommitted files in %s.", len(files), projectRoot))
	return files, nil
}

// parseGitStatus parses git status --porcelain output
func parseGitStatus(output string) []domain.FileStatus {
	var files []domain.FileStatus
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		if file := parseStatusLine(scanner.Text()); file != nil {
			files = append(files, *file)
		}
	}
	return files
}

// parseStatusLine parses a single git status line
func parseStatusLine(line string) *domain.FileStatus {
	if len(line) < 3 {
		return nil
	}
	status := strings.TrimSpace(line[:2])
	path := extractFilePath(status, strings.TrimSpace(line[3:]))
	return &domain.FileStatus{Path: path, Status: mapGitStatus(status)}
}

// extractFilePath handles renamed files
func extractFilePath(status, path string) string {
	if strings.HasPrefix(status, "R") {
		if parts := strings.Split(path, " -> "); len(parts) == 2 {
			return parts[1]
		}
	}
	return path
}

// mapGitStatus maps git status codes to simple codes
func mapGitStatus(status string) string {
	switch {
	case strings.HasPrefix(status, "M") || strings.Contains(status, "M"):
		return "M"
	case strings.HasPrefix(status, "A"):
		return "A"
	case strings.HasPrefix(status, "D"):
		return "D"
	case strings.HasPrefix(status, "R"):
		return "R"
	case strings.HasPrefix(status, "C"):
		return "C"
	case status == "??":
		return "U"
	case strings.HasPrefix(status, "U"):
		return "UM"
	default:
		return status
	}
}

func (r *Repository) GetRichCommitHistory(projectRoot, branchName string, limit int) ([]domain.CommitWithFiles, error) {
	args := []string{
		"log",
		"--pretty=format:COMMIT %H %P%n%s%n%an%n%cI", // Include Author Name (%an) and Committer Date, ISO 8601 format (%cI)
		"--name-status",
		"--topo-order",
		fmt.Sprintf("--max-count=%d", limit),
	}

	if branchName != "" {
		args = append(args, branchName)
	}

	cmd := exec.Command("git", args...)
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git log: %w", err)
	}

	commits, parseErr := ParseRichLogOutput(string(output))
	if parseErr != nil {
		return nil, parseErr
	}
	r.log.Info(fmt.Sprintf("Loaded and parsed %d commits for %s.", len(commits), projectRoot))
	return commits, nil
}

func (r *Repository) GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", commitHash, filePath)) //nolint:gosec // Git command
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file content at commit %s:%s: %w", commitHash, filePath, err)
	}

	return string(output), nil
}

func (r *Repository) GetGitignoreContent(projectRoot string) (string, error) {
	cmd := exec.Command("git", "show", "HEAD:.gitignore")
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		// .gitignore might not exist or not be tracked, which is not an error for us.
		r.log.Debug(fmt.Sprintf("Could not get .gitignore content for %s (might not exist or be tracked): %v", projectRoot, err))
		return "", nil
	}

	return string(output), nil
}

// ParseRichLogOutput parses the output of git log with --pretty and --name-status
func ParseRichLogOutput(gitOutput string) ([]domain.CommitWithFiles, error) {
	var commits []domain.CommitWithFiles
	var currentCommit *domain.CommitWithFiles
	scanner := bufio.NewScanner(strings.NewReader(gitOutput))

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "COMMIT ") {
			if currentCommit != nil {
				commits = append(commits, *currentCommit)
			}
			currentCommit = parseCommitHeader(line)
		} else if currentCommit != nil {
			parseCommitLine(currentCommit, line)
		}
	}

	if currentCommit != nil {
		commits = append(commits, *currentCommit)
	}
	return commits, nil
}

// parseCommitHeader parses a COMMIT line
func parseCommitHeader(line string) *domain.CommitWithFiles {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	return &domain.CommitWithFiles{
		Hash:    parts[1],
		IsMerge: len(parts) > 3,
		Files:   []string{},
	}
}

// parseCommitLine parses a line belonging to current commit
func parseCommitLine(commit *domain.CommitWithFiles, line string) {
	if strings.Contains(line, "\t") {
		parseFileLine(commit, line)
	} else {
		parseMetadataLine(commit, line)
	}
}

// parseFileLine parses a file status line
func parseFileLine(commit *domain.CommitWithFiles, line string) {
	parts := strings.Split(line, "\t")
	if len(parts) >= 2 && parts[1] != "" {
		commit.Files = append(commit.Files, parts[1])
	}
}

// parseMetadataLine parses subject/author/date lines
func parseMetadataLine(commit *domain.CommitWithFiles, line string) {
	switch {
	case commit.Subject == "":
		commit.Subject = line
	case commit.Author == "":
		commit.Author = line
	case commit.Date == "":
		if t, err := time.Parse(time.RFC3339, line); err == nil {
			commit.Date = t.Format("2006-01-02 15:04")
		} else {
			commit.Date = line
		}
	}
}

// GetBranches returns all git branches
func (r *Repository) GetBranches(projectRoot string) ([]string, error) {
	cmd := exec.Command("git", "branch", "-a")
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git branches: %w", err)
	}

	var branches []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Remove the current branch indicator (*)
		if strings.HasPrefix(line, "* ") {
			line = strings.TrimPrefix(line, "* ")
		} else if strings.HasPrefix(line, "  ") {
			line = strings.TrimPrefix(line, "  ")
		}

		// Skip remote tracking references for now
		if strings.HasPrefix(line, "remotes/") {
			continue
		}

		// Skip empty lines and special references
		if line != "" && !strings.Contains(line, "HEAD") {
			branches = append(branches, line)
		}
	}

	r.log.Info(fmt.Sprintf("Found %d branches in %s.", len(branches), projectRoot))
	return branches, nil
}

// GetCurrentBranch returns the current git branch
func (r *Repository) GetCurrentBranch(projectRoot string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current git branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	r.log.Info(fmt.Sprintf("Current branch in %s: %s", projectRoot, branch))
	return branch, nil
}

// GetAllFiles returns a list of all files in the repository.
func (r *Repository) GetAllFiles(projectPath string) ([]string, error) {
	cmd := exec.Command("git", "ls-files")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	return files, nil
}

// GenerateDiff generates a git diff between HEAD and HEAD~1.
func (r *Repository) GenerateDiff(projectPath string) (string, error) {
	cmd := exec.Command("git", "diff", "HEAD~1", "HEAD")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		// If HEAD~1 does not exist (e.g., first commit), try diffing against the empty tree
		cmd = exec.Command("git", "diff", "4b825dc642cb6eb9a060e54bf8d69288fbee4904", "HEAD") // magic empty tree hash
		executil.HideWindow(cmd)
		cmd.Dir = projectPath
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to generate diff: %w", err)
		}
	}

	return string(output), nil
}

// IsGitRepository checks if the given path is a git repository
func (r *Repository) IsGitRepository(projectPath string) bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath
	err := cmd.Run()
	return err == nil
}

// CloneRepository clones a remote repository to a local path (shallow clone for speed)
func (r *Repository) CloneRepository(url, targetPath string, depth int) error {
	args := []string{"clone"}
	if depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", depth))
	}
	args = append(args, url, targetPath)

	cmd := exec.Command("git", args...)
	executil.HideWindow(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to clone repository: %s - %w", string(output), err)
	}

	r.log.Info(fmt.Sprintf("Cloned repository %s to %s", url, targetPath))
	return nil
}

// CheckoutBranch switches to a specific branch
func (r *Repository) CheckoutBranch(projectPath, branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout branch %s: %s - %w", branch, string(output), err)
	}

	r.log.Info(fmt.Sprintf("Checked out branch %s in %s", branch, projectPath))
	return nil
}

// CheckoutCommit switches to a specific commit (detached HEAD)
func (r *Repository) CheckoutCommit(projectPath, commitHash string) error {
	cmd := exec.Command("git", "checkout", commitHash)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to checkout commit %s: %s - %w", commitHash, string(output), err)
	}

	r.log.Info(fmt.Sprintf("Checked out commit %s in %s", commitHash, projectPath))
	return nil
}

// ListFilesAtRef returns list of files at a specific branch or commit without checkout
func (r *Repository) ListFilesAtRef(projectPath, ref string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", ref)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list files at ref %s: %w", ref, err)
	}

	var files []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		files = append(files, scanner.Text())
	}

	r.log.Info(fmt.Sprintf("Listed %d files at ref %s in %s", len(files), ref, projectPath))
	return files, nil
}

// GetFileAtRef returns file content at a specific branch or commit without checkout
func (r *Repository) GetFileAtRef(projectPath, filePath, ref string) (string, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", ref, filePath)) //nolint:gosec // Git command
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file %s at ref %s: %w", filePath, ref, err)
	}

	return string(output), nil
}

// GetTreeAtRef returns file tree structure at a specific ref
func (r *Repository) GetTreeAtRef(projectPath, ref string) ([]GitTreeEntry, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--long", ref)
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get tree at ref %s: %w", ref, err)
	}

	var entries []GitTreeEntry
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Format: <mode> <type> <hash> <size>\t<path>
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			continue
		}

		meta := strings.Fields(parts[0])
		if len(meta) < 4 {
			continue
		}

		size := int64(0)
		if meta[3] != "-" {
			_, _ = fmt.Sscanf(meta[3], "%d", &size)
		}

		entries = append(entries, GitTreeEntry{
			Path:  parts[1],
			Type:  meta[1],
			Size:  size,
			IsDir: meta[1] == "tree",
		})
	}

	return entries, nil
}

// GitTreeEntry represents a file/folder in git tree
type GitTreeEntry struct {
	Path  string `json:"path"`
	Type  string `json:"type"`
	Size  int64  `json:"size"`
	IsDir bool   `json:"isDir"`
}

// GetCommitHistory returns recent commits with hash and subject
func (r *Repository) GetCommitHistory(projectPath string, limit int) ([]domain.CommitInfo, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%H|%s|%an|%cI", fmt.Sprintf("-n%d", limit)) //nolint:gosec // Git command
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %w", err)
	}

	var commits []domain.CommitInfo
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "|", 4)
		if len(parts) >= 4 {
			commits = append(commits, domain.CommitInfo{
				Hash:    parts[0],
				Subject: parts[1],
				Author:  parts[2],
				Date:    parts[3],
			})
		}
	}

	return commits, nil
}

// FetchRemoteBranches fetches and returns all remote branches
func (r *Repository) FetchRemoteBranches(projectPath string) ([]string, error) {
	// First fetch all remotes
	fetchCmd := exec.Command("git", "fetch", "--all")
	executil.HideWindow(fetchCmd)
	fetchCmd.Dir = projectPath
	_ = fetchCmd.Run() // Ignore errors, might not have network

	// Get remote branches
	cmd := exec.Command("git", "branch", "-r")
	executil.HideWindow(cmd)
	cmd.Dir = projectPath

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get remote branches: %w", err)
	}

	var branches []string
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.Contains(line, "HEAD") {
			continue
		}
		// Remove "origin/" prefix for cleaner display
		line = strings.TrimPrefix(line, "origin/")
		branches = append(branches, line)
	}

	return branches, nil
}
