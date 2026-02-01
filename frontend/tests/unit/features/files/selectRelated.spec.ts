/**
 * Tests for selectRelated functionality (Task 12)
 * Validates smart file selection based on naming patterns
 */

import { useDependencyGraph } from '@/composables/useDependencyGraph'
import type { FileNode } from '@/types/domain'
import { describe, expect, it } from 'vitest'

describe('selectRelated (Task 12: Smart Selection)', () => {
    const { findRelatedFiles } = useDependencyGraph()

    describe('findRelatedFiles', () => {
        it('should find related files with same basename in same directory', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/UserComponent.vue', name: 'UserComponent.vue', isDir: false, size: 1000 },
                { path: 'src/components/user.store.ts', name: 'user.store.ts', isDir: false, size: 500 },
                { path: 'src/components/user.types.ts', name: 'user.types.ts', isDir: false, size: 200 },
                { path: 'src/components/user.css', name: 'user.css', isDir: false, size: 300 },
            ]

            const related = findRelatedFiles('src/components/UserComponent.vue', allNodes)

            expect(related).toContain('src/components/user.store.ts')
            expect(related).toContain('src/components/user.types.ts')
            expect(related).toContain('src/components/user.css')
            expect(related).not.toContain('src/components/UserComponent.vue')
        })

        it('should find test files for TypeScript files', () => {
            const allNodes: FileNode[] = [
                { path: 'src/utils/helper.ts', name: 'helper.ts', isDir: false, size: 1000 },
                { path: 'src/utils/helper.test.ts', name: 'helper.test.ts', isDir: false, size: 500 },
                { path: 'src/utils/helper.spec.ts', name: 'helper.spec.ts', isDir: false, size: 500 },
            ]

            const related = findRelatedFiles('src/utils/helper.ts', allNodes)

            expect(related).toContain('src/utils/helper.test.ts')
            expect(related).toContain('src/utils/helper.spec.ts')
        })

        it('should find styles for Vue components', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/Button.vue', name: 'Button.vue', isDir: false, size: 1000 },
                { path: 'src/components/button.css', name: 'button.css', isDir: false, size: 300 },
                { path: 'src/components/button.scss', name: 'button.scss', isDir: false, size: 400 },
            ]

            const related = findRelatedFiles('src/components/Button.vue', allNodes)

            expect(related).toContain('src/components/button.css')
            expect(related).toContain('src/components/button.scss')
        })

        it('should return empty array when no related files found', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/UserComponent.vue', name: 'UserComponent.vue', isDir: false, size: 1000 },
                { path: 'src/components/OtherComponent.vue', name: 'OtherComponent.vue', isDir: false, size: 1000 },
            ]

            const related = findRelatedFiles('src/components/UserComponent.vue', allNodes)

            expect(related).toEqual([])
        })

        it('should not include the source file itself', () => {
            const allNodes: FileNode[] = [
                { path: 'src/utils/helper.ts', name: 'helper.ts', isDir: false, size: 1000 },
                { path: 'src/utils/helper.test.ts', name: 'helper.test.ts', isDir: false, size: 500 },
            ]

            const related = findRelatedFiles('src/utils/helper.ts', allNodes)

            expect(related).not.toContain('src/utils/helper.ts')
        })

        it('should not include directories', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/UserComponent.vue', name: 'UserComponent.vue', isDir: false, size: 1000 },
                { path: 'src/components/user', name: 'user', isDir: true },
                { path: 'src/components/user.store.ts', name: 'user.store.ts', isDir: false, size: 500 },
            ]

            const related = findRelatedFiles('src/components/UserComponent.vue', allNodes)

            expect(related).toContain('src/components/user.store.ts')
            expect(related).not.toContain('src/components/user')
        })

        it('should handle case-insensitive matching', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/UserComponent.vue', name: 'UserComponent.vue', isDir: false, size: 1000 },
                { path: 'src/components/usercomponent.css', name: 'usercomponent.css', isDir: false, size: 300 },
            ]

            const related = findRelatedFiles('src/components/UserComponent.vue', allNodes)

            expect(related).toContain('src/components/usercomponent.css')
        })

        it('should only find files in same directory', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/UserComponent.vue', name: 'UserComponent.vue', isDir: false, size: 1000 },
                { path: 'src/components/user.store.ts', name: 'user.store.ts', isDir: false, size: 500 },
                { path: 'src/stores/user.store.ts', name: 'user.store.ts', isDir: false, size: 500 },
            ]

            const related = findRelatedFiles('src/components/UserComponent.vue', allNodes)

            expect(related).toContain('src/components/user.store.ts')
            expect(related).not.toContain('src/stores/user.store.ts')
        })

        it('should return unique paths (no duplicates)', () => {
            const allNodes: FileNode[] = [
                { path: 'src/components/user.vue', name: 'user.vue', isDir: false, size: 1000 },
                { path: 'src/components/user.css', name: 'user.css', isDir: false, size: 300 },
            ]

            const related = findRelatedFiles('src/components/user.vue', allNodes)

            expect(related).toEqual(['src/components/user.css'])
            expect(new Set(related).size).toBe(related.length)
        })
    })
})
