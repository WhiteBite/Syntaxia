import { useLogger } from '@/composables/useLogger'
import { noiseReductionService, type NoiseReductionResult } from '@/services/noiseReduction.service'
import { useSettingsStore } from '@/stores/settings.store'
import { computed } from 'vue'

const logger = useLogger('useNoiseReduction')

/**
 * Composable for applying noise reduction to context content
 */
export function useNoiseReduction() {
    const settingsStore = useSettingsStore()
    const settings = computed(() => settingsStore.settings.context)

    /**
     * Check if noise reduction is enabled
     */
    const isEnabled = computed(() => settings.value.enableNoiseReduction)

    /**
     * Get current noise reduction options
     */
    const options = computed(() => ({
        collapseImports: settings.value.noiseReductionCollapseImports,
        removeComments: settings.value.noiseReductionRemoveComments,
        removeTypeDefinitions: settings.value.noiseReductionRemoveTypeDefinitions,
        collapseBoilerplate: true // Always collapse boilerplate
    }))

    /**
     * Apply noise reduction to file content
     */
    function cleanupContent(content: string, filename: string): NoiseReductionResult {
        if (!isEnabled.value) {
            const tokens = Math.ceil(content.length / 4)
            return {
                content,
                originalTokens: tokens,
                cleanedTokens: tokens,
                savings: 0,
                savingsPercent: 0
            }
        }

        try {
            const language = noiseReductionService.getLanguageFromExtension(filename)
            const result = noiseReductionService.cleanupFileContent(content, language, options.value)

            if (result.savings > 0) {
                logger.debug(`Noise reduction for ${filename}: saved ${result.savings} tokens (${result.savingsPercent.toFixed(1)}%)`)
            }

            return result
        } catch (error) {
            logger.error(`Failed to apply noise reduction to ${filename}:`, error)
            // Return original content on error
            const tokens = Math.ceil(content.length / 4)
            return {
                content,
                originalTokens: tokens,
                cleanedTokens: tokens,
                savings: 0,
                savingsPercent: 0
            }
        }
    }

    /**
     * Apply noise reduction to multiple files
     */
    function cleanupMultipleFiles(files: Array<{ path: string; content: string }>): {
        files: Array<{ path: string; content: string; result: NoiseReductionResult }>
        totalSavings: number
        totalSavingsPercent: number
    } {
        let totalOriginal = 0
        let totalCleaned = 0

        const processedFiles = files.map(file => {
            const result = cleanupContent(file.content, file.path)
            totalOriginal += result.originalTokens
            totalCleaned += result.cleanedTokens

            return {
                path: file.path,
                content: result.content,
                result
            }
        })

        const totalSavings = totalOriginal - totalCleaned
        const totalSavingsPercent = totalOriginal > 0 ? (totalSavings / totalOriginal) * 100 : 0

        return {
            files: processedFiles,
            totalSavings,
            totalSavingsPercent
        }
    }

    /**
     * Get estimated savings percentage based on current settings
     */
    function getEstimatedSavings(): number {
        if (!isEnabled.value) return 0

        let savingsPercent = 0
        if (options.value.collapseImports) savingsPercent += 10
        if (options.value.removeComments) savingsPercent += 8
        if (options.value.removeTypeDefinitions) savingsPercent += 5

        return Math.min(savingsPercent, 30) // Cap at 30%
    }

    return {
        isEnabled,
        options,
        cleanupContent,
        cleanupMultipleFiles,
        getEstimatedSavings
    }
}
