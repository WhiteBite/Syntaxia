package main

import (
	"context"
	"encoding/json"
	"fmt"

	appai "syntaxia/application/ai"
	"syntaxia/application/analysis"
	"syntaxia/application/diff"
	"syntaxia/application/export"
	"syntaxia/application/protocol"
	"syntaxia/application/sbom"
	"syntaxia/application/settings"
	"syntaxia/application/symbol"
	"syntaxia/application/taskflow"
	"syntaxia/cmd/app"
	"syntaxia/domain"
	"syntaxia/handlers"
	"syntaxia/infrastructure/wailsbridge"
	contextservice "syntaxia/internal/context"
	projectservice "syntaxia/internal/project"
)

type contextAnalysisService interface {
	domain.ContextAnalyzer
	AnalyzeTaskAndCollectContext(ctx context.Context, task string, allFiles []*domain.FileNode, rootDir string) (*domain.ContextAnalysisResult, error)
}

// App is the main application struct that exposes methods to the frontend via Wails.
// API methods are organized in separate files in the api/ directory:
// - api/project_api.go - file/directory operations
// - api/context_api.go - context building and export
// - api/ai_api.go - AI generation and semantic search
// - api/git_api.go - Git operations (local, GitHub, GitLab)
// - api/analysis_api.go - code analysis, testing, SBOM
// - api/settings_api.go - application settings
// - api/taskflow_api.go - task execution and protocols
// - api/window_api.go - window state management
// - api/ux_api.go - UX metrics and reports
// - api/apply_api.go - code edits and diffs
type App struct {
	ctx       context.Context
	log       domain.Logger
	bridge    *wailsbridge.Bridge
	container *app.AppContainer

	// Handlers (new architecture) - primary delegation targets
	projectHandler  *handlers.ProjectHandler
	contextHandler  *handlers.ContextHandler
	aiHandler       *handlers.AIHandler
	analysisHandler *handlers.AnalysisHandler
	settingsHandler *handlers.SettingsHandler
	taskflowHandler *handlers.TaskflowHandler

	// Services (kept for methods not yet migrated to handlers)
	projectService        *projectservice.Service
	contextService        *contextservice.Service
	aiService             *appai.Service
	settingsService       *settings.Service
	contextAnalysis       contextAnalysisService
	symbolGraph           *symbol.Service
	testService           domain.ITestService
	staticAnalyzerService domain.IStaticAnalyzerService
	sbomService           *sbom.Service
	repairService         domain.RepairService
	guardrailService      domain.GuardrailService
	taskflowService       domain.TaskflowService
	uxMetricsService      domain.UXMetricsService
	applyService          *diff.ApplyService
	diffService           *diff.Service
	buildService          domain.IBuildService
	fileWatcher           domain.FileSystemWatcher
	gitRepo               domain.GitRepository
	exportService         *export.Service
	fileReader            domain.FileContentReader
	reportService         *export.ReportService

	// Task Protocol Services
	taskProtocolService         domain.TaskProtocolService
	taskProtocolConfigService   *protocol.ConfigService
	taskflowProtocolIntegration *taskflow.ProtocolIntegration

	// Qwen Services
	qwenHandler *handlers.QwenHandler

	// Analysis Container (for smart analysis tools)
	analysisContainer *analysis.Container

	// Startup path from command line arguments
	startupPath string
}

func (a *App) startup(ctx context.Context, container *app.AppContainer) {
	a.ctx = ctx
	a.log = container.Log
	a.bridge = container.Bridge
	a.container = container

	// Auto-update shell integration if exe path changed
	if container.ShellIntegration != nil {
		if updated, err := container.ShellIntegration.UpdateIfNeeded(); err != nil {
			a.log.Warning("Failed to update shell integration: " + err.Error())
		} else if updated {
			a.log.Info("Shell integration path updated")
		}
	}

	// Initialize handlers from container
	a.projectHandler = container.ProjectHandler
	a.contextHandler = container.ContextHandler
	a.aiHandler = container.AIHandler
	a.analysisHandler = container.AnalysisHandler
	a.settingsHandler = container.SettingsHandler
	a.taskflowHandler = container.TaskflowHandler

	// Services (for methods not yet migrated)
	a.projectService = container.ProjectService
	a.contextService = container.ContextService
	a.aiService = container.AIService
	a.settingsService = container.SettingsService
	if analyzer, ok := container.ContextAnalysis.(contextAnalysisService); ok {
		a.contextAnalysis = analyzer
	} else {
		a.contextAnalysis = nil
		if container.ContextAnalysis != nil {
			a.log.Warning("context analysis service does not implement required AnalyzeTaskAndCollectContext method")
		}
	}
	a.symbolGraph = container.SymbolGraph
	a.testService = container.TestService
	a.staticAnalyzerService = container.StaticAnalyzerService
	a.sbomService = container.SBOMService
	a.repairService = container.RepairService
	a.guardrailService = container.GuardrailService
	a.taskflowService = container.TaskflowService
	a.uxMetricsService = container.UXMetricsService
	a.applyService = container.ApplyService
	a.diffService = container.DiffService
	a.buildService = container.BuildService
	a.fileWatcher = container.Watcher
	a.gitRepo = container.GitRepo
	a.exportService = container.ExportService
	a.fileReader = container.FileReader
	a.reportService = container.ReportService

	// Task Protocol Services
	a.taskProtocolService = container.TaskProtocolService
	a.taskProtocolConfigService = container.TaskProtocolConfigService
	a.taskflowProtocolIntegration = container.TaskflowProtocolIntegration

	// Qwen Handler
	a.qwenHandler = container.QwenHandler

	// Analysis Container
	a.analysisContainer = container.AnalysisContainer
}

func (a *App) domReady(ctx context.Context) {
	a.ctx = ctx
	a.bridge.SetWailsContext(ctx)

	// Restore window state after DOM is ready
	if err := a.LoadWindowState(); err != nil {
		a.log.Warning("Failed to load window state: " + err.Error())
	}
}

func (a *App) shutdown(ctx context.Context) {
	a.ctx = ctx

	// Save window state before shutdown
	if err := a.SaveWindowState(); err != nil {
		a.log.Warning("Failed to save window state: " + err.Error())
	}

	// Shutdown container services first
	if a.container != nil {
		if err := a.container.Shutdown(ctx); err != nil {
			a.log.Warning("Container shutdown error: " + err.Error())
		}
	}

	// Stop file watcher
	a.fileWatcher.Stop()
}

// transformError transforms any error into a structured domain error
func (a *App) transformError(err error) error {
	if err == nil {
		return nil
	}

	if domainErr, ok := err.(*domain.DomainError); ok {
		return a.transformDomainError(domainErr)
	}

	return a.transformDomainError(domain.NewInternalError("An unexpected error occurred", err))
}

// transformDomainError transforms domain errors into frontend-friendly format
func (a *App) transformDomainError(domainErr *domain.DomainError) error {
	frontendError := map[string]interface{}{
		"code":        domainErr.Code,
		"message":     domainErr.Message,
		"recoverable": domainErr.Recoverable,
		"context":     domainErr.Context,
	}

	if domainErr.Cause != nil {
		frontendError["cause"] = domainErr.Cause.Error()
	}

	errorJson, _ := json.Marshal(frontendError)
	return fmt.Errorf("domain_error:%s", string(errorJson))
}
