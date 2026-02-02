/**
 * Test runner composable - Reusable test execution logic
 */

import { useLogger } from '@/composables/useLogger'
import { computed, ref } from 'vue'
import { testingApi } from '../api/testing.api'
import type { TestResult, TestSuiteUI } from '../types'

const logger = useLogger('useTestRunner')

export function useTestRunner() {
    const isRunning = ref(false)
    const currentTest = ref<string | null>(null)
    const results = ref<TestResult[]>([])
    const error = ref<string | null>(null)

    const hasResults = computed(() => results.value.length > 0)
    const passedCount = computed(() => results.value.filter(t => t.status === 'passed').length)
    const failedCount = computed(() => results.value.filter(t => t.status === 'failed').length)
    const totalDuration = computed(() => results.value.reduce((sum, t) => sum + (t.duration || 0), 0))

    /**
     * Run tests with optional filter
     */
    async function runTests(
        projectPath: string,
        testFilter?: string[]
    ): Promise<{ results: TestResult[]; suites: TestSuiteUI[] }> {
        if (isRunning.value) {
            logger.warn('Tests already running')
            throw new Error('Tests already running')
        }

        isRunning.value = true
        error.value = null

        try {
            logger.info('Running tests', { projectPath, filter: testFilter })

            const result = await testingApi.runTests(projectPath, testFilter)
            results.value = result.results

            logger.info(`Tests completed: ${passedCount.value} passed, ${failedCount.value} failed`)

            return result
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to run tests'
            logger.error('Test execution failed:', err)
            throw err
        } finally {
            isRunning.value = false
            currentTest.value = null
        }
    }

    /**
     * Run a single test
     */
    async function runSingleTest(
        projectPath: string,
        testName: string
    ): Promise<TestResult> {
        if (isRunning.value) {
            logger.warn('Tests already running')
            throw new Error('Tests already running')
        }

        isRunning.value = true
        currentTest.value = testName
        error.value = null

        try {
            logger.info('Running single test', { testName })

            const result = await testingApi.runSingleTest(projectPath, testName)

            // Update results array
            const index = results.value.findIndex(t => t.name === testName)
            if (index !== -1) {
                results.value[index] = result
            } else {
                results.value.push(result)
            }

            return result
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to run test'
            logger.error('Single test execution failed:', err)
            throw err
        } finally {
            isRunning.value = false
            currentTest.value = null
        }
    }

    /**
     * Discover available tests
     */
    async function discoverTests(projectPath: string): Promise<TestSuiteUI[]> {
        logger.info('Discovering tests', { projectPath })

        try {
            const suites = await testingApi.discoverTests(projectPath)

            // Extract all tests from suites
            results.value = suites.flatMap(s => s.tests)

            return suites
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to discover tests'
            logger.error('Test discovery failed:', err)
            throw err
        }
    }

    /**
     * Stop running tests
     */
    async function stopTests(): Promise<void> {
        if (!isRunning.value) return

        logger.info('Stopping tests')

        try {
            await testingApi.stopTests()
        } catch (err) {
            logger.error('Failed to stop tests:', err)
        } finally {
            isRunning.value = false
            currentTest.value = null
        }
    }

    /**
     * Clear test results
     */
    function clearResults(): void {
        results.value = []
        error.value = null
    }

    /**
     * Get test result by name
     */
    function getTestResult(testName: string): TestResult | undefined {
        return results.value.find(t => t.name === testName)
    }

    return {
        // State
        isRunning,
        currentTest,
        results,
        error,
        // Computed
        hasResults,
        passedCount,
        failedCount,
        totalDuration,
        // Actions
        runTests,
        runSingleTest,
        discoverTests,
        stopTests,
        clearResults,
        getTestResult
    }
}
