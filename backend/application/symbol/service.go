package symbol

import (
	"context"
	"fmt"
	"syntaxia/domain"
	"sync"
	"time"
)

// Service предоставляет высокоуровневый API для работы с графом символов
type Service struct {
	log                 domain.Logger
	symbolGraphBuilders map[string]domain.SymbolGraphBuilder
	importGraphBuilders map[string]domain.ImportGraphBuilder
	cache               map[string]*domain.SymbolGraph
	cacheTimestamps     map[string]int64
	lastAccessed        map[string]int64
	cacheSize           int64
	cacheHits           int64
	cacheMisses         int64
	evictions           int64
	mu                  sync.RWMutex
}

const (
	maxCacheGraphs = 10
	maxCacheSizeMB = 50
)

// NewService создает новый сервис графа символов
func NewService(log domain.Logger, symbolGraphBuilders map[string]domain.SymbolGraphBuilder, importGraphBuilders map[string]domain.ImportGraphBuilder) *Service {
	return &Service{
		log:                 log,
		symbolGraphBuilders: symbolGraphBuilders,
		importGraphBuilders: importGraphBuilders,
		cache:               make(map[string]*domain.SymbolGraph),
		cacheTimestamps:     make(map[string]int64),
		lastAccessed:        make(map[string]int64),
	}
}

// RegisterSymbolGraphBuilder регистрирует builder для языка
func (s *Service) RegisterSymbolGraphBuilder(language string, builder domain.SymbolGraphBuilder) {
	s.symbolGraphBuilders[language] = builder
	s.log.Info(fmt.Sprintf("Registered symbol graph builder for language: %s", language))
}

// RegisterImportGraphBuilder регистрирует builder импортов для языка
func (s *Service) RegisterImportGraphBuilder(language string, builder domain.ImportGraphBuilder) {
	s.importGraphBuilders[language] = builder
	s.log.Info(fmt.Sprintf("Registered import graph builder for language: %s", language))
}

// BuildSymbolGraph строит граф символов для проекта
func (s *Service) BuildSymbolGraph(ctx context.Context, projectRoot, language string) (*domain.SymbolGraph, error) {
	builder, exists := s.symbolGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no symbol graph builder registered for language: %s", language)
	}

	// Проверяем кэш (только чтение под RLock)
	cacheKey := fmt.Sprintf("%s:%s", projectRoot, language)
	s.mu.RLock()
	_, cacheExists := s.cache[cacheKey]
	s.mu.RUnlock()

	if cacheExists {
		// Cache hit - обновляем метрики под Lock
		s.mu.Lock()
		// Повторная проверка после получения Lock (double-checked locking)
		if graph, stillExists := s.cache[cacheKey]; stillExists {
			s.lastAccessed[cacheKey] = time.Now().Unix()
			s.cacheHits++
			s.mu.Unlock()
			s.log.Info(fmt.Sprintf("Using cached symbol graph for project: %s (language: %s)", projectRoot, language))
			return graph, nil
		}
		s.mu.Unlock()
	}

	// Cache miss - инкрементируем счётчик под Lock
	s.mu.Lock()
	s.cacheMisses++
	s.mu.Unlock()

	s.log.Info(fmt.Sprintf("Building symbol graph for project: %s (language: %s)", projectRoot, language))
	graph, err := builder.BuildGraph(ctx, projectRoot)
	if err != nil {
		return nil, err
	}

	// Кэшируем результат с LRU eviction
	s.mu.Lock()
	graphSize := s.estimateGraphSize(graph)

	// Если ключ уже существует, вычитаем старый размер
	if oldGraph, exists := s.cache[cacheKey]; exists {
		oldSize := s.estimateGraphSize(oldGraph)
		s.cacheSize -= oldSize
	}

	s.evictOldestIfNeeded(graphSize)
	s.cache[cacheKey] = graph
	s.cacheTimestamps[cacheKey] = time.Now().Unix()
	s.lastAccessed[cacheKey] = time.Now().Unix()
	s.cacheSize += graphSize
	s.mu.Unlock()

	return graph, nil
}

// estimateGraphSize оценивает размер графа в байтах
func (s *Service) estimateGraphSize(graph *domain.SymbolGraph) int64 {
	if graph == nil {
		return 0
	}
	nodeCount := len(graph.Nodes)
	avgNodeSize := int64(500)
	return int64(nodeCount) * avgNodeSize
}

