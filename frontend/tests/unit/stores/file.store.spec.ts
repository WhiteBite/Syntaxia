import { useFileStore, type FileNode } from '@/features/files/model/file.store'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock filesApi
vi.mock('@/features/files/api/files.api', () => ({
    filesApi: {
        listFiles: vi.fn(),
        clearCache: vi.fn()
    }
}))

// Mock settings store
vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        settings: {
            fileExplorer: {
                autoSaveSelection: false
            }
        }
    })
}))

describe('FileStore', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    describe('removeNode', () => {
        it('should remove a file node from tree', () => {
            const store = useFileStore()

            // Setup tree
            const tree: FileNode[] = [{
                name: 'root',
                path: '/root',
                isDir: true,
                children: [
                    { name: 'file1.ts', path: '/root/file1.ts', isDir: false },
                    { name: 'file2.ts', path: '/root/file2.ts', isDir: false }
                ]
            }]
            store.setFileTree(tree)

            // Remove file1.ts
            const removed = store.removeNode('/root/file1.ts')

            expect(removed).toBe(true)
            expect(store.nodes[0].children).toHaveLength(1)
            expect(store.nodes[0].children![0].name).toBe('file2.ts')
        })

        it('should remove a directory node with all children', () => {
            const store = useFileStore()

            const tree: FileNode[] = [{
                name: 'root',
                path: '/root',
                isDir: true,
                children: [
                    {
                        name: 'docs',
                        path: '/root/docs',
                        isDir: true,
                        children: [
                            { name: 'readme.md', path: '/root/docs/readme.md', isDir: false },
                            { name: 'guide.md', path: '/root/docs/guide.md', isDir: false }
                        ]
                    },
                    { name: 'src', path: '/root/src', isDir: true, children: [] }
                ]
            }]
            store.setFileTree(tree)

            // Remove docs folder
            const removed = store.removeNode('/root/docs')

            expect(removed).toBe(true)
            expect(store.nodes[0].children).toHaveLength(1)
            expect(store.nodes[0].children![0].name).toBe('src')
        })

        it('should deselect files when removing a directory', () => {
            const store = useFileStore()

            const tree: FileNode[] = [{
                name: 'root',
                path: '/root',
                isDir: true,
                children: [
                    {
                        name: 'docs',
                        path: '/root/docs',
                        isDir: true,
                        children: [
                            { name: 'readme.md', path: '/root/docs/readme.md', isDir: false }
                        ]
                    }
                ]
            }]
            store.setFileTree(tree)

            // Select file
            store.selectPath('/root/docs/readme.md')
            expect(store.selectedPaths.has('/root/docs/readme.md')).toBe(true)

            // Remove docs folder
            store.removeNode('/root/docs')

            // File should be deselected
            expect(store.selectedPaths.has('/root/docs/readme.md')).toBe(false)
        })

        it('should return false when node not found', () => {
            const store = useFileStore()

            const tree: FileNode[] = [{
                name: 'root',
                path: '/root',
                isDir: true,
                children: []
            }]
            store.setFileTree(tree)

            const removed = store.removeNode('/nonexistent')

            expect(removed).toBe(false)
        })

        it('should clear caches after removal', () => {
            const store = useFileStore()

            const tree: FileNode[] = [{
                name: 'root',
                path: '/root',
                isDir: true,
                children: [
                    { name: 'file.ts', path: '/root/file.ts', isDir: false }
                ]
            }]
            store.setFileTree(tree)

            // Access getAllFilesInNode to populate cache
            store.getAllFilesInNode(store.nodes[0])

            // Remove file
            store.removeNode('/root/file.ts')

            // Cache should be invalidated - next call should return updated result
            const files = store.getAllFilesInNode(store.nodes[0])
            expect(files).toHaveLength(0)
        })
    })
})
