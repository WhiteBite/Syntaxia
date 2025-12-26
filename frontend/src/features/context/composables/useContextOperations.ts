/**
 * Context operations composable
 * Handles duplicate, merge, and reorder operations
 */

import { useLogger } from '@/composables/useLogger'
import { contextApi } from '../api/context.api'
import type { ContextSummary } from '../model/context.store'

const logger = useLogger('ContextOperations')

export interface DuplicateResult {
    id: string
    name: string
    content: string
    fileCount: number
    totalSize: number
    lineCount: number
    tokenCount: number
}

export interface MergeResult {
    id: string
    name: string
    content: string
    fileCount: number
}

/**
 * Duplicate a context
 */
export async function duplicateContextContent(
    ctxId: string,
    contexts: ContextSummary[]
): Promise<DuplicateResult | null> {
    try {
        const content = await contextApi.getContextContent(ctxId)
        const original = contexts.find(c => c.id === ctxId)

        if (!original) {
            logger.error('Original context not found:', ctxId)
            return null
        }

        const newId = `ctx-copy-${Date.now()}`
        const newName = `${original.name} (копия)`
        const lines = content.split('\n')

        return {
            id: newId,
            name: newName,
            content,
            fileCount: original.fileCount,
            totalSize: content.length,
            lineCount: lines.length,
            tokenCount: Math.round(content.length / 4)
        }
    } catch (err) {
        logger.error('Failed to duplicate context:', err)
        return null
    }
}

/**
 * Merge multiple contexts into one
 */
export async function mergeContextsContent(
    contextIds: string[],
    contexts: ContextSummary[],
    loadContent: (ctxId: string) => Promise<string>
): Promise<MergeResult | null> {
    if (contextIds.length < 2) {
        logger.warn('Need at least 2 contexts to merge')
        return null
    }

    try {
        const contents: string[] = []
        let totalFiles = 0
        const names: string[] = []

        for (const ctxId of contextIds) {
            const content = await loadContent(ctxId)
            contents.push(content)

            const ctx = contexts.find(c => c.id === ctxId)
            if (ctx) {
                totalFiles += ctx.fileCount
                names.push(ctx.name || ctx.id)
            }
        }

        // Merge contents with separator
        const separator = '\n\n' + '='.repeat(80) + '\n\n'
        const mergedContent = contents.join(separator)

        const newId = `merged-${Date.now()}`
        const namePreview = names.slice(0, 2).join(' + ')
        const suffix = names.length > 2 ? ` +${names.length - 2}` : ''
        const newName = `Merged: ${namePreview}${suffix}`

        return {
            id: newId,
            name: newName,
            content: mergedContent,
            fileCount: totalFiles
        }
    } catch (err) {
        logger.error('Failed to merge contexts:', err)
        return null
    }
}

/**
 * Reorder contexts in list (for drag & drop)
 */
export function reorderContextList(
    contexts: ContextSummary[],
    fromIndex: number,
    toIndex: number
): ContextSummary[] {
    const list = [...contexts]
    const [removed] = list.splice(fromIndex, 1)
    list.splice(toIndex, 0, removed)
    return list
}

/**
 * Create a context summary from raw content
 */
export function createContextSummary(
    id: string,
    name: string,
    content: string,
    fileCount: number
): ContextSummary {
    const lines = content.split('\n')
    return {
        id,
        name,
        fileCount,
        totalSize: content.length,
        lineCount: lines.length,
        tokenCount: Math.round(content.length / 4),
        createdAt: new Date().toISOString(),
        isFavorite: false
    }
}

/**
 * Composable for context operations
 */
export function useContextOperations() {
    return {
        duplicateContextContent,
        mergeContextsContent,
        reorderContextList,
        createContextSummary
    }
}
