package fsscanner

import (
	"path/filepath"
	"strings"
	"sync"
	"syntaxia/domain"
	"time"
)

// fileSearcher реализует domain.FileSearcher
type fileSearcher struct {
	treeBuilder domain.TreeBuilder
	logger      domain.Logger

	mu          sync.RWMutex
	searchCache map[string]*searchIndex
	cacheSize   int64
	cacheHits   int64
	cacheMisses int64
}

// searchIndex индекс для быстрого поиска файлов
type searchIndex struct {
	files     []indexedFile
	createdAt time.Time
	projectRoot string
}

// indexedFile файл в индексе поиска
type indexedFile struct {
	path        string
	name        string
	size        int64
	contentType string
	depth       int
	nameLower   string // для case-insensitive поиска
	pathLower   string
}

const (
	maxSearchCacheEntries = 5
	searchCacheDuration   = 5 * time.Minute
	maxSearchCacheSizeMB  = 50
)

// NewFileSearcher создает новый FileSearcher
func NewFileSearcher(treeBuilder domain.TreeBuilder, logger domain.Logger) domain.FileSearcher {
	return &fileSearcher{
		treeBuilder: treeBuilder,
		logger:      logger,
		searchCache: make(map[string]*searchIndex),
	}
}

// BuildSearchIndex создает индекс поиска для проекта
func (s *fileSearcher) BuildSearchIndex(projectRoot string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, есть ли уже актуальный индекс
	if idx, exists := s.searchCache[projectRoot]; exists {
		if time.Since(idx.createdAt) < searchCacheDuration {
			s.logger.Debug("Search index already exists for: " + projectRoot)
			return nil
		}
	}

	// Получаем дерево файлов (используем кэш tree builder)
	tree, err := s.treeBuilder.BuildTree(projectRoot, true, true)
	if err != nil {
		return err
	}

	// Строим индекс из дерева
	files := make([]indexedFile, 0, 1000)
	s.indexTree(tree, &files, 0)

	// Сохраняем индекс
	s.evictOldestIfNeeded()
	
	idx := &searchIndex{
		files:       files,
		createdAt:   time.Now(),
		projectRoot: projectRoot,
	}
	s.searchCache[projectRoot] = idx
	s.cacheSize += s.estimateIndexSize(idx)

	s.logger.Info("Built search index for " + projectRoot + " with " + string(rune(len(files))) + " files")
	return nil
}

// indexTree рекурсивно индексирует дерево файлов
func (s *fileSearcher) indexTree(nodes []*domain.FileNode, files *[]indexedFile, depth int) {
	for _, node := range nodes {
		if node.IsDir {
			// Рекурсивно обрабатываем директории
			if node.Children != nil {
				s.indexTree(node.Children, files, depth+1)
			}
		} else {
			// Индексируем файл
			*files = append(*files, indexedFile{
				path:        node.Path,
				name:        node.Name,
				size:        node.Size,
				contentType: node.ContentType,
				depth:       depth,
				nameLower:   strings.ToLower(node.Name),
				pathLower:   strings.ToLower(node.Path),
			})
		}
	}
}

// SearchFiles выполняет поиск файлов по запросу
func (s *fileSearcher) SearchFiles(projectRoot, query string, options domain.SearchOptions) ([]domain.FileSearchResult, error) {
	startTime := time.Now()

	// Получаем или создаем индекс
	idx, err := s.getOrCreateIndex(projectRoot)
	if err != nil {
		return nil, err
	}

	// Выполняем поиск
	results := s.performSearch(idx, query, options)

	s.logger.Debug("Search completed in " + time.Since(startTime).String())
	return results, nil
}

// getOrCreateIndex получает существующий индекс или создает новый
func (s *fileSearcher) getOrCreateIndex(projectRoot string) (*searchIndex, error) {
	s.mu.RLock()
	idx, exists := s.searchCache[projectRoot]
	s.mu.RUnlock()

	if exists && time.Since(idx.createdAt) < searchCacheDuration {
		s.mu.Lock()
		s.cacheHits++
		s.mu.Unlock()
		return idx, nil
	}

	s.mu.Lock()
	s.cacheMisses++
	s.mu.Unlock()

	// Создаем новый индекс
	if err := s.BuildSearchIndex(projectRoot); err != nil {
		return nil, err
	}

	s.mu.RLock()
	idx = s.searchCache[projectRoot]
	s.mu.RUnlock()

	return idx, nil
}

