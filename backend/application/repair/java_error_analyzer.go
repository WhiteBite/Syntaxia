package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// javaErrorPattern defines a Java error pattern with its type and suggestion
type javaErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Java error patterns - comprehensive list for javac output
var javaErrorPatterns = []javaErrorPattern{
	// Symbol errors
	{regexp.MustCompile(`error: cannot find symbol`), domain.ErrorTypeImport, "Add missing import or declaration"},
	{regexp.MustCompile(`symbol:\s+class (\w+)`), domain.ErrorTypeImport, "Import or declare class '%s'"},
	{regexp.MustCompile(`symbol:\s+variable (\w+)`), domain.ErrorTypeImport, "Declare variable '%s' before use"},
	{regexp.MustCompile(`symbol:\s+method (\w+)`), domain.ErrorTypeImport, "Import class containing method '%s' or declare it"},
	// Package/Import errors
	{regexp.MustCompile(`error: package (\S+) does not exist`), domain.ErrorTypeDependency, "Add dependency for package '%s'"},
	{regexp.MustCompile(`error: cannot access (\w+)`), domain.ErrorTypeDependency, "Check classpath for '%s'"},
	// Type errors
	{regexp.MustCompile(`error: incompatible types: (.+)`), domain.ErrorTypeTypeCheck, "Fix type mismatch: %s"},
	{regexp.MustCompile(`error: method (\w+) in class (\w+) cannot be applied`), domain.ErrorTypeTypeCheck, "Check method '%s' arguments in class '%s'"},
	{regexp.MustCompile(`error: no suitable method found for (\w+)`), domain.ErrorTypeTypeCheck, "No matching method '%s' - check argument types"},
	{regexp.MustCompile(`error: no suitable constructor found`), domain.ErrorTypeTypeCheck, "No matching constructor - check argument types"},
	{regexp.MustCompile(`error: (\w+) has private access in (\w+)`), domain.ErrorTypeTypeCheck, "'%s' is private in '%s' - use public accessor"},
	{regexp.MustCompile(`error: non-static .+ cannot be referenced from a static context`), domain.ErrorTypeTypeCheck, "Cannot access instance member from static context"},
	// Syntax errors
	{regexp.MustCompile(`error: ';' expected`), domain.ErrorTypeSyntax, "Add missing semicolon"},
	{regexp.MustCompile(`error: '\)' expected`), domain.ErrorTypeSyntax, "Add missing closing parenthesis"},
	{regexp.MustCompile(`error: '\}' expected`), domain.ErrorTypeSyntax, "Add missing closing brace"},
	{regexp.MustCompile(`error: reached end of file while parsing`), domain.ErrorTypeSyntax, "Check for unclosed braces or parentheses"},
	{regexp.MustCompile(`error: illegal start of`), domain.ErrorTypeSyntax, "Fix syntax - illegal start of expression/type"},
	// Compilation errors
	{regexp.MustCompile(`error: variable (\w+) might not have been initialized`), domain.ErrorTypeCompilation, "Initialize variable '%s' before use"},
	{regexp.MustCompile(`error: unreported exception (\S+)`), domain.ErrorTypeCompilation, "Handle or declare exception '%s'"},
	{regexp.MustCompile(`error: missing return statement`), domain.ErrorTypeCompilation, "Add return statement"},
	// Generic error fallback
	{regexp.MustCompile(`error: (.+)`), domain.ErrorTypeCompilation, "Java error: %s"},
}

// javaLocationPattern matches Java compiler error location format: File.java:line: error:
var javaLocationPattern = regexp.MustCompile(`([^:\s]+\.java):(\d+):`)

// JavaErrorAnalyzer analyzes Java compiler error output
type JavaErrorAnalyzer struct{}

// NewJavaErrorAnalyzer creates a new Java error analyzer
func NewJavaErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &JavaErrorAnalyzer{}
}

// AnalyzeError analyzes Java error output and extracts details
func (ja *JavaErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "javac",
		Suggestions: make([]string, 0),
	}

	// Extract file location from Java error format
	ja.extractLocationInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range javaErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := ja.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = ja.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Java compiler output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Java-specific correction suggestions
func (ja *JavaErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	if strings.Contains(msg, "cannot find symbol") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import statement or declare the symbol",
		})
	}

	if strings.Contains(msg, "package") && strings.Contains(msg, "does not exist") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing dependency to build configuration (pom.xml/build.gradle)",
		})
	}

	if strings.Contains(msg, "expected") || strings.Contains(msg, "illegal start") ||
		strings.Contains(msg, "reached end of file") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix Java syntax error",
		})
	}

	if strings.Contains(msg, "incompatible types") || strings.Contains(msg, "cannot be applied") ||
		strings.Contains(msg, "no suitable") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch or method signature",
		})
	}

	if strings.Contains(msg, "private access") || strings.Contains(msg, "static context") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix access modifier or static/instance context",
		})
	}

	if strings.Contains(msg, "unreported exception") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Add try-catch block or throws declaration",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Java output
func (ja *JavaErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range javaErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return ja.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (ja *JavaErrorAnalyzer) GetLanguage() string {
	return langJava
}

// extractLocationInfo extracts file and line information from Java error output
func (ja *JavaErrorAnalyzer) extractLocationInfo(errorOutput string, details *domain.ErrorDetails) {
	if matches := javaLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 3 {
		details.SourceFile = matches[1]
		if line, err := parseIntSafe(matches[2]); err == nil {
			details.LineNumber = line
		}
	}
}

// formatSuggestion formats a suggestion string with captured groups
func (ja *JavaErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (ja *JavaErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "cannot find symbol") || strings.Contains(errorLower, "symbol:"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "package") && strings.Contains(errorLower, "does not exist"):
		return domain.ErrorTypeDependency
	case strings.Contains(errorLower, "incompatible types") || strings.Contains(errorLower, "cannot be applied") ||
		strings.Contains(errorLower, "no suitable") || strings.Contains(errorLower, "private access"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "expected") || strings.Contains(errorLower, "illegal start") ||
		strings.Contains(errorLower, "reached end of file"):
		return domain.ErrorTypeSyntax
	default:
		return domain.ErrorTypeCompilation
	}
}
