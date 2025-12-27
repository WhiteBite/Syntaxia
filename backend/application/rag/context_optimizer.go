// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"regexp"
	"sort"
	"strings"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// OptimizeOptions configures context optimization
type OptimizeOptions struct {
	RemoveComments    bool // Remove comments from code
	RemoveEmptyLines  bool // Collapse multiple empty lines
	CollapseFunctions bool // Replace function bodies with "..."
	MaxTokens         int  // Maximum total tokens
	PositionalSort    bool // Sort by priority (important first/last)
}

// FileContent represents a file with its content and metadata
type FileContent struct {
	Path      string  `json:"path"`
	Content   string  `json:"content"`
	Tokens    int     `json:"tokens"`
	Relevance float64 `json:"relevance"`
	Priority  int     `json:"priority"` // for positional sorting (higher = more important)
}

// ContextOptimizer optimizes context for AI consumption
type ContextOptimizer interface {
	// Optimize applies optimizations to file contents
	Optimize(ctx context.Context, files []FileContent, opts OptimizeOptions) ([]FileContent, error)

	// CalculateSNR calculates signal-to-noise ratio for content
	CalculateSNR(content string) float64

	// OptimizeContent optimizes a single content string
	OptimizeContent(ctx context.Context, content, filePath string, opts OptimizeOptions) string
}

// ContextOptimizerImpl implements ContextOptimizer
type ContextOptimizerImpl struct {
	analyzerRegistry analysis.AnalyzerRegistry
	log              domain.Logger
}

// Optimization constants
const (
	snrSignalWeight    = 0.7
	snrNoiseWeight     = 0.3
	minSignalRatio     = 0.1
	maxEmptyLineRun    = 1
	collapsedBodyToken = "/* ... */"
)

// Comment patterns for different languages
var (
	singleLineCommentPattern = regexp.MustCompile(`(?m)^\s*//.*$`)
	hashCommentPattern       = regexp.MustCompile(`(?m)^\s*#.*$`)
	multiLineCommentPattern  = regexp.MustCompile(`(?s)/\*.*?\*/`)
	docCommentPattern        = regexp.MustCompile(`(?s)/\*\*.*?\*/`)
	emptyLinesPattern        = regexp.MustCompile(`\n{3,}`)
	trailingWhitespace       = regexp.MustCompile(`(?m)[ \t]+$`)
)

// NewContextOptimizer creates a new ContextOptimizer
func NewContextOptimizer(
	analyzerRegistry analysis.AnalyzerRegistry,
	log domain.Logger,
) *ContextOptimizerImpl {
	return &ContextOptimizerImpl{
		analyzerRegistry: analyzerRegistry,
		log:              log,
	}
}

// Optimize applies optimizations to file contents
func (o *ContextOptimizerImpl) Optimize(ctx context.Context, files []FileContent, opts OptimizeOptions) ([]FileContent, error) {
	result := make([]FileContent, 0, len(files))

	for _, file := range files {
		optimized := o.OptimizeContent(ctx, file.Content, file.Path, opts)
		file.Content = optimized
		file.Tokens = o.estimateTokens(optimized)
		result = append(result, file)
	}

	if opts.PositionalSort {
		result = o.sortByPosition(result)
	}

	if opts.MaxTokens > 0 {
		result = o.trimToTokenBudget(result, opts.MaxTokens)
	}

	return result, nil
}

// OptimizeContent optimizes a single content string
func (o *ContextOptimizerImpl) OptimizeContent(ctx context.Context, content, filePath string, opts OptimizeOptions) string {
	result := content

	if opts.RemoveComments {
		result = o.removeComments(result, filePath)
	}

	if opts.RemoveEmptyLines {
		result = o.collapseEmptyLines(result)
	}

	result = trailingWhitespace.ReplaceAllString(result, "")

	if opts.CollapseFunctions {
		result = o.collapseFunctions(ctx, result, filePath)
	}

	return result
}

