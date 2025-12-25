import { useQuickFilters } from '@/features/files/composables/useQuickFilters'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'; // Used in mocks

// Mock stores
const mockSetFilterExtensions = vi.fn()
const mockNodesData = [
    {
        name: 'src', isDir: true, children: [
            { name: 'app.ts', isDir: false },
            { name: 'utils.ts', isDir: false },
            { name: 'styles.css', isDir: false },
        ]
    },
    { name: 'README.md', isDir: false },
    { name: 'package.json', isDir: false },
]

vi.mock('@/features/files/model/file.store', () => ({
    useFileStore: () => ({
        nodes: mockNodesData,
        setFilterExtensions: mockSetFilterExtensions,
    }),
}))

vi.mock('@/stores/settings.store', () => ({
    useSettingsStore: () => ({
        settings: {
            fileExplorer: {
                quickFilters: [
                    { id: 'source', label: 'Source', extensions: ['.ts', '.js'] },
                    { id: 'config', label: 'Config', extensions: ['.json', '.yaml'] },
                    { id: 'docs', label: 'Docs', extensions: ['.md', '.txt'] },
                ],
            },
        },
        resetToDefaults: vi.fn(),
    }),
}))

vi.mock('@/stores/project.store', () => ({
    useProjectStore: () => ({
        currentPath: '/test/project',
    }),
}))

vi.mock('@/composables/useI18n', () => ({
    useI18n: () => ({
        t: (key: string) => key,
    }),
}))

// Mock sub-composables
vi.mock('@/features/files/composables/useFilterDropdown', () => ({
    useFilterDropdown: () => ({
        openDropdown: ref(null),
        dropdownStyle: ref({}),
        dropdownRefs: {},
        toggleDropdown: vi.fn(),
        closeDropdown: vi.fn(),
        setupListeners: vi.fn(),
        cleanupListeners: vi.fn(),
    }),
}))

vi.mock('@/features/files/composables/useFilterPersistence', () => ({
    useFilterPersistence: () => ({
        loadState: vi.fn(),
        clearState: vi.fn(),
    }),
}))

vi.mock('@/features/files/composables/useSmartFilters', () => ({
    useSmartFilters: () => ({
        languageFilters: ref([
            { id: 'lang-ts', label: 'TypeScript', extensions: ['.ts', '.tsx'], category: 'lang', shortLabel: 'TS' },
        ]),
        smartFilters: ref([
            { id: 'smart-vue', label: 'Vue', extensions: ['.vue'], category: 'smart', shortLabel: 'Vue' },
        ]),
        isLoading: ref(false),
        loadProjectStructure: vi.fn(),
    }),
}))

