package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// cppErrorPattern defines a C++ error pattern with its type and suggestion
type cppErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// C++ error patterns for better error classification
var cppErrorPatterns = []cppErrorPattern{
	{regexp.MustCompile(`'(\w+)' was not declared in this scope`), domain.ErrorTypeImport, "Include header or declare '%s'"},
	{regexp.MustCompile(`no matching function for call to`), domain.ErrorTypeTypeCheck, "Check function signature and arguments"},
	{regexp.MustCompile(`cannot convert '.*' to '.*'`), domain.ErrorTypeTypeCheck, "Check type compatibility"},
	{regexp.MustCompile(`expected .*;.* before`), domain.ErrorTypeSyntax, "Add missing semicolon"},
	{regexp.MustCompile(`expected '.*' before`), domain.ErrorTypeSyntax, "Fix syntax - check for missing tokens"},
	{regexp.MustCompile(`undefined reference to`), domain.ErrorTypeCompilation, "Link the required library or define the symbol"},
	{regexp.MustCompile(`redefinition of`), domain.ErrorTypeCompilation, "Remove duplicate definition"},
	{regexp.MustCompile(`'(\w+)' is not a member of`), domain.ErrorTypeTypeCheck, "Check class/struct member access"},
	{regexp.MustCompile(`invalid use of incomplete type`), domain.ErrorTypeImport, "Include the complete type definition"},
	{regexp.MustCompile(`no member named '(\w+)' in`), domain.ErrorTypeTypeCheck, "Check member name spelling or add declaration"},
}

// CPPErrorAnalyzer analyzes C++-specific errors
type CPPErrorAnalyzer struct{}

// NewCPPErrorAnalyzer creates a new C++ error analyzer
func NewCPPErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &CPPErrorAnalyzer{}
}

// AnalyzeError analyzes C++ error output and extracts details
func (c *CPPErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "g++",
		Suggestions: make([]string, 0),
	}

	for _, pattern := range cppErrorPatterns {
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
		if strings.Contains(errorOutput, "not declared") {
			details.ErrorType = domain.ErrorTypeImport
			details.Suggestions = append(details.Suggestions, "Check if the identifier is declared or include the required header")
		} else if strings.Contains(errorOutput, "cannot convert") {
			details.ErrorType = domain.ErrorTypeTypeCheck
			details.Suggestions = append(details.Suggestions, "Check type compatibility")
		}
	}

	return details, nil
}

// SuggestCorrections provides C++-specific correction suggestions
func (c *CPPErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Include-related errors
	if strings.Contains(msg, "not declared") || strings.Contains(msg, "incomplete type") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: fmt.Sprintf("Add missing #include directive: %s", msg),
		})
	}

	// Undefined reference - linker error
	if strings.Contains(msg, "undefined reference") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionAddMissingCode,
			Target:      errDetails.SourceFile,
			Description: "Link required library or implement missing function",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from C++ output
func (c *CPPErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range cppErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	if strings.Contains(errorOutput, "not declared") || strings.Contains(errorOutput, "incomplete type") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "cannot convert") || strings.Contains(errorOutput, "no matching function") {
		return domain.ErrorTypeTypeCheck
	}
	return domain.ErrorTypeCompilation
}

// GetLanguage returns the supported language
func (c *CPPErrorAnalyzer) GetLanguage() string {
	return "cpp"
}
