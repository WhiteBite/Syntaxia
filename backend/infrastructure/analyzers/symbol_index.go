package analyzers

import (
	"context"
	"crypto/md5" //nolint:gosec // MD5 used for content hashing, not security
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syntaxia/domain/analysis"
	"time"
)

// FileModTime tracks file modification times for incremental updates
type FileModTime struct {
	Path    string
	ModTime time.Time
}

// IncrementalUpdate represents a batch of file changes for incremental indexing
type IncrementalUpdate struct {
	AddedFiles    []string // New files to index
	ModifiedFiles []string // Files that have been modified
	DeletedFiles  []string // Files that have been deleted
}

// SymbolIndexImpl implements analysis.SymbolIndex
type SymbolIndexImpl struct {
	mu       sync.RWMutex
	symbols  []analysis.Symbol
	byName   map[string][]int
	byFile   map[string][]int
	byKind   map[analysis.SymbolKind][]int
	registry analysis.AnalyzerRegistry
	indexed  bool

	// Caching fields for one-time initialization
	indexOnce    sync.Once
	lastIndexErr error
	projectRoot  string

	// Incremental update tracking
	fileModTimes map[string]time.Time

	// Vue cross-file symbol resolution
	vueComponents map[string]*VueComponentInfo
}

// NewSymbolIndex creates a new symbol index
func NewSymbolIndex(registry analysis.AnalyzerRegistry) *SymbolIndexImpl {
	return &SymbolIndexImpl{
		symbols:       make([]analysis.Symbol, 0),
		byName:        make(map[string][]int),
		byFile:        make(map[string][]int),
		byKind:        make(map[analysis.SymbolKind][]int),
		registry:      registry,
		fileModTimes:  make(map[string]time.Time),
		vueComponents: make(map[string]*VueComponentInfo),
	}
}

// EnsureIndexed ensures the project is indexed exactly once.
// Subsequent calls return immediately with cached result.
// Use Invalidate() to force re-indexing.
func (idx *SymbolIndexImpl) EnsureIndexed(ctx context.Context, projectRoot string) error {
	idx.mu.RLock()
	needsReindex := idx.projectRoot != "" && idx.projectRoot != projectRoot
	idx.mu.RUnlock()

	if needsReindex {
		idx.Invalidate()
	}

	idx.indexOnce.Do(func() {
		idx.mu.Lock()
		idx.projectRoot = projectRoot
		idx.mu.Unlock()
		idx.lastIndexErr = idx.IndexProject(ctx, projectRoot)
	})
	return idx.lastIndexErr
}

// Invalidate resets the index, forcing re-indexing on next EnsureIndexed call.
func (idx *SymbolIndexImpl) Invalidate() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.symbols = make([]analysis.Symbol, 0)
	idx.byName = make(map[string][]int)
	idx.byFile = make(map[string][]int)
	idx.byKind = make(map[analysis.SymbolKind][]int)
	idx.indexed = false
	idx.indexOnce = sync.Once{}
	idx.lastIndexErr = nil
	idx.projectRoot = ""
	idx.fileModTimes = make(map[string]time.Time)
	idx.vueComponents = make(map[string]*VueComponentInfo)
}

// InvalidateFile removes symbols from a specific file.
func (idx *SymbolIndexImpl) InvalidateFile(filePath string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	fileIndices := idx.byFile[filePath]
	if len(fileIndices) == 0 {
		return
	}

	toRemove := make(map[int]bool)
	for _, i := range fileIndices {
		toRemove[i] = true
	}

	newSymbols := make([]analysis.Symbol, 0, len(idx.symbols)-len(toRemove))
	for i, sym := range idx.symbols {
		if !toRemove[i] {
			newSymbols = append(newSymbols, sym)
		}
	}
	idx.symbols = newSymbols

	idx.rebuildIndices()
}

// rebuildIndices rebuilds all indices from symbols slice.
func (idx *SymbolIndexImpl) rebuildIndices() {
	idx.byName = make(map[string][]int)
	idx.byFile = make(map[string][]int)
	idx.byKind = make(map[analysis.SymbolKind][]int)
	for i, sym := range idx.symbols {
		nameLower := strings.ToLower(sym.Name)
		idx.byName[nameLower] = append(idx.byName[nameLower], i)
		idx.byFile[sym.FilePath] = append(idx.byFile[sym.FilePath], i)
		idx.byKind[sym.Kind] = append(idx.byKind[sym.Kind], i)
	}
}

func (idx *SymbolIndexImpl) Clear() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.symbols = make([]analysis.Symbol, 0)
	idx.byName = make(map[string][]int)
	idx.byFile = make(map[string][]int)
	idx.byKind = make(map[analysis.SymbolKind][]int)
	idx.indexed = false
	idx.fileModTimes = make(map[string]time.Time)
	idx.vueComponents = make(map[string]*VueComponentInfo)
}

