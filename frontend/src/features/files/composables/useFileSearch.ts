/**
 * Composable for file search functionality
 * Extracted from useFileExplorer for better separation of concerns
 */
import { FILE_TREE } from '@/config/constants'
import { debounce } from '@/utils/performance'
import { ref } from 'vue'
import { useFileStore } from '../model/file.store'

export function useFileSearch() {
    const fileStore = useFileStore()

    const query = ref('')
    const isSearching = ref(false)

    // Debounced search for better performance
    const debouncedSearch = debounce(() => {
        fileStore.setSearchQuery(query.value)
        isSearching.value = false
    }, FILE_TREE.DEBOUNCE_MS)

    /**
     * Handle search input change
     */
    function handleSearch() {
        isSearching.value = true
        debouncedSearch()
    }

    /**
     * Clear search query
     */
    function clear() {
        query.value = ''
        fileStore.setSearchQuery('')
        isSearching.value = false
    }

    /**
     * Set search query programmatically
     */
    function setQuery(newQuery: string) {
        query.value = newQuery
        handleSearch()
    }

    /**
     * Check if search is active
     */
    function isActive(): boolean {
        return query.value.length > 0
    }

    return {
        query,
        isSearching,
        handleSearch,
        clear,
        setQuery,
        isActive,
    }
}
