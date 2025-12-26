package diff

import (
	"context"
	"errors"
	"testing"

	"syntaxia/domain"
)

// mockDiffEngine реализует domain.DiffEngine для тестов
type mockDiffEngine struct {
	generateDiffFunc            func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error)
	generateDiffFromResultsFunc func(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error)
	generateDiffFromEditsFunc   func(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error)
	publishDiffFunc             func(ctx context.Context, diff *domain.DiffResult) error
}

func (m *mockDiffEngine) GenerateDiff(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
	if m.generateDiffFunc != nil {
		return m.generateDiffFunc(ctx, beforePath, afterPath, format)
	}
	return &domain.DiffResult{ID: "test-diff-id", Format: format}, nil
}

func (m *mockDiffEngine) GenerateDiffFromResults(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
	if m.generateDiffFromResultsFunc != nil {
		return m.generateDiffFromResultsFunc(ctx, results, format)
	}
	return &domain.DiffResult{ID: "test-diff-id", Format: format}, nil
}

func (m *mockDiffEngine) GenerateDiffFromEdits(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
	if m.generateDiffFromEditsFunc != nil {
		return m.generateDiffFromEditsFunc(ctx, edits, format)
	}
	return &domain.DiffResult{ID: "test-diff-id", Format: format}, nil
}

func (m *mockDiffEngine) PublishDiff(ctx context.Context, diff *domain.DiffResult) error {
	if m.publishDiffFunc != nil {
		return m.publishDiffFunc(ctx, diff)
	}
	return nil
}

func TestNewService(t *testing.T) {
	logger := &domain.NoopLogger{}
	engine := &mockDiffEngine{}

	service := NewService(logger, engine)

	if service == nil {
		t.Fatal("expected service to be created")
	}
	if service.log != logger {
		t.Error("expected logger to be set")
	}
	if service.engine != engine {
		t.Error("expected engine to be set")
	}
}

func TestService_GenerateDiff(t *testing.T) {
	tests := []struct {
		name       string
		beforePath string
		afterPath  string
		format     domain.DiffFormat
		engineFunc func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error)
		wantErr    bool
		wantID     string
	}{
		{
			name:       "successful diff generation",
			beforePath: "/path/to/before.go",
			afterPath:  "/path/to/after.go",
			format:     domain.DiffFormatGit,
			engineFunc: func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
				return &domain.DiffResult{
					ID:     "diff-123",
					Format: format,
					Content: `--- a/file.go
+++ b/file.go
@@ -1,3 +1,4 @@
 package main
+import "fmt"
 func main() {}`,
				}, nil
			},
			wantErr: false,
			wantID:  "diff-123",
		},
		{
			name:       "empty paths",
			beforePath: "",
			afterPath:  "",
			format:     domain.DiffFormatUnified,
			engineFunc: func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
				return &domain.DiffResult{ID: "empty-diff", Format: format}, nil
			},
			wantErr: false,
			wantID:  "empty-diff",
		},
		{
			name:       "engine error",
			beforePath: "/invalid/path",
			afterPath:  "/another/invalid",
			format:     domain.DiffFormatJSON,
			engineFunc: func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("file not found")
			},
			wantErr: true,
		},
		{
			name:       "json format",
			beforePath: "/path/to/file.ts",
			afterPath:  "/path/to/file.ts",
			format:     domain.DiffFormatJSON,
			engineFunc: nil,
			wantErr:    false,
			wantID:     "test-diff-id",
		},
		{
			name:       "html format",
			beforePath: "/path/to/file.vue",
			afterPath:  "/path/to/file.vue",
			format:     domain.DiffFormatHTML,
			engineFunc: nil,
			wantErr:    false,
			wantID:     "test-diff-id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{generateDiffFunc: tt.engineFunc}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateDiff(context.Background(), tt.beforePath, tt.afterPath, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.ID != tt.wantID {
				t.Errorf("GenerateDiff() ID = %v, want %v", result.ID, tt.wantID)
			}
		})
	}
}

func TestService_GenerateDiffFromResults(t *testing.T) {
	tests := []struct {
		name       string
		results    []*domain.ApplyResult
		format     domain.DiffFormat
		engineFunc func(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error)
		wantErr    bool
	}{
		{
			name:    "empty results",
			results: []*domain.ApplyResult{},
			format:  domain.DiffFormatGit,
			wantErr: false,
		},
		{
			name: "single successful result",
			results: []*domain.ApplyResult{
				{Success: true, Path: "/path/to/file.go", OperationID: "op-1"},
			},
			format:  domain.DiffFormatUnified,
			wantErr: false,
		},
		{
			name: "multiple results mixed success",
			results: []*domain.ApplyResult{
				{Success: true, Path: "/path/to/file1.go", OperationID: "op-1"},
				{Success: false, Path: "/path/to/file2.go", OperationID: "op-2", Error: "failed"},
				{Success: true, Path: "/path/to/file3.go", OperationID: "op-3"},
			},
			format:  domain.DiffFormatJSON,
			wantErr: false,
		},
		{
			name: "engine error",
			results: []*domain.ApplyResult{
				{Success: true, Path: "/path/to/file.go", OperationID: "op-1"},
			},
			format: domain.DiffFormatGit,
			engineFunc: func(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("engine failure")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{generateDiffFromResultsFunc: tt.engineFunc}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateDiffFromResults(context.Background(), tt.results, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDiffFromResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("GenerateDiffFromResults() expected result, got nil")
			}
		})
	}
}

