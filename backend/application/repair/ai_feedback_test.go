package repair

import (
	"strings"
	"syntaxia/domain"
	"testing"
)

func TestNewAIFeedbackGenerator(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
	if len(gen.templates) == 0 {
		t.Error("expected templates to be registered")
	}
}

func TestAIFeedbackGenerator_GeneratePrompt(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	tests := []struct {
		name        string
		err         *domain.ErrorDetails
		codeContext string
		wantContain []string
	}{
		{
			name: "import error",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeImport,
				SourceFile: "main.go",
				LineNumber: 10,
				Message:    "undefined: fmt",
			},
			codeContext: "import \"os\"\n\nfunc main() {}",
			wantContain: []string{"import error", "main.go", "10", "undefined: fmt"},
		},
		{
			name: "syntax error",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeSyntax,
				SourceFile: "app.ts",
				LineNumber: 5,
				Message:    "unexpected token",
			},
			codeContext: "const x = {",
			wantContain: []string{"syntax error", "app.ts", "unexpected token"},
		},
		{
			name: "type error",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeTypeCheck,
				SourceFile: "handler.ts",
				LineNumber: 42,
				Message:    "Property 'user' does not exist on type 'Session'",
			},
			codeContext: "session.user",
			wantContain: []string{"type error", "handler.ts", "Property 'user'"},
		},
		{
			name: "linting error with tool",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeLinting,
				SourceFile: "utils.go",
				LineNumber: 15,
				Message:    "exported function should have comment",
				Tool:       "golangci-lint",
			},
			codeContext: "func ExportedFunc() {}",
			wantContain: []string{"linting", "utils.go", "golangci-lint"},
		},
		{
			name: "zero line number",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "main.go",
				LineNumber: 0,
				Message:    "compilation failed",
			},
			codeContext: "",
			wantContain: []string{"compilation error", "main.go:1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := gen.GeneratePrompt(tt.err, tt.codeContext)

			if prompt == "" {
				t.Error("expected non-empty prompt")
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(strings.ToLower(prompt), strings.ToLower(want)) {
					t.Errorf("prompt missing expected content: %q", want)
				}
			}
		})
	}
}

func TestAIFeedbackGenerator_FormatErrorsForAI(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	tests := []struct {
		name        string
		errors      []*domain.ErrorDetails
		wantEmpty   bool
		wantContain []string
	}{
		{
			name:      "empty errors",
			errors:    []*domain.ErrorDetails{},
			wantEmpty: true,
		},
		{
			name: "single error",
			errors: []*domain.ErrorDetails{
				{
					ErrorType:  domain.ErrorTypeImport,
					SourceFile: "main.go",
					LineNumber: 10,
					Message:    "undefined: fmt",
				},
			},
			wantContain: []string{"Error 1", "main.go", "import", "undefined: fmt"},
		},
		{
			name: "multiple errors",
			errors: []*domain.ErrorDetails{
				{
					ErrorType:  domain.ErrorTypeImport,
					SourceFile: "main.go",
					Message:    "error 1",
				},
				{
					ErrorType:  domain.ErrorTypeSyntax,
					SourceFile: "app.ts",
					Message:    "error 2",
				},
			},
			wantContain: []string{"Error 1", "Error 2", "main.go", "app.ts"},
		},
		{
			name: "error with suggestions",
			errors: []*domain.ErrorDetails{
				{
					ErrorType:   domain.ErrorTypeTypeCheck,
					SourceFile:  "handler.go",
					Message:     "type mismatch",
					Suggestions: []string{"Use type assertion", "Check interface"},
				},
			},
			wantContain: []string{"Suggestions", "type assertion", "Check interface"},
		},
		{
			name: "error with tool",
			errors: []*domain.ErrorDetails{
				{
					ErrorType:  domain.ErrorTypeLinting,
					SourceFile: "main.go",
					Message:    "lint error",
					Tool:       "eslint",
				},
			},
			wantContain: []string{"Tool", "eslint"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.FormatErrorsForAI(tt.errors)

			if tt.wantEmpty {
				if result != "" {
					t.Errorf("expected empty result, got: %s", result)
				}
				return
			}

			if result == "" {
				t.Error("expected non-empty result")
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("result missing expected content: %q", want)
				}
			}
		})
	}
}

