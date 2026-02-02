import { onMounted, onUnmounted, ref, type Ref } from 'vue'

export type DropdownPlacement = 'bottom-start' | 'bottom-end' | 'top-start' | 'top-end' | 'left' | 'right'

export interface DropdownPosition {
    top: string
    left: string
}

export interface UseDropdownOptions {
    placement?: DropdownPlacement
    offset?: number
    closeOnClick?: boolean
    closeOnEscape?: boolean
}

export interface UseDropdownReturn {
    isOpen: Ref<boolean>
    triggerRef: Ref<HTMLElement | null>
    dropdownRef: Ref<HTMLElement | null>
    position: Ref<DropdownPosition>
    open: () => void
    close: () => void
    toggle: () => void
    calculatePosition: () => void
}

/**
 * Composable for dropdown functionality with positioning, click outside, and keyboard handling
 */
export function useDropdown(options: UseDropdownOptions = {}): UseDropdownReturn {
    const {
        placement = 'bottom-start',
        offset = 8,
        closeOnClick = true,
        closeOnEscape = true
    } = options

    const isOpen = ref(false)
    const triggerRef = ref<HTMLElement | null>(null)
    const dropdownRef = ref<HTMLElement | null>(null)
    const position = ref<DropdownPosition>({ top: '0px', left: '0px' })

    /**
     * Calculate dropdown position based on trigger element and placement
     */
    function calculatePosition(): void {
        if (!triggerRef.value || !dropdownRef.value) return

        const triggerRect = triggerRef.value.getBoundingClientRect()
        const dropdownRect = dropdownRef.value.getBoundingClientRect()
        const viewportWidth = window.innerWidth
        const viewportHeight = window.innerHeight

        let top = 0
        let left = 0

        // Calculate base position based on placement
        switch (placement) {
            case 'bottom-start':
                top = triggerRect.bottom + offset
                left = triggerRect.left
                break
            case 'bottom-end':
                top = triggerRect.bottom + offset
                left = triggerRect.right - dropdownRect.width
                break
            case 'top-start':
                top = triggerRect.top - dropdownRect.height - offset
                left = triggerRect.left
                break
            case 'top-end':
                top = triggerRect.top - dropdownRect.height - offset
                left = triggerRect.right - dropdownRect.width
                break
            case 'left':
                top = triggerRect.top
                left = triggerRect.left - dropdownRect.width - offset
                break
            case 'right':
                top = triggerRect.top
                left = triggerRect.right + offset
                break
        }

        // Adjust if dropdown goes outside viewport (horizontal)
        if (left + dropdownRect.width > viewportWidth) {
            left = viewportWidth - dropdownRect.width - 8
        }
        if (left < 8) {
            left = 8
        }

        // Adjust if dropdown goes outside viewport (vertical)
        if (top + dropdownRect.height > viewportHeight) {
            // Flip to top if there's more space
            if (triggerRect.top > viewportHeight - triggerRect.bottom) {
                top = triggerRect.top - dropdownRect.height - offset
            } else {
                top = viewportHeight - dropdownRect.height - 8
            }
        }
        if (top < 8) {
            top = 8
        }

        position.value = {
            top: `${top}px`,
            left: `${left}px`
        }
    }

    /**
     * Open dropdown
     */
    function open(): void {
        isOpen.value = true
        // Calculate position on next tick after DOM update
        setTimeout(() => {
            calculatePosition()
        }, 0)
    }

    /**
     * Close dropdown
     */
    function close(): void {
        isOpen.value = false
    }

    /**
     * Toggle dropdown open/close
     */
    function toggle(): void {
        if (isOpen.value) {
            close()
        } else {
            open()
        }
    }

    /**
     * Handle click outside to close dropdown
     */
    function handleClickOutside(event: MouseEvent): void {
        if (!isOpen.value) return

        const target = event.target as Node

        // Check if click is outside both trigger and dropdown
        const isOutsideTrigger = triggerRef.value && !triggerRef.value.contains(target)
        const isOutsideDropdown = dropdownRef.value && !dropdownRef.value.contains(target)

        if (isOutsideTrigger && isOutsideDropdown) {
            close()
        }
    }

    /**
     * Handle click inside dropdown
     */
    function handleDropdownClick(): void {
        if (closeOnClick) {
            close()
        }
    }

    /**
     * Handle escape key to close dropdown
     */
    function handleEscape(event: KeyboardEvent): void {
        if (closeOnEscape && isOpen.value && event.key === 'Escape') {
            close()
            // Return focus to trigger
            triggerRef.value?.focus()
        }
    }

    /**
     * Recalculate position on window resize
     */
    function handleResize(): void {
        if (isOpen.value) {
            calculatePosition()
        }
    }

    // Setup event listeners
    onMounted(() => {
        document.addEventListener('click', handleClickOutside)
        document.addEventListener('keydown', handleEscape)
        window.addEventListener('resize', handleResize)
        window.addEventListener('scroll', handleResize, true)

        // Add click handler to dropdown if closeOnClick is enabled
        if (closeOnClick && dropdownRef.value) {
            dropdownRef.value.addEventListener('click', handleDropdownClick)
        }
    })

    // Cleanup event listeners
    onUnmounted(() => {
        document.removeEventListener('click', handleClickOutside)
        document.removeEventListener('keydown', handleEscape)
        window.removeEventListener('resize', handleResize)
        window.removeEventListener('scroll', handleResize, true)

        if (dropdownRef.value) {
            dropdownRef.value.removeEventListener('click', handleDropdownClick)
        }
    })

    return {
        isOpen,
        triggerRef,
        dropdownRef,
        position,
        open,
        close,
        toggle,
        calculatePosition
    }
}
