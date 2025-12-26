import type { domain } from '#wailsjs/go/models'
import type { ContextSummary } from '../model/context.store'

export interface ValidationResult {
    valid: boolean
    errors: string[]
}

export interface ContextStats {
    avgFileSize: number
    avgLinesPerFile: number
    estimatedReadTime: number // minutes
    complexity: 'low' | 'medium' | 'high'
}

export function estimateTokens(text: string): number {
    // Rough estimation: ~4 characters per token
    return Math.ceil(text.length / 4)
}

export function formatContextSize(bytes: number): string {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + ' ' + sizes[i]
}

export function validateContextOptions(options: Partial<domain.ContextBuildOptions>): ValidationResult {
    const errors: string[] = []

    if (options.maxTokens && options.maxTokens < 1000) {
        errors.push('Max tokens must be at least 1000')
    }

    if (options.maxTokens && options.maxTokens > 1000000) {
        errors.push('Max tokens cannot exceed 1,000,000')
    }

    if (options.splitStrategy && !['smart', 'file', 'token'].includes(options.splitStrategy)) {
        errors.push('Invalid split strategy')
    }

    return {
        valid: errors.length === 0,
        errors
    }
}

export function calculateContextStats(summary: ContextSummary): ContextStats {
    const avgFileSize = summary.fileCount > 0 ? summary.totalSize / summary.fileCount : 0
    const avgLinesPerFile = summary.fileCount > 0 ? summary.lineCount / summary.fileCount : 0
    const estimatedReadTime = Math.ceil(summary.lineCount / 50) // ~50 lines per minute

    let complexity: 'low' | 'medium' | 'high' = 'low'
    if (summary.fileCount > 20 || summary.lineCount > 5000) {
        complexity = 'high'
    } else if (summary.fileCount > 10 || summary.lineCount > 2000) {
        complexity = 'medium'
    }

    return {
        avgFileSize,
        avgLinesPerFile,
        estimatedReadTime,
        complexity
    }
}

export function splitContextIntoChunks(content: string, maxTokens: number): string[] {
    const lines = content.split('\n')
    const chunks: string[] = []
    let currentChunk: string[] = []
    let currentTokens = 0

    for (const line of lines) {
        const lineTokens = estimateTokens(line)

        if (currentTokens + lineTokens > maxTokens && currentChunk.length > 0) {
            chunks.push(currentChunk.join('\n'))
            currentChunk = []
            currentTokens = 0
        }

        currentChunk.push(line)
        currentTokens += lineTokens
    }

    if (currentChunk.length > 0) {
        chunks.push(currentChunk.join('\n'))
    }

    return chunks
}

export function formatTimestamp(isoString: string): string {
    const date = new Date(isoString)
    const now = new Date()
    const diffMs = now.getTime() - date.getTime()
    const diffMins = Math.floor(diffMs / 60000)

    if (diffMins < 1) return 'just now'
    if (diffMins < 60) return `${diffMins}m ago`

    const diffHours = Math.floor(diffMins / 60)
    if (diffHours < 24) return `${diffHours}h ago`

    const diffDays = Math.floor(diffHours / 24)
    if (diffDays < 7) return `${diffDays}d ago`

    return date.toLocaleDateString()
}
