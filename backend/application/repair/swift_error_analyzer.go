package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// swiftErrorPattern defines a Swift error pattern
type swiftErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Swift error patterns - comprehensive list for Swift compiler error classification
var swiftErrorPatterns = []swiftErrorPattern{
	// Undefined/unresolved errors
	{regexp.MustCompile(`cannot find '(\w+)' in scope`), domain.ErrorTypeImport, "Import module or declare '%s'"},
	{regexp.MustCompile(`use of unresolved identifier '(\w+)'`), domain.ErrorTypeImport, "Import module or declare '%s'"},
	{regexp.MustCompile(`cannot find type '(\w+)' in scope`), domain.ErrorTypeImport, "Import module containing type '%s'"},
	// Type errors
	{regexp.MustCompile(`cannot convert value of type '([^']+)' to expected argument type '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix type mismatch: expected '%s'"},
	{regexp.MustCompile(`cannot assign value of type '([^']+)' to type '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix type assignment"},
	{regexp.MustCompile(`value of type '([^']+)' has no member '(\w+)'`), domain.ErrorTypeTypeCheck, "Type '%s' has no member '%s'"},
	{regexp.MustCompile(`type '([^']+)' does not conform to protocol '([^']+)'`), domain.ErrorTypeTypeCheck, "Implement protocol '%s' conformance"},
	// Syntax errors
	{regexp.MustCompile(`expected '([^']+)' in`), domain.ErrorTypeSyntax, "Add missing '%s'"},
	{regexp.MustCompile(`expected expression`), domain.ErrorTypeSyntax, "Add missing expression"},
	{regexp.MustCompile(`expected declaration`), domain.ErrorTypeSyntax, "Add missing declaration"},
	{regexp.MustCompile(`unexpected '([^']+)' in`), domain.ErrorTypeSyntax, "Remove unexpected '%s'"},
	// Optional/nil errors
	{regexp.MustCompile(`value of optional type '([^']+)' must be unwrapped`), domain.ErrorTypeTypeCheck, "Unwrap optional with ? or !"},
	{regexp.MustCompile(`'nil' cannot be assigned to type '([^']+)'`), domain.ErrorTypeTypeCheck, "Make type optional or provide non-nil value"},
	{regexp.MustCompile(`cannot use optional chaining on non-optional value`), domain.ErrorTypeTypeCheck, "Remove optional chaining (?) from non-optional"},
	// Function/method errors
	{regexp.MustCompile(`missing argument for parameter '(\w+)'`), domain.ErrorTypeTypeCheck, "Add missing argument '%s'"},
	{regexp.MustCompile(`extra argument '(\w+)' in call`), domain.ErrorTypeTypeCheck, "Remove extra argument '%s'"},
	{regexp.MustCompile(`incorrect argument label in call \(have '(\w+)', expected '(\w+)'\)`), domain.ErrorTypeTypeCheck, "Change argument label from '%s' to '%s'"},
	{regexp.MustCompile(`missing return in a function expected to return '([^']+)'`), domain.ErrorTypeTypeCheck, "Add return statement"},
	// Access control errors
	{regexp.MustCompile(`'(\w+)' is inaccessible due to '(\w+)' protection level`), domain.ErrorTypeCompilation, "'%s' is %s - change access level"},
	{regexp.MustCompile(`cannot override '(\w+)' which has been marked unavailable`), domain.ErrorTypeCompilation, "Cannot override unavailable member '%s'"},
	// Initialization errors
	{regexp.MustCompile(`'self' used before all stored properties are initialized`), domain.ErrorTypeCompilation, "Initialize all stored properties before using self"},
	{regexp.MustCompile(`property '(\w+)' not initialized at super\.init call`), domain.ErrorTypeCompilation, "Initialize property '%s' before super.init"},
	{regexp.MustCompile(`return from initializer without initializing all stored properties`), domain.ErrorTypeCompilation, "Initialize all stored properties"},
	// Mutability errors
	{regexp.MustCompile(`cannot assign to property: '(\w+)' is a 'let' constant`), domain.ErrorTypeCompilation, "Change '%s' from 'let' to 'var'"},
	{regexp.MustCompile(`cannot use mutating member on immutable value`), domain.ErrorTypeCompilation, "Use 'var' instead of 'let' for mutable value"},
	// Concurrency errors (Swift 5.5+)
	{regexp.MustCompile(`actor-isolated property '(\w+)' can not be referenced from a non-isolated context`), domain.ErrorTypeCompilation, "Use 'await' or mark context as actor-isolated"},
	{regexp.MustCompile(`expression is 'async' but is not marked with 'await'`), domain.ErrorTypeCompilation, "Add 'await' before async expression"},
	// Generic errors
	{regexp.MustCompile(`error: (.+)`), domain.ErrorTypeCompilation, "Swift error: %s"},
}

// swiftLocationPattern matches Swift compiler error location format
var swiftLocationPattern = regexp.MustCompile(`([^:]+):(\d+):(\d+):`)

// SwiftErrorAnalyzer analyzes Swift-specific errors
type SwiftErrorAnalyzer struct{}

// NewSwiftErrorAnalyzer creates a new Swift error analyzer
func NewSwiftErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &SwiftErrorAnalyzer{}
}

// AnalyzeError analyzes Swift error output and extracts details
func (s *SwiftErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "swiftc",
		Suggestions: make([]string, 0),
	}

	// Extract file location
	s.extractSwiftLocation(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range swiftErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := s.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = s.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Swift compiler output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Swift-specific correction suggestions
func (s *SwiftErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Undefined/unresolved - suggest import
	if strings.Contains(msg, "cannot find") || strings.Contains(msg, "unresolved identifier") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import statement",
		})
	}

	// Type errors
	if strings.Contains(msg, "cannot convert") || strings.Contains(msg, "cannot assign") ||
		strings.Contains(msg, "has no member") || strings.Contains(msg, "does not conform") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "expected") || strings.Contains(msg, "unexpected") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	}

	// Optional handling
	if strings.Contains(msg, "optional") || strings.Contains(msg, "unwrapped") || strings.Contains(msg, "nil") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Handle optional value properly (use ?, !, or if let)",
		})
	}

	// Function arguments
	if strings.Contains(msg, "missing argument") || strings.Contains(msg, "extra argument") ||
		strings.Contains(msg, "argument label") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix function call arguments",
		})
	}

	// Mutability
	if strings.Contains(msg, "let constant") || strings.Contains(msg, "immutable") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Change 'let' to 'var' for mutable value",
		})
	}

	// Async/await
	if strings.Contains(msg, "async") || strings.Contains(msg, "await") || strings.Contains(msg, "actor") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Add 'await' or fix concurrency context",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Swift output
func (s *SwiftErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range swiftErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return s.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (s *SwiftErrorAnalyzer) GetLanguage() string {
	return langSwift
}

// extractSwiftLocation extracts file and line information from Swift error output
func (s *SwiftErrorAnalyzer) extractSwiftLocation(errorOutput string, details *domain.ErrorDetails) {
	if matches := swiftLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 4 {
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
func (s *SwiftErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (s *SwiftErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "cannot find") || strings.Contains(errorLower, "unresolved"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "cannot convert") || strings.Contains(errorLower, "cannot assign") ||
		strings.Contains(errorLower, "has no member") || strings.Contains(errorLower, "does not conform"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "expected") || strings.Contains(errorLower, "unexpected"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "optional") || strings.Contains(errorLower, "nil"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "async") || strings.Contains(errorLower, "await") || strings.Contains(errorLower, "actor"):
		return domain.ErrorTypeCompilation
	default:
		return domain.ErrorTypeCompilation
	}
}
