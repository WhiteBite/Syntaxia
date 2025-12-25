import { computed, ref, type Ref } from 'vue'
import { FILE_TYPE_FILTERS } from '../model/types'

/**
 * Composable for file filtering with search and type filters
 * Optimized with memoization for large file lists
 */
export function useGitFilters(files: Ref<string[]>) {
    const searchQuery = ref('')
    const activeFilters = ref<Set<string>>(new Set())

    // Memoized extension extraction
    const fileExtensionCache = new Map<string, string>()

    function getExtension(filePath: string): string {
        if (fileExtensionCache.has(filePath)) {
            return fileExtensionCache.get(filePath)!
        }
        const ext = '.' + (filePath.split('.').pop()?.toLowerCase() || '')
        fileExtensionCache.set(filePath, ext)
        return ext
    }

    // Clear cache when files change significantly
    function clearExtensionCache() {
        if (fileExtensionCache.size > 10000) {
            fileExtensionCache.clear()
        }
    }

    // Filtered files with optimized filtering
    const filteredFiles = computed(() => {
        let result = files.value || []

        // Apply search filter
        if (searchQuery.value) {
            const query = searchQuery.value.toLowerCase()
            result = result.filter(f => f.toLowerCase().includes(query))
        }

        // Apply type filters
        if (activeFilters.value.size > 0) {
            const activeExtensions = new Set<string>()
            activeFilters.value.forEach(filterId => {
                const filter = FILE_TYPE_FILTERS.find(f => f.id === filterId)
                if (filter) {
                    filter.extensions.forEach(ext => activeExtensions.add(ext))
                }
            })

            result = result.filter(f => activeExtensions.has(getExtension(f)))
        }

        return result
    })

    // Filter counts for UI badges
    const filterCounts = computed(() => {
        const counts: Record<string, number> = {}

        for (const filter of FILE_TYPE_FILTERS) {
            counts[filter.id] = files.value.filter(f =>
                filter.extensions.includes(getExtension(f))
            ).length
        }

        return counts
    })

    function toggleFilter(filterId: string) {
        if (activeFilters.value.has(filterId)) {
            activeFilters.value.delete(filterId)
        } else {
            activeFilters.value.add(filterId)
        }
        activeFilters.value = new Set(activeFilters.value)
    }

    function clearFilters() {
        activeFilters.value = new Set()
        searchQuery.value = ''
    }

    function setSearch(query: string) {
        searchQuery.value = query
    }

    // Debounced search for performance
    let searchTimeout: ReturnType<typeof setTimeout> | null = null

    function setSearchDebounced(query: string, delay: number = 150) {
        if (searchTimeout) {
            clearTimeout(searchTimeout)
        }
        searchTimeout = setTimeout(() => {
            searchQuery.value = query
        }, delay)
    }

    return {
        searchQuery,
        activeFilters,
        filteredFiles,
        filterCounts,
        filters: FILE_TYPE_FILTERS,
        toggleFilter,
        clearFilters,
        setSearch,
        setSearchDebounced,
        clearExtensionCache
    }
}

/**
 * Composable for file selection management
 */
export function useFileSelection(files: Ref<string[]>) {
    const selectedFiles = ref<Set<string>>(new Set())

    const selectedCount = computed(() => selectedFiles.value.size)
    const hasSelection = computed(() => selectedFiles.value.size > 0)
    const allSelected = computed(() =>
        files.value.length > 0 && files.value.every(f => selectedFiles.value.has(f))
    )

    function isSelected(path: string): boolean {
        return selectedFiles.value.has(path)
    }

    function toggle(path: string) {
        if (selectedFiles.value.has(path)) {
            selectedFiles.value.delete(path)
        } else {
            selectedFiles.value.add(path)
        }
        selectedFiles.value = new Set(selectedFiles.value)
    }

    function selectAll() {
        selectedFiles.value = new Set(files.value)
    }

    function clearSelection() {
        selectedFiles.value = new Set()
    }

    function selectFolder(folderPath: string) {
        const folderFiles = files.value.filter(f => f.startsWith(folderPath + '/'))
        const allFolderSelected = folderFiles.every(f => selectedFiles.value.has(f))

        if (allFolderSelected) {
            folderFiles.forEach(f => selectedFiles.value.delete(f))
        } else {
            folderFiles.forEach(f => selectedFiles.value.add(f))
        }
        selectedFiles.value = new Set(selectedFiles.value)
    }

    function getFilesInFolder(folderPath: string): string[] {
        return files.value.filter(f => f.startsWith(folderPath + '/'))
    }

    function hasSomeInFolder(folderPath: string): boolean {
        return getFilesInFolder(folderPath).some(f => selectedFiles.value.has(f))
    }

    function hasAllInFolder(folderPath: string): boolean {
        const folderFiles = getFilesInFolder(folderPath)
        return folderFiles.length > 0 && folderFiles.every(f => selectedFiles.value.has(f))
    }

    function setSelection(files: Set<string>) {
        selectedFiles.value = new Set(files)
    }

    function getSelectedArray(): string[] {
        return Array.from(selectedFiles.value)
    }

    return {
        selectedFiles,
        selectedCount,
        hasSelection,
        allSelected,
        isSelected,
        toggle,
        selectAll,
        clearSelection,
        selectFolder,
        getFilesInFolder,
        hasSomeInFolder,
        hasAllInFolder,
        setSelection,
        getSelectedArray
    }
}
