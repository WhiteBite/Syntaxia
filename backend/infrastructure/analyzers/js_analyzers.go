package analyzers

import (
	"context"
	"regexp"
	"syntaxia/domain/analysis"
	"strings"
)

// stripCommentsAndStrings removes comments and string literals from code
// to prevent false positives in regex-based parsing
func stripCommentsAndStrings(code string) string {
	var result strings.Builder
	lines := strings.Split(code, "\n")
	inBlockComment := false

	for _, line := range lines {
		processedLine := processLineForComments(line, &inBlockComment)
		result.WriteString(processedLine)
		result.WriteString("\n")
	}

	return result.String()
}

// handleBlockComment handles block comment state
func handleBlockComment(line string, i int, inBlockComment *bool) int {
	if i+1 < len(line) && line[i] == '*' && line[i+1] == '/' {
		*inBlockComment = false
		return i + 2
	}
	return i + 1
}

// handleStringChar handles string literal state
func handleStringChar(line string, i int, stringChar byte, inString *bool) (int, bool) {
	if line[i] == '\\' && i+1 < len(line) {
		return i + 2, false
	}
	if line[i] == stringChar {
		*inString = false
		return i + 1, true
	}
	return i + 1, false
}

// processLineForComments processes a single line, removing comments and strings
func processLineForComments(line string, inBlockComment *bool) string {
	result := make([]byte, 0, len(line))
	lineLen := len(line)
	inString := false
	var stringChar byte

	for i := 0; i < lineLen; {
		ch := line[i]

		if *inBlockComment {
			i = handleBlockComment(line, i, inBlockComment)
			continue
		}

		if inString {
			newI, emit := handleStringChar(line, i, stringChar, &inString)
			if emit {
				result = append(result, ch)
			}
			i = newI
			continue
		}

		nextCh := byte(0)
		if i+1 < lineLen {
			nextCh = line[i+1]
		}

		switch ch {
		case '"', '\'', '`':
			inString, stringChar = true, ch
			result = append(result, ch)
			i++
		case '/':
			if nextCh == '/' {
				return string(result)
			}
			if nextCh == '*' {
				*inBlockComment = true
				i += 2
			} else {
				result = append(result, ch)
				i++
			}
		default:
			result = append(result, ch)
			i++
		}
	}

	return string(result)
}

// TypeScriptAnalyzer analyzes TypeScript files
type TypeScriptAnalyzer struct {
	classRe     *regexp.Regexp
	interfaceRe *regexp.Regexp
	typeRe      *regexp.Regexp
	functionRe  *regexp.Regexp
	enumRe      *regexp.Regexp
	importRe    *regexp.Regexp
}

func NewTypeScriptAnalyzer() *TypeScriptAnalyzer {
	return &TypeScriptAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(abstract\s+)?class\s+(\w+)`),
		interfaceRe: regexp.MustCompile(`(?m)^[\t ]*(export\s+)?interface\s+(\w+)`),
		typeRe:      regexp.MustCompile(`(?m)^[\t ]*(export\s+)?type\s+(\w+)`),
		functionRe:  regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(async\s+)?function\s+(\w+)`),
		enumRe:      regexp.MustCompile(`(?m)^[\t ]*(export\s+)?enum\s+(\w+)`),
		importRe:    regexp.MustCompile(`(?m)^import\s+.*from\s+['"]([^'"]+)['"]`),
	}
}

func (a *TypeScriptAnalyzer) Language() string     { return "typescript" }
func (a *TypeScriptAnalyzer) Extensions() []string { return []string{".ts", ".tsx"} }
func (a *TypeScriptAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".ts") || strings.HasSuffix(filePath, ".tsx")
}

// symbolExtractContext holds context for symbol extraction
type symbolExtractContext struct {
	text          string
	content       []byte
	lines         []string
	strippedLines []string
	filePath      string
	language      string
}

// isValidMatch checks if a match position is in actual code (not comment/string)
func (ctx *symbolExtractContext) isValidMatch(matchStart int) bool {
	lineNum := countLines(ctx.content[:matchStart])
	if lineNum > 0 && lineNum <= len(ctx.strippedLines) && lineNum <= len(ctx.lines) {
		strippedLine := strings.TrimSpace(ctx.strippedLines[lineNum-1])
		originalLine := strings.TrimSpace(ctx.lines[lineNum-1])
		if len(strippedLine) == 0 && len(originalLine) > 0 {
			return false
		}
		if len(originalLine) > 0 && len(strippedLine) < len(originalLine)/3 {
			return false
		}
	}
	return true
}

