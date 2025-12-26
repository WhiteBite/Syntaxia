import { useFileSelection } from '@/composables/useFileSelection'
import type { FileNode } from '@/features/files/model/file.store'
import { beforeEach, describe, expect, it, vi } from 'vitest'

describe('useFileSelection', () => {
    // Mock data
    const mockTree: FileNode[] = [
        {
            name: 'src',
            path: '/project/src',
            isDir: true,
            children: [
                { name: 'index.ts', path: '/project/src/index.ts', isDir: false },
                { name: 'utils.ts', path: '/project/src/utils.ts', isDir: false },
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

    // Mock functions
    const findNode = vi.fn((path: string): FileNode | null => {
        const search = (nodes: FileNode[]): FileNode | null => {
            for (const node of nodes) {
                if (node.path === path) return node
                if (node.children) {
                    const found = search(node.children)
                    if (found) return found
                }
            }
            return null
        }
        return search(mockTree)
    })

    const getAllFilesInNode = vi.fn((node: FileNode): string[] => {
        const files: string[] = []
        const collect = (n: FileNode) => {
            if (!n.isDir) {
                files.push(n.path)
            } else if (n.children) {
                n.children.forEach(collect)
            }
        }
        collect(node)
        return files
    })

    let selection: ReturnType<typeof useFileSelection>

    beforeEach(() => {
        vi.clearAllMocks()
        selection = useFileSelection({ findNode, getAllFilesInNode })
    })

    it('toggleSelect добавляет файл в selection', () => {
        expect(selection.selectedPaths.value.size).toBe(0)

        selection.toggleSelect('/project/README.md')

        expect(selection.selectedPaths.value.has('/project/README.md')).toBe(true)
        expect(selection.selectedCount.value).toBe(1)
    })

    it('toggleSelect удаляет файл из selection', () => {
        // Add file first
        selection.toggleSelect('/project/README.md')
        expect(selection.selectedPaths.value.has('/project/README.md')).toBe(true)

        // Toggle again to remove
        selection.toggleSelect('/project/README.md')

        expect(selection.selectedPaths.value.has('/project/README.md')).toBe(false)
        expect(selection.selectedCount.value).toBe(0)
    })

    it('selectRecursive выбирает все файлы в папке', () => {
        selection.selectRecursive('/project/src')

        expect(selection.selectedPaths.value.has('/project/src/index.ts')).toBe(true)
        expect(selection.selectedPaths.value.has('/project/src/utils.ts')).toBe(true)
        expect(selection.selectedPaths.value.has('/project/src/components/Button.vue')).toBe(true)
        expect(selection.selectedCount.value).toBe(3)
    })

    it('clearSelection очищает выбор', () => {
        // Add some files
        selection.selectRecursive('/project/src')
        expect(selection.selectedCount.value).toBe(3)

        // Clear
        selection.clearSelection()

        expect(selection.selectedCount.value).toBe(0)
        expect(selection.hasSelectedFiles.value).toBe(false)
    })
})
