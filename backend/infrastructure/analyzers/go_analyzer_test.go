package analyzers

import (
	"context"
	"syntaxia/domain/analysis"
	"testing"
)

func TestGoAnalyzer_Language(t *testing.T) {
	a := NewGoAnalyzer()
	if a.Language() != "go" {
		t.Errorf("expected 'go', got %q", a.Language())
	}
}

func TestGoAnalyzer_Extensions(t *testing.T) {
	a := NewGoAnalyzer()
	exts := a.Extensions()
	if len(exts) != 1 || exts[0] != ".go" {
		t.Errorf("expected ['.go'], got %v", exts)
	}
}

func TestGoAnalyzer_CanAnalyze(t *testing.T) {
	a := NewGoAnalyzer()
	tests := []struct {
		path   string
		expect bool
	}{
		{"main.go", true},
		{"service_test.go", true},
		{"app.ts", false},
	}
	for _, tt := range tests {
		if got := a.CanAnalyze(tt.path); got != tt.expect {
			t.Errorf("CanAnalyze(%q) = %v, want %v", tt.path, got, tt.expect)
		}
	}
}

func TestGoAnalyzer_ExtractSymbols(t *testing.T) {
	a := NewGoAnalyzer()
	ctx := context.Background()
	code := "package main\n\ntype User struct {\n\tID int\n}\n\nfunc NewUser() *User {\n\treturn &User{}\n}\n\nfunc (u *User) GetID() int {\n\treturn u.ID\n}\n"

	symbols, err := a.ExtractSymbols(ctx, "main.go", []byte(code))
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	symbolMap := make(map[string]analysis.Symbol)
	for _, s := range symbols {
		symbolMap[s.Name] = s
	}

	if _, ok := symbolMap["main"]; !ok {
		t.Error("package 'main' not found")
	}
	if _, ok := symbolMap["User"]; !ok {
		t.Error("struct 'User' not found")
	}
	if _, ok := symbolMap["NewUser"]; !ok {
		t.Error("function 'NewUser' not found")
	}
	if sym, ok := symbolMap["GetID"]; ok {
		if sym.Kind != analysis.KindMethod {
			t.Errorf("GetID should be method, got %v", sym.Kind)
		}
		if sym.Parent != "User" {
			t.Errorf("GetID parent should be 'User', got %q", sym.Parent)
		}
	} else {
		t.Error("method 'GetID' not found")
	}
}

func TestGoAnalyzer_GetImports(t *testing.T) {
	a := NewGoAnalyzer()
	ctx := context.Background()
	code := "package main\n\nimport (\n\t\"fmt\"\n\t\"github.com/gin-gonic/gin\"\n\tmylog \"github.com/sirupsen/logrus\"\n)\n"

	imports, err := a.GetImports(ctx, "main.go", []byte(code))
	if err != nil {
		t.Fatalf("GetImports failed: %v", err)
	}

	if len(imports) != 3 {
		t.Errorf("expected 3 imports, got %d", len(imports))
	}

	importMap := make(map[string]analysis.Import)
	for _, imp := range imports {
		importMap[imp.Path] = imp
	}

	if imp, ok := importMap["fmt"]; ok {
		if !imp.IsLocal {
			t.Error("fmt should be local")
		}
	}
	if imp, ok := importMap["github.com/sirupsen/logrus"]; ok {
		if imp.Alias != "mylog" {
			t.Errorf("logrus alias should be 'mylog', got %q", imp.Alias)
		}
	}
}

func TestGoAnalyzer_GetExports(t *testing.T) {
	a := NewGoAnalyzer()
	ctx := context.Background()
	code := "package main\n\ntype User struct{}\ntype privateType struct{}\n\nfunc PublicFunc() {}\nfunc privateFunc() {}\n"

	exports, err := a.GetExports(ctx, "main.go", []byte(code))
	if err != nil {
		t.Fatalf("GetExports failed: %v", err)
	}

	exportNames := make(map[string]bool)
	for _, exp := range exports {
		exportNames[exp.Name] = true
	}

	if !exportNames["User"] {
		t.Error("User should be exported")
	}
	if !exportNames["PublicFunc"] {
		t.Error("PublicFunc should be exported")
	}
	if exportNames["privateType"] {
		t.Error("privateType should not be exported")
	}
}

func TestGoAnalyzer_GetFunctionBody(t *testing.T) {
	a := NewGoAnalyzer()
	ctx := context.Background()
	code := "package main\n\nfunc Hello() string {\n\treturn \"hello\"\n}\n"

	body, startLine, endLine, err := a.GetFunctionBody(ctx, "main.go", []byte(code), "Hello")
	if err != nil {
		t.Fatalf("GetFunctionBody failed: %v", err)
	}

	if body == "" {
		t.Error("expected function body")
	}
	if startLine != 3 {
		t.Errorf("expected startLine 3, got %d", startLine)
	}
	if endLine != 5 {
		t.Errorf("expected endLine 5, got %d", endLine)
	}
}
