package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestKtlintAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"kotlin", "kt"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestKtlintAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeKtlint {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeKtlint)
	}
}

func TestKtlintAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

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
			errMsg:  "ktlint analyzer only supports Kotlin language",
		},
		{
			name: "invalid language - java",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ktlint analyzer only supports Kotlin language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "kotlin",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid kotlin config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "kotlin",
				ProjectPath: ".",
			},
			wantErr: false,
		},
		{
			name: "valid kt config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "kt",
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

func TestKtlintAnalyzer_ParseKtlintJSONOutput(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

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
			name:          "single issue",
			output:        []byte(`[{"file":"src/Main.kt","line":10,"column":5,"message":"Unexpected indentation","rule":"indent"}]`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Main.kt", line: 10, severity: "warning", code: "indent"},
			},
		},
		{
			name: "multiple issues",
			output: []byte(`[
				{"file":"src/File1.kt","line":5,"column":1,"message":"No wildcard imports","rule":"no-wildcard-imports"},
				{"file":"src/File2.kt","line":15,"column":10,"message":"Missing spacing","rule":"spacing"}
			]`),
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "src/File1.kt", line: 5, severity: "warning", code: "no-wildcard-imports"},
				{file: "src/File2.kt", line: 15, severity: "warning", code: "spacing"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseKtlintJSONOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKtlintJSONOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseKtlintJSONOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestKtlintAnalyzer_ParseKtlintTextOutput(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

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
			name:          "single issue",
			output:        "src/Main.kt:10:5: Unexpected indentation (indent)",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Main.kt", line: 10, severity: "warning", code: "indent"},
			},
		},
		{
			name: "multiple issues",
			output: `src/File1.kt:5:1: No wildcard imports (no-wildcard-imports)
src/File2.kt:15:10: Missing spacing (spacing)`,
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "src/File1.kt", line: 5, severity: "warning", code: "no-wildcard-imports"},
				{file: "src/File2.kt", line: 15, severity: "warning", code: "spacing"},
			},
		},
		{
			name:          "non-kt line ignored",
			output:        "Some random text without .kt extension",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := analyzer.parseKtlintTextOutput(tt.output)

			if len(issues) != tt.expectedCount {
				t.Errorf("parseKtlintTextOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestKtlintAnalyzer_CategorizeRule(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		rule     string
		expected string
	}{
		{"", "other"},
		{"indent", "formatting"},
		{"spacing", "formatting"},
		{"no-wildcard-imports", "imports"},
		{"import-ordering", "imports"},
		{"naming", "naming"},
		{"package-naming", "naming"},
		{"comment-spacing", "formatting"},
		{"max-line-length", "formatting"},
		{"no-unused-imports", "imports"},
		{"unknown-rule", "style"},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			result := analyzer.categorizeRule(tt.rule)
			if result != tt.expected {
				t.Errorf("categorizeRule(%q) = %q, want %q", tt.rule, result, tt.expected)
			}
		})
	}
}

func TestKtlintAnalyzer_ParseInt(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected int
	}{
		{"10", 10},
		{"0", 0},
		{"123", 123},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := analyzer.parseInt(tt.input)
			if result != tt.expected {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestKtlintAnalyzer_ParseKtlintLine(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		expected *expectedIssue
	}{
		{
			name:     "valid line with rule",
			line:     "src/Main.kt:10:5: Unexpected indentation (indent)",
			expected: &expectedIssue{file: "src/Main.kt", line: 10, severity: "warning", code: "indent"},
		},
		{
			name:     "line without rule in parentheses",
			line:     "src/File.kt:20:1: Some message without rule",
			expected: &expectedIssue{file: "src/File.kt", line: 20, severity: "warning", code: ""},
		},
		{
			name:     "incomplete line",
			line:     "src/File.kt:10",
			expected: nil,
		},
		{
			name:     "empty line",
			line:     "",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.parseKtlintLine(tt.line)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseKtlintLine(%q) = %v, want nil", tt.line, result)
				}
				return
			}
			if result == nil {
				t.Errorf("parseKtlintLine(%q) = nil, want non-nil", tt.line)
				return
			}
			if result.File != tt.expected.file {
				t.Errorf("parseKtlintLine(%q).File = %q, want %q", tt.line, result.File, tt.expected.file)
			}
			if result.Line != tt.expected.line {
				t.Errorf("parseKtlintLine(%q).Line = %d, want %d", tt.line, result.Line, tt.expected.line)
			}
		})
	}
}

