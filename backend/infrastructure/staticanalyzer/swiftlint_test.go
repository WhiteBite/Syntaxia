package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestSwiftLintAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"swift"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestSwiftLintAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeSwiftLint {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeSwiftLint)
	}
}

func TestSwiftLintAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "swift",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "non-existent project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "swift",
				ProjectPath: "/non/existent/path/to/project",
			},
			wantErr: true,
			errMsg:  "project path does not exist",
		},
		{
			name: "valid config with existing path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "swift",
				ProjectPath: ".",
			},
			wantErr: false,
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

func TestSwiftLintAnalyzer_ParseJSONOutput(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         string
		projectPath    string
		expectedCount  int
		expectedIssues []expectedIssue
		wantErr        bool
	}{
		{
			name:          "empty array",
			output:        "[]",
			projectPath:   "/project",
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:        "no JSON found",
			output:      "some random text without json",
			projectPath: "/project",
			wantErr:     true,
		},
		{
			name:          "single warning",
			output:        `[{"rule_id":"trailing_whitespace","reason":"Lines should not have trailing whitespace.","line":10,"column":5,"file":"main.swift","severity":"warning"}]`,
			projectPath:   "/project",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "main.swift", line: 10, severity: "warning", code: "trailing_whitespace"},
			},
		},
		{
			name:          "single error",
			output:        `[{"rule_id":"force_cast","reason":"Force casts should be avoided.","line":25,"column":10,"file":"Utils.swift","severity":"error"}]`,
			projectPath:   "/project",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "Utils.swift", line: 25, severity: "error", code: "force_cast"},
			},
		},
		{
			name: "multiple issues",
			output: `[
				{"rule_id":"line_length","reason":"Line should be 120 characters or less","line":5,"column":121,"file":"File1.swift","severity":"warning"},
				{"rule_id":"force_unwrapping","reason":"Force unwrapping should be avoided.","line":15,"column":8,"file":"File2.swift","severity":"error"}
			]`,
			projectPath:   "/project",
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "File1.swift", line: 5, severity: "warning", code: "line_length"},
				{file: "File2.swift", line: 15, severity: "error", code: "force_unwrapping"},
			},
		},
		{
			name:          "JSON with prefix text",
			output:        `Loading configuration from .swiftlint.yml\n[{"rule_id":"test","reason":"Test","line":1,"file":"test.swift","severity":"warning"}]`,
			projectPath:   "/project",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseJSONOutput(tt.output, tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSONOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseJSONOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestSwiftLintAnalyzer_ParsePlainOutput(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         string
		projectPath    string
		expectedCount  int
		expectedIssues []expectedIssue
	}{
		{
			name:          "empty output",
			output:        "",
			projectPath:   ".",
			expectedCount: 0,
		},
		{
			name:          "single warning",
			output:        "main.swift:10:5: warning: Lines should not have trailing whitespace. (trailing_whitespace)",
			projectPath:   ".",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "main.swift", line: 10, severity: "warning", code: "trailing_whitespace"},
			},
		},
		{
			name:          "single error",
			output:        "Utils.swift:25:10: error: Force casts should be avoided. (force_cast)",
			projectPath:   ".",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "Utils.swift", line: 25, severity: "error", code: "force_cast"},
			},
		},
		{
			name: "multiple issues",
			output: `File1.swift:5:121: warning: Line should be 120 characters or less (line_length)
File2.swift:15:8: error: Force unwrapping should be avoided. (force_unwrapping)`,
			projectPath:   ".",
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "File1.swift", line: 5, severity: "warning", code: "line_length"},
				{file: "File2.swift", line: 15, severity: "error", code: "force_unwrapping"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := analyzer.parsePlainOutput(tt.output, tt.projectPath)

			if len(issues) != tt.expectedCount {
				t.Errorf("parsePlainOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestSwiftLintAnalyzer_MapSeverity(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"Error", "error"},
		{"warning", "warning"},
		{"WARNING", "warning"},
		{"Warning", "warning"},
		{"info", "info"},
		{"unknown", "info"},
		{"", "info"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := analyzer.mapSeverity(tt.input)
			if result != tt.expected {
				t.Errorf("mapSeverity(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSwiftLintAnalyzer_FindConfigFile(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	// Test with non-existent path - should return empty string
	result := analyzer.findConfigFile("/non/existent/path")
	if result != "" {
		t.Errorf("findConfigFile() for non-existent path = %q, want empty string", result)
	}

	// Test with current directory (no swiftlint config)
	result = analyzer.findConfigFile(".")
	if result != "" {
		t.Errorf("findConfigFile() for current dir = %q, want empty string", result)
	}
}

func TestSwiftLintAnalyzer_HasSwiftFiles(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	// Test with current directory (no swift files)
	hasFiles, err := analyzer.hasSwiftFiles(".")
	if err != nil {
		t.Errorf("hasSwiftFiles() error = %v", err)
	}
	if hasFiles {
		t.Errorf("hasSwiftFiles() for current dir = true, want false")
	}

	// Test with non-existent path
	hasFiles, err = analyzer.hasSwiftFiles("/non/existent/path")
	if err != nil {
		t.Errorf("hasSwiftFiles() error = %v", err)
	}
	if hasFiles {
		t.Errorf("hasSwiftFiles() for non-existent path = true, want false")
	}
}

func TestSwiftLintAnalyzer_ParseJSONOutput_EdgeCases(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		name        string
		output      string
		projectPath string
		wantErr     bool
	}{
		{
			name:        "malformed JSON",
			output:      `[{"rule_id": "test", "reason": "incomplete`,
			projectPath: "/project",
			wantErr:     true,
		},
		{
			name:        "JSON with trailing text",
			output:      `[{"rule_id":"test","reason":"Test","line":1,"file":"test.swift","severity":"warning"}]\nDone.`,
			projectPath: "/project",
			wantErr:     false,
		},
		{
			name:        "only closing bracket",
			output:      `]`,
			projectPath: "/project",
			wantErr:     true,
		},
		{
			name:        "brackets in wrong order",
			output:      `][`,
			projectPath: "/project",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := analyzer.parseJSONOutput(tt.output, tt.projectPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseJSONOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSwiftLintAnalyzer_ParsePlainOutput_EdgeCases(t *testing.T) {
	analyzer := NewSwiftLintAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        string
		expectedCount int
	}{
		{
			name:          "whitespace only",
			output:        "   \n\t\n   ",
			expectedCount: 0,
		},
		{
			name:          "incomplete line format",
			output:        "file.swift:10",
			expectedCount: 0,
		},
		{
			name:          "line without severity prefix",
			output:        "file.swift:10:5: Some message without severity",
			expectedCount: 1,
		},
		{
			name: "mixed valid and invalid lines",
			output: `file1.swift:10:5: warning: Valid warning (rule1)
invalid line
file2.swift:20:1: error: Valid error (rule2)`,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := analyzer.parsePlainOutput(tt.output, ".")
			if len(issues) != tt.expectedCount {
				t.Errorf("parsePlainOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}
