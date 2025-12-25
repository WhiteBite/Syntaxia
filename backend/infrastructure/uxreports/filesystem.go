package uxreports

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
)

// FileSystemUXReportRepository implements UXReportRepository using file system
type FileSystemUXReportRepository struct {
	reportsDir string
}

// NewFileSystemUXReportRepository creates a new file system UX report repository
func NewFileSystemUXReportRepository(reportsDir string) *FileSystemUXReportRepository {
	return &FileSystemUXReportRepository{
		reportsDir: reportsDir,
	}
}

// LoadReport loads a UX report by ID
func (r *FileSystemUXReportRepository) LoadReport(reportID string) (*domain.UXReport, error) {
	reportPath := filepath.Join(r.reportsDir, fmt.Sprintf("%s.json", reportID))
	data, err := os.ReadFile(reportPath)
	if err != nil {
		return nil, err
	}

	var report domain.UXReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, err
	}

	return &report, nil
}

// SaveReport saves a UX report
func (r *FileSystemUXReportRepository) SaveReport(report *domain.UXReport) error {
	// Create directory if needed
	if err := os.MkdirAll(r.reportsDir, 0o755); err != nil {
		return err
	}

	reportPath := filepath.Join(r.reportsDir, fmt.Sprintf("%s.json", report.ID))
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(reportPath, data, 0o600)
}

// ListReports returns all UX reports of a specific type
func (r *FileSystemUXReportRepository) ListReports(reportType domain.UXReportType) ([]*domain.UXReport, error) {
	// Create directory if needed
	if err := os.MkdirAll(r.reportsDir, 0o755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(r.reportsDir)
	if err != nil {
		return nil, err
	}

	reports := make([]*domain.UXReport, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a JSON file
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		reportPath := filepath.Join(r.reportsDir, entry.Name())
		data, err := os.ReadFile(reportPath)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		var report domain.UXReport
		if err := json.Unmarshal(data, &report); err != nil {
			// Skip files that can't be parsed
			continue
		}

		// Filter by report type if specified
		if reportType != "" && report.Type != reportType {
			continue
		}

		reports = append(reports, &report)
	}

	return reports, nil
}

// DeleteReport deletes a UX report
func (r *FileSystemUXReportRepository) DeleteReport(reportID string) error {
	reportPath := filepath.Join(r.reportsDir, fmt.Sprintf("%s.json", reportID))
	return os.Remove(reportPath)
}
