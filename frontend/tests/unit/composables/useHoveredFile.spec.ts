import { beforeEach, describe, expect, it, vi } from 'vitest'

describe('useHoveredFile', () => {
    beforeEach(() => {
        // Reset module state between tests by re-importing
        vi.resetModules()
    })

    it('should create hovered file state with initial null values', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state } = useHoveredFile()

        expect(state.path).toBeNull()
        expect(state.isDir).toBeNull()
    })

    it('should set hovered file path and isDir', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state, setHovered } = useHoveredFile()

        setHovered('/path/to/file.ts', false)

        expect(state.path).toBe('/path/to/file.ts')
        expect(state.isDir).toBe(false)
    })

    it('should set hovered directory', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state, setHovered } = useHoveredFile()

        setHovered('/path/to/folder', true)

        expect(state.path).toBe('/path/to/folder')
        expect(state.isDir).toBe(true)
    })

    it('should clear hovered state when path matches', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state, setHovered, clearHovered } = useHoveredFile()

        setHovered('/path/to/file.ts', false)
        clearHovered('/path/to/file.ts')

        expect(state.path).toBeNull()
        expect(state.isDir).toBeNull()
    })

    it('should not clear hovered state when path does not match', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state, setHovered, clearHovered } = useHoveredFile()

        setHovered('/path/to/file.ts', false)
        clearHovered('/different/path.ts')

        expect(state.path).toBe('/path/to/file.ts')
        expect(state.isDir).toBe(false)
    })

    it('should handle null values in setHovered', async () => {
        const { useHoveredFile } = await import('@/features/files/composables/useHoveredFile')
        const { state, setHovered } = useHoveredFile()

        setHovered('/path/to/file.ts', false)
        setHovered(null, null)

        expect(state.path).toBeNull()
        expect(state.isDir).toBeNull()
    })
})
