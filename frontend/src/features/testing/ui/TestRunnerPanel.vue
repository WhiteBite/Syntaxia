<template>
  <div class="test-runner-panel">
    <!-- Header -->
    <div class="panel-header">
      <h2 class="panel-title">{{ t('testing.title') }}</h2>

      <!-- Action Buttons -->
      <div class="header-actions">
        <button
          class="btn btn-primary"
          :disabled="running"
          @click="runAllTests"
        >
          <svg v-if="running" class="btn-icon animate-spin" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.25"/>
            <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M5 3l14 9-14 9V3z" fill="currentColor"/>
          </svg>
          {{ t('testing.runAll') }}
        </button>

        <button
          class="btn btn-ghost"
          :disabled="running || !hasFailedTests"
          @click="runFailedTests"
        >
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('testing.runFailed') }}
        </button>

        <button
          v-if="running"
          class="btn btn-ghost btn-danger"
          @click="stopTests"
        >
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none">
            <rect x="6" y="6" width="12" height="12" rx="2" fill="currentColor"/>
          </svg>
          {{ t('testing.stop') }}
        </button>
      </div>
    </div>

    <!-- Stats Bar -->
    <div class="stats-bar">
      <div class="stats-items">
        <div class="stat-item stat-item--passed">
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none">
            <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="stat-value">{{ passedCount }}</span>
          <span class="stat-label">{{ t('testing.passed') }}</span>
        </div>

        <div class="stat-item stat-item--failed">
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none">
            <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          <span class="stat-value">{{ failedCount }}</span>
          <span class="stat-label">{{ t('testing.failed') }}</span>
        </div>

        <div class="stat-item stat-item--skipped">
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none">
            <path d="M5 5l14 14M5 19L19 5" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" opacity="0.5"/>
          </svg>
          <span class="stat-value">{{ skippedCount }}</span>
          <span class="stat-label">{{ t('testing.skipped') }}</span>
        </div>
      </div>

      <!-- Progress Bar -->
      <div v-if="running" class="progress-container">
        <div class="progress-bar">
          <div 
            class="progress-fill"
            :style="{ width: `${progressPercent}%` }"
          ></div>
        </div>
        <span class="progress-text">{{ progressPercent }}%</span>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="filter-tabs">
      <button
        v-for="filterOption in filterOptions"
        :key="filterOption.value"
        class="filter-tab"
        :class="{ 'filter-tab--active': filter === filterOption.value }"
        @click="setFilter(filterOption.value)"
      >
        {{ filterOption.label }}
        <span v-if="filterOption.count > 0" class="filter-count">
          {{ filterOption.count }}
        </span>
      </button>

      <div class="filter-spacer"></div>

      <!-- Expand/Collapse -->
      <button
        class="toolbar-btn"
        :title="t('testing.expandAll')"
        @click="expandAllSuites"
      >
        <svg class="toolbar-icon" viewBox="0 0 24 24" fill="none">
          <path d="M19 9l-7 7-7-7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <button
        class="toolbar-btn"
        :title="t('testing.collapseAll')"
        @click="collapseAllSuites"
      >
        <svg class="toolbar-icon" viewBox="0 0 24 24" fill="none">
          <path d="M5 15l7-7 7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
      <button
        class="toolbar-btn"
        :title="t('testing.clearResults')"
        @click="clearResults"
      >
        <svg class="toolbar-icon" viewBox="0 0 24 24" fill="none">
          <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </button>
    </div>

    <!-- Content -->
    <div class="panel-content">
      <div class="content-split">
        <!-- Test Tree -->
        <div class="tree-section">
          <TestTree
            :suites="filteredSuites"
            :selected-test-id="selectedTestId"
            :running="running"
            @toggle-suite="toggleSuite"
            @select-test="selectTest"
            @run-test="runTest"
            @discover="discoverTests"
          />
        </div>

        <!-- Test Output -->
        <div class="output-section">
          <TestOutput
            :output="selectedTest?.output"
            :error="selectedTest?.error"
            :stack-trace="selectedTest?.stackTrace"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { storeToRefs } from 'pinia'
import { computed } from 'vue'
import { useTestingStore } from '../model/testing.store'
import type { TestFilter } from '../types'
import TestOutput from './TestOutput.vue'
import TestTree from './TestTree.vue'

const { t } = useI18n()
const store = useTestingStore()

const {
  running,
  filter,
  selectedTestId,
  selectedTest,
  filteredSuites,
  passedCount,
  failedCount,
  skippedCount,
  hasFailedTests,
  progressPercent,
  stats
} = storeToRefs(store)

const {
  runAllTests,
  runFailedTests,
  runTest,
  stopTests,
  discoverTests,
  setFilter,
  selectTest,
  toggleSuite,
  expandAllSuites,
  collapseAllSuites,
  clearResults
} = store

const filterOptions = computed(() => [
  { value: 'all' as TestFilter, label: t('testing.all'), count: stats.value.total },
  { value: 'passed' as TestFilter, label: t('testing.passed'), count: stats.value.passed },
  { value: 'failed' as TestFilter, label: t('testing.failed'), count: stats.value.failed },
  { value: 'skipped' as TestFilter, label: t('testing.skipped'), count: stats.value.skipped }
])
</script>

<style scoped>
.test-runner-panel {
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
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #e5e7eb;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.btn-primary {
  background: rgba(139, 92, 246, 0.9);
  border: 1px solid rgba(139, 92, 246, 0.5);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: rgba(139, 92, 246, 1);
}

.btn-ghost {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #d1d5db;
}

.btn-ghost:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.btn-danger {
  color: #f87171;
  border-color: rgba(239, 68, 68, 0.3);
}

.btn-danger:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.15);
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-icon {
  width: 14px;
  height: 14px;
}

/* Stats Bar */
.stats-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.stats-items {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-icon {
  width: 14px;
  height: 14px;
}

.stat-item--passed { color: #22c55e; }
.stat-item--failed { color: #ef4444; }
.stat-item--skipped { color: #6b7280; }

.stat-value {
  font-size: 16px;
  font-weight: 600;
}

.stat-label {
  font-size: 12px;
  opacity: 0.7;
}

.progress-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-bar {
  width: 120px;
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b5cf6, #a78bfa);
  border-radius: 3px;
  transition: width 0.3s ease-out;
}

.progress-text {
  font-size: 12px;
  color: #9ca3af;
  font-family: ui-monospace, monospace;
}

/* Filter Tabs */
.filter-tabs {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 8px 16px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.filter-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: none;
  border: 1px solid transparent;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.filter-tab:hover {
  background: rgba(255, 255, 255, 0.04);
  color: #e5e7eb;
}

.filter-tab--active {
  background: rgba(139, 92, 246, 0.15);
  border-color: rgba(139, 92, 246, 0.3);
  color: #a78bfa;
}

.filter-count {
  padding: 1px 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  font-size: 10px;
}

.filter-tab--active .filter-count {
  background: rgba(139, 92, 246, 0.3);
}

.filter-spacer {
  flex: 1;
}

.toolbar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: none;
  border: 1px solid transparent;
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.toolbar-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #e5e7eb;
}

.toolbar-icon {
  width: 16px;
  height: 16px;
}

/* Content */
.panel-content {
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.content-split {
  display: flex;
  height: 100%;
}

.tree-section {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 12px;
  border-right: 1px solid rgba(255, 255, 255, 0.06);
}

.output-section {
  width: 40%;
  min-width: 250px;
  max-width: min(500px, 45%);
}

/* Animation */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
