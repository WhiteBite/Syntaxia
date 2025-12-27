package verification

import (
	"context"
	"testing"
	"time"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
	"syntaxia/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestService() (*ChatIntegrationService, *testutils.MockLogger, *testutils.MockBuildService, *testutils.MockTestService, *testutils.MockStaticAnalyzerService) {
	logger := testutils.NewMockLogger()
	logger.On("Info", mock.Anything).Return()
	logger.On("Warning", mock.Anything).Return()
	logger.On("Error", mock.Anything).Return()
	logger.On("Debug", mock.Anything).Return()

	buildService := &testutils.MockBuildService{}
	testService := &testutils.MockTestService{}
	staticAnalyzer := &testutils.MockStaticAnalyzerService{}

	service := NewChatIntegrationService(logger, buildService, testService, staticAnalyzer)

	return service, logger, buildService, testService, staticAnalyzer
}

func TestNewChatIntegrationService(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	assert.NotNil(t, service)
	assert.NotNil(t, service.baselines)
}

func TestDetectLanguageFromFile(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		file     string
		expected string
	}{
		{"Go file", "main.go", "go"},
		{"TypeScript file", "app.ts", "typescript"},
		{"TSX file", "component.tsx", "typescript"},
		{"JavaScript file", "script.js", "javascript"},
		{"JSX file", "component.jsx", "javascript"},
		{"Python file", "main.py", "python"},
		{"Java file", "Main.java", "java"},
		{"Rust file", "main.rs", "rust"},
		{"Kotlin file", "Main.kt", "kotlin"},
		{"C# file", "Program.cs", "csharp"},
		{"C++ file", "main.cpp", "cpp"},
		{"C file", "main.c", "cpp"},
		{"Ruby file", "app.rb", "ruby"},
		{"PHP file", "index.php", "php"},
		{"Swift file", "main.swift", "swift"},
		{"Dart file", "main.dart", "dart"},
		{"Unknown file", "readme.md", ""},
		{"No extension", "Makefile", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.detectLanguageFromFile(tt.file)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectLanguagesFromFiles(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		files    []string
		fallback []string
		expected []string
	}{
		{
			name:     "Single Go file",
			files:    []string{"main.go"},
			fallback: []string{"typescript"},
			expected: []string{"go"},
		},
		{
			name:     "Multiple languages",
			files:    []string{"main.go", "app.ts", "script.py"},
			fallback: []string{},
			expected: []string{"go", "typescript", "python"},
		},
		{
			name:     "Empty files uses fallback",
			files:    []string{},
			fallback: []string{"go", "typescript"},
			expected: []string{"go", "typescript"},
		},
		{
			name:     "Unknown files uses fallback",
			files:    []string{"readme.md", "Makefile"},
			fallback: []string{"go"},
			expected: []string{"go"},
		},
		{
			name:     "Duplicate languages deduplicated",
			files:    []string{"main.go", "test.go", "util.go"},
			fallback: []string{},
			expected: []string{"go"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.detectLanguagesFromFiles(tt.files, tt.fallback)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestFingerprintKey(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name        string
		fingerprint IssueFingerprint
		expected    string
	}{
		{
			name: "Full fingerprint",
			fingerprint: IssueFingerprint{
				File:    "main.go",
				Line:    10,
				Code:    "E001",
				Message: "error message",
			},
			expected: "main.go:10:E001:error message",
		},
		{
			name: "Empty code",
			fingerprint: IssueFingerprint{
				File:    "app.ts",
				Line:    5,
				Code:    "",
				Message: "warning",
			},
			expected: "app.ts:5::warning",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.fingerprintKey(tt.fingerprint)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasErrors(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		issues   []VerificationIssue
		expected bool
	}{
		{
			name:     "Empty issues",
			issues:   []VerificationIssue{},
			expected: false,
		},
		{
			name: "Only warnings",
			issues: []VerificationIssue{
				{Severity: "warning", Message: "warning 1"},
				{Severity: "warning", Message: "warning 2"},
			},
			expected: false,
		},
		{
			name: "Has error",
			issues: []VerificationIssue{
				{Severity: "warning", Message: "warning"},
				{Severity: "error", Message: "error"},
			},
			expected: true,
		},
		{
			name: "Only errors",
			issues: []VerificationIssue{
				{Severity: "error", Message: "error 1"},
				{Severity: "error", Message: "error 2"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.hasErrors(tt.issues)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterNewErrors(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name          string
		current       []VerificationIssue
		baseline      *BaselineState
		expectedNew   int
		expectedFixed int
	}{
		{
			name:    "No baseline issues, all current are new",
			current: []VerificationIssue{
				{Type: "type", File: "main.go", Line: 10, Message: "error 1"},
				{Type: "lint", File: "app.go", Line: 20, Message: "warning 1"},
			},
			baseline: &BaselineState{
				TypeErrors: []IssueFingerprint{},
				LintIssues: []IssueFingerprint{},
			},
			expectedNew:   2,
			expectedFixed: 0,
		},
		{
			name:    "All current in baseline, none new",
			current: []VerificationIssue{
				{Type: "type", File: "main.go", Line: 10, Message: "error 1"},
			},
			baseline: &BaselineState{
				TypeErrors: []IssueFingerprint{
					{File: "main.go", Line: 10, Message: "error 1"},
				},
			},
			expectedNew:   0,
			expectedFixed: 0,
		},
		{
			name:    "Some fixed, some new",
			current: []VerificationIssue{
				{Type: "type", File: "main.go", Line: 15, Message: "new error"},
			},
			baseline: &BaselineState{
				TypeErrors: []IssueFingerprint{
					{File: "main.go", Line: 10, Message: "old error"},
				},
			},
			expectedNew:   1,
			expectedFixed: 1,
		},
		{
			name:    "Empty current, all baseline fixed",
			current: []VerificationIssue{},
			baseline: &BaselineState{
				TypeErrors: []IssueFingerprint{
					{File: "main.go", Line: 10, Message: "fixed error"},
				},
				LintIssues: []IssueFingerprint{
					{File: "app.go", Line: 5, Message: "fixed warning"},
				},
			},
			expectedNew:   0,
			expectedFixed: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newErrors, fixedErrors := service.filterNewErrors(tt.current, tt.baseline)
			assert.Len(t, newErrors, tt.expectedNew)
			assert.Len(t, fixedErrors, tt.expectedFixed)
		})
	}
}

func TestDetermineSuccess(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		result   *ChatVerificationResult
		expected bool
	}{
		{
			name: "All passed",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: true},
				TestResult:  &StepResult{Success: true},
				LintResult:  &StepResult{Success: true},
				NewErrors:   []VerificationIssue{},
			},
			expected: true,
		},
		{
			name: "Build failed",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: false},
				TestResult:  &StepResult{Skipped: true},
				LintResult:  &StepResult{Success: true},
			},
			expected: false,
		},
		{
			name: "Tests failed",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: true},
				TestResult:  &StepResult{Success: false},
				LintResult:  &StepResult{Success: true},
			},
			expected: false,
		},
		{
			name: "Timeout",
			result: &ChatVerificationResult{
				Timeout:     true,
				BuildResult: &StepResult{Success: true},
			},
			expected: false,
		},
		{
			name: "New errors with severity error",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: true},
				TestResult:  &StepResult{Skipped: true},
				LintResult:  &StepResult{Success: true},
				NewErrors: []VerificationIssue{
					{Severity: "error", Message: "new error"},
				},
			},
			expected: false,
		},
		{
			name: "New warnings only (success)",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: true},
				TestResult:  &StepResult{Skipped: true},
				LintResult:  &StepResult{Success: true},
				NewErrors: []VerificationIssue{
					{Severity: "warning", Message: "new warning"},
				},
			},
			expected: true,
		},
		{
			name: "Skipped steps are ignored",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Skipped: true},
				TestResult:  &StepResult{Skipped: true},
				LintResult:  &StepResult{Skipped: true},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.determineSuccess(tt.result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateSummary(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		result   *ChatVerificationResult
		contains []string
	}{
		{
			name: "All passed",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Success: true},
				TestResult:  &StepResult{Success: true},
				LintResult:  &StepResult{Success: true},
			},
			contains: []string{"✅ Build passed", "✅ Tests passed", "✅ Lint passed"},
		},
		{
			name: "Build failed",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{
					Success: false,
					Issues:  []VerificationIssue{{Message: "error"}},
				},
			},
			contains: []string{"❌ Build failed"},
		},
		{
			name: "Timeout",
			result: &ChatVerificationResult{
				Timeout: true,
			},
			contains: []string{"⏱️ Verification timed out"},
		},
		{
			name: "New and fixed errors",
			result: &ChatVerificationResult{
				NewErrors:   []VerificationIssue{{Message: "new"}},
				FixedErrors: []VerificationIssue{{Message: "fixed1"}, {Message: "fixed2"}},
			},
			contains: []string{"🆕 1 new issues", "🔧 2 issues fixed"},
		},
		{
			name: "Skipped steps not shown",
			result: &ChatVerificationResult{
				BuildResult: &StepResult{Skipped: true},
				TestResult:  &StepResult{Skipped: true},
				LintResult:  &StepResult{Skipped: true},
			},
			contains: []string{"No verification steps executed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := service.generateSummary(tt.result)
			for _, expected := range tt.contains {
				assert.Contains(t, summary, expected)
			}
		})
	}
}

