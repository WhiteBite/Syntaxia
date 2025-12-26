package repair

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"syntaxia/domain"
)

// CorrectionEngine implements the CorrectionEngine interface
type CorrectionEngine struct {
	log             domain.Logger
	fileSystem      domain.FileSystemProvider
	correctionRules map[domain.ErrorType][]domain.CorrectionRule
	toolChecker     *ToolChecker
}

// NewCorrectionEngine creates a new CorrectionEngine instance
func NewCorrectionEngine(log domain.Logger, fileSystem domain.FileSystemProvider) domain.CorrectionEngine {
	engine := &CorrectionEngine{
		log:             log,
		fileSystem:      fileSystem,
		correctionRules: make(map[domain.ErrorType][]domain.CorrectionRule),
		toolChecker:     NewToolChecker(),
	}
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

	return &domain.CorrectionResult{
		Success:      overallSuccess,
		Message:      strings.Join(allMessages, "; "),
		FilesChanged: removeDuplicates(allFilesChanged),
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


// applyImportFix fixes import issues using language-specific tools
func (c *CorrectionEngine) applyImportFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)
	ext := filepath.Ext(filePath)
	language := detectLanguageFromExtension(ext)

	switch language {
	case langGo:
		return c.fixGoImports(ctx, filePath, projectPath)
	case langTypeScript, langJavaScript:
		return c.fixTSImports(ctx, filePath, projectPath, step)
	default:
		return &domain.CorrectionResult{
			Success: false,
			Message: fmt.Sprintf("Import fix not supported for extension: %s", ext),
		}, nil
	}
}

// fixGoImports runs goimports on a Go file
func (c *CorrectionEngine) fixGoImports(ctx context.Context, filePath, projectPath string) (*domain.CorrectionResult, error) {
	if !c.toolChecker.IsAvailable("goimports") {
		c.log.Warning("goimports not available, trying gofmt")
		return c.formatGoFile(ctx, filePath, projectPath)
	}

	result, err := runCommandWithTimeout(ctx, projectPath, goimportsTimeout, "goimports", "-w", filePath)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to run goimports: %v", err)}, nil
	}
	if !result.Success {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("goimports failed: %s", result.Output)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: "Go imports fixed with goimports", FilesChanged: []string{filePath}}, nil
}

// fixTSImports attempts to fix TypeScript/JavaScript imports
func (c *CorrectionEngine) fixTSImports(ctx context.Context, filePath, projectPath string, step *domain.CorrectionStep) (*domain.CorrectionResult, error) {
	content, err := c.fileSystem.ReadFile(filePath)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to read file: %v", err)}, nil
	}

	missingIdent := extractMissingIdentifier(step.Description)
	if missingIdent == "" {
		return c.runESLintFix(ctx, filePath, projectPath)
	}

	modifiedContent, found := c.addTSImport(string(content), missingIdent)
	if !found {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Could not find export for: %s", missingIdent)}, nil
	}

	if err = c.fileSystem.WriteFile(filePath, []byte(modifiedContent), 0644); err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to write file: %v", err)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: fmt.Sprintf("Added import for: %s", missingIdent), FilesChanged: []string{filePath}}, nil
}


// addTSImport attempts to add a TypeScript import statement
func (c *CorrectionEngine) addTSImport(content, identifier string) (string, bool) {
	knownImports := map[string]string{
		"ref": "import { ref } from 'vue'", "computed": "import { computed } from 'vue'",
		"watch": "import { watch } from 'vue'", "onMounted": "import { onMounted } from 'vue'",
		"defineComponent": "import { defineComponent } from 'vue'",
		"useState": "import { useState } from 'react'", "useEffect": "import { useEffect } from 'react'",
	}

	importStmt, ok := knownImports[identifier]
	if !ok {
		return content, false
	}
	if strings.Contains(content, importStmt) {
		return content, true
	}

	lines := strings.Split(content, "\n")
	insertIdx := 0
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "import ") {
			insertIdx = i + 1
		} else if trimmed != "" && !strings.HasPrefix(trimmed, "//") && insertIdx > 0 {
			break
		}
	}

	newLines := make([]string, 0, len(lines)+1)
	newLines = append(newLines, lines[:insertIdx]...)
	newLines = append(newLines, importStmt)
	newLines = append(newLines, lines[insertIdx:]...)
	return strings.Join(newLines, "\n"), true
}