func TestKtlintAnalyzer_ParseGradleKtlintOutput(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        string
		expectedCount int
	}{
		{
			name:          "empty output",
			output:        "",
			expectedCount: 0,
		},
		{
			name:          "gradle task output only",
			output:        "> Task :ktlintCheck\nBUILD SUCCESSFUL",
			expectedCount: 0,
		},
		{
			name: "gradle output with issues",
			output: `> Task :ktlintCheck
src/Main.kt:10:5: Unexpected indentation (indent)
src/Utils.kt:20:1: No wildcard imports (no-wildcard-imports)
BUILD FAILED`,
			expectedCount: 2,
		},
		{
			name:          "line starting with > should be ignored",
			output:        "> src/Main.kt:10:5: This should be ignored",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := analyzer.parseGradleKtlintOutput(tt.output)
			if len(issues) != tt.expectedCount {
				t.Errorf("parseGradleKtlintOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestKtlintAnalyzer_CheckKtlintInstalled(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	// This test just verifies the method doesn't panic
	// The actual result depends on whether ktlint is installed
	_ = analyzer.checkKtlintInstalled()
}

func TestKtlintAnalyzer_HasGradleKtlint(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasGradleKtlint("/non/existent/path")
	if result {
		t.Errorf("hasGradleKtlint() for non-existent path = true, want false")
	}

	// Test with current directory (no ktlint in build.gradle)
	result = analyzer.hasGradleKtlint(".")
	if result {
		t.Errorf("hasGradleKtlint() for current dir = true, want false (no ktlint plugin)")
	}
}

func TestKtlintAnalyzer_CategorizeRule_AllCategories(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		rule     string
		expected string
	}{
		{"", "other"},
		{"indent", "formatting"},
		{"standard:indent", "formatting"},
		{"spacing", "formatting"},
		{"argument-list-spacing", "formatting"},
		{"import-ordering", "imports"},
		{"no-wildcard-imports", "imports"},
		{"package-naming", "naming"},
		{"class-naming", "naming"},
		{"comment-spacing", "formatting"},
		{"max-line-length", "formatting"},
		{"no-unused-imports", "imports"},
		{"unknown-rule", "style"},
		{"custom-rule", "style"},
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			result := analyzer.categorizeRule(tt.rule)
			if result != tt.expected {
				t.Errorf("categorizeRule(%q) = %q, want %q", tt.rule, result, tt.expected)
			}
		})
	}
}

func TestKtlintAnalyzer_GetGradleCommand(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	// Test with non-existent path - should return empty or gradle
	result := analyzer.getGradleCommand("/non/existent/path")
	// Result depends on whether gradle is installed globally
	t.Logf("getGradleCommand() for non-existent path = %q", result)
}

func TestKtlintAnalyzer_HasGradleWrapper(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasGradleWrapper("/non/existent/path")
	if result {
		t.Errorf("hasGradleWrapper() for non-existent path = true, want false")
	}

	// Test with current directory (no gradle wrapper)
	result = analyzer.hasGradleWrapper(".")
	if result {
		t.Errorf("hasGradleWrapper() for current dir = true, want false")
	}
}

