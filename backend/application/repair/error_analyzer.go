package repair

import (
	"fmt"
	"regexp"
	"syntaxia/domain"
	"strings"
)

// ErrorAnalyzer implements the ErrorAnalyzer interface
type ErrorAnalyzer struct {
	log               domain.Logger
	languageAnalyzers map[string]domain.LanguageErrorAnalyzer
}

// NewErrorAnalyzer creates a new ErrorAnalyzer instance
func NewErrorAnalyzer(log domain.Logger) domain.ErrorAnalyzer {
	analyzer := &ErrorAnalyzer{
		log:               log,
		languageAnalyzers: make(map[string]domain.LanguageErrorAnalyzer),
	}

	// Register language-specific analyzers
	analyzer.languageAnalyzers[langGo] = NewGoErrorAnalyzer()
	analyzer.languageAnalyzers[langTypeScript] = NewTypeScriptErrorAnalyzer()
	analyzer.languageAnalyzers[langJavaScript] = NewJavaScriptErrorAnalyzer()

	return analyzer
}

// AnalyzeError analyzes error output and provides detailed error information
func (e *ErrorAnalyzer) AnalyzeError(errorOutput string, stage domain.ProtocolStage) (*domain.ErrorDetails, error) {
	e.log.Debug(fmt.Sprintf("Analyzing error for stage %s: %s", stage, errorOutput))

	errorDetails := &domain.ErrorDetails{
		Stage:       stage,
		Message:     errorOutput,
		ErrorType:   e.ClassifyErrorType(errorOutput),
		Severity:    "error",
		Suggestions: make([]string, 0),
	}

	// Try to extract file and line information
	e.extractLocationInfo(errorOutput, errorDetails)

	// Try language-specific analysis
	for language, analyzer := range e.languageAnalyzers {
		if langDetails, err := analyzer.AnalyzeError(errorOutput); err == nil {
			e.log.Debug(fmt.Sprintf("Language-specific analysis successful for %s", language))
			e.mergeErrorDetails(errorDetails, langDetails)
			break
		}
	}

	// Add stage-specific analysis
	e.addStageSpecificAnalysis(errorDetails, stage)

	return errorDetails, nil
}

// SuggestCorrections provides correction suggestions for a given error
func (e *ErrorAnalyzer) SuggestCorrections(errDetails *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	// Try language-specific corrections first
	for _, analyzer := range e.languageAnalyzers {
		if langCorrections, err := analyzer.SuggestCorrections(errDetails); err == nil && len(langCorrections) > 0 {
			corrections = append(corrections, langCorrections...)
		}
	}

	// Add generic corrections based on error type
	genericCorrections := e.getGenericCorrections(errDetails)
	corrections = append(corrections, genericCorrections...)

	return corrections, nil
}

// ClassifyErrorType determines the type of error from output
func (e *ErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	errorLower := strings.ToLower(errorOutput)

	// TypeScript specific errors (check first for specificity)
	if strings.Contains(errorOutput, "TS2304") || strings.Contains(errorLower, "cannot find name") {
		return domain.ErrorTypeImport
	}

	// Import errors
	if strings.Contains(errorLower, "import") || strings.Contains(errorLower, "module") {
		return domain.ErrorTypeImport
	}

	// Compilation errors
	if strings.Contains(errorLower, "compile") || strings.Contains(errorLower, "syntax error") {
		return domain.ErrorTypeCompilation
	}

	// Type checking errors
	if strings.Contains(errorLower, "type") || strings.Contains(errorLower, "cannot use") {
		return domain.ErrorTypeTypeCheck
	}

	// Linting errors
	if strings.Contains(errorLower, "lint") || strings.Contains(errorLower, "style") {
		return domain.ErrorTypeLinting
	}

	// Test errors
	if strings.Contains(errorLower, "test") || strings.Contains(errorLower, "spec") {
		return domain.ErrorTypeTesting
	}

	// Default to compilation
	return domain.ErrorTypeCompilation
}

// Helper methods

