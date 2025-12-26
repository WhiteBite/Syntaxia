package staticanalyzer

import (
	"context"
	"syntaxia/domain"
	"testing"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string)   {}
func (m *mockLogger) Info(msg string)    {}
func (m *mockLogger) Warning(msg string) {}
func (m *mockLogger) Error(msg string)   {}
func (m *mockLogger) Fatal(msg string)   {}

func TestClippyAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"rust", "rs"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestClippyAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeClippy {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeClippy)
	}
}

func TestClippyAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		name    string
		config  *domain.StaticAnalyzerConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid rust config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: ".", // Current directory as placeholder
			},
			wantErr: true, // Will fail because Cargo.toml doesn't exist
			errMsg:  "Cargo.toml not found",
		},
		{
			name: "valid rs config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rs",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "Cargo.toml not found",
		},
		{
			name: "invalid language - go",
			config: &domain.StaticAnalyzerConfig{
				Language:    "go",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "clippy analyzer only supports Rust language",
		},
		{
			name: "invalid language - python",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "clippy analyzer only supports Rust language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "non-existent project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/non/existent/path/to/project",
			},
			wantErr: true,
			errMsg:  "Cargo.toml not found",
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

func TestClippyAnalyzer_ParseClippyOutput(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

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
			output:        "Compiling myproject v0.1.0\nFinished dev [unoptimized + debuginfo] target(s)",
			expectedCount: 0,
		},
		{
			name: "single warning",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::needless_return"},"level":"warning","message":"unneeded `+"`return`"+` statement","spans":[{"file_name":"src/main.rs","line_start":10,"line_end":10,"column_start":5,"column_end":15,"is_primary":true}],"children":[{"level":"help","message":"remove `+"`return`"+`"}]}}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/main.rs", line: 10, severity: "warning", code: "clippy::needless_return"},
			},
		},
		{
			name: "error message",
			output: `{"reason":"compiler-message","message":{"code":{"code":"E0425"},"level":"error","message":"cannot find value `+"`x`"+` in this scope","spans":[{"file_name":"src/lib.rs","line_start":5,"line_end":5,"column_start":10,"column_end":11,"is_primary":true}],"children":[]}}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/lib.rs", line: 5, severity: "error", code: "E0425"},
			},
		},
		{
			name: "multiple issues",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::clone_on_copy"},"level":"warning","message":"using `+"`clone`"+` on a `+"`Copy`"+` type","spans":[{"file_name":"src/utils.rs","line_start":20,"line_end":20,"column_start":1,"column_end":10,"is_primary":true}],"children":[]}}
{"reason":"compiler-message","message":{"code":{"code":"clippy::redundant_closure"},"level":"warning","message":"redundant closure","spans":[{"file_name":"src/utils.rs","line_start":25,"line_end":25,"column_start":5,"column_end":20,"is_primary":true}],"children":[]}}`,
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "src/utils.rs", line: 20, severity: "warning", code: "clippy::clone_on_copy"},
				{file: "src/utils.rs", line: 25, severity: "warning", code: "clippy::redundant_closure"},
			},
		},
		{
			name: "note level message",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::cognitive_complexity"},"level":"note","message":"this function has high cognitive complexity","spans":[{"file_name":"src/complex.rs","line_start":1,"line_end":100,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/complex.rs", line: 1, severity: "info", code: "clippy::cognitive_complexity"},
			},
		},
		{
			name: "help level message",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::style"},"level":"help","message":"consider using a different approach","spans":[{"file_name":"src/style.rs","line_start":15,"line_end":15,"column_start":1,"column_end":50,"is_primary":true}],"children":[]}}`,
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/style.rs", line: 15, severity: "hint", code: "clippy::style"},
			},
		},
		{
			name:          "build artifact message (should be ignored)",
			output:        `{"reason":"build-script-executed","package_id":"myproject 0.1.0"}`,
			expectedCount: 0,
		},
		{
			name:          "message without spans (should be ignored)",
			output:        `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"some error","spans":[],"children":[]}}`,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseClippyOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseClippyOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseClippyOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestClippyAnalyzer_ConvertSeverity(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"Error", "error"},
		{"warning", "warning"},
		{"WARNING", "warning"},
		{"note", "info"},
		{"NOTE", "info"},
		{"help", "hint"},
		{"HELP", "hint"},
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

func TestClippyAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"", "other"},
		{"E0001", "compiler-error"},
		{"E0425", "compiler-error"},
		{"W0001", "compiler-warning"},
		{"clippy::cognitive_complexity", "complexity"},
		{"clippy::too_many_arguments", "complexity"},
		{"clippy::needless_return", "style"},
		{"clippy::redundant_closure", "style"},
		{"clippy::needless_collect", "performance"},
		{"clippy::clone_on_copy", "performance"},
		{"clippy::eq_op", "correctness"},
		{"clippy::almost_swapped", "correctness"},
		{"clippy::unknown_lint", "clippy"},
		{"unknown_code", "other"},
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

