/**
 * Testing feature types
 */

export type TestStatus = 'passed' | 'failed' | 'skipped' | 'running' | 'pending'

export type TestFilter = 'all' | 'passed' | 'failed' | 'skipped'

export interface TestResult {
    id: string
    name: string
    file: string
    suite?: string
    status: TestStatus
    duration?: number
    output?: string
    error?: string
    stackTrace?: string
}

export interface TestSuiteUI {
    id: string
    name: string
    file: string
    tests: TestResult[]
    status: TestStatus
    expanded: boolean
}

export interface TestRunStats {
    total: number
    passed: number
    failed: number
    skipped: number
    running: number
    duration: number
}
