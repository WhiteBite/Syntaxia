/**
 * Context metadata persistence composable
 * Handles saving/loading context names and favorites to localStorage
 */

import { STORAGE_KEYS } from '@/config/constants'
import type { ContextSummary } from '../model/context.store'

interface ContextMetadataItem {
    id: string
    name: string
    isFavorite?: boolean
}

/**
 * Save context metadata to localStorage
 */
export function saveContextMetadata(contexts: ContextSummary[]): void {
    try {
        const metadata: ContextMetadataItem[] = contexts.map(c => ({
            id: c.id,
            name: c.name,
            isFavorite: c.isFavorite
        }))
        localStorage.setItem(STORAGE_KEYS.CONTEXT_METADATA, JSON.stringify(metadata))
    } catch (err) {
        console.warn('[ContextMetadata] Failed to save:', err)
    }
}

/**
 * Load and apply context metadata from localStorage
 */
export function loadContextMetadata(contexts: ContextSummary[]): void {
    try {
        const saved = localStorage.getItem(STORAGE_KEYS.CONTEXT_METADATA)
        if (!saved) return

        const metadata = JSON.parse(saved) as ContextMetadataItem[]

        for (const ctx of contexts) {
            const meta = metadata.find(m => m.id === ctx.id)
            if (meta) {
                ctx.name = meta.name
                ctx.isFavorite = meta.isFavorite
            }
        }
    } catch (err) {
        console.warn('[ContextMetadata] Failed to load:', err)
    }
}

/**
 * Composable for context metadata operations
 */
export function useContextMetadata() {
    return {
        save: saveContextMetadata,
        load: loadContextMetadata
    }
}
