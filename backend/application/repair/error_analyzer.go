package repair

import (
	"fmt"
	"regexp"
	"strings"
	"syntaxia/domain"
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
	analyzer.languageAnalyzers[langJava] = NewJavaErrorAnalyzer()
	analyzer.languageAnalyzers[langRust] = NewRustErrorAnalyzer()
	analyzer.languageAnalyzers[langPython] = NewPythonErrorAnalyzer()
	analyzer.languageAnalyzers[langKotlin] = NewKotlinErrorAnalyzer()
	analyzer.languageAnalyzers[langCSharp] = NewCSharpErrorAnalyzer()
	analyzer.languageAnalyzers[langDart] = NewDartErrorAnalyzer()
	analyzer.languageAnalyzers[langCpp] = NewCPPErrorAnalyzer()
	analyzer.languageAnalyzers[langPHP] = NewPHPErrorAnalyzer()
	analyzer.languageAnalyzers[langRuby] = NewRubyErrorAnalyzer()
	analyzer.languageAnalyzers[langSwift] = NewSwiftErrorAnalyzer()

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

// Utility functions

func parseIntSafe(s string) (int, error) {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	return result, err
}
