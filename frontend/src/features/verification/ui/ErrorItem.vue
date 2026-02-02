<template>
  <div 
    class="error-item"
    :class="[`error-item--${error.severity}`]"
  >
    <!-- Severity Icon -->
    <div class="error-icon">
      <svg v-if="error.severity === 'error'" class="icon-error" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <svg v-else-if="error.severity === 'warning'" class="icon-warning" viewBox="0 0 24 24" fill="none">
        <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <svg v-else class="icon-info" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        <path d="M12 16v-4m0-4h.01" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </div>

    <!-- Error Content -->
    <div class="error-content">
      <!-- Location -->
      <button 
        class="error-location"
        @click="$emit('go-to-file', error)"
        :title="t('verification.goToFile')"
      >
        <span class="error-file">{{ fileName }}</span>
        <span class="error-line">:{{ error.line }}</span>
        <span v-if="error.column" class="error-column">:{{ error.column }}</span>
      </button>

      <!-- Message -->
      <p class="error-message">{{ error.message }}</p>

      <!-- Rule (if present) -->
      <span v-if="error.rule" class="error-rule">
        {{ error.rule }}
      </span>
    </div>

    <!-- Actions -->
    <div class="error-actions">
      <BaseButton
        variant="secondary"
        size="sm"
        :disabled="isFixing"
        :loading="isFixing"
        @click="$emit('fix-with-ai', error.id)"
        :title="t('verification.fixWithAI')"
      >
        <template #icon>
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
            <path d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </template>
        <span class="action-text">{{ isFixing ? t('verification.fixing') : t('verification.fixWithAI') }}</span>
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import BaseButton from '@/components/ui/BaseButton.vue'
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'
import type { VerificationError } from '../api/verification.api'

const { t } = useI18n()

const props = defineProps<{
  error: VerificationError
  isFixing: boolean
}>()

defineEmits<{
  (e: 'go-to-file', error: VerificationError): void
  (e: 'fix-with-ai', errorId: string): void
}>()

// Extract filename from path
const fileName = computed(() => {
  const parts = props.error.file.split(/[/\\]/)
  return parts[parts.length - 1] || props.error.file
})
</script>

<style scoped>
.error-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.15s ease-out;
}

.error-item:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.1);
}

/* Severity variants */
.error-item--error {
  border-left: 3px solid #ef4444;
}

.error-item--warning {
  border-left: 3px solid #f59e0b;
}

.error-item--info {
  border-left: 3px solid #3b82f6;
}

/* Icon */
.error-icon {
  flex-shrink: 0;
  width: 20px;
  height: 20px;
  margin-top: 2px;
}

.error-icon svg {
  width: 100%;
  height: 100%;
}

.icon-error {
  color: #ef4444;
}

.icon-warning {
  color: #f59e0b;
}

.icon-info {
  color: #3b82f6;
}

/* Content */
.error-content {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.error-location {
  display: inline-flex;
  align-items: center;
  gap: 0;
  padding: 2px 8px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.2);
  border-radius: 4px;
  font-family: ui-monospace, monospace;
  font-size: 12px;
  color: #a78bfa;
  cursor: pointer;
  transition: all 0.15s ease-out;
  width: fit-content;
}

.error-location:hover {
  background: rgba(139, 92, 246, 0.2);
  border-color: rgba(139, 92, 246, 0.4);
  color: #c4b5fd;
}

.error-file {
  color: #e5e7eb;
}

.error-line,
.error-column {
  color: #9ca3af;
}

.error-message {
  font-size: 13px;
  color: #d1d5db;
  line-height: 1.5;
  margin: 0;
  word-break: break-word;
}

.error-rule {
  display: inline-block;
  padding: 2px 6px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
  font-family: ui-monospace, monospace;
  font-size: 10px;
  color: #6b7280;
  width: fit-content;
}

/* Actions */
.error-actions {
  flex-shrink: 0;
  opacity: 0;
  transform: translateX(8px);
  transition: all 0.15s ease-out;
}

.error-item:hover .error-actions {
  opacity: 1;
  transform: translateX(0);
}

.action-text {
  display: none;
}

@media (min-width: 640px) {
  .action-text {
    display: inline;
  }
}
</style>
