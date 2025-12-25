package analyzers

import (
	"testing"
)

func TestNewAnalyzerRegistry(t *testing.T) {
	registry := NewAnalyzerRegistry()
	if registry == nil {
		t.Fatal("NewAnalyzerRegistry returned nil")
	}

	langs := registry.SupportedLanguages()
	if len(langs) == 0 {
		t.Error("no languages registered")
	}

	exts := registry.SupportedExtensions()
	if len(exts) == 0 {
		t.Error("no extensions registered")
	}
}

func TestAnalyzerRegistry_GetAnalyzer(t *testing.T) {
	registry := NewAnalyzerRegistry()

	tests := []struct {
		filePath   string
		expectLang string
		expectNil  bool
	}{
		{"main.go", "go", false},
		{"app.ts", "typescript", false},
		{"app.tsx", "typescript", false},
		{"script.js", "javascript", false},
		{"Service.java", "java", false},
		{"Service.kt", "kotlin", false},
		{"unknown.xyz", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			analyzer := registry.GetAnalyzer(tt.filePath)
			if tt.expectNil {
				if analyzer != nil {
					t.Errorf("expected nil for %s", tt.filePath)
				}
				return
			}
			if analyzer == nil {
				t.Fatalf("expected analyzer for %s", tt.filePath)
			}
			if analyzer.Language() != tt.expectLang {
				t.Errorf("expected %q, got %q", tt.expectLang, analyzer.Language())
			}
		})
	}
}

func TestAnalyzerRegistry_GetAnalyzerByLanguage(t *testing.T) {
	registry := NewAnalyzerRegistry()

	tests := []struct {
		lang      string
		expectNil bool
	}{
		{"go", false},
		{"typescript", false},
		{"javascript", false},
		{"java", false},
		{"kotlin", false},
		{"unknown", true},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			analyzer := registry.GetAnalyzerByLanguage(tt.lang)
			if tt.expectNil && analyzer != nil {
				t.Errorf("expected nil for %s", tt.lang)
			}
			if !tt.expectNil && analyzer == nil {
				t.Errorf("expected analyzer for %s", tt.lang)
			}
		})
	}
}

func TestAnalyzerRegistry_SupportedLanguages(t *testing.T) {
	registry := NewAnalyzerRegistry()
	langs := registry.SupportedLanguages()

	expected := map[string]bool{"go": true, "typescript": true, "javascript": true, "java": true, "kotlin": true}
	for _, l := range langs {
		if !expected[l] {
			t.Logf("extra language: %s", l)
		}
	}

	langMap := make(map[string]bool)
	for _, l := range langs {
		langMap[l] = true
	}
	for exp := range expected {
		if !langMap[exp] {
			t.Errorf("expected language %q", exp)
		}
	}
}

func TestAnalyzerRegistry_SupportedExtensions(t *testing.T) {
	registry := NewAnalyzerRegistry()
	exts := registry.SupportedExtensions()

	expected := []string{".go", ".ts", ".tsx", ".js", ".java", ".kt"}
	extMap := make(map[string]bool)
	for _, e := range exts {
		extMap[e] = true
	}

	for _, exp := range expected {
		if !extMap[exp] {
			t.Errorf("expected extension %q", exp)
		}
	}
}
