package analyzers

import (
	"context"
	"crypto/md5" //nolint:gosec // MD5 used for content hashing, not security
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/domain/analysis"
	"strings"
	"sync"
	"time"

	_ "modernc.org/sqlite" // SQLite driver
)

// CachedSymbolIndex wraps SymbolIndexImpl with SQLite persistence
type CachedSymbolIndex struct {
	*SymbolIndexImpl
	db          *sql.DB
	dbPath      string
	projectRoot string
	mu          sync.Mutex
}

// NewCachedSymbolIndex creates a symbol index with SQLite caching
func NewCachedSymbolIndex(registry analysis.AnalyzerRegistry, cacheDir string) (*CachedSymbolIndex, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(cacheDir, "symbols.db")
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		return nil, err
	}

	idx := &CachedSymbolIndex{
		SymbolIndexImpl: NewSymbolIndex(registry),
		db:              db,
		dbPath:          dbPath,
	}

	if err := idx.initDB(); err != nil {
		db.Close()
		return nil, err
	}

	return idx, nil
}

func (idx *CachedSymbolIndex) initDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS files (
		path TEXT PRIMARY KEY,
		hash TEXT NOT NULL,
		indexed_at INTEGER NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS symbols (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		file_path TEXT NOT NULL,
		name TEXT NOT NULL,
		kind TEXT NOT NULL,
		language TEXT,
		start_line INTEGER,
		end_line INTEGER,
		signature TEXT,
		doc_comment TEXT,
		parent TEXT,
		extra TEXT,
		FOREIGN KEY (file_path) REFERENCES files(path) ON DELETE CASCADE
	);
	
	CREATE INDEX IF NOT EXISTS idx_symbols_name ON symbols(name);
	CREATE INDEX IF NOT EXISTS idx_symbols_file ON symbols(file_path);
	CREATE INDEX IF NOT EXISTS idx_symbols_kind ON symbols(kind);
	CREATE INDEX IF NOT EXISTS idx_symbols_name_lower ON symbols(lower(name));
	`
	_, err := idx.db.Exec(schema)
	return err
}

// IndexProject indexes project with caching
func (idx *CachedSymbolIndex) IndexProject(ctx context.Context, projectRoot string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.projectRoot = projectRoot
	idx.SymbolIndexImpl.Clear()

	cachedFiles := idx.loadCachedFileHashes()
	filesToIndex, filesToRemove := idx.scanProjectFiles(projectRoot, cachedFiles)

	idx.removeStaleFiles(filesToRemove)
	idx.indexNewFiles(ctx, projectRoot, filesToIndex)
	idx.loadFromCache()

	idx.indexed = true
	return nil
}

// loadCachedFileHashes loads file hashes from cache
func (idx *CachedSymbolIndex) loadCachedFileHashes() map[string]string {
	cachedFiles := make(map[string]string)
	rows, err := idx.db.Query("SELECT path, hash FROM files")
	if err != nil {
		return cachedFiles
	}
	defer rows.Close()
	for rows.Next() {
		var path, hash string
		if rows.Scan(&path, &hash) == nil {
			cachedFiles[path] = hash
		}
	}
	return cachedFiles
}

// scanProjectFiles scans project and determines what needs indexing
func (idx *CachedSymbolIndex) scanProjectFiles(projectRoot string, cachedFiles map[string]string) ([]string, []string) {
	var filesToIndex []string
	visited := make(map[string]bool)

	_ = filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			if info != nil && info.IsDir() && idx.shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if idx.registry.GetAnalyzer(path) == nil {
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		visited[relPath] = true

		if idx.needsReindex(path, relPath, cachedFiles) {
			filesToIndex = append(filesToIndex, relPath)
		}
		return nil
	})

	filesToRemove := make([]string, 0, len(cachedFiles))
	for path := range cachedFiles {
		if !visited[path] {
			filesToRemove = append(filesToRemove, path)
		}
	}
	return filesToIndex, filesToRemove
}

// shouldSkipDir checks if directory should be skipped
func (idx *CachedSymbolIndex) shouldSkipDir(name string) bool {
	return strings.HasPrefix(name, ".") || name == "node_modules" || name == "vendor" || name == "build" || name == "dist"
}

// needsReindex checks if file needs reindexing
func (idx *CachedSymbolIndex) needsReindex(fullPath, relPath string, cachedFiles map[string]string) bool {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return false
	}
	cachedHash, exists := cachedFiles[relPath]
	return !exists || cachedHash != hashContent(content)
}

// removeStaleFiles removes files no longer in project
func (idx *CachedSymbolIndex) removeStaleFiles(files []string) {
	for _, path := range files {
		idx.removeFileFromCache(path)
	}
}

// indexNewFiles indexes new or changed files
func (idx *CachedSymbolIndex) indexNewFiles(ctx context.Context, projectRoot string, files []string) {
	for _, relPath := range files {
		if content, err := os.ReadFile(filepath.Join(projectRoot, relPath)); err == nil {
			_ = idx.indexFileWithCache(ctx, relPath, content)
		}
	}
}

func (idx *CachedSymbolIndex) indexFileWithCache(ctx context.Context, filePath string, content []byte) error {
	analyzer := idx.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return err
	}

	hash := hashContent(content)

	// Begin transaction
	tx, err := idx.db.Begin()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Remove old symbols for this file
	_, _ = tx.Exec("DELETE FROM symbols WHERE file_path = ?", filePath)
	_, _ = tx.Exec("DELETE FROM files WHERE path = ?", filePath)

	// Insert file record
	_, _ = tx.Exec("INSERT INTO files (path, hash, indexed_at) VALUES (?, ?, ?)",
		filePath, hash, time.Now().Unix())

	// Insert symbols
	stmt, err := tx.Prepare(`
		INSERT INTO symbols (file_path, name, kind, language, start_line, end_line, signature, doc_comment, parent, extra)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, sym := range symbols {
		extra := ""
		if sym.Extra != nil {
			if b, err := json.Marshal(sym.Extra); err == nil {
				extra = string(b)
			}
		}
		_, _ = stmt.Exec(filePath, sym.Name, string(sym.Kind), sym.Language,
			sym.StartLine, sym.EndLine, sym.Signature, sym.DocComment, sym.Parent, extra)
	}

	return tx.Commit()
}

