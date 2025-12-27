<template>
  <div class="test-tree">
    <!-- Empty State -->
    <div v-if="suites.length === 0" class="tree-empty">
      <svg class="empty-icon" viewBox="0 0 24 24" fill="none">
        <path d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
      </svg>
      <p class="empty-text">{{ t('testing.noTests') }}</p>
      <button class="btn btn-ghost" @click="$emit('discover')">
        {{ t('testing.discoverTests') }}
      </button>
    </div>

    <!-- Test Suites -->
    <div v-else class="tree-content">
      <div 
        v-for="suite in suites" 
        :key="suite.id"
        class="suite-group"
      >
        <!-- Suite Header -->
        <button 
          class="suite-header"
          :class="[`suite-header--${suite.status}`]"
          @click="$emit('toggle-suite', suite.id)"
        >
          <svg 
            class="suite-chevron" 
            :class="{ 'suite-chevron--expanded': suite.expanded }"
            viewBox="0 0 24 24" 
            fill="none"
          >
            <path d="M9 5l7 7-7 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>

          <!-- Suite Status Icon -->
          <div class="suite-status">
            <svg v-if="suite.status === 'running'" class="status-icon animate-spin" viewBox="0 0 24 24" fill="none">
              <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2" opacity="0.25"/>
              <path d="M12 2a10 10 0 0110 10" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <svg v-else-if="suite.status === 'passed'" class="status-icon" viewBox="0 0 24 24" fill="none">
              <path d="M5 13l4 4L19 7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <svg v-else-if="suite.status === 'failed'" class="status-icon" viewBox="0 0 24 24" fill="none">
              <path d="M6 18L18 6M6 6l12 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
            <span v-else class="status-dot"></span>
          </div>

          <span class="suite-name">{{ suite.name }}</span>
          
          <div class="suite-badges">
            <span v-if="getPassedCount(suite)" class="suite-badge suite-badge--passed">
              {{ getPassedCount(suite) }}
            </span>
            <span v-if="getFailedCount(suite)" class="suite-badge suite-badge--failed">
              {{ getFailedCount(suite) }}
            </span>
            <span class="suite-total">{{ suite.tests.length }}</span>
          </div>
        </button>

        <!-- Suite Tests -->
        <div 
          v-if="suite.expanded"
          class="suite-tests"
        >
          <TestItem
            v-for="test in suite.tests"
            :key="test.id"
            :test="test"
            :is-selected="selectedTestId === test.id"
            :is-running="running"
            @select="$emit('select-test', $event)"
            @run="$emit('run-test', $event)"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { TestSuiteUI } from '../types'
import TestItem from './TestItem.vue'

const { t } = useI18n()

defineProps<{
  suites: TestSuiteUI[]
  selectedTestId: string | null
  running: boolean
}>()

defineEmits<{
  (e: 'toggle-suite', suiteId: string): void
  (e: 'select-test', testId: string): void
  (e: 'run-test', testId: string): void
  (e: 'discover'): void
}>()

function getPassedCount(suite: TestSuiteUI): number {
  return suite.tests.filter(t => t.status === 'passed').length
}

function getFailedCount(suite: TestSuiteUI): number {
  return suite.tests.filter(t => t.status === 'failed').length
}
</script>

<style scoped>
.test-tree {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.tree-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 32px;
  text-align: center;
}

.empty-icon {
  width: 48px;
  height: 48px;
  color: #4b5563;
  margin-bottom: 16px;
}

.empty-text {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 16px;
}

.tree-content {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.suite-group {
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  overflow: hidden;
}

.suite-header {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  background: none;
  border: none;
  cursor: pointer;
  transition: background 0.15s ease-out;
}

.suite-header:hover {
  background: rgba(255, 255, 255, 0.04);
}

.suite-chevron {
  width: 14px;
  height: 14px;
  color: #6b7280;
  transition: transform 0.15s ease-out;
  flex-shrink: 0;
}

.suite-chevron--expanded {
  transform: rotate(90deg);
}

.suite-status {
  flex-shrink: 0;
  width: 16px;
  height: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.suite-header--passed .suite-status {
  color: #22c55e;
}

.suite-header--failed .suite-status {
  color: #ef4444;
}

.suite-header--running .suite-status {
  color: #3b82f6;
}

.suite-header--pending .suite-status {
  color: #4b5563;
}

.status-icon {
  width: 14px;
  height: 14px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  opacity: 0.5;
}

.suite-name {
  flex: 1;
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  text-align: left;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.suite-badges {
  display: flex;
  align-items: center;
  gap: 6px;
}

.suite-badge {
  padding: 2px 6px;
  border-radius: 10px;
  font-size: 10px;
  font-weight: 600;
}

.suite-badge--passed {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

.suite-badge--failed {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

.suite-total {
  font-size: 11px;
  color: #6b7280;
}

.suite-tests {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
</style>
