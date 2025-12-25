package domain

// ProjectStructure represents the detected project structure and architecture
type ProjectStructure struct {
	Architecture *ArchitectureInfo `json:"architecture"`
	Conventions  *ConventionInfo   `json:"conventions"`
	Frameworks   []FrameworkInfo   `json:"frameworks"`
	BuildSystems []BuildSystemInfo `json:"buildSystems"`
	Languages    []LanguageInfo    `json:"languages"`
	Layers       []LayerInfo       `json:"layers"`
	ProjectType  string            `json:"projectType"` // web, cli, library, service, monorepo
	Confidence   float64           `json:"confidence"`
}

// ArchitectureType represents different architecture patterns
type ArchitectureType string

const (
	ArchCleanArchitecture ArchitectureType = "clean"
	ArchHexagonal         ArchitectureType = "hexagonal"
	ArchMVC               ArchitectureType = "mvc"
	ArchMVVM              ArchitectureType = "mvvm"
	ArchLayered           ArchitectureType = "layered"
	ArchMicroservices     ArchitectureType = "microservices"
	ArchMonolith          ArchitectureType = "monolith"
	ArchServerless        ArchitectureType = "serverless"
	ArchEventDriven       ArchitectureType = "event-driven"
	ArchDDD               ArchitectureType = "ddd"
	ArchUnknown           ArchitectureType = "unknown"
)

// ArchitectureInfo contains detected architecture information
type ArchitectureInfo struct {
	Type        ArchitectureType `json:"type"`
	Confidence  float64          `json:"confidence"`
	Description string           `json:"description"`
	Indicators  []string         `json:"indicators"` // What led to this detection
	Layers      []LayerInfo      `json:"layers"`
}

// LayerInfo represents an architectural layer
type LayerInfo struct {
	Name         string   `json:"name"` // domain, application, infrastructure, presentation, handlers, services
	Path         string   `json:"path"` // Directory path
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"` // Which layers this depends on
	Files        []string `json:"files"`        // Sample files in this layer
	Patterns     []string `json:"patterns"`     // Detected patterns (repository, service, handler, etc.)
}

// ConventionInfo contains detected naming and structure conventions
type ConventionInfo struct {
	NamingStyle     NamingStyle     `json:"namingStyle"`     // camelCase, snake_case, PascalCase, kebab-case
	FileNaming      FileNamingStyle `json:"fileNaming"`      // How files are named
	FolderStructure FolderStructure `json:"folderStructure"` // flat, by-feature, by-type, by-layer
	TestConventions TestConventions `json:"testConventions"`
	ImportStyle     ImportStyle     `json:"importStyle"`
	CodeStyle       CodeStyleInfo   `json:"codeStyle"`
}

// NamingStyle represents naming conventions
type NamingStyle string

const (
	NamingCamelCase  NamingStyle = "camelCase"
	NamingSnakeCase  NamingStyle = "snake_case"
	NamingPascalCase NamingStyle = "PascalCase"
	NamingKebabCase  NamingStyle = "kebab-case"
	NamingMixed      NamingStyle = "mixed"
)

// FileNamingStyle represents file naming conventions
type FileNamingStyle struct {
	Style    NamingStyle `json:"style"`
	Suffixes []string    `json:"suffixes"` // Common suffixes like _test, .spec, .service
	Prefixes []string    `json:"prefixes"` // Common prefixes
	Examples []string    `json:"examples"`
}

// FolderStructure represents folder organization style
type FolderStructure string

const (
	FolderFlat      FolderStructure = "flat"
	FolderByFeature FolderStructure = "by-feature"
	FolderByType    FolderStructure = "by-type"
	FolderByLayer   FolderStructure = "by-layer"
	FolderHybrid    FolderStructure = "hybrid"
)

// TestConventions represents testing conventions
type TestConventions struct {
	Location   string   `json:"location"`   // same-dir, __tests__, tests/, spec/
	FileSuffix string   `json:"fileSuffix"` // _test, .test, .spec, _spec
	Framework  string   `json:"framework"`  // jest, vitest, go test, pytest, junit
	Patterns   []string `json:"patterns"`   // BDD, TDD, etc.
}

// ImportStyle represents import organization style
type ImportStyle struct {
	AbsoluteImports bool     `json:"absoluteImports"` // Uses absolute imports
	AliasedImports  bool     `json:"aliasedImports"`  // Uses @ or ~ aliases
	GroupedImports  bool     `json:"groupedImports"`  // Groups imports by type
	ImportOrder     []string `json:"importOrder"`     // stdlib, external, internal
}

// CodeStyleInfo represents code style conventions
type CodeStyleInfo struct {
	IndentStyle   string `json:"indentStyle"` // tabs, spaces
	IndentSize    int    `json:"indentSize"`
	MaxLineLength int    `json:"maxLineLength"`
	TrailingComma bool   `json:"trailingComma"`
	Semicolons    bool   `json:"semicolons"`
	QuoteStyle    string `json:"quoteStyle"` // single, double
	ConfigFile    string `json:"configFile"` // .prettierrc, .editorconfig, etc.
}

// FrameworkInfo contains detected framework information
type FrameworkInfo struct {
	Name          string   `json:"name"`
	Version       string   `json:"version"`
	Category      string   `json:"category"` // web, cli, testing, orm, etc.
	Language      string   `json:"language"`
	ConfigFiles   []string `json:"configFiles"`   // Framework config files found
	Indicators    []string `json:"indicators"`    // What led to detection
	BestPractices []string `json:"bestPractices"` // Suggested best practices
}

// BuildSystemInfo contains build system information
type BuildSystemInfo struct {
	Name       string   `json:"name"`       // make, npm, gradle, maven, cargo, go
	ConfigFile string   `json:"configFile"` // Makefile, package.json, build.gradle
	Scripts    []string `json:"scripts"`    // Available build scripts
	Commands   []string `json:"commands"`   // Common commands
}

// LanguageInfo contains language-specific information
type LanguageInfo struct {
	Name       string  `json:"name"`
	Version    string  `json:"version"`
	FileCount  int     `json:"fileCount"`
	Percentage float64 `json:"percentage"` // Percentage of codebase
	Primary    bool    `json:"primary"`    // Is this the primary language
}
