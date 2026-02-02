<template>
  <BaseModal 
    :model-value="show"
    @update:model-value="(value) => !value && $emit('close')"
    :title="t('context.confirmDelete')"
    size="md"
    :close-on-backdrop="true"
    @close="$emit('close')"
  >
    <template #header>
      <div class="modal-header-content">
        <div class="modal-icon-danger">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </div>
        <div>
          <h3 class="modal-title">{{ t('context.confirmDelete') }}</h3>
          <p class="modal-subtitle">{{ message }}</p>
        </div>
      </div>
    </template>

    <template #footer>
      <button @click="$emit('close')" class="btn btn-ghost">
        {{ t('context.cancel') }}
      </button>
      <button @click="$emit('confirm')" class="btn btn-danger">
        {{ t('context.delete') }}
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n';
import BaseModal from '@/components/ui/BaseModal.vue';

const { t } = useI18n()

defineProps<{
  show: boolean
  message: string
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'confirm'): void
}>()
</script>

<style scoped>
.modal-header-content {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  width: 100%;
}

.modal-icon-danger {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-xl);
  background: rgba(239, 68, 68, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #f87171;
  flex-shrink: 0;
}

.modal-title {
  font-size: var(--font-size-lg);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0 0 var(--space-1) 0;
}

.modal-subtitle {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
  margin: 0;
}

.btn-danger {
  background: #ef4444;
  color: white;
  border: none;
}

.btn-danger:hover {
  background: #dc2626;
}
</style>
