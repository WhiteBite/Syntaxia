package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// JavaAnalyzer analyzes Java files
type JavaAnalyzer struct {
	classRe     *regexp.Regexp
	methodRe    *regexp.Regexp
	interfaceRe *regexp.Regexp
	enumRe      *regexp.Regexp
	importRe    *regexp.Regexp
}

func NewJavaAnalyzer() *JavaAnalyzer {
	return &JavaAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*(abstract|final)?\s*class\s+(\w+)`),
		interfaceRe: regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*interface\s+(\w+)`),
		methodRe:    regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*(static)?\s*[\w<>\[\]]+\s+(\w+)\s*\([^)]*\)\s*[{;]`),
		enumRe:      regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*enum\s+(\w+)`),
		importRe:    regexp.MustCompile(`(?m)^import\s+([\w.]+);`),
	}
}

func (a *JavaAnalyzer) Language() string     { return "java" }
func (a *JavaAnalyzer) Extensions() []string { return []string{".java"} }
func (a *JavaAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".java")
}

func (a *JavaAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	classMatches := a.classRe.FindAllStringSubmatchIndex(text, -1)
	ifaceMatches := a.interfaceRe.FindAllStringSubmatchIndex(text, -1)
	enumMatches := a.enumRe.FindAllStringSubmatchIndex(text, -1)
	symbols := make([]analysis.Symbol, 0, len(classMatches)+len(ifaceMatches)+len(enumMatches))
	lines := strings.Split(text, "\n")

	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindClass, Language: "java", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	for _, match := range ifaceMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindInterface, Language: "java", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	for _, match := range enumMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindEnum, Language: "java", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	return symbols, nil
}

func (a *JavaAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		imports = append(imports, analysis.Import{Path: match[1], IsLocal: false})
	}
	return imports, nil
}

// KotlinAnalyzer analyzes Kotlin files
type KotlinAnalyzer struct {
	classRe  *regexp.Regexp
	funRe    *regexp.Regexp
	objectRe *regexp.Regexp
	importRe *regexp.Regexp
}

func NewKotlinAnalyzer() *KotlinAnalyzer {
	return &KotlinAnalyzer{
		classRe:  regexp.MustCompile(`(?m)^[\t ]*(data\s+|sealed\s+|open\s+)?class\s+(\w+)`),
		funRe:    regexp.MustCompile(`(?m)^[\t ]*(private|public|internal)?\s*(suspend\s+)?fun\s+(\w+)`),
		objectRe: regexp.MustCompile(`(?m)^[\t ]*(companion\s+)?object\s+(\w+)?`),
		importRe: regexp.MustCompile(`(?m)^import\s+([\w.]+)`),
	}
}

func (a *KotlinAnalyzer) Language() string     { return "kotlin" }
func (a *KotlinAnalyzer) Extensions() []string { return []string{".kt", ".kts"} }
func (a *KotlinAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".kt") || strings.HasSuffix(filePath, ".kts")
}

func (a *KotlinAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	classMatches := a.classRe.FindAllStringSubmatchIndex(text, -1)
	funMatches := a.funRe.FindAllStringSubmatchIndex(text, -1)
	symbols := make([]analysis.Symbol, 0, len(classMatches)+len(funMatches))
	lines := strings.Split(text, "\n")

	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindClass, Language: "kotlin", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	for _, match := range funMatches {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindFunction, Language: "kotlin", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	return symbols, nil
}

func (a *KotlinAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		imports = append(imports, analysis.Import{Path: match[1], IsLocal: false})
	}
	return imports, nil
}

// Helper
func countLines(data []byte) int {
	return strings.Count(string(data), "\n") + 1
}

