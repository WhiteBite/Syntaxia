package analyzers

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"syntaxia/domain/analysis"
	"sort"
	"strings"
	"sync"
)

const (
	extGo          = ".go"
	dirVendor      = "vendor"
	dirNodeModules = "node_modules"
)

// CallGraphBuilderImpl builds call graphs
type CallGraphBuilderImpl struct {
	mu          sync.RWMutex
	registry    analysis.AnalyzerRegistry
	graph       *analysis.CallGraph
	depGraph    *analysis.DependencyGraph
	fileImports map[string][]importInfo // file -> imports

	// Caching fields for one-time initialization
	buildOnce    sync.Once
	lastBuildErr error
	projectRoot  string
	built        bool
}

type importInfo struct {
	path string
	line int
}

// NewCallGraphBuilder creates a new call graph builder
func NewCallGraphBuilder(registry analysis.AnalyzerRegistry) *CallGraphBuilderImpl {
	return &CallGraphBuilderImpl{
		registry: registry,
		graph: &analysis.CallGraph{
			Nodes: make(map[string]*analysis.CallNode),
			Edges: make([]analysis.CallEdge, 0),
		},
		depGraph: &analysis.DependencyGraph{
			Nodes: make(map[string]*analysis.DependencyNode),
			Edges: make([]analysis.DependencyEdge, 0),
		},
		fileImports: make(map[string][]importInfo),
	}
}

// EnsureBuilt ensures the call graph is built exactly once.
// Subsequent calls return immediately with cached result.
// Use Invalidate() to force rebuild.
func (b *CallGraphBuilderImpl) EnsureBuilt(projectRoot string) (*analysis.CallGraph, error) {
	// Check if project changed - need to rebuild
	b.mu.RLock()
	needsRebuild := b.projectRoot != "" && b.projectRoot != projectRoot
	b.mu.RUnlock()

	if needsRebuild {
		b.Invalidate()
	}

	b.buildOnce.Do(func() {
		b.mu.Lock()
		b.projectRoot = projectRoot
		b.mu.Unlock()
		_, b.lastBuildErr = b.Build(projectRoot)
		if b.lastBuildErr == nil {
			b.mu.Lock()
			b.built = true
			b.mu.Unlock()
		}
	})

	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.graph, b.lastBuildErr
}

// Invalidate resets the call graph, forcing rebuild on next EnsureBuilt call.
func (b *CallGraphBuilderImpl) Invalidate() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.graph = &analysis.CallGraph{
		Nodes: make(map[string]*analysis.CallNode),
		Edges: make([]analysis.CallEdge, 0),
	}
	b.depGraph = &analysis.DependencyGraph{
		Nodes: make(map[string]*analysis.DependencyNode),
		Edges: make([]analysis.DependencyEdge, 0),
	}
	b.fileImports = make(map[string][]importInfo)
	b.buildOnce = sync.Once{} // Reset sync.Once
	b.lastBuildErr = nil
	b.projectRoot = ""
	b.built = false
}

// IsBuilt returns whether the call graph has been built.
func (b *CallGraphBuilderImpl) IsBuilt() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.built
}

// GetProjectRoot returns the project root for which the graph was built.
func (b *CallGraphBuilderImpl) GetProjectRoot() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.projectRoot
}

// Build builds call graph for Go project
func (b *CallGraphBuilderImpl) Build(projectRoot string) (*analysis.CallGraph, error) {
	b.graph = &analysis.CallGraph{
		Nodes: make(map[string]*analysis.CallNode),
		Edges: make([]analysis.CallEdge, 0),
	}
	b.fileImports = make(map[string][]importInfo)

	// Walk project and analyze Go files
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == dirVendor || name == dirNodeModules {
				return filepath.SkipDir
			}
			return nil
		}

		ext := filepath.Ext(path)
		relPath, _ := filepath.Rel(projectRoot, path)

		switch ext {
		case extGo:
			b.analyzeGoFile(path, relPath)
		case ".ts", ".js", ".tsx", ".jsx":
			b.analyzeJSFile(path, relPath)
		case ".vue":
			b.analyzeVueFile(path, relPath)
		}

		return nil
	})

	return b.graph, err
}

// buildFuncSignature builds a function signature string
func buildFuncSignature(decl *ast.FuncDecl) string {
	var sig strings.Builder
	sig.WriteString("func ")
	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		sig.WriteString("(receiver) ")
	}
	sig.WriteString(decl.Name.Name + "(...)")
	return sig.String()
}