func TestService_GenerateDiffFromEdits(t *testing.T) {
	tests := []struct {
		name       string
		edits      *domain.EditsJSON
		format     domain.DiffFormat
		engineFunc func(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error)
		wantErr    bool
	}{
		{
			name: "empty edits",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits:         []*domain.Edit{},
			},
			format:  domain.DiffFormatGit,
			wantErr: false,
		},
		{
			name: "single edit",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits: []*domain.Edit{
					{ID: "edit-1", Path: "/path/to/file.go", Kind: "fullFile", Op: "modify", Content: "new content"},
				},
			},
			format:  domain.DiffFormatUnified,
			wantErr: false,
		},
		{
			name: "multiple edits with dependencies",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Metadata: &domain.EditsMetadata{
					Reason:     "refactoring",
					Confidence: 0.95,
				},
				Edits: []*domain.Edit{
					{ID: "edit-1", Path: "/path/to/file1.go", Kind: "fullFile", Op: "modify"},
					{ID: "edit-2", Path: "/path/to/file2.go", Kind: "anchorPatch", Op: "modify", DependsOn: []string{"edit-1"}},
				},
			},
			format:  domain.DiffFormatJSON,
			wantErr: false,
		},
		{
			name: "engine error",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits:         []*domain.Edit{{ID: "edit-1"}},
			},
			format: domain.DiffFormatGit,
			engineFunc: func(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("invalid edits")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{generateDiffFromEditsFunc: tt.engineFunc}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateDiffFromEdits(context.Background(), tt.edits, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateDiffFromEdits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result == nil {
				t.Error("GenerateDiffFromEdits() expected result, got nil")
			}
		})
	}
}

