package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestESLintAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"typescript", "ts", "javascript", "js"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestESLintAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeESLint {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeESLint)
	}
}

func TestESLintAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid language - go",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ESLint analyzer only supports TypeScript/JavaScript",
		},
		{
			name: "invalid language - rust",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ESLint analyzer only supports TypeScript/JavaScript",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "typescript",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid typescript config - no ts files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "typescript",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no TypeScript/JavaScript files found",
		},
		{
			name: "valid ts config - no ts files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ts",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no TypeScript/JavaScript files found",
		},
		{
			name: "valid javascript config - no js files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "javascript",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no TypeScript/JavaScript files found",
		},
		{
			name: "valid js config - no js files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "js",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no TypeScript/JavaScript files found",
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

func TestESLintAnalyzer_ParseESLintOutput(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         []byte
		expectedCount  int
		expectedIssues []expectedIssue
		wantErr        bool
	}{
		{
			name:          "empty array",
			output:        []byte(`[]`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:    "invalid json",
			output:  []byte(`not valid json`),
			wantErr: true,
		},
		{
			name: "single error",
			output: []byte(`[{
				"filePath": "src/index.ts",
				"messages": [{
					"ruleId": "no-unused-vars",
					"severity": 2,
					"message": "'x' is defined but never used",
					"line": 10,
					"column": 5
				}],
				"errorCount": 1,
				"warningCount": 0
			}]`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/index.ts", line: 10, severity: "error", code: "no-unused-vars"},
			},
		},
		{
			name: "single warning",
			output: []byte(`[{
				"filePath": "src/utils.ts",
				"messages": [{
					"ruleId": "prefer-const",
					"severity": 1,
					"message": "Use const instead of let",
					"line": 5,
					"column": 1
				}],
				"errorCount": 0,
				"warningCount": 1
			}]`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/utils.ts", line: 5, severity: "warning", code: "prefer-const"},
			},
		},
		{
			name: "info severity (0)",
			output: []byte(`[{
				"filePath": "src/info.ts",
				"messages": [{
					"ruleId": "some-rule",
					"severity": 0,
					"message": "Info message",
					"line": 1,
					"column": 1
				}],
				"errorCount": 0,
				"warningCount": 0
			}]`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/info.ts", line: 1, severity: "info", code: "some-rule"},
			},
		},
		{
			name: "multiple files with issues",
			output: []byte(`[
				{"filePath": "src/file1.ts", "messages": [{"ruleId": "rule1", "severity": 2, "message": "Error 1", "line": 1, "column": 1}], "errorCount": 1, "warningCount": 0},
				{"filePath": "src/file2.ts", "messages": [{"ruleId": "rule2", "severity": 1, "message": "Warning 1", "line": 2, "column": 2}], "errorCount": 0, "warningCount": 1}
			]`),
			expectedCount: 2,
		},
		{
			name: "issue with fix suggestion",
			output: []byte(`[{
				"filePath": "src/fixable.ts",
				"messages": [{
					"ruleId": "semi",
					"severity": 2,
					"message": "Missing semicolon",
					"line": 10,
					"column": 20,
					"fix": {
						"range": [100, 100],
						"text": ";"
					}
				}],
				"errorCount": 1,
				"warningCount": 0
			}]`),
			expectedCount: 1,
		},
		{
			name: "file with no messages",
			output: []byte(`[{
				"filePath": "src/clean.ts",
				"messages": [],
				"errorCount": 0,
				"warningCount": 0
			}]`),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseESLintOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseESLintOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseESLintOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestESLintAnalyzer_ConvertSeverity(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})

	tests := []struct {
		input    int
		expected string
	}{
		{0, "info"},
		{1, "warning"},
		{2, "error"},
		{3, "warning"},  // unknown defaults to warning
		{-1, "warning"}, // unknown defaults to warning
	}

	for _, tt := range tests {
		t.Run(string(rune(tt.input+'0')), func(t *testing.T) {
			result := analyzer.convertSeverity(tt.input)
			if result != tt.expected {
				t.Errorf("convertSeverity(%d) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestESLintAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})

	tests := []struct {
		ruleID   string
		expected string
	}{
		{"@typescript-eslint/no-unused-vars", "typescript"},
		{"@typescript-eslint/explicit-function-return-type", "typescript"},
		{"react/jsx-uses-react", "react"},
		{"react/prop-types", "react"},
		{"import/no-unresolved", "imports"},
		{"import/order", "imports"},
		{"prefer-const", "style"},
		{"prefer-arrow-callback", "style"},
		{"no-unused-vars", "best-practices"},
		{"no-console", "best-practices"},
		{"semi", "other"},
		{"indent", "other"},
		{"", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.ruleID, func(t *testing.T) {
			result := analyzer.getCategory(tt.ruleID)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.ruleID, result, tt.expected)
			}
		})
	}
}

func TestESLintAnalyzer_HasTypeScriptFilePaths(t *testing.T) {
	analyzer := NewESLintAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasTypeScriptFilePaths("/non/existent/path")
	if result {
		t.Errorf("hasTypeScriptFilePaths() for non-existent path = true, want false")
	}

	// Test with current directory (no ts/js files in staticanalyzer)
	result = analyzer.hasTypeScriptFilePaths(".")
	if result {
		t.Errorf("hasTypeScriptFilePaths() for current dir = true, want false")
	}
}
