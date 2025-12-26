package tools

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"syntaxia/domain"
	"syntaxia/domain/analysis"
)

// === Mock implementations ===

// MockLogger implements domain.Logger for testing
type MockLogger struct {
	messages []string
}

func (m *MockLogger) Debug(message string)   { m.messages = append(m.messages, "DEBUG: "+message) }
func (m *MockLogger) Info(message string)    { m.messages = append(m.messages, "INFO: "+message) }
func (m *MockLogger) Warning(message string) { m.messages = append(m.messages, "WARN: "+message) }
func (m *MockLogger) Error(message string)   { m.messages = append(m.messages, "ERROR: "+message) }
func (m *MockLogger) Fatal(message string)   { m.messages = append(m.messages, "FATAL: "+message) }

// MockLanguageAnalyzer implements analysis.LanguageAnalyzer for testing
type MockLanguageAnalyzer struct {
	language       string
	extensions     []string
	symbols        []analysis.Symbol
	imports        []analysis.Import
	exports        []analysis.Export
	extractErr     error
	importsErr     error
	exportsErr     error
	funcBody       string
	funcStartLine  int
	funcEndLine    int
	funcBodyErr    error
}

func (m *MockLanguageAnalyzer) Language() string                { return m.language }
func (m *MockLanguageAnalyzer) Extensions() []string            { return m.extensions }
func (m *MockLanguageAnalyzer) CanAnalyze(filePath string) bool { return true }

func (m *MockLanguageAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	if m.extractErr != nil {
		return nil, m.extractErr
	}
	return m.symbols, nil
}

func (m *MockLanguageAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	if m.importsErr != nil {
		return nil, m.importsErr
	}
	return m.imports, nil
}

func (m *MockLanguageAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	if m.exportsErr != nil {
		return nil, m.exportsErr
	}
	return m.exports, nil
}

func (m *MockLanguageAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	if m.funcBodyErr != nil {
		return "", 0, 0, m.funcBodyErr
	}
	return m.funcBody, m.funcStartLine, m.funcEndLine, nil
}

// MockAnalyzerRegistry implements analysis.AnalyzerRegistry for testing
type MockAnalyzerRegistry struct {
	analyzer           analysis.LanguageAnalyzer
	analyzerByLang     map[string]analysis.LanguageAnalyzer
	supportedLangs     []string
	supportedExts      []string
}

func NewMockAnalyzerRegistry() *MockAnalyzerRegistry {
	return &MockAnalyzerRegistry{
		analyzerByLang: make(map[string]analysis.LanguageAnalyzer),
	}
}

func (m *MockAnalyzerRegistry) Register(analyzer analysis.LanguageAnalyzer) {
	m.analyzerByLang[analyzer.Language()] = analyzer
}

func (m *MockAnalyzerRegistry) GetAnalyzer(filePath string) analysis.LanguageAnalyzer {
	return m.analyzer
}

func (m *MockAnalyzerRegistry) GetAnalyzerByLanguage(lang string) analysis.LanguageAnalyzer {
	return m.analyzerByLang[lang]
}

func (m *MockAnalyzerRegistry) SupportedLanguages() []string {
	return m.supportedLangs
}

func (m *MockAnalyzerRegistry) SupportedExtensions() []string {
	return m.supportedExts
}

// MockSymbolIndex implements analysis.SymbolIndex for testing
type MockSymbolIndex struct {
	indexed       bool
	indexErr      error
	symbols       []analysis.Symbol
	symbolsByKind map[analysis.SymbolKind][]analysis.Symbol
	definition    *analysis.Symbol
	stats         map[string]int
}

func NewMockSymbolIndex() *MockSymbolIndex {
	return &MockSymbolIndex{
		symbolsByKind: make(map[analysis.SymbolKind][]analysis.Symbol),
		stats:         map[string]int{"total_symbols": 0, "files": 0},
	}
}

func (m *MockSymbolIndex) IndexProject(ctx context.Context, projectRoot string) error {
	if m.indexErr != nil {
		return m.indexErr
	}
	m.indexed = true
	return nil
}

func (m *MockSymbolIndex) IndexFile(ctx context.Context, filePath string, content []byte) error {
	return nil
}

func (m *MockSymbolIndex) SearchByName(query string) []analysis.Symbol {
	var results []analysis.Symbol
	for _, s := range m.symbols {
		if strings.Contains(strings.ToLower(s.Name), strings.ToLower(query)) {
			results = append(results, s)
		}
	}
	return results
}

