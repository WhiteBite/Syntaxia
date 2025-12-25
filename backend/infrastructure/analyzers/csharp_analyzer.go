package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// CSharpAnalyzer analyzes C# files
type CSharpAnalyzer struct {
	classRe     *regexp.Regexp
	interfaceRe *regexp.Regexp
	structRe    *regexp.Regexp
	enumRe      *regexp.Regexp
	methodRe    *regexp.Regexp
	propertyRe  *regexp.Regexp
	namespaceRe *regexp.Regexp
	usingRe     *regexp.Regexp
	recordRe    *regexp.Regexp
}

func NewCSharpAnalyzer() *CSharpAnalyzer {
	return &CSharpAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*(static\s+)?(abstract\s+|sealed\s+)?class\s+(\w+)`),
		interfaceRe: regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*interface\s+(\w+)`),
		structRe:    regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*(readonly\s+)?struct\s+(\w+)`),
		enumRe:      regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*enum\s+(\w+)`),
		methodRe:    regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*(static\s+)?(async\s+)?(virtual\s+|override\s+|abstract\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\([^)]*\)`),
		propertyRe:  regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*(static\s+)?[\w<>\[\],\s]+\s+(\w+)\s*\{\s*(get|set)`),
		namespaceRe: regexp.MustCompile(`(?m)^namespace\s+([\w.]+)`),
		usingRe:     regexp.MustCompile(`(?m)^using\s+([\w.]+)\s*;`),
		recordRe:    regexp.MustCompile(`(?m)^[\t ]*(public|private|protected|internal)?\s*record\s+(\w+)`),
	}
}

func (a *CSharpAnalyzer) Language() string     { return "csharp" }
func (a *CSharpAnalyzer) Extensions() []string { return []string{".cs"} }
func (a *CSharpAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".cs")
}

func (a *CSharpAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Classes
	for _, match := range a.classRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[8]:match[9]]
		endLine := findBlockEndLine(lines, line-1)
		visibility := "internal"
		if match[2] != -1 {
			visibility = text[match[2]:match[3]]
		}
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindClass,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     map[string]string{"visibility": visibility},
		})
	}

	// Interfaces
	for _, match := range a.interfaceRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindInterface,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Structs
	for _, match := range a.structRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindStruct,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Enums
	for _, match := range a.enumRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindEnum,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Records (C# 9+)
	for _, match := range a.recordRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindClass,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     map[string]string{"record": "true"},
		})
	}

	// Methods (simplified - may have false positives)
	methodMatches := a.methodRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range methodMatches {
		line := countLines(content[:match[0]])
		name := text[match[10]:match[11]]
		// Skip common false positives
		if name == "if" || name == "for" || name == "while" || name == "switch" || name == "catch" || name == "using" {
			continue
		}
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindMethod,
			Language:  "csharp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	return symbols, nil
}

func (a *CSharpAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.usingRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := match[1]
		// System.* and Microsoft.* are external
		isLocal := !strings.HasPrefix(path, "System") && !strings.HasPrefix(path, "Microsoft")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

func (a *CSharpAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		// In C#, public symbols are exported
		if sym.Extra != nil && sym.Extra["visibility"] == "public" {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

func (a *CSharpAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	funcRe := regexp.MustCompile(`(?m)[\w<>\[\],\s]+\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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
