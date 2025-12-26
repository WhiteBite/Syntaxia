import { useI18n } from '@/composables/useI18n'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, onMounted, ref } from 'vue'

export interface ContextMenuState {
    visible: boolean
    x: number
    y: number
    project: { path: string; name: string } | null
}

export function useProjectSelector(emit: (e: 'opened', path: string) => void) {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const { t, locale, setLocale } = useI18n()

    const isDragging = ref(false)
    const contextMenu = ref<ContextMenuState>({
        visible: false,
        x: 0,
        y: 0,
        project: null
    })

    const recentProjects = computed(() => projectStore.recentProjects)

    onMounted(async () => {
        try {
            await projectStore.fetchRecentProjects()
        } catch {
            // Ignore fetch errors
        }
    })

    function toggleLanguage() {
        const newLocale = locale.value === 'ru' ? 'en' : 'ru'
        setLocale(newLocale)
        uiStore.addToast(
            newLocale === 'ru' ? 'Язык изменён на русский' : 'Language changed to English',
            'success'
        )
    }

    function onToggleAutoOpen(e: Event) {
        const target = e.target as HTMLInputElement
        projectStore.setAutoOpenLast(target.checked)
    }

    async function handleDrop(e: DragEvent) {
        isDragging.value = false
        const items = e.dataTransfer?.items
        if (!items) return

        for (let i = 0; i < items.length; i++) {
            const item = items[i]
            if (item.kind === 'file') {
                const entry = item.webkitGetAsEntry?.()
                if (entry?.isDirectory) {
                    const path = (entry as unknown as { fullPath: string }).fullPath
                    const success = await projectStore.openProjectByPath(path)
                    if (success) {
                        emit('opened', path)
                    }
                    return
                }
            }
        }
        uiStore.addToast('Please drop a folder, not a file', 'warning')
    }

    async function selectProject() {
        try {
            const dirPath = await apiService.selectDirectory()
            if (!dirPath || dirPath === '') return

            const success = await projectStore.openProjectByPath(dirPath)
            if (success && projectStore.currentPath) {
                emit('opened', projectStore.currentPath)
            }
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Unknown error'
            uiStore.addToast(`Failed to select directory: ${errorMessage}`, 'error')
        }
    }

    async function openRecentProject(path: string) {
        try {
            await projectStore.openProjectByPath(path)
            emit('opened', path)
        } catch {
            uiStore.addToast('Failed to open project', 'error')
        }
    }

    async function removeProject(path: string) {
        try {
            await projectStore.removeFromRecent(path)
            uiStore.addToast(t('welcome.projectRemoved'), 'success')
        } catch {
            // Ignore remove errors
        }
    }

    async function clearAllHistory() {
        if (confirm(t('welcome.confirmClearHistory'))) {
            projectStore.clearRecent()
            uiStore.addToast(t('welcome.historyCleared'), 'success')
        }
    }

    function showContextMenu(event: MouseEvent, project: { path: string; name: string }) {
        contextMenu.value = { visible: true, x: event.clientX, y: event.clientY, project }
    }

    function hideContextMenu() {
        contextMenu.value.visible = false
        contextMenu.value.project = null
    }

    async function copyProjectPath() {
        if (contextMenu.value.project) {
            try {
                await navigator.clipboard.writeText(contextMenu.value.project.path)
                uiStore.addToast(t('welcome.pathCopied'), 'success')
            } catch {
                // Ignore clipboard errors
            }
        }
        hideContextMenu()
    }

    async function removeProjectFromMenu() {
        if (contextMenu.value.project) {
            await removeProject(contextMenu.value.project.path)
        }
        hideContextMenu()
    }

    // Setup global click listener
    if (typeof window !== 'undefined') {
        window.addEventListener('click', hideContextMenu)
    }

    return {
        // State
        isDragging,
        contextMenu,
        recentProjects,
        locale,
        projectStore,

        // Methods
        toggleLanguage,
        onToggleAutoOpen,
        handleDrop,
        selectProject,
        openRecentProject,
        removeProject,
        clearAllHistory,
        showContextMenu,
        hideContextMenu,
        copyProjectPath,
        removeProjectFromMenu,
    }
}
