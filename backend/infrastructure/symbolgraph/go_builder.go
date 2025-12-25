package symbolgraph

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"syntaxia/domain"
	"sort"
	"strings"
)

// GoSymbolGraphBuilder реализует SymbolGraphBuilder для Go
type GoSymbolGraphBuilder struct {
	log domain.Logger
}

// NewGoSymbolGraphBuilder создает новый builder для Go
func NewGoSymbolGraphBuilder(log domain.Logger) *GoSymbolGraphBuilder {
	return &GoSymbolGraphBuilder{
		log: log,
	}
}

// BuildGraph строит граф символов для Go проекта
func (b *GoSymbolGraphBuilder) BuildGraph(ctx context.Context, projectRoot string) (*domain.SymbolGraph, error) {
	b.log.Info(fmt.Sprintf("Building symbol graph for Go project: %s", projectRoot))

	graph := &domain.SymbolGraph{
		Nodes: []*domain.SymbolNode{},
		Edges: []*domain.SymbolEdge{},
	}

	// Проходим по всем .go файлам в проекте
	err := filepath.Walk(projectRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Пропускаем vendor, node_modules и другие служебные директории
			if info.Name() == "vendor" || info.Name() == "node_modules" ||
				info.Name() == ".git" || info.Name() == "dist" ||
				strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Парсим Go файл
		nodes, edges, err := b.parseGoFile(path, projectRoot)
		if err != nil {
			b.log.Warning(fmt.Sprintf("Failed to parse %s: %v", path, err))
			return nil // Продолжаем с другими файлами
		}

		graph.Nodes = append(graph.Nodes, nodes...)
		graph.Edges = append(graph.Edges, edges...)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to build symbol graph: %w", err)
	}

	// Обеспечиваем детерминизм: сортируем узлы и ребра
	b.sortGraphForDeterminism(graph)

	b.log.Info(fmt.Sprintf("Built symbol graph with %d nodes and %d edges", len(graph.Nodes), len(graph.Edges)))
	return graph, nil
}

// parseContext holds parsing context to avoid repeated allocations
type parseContext struct {
	fset        *token.FileSet
	packageName string
	relPath     string
	nodes       []*domain.SymbolNode
	edges       []*domain.SymbolEdge
}

// processImports extracts import nodes and edges
func (b *GoSymbolGraphBuilder) processImports(ctx *parseContext, file *ast.File, packageNodeID string) {
	for _, imp := range file.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		importNode := &domain.SymbolNode{
			ID:         fmt.Sprintf("import:%s:%s", ctx.relPath, importPath),
			Name:       importPath,
			Type:       domain.SymbolTypeImport,
			Path:       ctx.relPath,
			Package:    ctx.packageName,
			Visibility: domain.VisibilityPublic,
		}
		ctx.nodes = append(ctx.nodes, importNode)
		ctx.edges = append(ctx.edges, &domain.SymbolEdge{
			From: packageNodeID, To: importNode.ID, Type: domain.EdgeTypeImports, Weight: 1.0,
		})
	}
}

// processFuncDecl processes function declarations
func (b *GoSymbolGraphBuilder) processFuncDecl(ctx *parseContext, x *ast.FuncDecl) {
	funcNode := &domain.SymbolNode{
		ID:         fmt.Sprintf("func:%s:%s", ctx.relPath, x.Name.Name),
		Name:       x.Name.Name,
		Type:       domain.SymbolTypeFunction,
		Path:       ctx.relPath,
		Package:    ctx.packageName,
		Visibility: b.getVisibility(x.Name.Name),
	}
	if x.Recv != nil {
		funcNode.Type = domain.SymbolTypeMethod
		if len(x.Recv.List) > 0 {
			if receiverType := b.getReceiverType(x.Recv.List[0].Type); receiverType != "" {
				ctx.edges = append(ctx.edges, &domain.SymbolEdge{
					From: funcNode.ID, To: fmt.Sprintf("type:%s:%s", ctx.relPath, receiverType),
					Type: domain.EdgeTypeReferences, Weight: 1.0,
				})
			}
		}
	}
	ctx.nodes = append(ctx.nodes, funcNode)
}