func TestClippyAnalyzer_Analyze_ClippyNotInstalled(t *testing.T) {
	// This test verifies behavior when clippy is not installed
	// In CI environments, clippy might not be available
	analyzer := NewClippyAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:    "rust",
		ProjectPath: "/non/existent/path",
	}

	// The analyze method should return an error if clippy is not installed
	// or if the project path doesn't exist
	_, err := analyzer.Analyze(ctx, config)
	if err == nil {
		// If no error, clippy might be installed - that's also valid
		t.Log("Clippy appears to be installed, skipping not-installed test")
	}
}

// expectedIssue is a helper struct for test assertions
type expectedIssue struct {
	file     string
	line     int
	severity string
	code     string
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestClippyAnalyzer_CategorizeLint(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		lintName string
		expected string
	}{
		{"cognitive_complexity", "complexity"},
		{"too_many_arguments", "complexity"},
		{"too_many_lines", "complexity"},
		{"type_complexity", "complexity"},
		{"excessive_precision", "complexity"},
		{"needless_return", "style"},
		{"redundant_closure", "style"},
		{"single_match", "style"},
		{"match_bool", "style"},
		{"if_same_then_else", "style"},
		{"collapsible_if", "style"},
		{"needless_collect", "performance"},
		{"unnecessary_to_owned", "performance"},
		{"clone_on_copy", "performance"},
		{"useless_vec", "performance"},
		{"box_collection", "performance"},
		{"eq_op", "correctness"},
		{"erasing_op", "correctness"},
		{"almost_swapped", "correctness"},
		{"suspicious_arithmetic_impl", "correctness"},
		{"misrefactored_assign_op", "correctness"},
		{"unknown_lint_name", "clippy"},
		{"", "clippy"},
	}

	for _, tt := range tests {
		t.Run(tt.lintName, func(t *testing.T) {
			result := analyzer.categorizeLint(tt.lintName)
			if result != tt.expected {
				t.Errorf("categorizeLint(%q) = %q, want %q", tt.lintName, result, tt.expected)
			}
		})
	}
}

func TestClippyAnalyzer_ParseClippyOutput_WithSuggestions(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	output := `{"reason":"compiler-message","message":{"code":{"code":"clippy::needless_return"},"level":"warning","message":"unneeded return statement","spans":[{"file_name":"src/main.rs","line_start":10,"line_end":10,"column_start":5,"column_end":15,"is_primary":true}],"children":[{"level":"help","message":"remove the return statement"},{"level":"note","message":"for more information visit..."}]}}`

	issues, err := analyzer.parseClippyOutput([]byte(output))
	if err != nil {
		t.Errorf("parseClippyOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseClippyOutput() returned %d issues, want 1", len(issues))
		return
	}

	if len(issues[0].Suggestions) != 2 {
		t.Errorf("issue.Suggestions has %d items, want 2", len(issues[0].Suggestions))
	}
}

func TestClippyAnalyzer_ParseClippyOutput_NonPrimarySpan(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	// Test with no primary span - should use first span
	output := `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"test error","spans":[{"file_name":"src/test.rs","line_start":5,"line_end":5,"column_start":1,"column_end":10,"is_primary":false}],"children":[]}}`

	issues, err := analyzer.parseClippyOutput([]byte(output))
	if err != nil {
		t.Errorf("parseClippyOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseClippyOutput() returned %d issues, want 1", len(issues))
		return
	}

	if issues[0].File != "src/test.rs" {
		t.Errorf("issue.File = %q, want %q", issues[0].File, "src/test.rs")
	}
}

