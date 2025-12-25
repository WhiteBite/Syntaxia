/**
 * Composable for managing ignore rules (gitignore and custom)
 * Extracted from useFileExplorer for better separation of concerns
 */
import { useLogger } from '@/composables/useLogger'
import { apiService } from '@/services/api.service'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { filesApi } from '../api/files.api'
import { useFileStore, type FileNode } from '../model/file.store'

const logger = useLogger('IgnoreRules')

export function useIgnoreRules() {
    const fileStore = useFileStore()
    const settingsStore = useSettingsStore()
    const uiStore = useUIStore()

    /**
     * Add a file or folder to custom ignore rules
     */
    async function addToIgnore(node: FileNode): Promise<boolean> {
        try {
            const currentRules = settingsStore.getCustomIgnoreRules()
            const newRule = node.isDir ? `${node.name}/` : node.name
            const updatedRules = currentRules ? `${currentRules}\n${newRule}` : newRule

            await apiService.updateCustomIgnoreRules(updatedRules)
            settingsStore.setCustomIgnoreRules(updatedRules)
            fileStore.removeNode(node.path)

            // Clear caches
            await apiService.clearFileTreeCache()
            filesApi.clearCache()

            uiStore.addToast('Добавлено в исключения', 'success')
            return true
        } catch (error) {
            logger.error('Failed to add to ignore:', error)
            uiStore.addToast('Ошибка добавления в исключения', 'error')
            return false
        }
    }

    /**
     * Remove a file or folder from custom ignore rules
     */
    async function removeFromIgnore(node: FileNode): Promise<boolean> {
        try {
            const currentIgnoreRules = settingsStore.getCustomIgnoreRules()
            const pattern = node.isDir ? `${node.name}/` : node.name

            const lines = currentIgnoreRules.split('\n').filter(line => {
                const trimmed = line.trim()
                return trimmed && !trimmed.includes(pattern) && !trimmed.startsWith('#')
            })

            const updatedRules = lines.join('\n')
            await apiService.updateCustomIgnoreRules(updatedRules)
            settingsStore.setCustomIgnoreRules(updatedRules)

            uiStore.addToast('Удалено из исключений', 'success')
            return true
        } catch (error) {
            logger.error('Failed to remove from ignore:', error)
            uiStore.addToast('Ошибка удаления из исключений', 'error')
            return false
        }
    }

    /**
     * Get current custom ignore rules
     */
    function getCustomRules(): string {
        return settingsStore.getCustomIgnoreRules()
    }

    /**
     * Update custom ignore rules
     */
    async function updateCustomRules(rules: string): Promise<boolean> {
        try {
            await apiService.updateCustomIgnoreRules(rules)
            settingsStore.setCustomIgnoreRules(rules)
            return true
        } catch (error) {
            logger.error('Failed to update ignore rules:', error)
            return false
        }
    }

    return {
        addToIgnore,
        removeFromIgnore,
        getCustomRules,
        updateCustomRules,
    }
}
