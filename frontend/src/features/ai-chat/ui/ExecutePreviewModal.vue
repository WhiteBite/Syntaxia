<template>
  <div class="modal-backdrop" @click.self="$emit('cancel')">
    <div class="modal-content">
      <div class="modal-header">
        <div class="modal-icon">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h3 class="modal-title">{{ t('chat.executePreview.title') }}</h3>
      </div>
      
      <div class="modal-body">
        <p class="modal-text">
          {{ t('chat.executePreview.description', { count: fileCount }) }}
        </p>
        <p class="modal-hint">
          {{ t('chat.executePreview.hint') }}
        </p>
      </div>
      
      <div class="modal-actions">
        <button 
          class="btn btn-secondary"
          @click="$emit('cancel')"
        >
          {{ t('chat.executePreview.cancel') }}
        </button>
        <button 
          class="btn btn-primary"
          @click="$emit('confirm')"
        >
          {{ t('chat.executePreview.confirm') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'

defineProps<{
  fileCount: number
}>()

defineEmits<{
  confirm: []
  cancel: []
}>()

const { t } = useI18n()
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal);
  padding: 1rem;
}

.modal-content {
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  width: 100%;
  max-width: min(400px, 90vw);
  padding: 1.5rem;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.modal-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 2.5rem;
  height: 2.5rem;
  border-radius: var(--radius-lg);
  background: var(--color-warning-soft);
  color: var(--color-warning);
  flex-shrink: 0;
}

.modal-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.modal-body {
  margin-bottom: 1.5rem;
}

.modal-text {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin: 0 0 0.5rem 0;
  line-height: 1.5;
}

.modal-hint {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin: 0;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
}
</style>
