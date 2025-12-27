package repair

import (
	"context"
	"errors"
	"syntaxia/domain"
	"testing"
	"time"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(message string)   {}
func (m *mockLogger) Info(message string)    {}
func (m *mockLogger) Warning(message string) {}
func (m *mockLogger) Error(message string)   {}
func (m *mockLogger) Fatal(message string)   {}

// mockErrorAnalyzer implements domain.ErrorAnalyzer for testing
type mockErrorAnalyzer struct {
	analyzeResult     *domain.ErrorDetails
	analyzeErr        error
	suggestResult     []*domain.CorrectionStep
	suggestErr        error
	classifyResult    domain.ErrorType
}

func (m *mockErrorAnalyzer) AnalyzeError(errorOutput string, stage domain.ProtocolStage) (*domain.ErrorDetails, error) {
	if m.analyzeErr != nil {
		return nil, m.analyzeErr
	}
	return m.analyzeResult, nil
}

func (m *mockErrorAnalyzer) SuggestCorrections(err *domain.ErrorDetails) ([]*domain.CorrectionStep, error) {
	if m.suggestErr != nil {
		return nil, m.suggestErr
	}
	return m.suggestResult, nil
}

func (m *mockErrorAnalyzer) ClassifyErrorType(errorOutput string) domain.ErrorType {
	return m.classifyResult
}

// mockCorrectionEngine implements domain.CorrectionEngine for testing
type mockCorrectionEngine struct {
	applyResult    *domain.CorrectionResult
	applyErr       error
	canHandleValue bool
	applyCalls     int
}

func (m *mockCorrectionEngine) ApplyCorrection(ctx context.Context, step *domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	m.applyCalls++
	if m.applyErr != nil {
		return nil, m.applyErr
	}
	return m.applyResult, nil
}

func (m *mockCorrectionEngine) ApplyCorrections(ctx context.Context, steps []*domain.CorrectionStep, projectPath string) (*domain.CorrectionResult, error) {
	return m.applyResult, m.applyErr
}

func (m *mockCorrectionEngine) CanHandle(err *domain.ErrorDetails) bool {
	return m.canHandleValue
}

// mockFileSystem implements domain.FileSystemProvider for testing
type mockFileSystem struct {
	readResult []byte
	readErr    error
	writeErr   error
}

func (m *mockFileSystem) ReadFile(filename string) ([]byte, error) {
	if m.readErr != nil {
		return nil, m.readErr
	}
	return m.readResult, nil
}

func (m *mockFileSystem) WriteFile(filename string, data []byte, perm int) error {
	return m.writeErr
}

func (m *mockFileSystem) MkdirAll(path string, perm int) error {
	return nil
}

func TestDefaultSelfCorrectionConfig(t *testing.T) {
	config := DefaultSelfCorrectionConfig()

	tests := []struct {
		name     string
		got      interface{}
		expected interface{}
	}{
		{"MaxRetries", config.MaxRetries, 5},
		{"InitialBackoff", config.InitialBackoff, 100 * time.Millisecond},
		{"MaxBackoff", config.MaxBackoff, 30 * time.Second},
		{"BackoffMultiplier", config.BackoffMultiplier, 2.0},
		{"EnableAIFallback", config.EnableAIFallback, true},
		{"CollectFullContext", config.CollectFullContext, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("DefaultSelfCorrectionConfig().%s = %v, want %v", tt.name, tt.got, tt.expected)
			}
		})
	}
}

