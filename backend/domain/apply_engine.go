package domain

import "context"

// ApplyStrategy определяет стратегию применения правок
type ApplyStrategy string

const (
	ApplyStrategyAnchor   ApplyStrategy = "anchor"
	ApplyStrategyFullFile ApplyStrategy = "fullFile"
	ApplyStrategyAST      ApplyStrategy = "ast"
	ApplyStrategyRecipe   ApplyStrategy = "recipe"
)

// EditType represents the type of edit operation
type EditType string

const (
	EditTypeReplace EditType = "replace"
	EditTypeInsert  EditType = "insert"
	EditTypeDelete  EditType = "delete"
)

// ApplyOperation определяет операцию применения
type ApplyOperation struct {
	ID           string                 `json:"id"`
	Path         string                 `json:"path"`
	Language     string                 `json:"language"`
	Strategy     ApplyStrategy          `json:"strategy"`
	Operation    string                 `json:"operation"` // "modify", "create", "delete", "move"
	Content      string                 `json:"content,omitempty"`
	AnchorBefore string                 `json:"anchorBefore,omitempty"`
	AnchorAfter  string                 `json:"anchorAfter,omitempty"`
	Hash         string                 `json:"hash,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ApplyResult представляет результат применения правки
type ApplyResult struct {
	Success      bool                   `json:"success"`
	Path         string                 `json:"path"`
	OperationID  string                 `json:"operationId"`
	Error        string                 `json:"error,omitempty"`
	AppliedLines int                    `json:"appliedLines"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

// ApplyEngine определяет интерфейс для применения правок
type ApplyEngine interface {
	// ApplyOperation применяет одну операцию
	ApplyOperation(ctx context.Context, op *ApplyOperation) (*ApplyResult, error)

	// ApplyOperations применяет несколько операций
	ApplyOperations(ctx context.Context, ops []*ApplyOperation) ([]*ApplyResult, error)

	// ApplyEdit применяет одну правку
	ApplyEdit(ctx context.Context, edit Edit) error

	// ValidateOperation проверяет корректность операции
	ValidateOperation(ctx context.Context, op *ApplyOperation) error

	// RollbackOperation откатывает операцию
	RollbackOperation(ctx context.Context, result *ApplyResult) error

	// RegisterFormatter регистрирует форматтер для языка
	RegisterFormatter(language string, formatter Formatter)

	// RegisterImportFixer регистрирует исправитель импортов для языка
	RegisterImportFixer(language string, fixer ImportFixer)
}

// Formatter определяет интерфейс для форматирования кода
type Formatter interface {
	// FormatFile форматирует файл
	FormatFile(ctx context.Context, path string) error

	// FormatContent форматирует содержимое
	FormatContent(ctx context.Context, content string, language string) (string, error)

	// GetSupportedLanguages возвращает поддерживаемые языки
	GetSupportedLanguages() []string
}

// ImportFixer определяет интерфейс для исправления импортов
type ImportFixer interface {
	// FixImports исправляет импорты в файле
	FixImports(ctx context.Context, path string) error

	// FixImportsInContent исправляет импорты в содержимом
	FixImportsInContent(ctx context.Context, content string, language string) (string, error)

	// GetSupportedLanguages возвращает поддерживаемые языки
	GetSupportedLanguages() []string
}

// ApplyEngineConfig определяет конфигурацию движка применения
type ApplyEngineConfig struct {
	AutoFormat     bool     `json:"autoFormat"`
	AutoFixImports bool     `json:"autoFixImports"`
	BackupFiles    bool     `json:"backupFiles"`
	ValidateAfter  bool     `json:"validateAfter"`
	Languages      []string `json:"languages"`
}

// EditsJSON представляет структуру правок
type EditsJSON struct {
	SchemaVersion    string         `json:"schemaVersion"`
	ToolchainVersion string         `json:"toolchainVersion"`
	Metadata         *EditsMetadata `json:"metadata"`
	Edits            []*Edit        `json:"edits"`
}

// EditsMetadata представляет метаданные правок
type EditsMetadata struct {
	Reason          string  `json:"reason"`
	TaskID          string  `json:"taskId"`
	StepID          string  `json:"stepId"`
	Confidence      float64 `json:"confidence"`
	EstimatedImpact string  `json:"estimatedImpact"`
}

// Edit представляет одну правку
type Edit struct {
	ID          string                 `json:"id"`
	AtomicGroup string                 `json:"atomicGroup,omitempty"`
	DependsOn   []string               `json:"dependsOn,omitempty"`
	Kind        string                 `json:"kind"`
	Op          string                 `json:"op"`
	Path        string                 `json:"path"`
	FilePath    string                 `json:"filePath"`
	Language    string                 `json:"language"`
	Content     string                 `json:"content,omitempty"`
	Anchor      interface{}            `json:"anchor,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`

	// Fields used by internal modification package
	Type       EditType `json:"-"` // Not serialized to JSON
	OldContent string   `json:"-"` // Not serialized to JSON
	NewContent string   `json:"-"` // Not serialized to JSON
	Position   int      `json:"-"` // Not serialized to JSON
}
