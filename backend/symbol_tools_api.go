package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

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

// ClassHierarchy represents class inheritance hierarchy
type ClassHierarchy struct {
	ClassName  string               `json:"className"`
	Kind       string               `json:"kind"`
	FilePath   string               `json:"filePath"`
	StartLine  int                  `json:"startLine"`
	Parents    []ClassHierarchyNode `json:"parents,omitempty"`
	Subclasses []ClassHierarchyNode `json:"subclasses,omitempty"`
}

// ClassHierarchyNode represents a node in class hierarchy
type ClassHierarchyNode struct {
	Name      string `json:"name"`
	FilePath  string `json:"filePath,omitempty"`
	StartLine int    `json:"startLine,omitempty"`
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
