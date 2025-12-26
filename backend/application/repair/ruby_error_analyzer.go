package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// rubyErrorPattern defines a Ruby error pattern
type rubyErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// Ruby traceback pattern matches Ruby error format
var rubyTracebackPattern = regexp.MustCompile(`([^:]+):(\d+):(?:in\s+` + "`" + `([^']+)` + "`" + `)?`)

// Ruby error patterns - comprehensive list for Ruby error classification
var rubyErrorPatterns = []rubyErrorPattern{
	// Syntax errors
	{regexp.MustCompile(`SyntaxError:\s*(.+)`), domain.ErrorTypeSyntax, "Fix syntax error: %s"},
	{regexp.MustCompile(`syntax error,\s*(.+)`), domain.ErrorTypeSyntax, "Fix syntax error: %s"},
	{regexp.MustCompile(`unexpected\s+(\w+),\s*expecting\s+(.+)`), domain.ErrorTypeSyntax, "Unexpected %s, expecting %s"},

	// Load/Require errors
	{regexp.MustCompile(`LoadError:\s*cannot load such file\s*--\s*([^\s]+)`), domain.ErrorTypeDependency, "Install missing gem: gem install %s"},
	{regexp.MustCompile(`LoadError:\s*(.+)`), domain.ErrorTypeDependency, "Fix load error: %s"},
	{regexp.MustCompile(`cannot load such file\s*--\s*([^\s]+)`), domain.ErrorTypeDependency, "Install missing gem '%s' or check require path"},

	// Name errors
	{regexp.MustCompile(`NameError:\s*uninitialized constant\s+(\w+)`), domain.ErrorTypeImport, "Define or require constant '%s'"},
	{regexp.MustCompile(`NameError:\s*undefined local variable or method\s+` + "`" + `(\w+)` + "`"), domain.ErrorTypeImport, "Define variable or method '%s' before use"},
	{regexp.MustCompile(`NameError:\s*(.+)`), domain.ErrorTypeImport, "Fix name error: %s"},

	// NoMethod errors
	{regexp.MustCompile(`NoMethodError:\s*undefined method\s+` + "`" + `(\w+)` + "`" + `\s+for\s+(.+)`), domain.ErrorTypeTypeCheck, "Method '%s' not defined for %s"},
	{regexp.MustCompile(`NoMethodError:\s*private method\s+` + "`" + `(\w+)` + "`" + `\s+called`), domain.ErrorTypeTypeCheck, "Method '%s' is private - use public method or change visibility"},
	{regexp.MustCompile(`NoMethodError:\s*(.+)`), domain.ErrorTypeTypeCheck, "Fix method error: %s"},

	// Type errors
	{regexp.MustCompile(`TypeError:\s*no implicit conversion of\s+(\w+)\s+into\s+(\w+)`), domain.ErrorTypeTypeCheck, "Cannot convert %s to %s - check types"},
	{regexp.MustCompile(`TypeError:\s*(.+)\s+can't be coerced into\s+(\w+)`), domain.ErrorTypeTypeCheck, "%s cannot be coerced into %s"},
	{regexp.MustCompile(`TypeError:\s*(.+)`), domain.ErrorTypeTypeCheck, "Fix type error: %s"},

	// Argument errors
	{regexp.MustCompile(`ArgumentError:\s*wrong number of arguments\s*\(given\s+(\d+),\s*expected\s+(\d+)\)`), domain.ErrorTypeTypeCheck, "Wrong argument count: given %s, expected %s"},
	{regexp.MustCompile(`ArgumentError:\s*(.+)`), domain.ErrorTypeTypeCheck, "Fix argument error: %s"},

	// Key/Index errors
	{regexp.MustCompile(`KeyError:\s*key not found:\s*(.+)`), domain.ErrorTypeCompilation, "Key not found: %s"},
	{regexp.MustCompile(`IndexError:\s*(.+)`), domain.ErrorTypeCompilation, "Fix index error: %s"},

	// Runtime errors
	{regexp.MustCompile(`RuntimeError:\s*(.+)`), domain.ErrorTypeCompilation, "Fix runtime error: %s"},
	{regexp.MustCompile(`ZeroDivisionError:\s*(.+)`), domain.ErrorTypeCompilation, "Fix division by zero: %s"},
	{regexp.MustCompile(`SystemStackError:\s*(.+)`), domain.ErrorTypeCompilation, "Fix stack overflow (possible infinite recursion): %s"},
	{regexp.MustCompile(`Errno::ENOENT:\s*No such file or directory`), domain.ErrorTypeCompilation, "File not found - check file path"},

	// RSpec test failures
	{regexp.MustCompile(`RSpec::Expectations::ExpectationNotMetError`), domain.ErrorTypeTesting, "RSpec expectation failed - check test assertions"},
	{regexp.MustCompile(`expected:\s*(.+)\s*got:\s*(.+)`), domain.ErrorTypeTesting, "Expected %s but got %s"},
	{regexp.MustCompile(`Failure/Error:\s*(.+)`), domain.ErrorTypeTesting, "Test failure: %s"},

	// Minitest failures
	{regexp.MustCompile(`Minitest::Assertion:\s*(.+)`), domain.ErrorTypeTesting, "Minitest assertion failed: %s"},
	{regexp.MustCompile(`Expected\s+(.+)\s+to\s+(.+)`), domain.ErrorTypeTesting, "Expected %s to %s"},

	// RuboCop linting
	{regexp.MustCompile(`Style/(\w+):\s*(.+)`), domain.ErrorTypeLinting, "Style issue (%s): %s"},
	{regexp.MustCompile(`Layout/(\w+):\s*(.+)`), domain.ErrorTypeLinting, "Layout issue (%s): %s"},
	{regexp.MustCompile(`Lint/(\w+):\s*(.+)`), domain.ErrorTypeLinting, "Lint issue (%s): %s"},
	{regexp.MustCompile(`Metrics/(\w+):\s*(.+)`), domain.ErrorTypeLinting, "Metrics issue (%s): %s"},

	// Bundler errors
	{regexp.MustCompile(`Bundler::GemNotFound`), domain.ErrorTypeDependency, "Run 'bundle install' to install missing gems"},
	{regexp.MustCompile(`Could not find gem\s+'([^']+)'`), domain.ErrorTypeDependency, "Add gem '%s' to Gemfile and run 'bundle install'"},
}

