package guardrails

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"syntaxia/domain"
	"time"
)

// OperationType defines the type of file operation
type OperationType string

const (
	OperationCreate OperationType = "create"
	OperationModify OperationType = "modify"
	OperationDelete OperationType = "delete"
)

// FileOperation represents a file operation request for validation
type FileOperation struct {
	Path         string        `json:"path"`
	Operation    OperationType `json:"operation"`
	LinesChanged int64         `json:"linesChanged"`
	Confirmed    bool          `json:"confirmed"`
}

// DirectoryPolicy defines whitelist/blacklist for directories
type DirectoryPolicy struct {
	Whitelist []string `json:"whitelist"`
	Blacklist []string `json:"blacklist"`
}

// SessionBudget tracks budget usage within a session
type SessionBudget struct {
	SessionID    string    `json:"sessionId"`
	TokensUsed   int64     `json:"tokensUsed"`
	FilesChanged int       `json:"filesChanged"`
	LinesChanged int64     `json:"linesChanged"`
	StartTime    time.Time `json:"startTime"`
	LastUpdated  time.Time `json:"lastUpdated"`
}

// Budget limit constants
const (
	DefaultMaxTokensPerRequest  = 100000
	DefaultMaxFilesPerSession   = 50
	DefaultMaxLinesPerOperation = 500
	DefaultMaxLinesPerSession   = 5000
)

// GuardrailsServiceImpl implements enhanced guardrails with OPA integration
type GuardrailsServiceImpl struct {
	mu               sync.RWMutex
	log              domain.Logger
	opaService       domain.OPAService
	fileStatProvider domain.FileStatProvider
	config           domain.GuardrailConfig
	directoryPolicy  DirectoryPolicy
	criticalFiles    map[string]bool
	policies         []domain.GuardrailPolicy
	budgets          []domain.BudgetPolicy
	sessions         map[string]*SessionBudget
	ephemeralMode    bool
	ephemeralEnd     time.Time
	taskTypeProvider domain.TaskTypeProvider
}

// NewGuardrailsService creates a new guardrails service with OPA integration
func NewGuardrailsService(
	log domain.Logger,
	opaService domain.OPAService,
	fileStatProvider domain.FileStatProvider,
) *GuardrailsServiceImpl {
	svc := &GuardrailsServiceImpl{
		log:              log,
		opaService:       opaService,
		fileStatProvider: fileStatProvider,
		config: domain.GuardrailConfig{
			FailClosed:           true,
			EnableEphemeralMode:  true,
			EphemeralTimeout:     5 * time.Minute,
			EnableTaskValidation: true,
			EnableBudgetTracking: true,
			EnablePathValidation: true,
		},
		directoryPolicy: DirectoryPolicy{
			Whitelist: []string{},
			Blacklist: []string{"node_modules", ".git", "vendor", "dist", "build", "tmp", "temp", "cache"},
		},
		criticalFiles: map[string]bool{
			"go.mod": true, "go.sum": true, "package.json": true, "package-lock.json": true,
			"yarn.lock": true, "pnpm-lock.yaml": true, ".env": true, ".env.local": true, ".env.production": true,
		},
		policies: make([]domain.GuardrailPolicy, 0),
		budgets:  make([]domain.BudgetPolicy, 0),
		sessions: make(map[string]*SessionBudget),
	}
	svc.initDefaultPolicies()
	svc.initDefaultBudgets()
	return svc
}