describe('useQuickFilters', () => {
    beforeEach(() => {
        vi.clearAllMocks()
    })

    it('should initialize with empty filter state', () => {
        const filters = useQuickFilters()

        expect(filters.allActiveFilters.value).toHaveLength(0)
        expect(filters.allExcludedFilters.value).toHaveLength(0)
    })

    it('should have type filters from settings', () => {
        const filters = useQuickFilters()

        expect(filters.typeFilters.value).toHaveLength(3)
        expect(filters.typeFilters.value[0].id).toBe('source')
    })

    it('should toggle filter on click', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const event = new MouseEvent('click')

        filters.toggleFilter(sourceFilter, event)

        expect(filters.isFilterActive('source')).toBe(true)
        expect(filters.allActiveFilters.value).toHaveLength(1)
        expect(mockSetFilterExtensions).toHaveBeenCalled()
    })

    it('should deselect filter when clicking active filter', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const event = new MouseEvent('click')

        // Activate
        filters.toggleFilter(sourceFilter, event)
        expect(filters.isFilterActive('source')).toBe(true)

        // Deactivate
        filters.toggleFilter(sourceFilter, event)
        expect(filters.isFilterActive('source')).toBe(false)
    })

    it('should support multi-select with Ctrl key', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const configFilter = filters.typeFilters.value[1]

        // First filter - normal click
        filters.toggleFilter(sourceFilter, new MouseEvent('click'))
        expect(filters.allActiveFilters.value).toHaveLength(1)

        // Second filter - Ctrl+click
        filters.toggleFilter(configFilter, new MouseEvent('click', { ctrlKey: true }))
        expect(filters.allActiveFilters.value).toHaveLength(2)
        expect(filters.isFilterActive('source')).toBe(true)
        expect(filters.isFilterActive('config')).toBe(true)
    })

    it('should exclude filter with Shift key', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const event = new MouseEvent('click', { shiftKey: true })

        filters.toggleFilter(sourceFilter, event)

        expect(filters.isFilterExcluded('source')).toBe(true)
        expect(filters.isFilterActive('source')).toBe(false)
        expect(filters.allExcludedFilters.value).toHaveLength(1)
    })

    it('should remove filter from active and excluded', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]

        // Activate
        filters.toggleFilter(sourceFilter, new MouseEvent('click'))
        expect(filters.isFilterActive('source')).toBe(true)

        // Remove
        filters.removeFilter(sourceFilter)
        expect(filters.isFilterActive('source')).toBe(false)
        expect(filters.isFilterExcluded('source')).toBe(false)
    })

    it('should clear type filters', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const configFilter = filters.typeFilters.value[1]

        filters.toggleFilter(sourceFilter, new MouseEvent('click', { ctrlKey: true }))
        filters.toggleFilter(configFilter, new MouseEvent('click', { ctrlKey: true }))
        expect(filters.activeTypeFilters.value).toHaveLength(2)

        filters.clearTypeFilters()
        expect(filters.activeTypeFilters.value).toHaveLength(0)
    })

    it('should clear all filters', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]

        filters.toggleFilter(sourceFilter, new MouseEvent('click'))
        expect(filters.allActiveFilters.value.length).toBeGreaterThan(0)

        filters.clearAllFilters()
        expect(filters.allActiveFilters.value).toHaveLength(0)
        expect(filters.allExcludedFilters.value).toHaveLength(0)
    })

    it('should calculate total files count', () => {
        const filters = useQuickFilters()

        // 3 files in src + README.md + package.json = 5
        expect(filters.totalFiles.value).toBe(5)
    })

    it('should get filter count for specific filter', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0] // .ts, .js

        const count = filters.getFilterCount(sourceFilter)
        // app.ts and utils.ts match .ts extension
        expect(count).toBe(2)
    })

    it('should get filter percentage', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]

        const percentage = filters.getFilterPercentage(sourceFilter)
        // 2 out of 5 files = 40%
        expect(percentage).toBe(40)
    })

    it('should have hasActiveTypeFilters flag', () => {
        const filters = useQuickFilters()

        expect(filters.hasActiveTypeFilters.value).toBe(false)

        const sourceFilter = filters.typeFilters.value[0]
        filters.toggleFilter(sourceFilter, new MouseEvent('click'))

        expect(filters.hasActiveTypeFilters.value).toBe(true)
    })

    it('should replace single filter on normal click', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]
        const configFilter = filters.typeFilters.value[1]

        // First filter
        filters.toggleFilter(sourceFilter, new MouseEvent('click'))
        expect(filters.isFilterActive('source')).toBe(true)

        // Second filter without Ctrl - should replace
        filters.toggleFilter(configFilter, new MouseEvent('click'))
        expect(filters.isFilterActive('source')).toBe(false)
        expect(filters.isFilterActive('config')).toBe(true)
        expect(filters.allActiveFilters.value).toHaveLength(1)
    })

    it('should toggle exclude off when clicking excluded filter with Shift', () => {
        const filters = useQuickFilters()
        const sourceFilter = filters.typeFilters.value[0]

        // Exclude
        filters.toggleFilter(sourceFilter, new MouseEvent('click', { shiftKey: true }))
        expect(filters.isFilterExcluded('source')).toBe(true)

        // Un-exclude
        filters.toggleFilter(sourceFilter, new MouseEvent('click', { shiftKey: true }))
        expect(filters.isFilterExcluded('source')).toBe(false)
    })
})