func (idx *CachedSymbolIndex) removeFileFromCache(filePath string) {
	_, _ = idx.db.Exec("DELETE FROM symbols WHERE file_path = ?", filePath)
	_, _ = idx.db.Exec("DELETE FROM files WHERE path = ?", filePath)
}

func (idx *CachedSymbolIndex) loadFromCache() {
	rows, err := idx.db.Query(`
		SELECT file_path, name, kind, language, start_line, end_line, signature, doc_comment, parent, extra
		FROM symbols
	`)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var sym analysis.Symbol
		var extra string
		if err := rows.Scan(&sym.FilePath, &sym.Name, &sym.Kind, &sym.Language,
			&sym.StartLine, &sym.EndLine, &sym.Signature, &sym.DocComment, &sym.Parent, &extra); err != nil {
			continue
		}
		sym.Line = sym.StartLine
		if extra != "" {
			_ = json.Unmarshal([]byte(extra), &sym.Extra)
		}
		idx.addSymbolLocked(sym)
	}
}

// IndexFile indexes a single file with caching
func (idx *CachedSymbolIndex) IndexFile(ctx context.Context, filePath string, content []byte) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Remove old symbols from memory index
	idx.removeSymbolsForFile(filePath)

	// Index with cache
	if err := idx.indexFileWithCache(ctx, filePath, content); err != nil {
		return err
	}

	// Add to memory index
	analyzer := idx.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil
	}

	symbols, err := analyzer.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return err
	}

	for _, sym := range symbols {
		idx.addSymbolLocked(sym)
	}
	return nil
}

