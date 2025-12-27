package analyzers

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"syntaxia/domain/analysis"
)

// TSCallGraphBuilder builds call graphs for TypeScript/JavaScript/Vue projects
type TSCallGraphBuilder struct {
	mu              sync.RWMutex
	graph           *analysis.CallGraph
	componentGraph  *VueComponentGraph
	fileImports     map[string][]tsImportInfo
	componentUsages map[string][]ComponentUsage // file -> component usages
}

// tsImportInfo holds import information for TS/JS files
type tsImportInfo struct {
	path       string
	line       int
	names      []string // imported names
	isDefault  bool
	isNamespace bool
}

// VueComponentGraph represents Vue component usage relationships
type VueComponentGraph struct {
	Components map[string]*VueComponentNode `json:"components"`
	Usages     []ComponentUsageEdge         `json:"usages"`
}

// VueComponentNode represents a Vue component in the graph
type VueComponentNode struct {
	ID        string   `json:"id"`        // unique identifier (file path)
	Name      string   `json:"name"`      // component name
	FilePath  string   `json:"filePath"`  // file where defined
	Props     []string `json:"props"`     // component props
	Emits     []string `json:"emits"`     // emitted events
	UsedBy    []string `json:"usedBy"`    // components that use this
	Uses      []string `json:"uses"`      // components this uses
}

// ComponentUsage represents where a component is used
type ComponentUsage struct {
	ComponentName string `json:"componentName"`
	Line          int    `json:"line"`
	InTemplate    bool   `json:"inTemplate"`
}

// ComponentUsageEdge represents a usage relationship
type ComponentUsageEdge struct {
	From     string `json:"from"`     // file using the component
	To       string `json:"to"`       // component being used
	Line     int    `json:"line"`     // line of usage
	UsageType string `json:"usageType"` // "template", "script", "dynamic"
}

// NewTSCallGraphBuilder creates a new TypeScript call graph builder
func NewTSCallGraphBuilder() *TSCallGraphBuilder {
	return &TSCallGraphBuilder{
		graph: &analysis.CallGraph{
			Nodes: make(map[string]*analysis.CallNode),
			Edges: make([]analysis.CallEdge, 0),
		},
		componentGraph: &VueComponentGraph{
			Components: make(map[string]*VueComponentNode),
			Usages:     make([]ComponentUsageEdge, 0),
		},
		fileImports:     make(map[string][]tsImportInfo),
		componentUsages: make(map[string][]ComponentUsage),
	}
}

// Build builds call graph for TypeScript/JavaScript/Vue project
func (b *TSCallGraphBuilder) Build(projectRoot string) (*analysis.CallGraph, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.reset()

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return b.shouldSkipDir(info.Name())
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		ext := filepath.Ext(path)

		switch ext {
		case ".ts", ".tsx":
			b.analyzeTypeScriptFile(path, relPath)
		case ".js", ".jsx":
			b.analyzeJavaScriptFile(path, relPath)
		case ".vue":
			b.analyzeVueFile(path, relPath)
		}

		return nil
	})

	return b.graph, err
}

// reset clears all internal state
func (b *TSCallGraphBuilder) reset() {
	b.graph = &analysis.CallGraph{
		Nodes: make(map[string]*analysis.CallNode),
		Edges: make([]analysis.CallEdge, 0),
	}
	b.componentGraph = &VueComponentGraph{
		Components: make(map[string]*VueComponentNode),
		Usages:     make([]ComponentUsageEdge, 0),
	}
	b.fileImports = make(map[string][]tsImportInfo)
	b.componentUsages = make(map[string][]ComponentUsage)
}

// shouldSkipDir returns filepath.SkipDir for directories to skip
func (b *TSCallGraphBuilder) shouldSkipDir(name string) error {
	skipDirs := map[string]bool{
		"node_modules": true, "vendor": true, "dist": true,
		"build": true, ".git": true, ".next": true, ".nuxt": true,
	}
	if strings.HasPrefix(name, ".") || skipDirs[name] {
		return filepath.SkipDir
	}
	return nil
}

