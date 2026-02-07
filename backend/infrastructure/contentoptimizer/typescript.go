package contentoptimizer

import (
	"regexp"
	"strings"
	"syntaxia/domain"
)

// optimizeTypeScript applies TypeScript/JavaScript-specific optimizations
func (o *ContentOptimizer) optimizeTypeScript(content string, opts domain.ContentOptimizeOptions) string {
	// Remove comments first (before collapsing imports)
	if opts.StripComments {
		content = stripTSComments(content)
	}

	// Remove type definitions
	if opts.StripComments {
		content = removeTSTypeDefinitions(content)
	}

	// Collapse imports
	if opts.StripComments {
		content = collapseTSImports(content)
	}

	// Remove boilerplate
	if opts.StripLicense {
		content = stripTSBoilerplate(content)
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

// stripTSComments removes TypeScript/JavaScript comments
func stripTSComments(content string) string {
	// Remove single-line comments
	lineCommentRE := regexp.MustCompile(`(?m)^\s*//.*$`)
	content = lineCommentRE.ReplaceAllString(content, "")

	// Remove multi-line comments (including JSDoc)
	blockCommentRE := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	content = blockCommentRE.ReplaceAllString(content, "")

	return content
}

// removeTSTypeDefinitions removes TypeScript type definitions
func removeTSTypeDefinitions(content string) string {
	// Remove interface declarations (multiline)
	interfaceRE := regexp.MustCompile(`(?m)^export\s+interface\s+\w+[^{]*\{[^}]*\}\s*$`)
	content = interfaceRE.ReplaceAllString(content, "// Type definition removed")

	// Remove type aliases
	typeAliasRE := regexp.MustCompile(`(?m)^export\s+type\s+\w+\s*=\s*.*$`)
	content = typeAliasRE.ReplaceAllString(content, "// Type definition removed")

	return content
}

// collapseTSImports collapses TypeScript/JavaScript import statements
func collapseTSImports(content string) string {
	importRE := regexp.MustCompile(`(?m)^import\s+.*?from\s+['"].*?['"];?\s*$`)
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

// stripTSBoilerplate removes common TypeScript/JavaScript boilerplate
func stripTSBoilerplate(content string) string {
	// Remove 'use strict'
	useStrictRE := regexp.MustCompile(`['"]use strict['"];?\s*`)
	content = useStrictRE.ReplaceAllString(content, "")

	// Remove license headers at the beginning
	licenseRE := regexp.MustCompile(`(?s)^(//.*?\n|/\*.*?\*/\n)+`)
	if strings.Contains(content, "Copyright") || strings.Contains(content, "License") {
		content = licenseRE.ReplaceAllString(content, "// License header removed\n")
	}

	return content
}
