<template>
  <BaseModal :model-value="visible" @update:model-value="val => !val && $emit('close')">
    <template #header>
      <div>
        <h3 class="popup-title">{{ t('context.impactTitle') }}</h3>
        <p class="popup-subtitle">{{ t('context.impactSubtitle') }}</p>
      </div>
    </template>

    <div class="popup-body">
          <!-- Risk Score -->
          <div v-if="impactResult" class="risk-section">
            <div class="risk-header">
              <span class="risk-label">{{ t('context.riskScore') }}</span>
              <span class="risk-value" :class="getRiskClass(impactResult.riskLevel)">
                {{ getRiskLabel(impactResult.riskLevel) }}
              </span>
            </div>
            <div class="risk-bar-bg">
              <div 
                class="risk-bar" 
                :class="getRiskBarClass(impactResult.riskLevel)" 
                :style="{ width: `${Math.round(impactResult.aggregateRisk * 100)}%` }" 
              />
            </div>
          </div>

          <!-- Affected Files -->
          <div v-if="impactResult?.affectedFiles.length" class="impact-section">
            <div class="impact-section-header">
              {{ t('context.affectedFiles') }} ({{ impactResult.affectedFiles.length }})
            </div>
            <div class="popup-list">
              <div 
                v-for="file in impactResult.affectedFiles" 
                :key="file.path" 
                class="popup-item popup-item-readonly"
              >
                <span class="popup-item-icon">📄</span>
                <span class="popup-item-path">{{ file.path }}</span>
                <span 
                  class="popup-item-type" 
                  :class="file.type === 'direct' ? 'type-direct' : 'type-transitive'"
                >
                  {{ file.type === 'direct' ? t('context.directDep') : t('context.transitiveDep') }}
                </span>
              </div>
            </div>
          </div>

          <!-- Related Tests -->
          <div v-if="impactResult?.relatedTests.length" class="impact-section">
            <div class="impact-section-header">
              🧪 {{ t('context.relatedTests') }} ({{ impactResult.relatedTests.length }})
            </div>
            <div class="popup-list">
              <div 
                v-for="test in impactResult.relatedTests" 
                :key="test" 
                class="popup-item popup-item-readonly"
              >
                <span class="popup-item-icon">🧪</span>
                <span class="popup-item-path">{{ test }}</span>
              </div>
            </div>
          </div>

          <div v-if="!impactResult || impactResult.totalDependents === 0" class="popup-empty">
            {{ t('context.noImpactFiles') }}
          </div>
        </div>
  </BaseModal>
</template>

<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import { useI18n } from '@/composables/useI18n'
import type { ImpactResult } from '../composables/useAnalysisStatus'

defineProps<{
  visible: boolean
  impactResult: ImpactResult | null
  getRiskClass: (level: string) => string
  getRiskBarClass: (level: string) => string
  getRiskLabel: (level: string) => string
}>()

defineEmits<{
  (e: 'close'): void
}>()

const { t } = useI18n()
</script>

<style scoped>
.popup-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0;
}

.popup-subtitle {
  font-size: 0.75rem;
  color: var(--text-muted);
  margin: 0.25rem 0 0;
}

.popup-body {
  overflow-y: auto;
  padding: 0.75rem;
}

.popup-empty {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--text-muted);
  font-size: 0.875rem;
}

.popup-list {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.popup-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background 150ms;
}

.popup-item:hover {
  background: var(--bg-2);
}

.popup-item-readonly {
  cursor: default;
}

.popup-item-icon {
  flex-shrink: 0;
  font-size: 0.875rem;
}

.popup-item-path {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: 0.8125rem;
  color: var(--text-primary);
}

.popup-item-type {
  flex-shrink: 0;
  font-size: 0.6875rem;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
}

.type-direct {
  background: rgba(251, 191, 36, 0.15);
  color: #fbbf24;
}

.type-transitive {
  background: var(--bg-3);
  color: var(--text-muted);
}

/* Risk Section */
.risk-section {
  padding: 0.75rem;
  background: var(--bg-2);
  border-radius: 0.5rem;
  margin-bottom: 0.75rem;
}

.risk-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.risk-label {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.risk-value {
  font-size: 0.75rem;
  font-weight: 600;
}

.risk-low { color: #22c55e; }
.risk-medium { color: #f59e0b; }
.risk-high { color: #ef4444; }

.risk-bar-bg {
  height: 0.375rem;
  background: var(--bg-3);
  border-radius: 0.25rem;
  overflow: hidden;
}

.risk-bar {
  height: 100%;
  border-radius: 0.25rem;
  transition: width 300ms ease-out;
}

/* Impact Section */
.impact-section {
  margin-bottom: 0.75rem;
}

.impact-section-header {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-secondary);
  padding: 0.5rem 0.25rem;
}
</style>
