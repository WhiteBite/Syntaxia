import { useSandboxStore, type SandboxChange } from '@/stores/sandbox.store'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// Mock apiService
vi.mock('@/services/api.service', () => ({
    apiService: {
        getSandboxChanges: vi.fn(),
        applySandboxChanges: vi.fn(),
        discardSandboxChanges: vi.fn(),
        discardSandboxFile: vi.fn(),
        getSandboxDiff: vi.fn(),
        getSandboxAllDiffs: vi.fn()
    }
}))

// Mock project store
vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        currentPath: '/test/project'
    })
}))

// Mock UI store
const mockAddToast = vi.fn()
vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: mockAddToast
    })
}))

import { apiService } from '@/services/api.service'

describe('SandboxStore', () => {
    let store: ReturnType<typeof useSandboxStore>

    beforeEach(() => {
        setActivePinia(createPinia())
        store = useSandboxStore()
        vi.clearAllMocks()
    })

    afterEach(() => {
        vi.restoreAllMocks()
    })

    describe('initial state', () => {
        it('should have empty changes array', () => {
            expect(store.changes).toEqual([])
        })

        it('should have isLoading as false', () => {
            expect(store.isLoading).toBe(false)
        })
    })

    describe('computed properties', () => {
        it('should return hasChanges as false when no changes', () => {
            expect(store.hasChanges).toBe(false)
        })

        it('should return hasChanges as true when changes exist', () => {
            store.changes = [
                { path: 'file.ts', operation: 'modify', diff: '+line' }
            ]
            expect(store.hasChanges).toBe(true)
        })

        it('should return correct changeCount', () => {
            expect(store.changeCount).toBe(0)

            store.changes = [
                { path: 'file1.ts', operation: 'create', diff: '' },
                { path: 'file2.ts', operation: 'modify', diff: '' },
                { path: 'file3.ts', operation: 'delete', diff: '' }
            ]
            expect(store.changeCount).toBe(3)
        })
    })

    describe('refresh', () => {
        it('should load changes from API', async () => {
            const mockChanges: SandboxChange[] = [
                { path: 'src/main.ts', operation: 'modify', diff: '+new line' },
                { path: 'src/utils.ts', operation: 'create', diff: '+content' }
            ]
            vi.mocked(apiService.getSandboxChanges).mockResolvedValue({
                changes: mockChanges
            })

            await store.refresh()

            expect(store.changes).toEqual(mockChanges)
            expect(store.isLoading).toBe(false)
        })

        it('should set isLoading during request', async () => {
            let resolvePromise: (value: { changes: SandboxChange[] }) => void
            const promise = new Promise<{ changes: SandboxChange[] }>((resolve) => {
                resolvePromise = resolve
            })
            vi.mocked(apiService.getSandboxChanges).mockReturnValue(promise)

            const refreshPromise = store.refresh()
            expect(store.isLoading).toBe(true)

            resolvePromise!({ changes: [] })
            await refreshPromise

            expect(store.isLoading).toBe(false)
        })

        it('should handle API error gracefully', async () => {
            vi.mocked(apiService.getSandboxChanges).mockRejectedValue(
                new Error('API Error')
            )

            await store.refresh()

            expect(store.changes).toEqual([])
            expect(store.isLoading).toBe(false)
        })

        it('should handle null changes in response', async () => {
            vi.mocked(apiService.getSandboxChanges).mockResolvedValue({
                changes: null
            })

            await store.refresh()

            expect(store.changes).toEqual([])
        })
    })

    describe('applyAll', () => {
        it('should apply changes and clear list', async () => {
            store.changes = [
                { path: 'file.ts', operation: 'modify', diff: '' }
            ]
            vi.mocked(apiService.applySandboxChanges).mockResolvedValue(undefined)

            const result = await store.applyAll()

            expect(result).toBe(true)
            expect(store.changes).toEqual([])
            expect(mockAddToast).toHaveBeenCalledWith(
                'Changes applied successfully',
                'success'
            )
        })

        it('should return false on error', async () => {
            vi.mocked(apiService.applySandboxChanges).mockRejectedValue(
                new Error('Apply failed')
            )

            const result = await store.applyAll()

            expect(result).toBe(false)
            expect(mockAddToast).toHaveBeenCalledWith(
                'Failed to apply changes',
                'error'
            )
        })

        it('should set isLoading during operation', async () => {
            let resolvePromise: () => void
            const promise = new Promise<void>((resolve) => {
                resolvePromise = resolve
            })
            vi.mocked(apiService.applySandboxChanges).mockReturnValue(promise)

            const applyPromise = store.applyAll()
            expect(store.isLoading).toBe(true)

            resolvePromise!()
            await applyPromise

            expect(store.isLoading).toBe(false)
        })
    })

    describe('discardAll', () => {
        it('should discard all changes', async () => {
            store.changes = [
                { path: 'file1.ts', operation: 'modify', diff: '' },
                { path: 'file2.ts', operation: 'create', diff: '' }
            ]
            vi.mocked(apiService.discardSandboxChanges).mockResolvedValue(undefined)

            await store.discardAll()

            expect(store.changes).toEqual([])
            expect(mockAddToast).toHaveBeenCalledWith('Changes discarded', 'info')
        })

        it('should show error toast on failure', async () => {
            vi.mocked(apiService.discardSandboxChanges).mockRejectedValue(
                new Error('Discard failed')
            )

            await store.discardAll()

            expect(mockAddToast).toHaveBeenCalledWith(
                'Failed to discard changes',
                'error'
            )
        })

        it('should set isLoading during operation', async () => {
            let resolvePromise: () => void
            const promise = new Promise<void>((resolve) => {
                resolvePromise = resolve
            })
            vi.mocked(apiService.discardSandboxChanges).mockReturnValue(promise)

            const discardPromise = store.discardAll()
            expect(store.isLoading).toBe(true)

            resolvePromise!()
            await discardPromise

            expect(store.isLoading).toBe(false)
        })
    })

    describe('discardFile', () => {
        it('should discard single file and remove from list', async () => {
            store.changes = [
                { path: 'file1.ts', operation: 'modify', diff: '' },
                { path: 'file2.ts', operation: 'create', diff: '' },
                { path: 'file3.ts', operation: 'delete', diff: '' }
            ]
            vi.mocked(apiService.discardSandboxFile).mockResolvedValue(undefined)

            await store.discardFile('file2.ts')

            expect(store.changes).toHaveLength(2)
            expect(store.changes.find(c => c.path === 'file2.ts')).toBeUndefined()
        })

        it('should show error toast on failure', async () => {
            vi.mocked(apiService.discardSandboxFile).mockRejectedValue(
                new Error('Discard file failed')
            )

            await store.discardFile('file.ts')

            expect(mockAddToast).toHaveBeenCalledWith(
                'Failed to discard file',
                'error'
            )
        })

        it('should call API with correct path', async () => {
            vi.mocked(apiService.discardSandboxFile).mockResolvedValue(undefined)

            await store.discardFile('src/components/Button.vue')

            expect(apiService.discardSandboxFile).toHaveBeenCalledWith(
                'src/components/Button.vue'
            )
        })
    })

    describe('getDiff', () => {
        it('should return diff for file', async () => {
            const expectedDiff = `@@ -1,3 +1,4 @@
 line1
+new line
 line2`
            vi.mocked(apiService.getSandboxDiff).mockResolvedValue(expectedDiff)

            const diff = await store.getDiff('file.ts')

            expect(diff).toBe(expectedDiff)
            expect(apiService.getSandboxDiff).toHaveBeenCalledWith('file.ts')
        })

        it('should return empty string on error', async () => {
            vi.mocked(apiService.getSandboxDiff).mockRejectedValue(
                new Error('Get diff failed')
            )

            const diff = await store.getDiff('file.ts')

            expect(diff).toBe('')
        })
    })

    describe('getAllDiffs', () => {
        it('should return all diffs', async () => {
            const expectedDiffs = `diff --git a/file1.ts
+content1
diff --git a/file2.ts
+content2`
            vi.mocked(apiService.getSandboxAllDiffs).mockResolvedValue(expectedDiffs)

            const diffs = await store.getAllDiffs()

            expect(diffs).toBe(expectedDiffs)
        })

        it('should return empty string on error', async () => {
            vi.mocked(apiService.getSandboxAllDiffs).mockRejectedValue(
                new Error('Get all diffs failed')
            )

            const diffs = await store.getAllDiffs()

            expect(diffs).toBe('')
        })
    })

    describe('SandboxChange type', () => {
        it('should correctly type SandboxChange interface', () => {
            const change: SandboxChange = {
                path: 'src/main.ts',
                operation: 'modify',
                diff: '+new line\n-old line'
            }

            expect(change.path).toBe('src/main.ts')
            expect(change.operation).toBe('modify')
            expect(change.diff).toContain('+new line')
        })

        it('should support all operation types', () => {
            const createChange: SandboxChange = {
                path: 'new.ts',
                operation: 'create',
                diff: '+content'
            }
            const modifyChange: SandboxChange = {
                path: 'existing.ts',
                operation: 'modify',
                diff: '+new\n-old'
            }
            const deleteChange: SandboxChange = {
                path: 'removed.ts',
                operation: 'delete',
                diff: '-content'
            }

            expect(createChange.operation).toBe('create')
            expect(modifyChange.operation).toBe('modify')
            expect(deleteChange.operation).toBe('delete')
        })
    })
})
