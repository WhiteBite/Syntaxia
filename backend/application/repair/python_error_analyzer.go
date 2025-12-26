package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// pythonErrorPattern defines a Python error pattern
type pythonErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// pythonTracebackPattern matches Python traceback format
var pythonTracebackPattern = regexp.MustCompile(`File "([^"]+)", line (\d+)(?:, in (\w+))?`)

// Python error patterns - comprehensive list for Python error classification
var pythonErrorPatterns = []pythonErrorPattern{
	// Syntax errors
	{regexp.MustCompile(`SyntaxError: (.+)`), domain.ErrorTypeSyntax, "Fix syntax error: %s"},
	{regexp.MustCompile(`IndentationError: (.+)`), domain.ErrorTypeSyntax, "Fix indentation: %s"},
	{regexp.MustCompile(`TabError: (.+)`), domain.ErrorTypeSyntax, "Fix tab/space mixing: %s"},
	// Import errors
	{regexp.MustCompile(`ImportError: No module named '([^']+)'`), domain.ErrorTypeDependency, "Install module '%s' with pip"},
	{regexp.MustCompile(`ModuleNotFoundError: No module named '([^']+)'`), domain.ErrorTypeDependency, "Install module '%s' with pip"},
	{regexp.MustCompile(`ImportError: cannot import name '(\w+)'`), domain.ErrorTypeImport, "Check import for '%s'"},
	{regexp.MustCompile(`ImportError: (.+)`), domain.ErrorTypeImport, "Fix import error: %s"},
	// Name errors
	{regexp.MustCompile(`NameError: name '(\w+)' is not defined`), domain.ErrorTypeImport, "Define or import '%s' before use"},
	// Type errors
	{regexp.MustCompile(`TypeError: (.+) got an unexpected keyword argument '(\w+)'`), domain.ErrorTypeTypeCheck, "Remove unexpected argument '%s'"},
	{regexp.MustCompile(`TypeError: (.+) missing (\d+) required positional argument`), domain.ErrorTypeTypeCheck, "Add missing required arguments"},
	{regexp.MustCompile(`TypeError: (.+)`), domain.ErrorTypeTypeCheck, "Fix type error: %s"},
	// Attribute errors
	{regexp.MustCompile(`AttributeError: '(\w+)' object has no attribute '(\w+)'`), domain.ErrorTypeTypeCheck, "Object '%s' has no attribute '%s'"},
	{regexp.MustCompile(`AttributeError: module '([^']+)' has no attribute '(\w+)'`), domain.ErrorTypeTypeCheck, "Module '%s' has no attribute '%s'"},
	{regexp.MustCompile(`AttributeError: (.+)`), domain.ErrorTypeTypeCheck, "Fix attribute error: %s"},
	// Value/Key/Index errors
	{regexp.MustCompile(`ValueError: (.+)`), domain.ErrorTypeCompilation, "Fix value error: %s"},
	{regexp.MustCompile(`KeyError: '([^']+)'`), domain.ErrorTypeCompilation, "Key '%s' not found"},
	{regexp.MustCompile(`IndexError: (.+)`), domain.ErrorTypeCompilation, "Fix index error: %s"},
	// Runtime errors
	{regexp.MustCompile(`ZeroDivisionError: (.+)`), domain.ErrorTypeCompilation, "Fix division by zero: %s"},
	{regexp.MustCompile(`FileNotFoundError: (.+)`), domain.ErrorTypeCompilation, "File not found: %s"},
	{regexp.MustCompile(`RuntimeError: (.+)`), domain.ErrorTypeCompilation, "Fix runtime error: %s"},
	{regexp.MustCompile(`RecursionError: (.+)`), domain.ErrorTypeCompilation, "Fix recursion error: %s"},
	// Assertion errors (testing)
	{regexp.MustCompile(`AssertionError: (.+)`), domain.ErrorTypeTesting, "Assertion failed: %s"},
	{regexp.MustCompile(`AssertionError`), domain.ErrorTypeTesting, "Assertion failed - check test conditions"},
	// pytest specific errors
	{regexp.MustCompile(`FAILED (.+)::([\w_]+)`), domain.ErrorTypeTesting, "Test failed: %s::%s"},
	{regexp.MustCompile(`fixture '(\w+)' not found`), domain.ErrorTypeTesting, "Missing pytest fixture '%s'"},
	// Type hints / mypy errors
	{regexp.MustCompile(`error: Incompatible types in assignment`), domain.ErrorTypeTypeCheck, "Fix type annotation"},
	{regexp.MustCompile(`error: Argument .+ has incompatible type`), domain.ErrorTypeTypeCheck, "Fix argument type annotation"},
	// Linting errors (pylint, flake8, ruff)
	{regexp.MustCompile(`undefined name '(\w+)'`), domain.ErrorTypeImport, "Undefined name '%s' - add import"},
	{regexp.MustCompile(`'(\w+)' imported but unused`), domain.ErrorTypeLinting, "Remove unused import '%s'"},
	{regexp.MustCompile(`local variable '(\w+)' is assigned .* but never used`), domain.ErrorTypeLinting, "Remove unused variable '%s'"},
}

