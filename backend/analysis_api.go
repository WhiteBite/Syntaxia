package main

import (
	"fmt"
	"path/filepath"
	"syntaxia/application/project"
	"syntaxia/domain"
	"syntaxia/infrastructure/git"
	"strings"
	"time"
)

// BuildSymbolGraph builds a symbol graph for a project
func (a *App) BuildSymbolGraph(projectRoot, language string) (*domain.SymbolGraph, error) {
	return a.analysisHandler.BuildSymbolGraph(a.ctx, projectRoot, language)
}

// GetSymbolSuggestions returns symbol suggestions
func (a *App) GetSymbolSuggestions(query, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return a.analysisHandler.GetSymbolSuggestions(a.ctx, query, language, graph)
}

// GetSymbolDependencies returns symbol dependencies
func (a *App) GetSymbolDependencies(symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return a.analysisHandler.GetSymbolDependencies(a.ctx, symbolID, language, graph)
}

// GetSymbolDependents returns symbols that depend on the specified symbol
func (a *App) GetSymbolDependents(symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return a.analysisHandler.GetSymbolDependents(a.ctx, symbolID, language, graph)
}

// Build executes project build
func (a *App) Build(projectPath, language string) (*domain.BuildResult, error) {
	return a.analysisHandler.Build(a.ctx, projectPath, language)
}

// TypeCheck performs type checking
func (a *App) TypeCheck(projectPath, language string) (*domain.TypeCheckResult, error) {
	return a.analysisHandler.TypeCheck(a.ctx, projectPath, language)
}

// BuildAndTypeCheck performs build and type checking
func (a *App) BuildAndTypeCheck(projectPath, language string) (*domain.BuildResult, *domain.TypeCheckResult, error) {
	return a.analysisHandler.BuildAndTypeCheck(a.ctx, projectPath, language)
}

// ValidateProject performs full project validation
func (a *App) ValidateProject(projectPath string, languages []string) (*domain.ProjectValidationResult, error) {
	return a.analysisHandler.ValidateProject(a.ctx, projectPath, languages)
}

// DetectLanguages detects languages in a project
func (a *App) DetectLanguages(projectPath string) ([]string, error) {
	return a.analysisHandler.DetectLanguages(a.ctx, projectPath)
}

// RunTests executes tests according to configuration
func (a *App) RunTests(config *domain.TestConfig) ([]*domain.TestResult, error) {
	return a.analysisHandler.RunTests(a.ctx, config)
}

// RunTargetedTests executes targeted tests for affected files
func (a *App) RunTargetedTests(config *domain.TestConfig, changedFiles []string) ([]*domain.TestResult, error) {
	return a.analysisHandler.RunTargetedTests(a.ctx, config, changedFiles)
}

// DiscoverTests discovers tests in a project
func (a *App) DiscoverTests(projectPath, language string) (*domain.TestSuite, error) {
	return a.analysisHandler.DiscoverTests(a.ctx, projectPath, language)
}

// BuildAffectedGraph builds a graph of affected files
func (a *App) BuildAffectedGraph(changedFiles []string, projectPath string) (*domain.AffectedGraph, error) {
	return a.analysisHandler.BuildAffectedGraph(a.ctx, changedFiles, projectPath)
}

// RunSmokeTests executes only smoke tests
func (a *App) RunSmokeTests(projectPath, language string) ([]*domain.TestResult, error) {
	return a.analysisHandler.RunSmokeTests(a.ctx, projectPath, language)
}

// RunUnitTests executes only unit tests
func (a *App) RunUnitTests(projectPath, language string) ([]*domain.TestResult, error) {
	return a.analysisHandler.RunUnitTests(a.ctx, projectPath, language)
}

// RunIntegrationTests executes only integration tests
func (a *App) RunIntegrationTests(projectPath, language string) ([]*domain.TestResult, error) {
	return a.analysisHandler.RunIntegrationTests(a.ctx, projectPath, language)
}

// ValidateTestResults validates test results
func (a *App) ValidateTestResults(results []*domain.TestResult) *domain.TestValidationResult {
	return a.analysisHandler.ValidateTestResults(results)
}

// AnalyzeProject performs static analysis on a project
func (a *App) AnalyzeProject(projectPath string, languages []string) (*domain.StaticAnalysisReport, error) {
	return a.analysisHandler.AnalyzeProject(a.ctx, projectPath, languages)
}

// AnalyzeFile performs static analysis on a single file
func (a *App) AnalyzeFile(filePath, language string) (*domain.StaticAnalysisResult, error) {
	return a.analysisHandler.AnalyzeFile(a.ctx, filePath, language)
}