// collectFuncDeclarations collects all function declarations from a Go file
func (b *CallGraphBuilderImpl) collectFuncDeclarations(file *ast.File, fset *token.FileSet, pkgName, relPath string) {
	ast.Inspect(file, func(n ast.Node) bool {
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		nodeID := b.makeFunctionID(pkgName, decl.Name.Name, relPath)
		pos := fset.Position(decl.Pos())
		b.graph.Nodes[nodeID] = &analysis.CallNode{
			ID: nodeID, Name: decl.Name.Name, FilePath: relPath, Line: pos.Line,
			Package: pkgName, Signature: buildFuncSignature(decl),
			Callers: make([]string, 0), Callees: make([]string, 0),
		}
		return true
	})
}

// collectFuncCalls collects all function calls from a Go file
func (b *CallGraphBuilderImpl) collectFuncCalls(file *ast.File, fset *token.FileSet, pkgName, relPath string) {
	ast.Inspect(file, func(n ast.Node) bool {
		decl, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		callerID := b.makeFunctionID(pkgName, decl.Name.Name, relPath)
		ast.Inspect(decl.Body, func(inner ast.Node) bool {
			call, ok := inner.(*ast.CallExpr)
			if !ok {
				return true
			}
			calleeName := b.extractCallName(call)
			if calleeName == "" {
				return true
			}
			calleeID := b.makeFunctionID(pkgName, calleeName, "")
			pos := fset.Position(call.Pos())
			b.graph.Edges = append(b.graph.Edges, analysis.CallEdge{
				From: callerID, To: calleeID, FilePath: relPath, Line: pos.Line, CallType: "direct",
			})
			if caller, ok := b.graph.Nodes[callerID]; ok {
				caller.Callees = append(caller.Callees, calleeID)
			}
			if callee, ok := b.graph.Nodes[calleeID]; ok {
				callee.Callers = append(callee.Callers, callerID)
			}
			return true
		})
		return true
	})
}

func (b *CallGraphBuilderImpl) analyzeGoFile(path, relPath string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return
	}

	pkgName := ""
	if file.Name != nil {
		pkgName = file.Name.Name
	}

	b.collectFuncDeclarations(file, fset, pkgName, relPath)
	b.collectFuncCalls(file, fset, pkgName, relPath)
}

func (b *CallGraphBuilderImpl) extractCallName(call *ast.CallExpr) string {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		return fn.Name
	case *ast.SelectorExpr:
		return fn.Sel.Name
	}
	return ""
}

func (b *CallGraphBuilderImpl) analyzeJSFile(path, relPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	text := string(content)

	// Simple regex-based analysis for JS/TS
	funcRe := regexp.MustCompile(`(?m)(?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\([^)]*\)\s*=>|(\w+)\s*\([^)]*\)\s*\{)`)

	// Find function definitions
	funcMatches := funcRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range funcMatches {
		var name string
		for i := 2; i < len(match); i += 2 {
			if match[i] != -1 {
				name = text[match[i]:match[i+1]]
				break
			}
		}
		if name == "" || isJSKeyword(name) {
			continue
		}

		line := strings.Count(text[:match[0]], "\n") + 1
		nodeID := b.makeFunctionID("", name, relPath)

		b.graph.Nodes[nodeID] = &analysis.CallNode{
			ID:       nodeID,
			Name:     name,
			FilePath: relPath,
			Line:     line,
			Callers:  make([]string, 0),
			Callees:  make([]string, 0),
		}
	}

	// Analyze function calls using the shared analyzeJSCalls method
	// This properly tracks which function makes each call and updates Edges/Callers/Callees
	b.analyzeJSCalls(text, relPath)
}

func (b *CallGraphBuilderImpl) makeFunctionID(pkg, name, file string) string {
	if pkg != "" {
		return pkg + "." + name
	}
	if file != "" {
		return file + ":" + name
	}
	return name
}

// GetCallers returns functions that call the given function
func (b *CallGraphBuilderImpl) GetCallers(functionID string) []analysis.CallNode {
	node, ok := b.graph.Nodes[functionID]
	if !ok {
		return nil
	}

	var callers []analysis.CallNode
	for _, callerID := range node.Callers {
		if caller, ok := b.graph.Nodes[callerID]; ok {
			callers = append(callers, *caller)
		}
	}
	return callers
}

