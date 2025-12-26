/**
 * Composable for keyboard navigation in file tree
 * Implements roving tabindex pattern for accessibility
 */
import { ref, type Ref } from 'vue'

export interface TreeNavigationOptions {
    containerRef: Ref<HTMLElement | null>
    onToggleExpand?: (path: string) => void
    onToggleSelect?: (path: string) => void
}

export function useTreeKeyboardNavigation(options: TreeNavigationOptions) {
    const focusedIndex = ref(-1)

    function getVisibleRows(): HTMLElement[] {
        if (!options.containerRef.value) return []
        return Array.from(options.containerRef.value.querySelectorAll('.tree-row[tabindex]'))
    }

    function focusRow(index: number): void {
        const rows = getVisibleRows()
        if (index < 0 || index >= rows.length) return

        // Update tabindex (roving tabindex pattern)
        rows.forEach((row, i) => {
            row.setAttribute('tabindex', i === index ? '0' : '-1')
        })

        rows[index].focus()
        focusedIndex.value = index
    }

    function focusFirst(): void {
        focusRow(0)
    }

    function focusLast(): void {
        const rows = getVisibleRows()
        focusRow(rows.length - 1)
    }

    function focusNext(): void {
        const rows = getVisibleRows()
        const currentIndex = focusedIndex.value
        if (currentIndex < rows.length - 1) {
            focusRow(currentIndex + 1)
        }
    }

    function focusPrev(): void {
        const currentIndex = focusedIndex.value
        if (currentIndex > 0) {
            focusRow(currentIndex - 1)
        }
    }

    function handleKeydown(event: KeyboardEvent): void {
        const target = event.target as HTMLElement
        if (!target.classList.contains('tree-row')) return

        const path = target.dataset.path
        const isDir = target.dataset.isDir === 'true'

        switch (event.key) {
            case 'ArrowDown':
                event.preventDefault()
                focusNext()
                break
            case 'ArrowUp':
                event.preventDefault()
                focusPrev()
                break
            case 'Home':
                event.preventDefault()
                focusFirst()
                break
            case 'End':
                event.preventDefault()
                focusLast()
                break
            case 'ArrowRight':
                if (isDir && path) {
                    event.preventDefault()
                    options.onToggleExpand?.(path)
                }
                break
            case 'ArrowLeft':
                if (isDir && path) {
                    event.preventDefault()
                    options.onToggleExpand?.(path)
                }
                break
            case 'Enter':
            case ' ':
                if (path) {
                    event.preventDefault()
                    options.onToggleSelect?.(path)
                }
                break
        }
    }

    function setupNavigation(): void {
        const container = options.containerRef.value
        if (!container) return

        container.addEventListener('keydown', handleKeydown)

        // Initialize first row as focusable
        const rows = getVisibleRows()
        rows.forEach((row, i) => {
            row.setAttribute('tabindex', i === 0 ? '0' : '-1')
        })
    }

    function cleanupNavigation(): void {
        const container = options.containerRef.value
        if (!container) return

        container.removeEventListener('keydown', handleKeydown)
    }

    return {
        focusedIndex,
        focusFirst,
        focusLast,
        focusNext,
        focusPrev,
        focusRow,
        handleKeydown,
        setupNavigation,
        cleanupNavigation,
    }
}
