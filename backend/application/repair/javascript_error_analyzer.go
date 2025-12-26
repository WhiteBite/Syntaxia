package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// jsErrorPattern defines a JavaScript error pattern
type jsErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// JavaScript error patterns
var jsErrorPatterns = []jsErrorPattern{
	{regexp.MustCompile(`'(\w+)' is not defined`), domain.ErrorTypeImport, "Import or declare '%s'"},
	{regexp.MustCompile(`SyntaxError:`), domain.ErrorTypeSyntax, "Fix syntax error"},
	{regexp.MustCompile(`Unexpected token`), domain.ErrorTypeSyntax, "Fix unexpected token"},
	{regexp.MustCompile(`Cannot read propert.* of (undefined|null)`), domain.ErrorTypeCompilation, "Check for null/undefined"},
	{regexp.MustCompile(`'(\w+)' is assigned .* but never used`), domain.ErrorTypeLinting, "Remove unused variable '%s'"},
	{regexp.MustCompile(`'(\w+)' is defined but never used`), domain.ErrorTypeLinting, "Remove unused '%s'"},
}

// JavaScriptErrorAnalyzer analyzes JavaScript-specific errors
type JavaScriptErrorAnalyzer struct{}

// NewJavaScriptErrorAnalyzer creates a new JavaScript error analyzer
func NewJavaScriptErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &JavaScriptErrorAnalyzer{}
}

// AnalyzeError analyzes JavaScript error output and extracts details
func (j *JavaScriptErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "eslint",
		Suggestions: make([]string, 0),
	}

	for _, pattern := range jsErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := pattern.suggestion
			if len(matches) > 1 && strings.Contains(suggestion, "%s") {
				suggestion = fmt.Sprintf(suggestion, matches[1])
			}
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback
	if details.ErrorType == "" && strings.Contains(errorOutput, "is not defined") {
		details.ErrorType = domain.ErrorTypeImport
		details.Suggestions = append(details.Suggestions, "Variable or function not defined - check imports")
	}

	return details, nil
}

// SuggestCorrections provides JavaScript-specific correction suggestions
func (j *JavaScriptErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	if strings.Contains(msg, "is not defined") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import or declaration",
		})
	}

	if strings.Contains(msg, "never used") || strings.Contains(msg, "but never used") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionRemoveCode,
			Target:      errDetails.SourceFile,
			Description: "Remove unused variable/import",
		})
	}

	if strings.Contains(msg, "SyntaxError") || strings.Contains(msg, "Unexpected token") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from JavaScript output
func (j *JavaScriptErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range jsErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	if strings.Contains(errorOutput, "is not defined") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "SyntaxError") {
		return domain.ErrorTypeSyntax
	}
	return domain.ErrorTypeCompilation
}

// GetLanguage returns the supported language
func (j *JavaScriptErrorAnalyzer) GetLanguage() string {
	return langJavaScript
}
