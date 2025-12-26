/**
 * useFileTree - File tree state and operations
 * Manages tree structure, expansion, and caching
 */

import { FILE_TREE } from '@/config/constants'
import type { FileNode } from '@/types/domain'
import {
    autoExpandToFiles,
    buildNodePathCache,
    convertDomainNodes,
    findNode as findNodeUtil,
    generateFlattened,
    getAllFilesInNode as getAllFilesUtil,
    walkTree,
    type DomainNode,
} from '@/utils/fileTreeUtils'
import { computed, ref, shallowRef } from 'vue'

export function useFileTree() {
    // State
    const nodes = ref<FileNode[]>([])
    const rootPath = ref<string>('')
    const currentDirectory = ref<string>('')
    const directoryHistory = ref<string[]>([])

    // Caches - NOT REACTIVE to avoid Reactivity Storm
    const allFilesCache = new Map<string, string[]>()
    const nodePathCache = new Map<string, FileNode>()
    const flattenedNodesCache = shallowRef<FileNode[] | null>(null)

    // Computed
    const flattenedNodes = computed(() => {
        if (nodes.value.length === 0) return []
        if (!flattenedNodesCache.value) {
            flattenedNodesCache.value = generateFlattened(nodes.value, FILE_TREE.MAX_FLATTEN_DEPTH)
        }
        return flattenedNodesCache.value
    })

    const projectName = computed(() => {
        if (!rootPath.value) return 'Project'
        return rootPath.value.split(/[\\/]/).pop() || 'Project'
    })

    const breadcrumbs = computed(() => {
        if (!currentDirectory.value || !rootPath.value) return []
        const relative = currentDirectory.value
            .replace(rootPath.value, '')
            .replace(/^[\\/]+/, '')
        if (!relative) return [projectName.value]
        const segments = relative.split(/[\\/]/)
        return [projectName.value, ...segments]
    })

    // Actions
    function setFileTree(tree: FileNode[] | DomainNode[]) {
        const expandedPaths = getExpandedPaths()
        nodes.value = convertDomainNodes(tree as DomainNode[])

        // Clear caches when tree changes
        allFilesCache.clear()
        nodePathCache.clear()
        flattenedNodesCache.value = null

        // Build path cache for O(1) node lookups
        buildNodePathCache(nodes.value, nodePathCache)

        // Restore expanded state
        if (expandedPaths.length > 0) {
            restoreExpandedPaths(expandedPaths)
        }
    }

    function removeNode(path: string, selectedPaths: Set<string>): boolean {
        const removeFromTree = (tree: FileNode[], targetPath: string): boolean => {
            for (let i = 0; i < tree.length; i++) {
                if (tree[i].path === targetPath) {
                    // Deselect all files in this node before removing
                    if (tree[i].isDir) {
                        const filesToDeselect = getAllFilesInNode(tree[i])
                        filesToDeselect.forEach((p) => selectedPaths.delete(p))
                    } else {
                        selectedPaths.delete(targetPath)
                    }
                    tree.splice(i, 1)
                    return true
                }
                if (tree[i].children) {
                    if (removeFromTree(tree[i].children!, targetPath)) {
                        return true
                    }
                }
            }
            return false
        }

        const removed = removeFromTree(nodes.value, path)
        if (removed) {
            // Clear caches for removed path and its children
            nodePathCache.delete(path)
            allFilesCache.delete(path)

            // Clear parent caches too
            const parentPath =
                path.substring(0, path.lastIndexOf('/')) ||
                path.substring(0, path.lastIndexOf('\\'))
            if (parentPath) {
                allFilesCache.delete(parentPath)
            }
            flattenedNodesCache.value = null

            // Force Vue to detect the change
            nodes.value = [...nodes.value]
        }
        return removed
    }

    function findNode(path: string): FileNode | null {
        return findNodeUtil(nodes.value, path, nodePathCache)
    }

    function nodeExists(path: string): boolean {
        return findNode(path) !== null
    }

    function toggleExpand(path: string) {
        const node = findNode(path)
        if (node && node.isDir) {
            node.isExpanded = !node.isExpanded
            // Note: We don't invalidate flattenedNodesCache here for performance
            // The cache is only used for search, not for rendering expanded state
        }
    }

    function expandPath(path: string) {
        const node = findNode(path)
        if (node && node.isDir) {
            node.isExpanded = true
        }
    }

    function collapsePath(path: string) {
        const node = findNode(path)
        if (node && node.isDir) {
            node.isExpanded = false
        }
    }

    function expandRecursive(path: string) {
        const node = findNode(path)
        if (!node || !node.isDir) return

        const expandNode = (n: FileNode) => {
            if (n.isDir) {
                n.isExpanded = true
                if (n.children) {
                    n.children.forEach(expandNode)
                }
            }
        }
        expandNode(node)
    }

    function collapseRecursive(path: string) {
        const node = findNode(path)
        if (!node || !node.isDir) return

        const collapseNode = (n: FileNode) => {
            if (n.isDir) {
                n.isExpanded = false
                if (n.children) {
                    n.children.forEach(collapseNode)
                }
            }
        }
        collapseNode(node)
    }

    function expandAll() {
        walkTree(nodes.value, (node) => {
            if (node.isDir) {
                node.isExpanded = true
            }
        })
    }

    function collapseAll() {
        walkTree(nodes.value, (node) => {
            if (node.isDir) {
                node.isExpanded = false
            }
        })
    }

    function getExpandedPaths(): string[] {
        const expanded: string[] = []
        walkTree(nodes.value, (node) => {
            if (node.isDir && node.isExpanded) {
                expanded.push(node.path)
            }
        })
        return expanded
    }

    function restoreExpandedPaths(paths: string[]) {
        const pathSet = new Set(paths)
        walkTree(nodes.value, (node) => {
            if (node.isDir) {
                node.isExpanded = pathSet.has(node.path)
            }
        })
    }

    function getAllFilesInNode(node: FileNode): string[] {
        return getAllFilesUtil(node, allFilesCache)
    }

    function getRecursiveFileCount(node: FileNode): number {
        if (!node.isDir) return 0
        return getAllFilesInNode(node).length
    }

    function isDirectory(path: string): boolean {
        const node = findNode(path)
        return node ? node.isDir : false
    }

    function getNodesByPaths(paths: string[]): Map<string, FileNode> {
        const result = new Map<string, FileNode>()
        const pathSet = new Set(paths)

        function walk(nodeList: FileNode[]) {
            for (const node of nodeList) {
                if (pathSet.has(node.path)) {
                    result.set(node.path, node)
                    if (result.size === paths.length) return
                }
                if (node.children) {
                    walk(node.children)
                }
            }
        }

        walk(nodes.value)
        return result
    }

    function getAvailableExtensions(): string[] {
        const extensions = new Set<string>()

        function collectExtensions(nodeList: FileNode[]) {
            for (const node of nodeList) {
                if (!node.isDir && node.name.includes('.')) {
                    const ext = '.' + node.name.split('.').pop()
                    extensions.add(ext)
                }
                if (node.children) {
                    collectExtensions(node.children)
                }
            }
        }

        collectExtensions(nodes.value)
        return Array.from(extensions).sort()
    }

    function setRootPath(path: string) {
        rootPath.value = path
        currentDirectory.value = path
        directoryHistory.value = []
    }

    function autoExpand(maxDepth: number = 3) {
        autoExpandToFiles(nodes.value, maxDepth)
    }

    function getMemoryUsage(): number {
        let size = 0
        let nodeCount = 0

        const countNodes = (nodeList: FileNode[]) => {
            nodeCount += nodeList.length
            nodeList.forEach((node) => {
                if (node.children) countNodes(node.children)
            })
        }
        countNodes(nodes.value)

        size += nodeCount * 250
        size += (flattenedNodesCache.value?.length || 0) * 200
        size += allFilesCache.size * 150
        size += nodePathCache.size * 50

        return size
    }

    function clearCaches() {
        allFilesCache.clear()
        nodePathCache.clear()
        flattenedNodesCache.value = null
    }

    function reset() {
        nodes.value = []
        rootPath.value = ''
        currentDirectory.value = ''
        directoryHistory.value = []
        clearCaches()
    }

    return {
        // State
        nodes,
        rootPath,
        currentDirectory,
        directoryHistory,
        // Computed
        flattenedNodes,
        projectName,
        breadcrumbs,
        // Actions
        setFileTree,
        removeNode,
        findNode,
        nodeExists,
        toggleExpand,
        expandPath,
        collapsePath,
        expandRecursive,
        collapseRecursive,
        expandAll,
        collapseAll,
        getExpandedPaths,
        restoreExpandedPaths,
        getAllFilesInNode,
        getRecursiveFileCount,
        isDirectory,
        getNodesByPaths,
        getAvailableExtensions,
        setRootPath,
        autoExpand,
        getMemoryUsage,
        clearCaches,
        reset,
    }
}
