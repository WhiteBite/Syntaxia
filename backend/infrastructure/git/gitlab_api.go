package git

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// GitLabAPI provides access to GitLab repositories via REST API
type GitLabAPI struct {
	client  *http.Client
	baseURL string
}

// GitLabRepo represents parsed GitLab repository info
type GitLabRepo struct {
	Host      string // gitlab.com or self-hosted
	Namespace string // user or group/subgroup
	Name      string
}

// GitLabBranch represents a branch from GitLab API
type GitLabBranch struct {
	Name   string `json:"name"`
	Commit struct {
		ID string `json:"id"`
	} `json:"commit"`
	Default bool `json:"default"`
}

// GitLabCommit represents a commit from GitLab API
type GitLabCommit struct {
	ID            string `json:"id"`
	ShortID       string `json:"short_id"`
	Title         string `json:"title"`
	Message       string `json:"message"`
	AuthorName    string `json:"author_name"`
	CommittedDate string `json:"committed_date"`
}

// GitLabTreeEntry represents a file/folder in GitLab tree
type GitLabTreeEntry struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"` // "blob" or "tree"
	Path string `json:"path"`
	Mode string `json:"mode"`
}

// NewGitLabAPI creates a new GitLab API client
func NewGitLabAPI() *GitLabAPI {
	return &GitLabAPI{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://gitlab.com/api/v4",
	}
}

// ParseGitLabURL extracts host, namespace and repo name from GitLab URL
func ParseGitLabURL(urlStr string) (*GitLabRepo, error) {
	// Support various GitLab URL formats:
	// https://gitlab.com/namespace/repo
	// https://gitlab.com/group/subgroup/repo
	// https://gitlab.com/namespace/repo.git
	// git@gitlab.com:namespace/repo.git
	// https://self-hosted.gitlab.com/namespace/repo

	// Try HTTPS format
	httpsPattern := regexp.MustCompile(`https?://([^/]+)/(.+?)(?:\.git)?/?$`)
	if matches := httpsPattern.FindStringSubmatch(urlStr); len(matches) >= 3 {
		host := matches[1]
		pathParts := strings.Split(strings.TrimSuffix(matches[2], ".git"), "/")
		if len(pathParts) >= 2 {
			name := pathParts[len(pathParts)-1]
			namespace := strings.Join(pathParts[:len(pathParts)-1], "/")
			return &GitLabRepo{
				Host:      host,
				Namespace: namespace,
				Name:      name,
			}, nil
		}
	}

	// Try SSH format
	sshPattern := regexp.MustCompile(`git@([^:]+):(.+?)(?:\.git)?$`)
	if matches := sshPattern.FindStringSubmatch(urlStr); len(matches) >= 3 {
		host := matches[1]
		pathParts := strings.Split(strings.TrimSuffix(matches[2], ".git"), "/")
		if len(pathParts) >= 2 {
			name := pathParts[len(pathParts)-1]
			namespace := strings.Join(pathParts[:len(pathParts)-1], "/")
			return &GitLabRepo{
				Host:      host,
				Namespace: namespace,
				Name:      name,
			}, nil
		}
	}

	return nil, fmt.Errorf("invalid GitLab URL: %s", urlStr)
}

// IsGitLabURL checks if the URL is a GitLab repository
func IsGitLabURL(urlStr string) bool {
	return strings.Contains(urlStr, "gitlab.com") || strings.Contains(urlStr, "gitlab.")
}

// getAPIBaseURL returns the API base URL for the given host
func (g *GitLabAPI) getAPIBaseURL(host string) string {
	if host == "gitlab.com" {
		return "https://gitlab.com/api/v4"
	}
	return fmt.Sprintf("https://%s/api/v4", host)
}

// getProjectPath returns URL-encoded project path
func getProjectPath(namespace, name string) string {
	return url.PathEscape(namespace + "/" + name)
}

// GetBranches returns all branches for a repository
func (g *GitLabAPI) GetBranches(host, namespace, name string) ([]GitLabBranch, error) {
	baseURL := g.getAPIBaseURL(host)
	projectPath := getProjectPath(namespace, name)
	apiURL := fmt.Sprintf("%s/projects/%s/repository/branches?per_page=100", baseURL, projectPath)

	resp, err := g.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch branches: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %s", resp.Status)
	}

	var branches []GitLabBranch
	if err := json.NewDecoder(resp.Body).Decode(&branches); err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}

	return branches, nil
}

// GetDefaultBranch returns the default branch name
func (g *GitLabAPI) GetDefaultBranch(host, namespace, name string) (string, error) {
	baseURL := g.getAPIBaseURL(host)
	projectPath := getProjectPath(namespace, name)
	apiURL := fmt.Sprintf("%s/projects/%s", baseURL, projectPath)

	resp, err := g.client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch project info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitLab API error: %s", resp.Status)
	}

	var projectInfo struct {
		DefaultBranch string `json:"default_branch"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&projectInfo); err != nil {
		return "", fmt.Errorf("failed to decode project info: %w", err)
	}

	return projectInfo.DefaultBranch, nil
}

// GetCommits returns recent commits for a branch
func (g *GitLabAPI) GetCommits(host, namespace, name, branch string, limit int) ([]GitLabCommit, error) {
	baseURL := g.getAPIBaseURL(host)
	projectPath := getProjectPath(namespace, name)
	apiURL := fmt.Sprintf("%s/projects/%s/repository/commits?ref_name=%s&per_page=%d",
		baseURL, projectPath, url.QueryEscape(branch), limit)

	resp, err := g.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch commits: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitLab API error: %s", resp.Status)
	}

	var commits []GitLabCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	return commits, nil
}

// GetTree returns the file tree for a specific ref (branch/commit)
func (g *GitLabAPI) GetTree(host, namespace, name, ref string) ([]GitLabTreeEntry, error) {
	baseURL := g.getAPIBaseURL(host)
	projectPath := getProjectPath(namespace, name)

	var allEntries []GitLabTreeEntry
	page := 1
	perPage := 100

	for {
		apiURL := fmt.Sprintf("%s/projects/%s/repository/tree?ref=%s&recursive=true&per_page=%d&page=%d",
			baseURL, projectPath, url.QueryEscape(ref), perPage, page)

		resp, err := g.client.Get(apiURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch tree: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("GitLab API error: %s", resp.Status)
		}

		var entries []GitLabTreeEntry
		if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("failed to decode tree: %w", err)
		}
		resp.Body.Close()

		allEntries = append(allEntries, entries...)

		if len(entries) < perPage {
			break
		}
		page++
	}

	return allEntries, nil
}

// GetFileContent returns the content of a file at a specific ref
func (g *GitLabAPI) GetFileContent(host, namespace, name, path, ref string) (string, error) {
	baseURL := g.getAPIBaseURL(host)
	projectPath := getProjectPath(namespace, name)
	encodedPath := url.PathEscape(path)
	apiURL := fmt.Sprintf("%s/projects/%s/repository/files/%s/raw?ref=%s",
		baseURL, projectPath, encodedPath, url.QueryEscape(ref))

	resp, err := g.client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitLab API error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body: %w", err)
	}

	return string(body), nil
}

// ListFiles returns list of file paths at a specific ref
func (g *GitLabAPI) ListFiles(host, namespace, name, ref string) ([]string, error) {
	entries, err := g.GetTree(host, namespace, name, ref)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.Type == "blob" {
			files = append(files, entry.Path)
		}
	}

	return files, nil
}
