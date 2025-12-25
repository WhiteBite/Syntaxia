package analyzers

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"syntaxia/domain/analysis"
	"strings"
)

// GoAnalyzer analyzes Go source files using go/ast
type GoAnalyzer struct{}

// NewGoAnalyzer creates a new Go analyzer
func NewGoAnalyzer() *GoAnalyzer {
	return &GoAnalyzer{}
}

func (a *GoAnalyzer) Language() string     { return "go" }
func (a *GoAnalyzer) Extensions() []string { return []string{".go"} }
func (a *GoAnalyzer) CanAnalyze(filePath string) bool {
	return strings.HasSuffix(filePath, ".go")
}

func (a *GoAnalyzer) ExtractSymbols(ctx context.Context, filePath string, content []byte) ([]analysis.Symbol, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	symbols := make([]analysis.Symbol, 0, 32)

	// Extract package
	if file.Name != nil {
		symbols = append(symbols, analysis.Symbol{
			Name:      file.Name.Name,
			Kind:      analysis.KindPackage,
			Language:  "go",
			FilePath:  filePath,
			StartLine: fset.Position(file.Name.Pos()).Line,
		})
	}

	// Walk AST
	ast.Inspect(file, func(n ast.Node) bool {
		switch decl := n.(type) {
		case *ast.FuncDecl:
			sym := a.extractFunction(fset, filePath, decl)
			symbols = append(symbols, sym)
		case *ast.GenDecl:
			syms := a.extractGenDecl(fset, filePath, decl)
			symbols = append(symbols, syms...)
		}
		return true
	})

	return symbols, nil
}

func (a *GoAnalyzer) extractFunction(fset *token.FileSet, filePath string, decl *ast.FuncDecl) analysis.Symbol {
	kind := analysis.KindFunction
	var parent string

	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		kind = analysis.KindMethod
		if t, ok := decl.Recv.List[0].Type.(*ast.StarExpr); ok {
			if ident, ok := t.X.(*ast.Ident); ok {
				parent = ident.Name
			}
		} else if ident, ok := decl.Recv.List[0].Type.(*ast.Ident); ok {
			parent = ident.Name
		}
	}

	var sig strings.Builder
	sig.WriteString("func ")
	if parent != "" {
		sig.WriteString("(" + parent + ") ")
	}
	sig.WriteString(decl.Name.Name)
	sig.WriteString("(")
	if decl.Type.Params != nil {
		sig.WriteString(a.formatFieldList(decl.Type.Params))
	}
	sig.WriteString(")")
	if decl.Type.Results != nil {
		results := a.formatFieldList(decl.Type.Results)
		if len(decl.Type.Results.List) > 1 {
			sig.WriteString(" (" + results + ")")
		} else {
			sig.WriteString(" " + results)
		}
	}

	var doc string
	if decl.Doc != nil {
		doc = decl.Doc.Text()
	}

	startPos := fset.Position(decl.Pos())
	endPos := fset.Position(decl.End())

	return analysis.Symbol{
		Name:       decl.Name.Name,
		Kind:       kind,
		Language:   "go",
		FilePath:   filePath,
		StartLine:  startPos.Line,
		EndLine:    endPos.Line,
		StartCol:   startPos.Column,
		EndCol:     endPos.Column,
		Signature:  sig.String(),
		DocComment: doc,
		Parent:     parent,
	}
}

func (a *GoAnalyzer) extractGenDecl(fset *token.FileSet, filePath string, decl *ast.GenDecl) []analysis.Symbol {
	var symbols []analysis.Symbol
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.TypeSpec:
			sym := a.extractTypeSpec(fset, filePath, decl, s)
			symbols = append(symbols, sym)
		case *ast.ValueSpec:
			syms := a.extractValueSpec(fset, filePath, decl, s)
			symbols = append(symbols, syms...)
		}
	}
	return symbols
}

