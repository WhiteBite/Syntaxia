package domain

import "time"

// FileDependency represents a dependency between two files
type FileDependency struct {
	SourcePath string `json:"sourcePath"` // File that imports
	TargetPath string `json:"targetPath"` // File being imported
	Type       string `json:"type"`       // "import", "style", "test", "type"
	ImportName string `json:"importName,omitempty"`
	Line       int    `json:"line,omitempty"`
}

// DependencyGraph represents the complete dependency graph of a project
type DependencyGraph struct {
	Files        map[string][]FileDependency `json:"files"`        // Map of file path to its dependencies
	LastAnalyzed time.Time                   `json:"lastAnalyzed"` // When the graph was last built
}

// DependencyGraphStats contains statistics about the dependency graph cache
type DependencyGraphStats struct {
	FileCount      int       `json:"fileCount"`      // Total files in graph
	DependencyCount int      `json:"dependencyCount"` // Total dependencies
	LastAnalyzed   time.Time `json:"lastAnalyzed"`   // When graph was built
	IsCached       bool      `json:"isCached"`       // Whether graph is cached
	CacheSize      int64     `json:"cacheSize"`      // Approximate memory size in bytes
}

// DependencyAnalyzer defines the interface for analyzing file dependencies
type DependencyAnalyzer interface {
	// AnalyzeFile extracts imports from a single file
	AnalyzeFile(projectRoot, filePath string) ([]FileDependency, error)

	// BuildGraph builds the complete dependency graph for a project
	BuildGraph(projectRoot string, ignorePatterns []string) (*DependencyGraph, error)

	// ClearCache clears the cached dependency graph
	ClearCache()

	// GetFileDependenciesBatch returns dependencies for multiple files (batch query)
	GetFileDependenciesBatch(projectRoot string, filePaths []string) (map[string][]FileDependency, error)

	// GetIncomingDependencies returns files that import this file
	GetIncomingDependencies(projectRoot, filePath string) ([]string, error)

	// GetOutgoingDependencies returns files that this file imports
	GetOutgoingDependencies(projectRoot, filePath string) ([]string, error)

	// IsDependencyGraphCached checks if graph is cached
	IsDependencyGraphCached(projectRoot string) bool

	// GetDependencyGraphStats returns cache statistics
	GetDependencyGraphStats(projectRoot string) (*DependencyGraphStats, error)

	// UpdateFile updates single file in cached graph (incremental update)
	UpdateFile(projectRoot, filePath string) error

	// RemoveFile removes file from cached graph
	RemoveFile(projectRoot, filePath string) error
}
