package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"syntaxia/domain"
	"syntaxia/infrastructure/settingsfs"
)

// GetSettings returns current application settings
func (a *App) GetSettings() (domain.SettingsDTO, error) {
	return a.settingsHandler.GetSettings()
}

// SaveSettings saves application settings
func (a *App) SaveSettings(settingsJson string) error {
	return a.settingsHandler.SaveSettings(settingsJson)
}

// RefreshAIModels refreshes available AI models for a provider
func (a *App) RefreshAIModels(provider, apiKey string) error {
	return a.settingsHandler.RefreshAIModels(provider, apiKey)
}

// GetRecentProjects returns the list of recently opened projects
func (a *App) GetRecentProjects() (string, error) {
	projects := a.settingsService.GetRecentProjects()
	result, err := json.Marshal(projects)
	if err != nil {
		return "", fmt.Errorf("failed to marshal recent projects: %w", err)
	}
	return string(result), nil
}

// AddRecentProject adds a project to the recent list and saves settings
func (a *App) AddRecentProject(path, name string) error {
	a.settingsService.AddRecentProject(path, name)
	return a.settingsService.Save()
}

// RemoveRecentProject removes a project from the recent list and saves settings
func (a *App) RemoveRecentProject(path string) error {
	a.settingsService.RemoveRecentProject(path)
	return a.settingsService.Save()
}

// GetCustomIgnoreRules returns custom ignore rules from settings
func (a *App) GetCustomIgnoreRules() (string, error) {
	dto, err := a.settingsService.GetSettingsDTO()
	if err != nil {
		return "", fmt.Errorf("failed to get settings: %w", err)
	}
	return dto.CustomIgnoreRules, nil
}

// UpdateCustomIgnoreRules updates custom ignore rules in settings
func (a *App) UpdateCustomIgnoreRules(rules string) error {
	dto, err := a.settingsService.GetSettingsDTO()
	if err != nil {
		return fmt.Errorf("failed to get settings: %w", err)
	}
	dto.CustomIgnoreRules = rules
	return a.settingsService.SaveSettingsDTO(dto)
}

// IgnorePreviewResult contains grouped preview data for performance
type IgnorePreviewResult struct {
	TotalFiles   int                       `json:"totalFiles"`
	ByDirectory  map[string]int            `json:"byDirectory"`
	ByRule       map[string]int            `json:"byRule"`
	SampleFiles  []string                  `json:"sampleFiles"`
	TopDirs      []DirCount                `json:"topDirs"`
}

// DirCount represents directory with file count
type DirCount struct {
	Dir   string `json:"dir"`
	Count int    `json:"count"`
}

// TestIgnoreRules tests ignore rules against project files (legacy, returns flat list)
func (a *App) TestIgnoreRules(projectPath string, rules string) ([]string, error) {
	result, err := a.TestIgnoreRulesDetailed(projectPath, rules)
	if err != nil {
		return nil, err
	}
	return result.SampleFiles, nil
}

// TestIgnoreRulesDetailed returns detailed preview with grouping for performance
func (a *App) TestIgnoreRulesDetailed(projectPath string, rules string) (*IgnorePreviewResult, error) {
	tree, err := a.projectService.ListFiles(projectPath, false, false)
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	// Parse rules once
	parsedRules := a.parseRules(rules)
	if len(parsedRules) == 0 {
		return &IgnorePreviewResult{
			TotalFiles:  0,
			ByDirectory: make(map[string]int),
			ByRule:      make(map[string]int),
			SampleFiles: []string{},
			TopDirs:     []DirCount{},
		}, nil
	}

	// Flatten tree and check matches in one pass
	byDir := make(map[string]int)
	byRule := make(map[string]int)
	var sampleFiles []string
	totalFiles := 0
	const maxSamples = 100

	var process func(nodes []*domain.FileNode)
	process = func(nodes []*domain.FileNode) {
		for _, node := range nodes {
			if !node.IsDir {
				normalizedPath := strings.ReplaceAll(node.RelPath, "\\", "/")
				if matchedRule := a.findMatchingRule(normalizedPath, parsedRules); matchedRule != "" {
					totalFiles++
					// Group by top-level directory
					dir := getTopDir(normalizedPath)
					byDir[dir]++
					byRule[matchedRule]++
					// Collect samples
					if len(sampleFiles) < maxSamples {
						sampleFiles = append(sampleFiles, node.RelPath)
					}
				}
			}
			if node.Children != nil {
				process(node.Children)
			}
		}
	}
	process(tree)

	// Sort directories by count (top 10)
	topDirs := sortDirCounts(byDir, 10)

	return &IgnorePreviewResult{
		TotalFiles:  totalFiles,
		ByDirectory: byDir,
		ByRule:      byRule,
		SampleFiles: sampleFiles,
		TopDirs:     topDirs,
	}, nil
}

