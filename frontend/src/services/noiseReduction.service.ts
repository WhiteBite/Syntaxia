import { useLogger } from '@/composables/useLogger'

const logger = useLogger('NoiseReductionService')

export interface NoiseReductionOptions {
    collapseImports: boolean
    removeComments: boolean
    removeTypeDefinitions: boolean
    collapseBoilerplate: boolean
}

export interface NoiseReductionResult {
    content: string
    originalTokens: number
    cleanedTokens: number
    savings: number
    savingsPercent: number
}

const BYTES_PER_TOKEN = 4 // Approximate

/**
 * Service for reducing noise in code files to optimize AI context
 */
export class NoiseReductionService {
    /**
     * Clean up file content based on language and options
     */
    cleanupFileContent(
        content: string,
        language: string,
        options: NoiseReductionOptions
    ): NoiseReductionResult {
        const originalTokens = Math.ceil(content.length / BYTES_PER_TOKEN)
        let cleaned = content

        try {
            // Apply language-specific cleanup
            switch (language.toLowerCase()) {
                case 'typescript':
                case 'javascript':
                case 'tsx':
                case 'jsx':
                    cleaned = this.cleanupJavaScriptLike(cleaned, options)
                    break
                case 'python':
                case 'py':
                    cleaned = this.cleanupPython(cleaned, options)
                    break
                case 'go':
                    cleaned = this.cleanupGo(cleaned, options)
                    break
                case 'vue':
                    cleaned = this.cleanupVue(cleaned, options)
                    break
                case 'css':
                case 'scss':
                case 'sass':
                case 'less':
                    cleaned = this.cleanupCSS(cleaned, options)
                    break
                default:
                    // Generic cleanup for unknown languages
                    cleaned = this.cleanupGeneric(cleaned, options)
            }

            const cleanedTokens = Math.ceil(cleaned.length / BYTES_PER_TOKEN)
            const savings = originalTokens - cleanedTokens
            const savingsPercent = originalTokens > 0 ? (savings / originalTokens) * 100 : 0

            return {
                content: cleaned,
                originalTokens,
                cleanedTokens,
                savings,
                savingsPercent
            }
        } catch (error) {
            logger.error('Failed to cleanup content:', error)
            // Return original content on error
            return {
                content,
                originalTokens,
                cleanedTokens: originalTokens,
                savings: 0,
                savingsPercent: 0
            }
        }
    }

