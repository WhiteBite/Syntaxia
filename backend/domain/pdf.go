package domain

// PDFOptions описывает опции генерации PDF.
type PDFOptions struct {
	Dark        bool
	LineNumbers bool
	PageNumbers bool
}

// PDFGenerator определяет контракт для генерации PDF.
type PDFGenerator interface {
	// Generate создаёт PDF и возвращает его содержимое в памяти.
	Generate(text string, opts PDFOptions) ([]byte, error)
	// WriteAtomic создаёт PDF и атомарно записывает его в файл (через временный файл и os.Rename).
	WriteAtomic(text string, opts PDFOptions, outputPath string) error
}
