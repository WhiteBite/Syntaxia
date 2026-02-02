<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="isOpen" class="modal-backdrop" @click="handleBackdropClick">
        <div class="modal-container" @click.stop>
          <div class="modal-header">
            <div class="modal-icon-warning">
              <AlertTriangleIcon class="w-6 h-6" />
            </div>
            <h2 class="modal-title">{{ t('files.sizeGuard.title') }}</h2>
          </div>

          <div class="modal-content">
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
          </div>

          <div class="modal-footer">
            <button class="btn btn-ghost" @click="handleCancel">
              {{ t('files.sizeGuard.cancel') }}
            </button>
            <button class="btn btn-warning" @click="handleConfirm">
              {{ t('files.sizeGuard.selectAnyway') }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useI18n } from '@/composables/useI18n'
import { TOKEN_THRESHOLDS } from '@/config/constants'
import { AlertTriangle as AlertTriangleIcon } from 'lucide-vue-next'

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

function handleBackdropClick() {
  handleCancel()
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  backdrop-filter: blur(2px);
}

.modal-container {
  background: var(--bg-secondary);
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border-subtle);
}

.modal-header {
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  gap: var(--space-4);
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

.modal-content {
  padding: var(--space-6);
  overflow-y: auto;
  flex: 1;
  min-height: 0;
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
  background: var(--bg-tertiary);
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

.modal-footer {
  padding: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
}

.btn-warning {
  background: #fb923c;
  color: white;
  border: none;
}

.btn-warning:hover {
  background: #f97316;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.2s ease;
}

.modal-enter-from .modal-container,
.modal-leave-to .modal-container {
  transform: scale(0.95);
}
</style>
