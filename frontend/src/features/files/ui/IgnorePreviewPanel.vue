<template>
  <div class="preview-panel">
    <!-- Header with stats -->
    <div class="preview-panel__header">
      <div class="preview-panel__title">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        <span>{{ t('ignoreModal.preview') }}</span>
        <span v-if="result?.totalFiles" class="preview-panel__count">{{ formatNumber(result.totalFiles) }}</span>
      </div>
      <!-- View toggle -->
      <div v-if="result?.totalFiles" class="preview-panel__tabs">
        <button 
          :class="['preview-tab', viewMode === 'dirs' ? 'preview-tab--active' : '']"
          @click="viewMode = 'dirs'"
        >Папки</button>
        <button 
          :class="['preview-tab', viewMode === 'rules' ? 'preview-tab--active' : '']"
          @click="viewMode = 'rules'"
        >Правила</button>
        <button 
          :class="['preview-tab', viewMode === 'files' ? 'preview-tab--active' : '']"
          @click="viewMode = 'files'"
        >Файлы</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="preview-panel__loading">
      <svg class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
      </svg>
      <span>Анализ...</span>
    </div>

    <!-- Empty state -->
    <div v-else-if="!result?.totalFiles" class="preview-panel__empty">
      <svg class="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <p>{{ t('ignoreModal.previewEmpty') }}</p>
    </div>

    <!-- Content -->
    <div v-else class="preview-panel__content">
      <!-- Directories view -->
      <div v-if="viewMode === 'dirs'" class="preview-list">
        <div v-for="dir in result.topDirs" :key="dir.dir" class="preview-dir-item">
          <svg class="w-4 h-4 text-amber-400/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <span class="preview-dir-item__name">{{ dir.dir }}</span>
          <span class="preview-dir-item__count">{{ formatNumber(dir.count) }}</span>
        </div>
      </div>

      <!-- Rules view -->
      <div v-else-if="viewMode === 'rules'" class="preview-list">
        <div v-for="(count, rule) in topRules" :key="rule" class="preview-rule-item">
          <code class="preview-rule-item__pattern">{{ rule }}</code>
          <span class="preview-rule-item__count">{{ formatNumber(count) }}</span>
        </div>
      </div>

      <!-- Files view (virtualized) -->
      <div v-else class="preview-list">
        <div v-for="file in result.sampleFiles" :key="file" class="preview-file-item">
          <svg class="w-3.5 h-3.5 text-red-400/70 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
          </svg>
          <span class="preview-file-item__path">{{ file }}</span>
        </div>
        <div v-if="result.totalFiles > result.sampleFiles.length" class="preview-more">
          и ещё {{ formatNumber(result.totalFiles - result.sampleFiles.length) }} файлов...
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { IgnorePreviewResult } from '@/services/api/settings.api'
import { computed, ref } from 'vue'

const { t } = useI18n()

const props = defineProps<{
  result: IgnorePreviewResult | null
  loading?: boolean
}>()

const viewMode = ref<'dirs' | 'rules' | 'files'>('dirs')

const topRules = computed(() => {
  if (!props.result?.byRule) return {}
  const entries = Object.entries(props.result.byRule)
  entries.sort((a, b) => b[1] - a[1])
  return Object.fromEntries(entries.slice(0, 10))
})

function formatNumber(n: number): string {
  return n.toLocaleString('ru-RU')
}
</script>

<style scoped>
.preview-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: rgba(0, 0, 0, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 0.75rem;
  overflow: hidden;
}

.preview-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  background: rgba(0, 0, 0, 0.2);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  gap: 0.5rem;
  flex-wrap: wrap;
}

.preview-panel__title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: #94a3b8;
}

.preview-panel__count {
  padding: 0.125rem 0.5rem;
  background: rgba(239, 68, 68, 0.15);
  border-radius: 9999px;
  font-size: 0.6875rem;
  color: #f87171;
}

.preview-panel__tabs {
  display: flex;
  gap: 2px;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 0.375rem;
  padding: 2px;
}

.preview-tab {
  padding: 0.25rem 0.5rem;
  font-size: 0.6875rem;
  color: #64748b;
  border-radius: 0.25rem;
  transition: all 150ms;
}

.preview-tab:hover { color: #94a3b8; }
.preview-tab--active {
  background: rgba(139, 92, 246, 0.2);
  color: #c084fc;
}

.preview-panel__loading,
.preview-panel__empty {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  color: #475569;
  font-size: 0.75rem;
}

.preview-panel__content {
  flex: 1;
  overflow: hidden;
}

.preview-list {
  height: 100%;
  overflow-y: auto;
  padding: 0.5rem;
}

.preview-dir-item,
.preview-rule-item,
.preview-file-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.375rem 0.5rem;
  border-radius: 0.375rem;
  font-size: 0.75rem;
}

.preview-dir-item:hover,
.preview-rule-item:hover,
.preview-file-item:hover {
  background: rgba(255, 255, 255, 0.03);
}

.preview-dir-item__name,
.preview-file-item__path {
  flex: 1;
  color: #94a3b8;
  font-family: ui-monospace, monospace;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-dir-item__count,
.preview-rule-item__count {
  padding: 0.125rem 0.375rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 0.25rem;
  font-size: 0.625rem;
  color: #64748b;
}

.preview-rule-item__pattern {
  flex: 1;
  color: #c084fc;
  font-size: 0.6875rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-more {
  padding: 0.5rem;
  text-align: center;
  font-size: 0.6875rem;
  color: #64748b;
}
</style>
