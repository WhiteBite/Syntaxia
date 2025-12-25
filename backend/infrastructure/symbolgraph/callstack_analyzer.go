package symbolgraph

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

// CallStackAnalyzer analyzes call stacks and dependencies between symbols
type CallStackAnalyzer struct {
	log          domain.Logger
	graphBuilder *GoSymbolGraphBuilder
}

// NewCallStackAnalyzer creates a new call stack analyzer
func NewCallStackAnalyzer(log domain.Logger) *CallStackAnalyzer {
	return &CallStackAnalyzer{
		log:          log,
		graphBuilder: NewGoSymbolGraphBuilder(log),
	}
}

// CallStackEntry represents a single entry in the call stack
type CallStackEntry struct {
	Symbol     *domain.SymbolNode `json:"symbol"`
	CalledFrom []*CallStackEntry  `json:"calledFrom,omitempty"` // Who calls this symbol
	CallsTo    []*CallStackEntry  `json:"callsTo,omitempty"`    // What this symbol calls
	Depth      int                `json:"depth"`
}

// CallStackResult contains the full call stack analysis
type CallStackResult struct {
	RootSymbol   *domain.SymbolNode   `json:"rootSymbol"`
	Callers      []*domain.SymbolNode `json:"callers"`      // Functions that call the root
	Callees      []*domain.SymbolNode `json:"callees"`      // Functions called by the root
	Dependencies []*domain.SymbolNode `json:"dependencies"` // Types/interfaces used
	RelatedFiles []string             `json:"relatedFiles"` // All files involved
	TotalSymbols int                  `json:"totalSymbols"`
}

// AnalyzeCallStack analyzes the call stack for a given symbol
func (a *CallStackAnalyzer) AnalyzeCallStack(
	ctx context.Context,
	projectRoot string,
	filePath string,
	symbolName string,
	maxDepth int,
) (*CallStackResult, error) {
	a.log.Info(fmt.Sprintf("Analyzing call stack for symbol: %s in %s", symbolName, filePath))

	// Build the symbol graph first
	graph, err := a.graphBuilder.BuildGraph(ctx, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to build symbol graph: %w", err)
	}

	// Find the root symbol
	rootSymbol := a.findSymbol(graph, filePath, symbolName, projectRoot)
	if rootSymbol == nil {
		return nil, fmt.Errorf("symbol %s not found in %s", symbolName, filePath)
	}

	result := &CallStackResult{
		RootSymbol:   rootSymbol,
		Callers:      make([]*domain.SymbolNode, 0),
		Callees:      make([]*domain.SymbolNode, 0),
		Dependencies: make([]*domain.SymbolNode, 0),
		RelatedFiles: make([]string, 0),
	}

	// Build call graph with function calls
	callGraph, err := a.buildCallGraph(ctx, projectRoot)
	if err != nil {
		a.log.Warning(fmt.Sprintf("Failed to build call graph: %v", err))
	} else {
		// Find callers (who calls this function)
		result.Callers = a.findCallers(rootSymbol, callGraph, graph, maxDepth)

		// Find callees (what this function calls)
		result.Callees = a.findCallees(rootSymbol, callGraph, graph, maxDepth)
	}

	// Find type dependencies
	result.Dependencies = a.findDependencies(rootSymbol, graph)

	// Collect all related files
	fileSet := make(map[string]bool)
	fileSet[rootSymbol.Path] = true
	for _, s := range result.Callers {
		fileSet[s.Path] = true
	}
	for _, s := range result.Callees {
		fileSet[s.Path] = true
	}
	for _, s := range result.Dependencies {
		fileSet[s.Path] = true
	}

	for f := range fileSet {
		result.RelatedFiles = append(result.RelatedFiles, f)
	}

	result.TotalSymbols = len(result.Callers) + len(result.Callees) + len(result.Dependencies) + 1

	a.log.Info(fmt.Sprintf("Call stack analysis complete: %d callers, %d callees, %d dependencies",
		len(result.Callers), len(result.Callees), len(result.Dependencies)))

	return result, nil
}

