import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Create mutable mock state
let mockHasContext = true

// Mock dependencies - must be before imports
vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn()
    })
}))

vi.mock('@/features/context/model/context.store', () => ({
    useContextStore: () => ({
        get hasContext() { return mockHasContext },
        getFullContextContent: vi.fn().mockResolvedValue('test context content')
    })
}))

vi.mock('@/services/api.service', () => ({
    apiService: {
        exportContext: vi.fn().mockResolvedValue({
            mode: 'clipboard',
            text: 'exported content'
        })
    }
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        projectPath: '/test/project'
    })
}))

vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        settings: {
            context: {
                outputFormat: 'xml',
                stripComments: false,
                includeManifest: true,
                includeFileTree: true,
                maxTokens: 100000,
                enableAutoSplit: false,
                maxTokensPerChunk: 50000,
                includeLineNumbers: true,
                splitStrategy: 'file'
            },
            aiModel: 'gpt-4'
        }
    })
}))

// Import after mocks
import { useExport } from '@/composables/useExport'

// Mock localStorage
const localStorageMock = {
    store: {} as Record<string, string>,
    getItem: vi.fn((key: string) => localStorageMock.store[key] || null),
    setItem: vi.fn((key: string, value: string) => {
        localStorageMock.store[key] = value
    }),
    clear: vi.fn(() => {
        localStorageMock.store = {}
    }),
}
Object.defineProperty(globalThis, 'localStorage', { value: localStorageMock })

// Mock navigator.clipboard
Object.defineProperty(navigator, 'clipboard', {
    value: {
        writeText: vi.fn().mockResolvedValue(undefined)
    },
    writable: true
})

describe('useExport', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        localStorageMock.clear()
        mockHasContext = true
        vi.clearAllMocks()
    })

    it('should initialize with clipboard mode', () => {
        const { selectedMode } = useExport()
        expect(selectedMode.value).toBe('clipboard')
    })

    it('should initialize with closed state', () => {
        const { isOpen } = useExport()
        expect(isOpen.value).toBe(false)
    })

    it('should initialize with not exporting state', () => {
        const { isExporting } = useExport()
        expect(isExporting.value).toBe(false)
    })

    it('should have includeFileTree in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('includeFileTree')
        expect(settings.value.includeFileTree).toBe(true)
    })

    it('should have includeManifest in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('includeManifest')
        expect(settings.value.includeManifest).toBe(true)
    })

    it('should have stripComments in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('stripComments')
        expect(settings.value.stripComments).toBe(false)
    })

    it('should have tokenLimit in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('tokenLimit')
        expect(settings.value.tokenLimit).toBe(100000)
    })

    it('should have exportFormat in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('exportFormat')
        expect(settings.value.exportFormat).toBe('manifest')
    })

    it('should have aiProfile in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('aiProfile')
        expect(settings.value.aiProfile).toBe('gpt-4')
    })

    it('should open modal when context exists', () => {
        mockHasContext = true
        const { open, isOpen } = useExport()

        const result = open()

        expect(result).toBe(true)
        expect(isOpen.value).toBe(true)
    })

    it('should close modal', () => {
        mockHasContext = true
        const { open, close, isOpen } = useExport()

        open()
        expect(isOpen.value).toBe(true)

        close()
        expect(isOpen.value).toBe(false)
    })

    it('should clear error and result on open', () => {
        mockHasContext = true
        const { open, error, exportResult } = useExport()

        open()

        expect(error.value).toBeNull()
        expect(exportResult.value).toBeNull()
    })

    it('should have enableAutoSplit in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('enableAutoSplit')
        expect(settings.value.enableAutoSplit).toBe(false)
    })

    it('should have maxTokensPerChunk in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('maxTokensPerChunk')
        expect(settings.value.maxTokensPerChunk).toBe(50000)
    })

    it('should have includeLineNumbers in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('includeLineNumbers')
        expect(settings.value.includeLineNumbers).toBe(true)
    })

    it('should allow changing selected mode', () => {
        const { selectedMode } = useExport()

        selectedMode.value = 'ai'
        expect(selectedMode.value).toBe('ai')

        selectedMode.value = 'human'
        expect(selectedMode.value).toBe('human')

        selectedMode.value = 'clipboard'
        expect(selectedMode.value).toBe('clipboard')
    })

    it('should expose contextStore', () => {
        mockHasContext = true
        const { contextStore } = useExport()
        expect(contextStore).toBeDefined()
        expect(contextStore.hasContext).toBe(true)
    })

    it('should have splitStrategy in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('splitStrategy')
        expect(settings.value.splitStrategy).toBe('file')
    })

    it('should have overlapTokens in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('overlapTokens')
        expect(settings.value.overlapTokens).toBe(200)
    })

    it('should have theme in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('theme')
        expect(settings.value.theme).toBe('default')
    })

    it('should have includePageNumbers in settings', () => {
        const { settings } = useExport()
        expect(settings.value).toHaveProperty('includePageNumbers')
        expect(settings.value.includePageNumbers).toBe(true)
    })

    it('should set error when opening without context', () => {
        mockHasContext = false
        const { open, error, isOpen } = useExport()

        const result = open()

        expect(result).toBe(false)
        expect(isOpen.value).toBe(false)
        expect(error.value).toBeTruthy()
    })
})
