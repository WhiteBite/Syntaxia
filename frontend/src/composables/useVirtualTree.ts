/**
 * useVirtualTree - Flatten tree for virtualization
 * Converts hierarchical tree to flat list of visible nodes
 */

import { useSettingsStore } from '@/stores/settings.store'
import type { FileNode } from '@/types/domain'
import { computed, type Ref } from 'vue'

export interface FlattenedNode {
    id: string // Flat key for RecycleScroller (same as node.path)
    node: FileNode
    depth: number
    isLast: boolean
    ancestorHasMoreSiblings: boolean[]
    displayName?: string // For compact mode (e.g. "src/components")
    relativePath?: string // For search mode (relative to root)
}

export interface UseVirtualTreeOptions {
    nodes: Ref<FileNode[]>
    isSelectedOnlyMode?: Ref<boolean>
    selectedPaths?: Ref<Set<string>>
    rootPath?: Ref<string>
    focusedFolderPath?: Ref<string | null>
}

/**
 * Flatten tree to visible nodes only (expanded folders)
 */
export function useVirtualTree(options: UseVirtualTreeOptions) {
    const { nodes, isSelectedOnlyMode, selectedPaths, rootPath, focusedFolderPath } = options
    const settingsStore = useSettingsStore()

    const flattenedVisibleNodes = computed<FlattenedNode[]>(() => {
        const result: FlattenedNode[] = []
        const isCompactEnabled = settingsStore.settings.fileExplorer.compactNestedFolders

        // Folder Focus Mode: show only the focused folder and its children
        if (focusedFolderPath?.value) {
            // Find the focused folder node
            function findFocusedNode(nodeList: FileNode[]): FileNode | null {
                for (const node of nodeList) {
                    if (focusedFolderPath && node.path === focusedFolderPath.value) {
                        return node
                    }
                    if (node.children) {
                        const found = findFocusedNode(node.children)
                        if (found) return found
                    }
                }
                return null
            }

            const focusedNode = findFocusedNode(nodes.value)
            if (focusedNode && focusedNode.isDir) {
                // Show the focused folder as root
                const focusedNodes = [focusedNode]

                function flatten(
                    nodeList: FileNode[],
                    depth: number,
                    ancestorHasMoreSiblings: boolean[]
                ) {
                    nodeList.forEach((node, index) => {
                        const isLast = index === nodeList.length - 1
                        let displayName = node.name
                        let currentNode = node

                        // Compact mode: merge folders with single subfolder
                        if (isCompactEnabled && node.isDir && node.children?.length === 1 && node.children[0].isDir) {
                            let nextNode = node.children[0]
                            while (nextNode.isDir && nextNode.children?.length === 1 && nextNode.children[0].isDir) {
                                displayName += '/' + nextNode.name
                                nextNode = nextNode.children[0]
                            }
                            displayName += '/' + nextNode.name
                            currentNode = nextNode
                        }

                        result.push({
                            id: node.path,
                            node: currentNode,
                            depth,
                            isLast,
                            ancestorHasMoreSiblings: [...ancestorHasMoreSiblings],
                            displayName
                        })

                        // Only recurse into expanded directories
                        if (currentNode.isDir && currentNode.isExpanded && currentNode.children?.length) {
                            flatten(
                                currentNode.children,
                                depth + 1,
                                [...ancestorHasMoreSiblings, !isLast]
                            )
                        }
                    })
                }

                flatten(focusedNodes, 0, [])
                return result
            }
        }

        // Selected Only Mode: show flat list of selected files only
        if (isSelectedOnlyMode?.value && selectedPaths?.value) {
            const allNodes: FileNode[] = []

            // Collect all nodes from tree
            function collectAllNodes(nodeList: FileNode[]) {
                nodeList.forEach((node) => {
                    allNodes.push(node)
                    if (node.children?.length) {
                        collectAllNodes(node.children)
                    }
                })
            }
            collectAllNodes(nodes.value)

            // Filter only selected files and create flat list
            allNodes
                .filter((node) => !node.isDir && selectedPaths.value.has(node.path))
                .forEach((node) => {
                    // Calculate relative path from root
                    let relativePath = node.path
                    if (rootPath?.value && node.path.startsWith(rootPath.value)) {
                        relativePath = node.path.slice(rootPath.value.length).replace(/^[/\\]/, '')
                    }

                    result.push({
                        id: node.path,
                        node,
                        depth: 0, // Flat list, no depth
                        isLast: false,
                        ancestorHasMoreSiblings: [],
                        displayName: node.name,
                        relativePath
                    })
                })

            return result
        }

        // Normal tree mode
        function flatten(
            nodeList: FileNode[],
            depth: number,
            ancestorHasMoreSiblings: boolean[]
        ) {
            nodeList.forEach((node, index) => {
                const isLast = index === nodeList.length - 1
                let displayName = node.name
                let currentNode = node

                // Compact mode: merge folders with single subfolder
                if (isCompactEnabled && node.isDir && node.children?.length === 1 && node.children[0].isDir) {
                    let nextNode = node.children[0]
                    // We also need to check if the node is expanded to continue compacting
                    // Usually compacting shows the chain even if some are not expanded, 
                    // but the LAST one in the chain is what we interact with.
                    while (nextNode.isDir && nextNode.children?.length === 1 && nextNode.children[0].isDir) {
                        displayName += '/' + nextNode.name
                        nextNode = nextNode.children[0]
                    }
                    displayName += '/' + nextNode.name
                    currentNode = nextNode
                }

                result.push({
                    id: node.path, // Use original node path as ID
                    node: currentNode,
                    depth,
                    isLast,
                    ancestorHasMoreSiblings: [...ancestorHasMoreSiblings],
                    displayName
                })

                // Only recurse into expanded directories
                // Note: in compact mode, we interact with the deepest folder in the chain
                if (currentNode.isDir && currentNode.isExpanded && currentNode.children?.length) {
                    flatten(
                        currentNode.children,
                        depth + 1,
                        [...ancestorHasMoreSiblings, !isLast]
                    )
                }
            })
        }

        flatten(nodes.value, 0, [])
        return result
    })


    const totalVisibleCount = computed(() => flattenedVisibleNodes.value.length)

    return {
        flattenedVisibleNodes,
        totalVisibleCount,
    }
}
