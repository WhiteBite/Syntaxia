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

func (m *mockSymbolGraphBuilder) UpdateGraph(ctx context.Context, projectRoot string, changedFiles []string) (*domain.SymbolGraph, error) {
	return m.graph, m.err
}

func (m *mockSymbolGraphBuilder) GetSuggestions(ctx context.Context, query string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return nil, nil
}

func (m *mockSymbolGraphBuilder) GetDependencies(ctx context.Context, symbolID string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return nil, nil
}

func (m *mockSymbolGraphBuilder) GetDependents(ctx context.Context, symbolID string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	return nil, nil
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


// TestTestEngine_RunTargetedTests_NoAnalyzer tests running targeted tests without analyzer
func TestTestEngine_RunTargetedTests_NoAnalyzer(t *testing.T) {
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
	// Note: No analyzer registered

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAffected,
	}

	affectedGraph := &domain.AffectedGraph{
		ChangedFiles:  []string{"file1.go"},
		AffectedFiles: []string{"file1.go"},
	}

	// Should fall back to running all tests
	testResults, err := engine.RunTargetedTests(context.Background(), config, affectedGraph)

	require.NoError(t, err)
	assert.NotEmpty(t, testResults)
}

// TestTestEngine_RunTargetedTests_WithSmokeTests tests running targeted tests with smoke scope
func TestTestEngine_RunTargetedTests_WithSmokeTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "smoke_test.go", Name: "smoke_test.go", Type: "smoke"},
		{Path: "unit_test.go", Name: "unit_test.go", Type: "unit"},
	}
	runner := &mockTestRunner{language: "go", tests: tests}
	analyzer := &mockTestAnalyzer{testsForFile: []string{"unit_test.go"}, isSmoke: true}

	engine.RegisterTestRunner("go", runner)
	engine.RegisterTestAnalyzer("go", analyzer)

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAffectedSmoke,
	}

	affectedGraph := &domain.AffectedGraph{
		ChangedFiles:  []string{"file1.go"},
		AffectedFiles: []string{"file1.go"},
	}

	results, err := engine.RunTargetedTests(context.Background(), config, affectedGraph)

	require.NoError(t, err)
	assert.NotEmpty(t, results)
}

// TestTestEngine_findTestsForFile tests finding tests for different file types
func TestTestEngine_findTestsForFile(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	analyzer := &mockTestAnalyzer{testsForFile: []string{"test.go"}}
	engine.RegisterTestAnalyzer("go", analyzer)

	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "go_file",
			filePath:    "main.go",
			expectError: false,
		},
		{
			name:        "typescript_file",
			filePath:    "main.ts",
			expectError: true, // No analyzer registered
		},
		{
			name:        "tsx_file",
			filePath:    "component.tsx",
			expectError: true,
		},
		{
			name:        "javascript_file",
			filePath:    "script.js",
			expectError: true,
		},
		{
			name:        "jsx_file",
			filePath:    "component.jsx",
			expectError: true,
		},
		{
			name:        "java_file",
			filePath:    "Main.java",
			expectError: true,
		},
		{
			name:        "python_file",
			filePath:    "main.py",
			expectError: true,
		},
		{
			name:        "unsupported_file",
			filePath:    "main.rb",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := engine.findTestsForFile(context.Background(), tt.filePath, "/project")
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTestEngine_filterTestsByScope_AffectedSmoke tests filtering with affected+smoke scope
func TestTestEngine_filterTestsByScope_AffectedSmoke(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "unit_test.go", Name: "unit_test.go", Type: "unit"},
		{Path: "smoke_test.go", Name: "smoke_test.go", Type: "smoke"},
		{Path: "integration_test.go", Name: "integration_test.go", Type: "integration"},
	}

	// AffectedSmoke should return all tests (filtering happens in RunTargetedTests)
	filtered := engine.filterTestsByScope(tests, domain.TestScopeAffectedSmoke)
	assert.Len(t, filtered, 3)
}

// TestTestEngine_findSmokeTests tests finding smoke tests
func TestTestEngine_findSmokeTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "smoke_test.go", Name: "smoke_test.go", Type: "smoke"},
		{Path: "unit_test.go", Name: "unit_test.go", Type: "unit"},
	}
	runner := &mockTestRunner{language: "go", tests: tests}
	analyzer := &mockTestAnalyzer{isSmoke: true}

	engine.RegisterTestRunner("go", runner)
	engine.RegisterTestAnalyzer("go", analyzer)

	smokeTests, err := engine.findSmokeTests(context.Background(), "/project", "go")

	require.NoError(t, err)
	assert.NotEmpty(t, smokeTests)
}

// TestTestEngine_findSmokeTests_NoAnalyzer tests finding smoke tests without analyzer
func TestTestEngine_findSmokeTests_NoAnalyzer(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	runner := &mockTestRunner{language: "go", tests: []*domain.TestInfo{}}
	engine.RegisterTestRunner("go", runner)
	// No analyzer registered

	_, err := engine.findSmokeTests(context.Background(), "/project", "go")

	assert.Error(t, err)
}

// TestTestEngine_findSmokeTests_NoRunner tests finding smoke tests without runner
func TestTestEngine_findSmokeTests_NoRunner(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	analyzer := &mockTestAnalyzer{isSmoke: true}
	engine.RegisterTestAnalyzer("go", analyzer)
	// No runner registered

	_, err := engine.findSmokeTests(context.Background(), "/project", "go")

	assert.Error(t, err)
}

