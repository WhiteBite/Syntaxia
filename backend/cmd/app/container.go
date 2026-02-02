package app

import (
	"context"
	"fmt"
	"syntaxia/application"
	appai "syntaxia/application/ai"
	"syntaxia/application/analysis"
	"syntaxia/application/diff"
	"syntaxia/application/export"
	"syntaxia/application/protocol"
	"syntaxia/application/rag"
	"syntaxia/application/router"
	"syntaxia/application/sbom"
	"syntaxia/application/settings"
	"syntaxia/application/symbol"
	"syntaxia/application/taskflow"
	"syntaxia/application/verification"
	"syntaxia/domain"
	domainanalysis "syntaxia/domain/analysis"
	"syntaxia/handlers"
	"syntaxia/infrastructure/shellintegration"
	"syntaxia/infrastructure/wailsbridge"
	"syntaxia/internal/initmanager"
	"sync"

	contextservice "syntaxia/internal/context"
	projectservice "syntaxia/internal/project"
)

// Use domain constants for default hosts
const openRouterHost = domain.OpenRouterDefaultHost

// AppContainer holds all the services and repositories for the application.
type AppContainer struct {
	Log                   domain.Logger
	Bus                   domain.EventBus
	SettingsRepo          domain.SettingsRepository
	FileReader            domain.FileContentReader
	GitRepo               domain.GitRepository
	TreeBuilder           domain.TreeBuilder
	ContextSplitter       domain.ContextSplitter
	Watcher               domain.FileSystemWatcher
	CommandRunner         domain.CommandRunner
	SettingsService       *settings.Service
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

	// Sandbox for AI file changes
	SandboxFS domain.SandboxFS

	// Smart Context Collector
	SmartContextCollector domain.SmartContextCollector

	// Unified internal services (new architecture)
	ContextService *contextservice.Service
	ProjectService *projectservice.Service

	// Context interfaces (implemented by ContextService)
	ContextBuilder    domain.ContextBuilder
	ContextAnalyzer   domain.ContextAnalyzer
	ContextRepository domain.ContextRepository

	ReportService    *export.ReportService
	RouterLLMService *router.LLMService
	Bridge           *wailsbridge.Bridge

	// Task Protocol Services
	TaskProtocolService         domain.TaskProtocolService
	ErrorAnalyzer               domain.ErrorAnalyzer
	CorrectionEngine            domain.CorrectionEngine
	TaskProtocolConfigService   *protocol.ConfigService
	TaskflowProtocolIntegration *taskflow.ProtocolIntegration
	VerificationPipelineService *verification.Service

	// Qwen Task Services
	SmartContextService *rag.SmartContextService
	QwenTaskService     *appai.QwenTaskService

	// Semantic Search Services
	SymbolIndex       domainanalysis.SymbolIndex
	EmbeddingProvider domain.EmbeddingProvider
	VectorStore       domain.VectorStore
	SemanticSearch    domain.SemanticSearchService
	RAGService        domain.RAGService
	SemanticHandler   *handlers.SemanticHandler

	// Handlers (new architecture)
	ProjectHandler  *handlers.ProjectHandler
	ContextHandler  *handlers.ContextHandler
	QwenHandler     *handlers.QwenHandler
	AIHandler       *handlers.AIHandler
	AnalysisHandler *handlers.AnalysisHandler
	SettingsHandler *handlers.SettingsHandler
	TaskflowHandler *handlers.TaskflowHandler

	// Analysis tools (shared across handlers)
	AnalysisContainer *analysis.Container
	ToolExecutor      *application.ToolExecutorImpl

	// Dependency Analysis
	DependencyAnalyzer domain.DependencyAnalyzer

	// Lazy initialization support
	lazyInitOnce              sync.Once
	testServiceOnce           sync.Once
	staticAnalyzerServiceOnce sync.Once
	sbomServiceOnce           sync.Once
	symbolGraphOnce           sync.Once

	// Lazy service manager for coordinated lifecycle management
	lazyManager *initmanager.LazyServiceManager

	// Cleanup goroutine control
	cleanupStopCh chan struct{}

	// Shell Integration (OS context menu)
	ShellIntegration *shellintegration.Service

	// Infrastructure components (internal use only)
	infrastructureComponents *infrastructureComponents
}

// NewContainer creates and wires up all the application dependencies.
func NewContainer(ctx context.Context, embeddedIgnoreGlob, defaultCustomPrompt string) (*AppContainer, error) {
	c := &AppContainer{}

	// Initialize infrastructure layer
	if err := c.initializeInfrastructure(ctx, embeddedIgnoreGlob, defaultCustomPrompt); err != nil {
		return nil, fmt.Errorf("failed to initialize infrastructure: %w", err)
	}

	// Initialize application services
	if err := c.initializeServices(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize handlers (new architecture)
	if err := c.initializeHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	return c, nil
}

// GetLazyServiceStats returns statistics about lazy-loaded services
func (c *AppContainer) GetLazyServiceStats() map[string]interface{} {
	if c.lazyManager == nil {
		return map[string]interface{}{
			"enabled": false,
			"message": "Lazy service manager not initialized",
		}
	}

	stats := c.lazyManager.GetInitializationStats()
	stats["enabled"] = true
	return stats
}

// Shutdown gracefully shuts down all services in the container
func (c *AppContainer) Shutdown(ctx context.Context) error {
	c.Log.Info("Starting container shutdown...")

	var shutdownErrors []error

	// Stop cleanup goroutine first
	if c.cleanupStopCh != nil {
		close(c.cleanupStopCh)
	}

	// Shutdown handlers that support it
	if c.AIHandler != nil {
		if err := c.AIHandler.Shutdown(ctx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("AIHandler shutdown: %w", err))
		}
	}

	// Shutdown services
	if c.AIService != nil {
		if err := c.AIService.Shutdown(ctx); err != nil {
			shutdownErrors = append(shutdownErrors, fmt.Errorf("AIService shutdown: %w", err))
		}
	}

	// Stop file watcher
	if c.Watcher != nil {
		c.Watcher.Stop()
	}

	// Close cached symbol index if it supports closing
	if c.SymbolIndex != nil {
		if closer, ok := c.SymbolIndex.(interface{ Close() error }); ok {
			if err := closer.Close(); err != nil {
				shutdownErrors = append(shutdownErrors, fmt.Errorf("SymbolIndex close: %w", err))
			}
		}
	}

	// Shutdown lazy service manager
	if c.lazyManager != nil {
		// Force unload all services on shutdown
		c.lazyManager.UnloadUnusedServices(0)
	}

	if len(shutdownErrors) > 0 {
		c.Log.Warning(fmt.Sprintf("Container shutdown completed with %d errors", len(shutdownErrors)))
		return fmt.Errorf("shutdown errors: %v", shutdownErrors)
	}

	c.Log.Info("Container shutdown complete")
	return nil
}
