<template>
  <div v-if="shouldShowBar" class="analysis-status-bar">
    <!-- Loading State -->
    <div v-if="isLoadingRelated" class="status-btn status-btn-loading">
      <span class="status-spinner"></span>
      <span class="status-label">{{ t('context.analyzingFiles') }}</span>
    </div>

    <!-- Related Files Button -->
    <button 
      v-else-if="relatedCount > 0"
      class="status-btn status-btn-related"
      :title="t('context.relatedFilesTooltip')"
      @click="showRelatedPopup = true"
    >
      <span class="status-icon">💡</span>
      <span class="status-count">{{ relatedCount }}</span>
      <span class="status-label">{{ t('context.relatedFiles') }}</span>
    </button>

    <!-- Dependent Files Button -->
    <button 
      v-if="dependentCount > 0"
      class="status-btn status-btn-impact"
      :title="t('context.dependentFilesTooltip')"
      @click="showImpactPopup = true"
    >
      <span class="status-icon">✨</span>
      <span class="status-count">{{ dependentCount }}</span>
      <span class="status-label">{{ t('context.recommendations') }}</span>
    </button>

    <!-- Smart Suggestions HUD -->
    <SmartSuggestionsHud
      :visible="showRelatedPopup"
      :suggestions="suggestions"
      :selected-paths="selectedRelated"
      :get-file-icon-class="getFileIconClass"
      :get-file-name="getFileName"
      :get-file-path="getFilePath"
      :get-source-badge-class="getSourceBadgeClass"
      :get-source-label="getSourceLabel"
      @close="showRelatedPopup = false"
      @toggle="toggleRelated"
      @add="addSelectedRelated"
    />

    <!-- Impact Analysis Popup -->
    <ImpactAnalysisPopup
      :visible="showImpactPopup"
      :impact-result="impactResult"
      :get-risk-class="getRiskClass"
      :get-risk-bar-class="getRiskBarClass"
      :get-risk-label="getRiskLabel"
      @close="showImpactPopup = false"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { toRef } from 'vue'
import { useAnalysisStatus } from '../composables/useAnalysisStatus'
import ImpactAnalysisPopup from './ImpactAnalysisPopup.vue'
import SmartSuggestionsHud from './SmartSuggestionsHud.vue'

const props = defineProps<{
  selectedFiles: string[]
}>()

const emit = defineEmits<{
  (e: 'add-files', files: string[]): void
}>()

const { t } = useI18n()

const {
  // State
  suggestions,
  impactResult,
  selectedRelated,
  isLoadingRelated,
  showRelatedPopup,
  showImpactPopup,

  // Computed
  shouldShowBar,
  relatedCount,
  dependentCount,

  // Actions
  toggleRelated,
  addSelectedRelated,

  // Helpers
  getSourceLabel,
  getSourceBadgeClass,
  getFileIconClass,
  getFileName,
  getFilePath,
  getRiskClass,
  getRiskBarClass,
  getRiskLabel,
} = useAnalysisStatus({
  selectedFiles: toRef(props, 'selectedFiles'),
  onAddFiles: (files) => emit('add-files', files),
})
</script>

<style scoped>
.analysis-status-bar {
  display: flex;
  gap: 0.5rem;
  padding: 0.5rem;
  border-top: 1px solid var(--border-default);
  background: var(--bg-2);
}

.status-btn {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.625rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
  transition: all 150ms ease-out;
  cursor: pointer;
  border: 1px solid transparent;
}

.status-btn-related {
  background: rgba(251, 191, 36, 0.1);
  color: #fbbf24;
  border-color: rgba(251, 191, 36, 0.2);
}

.status-btn-related:hover {
  background: rgba(251, 191, 36, 0.2);
  border-color: rgba(251, 191, 36, 0.3);
}

.status-btn-impact {
  background: rgba(139, 92, 246, 0.1);
  color: #a78bfa;
  border-color: rgba(139, 92, 246, 0.2);
}

.status-btn-impact:hover {
  background: rgba(139, 92, 246, 0.2);
  border-color: rgba(139, 92, 246, 0.3);
}

.status-btn-loading {
  opacity: 0.7;
  cursor: wait;
}

.status-icon {
  font-size: 0.875rem;
}

.status-count {
  font-weight: 600;
}

.status-label {
  color: var(--text-muted);
}

.status-spinner {
  width: 0.875rem;
  height: 0.875rem;
  border: 2px solid currentColor;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
