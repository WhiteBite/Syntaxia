<template>
  <!-- Delete Dialog -->
  <Transition name="fade">
    <div v-if="deleteConfirmId" class="tpl-dialog-overlay" @click.self="$emit('cancel-delete')">
      <div class="tpl-dialog danger">
        <Trash2 class="w-6 h-6" />
        <h3>{{ t('templates.deleteConfirm') }}</h3>
        <p>{{ t('templates.deleteConfirmText') }}</p>
        <div class="tpl-dialog-btns">
          <button @click="$emit('cancel-delete')" class="tpl-cancel-btn">
            {{ t('common.cancel') }}
          </button>
          <button @click="$emit('confirm-delete')" class="tpl-danger-btn">
            {{ t('templates.delete') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
  
  <!-- Unsaved Dialog -->
  <Transition name="fade">
    <div v-if="showUnsavedWarning" class="tpl-dialog-overlay" @click.self="$emit('cancel-unsaved')">
      <div class="tpl-dialog warning">
        <AlertTriangle class="w-6 h-6" />
        <h3>{{ t('templates.unsavedChanges') }}</h3>
        <p>{{ t('templates.unsavedChangesText') }}</p>
        <div class="tpl-dialog-btns">
          <button @click="$emit('discard')" class="tpl-cancel-btn">
            {{ t('templates.discard') }}
          </button>
          <button @click="$emit('save-and-close')" class="tpl-save-btn">
            {{ t('common.save') }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { AlertTriangle, Trash2 } from 'lucide-vue-next'

const { t } = useI18n()

defineProps<{
  deleteConfirmId: string | null
  showUnsavedWarning: boolean
}>()

defineEmits<{
  (e: 'cancel-delete'): void
  (e: 'confirm-delete'): void
  (e: 'cancel-unsaved'): void
  (e: 'discard'): void
  (e: 'save-and-close'): void
}>()
</script>

<style scoped>
.tpl-dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: calc(var(--z-modal) + 10);
}

.tpl-dialog {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  padding: 1.5rem 2rem;
  background: #0D0E12;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-lg);
  text-align: center;
  max-width: 320px;
}

.tpl-dialog.danger svg { color: var(--color-danger); }
.tpl-dialog.warning svg { color: var(--color-warning); }

.tpl-dialog h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.tpl-dialog p {
  font-size: 12px;
  color: var(--text-muted);
  margin: 0;
}

.tpl-dialog-btns {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.5rem;
}

.tpl-cancel-btn {
  padding: 0.5rem 0.75rem;
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-md);
  color: var(--text-muted);
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-cancel-btn:hover {
  background: rgba(255, 255, 255, 0.05);
  color: var(--text-primary);
}

.tpl-danger-btn {
  padding: 0.5rem 0.875rem;
  background: var(--color-danger-soft);
  border: 1px solid var(--color-danger-border);
  border-radius: var(--radius-md);
  color: var(--color-danger);
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-danger-btn:hover {
  background: rgba(248, 113, 113, 0.25);
}

.tpl-save-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.5rem 0.875rem;
  background: var(--accent-indigo);
  border: none;
  border-radius: var(--radius-md);
  color: white;
  font-size: 11px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
}

.tpl-save-btn:hover:not(:disabled) { 
  background: #818cf8;
  transform: translateY(-1px);
}

/* Fade Transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