func TestSelfCorrectionEngine_AttemptCorrection(t *testing.T) {
	tests := []struct {
		name              string
		errorOutput       string
		stage             domain.ProtocolStage
		analyzeResult     *domain.ErrorDetails
		analyzeErr        error
		suggestResult     []*domain.CorrectionStep
		suggestErr        error
		applyResult       *domain.CorrectionResult
		applyErr          error
		canHandle         bool
		expectedSuccess   bool
		expectedRequireAI bool
		expectedAttempts  int
		wantErr           bool
	}{
		{
			name:        "successful correction on first attempt",
			errorOutput: "undefined: someFunc",
			stage:       domain.StageBuilding,
			analyzeResult: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeImport,
				SourceFile: "main.go",
				Message:    "undefined: someFunc",
			},
			suggestResult: []*domain.CorrectionStep{
				{Action: domain.ActionFixImport, Description: "Fix import"},
			},
			applyResult:       &domain.CorrectionResult{Success: true, Message: "Fixed"},
			canHandle:         true,
			expectedSuccess:   true,
			expectedRequireAI: false,
			expectedAttempts:  1,
		},
		{
			name:        "all attempts fail - requires AI",
			errorOutput: "complex error",
			stage:       domain.StageBuilding,
			analyzeResult: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "main.go",
				Message:    "complex error",
			},
			suggestResult: []*domain.CorrectionStep{
				{Action: domain.ActionFixSyntax, Description: "Fix syntax"},
			},
			applyResult:       &domain.CorrectionResult{Success: false, Message: "Failed"},
			canHandle:         true,
			expectedSuccess:   false,
			expectedRequireAI: true,
			expectedAttempts:  5, // MaxRetries default
		},
		{
			name:        "analyze error fails",
			errorOutput: "some error",
			stage:       domain.StageBuilding,
			analyzeErr:  errors.New("analysis failed"),
			wantErr:     true,
		},
		{
			name:        "no corrections suggested",
			errorOutput: "unknown error",
			stage:       domain.StageBuilding,
			analyzeResult: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "main.go",
				Message:    "unknown error",
			},
			suggestResult:     []*domain.CorrectionStep{},
			expectedSuccess:   false,
			expectedRequireAI: true,
			expectedAttempts:  5,
		},
		{
			name:        "correction engine cannot handle error",
			errorOutput: "type error",
			stage:       domain.StageBuilding,
			analyzeResult: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeTypeCheck,
				SourceFile: "main.go",
				Message:    "type error",
			},
			suggestResult: []*domain.CorrectionStep{
				{Action: domain.ActionFixType, Description: "Fix type"},
			},
			canHandle:         false,
			expectedSuccess:   false,
			expectedRequireAI: true,
			expectedAttempts:  5,
		},
		{
			name:        "apply correction returns error",
			errorOutput: "syntax error",
			stage:       domain.StageBuilding,
			analyzeResult: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeSyntax,
				SourceFile: "main.go",
				Message:    "syntax error",
			},
			suggestResult: []*domain.CorrectionStep{
				{Action: domain.ActionFixSyntax, Description: "Fix syntax"},
			},
			applyErr:          errors.New("apply failed"),
			canHandle:         true,
			expectedSuccess:   false,
			expectedRequireAI: true,
			expectedAttempts:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockAnalyzer := &mockErrorAnalyzer{
				analyzeResult:  tt.analyzeResult,
				analyzeErr:     tt.analyzeErr,
				suggestResult:  tt.suggestResult,
				suggestErr:     tt.suggestErr,
				classifyResult: domain.ErrorTypeCompilation,
			}

			mockEngine := &mockCorrectionEngine{
				applyResult:    tt.applyResult,
				applyErr:       tt.applyErr,
				canHandleValue: tt.canHandle,
			}

			mockFS := &mockFileSystem{
				readResult: []byte("package main\n\nfunc main() {}\n"),
			}

			// Use minimal config for faster tests
			config := &SelfCorrectionConfig{
				MaxRetries:         5,
				InitialBackoff:     1 * time.Millisecond,
				MaxBackoff:         10 * time.Millisecond,
				BackoffMultiplier:  1.5,
				EnableAIFallback:   true,
				CollectFullContext: true,
			}

			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				mockAnalyzer,
				mockEngine,
				mockFS,
				config,
			)

			ctx := context.Background()
			result, err := engine.AttemptCorrection(ctx, tt.errorOutput, tt.stage, "/test/project")

			if tt.wantErr {
				if err == nil {
					t.Errorf("AttemptCorrection() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("AttemptCorrection() unexpected error: %v", err)
				return
			}

			if result.Success != tt.expectedSuccess {
				t.Errorf("AttemptCorrection() Success = %v, want %v", result.Success, tt.expectedSuccess)
			}

			if result.RequiresAI != tt.expectedRequireAI {
				t.Errorf("AttemptCorrection() RequiresAI = %v, want %v", result.RequiresAI, tt.expectedRequireAI)
			}

			if result.TotalAttempts != tt.expectedAttempts {
				t.Errorf("AttemptCorrection() TotalAttempts = %v, want %v", result.TotalAttempts, tt.expectedAttempts)
			}
		})
	}
}

