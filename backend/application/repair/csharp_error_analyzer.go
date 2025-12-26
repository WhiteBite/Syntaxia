package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// csharpErrorPattern defines a C# error pattern with its type and suggestion
type csharpErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	suggestion string
}

// C# error patterns - comprehensive list for dotnet/csc output
var csharpErrorPatterns = []csharpErrorPattern{
	// CS0103: The name 'x' does not exist in the current context
	{regexp.MustCompile(`CS0103: The name '(\w+)' does not exist`), domain.ErrorTypeImport, "Add using directive or declare '%s'"},

	// CS0246: The type or namespace name 'x' could not be found
	{regexp.MustCompile(`CS0246: The type or namespace name '(\w+)' could not be found`), domain.ErrorTypeImport, "Add using directive for '%s' or install NuGet package"},

	// CS0029: Cannot implicitly convert type 'x' to 'y'
	{regexp.MustCompile(`CS0029: Cannot implicitly convert type '([^']+)' to '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix type conversion from '%s' to '%s'"},

	// CS1002: ; expected
	{regexp.MustCompile(`CS1002: ; expected`), domain.ErrorTypeSyntax, "Add missing semicolon"},

	// CS0117: 'x' does not contain a definition for 'y'
	{regexp.MustCompile(`CS0117: '([^']+)' does not contain a definition for '([^']+)'`), domain.ErrorTypeTypeCheck, "'%s' has no member '%s' - check spelling or type"},

	// CS0234: The type or namespace name 'x' does not exist in the namespace 'y'
	{regexp.MustCompile(`CS0234: The type or namespace name '(\w+)' does not exist in the namespace '([^']+)'`), domain.ErrorTypeImport, "Type '%s' not found in namespace '%s' - check using directives"},

	// CS0019: Operator 'x' cannot be applied to operands of type 'y' and 'z'
	{regexp.MustCompile(`CS0019: Operator '([^']+)' cannot be applied to operands of type '([^']+)' and '([^']+)'`), domain.ErrorTypeTypeCheck, "Operator '%s' incompatible with types '%s' and '%s'"},

	// CS0161: not all code paths return a value
	{regexp.MustCompile(`CS0161: .+ not all code paths return a value`), domain.ErrorTypeCompilation, "Add return statement for all code paths"},

	// CS0168: The variable 'x' is declared but never used
	{regexp.MustCompile(`CS0168: The variable '(\w+)' is declared but never used`), domain.ErrorTypeLinting, "Remove unused variable '%s'"},

	// CS0219: The variable 'x' is assigned but its value is never used
	{regexp.MustCompile(`CS0219: The variable '(\w+)' is assigned but its value is never used`), domain.ErrorTypeLinting, "Remove unused variable '%s' or use its value"},

	// CS8600: Converting null literal or possible null value to non-nullable type
	{regexp.MustCompile(`CS8600: Converting null literal or possible null value to non-nullable type`), domain.ErrorTypeTypeCheck, "Handle nullable value - use null check or nullable type"},

	// CS8602: Dereference of a possibly null reference
	{regexp.MustCompile(`CS8602: Dereference of a possibly null reference`), domain.ErrorTypeTypeCheck, "Add null check before dereferencing"},

	// CS8604: Possible null reference argument
	{regexp.MustCompile(`CS8604: Possible null reference argument`), domain.ErrorTypeTypeCheck, "Argument may be null - add null check or use nullable parameter"},

	// CS0120: An object reference is required for the non-static field/method/property
	{regexp.MustCompile(`CS0120: An object reference is required for the non-static`), domain.ErrorTypeTypeCheck, "Cannot access instance member from static context - create instance or make member static"},

	// CS0176: Member 'x' cannot be accessed with an instance reference
	{regexp.MustCompile(`CS0176: Member '([^']+)' cannot be accessed with an instance reference`), domain.ErrorTypeTypeCheck, "Access static member '%s' using type name, not instance"},

	// CS1061: 'x' does not contain a definition for 'y'
	{regexp.MustCompile(`CS1061: '([^']+)' does not contain a definition for '([^']+)'`), domain.ErrorTypeTypeCheck, "'%s' has no member '%s' - check type or add extension method"},

	// CS0428: Cannot convert method group 'x' to non-delegate type
	{regexp.MustCompile(`CS0428: Cannot convert method group '(\w+)' to non-delegate type`), domain.ErrorTypeTypeCheck, "Add parentheses to invoke method '%s' or use delegate"},

	// CS0535: 'x' does not implement interface member 'y'
	{regexp.MustCompile(`CS0535: '([^']+)' does not implement interface member '([^']+)'`), domain.ErrorTypeTypeCheck, "Implement missing interface member '%s' in '%s'"},

	// CS0534: 'x' does not implement inherited abstract member 'y'
	{regexp.MustCompile(`CS0534: '([^']+)' does not implement inherited abstract member '([^']+)'`), domain.ErrorTypeTypeCheck, "Implement abstract member '%s' in '%s'"},

	// CS1503: Argument x: cannot convert from 'y' to 'z'
	{regexp.MustCompile(`CS1503: Argument \d+: cannot convert from '([^']+)' to '([^']+)'`), domain.ErrorTypeTypeCheck, "Fix argument type - cannot convert '%s' to '%s'"},

	// CS0266: Cannot implicitly convert type 'x' to 'y'. An explicit conversion exists
	{regexp.MustCompile(`CS0266: Cannot implicitly convert type '([^']+)' to '([^']+)'\. An explicit conversion exists`), domain.ErrorTypeTypeCheck, "Add explicit cast from '%s' to '%s'"},

	// Generic error fallback
	{regexp.MustCompile(`error CS\d+: (.+)`), domain.ErrorTypeCompilation, "C# error: %s"},
}