func (s *GuardrailsServiceImpl) initDefaultPolicies() {
	criticalFilesPolicy := domain.GuardrailPolicy{
		ID: "critical-files", Name: "Critical Files Protection",
		Description: "Prohibits changes to critical project files",
		Type: domain.GuardrailTypeForbiddenPath, Severity: domain.GuardrailSeverityBlock, Enabled: true,
		Rules: []domain.GuardrailRule{
			{ID: "go-mod", Pattern: "go.mod", Description: "go.mod is critical", Action: domain.GuardrailActionBlock, Message: "go.mod changes require confirmation"},
			{ID: "go-sum", Pattern: "go.sum", Description: "go.sum is critical", Action: domain.GuardrailActionBlock, Message: "go.sum changes require confirmation"},
			{ID: "package-json", Pattern: "package.json", Description: "package.json is critical", Action: domain.GuardrailActionBlock, Message: "package.json changes require confirmation"},
			{ID: "env-files", Pattern: ".env", Description: ".env files are critical", Action: domain.GuardrailActionBlock, Message: ".env file changes require confirmation"},
		},
	}
	secretsPolicy := domain.GuardrailPolicy{
		ID: "secrets-protection", Name: "Secrets Protection",
		Description: "Prohibits changes to files containing secrets",
		Type: domain.GuardrailTypeForbiddenPath, Severity: domain.GuardrailSeverityBlock, Enabled: true,
		Rules: []domain.GuardrailRule{
			{ID: "key-files", Pattern: ".key|.pem|.p12|.pfx", Description: "Key files", Action: domain.GuardrailActionBlock, Message: "Key file changes are prohibited"},
			{ID: "secret-files", Pattern: ".secret", Description: "Secret files", Action: domain.GuardrailActionBlock, Message: "Secret file changes are prohibited"},
		},
	}
	s.policies = append(s.policies, criticalFilesPolicy, secretsPolicy)
}

func (s *GuardrailsServiceImpl) initDefaultBudgets() {
	s.budgets = []domain.BudgetPolicy{
		{ID: "max-tokens-request", Name: "Max Tokens Per Request", Description: "Maximum tokens allowed per AI request", Type: domain.BudgetTypeTokens, Limit: DefaultMaxTokensPerRequest, Unit: domain.BudgetUnitCount, Enabled: true},
		{ID: "max-files-session", Name: "Max Files Per Session", Description: "Maximum files that can be changed per session", Type: domain.BudgetTypeFiles, Limit: DefaultMaxFilesPerSession, Unit: domain.BudgetUnitCount, TimeWindow: time.Hour, Enabled: true},
		{ID: "max-lines-operation", Name: "Max Lines Per Operation", Description: "Maximum lines changed per single operation", Type: domain.BudgetTypeLines, Limit: DefaultMaxLinesPerOperation, Unit: domain.BudgetUnitCount, Enabled: true},
		{ID: "max-lines-session", Name: "Max Lines Per Session", Description: "Maximum lines changed per session", Type: domain.BudgetTypeLines, Limit: DefaultMaxLinesPerSession, Unit: domain.BudgetUnitCount, TimeWindow: time.Hour, Enabled: true},
	}
}

// ValidateFileOperation validates a file operation against all guardrails
func (s *GuardrailsServiceImpl) ValidateFileOperation(op FileOperation) (*domain.TaskValidationResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := &domain.TaskValidationResult{
		Valid: true, Violations: make([]domain.GuardrailViolation, 0),
		BudgetViolations: make([]domain.BudgetViolation, 0), Timestamp: time.Now(),
	}

	if violations := s.checkCriticalFile(op); len(violations) > 0 {
		result.Violations = append(result.Violations, violations...)
		result.Valid = false
	}
	if violations := s.checkDirectoryPolicy(op.Path); len(violations) > 0 {
		result.Violations = append(result.Violations, violations...)
		result.Valid = false
	}
	if violations := s.checkDeletionConfirmation(op); len(violations) > 0 {
		result.Violations = append(result.Violations, violations...)
		result.Valid = false
	}
	if violations := s.checkLinesLimit(op); len(violations) > 0 {
		result.BudgetViolations = append(result.BudgetViolations, violations...)
		result.Valid = false
	}
	if s.opaService != nil {
		if opaViolations := s.validateWithOPA(op.Path); len(opaViolations) > 0 {
			result.Violations = append(result.Violations, opaViolations...)
			result.Valid = false
		}
	}
	if !result.Valid && s.config.FailClosed {
		result.Error = "File operation blocked by guardrails"
	}
	return result, nil
}

