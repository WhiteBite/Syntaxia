package fsscanner

import (
	"io/fs"
	"path/filepath"
	"syntaxia/domain"
	"syntaxia/infrastructure/textutils"
	"sort"
	"strings"
	"sync"
	"time"

	gitignore "github.com/sabhiram/go-gitignore"
)

// detectContentTypeByExt returns content type string based on file extension.
func detectContentTypeByExt(filename string) string {
	ct := textutils.DetectByExtension(filename)
	return ct.String()
}

type fileTreeBuilder struct {
	settingsRepo domain.SettingsRepository
	log          domain.Logger

	mu          sync.RWMutex
	giCache     map[string]*gitignore.GitIgnore // per-project .gitignore cache
	customCache *gitignore.GitIgnore            // compiled custom rules
	customHash  string                          // hash of custom rules content for cache invalidation

	// Cache for file trees with timestamps for invalidation
	treeCache        map[string]*cachedTree
	cacheAccessTimes map[string]time.Time
	cacheMutex       sync.RWMutex
	cacheDuration    time.Duration
	cacheSize        int64
	cacheHits        int64
	cacheMisses      int64
}

type cachedTree struct {
	nodes   []*domain.FileNode
	modTime time.Time
	size    int64
}

const (
	maxTreeCacheEntries = 5
	maxTreeCacheSizeMB  = 20
)

func New(settingsRepo domain.SettingsRepository, log domain.Logger) domain.TreeBuilder {
	return &fileTreeBuilder{
		settingsRepo:     settingsRepo,
		log:              log,
		giCache:          make(map[string]*gitignore.GitIgnore),
		treeCache:        make(map[string]*cachedTree),
		cacheAccessTimes: make(map[string]time.Time),
		cacheDuration:    2 * time.Minute, // Cache for 2 minutes (reduced from 5)
	}
}

func (b *fileTreeBuilder) BuildTree(dirPath string, useGitignore, useCustomIgnore bool) ([]*domain.FileNode, error) {
	if cached := b.getCachedTree(dirPath); cached != nil {
		b.log.Debug("Using cached tree for: " + dirPath)
		return cached, nil
	}

	gi, ci := b.getIgnoreMatchers(dirPath, useGitignore, useCustomIgnore)
	nodesMap := make(map[string]*domain.FileNode)
	root := b.createRootNode(dirPath)
	nodesMap[dirPath] = root

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		return b.processEntry(path, d, err, dirPath, gi, ci, nodesMap)
	})
	if err != nil {
		return nil, err
	}

	b.sortChildren(nodesMap)
	result := []*domain.FileNode{root}
	b.setCachedTree(dirPath, result)
	return result, nil
}

// getIgnoreMatchers returns gitignore and custom ignore matchers
func (b *fileTreeBuilder) getIgnoreMatchers(dirPath string, useGitignore, useCustomIgnore bool) (*gitignore.GitIgnore, *gitignore.GitIgnore) {
	var gi, ci *gitignore.GitIgnore
	if useGitignore {
		gi = b.getGitignore(dirPath)
	}
	if useCustomIgnore {
		ci = b.getCustomIgnore()
	}
	return gi, ci
}

// createRootNode creates the root node for the tree
func (b *fileTreeBuilder) createRootNode(dirPath string) *domain.FileNode {
	return &domain.FileNode{
		Name: filepath.Base(dirPath), Path: dirPath, RelPath: ".",
		IsDir: true, Children: []*domain.FileNode{},
	}
}

// processEntry processes a single directory entry during walk
func (b *fileTreeBuilder) processEntry(path string, d fs.DirEntry, err error, dirPath string, gi, ci *gitignore.GitIgnore, nodesMap map[string]*domain.FileNode) error {
	if err != nil {
		return err
	}
	if path == dirPath {
		return nil
	}

	// Always skip .git directory (not in .gitignore but should be ignored)
	if d.IsDir() && d.Name() == ".git" {
		return fs.SkipDir
	}

	relPath, _ := filepath.Rel(dirPath, path)
	isGi, isCi := b.checkIgnored(relPath, d.IsDir(), gi, ci)

	if d.IsDir() && (isGi || isCi) {
		return fs.SkipDir
	}
	if !d.IsDir() && (isGi || isCi) {
		return nil
	}

	node := b.createFileNode(path, relPath, d, isGi, isCi)
	nodesMap[path] = node

	if p, ok := nodesMap[filepath.Dir(path)]; ok {
		p.Children = append(p.Children, node)
	}
	return nil
}

