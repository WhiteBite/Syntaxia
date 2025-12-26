package staticanalyzer

import (
	"context"
	"syntaxia/domain"
	"testing"
)

func TestErrorProneAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"java"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestErrorProneAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeErrorProne {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeErrorProne)
	}
}

func TestErrorProneAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid java config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
				ProjectPath: ".",
			},
			wantErr: true, // Will fail because no Java files exist
			errMsg:  "no Java files found",
		},
		{
			name: "invalid language - go",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ErrorProne analyzer only supports Java language",
		},
		{
			name: "invalid language - python",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ErrorProne analyzer only supports Java language",
		},
		{
			name: "invalid language - rust",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ErrorProne analyzer only supports Java language",
		},
		{
			name: "invalid language - typescript",
			config: &domain.StaticAnalyzerConfig{
				Language:    "typescript",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ErrorProne analyzer only supports Java language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
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

func TestErrorProneAnalyzer_ParseErrorProneOutput(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})

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
			name:          "non-error output",
			output:        "Compiling 10 source files...\nBuild successful",
			expectedCount: 0,
		},
		{
			name:          "single error",
			output:        "src/Main.java:10:15: [error] NullPointerException may be thrown",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Main.java", line: 10, severity: "error"},
			},
		},
		{
			name:          "single warning",
			output:        "src/Utils.java:25:5: [warning] Unused variable 'temp'",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Utils.java", line: 25, severity: "warning"},
			},
		},
		{
			name: "multiple issues",
			output: `src/Main.java:10:15: [error] NullPointerException may be thrown
src/Main.java:20:10: [warning] Deprecated method usage
src/Utils.java:5:1: [error] Resource leak: stream is never closed`,
			expectedCount: 3,
			expectedIssues: []expectedIssue{
				{file: "src/Main.java", line: 10, severity: "error"},
				{file: "src/Main.java", line: 20, severity: "warning"},
				{file: "src/Utils.java", line: 5, severity: "error"},
			},
		},
		{
			name:          "null safety issue",
			output:        "src/Service.java:50:20: [error] null dereference",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Service.java", line: 50, severity: "error"},
			},
		},
		{
			name:          "unused code issue",
			output:        "src/Helper.java:30:5: [warning] unused private method 'helper'",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Helper.java", line: 30, severity: "warning"},
			},
		},
		{
			name:          "deprecation issue",
			output:        "src/Legacy.java:15:10: [warning] deprecated API usage",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Legacy.java", line: 15, severity: "warning"},
			},
		},
		{
			name:          "concurrency issue",
			output:        "src/ThreadPool.java:100:5: [error] concurrent modification detected",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/ThreadPool.java", line: 100, severity: "error"},
			},
		},
		{
			name: "mixed with non-error lines",
			output: `Note: Some input files use unchecked operations
src/Main.java:10:15: [error] NullPointerException may be thrown
Note: Recompile with -Xlint:unchecked for details
src/Utils.java:5:1: [warning] Unused import`,
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "src/Main.java", line: 10, severity: "error"},
				{file: "src/Utils.java", line: 5, severity: "warning"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseErrorProneOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseErrorProneOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseErrorProneOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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
			}
		})
	}
}

func TestErrorProneAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})

	tests := []struct {
		message  string
		expected string
	}{
		{"NullPointerException may be thrown", "null-safety"},
		{"null dereference detected", "null-safety"},
		{"Possible null value", "null-safety"},
		{"Unused variable 'temp'", "unused-code"},
		{"unused private method", "unused-code"},
		{"Deprecated method usage", "deprecation"},
		{"deprecated API call", "deprecation"},
		{"Concurrent modification", "concurrency"},
		{"concurrent access issue", "concurrency"},
		{"Resource leak detected", "resource-management"},
		{"resource not closed", "resource-management"},
		{"Some other error message", "other"},
		{"", "other"},
		{"Generic compilation error", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := analyzer.getCategory(tt.message)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.message, result, tt.expected)
			}
		})
	}
}

func TestErrorProneAnalyzer_Analyze_JavaNotInstalled(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:    "java",
		ProjectPath: "/non/existent/path",
	}

	result, err := analyzer.Analyze(ctx, config)
	if err != nil {
		t.Errorf("Analyze() should not return error, got %v", err)
		return
	}

	if result == nil {
		t.Error("Analyze() should return a result even when Java is not installed")
		return
	}

	t.Logf("Analyze result: Success=%v, Error=%q", result.Success, result.Error)
}

func TestErrorProneAnalyzer_ParseErrorProneOutput_EdgeCases(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        string
		expectedCount int
	}{
		{
			name:          "line with insufficient parts",
			output:        "src/Main.java:10",
			expectedCount: 0,
		},
		{
			name:          "line with non-numeric line number",
			output:        "src/Main.java:abc:10: [error] Some error",
			expectedCount: 0,
		},
		{
			name:          "line with non-numeric column",
			output:        "src/Main.java:10:abc: [error] Some error",
			expectedCount: 0,
		},
		{
			name:          "only whitespace",
			output:        "   \n\t\n   ",
			expectedCount: 0,
		},
		{
			name:          "windows path",
			output:        `C:\project\src\Main.java:10:15: [error] Some error`,
			expectedCount: 0, // Windows paths with drive letter won't parse correctly
		},
		{
			name:          "relative path with dots",
			output:        "../src/Main.java:10:15: [error] Some error",
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseErrorProneOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseErrorProneOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseErrorProneOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestErrorProneAnalyzer_ParseErrorProneOutput_MessageExtraction(t *testing.T) {
	analyzer := NewErrorProneAnalyzer(&mockLogger{})

	tests := []struct {
		name            string
		output          string
		expectedMessage string
	}{
		{
			name:            "error prefix stripped",
			output:          "src/Main.java:10:15: [error] NullPointerException",
			expectedMessage: "NullPointerException",
		},
		{
			name:            "warning prefix stripped",
			output:          "src/Main.java:10:15: [warning] Unused variable",
			expectedMessage: "Unused variable",
		},
		{
			name:            "no prefix",
			output:          "src/Main.java:10:15: Some message without prefix",
			expectedMessage: " Some message without prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseErrorProneOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseErrorProneOutput() error = %v", err)
				return
			}

			if len(issues) != 1 {
				t.Errorf("parseErrorProneOutput() returned %d issues, want 1", len(issues))
				return
			}

			if issues[0].Message != tt.expectedMessage {
				t.Errorf("Message = %q, want %q", issues[0].Message, tt.expectedMessage)
			}
		})
	}
}
