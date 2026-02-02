package app

import (
	"context"
	"fmt"
	"syntaxia/application"
	appai "syntaxia/application/ai"
	"syntaxia/application/analysis"
	"syntaxia/application/rag"
	"syntaxia/domain"
	domainanalysis "syntaxia/domain/analysis"
	"syntaxia/handlers"
	"syntaxia/infrastructure/analyzers"
	"syntaxia/infrastructure/embeddings"
	"syntaxia/infrastructure/git"
	"syntaxia/infrastructure/memory"
	"syntaxia/infrastructure/projectstructure"
)

// initializeHandlers creates all handlers with proper dependencies
func (c *AppContainer) initializeHandlers() error {
	// Initialize Analysis Container with factory functions for DI
	analysisConfig := analysis.ContainerConfig{
		RegistryFactory: func() domainanalysis.AnalyzerRegistry {
			return analyzers.NewAnalyzerRegistry()
		},
		SymbolIndexFactory: func(registry domainanalysis.AnalyzerRegistry) domainanalysis.SymbolIndex {
			return analyzers.NewSymbolIndex(registry)
		},
		CallGraphFactory: func(registry domainanalysis.AnalyzerRegistry) domain.CallGraphBuilder {
			return &callGraphAdapter{impl: analyzers.NewCallGraphBuilder(registry)}
		},
		GitContextFactory: func(projectRoot string) domain.GitContextBuilder {
			return &gitContextAdapter{impl: git.NewContextBuilder(projectRoot)}
		},
		ContextMemoryFactory: func(contextDir string) (domain.ContextMemory, error) {
			return memory.NewContextMemory(contextDir)
		},
		ProjectStructureFactory: func() domain.ProjectStructureDetector {
			return &projectStructureAdapter{impl: projectstructure.NewDetector()}
		},
		ReferenceFinderFactory: func(registry domainanalysis.AnalyzerRegistry) domain.ReferenceFinder {
			return &referenceFinderAdapter{impl: analyzers.NewReferenceFinder(registry)}
		},
	}
	c.AnalysisContainer = analysis.NewContainer(c.Log, analysisConfig)

	// Initialize Dependency Analyzer
	c.DependencyAnalyzer = analysis.NewDependencyAnalyzer(
		c.AnalysisContainer.GetRegistry(),
		c.Log,
	)

	// Initialize Tool Executor with all dependencies
	c.ToolExecutor = application.NewToolExecutor(
		c.Log,
		c.FileReader,
		c.AnalysisContainer.GetRegistry(),
		c.AnalysisContainer.GetSymbolIndex(),
		c.AnalysisContainer.GetCallGraph(),
		c.AnalysisContainer.GetReferenceFinder(),
	)
	c.ToolExecutor.SetAnalysisContainer(c.AnalysisContainer)
	c.ToolExecutor.SetContextMemory(c.AnalysisContainer.GetContextMemory())

	// Wire sandbox for file write operations
	if c.SandboxFS != nil {
		c.ToolExecutor.SetSandboxFS(c.SandboxFS)
	}

	// Wire semantic search if available
	if c.SemanticSearch != nil {
		// Create adapter for SemanticSearcher interface
		c.ToolExecutor.SetSemanticSearch(&semanticSearchAdapter{service: c.SemanticSearch})
	}

	// Project Handler - delegates to ProjectService
	c.ProjectHandler = handlers.NewProjectHandler(
		c.Log,
		c.Bus,
		c.ProjectService,
		c.Watcher,
		c.FileReader,
		c.GitRepo,
	)

	// Context Handler - uses unified ContextService
	c.ContextHandler = handlers.NewContextHandler(
		c.Log,
		c.Bus,
		c.ContextService,
	)

	// AI Handler with injected ToolExecutor
	c.AIHandler = handlers.NewAIHandlerWithTools(
		c.Log,
		c.AIService,
		c.ContextAnalysis,
		c.ToolExecutor,
	)

	// Analysis Handler
	c.AnalysisHandler = handlers.NewAnalysisHandler(
		c.Log,
		c.TestService,
		c.StaticAnalyzerService,
		c.BuildService,
		c.SBOMService,
		c.SymbolGraph,
	)

	// Settings Handler
	c.SettingsHandler = handlers.NewSettingsHandler(
		c.Log,
		c.SettingsService,
	)

	// Taskflow Handler
	c.TaskflowHandler = handlers.NewTaskflowHandler(
		c.Log,
		c.TaskflowService,
		c.GuardrailService,
		c.RepairService,
		c.TaskProtocolService,
		c.TaskProtocolConfigService,
		c.TaskflowProtocolIntegration,
		c.BuildService,
	)

	// Qwen Task Service and Handler
	c.QwenTaskService = appai.NewQwenTaskService(
		c.Log,
		c.AIService,
		&smartContextAdapter{svc: c.SmartContextService},
		c.SettingsService,
	)

	c.QwenHandler = handlers.NewQwenHandler(
		c.Log,
		c.QwenTaskService,
	)

	// Semantic Handler (if semantic search is available)
	if c.SemanticSearch != nil && c.RAGService != nil {
		c.SemanticHandler = handlers.NewSemanticHandler(
			c.SemanticSearch,
			c.RAGService,
			c.Log,
		)
	}

	return nil
}

