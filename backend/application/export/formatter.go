package export

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"time"
)

// Language constants for formatter
const (
	langGo         = "go"
	langTypeScript = "typescript"
	langJavaScript = "javascript"
)

// FormatterService provides high-level API for code formatting.
type FormatterService struct {
	log           domain.Logger
	commandRunner domain.CommandRunner
}

// NewFormatterService creates a new formatter service.
func NewFormatterService(log domain.Logger, commandRunner domain.CommandRunner) *FormatterService {
	return &FormatterService{
		log:           log,
		commandRunner: commandRunner,
	}
}

// FormatProject formats a project for the specified language.
func (s *FormatterService) FormatProject(ctx context.Context, projectPath, language string) (*domain.FormatResult, error) {
	s.log.Info(fmt.Sprintf("Formatting project: %s for language: %s", projectPath, language))

	startTime := time.Now()

	var result *domain.FormatResult
	var err error

	switch strings.ToLower(language) {
	case langGo:
		result, err = s.formatGoProject(ctx, projectPath)
	case langTypeScript, "ts", langJavaScript, "js":
		result, err = s.formatTypeScriptProject(ctx, projectPath)
	case "python", "py":
		result, err = s.formatPythonProject(ctx, projectPath)
	case "java":
		result, err = s.formatJavaProject(ctx, projectPath)
	case "cpp", "c":
		result, err = s.formatCppProject(ctx, projectPath)
	default:
		return &domain.FormatResult{
			Language:   language,
			Success:    true,
			FilesCount: 0,
			Error:      fmt.Sprintf("No formatter available for language: %s", language),
		}, nil
	}

	if err != nil {
		s.log.Warning(fmt.Sprintf("Formatting failed for %s: %v", language, err))
		return &domain.FormatResult{
			Language:   language,
			Success:    false,
			FilesCount: 0,
			Error:      err.Error(),
		}, nil
	}

	duration := time.Since(startTime).Seconds()
	s.log.Info(fmt.Sprintf("Formatting completed for %s in %.2fs", language, duration))

	return result, nil
}

// formatGoProject formats a Go project.
func (s *FormatterService) formatGoProject(ctx context.Context, projectPath string) (*domain.FormatResult, error) {
	_, err := s.commandRunner.RunCommandInDir(ctx, projectPath, "gofmt", "-w", ".")
	if err != nil {
		return nil, fmt.Errorf("gofmt failed: %w", err)
	}

	_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "goimports", "-w", ".")
	if err != nil {
		return nil, fmt.Errorf("goimports failed: %w", err)
	}

	files, err := s.countFiles(projectPath, "*.go")
	if err != nil {
		return nil, fmt.Errorf("failed to count Go files: %w", err)
	}

	return &domain.FormatResult{
		Language:   "go",
		Success:    true,
		FilesCount: files,
	}, nil
}

// formatTypeScriptProject formats a TypeScript/JavaScript project.
func (s *FormatterService) formatTypeScriptProject(ctx context.Context, projectPath string) (*domain.FormatResult, error) {
	_, err := s.commandRunner.RunCommand(ctx, "npx", "prettier", "--version")
	if err != nil {
		return &domain.FormatResult{
			Language:   "typescript",
			Success:    true,
			FilesCount: 0,
			Error:      "Prettier not available",
		}, nil
	}

	_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "npx", "prettier", "--write", ".")
	if err != nil {
		return nil, fmt.Errorf("prettier failed: %w", err)
	}

	files, err := s.countFiles(projectPath, "*.ts", "*.tsx", "*.js", "*.jsx")
	if err != nil {
		return nil, fmt.Errorf("failed to count TypeScript/JavaScript files: %w", err)
	}

	return &domain.FormatResult{
		Language:   "typescript",
		Success:    true,
		FilesCount: files,
	}, nil
}

// formatPythonProject formats a Python project.
func (s *FormatterService) formatPythonProject(ctx context.Context, projectPath string) (*domain.FormatResult, error) {
	_, err := s.commandRunner.RunCommand(ctx, "black", "--version")
	if err != nil {
		return &domain.FormatResult{
			Language:   "python",
			Success:    true,
			FilesCount: 0,
			Error:      "Black not available",
		}, nil
	}

	_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "black", ".")
	if err != nil {
		return nil, fmt.Errorf("black failed: %w", err)
	}

	files, err := s.countFiles(projectPath, "*.py")
	if err != nil {
		return nil, fmt.Errorf("failed to count Python files: %w", err)
	}

	return &domain.FormatResult{
		Language:   "python",
		Success:    true,
		FilesCount: files,
	}, nil
}

// formatJavaProject formats a Java project.
func (s *FormatterService) formatJavaProject(ctx context.Context, projectPath string) (*domain.FormatResult, error) {
	_, err := s.commandRunner.RunCommand(ctx, "java", "-version")
	if err != nil {
		return &domain.FormatResult{
			Language:   "java",
			Success:    true,
			FilesCount: 0,
			Error:      "Java not available",
		}, nil
	}

	_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "java", "-jar", "google-java-format.jar", "--replace", ".")
	if err != nil {
		s.log.Info("Google Java Format not available, skipping Java formatting")
	}

	files, err := s.countFiles(projectPath, "*.java")
	if err != nil {
		return nil, fmt.Errorf("failed to count Java files: %w", err)
	}

	return &domain.FormatResult{
		Language:   "java",
		Success:    true,
		FilesCount: files,
	}, nil
}

// formatCppProject formats a C/C++ project.
func (s *FormatterService) formatCppProject(ctx context.Context, projectPath string) (*domain.FormatResult, error) {
	_, err := s.commandRunner.RunCommand(ctx, "clang-format", "--version")
	if err != nil {
		return &domain.FormatResult{
			Language:   "cpp",
			Success:    true,
			FilesCount: 0,
			Error:      "clang-format not available",
		}, nil
	}

	_, err = s.commandRunner.RunCommandInDir(ctx, projectPath, "clang-format", "-i", "-r", ".")
	if err != nil {
		return nil, fmt.Errorf("clang-format failed: %w", err)
	}

	files, err := s.countFiles(projectPath, "*.c", "*.cpp", "*.cc", "*.h", "*.hpp")
	if err != nil {
		return nil, fmt.Errorf("failed to count C/C++ files: %w", err)
	}

	return &domain.FormatResult{
		Language:   "cpp",
		Success:    true,
		FilesCount: files,
	}, nil
}

// countFiles counts files matching patterns.
func (s *FormatterService) countFiles(projectPath string, patterns ...string) (int, error) {
	total := 0
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(projectPath, "**", pattern))
		if err != nil {
			return 0, err
		}
		total += len(matches)
	}
	return total, nil
}

// GetSupportedLanguages returns supported languages.
func (s *FormatterService) GetSupportedLanguages() []string {
	return []string{"go", "typescript", "javascript", "python", "java", "cpp", "c"}
}
