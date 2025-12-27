<template>
  <div 
    class="test-item"
    :class="[
      `test-item--${test.status}`,
      { 'test-item--selected': isSelected }
    ]"
    @click="$emit('select', test.id)"
  >
    <!-- Status Icon -->
    <div class="test-status">
      <!-- Running -->
      <svg v-if="test.status === 'running'" class="status-icon animate-spin" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.25"/>
        <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
      <!-- Passed -->
      <svg v-else-if="test.status === 'passed'" class="status-icon" viewBox="0 0 24 24" fill="none">
        <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <!-- Failed -->
      <svg v-else-if="test.status === 'failed'" class="status-icon" viewBox="0 0 24 24" fill="none">
        <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <!-- Skipped -->
      <svg v-else-if="test.status === 'skipped'" class="status-icon" viewBox="0 0 24 24" fill="none">
        <path d="M5 5l14 14M5 19L19 5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.5"/>
      </svg>
      <!-- Pending -->
      <span v-else class="status-dot"></span>
    </div>

    <!-- Test Info -->
    <div class="test-info">
      <span class="test-name">{{ test.name }}</span>
      <span v-if="test.duration !== undefined" class="test-duration">
        {{ formatDuration(test.duration) }}
      </span>
    </div>

    <!-- Run Button -->
    <button
      class="run-btn"
      :disabled="isRunning"
      :title="t('testing.runTest')"
      @click.stop="$emit('run', test.id)"
    >
      <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
        <path d="M5 3l14 9-14 9V3z" fill="currentColor"/>
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { TestResult } from '../types'

const { t } = useI18n()

defineProps<{
  test: TestResult
  isSelected: boolean
  isRunning: boolean
}>()

defineEmits<{
  (e: 'select', testId: string): void
  (e: 'run', testId: string): void
}>()

function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(2)}s`
}
</script>

<style scoped>
.test-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.test-item:hover {
  background: rgba(255, 255, 255, 0.04);
}

.test-item--selected {
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.2);
}

/* Status colors */
.test-item--passed .test-status {
  color: #22c55e;
}

.test-item--failed .test-status {
  color: #ef4444;
}

.test-item--running .test-status {
  color: #3b82f6;
}

.test-item--skipped .test-status {
  color: #6b7280;
}

.test-item--pending .test-status {
  color: #4b5563;
}

.test-status {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-icon {
  width: 16px;
  height: 16px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.5;
}

.test-info {
  flex: 1;
  min-width: 0;
  display: flex;
  align-items: center;
  gap: 8px;
}

.test-name {
  font-size: 13px;
  color: #e5e7eb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.test-duration {
  font-size: 11px;
  color: #6b7280;
  font-family: ui-monospace, monospace;
  flex-shrink: 0;
}

.run-btn {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: none;
  border: none;
  border-radius: 4px;
  color: #6b7280;
  cursor: pointer;
  opacity: 0;
  transition: all 0.15s ease-out;
}

.test-item:hover .run-btn {
  opacity: 1;
}

.run-btn:hover:not(:disabled) {
  background: rgba(139, 92, 246, 0.15);
  color: #a78bfa;
}

.run-btn:disabled {
  cursor: not-allowed;
  opacity: 0.3;
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
