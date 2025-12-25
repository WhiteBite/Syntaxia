/**
 * Main composable for FileExplorer component
 * Composes smaller composables for search, quicklook, ignore rules
 */
import { useContextMenu } from '@/composables/useContextMenu'
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useContextStore } from '@/features/context'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref, watch } from 'vue'
import { filesApi } from '../api/files.api'
import { copyToClipboard, getRelativePath } from '../lib/file-utils'
import { useFileStore, type FileNode } from '../model/file.store'
import { useFileSearch } from './useFileSearch'
import { useIgnoreRules } from './useIgnoreRules'
import { useQuickLook } from './useQuickLook'

const logger = useLogger('FileExplorer')

export function useFileExplorer() {
    const { t } = useI18n()
    const fileStore = useFileStore()
    const contextStore = useContextStore()
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const settingsStore = useSettingsStore()
    const contextMenu = useContextMenu()

    // Compose smaller composables
    const search = useFileSearch()
    const ignoreRules = useIgnoreRules()

    // QuickLook - uses injected hoveredFile from parent component
    // NOTE: provideHoveredFile() must be called in the parent component (FileExplorer.vue)
    // before this composable is used
    const quickLook = useQuickLook()

    // UI state
    const showSettings = ref(false)
    const filterExtensions = ref<string[]>([])
    const totalFileCount = ref(0)

    // Computed
    const availableExtensions = computed(() => fileStore.getAvailableExtensions())

    const hasSelectionHistory = computed(() => {
        if (!projectStore.currentPath) return false
        const stats = fileStore.getSelectionStats()
        return stats[projectStore.currentPath] > 0
    })

    const selectionHistoryCount = computed(() => {
        if (!projectStore.currentPath) return 0
        const stats = fileStore.getSelectionStats()
        return stats[projectStore.currentPath] || 0
    })

    const selectionProgress = computed(() => {
        if (totalFileCount.value === 0) return 0
        return Math.round((fileStore.selectedCount / totalFileCount.value) * 100)
    })

    // Watch nodes to update total file count
    watch(() => fileStore.nodes, (newNodes) => {
        let count = 0
        const countFiles = (nodes: FileNode[]) => {
            nodes.forEach(node => {
                if (!node.isDir) count++
                if (node.children) countFiles(node.children)
            })
        }
        countFiles(newNodes)
        totalFileCount.value = count
    }, { immediate: true, deep: false })

    // Actions
    function handleToggleSelect(path: string) {
        fileStore.toggleSelect(path)
    }

    function handleToggleExpand(path: string) {
        fileStore.toggleExpand(path)
    }

    async function handleRefresh() {
        if (!projectStore.currentPath) return
        try {
            fileStore.clearSelection()
            await fileStore.loadFileTree(projectStore.currentPath)
            uiStore.addToast(t('toast.refreshed'), 'success')
        } catch (error) {
            logger.error('Failed to refresh file tree:', error)
            uiStore.addToast(t('toast.refreshError'), 'error')
        }
    }

    async function handleRefreshPreserveState() {
        if (!projectStore.currentPath) return
        try {
            const expandedPaths = fileStore.getExpandedPaths()
            const selectedPaths = fileStore.selectedFilesList
            await apiService.clearFileTreeCache()
            filesApi.clearCache()
            await fileStore.loadFileTree(projectStore.currentPath)
            fileStore.restoreExpandedPaths(expandedPaths)
            for (const path of selectedPaths) {
                if (fileStore.nodeExists(path)) {
                    fileStore.toggleSelect(path)
                }
            }
        } catch (error) {
            logger.error('Failed to refresh file tree:', error)
            uiStore.addToast(t('toast.refreshError'), 'error')
        }
    }

    function handleFilterUpdate(selected: string[]) {
        fileStore.setFilterExtensions(selected)
    }

    async function handleSettingsChange() {
        try {
            const dto = await apiService.getSettings()
            dto.useGitignore = settingsStore.settings.fileExplorer.useGitignore
            dto.useCustomIgnore = settingsStore.settings.fileExplorer.useCustomIgnore
            await apiService.saveSettings(JSON.stringify(dto))
            if (projectStore.currentPath) {
                await handleRefresh()
            }
        } catch (error) {
            logger.error('Failed to save settings:', error)
            uiStore.addToast('Failed to save settings', 'error')
        }
    }

    function restorePreviousSelection() {
        if (projectStore.currentPath) {
            fileStore.loadSelectionFromStorage(projectStore.currentPath)
            uiStore.addToast('Selection restored', 'success')
        }
    }

    function handleContextMenuShow(node: FileNode, event: MouseEvent) {
        contextMenu.show(node, event)
    }

    /**
     * Handle global keydown events
     * Uses hoveredFile from provide/inject instead of window globals
     */
    function handleGlobalKeydown(event: KeyboardEvent) {
        if (event.key !== ' ') return

        const target = event.target as HTMLElement
        if (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA' || target.isContentEditable) {
            return
        }

        // Use QuickLook's spacebar handler
        if (quickLook.handleSpacebarPreview()) {
            event.preventDefault()
        }
    }

    async function handleContextMenuAction(payload: { type: string; node: FileNode }) {
        const { type, node } = payload
        try {
            switch (type) {
                case 'quickLook':
                    if (!node.isDir) quickLook.open(node.path)
                    break
                case 'selectAll':
                    if (node.isDir) {
                        fileStore.selectRecursive(node.path)
                        uiStore.addToast('All files selected in folder', 'success')
                    }
                    break
                case 'deselectAll':
                    if (node.isDir) {
                        fileStore.deselectRecursive(node.path)
                        uiStore.addToast('All files deselected in folder', 'success')
                    }
                    break
                case 'copyPath':
                    await copyToClipboard(node.path)
                    uiStore.addToast('Path copied to clipboard', 'success')
                    break
                case 'copyRelativePath': {
                    const relativePath = projectStore.currentPath
                        ? getRelativePath(node.path, projectStore.currentPath)
                        : node.path
                    await copyToClipboard(relativePath)
                    uiStore.addToast('Relative path copied to clipboard', 'success')
                    break
                }
                case 'addToCustomIgnore':
                    await ignoreRules.addToIgnore(node)
                    break
                case 'removeFromIgnore':
                    if (await ignoreRules.removeFromIgnore(node)) {
                        await handleRefreshPreserveState()
                    }
                    break
                case 'expandAll':
                    if (node.isDir) {
                        fileStore.expandRecursive(node.path)
                        uiStore.addToast('Expanded all folders', 'success')
                    }
                    break
                case 'collapseAll':
                    if (node.isDir) {
                        fileStore.collapseRecursive(node.path)
                        uiStore.addToast('Collapsed all folders', 'success')
                    }
                    break
            }
        } catch (error) {
            logger.error('Context menu action failed:', error)
            uiStore.addToast('Action failed', 'error')
        }
    }

    /**
     * Initialize the file explorer
     * Must be called in component setup (onMounted)
     */
    function initialize() {
        // Setup global keyboard handler
        window.addEventListener('keydown', handleGlobalKeydown)

        // Load file tree
        if (projectStore.currentPath) {
            fileStore.loadFileTree(projectStore.currentPath).catch(error => {
                logger.error('Failed to load file tree:', error)
                uiStore.addToast('Failed to load project files.', 'error')
            })
        } else {
            uiStore.addToast('No project selected', 'warning')
        }
    }

    function cleanup() {
        window.removeEventListener('keydown', handleGlobalKeydown)
    }

    function setupWatchers() {
        watch(() => projectStore.currentPath, async (newPath, oldPath) => {
            if (newPath && newPath !== oldPath) {
                fileStore.clearSelection()
                contextStore.clearContext()
                try {
                    await fileStore.loadFileTree(newPath)
                } catch (error) {
                    logger.error('Failed to load file tree:', error)
                    uiStore.addToast('Failed to load project files', 'error')
                }
            }
        })
    }

    return {
        // Search (delegated to useFileSearch)
        searchQuery: search.query,
        handleSearch: search.handleSearch,
        clearSearch: search.clear,

        // QuickLook
        quickLookVisible: quickLook.isVisible,
        quickLookPath: quickLook.currentPath,
        handleQuickLook: (path: string) => quickLook.toggle(path),
        handleAddToContext: (path: string) => quickLook.addToContext(path),

        // UI state
        showSettings,
        filterExtensions,
        contextMenu,

        // Computed
        availableExtensions,
        hasSelectionHistory,
        selectionHistoryCount,
        totalFileCount,
        selectionProgress,

        // Actions
        handleToggleSelect,
        handleToggleExpand,
        handleRefresh,
        handleFilterUpdate,
        handleSettingsChange,
        restorePreviousSelection,
        handleContextMenuShow,
        handleContextMenuAction,

        // Lifecycle
        initialize,
        cleanup,
        setupWatchers,
    }
}
