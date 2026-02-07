package contentoptimizer

import (
	"regexp"
	"strings"
	"syntaxia/domain"
)

// optimizeGo applies Go-specific optimizations
func (o *ContentOptimizer) optimizeGo(content string, opts domain.ContentOptimizeOptions) string {
	// Remove comments first (before collapsing imports)
	if opts.StripComments {
		content = stripGoComments(content)
	}

	// Collapse import blocks
	if opts.StripComments {
		content = collapseGoImports(content)
	}

	// Strip license headers
	if opts.StripLicense {
		content = stripGoLicense(content)
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

// stripGoComments removes Go comments (// and /* */)
func stripGoComments(content string) string {
	// Remove single-line comments
	lineCommentRE := regexp.MustCompile(`(?m)^\s*//.*$`)
	content = lineCommentRE.ReplaceAllString(content, "")

	// Remove multi-line comments
	blockCommentRE := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	content = blockCommentRE.ReplaceAllString(content, "")

	return content
}

// collapseGoImports collapses import blocks into a single comment
func collapseGoImports(content string) string {
	// Match import blocks: import ( ... )
	importBlockRE := regexp.MustCompile(`(?s)import\s*\(\s*\n(.*?)\n\s*\)`)
	content = importBlockRE.ReplaceAllStringFunc(content, func(match string) string {
		// Count number of import lines
		lines := strings.Split(match, "\n")
		importCount := 0
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && trimmed != "import (" && trimmed != ")" {
				importCount++
			}
		}
		if importCount > 0 {
			return "// " + strings.TrimSpace(strings.Split(match, "\n")[0]) + " - " + 
				   strings.TrimSpace(strings.Join([]string{string(rune(importCount + '0')), " imports collapsed"}, ""))
		}
		return match
	})

	// Match single imports: import "..."
	singleImportRE := regexp.MustCompile(`(?m)^import\s+"[^"]+"\s*$`)
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	importLines := make([]string, 0)
	inImportBlock := false

	for _, line := range lines {
		if singleImportRE.MatchString(line) {
			importLines = append(importLines, line)
			inImportBlock = true
		} else if strings.TrimSpace(line) == "" && inImportBlock {
			// Empty line within import block - continue collecting
			continue
		} else {
			// Non-import line - flush collected imports
			if len(importLines) > 0 {
				result = append(result, "// "+string(rune(len(importLines)+'0'))+" imports collapsed")
				importLines = importLines[:0]
				inImportBlock = false
			}
			result = append(result, line)
		}
	}

	// Flush remaining imports
	if len(importLines) > 0 {
		result = append(result, "// "+string(rune(len(importLines)+'0'))+" imports collapsed")
	}

	return strings.Join(result, "\n")
}

// stripGoLicense removes license headers from Go files
func stripGoLicense(content string) string {
	// Match license headers at the beginning of file (before package declaration)
	licenseRE := regexp.MustCompile(`(?s)^(//.*?\n)+(?=package)`)
	content = licenseRE.ReplaceAllString(content, "// License header removed\n")
	return content
}
