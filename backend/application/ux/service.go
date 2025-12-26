package ux

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	"sync"
	"time"
)

// ServiceImpl реализует UXMetricsService
type ServiceImpl struct {
	log     domain.Logger
	reports map[string]*domain.UXReport
	mu      sync.RWMutex
	repo    domain.UXReportRepository
}

// NewService создает новый сервис UX метрик
func NewService(log domain.Logger, repo domain.UXReportRepository) domain.UXMetricsService {
	return &ServiceImpl{
		log:     log,
		reports: make(map[string]*domain.UXReport),
		repo:    repo,
	}
}

// GetUXReport возвращает UX отчёт
func (s *ServiceImpl) GetUXReport(reportID string) (*domain.UXReport, error) {
	s.mu.RLock()
	report, exists := s.reports[reportID]
	s.mu.RUnlock()

	if !exists {
		loadedReport, err := s.repo.LoadReport(reportID)
		if err != nil {
			return nil, fmt.Errorf("UX report not found: %s", reportID)
		}

		s.mu.Lock()
		s.reports[reportID] = loadedReport
		s.mu.Unlock()

		return loadedReport, nil
	}

	return report, nil
}

// SaveUXReport сохраняет UX отчёт
func (s *ServiceImpl) SaveUXReport(report *domain.UXReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reports[report.ID] = report

	if err := s.repo.SaveReport(report); err != nil {
		return fmt.Errorf("failed to save report through repository: %w", err)
	}

	s.log.Info(fmt.Sprintf("Saved UX report: %s", report.ID))
	return nil
}

// GetUXReports возвращает все UX отчёты
func (s *ServiceImpl) GetUXReports(reportType domain.UXReportType) ([]*domain.UXReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.reports) > 0 {
		var reports []*domain.UXReport
		for _, report := range s.reports {
			if reportType == "" || report.Type == reportType {
				reports = append(reports, report)
			}
		}
		return reports, nil
	}

	reports, err := s.repo.ListReports(reportType)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	for _, report := range reports {
		s.reports[report.ID] = report
	}
	s.mu.Unlock()

	return reports, nil
}

// DeleteUXReport удаляет UX отчёт
func (s *ServiceImpl) DeleteUXReport(reportID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.repo.DeleteReport(reportID); err != nil {
		return fmt.Errorf("failed to delete report through repository: %w", err)
	}

	delete(s.reports, reportID)

	s.log.Info(fmt.Sprintf("Deleted UX report: %s", reportID))
	return nil
}

// GetMetricsSummary возвращает сводку метрик
func (s *ServiceImpl) GetMetricsSummary() (map[string]any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	summary := map[string]any{
		"totalReports":      len(s.reports),
		"reportsByType":     make(map[string]int),
		"lastReportTime":    nil,
		"averageConfidence": 0.0,
	}

	typeCounts := make(map[string]int)
	var totalConfidence float64
	var confidenceCount int
	var lastReportTime *time.Time

	for _, report := range s.reports {
		typeCounts[string(report.Type)]++

		if report.Type == domain.UXReportTypeWhyView {
			if whyView, ok := report.Content.(*domain.WhyViewReport); ok {
				totalConfidence += whyView.Confidence
				confidenceCount++
			}
		}

		if lastReportTime == nil || report.CreatedAt.After(*lastReportTime) {
			lastReportTime = &report.CreatedAt
		}
	}

	summary["reportsByType"] = typeCounts
	summary["lastReportTime"] = lastReportTime

	if confidenceCount > 0 {
		summary["averageConfidence"] = totalConfidence / float64(confidenceCount)
	}

	return summary, nil
}

// ListReports получает список отчетов
func (s *ServiceImpl) ListReports(_ context.Context, reportType string) ([]*domain.GenericReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reports := make([]*domain.GenericReport, 0, len(s.reports))

	for _, report := range s.reports {
		if reportType != "" && string(report.Type) != reportType {
			continue
		}

		genericReport := &domain.GenericReport{
			Id:        report.ID,
			TaskId:    s.extractTaskID(report),
			Title:     report.Title,
			Type:      string(report.Type),
			Summary:   report.Description,
			Content:   s.serializeContent(report.Content),
			CreatedAt: report.CreatedAt,
			UpdatedAt: report.CreatedAt,
		}

		reports = append(reports, genericReport)
	}

	return reports, nil
}

// GetReport получает конкретный отчет
func (s *ServiceImpl) GetReport(_ context.Context, reportID string) (*domain.GenericReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	report, exists := s.reports[reportID]
	if !exists {
		return nil, fmt.Errorf("report not found: %s", reportID)
	}

	genericReport := &domain.GenericReport{
		Id:        report.ID,
		TaskId:    s.extractTaskID(report),
		Title:     report.Title,
		Type:      string(report.Type),
		Summary:   report.Description,
		Content:   s.serializeContent(report.Content),
		CreatedAt: report.CreatedAt,
		UpdatedAt: report.CreatedAt,
	}

	return genericReport, nil
}

// extractTaskID извлекает TaskID из отчета
func (s *ServiceImpl) extractTaskID(report *domain.UXReport) string {
	if taskID, ok := report.Metadata["taskID"].(string); ok {
		return taskID
	}
	return ""
}

// serializeContent сериализует содержимое отчета в JSON
func (s *ServiceImpl) serializeContent(content any) string {
	if content == nil {
		return ""
	}

	data, err := json.Marshal(content)
	if err != nil {
		return fmt.Sprintf("Error serializing content: %v", err)
	}

	return string(data)
}

// logSaveWarning logs a warning when saving UX report fails (non-critical)
func (s *ServiceImpl) logSaveWarning(err error) {
	s.log.Warning(fmt.Sprintf("Failed to save UX report: %v", err))
}
