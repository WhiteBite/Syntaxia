/**
 * Unit tests for useFileExplorer.ts composable
 * Tests the main composable for FileExplorer component
 */
import { useFileExplorer } from '@/features/files/composables/useFileExplorer'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

// Mock dependencies
vi.mock('@/composables/useLogger', () => ({
    useLogger: () => ({
        debug: vi.fn(),
        info: vi.fn(),
        warn: vi.fn(),
        error: vi.fn(),
    }),
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key,
    }),
}))

vi.mock('@/composables/useContextMenu', () => ({
    useContextMenu: () => ({
        show: vi.fn(),
        hide: vi.fn(),
        isVisible: { value: false },
    }),
}))

vi.mock('@/features/files/api/files.api', () => ({
    filesApi: {
        listFiles: vi.fn().mockResolvedValue([]),
        clearCache: vi.fn(),
    },
}))

vi.mock('@/services/api.service', () => ({
    apiService: {
        clearFileTreeCache: vi.fn().mockResolvedValue(undefined),
        getSettings: vi.fn().mockResolvedValue({}),
        saveSettings: vi.fn().mockResolvedValue(undefined),
    },
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        currentPath: '/test/project',
    }),
}))

vi.mock('@/stores/ui.store', () => ({
    useUIStore: () => ({
        addToast: vi.fn(),
    }),
}))

vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        settings: {
            fileExplorer: {
                useGitignore: true,
                useCustomIgnore: true,
                autoSaveSelection: false,
            },
        },
    }),
}))

vi.mock('@/features/context', () => ({
    useContextStore: () => ({
        clearContext: vi.fn(),
    }),
}))

// Mock the file store
const mockFileStore = {
    nodes: [],
    selectedCount: 0,
    hasSelectedFiles: false,
    selectedFilesList: [],
    toggleSelect: vi.fn(),
    toggleExpand: vi.fn(),
    clearSelection: vi.fn(),
    loadFileTree: vi.fn().mockResolvedValue(undefined),
    getExpandedPaths: vi.fn().mockReturnValue([]),
    restoreExpandedPaths: vi.fn(),
    nodeExists: vi.fn().mockReturnValue(true),
    setFilterExtensions: vi.fn(),
    getAvailableExtensions: vi.fn().mockReturnValue(['.ts', '.vue', '.js']),
    getSelectionStats: vi.fn().mockReturnValue({}),
    loadSelectionFromStorage: vi.fn(),
    selectRecursive: vi.fn(),
    deselectRecursive: vi.fn(),
    expandRecursive: vi.fn(),
    collapseRecursive: vi.fn(),
    removeNode: vi.fn(),
    setSearchQuery: vi.fn(),
}

vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => mockFileStore,
}))

// Mock QuickLook composable
vi.mock('@/features/files/composables/useQuickLook', () => ({
    useQuickLook: () => ({
        isVisible: { value: false },
        currentPath: { value: '' },
        open: vi.fn(),
        close: vi.fn(),
        toggle: vi.fn(),
        addToContext: vi.fn(),
        handleSpacebarPreview: vi.fn().mockReturnValue(false),
    }),
}))

// Mock IgnoreRules composable
vi.mock('@/features/files/composables/useIgnoreRules', () => ({
    useIgnoreRules: () => ({
        addToIgnore: vi.fn().mockResolvedValue(true),
        removeFromIgnore: vi.fn().mockResolvedValue(true),
    }),
}))

// Mock FileSearch composable
vi.mock('@/features/files/composables/useFileSearch', () => ({
    useFileSearch: () => ({
        query: { value: '' },
        isSearching: { value: false },
        handleSearch: vi.fn(),
        clear: vi.fn(),
        setQuery: vi.fn(),
        isActive: vi.fn().mockReturnValue(false),
    }),
}))

