package domain

import "context"

// ContextPack представляет упакованный контекст для задачи
type ContextPack struct {
	PackVersion string               `json:"packVersion"`
	Target      *ContextTarget       `json:"target"`
	Snippets    []*ContextSnippet    `json:"snippets"`
	Deps        []*ContextDependency `json:"deps,omitempty"`
	Build       *ContextBuild        `json:"build,omitempty"`
	Tests       []*ContextTest       `json:"tests,omitempty"`
	Constraints *ContextConstraints  `json:"constraints"`
	Provenance  *ContextProvenance   `json:"provenance"`
}

// ContextTarget определяет цель контекста
type ContextTarget struct {
	Lang   string `json:"lang"`
	File   string `json:"file,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}

// ContextSnippet представляет фрагмент кода
type ContextSnippet struct {
	Path      string `json:"path"`
	Kind      string `json:"kind"` // "header", "imports", "body", "match"
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
	TextHash  string `json:"textHash"`
	Content   string `json:"content,omitempty"`
}

// ContextDependency представляет зависимость
type ContextDependency struct {
	Pkg     string `json:"pkg"`
	Version string `json:"version,omitempty"`
	License string `json:"license,omitempty"`
	Path    string `json:"path,omitempty"`
}

// ContextBuild представляет информацию о сборке
type ContextBuild struct {
	Toolchain string            `json:"toolchain"`
	Commands  []string          `json:"commands"`
	Env       map[string]string `json:"env,omitempty"`
}

// ContextTest представляет тест
type ContextTest struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

// ContextConstraints представляет ограничения
type ContextConstraints struct {
	MaxFiles   int    `json:"maxFiles"`
	MaxLines   int    `json:"maxLines"`
	TimeBudget string `json:"timeBudget,omitempty"`
}

// ContextProvenance представляет происхождение контекста
type ContextProvenance struct {
	IndexedAt    string `json:"indexedAt"`
	IndexID      string `json:"indexId"`
	SourceMirror string `json:"sourceMirror,omitempty"`
}

// ContextPackBuilder определяет интерфейс для построения упакованного контекста
type ContextPackBuilder interface {
	// BuildPack строит упакованный контекст для задачи
	BuildPack(ctx context.Context, task string, projectRoot string, options *ContextPackOptions) (*ContextPack, error)

	// UpdatePack обновляет упакованный контекст для измененных файлов
	UpdatePack(ctx context.Context, pack *ContextPack, changedFiles []string) (*ContextPack, error)

	// ValidatePack проверяет корректность упакованного контекста
	ValidatePack(ctx context.Context, pack *ContextPack) error
}

// ContextPackOptions определяет опции для построения контекста
type ContextPackOptions struct {
	MaxFiles        int      `json:"maxFiles"`
	MaxLines        int      `json:"maxLines"`
	IncludeTests    bool     `json:"includeTests"`
	IncludeDeps     bool     `json:"includeDeps"`
	IncludeBuild    bool     `json:"includeBuild"`
	SmartSnippets   bool     `json:"smartSnippets"`
	Deterministic   bool     `json:"deterministic"`
	Language        string   `json:"language"`
	TargetFile      string   `json:"targetFile,omitempty"`
	TargetSymbol    string   `json:"targetSymbol,omitempty"`
	ExcludePatterns []string `json:"excludePatterns,omitempty"`
}

// ContextPackSerializer определяет интерфейс для сериализации контекста
type ContextPackSerializer interface {
	// SerializePack сериализует упакованный контекст
	SerializePack(ctx context.Context, pack *ContextPack) ([]byte, error)

	// DeserializePack десериализует упакованный контекст
	DeserializePack(ctx context.Context, data []byte) (*ContextPack, error)

	// GetFormat возвращает формат сериализации
	GetFormat() string
}
