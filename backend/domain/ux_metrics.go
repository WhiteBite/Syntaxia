package domain

import (
	"context"
	"time"
)

// UXReport представляет UX отчёт
type UXReport struct {
	ID          string
	Type        UXReportType
	Title       string
	Description string
	Content     interface{}
	CreatedAt   time.Time
	Metadata    map[string]interface{}
}

// UXReportType тип UX отчёта
type UXReportType string

const (
	UXReportTypeDerivedDiff UXReportType = "derived_diff"
	UXReportTypeWhyView     UXReportType = "why_view"
	UXReportTypeTimeToGreen UXReportType = "time_to_green"
	UXReportTypeMetrics     UXReportType = "metrics"
	UXReportTypePerformance UXReportType = "performance"
)

// WhyViewReport отчёт "почему эти файлы"
type WhyViewReport struct {
	TaskID      string
	Files       []FileReason
	Context     string
	Explanation string
	Confidence  float64
	Suggestions []string
}

// FileReason объяснение, почему файл был изменён
type FileReason struct {
	FilePath     string
	Reason       string
	Impact       string
	Confidence   float64
	RelatedFiles []string
	Context      map[string]interface{}
	Category     string   `json:"category"`
	Importance   string   `json:"importance"`
	Suggestions  []string `json:"suggestions"`
}

// TimeToGreenMetrics метрики времени до "зелёного" статуса
type TimeToGreenMetrics struct {
	TaskID             string
	StartTime          time.Time
	EndTime            time.Time
	Duration           time.Duration
	Attempts           int
	RepairAttempts     int
	BuildTime          time.Duration
	TestTime           time.Duration
	StaticAnalysisTime time.Duration
	TotalTime          time.Duration
	Success            bool
	Bottlenecks        []Bottleneck
}

// Bottleneck узкое место в процессе
type Bottleneck struct {
	Type        string
	Description string
	Duration    time.Duration
	Impact      string
	Suggestions []string
}

// PerformanceMetrics метрики производительности
type PerformanceMetrics struct {
	TaskID         string
	MemoryUsage    int64
	CPUUsage       float64
	DiskIO         int64
	NetworkIO      int64
	FileOperations int
	APIRequests    int
	CacheHits      int
	CacheMisses    int
	Timestamps     []string // ISO 8601 formatted timestamps (Wails doesn't support time.Time in TS bindings)
	Values         []float64
}

// DerivedDiffReport отчёт о derived diff
type DerivedDiffReport struct {
	TaskID       string
	OriginalDiff string
	DerivedDiff  string
	Changes      []DiffChange
	Summary      *DiffSummary
	Impact       DiffImpact
}

// DiffChange изменение в diff
type DiffChange struct {
	Type       string
	FilePath   string
	LineNumber int
	OldContent string
	NewContent string
	Reason     string
	Confidence float64
}

// DiffImpact влияние изменений
type DiffImpact struct {
	RiskLevel         string
	AffectedTests     []string
	BreakingChanges   []string
	PerformanceImpact string
	SecurityImpact    string
}

// UXMetricsService интерфейс для сервиса UX метрик
type UXMetricsService interface {
	// GenerateWhyViewReport генерирует отчёт "почему эти файлы"
	GenerateWhyViewReport(taskID string, files []string, taskContext map[string]interface{}) (*WhyViewReport, error)

	// GenerateTimeToGreenMetrics генерирует метрики time_to_green
	GenerateTimeToGreenMetrics(taskID string) (*TimeToGreenMetrics, error)

	// GenerateDerivedDiffReport генерирует отчёт о derived diff
	GenerateDerivedDiffReport(taskID string, originalDiff, derivedDiff string) (*DerivedDiffReport, error)

	// GeneratePerformanceMetrics генерирует метрики производительности
	GeneratePerformanceMetrics(taskID string) (*PerformanceMetrics, error)

	// GetUXReport возвращает UX отчёт
	GetUXReport(reportID string) (*UXReport, error)

	// SaveUXReport сохраняет UX отчёт
	SaveUXReport(report *UXReport) error

	// GetUXReports возвращает все UX отчёты
	GetUXReports(reportType UXReportType) ([]*UXReport, error)

	// DeleteUXReport удаляет UX отчёт
	DeleteUXReport(reportID string) error

	// GetMetricsSummary возвращает сводку метрик
	GetMetricsSummary() (map[string]interface{}, error)

	// ListReports получает список отчетов
	ListReports(ctx context.Context, reportType string) ([]*GenericReport, error)

	// GetReport получает конкретный отчет
	GetReport(ctx context.Context, reportID string) (*GenericReport, error)
}

// UXReportRepository интерфейс для работы с UX отчетами
type UXReportRepository interface {
	// LoadReport загружает UX отчёт по ID
	LoadReport(reportID string) (*UXReport, error)

	// SaveReport сохраняет UX отчёт
	SaveReport(report *UXReport) error

	// ListReports возвращает список всех UX отчётов определенного типа
	ListReports(reportType UXReportType) ([]*UXReport, error)

	// DeleteReport удаляет UX отчёт
	DeleteReport(reportID string) error
}
