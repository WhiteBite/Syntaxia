package context

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

// applyContentOptimizations applies all content optimizations based on options
func (s *Service) applyContentOptimizations(content, filePath string, options *BuildOptions) string {
	if options == nil {
		return content
	}

	// 1. Strip comments (existing functionality)
	if options.StripComments {
		content = s.stripComments(content, filePath)
	}

	// 2. Strip license headers
	if options.StripLicense {
		content = s.stripLicenseHeader(content, filePath)
	}

	// 3. Compact data files (JSON/YAML)
	if options.CompactDataFiles {
		content = s.compactDataFile(content, filePath)
	}

	// 4. Trim trailing whitespace
	if options.TrimWhitespace {
		content = s.trimTrailingWhitespace(content)
	}

	// 5. Collapse empty lines (last, after all removals)
	if options.CollapseEmptyLines {
		content = s.collapseEmptyLines(content)
	}

	return content
}

// stripComments removes comments from code based on file extension
func (s *Service) stripComments(content, filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs":
		return s.stripCStyleComments(content)
	case ".py", ".sh":
		return s.stripHashComments(content)
	case ".html", ".xml":
		return s.stripXMLComments(content)
	default:
		return content
	}
}

func (s *Service) stripCStyleComments(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	inBlockComment := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Handle block comments
		if inBlockComment {
			if strings.Contains(line, "*/") {
				inBlockComment = false
				// Keep part after */
				parts := strings.SplitN(line, "*/", 2)
				if len(parts) > 1 {
					line = parts[1]
				} else {
					continue
				}
			} else {
				continue
			}
		}

		// Handle start of block comments
		if strings.Contains(line, "/*") {
			inBlockComment = true
			parts := strings.SplitN(line, "/*", 2)
			line = parts[0]
			if strings.TrimSpace(line) == "" {
				continue
			}
		}

		// Handle line comments
		if strings.HasPrefix(trimmed, "//") {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (s *Service) stripHashComments(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			continue
		}
		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

func (s *Service) stripXMLComments(content string) string {
	result := content
	for {
		start := strings.Index(result, "<!--")
		if start == -1 {
			break
		}
		end := strings.Index(result[start:], "-->")
		if end == -1 {
			break
		}
		result = result[:start] + result[start+end+3:]
	}
	return result
}

// stripLicenseHeader removes license/copyright headers from file content
func (s *Service) stripLicenseHeader(content, filePath string) string {
	if len(content) == 0 {
		return content
	}

	licenseKeywords := []string{
		"copyright", "license", "licensed", "spdx-license",
		"mit license", "apache license", "bsd license",
		"all rights reserved", "permission is hereby granted",
	}

	trimmed := strings.TrimLeft(content, " \t\n\r")
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check for block comment at start
	var blockStart, blockEnd string
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".swift", ".kt", ".rs":
		blockStart, blockEnd = "/*", "*/"
	case ".html", ".xml", ".vue", ".svelte":
		blockStart, blockEnd = "<!--", "-->"
	}

	if blockStart != "" && strings.HasPrefix(trimmed, blockStart) {
		endIdx := strings.Index(trimmed, blockEnd)
		if endIdx != -1 {
			commentBlock := strings.ToLower(trimmed[:endIdx+len(blockEnd)])
			for _, keyword := range licenseKeywords {
				if strings.Contains(commentBlock, keyword) {
					afterComment := trimmed[endIdx+len(blockEnd):]
					return strings.TrimLeft(afterComment, " \t\n\r")
				}
			}
		}
	}

	// Check for line comments at start
	var lineComment string
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".swift", ".kt", ".rs":
		lineComment = "//"
	case ".py", ".rb", ".sh", ".yaml", ".yml":
		lineComment = "#"
	case ".sql":
		lineComment = "--"
	}

	if lineComment != "" && strings.HasPrefix(trimmed, lineComment) {
		lines := strings.Split(trimmed, "\n")
		var commentLines []string
		lastCommentLine := 0

		for i, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, lineComment) || trimmedLine == "" {
				commentLines = append(commentLines, line)
				if strings.HasPrefix(trimmedLine, lineComment) {
					lastCommentLine = i
				}
			} else {
				break
			}
		}

		if len(commentLines) > 0 {
			commentBlock := strings.ToLower(strings.Join(commentLines[:lastCommentLine+1], "\n"))
			for _, keyword := range licenseKeywords {
				if strings.Contains(commentBlock, keyword) {
					remaining := lines[lastCommentLine+1:]
					return strings.TrimLeft(strings.Join(remaining, "\n"), "\n")
				}
			}
		}
	}

	return content
}

