package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	appai "syntaxia/application/ai"
	"syntaxia/application/analysis"
	appcontext "syntaxia/application/context"
	"syntaxia/application/build"
	"syntaxia/application/diff"
	"syntaxia/application/export"
	"syntaxia/application/guardrails"
	"syntaxia/application/protocol"
	"syntaxia/application/rag"
	"syntaxia/application/repair"
	"syntaxia/application/router"
	"syntaxia/application/sbom"
	"syntaxia/application/settings"
	"syntaxia/application/symbol"
	"syntaxia/application/taskflow"
	"syntaxia/application/ux"
	"syntaxia/application/verification"
	"syntaxia/domain"
	"syntaxia/infrastructure/ai"
	"syntaxia/infrastructure/analyzers"
	"syntaxia/infrastructure/embeddings"
	"syntaxia/infrastructure/policy"
	"syntaxia/infrastructure/sandbox"
	"syntaxia/infrastructure/symbolgraph"
	"syntaxia/internal/initmanager"
	"sync"
	"time"

	contextservice "syntaxia/internal/context"
	projectservice "syntaxia/internal/project"
)

// initializeServices initializes all application layer services
func (c *AppContainer) initializeServices(ctx context.Context) error {
	// Application Services
	modelFetchers := createModelFetchers(ctx, c.Log, c.SettingsRepo)
	var err error
	c.SettingsService, err = settings.NewService(c.Log, c.Bus, c.SettingsRepo, modelFetchers)
	if err != nil {
		return err
	}
	// Connect watcher to settings changes
	c.SettingsService.OnIgnoreRulesChanged(c.Watcher.RefreshAndRescan)

	// AI Service needs to be created before context service
	providerRegistry := createProviderRegistry(c.Log, c.SettingsService)

	// Create rate limiter and metrics collector
	rateLimiter := appai.NewRateLimiter()
	metrics := appai.NewMetricsCollector()

	// Create intelligent service with dependencies
	intelligentService := appai.NewIntelligentService(c.SettingsService, c.Log, rateLimiter, metrics)

	// Create AI service with intelligent service
	c.AIService = appai.NewService(c.SettingsService, c.Log, providerRegistry, intelligentService)

	// Set provider getter in IntelligentAIService (uses interface to break circular dependency)
	intelligentService.SetProviderGetter(c.AIService)

	// Connect SettingsService to AIService for cache invalidation on settings change
	c.SettingsService.SetAICacheInvalidator(c.AIService)

	// Create OPA service (will be used by ContextService internally)
	_ = policy.NewOPAService(c.Log)

	// Create file system writer (using standard os implementation)
	_ = &OSFileSystemWriter{} // Will be used by ContextService internally

	// Get context directory
	homeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return fmt.Errorf("failed to determine user home directory: %w", homeErr)
	}
	contextDir := filepath.Join(homeDir, ".syntaxia", "contexts")
	if mkErr := os.MkdirAll(contextDir, 0o755); mkErr != nil {
		return fmt.Errorf("failed to create context directory: %w", mkErr)
	}

	// Create unified ContextService (replaces ContextBuilder, ContextGenerator, ContextRepository)
	c.ContextService, err = contextservice.NewService(
		c.FileReader,
		&SimpleTokenCounter{},
		c.Bus,
		c.Log,
	)
	if err != nil {
		return fmt.Errorf("failed to create context service: %w", err)
	}

	// ContextService implements ContextRepository interface
	c.ContextRepository = c.ContextService

	// Create ContextBuilder adapter
	c.ContextBuilder = contextservice.NewContextBuilderAdapter(c.ContextService)

	// Create analyzer responsible for task-driven context suggestions
	contextAnalyzer := analysis.NewContextAnalyzer(c.Log, c.AIService)
	c.ContextAnalyzer = contextAnalyzer

	// Create unified ProjectService
	c.ProjectService = projectservice.NewService(
		c.Log,
		c.Bus,
		c.TreeBuilder,
		c.GitRepo,
		c.ContextService, // Pass context service interface
	)

	// Initialize remaining services
	c.ContextAnalysis = contextAnalyzer

	// Smart Context Collector
	c.SmartContextCollector = appcontext.NewSmartContextCollector(c.Log, c.FileReader, c.TreeBuilder)

	// Symbol Graph Service
	c.SymbolGraph = symbol.NewService(c.Log, c.infrastructureComponents.symbolGraphBuilders, c.infrastructureComponents.importGraphBuilders)

	// Create CallStack Analyzer and Smart Context Service for Qwen integration
	callStackAnalyzer := symbolgraph.NewCallStackAnalyzerAdapter(c.Log)
	c.SmartContextService = rag.NewSmartContextService(
		c.Log,
		c.FileReader,
		c.SymbolGraph,
		callStackAnalyzer,
	)

	// Static Analyzer Service
	c.StaticAnalyzerService = analysis.NewStaticAnalyzerService(c.Log, c.infrastructureComponents.staticAnalyzerEngine)

	// SBOM Service
	sbomFileStatProvider := &OSFileStatProvider{}
	c.SBOMService = sbom.NewService(c.Log, c.infrastructureComponents.sbomGenerator, c.infrastructureComponents.vulnScanner, c.infrastructureComponents.licenseScanner, sbomFileStatProvider)

	// Repair Service
	c.RepairService = repair.NewService(c.Log, c.CommandRunner)

	// Build Service
	c.BuildService = build.NewService(c.Log, c.infrastructureComponents.buildPipeline)

	// Test Service (already initialized in infrastructure via testServiceOnce)

	// Create RouterPlannerService
	planner := router.NewPlannerService(c.Log, c.BuildService, c.TestService, c.StaticAnalyzerService, c.RepairService)

	// Guardrail Service
	guardrailFileStatProvider := &OSFileStatProvider{}
	c.GuardrailService = guardrails.NewService(c.Log, c.infrastructureComponents.guardrailOPAService, guardrailFileStatProvider)

	// Taskflow Service
	c.TaskflowService = taskflow.NewService(c.Log, planner, c.RouterLLMService, c.GuardrailService, c.infrastructureComponents.taskflowRepo, c.GitRepo)

	// ⚠️ CRITICAL: Update GuardrailService with TaskTypeProvider to resolve circular dependency
	if taskTypeProvider, ok := c.TaskflowService.(domain.TaskTypeProvider); ok {
		c.GuardrailService.SetTaskTypeProvider(taskTypeProvider)
	} else {
		c.Log.Warning("TaskflowService does not implement TaskTypeProvider; guardrails may be limited")
	}

	// UX Metrics Service
	c.UXMetricsService = ux.NewService(c.Log, c.infrastructureComponents.uxReportRepo)

	// Apply Service
	c.ApplyService = diff.NewApplyService(c.Log, c.infrastructureComponents.applyConfig, c.infrastructureComponents.applyEngine, c.infrastructureComponents.formatterMap, c.infrastructureComponents.importFixerMap)

	// Diff Service
	c.DiffService = diff.NewService(c.Log, c.infrastructureComponents.diffEngine)

	// Export Service
	tempFileProvider := &OSTempFileProvider{}
	exportFileStatProvider := &OSFileStatProvider{}
	exportPathProvider := &FilePathProvider{}
	exportFileSystemWriter := &OSFileSystemWriter{}
	c.ExportService = export.NewService(c.Log, c.ContextSplitter, c.infrastructureComponents.contextFormatter, c.infrastructureComponents.pdfGen, c.infrastructureComponents.arch, tempFileProvider, exportPathProvider, exportFileSystemWriter, exportFileStatProvider)

	// Report Service
	c.ReportService = export.NewReportService(c.Log, c.infrastructureComponents.reportRepo)

	// Router LLM Service
	c.RouterLLMService = router.NewLLMServiceWithClient(c.infrastructureComponents.routerLLMConfig, c.Log, c.infrastructureComponents.llmClient, c.infrastructureComponents.fileReader)

	// Initialize Task Protocol Services
	if err := initializeTaskProtocolServices(c); err != nil {
		return fmt.Errorf("failed to initialize task protocol services: %w", err)
	}

	// Initialize lazy service manager for memory optimization
	c.lazyManager = initmanager.NewLazyServiceManager()

	// Start periodic cleanup of unused services (runs every 5 minutes)
	c.cleanupStopCh = make(chan struct{})
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-c.cleanupStopCh:
				return
			case <-ticker.C:
				// Unload services idle for more than 10 minutes
				unloaded := c.lazyManager.UnloadUnusedServices(10 * time.Minute)
				if unloaded > 0 {
					c.Log.Info(fmt.Sprintf("Unloaded %d idle services to free memory", unloaded))
				}
			}
		}
	}()

	// Initialize Semantic Search Services
	if err := c.initializeSemanticSearch(); err != nil {
		c.Log.Warning(fmt.Sprintf("Semantic search initialization failed (non-critical): %v", err))
		// Non-critical - continue without semantic search
	}

	return nil
}