// applySyntaxFix attempts to fix syntax errors
func (c *CorrectionEngine) applySyntaxFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		return c.formatGoFile(ctx, filePath, projectPath)
	case extTS, extTSX, extJS, extJSX, extVue:
		return c.formatJSFile(ctx, filePath, projectPath)
	default:
		return &domain.CorrectionResult{Success: false, Message: "Syntax fix requires manual intervention"}, nil
	}
}

// applyTypeFix attempts to fix type errors
func (c *CorrectionEngine) applyTypeFix(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		return c.runGoLintFix(ctx, projectPath)
	case extTS, extTSX, extJS, extJSX:
		return c.runESLintFix(ctx, filePath, projectPath)
	default:
		return &domain.CorrectionResult{Success: false, Message: "Type fix not automatically supported"}, nil
	}
}


// applyAddMissingCode handles adding missing code
func (c *CorrectionEngine) applyAddMissingCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{Success: false, Message: "Adding missing code requires AI assistance"}, nil
}

// applyRemoveCode handles removing code
func (c *CorrectionEngine) applyRemoveCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		return c.fixGoImports(ctx, filePath, projectPath)
	case extTS, extTSX, extJS, extJSX:
		return c.runESLintFix(ctx, filePath, projectPath)
	default:
		return &domain.CorrectionResult{Success: false, Message: "Remove code not automatically supported"}, nil
	}
}

// applyFormatCode formats code using language-specific tools
func (c *CorrectionEngine) applyFormatCode(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	filePath := filepath.Join(projectPath, step.Target)
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		return c.formatGoFile(ctx, filePath, projectPath)
	case extTS, extTSX, extJS, extJSX, extVue:
		return c.formatJSFile(ctx, filePath, projectPath)
	case extPy:
		return c.formatPythonFile(ctx, filePath, projectPath)
	default:
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Formatting not supported for: %s", ext)}, nil
	}
}

// applyUpdateTest handles test updates
func (c *CorrectionEngine) applyUpdateTest(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return &domain.CorrectionResult{Success: false, Message: "Test updates require AI assistance"}, nil
}


// formatGoFile formats a Go file using gofmt or goimports
func (c *CorrectionEngine) formatGoFile(ctx context.Context, filePath, projectPath string) (*domain.CorrectionResult, error) {
	var toolName string
	var args []string

	if c.toolChecker.IsAvailable("goimports") {
		toolName, args = "goimports", []string{"-w", filePath}
	} else if c.toolChecker.IsAvailable("gofmt") {
		toolName, args = "gofmt", []string{"-w", filePath}
	} else {
		return &domain.CorrectionResult{Success: false, Message: "No Go formatter available"}, nil
	}

	result, err := runCommand(ctx, projectPath, toolName, args...)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to run %s: %v", toolName, err)}, nil
	}
	if !result.Success {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("%s failed: %s", toolName, result.Output)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: fmt.Sprintf("Go file formatted with %s", toolName), FilesChanged: []string{filePath}}, nil
}

// formatJSFile formats a JavaScript/TypeScript file using prettier
func (c *CorrectionEngine) formatJSFile(ctx context.Context, filePath, projectPath string) (*domain.CorrectionResult, error) {
	var toolName string
	var args []string

	if c.toolChecker.IsAvailable("prettier") {
		toolName, args = "prettier", []string{"--write", filePath}
	} else if c.toolChecker.IsAvailable("npx") {
		toolName, args = "npx", []string{"prettier", "--write", filePath}
	} else {
		return &domain.CorrectionResult{Success: false, Message: "No JS/TS formatter available"}, nil
	}

	result, err := runCommand(ctx, projectPath, toolName, args...)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to run %s: %v", toolName, err)}, nil
	}
	if !result.Success {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("%s failed: %s", toolName, result.Output)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: "JS/TS file formatted with prettier", FilesChanged: []string{filePath}}, nil
}


// formatPythonFile formats a Python file using black or autopep8
func (c *CorrectionEngine) formatPythonFile(ctx context.Context, filePath, projectPath string) (*domain.CorrectionResult, error) {
	var toolName string
	var args []string

	if c.toolChecker.IsAvailable("black") {
		toolName, args = "black", []string{filePath}
	} else if c.toolChecker.IsAvailable("autopep8") {
		toolName, args = "autopep8", []string{"--in-place", filePath}
	} else {
		return &domain.CorrectionResult{Success: false, Message: "No Python formatter available"}, nil
	}

	result, err := runCommand(ctx, projectPath, toolName, args...)
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to run %s: %v", toolName, err)}, nil
	}
	if !result.Success {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("%s failed: %s", toolName, result.Output)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: fmt.Sprintf("Python file formatted with %s", toolName), FilesChanged: []string{filePath}}, nil
}

