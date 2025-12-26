package symbolgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// TSImportGraphBuilder builds import graphs for TypeScript/JavaScript/Vue projects
type TSImportGraphBuilder struct {
	log      domain.Logger
	registry analysis.AnalyzerRegistry
}

// NewTSImportGraphBuilder creates a new TypeScript import graph builder
func NewTSImportGraphBuilder(log domain.Logger, registry analysis.AnalyzerRegistry) domain.ImportGraphBuilder {
	return &TSImportGraphBuilder{
		log:      log,
		registry: registry,
	}
}

// tsConfigPaths represents the paths configuration from tsconfig.json
type tsConfigPaths struct {
	CompilerOptions struct {
		BaseURL string              `json:"baseUrl"`
		Paths   map[string][]string `json:"paths"`
	} `json:"compilerOptions"`
}

// BuildImportGraph builds the import graph for a TypeScript/JavaScript/Vue project
func (b *TSImportGraphBuilder) BuildImportGraph(ctx context.Context, projectRoot string) (*domain.ImportGraph, error) {
	b.log.Info(fmt.Sprintf("Building import graph for TS/JS/Vue project: %s", projectRoot))

	graph := &domain.ImportGraph{
		Packages: make(map[string]*domain.PackageNode),
		Imports:  make([]*domain.ImportEdge, 0),
	}

	// Load path aliases from tsconfig.json
	pathAliases := b.loadPathAliases(projectRoot)

	// Find all supported files
	files, err := b.findSourceFiles(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to find source files: %w", err)
	}

	b.log.Info(fmt.Sprintf("Found %d source files to analyze", len(files)))

	// Process each file
	for _, filePath := range files {
		if err := b.processFile(ctx, projectRoot, filePath, graph, pathAliases); err != nil {
			b.log.Warning(fmt.Sprintf("Failed to process %s: %v", filePath, err))
			continue
		}
	}

	// Sort imports for determinism
	b.sortGraphForDeterminism(graph)

	b.log.Info(fmt.Sprintf("Built import graph with %d packages and %d imports",
		len(graph.Packages), len(graph.Imports)))

	return graph, nil
}

// loadPathAliases loads path aliases from tsconfig.json or jsconfig.json
func (b *TSImportGraphBuilder) loadPathAliases(projectRoot string) map[string]string {
	aliases := make(map[string]string)

	// Try tsconfig.json first, then jsconfig.json
	configFiles := []string{"tsconfig.json", "jsconfig.json"}

	for _, configFile := range configFiles {
		configPath := filepath.Join(projectRoot, configFile)
		data, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}

		var config tsConfigPaths
		if err := json.Unmarshal(data, &config); err != nil {
			b.log.Warning(fmt.Sprintf("Failed to parse %s: %v", configFile, err))
			continue
		}

		baseURL := config.CompilerOptions.BaseURL
		if baseURL == "" {
			baseURL = "."
		}

		// Process path mappings
		for alias, paths := range config.CompilerOptions.Paths {
			if len(paths) == 0 {
				continue
			}

			// Remove wildcard from alias: "@/*" -> "@/"
			aliasPrefix := strings.TrimSuffix(alias, "*")

			// Get the first path and remove wildcard
			targetPath := strings.TrimSuffix(paths[0], "*")

			// Combine with baseURL
			fullPath := filepath.Join(baseURL, targetPath)

			aliases[aliasPrefix] = fullPath
		}

		b.log.Info(fmt.Sprintf("Loaded %d path aliases from %s", len(aliases), configFile))
		break
	}

	return aliases
}