// initializeSemanticSearch initializes semantic search services
func (c *AppContainer) initializeSemanticSearch() error {
	// Get data directory for embeddings storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	dataDir := filepath.Join(homeDir, ".syntaxia", "embeddings")

	// Create analyzer registry for symbol extraction
	analyzerRegistry := analyzers.NewAnalyzerRegistry()

	// Create symbol index with SQLite caching for incremental indexing
	symbolCacheDir := filepath.Join(dataDir, "symbol_cache")
	cachedSymbolIndex, err := analyzers.NewCachedSymbolIndex(analyzerRegistry, symbolCacheDir)
	if err != nil {
		c.Log.Warning(fmt.Sprintf("Failed to create cached symbol index, falling back to in-memory: %v", err))
		c.SymbolIndex = analyzers.NewSymbolIndex(analyzerRegistry)
	} else {
		c.SymbolIndex = cachedSymbolIndex
	}

	// Create vector store (SQLite-based)
	vectorStore, err := embeddings.NewSQLiteVectorStore(dataDir, c.Log)
	if err != nil {
		return fmt.Errorf("failed to create vector store: %w", err)
	}
	c.VectorStore = vectorStore

	// Create embedding provider (OpenAI by default)
	// Get API key from settings
	settings, err := c.SettingsService.GetSettingsDTO()
	if err != nil {
		c.Log.Warning("Failed to get settings for embedding provider: " + err.Error())
	}

	apiKey := ""
	if settings.OpenAIAPIKey != "" {
		apiKey = settings.OpenAIAPIKey
	}

	if apiKey != "" {
		embeddingProvider, err := embeddings.NewOpenAIEmbeddingProvider(
			apiKey,
			domain.EmbeddingModelOpenAI3S, // Use small model by default
			c.Log,
		)
		if err != nil {
			c.Log.Warning("Failed to create embedding provider: " + err.Error())
		} else {
			c.EmbeddingProvider = embeddingProvider
		}
	}

	// Create semantic search service (only if embedding provider is available)
	if c.EmbeddingProvider != nil {
		// Create code chunker for semantic indexing
		chunker := &codeChunkerAdapter{impl: embeddings.NewCodeChunker(embeddings.DefaultChunkerConfig())}

		c.SemanticSearch = rag.NewSemanticSearchService(
			c.EmbeddingProvider,
			c.VectorStore,
			c.SymbolIndex,
			c.Log,
			chunker,
		)

		// Create RAG service
		c.RAGService = rag.NewService(
			c.SemanticSearch,
			c.EmbeddingProvider,
			c.Log,
		)

		c.Log.Info("Semantic search services initialized successfully")
	} else {
		c.Log.Warning("Semantic search disabled: no embedding provider configured (set OpenAI API key)")
	}

	return nil
}

