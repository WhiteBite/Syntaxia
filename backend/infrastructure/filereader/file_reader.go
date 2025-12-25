package filereader

import (
	"os"
	"syntaxia/domain"
)

// fileReader implements the domain.FileReader interface
type fileReader struct {
}

// NewFileReader creates a new file reader that implements the domain.FileReader interface
func NewFileReader() domain.FileReader {
	return &fileReader{}
}

// ReadFile reads the content of a file
func (r *fileReader) ReadFile(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
