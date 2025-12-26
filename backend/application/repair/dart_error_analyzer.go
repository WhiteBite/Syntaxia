package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// dartErrorPattern defines a Dart error pattern
type dartErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Dart error patterns - comprehensive list for Dart/Flutter error classification
var dartErrorPatterns = []dartErrorPattern{
	// Undefined name errors
	{regexp.MustCompile(`Undefined name '(\w+)'`), domain.ErrorTypeImport, "Import or declare '%s'"},
	{regexp.MustCompile(`Undefined class '(\w+)'`), domain.ErrorTypeImport, "Import or declare class '%s'"},
	// Type errors
	{regexp.MustCompile(`The argument type '([^']+)' can't be assigned to the parameter type '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix type mismatch: expected '%s'"},
	{regexp.MustCompile(`A value of type '([^']+)' can't be assigned to a variable of type '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix type assignment"},
	{regexp.MustCompile(`The return type '([^']+)' isn't a '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix return type"},
	// Method/property errors
	{regexp.MustCompile(`The method '(\w+)' isn't defined for the type '([^']+)'`), domain.ErrorTypeTypeCheck, "Method '%s' not found on type '%s'"},
	{regexp.MustCompile(`The getter '(\w+)' isn't defined for the type '([^']+)'`), domain.ErrorTypeTypeCheck, "Getter '%s' not found on type '%s'"},
	{regexp.MustCompile(`The setter '(\w+)' isn't defined for the type '([^']+)'`), domain.ErrorTypeTypeCheck, "Setter '%s' not found on type '%s'"},
	// Syntax errors
	{regexp.MustCompile(`Expected ';' after this`), domain.ErrorTypeSyntax, "Add missing semicolon"},
	{regexp.MustCompile(`Expected '\)' before this`), domain.ErrorTypeSyntax, "Add missing closing parenthesis"},
	{regexp.MustCompile(`Expected '\}' before this`), domain.ErrorTypeSyntax, "Add missing closing brace"},
	{regexp.MustCompile(`Expected an identifier`), domain.ErrorTypeSyntax, "Add missing identifier"},
	// Missing implementation
	{regexp.MustCompile(`Missing concrete implementation of '([^']+)'`), domain.ErrorTypeCompilation, "Implement missing method '%s'"},
	{regexp.MustCompile(`'(\w+)' doesn't implement '(\w+)'`), domain.ErrorTypeCompilation, "Class '%s' must implement '%s'"},
	// Null safety errors
	{regexp.MustCompile(`The value 'null' can't be assigned to a variable of type '([^']+)'`), domain.ErrorTypeTypeCheck, "Handle null value for type '%s'"},
	{regexp.MustCompile(`A nullable expression can't be used as a condition`), domain.ErrorTypeTypeCheck, "Add null check before condition"},
	{regexp.MustCompile(`The property '(\w+)' can't be unconditionally accessed because the receiver can be 'null'`), domain.ErrorTypeTypeCheck, "Add null check before accessing '%s'"},
	// Import errors
	{regexp.MustCompile(`Target of URI doesn't exist: '([^']+)'`), domain.ErrorTypeDependency, "Fix import path '%s'"},
	{regexp.MustCompile(`Undefined name '(\w+)' in library`), domain.ErrorTypeImport, "Import library containing '%s'"},
	// Const/final errors
	{regexp.MustCompile(`'(\w+)' can't be used as a setter because it's final`), domain.ErrorTypeCompilation, "Cannot modify final variable '%s'"},
	{regexp.MustCompile(`Const variables must be initialized with a constant value`), domain.ErrorTypeCompilation, "Use constant value for const variable"},
	// Flutter widget errors
	{regexp.MustCompile(`The method 'build' must be implemented`), domain.ErrorTypeCompilation, "Implement build() method in widget"},
	{regexp.MustCompile(`The method 'createState' must be implemented`), domain.ErrorTypeCompilation, "Implement createState() method"},
}

// dartLocationPattern matches Dart analyzer error location format
var dartLocationPattern = regexp.MustCompile(`([^:]+):(\d+):(\d+):`)

// DartErrorAnalyzer analyzes Dart-specific errors
type DartErrorAnalyzer struct{}

// NewDartErrorAnalyzer creates a new Dart error analyzer
func NewDartErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &DartErrorAnalyzer{}
}

// AnalyzeError analyzes Dart error output and extracts details
func (d *DartErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "dart",
		Suggestions: make([]string, 0),
	}

	// Extract file location
	d.extractDartLocation(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range dartErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := d.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = d.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Dart analyzer output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Dart-specific correction suggestions
func (d *DartErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Undefined name - suggest import
	if strings.Contains(msg, "Undefined name") || strings.Contains(msg, "Undefined class") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import statement",
		})
	}

	// Type errors
	if strings.Contains(msg, "can't be assigned") || strings.Contains(msg, "type") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "Expected") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	}

	// Missing implementation
	if strings.Contains(msg, "Missing concrete implementation") || strings.Contains(msg, "must be implemented") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Implement missing method or property",
		})
	}

	// Null safety
	if strings.Contains(msg, "null") || strings.Contains(msg, "nullable") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Add null check or use null-aware operator",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Dart output
func (d *DartErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range dartErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return d.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (d *DartErrorAnalyzer) GetLanguage() string {
	return langDart
}

// extractDartLocation extracts file and line information from Dart error output
func (d *DartErrorAnalyzer) extractDartLocation(errorOutput string, details *domain.ErrorDetails) {
	if matches := dartLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 4 {
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
func (d *DartErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (d *DartErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "undefined name") || strings.Contains(errorLower, "undefined class"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "target of uri") || strings.Contains(errorLower, "package:"):
		return domain.ErrorTypeDependency
	case strings.Contains(errorLower, "can't be assigned") || strings.Contains(errorLower, "isn't defined"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "expected") || strings.Contains(errorLower, "unexpected"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "null") || strings.Contains(errorLower, "nullable"):
		return domain.ErrorTypeTypeCheck
	default:
		return domain.ErrorTypeCompilation
	}
}