// csharpLocationPattern matches C# compiler error location format: File.cs(10,5): error CS0103:
var csharpLocationPattern = regexp.MustCompile(`([^(]+)\((\d+),(\d+)\):`)

// CSharpErrorAnalyzer analyzes C# compiler error output
type CSharpErrorAnalyzer struct{}

// NewCSharpErrorAnalyzer creates a new C# error analyzer
func NewCSharpErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &CSharpErrorAnalyzer{}
}

// AnalyzeError analyzes C# error output and extracts details
func (c *CSharpErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "dotnet",
		Suggestions: make([]string, 0),
	}

	// Extract file location from C# error format
	c.extractLocationInfo(errorOutput, details)

	// Match against known error patterns
	for _, pattern := range csharpErrorPatterns {
		if matches := pattern.pattern.FindStringSubmatch(errorOutput); matches != nil {
			details.ErrorType = pattern.errType
			suggestion := c.formatSuggestion(pattern.suggestion, matches)
			details.Suggestions = append(details.Suggestions, suggestion)
			break
		}
	}

	// Fallback classification if no pattern matched
	if details.ErrorType == "" {
		details.ErrorType = c.classifyByKeywords(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check C# compiler output for details")
		}
	}

	return details, nil
}

// SuggestCorrections provides C#-specific correction suggestions
func (c *CSharpErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Name/type not found - suggest using directive
	if strings.Contains(msg, "CS0103") || strings.Contains(msg, "CS0246") ||
		strings.Contains(msg, "does not exist") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing using directive or install NuGet package",
		})
	}

	// Type conversion errors
	if strings.Contains(msg, "CS0029") || strings.Contains(msg, "CS0266") ||
		strings.Contains(msg, "cannot convert") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type conversion - add explicit cast or change type",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "CS1002") || strings.Contains(msg, "expected") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	}

	// Nullable reference warnings
	if strings.Contains(msg, "CS8600") || strings.Contains(msg, "CS8602") ||
		strings.Contains(msg, "CS8604") || strings.Contains(msg, "null") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Handle nullable reference - add null check or use nullable type",
		})
	}

	// Unused variable warnings
	if strings.Contains(msg, "CS0168") || strings.Contains(msg, "CS0219") ||
		strings.Contains(msg, "never used") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionRemoveCode,
			Target:      errDetails.SourceFile,
			Description: "Remove unused variable or use its value",
		})
	}

	// Interface/abstract implementation
	if strings.Contains(msg, "CS0535") || strings.Contains(msg, "CS0534") ||
		strings.Contains(msg, "does not implement") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Implement missing interface or abstract members",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from C# output
func (c *CSharpErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	for _, pattern := range csharpErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}

	return c.classifyByKeywords(errorOutput)
}

// GetLanguage returns the supported language
func (c *CSharpErrorAnalyzer) GetLanguage() string {
	return langCSharp
}

// extractLocationInfo extracts file and line information from C# error output
func (c *CSharpErrorAnalyzer) extractLocationInfo(errorOutput string, details *domain.ErrorDetails) {
	if matches := csharpLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 4 {
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
func (c *CSharpErrorAnalyzer) formatSuggestion(template string, matches []string) string {
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
func (c *CSharpErrorAnalyzer) classifyByKeywords(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)

	switch {
	case strings.Contains(errorLower, "cs0103") ||
		strings.Contains(errorLower, "cs0246") ||
		strings.Contains(errorLower, "cs0234") ||
		strings.Contains(errorLower, "does not exist"):
		return domain.ErrorTypeImport

	case strings.Contains(errorLower, "cs0029") ||
		strings.Contains(errorLower, "cs0266") ||
		strings.Contains(errorLower, "cs1503") ||
		strings.Contains(errorLower, "cannot convert") ||
		strings.Contains(errorLower, "cs0117") ||
		strings.Contains(errorLower, "cs1061"):
		return domain.ErrorTypeTypeCheck

	case strings.Contains(errorLower, "cs1002") ||
		strings.Contains(errorLower, "expected"):
		return domain.ErrorTypeSyntax

	case strings.Contains(errorLower, "cs0168") ||
		strings.Contains(errorLower, "cs0219") ||
		strings.Contains(errorLower, "never used"):
		return domain.ErrorTypeLinting

	default:
		return domain.ErrorTypeCompilation
	}
}
