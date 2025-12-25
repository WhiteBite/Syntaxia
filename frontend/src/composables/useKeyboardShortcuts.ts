/**
 * Composable for managing global keyboard shortcuts
 * 
 * Note: App.vue currently uses @vueuse/core useMagicKeys for shortcuts.
 * This composable provides an alternative approach with better structure
 * and can be used for feature-specific shortcuts.
 */
import { onMounted, onUnmounted, type Ref, unref } from 'vue'

export interface Shortcut {
    /** Key to listen for (case-insensitive) */
    key: string
    /** Require Ctrl/Cmd key */
    ctrl?: boolean
    /** Require Shift key */
    shift?: boolean
    /** Require Alt key */
    alt?: boolean
    /** Action to execute */
    action: () => void
    /** Human-readable description for help dialogs */
    description: string
    /** Optional: only active when condition is true */
    enabled?: boolean | Ref<boolean>
}

/**
 * Register global keyboard shortcuts
 * @param shortcuts Array of shortcut definitions
 * @example
 * useKeyboardShortcuts([
 *   { key: 'b', ctrl: true, action: buildContext, description: 'Build context' },
 *   { key: 'Escape', action: closeModal, description: 'Close modal' }
 * ])
 */
export function useKeyboardShortcuts(shortcuts: Shortcut[]) {
    function handleKeydown(e: KeyboardEvent) {
        // Skip if user is typing in an input
        const target = e.target as HTMLElement
        const isEditable = target.tagName === 'INPUT' ||
            target.tagName === 'TEXTAREA' ||
            target.isContentEditable

        for (const s of shortcuts) {
            // Check if shortcut is enabled
            if (s.enabled !== undefined && !unref(s.enabled)) {
                continue
            }

            // Allow Escape key even in inputs
            if (isEditable && s.key.toLowerCase() !== 'escape') {
                continue
            }

            const ctrlMatch = s.ctrl ? (e.ctrlKey || e.metaKey) : !(e.ctrlKey || e.metaKey)
            const shiftMatch = s.shift ? e.shiftKey : !e.shiftKey
            const altMatch = s.alt ? e.altKey : !e.altKey

            if (
                e.key.toLowerCase() === s.key.toLowerCase() &&
                ctrlMatch &&
                shiftMatch &&
                altMatch
            ) {
                e.preventDefault()
                s.action()
                return
            }
        }
    }

    onMounted(() => {
        document.addEventListener('keydown', handleKeydown)
    })

    onUnmounted(() => {
        document.removeEventListener('keydown', handleKeydown)
    })

    return {
        shortcuts
    }
}
