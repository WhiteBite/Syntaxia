package analysis

// CallGraph represents function call relationships
type CallGraph struct {
	Nodes map[string]*CallNode `json:"nodes"`
	Edges []CallEdge           `json:"edges"`
}

// CallNode represents a function in the call graph
type CallNode struct {
	ID        string   `json:"id"`        // unique identifier
	Name      string   `json:"name"`      // function name
	FilePath  string   `json:"filePath"`  // file where defined
	Line      int      `json:"line"`      // line number
	Package   string   `json:"package"`   // package/module name
	Signature string   `json:"signature"` // function signature
	Callers   []string `json:"callers"`   // functions that call this
	Callees   []string `json:"callees"`   // functions this calls
}

// CallEdge represents a call from one function to another
type CallEdge struct {
	From     string `json:"from"`     // caller function ID
	To       string `json:"to"`       // callee function ID
	FilePath string `json:"filePath"` // where the call happens
	Line     int    `json:"line"`     // line of the call
	CallType string `json:"callType"` // direct, method, callback, etc.
}

// DependencyGraph represents file/package dependencies
type DependencyGraph struct {
	Nodes map[string]*DependencyNode `json:"nodes"`
	Edges []DependencyEdge           `json:"edges"`
}

// DependencyNode represents a file or package in the dependency graph
type DependencyNode struct {
	ID           string   `json:"id"`           // unique identifier (file path or package name)
	Name         string   `json:"name"`         // display name
	Type         string   `json:"type"`         // "file" or "package"
	FilePath     string   `json:"filePath"`     // file path (for file nodes)
	Package      string   `json:"package"`      // package name
	Dependencies []string `json:"dependencies"` // what this depends on
	Dependents   []string `json:"dependents"`   // what depends on this
}

// DependencyEdge represents a dependency from one node to another
type DependencyEdge struct {
	From       string `json:"from"`       // source node ID
	To         string `json:"to"`         // target node ID
	ImportPath string `json:"importPath"` // import path/statement
	Line       int    `json:"line"`       // line of import
}

// CyclicDependency represents a cycle in the dependency graph
type CyclicDependency struct {
	Cycle []string `json:"cycle"` // nodes forming the cycle
	Type  string   `json:"type"`  // "file" or "package"
}

// CallGraphBuilder builds call graphs from source code
type CallGraphBuilder interface {
	// Build builds call graph for the project
	Build(projectRoot string) (*CallGraph, error)

	// GetCallers returns functions that call the given function
	GetCallers(functionID string) []CallNode

	// GetCallees returns functions called by the given function
	GetCallees(functionID string) []CallNode

	// GetCallChain returns the call chain from start to end
	GetCallChain(startID, endID string, maxDepth int) [][]string

	// GetImpact returns all functions affected if given function changes
	GetImpact(functionID string, maxDepth int) []CallNode

	// BuildDependencyGraph builds file/package dependency graph
	BuildDependencyGraph(projectRoot string) (*DependencyGraph, error)

	// FindCyclicDependencies finds all cyclic dependencies
	FindCyclicDependencies(projectRoot string) ([]CyclicDependency, error)

	// GetFileDependencies returns files that a file depends on
	GetFileDependencies(filePath string) []DependencyNode

	// GetFileDependents returns files that depend on a file
	GetFileDependents(filePath string) []DependencyNode
}
