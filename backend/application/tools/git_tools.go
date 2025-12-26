package tools

import (
	"fmt"
	"os/exec"
	"strings"

	"syntaxia/domain"
	"syntaxia/internal/executil"
)

const (
	// minLimit is the minimum allowed value for limit parameters
	minLimit = 1
	// maxLimit is the maximum allowed value for limit parameters
	maxLimit = 100
	// defaultLimit is the default value for limit parameters
	defaultLimit = 10
	// maxDiffSize is the maximum size of diff output before truncation
	maxDiffSize = 5000
)

// GitToolsHandler handles git-related tools
type GitToolsHandler struct {
	BaseHandler
	GitContext domain.GitContextBuilder
}

// NewGitToolsHandler creates a new git tools handler
func NewGitToolsHandler(logger domain.Logger, gitContext domain.GitContextBuilder) *GitToolsHandler {
	return &GitToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		GitContext:  gitContext,
	}
}

var gitToolNames = map[string]bool{
	"git_status":          true,
	"git_diff":            true,
	"git_log":             true,
	"git_changed_files":   true,
	"git_co_changed":      true,
	"git_suggest_context": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *GitToolsHandler) CanHandle(toolName string) bool {
	return gitToolNames[toolName]
}

// validateArgs validates that all required arguments are present and not nil
func (h *GitToolsHandler) validateArgs(args map[string]any, required []string) error {
	for _, key := range required {
		val, ok := args[key]
		if !ok {
			return fmt.Errorf("missing required argument: %s", key)
		}
		if val == nil {
			return fmt.Errorf("argument %s cannot be nil", key)
		}
	}
	return nil
}

// getStringArg safely extracts a string argument with a default value
func (h *GitToolsHandler) getStringArg(args map[string]any, key string, defaultVal string) string {
	if val, ok := args[key].(string); ok {
		return val
	}
	return defaultVal
}

// getIntArg safely extracts an integer argument with a default value
// Handles both float64 (JSON numbers) and int types
func (h *GitToolsHandler) getIntArg(args map[string]any, key string, defaultVal int) int {
	if val, ok := args[key].(float64); ok {
		return int(val)
	}
	if val, ok := args[key].(int); ok {
		return val
	}
	return defaultVal
}

// getBoolArg safely extracts a boolean argument with a default value
func (h *GitToolsHandler) getBoolArg(args map[string]any, key string, defaultVal bool) bool {
	if val, ok := args[key].(bool); ok {
		return val
	}
	return defaultVal
}

// getStringArrayArg safely extracts a string array argument
func (h *GitToolsHandler) getStringArrayArg(args map[string]any, key string) ([]string, error) {
	val, ok := args[key]
	if !ok {
		return nil, nil
	}
	if val == nil {
		return nil, nil
	}

	arr, ok := val.([]any)
	if !ok {
		return nil, fmt.Errorf("argument %s must be an array", key)
	}

	result := make([]string, 0, len(arr))
	for i, item := range arr {
		str, ok := item.(string)
		if !ok {
			return nil, fmt.Errorf("argument %s[%d] must be a string", key, i)
		}
		result = append(result, str)
	}
	return result, nil
}

// GetTools returns the list of git tools
func (h *GitToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "git_status",
			Description: "Get git status - list of modified, added, deleted files.",
			Parameters:  domain.ToolParameters{Type: "object", Properties: map[string]domain.ToolProperty{}},
		},
		{
			Name:        "git_diff",
			Description: "Get git diff for a file or all changes.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path":   {Type: "string", Description: "Path to file (optional, empty for all changes)"},
					"staged": {Type: "boolean", Description: "Show staged changes only", Default: false},
				},
			},
		},
		{
			Name:        "git_log",
			Description: "Get recent git commits.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"limit": {Type: "integer", Description: "Number of commits to show", Default: 10},
					"path":  {Type: "string", Description: "Filter by file path (optional)"},
				},
			},
		},
		{
			Name:        "git_changed_files",
			Description: "Get recently changed files from git history",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"since":       {Type: "string", Description: "Time period (e.g., '1 week ago', '2024-01-01')"},
					"path_filter": {Type: "string", Description: "Filter by path pattern"},
				},
			},
		},
		{
			Name:        "git_co_changed",
			Description: "Get files that are often changed together with the specified file",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"file_path": {Type: "string", Description: "Path to the file"},
					"limit":     {Type: "integer", Description: "Maximum results (default: 10)"},
				},
				Required: []string{"file_path"},
			},
		},
		{
			Name:        "git_suggest_context",
			Description: "Suggest files to include in context based on git history and task description",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"task":          {Type: "string", Description: "Task description"},
					"current_files": {Type: "array", Description: "Currently selected files"},
					"limit":         {Type: "integer", Description: "Maximum suggestions (default: 10)"},
				},
			},
		},
	}
}

