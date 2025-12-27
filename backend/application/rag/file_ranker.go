// Package rag provides Retrieval Augmented Generation services
package rag

import (
	"context"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// RankedFile represents a file with its PageRank score
type RankedFile struct {
	FilePath    string  `json:"filePath"`
	Rank        float64 `json:"rank"`
	ImportCount int     `json:"importCount"` // How many files import this
	ExportCount int     `json:"exportCount"` // How many symbols exported
	FileSize    int64   `json:"fileSize"`
}

// FileRanker ranks files using PageRank algorithm based on import graph
type FileRanker interface {
	// RankFiles returns files ranked by importance using PageRank
	RankFiles(ctx context.Context, projectRoot string) ([]RankedFile, error)

	// InvalidateCache clears the cache for a project
	InvalidateCache(projectRoot string)
}

// FileRankerImpl implements FileRanker using PageRank algorithm
type FileRankerImpl struct {
	analyzerRegistry analysis.AnalyzerRegistry
	log              domain.Logger

	mu    sync.RWMutex
	cache map[string]*fileRankCache
}

type fileRankCache struct {
	rankedFiles []RankedFile
	graph       *importGraph
}

// importGraph represents the file dependency graph
type importGraph struct {
	// nodes maps file path to node index
	nodes map[string]int
	// files maps node index to file path
	files []string
	// adjacency[i] contains indices of files that file i imports
	adjacency [][]int
	// reverseAdj[i] contains indices of files that import file i
	reverseAdj [][]int
	// exportCounts[i] is the number of exports in file i
	exportCounts []int
	// fileSizes[i] is the size of file i
	fileSizes []int64
}

// PageRank algorithm constants
const (
	dampingFactor    = 0.85
	maxIterations    = 100
	convergenceEps   = 1e-6
	minFileRank      = 0.0001
	entryPointBoost  = 2.0
	testFileDiscount = 0.3
)

// NewFileRanker creates a new FileRanker
func NewFileRanker(
	analyzerRegistry analysis.AnalyzerRegistry,
	log domain.Logger,
) *FileRankerImpl {
	return &FileRankerImpl{
		analyzerRegistry: analyzerRegistry,
		log:              log,
		cache:            make(map[string]*fileRankCache),
	}
}

// RankFiles returns files ranked by importance using PageRank
func (r *FileRankerImpl) RankFiles(ctx context.Context, projectRoot string) ([]RankedFile, error) {
	// Check cache
	if cached := r.getCached(projectRoot); cached != nil {
		return cached, nil
	}

	r.log.Info(fmt.Sprintf("Ranking files in %s using PageRank", projectRoot))

	// Build import graph
	graph, err := r.buildImportGraph(ctx, projectRoot)
	if err != nil {
		return nil, fmt.Errorf("failed to build import graph: %w", err)
	}

	if len(graph.files) == 0 {
		return []RankedFile{}, nil
	}

	// Run PageRank
	ranks := r.pageRank(graph)

	// Convert to RankedFile slice
	rankedFiles := make([]RankedFile, len(graph.files))
	for i, filePath := range graph.files {
		rankedFiles[i] = RankedFile{
			FilePath:    filePath,
			Rank:        ranks[i],
			ImportCount: len(graph.reverseAdj[i]),
			ExportCount: graph.exportCounts[i],
			FileSize:    graph.fileSizes[i],
		}
	}

	// Sort by rank descending
	sort.Slice(rankedFiles, func(i, j int) bool {
		return rankedFiles[i].Rank > rankedFiles[j].Rank
	})

	// Cache results
	r.setCache(projectRoot, rankedFiles, graph)

	r.log.Info(fmt.Sprintf("Ranked %d files", len(rankedFiles)))

	return rankedFiles, nil
}

// InvalidateCache clears the cache for a project
func (r *FileRankerImpl) InvalidateCache(projectRoot string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.cache, projectRoot)
}

// getCached returns cached ranked files if available
func (r *FileRankerImpl) getCached(projectRoot string) []RankedFile {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if cache, ok := r.cache[projectRoot]; ok {
		return cache.rankedFiles
	}
	return nil
}

// setCache stores ranked files in cache
func (r *FileRankerImpl) setCache(projectRoot string, rankedFiles []RankedFile, graph *importGraph) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache[projectRoot] = &fileRankCache{
		rankedFiles: rankedFiles,
		graph:       graph,
	}
}