// AnalyzeGoProject performs static analysis on a Go project
func (a *App) AnalyzeGoProject(projectPath string) (*domain.StaticAnalysisResult, error) {
	return a.analysisHandler.AnalyzeGoProject(a.ctx, projectPath)
}

// AnalyzeTypeScriptProject performs static analysis on a TypeScript project
func (a *App) AnalyzeTypeScriptProject(projectPath string) (*domain.StaticAnalysisResult, error) {
	return a.analysisHandler.AnalyzeTypeScriptProject(a.ctx, projectPath)
}

// AnalyzeJavaScriptProject performs static analysis on a JavaScript project
func (a *App) AnalyzeJavaScriptProject(projectPath string) (*domain.StaticAnalysisResult, error) {
	return a.analysisHandler.AnalyzeJavaScriptProject(a.ctx, projectPath)
}

// GetSupportedAnalyzers returns supported analyzers
func (a *App) GetSupportedAnalyzers() []domain.StaticAnalyzerType {
	return a.analysisHandler.GetSupportedAnalyzers()
}

// ValidateAnalysisResults validates static analysis results
func (a *App) ValidateAnalysisResults(results map[string]*domain.StaticAnalysisResult) *domain.StaticAnalysisValidationResult {
	return a.analysisHandler.ValidateAnalysisResults(results)
}

// GenerateSBOM generates SBOM for a project
func (a *App) GenerateSBOM(projectPath string, format domain.SBOMFormat) (*domain.SBOMResult, error) {
	return a.analysisHandler.GenerateSBOM(a.ctx, projectPath, format)
}

// ScanVulnerabilities scans vulnerabilities in a project
func (a *App) ScanVulnerabilities(projectPath string) (*domain.VulnerabilityScanResult, error) {
	return a.analysisHandler.ScanVulnerabilities(a.ctx, projectPath)
}

// ScanLicenses scans licenses in a project
func (a *App) ScanLicenses(projectPath string) (*domain.LicenseScanResult, error) {
	return a.analysisHandler.ScanLicenses(a.ctx, projectPath)
}

// GenerateComplianceReport generates a compliance report
func (a *App) GenerateComplianceReport(projectPath string, requirements *domain.ComplianceRequirements) (*domain.ComplianceReport, error) {
	return a.analysisHandler.GenerateComplianceReport(a.ctx, projectPath, requirements)
}

// GetSupportedSBOMFormats returns supported SBOM formats
func (a *App) GetSupportedSBOMFormats() []domain.SBOMFormat {
	return a.analysisHandler.GetSupportedSBOMFormats()
}

// ValidateSBOM validates SBOM
func (a *App) ValidateSBOM(sbomPath string, format domain.SBOMFormat) error {
	return a.analysisHandler.ValidateSBOM(a.ctx, sbomPath, format)
}

// === Project Structure Detection ===

// GetProjectStructure returns full project structure analysis
func (a *App) GetProjectStructure(projectPath string) (*domain.ProjectStructure, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.DetectStructure(projectPath)
}

// GetProjectStructureSummary returns text description of project structure
func (a *App) GetProjectStructureSummary(projectPath string) (string, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.GetArchitectureSummary(projectPath)
}

// DetectProjectArchitecture detects project architecture
func (a *App) DetectProjectArchitecture(projectPath string) (*domain.ArchitectureInfo, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.DetectArchitecture(projectPath)
}

// DetectProjectFrameworks detects project frameworks
func (a *App) DetectProjectFrameworks(projectPath string) ([]domain.FrameworkInfo, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.DetectFrameworks(projectPath)
}

// DetectProjectConventions detects project conventions
func (a *App) DetectProjectConventions(projectPath string) (*domain.ConventionInfo, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.DetectConventions(projectPath)
}

// GetRelatedLayers returns related architectural layers for a file
func (a *App) GetRelatedLayers(projectPath, filePath string) ([]domain.LayerInfo, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.GetRelatedLayers(projectPath, filePath)
}

// SuggestRelatedFiles suggests related files based on architecture
func (a *App) SuggestRelatedFiles(projectPath, filePath string) ([]string, error) {
	service := project.NewStructureServiceLazy(a.log)
	return service.SuggestRelatedFiles(projectPath, filePath)
}

// === Smart Context Suggestions ===

