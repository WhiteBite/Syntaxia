import { onClickOutside, onKeyStroke } from '@vueuse/core'
import { onUnmounted, ref, watch, type Ref } from 'vue'

/**
 * Position coordinates for overlay elements
 */
export interface OverlayPosition {
    top: string
    left: string
    transform?: string
}

/**
 * Configuration options for useOverlay composable
 */
export interface UseOverlayOptions {
    /**
     * Initial visibility state
     * @default false
     */
    initialState?: boolean

    /**
     * Close overlay on ESC key press
     * @default true
     */
    closeOnEscape?: boolean

    /**
     * Close overlay when clicking outside
     * @default true
     */
    closeOnClickOutside?: boolean

    /**
     * Lock body scroll when overlay is open
     * @default false
     */
    lockScroll?: boolean

    /**
     * Z-index for the overlay
     * @default 1000
     */
    zIndex?: number

    /**
     * Callback when overlay opens
     */
    onOpen?: () => void

    /**
     * Callback when overlay closes
     */
    onClose?: () => void

    /**
     * Elements to ignore when detecting click outside
     */
    ignoreElements?: Ref<HTMLElement | null>[]
}

/**
 * Return type for useOverlay composable
 */
export interface UseOverlayReturn {
    /**
     * Reactive visibility state
     */
    isVisible: Ref<boolean>

    /**
     * Reactive position state
     */
    position: Ref<OverlayPosition>

    /**
     * Reference to the overlay element
     */
    overlayRef: Ref<HTMLElement | null>

    /**
     * Open the overlay
     */
    open: () => void

    /**
     * Close the overlay
     */
    close: () => void

    /**
     * Toggle overlay visibility
     */
    toggle: () => void

    /**
     * Set overlay position
     */
    setPosition: (pos: Partial<OverlayPosition>) => void

    /**
     * Get current z-index
     */
    getZIndex: () => number
}

/**
 * Base composable for overlay components (modals, popovers, dropdowns)
 * 
 * Provides common functionality:
 * - Visibility state management
 * - Position management
 * - Click outside detection
 * - Escape key handling
 * - Scroll lock (optional)
 * - Z-index management
 * 
 * @example
 * ```ts
 * const {
 *   isVisible,
 *   overlayRef,
 *   open,
 *   close,
 *   setPosition
 * } = useOverlay({
 *   closeOnEscape: true,
 *   closeOnClickOutside: true,
 *   lockScroll: true,
 *   onOpen: () => console.log('Opened'),
 *   onClose: () => console.log('Closed')
 * })
 * ```
 */
export function useOverlay(options: UseOverlayOptions = {}): UseOverlayReturn {
    const {
        initialState = false,
        closeOnEscape = true,
        closeOnClickOutside = true,
        lockScroll = false,
        zIndex = 1000,
        onOpen,
        onClose,
        ignoreElements = []
    } = options

    const isVisible = ref(initialState)
    const position = ref<OverlayPosition>({
        top: '0px',
        left: '0px',
        transform: ''
    })
    const overlayRef = ref<HTMLElement | null>(null)

    /**
     * Lock body scroll to prevent background scrolling
     */
    function lockBodyScroll(): void {
        if (!lockScroll) return

        const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth
        document.body.style.overflow = 'hidden'
        document.body.style.paddingRight = `${scrollbarWidth}px`
    }

    /**
     * Unlock body scroll
     */
    function unlockBodyScroll(): void {
        if (!lockScroll) return

        document.body.style.overflow = ''
        document.body.style.paddingRight = ''
    }

    /**
     * Open the overlay
     */
    function open(): void {
        if (isVisible.value) return

        isVisible.value = true
        lockBodyScroll()
        onOpen?.()
    }

    /**
     * Close the overlay
     */
    function close(): void {
        if (!isVisible.value) return

        isVisible.value = false
        unlockBodyScroll()
        onClose?.()
    }

    /**
     * Toggle overlay visibility
     */
    function toggle(): void {
        if (isVisible.value) {
            close()
        } else {
            open()
        }
    }

    /**
     * Set overlay position
     */
    function setPosition(pos: Partial<OverlayPosition>): void {
        position.value = {
            ...position.value,
            ...pos
        }
    }

    /**
     * Get current z-index
     */
    function getZIndex(): number {
        return zIndex
    }

    // Click outside detection
    if (closeOnClickOutside) {
        onClickOutside(
            overlayRef,
            (event) => {
                // Check if click is on any ignored elements
                const isIgnored = ignoreElements.some(
                    (el) => el.value?.contains(event.target as Node)
                )
                if (!isIgnored && isVisible.value) {
                    close()
                }
            },
            { ignore: ignoreElements }
        )
    }

    // Escape key handling
    if (closeOnEscape) {
        onKeyStroke('Escape', (e) => {
            if (isVisible.value) {
                e.preventDefault()
                close()
            }
        })
    }

    // Watch for visibility changes
    watch(isVisible, (newValue) => {
        if (newValue) {
            lockBodyScroll()
        } else {
            unlockBodyScroll()
        }
    })

    // Cleanup on unmount
    onUnmounted(() => {
        if (isVisible.value) {
            unlockBodyScroll()
        }
    })

    return {
        isVisible,
        position,
        overlayRef,
        open,
        close,
        toggle,
        setPosition,
        getZIndex
    }
}