// GetCallees returns functions called by the given function
func (b *CallGraphBuilderImpl) GetCallees(functionID string) []analysis.CallNode {
	node, ok := b.graph.Nodes[functionID]
	if !ok {
		return nil
	}

	var callees []analysis.CallNode
	for _, calleeID := range node.Callees {
		if callee, ok := b.graph.Nodes[calleeID]; ok {
			callees = append(callees, *callee)
		}
	}
	return callees
}

// GetCallChain finds path between two functions
func (b *CallGraphBuilderImpl) GetCallChain(startID, endID string, maxDepth int) [][]string {
	var paths [][]string
	visited := make(map[string]bool)

	var dfs func(current string, path []string, depth int)
	dfs = func(current string, path []string, depth int) {
		if depth > maxDepth {
			return
		}
		if current == endID {
			pathCopy := make([]string, len(path))
			copy(pathCopy, path)
			paths = append(paths, pathCopy)
			return
		}
		if visited[current] {
			return
		}

		visited[current] = true
		node, ok := b.graph.Nodes[current]
		if ok {
			for _, calleeID := range node.Callees {
				dfs(calleeID, append(path, calleeID), depth+1)
			}
		}
		visited[current] = false
	}

	dfs(startID, []string{startID}, 0)
	return paths
}

// GetImpact returns all functions affected if given function changes
func (b *CallGraphBuilderImpl) GetImpact(functionID string, maxDepth int) []analysis.CallNode {
	affected := make(map[string]*analysis.CallNode)
	visited := make(map[string]bool)

	var traverse func(id string, depth int)
	traverse = func(id string, depth int) {
		if depth > maxDepth || visited[id] {
			return
		}
		visited[id] = true

		node, ok := b.graph.Nodes[id]
		if !ok {
			return
		}

		for _, callerID := range node.Callers {
			if caller, ok := b.graph.Nodes[callerID]; ok {
				affected[callerID] = caller
				traverse(callerID, depth+1)
			}
		}
	}

	traverse(functionID, 0)

	result := make([]analysis.CallNode, 0, len(affected))
	for _, node := range affected {
		result = append(result, *node)
	}
	return result
}

// GetTransitiveDependencies returns all functions that the given function depends on (transitively)
// direction: "callees" for what this function calls, "callers" for what calls this function
func (b *CallGraphBuilderImpl) GetTransitiveDependencies(functionID string, maxDepth int, direction string) []analysis.CallNode {
	deps := make(map[string]*analysis.CallNode)
	visited := make(map[string]bool)
	depths := make(map[string]int) // Track depth for each node

	var traverse func(id string, depth int)
	traverse = func(id string, depth int) {
		if depth > maxDepth || visited[id] {
			return
		}
		visited[id] = true

		node, ok := b.graph.Nodes[id]
		if !ok {
			return
		}

		var nextIDs []string
		if direction == "callers" {
			nextIDs = node.Callers
		} else {
			nextIDs = node.Callees
		}

		for _, nextID := range nextIDs {
			if next, ok := b.graph.Nodes[nextID]; ok {
				if _, exists := deps[nextID]; !exists {
					deps[nextID] = next
					depths[nextID] = depth + 1
				}
				traverse(nextID, depth+1)
			}
		}
	}

	traverse(functionID, 0)

	result := make([]analysis.CallNode, 0, len(deps))
	for _, node := range deps {
		result = append(result, *node)
	}

	// Sort by depth (closest first)
	sort.Slice(result, func(i, j int) bool {
		return depths[result[i].ID] < depths[result[j].ID]
	})

	return result
}

// GetTransitiveCallees returns all functions called by the given function (transitively)
func (b *CallGraphBuilderImpl) GetTransitiveCallees(functionID string, maxDepth int) []analysis.CallNode {
	return b.GetTransitiveDependencies(functionID, maxDepth, "callees")
}

// GetTransitiveCallers returns all functions that call the given function (transitively)
func (b *CallGraphBuilderImpl) GetTransitiveCallers(functionID string, maxDepth int) []analysis.CallNode {
	return b.GetTransitiveDependencies(functionID, maxDepth, "callers")
}

