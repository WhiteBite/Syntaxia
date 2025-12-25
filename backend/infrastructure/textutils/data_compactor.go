package textutils

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
)

// DataCompactor сжимает JSON и YAML файлы для экономии токенов
// Удаляет форматирование, сохраняя структуру данных
type DataCompactor struct{}

// NewDataCompactor создает новый компактор данных
func NewDataCompactor() *DataCompactor {
	return &DataCompactor{}
}

// Compact сжимает содержимое файла на основе расширения
func (c *DataCompactor) Compact(content, filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".json":
		return c.CompactJSON(content)
	case ".yaml", ".yml":
		return c.CompactYAML(content)
	default:
		return content
	}
}

// CompactJSON минимизирует JSON, удаляя форматирование
// Использует стандартный json пакет для надежности
func (c *DataCompactor) CompactJSON(content string) string {
	if len(content) == 0 {
		return content
	}

	// Быстрая проверка - если уже компактный (нет переносов), возвращаем как есть
	if !strings.Contains(content, "\n") {
		return content
	}

	// Парсим и сериализуем без форматирования
	var data interface{}
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		// Если не валидный JSON, возвращаем как есть
		return content
	}

	// Используем buffer для избежания аллокаций
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false) // Не экранируем HTML для читаемости

	if err := encoder.Encode(data); err != nil {
		return content
	}

	// Encode добавляет \n в конце, убираем его
	result := buf.String()
	return strings.TrimSuffix(result, "\n")
}

// CompactYAML минимизирует YAML
// Упрощенная реализация без внешних зависимостей:
// - Удаляет комментарии
// - Схлопывает пустые строки
// - Минимизирует отступы (опционально)
func (c *DataCompactor) CompactYAML(content string) string {
	if len(content) == 0 {
		return content
	}

	lines := strings.Split(content, "\n")
	var result []string
	prevEmpty := false

	for _, line := range lines {
		// Удаляем комментарии (но не в строках)
		trimmed := strings.TrimSpace(line)

		// Пропускаем пустые строки (оставляем максимум одну)
		if trimmed == "" {
			if !prevEmpty {
				result = append(result, "")
				prevEmpty = true
			}
			continue
		}
		prevEmpty = false

		// Пропускаем строки-комментарии
		if strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Удаляем inline комментарии (осторожно - не в строках)
		line = c.removeYAMLInlineComment(line)

		// Убираем trailing whitespace
		line = strings.TrimRight(line, " \t")

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// removeYAMLInlineComment удаляет inline комментарии из YAML строки
// Учитывает, что # может быть внутри строки
func (c *DataCompactor) removeYAMLInlineComment(line string) string {
	// Если нет #, возвращаем как есть
	if !strings.Contains(line, "#") {
		return line
	}

	// Простая эвристика: если # после пробела и не в кавычках
	inSingleQuote := false
	inDoubleQuote := false

	for i := 0; i < len(line); i++ {
		ch := line[i]

		switch ch {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote {
				// Проверяем, что перед # есть пробел (YAML требование)
				if i > 0 && (line[i-1] == ' ' || line[i-1] == '\t') {
					return strings.TrimRight(line[:i], " \t")
				}
			}
		}
	}

	return line
}

// IsDataFile проверяет, является ли файл файлом данных (JSON/YAML)
func IsDataFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json", ".yaml", ".yml":
		return true
	default:
		return false
	}
}

// CompactDataFile глобальная функция для удобства
func CompactDataFile(content, filePath string) string {
	return defaultDataCompactor.Compact(content, filePath)
}

var defaultDataCompactor = NewDataCompactor()