// performSearch выполняет поиск по индексу
func (s *fileSearcher) performSearch(idx *searchIndex, query string, options domain.SearchOptions) []domain.FileSearchResult {
	queryLower := strings.ToLower(query)
	results := make([]domain.FileSearchResult, 0, options.MaxResults)

	for _, file := range idx.files {
		// Фильтр по типам файлов
		if len(options.FileTypes) > 0 {
			ext := filepath.Ext(file.name)
			if !contains(options.FileTypes, ext) {
				continue
			}
		}

		// Поиск совпадений
		var score float64
		var matches []domain.MatchRange

		if options.CaseSensitive {
			score, matches = s.matchFile(file.name, file.path, query, options.IncludePath)
		} else {
			score, matches = s.matchFile(file.nameLower, file.pathLower, queryLower, options.IncludePath)
		}

		if score > 0 {
			results = append(results, domain.FileSearchResult{
				Path:        file.path,
				Name:        file.name,
				Score:       score,
				Matches:     matches,
				Size:        file.size,
				ContentType: file.contentType,
				Depth:       file.depth,
			})

			if len(results) >= options.MaxResults {
				break
			}
		}
	}

	return results
}

// matchFile проверяет совпадение файла с запросом
func (s *fileSearcher) matchFile(name, path, query string, includePath bool) (float64, []domain.MatchRange) {
	var matches []domain.MatchRange
	var score float64

	// Поиск в имени файла
	nameIdx := strings.Index(name, query)
	if nameIdx != -1 {
		matches = append(matches, domain.MatchRange{
			Start: nameIdx,
			End:   nameIdx + len(query),
			Field: "name",
		})
		// Точное совпадение в начале имени - высокий score
		if nameIdx == 0 {
			score = 1.0
		} else {
			score = 0.8
		}
	}

	// Поиск в пути (если включено)
	if includePath {
		pathIdx := strings.Index(path, query)
		if pathIdx != -1 {
			matches = append(matches, domain.MatchRange{
				Start: pathIdx,
				End:   pathIdx + len(query),
				Field: "path",
			})
			if score == 0 {
				score = 0.5
			}
		}
	}

	return score, matches
}

// InvalidateSearchIndex инвалидирует индекс поиска для проекта
func (s *fileSearcher) InvalidateSearchIndex(projectRoot string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if idx, exists := s.searchCache[projectRoot]; exists {
		s.cacheSize -= s.estimateIndexSize(idx)
		delete(s.searchCache, projectRoot)
		s.logger.Debug("Invalidated search index for: " + projectRoot)
	}
}

// GetSearchStats возвращает статистику индекса поиска
func (s *fileSearcher) GetSearchStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"cached_indices": len(s.searchCache),
		"cache_size_mb":  s.cacheSize / (1024 * 1024),
		"cache_hits":     s.cacheHits,
		"cache_misses":   s.cacheMisses,
	}
}

// evictOldestIfNeeded удаляет старые индексы при превышении лимитов
func (s *fileSearcher) evictOldestIfNeeded() {
	maxSize := int64(maxSearchCacheSizeMB * 1024 * 1024)

	for len(s.searchCache) >= maxSearchCacheEntries || s.cacheSize > maxSize {
		var oldestKey string
		oldestTime := time.Now()

		for key, idx := range s.searchCache {
			if idx.createdAt.Before(oldestTime) {
				oldestTime = idx.createdAt
				oldestKey = key
			}
		}

		if oldestKey != "" {
			idx := s.searchCache[oldestKey]
			s.cacheSize -= s.estimateIndexSize(idx)
			delete(s.searchCache, oldestKey)
			s.logger.Debug("Evicted search index: " + oldestKey)
		} else {
			break
		}
	}
}

// estimateIndexSize оценивает размер индекса в байтах
func (s *fileSearcher) estimateIndexSize(idx *searchIndex) int64 {
	var size int64
	for _, file := range idx.files {
		size += int64(len(file.path) + len(file.name) + len(file.contentType) + 100)
	}
	return size
}

// contains проверяет наличие элемента в слайсе
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
