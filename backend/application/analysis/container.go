// Package analysis provides code analysis services.
package analysis

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"syntaxia/application/tools"
	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// ContainerConfig holds factory functions for creating infrastructure implementations.
// This allows dependency injection without importing infrastructure packages.
type ContainerConfig struct {
	// RegistryFactory creates analyzer registry
	RegistryFactory func() analysis.AnalyzerRegistry
	// SymbolIndexFactory creates symbol index
	SymbolIndexFactory func(registry analysis.AnalyzerRegistry) analysis.SymbolIndex
	// CallGraphFactory creates call graph builder
	CallGraphFactory func(registry analysis.AnalyzerRegistry) domain.CallGraphBuilder
	// GitContextFactory creates git context builder
	GitContextFactory func(projectRoot string) domain.GitContextBuilder
	// ContextMemoryFactory creates context memory
	ContextMemoryFactory func(contextDir string) (domain.ContextMemory, error)
	// ProjectStructureFactory creates project structure detector
	ProjectStructureFactory func() domain.ProjectStructureDetector
	// ReferenceFinderFactory creates reference finder
	ReferenceFinderFactory func(registry analysis.AnalyzerRegistry) domain.ReferenceFinder
}

// Container manages analysis services with lazy initialization and caching.
// It provides a centralized way to access analysis tools that are expensive to create.
//
// Thread-safe: All methods can be called concurrently.
// Project-scoped: Call SetProject() when switching projects to invalidate caches.
type Container struct {
	mu          sync.RWMutex
	logger      domain.Logger
	projectRoot string
	cacheDir    string

	// Lazy-initialized services
	registry         analysis.AnalyzerRegistry
	symbolIndex      analysis.SymbolIndex
	callGraph        domain.CallGraphBuilder
	gitContext       domain.GitContextBuilder
	contextMemory    domain.ContextMemory
	projectStructure domain.ProjectStructureDetector
	referenceFinder  domain.ReferenceFinder

	// Optional services (may be nil)
	semanticSearch tools.SemanticSearcher

	// Initialization flags
	symbolIndexBuilt bool
	callGraphBuilt   bool

	// Factory functions for creating infrastructure implementations (injected via config)
	config ContainerConfig
}

// NewContainer creates a new analysis container with factory functions for DI.
func NewContainer(logger domain.Logger, config ContainerConfig) *Container {
	var registry analysis.AnalyzerRegistry
	if config.RegistryFactory != nil {
		registry = config.RegistryFactory()
	}

	homeDir, err := os.UserHomeDir()
	cacheDir := ""
	if err == nil {
		cacheDir = filepath.Join(homeDir, ".syntaxia", "analysis")
		_ = os.MkdirAll(cacheDir, 0o755)
	}

	return &Container{
		logger:   logger,
		registry: registry,
		cacheDir: cacheDir,
		config:   config,
	}
}

// SetProject sets the current project root and invalidates project-specific caches.
func (c *Container) SetProject(projectRoot string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.projectRoot == projectRoot {
		return
	}

	c.logger.Info("AnalysisContainer: switching project to " + projectRoot)
	c.projectRoot = projectRoot
	c.invalidateCachesLocked()
}

// GetProjectRoot returns the current project root.
func (c *Container) GetProjectRoot() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.projectRoot
}

// GetRegistry returns the analyzer registry.
func (c *Container) GetRegistry() analysis.AnalyzerRegistry {
	return c.registry
}

// GetSymbolIndex returns the symbol index, creating it if needed.
func (c *Container) GetSymbolIndex() analysis.SymbolIndex {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.symbolIndex == nil && c.config.SymbolIndexFactory != nil {
		c.symbolIndex = c.config.SymbolIndexFactory(c.registry)
	}
	return c.symbolIndex
}

// GetCallGraph returns the call graph builder, creating it if needed.
func (c *Container) GetCallGraph() domain.CallGraphBuilder {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.callGraph == nil && c.config.CallGraphFactory != nil {
		c.callGraph = c.config.CallGraphFactory(c.registry)
	}
	return c.callGraph
}