// =============================================================================
// Adapters for domain interfaces
// =============================================================================

// callGraphAdapter adapts analyzers.CallGraphBuilderImpl to domain.CallGraphBuilder
type callGraphAdapter struct {
	impl *analyzers.CallGraphBuilderImpl
}

func (a *callGraphAdapter) Build(projectRoot string) (*domain.CallGraph, error) {
	result, err := a.impl.Build(projectRoot)
	if err != nil {
		return nil, err
	}
	// Convert analysis.CallGraph to domain.CallGraph
	nodes := make(map[string]*domain.CallGraphNode)
	for k, v := range result.Nodes {
		nodes[k] = &domain.CallGraphNode{
			ID:       v.ID,
			Name:     v.Name,
			FilePath: v.FilePath,
			Line:     v.Line,
			Package:  v.Package,
		}
	}
	edges := make([]domain.CallGraphEdge, len(result.Edges))
	for i, e := range result.Edges {
		edges[i] = domain.CallGraphEdge{From: e.From, To: e.To, Line: e.Line}
	}
	return &domain.CallGraph{Nodes: nodes, Edges: edges}, nil
}

func (a *callGraphAdapter) GetCallers(functionID string) []domain.CallGraphNode {
	result := a.impl.GetCallers(functionID)
	nodes := make([]domain.CallGraphNode, len(result))
	for i, n := range result {
		nodes[i] = domain.CallGraphNode{ID: n.ID, Name: n.Name, FilePath: n.FilePath, Line: n.Line, Package: n.Package}
	}
	return nodes
}

func (a *callGraphAdapter) GetCallees(functionID string) []domain.CallGraphNode {
	result := a.impl.GetCallees(functionID)
	nodes := make([]domain.CallGraphNode, len(result))
	for i, n := range result {
		nodes[i] = domain.CallGraphNode{ID: n.ID, Name: n.Name, FilePath: n.FilePath, Line: n.Line, Package: n.Package}
	}
	return nodes
}

func (a *callGraphAdapter) GetImpact(functionID string, maxDepth int) []domain.CallGraphNode {
	result := a.impl.GetImpact(functionID, maxDepth)
	nodes := make([]domain.CallGraphNode, len(result))
	for i, n := range result {
		nodes[i] = domain.CallGraphNode{ID: n.ID, Name: n.Name, FilePath: n.FilePath, Line: n.Line, Package: n.Package}
	}
	return nodes
}

func (a *callGraphAdapter) GetCallChain(startID, endID string, maxDepth int) [][]string {
	return a.impl.GetCallChain(startID, endID, maxDepth)
}

// gitContextAdapter adapts git.ContextBuilder to domain.GitContextBuilder
type gitContextAdapter struct {
	impl *git.ContextBuilder
}

