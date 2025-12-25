package filesystem

import (
	"os"
	"path/filepath"

	"syntaxia/domain"
)

// FilePathProvider implements domain.PathProvider using standard filepath functions
type FilePathProvider struct{}

// NewFilePathProvider creates a new FilePathProvider
func NewFilePathProvider() domain.PathProvider {
	return &FilePathProvider{}
}

// Join соединяет элементы пути
func (p *FilePathProvider) Join(elem ...string) string {
	return filepath.Join(elem...)
}

// Base возвращает последний элемент пути
func (p *FilePathProvider) Base(path string) string {
	return filepath.Base(path)
}

// Dir возвращает все элементы пути кроме последнего
func (p *FilePathProvider) Dir(path string) string {
	return filepath.Dir(path)
}

// IsAbs проверяет, является ли путь абсолютным
func (p *FilePathProvider) IsAbs(path string) bool {
	return filepath.IsAbs(path)
}

// Clean возвращает очищенную версию пути
func (p *FilePathProvider) Clean(path string) string {
	return filepath.Clean(path)
}

// Getwd возвращает текущую рабочую директорию
func (p *FilePathProvider) Getwd() (string, error) {
	return os.Getwd()
}
