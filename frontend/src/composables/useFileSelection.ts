/**
 * useFileSelection - File selection state and operations
 * Manages selected files, recursive selection, selection stats, and undo/redo
 */

import type { FileNode } from '@/types/domain'
import { computed, shallowRef, triggerRef } from 'vue'

export interface UseFileSelectionOptions {
    findNode: (path: string) => FileNode | null
    getAllFilesInNode: (node: FileNode) => string[]
}

const MAX_HISTORY = 15

export function useFileSelection(options: UseFileSelectionOptions) {
    const { findNode, getAllFilesInNode } = options

    // State - Use shallowRef to avoid deep reactivity on large sets
    const selectedPaths = shallowRef<Set<string>>(new Set())

    // Undo/Redo history
    const history = shallowRef<string[][]>([[]])
    const historyIndex = shallowRef(0)

    // Computed
    const hasSelectedFiles = computed(() => selectedPaths.value.size > 0)
    const selectedCount = computed(() => selectedPaths.value.size)
    const selectedFilesList = computed(() => Array.from(selectedPaths.value))
    const canUndo = computed(() => historyIndex.value > 0)
    const canRedo = computed(() => historyIndex.value < history.value.length - 1)

    // Save current state to history
    function saveToHistory() {
        const currentState = Array.from(selectedPaths.value)
        const newHistory = history.value.slice(0, historyIndex.value + 1)
        newHistory.push(currentState)

        // Limit history size
        if (newHistory.length > MAX_HISTORY) {
            newHistory.shift()
        } else {
            historyIndex.value++
        }
        history.value = newHistory
    }

    // Undo selection
    function undoSelection(): boolean {
        if (!canUndo.value) return false
        historyIndex.value--
        const prevState = history.value[historyIndex.value]
        selectedPaths.value = new Set(prevState)
        triggerRef(selectedPaths)
        return true
    }

    // Redo selection
    function redoSelection(): boolean {
        if (!canRedo.value) return false
        historyIndex.value++
        const nextState = history.value[historyIndex.value]
        selectedPaths.value = new Set(nextState)
        triggerRef(selectedPaths)
        return true
    }

    // Actions
    function toggleSelect(path: string) {
        const node = findNode(path)
        if (!node) return

        saveToHistory()

        if (node.isDir) {
            toggleSelectRecursiveInternal(path)
        } else {
            if (selectedPaths.value.has(path)) {
                selectedPaths.value.delete(path)
            } else {
                selectedPaths.value.add(path)
            }
        }
        triggerRef(selectedPaths)
    }

    function selectPath(path: string) {
        selectedPaths.value.add(path)
        triggerRef(selectedPaths)
    }

    function deselectPath(path: string) {
        selectedPaths.value.delete(path)
        triggerRef(selectedPaths)
    }

    function selectMultiple(paths: string[]) {
        paths.forEach((p) => selectedPaths.value.add(p))
        triggerRef(selectedPaths)
    }

    function clearSelection() {
        if (selectedPaths.value.size > 0) {
            saveToHistory()
        }
        selectedPaths.value.clear()
        triggerRef(selectedPaths)
    }

    function selectRecursive(path: string) {
        const node = findNode(path)
        if (!node) return

        const filePaths = getAllFilesInNode(node)
        filePaths.forEach((p) => selectedPaths.value.add(p))
        triggerRef(selectedPaths)
    }

    function deselectRecursive(path: string) {
        const node = findNode(path)
        if (!node) return

        const filePaths = getAllFilesInNode(node)
        filePaths.forEach((p) => selectedPaths.value.delete(p))
        triggerRef(selectedPaths)
    }

    // Internal version without history save (called from toggleSelect which already saves)
    function toggleSelectRecursiveInternal(path: string) {
        const node = findNode(path)
        if (!node || !node.isDir) return

        const childFilePaths = getAllFilesInNode(node)
        const anySelected = childFilePaths.some((filePath) =>
            selectedPaths.value.has(filePath)
        )

        if (anySelected) {
            childFilePaths.forEach((filePath) => selectedPaths.value.delete(filePath))
        } else {
            childFilePaths.forEach((filePath) => selectedPaths.value.add(filePath))
        }
    }

    function toggleSelectRecursive(path: string) {
        saveToHistory()
        toggleSelectRecursiveInternal(path)
        triggerRef(selectedPaths)
    }

    function selectByExtension(
        extension: string,
        nodes: FileNode[],
        walkTree: (tree: FileNode[], fn: (node: FileNode) => void) => void
    ) {
        walkTree(nodes, (node) => {
            if (!node.isDir && node.name.endsWith(extension)) {
                selectedPaths.value.add(node.path)
            }
        })
        triggerRef(selectedPaths)
    }

    function getSelectedFileCountInNode(node: FileNode): number {
        const allFiles = getAllFilesInNode(node)
        return allFiles.filter((filePath) => selectedPaths.value.has(filePath)).length
    }

    function isSelected(path: string): boolean {
        return selectedPaths.value.has(path)
    }

    function getSelectionState(
        node: FileNode
    ): 'none' | 'partial' | 'full' {
        if (!node.isDir) {
            return selectedPaths.value.has(node.path) ? 'full' : 'none'
        }

        if (selectedCount.value === 0) {
            return 'none'
        }

        const allFiles = getAllFilesInNode(node)
        if (allFiles.length === 0) {
            return selectedPaths.value.has(node.path) ? 'full' : 'none'
        }

        let selectedFileCount = 0
        for (const filePath of allFiles) {
            if (selectedPaths.value.has(filePath)) {
                selectedFileCount++
            }
        }

        if (selectedFileCount === 0) return 'none'
        if (selectedFileCount === allFiles.length) return 'full'
        return 'partial'
    }

    return {
        // State
        selectedPaths,
        // Computed
        hasSelectedFiles,
        selectedCount,
        selectedFilesList,
        canUndo,
        canRedo,
        // Actions
        toggleSelect,
        selectPath,
        deselectPath,
        selectMultiple,
        clearSelection,
        selectRecursive,
        deselectRecursive,
        toggleSelectRecursive,
        selectByExtension,
        getSelectedFileCountInNode,
        isSelected,
        getSelectionState,
        undoSelection,
        redoSelection,
    }
}
