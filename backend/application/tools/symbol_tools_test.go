package tools

import (
	"strings"
	"testing"

	"syntaxia/testutils"
)

func TestSymbolToolsHandler_ListSymbols(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing path returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "path is required",
		},
		{
			name:        "empty path returns error",
			args:        map[string]any{"path": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "path is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.ListSymbols(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ListSymbols() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("ListSymbols() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("ListSymbols() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_SearchSymbols(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing query returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "query is required",
		},
		{
			name:        "empty query returns error",
			args:        map[string]any{"query": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "query is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.SearchSymbols(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SearchSymbols() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("SearchSymbols() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("SearchSymbols() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_FindDefinition(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing name returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:        "empty name returns error",
			args:        map[string]any{"name": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.FindDefinition(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FindDefinition() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("FindDefinition() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("FindDefinition() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_FindReferences(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing name returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:        "empty name returns error",
			args:        map[string]any{"name": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.FindReferences(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FindReferences() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("FindReferences() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("FindReferences() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_GetSymbolInfo(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing name returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
		{
			name:        "empty name returns error",
			args:        map[string]any{"name": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "name is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.GetSymbolInfo(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSymbolInfo() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetSymbolInfo() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("GetSymbolInfo() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_GetClassHierarchy(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		projectRoot string
		wantErr     bool
		errContains string
	}{
		{
			name:        "missing class_name returns error",
			args:        map[string]any{},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "class_name is required",
		},
		{
			name:        "empty class_name returns error",
			args:        map[string]any{"class_name": ""},
			projectRoot: "/tmp",
			wantErr:     true,
			errContains: "class_name is required",
		},
	}

	handler := &SymbolToolsHandler{
		logger: testutils.NewMockLogger(),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.GetClassHierarchy(tt.args, tt.projectRoot)
			if tt.wantErr {
				if err == nil {
					t.Errorf("GetClassHierarchy() expected error containing %q, got nil", tt.errContains)
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetClassHierarchy() error = %v, want error containing %q", err, tt.errContains)
				}
			} else if err != nil {
				t.Errorf("GetClassHierarchy() unexpected error = %v", err)
			}
		})
	}
}

func TestSymbolToolsHandler_NewSymbolToolsHandler(t *testing.T) {
	// Test constructor creates handler with all dependencies
	logger := testutils.NewMockLogger()
	handler := NewSymbolToolsHandler(nil, nil, logger, nil)

	if handler == nil {
		t.Fatal("NewSymbolToolsHandler() returned nil")
	}
	if handler.logger == nil {
		t.Error("NewSymbolToolsHandler() logger is nil")
	}
}
