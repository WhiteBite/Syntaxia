package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// tsErrorPattern defines a TypeScript error pattern
type tsErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// TypeScript error patterns
var tsErrorPatterns = []tsErrorPattern{
	{regexp.MustCompile(`TS2304: Cannot find name '(\w+)'`), domain.ErrorTypeImport, "Add import for '%s'"},
	{regexp.MustCompile(`TS2322: Type '(.+)' is not assignable`), domain.ErrorTypeTypeCheck, "Fix type assignment"},
	{regexp.MustCompile(`TS2339: Property '(\w+)' does not exist`), domain.ErrorTypeTypeCheck, "Property '%s' not found - check type definition"},
	{regexp.MustCompile(`TS2345: Argument of type '(.+)' is not assignable`), domain.ErrorTypeTypeCheck, "Fix argument type"},
	{regexp.MustCompile(`TS2307: Cannot find module '([^']+)'`), domain.ErrorTypeImport, "Install or import module '%s'"},
	{regexp.MustCompile(`TS1005: '(.+)' expected`), domain.ErrorTypeSyntax, "Syntax error - '%s' expected"},
	{regexp.MustCompile(`TS2551: Property '(\w+)' does not exist.*Did you mean '(\w+)'`), domain.ErrorTypeTypeCheck, "Typo: use '%s' instead"},
	{regexp.MustCompile(`TS6133: '(\w+)' is declared but`), domain.ErrorTypeLinting, "Remove unused variable '%s'"},
	{regexp.MustCompile(`TS6196: '(\w+)' is declared but never used`), domain.ErrorTypeLinting, "Remove unused '%s'"},
}

// TypeScriptErrorAnalyzer analyzes TypeScript-specific errors
type TypeScriptErrorAnalyzer struct{}

// NewTypeScriptErrorAnalyzer creates a new TypeScript error analyzer
func NewTypeScriptErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &TypeScriptErrorAnalyzer{}
}

// AnalyzeError analyzes TypeScript error output and extracts details
func (t *TypeScriptErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "tsc",
		Suggestions: make([]string, 0),
	}

	for _, pattern := range tsErrorPatterns {
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
	if details.ErrorType == "" {
		if strings.Contains(errorOutput, "TS2304") {
			details.ErrorType = domain.ErrorTypeImport
			details.Suggestions = append(details.Suggestions, "Cannot find name - check imports or declarations")
		} else if strings.Contains(errorOutput, "TS2322") {
			details.ErrorType = domain.ErrorTypeTypeCheck
			details.Suggestions = append(details.Suggestions, "Type assignment error - check type compatibility")
		}
	}

	return details, nil
}

// SuggestCorrections provides TypeScript-specific correction suggestions
func (t *TypeScriptErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	if strings.Contains(msg, "TS2304") || strings.Contains(msg, "Cannot find name") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: fmt.Sprintf("Add missing import: %s", msg),
		})
	}

	if strings.Contains(msg, "TS2307") || strings.Contains(msg, "Cannot find module") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Install missing module or fix import path",
		})
	}

	if strings.Contains(msg, "TS6133") || strings.Contains(msg, "TS6196") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionRemoveCode,
			Target:      errDetails.SourceFile,
			Description: "Remove unused declaration",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from TypeScript output
func (t *TypeScriptErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range tsErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	if strings.Contains(errorOutput, "TS2304") || strings.Contains(errorOutput, "Cannot find") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "TS2322") || strings.Contains(errorOutput, "Type") {
		return domain.ErrorTypeTypeCheck
	}
	return domain.ErrorTypeCompilation
}

// GetLanguage returns the supported language
func (t *TypeScriptErrorAnalyzer) GetLanguage() string {
	return langTypeScript
}