// GetExports returns public symbols for Java
func (a *JavaAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	text := string(content)

	// Public classes
	publicClassRe := regexp.MustCompile(`(?m)^[\t ]*public\s+(abstract\s+|final\s+)?class\s+(\w+)`)
	classMatches := publicClassRe.FindAllStringSubmatchIndex(text, -1)

	// Public interfaces
	publicInterfaceRe := regexp.MustCompile(`(?m)^[\t ]*public\s+interface\s+(\w+)`)
	ifaceMatches := publicInterfaceRe.FindAllStringSubmatchIndex(text, -1)

	// Public methods
	publicMethodRe := regexp.MustCompile(`(?m)^[\t ]*public\s+(static\s+)?[\w<>\[\]]+\s+(\w+)\s*\(`)
	methodMatches := publicMethodRe.FindAllStringSubmatchIndex(text, -1)

	exports := make([]analysis.Export, 0, len(classMatches)+len(ifaceMatches)+len(methodMatches))

	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		exports = append(exports, analysis.Export{Name: name, Kind: "class", Line: line})
	}

	for _, match := range ifaceMatches {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		exports = append(exports, analysis.Export{Name: name, Kind: "interface", Line: line})
	}

	for _, match := range methodMatches {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		exports = append(exports, analysis.Export{Name: name, Kind: "method", Line: line})
	}

	return exports, nil
}

// GetFunctionBody returns method body for Java
func (a *JavaAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	return extractJVMFunctionBody(content, funcName)
}

// GetExports returns public symbols for Kotlin
func (a *KotlinAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	text := string(content)

	// Public/internal classes (Kotlin default is public)
	classRe := regexp.MustCompile(`(?m)^[\t ]*(public\s+|internal\s+)?(data\s+|sealed\s+|open\s+)?class\s+(\w+)`)
	classMatches := classRe.FindAllStringSubmatchIndex(text, -1)

	// Public functions
	funRe := regexp.MustCompile(`(?m)^[\t ]*(public\s+|internal\s+)?(suspend\s+)?fun\s+(\w+)`)
	funMatches := funRe.FindAllStringSubmatchIndex(text, -1)

	exports := make([]analysis.Export, 0, len(classMatches)+len(funMatches))

	for _, match := range classMatches {
		line := countLines(content[:match[0]])
		// Check if not private
		lineStart := match[0]
		lineText := text[lineStart : lineStart+min(100, len(text)-lineStart)]
		if !strings.Contains(lineText, "private") {
			name := text[match[6]:match[7]]
			exports = append(exports, analysis.Export{Name: name, Kind: "class", Line: line})
		}
	}

	for _, match := range funMatches {
		line := countLines(content[:match[0]])
		lineStart := match[0]
		lineText := text[lineStart : lineStart+min(100, len(text)-lineStart)]
		if !strings.Contains(lineText, "private") {
			name := text[match[6]:match[7]]
			exports = append(exports, analysis.Export{Name: name, Kind: "function", Line: line})
		}
	}

	return exports, nil
}

// GetFunctionBody returns function body for Kotlin
func (a *KotlinAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	return extractJVMFunctionBody(content, funcName)
}

// extractJVMFunctionBody extracts function/method body for Java/Kotlin
func extractJVMFunctionBody(content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	// Find function start
	funcRe := regexp.MustCompile(`(?m)[\t ]*(?:public|private|protected|internal)?\s*(?:static|suspend|override)?\s*(?:fun|[\w<>\[\]]+)\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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

	// Find matching braces
	braceCount := 0
	started := false
	endLine := startLine

	for i := startLine; i < len(lines); i++ {
		line := lines[i]
		for _, ch := range line {
			if ch == '{' {
				braceCount++
				started = true
			} else if ch == '}' {
				braceCount--
			}
		}
		if started && braceCount == 0 {
			endLine = i
			break
		}
	}

	var body strings.Builder
	for i := startLine; i <= endLine && i < len(lines); i++ {
		body.WriteString(lines[i])
		if i < endLine {
			body.WriteString("\n")
		}
	}

	return body.String(), startLine + 1, endLine + 1, nil
}
