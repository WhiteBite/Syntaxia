/**
 * Testing API - Wails bridge for test operations
 */

import type { domain } from '#wailsjs/go/models'
import { useLogger } from '@/composables/useLogger'
import { buildApi } from '@/services/api/build.api'
import type { TestResult, TestSuiteUI } from '../types'

const logger = useLogger('TestingApi')

/**
 * Convert backend TestResult to UI TestResult
 */
function mapTestResult(result: domain.TestResult, index: number): TestResult {
    const fileName = result.testPath?.split(/[/\\]/).pop() || ''
    const suiteName = fileName.replace(/\.(test|spec)\.(ts|js|go|py)$/, '')

    return {
        id: `test-${index}-${Date.now()}`,
        name: result.testName || `Test ${index + 1}`,
        file: result.testPath || '',
        suite: suiteName,
        status: result.success ? 'passed' : 'failed',
        duration: result.duration,
        output: result.output,
        error: result.error,
        stackTrace: result.error
    }
}

/**
 * Group test results into suites by file
 */
function groupIntoSuites(results: TestResult[]): TestSuiteUI[] {
    const suiteMap = new Map<string, TestSuiteUI>()

    for (const test of results) {
        const suiteKey = test.file || 'Unknown'

        if (!suiteMap.has(suiteKey)) {
            const fileName = suiteKey.split(/[/\\]/).pop() || suiteKey
            suiteMap.set(suiteKey, {
                id: `suite-${suiteKey}`,
                name: fileName,
                file: suiteKey,
                tests: [],
                status: 'pending',
                expanded: true
            })
        }

        suiteMap.get(suiteKey)!.tests.push(test)
    }

    // Calculate suite status
    for (const suite of suiteMap.values()) {
        const hasRunning = suite.tests.some(t => t.status === 'running')
        const hasFailed = suite.tests.some(t => t.status === 'failed')
        const allPassed = suite.tests.every(t => t.status === 'passed')

        if (hasRunning) {
            suite.status = 'running'
        } else if (hasFailed) {
            suite.status = 'failed'
        } else if (allPassed) {
            suite.status = 'passed'
        } else {
            suite.status = 'pending'
        }
    }

    return Array.from(suiteMap.values())
}

export const testingApi = {
    /**
     * Run all tests in the project
     */
    async runTests(
        projectPath: string,
        testFilter?: string[]
    ): Promise<{ results: TestResult[]; suites: TestSuiteUI[] }> {
        logger.info(`Running tests for ${projectPath}`, { filter: testFilter })

        const config: domain.TestConfig = {
            projectPath,
            language: 'auto',
            scope: 'all',
            parallel: true,
            timeout: 60,
            coverage: false,
            verbose: true,
            testPatterns: testFilter
        }

        const backendResults = await buildApi.runTests(config)
        const results = backendResults.map(mapTestResult)
        const suites = groupIntoSuites(results)

        logger.info(`Tests completed: ${results.length} tests`)
        return { results, suites }
    },

    /**
     * Run a single test by name
     */
    async runSingleTest(
        projectPath: string,
        testName: string
    ): Promise<TestResult> {
        logger.info(`Running single test: ${testName}`)

        const config: domain.TestConfig = {
            projectPath,
            language: 'auto',
            scope: 'single',
            parallel: false,
            timeout: 30,
            coverage: false,
            verbose: true,
            testPatterns: [testName]
        }

        const backendResults = await buildApi.runTests(config)

        if (backendResults.length === 0) {
            return {
                id: `test-${Date.now()}`,
                name: testName,
                file: '',
                status: 'failed',
                error: 'Test not found or failed to run'
            }
        }

        return mapTestResult(backendResults[0], 0)
    },

    /**
     * Discover tests in the project
     */
    async discoverTests(
        projectPath: string
    ): Promise<TestSuiteUI[]> {
        logger.info(`Discovering tests for ${projectPath}`)

        const suite = await buildApi.discoverTests(projectPath, 'auto')

        if (!suite.tests || suite.tests.length === 0) {
            return []
        }

        const results: TestResult[] = suite.tests.map((test, index) => ({
            id: `test-${index}-${Date.now()}`,
            name: test.name,
            file: test.path,
            suite: test.type,
            status: 'pending' as const,
            duration: undefined,
            output: undefined,
            error: undefined
        }))

        return groupIntoSuites(results)
    },

    /**
     * Stop running tests (placeholder - actual implementation depends on backend)
     */
    async stopTests(): Promise<void> {
        logger.info('Stop tests requested')
        // Backend implementation needed for actual cancellation
    }
}
