package context

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"sync/atomic"
	"time"
)

// GetContext retrieves a context by ID (backward compatibility)
func (s *Service) GetContext(ctx context.Context, contextID string) (*domain.Context, error) {
	var domainCtx domain.Context
	if err := s.readAndUnmarshalJSON(filepath.Join(s.contextDir, contextID+".json"), "context: "+contextID, &domainCtx); err != nil {
		return nil, err
	}
	return &domainCtx, nil
}

// ListProjectContexts lists all contexts for a project
func (s *Service) ListProjectContexts(ctx context.Context, projectPath string) ([]*domain.Context, error) {
	entries, err := os.ReadDir(s.contextDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read context directory: %w", err)
	}

	var contexts []*domain.Context
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		contextID := strings.TrimSuffix(entry.Name(), ".json")
		domainCtx, err := s.GetContext(ctx, contextID)
		if err != nil {
			s.logger.Warning(fmt.Sprintf("Failed to load context %s: %v", contextID, err))
			continue
		}

		if domainCtx.ProjectPath == projectPath {
			contexts = append(contexts, domainCtx)
		}
	}

	return contexts, nil
}

// DeleteContext deletes a context by ID
func (s *Service) DeleteContext(ctx context.Context, contextID string) error {
	jsonPath := filepath.Join(s.contextDir, contextID+".json")
	streamPath := filepath.Join(s.contextDir, contextID+".ctx")

	if err := os.Remove(jsonPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete JSON context: %w", err)
	}

	if err := os.Remove(streamPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete streaming context: %w", err)
	}

	s.streamsMu.Lock()
	delete(s.streams, contextID)
	s.streamsMu.Unlock()

	s.logger.Info(fmt.Sprintf("Deleted context %s", contextID))
	return nil
}

// BuildContextSummary implements domain.ContextBuilder interface
func (s *Service) BuildContextSummary(ctx context.Context, projectPath string, includedPaths []string, options *domain.ContextBuildOptions) (*domain.ContextSummary, error) {
	atomic.AddInt64(&s.activeOperations, 1)
	defer atomic.AddInt64(&s.activeOperations, -1)
	atomic.AddInt64(&s.totalOperations, 1)

	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	buildOpts := s.convertBuildOptions(options)
	domainCtx, err := s.BuildContext(ctx, projectPath, includedPaths, buildOpts)
	if err != nil {
		return nil, err
	}

	// Use original includedPaths for SelectedFiles, not expanded domainCtx.Files
	// This preserves the exact files user selected in the UI
	summary := &domain.ContextSummary{
		ID:          domainCtx.ID,
		ProjectPath: domainCtx.ProjectPath,
		FileCount:   len(domainCtx.Files),
		TotalSize:   domainCtx.TotalChars,
		LineCount:   int(domainCtx.TotalLines),
		TokenCount:  domainCtx.TokenCount,
		CreatedAt:   domainCtx.CreatedAt,
		UpdatedAt:   domainCtx.UpdatedAt,
		Status:      "ready",
		Metadata: domain.ContextMetadata{
			SelectedFiles: includedPaths,
			ProjectPath:   domainCtx.ProjectPath,
		},
	}

	if err := s.SaveContextSummary(summary); err != nil {
		s.logger.Warning(fmt.Sprintf("Failed to save context summary: %v", err))
	}

	return summary, nil
}

func (s *Service) convertBuildOptions(opts *domain.ContextBuildOptions) *BuildOptions {
	if opts == nil {
		return nil
	}

	outputFormat := OutputFormat(opts.OutputFormat)
	if outputFormat == "" {
		outputFormat = FormatXML
	}

	s.logger.Info(fmt.Sprintf("[convertBuildOptions] Input format: '%s', Using format: '%s'", opts.OutputFormat, outputFormat))

	return &BuildOptions{
		MaxTokens:            opts.MaxTokens,
		MaxMemoryMB:          opts.MaxMemoryMB,
		StripComments:        opts.StripComments,
		IncludeManifest:      opts.IncludeManifest,
		IncludeLineNumbers:   opts.IncludeLineNumbers,
		ForceStream:          true,
		EnableProgressEvents: true,
		OutputFormat:         outputFormat,
		ExcludeTests:         opts.ExcludeTests,
		CollapseEmptyLines:   opts.CollapseEmptyLines,
		StripLicense:         opts.StripLicense,
		CompactDataFiles:     opts.CompactDataFiles,
		SkeletonMode:         opts.SkeletonMode,
		TrimWhitespace:       opts.TrimWhitespace,
	}
}

