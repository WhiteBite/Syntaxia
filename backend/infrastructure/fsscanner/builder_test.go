package fsscanner

import (
	"os"
	"path/filepath"
	"syntaxia/domain"
	"testing"
)

type fakeSettingsRepo struct {
	custom string
}

func (f *fakeSettingsRepo) GetCustomIgnoreRules() string    { return f.custom }
func (f *fakeSettingsRepo) SetCustomIgnoreRules(r string)   { f.custom = r }
func (f *fakeSettingsRepo) GetCustomPromptRules() string    { return "" }
func (f *fakeSettingsRepo) SetCustomPromptRules(string)     {}
func (f *fakeSettingsRepo) GetOpenAIKey() string            { return "" }
func (f *fakeSettingsRepo) SetOpenAIKey(string)             {}
func (f *fakeSettingsRepo) GetGeminiKey() string            { return "" }
func (f *fakeSettingsRepo) SetGeminiKey(string)             {}
func (f *fakeSettingsRepo) GetOpenRouterKey() string        { return "" }
func (f *fakeSettingsRepo) SetOpenRouterKey(string)         {}
func (f *fakeSettingsRepo) GetLocalAIKey() string           { return "" }
func (f *fakeSettingsRepo) SetLocalAIKey(string)            {}
func (f *fakeSettingsRepo) GetLocalAIHost() string          { return "" }
func (f *fakeSettingsRepo) SetLocalAIHost(string)           {}
func (f *fakeSettingsRepo) GetLocalAIModelName() string     { return "" }
func (f *fakeSettingsRepo) SetLocalAIModelName(string)      {}
func (f *fakeSettingsRepo) GetQwenKey() string              { return "" }
func (f *fakeSettingsRepo) SetQwenKey(string)               {}
func (f *fakeSettingsRepo) GetQwenHost() string             { return "" }
func (f *fakeSettingsRepo) SetQwenHost(string)              {}
func (f *fakeSettingsRepo) GetSelectedAIProvider() string   { return "" }
func (f *fakeSettingsRepo) SetSelectedAIProvider(string)    {}
func (f *fakeSettingsRepo) GetSelectedModel(string) string  { return "" }
func (f *fakeSettingsRepo) SetSelectedModel(string, string) {}
func (f *fakeSettingsRepo) GetModels(string) []string       { return nil }
func (f *fakeSettingsRepo) SetModels(string, []string)      {}
func (f *fakeSettingsRepo) GetUseGitignore() bool           { return true }
func (f *fakeSettingsRepo) SetUseGitignore(bool)            {}
func (f *fakeSettingsRepo) GetUseCustomIgnore() bool        { return true }
func (f *fakeSettingsRepo) SetUseCustomIgnore(bool)         {}
func (f *fakeSettingsRepo) GetRecentProjects() []domain.RecentProjectInfo {
	return nil
}
func (f *fakeSettingsRepo) AddRecentProject(path, name string) {}
func (f *fakeSettingsRepo) RemoveRecentProject(path string)    {}
func (f *fakeSettingsRepo) Save() error                        { return nil }
func (f *fakeSettingsRepo) GetSettingsDTO() (domain.SettingsDTO, error) {
	return domain.SettingsDTO{}, nil
}

func collectRelPaths(nodes []*domain.FileNode) map[string]bool {
	res := map[string]bool{}
	var walk func(n *domain.FileNode)
	walk = func(n *domain.FileNode) {
		if n.RelPath != "" {
			res[n.RelPath] = true
		}
		for _, c := range n.Children {
			walk(c)
		}
	}
	for _, n := range nodes {
		walk(n)
	}
	return res
}

