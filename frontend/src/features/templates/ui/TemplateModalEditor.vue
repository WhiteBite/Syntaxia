<template>
  <main class="tpl-editor">
    <template v-if="editingTemplate">
      <TemplateEditorHeader
        :template="editingTemplate"
        :is-current-template="isCurrentTemplate"
        @toggle-emoji-picker="$emit('toggle-emoji-picker')"
        @update-field="updateField"
      />
      
      <TemplateEditorContent
        ref="contentRef"
        :template="editingTemplate"
        @toggle-section="toggleSection"
        @update-field="updateField"
      />
    </template>
    
    <div v-else class="tpl-empty">
      <FileText class="w-10 h-10 opacity-15" />
      <p>{{ t('templates.selectToEdit') }}</p>
    </div>
  </main>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { FileText } from 'lucide-vue-next'
import { ref } from 'vue'
import type { PromptTemplate, TemplateSections } from '../model/template.types'
import TemplateEditorContent from './TemplateEditorContent.vue'
import TemplateEditorHeader from './TemplateEditorHeader.vue'

const { t } = useI18n()

const props = defineProps<{
  editingTemplate: PromptTemplate | null
  isCurrentTemplate: boolean
}>()

const emit = defineEmits<{
  (e: 'update:editingTemplate', template: PromptTemplate): void
  (e: 'toggle-emoji-picker'): void
}>()

const contentRef = ref<InstanceType<typeof TemplateEditorContent> | null>(null)

function updateField(field: keyof PromptTemplate, value: string) {
  if (!props.editingTemplate) return
  emit('update:editingTemplate', {
    ...props.editingTemplate,
    [field]: value
  })
}

function toggleSection(key: keyof TemplateSections) {
  if (!props.editingTemplate) return
  emit('update:editingTemplate', {
    ...props.editingTemplate,
    sections: {
      ...props.editingTemplate.sections,
      [key]: !props.editingTemplate.sections[key]
    }
  })
}

defineExpose({
  resizeTextareas: () => contentRef.value?.resizeTextareas()
})
</script>

<style scoped>
.tpl-editor {
  display: flex;
  flex-direction: column;
  background: #0D0E12;
  min-width: 0;
  overflow: hidden;
}

.tpl-empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  color: var(--text-subtle);
  font-size: 12px;
}
</style>