// SaveContextSummary persists context metadata
func (s *Service) SaveContextSummary(summary *domain.ContextSummary) error {
	if summary == nil {
		return fmt.Errorf("context summary is nil")
	}

	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context summary: %w", err)
	}

	summaryPath := filepath.Join(s.contextDir, summary.ID+".summary.json")
	if err := os.WriteFile(summaryPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write context summary: %w", err)
	}

	return nil
}

// GetContextSummary retrieves context metadata by ID
func (s *Service) GetContextSummary(ctx context.Context, contextID string) (*domain.ContextSummary, error) {
	var summary domain.ContextSummary
	if err := s.readAndUnmarshalJSON(filepath.Join(s.contextDir, contextID+".summary.json"), "context summary: "+contextID, &summary); err != nil {
		return nil, err
	}
	return &summary, nil
}

// GetProjectContextSummaries lists all context summaries for a project
func (s *Service) GetProjectContextSummaries(ctx context.Context, projectPath string) ([]*domain.ContextSummary, error) {
	entries, err := os.ReadDir(s.contextDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*domain.ContextSummary{}, nil
		}
		return nil, fmt.Errorf("failed to read context directory: %w", err)
	}

	summaries := make([]*domain.ContextSummary, 0)
	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".summary.json") {
			continue
		}

		contextID := strings.TrimSuffix(entry.Name(), ".summary.json")
		summary, err := s.GetContextSummary(ctx, contextID)
		if err != nil {
			s.logger.Warning(fmt.Sprintf("Failed to load context summary %s: %v", contextID, err))
			continue
		}

		if summary.ProjectPath == projectPath {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

// ReadContextChunk returns a chunk of context content (memory-safe pagination)
func (s *Service) ReadContextChunk(ctx context.Context, contextID string, startLine int, lineCount int) (*domain.ContextChunk, error) {
	if lineCount <= 0 {
		lineCount = 1000
	}
	if startLine < 1 {
		startLine = 1
	}

	const maxLinesPerRequest = 10000
	if lineCount > maxLinesPerRequest {
		lineCount = maxLinesPerRequest
	}

	s.streamsMu.RLock()
	stream, exists := s.streams[contextID]
	s.streamsMu.RUnlock()

	var contextPath string
	if exists {
		contextPath = stream.contextPath
	} else {
		contextPath = filepath.Join(s.contextDir, contextID+".ctx")
		if _, err := os.Stat(contextPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("context not found: %s", contextID)
		}
	}

	file, err := os.Open(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open context file: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReaderSize(file, 64*1024)
	scanner := bufio.NewScanner(reader)

	var lines []string
	currentLine := 0
	hasMore := false

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if len(lines) >= lineCount {
			hasMore = true
			break
		}
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading context file: %w", err)
	}

	endLine := startLine + len(lines) - 1
	if len(lines) == 0 {
		endLine = startLine - 1
	}

	return &domain.ContextChunk{
		Lines:     lines,
		StartLine: startLine,
		EndLine:   endLine,
		HasMore:   hasMore,
		ChunkID:   fmt.Sprintf("%s:%d", contextID, startLine),
		ContextID: contextID,
	}, nil
}

// ReadContextContent returns full context content as string
func (s *Service) ReadContextContent(ctx context.Context, contextID string) (string, error) {
	s.streamsMu.RLock()
	stream, exists := s.streams[contextID]
	s.streamsMu.RUnlock()

	var contextPath string
	if exists {
		contextPath = stream.contextPath
	} else {
		contextPath = filepath.Join(s.contextDir, contextID+".ctx")
	}

	data, err := os.ReadFile(contextPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("context content not found: %s", contextID)
		}
		return "", fmt.Errorf("failed to read context content: %w", err)
	}

	atomic.AddInt64(&s.totalBytesRead, int64(len(data)))
	return string(data), nil
}