func (s *GuardrailsServiceImpl) checkCriticalFile(op FileOperation) []domain.GuardrailViolation {
	var violations []domain.GuardrailViolation
	baseName := filepath.Base(op.Path)
	if s.criticalFiles[baseName] && !s.ephemeralMode {
		violations = append(violations, domain.GuardrailViolation{
			PolicyID: "critical-files", RuleID: "critical-file-protection", Severity: domain.GuardrailSeverityBlock,
			Message: fmt.Sprintf("Critical file '%s' cannot be modified without ephemeral mode", baseName),
			FilePath: op.Path, Timestamp: time.Now(),
			Context: map[string]interface{}{"operation": op.Operation, "ephemeralMode": s.ephemeralMode},
		})
	}
	return violations
}

func (s *GuardrailsServiceImpl) checkDirectoryPolicy(path string) []domain.GuardrailViolation {
	var violations []domain.GuardrailViolation
	normalizedPath := filepath.ToSlash(path)
	parts := strings.Split(normalizedPath, "/")

	for _, dir := range s.directoryPolicy.Blacklist {
		for _, part := range parts {
			if part == dir {
				violations = append(violations, domain.GuardrailViolation{
					PolicyID: "directory-policy", RuleID: "blacklist", Severity: domain.GuardrailSeverityBlock,
					Message: fmt.Sprintf("Directory '%s' is blacklisted", dir), FilePath: path, Timestamp: time.Now(),
					Context: map[string]interface{}{"blacklistedDir": dir},
				})
				return violations
			}
		}
	}

	if len(s.directoryPolicy.Whitelist) > 0 {
		allowed := false
		for _, dir := range s.directoryPolicy.Whitelist {
			if strings.HasPrefix(normalizedPath, dir+"/") || normalizedPath == dir {
				allowed = true
				break
			}
		}
		if !allowed {
			violations = append(violations, domain.GuardrailViolation{
				PolicyID: "directory-policy", RuleID: "whitelist", Severity: domain.GuardrailSeverityBlock,
				Message: fmt.Sprintf("Path '%s' is not in whitelist", path), FilePath: path, Timestamp: time.Now(),
				Context: map[string]interface{}{"whitelist": s.directoryPolicy.Whitelist},
			})
		}
	}
	return violations
}

func (s *GuardrailsServiceImpl) checkDeletionConfirmation(op FileOperation) []domain.GuardrailViolation {
	var violations []domain.GuardrailViolation
	if op.Operation == OperationDelete && !op.Confirmed {
		violations = append(violations, domain.GuardrailViolation{
			PolicyID: "deletion-policy", RuleID: "require-confirmation", Severity: domain.GuardrailSeverityBlock,
			Message: fmt.Sprintf("File deletion requires confirmation: %s", op.Path), FilePath: op.Path, Timestamp: time.Now(),
			Context: map[string]interface{}{"requiresConfirmation": true},
		})
	}
	return violations
}

func (s *GuardrailsServiceImpl) checkLinesLimit(op FileOperation) []domain.BudgetViolation {
	var violations []domain.BudgetViolation
	for _, budget := range s.budgets {
		if !budget.Enabled || budget.Type != domain.BudgetTypeLines || budget.ID != "max-lines-operation" {
			continue
		}
		if op.LinesChanged > budget.Limit {
			violations = append(violations, domain.BudgetViolation{
				PolicyID: budget.ID, Type: budget.Type, Current: op.LinesChanged, Limit: budget.Limit, Unit: budget.Unit,
				Message: fmt.Sprintf("Lines changed (%d) exceeds limit (%d)", op.LinesChanged, budget.Limit), Timestamp: time.Now(),
				Context: map[string]interface{}{"operation": op.Operation, "path": op.Path},
			})
		}
	}
	return violations
}

func (s *GuardrailsServiceImpl) validateWithOPA(path string) []domain.GuardrailViolation {
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
			Context: map[string]interface{}{"opaViolation": true, "violationType": v.Type},
		})
	}
	return violations
}