// extractSymbolsWithRegex extracts symbols using a regex pattern
func (ctx *symbolExtractContext) extractSymbolsWithRegex(re *regexp.Regexp, kind analysis.SymbolKind, nameIdx int, useBlockEnd bool) []analysis.Symbol {
	matches := re.FindAllStringSubmatchIndex(ctx.text, -1)
	symbols := make([]analysis.Symbol, 0, len(matches))
	for _, match := range matches {
		if !ctx.isValidMatch(match[0]) {
			continue
		}
		line := countLines(ctx.content[:match[0]])
		name := ctx.text[match[nameIdx]:match[nameIdx+1]]
		var endLine int
		if useBlockEnd {
			endLine = findBlockEndLine(ctx.lines, line-1)
		} else {
			endLine = findTypeEndLine(ctx.lines, line-1)
		}
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: kind, Language: ctx.language,
			FilePath: ctx.filePath, StartLine: line, EndLine: endLine,
		})
	}
	return symbols
}

func (a *TypeScriptAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	strippedText := stripCommentsAndStrings(text)

	extractCtx := &symbolExtractContext{
		text: text, content: content, lines: strings.Split(text, "\n"),
		strippedLines: strings.Split(strippedText, "\n"), filePath: filePath, language: "typescript",
	}

	symbols := make([]analysis.Symbol, 0, 32)
	symbols = append(symbols, extractCtx.extractSymbolsWithRegex(a.classRe, analysis.KindClass, 6, true)...)
	symbols = append(symbols, extractCtx.extractSymbolsWithRegex(a.interfaceRe, analysis.KindInterface, 4, true)...)
	symbols = append(symbols, extractCtx.extractSymbolsWithRegex(a.typeRe, analysis.KindType, 4, false)...)
	symbols = append(symbols, extractCtx.extractSymbolsWithRegex(a.functionRe, analysis.KindFunction, 6, true)...)
	symbols = append(symbols, extractCtx.extractSymbolsWithRegex(a.enumRe, analysis.KindEnum, 4, true)...)

	return symbols, nil
}

func (a *TypeScriptAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		imports = append(imports, analysis.Import{Path: match[1], IsLocal: strings.HasPrefix(match[1], ".")})
	}
	return imports, nil
}

// JavaScriptAnalyzer analyzes JavaScript files
type JavaScriptAnalyzer struct {
	classRe    *regexp.Regexp
	functionRe *regexp.Regexp
	importRe   *regexp.Regexp
}