describe('useFileExplorer', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
        vi.clearAllMocks()
    })

    describe('Initial State', () => {
        it('should return search state', () => {
            const explorer = useFileExplorer()

            expect(explorer.searchQuery).toBeDefined()
            expect(explorer.handleSearch).toBeDefined()
            expect(explorer.clearSearch).toBeDefined()
        })

        it('should return QuickLook state', () => {
            const explorer = useFileExplorer()

            expect(explorer.quickLookVisible).toBeDefined()
            expect(explorer.quickLookPath).toBeDefined()
            expect(explorer.handleQuickLook).toBeDefined()
            expect(explorer.handleAddToContext).toBeDefined()
        })

        it('should return UI state', () => {
            const explorer = useFileExplorer()

            expect(explorer.showSettings).toBeDefined()
            expect(explorer.filterExtensions).toBeDefined()
            expect(explorer.contextMenu).toBeDefined()
        })

        it('should return computed properties', () => {
            const explorer = useFileExplorer()

            expect(explorer.availableExtensions).toBeDefined()
            expect(explorer.hasSelectionHistory).toBeDefined()
            expect(explorer.selectionHistoryCount).toBeDefined()
            expect(explorer.totalFileCount).toBeDefined()
            expect(explorer.selectionProgress).toBeDefined()
        })

        it('should return action handlers', () => {
            const explorer = useFileExplorer()

            expect(explorer.handleToggleSelect).toBeDefined()
            expect(explorer.handleToggleExpand).toBeDefined()
            expect(explorer.handleRefresh).toBeDefined()
            expect(explorer.handleFilterUpdate).toBeDefined()
            expect(explorer.handleSettingsChange).toBeDefined()
            expect(explorer.restorePreviousSelection).toBeDefined()
            expect(explorer.handleContextMenuShow).toBeDefined()
            expect(explorer.handleContextMenuAction).toBeDefined()
        })

        it('should return lifecycle methods', () => {
            const explorer = useFileExplorer()

            expect(explorer.initialize).toBeDefined()
            expect(explorer.cleanup).toBeDefined()
            expect(explorer.setupWatchers).toBeDefined()
        })
    })

    describe('handleToggleSelect', () => {
        it('should call fileStore.toggleSelect with path', () => {
            const explorer = useFileExplorer()
            explorer.handleToggleSelect('/test/file.ts')

            expect(mockFileStore.toggleSelect).toHaveBeenCalledWith('/test/file.ts')
        })
    })

    describe('handleToggleExpand', () => {
        it('should call fileStore.toggleExpand with path', () => {
            const explorer = useFileExplorer()
            explorer.handleToggleExpand('/test/folder')

            expect(mockFileStore.toggleExpand).toHaveBeenCalledWith('/test/folder')
        })
    })

    describe('handleFilterUpdate', () => {
        it('should call fileStore.setFilterExtensions with selected extensions', () => {
            const explorer = useFileExplorer()
            const extensions = ['.ts', '.vue']
            explorer.handleFilterUpdate(extensions)

            expect(mockFileStore.setFilterExtensions).toHaveBeenCalledWith(extensions)
        })
    })

    describe('handleContextMenuAction', () => {
        it('should handle selectAll action for directory', async () => {
            const explorer = useFileExplorer()
            const node = { name: 'src', path: '/test/src', isDir: true }

            await explorer.handleContextMenuAction({ type: 'selectAll', node })

            expect(mockFileStore.selectRecursive).toHaveBeenCalledWith('/test/src')
        })

        it('should handle deselectAll action for directory', async () => {
            const explorer = useFileExplorer()
            const node = { name: 'src', path: '/test/src', isDir: true }

            await explorer.handleContextMenuAction({ type: 'deselectAll', node })

            expect(mockFileStore.deselectRecursive).toHaveBeenCalledWith('/test/src')
        })

        it('should handle expandAll action for directory', async () => {
            const explorer = useFileExplorer()
            const node = { name: 'src', path: '/test/src', isDir: true }

            await explorer.handleContextMenuAction({ type: 'expandAll', node })

            expect(mockFileStore.expandRecursive).toHaveBeenCalledWith('/test/src')
        })

        it('should handle collapseAll action for directory', async () => {
            const explorer = useFileExplorer()
            const node = { name: 'src', path: '/test/src', isDir: true }

            await explorer.handleContextMenuAction({ type: 'collapseAll', node })

            expect(mockFileStore.collapseRecursive).toHaveBeenCalledWith('/test/src')
        })
    })

    describe('Computed Properties', () => {
        it('should compute availableExtensions from fileStore', () => {
            const explorer = useFileExplorer()

            expect(explorer.availableExtensions.value).toEqual(['.ts', '.vue', '.js'])
        })

        it('should compute selectionProgress as 0 when no files', () => {
            const explorer = useFileExplorer()

            expect(explorer.selectionProgress.value).toBe(0)
        })
    })

    describe('Lifecycle', () => {
        it('should add keydown listener on initialize', () => {
            const addEventListenerSpy = vi.spyOn(window, 'addEventListener')
            const explorer = useFileExplorer()

            explorer.initialize()

            expect(addEventListenerSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
        })

        it('should remove keydown listener on cleanup', () => {
            const removeEventListenerSpy = vi.spyOn(window, 'removeEventListener')
            const explorer = useFileExplorer()

            explorer.initialize()
            explorer.cleanup()

            expect(removeEventListenerSpy).toHaveBeenCalledWith('keydown', expect.any(Function))
        })
    })
})
