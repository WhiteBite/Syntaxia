package repair

import (
	"context"
	"io/fs"
	"syntaxia/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testMainGoContent = "package main\n\nfunc main() {}"

// TestLogger is a mock logger for testing
type TestLogger struct{}

func (l *TestLogger) Debug(msg string)                          {}
func (l *TestLogger) Info(msg string)                           {}
func (l *TestLogger) Warning(msg string)                        {}
func (l *TestLogger) Error(msg string)                          {}
func (l *TestLogger) Fatal(msg string)                          {}
func (l *TestLogger) Debugf(format string, args ...interface{}) {}
func (l *TestLogger) Infof(format string, args ...interface{})  {}
func (l *TestLogger) Warnf(format string, args ...interface{})  {}
func (l *TestLogger) Errorf(format string, args ...interface{}) {}
func (l *TestLogger) Fatalf(format string, args ...interface{}) {}

// MockFileSystemProvider is a mock file system provider for testing
type MockFileSystemProvider struct {
	ReadFileContent string
	WriteFileError  error
}

func (m *MockFileSystemProvider) ReadFile(path string) ([]byte, error) {
	return []byte(m.ReadFileContent), nil
}

func (m *MockFileSystemProvider) WriteFile(path string, data []byte, perm int) error {
	return m.WriteFileError
}

func (m *MockFileSystemProvider) Stat(path string) (fs.FileInfo, error) {
	return &mockFileInfo{}, nil
}

func (m *MockFileSystemProvider) ReadDir(path string) ([]fs.DirEntry, error) {
	return nil, nil
}

func (m *MockFileSystemProvider) MkdirAll(path string, perm int) error {
	return nil
}

func (m *MockFileSystemProvider) Remove(path string) error {
	return nil
}

func (m *MockFileSystemProvider) RemoveAll(path string) error {
	return nil
}

func (m *MockFileSystemProvider) Rename(oldpath, newpath string) error {
	return nil
}

func (m *MockFileSystemProvider) Exists(path string) bool {
	return true
}

func (m *MockFileSystemProvider) IsDir(path string) bool {
	return false
}

func (m *MockFileSystemProvider) GetAbsolutePath(path string) (string, error) {
	return path, nil
}

type mockFileInfo struct{}

func (m *mockFileInfo) Name() string       { return "test" }
func (m *mockFileInfo) Size() int64        { return 0 }
func (m *mockFileInfo) Mode() fs.FileMode  { return 0 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() interface{}   { return nil }

// TestErrorAnalyzer_AnalyzeError tests error analysis functionality
func TestErrorAnalyzer_AnalyzeError(t *testing.T) {
	tests := []struct {
		name        string
		errorOutput string
		stage       domain.ProtocolStage
		expected    func(*domain.ErrorDetails) bool
	}{
		{
			name:        "go_compilation_error",
			errorOutput: "main.go:10:5: undefined: fmt",
			stage:       domain.StageBuilding,
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeCompilation &&
					details.SourceFile == "main.go" &&
					details.LineNumber == 10 &&
					details.Column == 5
			},
		},
		{
			name:        "typescript_type_error",
			errorOutput: "src/main.ts:25:10: error TS2304: Cannot find name 'SomeType'",
			stage:       domain.StageBuilding,
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeImport &&
					details.SourceFile == "src/main.ts" &&
					details.LineNumber == 25
			},
		},
		{
			name:        "linting_error",
			errorOutput: "style violation: missing semicolon",
			stage:       domain.StageLinting,
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeLinting &&
					details.Stage == domain.StageLinting
			},
		},
		{
			name:        "test_failure",
			errorOutput: "test failed: expected 5, got 3",
			stage:       domain.StageTesting,
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeTesting &&
					details.Stage == domain.StageTesting
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			analyzer := NewErrorAnalyzer(logger)

			result, err := analyzer.AnalyzeError(tt.errorOutput, tt.stage)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, tt.expected(result), "Error analysis result validation failed")
		})
	}
}

