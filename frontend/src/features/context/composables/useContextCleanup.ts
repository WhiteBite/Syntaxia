/**
 * Context auto-cleanup composable
 * Handles automatic deletion of old contexts based on settings
 */

import { useLogger } from '@/composables/useLogger'
import { useSettingsStore } from '@/stores/settings.store'
import type { ContextSummary } from '../model/context.store'

const logger = useLogger('ContextCleanup')

export interface CleanupOptions {
    maxContexts: number
    autoCleanupDays: number
    autoCleanupOnLimit: boolean
}

export interface CleanupResult {
    deletedByLimit: string[]
    deletedByAge: string[]
}

/**
 * Get non-favorite contexts sorted by creation date (oldest first)
 */
function getNonFavoritesSorted(contexts: ContextSummary[]): ContextSummary[] {
    return contexts
        .filter(c => !c.isFavorite)
        .sort((a, b) => {
            const dateA = new Date(a.createdAt || 0).getTime()
            const dateB = new Date(b.createdAt || 0).getTime()
            return dateA - dateB
        })
}

/**
 * Determine which contexts should be deleted based on count limit
 */
export function getContextsToDeleteByLimit(
    contexts: ContextSummary[],
    maxContexts: number
): ContextSummary[] {
    const totalCount = contexts.length
    if (totalCount <= maxContexts) return []

    const nonFavorites = getNonFavoritesSorted(contexts)
    const toDelete = totalCount - maxContexts

    return nonFavorites.slice(0, toDelete)
}

/**
 * Determine which contexts should be deleted based on age
 */
export function getContextsToDeleteByAge(
    contexts: ContextSummary[],
    maxAgeDays: number
): ContextSummary[] {
    if (maxAgeDays <= 0) return []

    const cutoffDate = Date.now() - maxAgeDays * 24 * 60 * 60 * 1000
    const nonFavorites = getNonFavoritesSorted(contexts)

    return nonFavorites.filter(ctx => {
        const ctxDate = new Date(ctx.createdAt || 0).getTime()
        return ctxDate < cutoffDate
    })
}

/**
 * Perform auto-cleanup of contexts
 */
export async function performAutoCleanup(
    contexts: ContextSummary[],
    deleteContext: (id: string) => Promise<void>
): Promise<CleanupResult> {
    const settingsStore = useSettingsStore()
    const storage = settingsStore.settings.contextStorage

    const result: CleanupResult = {
        deletedByLimit: [],
        deletedByAge: []
    }

    if (!storage.autoCleanupOnLimit) {
        return result
    }

    // Delete by count limit
    const byLimit = getContextsToDeleteByLimit(contexts, storage.maxContexts)
    for (const ctx of byLimit) {
        try {
            await deleteContext(ctx.id)
            result.deletedByLimit.push(ctx.id)
            logger.debug('Auto-deleted old context:', ctx.id)
        } catch (err) {
            logger.error('Failed to auto-delete context:', err)
        }
    }

    // Delete by age
    const byAge = getContextsToDeleteByAge(contexts, storage.autoCleanupDays)
    for (const ctx of byAge) {
        // Skip if already deleted by limit
        if (result.deletedByLimit.includes(ctx.id)) continue

        try {
            await deleteContext(ctx.id)
            result.deletedByAge.push(ctx.id)
            logger.debug('Auto-deleted expired context:', ctx.id)
        } catch (err) {
            logger.error('Failed to auto-delete expired context:', err)
        }
    }

    return result
}

/**
 * Composable for context cleanup operations
 */
export function useContextCleanup() {
    return {
        getContextsToDeleteByLimit,
        getContextsToDeleteByAge,
        performAutoCleanup
    }
}
