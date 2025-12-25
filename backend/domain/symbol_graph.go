package domain

import "context"

// SymbolNode представляет узел в графе символов
type SymbolNode struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       SymbolType        `json:"type"`
	Path       string            `json:"path"`
	Line       int               `json:"line"`
	Column     int               `json:"column"`
	Package    string            `json:"package,omitempty"`
	Visibility Visibility        `json:"visibility"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// SymbolType определяет тип символа
type SymbolType string

const (
	SymbolTypeFunction  SymbolType = "function"
	SymbolTypeMethod    SymbolType = "method"
	SymbolTypeStruct    SymbolType = "struct"
	SymbolTypeInterface SymbolType = "interface"
	SymbolTypeVariable  SymbolType = "variable"
	SymbolTypeConstant  SymbolType = "constant"
	SymbolTypePackage   SymbolType = "package"
	SymbolTypeImport    SymbolType = "import"
	SymbolTypeType      SymbolType = "type"
	SymbolTypeField     SymbolType = "field"
)

// Visibility определяет видимость символа
type Visibility string

const (
	VisibilityPublic    Visibility = "public"
	VisibilityPrivate   Visibility = "private"
	VisibilityProtected Visibility = "protected"
)

// SymbolEdge представляет связь между символами
type SymbolEdge struct {
	From     string            `json:"from"`
	To       string            `json:"to"`
	Type     EdgeType          `json:"type"`
	Weight   float64           `json:"weight"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// EdgeType определяет тип связи
type EdgeType string

const (
	EdgeTypeCalls      EdgeType = "calls"
	EdgeTypeImports    EdgeType = "imports"
	EdgeTypeExtends    EdgeType = "extends"
	EdgeTypeImplements EdgeType = "implements"
	EdgeTypeUses       EdgeType = "uses"
	EdgeTypeReferences EdgeType = "references"
	EdgeTypeDependsOn  EdgeType = "depends_on"
)

// SymbolGraph представляет граф символов проекта
type SymbolGraph struct {
	Nodes []*SymbolNode `json:"nodes"`
	Edges []*SymbolEdge `json:"edges"`
}

// SymbolGraphBuilder определяет интерфейс для построения графа символов
type SymbolGraphBuilder interface {
	// BuildGraph строит граф символов для указанного проекта
	BuildGraph(ctx context.Context, projectRoot string) (*SymbolGraph, error)

	// UpdateGraph обновляет граф для измененных файлов
	UpdateGraph(ctx context.Context, projectRoot string, changedFiles []string) (*SymbolGraph, error)

	// GetSuggestions возвращает предложения символов на основе запроса
	GetSuggestions(ctx context.Context, query string, graph *SymbolGraph) ([]*SymbolNode, error)

	// GetDependencies возвращает зависимости для указанного символа
	GetDependencies(ctx context.Context, symbolID string, graph *SymbolGraph) ([]*SymbolNode, error)

	// GetDependents возвращает символы, которые зависят от указанного
	GetDependents(ctx context.Context, symbolID string, graph *SymbolGraph) ([]*SymbolNode, error)
}

// ImportGraph представляет граф импортов
type ImportGraph struct {
	Packages map[string]*PackageNode `json:"packages"`
	Imports  []*ImportEdge           `json:"imports"`
}

// PackageNode представляет пакет в графе импортов
type PackageNode struct {
	Name    string   `json:"name"`
	Path    string   `json:"path"`
	Files   []string `json:"files"`
	Imports []string `json:"imports"`
	Exports []string `json:"exports"`
}

// ImportEdge представляет импорт между пакетами
type ImportEdge struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Type     string `json:"type"` // "direct", "indirect", "test"
	Optional bool   `json:"optional"`
}

// ImportGraphBuilder определяет интерфейс для построения графа импортов
type ImportGraphBuilder interface {
	// BuildImportGraph строит граф импортов для проекта
	BuildImportGraph(ctx context.Context, projectRoot string) (*ImportGraph, error)

	// GetImportPath возвращает путь импорта между двумя пакетами
	GetImportPath(ctx context.Context, from, to string, graph *ImportGraph) ([]string, error)

	// GetCircularImports возвращает циклические импорты
	GetCircularImports(ctx context.Context, graph *ImportGraph) ([][]string, error)
}
