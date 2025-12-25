package domain

import "time"

// TaskTypeProvider интерфейс для получения типа задачи
// Используется для разрыва циклической зависимости между GuardrailService и TaskflowService
type TaskTypeProvider interface {
	// GetTaskType возвращает тип задачи по ID
	GetTaskType(taskID string) (string, error)
}

// GuardrailPolicy определяет политику безопасности
type GuardrailPolicy struct {
	ID          string
	Name        string
	Description string
	Type        GuardrailType
	Severity    GuardrailSeverity
	Enabled     bool
	Rules       []GuardrailRule
}

// GuardrailType тип guardrail
type GuardrailType string

const (
	GuardrailTypeForbiddenPath GuardrailType = "forbidden_path"
)

// GuardrailSeverity уровень серьезности
type GuardrailSeverity string

const (
	GuardrailSeverityLow    GuardrailSeverity = "low"
	GuardrailSeverityMedium GuardrailSeverity = "medium"
	GuardrailSeverityBlock  GuardrailSeverity = "block"
)

// GuardrailRule правило guardrail
type GuardrailRule struct {
	ID          string
	Pattern     string // regex pattern или glob pattern
	Description string
	Action      GuardrailAction
	Message     string
}

// GuardrailAction действие при срабатывании правила
type GuardrailAction string

const (
	GuardrailActionBlock GuardrailAction = "block"
)

// BudgetPolicy определяет бюджетные ограничения
type BudgetPolicy struct {
	ID          string
	Name        string
	Description string
	Type        BudgetType
	Limit       int64
	Unit        BudgetUnit
	TimeWindow  time.Duration
	Enabled     bool
}

// BudgetType тип бюджета
type BudgetType string

const (
	BudgetTypeFiles  BudgetType = "files"
	BudgetTypeLines  BudgetType = "lines"
	BudgetTypeTokens BudgetType = "tokens"
)

// BudgetUnit единица измерения бюджета
type BudgetUnit string

const (
	BudgetUnitCount BudgetUnit = "count"
)

// GuardrailViolation нарушение guardrail
type GuardrailViolation struct {
	PolicyID   string
	RuleID     string
	Severity   GuardrailSeverity
	Message    string
	FilePath   string
	LineNumber int
	Timestamp  time.Time
	Context    map[string]interface{}
}

// BudgetViolation нарушение бюджета
type BudgetViolation struct {
	PolicyID  string
	Type      BudgetType
	Current   int64
	Limit     int64
	Unit      BudgetUnit
	Message   string
	Timestamp time.Time
	Context   map[string]interface{}
}

// GuardrailConfig конфигурация guardrails
type GuardrailConfig struct {
	FailClosed           bool
	EnableEphemeralMode  bool
	EphemeralTimeout     time.Duration
	EnableTaskValidation bool
	EnableBudgetTracking bool
	EnablePathValidation bool
}

// TaskValidationResult результат валидации задачи
type TaskValidationResult struct {
	TaskID           string
	Valid            bool
	Error            string
	Violations       []GuardrailViolation
	BudgetViolations []BudgetViolation
	Timestamp        time.Time
}

// OPAViolation нарушение OPA политики
type OPAViolation struct {
	Type    string
	Message string
}

// OPAValidationResult результат валидации OPA
type OPAValidationResult struct {
	Valid            bool
	Violations       []OPAViolation
	EphemeralAllowed bool
}

// GuardrailService интерфейс для сервиса guardrails
type GuardrailService interface {
	// ValidatePath проверяет путь на соответствие политикам
	ValidatePath(path string) ([]GuardrailViolation, error)

	// ValidateBudget проверяет бюджетные ограничения
	ValidateBudget(budgetType BudgetType, current int64) ([]BudgetViolation, error)

	// ValidateTask проверяет задачу на соответствие политикам
	ValidateTask(taskID string, files []string, linesChanged int64) (*TaskValidationResult, error)

	// EnableEphemeralMode включает ephemeral mode для критических путей
	EnableEphemeralMode(taskID string, taskType string, duration time.Duration) error

	// DisableEphemeralMode отключает ephemeral mode
	DisableEphemeralMode()

	// GetPolicies возвращает все политики
	GetPolicies() ([]GuardrailPolicy, error)

	// GetBudgetPolicies возвращает бюджетные политики
	GetBudgetPolicies() ([]BudgetPolicy, error)

	// AddPolicy добавляет новую политику
	AddPolicy(policy GuardrailPolicy) error

	// RemovePolicy удаляет политику
	RemovePolicy(policyID string) error

	// UpdatePolicy обновляет политику
	UpdatePolicy(policy GuardrailPolicy) error

	// AddBudgetPolicy добавляет бюджетную политику
	AddBudgetPolicy(policy BudgetPolicy) error

	// RemoveBudgetPolicy удаляет бюджетную политику
	RemoveBudgetPolicy(policyID string) error

	// UpdateBudgetPolicy обновляет бюджетную политику
	UpdateBudgetPolicy(policy BudgetPolicy) error

	// GetConfig возвращает конфигурацию guardrails
	GetConfig() GuardrailConfig

	// UpdateConfig обновляет конфигурацию guardrails
	UpdateConfig(config GuardrailConfig) error

	// SetTaskTypeProvider устанавливает провайдер типов задач для разрешения циклической зависимости
	SetTaskTypeProvider(taskTypeProvider TaskTypeProvider)
}
