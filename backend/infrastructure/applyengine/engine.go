package applyengine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

// Operation type constants
const (
	opModify = "modify"
)

// Impl реализует ApplyEngine
type Impl struct {
	log          domain.Logger
	config       *domain.ApplyEngineConfig
	formatters   map[string]domain.Formatter
	importFixers map[string]domain.ImportFixer
	backups      map[string]string // path -> backup content
}

// NewApplyEngine создает новый движок применения
func NewApplyEngine(log domain.Logger, config *domain.ApplyEngineConfig) *Impl {
	return &Impl{
		log:          log,
		config:       config,
		formatters:   make(map[string]domain.Formatter),
		importFixers: make(map[string]domain.ImportFixer),
		backups:      make(map[string]string),
	}
}

// RegisterFormatter регистрирует форматтер для языка
func (e *Impl) RegisterFormatter(language string, formatter domain.Formatter) {
	e.formatters[language] = formatter
	e.log.Info(fmt.Sprintf("Registered formatter for language: %s", language))
}

// RegisterImportFixer регистрирует исправитель импортов для языка
func (e *Impl) RegisterImportFixer(language string, fixer domain.ImportFixer) {
	e.importFixers[language] = fixer
	e.log.Info(fmt.Sprintf("Registered import fixer for language: %s", language))
}

// ApplyOperation применяет одну операцию
func (e *Impl) ApplyOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	e.log.Info(fmt.Sprintf("Applying operation %s to %s", op.ID, op.Path))

	// Валидируем операцию
	if err := e.ValidateOperation(ctx, op); err != nil {
		return &domain.ApplyResult{
			Success:     false,
			Path:        op.Path,
			OperationID: op.ID,
			Error:       err.Error(),
		}, nil
	}

	// Создаем резервную копию если нужно
	if e.config.BackupFiles {
		if err := e.createBackup(op.Path); err != nil {
			e.log.Warning(fmt.Sprintf("Failed to create backup for %s: %v", op.Path, err))
		}
	}

	var result *domain.ApplyResult
	var err error

	// Применяем операцию в зависимости от стратегии
	switch op.Strategy {
	case domain.ApplyStrategyAnchor:
		result, err = e.applyAnchorOperation(ctx, op)
	case domain.ApplyStrategyFullFile:
		result, err = e.applyFullFileOperation(ctx, op)
	case domain.ApplyStrategyAST:
		result, err = e.applyASTOperation(ctx, op)
	case domain.ApplyStrategyRecipe:
		result, err = e.applyRecipeOperation(ctx, op)
	default:
		err = fmt.Errorf("unsupported strategy: %s", op.Strategy)
	}

	if err != nil {
		return &domain.ApplyResult{
			Success:     false,
			Path:        op.Path,
			OperationID: op.ID,
			Error:       err.Error(),
		}, nil
	}

	// Применяем пост-обработку
	if result.Success {
		if err := e.postProcess(ctx, op); err != nil {
			e.log.Warning(fmt.Sprintf("Post-processing failed for %s: %v", op.Path, err))
		}
	}

	return result, nil
}

// ApplyOperations применяет несколько операций
func (e *Impl) ApplyOperations(ctx context.Context, ops []*domain.ApplyOperation) ([]*domain.ApplyResult, error) {
	results := make([]*domain.ApplyResult, 0, len(ops))

	for _, op := range ops {
		result, err := e.ApplyOperation(ctx, op)
		if err != nil {
			return results, err
		}
		results = append(results, result)

		// Если операция не удалась, останавливаемся
		if !result.Success {
			break
		}
	}

	return results, nil
}

// ValidateOperation проверяет корректность операции
func (e *Impl) ValidateOperation(ctx context.Context, op *domain.ApplyOperation) error {
	if op.ID == "" {
		return fmt.Errorf("operation ID is required")
	}

	if op.Path == "" {
		return fmt.Errorf("operation path is required")
	}

	if op.Language == "" {
		return fmt.Errorf("operation language is required")
	}

	// Проверяем существование файла для модификации
	if op.Operation == opModify {
		if _, err := os.Stat(op.Path); os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %s", op.Path)
		}
	}

	// Проверяем якоря для anchor стратегии
	if op.Strategy == domain.ApplyStrategyAnchor {
		if op.AnchorBefore == "" && op.AnchorAfter == "" {
			return fmt.Errorf("at least one anchor is required for anchor strategy")
		}

		// Валидируем hash окна если предоставлен
		if op.Hash != "" {
			if err := e.validateAnchorHash(ctx, op); err != nil {
				return fmt.Errorf("anchor hash validation failed: %w", err)
			}
		}
	}

	return nil
}

