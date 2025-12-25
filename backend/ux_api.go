package main

import (
	"encoding/json"
	"fmt"
	"os"
	"syntaxia/domain"
	"syntaxia/infrastructure/version"
)

// === UX Metrics ===

// GenerateWhyViewReport generates "why these files" report
func (a *App) GenerateWhyViewReport(taskID string, files []string) (*domain.WhyViewReport, error) {
	return a.uxMetricsService.GenerateWhyViewReport(taskID, files, nil)
}

// GenerateTimeToGreenMetrics generates time_to_green metrics
func (a *App) GenerateTimeToGreenMetrics(taskID string) (*domain.TimeToGreenMetrics, error) {
	return a.uxMetricsService.GenerateTimeToGreenMetrics(taskID)
}

// GenerateDerivedDiffReport generates derived diff report
func (a *App) GenerateDerivedDiffReport(taskID string, originalDiff, derivedDiff string) (*domain.DerivedDiffReport, error) {
	return a.uxMetricsService.GenerateDerivedDiffReport(taskID, originalDiff, derivedDiff)
}

// GeneratePerformanceMetrics generates performance metrics
func (a *App) GeneratePerformanceMetrics(taskID string) (*domain.PerformanceMetrics, error) {
	return a.uxMetricsService.GeneratePerformanceMetrics(taskID)
}

// GetUXReport returns UX report
func (a *App) GetUXReport(reportID string) (*domain.UXReport, error) {
	return a.uxMetricsService.GetUXReport(reportID)
}

// SaveUXReport saves UX report
func (a *App) SaveUXReport(report *domain.UXReport) error {
	return a.uxMetricsService.SaveUXReport(report)
}

// GetUXReports returns all UX reports
func (a *App) GetUXReports(reportType domain.UXReportType) ([]*domain.UXReport, error) {
	return a.uxMetricsService.GetUXReports(reportType)
}

// DeleteUXReport deletes UX report
func (a *App) DeleteUXReport(reportID string) error {
	return a.uxMetricsService.DeleteUXReport(reportID)
}

// GetMetricsSummary returns metrics summary
func (a *App) GetMetricsSummary() (map[string]interface{}, error) {
	return a.uxMetricsService.GetMetricsSummary()
}

// ListReports lists reports
func (a *App) ListReports(reportType string) (string, error) {
	reports, err := a.uxMetricsService.ListReports(a.ctx, reportType)
	if err != nil {
		return "", a.transformError(err)
	}

	reportsJson, err := json.Marshal(reports)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal reports", err)
		return "", a.transformError(marshalErr)
	}

	return string(reportsJson), nil
}

// GetReport gets a specific report
func (a *App) GetReport(reportId string) (string, error) {
	report, err := a.uxMetricsService.GetReport(a.ctx, reportId)
	if err != nil {
		return "", a.transformError(err)
	}

	reportJson, err := json.Marshal(report)
	if err != nil {
		marshalErr := domain.NewInternalError("failed to marshal report", err)
		return "", a.transformError(marshalErr)
	}

	return string(reportJson), nil
}

// GenerateReport generates a new report
func (a *App) GenerateReport(reportType string, parametersJson string) (string, error) {
	var parameters map[string]interface{}
	if parametersJson != "" {
		if err := json.Unmarshal([]byte(parametersJson), &parameters); err != nil {
			return "", fmt.Errorf("failed to parse parameters JSON: %w", err)
		}
	}

	report, err := a.reportService.GenerateReport(a.ctx, reportType, parameters)
	if err != nil {
		return "", fmt.Errorf("failed to generate report: %w", err)
	}

	reportJson, err := json.Marshal(report)
	if err != nil {
		return "", fmt.Errorf("failed to marshal report: %w", err)
	}

	return string(reportJson), nil
}

// DeleteReport deletes a report by ID
func (a *App) DeleteReport(reportId string) error {
	return a.reportService.DeleteReport(a.ctx, reportId)
}

// ExportReport exports a report in the specified format
func (a *App) ExportReport(reportId string, format string) (string, error) {
	result, err := a.reportService.ExportReport(a.ctx, reportId, format)
	if err != nil {
		return "", fmt.Errorf("failed to export report: %w", err)
	}

	resultJson, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal export result: %w", err)
	}

	return string(resultJson), nil
}

// === Version & Releases ===

// GetVersionInfo returns current app version info
func (a *App) GetVersionInfo() version.Info {
	return version.GetInfo()
}

// GetReleases returns GitHub releases with update check
func (a *App) GetReleases() (*version.ReleasesResponse, error) {
	service := version.NewReleasesService()
	return service.GetReleases(a.ctx)
}

// GetFileStats returns file statistics
func (a *App) GetFileStats(filePath string) (string, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file stats: %w", err)
	}

	stats := map[string]interface{}{
		"name":    fileInfo.Name(),
		"size":    fileInfo.Size(),
		"modTime": fileInfo.ModTime().Unix(),
		"isDir":   fileInfo.IsDir(),
		"mode":    fileInfo.Mode().String(),
	}

	statsJson, err := json.Marshal(stats)
	if err != nil {
		return "", fmt.Errorf("failed to marshal file stats: %w", err)
	}

	return string(statsJson), nil
}
