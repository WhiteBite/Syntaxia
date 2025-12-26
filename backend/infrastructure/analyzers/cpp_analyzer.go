package analyzers

import (
	"context"
	"regexp"
	"strings"
	"syntaxia/domain/analysis"
)

// CppAnalyzer analyzes C/C++ files using regex patterns
type CppAnalyzer struct {
	classRe     *regexp.Regexp
	structRe    *regexp.Regexp
	enumRe      *regexp.Regexp
	functionRe  *regexp.Regexp
	methodRe    *regexp.Regexp
	namespaceRe *regexp.Regexp
	includeRe   *regexp.Regexp
	defineRe    *regexp.Regexp
	typedefRe   *regexp.Regexp
	templateRe  *regexp.Regexp
}

// NewCppAnalyzer creates a new C/C++ analyzer
func NewCppAnalyzer() *CppAnalyzer {
	return &CppAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^[\t ]*(template\s*<[^>]*>\s*)?(class|struct)\s+(\w+)(?:\s*:\s*(?:public|private|protected)\s+\w+)?`),
		structRe:    regexp.MustCompile(`(?m)^[\t ]*struct\s+(\w+)\s*\{`),
		enumRe:      regexp.MustCompile(`(?m)^[\t ]*enum\s+(class\s+)?(\w+)`),
		functionRe:  regexp.MustCompile(`(?m)^[\t ]*(?:static\s+|inline\s+|virtual\s+|explicit\s+|constexpr\s+)*(?:[\w:*&<>,\s]+)\s+(\w+)\s*\([^)]*\)\s*(?:const\s*)?(?:noexcept\s*)?(?:override\s*)?(?:final\s*)?(?:\{|;)`),
		methodRe:    regexp.MustCompile(`(?m)^[\t ]*(?:[\w:*&<>,\s]+)\s+(\w+)::(\w+)\s*\([^)]*\)`),
		namespaceRe: regexp.MustCompile(`(?m)^[\t ]*namespace\s+(\w+)`),
		includeRe:   regexp.MustCompile(`(?m)^[\t ]*#include\s*[<"]([^>"]+)[>"]`),
		defineRe:    regexp.MustCompile(`(?m)^[\t ]*#define\s+(\w+)`),
		typedefRe:   regexp.MustCompile(`(?m)^[\t ]*typedef\s+.+\s+(\w+)\s*;`),
		templateRe:  regexp.MustCompile(`(?m)^[\t ]*template\s*<[^>]*>`),
	}
}

func (a *CppAnalyzer) Language() string     { return "cpp" }
func (a *CppAnalyzer) Extensions() []string { return []string{".cpp", ".cc", ".cxx", ".c", ".h", ".hpp", ".hxx"} }
func (a *CppAnalyzer) CanAnalyze(filePath string) bool {
	lower := strings.ToLower(filePath)
	for _, ext := range a.Extensions() {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (a *CppAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Classes and structs
	for _, match := range a.classRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		// Group 3 is the class/struct name
		nameStart, nameEnd := match[6], match[7]
		if nameStart < 0 || nameEnd < 0 {
			continue
		}
		name := text[nameStart:nameEnd]
		endLine := findBlockEndLine(lines, line-1)
		kind := analysis.KindClass
		// Check if it's a struct
		typeStart, typeEnd := match[4], match[5]
		if typeStart >= 0 && typeEnd >= 0 && text[typeStart:typeEnd] == "struct" {
			kind = analysis.KindStruct
		}
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      kind,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Enums
	for _, match := range a.enumRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		// Group 2 is the enum name
		nameStart, nameEnd := match[4], match[5]
		if nameStart < 0 || nameEnd < 0 {
			continue
		}
		name := text[nameStart:nameEnd]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindEnum,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Namespaces
	for _, match := range a.namespaceRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindModule, // namespace maps to module
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Functions (standalone)
	for _, match := range a.functionRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		// Skip common keywords that might be matched
		if isReservedKeyword(name) {
			continue
		}
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindFunction,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Methods (class::method)
	for _, match := range a.methodRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		className := text[match[2]:match[3]]
		methodName := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      methodName,
			Kind:      analysis.KindMethod,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Parent:    className,
		})
	}

	// Defines (macros)
	for _, match := range a.defineRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindConstant,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	// Typedefs
	for _, match := range a.typedefRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindType,
			Language:  "cpp",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	return symbols, nil
}

func (a *CppAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.includeRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := strings.TrimSpace(match[1])
		// Local if uses quotes or doesn't contain /
		isLocal := !strings.Contains(path, "/") || strings.HasPrefix(path, ".")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

func (a *CppAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	// In C++, symbols in header files are typically exported
	isHeader := strings.HasSuffix(filePath, ".h") ||
		strings.HasSuffix(filePath, ".hpp") ||
		strings.HasSuffix(filePath, ".hxx")

	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		// Export all symbols from headers, or public symbols from source files
		if isHeader || sym.Kind == analysis.KindClass || sym.Kind == analysis.KindStruct {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

func (a *CppAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	lines := strings.Split(string(content), "\n")

	funcRe := regexp.MustCompile(`(?m)^[\t ]*(?:[\w:*&<>,\s]+)\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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

// isReservedKeyword checks if a name is a C++ reserved keyword
func isReservedKeyword(name string) bool {
	keywords := map[string]bool{
		"if": true, "else": true, "for": true, "while": true, "do": true,
		"switch": true, "case": true, "default": true, "break": true, "continue": true,
		"return": true, "goto": true, "try": true, "catch": true, "throw": true,
		"new": true, "delete": true, "sizeof": true, "typeof": true,
		"class": true, "struct": true, "union": true, "enum": true,
		"public": true, "private": true, "protected": true,
		"virtual": true, "override": true, "final": true,
		"static": true, "const": true, "volatile": true, "mutable": true,
		"inline": true, "extern": true, "register": true,
		"namespace": true, "using": true, "typedef": true,
		"template": true, "typename": true,
		"true": true, "false": true, "nullptr": true,
		"and": true, "or": true, "not": true,
	}
	return keywords[name]
}