func (m *MockSymbolIndex) FindByExactName(name string) []analysis.Symbol {
	var results []analysis.Symbol
	for _, s := range m.symbols {
		if s.Name == name {
			results = append(results, s)
		}
	}
	return results
}

func (m *MockSymbolIndex) GetSymbolsInFile(filePath string) []analysis.Symbol {
	var results []analysis.Symbol
	for _, s := range m.symbols {
		if s.FilePath == filePath {
			results = append(results, s)
		}
	}
	return results
}

func (m *MockSymbolIndex) GetSymbolsByKind(kind analysis.SymbolKind) []analysis.Symbol {
	return m.symbolsByKind[kind]
}

func (m *MockSymbolIndex) FindDefinition(name string, kind analysis.SymbolKind) *analysis.Symbol {
	return m.definition
}

func (m *MockSymbolIndex) Stats() map[string]int {
	return m.stats
}

func (m *MockSymbolIndex) IsIndexed() bool {
	return m.indexed
}

func (m *MockSymbolIndex) Clear() {
	m.indexed = false
	m.symbols = nil
}

// MockReferenceFinder implements domain.ReferenceFinder for testing
type MockReferenceFinder struct {
	references []domain.SymbolReference
	findErr    error
}

func (m *MockReferenceFinder) FindReferences(ctx context.Context, projectRoot string, symbolName string, symbolKind string) ([]domain.SymbolReference, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.references, nil
}

func (m *MockReferenceFinder) FindUsages(ctx context.Context, projectRoot string, symbolName string) ([]domain.SymbolReference, error) {
	var usages []domain.SymbolReference
	for _, ref := range m.references {
		if !ref.IsDefinition {
			usages = append(usages, ref)
		}
	}
	return usages, nil
}

// === Test helper functions ===

func createTestHandler(registry *MockAnalyzerRegistry, index *MockSymbolIndex, refFinder *MockReferenceFinder) *SymbolToolsHandler {
	logger := &MockLogger{}
	if registry == nil {
		registry = NewMockAnalyzerRegistry()
	}
	if index == nil {
		index = NewMockSymbolIndex()
	}
	return NewSymbolToolsHandler(registry, index, logger, refFinder)
}

// === ListSymbols Tests ===

func TestListSymbols_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main\nfunc main() {}"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language:   "go",
		extensions: []string{".go"},
		symbols: []analysis.Symbol{
			{Name: "main", Kind: analysis.KindFunction, StartLine: 2},
			{Name: "helper", Kind: analysis.KindFunction, StartLine: 5},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "main") || !strings.Contains(resultStr, "helper") {
		t.Errorf("expected symbols in result, got: %s", resultStr)
	}
}

func TestListSymbols_MissingPath(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.ListSymbols(map[string]any{}, "/tmp")

	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "path is required") {
		t.Errorf("expected 'path is required' error, got: %v", err)
	}
}