func (a *gitContextAdapter) GetRecentChanges(since string, pathFilter string) ([]domain.RecentChange, error) {
	result, err := a.impl.GetRecentChanges(since, pathFilter)
	if err != nil {
		return nil, err
	}
	changes := make([]domain.RecentChange, len(result))
	for i, r := range result {
		changes[i] = domain.RecentChange{
			FilePath:    r.FilePath,
			ChangeCount: r.ChangeCount,
			LastChanged: r.LastChanged,
			Authors:     r.Authors,
		}
	}
	return changes, nil
}

func (a *gitContextAdapter) GetCoChangedFiles(filePath string, limit int) ([]string, error) {
	return a.impl.GetCoChangedFiles(filePath, limit)
}

func (a *gitContextAdapter) SuggestContextFiles(taskDescription string, currentFiles []string, limit int) ([]string, error) {
	return a.impl.SuggestContextFiles(taskDescription, currentFiles, limit)
}

func (a *gitContextAdapter) GetRelatedByAuthor(filePath string, limit int) ([]string, error) {
	return a.impl.GetRelatedByAuthor(filePath, limit)
}

// projectStructureAdapter adapts projectstructure.Detector to domain.ProjectStructureDetector
type projectStructureAdapter struct {
	impl *projectstructure.Detector
}

func (a *projectStructureAdapter) Detect(projectPath string) (*domain.ProjectStructureInfo, error) {
	result, err := a.impl.DetectStructure(projectPath)
	if err != nil {
		return nil, err
	}
	langs := make([]string, len(result.Languages))
	for i, l := range result.Languages {
		langs[i] = l.Name
	}
	frameworks := make([]string, len(result.Frameworks))
	for i, f := range result.Frameworks {
		frameworks[i] = f.Name
	}
	return &domain.ProjectStructureInfo{
		Languages:  langs,
		Frameworks: frameworks,
	}, nil
}

func (a *projectStructureAdapter) DetectLanguages(projectPath string) ([]string, error) {
	result, err := a.impl.DetectStructure(projectPath)
	if err != nil {
		return nil, err
	}
	langs := make([]string, len(result.Languages))
	for i, l := range result.Languages {
		langs[i] = l.Name
	}
	return langs, nil
}

func (a *projectStructureAdapter) DetectFrameworks(projectPath string) ([]domain.FrameworkInfo, error) {
	return a.impl.DetectFrameworks(projectPath)
}

func (a *projectStructureAdapter) DetectStructure(projectPath string) (*domain.ProjectStructure, error) {
	return a.impl.DetectStructure(projectPath)
}

func (a *projectStructureAdapter) DetectArchitecture(projectPath string) (*domain.ArchitectureInfo, error) {
	return a.impl.DetectArchitecture(projectPath)
}

func (a *projectStructureAdapter) DetectConventions(projectPath string) (*domain.ConventionInfo, error) {
	return a.impl.DetectConventions(projectPath)
}

func (a *projectStructureAdapter) GetRelatedLayers(projectPath, filePath string) ([]domain.LayerInfo, error) {
	return a.impl.GetRelatedLayers(projectPath, filePath)
}

func (a *projectStructureAdapter) SuggestRelatedFiles(projectPath, filePath string) ([]string, error) {
	return a.impl.SuggestRelatedFiles(projectPath, filePath)
}

// referenceFinderAdapter adapts analyzers.ReferenceFinder to domain.ReferenceFinder
type referenceFinderAdapter struct {
	impl *analyzers.ReferenceFinder
}

func (a *referenceFinderAdapter) FindReferences(ctx context.Context, projectRoot string, symbolName string, symbolKind string) ([]domain.SymbolReference, error) {
	kind := domainanalysis.SymbolKind(symbolKind)
	result, err := a.impl.FindReferences(ctx, projectRoot, symbolName, kind)
	if err != nil {
		return nil, err
	}
	refs := make([]domain.SymbolReference, len(result))
	for i, r := range result {
		refs[i] = domain.SymbolReference{
			FilePath:     r.FilePath,
			Line:         r.Line,
			Column:       r.Column,
			LineText:     r.LineText,
			Context:      r.Context,
			IsDefinition: r.IsDefinition,
		}
	}
	return refs, nil
}

