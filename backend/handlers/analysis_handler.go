package handlers

import (
	"context"
	"syntaxia/application/sbom"
	"syntaxia/application/symbol"
	"syntaxia/domain"
)

// AnalysisHandler handles static analysis, testing, and build operations
type AnalysisHandler struct {
	log                   domain.Logger
	testService           domain.ITestService
	staticAnalyzerService domain.IStaticAnalyzerService
	buildService          domain.IBuildService
	sbomService           *sbom.Service
	symbolGraph           *symbol.Service

	// Semaphore for limiting concurrent analysis operations
	sem chan struct{}
}

const maxConcurrentAnalysis = 4

// NewAnalysisHandler creates a new analysis handler
func NewAnalysisHandler(
	log domain.Logger,
	testService domain.ITestService,
	staticAnalyzerService domain.IStaticAnalyzerService,
	buildService domain.IBuildService,
	sbomService *sbom.Service,
	symbolGraph *symbol.Service,
) *AnalysisHandler {
	return &AnalysisHandler{
		log:                   log,
		testService:           testService,
		staticAnalyzerService: staticAnalyzerService,
		buildService:          buildService,
		sbomService:           sbomService,
		symbolGraph:           symbolGraph,
		sem:                   make(chan struct{}, maxConcurrentAnalysis),
	}
}

// acquireSem acquires a semaphore slot
func (h *AnalysisHandler) acquireSem(ctx context.Context) error {
	select {
	case h.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// releaseSem releases a semaphore slot
func (h *AnalysisHandler) releaseSem() {
	<-h.sem
}

// === Build Operations ===

// Build executes project build
func (h *AnalysisHandler) Build(ctx context.Context, projectPath, language string) (*domain.BuildResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.buildService.Build(ctx, projectPath, language)
}

// TypeCheck performs type checking
func (h *AnalysisHandler) TypeCheck(ctx context.Context, projectPath, language string) (*domain.TypeCheckResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.buildService.TypeCheck(ctx, projectPath, language)
}

// BuildAndTypeCheck executes build and type checking
func (h *AnalysisHandler) BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*domain.BuildResult, *domain.TypeCheckResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, nil, err
	}
	defer h.releaseSem()

	return h.buildService.BuildAndTypeCheck(ctx, projectPath, language)
}

// ValidateProject performs full project validation
func (h *AnalysisHandler) ValidateProject(ctx context.Context, projectPath string, languages []string) (*domain.ProjectValidationResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.buildService.ValidateProject(ctx, projectPath, languages)
}

// DetectLanguages detects languages in project
func (h *AnalysisHandler) DetectLanguages(ctx context.Context, projectPath string) ([]string, error) {
	return h.buildService.DetectLanguages(ctx, projectPath)
}

// === Test Operations ===

// RunTests executes tests
func (h *AnalysisHandler) RunTests(ctx context.Context, config *domain.TestConfig) ([]*domain.TestResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.testService.RunTests(ctx, config)
}

// RunTargetedTests executes targeted tests for affected files
func (h *AnalysisHandler) RunTargetedTests(ctx context.Context, config *domain.TestConfig, changedFiles []string) ([]*domain.TestResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.testService.RunTargetedTests(ctx, config, changedFiles)
}

// DiscoverTests discovers tests in project
func (h *AnalysisHandler) DiscoverTests(ctx context.Context, projectPath, language string) (*domain.TestSuite, error) {
	return h.testService.DiscoverTests(ctx, projectPath, language)
}

// BuildAffectedGraph builds affected files graph
func (h *AnalysisHandler) BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*domain.AffectedGraph, error) {
	return h.testService.BuildAffectedGraph(ctx, changedFiles, projectPath)
}

// RunSmokeTests executes smoke tests
func (h *AnalysisHandler) RunSmokeTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.testService.RunSmokeTests(ctx, projectPath, language)
}

// RunUnitTests executes unit tests
func (h *AnalysisHandler) RunUnitTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.testService.RunUnitTests(ctx, projectPath, language)
}

// RunIntegrationTests executes integration tests
func (h *AnalysisHandler) RunIntegrationTests(ctx context.Context, projectPath, language string) ([]*domain.TestResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.testService.RunIntegrationTests(ctx, projectPath, language)
}

// ValidateTestResults validates test results
func (h *AnalysisHandler) ValidateTestResults(results []*domain.TestResult) *domain.TestValidationResult {
	return h.testService.ValidateTestResults(results)
}

// === Static Analysis Operations ===

