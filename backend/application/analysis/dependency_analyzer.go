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

	// Cache for dependency graph
	mu    sync.RWMutex
	cache *domain.DependencyGraph
}

// NewDependencyAnalyzer creates a new dependency analyzer
func NewDependencyAnalyzer(registry analysis.AnalyzerRegistry, logger domain.Logger) *DependencyAnalyzerImpl {
	return &DependencyAnalyzerImpl{
		registry: registry,
		logger:   logger,
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
	d.cache = graph

	d.logger.Info(fmt.Sprintf("Built dependency graph with %d files", len(graph.Files)))
	return graph, nil
}

// ClearCache clears the cached dependency graph
func (d *DependencyAnalyzerImpl) ClearCache() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cache = nil
	d.logger.Debug("Dependency graph cache cleared")
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