func TestKtlintAnalyzer_ParseKtlintLine_AllFormats(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantFile string
		wantLine int
		wantCode string
	}{
		{
			name:     "standard format with rule",
			line:     "src/Main.kt:10:5: Unexpected indentation (indent)",
			wantNil:  false,
			wantFile: "src/Main.kt",
			wantLine: 10,
			wantCode: "indent",
		},
		{
			name:     "format without rule",
			line:     "src/File.kt:20:1: Some message",
			wantNil:  false,
			wantFile: "src/File.kt",
			wantLine: 20,
			wantCode: "",
		},
		{
			name:    "incomplete line - only file and line",
			line:    "src/File.kt:10",
			wantNil: true,
		},
		{
			name:    "empty line",
			line:    "",
			wantNil: true,
		},
		{
			name:    "line with only two parts",
			line:    "src/File.kt:10:",
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.parseKtlintLine(tt.line)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parseKtlintLine(%q) = %v, want nil", tt.line, result)
				}
				return
			}
			if result == nil {
				t.Errorf("parseKtlintLine(%q) = nil, want non-nil", tt.line)
				return
			}
			if result.File != tt.wantFile {
				t.Errorf("File = %q, want %q", result.File, tt.wantFile)
			}
			if result.Line != tt.wantLine {
				t.Errorf("Line = %d, want %d", result.Line, tt.wantLine)
			}
			if result.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", result.Code, tt.wantCode)
			}
		})
	}
}


func TestKtlintAnalyzer_ParseKtlintJSONOutput_AllCases(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        []byte
		expectedCount int
		wantErr       bool
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
			name: "multiple issues",
			output: []byte(`[
				{"file":"src/File1.kt","line":5,"column":1,"message":"No wildcard imports","rule":"no-wildcard-imports"},
				{"file":"src/File2.kt","line":15,"column":10,"message":"Missing spacing","rule":"spacing"},
				{"file":"src/File3.kt","line":25,"column":5,"message":"Unexpected indentation","rule":"indent"}
			]`),
			expectedCount: 3,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseKtlintJSONOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseKtlintJSONOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(issues) != tt.expectedCount {
				t.Errorf("parseKtlintJSONOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestKtlintAnalyzer_ParseKtlintTextOutput_AllCases(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        string
		expectedCount int
	}{
		{
			name:          "empty output",
			output:        "",
			expectedCount: 0,
		},
		{
			name: "multiple issues",
			output: `src/File1.kt:5:1: No wildcard imports (no-wildcard-imports)
src/File2.kt:15:10: Missing spacing (spacing)
src/File3.kt:25:5: Unexpected indentation (indent)`,
			expectedCount: 3,
		},
		{
			name:          "non-kt line ignored",
			output:        "Some random text without .kt extension",
			expectedCount: 0,
		},
		{
			name: "mixed valid and invalid lines",
			output: `src/File1.kt:5:1: Valid issue (rule1)
Invalid line without .kt
src/File2.kt:10:1: Another valid issue (rule2)`,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := analyzer.parseKtlintTextOutput(tt.output)
			if len(issues) != tt.expectedCount {
				t.Errorf("parseKtlintTextOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}


func TestKtlintAnalyzer_CategorizeRule_Documentation(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	tests := []struct {
		rule     string
		expected string
	}{
		{"comment-wrapping", "documentation"}, // contains "comment" but not "spacing"
		{"kdoc-wrapping", "style"},            // doesn't contain "comment"
		{"comment-spacing", "formatting"},     // "spacing" is checked before "comment"
	}

	for _, tt := range tests {
		t.Run(tt.rule, func(t *testing.T) {
			result := analyzer.categorizeRule(tt.rule)
			if result != tt.expected {
				t.Errorf("categorizeRule(%q) = %q, want %q", tt.rule, result, tt.expected)
			}
		})
	}
}

func TestKtlintAnalyzer_ParseKtlintLine_WithNestedParentheses(t *testing.T) {
	analyzer := NewKtlintAnalyzer(&mockLogger{})

	// Test with message containing parentheses
	line := "src/Main.kt:10:5: Function name should be lowercase (see docs) (naming)"

	result := analyzer.parseKtlintLine(line)
	if result == nil {
		t.Error("parseKtlintLine() returned nil")
		return
	}

	if result.Code != "naming" {
		t.Errorf("Code = %q, want %q", result.Code, "naming")
	}
}
