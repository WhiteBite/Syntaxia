/**
 * File Tree Utilities
 * Pure functions for tree traversal and manipulation
 */

import type { FileNode } from '@/types/domain'

/**
 * Domain node from backend API
 */
export interface DomainNode {
    name: string
    path: string
    isDir: boolean
    children?: DomainNode[]
    size?: number
    isIgnored?: boolean
    isGitignored?: boolean
    isCustomIgnored?: boolean
}

/**
 * Find a node in the tree by path
 * @param tree - Root nodes array
 * @param path - Path to find
 * @param cache - Optional cache map for O(1) lookups
 */
export function findNode(
    tree: FileNode[],
    path: string,
    cache?: Map<string, FileNode>
): FileNode | null {
    // Check cache first - O(1) lookup
    if (cache) {
        const cached = cache.get(path)
        if (cached) return cached
    }

    // Fallback to tree traversal
    for (const node of tree) {
        if (node.path === path) {
            cache?.set(path, node)
            return node
        }
        if (node.children) {
            const found = findNode(node.children, path, cache)
            if (found) return found
        }
    }
    return null
}

/**
 * Walk through all nodes in the tree
 * @param tree - Root nodes array
 * @param fn - Callback function for each node
 */
export function walkTree(tree: FileNode[], fn: (node: FileNode) => void): void {
    for (const node of tree) {
        fn(node)
        if (node.children) {
            walkTree(node.children, fn)
        }
    }
}

/**
 * Filter tree by file extensions
 * @param tree - Root nodes array
 * @param include - Extensions to include (e.g., ['.ts', '.vue'])
 * @param exclude - Extensions to exclude
 */
export function filterTreeByExtensions(
    tree: FileNode[],
    include: string[],
    exclude: string[] = []
): FileNode[] {
    return tree
        .filter((node) => {
            if (node.isDir) {
                const filteredChildren = node.children
                    ? filterTreeByExtensions(node.children, include, exclude)
                    : []
                return filteredChildren.length > 0
            }
            // If exclude list has items, exclude those extensions
            if (exclude.length > 0 && exclude.some((ext) => node.name.endsWith(ext))) {
                return false
            }
            // If include list is empty (only excluding), show all non-excluded
            if (include.length === 0) return true
            // Otherwise filter by include list
            return include.some((ext) => node.name.endsWith(ext))
        })
        .map((node) => ({
            ...node,
            children: node.children
                ? filterTreeByExtensions(node.children, include, exclude)
                : undefined,
        }))
}

/**
 * Convert domain nodes from backend to FileNode format
 * @param domainNodes - Nodes from backend API
 */
export function convertDomainNodes(domainNodes: DomainNode[]): FileNode[] {
    return domainNodes.map((node) => ({
        name: node.name,
        path: node.path,
        isDir: node.isDir,
        isExpanded: false,
        isSelected: false,
        children: node.children ? convertDomainNodes(node.children) : undefined,
        size: node.size,
        isIgnored: node.isIgnored || node.isGitignored || node.isCustomIgnored,
    }))
}

/**
 * Generate flattened list of nodes (for search/virtualization)
 * @param tree - Root nodes array
 * @param maxDepth - Maximum depth to traverse
 * @param currentDepth - Current depth (internal)
 */
export function generateFlattened(
    tree: FileNode[],
    maxDepth: number,
    currentDepth: number = 0
): FileNode[] {
    if (currentDepth >= maxDepth) return []

    const result: FileNode[] = []
    for (const node of tree) {
        result.push(node)
        if (node.children && currentDepth < maxDepth - 1) {
            result.push(...generateFlattened(node.children, maxDepth, currentDepth + 1))
        }
    }
    return result
}

/**
 * Build path cache for entire tree - O(1) node lookups
 * @param tree - Root nodes array
 * @param cache - Map to populate
 */
export function buildNodePathCache(
    tree: FileNode[],
    cache: Map<string, FileNode>
): void {
    for (const node of tree) {
        cache.set(node.path, node)
        if (node.children) {
            buildNodePathCache(node.children, cache)
        }
    }
}

/**
 * Get all file paths in a node (recursive)
 * @param node - Node to get files from
 * @param cache - Optional cache for results
 */
export function getAllFilesInNode(
    node: FileNode,
    cache?: Map<string, string[]>
): string[] {
    // Check cache first
    if (cache) {
        const cached = cache.get(node.path)
        if (cached) return cached
    }

    const files: string[] = []

    if (node.children) {
        for (const child of node.children) {
            if (!child.isDir) {
                files.push(child.path)
            } else {
                files.push(...getAllFilesInNode(child, cache))
            }
        }
    }

    // Cache the result
    cache?.set(node.path, files)
    return files
}

/**
 * Count total files in tree
 * @param nodes - Root nodes array
 */
export function countTotalFiles(nodes: FileNode[]): number {
    let count = 0
    walkTree(nodes, (node) => {
        if (!node.isDir) count++
    })
    return count
}

/**
 * Auto-expand folders to reveal files (smart expand)
 * @param tree - Root nodes array
 * @param maxDepth - Maximum depth to expand
 * @param currentDepth - Current depth (internal)
 * @param state - Current expansion state to limit total expanded nodes
 */
export function autoExpandToFiles(
    tree: FileNode[],
    maxDepth: number = 2,
    currentDepth: number = 0,
    state: { expandedCount: number } = { expandedCount: 0 }
): void {
    if (currentDepth >= maxDepth || state.expandedCount > 30) return

    for (const node of tree) {
        if (node.isDir && node.children && node.children.length > 0) {
            // Root level folders always expanded if they have children
            if (currentDepth === 0) {
                node.isExpanded = true
                state.expandedCount++
                autoExpandToFiles(node.children, maxDepth, currentDepth + 1, state)
                continue
            }

            // For deeper levels, expand only if:
            // 1. It has direct files and not too many children (avoid cluttering)
            // 2. OR it has only one child which is also a directory (chain)
            const hasFiles = node.children.some((child) => !child.isDir)
            const childCount = node.children.length

            if ((hasFiles && childCount < 15 && state.expandedCount < 25) ||
                (childCount === 1 && node.children[0].isDir)) {
                node.isExpanded = true
                state.expandedCount++
                autoExpandToFiles(node.children, maxDepth, currentDepth + 1, state)
            }
        }
    }
}

