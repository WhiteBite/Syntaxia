package textutils

import (
	"strings"
)

// LicenseStripper удаляет блоки лицензий и copyright из начала файлов
// Оптимизирован для минимального влияния на производительность
type LicenseStripper struct {
	// Предкомпилированные паттерны для быстрого поиска (lowercase)
	licenseKeywords []string
}

// NewLicenseStripper создает новый stripper с предзагруженными паттернами
func NewLicenseStripper() *LicenseStripper {
	return &LicenseStripper{
		licenseKeywords: []string{
			"copyright",
			"license",
			"licensed",
			"spdx-license",
			"mit license",
			"apache license",
			"bsd license",
			"gnu general public",
			"gpl",
			"lgpl",
			"mozilla public",
			"all rights reserved",
			"permission is hereby granted",
			"redistribution and use",
			"this file is part of",
			"автор:",
			"лицензия",
		},
	}
}

// Strip удаляет лицензионный заголовок из начала файла
// Возвращает контент без изменений если лицензия не найдена
func (l *LicenseStripper) Strip(content string) string {
	if len(content) == 0 {
		return content
	}

	// Быстрая проверка - ищем начало комментария в первых 50 символах
	firstChars := content
	if len(firstChars) > 50 {
		firstChars = content[:50]
	}

	// Определяем тип комментария
	commentStart, commentEnd, lineComment := l.detectCommentStyle(firstChars)

	if commentStart == "" && lineComment == "" {
		return content
	}

	// Ищем блок комментария в начале файла
	trimmed := strings.TrimLeft(content, " \t\n\r")

	if commentStart != "" && strings.HasPrefix(trimmed, commentStart) {
		return l.stripBlockComment(content, trimmed, commentStart, commentEnd)
	}

	if lineComment != "" && strings.HasPrefix(trimmed, lineComment) {
		return l.stripLineComments(content, trimmed, lineComment)
	}

	return content
}

// detectCommentStyle определяет стиль комментариев по первым символам
func (l *LicenseStripper) detectCommentStyle(sample string) (blockStart, blockEnd, lineComment string) {
	sample = strings.TrimLeft(sample, " \t\n\r")

	// Проверяем блочные комментарии
	if strings.HasPrefix(sample, "/*") {
		return "/*", "*/", "//"
	}
	if strings.HasPrefix(sample, "<!--") {
		return "<!--", "-->", ""
	}
	if strings.HasPrefix(sample, "\"\"\"") {
		return "\"\"\"", "\"\"\"", "#"
	}
	if strings.HasPrefix(sample, "'''") {
		return "'''", "'''", "#"
	}

	// Проверяем строчные комментарии
	if strings.HasPrefix(sample, "//") {
		return "", "", "//"
	}
	if strings.HasPrefix(sample, "#") && !strings.HasPrefix(sample, "#!") {
		return "", "", "#"
	}
	if strings.HasPrefix(sample, "--") {
		return "", "", "--"
	}

	return "", "", ""
}

// stripBlockComment удаляет блочный комментарий если он содержит лицензию
func (l *LicenseStripper) stripBlockComment(original, trimmed, start, end string) string {
	endIdx := strings.Index(trimmed, end)
	if endIdx == -1 {
		return original
	}

	commentBlock := trimmed[:endIdx+len(end)]

	// Проверяем, содержит ли комментарий лицензионные ключевые слова
	if !l.containsLicenseKeyword(commentBlock) {
		return original
	}

	// Удаляем комментарий и ведущие пробелы после него
	afterComment := trimmed[endIdx+len(end):]
	afterComment = strings.TrimLeft(afterComment, " \t\n\r")

	return afterComment
}

// stripLineComments удаляет последовательные строчные комментарии с лицензией
func (l *LicenseStripper) stripLineComments(original, trimmed, prefix string) string {
	lines := strings.Split(trimmed, "\n")
	commentLines := []string{}
	lastCommentLine := 0

	// Собираем последовательные строки комментариев в начале
	for i, line := range lines {
		trimmedLine := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmedLine, prefix) || trimmedLine == "" {
			commentLines = append(commentLines, line)
			if strings.HasPrefix(trimmedLine, prefix) {
				lastCommentLine = i
			}
		} else {
			break
		}
	}

	if len(commentLines) == 0 {
		return original
	}

	// Проверяем, содержит ли блок комментариев лицензию
	commentBlock := strings.Join(commentLines[:lastCommentLine+1], "\n")
	if !l.containsLicenseKeyword(commentBlock) {
		return original
	}

	// Возвращаем контент после блока комментариев
	remaining := lines[lastCommentLine+1:]
	result := strings.Join(remaining, "\n")
	return strings.TrimLeft(result, "\n")
}

// containsLicenseKeyword проверяет наличие лицензионных ключевых слов
// Использует lowercase сравнение для надежности
func (l *LicenseStripper) containsLicenseKeyword(text string) bool {
	lower := strings.ToLower(text)
	for _, keyword := range l.licenseKeywords {
		if strings.Contains(lower, keyword) {
			return true
		}
	}
	return false
}

// StripWithLanguageHint удаляет лицензию с подсказкой о языке для оптимизации
func (l *LicenseStripper) StripWithLanguageHint(content, ext string) string {
	// Для некоторых расширений можно сразу определить стиль комментариев
	// и пропустить автодетект
	switch ext {
	case ".go", ".js", ".ts", ".java", ".c", ".cpp", ".cs", ".swift", ".kt", ".rs":
		return l.stripWithStyle(content, "/*", "*/", "//")
	case ".py", ".rb", ".sh", ".yaml", ".yml":
		return l.stripWithStyle(content, "", "", "#")
	case ".html", ".xml", ".vue", ".svelte":
		return l.stripWithStyle(content, "<!--", "-->", "")
	case ".sql":
		return l.stripWithStyle(content, "/*", "*/", "--")
	default:
		return l.Strip(content)
	}
}

// stripWithStyle удаляет лицензию с известным стилем комментариев
func (l *LicenseStripper) stripWithStyle(content, blockStart, blockEnd, lineComment string) string {
	if len(content) == 0 {
		return content
	}

	trimmed := strings.TrimLeft(content, " \t\n\r")

	if blockStart != "" && strings.HasPrefix(trimmed, blockStart) {
		return l.stripBlockComment(content, trimmed, blockStart, blockEnd)
	}

	if lineComment != "" && strings.HasPrefix(trimmed, lineComment) {
		return l.stripLineComments(content, trimmed, lineComment)
	}

	return content
}
