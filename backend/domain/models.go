package domain

import (
	"context"
	"time"
)

type FileNode struct {
	Name            string      `json:"name"`
	Path            string      `json:"path"`
	RelPath         string      `json:"relPath"`
	IsDir           bool        `json:"isDir"`
	Size            int64       `json:"size"`
	ContentType     string      `json:"contentType"` // "text", "binary", "unknown"
	Children        []*FileNode `json:"children,omitempty"`
	IsGitignored    bool        `json:"isGitignored"`
	IsCustomIgnored bool        `json:"isCustomIgnored"`
	IsIgnored       bool        `json:"isIgnored"`
}

type FileStatus struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

type Commit struct {
	Hash    string `json:"hash"`
	Subject string `json:"subject"`
}

type CommitWithFiles struct {
	Hash    string   `json:"hash"`
	Subject string   `json:"subject"`
	Author  string   `json:"author"`
	Date    string   `json:"date"`
	Files   []string `json:"files"`
	IsMerge bool     `json:"isMerge"`
}

// CommitInfo represents a simplified commit for selection UI
type CommitInfo struct {
	Hash    string `json:"hash"`
	Subject string `json:"subject"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// TokenCounter defines a function type for counting tokens.
type TokenCounter func(text string) int

type ContextSummaryInfo struct {
	ID            string
	FileCount     int
	TotalFiles    int
	TotalSize     int64
	TotalLines    int
	TokenCount    int
	LineCount     int
	TotalChunks   int
	LanguageStats map[string]int
	LastModified  time.Time
	GitRepo       bool
	BuildSystem   string
	Frameworks    []string
	HasTests      bool
	HasDockerfile bool
	HasCICD       bool
}

type ParsedDiff struct {
	FileDiffs []FileDiff `json:"fileDiffs"`
}

type FileDiff struct {
	FilePath string `json:"filePath"`
	Hunks    []Hunk `json:"hunks"`
}

type Hunk struct {
	Header string   `json:"header"`
	Lines  []string `json:"lines"`
}

// AutonomousTaskRequest запрос на запуск автономной задачи
type AutonomousTaskRequest struct {
	Task        string                `json:"task"`
	SlaPolicy   string                `json:"slaPolicy"`
	ProjectPath string                `json:"projectPath"`
	Options     AutonomousTaskOptions `json:"options"`
}

// AutonomousTaskOptions опции автономной задачи
type AutonomousTaskOptions struct {
	MaxTokens            int     `json:"maxTokens"`
	Temperature          float64 `json:"temperature"`
	EnableStaticAnalysis bool    `json:"enableStaticAnalysis"`
	EnableTests          bool    `json:"enableTests"`
	EnableSBOM           bool    `json:"enableSBOM"`
}

// AutonomousTaskResponse ответ на запуск автономной задачи
type AutonomousTaskResponse struct {
	TaskId  string `json:"taskId"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// AutonomousTaskStatus статус автономной задачи
type AutonomousTaskStatus struct {
	TaskId                 string    `json:"taskId"`
	Status                 string    `json:"status"`
	CurrentStep            string    `json:"currentStep"`
	Progress               float64   `json:"progress"`
	EstimatedTimeRemaining int64     `json:"estimatedTimeRemaining"`
	StartedAt              time.Time `json:"startedAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	Error                  string    `json:"error"`
}

// GenericReport общий отчет
type GenericReport struct {
	Id        string    `json:"id"`
	TaskId    string    `json:"taskId"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ContextBuildOptions опции для построения контекста
type ContextBuildOptions struct {
	StripComments        bool   `json:"stripComments"`
	IncludeManifest      bool   `json:"includeManifest"`
	IncludeLineNumbers   bool   `json:"includeLineNumbers"`
	MaxTokens            int    `json:"maxTokens"`
	MaxMemoryMB          int    `json:"maxMemoryMB"`
	IncludeTests         bool   `json:"includeTests"`
	SplitStrategy        string `json:"splitStrategy"`
	ForceStream          bool   `json:"forceStream"`
	EnableProgressEvents bool   `json:"enableProgressEvents"`
	OutputFormat         string `json:"outputFormat"` // markdown, xml, json, plain

	// Content optimization options
	ExcludeTests       bool `json:"excludeTests"`       // Исключать тестовые файлы из контекста
	CollapseEmptyLines bool `json:"collapseEmptyLines"` // Схлопывать множественные пустые строки
	StripLicense       bool `json:"stripLicense"`       // Удалять лицензионные заголовки
	CompactDataFiles   bool `json:"compactDataFiles"`   // Сжимать JSON/YAML файлы
	SkeletonMode       bool `json:"skeletonMode"`       // Генерировать только скелет кода (AST-based)
	TrimWhitespace     bool `json:"trimWhitespace"`     // Удалять trailing whitespace
}

// Context представляет контекст проекта
type Context struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Files       []string  `json:"files"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ProjectPath string    `json:"projectPath"`
	TokenCount  int       `json:"tokenCount"`
	TotalLines  int64     `json:"totalLines"`
	TotalChars  int64     `json:"totalChars"`
}

// CRITICAL OOM FIX: ContextSummary replaces full Context content to prevent memory issues
// This lightweight summary contains metadata without storing large text content
type ContextSummary struct {
	ID          string          `json:"id"`
	ProjectPath string          `json:"projectPath"`
	FileCount   int             `json:"fileCount"`
	TotalSize   int64           `json:"totalSize"`
	TokenCount  int             `json:"tokenCount"`
	LineCount   int             `json:"lineCount"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	Status      string          `json:"status"`
	Metadata    ContextMetadata `json:"metadata"`
}

// ContextMetadata contains additional context information
type ContextMetadata struct {
	BuildDuration  int64                `json:"buildDuration"`
	LastModified   time.Time            `json:"lastModified"`
	SelectedFiles  []string             `json:"selectedFiles"`
	BuildOptions   *ContextBuildOptions `json:"buildOptions,omitempty"`
	Warnings       []string             `json:"warnings"`
	Errors         []string             `json:"errors"`
	ContentPath    string               `json:"contentPath"`
	SkippedFiles   []string             `json:"skippedFiles,omitempty"`
	SkippedReasons map[string]string    `json:"skippedReasons,omitempty"`
	ProjectPath    string               `json:"projectPath,omitempty"`
	ChunksCount    int                  `json:"chunksCount,omitempty"`
	SplitStrategy  string               `json:"splitStrategy,omitempty"`
}

// ContextChunk represents a paginated piece of context content
// This allows viewing context content without loading it all into memory
type ContextChunk struct {
	Lines     []string `json:"lines"`
	StartLine int      `json:"startLine"`
	EndLine   int      `json:"endLine"`
	HasMore   bool     `json:"hasMore"`
	ChunkID   string   `json:"chunkId"`
	ContextID string   `json:"contextId"`
}

// ContextStream represents a streaming context for large projects
// This allows working with contexts that are too large to fit in memory
type ContextStream struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	ProjectPath string   `json:"projectPath"`
	TotalLines  int64    `json:"totalLines"`
	TotalChars  int64    `json:"totalChars"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	TokenCount  int      `json:"tokenCount"`
}

// ContextLineRange represents a range of lines from a context
// This allows retrieving specific line ranges without loading the entire context
type ContextLineRange struct {
	StartLine int64    `json:"startLine"`
	EndLine   int64    `json:"endLine"`
	Lines     []string `json:"lines"`
}

// SuggestedFile represents a file suggested by the AI analyzer
type SuggestedFile struct {
	Path       string  `json:"path"`
	Reason     string  `json:"reason"`
	Confidence float64 `json:"confidence"`
}

// SLAPolicy определяет SLA политику для автономных задач
type SLAPolicy struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MaxTokens   int    `json:"maxTokens"`
	MaxFiles    int    `json:"maxFiles"`
	MaxTime     int64  `json:"maxTime"`   // seconds
	MaxMemory   int64  `json:"maxMemory"` // bytes
	MaxRetries  int    `json:"maxRetries"`
	Timeout     int64  `json:"timeout"` // seconds
}

// AutonomousTask представляет автономную задачу
type AutonomousTask struct {
	ID          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Status      string                `json:"status"`
	ProjectPath string                `json:"projectPath"`
	SlaPolicy   string                `json:"slaPolicy"`
	Options     AutonomousTaskOptions `json:"options"`
	CreatedAt   time.Time             `json:"createdAt"`
	UpdatedAt   time.Time             `json:"updatedAt"`
	CompletedAt *time.Time            `json:"completedAt,omitempty"`
	Progress    float64               `json:"progress"`
	Error       string                `json:"error,omitempty"`
}

// LogEntry представляет запись в логе
type LogEntry struct {
	ID        string                 `json:"id"`
	TaskID    string                 `json:"taskId"`
	Level     string                 `json:"level"` // INFO, WARN, ERROR, DEBUG
	Message   string                 `json:"message"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// ReportRepository определяет интерфейс для работы с отчетами
type ReportRepository interface {
	// GetReport retrieves a report by ID
	GetReport(ctx context.Context, reportID string) (*GenericReport, error)

	// ListReports lists all reports, optionally filtered by type
	ListReports(ctx context.Context, reportType string) ([]*GenericReport, error)

	// SaveReport saves a report
	SaveReport(ctx context.Context, report *GenericReport) error

	// DeleteReport deletes a report by ID
	DeleteReport(ctx context.Context, reportID string) error

	// GetReportsByTask retrieves all reports for a specific task
	GetReportsByTask(ctx context.Context, taskID string) ([]*GenericReport, error)
}

// LLMClient определяет интерфейс для работы с LLM
type LLMClient interface {
	// HealthCheck проверяет доступность LLM сервиса
	HealthCheck(ctx context.Context) error

	// GenerateWithGBNF генерирует ответ с использованием GBNF грамматики
	GenerateWithGBNF(ctx context.Context, prompt string, grammar string, options map[string]interface{}) (*LlamaCppResponse, error)

	// GenerateEditsJSON генерирует Edits JSON с помощью GBNF грамматики
	GenerateEditsJSON(ctx context.Context, prompt string, options map[string]interface{}) ([]byte, error)

	// GetModelInfo получает информацию о модели
	GetModelInfo(ctx context.Context) (map[string]interface{}, error)
}

// LlamaCppResponse represents the response from llama.cpp server
type LlamaCppResponse struct {
	Content    string `json:"content"`
	Stop       bool   `json:"stop"`
	StopReason string `json:"stop_reason"`
	Timings    struct {
		PredN    int     `json:"pred_n"`
		PredMS   float64 `json:"pred_ms"`
		PromptN  int     `json:"prompt_n"`
		PromptMS float64 `json:"prompt_ms"`
	} `json:"timings"`
	TokensEvaluated int  `json:"tokens_evaluated"`
	TokensPredicted int  `json:"tokens_predicted"`
	Truncated       bool `json:"truncated"`
}

// LLMConfig конфигурация для LLM клиента
type LLMConfig struct {
	BaseURL       string        `json:"base_url"`
	Timeout       time.Duration `json:"timeout"`
	MaxTokens     int           `json:"max_tokens"`
	Temperature   float64       `json:"temperature"`
	TopP          float64       `json:"top_p"`
	TopK          int           `json:"top_k"`
	RepeatPenalty float64       `json:"repeat_penalty"`
}

// FileReader определяет интерфейс для чтения файлов
type FileReader interface {
	// ReadFile reads the content of a file
	ReadFile(filename string) ([]byte, error)
}

// SBOMGenerator определяет интерфейс для генерации SBOM
type SBOMGenerator interface {
	// IsAvailable проверяет доступность инструмента
	IsAvailable() bool

	// GenerateSBOM генерирует SBOM для проекта
	GenerateSBOM(ctx context.Context, projectPath string, format SBOMFormat) (*SBOMResult, error)

	// ValidateSBOM валидирует SBOM файл
	ValidateSBOM(ctx context.Context, sbomPath string, format SBOMFormat) error
}

// VulnerabilityScanner определяет интерфейс для сканирования уязвимостей
type VulnerabilityScanner interface {
	// IsAvailable проверяет доступность инструмента
	IsAvailable() bool

	// ScanVulnerabilities сканирует уязвимости в проекте
	ScanVulnerabilities(ctx context.Context, projectPath string) (*VulnerabilityScanResult, error)
}

// LicenseScanner определяет интерфейс для сканирования лицензий
type LicenseScanner interface {
	// IsAvailable проверяет доступность инструмента
	IsAvailable() bool

	// ScanLicenses сканирует лицензии в проекте
	ScanLicenses(ctx context.Context, projectPath string) (*LicenseScanResult, error)
}

// Task Protocol Verification System Models

// ProtocolStage represents different stages of the verification protocol
type ProtocolStage string

const (
	StageAnalysis   ProtocolStage = "analysis"
	StageLinting    ProtocolStage = "linting"
	StageBuilding   ProtocolStage = "building"
	StageTesting    ProtocolStage = "testing"
	StageGuardrails ProtocolStage = "guardrails"
)

// ErrorType represents different types of verification errors
type ErrorType string

const (
	ErrorTypeLinting     ErrorType = "linting"
	ErrorTypeCompilation ErrorType = "compilation"
	ErrorTypeTypeCheck   ErrorType = "typecheck"
	ErrorTypeTesting     ErrorType = "testing"
	ErrorTypeGuardrail   ErrorType = "guardrail"
	ErrorTypeDependency  ErrorType = "dependency"
	ErrorTypeSyntax      ErrorType = "syntax"
	ErrorTypeImport      ErrorType = "import"
	ErrorTypeLogic       ErrorType = "logic"
)

// CorrectionAction represents different types of correction actions
type CorrectionAction string

const (
	ActionFixImport      CorrectionAction = "fix_import"
	ActionFixSyntax      CorrectionAction = "fix_syntax"
	ActionFixType        CorrectionAction = "fix_type"
	ActionAddMissingCode CorrectionAction = "add_missing_code"
	ActionRemoveCode     CorrectionAction = "remove_code"
	ActionUpdateTest     CorrectionAction = "update_test"
	ActionFormatCode     CorrectionAction = "format_code"
)

// TaskProtocolConfig represents configuration for the verification protocol
type TaskProtocolConfig struct {
	ProjectPath    string                   `json:"projectPath"`
	Languages      []string                 `json:"languages"`
	EnabledStages  []ProtocolStage          `json:"enabledStages"`
	MaxRetries     int                      `json:"maxRetries"`
	FailFast       bool                     `json:"failFast"`
	SelfCorrection SelfCorrectionConfig     `json:"selfCorrection"`
	Timeouts       map[string]time.Duration `json:"timeouts"`
}

// SelfCorrectionConfig represents configuration for self-correction capabilities
type SelfCorrectionConfig struct {
	Enabled      bool `json:"enabled"`
	MaxAttempts  int  `json:"maxAttempts"`
	AIAssistance bool `json:"aiAssistance"`
}

// TaskProtocolResult represents the result of protocol verification
type TaskProtocolResult struct {
	TaskID           string                 `json:"taskId"`
	Success          bool                   `json:"success"`
	StartedAt        time.Time              `json:"startedAt"`
	CompletedAt      time.Time              `json:"completedAt"`
	Stages           []*ProtocolStageResult `json:"stages"`
	CorrectionCycles int                    `json:"correctionCycles"`
	FinalError       string                 `json:"finalError,omitempty"`
}

// ProtocolStageResult represents the result of a single protocol stage
type ProtocolStageResult struct {
	Stage           ProtocolStage     `json:"stage"`
	Success         bool              `json:"success"`
	Attempts        int               `json:"attempts"`
	Duration        time.Duration     `json:"duration"`
	ErrorDetails    *ErrorDetails     `json:"errorDetails,omitempty"`
	CorrectionSteps []*CorrectionStep `json:"correctionSteps,omitempty"`
}

// ErrorDetails provides detailed information about verification errors
type ErrorDetails struct {
	Stage       ProtocolStage `json:"stage"`
	ErrorType   ErrorType     `json:"errorType"`
	Message     string        `json:"message"`
	SourceFile  string        `json:"sourceFile,omitempty"`
	LineNumber  int           `json:"lineNumber,omitempty"`
	Column      int           `json:"column,omitempty"`
	Tool        string        `json:"tool"` // staticcheck, eslint, go build, etc.
	Severity    string        `json:"severity"`
	Suggestions []string      `json:"suggestions"`
}

// CorrectionStep represents a single correction action
type CorrectionStep struct {
	Action      CorrectionAction `json:"action"`
	Target      string           `json:"target"` // file path, function name, etc.
	Description string           `json:"description"`
	Applied     bool             `json:"applied"`
	Result      string           `json:"result,omitempty"`
}

// CorrectionResult represents the outcome of applying a correction
type CorrectionResult struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	FilesChanged []string `json:"filesChanged"`
}

// CorrectionGuidance provides AI-generated guidance for error correction
type CorrectionGuidance struct {
	Error       *ErrorDetails     `json:"error"`
	Steps       []*CorrectionStep `json:"steps"`
	Explanation string            `json:"explanation"`
	Confidence  float64           `json:"confidence"`
}

// TaskContext provides context information for task protocol operations
type TaskContext struct {
	TaskID      string            `json:"taskId"`
	ProjectPath string            `json:"projectPath"`
	Languages   []string          `json:"languages"`
	Files       []string          `json:"files"`
	Metadata    map[string]string `json:"metadata"`
}
