package staticanalyzer

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"testing"
)

func TestGenerateSummary(t *testing.T) {
	tests := []struct {
		name           string
		issues         []*domain.StaticIssue
		expectedTotal  int
		expectedErrors int
		expectedWarns  int
		expectedInfo   int
		expectedHints  int
		expectedFiles  int
	}{
		{
			name:           "empty issues",
			issues:         []*domain.StaticIssue{},
			expectedTotal:  0,
			expectedErrors: 0,
			expectedWarns:  0,
			expectedInfo:   0,
			expectedHints:  0,
			expectedFiles:  0,
		},
		{
			name: "single error",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "error", Category: "lint"},
			},
			expectedTotal:  1,
			expectedErrors: 1,
			expectedWarns:  0,
			expectedInfo:   0,
			expectedHints:  0,
			expectedFiles:  1,
		},
		{
			name: "single warning",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "warning", Category: "style"},
			},
			expectedTotal:  1,
			expectedErrors: 0,
			expectedWarns:  1,
			expectedInfo:   0,
			expectedHints:  0,
			expectedFiles:  1,
		},
		{
			name: "single info",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "info", Category: "documentation"},
			},
			expectedTotal:  1,
			expectedErrors: 0,
			expectedWarns:  0,
			expectedInfo:   1,
			expectedHints:  0,
			expectedFiles:  1,
		},
		{
			name: "single hint",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "hint", Category: "style"},
			},
			expectedTotal:  1,
			expectedErrors: 0,
			expectedWarns:  0,
			expectedInfo:   0,
			expectedHints:  1,
			expectedFiles:  1,
		},
		{
			name: "mixed severities",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "error", Category: "lint"},
				{File: "file1.go", Severity: "warning", Category: "style"},
				{File: "file2.go", Severity: "info", Category: "documentation"},
				{File: "file3.go", Severity: "hint", Category: "style"},
			},
			expectedTotal:  4,
			expectedErrors: 1,
			expectedWarns:  1,
			expectedInfo:   1,
			expectedHints:  1,
			expectedFiles:  3,
		},
		{
			name: "multiple issues same file",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "error", Category: "lint"},
				{File: "file1.go", Severity: "error", Category: "lint"},
				{File: "file1.go", Severity: "warning", Category: "style"},
			},
			expectedTotal:  3,
			expectedErrors: 2,
			expectedWarns:  1,
			expectedInfo:   0,
			expectedHints:  0,
			expectedFiles:  1,
		},
		{
			name: "unknown severity",
			issues: []*domain.StaticIssue{
				{File: "file1.go", Severity: "unknown", Category: "other"},
			},
			expectedTotal:  1,
			expectedErrors: 0,
			expectedWarns:  0,
			expectedInfo:   0,
			expectedHints:  0,
			expectedFiles:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := generateSummary(tt.issues)

			if summary.TotalIssues != tt.expectedTotal {
				t.Errorf("TotalIssues = %d, want %d", summary.TotalIssues, tt.expectedTotal)
			}
			if summary.ErrorCount != tt.expectedErrors {
				t.Errorf("ErrorCount = %d, want %d", summary.ErrorCount, tt.expectedErrors)
			}
			if summary.WarningCount != tt.expectedWarns {
				t.Errorf("WarningCount = %d, want %d", summary.WarningCount, tt.expectedWarns)
			}
			if summary.InfoCount != tt.expectedInfo {
				t.Errorf("InfoCount = %d, want %d", summary.InfoCount, tt.expectedInfo)
			}
			if summary.HintCount != tt.expectedHints {
				t.Errorf("HintCount = %d, want %d", summary.HintCount, tt.expectedHints)
			}
			if summary.FilesWithIssues != tt.expectedFiles {
				t.Errorf("FilesWithIssues = %d, want %d", summary.FilesWithIssues, tt.expectedFiles)
			}
		})
	}
}

