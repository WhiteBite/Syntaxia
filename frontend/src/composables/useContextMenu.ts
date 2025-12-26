import type { FileNode } from '@/types/domain';
import { onMounted, onUnmounted, ref, type Ref } from 'vue';

export interface ContextMenuState {
    isVisible: Ref<boolean>
    position: Ref<{ x: number; y: number }>
    targetNode: Ref<FileNode | null>
}

export function useContextMenu() {
    const isVisible = ref(false)
    const position = ref({ x: 0, y: 0 })
    const targetNode = ref<FileNode | null>(null)

    function show(node: FileNode, event: MouseEvent) {
        targetNode.value = node

        // Calculate position with viewport bounds check
        const viewportWidth = window.innerWidth
        const viewportHeight = window.innerHeight
        const menuWidth = 250 // Approximate menu width
        const menuHeight = 300 // Approximate menu height

        let x = event.clientX
        let y = event.clientY

        // Prevent menu from going off-screen horizontally
        if (x + menuWidth > viewportWidth) {
            x = viewportWidth - menuWidth - 10
        }

        // Prevent menu from going off-screen vertically
        if (y + menuHeight > viewportHeight) {
            y = viewportHeight - menuHeight - 10
        }

        position.value = { x, y }
        isVisible.value = true
    }

    function hide() {
        isVisible.value = false
        targetNode.value = null
    }

    function handleClickOutside(event: MouseEvent) {
        const target = event.target as HTMLElement
        if (!target.closest('.context-menu')) {
            hide()
        }
    }

    function handleEscape(event: KeyboardEvent) {
        if (event.key === 'Escape') {
            hide()
        }
    }

    onMounted(() => {
        document.addEventListener('click', handleClickOutside)
        document.addEventListener('keydown', handleEscape)
    })

    onUnmounted(() => {
        document.removeEventListener('click', handleClickOutside)
        document.removeEventListener('keydown', handleEscape)
    })

    return {
        isVisible,
        position,
        targetNode,
        show,
        hide
    }
}
