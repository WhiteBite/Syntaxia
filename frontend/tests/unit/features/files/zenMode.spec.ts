import { useFileStore } from '@/features/files/model/file.store'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Zen Mode', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('should toggle zen mode on and off', () => {
        const fileStore = useFileStore()

        expect(fileStore.isZenMode).toBe(false)

        fileStore.toggleZenMode()
        expect(fileStore.isZenMode).toBe(true)

        fileStore.toggleZenMode()
        expect(fileStore.isZenMode).toBe(false)
    })

    it('should not dim nodes when zen mode is off', () => {
        const fileStore = useFileStore()

        const file: FileNode = {
            name: 'test.ts',
            path: '/project/test.ts',
            isDir: false,
            isExpanded: false,
            size: 1000,
            contentType: 'text',
        }

        expect(fileStore.shouldDimNode(file)).toBe(false)
    })

    it('should not dim selected files in zen mode', () => {
        const fileStore = useFileStore()

        const file: FileNode = {
            name: 'test.ts',
            path: '/project/test.ts',
            isDir: false,
            isExpanded: false,
            size: 1000,
            contentType: 'text',
        }

        // Select the file
        fileStore.selectPath(file.path)

        // Enable zen mode
        fileStore.toggleZenMode()

        // Selected file should not be dimmed
        expect(fileStore.shouldDimNode(file)).toBe(false)
    })

    it('should dim unselected files in zen mode', () => {
        const fileStore = useFileStore()

        const file: FileNode = {
            name: 'test.ts',
            path: '/project/test.ts',
            isDir: false,
            isExpanded: false,
            size: 1000,
            contentType: 'text',
        }

        // Enable zen mode without selecting the file
        fileStore.toggleZenMode()

        // Unselected file should be dimmed
        expect(fileStore.shouldDimNode(file)).toBe(true)
    })

    it('should not dim folders containing selected files', () => {
        const fileStore = useFileStore()

        const folder: FileNode = {
            name: 'src',
            path: '/project/src',
            isDir: true,
            isExpanded: true,
            children: [
                {
                    name: 'test.ts',
                    path: '/project/src/test.ts',
                    isDir: false,
                    isExpanded: false,
                    size: 1000,
                    contentType: 'text',
                },
            ],
        }

        // Set up the tree
        fileStore.setFileTree([folder])

        // Select a file in the folder
        fileStore.selectPath('/project/src/test.ts')

        // Enable zen mode
        fileStore.toggleZenMode()

        // Folder containing selected file should not be dimmed
        expect(fileStore.shouldDimNode(folder)).toBe(false)
    })

    it('should dim folders without selected files', () => {
        const fileStore = useFileStore()

        const folder: FileNode = {
            name: 'src',
            path: '/project/src',
            isDir: true,
            isExpanded: true,
            children: [
                {
                    name: 'test.ts',
                    path: '/project/src/test.ts',
                    isDir: false,
                    isExpanded: false,
                    size: 1000,
                    contentType: 'text',
                },
            ],
        }

        // Set up the tree
        fileStore.setFileTree([folder])

        // Enable zen mode without selecting any files
        fileStore.toggleZenMode()

        // Folder without selected files should be dimmed
        expect(fileStore.shouldDimNode(folder)).toBe(true)
    })

    it('should handle nested folders correctly', () => {
        const fileStore = useFileStore()

        const rootFolder: FileNode = {
            name: 'src',
            path: '/project/src',
            isDir: true,
            isExpanded: true,
            children: [
                {
                    name: 'components',
                    path: '/project/src/components',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'Button.vue',
                            path: '/project/src/components/Button.vue',
                            isDir: false,
                            isExpanded: false,
                            size: 1000,
                            contentType: 'text',
                        },
                    ],
                },
            ],
        }

        // Set up the tree
        fileStore.setFileTree([rootFolder])

        // Select a deeply nested file
        fileStore.selectPath('/project/src/components/Button.vue')

        // Enable zen mode
        fileStore.toggleZenMode()

        // All parent folders should not be dimmed
        expect(fileStore.shouldDimNode(rootFolder)).toBe(false)
        expect(fileStore.shouldDimNode(rootFolder.children![0])).toBe(false)

        // The selected file should not be dimmed
        expect(fileStore.shouldDimNode(rootFolder.children![0].children![0])).toBe(false)
    })
})
