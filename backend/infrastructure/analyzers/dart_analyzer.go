package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// DartAnalyzer analyzes Dart/Flutter files
type DartAnalyzer struct {
	classRe    *regexp.Regexp
	widgetRe   *regexp.Regexp
	functionRe *regexp.Regexp
	enumRe     *regexp.Regexp
	importRe   *regexp.Regexp
}

func NewDartAnalyzer() *DartAnalyzer {
	return &DartAnalyzer{
		classRe:    regexp.MustCompile(`(?m)^[\t ]*(abstract\s+)?class\s+(\w+)`),
		widgetRe:   regexp.MustCompile(`(?m)^[\t ]*class\s+(\w+)\s+extends\s+(StatelessWidget|StatefulWidget)`),
		functionRe: regexp.MustCompile(`(?m)^([\w<>?]+)\s+(\w+)\s*\([^)]*\)\s*(?:async\s*)?\{`),
		enumRe:     regexp.MustCompile(`(?m)^[\t ]*enum\s+(\w+)`),
		importRe:   regexp.MustCompile(`(?m)^import\s+['"]([^'"]+)['"]`),
	}
}

func (a *DartAnalyzer) Language() string     { return "dart" }
func (a *DartAnalyzer) Extensions() []string { return []string{".dart"} }
func (a *DartAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".dart")
}

func (a *DartAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	widgetMatches := a.widgetRe.FindAllStringSubmatchIndex(text, -1)
	classMatches := a.classRe.FindAllStringSubmatchIndex(text, -1)
	enumMatches := a.enumRe.FindAllStringSubmatchIndex(text, -1)
	symbols := make([]analysis.Symbol, 0, len(widgetMatches)+len(classMatches)+len(enumMatches))
	lines := strings.Split(text, "\n")

	// Widgets first
	widgetNames := make(map[string]bool)
	for _, match := range widgetMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		widgetType := text[match[4]:match[5]]
		widgetNames[name] = true
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindWidget,
			Language:  "dart",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     map[string]string{"widgetType": widgetType},
		})
	}

	// Classes (excluding widgets)
	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		if widgetNames[name] {
			continue
		}
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindClass, Language: "dart", FilePath: filePath, StartLine: line, EndLine: endLine})
	}

	// Enums
	for _, match := range enumMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindEnum, Language: "dart", FilePath: filePath, StartLine: line, EndLine: endLine})
	}

	return symbols, nil
}

func (a *DartAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := match[1]
		isLocal := strings.HasPrefix(path, "package:") || strings.HasPrefix(path, "./")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

// GetExports extracts exported symbols from Dart files
// In Dart, all public symbols (not starting with _) are exported by default
func (a *DartAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	exports := make([]analysis.Export, 0, len(symbols))
	for _, sym := range symbols {
		// In Dart, symbols starting with _ are private
		if !strings.HasPrefix(sym.Name, "_") {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

// GetFunctionBody returns the body of a function by name
func (a *DartAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	// Find function by name
	funcRe := regexp.MustCompile(`(?m)^[\t ]*[\w<>?]+\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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

	endLine := findBlockEndLine(lines, startLine)

	var body strings.Builder
	for i := startLine; i < endLine && i < len(lines); i++ {
		body.WriteString(lines[i])
		if i < endLine-1 {
			body.WriteString("\n")
		}
	}

	return body.String(), startLine + 1, endLine, nil
}
