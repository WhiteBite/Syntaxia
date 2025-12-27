/**
 * Build feature types
 */

export type BuildStatusType = 'idle' | 'building' | 'success' | 'failed' | 'cancelled'

export interface BuildError {
    file: string
    line: number
    column?: number
    message: string
}

export interface BuildWarning {
    file: string
    line: number
    column?: number
    message: string
}

export interface BuildResult {
    id: string
    status: 'success' | 'failed' | 'cancelled'
    startTime: string
    endTime: string
    duration: number
    output: string
    errors: BuildError[]
    warnings: BuildWarning[]
}

export interface BuildStatus {
    status: BuildStatusType
    progress?: number
    currentStep?: string
}

export interface BuildHistoryEntry {
    id: string
    status: 'success' | 'failed' | 'cancelled'
    startTime: string
    endTime: string
    duration: number
    errorsCount: number
    warningsCount: number
    output: string
}