// compactDataFile compacts JSON/YAML files
func (s *Service) compactDataFile(content, filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return s.compactJSON(content)
	case ".yaml", ".yml":
		return s.compactYAML(content)
	default:
		return content
	}
}

func (s *Service) compactJSON(content string) string {
	if !strings.Contains(content, "\n") {
		return content // Already compact
	}

	var data interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return content // Invalid JSON, return as-is
	}

	result, err := json.Marshal(data)
	if err != nil {
		return content
	}
	return string(result)
}

func (s *Service) compactYAML(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	prevEmpty := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines (keep max one)
		if trimmed == "" {
			if !prevEmpty {
				result = append(result, "")
				prevEmpty = true
			}
			continue
		}
		prevEmpty = false

		// Skip comment lines
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Remove inline comments (simple heuristic)
		if idx := strings.Index(line, " #"); idx != -1 {
			line = strings.TrimRight(line[:idx], " \t")
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// trimTrailingWhitespace removes trailing whitespace from each line
func (s *Service) trimTrailingWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}

// collapseEmptyLines collapses multiple empty lines into maximum two
func (s *Service) collapseEmptyLines(content string) string {
	if !strings.Contains(content, "\n\n\n") {
		return content // Fast path: no triple newlines
	}

	var result strings.Builder
	result.Grow(len(content))

	emptyCount := 0
	lineStart := 0

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			line := content[lineStart:i]
			isEmpty := strings.TrimSpace(line) == ""

			if isEmpty {
				emptyCount++
				if emptyCount <= 2 {
					result.WriteString(line)
					result.WriteByte('\n')
				}
			} else {
				emptyCount = 0
				result.WriteString(line)
				result.WriteByte('\n')
			}
			lineStart = i + 1
		}
	}

	// Handle last line
	if lineStart < len(content) {
		result.WriteString(content[lineStart:])
	}

	return result.String()
}

// filterTestFiles filters out test files from the list
func (s *Service) filterTestFiles(paths []string) []string {
	if len(paths) == 0 {
		return paths
	}

	result := make([]string, 0, len(paths))
	for _, path := range paths {
		if !s.isTestFile(path) {
			result = append(result, path)
		}
	}

	if len(result) < len(paths) {
		s.logger.Info(fmt.Sprintf("Filtered out %d test files", len(paths)-len(result)))
	}

	return result
}

// isTestFile checks if a file is a test file
func (s *Service) isTestFile(filePath string) bool {
	// Normalize path
	filePath = filepath.ToSlash(filePath)
	fileName := filepath.Base(filePath)
	lower := strings.ToLower(fileName)

	// Check test directories (both with and without leading slash)
	testDirs := []string{"test/", "tests/", "__tests__/", "spec/", "e2e/", "__mocks__/"}
	for _, dir := range testDirs {
		if strings.Contains(filePath, "/"+dir) || strings.HasPrefix(filePath, dir) {
			return true
		}
	}

	// Check test file patterns
	testPatterns := []string{
		"_test.go", ".test.js", ".test.ts", ".test.tsx",
		".spec.js", ".spec.ts", ".spec.tsx",
		"_test.py", "test_", "Test.java", "Tests.cs", "_spec.rb",
		".stories.tsx", ".stories.ts", ".stories.js",
	}

	for _, pattern := range testPatterns {
		if strings.Contains(lower, strings.ToLower(pattern)) {
			return true
		}
	}

	return strings.HasPrefix(lower, "test_")
}