// findSourceFiles finds all TypeScript, JavaScript, and Vue files in the project
func (b *TSImportGraphBuilder) findSourceFiles(projectRoot string) ([]string, error) {
	var files []string

	supportedExts := map[string]bool{
		".ts":  true,
		".tsx": true,
		".js":  true,
		".jsx": true,
		".mjs": true,
		".vue": true,
	}

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Skip common non-source directories
			name := info.Name()
			if name == "node_modules" || name == "dist" || name == "build" ||
				name == ".git" || name == "coverage" || name == ".next" ||
				name == ".nuxt" || name == "vendor" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if supportedExts[ext] {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// processFile processes a single file and adds its imports to the graph
func (b *TSImportGraphBuilder) processFile(
	ctx context.Context,
	projectRoot string,
	filePath string,
	graph *domain.ImportGraph,
	pathAliases map[string]string,
) error {
	// Get relative path for the file
	relPath, err := filepath.Rel(projectRoot, filePath)
	if err != nil {
		relPath = filePath
	}
	relPath = filepath.ToSlash(relPath)

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get appropriate analyzer
	analyzer := b.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return fmt.Errorf("no analyzer found for %s", filePath)
	}

	// Extract imports
	imports, err := analyzer.GetImports(ctx, filePath, content)
	if err != nil {
		return fmt.Errorf("failed to extract imports: %w", err)
	}

	// Extract exports
	exports, err := analyzer.GetExports(ctx, filePath, content)
	if err != nil {
		// Non-critical, continue without exports
		exports = nil
	}

	// Create or update package node for this file
	packageNode := b.getOrCreatePackageNode(graph, relPath)

	// Add exports to package node
	for _, exp := range exports {
		packageNode.Exports = append(packageNode.Exports, exp.Name)
	}

	// Process imports
	for _, imp := range imports {
		resolvedPath := b.resolveImportPath(projectRoot, filePath, imp.Path, pathAliases)

		// Add to package imports
		packageNode.Imports = append(packageNode.Imports, resolvedPath)

		// Create import edge
		edgeType := "direct"
		if !imp.IsLocal {
			edgeType = "external"
		}

		graph.Imports = append(graph.Imports, &domain.ImportEdge{
			From: relPath,
			To:   resolvedPath,
			Type: edgeType,
		})

		// Create package node for the imported file if it's local
		if imp.IsLocal || strings.HasPrefix(imp.Path, "@/") || strings.HasPrefix(imp.Path, "~/") {
			b.getOrCreatePackageNode(graph, resolvedPath)
		}
	}

	return nil
}

// getOrCreatePackageNode gets or creates a package node for a file path
func (b *TSImportGraphBuilder) getOrCreatePackageNode(graph *domain.ImportGraph, filePath string) *domain.PackageNode {
	if node, exists := graph.Packages[filePath]; exists {
		return node
	}

	node := &domain.PackageNode{
		Name:    filepath.Base(filePath),
		Path:    filePath,
		Files:   []string{filePath},
		Imports: make([]string, 0),
		Exports: make([]string, 0),
	}
	graph.Packages[filePath] = node
	return node
}

// resolveImportPath resolves an import path to an actual file path
func (b *TSImportGraphBuilder) resolveImportPath(
	projectRoot string,
	fromFile string,
	importPath string,
	pathAliases map[string]string,
) string {
	// Handle path aliases first
	for alias, target := range pathAliases {
		if strings.HasPrefix(importPath, alias) {
			importPath = target + strings.TrimPrefix(importPath, alias)
			break
		}
	}

	// Handle relative imports
	if strings.HasPrefix(importPath, ".") {
		// Get relative path of fromFile from projectRoot
		relFromFile, err := filepath.Rel(projectRoot, fromFile)
		if err != nil {
			relFromFile = fromFile
		}
		relFromFile = filepath.ToSlash(relFromFile)

		// Get directory of the source file
		fromDir := filepath.Dir(relFromFile)

		// Join with import path and clean
		resolved := filepath.Join(fromDir, importPath)
		resolved = filepath.Clean(resolved)
		resolved = filepath.ToSlash(resolved)

		return b.resolveFileExtension(projectRoot, resolved)
	}

	// Handle absolute imports (from project root or path aliases)
	if !strings.Contains(importPath, "node_modules") {
		// Try to resolve as a project file
		resolved := b.resolveFileExtension(projectRoot, importPath)
		if resolved != importPath {
			return resolved
		}
	}

	return importPath
}

// resolveFileExtension tries to resolve a path to an actual file with extension
func (b *TSImportGraphBuilder) resolveFileExtension(projectRoot, importPath string) string {
	// Clean the path
	importPath = filepath.Clean(importPath)
	importPath = filepath.ToSlash(importPath)

	// If already has extension, check if file exists
	if ext := filepath.Ext(importPath); ext != "" {
		fullPath := filepath.Join(projectRoot, importPath)
		if _, err := os.Stat(fullPath); err == nil {
			return importPath
		}
	}

	// Try common extensions
	extensions := []string{".ts", ".tsx", ".js", ".jsx", ".vue", ".mjs"}

	for _, ext := range extensions {
		candidate := importPath + ext
		fullPath := filepath.Join(projectRoot, candidate)
		if _, err := os.Stat(fullPath); err == nil {
			return candidate
		}
	}

	// Try index files (barrel exports)
	indexFiles := []string{
		"index.ts", "index.tsx", "index.js", "index.jsx", "index.vue",
	}

	for _, indexFile := range indexFiles {
		candidate := filepath.Join(importPath, indexFile)
		candidate = filepath.ToSlash(candidate)
		fullPath := filepath.Join(projectRoot, candidate)
		if _, err := os.Stat(fullPath); err == nil {
			return candidate
		}
	}

	return importPath
}

// GetImportPath finds the import path between two files using BFS
func (b *TSImportGraphBuilder) GetImportPath(
	ctx context.Context,
	from, to string,
	graph *domain.ImportGraph,
) ([]string, error) {
	// Normalize paths
	from = filepath.ToSlash(from)
	to = filepath.ToSlash(to)

	// BFS to find shortest path
	if from == to {
		return []string{from}, nil
	}

	// Build adjacency list
	adjacency := make(map[string][]string)
	for _, edge := range graph.Imports {
		adjacency[edge.From] = append(adjacency[edge.From], edge.To)
	}

	// BFS
	visited := make(map[string]bool)
	parent := make(map[string]string)
	queue := []string{from}
	visited[from] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, neighbor := range adjacency[current] {
			if visited[neighbor] {
				continue
			}

			visited[neighbor] = true
			parent[neighbor] = current

			if neighbor == to {
				// Reconstruct path
				path := []string{to}
				for node := to; node != from; node = parent[node] {
					if p, ok := parent[node]; ok {
						path = append([]string{p}, path...)
					} else {
						break
					}
				}
				return path, nil
			}

			queue = append(queue, neighbor)
		}
	}

	return nil, fmt.Errorf("no import path found from %s to %s", from, to)
}

