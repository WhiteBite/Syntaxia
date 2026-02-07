/**
 * useFileFuzzySearch - Fuzzy file search functionality with backend optimization
 * Used by file.store for search results computation
 * 
 * Note: This is different from features/files/composables/useFileSearch.ts
 * which handles search input with debounce for the UI
 * 
 * OPTIMIZATION: Uses backend search API for better performance on large projects
 */

import { searchFiles } from '@/services/search.api'
import type { FileNode } from '@/types/domain'
import Fuse, { type FuseResultMatch } from 'fuse.js'
import { computed, ref, watch, type ComputedRef } from 'vue'

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
    const backendResults = ref<SearchResultNode[]>([])
    const isSearching = ref(false)
    const useBackend = ref(true) // Feature flag for backend search

    // Watch for search query changes and trigger backend search
    watch(searchQuery, async (newQuery) => {
        if (!newQuery || !rootPath?.value || !useBackend.value) {
            backendResults.value = []
            return
        }

        isSearching.value = true
        try {
            const results = await searchFiles(rootPath.value as string, newQuery, {
                maxResults,
                fuzzyMatch: false,
                caseSensitive: false,
                includePath: true,
            })

            // Convert to SearchResultNode format
            backendResults.value = results.map(result => {
                const fuseMatches: FuseResultMatch[] = result.matches?.map(match => ({
                    indices: [[match.start, match.end - 1]] as [number, number][],
                    value: match.field === 'name' ? result.name : result.path,
                    key: match.field,
                    refIndex: 0
                })) || []

                return {
                    name: result.name,
                    path: result.path,
                    relPath: result.path.replace((rootPath?.value || '') + '/', ''),
                    isDir: false,
                    size: result.size,
                    contentType: result.contentType as 'text' | 'binary' | 'unknown',
                    children: [],
                    isGitignored: false,
                    isCustomIgnored: false,
                    isIgnored: false,
                    fileCount: 1,
                    totalSize: result.size,
                    depth: result.depth || 0,
                    directFileCount: 0,
                    relativePath: result.path.replace((rootPath?.value || '') + '/', ''),
                    score: result.score,
                    matches: fuseMatches
                } as SearchResultNode
            })
        } catch (error) {
            // Fallback to frontend search on error
            useBackend.value = false
            backendResults.value = []
        } finally {
            isSearching.value = false
        }
    })

    // Computed - returns backend results if available, otherwise frontend fallback
    const searchResults = computed((): SearchResultNode[] => {
        if (!searchQuery.value) return []

        // Use backend results if available
        if (useBackend.value && backendResults.value.length > 0) {
            return backendResults.value
        }

        // Fallback to frontend search
        const allFiles = flattenedNodes.value.filter(node => !node.isDir)
        const rawQuery = searchQuery.value.trim().toLowerCase()

        // Operator handling
        const isExclude = rawQuery.startsWith('!')
        const isExtensionOnly = rawQuery.startsWith('.') && !rawQuery.includes(' ')
        const searchTerm = isExclude ? rawQuery.substring(1).trim() : rawQuery

        if (isExtensionOnly) {
            return allFiles
                .filter(f => f.name.toLowerCase().endsWith(searchTerm))
                .map(f => ({
                    ...f,
                    depth: 0,
                    relativePath: rootPath?.value ? f.path.replace(rootPath.value + '/', '') : f.path
                }))
                .slice(0, maxResults)
        }

        // For large trees or complex queries, use simple string matching
        if (allFiles.length > 2000 || isExclude || searchTerm.includes(' ')) {
            const terms = searchTerm.split(' ').filter(Boolean)
            
            return allFiles
                .filter(file => {
                    const name = file.name.toLowerCase()
                    const path = file.path.toLowerCase()
                    
                    if (isExclude) {
                        return !name.includes(searchTerm) && !path.includes(searchTerm)
                    }
                    
                    // All terms must match (AND logic)
                    return terms.every(term => name.includes(term) || path.includes(term))
                })
                .map(file => {
                    const matches: FuseResultMatch[] = []
                    // Highlight first term for simplicity in non-fuse mode
                    const firstTerm = terms[0]
                    if (firstTerm) {
                        const nameIndex = file.name.toLowerCase().indexOf(firstTerm)
                        if (nameIndex !== -1) {
                            matches.push({
                                indices: [[nameIndex, nameIndex + firstTerm.length - 1]],
                                value: file.name,
                                key: 'name',
                                refIndex: 0
                            })
                        }
                    }

                    return {
                        ...file,
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
        backendResults.value = []
    }

    return {
        // State
        searchQuery,
        isSearching,
        // Computed
        searchResults,
        // Actions
        setSearchQuery,
        clearSearch,
    }
}
