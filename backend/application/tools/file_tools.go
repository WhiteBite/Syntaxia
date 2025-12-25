package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"syntaxia/domain"
	"strings"
)

// FileToolsHandler handles file-related tools
type FileToolsHandler struct {
	BaseHandler
	FileReader domain.FileContentReader
}

// NewFileToolsHandler creates a new file tools handler
func NewFileToolsHandler(logger domain.Logger, fileReader domain.FileContentReader) *FileToolsHandler {
	return &FileToolsHandler{
		BaseHandler: NewBaseHandler(logger),
		FileReader:  fileReader,
	}
}

var fileToolNames = map[string]bool{
	"search_files":   true,
	"search_content": true,
	"read_file":      true,
	"list_directory": true,
	"get_file_info":  true,
	"list_functions": true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *FileToolsHandler) CanHandle(toolName string) bool {
	return fileToolNames[toolName]
}

// GetTools returns the list of file tools
func (h *FileToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "search_files",
			Description: "Search for files by name pattern (glob). Returns list of matching file paths.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"pattern":   {Type: "string", Description: "Glob pattern or partial filename to search for"},
					"directory": {Type: "string", Description: "Directory to search in (relative to project root)"},
				},
				Required: []string{"pattern"},
			},
		},
		{
			Name:        "search_content",
			Description: "Search for text/regex pattern in file contents. Returns matching lines with context.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"pattern":      {Type: "string", Description: "Text or regex pattern to search for"},
					"file_pattern": {Type: "string", Description: "Glob pattern to filter files"},
					"max_results":  {Type: "integer", Description: "Maximum number of results", Default: 20},
				},
				Required: []string{"pattern"},
			},
		},
		{
			Name:        "read_file",
			Description: "Read the contents of a file. Use this to examine code in detail.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path":       {Type: "string", Description: "Path to the file (relative to project root)"},
					"start_line": {Type: "integer", Description: "Starting line number (1-based)"},
					"end_line":   {Type: "integer", Description: "Ending line number (inclusive)"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "list_directory",
			Description: "List files and directories in a path. Use to explore project structure.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path":      {Type: "string", Description: "Directory path (relative to project root)"},
					"recursive": {Type: "boolean", Description: "Whether to list recursively", Default: false},
					"max_depth": {Type: "integer", Description: "Maximum depth for recursive listing", Default: 2},
				},
			},
		},
		{
			Name:        "get_file_info",
			Description: "Get metadata about a file (size, type, last modified).",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the file (relative to project root)"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "list_functions",
			Description: "List all functions/methods in a file. Works for Go, TypeScript, JavaScript.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the source file"},
				},
				Required: []string{"path"},
			},
		},
	}
}

// Execute executes a file tool
func (h *FileToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	switch toolName {
	case "search_files":
		return h.searchFiles(args, projectRoot)
	case "search_content":
		return h.searchContent(args, projectRoot)
	case "read_file":
		return h.readFile(args, projectRoot)
	case "list_directory":
		return h.listDirectory(args, projectRoot)
	case "get_file_info":
		return h.getFileInfo(args, projectRoot)
	case "list_functions":
		return h.listFunctions(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown file tool: %s", toolName)
	}
}

func (h *FileToolsHandler) searchFiles(args map[string]any, projectRoot string) (string, error) {
	pattern, _ := args["pattern"].(string)
	directory, _ := args["directory"].(string)

	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}

	searchDir := projectRoot
	if directory != "" {
		searchDir = filepath.Join(projectRoot, directory)
	}

	var matches []string
	patternLower := strings.ToLower(pattern)

	err := filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if name == "node_modules" || name == ".git" || name == "vendor" || name == "dist" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		nameLower := strings.ToLower(info.Name())

		if strings.Contains(nameLower, patternLower) {
			matches = append(matches, relPath)
		} else if matched, _ := filepath.Match(strings.ToLower(pattern), nameLower); matched {
			matches = append(matches, relPath)
		}

		if len(matches) >= 50 {
			return fmt.Errorf("limit reached")
		}
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return "", err
	}

	if len(matches) == 0 {
		return "No files found matching pattern: " + pattern, nil
	}

	return fmt.Sprintf("Found %d files:\n%s", len(matches), strings.Join(matches, "\n")), nil
}

func (h *FileToolsHandler) searchContent(args map[string]any, projectRoot string) (string, error) {
	pattern, _ := args["pattern"].(string)
	filePattern, _ := args["file_pattern"].(string)
	maxResults := 20
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}

	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		regex = regexp.MustCompile(regexp.QuoteMeta(pattern))
	}

	var results []string
	count := 0

	_ = filepath.Walk(projectRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return nil
		}
		if filePattern != "" {
			if matched, _ := filepath.Match(filePattern, info.Name()); !matched {
				return nil
			}
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		for i, line := range strings.Split(string(content), "\n") {
			if regex.MatchString(line) {
				results = append(results, fmt.Sprintf("%s:%d: %s", relPath, i+1, strings.TrimSpace(line)))
				count++
				if count >= maxResults {
					return fmt.Errorf("limit reached")
				}
			}
		}
		return nil
	})

	if len(results) == 0 {
		return "No matches found for: " + pattern, nil
	}
	return fmt.Sprintf("Found %d matches:\n%s", len(results), strings.Join(results, "\n")), nil
}