// Precompiled regex patterns for TypeScript analysis
var (
	// Function declarations
	tsFuncDeclRe = regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(async\s+)?function\s+(\w+)\s*(<[^>]*>)?\s*\(`)
	// Arrow functions assigned to const/let/var
	tsArrowFuncRe = regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(const|let|var)\s+(\w+)\s*(?::\s*[^=]+)?\s*=\s*(async\s*)?\([^)]*\)\s*(?::\s*[^=]+)?\s*=>`)
	// Class declarations
	tsClassRe = regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(abstract\s+)?class\s+(\w+)`)
	// Method declarations in class
	tsMethodRe = regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|static|async|\s)*(\w+)\s*(<[^>]*>)?\s*\([^)]*\)\s*(?::\s*[^{]+)?\s*\{`)
	// Function calls
	tsFuncCallRe = regexp.MustCompile(`\b(\w+)\s*\(`)
	// Method calls (obj.method())
	tsMethodCallRe = regexp.MustCompile(`\b(\w+)\.(\w+)\s*\(`)
	// Import statements
	tsImportRe = regexp.MustCompile(`(?m)^import\s+(?:(\w+)\s*,?\s*)?(?:\{([^}]+)\})?\s*(?:\*\s+as\s+(\w+))?\s*from\s+['"]([^'"]+)['"]`)
)

// analyzeTypeScriptFile analyzes a TypeScript file for functions and calls
func (b *TSCallGraphBuilder) analyzeTypeScriptFile(path, relPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	text := string(content)
	strippedText := stripCommentsAndStrings(text)

	// Collect imports
	b.collectTSImports(text, relPath)

	// Collect function declarations
	b.collectTSFunctions(strippedText, text, relPath)

	// Collect class methods
	b.collectTSClassMethods(strippedText, text, relPath)

	// Analyze function calls
	b.analyzeTSCalls(strippedText, relPath)
}

// analyzeJavaScriptFile analyzes a JavaScript file (reuses TS logic)
func (b *TSCallGraphBuilder) analyzeJavaScriptFile(path, relPath string) {
	b.analyzeTypeScriptFile(path, relPath)
}

// collectTSImports collects import statements from a file
func (b *TSCallGraphBuilder) collectTSImports(text, relPath string) {
	matches := tsImportRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		imp := tsImportInfo{
			path: match[4],
			line: strings.Count(text[:strings.Index(text, match[0])], "\n") + 1,
		}

		// Default import
		if match[1] != "" {
			imp.names = append(imp.names, match[1])
			imp.isDefault = true
		}

		// Named imports
		if match[2] != "" {
			names := strings.Split(match[2], ",")
			for _, n := range names {
				n = strings.TrimSpace(n)
				if idx := strings.Index(n, " as "); idx != -1 {
					n = strings.TrimSpace(n[:idx])
				}
				if n != "" {
					imp.names = append(imp.names, n)
				}
			}
		}

		// Namespace import
		if match[3] != "" {
			imp.names = append(imp.names, match[3])
			imp.isNamespace = true
		}

		b.fileImports[relPath] = append(b.fileImports[relPath], imp)
	}
}

// collectTSFunctions collects function declarations from TypeScript code
func (b *TSCallGraphBuilder) collectTSFunctions(strippedText, originalText, relPath string) {
	lines := strings.Split(originalText, "\n")

	// Regular function declarations
	for _, match := range tsFuncDeclRe.FindAllStringSubmatchIndex(strippedText, -1) {
		name := strippedText[match[6]:match[7]]
		line := strings.Count(strippedText[:match[0]], "\n") + 1
		nodeID := b.makeFunctionID(name, relPath)

		b.graph.Nodes[nodeID] = &analysis.CallNode{
			ID:        nodeID,
			Name:      name,
			FilePath:  relPath,
			Line:      line,
			Signature: b.extractSignature(lines, line-1),
			Callers:   make([]string, 0),
			Callees:   make([]string, 0),
		}
	}

	// Arrow functions
	for _, match := range tsArrowFuncRe.FindAllStringSubmatchIndex(strippedText, -1) {
		name := strippedText[match[6]:match[7]]
		line := strings.Count(strippedText[:match[0]], "\n") + 1
		nodeID := b.makeFunctionID(name, relPath)

		b.graph.Nodes[nodeID] = &analysis.CallNode{
			ID:        nodeID,
			Name:      name,
			FilePath:  relPath,
			Line:      line,
			Signature: b.extractSignature(lines, line-1),
			Callers:   make([]string, 0),
			Callees:   make([]string, 0),
		}
	}
}