// ParsedRule represents a pre-parsed ignore rule
type ParsedRule struct {
	Original string
	Pattern  string
	IsDir    bool
}

// parseRules parses rules once for reuse
func (a *App) parseRules(rules string) []ParsedRule {
	var parsed []ParsedRule
	lines := strings.Split(rules, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pattern := strings.TrimSuffix(line, "/")
		isDir := strings.HasSuffix(line, "/")
		parsed = append(parsed, ParsedRule{Original: line, Pattern: pattern, IsDir: isDir})
	}
	return parsed
}

// findMatchingRule returns the first matching rule or empty string
func (a *App) findMatchingRule(path string, rules []ParsedRule) string {
	for _, rule := range rules {
		if a.matchPattern(path, rule.Pattern, rule.IsDir) {
			return rule.Original
		}
	}
	return ""
}

// getTopDir extracts top-level directory from path
func getTopDir(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		return parts[0]
	}
	return "(root)"
}

// sortDirCounts sorts directories by count and returns top N
func sortDirCounts(byDir map[string]int, limit int) []DirCount {
	dirs := make([]DirCount, 0, len(byDir))
	for dir, count := range byDir {
		dirs = append(dirs, DirCount{Dir: dir, Count: count})
	}
	// Simple bubble sort for small lists
	for i := 0; i < len(dirs)-1; i++ {
		for j := i + 1; j < len(dirs); j++ {
			if dirs[j].Count > dirs[i].Count {
				dirs[i], dirs[j] = dirs[j], dirs[i]
			}
		}
	}
	if len(dirs) > limit {
		return dirs[:limit]
	}
	return dirs
}

// matchesIgnorePattern checks if a path matches any ignore pattern
func (a *App) matchesIgnorePattern(path string, rules string) bool {
	if rules == "" {
		return false
	}

	// Normalize path separators
	normalizedPath := strings.ReplaceAll(path, "\\", "/")

	lines := strings.Split(rules, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Remove trailing slash for directory patterns
		pattern := strings.TrimSuffix(line, "/")
		isDir := strings.HasSuffix(line, "/")

		// Check different matching strategies
		if a.matchPattern(normalizedPath, pattern, isDir) {
			return true
		}
	}

	return false
}

// matchPattern checks if path matches a single pattern
func (a *App) matchPattern(path, pattern string, isDir bool) bool {
	// Exact filename match (e.g., "go.mod" matches "go.mod" or "subdir/go.mod")
	fileName := filepath.Base(path)
	if fileName == pattern {
		return true
	}

	// Directory pattern (e.g., ".git/" matches ".git/config")
	if isDir {
		if strings.HasPrefix(path, pattern+"/") || strings.Contains(path, "/"+pattern+"/") {
			return true
		}
		// Also match the directory itself
		if path == pattern || strings.HasSuffix(path, "/"+pattern) {
			return true
		}
	}

	// Wildcard patterns
	if strings.Contains(pattern, "*") {
		matched, _ := filepath.Match(pattern, fileName)
		if matched {
			return true
		}
		// Try matching full path for patterns like "dir/*.go"
		matched, _ = filepath.Match(pattern, path)
		if matched {
			return true
		}
	}

	// Path prefix match (e.g., "vendor" matches "vendor/pkg/file.go")
	if strings.HasPrefix(path, pattern+"/") {
		return true
	}

	// Contains match for paths (e.g., "node_modules" matches "a/node_modules/b")
	if strings.Contains(path, "/"+pattern+"/") || strings.HasPrefix(path, pattern+"/") {
		return true
	}

	return false
}

