package domain

import (
	"context"
	"io/fs"
	"time"
)

// Logger определяет интерфейс для логирования
type Logger interface {
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
	Fatal(message string)
}

// NoopLogger provides a dummy implementation of domain.Logger for cases where a logger
// is required but full logging infrastructure might not be available or desired.
type NoopLogger struct{}

func (l *NoopLogger) Debug(message string)   {}
func (l *NoopLogger) Info(message string)    {}
func (l *NoopLogger) Warning(message string) {}
func (l *NoopLogger) Error(message string)   {}
func (l *NoopLogger) Fatal(message string)   {}

// EventBus определяет интерфейс для событийной шины
type EventBus interface {
	Emit(eventName string, data ...interface{})
}

// TreeBuilder определяет интерфейс для построения дерева файлов
type TreeBuilder interface {
	BuildTree(dirPath string, useGitignore bool, useCustomIgnore bool) ([]*FileNode, error)
	InvalidateCache()
}

// FileContentReader определяет интерфейс для чтения содержимого файлов
type FileContentReader interface {
	ReadContents(
		ctx context.Context,
		filePaths []string,
		rootDir string,
		progress func(current, total int64),
	) (map[string]string, error)
}

// GitRepository определяет интерфейс для работы с Git
type GitRepository interface {
	GetUncommittedFiles(projectRoot string) ([]FileStatus, error)
	GetRichCommitHistory(projectRoot, branchName string, limit int) ([]CommitWithFiles, error)
	GetFileContentAtCommit(projectRoot, filePath, commitHash string) (string, error)
	GetGitignoreContent(projectRoot string) (string, error)
	IsGitAvailable() bool
	GetBranches(projectRoot string) ([]string, error)
	GetCurrentBranch(projectRoot string) (string, error)
	GetAllFiles(projectPath string) ([]string, error)
	GenerateDiff(projectPath string) (string, error)
	// New methods for remote/branch context building
	IsGitRepository(projectPath string) bool
	CloneRepository(url, targetPath string, depth int) error
	CheckoutBranch(projectPath, branch string) error
	CheckoutCommit(projectPath, commitHash string) error
	GetCommitHistory(projectPath string, limit int) ([]CommitInfo, error)
	FetchRemoteBranches(projectPath string) ([]string, error)
	// Read files at specific ref without checkout
	ListFilesAtRef(projectPath, ref string) ([]string, error)
	GetFileAtRef(projectPath, filePath, ref string) (string, error)
}

// SettingsRepository определяет интерфейс для работы с настройками
type SettingsRepository interface {
	GetCustomIgnoreRules() string
	SetCustomIgnoreRules(rules string)
	GetCustomPromptRules() string
	SetCustomPromptRules(rules string)
	GetOpenAIKey() string
	SetOpenAIKey(key string)
	GetGeminiKey() string
	SetGeminiKey(key string)
	GetOpenRouterKey() string
	SetOpenRouterKey(key string)
	GetLocalAIKey() string
	SetLocalAIKey(key string)
	GetLocalAIHost() string
	SetLocalAIHost(host string)
	GetLocalAIModelName() string
	SetLocalAIModelName(name string)
	GetQwenKey() string
	SetQwenKey(key string)
	GetQwenHost() string
	SetQwenHost(host string)
	GetSelectedAIProvider() string
	SetSelectedAIProvider(provider string)
	GetSelectedModel(provider string) string
	SetSelectedModel(provider, model string)
	GetModels(provider string) []string
	SetModels(provider string, models []string)
	GetUseGitignore() bool
	SetUseGitignore(use bool)
	GetUseCustomIgnore() bool
	SetUseCustomIgnore(use bool)
	GetRecentProjects() []RecentProjectInfo
	AddRecentProject(path, name string)
	RemoveRecentProject(path string)

	Save() error
	GetSettingsDTO() (SettingsDTO, error) // Added as per compilation error
}

// FileSystemWatcher определяет интерфейс для отслеживания изменений файловой системы
type FileSystemWatcher interface {
	Start(rootPath string) error
	Stop()
	RefreshAndRescan() error
}

