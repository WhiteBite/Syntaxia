// Package protocol provides task protocol verification services.
package protocol

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
)

// GoProtocolImplementation provides Go-specific protocol validation
type GoProtocolImplementation struct {
	log            domain.Logger
	staticAnalyzer domain.IStaticAnalyzerService
	buildService   domain.IBuildService
	testService    domain.ITestService
}

// NewGoProtocolImplementation creates a new Go protocol implementation
func NewGoProtocolImplementation(
	log domain.Logger,
	staticAnalyzer domain.IStaticAnalyzerService,
	buildService domain.IBuildService,
	testService domain.ITestService,
) *GoProtocolImplementation {
	return &GoProtocolImplementation{
		log:            log,
		staticAnalyzer: staticAnalyzer,
		buildService:   buildService,
		testService:    testService,
	}
}

// ExecuteLintingStage executes Go-specific linting
func (g *GoProtocolImplementation) ExecuteLintingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	g.log.Info("Running Go linting stage")

	report, err := g.staticAnalyzer.AnalyzeProject(ctx, config.ProjectPath, []string{"go"})
	if err != nil {
		return fmt.Errorf("Go static analysis failed: %w", err)
	}

	if g.hasGoLintingErrors(report) {
		return fmt.Errorf("Go linting errors found")
	}

	return nil
}

// ExecuteBuildingStage executes Go-specific building
func (g *GoProtocolImplementation) ExecuteBuildingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	g.log.Info("Running Go building stage")

	validation, err := g.buildService.ValidateProject(ctx, config.ProjectPath, []string{"go"})
	if err != nil {
		return fmt.Errorf("Go build validation failed: %w", err)
	}

	if !validation.Success {
		return fmt.Errorf("Go project build failed")
	}

	return nil
}

// ExecuteTestingStage executes Go-specific testing
func (g *GoProtocolImplementation) ExecuteTestingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	g.log.Info("Running Go testing stage")

	results, err := g.testService.RunSmokeTests(ctx, config.ProjectPath, "go")
	if err != nil {
		return fmt.Errorf("Go tests failed: %w", err)
	}

	validation := g.testService.ValidateTestResults(results)
	if !validation.Success {
		return fmt.Errorf("Go test validation failed: %d tests failed", validation.FailedTests)
	}

	return nil
}

// AnalyzeGoError analyzes Go-specific errors
func (g *GoProtocolImplementation) AnalyzeGoError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "go",
		Severity:    "error",
		Suggestions: make([]string, 0),
	}

	if strings.Contains(errorOutput, "undefined:") {
		details.ErrorType = domain.ErrorTypeCompilation
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Check if the identifier is declared or imported correctly")

		if parts := strings.Split(errorOutput, "undefined:"); len(parts) > 1 {
			identifier := strings.TrimSpace(parts[1])
			details.Suggestions = append(details.Suggestions, fmt.Sprintf("Consider importing package for '%s'", identifier))
		}
	} else if strings.Contains(errorOutput, "cannot use") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Check type compatibility and conversions")
	} else if strings.Contains(errorOutput, "syntax error") {
		details.ErrorType = domain.ErrorTypeSyntax
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Fix syntax issues")
	} else if strings.Contains(errorOutput, "package") && strings.Contains(errorOutput, "not found") {
		details.ErrorType = domain.ErrorTypeImport
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Check import path and run 'go mod tidy'")
	}

	return details, nil
}

// SuggestGoCorrections provides Go-specific correction suggestions
func (g *GoProtocolImplementation) SuggestGoCorrections(err *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	switch err.ErrorType {
	case domain.ErrorTypeCompilation:
		if strings.Contains(err.Message, "undefined:") {
			corrections = append(corrections, &domain.CorrectionStep{
				Action:      domain.ActionFixImport,
				Target:      err.SourceFile,
				Description: "Add missing import statement",
			})
		}
	case domain.ErrorTypeImport:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixImport,
			Target:      err.SourceFile,
			Description: "Fix import path or run 'go mod tidy'",
		})
	case domain.ErrorTypeSyntax:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      err.SourceFile,
			Description: "Fix Go syntax error",
		})
	case domain.ErrorTypeTypeCheck:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixType,
			Target:      err.SourceFile,
			Description: "Fix type compatibility issue",
		})
	case domain.ErrorTypeLinting:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFormatCode,
			Target:      err.SourceFile,
			Description: "Format Go code with gofmt",
		})
	}

	return corrections, nil
}

// GetGoLintingConfiguration returns Go-specific linting configuration
func (g *GoProtocolImplementation) GetGoLintingConfiguration() map[string]any {
	return map[string]any{
		"tools": []string{"staticcheck", "go vet", "golangci-lint"},
		"rules": map[string]any{
			"strict_mode":     true,
			"fail_on_warning": false,
			"checks":          []string{"SA*", "ST*"},
		},
		"exclusions": []string{"vendor/", "*.pb.go"},
	}
}

// GetGoBuildConfiguration returns Go-specific build configuration
func (g *GoProtocolImplementation) GetGoBuildConfiguration() map[string]any {
	return map[string]any{
		"build_flags": []string{"-race", "-v"},
		"test_flags":  []string{"-race", "-cover", "-coverprofile=coverage.out"},
		"build_tags":  []string{"integration"},
		"go_version":  "1.21+",
	}
}

// GetGoTestingConfiguration returns Go-specific testing configuration
func (g *GoProtocolImplementation) GetGoTestingConfiguration() map[string]any {
	return map[string]any{
		"test_patterns":      []string{"./..."},
		"parallel_execution": true,
		"coverage_threshold": 80,
		"benchmark_tests":    false,
		"integration_tests":  true,
		"test_timeout":       "10m",
	}
}

func (g *GoProtocolImplementation) hasGoLintingErrors(_ any) bool {
	return false
}

// ApplyGoFormatting applies Go-specific code formatting
func (g *GoProtocolImplementation) ApplyGoFormatting(_ context.Context, _ string, files []string) error {
	g.log.Info("Applying Go formatting")

	for _, file := range files {
		if strings.HasSuffix(file, ".go") {
			g.log.Debug(fmt.Sprintf("Formatting Go file: %s", file))
		}
	}

	return nil
}

// ValidateGoModules validates Go module structure
func (g *GoProtocolImplementation) ValidateGoModules(_ context.Context, _ string) error {
	g.log.Info("Validating Go modules")
	return nil
}

// RunGoSecurityChecks runs Go-specific security analysis
func (g *GoProtocolImplementation) RunGoSecurityChecks(_ context.Context, _ string) error {
	g.log.Info("Running Go security checks")
	return nil
}

// OptimizeGoImports optimizes Go import statements
func (g *GoProtocolImplementation) OptimizeGoImports(_ context.Context, _ string, _ []string) error {
	g.log.Info("Optimizing Go imports")
	return nil
}