// collectTSClassMethods collects class method declarations
func (b *TSCallGraphBuilder) collectTSClassMethods(strippedText, originalText, relPath string) {
	lines := strings.Split(originalText, "\n")

	// Find classes first
	classMatches := tsClassRe.FindAllStringSubmatchIndex(strippedText, -1)
	for _, classMatch := range classMatches {
		className := strippedText[classMatch[6]:classMatch[7]]
		classLine := strings.Count(strippedText[:classMatch[0]], "\n") + 1
		classEndLine := findBlockEndLine(lines, classLine-1)

		// Extract class body
		classBody := extractClassBody(lines, classLine-1, classEndLine)

		// Find methods in class body
		methodMatches := tsMethodRe.FindAllStringSubmatchIndex(classBody, -1)
		for _, methodMatch := range methodMatches {
			methodName := classBody[methodMatch[4]:methodMatch[5]]
			if isTSKeyword(methodName) {
				continue
			}

			methodLine := classLine + strings.Count(classBody[:methodMatch[0]], "\n")
			nodeID := b.makeMethodID(className, methodName, relPath)

			b.graph.Nodes[nodeID] = &analysis.CallNode{
				ID:        nodeID,
				Name:      methodName,
				FilePath:  relPath,
				Line:      methodLine,
				Package:   className,
				Signature: b.extractSignature(lines, methodLine-1),
				Callers:   make([]string, 0),
				Callees:   make([]string, 0),
			}
		}
	}
}

// extractClassBody extracts the body of a class between start and end lines
func extractClassBody(lines []string, startIdx, endIdx int) string {
	if startIdx < 0 || endIdx > len(lines) {
		return ""
	}
	return strings.Join(lines[startIdx:endIdx], "\n")
}

// analyzeTSCalls analyzes function and method calls in TypeScript code
func (b *TSCallGraphBuilder) analyzeTSCalls(strippedText, relPath string) {
	lines := strings.Split(strippedText, "\n")
	funcsInFile := b.buildTSFuncScopes(relPath)

	for lineNum, line := range lines {
		actualLine := lineNum + 1
		callerID := b.findContainingFunc(funcsInFile, actualLine)
		if callerID == "" {
			continue
		}

		// Regular function calls
		funcCalls := tsFuncCallRe.FindAllStringSubmatch(line, -1)
		for _, match := range funcCalls {
			calleeName := match[1]
			if isTSKeyword(calleeName) {
				continue
			}
			b.addTSCallEdge(callerID, calleeName, relPath, actualLine, "direct")
		}

		// Method calls
		methodCalls := tsMethodCallRe.FindAllStringSubmatch(line, -1)
		for _, match := range methodCalls {
			objName := match[1]
			methodName := match[2]
			if isTSKeyword(methodName) {
				continue
			}
			calleeID := objName + "." + methodName
			b.addTSCallEdge(callerID, calleeID, relPath, actualLine, "method")
		}
	}
}

// tsFuncScope represents a function's scope in a TypeScript file
type tsFuncScope struct {
	nodeID    string
	startLine int
	endLine   int
}

// buildTSFuncScopes builds sorted function scopes for a file
func (b *TSCallGraphBuilder) buildTSFuncScopes(relPath string) []tsFuncScope {
	var funcs []tsFuncScope
	for nodeID, node := range b.graph.Nodes {
		if node.FilePath == relPath {
			funcs = append(funcs, tsFuncScope{
				nodeID:    nodeID,
				startLine: node.Line,
				endLine:   node.Line + 100, // Will be adjusted
			})
		}
	}

	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].startLine < funcs[j].startLine
	})

	// Adjust end lines based on next function start
	for i := 0; i < len(funcs)-1; i++ {
		funcs[i].endLine = funcs[i+1].startLine - 1
	}

	return funcs
}

