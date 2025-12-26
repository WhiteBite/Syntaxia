import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

export interface SandboxChange {
    path: string
    operation: 'create' | 'modify' | 'delete'
    diff: string
}

export const useSandboxStore = defineStore('sandbox', () => {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()

    // State
    const changes = ref<SandboxChange[]>([])
    const isLoading = ref(false)

    // Computed
    const hasChanges = computed(() => changes.value.length > 0)
    const changeCount = computed(() => changes.value.length)

    // Actions
    async function refresh(): Promise<void> {
        if (!projectStore.currentPath) return

        try {
            isLoading.value = true
            const result = await apiService.getSandboxChanges()
            changes.value = result.changes || []
        } catch (error) {
            // Sandbox might not have changes, that's ok
            changes.value = []
        } finally {
            isLoading.value = false
        }
    }

    async function applyAll(): Promise<boolean> {
        try {
            isLoading.value = true
            await apiService.applySandboxChanges()
            changes.value = []
            uiStore.addToast('Changes applied successfully', 'success')
            return true
        } catch (error) {
            uiStore.addToast('Failed to apply changes', 'error')
            return false
        } finally {
            isLoading.value = false
        }
    }

    async function discardAll(): Promise<void> {
        try {
            isLoading.value = true
            await apiService.discardSandboxChanges()
            changes.value = []
            uiStore.addToast('Changes discarded', 'info')
        } catch (error) {
            uiStore.addToast('Failed to discard changes', 'error')
        } finally {
            isLoading.value = false
        }
    }

    async function discardFile(path: string): Promise<void> {
        try {
            await apiService.discardSandboxFile(path)
            changes.value = changes.value.filter(c => c.path !== path)
        } catch (error) {
            uiStore.addToast('Failed to discard file', 'error')
        }
    }

    async function getDiff(path: string): Promise<string> {
        try {
            return await apiService.getSandboxDiff(path)
        } catch {
            return ''
        }
    }

    async function getAllDiffs(): Promise<string> {
        try {
            return await apiService.getSandboxAllDiffs()
        } catch {
            return ''
        }
    }

    return {
        // State
        changes,
        isLoading,
        // Computed
        hasChanges,
        changeCount,
        // Actions
        refresh,
        applyAll,
        discardAll,
        discardFile,
        getDiff,
        getAllDiffs
    }
})