// processStructFields processes struct field declarations
func (b *GoSymbolGraphBuilder) processStructFields(ctx *parseContext, structType *ast.StructType, typeName string) {
	for _, field := range structType.Fields.List {
		for _, name := range field.Names {
			fieldNode := &domain.SymbolNode{
				ID:         fmt.Sprintf("field:%s:%s.%s", ctx.relPath, typeName, name.Name),
				Name:       name.Name,
				Type:       domain.SymbolTypeField,
				Path:       ctx.relPath,
				Package:    ctx.packageName,
				Visibility: b.getVisibility(name.Name),
			}
			ctx.nodes = append(ctx.nodes, fieldNode)
			ctx.edges = append(ctx.edges, &domain.SymbolEdge{
				From: fieldNode.ID, To: fmt.Sprintf("type:%s:%s", ctx.relPath, typeName),
				Type: domain.EdgeTypeReferences, Weight: 1.0,
			})
		}
	}
}

// processTypeSpec processes type specifications
func (b *GoSymbolGraphBuilder) processTypeSpec(ctx *parseContext, typeSpec *ast.TypeSpec) {
	symbolType := domain.SymbolTypeType
	switch t := typeSpec.Type.(type) {
	case *ast.StructType:
		symbolType = domain.SymbolTypeStruct
		b.processStructFields(ctx, t, typeSpec.Name.Name)
	case *ast.InterfaceType:
		symbolType = domain.SymbolTypeInterface
	}
	ctx.nodes = append(ctx.nodes, &domain.SymbolNode{
		ID:         fmt.Sprintf("type:%s:%s", ctx.relPath, typeSpec.Name.Name),
		Name:       typeSpec.Name.Name,
		Type:       symbolType,
		Path:       ctx.relPath,
		Package:    ctx.packageName,
		Visibility: b.getVisibility(typeSpec.Name.Name),
	})
}

// processValueSpec processes var/const specifications
func (b *GoSymbolGraphBuilder) processValueSpec(ctx *parseContext, valueSpec *ast.ValueSpec, tok token.Token) {
	symbolType := domain.SymbolTypeVariable
	if tok == token.CONST {
		symbolType = domain.SymbolTypeConstant
	}
	for _, name := range valueSpec.Names {
		ctx.nodes = append(ctx.nodes, &domain.SymbolNode{
			ID:         fmt.Sprintf("%s:%s:%s", strings.ToLower(tok.String()), ctx.relPath, name.Name),
			Name:       name.Name,
			Type:       symbolType,
			Path:       ctx.relPath,
			Package:    ctx.packageName,
			Visibility: b.getVisibility(name.Name),
		})
	}
}

// parseGoFile парсит отдельный Go файл и извлекает символы
func (b *GoSymbolGraphBuilder) parseGoFile(filePath, projectRoot string) ([]*domain.SymbolNode, []*domain.SymbolEdge, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	relPath, _ := filepath.Rel(projectRoot, filePath)
	ctx := &parseContext{
		fset:        fset,
		packageName: file.Name.Name,
		relPath:     relPath,
		nodes:       make([]*domain.SymbolNode, 0, 32),
		edges:       make([]*domain.SymbolEdge, 0, 16),
	}

	// Package node
	packageNode := &domain.SymbolNode{
		ID: fmt.Sprintf("%s:%s", ctx.packageName, relPath), Name: ctx.packageName,
		Type: domain.SymbolTypePackage, Path: relPath, Package: ctx.packageName, Visibility: domain.VisibilityPublic,
	}
	ctx.nodes = append(ctx.nodes, packageNode)

	b.processImports(ctx, file, packageNode.ID)

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			b.processFuncDecl(ctx, x)
		case *ast.GenDecl:
			for _, spec := range x.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					b.processTypeSpec(ctx, s)
				case *ast.ValueSpec:
					b.processValueSpec(ctx, s, x.Tok)
				}
			}
		}
		return true
	})

	return ctx.nodes, ctx.edges, nil
}

