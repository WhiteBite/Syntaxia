<template>
  <div class="tpl-editor-header">
    <BaseButton
      variant="ghost"
      size="sm"
      icon-only
      @click="$emit('toggle-emoji-picker')"
      :disabled="template.isBuiltIn"
      class="tpl-icon-btn"
    >
      <template #icon>
        <span class="tpl-icon-display">{{ template.icon }}</span>
      </template>
    </BaseButton>
    <div class="tpl-name-group">
      <input 
        :value="template.name"
        @input="$emit('update-field', 'name', ($event.target as HTMLInputElement).value)"
        class="tpl-name-input" 
        :placeholder="t('templates.namePlaceholder')" 
        :disabled="template.isBuiltIn" 
      />
      <input 
        :value="template.description"
        @input="$emit('update-field', 'description', ($event.target as HTMLInputElement).value)"
        class="tpl-desc-input" 
        :placeholder="t('templates.descriptionPlaceholder')" 
        :disabled="template.isBuiltIn" 
      />
    </div>
    <span v-if="isCurrentTemplate" class="tpl-active-badge">
      <Check class="w-3 h-3" />{{ t('templates.currentTemplate') }}
    </span>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'
import { Check } from 'lucide-vue-next'
import type { PromptTemplate } from '../model/template.types'

const { t } = useI18n()

defineProps<{
  template: PromptTemplate
  isCurrentTemplate: boolean
}>()

defineEmits<{
  (e: 'toggle-emoji-picker'): void
  (e: 'update-field', field: keyof PromptTemplate, value: string): void
}>()
</script>

<style scoped>
.tpl-editor-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.tpl-icon-btn {
  width: 2.5rem;
  height: 2.5rem;
}

.tpl-icon-display {
  font-size: 1.25rem;
  line-height: 1;
}

.tpl-name-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
  min-width: 0;
}

.tpl-name-input {
  background: transparent;
  border: none;
  border-bottom: 1px solid transparent;
  color: var(--text-primary);
  font-size: 15px;
  font-weight: 600;
  padding: 0.125rem 0;
  transition: border-color 0.15s;
}

.tpl-name-input:hover:not(:disabled) {
  border-bottom-color: rgba(255, 255, 255, 0.1);
}

.tpl-name-input:focus {
  outline: none;
  border-bottom-color: var(--accent-indigo);
}

.tpl-desc-input {
  background: transparent;
  border: none;
  color: var(--text-muted);
  font-size: 11px;
  padding: 0.125rem 0;
}

.tpl-desc-input:focus {
  outline: none;
  color: var(--text-secondary);
}

.tpl-active-badge {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.625rem;
  background: var(--color-success-soft);
  border: 1px solid var(--color-success-border);
  border-radius: var(--radius-full);
  font-size: 10px;
  font-weight: 500;
  color: var(--color-success);
  white-space: nowrap;
}
</style>
