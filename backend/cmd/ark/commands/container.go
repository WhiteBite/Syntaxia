package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	appai "syntaxia/application/ai"
	"syntaxia/application/analysis"
	"syntaxia/application/build"
	"syntaxia/application/diff"
	"syntaxia/application/export"
	"syntaxia/application/guardrails"
	"syntaxia/application/repair"
	"syntaxia/application/sbom"
	"syntaxia/application/settings"
	"syntaxia/application/symbol"
	"syntaxia/application/ux"
	"syntaxia/application/verification"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai"
	"syntaxia/infrastructure/contextbuilder"
	"syntaxia/infrastructure/exec"
	"syntaxia/infrastructure/filereader"
	"syntaxia/infrastructure/filesystem"
	"syntaxia/infrastructure/formatters"
	"syntaxia/infrastructure/fsscanner"
	"syntaxia/infrastructure/git"
	"syntaxia/infrastructure/policy"
	"syntaxia/infrastructure/sbomlicensing"
	"syntaxia/infrastructure/settingsfs"
	"syntaxia/infrastructure/textutils"
	"syntaxia/infrastructure/uxreports"

	// Internal services (unified architecture)
	contextservice "syntaxia/internal/context"
	projectservice "syntaxia/internal/project"

	// new wiring
	"syntaxia/infrastructure/applyengine"
	archiverinfra "syntaxia/infrastructure/archiver"
	"syntaxia/infrastructure/buildpipeline"
	"syntaxia/infrastructure/diffengine"
	"syntaxia/infrastructure/pdfgen"
	"syntaxia/infrastructure/staticanalyzer"
	"syntaxia/infrastructure/symbolgraph"
	"syntaxia/infrastructure/testengine"
)

const openRouterHost = "https://openrouter.ai/api/v1"

// CLIContainer holds all the services and repositories for the application.
type CLIContainer struct {
	Log                   domain.Logger
	EventBus              domain.EventBus
	SettingsRepo          domain.SettingsRepository
	FileReader            domain.FileContentReader
	GitRepo               domain.GitRepository
	TreeBuilder           domain.TreeBuilder
	ContextSplitter       domain.ContextSplitter
	CommandRunner         domain.CommandRunner
	SettingsService       *settings.Service
	ContextService        *contextservice.Service
	ProjectService        *projectservice.Service
	AIService             *appai.Service
	ContextAnalysis       domain.ContextAnalyzer
	SymbolGraph           *symbol.Service
	TestService           domain.ITestService
	StaticAnalyzerService domain.IStaticAnalyzerService
	SBOMService           *sbom.Service
	RepairService         domain.RepairService
	GuardrailService      domain.GuardrailService
	TaskflowService       domain.TaskflowService
	UXMetricsService      domain.UXMetricsService
	ApplyService          *diff.ApplyService
	DiffService           *diff.Service
	BuildService          domain.IBuildService
	ExportService         *export.Service
	VerificationService   *verification.Service
	opaService            domain.OPAService
}

