package uxreports

import (
	"syntaxia/domain"
	"sync"
	"time"
)

// InMemoryUXReportRepository implements UXReportRepository using in-memory storage
type InMemoryUXReportRepository struct {
	reports   map[string]*domain.UXReport
	createdAt map[string]time.Time
	mu        sync.RWMutex
}

const (
	maxReports       = 100
	maxReportsSizeMB = 10
	maxReportAge     = 24 * time.Hour
)

// NewInMemoryUXReportRepository creates a new in-memory UX report repository
func NewInMemoryUXReportRepository() *InMemoryUXReportRepository {
	repo := &InMemoryUXReportRepository{
		reports:   make(map[string]*domain.UXReport),
		createdAt: make(map[string]time.Time),
	}

	// Start periodic cleanup
	go repo.periodicCleanup()

	return repo
}

// LoadReport loads a UX report by ID from memory
func (r *InMemoryUXReportRepository) LoadReport(reportID string) (*domain.UXReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	report, exists := r.reports[reportID]
	if !exists {
		return nil, nil // Return nil if not found, consistent with filesystem implementation
	}

	return report, nil
}

// SaveReport saves a UX report to memory with size limits
func (r *InMemoryUXReportRepository) SaveReport(report *domain.UXReport) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check limits before saving
	if len(r.reports) >= maxReports {
		r.evictOldestReport()
	}

	// Check memory usage
	if r.getMemoryUsageUnlocked() > maxReportsSizeMB*1024*1024 {
		r.evictOldestReport()
	}

	r.reports[report.ID] = report
	r.createdAt[report.ID] = time.Now()
	return nil
}

// evictOldestReport удаляет самый старый отчет
func (r *InMemoryUXReportRepository) evictOldestReport() {
	var oldestID string
	oldestTime := time.Now()

	for id, createdTime := range r.createdAt {
		if createdTime.Before(oldestTime) {
			oldestTime = createdTime
			oldestID = id
		}
	}

	if oldestID != "" {
		delete(r.reports, oldestID)
		delete(r.createdAt, oldestID)
	}
}

// CleanupOldReports удаляет отчеты старше указанного возраста
func (r *InMemoryUXReportRepository) CleanupOldReports(maxAge time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, createdTime := range r.createdAt {
		if now.Sub(createdTime) > maxAge {
			delete(r.reports, id)
			delete(r.createdAt, id)
		}
	}
}

// periodicCleanup запускает периодическую очистку старых отчетов
func (r *InMemoryUXReportRepository) periodicCleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		r.CleanupOldReports(maxReportAge)
	}
}

// GetMemoryUsage возвращает приблизительный размер используемой памяти
func (r *InMemoryUXReportRepository) GetMemoryUsage() int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.getMemoryUsageUnlocked()
}

func (r *InMemoryUXReportRepository) getMemoryUsageUnlocked() int64 {
	var size int64
	for _, report := range r.reports {
		// Приблизительная оценка размера отчета
		size += int64(len(report.ID) + len(report.Title) + len(report.Description) + 500) // Базовый размер структуры
		// Добавляем размер метаданных
		for k, v := range report.Metadata {
			size += int64(len(k) + 100) // Приблизительный размер значения
			_ = v                       // Используем переменную
		}
	}
	return size
}

// ListReports returns all UX reports of a specific type from memory
func (r *InMemoryUXReportRepository) ListReports(reportType domain.UXReportType) ([]*domain.UXReport, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	reports := make([]*domain.UXReport, 0, len(r.reports))
	for _, report := range r.reports {
		// Filter by report type if specified
		if reportType != "" && report.Type != reportType {
			continue
		}

		reports = append(reports, report)
	}

	return reports, nil
}

// DeleteReport deletes a UX report from memory
func (r *InMemoryUXReportRepository) DeleteReport(reportID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.reports, reportID)
	return nil
}
