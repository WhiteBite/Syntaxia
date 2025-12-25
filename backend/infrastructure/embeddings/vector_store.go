package embeddings

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"sort"
	"sync"
	"time"

	_ "modernc.org/sqlite"
)

// SQLiteVectorStore implements VectorStore using SQLite
type SQLiteVectorStore struct {
	db     *sql.DB
	mu     sync.RWMutex
	dbPath string
	log    domain.Logger
}

// NewSQLiteVectorStore creates a new SQLite-based vector store
func NewSQLiteVectorStore(dataDir string, log domain.Logger) (*SQLiteVectorStore, error) {
	// Ensure data directory exists
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "embeddings.db")

	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_synchronous=NORMAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteVectorStore{
		db:     db,
		dbPath: dbPath,
		log:    log,
	}

	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// initSchema creates the database schema
func (s *SQLiteVectorStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS embeddings (
		id TEXT PRIMARY KEY,
		project_id TEXT NOT NULL,
		file_path TEXT NOT NULL,
		content TEXT NOT NULL,
		start_line INTEGER NOT NULL,
		end_line INTEGER NOT NULL,
		chunk_type TEXT NOT NULL,
		symbol_name TEXT,
		symbol_kind TEXT,
		language TEXT NOT NULL,
		token_count INTEGER NOT NULL,
		content_hash TEXT NOT NULL,
		embedding BLOB NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_embeddings_project ON embeddings(project_id);
	CREATE INDEX IF NOT EXISTS idx_embeddings_file ON embeddings(project_id, file_path);
	CREATE INDEX IF NOT EXISTS idx_embeddings_hash ON embeddings(project_id, content_hash);
	CREATE INDEX IF NOT EXISTS idx_embeddings_language ON embeddings(project_id, language);
	CREATE INDEX IF NOT EXISTS idx_embeddings_chunk_type ON embeddings(project_id, chunk_type);
	
	CREATE TABLE IF NOT EXISTS projects (
		id TEXT PRIMARY KEY,
		root_path TEXT NOT NULL,
		last_indexed DATETIME,
		total_chunks INTEGER DEFAULT 0,
		total_files INTEGER DEFAULT 0,
		dimensions INTEGER DEFAULT 0
	);
	`

	_, err := s.db.Exec(schema)
	return err
}

// Store stores an embedded chunk
func (s *SQLiteVectorStore) Store(ctx context.Context, projectID string, chunk domain.EmbeddedChunk) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	embeddingBytes, err := encodeEmbedding(chunk.Embedding)
	if err != nil {
		return fmt.Errorf("failed to encode embedding: %w", err)
	}

	query := `
	INSERT OR REPLACE INTO embeddings 
	(id, project_id, file_path, content, start_line, end_line, chunk_type, 
	 symbol_name, symbol_kind, language, token_count, content_hash, embedding, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = s.db.ExecContext(ctx, query,
		chunk.Chunk.ID,
		projectID,
		chunk.Chunk.FilePath,
		chunk.Chunk.Content,
		chunk.Chunk.StartLine,
		chunk.Chunk.EndLine,
		string(chunk.Chunk.ChunkType),
		chunk.Chunk.SymbolName,
		chunk.Chunk.SymbolKind,
		chunk.Chunk.Language,
		chunk.Chunk.TokenCount,
		chunk.Chunk.Hash,
		embeddingBytes,
		chunk.CreatedAt,
		chunk.UpdatedAt,
	)

	return err
}

// StoreBatch stores multiple embedded chunks efficiently
func (s *SQLiteVectorStore) StoreBatch(ctx context.Context, projectID string, chunks []domain.EmbeddedChunk) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx, `
	INSERT OR REPLACE INTO embeddings 
	(id, project_id, file_path, content, start_line, end_line, chunk_type, 
	 symbol_name, symbol_kind, language, token_count, content_hash, embedding, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, chunk := range chunks {
		embeddingBytes, err := encodeEmbedding(chunk.Embedding)
		if err != nil {
			return fmt.Errorf("failed to encode embedding: %w", err)
		}

		_, err = stmt.ExecContext(ctx,
			chunk.Chunk.ID,
			projectID,
			chunk.Chunk.FilePath,
			chunk.Chunk.Content,
			chunk.Chunk.StartLine,
			chunk.Chunk.EndLine,
			string(chunk.Chunk.ChunkType),
			chunk.Chunk.SymbolName,
			chunk.Chunk.SymbolKind,
			chunk.Chunk.Language,
			chunk.Chunk.TokenCount,
			chunk.Chunk.Hash,
			embeddingBytes,
			chunk.CreatedAt,
			chunk.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert chunk: %w", err)
		}
	}

	return tx.Commit()
}

// Search performs vector similarity search using cosine similarity
func (s *SQLiteVectorStore) Search(ctx context.Context, projectID string, query domain.EmbeddingVector, topK int, minScore float32) ([]domain.SemanticSearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Load all embeddings for the project (for small-medium projects)
	// For large projects, consider using approximate nearest neighbor algorithms
	rows, err := s.db.QueryContext(ctx, `
	SELECT id, file_path, content, start_line, end_line, chunk_type, 
	       symbol_name, symbol_kind, language, token_count, content_hash, embedding
	FROM embeddings 
	WHERE project_id = ?
	`, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query embeddings: %w", err)
	}
	defer rows.Close()

	type scoredResult struct {
		chunk domain.CodeChunk
		score float32
	}

	var results []scoredResult

	for rows.Next() {
		var chunk domain.CodeChunk
		var chunkType string
		var embeddingBytes []byte
		var symbolName, symbolKind sql.NullString

		err := rows.Scan(
			&chunk.ID,
			&chunk.FilePath,
			&chunk.Content,
			&chunk.StartLine,
			&chunk.EndLine,
			&chunkType,
			&symbolName,
			&symbolKind,
			&chunk.Language,
			&chunk.TokenCount,
			&chunk.Hash,
			&embeddingBytes,
		)
		if err != nil {
			continue
		}

		chunk.ChunkType = domain.ChunkType(chunkType)
		if symbolName.Valid {
			chunk.SymbolName = symbolName.String
		}
		if symbolKind.Valid {
			chunk.SymbolKind = symbolKind.String
		}

		embedding, err := decodeEmbedding(embeddingBytes)
		if err != nil {
			continue
		}

		// Calculate cosine similarity
		score := cosineSimilarity(query, embedding)

		if score >= minScore {
			results = append(results, scoredResult{chunk: chunk, score: score})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// Take top K
	if len(results) > topK {
		results = results[:topK]
	}

	// Convert to response format
	searchResults := make([]domain.SemanticSearchResult, len(results))
	for i, r := range results {
		searchResults[i] = domain.SemanticSearchResult{
			Chunk: r.chunk,
			Score: r.score,
		}
	}

	return searchResults, nil
}

// Delete removes embeddings for a file
func (s *SQLiteVectorStore) Delete(ctx context.Context, projectID, filePath string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.ExecContext(ctx,
		"DELETE FROM embeddings WHERE project_id = ? AND file_path = ?",
		projectID, filePath)
	return err
}

// DeleteProject removes all embeddings for a project
func (s *SQLiteVectorStore) DeleteProject(ctx context.Context, projectID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, "DELETE FROM embeddings WHERE project_id = ?", projectID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", projectID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetStats returns statistics about stored embeddings
func (s *SQLiteVectorStore) GetStats(ctx context.Context, projectID string) (*domain.VectorStoreStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var stats domain.VectorStoreStats

	// Get counts
	row := s.db.QueryRowContext(ctx, `
	SELECT 
		COUNT(*) as total_chunks,
		COUNT(DISTINCT file_path) as total_files,
		COALESCE(SUM(token_count), 0) as total_tokens,
		MAX(updated_at) as last_updated
	FROM embeddings 
	WHERE project_id = ?
	`, projectID)

	var lastUpdated sql.NullTime
	err := row.Scan(&stats.TotalChunks, &stats.TotalFiles, &stats.TotalTokens, &lastUpdated)
	if err != nil {
		return nil, err
	}

	if lastUpdated.Valid {
		stats.LastUpdated = lastUpdated.Time
	}

	// Get file size
	if info, err := os.Stat(s.dbPath); err == nil {
		stats.IndexSize = info.Size()
	}

	// Get dimensions from first embedding
	var embeddingBytes []byte
	row = s.db.QueryRowContext(ctx,
		"SELECT embedding FROM embeddings WHERE project_id = ? LIMIT 1", projectID)
	if err := row.Scan(&embeddingBytes); err == nil {
		if emb, err := decodeEmbedding(embeddingBytes); err == nil {
			stats.Dimensions = len(emb)
		}
	}

	return &stats, nil
}

// GetChunkByID retrieves a specific chunk
func (s *SQLiteVectorStore) GetChunkByID(ctx context.Context, projectID, chunkID string) (*domain.EmbeddedChunk, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	row := s.db.QueryRowContext(ctx, `
	SELECT id, file_path, content, start_line, end_line, chunk_type, 
	       symbol_name, symbol_kind, language, token_count, content_hash, 
	       embedding, created_at, updated_at
	FROM embeddings 
	WHERE project_id = ? AND id = ?
	`, projectID, chunkID)

	var chunk domain.CodeChunk
	var chunkType string
	var embeddingBytes []byte
	var symbolName, symbolKind sql.NullString
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&chunk.ID,
		&chunk.FilePath,
		&chunk.Content,
		&chunk.StartLine,
		&chunk.EndLine,
		&chunkType,
		&symbolName,
		&symbolKind,
		&chunk.Language,
		&chunk.TokenCount,
		&chunk.Hash,
		&embeddingBytes,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	chunk.ChunkType = domain.ChunkType(chunkType)
	if symbolName.Valid {
		chunk.SymbolName = symbolName.String
	}
	if symbolKind.Valid {
		chunk.SymbolKind = symbolKind.String
	}

	embedding, err := decodeEmbedding(embeddingBytes)
	if err != nil {
		return nil, err
	}

	return &domain.EmbeddedChunk{
		Chunk:     chunk,
		Embedding: embedding,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// ListChunks lists all chunks for a file
func (s *SQLiteVectorStore) ListChunks(ctx context.Context, projectID, filePath string) ([]domain.EmbeddedChunk, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.QueryContext(ctx, `
	SELECT id, file_path, content, start_line, end_line, chunk_type, 
	       symbol_name, symbol_kind, language, token_count, content_hash, 
	       embedding, created_at, updated_at
	FROM embeddings 
	WHERE project_id = ? AND file_path = ?
	ORDER BY start_line
	`, projectID, filePath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chunks []domain.EmbeddedChunk

	for rows.Next() {
		var chunk domain.CodeChunk
		var chunkType string
		var embeddingBytes []byte
		var symbolName, symbolKind sql.NullString
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&chunk.ID,
			&chunk.FilePath,
			&chunk.Content,
			&chunk.StartLine,
			&chunk.EndLine,
			&chunkType,
			&symbolName,
			&symbolKind,
			&chunk.Language,
			&chunk.TokenCount,
			&chunk.Hash,
			&embeddingBytes,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			continue
		}

		chunk.ChunkType = domain.ChunkType(chunkType)
		if symbolName.Valid {
			chunk.SymbolName = symbolName.String
		}
		if symbolKind.Valid {
			chunk.SymbolKind = symbolKind.String
		}

		embedding, err := decodeEmbedding(embeddingBytes)
		if err != nil {
			continue
		}

		chunks = append(chunks, domain.EmbeddedChunk{
			Chunk:     chunk,
			Embedding: embedding,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return chunks, nil
}

// GetFileHashes returns content hashes for all chunks in a file
func (s *SQLiteVectorStore) GetFileHashes(ctx context.Context, projectID, filePath string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rows, err := s.db.QueryContext(ctx,
		"SELECT id, content_hash FROM embeddings WHERE project_id = ? AND file_path = ?",
		projectID, filePath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashes := make(map[string]string)
	for rows.Next() {
		var id, hash string
		if err := rows.Scan(&id, &hash); err == nil {
			hashes[id] = hash
		}
	}

	return hashes, nil
}

// Close closes the database connection
func (s *SQLiteVectorStore) Close() error {
	return s.db.Close()
}

// Helper functions

func encodeEmbedding(embedding domain.EmbeddingVector) ([]byte, error) {
	return json.Marshal(embedding)
}

func decodeEmbedding(data []byte) (domain.EmbeddingVector, error) {
	var embedding domain.EmbeddingVector
	err := json.Unmarshal(data, &embedding)
	return embedding, err
}

func cosineSimilarity(a, b domain.EmbeddingVector) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct, normA, normB float64

	for i := range a {
		dotProduct += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return float32(dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)))
}
