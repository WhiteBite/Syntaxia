package context

import (
	"fmt"
	"path/filepath"
	"strings"
)

// OutputFormat defines the format for context output
type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown" // ## File: path\n```lang\ncontent\n```
	FormatXML      OutputFormat = "xml"      // <file path="..."><content>...</content></file>
	FormatPlain    OutputFormat = "plain"    // --- File: path ---\ncontent
)

// formatFileHeader returns the header for a file based on output format
func formatFileHeader(filePath string, format OutputFormat) string {
	switch format {
	case FormatXML:
		return fmt.Sprintf("<file path=\"%s\">\n<content>\n", filePath)
	case FormatPlain:
		return fmt.Sprintf("--- File: %s ---\n", filePath)
	default: // FormatMarkdown
		ext := filepath.Ext(filePath)
		lang := ""
		if len(ext) > 1 {
			lang = ext[1:]
		}
		return fmt.Sprintf("## File: %s\n\n```%s\n", filePath, lang)
	}
}

// formatFileFooter returns the footer for a file based on output format
func formatFileFooter(format OutputFormat) string {
	switch format {
	case FormatXML:
		return "\n</content>\n</file>\n\n"
	case FormatPlain:
		return "\n\n"
	default: // FormatMarkdown
		return "\n```\n\n"
	}
}

// escapeForFormat escapes content based on output format
func escapeForFormat(content string, format OutputFormat) string {
	switch format {
	case FormatXML:
		// Escape XML special characters
		content = strings.ReplaceAll(content, "&", "&amp;")
		content = strings.ReplaceAll(content, "<", "&lt;")
		content = strings.ReplaceAll(content, ">", "&gt;")
		return content
	default:
		return content
	}
}

// addLineNumbers adds line numbers to content
func addLineNumbers(content string) string {
	lines := strings.Split(content, "\n")
	width := len(fmt.Sprintf("%d", len(lines)))

	var result strings.Builder
	result.Grow(len(content) + len(lines)*(width+3)) // pre-allocate

	for i, line := range lines {
		if i > 0 {
			result.WriteByte('\n')
		}
		result.WriteString(fmt.Sprintf("%*d | %s", width, i+1, line))
	}
	return result.String()
}

// Service method wrappers for backward compatibility
func (s *Service) formatFileHeader(filePath string, format OutputFormat) string {
	return formatFileHeader(filePath, format)
}

func (s *Service) formatFileFooter(format OutputFormat) string {
	return formatFileFooter(format)
}

func (s *Service) escapeForFormat(content string, format OutputFormat) string {
	return escapeForFormat(content, format)
}