func TestGenerateSummary_CategoryBreakdown(t *testing.T) {
	issues := []*domain.StaticIssue{
		{File: "file1.go", Severity: "error", Category: "lint"},
		{File: "file2.go", Severity: "warning", Category: "style"},
		{File: "file3.go", Severity: "warning", Category: "style"},
		{File: "file4.go", Severity: "info", Category: "documentation"},
	}

	summary := generateSummary(issues)

	if summary.CategoryBreakdown["lint"] != 1 {
		t.Errorf("CategoryBreakdown[lint] = %d, want 1", summary.CategoryBreakdown["lint"])
	}
	if summary.CategoryBreakdown["style"] != 2 {
		t.Errorf("CategoryBreakdown[style] = %d, want 2", summary.CategoryBreakdown["style"])
	}
	if summary.CategoryBreakdown["documentation"] != 1 {
		t.Errorf("CategoryBreakdown[documentation] = %d, want 1", summary.CategoryBreakdown["documentation"])
	}
}

func TestGenerateSummary_SeverityBreakdown(t *testing.T) {
	issues := []*domain.StaticIssue{
		{File: "file1.go", Severity: "error", Category: "lint"},
		{File: "file2.go", Severity: "error", Category: "lint"},
		{File: "file3.go", Severity: "warning", Category: "style"},
	}

	summary := generateSummary(issues)

	if summary.SeverityBreakdown["error"] != 2 {
		t.Errorf("SeverityBreakdown[error] = %d, want 2", summary.SeverityBreakdown["error"])
	}
	if summary.SeverityBreakdown["warning"] != 1 {
		t.Errorf("SeverityBreakdown[warning] = %d, want 1", summary.SeverityBreakdown["warning"])
	}
}

func TestNewStaticAnalyzerEngine(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	if engine == nil {
		t.Fatal("NewStaticAnalyzerEngine() returned nil")
	}

	if engine.analyzers == nil {
		t.Error("engine.analyzers is nil")
	}

	if engine.languageMap == nil {
		t.Error("engine.languageMap is nil")
	}

	// Check language mappings
	expectedMappings := map[string]domain.StaticAnalyzerType{
		"go":         domain.StaticAnalyzerTypeStaticcheck,
		"typescript": domain.StaticAnalyzerTypeESLint,
		"ts":         domain.StaticAnalyzerTypeESLint,
		"javascript": domain.StaticAnalyzerTypeESLint,
		"js":         domain.StaticAnalyzerTypeESLint,
		"java":       domain.StaticAnalyzerTypeErrorProne,
		"python":     domain.StaticAnalyzerTypeRuff,
		"py":         domain.StaticAnalyzerTypeRuff,
		"rust":       domain.StaticAnalyzerTypeClippy,
		"rs":         domain.StaticAnalyzerTypeClippy,
		"kotlin":     domain.StaticAnalyzerTypeKtlint,
		"kt":         domain.StaticAnalyzerTypeKtlint,
	}

	for lang, expectedType := range expectedMappings {
		if engine.languageMap[lang] != expectedType {
			t.Errorf("languageMap[%s] = %v, want %v", lang, engine.languageMap[lang], expectedType)
		}
	}
}

func TestStaticAnalyzerEngine_RegisterAnalyzer(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})
	analyzer := NewClippyAnalyzer(&mockLogger{})

	engine.RegisterAnalyzer(analyzer)

	registered, exists := engine.analyzers[string(domain.StaticAnalyzerTypeClippy)]
	if !exists {
		t.Error("Analyzer was not registered")
	}
	if registered != analyzer {
		t.Error("Registered analyzer does not match")
	}
}

func TestStaticAnalyzerEngine_GetSupportedAnalyzers(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	// Register some analyzers
	engine.RegisterAnalyzer(NewClippyAnalyzer(&mockLogger{}))
	engine.RegisterAnalyzer(NewKtlintAnalyzer(&mockLogger{}))

	analyzers := engine.GetSupportedAnalyzers()

	if len(analyzers) != 2 {
		t.Errorf("GetSupportedAnalyzers() returned %d analyzers, want 2", len(analyzers))
	}
}

