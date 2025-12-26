package staticanalyzer

import (
	"context"
	"syntaxia/domain"
	"testing"
)

func TestRuffAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"python", "py"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestRuffAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeRuff {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeRuff)
	}
}

func TestRuffAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid python config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: ".",
			},
			wantErr: true, // Will fail because no Python files exist
			errMsg:  "no Python files found",
		},
		{
			name: "valid py config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "py",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no Python files found",
		},
		{
			name: "invalid language - go",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "Ruff analyzer only supports Python language",
		},
		{
			name: "invalid language - rust",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "Ruff analyzer only supports Python language",
		},
		{
			name: "invalid language - java",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "Ruff analyzer only supports Python language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := analyzer.ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateConfig() error = %q, want to contain %q", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestRuffAnalyzer_ParseRuffOutput(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         string
		expectedCount  int
		expectedIssues []expectedIssue
	}{
		{
			name:          "empty output",
			output:        "",
			expectedCount: 0,
		},
		{
			name:          "non-json output",
			output:        "All checks passed!",
			expectedCount: 0,
		},
		{
			name:          "single error",
			output:        `{"code": "E501", "message": "Line too long (120 > 88 characters)", "filename": "src/main.py", "row": 10, "column": 89}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/main.py", line: 10, severity: "error", code: "E501"},
			},
		},
		{
			name:          "single warning",
			output:        `{"code": "W291", "message": "Trailing whitespace", "filename": "src/utils.py", "row": 5, "column": 20}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/utils.py", line: 5, severity: "warning", code: "W291"},
			},
		},
		{
			name: "multiple issues",
			output: `{"code": "F401", "message": "Module imported but unused", "filename": "src/app.py", "row": 1, "column": 1}
{"code": "E302", "message": "Expected 2 blank lines, found 1", "filename": "src/app.py", "row": 10, "column": 1}
{"code": "W503", "message": "Line break before binary operator", "filename": "src/app.py", "row": 15, "column": 5}`,
			expectedCount: 3,
			expectedIssues: []expectedIssue{
				{file: "src/app.py", line: 1, severity: "warning", code: "F401"},
				{file: "src/app.py", line: 10, severity: "error", code: "E302"},
				{file: "src/app.py", line: 15, severity: "warning", code: "W503"},
			},
		},
		{
			name:          "import issue",
			output:        `{"code": "I001", "message": "Import block is un-sorted or un-formatted", "filename": "src/imports.py", "row": 1, "column": 1}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/imports.py", line: 1, severity: "warning", code: "I001"},
			},
		},
		{
			name:          "naming convention issue",
			output:        `{"code": "N802", "message": "Function name should be lowercase", "filename": "src/naming.py", "row": 5, "column": 5}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/naming.py", line: 5, severity: "warning", code: "N802"},
			},
		},
		{
			name:          "pyupgrade issue",
			output:        `{"code": "UP006", "message": "Use `+"`list`"+` instead of `+"`List`"+` for type annotation", "filename": "src/types.py", "row": 3, "column": 10}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/types.py", line: 3, severity: "warning", code: "UP006"},
			},
		},
		{
			name:          "flake8 issue",
			output:        `{"code": "F841", "message": "Local variable is assigned but never used", "filename": "src/unused.py", "row": 10, "column": 5}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/unused.py", line: 10, severity: "warning", code: "F841"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseRuffOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseRuffOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseRuffOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
				return
			}

			for i, expected := range tt.expectedIssues {
				if i >= len(issues) {
					break
				}
				issue := issues[i]
				if issue.File != expected.file {
					t.Errorf("issue[%d].File = %q, want %q", i, issue.File, expected.file)
				}
				if issue.Line != expected.line {
					t.Errorf("issue[%d].Line = %d, want %d", i, issue.Line, expected.line)
				}
				if issue.Severity != expected.severity {
					t.Errorf("issue[%d].Severity = %q, want %q", i, issue.Severity, expected.severity)
				}
				if issue.Code != expected.code {
					t.Errorf("issue[%d].Code = %q, want %q", i, issue.Code, expected.code)
				}
			}
		})
	}
}

