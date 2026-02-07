package analysis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
	"time"
)

// DependencyAnalyzerImpl implements domain.DependencyAnalyzer
type DependencyAnalyzerImpl struct {
	registry analysis.AnalyzerRegistry
	logger   domain.Logger

	// Cache for dependency graph (keyed by project root)
	mu    sync.RWMutex
	cache map[string]*domain.DependencyGraph
}

// NewDependencyAnalyzer creates a new dependency analyzer
func NewDependencyAnalyzer(registry analysis.AnalyzerRegistry, logger domain.Logger) *DependencyAnalyzerImpl {
	return &DependencyAnalyzerImpl{
		registry: registry,
		logger:   logger,
		cache:    make(map[string]*domain.DependencyGraph),
	}
}

// AnalyzeFile extracts imports from a single file
func (d *DependencyAnalyzerImpl) AnalyzeFile(projectRoot, filePath string) ([]domain.FileDependency, error) {
	// Get analyzer for this file type
	analyzer := d.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil, fmt.Errorf("no analyzer found for file: %s", filePath)
	}

	// Read file content
	fullPath := filepath.Join(projectRoot, filePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	// Extract imports
	ctx := context.Background()
	imports, err := analyzer.GetImports(ctx, filePath, content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract imports from %s: %w", filePath, err)
	}

	// Convert imports to file dependencies
	dependencies := make([]domain.FileDependency, 0, len(imports))
	for _, imp := range imports {
		targetPath, err := d.resolveImportPath(projectRoot, filePath, imp.Path)
		if err != nil {
			// Skip unresolvable imports (external packages, etc.)
			d.logger.Debug(fmt.Sprintf("Could not resolve import %s in %s: %v", imp.Path, filePath, err))
			continue
		}

		depType := d.determineImportType(imp.Path, analyzer.Language())
		dependencies = append(dependencies, domain.FileDependency{
			SourcePath: filePath,
			TargetPath: targetPath,
			Type:       depType,
			ImportName: imp.Alias,
			Line:       imp.Line,
		})
	}

	return dependencies, nil
}

// BuildGraph builds the complete dependency graph for a project
func (d *DependencyAnalyzerImpl) BuildGraph(projectRoot string, ignorePatterns []string) (*domain.DependencyGraph, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	graph := &domain.DependencyGraph{
		Files:        make(map[string][]domain.FileDependency),
		LastAnalyzed: time.Now(),
	}

	// Walk the project directory
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Check if directory should be ignored
			if d.shouldIgnore(path, projectRoot, ignorePatterns) {
				return filepath.SkipDir
			}
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(projectRoot, path)
		if err != nil {
			return err
		}

		// Normalize path separators to forward slashes
		relPath = filepath.ToSlash(relPath)

		// Check if file should be ignored
		if d.shouldIgnore(path, projectRoot, ignorePatterns) {
			return nil
		}

		// Check if we have an analyzer for this file
		analyzer := d.registry.GetAnalyzer(relPath)
		if analyzer == nil {
			return nil
		}

		// Analyze file dependencies
		deps, err := d.AnalyzeFile(projectRoot, relPath)
		if err != nil {
			d.logger.Warning(fmt.Sprintf("Failed to analyze %s: %v", relPath, err))
			return nil // Continue with other files
		}

		// Add file to graph even if it has no dependencies
		graph.Files[relPath] = deps

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk project directory: %w", err)
	}

	// Cache the graph
	d.cache[projectRoot] = graph

	d.logger.Info(fmt.Sprintf("Built dependency graph with %d files", len(graph.Files)))
	return graph, nil
}

// ClearCache clears the cached dependency graph
func (d *DependencyAnalyzerImpl) ClearCache() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = make(map[string]*domain.DependencyGraph)
	d.logger.Debug("Dependency graph cache cleared")
}