// ValidateSessionBudget validates budget limits for a session
func (s *GuardrailsServiceImpl) ValidateSessionBudget(sessionID string, filesChanged int, linesChanged int64) ([]domain.BudgetViolation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := s.getOrCreateSession(sessionID)
	var violations []domain.BudgetViolation
	totalFiles := session.FilesChanged + filesChanged

	for _, budget := range s.budgets {
		if !budget.Enabled {
			continue
		}
		switch budget.ID {
		case "max-files-session":
			if int64(totalFiles) > budget.Limit {
				violations = append(violations, domain.BudgetViolation{
					PolicyID: budget.ID, Type: budget.Type, Current: int64(totalFiles), Limit: budget.Limit, Unit: budget.Unit,
					Message: fmt.Sprintf("Session files (%d) would exceed limit (%d)", totalFiles, budget.Limit), Timestamp: time.Now(),
					Context: map[string]interface{}{"sessionId": sessionID},
				})
			}
		case "max-lines-session":
			totalLines := session.LinesChanged + linesChanged
			if totalLines > budget.Limit {
				violations = append(violations, domain.BudgetViolation{
					PolicyID: budget.ID, Type: budget.Type, Current: totalLines, Limit: budget.Limit, Unit: budget.Unit,
					Message: fmt.Sprintf("Session lines (%d) would exceed limit (%d)", totalLines, budget.Limit), Timestamp: time.Now(),
					Context: map[string]interface{}{"sessionId": sessionID},
				})
			}
		}
	}
	return violations, nil
}

// ValidateTokenBudget validates token budget for a request
func (s *GuardrailsServiceImpl) ValidateTokenBudget(tokens int64) ([]domain.BudgetViolation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var violations []domain.BudgetViolation
	for _, budget := range s.budgets {
		if !budget.Enabled || budget.ID != "max-tokens-request" {
			continue
		}
		if tokens > budget.Limit {
			violations = append(violations, domain.BudgetViolation{
				PolicyID: budget.ID, Type: budget.Type, Current: tokens, Limit: budget.Limit, Unit: budget.Unit,
				Message: fmt.Sprintf("Request tokens (%d) exceeds limit (%d)", tokens, budget.Limit), Timestamp: time.Now(),
			})
		}
	}
	return violations, nil
}

// UpdateSessionBudget updates session budget after successful operation
func (s *GuardrailsServiceImpl) UpdateSessionBudget(sessionID string, filesChanged int, linesChanged int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session := s.getOrCreateSession(sessionID)
	session.FilesChanged += filesChanged
	session.LinesChanged += linesChanged
	session.LastUpdated = time.Now()
}

func (s *GuardrailsServiceImpl) getOrCreateSession(sessionID string) *SessionBudget {
	session, exists := s.sessions[sessionID]
	if !exists {
		session = &SessionBudget{SessionID: sessionID, StartTime: time.Now(), LastUpdated: time.Now()}
		s.sessions[sessionID] = session
	}
	return session
}

// ResetSession resets session budget
func (s *GuardrailsServiceImpl) ResetSession(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, sessionID)
}

// GetSessionBudget returns current session budget usage
func (s *GuardrailsServiceImpl) GetSessionBudget(sessionID string) *SessionBudget {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if session, exists := s.sessions[sessionID]; exists {
		return &SessionBudget{
			SessionID: session.SessionID, TokensUsed: session.TokensUsed, FilesChanged: session.FilesChanged,
			LinesChanged: session.LinesChanged, StartTime: session.StartTime, LastUpdated: session.LastUpdated,
		}
	}
	return nil
}

// SetDirectoryWhitelist sets the directory whitelist
func (s *GuardrailsServiceImpl) SetDirectoryWhitelist(dirs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.directoryPolicy.Whitelist = dirs
}

// SetDirectoryBlacklist sets the directory blacklist
func (s *GuardrailsServiceImpl) SetDirectoryBlacklist(dirs []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.directoryPolicy.Blacklist = dirs
}

// AddCriticalFile adds a file to the critical files list
func (s *GuardrailsServiceImpl) AddCriticalFile(filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.criticalFiles[filename] = true
}

// RemoveCriticalFile removes a file from the critical files list
func (s *GuardrailsServiceImpl) RemoveCriticalFile(filename string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.criticalFiles, filename)
}

