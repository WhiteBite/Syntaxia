import { NoiseReductionService } from '@/services/noiseReduction.service'
import { describe, expect, it } from 'vitest'

describe('NoiseReductionService', () => {
    const service = new NoiseReductionService()

    describe('getLanguageFromExtension', () => {
        it('should detect TypeScript files', () => {
            expect(service.getLanguageFromExtension('file.ts')).toBe('typescript')
            expect(service.getLanguageFromExtension('component.tsx')).toBe('tsx')
        })

        it('should detect JavaScript files', () => {
            expect(service.getLanguageFromExtension('file.js')).toBe('javascript')
            expect(service.getLanguageFromExtension('component.jsx')).toBe('jsx')
        })

        it('should detect Python files', () => {
            expect(service.getLanguageFromExtension('script.py')).toBe('python')
        })

        it('should detect Go files', () => {
            expect(service.getLanguageFromExtension('main.go')).toBe('go')
        })

        it('should detect Vue files', () => {
            expect(service.getLanguageFromExtension('Component.vue')).toBe('vue')
        })

        it('should detect CSS files', () => {
            expect(service.getLanguageFromExtension('styles.css')).toBe('css')
            expect(service.getLanguageFromExtension('styles.scss')).toBe('scss')
        })

        it('should return unknown for unsupported extensions', () => {
            expect(service.getLanguageFromExtension('file.xyz')).toBe('unknown')
        })
    })

    describe('cleanupFileContent - TypeScript/JavaScript', () => {
        it('should collapse import statements', () => {
            const content = `import { foo } from 'foo'
import { bar } from 'bar'
import { baz } from 'baz'

export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 3 imports collapsed')
            expect(result.content).toContain('export function test()')
            expect(result.content).not.toContain('import { foo }')
            expect(result.savings).toBeGreaterThan(0)
        })

        it('should remove single-line comments', () => {
            const content = `// This is a comment
export function test() {
    // Another comment
    return 'hello' // Inline comment
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('// This is a comment')
            expect(result.content).not.toContain('// Another comment')
            expect(result.content).toContain('export function test()')
        })

        it('should remove multi-line comments', () => {
            const content = `/* This is a
   multi-line comment */
export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('multi-line comment')
            expect(result.content).toContain('export function test()')
        })

        it('should remove type definitions', () => {
            const content = `export interface User {
    name: string
    age: number
}

export type Status = 'active' | 'inactive'

export function getUser(): User {
    return { name: 'John', age: 30 }
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: false,
                removeComments: false,
                removeTypeDefinitions: true,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('export interface User')
            expect(result.content).not.toContain('export type Status')
            expect(result.content).toContain('export function getUser()')
        })

        it('should remove "use strict"', () => {
            const content = `"use strict";

export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'javascript', {
                collapseImports: false,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: true
            })

            expect(result.content).not.toContain('"use strict"')
            expect(result.content).toContain('export function test()')
        })

        it('should calculate token savings correctly', () => {
            const content = `import { foo } from 'foo'
import { bar } from 'bar'
import { baz } from 'baz'
import { qux } from 'qux'
import { quux } from 'quux'

// This is a long comment that takes up space
// Another comment
// Yet another comment

export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.savings).toBeGreaterThan(0)
            expect(result.savingsPercent).toBeGreaterThan(0)
            expect(result.cleanedTokens).toBeLessThan(result.originalTokens)
        })
    })

    describe('cleanupFileContent - Python', () => {
        it('should collapse Python imports', () => {
            const content = `import os
import sys
from typing import List, Dict
from pathlib import Path

def test():
    return "hello"`

            const result = service.cleanupFileContent(content, 'python', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 4 imports collapsed')
            expect(result.content).toContain('def test()')
            expect(result.content).not.toContain('import os')
        })

        it('should remove Python comments and docstrings', () => {
            const content = `# This is a comment
def test():
    """This is a docstring"""
    # Another comment
    return "hello"`

            const result = service.cleanupFileContent(content, 'python', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('# This is a comment')
            expect(result.content).not.toContain('This is a docstring')
            expect(result.content).toContain('def test()')
        })

        it('should remove Python type hints', () => {
            const content = `def greet(name: str) -> str:
    result: str = "Hello"
    return result`

            const result = service.cleanupFileContent(content, 'python', {
                collapseImports: false,
                removeComments: false,
                removeTypeDefinitions: true,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('def greet(name)')
            expect(result.content).not.toContain(': str')
            expect(result.content).not.toContain('-> str')
        })
    })

    describe('cleanupFileContent - Go', () => {
        it('should collapse Go import blocks', () => {
            const content = `package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    fmt.Println("hello")
}`

            const result = service.cleanupFileContent(content, 'go', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 3 imports collapsed')
            expect(result.content).toContain('func main()')
            expect(result.content).not.toContain('"fmt"')
        })

        it('should remove Go comments', () => {
            const content = `package main

// This is a comment
func main() {
    /* Multi-line
       comment */
    println("hello")
}`

            const result = service.cleanupFileContent(content, 'go', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('// This is a comment')
            expect(result.content).not.toContain('Multi-line')
            expect(result.content).toContain('func main()')
        })
    })

    describe('cleanupFileContent - Vue', () => {
        it('should clean script section in Vue files', () => {
            const content = `<template>
  <div>{{ message }}</div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'pinia'

export default {
  setup() {
    const message = ref('Hello')
    return { message }
  }
}
</script>

<style scoped>
div { color: red; }
</style>`

            const result = service.cleanupFileContent(content, 'vue', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('<template>')
            expect(result.content).toContain('// 2 imports collapsed')
            expect(result.content).not.toContain('import { ref }')
            expect(result.content).toContain('<style scoped>')
        })

        it('should remove HTML comments from Vue template', () => {
            const content = `<template>
  <!-- This is a comment -->
  <div>{{ message }}</div>
</template>

<script>
export default {}
</script>`

            const result = service.cleanupFileContent(content, 'vue', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('<!-- This is a comment -->')
            expect(result.content).toContain('<div>{{ message }}</div>')
        })
    })

    describe('cleanupFileContent - CSS', () => {
        it('should collapse CSS @import statements', () => {
            const content = `@import url('reset.css');
@import url('typography.css');
@import url('layout.css');

.container {
    display: flex;
}`

            const result = service.cleanupFileContent(content, 'css', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 3 imports collapsed')
            expect(result.content).toContain('.container')
            expect(result.content).not.toContain('@import')
        })

        it('should remove CSS comments', () => {
            const content = `/* Main styles */
.container {
    /* Flexbox layout */
    display: flex;
}`

            const result = service.cleanupFileContent(content, 'css', {
                collapseImports: false,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).not.toContain('/* Main styles */')
            expect(result.content).not.toContain('/* Flexbox layout */')
            expect(result.content).toContain('.container')
        })
    })

    describe('whitespace normalization', () => {
        it('should collapse excessive blank lines', () => {
            const content = `function test() {



    return 'hello'
}


function another() {
    return 'world'
}`

            const result = service.cleanupFileContent(content, 'javascript', {
                collapseImports: false,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: true
            })

            // Should not have 3+ consecutive newlines
            expect(result.content).not.toMatch(/\n{3,}/)
            expect(result.content).toContain('function test()')
            expect(result.content).toContain('function another()')
        })
    })

    describe('error handling', () => {
        it('should return original content on error', () => {
            const content = 'valid content'

            // Test with unknown language - should still work
            const result = service.cleanupFileContent(content, 'unknown', {
                collapseImports: true,
                removeComments: true,
                removeTypeDefinitions: true,
                collapseBoilerplate: true
            })

            expect(result.content).toBe(content.trim())
            expect(result.savings).toBe(0)
        })
    })

    describe('combined options', () => {
        it('should apply collapse imports and remove comments', () => {
            const content = `import { foo } from 'foo'
import { bar } from 'bar'

// This is a comment
export function getUser() {
    // Implementation
    return { name: 'John' }
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 2 imports collapsed')
            expect(result.content).not.toContain('// This is a comment')
            expect(result.content).not.toContain('// Implementation')
            expect(result.content).toContain('export function getUser')
            expect(result.savings).toBeGreaterThan(0)
        })

        it('should handle imports with empty lines between them', () => {
            const content = `import { foo } from 'foo'

import { bar } from 'bar'

import { baz } from 'baz'

export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            expect(result.content).toContain('// 3 imports collapsed')
            expect(result.content).toContain('export function test()')
            expect(result.content).not.toContain('import { foo }')
        })

        it('should handle mixed imports and code', () => {
            const content = `import { foo } from 'foo'
import { bar } from 'bar'

const config = { test: true }

import { baz } from 'baz'

export function test() {
    return 'hello'
}`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: false,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            // Should collapse first block of 2 imports
            expect(result.content).toContain('// 2 imports collapsed')
            // Should collapse second import separately
            expect(result.content).toContain('// 1 imports collapsed')
            expect(result.content).toContain('const config')
            expect(result.content).toContain('export function test()')
        })
    })

    describe('real-world scenarios', () => {
        it('should handle complex TypeScript file with imports and comments', () => {
            const content = `import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Settings } from '@/types'

// Store definition
export const useUserStore = defineStore('user', () => {
    const user = ref<User | null>(null)
    const settings = ref<Settings>({})
    
    // Computed properties
    const isAuthenticated = computed(() => user.value !== null)
    
    // Actions
    async function login(email: string, password: string) {
        // Login implementation
        return true
    }
    
    return { user, settings, isAuthenticated, login }
})`

            const result = service.cleanupFileContent(content, 'typescript', {
                collapseImports: true,
                removeComments: true,
                removeTypeDefinitions: false,
                collapseBoilerplate: false
            })

            // Should collapse imports
            expect(result.content).toContain('// 3 imports collapsed')

            // Should remove comments
            expect(result.content).not.toContain('// Store definition')
            expect(result.content).not.toContain('// Computed properties')
            expect(result.content).not.toContain('// Login implementation')

            // Should preserve business logic
            expect(result.content).toContain('export const useUserStore')
            expect(result.content).toContain('async function login')

            // Should have significant savings
            expect(result.savingsPercent).toBeGreaterThan(10)
        })
    })
})
