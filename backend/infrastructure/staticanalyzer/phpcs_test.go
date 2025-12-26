package staticanalyzer

import (
	"context"
	"syntaxia/domain"
	"testing"
)

func TestPHPCSAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"php"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestPHPCSAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypePHPCS {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypePHPCS)
	}
}

func TestPHPCSAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

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
			errMsg:  "PHPCS analyzer only supports PHP language",
		},
		{
			name: "invalid language - java",
			config: &domain.StaticAnalyzerConfig{
				Language:    "java",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "PHPCS analyzer only supports PHP language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "php",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid php config - no php files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "php",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no PHP files found",
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

func TestPHPCSAnalyzer_ParsePHPCSOutput(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         []byte
		expectedCount  int
		expectedIssues []expectedIssue
		wantErr        bool
	}{
		{
			name:          "empty files",
			output:        []byte(`{"totals":{"errors":0,"warnings":0,"fixable":0},"files":{}}`),
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
			output: []byte(`{
				"totals":{"errors":1,"warnings":0,"fixable":0},
				"files":{
					"src/Controller.php":{
						"errors":1,
						"warnings":0,
						"messages":[{
							"message":"Missing class doc comment",
							"source":"PSR12.Classes.ClassInstantiation.MissingParentheses",
							"severity":5,
							"fixable":false,
							"type":"ERROR",
							"line":10,
							"column":5
						}]
					}
				}
			}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Controller.php", line: 10, severity: "error", code: "PSR12.Classes.ClassInstantiation.MissingParentheses"},
			},
		},
		{
			name: "single warning",
			output: []byte(`{
				"totals":{"errors":0,"warnings":1,"fixable":1},
				"files":{
					"src/Helper.php":{
						"errors":0,
						"warnings":1,
						"messages":[{
							"message":"Line exceeds 120 characters",
							"source":"Generic.Files.LineLength.TooLong",
							"severity":3,
							"fixable":true,
							"type":"WARNING",
							"line":25,
							"column":121
						}]
					}
				}
			}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "src/Helper.php", line: 25, severity: "warning", code: "Generic.Files.LineLength.TooLong"},
			},
		},
		{
			name: "multiple files with issues",
			output: []byte(`{
				"totals":{"errors":2,"warnings":1,"fixable":1},
				"files":{
					"src/File1.php":{
						"errors":1,
						"warnings":0,
						"messages":[{
							"message":"Error 1",
							"source":"PSR12.Error1",
							"severity":5,
							"fixable":false,
							"type":"ERROR",
							"line":5,
							"column":1
						}]
					},
					"src/File2.php":{
						"errors":1,
						"warnings":1,
						"messages":[
							{"message":"Error 2","source":"PSR12.Error2","severity":5,"fixable":false,"type":"ERROR","line":10,"column":1},
							{"message":"Warning 1","source":"Generic.Warning1","severity":3,"fixable":true,"type":"WARNING","line":15,"column":1}
						]
					}
				}
			}`),
			expectedCount: 3,
		},
		{
			name: "fixable issue adds suggestion",
			output: []byte(`{
				"totals":{"errors":0,"warnings":1,"fixable":1},
				"files":{
					"src/Fixable.php":{
						"errors":0,
						"warnings":1,
						"messages":[{
							"message":"Fixable issue",
							"source":"PSR12.Fixable",
							"severity":3,
							"fixable":true,
							"type":"WARNING",
							"line":1,
							"column":1
						}]
					}
				}
			}`),
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parsePHPCSOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePHPCSOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parsePHPCSOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestPHPCSAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	tests := []struct {
		source   string
		expected string
	}{
		{"", "other"},
		{"PSR1.Classes.ClassDeclaration", "style"},
		{"PSR2.Methods.MethodDeclaration", "style"},
		{"PSR12.Files.FileHeader", "style"},
		{"Generic.CodeAnalysis.UnusedFunctionParameter", "complexity"},
		{"Generic.Commenting.DocComment", "documentation"},
		{"Generic.Files.LineLength", "style"},
		{"Generic.Formatting.SpaceAfterCast", "style"},
		{"Generic.Functions.OpeningFunctionBraceBsdAllman", "style"},
		{"Generic.Metrics.CyclomaticComplexity", "complexity"},
		{"Generic.NamingConventions.UpperCaseConstantName", "naming"},
		{"Generic.PHP.DeprecatedFunctions", "correctness"},
		{"Generic.Strings.UnnecessaryStringConcat", "style"},
		{"Generic.WhiteSpace.DisallowTabIndent", "style"},
		{"Squiz.Arrays.ArrayDeclaration", "squiz"},
		{"PEAR.Commenting.FileComment", "pear"},
		{"Zend.Files.ClosingTag", "zend"},
		{"Unknown.Rule.Name", "other"},
		{"Generic", "generic"},
		{"Generic.Unknown", "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			result := analyzer.getCategory(tt.source)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.source, result, tt.expected)
			}
		})
	}
}

func TestPHPCSAnalyzer_FindPHPCSPath(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	// Test with non-existent path - should return global phpcs
	result := analyzer.findPHPCSPath("/non/existent/path")
	if result != "phpcs" {
		t.Errorf("findPHPCSPath() for non-existent path = %q, want %q", result, "phpcs")
	}
}