// GetDependencyPath finds the shortest path between two functions
func (b *CallGraphBuilderImpl) GetDependencyPath(fromID, toID string, maxDepth int) []string {
	if fromID == toID {
		return []string{fromID}
	}

	// BFS for shortest path
	type queueItem struct {
		id   string
		path []string
	}

	visited := make(map[string]bool)
	queue := []queueItem{{id: fromID, path: []string{fromID}}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if len(current.path) > maxDepth+1 {
			continue
		}

		if visited[current.id] {
			continue
		}
		visited[current.id] = true

		node, ok := b.graph.Nodes[current.id]
		if !ok {
			continue
		}

		for _, calleeID := range node.Callees {
			newPath := make([]string, len(current.path)+1)
			copy(newPath, current.path)
			newPath[len(current.path)] = calleeID

			if calleeID == toID {
				return newPath
			}

			queue = append(queue, queueItem{id: calleeID, path: newPath})
		}
	}

	return nil // No path found
}

// BuildForFile builds call graph for a single file
func (b *CallGraphBuilderImpl) BuildForFile(ctx context.Context, filePath string, content []byte) error {
	ext := filepath.Ext(filePath)

	switch ext {
	case extGo:
		b.analyzeGoFileContent(filePath, content)
	}

	return nil
}

func (b *CallGraphBuilderImpl) analyzeGoFileContent(relPath string, content []byte) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, relPath, content, parser.ParseComments)
	if err != nil {
		return
	}

	pkgName := ""
	if file.Name != nil {
		pkgName = file.Name.Name
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			nodeID := b.makeFunctionID(pkgName, decl.Name.Name, relPath)
			pos := fset.Position(decl.Pos())

			b.graph.Nodes[nodeID] = &analysis.CallNode{
				ID:       nodeID,
				Name:     decl.Name.Name,
				FilePath: relPath,
				Line:     pos.Line,
				Package:  pkgName,
				Callers:  make([]string, 0),
				Callees:  make([]string, 0),
			}
		}
		return true
	})
}

// analyzeVueFile analyzes Vue SFC files
func (b *CallGraphBuilderImpl) analyzeVueFile(path, relPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	text := string(content)

	// Extract script section
	scriptRe := regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	scriptMatch := scriptRe.FindStringSubmatch(text)
	if len(scriptMatch) < 2 {
		return
	}
	scriptContent := scriptMatch[1]

	// Find imports for dependency graph
	importRe := regexp.MustCompile(`import\s+(?:{[^}]+}|\w+|\*\s+as\s+\w+)\s+from\s+['"]([^'"]+)['"]`)
	importMatches := importRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range importMatches {
		if match[2] != -1 && match[3] != -1 {
			importPath := text[match[2]:match[3]]
			line := strings.Count(text[:match[0]], "\n") + 1
			b.fileImports[relPath] = append(b.fileImports[relPath], importInfo{path: importPath, line: line})
		}
	}

	// Find function definitions in script
	funcRe := regexp.MustCompile(`(?m)(?:function\s+(\w+)|(?:const|let|var)\s+(\w+)\s*=\s*(?:async\s*)?\([^)]*\)\s*=>|(\w+)\s*\([^)]*\)\s*\{)`)
	funcMatches := funcRe.FindAllStringSubmatchIndex(scriptContent, -1)

	for _, match := range funcMatches {
		var name string
		for i := 2; i < len(match); i += 2 {
			if match[i] != -1 {
				name = scriptContent[match[i]:match[i+1]]
				break
			}
		}
		if name == "" || isJSKeyword(name) {
			continue
		}

		line := strings.Count(text[:match[0]], "\n") + 1
		nodeID := b.makeFunctionID("", name, relPath)

		b.graph.Nodes[nodeID] = &analysis.CallNode{
			ID:       nodeID,
			Name:     name,
			FilePath: relPath,
			Line:     line,
			Callers:  make([]string, 0),
			Callees:  make([]string, 0),
		}
	}

	// Find function calls within script
	b.analyzeJSCalls(scriptContent, relPath)
}

// funcScope represents a function's scope in a file
type funcScope struct {
	nodeID    string
	startLine int
	endLine   int
}

// buildFuncScopes builds sorted function scopes for a file
func (b *CallGraphBuilderImpl) buildFuncScopes(relPath string) []funcScope {
	var funcsInFile []funcScope
	for nodeID, node := range b.graph.Nodes {
		if node.FilePath == relPath {
			funcsInFile = append(funcsInFile, funcScope{nodeID: nodeID, startLine: node.Line, endLine: node.Line + 100})
		}
	}
	sort.Slice(funcsInFile, func(i, j int) bool { return funcsInFile[i].startLine < funcsInFile[j].startLine })
	for i := 0; i < len(funcsInFile)-1; i++ {
		funcsInFile[i].endLine = funcsInFile[i+1].startLine - 1
	}
	return funcsInFile
}

