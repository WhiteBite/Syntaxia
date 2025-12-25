package reportfs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

// ReportFileSystemRepository implements ReportRepository using the file system
type ReportFileSystemRepository struct {
	logger    domain.Logger
	reportDir string
}

// NewReportFileSystemRepository creates a new file system report repository
func NewReportFileSystemRepository(logger domain.Logger) (*ReportFileSystemRepository, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	reportDir := filepath.Join(homeDir, ".syntaxia", "reports")
	if err := os.MkdirAll(reportDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create report directory: %w", err)
	}

	return &ReportFileSystemRepository{
		logger:    logger,
		reportDir: reportDir,
	}, nil
}

// GetReport retrieves a report by ID
func (r *ReportFileSystemRepository) GetReport(ctx context.Context, reportID string) (*domain.GenericReport, error) {
	reportPath := filepath.Join(r.reportDir, reportID+".json")

	data, err := os.ReadFile(reportPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("report not found: %s", reportID)
		}
		return nil, fmt.Errorf("failed to read report file: %w", err)
	}

	var report domain.GenericReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to unmarshal report: %w", err)
	}

	return &report, nil
}

// ListReports lists all reports, optionally filtered by type
func (r *ReportFileSystemRepository) ListReports(ctx context.Context, reportType string) ([]*domain.GenericReport, error) {
	entries, err := os.ReadDir(r.reportDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read report directory: %w", err)
	}

	reports := make([]*domain.GenericReport, 0, len(entries))

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		reportID := strings.TrimSuffix(entry.Name(), ".json")
		report, err := r.GetReport(ctx, reportID)
		if err != nil {
			r.logger.Warning(fmt.Sprintf("Failed to load report %s: %v", reportID, err))
			continue
		}

		// Filter by type if specified
		if reportType != "" && report.Type != reportType {
			continue
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// SaveReport saves a report
func (r *ReportFileSystemRepository) SaveReport(ctx context.Context, report *domain.GenericReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report: %w", err)
	}

	reportPath := filepath.Join(r.reportDir, report.Id+".json")
	if err := os.WriteFile(reportPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write report file: %w", err)
	}

	return nil
}

// DeleteReport deletes a report by ID
func (r *ReportFileSystemRepository) DeleteReport(ctx context.Context, reportID string) error {
	reportPath := filepath.Join(r.reportDir, reportID+".json")

	if err := os.Remove(reportPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("report not found: %s", reportID)
		}
		return fmt.Errorf("failed to delete report file: %w", err)
	}

	return nil
}

// GetReportsByTask retrieves all reports for a specific task
func (r *ReportFileSystemRepository) GetReportsByTask(ctx context.Context, taskID string) ([]*domain.GenericReport, error) {
	entries, err := os.ReadDir(r.reportDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read report directory: %w", err)
	}

	var reports []*domain.GenericReport

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		reportID := strings.TrimSuffix(entry.Name(), ".json")
		report, err := r.GetReport(ctx, reportID)
		if err != nil {
			r.logger.Warning(fmt.Sprintf("Failed to load report %s: %v", reportID, err))
			continue
		}

		if report.TaskId == taskID {
			reports = append(reports, report)
		}
	}

	return reports, nil
}