func TestGetDefaultTimeout(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		mode     VerificationMode
		expected time.Duration
	}{
		{"Fast mode", ModeFast, DefaultFastModeTimeout},
		{"Incremental mode", ModeIncremental, DefaultBuildTimeout + DefaultTestTimeout},
		{"Full mode", ModeFull, DefaultBuildTimeout + DefaultTestTimeout + DefaultAnalysisTimeout},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getDefaultTimeout(tt.mode)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCollectAllIssues(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	tests := []struct {
		name     string
		build    *StepResult
		test     *StepResult
		lint     *StepResult
		expected int
	}{
		{
			name:     "All nil",
			build:    nil,
			test:     nil,
			lint:     nil,
			expected: 0,
		},
		{
			name: "All skipped",
			build: &StepResult{Skipped: true, Issues: []VerificationIssue{{Message: "ignored"}}},
			test:  &StepResult{Skipped: true},
			lint:  &StepResult{Skipped: true},
			expected: 0,
		},
		{
			name: "Build issues only",
			build: &StepResult{
				Issues: []VerificationIssue{
					{Type: "type", Message: "error 1"},
					{Type: "type", Message: "error 2"},
				},
			},
			test:     &StepResult{Skipped: true},
			lint:     &StepResult{Skipped: true},
			expected: 2,
		},
		{
			name: "All steps have issues",
			build: &StepResult{
				Issues: []VerificationIssue{{Type: "type", Message: "build error"}},
			},
			test: &StepResult{
				Issues: []VerificationIssue{{Type: "test", Message: "test failure"}},
			},
			lint: &StepResult{
				Issues: []VerificationIssue{
					{Type: "lint", Message: "lint 1"},
					{Type: "lint", Message: "lint 2"},
				},
			},
			expected: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.collectAllIssues(tt.build, tt.test, tt.lint)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestBaselineManagement(t *testing.T) {
	service, _, _, _, _ := setupTestService()

	projectPath := "/test/project"

	// Initially no baseline
	baseline := service.GetBaseline(projectPath)
	assert.Nil(t, baseline)

	// Store baseline manually
	service.baselineMu.Lock()
	service.baselines[projectPath] = &BaselineState{
		CapturedAt: time.Now(),
		TypeErrors: []IssueFingerprint{{File: "test.go", Line: 1}},
	}
	service.baselineMu.Unlock()

	// Get baseline
	baseline = service.GetBaseline(projectPath)
	assert.NotNil(t, baseline)
	assert.Len(t, baseline.TypeErrors, 1)

	// Clear baseline
	service.ClearBaseline(projectPath)
	baseline = service.GetBaseline(projectPath)
	assert.Nil(t, baseline)
}

func TestCaptureBaseline(t *testing.T) {
	service, _, buildService, _, staticAnalyzer := setupTestService()

	ctx := context.Background()
	projectPath := "/test/project"
	languages := []string{"go"}

	// Setup mocks
	buildService.On("TypeCheck", ctx, projectPath, "go").Return(&domain.TypeCheckResult{
		Success: true,
		Issues: []*domain.TypeIssue{
			{File: "main.go", Line: 10, Message: "existing error", Severity: "error"},
		},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", ctx, projectPath, languages).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {
				Success: true,
				Issues: []*domain.StaticIssue{
					{File: "main.go", Line: 20, Message: "lint warning", Severity: "warning"},
				},
			},
		},
	}, nil)

	// Capture baseline
	baseline, err := service.CaptureBaseline(ctx, projectPath, languages)

	assert.NoError(t, err)
	assert.NotNil(t, baseline)
	assert.Len(t, baseline.TypeErrors, 1)
	assert.Len(t, baseline.LintIssues, 1)
	assert.Equal(t, "main.go", baseline.TypeErrors[0].File)
	assert.Equal(t, 10, baseline.TypeErrors[0].Line)

	// Verify baseline is stored
	storedBaseline := service.GetBaseline(projectPath)
	assert.NotNil(t, storedBaseline)
	assert.Equal(t, baseline, storedBaseline)

	buildService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestRunAfterChanges_FastMode(t *testing.T) {
	service, _, buildService, _, staticAnalyzer := setupTestService()

	ctx := context.Background()
	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go", "util.go"},
		Mode:         ModeFast,
		SkipTests:    true,
	}

	// Setup mocks
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: true,
		Issues:  []*domain.TypeIssue{},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, []string{"go"}).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil)

	// Run verification
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ModeFast, result.Mode)
	assert.True(t, result.Success)
	assert.ElementsMatch(t, config.ChangedFiles, result.ChangedFiles)

	buildService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestRunAfterChanges_WithNewErrors(t *testing.T) {
	service, _, buildService, _, staticAnalyzer := setupTestService()

	ctx := context.Background()
	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go"},
		Mode:         ModeFast,
		SkipTests:    true,
		BaselineState: &BaselineState{
			TypeErrors: []IssueFingerprint{},
			LintIssues: []IssueFingerprint{},
		},
	}

	// Setup mocks - return new error
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: false,
		Issues: []*domain.TypeIssue{
			{File: "main.go", Line: 15, Message: "undefined: foo", Severity: "error"},
		},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, []string{"go"}).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil)

	// Run verification
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	assert.Len(t, result.NewErrors, 1)
	assert.Equal(t, "main.go", result.NewErrors[0].File)
	assert.Equal(t, 15, result.NewErrors[0].Line)

	buildService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestRunAfterChanges_IncrementalMode(t *testing.T) {
	service, _, buildService, testService, staticAnalyzer := setupTestService()

	ctx := context.Background()
	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go"},
		Mode:         ModeIncremental,
	}

	// Setup mocks
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: true,
		Issues:  []*domain.TypeIssue{},
	}, nil)

	testService.On("RunTargetedTests", mock.Anything, mock.Anything, config.ChangedFiles).Return([]*domain.TestResult{
		{Success: true, TestPath: "main_test.go", TestName: "TestMain"},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, config.Languages).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil)

	// Run verification
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ModeIncremental, result.Mode)
	assert.True(t, result.Success)
	assert.NotNil(t, result.TestResult)
	assert.False(t, result.TestResult.Skipped)

	buildService.AssertExpectations(t)
	testService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestRunAfterChanges_FullMode(t *testing.T) {
	service, _, buildService, testService, staticAnalyzer := setupTestService()

	ctx := context.Background()
	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go"},
		Mode:         ModeFull,
	}

	// Setup mocks
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: true,
		Issues:  []*domain.TypeIssue{},
	}, nil)

	testService.On("RunSmokeTests", mock.Anything, config.ProjectPath, "go").Return([]*domain.TestResult{
		{Success: true, TestPath: "main_test.go", TestName: "TestMain"},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, config.Languages).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil)

	// Run verification
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, ModeFull, result.Mode)
	assert.True(t, result.Success)

	buildService.AssertExpectations(t)
	testService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestRunAfterChanges_BuildFails(t *testing.T) {
	service, _, buildService, _, staticAnalyzer := setupTestService()

	ctx := context.Background()
	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go"},
		Mode:         ModeFull,
		SkipTests:    true,
	}

	// Setup mocks - build fails
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: false,
		Issues: []*domain.TypeIssue{
			{File: "main.go", Line: 10, Message: "undefined: foo", Severity: "error"},
		},
	}, nil)

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, config.Languages).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil)

	// Run verification
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.False(t, result.Success)
	assert.NotNil(t, result.BuildResult)
	assert.False(t, result.BuildResult.Success)
	assert.Len(t, result.BuildResult.Issues, 1)

	buildService.AssertExpectations(t)
	staticAnalyzer.AssertExpectations(t)
}

