package analysis

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"syntaxia/infrastructure/analyzers"
)

// mockLogger implements domain.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Info(msg string)                                         {}
func (m *mockLogger) Warning(msg string)                                      {}
func (m *mockLogger) Error(msg string)                                        {}
func (m *mockLogger) Debug(msg string)                                        {}
func (m *mockLogger) Fatal(msg string)                                        {}
func (m *mockLogger) InfoWithFields(msg string, fields map[string]interface{}) {}
func (m *mockLogger) WarningWithFields(msg string, fields map[string]interface{}) {}
func (m *mockLogger) ErrorWithFields(msg string, fields map[string]interface{}) {}
func (m *mockLogger) DebugWithFields(msg string, fields map[string]interface{}) {}

func TestDependencyAnalyzer_AnalyzeFile_TypeScript(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test files
	componentFile := filepath.Join(tmpDir, "Component.ts")
	componentContent := `import { ref } from 'vue'
import { useStore } from './store'
import type { User } from './types'

export default {
  name: 'Component'
}`

	utilFile := filepath.Join(tmpDir, "store.ts")
	utilContent := `export const useStore = () => {
  return { state: {} }
}`

	typesFile := filepath.Join(tmpDir, "types.ts")
	typesContent := `export interface User {
  id: number
  name: string
}`

	if err := os.WriteFile(componentFile, []byte(componentContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(utilFile, []byte(utilContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(typesFile, []byte(typesContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create analyzer
	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Analyze file
	deps, err := analyzer.AnalyzeFile(tmpDir, "Component.ts")
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Verify dependencies
	if len(deps) < 2 {
		t.Errorf("Expected at least 2 dependencies, got %d", len(deps))
	}

	// Check for store import
	foundStore := false
	for _, dep := range deps {
		if dep.TargetPath == "store.ts" {
			foundStore = true
			if dep.Type != "import" {
				t.Errorf("Expected import type, got %s", dep.Type)
			}
		}
		// Note: TypeScript type imports are detected as regular imports
		// since the Import struct doesn't distinguish type-only imports
	}

	if !foundStore {
		t.Error("Expected to find store.ts dependency")
	}
}

func TestDependencyAnalyzer_AnalyzeFile_JavaScript(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	mainFile := filepath.Join(tmpDir, "main.js")
	mainContent := `import { helper } from './utils'
import './styles.css'

export function main() {
  return helper()
}`

	utilFile := filepath.Join(tmpDir, "utils.js")
	utilContent := `export function helper() {
  return 'hello'
}`

	styleFile := filepath.Join(tmpDir, "styles.css")
	styleContent := `.container { color: red; }`

	if err := os.WriteFile(mainFile, []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(utilFile, []byte(utilContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(styleFile, []byte(styleContent), 0644); err != nil {
		t.Fatal(err)
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	deps, err := analyzer.AnalyzeFile(tmpDir, "main.js")
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Should have at least utils dependency (CSS might not resolve if no analyzer)
	if len(deps) < 1 {
		t.Errorf("Expected at least 1 dependency, got %d", len(deps))
	}

	// Check for utils import
	foundUtils := false
	for _, dep := range deps {
		if dep.TargetPath == "utils.js" {
			foundUtils = true
		}
	}

	if !foundUtils {
		t.Error("Expected to find utils.js dependency")
	}
}

func TestDependencyAnalyzer_BuildGraph(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a small project structure
	files := map[string]string{
		"index.ts": `import { App } from './App'
export { App }`,
		"App.ts": `import { Component } from './components/Component'
export const App = { Component }`,
		"components/Component.ts": `import { helper } from '../utils'
export const Component = {}`,
		"utils.ts": `export function helper() {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	graph, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	if graph == nil {
		t.Fatal("Expected non-nil graph")
	}

	// Check that all files are in the graph
	expectedFiles := []string{"index.ts", "App.ts", "components/Component.ts", "utils.ts"}
	for _, file := range expectedFiles {
		if _, ok := graph.Files[file]; !ok {
			t.Errorf("Expected file %s in graph", file)
		}
	}

	// Verify index.ts depends on App.ts
	indexDeps := graph.Files["index.ts"]
	foundApp := false
	for _, dep := range indexDeps {
		if dep.TargetPath == "App.ts" {
			foundApp = true
			break
		}
	}
	if !foundApp {
		t.Error("Expected index.ts to depend on App.ts")
	}
}

func TestDependencyAnalyzer_ClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a simple file
	testFile := filepath.Join(tmpDir, "test.ts")
	if err := os.WriteFile(testFile, []byte("export const x = 1"), 0644); err != nil {
		t.Fatal(err)
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph to populate cache
	_, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Verify cache exists
	if !analyzer.IsDependencyGraphCached(tmpDir) {
		t.Error("Expected graph to be cached")
	}

	// Clear cache - should not error
	analyzer.ClearCache()

	// Verify cache is cleared
	if analyzer.IsDependencyGraphCached(tmpDir) {
		t.Error("Expected cache to be cleared")
	}

	// Build again to verify it still works
	graph, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph after ClearCache failed: %v", err)
	}
	if graph == nil {
		t.Error("Expected non-nil graph after cache clear")
	}
}

func TestDependencyAnalyzer_GetFileDependenciesBatch(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.ts": `import { App } from './App'
export { App }`,
		"App.ts": `import { Component } from './components/Component'
export const App = { Component }`,
		"components/Component.ts": `import { helper } from '../utils'
export const Component = {}`,
		"utils.ts": `export function helper() {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph first
	_, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Test batch query
	filePaths := []string{"index.ts", "App.ts", "utils.ts"}
	result, err := analyzer.GetFileDependenciesBatch(tmpDir, filePaths)
	if err != nil {
		t.Fatalf("GetFileDependenciesBatch failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// Verify index.ts has dependencies
	if deps, ok := result["index.ts"]; !ok || len(deps) == 0 {
		t.Error("Expected index.ts to have dependencies")
	}

	// Verify utils.ts has no dependencies
	if deps, ok := result["utils.ts"]; !ok || len(deps) != 0 {
		t.Errorf("Expected utils.ts to have 0 dependencies, got %d", len(deps))
	}
}

func TestDependencyAnalyzer_GetIncomingDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.ts": `import { App } from './App'`,
		"App.ts":   `import { helper } from './utils'`,
		"main.ts":  `import { helper } from './utils'`,
		"utils.ts": `export function helper() {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph
	_, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Get incoming dependencies for utils.ts
	incoming, err := analyzer.GetIncomingDependencies(tmpDir, "utils.ts")
	if err != nil {
		t.Fatalf("GetIncomingDependencies failed: %v", err)
	}

	// Should have 2 files importing utils.ts (App.ts and main.ts)
	if len(incoming) != 2 {
		t.Errorf("Expected 2 incoming dependencies, got %d", len(incoming))
	}

	// Verify both App.ts and main.ts are in the list
	hasApp := false
	hasMain := false
	for _, file := range incoming {
		if file == "App.ts" {
			hasApp = true
		}
		if file == "main.ts" {
			hasMain = true
		}
	}

	if !hasApp || !hasMain {
		t.Error("Expected both App.ts and main.ts in incoming dependencies")
	}
}

func TestDependencyAnalyzer_GetOutgoingDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.ts": `import { App } from './App'
import { helper } from './utils'`,
		"App.ts":   `export const App = {}`,
		"utils.ts": `export function helper() {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph
	_, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Get outgoing dependencies for index.ts
	outgoing, err := analyzer.GetOutgoingDependencies(tmpDir, "index.ts")
	if err != nil {
		t.Fatalf("GetOutgoingDependencies failed: %v", err)
	}

	// Should have 2 outgoing dependencies
	if len(outgoing) != 2 {
		t.Errorf("Expected 2 outgoing dependencies, got %d", len(outgoing))
	}

	// Verify App.ts and utils.ts are in the list
	hasApp := false
	hasUtils := false
	for _, file := range outgoing {
		if file == "App.ts" {
			hasApp = true
		}
		if file == "utils.ts" {
			hasUtils = true
		}
	}

	if !hasApp || !hasUtils {
		t.Error("Expected both App.ts and utils.ts in outgoing dependencies")
	}
}

func TestDependencyAnalyzer_GetDependencyGraphStats(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.ts": `import { App } from './App'`,
		"App.ts":   `import { helper } from './utils'`,
		"utils.ts": `export function helper() {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Get stats before building graph
	stats, err := analyzer.GetDependencyGraphStats(tmpDir)
	if err != nil {
		t.Fatalf("GetDependencyGraphStats failed: %v", err)
	}

	if stats.IsCached {
		t.Error("Expected graph to not be cached initially")
	}

	// Build graph
	_, err = analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Get stats after building
	stats, err = analyzer.GetDependencyGraphStats(tmpDir)
	if err != nil {
		t.Fatalf("GetDependencyGraphStats failed: %v", err)
	}

	if !stats.IsCached {
		t.Error("Expected graph to be cached")
	}

	if stats.FileCount != 3 {
		t.Errorf("Expected 3 files, got %d", stats.FileCount)
	}

	if stats.DependencyCount != 2 {
		t.Errorf("Expected 2 dependencies, got %d", stats.DependencyCount)
	}

	if stats.CacheSize <= 0 {
		t.Error("Expected positive cache size")
	}
}

func TestDependencyAnalyzer_UpdateFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create initial file
	testFile := filepath.Join(tmpDir, "test.ts")
	if err := os.WriteFile(testFile, []byte("export const x = 1"), 0644); err != nil {
		t.Fatal(err)
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph
	graph, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	initialTime := graph.LastAnalyzed

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	// Update file content
	newContent := `import { helper } from './utils'
export const x = 1`
	if err := os.WriteFile(testFile, []byte(newContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Create utils file
	utilsFile := filepath.Join(tmpDir, "utils.ts")
	if err := os.WriteFile(utilsFile, []byte("export function helper() {}"), 0644); err != nil {
		t.Fatal(err)
	}

	// Update file in cache
	err = analyzer.UpdateFile(tmpDir, "test.ts")
	if err != nil {
		t.Fatalf("UpdateFile failed: %v", err)
	}

	// Get updated stats
	stats, err := analyzer.GetDependencyGraphStats(tmpDir)
	if err != nil {
		t.Fatalf("GetDependencyGraphStats failed: %v", err)
	}

	// Verify LastAnalyzed was updated
	if !stats.LastAnalyzed.After(initialTime) {
		t.Error("Expected LastAnalyzed to be updated")
	}

	// Verify dependencies were updated
	outgoing, err := analyzer.GetOutgoingDependencies(tmpDir, "test.ts")
	if err != nil {
		t.Fatalf("GetOutgoingDependencies failed: %v", err)
	}

	if len(outgoing) != 1 || outgoing[0] != "utils.ts" {
		t.Error("Expected test.ts to depend on utils.ts after update")
	}
}

func TestDependencyAnalyzer_RemoveFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"index.ts": `import { App } from './App'`,
		"App.ts":   `export const App = {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Build graph
	_, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Verify initial file count
	stats, err := analyzer.GetDependencyGraphStats(tmpDir)
	if err != nil {
		t.Fatalf("GetDependencyGraphStats failed: %v", err)
	}

	if stats.FileCount != 2 {
		t.Errorf("Expected 2 files initially, got %d", stats.FileCount)
	}

	// Remove App.ts from cache
	err = analyzer.RemoveFile(tmpDir, "App.ts")
	if err != nil {
		t.Fatalf("RemoveFile failed: %v", err)
	}

	// Verify file count decreased
	stats, err = analyzer.GetDependencyGraphStats(tmpDir)
	if err != nil {
		t.Fatalf("GetDependencyGraphStats failed: %v", err)
	}

	if stats.FileCount != 1 {
		t.Errorf("Expected 1 file after removal, got %d", stats.FileCount)
	}

	// Verify App.ts is not in outgoing dependencies
	outgoing, err := analyzer.GetOutgoingDependencies(tmpDir, "App.ts")
	if err != nil {
		t.Fatalf("GetOutgoingDependencies failed: %v", err)
	}

	if len(outgoing) != 0 {
		t.Error("Expected no dependencies for removed file")
	}
}

func TestDependencyAnalyzer_IgnorePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files in different directories
	files := map[string]string{
		"src/index.ts":              "export const x = 1",
		"node_modules/package/index.js": "module.exports = {}",
		"dist/bundle.js":            "// compiled",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	graph, err := analyzer.BuildGraph(tmpDir, nil)
	if err != nil {
		t.Fatalf("BuildGraph failed: %v", err)
	}

	// Should only include src/index.ts, not node_modules or dist
	if _, ok := graph.Files["node_modules/package/index.js"]; ok {
		t.Error("node_modules should be ignored")
	}
	if _, ok := graph.Files["dist/bundle.js"]; ok {
		t.Error("dist should be ignored")
	}
	if _, ok := graph.Files["src/index.ts"]; !ok {
		t.Error("src/index.ts should be included")
	}
}

func TestDependencyAnalyzer_ConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test file
	testFile := filepath.Join(tmpDir, "test.ts")
	if err := os.WriteFile(testFile, []byte("export const x = 1"), 0644); err != nil {
		t.Fatal(err)
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	// Test concurrent BuildGraph and ClearCache calls
	done := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			_, _ = analyzer.BuildGraph(tmpDir, nil)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 10; i++ {
			analyzer.ClearCache()
		}
		done <- true
	}()

	<-done
	<-done
}

func TestDependencyAnalyzer_UnsupportedFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create unsupported file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("plain text"), 0644); err != nil {
		t.Fatal(err)
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	_, err := analyzer.AnalyzeFile(tmpDir, "test.txt")
	if err == nil {
		t.Error("Expected error for unsupported file type")
	}
}

func TestDependencyAnalyzer_RelativeImports(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested structure
	files := map[string]string{
		"src/features/auth/Login.ts": `import { Button } from '../../components/Button'
export const Login = {}`,
		"src/components/Button.ts": `export const Button = {}`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	registry := analyzers.NewAnalyzerRegistry()
	logger := &mockLogger{}
	analyzer := NewDependencyAnalyzer(registry, logger)

	deps, err := analyzer.AnalyzeFile(tmpDir, "src/features/auth/Login.ts")
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	// Should resolve relative import correctly
	found := false
	for _, dep := range deps {
		if dep.TargetPath == "src/components/Button.ts" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected to resolve relative import to src/components/Button.ts")
	}
}
