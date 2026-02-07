/**
 * Search API - Backend search and filter operations
 * Replaces frontend Fuse.js with backend search index
 */

import { FilterFilesByExtension, FilterFilesByWeight, SearchFiles } from '#wailsjs/go/main/App'
import type { domain } from '#wailsjs/go/models'
import type { FileNode } from '@/types/domain'

export interface SearchOptions {
    maxResults?: number
    fuzzyMatch?: boolean
    caseSensitive?: boolean
    includePath?: boolean
    fileTypes?: string[]
}

export interface FileSearchResult {
    path: string
    name: string
    score: number
    matches?: MatchRange[]
    size: number
    contentType: string
    depth: number
}

export interface MatchRange {
    start: number
    end: number
    field: string // "name" or "path"
}

/**
 * Search files using backend search index
 */
export async function searchFiles(
    projectRoot: string,
    query: string,
    options: SearchOptions = {}
): Promise<FileSearchResult[]> {
    const searchOptions: domain.SearchOptions = {
        maxResults: options.maxResults || 100,
        fuzzyMatch: options.fuzzyMatch ?? false,
        caseSensitive: options.caseSensitive ?? false,
        includePath: options.includePath ?? true,
        fileTypes: options.fileTypes || []
    }

    try {
        const results = await SearchFiles(projectRoot, query, searchOptions)
        return results as FileSearchResult[]
    } catch (error) {
        console.error('Backend search failed:', error)
        throw error
    }
}

/**
 * Filter files by extension using backend pre-computed metadata
 */
export async function filterByExtension(
    projectRoot: string,
    includeExts: string[],
    excludeExts: string[] = []
): Promise<FileNode[]> {
    try {
        const filtered = await FilterFilesByExtension(projectRoot, includeExts, excludeExts)
        return filtered as FileNode[]
    } catch (error) {
        console.error('Backend filter by extension failed:', error)
        throw error
    }
}

/**
 * Filter files by token weight using backend pre-computed TotalSize
 */
export async function filterByWeight(
    projectRoot: string,
    minTokens: number
): Promise<FileNode[]> {
    try {
        const filtered = await FilterFilesByWeight(projectRoot, minTokens)
        return filtered as FileNode[]
    } catch (error) {
        console.error('Backend filter by weight failed:', error)
        throw error
    }
}
