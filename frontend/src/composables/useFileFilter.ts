/**
 * useFileFilter - File filtering by extensions and token weight
 * Manages include/exclude extension filters, weight filters, and sorting
 */

import { TOKEN_THRESHOLDS } from '@/config/constants'
import { useSettingsStore } from '@/stores/settings.store'
import type { FileNode } from '@/types/domain'
import { filterTreeByExtensions } from '@/utils/fileTreeUtils'
import { computed, ref, type Ref } from 'vue'

export type WeightFilterLevel = 'none' | 'medium' | 'heavy' | 'critical'

export interface UseFileFilterOptions {
    nodes: Ref<FileNode[]>
    getAllFilesInNode?: (node: FileNode) => string[]
}

/**
 * Calculate token count for a file
 */
function getFileTokens(node: FileNode): number {
    if (node.isDir || !node.size) return 0
    return Math.round(node.size / TOKEN_THRESHOLDS.BYTES_PER_TOKEN)
}

/**
 * Get weight threshold value for a filter level
 */
function getWeightThreshold(level: WeightFilterLevel): number {
    switch (level) {
        case 'medium': return TOKEN_THRESHOLDS.MEDIUM
        case 'heavy': return TOKEN_THRESHOLDS.HEAVY
        case 'critical': return TOKEN_THRESHOLDS.CRITICAL
        default: return 0
    }
}

/**
 * Check if a node or its children contain files above the weight threshold
 */
function hasFilesAboveWeight(
    node: FileNode,
    threshold: number,
    getAllFilesInNode?: (node: FileNode) => string[]
): boolean {
    if (!node.isDir) {
        return getFileTokens(node) >= threshold
    }

    // For directories, check if any child file meets the threshold
    // Note: getAllFilesInNode returns paths, we need to check children directly
    if (node.children) {
        return node.children.some(child => hasFilesAboveWeight(child, threshold, getAllFilesInNode))
    }

    return false
}

/**
 * Filter tree by weight threshold
 */
function filterTreeByWeight(
    nodes: FileNode[],
    threshold: number,
    getAllFilesInNode?: (node: FileNode) => string[]
): FileNode[] {
    if (threshold === 0) return nodes

    return nodes.reduce<FileNode[]>((acc, node) => {
        if (!node.isDir) {
            // File: include if it meets threshold
            if (getFileTokens(node) >= threshold) {
                acc.push(node)
            }
        } else {
            // Directory: include if it has files meeting threshold
            if (hasFilesAboveWeight(node, threshold, getAllFilesInNode)) {
                const filteredChildren = node.children
                    ? filterTreeByWeight(node.children, threshold, getAllFilesInNode)
                    : []

                acc.push({
                    ...node,
                    children: filteredChildren
                })
            }
        }
        return acc
    }, [])
}

/**
 * Sort nodes with folders first, then alphabetically
 * Optimized to avoid unnecessary object creation
 */
function sortFoldersFirst(nodes: FileNode[]): FileNode[] {
    // Check if already sorted (optimization for repeated calls)
    let needsSort = false
    for (let i = 1; i < nodes.length; i++) {
        const prev = nodes[i - 1]
        const curr = nodes[i]
        // Check folder order
        if (!prev.isDir && curr.isDir) {
            needsSort = true
            break
        }
        // Check alphabetical order within same type
        if (prev.isDir === curr.isDir &&
            prev.name.toLowerCase().localeCompare(curr.name.toLowerCase()) > 0) {
            needsSort = true
            break
        }
    }

    const sorted = needsSort
        ? [...nodes].sort((a, b) => {
            if (a.isDir && !b.isDir) return -1
            if (!a.isDir && b.isDir) return 1
            return a.name.toLowerCase().localeCompare(b.name.toLowerCase())
        })
        : nodes

    // Recursively sort children only if they have children
    const withSortedChildren = sorted.map(node => {
        if (node.isDir && node.children && node.children.length > 1) {
            const sortedChildren = sortFoldersFirst(node.children)
            // Only create new object if children actually changed
            if (sortedChildren !== node.children) {
                return { ...node, children: sortedChildren }
            }
        }
        return node
    })

    // Avoid returning a new array if nothing changed
    const anyChildChanged = withSortedChildren.some((node, i) => node !== sorted[i])
    return anyChildChanged ? withSortedChildren : sorted
}


