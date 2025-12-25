package analyzers

import (
	"context"
	"os"
	"path/filepath"
	"syntaxia/domain/analysis"
	"strings"
	"sync"
)

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
}

// NewSymbolIndex creates a new symbol index
func NewSymbolIndex(registry analysis.AnalyzerRegistry) *SymbolIndexImpl {
	return &SymbolIndexImpl{
		symbols:  make([]analysis.Symbol, 0),
		byName:   make(map[string][]int),
		byFile:   make(map[string][]int),
		byKind:   make(map[analysis.SymbolKind][]int),
		registry: registry,
	}
}

// EnsureIndexed ensures the project is indexed exactly once.
// Subsequent calls return immediately with cached result.
// Use Invalidate() to force re-indexing.
func (idx *SymbolIndexImpl) EnsureIndexed(ctx context.Context, projectRoot string) error {
	// Check if project changed - need to re-index
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
	idx.indexOnce = sync.Once{} // Reset sync.Once
	idx.lastIndexErr = nil
	idx.projectRoot = ""
}

// InvalidateFile removes symbols from a specific file and marks for partial re-index.
// Useful for incremental updates when a single file changes.
func (idx *SymbolIndexImpl) InvalidateFile(filePath string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Get indices of symbols in this file
	fileIndices := idx.byFile[filePath]
	if len(fileIndices) == 0 {
		return
	}

	// Mark symbols as removed (we'll rebuild indices)
	toRemove := make(map[int]bool)
	for _, i := range fileIndices {
		toRemove[i] = true
	}

	// Rebuild symbols slice without removed ones
	newSymbols := make([]analysis.Symbol, 0, len(idx.symbols)-len(toRemove))
	oldToNew := make(map[int]int) // old index -> new index
	for i, sym := range idx.symbols {
		if !toRemove[i] {
			oldToNew[i] = len(newSymbols)
			newSymbols = append(newSymbols, sym)
		}
	}
	idx.symbols = newSymbols

	// Rebuild all indices
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
}

func (idx *SymbolIndexImpl) IndexProject(ctx context.Context, projectRoot string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Clear index while holding the lock to prevent race condition
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
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "build" || name == "dist" {
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