func TestVerifyWithTimeout_Cancellation(t *testing.T) {
	service, _, buildService, _, staticAnalyzer := setupTestService()

	// Create a context with very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	config := &ChatIntegrationConfig{
		ProjectPath:  "/test/project",
		Languages:    []string{"go"},
		ChangedFiles: []string{"main.go"},
		Mode:         ModeFast,
		Timeout:      1 * time.Millisecond,
	}

	// Setup mocks - they may or may not be called depending on timing
	buildService.On("TypeCheck", mock.Anything, config.ProjectPath, "go").Return(&domain.TypeCheckResult{
		Success: true,
		Issues:  []*domain.TypeIssue{},
	}, nil).Maybe()

	staticAnalyzer.On("AnalyzeProject", mock.Anything, config.ProjectPath, []string{"go"}).Return(&domain.StaticAnalysisReport{
		Results: map[string]*domain.StaticAnalysisResult{
			"go": {Success: true, Issues: []*domain.StaticIssue{}},
		},
	}, nil).Maybe()

	// Run verification with short timeout context
	result, err := service.RunAfterChanges(ctx, config)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	// Result should complete (either with timeout flag or normally)
}

// ErrorFilter tests
func TestErrorFilter_NewErrorFilter(t *testing.T) {
	filter := NewErrorFilter()

	assert.NotNil(t, filter)
	assert.NotNil(t, filter.baselineErrors)
	assert.NotNil(t, filter.baselineKeys)
	assert.Equal(t, 0, filter.GetBaselineCount())
}

