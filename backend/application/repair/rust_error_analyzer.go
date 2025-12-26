package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
)

// rustErrorPattern defines a Rust error pattern with its type and suggestion
type rustErrorPattern struct {
	pattern    *regexp.Regexp
	errType    domain.ErrorType
	category   string
	suggestion string
}

// Rust error patterns
var rustErrorPatterns = []rustErrorPattern{
	// Import/Resolution errors (E04xx)
	{regexp.MustCompile(`error\[E0425\]: cannot find value ` + "`" + `(\w+)` + "`"), domain.ErrorTypeImport, "undefined_value", "Declare or import '%s'"},
	{regexp.MustCompile(`error\[E0412\]: cannot find type ` + "`" + `(\w+)` + "`"), domain.ErrorTypeImport, "undefined_type", "Import or define type '%s'"},
	{regexp.MustCompile(`error\[E0433\]: failed to resolve: use of undeclared`), domain.ErrorTypeImport, "undeclared", "Add missing use statement"},
	{regexp.MustCompile(`error\[E0432\]: unresolved import ` + "`" + `([^` + "`" + `]+)` + "`"), domain.ErrorTypeDependency, "unresolved_import", "Fix import path '%s' or add dependency"},
	{regexp.MustCompile(`error\[E0463\]: can't find crate for ` + "`" + `(\w+)` + "`"), domain.ErrorTypeDependency, "missing_crate", "Add crate '%s' to Cargo.toml"},
	// Type errors (E03xx)
	{regexp.MustCompile(`error\[E0308\]: mismatched types`), domain.ErrorTypeTypeCheck, "type_mismatch", "Fix type mismatch - check expected vs actual types"},
	{regexp.MustCompile(`error\[E0277\]: .+ doesn't implement`), domain.ErrorTypeTypeCheck, "trait_not_impl", "Implement required trait or add trait bound"},
	{regexp.MustCompile(`error\[E0599\]: no method named ` + "`" + `(\w+)` + "`"), domain.ErrorTypeTypeCheck, "no_method", "Method '%s' not found - check type or import trait"},
	{regexp.MustCompile(`error\[E0609\]: no field ` + "`" + `(\w+)` + "`"), domain.ErrorTypeTypeCheck, "no_field", "Field '%s' not found - check struct definition"},
	// Borrow checker errors (E05xx)
	{regexp.MustCompile(`error\[E0382\]: borrow of moved value`), domain.ErrorTypeCompilation, "moved_value", "Value was moved - use clone() or restructure ownership"},
	{regexp.MustCompile(`error\[E0502\]: cannot borrow .+ as mutable`), domain.ErrorTypeCompilation, "borrow_conflict", "Cannot borrow as mutable while immutably borrowed"},
	{regexp.MustCompile(`error\[E0507\]: cannot move out of`), domain.ErrorTypeCompilation, "cannot_move", "Cannot move out - use clone() or reference"},
	// Lifetime errors (E06xx)
	{regexp.MustCompile(`error\[E0106\]: missing lifetime specifier`), domain.ErrorTypeTypeCheck, "missing_lifetime", "Add lifetime parameter"},
	{regexp.MustCompile(`error\[E0621\]: explicit lifetime required`), domain.ErrorTypeTypeCheck, "lifetime_required", "Add explicit lifetime annotation"},
	// Syntax errors
	{regexp.MustCompile(`error: expected .+, found`), domain.ErrorTypeSyntax, "expected_token", "Fix syntax - check for missing tokens"},
	{regexp.MustCompile(`error: unexpected token`), domain.ErrorTypeSyntax, "unexpected_token", "Remove or fix unexpected token"},
	{regexp.MustCompile(`error: cannot find macro ` + "`" + `(\w+)` + "`"), domain.ErrorTypeImport, "missing_macro", "Import macro '%s' or add macro_use"},
	// Generic errors (catch-all patterns - should be last)
	{regexp.MustCompile(`error\[E\d+\]: (.+)`), domain.ErrorTypeCompilation, "generic", "Check rustc error message for details"},
	{regexp.MustCompile(`error: (.+)`), domain.ErrorTypeCompilation, "error", "Check error message for details"},
}

