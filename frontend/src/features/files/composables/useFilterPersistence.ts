/**
 * Composable for filter state persistence per project
 */
import { useProjectStore } from '@/stores/project.store'
import { watch, type Ref } from 'vue'
import { STORAGE_KEY } from '../constants/filterConfig'
import type { FilterState } from '../model/types'

export function useFilterPersistence(filterState: Ref<FilterState>) {
    const projectStore = useProjectStore()

    function getProjectKey(): string {
        return projectStore.currentPath?.replace(/[\\/:]/g, '_') || 'default'
    }

    function saveState(): void {
        try {
            const allStates = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
            allStates[getProjectKey()] = {
                active: Array.from(filterState.value.active),
                excluded: Array.from(filterState.value.excluded),
            }
            localStorage.setItem(STORAGE_KEY, JSON.stringify(allStates))
        } catch {
            // Silently fail - persistence is optional
        }
    }

    function loadState(): void {
        try {
            const allStates = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}')
            const saved = allStates[getProjectKey()]
            if (saved) {
                if (Array.isArray(saved.active)) {
                    filterState.value.active = new Set(saved.active)
                } else if (Array.isArray(saved)) {
                    // Legacy format support
                    filterState.value.active = new Set(saved)
                }
                if (Array.isArray(saved.excluded)) {
                    filterState.value.excluded = new Set(saved.excluded)
                }
            }
        } catch {
            // Silently fail
        }
    }

    function clearState(): void {
        filterState.value.active.clear()
        filterState.value.excluded.clear()
    }

    // Auto-save on changes
    watch(
        () => [filterState.value.active.size, filterState.value.excluded.size],
        () => saveState(),
        { deep: true }
    )

    return {
        saveState,
        loadState,
        clearState,
        getProjectKey,
    }
}
