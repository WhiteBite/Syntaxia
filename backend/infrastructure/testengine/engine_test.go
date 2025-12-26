package testengine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// mockSymbolGraphBuilder implements domain.SymbolGraphBuilder for testing
type mockSymbolGraphBuilder struct {
	graph *domain.SymbolGraph
	err   error
}

func (m *mockSymbolGraphBuilder) BuildGraph(ctx context.Context, projectPath string) (*domain.SymbolGraph, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.graph, nil
}

// mockTestRunner implements domain.TestRunner for testing
type mockTestRunner struct {
	language string
	tests    []*domain.TestInfo
	results  []*domain.TestResult
}

func (m *mockTestRunner) GetLanguage() string {
	return m.language
}

func (m *mockTestRunner) RunTest(ctx context.Context, testPath string, config *domain.TestConfig) (*domain.TestResult, error) {
	return &domain.TestResult{
		TestPath: testPath,
		TestName: testPath,
		Language: m.language,
		Success:  true,
	}, nil
}

func (m *mockTestRunner) RunTestSuite(ctx context.Context, suite *domain.TestSuite) ([]*domain.TestResult, error) {
	return m.results, nil
}

func (m *mockTestRunner) DiscoverTests(ctx context.Context, projectPath string) ([]*domain.TestInfo, error) {
	return m.tests, nil
}

// mockTestAnalyzer implements domain.TestAnalyzer for testing
type mockTestAnalyzer struct {
	dependencies []string
	testsForFile []string
	isSmoke      bool
}

func (m *mockTestAnalyzer) AnalyzeTestDependencies(ctx context.Context, testPath string) ([]string, error) {
	return m.dependencies, nil
}

func (m *mockTestAnalyzer) FindTestsForFile(ctx context.Context, filePath, projectPath string) ([]string, error) {
	return m.testsForFile, nil
}

func (m *mockTestAnalyzer) IsSmokeTest(ctx context.Context, testPath string) (bool, error) {
	return m.isSmoke, nil
}

// TestNewTestEngine tests engine creation
func TestNewTestEngine(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}

	engine := NewTestEngine(log, symbolGraph)

	assert.NotNil(t, engine)
	assert.NotNil(t, engine.testRunners)
	assert.NotNil(t, engine.testAnalyzers)
}

// TestTestEngine_RegisterTestRunner tests runner registration
func TestTestEngine_RegisterTestRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	runner := &mockTestRunner{language: "go"}
	engine.RegisterTestRunner("go", runner)

	assert.Len(t, engine.testRunners, 1)
}

// TestTestEngine_RegisterTestAnalyzer tests analyzer registration
func TestTestEngine_RegisterTestAnalyzer(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	analyzer := &mockTestAnalyzer{}
	engine.RegisterTestAnalyzer("go", analyzer)

	assert.Len(t, engine.testAnalyzers, 1)
}

// TestTestEngine_GetSupportedLanguages tests getting supported languages
func TestTestEngine_GetSupportedLanguages(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	engine.RegisterTestRunner("go", &mockTestRunner{language: "go"})
	engine.RegisterTestRunner("python", &mockTestRunner{language: "python"})

	languages := engine.GetSupportedLanguages()

	assert.Len(t, languages, 2)
}

// TestTestEngine_DiscoverTests tests test discovery
func TestTestEngine_DiscoverTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "test1.go", Name: "test1.go", Type: "unit"},
		{Path: "test2.go", Name: "test2.go", Type: "unit"},
	}
	runner := &mockTestRunner{language: "go", tests: tests}
	engine.RegisterTestRunner("go", runner)

	suite, err := engine.DiscoverTests(context.Background(), "/project", "go")

	require.NoError(t, err)
	assert.NotNil(t, suite)
	assert.Len(t, suite.Tests, 2)
	assert.Equal(t, "go", suite.Language)
}

// TestTestEngine_DiscoverTests_NoRunner tests discovery without runner
func TestTestEngine_DiscoverTests_NoRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	_, err := engine.DiscoverTests(context.Background(), "/project", "unknown")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no test runner registered")
}

// TestTestEngine_RunTests tests running tests
func TestTestEngine_RunTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "test1.go", Name: "test1.go", Type: "unit"},
	}
	results := []*domain.TestResult{
		{TestPath: "test1.go", Success: true},
	}
	runner := &mockTestRunner{language: "go", tests: tests, results: results}
	engine.RegisterTestRunner("go", runner)

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAll,
	}

	testResults, err := engine.RunTests(context.Background(), config)

	require.NoError(t, err)
	assert.Len(t, testResults, 1)
}

// TestTestEngine_RunTests_NoRunner tests running tests without runner
func TestTestEngine_RunTests_NoRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	config := &domain.TestConfig{
		Language:    "unknown",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAll,
	}

	_, err := engine.RunTests(context.Background(), config)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no test runner registered")
}

// TestTestEngine_GetTestCoverage tests getting test coverage
func TestTestEngine_GetTestCoverage(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	coverage, err := engine.GetTestCoverage(context.Background(), "test.go")

	require.NoError(t, err)
	assert.NotNil(t, coverage)
	assert.Equal(t, 0.0, coverage.Percentage)
}