// buildImportGraph builds the import dependency graph
func (r *FileRankerImpl) buildImportGraph(ctx context.Context, projectRoot string) (*importGraph, error) {
	graph := &importGraph{
		nodes:        make(map[string]int),
		files:        make([]string, 0),
		adjacency:    make([][]int, 0),
		reverseAdj:   make([][]int, 0),
		exportCounts: make([]int, 0),
		fileSizes:    make([]int64, 0),
	}

	// Collect all analyzable files
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		if info.IsDir() {
			// Skip common non-source directories
			name := info.Name()
			if name == "node_modules" || name == ".git" || name == "vendor" ||
				name == "dist" || name == "build" || name == "__pycache__" ||
				name == ".idea" || name == ".vscode" {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if we have an analyzer for this file
		relPath, _ := filepath.Rel(projectRoot, path)
		relPath = filepath.ToSlash(relPath)

		if r.analyzerRegistry.GetAnalyzer(relPath) != nil {
			idx := len(graph.files)
			graph.nodes[relPath] = idx
			graph.files = append(graph.files, relPath)
			graph.adjacency = append(graph.adjacency, []int{})
			graph.reverseAdj = append(graph.reverseAdj, []int{})
			graph.exportCounts = append(graph.exportCounts, 0)
			graph.fileSizes = append(graph.fileSizes, info.Size())
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	// Build edges by analyzing imports
	for filePath, nodeIdx := range graph.nodes {
		imports, exports, err := r.analyzeFileImportsExports(ctx, projectRoot, filePath)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Failed to analyze %s: %v", filePath, err))
			continue
		}

		graph.exportCounts[nodeIdx] = exports

		// Resolve imports to file indices
		for _, imp := range imports {
			targetPath := r.resolveImport(filePath, imp, graph.nodes)
			if targetIdx, ok := graph.nodes[targetPath]; ok && targetIdx != nodeIdx {
				graph.adjacency[nodeIdx] = append(graph.adjacency[nodeIdx], targetIdx)
				graph.reverseAdj[targetIdx] = append(graph.reverseAdj[targetIdx], nodeIdx)
			}
		}
	}

	return graph, nil
}

// analyzeFileImportsExports extracts imports and export count from a file
func (r *FileRankerImpl) analyzeFileImportsExports(ctx context.Context, projectRoot, filePath string) ([]string, int, error) {
	fullPath := filepath.Join(projectRoot, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, 0, err
	}

	analyzer := r.analyzerRegistry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil, 0, nil
	}

	// Get imports
	imports, err := analyzer.GetImports(ctx, filePath, content)
	if err != nil {
		return nil, 0, err
	}

	importPaths := make([]string, 0, len(imports))
	for _, imp := range imports {
		if imp.IsLocal {
			importPaths = append(importPaths, imp.Path)
		}
	}

	// Get exports count
	exports, err := analyzer.GetExports(ctx, filePath, content)
	if err != nil {
		return importPaths, 0, nil
	}

	return importPaths, len(exports), nil
}

// resolveImport resolves an import path to a file path
func (r *FileRankerImpl) resolveImport(fromFile, importPath string, nodes map[string]int) string {
	if importPath == "" {
		return ""
	}

	// Handle relative imports
	if strings.HasPrefix(importPath, ".") {
		dir := filepath.Dir(fromFile)
		resolved := filepath.Join(dir, importPath)
		resolved = filepath.ToSlash(filepath.Clean(resolved))

		// Try with common extensions
		extensions := []string{"", ".ts", ".tsx", ".js", ".jsx", ".go", ".py", ".vue"}
		for _, ext := range extensions {
			candidate := resolved + ext
			if _, ok := nodes[candidate]; ok {
				return candidate
			}
		}

		// Try index files
		indexFiles := []string{"/index.ts", "/index.tsx", "/index.js", "/index.jsx"}
		for _, idx := range indexFiles {
			candidate := resolved + idx
			if _, ok := nodes[candidate]; ok {
				return candidate
			}
		}

		return resolved
	}

	// Handle alias imports (e.g., @/components/...)
	if strings.HasPrefix(importPath, "@/") {
		// Common alias for src/
		resolved := "src/" + importPath[2:]
		extensions := []string{"", ".ts", ".tsx", ".js", ".jsx", ".vue"}
		for _, ext := range extensions {
			candidate := resolved + ext
			if _, ok := nodes[candidate]; ok {
				return candidate
			}
		}
	}

	return importPath
}

// pageRank runs the PageRank algorithm on the import graph
func (r *FileRankerImpl) pageRank(graph *importGraph) []float64 {
	n := len(graph.files)
	if n == 0 {
		return []float64{}
	}

	// Initialize ranks
	ranks := make([]float64, n)
	newRanks := make([]float64, n)
	initialRank := 1.0 / float64(n)

	for i := range ranks {
		ranks[i] = initialRank
	}

	// Identify entry points and test files for boosting/discounting
	entryPoints := make([]bool, n)
	testFiles := make([]bool, n)

	for i, filePath := range graph.files {
		lower := strings.ToLower(filePath)

		// Entry points: main files, index files, app files
		if strings.Contains(lower, "main.") || strings.Contains(lower, "index.") ||
			strings.Contains(lower, "app.") || strings.HasSuffix(lower, "/mod.rs") {
			entryPoints[i] = true
		}

		// Test files
		if strings.Contains(lower, "_test.") || strings.Contains(lower, ".test.") ||
			strings.Contains(lower, ".spec.") || strings.Contains(lower, "/test/") ||
			strings.Contains(lower, "/tests/") {
			testFiles[i] = true
		}
	}

	// Iterate until convergence
	for iter := 0; iter < maxIterations; iter++ {
		// Calculate new ranks
		for i := range newRanks {
			newRanks[i] = (1 - dampingFactor) / float64(n)
		}

		// Distribute rank through edges
		for i := 0; i < n; i++ {
			outDegree := len(graph.adjacency[i])
			if outDegree == 0 {
				// Dangling node: distribute to all nodes
				contribution := dampingFactor * ranks[i] / float64(n)
				for j := range newRanks {
					newRanks[j] += contribution
				}
			} else {
				// Distribute to linked nodes
				contribution := dampingFactor * ranks[i] / float64(outDegree)
				for _, j := range graph.adjacency[i] {
					newRanks[j] += contribution
				}
			}
		}

		// Check convergence
		diff := 0.0
		for i := range ranks {
			diff += math.Abs(newRanks[i] - ranks[i])
		}

		// Swap ranks
		ranks, newRanks = newRanks, ranks

		if diff < convergenceEps {
			r.log.Info(fmt.Sprintf("PageRank converged after %d iterations", iter+1))
			break
		}
	}

	// Apply boosts and discounts
	for i := range ranks {
		// Boost entry points
		if entryPoints[i] {
			ranks[i] *= entryPointBoost
		}

		// Discount test files
		if testFiles[i] {
			ranks[i] *= testFileDiscount
		}

		// Boost files with many exports (they're likely important)
		if graph.exportCounts[i] > 5 {
			ranks[i] *= 1.0 + float64(graph.exportCounts[i])*0.02
		}

		// Slight boost for files that are imported by many
		importedBy := len(graph.reverseAdj[i])
		if importedBy > 3 {
			ranks[i] *= 1.0 + float64(importedBy)*0.05
		}

		// Ensure minimum rank
		if ranks[i] < minFileRank {
			ranks[i] = minFileRank
		}
	}

	// Normalize ranks
	maxRank := 0.0
	for _, rank := range ranks {
		if rank > maxRank {
			maxRank = rank
		}
	}

	if maxRank > 0 {
		for i := range ranks {
			ranks[i] /= maxRank
		}
	}

	return ranks
}

// GetTopFiles returns the top N ranked files
func (r *FileRankerImpl) GetTopFiles(ctx context.Context, projectRoot string, n int) ([]RankedFile, error) {
	rankedFiles, err := r.RankFiles(ctx, projectRoot)
	if err != nil {
		return nil, err
	}

	if n > len(rankedFiles) {
		n = len(rankedFiles)
	}

	return rankedFiles[:n], nil
}

// GetFileRank returns the rank of a specific file
func (r *FileRankerImpl) GetFileRank(ctx context.Context, projectRoot, filePath string) (float64, error) {
	rankedFiles, err := r.RankFiles(ctx, projectRoot)
	if err != nil {
		return 0, err
	}

	for _, rf := range rankedFiles {
		if rf.FilePath == filePath {
			return rf.Rank, nil
		}
	}

	return 0, nil
}

// GetImportGraph returns the import graph for debugging/visualization
func (r *FileRankerImpl) GetImportGraph(ctx context.Context, projectRoot string) (map[string][]string, error) {
	// Ensure graph is built
	_, err := r.RankFiles(ctx, projectRoot)
	if err != nil {
		return nil, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	cache, ok := r.cache[projectRoot]
	if !ok || cache.graph == nil {
		return nil, fmt.Errorf("graph not available")
	}

	graph := cache.graph
	result := make(map[string][]string)

	for i, filePath := range graph.files {
		imports := make([]string, 0, len(graph.adjacency[i]))
		for _, j := range graph.adjacency[i] {
			imports = append(imports, graph.files[j])
		}
		result[filePath] = imports
	}

	return result, nil
}
