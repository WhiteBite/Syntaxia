package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// RequestsyntaxiaContextGeneration generates context for selected files
func (a *App) RequestsyntaxiaContextGeneration(rootDir string, includedPaths []string) {
	a.projectHandler.GenerateContext(a.ctx, rootDir, includedPaths)
}

// ExportContext exports context with specified settings
func (a *App) ExportContext(settingsJson string) (domain.ExportResult, error) {
	var settings domain.ExportSettings
	if err := json.Unmarshal([]byte(settingsJson), &settings); err != nil {
		validationErr := domain.NewValidationError("failed to parse export settings", map[string]interface{}{
			"originalError": err.Error(),
			"settingsJson":  settingsJson,
		})
		return domain.ExportResult{}, a.transformError(validationErr)
	}

	result, err := a.exportService.Export(a.ctx, settings)
	if err != nil {
		return domain.ExportResult{}, a.transformError(err)
	}

	return result, nil
}

// CleanupTempFiles cleans up temporary export files
func (a *App) CleanupTempFiles(filePath string) error {
	if filePath == "" {
		return nil
	}

	if !strings.Contains(filePath, "syntaxia-export-") {
		return fmt.Errorf("not a temp export file")
	}

	tempDir := filepath.Dir(filePath)
	return os.RemoveAll(tempDir)
}

// BuildContext builds context and returns ContextSummary (OOM-safe)
func (a *App) BuildContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	var options domain.ContextBuildOptions
	if strings.TrimSpace(optionsJson) != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			validationErr := domain.NewValidationError("failed to parse options JSON", map[string]interface{}{
				"originalError": err.Error(),
				"optionsJson":   optionsJson,
			})
			return "", a.transformError(validationErr)
		}
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContext called with empty includedPaths - this may include all project files")
	}

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, &options)
	if err != nil {
		return "", a.transformError(err)
	}

	contextJson, err := json.Marshal(summary)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summary", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextJson), nil
}

// GetContextContent returns paginated context content for memory-safe viewing
func (a *App) GetContextContent(contextID string, startLine int, lineCount int) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if lineCount <= 0 {
		lineCount = 1000
	}

	chunk, err := a.contextService.ReadContextChunk(a.ctx, contextID, startLine, lineCount)
	if err != nil {
		return "", a.transformError(err)
	}

	chunkJson, err := json.Marshal(chunk)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context chunk", err)
		return "", a.transformError(marshalErr)
	}

	return string(chunkJson), nil
}

// GetFullContextContent returns the full context content as a string
func (a *App) GetFullContextContent(contextID string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	content, err := a.contextService.ReadContextContent(a.ctx, contextID)
	if err != nil {
		return "", a.transformError(err)
	}

	return content, nil
}

// BuildContextLegacy is deprecated - use BuildContext instead
func (a *App) BuildContextLegacy() (string, error) {
	return "", a.transformError(domain.NewConfigurationError("legacy context building is no longer supported", nil))
}

// GetContext retrieves context metadata by ID
func (a *App) GetContext(contextID string) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	summary, err := a.contextService.GetContextSummary(a.ctx, contextID)
	if err != nil {
		return "", a.transformError(err)
	}

	contextJson, err := json.Marshal(summary)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summary", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextJson), nil
}

// contextSummaryJSON is a JSON-friendly version of ContextSummary with name field
type contextSummaryJSON struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name,omitempty"`
	ProjectPath string                 `json:"projectPath"`
	FileCount   int                    `json:"fileCount"`
	TotalSize   int64                  `json:"totalSize"`
	TokenCount  int                    `json:"tokenCount"`
	LineCount   int                    `json:"lineCount"`
	CreatedAt   string                 `json:"createdAt"`
	Metadata    *domain.ContextMetadata `json:"metadata,omitempty"`
}