// runGoLintFix runs golangci-lint with --fix flag
func (c *CorrectionEngine) runGoLintFix(ctx context.Context, projectPath string) (*domain.CorrectionResult, error) {
	if !c.toolChecker.IsAvailable("golangci-lint") {
		return &domain.CorrectionResult{Success: false, Message: "golangci-lint not available"}, nil
	}

	_, err := runCommandWithTimeout(ctx, projectPath, linterTimeout, "golangci-lint", "run", "--fix", "./...")
	if err != nil {
		return &domain.CorrectionResult{Success: false, Message: fmt.Sprintf("Failed to run golangci-lint: %v", err)}, nil
	}
	return &domain.CorrectionResult{Success: true, Message: "Ran golangci-lint --fix", FilesChanged: []string{"*.go"}}, nil
}

// runESLintFix runs eslint with --fix flag
func (c *CorrectionEngine) runESLintFix(ctx context.Context, filePath, projectPath string) (*domain.CorrectionResult, error) {
	var toolName string
	var args []string

	if c.toolChecker.IsAvailable("eslint") {
		toolName, args = "eslint", []string{"--fix", filePath}
	} else if c.toolChecker.IsAvailable("npx") {
		toolName, args = "npx", []string{"eslint", "--fix", filePath}
	} else {
		return &domain.CorrectionResult{Success: false, Message: "eslint not available"}, nil
	}

	_, _ = runCommandWithTimeout(ctx, projectPath, linterTimeout, toolName, args...)
	return &domain.CorrectionResult{Success: true, Message: "Ran eslint --fix", FilesChanged: []string{filePath}}, nil
}


// Helper methods

func (c *CorrectionEngine) sortCorrectionSteps(steps []*domain.CorrectionStep) []*domain.CorrectionStep {
	sorted := make([]*domain.CorrectionStep, len(steps))
	copy(sorted, steps)
	sort.Slice(sorted, func(i, j int) bool {
		return c.getCorrectionPriority(sorted[i].Action) > c.getCorrectionPriority(sorted[j].Action)
	})
	return sorted
}

func (c *CorrectionEngine) getCorrectionPriority(action domain.CorrectionAction) int {
	priorities := map[domain.CorrectionAction]int{
		domain.ActionFixImport: 100, domain.ActionFixSyntax: 90, domain.ActionFormatCode: 80,
		domain.ActionFixType: 70, domain.ActionAddMissingCode: 60, domain.ActionUpdateTest: 50,
		domain.ActionRemoveCode: 40,
	}
	if p, ok := priorities[action]; ok {
		return p
	}
	return 0
}

func (c *CorrectionEngine) registerCorrectionRules() {
	c.correctionRules[domain.ErrorTypeImport] = []domain.CorrectionRule{NewImportCorrectionRule(c.log, c.toolChecker)}
	c.correctionRules[domain.ErrorTypeSyntax] = []domain.CorrectionRule{NewSyntaxCorrectionRule(c.log)}
	c.correctionRules[domain.ErrorTypeTypeCheck] = []domain.CorrectionRule{NewTypeCorrectionRule(c.log)}
	c.correctionRules[domain.ErrorTypeLinting] = []domain.CorrectionRule{NewLintingCorrectionRule(c.log, c.toolChecker)}
	c.correctionRules[domain.ErrorTypeCompilation] = []domain.CorrectionRule{NewCompilationCorrectionRule(c.log)}
}

// extractMissingIdentifier extracts the missing identifier from error description
func extractMissingIdentifier(description string) string {
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`Cannot find name '(\w+)'`),
		regexp.MustCompile(`undefined: (\w+)`),
		regexp.MustCompile(`'(\w+)' is not defined`),
		regexp.MustCompile(`missing import for: (\w+)`),
	}
	for _, pattern := range patterns {
		if matches := pattern.FindStringSubmatch(description); len(matches) > 1 {
			return matches[1]
		}
	}
	return ""
}

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