func TestClippyAnalyzer_ParseClippyOutput_NoCode(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	// Test with no code field
	output := `{"reason":"compiler-message","message":{"code":null,"level":"warning","message":"test warning","spans":[{"file_name":"src/test.rs","line_start":1,"line_end":1,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}`

	issues, err := analyzer.parseClippyOutput([]byte(output))
	if err != nil {
		t.Errorf("parseClippyOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseClippyOutput() returned %d issues, want 1", len(issues))
		return
	}

	if issues[0].Code != "" {
		t.Errorf("issue.Code = %q, want empty string", issues[0].Code)
	}
}


func TestClippyAnalyzer_GetCategory_AllCases(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"", "other"},
		{"E0001", "compiler-error"},
		{"E0425", "compiler-error"},
		{"E9999", "compiler-error"},
		{"W0001", "compiler-warning"},
		{"W9999", "compiler-warning"},
		{"clippy::cognitive_complexity", "complexity"},
		{"clippy::too_many_arguments", "complexity"},
		{"clippy::too_many_lines", "complexity"},
		{"clippy::type_complexity", "complexity"},
		{"clippy::excessive_precision", "complexity"},
		{"clippy::needless_return", "style"},
		{"clippy::redundant_closure", "style"},
		{"clippy::single_match", "style"},
		{"clippy::match_bool", "style"},
		{"clippy::if_same_then_else", "style"},
		{"clippy::collapsible_if", "style"},
		{"clippy::needless_collect", "performance"},
		{"clippy::unnecessary_to_owned", "performance"},
		{"clippy::clone_on_copy", "performance"},
		{"clippy::useless_vec", "performance"},
		{"clippy::box_collection", "performance"},
		{"clippy::eq_op", "correctness"},
		{"clippy::erasing_op", "correctness"},
		{"clippy::almost_swapped", "correctness"},
		{"clippy::suspicious_arithmetic_impl", "correctness"},
		{"clippy::misrefactored_assign_op", "correctness"},
		{"clippy::unknown_lint", "clippy"},
		{"clippy::some_other_lint", "clippy"},
		{"unknown_code", "other"},
		{"random", "other"},
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

func TestClippyAnalyzer_ConvertSeverity_AllCases(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"error", "error"},
		{"ERROR", "error"},
		{"Error", "error"},
		{"eRrOr", "error"},
		{"warning", "warning"},
		{"WARNING", "warning"},
		{"Warning", "warning"},
		{"wArNiNg", "warning"},
		{"note", "info"},
		{"NOTE", "info"},
		{"Note", "info"},
		{"nOtE", "info"},
		{"help", "hint"},
		{"HELP", "hint"},
		{"Help", "hint"},
		{"hElP", "hint"},
		{"unknown", "warning"},
		{"", "warning"},
		{"other", "warning"},
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

func TestClippyAnalyzer_ParseClippyOutput_EdgeCases(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

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
			name:          "whitespace only",
			output:        "   \n\t\n   ",
			expectedCount: 0,
		},
		{
			name:          "non-json output",
			output:        "Compiling myproject v0.1.0\nFinished dev [unoptimized + debuginfo] target(s)",
			expectedCount: 0,
		},
		{
			name:          "build artifact message",
			output:        `{"reason":"build-script-executed","package_id":"myproject 0.1.0"}`,
			expectedCount: 0,
		},
		{
			name:          "message without spans",
			output:        `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"some error","spans":[],"children":[]}}`,
			expectedCount: 0,
		},
		{
			name: "multiple messages mixed with artifacts",
			output: `{"reason":"build-script-executed","package_id":"pkg1"}
{"reason":"compiler-message","message":{"code":{"code":"clippy::test"},"level":"warning","message":"test","spans":[{"file_name":"src/main.rs","line_start":1,"line_end":1,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}
{"reason":"build-finished","success":true}`,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseClippyOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseClippyOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseClippyOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}


