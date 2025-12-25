/**
 * Context content loading and caching composable
 * Handles loading, caching, and chunking of context content
 */

import { useLogger } from '@/composables/useLogger'
import { ref } from 'vue'
import { contextApi } from '../api/context.api'

const logger = useLogger('ContextContent')

export const MAX_CONTEXT_SIZE = 20 * 1024 * 1024 // 20 MB

export interface ContextChunk {
    lines: string[]
    startLine: number
    endLine: number
    hasMore: boolean
}

export interface ContentCache {
    contextId: string | null
    content: string | null
}

/**
 * Composable for context content operations
 */
export function useContextContent() {
    const cache = ref<ContentCache>({
        contextId: null,
        content: null
    })

    /**
     * Clear the content cache
     */
    function clearCache() {
        cache.value = { contextId: null, content: null }
    }

    /**
     * Check if content is cached for given context
     */
    function isCached(contextId: string): boolean {
        return cache.value.contextId === contextId && cache.value.content !== null
    }

    /**
     * Get cached content or null
     */
    function getCachedContent(contextId: string): string | null {
        if (isCached(contextId)) {
            return cache.value.content
        }
        return null
    }

    /**
     * Load full context content with caching
     */
    async function loadFullContent(contextId: string): Promise<string> {
        // Return cached if available
        const cached = getCachedContent(contextId)
        if (cached !== null) {
            logger.debug('Using cached content')
            return cached
        }

        logger.debug('Loading full content from API...')
        const content = await contextApi.getContextContent(contextId)

        // Check size limit
        if (content.length > MAX_CONTEXT_SIZE) {
            const sizeMB = Math.round(content.length / (1024 * 1024))
            const limitMB = Math.round(MAX_CONTEXT_SIZE / (1024 * 1024))
            throw new Error(`Context content (${sizeMB}MB) exceeds maximum allowed size (${limitMB}MB)`)
        }

        // Update cache
        cache.value = { contextId, content }
        logger.debug('Content cached, size:', Math.round(content.length / 1024), 'KB')

        return content
    }

    /**
     * Load a chunk of context content
     * If lineCount is 0 or negative, loads all lines from startLine
     */
    async function loadChunk(
        contextId: string,
        startLine: number = 0,
        lineCount: number = 0
    ): Promise<ContextChunk> {
        const content = await loadFullContent(contextId)
        const lines = content.split('\n')
        // lineCount <= 0 means load all remaining lines
        const endLine = lineCount > 0
            ? Math.min(startLine + lineCount, lines.length)
            : lines.length

        logger.debug('Loaded chunk:', startLine, '-', endLine, 'of', lines.length, 'lines')

        return {
            lines: lines.slice(startLine, endLine),
            startLine,
            endLine,
            hasMore: endLine < lines.length
        }
    }

    /**
     * Estimate memory usage of cached content
     */
    function getMemoryUsage(): number {
        if (!cache.value.content) return 0
        // UTF-16 encoding: 2 bytes per character
        return cache.value.content.length * 2
    }

    return {
        cache,
        clearCache,
        isCached,
        getCachedContent,
        loadFullContent,
        loadChunk,
        getMemoryUsage
    }
}
