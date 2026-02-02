/**
 * Verification composable - Reusable verification logic
 */

import { useLogger } from '@/composables/useLogger'
import { computed, ref } from 'vue'
import { verificationApi, type VerificationError, type VerificationResult, type VerificationType } from '../api/verification.api'

const logger = useLogger('useVerification')

export function useVerification() {
    const isRunning = ref(false)
    const currentType = ref<VerificationType | null>(null)
    const results = ref<Map<VerificationType, VerificationResult>>(new Map())
    const errors = ref<VerificationError[]>([])

    const hasErrors = computed(() => errors.value.some(e => e.severity === 'error'))
    const hasWarnings = computed(() => errors.value.some(e => e.severity === 'warning'))
    const errorCount = computed(() => errors.value.filter(e => e.severity === 'error').length)
    const warningCount = computed(() => errors.value.filter(e => e.severity === 'warning').length)

    /**
     * Run verification of specified type
     */
    async function runVerification(
        type: VerificationType,
        projectPath: string
    ): Promise<VerificationResult> {
        if (isRunning.value) {
            logger.warn('Verification already running')
            throw new Error('Verification already running')
        }

        isRunning.value = true
        currentType.value = type

        try {
            logger.info(`Running ${type} verification`)

            const result = await verificationApi.runVerification(type, projectPath)
            results.value.set(type, result)

            // Update errors - remove old errors of this type and add new ones
            errors.value = errors.value.filter(e => e.source !== type)
            errors.value.push(...result.errors)

            logger.info(`${type} verification completed: ${result.success ? 'passed' : 'failed'}`)

            return result
        } catch (err) {
            logger.error(`${type} verification failed:`, err)
            throw err
        } finally {
            isRunning.value = false
            currentType.value = null
        }
    }

    /**
     * Run all verifications sequentially
     */
    async function runAllVerifications(projectPath: string): Promise<void> {
        const types: VerificationType[] = ['build', 'lint', 'test']

        for (const type of types) {
            try {
                await runVerification(type, projectPath)
            } catch (err) {
                logger.error(`Failed to run ${type} verification:`, err)
                // Continue with other verifications even if one fails
            }
        }
    }

    /**
     * Get verification result for a specific type
     */
    function getResult(type: VerificationType): VerificationResult | undefined {
        return results.value.get(type)
    }

    /**
     * Get errors for a specific type
     */
    function getErrorsByType(type: VerificationType): VerificationError[] {
        return errors.value.filter(e => e.source === type)
    }

    /**
     * Get errors for a specific file
     */
    function getErrorsByFile(file: string): VerificationError[] {
        return errors.value.filter(e => e.file === file)
    }

    /**
     * Request AI fix for an error
     */
    async function requestAIFix(
        error: VerificationError,
        projectPath: string
    ): Promise<string> {
        logger.info(`Requesting AI fix for ${error.file}:${error.line}`)

        try {
            const prompt = await verificationApi.requestAIFix(error, projectPath)
            return prompt
        } catch (err) {
            logger.error('AI fix request failed:', err)
            throw err
        }
    }

    /**
     * Clear all results and errors
     */
    function clearResults(): void {
        results.value.clear()
        errors.value = []
    }

    /**
     * Remove a specific error
     */
    function removeError(errorId: string): void {
        errors.value = errors.value.filter(e => e.id !== errorId)
    }

    return {
        // State
        isRunning,
        currentType,
        results,
        errors,
        // Computed
        hasErrors,
        hasWarnings,
        errorCount,
        warningCount,
        // Actions
        runVerification,
        runAllVerifications,
        getResult,
        getErrorsByType,
        getErrorsByFile,
        requestAIFix,
        clearResults,
        removeError
    }
}