// checkIgnored checks if path matches gitignore or custom ignore
func (b *fileTreeBuilder) checkIgnored(relPath string, isDir bool, gi, ci *gitignore.GitIgnore) (isGitIgnored, isCustomIgnored bool) {
	matchPath := relPath
	if isDir && !strings.HasSuffix(matchPath, string(filepath.Separator)) {
		matchPath += string(filepath.Separator)
	}
	isGi := gi != nil && gi.MatchesPath(matchPath)
	isCi := ci != nil && ci.MatchesPath(matchPath)
	return isGi, isCi
}

// createFileNode creates a FileNode from directory entry
func (b *fileTreeBuilder) createFileNode(path, relPath string, d fs.DirEntry, isGi, isCi bool) *domain.FileNode {
	var fsize int64
	contentType := "unknown"

	if d.IsDir() {
		contentType = ""
	} else {
		if info, e := d.Info(); e == nil {
			fsize = info.Size()
		}
		// Detect content type by extension (fast, no file read)
		contentType = detectContentTypeByExt(d.Name())
	}

	return &domain.FileNode{
		Name: d.Name(), Path: path, RelPath: relPath, IsDir: d.IsDir(),
		IsGitignored: isGi, IsCustomIgnored: isCi, IsIgnored: isGi || isCi,
		Children: []*domain.FileNode{}, Size: fsize, ContentType: contentType,
	}
}

// sortChildren sorts children of all nodes
func (b *fileTreeBuilder) sortChildren(nodesMap map[string]*domain.FileNode) {
	for _, n := range nodesMap {
		sort.Slice(n.Children, func(i, j int) bool {
			if n.Children[i].IsDir != n.Children[j].IsDir {
				return n.Children[i].IsDir
			}
			return strings.ToLower(n.Children[i].Name) < strings.ToLower(n.Children[j].Name)
		})
	}
}

// getCachedTree retrieves a cached tree if it's still valid
func (b *fileTreeBuilder) getCachedTree(dirPath string) []*domain.FileNode {
	// Только чтение под RLock
	b.cacheMutex.RLock()
	cached, exists := b.treeCache[dirPath]
	b.cacheMutex.RUnlock()

	if !exists {
		// Cache miss
		b.cacheMutex.Lock()
		b.cacheMisses++
		b.cacheMutex.Unlock()
		return nil
	}

	// Check if cache is still valid
	if time.Since(cached.modTime) < b.cacheDuration {
		// Cache hit - обновляем метрики под Lock
		b.cacheMutex.Lock()
		// Повторная проверка после получения Lock
		if stillCached, stillExists := b.treeCache[dirPath]; stillExists {
			if time.Since(stillCached.modTime) < b.cacheDuration {
				b.cacheAccessTimes[dirPath] = time.Now()
				b.cacheHits++
				b.cacheMutex.Unlock()
				return stillCached.nodes
			}
		}
		b.cacheMutex.Unlock()
	}

	// Cache expired - удаляем под Lock
	b.cacheMutex.Lock()
	if expiredCached, stillExists := b.treeCache[dirPath]; stillExists {
		b.cacheSize -= expiredCached.size
		delete(b.treeCache, dirPath)
		delete(b.cacheAccessTimes, dirPath)
	}
	b.cacheMisses++
	b.cacheMutex.Unlock()

	return nil
}

// setCachedTree stores a tree in the cache with LRU eviction
func (b *fileTreeBuilder) setCachedTree(dirPath string, nodes []*domain.FileNode) {
	b.cacheMutex.Lock()
	defer b.cacheMutex.Unlock()

	treeSize := b.estimateTreeSize(nodes)
	b.evictOldestCacheEntry(treeSize)

	b.treeCache[dirPath] = &cachedTree{
		nodes:   nodes,
		modTime: time.Now(),
		size:    treeSize,
	}
	b.cacheAccessTimes[dirPath] = time.Now()
	b.cacheSize += treeSize
}

// estimateTreeSize оценивает размер дерева в байтах
func (b *fileTreeBuilder) estimateTreeSize(nodes []*domain.FileNode) int64 {
	var size int64
	for _, node := range nodes {
		size += int64(len(node.Name) + len(node.Path) + len(node.RelPath) + 100) // Базовый размер структуры
		if node.Children != nil {
			size += b.estimateTreeSize(node.Children)
		}
	}
	return size
}

