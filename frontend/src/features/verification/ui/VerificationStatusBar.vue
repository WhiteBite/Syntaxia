<template>
  <div class="status-bar">
    <!-- Status Icons -->
    <div class="status-icons">
      <!-- Build Status -->
      <div 
        class="status-item"
        :class="[`status-item--${status.build}`]"
        :title="`${t('verification.build')}: ${t(`verification.${status.build}`)}`"
      >
        <BaseSpinner v-if="status.build === 'running'" size="sm" class="status-icon" />
        <svg v-else-if="status.build === 'passed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else-if="status.build === 'failed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <svg v-else class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        </svg>
        <span class="status-label">{{ t('verification.build') }}</span>
      </div>

      <!-- Lint Status -->
      <div 
        class="status-item"
        :class="[`status-item--${status.lint}`]"
        :title="`${t('verification.lint')}: ${t(`verification.${status.lint}`)}`"
      >
        <BaseSpinner v-if="status.lint === 'running'" size="sm" class="status-icon" />
        <svg v-else-if="status.lint === 'passed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else-if="status.lint === 'failed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <svg v-else class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        </svg>
        <span class="status-label">{{ t('verification.lint') }}</span>
      </div>

      <!-- Test Status -->
      <div 
        class="status-item"
        :class="[`status-item--${status.test}`]"
        :title="`${t('verification.tests')}: ${t(`verification.${status.test}`)}`"
      >
        <BaseSpinner v-if="status.test === 'running'" size="sm" class="status-icon" />
        <svg v-else-if="status.test === 'passed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else-if="status.test === 'failed'" class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
          <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
        </svg>
        <svg v-else class="status-icon" viewBox="0 0 24 24" fill="none">
          <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
        </svg>
        <span class="status-label">{{ t('verification.tests') }}</span>
      </div>
    </div>

    <!-- Divider -->
    <div class="status-divider"></div>

    <!-- Current Action / Counters -->
    <div class="status-info">
      <!-- Running Status -->
      <span v-if="currentStatusText" class="status-running">
        {{ currentStatusText }}
      </span>
      
      <!-- Counters -->
      <template v-else>
        <span v-if="errorCount > 0" class="status-counter status-counter--error">
          <svg class="counter-icon" viewBox="0 0 24 24" fill="none">
            <circle cx="12" cy="12" r="10" stroke="currentColor" stroke-width="2"/>
            <path d="M15 9l-6 6M9 9l6 6" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          </svg>
          {{ errorCount }} {{ t('verification.errors') }}
        </span>
        
        <span v-if="warningCount > 0" class="status-counter status-counter--warning">
          <svg class="counter-icon" viewBox="0 0 24 24" fill="none">
            <path d="M12 9v4m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ warningCount }} {{ t('verification.warnings') }}
        </span>
        
        <span v-if="errorCount === 0 && warningCount === 0 && allPassed" class="status-counter status-counter--success">
          <svg class="counter-icon" viewBox="0 0 24 24" fill="none">
            <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ t('verification.allPassed') }}
        </span>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseSpinner } from '@/components/ui'
import { useI18n } from '@/composables/useI18n'
import type { VerificationStatus } from '../model/verification.store'

const { t } = useI18n()

defineProps<{
  status: VerificationStatus
  errorCount: number
  warningCount: number
  currentStatusText: string
  allPassed: boolean
}>()
</script>

<style scoped>
.status-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

/* Status Icons */
.status-icons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  transition: all 0.15s ease-out;
}

.status-icon {
  width: 16px;
  height: 16px;
}

.status-label {
  font-size: 12px;
  font-weight: 500;
  color: #9ca3af;
}

/* Status variants */
.status-item--idle {
  color: #6b7280;
}

.status-item--running {
  color: #3b82f6;
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.2);
}

.status-item--passed {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.1);
  border-color: rgba(34, 197, 94, 0.2);
}

.status-item--failed {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

/* Divider */
.status-divider {
  width: 1px;
  height: 24px;
  background: rgba(255, 255, 255, 0.1);
}

/* Info */
.status-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
}

.status-running {
  font-size: 12px;
  color: #3b82f6;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.status-counter {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 500;
}

.counter-icon {
  width: 14px;
  height: 14px;
}

.status-counter--error {
  color: #f87171;
}

.status-counter--warning {
  color: #fbbf24;
}

.status-counter--success {
  color: #4ade80;
}

</style>
