package filesystem

import (
	"os"

	"syntaxia/domain"
)

// OSFileStatProvider implements domain.FileStatProvider using standard os functions
type OSFileStatProvider struct{}

// NewOSFileStatProvider creates a new OSFileStatProvider
func NewOSFileStatProvider() domain.FileStatProvider {
	return &OSFileStatProvider{}
}

// Stat возвращает информацию о файле
func (p *OSFileStatProvider) Stat(name string) (domain.FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return fi, nil
}
