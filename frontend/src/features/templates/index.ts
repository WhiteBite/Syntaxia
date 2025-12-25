// Template feature exports
export { useTemplateStore } from './model/template.store'
export * from './model/template.types'

// Composables
export { detectLanguages, generateFileTree } from './composables/useFileTree'

// UI Components
export { default as TemplateCard } from './ui/TemplateCard.vue'
export { default as TemplateListItem } from './ui/TemplateListItem.vue'
export { default as TemplateModal } from './ui/TemplateModal.vue'
export { default as TemplateOptionTile } from './ui/TemplateOptionTile.vue'
export { default as TemplatePreviewBlock } from './ui/TemplatePreviewBlock.vue'
export { default as TemplateSection } from './ui/TemplateSection.vue'
export { default as TemplateSelector } from './ui/TemplateSelector.vue'