// TestErrorAnalyzer_SuggestCorrections tests correction suggestion functionality
func TestErrorAnalyzer_SuggestCorrections(t *testing.T) {
	tests := []struct {
		name     string
		error    *domain.ErrorDetails
		expected func([]*domain.CorrectionStep) bool
	}{
		{
			name: "import_error_corrections",
			error: &domain.ErrorDetails{
				Stage:      domain.StageBuilding,
				ErrorType:  domain.ErrorTypeImport,
				Message:    "undefined: fmt",
				SourceFile: "main.go",
			},
			expected: func(steps []*domain.CorrectionStep) bool {
				if len(steps) == 0 {
					return false
				}
				for _, step := range steps {
					if step.Action == domain.ActionFixImport {
						return true
					}
				}
				return false
			},
		},
		{
			name: "syntax_error_corrections",
			error: &domain.ErrorDetails{
				Stage:      domain.StageBuilding,
				ErrorType:  domain.ErrorTypeSyntax,
				Message:    "syntax error: missing )",
				SourceFile: "main.go",
			},
			expected: func(steps []*domain.CorrectionStep) bool {
				return len(steps) > 0 &&
					steps[0].Action == domain.ActionFixSyntax
			},
		},
		{
			name: "type_error_corrections",
			error: &domain.ErrorDetails{
				Stage:      domain.StageBuilding,
				ErrorType:  domain.ErrorTypeTypeCheck,
				Message:    "cannot use string as int",
				SourceFile: "main.go",
			},
			expected: func(steps []*domain.CorrectionStep) bool {
				return len(steps) > 0 &&
					steps[0].Action == domain.ActionFixType
			},
		},
		{
			name: "linting_error_corrections",
			error: &domain.ErrorDetails{
				Stage:      domain.StageLinting,
				ErrorType:  domain.ErrorTypeLinting,
				Message:    "formatting issue",
				SourceFile: "main.go",
			},
			expected: func(steps []*domain.CorrectionStep) bool {
				return len(steps) > 0 &&
					steps[0].Action == domain.ActionFormatCode
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			analyzer := NewErrorAnalyzer(logger)

			result, err := analyzer.SuggestCorrections(tt.error)

			assert.NoError(t, err)
			assert.True(t, tt.expected(result), "Correction suggestions validation failed")
		})
	}
}

// TestErrorAnalyzer_ClassifyErrorType tests error type classification
func TestErrorAnalyzer_ClassifyErrorType(t *testing.T) {
	tests := []struct {
		name        string
		errorOutput string
		expected    domain.ErrorType
	}{
		{
			name:        "compilation_error",
			errorOutput: "compile error: syntax error",
			expected:    domain.ErrorTypeCompilation,
		},
		{
			name:        "type_check_error",
			errorOutput: "type error: cannot use string as int",
			expected:    domain.ErrorTypeTypeCheck,
		},
		{
			name:        "import_error",
			errorOutput: "import error: module not found",
			expected:    domain.ErrorTypeImport,
		},
		{
			name:        "linting_error",
			errorOutput: "lint error: style violation",
			expected:    domain.ErrorTypeLinting,
		},
		{
			name:        "test_error",
			errorOutput: "test failed: assertion error",
			expected:    domain.ErrorTypeTesting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			analyzer := NewErrorAnalyzer(logger)

			result := analyzer.ClassifyErrorType(tt.errorOutput)

			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCorrectionEngine_ApplyCorrection tests single correction application
func TestCorrectionEngine_ApplyCorrection(t *testing.T) {
	tests := []struct {
		name      string
		step      *domain.CorrectionStep
		mockSetup func(*MockFileSystemProvider)
		expected  func(*domain.CorrectionResult) bool
	}{
		{
			name: "format_code_correction",
			step: &domain.CorrectionStep{
				Action:      domain.ActionFormatCode,
				Target:      "main.go",
				Description: "Format Go file",
			},
			mockSetup: func(mock *MockFileSystemProvider) {},
			expected: func(result *domain.CorrectionResult) bool {
				return result.Success && result.Message != ""
			},
		},
		{
			name: "import_fix_correction",
			step: &domain.CorrectionStep{
				Action:      domain.ActionFixImport,
				Target:      "main.go",
				Description: "Fix import statement",
			},
			mockSetup: func(mock *MockFileSystemProvider) {
				mock.ReadFileContent = testMainGoContent
			},
			expected: func(result *domain.CorrectionResult) bool {
				return result != nil && result.Message != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			fileSystemProvider := &MockFileSystemProvider{}
			tt.mockSetup(fileSystemProvider)

			engine := NewCorrectionEngine(logger, fileSystemProvider)

			result, err := engine.ApplyCorrection(context.Background(), tt.step, "/test/project")

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, tt.expected(result), "Correction result validation failed")
		})
	}
}

// TestCorrectionEngine_ApplyCorrections tests multiple correction application
func TestCorrectionEngine_ApplyCorrections(t *testing.T) {
	tests := []struct {
		name      string
		steps     []*domain.CorrectionStep
		mockSetup func(*MockFileSystemProvider)
		expected  func(*domain.CorrectionResult) bool
	}{
		{
			name: "multiple_corrections_success",
			steps: []*domain.CorrectionStep{
				{
					Action:      domain.ActionFormatCode,
					Target:      "main.go",
					Description: "Format code",
				},
				{
					Action:      domain.ActionFixSyntax,
					Target:      "main.go",
					Description: "Fix syntax",
				},
			},
			mockSetup: func(mock *MockFileSystemProvider) {
				mock.ReadFileContent = testMainGoContent
			},
			expected: func(result *domain.CorrectionResult) bool {
				return result.Success && result.Message != ""
			},
		},
		{
			name: "mixed_success_failure",
			steps: []*domain.CorrectionStep{
				{
					Action:      domain.ActionFixImport,
					Target:      "main.go",
					Description: "Fix import",
				},
				{
					Action:      domain.CorrectionAction("unsupported"),
					Target:      "main.go",
					Description: "Unsupported action",
				},
			},
			mockSetup: func(mock *MockFileSystemProvider) {
				mock.ReadFileContent = testMainGoContent
			},
			expected: func(result *domain.CorrectionResult) bool {
				return !result.Success
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			fileSystemProvider := &MockFileSystemProvider{}
			tt.mockSetup(fileSystemProvider)

			engine := NewCorrectionEngine(logger, fileSystemProvider)

			result, err := engine.ApplyCorrections(context.Background(), tt.steps, "/test/project")

			if tt.expected(result) {
				assert.NoError(t, err)
			}
			assert.NotNil(t, result)
			assert.True(t, tt.expected(result), "Multiple corrections result validation failed")
		})
	}
}

// TestCorrectionEngine_CanHandle tests error handling capability check
func TestCorrectionEngine_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		error    *domain.ErrorDetails
		expected bool
	}{
		{
			name: "import_error_can_handle",
			error: &domain.ErrorDetails{
				ErrorType: domain.ErrorTypeImport,
			},
			expected: true,
		},
		{
			name: "syntax_error_can_handle",
			error: &domain.ErrorDetails{
				ErrorType: domain.ErrorTypeSyntax,
			},
			expected: true,
		},
		{
			name: "type_error_can_handle",
			error: &domain.ErrorDetails{
				ErrorType: domain.ErrorTypeTypeCheck,
			},
			expected: true,
		},
		{
			name: "linting_error_can_handle",
			error: &domain.ErrorDetails{
				ErrorType: domain.ErrorTypeLinting,
			},
			expected: true,
		},
		{
			name: "unsupported_error_cannot_handle",
			error: &domain.ErrorDetails{
				ErrorType: domain.ErrorType("unsupported"),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := &TestLogger{}
			fileSystemProvider := &MockFileSystemProvider{}
			engine := NewCorrectionEngine(logger, fileSystemProvider)

			result := engine.CanHandle(tt.error)

			assert.Equal(t, tt.expected, result)
		})
	}
}

