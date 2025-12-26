package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// kotlinErrorPattern defines a Kotlin error pattern with its type and suggestion
type kotlinErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Kotlin error patterns - comprehensive list for kotlinc output
var kotlinErrorPatterns = []kotlinErrorPattern{
	// Unresolved reference errors
	{regexp.MustCompile(`Unresolved reference: (\w+)`), domain.ErrorTypeImport, "Import or declare '%s'"},
	{regexp.MustCompile(`Unresolved reference\. None of the following candidates is applicable`), domain.ErrorTypeTypeCheck, "Check function signature and argument types"},
	// Type mismatch errors
	{regexp.MustCompile(`Type mismatch: inferred type is (.+) but (.+) was expected`), domain.ErrorTypeTypeCheck, "Fix type mismatch: expected %s"},
	{regexp.MustCompile(`Type mismatch\. Required: (.+) Found: (.+)`), domain.ErrorTypeTypeCheck, "Fix type mismatch: required %s"},
	// Syntax errors
	{regexp.MustCompile(`Expecting '([^']+)'`), domain.ErrorTypeSyntax, "Add missing '%s'"},
	{regexp.MustCompile(`Unexpected tokens`), domain.ErrorTypeSyntax, "Fix syntax - unexpected tokens"},
	// Val/var reassignment errors
	{regexp.MustCompile(`Val cannot be reassigned`), domain.ErrorTypeCompilation, "Use 'var' instead of 'val' for mutable variable"},
	// Function call errors
	{regexp.MustCompile(`None of the following functions can be called with the arguments supplied`), domain.ErrorTypeTypeCheck, "Check function arguments - no matching overload"},
	{regexp.MustCompile(`Too many arguments for`), domain.ErrorTypeTypeCheck, "Remove extra arguments"},
	{regexp.MustCompile(`No value passed for parameter '(\w+)'`), domain.ErrorTypeTypeCheck, "Add missing argument '%s'"},
	// Null safety errors
	{regexp.MustCompile(`Only safe \(\?\.\) or non-null asserted \(!!\.\) calls are allowed on a nullable receiver`), domain.ErrorTypeTypeCheck, "Use safe call (?.) or assert non-null (!!)"},
	{regexp.MustCompile(`Null can not be a value of a non-null type`), domain.ErrorTypeTypeCheck, "Handle nullable type properly"},
	// Visibility errors
	{regexp.MustCompile(`Cannot access '(\w+)': it is (\w+) in '(\w+)'`), domain.ErrorTypeTypeCheck, "'%s' is %s - change visibility or use accessor"},
	// Override errors
	{regexp.MustCompile(`'(\w+)' overrides nothing`), domain.ErrorTypeTypeCheck, "Remove 'override' or check parent class"},
	{regexp.MustCompile(`'(\w+)' hides member of supertype`), domain.ErrorTypeTypeCheck, "Add 'override' modifier"},
	// Class/interface errors
	{regexp.MustCompile(`Class '(\w+)' is not abstract and does not implement abstract member`), domain.ErrorTypeTypeCheck, "Implement abstract member or make class abstract"},
	// Generic errors
	{regexp.MustCompile(`e: (.+)`), domain.ErrorTypeCompilation, "Kotlin error: %s"},
}

// kotlinLocationPattern matches Kotlin error location format: e: file.kt:10:5:
var kotlinLocationPattern = regexp.MustCompile(`e:\s*([^:]+):(\d+):(\d+):`)

// KotlinErrorAnalyzer analyzes Kotlin compiler error output
type KotlinErrorAnalyzer struct{}

// NewKotlinErrorAnalyzer creates a new Kotlin error analyzer
func NewKotlinErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &KotlinErrorAnalyzer{}
}

// AnalyzeError analyzes Kotlin error output and extracts details
func (k *KotlinErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "kotlinc",
		Suggestions: make([]string, 0),
	}

	// Extract file location from Kotlin error format
	k.extractLocationInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range kotlinErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := k.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = k.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Kotlin compiler output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Kotlin-specific correction suggestions
func (k *KotlinErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Unresolved reference - suggest import
	if strings.Contains(msg, "Unresolved reference") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import statement or declare the symbol",
		})
	}

	// Type mismatch
	if strings.Contains(msg, "Type mismatch") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch - check expected vs actual types",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "Expecting") || strings.Contains(msg, "Unexpected") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix Kotlin syntax error",
		})
	}

	// Val reassignment
	if strings.Contains(msg, "Val cannot be reassigned") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Change 'val' to 'var' for mutable variable",
		})
	}

	// Null safety
	if strings.Contains(msg, "nullable") || strings.Contains(msg, "Null can not be") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Handle nullable type with ?. or !! or ?:",
		})
	}

	// Function arguments
	if strings.Contains(msg, "None of the following functions") || strings.Contains(msg, "Too many arguments") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix function call arguments",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Kotlin output
func (k *KotlinErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range kotlinErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return k.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (k *KotlinErrorAnalyzer) GetLanguage() string {
	return langKotlin
}

// extractLocationInfo extracts file and line information from Kotlin error output
func (k *KotlinErrorAnalyzer) extractLocationInfo(errorOutput string, details *domain.ErrorDetails) {
	if matches := kotlinLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 4 {
		details.SourceFile = matches[1]
		if line, err := parseIntSafe(matches[2]); err == nil {
			details.LineNumber = line
		}
		if col, err := parseIntSafe(matches[3]); err == nil {
			details.Column = col
		}
	}
}

// formatSuggestion formats a suggestion string with captured groups
func (k *KotlinErrorAnalyzer) formatSuggestion(template string, matches []string) string {
	if len(matches) <= 1 || !strings.Contains(template, "%s") {
		return template
	}
	placeholderCount := strings.Count(template, "%s")
	args := make([]interface{}, 0, placeholderCount)
	for i := 1; i < len(matches) && len(args) < placeholderCount; i++ {
		args = append(args, matches[i])
	}
	for len(args) < placeholderCount {
		args = append(args, "")
	}
	return fmt.Sprintf(template, args...)
}

// classifyByKeywords classifies error type based on keywords
func (k *KotlinErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "unresolved reference"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "type mismatch") ||
		strings.Contains(errorLower, "none of the following functions") ||
		strings.Contains(errorLower, "nullable"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "expecting") || strings.Contains(errorLower, "unexpected"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "val cannot be reassigned") || strings.Contains(errorLower, "override"):
		return domain.ErrorTypeCompilation
	default:
		return domain.ErrorTypeCompilation
	}
}