func TestStaticAnalyzerEngine_GetAnalyzerForLanguage(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})
	clippyAnalyzer := NewClippyAnalyzer(&mockLogger{})
	engine.RegisterAnalyzer(clippyAnalyzer)

	tests := []struct {
		name     string
		language string
		wantErr  bool
	}{
		{
			name:     "rust language",
			language: "rust",
			wantErr:  false,
		},
		{
			name:     "rs extension",
			language: "rs",
			wantErr:  false,
		},
		{
			name:     "uppercase RUST",
			language: "RUST",
			wantErr:  false,
		},
		{
			name:     "unmapped language",
			language: "cobol",
			wantErr:  true,
		},
		{
			name:     "mapped but not registered",
			language: "go",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analyzer, err := engine.GetAnalyzerForLanguage(tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalyzerForLanguage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && analyzer == nil {
				t.Error("GetAnalyzerForLanguage() returned nil analyzer")
			}
		})
	}
}

func TestStaticAnalyzerEngine_GenerateReport(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.5,
			Issues: []*domain.StaticIssue{
				{File: "main.rs", Line: 10, Severity: "error", Message: "Error 1"},
				{File: "lib.rs", Line: 20, Severity: "warning", Message: "Warning 1"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  2,
				ErrorCount:   1,
				WarningCount: 1,
			},
		},
		"kotlin": {
			Success:     true,
			Language:    "kotlin",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeKtlint,
			Duration:    0.5,
			Issues: []*domain.StaticIssue{
				{File: "Main.kt", Line: 5, Severity: "warning", Message: "Warning 2"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  1,
				WarningCount: 1,
			},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if report == nil {
		t.Fatal("GenerateReport() returned nil")
	}

	if report.ProjectPath != "/project" {
		t.Errorf("report.ProjectPath = %q, want %q", report.ProjectPath, "/project")
	}

	if report.Summary.TotalIssues != 3 {
		t.Errorf("report.Summary.TotalIssues = %d, want 3", report.Summary.TotalIssues)
	}

	if report.Summary.TotalErrors != 1 {
		t.Errorf("report.Summary.TotalErrors = %d, want 1", report.Summary.TotalErrors)
	}

	if report.Summary.TotalWarnings != 2 {
		t.Errorf("report.Summary.TotalWarnings = %d, want 2", report.Summary.TotalWarnings)
	}

	if len(report.Summary.CriticalIssues) != 1 {
		t.Errorf("len(report.Summary.CriticalIssues) = %d, want 1", len(report.Summary.CriticalIssues))
	}

	if report.Summary.Success {
		t.Error("report.Summary.Success = true, want false (has errors)")
	}
}

func TestStaticAnalyzerEngine_GenerateReport_NoErrors(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.0,
			Issues:      []*domain.StaticIssue{},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues: 0,
			},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if !report.Summary.Success {
		t.Error("report.Summary.Success = false, want true (no errors)")
	}

	if len(report.Recommendations) == 0 {
		t.Error("report.Recommendations is empty, want at least one recommendation")
	}
}

func TestStaticAnalyzerEngine_GenerateReport_WithWarnings(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.0,
			Issues: []*domain.StaticIssue{
				{File: "main.rs", Line: 10, Severity: "warning", Message: "Warning"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  1,
				WarningCount: 1,
			},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if report.Summary.TotalWarnings != 1 {
		t.Errorf("report.Summary.TotalWarnings = %d, want 1", report.Summary.TotalWarnings)
	}

	// Should have recommendation about warnings
	hasWarningRecommendation := false
	for _, rec := range report.Recommendations {
		if contains(rec, "warning") {
			hasWarningRecommendation = true
			break
		}
	}
	if !hasWarningRecommendation {
		t.Error("Expected recommendation about warnings")
	}
}

func TestStaticAnalyzerEngine_GenerateReport_FailedResults(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     false,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    0.5,
			Error:       "Analysis failed",
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if report == nil {
		t.Fatal("GenerateReport() returned nil")
	}

	// Failed results should not be counted in languages analyzed
	if len(report.Summary.LanguagesAnalyzed) != 0 {
		t.Errorf("len(report.Summary.LanguagesAnalyzed) = %d, want 0", len(report.Summary.LanguagesAnalyzed))
	}
}

func TestStaticAnalyzerEngine_GenerateReport_MixedResults(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.0,
			Issues: []*domain.StaticIssue{
				{File: "main.rs", Line: 10, Severity: "error", Message: "Error 1"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  1,
				ErrorCount:   1,
				WarningCount: 0,
			},
		},
		"kotlin": {
			Success:     false,
			Language:    "kotlin",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeKtlint,
			Duration:    0.5,
			Error:       "Ktlint not found",
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if len(report.Summary.LanguagesAnalyzed) != 1 {
		t.Errorf("len(report.Summary.LanguagesAnalyzed) = %d, want 1", len(report.Summary.LanguagesAnalyzed))
	}

	if report.Summary.TotalErrors != 1 {
		t.Errorf("report.Summary.TotalErrors = %d, want 1", report.Summary.TotalErrors)
	}
}

func TestStaticAnalyzerEngine_GenerateReport_WithCriticalIssues(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"go": {
			Success:     true,
			Language:    "go",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
			Duration:    2.0,
			Issues: []*domain.StaticIssue{
				{File: "main.go", Line: 10, Severity: "error", Message: "Critical error 1"},
				{File: "utils.go", Line: 20, Severity: "error", Message: "Critical error 2"},
				{File: "helper.go", Line: 30, Severity: "warning", Message: "Warning"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  3,
				ErrorCount:   2,
				WarningCount: 1,
			},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if len(report.Summary.CriticalIssues) != 2 {
		t.Errorf("len(report.Summary.CriticalIssues) = %d, want 2", len(report.Summary.CriticalIssues))
	}

	// Should have recommendations about critical issues
	hasCriticalRecommendation := false
	for _, rec := range report.Recommendations {
		if contains(rec, "critical") || contains(rec, "error") {
			hasCriticalRecommendation = true
			break
		}
	}
	if !hasCriticalRecommendation {
		t.Error("Expected recommendation about critical issues")
	}
}

func TestStaticAnalyzerEngine_GenerateReport_Timestamp(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{}

	report := engine.GenerateReport(results, "/project")

	if report.Timestamp == "" {
		t.Error("report.Timestamp is empty")
	}
}

func TestStaticAnalyzerEngine_AnalyzeProject_NoAnalyzers(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})
	// Don't register any analyzers

	results, err := engine.AnalyzeProject(nil, "/project", []string{"go", "rust"})
	if err != nil {
		t.Errorf("AnalyzeProject() error = %v", err)
	}

	// Should return empty results since no analyzers are registered
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}

func TestStaticAnalyzerEngine_AnalyzeProject_UnknownLanguage(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results, err := engine.AnalyzeProject(nil, "/project", []string{"cobol", "fortran"})
	if err != nil {
		t.Errorf("AnalyzeProject() error = %v", err)
	}

	// Should return empty results for unknown languages
	if len(results) != 0 {
		t.Errorf("len(results) = %d, want 0", len(results))
	}
}


func TestStaticAnalyzerEngine_AnalyzeFile_NoAnalyzer(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})
	// Don't register any analyzers

	config := &domain.StaticAnalyzerConfig{
		Language:    "go",
		ProjectPath: "/project",
	}

	_, err := engine.AnalyzeFile(nil, "/project/main.go", config)
	if err == nil {
		t.Error("AnalyzeFile() should return error when no analyzer is registered")
	}
}

func TestStaticAnalyzerEngine_AnalyzeFile_UnknownLanguage(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	config := &domain.StaticAnalyzerConfig{
		Language:    "cobol",
		ProjectPath: "/project",
	}

	_, err := engine.AnalyzeFile(nil, "/project/main.cob", config)
	if err == nil {
		t.Error("AnalyzeFile() should return error for unknown language")
	}
}

func TestStaticAnalyzerEngine_GetAnalyzerForLanguage_CaseInsensitive(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})
	clippyAnalyzer := NewClippyAnalyzer(&mockLogger{})
	engine.RegisterAnalyzer(clippyAnalyzer)

	tests := []struct {
		language string
		wantErr  bool
	}{
		{"rust", false},
		{"RUST", false},
		{"Rust", false},
		{"rUsT", false},
		{"rs", false},
		{"RS", false},
		{"Rs", false},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			_, err := engine.GetAnalyzerForLanguage(tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAnalyzerForLanguage(%q) error = %v, wantErr %v", tt.language, err, tt.wantErr)
			}
		})
	}
}

