/**
 * Undo/Redo Composable
 * Generic history management for any state type
 */

import { computed, ref, shallowRef } from 'vue'

export function useUndoRedo<T>(initialState: T, maxHistory = 20) {
    const history = shallowRef<T[]>([initialState])
    const currentIndex = ref(0)

    const current = computed(() => history.value[currentIndex.value])
    const canUndo = computed(() => currentIndex.value > 0)
    const canRedo = computed(() => currentIndex.value < history.value.length - 1)

    function push(state: T) {
        // Remove "future" states if we're in the middle of history
        const newHistory = history.value.slice(0, currentIndex.value + 1)
        newHistory.push(state)

        // Limit history size
        if (newHistory.length > maxHistory) {
            newHistory.shift()
        } else {
            currentIndex.value++
        }

        history.value = newHistory
    }

    function undo(): T | undefined {
        if (canUndo.value) {
            currentIndex.value--
            return current.value
        }
        return undefined
    }

    function redo(): T | undefined {
        if (canRedo.value) {
            currentIndex.value++
            return current.value
        }
        return undefined
    }

    function clear() {
        history.value = [history.value[currentIndex.value]]
        currentIndex.value = 0
    }

    function reset(newInitialState: T) {
        history.value = [newInitialState]
        currentIndex.value = 0
    }

    return {
        current,
        canUndo,
        canRedo,
        historyLength: computed(() => history.value.length),
        push,
        undo,
        redo,
        clear,
        reset,
    }
}