func TestSelfCorrectionEngine_CalculateBackoff(t *testing.T) {
	tests := []struct {
		name            string
		attempt         int
		initialBackoff  time.Duration
		maxBackoff      time.Duration
		multiplier      float64
		expectedBackoff time.Duration
	}{
		{
			name:            "first attempt",
			attempt:         1,
			initialBackoff:  100 * time.Millisecond,
			maxBackoff:      30 * time.Second,
			multiplier:      2.0,
			expectedBackoff: 100 * time.Millisecond,
		},
		{
			name:            "second attempt",
			attempt:         2,
			initialBackoff:  100 * time.Millisecond,
			maxBackoff:      30 * time.Second,
			multiplier:      2.0,
			expectedBackoff: 200 * time.Millisecond,
		},
		{
			name:            "third attempt",
			attempt:         3,
			initialBackoff:  100 * time.Millisecond,
			maxBackoff:      30 * time.Second,
			multiplier:      2.0,
			expectedBackoff: 400 * time.Millisecond,
		},
		{
			name:            "capped at max backoff",
			attempt:         10,
			initialBackoff:  100 * time.Millisecond,
			maxBackoff:      1 * time.Second,
			multiplier:      2.0,
			expectedBackoff: 1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SelfCorrectionConfig{
				InitialBackoff:    tt.initialBackoff,
				MaxBackoff:        tt.maxBackoff,
				BackoffMultiplier: tt.multiplier,
			}

			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				config,
			)

			got := engine.calculateBackoff(tt.attempt)
			if got != tt.expectedBackoff {
				t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.expectedBackoff)
			}
		})
	}
}

func TestSelfCorrectionEngine_ExtractSurroundingCode(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineNumber int
		before     int
		after      int
		wantEmpty  bool
		wantMarker bool
	}{
		{
			name:       "extract code around line 3",
			source:     "line1\nline2\nline3\nline4\nline5",
			lineNumber: 3,
			before:     1,
			after:      1,
			wantEmpty:  false,
			wantMarker: true,
		},
		{
			name:       "empty source",
			source:     "",
			lineNumber: 1,
			before:     2,
			after:      2,
			wantEmpty:  true,
		},
		{
			name:       "line number zero",
			source:     "line1\nline2",
			lineNumber: 0,
			before:     1,
			after:      1,
			wantEmpty:  true,
		},
		{
			name:       "line number exceeds source",
			source:     "line1\nline2",
			lineNumber: 10,
			before:     1,
			after:      1,
			wantEmpty:  true,
		},
		{
			name:       "first line",
			source:     "line1\nline2\nline3",
			lineNumber: 1,
			before:     2,
			after:      2,
			wantEmpty:  false,
			wantMarker: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				nil,
			)

			got := engine.extractSurroundingCode(tt.source, tt.lineNumber, tt.before, tt.after)

			if tt.wantEmpty && got != "" {
				t.Errorf("extractSurroundingCode() = %q, want empty", got)
			}

			if !tt.wantEmpty && got == "" {
				t.Errorf("extractSurroundingCode() = empty, want non-empty")
			}

			if tt.wantMarker && got != "" {
				if len(got) < 2 || got[0:2] != "> " && !containsMarker(got) {
					// Check if marker exists somewhere in output
					if !containsMarker(got) {
						t.Errorf("extractSurroundingCode() missing '> ' marker for error line")
					}
				}
			}
		})
	}
}