func (a *referenceFinderAdapter) FindUsages(ctx context.Context, projectRoot string, symbolName string) ([]domain.SymbolReference, error) {
	refs, err := a.FindReferences(ctx, projectRoot, symbolName, "")
	if err != nil {
		return nil, err
	}
	// Filter out definitions
	usages := make([]domain.SymbolReference, 0, len(refs))
	for _, r := range refs {
		if !r.IsDefinition {
			usages = append(usages, r)
		}
	}
	return usages, nil
}

// codeChunkerAdapter adapts embeddings.CodeChunker to domain.CodeChunker
type codeChunkerAdapter struct {
	impl *embeddings.CodeChunker
}

func (a *codeChunkerAdapter) ChunkFile(filePath string, content []byte, symbols []domain.ChunkSymbolInfo) []domain.CodeChunk {
	// Convert domain symbols to embeddings symbols
	embSymbols := make([]embeddings.SymbolInfo, len(symbols))
	for i, s := range symbols {
		embSymbols[i] = embeddings.SymbolInfo{
			Name:      s.Name,
			Kind:      s.Kind,
			StartLine: s.StartLine,
			EndLine:   s.EndLine,
		}
	}
	return a.impl.ChunkFile(filePath, content, embSymbols)
}

// semanticSearchAdapter adapts domain.SemanticSearchService to tools.SemanticSearcher
type semanticSearchAdapter struct {
	service     domain.SemanticSearchService
	projectRoot string
}

func (a *semanticSearchAdapter) Search(query string, limit int) ([]domain.SemanticSearchResult, error) {
	if a.service == nil {
		return nil, fmt.Errorf("semantic search service not available")
	}

	req := domain.SemanticSearchRequest{
		Query:       query,
		ProjectRoot: a.projectRoot,
		TopK:        limit,
		MinScore:    0.5,
		SearchType:  domain.SearchTypeSemantic,
	}

	resp, err := a.service.Search(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.Results) == 0 {
		return nil, nil
	}

	return resp.Results, nil
}

// SetProjectRoot updates the project root for semantic search
func (a *semanticSearchAdapter) SetProjectRoot(projectRoot string) {
	a.projectRoot = projectRoot
}

// smartContextAdapter adapts rag.SmartContextService to appai.SmartContextProvider
type smartContextAdapter struct {
	svc *rag.SmartContextService
}

func (a *smartContextAdapter) CollectContext(ctx context.Context, req appai.SmartContextRequest) (*appai.SmartContextResult, error) {
	ragReq := rag.SmartContextRequest{
		ProjectRoot:   req.ProjectRoot,
		Task:          req.Task,
		SelectedFiles: req.SelectedFiles,
		SelectedCode:  req.SelectedCode,
		SourceFile:    req.SourceFile,
		MaxTokens:     req.MaxTokens,
		MaxDepth:      req.MaxDepth,
		Language:      req.Language,
	}
	result, err := a.svc.CollectContext(ctx, ragReq)
	if err != nil {
		return nil, err
	}
	files := make([]appai.ContextFile, len(result.Files))
	for i, f := range result.Files {
		files[i] = appai.ContextFile{Path: f.Path, Content: f.Content, Tokens: f.Tokens, Relevance: f.Relevance, Reason: f.Reason}
	}
	var callStack *appai.CallStackResult
	if result.CallStack != nil {
		callStack = &appai.CallStackResult{
			RootSymbol:   result.CallStack.RootSymbol,
			Callers:      result.CallStack.Callers,
			Callees:      result.CallStack.Callees,
			Dependencies: result.CallStack.Dependencies,
			RelatedFiles: result.CallStack.RelatedFiles,
			TotalSymbols: result.CallStack.TotalSymbols,
		}
	}
	return &appai.SmartContextResult{
		Context:         result.Context,
		Files:           files,
		Symbols:         result.Symbols,
		CallStack:       callStack,
		TokenEstimate:   result.TokenEstimate,
		TruncatedFiles:  result.TruncatedFiles,
		ExcludedFiles:   result.ExcludedFiles,
		RelevanceScores: result.RelevanceScores,
	}, nil
}