// Execute executes a git tool
func (h *GitToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	// Validate args is not nil
	if args == nil {
		args = make(map[string]any)
	}

	switch toolName {
	case "git_status":
		return h.gitStatus(projectRoot)
	case "git_diff":
		return h.gitDiff(args, projectRoot)
	case "git_log":
		return h.gitLog(args, projectRoot)
	case "git_changed_files":
		return h.gitChangedFiles(args, projectRoot)
	case "git_co_changed":
		return h.gitCoChanged(args, projectRoot)
	case "git_suggest_context":
		return h.gitSuggestContext(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown git tool: %s", toolName)
	}
}

func (h *GitToolsHandler) gitStatus(projectRoot string) (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git status failed: %w", err)
	}

	if len(output) == 0 {
		return "Working directory clean - no changes", nil
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		status := line[:2]
		file := strings.TrimSpace(line[2:])

		var statusText string
		switch strings.TrimSpace(status) {
		case "M":
			statusText = "modified"
		case "A":
			statusText = "added"
		case "D":
			statusText = "deleted"
		case "??":
			statusText = "untracked"
		case "MM":
			statusText = "modified (staged + unstaged)"
		default:
			statusText = status
		}
		result = append(result, fmt.Sprintf("%s: %s", statusText, file))
	}

	return fmt.Sprintf("Git status (%d files):\n%s", len(result), strings.Join(result, "\n")), nil
}

func (h *GitToolsHandler) gitDiff(args map[string]any, projectRoot string) (string, error) {
	path := h.getStringArg(args, "path", "")
	staged := h.getBoolArg(args, "staged", false)

	cmdArgs := []string{"diff"}
	if staged {
		cmdArgs = append(cmdArgs, "--staged")
	}
	if path != "" {
		cmdArgs = append(cmdArgs, "--", path)
	}

	cmd := exec.Command("git", cmdArgs...)
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}

	if len(output) == 0 {
		return "No differences found", nil
	}

	result := string(output)
	if len(result) > maxDiffSize {
		result = result[:maxDiffSize] + "\n... (truncated)"
	}

	return result, nil
}

func (h *GitToolsHandler) gitLog(args map[string]any, projectRoot string) (string, error) {
	limit := h.getIntArg(args, "limit", defaultLimit)

	// Validate limit range
	if limit < minLimit {
		limit = minLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	path := h.getStringArg(args, "path", "")

	cmdArgs := []string{"log", fmt.Sprintf("-n%d", limit), "--oneline", "--format=%h %s (%an, %ar)"}
	if path != "" {
		cmdArgs = append(cmdArgs, "--", path)
	}

	cmd := exec.Command("git", cmdArgs...)
	executil.HideWindow(cmd)
	cmd.Dir = projectRoot
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git log failed: %w", err)
	}

	if len(output) == 0 {
		return "No commits found", nil
	}

	return fmt.Sprintf("Recent commits:\n%s", strings.TrimSpace(string(output))), nil
}

func (h *GitToolsHandler) gitChangedFiles(args map[string]any, projectRoot string) (string, error) {
	if h.GitContext == nil {
		return "", fmt.Errorf("git context not initialized")
	}

	since := h.getStringArg(args, "since", "")
	pathFilter := h.getStringArg(args, "path_filter", "")

	changes, err := h.GitContext.GetRecentChanges(since, pathFilter)
	if err != nil {
		return "", fmt.Errorf("failed to get recent changes: %w", err)
	}

	if len(changes) == 0 {
		return "No recent changes found", nil
	}

	var result strings.Builder
	result.WriteString("Recently changed files:\n\n")
	for _, c := range changes {
		result.WriteString(fmt.Sprintf("  - %s (%d changes, last: %s)\n", c.FilePath, c.ChangeCount, c.LastChanged.Format("2006-01-02")))
	}
	return result.String(), nil
}

func (h *GitToolsHandler) gitCoChanged(args map[string]any, projectRoot string) (string, error) {
	if h.GitContext == nil {
		return "", fmt.Errorf("git context not initialized")
	}

	// Validate required argument
	if err := h.validateArgs(args, []string{"file_path"}); err != nil {
		return "", err
	}

	filePath := h.getStringArg(args, "file_path", "")
	if filePath == "" {
		return "", fmt.Errorf("file_path cannot be empty")
	}

	limit := h.getIntArg(args, "limit", defaultLimit)
	if limit < minLimit {
		limit = minLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	coChanged, err := h.GitContext.GetCoChangedFiles(filePath, limit)
	if err != nil {
		return "", fmt.Errorf("failed to get co-changed files: %w", err)
	}

	if len(coChanged) == 0 {
		return fmt.Sprintf("No co-changed files found for %s", filePath), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Files often changed with %s:\n\n", filePath))
	for i, f := range coChanged {
		result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, f))
	}
	return result.String(), nil
}

func (h *GitToolsHandler) gitSuggestContext(args map[string]any, projectRoot string) (string, error) {
	if h.GitContext == nil {
		return "", fmt.Errorf("git context not initialized")
	}

	task := h.getStringArg(args, "task", "")

	currentFiles, err := h.getStringArrayArg(args, "current_files")
	if err != nil {
		return "", err
	}

	limit := h.getIntArg(args, "limit", defaultLimit)
	if limit < minLimit {
		limit = minLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}

	suggestions, err := h.GitContext.SuggestContextFiles(task, currentFiles, limit)
	if err != nil {
		return "", fmt.Errorf("failed to suggest context files: %w", err)
	}

	if len(suggestions) == 0 {
		return "No context suggestions found", nil
	}

	var result strings.Builder
	result.WriteString("Suggested files for context:\n\n")
	for i, f := range suggestions {
		result.WriteString(fmt.Sprintf("  %d. %s\n", i+1, f))
	}
	return result.String(), nil
}
