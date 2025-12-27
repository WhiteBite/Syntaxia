/**
 * Composable for managing favorite files
 * Persists favorites to localStorage with a maximum limit
 */

import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { ref, watch } from 'vue'

const STORAGE_KEY = 'syntaxia-favorites'
const MAX_FAVORITES = 20

// Singleton reactive state
const favorites = ref<string[]>([])
let initialized = false

function loadFromStorage(): string[] {
    try {
        const stored = localStorage.getItem(STORAGE_KEY)
        return stored ? JSON.parse(stored) : []
    } catch {
        return []
    }
}

function saveToStorage(paths: string[]): void {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(paths))
}

export function useFavorites() {
    const { t } = useI18n()
    const uiStore = useUIStore()

    // Initialize once
    if (!initialized) {
        favorites.value = loadFromStorage()
        initialized = true

        // Auto-save on changes
        watch(favorites, (newValue) => {
            saveToStorage(newValue)
        }, { deep: true })
    }

    function isFavorite(path: string): boolean {
        return favorites.value.includes(path)
    }

    function addFavorite(path: string): boolean {
        if (isFavorite(path)) return false
        if (favorites.value.length >= MAX_FAVORITES) {
            uiStore.addToast(t('files.favorites.maxReached'), 'warning')
            return false
        }
        favorites.value.push(path)
        uiStore.addToast(t('files.favorites.added'), 'success')
        return true
    }

    function removeFavorite(path: string): boolean {
        const index = favorites.value.indexOf(path)
        if (index === -1) return false
        favorites.value.splice(index, 1)
        uiStore.addToast(t('files.favorites.removed'), 'info')
        return true
    }

    function toggleFavorite(path: string): void {
        if (isFavorite(path)) {
            removeFavorite(path)
        } else {
            addFavorite(path)
        }
    }

    return {
        favorites,
        isFavorite,
        addFavorite,
        removeFavorite,
        toggleFavorite,
        MAX_FAVORITES,
    }
}
