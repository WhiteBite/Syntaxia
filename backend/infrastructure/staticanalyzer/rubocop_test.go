package staticanalyzer

import (
	"context"
	"syntaxia/domain"
	"testing"
)

func TestRuboCopAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"ruby", "rb"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestRuboCopAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeRuboCop {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeRuboCop)
	}
}

func TestRuboCopAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

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
			errMsg:  "RuboCop analyzer only supports Ruby language",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ruby",
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

func TestRuboCopAnalyzer_ParseRuboCopOutput(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		output         []byte
		expectedCount  int
		expectedIssues []expectedIssue
		wantErr        bool
	}{
		{
			name:          "empty files",
			output:        []byte(`{"metadata":{},"files":[],"summary":{"offense_count":0}}`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:    "invalid json",
			output:  []byte(`not valid json`),
			wantErr: true,
		},
		{
			name: "single offense",
			output: []byte(`{
				"metadata":{"rubocop_version":"1.50.0"},
				"files":[{
					"path":"app/models/user.rb",
					"offenses":[{
						"severity":"warning",
						"message":"Line is too long. [125/120]",
						"cop_name":"Layout/LineLength",
						"corrected":false,
						"correctable":true,
						"location":{"start_line":10,"start_column":1,"last_line":10,"last_column":125,"length":125,"line":10,"column":1}
					}]
				}],
				"summary":{"offense_count":1}
			}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "app/models/user.rb", line: 10, severity: "warning", code: "Layout/LineLength"},
			},
		},
		{
			name: "error severity",
			output: []byte(`{
				"metadata":{},
				"files":[{
					"path":"lib/parser.rb",
					"offenses":[{
						"severity":"error",
						"message":"Syntax error",
						"cop_name":"Lint/Syntax",
						"corrected":false,
						"correctable":false,
						"location":{"line":5,"column":10}
					}]
				}],
				"summary":{"offense_count":1}
			}`),
			expectedCount: 1,
			expectedIssues: []expectedIssue{
				{file: "lib/parser.rb", line: 5, severity: "error", code: "Lint/Syntax"},
			},
		},
		{
			name: "multiple files with offenses",
			output: []byte(`{
				"metadata":{},
				"files":[
					{"path":"file1.rb","offenses":[{"severity":"warning","message":"Msg1","cop_name":"Style/Cop1","corrected":false,"correctable":false,"location":{"line":1,"column":1}}]},
					{"path":"file2.rb","offenses":[{"severity":"convention","message":"Msg2","cop_name":"Style/Cop2","corrected":false,"correctable":true,"location":{"line":2,"column":2}}]}
				],
				"summary":{"offense_count":2}
			}`),
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseRuboCopOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRuboCopOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseRuboCopOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
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

func TestRuboCopAnalyzer_ConvertSeverity(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"fatal", "error"},
		{"error", "error"},
		{"warning", "warning"},
		{"convention", "info"},
		{"refactor", "info"},
		{"info", "hint"},
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

func TestRuboCopAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		copName  string
		expected string
	}{
		{"", "other"},
		{"Layout/LineLength", "formatting"},
		{"Lint/UselessAssignment", "lint"},
		{"Metrics/CyclomaticComplexity", "complexity"},
		{"Naming/VariableName", "naming"},
		{"Security/Eval", "security"},
		{"Style/StringLiterals", "style"},
		{"Performance/Count", "performance"},
		{"Bundler/DuplicatedGem", "bundler"},
		{"Gemspec/RequiredRubyVersion", "gemspec"},
		{"Rails/ActiveRecordAliases", "rails"},
		{"RSpec/DescribeClass", "rspec"},
		{"CustomDepartment/SomeRule", "customdepartment"},
	}

	for _, tt := range tests {
		t.Run(tt.copName, func(t *testing.T) {
			result := analyzer.getCategory(tt.copName)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.copName, result, tt.expected)
			}
		})
	}
}

func TestRuboCopAnalyzer_BuildCommand(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		name           string
		config         *domain.StaticAnalyzerConfig
		expectedArgs   []string
		unexpectedArgs []string
	}{
		{
			name: "basic config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ruby",
				ProjectPath: "/project",
			},
			expectedArgs:   []string{"--format", "json", "."},
			unexpectedArgs: []string{"--only", "--except"},
		},
		{
			name: "with specific rules",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ruby",
				ProjectPath: "/project",
				Rules:       []string{"Style/StringLiterals", "Layout/LineLength"},
			},
			expectedArgs: []string{"--format", "json", "--only", "."},
		},
		{
			name: "with excluded rules",
			config: &domain.StaticAnalyzerConfig{
				Language:     "ruby",
				ProjectPath:  "/project",
				ExcludeRules: []string{"Metrics/MethodLength"},
			},
			expectedArgs: []string{"--format", "json", "--except", "."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := analyzer.buildCommand(tt.config)

			for _, expected := range tt.expectedArgs {
				found := false
				for _, arg := range args {
					if arg == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("buildCommand() missing expected arg %q in %v", expected, args)
				}
			}
		})
	}
}

func TestRuboCopAnalyzer_HasBundler(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasBundler("/non/existent/path")
	if result {
		t.Errorf("hasBundler() for non-existent path = true, want false")
	}

	// Test with current directory (no Gemfile with rubocop)
	result = analyzer.hasBundler(".")
	if result {
		t.Errorf("hasBundler() for current dir = true, want false")
	}
}

func TestRuboCopAnalyzer_ParseRuboCopOutput_WithCorrectableIssue(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	output := []byte(`{
		"metadata":{"rubocop_version":"1.50.0"},
		"files":[{
			"path":"app/models/user.rb",
			"offenses":[{
				"severity":"warning",
				"message":"Line is too long. [125/120]",
				"cop_name":"Layout/LineLength",
				"corrected":false,
				"correctable":true,
				"location":{"start_line":10,"start_column":1,"last_line":10,"last_column":125,"length":125,"line":10,"column":1}
			}]
		}],
		"summary":{"offense_count":1}
	}`)

	issues, err := analyzer.parseRuboCopOutput(output)
	if err != nil {
		t.Errorf("parseRuboCopOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseRuboCopOutput() returned %d issues, want 1", len(issues))
		return
	}

	if len(issues[0].Suggestions) == 0 {
		t.Errorf("issue.Suggestions is empty, want auto-correct suggestion")
	}
}

func TestRuboCopAnalyzer_ParseRuboCopOutput_JSONInMiddleOfOutput(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	// Test with JSON embedded in other output
	output := []byte(`Inspecting 5 files
{"metadata":{},"files":[],"summary":{"offense_count":0}}
5 files inspected, no offenses detected`)

	issues, err := analyzer.parseRuboCopOutput(output)
	if err != nil {
		t.Errorf("parseRuboCopOutput() error = %v", err)
		return
	}

	if len(issues) != 0 {
		t.Errorf("parseRuboCopOutput() returned %d issues, want 0", len(issues))
	}
}

func TestRuboCopAnalyzer_ValidateConfig_NoRubyFiles(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	config := &domain.StaticAnalyzerConfig{
		Language:    "ruby",
		ProjectPath: ".",
	}

	err := analyzer.ValidateConfig(config)
	if err == nil {
		t.Log("ValidateConfig() returned nil - Ruby files may exist in current directory")
	} else if !contains(err.Error(), "no Ruby files found") {
		t.Errorf("ValidateConfig() error = %q, want to contain 'no Ruby files found'", err.Error())
	}
}


func TestRuboCopAnalyzer_GetCategory_AllDepartments(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		copName  string
		expected string
	}{
		{"", "other"},
		{"Layout/LineLength", "formatting"},
		{"Layout/IndentationWidth", "formatting"},
		{"Lint/UselessAssignment", "lint"},
		{"Lint/Debugger", "lint"},
		{"Metrics/CyclomaticComplexity", "complexity"},
		{"Metrics/MethodLength", "complexity"},
		{"Naming/VariableName", "naming"},
		{"Naming/MethodName", "naming"},
		{"Security/Eval", "security"},
		{"Security/Open", "security"},
		{"Style/StringLiterals", "style"},
		{"Style/FrozenStringLiteralComment", "style"},
		{"Performance/Count", "performance"},
		{"Performance/Detect", "performance"},
		{"Bundler/DuplicatedGem", "bundler"},
		{"Bundler/OrderedGems", "bundler"},
		{"Gemspec/RequiredRubyVersion", "gemspec"},
		{"Rails/ActiveRecordAliases", "rails"},
		{"Rails/HttpPositionalArguments", "rails"},
		{"RSpec/DescribeClass", "rspec"},
		{"RSpec/ExampleLength", "rspec"},
		{"CustomDepartment/SomeRule", "customdepartment"},
		{"UnknownDept/Rule", "unknowndept"},
	}

	for _, tt := range tests {
		t.Run(tt.copName, func(t *testing.T) {
			result := analyzer.getCategory(tt.copName)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.copName, result, tt.expected)
			}
		})
	}
}

func TestRuboCopAnalyzer_ConvertSeverity_AllCases(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		input    string
		expected string
	}{
		{"fatal", "error"},
		{"FATAL", "error"},
		{"Fatal", "error"},
		{"error", "error"},
		{"ERROR", "error"},
		{"Error", "error"},
		{"warning", "warning"},
		{"WARNING", "warning"},
		{"Warning", "warning"},
		{"convention", "info"},
		{"CONVENTION", "info"},
		{"Convention", "info"},
		{"refactor", "info"},
		{"REFACTOR", "info"},
		{"Refactor", "info"},
		{"info", "hint"},
		{"INFO", "hint"},
		{"Info", "hint"},
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

func TestRuboCopAnalyzer_ParseRuboCopOutput_EdgeCases(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        []byte
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "empty JSON object",
			output:        []byte(`{"metadata":{},"files":[],"summary":{"offense_count":0}}`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:    "completely invalid JSON",
			output:  []byte(`not json at all`),
			wantErr: true,
		},
		{
			name: "JSON with text before",
			output: []byte(`Inspecting 5 files
{"metadata":{},"files":[],"summary":{"offense_count":0}}`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "JSON with text after",
			output: []byte(`{"metadata":{},"files":[],"summary":{"offense_count":0}}
5 files inspected, no offenses detected`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "multiple offenses in single file",
			output: []byte(`{
				"metadata":{},
				"files":[{
					"path":"app/models/user.rb",
					"offenses":[
						{"severity":"warning","message":"Msg1","cop_name":"Style/Cop1","corrected":false,"correctable":false,"location":{"line":1,"column":1}},
						{"severity":"error","message":"Msg2","cop_name":"Lint/Cop2","corrected":false,"correctable":false,"location":{"line":5,"column":10}},
						{"severity":"convention","message":"Msg3","cop_name":"Layout/Cop3","corrected":false,"correctable":true,"location":{"line":10,"column":1}}
					]
				}],
				"summary":{"offense_count":3}
			}`),
			expectedCount: 3,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseRuboCopOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRuboCopOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(issues) != tt.expectedCount {
				t.Errorf("parseRuboCopOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}


func TestRuboCopAnalyzer_ParseRuboCopOutput_AllCases(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        []byte
		expectedCount int
		wantErr       bool
	}{
		{
			name:          "empty files",
			output:        []byte(`{"metadata":{},"files":[],"summary":{"offense_count":0}}`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name:    "completely invalid JSON",
			output:  []byte(`not json at all`),
			wantErr: true,
		},
		{
			name: "JSON with text before and after",
			output: []byte(`Inspecting 5 files
{"metadata":{},"files":[],"summary":{"offense_count":0}}
5 files inspected, no offenses detected`),
			expectedCount: 0,
			wantErr:       false,
		},
		{
			name: "multiple files with offenses",
			output: []byte(`{
				"metadata":{"rubocop_version":"1.50.0"},
				"files":[
					{"path":"file1.rb","offenses":[
						{"severity":"warning","message":"Msg1","cop_name":"Style/Cop1","corrected":false,"correctable":false,"location":{"line":1,"column":1}},
						{"severity":"error","message":"Msg2","cop_name":"Lint/Cop2","corrected":false,"correctable":true,"location":{"line":5,"column":10}}
					]},
					{"path":"file2.rb","offenses":[
						{"severity":"convention","message":"Msg3","cop_name":"Layout/Cop3","corrected":false,"correctable":false,"location":{"line":10,"column":1}}
					]}
				],
				"summary":{"offense_count":3}
			}`),
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name: "correctable offense adds suggestion",
			output: []byte(`{
				"metadata":{},
				"files":[{
					"path":"app/models/user.rb",
					"offenses":[{
						"severity":"warning",
						"message":"Line is too long",
						"cop_name":"Layout/LineLength",
						"corrected":false,
						"correctable":true,
						"location":{"line":10,"column":1}
					}]
				}],
				"summary":{"offense_count":1}
			}`),
			expectedCount: 1,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseRuboCopOutput(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRuboCopOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(issues) != tt.expectedCount {
				t.Errorf("parseRuboCopOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestRuboCopAnalyzer_BuildCommand_AllCases(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	tests := []struct {
		name         string
		config       *domain.StaticAnalyzerConfig
		expectedArgs []string
	}{
		{
			name: "basic config",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ruby",
				ProjectPath: "/project",
			},
			expectedArgs: []string{"--format", "json", "."},
		},
		{
			name: "with specific rules",
			config: &domain.StaticAnalyzerConfig{
				Language:    "ruby",
				ProjectPath: "/project",
				Rules:       []string{"Style/StringLiterals", "Layout/LineLength"},
			},
			expectedArgs: []string{"--format", "json", "--only", "."},
		},
		{
			name: "with excluded rules",
			config: &domain.StaticAnalyzerConfig{
				Language:     "ruby",
				ProjectPath:  "/project",
				ExcludeRules: []string{"Metrics/MethodLength"},
			},
			expectedArgs: []string{"--format", "json", "--except", "."},
		},
		{
			name: "with both rules and excludes",
			config: &domain.StaticAnalyzerConfig{
				Language:     "ruby",
				ProjectPath:  "/project",
				Rules:        []string{"Style/StringLiterals"},
				ExcludeRules: []string{"Metrics/MethodLength"},
			},
			expectedArgs: []string{"--format", "json", "--only", "--except", "."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := analyzer.buildCommand(tt.config)

			for _, expected := range tt.expectedArgs {
				found := false
				for _, arg := range args {
					if arg == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("buildCommand() missing expected arg %q in %v", expected, args)
				}
			}
		})
	}
}


func TestRuboCopAnalyzer_Analyze_RuboCopNotInstalled(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})
	ctx := context.Background()

	config := &domain.StaticAnalyzerConfig{
		Language:    "ruby",
		ProjectPath: "/non/existent/path",
		Analyzer:    domain.StaticAnalyzerTypeRuboCop,
	}

	result, err := analyzer.Analyze(ctx, config)
	if err != nil {
		t.Errorf("Analyze() should not return error, got %v", err)
		return
	}

	if result == nil {
		t.Error("Analyze() should return a result")
		return
	}

	t.Logf("Analyze result: Success=%v, Error=%q", result.Success, result.Error)
}

func TestRuboCopAnalyzer_ParseRuboCopOutput_NonCorrectableIssue(t *testing.T) {
	analyzer := NewRuboCopAnalyzer(&mockLogger{})

	output := []byte(`{
		"metadata":{"rubocop_version":"1.50.0"},
		"files":[{
			"path":"app/models/user.rb",
			"offenses":[{
				"severity":"error",
				"message":"Syntax error",
				"cop_name":"Lint/Syntax",
				"corrected":false,
				"correctable":false,
				"location":{"line":10,"column":1}
			}]
		}],
		"summary":{"offense_count":1}
	}`)

	issues, err := analyzer.parseRuboCopOutput(output)
	if err != nil {
		t.Errorf("parseRuboCopOutput() error = %v", err)
		return
	}

	if len(issues) != 1 {
		t.Errorf("parseRuboCopOutput() returned %d issues, want 1", len(issues))
		return
	}

	// Non-correctable issue should have no suggestions
	if len(issues[0].Suggestions) != 0 {
		t.Errorf("issue.Suggestions has %d items, want 0", len(issues[0].Suggestions))
	}
}
