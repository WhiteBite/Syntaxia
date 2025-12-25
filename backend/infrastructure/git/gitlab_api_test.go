package git

import (
	"testing"
)

func TestParseGitLabURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantHost  string
		wantNS    string
		wantName  string
		wantError bool
	}{
		{
			name:     "Standard HTTPS URL",
			url:      "https://gitlab.com/user/repo",
			wantHost: "gitlab.com",
			wantNS:   "user",
			wantName: "repo",
		},
		{
			name:     "HTTPS URL with .git",
			url:      "https://gitlab.com/user/repo.git",
			wantHost: "gitlab.com",
			wantNS:   "user",
			wantName: "repo",
		},
		{
			name:     "Group/subgroup URL",
			url:      "https://gitlab.com/group/subgroup/repo",
			wantHost: "gitlab.com",
			wantNS:   "group/subgroup",
			wantName: "repo",
		},
		{
			name:     "SSH URL",
			url:      "git@gitlab.com:user/repo.git",
			wantHost: "gitlab.com",
			wantNS:   "user",
			wantName: "repo",
		},
		{
			name:     "Self-hosted GitLab",
			url:      "https://gitlab.company.com/team/project",
			wantHost: "gitlab.company.com",
			wantNS:   "team",
			wantName: "project",
		},
		{
			name:      "Invalid URL",
			url:       "not-a-url",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo, err := ParseGitLabURL(tt.url)

			if tt.wantError {
				if err == nil {
					t.Errorf("ParseGitLabURL() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseGitLabURL() unexpected error: %v", err)
				return
			}

			if repo.Host != tt.wantHost {
				t.Errorf("ParseGitLabURL() Host = %v, want %v", repo.Host, tt.wantHost)
			}
			if repo.Namespace != tt.wantNS {
				t.Errorf("ParseGitLabURL() Namespace = %v, want %v", repo.Namespace, tt.wantNS)
			}
			if repo.Name != tt.wantName {
				t.Errorf("ParseGitLabURL() Name = %v, want %v", repo.Name, tt.wantName)
			}
		})
	}
}

func TestIsGitLabURL(t *testing.T) {
	tests := []struct {
		url  string
		want bool
	}{
		{"https://gitlab.com/user/repo", true},
		{"https://gitlab.company.com/user/repo", true},
		{"git@gitlab.com:user/repo.git", true},
		{"https://github.com/user/repo", false},
		{"https://bitbucket.org/user/repo", false},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			if got := IsGitLabURL(tt.url); got != tt.want {
				t.Errorf("IsGitLabURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestGetProjectPath(t *testing.T) {
	tests := []struct {
		namespace string
		name      string
		want      string
	}{
		{"user", "repo", "user%2Frepo"},
		{"group/subgroup", "project", "group%2Fsubgroup%2Fproject"},
	}

	for _, tt := range tests {
		t.Run(tt.namespace+"/"+tt.name, func(t *testing.T) {
			got := getProjectPath(tt.namespace, tt.name)
			if got != tt.want {
				t.Errorf("getProjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
