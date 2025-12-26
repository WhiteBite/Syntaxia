package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestDotnetFormatAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"csharp", "cs"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestDotnetFormatAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeDotnetFormat {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeDotnetFormat)
	}
}

func TestDotnetFormatAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

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
			errMsg:  "dotnet format analyzer only supports C# language",
		},
		{
			name: "invalid language - java",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "dotnet format analyzer only supports C# language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "csharp",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid csharp config - no project file",
			config: &domain.StaticAnalyzerConfig{
				Language:    "csharp",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no .csproj or .sln file found",
		},
		{
			name: "valid cs config - no project file",
			config: &domain.StaticAnalyzerConfig{
				Language:    "cs",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no .csproj or .sln file found",
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

func TestDotnetFormatAnalyzer_ParseDotnetFormatOutput(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

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
			name:          "standard issue format",
			output:        "src/Program.cs(10,5): warning IDE0001: Simplify name",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Program.cs", line: 10, severity: "warning", code: "IDE0001"},
			},
		},
		{
			name:          "error issue",
			output:        "src/Utils.cs(25,10): error CS0103: The name 'x' does not exist",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Utils.cs", line: 25, severity: "error", code: "CS0103"},
			},
		},
		{
			name:          "info issue",
			output:        "src/Helper.cs(5,1): info CA1000: Consider making method static",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Helper.cs", line: 5, severity: "info", code: "CA1000"},
			},
		},
		{
			name: "multiple issues",
			output: `src/File1.cs(10,5): warning IDE0001: Message 1
src/File2.cs(20,10): error CS0001: Message 2`,
			expectedCount: 2,
			expectedIssues: []expectedIssue{
				{file: "src/File1.cs", line: 10, severity: "warning", code: "IDE0001"},
				{file: "src/File2.cs", line: 20, severity: "error", code: "CS0001"},
			},
		},
		{
			name:          "would format output",
			output:        "Would format: src/NeedsFormat.cs",
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/NeedsFormat.cs", severity: "warning", code: "FORMAT001"},
			},
		},
		{
			name: "mixed output",
			output: `Building project...
src/Program.cs(10,5): warning IDE0001: Simplify name
Would format: src/Other.cs
Done.`,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseDotnetFormatOutput(tt.output)
			if err != nil {
				t.Errorf("parseDotnetFormatOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseDotnetFormatOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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
				if expected.line != 0 && issue.Line != expected.line {
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

func TestDotnetFormatAnalyzer_CategorizeIssue(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

	tests := []struct {
		code     string
		expected string
	}{
		{"", "other"},
		{"IDE0001", "style"},
		{"IDE0055", "style"},
		{"CS0001", "compiler"},
		{"CS0103", "compiler"},
		{"CA1000", "analysis"},
		{"CA2000", "analysis"},
		{"SA1000", "style"},
		{"SA1200", "style"},
		{"FORMAT001", "formatting"},
		{"UNKNOWN001", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			result := analyzer.categorizeIssue(tt.code)
			if result != tt.expected {
				t.Errorf("categorizeIssue(%q) = %q, want %q", tt.code, result, tt.expected)
			}
		})
	}
}

func TestDotnetFormatAnalyzer_ParseInt(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

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

func TestDotnetFormatAnalyzer_HasDotnetProject(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "non-existent path",
			path:     "/non/existent/path/to/project",
			expected: false,
		},
		{
			name:     "current directory without csproj",
			path:     ".",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.hasDotnetProject(tt.path)
			if result != tt.expected {
				t.Errorf("hasDotnetProject(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestDotnetFormatAnalyzer_ParseDotnetFormatOutput_EdgeCases(t *testing.T) {
	analyzer := NewDotnetFormatAnalyzer(&mockLogger{})

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
			name:          "build output without issues",
			output:        "Build started...\nBuild succeeded.\n0 Warning(s)\n0 Error(s)",
			expectedCount: 0,
		},
		{
			name: "multiple would format entries",
			output: `Would format: src/File1.cs
Would format: src/File2.cs
Would format: src/File3.cs`,
			expectedCount: 3,
		},
		{
			name: "mixed standard and would format",
			output: `src/Program.cs(10,5): warning IDE0001: Simplify name
Would format: src/Other.cs
src/Utils.cs(20,1): error CS0001: Error message`,
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseDotnetFormatOutput(tt.output)
			if err != nil {
				t.Errorf("parseDotnetFormatOutput() error = %v", err)
				return
			}
			if len(issues) != tt.expectedCount {
				t.Errorf("parseDotnetFormatOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}