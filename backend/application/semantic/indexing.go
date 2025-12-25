package semantic

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"time"
)

// IndexProject indexes all files in a project
func (s *ServiceImpl) IndexProject(ctx context.Context, projectRoot string) error {
	projectID := generateProjectID(projectRoot)

	state, err := s.startIndexingState(projectID)
	if err != nil {
		return err
	}

	defer func() {
		s.indexingMu.Lock()
		state.InProgress = false
		s.indexingMu.Unlock()
	}()

	s.log.Info(fmt.Sprintf("Starting semantic indexing for project: %s", projectRoot))

	if s.symbolIndex != nil {
		s.log.Info("Indexing symbols for project...")
		if err := s.symbolIndex.IndexProject(ctx, projectRoot); err != nil {
			s.log.Warning(fmt.Sprintf("Symbol indexing failed (non-critical): %v", err))
		}
	}

	files, err := s.collectCodeFiles(projectRoot)
	if err != nil {
		state.Error = err
		return fmt.Errorf("failed to walk project: %w", err)
	}

	state.TotalFiles = len(files)
	s.log.Info(fmt.Sprintf("Found %d files to index", len(files)))

	// Process files in batches
	batchSize := 10
	for i := 0; i < len(files); i += batchSize {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		end := min(i+batchSize, len(files))

		batch := files[i:end]
		if err := s.indexFileBatch(ctx, projectRoot, projectID, batch); err != nil {
			s.log.Warning(fmt.Sprintf("Failed to index batch: %v", err))
		}

		state.IndexedFiles = end
		state.Progress = float64(end) / float64(len(files))
	}

	s.log.Info(fmt.Sprintf("Completed semantic indexing: %d files indexed", state.IndexedFiles))
	return nil
}

// indexFileBatch indexes a batch of files
func (s *ServiceImpl) indexFileBatch(ctx context.Context, projectRoot, projectID string, files []string) error {
	var allChunks []domain.CodeChunk

	for _, relPath := range files {
		fullPath := filepath.Join(projectRoot, relPath)

		content, err := os.ReadFile(fullPath)
		if err != nil {
			continue
		}

		// Get symbols for better chunking
		var symbols []SymbolInfoForChunking
		if s.symbolIndex != nil {
			syms := s.symbolIndex.GetSymbolsInFile(relPath)
			for _, sym := range syms {
				symbols = append(symbols, SymbolInfoForChunking{
					Name:      sym.Name,
					Kind:      string(sym.Kind),
					StartLine: sym.StartLine,
					EndLine:   sym.EndLine,
				})
			}
		}

		// Chunk the file
		chunks := s.chunkFile(relPath, content, symbols)
		allChunks = append(allChunks, chunks...)
	}

	if len(allChunks) == 0 {
		return nil
	}

	// Generate embeddings for all chunks
	texts := make([]string, len(allChunks))
	for i, chunk := range allChunks {
		texts[i] = createEmbeddingText(chunk)
	}

	// Generate embeddings with retry logic
	resp, err := s.generateEmbeddingsWithRetry(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Create embedded chunks
	now := time.Now()
	embeddedChunks := make([]domain.EmbeddedChunk, len(allChunks))
	for i, chunk := range allChunks {
		embeddedChunks[i] = domain.EmbeddedChunk{
			Chunk:     chunk,
			Embedding: resp.Embeddings[i],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	return s.vectorStore.StoreBatch(ctx, projectID, embeddedChunks)
}

// IndexFile indexes a single file
func (s *ServiceImpl) IndexFile(ctx context.Context, projectRoot string, filePath string) error {
	projectID := generateProjectID(projectRoot)
	fullPath := filepath.Join(projectRoot, filePath)

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get symbols for better chunking
	var symbols []SymbolInfoForChunking
	if s.symbolIndex != nil {
		syms := s.symbolIndex.GetSymbolsInFile(filePath)
		for _, sym := range syms {
			symbols = append(symbols, SymbolInfoForChunking{
				Name:      sym.Name,
				Kind:      string(sym.Kind),
				StartLine: sym.StartLine,
				EndLine:   sym.EndLine,
			})
		}
	}

	// Chunk the file
	chunks := s.chunkFile(filePath, content, symbols)
	if len(chunks) == 0 {
		return nil
	}

	// Generate embeddings
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = createEmbeddingText(chunk)
	}

	resp, err := s.generateEmbeddingsWithRetry(ctx, texts)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Create embedded chunks
	now := time.Now()
	embeddedChunks := make([]domain.EmbeddedChunk, len(chunks))
	for i, chunk := range chunks {
		embeddedChunks[i] = domain.EmbeddedChunk{
			Chunk:     chunk,
			Embedding: resp.Embeddings[i],
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	// Delete old chunks for this file first
	if err := s.vectorStore.Delete(ctx, projectID, filePath); err != nil {
		s.log.Warning(fmt.Sprintf("Failed to delete old chunks: %v", err))
	}

	return s.vectorStore.StoreBatch(ctx, projectID, embeddedChunks)
}

// collectCodeFiles collects all code files from project
func (s *ServiceImpl) collectCodeFiles(projectRoot string) ([]string, error) {
	var files []string
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if shouldSkipDir(info.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if isCodeFile(path) {
			relPath, _ := filepath.Rel(projectRoot, path)
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}

// chunkFile chunks a file into code chunks
func (s *ServiceImpl) chunkFile(filePath string, content []byte, symbols []SymbolInfoForChunking) []domain.CodeChunk {
	// Convert to domain.ChunkSymbolInfo format
	symbolInfos := make([]domain.ChunkSymbolInfo, 0, len(symbols))
	for _, sym := range symbols {
		symbolInfos = append(symbolInfos, domain.ChunkSymbolInfo{
			Name:      sym.Name,
			Kind:      sym.Kind,
			StartLine: sym.StartLine,
			EndLine:   sym.EndLine,
		})
	}

	return s.chunker.ChunkFile(filePath, content, symbolInfos)
}