func TestErrorFilter_SetBaseline(t *testing.T) {
	filter := NewErrorFilter()

	errors := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "error 1"},
		{Type: "lint", File: "app.go", Line: 20, Message: "warning 1"},
		{Type: "type", File: "main.go", Line: 15, Message: "error 2"},
	}

	filter.SetBaseline(errors)

	assert.Equal(t, 3, filter.GetBaselineCount())
	assert.Len(t, filter.GetBaselineForFile("main.go"), 2)
	assert.Len(t, filter.GetBaselineForFile("app.go"), 1)
	assert.Nil(t, filter.GetBaselineForFile("unknown.go"))
}

func TestErrorFilter_SetBaselineFromState(t *testing.T) {
	filter := NewErrorFilter()

	state := &BaselineState{
		TypeErrors: []IssueFingerprint{
			{File: "main.go", Line: 10, Message: "type error", Severity: "error"},
		},
		LintIssues: []IssueFingerprint{
			{File: "app.go", Line: 5, Message: "lint warning", Severity: "warning"},
		},
	}

	filter.SetBaselineFromState(state)

	assert.Equal(t, 2, filter.GetBaselineCount())
}

func TestErrorFilter_FilterNewErrors(t *testing.T) {
	filter := NewErrorFilter()

	baseline := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "existing error"},
	}
	filter.SetBaseline(baseline)

	current := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "existing error"}, // not new
		{Type: "type", File: "main.go", Line: 20, Message: "new error"},      // new
		{Type: "lint", File: "app.go", Line: 5, Message: "new warning"},      // new
	}

	newErrors := filter.FilterNewErrors(current)

	assert.Len(t, newErrors, 2)
	assert.Equal(t, "new error", newErrors[0].Message)
	assert.Equal(t, "new warning", newErrors[1].Message)
}