// GetProjectContexts lists all stored context summaries for a project path
// Combines contexts from both JSON files (built contexts) and SQLite (saved contexts)
func (a *App) GetProjectContexts(projectPath string) (string, error) {
	var allSummaries []contextSummaryJSON

	// 1. Get contexts from JSON files (built contexts)
	if a.contextService != nil {
		summaries, err := a.contextService.GetProjectContextSummaries(a.ctx, projectPath)
		if err == nil && summaries != nil {
			for _, s := range summaries {
				allSummaries = append(allSummaries, contextSummaryJSON{
					ID:          s.ID,
					ProjectPath: s.ProjectPath,
					FileCount:   s.FileCount,
					TotalSize:   s.TotalSize,
					TokenCount:  s.TokenCount,
					LineCount:   s.LineCount,
					CreatedAt:   s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
					Metadata:    &s.Metadata,
				})
			}
		}
	}

	// 2. Get contexts from SQLite (saved contexts via SaveContextMemory)
	if a.analysisContainer != nil {
		contextMemory := a.analysisContainer.GetContextMemory()
		if contextMemory != nil {
			savedContexts, err := contextMemory.GetRecentContexts(projectPath, 100)
			if err == nil && savedContexts != nil {
				// Build set of existing IDs to avoid duplicates
				existingIDs := make(map[string]bool)
				for _, s := range allSummaries {
					existingIDs[s.ID] = true
				}

				for _, ctx := range savedContexts {
					// Skip if already exists in JSON summaries
					if existingIDs[ctx.ID] {
						continue
					}

					allSummaries = append(allSummaries, contextSummaryJSON{
						ID:          ctx.ID,
						Name:        ctx.Topic, // Topic becomes Name for frontend
						ProjectPath: ctx.ProjectRoot,
						FileCount:   len(ctx.Files),
						TotalSize:   0,
						LineCount:   0,
						CreatedAt:   ctx.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
						Metadata: &domain.ContextMetadata{
							SelectedFiles: ctx.Files,
						},
					})
				}
			}
		}
	}

	contextsJson, err := json.Marshal(allSummaries)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context summaries", err)
		return "", a.transformError(marshalErr)
	}

	return string(contextsJson), nil
}

// DeleteContext removes context metadata and associated content from disk
func (a *App) DeleteContext(contextID string) error {
	if a.contextService == nil {
		return a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if err := a.contextService.DeleteContext(a.ctx, contextID); err != nil {
		return a.transformError(err)
	}

	return nil
}

// BuildContextFromRequest builds context using provided options and returns ContextSummary
func (a *App) BuildContextFromRequest(projectPath string, includedPaths []string, options *domain.ContextBuildOptions) (*domain.ContextSummary, error) {
	if a.contextService == nil {
		return nil, a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if len(includedPaths) == 0 {
		a.log.Warning("BuildContextFromRequest called with empty includedPaths")
	}

	if options == nil {
		options = &domain.ContextBuildOptions{}
	}

	a.log.Info(fmt.Sprintf("[BuildContextFromRequest] OutputFormat received: '%s', StripComments: %v, MaxTokens: %d",
		options.OutputFormat, options.StripComments, options.MaxTokens))

	summary, err := a.contextService.BuildContextSummary(a.ctx, projectPath, includedPaths, options)
	if err != nil {
		return nil, a.transformError(err)
	}

	return summary, nil
}

// GetContextLines returns a chunk of context content between startLine and endLine inclusive
func (a *App) GetContextLines(contextID string, startLine, endLine int64) (string, error) {
	if a.contextService == nil {
		return "", a.transformError(domain.NewConfigurationError("context service not available", nil))
	}

	if endLine < startLine {
		return "", a.transformError(domain.NewValidationError("invalid line range", map[string]interface{}{
			"startLine": startLine,
			"endLine":   endLine,
		}))
	}

	lineCount := int(endLine-startLine) + 1
	chunk, err := a.contextService.ReadContextChunk(a.ctx, contextID, int(startLine), lineCount)
	if err != nil {
		return "", a.transformError(err)
	}

	chunkJson, err := json.Marshal(chunk)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal context chunk", err)
		return "", a.transformError(marshalErr)
	}

	return string(chunkJson), nil
}

// CreateStreamingContext delegates to BuildContext to create a disk-backed context summary
func (a *App) CreateStreamingContext(projectPath string, includedPaths []string, optionsJson string) (string, error) {
	return a.BuildContext(projectPath, includedPaths, optionsJson)
}

// GetStreamingContext returns context summary metadata for streaming compatibility
func (a *App) GetStreamingContext(contextID string) (string, error) {
	return a.GetContext(contextID)
}

// CloseStreamingContext removes a streaming context and associated resources
func (a *App) CloseStreamingContext(contextID string) error {
	return a.DeleteContext(contextID)
}

// ExportProject exports an entire project
func (a *App) ExportProject(projectPath string, format string, optionsJson string) (string, error) {
	var options map[string]interface{}
	if optionsJson != "" {
		if err := json.Unmarshal([]byte(optionsJson), &options); err != nil {
			return "", fmt.Errorf("failed to parse options JSON: %w", err)
		}
	}

	exportSettings := domain.ExportSettings{
		ProjectPath: projectPath,
		Format:      format,
		Options:     options,
	}

	result, err := a.exportService.Export(a.ctx, exportSettings)
	if err != nil {
		return "", fmt.Errorf("failed to export project: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export result: %w", err)
	}

	return string(resultJson), nil
}

// GetExportHistory returns export history
func (a *App) GetExportHistory(projectPath string) (string, error) {
	history, err := a.exportService.GetExportHistory(a.ctx, projectPath)
	if err != nil {
		return "", fmt.Errorf("failed to get export history: %w", err)
	}

	historyJson, err := json.Marshal(history)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export history: %w", err)
	}

	return string(historyJson), nil
}

