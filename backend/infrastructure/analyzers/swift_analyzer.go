package analyzers

import (
	"context"
	"regexp"
	"strings"

	"syntaxia/domain/analysis"
)

// SwiftAnalyzer analyzes Swift files
type SwiftAnalyzer struct {
	classRe     *regexp.Regexp
	structRe    *regexp.Regexp
	enumRe      *regexp.Regexp
	protocolRe  *regexp.Regexp
	funcRe      *regexp.Regexp
	extensionRe *regexp.Regexp
	importRe    *regexp.Regexp
	typealiasRe *regexp.Regexp
}

func NewSwiftAnalyzer() *SwiftAnalyzer {
	return &SwiftAnalyzer{
		classRe:     regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+|fileprivate\s+|open\s+)?(final\s+)?class\s+(\w+)`),
		structRe:    regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+|fileprivate\s+)?struct\s+(\w+)`),
		enumRe:      regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+|fileprivate\s+)?enum\s+(\w+)`),
		protocolRe:  regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+)?protocol\s+(\w+)`),
		funcRe:      regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+|fileprivate\s+|open\s+)?(static\s+|class\s+)?(override\s+)?func\s+(\w+)`),
		extensionRe: regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+)?extension\s+(\w+)`),
		importRe:    regexp.MustCompile(`(?m)^[\t ]*import\s+(\w+)`),
		typealiasRe: regexp.MustCompile(`(?m)^[\t ]*(public\s+|private\s+|internal\s+)?typealias\s+(\w+)`),
	}
}

func (a *SwiftAnalyzer) Language() string     { return "swift" }
func (a *SwiftAnalyzer) Extensions() []string { return []string{".swift"} }

func (a *SwiftAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(strings.ToLower(filePath), ".swift")
}


func (a *SwiftAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	text := string(content)
	lines := strings.Split(text, "\n")
	symbols := make([]analysis.Symbol, 0)

	// Classes
	for _, match := range a.classRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[6]:match[7]]
		endLine := findBlockEndLine(lines, line-1)
		vis := extractSwiftVisibility(text, match)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindClass, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: endLine,
			Extra: map[string]string{"visibility": vis},
		})
	}

	// Structs
	for _, match := range a.structRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		vis := extractSwiftVisibility(text, match)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindStruct, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: endLine,
			Extra: map[string]string{"visibility": vis},
		})
	}

	// Enums
	for _, match := range a.enumRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindEnum, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: endLine,
		})
	}

	// Protocols
	for _, match := range a.protocolRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		endLine := findBlockEndLine(lines, line-1)
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindInterface, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: endLine,
		})
	}

	// Functions
	for _, match := range a.funcRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[8]:match[9]]
		endLine := findBlockEndLine(lines, line-1)
		vis := extractSwiftVisibility(text, match)
		isStatic := match[4] != -1 && match[5] != -1
		extra := map[string]string{"visibility": vis}
		if isStatic {
			extra["static"] = "true"
		}
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindFunction, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: endLine, Extra: extra,
		})
	}

	// Type aliases
	for _, match := range a.typealiasRe.FindAllStringSubmatchIndex(text, -1) {
		line := countLines(content[:match[0]])
		name := text[match[4]:match[5]]
		symbols = append(symbols, analysis.Symbol{
			Name: name, Kind: analysis.KindType, Language: "swift",
			FilePath: filePath, StartLine: line, EndLine: line,
		})
	}

	return symbols, nil
}

func extractSwiftVisibility(text string, match []int) string {
	if match[2] != -1 && match[3] != -1 {
		return strings.TrimSpace(text[match[2]:match[3]])
	}
	return "internal"
}

func (a *SwiftAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	matches := a.importRe.FindAllStringSubmatch(string(content), -1)
	imports := make([]analysis.Import, 0, len(matches))
	for _, match := range matches {
		path := match[1]
		isLocal := !isSwiftFramework(path)
		imports = append(imports, analysis.Import{Path: path, IsLocal: isLocal})
	}
	return imports, nil
}

func isSwiftFramework(name string) bool {
	frameworks := map[string]bool{
		"Foundation": true, "UIKit": true, "SwiftUI": true, "Combine": true,
		"CoreData": true, "CoreGraphics": true, "CoreLocation": true,
		"MapKit": true, "AVFoundation": true, "Photos": true, "WebKit": true,
		"StoreKit": true, "CloudKit": true, "HealthKit": true, "HomeKit": true,
		"WatchKit": true, "AppKit": true, "Cocoa": true, "Darwin": true,
		"Dispatch": true, "ObjectiveC": true, "Swift": true, "XCTest": true,
	}
	return frameworks[name]
}

func (a *SwiftAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}
	exports := make([]analysis.Export, 0)
	for _, sym := range symbols {
		if sym.Extra != nil {
			vis := sym.Extra["visibility"]
			if vis == "public" || vis == "open" {
				exports = append(exports, analysis.Export{Name: sym.Name, Kind: string(sym.Kind), Line: sym.StartLine})
			}
		}
	}
	return exports, nil
}

func (a *SwiftAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	lines := strings.Split(string(content), "\n")
	funcRe := regexp.MustCompile(`(?m)func\s+` + regexp.QuoteMeta(funcName) + `\s*[(<]`)
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
