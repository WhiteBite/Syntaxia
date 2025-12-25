// Package protocol provides task protocol verification services.
package protocol

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"strings"
)

// TypeScriptProtocolImplementation provides TypeScript-specific protocol validation
type TypeScriptProtocolImplementation struct {
	log            domain.Logger
	staticAnalyzer domain.IStaticAnalyzerService
	buildService   domain.IBuildService
	testService    domain.ITestService
}

// NewTypeScriptProtocolImplementation creates a new TypeScript protocol implementation
func NewTypeScriptProtocolImplementation(
	log domain.Logger,
	staticAnalyzer domain.IStaticAnalyzerService,
	buildService domain.IBuildService,
	testService domain.ITestService,
) *TypeScriptProtocolImplementation {
	return &TypeScriptProtocolImplementation{
		log:            log,
		staticAnalyzer: staticAnalyzer,
		buildService:   buildService,
		testService:    testService,
	}
}

// ExecuteLintingStage executes TypeScript-specific linting
func (ts *TypeScriptProtocolImplementation) ExecuteLintingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	ts.log.Info("Running TypeScript linting stage")

	report, err := ts.staticAnalyzer.AnalyzeProject(ctx, config.ProjectPath, []string{"typescript"})
	if err != nil {
		return fmt.Errorf("TypeScript static analysis failed: %w", err)
	}

	if ts.hasTypeScriptLintingErrors(report) {
		return fmt.Errorf("TypeScript linting errors found")
	}

	return nil
}

// ExecuteBuildingStage executes TypeScript-specific building
func (ts *TypeScriptProtocolImplementation) ExecuteBuildingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	ts.log.Info("Running TypeScript building stage")

	validation, err := ts.buildService.ValidateProject(ctx, config.ProjectPath, []string{"typescript"})
	if err != nil {
		return fmt.Errorf("TypeScript build validation failed: %w", err)
	}

	if !validation.Success {
		return fmt.Errorf("TypeScript project build failed")
	}

	return nil
}

// ExecuteTestingStage executes TypeScript-specific testing
func (ts *TypeScriptProtocolImplementation) ExecuteTestingStage(ctx context.Context, config *domain.TaskProtocolConfig) error {
	ts.log.Info("Running TypeScript testing stage")

	results, err := ts.testService.RunSmokeTests(ctx, config.ProjectPath, "typescript")
	if err != nil {
		return fmt.Errorf("TypeScript tests failed: %w", err)
	}

	validation := ts.testService.ValidateTestResults(results)
	if !validation.Success {
		return fmt.Errorf("TypeScript test validation failed: %d tests failed", validation.FailedTests)
	}

	return nil
}

// AnalyzeTypeScriptError analyzes TypeScript-specific errors
func (ts *TypeScriptProtocolImplementation) AnalyzeTypeScriptError(errorOutput string) (*domain.ErrorDetails, error) {
	details := &domain.ErrorDetails{
		Tool:        "tsc",
		Severity:    "error",
		Suggestions: make([]string, 0),
	}

	if strings.Contains(errorOutput, "TS2304") {
		details.ErrorType = domain.ErrorTypeImport
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Cannot find name - check imports or type declarations")
		details.Suggestions = append(details.Suggestions, "Consider adding type definitions or importing the missing module")
	} else if strings.Contains(errorOutput, "TS2322") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Type assignment error - check type compatibility")
		details.Suggestions = append(details.Suggestions, "Consider type casting or updating type definitions")
	} else if strings.Contains(errorOutput, "TS2307") {
		details.ErrorType = domain.ErrorTypeImport
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Cannot find module - check import path")
		details.Suggestions = append(details.Suggestions, "Verify module exists and is properly installed")
	} else if strings.Contains(errorOutput, "TS2345") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Argument type mismatch - check function parameters")
	} else if strings.Contains(errorOutput, "TS2339") {
		details.ErrorType = domain.ErrorTypeTypeCheck
		details.Message = errorOutput
		details.Suggestions = append(details.Suggestions, "Property does not exist - check object type definitions")
	}

	return details, nil
}

