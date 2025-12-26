package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// PythonAnalyzer analyzes Python files
type PythonAnalyzer struct {
	classRe    *regexp.Regexp
	functionRe *regexp.Regexp
	methodRe   *regexp.Regexp
	importRe   *regexp.Regexp
	fromImport *regexp.Regexp
	decoratorRe *regexp.Regexp
}

func NewPythonAnalyzer() *PythonAnalyzer {
	return &PythonAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^class\s+(\w+)(?:\s*\([^)]*\))?\s*:`),
		functionRe:  regexp.MustCompile(`(?m)^def\s+(\w+)\s*\([^)]*\)\s*(?:->\s*[^:]+)?\s*:`),
		methodRe:    regexp.MustCompile(`(?m)^[\t ]+def\s+(\w+)\s*\([^)]*\)\s*(?:->\s*[^:]+)?\s*:`),
		importRe:    regexp.MustCompile(`(?m)^import\s+(\S+)`),
		fromImport:  regexp.MustCompile(`(?m)^from\s+(\S+)\s+import\s+(.+)`),
		decoratorRe: regexp.MustCompile(`(?m)^@(\w+)`),
	}
}

func (a *PythonAnalyzer) Language() string     { return "python" }
func (a *PythonAnalyzer) Extensions() []string { return []string{".py", ".pyw", ".pyi"} }
func (a *PythonAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".py") ||
		strings.HasSuffix(filePath, ".pyw") ||
		strings.HasSuffix(filePath, ".pyi")
}

func (a *PythonAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Extract classes
	classMatches := a.classRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findPythonBlockEnd(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindClass,
			Language:  "python",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Extract top-level functions (not methods)
	funcMatches := a.functionRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range funcMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findPythonBlockEnd(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindFunction,
			Language:  "python",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Extract methods (indented def)
	methodMatches := a.methodRe.FindAllStringSubmatchIndex(text, -1)
	for _, match := range methodMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findPythonBlockEnd(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindMethod,
			Language:  "python",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	return symbols, nil
}

func (a *PythonAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	text := string(content)
	imports := make([]analysis.Import, 0)

	// import module
	matches := a.importRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		path := match[1]
		isLocal := !strings.Contains(path, ".")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}

	// from module import ...
	fromMatches := a.fromImport.FindAllStringSubmatch(text, -1)
	for _, match := range fromMatches {
		path := match[1]
		isLocal := strings.HasPrefix(path, ".")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}

	return imports, nil
}

func (a *PythonAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	exports := make([]analysis.Export, 0, len(symbols))
	for _, sym := range symbols {
		// In Python, symbols not starting with _ are public
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

func (a *PythonAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	funcRe := regexp.MustCompile(`(?m)^[\t ]*def\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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

	endLine := findPythonBlockEnd(lines, startLine)

	var body strings.Builder
	for i := startLine; i < endLine && i < len(lines); i++ {
		body.WriteString(lines[i])
		if i < endLine-1 {
			body.WriteString("\n")
		}
	}

	return body.String(), startLine + 1, endLine, nil
}

// findPythonBlockEnd finds the end of a Python block based on indentation
func findPythonBlockEnd(lines []string, startLine int) int {
	if startLine >= len(lines) {
		return startLine + 1
	}

	// Get the indentation of the starting line
	startIndent := getIndentation(lines[startLine])

	for i := startLine + 1; i < len(lines); i++ {
		line := lines[i]
		// Skip empty lines and comments
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		currentIndent := getIndentation(line)
		// If we find a line with same or less indentation, block ends
		if currentIndent <= startIndent {
			return i
		}
	}

	return len(lines)
}

// getIndentation returns the number of leading spaces/tabs
func getIndentation(line string) int {
	count := 0
	for _, ch := range line {
		if ch == ' ' {
			count++
		} else if ch == '\t' {
			count += 4 // Treat tab as 4 spaces
		} else {
			break
		}
	}
	return count
}
