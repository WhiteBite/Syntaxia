package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestStaticcheckAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"go"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestStaticcheckAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeStaticcheck {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeStaticcheck)
	}
}

func TestStaticcheckAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "invalid language - rust",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "staticcheck analyzer only supports Go language",
		},
		{
			name: "invalid language - python",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "staticcheck analyzer only supports Go language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid go config - no go files in staticanalyzer dir",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "/non/existent/path",
			},
			wantErr: true,
			errMsg:  "no Go files found",
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

func TestStaticcheckAnalyzer_ParseStaticcheckOutput(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         []byte
		expectedCount  int
		expectedIssues []expectedIssue
	}{
		{
			name:          "empty output",
			output:        []byte(``),
			expectedCount: 0,
		},
		{
			name:          "whitespace only",
			output:        []byte(`   `),
			expectedCount: 0,
		},
		{
			name: "single issue",
			output: []byte(`{"code":"SA1000","severity":"error","location":{"file":"main.go","line":10,"column":5},"message":"Invalid regexp"}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "main.go", line: 10, severity: "error", code: "SA1000"},
			},
		},
		{
			name: "warning severity",
			output: []byte(`{"code":"ST1000","severity":"warning","location":{"file":"utils.go","line":20,"column":1},"message":"Style issue"}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "utils.go", line: 20, severity: "warning", code: "ST1000"},
			},
		},
		{
			name: "info severity",
			output: []byte(`{"code":"U1000","severity":"info","location":{"file":"unused.go","line":5,"column":1},"message":"Unused code"}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "unused.go", line: 5, severity: "info", code: "U1000"},
			},
		},
		{
			name: "multiple issues",
			output: []byte(`{"code":"SA1000","severity":"error","location":{"file":"file1.go","line":10,"column":5},"message":"Error 1"}
{"code":"ST1000","severity":"warning","location":{"file":"file2.go","line":20,"column":1},"message":"Warning 1"}`),
			expectedCount: 2,
		},
		{
			name:          "invalid json line (skipped)",
			output:        []byte(`not valid json`),
			expectedCount: 0,
		},
		{
			name: "mixed valid and invalid lines",
			output: []byte(`invalid line
{"code":"SA1000","severity":"error","location":{"file":"main.go","line":10,"column":5},"message":"Error"}
another invalid line`),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseStaticcheckOutput(tt.output)
			if err != nil {
				t.Errorf("parseStaticcheckOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseStaticcheckOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestStaticcheckAnalyzer_ConvertSeverity(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"Error", "error"},
		{"warning", "warning"},
		{"WARNING", "warning"},
		{"info", "info"},
		{"INFO", "info"},
		{"unknown", "warning"},
		{"", "warning"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := analyzer.convertSeverity(tt.input)
			if result != tt.expected {
				t.Errorf("convertSeverity(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStaticcheckAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"SA1000", "static-analysis"},
		{"SA4006", "static-analysis"},
		{"ST1000", "style"},
		{"ST1003", "style"},
		{"U1000", "unused"},
		{"U1001", "unused"},
		{"QF1001", "quickfix"},
		{"QF1002", "quickfix"},
		{"unknown", "other"},
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

func TestStaticcheckAnalyzer_HasGoFilePaths(t *testing.T) {
	analyzer := NewStaticcheckAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasGoFilePaths("/non/existent/path")
	if result {
		t.Errorf("hasGoFilePaths() for non-existent path = true, want false")
	}

	// Note: The glob pattern **/*.go doesn't work recursively in Go's filepath.Glob
	// So hasGoFilePaths will return false for directories without .go files at root level
}
