/**
 * useFileFuzzySearch - Fuzzy file search functionality using Fuse.js
 * Used by file.store for search results computation
 * 
 * Note: This is different from features/files/composables/useFileSearch.ts
 * which handles search input with debounce for the UI
 */

import type { FileNode } from '@/types/domain'
import Fuse from 'fuse.js'
import { computed, ref, type ComputedRef } from 'vue'

export interface UseFileFuzzySearchOptions {
    flattenedNodes: ComputedRef<FileNode[]>
    maxResults?: number
    fuseThreshold?: number
}

export function useFileFuzzySearch(options: UseFileFuzzySearchOptions) {
    const { flattenedNodes, maxResults = 100, fuseThreshold = 0.3 } = options

    // State
    const searchQuery = ref('')

    // Computed
    const searchResults = computed(() => {
        if (!searchQuery.value) return []

        const allFiles = flattenedNodes.value

        // For large trees, use simple string matching
        if (allFiles.length > 2000) {
            const query = searchQuery.value.toLowerCase()
            return allFiles
                .filter(
                    (file) =>
                        file.name.toLowerCase().includes(query) ||
                        file.path.toLowerCase().includes(query)
                )
                .slice(0, maxResults)
        }

        // For smaller trees, use Fuse.js for fuzzy search
        const fuse = new Fuse(allFiles, {
            keys: ['name', 'path'],
            threshold: fuseThreshold,
        })

        return fuse
            .search(searchQuery.value)
            .map((result) => result.item)
            .slice(0, maxResults)
    })

    // Actions
    function setSearchQuery(query: string) {
        searchQuery.value = query
    }

    function clearSearch() {
        searchQuery.value = ''
    }

    return {
        // State
        searchQuery,
        // Computed
        searchResults,
        // Actions
        setSearchQuery,
        clearSearch,
    }
}
