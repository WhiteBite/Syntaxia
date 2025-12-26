package textutils

import (
	"strings"
	"testing"
)

func TestLicenseStripper_Strip_BlockComment(t *testing.T) {
	l := NewLicenseStripper()

	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name: "MIT license block comment",
			input: `/*
 * Copyright (c) 2024 Example Corp
 * MIT License
 * Permission is hereby granted...
 */

package main

func main() {}`,
			contains: "package main",
			excludes: "Copyright",
		},
		{
			name: "Apache license",
			input: `/*
 * Licensed under the Apache License, Version 2.0
 * You may not use this file except in compliance
 */
package main`,
			contains: "package main",
			excludes: "Apache License",
		},
		{
			name: "non-license comment preserved",
			input: `/*
 * This function does something important
 * It processes data efficiently
 */
package main`,
			contains: "This function does something",
			excludes: "",
		},
		{
			name: "HTML license comment",
			input: `<!--
  Copyright 2024 Example Corp
  All rights reserved
-->
<html>
<body>Hello</body>
</html>`,
			contains: "<html>",
			excludes: "Copyright",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := l.Strip(tt.input)
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("result should contain %q, got:\n%s", tt.contains, result)
			}
			if tt.excludes != "" && strings.Contains(result, tt.excludes) {
				t.Errorf("result should not contain %q, got:\n%s", tt.excludes, result)
			}
		})
	}
}

func TestLicenseStripper_Strip_LineComments(t *testing.T) {
	l := NewLicenseStripper()

	tests := []struct {
		name     string
		input    string
		contains string
		excludes string
	}{
		{
			name: "Python license header",
			input: `# Copyright (c) 2024 Example Corp
# Licensed under MIT License
# All rights reserved

def main():
    pass`,
			contains: "def main():",
			excludes: "Copyright",
		},
		{
			name: "Go line comments license",
			input: `// Copyright 2024 Example Corp
// SPDX-License-Identifier: MIT
// 
// Permission is hereby granted

package main`,
			contains: "package main",
			excludes: "Copyright",
		},
		{
			name: "SQL license",
			input: `-- Copyright 2024 Example Corp
-- Licensed under BSD License

SELECT * FROM users;`,
			contains: "SELECT * FROM users",
			excludes: "Copyright",
		},
		{
			name: "non-license line comments preserved",
			input: `# This script does something
# It processes files

def main():
    pass`,
			contains: "# This script does something",
			excludes: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := l.Strip(tt.input)
			if tt.contains != "" && !strings.Contains(result, tt.contains) {
				t.Errorf("result should contain %q, got:\n%s", tt.contains, result)
			}
			if tt.excludes != "" && strings.Contains(result, tt.excludes) {
				t.Errorf("result should not contain %q, got:\n%s", tt.excludes, result)
			}
		})
	}
}

func TestLicenseStripper_Strip_NoLicense(t *testing.T) {
	l := NewLicenseStripper()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no comments",
			input: "package main\n\nfunc main() {}",
		},
		{
			name:  "empty string",
			input: "",
		},
		{
			name:  "shebang preserved",
			input: "#!/bin/bash\necho hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := l.Strip(tt.input)
			if result != tt.input {
				t.Errorf("expected unchanged input, got:\n%s", result)
			}
		})
	}
}

func TestLicenseStripper_StripWithLanguageHint(t *testing.T) {
	l := NewLicenseStripper()

	input := `/*
 * Copyright 2024 Example
 * MIT License
 */
package main`

	result := l.StripWithLanguageHint(input, ".go")
	if strings.Contains(result, "Copyright") {
		t.Error("should strip license with language hint")
	}
	if !strings.Contains(result, "package main") {
		t.Error("should preserve code")
	}
}

func TestLicenseStripper_ContainsLicenseKeyword(t *testing.T) {
	l := NewLicenseStripper()

	tests := []struct {
		text     string
		expected bool
	}{
		{"Copyright 2024", true},
		{"MIT License", true},
		{"Apache License", true},
		{"All rights reserved", true},
		{"Permission is hereby granted", true},
		{"This is a normal comment", false},
		{"Function documentation", false},
		{"TODO: fix this", false},
		{"ЛИЦЕНЗИЯ MIT", true}, // Russian
	}

	for _, tt := range tests {
		result := l.containsLicenseKeyword(tt.text)
		if result != tt.expected {
			t.Errorf("containsLicenseKeyword(%q) = %v, want %v", tt.text, result, tt.expected)
		}
	}
}

func BenchmarkLicenseStripper_Strip(b *testing.B) {
	l := NewLicenseStripper()
	input := `/*
 * Copyright (c) 2024 Example Corporation
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * 
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 */

package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Strip(input)
	}
}

func BenchmarkLicenseStripper_NoLicense(b *testing.B) {
	l := NewLicenseStripper()
	input := `package main

import "fmt"

func main() {
    fmt.Println("Hello")
}
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Strip(input)
	}
}
