export { contextApi } from './api/context.api'
export * from './lib/context-utils'
export { generateSmartName, useContextStore } from './model/context.store'

// UI Components
export { default as ContextBuilder } from './ui/ContextBuilder.vue'
export { default as ContextDeleteModal } from './ui/ContextDeleteModal.vue'
export { default as ContextItemActions } from './ui/ContextItemActions.vue'
export { default as ContextList } from './ui/ContextList.vue'
export { default as ContextListEmpty } from './ui/ContextListEmpty.vue'
export { default as ContextListItem } from './ui/ContextListItem.vue'
export { default as ContextListToolbar } from './ui/ContextListToolbar.vue'
export { default as ContextPanel } from './ui/ContextPanel.vue'

// Composables
export { useChunking } from './composables/useChunking'
export { useContextDragDrop } from './composables/useContextDragDrop'
export { useContextList } from './composables/useContextList'
export { useContextSearch } from './composables/useContextSearch'
export * from './composables/useSyntaxHighlight'

// Chunking UI
export { default as ChunkCutLine } from './ui/ChunkCutLine.vue'

