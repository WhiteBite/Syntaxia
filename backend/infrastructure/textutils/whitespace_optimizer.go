package textutils

import (
	"strings"
	"unicode"
)

// WhitespaceOptimizer оптимизирует пробельные символы в тексте
// Использует однопроходный алгоритм O(n) без регулярных выражений для максимальной производительности
type WhitespaceOptimizer struct{}

// NewWhitespaceOptimizer создает новый оптимизатор пробелов
func NewWhitespaceOptimizer() *WhitespaceOptimizer {
	return &WhitespaceOptimizer{}
}

// CollapseEmptyLines схлопывает множественные пустые строки (2+) в одну
// Алгоритм: однопроходный O(n), zero-allocation для небольших файлов
func (w *WhitespaceOptimizer) CollapseEmptyLines(content string) string {
	if len(content) == 0 {
		return content
	}

	// Предварительная проверка - если нет 3+ переносов подряд, возвращаем как есть
	if !strings.Contains(content, "\n\n\n") {
		return content
	}

	// Используем strings.Builder для эффективной конкатенации
	// Предаллоцируем примерный размер (обычно результат меньше исходника)
	var b strings.Builder
	b.Grow(len(content))

	emptyLineCount := 0
	lineStart := 0

	for i := 0; i < len(content); i++ {
		if content[i] == '\n' {
			line := content[lineStart:i]
			isEmptyLine := isWhitespaceOnly(line)

			if isEmptyLine {
				emptyLineCount++
				if emptyLineCount <= 2 {
					b.WriteString(line)
					b.WriteByte('\n')
				}
				// Пропускаем если уже 2+ пустых строк
			} else {
				emptyLineCount = 0
				b.WriteString(line)
				b.WriteByte('\n')
			}
			lineStart = i + 1
		}
	}

	// Обработка последней строки (без \n в конце)
	if lineStart < len(content) {
		line := content[lineStart:]
		if !isWhitespaceOnly(line) || emptyLineCount < 2 {
			b.WriteString(line)
		}
	}

	return b.String()
}

// isWhitespaceOnly проверяет, содержит ли строка только пробельные символы
// Инлайн-оптимизация для горячего пути
func isWhitespaceOnly(s string) bool {
	for i := 0; i < len(s); i++ {
		if !unicode.IsSpace(rune(s[i])) {
			return false
		}
	}
	return true
}

// TrimTrailingWhitespace удаляет пробелы в конце каждой строки
// Полезно для дополнительной экономии токенов
func (w *WhitespaceOptimizer) TrimTrailingWhitespace(content string) string {
	if len(content) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var b strings.Builder
	b.Grow(len(content))

	for i, line := range lines {
		b.WriteString(strings.TrimRightFunc(line, unicode.IsSpace))
		if i < len(lines)-1 {
			b.WriteByte('\n')
		}
	}

	return b.String()
}

// OptimizeWhitespace применяет все оптимизации пробелов
func (w *WhitespaceOptimizer) OptimizeWhitespace(content string) string {
	result := w.TrimTrailingWhitespace(content)
	result = w.CollapseEmptyLines(result)
	return result
}