// getVisibility определяет видимость символа по его имени
func (b *GoSymbolGraphBuilder) getVisibility(name string) domain.Visibility {
	if name == "" {
		return domain.VisibilityPrivate
	}

	// В Go: заглавная буква = публичный, строчная = приватный
	if name != "" && name[0] >= 'A' && name[0] <= 'Z' {
		return domain.VisibilityPublic
	}
	return domain.VisibilityPrivate
}

// getReceiverType извлекает тип получателя из AST
func (b *GoSymbolGraphBuilder) getReceiverType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name
		}
	}
	return ""
}

// UpdateGraph обновляет граф для измененных файлов
func (b *GoSymbolGraphBuilder) UpdateGraph(ctx context.Context, projectRoot string, changedFiles []string) (*domain.SymbolGraph, error) {
	// Для простоты перестраиваем весь граф
	// В будущем можно оптимизировать для инкрементальных обновлений
	return b.BuildGraph(ctx, projectRoot)
}

// GetSuggestions возвращает предложения символов на основе запроса
func (b *GoSymbolGraphBuilder) GetSuggestions(ctx context.Context, query string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	var suggestions []*domain.SymbolNode

	query = strings.ToLower(query)
	for _, node := range graph.Nodes {
		if strings.Contains(strings.ToLower(node.Name), query) ||
			strings.Contains(strings.ToLower(node.Package), query) {
			suggestions = append(suggestions, node)
		}
	}

	return suggestions, nil
}

// GetDependencies возвращает зависимости для указанного символа
func (b *GoSymbolGraphBuilder) GetDependencies(ctx context.Context, symbolID string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	var dependencies []*domain.SymbolNode
	dependencyIDs := make(map[string]bool)

	// Находим все ребра, где To = symbolID
	for _, edge := range graph.Edges {
		if edge.To == symbolID {
			dependencyIDs[edge.From] = true
		}
	}

	// Находим узлы по ID
	for _, node := range graph.Nodes {
		if dependencyIDs[node.ID] {
			dependencies = append(dependencies, node)
		}
	}

	return dependencies, nil
}

// GetDependents возвращает символы, которые зависят от указанного
func (b *GoSymbolGraphBuilder) GetDependents(ctx context.Context, symbolID string, graph *domain.SymbolGraph) ([]*domain.SymbolNode, error) {
	var dependents []*domain.SymbolNode
	dependentIDs := make(map[string]bool)

	// Находим все ребра, где From = symbolID
	for _, edge := range graph.Edges {
		if edge.From == symbolID {
			dependentIDs[edge.To] = true
		}
	}

	// Находим узлы по ID
	for _, node := range graph.Nodes {
		if dependentIDs[node.ID] {
			dependents = append(dependents, node)
		}
	}

	return dependents, nil
}

// sortGraphForDeterminism сортирует узлы и ребра для обеспечения детерминизма
func (b *GoSymbolGraphBuilder) sortGraphForDeterminism(graph *domain.SymbolGraph) {
	// Сортируем узлы по ID для детерминизма
	sort.Slice(graph.Nodes, func(i, j int) bool {
		return graph.Nodes[i].ID < graph.Nodes[j].ID
	})

	// Сортируем ребра по From, затем по To для детерминизма
	sort.Slice(graph.Edges, func(i, j int) bool {
		if graph.Edges[i].From != graph.Edges[j].From {
			return graph.Edges[i].From < graph.Edges[j].From
		}
		return graph.Edges[i].To < graph.Edges[j].To
	})
}
