package diff

import (
	"context"
	"fmt"

	"syntaxia/domain"
)

// Service предоставляет высокоуровневый API для работы с diff
type Service struct {
	log    domain.Logger
	engine domain.DiffEngine
}

// NewService создает новый сервис diff
func NewService(log domain.Logger, engine domain.DiffEngine) *Service {
	return &Service{log: log, engine: engine}
}

// GenerateDiff генерирует diff между двумя состояниями
func (s *Service) GenerateDiff(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	s.log.Info(fmt.Sprintf("Generating diff between %s and %s", beforePath, afterPath))
	return s.engine.GenerateDiff(ctx, beforePath, afterPath, format)
}

// GenerateDiffFromResults генерирует diff из результатов применения правок
func (s *Service) GenerateDiffFromResults(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
	s.log.Info(fmt.Sprintf("Generating diff from %d apply results", len(results)))
	return s.engine.GenerateDiffFromResults(ctx, results, format)
}

// GenerateDiffFromEdits генерирует diff из Edits JSON
func (s *Service) GenerateDiffFromEdits(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
	s.log.Info(fmt.Sprintf("Generating diff from %d edits", len(edits.Edits)))
	return s.engine.GenerateDiffFromEdits(ctx, edits, format)
}

// PublishDiff публикует diff
func (s *Service) PublishDiff(ctx context.Context, diff *domain.DiffResult) error {
	s.log.Info(fmt.Sprintf("Publishing diff %s", diff.ID))
	return s.engine.PublishDiff(ctx, diff)
}

// GenerateAndPublishDiff генерирует и публикует diff
func (s *Service) GenerateAndPublishDiff(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	diff, err := s.GenerateDiff(ctx, beforePath, afterPath, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate diff: %w", err)
	}
	if err := s.PublishDiff(ctx, diff); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to publish diff: %v", err))
	}
	return diff, nil
}

// GenerateAndPublishDiffFromResults генерирует и публикует diff из результатов
func (s *Service) GenerateAndPublishDiffFromResults(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
	diff, err := s.GenerateDiffFromResults(ctx, results, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate diff from results: %w", err)
	}
	if err := s.PublishDiff(ctx, diff); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to publish diff: %v", err))
	}
	return diff, nil
}

// GenerateAndPublishDiffFromEdits генерирует и публикует diff из правок
func (s *Service) GenerateAndPublishDiffFromEdits(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
	diff, err := s.GenerateDiffFromEdits(ctx, edits, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate diff from edits: %w", err)
	}
	if err := s.PublishDiff(ctx, diff); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to publish diff: %v", err))
	}
	return diff, nil
}

// GetSupportedFormats возвращает поддерживаемые форматы
func (s *Service) GetSupportedFormats() []domain.DiffFormat {
	return []domain.DiffFormat{
		domain.DiffFormatGit,
		domain.DiffFormatUnified,
		domain.DiffFormatJSON,
		domain.DiffFormatHTML,
	}
}
