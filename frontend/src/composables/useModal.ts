import { useOverlay } from './ui/useOverlay'

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
 * Built on top of useOverlay composable for consistent overlay behavior.
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

  // Use base overlay functionality
  const overlay = useOverlay({
    initialState,
    closeOnEscape: closeOnEsc,
    closeOnClickOutside: false, // Modals typically don't close on outside click
    lockScroll,
    zIndex: 1000,
    onOpen,
    onClose
  })

  return {
    isOpen: overlay.isVisible,
    open: overlay.open,
    close: overlay.close,
    toggle: overlay.toggle,
  }
}