// TestTestEngine_collectTargetedTests tests collecting targeted tests
func TestTestEngine_collectTargetedTests(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	analyzer := &mockTestAnalyzer{testsForFile: []string{"test1.go", "test2.go"}}

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
	}

	affectedGraph := &domain.AffectedGraph{
		AffectedFiles: []string{"file1.go", "file2.go"},
	}

	targetTests, testSet := engine.collectTargetedTests(context.Background(), config, affectedGraph, analyzer)

	assert.NotEmpty(t, targetTests)
	assert.NotEmpty(t, testSet)
}

// TestTestEngine_addSmokeTests_NotSmokeScope tests addSmokeTests with non-smoke scope
func TestTestEngine_addSmokeTests_NotSmokeScope(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeUnit, // Not smoke scope
	}

	targetTests := []string{"test1.go"}
	testSet := map[string]bool{"test1.go": true}

	result := engine.addSmokeTests(context.Background(), config, targetTests, testSet)

	// Should return unchanged list
	assert.Equal(t, targetTests, result)
}

// TestTestEngine_BuildAffectedGraph_WithDependencies tests building affected graph with dependencies
func TestTestEngine_BuildAffectedGraph_WithDependencies(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{
		graph: &domain.SymbolGraph{
			Nodes: []*domain.SymbolNode{
				{ID: "1", Path: "file1.go"},
				{ID: "2", Path: "file2.go"},
				{ID: "3", Path: "file3.go"},
				{ID: "4", Path: "file4.go"},
			},
			Edges: []*domain.SymbolEdge{
				{From: "1", To: "2"},
				{From: "2", To: "3"},
				{From: "3", To: "4"},
			},
		},
	}
	engine := NewTestEngine(log, symbolGraph)

	// Register analyzer for Go
	analyzer := &mockTestAnalyzer{testsForFile: []string{"test.go"}}
	engine.RegisterTestAnalyzer("go", analyzer)

	changedFiles := []string{"file1.go"}
	graph, err := engine.BuildAffectedGraph(context.Background(), changedFiles, "/project")

	require.NoError(t, err)
	assert.NotNil(t, graph)
	assert.Contains(t, graph.ChangedFiles, "file1.go")
	// Should include indirect dependencies
	assert.GreaterOrEqual(t, len(graph.AffectedFiles), 1)
}

// TestTestEngine_findFileDependencies_SelfReference tests finding dependencies with self-reference
func TestTestEngine_findFileDependencies_SelfReference(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{
			{ID: "1", Path: "file1.go"},
		},
		Edges: []*domain.SymbolEdge{
			{From: "1", To: "1"}, // Self-reference
		},
	}

	deps, err := engine.findFileDependencies("file1.go", graph)

	require.NoError(t, err)
	// Should not include self in dependencies
	assert.Empty(t, deps)
}

// TestTestEngine_findIndirectDependencies_CyclicDependencies tests handling cyclic dependencies
func TestTestEngine_findIndirectDependencies_CyclicDependencies(t *testing.T) {
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
			{From: "3", To: "1"}, // Cyclic dependency
		},
	}

	directDeps := []string{"file2.go"}
	indirectDeps := engine.findIndirectDependencies(directDeps, graph)

	// Should handle cyclic dependencies without infinite loop
	assert.NotNil(t, indirectDeps)
}

// TestTestEngine_RunTests_FilterByScope tests running tests with different scopes
func TestTestEngine_RunTests_FilterByScope(t *testing.T) {
	log := &domain.NoopLogger{}
	symbolGraph := &mockSymbolGraphBuilder{}
	engine := NewTestEngine(log, symbolGraph)

	tests := []*domain.TestInfo{
		{Path: "unit_test.go", Name: "unit_test.go", Type: "unit"},
		{Path: "smoke_test.go", Name: "smoke_test.go", Type: "smoke"},
		{Path: "integration_test.go", Name: "integration_test.go", Type: "integration"},
	}
	results := []*domain.TestResult{
		{TestPath: "unit_test.go", Success: true},
	}
	runner := &mockTestRunner{language: "go", tests: tests, results: results}
	engine.RegisterTestRunner("go", runner)

	testCases := []struct {
		name  string
		scope domain.TestScope
	}{
		{"all", domain.TestScopeAll},
		{"unit", domain.TestScopeUnit},
		{"smoke", domain.TestScopeSmoke},
		{"integration", domain.TestScopeIntegration},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &domain.TestConfig{
				Language:    "go",
				ProjectPath: "/project",
				Scope:       tc.scope,
			}

			testResults, err := engine.RunTests(context.Background(), config)

			require.NoError(t, err)
			assert.NotNil(t, testResults)
		})
	}
}

// TestTestEngine_RunTargetedTests_EmptyAffectedFiles tests with empty affected files
func TestTestEngine_RunTargetedTests_EmptyAffectedFiles(t *testing.T) {
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
	analyzer := &mockTestAnalyzer{testsForFile: []string{}} // No tests found

	engine.RegisterTestRunner("go", runner)
	engine.RegisterTestAnalyzer("go", analyzer)

	config := &domain.TestConfig{
		Language:    "go",
		ProjectPath: "/project",
		Scope:       domain.TestScopeAffected,
	}

	affectedGraph := &domain.AffectedGraph{
		ChangedFiles:  []string{},
		AffectedFiles: []string{},
	}

	// Should fall back to running all tests
	testResults, err := engine.RunTargetedTests(context.Background(), config, affectedGraph)

	require.NoError(t, err)
	assert.NotEmpty(t, testResults)
}