func TestAIFeedbackGenerator_FormatErrorWithContext(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	tests := []struct {
		name        string
		err         *domain.ErrorDetails
		codeSnippet string
		wantContain []string
	}{
		{
			name: "full error with context",
			err: &domain.ErrorDetails{
				ErrorType:   domain.ErrorTypeTypeCheck,
				SourceFile:  "handler.ts",
				LineNumber:  42,
				Column:      10,
				Message:     "Property 'user' does not exist",
				Tool:        "tsc",
				Severity:    "error",
				Suggestions: []string{"Add user property"},
			},
			codeSnippet: "const user = session.user;",
			wantContain: []string{
				"handler.ts:42:10",
				"typecheck",
				"Property 'user'",
				"tsc",
				"error",
				"session.user",
				"Add user property",
			},
		},
		{
			name: "error without optional fields",
			err: &domain.ErrorDetails{
				ErrorType:  domain.ErrorTypeCompilation,
				SourceFile: "main.go",
				Message:    "compilation failed",
			},
			codeSnippet: "",
			wantContain: []string{"main.go", "compilation", "compilation failed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.FormatErrorWithContext(tt.err, tt.codeSnippet)

			if result == "" {
				t.Error("expected non-empty result")
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("result missing expected content: %q", want)
				}
			}
		})
	}
}

func TestAIFeedbackGenerator_GenerateBatchFixPrompt(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	tests := []struct {
		name         string
		errors       []*domain.ErrorDetails
		fileContents map[string]string
		wantContain  []string
	}{
		{
			name: "single file multiple errors",
			errors: []*domain.ErrorDetails{
				{ErrorType: domain.ErrorTypeImport, SourceFile: "main.go", LineNumber: 5, Message: "error 1"},
				{ErrorType: domain.ErrorTypeSyntax, SourceFile: "main.go", LineNumber: 10, Message: "error 2"},
			},
			fileContents: map[string]string{
				"main.go": "package main\n\nfunc main() {}",
			},
			wantContain: []string{"main.go", "error 1", "error 2", "package main"},
		},
		{
			name: "multiple files",
			errors: []*domain.ErrorDetails{
				{ErrorType: domain.ErrorTypeImport, SourceFile: "main.go", LineNumber: 5, Message: "go error"},
				{ErrorType: domain.ErrorTypeSyntax, SourceFile: "app.ts", LineNumber: 10, Message: "ts error"},
			},
			fileContents: map[string]string{
				"main.go": "package main",
				"app.ts":  "const x = 1;",
			},
			wantContain: []string{"main.go", "app.ts", "go error", "ts error"},
		},
		{
			name: "no file contents",
			errors: []*domain.ErrorDetails{
				{ErrorType: domain.ErrorTypeCompilation, SourceFile: "unknown.go", Message: "error"},
			},
			fileContents: map[string]string{},
			wantContain:  []string{"unknown.go", "error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.GenerateBatchFixPrompt(tt.errors, tt.fileContents)

			if result == "" {
				t.Error("expected non-empty result")
			}

			for _, want := range tt.wantContain {
				if !strings.Contains(result, want) {
					t.Errorf("result missing expected content: %q", want)
				}
			}
		})
	}
}

func TestAIFeedbackGenerator_SetTemplate(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	customTemplate := "Custom template: %s:%d - %s\n%s"
	gen.SetTemplate(domain.ErrorTypeImport, customTemplate)

	got := gen.GetTemplate(domain.ErrorTypeImport)
	if got != customTemplate {
		t.Errorf("GetTemplate() = %q, want %q", got, customTemplate)
	}
}

func TestAIFeedbackGenerator_GetTemplate(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	tests := []struct {
		name      string
		errorType domain.ErrorType
		wantEmpty bool
	}{
		{"import template", domain.ErrorTypeImport, false},
		{"syntax template", domain.ErrorTypeSyntax, false},
		{"typecheck template", domain.ErrorTypeTypeCheck, false},
		{"compilation template", domain.ErrorTypeCompilation, false},
		{"linting template", domain.ErrorTypeLinting, false},
		{"testing template", domain.ErrorTypeTesting, false},
		{"dependency template", domain.ErrorTypeDependency, false},
		{"logic template", domain.ErrorTypeLogic, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			template := gen.GetTemplate(tt.errorType)
			if tt.wantEmpty && template != "" {
				t.Errorf("expected empty template, got: %s", template)
			}
			if !tt.wantEmpty && template == "" {
				t.Error("expected non-empty template")
			}
		})
	}
}

func TestAIFeedbackGenerator_UnknownErrorType(t *testing.T) {
	gen := NewAIFeedbackGenerator()

	err := &domain.ErrorDetails{
		ErrorType:  "unknown_type",
		SourceFile: "test.go",
		LineNumber: 1,
		Message:    "unknown error",
	}

	// Should fall back to compilation template
	prompt := gen.GeneratePrompt(err, "code context")

	if prompt == "" {
		t.Error("expected non-empty prompt for unknown error type")
	}
	if !strings.Contains(prompt, "test.go") {
		t.Error("prompt should contain file name")
	}
}
