import { useEventListener } from '@vueuse/core'
import { nextTick, ref, watch, type Ref } from 'vue'
import { useOverlay } from './useOverlay'

export type PopoverPlacement =
    | 'top'
    | 'bottom'
    | 'left'
    | 'right'
    | 'top-start'
    | 'top-end'
    | 'bottom-start'
    | 'bottom-end'

export type PopoverTrigger = 'click' | 'hover' | 'manual'

export interface PopoverPosition {
    top: string
    left: string
    transform: string
}

export interface ArrowPosition {
    top?: string
    left?: string
    bottom?: string
    right?: string
    transform?: string
}

export interface UsePopoverOptions {
    placement?: PopoverPlacement
    offset?: number
    trigger?: PopoverTrigger
    delay?: number
    arrow?: boolean
    onOpen?: () => void
    onClose?: () => void
}

export interface UsePopoverReturn {
    isOpen: Ref<boolean>
    open: () => void
    close: () => void
    toggle: () => void
    triggerRef: Ref<HTMLElement | null>
    popoverRef: Ref<HTMLElement | null>
    position: Ref<PopoverPosition>
    arrowPosition: Ref<ArrowPosition>
    calculatePosition: () => void
}

export function usePopover(options: UsePopoverOptions = {}): UsePopoverReturn {
    const {
        placement = 'top',
        offset = 8,
        trigger = 'click',
        delay = 0,
        onOpen,
        onClose
    } = options

    const triggerRef = ref<HTMLElement | null>(null)
    const popoverRef = ref<HTMLElement | null>(null)
    const arrowPosition = ref<ArrowPosition>({})

    let hoverTimeout: ReturnType<typeof setTimeout> | null = null
    let closeTimeout: ReturnType<typeof setTimeout> | null = null

    // Use base overlay functionality
    const overlay = useOverlay({
        closeOnEscape: true,
        closeOnClickOutside: trigger !== 'manual',
        lockScroll: false,
        zIndex: 1000,
        ignoreElements: [triggerRef],
        onOpen: () => {
            onOpen?.()
            nextTick(() => {
                calculatePosition()
            })
        },
        onClose
    })

    // Assign popoverRef to overlayRef for click outside detection
    watch(popoverRef, (newRef) => {
        overlay.overlayRef.value = newRef
    })

    function calculatePosition() {
        if (!triggerRef.value || !popoverRef.value) return

        const triggerRect = triggerRef.value.getBoundingClientRect()
        const popoverRect = popoverRef.value.getBoundingClientRect()

        const viewport = {
            width: window.innerWidth,
            height: window.innerHeight
        }

        let top = 0
        let left = 0
        let arrowTop: string | undefined
        let arrowLeft: string | undefined
        let arrowBottom: string | undefined
        let arrowRight: string | undefined
        let arrowTransform: string | undefined

        // Calculate base position
        switch (placement) {
            case 'top':
                top = triggerRect.top - popoverRect.height - offset
                left = triggerRect.left + triggerRect.width / 2 - popoverRect.width / 2
                arrowTop = '100%'
                arrowLeft = '50%'
                arrowTransform = 'translateX(-50%)'
                break

            case 'top-start':
                top = triggerRect.top - popoverRect.height - offset
                left = triggerRect.left
                arrowTop = '100%'
                arrowLeft = `${Math.min(triggerRect.width / 2, 20)}px`
                break

            case 'top-end':
                top = triggerRect.top - popoverRect.height - offset
                left = triggerRect.right - popoverRect.width
                arrowTop = '100%'
                arrowRight = `${Math.min(triggerRect.width / 2, 20)}px`
                break

            case 'bottom':
                top = triggerRect.bottom + offset
                left = triggerRect.left + triggerRect.width / 2 - popoverRect.width / 2
                arrowBottom = '100%'
                arrowLeft = '50%'
                arrowTransform = 'translateX(-50%)'
                break

            case 'bottom-start':
                top = triggerRect.bottom + offset
                left = triggerRect.left
                arrowBottom = '100%'
                arrowLeft = `${Math.min(triggerRect.width / 2, 20)}px`
                break

            case 'bottom-end':
                top = triggerRect.bottom + offset
                left = triggerRect.right - popoverRect.width
                arrowBottom = '100%'
                arrowRight = `${Math.min(triggerRect.width / 2, 20)}px`
                break

            case 'left':
                top = triggerRect.top + triggerRect.height / 2 - popoverRect.height / 2
                left = triggerRect.left - popoverRect.width - offset
                arrowLeft = '100%'
                arrowTop = '50%'
                arrowTransform = 'translateY(-50%)'
                break

            case 'right':
                top = triggerRect.top + triggerRect.height / 2 - popoverRect.height / 2
                left = triggerRect.right + offset
                arrowRight = '100%'
                arrowTop = '50%'
                arrowTransform = 'translateY(-50%)'
                break
        }

        // Auto-correct position if out of viewport
        const padding = 8

        // Horizontal correction
        if (left < padding) {
            left = padding
        } else if (left + popoverRect.width > viewport.width - padding) {
            left = viewport.width - popoverRect.width - padding
        }

        // Vertical correction
        if (top < padding) {
            top = padding
        } else if (top + popoverRect.height > viewport.height - padding) {
            top = viewport.height - popoverRect.height - padding
        }

        overlay.setPosition({
            top: `${top}px`,
            left: `${left}px`,
            transform: ''
        })

        arrowPosition.value = {
            top: arrowTop,
            left: arrowLeft,
            bottom: arrowBottom,
            right: arrowRight,
            transform: arrowTransform
        }
    }

    function open() {
        if (overlay.isVisible.value) return

        if (closeTimeout) {
            clearTimeout(closeTimeout)
            closeTimeout = null
        }

        const doOpen = () => {
            overlay.open()
        }

        if (delay > 0 && trigger === 'hover') {
            hoverTimeout = setTimeout(doOpen, delay)
        } else {
            doOpen()
        }
    }

    function close() {
        if (!overlay.isVisible.value) return

        if (hoverTimeout) {
            clearTimeout(hoverTimeout)
            hoverTimeout = null
        }

        const doClose = () => {
            overlay.close()
        }

        if (trigger === 'hover') {
            closeTimeout = setTimeout(doClose, 100)
        } else {
            doClose()
        }
    }

    function toggle() {
        if (overlay.isVisible.value) {
            close()
        } else {
            open()
        }
    }

    // Recalculate position on scroll/resize
    watch(overlay.isVisible, (open) => {
        if (open) {
            useEventListener(window, 'scroll', calculatePosition, { passive: true })
            useEventListener(window, 'resize', calculatePosition, { passive: true })
        }
    })

    return {
        isOpen: overlay.isVisible,
        open,
        close,
        toggle,
        triggerRef,
        popoverRef,
        position: overlay.position as Ref<PopoverPosition>,
        arrowPosition,
        calculatePosition
    }
}
