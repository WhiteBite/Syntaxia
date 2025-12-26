package textutils

import (
	"context"
	"fmt"
	"path/filepath"
	"syntaxia/domain/analysis"
	"sort"
	"strings"
)

// AnalyzerRegistry интерфейс для получения анализаторов
// Определен здесь чтобы избежать циклических зависимостей
type AnalyzerRegistry interface {
	GetAnalyzer(filePath string) analysis.LanguageAnalyzer
}

// SkeletonGenerator генерирует скелет кода из исходного файла
// Использует AST-анализаторы для извлечения структуры без тел функций
type SkeletonGenerator struct {
	registry AnalyzerRegistry
}

// NewSkeletonGenerator создает новый генератор скелетов
func NewSkeletonGenerator(registry AnalyzerRegistry) *SkeletonGenerator {
	return &SkeletonGenerator{
		registry: registry,
	}
}

// Generate генерирует скелет кода из содержимого файла
// Возвращает пустую строку если анализатор не найден или произошла ошибка
func (g *SkeletonGenerator) Generate(ctx context.Context, content, filePath string) string {
	if g.registry == nil {
		return ""
	}

	analyzer := g.registry.GetAnalyzer(filePath)
	if analyzer == nil {
		return ""
	}

	// Извлекаем символы
	symbols, err := analyzer.ExtractSymbols(ctx, filePath, []byte(content))
	if err != nil || len(symbols) == 0 {
		return ""
	}

	// Извлекаем импорты
	imports, _ := analyzer.GetImports(ctx, filePath, []byte(content))

	// Генерируем скелет в зависимости от языка
	lang := analyzer.Language()
	return g.generateSkeleton(lang, filePath, symbols, imports)
}

// generateSkeleton генерирует скелет для конкретного языка
func (g *SkeletonGenerator) generateSkeleton(lang, filePath string, symbols []analysis.Symbol, imports []analysis.Import) string {
	switch lang {
	case "go":
		return g.generateGoSkeleton(symbols, imports)
	case "typescript", "javascript":
		return g.generateTSSkeleton(symbols, imports)
	case "python":
		return g.generatePythonSkeleton(symbols, imports)
	case "java", "kotlin":
		return g.generateJavaSkeleton(symbols, imports)
	case "rust":
		return g.generateRustSkeleton(symbols, imports)
	case "csharp":
		return g.generateCSharpSkeleton(symbols, imports)
	default:
		return g.generateGenericSkeleton(symbols)
	}
}

// generateGoSkeleton генерирует скелет для Go
func (g *SkeletonGenerator) generateGoSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Package
	for _, sym := range symbols {
		if sym.Kind == analysis.KindPackage {
			b.WriteString("package ")
			b.WriteString(sym.Name)
			b.WriteString("\n\n")
			break
		}
	}

	// Imports (компактно)
	if len(imports) > 0 {
		b.WriteString("import (")
		for _, imp := range imports {
			b.WriteString(" \"")
			b.WriteString(imp.Path)
			b.WriteString("\"")
		}
		b.WriteString(" )\n\n")
	}

	// Группируем символы по типу
	types := g.filterSymbols(symbols, analysis.KindStruct, analysis.KindInterface, analysis.KindType)
	funcs := g.filterSymbols(symbols, analysis.KindFunction)
	methods := g.filterSymbols(symbols, analysis.KindMethod)
	consts := g.filterSymbols(symbols, analysis.KindConstant)
	vars := g.filterSymbols(symbols, analysis.KindVariable)

	// Constants
	if len(consts) > 0 {
		b.WriteString("const ( ")
		for _, c := range consts {
			b.WriteString(c.Name)
			b.WriteString("; ")
		}
		b.WriteString(")\n\n")
	}

	// Variables
	if len(vars) > 0 {
		b.WriteString("var ( ")
		for _, v := range vars {
			b.WriteString(v.Name)
			b.WriteString("; ")
		}
		b.WriteString(")\n\n")
	}

	// Types
	for _, t := range types {
		if t.Kind == analysis.KindStruct {
			b.WriteString("type ")
			b.WriteString(t.Name)
			b.WriteString(" struct { /* ... */ }\n")
		} else if t.Kind == analysis.KindInterface {
			b.WriteString("type ")
			b.WriteString(t.Name)
			b.WriteString(" interface { /* ... */ }\n")
		} else {
			b.WriteString("type ")
			b.WriteString(t.Name)
			b.WriteString(" /* ... */\n")
		}
	}
	if len(types) > 0 {
		b.WriteString("\n")
	}

	// Functions
	for _, f := range funcs {
		if f.Signature != "" {
			b.WriteString(f.Signature)
		} else {
			b.WriteString("func ")
			b.WriteString(f.Name)
			b.WriteString("()")
		}
		b.WriteString("\n")
	}

	// Methods (группируем по receiver)
	methodsByReceiver := make(map[string][]analysis.Symbol)
	for _, m := range methods {
		methodsByReceiver[m.Parent] = append(methodsByReceiver[m.Parent], m)
	}
	for receiver, ms := range methodsByReceiver {
		b.WriteString("// Methods for ")
		b.WriteString(receiver)
		b.WriteString(": ")
		for i, m := range ms {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(m.Name)
		}
		b.WriteString("\n")
	}

	return b.String()
}

