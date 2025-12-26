<template>
  <div v-if="canSave" class="context-footer">
    <button @click="$emit('save')" class="context-save-btn">
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
      </svg>
      {{ t('context.saveCurrentContext') }}
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'

const { t } = useI18n()

defineProps<{
  canSave: boolean
}>()

defineEmits<{
  (e: 'save'): void
}>()
</script>

<style scoped>
.context-footer {
  padding: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: linear-gradient(180deg, transparent 0%, rgba(139, 92, 246, 0.03) 100%);
}

.context-save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 
    0 4px 12px rgba(139, 92, 246, 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.context-save-btn:hover {
  background: linear-gradient(135deg, #a78bfa 0%, #8b5cf6 100%);
  transform: translateY(-1px);
  box-shadow: 
    0 6px 20px rgba(139, 92, 246, 0.45),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}
</style>
