/**
 * File Store - Unified file tree management
 * Composes: useFileTree, useFileSelection, useFileFuzzySearch, useFileFilter, useFilePersistence
 */

import { useFileFilter } from '@/composables/useFileFilter'
import { useFileFuzzySearch } from '@/composables/useFileFuzzySearch'
import { useFilePersistence } from '@/composables/useFilePersistence'
import { useFileSelection } from '@/composables/useFileSelection'
import { useFileTree } from '@/composables/useFileTree'
import { useLogger } from '@/composables/useLogger'
import { useSettingsStore } from '@/stores/settings.store'
import { walkTree } from '@/utils/fileTreeUtils'
import { defineStore } from 'pinia'
import { computed, ref, triggerRef } from 'vue'
import { filesApi } from '../api/files.api'

const logger = useLogger('FileStore')

// Re-export FileNode for backward compatibility
export type { FileNode } from '@/types/domain'

export const useFileStore = defineStore('file', () => {
    // Compose: File Tree
    const tree = useFileTree()

    // Compose: File Selection (depends on tree)
    const selection = useFileSelection({
        findNode: tree.findNode,
        getAllFilesInNode: tree.getAllFilesInNode,
    })

    // Compose: File Fuzzy Search (depends on tree)
    const search = useFileFuzzySearch({
        flattenedNodes: tree.flattenedNodes,
    })

    // Compose: File Filter (depends on tree)
    const filter = useFileFilter({
        nodes: tree.nodes,
    })

    // Compose: Persistence (depends on tree and selection)
    const persistence = useFilePersistence({
        nodes: tree.nodes,
        selectedPaths: selection.selectedPaths,
        rootPath: tree.rootPath,
        findNode: tree.findNode,
    })

    // Additional state
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    // Settings
    const settingsStore = useSettingsStore()
    const autoSaveSelection = computed(() => settingsStore.settings.fileExplorer.autoSaveSelection)

    // Computed: Selected files total size
    const selectedFilesTotalSize = computed(() => {
        let totalSize = 0
        selection.selectedPaths.value.forEach((path) => {
            const node = tree.findNode(path)
            if (node && !node.isDir && node.size) {
                totalSize += node.size
            }
        })
        return totalSize
    })

    const estimatedTokenCount = computed(() => Math.round(selectedFilesTotalSize.value / 4))
    const estimatedContextSize = computed(() => selectedFilesTotalSize.value / (1024 * 1024))

    function getSelectedFilesSize(): number {
        return selectedFilesTotalSize.value
    }

    // Actions

    async function loadFileTree(projectPath: string, directory?: string) {
        isLoading.value = true
        error.value = null

        try {
            const targetPath = directory || projectPath
            const files = await filesApi.listFiles(targetPath, true, true)
            tree.setFileTree(files)

            // Set root path on first load
            if (!tree.rootPath.value) {
                tree.setRootPath(projectPath)

                // Load expanded state or auto-expand
                const loadedPaths = persistence.loadExpandedState()
                if (loadedPaths.length === 0) {
                    tree.autoExpand(3)
                }

                // Load saved selection
                const savedSelection = persistence.loadSelectionFromStorage(projectPath)
                if (savedSelection.length > 0) {
                    selection.selectMultiple(savedSelection)
                }
            } else if (directory) {
                tree.currentDirectory.value = directory
            }
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to load files'
            throw err
        } finally {
            isLoading.value = false
        }
    }

    function toggleSelect(path: string) {
        selection.toggleSelect(path)

        if (autoSaveSelection.value) {
            persistence.debouncedSaveSelection()
        }
    }

    function clearSelection() {
        selection.clearSelection()

        if (autoSaveSelection.value) {
            persistence.debouncedSaveSelection()
        }
    }

    function toggleExpand(pathOrCompact: string) {
        // Check if this is a compact mode batch operation
        if (pathOrCompact.startsWith('{')) {
            try {
                const { paths, expand } = JSON.parse(pathOrCompact) as { paths: string[], expand: boolean }
                for (const p of paths) {
                    if (expand) {
                        tree.expandPath(p)
                    } else {
                        tree.collapsePath(p)
                    }
                }
            } catch {
                // Fallback to simple toggle if JSON parse fails
                tree.toggleExpand(pathOrCompact)
            }
        } else {
            tree.toggleExpand(pathOrCompact)
        }
        persistence.debouncedSaveExpandedState()
    }

    function expandRecursive(path: string) {
        tree.expandRecursive(path)
        persistence.debouncedSaveExpandedState()
    }

    function collapseRecursive(path: string) {
        tree.collapseRecursive(path)
        persistence.debouncedSaveExpandedState()
    }

    function expandAll() {
        tree.expandAll()
        persistence.debouncedSaveExpandedState()
    }

    function collapseAll() {
        tree.collapseAll()
        persistence.debouncedSaveExpandedState()
    }

    function removeNode(path: string): boolean {
        const removed = tree.removeNode(path, selection.selectedPaths.value)
        if (removed) {
            triggerRef(selection.selectedPaths)
        }
        return removed
    }

    function selectByExtension(extension: string) {
        walkTree(tree.nodes.value, (node) => {
            if (!node.isDir && node.name.endsWith(extension)) {
                selection.selectedPaths.value.add(node.path)
            }
        })
        triggerRef(selection.selectedPaths)
    }

    async function refreshFileTree(): Promise<void> {
        filesApi.clearCache()
    }

    function resetStore() {
        tree.reset()
        selection.clearSelection()
        search.clearSearch()
        filter.clearFilters()
        error.value = null
        isLoading.value = false
        persistence.dispose()

        // Force garbage collection
        if (typeof window !== 'undefined' && 'gc' in window) {
            try {
                (window as unknown as { gc?: () => void }).gc?.()
            } catch {
                // Ignore
            }
        }
    }

    function getMemoryUsage(): number {
        let size = tree.getMemoryUsage()
        size += selection.selectedPaths.value.size * 100
        return size
    }

    function pruneUnusedBranches() {
        logger.debug('Pruning unused branches...')
    }

    return {
        // State (from tree)
        nodes: tree.nodes,
        rootPath: tree.rootPath,
        currentDirectory: tree.currentDirectory,
        directoryHistory: tree.directoryHistory,
        isLoading,
        error,

        // State (from selection)
        selectedPaths: selection.selectedPaths,

        // State (from search)
        searchQuery: search.searchQuery,

        // State (from filter)
        filterExtensions: filter.filterExtensions,
        excludeExtensions: filter.excludeExtensions,

        // Computed (from tree)
        projectName: tree.projectName,
        breadcrumbs: tree.breadcrumbs,
        flattenedNodes: tree.flattenedNodes,

        // Computed (from selection)
        hasSelectedFiles: selection.hasSelectedFiles,
        selectedCount: selection.selectedCount,
        selectedFilesList: selection.selectedFilesList,
        canUndoSelection: selection.canUndo,
        canRedoSelection: selection.canRedo,

        // Computed (from search)
        searchResults: search.searchResults,

        // Computed (from filter)
        filteredNodes: filter.filteredNodes,

        // Computed (local)
        estimatedTokenCount,
        estimatedContextSize,

        // Actions (tree)
        setFileTree: tree.setFileTree,
        loadFileTree,
        removeNode,
        toggleExpand,
        expandPath: tree.expandPath,
        collapsePath: tree.collapsePath,
        expandRecursive,
        collapseRecursive,
        expandAll,
        collapseAll,
        setRootPath: tree.setRootPath,
        getAvailableExtensions: tree.getAvailableExtensions,
        nodeExists: tree.nodeExists,
        getExpandedPaths: tree.getExpandedPaths,
        restoreExpandedPaths: tree.restoreExpandedPaths,
        autoExpandToFiles: () => tree.autoExpand(3),

        // Actions (selection)
        toggleSelect,
        selectPath: selection.selectPath,
        deselectPath: selection.deselectPath,
        selectMultiple: selection.selectMultiple,
        clearSelection,
        selectRecursive: selection.selectRecursive,
        deselectRecursive: selection.deselectRecursive,
        selectByExtension,
        undoSelection: selection.undoSelection,
        redoSelection: selection.redoSelection,

        // Actions (search)
        setSearchQuery: search.setSearchQuery,

        // Actions (filter)
        setFilterExtensions: filter.setFilterExtensions,

        // Actions (persistence)
        autoSaveSelection,
        saveSelectionToStorage: persistence.saveSelectionToStorage,
        loadSelectionFromStorage: persistence.loadSelectionFromStorage,
        clearSelectionHistory: persistence.clearSelectionHistory,
        getSelectionStats: persistence.getSelectionStats,
        saveExpandedState: persistence.saveExpandedState,
        loadExpandedState: persistence.loadExpandedState,

        // Actions (other)
        refreshFileTree,
        getSelectedFilesSize,
        resetStore,
        getMemoryUsage,
        pruneUnusedBranches,

        // Public utility methods for UI components
        getRecursiveFileCount: tree.getRecursiveFileCount,
        getAllFilesInNode: tree.getAllFilesInNode,
        getSelectedFileCountInNode: selection.getSelectedFileCountInNode,
        isDirectory: tree.isDirectory,
        getNodesByPaths: tree.getNodesByPaths,
    }
})