// GetGitContext returns the git context builder for the current project.
func (c *Container) GetGitContext() domain.GitContextBuilder {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.projectRoot == "" {
		return nil
	}

	if c.gitContext == nil && c.config.GitContextFactory != nil {
		c.gitContext = c.config.GitContextFactory(c.projectRoot)
	}
	return c.gitContext
}

// GetContextMemory returns the context memory service, creating it if needed.
func (c *Container) GetContextMemory() domain.ContextMemory {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.contextMemory == nil && c.config.ContextMemoryFactory != nil {
		cacheDir := c.cacheDir
		if cacheDir == "" {
			cacheDir = os.TempDir()
		}
		contextDir := filepath.Join(cacheDir, "contexts")
		var err error
		c.contextMemory, err = c.config.ContextMemoryFactory(contextDir)
		if err != nil {
			c.logger.Warning("Failed to create context memory: " + err.Error())
			return nil
		}
	}
	return c.contextMemory
}

// GetProjectStructure returns the project structure detector, creating it if needed.
func (c *Container) GetProjectStructure() domain.ProjectStructureDetector {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.projectStructure == nil && c.config.ProjectStructureFactory != nil {
		c.projectStructure = c.config.ProjectStructureFactory()
	}
	return c.projectStructure
}

// GetReferenceFinder returns the reference finder, creating it if needed.
func (c *Container) GetReferenceFinder() domain.ReferenceFinder {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.referenceFinder == nil && c.config.ReferenceFinderFactory != nil {
		c.referenceFinder = c.config.ReferenceFinderFactory(c.registry)
	}
	return c.referenceFinder
}

// SetSemanticSearch sets the semantic search service (optional).
func (c *Container) SetSemanticSearch(ss tools.SemanticSearcher) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.semanticSearch = ss
}

// GetSemanticSearch returns the semantic search service (may be nil).
func (c *Container) GetSemanticSearch() tools.SemanticSearcher {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.semanticSearch
}

// EnsureSymbolIndexBuilt ensures the symbol index is built for the current project.
func (c *Container) EnsureSymbolIndexBuilt(ctx context.Context) error {
	c.mu.Lock()
	if c.symbolIndexBuilt || c.projectRoot == "" {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	index := c.GetSymbolIndex()
	if err := index.IndexProject(ctx, c.projectRoot); err != nil {
		return err
	}

	c.mu.Lock()
	c.symbolIndexBuilt = true
	c.mu.Unlock()

	return nil
}

// EnsureCallGraphBuilt ensures the call graph is built for the current project.
func (c *Container) EnsureCallGraphBuilt(projectRoot string) error {
	c.mu.Lock()
	if c.callGraphBuilt {
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()

	cg := c.GetCallGraph()
	if _, err := cg.Build(projectRoot); err != nil {
		return err
	}

	c.mu.Lock()
	c.callGraphBuilt = true
	c.mu.Unlock()

	return nil
}

// InvalidateCache invalidates all caches.
func (c *Container) InvalidateCache() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.invalidateCachesLocked()
}

func (c *Container) invalidateCachesLocked() {
	if c.symbolIndex != nil {
		c.symbolIndex = nil
	}
	c.symbolIndexBuilt = false

	if c.callGraph != nil {
		c.callGraph = nil
	}
	c.callGraphBuilt = false

	c.gitContext = nil
	c.projectStructure = nil
}

// Stats returns statistics about the container state.
func (c *Container) Stats() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return map[string]interface{}{
		"projectRoot":       c.projectRoot,
		"symbolIndexBuilt":  c.symbolIndexBuilt,
		"callGraphBuilt":    c.callGraphBuilt,
		"hasGitContext":     c.gitContext != nil,
		"hasContextMemory":  c.contextMemory != nil,
		"hasSemanticSearch": c.semanticSearch != nil,
	}
}