func TestRuffAnalyzer_ParseRuffLine(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		expected *domain.StaticIssue
	}{
		{
			name:     "empty line",
			line:     "",
			expected: nil,
		},
		{
			name:     "non-json line",
			line:     "Some random text",
			expected: nil,
		},
		{
			name: "valid error line",
			line: `{"code": "E501", "message": "Line too long", "filename": "test.py", "row": 10, "column": 89}`,
			expected: &domain.StaticIssue{
				File:     "test.py",
				Line:     10,
				Column:   89,
				Severity: "error",
				Code:     "E501",
				Message:  "Line too long",
				Category: "error",
			},
		},
		{
			name: "valid warning line",
			line: `{"code": "W291", "message": "Trailing whitespace", "filename": "utils.py", "row": 5, "column": 20}`,
			expected: &domain.StaticIssue{
				File:     "utils.py",
				Line:     5,
				Column:   20,
				Severity: "warning",
				Code:     "W291",
				Message:  "Trailing whitespace",
				Category: "warning",
			},
		},
		{
			name:     "line without code field",
			line:     `{"message": "Some message", "filename": "test.py"}`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.parseRuffLine(tt.line)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseRuffLine(%q) = %+v, want nil", tt.line, result)
				}
				return
			}

			if result == nil {
				t.Errorf("parseRuffLine(%q) = nil, want %+v", tt.line, tt.expected)
				return
			}

			if result.File != tt.expected.File {
				t.Errorf("File = %q, want %q", result.File, tt.expected.File)
			}
			if result.Line != tt.expected.Line {
				t.Errorf("Line = %d, want %d", result.Line, tt.expected.Line)
			}
			if result.Column != tt.expected.Column {
				t.Errorf("Column = %d, want %d", result.Column, tt.expected.Column)
			}
			if result.Severity != tt.expected.Severity {
				t.Errorf("Severity = %q, want %q", result.Severity, tt.expected.Severity)
			}
			if result.Code != tt.expected.Code {
				t.Errorf("Code = %q, want %q", result.Code, tt.expected.Code)
			}
		})
	}
}

func TestRuffAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewRuffAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"E501", "error"},
		{"E302", "error"},
		{"W291", "warning"},
		{"W503", "warning"},
		{"F401", "flake8"},
		{"F841", "flake8"},
		{"I001", "imports"},
		{"I002", "imports"},
		{"N802", "naming"},
		{"N806", "naming"},
		{"UP006", "pyupgrade"},
		{"UP035", "pyupgrade"},
		{"X001", "other"},
		{"", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := analyzer.getCategory(tt.code)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestRuffAnalyzer_Analyze_RuffNotInstalled(t *testing.T) {
	// This test verifies behavior when ruff is not installed
	analyzer := NewRuffAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:    "python",
		ProjectPath: "/non/existent/path",
	}

	// The analyze method should return a result with error message if ruff is not installed
	result, err := analyzer.Analyze(ctx, config)
	if err != nil {
		t.Errorf("Analyze() should not return error, got %v", err)
		return
	}

	// Result should be returned even if ruff is not installed
	if result == nil {
		t.Error("Analyze() should return a result even when ruff is not installed")
		return
	}

	// If ruff is not installed, there should be an error message in the result
	// If ruff is installed, the result might have different content
	t.Logf("Analyze result: Success=%v, Error=%q", result.Success, result.Error)
}

func TestExtractJSONString(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		prefix   string
		expected string
	}{
		{
			name:     "extract code",
			line:     `{"code": "E501", "message": "test"}`,
			prefix:   `"code": "`,
			expected: "E501",
		},
		{
			name:     "extract message",
			line:     `{"code": "E501", "message": "Line too long"}`,
			prefix:   `"message": "`,
			expected: "Line too long",
		},
		{
			name:     "extract filename",
			line:     `{"filename": "src/main.py", "row": 10}`,
			prefix:   `"filename": "`,
			expected: "src/main.py",
		},
		{
			name:     "prefix not found",
			line:     `{"code": "E501"}`,
			prefix:   `"missing": "`,
			expected: "",
		},
		{
			name:     "empty line",
			line:     "",
			prefix:   `"code": "`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSONString(tt.line, tt.prefix)
			if result != tt.expected {
				t.Errorf("extractJSONString(%q, %q) = %q, want %q", tt.line, tt.prefix, result, tt.expected)
			}
		})
	}
}

func TestExtractJSONInt(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		prefix   string
		expected int
	}{
		{
			name:     "extract row",
			line:     `{"row": 10, "column": 5}`,
			prefix:   `"row": `,
			expected: 10,
		},
		{
			name:     "extract column",
			line:     `{"row": 10, "column": 5}`,
			prefix:   `"column": `,
			expected: 5,
		},
		{
			name:     "extract at end of object",
			line:     `{"row": 100}`,
			prefix:   `"row": `,
			expected: 100,
		},
		{
			name:     "prefix not found",
			line:     `{"row": 10}`,
			prefix:   `"missing": `,
			expected: 0,
		},
		{
			name:     "empty line",
			line:     "",
			prefix:   `"row": `,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractJSONInt(tt.line, tt.prefix)
			if result != tt.expected {
				t.Errorf("extractJSONInt(%q, %q) = %d, want %d", tt.line, tt.prefix, result, tt.expected)
			}
		})
	}
}