// findContainingFunc finds which function contains a given line
func (b *TSCallGraphBuilder) findContainingFunc(funcs []tsFuncScope, lineNum int) string {
	for _, f := range funcs {
		if lineNum >= f.startLine && lineNum <= f.endLine {
			return f.nodeID
		}
	}
	return ""
}

// addTSCallEdge adds a call edge and updates caller/callee lists
func (b *TSCallGraphBuilder) addTSCallEdge(callerID, calleeName, relPath string, line int, callType string) {
	// Try to find callee in graph
	calleeID := b.findCalleeID(calleeName, relPath)
	if calleeID == "" {
		return
	}

	b.graph.Edges = append(b.graph.Edges, analysis.CallEdge{
		From:     callerID,
		To:       calleeID,
		FilePath: relPath,
		Line:     line,
		CallType: callType,
	})

	if caller, ok := b.graph.Nodes[callerID]; ok {
		caller.Callees = append(caller.Callees, calleeID)
	}
	if callee, ok := b.graph.Nodes[calleeID]; ok {
		callee.Callers = append(callee.Callers, callerID)
	}
}

// findCalleeID finds the callee node ID in the graph
func (b *TSCallGraphBuilder) findCalleeID(calleeName, relPath string) string {
	// Try exact match with file
	fileID := b.makeFunctionID(calleeName, relPath)
	if _, exists := b.graph.Nodes[fileID]; exists {
		return fileID
	}

	// Try without file (global/imported)
	globalID := b.makeFunctionID(calleeName, "")
	if _, exists := b.graph.Nodes[globalID]; exists {
		return globalID
	}

	// Search in all nodes
	for nodeID, node := range b.graph.Nodes {
		if node.Name == calleeName {
			return nodeID
		}
	}

	return ""
}

