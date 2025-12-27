/**
 * Symbols Feature Module - Public API
 * 
 * Provides UI for navigating code symbols (classes, functions, methods, interfaces)
 */

// Store
export { useSymbolsStore } from './model/symbols.store'

// API
export { symbolsApi } from './api/symbols.api'

// Types
export type { Symbol, SymbolFilter, SymbolGroup, SymbolKind, SymbolTreeNode } from './types'

// UI Components
export { default as SymbolBrowser } from './ui/SymbolBrowser.vue'
export { default as SymbolItem } from './ui/SymbolItem.vue'
export { default as SymbolSearch } from './ui/SymbolSearch.vue'
export { default as SymbolTree } from './ui/SymbolTree.vue'

