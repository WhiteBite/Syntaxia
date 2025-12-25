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