// ContextSplitter определяет интерфейс для разбиения большого контекста на части.
type ContextSplitter interface {
	// SplitContext разбивает текст контекста на части, учитывая лимиты токенов и стратегии.
	SplitContext(ctxText string, settings SplitSettings) ([]string, error)
}

// CommentStripper определяет интерфейс для удаления комментариев из исходного кода
type CommentStripper interface {
	// Strip удаляет комментарии из содержимого файла, опираясь на расширение filePath
	Strip(content string, filePath string) string
}

// ContentOptimizer определяет интерфейс для оптимизации контента файлов для AI-контекста
type ContentOptimizer interface {
	// Optimize применяет оптимизации к контенту файла
	Optimize(ctx context.Context, content, filePath string, opts ContentOptimizeOptions) string

	// OptimizeWithDefaults применяет оптимизации с настройками по умолчанию
	OptimizeWithDefaults(ctx context.Context, content, filePath string) string

	// CanGenerateSkeleton проверяет возможность генерации скелета для файла
	CanGenerateSkeleton(filePath string) bool
}

// ContentOptimizeOptions опции оптимизации контента для AI
type ContentOptimizeOptions struct {
	CollapseEmptyLines bool `json:"collapseEmptyLines"` // Схлопывать множественные пустые строки
	StripLicense       bool `json:"stripLicense"`       // Удалять лицензионные заголовки
	StripComments      bool `json:"stripComments"`      // Удалять комментарии
	CompactDataFiles   bool `json:"compactDataFiles"`   // Сжимать JSON/YAML файлы
	SkeletonMode       bool `json:"skeletonMode"`       // Генерировать только скелет кода (AST-based)
	TrimWhitespace     bool `json:"trimWhitespace"`     // Удалять trailing whitespace
}

// AIProviderFactory is a function type that creates an AIProvider.
type AIProviderFactory func(providerType, apiKey string) (AIProvider, error)

// ModelFetcher is a function type that fetches available models for a provider.
type ModelFetcher func(apiKey string) ([]string, error)

// ModelFetcherRegistry is a map of provider types to their model fetchers.
type ModelFetcherRegistry map[string]ModelFetcher

// KeyResolver defines the interface for resolving API keys for different providers
type KeyResolver interface {
	// GetKey retrieves the API key for a given provider type from settings
	GetKey(providerType string, settings SettingsDTO) (string, error)
}

// PathProvider определяет интерфейс для работы с путями файловой системы
type PathProvider interface {
	// Join соединяет элементы пути
	Join(elem ...string) string

	// Base возвращает последний элемент пути
	Base(path string) string

	// Dir возвращает все элементы пути кроме последнего
	Dir(path string) string

	// IsAbs проверяет, является ли путь абсолютным
	IsAbs(path string) bool

	// Clean возвращает очищенную версию пути
	Clean(path string) string

	// Getwd возвращает текущую рабочую директорию
	Getwd() (string, error)
}

// FileSystemWriter определяет интерфейс для записи файлов
type FileSystemWriter interface {
	// WriteFile записывает данные в файл
	WriteFile(filename string, data []byte, perm int) error

	// MkdirAll создает директорию вместе со всеми родительскими директориями
	MkdirAll(path string, perm int) error

	// Remove удаляет файл или директорию
	Remove(name string) error

	// RemoveAll удаляет файл или директорию и все её содержимое
	RemoveAll(path string) error
}

// OPAService определяет интерфейс для работы с OPA политиками
type OPAService interface {
	// ValidatePath проверяет путь через OPA политики
	ValidatePath(path string) (*OPAValidationResult, error)

	// ValidateBudget проверяет бюджет через OPA политики
	ValidateBudget(budgetType string, current, limit int64) (*OPAValidationResult, error)

	// ValidateTask проверяет задачу через OPA политики
	ValidateTask(taskID string, files []string, linesChanged int64, ephemeralMode bool) (*OPAValidationResult, error)

	// ValidateConfig проверяет конфигурацию через OPA политики
	ValidateConfig(config GuardrailConfig) (*OPAValidationResult, error)
}

// FileStatProvider определяет интерфейс для получения информации о файлах
type FileStatProvider interface {
	// Stat возвращает информацию о файле
	Stat(name string) (FileInfo, error)
}

