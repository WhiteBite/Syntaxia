<template>
  <div class="context-stats">
    <label class="text-xs text-gray-400 mb-2 block">{{ t('context.outputFormat') }}</label>
    <div class="grid grid-cols-2 gap-2">
      <button 
        v-for="format in formats" 
        :key="format.id" 
        @click="$emit('update:modelValue', format.id)" 
        :title="format.tooltip"
        :class="[
          'btn-unified text-xs',
          modelValue === format.id ? 'btn-unified-primary' : 'btn-unified-secondary'
        ]"
      >
        {{ format.label }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import type { OutputFormat } from '@/stores/settings.store';

const { t } = useI18n()

defineProps<{
  modelValue: OutputFormat
}>()

defineEmits<{
  'update:modelValue': [value: OutputFormat]
}>()

const formats: { id: OutputFormat; label: string; tooltip: string }[] = [
  { id: 'xml', label: 'XML', tooltip: 'Best for AI - structured, easy to parse' },
  { id: 'markdown', label: 'Markdown', tooltip: 'Good for documentation and readability' },
  { id: 'plain', label: 'Plain Text', tooltip: 'Simple format with separators' },
]
</script>
