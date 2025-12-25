package git

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GitHubAPI provides access to GitHub repositories via REST API
type GitHubAPI struct {
	client  *http.Client
	baseURL string
}

// GitHubRepo represents parsed GitHub repository info
type GitHubRepo struct {
	Owner string
	Name  string
}

// GitHubBranch represents a branch from GitHub API
type GitHubBranch struct {
	Name   string `json:"name"`
	Commit struct {
		SHA string `json:"sha"`
	} `json:"commit"`
}

// GitHubCommit represents a commit from GitHub API
type GitHubCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name string `json:"name"`
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

// GitHubTreeEntry represents a file/folder in GitHub tree
type GitHubTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"` // "blob" or "tree"
	SHA  string `json:"sha"`
	Size int64  `json:"size,omitempty"`
}

// GitHubTree represents the tree response from GitHub API
type GitHubTree struct {
	SHA       string            `json:"sha"`
	Tree      []GitHubTreeEntry `json:"tree"`
	Truncated bool              `json:"truncated"`
}

// GitHubContent represents file content from GitHub API
type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	SHA         string `json:"sha"`
	Size        int64  `json:"size"`
	Type        string `json:"type"`
	Content     string `json:"content"`
	Encoding    string `json:"encoding"`
	DownloadURL string `json:"download_url"`
}

// NewGitHubAPI creates a new GitHub API client
func NewGitHubAPI() *GitHubAPI {
	return &GitHubAPI{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.github.com",
	}
}

// ParseGitHubURL extracts owner and repo name from GitHub URL
func ParseGitHubURL(url string) (*GitHubRepo, error) {
	// Support various GitHub URL formats:
	// https://github.com/owner/repo
	// https://github.com/owner/repo.git
	// git@github.com:owner/repo.git
	// github.com/owner/repo

	patterns := []string{
		`github\.com[/:]([^/]+)/([^/\.]+)(?:\.git)?`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(url)
		if len(matches) >= 3 {
			return &GitHubRepo{
				Owner: matches[1],
				Name:  strings.TrimSuffix(matches[2], ".git"),
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid GitHub URL: %s", url)
}

// IsGitHubURL checks if the URL is a GitHub repository
func IsGitHubURL(url string) bool {
	return strings.Contains(url, "github.com")
}

// GetBranches returns all branches for a repository
func (g *GitHubAPI) GetBranches(owner, repo string) ([]GitHubBranch, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/branches?per_page=100", g.baseURL, owner, repo)

	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var branches []GitHubBranch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}

	return branches, nil
}

// GetDefaultBranch returns the default branch name
func (g *GitHubAPI) GetDefaultBranch(owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", g.baseURL, owner, repo)

	resp, err := g.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch repo info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var repoInfo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return "", fmt.Errorf("failed to decode repo info: %w", err)
	}

	return repoInfo.DefaultBranch, nil
}

// GetCommits returns recent commits for a branch
func (g *GitHubAPI) GetCommits(owner, repo, branch string, limit int) ([]GitHubCommit, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?sha=%s&per_page=%d", g.baseURL, owner, repo, branch, limit)

	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var commits []GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	return commits, nil
}

// GetTree returns the file tree for a specific ref (branch/commit)
func (g *GitHubAPI) GetTree(owner, repo, ref string) (*GitHubTree, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", g.baseURL, owner, repo, ref)

	resp, err := g.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tree: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var tree GitHubTree
	if err := json.NewDecoder(resp.Body).Decode(&tree); err != nil {
		return nil, fmt.Errorf("failed to decode tree: %w", err)
	}

	return &tree, nil
}

// GetFileContent returns the content of a file at a specific ref
func (g *GitHubAPI) GetFileContent(owner, repo, path, ref string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/contents/%s?ref=%s", g.baseURL, owner, repo, path, ref)

	resp, err := g.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	var content GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	// Content is base64 encoded
	if content.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(strings.ReplaceAll(content.Content, "\n", ""))
		if err != nil {
			return "", fmt.Errorf("failed to decode base64 content: %w", err)
		}
		return string(decoded), nil
	}

	return content.Content, nil
}

// GetRawFileContent returns raw file content using download URL (for larger files)
func (g *GitHubAPI) GetRawFileContent(owner, repo, path, ref string) (string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, path)

	resp, err := g.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch raw file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("raw content error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}

// ListFiles returns list of file paths at a specific ref
func (g *GitHubAPI) ListFiles(owner, repo, ref string) ([]string, error) {
	tree, err := g.GetTree(owner, repo, ref)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range tree.Tree {
		if entry.Type == "blob" {
			files = append(files, entry.Path)
		}
	}

	return files, nil
}
