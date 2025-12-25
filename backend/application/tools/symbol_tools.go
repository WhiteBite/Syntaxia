package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// SymbolToolsHandler handles symbol-related tools
type SymbolToolsHandler struct {
	registry        analysis.AnalyzerRegistry
	symbolIndex     analysis.SymbolIndex
	logger          domain.Logger
	referenceFinder domain.ReferenceFinder
}

// NewSymbolToolsHandler creates a new SymbolToolsHandler
func NewSymbolToolsHandler(registry analysis.AnalyzerRegistry, symbolIndex analysis.SymbolIndex, logger domain.Logger, referenceFinder domain.ReferenceFinder) *SymbolToolsHandler {
	return &SymbolToolsHandler{
		registry:        registry,
		symbolIndex:     symbolIndex,
		logger:          logger,
		referenceFinder: referenceFinder,
	}
}

// ListSymbols extracts symbols from a file using language analyzers
func (h *SymbolToolsHandler) ListSymbols(args map[string]any, projectRoot string) (interface{}, error) {
	path, _ := args["path"].(string)
	kindFilter, _ := args["kind"].(string)

	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	fullPath := filepath.Join(projectRoot, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	analyzer := h.registry.GetAnalyzer(path)
	if analyzer == nil {
		return fmt.Sprintf("No analyzer available for file type: %s", filepath.Ext(path)), nil
	}

	symbols, err := analyzer.ExtractSymbols(context.Background(), path, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract symbols: %w", err)
	}

	// Filter by kind if specified
	if kindFilter != "" {
		var filtered []analysis.Symbol
		for _, s := range symbols {
			if strings.EqualFold(string(s.Kind), kindFilter) {
				filtered = append(filtered, s)
			}
		}
		symbols = filtered
	}

	if len(symbols) == 0 {
		return fmt.Sprintf("No symbols found in %s", path), nil
	}

	// Format output
	lines := make([]string, 0, len(symbols)+1)
	lines = append(lines, fmt.Sprintf("Symbols in %s (%s):", path, analyzer.Language()))
	for _, s := range symbols {
		line := fmt.Sprintf("  [%s] %s", s.Kind, s.Name)
		if s.StartLine > 0 {
			line += fmt.Sprintf(" (line %d)", s.StartLine)
		}
		if s.Parent != "" {
			line += fmt.Sprintf(" <- %s", s.Parent)
		}
		if s.Signature != "" {
			line += fmt.Sprintf("\n       %s", s.Signature)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// SearchSymbols searches for symbols across the project
func (h *SymbolToolsHandler) SearchSymbols(args map[string]any, projectRoot string) (interface{}, error) {
	query, _ := args["query"].(string)
	kindFilter, _ := args["kind"].(string)

	if query == "" {
		return nil, fmt.Errorf("query is required")
	}

	// Ensure index is built
	if !h.symbolIndex.IsIndexed() {
		h.logger.Info("Building symbol index...")
		if err := h.symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
		stats := h.symbolIndex.Stats()
		h.logger.Info(fmt.Sprintf("Indexed %d symbols in %d files", stats["total_symbols"], stats["files"]))
	}

	// Search
	results := h.symbolIndex.SearchByName(query)

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

	if len(results) == 0 {
		return fmt.Sprintf("No symbols found matching '%s'", query), nil
	}

	// Limit results
	if len(results) > 30 {
		results = results[:30]
	}

	// Format output
	lines := make([]string, 0, len(results)+1)
	lines = append(lines, fmt.Sprintf("Found %d symbols matching '%s':", len(results), query))
	for _, s := range results {
		line := fmt.Sprintf("  [%s] %s", s.Kind, s.Name)
		if s.FilePath != "" {
			line += fmt.Sprintf(" in %s", s.FilePath)
		}
		if s.StartLine > 0 {
			line += fmt.Sprintf(":%d", s.StartLine)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

// FindDefinition finds where a symbol is defined
func (h *SymbolToolsHandler) FindDefinition(args map[string]any, projectRoot string) (interface{}, error) {
	name, _ := args["name"].(string)
	kindFilter, _ := args["kind"].(string)

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Ensure index is built
	if !h.symbolIndex.IsIndexed() {
		h.logger.Info("Building symbol index...")
		if err := h.symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	// Find exact match
	var kind analysis.SymbolKind
	if kindFilter != "" {
		kind = analysis.SymbolKind(kindFilter)
	}

	sym := h.symbolIndex.FindDefinition(name, kind)
	if sym == nil {
		return fmt.Sprintf("No definition found for '%s'", name), nil
	}

	// Read the relevant code
	fullPath := filepath.Join(projectRoot, sym.FilePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Sprintf("Found %s '%s' in %s:%d but couldn't read file", sym.Kind, sym.Name, sym.FilePath, sym.StartLine), nil
	}

	lines := strings.Split(string(content), "\n")
	startLine := sym.StartLine - 1
	if startLine < 0 {
		startLine = 0
	}
	endLine := sym.EndLine
	if endLine == 0 || endLine > len(lines) {
		endLine = startLine + 20 // Show 20 lines by default
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Format output
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Definition of %s '%s':\n", sym.Kind, sym.Name))
	result.WriteString(fmt.Sprintf("File: %s\n", sym.FilePath))
	result.WriteString(fmt.Sprintf("Lines: %d-%d\n\n", sym.StartLine, endLine))

	for i := startLine; i < endLine; i++ {
		result.WriteString(fmt.Sprintf("%4d | %s\n", i+1, lines[i]))
	}

	return result.String(), nil
}

// FindReferences finds all references to a symbol across the project
func (h *SymbolToolsHandler) FindReferences(args map[string]any, projectRoot string) (interface{}, error) {
	name, _ := args["name"].(string)
	kindFilter, _ := args["kind"].(string)
	includeDefinition := true
	if incDef, ok := args["include_definition"].(bool); ok {
		includeDefinition = incDef
	}

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if h.referenceFinder == nil {
		return nil, fmt.Errorf("reference finder not initialized")
	}

	refs, err := h.referenceFinder.FindReferences(context.Background(), projectRoot, name, kindFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to find references: %w", err)
	}

	// Filter out definitions if requested
	if !includeDefinition {
		var filtered []domain.SymbolReference
		for _, ref := range refs {
			if !ref.IsDefinition {
				filtered = append(filtered, ref)
			}
		}
		refs = filtered
	}

	if len(refs) == 0 {
		return fmt.Sprintf("No references found for '%s'", name), nil
	}

	// Format output
	var result strings.Builder
	defCount := 0
	usageCount := 0
	for _, ref := range refs {
		if ref.IsDefinition {
			defCount++
		} else {
			usageCount++
		}
	}

	result.WriteString(fmt.Sprintf("Found %d references to '%s' (%d definitions, %d usages):\n\n",
		len(refs), name, defCount, usageCount))

	for _, ref := range refs {
		marker := "  "
		if ref.IsDefinition {
			marker = "* "
		}
		result.WriteString(fmt.Sprintf("%s%s:%d:%d\n", marker, ref.FilePath, ref.Line, ref.Column))
		result.WriteString(fmt.Sprintf("    %s\n", ref.LineText))
	}

	return result.String(), nil
}

// GetSymbolInfo returns detailed information about a symbol
func (h *SymbolToolsHandler) GetSymbolInfo(args map[string]any, projectRoot string) (interface{}, error) {
	name, _ := args["name"].(string)
	kindFilter, _ := args["kind"].(string)
	filePath, _ := args["file_path"].(string)

	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	// Ensure index is built
	if !h.symbolIndex.IsIndexed() {
		h.logger.Info("Building symbol index...")
		if err := h.symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	// Find symbol
	var kind analysis.SymbolKind
	if kindFilter != "" {
		kind = analysis.SymbolKind(kindFilter)
	}

	sym := h.symbolIndex.FindDefinition(name, kind)
	if sym == nil {
		// Try searching by name
		results := h.symbolIndex.SearchByName(name)
		if len(results) == 0 {
			return fmt.Sprintf("No symbol found: '%s'", name), nil
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

	// Read source code
	fullPath := filepath.Join(projectRoot, sym.FilePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		content = nil // Continue without source
	}

	// Format output
	var result strings.Builder
	result.WriteString(fmt.Sprintf("Symbol: %s\n", sym.Name))
	result.WriteString(fmt.Sprintf("Kind: %s\n", sym.Kind))
	result.WriteString(fmt.Sprintf("File: %s:%d\n", sym.FilePath, sym.StartLine))

	if sym.Parent != "" {
		result.WriteString(fmt.Sprintf("Parent: %s\n", sym.Parent))
	}

	if len(sym.Modifiers) > 0 {
		result.WriteString(fmt.Sprintf("Modifiers: %s\n", strings.Join(sym.Modifiers, ", ")))
	}

	if sym.Signature != "" {
		result.WriteString(fmt.Sprintf("\nSignature:\n  %s\n", sym.Signature))
	}

	if sym.DocComment != "" {
		result.WriteString(fmt.Sprintf("\nDocumentation:\n  %s\n", sym.DocComment))
	}

	// Show children if any
	if len(sym.Children) > 0 {
		result.WriteString(fmt.Sprintf("\nMembers (%d):\n", len(sym.Children)))
		for _, child := range sym.Children {
			result.WriteString(fmt.Sprintf("  [%s] %s", child.Kind, child.Name))
			if child.Signature != "" {
				result.WriteString(fmt.Sprintf(" - %s", child.Signature))
			}
			result.WriteString("\n")
		}
	}

	// Show source code
	if content != nil && sym.StartLine > 0 {
		lines := strings.Split(string(content), "\n")
		startLine := sym.StartLine - 1
		endLine := sym.EndLine
		if endLine == 0 || endLine > len(lines) {
			endLine = startLine + 15
		}
		if endLine > len(lines) {
			endLine = len(lines)
		}

		result.WriteString(fmt.Sprintf("\nSource (lines %d-%d):\n", sym.StartLine, endLine))
		for i := startLine; i < endLine; i++ {
			result.WriteString(fmt.Sprintf("%4d | %s\n", i+1, lines[i]))
		}
	}

	return result.String(), nil
}

// GetClassHierarchy returns class inheritance hierarchy
func (h *SymbolToolsHandler) GetClassHierarchy(args map[string]any, projectRoot string) (interface{}, error) {
	className, _ := args["class_name"].(string)
	direction, _ := args["direction"].(string)

	if className == "" {
		return nil, fmt.Errorf("class_name is required")
	}
	if direction == "" {
		direction = "both"
	}

	// Ensure index is built
	if !h.symbolIndex.IsIndexed() {
		if err := h.symbolIndex.IndexProject(context.Background(), projectRoot); err != nil {
			return nil, fmt.Errorf("failed to build index: %w", err)
		}
	}

	// Find the class
	sym := h.symbolIndex.FindDefinition(className, analysis.KindClass)
	if sym == nil {
		// Try interface
		sym = h.symbolIndex.FindDefinition(className, analysis.KindInterface)
	}
	if sym == nil {
		return fmt.Sprintf("Class or interface '%s' not found", className), nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Class Hierarchy for: %s\n", className))
	result.WriteString(fmt.Sprintf("Type: %s\n", sym.Kind))
	result.WriteString(fmt.Sprintf("File: %s:%d\n\n", sym.FilePath, sym.StartLine))

	// Get parent from symbol
	if (direction == "up" || direction == "both") && sym.Parent != "" {
		result.WriteString("Extends/Implements:\n")
		parents := strings.Split(sym.Parent, ",")
		for _, p := range parents {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			result.WriteString(fmt.Sprintf("  ↑ %s", p))
			// Try to find parent definition
			parentSym := h.symbolIndex.FindDefinition(p, "")
			if parentSym != nil {
				result.WriteString(fmt.Sprintf(" (%s:%d)", parentSym.FilePath, parentSym.StartLine))
			}
			result.WriteString("\n")
		}
		result.WriteString("\n")
	}

	// Find subclasses (classes that extend this one)
	if direction == "down" || direction == "both" {
		allClasses := h.symbolIndex.GetSymbolsByKind(analysis.KindClass)
		var subclasses []analysis.Symbol
		for _, c := range allClasses {
			if c.Parent != "" && strings.Contains(c.Parent, className) {
				subclasses = append(subclasses, c)
			}
		}

		if len(subclasses) > 0 {
			result.WriteString(fmt.Sprintf("Subclasses (%d):\n", len(subclasses)))
			for _, sub := range subclasses {
				result.WriteString(fmt.Sprintf("  ↓ %s (%s:%d)\n", sub.Name, sub.FilePath, sub.StartLine))
			}
		} else {
			result.WriteString("Subclasses: none found\n")
		}
	}

	return result.String(), nil
}

// symbolToolNames defines which tools this handler manages
var symbolToolNames = map[string]bool{
	"list_symbols":        true,
	"search_symbols":      true,
	"find_definition":     true,
	"find_references":     true,
	"get_symbol_info":     true,
	"get_class_hierarchy": true,
	"get_imports":         true,
}

// CanHandle returns true if this handler can handle the given tool
func (h *SymbolToolsHandler) CanHandle(toolName string) bool {
	return symbolToolNames[toolName]
}

// GetTools returns the list of symbol tools
func (h *SymbolToolsHandler) GetTools() []domain.Tool {
	return []domain.Tool{
		{
			Name:        "list_symbols",
			Description: "List all symbols (classes, functions, interfaces, types) in a file.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the source file"},
					"kind": {Type: "string", Description: "Filter by symbol kind: class, function, interface, type, method, enum"},
				},
				Required: []string{"path"},
			},
		},
		{
			Name:        "search_symbols",
			Description: "Search for symbols across the entire project by name.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"query": {Type: "string", Description: "Symbol name to search for (partial match supported)"},
					"kind":  {Type: "string", Description: "Filter by kind: class, function, interface, type, method, enum"},
				},
				Required: []string{"query"},
			},
		},
		{
			Name:        "find_definition",
			Description: "Find where a symbol is defined in the project.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"name": {Type: "string", Description: "Exact symbol name to find"},
					"kind": {Type: "string", Description: "Symbol kind: class, function, interface, type"},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "find_references",
			Description: "Find all references to a symbol across the project.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"name":               {Type: "string", Description: "Symbol name to find references for"},
					"kind":               {Type: "string", Description: "Symbol kind: class, function, interface, type"},
					"include_definition": {Type: "boolean", Description: "Whether to include the definition in results", Default: true},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "get_symbol_info",
			Description: "Get detailed information about a symbol including signature, documentation, modifiers.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"name":      {Type: "string", Description: "Symbol name to get info for"},
					"kind":      {Type: "string", Description: "Symbol kind: class, function, interface, type, method"},
					"file_path": {Type: "string", Description: "File path to narrow search"},
				},
				Required: []string{"name"},
			},
		},
		{
			Name:        "get_class_hierarchy",
			Description: "Get class inheritance hierarchy - parent classes, implemented interfaces, and subclasses.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"class_name": {Type: "string", Description: "Class name to analyze"},
					"direction":  {Type: "string", Description: "Direction: 'up' (ancestors), 'down' (descendants), 'both'"},
				},
				Required: []string{"class_name"},
			},
		},
		{
			Name:        "get_imports",
			Description: "Get all imports/dependencies of a file.",
			Parameters: domain.ToolParameters{
				Type: "object",
				Properties: map[string]domain.ToolProperty{
					"path": {Type: "string", Description: "Path to the source file"},
				},
				Required: []string{"path"},
			},
		},
	}
}

// Execute executes a symbol tool
func (h *SymbolToolsHandler) Execute(toolName string, args map[string]any, projectRoot string) (string, error) {
	var result any
	var err error

	switch toolName {
	case "list_symbols":
		result, err = h.ListSymbols(args, projectRoot)
	case "search_symbols":
		result, err = h.SearchSymbols(args, projectRoot)
	case "find_definition":
		result, err = h.FindDefinition(args, projectRoot)
	case "find_references":
		result, err = h.FindReferences(args, projectRoot)
	case "get_symbol_info":
		result, err = h.GetSymbolInfo(args, projectRoot)
	case "get_class_hierarchy":
		result, err = h.GetClassHierarchy(args, projectRoot)
	case "get_imports":
		result, err = h.GetImports(args, projectRoot)
	default:
		return "", fmt.Errorf("unknown symbol tool: %s", toolName)
	}

	if err != nil {
		return "", err
	}
	if str, ok := result.(string); ok {
		return str, nil
	}
	return fmt.Sprintf("%v", result), nil
}

// GetImports extracts imports from a file using language analyzers
func (h *SymbolToolsHandler) GetImports(args map[string]any, projectRoot string) (any, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	fullPath := filepath.Join(projectRoot, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}

	analyzer := h.registry.GetAnalyzer(path)
	if analyzer == nil {
		return fmt.Sprintf("No analyzer available for file type: %s", filepath.Ext(path)), nil
	}

	imports, err := analyzer.GetImports(context.Background(), path, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract imports: %w", err)
	}

	if len(imports) == 0 {
		return fmt.Sprintf("No imports found in %s", path), nil
	}

	// Separate local and external imports
	var local, external []analysis.Import
	for _, imp := range imports {
		if imp.IsLocal {
			local = append(local, imp)
		} else {
			external = append(external, imp)
		}
	}

	// Format output
	var lines []string
	lines = append(lines, fmt.Sprintf("Imports in %s:", path))
	if len(external) > 0 {
		lines = append(lines, "\nExternal:")
		for _, imp := range external {
			line := fmt.Sprintf("  %s", imp.Path)
			if imp.Alias != "" {
				line += fmt.Sprintf(" (as %s)", imp.Alias)
			}
			lines = append(lines, line)
		}
	}
	if len(local) > 0 {
		lines = append(lines, "\nLocal:")
		for _, imp := range local {
			line := fmt.Sprintf("  %s", imp.Path)
			if imp.Alias != "" {
				line += fmt.Sprintf(" (as %s)", imp.Alias)
			}
			lines = append(lines, line)
		}
	}

	return strings.Join(lines, "\n"), nil
}
