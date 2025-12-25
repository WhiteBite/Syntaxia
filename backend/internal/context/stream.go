package context

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Stream represents a memory-safe streaming context
type Stream struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Files       []string  `json:"files"`
	ProjectPath string    `json:"projectPath"`
	TotalLines  int64     `json:"totalLines"`
	TotalChars  int64     `json:"totalChars"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	TokenCount  int       `json:"tokenCount"`
	contextPath string    `json:"-"`
}

// LineRange represents a range of lines from a streaming context
type LineRange struct {
	StartLine int64    `json:"startLine"`
	EndLine   int64    `json:"endLine"`
	Lines     []string `json:"lines"`
}

// streamWriteState holds state during stream writing
type streamWriteState struct {
	totalLines int64
	totalChars int64
	tokenCount int
	files      []string
}

// createProgressCallback creates a progress callback for file reading
func (s *Service) createProgressCallback(ctx context.Context, options *BuildOptions) func(int64, int64) {
	return func(current, total int64) {
		if options.EnableProgressEvents && s.eventBus != nil {
			select {
			case <-ctx.Done():
			default:
				s.eventBus.Emit("syntaxiaContextGenerationProgress", map[string]interface{}{"current": current, "total": total})
			}
		}
	}
}

// writeStreamHeader writes the manifest header if requested
func (s *Service) writeStreamHeader(writer *bufio.Writer, projectPath string, options *BuildOptions, state *streamWriteState) error {
	if !options.IncludeManifest {
		return nil
	}
	header := fmt.Sprintf("# Streaming Context\nProject Path: %s\nGenerated: %s\n\n", projectPath, time.Now().Format(time.RFC3339))
	if _, err := writer.WriteString(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}
	state.totalLines += int64(strings.Count(header, "\n"))
	state.totalChars += int64(len(header))
	return nil
}

// writeFileToStream writes a single file to the stream
func (s *Service) writeFileToStream(writer *bufio.Writer, filePath, content string, options *BuildOptions, state *streamWriteState) error {
	// Apply content optimizations
	content = s.applyContentOptimizations(content, filePath, options)

	// Add line numbers if requested
	if options.IncludeLineNumbers {
		content = addLineNumbers(content)
	}

	fileTokens := s.tokenCounter.CountTokens(content)
	state.tokenCount += fileTokens

	if options.MaxTokens > 0 && state.tokenCount > options.MaxTokens {
		return fmt.Errorf("context would exceed token limit: %d > %d", state.tokenCount, options.MaxTokens)
	}

	format := options.OutputFormat
	if format == "" {
		format = FormatXML // Default to XML - best for AI context
	}

	// Log format for first file only to avoid spam
	if len(state.files) == 0 {
		s.logger.Info(fmt.Sprintf("[writeFileToStream] Using output format: '%s', lineNumbers: %v", format, options.IncludeLineNumbers))
	}

	escapedContent := s.escapeForFormat(content, format)
	parts := []string{s.formatFileHeader(filePath, format), escapedContent, s.formatFileFooter(format)}

	for _, part := range parts {
		if _, err := writer.WriteString(part); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}
		state.totalLines += int64(strings.Count(part, "\n"))
		state.totalChars += int64(len(part))
	}
	state.totalLines++ // for content without trailing newline

	if options.MaxMemoryMB > 0 && state.totalChars > int64(options.MaxMemoryMB*1024*1024)/2 {
		return writer.Flush()
	}
	return nil
}

// CreateStream creates a memory-safe streaming context
func (s *Service) CreateStream(ctx context.Context, projectPath string, includedPaths []string, options *BuildOptions) (stream *Stream, err error) {
	if options == nil {
		options = &BuildOptions{MaxMemoryMB: s.defaultMaxMemoryMB, MaxTokens: s.defaultMaxTokens}
	}

	if err := s.validateLimits(options); err != nil {
		return nil, err
	}

	// Filter out test files if requested
	if options.ExcludeTests {
		includedPaths = s.filterTestFiles(includedPaths)
	}

	s.logger.Info(fmt.Sprintf("Creating streaming context for project: %s, files: %d", projectPath, len(includedPaths)))

	totalSize, oversizedFiles, err := s.estimateTotalSize(projectPath, includedPaths)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate file sizes: %w", err)
	}

	if err := s.checkMemoryLimit(options, totalSize, oversizedFiles); err != nil {
		return nil, err
	}

	contents, err := s.fileReader.ReadContents(ctx, includedPaths, projectPath, s.createProgressCallback(ctx, options))
	if err != nil {
		return nil, fmt.Errorf("failed to read file contents: %w", err)
	}

	contextID := fmt.Sprintf("stream_%s", uuid.New().String())
	contextPath := filepath.Join(s.contextDir, contextID+".ctx")

	file, err := os.Create(contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create context file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close context file: %w", closeErr)
		}
	}()

	writer := bufio.NewWriter(file)
	defer func() {
		if flushErr := writer.Flush(); flushErr != nil && err == nil {
			err = fmt.Errorf("failed to flush writer: %w", flushErr)
		}
	}()

	state := &streamWriteState{files: make([]string, 0, len(includedPaths))}

	if err := s.writeStreamHeader(writer, projectPath, options, state); err != nil {
		return nil, err
	}

	for _, filePath := range includedPaths {
		content, exists := contents[filePath]
		if !exists {
			s.logger.Warning(fmt.Sprintf("[CreateStream] File not found: %s", filePath))
			continue
		}

		state.files = append(state.files, filePath)
		if err := s.writeFileToStream(writer, filePath, content, options, state); err != nil {
			_ = file.Close()
			_ = os.Remove(contextPath)
			return nil, err
		}
	}

	now := time.Now()
	stream = &Stream{
		ID: contextID, Name: s.generateContextName(projectPath, state.files),
		Description: fmt.Sprintf("Streaming context with %d files from %s", len(state.files), filepath.Base(projectPath)),
		Files:       state.files, ProjectPath: projectPath, TotalLines: state.totalLines, TotalChars: state.totalChars,
		CreatedAt: now, UpdatedAt: now, TokenCount: state.tokenCount, contextPath: contextPath,
	}

	s.streamsMu.Lock()
	s.streams[contextID] = stream
	s.streamsMu.Unlock()

	s.logger.Info(fmt.Sprintf("Created streaming context %s with %d lines, %d tokens", contextID, state.totalLines, state.tokenCount))
	return stream, nil
}

// GetContextLines retrieves a range of lines from a streaming context with limits
func (s *Service) GetContextLines(ctx context.Context, contextID string, startLine, endLine int64) (*LineRange, error) {
	s.streamsMu.RLock()
	stream, exists := s.streams[contextID]
	s.streamsMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("streaming context not found: %s", contextID)
	}

	// Limit maximum lines per request
	const maxLinesPerRequest = 10000
	if endLine-startLine > maxLinesPerRequest {
		return nil, fmt.Errorf("requested line range too large: %d lines (max: %d)", endLine-startLine, maxLinesPerRequest)
	}

	file, err := os.Open(stream.contextPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open context file: %w", err)
	}
	defer file.Close()

	// Use buffered reader for better performance
	reader := bufio.NewReaderSize(file, 64*1024) // 64KB buffer
	scanner := bufio.NewScanner(reader)
	var lines []string
	var currentLine int64

	// Skip to start line
	for currentLine < startLine && scanner.Scan() {
		currentLine++
	}

	// Read lines in range
	for currentLine >= startLine && currentLine <= endLine && scanner.Scan() {
		lines = append(lines, scanner.Text())
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading context file: %w", err)
	}

	return &LineRange{
		StartLine: startLine,
		EndLine:   endLine,
		Lines:     lines,
	}, nil
}

// GetStream retrieves a streaming context by ID
func (s *Service) GetStream(ctx context.Context, contextID string) (*Stream, error) {
	s.streamsMu.RLock()
	stream, exists := s.streams[contextID]
	s.streamsMu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("streaming context not found: %s", contextID)
	}

	return stream, nil
}
