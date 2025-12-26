import { useQuickLook } from '@/features/files/composables/useQuickLook'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock stores
const mockToggleSelect = vi.fn()
const mockToggleExpand = vi.fn()
const mockAddToast = vi.fn()

vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        toggleSelect: mockToggleSelect,
        toggleExpand: mockToggleExpand,
    }),
}))

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: mockAddToast,
    }),
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key,
    }),
}))

// Mock hoveredFile state
let mockHoveredPath: string | null = null
let mockHoveredIsDir: boolean | null = null

vi.mock('@/features/files/composables/useHoveredFile', () => ({
    useHoveredFile: () => ({
        state: {
            get path() { return mockHoveredPath },
            get isDir() { return mockHoveredIsDir },
        },
        setHovered: (path: string | null, isDir: boolean | null) => {
            mockHoveredPath = path
            mockHoveredIsDir = isDir
        },
        clearHovered: () => {
            mockHoveredPath = null
            mockHoveredIsDir = null
        },
    }),
}))

describe('useQuickLook', () => {
    beforeEach(() => {
        vi.clearAllMocks()
        mockHoveredPath = null
        mockHoveredIsDir = null
        setActivePinia(createPinia())
    })

    it('should initialize with closed state', () => {
        const quickLook = useQuickLook()

        expect(quickLook.isVisible.value).toBe(false)
        expect(quickLook.currentPath.value).toBe('')
    })

    it('should open with specified path', () => {
        const quickLook = useQuickLook()

        quickLook.open('/path/to/file.ts')

        expect(quickLook.isVisible.value).toBe(true)
        expect(quickLook.currentPath.value).toBe('/path/to/file.ts')
    })

    it('should close when same path is opened again', () => {
        const quickLook = useQuickLook()

        quickLook.open('/path/to/file.ts')
        quickLook.open('/path/to/file.ts')

        expect(quickLook.isVisible.value).toBe(false)
    })

    it('should switch to different path', () => {
        const quickLook = useQuickLook()

        quickLook.open('/path/to/file1.ts')
        quickLook.open('/path/to/file2.ts')

        expect(quickLook.isVisible.value).toBe(true)
        expect(quickLook.currentPath.value).toBe('/path/to/file2.ts')
    })

    it('should close explicitly', () => {
        const quickLook = useQuickLook()

        quickLook.open('/path/to/file.ts')
        quickLook.close()

        expect(quickLook.isVisible.value).toBe(false)
    })

    it('should add file to context', () => {
        const quickLook = useQuickLook()

        quickLook.addToContext('/path/to/file.ts')

        expect(mockToggleSelect).toHaveBeenCalledWith('/path/to/file.ts')
        expect(mockAddToast).toHaveBeenCalledWith('files.addToContext', 'success')
    })

    it('should handle spacebar for hovered file', () => {
        mockHoveredPath = '/path/to/file.ts'
        mockHoveredIsDir = false

        const quickLook = useQuickLook()
        const handled = quickLook.handleSpacebarPreview()

        expect(handled).toBe(true)
        expect(quickLook.isVisible.value).toBe(true)
    })

    it('should expand directory on spacebar', () => {
        mockHoveredPath = '/path/to/folder'
        mockHoveredIsDir = true

        const quickLook = useQuickLook()
        const handled = quickLook.handleSpacebarPreview()

        expect(handled).toBe(true)
        expect(mockToggleExpand).toHaveBeenCalledWith('/path/to/folder')
    })

    it('should close quicklook on spacebar when nothing hovered', () => {
        mockHoveredPath = null
        mockHoveredIsDir = null

        const quickLook = useQuickLook()
        quickLook.open('/path/to/file.ts')

        const handled = quickLook.handleSpacebarPreview()

        expect(handled).toBe(true)
        expect(quickLook.isVisible.value).toBe(false)
    })

    it('should return false when nothing to do', () => {
        mockHoveredPath = null
        mockHoveredIsDir = null

        const quickLook = useQuickLook()
        const handled = quickLook.handleSpacebarPreview()

        expect(handled).toBe(false)
    })
})
