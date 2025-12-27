/**
 * Testing Store - Pinia store for test runner state management
 */

import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { testingApi } from '../api/testing.api'
import type { TestFilter, TestResult, TestRunStats, TestSuiteUI } from '../types'

const logger = useLogger('TestingStore')

export const useTestingStore = defineStore('testing', () => {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const { t } = useI18n()

    // State
    const tests = ref<TestResult[]>([])
    const suites = ref<TestSuiteUI[]>([])
    const running = ref(false)
    const filter = ref<TestFilter>('all')
    const selectedTestId = ref<string | null>(null)
    const lastRunTimestamp = ref<string | null>(null)

    // Getters
    const stats = computed<TestRunStats>(() => {
        const all = tests.value
        return {
            total: all.length,
            passed: all.filter(t => t.status === 'passed').length,
            failed: all.filter(t => t.status === 'failed').length,
            skipped: all.filter(t => t.status === 'skipped').length,
            running: all.filter(t => t.status === 'running').length,
            duration: all.reduce((sum, t) => sum + (t.duration || 0), 0)
        }
    })

    const passedCount = computed(() => stats.value.passed)
    const failedCount = computed(() => stats.value.failed)
    const skippedCount = computed(() => stats.value.skipped)

    const filteredTests = computed(() => {
        if (filter.value === 'all') return tests.value
        return tests.value.filter(t => t.status === filter.value)
    })

    const filteredSuites = computed(() => {
        if (filter.value === 'all') return suites.value

        return suites.value
            .map(suite => ({
                ...suite,
                tests: suite.tests.filter(t => t.status === filter.value)
            }))
            .filter(suite => suite.tests.length > 0)
    })

    const selectedTest = computed(() => {
        if (!selectedTestId.value) return null
        return tests.value.find(t => t.id === selectedTestId.value) || null
    })

    const hasTests = computed(() => tests.value.length > 0)
    const hasFailedTests = computed(() => stats.value.failed > 0)

    const progressPercent = computed(() => {
        if (stats.value.total === 0) return 0
        const completed = stats.value.passed + stats.value.failed + stats.value.skipped
        return Math.round((completed / stats.value.total) * 100)
    })

    // Actions
    async function runAllTests(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('testing.selectProject'), 'warning')
            return
        }

        running.value = true
        logger.info('Running all tests')

        try {
            // Mark all tests as running
            tests.value = tests.value.map(t => ({ ...t, status: 'running' as const }))
            suites.value = suites.value.map(s => ({ ...s, status: 'running' as const }))

            const result = await testingApi.runTests(projectPath)
            tests.value = result.results
            suites.value = result.suites
            lastRunTimestamp.value = new Date().toISOString()

            const { passed, failed } = stats.value
            if (failed > 0) {
                uiStore.addToast(t('testing.testsFailed', { passed, failed }), 'error')
            } else {
                uiStore.addToast(t('testing.testsPassed', { count: passed }), 'success')
            }

            logger.info(`Tests completed: ${passed} passed, ${failed} failed`)
        } catch (error) {
            logger.error('Failed to run tests:', error)
            uiStore.addToast(t('testing.runFailed'), 'error')
        } finally {
            running.value = false
        }
    }

    async function runFailedTests(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('testing.selectProject'), 'warning')
            return
        }

        const failedTestNames = tests.value
            .filter(t => t.status === 'failed')
            .map(t => t.name)

        if (failedTestNames.length === 0) {
            uiStore.addToast(t('testing.noFailedTests'), 'info')
            return
        }

        running.value = true
        logger.info(`Running ${failedTestNames.length} failed tests`)

        try {
            // Mark failed tests as running
            tests.value = tests.value.map(t =>
                t.status === 'failed' ? { ...t, status: 'running' as const } : t
            )

            const result = await testingApi.runTests(projectPath, failedTestNames)

            // Merge results
            for (const newResult of result.results) {
                const index = tests.value.findIndex(t => t.name === newResult.name)
                if (index !== -1) {
                    tests.value[index] = newResult
                }
            }

            // Recalculate suites
            updateSuiteStatuses()
            lastRunTimestamp.value = new Date().toISOString()

            logger.info('Failed tests re-run completed')
        } catch (error) {
            logger.error('Failed to run failed tests:', error)
            uiStore.addToast(t('testing.runFailed'), 'error')
        } finally {
            running.value = false
        }
    }

    async function runTest(testId: string): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('testing.selectProject'), 'warning')
            return
        }

        const test = tests.value.find(t => t.id === testId)
        if (!test) return

        running.value = true
        logger.info(`Running single test: ${test.name}`)

        try {
            // Mark test as running
            const index = tests.value.findIndex(t => t.id === testId)
            if (index !== -1) {
                tests.value[index] = { ...tests.value[index], status: 'running' }
            }

            const result = await testingApi.runSingleTest(projectPath, test.name)

            // Update test result
            if (index !== -1) {
                tests.value[index] = { ...result, id: testId }
            }

            updateSuiteStatuses()
        } catch (error) {
            logger.error(`Failed to run test ${test.name}:`, error)
            uiStore.addToast(t('testing.runFailed'), 'error')
        } finally {
            running.value = false
        }
    }

    async function stopTests(): Promise<void> {
        logger.info('Stopping tests')
        await testingApi.stopTests()
        running.value = false

        // Mark running tests as pending
        tests.value = tests.value.map(t =>
            t.status === 'running' ? { ...t, status: 'pending' as const } : t
        )
        updateSuiteStatuses()
    }

    async function discoverTests(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return

        logger.info('Discovering tests')

        try {
            const discovered = await testingApi.discoverTests(projectPath)
            suites.value = discovered
            tests.value = discovered.flatMap(s => s.tests)
        } catch (error) {
            logger.error('Failed to discover tests:', error)
        }
    }

    function setFilter(newFilter: TestFilter): void {
        filter.value = newFilter
    }

    function selectTest(testId: string | null): void {
        selectedTestId.value = testId
    }

    function toggleSuite(suiteId: string): void {
        const suite = suites.value.find(s => s.id === suiteId)
        if (suite) {
            suite.expanded = !suite.expanded
        }
    }

    function expandAllSuites(): void {
        suites.value.forEach(s => { s.expanded = true })
    }

    function collapseAllSuites(): void {
        suites.value.forEach(s => { s.expanded = false })
    }

    function clearResults(): void {
        tests.value = []
        suites.value = []
        selectedTestId.value = null
        lastRunTimestamp.value = null
    }

    function updateSuiteStatuses(): void {
        for (const suite of suites.value) {
            const suiteTests = tests.value.filter(t => t.file === suite.file)
            suite.tests = suiteTests

            const hasRunning = suiteTests.some(t => t.status === 'running')
            const hasFailed = suiteTests.some(t => t.status === 'failed')
            const allPassed = suiteTests.every(t => t.status === 'passed')

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
    }

    return {
        // State
        tests,
        suites,
        running,
        filter,
        selectedTestId,
        lastRunTimestamp,

        // Getters
        stats,
        passedCount,
        failedCount,
        skippedCount,
        filteredTests,
        filteredSuites,
        selectedTest,
        hasTests,
        hasFailedTests,
        progressPercent,

        // Actions
        runAllTests,
        runFailedTests,
        runTest,
        stopTests,
        discoverTests,
        setFilter,
        selectTest,
        toggleSuite,
        expandAllSuites,
        collapseAllSuites,
        clearResults
    }
})
