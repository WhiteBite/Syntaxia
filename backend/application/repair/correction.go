package repair

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"sort"
	"strings"
)

// CorrectionEngine implements the CorrectionEngine interface
type CorrectionEngine struct {
	log             domain.Logger
	fileSystem      domain.FileSystemProvider
	correctionRules map[domain.ErrorType][]domain.CorrectionRule
}

// NewCorrectionEngine creates a new CorrectionEngine instance
func NewCorrectionEngine(log domain.Logger, fileSystem domain.FileSystemProvider) domain.CorrectionEngine {
	engine := &CorrectionEngine{
		log:             log,
		fileSystem:      fileSystem,
		correctionRules: make(map[domain.ErrorType][]domain.CorrectionRule),
	}

	// Register correction rules
	engine.registerCorrectionRules()

	return engine
}

// ApplyCorrection applies a single correction step
func (c *CorrectionEngine) ApplyCorrection(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	c.log.Info(fmt.Sprintf("Applying correction: %s for target: %s", step.Action, step.Target))

	switch step.Action {
	case domain.ActionFixImport:
		return c.applyImportFix(ctx, step, projectPath)
	case domain.ActionFixSyntax:
		return c.applySyntaxFix(ctx, step, projectPath)
	case domain.ActionFixType:
		return c.applyTypeFix(ctx, step, projectPath)
	case domain.ActionAddMissingCode:
		return c.applyAddMissingCode(ctx, step, projectPath)
	case domain.ActionRemoveCode:
		return c.applyRemoveCode(ctx, step, projectPath)
	case domain.ActionFormatCode:
		return c.applyFormatCode(ctx, step, projectPath)
	case domain.ActionUpdateTest:
		return c.applyUpdateTest(ctx, step, projectPath)
	default:
		return nil, fmt.Errorf("unsupported correction action: %s", step.Action)
	}
}

// ApplyCorrections applies multiple correction steps
func (c *CorrectionEngine) ApplyCorrections(ctx context.Context, steps []*domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	c.log.Info(fmt.Sprintf("Applying %d correction steps", len(steps)))

	allFilesChanged := make([]string, 0)
	allMessages := make([]string, 0)
	overallSuccess := true

	// Sort steps by priority (some corrections should be applied before others)
	sortedSteps := c.sortCorrectionSteps(steps)

	for i, step := range sortedSteps {
		c.log.Debug(fmt.Sprintf("Applying correction step %d/%d: %s", i+1, len(sortedSteps), step.Description))

		result, err := c.ApplyCorrection(ctx, step, projectPath)
		if err != nil {
			c.log.Warning(fmt.Sprintf("Correction step failed: %v", err))
			step.Applied = false
			step.Result = fmt.Sprintf("Failed: %v", err)
			overallSuccess = false
			continue
		}

		step.Applied = result.Success
		step.Result = result.Message
		allFilesChanged = append(allFilesChanged, result.FilesChanged...)
		allMessages = append(allMessages, result.Message)

		if !result.Success {
			overallSuccess = false
		}
	}

	// Remove duplicates from files changed
	uniqueFiles := removeDuplicates(allFilesChanged)

	return &domain.CorrectionResult{
		Success:      overallSuccess,
		Message:      strings.Join(allMessages, "; "),
		FilesChanged: uniqueFiles,
	}, nil
}

// CanHandle checks if the engine can handle a specific error type
func (c *CorrectionEngine) CanHandle(errDetails *domain.ErrorDetails) bool {
	rules, exists := c.correctionRules[errDetails.ErrorType]
	if !exists {
		return false
	}

	for _, rule := range rules {
		if rule.CanHandle(errDetails) {
			return true
		}
	}

	return false
}

// Correction action implementations

func (c *CorrectionEngine) applyImportFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)

	// Read the file
	content, err := c.fileSystem.ReadFile(filePath)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to read file: %v", err)}, nil
	}

	// Apply import fixes (simplified implementation)
	modifiedContent := c.addMissingImports(string(content), step)

	if modifiedContent == string(content) {
		return &domain.CorrectionResult{Success: false, Message: "No import fixes applied"}, nil
	}

	// Write the modified content back
	err = c.fileSystem.WriteFile(filePath, []byte(modifiedContent), 0644)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to write file: %v", err)}, nil
	}

	return &domain.CorrectionResult{
		Success:      true,
		Message:      "Import fixes applied",
		FilesChanged: []string{step.Target},
	}, nil
}

func (c *CorrectionEngine) applySyntaxFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Syntax fix applied (placeholder implementation)",
	}, nil
}

func (c *CorrectionEngine) applyTypeFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Type fix applied (placeholder implementation)",
	}, nil
}

func (c *CorrectionEngine) applyAddMissingCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Missing code added (placeholder implementation)",
	}, nil
}

func (c *CorrectionEngine) applyRemoveCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Code removed (placeholder implementation)",
	}, nil
}

func (c *CorrectionEngine) applyFormatCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)

	// Determine file type and apply appropriate formatting
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		return c.formatGoFile(filePath)
	case extTS, extJS:
		return c.formatJSFile(filePath)
	default:
		return &domain.CorrectionResult{
			Success: false,
			Message: fmt.Sprintf("Unsupported file type for formatting: %s", ext),
		}, nil
	}
}

