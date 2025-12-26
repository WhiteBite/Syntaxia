package analyzers

import (
	"context"
	"regexp"
	"strings"

	"syntaxia/domain/analysis"
)

// PHPAnalyzer analyzes PHP files
type PHPAnalyzer struct {
	classRe     *regexp.Regexp
	interfaceRe *regexp.Regexp
	traitRe     *regexp.Regexp
	functionRe  *regexp.Regexp
	methodRe    *regexp.Regexp
	namespaceRe *regexp.Regexp
	useRe       *regexp.Regexp
	constRe     *regexp.Regexp
	enumRe      *regexp.Regexp
}

func NewPHPAnalyzer() *PHPAnalyzer {
	return &PHPAnalyzer{
		// class ClassName extends Base implements Interface
		classRe: regexp.MustCompile(`(?m)^[\t ]*(abstract\s+|final\s+)?(class)\s+(\w+)`),
		// interface InterfaceName
		interfaceRe: regexp.MustCompile(`(?m)^[\t ]*interface\s+(\w+)`),
		// trait TraitName
		traitRe: regexp.MustCompile(`(?m)^[\t ]*trait\s+(\w+)`),
		// function functionName(params)
		functionRe: regexp.MustCompile(`(?m)^[\t ]*function\s+(\w+)\s*\(`),
		// public/private/protected function methodName(params)
		methodRe: regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*(static\s+)?function\s+(\w+)\s*\(`),
		// namespace Name\Space;
		namespaceRe: regexp.MustCompile(`(?m)^[\t ]*namespace\s+([\w\\]+)\s*;`),
		// use Namespace\Class;
		useRe: regexp.MustCompile(`(?m)^[\t ]*use\s+([\w\\]+)(?:\s+as\s+(\w+))?\s*;`),
		// const CONSTANT_NAME = value;
		constRe: regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*const\s+(\w+)\s*=`),
		// enum EnumName (PHP 8.1+)
		enumRe: regexp.MustCompile(`(?m)^[\t ]*enum\s+(\w+)`),
	}
}

func (a *PHPAnalyzer) Language() string     { return "php" }
func (a *PHPAnalyzer) Extensions() []string { return []string{".php", ".phtml", ".php3", ".php4", ".php5", ".phps"} }
func (a *PHPAnalyzer) CanAnalyze(filePath string) bool {
	lower := strings.ToLower(filePath)
	for _, ext := range a.Extensions() {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

func (a *PHPAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Namespace
	if match := a.namespaceRe.FindStringSubmatch(text); match != nil {
		line := countLines(content[:a.namespaceRe.FindStringIndex(text)[0]])
		symbols = append(symbols, analysis.Symbol{
			Name:      match[1],
			Kind:      analysis.KindModule,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	// Classes
	for _, match := range a.classRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		isAbstract := match[2] != -1 && strings.Contains(text[match[2]:match[3]], "abstract")
		extra := make(map[string]string)
		if isAbstract {
			extra["abstract"] = "true"
		}
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindClass,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     extra,
		})
	}

	// Interfaces
	for _, match := range a.interfaceRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindInterface,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Traits
	for _, match := range a.traitRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindClass,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     map[string]string{"trait": "true"},
		})
	}

	// Enums (PHP 8.1+)
	for _, match := range a.enumRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[2]:match[3]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindEnum,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
		})
	}

	// Methods (with visibility modifiers)
	for _, match := range a.methodRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		// Skip constructor/destructor magic methods for cleaner output
		if name == "__construct" || name == "__destruct" {
			continue
		}
		endLine := findBlockEndLine(lines, line-1)
		visibility := "public"
		if match[2] != -1 {
			visibility = text[match[2]:match[3]]
		}
		isStatic := match[4] != -1 && match[5] != -1
		extra := map[string]string{"visibility": visibility}
		if isStatic {
			extra["static"] = "true"
		}
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindMethod,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   endLine,
			Extra:     extra,
		})
	}

	// Constants
	for _, match := range a.constRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		symbols = append(symbols, analysis.Symbol{
			Name:      name,
			Kind:      analysis.KindConstant,
			Language:  "php",
			FilePath:  filePath,
			StartLine: line,
			EndLine:   line,
		})
	}

	return symbols, nil
}

func (a *PHPAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.useRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := match[1]
		alias := ""
		if len(match) > 2 && match[2] != "" {
			alias = match[2]
		}
		// Local if it's in the same namespace or starts with App\
		isLocal := strings.HasPrefix(path, "App\\") || !strings.Contains(path, "\\")
		imports = append(imports, analysis.Import{
			Path:    path,
			Alias:   alias,
			IsLocal: isLocal,
		})
	}
	return imports, nil
}

func (a *PHPAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		// In PHP, public classes/interfaces/traits are exported
		// Methods with public visibility are also exported
		isPublic := true
		if sym.Extra != nil {
			if vis, ok := sym.Extra["visibility"]; ok && vis != "public" {
				isPublic = false
			}
		}
		if isPublic && (sym.Kind == analysis.KindClass || sym.Kind == analysis.KindInterface ||
			sym.Kind == analysis.KindEnum || sym.Kind == analysis.KindMethod) {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

func (a *PHPAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	lines := strings.Split(text, "\n")

	// Find function or method
	funcRe := regexp.MustCompile(`(?m)^[\t ]*(public|private|protected)?\s*(static\s+)?function\s+` + regexp.QuoteMeta(funcName) + `\s*\(`)

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
