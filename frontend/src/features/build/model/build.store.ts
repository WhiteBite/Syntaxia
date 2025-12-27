/**
 * Build feature store
 */

import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { buildFeatureApi } from '../api/build.api'
import type { BuildHistoryEntry, BuildResult, BuildStatusType } from '../types'

const logger = useLogger('BuildStore')

export const useBuildStore = defineStore('build', () => {
    const projectStore = useProjectStore()
    const uiStore = useUIStore()
    const { t } = useI18n()

    // State
    const status = ref<BuildStatusType>('idle')
    const progress = ref<number | undefined>(undefined)
    const currentStep = ref<string | undefined>(undefined)
    const output = ref<string>('')
    const history = ref<BuildHistoryEntry[]>([])
    const currentBuild = ref<BuildResult | null>(null)
    const selectedHistoryId = ref<string | null>(null)
    const isCancelling = ref(false)

    // Getters
    const isBuilding = computed(() => status.value === 'building')

    const lastBuild = computed(() => {
        if (history.value.length > 0) {
            return history.value[0]
        }
        return null
    })

    const lastBuildTime = computed(() => {
        if (lastBuild.value) {
            return new Date(lastBuild.value.endTime).toLocaleString()
        }
        return null
    })

    const selectedHistoryEntry = computed(() => {
        if (!selectedHistoryId.value) return null
        return history.value.find(h => h.id === selectedHistoryId.value) || null
    })

    const displayOutput = computed(() => {
        if (selectedHistoryEntry.value) {
            return selectedHistoryEntry.value.output
        }
        return output.value
    })

    // Actions
    async function build(): Promise<BuildResult | null> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('build.selectProject'), 'warning')
            return null
        }

        if (isBuilding.value) {
            uiStore.addToast(t('build.alreadyBuilding'), 'warning')
            return null
        }

        status.value = 'building'
        output.value = ''
        currentStep.value = t('build.compiling')
        selectedHistoryId.value = null

        try {
            logger.info(`Starting build for ${projectPath}`)

            // Simulate progress updates
            const progressInterval = setInterval(() => {
                if (progress.value === undefined) {
                    progress.value = 0
                }
                if (progress.value < 90) {
                    progress.value += Math.random() * 15
                }
            }, 500)

            const result = await buildFeatureApi.build(projectPath)

            clearInterval(progressInterval)
            progress.value = 100

            currentBuild.value = result
            output.value = result.output
            status.value = result.status === 'success' ? 'success' : 'failed'

            // Reload history
            await loadHistory()

            if (result.status === 'success') {
                uiStore.addToast(t('build.success'), 'success')
            } else {
                uiStore.addToast(t('build.failed'), 'error')
            }

            logger.info(`Build completed: ${result.status}`)

            // Reset progress after a delay
            setTimeout(() => {
                progress.value = undefined
                currentStep.value = undefined
            }, 1000)

            return result
        } catch (error) {
            logger.error('Build failed:', error)
            status.value = 'failed'
            output.value = error instanceof Error ? error.message : 'Build failed'
            uiStore.addToast(t('build.failed'), 'error')

            progress.value = undefined
            currentStep.value = undefined

            return null
        }
    }

    async function clean(): Promise<BuildResult | null> {
        const projectPath = projectStore.currentPath
        if (!projectPath) {
            uiStore.addToast(t('build.selectProject'), 'warning')
            return null
        }

        if (isBuilding.value) {
            uiStore.addToast(t('build.alreadyBuilding'), 'warning')
            return null
        }

        status.value = 'building'
        output.value = ''
        currentStep.value = t('build.cleaning')
        selectedHistoryId.value = null

        try {
            logger.info(`Cleaning build for ${projectPath}`)

            const result = await buildFeatureApi.clean(projectPath)

            currentBuild.value = result
            output.value = result.output
            status.value = result.status === 'success' ? 'success' : 'failed'

            if (result.status === 'success') {
                uiStore.addToast(t('build.cleanSuccess'), 'success')
            } else {
                uiStore.addToast(t('build.cleanFailed'), 'error')
            }

            logger.info(`Clean completed: ${result.status}`)

            currentStep.value = undefined

            return result
        } catch (error) {
            logger.error('Clean failed:', error)
            status.value = 'failed'
            output.value = error instanceof Error ? error.message : 'Clean failed'
            uiStore.addToast(t('build.cleanFailed'), 'error')

            currentStep.value = undefined

            return null
        }
    }

    async function rebuild(): Promise<BuildResult | null> {
        await clean()
        return await build()
    }

    async function cancelBuild(): Promise<void> {
        if (!isBuilding.value) return

        isCancelling.value = true

        try {
            await buildFeatureApi.cancelBuild()
            status.value = 'cancelled'
            output.value += '\n\nBuild cancelled by user.'
            uiStore.addToast(t('build.cancelled'), 'info')
            logger.info('Build cancelled')
        } catch (error) {
            logger.error('Failed to cancel build:', error)
        } finally {
            isCancelling.value = false
            progress.value = undefined
            currentStep.value = undefined
        }
    }

    async function loadHistory(): Promise<void> {
        const projectPath = projectStore.currentPath
        if (!projectPath) return

        try {
            history.value = await buildFeatureApi.getBuildHistory(projectPath)
        } catch (error) {
            logger.warn('Failed to load build history:', error)
        }
    }

    function selectHistoryEntry(id: string | null): void {
        selectedHistoryId.value = id
    }

    function clearOutput(): void {
        output.value = ''
        selectedHistoryId.value = null
    }

    function copyOutput(): void {
        const textToCopy = displayOutput.value
        if (textToCopy) {
            navigator.clipboard.writeText(textToCopy)
                .then(() => {
                    uiStore.addToast(t('build.outputCopied'), 'success')
                })
                .catch(() => {
                    uiStore.addToast(t('build.copyFailed'), 'error')
                })
        }
    }

    // Initialize history on project change
    function init(): void {
        if (projectStore.currentPath) {
            loadHistory()
        }
    }

    return {
        // State
        status,
        progress,
        currentStep,
        output,
        history,
        currentBuild,
        selectedHistoryId,
        isCancelling,

        // Getters
        isBuilding,
        lastBuild,
        lastBuildTime,
        selectedHistoryEntry,
        displayOutput,

        // Actions
        build,
        clean,
        rebuild,
        cancelBuild,
        loadHistory,
        selectHistoryEntry,
        clearOutput,
        copyOutput,
        init
    }
})
