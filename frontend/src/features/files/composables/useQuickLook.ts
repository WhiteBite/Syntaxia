/**
 * Composable for QuickLook file preview functionality
 * Extracted from useFileExplorer for better separation of concerns
 */
import { useI18n } from '@/composables/useI18n'
import { useUIStore } from '@/stores/ui.store'
import { ref } from 'vue'
import { useFileStore } from '../model/file.store'
import { useHoveredFile } from './useHoveredFile'

export function useQuickLook() {
    const { t } = useI18n()
    const fileStore = useFileStore()
    const uiStore = useUIStore()

    // Get hovered file state from singleton
    const { state: hoveredState } = useHoveredFile()

    const isVisible = ref(false)
    const currentPath = ref('')

    function open(path: string) {
        // Toggle if same file
        if (isVisible.value && currentPath.value === path) {
            close()
            return
        }
        currentPath.value = path
        isVisible.value = true
    }

    function close() {
        isVisible.value = false
    }

    function toggle(path: string) {
        if (isVisible.value && currentPath.value === path) {
            close()
        } else {
            open(path)
        }
    }

    function addToContext(path: string) {
        fileStore.toggleSelect(path)
        uiStore.addToast(t('files.addToContext'), 'success')
    }

    /**
     * Handle spacebar press for QuickLook
     * Opens preview for hovered file or closes if already open
     */
    function handleSpacebarPreview(): boolean {
        const path = hoveredState.path
        const isDir = hoveredState.isDir

        if (path) {
            if (isDir) {
                // For directories, toggle expand
                fileStore.toggleExpand(path)
            } else {
                // For files, toggle QuickLook
                toggle(path)
            }
            return true
        } else if (isVisible.value) {
            // Close QuickLook if open and no file hovered
            close()
            return true
        }

        return false
    }

    return {
        isVisible,
        currentPath,
        open,
        close,
        toggle,
        addToContext,
        handleSpacebarPreview,
    }
}
