package guardrails

import (
	"fmt"
	"syntaxia/domain"
	"time"
)

// ValidatePath проверяет путь на соответствие политикам
func (s *ServiceImpl) ValidatePath(path string) ([]domain.GuardrailViolation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.checkEphemeralExpiry()

	if s.ephemeralMode && s.isCriticalPath(path) {
		s.log.Info(fmt.Sprintf("Critical path %s allowed in ephemeral mode", path))
		return nil, nil
	}

	var violations []domain.GuardrailViolation
	violations = append(violations, s.validateWithOPA(path)...)
	policyViolations, err := s.validateWithPolicies(path)
	if err != nil {
		return append(violations, policyViolations...), err
	}
	return append(violations, policyViolations...), nil
}

// checkEphemeralExpiry checks and disables expired ephemeral mode
func (s *ServiceImpl) checkEphemeralExpiry() {
	if s.ephemeralMode && s.isEphemeralExpired() {
		s.disableEphemeralModeInternal()
	}
}

// validateWithOPA validates path using OPA service
func (s *ServiceImpl) validateWithOPA(path string) []domain.GuardrailViolation {
	if s.opaService == nil {
		return nil
	}
	opaResult, err := s.opaService.ValidatePath(path)
	if err != nil {
		s.log.Warning(fmt.Sprintf("OPA validation failed: %v", err))
		return nil
	}
	if opaResult.Valid {
		return nil
	}

	violations := make([]domain.GuardrailViolation, 0, len(opaResult.Violations))
	for _, v := range opaResult.Violations {
		violations = append(violations, domain.GuardrailViolation{
			PolicyID: "opa-policy", RuleID: v.Type, Severity: domain.GuardrailSeverityBlock,
			Message: v.Message, FilePath: path, Timestamp: time.Now(),
			Context: map[string]any{"opa_violation": true, "violation_type": v.Type},
		})
	}
	return violations
}

// validateWithPolicies validates path against configured policies
func (s *ServiceImpl) validateWithPolicies(path string) ([]domain.GuardrailViolation, error) {
	violations := make([]domain.GuardrailViolation, 0, len(s.policies))
	for _, policy := range s.policies {
		if !policy.Enabled || policy.Type != domain.GuardrailTypeForbiddenPath {
			continue
		}
		v, err := s.checkPolicyRules(path, policy)
		violations = append(violations, v...)
		if err != nil {
			return violations, err
		}
	}
	return violations, nil
}

// checkPolicyRules checks path against policy rules
func (s *ServiceImpl) checkPolicyRules(path string, policy domain.GuardrailPolicy) ([]domain.GuardrailViolation, error) {
	violations := make([]domain.GuardrailViolation, 0, len(policy.Rules))
	for _, rule := range policy.Rules {
		if !s.matchesPath(path, rule.Pattern) {
			continue
		}
		violation := domain.GuardrailViolation{
			PolicyID: policy.ID, RuleID: rule.ID, Severity: policy.Severity,
			Message: rule.Message, FilePath: path, Timestamp: time.Now(),
			Context: map[string]any{"pattern": rule.Pattern, "ephemeralMode": s.ephemeralMode, "failClosed": s.config.FailClosed},
		}
		violations = append(violations, violation)

		if s.config.FailClosed && policy.Severity == domain.GuardrailSeverityBlock {
			s.log.Error(fmt.Sprintf("Guardrail violation blocked: %s - %s", path, rule.Message))
			return violations, fmt.Errorf("guardrail violation: %s", rule.Message)
		}
	}
	return violations, nil
}

// ValidateBudget проверяет бюджетные ограничения
func (s *ServiceImpl) ValidateBudget(budgetType domain.BudgetType, current int64) ([]domain.BudgetViolation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var violations []domain.BudgetViolation

	for _, budget := range s.budgets {
		if !budget.Enabled || budget.Type != budgetType {
			continue
		}

		if current > budget.Limit {
			violation := domain.BudgetViolation{
				PolicyID:  budget.ID,
				Type:      budget.Type,
				Current:   current,
				Limit:     budget.Limit,
				Unit:      budget.Unit,
				Message:   fmt.Sprintf("Budget exceeded: %d %s (limit: %d)", current, budget.Unit, budget.Limit),
				Timestamp: time.Now(),
				Context: map[string]any{
					"timeWindow":   budget.TimeWindow,
					"failClosed":   s.config.FailClosed,
					"budgetType":   budgetType,
					"currentValue": current,
				},
			}
			violations = append(violations, violation)

			if s.config.FailClosed {
				s.log.Error(fmt.Sprintf("Budget violation blocked: %s", violation.Message))
				return violations, fmt.Errorf("budget violation: %s", violation.Message)
			}
		}
	}

	return violations, nil
}

// ValidateTask проверяет задачу на соответствие политикам
func (s *ServiceImpl) ValidateTask(taskID string, files []string, linesChanged int64) (*domain.TaskValidationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := &domain.TaskValidationResult{
		TaskID:           taskID,
		Valid:            true,
		Violations:       make([]domain.GuardrailViolation, 0),
		BudgetViolations: make([]domain.BudgetViolation, 0),
		Timestamp:        time.Now(),
	}

	for _, file := range files {
		pathViolations, err := s.ValidatePath(file)
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
			return result, err
		}
		result.Violations = append(result.Violations, pathViolations...)
	}

	if s.config.EnableBudgetTracking {
		fileBudgetViolations, err := s.ValidateBudget(domain.BudgetTypeFiles, int64(len(files)))
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
			return result, err
		}
		result.BudgetViolations = append(result.BudgetViolations, fileBudgetViolations...)
	}

	if s.config.EnableBudgetTracking {
		lineBudgetViolations, err := s.ValidateBudget(domain.BudgetTypeLines, linesChanged)
		if err != nil {
			result.Valid = false
			result.Error = err.Error()
			return result, err
		}
		result.BudgetViolations = append(result.BudgetViolations, lineBudgetViolations...)
	}

	if len(result.Violations) > 0 || len(result.BudgetViolations) > 0 {
		result.Valid = false
		if s.config.FailClosed {
			result.Error = "Task validation failed due to guardrail violations"
		}
	}

	return result, nil
}
