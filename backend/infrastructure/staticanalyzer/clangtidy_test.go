package staticanalyzer

import (
	"syntaxia/domain"
	"testing"
)

func TestClangTidyAnalyzer_GetSupportedLanguages(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})
	languages := analyzer.GetSupportedLanguages()

	expected := []string{"c", "cpp", "cc"}
	if len(languages) != len(expected) {
		t.Errorf("GetSupportedLanguages() returned %d languages, want %d", len(languages), len(expected))
	}

	for i, lang := range languages {
		if lang != expected[i] {
			t.Errorf("GetSupportedLanguages()[%d] = %q, want %q", i, lang, expected[i])
		}
	}
}

func TestClangTidyAnalyzer_GetAnalyzerType(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})
	analyzerType := analyzer.GetAnalyzerType()

	if analyzerType != domain.StaticAnalyzerTypeClangTidy {
		t.Errorf("GetAnalyzerType() = %q, want %q", analyzerType, domain.StaticAnalyzerTypeClangTidy)
	}
}

func TestClangTidyAnalyzer_ValidateConfig(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

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
			errMsg:  "ClangTidy analyzer only supports C/C++ languages",
		},
		{
			name: "invalid language - rust",
			config: &domain.StaticAnalyzerConfig{
				Language:    "rust",
				ProjectPath: "/some/path",
			},
			wantErr: true,
			errMsg:  "ClangTidy analyzer only supports C/C++ languages",
		},
		{
			name: "empty project path",
			config: &domain.StaticAnalyzerConfig{
				Language:    "cpp",
				ProjectPath: "",
			},
			wantErr: true,
			errMsg:  "project path is required",
		},
		{
			name: "valid c config - no c files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "c",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no C/C++ files found",
		},
		{
			name: "valid cpp config - no cpp files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "cpp",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no C/C++ files found",
		},
		{
			name: "valid cc config - no cc files",
			config: &domain.StaticAnalyzerConfig{
				Language:    "cc",
				ProjectPath: ".",
			},
			wantErr: true,
			errMsg:  "no C/C++ files found",
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

func TestClangTidyAnalyzer_ParseClangTidyOutput(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	tests := []struct {
		name          string
		output        []byte
		expectedCount int
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
			name:          "non-diagnostic line",
			output:        []byte(`Checking main.cpp...`),
			expectedCount: 0,
		},
		{
			name:          "line without DiagnosticName",
			output:        []byte(`{"Message": "Some message"}`),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues, err := analyzer.parseClangTidyOutput(tt.output)
			if err != nil {
				t.Errorf("parseClangTidyOutput() error = %v", err)
				return
			}

			if len(issues) != tt.expectedCount {
				t.Errorf("parseClangTidyOutput() returned %d issues, want %d", len(issues), tt.expectedCount)
			}
		})
	}
}

func TestClangTidyAnalyzer_GetCategory(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	tests := []struct {
		diagnosticName string
		expected       string
	}{
		{"unused-variable", "unused-code"},
		{"clang-analyzer-core.NullDereference", "null-safety"},
		{"memory-leak", "memory-management"},
		{"performance-unnecessary-copy", "performance"},
		{"modernize-use-auto", "modernization"},
		{"bugprone-use-after-move", "bug-prone"},
		{"misc-redundant-expression", "other"},
		{"", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.diagnosticName, func(t *testing.T) {
			result := analyzer.getCategory(tt.diagnosticName)
			if result != tt.expected {
				t.Errorf("getCategory(%q) = %q, want %q", tt.diagnosticName, result, tt.expected)
			}
		})
	}
}

func TestClangTidyAnalyzer_ExtractClangString(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		prefix   string
		expected string
	}{
		{
			name:     "extract DiagnosticName",
			line:     `{"DiagnosticName": "unused-variable", "Message": "test"}`,
			prefix:   "\"DiagnosticName\": \"",
			expected: "unused-variable",
		},
		{
			name:     "extract Message",
			line:     `{"DiagnosticName": "test", "Message": "This is a message"}`,
			prefix:   "\"Message\": \"",
			expected: "This is a message",
		},
		{
			name:     "prefix not found",
			line:     `{"DiagnosticName": "test"}`,
			prefix:   "\"NotFound\": \"",
			expected: "",
		},
		{
			name:     "empty line",
			line:     "",
			prefix:   "\"DiagnosticName\": \"",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractClangString(tt.line, tt.prefix)
			if result != tt.expected {
				t.Errorf("extractClangString() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestClangTidyAnalyzer_ExtractClangInt(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		prefix   string
		expected int
	}{
		{
			name:     "extract line number with comma",
			line:     `{"FilePathLineNumber": 42, "Column": 5}`,
			prefix:   "\"FilePathLineNumber\": ",
			expected: 42,
		},
		{
			name:     "extract line number with brace",
			line:     `{"FilePathLineNumber": 100}`,
			prefix:   "\"FilePathLineNumber\": ",
			expected: 100,
		},
		{
			name:     "prefix not found",
			line:     `{"Line": 10}`,
			prefix:   "\"NotFound\": ",
			expected: 0,
		},
		{
			name:     "empty line",
			line:     "",
			prefix:   "\"FilePathLineNumber\": ",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.extractClangInt(tt.line, tt.prefix)
			if result != tt.expected {
				t.Errorf("extractClangInt() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestClangTidyAnalyzer_ParseClangTidyLine(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	tests := []struct {
		name     string
		line     string
		wantNil  bool
		wantCode string
	}{
		{
			name:    "empty line",
			line:    "",
			wantNil: true,
		},
		{
			name:    "line without DiagnosticName",
			line:    `{"Message": "test"}`,
			wantNil: true,
		},
		{
			name:     "valid diagnostic line",
			line:     `{"DiagnosticName": "unused-variable", "Message": "Variable is unused", "FilePathPath": "main.cpp", "FilePathLineNumber": 10, "FilePathColumnStart": 5}`,
			wantNil:  false,
			wantCode: "unused-variable",
		},
		{
			name:    "diagnostic without file path",
			line:    `{"DiagnosticName": "test", "Message": "test"}`,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.parseClangTidyLine(tt.line)
			if tt.wantNil {
				if result != nil {
					t.Errorf("parseClangTidyLine() = %v, want nil", result)
				}
			} else {
				if result == nil {
					t.Errorf("parseClangTidyLine() = nil, want non-nil")
				} else if result.Code != tt.wantCode {
					t.Errorf("parseClangTidyLine().Code = %q, want %q", result.Code, tt.wantCode)
				}
			}
		})
	}
}

func TestClangTidyAnalyzer_HasCppFilePaths(t *testing.T) {
	analyzer := NewClangTidyAnalyzer(&mockLogger{})

	// Test with non-existent path
	result := analyzer.hasCppFilePaths("/non/existent/path")
	if result {
		t.Errorf("hasCppFilePaths() for non-existent path = true, want false")
	}

	// Test with current directory (no cpp files)
	result = analyzer.hasCppFilePaths(".")
	if result {
		t.Errorf("hasCppFilePaths() for current dir = true, want false")
	}
}