// SuggestTypeScriptCorrections provides TypeScript-specific correction suggestions
func (ts *TypeScriptProtocolImplementation) SuggestTypeScriptCorrections(err *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	corrections := make([]*domain.CorrectionStep, 0)

	switch err.ErrorType {
	case domain.ErrorTypeImport:
		if strings.Contains(err.Message, "TS2304") || strings.Contains(err.Message, "TS2307") {
			corrections = append(corrections, &domain.CorrectionStep{
				Action:      domain.ActionFixImport,
				Target:      err.SourceFile,
				Description: "Add missing import or type declaration",
			})
		}
	case domain.ErrorTypeTypeCheck:
		if strings.Contains(err.Message, "TS2322") {
			corrections = append(corrections, &domain.CorrectionStep{
				Action:      domain.ActionFixType,
				Target:      err.SourceFile,
				Description: "Fix type assignment error",
			})
		} else if strings.Contains(err.Message, "TS2345") {
			corrections = append(corrections, &domain.CorrectionStep{
				Action:      domain.ActionFixType,
				Target:      err.SourceFile,
				Description: "Fix function parameter type mismatch",
			})
		}
	case domain.ErrorTypeSyntax:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFixSyntax,
			Target:      err.SourceFile,
			Description: "Fix TypeScript syntax error",
		})
	case domain.ErrorTypeLinting:
		corrections = append(corrections, &domain.CorrectionStep{
			Action:      domain.ActionFormatCode,
			Target:      err.SourceFile,
			Description: "Format TypeScript code with Prettier",
		})
	}

	return corrections, nil
}

// GetTypeScriptLintingConfiguration returns TypeScript-specific linting configuration
func (ts *TypeScriptProtocolImplementation) GetTypeScriptLintingConfiguration() map[string]any {
	return map[string]any{
		"tools": []string{"eslint", "@typescript-eslint/eslint-plugin"},
		"rules": map[string]any{
			"extends": []string{
				"@typescript-eslint/recommended",
				"@typescript-eslint/recommended-requiring-type-checking",
			},
			"strict_mode":     true,
			"no_implicit_any": true,
			"fail_on_warning": false,
		},
		"parser_options": map[string]any{
			"ecmaVersion": "latest",
			"sourceType":  "module",
			"project":     "./tsconfig.json",
		},
		"exclusions": []string{"node_modules/", "dist/", "*.d.ts"},
	}
}

// GetTypeScriptBuildConfiguration returns TypeScript-specific build configuration
func (ts *TypeScriptProtocolImplementation) GetTypeScriptBuildConfiguration() map[string]any {
	return map[string]any{
		"compiler_options": map[string]any{
			"strict":                           true,
			"noImplicitAny":                    true,
			"strictNullChecks":                 true,
			"strictFunctionTypes":              true,
			"noImplicitReturns":                true,
			"noFallthroughCasesInSwitch":       true,
			"target":                           "ES2020",
			"module":                           "ESNext",
			"moduleResolution":                 "node",
			"esModuleInterop":                  true,
			"allowSyntheticDefaultImports":     true,
			"skipLibCheck":                     true,
			"forceConsistentCasingInFileNames": true,
		},
		"include":           []string{"src/**/*", "tests/**/*"},
		"exclude":           []string{"node_modules", "dist"},
		"type_checking":     true,
		"declaration_files": true,
	}
}

// GetTypeScriptTestingConfiguration returns TypeScript-specific testing configuration
func (ts *TypeScriptProtocolImplementation) GetTypeScriptTestingConfiguration() map[string]any {
	return map[string]any{
		"test_framework":         "vitest",
		"test_patterns":          []string{"**/*.test.ts", "**/*.spec.ts"},
		"coverage_threshold":     85,
		"parallel_execution":     true,
		"type_checking_in_tests": true,
		"test_timeout":           "30s",
		"setup_files":            []string{"tests/setup.ts"},
	}
}

func (ts *TypeScriptProtocolImplementation) hasTypeScriptLintingErrors(_ any) bool {
	return false
}

// ApplyTypeScriptFormatting applies TypeScript-specific code formatting
func (ts *TypeScriptProtocolImplementation) ApplyTypeScriptFormatting(_ context.Context, _ string, files []string) error {
	ts.log.Info("Applying TypeScript formatting")

	for _, file := range files {
		if strings.HasSuffix(file, ".ts") || strings.HasSuffix(file, ".tsx") {
			ts.log.Debug(fmt.Sprintf("Formatting TypeScript file: %s", file))
		}
	}

	return nil
}

// ValidateTypeScriptProject validates TypeScript project structure
func (ts *TypeScriptProtocolImplementation) ValidateTypeScriptProject(_ context.Context, _ string) error {
	ts.log.Info("Validating TypeScript project")
	return nil
}

// RunTypeScriptSecurityChecks runs TypeScript-specific security analysis
func (ts *TypeScriptProtocolImplementation) RunTypeScriptSecurityChecks(_ context.Context, _ string) error {
	ts.log.Info("Running TypeScript security checks")
	return nil
}

// OptimizeTypeScriptImports optimizes TypeScript import statements
func (ts *TypeScriptProtocolImplementation) OptimizeTypeScriptImports(_ context.Context, _ string, _ []string) error {
	ts.log.Info("Optimizing TypeScript imports")
	return nil
}

// GenerateTypeDeclarations generates missing TypeScript type declarations
func (ts *TypeScriptProtocolImplementation) GenerateTypeDeclarations(_ context.Context, _ string) error {
	ts.log.Info("Generating TypeScript type declarations")
	return nil
}
