<template>
  <BaseModal 
    :model-value="true" 
    @update:model-value="handleClose"
    size="sm"
    :show-close="false"
  >
    <template #header>
      <div class="modal-header-content">
        <div class="modal-icon">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h3 class="modal-title">{{ t('chat.executePreview.title') }}</h3>
      </div>
    </template>

    <template #default>
      <p class="modal-text">
        {{ t('chat.executePreview.description', { count: fileCount }) }}
      </p>
      <p class="modal-hint">
        {{ t('chat.executePreview.hint') }}
      </p>
    </template>

    <template #footer>
      <BaseButton 
        variant="secondary" 
        @click="$emit('cancel')"
      >
        {{ t('chat.executePreview.cancel') }}
      </BaseButton>
      <BaseButton 
        variant="primary" 
        @click="$emit('confirm')"
      >
        {{ t('chat.executePreview.confirm') }}
      </BaseButton>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import BaseModal from '@/components/ui/BaseModal.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

defineProps<{
  fileCount: number
}>()

const emit = defineEmits<{
  confirm: []
  cancel: []
}>()

const { t } = useI18n()

function handleClose(value: boolean) {
  if (!value) {
    emit('cancel')
  }
}
</script>

<style scoped>
.modal-header-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
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
</style>
