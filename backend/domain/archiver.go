package domain

// Archiver определяет контракт для упаковки набора файлов в ZIP.
type Archiver interface {
	// ZipFilesAtomic принимает набор (имя -> содержимое) и атомарно записывает ZIP на диск.
	ZipFilesAtomic(files map[string][]byte, outputPath string) error
}
