package application

import (
	"fmt"

	appanalysis "syntaxia/application/analysis"
	"syntaxia/application/project"
	"syntaxia/application/tools"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// ToolExecutorImpl implements the ToolExecutor interface using HandlerRegistry
type ToolExecutorImpl struct {
	logger                  domain.Logger
	fileReader              domain.FileContentReader
	registry                analysis.AnalyzerRegistry
	symbolIndex             analysis.SymbolIndex
	callGraph               domain.CallGraphBuilder
	gitContext              domain.GitContextBuilder
	contextMemory           domain.ContextMemory
	referenceFinder         domain.ReferenceFinder
	projectStructure        domain.ProjectStructureDetector
	hasSemanticSearch       bool
	semanticSearcherService tools.SemanticSearcher
	handlerRegistry         *tools.HandlerRegistry
}

// NewToolExecutor creates a new ToolExecutor with all handlers registered
func NewToolExecutor(
	logger domain.Logger,
	fileReader domain.FileContentReader,
	registry analysis.AnalyzerRegistry,
	symbolIndex analysis.SymbolIndex,
	callGraph domain.CallGraphBuilder,
	referenceFinder domain.ReferenceFinder,
) *ToolExecutorImpl {
	te := &ToolExecutorImpl{
		logger:          logger,
		fileReader:      fileReader,
		registry:        registry,
		symbolIndex:     symbolIndex,
		callGraph:       callGraph,
		referenceFinder: referenceFinder,
		handlerRegistry: tools.NewHandlerRegistry(logger),
	}

	te.registerHandlers()
	return te
}

// registerHandlers registers all tool handlers
func (te *ToolExecutorImpl) registerHandlers() {
	// File tools
	te.handlerRegistry.Register(tools.NewFileToolsHandler(te.logger, te.fileReader))

	// Symbol tools
	te.handlerRegistry.Register(tools.NewSymbolToolsHandler(te.registry, te.symbolIndex, te.logger, te.referenceFinder))

	// Call graph tools
	te.handlerRegistry.Register(tools.NewCallGraphToolsHandler(te.logger, te.callGraph))

	// Git tools
	te.handlerRegistry.Register(tools.NewGitToolsHandler(te.logger, te.gitContext))

	// Memory tools
	te.handlerRegistry.Register(tools.NewMemoryToolsHandler(te.logger, te.contextMemory))

	// Preferences tools
	te.handlerRegistry.Register(tools.NewPreferencesToolsHandler(te.logger, te.contextMemory))

	// Project structure tools
	if te.projectStructure != nil {
		projectStructureService := project.NewStructureService(te.logger, te.projectStructure)
		te.handlerRegistry.Register(tools.NewProjectStructureToolsHandler(te.logger, projectStructureService))
	}
}

// SetGitContext sets the git context builder for git-related tools
func (te *ToolExecutorImpl) SetGitContext(gitContext domain.GitContextBuilder) {
	te.gitContext = gitContext
	te.rebuildHandlerRegistry()
}

// SetContextMemory sets the context memory for memory-related tools
func (te *ToolExecutorImpl) SetContextMemory(cm domain.ContextMemory) {
	te.contextMemory = cm
	te.rebuildHandlerRegistry()
}

// SetSemanticSearch sets the semantic search service
func (te *ToolExecutorImpl) SetSemanticSearch(ss tools.SemanticSearcher) {
	te.semanticSearcherService = ss
	te.hasSemanticSearch = ss != nil
	te.rebuildHandlerRegistry()
}

// SetAnalysisContainer configures the tool executor with all services from the container
func (te *ToolExecutorImpl) SetAnalysisContainer(container *appanalysis.Container) {
	if container == nil {
		return
	}

	te.registry = container.GetRegistry()
	te.symbolIndex = container.GetSymbolIndex()
	te.callGraph = container.GetCallGraph()
	te.gitContext = container.GetGitContext()
	te.contextMemory = container.GetContextMemory()
	te.referenceFinder = container.GetReferenceFinder()
	te.projectStructure = container.GetProjectStructure()

	if ss := container.GetSemanticSearch(); ss != nil {
		te.SetSemanticSearch(ss)
	}

	te.rebuildHandlerRegistry()
}

// rebuildHandlerRegistry rebuilds the handler registry with current services
func (te *ToolExecutorImpl) rebuildHandlerRegistry() {
	te.handlerRegistry = tools.NewHandlerRegistry(te.logger)

	// File tools
	te.handlerRegistry.Register(tools.NewFileToolsHandler(te.logger, te.fileReader))

	// Symbol tools
	te.handlerRegistry.Register(tools.NewSymbolToolsHandler(te.registry, te.symbolIndex, te.logger, te.referenceFinder))

	// Call graph tools
	te.handlerRegistry.Register(tools.NewCallGraphToolsHandler(te.logger, te.callGraph))

	// Git tools
	te.handlerRegistry.Register(tools.NewGitToolsHandler(te.logger, te.gitContext))

	// Memory tools
	te.handlerRegistry.Register(tools.NewMemoryToolsHandler(te.logger, te.contextMemory))

	// Preferences tools
	te.handlerRegistry.Register(tools.NewPreferencesToolsHandler(te.logger, te.contextMemory))

	// Project structure tools
	if te.projectStructure != nil {
		projectStructureService := project.NewStructureService(te.logger, te.projectStructure)
		te.handlerRegistry.Register(tools.NewProjectStructureToolsHandler(te.logger, projectStructureService))
	}

	// Semantic tools (if available)
	if te.hasSemanticSearch && te.semanticSearcherService != nil {
		te.handlerRegistry.Register(tools.NewSemanticToolsHandler(te.logger, te.semanticSearcherService))
	}
}

// GetAvailableTools returns all available tools from registered handlers
func (te *ToolExecutorImpl) GetAvailableTools() []domain.Tool {
	return te.handlerRegistry.GetAllTools()
}

// ExecuteTool executes a tool and returns the result
func (te *ToolExecutorImpl) ExecuteTool(call domain.ToolCall, projectRoot string) domain.ToolResult {
	te.logger.Info(fmt.Sprintf("Executing tool: %s", call.Name))

	content, err := te.handlerRegistry.Execute(call.Name, call.Arguments, projectRoot)

	result := domain.ToolResult{ToolCallID: call.ID, Content: content}
	if err != nil {
		result.Error = err.Error()
		result.Content = fmt.Sprintf("Error: %s", err.Error())
	}
	return result
}
