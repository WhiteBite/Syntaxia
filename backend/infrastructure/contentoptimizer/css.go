package contentoptimizer

import (
	"regexp"
	"strings"
	"syntaxia/domain"
)

// optimizeCSS applies CSS-specific optimizations
func (o *ContentOptimizer) optimizeCSS(content string, opts domain.ContentOptimizeOptions) string {
	// Remove comments
	if opts.StripComments {
		content = stripCSSComments(content)
	}

	// Collapse @import statements
	if opts.StripComments {
		content = collapseCSSImports(content)
	}

	// Collapse empty lines
	if opts.CollapseEmptyLines {
		content = collapseEmptyLines(content)
	}

	// Trim trailing whitespace
	if opts.TrimWhitespace {
		content = trimTrailingWhitespace(content)
	}

	return content
}

// stripCSSComments removes CSS comments
func stripCSSComments(content string) string {
	commentRE := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	return commentRE.ReplaceAllString(content, "")
}

// collapseCSSImports collapses CSS @import statements
func collapseCSSImports(content string) string {
	importRE := regexp.MustCompile(`(?m)^@import\s+.*?;$`)
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	importLines := make([]string, 0)
	inImportBlock := false

	for _, line := range lines {
		if importRE.MatchString(line) {
			importLines = append(importLines, line)
			inImportBlock = true
		} else if strings.TrimSpace(line) == "" && inImportBlock {
			// Empty line within import block - continue collecting
			continue
		} else {
			// Non-import line - flush collected imports
			if len(importLines) > 0 {
				result = append(result, "/* "+string(rune(len(importLines)+'0'))+" imports collapsed */")
				importLines = importLines[:0]
				inImportBlock = false
			}
			result = append(result, line)
		}
	}

	// Flush remaining imports
	if len(importLines) > 0 {
		result = append(result, "/* "+string(rune(len(importLines)+'0'))+" imports collapsed */")
	}

	return strings.Join(result, "\n")
}
