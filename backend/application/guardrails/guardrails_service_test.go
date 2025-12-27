package guardrails

import (
	"syntaxia/domain"
	"testing"
	"time"
)

type mockLogger struct{}

func (m *mockLogger) Debug(msg string)   {}
func (m *mockLogger) Info(msg string)    {}
func (m *mockLogger) Warning(msg string) {}
func (m *mockLogger) Error(msg string)   {}
func (m *mockLogger) Fatal(msg string)   {}

type mockOPAService struct {
	pathResult   *domain.OPAValidationResult
	budgetResult *domain.OPAValidationResult
	configResult *domain.OPAValidationResult
}

func (m *mockOPAService) ValidatePath(path string) (*domain.OPAValidationResult, error) {
	if m.pathResult != nil {
		return m.pathResult, nil
	}
	return &domain.OPAValidationResult{Valid: true}, nil
}

func (m *mockOPAService) ValidateBudget(budgetType string, current, limit int64) (*domain.OPAValidationResult, error) {
	if m.budgetResult != nil {
		return m.budgetResult, nil
	}
	return &domain.OPAValidationResult{Valid: true}, nil
}

func (m *mockOPAService) ValidateTask(taskID string, files []string, linesChanged int64, ephemeralMode bool) (*domain.OPAValidationResult, error) {
	return &domain.OPAValidationResult{Valid: true}, nil
}

func (m *mockOPAService) ValidateConfig(config domain.GuardrailConfig) (*domain.OPAValidationResult, error) {
	if m.configResult != nil {
		return m.configResult, nil
	}
	return &domain.OPAValidationResult{Valid: true}, nil
}

func newTestService() *GuardrailsServiceImpl {
	return NewGuardrailsService(&mockLogger{}, &mockOPAService{}, nil)
}

