/**
 * Main composable for Quick Filters feature
 * Composes smaller composables for dropdown, persistence, and smart filters
 */
import { useI18n } from '@/composables/useI18n'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore, type QuickFilterConfig } from '@/stores/settings.store'
import { computed, ref, watch } from 'vue'
import { useFileStore } from '../model/file.store'
import type { AnyFilter, FilterCategory, FilterState, TypeFilter } from '../model/types'
import { useFilterDropdown } from './useFilterDropdown'
import { useFilterPersistence } from './useFilterPersistence'
import { useSmartFilters } from './useSmartFilters'

interface FileNode {
    name: string
    isDir: boolean
    children?: FileNode[]
}

const CATEGORY_MAP: Record<string, FilterCategory> = {
    source: 'code',
    tests: 'test',
    config: 'config',
    docs: 'docs',
    styles: 'styles',
}

// i18n keys for filter labels
const LABEL_I18N_MAP: Record<string, string> = {
    source: 'filterLabels.source',
    tests: 'filterLabels.tests',
    config: 'filterLabels.config',
    docs: 'filterLabels.docs',
    styles: 'filterLabels.styles',
}

export function useQuickFilters() {
    const { t } = useI18n()
    const fileStore = useFileStore()
    const settingsStore = useSettingsStore()
    const projectStore = useProjectStore()

    // Compose smaller composables
    const dropdown = useFilterDropdown()
    const smartFiltersComposable = useSmartFilters()
    const { languageFilters, smartFilters, isLoading, loadProjectStructure } = smartFiltersComposable

    // Filter state
    const filterState = ref<FilterState>({
        active: new Set<string>(),
        excluded: new Set<string>(),
    })

    const persistence = useFilterPersistence(filterState)

    // UI state
    const showSettingsModal = ref(false)

    // Type filters from settings with i18n labels
    const typeFilters = computed<TypeFilter[]>(() =>
        settingsStore.settings.fileExplorer.quickFilters.map((f) => ({
            ...f,
            category: CATEGORY_MAP[f.id] || 'code',
            shortLabel: LABEL_I18N_MAP[f.id] ? t(LABEL_I18N_MAP[f.id]) : f.label,
        }))
    )

    const editableFilters = computed(() => settingsStore.settings.fileExplorer.quickFilters)

    // Active filters by category
    const activeTypeFilters = computed(() =>
        typeFilters.value.filter((f) => filterState.value.active.has(f.id))
    )
    const activeLanguageFilters = computed(() =>
        languageFilters.value.filter((f) => filterState.value.active.has(f.id))
    )
    const activeSmartFilters = computed(() =>
        smartFilters.value.filter((f) => filterState.value.active.has(f.id))
    )
    const allActiveFilters = computed<AnyFilter[]>(() => [
        ...activeTypeFilters.value,
        ...activeLanguageFilters.value,
        ...activeSmartFilters.value,
    ])

    // Excluded filters by category
    const excludedTypeFilters = computed(() =>
        typeFilters.value.filter((f) => filterState.value.excluded.has(f.id))
    )
    const excludedLanguageFilters = computed(() =>
        languageFilters.value.filter((f) => filterState.value.excluded.has(f.id))
    )
    const excludedSmartFilters = computed(() =>
        smartFilters.value.filter((f) => filterState.value.excluded.has(f.id))
    )
    const allExcludedFilters = computed<AnyFilter[]>(() => [
        ...excludedTypeFilters.value,
        ...excludedLanguageFilters.value,
        ...excludedSmartFilters.value,
    ])

    // Has active flags
    const hasActiveTypeFilters = computed(
        () => activeTypeFilters.value.length > 0 || excludedTypeFilters.value.length > 0
    )
    const hasActiveLanguageFilters = computed(
        () => activeLanguageFilters.value.length > 0 || excludedLanguageFilters.value.length > 0
    )
    const hasActiveSmartFilters = computed(
        () => activeSmartFilters.value.length > 0 || excludedSmartFilters.value.length > 0
    )

    // Total files count
    const totalFiles = computed(() => {
        let count = 0
        const countFiles = (nodes: FileNode[]) => {
            nodes.forEach((node) => {
                if (!node.isDir) count++
                if (node.children) countFiles(node.children)
            })
        }
        countFiles(fileStore.nodes as FileNode[])
        return count
    })


    // Filter management
    function toggleFilter(filter: AnyFilter, event: MouseEvent): void {
        const isMulti = event.ctrlKey || event.metaKey
        const isExclude = event.shiftKey

        if (isExclude) {
            if (filterState.value.excluded.has(filter.id)) {
                filterState.value.excluded.delete(filter.id)
            } else {
                filterState.value.active.delete(filter.id)
                filterState.value.excluded.add(filter.id)
            }
        } else if (isMulti) {
            filterState.value.excluded.delete(filter.id)
            if (filterState.value.active.has(filter.id)) {
                filterState.value.active.delete(filter.id)
            } else {
                filterState.value.active.add(filter.id)
            }
        } else {
            filterState.value.excluded.clear()
            if (filterState.value.active.has(filter.id) && filterState.value.active.size === 1) {
                filterState.value.active.clear()
            } else {
                filterState.value.active.clear()
                filterState.value.active.add(filter.id)
            }
        }

        applyFilters()
        if (!isMulti && !isExclude) dropdown.closeDropdown()
    }

    function removeFilter(filter: AnyFilter): void {
        filterState.value.active.delete(filter.id)
        filterState.value.excluded.delete(filter.id)
        applyFilters()
    }

    function isFilterActive(id: string): boolean {
        return filterState.value.active.has(id)
    }

    function isFilterExcluded(id: string): boolean {
        return filterState.value.excluded.has(id)
    }

    function clearTypeFilters(): void {
        typeFilters.value.forEach((f) => {
            filterState.value.active.delete(f.id)
            filterState.value.excluded.delete(f.id)
        })
        applyFilters()
    }

    function clearLanguageFilters(): void {
        languageFilters.value.forEach((f) => {
            filterState.value.active.delete(f.id)
            filterState.value.excluded.delete(f.id)
        })
        applyFilters()
    }

    function clearSmartFilters(): void {
        smartFilters.value.forEach((f) => {
            filterState.value.active.delete(f.id)
            filterState.value.excluded.delete(f.id)
        })
        applyFilters()
    }

    function clearAllFilters(): void {
        filterState.value.active.clear()
        filterState.value.excluded.clear()
        applyFilters()
    }

    // Apply filters to file store
    function applyFilters(): void {
        const includeExts: string[] = []
        const excludeExts: string[] = []

        activeTypeFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions))
        activeLanguageFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions))
        activeSmartFilters.value.forEach((f) => f.extensions && includeExts.push(...f.extensions))

        excludedTypeFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions))
        excludedLanguageFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions))
        excludedSmartFilters.value.forEach((f) => f.extensions && excludeExts.push(...f.extensions))

        fileStore.setFilterExtensions([...new Set(includeExts)], [...new Set(excludeExts)])
    }

    // Filter info helpers
    function getFilterCount(filter: QuickFilterConfig | AnyFilter): number {
        let count = 0
        const countFiles = (nodes: FileNode[]) => {
            nodes.forEach((node) => {
                if (!node.isDir && filter.extensions?.some((e) => node.name.endsWith(e))) count++
                if (node.children) countFiles(node.children)
            })
        }
        countFiles(fileStore.nodes as FileNode[])
        return count
    }

    function getFilterPercentage(filter: QuickFilterConfig | AnyFilter): number {
        return totalFiles.value === 0
            ? 0
            : Math.min(100, (getFilterCount(filter) / totalFiles.value) * 100)
    }

    // Settings helpers
    function updateFilterExtensions(filter: QuickFilterConfig, value: string): void {
        filter.extensions = value
            .split(',')
            .map((s) => s.trim())
            .filter((s) => s)
    }

    function resetFilters(): void {
        settingsStore.resetToDefaults()
    }


    // Lifecycle
    function setupEventListeners(): void {
        dropdown.setupListeners()
        if (projectStore.currentPath) {
            loadProjectStructure()
            setTimeout(() => persistence.loadState(), 100)
        }
    }

    function cleanupEventListeners(): void {
        dropdown.cleanupListeners()
    }

    // Watch project changes
    watch(
        () => projectStore.currentPath,
        (path) => {
            if (path) {
                persistence.clearState()
                setTimeout(() => {
                    persistence.loadState()
                    applyFilters()
                }, 150)
            }
        }
    )

    return {
        // Dropdown state (from useFilterDropdown)
        openDropdown: dropdown.openDropdown,
        dropdownStyle: dropdown.dropdownStyle,
        dropdownRefs: dropdown.dropdownRefs,
        toggleDropdown: dropdown.toggleDropdown,

        // UI state
        showSettingsModal,
        isLoading,

        // Filters
        typeFilters,
        languageFilters,
        smartFilters,
        editableFilters,

        // Active filters
        activeTypeFilters,
        activeLanguageFilters,
        activeSmartFilters,
        allActiveFilters,

        // Excluded filters
        allExcludedFilters,

        // Has active flags
        hasActiveTypeFilters,
        hasActiveLanguageFilters,
        hasActiveSmartFilters,

        // Totals
        totalFiles,

        // Filter management
        toggleFilter,
        removeFilter,
        isFilterActive,
        isFilterExcluded,
        clearTypeFilters,
        clearLanguageFilters,
        clearSmartFilters,
        clearAllFilters,

        // Filter info
        getFilterCount,
        getFilterPercentage,

        // Settings
        updateFilterExtensions,
        resetFilters,

        // Lifecycle
        setupEventListeners,
        cleanupEventListeners,
    }
}

// Re-export types
export type { AnyFilter, LanguageFilter, SmartFilter, TypeFilter } from '../model/types'