// RubyErrorAnalyzer analyzes Ruby-specific errors
type RubyErrorAnalyzer struct{}

// NewRubyErrorAnalyzer creates a new Ruby error analyzer
func NewRubyErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &RubyErrorAnalyzer{}
}

// AnalyzeError analyzes Ruby error output and extracts details
func (r *RubyErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "ruby",
		Suggestions: make([]string, 0),
	}

	// Extract file location from traceback
	r.extractTracebackInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range rubyErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := r.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = r.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Ruby error output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Ruby-specific correction suggestions
func (r *RubyErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// LoadError - suggest gem install or bundle install
	if strings.Contains(msg, "LoadError") || strings.Contains(msg, "cannot load such file") {
		gemName := r.extractGemName(msg)
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: fmt.Sprintf("Install missing gem: gem install %s or add to Gemfile", gemName),
		})
	}

	// NameError - suggest require or definition
	if strings.Contains(msg, "NameError") || strings.Contains(msg, "uninitialized constant") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing require statement or define the constant/variable",
		})
	}

	// SyntaxError - suggest syntax fix
	if strings.Contains(msg, "SyntaxError") || strings.Contains(msg, "syntax error") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix Ruby syntax error - check for missing end, brackets, or quotes",
		})
	}

	// NoMethodError / TypeError - suggest type fix
	if strings.Contains(msg, "NoMethodError") || strings.Contains(msg, "TypeError") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type error - check method exists and argument types",
		})
	}

	// ArgumentError - suggest argument fix
	if strings.Contains(msg, "ArgumentError") || strings.Contains(msg, "wrong number of arguments") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix argument count - check method signature",
		})
	}

	// RuboCop linting issues - suggest auto-correct
	if strings.Contains(msg, "Style/") || strings.Contains(msg, "Layout/") ||
		strings.Contains(msg, "Lint/") || strings.Contains(msg, "Metrics/") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFormatCode,
			Target:      errDetails.SourceFile,
			Description: "Run 'rubocop -a' to auto-correct style issues",
		})
	}

	// Bundler errors - suggest bundle install
	if strings.Contains(msg, "Bundler") || strings.Contains(msg, "Could not find gem") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      "Gemfile",
			Description: "Run 'bundle install' to install missing dependencies",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Ruby output
func (r *RubyErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range rubyErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return r.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (r *RubyErrorAnalyzer) GetLanguage() string {
	return langRuby
}

// extractTracebackInfo extracts file and line information from Ruby traceback
func (r *RubyErrorAnalyzer) extractTracebackInfo(errorOutput string, details *domain.ErrorDetails) {
	matches := rubyTracebackPattern.FindAllStringSubmatch(errorOutput, -1)
	if len(matches) > 0 {
		// Use the first match (usually the most relevant)
		firstMatch := matches[0]
		if len(firstMatch) >= 2 {
			details.SourceFile = firstMatch[1]
		}
		if len(firstMatch) >= 3 {
			if line, err := parseIntSafe(firstMatch[2]); err == nil {
				details.LineNumber = line
			}
		}
	}
}

// formatSuggestion formats a suggestion string with captured groups
func (r *RubyErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (r *RubyErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "syntaxerror") || strings.Contains(errorLower, "syntax error"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "loaderror") || strings.Contains(errorLower, "cannot load such file"):
		return domain.ErrorTypeDependency
	case strings.Contains(errorLower, "nameerror") || strings.Contains(errorLower, "uninitialized constant"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "nomethoderror") || strings.Contains(errorLower, "typeerror"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "argumenterror"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "rspec") || strings.Contains(errorLower, "minitest") ||
		strings.Contains(errorLower, "failure/error") || strings.Contains(errorLower, "expected"):
		return domain.ErrorTypeTesting
	case strings.Contains(errorLower, "style/") || strings.Contains(errorLower, "layout/") ||
		strings.Contains(errorLower, "lint/") || strings.Contains(errorLower, "rubocop"):
		return domain.ErrorTypeLinting
	case strings.Contains(errorLower, "bundler") || strings.Contains(errorLower, "could not find gem"):
		return domain.ErrorTypeDependency
	default:
		return domain.ErrorTypeCompilation
	}
}

// extractGemName extracts gem name from load error message
func (r *RubyErrorAnalyzer) extractGemName(msg string) string {
	// Pattern: cannot load such file -- gem_name
	pattern := regexp.MustCompile(`cannot load such file\s*--\s*([^\s\n]+)`)
	if matches := pattern.FindStringSubmatch(msg); len(matches) > 1 {
		gemName := matches[1]
		// Remove path components if present
		if idx := strings.LastIndex(gemName, "/"); idx >= 0 {
			gemName = gemName[idx+1:]
		}
		return gemName
	}

	// Pattern: Could not find gem 'gem_name'
	pattern2 := regexp.MustCompile(`Could not find gem\s+'([^']+)'`)
	if matches := pattern2.FindStringSubmatch(msg); len(matches) > 1 {
		return matches[1]
	}

	return "unknown-gem"
}
