package textutils

import (
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

type commentStripperImpl struct {
	log domain.Logger
}

// NewCommentStripper creates a new implementation of domain.CommentStripper
func NewCommentStripper(log domain.Logger) domain.CommentStripper {
	return &commentStripperImpl{log: log}
}

// Strip removes comments from code content based on file extension
func (c *commentStripperImpl) Strip(content, filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs":
		return stripCStyleComments(content)
	case ".py", ".sh":
		return stripHashComments(content)
	case ".html", ".xml":
		return stripXMLComments(content)
	default:
		return content
	}
}

// stripCStyleComments removes C-style comments (// and /* */)
func stripCStyleComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inBlockComment := false

	for _, line := range lines {
		line, inBlockComment = processCommentLine(line, inBlockComment)
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, line)
		}
	}
	return strings.Join(result, "\n")
}

// processCommentLine processes a single line for C-style comments
func processCommentLine(line string, inBlock bool) (string, bool) {
	if inBlock {
		return handleBlockCommentContinuation(line)
	}
	line = handleBlockCommentStart(line, &inBlock)
	line = stripLineComment(line)
	line = stripInlineBlockComments(line)
	return line, inBlock
}

// handleBlockCommentContinuation handles lines inside a block comment
func handleBlockCommentContinuation(line string) (string, bool) {
	if idx := strings.Index(line, "*/"); idx != -1 {
		return line[idx+2:], false
	}
	return "", true
}

// handleBlockCommentStart handles start of block comments
func handleBlockCommentStart(line string, inBlock *bool) string {
	if strings.Contains(line, "/*") && !strings.Contains(line, "*/") {
		parts := strings.SplitN(line, "/*", 2)
		*inBlock = true
		return parts[0]
	}
	return line
}

// stripLineComment removes // comments
func stripLineComment(line string) string {
	if idx := strings.Index(line, "//"); idx != -1 {
		return line[:idx]
	}
	return line
}

// stripInlineBlockComments removes inline /* */ comments
func stripInlineBlockComments(line string) string {
	for strings.Contains(line, "/*") && strings.Contains(line, "*/") {
		start, end := strings.Index(line, "/*"), strings.Index(line, "*/")
		if end > start {
			line = line[:start] + line[end+2:]
		} else {
			break
		}
	}
	return line
}

// stripHashComments removes hash-style comments (#)
func stripHashComments(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	for _, line := range lines {
		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// stripXMLComments removes XML-style comments (<!-- -->)
func stripXMLComments(content string) string {
	for strings.Contains(content, "<!--") && strings.Contains(content, "-->") {
		start := strings.Index(content, "<!--")
		end := strings.Index(content, "-->")
		if end > start {
			content = content[:start] + content[end+3:]
		} else {
			// Malformed comment, break to avoid infinite loop
			break
		}
	}
	return content
}
