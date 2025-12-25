package domain

import "time"

// BudgetTypeLinesChanged represents budget type for lines changed
type BudgetTypeLinesChanged string

const (
	BudgetTypeLinesChangedDefault BudgetTypeLinesChanged = "lines_changed"
)

// ViolationReport represents a guardrail violation report
type ViolationReport struct {
	ID               string               `json:"id"`
	TaskID           string               `json:"taskId"`
	PolicyName       string               `json:"policyName"`
	Severity         string               `json:"severity"`
	Message          string               `json:"message"`
	FilePath         string               `json:"filePath,omitempty"`
	LineNumber       int                  `json:"lineNumber,omitempty"`
	Suggestion       string               `json:"suggestion,omitempty"`
	Violations       []GuardrailViolation `json:"violations"`
	BudgetViolations []BudgetViolation    `json:"budgetViolations"`
	GeneratedAt      time.Time            `json:"generatedAt"`
	CreatedAt        time.Time            `json:"createdAt"`
}

// SBOMToolsStatus represents status of SBOM tools
type SBOMToolsStatus struct {
	SyftAvailable  bool      `json:"syftAvailable"`
	GrypeAvailable bool      `json:"grypeAvailable"`
	SyftVersion    string    `json:"syftVersion,omitempty"`
	GrypeVersion   string    `json:"grypeVersion,omitempty"`
	LastChecked    time.Time `json:"lastChecked"`
}

// VulnerabilityStats represents vulnerability statistics
type VulnerabilityStats struct {
	Total    int `json:"total"`
	Critical int `json:"critical"`
	High     int `json:"high"`
	Medium   int `json:"medium"`
	Low      int `json:"low"`
	Unknown  int `json:"unknown"`
}

// LicenseStats represents license statistics
type LicenseStats struct {
	Total       int            `json:"total"`
	ByLicense   map[string]int `json:"byLicense"`
	Permissive  int            `json:"permissive"`
	Copyleft    int            `json:"copyleft"`
	Proprietary int            `json:"proprietary"`
	Unknown     int            `json:"unknown"`
}

// StaticAnalysisMetrics represents static analysis metrics
type StaticAnalysisMetrics struct {
	TotalLanguages   int                    `json:"totalLanguages"`
	FilesAnalyzed    int                    `json:"filesAnalyzed"`
	TotalIssues      int                    `json:"totalIssues"`
	IssuesBySeverity map[string]int         `json:"issuesBySeverity"`
	IssuesByCategory map[string]int         `json:"issuesByCategory"`
	IssuesByType     map[string]int         `json:"issuesByType"`
	IssuesByLevel    map[string]int         `json:"issuesByLevel"`
	Duration         time.Duration          `json:"duration"`
	ToolsUsed        []string               `json:"toolsUsed"`
	Coverage         StaticAnalysisCoverage `json:"coverage"`
	CriticalIssues   []*StaticAnalysisIssue `json:"criticalIssues"`
	Success          bool                   `json:"success"`
}

// StaticAnalysisCoverage represents static analysis coverage
type StaticAnalysisCoverage struct {
	LinesOfCode     int     `json:"linesOfCode"`
	LinesAnalyzed   int     `json:"linesAnalyzed"`
	CoveragePercent float64 `json:"coveragePercent"`
}

// StaticAnalysisIssue represents a single static analysis issue
type StaticAnalysisIssue struct {
	ID          string    `json:"id"`
	Rule        string    `json:"rule"`
	Severity    string    `json:"severity"`
	Category    string    `json:"category"`
	Message     string    `json:"message"`
	FilePath    string    `json:"filePath"`
	LineNumber  int       `json:"lineNumber"`
	ColumnStart int       `json:"columnStart,omitempty"`
	ColumnEnd   int       `json:"columnEnd,omitempty"`
	Tool        string    `json:"tool"`
	Suggestion  string    `json:"suggestion,omitempty"`
	CodeSnippet string    `json:"codeSnippet,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