// initializeTaskProtocolServices initializes the Task Protocol services
func initializeTaskProtocolServices(c *AppContainer) error {
	// Initialize Error Analyzer
	c.ErrorAnalyzer = repair.NewErrorAnalyzer(c.Log)

	// Initialize Correction Engine with file system provider
	fileSystemProvider := &OSFileSystemProvider{}
	c.CorrectionEngine = repair.NewCorrectionEngine(c.Log, fileSystemProvider)

	// Initialize SandboxFS for AI file changes
	// Note: projectRoot will be set when project is opened via SetProjectRoot
	c.SandboxFS = sandbox.NewSandboxFS("", fileSystemProvider, c.Log)

	// Initialize Task Protocol Config Service
	c.TaskProtocolConfigService = protocol.NewConfigService(c.Log, fileSystemProvider)

	// Initialize Task Protocol Service
	c.TaskProtocolService = protocol.NewService(
		c.Log,
		nil, // Will be set after VerificationPipelineService is created
		c.StaticAnalyzerService,
		c.TestService,
		c.BuildService,
		c.GuardrailService,
		c.AIService.GetIntelligentService(),
		c.ErrorAnalyzer,
		c.CorrectionEngine,
	)

	// Create VerificationPipelineService with Task Protocol integration
	formatterService := export.NewFormatterService(c.Log, &CommandRunnerImpl{})
	c.VerificationPipelineService = verification.NewService(
		c.Log,
		c.BuildService,
		c.TestService,
		c.StaticAnalyzerService,
		formatterService,
		&OSFileSystemWriter{},
		c.TaskProtocolService,
	)

	// Initialize Taskflow Protocol Integration
	c.TaskflowProtocolIntegration = taskflow.NewProtocolIntegration(
		c.Log,
		c.TaskflowService,
		c.TaskProtocolService,
		c.TaskProtocolConfigService,
		c.AIService.GetIntelligentService(),
	)

	return nil
}

