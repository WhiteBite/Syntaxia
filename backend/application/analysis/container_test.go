package analysis

import (
	"context"
	"syntaxia/domain"
	"testing"
)

func TestNewContainer(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	if container == nil {
		t.Fatal("expected container to be created")
	}

	if container.GetProjectRoot() != "" {
		t.Error("expected empty project root initially")
	}
}

func TestContainer_SetProject(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	container.SetProject("/test/project")

	if container.GetProjectRoot() != "/test/project" {
		t.Errorf("expected project root '/test/project', got '%s'", container.GetProjectRoot())
	}

	container.SetProject("/test/project")

	container.SetProject("/test/project2")

	if container.GetProjectRoot() != "/test/project2" {
		t.Errorf("expected project root '/test/project2', got '%s'", container.GetProjectRoot())
	}
}

func TestContainer_GetSymbolIndex(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	// Without factory, should return nil
	index1 := container.GetSymbolIndex()
	if index1 != nil {
		t.Error("expected nil symbol index without factory")
	}
}

func TestContainer_GetCallGraph(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	// Without factory, should return nil
	cg1 := container.GetCallGraph()
	if cg1 != nil {
		t.Error("expected nil call graph without factory")
	}
}

func TestContainer_GetGitContext(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	gc := container.GetGitContext()
	if gc != nil {
		t.Error("expected nil git context without project")
	}

	container.SetProject("/test/project")
	gc = container.GetGitContext()
	// Without factory, should still be nil
	if gc != nil {
		t.Error("expected nil git context without factory")
	}
}

func TestContainer_GetContextMemory(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	// Without factory, should return nil
	mem1 := container.GetContextMemory()
	if mem1 != nil {
		t.Error("expected nil context memory without factory")
	}
}

func TestContainer_InvalidateCache(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	container.SetProject("/test/project")

	container.InvalidateCache()

	stats := container.Stats()
	if stats["symbolIndexBuilt"].(bool) {
		t.Error("expected symbolIndexBuilt to be false after invalidation")
	}
	if stats["callGraphBuilt"].(bool) {
		t.Error("expected callGraphBuilt to be false after invalidation")
	}
}

func TestContainer_SetProjectInvalidatesCache(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	container.SetProject("/project1")
	container.SetProject("/project2")

	if container.GetProjectRoot() != "/project2" {
		t.Error("expected project root to be updated")
	}
}

func TestContainer_SemanticSearch(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	if container.GetSemanticSearch() != nil {
		t.Error("expected nil semantic search initially")
	}

	mockSS := &mockSemanticSearcher{}
	container.SetSemanticSearch(mockSS)

	if container.GetSemanticSearch() != mockSS {
		t.Error("expected semantic search to be set")
	}
}

func TestContainer_Stats(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	stats := container.Stats()

	if _, ok := stats["projectRoot"]; !ok {
		t.Error("expected projectRoot in stats")
	}
	if _, ok := stats["symbolIndexBuilt"]; !ok {
		t.Error("expected symbolIndexBuilt in stats")
	}
	if _, ok := stats["callGraphBuilt"]; !ok {
		t.Error("expected callGraphBuilt in stats")
	}
}

func TestContainer_ConcurrentAccess(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})
	container.SetProject("/test/project")

	done := make(chan bool, 10)

	for range 10 {
		go func() {
			_ = container.GetSymbolIndex()
			_ = container.GetCallGraph()
			_ = container.GetContextMemory()
			_ = container.Stats()
			done <- true
		}()
	}

	for range 10 {
		<-done
	}
}

func TestContainer_EnsureSymbolIndexBuilt_NoProject(t *testing.T) {
	logger := &domain.NoopLogger{}
	container := NewContainer(logger, ContainerConfig{})

	err := container.EnsureSymbolIndexBuilt(context.Background())
	if err != nil {
		t.Errorf("expected nil error without project, got %v", err)
	}
}

type mockSemanticSearcher struct{}

func (m *mockSemanticSearcher) Search(query string, limit int) ([]domain.SemanticSearchResult, error) {
	return nil, nil
}
