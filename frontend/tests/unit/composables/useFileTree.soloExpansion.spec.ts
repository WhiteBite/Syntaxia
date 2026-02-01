import { useFileTree } from '@/composables/useFileTree'
import type { FileNode } from '@/types/domain'
import { beforeEach, describe, expect, it } from 'vitest'

describe('useFileTree - Solo Expansion Mode', () => {
    let tree: ReturnType<typeof useFileTree>

    const mockTree: FileNode[] = [
        {
            name: 'folder1',
            path: 'folder1',
            isDir: true,
            isExpanded: false,
            children: [
                { name: 'file1.ts', path: 'folder1/file1.ts', isDir: false },
            ],
        },
        {
            name: 'folder2',
            path: 'folder2',
            isDir: true,
            isExpanded: false,
            children: [
                { name: 'file2.ts', path: 'folder2/file2.ts', isDir: false },
            ],
        },
        {
            name: 'folder3',
            path: 'folder3',
            isDir: true,
            isExpanded: false,
            children: [
                {
                    name: 'subfolder1',
                    path: 'folder3/subfolder1',
                    isDir: true,
                    isExpanded: false,
                    children: [
                        { name: 'file3.ts', path: 'folder3/subfolder1/file3.ts', isDir: false },
                    ],
                },
                {
                    name: 'subfolder2',
                    path: 'folder3/subfolder2',
                    isDir: true,
                    isExpanded: false,
                    children: [
                        { name: 'file4.ts', path: 'folder3/subfolder2/file4.ts', isDir: false },
                    ],
                },
            ],
        },
    ]

    beforeEach(() => {
        tree = useFileTree()
        tree.setFileTree(JSON.parse(JSON.stringify(mockTree)))
    })

    it('should collapse siblings when expanding a folder in solo mode', () => {
        // Verify initial state
        expect(tree.nodes.value.length).toBe(3)

        // Expand folder1 and folder2 normally (without solo mode)
        tree.toggleExpand('folder1', false)
        tree.toggleExpand('folder2', false)

        expect(tree.nodes.value[0].isExpanded).toBe(true)
        expect(tree.nodes.value[1].isExpanded).toBe(true)
        expect(tree.nodes.value[2].isExpanded).toBe(false)

        // Now expand folder3 in solo mode - should collapse folder1 and folder2
        tree.toggleExpand('folder3', true)

        // Check directly on nodes array
        expect(tree.nodes.value[0].isExpanded).toBe(false)
        expect(tree.nodes.value[1].isExpanded).toBe(false)
        expect(tree.nodes.value[2].isExpanded).toBe(true)
    })

    it('should only collapse siblings at the same level', () => {
        // Expand folder3 and its subfolder1
        tree.toggleExpand('folder3', false)
        tree.toggleExpand('folder3/subfolder1', false)

        expect(tree.findNode('folder3')?.isExpanded).toBe(true)
        expect(tree.findNode('folder3/subfolder1')?.isExpanded).toBe(true)

        // Expand subfolder2 in solo mode - should only collapse subfolder1, not folder3
        tree.toggleExpand('folder3/subfolder2', true)

        expect(tree.findNode('folder3')?.isExpanded).toBe(true) // Parent stays expanded
        expect(tree.findNode('folder3/subfolder1')?.isExpanded).toBe(false) // Sibling collapsed
        expect(tree.findNode('folder3/subfolder2')?.isExpanded).toBe(true) // Current expanded
    })

    it('should not collapse siblings when collapsing a folder in solo mode', () => {
        // Expand all root folders
        tree.toggleExpand('folder1', false)
        tree.toggleExpand('folder2', false)
        tree.toggleExpand('folder3', false)

        expect(tree.findNode('folder1')?.isExpanded).toBe(true)
        expect(tree.findNode('folder2')?.isExpanded).toBe(true)
        expect(tree.findNode('folder3')?.isExpanded).toBe(true)

        // Collapse folder2 in solo mode - should not affect others
        tree.toggleExpand('folder2', true)

        expect(tree.findNode('folder1')?.isExpanded).toBe(true)
        expect(tree.findNode('folder2')?.isExpanded).toBe(false)
        expect(tree.findNode('folder3')?.isExpanded).toBe(true)
    })

    it('should work correctly without solo mode', () => {
        // Expand all folders without solo mode
        tree.toggleExpand('folder1', false)
        tree.toggleExpand('folder2', false)
        tree.toggleExpand('folder3', false)

        // All should remain expanded
        expect(tree.findNode('folder1')?.isExpanded).toBe(true)
        expect(tree.findNode('folder2')?.isExpanded).toBe(true)
        expect(tree.findNode('folder3')?.isExpanded).toBe(true)
    })

    it('should handle root level folders correctly in solo mode', () => {
        // Expand folder1 in solo mode
        tree.toggleExpand('folder1', true)
        expect(tree.findNode('folder1')?.isExpanded).toBe(true)

        // Expand folder2 in solo mode - should collapse folder1
        tree.toggleExpand('folder2', true)
        expect(tree.findNode('folder1')?.isExpanded).toBe(false)
        expect(tree.findNode('folder2')?.isExpanded).toBe(true)

        // Expand folder3 in solo mode - should collapse folder2
        tree.toggleExpand('folder3', true)
        expect(tree.findNode('folder1')?.isExpanded).toBe(false)
        expect(tree.findNode('folder2')?.isExpanded).toBe(false)
        expect(tree.findNode('folder3')?.isExpanded).toBe(true)
    })
})
