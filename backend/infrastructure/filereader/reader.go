package filereader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"strings"
)

const maxFileSize = 5 * 1024 * 1024 // 5 MB limit per file
const expandDirectories = true      // Flag to enable directory expansion

type secureFileReader struct {
	log domain.Logger
}

func NewSecureFileReader(log domain.Logger) domain.FileContentReader {
	return &secureFileReader{log: log}
}

// pathValidationResult holds the result of path validation
type pathValidationResult struct {
	inputPath string
	absPath   string
	size      int64
}

// resolveAbsPath resolves the absolute path for an input path
func (r *secureFileReader) resolveAbsPath(inputPath, rootDir string) (string, error) {
	if filepath.IsAbs(inputPath) {
		return filepath.Abs(inputPath)
	}
	return r.sanitizeAndAbs(rootDir, inputPath)
}

// processExpandedFile processes a single expanded file from directory
func (r *secureFileReader) processExpandedFile(expandedPath, rootDir string) *pathValidationResult {
	fileInfo, err := os.Stat(expandedPath)
	if err != nil {
		r.log.Warning(fmt.Sprintf("Skipping file %s: cannot stat - %v", expandedPath, err))
		return nil
	}
	if fileInfo.Size() > maxFileSize {
		r.log.Warning(fmt.Sprintf("Skipping large file %s: size %d > %d", expandedPath, fileInfo.Size(), maxFileSize))
		return nil
	}

	relPath, err := filepath.Rel(rootDir, expandedPath)
	if err != nil {
		relPath = expandedPath
	}
	relPath = filepath.ToSlash(relPath)

	return &pathValidationResult{inputPath: relPath, absPath: expandedPath, size: fileInfo.Size()}
}

// validateAndExpandPath validates a single path and expands if directory
func (r *secureFileReader) validateAndExpandPath(inputPath, rootDir string) ([]pathValidationResult, error) {
	absPath, err := r.resolveAbsPath(inputPath, rootDir)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		if info.Size() > maxFileSize {
			return nil, fmt.Errorf("file too large: %d > %d", info.Size(), maxFileSize)
		}
		return []pathValidationResult{{inputPath: inputPath, absPath: absPath, size: info.Size()}}, nil
	}

	if !expandDirectories {
		return nil, fmt.Errorf("directory expansion disabled")
	}

	expandedPaths, err := r.expandDirectory(absPath, rootDir)
	if err != nil {
		return nil, err
	}

	var results []pathValidationResult
	for _, ep := range expandedPaths {
		if result := r.processExpandedFile(ep, rootDir); result != nil {
			results = append(results, *result)
		}
	}
	r.log.Info(fmt.Sprintf("Expanded directory %s: found %d files", inputPath, len(results)))
	return results, nil
}

// readFileContents reads file contents and updates progress
func (r *secureFileReader) readFileContents(ctx context.Context, validated []pathValidationResult, totalSize int64, progress func(int64, int64)) (map[string]string, error) {
	contents := make(map[string]string, len(validated))
	var currentSize int64

	for _, v := range validated {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		data, err := os.ReadFile(v.absPath)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Skipping file %s: read error - %v", v.inputPath, err))
			continue
		}

		contents[v.inputPath] = string(data)
		currentSize += int64(len(data))
		progress(currentSize, totalSize)
	}
	return contents, nil
}

func (r *secureFileReader) ReadContents(
	ctx context.Context,
	filePaths []string,
	rootDir string,
	progress func(current, total int64),
) (map[string]string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if progress == nil {
		progress = func(current, total int64) {}
	}

	var validated []pathValidationResult
	var totalSize int64

	for _, inputPath := range filePaths {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		results, err := r.validateAndExpandPath(inputPath, rootDir)
		if err != nil {
			r.log.Warning(fmt.Sprintf("Skipping path %s: %v", inputPath, err))
			continue
		}

		for _, res := range results {
			validated = append(validated, res)
			totalSize += res.size
		}
	}

	return r.readFileContents(ctx, validated, totalSize, progress)
}

// sanitizeAndAbs converts path to absolute, allowing files from any location
func (r *secureFileReader) sanitizeAndAbs(rootDir, relPath string) (string, error) {
	// If path is already absolute, use it directly
	if filepath.IsAbs(relPath) {
		absPath, err := filepath.Abs(relPath)
		if err != nil {
			return "", fmt.Errorf("could not get absolute path: %w", err)
		}
		return absPath, nil
	}

	// Otherwise, join with rootDir
	cleanRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path for root: %w", err)
	}

	absPath := filepath.Join(cleanRootDir, relPath)
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		return "", fmt.Errorf("could not get absolute path: %w", err)
	}

	return absPath, nil
}

// expandDirectory recursively walks through a directory and returns all file paths
func (r *secureFileReader) expandDirectory(dirPath, rootDir string) ([]string, error) {
	var files []string

	// Limit recursion depth to prevent DoS
	const maxDepth = 20
	const maxFiles = 10000 // Limit total files to prevent DoS

	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			r.log.Warning(fmt.Sprintf("Error accessing path %s: %v", path, err))
			return nil // Continue walking despite errors
		}

		// Check if we've reached the file limit
		if len(files) >= maxFiles {
			r.log.Warning(fmt.Sprintf("Stopping directory expansion: max files limit (%d) reached", maxFiles))
			return filepath.SkipAll // Stop walking completely
		}

		// Skip directories, only add files
		if d.IsDir() {
			// Check recursion depth
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return nil // Continue walking
			}

			// Count separators to determine depth
			depth := strings.Count(relPath, string(filepath.Separator)) + 1
			if depth > maxDepth {
				r.log.Warning(fmt.Sprintf("Skipping directory %s: max depth exceeded", path))
				return filepath.SkipDir
			}

			return nil
		}

		// Add the file path to the list
		files = append(files, path)

		return nil
	})

	return files, err
}
