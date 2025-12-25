/**
 * Composable for handling keyboard shortcuts
 * Automatically registers and unregisters event listeners
 */

import { onMounted, onUnmounted } from 'vue'

export type KeyboardShortcut = string // e.g., "Ctrl+K", "Ctrl+Shift+P"
export type ShortcutHandler = () => void

export function useKeyboard(shortcuts: Record<KeyboardShortcut, ShortcutHandler>) {
  const handleKeyDown = (e: KeyboardEvent) => {
    // Build key combination string
    const parts: string[] = []
    
    if (e.ctrlKey) parts.push('Ctrl')
    if (e.shiftKey) parts.push('Shift')
    if (e.altKey) parts.push('Alt')
    if (e.metaKey) parts.push('Meta')
    
    // Add the actual key (lowercase for consistency)
    const key = e.key.length === 1 ? e.key.toLowerCase() : e.key
    parts.push(key)
    
    const combination = parts.join('+')
    
    // Check if we have a handler for this combination
    const handler = shortcuts[combination]
    if (handler) {
      e.preventDefault()
      handler()
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
  })

  return {
    // Could expose methods to add/remove shortcuts dynamically if needed
  }
}