// AnalyzeProject performs static analysis
func (h *AnalysisHandler) AnalyzeProject(ctx context.Context, projectPath string, languages []string) (*domain.StaticAnalysisReport, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.staticAnalyzerService.AnalyzeProject(ctx, projectPath, languages)
}

// AnalyzeFile analyzes a single file
func (h *AnalysisHandler) AnalyzeFile(ctx context.Context, filePath, language string) (*domain.StaticAnalysisResult, error) {
	return h.staticAnalyzerService.AnalyzeFile(ctx, filePath, language)
}

// AnalyzeGoProject analyzes Go project
func (h *AnalysisHandler) AnalyzeGoProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.staticAnalyzerService.AnalyzeGoProject(ctx, projectPath)
}

// AnalyzeTypeScriptProject analyzes TypeScript project
func (h *AnalysisHandler) AnalyzeTypeScriptProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.staticAnalyzerService.AnalyzeTypeScriptProject(ctx, projectPath)
}

// AnalyzeJavaScriptProject analyzes JavaScript project
func (h *AnalysisHandler) AnalyzeJavaScriptProject(ctx context.Context, projectPath string) (*domain.StaticAnalysisResult, error) {
	if err := h.acquireSem(ctx); err != nil {
		return nil, err
	}
	defer h.releaseSem()

	return h.staticAnalyzerService.AnalyzeJavaScriptProject(ctx, projectPath)
}

// GetSupportedAnalyzers returns supported analyzers
func (h *AnalysisHandler) GetSupportedAnalyzers() []domain.StaticAnalyzerType {
	return h.staticAnalyzerService.GetSupportedAnalyzers()
}

// ValidateAnalysisResults validates analysis results
func (h *AnalysisHandler) ValidateAnalysisResults(results map[string]*domain.StaticAnalysisResult) *domain.StaticAnalysisValidationResult {
	return h.staticAnalyzerService.ValidateAnalysisResults(results)
}

// === SBOM Operations ===

// GenerateSBOM generates SBOM for project
func (h *AnalysisHandler) GenerateSBOM(ctx context.Context, projectPath string, format domain.SBOMFormat) (*domain.SBOMResult, error) {
	return h.sbomService.GenerateSBOM(ctx, projectPath, format)
}

// ScanVulnerabilities scans vulnerabilities
func (h *AnalysisHandler) ScanVulnerabilities(ctx context.Context, projectPath string) (*domain.VulnerabilityScanResult, error) {
	return h.sbomService.ScanVulnerabilities(ctx, projectPath)
}

// ScanLicenses scans licenses
func (h *AnalysisHandler) ScanLicenses(ctx context.Context, projectPath string) (*domain.LicenseScanResult, error) {
	return h.sbomService.ScanLicenses(ctx, projectPath)
}

// GenerateComplianceReport generates compliance report
func (h *AnalysisHandler) GenerateComplianceReport(ctx context.Context, projectPath string, requirements *domain.ComplianceRequirements) (*domain.ComplianceReport, error) {
	return h.sbomService.GenerateComplianceReport(ctx, projectPath, requirements)
}

// GetSupportedSBOMFormats returns supported SBOM formats
func (h *AnalysisHandler) GetSupportedSBOMFormats() []domain.SBOMFormat {
	return h.sbomService.GetSupportedSBOMFormats()
}

// ValidateSBOM validates SBOM
func (h *AnalysisHandler) ValidateSBOM(ctx context.Context, sbomPath string, format domain.SBOMFormat) error {
	return h.sbomService.ValidateSBOM(ctx, sbomPath, format)
}

// === Symbol Graph Operations ===

// BuildSymbolGraph builds symbol graph for project
func (h *AnalysisHandler) BuildSymbolGraph(ctx context.Context, projectRoot, language string) (*domain.SymbolGraph, error) {
	return h.symbolGraph.BuildSymbolGraph(ctx, projectRoot, language)
}

// GetSymbolSuggestions returns symbol suggestions
func (h *AnalysisHandler) GetSymbolSuggestions(ctx context.Context, query, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return h.symbolGraph.GetSuggestions(ctx, query, language, graph)
}

// GetSymbolDependencies returns symbol dependencies
func (h *AnalysisHandler) GetSymbolDependencies(ctx context.Context, symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return h.symbolGraph.GetDependencies(ctx, symbolID, language, graph)
}

// GetSymbolDependents returns symbols depending on the specified one
func (h *AnalysisHandler) GetSymbolDependents(ctx context.Context, symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return h.symbolGraph.GetDependents(ctx, symbolID, language, graph)
}