// CalculateSNR calculates signal-to-noise ratio for content
func (o *ContextOptimizerImpl) CalculateSNR(content string) float64 {
	if len(content) == 0 {
		return 0
	}

	lines := strings.Split(content, "\n")
	totalLines := len(lines)
	if totalLines == 0 {
		return 0
	}

	signalLines := 0
	noiseLines := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if o.isSignalLine(trimmed) {
			signalLines++
		} else if o.isNoiseLine(trimmed) {
			noiseLines++
		}
	}

	signalRatio := float64(signalLines) / float64(totalLines)
	noiseRatio := float64(noiseLines) / float64(totalLines)

	snr := signalRatio*snrSignalWeight - noiseRatio*snrNoiseWeight

	if snr < minSignalRatio {
		snr = minSignalRatio
	}
	if snr > 1.0 {
		snr = 1.0
	}

	return snr
}

func (o *ContextOptimizerImpl) isSignalLine(line string) bool {
	if len(line) == 0 {
		return false
	}

	signalPatterns := []string{
		"func ", "function ", "def ", "class ", "struct ", "interface ",
		"type ", "const ", "var ", "let ", "import ", "export ",
		"return ", "if ", "for ", "while ", "switch ", "case ",
		"async ", "await ", "pub ", "private ", "public ", "protected ",
	}

	lineLower := strings.ToLower(line)
	for _, pattern := range signalPatterns {
		if strings.HasPrefix(lineLower, pattern) || strings.Contains(lineLower, " "+pattern) {
			return true
		}
	}

	if strings.Contains(line, "(") && strings.Contains(line, ")") {
		return true
	}

	if strings.Contains(line, "=") && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "#") {
		return true
	}

	return false
}

func (o *ContextOptimizerImpl) isNoiseLine(line string) bool {
	if len(line) == 0 {
		return true
	}

	if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
		return true
	}

	if line == "{" || line == "}" || line == "(" || line == ")" ||
		line == "[" || line == "]" || line == ";" {
		return false
	}

	return false
}

func (o *ContextOptimizerImpl) removeComments(content, filePath string) string {
	ext := strings.ToLower(getFileExtension(filePath))

	result := content

	switch ext {
	case "go", "java", "js", "ts", "tsx", "jsx", "c", "cpp", "cs", "swift", "kt", "rs", "dart":
		result = multiLineCommentPattern.ReplaceAllString(result, "")
		result = singleLineCommentPattern.ReplaceAllString(result, "")

	case "py", "rb", "sh", "bash", "yaml", "yml":
		result = hashCommentPattern.ReplaceAllString(result, "")

	case "php":
		result = multiLineCommentPattern.ReplaceAllString(result, "")
		result = singleLineCommentPattern.ReplaceAllString(result, "")
		result = hashCommentPattern.ReplaceAllString(result, "")

	default:
		result = multiLineCommentPattern.ReplaceAllString(result, "")
		result = singleLineCommentPattern.ReplaceAllString(result, "")
	}

	return result
}

func (o *ContextOptimizerImpl) collapseEmptyLines(content string) string {
	result := emptyLinesPattern.ReplaceAllString(content, "\n\n")
	return strings.TrimSpace(result)
}

func (o *ContextOptimizerImpl) collapseFunctions(ctx context.Context, content, filePath string) string {
	analyzer := o.analyzerRegistry.GetAnalyzer(filePath)
	if analyzer == nil {
		return content
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, []byte(content))
	if err != nil {
		return content
	}

	functions := o.filterFunctions(symbols)
	if len(functions) == 0 {
		return content
	}

	sort.Slice(functions, func(i, j int) bool {
		return functions[i].EndLine > functions[j].EndLine
	})

	lines := strings.Split(content, "\n")

	for _, fn := range functions {
		if fn.EndLine-fn.StartLine < 5 {
			continue
		}

		lines = o.collapseFunctionBody(lines, fn)
	}

	return strings.Join(lines, "\n")
}

func (o *ContextOptimizerImpl) filterFunctions(symbols []analysis.Symbol) []analysis.Symbol {
	functions := make([]analysis.Symbol, 0)

	for _, sym := range symbols {
		if sym.Kind == analysis.KindFunction || sym.Kind == analysis.KindMethod {
			functions = append(functions, sym)
		}
	}

	return functions
}