// Vue component analysis patterns
var (
	vueScriptRe       = regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`)
	vueTemplateRe     = regexp.MustCompile(`(?s)<template[^>]*>(.*?)</template>`)
	vueComponentTagRe = regexp.MustCompile(`<([A-Z][a-zA-Z0-9]*|[a-z]+-[a-z-]+)`)
	vueDefinePropsRe  = regexp.MustCompile(`defineProps\s*<?\s*\{([^}]+)\}`)
	vueDefineEmitsRe  = regexp.MustCompile(`defineEmits\s*<?\s*\[([^\]]+)\]`)
	vueComponentsRe   = regexp.MustCompile(`components\s*:\s*\{([^}]+)\}`)
)

// analyzeVueFile analyzes a Vue SFC file
func (b *TSCallGraphBuilder) analyzeVueFile(path, relPath string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	text := string(content)

	// Extract component name from filename
	componentName := extractComponentName(relPath)

	// Create component node
	b.componentGraph.Components[relPath] = &VueComponentNode{
		ID:       relPath,
		Name:     componentName,
		FilePath: relPath,
		Props:    make([]string, 0),
		Emits:    make([]string, 0),
		UsedBy:   make([]string, 0),
		Uses:     make([]string, 0),
	}

	// Analyze script section
	if scriptMatch := vueScriptRe.FindStringSubmatch(text); len(scriptMatch) > 1 {
		scriptContent := scriptMatch[1]
		strippedScript := stripCommentsAndStrings(scriptContent)

		// Collect imports
		b.collectTSImports(scriptContent, relPath)

		// Collect functions
		b.collectTSFunctions(strippedScript, scriptContent, relPath)

		// Analyze calls
		b.analyzeTSCalls(strippedScript, relPath)

		// Extract props and emits
		b.extractVuePropsEmits(scriptContent, relPath)
	}

	// Analyze template for component usage
	if templateMatch := vueTemplateRe.FindStringSubmatch(text); len(templateMatch) > 1 {
		b.analyzeVueTemplate(templateMatch[1], relPath)
	}
}

// extractComponentName extracts Vue component name from file path
func extractComponentName(filePath string) string {
	base := filepath.Base(filePath)
	return strings.TrimSuffix(base, ".vue")
}

// extractVuePropsEmits extracts props and emits from Vue script
func (b *TSCallGraphBuilder) extractVuePropsEmits(script, relPath string) {
	component, exists := b.componentGraph.Components[relPath]
	if !exists {
		return
	}

	// Extract props
	if propsMatch := vueDefinePropsRe.FindStringSubmatch(script); len(propsMatch) > 1 {
		props := extractPropNames(propsMatch[1])
		component.Props = append(component.Props, props...)
	}

	// Extract emits
	if emitsMatch := vueDefineEmitsRe.FindStringSubmatch(script); len(emitsMatch) > 1 {
		emits := extractEmitNames(emitsMatch[1])
		component.Emits = append(component.Emits, emits...)
	}
}

// extractPropNames extracts prop names from defineProps content
func extractPropNames(content string) []string {
	propRe := regexp.MustCompile(`(\w+)\s*[?:]`)
	matches := propRe.FindAllStringSubmatch(content, -1)
	props := make([]string, 0, len(matches))
	for _, m := range matches {
		props = append(props, m[1])
	}
	return props
}

// extractEmitNames extracts emit names from defineEmits content
func extractEmitNames(content string) []string {
	emitRe := regexp.MustCompile(`['"]([^'"]+)['"]`)
	matches := emitRe.FindAllStringSubmatch(content, -1)
	emits := make([]string, 0, len(matches))
	for _, m := range matches {
		emits = append(emits, m[1])
	}
	return emits
}

// analyzeVueTemplate analyzes Vue template for component usage
func (b *TSCallGraphBuilder) analyzeVueTemplate(template, relPath string) {
	lines := strings.Split(template, "\n")

	for lineNum, line := range lines {
		matches := vueComponentTagRe.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			componentTag := match[1]

			// Skip HTML elements
			if isHTMLElement(componentTag) {
				continue
			}

			// Convert kebab-case to PascalCase if needed
			componentName := toPascalCase(componentTag)

			usage := ComponentUsage{
				ComponentName: componentName,
				Line:          lineNum + 1,
				InTemplate:    true,
			}
			b.componentUsages[relPath] = append(b.componentUsages[relPath], usage)

			// Find the component file and create edge
			componentFile := b.findComponentFile(componentName)
			if componentFile != "" {
				b.componentGraph.Usages = append(b.componentGraph.Usages, ComponentUsageEdge{
					From:      relPath,
					To:        componentFile,
					Line:      lineNum + 1,
					UsageType: "template",
				})

				// Update component relationships
				if comp, exists := b.componentGraph.Components[componentFile]; exists {
					comp.UsedBy = append(comp.UsedBy, relPath)
				}
				if comp, exists := b.componentGraph.Components[relPath]; exists {
					comp.Uses = append(comp.Uses, componentFile)
				}
			}
		}
	}
}

// findComponentFile finds the file path for a component name
func (b *TSCallGraphBuilder) findComponentFile(componentName string) string {
	for filePath, comp := range b.componentGraph.Components {
		if comp.Name == componentName {
			return filePath
		}
	}
	return ""
}

// Helper functions

// makeFunctionID creates a unique function ID
func (b *TSCallGraphBuilder) makeFunctionID(name, filePath string) string {
	if filePath != "" {
		return filePath + ":" + name
	}
	return name
}

// makeMethodID creates a unique method ID
func (b *TSCallGraphBuilder) makeMethodID(className, methodName, filePath string) string {
	return filePath + ":" + className + "." + methodName
}

// extractSignature extracts function signature from a line
func (b *TSCallGraphBuilder) extractSignature(lines []string, lineIdx int) string {
	if lineIdx < 0 || lineIdx >= len(lines) {
		return ""
	}
	line := strings.TrimSpace(lines[lineIdx])
	// Truncate at opening brace
	if idx := strings.Index(line, "{"); idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}
	return line
}

// isTSKeyword checks if a name is a TypeScript/JavaScript keyword
func isTSKeyword(name string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true,
		"do": true, "switch": true, "case": true, "break": true,
		"continue": true, "return": true, "throw": true, "try": true,
		"catch": true, "finally": true, "new": true, "delete": true,
		"typeof": true, "instanceof": true, "void": true, "in": true,
		"function": true, "class": true, "extends": true, "super": true,
		"this": true, "import": true, "export": true, "default": true,
		"const": true, "let": true, "var": true, "async": true,
		"await": true, "yield": true, "static": true, "get": true,
		"set": true, "constructor": true, "true": true, "false": true,
		"null": true, "undefined": true, "NaN": true, "Infinity": true,
	}
	return keywords[name]
}

// isHTMLElement checks if a tag is a standard HTML element
func isHTMLElement(tag string) bool {
	htmlElements := map[string]bool{
		"div": true, "span": true, "p": true, "a": true, "button": true,
		"input": true, "form": true, "label": true, "select": true,
		"option": true, "textarea": true, "table": true, "tr": true,
		"td": true, "th": true, "thead": true, "tbody": true, "ul": true,
		"ol": true, "li": true, "h1": true, "h2": true, "h3": true,
		"h4": true, "h5": true, "h6": true, "img": true, "video": true,
		"audio": true, "canvas": true, "svg": true, "path": true,
		"header": true, "footer": true, "nav": true, "main": true,
		"section": true, "article": true, "aside": true, "template": true,
		"slot": true, "component": true, "transition": true, "keep-alive": true,
		"teleport": true, "suspense": true,
	}
	return htmlElements[strings.ToLower(tag)]
}

// toPascalCase converts kebab-case to PascalCase
func toPascalCase(s string) string {
	// Already PascalCase
	if len(s) > 0 && s[0] >= 'A' && s[0] <= 'Z' {
		return s
	}

	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

// Public API methods

// GetCallGraph returns the built call graph
func (b *TSCallGraphBuilder) GetCallGraph() *analysis.CallGraph {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.graph
}

// GetComponentGraph returns the Vue component usage graph
func (b *TSCallGraphBuilder) GetComponentGraph() *VueComponentGraph {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.componentGraph
}

// GetCallers returns functions that call the given function
func (b *TSCallGraphBuilder) GetCallers(functionID string) []analysis.CallNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	node, ok := b.graph.Nodes[functionID]
	if !ok {
		return nil
	}

	callers := make([]analysis.CallNode, 0, len(node.Callers))
	for _, callerID := range node.Callers {
		if caller, ok := b.graph.Nodes[callerID]; ok {
			callers = append(callers, *caller)
		}
	}
	return callers
}

// GetCallees returns functions called by the given function
func (b *TSCallGraphBuilder) GetCallees(functionID string) []analysis.CallNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	node, ok := b.graph.Nodes[functionID]
	if !ok {
		return nil
	}

	callees := make([]analysis.CallNode, 0, len(node.Callees))
	for _, calleeID := range node.Callees {
		if callee, ok := b.graph.Nodes[calleeID]; ok {
			callees = append(callees, *callee)
		}
	}
	return callees
}

// GetComponentUsages returns components used by a file
func (b *TSCallGraphBuilder) GetComponentUsages(filePath string) []ComponentUsage {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.componentUsages[filePath]
}

// GetComponentDependents returns files that use a component
func (b *TSCallGraphBuilder) GetComponentDependents(componentPath string) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	comp, exists := b.componentGraph.Components[componentPath]
	if !exists {
		return nil
	}
	return comp.UsedBy
}

// GetComponentDependencies returns components used by a component
func (b *TSCallGraphBuilder) GetComponentDependencies(componentPath string) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	comp, exists := b.componentGraph.Components[componentPath]
	if !exists {
		return nil
	}
	return comp.Uses
}

// GetImpact returns all functions affected if given function changes
func (b *TSCallGraphBuilder) GetImpact(functionID string, maxDepth int) []analysis.CallNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

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

// GetCallersEdges returns call edges where the given function is called.
func (b *TSCallGraphBuilder) GetCallersEdges(symbol string) []analysis.CallEdge {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var edges []analysis.CallEdge
	for _, edge := range b.graph.Edges {
		if edge.To == symbol || strings.HasSuffix(edge.To, ":"+symbol) {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetCalleesEdges returns call edges where the given function calls others.
func (b *TSCallGraphBuilder) GetCalleesEdges(symbol string) []analysis.CallEdge {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var edges []analysis.CallEdge
	for _, edge := range b.graph.Edges {
		if edge.From == symbol || strings.HasSuffix(edge.From, ":"+symbol) {
			edges = append(edges, edge)
		}
	}
	return edges
}

// BuildComponentUsageGraph builds a graph of Vue component usage relationships.
func (b *TSCallGraphBuilder) BuildComponentUsageGraph(projectRoot string) (*VueComponentGraph, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Reset component graph
	b.componentGraph = &VueComponentGraph{
		Components: make(map[string]*VueComponentNode),
		Usages:     make([]ComponentUsageEdge, 0),
	}

	// Walk project and analyze Vue files
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return b.shouldSkipDir(info.Name())
		}

		if filepath.Ext(path) != ".vue" {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		b.analyzeVueFile(path, relPath)
		return nil
	})

	return b.componentGraph, err
}

// GetComponentUsagesList returns list of component names used in a file.
func (b *TSCallGraphBuilder) GetComponentUsagesList(componentName string) []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var files []string
	for filePath, usages := range b.componentUsages {
		for _, usage := range usages {
			if usage.ComponentName == componentName {
				files = append(files, filePath)
				break
			}
		}
	}
	return files
}

// GetCallChain finds a call chain from start to end function.
func (b *TSCallGraphBuilder) GetCallChain(startID, endID string, maxDepth int) [][]string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var chains [][]string
	visited := make(map[string]bool)

	var dfs func(current string, path []string, depth int)
	dfs = func(current string, path []string, depth int) {
		if depth > maxDepth {
			return
		}
		if current == endID {
			chain := make([]string, len(path))
			copy(chain, path)
			chains = append(chains, chain)
			return
		}
		if visited[current] {
			return
		}
		visited[current] = true
		defer func() { visited[current] = false }()

		node, ok := b.graph.Nodes[current]
		if !ok {
			return
		}

		for _, calleeID := range node.Callees {
			newPath := append(path, calleeID)
			dfs(calleeID, newPath, depth+1)
		}
	}

	dfs(startID, []string{startID}, 0)
	return chains
}

// FindFunctionByName finds a function node by name (partial match).
func (b *TSCallGraphBuilder) FindFunctionByName(name string) []*analysis.CallNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	nameLower := strings.ToLower(name)
	var results []*analysis.CallNode

	for _, node := range b.graph.Nodes {
		if strings.Contains(strings.ToLower(node.Name), nameLower) {
			results = append(results, node)
		}
	}
	return results
}

// GetFunctionsInFile returns all functions defined in a file.
func (b *TSCallGraphBuilder) GetFunctionsInFile(filePath string) []*analysis.CallNode {
	b.mu.RLock()
	defer b.mu.RUnlock()

	var results []*analysis.CallNode
	for _, node := range b.graph.Nodes {
		if node.FilePath == filePath {
			results = append(results, node)
		}
	}
	return results
}

// Stats returns statistics about the call graph.
func (b *TSCallGraphBuilder) Stats() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()

	stats := map[string]int{
		"total_nodes":      len(b.graph.Nodes),
		"total_edges":      len(b.graph.Edges),
		"total_components": len(b.componentGraph.Components),
		"component_usages": len(b.componentGraph.Usages),
		"files_with_imports": len(b.fileImports),
	}

	// Count files
	files := make(map[string]bool)
	for _, node := range b.graph.Nodes {
		files[node.FilePath] = true
	}
	stats["files_analyzed"] = len(files)

	return stats
}
