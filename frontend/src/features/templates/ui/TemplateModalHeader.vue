<template>
  <div class="tpl-header">
    <div class="tpl-header-left">
      <FileText class="w-4 h-4 text-purple-400" />
      <h2>{{ t('templates.manage') }}</h2>
    </div>
    <div class="tpl-header-center">
      <span v-if="hasChanges && !isEditingBuiltIn" class="tpl-unsaved-indicator" />
    </div>
    <button @click="$emit('close')" class="tpl-close">
      <X class="w-4 h-4" />
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { FileText, X } from 'lucide-vue-next'

const { t } = useI18n()

defineProps<{
  hasChanges: boolean
  isEditingBuiltIn: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>

<style scoped>
.tpl-header {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);
}

.tpl-header-left {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.tpl-header h2 {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.tpl-header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.tpl-unsaved-indicator {
  width: 6px;
  height: 6px;
  background: var(--color-warning);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

.tpl-close {
  padding: 0.375rem;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--text-muted);
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-close:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--text-primary);
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