// SetSLAPolicy sets SLA policy
func (a *App) SetSLAPolicy(policyJson string) error {
	var policy domain.SLAPolicy
	if err := json.Unmarshal([]byte(policyJson), &policy); err != nil {
		return fmt.Errorf("failed to parse SLA policy JSON: %w", err)
	}

	a.log.Info(fmt.Sprintf("SLA Policy set: %s", policy.Name))
	return nil
}

// GetSLAPolicy returns current SLA policy
func (a *App) GetSLAPolicy() (string, error) {
	defaultPolicy := domain.SLAPolicy{
		Name:        "standard",
		Description: "Standard SLA policy",
		MaxTokens:   10000,
		MaxFiles:    50,
		MaxTime:     300,
		MaxMemory:   1024 * 1024 * 100,
		MaxRetries:  3,
		Timeout:     30,
	}

	policyJson, err := json.Marshal(defaultPolicy)
	if err != nil {
		return "", fmt.Errorf("failed to marshal SLA policy: %w", err)
	}

	return string(policyJson), nil
}

// ============================================
// Secure API Key Storage
// ============================================

// SaveAPIKey saves an API key securely using encrypted storage
func (a *App) SaveAPIKey(provider, apiKey string) error {
	storage, err := a.getSecureStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize secure storage: %w", err)
	}

	key := formatAPIKeyName(provider)
	if err := storage.SaveCredential(key, apiKey); err != nil {
		return fmt.Errorf("failed to save API key for %s: %w", provider, err)
	}

	a.log.Info(fmt.Sprintf("API key saved securely for provider: %s", provider))
	return nil
}

// HasAPIKey checks if an API key exists for the given provider
func (a *App) HasAPIKey(provider string) bool {
	storage, err := a.getSecureStorage()
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to check API key: %v", err))
		return false
	}

	key := formatAPIKeyName(provider)
	return storage.HasCredential(key)
}

// DeleteAPIKey removes an API key for the given provider
func (a *App) DeleteAPIKey(provider string) error {
	storage, err := a.getSecureStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize secure storage: %w", err)
	}

	key := formatAPIKeyName(provider)
	if err := storage.DeleteCredential(key); err != nil {
		return fmt.Errorf("failed to delete API key for %s: %w", provider, err)
	}

	a.log.Info(fmt.Sprintf("API key deleted for provider: %s", provider))
	return nil
}

// GetAPIKeyStatus returns status of all API keys (without exposing actual keys)
func (a *App) GetAPIKeyStatus() (string, error) {
	storage, err := a.getSecureStorage()
	if err != nil {
		return "", fmt.Errorf("failed to initialize secure storage: %w", err)
	}

	providers := []string{"openai", "anthropic", "gemini", "openrouter", "ollama"}
	status := make(map[string]bool)

	for _, provider := range providers {
		key := formatAPIKeyName(provider)
		status[provider] = storage.HasCredential(key)
	}

	result, err := json.Marshal(status)
	if err != nil {
		return "", fmt.Errorf("failed to marshal API key status: %w", err)
	}

	return string(result), nil
}

// LoadAPIKey loads an API key for internal use (not exposed to frontend)
func (a *App) LoadAPIKey(provider string) (string, error) {
	storage, err := a.getSecureStorage()
	if err != nil {
		return "", fmt.Errorf("failed to initialize secure storage: %w", err)
	}

	key := formatAPIKeyName(provider)
	apiKey, err := storage.LoadCredential(key)
	if err != nil {
		return "", fmt.Errorf("failed to load API key for %s: %w", provider, err)
	}

	return apiKey, nil
}

// formatAPIKeyName formats the credential key name for a provider
func formatAPIKeyName(provider string) string {
	return fmt.Sprintf("api_key_%s", strings.ToLower(provider))
}

// getSecureStorage returns a SecureStorage instance for API key management
func (a *App) getSecureStorage() (*settingsfs.SecureStorage, error) {
	storage, err := settingsfs.NewSecureStorage()
	if err != nil {
		return nil, fmt.Errorf("failed to create secure storage: %w", err)
	}
	return storage, nil
}