// CallEdge represents a function call relationship
type CallEdge struct {
	Caller   string // Caller function ID
	Callee   string // Called function name
	FilePath string // File where the call occurs
	Line     int    // Line number of the call
}

// buildCallGraph builds a graph of function calls
func (a *CallStackAnalyzer) buildCallGraph(ctx context.Context, projectRoot string) ([]CallEdge, error) {
	var edges []CallEdge

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == "vendor" || info.Name() == "node_modules" ||
				info.Name() == ".git" || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fileEdges, err := a.parseCallsInFile(path, projectRoot)
		if err != nil {
			a.log.Warning(fmt.Sprintf("Failed to parse calls in %s: %v", path, err))
			return nil
		}

		edges = append(edges, fileEdges...)
		return nil
	})

	return edges, err
}

// parseCallsInFile extracts function calls from a Go file
func (a *CallStackAnalyzer) parseCallsInFile(filePath, projectRoot string) ([]CallEdge, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	relPath, _ := filepath.Rel(projectRoot, filePath)
	var edges []CallEdge
	var currentFunc string

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			// Track current function
			currentFunc = fmt.Sprintf("func:%s:%s", relPath, x.Name.Name)

		case *ast.CallExpr:
			if currentFunc == "" {
				return true
			}

			// Extract called function name
			calleeName := a.extractCalleeName(x)
			if calleeName != "" {
				pos := fset.Position(x.Pos())
				edges = append(edges, CallEdge{
					Caller:   currentFunc,
					Callee:   calleeName,
					FilePath: relPath,
					Line:     pos.Line,
				})
			}
		}
		return true
	})

	return edges, nil
}

// extractCalleeName extracts the function name from a call expression
func (a *CallStackAnalyzer) extractCalleeName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		// Simple function call: foo()
		return fn.Name
	case *ast.SelectorExpr:
		// Method call: obj.Method() or pkg.Func()
		if ident, ok := fn.X.(*ast.Ident); ok {
			return ident.Name + "." + fn.Sel.Name
		}
		return fn.Sel.Name
	}
	return ""
}

// findSymbol finds a symbol in the graph by file path and name
func (a *CallStackAnalyzer) findSymbol(graph *domain.SymbolGraph, filePath, symbolName, projectRoot string) *domain.SymbolNode {
	relPath, _ := filepath.Rel(projectRoot, filePath)
	if relPath == "" {
		relPath = filePath
	}

	// Normalize path separators
	relPath = strings.ReplaceAll(relPath, "\\", "/")

	for _, node := range graph.Nodes {
		nodePath := strings.ReplaceAll(node.Path, "\\", "/")
		if nodePath == relPath && node.Name == symbolName {
			return node
		}
	}

	// Try partial match
	for _, node := range graph.Nodes {
		if node.Name == symbolName && strings.HasSuffix(node.Path, filepath.Base(filePath)) {
			return node
		}
	}

	return nil
}

// findCallers finds all functions that call the given symbol
func (a *CallStackAnalyzer) findCallers(
	symbol *domain.SymbolNode,
	callGraph []CallEdge,
	symbolGraph *domain.SymbolGraph,
	maxDepth int,
) []*domain.SymbolNode {
	callers := make([]*domain.SymbolNode, 0)
	seen := make(map[string]bool)

	// Find direct callers
	for _, edge := range callGraph {
		if strings.HasSuffix(edge.Callee, symbol.Name) || edge.Callee == symbol.Name {
			if !seen[edge.Caller] {
				seen[edge.Caller] = true
				// Find the caller symbol in the graph
				for _, node := range symbolGraph.Nodes {
					if node.ID == edge.Caller {
						callers = append(callers, node)
						break
					}
				}
			}
		}
	}

	return callers
}

