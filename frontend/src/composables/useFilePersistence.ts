/**
 * useFilePersistence - Selection and expanded state persistence
 * Saves/loads state to localStorage
 */

import { useLogger } from '@/composables/useLogger'
import type { FileNode } from '@/types/domain'
import { walkTree } from '@/utils/fileTreeUtils'
import { type Ref, type ShallowRef } from 'vue'

const logger = useLogger('FilePersistence')

export interface UseFilePersistenceOptions {
    nodes: Ref<FileNode[]>
    selectedPaths: ShallowRef<Set<string>>
    rootPath: Ref<string>
    findNode: (path: string) => FileNode | null
}

const SELECTION_PREFIX = 'file-selection-'
const EXPANDED_PREFIX = 'file-expanded-'
const MAX_SAVED_SELECTIONS = 100

export function useFilePersistence(options: UseFilePersistenceOptions) {
    const { nodes, selectedPaths, rootPath, findNode } = options

    // Debounce timers
    let saveExpandedStateTimer: ReturnType<typeof setTimeout> | null = null
    let saveSelectionTimer: ReturnType<typeof setTimeout> | null = null

    // Selection persistence
    function saveSelectionToStorage() {
        if (!rootPath.value) return

        try {
            const key = `${SELECTION_PREFIX}${rootPath.value}`
            const selection = Array.from(selectedPaths.value).slice(0, MAX_SAVED_SELECTIONS)
            localStorage.setItem(key, JSON.stringify(selection))
            logger.debug(`Saved selection: ${selection.length} files`)
        } catch (err) {
            console.warn('[FilePersistence] Failed to save selection:', err)
        }
    }

    function debouncedSaveSelection(delay: number = 300) {
        if (saveSelectionTimer) {
            clearTimeout(saveSelectionTimer)
        }
        saveSelectionTimer = setTimeout(() => {
            saveSelectionToStorage()
            saveSelectionTimer = null
        }, delay)
    }

    function loadSelectionFromStorage(projectPath: string): string[] {
        try {
            const key = `${SELECTION_PREFIX}${projectPath}`
            const saved = localStorage.getItem(key)
            if (saved) {
                const selection = JSON.parse(saved) as string[]
                logger.debug(`Loaded selection: ${selection.length} files`)
                return selection
            }
        } catch (err) {
            console.warn('[FilePersistence] Failed to load selection:', err)
        }
        return []
    }

    function clearSelectionHistory(projectPath?: string) {
        try {
            if (projectPath) {
                const key = `${SELECTION_PREFIX}${projectPath}`
                localStorage.removeItem(key)
            } else {
                const keys = Object.keys(localStorage).filter((k) =>
                    k.startsWith(SELECTION_PREFIX)
                )
                keys.forEach((k) => localStorage.removeItem(k))
            }
        } catch (err) {
            console.warn('[FilePersistence] Failed to clear selection history:', err)
        }
    }

    function getSelectionStats(): Record<string, number> {
        const stats: Record<string, number> = {}
        try {
            const keys = Object.keys(localStorage).filter((k) =>
                k.startsWith(SELECTION_PREFIX)
            )
            keys.forEach((key) => {
                const projectPath = key.replace(SELECTION_PREFIX, '')
                const saved = localStorage.getItem(key)
                if (saved) {
                    const selection = JSON.parse(saved) as string[]
                    stats[projectPath] = selection.length
                }
            })
        } catch (err) {
            console.warn('[FilePersistence] Failed to get selection stats:', err)
        }
        return stats
    }

    // Expanded state persistence
    function saveExpandedState() {
        if (!rootPath.value) return

        try {
            const expandedPaths: string[] = []
            walkTree(nodes.value, (node) => {
                if (node.isDir && node.isExpanded) {
                    expandedPaths.push(node.path)
                }
            })

            const key = `${EXPANDED_PREFIX}${rootPath.value}`
            localStorage.setItem(key, JSON.stringify(expandedPaths))
        } catch (err) {
            console.warn('[FilePersistence] Failed to save expanded state:', err)
        }
    }

    function debouncedSaveExpandedState(delay: number = 500) {
        if (saveExpandedStateTimer) {
            clearTimeout(saveExpandedStateTimer)
        }
        saveExpandedStateTimer = setTimeout(() => {
            saveExpandedState()
            saveExpandedStateTimer = null
        }, delay)
    }

    function loadExpandedState(): string[] {
        if (!rootPath.value) return []

        try {
            const key = `${EXPANDED_PREFIX}${rootPath.value}`
            const saved = localStorage.getItem(key)
            if (saved) {
                const expandedPaths = JSON.parse(saved) as string[]
                expandedPaths.forEach((path) => {
                    const node = findNode(path)
                    if (node && node.isDir) {
                        node.isExpanded = true
                    }
                })
                return expandedPaths
            }
        } catch (err) {
            console.warn('[FilePersistence] Failed to load expanded state:', err)
        }
        return []
    }

    function clearExpandedHistory(projectPath?: string) {
        try {
            if (projectPath) {
                const key = `${EXPANDED_PREFIX}${projectPath}`
                localStorage.removeItem(key)
            } else {
                const keys = Object.keys(localStorage).filter((k) =>
                    k.startsWith(EXPANDED_PREFIX)
                )
                keys.forEach((k) => localStorage.removeItem(k))
            }
        } catch (err) {
            console.warn('[FilePersistence] Failed to clear expanded history:', err)
        }
    }

    // Cleanup
    function dispose() {
        if (saveExpandedStateTimer) {
            clearTimeout(saveExpandedStateTimer)
            saveExpandedStateTimer = null
        }
        if (saveSelectionTimer) {
            clearTimeout(saveSelectionTimer)
            saveSelectionTimer = null
        }
    }

    return {
        // Selection persistence
        saveSelectionToStorage,
        debouncedSaveSelection,
        loadSelectionFromStorage,
        clearSelectionHistory,
        getSelectionStats,
        // Expanded state persistence
        saveExpandedState,
        debouncedSaveExpandedState,
        loadExpandedState,
        clearExpandedHistory,
        // Cleanup
        dispose,
    }
}