    /**
     * Cleanup TypeScript/JavaScript files
     */
    private cleanupJavaScriptLike(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Remove comments FIRST (before collapsing imports)
        if (options.removeComments) {
            // Remove single-line comments
            result = result.replace(/\/\/.*$/gm, '')
            // Remove multi-line comments (but preserve JSDoc if not removing type definitions)
            if (options.removeTypeDefinitions) {
                result = result.replace(/\/\*[\s\S]*?\*\//g, '')
            } else {
                // Only remove non-JSDoc comments
                result = result.replace(/\/\*(?!\*)[\s\S]*?\*\//g, '')
            }
        }

        // Remove type definitions
        if (options.removeTypeDefinitions) {
            // Remove interface declarations (multiline)
            result = result.replace(/^export\s+interface\s+\w+[^{]*\{[^}]*\}\s*$/gm, '')
            // Remove type aliases - be more careful with multiline
            result = result.replace(/^export\s+type\s+\w+\s*=\s*'[^']*'\s*\|\s*'[^']*'(?:\s*\|\s*'[^']*')*\s*$/gm, '')
            // Remove simple type aliases
            result = result.replace(/^export\s+type\s+\w+\s*=\s*\w+\s*$/gm, '')
        }

        // Collapse boilerplate
        if (options.collapseBoilerplate) {
            // Remove 'use strict'
            result = result.replace(/['"]use strict['"];?\s*/g, '')
        }

        // Collapse imports LAST (after comments are removed)
        if (options.collapseImports) {
            result = this.collapseImportBlock(result, /^import\s+.*?from\s+['"].*?['"];?\s*$/m)
        }

        return this.normalizeWhitespace(result)
    }

    /**
     * Cleanup Python files
     */
    private cleanupPython(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Remove comments FIRST
        if (options.removeComments) {
            // Remove single-line comments
            result = result.replace(/#.*$/gm, '')
            // Remove docstrings
            result = result.replace(/"""[\s\S]*?"""/g, '')
            result = result.replace(/'''[\s\S]*?'''/g, '')
        }

        // Remove type hints
        if (options.removeTypeDefinitions) {
            // Remove function parameter type hints: name: str -> name
            result = result.replace(/(\w+):\s*\w+(\[.*?\])?(?=\s*[,)])/g, '$1')
            // Remove function return type hints
            result = result.replace(/\)\s*->\s*[^:]+:/g, '):')
            // Remove variable type hints
            result = result.replace(/(\w+):\s*\w+(\[.*?\])?\s*=/g, '$1 =')
        }

        // Collapse imports LAST
        if (options.collapseImports) {
            result = this.collapseImportBlock(result, /^(?:from\s+\S+\s+)?import\s+.*$/m)
        }

        return this.normalizeWhitespace(result)
    }

    /**
     * Cleanup Go files
     */
    private cleanupGo(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Collapse imports
        if (options.collapseImports) {
            // Collapse import blocks
            result = result.replace(/import\s*\(\s*\n([\s\S]*?)\n\s*\)/g, (_match, imports: string) => {
                const lines = imports.split('\n').filter((l: string) => l.trim())
                return `// ${lines.length} imports collapsed`
            })
            // Collapse single imports
            result = this.collapseImportBlock(result, /^import\s+".*?"$/gm)
        }

        // Remove comments
        if (options.removeComments) {
            // Remove single-line comments
            result = result.replace(/\/\/.*$/gm, '')
            // Remove multi-line comments
            result = result.replace(/\/\*[\s\S]*?\*\//g, '')
        }

        // Collapse package comments
        if (options.collapseBoilerplate) {
            result = result.replace(/^\/\/\s*Package\s+\w+[\s\S]*?(?=\npackage)/m, '// Package comment collapsed\n')
        }

        return this.normalizeWhitespace(result)
    }

    /**
     * Cleanup Vue files
     */
    private cleanupVue(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Extract and clean script section
        const scriptMatch = result.match(/<script[^>]*>([\s\S]*?)<\/script>/)
        if (scriptMatch) {
            const scriptContent = scriptMatch[1]

            const cleanedScript = this.cleanupJavaScriptLike(scriptContent, options)
            result = result.replace(scriptMatch[0], `<script>${cleanedScript}</script>`)
        }

        // Remove HTML comments if option enabled
        if (options.removeComments) {
            result = result.replace(/<!--[\s\S]*?-->/g, '')
        }

        return this.normalizeWhitespace(result)
    }

    /**
     * Cleanup CSS files
     */
    private cleanupCSS(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Remove comments
        if (options.removeComments) {
            result = result.replace(/\/\*[\s\S]*?\*\//g, '')
        }

        // Collapse @import statements
        if (options.collapseImports) {
            result = this.collapseImportBlock(result, /^@import\s+.*?;$/m)
        }

        return this.normalizeWhitespace(result)
    }

    /**
     * Generic cleanup for unknown file types
     */
    private cleanupGeneric(content: string, options: NoiseReductionOptions): string {
        let result = content

        // Only apply basic whitespace normalization
        if (options.collapseBoilerplate) {
            result = this.normalizeWhitespace(result)
        }

        return result
    }

    /**
     * Collapse import blocks into a single comment line
     */
    private collapseImportBlock(content: string, pattern: RegExp): string {
        const lines = content.split('\n')
        const result: string[] = []
        const importLines: string[] = []
        let inImportBlock = false

        for (let i = 0; i < lines.length; i++) {
            const line = lines[i]
            const isImport = pattern.test(line)
            const isEmpty = line.trim() === ''

            if (isImport) {
                // This is an import line - collect it
                importLines.push(line)
                inImportBlock = true
            } else if (isEmpty && inImportBlock) {
                // Empty line within import block - continue collecting
                continue
            } else {
                // Non-import, non-empty line - flush collected imports
                if (importLines.length > 0) {
                    result.push(`// ${importLines.length} imports collapsed`)
                    importLines.length = 0
                    inImportBlock = false
                }
                result.push(line)
            }
        }

        // Flush any remaining imports at end of file
        if (importLines.length > 0) {
            result.push(`// ${importLines.length} imports collapsed`)
        }

        return result.join('\n')
    }

    /**
     * Normalize whitespace - remove excessive blank lines
     */
    private normalizeWhitespace(content: string): string {
        // Replace 3+ consecutive newlines with 2 newlines
        return content
            .replace(/\n{3,}/g, '\n\n')
            .replace(/^\s+$/gm, '') // Remove whitespace-only lines
            .trim()
    }

    /**
     * Get language from file extension
     */
    getLanguageFromExtension(filename: string): string {
        const ext = filename.split('.').pop()?.toLowerCase() || ''

        const languageMap: Record<string, string> = {
            'ts': 'typescript',
            'tsx': 'tsx',
            'js': 'javascript',
            'jsx': 'jsx',
            'py': 'python',
            'go': 'go',
            'vue': 'vue',
            'css': 'css',
            'scss': 'scss',
            'sass': 'sass',
            'less': 'less',
            'java': 'java',
            'cpp': 'cpp',
            'c': 'c',
            'rs': 'rust',
            'rb': 'ruby',
            'php': 'php',
            'swift': 'swift',
            'kt': 'kotlin',
            'cs': 'csharp'
        }

        return languageMap[ext] || 'unknown'
    }
}

export const noiseReductionService = new NoiseReductionService()
