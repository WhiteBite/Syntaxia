package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai"
	"syntaxia/infrastructure/analyzers"
	"syntaxia/infrastructure/applyengine"
	"syntaxia/infrastructure/contextbuilder"
	"syntaxia/infrastructure/contentoptimizer"
	execinfra "syntaxia/infrastructure/exec"
	"syntaxia/infrastructure/filereader"
	"syntaxia/infrastructure/formatters"
	"syntaxia/infrastructure/fsscanner"
	"syntaxia/infrastructure/fswatcher"
	"syntaxia/infrastructure/git"
	"syntaxia/infrastructure/reportfs"
	"syntaxia/infrastructure/sbomlicensing"
	"syntaxia/infrastructure/settingsfs"
	"syntaxia/infrastructure/shellintegration"
	"syntaxia/infrastructure/staticanalyzer"
	"syntaxia/infrastructure/taskflowrepo"
	"syntaxia/infrastructure/testengine"
	"syntaxia/infrastructure/textutils"
	"syntaxia/infrastructure/uxreports"
	"syntaxia/infrastructure/wailsbridge"
	"syntaxia/internal/executil"
	"time"

	archiverinfra "syntaxia/infrastructure/archiver"
	"syntaxia/infrastructure/buildpipeline"
	"syntaxia/infrastructure/diffengine"
	"syntaxia/infrastructure/pdfgen"
	"syntaxia/infrastructure/policy"
	"syntaxia/infrastructure/symbolgraph"
	"syntaxia/application/router"
)

// createRouterLLMConfig creates configuration for router LLM service
func createRouterLLMConfig() router.LLMConfig {
	return router.LLMConfig{
		Enabled: false, // Disabled by default
		LLMConfig: domain.LLMConfig{
			BaseURL:       "http://localhost:8080",
			Timeout:       30 * time.Second,
			MaxTokens:     2048,
			Temperature:   0.7,
			TopP:          0.9,
			TopK:          40,
			RepeatPenalty: 1.1,
		},
		FallbackToHeuristic: true,
		MaxRetries:          3,
		Timeout:             30 * time.Second,
	}
}

