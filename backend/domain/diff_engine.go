package domain

import "context"

// DiffFormat определяет формат diff
type DiffFormat string

const (
	DiffFormatGit     DiffFormat = "git"
	DiffFormatUnified DiffFormat = "unified"
	DiffFormatJSON    DiffFormat = "json"
	DiffFormatHTML    DiffFormat = "html"
)

// DiffEntry представляет одну запись в diff
type DiffEntry struct {
	Path       string            `json:"path"`
	Operation  string            `json:"operation"` // "added", "modified", "deleted"
	OldContent string            `json:"oldContent,omitempty"`
	NewContent string            `json:"newContent,omitempty"`
	Hunks      []*DiffHunk       `json:"hunks,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// DiffHunk представляет блок изменений
type DiffHunk struct {
	OldStart int      `json:"oldStart"`
	OldCount int      `json:"oldCount"`
	NewStart int      `json:"newStart"`
	NewCount int      `json:"newCount"`
	Lines    []string `json:"lines"`
}

// DiffResult представляет результат генерации diff
type DiffResult struct {
	ID          string       `json:"id"`
	Format      DiffFormat   `json:"format"`
	Content     string       `json:"content"`
	Entries     []*DiffEntry `json:"entries"`
	Summary     *DiffSummary `json:"summary"`
	GeneratedAt string       `json:"generatedAt"`
}

// DiffSummary представляет сводку изменений
type DiffSummary struct {
	TotalFiles    int             `json:"totalFiles"`
	AddedFiles    int             `json:"addedFiles"`
	ModifiedFiles int             `json:"modifiedFiles"`
	DeletedFiles  int             `json:"deletedFiles"`
	TotalLines    int             `json:"totalLines"`
	AddedLines    int             `json:"addedLines"`
	RemovedLines  int             `json:"removedLines"`
	WhyView       *WhyView        `json:"whyView,omitempty"`
	Impact        *ImpactAnalysis `json:"impact,omitempty"`
	Risk          *RiskAssessment `json:"risk,omitempty"`
}

// WhyView представляет анализ причин изменений
type WhyView struct {
	Reason       string            `json:"reason"`
	TaskID       string            `json:"taskId"`
	StepID       string            `json:"stepId"`
	Confidence   float64           `json:"confidence"`
	RelatedFiles []string          `json:"relatedFiles"`
	Context      map[string]string `json:"context,omitempty"`
}

// ImpactAnalysis представляет анализ влияния изменений
type ImpactAnalysis struct {
	Level           string   `json:"level"` // "low", "medium", "high", "critical"
	AffectedAreas   []string `json:"affectedAreas"`
	Breaking        bool     `json:"breaking"`
	Performance     string   `json:"performance"`     // "improved", "degraded", "unchanged"
	Security        string   `json:"security"`        // "improved", "degraded", "unchanged"
	Maintainability string   `json:"maintainability"` // "improved", "degraded", "unchanged"
}

// RiskAssessment представляет оценку рисков
type RiskAssessment struct {
	Level        string   `json:"level"` // "low", "medium", "high", "critical"
	Risks        []string `json:"risks"`
	Mitigations  []string `json:"mitigations"`
	TestCoverage string   `json:"testCoverage"` // "adequate", "insufficient", "unknown"
	ReviewNeeded bool     `json:"reviewNeeded"`
}

// DiffEngine определяет интерфейс для генерации diff
type DiffEngine interface {
	// GenerateDiff генерирует diff между двумя состояниями
	GenerateDiff(ctx context.Context, beforePath, afterPath string, format DiffFormat) (*DiffResult, error)

	// GenerateDiffFromResults генерирует diff из результатов применения правок
	GenerateDiffFromResults(ctx context.Context, results []*ApplyResult, format DiffFormat) (*DiffResult, error)

	// GenerateDiffFromEdits генерирует diff из Edits JSON
	GenerateDiffFromEdits(ctx context.Context, edits *EditsJSON, format DiffFormat) (*DiffResult, error)

	// PublishDiff публикует diff
	PublishDiff(ctx context.Context, diff *DiffResult) error
}

// DiffPublisher определяет интерфейс для публикации diff
type DiffPublisher interface {
	// Publish публикует diff в указанном формате
	Publish(ctx context.Context, diff *DiffResult, target string) error

	// GetSupportedTargets возвращает поддерживаемые цели публикации
	GetSupportedTargets() []string

	// GetSupportedFormats возвращает поддерживаемые форматы
	GetSupportedFormats() []DiffFormat
}
