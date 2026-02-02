/**
 * Task feature types
 */

export type TaskType = 'feature' | 'bugfix' | 'refactor' | 'test' | 'docs'

export type TaskStatus = 'draft' | 'analyzing' | 'ready' | 'executing' | 'completed' | 'failed'

export type TaskComplexity = 'low' | 'medium' | 'high'

export interface Task {
    id: string
    description: string
    type: TaskType
    status: TaskStatus
    createdAt: string
    updatedAt: string
    suggestedFiles?: string[]
    selectedFiles?: string[]
    complexity?: TaskComplexity
    estimatedTime?: number
    recommendations?: string[]
}

export interface TaskDraft {
    description: string
    type?: TaskType
    timestamp: string
}

export interface AnalysisResult {
    suggestedFiles: string[]
    complexity: TaskComplexity
    estimatedTime?: number
    recommendations: string[]
}

export interface TaskExecutionResult {
    success: boolean
    filesModified: string[]
    error?: string
    duration: number
}

export interface TaskSuggestion {
    id: string
    text: string
    category: TaskType
}