func containsMarker(s string) bool {
	return len(s) > 0 && (s[0] == '>' || containsSubstring(s, "\n> "))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSelfCorrectionEngine_ShouldRetry(t *testing.T) {
	tests := []struct {
		name       string
		attempt    int
		maxRetries int
		err        error
		expected   bool
	}{
		{
			name:       "should retry - under max",
			attempt:    2,
			maxRetries: 5,
			err:        errors.New("some error"),
			expected:   true,
		},
		{
			name:       "should not retry - at max",
			attempt:    5,
			maxRetries: 5,
			err:        errors.New("some error"),
			expected:   false,
		},
		{
			name:       "should not retry - context canceled",
			attempt:    1,
			maxRetries: 5,
			err:        context.Canceled,
			expected:   false,
		},
		{
			name:       "should not retry - deadline exceeded",
			attempt:    1,
			maxRetries: 5,
			err:        context.DeadlineExceeded,
			expected:   false,
		},
		{
			name:       "should retry - nil error",
			attempt:    1,
			maxRetries: 5,
			err:        nil,
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &SelfCorrectionConfig{
				MaxRetries: tt.maxRetries,
			}

			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				config,
			)

			got := engine.ShouldRetry(tt.attempt, tt.err)
			if got != tt.expected {
				t.Errorf("ShouldRetry(%d, %v) = %v, want %v", tt.attempt, tt.err, got, tt.expected)
			}
		})
	}
}

func TestSelfCorrectionEngine_FormatAIPrompt(t *testing.T) {
	tests := []struct {
		name           string
		errCtx         *ErrorContext
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "full context with all fields",
			errCtx: &ErrorContext{
				OriginalError: &domain.ErrorDetails{
					ErrorType:   domain.ErrorTypeImport,
					Stage:       domain.StageBuilding,
					SourceFile:  "main.go",
					LineNumber:  10,
					Message:     "undefined: someFunc",
					Suggestions: []string{"Add import statement"},
				},
				AttemptedFixes: []AttemptedFix{
					{Attempt: 1, Description: "Fix import", Success: false, Error: "failed"},
				},
				SurroundingCode: "> 10 | someFunc()",
				TotalAttempts:   1,
			},
			wantContains: []string{
				"Error Correction Request",
				"import",
				"building",
				"main.go",
				"undefined: someFunc",
				"Previously Attempted Fixes",
				"Fix import",
				"Code Context",
				"someFunc()",
				"Analyzer Suggestions",
				"Add import statement",
			},
		},
		{
			name: "minimal context",
			errCtx: &ErrorContext{
				OriginalError: &domain.ErrorDetails{
					ErrorType: domain.ErrorTypeCompilation,
					Message:   "error",
				},
				AttemptedFixes: []AttemptedFix{},
			},
			wantContains: []string{
				"Error Correction Request",
				"compilation",
				"error",
			},
			wantNotContain: []string{
				"Previously Attempted Fixes",
				"Code Context",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				nil,
			)

			got := engine.FormatAIPrompt(tt.errCtx)

			for _, want := range tt.wantContains {
				if !containsSubstring(got, want) {
					t.Errorf("FormatAIPrompt() missing expected content: %q", want)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if containsSubstring(got, notWant) {
					t.Errorf("FormatAIPrompt() contains unexpected content: %q", notWant)
				}
			}
		})
	}
}