func NewJavaScriptAnalyzer() *JavaScriptAnalyzer {
	return &JavaScriptAnalyzer{
		classRe:    regexp.MustCompile(`(?m)^[\t ]*(export\s+)?class\s+(\w+)`),
		functionRe: regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(async\s+)?function\s+(\w+)`),
		importRe:   regexp.MustCompile(`(?m)^import\s+.*from\s+['"]([^'"]+)['"]`),
	}
}

func (a *JavaScriptAnalyzer) Language() string     { return "javascript" }
func (a *JavaScriptAnalyzer) Extensions() []string { return []string{".js", ".jsx", ".mjs"} }
func (a *JavaScriptAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".js") || strings.HasSuffix(filePath, ".jsx") || strings.HasSuffix(filePath, ".mjs")
}

func (a *JavaScriptAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	classMatches := a.classRe.FindAllStringSubmatchIndex(text, -1)
	funcMatches := a.functionRe.FindAllStringSubmatchIndex(text, -1)
	symbols := make([]analysis.Symbol, 0, len(classMatches)+len(funcMatches))
	lines := strings.Split(text, "\n")

	// Strip comments and strings to avoid false positives
	strippedText := stripCommentsAndStrings(text)
	strippedLines := strings.Split(strippedText, "\n")

	// Helper to check if a match position is in actual code (not comment/string)
	isValidMatch := func(matchStart int) bool {
		lineNum := countLines(content[:matchStart])
		if lineNum > 0 && lineNum <= len(strippedLines) && lineNum <= len(lines) {
			strippedLine := strings.TrimSpace(strippedLines[lineNum-1])
			originalLine := strings.TrimSpace(lines[lineNum-1])
			// If stripped line is empty but original is not, the match is in a comment/string
			if len(strippedLine) == 0 && len(originalLine) > 0 {
				return false
			}
			// If stripped line is significantly shorter, the match might be in a comment
			if len(originalLine) > 0 && len(strippedLine) < len(originalLine)/3 {
				return false
			}
		}
		return true
	}

	for _, match := range classMatches {
		if !isValidMatch(match[0]) {
			continue
		}
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindClass, Language: "javascript", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	for _, match := range funcMatches {
		if !isValidMatch(match[0]) {
			continue
		}
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{Name: name, Kind: analysis.KindFunction, Language: "javascript", FilePath: filePath, StartLine: line, EndLine: endLine})
	}
	return symbols, nil
}

func (a *JavaScriptAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		imports = append(imports, analysis.Import{Path: match[1], IsLocal: strings.HasPrefix(match[1], ".")})
	}
	return imports, nil
}

// VueAnalyzer analyzes Vue SFC files
type VueAnalyzer struct {
	scriptRe     *regexp.Regexp
	composableRe *regexp.Regexp
	importRe     *regexp.Regexp
}

func NewVueAnalyzer() *VueAnalyzer {
	return &VueAnalyzer{
		scriptRe:     regexp.MustCompile(`(?s)<script[^>]*>(.*?)</script>`),
		composableRe: regexp.MustCompile(`(?m)(const|function)\s+(use\w+)`),
		importRe:     regexp.MustCompile(`(?m)^import\s+.*from\s+['"]([^'"]+)['"]`),
	}
}

func (a *VueAnalyzer) Language() string     { return "vue" }
func (a *VueAnalyzer) Extensions() []string { return []string{".vue"} }
func (a *VueAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".vue")
}

func (a *VueAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	var symbols []analysis.Symbol
	text := string(content)
	totalLines := strings.Count(text, "\n") + 1

	// Component name from filename
	parts := strings.Split(filePath, "/")
	if len(parts) == 0 {
		parts = strings.Split(filePath, "\\")
	}
	fileName := parts[len(parts)-1]
	componentName := strings.TrimSuffix(fileName, ".vue")

	symbols = append(symbols, analysis.Symbol{
		Name:      componentName,
		Kind:      analysis.KindComponent,
		Language:  "vue",
		FilePath:  filePath,
		StartLine: 1,
		EndLine:   totalLines,
	})

	// Extract composables from script
	if match := a.scriptRe.FindStringSubmatch(text); match != nil {
		script := match[1]
		scriptLines := strings.Split(script, "\n")
		for _, m := range a.composableRe.FindAllStringSubmatchIndex(script, -1) {
			line := countLines([]byte(script[:m[0]]))
			name := script[m[4]:m[5]]
			endLine := findBlockEndLine(scriptLines, line-1)
			symbols = append(symbols, analysis.Symbol{
				Name:      name,
				Kind:      analysis.KindComposable,
				Language:  "vue",
				FilePath:  filePath,
				StartLine: line,
				EndLine:   endLine,
			})
		}
	}
	return symbols, nil
}

func (a *VueAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	var imports []analysis.Import
	text := string(content)
	if match := a.scriptRe.FindStringSubmatch(text); match != nil {
		for _, m := range a.importRe.FindAllStringSubmatch(match[1], -1) {
			imports = append(imports, analysis.Import{Path: m[1], IsLocal: strings.HasPrefix(m[1], ".") || strings.HasPrefix(m[1], "@/")})
		}
	}
	return imports, nil
}

// Precompiled regexes for export extraction
var (
	exportDeclRe    = regexp.MustCompile(`(?m)^[\t ]*export\s+(default\s+)?(async\s+)?(function|class|const|let|var|type|interface|enum)\s+(\w+)`)
	exportListRe    = regexp.MustCompile(`(?m)^[\t ]*export\s*\{([^}]+)\}`)
	exportDefaultRe = regexp.MustCompile(`(?m)^[\t ]*export\s+default\s+(\w+)`)
	reExportRe      = regexp.MustCompile(`(?m)^[\t ]*export\s*\{([^}]+)\}\s*from\s*['"]([^'"]+)['"]`)
)

// parseExportList parses "name as alias" style export lists
func parseExportList(listStr string, lineNum int, kind string, fromPath string, isReExport bool) []analysis.Export {
	parts := strings.Split(listStr, ",")
	exports := make([]analysis.Export, 0, len(parts))
	for _, n := range parts {
		n = strings.TrimSpace(n)
		parts := strings.Split(n, " as ")
		name := strings.TrimSpace(parts[0])
		if name == "" {
			continue
		}
		exp := analysis.Export{Name: name, Kind: kind, Line: lineNum, IsReExport: isReExport, FromPath: fromPath}
		if len(parts) > 1 {
			exp.Alias = strings.TrimSpace(parts[1])
		}
		exports = append(exports, exp)
	}
	return exports
}

// GetExports extracts exported symbols from TypeScript files
func (a *TypeScriptAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	lines := strings.Split(string(content), "\n")
	exports := make([]analysis.Export, 0, len(lines)/10+1)

	for i, line := range lines {
		lineNum := i + 1

		if match := exportDeclRe.FindStringSubmatch(line); match != nil {
			exports = append(exports, analysis.Export{Name: match[4], Kind: match[3], IsDefault: match[1] != "", Line: lineNum})
		}

		if match := exportListRe.FindStringSubmatch(line); match != nil {
			exports = append(exports, parseExportList(match[1], lineNum, "named", "", false)...)
		}

		if match := exportDefaultRe.FindStringSubmatch(line); match != nil {
			exports = append(exports, analysis.Export{Name: match[1], Kind: "default", IsDefault: true, Line: lineNum})
		}

		if match := reExportRe.FindStringSubmatch(line); match != nil {
			exports = append(exports, parseExportList(match[1], lineNum, "reexport", match[2], true)...)
		}
	}

	return exports, nil
}

// GetFunctionBody returns the full body of a function
func (a *TypeScriptAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	return extractFunctionBody(content, funcName)
}

// GetExports for JavaScript
func (a *JavaScriptAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	// Reuse TypeScript implementation
	ts := NewTypeScriptAnalyzer()
	return ts.GetExports(ctx, filePath, content)
}

// GetFunctionBody for JavaScript
func (a *JavaScriptAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	return extractFunctionBody(content, funcName)
}

// GetExports for Vue
func (a *VueAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	text := string(content)
	if match := a.scriptRe.FindStringSubmatch(text); match != nil {
		ts := NewTypeScriptAnalyzer()
		return ts.GetExports(ctx, filePath, []byte(match[1]))
	}
	return nil, nil
}

// GetFunctionBody for Vue
func (a *VueAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	text := string(content)
	if match := a.scriptRe.FindStringSubmatch(text); match != nil {
		return extractFunctionBody([]byte(match[1]), funcName)
	}
	return "", 0, 0, nil
}

// findTypeEndLine finds the end line of a type alias (usually ends with semicolon or next line)
func findTypeEndLine(lines []string, startLineIdx int) int {
	if startLineIdx < 0 || startLineIdx >= len(lines) {
		return startLineIdx + 1
	}

	// Check if type is on single line
	if strings.Contains(lines[startLineIdx], ";") {
		return startLineIdx + 1
	}

	// For multi-line types, look for semicolon or closing brace
	braceCount := 0
	for i := startLineIdx; i < len(lines) && i < startLineIdx+50; i++ {
		line := lines[i]
		for _, ch := range line {
			if ch == '{' || ch == '<' {
				braceCount++
			} else if ch == '}' || ch == '>' {
				braceCount--
			}
		}
		if strings.Contains(line, ";") && braceCount <= 0 {
			return i + 1
		}
	}

	return startLineIdx + 1
}

// funcBodyPatterns returns compiled patterns for finding function starts
func funcBodyPatterns(funcName string) []*regexp.Regexp {
	quotedName := regexp.QuoteMeta(funcName)
	return []*regexp.Regexp{
		regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(async\s+)?function\s+` + quotedName + `\s*\(`),
		regexp.MustCompile(`(?m)^[\t ]*(export\s+)?(const|let|var)\s+` + quotedName + `\s*=\s*(async\s*)?\(`),
		regexp.MustCompile(`(?m)^[\t ]*` + quotedName + `\s*\([^)]*\)\s*\{`),
		regexp.MustCompile(`(?m)^[\t ]*` + quotedName + `\s*:\s*(async\s*)?\(`),
	}
}