// FileInfo описывает файл и возвращается Stat
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}

// FileMode represents a file's mode and permission bits.
type FileMode = fs.FileMode

// TempFileProvider определяет интерфейс для работы с временными файлами
type TempFileProvider interface {
	// MkdirTemp создает временную директорию
	MkdirTemp(dir, pattern string) (string, error)
}

// CommandRunner определяет интерфейс для выполнения команд
type CommandRunner interface {
	// RunCommand выполняет команду с заданным контекстом и аргументами
	RunCommand(ctx context.Context, name string, args ...string) ([]byte, error)

	// RunCommandInDir выполняет команду в указанной директории
	RunCommandInDir(ctx context.Context, dir, name string, args ...string) ([]byte, error)
}

// Task Protocol Verification System Interfaces

// TaskProtocolService defines the interface for task protocol verification
type TaskProtocolService interface {
	// ExecuteProtocol executes the full verification protocol for a task
	ExecuteProtocol(ctx context.Context, config *TaskProtocolConfig) (*TaskProtocolResult, error)

	// ValidateStage executes a single verification stage
	ValidateStage(ctx context.Context, stage ProtocolStage, config *TaskProtocolConfig) (*ProtocolStageResult, error)

	// RequestCorrectionGuidance requests AI-generated correction guidance for errors
	RequestCorrectionGuidance(ctx context.Context, error *ErrorDetails, context *TaskContext) (*CorrectionGuidance, error)
}

// ErrorAnalyzer defines the interface for analyzing verification errors
type ErrorAnalyzer interface {
	// AnalyzeError analyzes error output and provides detailed error information
	AnalyzeError(errorOutput string, stage ProtocolStage) (*ErrorDetails, error)

	// SuggestCorrections provides correction suggestions for a given error
	SuggestCorrections(error *ErrorDetails) ([]*CorrectionStep, error)

	// ClassifyErrorType determines the type of error from output
	ClassifyErrorType(errorOutput string) ErrorType
}

// CorrectionEngine defines the interface for applying corrections
type CorrectionEngine interface {
	// ApplyCorrection applies a single correction step
	ApplyCorrection(ctx context.Context, step *CorrectionStep, projectPath string) (*CorrectionResult, error)

	// ApplyCorrections applies multiple correction steps
	ApplyCorrections(ctx context.Context, steps []*CorrectionStep, projectPath string) (*CorrectionResult, error)

	// CanHandle checks if the engine can handle a specific error type
	CanHandle(error *ErrorDetails) bool
}

// LanguageErrorAnalyzer defines language-specific error analysis
type LanguageErrorAnalyzer interface {
	// AnalyzeError analyzes error output for a specific language
	AnalyzeError(errorOutput string) (*ErrorDetails, error)

	// SuggestCorrections provides language-specific correction suggestions
	SuggestCorrections(error *ErrorDetails) ([]*CorrectionStep, error)

	// ClassifyErrorType determines the error type for the language
	ClassifyErrorType(errorOutput string) ErrorType

	// GetLanguage returns the supported language
	GetLanguage() string
}

// CorrectionRule defines rules for applying corrections
type CorrectionRule interface {
	// CanHandle checks if this rule can handle the given error
	CanHandle(error *ErrorDetails) bool

	// ApplyCorrection applies the correction rule
	ApplyCorrection(error *ErrorDetails, projectPath string) (*CorrectionResult, error)

	// GetPriority returns the priority of this rule (higher = more important)
	GetPriority() int

	// GetErrorTypes returns the error types this rule can handle
	GetErrorTypes() []ErrorType
}

// FileSystemProvider defines file system operations for task protocol
type FileSystemProvider interface {
	// ReadFile reads the contents of a file
	ReadFile(filename string) ([]byte, error)

	// WriteFile writes data to a file
	WriteFile(filename string, data []byte, perm int) error

	// MkdirAll creates directories
	MkdirAll(path string, perm int) error
}

// ContextBuilder определяет интерфейс для построения контекста
type ContextBuilder interface {
	// BuildContext builds a context from project files and returns a ContextSummary to prevent OOM issues
	BuildContext(ctx context.Context, projectPath string, includedPaths []string, options *ContextBuildOptions) (*ContextSummary, error)
}