func TestSelfCorrectionEngine_CreateAIAssistanceRequest(t *testing.T) {
	tests := []struct {
		name             string
		result           *SelfCorrectionResult
		projectPath      string
		expectNil        bool
		expectedPriority string
	}{
		{
			name: "creates request when AI required",
			result: &SelfCorrectionResult{
				RequiresAI: true,
				AIContext: &ErrorContext{
					OriginalError: &domain.ErrorDetails{
						Severity: "error",
					},
				},
			},
			projectPath:      "/test/project",
			expectNil:        false,
			expectedPriority: "high",
		},
		{
			name: "creates request with normal priority",
			result: &SelfCorrectionResult{
				RequiresAI: true,
				AIContext: &ErrorContext{
					OriginalError: &domain.ErrorDetails{
						Severity: "warning",
					},
				},
			},
			projectPath:      "/test/project",
			expectNil:        false,
			expectedPriority: "normal",
		},
		{
			name: "returns nil when AI not required",
			result: &SelfCorrectionResult{
				RequiresAI: false,
			},
			projectPath: "/test/project",
			expectNil:   true,
		},
		{
			name: "returns nil when no AI context",
			result: &SelfCorrectionResult{
				RequiresAI: true,
				AIContext:  nil,
			},
			projectPath: "/test/project",
			expectNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				nil,
			)

			got := engine.CreateAIAssistanceRequest(tt.result, tt.projectPath)

			if tt.expectNil {
				if got != nil {
					t.Errorf("CreateAIAssistanceRequest() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Errorf("CreateAIAssistanceRequest() = nil, want non-nil")
				return
			}

			if got.Priority != tt.expectedPriority {
				t.Errorf("CreateAIAssistanceRequest().Priority = %v, want %v", got.Priority, tt.expectedPriority)
			}

			if got.ProjectPath != tt.projectPath {
				t.Errorf("CreateAIAssistanceRequest().ProjectPath = %v, want %v", got.ProjectPath, tt.projectPath)
			}
		})
	}
}

func TestSelfCorrectionEngine_ContextCancellation(t *testing.T) {
	mockAnalyzer := &mockErrorAnalyzer{
		analyzeResult: &domain.ErrorDetails{
			ErrorType:  domain.ErrorTypeCompilation,
			SourceFile: "main.go",
			Message:    "error",
		},
		suggestResult: []*domain.CorrectionStep{
			{Action: domain.ActionFixSyntax, Description: "Fix"},
		},
	}

	mockEngine := &mockCorrectionEngine{
		applyResult:    &domain.CorrectionResult{Success: false},
		canHandleValue: true,
	}

	config := &SelfCorrectionConfig{
		MaxRetries:        5,
		InitialBackoff:    100 * time.Millisecond,
		MaxBackoff:        1 * time.Second,
		BackoffMultiplier: 2.0,
		EnableAIFallback:  true,
	}

	engine := NewSelfCorrectionEngine(
		&mockLogger{},
		mockAnalyzer,
		mockEngine,
		&mockFileSystem{},
		config,
	)

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel context after short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	_, err := engine.AttemptCorrection(ctx, "error", domain.StageBuilding, "/test")

	if err != context.Canceled {
		t.Errorf("AttemptCorrection() with cancelled context should return context.Canceled, got %v", err)
	}
}

func TestSelfCorrectionEngine_UpdateConfig(t *testing.T) {
	engine := NewSelfCorrectionEngine(
		&mockLogger{},
		&mockErrorAnalyzer{},
		&mockCorrectionEngine{},
		&mockFileSystem{},
		nil,
	)

	// Check default config
	defaultConfig := engine.GetConfig()
	if defaultConfig.MaxRetries != 5 {
		t.Errorf("Default MaxRetries = %d, want 5", defaultConfig.MaxRetries)
	}

	// Update config
	newConfig := &SelfCorrectionConfig{
		MaxRetries:        10,
		InitialBackoff:    200 * time.Millisecond,
		MaxBackoff:        60 * time.Second,
		BackoffMultiplier: 3.0,
		EnableAIFallback:  false,
	}

	engine.UpdateConfig(newConfig)

	updatedConfig := engine.GetConfig()
	if updatedConfig.MaxRetries != 10 {
		t.Errorf("Updated MaxRetries = %d, want 10", updatedConfig.MaxRetries)
	}
	if updatedConfig.EnableAIFallback != false {
		t.Errorf("Updated EnableAIFallback = %v, want false", updatedConfig.EnableAIFallback)
	}

	// Test nil config update (should not change)
	engine.UpdateConfig(nil)
	if engine.GetConfig().MaxRetries != 10 {
		t.Errorf("Config changed after nil update")
	}
}