// PythonErrorAnalyzer analyzes Python-specific errors
type PythonErrorAnalyzer struct{}

// NewPythonErrorAnalyzer creates a new Python error analyzer
func NewPythonErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &PythonErrorAnalyzer{}
}

// AnalyzeError analyzes Python error output and extracts details
func (p *PythonErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "python",
		Suggestions: make([]string, 0),
	}

	// Extract file location from traceback
	p.extractTracebackInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range pythonErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := p.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = p.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check Python error output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides Python-specific correction suggestions
func (p *PythonErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Import/Module errors - suggest pip install
	if strings.Contains(msg, "ModuleNotFoundError") || strings.Contains(msg, "No module named") {
		moduleName := p.extractModuleName(msg)
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: fmt.Sprintf("Install missing module: pip install %s", moduleName),
		})
	}

	// NameError - suggest import or definition
	if strings.Contains(msg, "NameError") || strings.Contains(msg, "is not defined") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing import or define the variable/function",
		})
	}

	// Syntax errors - suggest formatting
	if strings.Contains(msg, "SyntaxError") || strings.Contains(msg, "IndentationError") || strings.Contains(msg, "TabError") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix Python syntax error - check indentation and syntax",
		})
	}

	// Type errors
	if strings.Contains(msg, "TypeError") || strings.Contains(msg, "AttributeError") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type error - check argument types and object attributes",
		})
	}

	// Linting issues - suggest formatter
	if strings.Contains(msg, "imported but unused") || strings.Contains(msg, "never used") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionRemoveCode,
			Target:      errDetails.SourceFile,
			Description: "Remove unused import/variable",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Python output
func (p *PythonErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range pythonErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return p.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (p *PythonErrorAnalyzer) GetLanguage() string {
	return langPython
}

// extractTracebackInfo extracts file and line information from Python traceback
func (p *PythonErrorAnalyzer) extractTracebackInfo(errorOutput string, details *domain.ErrorDetails) {
	matches := pythonTracebackPattern.FindAllStringSubmatch(errorOutput, -1)
	if len(matches) > 0 {
		lastMatch := matches[len(matches)-1]
		if len(lastMatch) >= 2 {
			details.SourceFile = lastMatch[1]
		}
		if len(lastMatch) >= 3 {
			if line, err := parseIntSafe(lastMatch[2]); err == nil {
				details.LineNumber = line
			}
		}
	}
}

// formatSuggestion formats a suggestion string with captured groups
func (p *PythonErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (p *PythonErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "syntaxerror") || strings.Contains(errorLower, "indentationerror"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "modulenotfounderror") || strings.Contains(errorLower, "no module named"):
		return domain.ErrorTypeDependency
	case strings.Contains(errorLower, "importerror") || strings.Contains(errorLower, "nameerror"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "typeerror") || strings.Contains(errorLower, "attributeerror"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "assertionerror") || strings.Contains(errorLower, "pytest"):
		return domain.ErrorTypeTesting
	case strings.Contains(errorLower, "unused"):
		return domain.ErrorTypeLinting
	default:
		return domain.ErrorTypeCompilation
	}
}

// extractModuleName extracts module name from import error message
func (p *PythonErrorAnalyzer) extractModuleName(msg string) string {
	pattern := regexp.MustCompile(`No module named '([^']+)'`)
	if matches := pattern.FindStringSubmatch(msg); len(matches) > 1 {
		moduleName := matches[1]
		if idx := strings.Index(moduleName, "."); idx > 0 {
			return moduleName[:idx]
		}
		return moduleName
	}
	return "unknown-module"
}
