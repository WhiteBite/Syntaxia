/**
 * Global component type definitions
 * Enables TypeScript autocomplete for global components
 */

import CommandPalette from '../components/CommandPalette.vue'
import KeyboardShortcutsModal from '../components/KeyboardShortcutsModal.vue'
import ThemeToggle from '../components/ThemeToggle.vue'

declare module 'vue' {
  export interface GlobalComponents {
    CommandPalette: typeof CommandPalette
    KeyboardShortcutsModal: typeof KeyboardShortcutsModal
    ThemeToggle: typeof ThemeToggle
  }
}