// Package ai provides context tools adapters for agentic context gathering.
package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// ContextFileTools implements FileToolsProvider for context gathering
type ContextFileTools struct {
	log domain.Logger
}

// NewContextFileTools creates a new ContextFileTools
func NewContextFileTools(log domain.Logger) *ContextFileTools {
	return &ContextFileTools{log: log}
}

// SearchFiles searches for files by pattern
func (t *ContextFileTools) SearchFiles(pattern, directory, projectRoot string) ([]string, error) {
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
			if name == "node_modules" || name == ".git" || name == "vendor" || name == "dist" || name == "__pycache__" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		nameLower := strings.ToLower(info.Name())

		// Check partial match
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
		return nil, err
	}

	return matches, nil
}

// SearchContent searches file contents for a pattern
func (t *ContextFileTools) SearchContent(pattern, filePattern, projectRoot string, maxResults int) ([]ContentMatch, error) {
	if maxResults <= 0 {
		maxResults = 20
	}

	regex, err := regexp.Compile("(?i)" + pattern)
	if err != nil {
		// Fall back to literal match
		regex = regexp.MustCompile(regexp.QuoteMeta(pattern))
	}

	var results []ContentMatch

	err = filepath.Walk(projectRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}

		// Skip common non-code directories
		if strings.Contains(path, "node_modules") ||
			strings.Contains(path, ".git") ||
			strings.Contains(path, "vendor") ||
			strings.Contains(path, "__pycache__") {
			return nil
		}

		// Apply file pattern filter
		if filePattern != "" {
			if matched, _ := filepath.Match(filePattern, info.Name()); !matched {
				return nil
			}
		}

		// Skip binary files
		if isBinaryExtension(filepath.Ext(path)) {
			return nil
		}

		// Skip large files
		if info.Size() > 1024*1024 { // 1MB
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		lines := strings.Split(string(content), "\n")

		for i, line := range lines {
			if regex.MatchString(line) {
				results = append(results, ContentMatch{
					FilePath: relPath,
					Line:     i + 1,
					Content:  strings.TrimSpace(line),
				})
				if len(results) >= maxResults {
					return fmt.Errorf("limit reached")
				}
			}
		}
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, err
	}

	return results, nil
}

// ReadFile reads a file's content
func (t *ContextFileTools) ReadFile(path, projectRoot string) (string, error) {
	fullPath := filepath.Join(projectRoot, path)

	// Security check - prevent path traversal
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}

	if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
		return "", fmt.Errorf("path traversal not allowed")
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

func isBinaryExtension(ext string) bool {
	binaryExts := map[string]bool{
		".exe": true, ".dll": true, ".so": true, ".dylib": true,
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".zip": true, ".tar": true, ".gz": true, ".rar": true, ".7z": true,
		".mp3": true, ".mp4": true, ".avi": true, ".mov": true, ".wav": true,
		".woff": true, ".woff2": true, ".ttf": true, ".eot": true,
		".pyc": true, ".pyo": true, ".class": true,
		".o": true, ".a": true, ".lib": true,
	}
	return binaryExts[strings.ToLower(ext)]
}

// ContextSymbolTools implements SymbolToolsProvider for context gathering
type ContextSymbolTools struct {
	registry analysis.AnalyzerRegistry
	log      domain.Logger
}

// NewContextSymbolTools creates a new ContextSymbolTools
func NewContextSymbolTools(registry analysis.AnalyzerRegistry, log domain.Logger) *ContextSymbolTools {
	return &ContextSymbolTools{
		registry: registry,
		log:      log,
	}
}