func TestErrorFilter_IsNewError(t *testing.T) {
	filter := NewErrorFilter()

	baseline := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "existing error"},
	}
	filter.SetBaseline(baseline)

	existingError := VerificationIssue{Type: "type", File: "main.go", Line: 10, Message: "existing error"}
	newError := VerificationIssue{Type: "type", File: "main.go", Line: 20, Message: "new error"}

	assert.False(t, filter.IsNewError(existingError))
	assert.True(t, filter.IsNewError(newError))
}

func TestErrorFilter_GetFixedErrors(t *testing.T) {
	filter := NewErrorFilter()

	baseline := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "fixed error"},
		{Type: "type", File: "main.go", Line: 20, Message: "still exists"},
	}
	filter.SetBaseline(baseline)

	current := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 20, Message: "still exists"},
	}

	fixedErrors := filter.GetFixedErrors(current)

	assert.Len(t, fixedErrors, 1)
	assert.Equal(t, "fixed error", fixedErrors[0].Message)
}

func TestErrorFilter_Clear(t *testing.T) {
	filter := NewErrorFilter()

	baseline := []VerificationIssue{
		{Type: "type", File: "main.go", Line: 10, Message: "error"},
	}
	filter.SetBaseline(baseline)
	assert.Equal(t, 1, filter.GetBaselineCount())

	filter.Clear()
	assert.Equal(t, 0, filter.GetBaselineCount())
}


