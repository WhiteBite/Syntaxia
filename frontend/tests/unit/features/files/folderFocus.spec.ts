/**
 * Tests for Folder Focus Mode (Task 19)
 * Tests the ability to isolate a folder as temporary root
 */
import { useFileStore } from '@/features/files/model/file.store'
import type { FileNode } from '@/types/domain'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it } from 'vitest'

describe('Folder Focus Mode', () => {
    let fileStore: ReturnType<typeof useFileStore>

    beforeEach(() => {
        setActivePinia(createPinia())
        fileStore = useFileStore()
    })

    const mockTree: FileNode[] = [
        {
            name: 'frontend',
            path: '/project/frontend',
            isDir: true,
            isExpanded: true,
            children: [
                {
                    name: 'src',
                    path: '/project/frontend/src',
                    isDir: true,
                    isExpanded: true,
                    children: [
                        {
                            name: 'features',
                            path: '/project/frontend/src/features',
                            isDir: true,
                            isExpanded: true,
                            children: [
                                {
                                    name: 'auth',
                                    path: '/project/frontend/src/features/auth',
                                    isDir: true,
                                    isExpanded: false,
                                    children: [
                                        {
                                            name: 'login.ts',
                                            path: '/project/frontend/src/features/auth/login.ts',
                                            isDir: false,
                                            size: 1000,
                                        },
                                    ],
                                },
                                {
                                    name: 'files',
                                    path: '/project/frontend/src/features/files',
                                    isDir: true,
                                    isExpanded: false,
                                    children: [
                                        {
                                            name: 'FileTree.vue',
                                            path: '/project/frontend/src/features/files/FileTree.vue',
                                            isDir: false,
                                            size: 2000,
                                        },
                                    ],
                                },
                            ],
                        },
                        {
                            name: 'components',
                            path: '/project/frontend/src/components',
                            isDir: true,
                            isExpanded: false,
                            children: [
                                {
                                    name: 'Button.vue',
                                    path: '/project/frontend/src/components/Button.vue',
                                    isDir: false,
                                    size: 500,
                                },
                            ],
                        },
                    ],
                },
            ],
        },
        {
            name: 'backend',
            path: '/project/backend',
            isDir: true,
            isExpanded: false,
            children: [
                {
                    name: 'main.go',
                    path: '/project/backend/main.go',
                    isDir: false,
                    size: 3000,
                },
            ],
        },
    ]

    it('should set focused folder path', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        expect(fileStore.focusedFolderPath).toBe('/project/frontend/src/features')
    })

    it('should auto-expand focused folder', () => {
        fileStore.setFileTree(mockTree)

        // Collapse the features folder first
        fileStore.collapsePath('/project/frontend/src/features')
        const featuresFolder = fileStore.findNode('/project/frontend/src/features')
        expect(featuresFolder?.isExpanded).toBe(false)

        // Focus should auto-expand
        fileStore.setFocusOnFolder('/project/frontend/src/features')
        expect(featuresFolder?.isExpanded).toBe(true)
    })

    it('should clear focused folder path', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')
        expect(fileStore.focusedFolderPath).toBe('/project/frontend/src/features')

        fileStore.clearFolderFocus()
        expect(fileStore.focusedFolderPath).toBeNull()
    })

    it('should not set focus on a file (only folders)', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features/auth/login.ts')

        // Should not set focus on a file
        expect(fileStore.focusedFolderPath).toBeNull()
    })

    it('should not set focus on non-existent path', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/nonexistent')

        expect(fileStore.focusedFolderPath).toBeNull()
    })

    it('should work with selection in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Select files within focused folder
        fileStore.selectPath('/project/frontend/src/features/auth/login.ts')
        fileStore.selectPath('/project/frontend/src/features/files/FileTree.vue')

        expect(fileStore.selectedCount).toBe(2)
        expect(fileStore.selectedPaths.has('/project/frontend/src/features/auth/login.ts')).toBe(true)
        expect(fileStore.selectedPaths.has('/project/frontend/src/features/files/FileTree.vue')).toBe(true)
    })

    it('should work with search in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Search should still work
        fileStore.setSearchQuery('login')
        expect(fileStore.searchQuery).toBe('login')
    })

    it('should allow switching focus between folders', () => {
        fileStore.setFileTree(mockTree)

        // Focus on features
        fileStore.setFocusOnFolder('/project/frontend/src/features')
        expect(fileStore.focusedFolderPath).toBe('/project/frontend/src/features')

        // Switch focus to components
        fileStore.setFocusOnFolder('/project/frontend/src/components')
        expect(fileStore.focusedFolderPath).toBe('/project/frontend/src/components')
    })

    it('should maintain selection when clearing focus', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Select files
        fileStore.selectPath('/project/frontend/src/features/auth/login.ts')
        expect(fileStore.selectedCount).toBe(1)

        // Clear focus
        fileStore.clearFolderFocus()

        // Selection should be maintained
        expect(fileStore.selectedCount).toBe(1)
        expect(fileStore.selectedPaths.has('/project/frontend/src/features/auth/login.ts')).toBe(true)
    })

    it('should work with expand/collapse in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Expand auth folder
        fileStore.expandPath('/project/frontend/src/features/auth')
        const authFolder = fileStore.findNode('/project/frontend/src/features/auth')
        expect(authFolder?.isExpanded).toBe(true)

        // Collapse auth folder
        fileStore.collapsePath('/project/frontend/src/features/auth')
        expect(authFolder?.isExpanded).toBe(false)
    })

    it('should work with recursive selection in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Select all files in auth folder
        fileStore.selectRecursive('/project/frontend/src/features/auth')

        expect(fileStore.selectedCount).toBe(1)
        expect(fileStore.selectedPaths.has('/project/frontend/src/features/auth/login.ts')).toBe(true)
    })

    it('should work with zen mode in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Enable zen mode
        fileStore.toggleZenMode()
        expect(fileStore.isZenMode).toBe(true)

        // Select a file
        fileStore.selectPath('/project/frontend/src/features/auth/login.ts')

        // Check if dimming works correctly
        const loginFile = fileStore.findNode('/project/frontend/src/features/auth/login.ts')
        const fileTreeFile = fileStore.findNode('/project/frontend/src/features/files/FileTree.vue')

        expect(fileStore.shouldDimNode(loginFile!)).toBe(false) // Selected, should not dim
        expect(fileStore.shouldDimNode(fileTreeFile!)).toBe(true) // Not selected, should dim
    })

    it('should work with selected-only mode in focused mode', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features')

        // Select files
        fileStore.selectPath('/project/frontend/src/features/auth/login.ts')
        fileStore.selectPath('/project/frontend/src/features/files/FileTree.vue')

        // Enable selected-only mode
        fileStore.toggleSelectedOnlyMode()
        expect(fileStore.isSelectedOnlyMode).toBe(true)

        // Filtered nodes should only show selected files
        expect(fileStore.filteredNodes.length).toBe(2)
    })

    it('should handle focus on root-level folder', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend')

        expect(fileStore.focusedFolderPath).toBe('/project/frontend')

        // Should be able to access nested folders
        const srcFolder = fileStore.findNode('/project/frontend/src')
        expect(srcFolder).toBeDefined()
    })

    it('should handle focus on deeply nested folder', () => {
        fileStore.setFileTree(mockTree)
        fileStore.setFocusOnFolder('/project/frontend/src/features/auth')

        expect(fileStore.focusedFolderPath).toBe('/project/frontend/src/features/auth')

        // Should be able to access files in focused folder
        const loginFile = fileStore.findNode('/project/frontend/src/features/auth/login.ts')
        expect(loginFile).toBeDefined()
    })
})