// GetFileDependenciesBatch returns dependencies for multiple files (batch query)
func (d *DependencyAnalyzerImpl) GetFileDependenciesBatch(projectRoot string, filePaths []string) (map[string][]domain.FileDependency, error) {
	d.mu.RLock()
	graph, exists := d.cache[projectRoot]
	d.mu.RUnlock()

	// If no cache, build it first
	if !exists || graph == nil {
		var err error
		graph, err = d.BuildGraph(projectRoot, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build dependency graph: %w", err)
		}
	}

	result := make(map[string][]domain.FileDependency)
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, filePath := range filePaths {
		if deps, ok := graph.Files[filePath]; ok {
			result[filePath] = deps
		} else {
			// File not in cache, return empty slice
			result[filePath] = []domain.FileDependency{}
		}
	}

	return result, nil
}

// GetIncomingDependencies returns files that import this file
func (d *DependencyAnalyzerImpl) GetIncomingDependencies(projectRoot, filePath string) ([]string, error) {
	d.mu.RLock()
	graph, exists := d.cache[projectRoot]
	d.mu.RUnlock()

	if !exists || graph == nil {
		var err error
		graph, err = d.BuildGraph(projectRoot, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build dependency graph: %w", err)
		}
	}

	incoming := make([]string, 0)
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Iterate through all files and find those that depend on this file
	for sourceFile, deps := range graph.Files {
		for _, dep := range deps {
			if dep.TargetPath == filePath {
				incoming = append(incoming, sourceFile)
				break
			}
		}
	}

	return incoming, nil
}

// GetOutgoingDependencies returns files that this file imports
func (d *DependencyAnalyzerImpl) GetOutgoingDependencies(projectRoot, filePath string) ([]string, error) {
	d.mu.RLock()
	graph, exists := d.cache[projectRoot]
	d.mu.RUnlock()

	if !exists || graph == nil {
		var err error
		graph, err = d.BuildGraph(projectRoot, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to build dependency graph: %w", err)
		}
	}

	d.mu.RLock()
	defer d.mu.RUnlock()

	deps, ok := graph.Files[filePath]
	if !ok {
		return []string{}, nil
	}

	outgoing := make([]string, 0, len(deps))
	for _, dep := range deps {
		outgoing = append(outgoing, dep.TargetPath)
	}

	return outgoing, nil
}

// IsDependencyGraphCached checks if graph is cached
func (d *DependencyAnalyzerImpl) IsDependencyGraphCached(projectRoot string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	_, exists := d.cache[projectRoot]
	return exists
}

// GetDependencyGraphStats returns cache statistics
func (d *DependencyAnalyzerImpl) GetDependencyGraphStats(projectRoot string) (*domain.DependencyGraphStats, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	graph, exists := d.cache[projectRoot]
	if !exists || graph == nil {
		return &domain.DependencyGraphStats{
			FileCount:       0,
			DependencyCount: 0,
			IsCached:        false,
			CacheSize:       0,
		}, nil
	}

	// Calculate total dependencies
	totalDeps := 0
	for _, deps := range graph.Files {
		totalDeps += len(deps)
	}

	// Estimate cache size (rough approximation)
	// Each file path ~100 bytes, each dependency ~150 bytes
	cacheSize := int64(len(graph.Files)*100 + totalDeps*150)

	return &domain.DependencyGraphStats{
		FileCount:       len(graph.Files),
		DependencyCount: totalDeps,
		LastAnalyzed:    graph.LastAnalyzed,
		IsCached:        true,
		CacheSize:       cacheSize,
	}, nil
}

// UpdateFile updates single file in cached graph (incremental update)
func (d *DependencyAnalyzerImpl) UpdateFile(projectRoot, filePath string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	graph, exists := d.cache[projectRoot]
	if !exists || graph == nil {
		return fmt.Errorf("no cached graph for project: %s", projectRoot)
	}

	// Analyze the file
	deps, err := d.AnalyzeFile(projectRoot, filePath)
	if err != nil {
		return fmt.Errorf("failed to analyze file %s: %w", filePath, err)
	}

	// Update the graph
	graph.Files[filePath] = deps
	graph.LastAnalyzed = time.Now()

	d.logger.Debug(fmt.Sprintf("Updated file in dependency graph: %s", filePath))
	return nil
}

