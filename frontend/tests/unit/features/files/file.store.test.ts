/**
 * Unit tests for file.store.ts
 * Tests the unified file tree management store
 */
import { useFileStore } from '@/features/files/model/file.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock the files API
vi.mock('@/features/files/api/files.api', () => ({
    filesApi: {
        listFiles: vi.fn(),
        readFileContent: vi.fn(),
        getFileStats: vi.fn(),
        clearCache: vi.fn(),
    },
}))

// Mock the logger
vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn(),
    }),
}))

// Mock settings store
vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        settings: {
            fileExplorer: {
                autoSaveSelection: false,
            },
        },
    }),
}))

describe('FileStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorage.clear()
    })

    describe('Initial State', () => {
        it('should have empty nodes initially', () => {
            const store = useFileStore()
            expect(store.nodes).toEqual([])
        })

        it('should have empty rootPath initially', () => {
            const store = useFileStore()
            expect(store.rootPath).toBe('')
        })

        it('should have no selected files initially', () => {
            const store = useFileStore()
            expect(store.hasSelectedFiles).toBe(false)
            expect(store.selectedCount).toBe(0)
        })

        it('should not be loading initially', () => {
            const store = useFileStore()
            expect(store.isLoading).toBe(false)
        })

        it('should have no error initially', () => {
            const store = useFileStore()
            expect(store.error).toBeNull()
        })
    })

    describe('setFileTree', () => {
        it('should set file tree from domain nodes', () => {
            const store = useFileStore()
            const mockTree = [
                {
                    name: 'src',
                    path: '/project/src',
                    isDir: true,
                    children: [
                        { name: 'main.ts', path: '/project/src/main.ts', isDir: false },
                    ],
                },
            ]

            store.setFileTree(mockTree)

            expect(store.nodes.length).toBe(1)
            expect(store.nodes[0].name).toBe('src')
            expect(store.nodes[0].isDir).toBe(true)
        })

        it('should handle empty tree', () => {
            const store = useFileStore()
            store.setFileTree([])
            expect(store.nodes).toEqual([])
        })
    })

    describe('Selection Operations', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'main.ts', path: '/project/src/main.ts', isDir: false, size: 100 },
                    { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false, size: 200 },
                ],
            },
            { name: 'README.md', path: '/project/README.md', isDir: false, size: 50 },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should toggle select a file', () => {
            const store = useFileStore()
            store.toggleSelect('/project/src/main.ts')

            expect(store.hasSelectedFiles).toBe(true)
            expect(store.selectedCount).toBe(1)
            expect(store.selectedFilesList).toContain('/project/src/main.ts')
        })

        it('should toggle deselect a file', () => {
            const store = useFileStore()
            store.toggleSelect('/project/src/main.ts')
            store.toggleSelect('/project/src/main.ts')

            expect(store.hasSelectedFiles).toBe(false)
            expect(store.selectedCount).toBe(0)
        })

        it('should select multiple files', () => {
            const store = useFileStore()
            store.selectMultiple(['/project/src/main.ts', '/project/src/utils.ts'])

            expect(store.selectedCount).toBe(2)
        })

        it('should clear selection', () => {
            const store = useFileStore()
            store.selectMultiple(['/project/src/main.ts', '/project/src/utils.ts'])
            store.clearSelection()

            expect(store.hasSelectedFiles).toBe(false)
            expect(store.selectedCount).toBe(0)
        })

        it('should select recursively when toggling a directory', () => {
            const store = useFileStore()
            store.toggleSelect('/project/src')

            expect(store.selectedCount).toBe(2)
            expect(store.selectedFilesList).toContain('/project/src/main.ts')
            expect(store.selectedFilesList).toContain('/project/src/utils.ts')
        })

        it('should deselect recursively when toggling a selected directory', () => {
            const store = useFileStore()
            store.toggleSelect('/project/src')
            store.toggleSelect('/project/src')

            expect(store.selectedCount).toBe(0)
        })

        it('should select by extension', () => {
            const store = useFileStore()
            store.selectByExtension('.ts')

            expect(store.selectedCount).toBe(2)
            expect(store.selectedFilesList).toContain('/project/src/main.ts')
            expect(store.selectedFilesList).toContain('/project/src/utils.ts')
            expect(store.selectedFilesList).not.toContain('/project/README.md')
        })
    })

    describe('Expand/Collapse Operations', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
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
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should toggle expand a directory', () => {
            const store = useFileStore()
            const expandedBefore = store.getExpandedPaths()

            store.toggleExpand('/project/src')

            const expandedAfter = store.getExpandedPaths()
            expect(expandedAfter.includes('/project/src')).not.toBe(expandedBefore.includes('/project/src'))
        })

        it('should expand all directories', () => {
            const store = useFileStore()
            store.expandAll()

            const expanded = store.getExpandedPaths()
            expect(expanded).toContain('/project/src')
            expect(expanded).toContain('/project/src/components')
        })

        it('should collapse all directories', () => {
            const store = useFileStore()
            store.expandAll()
            store.collapseAll()

            const expanded = store.getExpandedPaths()
            expect(expanded.length).toBe(0)
        })

        it('should expand recursively', () => {
            const store = useFileStore()
            store.collapseAll()
            store.expandRecursive('/project/src')

            const expanded = store.getExpandedPaths()
            expect(expanded).toContain('/project/src')
            expect(expanded).toContain('/project/src/components')
        })

        it('should collapse recursively', () => {
            const store = useFileStore()
            store.expandAll()
            store.collapseRecursive('/project/src')

            const expanded = store.getExpandedPaths()
            expect(expanded).not.toContain('/project/src')
            expect(expanded).not.toContain('/project/src/components')
        })
    })

    describe('Computed Properties', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'main.ts', path: '/project/src/main.ts', isDir: false, size: 1000 },
                    { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false, size: 2000 },
                ],
            },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setRootPath('/project')
            store.setFileTree(mockTree)
        })

        it('should compute project name from root path', () => {
            const store = useFileStore()
            expect(store.projectName).toBe('project')
        })

        it('should compute estimated token count', () => {
            const store = useFileStore()
            store.selectMultiple(['/project/src/main.ts', '/project/src/utils.ts'])

            // Token count = total size / 4
            expect(store.estimatedTokenCount).toBe(750) // (1000 + 2000) / 4
        })

        it('should compute estimated context size in MB', () => {
            const store = useFileStore()
            store.selectMultiple(['/project/src/main.ts', '/project/src/utils.ts'])

            // Context size = total size / (1024 * 1024)
            const expectedSize = 3000 / (1024 * 1024)
            expect(store.estimatedContextSize).toBeCloseTo(expectedSize, 6)
        })

        it('should return available extensions', () => {
            const store = useFileStore()
            const extensions = store.getAvailableExtensions()

            expect(extensions).toContain('.ts')
        })
    })

    describe('Search Operations', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'main.ts', path: '/project/src/main.ts', isDir: false },
                    { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false },
                ],
            },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should set search query', () => {
            const store = useFileStore()
            store.setSearchQuery('main')

            expect(store.searchQuery).toBe('main')
        })

        it('should clear search query', () => {
            const store = useFileStore()
            store.setSearchQuery('main')
            store.setSearchQuery('')

            expect(store.searchQuery).toBe('')
        })
    })

    describe('Node Operations', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    { name: 'main.ts', path: '/project/src/main.ts', isDir: false },
                ],
            },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should check if node exists', () => {
            const store = useFileStore()

            expect(store.nodeExists('/project/src')).toBe(true)
            expect(store.nodeExists('/project/src/main.ts')).toBe(true)
            expect(store.nodeExists('/project/nonexistent')).toBe(false)
        })

        it('should check if path is directory', () => {
            const store = useFileStore()

            expect(store.isDirectory('/project/src')).toBe(true)
            expect(store.isDirectory('/project/src/main.ts')).toBe(false)
        })

        it('should remove node from tree', () => {
            const store = useFileStore()
            const removed = store.removeNode('/project/src/main.ts')

            expect(removed).toBe(true)
            expect(store.nodeExists('/project/src/main.ts')).toBe(false)
        })

        it('should deselect file when removing selected node', () => {
            const store = useFileStore()
            store.toggleSelect('/project/src/main.ts')
            expect(store.selectedCount).toBe(1)

            store.removeNode('/project/src/main.ts')
            expect(store.selectedCount).toBe(0)
        })
    })

    describe('Reset Store', () => {
        it('should reset all state', () => {
            const store = useFileStore()
            store.setFileTree([
                { name: 'test.ts', path: '/test.ts', isDir: false },
            ])
            store.setRootPath('/project')
            store.toggleSelect('/test.ts')

            store.resetStore()

            expect(store.nodes).toEqual([])
            expect(store.rootPath).toBe('')
            expect(store.selectedCount).toBe(0)
            expect(store.error).toBeNull()
            expect(store.isLoading).toBe(false)
        })
    })

    describe('Undo/Redo Selection', () => {
        const mockTree = [
            { name: 'file1.ts', path: '/file1.ts', isDir: false },
            { name: 'file2.ts', path: '/file2.ts', isDir: false },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should undo selection', () => {
            const store = useFileStore()
            store.toggleSelect('/file1.ts')
            expect(store.selectedCount).toBe(1)

            store.undoSelection()
            expect(store.selectedCount).toBe(0)
        })

        it('should track canUndo after selection', () => {
            const store = useFileStore()
            expect(store.canUndoSelection).toBe(false)

            store.toggleSelect('/file1.ts')
            expect(store.canUndoSelection).toBe(true)
        })

        it('should track canRedo after undo', () => {
            const store = useFileStore()
            store.toggleSelect('/file1.ts')
            expect(store.canRedoSelection).toBe(false)

            store.undoSelection()
            expect(store.canRedoSelection).toBe(true)
        })
    })

    describe('Weight Filter Operations', () => {
        const mockTree = [
            {
                name: 'src',
                path: '/project/src',
                isDir: true,
                children: [
                    // Small file: 1000 bytes = ~250 tokens
                    { name: 'small.ts', path: '/project/src/small.ts', isDir: false, size: 1000 },
                    // Medium file: 40000 bytes = ~10000 tokens
                    { name: 'medium.ts', path: '/project/src/medium.ts', isDir: false, size: 40000 },
                    // Heavy file: 200000 bytes = ~50000 tokens
                    { name: 'heavy.ts', path: '/project/src/heavy.ts', isDir: false, size: 200000 },
                    // Critical file: 400000 bytes = ~100000 tokens
                    { name: 'critical.ts', path: '/project/src/critical.ts', isDir: false, size: 400000 },
                ],
            },
        ]

        beforeEach(() => {
            const store = useFileStore()
            store.setFileTree(mockTree)
        })

        it('should have no weight filter by default', () => {
            const store = useFileStore()
            expect(store.weightFilter).toBe('none')
        })

        it('should set weight filter to medium', () => {
            const store = useFileStore()
            store.setWeightFilter('medium')
            expect(store.weightFilter).toBe('medium')
        })

        it('should set weight filter to heavy', () => {
            const store = useFileStore()
            store.setWeightFilter('heavy')
            expect(store.weightFilter).toBe('heavy')
        })

        it('should set weight filter to critical', () => {
            const store = useFileStore()
            store.setWeightFilter('critical')
            expect(store.weightFilter).toBe('critical')
        })

        it('should clear weight filter', () => {
            const store = useFileStore()
            store.setWeightFilter('heavy')
            store.clearWeightFilter()
            expect(store.weightFilter).toBe('none')
        })

        it('should filter files by medium weight (10K+ tokens)', () => {
            const store = useFileStore()
            store.setWeightFilter('medium')

            const filtered = store.filteredNodes
            expect(filtered.length).toBe(1)
            expect(filtered[0].isDir).toBe(true)
            expect(filtered[0].children?.length).toBe(3) // medium, heavy, critical
        })

        it('should filter files by heavy weight (50K+ tokens)', () => {
            const store = useFileStore()
            store.setWeightFilter('heavy')

            const filtered = store.filteredNodes
            expect(filtered.length).toBe(1)
            expect(filtered[0].isDir).toBe(true)
            expect(filtered[0].children?.length).toBe(2) // heavy, critical
        })

        it('should filter files by critical weight (100K+ tokens)', () => {
            const store = useFileStore()
            store.setWeightFilter('critical')

            const filtered = store.filteredNodes
            expect(filtered.length).toBe(1)
            expect(filtered[0].isDir).toBe(true)
            expect(filtered[0].children?.length).toBe(1) // critical only
        })

        it('should show all files when weight filter is none', () => {
            const store = useFileStore()
            store.setWeightFilter('none')

            const filtered = store.filteredNodes
            expect(filtered.length).toBe(1)
            expect(filtered[0].children?.length).toBe(4) // all files
        })

        it('should combine weight filter with extension filter', () => {
            const store = useFileStore()
            store.setWeightFilter('medium')
            store.setFilterExtensions(['.ts'], [])

            const filtered = store.filteredNodes
            expect(filtered.length).toBe(1)
            expect(filtered[0].children?.length).toBe(3) // medium, heavy, critical (all .ts)
        })
    })
})
