package testengine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"syntaxia/domain"
)

// TestNewTypeScriptTestAnalyzer tests analyzer creation
func TestNewTypeScriptTestAnalyzer(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.log)
}

// TestTypeScriptTestAnalyzer_AnalyzeTestDependencies tests dependency analysis
func TestTypeScriptTestAnalyzer_AnalyzeTestDependencies(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ts-deps-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests := []struct {
		name            string
		content         string
		expectedImports []string
	}{
		{
			name: "es6_imports",
			content: `import { Component } from 'react';
import axios from 'axios';
import * as utils from './utils';

describe('test', () => {});
`,
			expectedImports: []string{"react", "axios", "./utils"},
		},
		{
			name: "commonjs_require",
			content: `const fs = require('fs');
const path = require('path');

describe('test', () => {});
`,
			expectedImports: []string{"fs", "path"},
		},
		{
			name: "dynamic_imports",
			content: `const module = await import('./dynamic-module');
const lazy = import('lazy-module');

describe('test', () => {});
`,
			expectedImports: []string{"./dynamic-module", "lazy-module"},
		},
		{
			name: "mixed_imports",
			content: `import React from 'react';
const lodash = require('lodash');
const dynamic = await import('./dynamic');

describe('test', () => {});
`,
			expectedImports: []string{"react", "lodash", "./dynamic"},
		},
		{
			name:            "no_imports",
			content:         `describe('test', () => { it('works', () => {}); });`,
			expectedImports: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.name+".test.ts")
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			deps, err := analyzer.AnalyzeTestDependencies(context.Background(), testFile)
			require.NoError(t, err)

			for _, expected := range tt.expectedImports {
				assert.Contains(t, deps, expected, "Should contain import: %s", expected)
			}
		})
	}
}

// TestTypeScriptTestAnalyzer_AnalyzeTestDependencies_NonExistent tests with non-existent file
func TestTypeScriptTestAnalyzer_AnalyzeTestDependencies_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	_, err := analyzer.AnalyzeTestDependencies(context.Background(), "/nonexistent/file.test.ts")
	assert.Error(t, err)
}

// TestTypeScriptTestAnalyzer_FindTestsForFile tests finding tests for source file
func TestTypeScriptTestAnalyzer_FindTestsForFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ts-find-tests-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create source and test files
	files := map[string]string{
		"utils.ts":           `export const helper = () => {};`,
		"utils.test.ts":      `import { helper } from './utils'; describe('utils', () => {});`,
		"utils.spec.ts":      `import { helper } from './utils'; describe('utils spec', () => {});`,
		"component.tsx":      `export const Component = () => <div/>;`,
		"component.test.tsx": `import { Component } from './component'; describe('component', () => {});`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.ts", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
}

// TestTypeScriptTestAnalyzer_FindTestsForFile_InTestsDir tests finding tests in __tests__ directory
func TestTypeScriptTestAnalyzer_FindTestsForFile_InTestsDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ts-tests-dir-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create __tests__ directory
	testsDir := filepath.Join(tempDir, "__tests__")
	err = os.MkdirAll(testsDir, 0o755)
	require.NoError(t, err)

	files := map[string]string{
		"utils.ts":                `export const helper = () => {};`,
		"__tests__/utils.test.ts": `import { helper } from '../utils'; describe('utils', () => {});`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.MkdirAll(filepath.Dir(fullPath), 0o755)
		require.NoError(t, err)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests, err := analyzer.FindTestsForFile(context.Background(), "utils.ts", tempDir)
	require.NoError(t, err)

	assert.NotEmpty(t, tests)
}

// TestTypeScriptTestAnalyzer_IsSmokeTest tests smoke test detection
func TestTypeScriptTestAnalyzer_IsSmokeTest(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ts-smoke-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests := []struct {
		name     string
		fileName string
		content  string
		expected bool
	}{
		{
			name:     "smoke_in_content",
			fileName: "basic.test.ts",
			content:  `describe('smoke tests', () => { it('works', () => {}); });`,
			expected: true,
		},
		{
			name:     "smoke_comment",
			fileName: "api.test.ts",
			content:  `// smoke test\ndescribe('api', () => {});`,
			expected: true,
		},
		{
			name:     "smoke_in_filename",
			fileName: "smoke.test.ts",
			content:  `describe('test', () => {});`,
			expected: true,
		},
		{
			name:     "smoke_decorator",
			fileName: "feature.test.ts",
			content:  `@smoke\ndescribe('feature', () => {});`,
			expected: true,
		},
		{
			name:     "not_smoke_test",
			fileName: "unit.test.ts",
			content:  `describe('unit tests', () => { it('adds numbers', () => {}); });`,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, []byte(tt.content), 0o644)
			require.NoError(t, err)

			isSmoke, err := analyzer.IsSmokeTest(context.Background(), testFile)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, isSmoke)
		})
	}
}

// TestTypeScriptTestAnalyzer_IsSmokeTest_NonExistent tests smoke detection for non-existent file
func TestTypeScriptTestAnalyzer_IsSmokeTest_NonExistent(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	_, err := analyzer.IsSmokeTest(context.Background(), "/nonexistent/file.test.ts")
	assert.Error(t, err)
}

// TestTypeScriptTestAnalyzer_shouldSkipDir tests directory skip logic
func TestTypeScriptTestAnalyzer_shouldSkipDir(t *testing.T) {
	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests := []struct {
		name     string
		dirName  string
		expected bool
	}{
		{"node_modules", "node_modules", true},
		{"dist", "dist", true},
		{"build", "build", true},
		{"coverage", "coverage", true},
		{".git", ".git", true},
		{".next", ".next", true},
		{".nuxt", ".nuxt", true},
		{"out", "out", true},
		{".hidden", ".hidden", true},
		{"src", "src", false},
		{"tests", "tests", false},
		{"__tests__", "__tests__", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := analyzer.shouldSkipDir(tt.dirName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestTypeScriptTestAnalyzer_findTestsImportingFile tests finding tests that import a file
func TestTypeScriptTestAnalyzer_findTestsImportingFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "ts-import-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	files := map[string]string{
		"utils.ts":      `export const helper = () => {};`,
		"other.test.ts": `import { helper } from './utils'; describe('other', () => {});`,
		"unrelated.test.ts": `describe('unrelated', () => {});`,
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err = os.WriteFile(fullPath, []byte(content), 0o644)
		require.NoError(t, err)
	}

	log := &domain.NoopLogger{}
	analyzer := NewTypeScriptTestAnalyzer(log)

	tests, err := analyzer.findTestsImportingFile(context.Background(), "utils.ts", tempDir)
	require.NoError(t, err)

	// Should find other.test.ts but not unrelated.test.ts
	assert.NotEmpty(t, tests)
}