func TestStaticAnalyzerEngine_GenerateReport_EmptyResults(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{}

	report := engine.GenerateReport(results, "/project")

	if report == nil {
		t.Fatal("GenerateReport() returned nil")
	}

	if report.Summary.TotalIssues != 0 {
		t.Errorf("report.Summary.TotalIssues = %d, want 0", report.Summary.TotalIssues)
	}

	if !report.Summary.Success {
		t.Error("report.Summary.Success = false, want true (no errors)")
	}
}

func TestStaticAnalyzerEngine_GenerateReport_NilSummary(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.0,
			Issues:      []*domain.StaticIssue{},
			Summary:     nil, // nil summary
		},
	}

	report := engine.GenerateReport(results, "/project")

	if report == nil {
		t.Fatal("GenerateReport() returned nil")
	}

	// Should handle nil summary gracefully
	if report.Summary.TotalIssues != 0 {
		t.Errorf("report.Summary.TotalIssues = %d, want 0", report.Summary.TotalIssues)
	}
}

func TestStaticAnalyzerEngine_LanguageMap_AllMappings(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	expectedMappings := map[string]domain.StaticAnalyzerType{
		"go":         domain.StaticAnalyzerTypeStaticcheck,
		"typescript": domain.StaticAnalyzerTypeESLint,
		"ts":         domain.StaticAnalyzerTypeESLint,
		"javascript": domain.StaticAnalyzerTypeESLint,
		"js":         domain.StaticAnalyzerTypeESLint,
		"java":       domain.StaticAnalyzerTypeErrorProne,
		"python":     domain.StaticAnalyzerTypeRuff,
		"py":         domain.StaticAnalyzerTypeRuff,
		"c":          domain.StaticAnalyzerTypeClangTidy,
		"cpp":        domain.StaticAnalyzerTypeClangTidy,
		"cc":         domain.StaticAnalyzerTypeClangTidy,
		"rust":       domain.StaticAnalyzerTypeClippy,
		"rs":         domain.StaticAnalyzerTypeClippy,
		"kotlin":     domain.StaticAnalyzerTypeKtlint,
		"kt":         domain.StaticAnalyzerTypeKtlint,
	}

	for lang, expectedType := range expectedMappings {
		t.Run(lang, func(t *testing.T) {
			if engine.languageMap[lang] != expectedType {
				t.Errorf("languageMap[%q] = %v, want %v", lang, engine.languageMap[lang], expectedType)
			}
		})
	}
}


