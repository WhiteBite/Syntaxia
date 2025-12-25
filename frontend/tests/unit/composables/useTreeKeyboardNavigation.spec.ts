import { useTreeKeyboardNavigation } from '@/features/files/composables/useTreeKeyboardNavigation'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { ref } from 'vue'

describe('useTreeKeyboardNavigation', () => {
    let container: HTMLElement
    let onToggleExpand: ReturnType<typeof vi.fn>
    let onToggleSelect: ReturnType<typeof vi.fn>

    beforeEach(() => {
        // Create mock container with tree rows
        container = document.createElement('div')
        for (let i = 0; i < 5; i++) {
            const row = document.createElement('div')
            row.className = 'tree-row'
            row.setAttribute('tabindex', i === 0 ? '0' : '-1')
            row.dataset.path = `/path/to/file${i}`
            row.dataset.isDir = i % 2 === 0 ? 'true' : 'false'
            container.appendChild(row)
        }
        document.body.appendChild(container)

        onToggleExpand = vi.fn()
        onToggleSelect = vi.fn()
    })

    it('should initialize with focusedIndex -1', () => {
        const containerRef = ref(container)
        const { focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        expect(focusedIndex.value).toBe(-1)
    })

    it('should focus first row', () => {
        const containerRef = ref(container)
        const { focusFirst, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusFirst()
        expect(focusedIndex.value).toBe(0)
    })

    it('should focus last row', () => {
        const containerRef = ref(container)
        const { focusLast, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusLast()
        expect(focusedIndex.value).toBe(4)
    })

    it('should navigate to next row', () => {
        const containerRef = ref(container)
        const { focusFirst, focusNext, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusFirst()
        focusNext()
        expect(focusedIndex.value).toBe(1)
    })

    it('should navigate to previous row', () => {
        const containerRef = ref(container)
        const { focusRow, focusPrev, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusRow(2)
        focusPrev()
        expect(focusedIndex.value).toBe(1)
    })

    it('should not go below 0 when navigating up', () => {
        const containerRef = ref(container)
        const { focusFirst, focusPrev, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusFirst()
        focusPrev()
        expect(focusedIndex.value).toBe(0)
    })

    it('should not exceed max index when navigating down', () => {
        const containerRef = ref(container)
        const { focusLast, focusNext, focusedIndex } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        focusLast()
        focusNext()
        expect(focusedIndex.value).toBe(4)
    })

    it('should handle Enter key to toggle select', () => {
        const containerRef = ref(container)
        const { handleKeydown } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        const row = container.querySelector('.tree-row') as HTMLElement
        const event = new KeyboardEvent('keydown', { key: 'Enter' })
        Object.defineProperty(event, 'target', { value: row })
        event.preventDefault = vi.fn()

        handleKeydown(event)
        expect(onToggleSelect).toHaveBeenCalledWith('/path/to/file0')
    })

    it('should handle Space key to toggle select', () => {
        const containerRef = ref(container)
        const { handleKeydown } = useTreeKeyboardNavigation({
            containerRef,
            onToggleExpand,
            onToggleSelect,
        })

        const row = container.querySelector('.tree-row') as HTMLElement
        const event = new KeyboardEvent('keydown', { key: ' ' })
        Object.defineProperty(event, 'target', { value: row })
        event.preventDefault = vi.fn()

        handleKeydown(event)
        expect(onToggleSelect).toHaveBeenCalledWith('/path/to/file0')
    })
})