// generateTSSkeleton генерирует скелет для TypeScript/JavaScript
func (g *SkeletonGenerator) generateTSSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Imports (компактно)
	if len(imports) > 0 {
		for _, imp := range imports {
			b.WriteString("import ")
			if len(imp.Names) > 0 {
				b.WriteString("{ ")
				b.WriteString(strings.Join(imp.Names, ", "))
				b.WriteString(" }")
			} else if imp.Alias != "" {
				b.WriteString(imp.Alias)
			} else {
				b.WriteString("...")
			}
			b.WriteString(" from '")
			b.WriteString(imp.Path)
			b.WriteString("'\n")
		}
		b.WriteString("\n")
	}

	// Interfaces and Types
	interfaces := g.filterSymbols(symbols, analysis.KindInterface, analysis.KindType)
	for _, i := range interfaces {
		b.WriteString("interface ")
		b.WriteString(i.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	// Classes
	classes := g.filterSymbols(symbols, analysis.KindClass)
	for _, c := range classes {
		b.WriteString("class ")
		b.WriteString(c.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	// Functions
	funcs := g.filterSymbols(symbols, analysis.KindFunction)
	for _, f := range funcs {
		if f.Signature != "" {
			b.WriteString(f.Signature)
		} else {
			b.WriteString("function ")
			b.WriteString(f.Name)
			b.WriteString("()")
		}
		b.WriteString("\n")
	}

	// Components (Vue/React)
	components := g.filterSymbols(symbols, analysis.KindComponent, analysis.KindComposable)
	for _, c := range components {
		b.WriteString("// Component: ")
		b.WriteString(c.Name)
		b.WriteString("\n")
	}

	return b.String()
}

// generatePythonSkeleton генерирует скелет для Python
func (g *SkeletonGenerator) generatePythonSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Imports
	if len(imports) > 0 {
		for _, imp := range imports {
			if len(imp.Names) > 0 {
				b.WriteString("from ")
				b.WriteString(imp.Path)
				b.WriteString(" import ")
				b.WriteString(strings.Join(imp.Names, ", "))
			} else {
				b.WriteString("import ")
				b.WriteString(imp.Path)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	// Classes
	classes := g.filterSymbols(symbols, analysis.KindClass)
	for _, c := range classes {
		b.WriteString("class ")
		b.WriteString(c.Name)
		b.WriteString(": ...\n")
	}

	// Functions
	funcs := g.filterSymbols(symbols, analysis.KindFunction)
	for _, f := range funcs {
		if f.Signature != "" {
			b.WriteString(f.Signature)
		} else {
			b.WriteString("def ")
			b.WriteString(f.Name)
			b.WriteString("(): ...")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// generateJavaSkeleton генерирует скелет для Java/Kotlin
func (g *SkeletonGenerator) generateJavaSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Package
	for _, sym := range symbols {
		if sym.Kind == analysis.KindPackage {
			b.WriteString("package ")
			b.WriteString(sym.Name)
			b.WriteString(";\n\n")
			break
		}
	}

	// Imports (компактно - только количество)
	if len(imports) > 0 {
		b.WriteString(fmt.Sprintf("// %d imports\n\n", len(imports)))
	}

	// Classes and Interfaces
	classes := g.filterSymbols(symbols, analysis.KindClass)
	interfaces := g.filterSymbols(symbols, analysis.KindInterface)

	for _, i := range interfaces {
		b.WriteString("interface ")
		b.WriteString(i.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	for _, c := range classes {
		b.WriteString("class ")
		b.WriteString(c.Name)
		b.WriteString(" {\n")

		// Methods для этого класса
		methods := g.filterSymbolsByParent(symbols, c.Name, analysis.KindMethod)
		for _, m := range methods {
			b.WriteString("  ")
			if m.Signature != "" {
				b.WriteString(m.Signature)
			} else {
				b.WriteString(m.Name)
				b.WriteString("()")
			}
			b.WriteString(";\n")
		}
		b.WriteString("}\n")
	}

	return b.String()
}

// generateRustSkeleton генерирует скелет для Rust
func (g *SkeletonGenerator) generateRustSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Use statements
	if len(imports) > 0 {
		for _, imp := range imports {
			b.WriteString("use ")
			b.WriteString(imp.Path)
			b.WriteString(";\n")
		}
		b.WriteString("\n")
	}

	// Structs
	structs := g.filterSymbols(symbols, analysis.KindStruct)
	for _, s := range structs {
		b.WriteString("struct ")
		b.WriteString(s.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	// Traits (interfaces)
	traits := g.filterSymbols(symbols, analysis.KindInterface)
	for _, t := range traits {
		b.WriteString("trait ")
		b.WriteString(t.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	// Functions
	funcs := g.filterSymbols(symbols, analysis.KindFunction)
	for _, f := range funcs {
		if f.Signature != "" {
			b.WriteString(f.Signature)
		} else {
			b.WriteString("fn ")
			b.WriteString(f.Name)
			b.WriteString("()")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// generateCSharpSkeleton генерирует скелет для C#
func (g *SkeletonGenerator) generateCSharpSkeleton(symbols []analysis.Symbol, imports []analysis.Import) string {
	var b strings.Builder

	// Usings
	if len(imports) > 0 {
		for _, imp := range imports {
			b.WriteString("using ")
			b.WriteString(imp.Path)
			b.WriteString(";\n")
		}
		b.WriteString("\n")
	}

	// Namespace
	for _, sym := range symbols {
		if sym.Kind == analysis.KindModule || sym.Kind == analysis.KindPackage {
			b.WriteString("namespace ")
			b.WriteString(sym.Name)
			b.WriteString(";\n\n")
			break
		}
	}

	// Interfaces
	interfaces := g.filterSymbols(symbols, analysis.KindInterface)
	for _, i := range interfaces {
		b.WriteString("interface ")
		b.WriteString(i.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	// Classes
	classes := g.filterSymbols(symbols, analysis.KindClass)
	for _, c := range classes {
		b.WriteString("class ")
		b.WriteString(c.Name)
		b.WriteString(" { /* ... */ }\n")
	}

	return b.String()
}

// generateGenericSkeleton генерирует универсальный скелет
func (g *SkeletonGenerator) generateGenericSkeleton(symbols []analysis.Symbol) string {
	var b strings.Builder

	// Группируем по типу
	byKind := make(map[analysis.SymbolKind][]analysis.Symbol)
	for _, sym := range symbols {
		byKind[sym.Kind] = append(byKind[sym.Kind], sym)
	}

	// Выводим по категориям
	kindOrder := []analysis.SymbolKind{
		analysis.KindPackage,
		analysis.KindInterface,
		analysis.KindClass,
		analysis.KindStruct,
		analysis.KindType,
		analysis.KindFunction,
		analysis.KindMethod,
		analysis.KindConstant,
		analysis.KindVariable,
	}

	for _, kind := range kindOrder {
		if syms, ok := byKind[kind]; ok && len(syms) > 0 {
			b.WriteString("// ")
			b.WriteString(string(kind))
			b.WriteString(": ")
			names := make([]string, len(syms))
			for i, s := range syms {
				names[i] = s.Name
			}
			b.WriteString(strings.Join(names, ", "))
			b.WriteString("\n")
		}
	}

	return b.String()
}

// filterSymbols фильтрует символы по типам
func (g *SkeletonGenerator) filterSymbols(symbols []analysis.Symbol, kinds ...analysis.SymbolKind) []analysis.Symbol {
	kindSet := make(map[analysis.SymbolKind]bool)
	for _, k := range kinds {
		kindSet[k] = true
	}

	var result []analysis.Symbol
	for _, sym := range symbols {
		if kindSet[sym.Kind] {
			result = append(result, sym)
		}
	}

	// Сортируем по имени для детерминированного вывода
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// filterSymbolsByParent фильтрует символы по родителю и типам
func (g *SkeletonGenerator) filterSymbolsByParent(symbols []analysis.Symbol, parent string, kinds ...analysis.SymbolKind) []analysis.Symbol {
	kindSet := make(map[analysis.SymbolKind]bool)
	for _, k := range kinds {
		kindSet[k] = true
	}

	var result []analysis.Symbol
	for _, sym := range symbols {
		if sym.Parent == parent && kindSet[sym.Kind] {
			result = append(result, sym)
		}
	}

	return result
}

// CanGenerateSkeleton проверяет, можно ли сгенерировать скелет для файла
func (g *SkeletonGenerator) CanGenerateSkeleton(filePath string) bool {
	if g.registry == nil {
		return false
	}
	return g.registry.GetAnalyzer(filePath) != nil
}

// SupportedExtensions возвращает поддерживаемые расширения
func SupportedSkeletonExtensions() []string {
	return []string{
		".go", ".ts", ".tsx", ".js", ".jsx",
		".py", ".java", ".kt", ".rs", ".cs",
		".vue", ".dart",
	}
}

// IsSkeletonSupported проверяет поддержку расширения
func IsSkeletonSupported(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	for _, supported := range SupportedSkeletonExtensions() {
		if ext == supported {
			return true
		}
	}
	return false
}
