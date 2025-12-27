<template>
  <button
    class="history-item"
    :class="{
      'history-item--selected': isSelected,
      [`history-item--${entry.status}`]: true
    }"
    @click="$emit('select', entry.id)"
  >
    <!-- Status Icon -->
    <div class="history-icon">
      <svg v-if="entry.status === 'success'" class="icon-success" viewBox="0 0 24 24" fill="none">
        <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <svg v-else-if="entry.status === 'failed'" class="icon-failed" viewBox="0 0 24 24" fill="none">
        <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <svg v-else class="icon-cancelled" viewBox="0 0 24 24" fill="none">
        <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </div>

    <!-- Content -->
    <div class="history-content">
      <div class="history-time">{{ formattedTime }}</div>
      <div class="history-duration">{{ formattedDuration }}</div>
    </div>

    <!-- Stats -->
    <div class="history-stats">
      <span v-if="entry.errorsCount > 0" class="stat stat--error">
        {{ entry.errorsCount }} {{ t('build.errors') }}
      </span>
      <span v-if="entry.warningsCount > 0" class="stat stat--warning">
        {{ entry.warningsCount }} {{ t('build.warnings') }}
      </span>
    </div>
  </button>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed } from 'vue'
import type { BuildHistoryEntry } from '../types'

const { t } = useI18n()

const props = defineProps<{
  entry: BuildHistoryEntry
  isSelected: boolean
}>()

defineEmits<{
  (e: 'select', id: string): void
}>()

const formattedTime = computed(() => {
  const date = new Date(props.entry.startTime)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
})

const formattedDuration = computed(() => {
  const seconds = Math.round(props.entry.duration / 1000)
  if (seconds < 60) {
    return `${seconds}s`
  }
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}m ${remainingSeconds}s`
})
</script>

<style scoped>
.history-item {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease-out;
  text-align: left;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.1);
}

.history-item--selected {
  background: rgba(139, 92, 246, 0.1);
  border-color: rgba(139, 92, 246, 0.3);
}

/* Status variants */
.history-item--success .history-icon {
  color: #22c55e;
}

.history-item--failed .history-icon {
  color: #ef4444;
}

.history-item--cancelled .history-icon {
  color: #6b7280;
}

/* Icon */
.history-icon {
  flex-shrink: 0;
  width: 18px;
  height: 18px;
}

.history-icon svg {
  width: 100%;
  height: 100%;
}

/* Content */
.history-content {
  flex: 1;
  min-width: 0;
}

.history-time {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
}

.history-duration {
  font-size: 11px;
  color: #6b7280;
  margin-top: 2px;
}

/* Stats */
.history-stats {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
}

.stat {
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 10px;
  font-weight: 500;
}

.stat--error {
  background: rgba(239, 68, 68, 0.15);
  color: #ef4444;
}

.stat--warning {
  background: rgba(245, 158, 11, 0.15);
  color: #f59e0b;
}
</style>
