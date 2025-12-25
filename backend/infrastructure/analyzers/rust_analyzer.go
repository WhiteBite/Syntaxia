package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// RustAnalyzer analyzes Rust files
type RustAnalyzer struct {
	structRe   *regexp.Regexp
	enumRe     *regexp.Regexp
	traitRe    *regexp.Regexp
	implRe     *regexp.Regexp
	functionRe *regexp.Regexp
	modRe      *regexp.Regexp
	useRe      *regexp.Regexp
	constRe    *regexp.Regexp
	typeRe     *regexp.Regexp
}

func NewRustAnalyzer() *RustAnalyzer {
	return &RustAnalyzer{
		structRe:   regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?struct\s+(\w+)`),
		enumRe:     regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?enum\s+(\w+)`),
		traitRe:    regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?trait\s+(\w+)`),
		implRe:     regexp.MustCompile(`(?m)^[\t ]*impl(?:<[^>]+>)?\s+(?:(\w+)\s+for\s+)?(\w+)`),
		functionRe: regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?(async\s+)?fn\s+(\w+)`),
		modRe:      regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?mod\s+(\w+)`),
		useRe:      regexp.MustCompile(`(?m)^use\s+([^;]+);`),
		constRe:    regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?const\s+(\w+)`),
		typeRe:     regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?type\s+(\w+)`),
	}
}

func (a *RustAnalyzer) Language() string     { return "rust" }
func (a *RustAnalyzer) Extensions() []string { return []string{".rs"} }
func (a *RustAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".rs")
}

func (a *RustAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Structs
	for _, match := range a.structRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		isPublic := match[2] != -1 && match[3] != -1
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindStruct,
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     map[string]string{"public": boolToStr(isPublic)},
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
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Traits
	for _, match := range a.traitRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindInterface,
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Functions
	for _, match := range a.functionRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		isPublic := match[2] != -1 && match[3] != -1
		isAsync := match[4] != -1 && match[5] != -1
		extra := map[string]string{"public": boolToStr(isPublic)}
		if isAsync {
			extra["async"] = "true"
		}
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindFunction,
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     extra,
		})
	}

	// Consts
	for _, match := range a.constRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindConstant,
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	// Type aliases
	for _, match := range a.typeRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindType,
			Language:  "rust",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	return symbols, nil
}

func (a *RustAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.useRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := strings.TrimSpace(match[1])
		// Local if starts with crate::, self::, super::
		isLocal := strings.HasPrefix(path, "crate::") ||
			strings.HasPrefix(path, "self::") ||
			strings.HasPrefix(path, "super::")
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

func (a *RustAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		// In Rust, pub symbols are exported
		if sym.Extra != nil && sym.Extra["public"] == "true" {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

func (a *RustAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	funcRe := regexp.MustCompile(`(?m)^[\t ]*(pub\s+)?(async\s+)?fn\s+` + regexp.QuoteMeta(funcName) + `\s*[<(]`)

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

func boolToStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