// evictOldestIfNeeded удаляет старые записи при превышении лимитов
func (s *Service) evictOldestIfNeeded(newGraphSize int64) {
	maxSize := int64(maxCacheSizeMB * 1024 * 1024)

	for len(s.cache) >= maxCacheGraphs || (s.cacheSize+newGraphSize) > maxSize {
		if len(s.cache) == 0 {
			break
		}

		var oldestKey string
		oldestTime := time.Now().Unix()

		for key, accessTime := range s.lastAccessed {
			if accessTime < oldestTime {
				oldestTime = accessTime
				oldestKey = key
			}
		}

		if oldestKey == "" {
			break
		}

		if graph, exists := s.cache[oldestKey]; exists {
			s.cacheSize -= s.estimateGraphSize(graph)
		}
		delete(s.cache, oldestKey)
		delete(s.cacheTimestamps, oldestKey)
		delete(s.lastAccessed, oldestKey)
		s.evictions++

		s.log.Info(fmt.Sprintf("Evicted symbol graph from cache: %s", oldestKey))
	}
}

// BuildImportGraph строит граф импортов для проекта
func (s *Service) BuildImportGraph(ctx context.Context, projectRoot, language string) (*domain.ImportGraph, error) {
	builder, exists := s.importGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no import graph builder registered for language: %s", language)
	}

	s.log.Info(fmt.Sprintf("Building import graph for project: %s (language: %s)", projectRoot, language))
	return builder.BuildImportGraph(ctx, projectRoot)
}

// GetSuggestions возвращает предложения символов
func (s *Service) GetSuggestions(ctx context.Context, query, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	builder, exists := s.symbolGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no symbol graph builder registered for language: %s", language)
	}

	return builder.GetSuggestions(ctx, query, graph)
}

// GetDependencies возвращает зависимости символа
func (s *Service) GetDependencies(ctx context.Context, symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	builder, exists := s.symbolGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no symbol graph builder registered for language: %s", language)
	}

	return builder.GetDependencies(ctx, symbolID, graph)
}

// GetDependents возвращает символы, зависящие от указанного
func (s *Service) GetDependents(ctx context.Context, symbolID, language string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	builder, exists := s.symbolGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no symbol graph builder registered for language: %s", language)
	}

	return builder.GetDependents(ctx, symbolID, graph)
}

// GetCircularImports возвращает циклические импорты
func (s *Service) GetCircularImports(ctx context.Context, language string, graph *domain.ImportGraph) ([][]string, error) {
	builder, exists := s.importGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no import graph builder registered for language: %s", language)
	}

	return builder.GetCircularImports(ctx, graph)
}

// GetImportPath возвращает путь импорта между пакетами
func (s *Service) GetImportPath(ctx context.Context, from, to, language string, graph *domain.ImportGraph) ([]string, error) {
	builder, exists := s.importGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no import graph builder registered for language: %s", language)
	}

	return builder.GetImportPath(ctx, from, to, graph)
}

// GetSupportedLanguages возвращает список поддерживаемых языков
func (s *Service) GetSupportedLanguages() []string {
	languages := make([]string, 0, len(s.symbolGraphBuilders))
	for lang := range s.symbolGraphBuilders {
		languages = append(languages, lang)
	}
	return languages
}

// UpdateSymbolGraph инкрементально обновляет граф символов
func (s *Service) UpdateSymbolGraph(ctx context.Context, projectRoot, language string, changedFiles []string) (*domain.SymbolGraph, error) {
	builder, exists := s.symbolGraphBuilders[language]
	if !exists {
		return nil, fmt.Errorf("no symbol graph builder registered for language: %s", language)
	}

	s.log.Info(fmt.Sprintf("Updating symbol graph for project: %s (language: %s, changed files: %d)", projectRoot, language, len(changedFiles)))

	graph, err := builder.UpdateGraph(ctx, projectRoot, changedFiles)
	if err != nil {
		return nil, err
	}

	// Обновляем кэш с LRU eviction
	cacheKey := fmt.Sprintf("%s:%s", projectRoot, language)
	s.mu.Lock()
	graphSize := s.estimateGraphSize(graph)

	if oldGraph, exists := s.cache[cacheKey]; exists {
		oldSize := s.estimateGraphSize(oldGraph)
		s.cacheSize -= oldSize
	}

	s.evictOldestIfNeeded(graphSize)
	s.cache[cacheKey] = graph
	s.cacheTimestamps[cacheKey] = time.Now().Unix()
	s.lastAccessed[cacheKey] = time.Now().Unix()
	s.cacheSize += graphSize
	s.mu.Unlock()

	return graph, nil
}

// ClearCache очищает кэш графов символов
func (s *Service) ClearCache() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache = make(map[string]*domain.SymbolGraph)
	s.cacheTimestamps = make(map[string]int64)
	s.lastAccessed = make(map[string]int64)
	s.cacheSize = 0
	s.log.Info("Symbol graph cache cleared")
}

// GetCacheStats возвращает статистику кэша
func (s *Service) GetCacheStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return map[string]interface{}{
		"cached_graphs":    len(s.cache),
		"cache_timestamps": len(s.cacheTimestamps),
		"cache_size_mb":    s.cacheSize / (1024 * 1024),
		"cache_hits":       s.cacheHits,
		"cache_misses":     s.cacheMisses,
		"evictions":        s.evictions,
	}
}
