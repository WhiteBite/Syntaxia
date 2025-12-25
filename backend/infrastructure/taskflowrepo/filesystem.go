package taskflowrepo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"syntaxia/domain"
)

// FileSystemTaskflowRepository implements TaskflowRepository using file system
type FileSystemTaskflowRepository struct {
	statusPath string
}

// NewFileSystemTaskflowRepository creates a new file system taskflow repository
func NewFileSystemTaskflowRepository(statusPath string) *FileSystemTaskflowRepository {
	return &FileSystemTaskflowRepository{
		statusPath: statusPath,
	}
}

// LoadStatuses loads task statuses from file
func (r *FileSystemTaskflowRepository) LoadStatuses() (map[string]domain.TaskState, error) {
	data, err := os.ReadFile(r.statusPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]domain.TaskState), nil
		}
		return nil, err
	}

	var statusData struct {
		Tasks []struct {
			ID    string `json:"id"`
			State string `json:"state"`
		} `json:"tasks"`
	}

	if err := json.Unmarshal(data, &statusData); err != nil {
		return nil, err
	}

	statuses := make(map[string]domain.TaskState)
	for _, task := range statusData.Tasks {
		statuses[task.ID] = domain.TaskState(task.State)
	}

	return statuses, nil
}

// SaveStatuses saves task statuses to file
func (r *FileSystemTaskflowRepository) SaveStatuses(statuses map[string]domain.TaskState) error {
	tasks := make([]struct {
		ID    string `json:"id"`
		State string `json:"state"`
	}, 0, len(statuses))

	for taskID, state := range statuses {
		tasks = append(tasks, struct {
			ID    string `json:"id"`
			State string `json:"state"`
		}{
			ID:    taskID,
			State: string(state),
		})
	}

	statusData := struct {
		Version int         `json:"version"`
		Tasks   interface{} `json:"tasks"`
		History interface{} `json:"history"`
	}{
		Version: 1,
		Tasks:   tasks,
		History: []interface{}{},
	}

	data, err := json.MarshalIndent(statusData, "", "  ")
	if err != nil {
		return err
	}

	// Create directory if needed
	dir := filepath.Dir(r.statusPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	return os.WriteFile(r.statusPath, data, 0o600)
}