// Language-specific analyzer tests

func TestGoErrorAnalyzer(t *testing.T) {
	analyzer := NewGoErrorAnalyzer()

	tests := []struct {
		name        string
		errorOutput string
		expected    func(*domain.ErrorDetails) bool
	}{
		{
			name:        "go_undefined_error",
			errorOutput: "undefined: fmt",
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeCompilation &&
					details.Tool == "go" &&
					len(details.Suggestions) > 0
			},
		},
		{
			name:        "go_type_error",
			errorOutput: "cannot use string as int in assignment",
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeTypeCheck &&
					details.Tool == "go" &&
					len(details.Suggestions) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.AnalyzeError(tt.errorOutput)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, tt.expected(result), "Go error analysis validation failed")
			assert.Equal(t, "go", analyzer.GetLanguage())
		})
	}
}

func TestTypeScriptErrorAnalyzer(t *testing.T) {
	analyzer := NewTypeScriptErrorAnalyzer()

	tests := []struct {
		name        string
		errorOutput string
		expected    func(*domain.ErrorDetails) bool
	}{
		{
			name:        "typescript_cannot_find_name",
			errorOutput: "error TS2304: Cannot find name 'SomeType'",
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeImport &&
					details.Tool == "tsc" &&
					len(details.Suggestions) > 0
			},
		},
		{
			name:        "typescript_type_assignment",
			errorOutput: "error TS2322: Type 'string' is not assignable to type 'number'",
			expected: func(details *domain.ErrorDetails) bool {
				return details.ErrorType == domain.ErrorTypeTypeCheck &&
					details.Tool == "tsc" &&
					len(details.Suggestions) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := analyzer.AnalyzeError(tt.errorOutput)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.True(t, tt.expected(result), "TypeScript error analysis validation failed")
			assert.Equal(t, "typescript", analyzer.GetLanguage())
		})
	}
}

// Integration tests

func TestErrorAnalyzerAndCorrectionEngine_Integration(t *testing.T) {
	logger := &TestLogger{}
	analyzer := NewErrorAnalyzer(logger)
	fileSystemProvider := &MockFileSystemProvider{}
	fileSystemProvider.ReadFileContent = "package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}"
	engine := NewCorrectionEngine(logger, fileSystemProvider)

	errorOutput := "main.go:4:2: undefined: fmt"
	stage := domain.StageBuilding

	errorDetails, err := analyzer.AnalyzeError(errorOutput, stage)
	assert.NoError(t, err)
	assert.NotNil(t, errorDetails)

	corrections, err := analyzer.SuggestCorrections(errorDetails)
	assert.NoError(t, err)
	assert.True(t, len(corrections) > 0)

	canHandle := engine.CanHandle(errorDetails)
	assert.True(t, canHandle)

	result, err := engine.ApplyCorrections(context.Background(), corrections, "/test/project")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	t.Logf("Integration test completed: error analyzed, corrections suggested, and applied")
}
