package version

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	// GitHub repository info
	RepoOwner = "WhiteBite"
	RepoName  = "syntaxia"

	// Cache duration
	CacheDuration = 1 * time.Hour

	// API timeout
	APITimeout = 10 * time.Second
)

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"` // Markdown changelog
	PublishedAt time.Time `json:"published_at"`
	HTMLURL     string    `json:"html_url"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
}

// ReleasesResponse for frontend
type ReleasesResponse struct {
	CurrentVersion string    `json:"currentVersion"`
	LatestVersion  string    `json:"latestVersion"`
	HasUpdate      bool      `json:"hasUpdate"`
	Releases       []Release `json:"releases"`
	LastChecked    time.Time `json:"lastChecked"`
	Error          string    `json:"error,omitempty"`
}

// ReleasesService handles GitHub releases
type ReleasesService struct {
	httpClient *http.Client
	cache      *releasesCache
}

type releasesCache struct {
	mu        sync.RWMutex
	releases  []Release
	fetchedAt time.Time
}

// NewReleasesService creates a new releases service
func NewReleasesService() *ReleasesService {
	return &ReleasesService{
		httpClient: &http.Client{
			Timeout: APITimeout,
		},
		cache: &releasesCache{},
	}
}

// GetReleases returns releases with caching
func (s *ReleasesService) GetReleases(ctx context.Context) (*ReleasesResponse, error) {
	// Check cache first
	s.cache.mu.RLock()
	if time.Since(s.cache.fetchedAt) < CacheDuration && len(s.cache.releases) > 0 {
		releases := s.cache.releases
		s.cache.mu.RUnlock()
		return s.buildResponse(releases, nil), nil
	}
	s.cache.mu.RUnlock()

	// Fetch from GitHub
	releases, err := s.fetchReleases(ctx)
	if err != nil {
		// Return cached data if available, with error
		s.cache.mu.RLock()
		cached := s.cache.releases
		s.cache.mu.RUnlock()
		if len(cached) > 0 {
			return s.buildResponse(cached, err), nil
		}
		return s.buildResponse(nil, err), nil
	}

	// Update cache
	s.cache.mu.Lock()
	s.cache.releases = releases
	s.cache.fetchedAt = time.Now()
	s.cache.mu.Unlock()

	return s.buildResponse(releases, nil), nil
}

// fetchReleases fetches releases from GitHub API
func (s *ReleasesService) fetchReleases(ctx context.Context) ([]Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", RepoOwner, RepoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Syntaxia/"+Version)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("failed to decode releases: %w", err)
	}

	// Filter out drafts
	filtered := make([]Release, 0, len(releases))
	for _, r := range releases {
		if !r.Draft {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

// buildResponse builds the response with version comparison
func (s *ReleasesService) buildResponse(releases []Release, fetchErr error) *ReleasesResponse {
	resp := &ReleasesResponse{
		CurrentVersion: Version,
		Releases:       releases,
		LastChecked:    time.Now(),
	}

	if fetchErr != nil {
		resp.Error = fetchErr.Error()
	}

	// Find latest non-prerelease version
	for _, r := range releases {
		if !r.Prerelease {
			resp.LatestVersion = r.TagName
			resp.HasUpdate = compareVersions(Version, r.TagName)
			break
		}
	}

	return resp
}

// compareVersions returns true if latest > current
func compareVersions(current, latest string) bool {
	// Simple comparison - if current is "dev", always show update
	if current == "dev" {
		return false // Don't show update for dev builds
	}

	// Strip 'v' prefix for comparison
	if current != "" && current[0] == 'v' {
		current = current[1:]
	}
	if latest != "" && latest[0] == 'v' {
		latest = latest[1:]
	}

	return latest > current
}

// GetLatestRelease returns only the latest release
func (s *ReleasesService) GetLatestRelease(ctx context.Context) (*Release, error) {
	resp, err := s.GetReleases(ctx)
	if err != nil {
		return nil, err
	}

	for _, r := range resp.Releases {
		if !r.Prerelease {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("no releases found")
}

// ClearCache clears the releases cache
func (s *ReleasesService) ClearCache() {
	s.cache.mu.Lock()
	s.cache.releases = nil
	s.cache.fetchedAt = time.Time{}
	s.cache.mu.Unlock()
}
