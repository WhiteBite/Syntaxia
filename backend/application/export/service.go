// Package export provides context export functionality.
package export

import (
	"context"
	"encoding/base64"
	"fmt"
	"syntaxia/domain"
	"strings"
)

const (
	maxInMemorySize = 50 * 1024 * 1024 // 50MB
)

// Service handles context export operations.
type Service struct {
	contextSplitter  domain.ContextSplitter
	contextFormatter domain.ContextFormatter
	log              domain.Logger
	pdf              domain.PDFGenerator
	archiver         domain.Archiver
	tempFileProvider domain.TempFileProvider
	pathProvider     domain.PathProvider
	fileSystemWriter domain.FileSystemWriter
	fileStatProvider domain.FileStatProvider
}

// NewService creates a new export service.
func NewService(
	log domain.Logger,
	splitter domain.ContextSplitter,
	formatter domain.ContextFormatter,
	pdf domain.PDFGenerator,
	arch domain.Archiver,
	tempFileProvider domain.TempFileProvider,
	pathProvider domain.PathProvider,
	fileSystemWriter domain.FileSystemWriter,
	fileStatProvider domain.FileStatProvider,
) *Service {
	return &Service{
		contextSplitter:  splitter,
		contextFormatter: formatter,
		log:              log,
		pdf:              pdf,
		archiver:         arch,
		tempFileProvider: tempFileProvider,
		pathProvider:     pathProvider,
		fileSystemWriter: fileSystemWriter,
		fileStatProvider: fileStatProvider,
	}
}

// approxTokens provides rough token count estimation (~ quarter of rune count)
func approxTokens(s string) int { return len([]rune(s)) / 4 }

// exportClipboard handles clipboard export mode
func (s *Service) exportClipboard(settings domain.ExportSettings) (domain.ExportResult, error) {
	format := settings.ExportFormat
	if format == "" {
		format = "manifest"
	}
	out, err := s.contextFormatter.Format(format, settings.Context, domain.ContextFormatOptions{
		StripComments:   settings.StripComments,
		IncludeManifest: settings.IncludeManifest,
	})
	if err != nil {
		return domain.ExportResult{}, fmt.Errorf("failed to build clipboard context: %w", err)
	}
	return domain.ExportResult{Mode: settings.Mode, Text: out}, nil
}

// exportAISmallPDF generates small in-memory PDF for AI export
func (s *Service) exportAISmallPDF(settings domain.ExportSettings, content string) (domain.ExportResult, error) {
	pdfBytes, err := s.pdf.Generate(content, domain.PDFOptions{})
	if err != nil {
		return domain.ExportResult{}, fmt.Errorf("failed to generate AI PDF: %w", err)
	}
	return domain.ExportResult{
		Mode:       settings.Mode,
		FileName:   "context-ai.pdf",
		DataBase64: base64.StdEncoding.EncodeToString(pdfBytes),
		SizeBytes:  int64(len(pdfBytes)),
	}, nil
}

// exportAILargePDF generates large PDF to file for AI export
func (s *Service) exportAILargePDF(settings domain.ExportSettings, chunks []string) (domain.ExportResult, error) {
	tempDir, err := s.tempFileProvider.MkdirTemp("", "syntaxia-export-*")
	if err != nil {
		return domain.ExportResult{}, fmt.Errorf("failed to create temp dir: %w", err)
	}

	var outputPath, fileName string

	if len(chunks) == 1 {
		fileName = "context-ai.pdf"
		outputPath = s.pathProvider.Join(tempDir, fileName)
		if err := s.pdf.WriteAtomic(chunks[0], domain.PDFOptions{}, outputPath); err != nil {
			_ = s.fileSystemWriter.RemoveAll(tempDir)
			return domain.ExportResult{}, fmt.Errorf("failed to generate AI PDF: %w", err)
		}
	} else {
		files := make(map[string][]byte, len(chunks))
		for i, chunk := range chunks {
			b, err := s.pdf.Generate(chunk, domain.PDFOptions{})
			if err != nil {
				_ = s.fileSystemWriter.RemoveAll(tempDir)
				return domain.ExportResult{}, fmt.Errorf("failed to generate PDF chunk %d: %w", i+1, err)
			}
			files[fmt.Sprintf("context-ai-part-%02d.pdf", i+1)] = b
		}
		fileName = "context-ai.zip"
		outputPath = s.pathProvider.Join(tempDir, fileName)
		if err := s.archiver.ZipFilesAtomic(files, outputPath); err != nil {
			_ = s.fileSystemWriter.RemoveAll(tempDir)
			return domain.ExportResult{}, fmt.Errorf("failed to create ZIP: %w", err)
		}
	}

	fi, err := s.fileStatProvider.Stat(outputPath)
	if err != nil {
		_ = s.fileSystemWriter.RemoveAll(tempDir)
		return domain.ExportResult{}, fmt.Errorf("failed to stat output file: %w", err)
	}
	return domain.ExportResult{Mode: settings.Mode, FileName: fileName, FilePath: outputPath, IsLarge: true, SizeBytes: fi.Size()}, nil
}

