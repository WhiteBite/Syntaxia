package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestDartAnalyzeAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewDartAnalyzeAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"dart"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestDartAnalyzeAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewDartAnalyzeAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeDartAnalyze {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeDartAnalyze)
	}
}

func TestDartAnalyzeAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewDartAnalyzeAnalyzer(&mockLogger{})

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
			errMsg:  "DartAnalyzeAnalyzer only supports Dart language",
		},
		{
			name: "invalid language - python",
			config: &domain.StaticAnalyzerConfig{
				Language:    "python",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "DartAnalyzeAnalyzer only supports Dart language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "dart",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid dart config - no pubspec.yaml",
			config: &domain.StaticAnalyzerConfig{
				Language:    "dart",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "pubspec.yaml not found",
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

func TestDartAnalyzeAnalyzer_ParseDartAnalyzeOutput(t *testing.T) {
	analyzer := NewDartAnalyzeAnalyzer(&mockLogger{})

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
			name:          "single error",
			output:        []byte(`ERROR|COMPILE_TIME_ERROR|UNDEFINED_IDENTIFIER|lib/main.dart|10|5|3|Undefined name 'foo'`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "lib/main.dart", line: 10, severity: "error", code: "UNDEFINED_IDENTIFIER"},
			},
		},
		{
			name:          "single warning",
			output:        []byte(`WARNING|STATIC_WARNING|UNUSED_LOCAL_VARIABLE|lib/utils.dart|20|3|5|The value of the local variable 'x' isn't used`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "lib/utils.dart", line: 20, severity: "warning", code: "UNUSED_LOCAL_VARIABLE"},
			},
		},
		{
			name:          "single info",
			output:        []byte(`INFO|HINT|UNNECESSARY_CAST|lib/helper.dart|5|10|8|Unnecessary cast`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "lib/helper.dart", line: 5, severity: "info", code: "UNNECESSARY_CAST"},
			},
		},
		{
			name: "multiple issues",
			output: []byte(`ERROR|COMPILE_TIME_ERROR|UNDEFINED_IDENTIFIER|lib/main.dart|10|5|3|Error 1
WARNING|STATIC_WARNING|UNUSED_LOCAL_VARIABLE|lib/utils.dart|20|3|5|Warning 1`),
			expectedCount: 2,
		},
		{
			name:          "non-matching line (skipped)",
			output:        []byte(`Analyzing project...`),
			expectedCount: 0,
		},
		{
			name: "mixed valid and invalid lines",
			output: []byte(`Analyzing project...
ERROR|COMPILE_TIME_ERROR|UNDEFINED_IDENTIFIER|lib/main.dart|10|5|3|Error 1
Done.`),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseDartAnalyzeOutput(tt.output)
			if err != nil {
				t.Errorf("parseDartAnalyzeOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseDartAnalyzeOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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
