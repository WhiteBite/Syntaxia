<template>
  <BaseModal 
    :model-value="isOpen"
    @update:model-value="(value) => !value && handleCancel()"
    :title="t('files.sizeGuard.title')"
    size="md"
    :close-on-backdrop="true"
    :show-close="false"
  >
    <template #header>
      <div class="modal-header-content">
        <div class="modal-icon-warning">
          <AlertTriangleIcon class="w-6 h-6" />
        </div>
        <h2 class="modal-title">{{ t('files.sizeGuard.title') }}</h2>
      </div>
    </template>

    <div class="warning-message">
      <p class="warning-text">
        {{ t('files.sizeGuard.message') }}
      </p>
    </div>

    <div class="file-info">
      <div class="file-info-row">
        <span class="file-info-label">{{ t('files.sizeGuard.fileName') }}:</span>
        <span class="file-info-value file-name">{{ fileName }}</span>
      </div>
      <div class="file-info-row">
        <span class="file-info-label">{{ t('files.sizeGuard.fileSize') }}:</span>
        <span class="file-info-value file-size" :class="sizeClass">
          {{ formatTokens(tokens) }} {{ t('files.sizeGuard.tokens') }}
        </span>
      </div>
      <div class="file-info-row">
        <span class="file-info-label">{{ t('files.sizeGuard.contextUsage') }}:</span>
        <span class="file-info-value" :class="sizeClass">
          {{ percentage }}% {{ t('files.sizeGuard.ofLimit') }}
        </span>
      </div>
    </div>

    <div class="warning-details">
      <p class="warning-detail-text">
        {{ t('files.sizeGuard.impact') }}
      </p>
    </div>

    <div class="dont-show-again">
      <label class="checkbox-label">
        <input
          type="checkbox"
          v-model="dontShowAgain"
          class="checkbox-input"
        />
        <span class="checkbox-text">{{ t('files.sizeGuard.dontShowAgain') }}</span>
      </label>
    </div>

    <template #footer>
      <button class="btn btn-ghost" @click="handleCancel">
        {{ t('files.sizeGuard.cancel') }}
      </button>
      <button class="btn btn-warning" @click="handleConfirm">
        {{ t('files.sizeGuard.selectAnyway') }}
      </button>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { TOKEN_THRESHOLDS } from '@/config/constants'
import { AlertTriangle as AlertTriangleIcon } from 'lucide-vue-next'
import BaseModal from '@/components/ui/BaseModal.vue'

const { t } = useI18n()

interface Props {
  isOpen: boolean
  fileName: string
  tokens: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'confirm', dontShowAgain: boolean): void
  (e: 'cancel'): void
}>()

const dontShowAgain = ref(false)

const percentage = computed(() => {
  return Math.round((props.tokens / TOKEN_THRESHOLDS.MAX_CONTEXT) * 100)
})

const sizeClass = computed(() => {
  if (props.tokens >= TOKEN_THRESHOLDS.CRITICAL) return 'size-critical'
  if (props.tokens >= TOKEN_THRESHOLDS.HEAVY) return 'size-heavy'
  return 'size-medium'
})

function formatTokens(tokens: number): string {
  if (tokens < 1000) return tokens.toString()
  if (tokens < 10000) return (tokens / 1000).toFixed(1) + 'k'
  return Math.round(tokens / 1000) + 'k'
}

function handleConfirm() {
  emit('confirm', dontShowAgain.value)
}

function handleCancel() {
  emit('cancel')
}
</script>

<style scoped>
.modal-header-content {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  width: 100%;
}

.modal-icon-warning {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: rgba(251, 146, 60, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fb923c;
  flex-shrink: 0;
}

.modal-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
}

.warning-message {
  margin-bottom: var(--space-6);
}

.warning-text {
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  line-height: 1.6;
  margin: 0;
}

.file-info {
  background: var(--bg-2);
  border-radius: var(--radius-md);
  padding: var(--space-4);
  margin-bottom: var(--space-6);
  border: 1px solid var(--border-subtle);
}

.file-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-2) 0;
}

.file-info-row:not(:last-child) {
  border-bottom: 1px solid var(--border-subtle);
}

.file-info-label {
  font-size: var(--font-size-sm);
  color: var(--text-muted);
  font-weight: var(--font-weight-medium);
}

.file-info-value {
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  font-weight: var(--font-weight-semibold);
  font-family: ui-monospace, monospace;
}

.file-name {
  max-width: 60%;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  text-align: right;
}

.file-size {
  font-size: var(--font-size-base);
}

.size-medium {
  color: #fcd34d;
}

.size-heavy {
  color: #fb923c;
}

.size-critical {
  color: #f87171;
}

.warning-details {
  background: rgba(251, 146, 60, 0.1);
  border-left: 3px solid #fb923c;
  padding: var(--space-4);
  border-radius: var(--radius-sm);
  margin-bottom: var(--space-6);
}

.warning-detail-text {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  line-height: 1.5;
  margin: 0;
}

.dont-show-again {
  display: flex;
  align-items: center;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  width: 16px;
  height: 16px;
  border-radius: 4px;
  border: 1.5px solid #64748b;
  cursor: pointer;
  transition: all var(--transition-fast);
  appearance: none;
  background: transparent;
  position: relative;
}

.checkbox-input:checked {
  background: #6366f1;
  border-color: #6366f1;
}

.checkbox-input:checked::after {
  content: '✓';
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.checkbox-input:hover {
  border-color: #818cf8;
}

.checkbox-text {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
}

.btn-warning {
  background: #fb923c;
  color: white;
  border: none;
}

.btn-warning:hover {
  background: #f97316;
}
</style>