// findContainingFunc finds which function contains a given line
func findContainingFunc(funcsInFile []funcScope, lineNum int) string {
	for _, f := range funcsInFile {
		if lineNum >= f.startLine && lineNum <= f.endLine {
			return f.nodeID
		}
	}
	return ""
}

// addCallEdge adds a call edge and updates caller/callee lists
func (b *CallGraphBuilderImpl) addCallEdge(callerID, calleeID, relPath string, line int) {
	b.graph.Edges = append(b.graph.Edges, analysis.CallEdge{From: callerID, To: calleeID, FilePath: relPath, Line: line, CallType: "direct"})
	if caller, ok := b.graph.Nodes[callerID]; ok {
		caller.Callees = append(caller.Callees, calleeID)
	}
	if callee, ok := b.graph.Nodes[calleeID]; ok {
		callee.Callers = append(callee.Callers, callerID)
	}
}

// analyzeJSCalls finds function calls in JS/TS content
func (b *CallGraphBuilderImpl) analyzeJSCalls(content, relPath string) {
	callRe := regexp.MustCompile(`\b(\w+)\s*\(`)
	lines := strings.Split(content, "\n")
	funcsInFile := b.buildFuncScopes(relPath)

	for lineNum, line := range lines {
		actualLine := lineNum + 1
		matches := callRe.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			calleeName := match[1]
			if isJSKeyword(calleeName) {
				continue
			}

			callerID := findContainingFunc(funcsInFile, actualLine)
			if callerID == "" {
				continue
			}

			calleeID := b.makeFunctionID("", calleeName, "")
			calleeIDWithFile := b.makeFunctionID("", calleeName, relPath)

			if _, exists := b.graph.Nodes[calleeID]; exists {
				b.addCallEdge(callerID, calleeID, relPath, actualLine)
			} else if _, exists := b.graph.Nodes[calleeIDWithFile]; exists && calleeIDWithFile != callerID {
				b.addCallEdge(callerID, calleeIDWithFile, relPath, actualLine)
			}
		}
	}
}

func isJSKeyword(name string) bool {
	keywords := map[string]bool{
		"if": true, "for": true, "while": true, "function": true,
		"switch": true, "catch": true, "return": true, "throw": true,
		"new": true, "typeof": true, "instanceof": true, "await": true,
		"import": true, "export": true, "class": true, "extends": true,
	}
	return keywords[name]
}

// collectImportsFromProject walks project and collects imports
func (b *CallGraphBuilderImpl) collectImportsFromProject(projectRoot string) error {
	return filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == dirVendor || name == dirNodeModules {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, _ := filepath.Rel(projectRoot, path)
		switch filepath.Ext(path) {
		case extGo:
			b.collectGoImports(path, relPath)
		case ".ts", ".tsx", ".js", ".jsx", ".vue":
			b.collectJSImports(path, relPath)
		}
		return nil
	})
}

// ensureDepNode ensures a dependency node exists
func (b *CallGraphBuilderImpl) ensureDepNode(filePath string) {
	if _, exists := b.depGraph.Nodes[filePath]; !exists {
		b.depGraph.Nodes[filePath] = &analysis.DependencyNode{
			ID: filePath, Name: filepath.Base(filePath), Type: "file",
			FilePath: filePath, Package: filepath.Dir(filePath),
			Dependencies: make([]string, 0), Dependents: make([]string, 0),
		}
	}
}

// buildDepEdges builds dependency edges from collected imports
func (b *CallGraphBuilderImpl) buildDepEdges(projectRoot string) {
	for filePath, imports := range b.fileImports {
		b.ensureDepNode(filePath)
		for _, imp := range imports {
			targetPath := b.resolveImportPath(filePath, imp.path, projectRoot)
			if targetPath == "" {
				continue
			}
			b.ensureDepNode(targetPath)
			b.depGraph.Edges = append(b.depGraph.Edges, analysis.DependencyEdge{
				From: filePath, To: targetPath, ImportPath: imp.path, Line: imp.line,
			})
			b.depGraph.Nodes[filePath].Dependencies = append(b.depGraph.Nodes[filePath].Dependencies, targetPath)
			b.depGraph.Nodes[targetPath].Dependents = append(b.depGraph.Nodes[targetPath].Dependents, filePath)
		}
	}
}