// initializeInfrastructure initializes all infrastructure layer components
func (c *AppContainer) initializeInfrastructure(ctx context.Context, embeddedIgnoreGlob, defaultCustomPrompt string) error {
	var err error

	// Bridge for Wails (Logger and EventBus)
	bridge := wailsbridge.New(ctx)
	c.Bridge = bridge
	c.Log = bridge
	c.Bus = bridge

	// Repositories and Infrastructure
	c.SettingsRepo, err = settingsfs.New(c.Log, embeddedIgnoreGlob, defaultCustomPrompt)
	if err != nil {
		return err
	}
	c.FileReader = filereader.NewSecureFileReader(c.Log)
	c.GitRepo = git.New(c.Log)
	c.TreeBuilder = fsscanner.New(c.SettingsRepo, c.Log)
	c.FileSearcher = fsscanner.NewFileSearcher(c.TreeBuilder, c.Log)
	c.ContextSplitter = textutils.NewContextSplitter(c.Log)

	c.Watcher, err = fswatcher.New(ctx, c.Bus)
	if err != nil {
		return err
	}
	c.CommandRunner = execinfra.NewCommandRunnerImpl(c.Log)

	// Create symbol graph builders
	goSymbolGraphBuilder := symbolgraph.NewGoSymbolGraphBuilder(c.Log)
	symbolGraphBuilders := make(map[string]domain.SymbolGraphBuilder)
	symbolGraphBuilders["go"] = goSymbolGraphBuilder

	// Create import graph builders
	importGraphBuilders := make(map[string]domain.ImportGraphBuilder)

	// Create analyzer registry for TS/JS/Vue import graph builder
	tsAnalyzerRegistry := analyzers.NewAnalyzerRegistry()
	tsImportGraphBuilder := symbolgraph.NewTSImportGraphBuilder(c.Log, tsAnalyzerRegistry)
	importGraphBuilders["typescript"] = tsImportGraphBuilder
	importGraphBuilders["javascript"] = tsImportGraphBuilder
	importGraphBuilders["vue"] = tsImportGraphBuilder

	// Create TestService with lazy initialization
	c.testServiceOnce.Do(func() {
		testEngine := testengine.NewTestEngine(c.Log, goSymbolGraphBuilder)
		testEngine.RegisterTestRunner("go", testengine.NewGoTestRunner(c.Log))
		testEngine.RegisterTestRunner("python", testengine.NewPythonTestRunner(c.Log))

		// TypeScript/JavaScript test runner (same runner for both)
		tsRunner := testengine.NewTypeScriptTestRunner(c.Log)
		testEngine.RegisterTestRunner("typescript", tsRunner)
		testEngine.RegisterTestRunner("javascript", tsRunner)

		// Java test runner
		testEngine.RegisterTestRunner("java", testengine.NewJavaTestRunner(c.Log))

		// Rust test runner
		testEngine.RegisterTestRunner("rust", testengine.NewRustTestRunner(c.Log))

		// Kotlin test runner
		testEngine.RegisterTestRunner("kotlin", testengine.NewKotlinTestRunner(c.Log))

		// C# test runner
		testEngine.RegisterTestRunner("csharp", testengine.NewCSharpTestRunner(c.Log))

		// Register test analyzers for supported languages
		testEngine.RegisterTestAnalyzer("go", testengine.NewGoTestAnalyzer(c.Log))
		testEngine.RegisterTestAnalyzer("python", testengine.NewPythonTestAnalyzer(c.Log))

		// TypeScript/JavaScript test analyzer (same analyzer for both)
		tsAnalyzer := testengine.NewTypeScriptTestAnalyzer(c.Log)
		testEngine.RegisterTestAnalyzer("typescript", tsAnalyzer)
		testEngine.RegisterTestAnalyzer("javascript", tsAnalyzer)

		// Java test analyzer
		testEngine.RegisterTestAnalyzer("java", testengine.NewJavaTestAnalyzer(c.Log))

		// Rust test analyzer
		testEngine.RegisterTestAnalyzer("rust", testengine.NewRustTestAnalyzer(c.Log))

		// Kotlin test analyzer
		testEngine.RegisterTestAnalyzer("kotlin", testengine.NewKotlinTestAnalyzer(c.Log))

		// C# test analyzer
		testEngine.RegisterTestAnalyzer("csharp", testengine.NewCSharpTestAnalyzer(c.Log))

		// Dart test runner and analyzer
		testEngine.RegisterTestRunner("dart", testengine.NewDartTestRunner(c.Log))
		testEngine.RegisterTestAnalyzer("dart", testengine.NewDartTestAnalyzer(c.Log))

		// Ruby test runner and analyzer
		testEngine.RegisterTestRunner("ruby", testengine.NewRubyTestRunner(c.Log))
		testEngine.RegisterTestAnalyzer("ruby", testengine.NewRubyTestAnalyzer(c.Log))

		// C++ test runner and analyzer
		testEngine.RegisterTestRunner("cpp", testengine.NewCppTestRunner(c.Log))
		testEngine.RegisterTestAnalyzer("cpp", testengine.NewCppTestAnalyzer(c.Log))

		// Swift test runner and analyzer
		testEngine.RegisterTestRunner("swift", testengine.NewSwiftTestRunner(c.Log))
		testEngine.RegisterTestAnalyzer("swift", testengine.NewSwiftTestAnalyzer(c.Log))

		// PHP test runner and analyzer
		testEngine.RegisterTestRunner("php", testengine.NewPHPTestRunner(c.Log))
		testEngine.RegisterTestAnalyzer("php", testengine.NewPHPTestAnalyzer(c.Log))
	})

	// Create Static Analyzer Engine and infrastructure components
	staticAnalyzerEngine := staticanalyzer.NewStaticAnalyzerEngine(c.Log)
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewStaticcheckAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewESLintAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewErrorProneAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewRuffAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewClangTidyAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewClippyAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewKtlintAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewDotnetFormatAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewDartAnalyzeAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewSwiftLintAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewPHPCSAnalyzer(c.Log))
	staticAnalyzerEngine.RegisterAnalyzer(staticanalyzer.NewRuboCopAnalyzer(c.Log))

	// Create SBOM infrastructure components
	sbomGenerator := sbomlicensing.NewSyftGenerator(c.Log)
	vulnScanner := sbomlicensing.NewGrypeScanner(c.Log)
	licenseScanner := sbomlicensing.NewLicenseScanner(c.Log)

	// Create TaskflowRepository
	taskflowRepo := taskflowrepo.NewFileSystemTaskflowRepository("tasks/status.json")

	// Create OPA service and file stat provider for GuardrailService
	guardrailOPAService := policy.NewOPAService(c.Log)

	// Create apply engine
	applyConfig := &domain.ApplyEngineConfig{
		AutoFormat:     true,
		AutoFixImports: true,
		BackupFiles:    true,
		ValidateAfter:  true,
		Languages:      []string{"go", "typescript", "ts"},
	}
	applyEngine := applyengine.NewApplyEngine(c.Log, applyConfig)

	// Create formatters
	formatterMap := map[string]domain.Formatter{
		"go":         formatters.NewGoFormatter(c.Log),
		"typescript": formatters.NewTypeScriptFormatter(c.Log),
		"ts":         formatters.NewTypeScriptFormatter(c.Log),
	}

	// Create import fixers
	importFixerMap := map[string]domain.ImportFixer{
		"go":         formatters.NewGoFormatter(c.Log),
		"typescript": formatters.NewTypeScriptFormatter(c.Log),
		"ts":         formatters.NewTypeScriptFormatter(c.Log),
	}

	// Create diff engine
	diffEngine := diffengine.NewDiffEngine(c.Log)

	// Create build pipeline
	buildPipeline := buildpipeline.NewBuildPipeline(c.Log)

	// Create PDF and ZIP implementations
	pdfGen := pdfgen.NewGofpdfGenerator(c.Log)
	arch := archiverinfra.NewZipArchiver(c.Log)
	contextFormatter := contextbuilder.NewContextFormatter()

	// Create ContentOptimizer for noise reduction
	contentOptimizer := contentoptimizer.NewContentOptimizer(c.Log)

	// Create UXReportRepository
	uxReportRepo := uxreports.NewFileSystemUXReportRepository("reports/ux")

	// Create report repository
	reportRepo, err := reportfs.NewReportFileSystemRepository(c.Log)
	if err != nil {
		return fmt.Errorf("failed to initialize report repository: %w", err)
	}

	// Create LLM client for router
	routerLLMConfig := createRouterLLMConfig()
	llamaClient := ai.NewLlamaCppClient(ai.LlamaCppConfig{
		BaseURL:       routerLLMConfig.LLMConfig.BaseURL,
		Timeout:       routerLLMConfig.LLMConfig.Timeout,
		MaxTokens:     routerLLMConfig.LLMConfig.MaxTokens,
		Temperature:   routerLLMConfig.LLMConfig.Temperature,
		TopP:          routerLLMConfig.LLMConfig.TopP,
		TopK:          routerLLMConfig.LLMConfig.TopK,
		RepeatPenalty: routerLLMConfig.LLMConfig.RepeatPenalty,
	}, c.Log)
	llmClient := ai.NewLlamaCppClientAdapter(llamaClient)
	fileReader := filereader.NewFileReader()

	// Initialize Shell Integration (OS context menu)
	c.ShellIntegration = shellintegration.NewService("Syntaxia")

	// Store infrastructure components for service initialization
	c.infrastructureComponents = &infrastructureComponents{
		goSymbolGraphBuilder:   goSymbolGraphBuilder,
		symbolGraphBuilders:    symbolGraphBuilders,
		importGraphBuilders:    importGraphBuilders,
		staticAnalyzerEngine:   staticAnalyzerEngine,
		sbomGenerator:          sbomGenerator,
		vulnScanner:            vulnScanner,
		licenseScanner:         licenseScanner,
		taskflowRepo:           taskflowRepo,
		guardrailOPAService:    guardrailOPAService,
		applyConfig:            applyConfig,
		applyEngine:            applyEngine,
		formatterMap:           formatterMap,
		importFixerMap:         importFixerMap,
		diffEngine:             diffEngine,
		buildPipeline:          buildPipeline,
		pdfGen:                 pdfGen,
		arch:                   arch,
		contextFormatter:       contextFormatter,
		contentOptimizer:       contentOptimizer,
		uxReportRepo:           uxReportRepo,
		reportRepo:             reportRepo,
		routerLLMConfig:        routerLLMConfig,
		llmClient:              llmClient,
		fileReader:             fileReader,
	}

	return nil
}

