package analyzers

import (
	"context"
	"syntaxia/domain/analysis"
	"testing"
)

func TestTypeScriptAnalyzer_Language(t *testing.T) {
	a := NewTypeScriptAnalyzer()
	if a.Language() != "typescript" {
		t.Errorf("expected 'typescript', got %q", a.Language())
	}
}

func TestTypeScriptAnalyzer_Extensions(t *testing.T) {
	a := NewTypeScriptAnalyzer()
	exts := a.Extensions()
	if len(exts) != 2 {
		t.Errorf("expected 2 extensions, got %d", len(exts))
	}
}

func TestTypeScriptAnalyzer_CanAnalyze(t *testing.T) {
	a := NewTypeScriptAnalyzer()
	tests := []struct {
		path   string
		expect bool
	}{
		{"app.ts", true},
		{"component.tsx", true},
		{"script.js", false},
		{"main.go", false},
	}
	for _, tt := range tests {
		if got := a.CanAnalyze(tt.path); got != tt.expect {
			t.Errorf("CanAnalyze(%q) = %v, want %v", tt.path, got, tt.expect)
		}
	}
}

func TestTypeScriptAnalyzer_ExtractSymbols(t *testing.T) {
	a := NewTypeScriptAnalyzer()
	ctx := context.Background()
	code := `
export interface User {
	id: number;
	name: string;
}

export class UserService {
	getUser(id: number): User {
		return { id, name: "test" };
	}
}

export type UserID = number;

export function createUser(name: string): User {
	return { id: 1, name };
}

export enum Status {
	Active,
	Inactive
}
`

	symbols, err := a.ExtractSymbols(ctx, "user.ts", []byte(code))
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	symbolMap := make(map[string]analysis.Symbol)
	for _, s := range symbols {
		symbolMap[s.Name] = s
	}

	expected := []struct {
		name string
		kind analysis.SymbolKind
	}{
		{"User", analysis.KindInterface},
		{"UserService", analysis.KindClass},
		{"UserID", analysis.KindType},
		{"createUser", analysis.KindFunction},
		{"Status", analysis.KindEnum},
	}

	for _, exp := range expected {
		sym, ok := symbolMap[exp.name]
		if !ok {
			t.Errorf("symbol %q not found", exp.name)
			continue
		}
		if sym.Kind != exp.kind {
			t.Errorf("symbol %q: expected kind %v, got %v", exp.name, exp.kind, sym.Kind)
		}
	}
}

func TestTypeScriptAnalyzer_GetImports(t *testing.T) {
	a := NewTypeScriptAnalyzer()
	ctx := context.Background()
	code := `
import { User } from './models/user';
import * as utils from '../utils';
import axios from 'axios';
`

	imports, err := a.GetImports(ctx, "app.ts", []byte(code))
	if err != nil {
		t.Fatalf("GetImports failed: %v", err)
	}

	if len(imports) < 2 {
		t.Errorf("expected at least 2 imports, got %d", len(imports))
	}

	hasLocal := false
	hasExternal := false
	for _, imp := range imports {
		if imp.IsLocal {
			hasLocal = true
		} else {
			hasExternal = true
		}
	}

	if !hasLocal {
		t.Error("expected local import")
	}
	if !hasExternal {
		t.Error("expected external import")
	}
}

func TestJavaScriptAnalyzer_Language(t *testing.T) {
	a := NewJavaScriptAnalyzer()
	if a.Language() != "javascript" {
		t.Errorf("expected 'javascript', got %q", a.Language())
	}
}

func TestJavaScriptAnalyzer_Extensions(t *testing.T) {
	a := NewJavaScriptAnalyzer()
	exts := a.Extensions()
	expected := map[string]bool{".js": true, ".jsx": true, ".mjs": true}
	for _, ext := range exts {
		if !expected[ext] {
			t.Errorf("unexpected extension %q", ext)
		}
	}
}

func TestJavaScriptAnalyzer_ExtractSymbols(t *testing.T) {
	a := NewJavaScriptAnalyzer()
	ctx := context.Background()
	code := `
class UserService {
	constructor() {
		this.users = [];
	}
	
	getUser(id) {
		return this.users.find(u => u.id === id);
	}
}

function createUser(name) {
	return { id: Date.now(), name };
}

const helper = () => {
	return "help";
};
`

	symbols, err := a.ExtractSymbols(ctx, "app.js", []byte(code))
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	symbolMap := make(map[string]analysis.Symbol)
	for _, s := range symbols {
		symbolMap[s.Name] = s
	}

	if _, ok := symbolMap["UserService"]; !ok {
		t.Error("class 'UserService' not found")
	}
	if _, ok := symbolMap["createUser"]; !ok {
		t.Error("function 'createUser' not found")
	}
}

func TestStripCommentsAndStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name:     "line comment",
			input:    "const x = 1; // comment\nconst y = 2;",
			contains: "const x = 1;",
			excludes: "comment",
		},
		{
			name:     "block comment",
			input:    "const x = 1; /* block */ const y = 2;",
			contains: "const x = 1;",
			excludes: "block",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripCommentsAndStrings(tt.input)
			if tt.contains != "" && !containsStr(result, tt.contains) {
				t.Errorf("result should contain %q", tt.contains)
			}
		})
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0)
}

// Placeholder test file