func TestSelfCorrectionEngine_FindRelatedFiles(t *testing.T) {
	tests := []struct {
		name          string
		errorDetails  *domain.ErrorDetails
		expectedFiles []string
	}{
		{
			name: "import error - Go files",
			errorDetails: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeImport,
				SourceFile: "main.go",
			},
			expectedFiles: []string{"main.go", "go.mod", "go.sum", "package.json"},
		},
		{
			name: "type error",
			errorDetails: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeTypeCheck,
				SourceFile: "handler.go",
			},
			expectedFiles: []string{"handler.go", "types.go", "types.ts", "index.d.ts"},
		},
		{
			name: "compilation error - no extra files",
			errorDetails: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "main.go",
			},
			expectedFiles: []string{"main.go"},
		},
		{
			name: "no source file",
			errorDetails: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "",
			},
			expectedFiles: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewSelfCorrectionEngine(
				&mockLogger{},
				&mockErrorAnalyzer{},
				&mockCorrectionEngine{},
				&mockFileSystem{},
				nil,
			)

			got := engine.findRelatedFiles(tt.errorDetails, "/test")

			if len(got) != len(tt.expectedFiles) {
				t.Errorf("findRelatedFiles() returned %d files, want %d", len(got), len(tt.expectedFiles))
				return
			}

			for i, expected := range tt.expectedFiles {
				if got[i] != expected {
					t.Errorf("findRelatedFiles()[%d] = %s, want %s", i, got[i], expected)
				}
			}
		})
	}
}


// === Tests for SelfCorrectionService ===

func TestNewSelfCorrectionService(t *testing.T) {
	tests := []struct {
		name       string
		maxRetries int
		expected   int
	}{
		{"default retries", 0, 3},
		{"negative retries", -1, 3},
		{"custom retries", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewSelfCorrectionService(
				&mockLogger{},
				&mockCorrectionEngine{},
				&mockErrorAnalyzer{},
				&mockFileSystem{},
				tt.maxRetries,
			)
			if svc.maxRetries != tt.expected {
				t.Errorf("maxRetries = %d, want %d", svc.maxRetries, tt.expected)
			}
		})
	}
}

func TestSelfCorrectionService_AttemptCorrection_AutoFixSuccess(t *testing.T) {
	mockEngine := &mockCorrectionEngine{
		applyResult:    &domain.CorrectionResult{Success: true, Message: "Fixed"},
		canHandleValue: true,
	}
	mockAnalyzer := &mockErrorAnalyzer{
		suggestResult: []*domain.CorrectionStep{
			{Action: domain.ActionFixImport, Description: "Fix import"},
		},
	}

	svc := NewSelfCorrectionService(
		&mockLogger{},
		mockEngine,
		mockAnalyzer,
		&mockFileSystem{},
		3,
	)

	errors := []*domain.ErrorDetails{
		{ErrorType: domain.ErrorTypeImport, SourceFile: "main.go", Message: "undefined: fmt"},
	}

	result, err := svc.AttemptCorrection(context.Background(), errors, "/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !result.AllFixed {
		t.Error("expected AllFixed to be true")
	}
	if len(result.RemainingErrors) != 0 {
		t.Errorf("expected no remaining errors, got %d", len(result.RemainingErrors))
	}
	if len(result.Attempts) != 1 {
		t.Errorf("expected 1 attempt, got %d", len(result.Attempts))
	}
	if !result.Attempts[0].AutoFixSuccess {
		t.Error("expected AutoFixSuccess to be true")
	}
}