// createModelFetchers creates model fetchers for all AI providers
func createModelFetchers(ctx context.Context, log domain.Logger, repo domain.SettingsRepository) domain.ModelFetcherRegistry {
	registry := ai.GetProviderRegistry(openRouterHost)
	fetchers := make(domain.ModelFetcherRegistry)

	for providerType, config := range registry {
		// Capture variables for closure
		providerType := providerType
		config := config

		// Create a cached fetcher
		cachedFetcher := &cachedModelFetcher{
			fetcher: config.ModelFetcher,
			log:     log,
			repo:    repo,
			cache:   make(map[string][]string),
		}

		fetchers[providerType] = func(apiKey string) ([]string, error) {
			// For the model fetchers, we need to provide the host and logger
			// We'll use the container's logger and get the host based on provider type
			host := ""
			if providerType == "openrouter" {
				host = openRouterHost
			} else if providerType == "localai" {
				host = repo.GetLocalAIHost()
			} else if providerType == "qwen" {
				host = repo.GetQwenHost()
			}

			models, err := cachedFetcher.FetchModels(ctx, apiKey, host, log)
			if err != nil {
				log.Warning("Failed to create " + providerType + " client for model listing: " + err.Error())
				return nil, err
			}
			return models, nil
		}
	}

	return fetchers
}

// cachedModelFetcher adds caching to model fetchers
type cachedModelFetcher struct {
	fetcher func(context.Context, string, string, domain.Logger) ([]string, error)
	log     domain.Logger
	repo    domain.SettingsRepository
	cache   map[string][]string
	mu      sync.RWMutex
}

func (c *cachedModelFetcher) FetchModels(ctx context.Context, apiKey, host string, log domain.Logger) ([]string, error) {
	// Create a cache key based on API key and host
	cacheKey := apiKey + "|" + host

	// Check if we have cached models
	c.mu.RLock()
	if models, exists := c.cache[cacheKey]; exists {
		c.mu.RUnlock()
		log.Debug("Using cached models for provider")
		return models, nil
	}
	c.mu.RUnlock()

	// Fetch models and cache them
	models, err := c.fetcher(ctx, apiKey, host, log)
	if err != nil {
		return nil, err
	}

	// Store in cache
	c.mu.Lock()
	c.cache[cacheKey] = models
	c.mu.Unlock()

	return models, nil
}

// createProviderRegistry creates AI provider factory registry
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
		case "qwen":
			dto, err := settingsService.GetSettingsDTO()
			if err != nil {
				return "", err
			}
			if dto.QwenHost != "" {
				return dto.QwenHost, nil
			}
			return domain.QwenDefaultHost, nil
		default:
			return "", nil
		}
	}

	registry := ai.NewAIProviderFactoryRegistry(log, openRouterHost, resolveHost)

	// Override qwen-cli factory to pass settings
	registry["qwen-cli"] = func(_, _ string) (domain.AIProvider, error) {
		dto, err := settingsService.GetSettingsDTO()
		if err != nil {
			// Use defaults if settings unavailable
			return ai.NewQwenCLI(log)
		}
		return ai.NewQwenCLIWithSettings(log, dto.QwenCLISettings)
	}

	return registry
}
