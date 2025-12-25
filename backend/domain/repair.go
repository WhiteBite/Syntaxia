package domain

import (
	"context"
	"time"
)

// RepairRule определяет правило для автоматического исправления
type RepairRule struct {
	ID          string
	Name        string
	Description string
	Pattern     string // regex pattern для поиска ошибок
	Fix         string // шаблон исправления
	Priority    int    // приоритет правила (выше = важнее)
	Language    string // язык программирования
	Category    string // категория ошибки (syntax, lint, build, etc.)
}

// RepairResult результат выполнения repair операции
type RepairResult struct {
	Success    bool
	RuleID     string
	FixedFiles []string
	Error      string
	Duration   time.Duration
	Attempts   int
}

// RepairRequest запрос на выполнение repair
type RepairRequest struct {
	ProjectPath string
	ErrorOutput string
	Language    string
	MaxAttempts int
	Rules       []RepairRule
}

// RepairService интерфейс для сервиса repair
type RepairService interface {
	// ExecuteRepair выполняет repair цикл
	ExecuteRepair(ctx context.Context, req RepairRequest) (*RepairResult, error)

	// GetAvailableRules возвращает доступные правила для языка
	GetAvailableRules(language string) ([]RepairRule, error)

	// AddRule добавляет новое правило
	AddRule(rule RepairRule) error

	// RemoveRule удаляет правило
	RemoveRule(ruleID string) error

	// ValidateRule проверяет корректность правила
	ValidateRule(rule RepairRule) error
}