func TestService_PublishDiff(t *testing.T) {
	tests := []struct {
		name       string
		diff       *domain.DiffResult
		engineFunc func(ctx context.Context, diff *domain.DiffResult) error
		wantErr    bool
	}{
		{
			name: "successful publish",
			diff: &domain.DiffResult{
				ID:      "diff-123",
				Format:  domain.DiffFormatGit,
				Content: "diff content",
			},
			wantErr: false,
		},
		{
			name: "publish error",
			diff: &domain.DiffResult{
				ID:     "diff-456",
				Format: domain.DiffFormatJSON,
			},
			engineFunc: func(ctx context.Context, diff *domain.DiffResult) error {
				return errors.New("publish failed")
			},
			wantErr: true,
		},
		{
			name: "empty diff",
			diff: &domain.DiffResult{
				ID:      "",
				Format:  domain.DiffFormatUnified,
				Content: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{publishDiffFunc: tt.engineFunc}
			service := NewService(&domain.NoopLogger{}, engine)

			err := service.PublishDiff(context.Background(), tt.diff)

			if (err != nil) != tt.wantErr {
				t.Errorf("PublishDiff() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_GenerateAndPublishDiff(t *testing.T) {
	tests := []struct {
		name           string
		beforePath     string
		afterPath      string
		format         domain.DiffFormat
		generateFunc   func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error)
		publishFunc    func(ctx context.Context, diff *domain.DiffResult) error
		wantErr        bool
		wantResultNil  bool
	}{
		{
			name:       "successful generate and publish",
			beforePath: "/before.go",
			afterPath:  "/after.go",
			format:     domain.DiffFormatGit,
			wantErr:    false,
		},
		{
			name:       "generate fails",
			beforePath: "/invalid",
			afterPath:  "/invalid",
			format:     domain.DiffFormatGit,
			generateFunc: func(ctx context.Context, beforePath, afterPath string, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("generate failed")
			},
			wantErr:       true,
			wantResultNil: true,
		},
		{
			name:       "publish fails but returns result",
			beforePath: "/before.go",
			afterPath:  "/after.go",
			format:     domain.DiffFormatUnified,
			publishFunc: func(ctx context.Context, diff *domain.DiffResult) error {
				return errors.New("publish failed")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{
				generateDiffFunc: tt.generateFunc,
				publishDiffFunc:  tt.publishFunc,
			}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateAndPublishDiff(context.Background(), tt.beforePath, tt.afterPath, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAndPublishDiff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResultNil && result != nil {
				t.Error("GenerateAndPublishDiff() expected nil result")
			}
			if !tt.wantResultNil && !tt.wantErr && result == nil {
				t.Error("GenerateAndPublishDiff() expected result, got nil")
			}
		})
	}
}

func TestService_GenerateAndPublishDiffFromResults(t *testing.T) {
	tests := []struct {
		name          string
		results       []*domain.ApplyResult
		format        domain.DiffFormat
		generateFunc  func(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error)
		publishFunc   func(ctx context.Context, diff *domain.DiffResult) error
		wantErr       bool
		wantResultNil bool
	}{
		{
			name: "successful generate and publish",
			results: []*domain.ApplyResult{
				{Success: true, Path: "/file.go", OperationID: "op-1"},
			},
			format:  domain.DiffFormatGit,
			wantErr: false,
		},
		{
			name:    "generate fails",
			results: []*domain.ApplyResult{},
			format:  domain.DiffFormatJSON,
			generateFunc: func(ctx context.Context, results []*domain.ApplyResult, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("generate failed")
			},
			wantErr:       true,
			wantResultNil: true,
		},
		{
			name: "publish fails but returns result",
			results: []*domain.ApplyResult{
				{Success: true, Path: "/file.go", OperationID: "op-1"},
			},
			format: domain.DiffFormatUnified,
			publishFunc: func(ctx context.Context, diff *domain.DiffResult) error {
				return errors.New("publish failed")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{
				generateDiffFromResultsFunc: tt.generateFunc,
				publishDiffFunc:             tt.publishFunc,
			}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateAndPublishDiffFromResults(context.Background(), tt.results, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAndPublishDiffFromResults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResultNil && result != nil {
				t.Error("GenerateAndPublishDiffFromResults() expected nil result")
			}
			if !tt.wantResultNil && !tt.wantErr && result == nil {
				t.Error("GenerateAndPublishDiffFromResults() expected result, got nil")
			}
		})
	}
}

func TestService_GenerateAndPublishDiffFromEdits(t *testing.T) {
	tests := []struct {
		name          string
		edits         *domain.EditsJSON
		format        domain.DiffFormat
		generateFunc  func(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error)
		publishFunc   func(ctx context.Context, diff *domain.DiffResult) error
		wantErr       bool
		wantResultNil bool
	}{
		{
			name: "successful generate and publish",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits:         []*domain.Edit{{ID: "edit-1", Path: "/file.go"}},
			},
			format:  domain.DiffFormatGit,
			wantErr: false,
		},
		{
			name: "generate fails",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits:         []*domain.Edit{},
			},
			format: domain.DiffFormatJSON,
			generateFunc: func(ctx context.Context, edits *domain.EditsJSON, format domain.DiffFormat) (*domain.DiffResult, error) {
				return nil, errors.New("generate failed")
			},
			wantErr:       true,
			wantResultNil: true,
		},
		{
			name: "publish fails but returns result",
			edits: &domain.EditsJSON{
				SchemaVersion: "1.0",
				Edits:         []*domain.Edit{{ID: "edit-1"}},
			},
			format: domain.DiffFormatHTML,
			publishFunc: func(ctx context.Context, diff *domain.DiffResult) error {
				return errors.New("publish failed")
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiffEngine{
				generateDiffFromEditsFunc: tt.generateFunc,
				publishDiffFunc:           tt.publishFunc,
			}
			service := NewService(&domain.NoopLogger{}, engine)

			result, err := service.GenerateAndPublishDiffFromEdits(context.Background(), tt.edits, tt.format)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAndPublishDiffFromEdits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResultNil && result != nil {
				t.Error("GenerateAndPublishDiffFromEdits() expected nil result")
			}
			if !tt.wantResultNil && !tt.wantErr && result == nil {
				t.Error("GenerateAndPublishDiffFromEdits() expected result, got nil")
			}
		})
	}
}

func TestService_GetSupportedFormats(t *testing.T) {
	service := NewService(&domain.NoopLogger{}, &mockDiffEngine{})

	formats := service.GetSupportedFormats()

	expectedFormats := []domain.DiffFormat{
		domain.DiffFormatGit,
		domain.DiffFormatUnified,
		domain.DiffFormatJSON,
		domain.DiffFormatHTML,
	}

	if len(formats) != len(expectedFormats) {
		t.Errorf("GetSupportedFormats() returned %d formats, want %d", len(formats), len(expectedFormats))
		return
	}

	for i, format := range formats {
		if format != expectedFormats[i] {
			t.Errorf("GetSupportedFormats()[%d] = %v, want %v", i, format, expectedFormats[i])
		}
	}
}