// RollbackOperation откатывает операцию
func (e *Impl) RollbackOperation(ctx context.Context, result *domain.ApplyResult) error {
	if !e.config.BackupFiles {
		return fmt.Errorf("backup not available")
	}

	backup, exists := e.backups[result.Path]
	if !exists {
		return fmt.Errorf("backup not found for %s", result.Path)
	}

	// Восстанавливаем файл из резервной копии
	if err := os.WriteFile(result.Path, []byte(backup), 0o600); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	delete(e.backups, result.Path)
	e.log.Info(fmt.Sprintf("Rolled back operation %s for %s", result.OperationID, result.Path))

	return nil
}

// applyAnchorOperation применяет операцию с якорями
func (e *Impl) applyAnchorOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	content, err := os.ReadFile(op.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	newLines := make([]string, 0, len(lines)+10)
	anchorFound := false

	for _, line := range lines {
		newLines = append(newLines, line)

		// Ищем якорь "до"
		if op.AnchorBefore != "" && strings.Contains(line, op.AnchorBefore) {
			anchorFound = true
			// Вставляем новый контент после этой строки
			newContentLines := strings.Split(op.Content, "\n")
			newLines = append(newLines, newContentLines...)
		}

		// Ищем якорь "после"
		if op.AnchorAfter != "" && strings.Contains(line, op.AnchorAfter) {
			anchorFound = true
			// Вставляем новый контент перед этой строкой
			newContentLines := strings.Split(op.Content, "\n")
			// Удаляем последнюю добавленную строку и вставляем контент
			newLines = newLines[:len(newLines)-1]
			newLines = append(newLines, newContentLines...)
			newLines = append(newLines, line)
		}
	}

	// Проверяем, что якорь был найден
	if !anchorFound {
		return &domain.ApplyResult{
			Success:     false,
			Path:        op.Path,
			OperationID: op.ID,
			Error:       "anchor not found in file",
		}, nil
	}

	newContent := strings.Join(newLines, "\n")

	// Проверяем, что файл изменился
	if newContent == string(content) {
		return &domain.ApplyResult{
			Success:     false,
			Path:        op.Path,
			OperationID: op.ID,
			Error:       "no changes applied",
		}, nil
	}

	if err := os.WriteFile(op.Path, []byte(newContent), 0o600); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &domain.ApplyResult{
		Success:      true,
		Path:         op.Path,
		OperationID:  op.ID,
		AppliedLines: len(strings.Split(op.Content, "\n")),
	}, nil
}

// applyFullFileOperation применяет операцию замены всего файла
func (e *Impl) applyFullFileOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	// Создаем директорию если нужно
	dir := filepath.Dir(op.Path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(op.Path, []byte(op.Content), 0o600); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	return &domain.ApplyResult{
		Success:      true,
		Path:         op.Path,
		OperationID:  op.ID,
		AppliedLines: len(strings.Split(op.Content, "\n")),
	}, nil
}

// applyASTOperation применяет операцию на уровне AST
func (e *Impl) applyASTOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	// Базовая реализация - используем fullFile стратегию
	// В будущем можно добавить поддержку AST трансформаций
	return e.applyFullFileOperation(ctx, op)
}

// applyRecipeOperation применяет операцию рецепта
func (e *Impl) applyRecipeOperation(ctx context.Context, op *domain.ApplyOperation) (*domain.ApplyResult, error) {
	// Базовая реализация - используем fullFile стратегию
	// В будущем можно добавить поддержку рецептов
	return e.applyFullFileOperation(ctx, op)
}

// createBackup создает резервную копию файла
func (e *Impl) createBackup(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	e.backups[path] = string(content)
	return nil
}

// postProcess выполняет пост-обработку файла
func (e *Impl) postProcess(ctx context.Context, op *domain.ApplyOperation) error {
	// Форматирование
	if e.config.AutoFormat {
		if formatter, exists := e.formatters[op.Language]; exists {
			if err := formatter.FormatFile(ctx, op.Path); err != nil {
				e.log.Warning(fmt.Sprintf("Formatting failed for %s: %v", op.Path, err))
				// Не прерываем выполнение, но логируем ошибку
			} else {
				e.log.Info(fmt.Sprintf("Formatted file: %s", op.Path))
			}
		} else {
			e.log.Warning(fmt.Sprintf("No formatter registered for language: %s", op.Language))
		}
	}

	// Исправление импортов
	if e.config.AutoFixImports {
		if fixer, exists := e.importFixers[op.Language]; exists {
			if err := fixer.FixImports(ctx, op.Path); err != nil {
				e.log.Warning(fmt.Sprintf("Import fixing failed for %s: %v", op.Path, err))
				// Не прерываем выполнение, но логируем ошибку
			} else {
				e.log.Info(fmt.Sprintf("Fixed imports in file: %s", op.Path))
			}
		} else {
			e.log.Warning(fmt.Sprintf("No import fixer registered for language: %s", op.Language))
		}
	}

	// Валидация после применения
	if e.config.ValidateAfter {
		if err := e.validateAppliedChanges(ctx, op); err != nil {
			e.log.Warning(fmt.Sprintf("Post-application validation failed for %s: %v", op.Path, err))
		}
	}

	return nil
}