// ContextFormatOptions options for context formatting
type ContextFormatOptions struct {
	StripComments   bool `json:"stripComments"`
	IncludeManifest bool `json:"includeManifest"`
}

// ContextFormatter defines interface for formatting context output
type ContextFormatter interface {
	// Format formats context string according to specified format (plain, manifest, json, markdown, xml)
	Format(format string, contextContent string, opts ContextFormatOptions) (string, error)
}

// ContextRepository определяет интерфейс для хранения контекста
type ContextRepository interface {
	// SaveContextSummary persists lightweight context metadata on disk
	SaveContextSummary(summary *ContextSummary) error

	// GetContextSummary retrieves persisted context metadata by ID
	GetContextSummary(ctx context.Context, contextID string) (*ContextSummary, error)

	// GetProjectContextSummaries lists context metadata for a project path
	GetProjectContextSummaries(ctx context.Context, projectPath string) ([]*ContextSummary, error)

	// DeleteContext removes context metadata and content from disk
	DeleteContext(ctx context.Context, contextID string) error

	// ReadContextChunk returns a memory-safe chunk of context content
	ReadContextChunk(ctx context.Context, contextID string, startLine int, lineCount int) (*ContextChunk, error)

	// ReadContextContent returns the full context content as a string
	ReadContextContent(ctx context.Context, contextID string) (string, error)
}

// IBuildService defines the interface for build service operations
type IBuildService interface {
	// Build executes project build
	Build(ctx context.Context, projectPath, language string) (*BuildResult, error)

	// TypeCheck performs type checking
	TypeCheck(ctx context.Context, projectPath, language string) (*TypeCheckResult, error)

	// BuildAndTypeCheck executes build and type checking
	BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*BuildResult, *TypeCheckResult, error)

	// BuildMultiLanguage executes build for multiple languages
	BuildMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*BuildResult, error)

	// TypeCheckMultiLanguage executes type checking for multiple languages
	TypeCheckMultiLanguage(ctx context.Context, projectPath string, languages []string) (map[string]*TypeCheckResult, error)

	// ValidateProject performs full project validation
	ValidateProject(ctx context.Context, projectPath string, languages []string) (*ProjectValidationResult, error)

	// GetSupportedLanguages returns supported languages
	GetSupportedLanguages() []string

	// DetectLanguages detects languages in project
	DetectLanguages(ctx context.Context, projectPath string) ([]string, error)
}

// ITestService defines the interface for test service operations
type ITestService interface {
	// RunTests executes tests according to configuration
	RunTests(ctx context.Context, config *TestConfig) ([]*TestResult, error)

	// RunTargetedTests executes targeted tests for affected files
	RunTargetedTests(ctx context.Context, config *TestConfig, changedFiles []string) ([]*TestResult, error)

	// DiscoverTests discovers tests in project
	DiscoverTests(ctx context.Context, projectPath, language string) (*TestSuite, error)

	// BuildAffectedGraph builds affected files graph
	BuildAffectedGraph(ctx context.Context, changedFiles []string, projectPath string) (*AffectedGraph, error)

	// GetTestCoverage gets test coverage
	GetTestCoverage(ctx context.Context, testPath string) (*TestCoverage, error)

	// GetSupportedLanguages returns supported languages
	GetSupportedLanguages() []string

	// RunSmokeTests executes only smoke tests
	RunSmokeTests(ctx context.Context, projectPath, language string) ([]*TestResult, error)

	// RunUnitTests executes only unit tests
	RunUnitTests(ctx context.Context, projectPath, language string) ([]*TestResult, error)

	// RunIntegrationTests executes only integration tests
	RunIntegrationTests(ctx context.Context, projectPath, language string) ([]*TestResult, error)

	// ValidateTestResults validates test results
	ValidateTestResults(results []*TestResult) *TestValidationResult
}

