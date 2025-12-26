package repair

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
)

// ImportCorrectionRule handles import-related corrections
type ImportCorrectionRule struct {
	log         domain.Logger
	toolChecker *ToolChecker
}

// NewImportCorrectionRule creates a new ImportCorrectionRule
func NewImportCorrectionRule(log domain.Logger, toolChecker *ToolChecker) domain.CorrectionRule {
	return &ImportCorrectionRule{log: log, toolChecker: toolChecker}
}

func (r *ImportCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeImport
}

func (r *ImportCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	if errDetails.SourceFile == "" {
		return &domain.CorrectionResult{
			Success: false,
			Message: "No source file specified for import correction",
		}, nil
	}

	ext := filepath.Ext(errDetails.SourceFile)
	filePath := filepath.Join(projectPath, errDetails.SourceFile)

	switch ext {
	case extGo:
		return r.applyGoImportFix(filePath, projectPath, errDetails.SourceFile)
	case extTS, extTSX, extJS, extJSX:
		return r.applyJSImportFix(filePath, projectPath, errDetails.SourceFile)
	}

	return &domain.CorrectionResult{
		Success: false,
		Message: fmt.Sprintf("Import fix not supported for: %s", ext),
	}, nil
}

func (r *ImportCorrectionRule) applyGoImportFix(filePath, projectPath, sourceFile string) (*domain.CorrectionResult, error) {
	if !r.toolChecker.IsAvailable("goimports") {
		return &domain.CorrectionResult{
			Success: false,
			Message: "goimports not available",
		}, nil
	}
	result, err := runCommand(context.Background(), projectPath, "goimports", "-w", filePath)
	if err != nil || !result.Success {
		return &domain.CorrectionResult{
			Success: false,
			Message: fmt.Sprintf("goimports failed: %v", err),
		}, nil
	}
	return &domain.CorrectionResult{
		Success:      true,
		Message:      "Import fixed with goimports",
		FilesChanged: []string{sourceFile},
	}, nil
}

func (r *ImportCorrectionRule) applyJSImportFix(filePath, projectPath, sourceFile string) (*domain.CorrectionResult, error) {
	if r.toolChecker.IsAvailable("eslint") || r.toolChecker.IsAvailable("npx") {
		var toolName string
		var args []string
		if r.toolChecker.IsAvailable("eslint") {
			toolName = "eslint"
			args = []string{"--fix", filePath}
		} else {
			toolName = "npx"
			args = []string{"eslint", "--fix", filePath}
		}
		_, _ = runCommand(context.Background(), projectPath, toolName, args...)
		return &domain.CorrectionResult{
			Success:      true,
			Message:      "Attempted import fix with eslint",
			FilesChanged: []string{sourceFile},
		}, nil
	}
	return &domain.CorrectionResult{
		Success: false,
		Message: "No JS/TS import fixer available",
	}, nil
}

func (r *ImportCorrectionRule) GetPriority() int {
	return 100
}

func (r *ImportCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeImport}
}


// SyntaxCorrectionRule handles syntax-related corrections
type SyntaxCorrectionRule struct {
	log domain.Logger
}

// NewSyntaxCorrectionRule creates a new SyntaxCorrectionRule
func NewSyntaxCorrectionRule(log domain.Logger) domain.CorrectionRule {
	return &SyntaxCorrectionRule{log: log}
}

func (r *SyntaxCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeSyntax
}

func (r *SyntaxCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: false,
		Message: "Syntax errors require manual correction",
	}, nil
}

func (r *SyntaxCorrectionRule) GetPriority() int {
	return 90
}

func (r *SyntaxCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeSyntax}
}

// TypeCorrectionRule handles type-related corrections
type TypeCorrectionRule struct {
	log domain.Logger
}

// NewTypeCorrectionRule creates a new TypeCorrectionRule
func NewTypeCorrectionRule(log domain.Logger) domain.CorrectionRule {
	return &TypeCorrectionRule{log: log}
}

