package analyzers

import (
	"context"
	"testing"
)

func TestRustAnalyzer_ExtractSymbols(t *testing.T) {
	analyzer := NewRustAnalyzer()
	ctx := context.Background()

	content := []byte(`
pub struct MyStruct {
    name: String,
}

enum Status {
    Active,
    Inactive,
}

pub trait Drawable {
    fn draw(&self);
}

impl Drawable for MyStruct {
    fn draw(&self) {
        println!("{}", self.name);
    }
}

pub fn public_function() -> i32 {
    42
}

async fn async_function() {
    // async code
}

const MAX_SIZE: usize = 100;

type Result<T> = std::result::Result<T, Error>;
`)

	symbols, err := analyzer.ExtractSymbols(ctx, "test.rs", content)
	if err != nil {
		t.Fatalf("ExtractSymbols failed: %v", err)
	}

	counts := make(map[string]int)
	for _, sym := range symbols {
		counts[string(sym.Kind)]++
	}

	if counts["struct"] != 1 {
		t.Errorf("Expected 1 struct, got %d", counts["struct"])
	}
	if counts["enum"] != 1 {
		t.Errorf("Expected 1 enum, got %d", counts["enum"])
	}
	if counts["interface"] != 1 { // trait
		t.Errorf("Expected 1 trait (interface), got %d", counts["interface"])
	}
	if counts["function"] < 2 {
		t.Errorf("Expected at least 2 functions, got %d", counts["function"])
	}
	if counts["constant"] != 1 {
		t.Errorf("Expected 1 constant, got %d", counts["constant"])
	}
}

func TestRustAnalyzer_GetImports(t *testing.T) {
	analyzer := NewRustAnalyzer()
	ctx := context.Background()

	content := []byte(`
use std::collections::HashMap;
use crate::module::something;
use self::local::item;
use super::parent::thing;
use external_crate::Type;
`)

	imports, err := analyzer.GetImports(ctx, "test.rs", content)
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
	if localCount != 3 { // crate::, self::, super::
		t.Errorf("Expected 3 local imports, got %d", localCount)
	}
}

func TestRustAnalyzer_GetExports(t *testing.T) {
	analyzer := NewRustAnalyzer()
	ctx := context.Background()

	content := []byte(`
pub struct PublicStruct {}
struct PrivateStruct {}
pub fn public_fn() {}
fn private_fn() {}
`)

	exports, err := analyzer.GetExports(ctx, "test.rs", content)
	if err != nil {
		t.Fatalf("GetExports failed: %v", err)
	}

	if len(exports) != 2 {
		t.Errorf("Expected 2 exports (pub items), got %d", len(exports))
	}
}

func TestRustAnalyzer_CanAnalyze(t *testing.T) {
	analyzer := NewRustAnalyzer()

	tests := []struct {
		path     string
		expected bool
	}{
		{"main.rs", true},
		{"lib.rs", true},
		{"test.go", false},
		{"test.py", false},
	}

	for _, tt := range tests {
		if got := analyzer.CanAnalyze(tt.path); got != tt.expected {
			t.Errorf("CanAnalyze(%s) = %v, want %v", tt.path, got, tt.expected)
		}
	}
}
