package analyzers

import (
	"context"
	"testing"
)

func TestPythonAnalyzer_ExtractSymbols(t *testing.T) {
	analyzer := NewPythonAnalyzer()
	ctx := context.Background()

	content := []byte(`
class MyClass:
    def __init__(self, name):
        self.name = name
    
    def greet(self):
        return f"Hello, {self.name}"

def standalone_function(x, y):
    return x + y

class AnotherClass(BaseClass):
    pass
`)

	symbols, err := analyzer.ExtractSymbols(ctx, "test.py", content)
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	// Should find 2 classes, 1 standalone function, 2 methods
	classCount := 0
	funcCount := 0
	methodCount := 0
	for _, sym := range symbols {
		switch sym.Kind {
		case "class":
			classCount++
		case "function":
			funcCount++
		case "method":
			methodCount++
		}
	}

	if classCount != 2 {
		t.Errorf("Expected 2 classes, got %d", classCount)
	}
	if funcCount != 1 {
		t.Errorf("Expected 1 function, got %d", funcCount)
	}
	if methodCount != 2 {
		t.Errorf("Expected 2 methods, got %d", methodCount)
	}
}

func TestPythonAnalyzer_GetImports(t *testing.T) {
	analyzer := NewPythonAnalyzer()
	ctx := context.Background()

	content := []byte(`
import os
import sys
from typing import List, Optional
from .local_module import something
from ..parent import other
`)

	imports, err := analyzer.GetImports(ctx, "test.py", content)
	if err != nil {
		t.Fatalf("GetImports failed: %v", err)
	}

	if len(imports) != 5 {
		t.Errorf("Expected 5 imports, got %d", len(imports))
	}

	// Check local imports (relative imports starting with .)
	localCount := 0
	for _, imp := range imports {
		if imp.IsLocal {
			localCount++
		}
	}
	// .local_module and ..parent are local (relative imports)
	if localCount < 2 {
		t.Errorf("Expected at least 2 local imports, got %d", localCount)
	}
}

func TestPythonAnalyzer_GetFunctionBody(t *testing.T) {
	analyzer := NewPythonAnalyzer()
	ctx := context.Background()

	content := []byte(`
def my_function(x):
    result = x * 2
    return result

def another():
    pass
`)

	body, startLine, endLine, err := analyzer.GetFunctionBody(ctx, "test.py", content, "my_function")
	if err != nil {
		t.Fatalf("GetFunctionBody failed: %v", err)
	}

	if startLine == 0 {
		t.Error("Expected non-zero start line")
	}
	if endLine <= startLine {
		t.Error("Expected end line > start line")
	}
	if body == "" {
		t.Error("Expected non-empty body")
	}
}

func TestPythonAnalyzer_CanAnalyze(t *testing.T) {
	analyzer := NewPythonAnalyzer()

	tests := []struct {
		path     string
		expected bool
	}{
		{"test.py", true},
		{"module.pyw", true},
		{"types.pyi", true},
		{"test.go", false},
		{"test.js", false},
	}

	for _, tt := range tests {
		if got := analyzer.CanAnalyze(tt.path); got != tt.expected {
			t.Errorf("CanAnalyze(%s) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}