func TestPHPCSAnalyzer_GetCategory_AllStandards(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	tests := []struct {
		source   string
		expected string
	}{
		{"", "other"},
		{"PSR1.Classes.ClassDeclaration", "style"},
		{"PSR2.Methods.MethodDeclaration", "style"},
		{"PSR12.Files.FileHeader", "style"},
		{"Generic.CodeAnalysis.UnusedFunctionParameter", "complexity"},
		{"Generic.Commenting.DocComment", "documentation"},
		{"Generic.Files.LineLength", "style"},
		{"Generic.Formatting.SpaceAfterCast", "style"},
		{"Generic.Functions.OpeningFunctionBraceBsdAllman", "style"},
		{"Generic.Metrics.CyclomaticComplexity", "complexity"},
		{"Generic.NamingConventions.UpperCaseConstantName", "naming"},
		{"Generic.PHP.DeprecatedFunctions", "correctness"},
		{"Generic.Strings.UnnecessaryStringConcat", "style"},
		{"Generic.WhiteSpace.DisallowTabIndent", "style"},
		{"Squiz.Arrays.ArrayDeclaration", "squiz"},
		{"PEAR.Commenting.FileComment", "pear"},
		{"Zend.Files.ClosingTag", "zend"},
		{"Unknown.Rule.Name", "other"},
		{"Generic", "generic"},
		{"Generic.Unknown", "generic"},
	}

	for _, tt := range tests {
		t.Run(tt.source, func(t *testing.T) {
			result := analyzer.getCategory(tt.source)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.source, result, tt.expected)
			}
		})
	}
}

func TestPHPCSAnalyzer_ParsePHPCSOutput_FixableIssue(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	output := []byte(`{
		"totals":{"errors":0,"warnings":1,"fixable":1},
		"files":{
			"src/Fixable.php":{
				"errors":0,
				"warnings":1,
				"messages":[{
					"message":"Fixable issue",
					"source":"PSR12.Fixable",
					"severity":3,
					"fixable":true,
					"type":"WARNING",
					"line":1,
					"column":1
				}]
			}
		}
	}`)

	issues, err := analyzer.parsePHPCSOutput(output)
	if err != nil {
		t.Errorf("parsePHPCSOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parsePHPCSOutput() returned %d issues, want 1", len(issues))
		return
	}

	if len(issues[0].Suggestions) == 0 {
		t.Error("issue.Suggestions is empty, want auto-fix suggestion")
	}
}

func TestPHPCSAnalyzer_HasPHPFiles(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasPHPFiles("/non/existent/path")
	if result {
		t.Errorf("hasPHPFiles() for non-existent path = true, want false")
	}

	// Test with current directory (no PHP files)
	result = analyzer.hasPHPFiles(".")
	if result {
		t.Errorf("hasPHPFiles() for current dir = true, want false")
	}
}


func TestPHPCSAnalyzer_ParsePHPCSOutput_AllTypes(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        []byte
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "empty files",
			output:        []byte(`{"totals":{"errors":0,"warnings":0,"fixable":0},"files":{}}`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:    "invalid json",
			output:  []byte(`not valid json`),
			wantErr: true,
		},
		{
			name: "multiple files with multiple messages",
			output: []byte(`{
				"totals":{"errors":3,"warnings":2,"fixable":2},
				"files":{
					"src/File1.php":{
						"errors":2,
						"warnings":1,
						"messages":[
							{"message":"Error 1","source":"PSR12.Error1","severity":5,"fixable":false,"type":"ERROR","line":5,"column":1},
							{"message":"Error 2","source":"PSR12.Error2","severity":5,"fixable":true,"type":"ERROR","line":10,"column":1},
							{"message":"Warning 1","source":"Generic.Warning1","severity":3,"fixable":true,"type":"WARNING","line":15,"column":1}
						]
					},
					"src/File2.php":{
						"errors":1,
						"warnings":1,
						"messages":[
							{"message":"Error 3","source":"PSR12.Error3","severity":5,"fixable":false,"type":"ERROR","line":20,"column":1},
							{"message":"Warning 2","source":"Generic.Warning2","severity":3,"fixable":false,"type":"WARNING","line":25,"column":1}
						]
					}
				}
			}`),
			expectedCount: 5,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parsePHPCSOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePHPCSOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(issues) != tt.expectedCount {
				t.Errorf("parsePHPCSOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}


func TestPHPCSAnalyzer_Analyze_PHPCSNotInstalled(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:    "php",
		ProjectPath: "/non/existent/path",
		Analyzer:    domain.StaticAnalyzerTypePHPCS,
	}

	result, err := analyzer.Analyze(ctx, config)
	if err != nil {
		t.Errorf("Analyze() should not return error, got %v", err)
		return
	}

	if result == nil {
		t.Error("Analyze() should return a result even when PHPCS is not installed")
		return
	}

	// Result should be returned with error message if PHPCS is not installed
	t.Logf("Analyze result: Success=%v, Error=%q", result.Success, result.Error)
}

func TestPHPCSAnalyzer_ParsePHPCSOutput_EmptyOutput(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})

	// Test with empty byte slice
	issues, err := analyzer.parsePHPCSOutput([]byte{})
	if err == nil {
		t.Log("parsePHPCSOutput() with empty output returned no error")
	}
	if issues != nil && len(issues) > 0 {
		t.Errorf("parsePHPCSOutput() with empty output returned %d issues, want 0", len(issues))
	}
}

func TestPHPCSAnalyzer_Analyze_WithRules(t *testing.T) {
	analyzer := NewPHPCSAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:     "php",
		ProjectPath:  "/non/existent/path",
		Analyzer:     domain.StaticAnalyzerTypePHPCS,
		Rules:        []string{"PSR12", "PSR1"},
		ExcludeRules: []string{"Generic.Files.LineLength"},
	}

	result, err := analyzer.Analyze(ctx, config)
	if err != nil {
		t.Errorf("Analyze() should not return error, got %v", err)
		return
	}

	if result == nil {
		t.Error("Analyze() should return a result")
	}
}
