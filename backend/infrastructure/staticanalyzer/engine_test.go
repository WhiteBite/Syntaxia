package staticanalyzer

import (
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
