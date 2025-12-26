package guardrails

import (
	"fmt"
	"syntaxia/domain"
	"sync"
	"time"
)

// ServiceImpl реализует GuardrailService
type ServiceImpl struct {
	log              domain.Logger
	policies         []domain.GuardrailPolicy
	budgets          []domain.BudgetPolicy
	mu               sync.RWMutex
	config           domain.GuardrailConfig
	ephemeralMode    bool
	ephemeralEnd     time.Time
	opaService       domain.OPAService
	fileStatProvider domain.FileStatProvider
	taskTypeProvider domain.TaskTypeProvider
}

// NewService создает новый сервис guardrails
func NewService(log domain.Logger, opaService domain.OPAService, fileStatProvider domain.FileStatProvider) domain.GuardrailService {
	service := &ServiceImpl{
		log:      log,
		policies: make([]domain.GuardrailPolicy, 0),
		budgets:  make([]domain.BudgetPolicy, 0),
		config: domain.GuardrailConfig{
			FailClosed:           true,
			EnableEphemeralMode:  true,
			EphemeralTimeout:     5 * time.Minute,
			EnableTaskValidation: true,
			EnableBudgetTracking: true,
			EnablePathValidation: true,
		},
		opaService:       opaService,
		fileStatProvider: fileStatProvider,
	}

	service.initializeDefaultPolicies()

	return service
}

// SetTaskTypeProvider устанавливает провайдер типов задач
func (s *ServiceImpl) SetTaskTypeProvider(taskTypeProvider domain.TaskTypeProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.taskTypeProvider = taskTypeProvider
}

// EnableEphemeralMode включает ephemeral mode для критических путей
func (s *ServiceImpl) EnableEphemeralMode(taskID, taskType string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.taskTypeProvider == nil {
		s.log.Warning("TaskTypeProvider not initialized - ephemeral mode validation skipped")
		if taskType != string(domain.TaskTypeScaffold) && taskType != string(domain.TaskTypeDepsFix) {
			return fmt.Errorf("ephemeral mode only allowed for scaffold/deps_fix tasks")
		}
	}

	if taskType != string(domain.TaskTypeScaffold) && taskType != string(domain.TaskTypeDepsFix) {
		return fmt.Errorf("ephemeral mode only allowed for scaffold/deps_fix tasks")
	}

	s.ephemeralMode = true
	s.ephemeralEnd = time.Now().Add(duration)
	s.log.Info(fmt.Sprintf("Ephemeral mode enabled for task %s until %s", taskID, s.ephemeralEnd.Format(time.RFC3339)))

	return nil
}

// DisableEphemeralMode отключает ephemeral mode
func (s *ServiceImpl) DisableEphemeralMode() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.ephemeralMode = false
	s.ephemeralEnd = time.Time{}
	s.log.Info("Ephemeral mode disabled")
}

// GetPolicies возвращает все политики
func (s *ServiceImpl) GetPolicies() ([]domain.GuardrailPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	policies := make([]domain.GuardrailPolicy, len(s.policies))
	copy(policies, s.policies)
	return policies, nil
}

// GetBudgetPolicies возвращает бюджетные политики
func (s *ServiceImpl) GetBudgetPolicies() ([]domain.BudgetPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	budgets := make([]domain.BudgetPolicy, len(s.budgets))
	copy(budgets, s.budgets)
	return budgets, nil
}

// AddPolicy добавляет новую политику
func (s *ServiceImpl) AddPolicy(policy domain.GuardrailPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.policies {
		if existing.ID == policy.ID {
			return fmt.Errorf("policy with ID %s already exists", policy.ID)
		}
	}

	s.policies = append(s.policies, policy)
	s.log.Info(fmt.Sprintf("Added guardrail policy: %s", policy.Name))
	return nil
}

// RemovePolicy удаляет политику
func (s *ServiceImpl) RemovePolicy(policyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, policy := range s.policies {
		if policy.ID == policyID {
			s.policies = append(s.policies[:i], s.policies[i+1:]...)
			s.log.Info(fmt.Sprintf("Removed guardrail policy: %s", policy.Name))
			return nil
		}
	}

	return fmt.Errorf("policy with ID %s not found", policyID)
}

// UpdatePolicy обновляет политику
func (s *ServiceImpl) UpdatePolicy(policy domain.GuardrailPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.policies {
		if existing.ID == policy.ID {
			s.policies[i] = policy
			s.log.Info(fmt.Sprintf("Updated guardrail policy: %s", policy.Name))
			return nil
		}
	}

	return fmt.Errorf("policy with ID %s not found", policy.ID)
}

// AddBudgetPolicy добавляет бюджетную политику
func (s *ServiceImpl) AddBudgetPolicy(policy domain.BudgetPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.budgets {
		if existing.ID == policy.ID {
			return fmt.Errorf("budget policy with ID %s already exists", policy.ID)
		}
	}

	s.budgets = append(s.budgets, policy)
	s.log.Info(fmt.Sprintf("Added budget policy: %s", policy.Name))
	return nil
}

// RemoveBudgetPolicy удаляет бюджетную политику
func (s *ServiceImpl) RemoveBudgetPolicy(policyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, policy := range s.budgets {
		if policy.ID == policyID {
			s.budgets = append(s.budgets[:i], s.budgets[i+1:]...)
			s.log.Info(fmt.Sprintf("Removed budget policy: %s", policy.Name))
			return nil
		}
	}

	return fmt.Errorf("budget policy with ID %s not found", policyID)
}

// UpdateBudgetPolicy обновляет бюджетную политику
func (s *ServiceImpl) UpdateBudgetPolicy(policy domain.BudgetPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.budgets {
		if existing.ID == policy.ID {
			s.budgets[i] = policy
			s.log.Info(fmt.Sprintf("Updated budget policy: %s", policy.Name))
			return nil
		}
	}

	return fmt.Errorf("budget policy with ID %s not found", policy.ID)
}

// GetConfig возвращает конфигурацию guardrails
func (s *ServiceImpl) GetConfig() domain.GuardrailConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.config
}

// UpdateConfig обновляет конфигурацию guardrails
func (s *ServiceImpl) UpdateConfig(config domain.GuardrailConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.config = config
	s.log.Info("Updated guardrail configuration")
	return nil
}
