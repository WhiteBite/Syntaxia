package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BatchWriteParams represents parameters for batch file write
type BatchWriteParams struct {
	Files []FileWrite `json:"files"`
}

// FileWrite represents a single file write operation
type FileWrite struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// BatchWriteResult represents the result of batch file write
type BatchWriteResult struct {
	Success []string          `json:"success"`
	Failed  []BatchWriteError `json:"failed"`
}

// BatchWriteError represents a failed file write
type BatchWriteError struct {
	Path  string `json:"path"`
	Error string `json:"error"`
}

// DiffPreviewParams represents parameters for diff preview
type DiffPreviewParams struct {
	Path       string `json:"path"`
	NewContent string `json:"new_content"`
}

// DiffPreviewResult represents the result of diff preview
type DiffPreviewResult struct {
	Path         string   `json:"path"`
	Diff         string   `json:"diff"`
	LinesAdded   int      `json:"lines_added"`
	LinesRemoved int      `json:"lines_removed"`
	Preview      []string `json:"preview"`
}

func (h *FileToolsHandler) batchWriteFiles(args map[string]any, projectRoot string) (string, error) {
	filesArg, ok := args["files"]
	if !ok {
		return "", fmt.Errorf("files parameter is required")
	}

	filesSlice, ok := filesArg.([]interface{})
	if !ok {
		return "", fmt.Errorf("files must be an array")
	}

	if len(filesSlice) == 0 {
		return "", fmt.Errorf("files array cannot be empty")
	}

	// Check if sandbox is available
	if h.SandboxFS == nil {
		return "", fmt.Errorf("sandbox not initialized - cannot write files")
	}

	// Parse and validate all files first
	var files []FileWrite
	for i, f := range filesSlice {
		fileMap, ok := f.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("file at index %d is not a valid object", i)
		}

		path, _ := fileMap["path"].(string)
		content, _ := fileMap["content"].(string)

		if path == "" {
			return "", fmt.Errorf("file at index %d missing 'path'", i)
		}
		if content == "" {
			return "", fmt.Errorf("file at index %d missing 'content'", i)
		}

		// Check content size
		if len(content) > MaxFileSize {
			return "", fmt.Errorf("file '%s' content exceeds maximum size (%d bytes)", path, MaxFileSize)
		}

		files = append(files, FileWrite{Path: path, Content: content})
	}

	// Validate all paths before writing (atomic check)
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root: %w", err)
	}
	absProjectRoot = filepath.Clean(absProjectRoot)

	var validatedPaths []string
	for _, file := range files {
		fullPath := filepath.Join(projectRoot, file.Path)
		absFullPath, err := filepath.Abs(fullPath)
		if err != nil {
			return "", fmt.Errorf("failed to resolve path '%s': %w", file.Path, err)
		}
		absFullPath = filepath.Clean(absFullPath)

		if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
			return "", fmt.Errorf("path traversal not allowed for '%s'", file.Path)
		}
		validatedPaths = append(validatedPaths, absFullPath)
	}

	// Write all files atomically to sandbox
	result := BatchWriteResult{
		Success: make([]string, 0, len(files)),
		Failed:  make([]BatchWriteError, 0),
	}

	for i, file := range files {
		err := h.SandboxFS.WriteFile(validatedPaths[i], []byte(file.Content), 0644)
		if err != nil {
			result.Failed = append(result.Failed, BatchWriteError{
				Path:  file.Path,
				Error: err.Error(),
			})
		} else {
			result.Success = append(result.Success, file.Path)
		}
	}

	// Format result
	var sb strings.Builder
	if len(result.Failed) > 0 {
		sb.WriteString(fmt.Sprintf("Batch write partially failed. %d succeeded, %d failed.\n\n",
			len(result.Success), len(result.Failed)))
		sb.WriteString("Failed files:\n")
		for _, f := range result.Failed {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", f.Path, f.Error))
		}
		if len(result.Success) > 0 {
			sb.WriteString("\nSuccessful files:\n")
			for _, p := range result.Success {
				sb.WriteString(fmt.Sprintf("- %s\n", p))
			}
		}
	} else {
		sb.WriteString(fmt.Sprintf("Successfully wrote %d files to sandbox:\n", len(result.Success)))
		for _, p := range result.Success {
			sb.WriteString(fmt.Sprintf("- %s\n", p))
		}
		sb.WriteString("\nUse 'Review Changes' to preview and apply.")
	}

	return sb.String(), nil
}