func TestSelfCorrectionService_AttemptCorrection_NeedsAIHelp(t *testing.T) {
	mockEngine := &mockCorrectionEngine{
		canHandleValue: false, // Cannot handle
	}
	mockAnalyzer := &mockErrorAnalyzer{}

	svc := NewSelfCorrectionService(
		&mockLogger{},
		mockEngine,
		mockAnalyzer,
		&mockFileSystem{readResult: []byte("package main\n\nfunc main() {}\n")},
		3,
	)

	errors := []*domain.ErrorDetails{
		{ErrorType: domain.ErrorTypeCompilation, SourceFile: "main.go", LineNumber: 5, Message: "complex error"},
	}

	result, err := svc.AttemptCorrection(context.Background(), errors, "/test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.AllFixed {
		t.Error("expected AllFixed to be false")
	}
	if len(result.RemainingErrors) != 1 {
		t.Errorf("expected 1 remaining error, got %d", len(result.RemainingErrors))
	}
	if !result.Attempts[0].NeedsAIHelp {
		t.Error("expected NeedsAIHelp to be true")
	}
	if result.Attempts[0].SuggestedPrompt == "" {
		t.Error("expected SuggestedPrompt to be non-empty")
	}
}

func TestSelfCorrectionService_GenerateAIPrompt(t *testing.T) {
	svc := NewSelfCorrectionService(
		&mockLogger{},
		nil,
		nil,
		&mockFileSystem{readResult: []byte("line1\nline2\nline3\nline4\nline5\n")},
		3,
	)

	err := &domain.ErrorDetails{
		ErrorType:  domain.ErrorTypeTypeCheck,
		SourceFile: "test.ts",
		LineNumber: 3,
		Message:    "Property 'user' does not exist on type 'Session'",
	}

	prompt := svc.GenerateAIPrompt(err)

	if prompt == "" {
		t.Error("expected non-empty prompt")
	}
	if !containsSubstring(prompt, "type error") {
		t.Error("expected prompt to mention type error")
	}
	if !containsSubstring(prompt, "test.ts") {
		t.Error("expected prompt to mention file name")
	}
}