// GetCircularImports finds all circular import dependencies using Tarjan's algorithm
func (b *TSImportGraphBuilder) GetCircularImports(
	ctx context.Context,
	graph *domain.ImportGraph,
) ([][]string, error) {
	// Build adjacency list
	adjacency := make(map[string][]string)
	allNodes := make(map[string]bool)

	for _, edge := range graph.Imports {
		adjacency[edge.From] = append(adjacency[edge.From], edge.To)
		allNodes[edge.From] = true
		allNodes[edge.To] = true
	}

	// Tarjan's algorithm for finding strongly connected components
	var (
		index      = 0
		stack      []string
		onStack    = make(map[string]bool)
		indices    = make(map[string]int)
		lowLinks   = make(map[string]int)
		cycles     [][]string
	)

	var strongConnect func(v string)
	strongConnect = func(v string) {
		indices[v] = index
		lowLinks[v] = index
		index++
		stack = append(stack, v)
		onStack[v] = true

		for _, w := range adjacency[v] {
			if _, ok := indices[w]; !ok {
				strongConnect(w)
				if lowLinks[w] < lowLinks[v] {
					lowLinks[v] = lowLinks[w]
				}
			} else if onStack[w] {
				if indices[w] < lowLinks[v] {
					lowLinks[v] = indices[w]
				}
			}
		}

		// If v is a root node, pop the stack and generate an SCC
		if lowLinks[v] == indices[v] {
			var scc []string
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				scc = append(scc, w)
				if w == v {
					break
				}
			}

			// Only report SCCs with more than one node (actual cycles)
			if len(scc) > 1 {
				// Reverse to get correct order
				for i, j := 0, len(scc)-1; i < j; i, j = i+1, j-1 {
					scc[i], scc[j] = scc[j], scc[i]
				}
				cycles = append(cycles, scc)
			}
		}
	}

	// Run Tarjan's algorithm on all nodes
	for node := range allNodes {
		if _, ok := indices[node]; !ok {
			strongConnect(node)
		}
	}

	// Sort cycles for determinism
	for _, cycle := range cycles {
		sort.Strings(cycle)
	}
	sort.Slice(cycles, func(i, j int) bool {
		if len(cycles[i]) != len(cycles[j]) {
			return len(cycles[i]) < len(cycles[j])
		}
		for k := 0; k < len(cycles[i]); k++ {
			if cycles[i][k] != cycles[j][k] {
				return cycles[i][k] < cycles[j][k]
			}
		}
		return false
	})

	b.log.Info(fmt.Sprintf("Found %d circular import cycles", len(cycles)))

	return cycles, nil
}

// sortGraphForDeterminism sorts the graph elements for consistent output
func (b *TSImportGraphBuilder) sortGraphForDeterminism(graph *domain.ImportGraph) {
	// Sort imports
	sort.Slice(graph.Imports, func(i, j int) bool {
		if graph.Imports[i].From != graph.Imports[j].From {
			return graph.Imports[i].From < graph.Imports[j].From
		}
		return graph.Imports[i].To < graph.Imports[j].To
	})

	// Sort exports and imports within each package
	for _, pkg := range graph.Packages {
		sort.Strings(pkg.Imports)
		sort.Strings(pkg.Exports)
	}
}