// TestTestEngine_filterTestsByScope tests filtering tests by scope
func TestTestEngine_filterTestsByScope(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "unit_test.go", Name: "unit_test.go", Type: "unit"},
		{Path: "smoke_test.go", Name: "smoke_test.go", Type: "smoke"},
		{Path: "integration_test.go", Name: "integration_test.go", Type: "integration"},
	}

	testCases := []struct {
		name     string
		scope    domain.TestScope
		expected int
	}{
		{"all", domain.TestScopeAll, 3},
		{"unit", domain.TestScopeUnit, 1},
		{"smoke", domain.TestScopeSmoke, 1},
		{"integration", domain.TestScopeIntegration, 1},
		{"affected", domain.TestScopeAffected, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filtered := engine.filterTestsByScope(tests, tc.scope)
			assert.Len(t, filtered, tc.expected)
		})
	}
}

// TestTestEngine_BuildAffectedGraph tests building affected graph
func TestTestEngine_BuildAffectedGraph(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{
		graph: &domain.SymbolGraph{
			Nodes: []*domain.SymbolNode{
				{ID: "1", Path: "file1.go"},
				{ID: "2", Path: "file2.go"},
			},
			Edges: []*domain.SymbolEdge{
				{From: "1", To: "2"},
			},
		},
	}
	engine := NewTestEngine(log, symbolGraph)

	changedFiles := []string{"file1.go"}
	graph, err := engine.BuildAffectedGraph(context.Background(), changedFiles, "/project")

	require.NoError(t, err)
	assert.NotNil(t, graph)
	assert.Contains(t, graph.ChangedFiles, "file1.go")
}

// TestTestEngine_BuildAffectedGraph_SymbolGraphError tests building affected graph with error
func TestTestEngine_BuildAffectedGraph_SymbolGraphError(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{
		err: assert.AnError,
	}
	engine := NewTestEngine(log, symbolGraph)

	changedFiles := []string{"file1.go"}
	graph, err := engine.BuildAffectedGraph(context.Background(), changedFiles, "/project")

	require.NoError(t, err)
	assert.NotNil(t, graph)
	// Falls back to simple graph with only changed files
	assert.Equal(t, changedFiles, graph.AffectedFiles)
}

// TestTestEngine_findFileDependencies tests finding file dependencies
func TestTestEngine_findFileDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{
			{ID: "1", Path: "file1.go"},
			{ID: "2", Path: "file2.go"},
			{ID: "3", Path: "file3.go"},
		},
		Edges: []*domain.SymbolEdge{
			{From: "1", To: "2"},
			{From: "1", To: "3"},
		},
	}

	deps, err := engine.findFileDependencies("file1.go", graph)

	require.NoError(t, err)
	assert.Len(t, deps, 2)
}

// TestTestEngine_findFileDependencies_NoDeps tests finding dependencies with no deps
func TestTestEngine_findFileDependencies_NoDeps(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{
			{ID: "1", Path: "file1.go"},
		},
		Edges: []*domain.SymbolEdge{},
	}

	deps, err := engine.findFileDependencies("file1.go", graph)

	require.NoError(t, err)
	assert.Empty(t, deps)
}

// TestTestEngine_findIndirectDependencies tests finding indirect dependencies
func TestTestEngine_findIndirectDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{
			{ID: "1", Path: "file1.go"},
			{ID: "2", Path: "file2.go"},
			{ID: "3", Path: "file3.go"},
		},
		Edges: []*domain.SymbolEdge{
			{From: "1", To: "2"},
			{From: "2", To: "3"},
		},
	}

	directDeps := []string{"file2.go"}
	indirectDeps := engine.findIndirectDependencies(directDeps, graph)

	assert.Contains(t, indirectDeps, "file3.go")
}

// TestTestEngine_findIndirectDependencies_NilGraph tests with nil graph
func TestTestEngine_findIndirectDependencies_NilGraph(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	directDeps := []string{"file1.go"}
	indirectDeps := engine.findIndirectDependencies(directDeps, nil)

	assert.Empty(t, indirectDeps)
}

// TestTestEngine_RunTargetedTests tests running targeted tests
func TestTestEngine_RunTargetedTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "test1.go", Name: "test1.go", Type: "unit"},
	}
	runner := &mockTestRunner{language: "go", tests: tests}
	analyzer := &mockTestAnalyzer{testsForFile: []string{"test1.go"}}

	engine.RegisterTestRunner("go", runner)
	engine.RegisterTestAnalyzer("go", analyzer)

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAffected,
	}

	affectedGraph := &domain.AffectedGraph{
		ChangedFiles:  []string{"file1.go"},
		AffectedFiles: []string{"file1.go"},
	}

	results, err := engine.RunTargetedTests(context.Background(), config, affectedGraph)

	require.NoError(t, err)
	assert.NotEmpty(t, results)
}

// TestTestEngine_RunTargetedTests_NoRunner tests running targeted tests without runner
func TestTestEngine_RunTargetedTests_NoRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	config := &domain.TestConfig{
		Language:    "unknown",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAffected,
	}

	affectedGraph := &domain.AffectedGraph{
		ChangedFiles:  []string{"file1.go"},
		AffectedFiles: []string{"file1.go"},
	}

	_, err := engine.RunTargetedTests(context.Background(), config, affectedGraph)

	assert.Error(t, err)
}

// TestBuildNodeIndex tests building node index
func TestBuildNodeIndex(t *testing.T) {
	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{
			{ID: "1", Path: "file1.go"},
			{ID: "2", Path: "file2.go"},
		},
	}

	index := buildNodeIndex(graph)

	assert.Len(t, index, 2)
	assert.NotNil(t, index["1"])
	assert.NotNil(t, index["2"])
}

// TestBuildNodeIndex_EmptyGraph tests building node index with empty graph
func TestBuildNodeIndex_EmptyGraph(t *testing.T) {
	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{},
	}

	index := buildNodeIndex(graph)

	assert.Empty(t, index)
}