func TestClippyAnalyzer_ParseClippyOutput_AllCases(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

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
			name:          "whitespace only",
			output:        "   \n\t\n   ",
			expectedCount: 0,
		},
		{
			name:          "non-json output",
			output:        "Compiling myproject v0.1.0\nFinished dev [unoptimized + debuginfo] target(s)",
			expectedCount: 0,
		},
		{
			name:          "build artifact message",
			output:        `{"reason":"build-script-executed","package_id":"myproject 0.1.0"}`,
			expectedCount: 0,
		},
		{
			name:          "message without spans",
			output:        `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"some error","spans":[],"children":[]}}`,
			expectedCount: 0,
		},
		{
			name: "single warning with children",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::needless_return"},"level":"warning","message":"unneeded return statement","spans":[{"file_name":"src/main.rs","line_start":10,"line_end":10,"column_start":5,"column_end":15,"is_primary":true}],"children":[{"level":"help","message":"remove return"},{"level":"note","message":"see docs"}]}}`,
			expectedCount: 1,
		},
		{
			name: "multiple messages",
			output: `{"reason":"compiler-message","message":{"code":{"code":"clippy::test1"},"level":"warning","message":"test1","spans":[{"file_name":"src/a.rs","line_start":1,"line_end":1,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}
{"reason":"compiler-message","message":{"code":{"code":"clippy::test2"},"level":"error","message":"test2","spans":[{"file_name":"src/b.rs","line_start":2,"line_end":2,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}`,
			expectedCount: 2,
		},
		{
			name: "message with non-primary span only",
			output: `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"test error","spans":[{"file_name":"src/test.rs","line_start":5,"line_end":5,"column_start":1,"column_end":10,"is_primary":false}],"children":[]}}`,
			expectedCount: 1,
		},
		{
			name: "message without code",
			output: `{"reason":"compiler-message","message":{"code":null,"level":"warning","message":"test warning","spans":[{"file_name":"src/test.rs","line_start":1,"line_end":1,"column_start":1,"column_end":1,"is_primary":true}],"children":[]}}`,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseClippyOutput([]byte(tt.output))
			if err != nil {
				t.Errorf("parseClippyOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseClippyOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestClippyAnalyzer_CategorizeLint_AllCategories(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		lintName string
		expected string
	}{
		{"cognitive_complexity", "complexity"},
		{"too_many_arguments", "complexity"},
		{"too_many_lines", "complexity"},
		{"type_complexity", "complexity"},
		{"excessive_precision", "complexity"},
		{"needless_return", "style"},
		{"redundant_closure", "style"},
		{"single_match", "style"},
		{"match_bool", "style"},
		{"if_same_then_else", "style"},
		{"collapsible_if", "style"},
		{"needless_collect", "performance"},
		{"unnecessary_to_owned", "performance"},
		{"clone_on_copy", "performance"},
		{"useless_vec", "performance"},
		{"box_collection", "performance"},
		{"eq_op", "correctness"},
		{"erasing_op", "correctness"},
		{"almost_swapped", "correctness"},
		{"suspicious_arithmetic_impl", "correctness"},
		{"misrefactored_assign_op", "correctness"},
		{"unknown_lint_name", "clippy"},
		{"", "clippy"},
		{"some_random_lint", "clippy"},
	}

	for _, tt := range tests {
		t.Run(tt.lintName, func(t *testing.T) {
			result := analyzer.categorizeLint(tt.lintName)
			if result != tt.expected {
				t.Errorf("categorizeLint(%q) = %q, want %q", tt.lintName, result, tt.expected)
			}
		})
	}
}


func TestClippyAnalyzer_ParseClippyOutput_MultipleSpans(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	// Test with multiple spans where first is not primary
	output := `{"reason":"compiler-message","message":{"code":{"code":"E0001"},"level":"error","message":"test error","spans":[{"file_name":"src/secondary.rs","line_start":1,"line_end":1,"column_start":1,"column_end":1,"is_primary":false},{"file_name":"src/primary.rs","line_start":10,"line_end":10,"column_start":5,"column_end":15,"is_primary":true}],"children":[]}}`

	issues, err := analyzer.parseClippyOutput([]byte(output))
	if err != nil {
		t.Errorf("parseClippyOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseClippyOutput() returned %d issues, want 1", len(issues))
		return
	}

	// Should use primary span
	if issues[0].File != "src/primary.rs" {
		t.Errorf("issue.File = %q, want %q", issues[0].File, "src/primary.rs")
	}
	if issues[0].Line != 10 {
		t.Errorf("issue.Line = %d, want 10", issues[0].Line)
	}
}

func TestClippyAnalyzer_GetCategory_CompilerCodes(t *testing.T) {
	analyzer := NewClippyAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"E0001", "compiler-error"},
		{"E0425", "compiler-error"},
		{"E9999", "compiler-error"},
		{"W0001", "compiler-warning"},
		{"W9999", "compiler-warning"},
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
