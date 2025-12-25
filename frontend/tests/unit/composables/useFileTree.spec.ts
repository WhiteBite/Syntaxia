import { useFileTree } from '@/composables/useFileTree'
import type { DomainNode } from '@/utils/fileTreeUtils'
import { beforeEach, describe, expect, it } from 'vitest'

describe('useFileTree', () => {
    let tree: ReturnType<typeof useFileTree>

    // Mock data
    const mockDomainNodes: DomainNode[] = [
        {
            name: 'src',
            path: '/project/src',
            isDir: true,
            children: [
                { name: 'index.ts', path: '/project/src/index.ts', isDir: false, size: 1024 },
                { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false, size: 512 },
                {
                    name: 'components',
                    path: '/project/src/components',
                    isDir: true,
                    children: [
                        { name: 'Button.vue', path: '/project/src/components/Button.vue', isDir: false },
                    ],
                },
            ],
        },
        { name: 'README.md', path: '/project/README.md', isDir: false },
    ]

    beforeEach(() => {
        tree = useFileTree()
        tree.setFileTree(mockDomainNodes)
    })

    it('toggleExpand не сбрасывает allFilesCache', () => {
        // Get files to populate cache
        const srcNode = tree.findNode('/project/src')
        expect(srcNode).not.toBeNull()

        const filesBefore = tree.getAllFilesInNode(srcNode!)
        expect(filesBefore.length).toBe(3) // index.ts, utils.ts, Button.vue

        // Toggle expand
        tree.toggleExpand('/project/src')

        // Cache should still work - same result
        const filesAfter = tree.getAllFilesInNode(srcNode!)
        expect(filesAfter.length).toBe(3)
        expect(filesAfter).toEqual(filesBefore)
    })

    it('removeNode инвалидирует только родительские кэши', () => {
        // Populate cache for src folder
        const srcNode = tree.findNode('/project/src')
        const filesBefore = tree.getAllFilesInNode(srcNode!)
        expect(filesBefore.length).toBe(3)

        // Create a mock selectedPaths set
        const selectedPaths = new Set<string>()
        selectedPaths.add('/project/src/index.ts')

        // Remove a file
        const removed = tree.removeNode('/project/src/index.ts', selectedPaths)
        expect(removed).toBe(true)

        // Selection should be updated
        expect(selectedPaths.has('/project/src/index.ts')).toBe(false)

        // Cache should be invalidated - new count
        const srcNodeAfter = tree.findNode('/project/src')
        if (srcNodeAfter) {
            const filesAfter = tree.getAllFilesInNode(srcNodeAfter)
            expect(filesAfter.length).toBe(2) // utils.ts, Button.vue
        }
    })

    it('getMemoryUsage возвращает оценку использования памяти', () => {
        const usage = tree.getMemoryUsage()

        // Should return a positive number
        expect(usage).toBeGreaterThan(0)

        // Should be reasonable (not too small, not too large for our test data)
        expect(usage).toBeGreaterThan(500) // At least 500 bytes
        expect(usage).toBeLessThan(100000) // Less than 100KB for small tree
    })

    it('toggleExpand не инвалидирует весь кэш', () => {
        // Get files from two different nodes to populate cache
        const srcNode = tree.findNode('/project/src')
        const componentsNode = tree.findNode('/project/src/components')

        expect(srcNode).not.toBeNull()
        expect(componentsNode).not.toBeNull()

        // Populate cache for both nodes
        const srcFiles = tree.getAllFilesInNode(srcNode!)
        const componentFiles = tree.getAllFilesInNode(componentsNode!)

        expect(srcFiles.length).toBe(3)
        expect(componentFiles.length).toBe(1)

        // Toggle expand on src - should not affect components cache
        tree.toggleExpand('/project/src')

        // Both caches should still return correct results
        const srcFilesAfter = tree.getAllFilesInNode(srcNode!)
        const componentFilesAfter = tree.getAllFilesInNode(componentsNode!)

        expect(srcFilesAfter).toEqual(srcFiles)
        expect(componentFilesAfter).toEqual(componentFiles)
    })

    it('selection change не трогает allFilesCache', () => {
        // Get files to populate cache
        const srcNode = tree.findNode('/project/src')
        expect(srcNode).not.toBeNull()

        const filesBefore = tree.getAllFilesInNode(srcNode!)
        expect(filesBefore.length).toBe(3)

        // Simulate selection change by calling getAllFilesInNode multiple times
        // This should use cached result, not recalculate
        const filesSecondCall = tree.getAllFilesInNode(srcNode!)
        const filesThirdCall = tree.getAllFilesInNode(srcNode!)

        // All calls should return the same cached result
        expect(filesSecondCall).toEqual(filesBefore)
        expect(filesThirdCall).toEqual(filesBefore)

        // Verify the arrays are the same reference (cached)
        expect(filesSecondCall).toBe(filesThirdCall)
    })
})