func (idx *SymbolIndexImpl) IndexProject(ctx context.Context, projectRoot string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.symbols = make([]analysis.Symbol, 0)
	idx.byName = make(map[string][]int)
	idx.byFile = make(map[string][]int)
	idx.byKind = make(map[analysis.SymbolKind][]int)
	idx.indexed = false

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		analyzer := idx.registry.GetAnalyzer(path)
		if analyzer == nil {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		symbols, err := analyzer.ExtractSymbols(ctx, relPath, content)
		if err != nil {
			return nil
		}

		for _, sym := range symbols {
			idx.addSymbolLocked(sym)
		}
		return nil
	})

	idx.indexed = true
	return err
}

func (idx *SymbolIndexImpl) IndexFile(ctx context.Context, filePath string, content []byte) error {
	analyzer := idx.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()
	for _, sym := range symbols {
		idx.addSymbolLocked(sym)
	}
	return nil
}

func (idx *SymbolIndexImpl) addSymbolLocked(sym analysis.Symbol) {
	i := len(idx.symbols)
	idx.symbols = append(idx.symbols, sym)
	nameLower := strings.ToLower(sym.Name)
	idx.byName[nameLower] = append(idx.byName[nameLower], i)
	idx.byFile[sym.FilePath] = append(idx.byFile[sym.FilePath], i)
	idx.byKind[sym.Kind] = append(idx.byKind[sym.Kind], i)
}

func (idx *SymbolIndexImpl) SearchByName(query string) []analysis.Symbol {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	query = strings.ToLower(query)
	var results []analysis.Symbol
	for name, indices := range idx.byName {
		if strings.Contains(name, query) {
			for _, i := range indices {
				results = append(results, idx.symbols[i])
			}
		}
	}
	return results
}

func (idx *SymbolIndexImpl) FindByExactName(name string) []analysis.Symbol {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	indices := idx.byName[strings.ToLower(name)]
	results := make([]analysis.Symbol, len(indices))
	for i, j := range indices {
		results[i] = idx.symbols[j]
	}
	return results
}

func (idx *SymbolIndexImpl) GetSymbolsInFile(filePath string) []analysis.Symbol {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	indices := idx.byFile[filePath]
	results := make([]analysis.Symbol, len(indices))
	for i, j := range indices {
		results[i] = idx.symbols[j]
	}
	return results
}

func (idx *SymbolIndexImpl) GetSymbolsByKind(kind analysis.SymbolKind) []analysis.Symbol {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	indices := idx.byKind[kind]
	results := make([]analysis.Symbol, len(indices))
	for i, j := range indices {
		results[i] = idx.symbols[j]
	}
	return results
}

func (idx *SymbolIndexImpl) FindDefinition(name string, kind analysis.SymbolKind) *analysis.Symbol {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	indices := idx.byName[strings.ToLower(name)]
	for _, i := range indices {
		sym := idx.symbols[i]
		if kind == "" || sym.Kind == kind {
			return &sym
		}
	}
	return nil
}

func (idx *SymbolIndexImpl) Stats() map[string]int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	stats := map[string]int{
		"total_symbols": len(idx.symbols),
		"unique_names":  len(idx.byName),
		"files":         len(idx.byFile),
	}
	for kind, indices := range idx.byKind {
		stats[string(kind)] = len(indices)
	}
	return stats
}

func (idx *SymbolIndexImpl) IsIndexed() bool {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return idx.indexed
}

// IncrementalUpdate performs incremental indexing - only re-indexes changed files.
// Returns the number of files that were re-indexed.
func (idx *SymbolIndexImpl) IncrementalUpdate(ctx context.Context, projectRoot string) (int, error) {
	idx.mu.Lock()
	if idx.projectRoot == "" {
		idx.projectRoot = projectRoot
	}
	idx.mu.Unlock()

	changedFiles, deletedFiles := idx.detectChanges(projectRoot)

	for _, relPath := range deletedFiles {
		idx.InvalidateFile(relPath)
		idx.mu.Lock()
		delete(idx.fileModTimes, relPath)
		delete(idx.vueComponents, relPath)
		idx.mu.Unlock()
	}

	for _, relPath := range changedFiles {
		fullPath := filepath.Join(projectRoot, relPath)
		_ = idx.reindexFile(ctx, fullPath, relPath)
	}

	return len(changedFiles) + len(deletedFiles), nil
}

// UpdateIncremental applies a batch of file changes to the index.
// This is useful for file watcher integration where changes are batched.
func (idx *SymbolIndexImpl) UpdateIncremental(ctx context.Context, update IncrementalUpdate) error {
	idx.mu.RLock()
	projectRoot := idx.projectRoot
	idx.mu.RUnlock()

	// Process deleted files first
	for _, relPath := range update.DeletedFiles {
		idx.InvalidateFile(relPath)
		idx.mu.Lock()
		delete(idx.fileModTimes, relPath)
		delete(idx.vueComponents, relPath)
		idx.mu.Unlock()
	}

	// Process modified files (remove old symbols, then reindex)
	for _, relPath := range update.ModifiedFiles {
		fullPath := filepath.Join(projectRoot, relPath)
		if err := idx.reindexFile(ctx, fullPath, relPath); err != nil {
			continue // Skip files that fail to reindex
		}
	}

	// Process added files
	for _, relPath := range update.AddedFiles {
		fullPath := filepath.Join(projectRoot, relPath)
		if err := idx.reindexFile(ctx, fullPath, relPath); err != nil {
			continue // Skip files that fail to index
		}
	}

	return nil
}