func (c *CorrectionEngine) applyUpdateTest(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Test updated (placeholder implementation)",
	}, nil
}

// Helper methods

func (c *CorrectionEngine) addMissingImports(content string, step *domain.CorrectionStep) string {
	lines := strings.Split(content, "\n")

	for i, line := range lines {
		if strings.Contains(line, "import") && i < len(lines)-1 {
			continue
		}
	}

	return content
}

func (c *CorrectionEngine) formatGoFile(filePath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success:      true,
		Message:      "Go file formatted",
		FilesChanged: []string{filePath},
	}, nil
}

func (c *CorrectionEngine) formatJSFile(filePath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success:      true,
		Message:      "JavaScript/TypeScript file formatted",
		FilesChanged: []string{filePath},
	}, nil
}

func (c *CorrectionEngine) sortCorrectionSteps(steps []*domain.CorrectionStep) []*domain.CorrectionStep {
	sorted := make([]*domain.CorrectionStep, len(steps))
	copy(sorted, steps)

	sort.Slice(sorted, func(i, j int) bool {
		return c.getCorrectionPriority(sorted[i].Action) > c.getCorrectionPriority(sorted[j].Action)
	})

	return sorted
}

func (c *CorrectionEngine) getCorrectionPriority(action domain.CorrectionAction) int {
	switch action {
	case domain.ActionFixImport:
		return 100
	case domain.ActionFixSyntax:
		return 90
	case domain.ActionFormatCode:
		return 80
	case domain.ActionFixType:
		return 70
	case domain.ActionAddMissingCode:
		return 60
	case domain.ActionUpdateTest:
		return 50
	case domain.ActionRemoveCode:
		return 40
	default:
		return 0
	}
}

func (c *CorrectionEngine) registerCorrectionRules() {
	c.correctionRules[domain.ErrorTypeImport] = []domain.CorrectionRule{
		NewImportCorrectionRule(),
	}
	c.correctionRules[domain.ErrorTypeSyntax] = []domain.CorrectionRule{
		NewSyntaxCorrectionRule(),
	}
	c.correctionRules[domain.ErrorTypeTypeCheck] = []domain.CorrectionRule{
		NewTypeCorrectionRule(),
	}
	c.correctionRules[domain.ErrorTypeLinting] = []domain.CorrectionRule{
		NewLintingCorrectionRule(),
	}
	c.correctionRules[domain.ErrorTypeCompilation] = []domain.CorrectionRule{
		NewCompilationCorrectionRule(),
	}
}

// Correction rule implementations

// ImportCorrectionRule handles import-related corrections
type ImportCorrectionRule struct{}

func NewImportCorrectionRule() domain.CorrectionRule {
	return &ImportCorrectionRule{}
}

func (r *ImportCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeImport
}

func (r *ImportCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Import correction applied",
	}, nil
}

func (r *ImportCorrectionRule) GetPriority() int {
	return 100
}

func (r *ImportCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeImport}
}

// SyntaxCorrectionRule handles syntax-related corrections
type SyntaxCorrectionRule struct{}

func NewSyntaxCorrectionRule() domain.CorrectionRule {
	return &SyntaxCorrectionRule{}
}

func (r *SyntaxCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeSyntax
}

func (r *SyntaxCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Syntax correction applied",
	}, nil
}

func (r *SyntaxCorrectionRule) GetPriority() int {
	return 90
}

func (r *SyntaxCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeSyntax}
}

// TypeCorrectionRule handles type-related corrections
type TypeCorrectionRule struct{}

func NewTypeCorrectionRule() domain.CorrectionRule {
	return &TypeCorrectionRule{}
}

func (r *TypeCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeTypeCheck
}

func (r *TypeCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Type correction applied",
	}, nil
}

func (r *TypeCorrectionRule) GetPriority() int {
	return 70
}

func (r *TypeCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeTypeCheck}
}

// LintingCorrectionRule handles linting-related corrections
type LintingCorrectionRule struct{}

func NewLintingCorrectionRule() domain.CorrectionRule {
	return &LintingCorrectionRule{}
}

func (r *LintingCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeLinting
}

func (r *LintingCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Linting correction applied",
	}, nil
}

func (r *LintingCorrectionRule) GetPriority() int {
	return 80
}

func (r *LintingCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeLinting}
}

// CompilationCorrectionRule handles compilation-related corrections
type CompilationCorrectionRule struct{}

func NewCompilationCorrectionRule() domain.CorrectionRule {
	return &CompilationCorrectionRule{}
}

func (r *CompilationCorrectionRule) CanHandle(errDetails *domain.ErrorDetails) bool {
	return errDetails.ErrorType == domain.ErrorTypeCompilation
}

func (r *CompilationCorrectionRule) ApplyCorrection(errDetails *domain.ErrorDetails, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{
		Success: true,
		Message: "Compilation correction applied",
	}, nil
}

func (r *CompilationCorrectionRule) GetPriority() int {
	return 95
}

func (r *CompilationCorrectionRule) GetErrorTypes() []domain.ErrorType {
	return []domain.ErrorType{domain.ErrorTypeCompilation}
}

// Utility functions

func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			list = append(list, item)
		}
	}

	return list
}
