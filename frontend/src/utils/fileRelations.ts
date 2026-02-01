/**
 * File Relations Utility
 * Finds related files based on naming patterns and location
 */

import type { FileNode } from '@/types/domain'

/**
 * Related file patterns to search for
 */
const RELATED_EXTENSIONS = [
    '.css',
    '.scss',
    '.sass',
    '.less',
    '.test.ts',
    '.test.js',
    '.spec.ts',
    '.spec.js',
    '.types.ts',
    '.d.ts',
    '.stories.ts',
    '.stories.js',
]

/**
 * Find related files for a given file path
 * @param filePath - The path of the file to find relations for
 * @param allNodes - All file nodes in the tree
 * @returns Array of related file paths
 */
export function findRelatedFiles(
    filePath: string,
    allNodes: FileNode[]
): string[] {
    const relatedPaths: string[] = []

    // Parse the file path
    const lastSlash = filePath.lastIndexOf('/')
    const directory = lastSlash >= 0 ? filePath.substring(0, lastSlash) : ''
    const fileName = lastSlash >= 0 ? filePath.substring(lastSlash + 1) : filePath

    // Get basename without extension
    const basename = getBasename(fileName)

    // Find files in the same directory with related patterns
    for (const node of allNodes) {
        if (node.isDir || node.path === filePath) continue

        // Check if in same directory
        const nodeLastSlash = node.path.lastIndexOf('/')
        const nodeDirectory = nodeLastSlash >= 0 ? node.path.substring(0, nodeLastSlash) : ''

        if (nodeDirectory !== directory) continue

        // Check if the file matches related patterns
        const nodeName = nodeLastSlash >= 0 ? node.path.substring(nodeLastSlash + 1) : node.path

        if (isRelatedFile(basename, nodeName)) {
            relatedPaths.push(node.path)
        }
    }

    return relatedPaths
}

/**
 * Get basename without extension
 * Examples:
 * - "Component.vue" -> "Component"
 * - "file.test.ts" -> "file"
 * - "styles.module.css" -> "styles"
 */
function getBasename(fileName: string): string {
    // Handle compound extensions like .test.ts, .spec.js, .d.ts
    const compoundPatterns = [
        '.test.ts',
        '.test.js',
        '.spec.ts',
        '.spec.js',
        '.types.ts',
        '.d.ts',
        '.module.css',
        '.module.scss',
        '.stories.ts',
        '.stories.js',
    ]

    for (const pattern of compoundPatterns) {
        if (fileName.endsWith(pattern)) {
            return fileName.substring(0, fileName.length - pattern.length)
        }
    }

    // Handle simple extensions
    const lastDot = fileName.lastIndexOf('.')
    return lastDot > 0 ? fileName.substring(0, lastDot) : fileName
}

/**
 * Check if a file is related to the base file
 */
function isRelatedFile(basename: string, candidateFileName: string): boolean {
    // Check if candidate starts with the same basename
    if (!candidateFileName.startsWith(basename)) {
        return false
    }

    // Get the part after basename
    const suffix = candidateFileName.substring(basename.length)

    // Check if suffix matches any related pattern
    for (const ext of RELATED_EXTENSIONS) {
        if (suffix === ext || suffix.startsWith(ext)) {
            return true
        }
    }

    // Check for module patterns like .module.css
    if (suffix.match(/^\.(module|component)\.(css|scss|sass|less)$/)) {
        return true
    }

    return false
}

/**
 * Get all file nodes from a tree structure
 */
export function getAllFileNodes(nodes: FileNode[]): FileNode[] {
    const result: FileNode[] = []

    function traverse(node: FileNode) {
        if (!node.isDir) {
            result.push(node)
        }
        if (node.children) {
            for (const child of node.children) {
                traverse(child)
            }
        }
    }

    for (const node of nodes) {
        traverse(node)
    }

    return result
}
