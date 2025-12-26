package symbolgraph

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"syntaxia/infrastructure/analyzers"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string)                          {}
func (m *mockLogger) Info(msg string)                           {}
func (m *mockLogger) Warning(msg string)                        {}
func (m *mockLogger) Error(msg string)                          {}
func (m *mockLogger) Fatal(msg string)                          {}
func (m *mockLogger) Debugf(format string, args ...interface{}) {}
func (m *mockLogger) Infof(format string, args ...interface{})  {}
func (m *mockLogger) Warnf(format string, args ...interface{})  {}
func (m *mockLogger) Errorf(format string, args ...interface{}) {}

func setupTestProject(t *testing.T) string {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "ts-import-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create directory structure
	dirs := []string{
		"src",
		"src/components",
		"src/composables",
		"src/features",
		"src/stores",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, dir), 0o755); err != nil {
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	return tmpDir
}

func writeFile(t *testing.T, dir, path, content string) {
	t.Helper()
	fullPath := filepath.Join(dir, path)
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}

func TestTSImportGraphBuilder_SimpleImport(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create test files
	writeFile(t, tmpDir, "src/main.ts", `
import { helper } from './utils'

export function main() {
    return helper()
}
`)

	writeFile(t, tmpDir, "src/utils.ts", `
export function helper() {
    return 'hello'
}
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Verify graph
	if len(graph.Packages) < 2 {
		t.Errorf("Expected at least 2 packages, got %d", len(graph.Packages))
	}

	// Check that main.ts imports utils.ts
	foundImport := false
	for _, imp := range graph.Imports {
		if imp.From == "src/main.ts" && (imp.To == "src/utils.ts" || imp.To == "src/utils") {
			foundImport = true
			break
		}
	}

	if !foundImport {
		t.Error("Expected import from main.ts to utils.ts not found")
	}
}

func TestTSImportGraphBuilder_PathAlias(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create tsconfig.json with path aliases
	writeFile(t, tmpDir, "tsconfig.json", `{
    "compilerOptions": {
        "baseUrl": ".",
        "paths": {
            "@/*": ["src/*"]
        }
    }
}`)

	// Create test files
	writeFile(t, tmpDir, "src/main.ts", `
import { Button } from '@/components/Button'

export function App() {
    return Button()
}
`)

	writeFile(t, tmpDir, "src/components/Button.ts", `
export function Button() {
    return 'button'
}
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Check that path alias was resolved
	foundImport := false
	for _, imp := range graph.Imports {
		if imp.From == "src/main.ts" {
			// Should resolve @/components/Button to src/components/Button.ts
			if strings.Contains(imp.To, "components/Button") {
				foundImport = true
				break
			}
		}
	}

	if !foundImport {
		t.Error("Expected path alias @/components/Button to be resolved")
		for _, imp := range graph.Imports {
			t.Logf("Import: %s -> %s", imp.From, imp.To)
		}
	}
}

func TestTSImportGraphBuilder_BarrelExport(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create tsconfig.json
	writeFile(t, tmpDir, "tsconfig.json", `{
    "compilerOptions": {
        "baseUrl": ".",
        "paths": {
            "@/*": ["src/*"]
        }
    }
}`)

	// Create barrel export (index.ts)
	writeFile(t, tmpDir, "src/features/index.ts", `
export { useFileStore } from './files'
export { useContextStore } from './context'
`)

	writeFile(t, tmpDir, "src/features/files.ts", `
export function useFileStore() {
    return {}
}
`)

	writeFile(t, tmpDir, "src/features/context.ts", `
export function useContextStore() {
    return {}
}
`)

	// Create file that imports from barrel
	writeFile(t, tmpDir, "src/main.ts", `
import { useFileStore } from '@/features'

export function App() {
    return useFileStore()
}
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Check that barrel import was resolved to index.ts or features directory
	foundBarrelImport := false
	for _, imp := range graph.Imports {
		if imp.From == "src/main.ts" {
			if strings.Contains(imp.To, "features") {
				foundBarrelImport = true
				break
			}
		}
	}

	if !foundBarrelImport {
		t.Error("Expected barrel import @/features to be resolved")
		for _, imp := range graph.Imports {
			t.Logf("Import: %s -> %s", imp.From, imp.To)
		}
	}
}

func TestTSImportGraphBuilder_VueSFC(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create tsconfig.json
	writeFile(t, tmpDir, "tsconfig.json", `{
    "compilerOptions": {
        "baseUrl": ".",
        "paths": {
            "@/*": ["src/*"]
        }
    }
}`)

	// Create Vue SFC
	writeFile(t, tmpDir, "src/components/ChatPanel.vue", `
<script setup lang="ts">
import { ref } from 'vue'
import { useChatStore } from '@/stores/chat'
import MessageItem from './MessageItem.vue'
import { useI18n } from '@/composables/useI18n'

const store = useChatStore()
const { t } = useI18n()
</script>

<template>
    <div>
        <MessageItem v-for="msg in store.messages" :key="msg.id" :message="msg" />
    </div>
</template>
`)

	writeFile(t, tmpDir, "src/stores/chat.ts", `
export function useChatStore() {
    return { messages: [] }
}
`)

	writeFile(t, tmpDir, "src/components/MessageItem.vue", `
<script setup lang="ts">
defineProps<{ message: any }>()
</script>

<template>
    <div>{{ message.text }}</div>
</template>
`)

	writeFile(t, tmpDir, "src/composables/useI18n.ts", `
export function useI18n() {
    return { t: (key: string) => key }
}
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Check that Vue SFC imports are extracted
	chatPanelImports := make([]string, 0)
	for _, imp := range graph.Imports {
		if imp.From == "src/components/ChatPanel.vue" {
			chatPanelImports = append(chatPanelImports, imp.To)
		}
	}

	if len(chatPanelImports) < 3 {
		t.Errorf("Expected at least 3 imports from ChatPanel.vue, got %d: %v", len(chatPanelImports), chatPanelImports)
	}
}

func TestTSImportGraphBuilder_CircularImports(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create circular dependency: A -> B -> C -> A
	writeFile(t, tmpDir, "src/a.ts", `
import { b } from './b'
export const a = () => b()
`)

	writeFile(t, tmpDir, "src/b.ts", `
import { c } from './c'
export const b = () => c()
`)

	writeFile(t, tmpDir, "src/c.ts", `
import { a } from './a'
export const c = () => a()
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Find circular imports
	cycles, err := builder.GetCircularImports(ctx, graph)
	if err != nil {
		t.Fatalf("GetCircularImports failed: %v", err)
	}

	if len(cycles) == 0 {
		t.Error("Expected to find circular imports, but found none")
		t.Logf("Graph imports: %+v", graph.Imports)
	} else {
		t.Logf("Found %d circular import cycles", len(cycles))
		for i, cycle := range cycles {
			t.Logf("Cycle %d: %v", i+1, cycle)
		}
	}
}

func TestTSImportGraphBuilder_GetImportPath(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create chain: main -> utils -> helpers
	writeFile(t, tmpDir, "src/main.ts", `
import { util } from './utils'
export const main = () => util()
`)

	writeFile(t, tmpDir, "src/utils.ts", `
import { helper } from './helpers'
export const util = () => helper()
`)

	writeFile(t, tmpDir, "src/helpers.ts", `
export const helper = () => 'help'
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Find path from main.ts to helpers.ts
	path, err := builder.GetImportPath(ctx, "src/main.ts", "src/helpers.ts", graph)
	if err != nil {
		t.Fatalf("GetImportPath failed: %v", err)
	}

	if len(path) != 3 {
		t.Errorf("Expected path length 3, got %d: %v", len(path), path)
	}

	// Verify path order
	if len(path) >= 3 {
		if path[0] != "src/main.ts" {
			t.Errorf("Expected path to start with src/main.ts, got %s", path[0])
		}
		if path[len(path)-1] != "src/helpers.ts" {
			t.Errorf("Expected path to end with src/helpers.ts, got %s", path[len(path)-1])
		}
	}
}

func TestTSImportGraphBuilder_NoPathFound(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create unconnected files
	writeFile(t, tmpDir, "src/a.ts", `
export const a = 1
`)

	writeFile(t, tmpDir, "src/b.ts", `
export const b = 2
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Try to find path between unconnected files
	_, err = builder.GetImportPath(ctx, "src/a.ts", "src/b.ts", graph)
	if err == nil {
		t.Error("Expected error for no path found, but got nil")
	}
}

func TestTSImportGraphBuilder_ExternalImports(t *testing.T) {
	tmpDir := setupTestProject(t)
	defer os.RemoveAll(tmpDir)

	// Create file with external imports
	writeFile(t, tmpDir, "src/main.ts", `
import { ref, computed } from 'vue'
import axios from 'axios'
import { helper } from './utils'

export function App() {
    const count = ref(0)
    return { count }
}
`)

	writeFile(t, tmpDir, "src/utils.ts", `
export function helper() {
    return 'help'
}
`)

	// Build import graph
	log := &mockLogger{}
	registry := analyzers.NewAnalyzerRegistry()
	builder := NewTSImportGraphBuilder(log, registry)

	ctx := context.Background()
	graph, err := builder.BuildImportGraph(ctx, tmpDir)
	if err != nil {
		t.Fatalf("BuildImportGraph failed: %v", err)
	}

	// Check that external imports are marked correctly
	externalCount := 0
	localCount := 0
	for _, imp := range graph.Imports {
		if imp.From == "src/main.ts" {
			if imp.Type == "external" {
				externalCount++
			} else if imp.Type == "direct" {
				localCount++
			}
		}
	}

	if externalCount < 2 {
		t.Errorf("Expected at least 2 external imports (vue, axios), got %d", externalCount)
	}

	if localCount < 1 {
		t.Errorf("Expected at least 1 local import (utils), got %d", localCount)
	}
}