// mockAnalyzer implements domain.StaticAnalyzer for testing
type mockAnalyzer struct {
	analyzerType domain.StaticAnalyzerType
	languages    []string
	analyzeFunc  func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error)
	validateFunc func(config *domain.StaticAnalyzerConfig) error
}

func (m *mockAnalyzer) Analyze(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
	if m.analyzeFunc != nil {
		return m.analyzeFunc(ctx, config)
	}
	return &domain.StaticAnalysisResult{
		Success:     true,
		Language:    config.Language,
		ProjectPath: config.ProjectPath,
		Analyzer:    m.analyzerType,
		Issues:      []*domain.StaticIssue{},
		Summary:     &domain.StaticAnalysisSummary{},
	}, nil
}

func (m *mockAnalyzer) GetSupportedLanguages() []string {
	return m.languages
}

func (m *mockAnalyzer) GetAnalyzerType() domain.StaticAnalyzerType {
	return m.analyzerType
}

func (m *mockAnalyzer) ValidateConfig(config *domain.StaticAnalyzerConfig) error {
	if m.validateFunc != nil {
		return m.validateFunc(config)
	}
	return nil
}

func TestStaticAnalyzerEngine_AnalyzeProject_WithMockAnalyzer(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return &domain.StaticAnalysisResult{
				Success:     true,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
				Issues: []*domain.StaticIssue{
					{File: "main.go", Line: 10, Severity: "warning", Message: "Test warning"},
				},
				Summary: &domain.StaticAnalysisSummary{
					TotalIssues:  1,
					WarningCount: 1,
				},
			}, nil
		},
	}

	engine.RegisterAnalyzer(analyzer)

	results, err := engine.AnalyzeProject(nil, "/project", []string{"go"})
	if err != nil {
		t.Errorf("AnalyzeProject() error = %v", err)
		return
	}

	if len(results) != 1 {
		t.Errorf("len(results) = %d, want 1", len(results))
		return
	}

	if results["go"] == nil {
		t.Error("results[\"go\"] is nil")
		return
	}

	if len(results["go"].Issues) != 1 {
		t.Errorf("len(results[\"go\"].Issues) = %d, want 1", len(results["go"].Issues))
	}
}

