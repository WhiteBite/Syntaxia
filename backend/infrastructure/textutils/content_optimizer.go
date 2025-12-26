package textutils

import (
	"context"
	"path/filepath"
	"strings"
)

// OptimizeOptions опции оптимизации контента
type OptimizeOptions struct {
	CollapseEmptyLines bool // Схлопывать множественные пустые строки
	StripLicense       bool // Удалять лицензионные заголовки
	StripComments      bool // Удалять комментарии (использует существующий CommentStripper)
	CompactDataFiles   bool // Сжимать JSON/YAML файлы
	SkeletonMode       bool // Генерировать только скелет кода (AST-based)
	TrimWhitespace     bool // Удалять trailing whitespace
}

// DefaultOptimizeOptions возвращает опции по умолчанию (минимальные оптимизации)
func DefaultOptimizeOptions() OptimizeOptions {
	return OptimizeOptions{
		CollapseEmptyLines: true,
		TrimWhitespace:     true,
	}
}

// AggressiveOptimizeOptions возвращает агрессивные опции для максимальной экономии
func AggressiveOptimizeOptions() OptimizeOptions {
	return OptimizeOptions{
		CollapseEmptyLines: true,
		StripLicense:       true,
		StripComments:      true,
		CompactDataFiles:   true,
		TrimWhitespace:     true,
	}
}

// SkeletonOptimizeOptions возвращает опции для режима скелета
func SkeletonOptimizeOptions() OptimizeOptions {
	return OptimizeOptions{
		SkeletonMode: true,
	}
}

// ContentOptimizer оптимизирует контент файлов для AI-контекста
// Объединяет все оптимизации в единый pipeline
type ContentOptimizer struct {
	whitespace  *WhitespaceOptimizer
	license     *LicenseStripper
	dataCompact *DataCompactor
	skeleton    *SkeletonGenerator
	comments    CommentStripperInterface
}

// CommentStripperInterface интерфейс для удаления комментариев
// Позволяет использовать существующий CommentStripper
type CommentStripperInterface interface {
	Strip(content string, filePath string) string
}

// NewContentOptimizer создает новый оптимизатор контента
func NewContentOptimizer(registry AnalyzerRegistry, commentStripper CommentStripperInterface) *ContentOptimizer {
	return &ContentOptimizer{
		whitespace:  NewWhitespaceOptimizer(),
		license:     NewLicenseStripper(),
		dataCompact: NewDataCompactor(),
		skeleton:    NewSkeletonGenerator(registry),
		comments:    commentStripper,
	}
}

// NewContentOptimizerSimple создает оптимизатор без skeleton generator
// Используется когда registry недоступен
func NewContentOptimizerSimple(commentStripper CommentStripperInterface) *ContentOptimizer {
	return &ContentOptimizer{
		whitespace:  NewWhitespaceOptimizer(),
		license:     NewLicenseStripper(),
		dataCompact: NewDataCompactor(),
		skeleton:    nil,
		comments:    commentStripper,
	}
}

// Optimize применяет оптимизации к контенту файла
// Порядок оптимизаций важен для максимальной эффективности
func (o *ContentOptimizer) Optimize(ctx context.Context, content, filePath string, opts OptimizeOptions) string {
	if len(content) == 0 {
		return content
	}

	// 1. Skeleton mode - самый агрессивный, применяется первым
	// Если успешно сгенерирован скелет, остальные оптимизации не нужны
	if opts.SkeletonMode && o.skeleton != nil {
		if skeleton := o.skeleton.Generate(ctx, content, filePath); skeleton != "" {
			return skeleton
		}
		// Если скелет не сгенерирован, продолжаем с обычными оптимизациями
	}

	result := content

	// 2. Удаление лицензий (до удаления комментариев, т.к. лицензия часто в комментарии)
	if opts.StripLicense {
		ext := strings.ToLower(filepath.Ext(filePath))
		result = o.license.StripWithLanguageHint(result, ext)
	}

	// 3. Удаление комментариев
	if opts.StripComments && o.comments != nil {
		result = o.comments.Strip(result, filePath)
	}

	// 4. Сжатие JSON/YAML
	if opts.CompactDataFiles && IsDataFile(filePath) {
		result = o.dataCompact.Compact(result, filePath)
	}

	// 5. Удаление trailing whitespace
	if opts.TrimWhitespace {
		result = o.whitespace.TrimTrailingWhitespace(result)
	}

	// 6. Схлопывание пустых строк (последним, после всех удалений)
	if opts.CollapseEmptyLines {
		result = o.whitespace.CollapseEmptyLines(result)
	}

	return result
}

