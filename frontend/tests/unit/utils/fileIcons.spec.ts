import { FILE_TYPE_CONFIG, getFileIcon, getFileTypeInfo } from '@/utils/fileIcons'
import { describe, expect, it } from 'vitest'

describe('fileIcons', () => {
    describe('getFileIcon', () => {
        it('should return correct icon for TypeScript files', () => {
            expect(getFileIcon('app.ts')).toBe('🔷')
            expect(getFileIcon('component.tsx')).toBe('⚛️')
        })

        it('should return correct icon for JavaScript files', () => {
            expect(getFileIcon('index.js')).toBe('🟨')
            expect(getFileIcon('App.jsx')).toBe('⚛️')
        })

        it('should return correct icon for Vue files', () => {
            expect(getFileIcon('Component.vue')).toBe('💚')
        })

        it('should return correct icon for Go files', () => {
            expect(getFileIcon('main.go')).toBe('🐹')
        })

        it('should return correct icon for Python files', () => {
            expect(getFileIcon('script.py')).toBe('🐍')
        })

        it('should return correct icon for config files', () => {
            expect(getFileIcon('config.json')).toBe('📋')
            expect(getFileIcon('settings.yaml')).toBe('⚙️')
            expect(getFileIcon('config.yml')).toBe('⚙️')
        })

        it('should return correct icon for markdown files', () => {
            expect(getFileIcon('README.md')).toBe('📖')
            expect(getFileIcon('docs.md')).toBe('📝')
        })

        it('should return correct icon for special files', () => {
            expect(getFileIcon('package.json')).toBe('📦')
            expect(getFileIcon('Dockerfile')).toBe('🐳')
            expect(getFileIcon('Makefile')).toBe('🔨')
            expect(getFileIcon('.gitignore')).toBe('🚫')
            expect(getFileIcon('go.mod')).toBe('🐹')
        })

        it('should return default icon for unknown extensions', () => {
            expect(getFileIcon('unknown.xyz')).toBe('📄')
            expect(getFileIcon('file.unknown')).toBe('📄')
        })

        it('should handle files without extension', () => {
            expect(getFileIcon('LICENSE')).toBe('📜')
        })
    })

    describe('getFileTypeInfo', () => {
        it('should return icon and colorClass for known types', () => {
            const info = getFileTypeInfo('app.ts')
            expect(info.icon).toBe('🔷')
            expect(info.colorClass).toBe('bg-blue-500')
        })

        it('should return info for special files', () => {
            const info = getFileTypeInfo('package.json')
            expect(info.icon).toBe('📦')
            expect(info.colorClass).toBe('bg-red-500')
        })

        it('should return default info for unknown types', () => {
            const info = getFileTypeInfo('unknown.xyz')
            expect(info.icon).toBe('📄')
            expect(info.colorClass).toBe('bg-gray-500')
        })
    })

    describe('FILE_TYPE_CONFIG', () => {
        it('should have default entry', () => {
            expect(FILE_TYPE_CONFIG.default).toBeDefined()
            expect(FILE_TYPE_CONFIG.default.icon).toBe('📄')
        })

        it('should have entries for common file types', () => {
            expect(FILE_TYPE_CONFIG.ts).toBeDefined()
            expect(FILE_TYPE_CONFIG.js).toBeDefined()
            expect(FILE_TYPE_CONFIG.vue).toBeDefined()
            expect(FILE_TYPE_CONFIG.go).toBeDefined()
            expect(FILE_TYPE_CONFIG.py).toBeDefined()
        })
    })
})
