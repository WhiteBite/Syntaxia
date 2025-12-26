import { onMounted, onUnmounted, ref } from 'vue'

export type TooltipPlacement = 'top' | 'bottom' | 'left' | 'right'

export interface TooltipOptions {
    delay?: number
    placement?: TooltipPlacement
    offset?: number
}

export function useTooltip(options: TooltipOptions = {}) {
    const isVisible = ref(false)
    const position = ref({ x: 0, y: 0 })
    const content = ref('')
    const placement = ref(options.placement || 'top')
    const delay = options.delay || 300
    const offset = options.offset || 8

    let showTimer: ReturnType<typeof setTimeout> | null = null
    let targetElement: HTMLElement | null = null

    function calculatePosition(element: HTMLElement, preferredPlacement: TooltipPlacement) {
        const rect = element.getBoundingClientRect()
        const tooltipWidth = 320 // max-width
        const tooltipHeight = 200 // approximate height

        let x = 0
        let y = 0
        let finalPlacement: TooltipPlacement = preferredPlacement

        // Calculate position based on placement
        switch (preferredPlacement) {
            case 'top':
                x = rect.left + rect.width / 2
                y = rect.top - offset
                // Check if tooltip fits above
                if (y - tooltipHeight < 0) {
                    finalPlacement = 'bottom'
                    y = rect.bottom + offset
                }
                break
            case 'bottom':
                x = rect.left + rect.width / 2
                y = rect.bottom + offset
                // Check if tooltip fits below
                if (y + tooltipHeight > window.innerHeight) {
                    finalPlacement = 'top'
                    y = rect.top - offset
                }
                break
            case 'left':
                x = rect.left - offset
                y = rect.top + rect.height / 2
                // Check if tooltip fits on left
                if (x - tooltipWidth < 0) {
                    finalPlacement = 'right'
                    x = rect.right + offset
                }
                break
            case 'right':
                x = rect.right + offset
                y = rect.top + rect.height / 2
                // Check if tooltip fits on right
                if (x + tooltipWidth > window.innerWidth) {
                    finalPlacement = 'left'
                    x = rect.left - offset
                }
                break
        }

        // Ensure tooltip doesn't go off-screen horizontally
        if (x + tooltipWidth / 2 > window.innerWidth) {
            x = window.innerWidth - tooltipWidth / 2 - 10
        }
        if (x - tooltipWidth / 2 < 0) {
            x = tooltipWidth / 2 + 10
        }

        placement.value = finalPlacement
        position.value = { x, y }
    }

    function show(element: HTMLElement, tooltipContent: string, preferredPlacement?: TooltipPlacement) {
        targetElement = element

        // Clear any existing timer
        if (showTimer) {
            clearTimeout(showTimer)
        }

        // Set content immediately
        content.value = tooltipContent

        // Show with delay
        showTimer = setTimeout(() => {
            calculatePosition(element, preferredPlacement || placement.value)
            isVisible.value = true
            showTimer = null
        }, delay)
    }

    function hide() {
        // Clear timer if tooltip hasn't shown yet
        if (showTimer) {
            clearTimeout(showTimer)
            showTimer = null
        }

        isVisible.value = false
        targetElement = null
    }

    function updatePosition() {
        if (targetElement && isVisible.value) {
            calculatePosition(targetElement, placement.value)
        }
    }

    // Update position on scroll/resize
    onMounted(() => {
        window.addEventListener('scroll', updatePosition, true)
        window.addEventListener('resize', updatePosition)
    })

    onUnmounted(() => {
        window.removeEventListener('scroll', updatePosition, true)
        window.removeEventListener('resize', updatePosition)
        if (showTimer) {
            clearTimeout(showTimer)
        }
    })

    return {
        isVisible,
        position,
        content,
        placement,
        show,
        hide,
        updatePosition
    }
}
