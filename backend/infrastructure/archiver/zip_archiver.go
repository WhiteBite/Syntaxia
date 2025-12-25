package archiver

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"syntaxia/domain"
)

type ZipArchiver struct {
	log domain.Logger
}

func NewZipArchiver(log domain.Logger) domain.Archiver {
	return &ZipArchiver{log: log}
}

// ZipFilesAtomic пишет ZIP с файлами (имя -> содержимое) атомарно.
func (a *ZipArchiver) ZipFilesAtomic(files map[string][]byte, outputPath string) error {
	dir := filepath.Dir(outputPath)
	tmp, err := os.CreateTemp(dir, "zip-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	zw := zip.NewWriter(tmp)

	// детерминированный порядок
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		b := files[name]
		f, err := zw.Create(name)
		if err != nil {
			_ = zw.Close()
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
			return fmt.Errorf("zip create %s: %w", name, err)
		}
		if _, err := f.Write(b); err != nil {
			_ = zw.Close()
			_ = tmp.Close()
			_ = os.Remove(tmpPath)
			return fmt.Errorf("zip write %s: %w", name, err)
		}
	}
	if err := zw.Close(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpPath)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := os.Rename(tmpPath, outputPath); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return nil
}
