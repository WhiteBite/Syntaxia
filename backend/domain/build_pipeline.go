package domain

import "context"

// BuildResult представляет результат сборки
type BuildResult struct {
	Success     bool                   `json:"success"`
	Language    string                 `json:"language"`
	ProjectPath string                 `json:"projectPath"`
	Output      string                 `json:"output"`
	Error       string                 `json:"error,omitempty"`
	Duration    float64                `json:"duration"`
	Artifacts   []string               `json:"artifacts,omitempty"`
	Warnings    []string               `json:"warnings,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TypeCheckResult представляет результат проверки типов
type TypeCheckResult struct {
	Success     bool                   `json:"success"`
	Language    string                 `json:"language"`
	ProjectPath string                 `json:"projectPath"`
	Output      string                 `json:"output"`
	Error       string                 `json:"error,omitempty"`
	Duration    float64                `json:"duration"`
	Issues      []*TypeIssue           `json:"issues,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// TypeIssue представляет проблему с типами
type TypeIssue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "error", "warning", "info"
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
}

// BuildPipeline определяет интерфейс для build/type-check pipeline
type BuildPipeline interface {
	// Build выполняет сборку проекта
	Build(ctx context.Context, projectPath, language string) (*BuildResult, error)

	// TypeCheck выполняет проверку типов
	TypeCheck(ctx context.Context, projectPath, language string) (*TypeCheckResult, error)

	// BuildAndTypeCheck выполняет сборку и проверку типов
	BuildAndTypeCheck(ctx context.Context, projectPath, language string) (*BuildResult, *TypeCheckResult, error)

	// GetSupportedLanguages возвращает поддерживаемые языки
	GetSupportedLanguages() []string
}

// BuildConfig определяет конфигурацию сборки
type BuildConfig struct {
	Language    string            `json:"language"`
	ProjectPath string            `json:"projectPath"`
	BuildArgs   []string          `json:"buildArgs,omitempty"`
	EnvVars     map[string]string `json:"envVars,omitempty"`
	Timeout     int               `json:"timeout"` // в секундах
	Parallel    bool              `json:"parallel"`
	CleanBefore bool              `json:"cleanBefore"`
	OutputDir   string            `json:"outputDir,omitempty"`
}

// TypeCheckConfig определяет конфигурацию проверки типов
type TypeCheckConfig struct {
	Language     string            `json:"language"`
	ProjectPath  string            `json:"projectPath"`
	Strict       bool              `json:"strict"`
	IncludeTests bool              `json:"includeTests"`
	EnvVars      map[string]string `json:"envVars,omitempty"`
	Timeout      int               `json:"timeout"` // в секундах
	ConfigFile   string            `json:"configFile,omitempty"`
}

// ProjectValidationResult представляет результат валидации проекта
type ProjectValidationResult struct {
	Success     bool                                 `json:"success"`
	ProjectPath string                               `json:"projectPath"`
	Languages   []string                             `json:"languages"`
	Results     map[string]*LanguageValidationResult `json:"results"`
}

// VerificationConfig представляет конфигурацию verification pipeline
type VerificationConfig struct {
	ProjectPath string   `json:"projectPath"`
	Languages   []string `json:"languages"`
	Timeout     int      `json:"timeout"` // в секундах
	Verbose     bool     `json:"verbose"`
}

// VerificationResult представляет результат verification pipeline
type VerificationResult struct {
	ProjectPath string              `json:"projectPath"`
	Languages   []string            `json:"languages"`
	Success     bool                `json:"success"`
	StartedAt   string              `json:"startedAt"`
	CompletedAt string              `json:"completedAt"`
	Steps       []*VerificationStep `json:"steps"`
}

// SandboxConfig определяет конфигурацию песочницы
type SandboxConfig struct {
	Engine      string            `json:"engine"`      // "docker" или "podman"
	Image       string            `json:"image"`       // базовый образ
	Timeout     int               `json:"timeout"`     // таймаут в секундах
	MemoryLimit string            `json:"memoryLimit"` // лимит памяти (например, "512m")
	CPULimit    string            `json:"cpuLimit"`    // лимит CPU (например, "1.0")
	NetworkMode string            `json:"networkMode"` // режим сети ("none", "host", "bridge")
	ReadOnly    bool              `json:"readOnly"`    // только для чтения
	Mounts      []SandboxMount    `json:"mounts"`      // точки монтирования
	EnvVars     map[string]string `json:"envVars"`     // переменные окружения
	User        string            `json:"user"`        // пользователь для запуска
	WorkingDir  string            `json:"workingDir"`  // рабочая директория
}

// SandboxMount определяет точку монтирования
type SandboxMount struct {
	Source   string `json:"source"`   // источник на хосте
	Target   string `json:"target"`   // цель в контейнере
	ReadOnly bool   `json:"readOnly"` // только для чтения
	Type     string `json:"type"`     // тип монтирования ("bind", "volume")
}

// SandboxResult представляет результат выполнения в песочнице
type SandboxResult struct {
	Success     bool                   `json:"success"`
	ExitCode    int                    `json:"exitCode"`
	Output      string                 `json:"output"`
	Error       string                 `json:"error,omitempty"`
	Duration    float64                `json:"duration"`
	ContainerID string                 `json:"containerId,omitempty"`
	Logs        string                 `json:"logs,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SandboxRunner определяет интерфейс для запуска команд в песочнице
type SandboxRunner interface {
	// Run выполняет команду в песочнице
	Run(ctx context.Context, config SandboxConfig, command []string) (*SandboxResult, error)

	// Cleanup очищает ресурсы песочницы
	Cleanup(ctx context.Context, containerID string) error

	// IsAvailable проверяет доступность движка песочницы
	IsAvailable(ctx context.Context) bool

	// GetInfo возвращает информацию о движке
	GetInfo(ctx context.Context) (map[string]interface{}, error)
}

// VerificationStep представляет один шаг verification pipeline
type VerificationStep struct {
	Name        string      `json:"name"`
	Success     bool        `json:"success"`
	StartedAt   string      `json:"startedAt"`
	CompletedAt string      `json:"completedAt"`
	Result      interface{} `json:"result,omitempty"`
	Error       error       `json:"error,omitempty"`
}

// FormatResult представляет результат форматирования
type FormatResult struct {
	Language   string `json:"language"`
	Success    bool   `json:"success"`
	FilesCount int    `json:"filesCount"`
	Error      string `json:"error,omitempty"`
}

// LanguageValidationResult представляет результат валидации языка
type LanguageValidationResult struct {
	Success         bool             `json:"success"`
	Language        string           `json:"language"`
	TypeCheckResult *TypeCheckResult `json:"typeCheckResult,omitempty"`
	BuildResult     *BuildResult     `json:"buildResult,omitempty"`
	TypeCheckError  string           `json:"typeCheckError,omitempty"`
	BuildError      string           `json:"buildError,omitempty"`
}