// rustLocationPattern matches Rust error location format: --> src/main.rs:5:13
var rustLocationPattern = regexp.MustCompile(`-->\s*([^:]+):(\d+):(\d+)`)

// RustErrorAnalyzer analyzes Rust compiler error output
type RustErrorAnalyzer struct{}

// NewRustErrorAnalyzer creates a new Rust error analyzer
func NewRustErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &RustErrorAnalyzer{}
}

// AnalyzeError analyzes Rust error output and extracts details
func (r *RustErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "rustc",
		Suggestions: make([]string, 0),
	}

	// Extract location from Rust error format
	r.extractRustLocation(errorOutput, details)

	// Match against patterns
	for _, pattern := range rustErrorPatterns {
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
		details.ErrorType = r.ClassifyErrorType(errorOutput)
		if len(details.Suggestions) == 0 {
			details.Suggestions = append(details.Suggestions, "Check rustc error output for details")
		}
	}

	return details, nil
}

// extractRustLocation extracts file and line information from Rust error output
func (r *RustErrorAnalyzer) extractRustLocation(errorOutput string, details *domain.ErrorDetails) {
	if matches := rustLocationPattern.FindStringSubmatch(errorOutput); len(matches) >= 4 {
		details.SourceFile = matches[1]
		if line, err := parseIntSafe(matches[2]); err == nil {
			details.LineNumber = line
		}
		if col, err := parseIntSafe(matches[3]); err == nil {
			details.Column = col
		}
	}
}

// SuggestCorrections provides Rust-specific correction suggestions
func (r *RustErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)
	msg := errDetails.Message

	// Import/resolution errors
	if strings.Contains(msg, "E0425") || strings.Contains(msg, "E0412") ||
		strings.Contains(msg, "E0433") || strings.Contains(msg, "cannot find") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Add missing use statement or declaration",
		})
	}

	// Dependency errors
	if strings.Contains(msg, "E0432") || strings.Contains(msg, "E0463") ||
		strings.Contains(msg, "unresolved import") || strings.Contains(msg, "can't find crate") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Fix import path or add dependency to Cargo.toml",
		})
	}

	// Type errors
	if strings.Contains(msg, "E0308") || strings.Contains(msg, "mismatched types") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type mismatch",
		})
	}

	// Borrow checker errors
	if strings.Contains(msg, "E0382") || strings.Contains(msg, "E0502") ||
		strings.Contains(msg, "E0507") || strings.Contains(msg, "borrow") || strings.Contains(msg, "moved") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix ownership/borrowing issue - consider clone() or restructuring",
		})
	}

	// Lifetime errors
	if strings.Contains(msg, "E0106") || strings.Contains(msg, "E0621") || strings.Contains(msg, "lifetime") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Add or fix lifetime annotations",
		})
	}

	// Syntax errors
	if strings.Contains(msg, "expected") || strings.Contains(msg, "unexpected token") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	}

	return corrections, nil
}

// ClassifyErrorType determines the error type from Rust output
func (r *RustErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	// Check patterns first
	for _, pattern := range rustErrorPatterns {
		if pattern.pattern.MatchString(errorOutput) {
			return pattern.errType
		}
	}

	// Fallback keyword-based classification
	errorLower := strings.ToLower(errorOutput)

	if strings.Contains(errorLower, "cannot find") || strings.Contains(errorLower, "unresolved") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorLower, "crate") || strings.Contains(errorLower, "cargo") {
		return domain.ErrorTypeDependency
	}
	if strings.Contains(errorLower, "type") || strings.Contains(errorLower, "trait") ||
		strings.Contains(errorLower, "lifetime") {
		return domain.ErrorTypeTypeCheck
	}
	if strings.Contains(errorLower, "borrow") || strings.Contains(errorLower, "move") {
		return domain.ErrorTypeCompilation
	}
	if strings.Contains(errorLower, "expected") || strings.Contains(errorLower, "unexpected") {
		return domain.ErrorTypeSyntax
	}

	return domain.ErrorTypeCompilation
}

// GetLanguage returns the supported language
func (r *RustErrorAnalyzer) GetLanguage() string {
	return langRust
}
