import { ref } from 'vue'

export function useContextDragDrop() {
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

    function handleDrop(e: DragEvent): string[] | null {
        isDragging.value = false
        dragCounter = 0

        if (!e.dataTransfer) return null

        // Handle file paths from internal drag
        const textData = e.dataTransfer.getData('text/plain')
        if (textData) {
            const paths = textData.split('\n').filter(p => p.trim())
            if (paths.length > 0) {
                return paths
            }
        }

        return null
    }

    return {
        isDragging,
        handleDragOver,
        handleDragLeave,
        handleDrop,
    }
}
