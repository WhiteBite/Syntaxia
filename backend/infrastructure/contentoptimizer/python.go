package contentoptimizer

import (
	"regexp"
	"strings"
	"syntaxia/domain"
)

// optimizePython applies Python-specific optimizations
func (o *ContentOptimizer) optimizePython(content string, opts domain.ContentOptimizeOptions) string {
	// Remove comments and docstrings first
	if opts.StripComments {
		content = stripPythonComments(content)
	}

	// Remove type hints
	if opts.StripComments {
		content = removePythonTypeHints(content)
	}

	// Collapse imports
	if opts.StripComments {
		content = collapsePythonImports(content)
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

// stripPythonComments removes Python comments and docstrings
func stripPythonComments(content string) string {
	// Remove single-line comments
	lineCommentRE := regexp.MustCompile(`(?m)^\s*#.*$`)
	content = lineCommentRE.ReplaceAllString(content, "")

	// Remove triple-quoted docstrings (""")
	docstringRE1 := regexp.MustCompile(`"""[\s\S]*?"""`)
	content = docstringRE1.ReplaceAllString(content, "")

	// Remove triple-quoted docstrings (''')
	docstringRE2 := regexp.MustCompile(`'''[\s\S]*?'''`)
	content = docstringRE2.ReplaceAllString(content, "")

	return content
}

// removePythonTypeHints removes Python type hints
func removePythonTypeHints(content string) string {
	// Remove function parameter type hints: name: str -> name
	// Note: Go regexp doesn't support lookahead, so we use a simpler approach
	paramTypeRE := regexp.MustCompile(`(\w+):\s*\w+(\[.*?\])?\s*([,)])`)
	content = paramTypeRE.ReplaceAllString(content, "$1$3")

	// Remove function return type hints: ) -> Type:
	returnTypeRE := regexp.MustCompile(`\)\s*->\s*[^:]+:`)
	content = returnTypeRE.ReplaceAllString(content, "):")

	// Remove variable type hints: name: Type = value -> name = value
	varTypeRE := regexp.MustCompile(`(\w+):\s*\w+(\[.*?\])?\s*=`)
	content = varTypeRE.ReplaceAllString(content, "$1 =")

	return content
}

// collapsePythonImports collapses Python import statements
func collapsePythonImports(content string) string {
	importRE := regexp.MustCompile(`(?m)^(?:from\s+\S+\s+)?import\s+.*$`)
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
				result = append(result, "# "+string(rune(len(importLines)+'0'))+" imports collapsed")
				importLines = importLines[:0]
				inImportBlock = false
			}
			result = append(result, line)
		}
	}

	// Flush remaining imports
	if len(importLines) > 0 {
		result = append(result, "# "+string(rune(len(importLines)+'0'))+" imports collapsed")
	}

	return strings.Join(result, "\n")
}
