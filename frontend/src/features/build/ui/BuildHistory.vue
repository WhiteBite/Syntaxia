<template>
  <div class="build-history">
    <div class="history-header">
      <h3 class="history-title">{{ t('build.history') }}</h3>
      <span class="history-count">{{ history.length }}</span>
    </div>

    <div v-if="history.length === 0" class="history-empty">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none">
        <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <p class="empty-text">{{ t('build.noHistory') }}</p>
    </div>

    <div v-else class="history-list">
      <BuildHistoryItem
        v-for="entry in history"
        :key="entry.id"
        :entry="entry"
        :is-selected="selectedId === entry.id"
        @select="$emit('select', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { BuildHistoryEntry } from '../types'
import BuildHistoryItem from './BuildHistoryItem.vue'

const { t } = useI18n()

defineProps<{
  history: BuildHistoryEntry[]
  selectedId: string | null
}>()

defineEmits<{
  (e: 'select', id: string | null): void
}>()
</script>

<style scoped>
.build-history {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.history-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.history-title {
  font-size: 13px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.history-count {
  padding: 2px 8px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  font-size: 11px;
  font-weight: 500;
  color: #9ca3af;
}

.history-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
  text-align: center;
}

.empty-icon {
  width: 40px;
  height: 40px;
  color: #4b5563;
  margin-bottom: 12px;
}

.empty-text {
  font-size: 13px;
  color: #6b7280;
  margin: 0;
}

.history-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
