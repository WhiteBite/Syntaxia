package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// phpErrorPattern defines a PHP error pattern with its type and suggestion
type phpErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// PHP error patterns - comprehensive list for PHP error classification
var phpErrorPatterns = []phpErrorPattern{
	// Parse/Syntax errors
	{regexp.MustCompile(`Parse error: (.+) in (.+) on line (\d+)`), domain.ErrorTypeSyntax, "Fix syntax error: %s"},
	{regexp.MustCompile(`syntax error, unexpected (.+)`), domain.ErrorTypeSyntax, "Fix unexpected %s"},
	{regexp.MustCompile(`syntax error, expecting (.+)`), domain.ErrorTypeSyntax, "Add expected %s"},
	// Fatal errors
	{regexp.MustCompile(`Fatal error: (.+) in (.+) on line (\d+)`), domain.ErrorTypeCompilation, "Fix fatal error: %s"},
	{regexp.MustCompile(`Fatal error: Uncaught Error: (.+)`), domain.ErrorTypeCompilation, "Handle error: %s"},
	{regexp.MustCompile(`Fatal error: Class '([^']+)' not found`), domain.ErrorTypeImport, "Import or autoload class '%s'"},
	{regexp.MustCompile(`Fatal error: Interface '([^']+)' not found`), domain.ErrorTypeImport, "Import interface '%s'"},
	{regexp.MustCompile(`Fatal error: Trait '([^']+)' not found`), domain.ErrorTypeImport, "Import trait '%s'"},
	// Undefined errors
	{regexp.MustCompile(`Undefined variable:? \$?(\w+)`), domain.ErrorTypeImport, "Define variable '$%s' before use"},
	{regexp.MustCompile(`Undefined constant "?(\w+)"?`), domain.ErrorTypeImport, "Define constant '%s'"},
	{regexp.MustCompile(`Call to undefined function (\w+)\(\)`), domain.ErrorTypeImport, "Import or define function '%s'"},
	{regexp.MustCompile(`Call to undefined method ([^:]+)::(\w+)\(\)`), domain.ErrorTypeImport, "Define method '%s' in class '%s'"},
	// Type errors
	{regexp.MustCompile(`TypeError: (.+)`), domain.ErrorTypeTypeCheck, "Fix type error: %s"},
	{regexp.MustCompile(`Argument (\d+) .+ must be of type (\w+), (\w+) given`), domain.ErrorTypeTypeCheck, "Argument %s expects %s, got %s"},
	{regexp.MustCompile(`Return value .+ must be of type (\w+), (\w+) returned`), domain.ErrorTypeTypeCheck, "Return type must be %s, not %s"},
	{regexp.MustCompile(`Cannot assign (\w+) to .+ of type (\w+)`), domain.ErrorTypeTypeCheck, "Cannot assign %s to type %s"},
	// Warning errors
	{regexp.MustCompile(`Warning: (.+) in (.+) on line (\d+)`), domain.ErrorTypeLinting, "Warning: %s"},
	{regexp.MustCompile(`Warning: Undefined array key "?([^"]+)"?`), domain.ErrorTypeLinting, "Check array key '%s' exists"},
	{regexp.MustCompile(`Warning: Trying to access array offset on`), domain.ErrorTypeLinting, "Check variable is array before access"},
	// Deprecated warnings
	{regexp.MustCompile(`Deprecated: (.+) in (.+) on line (\d+)`), domain.ErrorTypeLinting, "Deprecated: %s"},
	// PHPUnit/Testing errors
	{regexp.MustCompile(`Failed asserting that (.+)`), domain.ErrorTypeTesting, "Assertion failed: %s"},
	{regexp.MustCompile(`PHPUnit\\Framework\\ExpectationFailedException`), domain.ErrorTypeTesting, "PHPUnit expectation failed"},
	{regexp.MustCompile(`Error: (.+) in (.+):(\d+)`), domain.ErrorTypeCompilation, "Error: %s"},
	// Composer/Autoload errors
	{regexp.MustCompile(`require\(([^)]+)\): Failed to open stream`), domain.ErrorTypeDependency, "File not found: %s - run composer install"},
	{regexp.MustCompile(`include\(([^)]+)\): Failed to open stream`), domain.ErrorTypeDependency, "File not found: %s"},
}

// phpLocationPattern matches PHP error location format
var phpLocationPattern = regexp.MustCompile(`in ([^\s]+\.php)(?: on line |:)(\d+)`)

// PHPErrorAnalyzer analyzes PHP-specific errors
type PHPErrorAnalyzer struct{}

// NewPHPErrorAnalyzer creates a new PHP error analyzer
func NewPHPErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &PHPErrorAnalyzer{}
}

// AnalyzeError analyzes PHP error output and extracts details
func (p *PHPErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "php",
		Suggestions: make([]string, 0),
	}

	// Extract file location from PHP error format
	p.extractLocationInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range phpErrorPatterns {
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
			details.Suggestions = append(details.Suggestions, "Check PHP error output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides PHP-specific correction suggestions
func (p *PHPErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Class/Interface/Trait not found - suggest autoload
	if strings.Contains(msg, "not found") && (strings.Contains(msg, "Class") ||
		strings.Contains(msg, "Interface") || strings.Contains(msg, "Trait")) {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add use statement or run composer dump-autoload",
		})
	}

	// Undefined variable/function
	if strings.Contains(msg, "Undefined") || strings.Contains(msg, "undefined") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Define the variable/function before use or add use statement",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "Parse error") || strings.Contains(msg, "syntax error") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix PHP syntax error - check brackets, semicolons, quotes",
		})
	}

	// Type errors
	if strings.Contains(msg, "TypeError") || strings.Contains(msg, "must be of type") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch - check argument and return types",
		})
	}

	// Dependency errors
	if strings.Contains(msg, "Failed to open stream") || strings.Contains(msg, "require") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      "composer.json",
			Description: "Run composer install or add missing dependency",
		})
	}

	// Deprecated warnings
	if strings.Contains(msg, "Deprecated") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Update deprecated code to use modern alternatives",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from PHP output
func (p *PHPErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range phpErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}
	return p.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (p *PHPErrorAnalyzer) GetLanguage() string {
	return langPHP
}

// extractLocationInfo extracts file and line information from PHP error output
func (p *PHPErrorAnalyzer) extractLocationInfo(errorOutput string, details *domain.ErrorDetails) {
	if matches := phpLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 3 {
		details.SourceFile = matches[1]
		if line, err := parseIntSafe(matches[2]); err == nil {
			details.LineNumber = line
		}
	}
}

// formatSuggestion formats a suggestion string with captured groups
func (p *PHPErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (p *PHPErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)
	switch {
	case strings.Contains(errorLower, "parse error") || strings.Contains(errorLower, "syntax error"):
		return domain.ErrorTypeSyntax
	case strings.Contains(errorLower, "not found") || strings.Contains(errorLower, "undefined"):
		return domain.ErrorTypeImport
	case strings.Contains(errorLower, "failed to open stream") || strings.Contains(errorLower, "composer"):
		return domain.ErrorTypeDependency
	case strings.Contains(errorLower, "typeerror") || strings.Contains(errorLower, "must be of type"):
		return domain.ErrorTypeTypeCheck
	case strings.Contains(errorLower, "phpunit") || strings.Contains(errorLower, "asserting"):
		return domain.ErrorTypeTesting
	case strings.Contains(errorLower, "warning") || strings.Contains(errorLower, "deprecated"):
		return domain.ErrorTypeLinting
	default:
		return domain.ErrorTypeCompilation
	}
}
