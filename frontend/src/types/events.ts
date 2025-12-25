/**
 * Custom Event Types for Syntaxia
 * Provides type safety for window.addEventListener/dispatchEvent
 */

// Event detail types
export interface AddFilesToContextDetail {
    paths: string[]
}

export interface AIStreamChunkDetail {
    content?: string
    done?: boolean
    error?: string
}

// Custom event map
export interface AppCustomEvents {
    'global-build-context': CustomEvent<void>
    'global-open-export': CustomEvent<void>
    'global-copy-context': CustomEvent<void>
    'global-undo-selection': CustomEvent<void>
    'global-redo-selection': CustomEvent<void>
    'add-files-to-context': CustomEvent<AddFilesToContextDetail>
    'ai:stream:chunk': CustomEvent<AIStreamChunkDetail>
}

// Extend WindowEventMap for type-safe event listeners
declare global {
    interface WindowEventMap extends AppCustomEvents { }
}

// Helper to create typed custom events
export function createAppEvent<K extends keyof AppCustomEvents>(
    type: K,
    detail?: AppCustomEvents[K] extends CustomEvent<infer D> ? D : never
): AppCustomEvents[K] {
    return new CustomEvent(type, { detail }) as AppCustomEvents[K]
}

// Type-safe dispatch helper
export function dispatchAppEvent<K extends keyof AppCustomEvents>(
    type: K,
    detail?: AppCustomEvents[K] extends CustomEvent<infer D> ? D : never
): void {
    window.dispatchEvent(createAppEvent(type, detail))
}
