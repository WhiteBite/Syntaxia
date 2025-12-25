package policy

import (
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
)

// OPAService реализует сервис для работы с OPA политиками
type OPAService struct {
	log       domain.Logger
	policyDir string
}

// NewOPAService создает новый OPA сервис
func NewOPAService(log domain.Logger) *OPAService {
	return &OPAService{
		log:       log,
		policyDir: "backend/infrastructure/policy",
	}
}

// ValidatePath проверяет путь через OPA политики
func (s *OPAService) ValidatePath(path string) (*domain.OPAValidationResult, error) {
	// Загружаем политики
	policies, err := s.loadPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	// Создаем входные данные для OPA
	input := map[string]interface{}{
		"path": path,
	}

	// Выполняем валидацию через OPA
	result, err := s.evaluatePolicy(policies, "guardrails.validate_path", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	// Парсим результат
	validationResult := &domain.OPAValidationResult{
		Valid:      true,
		Violations: make([]domain.OPAViolation, 0),
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if valid, exists := resultMap["valid"].(bool); exists {
			validationResult.Valid = valid
		}

		if violations, exists := resultMap["violations"].([]interface{}); exists {
			for _, v := range violations {
				if violationMap, ok := v.(map[string]interface{}); ok {
					violation := domain.OPAViolation{
						Type:    s.getString(violationMap, "type"),
						Message: s.getString(violationMap, "message"),
					}
					validationResult.Violations = append(validationResult.Violations, violation)
				}
			}
		}
	}

	return validationResult, nil
}

// ValidateBudget проверяет бюджет через OPA политики
func (s *OPAService) ValidateBudget(budgetType string, current, limit int64) (*domain.OPAValidationResult, error) {
	// Загружаем политики
	policies, err := s.loadPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	// Создаем входные данные для OPA
	input := map[string]interface{}{
		"budget_type": budgetType,
		"current":     current,
		"limit":       limit,
	}

	// Выполняем валидацию через OPA
	result, err := s.evaluatePolicy(policies, "guardrails.validate_budget", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	// Парсим результат
	validationResult := &domain.OPAValidationResult{
		Valid:      true,
		Violations: make([]domain.OPAViolation, 0),
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if valid, exists := resultMap["valid"].(bool); exists {
			validationResult.Valid = valid
		}

		if violations, exists := resultMap["violations"].([]interface{}); exists {
			for _, v := range violations {
				if violationMap, ok := v.(map[string]interface{}); ok {
					violation := domain.OPAViolation{
						Type:    s.getString(violationMap, "type"),
						Message: s.getString(violationMap, "message"),
					}
					validationResult.Violations = append(validationResult.Violations, violation)
				}
			}
		}
	}

	return validationResult, nil
}

// ValidateTask проверяет задачу через OPA политики
func (s *OPAService) ValidateTask(taskID string, files []string, linesChanged int64, ephemeralMode bool) (*domain.OPAValidationResult, error) {
	// Загружаем политики
	policies, err := s.loadPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	// Создаем входные данные для OPA
	input := map[string]interface{}{
		"task_id":        taskID,
		"files":          files,
		"lines_changed":  linesChanged,
		"ephemeral_mode": ephemeralMode,
	}

	// Выполняем валидацию через OPA
	result, err := s.evaluatePolicy(policies, "guardrails.validate_task", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	// Парсим результат
	validationResult := &domain.OPAValidationResult{
		Valid:      true,
		Violations: make([]domain.OPAViolation, 0),
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if valid, exists := resultMap["valid"].(bool); exists {
			validationResult.Valid = valid
		}

		if violations, exists := resultMap["violations"].([]interface{}); exists {
			for _, v := range violations {
				if violationMap, ok := v.(map[string]interface{}); ok {
					violation := domain.OPAViolation{
						Type:    s.getString(violationMap, "type"),
						Message: s.getString(violationMap, "message"),
					}
					validationResult.Violations = append(validationResult.Violations, violation)
				}
			}
		}

		if ephemeralAllowed, exists := resultMap["ephemeral_allowed"].(bool); exists {
			validationResult.EphemeralAllowed = ephemeralAllowed
		}
	}

	return validationResult, nil
}

// ValidateConfig проверяет конфигурацию через OPA политики
func (s *OPAService) ValidateConfig(config domain.GuardrailConfig) (*domain.OPAValidationResult, error) {
	// Загружаем политики
	policies, err := s.loadPolicies()
	if err != nil {
		return nil, fmt.Errorf("failed to load policies: %w", err)
	}

	// Создаем входные данные для OPA
	input := map[string]interface{}{
		"fail_closed":            config.FailClosed,
		"enable_ephemeral_mode":  config.EnableEphemeralMode,
		"enable_task_validation": config.EnableTaskValidation,
		"enable_budget_tracking": config.EnableBudgetTracking,
		"enable_path_validation": config.EnablePathValidation,
	}

	// Выполняем валидацию через OPA
	result, err := s.evaluatePolicy(policies, "guardrails.validate_config", input)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate policy: %w", err)
	}

	// Парсим результат
	validationResult := &domain.OPAValidationResult{
		Valid:      true,
		Violations: make([]domain.OPAViolation, 0),
	}

	if resultMap, ok := result.(map[string]interface{}); ok {
		if valid, exists := resultMap["valid"].(bool); exists {
			validationResult.Valid = valid
		}

		if warnings, exists := resultMap["warnings"].([]interface{}); exists {
			for _, w := range warnings {
				if warning, ok := w.(string); ok {
					violation := domain.OPAViolation{
						Type:    "config_warning",
						Message: warning,
					}
					validationResult.Violations = append(validationResult.Violations, violation)
				}
			}
		}
	}

	return validationResult, nil
}

// loadPolicies загружает политики из файлов
func (s *OPAService) loadPolicies() (string, error) {
	policyFile := filepath.Join(s.policyDir, "guardrails.rego")

	data, err := os.ReadFile(policyFile)
	if err != nil {
		return "", fmt.Errorf("failed to read policy file: %w", err)
	}

	return string(data), nil
}

// evaluatePolicy выполняет оценку политики
// В реальной реализации здесь будет вызов OPA engine
func (s *OPAService) evaluatePolicy(policies, query string, input map[string]interface{}) (interface{}, error) {
	// Это упрощенная реализация для демонстрации
	// В реальном проекте здесь будет интеграция с OPA engine

	// Логируем вызов для отладки
	s.log.Info(fmt.Sprintf("OPA evaluation: query=%s, input=%v", query, input))

	// Возвращаем заглушку - в реальной реализации здесь будет вызов OPA
	return map[string]interface{}{
		"valid":      true,
		"violations": []interface{}{},
	}, nil
}

// getString извлекает строку из map
func (s *OPAService) getString(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
