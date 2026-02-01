/**
 * useFileFuzzySearch - Fuzzy file search functionality using Fuse.js
 * Used by file.store for search results computation
 * 
 * Note: This is different from features/files/composables/useFileSearch.ts
 * which handles search input with debounce for the UI
 */

import type { FileNode } from '@/types/domain'
import Fuse, { type FuseResultMatch } from 'fuse.js'
import { computed, ref, type ComputedRef } from 'vue'

export interface UseFileFuzzySearchOptions {
    flattenedNodes: ComputedRef<FileNode[]>
    rootPath?: ComputedRef<string>
    maxResults?: number
    fuseThreshold?: number
}

export interface SearchResultNode extends FileNode {
    depth: number
    relativePath?: string
    score?: number
    matches?: ReadonlyArray<FuseResultMatch>
}

export function useFileFuzzySearch(options: UseFileFuzzySearchOptions) {
    const { flattenedNodes, rootPath, maxResults = 100, fuseThreshold = 0.3 } = options

    // State
    const searchQuery = ref('')

    // Computed
    const searchResults = computed((): SearchResultNode[] => {
        if (!searchQuery.value) return []

        // FLAT SEARCH: Only files, no hierarchy
        const allFiles = flattenedNodes.value.filter(node => !node.isDir)

        // For large trees, use simple string matching
        if (allFiles.length > 2000) {
            const query = searchQuery.value.toLowerCase()
            return allFiles
                .filter(
                    (file) =>
                        file.name.toLowerCase().includes(query) ||
                        file.path.toLowerCase().includes(query)
                )
                .map(file => {
                    // Create simple match indices for highlighting
                    const nameIndex = file.name.toLowerCase().indexOf(query)
                    const matches: FuseResultMatch[] = []

                    if (nameIndex !== -1) {
                        matches.push({
                            indices: [[nameIndex, nameIndex + query.length - 1]],
                            value: file.name,
                            key: 'name',
                            refIndex: 0
                        })
                    }

                    return {
                        ...file,
                        // Flat list: no depth, show relative path
                        depth: 0,
                        relativePath: rootPath?.value ? file.path.replace(rootPath.value + '/', '') : file.path,
                        matches: matches.length > 0 ? matches : undefined
                    }
                })
                .slice(0, maxResults)
        }

        // For smaller trees, use Fuse.js for fuzzy search
        const fuse = new Fuse(allFiles, {
            keys: ['name', 'path'],
            threshold: fuseThreshold,
            includeScore: true,
            includeMatches: true
        })

        return fuse
            .search(searchQuery.value)
            .map((result) => ({
                ...result.item,
                depth: 0,
                relativePath: rootPath?.value ? result.item.path.replace(rootPath.value + '/', '') : result.item.path,
                score: result.score,
                matches: result.matches
            }))
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
