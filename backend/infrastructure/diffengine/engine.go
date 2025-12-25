package diffengine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syntaxia/domain"
	"strings"
	"time"
)

// Operation constants
const (
	opAdded    = "added"
	opModified = "modified"
	opDeleted  = "deleted"
)

// DiffEngineImpl реализует DiffEngine
type DiffEngineImpl struct {
	log        domain.Logger
	publishers map[string]domain.DiffPublisher
}

// NewDiffEngine создает новый движок diff
func NewDiffEngine(log domain.Logger) *DiffEngineImpl {
	return &DiffEngineImpl{
		log:        log,
		publishers: make(map[string]domain.DiffPublisher),
	}
}

// RegisterPublisher регистрирует издателя diff
func (e *DiffEngineImpl) RegisterPublisher(name string, publisher domain.DiffPublisher) {
	e.publishers[name] = publisher
	e.log.Info(fmt.Sprintf("Registered diff publisher: %s", name))
}

// GenerateDiff генерирует diff между двумя состояниями
func (e *DiffEngineImpl) GenerateDiff(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	e.log.Info(fmt.Sprintf("Generating diff between %s and %s in %s format", beforePath, afterPath, format))

	var content string
	var err error

	switch format {
	case domain.DiffFormatGit:
		content, err = e.generateGitDiff(ctx, beforePath, afterPath)
	case domain.DiffFormatUnified:
		content, err = e.generateUnifiedDiff(ctx, beforePath, afterPath)
	case domain.DiffFormatJSON:
		content, err = e.generateJSONDiff(ctx, beforePath, afterPath)
	case domain.DiffFormatHTML:
		content, err = e.generateHTMLDiff(ctx, beforePath, afterPath)
	default:
		return nil, fmt.Errorf("unsupported diff format: %s", format)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate diff: %w", err)
	}

	// Создаем записи diff
	entries, err := e.createDiffEntries(ctx, beforePath, afterPath)
	if err != nil {
		e.log.Warning(fmt.Sprintf("Failed to create diff entries: %v", err))
	}

	// Создаем сводку
	summary := e.createDiffSummary(entries)

	result := &domain.DiffResult{
		ID:          e.generateDiffID(beforePath, afterPath),
		Format:      format,
		Content:     content,
		Entries:     entries,
		Summary:     summary,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	e.log.Info(fmt.Sprintf("Generated diff with %d entries", len(entries)))
	return result, nil
}

// GenerateDiffFromResults генерирует diff из результатов применения правок
func (e *DiffEngineImpl) GenerateDiffFromResults(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
	e.log.Info(fmt.Sprintf("Generating diff from %d apply results", len(results)))

	var entries []*domain.DiffEntry

	for _, result := range results {
		if result.Success {
			entry := &domain.DiffEntry{
				Path:      result.Path,
				Operation: "modified",
				Metadata: map[string]string{
					"operationId":  result.OperationID,
					"appliedLines": fmt.Sprintf("%d", result.AppliedLines),
				},
			}
			entries = append(entries, entry)
		}
	}

	// Создаем сводку
	summary := e.createDiffSummary(entries)

	// Генерируем контент в нужном формате
	content, err := e.generateContentFromEntries(entries, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	result := &domain.DiffResult{
		ID:          e.generateDiffIDFromResults(results),
		Format:      format,
		Content:     content,
		Entries:     entries,
		Summary:     summary,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	return result, nil
}

// GenerateDiffFromEdits генерирует diff из Edits JSON
func (e *DiffEngineImpl) GenerateDiffFromEdits(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
	e.log.Info(fmt.Sprintf("Generating diff from %d edits", len(edits.Edits)))

	entries := make([]*domain.DiffEntry, 0, len(edits.Edits))

	for _, edit := range edits.Edits {
		operation := opModified
		switch edit.Op {
		case "create":
			operation = opAdded
		case "delete":
			operation = opDeleted
		}

		entry := &domain.DiffEntry{
			Path:       edit.Path,
			Operation:  operation,
			NewContent: edit.Content,
			Metadata: map[string]string{
				"editId":   edit.ID,
				"kind":     edit.Kind,
				"language": edit.Language,
			},
		}
		entries = append(entries, entry)
	}

	// Создаем сводку
	summary := e.createDiffSummary(entries)

	// Генерируем контент в нужном формате
	content, err := e.generateContentFromEntries(entries, format)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	result := &domain.DiffResult{
		ID:          e.generateDiffIDFromEdits(edits),
		Format:      format,
		Content:     content,
		Entries:     entries,
		Summary:     summary,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	return result, nil
}

// PublishDiff публикует diff
func (e *DiffEngineImpl) PublishDiff(ctx context.Context, diff *domain.DiffResult) error {
	e.log.Info(fmt.Sprintf("Publishing diff %s", diff.ID))

	// Сохраняем обогащенный отчет
	if err := e.saveEnrichedReport(ctx, diff); err != nil {
		e.log.Warning(fmt.Sprintf("Failed to save enriched report: %v", err))
	}

	// Публикуем во все зарегистрированные издатели
	for name, publisher := range e.publishers {
		if err := publisher.Publish(ctx, diff, name); err != nil {
			e.log.Warning(fmt.Sprintf("Failed to publish to %s: %v", name, err))
		} else {
			e.log.Info(fmt.Sprintf("Successfully published to %s", name))
		}
	}

	return nil
}

// generateGitDiff генерирует git diff
func (e *DiffEngineImpl) generateGitDiff(ctx context.Context, beforePath, afterPath string) (string, error) {
	// Используем git diff для генерации
	cmd := exec.CommandContext(ctx, "git", "diff", "--no-index", beforePath, afterPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Git diff может вернуть ненулевой код при наличии различий
		if strings.Contains(string(output), "diff --git") {
			return string(output), nil
		}
		return "", fmt.Errorf("git diff failed: %w", err)
	}

	return string(output), nil
}

// generateUnifiedDiff генерирует unified diff
func (e *DiffEngineImpl) generateUnifiedDiff(ctx context.Context, beforePath, afterPath string) (string, error) {
	// Используем diff для генерации unified diff
	cmd := exec.CommandContext(ctx, "diff", "-u", beforePath, afterPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// diff может вернуть ненулевой код при наличии различий
		if strings.Contains(string(output), "---") {
			return string(output), nil
		}
		return "", fmt.Errorf("diff failed: %w", err)
	}

	return string(output), nil
}

// generateJSONDiff генерирует JSON diff
func (e *DiffEngineImpl) generateJSONDiff(ctx context.Context, beforePath, afterPath string) (string, error) {
	// Создаем записи diff
	entries, err := e.createDiffEntries(ctx, beforePath, afterPath)
	if err != nil {
		return "", err
	}

	return e.generateContentFromEntries(entries, domain.DiffFormatJSON)
}

// generateHTMLDiff генерирует HTML diff
func (e *DiffEngineImpl) generateHTMLDiff(ctx context.Context, beforePath, afterPath string) (string, error) {
	// Создаем записи diff
	entries, err := e.createDiffEntries(ctx, beforePath, afterPath)
	if err != nil {
		return "", err
	}

	return e.generateContentFromEntries(entries, domain.DiffFormatHTML)
}

// createDiffEntries создает записи diff
func (e *DiffEngineImpl) createDiffEntries(ctx context.Context, beforePath, afterPath string) ([]*domain.DiffEntry, error) {
	var entries []*domain.DiffEntry

	// Простая реализация - сравниваем файлы
	beforeExists := false
	afterExists := false

	if _, err := os.Stat(beforePath); err == nil {
		beforeExists = true
	}

	if _, err := os.Stat(afterPath); err == nil {
		afterExists = true
	}

	if !beforeExists && afterExists {
		// Файл добавлен
		content, err := os.ReadFile(afterPath)
		if err != nil {
			return nil, err
		}

		entry := &domain.DiffEntry{
			Path:       afterPath,
			Operation:  "added",
			NewContent: string(content),
		}
		entries = append(entries, entry)
	} else if beforeExists && !afterExists {
		// Файл удален
		content, err := os.ReadFile(beforePath)
		if err != nil {
			return nil, err
		}

		entry := &domain.DiffEntry{
			Path:       beforePath,
			Operation:  "deleted",
			OldContent: string(content),
		}
		entries = append(entries, entry)
	} else if beforeExists && afterExists {
		// Файл изменен
		beforeContent, err := os.ReadFile(beforePath)
		if err != nil {
			return nil, err
		}

		afterContent, err := os.ReadFile(afterPath)
		if err != nil {
			return nil, err
		}

		if string(beforeContent) != string(afterContent) {
			entry := &domain.DiffEntry{
				Path:       afterPath,
				Operation:  "modified",
				OldContent: string(beforeContent),
				NewContent: string(afterContent),
			}
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// createDiffSummary создает сводку diff
func (e *DiffEngineImpl) createDiffSummary(entries []*domain.DiffEntry) *domain.DiffSummary {
	summary := &domain.DiffSummary{
		TotalFiles: len(entries),
	}

	for _, entry := range entries {
		switch entry.Operation {
		case "added":
			summary.AddedFiles++
			summary.AddedLines += len(strings.Split(entry.NewContent, "\n"))
		case "deleted":
			summary.DeletedFiles++
			summary.RemovedLines += len(strings.Split(entry.OldContent, "\n"))
		case "modified":
			summary.ModifiedFiles++
			oldLines := len(strings.Split(entry.OldContent, "\n"))
			newLines := len(strings.Split(entry.NewContent, "\n"))
			if newLines > oldLines {
				summary.AddedLines += newLines - oldLines
			} else if oldLines > newLines {
				summary.RemovedLines += oldLines - newLines
			}
		}
	}

	summary.TotalLines = summary.AddedLines + summary.RemovedLines

	// Добавляем обогащенную информацию
	summary.WhyView = e.createWhyView(entries)
	summary.Impact = e.createImpactAnalysis(entries)
	summary.Risk = e.createRiskAssessment(entries)

	return summary
}

// generateContentFromEntries генерирует контент из записей
func (e *DiffEngineImpl) generateContentFromEntries(entries []*domain.DiffEntry, format domain.DiffFormat) (string, error) {
	switch format {
	case domain.DiffFormatJSON:
		return e.generateJSONContent(entries)
	case domain.DiffFormatHTML:
		return e.generateHTMLContent(entries)
	default:
		return "", fmt.Errorf("unsupported format for content generation: %s", format)
	}
}

// generateJSONContent генерирует JSON контент
func (e *DiffEngineImpl) generateJSONContent(entries []*domain.DiffEntry) (string, error) {
	// Простая реализация - возвращаем JSON структуру
	content := fmt.Sprintf(`{
  "entries": %d,
  "files": [
`, len(entries))

	for i, entry := range entries {
		content += fmt.Sprintf(`    {
      "path": "%s",
      "operation": "%s"
    }`, entry.Path, entry.Operation)

		if i < len(entries)-1 {
			content += ","
		}
		content += "\n"
	}

	content += "  ]\n}"
	return content, nil
}

// createWhyView создает анализ причин изменений
func (e *DiffEngineImpl) createWhyView(entries []*domain.DiffEntry) *domain.WhyView {
	whyView := &domain.WhyView{
		Reason:       "Automatic code changes",
		Confidence:   0.8,
		RelatedFiles: make([]string, 0, len(entries)),
		Context:      make(map[string]string),
	}

	// Собираем связанные файлы
	for _, entry := range entries {
		whyView.RelatedFiles = append(whyView.RelatedFiles, entry.Path)
	}

	// Анализируем контекст изменений
	if len(entries) > 0 {
		// Определяем тип изменений
		hasNewFiles := false
		hasModifications := false
		hasDeletions := false

		for _, entry := range entries {
			switch entry.Operation {
			case "added":
				hasNewFiles = true
			case "modified":
				hasModifications = true
			case "deleted":
				hasDeletions = true
			}
		}

		if hasNewFiles {
			whyView.Reason = "Adding new functionality"
			whyView.Confidence = 0.9
		} else if hasModifications {
			whyView.Reason = "Improving existing code"
			whyView.Confidence = 0.85
		} else if hasDeletions {
			whyView.Reason = "Cleaning up code"
			whyView.Confidence = 0.75
		}

		// Добавляем контекст
		whyView.Context["totalChanges"] = fmt.Sprintf("%d", len(entries))
		whyView.Context["changeTypes"] = fmt.Sprintf("added:%d,modified:%d,deleted:%d",
			e.countOperation(entries, "added"),
			e.countOperation(entries, "modified"),
			e.countOperation(entries, "deleted"))
	}

	return whyView
}

// createImpactAnalysis создает анализ влияния изменений
func (e *DiffEngineImpl) createImpactAnalysis(entries []*domain.DiffEntry) *domain.ImpactAnalysis {
	impact := &domain.ImpactAnalysis{
		Level:           "low",
		AffectedAreas:   []string{},
		Breaking:        false,
		Performance:     "unchanged",
		Security:        "unchanged",
		Maintainability: "unchanged",
	}

	// Анализируем влияние на основе типов изменений
	totalChanges := len(entries)
	criticalFiles := 0

	for _, entry := range entries {
		// Определяем критические файлы
		if strings.Contains(entry.Path, "main.go") ||
			strings.Contains(entry.Path, "app.go") ||
			strings.Contains(entry.Path, "config") {
			criticalFiles++
		}

		// Анализируем области влияния
		if strings.Contains(entry.Path, "domain") {
			impact.AffectedAreas = append(impact.AffectedAreas, "domain")
		}
		if strings.Contains(entry.Path, "infrastructure") {
			impact.AffectedAreas = append(impact.AffectedAreas, "infrastructure")
		}
		if strings.Contains(entry.Path, "application") {
			impact.AffectedAreas = append(impact.AffectedAreas, "application")
		}
	}

	// Определяем уровень влияния
	if criticalFiles > 0 || totalChanges > 10 {
		impact.Level = "high"
	} else if totalChanges > 5 {
		impact.Level = "medium"
	}

	// Определяем breaking changes
	if e.countOperation(entries, "deleted") > 0 {
		impact.Breaking = true
	}

	// Оценка производительности, безопасности и поддерживаемости
	if e.countOperation(entries, "added") > e.countOperation(entries, "deleted") {
		impact.Maintainability = "improved"
	}

	return impact
}

// createRiskAssessment создает оценку рисков
func (e *DiffEngineImpl) createRiskAssessment(entries []*domain.DiffEntry) *domain.RiskAssessment {
	risk := &domain.RiskAssessment{
		Level:        "low",
		Risks:        []string{},
		Mitigations:  []string{},
		TestCoverage: "unknown",
		ReviewNeeded: false,
	}

	// Анализируем риски
	totalChanges := len(entries)
	hasDeletions := e.countOperation(entries, "deleted") > 0
	hasNewFiles := e.countOperation(entries, "new") > 0

	// Определяем уровень риска
	if hasDeletions || totalChanges > 15 {
		risk.Level = "high"
		risk.ReviewNeeded = true
	} else if totalChanges > 8 {
		risk.Level = "medium"
		risk.ReviewNeeded = true
	}

	// Добавляем риски
	if hasDeletions {
		risk.Risks = append(risk.Risks, "Potential breaking changes due to deletions")
		risk.Mitigations = append(risk.Mitigations, "Run comprehensive tests")
	}

	if hasNewFiles {
		risk.Risks = append(risk.Risks, "New functionality may introduce bugs")
		risk.Mitigations = append(risk.Mitigations, "Add unit tests for new code")
	}

	if totalChanges > 10 {
		risk.Risks = append(risk.Risks, "Large number of changes increases complexity")
		risk.Mitigations = append(risk.Mitigations, "Review changes in smaller batches")
	}

	// Оценка покрытия тестами
	if hasNewFiles {
		risk.TestCoverage = "insufficient"
	} else {
		risk.TestCoverage = "adequate"
	}

	return risk
}

// countOperation подсчитывает количество операций определенного типа
func (e *DiffEngineImpl) countOperation(entries []*domain.DiffEntry, operation string) int {
	count := 0
	for _, entry := range entries {
		if entry.Operation == operation {
			count++
		}
	}
	return count
}

// saveEnrichedReport сохраняет обогащенный отчет
func (e *DiffEngineImpl) saveEnrichedReport(ctx context.Context, diff *domain.DiffResult) error {
	// Создаем директорию reports/ux если не существует
	reportsDir := "reports/ux"
	if err := os.MkdirAll(reportsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create reports directory: %w", err)
	}

	// Сохраняем текстовый отчет
	txtPath := filepath.Join(reportsDir, "derived-diff.txt")
	txtContent := e.generateTextReport(diff)
	if err := os.WriteFile(txtPath, []byte(txtContent), 0o600); err != nil {
		return fmt.Errorf("failed to write text report: %w", err)
	}

	// Сохраняем JSON отчет
	jsonPath := filepath.Join(reportsDir, "derived-diff.json")
	jsonData, err := e.generateJSONReport(diff)
	if err != nil {
		return fmt.Errorf("failed to generate JSON report: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0o600); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	e.log.Info(fmt.Sprintf("Saved enriched reports to %s and %s", txtPath, jsonPath))
	return nil
}

// generateTextReport генерирует текстовый отчет
func (e *DiffEngineImpl) generateTextReport(diff *domain.DiffResult) string {
	content := fmt.Sprintf(`Derived Diff Report
Generated: %s
ID: %s

SUMMARY:
- Total Files: %d
- Added Files: %d
- Modified Files: %d
- Deleted Files: %d
- Total Lines: %d
- Added Lines: %d
- Removed Lines: %d

`, diff.GeneratedAt, diff.ID, diff.Summary.TotalFiles, diff.Summary.AddedFiles,
		diff.Summary.ModifiedFiles, diff.Summary.DeletedFiles, diff.Summary.TotalLines,
		diff.Summary.AddedLines, diff.Summary.RemovedLines)

	// Добавляем Why View
	if diff.Summary.WhyView != nil {
		content += fmt.Sprintf(`WHY VIEW:
- Reason: %s
- Confidence: %.2f
- Related Files: %d
- Context: %v

`, diff.Summary.WhyView.Reason, diff.Summary.WhyView.Confidence,
			len(diff.Summary.WhyView.RelatedFiles), diff.Summary.WhyView.Context)
	}

	// Добавляем Impact Analysis
	if diff.Summary.Impact != nil {
		content += fmt.Sprintf(`IMPACT ANALYSIS:
- Level: %s
- Breaking Changes: %t
- Affected Areas: %v
- Performance: %s
- Security: %s
- Maintainability: %s

`, diff.Summary.Impact.Level, diff.Summary.Impact.Breaking,
			diff.Summary.Impact.AffectedAreas, diff.Summary.Impact.Performance,
			diff.Summary.Impact.Security, diff.Summary.Impact.Maintainability)
	}

	// Добавляем Risk Assessment
	if diff.Summary.Risk != nil {
		content += fmt.Sprintf(`RISK ASSESSMENT:
- Level: %s
- Review Needed: %t
- Test Coverage: %s
- Risks: %v
- Mitigations: %v

`, diff.Summary.Risk.Level, diff.Summary.Risk.ReviewNeeded,
			diff.Summary.Risk.TestCoverage, diff.Summary.Risk.Risks,
			diff.Summary.Risk.Mitigations)
	}

	// Добавляем детали изменений
	content += "CHANGES:\n"
	for _, entry := range diff.Entries {
		content += fmt.Sprintf("- %s (%s)\n", entry.Path, entry.Operation)
	}

	return content
}

// generateJSONReport генерирует JSON отчет
func (e *DiffEngineImpl) generateJSONReport(diff *domain.DiffResult) ([]byte, error) {
	return json.MarshalIndent(diff, "", "  ")
}

// generateHTMLContent генерирует HTML контент
func (e *DiffEngineImpl) generateHTMLContent(entries []*domain.DiffEntry) (string, error) {
	content := `<!DOCTYPE html>
<html>
<head>
    <title>Diff Report</title>
    <style>
        body { font-family: monospace; margin: 20px; }
        .file { margin: 10px 0; padding: 10px; border: 1px solid #ccc; }
        .added { background-color: #e6ffe6; }
        .deleted { background-color: #ffe6e6; }
        .modified { background-color: #fff2e6; }
    </style>
</head>
<body>
    <h1>Diff Report</h1>
`

	for _, entry := range entries {
		class := entry.Operation
		content += fmt.Sprintf(`    <div class="file %s">
        <h3>%s (%s)</h3>
        <p>Operation: %s</p>
    </div>
`, class, entry.Path, entry.Operation, entry.Operation)
	}

	content += `</body>
</html>`

	return content, nil
}

// generateDiffID генерирует ID для diff
func (e *DiffEngineImpl) generateDiffID(beforePath, afterPath string) string {
	data := fmt.Sprintf("%s:%s:%s", beforePath, afterPath, time.Now().UTC().Format("2006-01-02T15:04:05"))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

// generateDiffIDFromResults генерирует ID для diff из результатов
func (e *DiffEngineImpl) generateDiffIDFromResults(results []*domain.ApplyResult) string {
	paths := make([]string, 0, len(results))
	for _, result := range results {
		paths = append(paths, result.Path)
	}

	data := fmt.Sprintf("%s:%s", strings.Join(paths, ","), time.Now().UTC().Format("2006-01-02T15:04:05"))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

// generateDiffIDFromEdits генерирует ID для diff из правок
func (e *DiffEngineImpl) generateDiffIDFromEdits(edits *domain.EditsJSON) string {
	paths := make([]string, 0, len(edits.Edits))
	for _, edit := range edits.Edits {
		paths = append(paths, edit.Path)
	}

	data := fmt.Sprintf("%s:%s", strings.Join(paths, ","), time.Now().UTC().Format("2006-01-02T15:04:05"))
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}