func (o *ContextOptimizerImpl) collapseFunctionBody(lines []string, fn analysis.Symbol) []string {
	if fn.StartLine < 1 || fn.EndLine > len(lines) {
		return lines
	}

	bodyStart := -1
	for i := fn.StartLine - 1; i < fn.EndLine && i < len(lines); i++ {
		if strings.Contains(lines[i], "{") {
			bodyStart = i
			break
		}
	}

	if bodyStart == -1 || bodyStart >= fn.EndLine-1 {
		return lines
	}

	result := make([]string, 0, len(lines))
	result = append(result, lines[:bodyStart+1]...)
	result = append(result, "    "+collapsedBodyToken)

	closingBrace := fn.EndLine - 1
	if closingBrace < len(lines) {
		result = append(result, lines[closingBrace:]...)
	}

	return result
}

func (o *ContextOptimizerImpl) sortByPosition(files []FileContent) []FileContent {
	sort.SliceStable(files, func(i, j int) bool {
		if files[i].Priority != files[j].Priority {
			return files[i].Priority > files[j].Priority
		}
		return files[i].Relevance > files[j].Relevance
	})

	return files
}

func (o *ContextOptimizerImpl) trimToTokenBudget(files []FileContent, maxTokens int) []FileContent {
	result := make([]FileContent, 0, len(files))
	totalTokens := 0

	for _, file := range files {
		if totalTokens+file.Tokens > maxTokens {
			remaining := maxTokens - totalTokens
			if remaining > 100 {
				truncated := o.truncateContent(file.Content, remaining)
				file.Content = truncated
				file.Tokens = o.estimateTokens(truncated)
				result = append(result, file)
			}
			break
		}

		result = append(result, file)
		totalTokens += file.Tokens
	}

	return result
}

func (o *ContextOptimizerImpl) truncateContent(content string, maxTokens int) string {
	maxChars := maxTokens * charsPerToken
	if len(content) <= maxChars {
		return content
	}

	truncated := content[:maxChars]

	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxChars/2 {
		truncated = truncated[:lastNewline]
	}

	return truncated + "\n// ... truncated ..."
}

func (o *ContextOptimizerImpl) estimateTokens(text string) int {
	return len(text) / charsPerToken
}

func getFileExtension(filePath string) string {
	idx := strings.LastIndex(filePath, ".")
	if idx == -1 {
		return ""
	}
	return filePath[idx+1:]
}

// OptimizeForContext applies standard optimizations for AI context
func (o *ContextOptimizerImpl) OptimizeForContext(ctx context.Context, files []FileContent, maxTokens int) ([]FileContent, error) {
	opts := OptimizeOptions{
		RemoveComments:    false,
		RemoveEmptyLines:  true,
		CollapseFunctions: false,
		MaxTokens:         maxTokens,
		PositionalSort:    true,
	}

	return o.Optimize(ctx, files, opts)
}

// OptimizeAggressive applies aggressive optimizations to fit more content
func (o *ContextOptimizerImpl) OptimizeAggressive(ctx context.Context, files []FileContent, maxTokens int) ([]FileContent, error) {
	opts := OptimizeOptions{
		RemoveComments:    true,
		RemoveEmptyLines:  true,
		CollapseFunctions: true,
		MaxTokens:         maxTokens,
		PositionalSort:    true,
	}

	return o.Optimize(ctx, files, opts)
}

// CalculateOptimalStrategy determines the best optimization strategy
func (o *ContextOptimizerImpl) CalculateOptimalStrategy(files []FileContent, targetTokens int) OptimizeOptions {
	totalTokens := 0
	for _, f := range files {
		totalTokens += f.Tokens
	}

	ratio := float64(totalTokens) / float64(targetTokens)

	opts := OptimizeOptions{
		MaxTokens:      targetTokens,
		PositionalSort: true,
	}

	if ratio > 2.0 {
		opts.RemoveComments = true
		opts.RemoveEmptyLines = true
		opts.CollapseFunctions = true
	} else if ratio > 1.5 {
		opts.RemoveComments = true
		opts.RemoveEmptyLines = true
	} else if ratio > 1.0 {
		opts.RemoveEmptyLines = true
	}

	return opts
}
