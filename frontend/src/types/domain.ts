/**
 * Domain Types - Shared interfaces for the application
 */

/**
 * Content type for files
 */
export type FileContentType = 'text' | 'binary' | 'unknown'

/**
 * File node for UI tree representation
 */
export interface FileNode {
    name: string
    path: string
    isDir: boolean
    isExpanded?: boolean
    isSelected?: boolean
    children?: FileNode[]
    size?: number
    isIgnored?: boolean
    contentType?: FileContentType

    // Pre-computed metadata from backend (for performance optimization)
    fileCount?: number              // Total files in directory (recursive)
    totalSize?: number              // Total size in bytes (recursive)
    tokenCount?: number             // Estimated token count for AI context
    depth?: number                  // Depth from root
    directFileCount?: number        // Files directly in this folder (non-recursive)
    extensionStats?: Record<string, number> // Extension counts (e.g., {".ts": 45, ".vue": 12})
}

/**
 * File node from backend API (raw format)
 */
export interface DomainFileNode {
    name: string
    path: string
    isDir: boolean
    children?: DomainFileNode[]
    size?: number
    isIgnored?: boolean
    isGitignored?: boolean
    isCustomIgnored?: boolean
    contentType?: FileContentType

    // Pre-computed metadata from backend (for performance optimization)
    fileCount?: number              // Total files in directory (recursive)
    totalSize?: number              // Total size in bytes (recursive)
    tokenCount?: number             // Estimated token count for AI context
    depth?: number                  // Depth from root
    directFileCount?: number        // Files directly in this folder (non-recursive)
    extensionStats?: Record<string, number> // Extension counts
}

/**
 * Context summary from backend
 */
export interface ContextSummary {
    id: string
    name: string
    files: string[]
    totalSize: number
    totalTokens: number
    createdAt: string
    updatedAt?: string
}

/**
 * Export settings for context export
 */
export interface ExportSettings {
    format: 'xml' | 'markdown' | 'json' | 'plain'
    includeLineNumbers: boolean
    includeMetadata: boolean
    maxTokens?: number
    chunkSize?: number
}

/**
 * Export result from backend
 */
export interface ExportResult {
    content: string
    format: string
    totalFiles: number
    totalSize: number
    totalTokens: number
    chunks?: string[]
}

/**
 * Command definition for command palette
 */
export interface Command {
    id: string
    name: string
    description: string
    shortcut?: string
    icon: CommandIcon
    action: () => void | Promise<void>
    category?: string
}

/**
 * Command icon - can be a component or string
 */
export type CommandIcon = string | { render: () => unknown }

/**
 * Cache entry with TTL
 */
export interface CacheEntry<T> {
    value: T
    timestamp: number
    ttl: number
}

/**
 * Memory diagnostic entry
 */
export interface DiagnosticEntry {
    category: 'memory' | 'performance' | 'store' | 'cache'
    message: string
    details: Record<string, unknown>
    timestamp: string
}

/**
 * Memory error context
 */
export interface MemoryErrorContext {
    component?: string
    operation?: string
    memoryUsage?: number
    [key: string]: unknown
}

/**
 * Batch request handler
 */
export interface BatchRequest<T, D = unknown> {
    resolve: (value: T) => void
    reject: (error: Error) => void
    data: D
}

/**
 * Generic function type for debounce/throttle
 */
export type AnyFunction = (...args: unknown[]) => unknown

/**
 * File dependency information
 */
export interface FileDependency {
    sourcePath: string
    targetPath: string
    type: 'import' | 'style' | 'test' | 'type'
    importName?: string
    line?: number
}

/**
 * Dependency graph for the project
 */
export interface DependencyGraph {
    files: Record<string, FileDependency[]>
    lastAnalyzed: string // ISO date
}
