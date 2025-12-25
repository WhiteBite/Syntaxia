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

// TypeScriptFormatter реализует Formatter для TypeScript
type TypeScriptFormatter struct {
	log domain.Logger
}

// NewTypeScriptFormatter создает новый форматтер для TypeScript
func NewTypeScriptFormatter(log domain.Logger) *TypeScriptFormatter {
	return &TypeScriptFormatter{
		log: log,
	}
}

// FormatFile форматирует TypeScript файл
func (f *TypeScriptFormatter) FormatFile(ctx context.Context, path string) error {
	if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".tsx") {
		return fmt.Errorf("not a TypeScript file: %s", path)
	}

	// Ищем package.json для определения рабочей директории
	workDir := f.findPackageJsonDir(path)
	if workDir == "" {
		workDir = filepath.Dir(path)
	}

	// Запускаем prettier
	cmd := exec.CommandContext(ctx, "npx", "prettier", "--write", path)
	cmd.Dir = workDir

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("prettier failed for %s: %v, output: %s", path, err, string(output)))
		return fmt.Errorf("prettier failed: %w", err)
	}

	// Запускаем eslint --fix
	cmd = exec.CommandContext(ctx, "npx", "eslint", "--fix", path)
	cmd.Dir = workDir

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("eslint failed for %s: %v, output: %s", path, err, string(output)))
		// Не возвращаем ошибку, так как prettier уже отработал
	}

	f.log.Info(fmt.Sprintf("Formatted TypeScript file: %s", path))
	return nil
}

// processContentWithTempFile is a helper that writes content to temp file, processes it, and returns result
func (f *TypeScriptFormatter) processContentWithTempFile(ctx context.Context, content, prefix string, processor func(ctx context.Context, path string) error) (string, error) {
	ext := ".ts"
	if strings.Contains(content, "jsx") || strings.Contains(content, "tsx") {
		ext = ".tsx"
	}

	tmpFile, err := os.CreateTemp("", prefix+"*"+ext)
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

// FormatContent форматирует содержимое TypeScript кода
func (f *TypeScriptFormatter) FormatContent(ctx context.Context, content string, language string) (string, error) {
	if language != "typescript" && language != "ts" {
		return content, fmt.Errorf("unsupported language: %s", language)
	}
	return f.processContentWithTempFile(ctx, content, "prettier-", f.FormatFile)
}

// GetSupportedLanguages возвращает поддерживаемые языки
func (f *TypeScriptFormatter) GetSupportedLanguages() []string {
	return []string{"typescript", "ts"}
}

// findPackageJsonDir ищет директорию с package.json
func (f *TypeScriptFormatter) findPackageJsonDir(filePath string) string {
	dir := filepath.Dir(filePath)

	for {
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Достигли корня
		}
		dir = parent
	}

	return ""
}

// FixImports исправляет импорты в TypeScript файле
func (f *TypeScriptFormatter) FixImports(ctx context.Context, path string) error {
	if !strings.HasSuffix(path, ".ts") && !strings.HasSuffix(path, ".tsx") {
		return fmt.Errorf("not a TypeScript file: %s", path)
	}

	// Ищем package.json для определения рабочей директории
	workDir := f.findPackageJsonDir(path)
	if workDir == "" {
		workDir = filepath.Dir(path)
	}

	// Запускаем eslint --fix для исправления импортов
	cmd := exec.CommandContext(ctx, "npx", "eslint", "--fix", path)
	cmd.Dir = workDir

	if output, err := cmd.CombinedOutput(); err != nil {
		f.log.Warning(fmt.Sprintf("eslint failed for %s: %v, output: %s", path, err, string(output)))
		return fmt.Errorf("eslint failed: %w", err)
	}

	f.log.Info(fmt.Sprintf("Fixed imports in TypeScript file: %s", path))
	return nil
}

// FixImportsInContent исправляет импорты в содержимом TypeScript кода
func (f *TypeScriptFormatter) FixImportsInContent(ctx context.Context, content string, language string) (string, error) {
	if language != "typescript" && language != "ts" {
		return content, fmt.Errorf("unsupported language: %s", language)
	}
	return f.processContentWithTempFile(ctx, content, "eslint-", f.FixImports)
}