func (r *TypeCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeTypeCheck
}

func (r *TypeCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: false,
		Message: "Type errors require manual correction or AI assistance",
	}, nil
}

func (r *TypeCorrectionRule) GetPriority() int {
	return 70
}

func (r *TypeCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeTypeCheck}
}


// LintingCorrectionRule handles linting-related corrections
type LintingCorrectionRule struct {
	log         domain.Logger
	toolChecker *ToolChecker
}

// NewLintingCorrectionRule creates a new LintingCorrectionRule
func NewLintingCorrectionRule(log domain.Logger, toolChecker *ToolChecker) domain.CorrectionRule {
	return &LintingCorrectionRule{log: log, toolChecker: toolChecker}
}

func (r *LintingCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeLinting
}

func (r *LintingCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	ext := filepath.Ext(errDetails.SourceFile)
	filePath := filepath.Join(projectPath, errDetails.SourceFile)

	switch ext {
	case extGo:
		return r.applyGoLintFix(filePath, projectPath, errDetails.SourceFile)
	case extTS, extTSX, extJS, extJSX:
		return r.applyJSLintFix(filePath, projectPath, errDetails.SourceFile)
	}

	return &domain.CorrectionResult{
		Success: false,
		Message: "No linter with --fix available",
	}, nil
}

func (r *LintingCorrectionRule) applyGoLintFix(filePath, projectPath, sourceFile string) (*domain.CorrectionResult, error) {
	if r.toolChecker.IsAvailable("golangci-lint") {
		result, _ := runCommandWithTimeout(context.Background(), projectPath, linterTimeout,
			"golangci-lint", "run", "--fix", "./...")
		return &domain.CorrectionResult{
			Success:      result != nil && result.Success,
			Message:      "Ran golangci-lint --fix",
			FilesChanged: []string{"*.go"},
		}, nil
	}
	if r.toolChecker.IsAvailable("gofmt") {
		_, _ = runCommand(context.Background(), projectPath, "gofmt", "-w", filePath)
		return &domain.CorrectionResult{
			Success:      true,
			Message:      "Formatted with gofmt",
			FilesChanged: []string{sourceFile},
		}, nil
	}
	return &domain.CorrectionResult{
		Success: false,
		Message: "No Go linter available",
	}, nil
}

func (r *LintingCorrectionRule) applyJSLintFix(filePath, projectPath, sourceFile string) (*domain.CorrectionResult, error) {
	if r.toolChecker.IsAvailable("eslint") {
		_, _ = runCommand(context.Background(), projectPath, "eslint", "--fix", filePath)
		return &domain.CorrectionResult{
			Success:      true,
			Message:      "Ran eslint --fix",
			FilesChanged: []string{sourceFile},
		}, nil
	}
	if r.toolChecker.IsAvailable("npx") {
		_, _ = runCommand(context.Background(), projectPath, "npx", "eslint", "--fix", filePath)
		return &domain.CorrectionResult{
			Success:      true,
			Message:      "Ran npx eslint --fix",
			FilesChanged: []string{sourceFile},
		}, nil
	}
	return &domain.CorrectionResult{
		Success: false,
		Message: "No JS/TS linter available",
	}, nil
}

func (r *LintingCorrectionRule) GetPriority() int {
	return 80
}

func (r *LintingCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeLinting}
}


// CompilationCorrectionRule handles compilation-related corrections
type CompilationCorrectionRule struct {
	log domain.Logger
}

// NewCompilationCorrectionRule creates a new CompilationCorrectionRule
func NewCompilationCorrectionRule(log domain.Logger) domain.CorrectionRule {
	return &CompilationCorrectionRule{log: log}
}

func (r *CompilationCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeCompilation
}

func (r *CompilationCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: false,
		Message: "Compilation errors require manual correction or AI assistance",
	}, nil
}

func (r *CompilationCorrectionRule) GetPriority() int {
	return 95
}

func (r *CompilationCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeCompilation}
}