func (h *FileToolsHandler) readFile(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	fullPath := filepath.Join(projectRoot, path)

	// Security check
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}
	absProjectRoot = filepath.Clean(absProjectRoot)
	absFullPath = filepath.Clean(absFullPath)

	if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
		return "", fmt.Errorf("path traversal not allowed")
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	startLine := 1
	endLine := len(lines)

	if sl, ok := args["start_line"].(float64); ok && sl > 0 {
		startLine = int(sl)
	}
	if el, ok := args["end_line"].(float64); ok && el > 0 {
		endLine = int(el)
	}

	if startLine > len(lines) {
		return "", fmt.Errorf("start_line %d exceeds file length %d", startLine, len(lines))
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}

	var result []string
	for i := startLine - 1; i < endLine; i++ {
		result = append(result, fmt.Sprintf("%4d | %s", i+1, lines[i]))
	}

	header := fmt.Sprintf("=== %s (lines %d-%d of %d) ===\n", path, startLine, endLine, len(lines))
	return header + strings.Join(result, "\n"), nil
}

func (h *FileToolsHandler) listDirectory(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	recursive, _ := args["recursive"].(bool)
	maxDepth := 2
	if md, ok := args["max_depth"].(float64); ok {
		maxDepth = int(md)
	}

	targetDir := projectRoot
	if path != "" {
		targetDir = filepath.Join(projectRoot, path)
	}

	var entries []string
	if recursive {
		baseDepth := strings.Count(targetDir, string(os.PathSeparator))
		_ = filepath.Walk(targetDir, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			depth := strings.Count(p, string(os.PathSeparator)) - baseDepth
			if depth > maxDepth {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if relPath, _ := filepath.Rel(projectRoot, p); relPath == "." {
				return nil
			}
			if info.IsDir() && (info.Name() == "node_modules" || info.Name() == ".git") {
				return filepath.SkipDir
			}
			prefix := strings.Repeat("  ", depth)
			if info.IsDir() {
				entries = append(entries, fmt.Sprintf("%s📁 %s/", prefix, info.Name()))
			} else {
				entries = append(entries, fmt.Sprintf("%s📄 %s (%d bytes)", prefix, info.Name(), info.Size()))
			}
			return nil
		})
	} else {
		files, err := os.ReadDir(targetDir)
		if err != nil {
			return "", err
		}
		for _, f := range files {
			info, _ := f.Info()
			if f.IsDir() {
				entries = append(entries, fmt.Sprintf("📁 %s/", f.Name()))
			} else if info != nil {
				entries = append(entries, fmt.Sprintf("📄 %s (%d bytes)", f.Name(), info.Size()))
			}
		}
	}

	if len(entries) == 0 {
		return "Directory is empty", nil
	}
	return strings.Join(entries, "\n"), nil
}

func (h *FileToolsHandler) getFileInfo(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	fullPath := filepath.Join(projectRoot, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return "", err
	}

	result := map[string]any{
		"name":     info.Name(),
		"path":     path,
		"size":     info.Size(),
		"isDir":    info.IsDir(),
		"modified": info.ModTime().Format("2006-01-02 15:04:05"),
		"ext":      filepath.Ext(path),
	}

	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
	return string(jsonBytes), nil
}

func (h *FileToolsHandler) listFunctions(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}

	fullPath := filepath.Join(projectRoot, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(path))
	var functions []string

	switch ext {
	case ".go":
		re := regexp.MustCompile(`func\s+(\([^)]+\)\s+)?(\w+)\s*\(`)
		for _, m := range re.FindAllStringSubmatch(string(content), -1) {
			if len(m) >= 3 {
				functions = append(functions, m[2])
			}
		}
	case ".ts", ".js", ".tsx", ".jsx":
		patterns := []string{
			`function\s+(\w+)\s*\(`,
			`(\w+)\s*=\s*(?:async\s+)?function`,
			`(\w+)\s*=\s*(?:async\s+)?\([^)]*\)\s*=>`,
		}
		for _, p := range patterns {
			re := regexp.MustCompile(p)
			for _, m := range re.FindAllStringSubmatch(string(content), -1) {
				if len(m) >= 2 && m[1] != "" {
					functions = append(functions, m[1])
				}
			}
		}
	default:
		return "Unsupported file type for function extraction: " + ext, nil
	}

	if len(functions) == 0 {
		return "No functions found in " + path, nil
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, f := range functions {
		if !seen[f] {
			seen[f] = true
			unique = append(unique, f)
		}
	}

	return fmt.Sprintf("Functions in %s:\n- %s", path, strings.Join(unique, "\n- ")), nil
}