// GetCriticalFiles returns the list of critical files
func (s *GuardrailsServiceImpl) GetCriticalFiles() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	files := make([]string, 0, len(s.criticalFiles))
	for f := range s.criticalFiles {
		files = append(files, f)
	}
	return files
}

// GetDirectoryPolicy returns the current directory policy
func (s *GuardrailsServiceImpl) GetDirectoryPolicy() DirectoryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return DirectoryPolicy{
		Whitelist: append([]string{}, s.directoryPolicy.Whitelist...),
		Blacklist: append([]string{}, s.directoryPolicy.Blacklist...),
	}
}

// SetBudgetLimit updates a budget limit
func (s *GuardrailsServiceImpl) SetBudgetLimit(budgetID string, limit int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, budget := range s.budgets {
		if budget.ID == budgetID {
			s.budgets[i].Limit = limit
			s.log.Info(fmt.Sprintf("Updated budget %s limit to %d", budgetID, limit))
			return nil
		}
	}
	return fmt.Errorf("budget policy %s not found", budgetID)
}

// GetBudgetLimits returns all budget limits
func (s *GuardrailsServiceImpl) GetBudgetLimits() map[string]int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	limits := make(map[string]int64)
	for _, budget := range s.budgets {
		limits[budget.ID] = budget.Limit
	}
	return limits
}

// EnableEphemeralMode enables ephemeral mode for critical operations
func (s *GuardrailsServiceImpl) EnableEphemeralMode(taskID, taskType string, duration time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.config.EnableEphemeralMode {
		return fmt.Errorf("ephemeral mode is disabled in configuration")
	}
	allowedTypes := map[string]bool{string(domain.TaskTypeScaffold): true, string(domain.TaskTypeDepsFix): true}
	if !allowedTypes[taskType] {
		return fmt.Errorf("ephemeral mode only allowed for scaffold/deps_fix tasks, got: %s", taskType)
	}
	s.ephemeralMode = true
	s.ephemeralEnd = time.Now().Add(duration)
	s.log.Info(fmt.Sprintf("Ephemeral mode enabled for task %s (type: %s) until %s", taskID, taskType, s.ephemeralEnd.Format(time.RFC3339)))
	return nil
}

// DisableEphemeralMode disables ephemeral mode
func (s *GuardrailsServiceImpl) DisableEphemeralMode() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ephemeralMode = false
	s.ephemeralEnd = time.Time{}
	s.log.Info("Ephemeral mode disabled")
}

// IsEphemeralModeActive returns whether ephemeral mode is active
func (s *GuardrailsServiceImpl) IsEphemeralModeActive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.ephemeralMode {
		return false
	}
	return !time.Now().After(s.ephemeralEnd)
}

// SetTaskTypeProvider sets the task type provider for ephemeral mode validation
func (s *GuardrailsServiceImpl) SetTaskTypeProvider(provider domain.TaskTypeProvider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.taskTypeProvider = provider
}

// GetConfig returns the current guardrail configuration
func (s *GuardrailsServiceImpl) GetConfig() domain.GuardrailConfig {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// UpdateConfig updates the guardrail configuration
func (s *GuardrailsServiceImpl) UpdateConfig(config domain.GuardrailConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.opaService != nil {
		result, err := s.opaService.ValidateConfig(config)
		if err != nil {
			s.log.Warning(fmt.Sprintf("OPA config validation failed: %v", err))
		} else if !result.Valid {
			for _, v := range result.Violations {
				s.log.Warning(fmt.Sprintf("Config warning: %s", v.Message))
			}
		}
	}
	s.config = config
	s.log.Info("Guardrail configuration updated")
	return nil
}

// IsCriticalFile checks if a file is in the critical files list
func (s *GuardrailsServiceImpl) IsCriticalFile(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	baseName := filepath.Base(path)
	return s.criticalFiles[baseName]
}

// IsPathAllowed checks if a path is allowed by directory policies
func (s *GuardrailsServiceImpl) IsPathAllowed(path string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	violations := s.checkDirectoryPolicy(path)
	return len(violations) == 0
}
