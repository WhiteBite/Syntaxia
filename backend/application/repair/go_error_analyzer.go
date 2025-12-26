package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// goErrorPattern defines a Go error pattern with its type and suggestion
type goErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Go error patterns for better error classification
var goErrorPatterns = []goErrorPattern{
	{regexp.MustCompile(`undefined: (\w+)`), domain.ErrorTypeImport, "Add missing import or declaration for '%s'"},
	{regexp.MustCompile(`cannot use .* as .* in`), domain.ErrorTypeTypeCheck, "Check type compatibility"},
	{regexp.MustCompile(`imported and not used: "([^"]+)"`), domain.ErrorTypeLinting, "Remove unused import '%s'"},
	{regexp.MustCompile(`(\w+) declared (but|and) not used`), domain.ErrorTypeLinting, "Remove or use the variable '%s'"},
	{regexp.MustCompile(`undeclared name: (\w+)`), domain.ErrorTypeImport, "Declare or import '%s'"},
	{regexp.MustCompile(`syntax error:`), domain.ErrorTypeSyntax, "Fix syntax error"},
	{regexp.MustCompile(`expected .*, found`), domain.ErrorTypeSyntax, "Fix syntax - check for missing tokens"},
	{regexp.MustCompile(`missing return`), domain.ErrorTypeCompilation, "Add missing return statement"},
	{regexp.MustCompile(`too many arguments`), domain.ErrorTypeTypeCheck, "Check function signature - too many arguments"},
	{regexp.MustCompile(`not enough arguments`), domain.ErrorTypeTypeCheck, "Check function signature - missing arguments"},
	{regexp.MustCompile(`cannot assign to`), domain.ErrorTypeTypeCheck, "Cannot assign - check if variable is assignable"},
	{regexp.MustCompile(`invalid operation:`), domain.ErrorTypeTypeCheck, "Invalid operation - check operand types"},
}

// GoErrorAnalyzer analyzes Go-specific errors
type GoErrorAnalyzer struct{}

// NewGoErrorAnalyzer creates a new Go error analyzer
func NewGoErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &GoErrorAnalyzer{}
}

// AnalyzeError analyzes Go error output and extracts details
func (g *GoErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "go",
		Suggestions: make([]string, 0),
	}

	for _, pattern := range goErrorPatterns {
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

	// Fallback classification
	if details.ErrorType == "" {
		if strings.Contains(errorOutput, "undefined:") {
			details.ErrorType = domain.ErrorTypeImport
			details.Suggestions = append(details.Suggestions, "Check if the identifier is declared or imported")
		} else if strings.Contains(errorOutput, "cannot use") {
			details.ErrorType = domain.ErrorTypeTypeCheck
			details.Suggestions = append(details.Suggestions, "Check type compatibility")
		}
	}

	return details, nil
}

// SuggestCorrections provides Go-specific correction suggestions
func (g *GoErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Import-related errors - use goimports
	if strings.Contains(msg, "undefined:") || strings.Contains(msg, "undeclared name:") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: fmt.Sprintf("Fix imports with goimports: %s", msg),
		})
	}

	// Unused imports/variables - use goimports or golangci-lint --fix
	if strings.Contains(msg, "imported and not used") || strings.Contains(msg, "declared but not used") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionRemoveCode,
			Target:      errDetails.SourceFile,
			Description: "Remove unused import/variable with goimports",
		})
	}

	// Formatting issues
	if strings.Contains(msg, "gofmt") || strings.Contains(msg, "format") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFormatCode,
			Target:      errDetails.SourceFile,
			Description: "Format code with gofmt",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Go output
func (g *GoErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range goErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	if strings.Contains(errorOutput, "undefined:") || strings.Contains(errorOutput, "undeclared name:") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "cannot use") {
		return domain.ErrorTypeTypeCheck
	}
	return domain.ErrorTypeCompilation
}

// GetLanguage returns the supported language
func (g *GoErrorAnalyzer) GetLanguage() string {
	return langGo
}