// SmartSuggestion represents a file suggestion with source and reason
type SmartSuggestion struct {
	Path       string  `json:"path"`
	Source     string  `json:"source"` // "git", "arch", "semantic"
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// SmartSuggestionsResult contains all suggestions grouped by source
type SmartSuggestionsResult struct {
	Suggestions []SmartSuggestion `json:"suggestions"`
	Total       int               `json:"total"`
}

// GetSmartSuggestions returns file suggestions from multiple sources
func (a *App) GetSmartSuggestions(projectPath string, currentFiles []string, task string) (*SmartSuggestionsResult, error) {
	var suggestions []SmartSuggestion
	seen := make(map[string]bool)

	// Mark current files as seen to avoid suggesting them
	for _, f := range currentFiles {
		seen[f] = true
	}

	// 1. Git-based suggestions (co-changed files)
	gitSuggestions := a.getGitSuggestions(projectPath, currentFiles, seen)
	suggestions = append(suggestions, gitSuggestions...)

	// 2. Architecture-based suggestions
	archSuggestions := a.getArchSuggestions(projectPath, currentFiles, seen)
	suggestions = append(suggestions, archSuggestions...)

	// Sort by confidence (highest first)
	sortSuggestionsByConfidence(suggestions)

	// Limit to top 15
	if len(suggestions) > 15 {
		suggestions = suggestions[:15]
	}

	return &SmartSuggestionsResult{
		Suggestions: suggestions,
		Total:       len(suggestions),
	}, nil
}

// getGitSuggestions returns suggestions based on git co-change history
func (a *App) getGitSuggestions(projectPath string, currentFiles []string, seen map[string]bool) []SmartSuggestion {
	var suggestions []SmartSuggestion
	gitContext := git.NewContextBuilder(projectPath)

	for _, file := range currentFiles {
		coChanged, err := gitContext.GetCoChangedFiles(file, 5)
		if err != nil {
			continue
		}
		for _, coFile := range coChanged {
			if seen[coFile] {
				continue
			}
			seen[coFile] = true
			suggestions = append(suggestions, SmartSuggestion{
				Path:       coFile,
				Source:     "git",
				Reason:     fmt.Sprintf("Often changed with %s", filepath.Base(file)),
				Confidence: 0.8,
			})
		}
	}
	return suggestions
}

// getArchSuggestions returns suggestions based on architecture analysis
func (a *App) getArchSuggestions(projectPath string, currentFiles []string, seen map[string]bool) []SmartSuggestion {
	var suggestions []SmartSuggestion
	service := project.NewStructureServiceLazy(a.log)

	for _, file := range currentFiles {
		related, err := service.SuggestRelatedFiles(projectPath, file)
		if err != nil {
			continue
		}
		for _, relFile := range related {
			if seen[relFile] {
				continue
			}
			seen[relFile] = true
			suggestions = append(suggestions, SmartSuggestion{
				Path:       relFile,
				Source:     "arch",
				Reason:     "Related architectural layer",
				Confidence: 0.7,
			})
		}
	}
	return suggestions
}

// sortSuggestionsByConfidence sorts suggestions by confidence descending
func sortSuggestionsByConfidence(suggestions []SmartSuggestion) {
	for i := 0; i < len(suggestions)-1; i++ {
		for j := i + 1; j < len(suggestions); j++ {
			if suggestions[j].Confidence > suggestions[i].Confidence {
				suggestions[i], suggestions[j] = suggestions[j], suggestions[i]
			}
		}
	}
}

// === File Quick Info (Phase 4) ===

// FileQuickInfo contains quick statistics about a file
type FileQuickInfo struct {
	SymbolCount    int     `json:"symbolCount"`
	ImportCount    int     `json:"importCount"`
	DependentCount int     `json:"dependentCount"`
	ChangeRisk     float64 `json:"changeRisk"`
	RiskLevel      string  `json:"riskLevel"` // "low", "medium", "high"
}

// GetFileQuickInfo returns quick statistics for a file
func (a *App) GetFileQuickInfo(projectPath, filePath string) (*FileQuickInfo, error) {
	info := &FileQuickInfo{}

	// Get symbol count using symbol index
	service := project.NewStructureServiceLazy(a.log)
	symbols, _ := service.GetFileSymbols(projectPath, filePath)
	info.SymbolCount = len(symbols)

	// Get import count
	imports, _ := service.GetFileImports(projectPath, filePath)
	info.ImportCount = len(imports)

	// Get dependent files count
	dependents, _ := service.GetDependentFiles(projectPath, filePath)
	info.DependentCount = len(dependents)

	// Calculate change risk based on dependents
	info.ChangeRisk = calculateRisk(info.DependentCount, info.SymbolCount)
	info.RiskLevel = getRiskLevel(info.ChangeRisk)

	return info, nil
}

func calculateRisk(dependents, symbols int) float64 {
	if dependents == 0 {
		return 0.1
	}
	risk := float64(dependents) / 20.0 // normalize to 0-1 range
	if risk > 1.0 {
		risk = 1.0
	}
	// Adjust by symbol count (more symbols = more risk)
	if symbols > 20 {
		risk += 0.1
	}
	if risk > 1.0 {
		risk = 1.0
	}
	return risk
}

func getRiskLevel(risk float64) string {
	if risk < 0.3 {
		return "low"
	}
	if risk < 0.7 {
		return "medium"
	}
	return "high"
}

// === Impact Preview (Phase 5) ===

// ImpactPreviewResult contains impact analysis for selected files
type ImpactPreviewResult struct {
	TotalDependents int            `json:"totalDependents"`
	AggregateRisk   float64        `json:"aggregateRisk"`
	RiskLevel       string         `json:"riskLevel"`
	AffectedFiles   []AffectedFile `json:"affectedFiles"`
	RelatedTests    []string       `json:"relatedTests"`
}

// AffectedFile represents a file affected by changes
type AffectedFile struct {
	Path       string `json:"path"`
	Type       string `json:"type"` // "direct", "transitive"
	Dependents int    `json:"dependents"`
}

// GetImpactPreview returns impact analysis for selected files
func (a *App) GetImpactPreview(projectPath string, filePaths []string) (*ImpactPreviewResult, error) {
	result := &ImpactPreviewResult{
		AffectedFiles: []AffectedFile{},
		RelatedTests:  []string{},
	}

	service := project.NewStructureServiceLazy(a.log)
	seen := make(map[string]bool)
	var totalRisk float64

	for _, filePath := range filePaths {
		seen[filePath] = true
	}

	// Collect all dependents
	for _, filePath := range filePaths {
		dependents, err := service.GetDependentFiles(projectPath, filePath)
		if err != nil {
			continue
		}

		for _, dep := range dependents {
			if seen[dep] {
				continue
			}
			seen[dep] = true

			depType := "direct"
			result.AffectedFiles = append(result.AffectedFiles, AffectedFile{
				Path: dep,
				Type: depType,
			})

			// Check if it's a test file
			if isTestFile(dep) {
				result.RelatedTests = append(result.RelatedTests, dep)
			}
		}

		// Calculate risk for this file
		info, _ := a.GetFileQuickInfo(projectPath, filePath)
		if info != nil {
			totalRisk += info.ChangeRisk
		}
	}

	result.TotalDependents = len(result.AffectedFiles)
	if len(filePaths) > 0 {
		result.AggregateRisk = totalRisk / float64(len(filePaths))
	}
	result.RiskLevel = getRiskLevel(result.AggregateRisk)

	// Limit affected files to 20
	if len(result.AffectedFiles) > 20 {
		result.AffectedFiles = result.AffectedFiles[:20]
	}

	return result, nil
}

func isTestFile(path string) bool {
	return strings.Contains(path, "_test.") ||
		strings.Contains(path, ".test.") ||
		strings.Contains(path, ".spec.") ||
		strings.Contains(path, "/tests/") ||
		strings.Contains(path, "/__tests__/")
}

// === Memory/Context API (Phase 6) ===

// ContextMemoryEntry represents a saved context
type ContextMemoryEntry struct {
	ID        string   `json:"id"`
	Topic     string   `json:"topic"`
	Summary   string   `json:"summary"`
	Files     []string `json:"files"`
	CreatedAt string   `json:"createdAt"`
}

// GetRecentContexts returns recently saved contexts
func (a *App) GetRecentContexts(projectPath string, limit int) ([]ContextMemoryEntry, error) {
	if a.analysisContainer == nil {
		return []ContextMemoryEntry{}, nil
	}

	contextMemory := a.analysisContainer.GetContextMemory()
	if contextMemory == nil {
		return []ContextMemoryEntry{}, nil
	}

	contexts, err := contextMemory.GetRecentContexts(projectPath, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent contexts: %w", err)
	}

	result := make([]ContextMemoryEntry, 0, len(contexts))
	for _, ctx := range contexts {
		result = append(result, ContextMemoryEntry{
			ID:        ctx.ID,
			Topic:     ctx.Topic,
			Summary:   ctx.Summary,
			Files:     ctx.Files,
			CreatedAt: ctx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

// FindContextByTopic searches contexts by topic
func (a *App) FindContextByTopic(projectPath, topic string) ([]ContextMemoryEntry, error) {
	if a.analysisContainer == nil {
		return []ContextMemoryEntry{}, nil
	}

	contextMemory := a.analysisContainer.GetContextMemory()
	if contextMemory == nil {
		return []ContextMemoryEntry{}, nil
	}

	contexts, err := contextMemory.FindContextByTopic(projectPath, topic)
	if err != nil {
		return nil, fmt.Errorf("failed to find contexts: %w", err)
	}

	result := make([]ContextMemoryEntry, 0, len(contexts))
	for _, ctx := range contexts {
		result = append(result, ContextMemoryEntry{
			ID:        ctx.ID,
			Topic:     ctx.Topic,
			Summary:   ctx.Summary,
			Files:     ctx.Files,
			CreatedAt: ctx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return result, nil
}

// SaveContextMemory saves current context to memory
func (a *App) SaveContextMemory(projectPath, topic, summary string, files []string) error {
	if a.analysisContainer == nil {
		return fmt.Errorf("analysis container not initialized")
	}

	contextMemory := a.analysisContainer.GetContextMemory()
	if contextMemory == nil {
		return fmt.Errorf("context memory not initialized")
	}

	ctx := &domain.ConversationContext{
		ID:          fmt.Sprintf("ctx_%d", time.Now().UnixNano()),
		ProjectRoot: projectPath,
		Topic:       topic,
		Summary:     summary,
		Files:       files,
		CreatedAt:   time.Now(),
	}

	return contextMemory.SaveContext(ctx)
}

// === Repair Service ===

// ExecuteRepair executes repair cycle
func (a *App) ExecuteRepair(projectPath, errorOutput, language string, maxAttempts int) (*domain.RepairResult, error) {
	req := domain.RepairRequest{
		ProjectPath: projectPath,
		ErrorOutput: errorOutput,
		Language:    language,
		MaxAttempts: maxAttempts,
	}
	return a.repairService.ExecuteRepair(a.ctx, req)
}

// GetAvailableRepairRules returns available rules for a language
func (a *App) GetAvailableRepairRules(language string) ([]domain.RepairRule, error) {
	return a.repairService.GetAvailableRules(language)
}

// AddRepairRule adds a new rule
func (a *App) AddRepairRule(rule domain.RepairRule) error {
	return a.repairService.AddRule(rule)
}

// RemoveRepairRule removes a rule
func (a *App) RemoveRepairRule(ruleID string) error {
	return a.repairService.RemoveRule(ruleID)
}

// ValidateRepairRule validates a rule
func (a *App) ValidateRepairRule(rule domain.RepairRule) error {
	return a.repairService.ValidateRule(rule)
}

// === Guardrail Service ===

// ValidatePath validates a path against policies
func (a *App) ValidatePath(path string) ([]domain.GuardrailViolation, error) {
	return a.guardrailService.ValidatePath(path)
}

// ValidateBudget validates budget constraints
func (a *App) ValidateBudget(budgetType domain.BudgetType, current int64) ([]domain.BudgetViolation, error) {
	return a.guardrailService.ValidateBudget(budgetType, current)
}

// GetGuardrailPolicies returns all policies
func (a *App) GetGuardrailPolicies() ([]domain.GuardrailPolicy, error) {
	return a.guardrailService.GetPolicies()
}

// GetBudgetPolicies returns budget policies
func (a *App) GetBudgetPolicies() ([]domain.BudgetPolicy, error) {
	return a.guardrailService.GetBudgetPolicies()
}

// AddGuardrailPolicy adds a new policy
func (a *App) AddGuardrailPolicy(policy domain.GuardrailPolicy) error {
	return a.guardrailService.AddPolicy(policy)
}

// RemoveGuardrailPolicy removes a policy
func (a *App) RemoveGuardrailPolicy(policyID string) error {
	return a.guardrailService.RemovePolicy(policyID)
}

// UpdateGuardrailPolicy updates a policy
func (a *App) UpdateGuardrailPolicy(policy domain.GuardrailPolicy) error {
	return a.guardrailService.UpdatePolicy(policy)
}

// AddBudgetPolicy adds a budget policy
func (a *App) AddBudgetPolicy(policy domain.BudgetPolicy) error {
	return a.guardrailService.AddBudgetPolicy(policy)
}

// RemoveBudgetPolicy removes a budget policy
func (a *App) RemoveBudgetPolicy(policyID string) error {
	return a.guardrailService.RemoveBudgetPolicy(policyID)
}

// UpdateBudgetPolicy updates a budget policy
func (a *App) UpdateBudgetPolicy(policy domain.BudgetPolicy) error {
	return a.guardrailService.UpdateBudgetPolicy(policy)
}
