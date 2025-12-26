package textutils

import (
	"path/filepath"
	"strings"
)

// TestFilter определяет, является ли файл тестовым
// Оптимизирован для быстрой проверки без регулярных выражений
type TestFilter struct {
	// Кэш расширений для быстрого lookup
	testExtensions map[string]bool
	// Суффиксы имен файлов (без расширения)
	testSuffixes []string
	// Префиксы имен файлов
	testPrefixes []string
	// Директории с тестами
	testDirs map[string]bool
}

// NewTestFilter создает новый фильтр тестовых файлов
func NewTestFilter() *TestFilter {
	return &TestFilter{
		testExtensions: map[string]bool{
			".test.js":     true,
			".test.ts":     true,
			".test.jsx":    true,
			".test.tsx":    true,
			".spec.js":     true,
			".spec.ts":     true,
			".spec.jsx":    true,
			".spec.tsx":    true,
			".test.go":     true, // Не стандартно, но встречается
			".stories.js":  true, // Storybook
			".stories.ts":  true,
			".stories.tsx": true,
		},
		testSuffixes: []string{
			"_test",  // Go: file_test.go
			".test",  // JS/TS: file.test.js
			".spec",  // JS/TS: file.spec.ts
			"_spec",  // Ruby: file_spec.rb
			"Test",   // Java: FileTest.java
			"Tests",  // C#: FileTests.cs
			"_tests", // Python: file_tests.py
		},
		testPrefixes: []string{
			"test_", // Python: test_file.py
			"Test",  // Java: TestFile.java
		},
		testDirs: map[string]bool{
			"test":         true,
			"tests":        true,
			"__tests__":    true,
			"__test__":     true,
			"spec":         true,
			"specs":        true,
			"testing":      true,
			"testdata":     true,
			"test-data":    true,
			"fixtures":     true,
			"__fixtures__": true,
			"__mocks__":    true,
			"mocks":        true,
			"e2e":          true,
			"integration":  true,
			"unit":         true,
		},
	}
}

// IsTestFile проверяет, является ли файл тестовым
// Использует быстрые строковые операции вместо regex
func (f *TestFilter) IsTestFile(filePath string) bool {
	// Нормализуем путь
	filePath = filepath.ToSlash(filePath)

	// Проверяем директории в пути
	if f.isInTestDirectory(filePath) {
		return true
	}

	// Получаем имя файла и расширение
	fileName := filepath.Base(filePath)

	// Проверяем составные расширения (.test.js, .spec.ts)
	if f.hasTestExtension(fileName) {
		return true
	}

	// Проверяем суффиксы и префиксы
	nameWithoutExt := f.removeExtension(fileName)

	return f.hasTestSuffix(nameWithoutExt) || f.hasTestPrefix(nameWithoutExt)
}

// isInTestDirectory проверяет, находится ли файл в тестовой директории
func (f *TestFilter) isInTestDirectory(filePath string) bool {
	parts := strings.Split(filePath, "/")
	for _, part := range parts {
		if f.testDirs[strings.ToLower(part)] {
			return true
		}
	}
	return false
}

// hasTestExtension проверяет составные тестовые расширения
func (f *TestFilter) hasTestExtension(fileName string) bool {
	lower := strings.ToLower(fileName)
	for ext := range f.testExtensions {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// hasTestSuffix проверяет суффиксы имени файла
func (f *TestFilter) hasTestSuffix(nameWithoutExt string) bool {
	for _, suffix := range f.testSuffixes {
		if strings.HasSuffix(nameWithoutExt, suffix) {
			return true
		}
		// Case-insensitive для Java/C# стиля
		if strings.HasSuffix(strings.ToLower(nameWithoutExt), strings.ToLower(suffix)) {
			return true
		}
	}
	return false
}

// hasTestPrefix проверяет префиксы имени файла
func (f *TestFilter) hasTestPrefix(nameWithoutExt string) bool {
	for _, prefix := range f.testPrefixes {
		if strings.HasPrefix(nameWithoutExt, prefix) {
			return true
		}
	}
	return false
}

// removeExtension удаляет расширение из имени файла
func (f *TestFilter) removeExtension(fileName string) string {
	ext := filepath.Ext(fileName)
	if ext == "" {
		return fileName
	}
	return fileName[:len(fileName)-len(ext)]
}

// FilterTestFiles фильтрует список путей, исключая тестовые файлы
// Возвращает новый slice без тестовых файлов
func (f *TestFilter) FilterTestFiles(paths []string) []string {
	if len(paths) == 0 {
		return paths
	}

	// Предаллоцируем с оптимистичным размером
	result := make([]string, 0, len(paths))

	for _, path := range paths {
		if !f.IsTestFile(path) {
			result = append(result, path)
		}
	}

	return result
}

// CountTestFiles подсчитывает количество тестовых файлов в списке
func (f *TestFilter) CountTestFiles(paths []string) int {
	count := 0
	for _, path := range paths {
		if f.IsTestFile(path) {
			count++
		}
	}
	return count
}

// IsTestFile глобальная функция для удобства использования
func IsTestFile(filePath string) bool {
	return defaultTestFilter.IsTestFile(filePath)
}

// FilterTestFiles глобальная функция для удобства использования
func FilterTestFiles(paths []string) []string {
	return defaultTestFilter.FilterTestFiles(paths)
}

// defaultTestFilter синглтон для глобальных функций
var defaultTestFilter = NewTestFilter()
