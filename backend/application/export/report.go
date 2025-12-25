package export

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"time"

	"github.com/google/uuid"
)

// ReportService handles report management operations.
type ReportService struct {
	logger     domain.Logger
	reportRepo domain.ReportRepository
}

// NewReportService creates a new report service.
func NewReportService(logger domain.Logger, reportRepo domain.ReportRepository) *ReportService {
	return &ReportService{
		logger:     logger,
		reportRepo: reportRepo,
	}
}

// GetReport retrieves a report by ID.
func (s *ReportService) GetReport(ctx context.Context, reportID string) (*domain.GenericReport, error) {
	return s.reportRepo.GetReport(ctx, reportID)
}

// ListReports lists all reports, optionally filtered by type.
func (s *ReportService) ListReports(ctx context.Context, reportType string) ([]*domain.GenericReport, error) {
	return s.reportRepo.ListReports(ctx, reportType)
}

// GenerateReport generates a new report based on type and parameters.
func (s *ReportService) GenerateReport(ctx context.Context, reportType string, parameters map[string]interface{}) (*domain.GenericReport, error) {
	reportID := uuid.New().String()
	now := time.Now()

	title := "Generated Report"
	if titleParam, ok := parameters["title"].(string); ok {
		title = titleParam
	}

	taskID := ""
	if taskIDParam, ok := parameters["taskId"].(string); ok {
		taskID = taskIDParam
	}

	var content, summary string
	switch reportType {
	case "analysis":
		content = s.generateAnalysisReport(parameters)
		summary = "Analysis report generated"
	case "export":
		content = s.generateExportReport(parameters)
		summary = "Export report generated"
	case "autonomous":
		content = s.generateAutonomousReport(parameters)
		summary = "Autonomous task report generated"
	default:
		content = "General report content"
		summary = "General report generated"
	}

	report := &domain.GenericReport{
		Id:        reportID,
		TaskId:    taskID,
		Type:      reportType,
		Title:     title,
		Summary:   summary,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.saveReport(report); err != nil {
		return nil, fmt.Errorf("failed to save generated report: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Generated report %s of type %s", reportID, reportType))
	return report, nil
}

// ExportReport exports a report in the specified format.
func (s *ReportService) ExportReport(ctx context.Context, reportID, format string) (*domain.ExportResult, error) {
	report, err := s.GetReport(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get report for export: %w", err)
	}

	exportResult := &domain.ExportResult{}

	switch format {
	case "json":
		exportData, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal report to JSON: %w", err)
		}
		exportResult.Text = string(exportData)
		exportResult.FileName = fmt.Sprintf("%s.json", reportID)

	case "txt":
		textContent := fmt.Sprintf("Report: %s\\nType: %s\\nTitle: %s\\n\\nSummary:\\n%s\\n\\nContent:\\n%s",
			report.Id, report.Type, report.Title, report.Summary, report.Content)
		exportResult.Text = textContent
		exportResult.FileName = fmt.Sprintf("%s.txt", reportID)

	case "csv":
		csvContent := fmt.Sprintf("ID,Type,Title,Summary,CreatedAt\\n\\\"%s\\\",\\\"%s\\\",\\\"%s\\\",\\\"%s\\\",\\\"%s\\\"",
			report.Id, report.Type, report.Title, report.Summary, report.CreatedAt.Format(time.RFC3339))
		exportResult.Text = csvContent
		exportResult.FileName = fmt.Sprintf("%s.csv", reportID)

	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}

	s.logger.Info(fmt.Sprintf("Exported report %s to format %s", reportID, format))
	return exportResult, nil
}

func (s *ReportService) generateAnalysisReport(parameters map[string]interface{}) string {
	return "Analysis report content based on provided parameters"
}

func (s *ReportService) generateExportReport(parameters map[string]interface{}) string {
	return "Export operation report with statistics and details"
}

func (s *ReportService) generateAutonomousReport(parameters map[string]interface{}) string {
	return "Autonomous task execution report with progress and results"
}

// CreateReport creates a new report.
func (s *ReportService) CreateReport(ctx context.Context, taskID, reportType, title, summary, content string) (*domain.GenericReport, error) {
	reportID := uuid.New().String()
	now := time.Now()

	report := &domain.GenericReport{
		Id:        reportID,
		TaskId:    taskID,
		Type:      reportType,
		Title:     title,
		Summary:   summary,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.saveReport(report); err != nil {
		return nil, fmt.Errorf("failed to save report: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Created report %s of type %s", reportID, reportType))
	return report, nil
}

// UpdateReport updates an existing report.
func (s *ReportService) UpdateReport(ctx context.Context, reportID, title, summary, content string) (*domain.GenericReport, error) {
	report, err := s.GetReport(ctx, reportID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		report.Title = title
	}
	if summary != "" {
		report.Summary = summary
	}
	if content != "" {
		report.Content = content
	}
	report.UpdatedAt = time.Now()

	if err := s.saveReport(report); err != nil {
		return nil, fmt.Errorf("failed to save updated report: %w", err)
	}

	s.logger.Info(fmt.Sprintf("Updated report %s", reportID))
	return report, nil
}

// DeleteReport deletes a report by ID.
func (s *ReportService) DeleteReport(ctx context.Context, reportID string) error {
	return s.reportRepo.DeleteReport(ctx, reportID)
}

func (s *ReportService) saveReport(report *domain.GenericReport) error {
	return s.reportRepo.SaveReport(context.Background(), report)
}

// GetReportsByTask retrieves all reports for a specific task.
func (s *ReportService) GetReportsByTask(ctx context.Context, taskID string) ([]*domain.GenericReport, error) {
	return s.reportRepo.GetReportsByTask(ctx, taskID)
}
