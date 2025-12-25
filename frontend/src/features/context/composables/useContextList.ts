import { useI18n } from '@/composables/useI18n'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, nextTick, reactive, ref, watch } from 'vue'
import { formatContextSize, formatTimestamp as formatTs } from '../lib/context-utils'
import { useContextStore, type ContextSummary } from '../model/context.store'

type SortBy = 'date' | 'name' | 'size'

export function useContextList() {
    const { t } = useI18n()
    const contextStore = useContextStore()
    const settingsStore = useSettingsStore()
    const uiStore = useUIStore()

    // State
    const showSettings = ref(false)
    const searchQuery = ref('')
    const sortBy = ref<SortBy>('date')
    const showFavoritesOnly = ref(false)
    const selectedContexts = ref(new Set<string>())
    const editingId = ref<string | null>(null)
    const editingName = ref('')
    const renameInput = ref<HTMLInputElement | null>(null)

    // Drag & drop state
    const dragIndex = ref<number | null>(null)
    const dragOverIndex = ref<number | null>(null)

    // Delete modal
    const deleteModal = reactive({
        show: false,
        contextId: null as string | null,
        contextIds: [] as string[],
        message: ''
    })

    // Storage settings (reactive binding)
    const storageSettings = computed({
        get: () => settingsStore.settings.contextStorage,
        set: (val) => settingsStore.updateContextStorageSettings(val)
    })

    // Watch for settings changes - auto-saved by store
    watch(() => settingsStore.settings.contextStorage, () => { }, { deep: true })

    // Filtered and sorted contexts
    const filteredContexts = computed(() => {
        let list = [...contextStore.contextList]

        // Filter by search
        if (searchQuery.value) {
            const q = searchQuery.value.toLowerCase()
            list = list.filter(c => (c.name || c.id).toLowerCase().includes(q))
        }

        // Filter favorites
        if (showFavoritesOnly.value) {
            list = list.filter(c => c.isFavorite)
        }

        // Sort
        list.sort((a, b) => {
            // Favorites always first
            if (a.isFavorite && !b.isFavorite) return -1
            if (!a.isFavorite && b.isFavorite) return 1

            switch (sortBy.value) {
                case 'name':
                    return (a.name || a.id).localeCompare(b.name || b.id)
                case 'size':
                    return b.totalSize - a.totalSize
                case 'date':
                default:
                    return new Date(b.createdAt || 0).getTime() - new Date(a.createdAt || 0).getTime()
            }
        })

        return list
    })

    // Format helpers
    function formatSize(bytes: number): string {
        return formatContextSize(bytes)
    }

    function formatTimestamp(ts: string): string {
        return formatTs(ts)
    }

    // Refresh
    async function refresh() {
        try {
            await contextStore.listProjectContexts()
        } catch (error) {
            console.error('[ContextList] Failed to refresh:', error)
        }
    }

    // Selection
    function toggleSelect(contextId: string) {
        if (selectedContexts.value.has(contextId)) {
            selectedContexts.value.delete(contextId)
        } else {
            selectedContexts.value.add(contextId)
        }
        selectedContexts.value = new Set(selectedContexts.value)
    }

    function clearSelection() {
        selectedContexts.value.clear()
        selectedContexts.value = new Set(selectedContexts.value)
    }

    // Favorite
    function toggleFavorite(contextId: string) {
        contextStore.toggleFavorite(contextId)
    }

    // Rename
    function startRename(context: ContextSummary) {
        editingId.value = context.id
        editingName.value = context.name || context.id
        nextTick(() => {
            renameInput.value?.focus()
            renameInput.value?.select()
        })
    }

    // Keyboard shortcuts handler
    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Delete' && selectedContexts.value.size > 0) {
            e.preventDefault()
            deleteSelected()
        }
        if (e.key === 'a' && (e.ctrlKey || e.metaKey) && filteredContexts.value.length > 0) {
            e.preventDefault()
            filteredContexts.value.forEach(c => selectedContexts.value.add(c.id))
            selectedContexts.value = new Set(selectedContexts.value)
        }
        if (e.key === 'c' && (e.ctrlKey || e.metaKey) && selectedContexts.value.size > 0) {
            e.preventDefault()
            copySelectedContext()
        }
        if (e.key === 'm' && (e.ctrlKey || e.metaKey) && selectedContexts.value.size >= 2) {
            e.preventDefault()
            mergeSelected()
        }
        if (e.key === 'Escape') {
            clearSelection()
            editingId.value = null
        }
    }

    function saveRename(contextId: string) {
        if (editingName.value.trim()) {
            contextStore.renameContext(contextId, editingName.value.trim())
            uiStore.addToast(t('context.renamed'), 'success')
        }
        cancelRename()
    }

    function cancelRename() {
        editingId.value = null
        editingName.value = ''
    }

    // Copy
    async function copyContext(contextId: string) {
        try {
            if (contextStore.contextId !== contextId) {
                await contextStore.loadContextContent(contextId, 0, 0)
            }
            const fullContent = await contextStore.getFullContextContent()
            await navigator.clipboard.writeText(fullContent)
            uiStore.addToast(t('toast.contextCopied'), 'success')
        } catch (error) {
            console.error('[ContextList] Failed to copy:', error)
            uiStore.addToast(t('toast.copyError'), 'error')
        }
    }

    async function copySelectedContext() {
        if (selectedContexts.value.size === 1) {
            const contextId = [...selectedContexts.value][0]
            await copyContext(contextId)
        } else if (selectedContexts.value.size > 1) {
            try {
                const contents: string[] = []
                for (const ctxId of selectedContexts.value) {
                    if (contextStore.contextId !== ctxId) {
                        await contextStore.loadContextContent(ctxId, 0, 0)
                    }
                    const content = await contextStore.getFullContextContent()
                    contents.push(content)
                }
                const merged = contents.join('\n\n' + '='.repeat(80) + '\n\n')
                await navigator.clipboard.writeText(merged)
                uiStore.addToast(t('toast.contextCopied'), 'success')
            } catch (error) {
                console.error('[ContextList] Failed to copy multiple:', error)
                uiStore.addToast(t('toast.copyError'), 'error')
            }
        }
    }

    // Duplicate
    async function duplicateContext(contextId: string) {
        try {
            const newId = await contextStore.duplicateContext(contextId)
            if (newId) {
                uiStore.addToast(t('context.duplicated'), 'success')
            }
        } catch (error) {
            console.error('[ContextList] Failed to duplicate:', error)
        }
    }

    // Export
    async function exportContext(contextId: string) {
        try {
            await contextStore.exportContext(contextId)
        } catch (error) {
            console.error('[ContextList] Export failed:', error)
        }
    }

    // Delete
    function confirmDelete(context: ContextSummary) {
        deleteModal.contextId = context.id
        deleteModal.contextIds = []
        deleteModal.message = t('context.confirmDeleteMessage').replace('{name}', context.name || context.id)
        deleteModal.show = true
    }

    function deleteSelected() {
        if (selectedContexts.value.size === 0) return

        deleteModal.contextId = null
        deleteModal.contextIds = [...selectedContexts.value]
        deleteModal.message = t('context.confirmDeleteMultiple').replace('{count}', String(selectedContexts.value.size))
        deleteModal.show = true
    }

    async function executeDelete() {
        try {
            if (deleteModal.contextId) {
                await contextStore.deleteContext(deleteModal.contextId)
                uiStore.addToast(t('context.deleted'), 'success')
            } else if (deleteModal.contextIds.length > 0) {
                for (const id of deleteModal.contextIds) {
                    await contextStore.deleteContext(id)
                }
                selectedContexts.value.clear()
                uiStore.addToast(t('context.deleted'), 'success')
            }
        } catch (error) {
            console.error('[ContextList] Delete failed:', error)
        } finally {
            deleteModal.show = false
            deleteModal.contextId = null
            deleteModal.contextIds = []
        }
    }

    function closeDeleteModal() {
        deleteModal.show = false
    }

    // Merge
    async function mergeSelected() {
        if (selectedContexts.value.size < 2) {
            uiStore.addToast(t('context.selectToMerge'), 'warning')
            return
        }

        try {
            const ids = [...selectedContexts.value]
            const newId = await contextStore.mergeContexts(ids)
            if (newId) {
                selectedContexts.value.clear()
                selectedContexts.value = new Set(selectedContexts.value)
                uiStore.addToast(t('context.merged'), 'success')
            }
        } catch (error) {
            console.error('[ContextList] Failed to merge:', error)
        }
    }

    // Drag & drop
    function handleDragStart(e: DragEvent, index: number) {
        dragIndex.value = index
        if (e.dataTransfer) {
            e.dataTransfer.effectAllowed = 'move'
            e.dataTransfer.setData('text/plain', String(index))
        }
    }

    function handleDragOver(e: DragEvent, index: number) {
        e.preventDefault()
        if (e.dataTransfer) {
            e.dataTransfer.dropEffect = 'move'
        }
        dragOverIndex.value = index
    }

    function handleDragLeave() {
        dragOverIndex.value = null
    }

    function handleDrop(e: DragEvent, toIndex: number) {
        e.preventDefault()
        if (dragIndex.value !== null && dragIndex.value !== toIndex) {
            contextStore.reorderContexts(dragIndex.value, toIndex)
        }
        dragIndex.value = null
        dragOverIndex.value = null
    }

    function handleDragEnd() {
        dragIndex.value = null
        dragOverIndex.value = null
    }

    return {
        // State
        showSettings,
        searchQuery,
        sortBy,
        showFavoritesOnly,
        selectedContexts,
        editingId,
        editingName,
        renameInput,
        dragIndex,
        dragOverIndex,
        deleteModal,
        storageSettings,
        // Computed
        filteredContexts,
        // Format helpers
        formatSize,
        formatTimestamp,
        // Actions
        refresh,
        toggleSelect,
        clearSelection,
        toggleFavorite,
        startRename,
        saveRename,
        cancelRename,
        copyContext,
        copySelectedContext,
        duplicateContext,
        exportContext,
        confirmDelete,
        deleteSelected,
        executeDelete,
        closeDeleteModal,
        mergeSelected,
        // Drag & drop
        handleDragStart,
        handleDragOver,
        handleDragLeave,
        handleDrop,
        handleDragEnd,
        // Keyboard
        handleKeydown
    }
}
