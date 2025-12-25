package analyzers

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// ReferenceFinder finds references to symbols across the project
type ReferenceFinder struct {
	registry analysis.AnalyzerRegistry
}

// NewReferenceFinder creates a new reference finder
func NewReferenceFinder(registry analysis.AnalyzerRegistry) *ReferenceFinder {
	return &ReferenceFinder{registry: registry}
}

// Reference represents a reference to a symbol
type Reference struct {
	FilePath     string `json:"filePath"`
	Line         int    `json:"line"`
	Column       int    `json:"column"`
	LineText     string `json:"lineText"`
	Context      string `json:"context"` // surrounding context
	IsDefinition bool   `json:"isDefinition"`
}

// skipDirs contains directories to skip during reference search
var refFinderSkipDirs = map[string]bool{
	"node_modules": true, "vendor": true, "build": true, "dist": true,
}

// getLineContext extracts surrounding lines for context
func getLineContext(lines []string, lineIdx int) string {
	start := lineIdx - 2
	if start < 0 {
		start = 0
	}
	end := lineIdx + 3
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[start:end], "\n")
}

// isDefinition checks if a reference is the symbol definition
func (rf *ReferenceFinder) isDefinition(ctx context.Context, analyzer analysis.LanguageAnalyzer, relPath string, content []byte, symbolName string, lineNum int) bool {
	symbols, _ := analyzer.ExtractSymbols(ctx, relPath, content)
	for _, sym := range symbols {
		if sym.Name == symbolName && sym.StartLine == lineNum {
			return true
		}
	}
	return false
}

// findReferencesInFile finds references in a single file
func (rf *ReferenceFinder) findReferencesInFile(ctx context.Context, pattern *regexp.Regexp, path, relPath string, symbolName string, symbolKind analysis.SymbolKind, maxRefs int) ([]Reference, bool) {
	analyzer := rf.registry.GetAnalyzer(path)
	if analyzer == nil {
		return nil, false
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	lines := strings.Split(string(content), "\n")
	var refs []Reference

	for i, line := range lines {
		for _, match := range pattern.FindAllStringIndex(line, -1) {
			ref := Reference{
				FilePath: relPath,
				Line:     i + 1,
				Column:   match[0] + 1,
				LineText: strings.TrimSpace(line),
				Context:  getLineContext(lines, i),
			}

			if symbolKind != "" {
				ref.IsDefinition = rf.isDefinition(ctx, analyzer, relPath, content, symbolName, i+1)
			}

			refs = append(refs, ref)
			if len(refs) >= maxRefs {
				return refs, true
			}
		}
	}
	return refs, false
}

// FindReferences finds all references to a symbol in the project
func (rf *ReferenceFinder) FindReferences(ctx context.Context, projectRoot string, symbolName string, symbolKind analysis.SymbolKind) ([]Reference, error) {
	pattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(symbolName) + `\b`)
	var references []Reference
	const maxResults = 50

	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") || refFinderSkipDirs[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, _ := filepath.Rel(projectRoot, path)
		refs, limitReached := rf.findReferencesInFile(ctx, pattern, path, relPath, symbolName, symbolKind, maxResults-len(references))
		references = append(references, refs...)

		if limitReached || len(references) >= maxResults {
			return filepath.SkipAll
		}
		return nil
	})

	if err != nil && !errors.Is(err, filepath.SkipAll) {
		return nil, err
	}

	return references, nil
}

// FindUsages finds where a symbol is used (excluding definition)
func (rf *ReferenceFinder) FindUsages(ctx context.Context, projectRoot string, symbolName string) ([]Reference, error) {
	refs, err := rf.FindReferences(ctx, projectRoot, symbolName, "")
	if err != nil {
		return nil, err
	}

	// Filter out definitions
	var usages []Reference
	for _, ref := range refs {
		if !ref.IsDefinition {
			usages = append(usages, ref)
		}
	}

	return usages, nil
}