export function useFileFilter(options: UseFileFilterOptions) {
    const { nodes, getAllFilesInNode } = options
    const settingsStore = useSettingsStore()

    // State
    const filterExtensions = ref<string[]>([])
    const excludeExtensions = ref<string[]>([])
    const weightFilter = ref<WeightFilterLevel>('none')

    // Computed
    const filteredNodes = computed(() => {
        let result = nodes.value

        // Apply extension filters
        if (filterExtensions.value.length > 0 || excludeExtensions.value.length > 0) {
            result = filterTreeByExtensions(
                result,
                filterExtensions.value,
                excludeExtensions.value
            )
        }

        // Apply weight filter
        if (weightFilter.value !== 'none') {
            const threshold = getWeightThreshold(weightFilter.value)
            result = filterTreeByWeight(result, threshold, getAllFilesInNode)
        }

        // Apply System Filters (tests, node_modules, hidden)
        const s = settingsStore.settings.fileExplorer
        if (s.hideNodeModules || s.hideHiddenFiles || s.hideTestFiles) {
            const filterBySystem = (nodes: FileNode[]): FileNode[] => {
                return nodes.reduce<FileNode[]>((acc, node) => {
                    const name = node.name.toLowerCase()
                    
                    if (s.hideNodeModules && name === 'node_modules') return acc
                    if (s.hideHiddenFiles && node.name.startsWith('.') && node.name !== '.midas') return acc
                    if (s.hideTestFiles && (
                        name.includes('test') || 
                        name.includes('spec') || 
                        name === '__tests__'
                    )) return acc

                    if (node.isDir && node.children) {
                        const filteredChildren = filterBySystem(node.children)
                        acc.push({ ...node, children: filteredChildren })
                    } else {
                        acc.push(node)
                    }
                    return acc
                }, [])
            }
            result = filterBySystem(result)
        }

        // Apply folders first sorting
        if (settingsStore.settings.fileExplorer.foldersFirst) {
            result = sortFoldersFirst(result)
        }

        return result
    })

    // Actions
    function setFilterExtensions(include: string[], exclude: string[] = []) {
        filterExtensions.value = include
        excludeExtensions.value = exclude
    }

    function clearFilters() {
        filterExtensions.value = []
        excludeExtensions.value = []
        weightFilter.value = 'none'
    }

    function addIncludeExtension(ext: string) {
        if (!filterExtensions.value.includes(ext)) {
            filterExtensions.value.push(ext)
        }
    }

    function removeIncludeExtension(ext: string) {
        const index = filterExtensions.value.indexOf(ext)
        if (index > -1) {
            filterExtensions.value.splice(index, 1)
        }
    }

    function addExcludeExtension(ext: string) {
        if (!excludeExtensions.value.includes(ext)) {
            excludeExtensions.value.push(ext)
        }
    }

    function removeExcludeExtension(ext: string) {
        const index = excludeExtensions.value.indexOf(ext)
        if (index > -1) {
            excludeExtensions.value.splice(index, 1)
        }
    }

    function setWeightFilter(level: WeightFilterLevel) {
        weightFilter.value = level
    }

    function clearWeightFilter() {
        weightFilter.value = 'none'
    }

    return {
        // State
        filterExtensions,
        excludeExtensions,
        weightFilter,
        // Computed
        filteredNodes,
        // Actions
        setFilterExtensions,
        clearFilters,
        addIncludeExtension,
        removeIncludeExtension,
        addExcludeExtension,
        removeExcludeExtension,
        setWeightFilter,
        clearWeightFilter,
    }
}
