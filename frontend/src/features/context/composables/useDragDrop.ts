import { useLogger } from '@/composables/useLogger'
import { ref } from 'vue'

/**
 * Composable for handling drag & drop file operations in context panel.
 */
export function useDragDrop() {
    const logger = useLogger('useDragDrop')

    const isDragging = ref(false)
    let dragCounter = 0

    function handleDragOver(e: DragEvent) {
        dragCounter++
        isDragging.value = true

        if (e.dataTransfer?.types.includes('Files') || e.dataTransfer?.types.includes('text/plain')) {
            e.dataTransfer.dropEffect = 'copy'
        }
    }

    function handleDragLeave() {
        dragCounter--
        if (dragCounter <= 0) {
            dragCounter = 0
            isDragging.value = false
        }
    }

    function handleDrop(e: DragEvent) {
        isDragging.value = false
        dragCounter = 0

        if (!e.dataTransfer) return

        const textData = e.dataTransfer.getData('text/plain')
        if (textData) {
            const paths = textData.split('\n').filter(p => p.trim())
            if (paths.length > 0) {
                logger.debug('Dropped file paths:', paths.length)
                window.dispatchEvent(new CustomEvent('add-files-to-context', { detail: { paths } }))
            }
        }
    }

    return {
        isDragging,
        handleDragOver,
        handleDragLeave,
        handleDrop
    }
}
