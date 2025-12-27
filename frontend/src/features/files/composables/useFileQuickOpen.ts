/**
 * useFileQuickOpen - Composable for quick file search (Ctrl+P)
 * Handles fuzzy search, keyboard navigation, and file selection
 */

import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import type { FileNode } from '@/types/domain'
import Fuse, { type FuseResultMatch } from 'fuse.js'
import { computed, ref, watch, type Ref } from 'vue'
import { useFileStore } from '../model/file.store'

export interface QuickOpenResult {
    item: FileNode
    matches?: ReadonlyArray<FuseResultMatch>
}

export interface UseFileQuickOpenOptions {
    maxResults?: number
    debounceMs?: number
}

const MAX_RESULTS = 20
const DEBOUNCE_MS = 150

export function useFileQuickOpen(options: UseFileQuickOpenOptions = {}) {
    const { maxResults = MAX_RESULTS, debounceMs = DEBOUNCE_MS } = options

    const fileStore = useFileStore()
    const uiStore = useUIStore()
    const { t } = useI18n()

    // State
    const query = ref('')
    const selectedIndex = ref(0)
    const isSearching = ref(false)
    const inputRef = ref<HTMLInputElement | null>(null)

    // Debounce timer
    let debounceTimer: ReturnType<typeof setTimeout> | null = null

    // Fuse.js instance (recreated when files change)
    const fuse = computed(() => {
        const files = fileStore.flattenedNodes.filter(node => !node.isDir)
        return new Fuse(files, {
            keys: [
                { name: 'name', weight: 0.7 },
                { name: 'path', weight: 0.3 }
            ],
            threshold: 0.4,
            includeMatches: true,
            ignoreLocation: true,
            minMatchCharLength: 1
        })
    })

    // Search results with match info
    const results: Ref<QuickOpenResult[]> = ref([])

    // Perform search with debounce
    function performSearch() {
        if (!query.value.trim()) {
            results.value = []
            selectedIndex.value = 0
            isSearching.value = false
            return
        }

        isSearching.value = true

        if (debounceTimer) {
            clearTimeout(debounceTimer)
        }

        debounceTimer = setTimeout(() => {
            const searchResults = fuse.value.search(query.value)
            results.value = searchResults.slice(0, maxResults).map(result => ({
                item: result.item,
                matches: result.matches
            }))
            selectedIndex.value = 0
            isSearching.value = false
        }, debounceMs)
    }

    // Watch query changes
    watch(query, performSearch)

    // Keyboard navigation
    function handleKeyDown(event: KeyboardEvent) {
        switch (event.key) {
            case 'ArrowDown':
                event.preventDefault()
                if (results.value.length > 0) {
                    selectedIndex.value = (selectedIndex.value + 1) % results.value.length
                }
                break
            case 'ArrowUp':
                event.preventDefault()
                if (results.value.length > 0) {
                    selectedIndex.value = selectedIndex.value === 0
                        ? results.value.length - 1
                        : selectedIndex.value - 1
                }
                break
            case 'Enter':
                event.preventDefault()
                selectCurrentItem()
                break
            case 'Escape':
                event.preventDefault()
                close()
                break
        }
    }

    // Select file and add to context
    function selectItem(index: number) {
        const result = results.value[index]
        if (!result) return

        addFileToContext(result.item.path)
        close()
    }

    function selectCurrentItem() {
        if (results.value.length > 0) {
            selectItem(selectedIndex.value)
        }
    }

    // Add file to context (toggle selection in file store)
    function addFileToContext(filePath: string) {
        if (!fileStore.selectedPaths.has(filePath)) {
            fileStore.toggleSelect(filePath)
            uiStore.addToast(t('fileSearch.fileAdded'), 'success', 2000)
        } else {
            uiStore.addToast(t('fileSearch.fileAlreadySelected'), 'info', 2000)
        }
    }

    // Close modal
    function close() {
        uiStore.closeFileSearchModal()
        reset()
    }

    // Reset state
    function reset() {
        query.value = ''
        results.value = []
        selectedIndex.value = 0
        isSearching.value = false
        if (debounceTimer) {
            clearTimeout(debounceTimer)
            debounceTimer = null
        }
    }

    // Focus input
    function focusInput() {
        inputRef.value?.focus()
    }

    return {
        // State
        query,
        results,
        selectedIndex,
        isSearching,
        inputRef,
        // Actions
        handleKeyDown,
        selectItem,
        selectCurrentItem,
        close,
        reset,
        focusInput
    }
}