func (idx *CachedSymbolIndex) removeSymbolsForFile(filePath string) {
	// Get indices of symbols to remove
	indices := idx.byFile[filePath]
	if len(indices) == 0 {
		return
	}

	// Create a set of indices to remove for O(1) lookup
	toRemove := make(map[int]bool)
	for _, i := range indices {
		toRemove[i] = true
	}

	// Remove from byName index
	for name, nameIndices := range idx.byName {
		var filtered []int
		for _, i := range nameIndices {
			if !toRemove[i] {
				filtered = append(filtered, i)
			}
		}
		if len(filtered) == 0 {
			delete(idx.byName, name)
		} else {
			idx.byName[name] = filtered
		}
	}

	// Remove from byKind index
	for kind, kindIndices := range idx.byKind {
		var filtered []int
		for _, i := range kindIndices {
			if !toRemove[i] {
				filtered = append(filtered, i)
			}
		}
		if len(filtered) == 0 {
			delete(idx.byKind, kind)
		} else {
			idx.byKind[kind] = filtered
		}
	}

	// Remove from byFile index
	delete(idx.byFile, filePath)

	// Note: We don't remove from symbols slice to avoid index invalidation
	// The symbols will be orphaned but won't be returned in queries
}

// Close closes the database connection
func (idx *CachedSymbolIndex) Close() error {
	if idx.db != nil {
		return idx.db.Close()
	}
	return nil
}

// GetCacheStats returns cache statistics
func (idx *CachedSymbolIndex) GetCacheStats() map[string]int {
	stats := idx.Stats()

	var fileCount, symbolCount int
	_ = idx.db.QueryRow("SELECT COUNT(*) FROM files").Scan(&fileCount)
	_ = idx.db.QueryRow("SELECT COUNT(*) FROM symbols").Scan(&symbolCount)

	stats["cached_files"] = fileCount
	stats["cached_symbols"] = symbolCount
	return stats
}

// InvalidateCache clears the entire cache
func (idx *CachedSymbolIndex) InvalidateCache() error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.SymbolIndexImpl.Clear()
	_, err := idx.db.Exec("DELETE FROM symbols; DELETE FROM files;")
	return err
}

func hashContent(content []byte) string {
	h := md5.Sum(content) //nolint:gosec // MD5 used for content hashing, not security
	return hex.EncodeToString(h[:])
}

// OnFileChanged handles file change events from file watcher
// This method is designed to be called when a file is created, modified, or deleted
func (idx *CachedSymbolIndex) OnFileChanged(ctx context.Context, filePath string, projectRoot string) error {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Get relative path
	relPath, err := filepath.Rel(projectRoot, filePath)
	if err != nil {
		relPath = filePath
	}

	// Check if file exists
	content, err := os.ReadFile(filePath)
	if err != nil {
		// File was deleted - remove from index
		idx.removeSymbolsForFile(relPath)
		idx.removeFileFromCache(relPath)
		return nil
	}

	// Check if we can analyze this file
	analyzer := idx.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return nil
	}

	// Check if file actually changed (compare hash)
	cachedFiles := idx.loadCachedFileHashes()
	currentHash := hashContent(content)
	if cachedHash, exists := cachedFiles[relPath]; exists && cachedHash == currentHash {
		// File hasn't changed
		return nil
	}

	// Remove old symbols from memory index
	idx.removeSymbolsForFile(relPath)

	// Index with cache
	if err := idx.indexFileWithCache(ctx, relPath, content); err != nil {
		return err
	}

	// Add to memory index
	symbols, err := analyzer.ExtractSymbols(ctx, relPath, content)
	if err != nil {
		return err
	}

	for _, sym := range symbols {
		idx.addSymbolLocked(sym)
	}

	return nil
}

// OnFileDeleted handles file deletion events
func (idx *CachedSymbolIndex) OnFileDeleted(filePath string, projectRoot string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	relPath, err := filepath.Rel(projectRoot, filePath)
	if err != nil {
		relPath = filePath
	}

	idx.removeSymbolsForFile(relPath)
	idx.removeFileFromCache(relPath)
}

// OnDirectoryChanged handles directory change events (batch reindex)
func (idx *CachedSymbolIndex) OnDirectoryChanged(ctx context.Context, dirPath string, projectRoot string) error {
	// Walk directory and reindex all files
	return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		return idx.OnFileChanged(ctx, path, projectRoot)
	})
}