// evictOldestCacheEntry удаляет старые записи при превышении лимитов
// ВАЖНО: должен вызываться только под b.cacheMutex.Lock()
func (b *fileTreeBuilder) evictOldestCacheEntry(newTreeSize int64) {
	maxSize := int64(maxTreeCacheSizeMB * 1024 * 1024)

	for b.needsEviction(newTreeSize, maxSize) {
		oldestKey := b.findOldestCacheKey()
		if oldestKey == "" {
			break
		}
		b.evictCacheKey(oldestKey)
	}
}

// needsEviction checks if cache eviction is needed
func (b *fileTreeBuilder) needsEviction(newTreeSize, maxSize int64) bool {
	return len(b.treeCache) > 0 &&
		(len(b.treeCache) >= maxTreeCacheEntries || (b.cacheSize+newTreeSize) > maxSize)
}

// findOldestCacheKey finds the oldest cache entry key
func (b *fileTreeBuilder) findOldestCacheKey() string {
	var oldestKey string
	oldestTime := time.Now()

	// Clean up orphaned access times and find oldest
	for key, accessTime := range b.cacheAccessTimes {
		if _, exists := b.treeCache[key]; !exists {
			delete(b.cacheAccessTimes, key)
			continue
		}
		if accessTime.Before(oldestTime) {
			oldestTime = accessTime
			oldestKey = key
		}
	}

	// Fallback: pick any key from treeCache
	if oldestKey == "" {
		for key := range b.treeCache {
			return key
		}
	}
	return oldestKey
}

// evictCacheKey removes a specific key from cache
func (b *fileTreeBuilder) evictCacheKey(key string) {
	if cached, exists := b.treeCache[key]; exists {
		b.cacheSize -= cached.size
	}
	delete(b.treeCache, key)
	delete(b.cacheAccessTimes, key)
	b.log.Debug("Evicted file tree from cache: " + key)
}

// GetCacheStats возвращает статистику кэша
func (b *fileTreeBuilder) GetCacheStats() map[string]interface{} {
	b.cacheMutex.RLock()
	defer b.cacheMutex.RUnlock()

	return map[string]interface{}{
		"cached_trees":  len(b.treeCache),
		"cache_size_mb": b.cacheSize / (1024 * 1024),
		"cache_hits":    b.cacheHits,
		"cache_misses":  b.cacheMisses,
	}
}

// InvalidateCache clears the tree cache
func (b *fileTreeBuilder) InvalidateCache() {
	b.cacheMutex.Lock()
	defer b.cacheMutex.Unlock()

	b.treeCache = make(map[string]*cachedTree)
	b.cacheAccessTimes = make(map[string]time.Time)
	b.cacheSize = 0
	b.cacheHits = 0
	b.cacheMisses = 0
}

// InvalidateCacheForPath clears the cache for a specific path
func (b *fileTreeBuilder) InvalidateCacheForPath(path string) {
	b.cacheMutex.Lock()
	defer b.cacheMutex.Unlock()

	if cached, exists := b.treeCache[path]; exists {
		b.cacheSize -= cached.size
	}
	delete(b.treeCache, path)
	delete(b.cacheAccessTimes, path)
}

func (b *fileTreeBuilder) getGitignore(root string) *gitignore.GitIgnore {
	b.mu.RLock()
	if gi, ok := b.giCache[root]; ok {
		b.mu.RUnlock()
		return gi
	}
	b.mu.RUnlock()

	ig, err := gitignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		return nil
	}

	b.mu.Lock()
	b.giCache[root] = ig
	b.mu.Unlock()
	return ig
}

func (b *fileTreeBuilder) getCustomIgnore() *gitignore.GitIgnore {
	rules := strings.ReplaceAll(b.settingsRepo.GetCustomIgnoreRules(), "\r\n", "\n")
	trimmed := []string{}
	for _, line := range strings.Split(rules, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		trimmed = append(trimmed, line)
	}
	hash := strings.Join(trimmed, "\n")

	b.mu.RLock()
	if b.customCache != nil && b.customHash == hash {
		cc := b.customCache
		b.mu.RUnlock()
		return cc
	}
	b.mu.RUnlock()

	if len(trimmed) == 0 {
		return nil
	}
	ci := gitignore.CompileIgnoreLines(trimmed...)

	b.mu.Lock()
	b.customCache = ci
	b.customHash = hash
	b.mu.Unlock()
	return ci
}
