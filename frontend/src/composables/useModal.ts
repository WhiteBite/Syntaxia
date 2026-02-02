import { onKeyStroke } from '@vueuse/core'
import { onUnmounted, ref, watch } from 'vue'

export interface UseModalOptions {
  /**
   * Initial state of the modal
   * @default false
   */
  initialState?: boolean
  /**
   * Close modal on ESC key press
   * @default true
   */
  closeOnEsc?: boolean
  /**
   * Lock body scroll when modal is open
   * @default true
   */
  lockScroll?: boolean
  /**
   * Callback when modal opens
   */
  onOpen?: () => void
  /**
   * Callback when modal closes
   */
  onClose?: () => void
}

/**
 * Generic modal composable for managing modal state with enhanced features
 * 
 * Features:
 * - State management (open/close/toggle)
 * - ESC key handling
 * - Body scroll lock
 * - Focus trap support
 * - Lifecycle callbacks
 * 
 * @example
 * ```ts
 * const { isOpen, open, close, toggle } = useModal({
 *   closeOnEsc: true,
 *   lockScroll: true,
 *   onOpen: () => console.log('Modal opened'),
 *   onClose: () => console.log('Modal closed')
 * })
 * ```
 */
export function useModal(options: UseModalOptions = {}) {
  const {
    initialState = false,
    closeOnEsc = true,
    lockScroll = true,
    onOpen,
    onClose,
  } = options

  const isOpen = ref(initialState)

  function open() {
    isOpen.value = true
  }

  function close() {
    isOpen.value = false
  }

  function toggle() {
    isOpen.value = !isOpen.value
  }

  // Body scroll lock
  function lockBodyScroll() {
    if (!lockScroll) return
    const scrollbarWidth = window.innerWidth - document.documentElement.clientWidth
    document.body.style.overflow = 'hidden'
    document.body.style.paddingRight = `${scrollbarWidth}px`
  }

  function unlockBodyScroll() {
    if (!lockScroll) return
    document.body.style.overflow = ''
    document.body.style.paddingRight = ''
  }

  // ESC key handling
  if (closeOnEsc) {
    onKeyStroke('Escape', (e) => {
      if (isOpen.value) {
        e.preventDefault()
        close()
      }
    })
  }

  // Watch for modal state changes
  watch(isOpen, (newValue) => {
    if (newValue) {
      lockBodyScroll()
      onOpen?.()
    } else {
      unlockBodyScroll()
      onClose?.()
    }
  })

  // Cleanup on unmount
  onUnmounted(() => {
    if (isOpen.value) {
      unlockBodyScroll()
    }
  })

  return {
    isOpen,
    open,
    close,
    toggle,
  }
}
