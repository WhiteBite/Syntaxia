package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"syntaxia/domain"
	contextservice "syntaxia/internal/context"
	"sync/atomic"
	"time"
)

// ContextHandler handles all context-related operations with memory safety
type ContextHandler struct {
	log            domain.Logger
	bus            domain.EventBus
	contextService *contextservice.Service

	// Metrics
	activeBuilds   int64
	totalBuilds    int64
	totalBytesRead int64
}

// NewContextHandler creates a new context handler
func NewContextHandler(
	log domain.Logger,
	bus domain.EventBus,
	contextService *contextservice.Service,
) *ContextHandler {
	return &ContextHandler{
		log:            log,
		bus:            bus,
		contextService: contextService,
	}
}

// BuildContext builds a context and returns a summary (prevents OOM)
func (h *ContextHandler) BuildContext(ctx context.Context, projectPath string, includedPaths []string, optionsJSON string) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	atomic.AddInt64(&h.activeBuilds, 1)
	defer atomic.AddInt64(&h.activeBuilds, -1)
	atomic.AddInt64(&h.totalBuilds, 1)

	var options domain.ContextBuildOptions
	if optionsJSON != "" {
		if err := json.Unmarshal([]byte(optionsJSON), &options); err != nil {
			return "", fmt.Errorf("failed to parse options JSON: %w", err)
		}
	}

	if len(includedPaths) == 0 {
		return "", fmt.Errorf("no files provided for context build")
	}

	// Add timeout for safety
	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	summary, err := h.contextService.BuildContextSummary(ctx, projectPath, includedPaths, &options)
	if err != nil {
		return "", err
	}

	contextJSON, err := json.Marshal(summary)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context summary: %w", err)
	}

	return string(contextJSON), nil
}

// GetContextContent returns paginated context content (memory-safe)
func (h *ContextHandler) GetContextContent(ctx context.Context, contextID string, startLine, lineCount int) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	if lineCount <= 0 {
		lineCount = 1000
	}

	chunk, err := h.contextService.ReadContextChunk(ctx, contextID, startLine, lineCount)
	if err != nil {
		return "", err
	}

	chunkJSON, err := json.Marshal(chunk)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context chunk: %w", err)
	}

	return string(chunkJSON), nil
}

// GetFullContextContent returns full context content as string
func (h *ContextHandler) GetFullContextContent(ctx context.Context, contextID string) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	content, err := h.contextService.ReadContextContent(ctx, contextID)
	if err != nil {
		return "", err
	}

	atomic.AddInt64(&h.totalBytesRead, int64(len(content)))
	return content, nil
}

// GetContext retrieves context metadata by ID
func (h *ContextHandler) GetContext(ctx context.Context, contextID string) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	summary, err := h.contextService.GetContextSummary(ctx, contextID)
	if err != nil {
		return "", err
	}

	contextJSON, err := json.Marshal(summary)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context summary: %w", err)
	}

	return string(contextJSON), nil
}

// GetProjectContexts lists all contexts for a project
func (h *ContextHandler) GetProjectContexts(ctx context.Context, projectPath string) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	summaries, err := h.contextService.GetProjectContextSummaries(ctx, projectPath)
	if err != nil {
		return "", err
	}

	contextsJSON, err := json.Marshal(summaries)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context summaries: %w", err)
	}

	return string(contextsJSON), nil
}

// DeleteContext removes a context
func (h *ContextHandler) DeleteContext(ctx context.Context, contextID string) error {
	if h.contextService == nil {
		return fmt.Errorf("context service not available")
	}

	return h.contextService.DeleteContext(ctx, contextID)
}

// GetContextLines returns a chunk of context content between lines
func (h *ContextHandler) GetContextLines(ctx context.Context, contextID string, startLine, endLine int64) (string, error) {
	if h.contextService == nil {
		return "", fmt.Errorf("context service not available")
	}

	if endLine < startLine {
		return "", fmt.Errorf("invalid line range: start=%d, end=%d", startLine, endLine)
	}

	lineCount := int(endLine-startLine) + 1
	chunk, err := h.contextService.ReadContextChunk(ctx, contextID, int(startLine), lineCount)
	if err != nil {
		return "", err
	}

	chunkJSON, err := json.Marshal(chunk)
	if err != nil {
		return "", fmt.Errorf("failed to marshal context chunk: %w", err)
	}

	return string(chunkJSON), nil
}

// CreateStreamingContext creates a streaming context (delegates to BuildContext)
func (h *ContextHandler) CreateStreamingContext(ctx context.Context, projectPath string, includedPaths []string, optionsJSON string) (string, error) {
	return h.BuildContext(ctx, projectPath, includedPaths, optionsJSON)
}

// GetStreamingContext returns context summary for streaming compatibility
func (h *ContextHandler) GetStreamingContext(ctx context.Context, contextID string) (string, error) {
	return h.GetContext(ctx, contextID)
}

// CloseStreamingContext removes a streaming context
func (h *ContextHandler) CloseStreamingContext(ctx context.Context, contextID string) error {
	return h.DeleteContext(ctx, contextID)
}

// GetMetrics returns handler metrics
func (h *ContextHandler) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"active_builds":    atomic.LoadInt64(&h.activeBuilds),
		"total_builds":     atomic.LoadInt64(&h.totalBuilds),
		"total_bytes_read": atomic.LoadInt64(&h.totalBytesRead),
	}
}
