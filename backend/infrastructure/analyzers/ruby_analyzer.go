package analyzers

import (
	"context"
	"regexp"
	"strings"

	"syntaxia/domain/analysis"
)

// RubyAnalyzer analyzes Ruby files
type RubyAnalyzer struct {
	classRe    *regexp.Regexp
	moduleRe   *regexp.Regexp
	methodRe   *regexp.Regexp
	attrRe     *regexp.Regexp
	requireRe  *regexp.Regexp
	constantRe *regexp.Regexp
}

func NewRubyAnalyzer() *RubyAnalyzer {
	return &RubyAnalyzer{
		classRe:    regexp.MustCompile(`(?m)^[\t ]*class\s+(\w+)(?:\s*<\s*\w+)?`),
		moduleRe:   regexp.MustCompile(`(?m)^[\t ]*module\s+(\w+)`),
		methodRe:   regexp.MustCompile(`(?m)^[\t ]*def\s+(self\.)?(\w+[?!=]?)`),
		attrRe:     regexp.MustCompile(`(?m)^[\t ]*attr_(reader|writer|accessor)\s+(.+)`),
		requireRe:  regexp.MustCompile(`(?m)^[\t ]*require(?:_relative)?\s+['"]([^'"]+)['"]`),
		constantRe: regexp.MustCompile(`(?m)^[\t ]*([A-Z][A-Z0-9_]*)\s*=`),
	}
}

func (a *RubyAnalyzer) Language() string     { return "ruby" }
func (a *RubyAnalyzer) Extensions() []string { return []string{".rb", ".rake", ".gemspec"} }

func (a *RubyAnalyzer) CanAnalyze(filePath string) bool {
	lower := strings.ToLower(filePath)
	for _, ext := range a.Extensions() {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}


func (a *RubyAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Classes
	for _, match := range a.classRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findRubyBlockEnd(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindClass, Language: "ruby",
			FilePath: filePath, StartLine: line, EndLine: endLine,
		})
	}

	// Modules
	for _, match := range a.moduleRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findRubyBlockEnd(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindModule, Language: "ruby",
			FilePath: filePath, StartLine: line, EndLine: endLine,
		})
	}

	// Methods
	for _, match := range a.methodRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		isClassMethod := match[2] != -1 && match[3] != -1
		endLine := findRubyBlockEnd(lines, line-1)
		extra := make(map[string]string)
		if isClassMethod {
			extra["class_method"] = "true"
		}
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindMethod, Language: "ruby",
			FilePath: filePath, StartLine: line, EndLine: endLine, Extra: extra,
		})
	}

	// Constants
	for _, match := range a.constantRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindConstant, Language: "ruby",
			FilePath: filePath, StartLine: line, EndLine: line,
		})
	}

	return symbols, nil
}

func (a *RubyAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.requireRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := match[1]
		isLocal := strings.HasPrefix(path, "./") || strings.HasPrefix(path, "../")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

func (a *RubyAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}
	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		if sym.Kind == analysis.KindClass || sym.Kind == analysis.KindModule {
			exports = append(exports, analysis.Export{Name: sym.Name, Kind: string(sym.Kind), Line: sym.StartLine})
		}
	}
	return exports, nil
}

func (a *RubyAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	lines := strings.Split(string(content), "\n")
	funcRe := regexp.MustCompile(`(?m)^[\t ]*def\s+(?:self\.)?` + regexp.QuoteMeta(funcName) + `\b`)
	startLine := -1
	for i, line := range lines {
		if funcRe.MatchString(line) {
			startLine = i
			break
		}
	}
	if startLine < 0 {
		return "", 0, 0, nil
	}
	endLine := findRubyBlockEnd(lines, startLine)
	var body strings.Builder
	for i := startLine; i < endLine && i < len(lines); i++ {
		body.WriteString(lines[i])
		if i < endLine-1 {
			body.WriteString("\n")
		}
	}
	return body.String(), startLine + 1, endLine, nil
}

// findRubyBlockEnd finds the end of a Ruby block (def/class/module...end)
func findRubyBlockEnd(lines []string, startLineIdx int) int {
	if startLineIdx < 0 || startLineIdx >= len(lines) {
		return startLineIdx + 1
	}
	blockRe := regexp.MustCompile(`\b(def|class|module|if|unless|case|while|until|for|begin|do)\b`)
	endRe := regexp.MustCompile(`\bend\b`)
	depth := 0
	for i := startLineIdx; i < len(lines); i++ {
		line := lines[i]
		depth += len(blockRe.FindAllString(line, -1))
		depth -= len(endRe.FindAllString(line, -1))
		if depth <= 0 {
			return i + 1
		}
	}
	return len(lines)
}