// infrastructureComponents holds infrastructure layer components
type infrastructureComponents struct {
	goSymbolGraphBuilder   domain.SymbolGraphBuilder
	symbolGraphBuilders    map[string]domain.SymbolGraphBuilder
	importGraphBuilders    map[string]domain.ImportGraphBuilder
	staticAnalyzerEngine   *staticanalyzer.StaticAnalyzerEngineImpl
	sbomGenerator          domain.SBOMGenerator
	vulnScanner            domain.VulnerabilityScanner
	licenseScanner         domain.LicenseScanner
	taskflowRepo           domain.TaskflowRepository
	guardrailOPAService    *policy.OPAService
	applyConfig            *domain.ApplyEngineConfig
	applyEngine            domain.ApplyEngine
	formatterMap           map[string]domain.Formatter
	importFixerMap         map[string]domain.ImportFixer
	diffEngine             domain.DiffEngine
	buildPipeline          domain.BuildPipeline
	pdfGen                 domain.PDFGenerator
	arch                   domain.Archiver
	contextFormatter       domain.ContextFormatter
	contentOptimizer       domain.ContentOptimizer
	uxReportRepo           domain.UXReportRepository
	reportRepo             domain.ReportRepository
	routerLLMConfig        router.LLMConfig
	llmClient              domain.LLMClient
	fileReader             domain.FileReader
}