// GetFileHash returns the content hash for a file (MD5).
// Returns empty string if file is not tracked.
func (idx *SymbolIndexImpl) GetFileHash(path string) string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if _, exists := idx.fileModTimes[path]; !exists {
		return ""
	}

	// Read file and compute hash
	fullPath := filepath.Join(idx.projectRoot, path)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return ""
	}

	return computeFileHash(content)
}

// SetFileHash is a no-op for SymbolIndexImpl as it uses mod times.
// Provided for interface compatibility with cached implementations.
func (idx *SymbolIndexImpl) SetFileHash(path, hash string) {
	// SymbolIndexImpl uses mod times, not hashes
	// This method exists for interface compatibility
}

// computeFileHash computes MD5 hash of content.
func computeFileHash(content []byte) string {
	h := md5.Sum(content) //nolint:gosec // MD5 used for content hashing, not security
	return hex.EncodeToString(h[:])
}

// detectChanges finds changed and deleted files.
func (idx *SymbolIndexImpl) detectChanges(projectRoot string) (changed, deleted []string) {
	_ = filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() && shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if idx.registry.GetAnalyzer(path) == nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)

		idx.mu.RLock()
		lastMod, exists := idx.fileModTimes[relPath]
		idx.mu.RUnlock()

		if !exists || info.ModTime().After(lastMod) {
			changed = append(changed, relPath)
		}
		return nil
	})

	idx.mu.RLock()
	trackedFiles := make(map[string]bool)
	for f := range idx.fileModTimes {
		trackedFiles[f] = true
	}
	idx.mu.RUnlock()

	for relPath := range trackedFiles {
		fullPath := filepath.Join(projectRoot, relPath)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			deleted = append(deleted, relPath)
		}
	}

	return changed, deleted
}

// reindexFile re-indexes a single file, removing old symbols first.
func (idx *SymbolIndexImpl) reindexFile(ctx context.Context, fullPath, relPath string) error {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	analyzer := idx.registry.GetAnalyzer(fullPath)
	if analyzer == nil {
		return nil
	}

	idx.InvalidateFile(relPath)

	symbols, err := analyzer.ExtractSymbols(ctx, relPath, content)
	if err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	for _, sym := range symbols {
		idx.addSymbolLocked(sym)
	}

	idx.fileModTimes[relPath] = info.ModTime()

	if strings.HasSuffix(relPath, ".vue") {
		idx.extractVueComponentInfo(relPath, content)
	}

	return nil
}

// UpdateFile updates a single file in the index (for file watcher integration).
func (idx *SymbolIndexImpl) UpdateFile(ctx context.Context, filePath string, content []byte) error {
	idx.InvalidateFile(filePath)

	analyzer := idx.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return err
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	for _, sym := range symbols {
		idx.addSymbolLocked(sym)
	}

	idx.fileModTimes[filePath] = time.Now()

	if strings.HasSuffix(filePath, ".vue") {
		idx.extractVueComponentInfo(filePath, content)
	}

	return nil
}

// RemoveFile removes a file from the index.
func (idx *SymbolIndexImpl) RemoveFile(filePath string) {
	idx.InvalidateFile(filePath)
	idx.mu.Lock()
	delete(idx.fileModTimes, filePath)
	delete(idx.vueComponents, filePath)
	idx.mu.Unlock()
}

// GetChangedFiles returns list of files that have changed since last index.
func (idx *SymbolIndexImpl) GetChangedFiles(projectRoot string) ([]string, error) {
	var changed []string

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() && shouldSkipDirectory(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}

		if idx.registry.GetAnalyzer(path) == nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)

		idx.mu.RLock()
		lastMod, exists := idx.fileModTimes[relPath]
		idx.mu.RUnlock()

		if !exists || info.ModTime().After(lastMod) {
			changed = append(changed, relPath)
		}
		return nil
	})

	return changed, err
}

// shouldSkipDirectory checks if a directory should be skipped during indexing.
func shouldSkipDirectory(name string) bool {
	skipDirs := map[string]bool{
		"node_modules": true,
		"vendor":       true,
		"build":        true,
		"dist":         true,
		".git":         true,
		".idea":        true,
		".vscode":      true,
		"__pycache__":  true,
		".next":        true,
		".nuxt":        true,
	}
	return strings.HasPrefix(name, ".") || skipDirs[name]
}
