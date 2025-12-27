import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useFileStore } from '@/features/files'
import { useUIStore } from '@/stores/ui.store'
import { ref } from 'vue'

const DRAG_DATA_TYPE = 'application/x-file-paths'

/**
 * Composable for handling drag & drop files into chat panel.
 * Adds dropped files to context via file store selection.
 */
export function useChatDragDrop() {
    const logger = useLogger('useChatDragDrop')
    const { t } = useI18n()
    const fileStore = useFileStore()
    const uiStore = useUIStore()

    const isDragOver = ref(false)
    let dragCounter = 0

    function handleDragEnter(e: DragEvent) {
        e.preventDefault()
        dragCounter++

        if (e.dataTransfer?.types.includes(DRAG_DATA_TYPE)) {
            isDragOver.value = true
            e.dataTransfer.dropEffect = 'copy'
        }
    }

    function handleDragOver(e: DragEvent) {
        e.preventDefault()

        if (e.dataTransfer?.types.includes(DRAG_DATA_TYPE)) {
            e.dataTransfer.dropEffect = 'copy'
        }
    }

    function handleDragLeave(e: DragEvent) {
        e.preventDefault()
        dragCounter--

        if (dragCounter <= 0) {
            dragCounter = 0
            isDragOver.value = false
        }
    }

    function handleDrop(e: DragEvent) {
        e.preventDefault()
        isDragOver.value = false
        dragCounter = 0

        if (!e.dataTransfer) return

        const data = e.dataTransfer.getData(DRAG_DATA_TYPE)
        if (!data) return

        try {
            const paths: string[] = JSON.parse(data)
            if (!Array.isArray(paths) || paths.length === 0) return

            addFilesToContext(paths)
        } catch (err) {
            logger.error('Failed to parse dropped file paths:', err)
        }
    }

    function addFilesToContext(paths: string[]) {
        const currentFiles = fileStore.selectedPaths
        const newFiles: string[] = []
        const duplicates: string[] = []

        for (const path of paths) {
            if (currentFiles.has(path)) {
                duplicates.push(path)
            } else {
                newFiles.push(path)
            }
        }

        if (newFiles.length > 0) {
            fileStore.selectMultiple(newFiles)
            uiStore.addToast(
                t('chat.dragDrop.filesAdded', { count: newFiles.length }),
                'success'
            )
            logger.debug('Added files to context:', newFiles.length)
        }

        if (duplicates.length > 0 && newFiles.length === 0) {
            uiStore.addToast(
                t('chat.dragDrop.filesAlreadyInContext', { count: duplicates.length }),
                'warning'
            )
        }
    }

    return {
        isDragOver,
        handleDragEnter,
        handleDragOver,
        handleDragLeave,
        handleDrop
    }
}

export { DRAG_DATA_TYPE }