// CollectSmartContext собирает релевантный контекст для AI задачи
func (a *App) CollectSmartContext(requestJson string) (string, error) {
	var req domain.SmartContextRequest
	if err := json.Unmarshal([]byte(requestJson), &req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	if a.container.SmartContextCollector == nil {
		return "", fmt.Errorf("smart context collector not available")
	}

	result, err := a.container.SmartContextCollector.CollectContext(a.ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to collect context: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJson), nil
}


// ============ Symbol Tools API ============

// SymbolInfo represents a symbol for frontend consumption
type SymbolInfo struct {
	Name      string   `json:"name"`
	Kind      string   `json:"kind"`
	FilePath  string   `json:"filePath"`
	StartLine int      `json:"startLine"`
	EndLine   int      `json:"endLine"`
	Signature string   `json:"signature,omitempty"`
	Parent    string   `json:"parent,omitempty"`
	Modifiers []string `json:"modifiers,omitempty"`
}

// SymbolLocation represents a symbol definition location
type SymbolLocation struct {
	FilePath  string `json:"filePath"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	StartCol  int    `json:"startCol"`
	EndCol    int    `json:"endCol"`
}

// ListSymbols returns symbols for a file
func (a *App) ListSymbols(projectRoot, filePath string) ([]SymbolInfo, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if filePath == "" {
		return nil, fmt.Errorf("filePath is required")
	}

	// Validate path doesn't escape project root
	if err := validatePathWithinProject(projectRoot, filePath); err != nil {
		return nil, err
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	registry := a.analysisContainer.GetRegistry()
	if registry == nil {
		return nil, fmt.Errorf("analyzer registry not available")
	}

	analyzer := registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil, fmt.Errorf("no analyzer available for file type: %s", filepath.Ext(filePath))
	}

	fullPath := filepath.Join(projectRoot, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	symbols, err := analyzer.ExtractSymbols(context.Background(), filePath, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract symbols: %w", err)
	}

	result := make([]SymbolInfo, 0, len(symbols))
	for _, s := range symbols {
		result = append(result, SymbolInfo{
			Name:      s.Name,
			Kind:      string(s.Kind),
			FilePath:  s.FilePath,
			StartLine: s.StartLine,
			EndLine:   s.EndLine,
			Signature: s.Signature,
			Parent:    s.Parent,
			Modifiers: s.Modifiers,
		})
	}

	return result, nil
}

// SearchSymbols searches for symbols by name across the project
func (a *App) SearchSymbols(projectRoot, query, kindFilter string) ([]SymbolInfo, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	symbolIndex := a.analysisContainer.GetSymbolIndex()
	if symbolIndex == nil {
		return nil, fmt.Errorf("symbol index not available")
	}

	// Ensure index is built
	if !symbolIndex.IsIndexed() {
		a.log.Info("Building symbol index...")
		if err := symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	results := symbolIndex.SearchByName(query)

	// Filter by kind if specified
	if kindFilter != "" {
		var filtered []analysis.Symbol
		for _, s := range results {
			if strings.EqualFold(string(s.Kind), kindFilter) {
				filtered = append(filtered, s)
			}
		}
		results = filtered
	}

	// Limit results
	const maxResults = 100
	if len(results) > maxResults {
		results = results[:maxResults]
	}

	result := make([]SymbolInfo, 0, len(results))
	for _, s := range results {
		result = append(result, SymbolInfo{
			Name:      s.Name,
			Kind:      string(s.Kind),
			FilePath:  s.FilePath,
			StartLine: s.StartLine,
			EndLine:   s.EndLine,
			Signature: s.Signature,
			Parent:    s.Parent,
			Modifiers: s.Modifiers,
		})
	}

	return result, nil
}

// GetSymbolDefinition returns definition location for a symbol
func (a *App) GetSymbolDefinition(projectRoot, symbolName, kindFilter string) (*SymbolLocation, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if symbolName == "" {
		return nil, fmt.Errorf("symbolName is required")
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	symbolIndex := a.analysisContainer.GetSymbolIndex()
	if symbolIndex == nil {
		return nil, fmt.Errorf("symbol index not available")
	}

	// Ensure index is built
	if !symbolIndex.IsIndexed() {
		a.log.Info("Building symbol index...")
		if err := symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	var kind analysis.SymbolKind
	if kindFilter != "" {
		kind = analysis.SymbolKind(kindFilter)
	}

	sym := symbolIndex.FindDefinition(symbolName, kind)
	if sym == nil {
		return nil, fmt.Errorf("no definition found for '%s'", symbolName)
	}

	return &SymbolLocation{
		FilePath:  sym.FilePath,
		StartLine: sym.StartLine,
		EndLine:   sym.EndLine,
		StartCol:  sym.StartCol,
		EndCol:    sym.EndCol,
	}, nil
}

// FindReferences finds all references to a symbol
func (a *App) FindReferences(projectRoot, symbolName, kindFilter string) ([]domain.SymbolReference, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if symbolName == "" {
		return nil, fmt.Errorf("symbolName is required")
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	referenceFinder := a.analysisContainer.GetReferenceFinder()
	if referenceFinder == nil {
		return nil, fmt.Errorf("reference finder not available")
	}

	refs, err := referenceFinder.FindReferences(context.Background(), projectRoot, symbolName, kindFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to find references: %w", err)
	}

	return refs, nil
}

// validatePathWithinProject validates that a path doesn't escape the project root
func validatePathWithinProject(projectRoot, path string) error {
	if strings.Contains(path, "..") {
		fullPath := filepath.Join(projectRoot, path)
		absProjectRoot, err := filepath.Abs(projectRoot)
		if err != nil {
			return fmt.Errorf("failed to resolve project root: %w", err)
		}
		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			return fmt.Errorf("failed to resolve file path: %w", err)
		}
		absProjectRoot = filepath.Clean(absProjectRoot)
		absFullPath = filepath.Clean(absFullPath)

		if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
			return fmt.Errorf("path traversal not allowed")
		}
	}
	return nil
}

// SymbolDetails represents detailed information about a symbol
type SymbolDetails struct {
	Name       string       `json:"name"`
	Kind       string       `json:"kind"`
	FilePath   string       `json:"filePath"`
	StartLine  int          `json:"startLine"`
	EndLine    int          `json:"endLine"`
	Signature  string       `json:"signature,omitempty"`
	Parent     string       `json:"parent,omitempty"`
	Modifiers  []string     `json:"modifiers,omitempty"`
	DocComment string       `json:"docComment,omitempty"`
	Children   []SymbolInfo `json:"children,omitempty"`
	SourceCode string       `json:"sourceCode,omitempty"`
}

// GetSymbolInfo returns detailed information about a symbol
func (a *App) GetSymbolInfo(projectRoot, symbolName, kindFilter, filePath string) (*SymbolDetails, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if symbolName == "" {
		return nil, fmt.Errorf("symbolName is required")
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	symbolIndex := a.analysisContainer.GetSymbolIndex()
	if symbolIndex == nil {
		return nil, fmt.Errorf("symbol index not available")
	}

	// Ensure index is built
	if !symbolIndex.IsIndexed() {
		a.log.Info("Building symbol index...")
		if err := symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	// Find symbol
	var kind analysis.SymbolKind
	if kindFilter != "" {
		kind = analysis.SymbolKind(kindFilter)
	}

	sym := symbolIndex.FindDefinition(symbolName, kind)
	if sym == nil {
		// Try searching by name
		results := symbolIndex.SearchByName(symbolName)
		if len(results) == 0 {
			return nil, fmt.Errorf("no symbol found: '%s'", symbolName)
		}
		// Filter by file path if provided
		if filePath != "" {
			for _, s := range results {
				if s.FilePath == filePath {
					sym = &s
					break
				}
			}
		}
		if sym == nil {
			sym = &results[0]
		}
	}

	// Build result
	result := &SymbolDetails{
		Name:       sym.Name,
		Kind:       string(sym.Kind),
		FilePath:   sym.FilePath,
		StartLine:  sym.StartLine,
		EndLine:    sym.EndLine,
		Signature:  sym.Signature,
		Parent:     sym.Parent,
		Modifiers:  sym.Modifiers,
		DocComment: sym.DocComment,
	}

	// Add children if any
	if len(sym.Children) > 0 {
		result.Children = make([]SymbolInfo, 0, len(sym.Children))
		for _, child := range sym.Children {
			result.Children = append(result.Children, SymbolInfo{
				Name:      child.Name,
				Kind:      string(child.Kind),
				FilePath:  child.FilePath,
				StartLine: child.StartLine,
				EndLine:   child.EndLine,
				Signature: child.Signature,
				Parent:    child.Parent,
				Modifiers: child.Modifiers,
			})
		}
	}

	// Read source code snippet
	if sym.StartLine > 0 {
		fullPath := filepath.Join(projectRoot, sym.FilePath)
		content, err := os.ReadFile(fullPath)
		if err == nil {
			lines := strings.Split(string(content), "\n")
			startLine := sym.StartLine - 1
			if startLine < 0 {
				startLine = 0
			}
			endLine := sym.EndLine
			if endLine == 0 || endLine > len(lines) {
				endLine = startLine + 15 // Show 15 lines by default
			}
			if endLine > len(lines) {
				endLine = len(lines)
			}
			if startLine < endLine {
				result.SourceCode = strings.Join(lines[startLine:endLine], "\n")
			}
		}
	}

	return result, nil
}

// ClassHierarchy represents class inheritance hierarchy
type ClassHierarchy struct {
	ClassName  string            `json:"className"`
	Kind       string            `json:"kind"`
	FilePath   string            `json:"filePath"`
	StartLine  int               `json:"startLine"`
	Parents    []ClassHierarchyNode `json:"parents,omitempty"`
	Subclasses []ClassHierarchyNode `json:"subclasses,omitempty"`
}

// ClassHierarchyNode represents a node in class hierarchy
type ClassHierarchyNode struct {
	Name      string `json:"name"`
	FilePath  string `json:"filePath,omitempty"`
	StartLine int    `json:"startLine,omitempty"`
}

// GetClassHierarchy returns class inheritance hierarchy
func (a *App) GetClassHierarchy(projectRoot, className, direction string) (*ClassHierarchy, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if className == "" {
		return nil, fmt.Errorf("className is required")
	}
	if direction == "" {
		direction = "both"
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	symbolIndex := a.analysisContainer.GetSymbolIndex()
	if symbolIndex == nil {
		return nil, fmt.Errorf("symbol index not available")
	}

	// Ensure index is built
	if !symbolIndex.IsIndexed() {
		a.log.Info("Building symbol index...")
		if err := symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	// Find the class
	sym := symbolIndex.FindDefinition(className, analysis.KindClass)
	if sym == nil {
		// Try interface
		sym = symbolIndex.FindDefinition(className, analysis.KindInterface)
	}
	if sym == nil {
		return nil, fmt.Errorf("class or interface '%s' not found", className)
	}

	result := &ClassHierarchy{
		ClassName: className,
		Kind:      string(sym.Kind),
		FilePath:  sym.FilePath,
		StartLine: sym.StartLine,
		Parents:   []ClassHierarchyNode{},
		Subclasses: []ClassHierarchyNode{},
	}

	// Get parents from symbol
	if (direction == "up" || direction == "both") && sym.Parent != "" {
		parents := strings.Split(sym.Parent, ",")
		for _, p := range parents {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			node := ClassHierarchyNode{Name: p}
			// Try to find parent definition
			parentSym := symbolIndex.FindDefinition(p, "")
			if parentSym != nil {
				node.FilePath = parentSym.FilePath
				node.StartLine = parentSym.StartLine
			}
			result.Parents = append(result.Parents, node)
		}
	}

	// Find subclasses (classes that extend this one)
	if direction == "down" || direction == "both" {
		allClasses := symbolIndex.GetSymbolsByKind(analysis.KindClass)
		for _, c := range allClasses {
			if c.Parent != "" && strings.Contains(c.Parent, className) {
				result.Subclasses = append(result.Subclasses, ClassHierarchyNode{
					Name:      c.Name,
					FilePath:  c.FilePath,
					StartLine: c.StartLine,
				})
			}
		}
	}

	return result, nil
}

// ImportInfo represents imports information for a file
type ImportInfo struct {
	FilePath        string       `json:"filePath"`
	ExternalImports []ImportItem `json:"externalImports"`
	LocalImports    []ImportItem `json:"localImports"`
	TotalCount      int          `json:"totalCount"`
}

// ImportItem represents a single import
type ImportItem struct {
	Path  string `json:"path"`
	Alias string `json:"alias,omitempty"`
}

// GetImports returns all imports/dependencies of a file
func (a *App) GetImports(projectRoot, filePath string) (*ImportInfo, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("projectRoot is required")
	}
	if filePath == "" {
		return nil, fmt.Errorf("filePath is required")
	}

	// Validate path doesn't escape project root
	if err := validatePathWithinProject(projectRoot, filePath); err != nil {
		return nil, err
	}

	if a.analysisContainer == nil {
		return nil, fmt.Errorf("analysis container not initialized")
	}

	registry := a.analysisContainer.GetRegistry()
	if registry == nil {
		return nil, fmt.Errorf("analyzer registry not available")
	}

	analyzer := registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil, fmt.Errorf("no analyzer available for file type: %s", filepath.Ext(filePath))
	}

	fullPath := filepath.Join(projectRoot, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	imports, err := analyzer.GetImports(context.Background(), filePath, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract imports: %w", err)
	}

	result := &ImportInfo{
		FilePath:        filePath,
		ExternalImports: []ImportItem{},
		LocalImports:    []ImportItem{},
		TotalCount:      len(imports),
	}

	for _, imp := range imports {
		item := ImportItem{
			Path:  imp.Path,
			Alias: imp.Alias,
		}
		if imp.IsLocal {
			result.LocalImports = append(result.LocalImports, item)
		} else {
			result.ExternalImports = append(result.ExternalImports, item)
		}
	}

	return result, nil
}
