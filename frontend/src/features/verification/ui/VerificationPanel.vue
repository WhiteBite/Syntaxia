<template>
  <div class="verification-panel">
    <!-- Header -->
    <div class="panel-header">
      <h2 class="panel-title">{{ t('verification.title') }}</h2>
      
      <!-- Tabs -->
      <div class="panel-tabs">
        <BaseButton
          v-for="tab in tabs"
          :key="tab.type"
          variant="ghost"
          size="sm"
          :class="{ 
            'tab-btn--active': activeTab === tab.type,
            [`tab-btn--${status[tab.type]}`]: true
          }"
          @click="setActiveTab(tab.type)"
        >
          <span class="tab-label">{{ t(`verification.${tab.type === 'test' ? 'tests' : tab.type}`) }}</span>
          <span class="tab-status">
            <BaseSpinner v-if="status[tab.type] === 'running'" size="sm" class="tab-icon" />
            <svg v-else-if="status[tab.type] === 'passed'" class="tab-icon" viewBox="0 0 24 24" fill="none">
              <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <svg v-else-if="status[tab.type] === 'failed'" class="tab-icon" viewBox="0 0 24 24" fill="none">
              <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span v-else class="tab-dot"></span>
          </span>
        </BaseButton>
      </div>

      <!-- Run All Button -->
      <BaseButton
        variant="primary"
        size="sm"
        :disabled="isRunning"
        :loading="isRunning"
        @click="runAll"
      >
        <template #icon>
          <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
            <path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" fill="currentColor"/>
            <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2"/>
          </svg>
        </template>
        {{ t('verification.runAll') }}
      </BaseButton>
    </div>

    <!-- Status Bar -->
    <VerificationStatusBar
      :status="status"
      :error-count="errorCount"
      :warning-count="warningCount"
      :current-status-text="currentStatusText"
      :all-passed="allPassed"
    />

    <!-- Toolbar -->
    <div class="panel-toolbar">
      <div class="toolbar-left">
        <BaseButton
          variant="ghost"
          size="sm"
          icon-only
          :class="{ 'toolbar-btn--active': groupBy === 'file' }"
          @click="groupBy = 'file'"
          :title="t('verification.groupByFile')"
        >
          <template #icon>
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
              <path d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </template>
        </BaseButton>
        <BaseButton
          variant="ghost"
          size="sm"
          icon-only
          :class="{ 'toolbar-btn--active': groupBy === 'severity' }"
          @click="groupBy = 'severity'"
          :title="t('verification.groupBySeverity')"
        >
          <template #icon>
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
              <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </template>
        </BaseButton>
      </div>

      <div class="toolbar-right">
        <BaseButton
          variant="ghost"
          size="sm"
          icon-only
          @click="errorListRef?.expandAll()"
          :title="t('verification.expandAll')"
        >
          <template #icon>
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
              <path d="M19 9l-7 7-7-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </template>
        </BaseButton>
        <BaseButton
          variant="ghost"
          size="sm"
          icon-only
          @click="errorListRef?.collapseAll()"
          :title="t('verification.collapseAll')"
        >
          <template #icon>
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
              <path d="M5 15l7-7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </template>
        </BaseButton>
        <BaseButton
          variant="ghost"
          size="sm"
          icon-only
          @click="clearResults"
          :title="t('verification.clearResults')"
        >
          <template #icon>
            <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none">
              <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </template>
        </BaseButton>
      </div>
    </div>

    <!-- Content -->
    <div class="panel-content">
      <ErrorList
        ref="errorListRef"
        :errors="filteredErrors"
        :group-by="groupBy"
        :is-fixing="isFixing"
        @go-to-file="handleGoToFile"
        @fix-with-ai="fixWithAI"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'
import { storeToRefs } from 'pinia'
import { computed, ref } from 'vue'
import type { VerificationError, VerificationType } from '../api/verification.api'
import { useVerificationStore } from '../model/verification.store'
import ErrorList from './ErrorList.vue'
import VerificationStatusBar from './VerificationStatusBar.vue'

const { t } = useI18n()
const store = useVerificationStore()

const {
  status,
  errors,
  activeTab,
  isFixing,
  errorCount,
  warningCount,
  isRunning,
  allPassed,
  currentStatusText
} = storeToRefs(store)

const { runAll, fixWithAI, setActiveTab, clearResults } = store

// Local state
const groupBy = ref<'file' | 'severity'>('file')
const errorListRef = ref<InstanceType<typeof ErrorList> | null>(null)

// Tabs config
const tabs: Array<{ type: VerificationType }> = [
  { type: 'build' },
  { type: 'lint' },
  { type: 'test' }
]

// Filter errors by active tab
const filteredErrors = computed(() => {
  return errors.value.filter(e => e.source === activeTab.value)
})

// Handle navigation to file
function handleGoToFile(error: VerificationError): void {
  window.dispatchEvent(new CustomEvent('verification:go-to-file', {
    detail: { file: error.file, line: error.line, column: error.column }
  }))
}
</script>

<style scoped>
.verification-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: var(--bg-primary, #0f1117);
  border-radius: 12px;
  overflow: hidden;
}

/* Header */
.panel-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

/* Tabs */
.panel-tabs {
  display: flex;
  gap: 4px;
  flex: 1;
}

.tab-btn--active {
  background: rgba(139, 92, 246, 0.15) !important;
  border-color: rgba(139, 92, 246, 0.3) !important;
  color: #a78bfa !important;
}

.tab-btn--passed .tab-status {
  color: #22c55e;
}

.tab-btn--failed .tab-status {
  color: #ef4444;
}

.tab-btn--running .tab-status {
  color: #3b82f6;
}

.tab-label {
  font-size: 12px;
  font-weight: 500;
}

.tab-status {
  display: flex;
  align-items: center;
}

.tab-icon {
  width: 14px;
  height: 14px;
}

.tab-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.5;
}

/* Toolbar */
.panel-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.toolbar-left,
.toolbar-right {
  display: flex;
  gap: 4px;
}

.toolbar-btn--active {
  background: rgba(139, 92, 246, 0.1) !important;
  border-color: rgba(139, 92, 246, 0.2) !important;
  color: #a78bfa !important;
}

/* Content */
.panel-content {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 16px;
}
</style>
