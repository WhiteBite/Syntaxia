package contentoptimizer

import (
	"regexp"
	"syntaxia/domain"
)

// optimizeVue applies Vue-specific optimizations
func (o *ContentOptimizer) optimizeVue(content string, opts domain.ContentOptimizeOptions) string {
	// Extract and clean script section
	scriptRE := regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	content = scriptRE.ReplaceAllStringFunc(content, func(match string) string {
		// Extract script content
		scriptMatch := scriptRE.FindStringSubmatch(match)
		if len(scriptMatch) < 2 {
			return match
		}
		scriptContent := scriptMatch[1]

		// Apply TypeScript optimizations to script content
		cleanedScript := o.optimizeTypeScript(scriptContent, opts)

		// Reconstruct script tag
		return "<script>" + cleanedScript + "</script>"
	})

	// Remove HTML comments if option enabled
	if opts.StripComments {
		htmlCommentRE := regexp.MustCompile(`<!--[\s\S]*?-->`)
		content = htmlCommentRE.ReplaceAllString(content, "")
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