// TestImpactAnalyzer tests
type MockSymbolIndex struct {
	mock.Mock
}

func (m *MockSymbolIndex) IndexProject(ctx context.Context, projectRoot string) error {
	args := m.Called(ctx, projectRoot)
	return args.Error(0)
}

func (m *MockSymbolIndex) IndexFile(ctx context.Context, filePath string, content []byte) error {
	args := m.Called(ctx, filePath, content)
	return args.Error(0)
}

func (m *MockSymbolIndex) SearchByName(query string) []analysis.Symbol {
	args := m.Called(query)
	return args.Get(0).([]analysis.Symbol)
}

func (m *MockSymbolIndex) FindByExactName(name string) []analysis.Symbol {
	args := m.Called(name)
	return args.Get(0).([]analysis.Symbol)
}

func (m *MockSymbolIndex) GetSymbolsInFile(filePath string) []analysis.Symbol {
	args := m.Called(filePath)
	return args.Get(0).([]analysis.Symbol)
}

func (m *MockSymbolIndex) GetSymbolsByKind(kind analysis.SymbolKind) []analysis.Symbol {
	args := m.Called(kind)
	return args.Get(0).([]analysis.Symbol)
}

func (m *MockSymbolIndex) FindDefinition(name string, kind analysis.SymbolKind) *analysis.Symbol {
	args := m.Called(name, kind)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*analysis.Symbol)
}

func (m *MockSymbolIndex) Stats() map[string]int {
	args := m.Called()
	return args.Get(0).(map[string]int)
}

func (m *MockSymbolIndex) IsIndexed() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockSymbolIndex) Clear() {
	m.Called()
}

func setupTestImpactAnalyzer() (*TestImpactAnalyzer, *MockSymbolIndex, *testutils.MockTestService, *testutils.MockLogger) {
	logger := testutils.NewMockLogger()
	logger.On("Info", mock.Anything).Return()
	logger.On("Warning", mock.Anything).Return()
	logger.On("Error", mock.Anything).Return()
	logger.On("Debug", mock.Anything).Return()

	symbolIndex := &MockSymbolIndex{}
	testService := &testutils.MockTestService{}

	analyzer := NewTestImpactAnalyzer(symbolIndex, testService, logger)

	return analyzer, symbolIndex, testService, logger
}

func TestNewTestImpactAnalyzer(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.testIndex)
	assert.False(t, analyzer.IsIndexed())
}