// exportAI handles AI export mode
func (s *Service) exportAI(settings domain.ExportSettings) (domain.ExportResult, error) {
	var chunks []string
	if settings.EnableAutoSplit {
		var err error
		chunks, err = s.contextSplitter.SplitContext(settings.Context, domain.SplitSettings{
			MaxTokensPerChunk: settings.MaxTokensPerChunk,
			OverlapTokens:     settings.OverlapTokens,
			SplitStrategy:     settings.SplitStrategy,
		})
		if err != nil {
			return domain.ExportResult{}, fmt.Errorf("failed to split context for AI export: %w", err)
		}
	} else {
		if totalTokens := approxTokens(settings.Context); totalTokens > settings.TokenLimit && settings.TokenLimit > 0 {
			s.log.Warning(fmt.Sprintf("Context (%d tokens) exceeds limit (%d), exporting as single PDF", totalTokens, settings.TokenLimit))
		}
		chunks = []string{settings.Context}
	}

	estimatedSize := int64(len(settings.Context) * 2)
	if len(chunks) == 1 && estimatedSize < maxInMemorySize {
		return s.exportAISmallPDF(settings, chunks[0])
	}
	return s.exportAILargePDF(settings, chunks)
}

// exportHuman handles human-readable export mode
func (s *Service) exportHuman(settings domain.ExportSettings) (domain.ExportResult, error) {
	opts := domain.PDFOptions{
		Dark:        strings.EqualFold(settings.Theme, "Dark"),
		LineNumbers: settings.IncludeLineNumbers,
		PageNumbers: settings.IncludePageNumbers,
	}

	if estimatedSize := int64(len(settings.Context) * 2); estimatedSize < maxInMemorySize {
		pdfBytes, err := s.pdf.Generate(settings.Context, opts)
		if err != nil {
			return domain.ExportResult{}, fmt.Errorf("failed to generate human-readable PDF: %w", err)
		}
		return domain.ExportResult{
			Mode: settings.Mode, FileName: "context-human.pdf",
			DataBase64: base64.StdEncoding.EncodeToString(pdfBytes), SizeBytes: int64(len(pdfBytes)),
		}, nil
	}

	tempDir, err := s.tempFileProvider.MkdirTemp("", "syntaxia-export-*")
	if err != nil {
		return domain.ExportResult{}, fmt.Errorf("failed to create temp dir: %w", err)
	}

	outputPath := s.pathProvider.Join(tempDir, "context-human.pdf")
	if err := s.pdf.WriteAtomic(settings.Context, opts, outputPath); err != nil {
		_ = s.fileSystemWriter.RemoveAll(tempDir)
		return domain.ExportResult{}, fmt.Errorf("failed to generate human-readable PDF: %w", err)
	}

	fi, err := s.fileStatProvider.Stat(outputPath)
	if err != nil {
		_ = s.fileSystemWriter.RemoveAll(tempDir)
		return domain.ExportResult{}, fmt.Errorf("failed to stat output file: %w", err)
	}
	return domain.ExportResult{Mode: settings.Mode, FileName: "context-human.pdf", FilePath: outputPath, IsLarge: true, SizeBytes: fi.Size()}, nil
}

// Export exports context based on settings.
func (s *Service) Export(_ context.Context, settings domain.ExportSettings) (domain.ExportResult, error) {
	if settings.Context == "" {
		s.log.Warning("Attempted to export empty context.")
		return domain.ExportResult{}, fmt.Errorf("context is empty, nothing to export")
	}

	switch settings.Mode {
	case domain.ExportModeClipboard:
		return s.exportClipboard(settings)
	case domain.ExportModeAI:
		return s.exportAI(settings)
	case domain.ExportModeHuman:
		return s.exportHuman(settings)
	default:
		return domain.ExportResult{}, fmt.Errorf("unknown export mode: %s", settings.Mode)
	}
}

// GetExportHistory returns export history for a project.
func (s *Service) GetExportHistory(ctx context.Context, projectPath string) ([]domain.ExportHistoryItem, error) {
	s.log.Info(fmt.Sprintf("Getting export history for project: %s", projectPath))

	// For now, return an empty history as this is a new feature
	var history []domain.ExportHistoryItem

	s.log.Debug(fmt.Sprintf("Found %d export history items for project %s", len(history), projectPath))
	return history, nil
}