func TestValidateFileOperation_CriticalFiles(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		operation OperationType
		wantValid bool
	}{
		{"go.mod blocked", "go.mod", OperationModify, false},
		{"go.sum blocked", "go.sum", OperationModify, false},
		{"package.json blocked", "package.json", OperationModify, false},
		{".env blocked", ".env", OperationModify, false},
		{".env.local blocked", ".env.local", OperationModify, false},
		{"regular file allowed", "main.go", OperationModify, true},
		{"nested regular file allowed", "src/utils/helper.go", OperationModify, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			op := FileOperation{Path: tt.path, Operation: tt.operation, LinesChanged: 10, Confirmed: true}
			result, err := svc.ValidateFileOperation(op)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Valid != tt.wantValid {
				t.Errorf("got Valid=%v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

func TestValidateFileOperation_DirectoryBlacklist(t *testing.T) {
	tests := []struct {
		name      string
		path      string
		wantValid bool
	}{
		{"node_modules blocked", "node_modules/package/index.js", false},
		{".git blocked", ".git/config", false},
		{"vendor blocked", "vendor/github.com/pkg/file.go", false},
		{"dist blocked", "dist/bundle.js", false},
		{"tmp blocked", "tmp/cache.txt", false},
		{"src allowed", "src/main.go", true},
		{"backend allowed", "backend/service.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			op := FileOperation{Path: tt.path, Operation: OperationModify, LinesChanged: 10, Confirmed: true}
			result, err := svc.ValidateFileOperation(op)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Valid != tt.wantValid {
				t.Errorf("got Valid=%v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

func TestValidateFileOperation_DeletionConfirmation(t *testing.T) {
	tests := []struct {
		name      string
		operation OperationType
		confirmed bool
		wantValid bool
	}{
		{"delete without confirmation blocked", OperationDelete, false, false},
		{"delete with confirmation allowed", OperationDelete, true, true},
		{"modify without confirmation allowed", OperationModify, false, true},
		{"create without confirmation allowed", OperationCreate, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			op := FileOperation{Path: "src/file.go", Operation: tt.operation, LinesChanged: 10, Confirmed: tt.confirmed}
			result, err := svc.ValidateFileOperation(op)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Valid != tt.wantValid {
				t.Errorf("got Valid=%v, want %v", result.Valid, tt.wantValid)
			}
		})
	}
}

func TestValidateFileOperation_LinesLimit(t *testing.T) {
	tests := []struct {
		name         string
		linesChanged int64
		wantValid    bool
	}{
		{"within limit", 100, true},
		{"at limit", DefaultMaxLinesPerOperation, true},
		{"exceeds limit", DefaultMaxLinesPerOperation + 1, false},
		{"way over limit", 10000, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			op := FileOperation{Path: "src/file.go", Operation: OperationModify, LinesChanged: tt.linesChanged, Confirmed: true}
			result, err := svc.ValidateFileOperation(op)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.Valid != tt.wantValid {
				t.Errorf("got Valid=%v, want %v for lines=%d", result.Valid, tt.wantValid, tt.linesChanged)
			}
		})
	}
}

func TestValidateSessionBudget(t *testing.T) {
	tests := []struct {
		name          string
		existingFiles int
		existingLines int64
		newFiles      int
		newLines      int64
		wantViolation bool
	}{
		{"within limits", 10, 1000, 5, 500, false},
		{"files exceed limit", 45, 1000, 10, 100, true},
		{"lines exceed limit", 10, 4500, 5, 1000, true},
		{"both exceed", 45, 4500, 10, 1000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			sessionID := "test-session"
			svc.UpdateSessionBudget(sessionID, tt.existingFiles, tt.existingLines)
			violations, err := svc.ValidateSessionBudget(sessionID, tt.newFiles, tt.newLines)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			hasViolation := len(violations) > 0
			if hasViolation != tt.wantViolation {
				t.Errorf("got violation=%v, want %v", hasViolation, tt.wantViolation)
			}
		})
	}
}

func TestValidateTokenBudget(t *testing.T) {
	tests := []struct {
		name          string
		tokens        int64
		wantViolation bool
	}{
		{"within limit", 50000, false},
		{"at limit", DefaultMaxTokensPerRequest, false},
		{"exceeds limit", DefaultMaxTokensPerRequest + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			violations, err := svc.ValidateTokenBudget(tt.tokens)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			hasViolation := len(violations) > 0
			if hasViolation != tt.wantViolation {
				t.Errorf("got violation=%v, want %v", hasViolation, tt.wantViolation)
			}
		})
	}
}

func TestDirectoryWhitelist(t *testing.T) {
	tests := []struct {
		name      string
		whitelist []string
		path      string
		wantValid bool
	}{
		{"path in whitelist", []string{"src", "backend"}, "src/main.go", true},
		{"path not in whitelist", []string{"src", "backend"}, "frontend/app.vue", false},
		{"empty whitelist allows all", []string{}, "any/path/file.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			svc.SetDirectoryWhitelist(tt.whitelist)
			svc.SetDirectoryBlacklist([]string{}) // Clear blacklist for this test
			allowed := svc.IsPathAllowed(tt.path)
			if allowed != tt.wantValid {
				t.Errorf("got allowed=%v, want %v", allowed, tt.wantValid)
			}
		})
	}
}

func TestEphemeralMode(t *testing.T) {
	tests := []struct {
		name        string
		taskType    string
		wantEnabled bool
		wantError   bool
	}{
		{"scaffold allowed", string(domain.TaskTypeScaffold), true, false},
		{"deps_fix allowed", string(domain.TaskTypeDepsFix), true, false},
		{"feature not allowed", string(domain.TaskTypeFeature), false, true},
		{"bugfix not allowed", string(domain.TaskTypeBugFix), false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			err := svc.EnableEphemeralMode("task-1", tt.taskType, 5*time.Minute)
			if tt.wantError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantError && !svc.IsEphemeralModeActive() {
				t.Error("ephemeral mode should be active")
			}
		})
	}
}

func TestEphemeralModeAllowsCriticalFiles(t *testing.T) {
	svc := newTestService()
	
	// Without ephemeral mode, go.mod should be blocked
	op := FileOperation{Path: "go.mod", Operation: OperationModify, LinesChanged: 10, Confirmed: true}
	result, _ := svc.ValidateFileOperation(op)
	if result.Valid {
		t.Error("go.mod should be blocked without ephemeral mode")
	}

	// Enable ephemeral mode
	err := svc.EnableEphemeralMode("task-1", string(domain.TaskTypeScaffold), 5*time.Minute)
	if err != nil {
		t.Fatalf("failed to enable ephemeral mode: %v", err)
	}

	// With ephemeral mode, go.mod should be allowed
	result, _ = svc.ValidateFileOperation(op)
	if !result.Valid {
		t.Error("go.mod should be allowed with ephemeral mode")
	}
}

func TestCriticalFileManagement(t *testing.T) {
	svc := newTestService()

	// Add custom critical file
	svc.AddCriticalFile("custom.config")
	files := svc.GetCriticalFiles()
	found := false
	for _, f := range files {
		if f == "custom.config" {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom.config should be in critical files")
	}

	// Remove critical file
	svc.RemoveCriticalFile("custom.config")
	files = svc.GetCriticalFiles()
	for _, f := range files {
		if f == "custom.config" {
			t.Error("custom.config should be removed from critical files")
		}
	}
}

func TestBudgetLimitManagement(t *testing.T) {
	svc := newTestService()

	// Update budget limit
	err := svc.SetBudgetLimit("max-tokens-request", 200000)
	if err != nil {
		t.Fatalf("failed to set budget limit: %v", err)
	}

	limits := svc.GetBudgetLimits()
	if limits["max-tokens-request"] != 200000 {
		t.Errorf("got limit=%d, want 200000", limits["max-tokens-request"])
	}

	// Try to update non-existent budget
	err = svc.SetBudgetLimit("non-existent", 100)
	if err == nil {
		t.Error("expected error for non-existent budget")
	}
}

func TestSessionManagement(t *testing.T) {
	svc := newTestService()
	sessionID := "test-session"

	// Update session
	svc.UpdateSessionBudget(sessionID, 5, 100)
	session := svc.GetSessionBudget(sessionID)
	if session == nil {
		t.Fatal("session should exist")
	}
	if session.FilesChanged != 5 {
		t.Errorf("got files=%d, want 5", session.FilesChanged)
	}
	if session.LinesChanged != 100 {
		t.Errorf("got lines=%d, want 100", session.LinesChanged)
	}

	// Reset session
	svc.ResetSession(sessionID)
	session = svc.GetSessionBudget(sessionID)
	if session != nil {
		t.Error("session should be nil after reset")
	}
}

func TestIsCriticalFile(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantCritical bool
	}{
		{"go.mod is critical", "go.mod", true},
		{"nested go.mod is critical", "subdir/go.mod", true},
		{"package.json is critical", "package.json", true},
		{".env is critical", ".env", true},
		{"regular file not critical", "main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestService()
			if svc.IsCriticalFile(tt.path) != tt.wantCritical {
				t.Errorf("got %v, want %v", svc.IsCriticalFile(tt.path), tt.wantCritical)
			}
		})
	}
}