func TestStaticAnalyzerEngine_AnalyzeProject_AnalyzerError(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return nil, fmt.Errorf("analysis failed")
		},
	}

	engine.RegisterAnalyzer(analyzer)

	results, err := engine.AnalyzeProject(nil, "/project", []string{"go"})
	if err != nil {
		t.Errorf("AnalyzeProject() error = %v", err)
		return
	}

	// Should return result with error message
	if results["go"] == nil {
		t.Error("results[\"go\"] is nil")
		return
	}

	if results["go"].Success {
		t.Error("results[\"go\"].Success = true, want false")
	}

	if results["go"].Error == "" {
		t.Error("results[\"go\"].Error is empty")
	}
}

func TestStaticAnalyzerEngine_AnalyzeFile_WithMockAnalyzer(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return &domain.StaticAnalysisResult{
				Success:     true,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
				Issues: []*domain.StaticIssue{
					{File: "main.go", Line: 10, Severity: "error", Message: "Test error"},
				},
				Summary: &domain.StaticAnalysisSummary{
					TotalIssues: 1,
					ErrorCount:  1,
				},
			}, nil
		},
	}

	engine.RegisterAnalyzer(analyzer)

	config := &domain.StaticAnalyzerConfig{
		Language:    "go",
		ProjectPath: "/project",
		Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
	}

	result, err := engine.AnalyzeFile(nil, "/project/main.go", config)
	if err != nil {
		t.Errorf("AnalyzeFile() error = %v", err)
		return
	}

	if result == nil {
		t.Error("AnalyzeFile() returned nil")
		return
	}

	if len(result.Issues) != 1 {
		t.Errorf("len(result.Issues) = %d, want 1", len(result.Issues))
	}
}

func TestStaticAnalyzerEngine_AnalyzeFile_ValidationError(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		validateFunc: func(config *domain.StaticAnalyzerConfig) error {
			return fmt.Errorf("validation failed")
		},
	}

	engine.RegisterAnalyzer(analyzer)

	config := &domain.StaticAnalyzerConfig{
		Language:    "go",
		ProjectPath: "/project",
		Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
	}

	_, err := engine.AnalyzeFile(nil, "/project/main.go", config)
	if err == nil {
		t.Error("AnalyzeFile() should return error when validation fails")
	}
}

func TestStaticAnalyzerEngine_AnalyzeFile_AnalyzerError(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return nil, fmt.Errorf("analysis failed")
		},
	}

	engine.RegisterAnalyzer(analyzer)

	config := &domain.StaticAnalyzerConfig{
		Language:    "go",
		ProjectPath: "/project",
		Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
	}

	result, err := engine.AnalyzeFile(nil, "/project/main.go", config)
	// Should return result with error, not error
	if result == nil {
		t.Error("AnalyzeFile() should return result even on error")
	}

	if result != nil && result.Success {
		t.Error("result.Success = true, want false")
	}

	if result != nil && result.Error == "" {
		t.Error("result.Error is empty")
	}

	// err can be nil since we return result with error message
	_ = err
}

func TestStaticAnalyzerEngine_AnalyzeFile_WithDuration(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	analyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return &domain.StaticAnalysisResult{
				Success:     true,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
				Issues:      []*domain.StaticIssue{},
				Summary:     &domain.StaticAnalysisSummary{},
			}, nil
		},
	}

	engine.RegisterAnalyzer(analyzer)

	config := &domain.StaticAnalyzerConfig{
		Language:    "go",
		ProjectPath: "/project",
		Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
	}

	result, err := engine.AnalyzeFile(context.Background(), "/project/main.go", config)
	if err != nil {
		t.Errorf("AnalyzeFile() error = %v", err)
		return
	}

	if result == nil {
		t.Error("AnalyzeFile() returned nil")
		return
	}

	// Duration should be set
	if result.Duration < 0 {
		t.Errorf("result.Duration = %v, want >= 0", result.Duration)
	}
}