// FilePathProvider implements domain.PathProvider using standard filepath functions
type FilePathProvider struct{}

func (p *FilePathProvider) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (p *FilePathProvider) Base(path string) string {
	return filepath.Base(path)
}

func (p *FilePathProvider) Dir(path string) string {
	return filepath.Dir(path)
}

func (p *FilePathProvider) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

func (p *FilePathProvider) Clean(path string) string {
	return filepath.Clean(path)
}

func (p *FilePathProvider) Getwd() (string, error) {
	return os.Getwd()
}

// OSFileSystemWriter implements domain.FileSystemWriter using standard os functions
type OSFileSystemWriter struct{}

func (w *OSFileSystemWriter) WriteFile(filename string, data []byte, perm int) error {
	return os.WriteFile(filename, data, os.FileMode(perm))
}

func (w *OSFileSystemWriter) MkdirAll(path string, perm int) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

func (w *OSFileSystemWriter) Remove(name string) error {
	return os.Remove(name)
}

func (w *OSFileSystemWriter) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

// OSTempFileProvider implements domain.TempFileProvider using standard os functions
type OSTempFileProvider struct{}

func (p *OSTempFileProvider) MkdirTemp(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// OSFileStatProvider implements domain.FileStatProvider using standard os functions
type OSFileStatProvider struct{}

func (p *OSFileStatProvider) Stat(name string) (domain.FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

// OSFileSystemProvider implements domain.FileSystemProvider
type OSFileSystemProvider struct{}

func (o *OSFileSystemProvider) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

func (o *OSFileSystemProvider) WriteFile(filename string, data []byte, perm int) error {
	return os.WriteFile(filename, data, os.FileMode(perm))
}

func (o *OSFileSystemProvider) MkdirAll(path string, perm int) error {
	return os.MkdirAll(path, os.FileMode(perm))
}

// CommandRunnerImpl implements domain.CommandRunner
type CommandRunnerImpl struct{}

func (c *CommandRunnerImpl) RunCommand(_ context.Context, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	executil.HideWindow(cmd)
	return cmd.Output()
}

func (c *CommandRunnerImpl) RunCommandInDir(_ context.Context, dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	executil.HideWindow(cmd)
	cmd.Dir = dir
	return cmd.Output()
}

// SimpleTokenCounter provides basic token estimation
type SimpleTokenCounter struct{}

func (s *SimpleTokenCounter) CountTokens(text string) int {
	// Simple approximation: 1 token ≈ 4 characters
	return len(text) / 4
}