// RemoveFile removes file from cached graph
func (d *DependencyAnalyzerImpl) RemoveFile(projectRoot, filePath string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	graph, exists := d.cache[projectRoot]
	if !exists || graph == nil {
		return fmt.Errorf("no cached graph for project: %s", projectRoot)
	}

	delete(graph.Files, filePath)
	graph.LastAnalyzed = time.Now()

	d.logger.Debug(fmt.Sprintf("Removed file from dependency graph: %s", filePath))
	return nil
}

// resolveImportPath resolves an import path to a project-relative file path
func (d *DependencyAnalyzerImpl) resolveImportPath(projectRoot, sourceFile, importPath string) (string, error) {
	// Handle different import styles
	if strings.HasPrefix(importPath, "@/") {
		// Vue/TypeScript alias import (@/components/...)
		return strings.TrimPrefix(importPath, "@/"), nil
	}

	if strings.HasPrefix(importPath, "./") || strings.HasPrefix(importPath, "../") {
		// Relative import
		sourceDir := filepath.Dir(sourceFile)
		resolved := filepath.Join(sourceDir, importPath)
		resolved = filepath.Clean(resolved)
		resolved = filepath.ToSlash(resolved)

		// Try to find the actual file with common extensions
		targetPath, err := d.findFileWithExtension(projectRoot, resolved)
		if err != nil {
			return "", err
		}
		return targetPath, nil
	}

	// Absolute imports or external packages - skip
	return "", fmt.Errorf("external or absolute import: %s", importPath)
}

// findFileWithExtension tries to find a file with common extensions
func (d *DependencyAnalyzerImpl) findFileWithExtension(projectRoot, basePath string) (string, error) {
	// Common extensions to try (including non-code files like CSS)
	extensions := []string{"", ".ts", ".tsx", ".js", ".jsx", ".vue", ".go", ".py", ".java", ".css", ".scss", ".sass", ".less"}

	for _, ext := range extensions {
		testPath := basePath + ext
		fullPath := filepath.Join(projectRoot, testPath)

		if _, err := os.Stat(fullPath); err == nil {
			return testPath, nil
		}

		// Also try index files in directories (skip for style files)
		if ext != "" && !strings.HasPrefix(ext, ".css") && !strings.HasPrefix(ext, ".scss") && 
			!strings.HasPrefix(ext, ".sass") && !strings.HasPrefix(ext, ".less") {
			indexPath := filepath.Join(basePath, "index"+ext)
			fullIndexPath := filepath.Join(projectRoot, indexPath)
			if _, err := os.Stat(fullIndexPath); err == nil {
				return indexPath, nil
			}
		}
	}

	return "", fmt.Errorf("file not found: %s", basePath)
}

// determineImportType determines the type of import
func (d *DependencyAnalyzerImpl) determineImportType(importPath, language string) string {
	// Check for style imports
	if strings.HasSuffix(importPath, ".css") ||
		strings.HasSuffix(importPath, ".scss") ||
		strings.HasSuffix(importPath, ".sass") ||
		strings.HasSuffix(importPath, ".less") {
		return "style"
	}

	// Check for test imports
	if strings.Contains(importPath, ".test.") ||
		strings.Contains(importPath, ".spec.") ||
		strings.Contains(importPath, "_test") {
		return "test"
	}

	// Check for type-only imports (TypeScript)
	if language == "typescript" && strings.Contains(importPath, ".d.ts") {
		return "type"
	}

	return "import"
}

// shouldIgnore checks if a path should be ignored based on patterns
func (d *DependencyAnalyzerImpl) shouldIgnore(path, projectRoot string, ignorePatterns []string) bool {
	relPath, err := filepath.Rel(projectRoot, path)
	if err != nil {
		return false
	}

	// Normalize to forward slashes
	relPath = filepath.ToSlash(relPath)

	// Common directories to always ignore
	commonIgnores := []string{
		"node_modules",
		".git",
		"dist",
		"build",
		"vendor",
		".next",
		".nuxt",
		"coverage",
		".cache",
	}

	for _, ignore := range commonIgnores {
		if strings.Contains(relPath, ignore) {
			return true
		}
	}

	// Check custom ignore patterns
	for _, pattern := range ignorePatterns {
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}
	}

	return false
}
