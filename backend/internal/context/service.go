package context

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"syntaxia/domain"
	"sort"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

// Service handles all context management operations with memory-safe streaming by default
type Service struct {
	fileReader   domain.FileContentReader
	tokenCounter TokenCounter
	eventBus     domain.EventBus
	logger       domain.Logger
	contextDir   string

	// Streaming support with RWMutex for concurrent reads
	streams   map[string]*Stream
	streamsMu sync.RWMutex

	// Memory limits (reduced for safety)
	defaultMaxMemoryMB int
	defaultMaxTokens   int

	// Cleanup tracking
	lastCleanup time.Time

	// Worker pool for file scanning (fixed goroutine count)
	workerCount int

	// Shutdown coordination
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
	wg           sync.WaitGroup

	// Metrics for monitoring
	activeOperations int64
	totalOperations  int64
	totalBytesRead   int64
}

// TokenCounter interface for token estimation
type TokenCounter interface {
	CountTokens(text string) int
}

// BuildOptions controls how context is built
type BuildOptions struct {
	MaxTokens            int          `json:"maxTokens,omitempty"`
	MaxMemoryMB          int          `json:"maxMemoryMB,omitempty"`
	StripComments        bool         `json:"stripComments,omitempty"`
	IncludeManifest      bool         `json:"includeManifest,omitempty"`
	IncludeLineNumbers   bool         `json:"includeLineNumbers,omitempty"`
	ForceStream          bool         `json:"forceStream,omitempty"`
	EnableProgressEvents bool         `json:"enableProgressEvents,omitempty"`
	OutputFormat         OutputFormat `json:"outputFormat,omitempty"`
	ExcludeTests         bool         `json:"excludeTests,omitempty"`
	CollapseEmptyLines   bool         `json:"collapseEmptyLines,omitempty"`
	StripLicense         bool         `json:"stripLicense,omitempty"`
	CompactDataFiles     bool         `json:"compactDataFiles,omitempty"`
	SkeletonMode         bool         `json:"skeletonMode,omitempty"`
	TrimWhitespace       bool         `json:"trimWhitespace,omitempty"`
}

// Context is an alias for domain.Context used internally
type Context = domain.Context

// NewService creates a new unified context service
func NewService(
	fileReader domain.FileContentReader,
	tokenCounter TokenCounter,
	eventBus domain.EventBus,
	logger domain.Logger,
) (*Service, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}
	contextDir := filepath.Join(homeDir, ".syntaxia", "contexts")
	if err := os.MkdirAll(contextDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create context directory: %w", err)
	}

	workerCount := runtime.NumCPU()
	if workerCount > 16 {
		workerCount = 16
	}
	if workerCount < 2 {
		workerCount = 2
	}

	svc := &Service{
		fileReader:         fileReader,
		tokenCounter:       tokenCounter,
		eventBus:           eventBus,
		logger:             logger,
		contextDir:         contextDir,
		streams:            make(map[string]*Stream),
		defaultMaxMemoryMB: 30,
		defaultMaxTokens:   5000,
		lastCleanup:        time.Now(),
		workerCount:        workerCount,
		shutdownCh:         make(chan struct{}),
	}

	svc.wg.Add(1)
	go svc.periodicCleanup()

	return svc, nil
}

// Shutdown gracefully stops the service
func (s *Service) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down context service...")

	s.shutdownOnce.Do(func() {
		close(s.shutdownCh)
	})

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Info("Context service shutdown complete")
		return nil
	case <-ctx.Done():
		s.logger.Warning("Context service shutdown timed out")
		return ctx.Err()
	}
}

// BuildContext builds a context from project files with memory-safe streaming
func (s *Service) BuildContext(ctx context.Context, projectPath string, includedPaths []string, options *BuildOptions) (*domain.Context, error) {
	if options == nil {
		options = &BuildOptions{
			MaxMemoryMB: s.defaultMaxMemoryMB,
			MaxTokens:   s.defaultMaxTokens,
			ForceStream: true,
		}
	} else {
		if options.MaxMemoryMB <= 0 {
			options.MaxMemoryMB = s.defaultMaxMemoryMB
		}
		if options.MaxTokens <= 0 {
			options.MaxTokens = s.defaultMaxTokens
		}
		options.ForceStream = true
	}

	if err := s.validateLimits(options); err != nil {
		return nil, err
	}

	return s.buildStreamingContext(ctx, projectPath, includedPaths, options)
}