func TestListSymbols_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	handler := createTestHandler(nil, nil, nil)

	_, err := handler.ListSymbols(map[string]any{"path": "nonexistent.go"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestListSymbols_NoAnalyzer(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.xyz")
	os.WriteFile(testFile, []byte("content"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = nil

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{"path": "test.xyz"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No analyzer available") {
		t.Errorf("expected 'No analyzer available' message, got: %s", resultStr)
	}
}

func TestListSymbols_WithKindFilter(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		symbols: []analysis.Symbol{
			{Name: "MyClass", Kind: analysis.KindClass, StartLine: 1},
			{Name: "myFunc", Kind: analysis.KindFunction, StartLine: 5},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{
		"path": "test.go",
		"kind": "function",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "myFunc") {
		t.Errorf("expected myFunc in result, got: %s", resultStr)
	}
	if strings.Contains(resultStr, "MyClass") {
		t.Errorf("should not contain MyClass when filtering by function, got: %s", resultStr)
	}
}

func TestListSymbols_NoSymbolsFound(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		symbols:  []analysis.Symbol{},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{"path": "empty.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No symbols found") {
		t.Errorf("expected 'No symbols found' message, got: %s", resultStr)
	}
}

func TestListSymbols_ExtractError(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language:   "go",
		extractErr: errors.New("parse error"),
	}

	handler := createTestHandler(registry, nil, nil)
	_, err := handler.ListSymbols(map[string]any{"path": "test.go"}, tmpDir)

	if err == nil {
		t.Fatal("expected error from analyzer")
	}
	if !strings.Contains(err.Error(), "failed to extract symbols") {
		t.Errorf("expected 'failed to extract symbols' error, got: %v", err)
	}
}

// === SearchSymbols Tests ===

func TestSearchSymbols_Success(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{
		{Name: "UserService", Kind: analysis.KindClass, FilePath: "user.go", StartLine: 10},
		{Name: "UserRepository", Kind: analysis.KindInterface, FilePath: "repo.go", StartLine: 5},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "User"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "UserService") || !strings.Contains(resultStr, "UserRepository") {
		t.Errorf("expected User symbols in result, got: %s", resultStr)
	}
}

func TestSearchSymbols_MissingQuery(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.SearchSymbols(map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing query")
	}
	if !strings.Contains(err.Error(), "query is required") {
		t.Errorf("expected 'query is required' error, got: %v", err)
	}
}

func TestSearchSymbols_BuildsIndexIfNeeded(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = false
	index.symbols = []analysis.Symbol{
		{Name: "TestFunc", Kind: analysis.KindFunction, FilePath: "test.go"},
	}
	index.stats = map[string]int{"total_symbols": 1, "files": 1}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "Test"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !index.indexed {
		t.Error("expected index to be built")
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "TestFunc") {
		t.Errorf("expected TestFunc in result, got: %s", resultStr)
	}
}

func TestSearchSymbols_IndexError(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = false
	index.indexErr = errors.New("index failed")

	handler := createTestHandler(nil, index, nil)
	_, err := handler.SearchSymbols(map[string]any{"query": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error from index")
	}
	if !strings.Contains(err.Error(), "failed to build index") {
		t.Errorf("expected 'failed to build index' error, got: %v", err)
	}
}

func TestSearchSymbols_NoResults(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "nonexistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No symbols found") {
		t.Errorf("expected 'No symbols found' message, got: %s", resultStr)
	}
}

func TestSearchSymbols_WithKindFilter(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{
		{Name: "UserClass", Kind: analysis.KindClass, FilePath: "user.go"},
		{Name: "UserFunc", Kind: analysis.KindFunction, FilePath: "user.go"},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{
		"query": "User",
		"kind":  "class",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "UserClass") {
		t.Errorf("expected UserClass in result, got: %s", resultStr)
	}
	if strings.Contains(resultStr, "UserFunc") {
		t.Errorf("should not contain UserFunc when filtering by class, got: %s", resultStr)
	}
}

// === FindDefinition Tests ===

func TestFindDefinition_Success(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\n\nfunc MyFunction() {\n\treturn\n}\n"
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyFunction",
		Kind:      analysis.KindFunction,
		FilePath:  "main.go",
		StartLine: 3,
		EndLine:   5,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "MyFunction"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MyFunction") || !strings.Contains(resultStr, "main.go") {
		t.Errorf("expected definition info, got: %s", resultStr)
	}
}

func TestFindDefinition_MissingName(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.FindDefinition(map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected 'name is required' error, got: %v", err)
	}
}

func TestFindDefinition_NotFound(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = nil

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "NonExistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No definition found") {
		t.Errorf("expected 'No definition found' message, got: %s", resultStr)
	}
}

func TestFindDefinition_WithKindFilter(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main\ntype MyType struct{}"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyType",
		Kind:      analysis.KindType,
		FilePath:  "test.go",
		StartLine: 2,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{
		"name": "MyType",
		"kind": "type",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MyType") {
		t.Errorf("expected MyType in result, got: %s", resultStr)
	}
}

// === FindReferences Tests ===

func TestFindReferences_Success(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{
			{FilePath: "main.go", Line: 10, Column: 5, LineText: "func main() { MyFunc() }", IsDefinition: false},
			{FilePath: "utils.go", Line: 5, Column: 1, LineText: "func MyFunc() {}", IsDefinition: true},
			{FilePath: "test.go", Line: 15, Column: 8, LineText: "result := MyFunc()", IsDefinition: false},
		},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.FindReferences(map[string]any{"name": "MyFunc"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "3 references") {
		t.Errorf("expected 3 references, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, "main.go") || !strings.Contains(resultStr, "utils.go") {
		t.Errorf("expected file paths in result, got: %s", resultStr)
	}
}

func TestFindReferences_MissingName(t *testing.T) {
	handler := createTestHandler(nil, nil, &MockReferenceFinder{})
	_, err := handler.FindReferences(map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected 'name is required' error, got: %v", err)
	}
}

func TestFindReferences_NoReferenceFinder(t *testing.T) {
	// Create handler with explicitly nil referenceFinder (not a nil pointer in interface)
	logger := &MockLogger{}
	registry := NewMockAnalyzerRegistry()
	index := NewMockSymbolIndex()
	handler := NewSymbolToolsHandler(registry, index, logger, nil)

	_, err := handler.FindReferences(map[string]any{"name": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error when reference finder not initialized")
	}
	if !strings.Contains(err.Error(), "reference finder not initialized") {
		t.Errorf("expected 'reference finder not initialized' error, got: %v", err)
	}
}

func TestFindReferences_FinderError(t *testing.T) {
	refFinder := &MockReferenceFinder{
		findErr: errors.New("search failed"),
	}

	handler := createTestHandler(nil, nil, refFinder)
	_, err := handler.FindReferences(map[string]any{"name": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error from finder")
	}
	if !strings.Contains(err.Error(), "failed to find references") {
		t.Errorf("expected 'failed to find references' error, got: %v", err)
	}
}

func TestFindReferences_NoResults(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.FindReferences(map[string]any{"name": "NonExistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No references found") {
		t.Errorf("expected 'No references found' message, got: %s", resultStr)
	}
}

func TestFindReferences_ExcludeDefinition(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{
			{FilePath: "def.go", Line: 1, IsDefinition: true},
			{FilePath: "use.go", Line: 10, IsDefinition: false},
		},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.FindReferences(map[string]any{
		"name":               "MyFunc",
		"include_definition": false,
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "1 references") {
		t.Errorf("expected 1 reference (excluding definition), got: %s", resultStr)
	}
}

// === GetSymbolInfo Tests ===

func TestGetSymbolInfo_Success(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\n\n// MyFunc does something\nfunc MyFunc(x int) string {\n\treturn \"\"\n}\n"
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:       "MyFunc",
		Kind:       analysis.KindFunction,
		FilePath:   "main.go",
		StartLine:  4,
		EndLine:    6,
		Signature:  "func MyFunc(x int) string",
		DocComment: "MyFunc does something",
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "MyFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MyFunc") || !strings.Contains(resultStr, "function") {
		t.Errorf("expected symbol info, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_MissingName(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.GetSymbolInfo(map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing name")
	}
	if !strings.Contains(err.Error(), "name is required") {
		t.Errorf("expected 'name is required' error, got: %v", err)
	}
}

func TestGetSymbolInfo_NotFound(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = nil
	index.symbols = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "NonExistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No symbol found") {
		t.Errorf("expected 'No symbol found' message, got: %s", resultStr)
	}
}

// === GetClassHierarchy Tests ===

func TestGetClassHierarchy_Success(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "UserService",
		Kind:      analysis.KindClass,
		FilePath:  "service.go",
		StartLine: 10,
		Parent:    "BaseService",
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{
		{Name: "AdminService", Kind: analysis.KindClass, Parent: "UserService", FilePath: "admin.go", StartLine: 5},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{"class_name": "UserService"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "UserService") {
		t.Errorf("expected UserService in result, got: %s", resultStr)
	}
	if !strings.Contains(resultStr, "BaseService") {
		t.Errorf("expected parent BaseService in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_MissingClassName(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.GetClassHierarchy(map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for missing class_name")
	}
	if !strings.Contains(err.Error(), "class_name is required") {
		t.Errorf("expected 'class_name is required' error, got: %v", err)
	}
}

func TestGetClassHierarchy_NotFound(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = nil

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{"class_name": "NonExistent"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "not found") {
		t.Errorf("expected 'not found' message, got: %s", resultStr)
	}
}

// === GetImports Tests ===

func TestGetImports_Success(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(testFile, []byte("package main\nimport \"fmt\""), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports: []analysis.Import{
			{Path: "fmt", IsLocal: false},
			{Path: "./utils", IsLocal: true, Alias: "u"},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.GetImports(map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "fmt") || !strings.Contains(resultStr, "./utils") {
		t.Errorf("expected imports in result, got: %s", resultStr)
	}
}

func TestGetImports_MissingPath(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.GetImports(map[string]any{}, "/tmp")

	if err == nil {
		t.Fatal("expected error for missing path")
	}
	if !strings.Contains(err.Error(), "path is required") {
		t.Errorf("expected 'path is required' error, got: %v", err)
	}
}

func TestGetImports_NoImports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "simple.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports:  []analysis.Import{},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.GetImports(map[string]any{"path": "simple.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "No imports found") {
		t.Errorf("expected 'No imports found' message, got: %s", resultStr)
	}
}

// === CanHandle Tests ===

func TestSymbolToolsHandler_CanHandle(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		want     bool
	}{
		{"list_symbols", "list_symbols", true},
		{"search_symbols", "search_symbols", true},
		{"find_definition", "find_definition", true},
		{"find_references", "find_references", true},
		{"get_symbol_info", "get_symbol_info", true},
		{"get_class_hierarchy", "get_class_hierarchy", true},
		{"get_imports", "get_imports", true},
		{"unknown_tool", "unknown_tool", false},
		{"read_file", "read_file", false},
		{"git_status", "git_status", false},
	}

	handler := createTestHandler(nil, nil, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handler.CanHandle(tt.toolName); got != tt.want {
				t.Errorf("CanHandle(%q) = %v, want %v", tt.toolName, got, tt.want)
			}
		})
	}
}

// === GetTools Tests ===

func TestSymbolToolsHandler_GetTools(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	tools := handler.GetTools()

	if len(tools) == 0 {
		t.Fatal("expected non-empty tools list")
	}

	expectedTools := []string{
		"list_symbols",
		"search_symbols",
		"find_definition",
		"find_references",
		"get_symbol_info",
		"get_class_hierarchy",
		"get_imports",
	}

	toolNames := make(map[string]bool)
	for _, tool := range tools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("expected tool %q not found", expected)
		}
	}
}

// === Execute Tests ===

func TestSymbolToolsHandler_Execute_UnknownTool(t *testing.T) {
	handler := createTestHandler(nil, nil, nil)
	_, err := handler.Execute("unknown_tool", map[string]any{}, "/project")

	if err == nil {
		t.Fatal("expected error for unknown tool")
	}
	if !strings.Contains(err.Error(), "unknown symbol tool") {
		t.Errorf("expected 'unknown symbol tool' error, got: %v", err)
	}
}

func TestSymbolToolsHandler_Execute_ListSymbols(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main\nfunc test() {}"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		symbols:  []analysis.Symbol{{Name: "test", Kind: analysis.KindFunction}},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.Execute("list_symbols", map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "test") {
		t.Errorf("expected 'test' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_SearchSymbols(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{{Name: "TestFunc", Kind: analysis.KindFunction}}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.Execute("search_symbols", map[string]any{"query": "Test"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "TestFunc") {
		t.Errorf("expected 'TestFunc' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_FindDefinition(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc MyFunc() {}"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name: "MyFunc", Kind: analysis.KindFunction, FilePath: "main.go", StartLine: 2,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.Execute("find_definition", map[string]any{"name": "MyFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "MyFunc") {
		t.Errorf("expected 'MyFunc' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_FindReferences(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{
			{FilePath: "main.go", Line: 10, LineText: "MyFunc()"},
		},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.Execute("find_references", map[string]any{"name": "MyFunc"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "main.go") {
		t.Errorf("expected 'main.go' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_GetSymbolInfo(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main\nfunc Info() {}"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name: "Info", Kind: analysis.KindFunction, FilePath: "test.go", StartLine: 2,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.Execute("get_symbol_info", map[string]any{"name": "Info"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "Info") {
		t.Errorf("expected 'Info' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_GetClassHierarchy(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name: "MyClass", Kind: analysis.KindClass, FilePath: "class.go", StartLine: 1,
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.Execute("get_class_hierarchy", map[string]any{"class_name": "MyClass"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "MyClass") {
		t.Errorf("expected 'MyClass' in result, got: %s", result)
	}
}

func TestSymbolToolsHandler_Execute_GetImports(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nimport \"fmt\""), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports:  []analysis.Import{{Path: "fmt", IsLocal: false}},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.Execute("get_imports", map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result, "fmt") {
		t.Errorf("expected 'fmt' in result, got: %s", result)
	}
}


// === Additional Edge Case Tests ===

func TestListSymbols_WithParentAndSignature(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main\ntype Service struct{}\nfunc (s *Service) Method() {}"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		symbols: []analysis.Symbol{
			{Name: "Method", Kind: analysis.KindMethod, StartLine: 3, Parent: "Service", Signature: "func (s *Service) Method()"},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Service") || !strings.Contains(resultStr, "Method") {
		t.Errorf("expected parent and method in result, got: %s", resultStr)
	}
}

func TestSearchSymbols_LimitResults(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	// Add more than 30 symbols
	for i := 0; i < 50; i++ {
		index.symbols = append(index.symbols, analysis.Symbol{
			Name:     fmt.Sprintf("TestFunc%d", i),
			Kind:     analysis.KindFunction,
			FilePath: "test.go",
		})
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "TestFunc"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	// Should be limited to 30 results
	if strings.Contains(resultStr, "TestFunc40") {
		t.Errorf("should limit results to 30, got: %s", resultStr)
	}
}

func TestFindDefinition_FileReadError(t *testing.T) {
	tmpDir := t.TempDir()
	// Don't create the file - it should handle missing file gracefully

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MissingFunc",
		Kind:      analysis.KindFunction,
		FilePath:  "nonexistent.go",
		StartLine: 1,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "MissingFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "couldn't read file") {
		t.Errorf("expected file read error message, got: %s", resultStr)
	}
}

func TestFindDefinition_IndexError(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = false
	index.indexErr = errors.New("index build failed")

	handler := createTestHandler(nil, index, nil)
	_, err := handler.FindDefinition(map[string]any{"name": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error from index")
	}
	if !strings.Contains(err.Error(), "failed to build index") {
		t.Errorf("expected 'failed to build index' error, got: %v", err)
	}
}

func TestGetSymbolInfo_WithChildren(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "class.go"), []byte("package main\ntype MyClass struct{}"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyClass",
		Kind:      analysis.KindClass,
		FilePath:  "class.go",
		StartLine: 2,
		Children: []analysis.Symbol{
			{Name: "Method1", Kind: analysis.KindMethod, Signature: "func Method1()"},
			{Name: "Field1", Kind: analysis.KindField},
		},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "MyClass"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Members") || !strings.Contains(resultStr, "Method1") {
		t.Errorf("expected children in result, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_WithModifiers(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "PublicFunc",
		Kind:      analysis.KindFunction,
		FilePath:  "test.go",
		StartLine: 1,
		Modifiers: []string{"public", "static"},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "PublicFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Modifiers") {
		t.Errorf("expected modifiers in result, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_WithDocComment(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main\n// DocFunc does something\nfunc DocFunc() {}"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:       "DocFunc",
		Kind:       analysis.KindFunction,
		FilePath:   "test.go",
		StartLine:  3,
		DocComment: "DocFunc does something",
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "DocFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Documentation") {
		t.Errorf("expected documentation in result, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_FallbackToSearch(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = nil // No exact definition
	index.symbols = []analysis.Symbol{
		{Name: "SearchedFunc", Kind: analysis.KindFunction, FilePath: "test.go", StartLine: 1},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "SearchedFunc"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "SearchedFunc") {
		t.Errorf("expected to find via search, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_WithFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "specific.go"), []byte("package main"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = nil
	index.symbols = []analysis.Symbol{
		{Name: "Func", Kind: analysis.KindFunction, FilePath: "other.go", StartLine: 1},
		{Name: "Func", Kind: analysis.KindFunction, FilePath: "specific.go", StartLine: 5},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{
		"name":      "Func",
		"file_path": "specific.go",
	}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "specific.go") {
		t.Errorf("expected specific.go in result, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_IndexError(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = false
	index.indexErr = errors.New("index error")

	handler := createTestHandler(nil, index, nil)
	_, err := handler.GetSymbolInfo(map[string]any{"name": "test"}, "/project")

	if err == nil {
		t.Fatal("expected error from index")
	}
}

func TestGetClassHierarchy_Interface(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	// First lookup for class returns nil, then interface lookup succeeds
	index.definition = &analysis.Symbol{
		Name:      "MyInterface",
		Kind:      analysis.KindInterface,
		FilePath:  "interface.go",
		StartLine: 1,
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{"class_name": "MyInterface"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MyInterface") {
		t.Errorf("expected interface in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_DirectionUp(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "ChildClass",
		Kind:      analysis.KindClass,
		FilePath:  "child.go",
		StartLine: 1,
		Parent:    "ParentClass, Interface1",
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "ChildClass",
		"direction":  "up",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "ParentClass") {
		t.Errorf("expected parent in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_DirectionDown(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "BaseClass",
		Kind:      analysis.KindClass,
		FilePath:  "base.go",
		StartLine: 1,
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{
		{Name: "DerivedClass", Kind: analysis.KindClass, Parent: "BaseClass", FilePath: "derived.go", StartLine: 1},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "BaseClass",
		"direction":  "down",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "DerivedClass") {
		t.Errorf("expected subclass in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_IndexError(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = false
	index.indexErr = errors.New("index error")

	handler := createTestHandler(nil, index, nil)
	_, err := handler.GetClassHierarchy(map[string]any{"class_name": "Test"}, "/project")

	if err == nil {
		t.Fatal("expected error from index")
	}
}

func TestGetImports_Error(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language:   "go",
		importsErr: errors.New("parse error"),
	}

	handler := createTestHandler(registry, nil, nil)
	_, err := handler.GetImports(map[string]any{"path": "test.go"}, tmpDir)

	if err == nil {
		t.Fatal("expected error from analyzer")
	}
	if !strings.Contains(err.Error(), "failed to extract imports") {
		t.Errorf("expected 'failed to extract imports' error, got: %v", err)
	}
}

func TestGetImports_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	handler := createTestHandler(nil, nil, nil)

	_, err := handler.GetImports(map[string]any{"path": "nonexistent.go"}, tmpDir)

	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestGetImports_WithAlias(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(testFile, []byte("package main\nimport alias \"fmt\""), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports: []analysis.Import{
			{Path: "fmt", IsLocal: false, Alias: "alias"},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.GetImports(map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "alias") {
		t.Errorf("expected alias in result, got: %s", resultStr)
	}
}

func TestFindReferences_WithKindFilter(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{
			{FilePath: "main.go", Line: 10, LineText: "MyFunc()"},
		},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.FindReferences(map[string]any{
		"name": "MyFunc",
		"kind": "function",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "main.go") {
		t.Errorf("expected references, got: %s", resultStr)
	}
}


func TestGetClassHierarchy_WithParentDefinition(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	
	// Setup: ChildClass extends ParentClass, and ParentClass is also in the index
	parentSym := &analysis.Symbol{
		Name:      "ParentClass",
		Kind:      analysis.KindClass,
		FilePath:  "parent.go",
		StartLine: 1,
	}
	
	index.definition = &analysis.Symbol{
		Name:      "ChildClass",
		Kind:      analysis.KindClass,
		FilePath:  "child.go",
		StartLine: 1,
		Parent:    "ParentClass",
	}
	
	// Add parent to symbols so FindDefinition can find it
	index.symbols = []analysis.Symbol{*parentSym}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "ChildClass",
		"direction":  "both",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "ParentClass") {
		t.Errorf("expected parent class in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_NoSubclasses(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "LeafClass",
		Kind:      analysis.KindClass,
		FilePath:  "leaf.go",
		StartLine: 1,
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{} // No subclasses

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "LeafClass",
		"direction":  "down",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "none found") {
		t.Errorf("expected 'none found' for subclasses, got: %s", resultStr)
	}
}

func TestFindDefinition_WithEndLine(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\n\nfunc MyFunction() {\n\treturn\n}\n\nfunc Other() {}"
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyFunction",
		Kind:      analysis.KindFunction,
		FilePath:  "main.go",
		StartLine: 3,
		EndLine:   5,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "MyFunction"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MyFunction") {
		t.Errorf("expected function definition, got: %s", resultStr)
	}
}

func TestFindDefinition_StartLineNegative(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\nfunc Test() {}"
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "Test",
		Kind:      analysis.KindFunction,
		FilePath:  "main.go",
		StartLine: -1, // Invalid start line
		EndLine:   0,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "Test"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should handle gracefully
	resultStr := result.(string)
	if resultStr == "" {
		t.Error("expected non-empty result")
	}
}

func TestGetSymbolInfo_NoSourceFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Don't create the file

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MissingSymbol",
		Kind:      analysis.KindFunction,
		FilePath:  "missing.go",
		StartLine: 1,
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "MissingSymbol"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should still return symbol info without source
	resultStr := result.(string)
	if !strings.Contains(resultStr, "MissingSymbol") {
		t.Errorf("expected symbol info without source, got: %s", resultStr)
	}
}

func TestGetImports_OnlyLocalImports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports: []analysis.Import{
			{Path: "./utils", IsLocal: true},
			{Path: "../common", IsLocal: true},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.GetImports(map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Local") {
		t.Errorf("expected Local section, got: %s", resultStr)
	}
}

func TestGetImports_OnlyExternalImports(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "main.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		imports: []analysis.Import{
			{Path: "fmt", IsLocal: false},
			{Path: "os", IsLocal: false},
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.GetImports(map[string]any{"path": "main.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "External") {
		t.Errorf("expected External section, got: %s", resultStr)
	}
}


func TestFindReferences_CountsDefinitionsAndUsages(t *testing.T) {
	refFinder := &MockReferenceFinder{
		references: []domain.SymbolReference{
			{FilePath: "def.go", Line: 1, LineText: "func MyFunc() {}", IsDefinition: true},
			{FilePath: "use1.go", Line: 10, LineText: "MyFunc()", IsDefinition: false},
			{FilePath: "use2.go", Line: 20, LineText: "result := MyFunc()", IsDefinition: false},
		},
	}

	handler := createTestHandler(nil, nil, refFinder)
	result, err := handler.FindReferences(map[string]any{"name": "MyFunc"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	// Should show counts: 1 definition, 2 usages
	if !strings.Contains(resultStr, "1 definitions") || !strings.Contains(resultStr, "2 usages") {
		t.Errorf("expected definition and usage counts, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_WithSignature(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main\nfunc Test(x int) string { return \"\" }"), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "Test",
		Kind:      analysis.KindFunction,
		FilePath:  "test.go",
		StartLine: 2,
		Signature: "func Test(x int) string",
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "Test"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Signature") {
		t.Errorf("expected signature in result, got: %s", resultStr)
	}
}

func TestListSymbols_WithStartLineZero(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	registry := NewMockAnalyzerRegistry()
	registry.analyzer = &MockLanguageAnalyzer{
		language: "go",
		symbols: []analysis.Symbol{
			{Name: "NoLine", Kind: analysis.KindFunction, StartLine: 0}, // No line info
		},
	}

	handler := createTestHandler(registry, nil, nil)
	result, err := handler.ListSymbols(map[string]any{"path": "test.go"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "NoLine") {
		t.Errorf("expected symbol without line info, got: %s", resultStr)
	}
}

func TestSearchSymbols_WithFilePathAndStartLine(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{
		{Name: "TestFunc", Kind: analysis.KindFunction, FilePath: "test.go", StartLine: 10},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "Test"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "test.go") || !strings.Contains(resultStr, ":10") {
		t.Errorf("expected file path and line number, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_MultipleParents(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyClass",
		Kind:      analysis.KindClass,
		FilePath:  "class.go",
		StartLine: 1,
		Parent:    "BaseClass, Interface1, Interface2",
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "MyClass",
		"direction":  "up",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "BaseClass") || !strings.Contains(resultStr, "Interface1") {
		t.Errorf("expected multiple parents, got: %s", resultStr)
	}
}


func TestGetClassHierarchy_EmptyParent(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "StandaloneClass",
		Kind:      analysis.KindClass,
		FilePath:  "standalone.go",
		StartLine: 1,
		Parent:    "", // No parent
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "StandaloneClass",
		"direction":  "both",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "StandaloneClass") {
		t.Errorf("expected class name in result, got: %s", resultStr)
	}
}

func TestGetClassHierarchy_ParentWithSpaces(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "MyClass",
		Kind:      analysis.KindClass,
		FilePath:  "class.go",
		StartLine: 1,
		Parent:    "  BaseClass  ,  Interface1  ", // With extra spaces
	}
	index.symbolsByKind[analysis.KindClass] = []analysis.Symbol{}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetClassHierarchy(map[string]any{
		"class_name": "MyClass",
		"direction":  "up",
	}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "BaseClass") {
		t.Errorf("expected trimmed parent name, got: %s", resultStr)
	}
}

func TestFindDefinition_EndLineExceedsFileLength(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\nfunc Short() {}"
	os.WriteFile(filepath.Join(tmpDir, "short.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "Short",
		Kind:      analysis.KindFunction,
		FilePath:  "short.go",
		StartLine: 2,
		EndLine:   100, // Way beyond file length
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.FindDefinition(map[string]any{"name": "Short"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Short") {
		t.Errorf("expected function in result, got: %s", resultStr)
	}
}

func TestGetSymbolInfo_EndLineZero(t *testing.T) {
	tmpDir := t.TempDir()
	content := "package main\nfunc Test() {}\nfunc Other() {}"
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(content), 0644)

	index := NewMockSymbolIndex()
	index.indexed = true
	index.definition = &analysis.Symbol{
		Name:      "Test",
		Kind:      analysis.KindFunction,
		FilePath:  "test.go",
		StartLine: 2,
		EndLine:   0, // No end line specified
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.GetSymbolInfo(map[string]any{"name": "Test"}, tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "Test") {
		t.Errorf("expected symbol info, got: %s", resultStr)
	}
}

func TestSearchSymbols_NoFilePath(t *testing.T) {
	index := NewMockSymbolIndex()
	index.indexed = true
	index.symbols = []analysis.Symbol{
		{Name: "OrphanFunc", Kind: analysis.KindFunction, FilePath: "", StartLine: 0},
	}

	handler := createTestHandler(nil, index, nil)
	result, err := handler.SearchSymbols(map[string]any{"query": "Orphan"}, "/project")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resultStr := result.(string)
	if !strings.Contains(resultStr, "OrphanFunc") {
		t.Errorf("expected symbol without file path, got: %s", resultStr)
	}
}
