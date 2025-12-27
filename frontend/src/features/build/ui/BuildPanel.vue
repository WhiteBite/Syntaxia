<template>
  <div class="build-panel">
    <!-- Header -->
    <div class="panel-header">
      <h2 class="panel-title">{{ t('build.title') }}</h2>
      
      <!-- Action Buttons -->
      <div class="panel-actions">
        <button
          class="btn btn-primary"
          :disabled="isBuilding"
          @click="build"
        >
          <svg v-if="isBuilding" class="btn-icon animate-spin" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.25"/>
            <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          <svg v-else class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" fill="currentColor"/>
            <path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2"/>
          </svg>
          {{ t('build.build') }}
        </button>

        <button
          class="btn btn-ghost"
          :disabled="isBuilding"
          @click="clean"
        >
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('build.clean') }}
        </button>

        <button
          class="btn btn-ghost"
          :disabled="isBuilding"
          @click="rebuild"
        >
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('build.rebuild') }}
        </button>

        <button
          v-if="isBuilding"
          class="btn btn-danger"
          :disabled="isCancelling"
          @click="cancelBuild"
        >
          <svg class="btn-icon" viewBox="0 0 24 24" fill="none">
            <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('build.cancel') }}
        </button>
      </div>
    </div>

    <!-- Status Bar -->
    <div class="status-bar" :class="[`status-bar--${status}`]">
      <div class="status-indicator">
        <span class="status-dot" :class="[`status-dot--${status}`]"></span>
        <span class="status-text">{{ statusText }}</span>
      </div>

      <div v-if="progress !== undefined && isBuilding" class="progress-container">
        <div class="progress-bar">
          <div class="progress-fill" :style="{ width: `${Math.min(progress, 100)}%` }"></div>
        </div>
        <span class="progress-text">{{ Math.round(progress) }}%</span>
      </div>

      <div v-if="lastBuildTime" class="last-build">
        <svg class="last-build-icon" viewBox="0 0 24 24" fill="none">
          <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span class="last-build-text">{{ lastBuildTime }}</span>
      </div>
    </div>

    <!-- Content -->
    <div class="panel-content">
      <div class="content-main">
        <BuildOutput
          :output="displayOutput"
          @copy="copyOutput"
          @clear="clearOutput"
        />
      </div>

      <div class="content-sidebar">
        <BuildHistory
          :history="history"
          :selected-id="selectedHistoryId"
          @select="selectHistoryEntry"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { storeToRefs } from 'pinia'
import { computed, onMounted } from 'vue'
import { useBuildStore } from '../model/build.store'
import BuildHistory from './BuildHistory.vue'
import BuildOutput from './BuildOutput.vue'

const { t } = useI18n()
const store = useBuildStore()

const {
  status,
  progress,
  currentStep,
  history,
  selectedHistoryId,
  isBuilding,
  isCancelling,
  lastBuildTime,
  displayOutput
} = storeToRefs(store)

const {
  build,
  clean,
  rebuild,
  cancelBuild,
  selectHistoryEntry,
  clearOutput,
  copyOutput,
  init
} = store

const statusText = computed(() => {
  if (currentStep.value) {
    return currentStep.value
  }

  switch (status.value) {
    case 'idle':
      return t('build.statusIdle')
    case 'building':
      return t('build.statusBuilding')
    case 'success':
      return t('build.statusSuccess')
    case 'failed':
      return t('build.statusFailed')
    case 'cancelled':
      return t('build.statusCancelled')
    default:
      return ''
  }
})

onMounted(() => {
  init()
})
</script>

<style scoped>
.build-panel {
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

.panel-actions {
  display: flex;
  gap: 8px;
}

/* Buttons */
.btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease-out;
  border: 1px solid transparent;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: rgba(139, 92, 246, 0.2);
  border-color: rgba(139, 92, 246, 0.4);
  color: #a78bfa;
}

.btn-primary:hover:not(:disabled) {
  background: rgba(139, 92, 246, 0.3);
  border-color: rgba(139, 92, 246, 0.6);
}

.btn-ghost {
  background: rgba(255, 255, 255, 0.04);
  border-color: rgba(255, 255, 255, 0.08);
  color: #9ca3af;
}

.btn-ghost:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.btn-danger {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.btn-danger:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.25);
  border-color: rgba(239, 68, 68, 0.5);
}

.btn-icon {
  width: 16px;
  height: 16px;
}

/* Status Bar */
.status-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #6b7280;
}

.status-dot--idle {
  background: #6b7280;
}

.status-dot--building {
  background: #3b82f6;
  animation: pulse 1.5s ease-in-out infinite;
}

.status-dot--success {
  background: #22c55e;
}

.status-dot--failed {
  background: #ef4444;
}

.status-dot--cancelled {
  background: #f59e0b;
}

.status-text {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
}

.progress-container {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  max-width: 30%;
}

.progress-bar {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #8b5cf6, #a78bfa);
  border-radius: 2px;
  transition: width 0.3s ease-out;
}

.progress-text {
  font-size: 11px;
  font-weight: 500;
  color: #9ca3af;
  min-width: 36px;
}

.last-build {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-left: auto;
}

.last-build-icon {
  width: 14px;
  height: 14px;
  color: #6b7280;
}

.last-build-text {
  font-size: 12px;
  color: #6b7280;
}

/* Content */
.panel-content {
  flex: 1;
  min-height: 0;
  display: flex;
  gap: 1px;
  background: rgba(255, 255, 255, 0.06);
}

.content-main {
  flex: 1;
  min-width: 0;
  background: var(--bg-primary, #0f1117);
}

.content-sidebar {
  width: 30%;
  min-width: 200px;
  max-width: min(300px, 35%);
  background: var(--bg-primary, #0f1117);
}

/* Animations */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}
</style>
