package git

import (
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantOwner string
		wantRepo  string
		wantErr   bool
	}{
		{
			name:      "https URL",
			url:       "https://github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "https URL with .git",
			url:       "https://github.com/owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "ssh URL",
			url:       "git@github.com:owner/repo.git",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:      "short URL",
			url:       "github.com/owner/repo",
			wantOwner: "owner",
			wantRepo:  "repo",
			wantErr:   false,
		},
		{
			name:    "invalid URL",
			url:     "https://gitlab.com/owner/repo",
			wantErr: true,
		},
		{
			name:    "empty URL",
			url:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := ParseGitHubURL(tt.url)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseGitHubURL() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseGitHubURL() error = %v", err)
				return
			}

			if repo.Owner != tt.wantOwner {
				t.Errorf("ParseGitHubURL() owner = %v, want %v", repo.Owner, tt.wantOwner)
			}

			if repo.Name != tt.wantRepo {
				t.Errorf("ParseGitHubURL() repo = %v, want %v", repo.Name, tt.wantRepo)
			}
		})
	}
}

func TestIsGitHubURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://github.com/owner/repo", true},
		{"https://github.com/owner/repo.git", true},
		{"git@github.com:owner/repo.git", true},
		{"github.com/owner/repo", true},
		{"https://gitlab.com/owner/repo", false},
		{"https://bitbucket.org/owner/repo", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := IsGitHubURL(tt.url); got != tt.want {
				t.Errorf("IsGitHubURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestNewGitHubAPI(t *testing.T) {
	api := NewGitHubAPI()

	if api == nil {
		t.Fatal("NewGitHubAPI() returned nil")
	}

	if api.client == nil {
		t.Error("NewGitHubAPI() client is nil")
	}

	if api.baseURL != "https://api.github.com" {
		t.Errorf("NewGitHubAPI() baseURL = %v, want https://api.github.com", api.baseURL)
	}
}

// Integration tests - skip if no network
func TestGitHubAPI_GetBranches_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	api := NewGitHubAPI()

	// Use a well-known public repo
	branches, err := api.GetBranches("golang", "go")
	if err != nil {
		t.Fatalf("GetBranches() error = %v", err)
	}

	if len(branches) == 0 {
		t.Error("GetBranches() returned empty list")
	}

	// Check that master branch exists
	found := false
	for _, b := range branches {
		if b.Name == "master" {
			found = true
			break
		}
	}

	if !found {
		t.Error("GetBranches() did not find 'master' branch")
	}
}

func TestGitHubAPI_GetDefaultBranch_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	api := NewGitHubAPI()

	branch, err := api.GetDefaultBranch("golang", "go")
	if err != nil {
		t.Fatalf("GetDefaultBranch() error = %v", err)
	}

	if branch == "" {
		t.Error("GetDefaultBranch() returned empty string")
	}

	t.Logf("Default branch: %s", branch)
}

func TestGitHubAPI_ListFiles_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	api := NewGitHubAPI()

	files, err := api.ListFiles("golang", "go", "master")
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}

	if len(files) == 0 {
		t.Error("ListFiles() returned empty list")
	}

	// Check for some expected files
	foundReadme := false
	for _, f := range files {
		if f == "README.md" {
			foundReadme = true
			break
		}
	}

	if !foundReadme {
		t.Log("README.md not found in root (might be in different location)")
	}

	t.Logf("Found %d files", len(files))
}
