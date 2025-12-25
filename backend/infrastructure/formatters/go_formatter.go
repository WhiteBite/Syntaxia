package formatters

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

// GoFormatter реализует Formatter для Go
type GoFormatter struct {
	log domain.Logger
}

// NewGoFormatter создает новый форматтер для Go
func NewGoFormatter(log domain.Logger) *GoFormatter {
	return &GoFormatter{
		log: log,
	}
}

// FormatFile форматирует Go файл
func (f *GoFormatter) FormatFile(ctx context.Context, path string) error {
	if !strings.HasSuffix(path, ".go") {
		return fmt.Errorf("not a Go file: %s", path)
	}

	// Запускаем gofmt
	cmd := exec.CommandContext(ctx, "gofmt", "-w", path)
	cmd.Dir = filepath.Dir(path)

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("gofmt failed for %s: %v, output: %s", path, err, string(output)))
		return fmt.Errorf("gofmt failed: %w", err)
	}

	// Запускаем goimports для исправления импортов
	cmd = exec.CommandContext(ctx, "goimports", "-w", path)
	cmd.Dir = filepath.Dir(path)

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("goimports failed for %s: %v, output: %s", path, err, string(output)))
		// Не возвращаем ошибку, так как gofmt уже отработал
	}

	f.log.Info(fmt.Sprintf("Formatted Go file: %s", path))
	return nil
}

// processContentWithTempFile is a helper that writes content to temp file, processes it, and returns result
func (f *GoFormatter) processContentWithTempFile(ctx context.Context, content, prefix string, processor func(ctx context.Context, path string) error) (string, error) {
	tmpFile, err := os.CreateTemp("", prefix+"*.go")
	if err != nil {
		return content, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return content, fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := processor(ctx, tmpFile.Name()); err != nil {
		return content, err
	}

	result, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return content, fmt.Errorf("failed to read processed file: %w", err)
	}

	return string(result), nil
}

// FormatContent форматирует содержимое Go кода
func (f *GoFormatter) FormatContent(ctx context.Context, content string, language string) (string, error) {
	if language != "go" {
		return content, fmt.Errorf("unsupported language: %s", language)
	}
	return f.processContentWithTempFile(ctx, content, "gofmt-", f.FormatFile)
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (f *GoFormatter) GetSupportedLanguages() []string {
	return []string{"go"}
}

// FixImports исправляет импорты в Go файле
func (f *GoFormatter) FixImports(ctx context.Context, path string) error {
	if !strings.HasSuffix(path, ".go") {
		return fmt.Errorf("not a Go file: %s", path)
	}

	// Запускаем goimports
	cmd := exec.CommandContext(ctx, "goimports", "-w", path)
	cmd.Dir = filepath.Dir(path)

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("goimports failed for %s: %v, output: %s", path, err, string(output)))
		return fmt.Errorf("goimports failed: %w", err)
	}

	f.log.Info(fmt.Sprintf("Fixed imports in Go file: %s", path))
	return nil
}

// FixImportsInContent исправляет импорты в содержимом Go кода
func (f *GoFormatter) FixImportsInContent(ctx context.Context, content string, language string) (string, error) {
	if language != "go" {
		return content, fmt.Errorf("unsupported language: %s", language)
	}
	return f.processContentWithTempFile(ctx, content, "goimports-", f.FixImports)
}
