package analysis

import "context"

// LanguageAnalyzer defines interface for language-specific code analyzers
type LanguageAnalyzer interface {
	// Language returns the language name
	Language() string

	// Extensions returns supported file extensions
	Extensions() []string

	// CanAnalyze checks if this analyzer can handle the file
	CanAnalyze(filePath string) bool

	// ExtractSymbols extracts all symbols from file content
	ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]Symbol, error)

	// GetImports extracts import statements
	GetImports(ctx context.Context, filePath string, content []byte) ([]Import, error)

	// GetExports extracts exported symbols (for JS/TS/Vue)
	GetExports(ctx context.Context, filePath string, content []byte) ([]Export, error)

	// GetFunctionBody returns the full body of a function by name
	GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error)
}

// AnalyzerRegistry manages language analyzers
type AnalyzerRegistry interface {
	// Register adds an analyzer
	Register(analyzer LanguageAnalyzer)

	// GetAnalyzer returns analyzer for file path
	GetAnalyzer(filePath string) LanguageAnalyzer

	// GetAnalyzerByLanguage returns analyzer by language name
	GetAnalyzerByLanguage(lang string) LanguageAnalyzer

	// SupportedLanguages returns list of supported languages
	SupportedLanguages() []string

	// SupportedExtensions returns list of supported extensions
	SupportedExtensions() []string
}

// SymbolIndex provides fast symbol lookup
type SymbolIndex interface {
	// IndexProject indexes all files in project
	IndexProject(ctx context.Context, projectRoot string) error

	// IndexFile indexes a single file
	IndexFile(ctx context.Context, filePath string, content []byte) error

	// SearchByName finds symbols by name (partial match)
	SearchByName(query string) []Symbol

	// FindByExactName finds symbols with exact name
	FindByExactName(name string) []Symbol

	// GetSymbolsInFile returns all symbols in a file
	GetSymbolsInFile(filePath string) []Symbol

	// GetSymbolsByKind returns symbols of specific kind
	GetSymbolsByKind(kind SymbolKind) []Symbol

	// FindDefinition finds where a symbol is defined
	FindDefinition(name string, kind SymbolKind) *Symbol

	// Stats returns index statistics
	Stats() map[string]int

	// IsIndexed returns whether project is indexed
	IsIndexed() bool

	// Clear clears the index
	Clear()
}
