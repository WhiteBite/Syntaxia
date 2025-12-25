package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// validateLimits validates memory and token limits
func (s *Service) validateLimits(options *BuildOptions) error {
	// Strict memory limit validation
	if options.MaxMemoryMB > 0 && options.MaxMemoryMB > 500 {
		return fmt.Errorf("memory limit cannot exceed 500MB for safety")
	}

	// Token limit validation - very permissive, frontend controls the actual limit
	const maxTokenLimit = 10000000 // 10M tokens - matches frontend max
	if options.MaxTokens > 0 && options.MaxTokens > maxTokenLimit {
		return fmt.Errorf("token limit cannot exceed %d (requested: %d), please adjust settings on frontend", maxTokenLimit, options.MaxTokens)
	}

	return nil
}

// estimateTotalSize estimates total size of files and identifies oversized files
func (s *Service) estimateTotalSize(projectPath string, includedPaths []string) (int64, []string, error) {
	var totalSize int64
	var oversizedFiles []string

	for _, filePath := range includedPaths {
		fullPath := filepath.Join(projectPath, filePath)
		if info, err := os.Stat(fullPath); err == nil {
			totalSize += info.Size()
			// Flag files larger than 1MB
			if info.Size() > 1024*1024 {
				oversizedFiles = append(oversizedFiles, filePath)
			}
		}
	}

	return totalSize, oversizedFiles, nil
}

// generateContextName generates a human-readable name for the context
func (s *Service) generateContextName(projectPath string, files []string) string {
	projectName := filepath.Base(projectPath)

	if len(files) == 1 {
		fileName := filepath.Base(files[0])
		return fmt.Sprintf("%s - %s", projectName, fileName)
	}

	return fmt.Sprintf("%s - %d files", projectName, len(files))
}

// checkMemoryLimit validates memory constraints
func (s *Service) checkMemoryLimit(options *BuildOptions, totalSize int64, oversizedFiles []string) error {
	if options.MaxMemoryMB > 0 {
		maxBytes := int64(options.MaxMemoryMB) * 1024 * 1024
		if totalSize > maxBytes {
			return fmt.Errorf("context would exceed memory limit: %d MB > %d MB. Oversized files: %v",
				totalSize/(1024*1024), options.MaxMemoryMB, oversizedFiles)
		}
	}
	return nil
}

// readAndUnmarshalJSON is a helper for reading and unmarshaling JSON files
func (s *Service) readAndUnmarshalJSON(filePath, entityName string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s not found", entityName)
		}
		return fmt.Errorf("failed to read %s: %w", entityName, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal %s: %w", entityName, err)
	}
	return nil
}

// emitEvent safely emits an event if eventBus is available
func (s *Service) emitEvent(event string, data interface{}) {
	if s.eventBus != nil {
		s.eventBus.Emit(event, data)
	}
}

// saveContext saves a context to disk
func (s *Service) saveContext(ctx *Context) error {
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	contextPath := filepath.Join(s.contextDir, ctx.ID+".json")
	if err := os.WriteFile(contextPath, data, 0o600); err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}