// findCallees finds all functions called by the given symbol
func (a *CallStackAnalyzer) findCallees(
	symbol *domain.SymbolNode,
	callGraph []CallEdge,
	symbolGraph *domain.SymbolGraph,
	maxDepth int,
) []*domain.SymbolNode {
	callees := make([]*domain.SymbolNode, 0)
	seen := make(map[string]bool)

	// Find direct callees
	for _, edge := range callGraph {
		if edge.Caller == symbol.ID {
			if !seen[edge.Callee] {
				seen[edge.Callee] = true
				// Find the callee symbol in the graph
				for _, node := range symbolGraph.Nodes {
					if node.Name == edge.Callee || strings.HasSuffix(node.Name, edge.Callee) {
						callees = append(callees, node)
						break
					}
				}
			}
		}
	}

	return callees
}

// buildNodeIDIndex creates a map of node ID to node for fast lookup
func buildNodeIDIndex(graph *domain.SymbolGraph) map[string]*domain.SymbolNode {
	index := make(map[string]*domain.SymbolNode, len(graph.Nodes))
	for _, node := range graph.Nodes {
		index[node.ID] = node
	}
	return index
}

// findEdgeDependencies finds dependencies from edges
func findEdgeDependencies(symbol *domain.SymbolNode, graph *domain.SymbolGraph, nodeIndex map[string]*domain.SymbolNode, seen map[string]bool) []*domain.SymbolNode {
	var deps []*domain.SymbolNode
	for _, edge := range graph.Edges {
		if edge.From != symbol.ID || (edge.Type != domain.EdgeTypeReferences && edge.Type != domain.EdgeTypeUses) {
			continue
		}
		if seen[edge.To] {
			continue
		}
		if node, ok := nodeIndex[edge.To]; ok {
			seen[edge.To] = true
			deps = append(deps, node)
		}
	}
	return deps
}

// findPackageTypeDependencies finds types in the same package
func findPackageTypeDependencies(symbol *domain.SymbolNode, graph *domain.SymbolGraph, seen map[string]bool) []*domain.SymbolNode {
	var deps []*domain.SymbolNode
	for _, node := range graph.Nodes {
		if node.Package != symbol.Package || seen[node.ID] {
			continue
		}
		if node.Type == domain.SymbolTypeStruct || node.Type == domain.SymbolTypeInterface {
			seen[node.ID] = true
			deps = append(deps, node)
		}
	}
	return deps
}

// findDependencies finds type dependencies for a symbol
func (a *CallStackAnalyzer) findDependencies(symbol *domain.SymbolNode, graph *domain.SymbolGraph) []*domain.SymbolNode {
	seen := make(map[string]bool)
	nodeIndex := buildNodeIDIndex(graph)

	deps := findEdgeDependencies(symbol, graph, nodeIndex, seen)
	deps = append(deps, findPackageTypeDependencies(symbol, graph, seen)...)
	return deps
}

// GetTransitiveDependencies returns all transitive dependencies up to maxDepth
func (a *CallStackAnalyzer) GetTransitiveDependencies(
	ctx context.Context,
	projectRoot string,
	filePath string,
	symbolName string,
	maxDepth int,
) ([]*domain.SymbolNode, error) {
	result, err := a.AnalyzeCallStack(ctx, projectRoot, filePath, symbolName, maxDepth)
	if err != nil {
		return nil, err
	}

	// Combine all related symbols
	allSymbols := make([]*domain.SymbolNode, 0)
	seen := make(map[string]bool)

	addUnique := func(symbols []*domain.SymbolNode) {
		for _, s := range symbols {
			if !seen[s.ID] {
				seen[s.ID] = true
				allSymbols = append(allSymbols, s)
			}
		}
	}

	addUnique(result.Callers)
	addUnique(result.Callees)
	addUnique(result.Dependencies)

	return allSymbols, nil
}