// OptimizeWithDefaults применяет оптимизации с настройками по умолчанию
func (o *ContentOptimizer) OptimizeWithDefaults(ctx context.Context, content, filePath string) string {
	return o.Optimize(ctx, content, filePath, DefaultOptimizeOptions())
}

// OptimizeAggressive применяет агрессивные оптимизации
func (o *ContentOptimizer) OptimizeAggressive(ctx context.Context, content, filePath string) string {
	return o.Optimize(ctx, content, filePath, AggressiveOptimizeOptions())
}

// GenerateSkeleton генерирует скелет кода
func (o *ContentOptimizer) GenerateSkeleton(ctx context.Context, content, filePath string) string {
	if o.skeleton == nil {
		return ""
	}
	return o.skeleton.Generate(ctx, content, filePath)
}

// CanGenerateSkeleton проверяет возможность генерации скелета
func (o *ContentOptimizer) CanGenerateSkeleton(filePath string) bool {
	if o.skeleton == nil {
		return false
	}
	return o.skeleton.CanGenerateSkeleton(filePath)
}

// EstimateSavings оценивает экономию токенов для данных опций
// Возвращает примерный процент экономии
func EstimateSavings(opts OptimizeOptions) int {
	savings := 0

	if opts.SkeletonMode {
		return 60 // ~60% экономия в skeleton mode
	}

	if opts.CollapseEmptyLines {
		savings += 5
	}
	if opts.StripLicense {
		savings += 3
	}
	if opts.StripComments {
		savings += 15
	}
	if opts.CompactDataFiles {
		savings += 10
	}
	if opts.TrimWhitespace {
		savings += 2
	}

	return savings
}

// OptimizeOptionsFromFlags создает опции из флагов
func OptimizeOptionsFromFlags(
	collapseEmpty, stripLicense, stripComments, compactData, skeleton, trimWS bool,
) OptimizeOptions {
	return OptimizeOptions{
		CollapseEmptyLines: collapseEmpty,
		StripLicense:       stripLicense,
		StripComments:      stripComments,
		CompactDataFiles:   compactData,
		SkeletonMode:       skeleton,
		TrimWhitespace:     trimWS,
	}
}

// OptimizeBatch оптимизирует несколько файлов параллельно
// Возвращает map[filePath]optimizedContent
func (o *ContentOptimizer) OptimizeBatch(
	ctx context.Context,
	files map[string]string,
	opts OptimizeOptions,
) map[string]string {
	result := make(map[string]string, len(files))

	// Для небольшого количества файлов - последовательно
	// Для большого - можно добавить параллелизм
	for path, content := range files {
		select {
		case <-ctx.Done():
			return result
		default:
			result[path] = o.Optimize(ctx, content, path, opts)
		}
	}

	return result
}

// Stats содержит статистику оптимизации
type OptimizeStats struct {
	OriginalSize   int
	OptimizedSize  int
	SavedBytes     int
	SavedPercent   float64
	FilesProcessed int
	SkeletonsUsed  int
}

// OptimizeBatchWithStats оптимизирует файлы и возвращает статистику
func (o *ContentOptimizer) OptimizeBatchWithStats(
	ctx context.Context,
	files map[string]string,
	opts OptimizeOptions,
) (map[string]string, OptimizeStats) {
	result := make(map[string]string, len(files))
	stats := OptimizeStats{}

	for path, content := range files {
		select {
		case <-ctx.Done():
			return result, stats
		default:
			stats.OriginalSize += len(content)
			optimized := o.Optimize(ctx, content, path, opts)
			result[path] = optimized
			stats.OptimizedSize += len(optimized)
			stats.FilesProcessed++

			// Проверяем, был ли использован skeleton
			if opts.SkeletonMode && o.skeleton != nil && o.skeleton.CanGenerateSkeleton(path) {
				stats.SkeletonsUsed++
			}
		}
	}

	stats.SavedBytes = stats.OriginalSize - stats.OptimizedSize
	if stats.OriginalSize > 0 {
		stats.SavedPercent = float64(stats.SavedBytes) / float64(stats.OriginalSize) * 100
	}

	return result, stats
}