func TestStaticAnalyzerEngine_GenerateReport_TotalDuration(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.5,
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		},
		"kotlin": {
			Success:     true,
			Language:    "kotlin",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeKtlint,
			Duration:    2.5,
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if report.TotalDuration != 4.0 {
		t.Errorf("report.TotalDuration = %v, want 4.0", report.TotalDuration)
	}
}

func TestStaticAnalyzerEngine_GenerateReport_AnalyzersUsed(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"rust": {
			Success:     true,
			Language:    "rust",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeClippy,
			Duration:    1.0,
			Issues:      []*domain.StaticIssue{},
			Summary:     &domain.StaticAnalysisSummary{},
		},
	}

	report := engine.GenerateReport(results, "/project")

	if len(report.Summary.AnalyzersUsed) != 1 {
		t.Errorf("len(report.Summary.AnalyzersUsed) = %d, want 1", len(report.Summary.AnalyzersUsed))
	}
}

func TestStaticAnalyzerEngine_GenerateRecommendations_ErrorsAndWarnings(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	results := map[string]*domain.StaticAnalysisResult{
		"go": {
			Success:     true,
			Language:    "go",
			ProjectPath: "/project",
			Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
			Duration:    1.0,
			Issues: []*domain.StaticIssue{
				{File: "main.go", Line: 10, Severity: "error", Message: "Critical error"},
				{File: "utils.go", Line: 20, Severity: "warning", Message: "Warning message"},
			},
			Summary: &domain.StaticAnalysisSummary{
				TotalIssues:  2,
				ErrorCount:   1,
				WarningCount: 1,
			},
		},
	}

	report := engine.GenerateReport(results, "/project")

	// Should have recommendations for both errors and warnings
	hasErrorRec := false
	hasWarningRec := false
	for _, rec := range report.Recommendations {
		if contains(rec, "error") || contains(rec, "critical") {
			hasErrorRec = true
		}
		if contains(rec, "warning") {
			hasWarningRec = true
		}
	}

	if !hasErrorRec {
		t.Error("Expected recommendation about errors")
	}
	if !hasWarningRec {
		t.Error("Expected recommendation about warnings")
	}
}

func TestStaticAnalyzerEngine_AnalyzeProject_MultipleLanguages(t *testing.T) {
	engine := NewStaticAnalyzerEngine(&mockLogger{})

	goAnalyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeStaticcheck,
		languages:    []string{"go"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return &domain.StaticAnalysisResult{
				Success:     true,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    domain.StaticAnalyzerTypeStaticcheck,
				Issues:      []*domain.StaticIssue{{File: "main.go", Line: 1, Severity: "warning", Message: "Go warning"}},
				Summary:     &domain.StaticAnalysisSummary{TotalIssues: 1, WarningCount: 1},
			}, nil
		},
	}

	rustAnalyzer := &mockAnalyzer{
		analyzerType: domain.StaticAnalyzerTypeClippy,
		languages:    []string{"rust"},
		analyzeFunc: func(ctx context.Context, config *domain.StaticAnalyzerConfig) (*domain.StaticAnalysisResult, error) {
			return &domain.StaticAnalysisResult{
				Success:     true,
				Language:    config.Language,
				ProjectPath: config.ProjectPath,
				Analyzer:    domain.StaticAnalyzerTypeClippy,
				Issues:      []*domain.StaticIssue{{File: "main.rs", Line: 1, Severity: "error", Message: "Rust error"}},
				Summary:     &domain.StaticAnalysisSummary{TotalIssues: 1, ErrorCount: 1},
			}, nil
		},
	}

	engine.RegisterAnalyzer(goAnalyzer)
	engine.RegisterAnalyzer(rustAnalyzer)

	results, err := engine.AnalyzeProject(context.Background(), "/project", []string{"go", "rust"})
	if err != nil {
		t.Errorf("AnalyzeProject() error = %v", err)
		return
	}

	if len(results) != 2 {
		t.Errorf("len(results) = %d, want 2", len(results))
	}

	if results["go"] == nil {
		t.Error("results[\"go\"] is nil")
	}
	if results["rust"] == nil {
		t.Error("results[\"rust\"] is nil")
	}
}
