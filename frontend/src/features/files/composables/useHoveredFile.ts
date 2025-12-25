/**
 * Composable for managing hovered file state using reactive singleton pattern.
 * 
 * NOTE: We use a singleton reactive state instead of provide/inject because
 * vue-virtual-scroller's RecycleScroller reuses DOM elements and can call
 * provide() outside of setup() context, causing Vue warnings.
 * 
 * This approach avoids:
 * - "[Vue warn]: provide() can only be used inside setup()"
 * - "[Vue warn]: injection "Symbol(HoveredFile)" not found"
 */

import { reactive } from 'vue'

export interface HoveredFileState {
    path: string | null
    isDir: boolean | null
}

interface HoveredFileActions {
    readonly state: HoveredFileState
    setHovered: (path: string | null, isDir: boolean | null) => void
    clearHovered: (currentPath: string) => void
}

// Singleton reactive state - shared across all components
const hoveredState = reactive<HoveredFileState>({
    path: null,
    isDir: null,
})

/**
 * Provides hovered file state to child components.
 * Now a no-op for backwards compatibility - state is managed via singleton.
 * @deprecated Use useHoveredFile() directly instead
 */
export function provideHoveredFile(): HoveredFileActions {
    return useHoveredFile()
}

/**
 * Returns hovered file state and actions.
 * Uses singleton reactive state to avoid provide/inject issues with virtual scrolling.
 */
export function useHoveredFile(): HoveredFileActions {
    const setHovered = (newPath: string | null, newIsDir: boolean | null) => {
        hoveredState.path = newPath
        hoveredState.isDir = newIsDir
    }

    const clearHovered = (currentPath: string) => {
        if (hoveredState.path === currentPath) {
            hoveredState.path = null
            hoveredState.isDir = null
        }
    }

    return {
        state: hoveredState,
        setHovered,
        clearHovered,
    }
}