// ListSymbols lists symbols in a file
func (t *ContextSymbolTools) ListSymbols(path, projectRoot string) ([]SymbolInfo, error) {
	fullPath := filepath.Join(projectRoot, path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	analyzer := t.registry.GetAnalyzer(path)
	if analyzer == nil {
		// Return empty list for unsupported file types
		return []SymbolInfo{}, nil
	}

	symbols, err := analyzer.ExtractSymbols(context.Background(), path, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract symbols: %w", err)
	}

	result := make([]SymbolInfo, 0, len(symbols))
	for _, s := range symbols {
		result = append(result, SymbolInfo{
			Name:      s.Name,
			Kind:      string(s.Kind),
			FilePath:  path,
			Line:      s.StartLine,
			Signature: s.Signature,
		})
	}

	return result, nil
}

// SearchSymbols searches for symbols across the project
func (t *ContextSymbolTools) SearchSymbols(query, kind, projectRoot string) ([]SymbolInfo, error) {
	var results []SymbolInfo
	queryLower := strings.ToLower(query)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}

		// Skip common non-code directories
		if strings.Contains(path, "node_modules") ||
			strings.Contains(path, ".git") ||
			strings.Contains(path, "vendor") {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		analyzer := t.registry.GetAnalyzer(relPath)
		if analyzer == nil {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		symbols, err := analyzer.ExtractSymbols(context.Background(), relPath, content)
		if err != nil {
			return nil
		}

		for _, s := range symbols {
			// Filter by name
			if !strings.Contains(strings.ToLower(s.Name), queryLower) {
				continue
			}
			// Filter by kind if specified
			if kind != "" && !strings.EqualFold(string(s.Kind), kind) {
				continue
			}

			results = append(results, SymbolInfo{
				Name:      s.Name,
				Kind:      string(s.Kind),
				FilePath:  relPath,
				Line:      s.StartLine,
				Signature: s.Signature,
			})

			if len(results) >= 30 {
				return fmt.Errorf("limit reached")
			}
		}
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, err
	}

	return results, nil
}

// SimpleSymbolTools provides basic symbol extraction without analyzer registry
// Useful for testing or when full analyzer infrastructure is not available
type SimpleSymbolTools struct {
	log domain.Logger
}

// NewSimpleSymbolTools creates a new SimpleSymbolTools
func NewSimpleSymbolTools(log domain.Logger) *SimpleSymbolTools {
	return &SimpleSymbolTools{log: log}
}

// ListSymbols extracts symbols using regex patterns
func (t *SimpleSymbolTools) ListSymbols(path, projectRoot string) ([]SymbolInfo, error) {
	fullPath := filepath.Join(projectRoot, path)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	ext := strings.ToLower(filepath.Ext(path))
	var symbols []SymbolInfo

	switch ext {
	case ".go":
		symbols = extractGoSymbols(string(content), path)
	case ".ts", ".tsx", ".js", ".jsx":
		symbols = extractTSSymbols(string(content), path)
	case ".py":
		symbols = extractPythonSymbols(string(content), path)
	default:
		return []SymbolInfo{}, nil
	}

	return symbols, nil
}

// SearchSymbols searches for symbols (simplified implementation)
func (t *SimpleSymbolTools) SearchSymbols(query, kind, projectRoot string) ([]SymbolInfo, error) {
	var results []SymbolInfo
	queryLower := strings.ToLower(query)

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() {
			return nil
		}

		// Skip common non-code directories
		if strings.Contains(path, "node_modules") ||
			strings.Contains(path, ".git") ||
			strings.Contains(path, "vendor") {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		symbols, err := t.ListSymbols(relPath, projectRoot)
		if err != nil {
			return nil
		}

		for _, s := range symbols {
			if !strings.Contains(strings.ToLower(s.Name), queryLower) {
				continue
			}
			if kind != "" && !strings.EqualFold(s.Kind, kind) {
				continue
			}

			results = append(results, s)
			if len(results) >= 30 {
				return fmt.Errorf("limit reached")
			}
		}
		return nil
	})

	if err != nil && err.Error() != "limit reached" {
		return nil, err
	}

	return results, nil
}

func extractGoSymbols(content, path string) []SymbolInfo {
	var symbols []SymbolInfo
	lines := strings.Split(content, "\n")

	funcRe := regexp.MustCompile(`^func\s+(\([^)]+\)\s+)?(\w+)\s*\(`)
	typeRe := regexp.MustCompile(`^type\s+(\w+)\s+(struct|interface)`)
	constRe := regexp.MustCompile(`^const\s+(\w+)`)
	varRe := regexp.MustCompile(`^var\s+(\w+)`)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if m := funcRe.FindStringSubmatch(trimmed); len(m) >= 3 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[2],
				Kind:     "function",
				FilePath: path,
				Line:     i + 1,
			})
		} else if m := typeRe.FindStringSubmatch(trimmed); len(m) >= 3 {
			kind := "type"
			if m[2] == "interface" {
				kind = "interface"
			} else if m[2] == "struct" {
				kind = "class"
			}
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     kind,
				FilePath: path,
				Line:     i + 1,
			})
		} else if m := constRe.FindStringSubmatch(trimmed); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "constant",
				FilePath: path,
				Line:     i + 1,
			})
		} else if m := varRe.FindStringSubmatch(trimmed); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "variable",
				FilePath: path,
				Line:     i + 1,
			})
		}
	}

	return symbols
}

func extractTSSymbols(content, path string) []SymbolInfo {
	var symbols []SymbolInfo
	lines := strings.Split(content, "\n")

	funcRe := regexp.MustCompile(`(?:export\s+)?(?:async\s+)?function\s+(\w+)`)
	classRe := regexp.MustCompile(`(?:export\s+)?class\s+(\w+)`)
	interfaceRe := regexp.MustCompile(`(?:export\s+)?interface\s+(\w+)`)
	typeRe := regexp.MustCompile(`(?:export\s+)?type\s+(\w+)`)
	constRe := regexp.MustCompile(`(?:export\s+)?const\s+(\w+)`)

	for i, line := range lines {
		if m := funcRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "function",
				FilePath: path,
				Line:     i + 1,
			})
		}
		if m := classRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "class",
				FilePath: path,
				Line:     i + 1,
			})
		}
		if m := interfaceRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "interface",
				FilePath: path,
				Line:     i + 1,
			})
		}
		if m := typeRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "type",
				FilePath: path,
				Line:     i + 1,
			})
		}
		if m := constRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "constant",
				FilePath: path,
				Line:     i + 1,
			})
		}
	}

	return symbols
}

func extractPythonSymbols(content, path string) []SymbolInfo {
	var symbols []SymbolInfo
	lines := strings.Split(content, "\n")

	funcRe := regexp.MustCompile(`^(?:async\s+)?def\s+(\w+)`)
	classRe := regexp.MustCompile(`^class\s+(\w+)`)

	for i, line := range lines {
		if m := funcRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "function",
				FilePath: path,
				Line:     i + 1,
			})
		}
		if m := classRe.FindStringSubmatch(line); len(m) >= 2 {
			symbols = append(symbols, SymbolInfo{
				Name:     m[1],
				Kind:     "class",
				FilePath: path,
				Line:     i + 1,
			})
		}
	}

	return symbols
}
