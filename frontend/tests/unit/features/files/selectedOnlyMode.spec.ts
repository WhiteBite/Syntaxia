/**
 * Tests for Selected Only Mode feature
 */

import { useFileStore } from '@/features/files/model/file.store'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Selected Only Mode', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('should toggle selected-only mode', () => {
        const store = useFileStore()

        expect(store.isSelectedOnlyMode).toBe(false)

        store.toggleSelectedOnlyMode()
        expect(store.isSelectedOnlyMode).toBe(true)

        store.toggleSelectedOnlyMode()
        expect(store.isSelectedOnlyMode).toBe(false)
    })

    it('should show only selected files in flat list when mode is active', () => {
        const store = useFileStore()

        // Setup test tree
        const testTree: FileNode[] = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'file1.ts', path: '/project/src/file1.ts', isDir: false, size: 100 },
                    { name: 'file2.ts', path: '/project/src/file2.ts', isDir: false, size: 200 },
                ]
            },
            {
                name: 'tests',
                path: '/project/tests',
                isDir: true,
                children: [
                    { name: 'test1.spec.ts', path: '/project/tests/test1.spec.ts', isDir: false, size: 150 },
                ]
            }
        ]

        store.setFileTree(testTree)
        store.setRootPath('/project')

        // Select some files
        store.selectPath('/project/src/file1.ts')
        store.selectPath('/project/tests/test1.spec.ts')

        expect(store.selectedCount).toBe(2)

        // Enable selected-only mode
        store.toggleSelectedOnlyMode()

        // filteredNodes should now only show selected files in flat list
        const filtered = store.filteredNodes

        // Should only have 2 files (no directories)
        expect(filtered.length).toBe(2)

        // All should be files (not directories)
        expect(filtered.every(node => !node.isDir)).toBe(true)

        // All should have depth 0 (flat list)
        expect(filtered.every(node => node.depth === 0)).toBe(true)

        // Should contain only selected files
        const paths = filtered.map(node => node.path)
        expect(paths).toContain('/project/src/file1.ts')
        expect(paths).toContain('/project/tests/test1.spec.ts')
    })

    it('should show normal tree when mode is disabled', () => {
        const store = useFileStore()

        const testTree: FileNode[] = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'file1.ts', path: '/project/src/file1.ts', isDir: false, size: 100 },
                ]
            }
        ]

        store.setFileTree(testTree)
        store.selectPath('/project/src/file1.ts')

        // Mode disabled - should show normal tree with folders
        expect(store.isSelectedOnlyMode).toBe(false)
        const filtered = store.filteredNodes

        // Should include directory
        expect(filtered.some(node => node.isDir)).toBe(true)
    })

    it('should update filtered list when selection changes in selected-only mode', () => {
        const store = useFileStore()

        const testTree: FileNode[] = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'file1.ts', path: '/project/src/file1.ts', isDir: false, size: 100 },
                    { name: 'file2.ts', path: '/project/src/file2.ts', isDir: false, size: 200 },
                ]
            }
        ]

        store.setFileTree(testTree)
        store.setRootPath('/project')
        store.selectPath('/project/src/file1.ts')

        // Enable mode
        store.toggleSelectedOnlyMode()
        expect(store.filteredNodes.length).toBe(1)

        // Add another file
        store.selectPath('/project/src/file2.ts')
        expect(store.filteredNodes.length).toBe(2)

        // Remove a file
        store.deselectPath('/project/src/file1.ts')
        expect(store.filteredNodes.length).toBe(1)
    })

    it('should show empty list when no files selected in selected-only mode', () => {
        const store = useFileStore()

        const testTree: FileNode[] = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'file1.ts', path: '/project/src/file1.ts', isDir: false, size: 100 },
                ]
            }
        ]

        store.setFileTree(testTree)

        // Enable mode with no selection
        store.toggleSelectedOnlyMode()

        expect(store.selectedCount).toBe(0)
        expect(store.filteredNodes.length).toBe(0)
    })
})
