import { useLogger } from '@/composables/useLogger'
import { buildApi } from '@/services/api/build.api'

const logger = useLogger('VerificationApi')

export type VerificationType = 'build' | 'lint' | 'test'
export type Severity = 'error' | 'warning' | 'info'

export interface VerificationError {
    id: string
    file: string
    line: number
    column?: number
    message: string
    severity: Severity
    source: VerificationType
    rule?: string
}

export interface VerificationResult {
    type: VerificationType
    success: boolean
    errors: VerificationError[]
    duration: number
    filesChecked: number
    timestamp: string
}

export class VerificationApi {
    /**
     * Run verification of specified type
     */
    async runVerification(
        type: VerificationType,
        projectPath: string
    ): Promise<VerificationResult> {
        const startTime = Date.now()

        try {
            logger.info(`Running ${type} verification for ${projectPath}`)

            const errors: VerificationError[] = []
            let success = true
            let filesChecked = 0

            // Use existing build API methods
            if (type === 'build') {
                const result = await buildApi.build(projectPath, 'auto')
                success = result.success
                filesChecked = 1

                // BuildResult has error (string) and warnings (string[])
                if (result.error) {
                    errors.push({
                        id: `build-error-${Date.now()}`,
                        file: projectPath,
                        line: 0,
                        message: result.error,
                        severity: 'error',
                        source: 'build'
                    })
                }

                if (result.warnings) {
                    result.warnings.forEach((warning, index) => {
                        errors.push({
                            id: `build-warning-${index}-${Date.now()}`,
                            file: projectPath,
                            line: 0,
                            message: warning,
                            severity: 'warning',
                            source: 'build'
                        })
                    })
                }
            } else if (type === 'lint') {
                const result = await buildApi.typeCheck(projectPath, 'auto')
                success = result.success

                // TypeCheckResult has issues (TypeIssue[])
                if (result.issues) {
                    filesChecked = new Set(result.issues.map(i => i.file)).size

                    result.issues.forEach((issue, index) => {
                        errors.push({
                            id: `lint-${index}-${Date.now()}`,
                            file: issue.file,
                            line: issue.line,
                            column: issue.column,
                            message: issue.message,
                            severity: (issue.severity as Severity) || 'warning',
                            source: 'lint',
                            rule: issue.code
                        })
                    })
                }

                // Also check for general error
                if (result.error) {
                    errors.push({
                        id: `lint-error-${Date.now()}`,
                        file: projectPath,
                        line: 0,
                        message: result.error,
                        severity: 'error',
                        source: 'lint'
                    })
                }
            } else if (type === 'test') {
                const results = await buildApi.runTests({
                    projectPath,
                    language: 'auto',
                    scope: 'all',
                    parallel: true,
                    timeout: 60,
                    coverage: false,
                    verbose: false
                })

                const failedTests = results.filter(r => !r.success)
                success = failedTests.length === 0
                filesChecked = results.length

                failedTests.forEach((r, index) => {
                    errors.push({
                        id: `test-${index}-${Date.now()}`,
                        file: r.testPath || r.testName,
                        line: 1,
                        message: r.error || `Test failed: ${r.testName}`,
                        severity: 'error',
                        source: 'test'
                    })
                })
            }

            const duration = Date.now() - startTime

            return {
                type,
                success,
                errors,
                duration,
                filesChecked,
                timestamp: new Date().toISOString()
            }
        } catch (error) {
            logger.error(`Failed to run ${type} verification:`, error)

            // Return empty result on error
            return {
                type,
                success: false,
                errors: [{
                    id: `${type}-error-${Date.now()}`,
                    file: '',
                    line: 0,
                    message: error instanceof Error ? error.message : 'Verification failed',
                    severity: 'error',
                    source: type
                }],
                duration: Date.now() - startTime,
                filesChecked: 0,
                timestamp: new Date().toISOString()
            }
        }
    }

    /**
     * Get cached verification results (returns empty for now)
     */
    async getVerificationResults(_projectPath: string): Promise<VerificationResult[]> {
        // Cached results not implemented yet
        return []
    }

    /**
     * Request AI fix for an error - emits event for AI chat to handle
     */
    async requestAIFix(error: VerificationError, _projectPath: string): Promise<string> {
        logger.info(`Requesting AI fix for ${error.file}:${error.line}`)

        // Generate a prompt for AI to fix the error
        const prompt = `Fix the following ${error.severity} in ${error.file} at line ${error.line}:
${error.message}${error.rule ? `\nRule: ${error.rule}` : ''}`

        return prompt
    }
}

export const verificationApi = new VerificationApi()