func (e *ErrorAnalyzer) extractLocationInfo(errorOutput string, details *domain.ErrorDetails) {
	patterns := []string{
		`([^:]+):(\d+):(\d+):`,
		`([^:]+):(\d+):`,
		`at ([^:]+):(\d+):(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(errorOutput)
		if len(matches) >= 3 {
			details.SourceFile = matches[1]
			if len(matches) >= 3 {
				if line, err := parseIntSafe(matches[2]); err == nil {
					details.LineNumber = line
				}
			}
			if len(matches) >= 4 {
				if col, err := parseIntSafe(matches[3]); err == nil {
					details.Column = col
				}
			}
			break
		}
	}
}

func (e *ErrorAnalyzer) mergeErrorDetails(target, source *domain.ErrorDetails) {
	if source.SourceFile != "" && target.SourceFile == "" {
		target.SourceFile = source.SourceFile
	}
	if source.LineNumber > 0 && target.LineNumber == 0 {
		target.LineNumber = source.LineNumber
	}
	if source.Column > 0 && target.Column == 0 {
		target.Column = source.Column
	}
	if source.Tool != "" {
		target.Tool = source.Tool
	}
	if len(source.Suggestions) > 0 {
		target.Suggestions = append(target.Suggestions, source.Suggestions...)
	}
}

func (e *ErrorAnalyzer) addStageSpecificAnalysis(details *domain.ErrorDetails, stage domain.ProtocolStage) {
	switch stage {
	case domain.StageLinting:
		details.Tool = "static-analyzer"
		if details.ErrorType == "" {
			details.ErrorType = domain.ErrorTypeLinting
		}
	case domain.StageBuilding:
		details.Tool = "compiler"
		if details.ErrorType == "" {
			details.ErrorType = domain.ErrorTypeCompilation
		}
	case domain.StageTesting:
		details.Tool = "test-runner"
		if details.ErrorType == "" {
			details.ErrorType = domain.ErrorTypeTesting
		}
	case domain.StageGuardrails:
		details.Tool = "guardrails"
		if details.ErrorType == "" {
			details.ErrorType = domain.ErrorTypeGuardrail
		}
	}
}

func (e *ErrorAnalyzer) getGenericCorrections(errDetails *domain.ErrorDetails) []*domain.CorrectionStep {
	corrections := make([]*domain.CorrectionStep, 0)

	switch errDetails.ErrorType {
	case domain.ErrorTypeImport:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      errDetails.SourceFile,
			Description: "Fix import statement",
		})
	case domain.ErrorTypeSyntax:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      errDetails.SourceFile,
			Description: "Fix syntax error",
		})
	case domain.ErrorTypeTypeCheck:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      errDetails.SourceFile,
			Description: "Fix type error",
		})
	case domain.ErrorTypeLinting:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFormatCode,
			Target:      errDetails.SourceFile,
			Description: "Format code to fix linting issues",
		})
	}

	return corrections
}

// Language-specific analyzers

// GoErrorAnalyzer analyzes Go-specific errors
type GoErrorAnalyzer struct{}

func NewGoErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &GoErrorAnalyzer{}
}

func (g *GoErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool: "go",
	}

	if strings.Contains(errorOutput, "undefined:") {
		details.ErrorType = domain.ErrorTypeCompilation
		details.Suggestions = append(details.Suggestions, "Check if the identifier is declared or imported")
	}

	if strings.Contains(errorOutput, "cannot use") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Suggestions = append(details.Suggestions, "Check type compatibility")
	}

	return details, nil
}

func (g *GoErrorAnalyzer) SuggestCorrections(error *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	if strings.Contains(error.Message, "undefined:") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionAddMissingCode,
			Target:      error.SourceFile,
			Description: "Add missing declaration or import",
		})
	}

	return corrections, nil
}

func (g *GoErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	if strings.Contains(errorOutput, "undefined:") || strings.Contains(errorOutput, "undeclared name:") {
		return domain.ErrorTypeCompilation
	}
	if strings.Contains(errorOutput, "cannot use") || strings.Contains(errorOutput, "type") {
		return domain.ErrorTypeTypeCheck
	}
	return domain.ErrorTypeCompilation
}

func (g *GoErrorAnalyzer) GetLanguage() string {
	return langGo
}

// TypeScriptErrorAnalyzer analyzes TypeScript-specific errors
type TypeScriptErrorAnalyzer struct{}

func NewTypeScriptErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &TypeScriptErrorAnalyzer{}
}

func (t *TypeScriptErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool: "tsc",
	}

	if strings.Contains(errorOutput, "TS2304") {
		details.ErrorType = domain.ErrorTypeImport
		details.Suggestions = append(details.Suggestions, "Cannot find name - check imports or declarations")
	}

	if strings.Contains(errorOutput, "TS2322") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Suggestions = append(details.Suggestions, "Type assignment error - check type compatibility")
	}

	return details, nil
}

func (t *TypeScriptErrorAnalyzer) SuggestCorrections(error *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	if strings.Contains(error.Message, "TS2304") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      error.SourceFile,
			Description: "Add missing import or type declaration",
		})
	}

	return corrections, nil
}

func (t *TypeScriptErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	if strings.Contains(errorOutput, "TS2304") || strings.Contains(errorOutput, "Cannot find") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "TS2322") || strings.Contains(errorOutput, "Type") {
		return domain.ErrorTypeTypeCheck
	}
	return domain.ErrorTypeCompilation
}

func (t *TypeScriptErrorAnalyzer) GetLanguage() string {
	return langTypeScript
}

// JavaScriptErrorAnalyzer analyzes JavaScript-specific errors
type JavaScriptErrorAnalyzer struct{}

func NewJavaScriptErrorAnalyzer() domain.LanguageErrorAnalyzer {
	return &JavaScriptErrorAnalyzer{}
}

func (j *JavaScriptErrorAnalyzer) AnalyzeError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool: "eslint",
	}

	if strings.Contains(errorOutput, "is not defined") {
		details.ErrorType = domain.ErrorTypeImport
		details.Suggestions = append(details.Suggestions, "Variable or function not defined - check imports")
	}

	return details, nil
}

func (j *JavaScriptErrorAnalyzer) SuggestCorrections(error *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	if strings.Contains(error.Message, "is not defined") {
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      error.SourceFile,
			Description: "Add missing import or declaration",
		})
	}

	return corrections, nil
}

func (j *JavaScriptErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	if strings.Contains(errorOutput, "is not defined") {
		return domain.ErrorTypeImport
	}
	if strings.Contains(errorOutput, "SyntaxError") {
		return domain.ErrorTypeSyntax
	}
	return domain.ErrorTypeCompilation
}

func (j *JavaScriptErrorAnalyzer) GetLanguage() string {
	return langJavaScript
}

// Utility functions

func parseIntSafe(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
