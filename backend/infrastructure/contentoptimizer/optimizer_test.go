package contentoptimizer

import (
	"context"
	"strings"
	"syntaxia/domain"
	"testing"
)

type mockLogger struct{}

func (m *mockLogger) Debug(message string)   {}
func (m *mockLogger) Info(message string)    {}
func (m *mockLogger) Warning(message string) {}
func (m *mockLogger) Error(message string)   {}
func (m *mockLogger) Fatal(message string)   {}

func TestOptimizeGo(t *testing.T) {
	optimizer := NewContentOptimizer(&mockLogger{})
	ctx := context.Background()

	input := `package main

import (
	"fmt"
	"os"
	"strings"
)

// This is a comment
func main() {
	fmt.Println("Hello")
}
`

	opts := domain.ContentOptimizeOptions{
		StripComments:      true,
		CollapseEmptyLines: true,
	}

	result := optimizer.Optimize(ctx, input, "main.go", opts)

	if strings.Contains(result, "// This is a comment") {
		t.Error("Comments should be removed")
	}

	if strings.Contains(result, "imports collapsed") {
		t.Log("Imports were collapsed successfully")
	}
}

func TestOptimizeTypeScript(t *testing.T) {
	optimizer := NewContentOptimizer(&mockLogger{})
	ctx := context.Background()

	input := `import { Component } from 'vue'
import { ref } from 'vue'

// This is a comment
export interface MyInterface {
	name: string
}

export function myFunc() {
	return 42
}
`

	opts := domain.ContentOptimizeOptions{
		StripComments:      true,
		CollapseEmptyLines: true,
	}

	result := optimizer.Optimize(ctx, input, "component.ts", opts)

	if strings.Contains(result, "// This is a comment") {
		t.Error("Comments should be removed")
	}

	if strings.Contains(result, "imports collapsed") {
		t.Log("Imports were collapsed successfully")
	}
}

func TestOptimizePython(t *testing.T) {
	optimizer := NewContentOptimizer(&mockLogger{})
	ctx := context.Background()

	input := `import os
import sys
from typing import List

# This is a comment
def my_func(name: str) -> str:
	"""This is a docstring"""
	return name
`

	opts := domain.ContentOptimizeOptions{
		StripComments:      true,
		CollapseEmptyLines: true,
	}

	result := optimizer.Optimize(ctx, input, "script.py", opts)

	if strings.Contains(result, "# This is a comment") {
		t.Error("Comments should be removed")
	}

	if strings.Contains(result, "docstring") {
		t.Error("Docstrings should be removed")
	}
}

func TestCollapseEmptyLines(t *testing.T) {
	input := `line1


line2



line3`

	result := collapseEmptyLines(input)

	// Should have at most 2 consecutive newlines
	if strings.Contains(result, "\n\n\n") {
		t.Error("Should not have more than 2 consecutive newlines")
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"main.go", "go"},
		{"component.ts", "typescript"},
		{"script.js", "javascript"},
		{"app.py", "python"},
		{"Component.vue", "vue"},
		{"styles.css", "css"},
		{"unknown.xyz", "unknown"},
	}

	for _, tt := range tests {
		result := detectLanguage(tt.path)
		if result != tt.expected {
			t.Errorf("detectLanguage(%s) = %s, want %s", tt.path, result, tt.expected)
		}
	}
}

func TestCanGenerateSkeleton(t *testing.T) {
	optimizer := NewContentOptimizer(&mockLogger{})

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"component.ts", false},
		{"script.py", false},
	}

	for _, tt := range tests {
		result := optimizer.CanGenerateSkeleton(tt.path)
		if result != tt.expected {
			t.Errorf("CanGenerateSkeleton(%s) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}
