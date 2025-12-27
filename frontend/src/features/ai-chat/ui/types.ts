// Types for Context Used components

export type ContextOperation = 'read' | 'write' | 'search' | 'created' | 'deleted'

export interface ContextUsedFile {
    path: string
    operation: ContextOperation
    tokens?: number
}

export interface ContextUsed {
    files: ContextUsedFile[]
    totalTokens: number
}
