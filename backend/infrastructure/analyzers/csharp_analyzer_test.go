package analyzers

import (
	"context"
	"testing"
)

func TestCSharpAnalyzer_ExtractSymbols(t *testing.T) {
	analyzer := NewCSharpAnalyzer()
	ctx := context.Background()

	content := []byte(`
namespace MyApp
{
    public class MyClass
    {
        public void DoSomething()
        {
            Console.WriteLine("Hello");
        }
    }

    public interface IService
    {
        void Execute();
    }

    public struct Point
    {
        public int X;
        public int Y;
    }

    public enum Status
    {
        Active,
        Inactive
    }

    public record Person(string Name, int Age);
}
`)

	symbols, err := analyzer.ExtractSymbols(ctx, "test.cs", content)
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	counts := make(map[string]int)
	for _, sym := range symbols {
		counts[string(sym.Kind)]++
	}

	if counts["class"] < 1 {
		t.Errorf("Expected at least 1 class, got %d", counts["class"])
	}
	if counts["interface"] != 1 {
		t.Errorf("Expected 1 interface, got %d", counts["interface"])
	}
	if counts["struct"] != 1 {
		t.Errorf("Expected 1 struct, got %d", counts["struct"])
	}
	if counts["enum"] != 1 {
		t.Errorf("Expected 1 enum, got %d", counts["enum"])
	}
}

func TestCSharpAnalyzer_GetImports(t *testing.T) {
	analyzer := NewCSharpAnalyzer()
	ctx := context.Background()

	content := []byte(`
using System;
using System.Collections.Generic;
using Microsoft.Extensions.Logging;
using MyApp.Services;
using MyApp.Models;
`)

	imports, err := analyzer.GetImports(ctx, "test.cs", content)
	if err != nil {
		t.Fatalf("GetImports failed: %v", err)
	}

	if len(imports) != 5 {
		t.Errorf("Expected 5 imports, got %d", len(imports))
	}

	localCount := 0
	for _, imp := range imports {
		if imp.IsLocal {
			localCount++
		}
	}
	// System.* and Microsoft.* are external
	if localCount != 2 {
		t.Errorf("Expected 2 local imports (MyApp.*), got %d", localCount)
	}
}

func TestCSharpAnalyzer_CanAnalyze(t *testing.T) {
	analyzer := NewCSharpAnalyzer()

	tests := []struct {
		path     string
		expected bool
	}{
		{"Program.cs", true},
		{"MyClass.cs", true},
		{"test.go", false},
		{"test.java", false},
	}

	for _, tt := range tests {
		if got := analyzer.CanAnalyze(tt.path); got != tt.expected {
			t.Errorf("CanAnalyze(%s) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}

func TestCSharpAnalyzer_GetFunctionBody(t *testing.T) {
	analyzer := NewCSharpAnalyzer()
	ctx := context.Background()

	content := []byte(`
public class MyClass
{
    public void MyMethod()
    {
        Console.WriteLine("Hello");
        Console.WriteLine("World");
    }

    public void AnotherMethod()
    {
        // do nothing
    }
}
`)

	body, startLine, endLine, err := analyzer.GetFunctionBody(ctx, "test.cs", content, "MyMethod")
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
