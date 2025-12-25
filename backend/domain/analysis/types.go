// Package analysis defines domain types for code analysis
package analysis

// Symbol represents a code symbol (function, class, type, etc.)
type Symbol struct {
	Name       string            `json:"name"`
	Kind       SymbolKind        `json:"kind"`
	Language   string            `json:"language"`
	FilePath   string            `json:"filePath"`
	Line       int               `json:"line"` // alias for StartLine
	StartLine  int               `json:"startLine"`
	EndLine    int               `json:"endLine"`
	StartCol   int               `json:"startCol"`
	EndCol     int               `json:"endCol"`
	Signature  string            `json:"signature,omitempty"`
	DocComment string            `json:"docComment,omitempty"`
	Parent     string            `json:"parent,omitempty"`
	Modifiers  []string          `json:"modifiers,omitempty"`
	Children   []Symbol          `json:"children,omitempty"`
	Extra      map[string]string `json:"extra,omitempty"`
}

// SymbolKind represents the type of symbol
type SymbolKind string

const (
	KindFunction   SymbolKind = "function"
	KindMethod     SymbolKind = "method"
	KindClass      SymbolKind = "class"
	KindInterface  SymbolKind = "interface"
	KindStruct     SymbolKind = "struct"
	KindType       SymbolKind = "type"
	KindVariable   SymbolKind = "variable"
	KindConstant   SymbolKind = "constant"
	KindProperty   SymbolKind = "property"
	KindField      SymbolKind = "field"
	KindEnum       SymbolKind = "enum"
	KindEnumMember SymbolKind = "enumMember"
	KindModule     SymbolKind = "module"
	KindPackage    SymbolKind = "package"
	KindImport     SymbolKind = "import"
	KindWidget     SymbolKind = "widget"
	KindComponent  SymbolKind = "component"
	KindComposable SymbolKind = "composable"

	// Aliases for backward compatibility
	SymbolFunction  SymbolKind = "function"
	SymbolMethod    SymbolKind = "method"
	SymbolClass     SymbolKind = "class"
	SymbolInterface SymbolKind = "interface"
	SymbolType      SymbolKind = "type"
	SymbolVariable  SymbolKind = "variable"
	SymbolConstant  SymbolKind = "constant"
)

// Import represents an import statement
type Import struct {
	Path    string   `json:"path"`
	Alias   string   `json:"alias,omitempty"`
	Names   []string `json:"names,omitempty"`
	IsLocal bool     `json:"isLocal"`
	Line    int      `json:"line,omitempty"`
}

// Export represents an exported symbol
type Export struct {
	Name       string `json:"name"`
	Alias      string `json:"alias,omitempty"`
	Kind       string `json:"kind"` // function, class, const, type, default
	IsDefault  bool   `json:"isDefault"`
	IsReExport bool   `json:"isReExport"`
	FromPath   string `json:"fromPath,omitempty"` // for re-exports
	Line       int    `json:"line,omitempty"`
}

// Location represents a position in code
type Location struct {
	FilePath  string `json:"filePath"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	StartCol  int    `json:"startCol"`
	EndCol    int    `json:"endCol"`
}
