import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import {
    verificationApi,
    type Severity,
    type VerificationError,
    type VerificationResult,
    type VerificationType
} from '../api/verification.api'

const logger = useLogger('VerificationStore')

export type VerificationStatusType = 'idle' | 'running' | 'passed' | 'failed'

export interface VerificationStatus {
    build: VerificationStatusType
    lint: VerificationStatusType
    test: VerificationStatusType
}

export const useVerificationStore = defineStore('verification', () => {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const { t } = useI18n()

    // State
    const status = ref<VerificationStatus>({
        build: 'idle',
        lint: 'idle',
        test: 'idle'
    })

    const errors = ref<VerificationError[]>([])
    const results = ref<Map<VerificationType, VerificationResult>>(new Map())
    const activeTab = ref<VerificationType>('build')
    const isFixing = ref<string | null>(null) // Error ID being fixed
    const lastRunTimestamp = ref<string | null>(null)

    // Getters
    const hasErrors = computed(() =>
        errors.value.some(e => e.severity === 'error')
    )

    const hasWarnings = computed(() =>
        errors.value.some(e => e.severity === 'warning')
    )

    const errorCount = computed(() =>
        errors.value.filter(e => e.severity === 'error').length
    )

    const warningCount = computed(() =>
        errors.value.filter(e => e.severity === 'warning').length
    )

    const infoCount = computed(() =>
        errors.value.filter(e => e.severity === 'info').length
    )

    const isRunning = computed(() =>
        status.value.build === 'running' ||
        status.value.lint === 'running' ||
        status.value.test === 'running'
    )

    const allPassed = computed(() =>
        status.value.build === 'passed' &&
        status.value.lint === 'passed' &&
        status.value.test === 'passed'
    )

    const currentStatusText = computed(() => {
        if (status.value.build === 'running') return t('verification.building')
        if (status.value.lint === 'running') return t('verification.linting')
        if (status.value.test === 'running') return t('verification.testing')
        return ''
    })

    // Filter errors by source type
    const errorsBySource = computed(() => {
        return (source: VerificationType) =>
            errors.value.filter(e => e.source === source)
    })

    // Group errors by file
    const errorsByFile = computed(() => {
        const grouped = new Map<string, VerificationError[]>()

        for (const error of errors.value) {
            const existing = grouped.get(error.file) || []
            existing.push(error)
            grouped.set(error.file, existing)
        }

        return grouped
    })

    // Group errors by severity
    const errorsBySeverity = computed(() => {
        const grouped: Record<Severity, VerificationError[]> = {
            error: [],
            warning: [],
            info: []
        }

        for (const error of errors.value) {
            grouped[error.severity].push(error)
        }

        return grouped
    })

    // Actions
    async function runBuild(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('verification.selectProject'), 'warning')
            return
        }

        status.value.build = 'running'

        try {
            const result = await verificationApi.runVerification('build', projectPath)
            results.value.set('build', result)

            // Remove old build errors and add new ones
            errors.value = errors.value.filter(e => e.source !== 'build')
            errors.value.push(...result.errors)

            status.value.build = result.success ? 'passed' : 'failed'
            lastRunTimestamp.value = result.timestamp

            logger.info(`Build completed: ${result.success ? 'passed' : 'failed'}`)
        } catch (error) {
            logger.error('Build failed:', error)
            status.value.build = 'failed'
            uiStore.addToast(t('verification.failed'), 'error')
        }
    }

    async function runLint(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('verification.selectProject'), 'warning')
            return
        }

        status.value.lint = 'running'

        try {
            const result = await verificationApi.runVerification('lint', projectPath)
            results.value.set('lint', result)

            // Remove old lint errors and add new ones
            errors.value = errors.value.filter(e => e.source !== 'lint')
            errors.value.push(...result.errors)

            status.value.lint = result.success ? 'passed' : 'failed'
            lastRunTimestamp.value = result.timestamp

            logger.info(`Lint completed: ${result.success ? 'passed' : 'failed'}`)
        } catch (error) {
            logger.error('Lint failed:', error)
            status.value.lint = 'failed'
            uiStore.addToast(t('verification.failed'), 'error')
        }
    }

    async function runTests(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('verification.selectProject'), 'warning')
            return
        }

        status.value.test = 'running'

        try {
            const result = await verificationApi.runVerification('test', projectPath)
            results.value.set('test', result)

            // Remove old test errors and add new ones
            errors.value = errors.value.filter(e => e.source !== 'test')
            errors.value.push(...result.errors)

            status.value.test = result.success ? 'passed' : 'failed'
            lastRunTimestamp.value = result.timestamp

            logger.info(`Tests completed: ${result.success ? 'passed' : 'failed'}`)
        } catch (error) {
            logger.error('Tests failed:', error)
            status.value.test = 'failed'
            uiStore.addToast(t('verification.failed'), 'error')
        }
    }

    async function runAll(): Promise<void> {
        await runBuild()
        await runLint()
        await runTests()
    }

    async function fixWithAI(errorId: string): Promise<void> {
        const error = errors.value.find(e => e.id === errorId)
        if (!error) return

        const projectPath = projectStore.currentPath
        if (!projectPath) return

        isFixing.value = errorId

        try {
            const suggestion = await verificationApi.requestAIFix(error, projectPath)

            // Emit event for AI chat to handle
            window.dispatchEvent(new CustomEvent('verification:fix-request', {
                detail: { error, suggestion }
            }))

            logger.info(`AI fix requested for ${error.file}:${error.line}`)
        } catch (err) {
            logger.error('AI fix request failed:', err)
            uiStore.addToast(t('verification.failed'), 'error')
        } finally {
            isFixing.value = null
        }
    }

    function setActiveTab(tab: VerificationType): void {
        activeTab.value = tab
    }

    function clearResults(): void {
        errors.value = []
        results.value.clear()
        status.value = {
            build: 'idle',
            lint: 'idle',
            test: 'idle'
        }
        lastRunTimestamp.value = null
    }

    function removeError(errorId: string): void {
        errors.value = errors.value.filter(e => e.id !== errorId)
    }

    // Load cached results on init
    async function loadCachedResults(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return

        try {
            const cachedResults = await verificationApi.getVerificationResults(projectPath)

            for (const result of cachedResults) {
                results.value.set(result.type, result)
                errors.value.push(...result.errors)
                status.value[result.type] = result.success ? 'passed' : 'failed'
            }

            if (cachedResults.length > 0) {
                lastRunTimestamp.value = cachedResults[0].timestamp
            }
        } catch (error) {
            logger.warn('Failed to load cached results:', error)
        }
    }

    return {
        // State
        status,
        errors,
        results,
        activeTab,
        isFixing,
        lastRunTimestamp,

        // Getters
        hasErrors,
        hasWarnings,
        errorCount,
        warningCount,
        infoCount,
        isRunning,
        allPassed,
        currentStatusText,
        errorsBySource,
        errorsByFile,
        errorsBySeverity,

        // Actions
        runBuild,
        runLint,
        runTests,
        runAll,
        fixWithAI,
        setActiveTab,
        clearResults,
        removeError,
        loadCachedResults
    }
})
