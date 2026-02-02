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

// DependencyAnalyzer defines the interface for analyzing file dependencies
type DependencyAnalyzer interface {
	// AnalyzeFile extracts imports from a single file
	AnalyzeFile(projectRoot, filePath string) ([]FileDependency, error)

	// BuildGraph builds the complete dependency graph for a project
	BuildGraph(projectRoot string, ignorePatterns []string) (*DependencyGraph, error)

	// ClearCache clears the cached dependency graph
	ClearCache()
}