func (h *FileToolsHandler) previewDiff(args map[string]any, projectRoot string) (string, error) {
	path, _ := args["path"].(string)
	newContent, _ := args["new_content"].(string)

	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if newContent == "" {
		return "", fmt.Errorf("new_content is required")
	}

	fullPath := filepath.Join(projectRoot, path)

	// Security check
	absProjectRoot, err := filepath.Abs(projectRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve project root: %w", err)
	}
	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve file path: %w", err)
	}
	absProjectRoot = filepath.Clean(absProjectRoot)
	absFullPath = filepath.Clean(absFullPath)

	if !strings.HasPrefix(absFullPath, absProjectRoot+string(filepath.Separator)) && absFullPath != absProjectRoot {
		return "", fmt.Errorf("path traversal not allowed")
	}

	// Read current content (may not exist for new files)
	var oldContent string
	oldBytes, err := os.ReadFile(fullPath)
	if err == nil {
		oldContent = string(oldBytes)
	}

	// Generate unified diff
	diff, linesAdded, linesRemoved := generateUnifiedDiff(path, oldContent, newContent)

	// Create preview (first 20 lines of diff)
	diffLines := strings.Split(diff, "\n")
	previewLines := diffLines
	const maxPreviewLines = 20
	if len(diffLines) > maxPreviewLines {
		previewLines = diffLines[:maxPreviewLines]
	}

	result := DiffPreviewResult{
		Path:         path,
		Diff:         diff,
		LinesAdded:   linesAdded,
		LinesRemoved: linesRemoved,
		Preview:      previewLines,
	}

	// Format output
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== Diff Preview: %s ===\n", path))
	sb.WriteString(fmt.Sprintf("Lines added: %d, Lines removed: %d\n\n", result.LinesAdded, result.LinesRemoved))

	if diff == "" {
		sb.WriteString("No changes detected.\n")
	} else {
		sb.WriteString(diff)
		if len(diffLines) > maxPreviewLines {
			sb.WriteString(fmt.Sprintf("\n... and %d more lines\n", len(diffLines)-maxPreviewLines))
		}
	}

	return sb.String(), nil
}

// generateUnifiedDiff generates a unified diff between old and new content
func generateUnifiedDiff(filename, oldContent, newContent string) (string, int, int) {
	oldLines := strings.Split(oldContent, "\n")
	newLines := strings.Split(newContent, "\n")

	// Simple line-by-line diff
	var diffLines []string
	diffLines = append(diffLines, fmt.Sprintf("--- a/%s", filename))
	diffLines = append(diffLines, fmt.Sprintf("+++ b/%s", filename))

	linesAdded := 0
	linesRemoved := 0

	// Use simple LCS-based diff
	hunks := computeDiffHunks(oldLines, newLines)

	for _, hunk := range hunks {
		diffLines = append(diffLines, hunk.header)
		for _, line := range hunk.lines {
			diffLines = append(diffLines, line)
			if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
				linesAdded++
			} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
				linesRemoved++
			}
		}
	}

	if len(hunks) == 0 {
		return "", 0, 0
	}

	return strings.Join(diffLines, "\n"), linesAdded, linesRemoved
}

type diffHunk struct {
	header string
	lines  []string
}

// computeDiffHunks computes diff hunks between old and new lines
func computeDiffHunks(oldLines, newLines []string) []diffHunk {
	var hunks []diffHunk

	const contextLines = 3
	var currentHunk *diffHunk
	hunkStart := -1

	i, j := 0, 0
	for i < len(oldLines) || j < len(newLines) {
		if i < len(oldLines) && j < len(newLines) && oldLines[i] == newLines[j] {
			// Lines match - context or end of hunk
			if currentHunk != nil {
				currentHunk.lines = append(currentHunk.lines, " "+oldLines[i])
				// Check if we should close the hunk
				if len(currentHunk.lines) > 0 {
					lastNonContext := -1
					for k := len(currentHunk.lines) - 1; k >= 0; k-- {
						if !strings.HasPrefix(currentHunk.lines[k], " ") {
							lastNonContext = k
							break
						}
					}
					if lastNonContext >= 0 && len(currentHunk.lines)-lastNonContext > contextLines {
						// Trim trailing context and close hunk
						currentHunk.lines = currentHunk.lines[:lastNonContext+contextLines+1]
						hunks = append(hunks, *currentHunk)
						currentHunk = nil
					}
				}
			}
			i++
			j++
		} else {
			// Lines differ
			if currentHunk == nil {
				// Start new hunk with context
				hunkStart = i
				startContext := i - contextLines
				if startContext < 0 {
					startContext = 0
				}
				currentHunk = &diffHunk{
					lines: make([]string, 0),
				}
				for k := startContext; k < i; k++ {
					currentHunk.lines = append(currentHunk.lines, " "+oldLines[k])
				}
			}

			// Find how many lines differ
			if i < len(oldLines) && (j >= len(newLines) || !lineExistsAhead(oldLines[i], newLines, j)) {
				currentHunk.lines = append(currentHunk.lines, "-"+oldLines[i])
				i++
			} else if j < len(newLines) {
				currentHunk.lines = append(currentHunk.lines, "+"+newLines[j])
				j++
			}
		}
	}

	// Close any remaining hunk
	if currentHunk != nil && len(currentHunk.lines) > 0 {
		// Update header
		currentHunk.header = fmt.Sprintf("@@ -%d,%d +%d,%d @@",
			hunkStart+1, len(oldLines)-hunkStart, hunkStart+1, len(newLines)-hunkStart)
		hunks = append(hunks, *currentHunk)
	}

	// Update headers for all hunks
	for i := range hunks {
		if hunks[i].header == "" {
			hunks[i].header = "@@ -1 +1 @@"
		}
	}

	return hunks
}

// lineExistsAhead checks if a line exists in the slice starting from index
func lineExistsAhead(line string, lines []string, startIdx int) bool {
	const lookAhead = 5
	end := startIdx + lookAhead
	if end > len(lines) {
		end = len(lines)
	}
	for i := startIdx; i < end; i++ {
		if lines[i] == line {
			return true
		}
	}
	return false
}
