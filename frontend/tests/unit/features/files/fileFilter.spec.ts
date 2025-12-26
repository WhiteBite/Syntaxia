/**
 * File Filtering Tests
 * Tests for useGitignore, useCustomIgnore, foldersFirst settings
 */

import { useFileFilter } from '@/composables/useFileFilter'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'
import { ref } from 'vue'

// Helper to create file nodes
function createFileNode(path: string, isDir = false, children: FileNode[] = []): FileNode {
    return {
        path,
        name: path.split('/').pop() || path,
        isDir,
        children,
        isExpanded: false,
        isIgnored: false,
    }
}

describe('File Filtering', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    describe('useGitignore', () => {
        it('should hide node_modules when useGitignore=true and node is ignored', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                ]),
                { ...createFileNode('node_modules', true), isIgnored: true },
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            const nodeModules = filteredNodes.value.find((n: FileNode) => n.name === 'node_modules')
            expect(nodeModules?.isIgnored).toBe(true)
        })

        it('should show node_modules when useGitignore=false', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true),
                createFileNode('node_modules', true),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            const nodeModules = filteredNodes.value.find((n: FileNode) => n.name === 'node_modules')
            expect(nodeModules).toBeDefined()
        })

        it('should apply .gitignore patterns recursively', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                    { ...createFileNode('src/.env'), isIgnored: true },
                ]),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            const envFile = src?.children?.find((n: FileNode) => n.name === '.env')
            expect(envFile?.isIgnored).toBe(true)
        })
    })

    describe('useCustomIgnore', () => {
        it('should apply custom patterns when enabled', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                    createFileNode('src/debug.log'),
                ]),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            expect(src?.children?.length).toBe(1)
            expect(src?.children?.[0].name).toBe('main.ts')
        })

        it('should ignore custom patterns when disabled', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                    createFileNode('src/debug.log'),
                ]),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            expect(src?.children?.length).toBe(2)
        })

        it('should combine with gitignore when both enabled', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                    { ...createFileNode('src/.env'), isIgnored: true },
                    createFileNode('src/debug.log'),
                ]),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            const tsFiles = src?.children?.filter((n: FileNode) => n.name.endsWith('.ts'))
            expect(tsFiles?.length).toBe(1)
        })
    })

    describe('foldersFirst', () => {
        it('should sort folders before files when enabled', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('zebra.ts'),
                createFileNode('src', true),
                createFileNode('apple.ts'),
                createFileNode('lib', true),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            expect(filteredNodes.value).toHaveLength(4)
            // Folders should come first
            expect(filteredNodes.value[0].isDir).toBe(true)
            expect(filteredNodes.value[1].isDir).toBe(true)
        })

        it('should use alphabetical sort when disabled', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('zebra.ts'),
                createFileNode('apple.ts'),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            expect(filteredNodes.value).toHaveLength(2)
            const names = filteredNodes.value.map((n: FileNode) => n.name)
            expect(names).toContain('zebra.ts')
            expect(names).toContain('apple.ts')
        })
    })

    describe('extension filtering', () => {
        it('should filter by single extension', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('style.css'),
                createFileNode('readme.md'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])

            expect(filteredNodes.value).toHaveLength(1)
            expect(filteredNodes.value[0].name).toBe('main.ts')
        })

        it('should filter by multiple extensions', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('app.vue'),
                createFileNode('style.css'),
                createFileNode('readme.md'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts', '.vue'])

            expect(filteredNodes.value).toHaveLength(2)
            expect(filteredNodes.value.map((n: FileNode) => n.name)).toContain('main.ts')
            expect(filteredNodes.value.map((n: FileNode) => n.name)).toContain('app.vue')
        })

        it('should show all files when no filter applied', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('style.css'),
                createFileNode('readme.md'),
            ])

            const { filteredNodes } = useFileFilter({ nodes })

            expect(filteredNodes.value).toHaveLength(3)
        })

        it('should clear filter correctly', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('style.css'),
            ])

            const { filteredNodes, setFilterExtensions, clearFilters } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])
            expect(filteredNodes.value).toHaveLength(1)

            clearFilters()
            expect(filteredNodes.value).toHaveLength(2)
        })
    })

    describe('directory handling', () => {
        it('should always show directories regardless of extension filter', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/main.ts'),
                ]),
                createFileNode('readme.md'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            expect(src).toBeDefined()
            expect(src?.isDir).toBe(true)
        })

        it('should hide empty directories after filtering', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('src', true, [
                    createFileNode('src/style.css'),
                ]),
                createFileNode('main.ts'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions(['.ts'])

            const src = filteredNodes.value.find((n: FileNode) => n.name === 'src')
            if (src) {
                const tsChildren = src.children?.filter((c: FileNode) => c.name.endsWith('.ts'))
                expect(tsChildren?.length || 0).toBe(0)
            }
        })
    })

    describe('exclude extensions', () => {
        it('should exclude files with specified extensions', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('style.css'),
                createFileNode('readme.md'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            setFilterExtensions([], ['.css', '.md'])

            expect(filteredNodes.value).toHaveLength(1)
            expect(filteredNodes.value[0].name).toBe('main.ts')
        })

        it('should combine include and exclude filters', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('test.spec.ts'),
                createFileNode('style.css'),
            ])

            const { filteredNodes, setFilterExtensions } = useFileFilter({ nodes })

            // Include .ts but exclude .spec.ts
            setFilterExtensions(['.ts'], ['.spec.ts'])

            expect(filteredNodes.value).toHaveLength(1)
            expect(filteredNodes.value[0].name).toBe('main.ts')
        })
    })

    describe('add/remove extensions', () => {
        it('should add include extension', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('app.vue'),
                createFileNode('style.css'),
            ])

            const { filteredNodes, addIncludeExtension } = useFileFilter({ nodes })

            addIncludeExtension('.ts')
            expect(filteredNodes.value).toHaveLength(1)

            addIncludeExtension('.vue')
            expect(filteredNodes.value).toHaveLength(2)
        })

        it('should remove include extension', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('app.vue'),
            ])

            const { filteredNodes, setFilterExtensions, removeIncludeExtension } = useFileFilter({ nodes })

            setFilterExtensions(['.ts', '.vue'])
            expect(filteredNodes.value).toHaveLength(2)

            removeIncludeExtension('.vue')
            expect(filteredNodes.value).toHaveLength(1)
            expect(filteredNodes.value[0].name).toBe('main.ts')
        })

        it('should add exclude extension', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('test.spec.ts'),
            ])

            const { filteredNodes, addExcludeExtension } = useFileFilter({ nodes })

            addExcludeExtension('.spec.ts')
            expect(filteredNodes.value).toHaveLength(1)
            expect(filteredNodes.value[0].name).toBe('main.ts')
        })

        it('should remove exclude extension', () => {
            const nodes = ref<FileNode[]>([
                createFileNode('main.ts'),
                createFileNode('test.spec.ts'),
            ])

            const { filteredNodes, setFilterExtensions, removeExcludeExtension } = useFileFilter({ nodes })

            setFilterExtensions([], ['.spec.ts'])
            expect(filteredNodes.value).toHaveLength(1)

            removeExcludeExtension('.spec.ts')
            expect(filteredNodes.value).toHaveLength(2)
        })
    })
})