// IStaticAnalyzerService defines the interface for static analyzer service operations
type IStaticAnalyzerService interface {
	// AnalyzeProject analyzes project
	AnalyzeProject(ctx context.Context, projectPath string, languages []string) (*StaticAnalysisReport, error)

	// AnalyzeFile analyzes single file
	AnalyzeFile(ctx context.Context, filePath, language string) (*StaticAnalysisResult, error)

	// GetSupportedAnalyzers returns supported analyzers
	GetSupportedAnalyzers() []StaticAnalyzerType

	// GetAnalyzerForLanguage returns analyzer for language
	GetAnalyzerForLanguage(language string) (StaticAnalyzer, error)

	// ValidateAnalysisResults validates analysis results
	ValidateAnalysisResults(results map[string]*StaticAnalysisResult) *StaticAnalysisValidationResult

	// AnalyzeGoProject analyzes Go project
	AnalyzeGoProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)

	// AnalyzeTypeScriptProject analyzes TypeScript project
	AnalyzeTypeScriptProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)

	// AnalyzeJavaScriptProject analyzes JavaScript project
	AnalyzeJavaScriptProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)

	// AnalyzeJavaProject analyzes Java project
	AnalyzeJavaProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)

	// AnalyzePythonProject analyzes Python project
	AnalyzePythonProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)

	// AnalyzeCppProject analyzes C/C++ project
	AnalyzeCppProject(ctx context.Context, projectPath string) (*StaticAnalysisResult, error)
}

// TaskAnalysis contains analysis results for a task
type TaskAnalysis struct {
	Type         string
	Priority     string
	Technologies []string
	FileTypes    []string
	Keywords     []string
	Reasoning    string
}

// ScoredFile represents a file with its relevance score
type ScoredFile struct {
	RelPath   string  `json:"relPath"`
	Name      string  `json:"name"`
	Size      int64   `json:"size"`
	Relevance float64 `json:"relevance"`
	Reason    string  `json:"reason,omitempty"`
}

// ContextAnalysisResult contains the results of context analysis
type ContextAnalysisResult struct {
	Task            string        `json:"task"`
	TaskType        string        `json:"taskType"`
	Priority        string        `json:"priority"`
	SelectedFiles   []ScoredFile  `json:"selectedFiles"`
	DependencyFiles []*FileNode   `json:"dependencyFiles"`
	Context         string        `json:"context"`
	AnalysisTime    time.Duration `json:"analysisTime"`
	Recommendations []string      `json:"recommendations"`
	EstimatedTokens int           `json:"estimatedTokens"`
	Confidence      float64       `json:"confidence"`
	Reasoning       string        `json:"reasoning,omitempty"`
}

// =============================================================================
// TextUtils Interfaces
// =============================================================================

// LicenseStripper defines interface for removing license headers from source code
type LicenseStripper interface {
	// Strip removes license header from the beginning of file content
	Strip(content string) string

	// StripWithLanguageHint removes license with language hint for optimization
	StripWithLanguageHint(content, ext string) string
}

// TextTokenizer defines interface for token counting and truncation
type TextTokenizer interface {
	// CountTokens returns the estimated token count for text
	CountTokens(text string) int

	// TruncateToTokens truncates text to fit within maxTokens
	TruncateToTokens(text string, maxTokens int) string
}

// SkeletonGenerator defines interface for generating code skeletons
type SkeletonGenerator interface {
	// Generate creates a skeleton representation of code (signatures without bodies)
	Generate(ctx context.Context, content, filePath string) string

	// CanGenerateSkeleton checks if skeleton generation is supported for file
	CanGenerateSkeleton(filePath string) bool
}

// =============================================================================
// Embeddings Interfaces
// =============================================================================

// CodeChunker defines interface for splitting code into chunks for embedding
type CodeChunker interface {
	// ChunkFile splits a file into chunks suitable for embedding
	ChunkFile(filePath string, content []byte, symbols []ChunkSymbolInfo) []CodeChunk
}

// ChunkSymbolInfo represents symbol information for chunking
type ChunkSymbolInfo struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"`
	StartLine int    `json:"startLine"`
	EndLine   int    `json:"endLine"`
}

// =============================================================================
// Reference Finder Interface
// =============================================================================

// ReferenceFinder defines interface for finding symbol references across project
type ReferenceFinder interface {
	// FindReferences finds all references to a symbol in the project
	FindReferences(ctx context.Context, projectRoot string, symbolName string, symbolKind string) ([]SymbolReference, error)

	// FindUsages finds where a symbol is used (excluding definition)
	FindUsages(ctx context.Context, projectRoot string, symbolName string) ([]SymbolReference, error)
}