func (a *GoAnalyzer) extractTypeSpec(fset *token.FileSet, filePath string, decl *ast.GenDecl, spec *ast.TypeSpec) analysis.Symbol {
	kind := analysis.KindType
	switch spec.Type.(type) {
	case *ast.StructType:
		kind = analysis.KindStruct
	case *ast.InterfaceType:
		kind = analysis.KindInterface
	}

	var doc string
	if decl.Doc != nil {
		doc = decl.Doc.Text()
	}

	startPos := fset.Position(spec.Pos())
	endPos := fset.Position(spec.End())

	return analysis.Symbol{
		Name:       spec.Name.Name,
		Kind:       kind,
		Language:   "go",
		FilePath:   filePath,
		StartLine:  startPos.Line,
		EndLine:    endPos.Line,
		DocComment: doc,
	}
}

func (a *GoAnalyzer) extractValueSpec(fset *token.FileSet, filePath string, decl *ast.GenDecl, spec *ast.ValueSpec) []analysis.Symbol {
	symbols := make([]analysis.Symbol, 0, len(spec.Names))
	kind := analysis.KindVariable
	if decl.Tok == token.CONST {
		kind = analysis.KindConstant
	}
	for _, name := range spec.Names {
		startPos := fset.Position(name.Pos())
		symbols = append(symbols, analysis.Symbol{
			Name:      name.Name,
			Kind:      kind,
			Language:  "go",
			FilePath:  filePath,
			StartLine: startPos.Line,
		})
	}
	return symbols
}

func (a *GoAnalyzer) formatFieldList(fl *ast.FieldList) string {
	var parts []string
	for _, field := range fl.List {
		var names []string
		for _, name := range field.Names {
			names = append(names, name.Name)
		}
		typeStr := a.typeToString(field.Type)
		if len(names) > 0 {
			parts = append(parts, strings.Join(names, ", ")+" "+typeStr)
		} else {
			parts = append(parts, typeStr)
		}
	}
	return strings.Join(parts, ", ")
}

func (a *GoAnalyzer) typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + a.typeToString(t.X)
	case *ast.ArrayType:
		return "[]" + a.typeToString(t.Elt)
	case *ast.MapType:
		return "map[" + a.typeToString(t.Key) + "]" + a.typeToString(t.Value)
	case *ast.SelectorExpr:
		return a.typeToString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "any"
	}
}

func (a *GoAnalyzer) GetImports(ctx context.Context, filePath string, content []byte) ([]analysis.Import, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]analysis.Import, 0, len(file.Imports))
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		var alias string
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		imports = append(imports, analysis.Import{
			Path:    path,
			Alias:   alias,
			IsLocal: !strings.Contains(path, "."),
		})
	}
	return imports, nil
}

// GetExports returns exported symbols (public symbols in Go - capitalized names)
func (a *GoAnalyzer) GetExports(ctx context.Context, filePath string, content []byte) ([]analysis.Export, error) {
	symbols, err := a.ExtractSymbols(ctx, filePath, content)
	if err != nil {
		return nil, err
	}

	var exports []analysis.Export
	for _, sym := range symbols {
		// In Go, exported symbols start with uppercase
		if len(sym.Name) > 0 && sym.Name[0] >= 'A' && sym.Name[0] <= 'Z' {
			exports = append(exports, analysis.Export{
				Name: sym.Name,
				Kind: string(sym.Kind),
				Line: sym.StartLine,
			})
		}
	}
	return exports, nil
}

// GetFunctionBody returns the full body of a function by name
func (a *GoAnalyzer) GetFunctionBody(ctx context.Context, filePath string, content []byte, funcName string) (string, int, int, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, content, parser.ParseComments)
	if err != nil {
		return "", 0, 0, err
	}

	var result string
	var startLine, endLine int

	ast.Inspect(file, func(n ast.Node) bool {
		if decl, ok := n.(*ast.FuncDecl); ok {
			if decl.Name.Name == funcName {
				startPos := fset.Position(decl.Pos())
				endPos := fset.Position(decl.End())
				startLine = startPos.Line
				endLine = endPos.Line

				// Extract function body from content
				lines := strings.Split(string(content), "\n")
				var body strings.Builder
				for i := startLine - 1; i < endLine && i < len(lines); i++ {
					body.WriteString(lines[i])
					if i < endLine-1 {
						body.WriteString("\n")
					}
				}
				result = body.String()
				return false
			}
		}
		return true
	})

	return result, startLine, endLine, nil
}