// NewCLIContainer creates and wires up all the application dependencies.
func NewCLIContainer(ctx context.Context, embeddedIgnoreGlob, defaultCustomPrompt string, verbose bool) (*CLIContainer, error) {
	c := &CLIContainer{}
	var err error

	// Logger for CLI
	logger := NewCLILogger(verbose)
	c.Log = logger

	// Repositories and Infrastructure
	c.SettingsRepo, err = settingsfs.New(c.Log, embeddedIgnoreGlob, defaultCustomPrompt)
	if err != nil {
		return nil, err
	}
	c.FileReader = filereader.NewSecureFileReader(c.Log)
	c.GitRepo = git.New(c.Log)
	c.TreeBuilder = fsscanner.New(c.SettingsRepo, c.Log)
	c.ContextSplitter = textutils.NewContextSplitter(c.Log)
	c.CommandRunner = exec.NewCommandRunnerImpl(c.Log)

	// Application Services
	modelFetchers := createModelFetchers(ctx, c.Log, c.SettingsRepo)
	c.SettingsService, err = settings.NewService(c.Log, nil, c.SettingsRepo, modelFetchers)
	if err != nil {
		return nil, err
	}

	// AI Service needs to be created before context service
	providerRegistry := createProviderRegistry(c.Log, c.SettingsService)

	// Create rate limiter and metrics collector
	rateLimiter := appai.NewRateLimiter()
	metrics := appai.NewMetricsCollector()

	// Create intelligent service with dependencies
	intelligentService := appai.NewIntelligentService(c.SettingsService, c.Log, rateLimiter, metrics)

	// Create AI service with intelligent service
	c.AIService = appai.NewService(c.SettingsService, c.Log, providerRegistry, intelligentService)

	// Create OPA service
	c.opaService = policy.NewOPAService(c.Log)

	// Create file stat provider (using standard os implementation)
	fileStatProvider := filesystem.NewOSFileStatProvider()

	// Create unified ContextService
	c.ContextService, err = contextservice.NewService(
		c.FileReader,
		&SimpleTokenCounter{},
		nil, // No event bus for CLI
		c.Log,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create context service: %w", err)
	}

	// Create unified ProjectService
	c.ProjectService = projectservice.NewService(
		c.Log,
		nil, // No event bus for CLI
		c.TreeBuilder,
		c.GitRepo,
		c.ContextService,
	)
	_ = fileStatProvider // Used by other services
	c.ContextAnalysis = analysis.NewKeywordAnalyzer(c.Log)
	// Create symbol graph builders
	goSymbolGraphBuilder := symbolgraph.NewGoSymbolGraphBuilder(c.Log)
	symbolGraphBuilders := make(map[string]domain.SymbolGraphBuilder)
	symbolGraphBuilders["go"] = goSymbolGraphBuilder

	// Create import graph builders (currently no implementation, using nil map)
	importGraphBuilders := make(map[string]domain.ImportGraphBuilder)

	c.SymbolGraph = symbol.NewService(c.Log, symbolGraphBuilders, importGraphBuilders)
	testEngine := testengine.NewTestEngine(c.Log, goSymbolGraphBuilder)
	c.TestService = build.NewTestService(c.Log, testEngine)
	staticAnalyzerEngine := staticanalyzer.NewStaticAnalyzerEngine(c.Log)
	c.StaticAnalyzerService = analysis.NewStaticAnalyzerService(c.Log, staticAnalyzerEngine)

	// Create SBOM infrastructure components
	sbomGenerator := sbomlicensing.NewSyftGenerator(c.Log)
	vulnScanner := sbomlicensing.NewGrypeScanner(c.Log)
	licenseScanner := sbomlicensing.NewLicenseScanner(c.Log)

	// Create SBOM service with all required dependencies
	c.SBOMService = sbom.NewService(c.Log, sbomGenerator, vulnScanner, licenseScanner, fileStatProvider)

	c.RepairService = repair.NewService(c.Log, c.CommandRunner)

	// Taskflow components not used in CLI currently

	// Create Guardrail service with required dependencies
	c.GuardrailService = guardrails.NewService(c.Log, c.opaService, fileStatProvider)

	// Create UX Metrics infrastructure components
	uxRepo := uxreports.NewInMemoryUXReportRepository()

	// Create UX Metrics service
	c.UXMetricsService = ux.NewService(c.Log, uxRepo)

	// Create Apply service infrastructure components
	applyEngine := applyengine.NewApplyEngine(c.Log, &domain.ApplyEngineConfig{
		AutoFormat:     true,
		AutoFixImports: true,
		BackupFiles:    true,
		ValidateAfter:  true,
		Languages:      []string{"go", "typescript", "ts"},
	})

	// Create formatters and import fixers
	formattersMap := make(map[string]domain.Formatter)
	importFixers := make(map[string]domain.ImportFixer)

	// Create and register Go formatter and import fixer
	goFormatter := formatters.NewGoFormatter(c.Log)
	formattersMap["go"] = goFormatter
	importFixers["go"] = goFormatter

	// Create and register TypeScript formatter and import fixer
	tsFormatter := formatters.NewTypeScriptFormatter(c.Log)
	formattersMap["typescript"] = tsFormatter
	formattersMap["ts"] = tsFormatter
	importFixers["typescript"] = tsFormatter
	importFixers["ts"] = tsFormatter

	// Register formatters and import fixers with the engine
	for lang, formatter := range formattersMap {
		applyEngine.RegisterFormatter(lang, formatter)
	}

	for lang, fixer := range importFixers {
		applyEngine.RegisterImportFixer(lang, fixer)
	}

	// Create Apply service with all required dependencies
	applyConfig := &domain.ApplyEngineConfig{
		AutoFormat:     true,
		AutoFixImports: true,
		BackupFiles:    true,
		ValidateAfter:  true,
		Languages:      []string{"go", "typescript", "ts"},
	}
	c.ApplyService = diff.NewApplyService(c.Log, applyConfig, applyEngine, formattersMap, importFixers)

	// Create Diff service
	diffEngine := diffengine.NewDiffEngine(c.Log)
	c.DiffService = diff.NewService(c.Log, diffEngine)

	buildPipeline := buildpipeline.NewBuildPipeline(c.Log)
	c.BuildService = build.NewService(c.Log, buildPipeline)

	// Create formatter service
	formatterService := export.NewFormatterService(c.Log, c.CommandRunner)

	// Create verification pipeline service
	c.VerificationService = verification.NewService(
		c.Log,
		c.BuildService,
		c.TestService,
		c.StaticAnalyzerService,
		formatterService,
		&OSFileSystemWriter{},
		nil, // Task Protocol Service not needed for CLI
	)

	// new: wire PDF and ZIP implementations
	pdfGen := pdfgen.NewGofpdfGenerator(c.Log)
	arch := archiverinfra.NewZipArchiver(c.Log)
	contextFormatter := contextbuilder.NewContextFormatter()

	// Create Export service with all required dependencies
	c.ExportService = export.NewService(
		c.Log,
		c.ContextSplitter,
		contextFormatter,
		pdfGen,
		arch,
		&OSTempFileProvider{}, // Temp file provider
		&FilePathProvider{},   // Path provider
		&OSFileSystemWriter{}, // File system writer
		fileStatProvider,      // File stat provider
	)

	return c, nil
}

