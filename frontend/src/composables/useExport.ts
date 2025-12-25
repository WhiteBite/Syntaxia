import type { domain } from '#wailsjs/go/models'
import { useLogger } from '@/composables/useLogger'
import { useContextStore } from '@/features/context/model/context.store'
import { apiService } from '@/services/api.service'
import { useSettingsStore } from '@/stores/settings.store'
import { computed, ref } from 'vue'

const logger = useLogger('Export')

export type ExportMode = 'clipboard' | 'ai' | 'human'

/**
 * Export modal composable for managing export functionality
 * Uses settingsStore as single source of truth for export settings
 */
export function useExport() {
  const contextStore = useContextStore()
  const settingsStore = useSettingsStore()

  const isOpen = ref(false)
  const isExporting = ref(false)
  const selectedMode = ref<ExportMode>('clipboard')
  const exportResult = ref<domain.ExportResult | null>(null)
  const error = ref<string | null>(null)

  // Computed settings from settingsStore (single source of truth)
  const settings = computed(() => ({
    // From settingsStore.context
    exportFormat: settingsStore.settings.context.outputFormat === 'xml' ? 'manifest' : 'plain',
    stripComments: settingsStore.settings.context.stripComments,
    includeManifest: settingsStore.settings.context.includeManifest,
    tokenLimit: settingsStore.settings.context.maxTokens,
    enableAutoSplit: settingsStore.settings.context.enableAutoSplit,
    maxTokensPerChunk: settingsStore.settings.context.maxTokensPerChunk,
    includeLineNumbers: settingsStore.settings.context.includeLineNumbers,
    // Fixed values
    aiProfile: settingsStore.settings.aiModel,
    overlapTokens: 200,
    splitStrategy: settingsStore.settings.context.splitStrategy,
    theme: 'default',
    includePageNumbers: true
  }))

  /**
   * Open export modal
   */
  function open() {
    if (!contextStore.hasContext) {
      error.value = 'Сначала соберите контекст'
      return false
    }

    isOpen.value = true
    exportResult.value = null
    error.value = null
    return true
  }

  /**
   * Close export modal
   */
  function close() {
    isOpen.value = false
  }

  /**
   * Execute export with current settings
   */
  async function executeExport() {
    if (!contextStore.hasContext) {
      error.value = 'Контекст не найден'
      return false
    }

    isExporting.value = true
    error.value = null

    try {
      // Get full context content (warning: can be large)
      const contextContent = await contextStore.getFullContextContent()

      // Build export settings JSON
      const s = settings.value
      const exportSettingsJson = {
        mode: selectedMode.value,
        context: contextContent,

        // Clipboard settings
        stripComments: s.stripComments,
        includeManifest: s.includeManifest,
        exportFormat: s.exportFormat,

        // AI settings
        aiProfile: s.aiProfile,
        tokenLimit: s.tokenLimit,
        fileSizeLimitKB: 5120, // 5 MB
        enableAutoSplit: s.enableAutoSplit,
        maxTokensPerChunk: s.maxTokensPerChunk,
        overlapTokens: s.overlapTokens,
        splitStrategy: s.splitStrategy,

        // Human settings
        theme: s.theme,
        includeLineNumbers: s.includeLineNumbers,
        includePageNumbers: s.includePageNumbers
      }

      const result = await apiService.exportContext(exportSettingsJson)
      exportResult.value = result

      // Handle result based on mode
      if (result.mode === 'clipboard' && result.text) {
        await navigator.clipboard.writeText(result.text)
        logger.debug('Content copied to clipboard')
      } else if (result.filePath) {
        logger.debug('File exported to:', result.filePath)
      } else if (result.dataBase64) {
        // Download file
        downloadBase64File(result.dataBase64, result.fileName || 'export.zip')
      }

      // Auto close after success
      setTimeout(() => {
        close()
      }, 1500)

      return true
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Ошибка экспорта'
      console.error('Export failed:', err)
      return false
    } finally {
      isExporting.value = false
    }
  }

  /**
   * Download base64 encoded file
   */
  function downloadBase64File(base64: string, filename: string) {
    const byteCharacters = atob(base64)
    const byteNumbers = new Array(byteCharacters.length)
    for (let i = 0; i < byteCharacters.length; i++) {
      byteNumbers[i] = byteCharacters.charCodeAt(i)
    }
    const byteArray = new Uint8Array(byteNumbers)
    const blob = new Blob([byteArray])

    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return {
    // State
    isOpen,
    isExporting,
    selectedMode,
    settings,
    exportResult,
    error,
    // Include context store for use in components
    contextStore,

    // Actions
    open,
    close,
    executeExport
  }
}