// calculateHash вычисляет хеш содержимого
func (e *Impl) calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// validateAnchorHash проверяет hash окна для якоря
func (e *Impl) validateAnchorHash(ctx context.Context, op *domain.ApplyOperation) error {
	content, err := os.ReadFile(op.Path)
	if err != nil {
		return fmt.Errorf("failed to read file for hash validation: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Находим окно вокруг якоря
	var windowLines []string
	contextLines := 5 // По умолчанию 5 строк контекста

	// Ищем якорь и собираем окно
	for i, line := range lines {
		if strings.Contains(line, op.AnchorBefore) || strings.Contains(line, op.AnchorAfter) {
			start := max(0, i-contextLines)
			end := min(len(lines), i+contextLines+1)
			windowLines = lines[start:end]
			break
		}
	}

	if len(windowLines) == 0 {
		return fmt.Errorf("anchor not found in file")
	}

	// Вычисляем hash окна
	windowContent := strings.Join(windowLines, "\n")
	calculatedHash := e.calculateHash(windowContent)

	if calculatedHash != op.Hash {
		return fmt.Errorf("anchor hash mismatch: expected %s, got %s", op.Hash, calculatedHash)
	}

	return nil
}

// validateAppliedChanges проверяет примененные изменения
func (e *Impl) validateAppliedChanges(ctx context.Context, op *domain.ApplyOperation) error {
	// Проверяем, что файл существует и читается
	if _, err := os.Stat(op.Path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist after application: %s", op.Path)
	}

	// Проверяем, что файл не пустой (если это не операция удаления)
	if op.Operation != "delete" {
		content, err := os.ReadFile(op.Path)
		if err != nil {
			return fmt.Errorf("failed to read file after application: %w", err)
		}

		if len(content) == 0 {
			return fmt.Errorf("file is empty after application: %s", op.Path)
		}

		// Проверяем базовую синтаксическую корректность для Go
		if op.Language == "go" {
			if err := e.validateGoSyntax(string(content)); err != nil {
				return fmt.Errorf("Go syntax validation failed: %w", err)
			}
		}
	}

	return nil
}

// validateGoSyntax проверяет базовую синтаксическую корректность Go кода
func (e *Impl) validateGoSyntax(content string) error {
	// Простая проверка на наличие базовых Go конструкций
	lines := strings.Split(content, "\n")

	// Проверяем, что есть хотя бы одна строка с кодом
	hasCode := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "//") && !strings.HasPrefix(trimmed, "/*") {
			hasCode = true
			break
		}
	}

	if !hasCode {
		return fmt.Errorf("file contains no Go code")
	}

	return nil
}

// ApplyEdit applies a single edit
func (e *Impl) ApplyEdit(ctx context.Context, edit domain.Edit) error {
	e.log.Info(fmt.Sprintf("Applying edit to file: %s", edit.FilePath))

	// Convert Edit to ApplyOperation
	op := &domain.ApplyOperation{
		ID:        edit.ID,
		Path:      edit.FilePath, // Use FilePath instead of Path for consistency with internal package
		Language:  edit.Language,
		Strategy:  e.determineStrategy(edit),
		Operation: e.determineOperation(edit),
		Content:   edit.NewContent, // Use NewContent for the operation content
	}

	// Add anchors if present
	if edit.Type == domain.EditTypeReplace {
		// For replace operations, we might need to use anchor-based strategy
		op.Strategy = domain.ApplyStrategyAnchor
		// Note: We don't have anchor information in the Edit struct, so we'll need to implement
		// a different approach for replace operations
	}

	// Apply the operation
	result, err := e.ApplyOperation(ctx, op)
	if err != nil {
		return fmt.Errorf("failed to apply operation: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("operation failed: %s", result.Error)
	}

	// Apply post-processing if configured
	if err := e.postProcess(ctx, op); err != nil {
		e.log.Warning(fmt.Sprintf("Post-processing failed for %s: %v", edit.FilePath, err))
	}

	return nil
}

// determineStrategy determines the appropriate strategy based on the edit type
func (e *Impl) determineStrategy(edit domain.Edit) domain.ApplyStrategy {
	// Default to full file strategy
	strategy := domain.ApplyStrategyFullFile

	// If we have anchor information, use anchor strategy
	if edit.Anchor != nil {
		strategy = domain.ApplyStrategyAnchor
	}

	return strategy
}

// determineOperation determines the operation type based on the edit type
func (e *Impl) determineOperation(edit domain.Edit) string {
	switch edit.Type {
	case domain.EditTypeReplace:
		return opModify
	case domain.EditTypeInsert:
		return "create"
	case domain.EditTypeDelete:
		return "delete"
	default:
		// Default to modify operation
		return opModify
	}
}