// CLILogger реализует простой логгер для CLI
type CLILogger struct {
	verbose bool
}

// NewCLILogger создает новый CLI логгер
func NewCLILogger(verbose bool) *CLILogger {
	return &CLILogger{
		verbose: verbose,
	}
}

// Info логирует информационное сообщение
func (l *CLILogger) Info(message string) {
	if l.verbose {
		fmt.Printf("[INFO] %s\n", message)
	}
}

// Warning логирует предупреждение
func (l *CLILogger) Warning(message string) {
	fmt.Printf("[WARN] %s\n", message)
}

// Error логирует ошибку
func (l *CLILogger) Error(message string) {
	fmt.Printf("[ERROR] %s\n", message)
}

// Debug логирует отладочное сообщение
func (l *CLILogger) Debug(message string) {
	if l.verbose {
		fmt.Printf("[DEBUG] %s\n", message)
	}
}

// Fatal логирует фатальную ошибку и завершает программу
func (l *CLILogger) Fatal(message string) {
	fmt.Printf("[FATAL] %s\n", message)
	os.Exit(1)
}

func createModelFetchers(ctx context.Context, log domain.Logger, repo domain.SettingsRepository) domain.ModelFetcherRegistry {
	fetchers := make(domain.ModelFetcherRegistry)

	// Gemini
	fetchers["gemini"] = func(apiKey string) ([]string, error) {
		p, err := ai.NewGemini(apiKey, "", log)
		if err != nil {
			log.Warning("Failed to create Gemini client for model listing: " + err.Error())
			return nil, err
		}
		return p.(*ai.GeminiProviderImpl).ListModels(ctx)
	}

	// OpenAI
	fetchers["openai"] = func(apiKey string) ([]string, error) {
		p, err := ai.NewOpenAI(apiKey, "", log)
		if err != nil {
			log.Warning("Failed to create OpenAI client for model listing: " + err.Error())
			return nil, err
		}
		return p.(*ai.OpenAIProviderImpl).ListModels(ctx)
	}

	// OpenRouter
	fetchers["openrouter"] = func(apiKey string) ([]string, error) {
		p, err := ai.NewOpenAI(apiKey, openRouterHost, log)
		if err != nil {
			log.Warning("Failed to create OpenRouter client for model listing: " + err.Error())
			return nil, err
		}
		return p.(*ai.OpenAIProviderImpl).ListModels(ctx)
	}

	// LocalAI
	fetchers["localai"] = func(apiKey string) ([]string, error) {
		localAIHost := repo.GetLocalAIHost()
		p, err := ai.NewLocalAI(apiKey, localAIHost, log)
		if err != nil {
			log.Warning("Failed to create LocalAI client for model listing: " + err.Error())
			return nil, err
		}
		return p.(*ai.LocalAIProviderImpl).ListModels(ctx)
	}

	return fetchers
}

func createProviderRegistry(log domain.Logger, settingsService *settings.Service) map[string]domain.AIProviderFactory {
	resolveHost := func(providerType string) (string, error) {
		switch providerType {
		case "openrouter":
			return openRouterHost, nil
		case "localai":
			dto, err := settingsService.GetSettingsDTO()
			if err != nil {
				return "", err
			}
			return dto.LocalAIHost, nil
		default:
			return "", nil
		}
	}

	return ai.NewAIProviderFactoryRegistry(log, openRouterHost, resolveHost)
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

func (t *OSTempFileProvider) MkdirTemp(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}

// SimpleTokenCounter provides basic token estimation
type SimpleTokenCounter struct{}

func (s *SimpleTokenCounter) CountTokens(text string) int {
	// Simple approximation: 1 token ≈ 4 characters
	return len(text) / 4
}