// findFuncStartLine finds the starting line of a function
func findFuncStartLine(lines []string, patterns []*regexp.Regexp) int {
	for i, line := range lines {
		for _, p := range patterns {
			if p.MatchString(line) {
				return i
			}
		}
	}
	return -1
}

// findFuncEndLine finds the ending line by matching braces
func findFuncEndLine(lines []string, startLine int) int {
	braceCount := 0
	started := false

	for i := startLine; i < len(lines); i++ {
		for _, ch := range lines[i] {
			if ch == '{' {
				braceCount++
				started = true
			} else if ch == '}' {
				braceCount--
			}
		}
		if started && braceCount == 0 {
			return i
		}
	}
	return startLine
}

// extractFunctionBody is a helper to extract function body from JS/TS code
func extractFunctionBody(content []byte, funcName string) (string, int, int, error) {
	lines := strings.Split(string(content), "\n")
	patterns := funcBodyPatterns(funcName)

	startLine := findFuncStartLine(lines, patterns)
	if startLine < 0 {
		return "", 0, 0, nil
	}

	endLine := findFuncEndLine(lines, startLine)

	var body strings.Builder
	body.Grow((endLine - startLine + 1) * 80)
	for i := startLine; i <= endLine && i < len(lines); i++ {
		body.WriteString(lines[i])
		if i < endLine {
			body.WriteByte('\n')
		}
	}

	return body.String(), startLine + 1, endLine + 1, nil
}
