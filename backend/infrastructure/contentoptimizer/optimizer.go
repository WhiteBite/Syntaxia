package contentoptimizer

import (
	"context"
	"path/filepath"
	"strings"
	"syntaxia/domain"
)

// ContentOptimizer implements domain.ContentOptimizer for optimizing file content for AI context
type ContentOptimizer struct {
	log domain.Logger
}

// NewContentOptimizer creates a new content optimizer
func NewContentOptimizer(log domain.Logger) domain.ContentOptimizer {
	return &ContentOptimizer{
		log: log,
	}
}

// Optimize applies optimizations to file content based on options
func (o *ContentOptimizer) Optimize(ctx context.Context, content, filePath string, opts domain.ContentOptimizeOptions) string {
	// Detect language from file extension
	lang := detectLanguage(filePath)

	// Apply language-specific optimization
	switch lang {
	case "go":
		return o.optimizeGo(content, opts)
	case "typescript", "javascript":
		return o.optimizeTypeScript(content, opts)
	case "python":
		return o.optimizePython(content, opts)
	case "vue":
		return o.optimizeVue(content, opts)
	case "css", "scss", "sass", "less":
		return o.optimizeCSS(content, opts)
	default:
		return o.optimizeGeneric(content, opts)
	}
}

// OptimizeWithDefaults applies optimizations with default settings
func (o *ContentOptimizer) OptimizeWithDefaults(ctx context.Context, content, filePath string) string {
	opts := domain.ContentOptimizeOptions{
		StripComments:      true,
		CollapseEmptyLines: true,
		StripLicense:       true,
		TrimWhitespace:     true,
		CompactDataFiles:   false,
		SkeletonMode:       false,
	}
	return o.Optimize(ctx, content, filePath, opts)
}

// CanGenerateSkeleton checks if skeleton generation is supported for the file
func (o *ContentOptimizer) CanGenerateSkeleton(filePath string) bool {
	lang := detectLanguage(filePath)
	// Currently only Go supports skeleton mode (via go/parser)
	return lang == "go"
}

// detectLanguage detects programming language from file extension
func detectLanguage(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".py":
		return "python"
	case ".vue":
		return "vue"
	case ".css":
		return "css"
	case ".scss":
		return "scss"
	case ".sass":
		return "sass"
	case ".less":
		return "less"
	default:
		return "unknown"
	}
}

// optimizeGeneric applies basic optimizations for unknown file types
func (o *ContentOptimizer) optimizeGeneric(content string, opts domain.ContentOptimizeOptions) string {
	if opts.CollapseEmptyLines {
		content = collapseEmptyLines(content)
	}
	if opts.TrimWhitespace {
		content = trimTrailingWhitespace(content)
	}
	return content
}

// collapseEmptyLines replaces 3+ consecutive newlines with 2 newlines
func collapseEmptyLines(content string) string {
	lines := strings.Split(content, "\n")
	result := make([]string, 0, len(lines))
	emptyCount := 0

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			emptyCount++
			if emptyCount <= 1 {
				result = append(result, line)
			}
		} else {
			emptyCount = 0
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// trimTrailingWhitespace removes trailing whitespace from each line
func trimTrailingWhitespace(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " \t")
	}
	return strings.Join(lines, "\n")
}