func TestTestImpactAnalyzer_IsTestFile(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	tests := []struct {
		name     string
		file     string
		expected bool
	}{
		{"Go test file", "main_test.go", true},
		{"Go source file", "main.go", false},
		{"JS test file", "app.test.js", true},
		{"JS spec file", "app.spec.js", true},
		{"JS source file", "app.js", false},
		{"TS test file", "component.test.ts", true},
		{"TS spec file", "component.spec.ts", true},
		{"Python test file", "test_module.py", true},
		{"Python source file", "module.py", false},
		{"File in test directory unix", "tests/helper.go", true},
		{"File in __tests__ directory unix", "__tests__/util.ts", true},
		{"File in test directory windows", "tests\\helper.go", true},
		{"File in __tests__ directory windows", "__tests__\\util.ts", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.isTestFile(tt.file)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTestImpactAnalyzer_BuildTestIndex(t *testing.T) {
	analyzer, _, testService, _ := setupTestImpactAnalyzer()

	ctx := context.Background()
	projectRoot := "/test/project"

	testService.On("GetSupportedLanguages").Return([]string{"go"})
	testService.On("DiscoverTests", ctx, projectRoot, "go").Return(&domain.TestSuite{
		Name:     "go-tests",
		Language: "go",
		Tests: []*domain.TestInfo{
			{
				Path:        "main_test.go",
				Name:        "TestMain",
				TargetFiles: []string{"main.go"},
			},
			{
				Path:        "util_test.go",
				Name:        "TestUtil",
				TargetFiles: []string{"util.go", "helper.go"},
			},
		},
	}, nil)

	err := analyzer.BuildTestIndex(ctx, projectRoot)

	assert.NoError(t, err)
	assert.True(t, analyzer.IsIndexed())

	index := analyzer.GetTestIndex()
	assert.Contains(t, index, "main.go")
	assert.Contains(t, index["main.go"], "main_test.go")
	assert.Contains(t, index, "util.go")
	assert.Contains(t, index["util.go"], "util_test.go")

	testService.AssertExpectations(t)
}

func TestTestImpactAnalyzer_GetAffectedTests(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	// Manually set up test index
	analyzer.mu.Lock()
	analyzer.testIndex = map[string][]string{
		"main.go":      {"main_test.go"},
		"util.go":      {"util_test.go", "integration_test.go"},
		"main_test.go": {"main_test.go"},
	}
	analyzer.indexed = true
	analyzer.mu.Unlock()

	tests := []struct {
		name         string
		changedFiles []string
		expected     []string
	}{
		{
			name:         "Single file change",
			changedFiles: []string{"main.go"},
			expected:     []string{"main_test.go"},
		},
		{
			name:         "Multiple file changes",
			changedFiles: []string{"main.go", "util.go"},
			expected:     []string{"main_test.go", "util_test.go", "integration_test.go"},
		},
		{
			name:         "Test file changed",
			changedFiles: []string{"main_test.go"},
			expected:     []string{"main_test.go"},
		},
		{
			name:         "Unknown file",
			changedFiles: []string{"unknown.go"},
			expected:     []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.GetAffectedTests(tt.changedFiles)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestTestImpactAnalyzer_RunTargetedTests(t *testing.T) {
	analyzer, _, testService, _ := setupTestImpactAnalyzer()

	// Set up test index
	analyzer.mu.Lock()
	analyzer.testIndex = map[string][]string{
		"main.go": {"main_test.go"},
	}
	analyzer.indexed = true
	analyzer.mu.Unlock()

	ctx := context.Background()
	projectRoot := "/test/project"
	changedFiles := []string{"main.go"}

	testService.On("RunTargetedTests", ctx, mock.Anything, changedFiles).Return([]*domain.TestResult{
		{Success: true, TestPath: "main_test.go", TestName: "TestMain"},
	}, nil)

	results, err := analyzer.RunTargetedTests(ctx, projectRoot, changedFiles)

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.True(t, results[0].Success)

	testService.AssertExpectations(t)
}

func TestTestImpactAnalyzer_RunTargetedTests_NoAffectedTests(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	// Empty test index
	analyzer.mu.Lock()
	analyzer.testIndex = map[string][]string{}
	analyzer.indexed = true
	analyzer.mu.Unlock()

	ctx := context.Background()
	projectRoot := "/test/project"
	changedFiles := []string{"unknown.go"}

	results, err := analyzer.RunTargetedTests(ctx, projectRoot, changedFiles)

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestTestImpactAnalyzer_Clear(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	// Set up test index
	analyzer.mu.Lock()
	analyzer.testIndex = map[string][]string{
		"main.go": {"main_test.go"},
	}
	analyzer.indexed = true
	analyzer.mu.Unlock()

	assert.True(t, analyzer.IsIndexed())

	analyzer.Clear()

	assert.False(t, analyzer.IsIndexed())
	assert.Empty(t, analyzer.GetTestIndex())
}

func TestTestImpactAnalyzer_FindTestsByConvention(t *testing.T) {
	analyzer, _, _, _ := setupTestImpactAnalyzer()

	// Set up test index with convention-based test files
	analyzer.mu.Lock()
	analyzer.testIndex = map[string][]string{
		"main_test.go":           {"main_test.go"},
		"component.test.ts":      {"component.test.ts"},
		"component.spec.ts":      {"component.spec.ts"},
		"test_module.py":         {"test_module.py"},
		"__tests__/component.ts": {"__tests__/component.ts"},
	}
	analyzer.indexed = true
	analyzer.mu.Unlock()

	tests := []struct {
		name     string
		file     string
		expected []string
	}{
		{
			name:     "Go file convention",
			file:     "main.go",
			expected: []string{"main_test.go"},
		},
		{
			name:     "TypeScript file convention",
			file:     "component.ts",
			expected: []string{"component.test.ts", "component.spec.ts"},
		},
		{
			name:     "Python file convention",
			file:     "module.py",
			expected: []string{"test_module.py"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.findTestsByConvention(tt.file)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