// GenerateContextAsync generates context asynchronously with progress events
func (s *Service) GenerateContextAsync(ctx context.Context, rootDir string, includedPaths []string) {
	go s.generateContextSafe(ctx, rootDir, includedPaths)
}

func (s *Service) buildStreamingContext(ctx context.Context, projectPath string, includedPaths []string, options *BuildOptions) (*domain.Context, error) {
	stream, err := s.CreateStream(ctx, projectPath, includedPaths, options)
	if err != nil {
		return nil, err
	}

	return &domain.Context{
		ID:          stream.ID,
		Name:        stream.Name,
		Description: stream.Description,
		Content:     fmt.Sprintf("STREAMING_CONTEXT:%s", stream.ID),
		Files:       stream.Files,
		CreatedAt:   stream.CreatedAt,
		UpdatedAt:   stream.UpdatedAt,
		ProjectPath: stream.ProjectPath,
		TokenCount:  stream.TokenCount,
		TotalLines:  stream.TotalLines,
		TotalChars:  stream.TotalChars,
	}, nil
}

func (s *Service) handleGenerationError(err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		s.logger.Error("Context generation timed out")
		s.emitEvent("app:error", "Context generation timed out after 30 seconds")
		s.emitEvent("syntaxiaContextGenerationTimeout", nil)
	} else {
		s.logger.Error(fmt.Sprintf("Failed to read file contents: %v", err))
		s.emitEvent("app:error", fmt.Sprintf("Context generation failed: %v", err))
		s.emitEvent("syntaxiaContextGenerationFailed", fmt.Sprintf("%v", err))
	}
}

func (s *Service) buildContextString(contents map[string]string) string {
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Manifest:\n")

	manifestPaths := make([]string, 0, len(contents))
	for path := range contents {
		manifestPaths = append(manifestPaths, path)
	}
	sort.Strings(manifestPaths)

	contextBuilder.WriteString(s.buildSimpleTree(manifestPaths))
	contextBuilder.WriteString("\n")

	for _, relPath := range manifestPaths {
		if content, exists := contents[relPath]; exists {
			contextBuilder.WriteString(fmt.Sprintf("--- File: %s ---\n", relPath))
			contextBuilder.WriteString(content)
			contextBuilder.WriteString("\n\n")
		}
	}
	return strings.TrimSpace(contextBuilder.String())
}

func (s *Service) generateContextSafe(ctx context.Context, rootDir string, includedPaths []string) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error(fmt.Sprintf("PANIC recovered in GenerateContext: %v\nStack: %s", r, debug.Stack()))
			s.emitEvent("app:error", fmt.Sprintf("Context generation failed: %v", r))
			s.emitEvent("syntaxiaContextGenerationFailed", fmt.Sprintf("%v", r))
		}
	}()

	s.emitEvent("syntaxiaContextGenerationStarted", map[string]interface{}{"fileCount": len(includedPaths), "rootDir": rootDir})

	if ctx == nil {
		ctx = context.Background()
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	s.logger.Info(fmt.Sprintf("Starting async context generation for %d files", len(includedPaths)))

	sortedPaths := make([]string, len(includedPaths))
	copy(sortedPaths, includedPaths)
	sort.Strings(sortedPaths)

	g, gctx := errgroup.WithContext(ctx)
	var contents map[string]string

	g.Go(func() error {
		var err error
		contents, err = s.fileReader.ReadContents(gctx, sortedPaths, rootDir, func(current, total int64) {
			select {
			case <-gctx.Done():
			default:
				s.emitEvent("syntaxiaContextGenerationProgress", map[string]interface{}{"current": current, "total": total})
			}
		})
		return err
	})

	if err := g.Wait(); err != nil {
		s.handleGenerationError(err)
		return
	}

	finalContext := s.buildContextString(contents)
	s.logger.Info(fmt.Sprintf("Async context generation completed. Length: %d characters", len(finalContext)))
	s.emitEvent("syntaxiaContextGenerated", finalContext)
}