// BuildDependencyGraph builds file/package dependency graph
func (b *CallGraphBuilderImpl) BuildDependencyGraph(projectRoot string) (*analysis.DependencyGraph, error) {
	b.depGraph = &analysis.DependencyGraph{
		Nodes: make(map[string]*analysis.DependencyNode),
		Edges: make([]analysis.DependencyEdge, 0),
	}
	b.fileImports = make(map[string][]importInfo)

	if err := b.collectImportsFromProject(projectRoot); err != nil {
		return nil, err
	}
	b.buildDepEdges(projectRoot)

	return b.depGraph, nil
}

func (b *CallGraphBuilderImpl) collectGoImports(path, relPath string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
	if err != nil {
		return
	}

	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		pos := fset.Position(imp.Pos())
		b.fileImports[relPath] = append(b.fileImports[relPath], importInfo{
			path: importPath,
			line: pos.Line,
		})
	}
}

func (b *CallGraphBuilderImpl) collectJSImports(path, relPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	text := string(content)

	// Match various import patterns
	importRe := regexp.MustCompile(`(?m)import\s+(?:{[^}]+}|\w+|\*\s+as\s+\w+)?\s*(?:,\s*(?:{[^}]+}|\w+))?\s*from\s+['"]([^'"]+)['"]`)
	matches := importRe.FindAllStringSubmatchIndex(text, -1)

	for _, match := range matches {
		if len(match) >= 4 && match[2] != -1 {
			importPath := text[match[2]:match[3]]
			line := strings.Count(text[:match[0]], "\n") + 1
			b.fileImports[relPath] = append(b.fileImports[relPath], importInfo{
				path: importPath,
				line: line,
			})
		}
	}

	// Also match require() calls
	requireRe := regexp.MustCompile(`require\s*\(\s*['"]([^'"]+)['"]\s*\)`)
	reqMatches := requireRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range reqMatches {
		if len(match) >= 4 && match[2] != -1 {
			importPath := text[match[2]:match[3]]
			line := strings.Count(text[:match[0]], "\n") + 1
			b.fileImports[relPath] = append(b.fileImports[relPath], importInfo{
				path: importPath,
				line: line,
			})
		}
	}
}

func (b *CallGraphBuilderImpl) resolveImportPath(fromFile, importPath, projectRoot string) string {
	// Skip external packages
	if !strings.HasPrefix(importPath, ".") && !strings.HasPrefix(importPath, "@/") && !strings.HasPrefix(importPath, "~/") {
		return ""
	}

	fromDir := filepath.Dir(fromFile)

	// Handle alias imports (@/, ~/)
	if strings.HasPrefix(importPath, "@/") {
		importPath = strings.TrimPrefix(importPath, "@/")
		fromDir = "src" // Common convention
	} else if strings.HasPrefix(importPath, "~/") {
		importPath = strings.TrimPrefix(importPath, "~/")
		fromDir = ""
	}

	// Resolve relative path
	resolved := filepath.Join(fromDir, importPath)
	resolved = filepath.Clean(resolved)

	// Try common extensions
	extensions := []string{"", ".ts", ".tsx", ".js", ".jsx", ".vue", "/index.ts", "/index.js"}
	for _, ext := range extensions {
		candidate := resolved + ext
		fullPath := filepath.Join(projectRoot, candidate)
		if _, err := os.Stat(fullPath); err == nil {
			return candidate
		}
	}

	return ""
}

// cycleDFSState holds state for cycle detection DFS
type cycleDFSState struct {
	visited  map[string]bool
	recStack map[string]bool
	path     []string
	cycles   []analysis.CyclicDependency
}

// extractCycle extracts a cycle from the current path
func (s *cycleDFSState) extractCycle(dep string) {
	cycleStart := -1
	for i, p := range s.path {
		if p == dep {
			cycleStart = i
			break
		}
	}
	if cycleStart >= 0 {
		cycle := make([]string, len(s.path)-cycleStart+1)
		copy(cycle, s.path[cycleStart:])
		cycle[len(cycle)-1] = dep
		s.cycles = append(s.cycles, analysis.CyclicDependency{Cycle: cycle, Type: "file"})
	}
}