// SymbolReference represents a reference to a symbol in code
type SymbolReference struct {
	FilePath     string `json:"filePath"`
	Line         int    `json:"line"`
	Column       int    `json:"column"`
	LineText     string `json:"lineText"`
	Context      string `json:"context"`
	IsDefinition bool   `json:"isDefinition"`
}

// =============================================================================
// Git Context Builder Interface
// =============================================================================

// GitContextBuilder defines interface for building context from git history
type GitContextBuilder interface {
	// GetRecentChanges returns files changed recently, sorted by relevance
	GetRecentChanges(since string, pathFilter string) ([]RecentChange, error)

	// GetCoChangedFiles returns files that are often changed together with the given file
	GetCoChangedFiles(filePath string, limit int) ([]string, error)

	// SuggestContextFiles suggests files to include in context based on git history
	SuggestContextFiles(taskDescription string, currentFiles []string, limit int) ([]string, error)

	// GetRelatedByAuthor returns files frequently changed by the same author
	GetRelatedByAuthor(filePath string, limit int) ([]string, error)
}

// RecentChange represents a recently changed file from git history
type RecentChange struct {
	FilePath    string    `json:"filePath"`
	ChangeCount int       `json:"changeCount"`
	LastChanged time.Time `json:"lastChanged"`
	Authors     []string  `json:"authors"`
}

// =============================================================================
// Call Graph Builder Interface
// =============================================================================

// CallGraphBuilder defines interface for building and querying call graphs
type CallGraphBuilder interface {
	// Build builds the call graph for a project
	Build(projectRoot string) (*CallGraph, error)

	// GetCallers returns functions that call the specified function
	GetCallers(functionID string) []CallGraphNode

	// GetCallees returns functions called by the specified function
	GetCallees(functionID string) []CallGraphNode

	// GetImpact returns all functions affected if the specified function changes
	GetImpact(functionID string, maxDepth int) []CallGraphNode

	// GetCallChain finds call chains between two functions
	GetCallChain(startID, endID string, maxDepth int) [][]string
}

// CallGraph represents a call graph for a project
type CallGraph struct {
	Nodes map[string]*CallGraphNode `json:"nodes"`
	Edges []CallGraphEdge           `json:"edges"`
}

// CallGraphNode represents a function in the call graph
type CallGraphNode struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FilePath string `json:"filePath"`
	Line     int    `json:"line"`
	Package  string `json:"package,omitempty"`
}

// CallGraphEdge represents a call relationship
type CallGraphEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Line int    `json:"line,omitempty"`
}

// =============================================================================
// Project Structure Interface
// =============================================================================

// ProjectStructureDetector defines interface for detecting project structure
type ProjectStructureDetector interface {
	// Detect analyzes project structure (basic detection)
	Detect(projectPath string) (*ProjectStructureInfo, error)

	// DetectLanguages detects programming languages used in project
	DetectLanguages(projectPath string) ([]string, error)

	// DetectFrameworks detects frameworks used in project
	DetectFrameworks(projectPath string) ([]FrameworkInfo, error)

	// DetectStructure analyzes project and returns complete structure info
	DetectStructure(projectPath string) (*ProjectStructure, error)

	// DetectArchitecture detects architecture pattern
	DetectArchitecture(projectPath string) (*ArchitectureInfo, error)

	// DetectConventions detects naming and code conventions
	DetectConventions(projectPath string) (*ConventionInfo, error)

	// GetRelatedLayers returns layers related to a file
	GetRelatedLayers(projectPath, filePath string) ([]LayerInfo, error)

	// SuggestRelatedFiles suggests related files based on architecture
	SuggestRelatedFiles(projectPath, filePath string) ([]string, error)
}

// ProjectStructureInfo contains detected project structure information
type ProjectStructureInfo struct {
	Languages   []string          `json:"languages"`
	Frameworks  []string          `json:"frameworks"`
	BuildTools  []string          `json:"buildTools"`
	PackageFile string            `json:"packageFile,omitempty"`
	EntryPoints []string          `json:"entryPoints,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}