func TestSelfCorrectionService_ShouldRetryWithAI(t *testing.T) {
	svc := NewSelfCorrectionService(&mockLogger{}, nil, nil, nil, 3)

	tests := []struct {
		name     string
		result   *CorrectionResult
		expected bool
	}{
		{
			name:     "all fixed - no retry",
			result:   &CorrectionResult{AllFixed: true},
			expected: false,
		},
		{
			name: "has AI prompts - should retry",
			result: &CorrectionResult{
				AllFixed:  false,
				AIPrompts: []string{"Fix this error"},
			},
			expected: true,
		},
		{
			name: "not fixed but no prompts - no retry",
			result: &CorrectionResult{
				AllFixed:  false,
				AIPrompts: []string{},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.ShouldRetryWithAI(tt.result)
			if got != tt.expected {
				t.Errorf("ShouldRetryWithAI() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSelfCorrectionService_GetCombinedAIPrompt(t *testing.T) {
	svc := NewSelfCorrectionService(&mockLogger{}, nil, nil, nil, 3)

	tests := []struct {
		name        string
		result      *CorrectionResult
		wantEmpty   bool
		wantContain string
	}{
		{
			name:      "no remaining errors",
			result:    &CorrectionResult{RemainingErrors: []*domain.ErrorDetails{}},
			wantEmpty: true,
		},
		{
			name: "has remaining errors",
			result: &CorrectionResult{
				RemainingErrors: []*domain.ErrorDetails{
					{ErrorType: domain.ErrorTypeImport, SourceFile: "main.go", Message: "undefined: fmt"},
				},
			},
			wantEmpty:   false,
			wantContain: "main.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := svc.GetCombinedAIPrompt(tt.result)
			if tt.wantEmpty && got != "" {
				t.Errorf("GetCombinedAIPrompt() = %q, want empty", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("GetCombinedAIPrompt() = empty, want non-empty")
			}
			if tt.wantContain != "" && !containsSubstring(got, tt.wantContain) {
				t.Errorf("GetCombinedAIPrompt() missing %q", tt.wantContain)
			}
		})
	}
}

func TestExtractCodeContext(t *testing.T) {
	tests := []struct {
		name       string
		source     string
		lineNumber int
		before     int
		after      int
		wantEmpty  bool
		wantMarker bool
	}{
		{
			name:       "extract around line 3",
			source:     "line1\nline2\nline3\nline4\nline5",
			lineNumber: 3,
			before:     1,
			after:      1,
			wantEmpty:  false,
			wantMarker: true,
		},
		{
			name:       "empty source",
			source:     "",
			lineNumber: 1,
			wantEmpty:  true,
		},
		{
			name:       "line number zero",
			source:     "line1\nline2",
			lineNumber: 0,
			wantEmpty:  true,
		},
		{
			name:       "line exceeds length",
			source:     "line1\nline2",
			lineNumber: 10,
			wantEmpty:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractCodeContext(tt.source, tt.lineNumber, tt.before, tt.after)
			if tt.wantEmpty && got != "" {
				t.Errorf("extractCodeContext() = %q, want empty", got)
			}
			if !tt.wantEmpty && got == "" {
				t.Error("extractCodeContext() = empty, want non-empty")
			}
			if tt.wantMarker && !containsSubstring(got, "> ") {
				t.Error("extractCodeContext() missing '> ' marker")
			}
		})
	}
}

func TestCorrectionAttempt_Fields(t *testing.T) {
	attempt := CorrectionAttempt{
		Error: &domain.ErrorDetails{
			ErrorType: domain.ErrorTypeImport,
			Message:   "test error",
		},
		AutoFixApplied:  true,
		AutoFixSuccess:  false,
		NeedsAIHelp:     true,
		SuggestedPrompt: "Fix this",
		AttemptNumber:   2,
	}

	if attempt.Error.ErrorType != domain.ErrorTypeImport {
		t.Error("Error field not set correctly")
	}
	if !attempt.AutoFixApplied {
		t.Error("AutoFixApplied should be true")
	}
	if attempt.AutoFixSuccess {
		t.Error("AutoFixSuccess should be false")
	}
	if !attempt.NeedsAIHelp {
		t.Error("NeedsAIHelp should be true")
	}
	if attempt.SuggestedPrompt != "Fix this" {
		t.Error("SuggestedPrompt not set correctly")
	}
	if attempt.AttemptNumber != 2 {
		t.Error("AttemptNumber should be 2")
	}
}

func TestCorrectionResult_Fields(t *testing.T) {
	result := CorrectionResult{
		Attempts: []CorrectionAttempt{
			{AttemptNumber: 1},
			{AttemptNumber: 2},
		},
		AllFixed: false,
		RemainingErrors: []*domain.ErrorDetails{
			{Message: "error1"},
		},
		AIPrompts:     []string{"prompt1", "prompt2"},
		TotalAttempts: 2,
	}

	if len(result.Attempts) != 2 {
		t.Errorf("Attempts length = %d, want 2", len(result.Attempts))
	}
	if result.AllFixed {
		t.Error("AllFixed should be false")
	}
	if len(result.RemainingErrors) != 1 {
		t.Errorf("RemainingErrors length = %d, want 1", len(result.RemainingErrors))
	}
	if len(result.AIPrompts) != 2 {
		t.Errorf("AIPrompts length = %d, want 2", len(result.AIPrompts))
	}
	if result.TotalAttempts != 2 {
		t.Errorf("TotalAttempts = %d, want 2", result.TotalAttempts)
	}
}
