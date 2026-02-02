/**
 * File Store - Unified file tree management
 * Composes: useFileTree, useFileSelection, useFileFuzzySearch, useFileFilter, useFilePersistence
 */

import { useDependencyGraph } from '@/composables/useDependencyGraph'
import { useFileFilter } from '@/composables/useFileFilter'
import { useFileFuzzySearch } from '@/composables/useFileFuzzySearch'
import { useFilePersistence } from '@/composables/useFilePersistence'
import { useFileSelection } from '@/composables/useFileSelection'
import { useFileTree } from '@/composables/useFileTree'
import { useLogger } from '@/composables/useLogger'
import { useSettingsStore } from '@/stores/settings.store'
import type { FileNode } from '@/types/domain'
import { walkTree } from '@/utils/fileTreeUtils'
import { defineStore } from 'pinia'
import { computed, ref, triggerRef } from 'vue'
import { filesApi } from '../api/files.api'

const logger = useLogger('FileStore')

// Re-export types for backward compatibility
export type { WeightFilterLevel } from '@/composables/useFileFilter'
export type { SelectionPreset } from '@/composables/useFilePersistence'
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
        rootPath: computed(() => tree.rootPath.value),
    })

    // Compose: File Filter (depends on tree)
    const filter = useFileFilter({
        nodes: tree.nodes,
        getAllFilesInNode: tree.getAllFilesInNode,
    })

    // Selected Only mode: show only selected files in flat list
    const isSelectedOnlyMode = ref(false)

    // Compose: Persistence (depends on tree and selection)
    const persistence = useFilePersistence({
        nodes: tree.nodes,
        selectedPaths: selection.selectedPaths,
        rootPath: tree.rootPath,
        findNode: tree.findNode,
    })

    // Compose: Dependency Graph
    const dependencyGraph = useDependencyGraph()

    // State for range selection
    const lastSelectedPath = ref<string | null>(null)

    // Additional state
    const isLoading = ref(false)

    const error = ref<string | null>(null)

    // Zen mode: dim unselected nodes
    const isZenMode = ref(false)

    // Keyboard navigation: focused path for roving tabindex
    const focusedPath = ref<string | null>(null)

    // Solo Expansion mode: only one folder expanded per level (accordion)
    const isSoloExpansionMode = ref(false)

    // Folder Focus mode: isolate a single folder as temporary root
    const focusedFolderPath = ref<string | null>(null)

    // Quick Open modal state
    const isQuickOpenModalVisible = ref(false)

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

    // Flat search mode: when search query is active, show flat list without hierarchy
    const isFlatSearchMode = computed(() => {
        return search.searchQuery.value.trim().length > 0
    })

    // Computed: Filtered nodes with selected-only mode support
    const filteredNodesWithSelectedOnly = computed(() => {
        // If selected-only mode is active, show only selected files in flat list
        if (isSelectedOnlyMode.value) {
            return tree.flattenedNodes.value
                .filter(node => !node.isDir && selection.isSelected(node.path))
                .map(node => ({
                    ...node,
                    depth: 0, // Flat display
                    relativePath: tree.rootPath.value
                        ? node.path.replace(tree.rootPath.value + '/', '')
                        : node.path
                }))
        }

        // Otherwise, use normal filtered nodes from filter composable
        return filter.filteredNodes.value
    })

    // Computed: Get dependencies for selected files
    // Returns a map of file paths that are dependencies -> array of selected files that depend on them
    const selectedFileDependencies = computed(() => {
        const deps = new Map<string, string[]>()

        for (const selectedPath of selection.selectedPaths.value) {
            const fileDeps = dependencyGraph.findDependencies(
                selectedPath,
                getAllFileNodes()
            )

            for (const dep of fileDeps) {
                if (!deps.has(dep.path)) {
                    deps.set(dep.path, [])
                }
                deps.get(dep.path)!.push(selectedPath)
            }
        }

        return deps
    })

    // Computed: Get all dependencies (incoming + outgoing) for each file
    const allFileDependencies = computed(() => {
        const allNodes = getAllFileNodes()
        const depsMap = new Map<string, { incoming: string[], outgoing: string[] }>()

        for (const selectedPath of selection.selectedPaths.value) {
            const { incoming, outgoing } = dependencyGraph.findAllDependencies(selectedPath, allNodes)

            // Track outgoing dependencies (files this file depends on)
            for (const dep of outgoing) {
                if (!depsMap.has(dep.path)) {
                    depsMap.set(dep.path, { incoming: [], outgoing: [] })
                }
                depsMap.get(dep.path)!.incoming.push(selectedPath)
            }

            // Track incoming dependencies (files that depend on this file)
            for (const dep of incoming) {
                if (!depsMap.has(dep.path)) {
                    depsMap.set(dep.path, { incoming: [], outgoing: [] })
                }
                depsMap.get(dep.path)!.outgoing.push(selectedPath)
            }
        }

        return depsMap
    })

    function getAllFileNodes(): FileNode[] {
        const result: FileNode[] = []
        walkTree(tree.nodes.value, (node) => {
            if (!node.isDir) result.push(node)
        })
        return result
    }

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
                    tree.autoExpand(1) // Changed from 3 to 1 - show only root level folders
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
        lastSelectedPath.value = path

        if (autoSaveSelection.value) {
            persistence.debouncedSaveSelection()
        }
    }

    function batchToggleSelect(paths: string[], shouldSelect: boolean) {
        for (const path of paths) {
            if (shouldSelect) {
                selection.selectPath(path)
            } else {
                selection.deselectPath(path)
            }
        }

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
                tree.toggleExpand(pathOrCompact, isSoloExpansionMode.value)
            }
        } else {
            tree.toggleExpand(pathOrCompact, isSoloExpansionMode.value)
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

    function toggleZenMode() {
        isZenMode.value = !isZenMode.value
    }

    function shouldDimNode(node: FileNode): boolean {
        if (!isZenMode.value) return false

        // Don't dim selected files
        if (!node.isDir && selection.selectedPaths.value.has(node.path)) {
            return false
        }

        // Don't dim folders that contain selected files
        if (node.isDir) {
            const allFiles = tree.getAllFilesInNode(node)
            const hasSelectedFiles = allFiles.some(path => selection.selectedPaths.value.has(path))
            if (hasSelectedFiles) return false
        }

        // Dim everything else
        return true
    }

    function toggleSelectedOnlyMode() {
        isSelectedOnlyMode.value = !isSelectedOnlyMode.value
    }

    function toggleSoloExpansionMode() {
        isSoloExpansionMode.value = !isSoloExpansionMode.value
    }

    function setFocusedPath(path: string | null) {
        focusedPath.value = path
    }

    function setFocusOnFolder(path: string) {
        const node = tree.findNode(path)
        if (node && node.isDir) {
            focusedFolderPath.value = path
            // Auto-expand the focused folder
            tree.expandPath(path)
        }
    }

    function clearFolderFocus() {
        focusedFolderPath.value = null
    }

    function selectRelated(path: string): number {
        // Get all file nodes from the tree
        const allFileNodes: FileNode[] = []
        walkTree(tree.nodes.value, (node) => {
            if (!node.isDir) allFileNodes.push(node)
        })

        // Find related files using dependency graph
        const relatedPaths = dependencyGraph.findRelatedFiles(path, allFileNodes)

        // Select all related files
        let selectedCount = 0
        for (const relatedPath of relatedPaths) {
            if (!selection.selectedPaths.value.has(relatedPath)) {
                selection.selectPath(relatedPath)
                selectedCount++
            }
        }

        if (autoSaveSelection.value && selectedCount > 0) {
            persistence.debouncedSaveSelection()
        }

        return selectedCount
    }

    function addDependencies(path: string): number {
        // Get all file nodes from the tree
        const allFileNodes: FileNode[] = []
        walkTree(tree.nodes.value, (node) => {
            if (!node.isDir) allFileNodes.push(node)
        })

        // Find all dependencies (incoming + outgoing)
        const { incoming, outgoing } = dependencyGraph.findAllDependencies(path, allFileNodes)
        const allDeps = [...incoming, ...outgoing]

        // Select all dependency files that aren't already selected
        let selectedCount = 0
        for (const dep of allDeps) {
            if (!selection.selectedPaths.value.has(dep.path)) {
                selection.selectPath(dep.path)
                selectedCount++
            }
        }

        if (autoSaveSelection.value && selectedCount > 0) {
            persistence.debouncedSaveSelection()
        }

        return selectedCount
    }

    function getPresets() {
        return persistence.getPresets()
    }

    function savePreset(name: string, description?: string) {
        persistence.savePreset(name, description)
    }

    function loadPreset(name: string): number {
        const count = persistence.loadPreset(name)
        triggerRef(selection.selectedPaths)

        if (autoSaveSelection.value && count > 0) {
            persistence.debouncedSaveSelection()
        }

        return count
    }

    function deletePreset(name: string) {
        persistence.deletePreset(name)
    }

    function clearAllPresets() {
        persistence.clearAllPresets()
    }

    function toggleQuickOpenModal() {
        isQuickOpenModalVisible.value = !isQuickOpenModalVisible.value
    }

    function openQuickOpenModal() {
        isQuickOpenModalVisible.value = true
    }

    function closeQuickOpenModal() {
        isQuickOpenModalVisible.value = false
    }

    return {
        // State (from tree)
        nodes: tree.nodes,
        rootPath: tree.rootPath,
        currentDirectory: tree.currentDirectory,
        directoryHistory: tree.directoryHistory,
        isLoading,
        error,
        isZenMode,
        isSelectedOnlyMode,
        isSoloExpansionMode,
        focusedPath,
        focusedFolderPath,
        isQuickOpenModalVisible,

        // State (from selection)
        selectedPaths: selection.selectedPaths,

        // State (from search)
        searchQuery: search.searchQuery,

        // State (from filter)
        filterExtensions: filter.filterExtensions,
        excludeExtensions: filter.excludeExtensions,
        weightFilter: filter.weightFilter,
        lastSelectedPath,

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
        filteredNodes: filteredNodesWithSelectedOnly,

        // Computed (local)
        estimatedTokenCount,
        estimatedContextSize,
        selectedFileDependencies,
        allFileDependencies,
        isFlatSearchMode,

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
        batchToggleSelect,
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
        setWeightFilter: filter.setWeightFilter,
        clearWeightFilter: filter.clearWeightFilter,

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
        toggleZenMode,
        shouldDimNode,
        toggleSelectedOnlyMode,
        toggleSoloExpansionMode,
        setFocusedPath,
        setFocusOnFolder,
        clearFolderFocus,
        selectRelated,
        addDependencies,

        // Actions (presets)
        getPresets,
        savePreset,
        loadPreset,
        deletePreset,
        clearAllPresets,
        toggleQuickOpenModal,
        openQuickOpenModal,
        closeQuickOpenModal,

        // Public utility methods for UI components
        findNode: tree.findNode,
        getRecursiveFileCount: tree.getRecursiveFileCount,
        getAllFilesInNode: tree.getAllFilesInNode,
        getSelectedFileCountInNode: selection.getSelectedFileCountInNode,
        isDirectory: tree.isDirectory,
        getNodesByPaths: tree.getNodesByPaths,
    }
})
