package analyzers

import (
	"path/filepath"
	"syntaxia/domain/analysis"
	"strings"
	"sync"
)

// AnalyzerRegistryImpl implements analysis.AnalyzerRegistry
type AnalyzerRegistryImpl struct {
	mu        sync.RWMutex
	analyzers []analysis.LanguageAnalyzer
	byLang    map[string]analysis.LanguageAnalyzer
	byExt     map[string]analysis.LanguageAnalyzer
}

// NewAnalyzerRegistry creates a new analyzer registry with default analyzers
func NewAnalyzerRegistry() *AnalyzerRegistryImpl {
	registry := &AnalyzerRegistryImpl{
		analyzers: make([]analysis.LanguageAnalyzer, 0),
		byLang:    make(map[string]analysis.LanguageAnalyzer),
		byExt:     make(map[string]analysis.LanguageAnalyzer),
	}

	// Register default analyzers (using existing implementations)
	registry.Register(NewGoAnalyzer())
	registry.Register(NewTypeScriptAnalyzer())
	registry.Register(NewJavaScriptAnalyzer())
	registry.Register(NewJavaAnalyzer())
	registry.Register(NewKotlinAnalyzer())
	// Phase 2: Vue and Dart analyzers
	registry.Register(NewVueAnalyzer())
	registry.Register(NewDartAnalyzer())
	// New language analyzers
	registry.Register(NewPythonAnalyzer())
	registry.Register(NewRustAnalyzer())
	registry.Register(NewCSharpAnalyzer())

	return registry
}

// Register adds an analyzer to the registry
func (r *AnalyzerRegistryImpl) Register(analyzer analysis.LanguageAnalyzer) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.analyzers = append(r.analyzers, analyzer)
	r.byLang[analyzer.Language()] = analyzer

	for _, ext := range analyzer.Extensions() {
		r.byExt[ext] = analyzer
	}
}

// GetAnalyzer returns analyzer for file path
func (r *AnalyzerRegistryImpl) GetAnalyzer(filePath string) analysis.LanguageAnalyzer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ext := strings.ToLower(filepath.Ext(filePath))
	if analyzer, ok := r.byExt[ext]; ok {
		return analyzer
	}
	return nil
}

// GetAnalyzerByLanguage returns analyzer by language name
func (r *AnalyzerRegistryImpl) GetAnalyzerByLanguage(lang string) analysis.LanguageAnalyzer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.byLang[strings.ToLower(lang)]
}

// SupportedLanguages returns list of supported languages
func (r *AnalyzerRegistryImpl) SupportedLanguages() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	langs := make([]string, 0, len(r.byLang))
	for lang := range r.byLang {
		langs = append(langs, lang)
	}
	return langs
}

// SupportedExtensions returns list of supported extensions
func (r *AnalyzerRegistryImpl) SupportedExtensions() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exts := make([]string, 0, len(r.byExt))
	for ext := range r.byExt {
		exts = append(exts, ext)
	}
	return exts
}