func TestBuildTree_CustomIgnore(t *testing.T) {
	dir := t.TempDir()
	// Structure:
	// node_modules/pkg/mod.go (ignored by custom)
	// kept.txt (kept)
	// ignored.txt (ignored by custom)
	if err := os.MkdirAll(filepath.Join(dir, "node_modules", "pkg"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "node_modules", "pkg", "mod.go"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "kept.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "ignored.txt"), []byte("no"), 0o644); err != nil {
		t.Fatal(err)
	}

	repo := &fakeSettingsRepo{custom: "node_modules/\nignored.txt\n"}
	builder := New(repo, &domain.NoopLogger{})
	nodes, err := builder.BuildTree(dir, true, true)
	if err != nil {
		t.Fatalf("BuildTree error: %v", err)
	}
	paths := collectRelPaths(nodes)
	if paths["node_modules"] || paths["node_modules/pkg/mod.go"] {
		t.Errorf("node_modules should be ignored by custom rules")
	}
	if paths["ignored.txt"] {
		t.Errorf("ignored.txt should be ignored by custom rules")
	}
	if !paths["kept.txt"] {
		t.Errorf("kept.txt should exist")
	}
}

func TestBuildTree_MetadataComputation(t *testing.T) {
	dir := t.TempDir()
	
	// Create test structure:
	// root/
	//   ├── file1.go (100 bytes)
	//   ├── file2.ts (200 bytes)
	//   └── subdir/
	//       ├── file3.go (150 bytes)
	//       ├── file4.vue (250 bytes)
	//       └── nested/
	//           └── file5.ts (300 bytes)
	
	if err := os.MkdirAll(filepath.Join(dir, "subdir", "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	
	files := map[string]int{
		"file1.go":                 100,
		"file2.ts":                 200,
		"subdir/file3.go":          150,
		"subdir/file4.vue":         250,
		"subdir/nested/file5.ts":   300,
	}
	
	for path, size := range files {
		content := make([]byte, size)
		if err := os.WriteFile(filepath.Join(dir, path), content, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	
	repo := &fakeSettingsRepo{}
	builder := New(repo, &domain.NoopLogger{})
	nodes, err := builder.BuildTree(dir, false, false)
	if err != nil {
		t.Fatalf("BuildTree error: %v", err)
	}
	
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(nodes))
	}
	
	root := nodes[0]
	
	// Test root metadata
	if root.Depth != 0 {
		t.Errorf("Root depth should be 0, got %d", root.Depth)
	}
	if root.FileCount != 5 {
		t.Errorf("Root FileCount should be 5, got %d", root.FileCount)
	}
	if root.TotalSize != 1000 {
		t.Errorf("Root TotalSize should be 1000, got %d", root.TotalSize)
	}
	if root.DirectFileCount != 2 {
		t.Errorf("Root DirectFileCount should be 2, got %d", root.DirectFileCount)
	}
	
	// Test extension stats
	expectedExts := map[string]int{
		".go":  2,
		".ts":  2,
		".vue": 1,
	}
	if len(root.ExtensionStats) != len(expectedExts) {
		t.Errorf("Root ExtensionStats length mismatch: expected %d, got %d", len(expectedExts), len(root.ExtensionStats))
	}
	for ext, count := range expectedExts {
		if root.ExtensionStats[ext] != count {
			t.Errorf("Root ExtensionStats[%s] should be %d, got %d", ext, count, root.ExtensionStats[ext])
		}
	}
	
	// Find subdir node
	var subdir *domain.FileNode
	for _, child := range root.Children {
		if child.Name == "subdir" && child.IsDir {
			subdir = child
			break
		}
	}
	if subdir == nil {
		t.Fatal("subdir not found")
	}
	
	// Test subdir metadata
	if subdir.Depth != 1 {
		t.Errorf("Subdir depth should be 1, got %d", subdir.Depth)
	}
	if subdir.FileCount != 3 {
		t.Errorf("Subdir FileCount should be 3, got %d", subdir.FileCount)
	}
	if subdir.TotalSize != 700 {
		t.Errorf("Subdir TotalSize should be 700, got %d", subdir.TotalSize)
	}
	if subdir.DirectFileCount != 2 {
		t.Errorf("Subdir DirectFileCount should be 2, got %d", subdir.DirectFileCount)
	}
	
	// Find nested node
	var nested *domain.FileNode
	for _, child := range subdir.Children {
		if child.Name == "nested" && child.IsDir {
			nested = child
			break
		}
	}
	if nested == nil {
		t.Fatal("nested not found")
	}
	
	// Test nested metadata
	if nested.Depth != 2 {
		t.Errorf("Nested depth should be 2, got %d", nested.Depth)
	}
	if nested.FileCount != 1 {
		t.Errorf("Nested FileCount should be 1, got %d", nested.FileCount)
	}
	if nested.TotalSize != 300 {
		t.Errorf("Nested TotalSize should be 300, got %d", nested.TotalSize)
	}
	if nested.DirectFileCount != 1 {
		t.Errorf("Nested DirectFileCount should be 1, got %d", nested.DirectFileCount)
	}
	
	// Test file metadata
	for _, child := range root.Children {
		if !child.IsDir {
			if child.FileCount != 1 {
				t.Errorf("File %s FileCount should be 1, got %d", child.Name, child.FileCount)
			}
			if child.TotalSize != child.Size {
				t.Errorf("File %s TotalSize should equal Size, got %d vs %d", child.Name, child.TotalSize, child.Size)
			}
			if child.DirectFileCount != 0 {
				t.Errorf("File %s DirectFileCount should be 0, got %d", child.Name, child.DirectFileCount)
			}
			if child.ExtensionStats != nil {
				t.Errorf("File %s ExtensionStats should be nil, got %v", child.Name, child.ExtensionStats)
			}
		}
	}
}

func TestBuildTree_EmptyDirectory(t *testing.T) {
	dir := t.TempDir()
	
	repo := &fakeSettingsRepo{}
	builder := New(repo, &domain.NoopLogger{})
	nodes, err := builder.BuildTree(dir, false, false)
	if err != nil {
		t.Fatalf("BuildTree error: %v", err)
	}
	
	if len(nodes) != 1 {
		t.Fatalf("Expected 1 root node, got %d", len(nodes))
	}
	
	root := nodes[0]
	
	// Empty directory should have zero counts
	if root.FileCount != 0 {
		t.Errorf("Empty dir FileCount should be 0, got %d", root.FileCount)
	}
	if root.TotalSize != 0 {
		t.Errorf("Empty dir TotalSize should be 0, got %d", root.TotalSize)
	}
	if root.DirectFileCount != 0 {
		t.Errorf("Empty dir DirectFileCount should be 0, got %d", root.DirectFileCount)
	}
	if root.ExtensionStats != nil && len(root.ExtensionStats) > 0 {
		t.Errorf("Empty dir ExtensionStats should be nil or empty, got %v", root.ExtensionStats)
	}
}

func TestBuildTree_FilesWithoutExtension(t *testing.T) {
	dir := t.TempDir()
	
	// Create files without extensions
	files := []string{"Makefile", "Dockerfile", "README"}
	for _, name := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte("content"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	
	repo := &fakeSettingsRepo{}
	builder := New(repo, &domain.NoopLogger{})
	nodes, err := builder.BuildTree(dir, false, false)
	if err != nil {
		t.Fatalf("BuildTree error: %v", err)
	}
	
	root := nodes[0]
	
	// Files without extensions should not be counted in ExtensionStats
	if root.FileCount != 3 {
		t.Errorf("FileCount should be 3, got %d", root.FileCount)
	}
	if root.ExtensionStats != nil && len(root.ExtensionStats) > 0 {
		t.Errorf("ExtensionStats should be empty for files without extensions, got %v", root.ExtensionStats)
	}
}