// dfsForCycles performs DFS to find cycles
func (b *CallGraphBuilderImpl) dfsForCycles(node string, state *cycleDFSState) bool {
	state.visited[node] = true
	state.recStack[node] = true
	state.path = append(state.path, node)

	if depNode, exists := b.depGraph.Nodes[node]; exists {
		for _, dep := range depNode.Dependencies {
			if !state.visited[dep] {
				if b.dfsForCycles(dep, state) {
					return true
				}
			} else if state.recStack[dep] {
				state.extractCycle(dep)
				return true
			}
		}
	}

	state.path = state.path[:len(state.path)-1]
	state.recStack[node] = false
	return false
}

// FindCyclicDependencies finds all cyclic dependencies in the project
func (b *CallGraphBuilderImpl) FindCyclicDependencies(projectRoot string) ([]analysis.CyclicDependency, error) {
	if _, err := b.BuildDependencyGraph(projectRoot); err != nil {
		return nil, err
	}

	state := &cycleDFSState{
		visited:  make(map[string]bool),
		recStack: make(map[string]bool),
		path:     make([]string, 0),
		cycles:   make([]analysis.CyclicDependency, 0),
	}

	for nodeID := range b.depGraph.Nodes {
		if !state.visited[nodeID] {
			b.dfsForCycles(nodeID, state)
		}
	}

	return state.cycles, nil
}

// GetFileDependencies returns files that a file depends on
func (b *CallGraphBuilderImpl) GetFileDependencies(filePath string) []analysis.DependencyNode {
	node, exists := b.depGraph.Nodes[filePath]
	if !exists {
		return nil
	}

	var deps []analysis.DependencyNode
	for _, depID := range node.Dependencies {
		if dep, ok := b.depGraph.Nodes[depID]; ok {
			deps = append(deps, *dep)
		}
	}
	return deps
}

// GetFileDependents returns files that depend on a file
func (b *CallGraphBuilderImpl) GetFileDependents(filePath string) []analysis.DependencyNode {
	node, exists := b.depGraph.Nodes[filePath]
	if !exists {
		return nil
	}

	var dependents []analysis.DependencyNode
	for _, depID := range node.Dependents {
		if dep, ok := b.depGraph.Nodes[depID]; ok {
			dependents = append(dependents, *dep)
		}
	}
	return dependents
}

// GetDependencyGraph returns the current dependency graph
func (b *CallGraphBuilderImpl) GetDependencyGraph() *analysis.DependencyGraph {
	return b.depGraph
}

// ExportMermaid exports the call graph as Mermaid diagram
func (b *CallGraphBuilderImpl) ExportMermaid(maxNodes int) string {
	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// Limit nodes for readability
	nodeCount := 0
	nodeIDs := make([]string, 0, len(b.graph.Nodes))
	for id := range b.graph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	nodeMap := make(map[string]string) // original ID -> safe ID
	for _, id := range nodeIDs {
		if nodeCount >= maxNodes {
			break
		}
		node := b.graph.Nodes[id]
		safeID := fmt.Sprintf("N%d", nodeCount)
		nodeMap[id] = safeID
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", safeID, node.Name))
		nodeCount++
	}

	// Add edges
	for _, edge := range b.graph.Edges {
		fromSafe, fromOK := nodeMap[edge.From]
		toSafe, toOK := nodeMap[edge.To]
		if fromOK && toOK {
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromSafe, toSafe))
		}
	}

	return sb.String()
}

// ExportDependencyMermaid exports the dependency graph as Mermaid diagram
func (b *CallGraphBuilderImpl) ExportDependencyMermaid(maxNodes int) string {
	var sb strings.Builder
	sb.WriteString("graph TD\n")

	nodeCount := 0
	nodeIDs := make([]string, 0, len(b.depGraph.Nodes))
	for id := range b.depGraph.Nodes {
		nodeIDs = append(nodeIDs, id)
	}
	sort.Strings(nodeIDs)

	nodeMap := make(map[string]string)
	for _, id := range nodeIDs {
		if nodeCount >= maxNodes {
			break
		}
		node := b.depGraph.Nodes[id]
		safeID := fmt.Sprintf("F%d", nodeCount)
		nodeMap[id] = safeID
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", safeID, node.Name))
		nodeCount++
	}

	for _, edge := range b.depGraph.Edges {
		fromSafe, fromOK := nodeMap[edge.From]
		toSafe, toOK := nodeMap[edge.To]
		if fromOK && toOK {
			sb.WriteString(fmt.Sprintf("    %s --> %s\n", fromSafe, toSafe))
		}
	}

	return sb.String()
}
