<template>
  <BaseModal :model-value="visible" @update:model-value="val => !val && $emit('close')">
    <template #header>
      <div class="smart-hud-header-content">
        <div class="smart-hud-title-row">
          <div class="smart-hud-sparkle-icon">✨</div>
          <div class="smart-hud-title-block">
            <span class="smart-hud-title">{{ t('context.aiRecommendations') }}</span>
            <span class="smart-hud-subtitle-inline">{{ t('context.foundFilesAnalysis', { count: suggestions.length }) }}</span>
          </div>
        </div>
      </div>
    </template>

    <!-- File List -->
    <div class="smart-hud-body">
      <div v-if="suggestions.length === 0" class="smart-hud-empty">
        {{ t('context.noRelatedFiles') }}
      </div>
      <div v-else class="smart-hud-list">
        <SmartSuggestionItem
          v-for="item in suggestions"
          :key="item.path"
          :is-selected="selectedPaths.has(item.path)"
          :icon-class="getFileIconClass(item.path)"
          :file-name="getFileName(item.path)"
          :file-path="getFilePath(item.path)"
          :badge-class="getSourceBadgeClass(item.source)"
          :source-label="getSourceLabel(item.source)"
          @toggle="$emit('toggle', item.path)"
        />
      </div>
    </div>

    <template #footer>
      <!-- Premium Footer -->
      <div v-if="suggestions.length > 0" class="smart-hud-footer-content">
        <button 
          class="smart-hud-action group"
          :disabled="selectedPaths.size === 0"
          @click="$emit('add')"
        >
          <svg class="smart-hud-action-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          <span>{{ t('context.addFiles') }} ({{ selectedPaths.size }})</span>
          <div class="smart-hud-shimmer"></div>
        </button>
      </div>
    </template>
  </BaseModal>
</template>

<script setup lang="ts">
import BaseModal from '@/components/ui/BaseModal.vue'
import { useI18n } from '@/composables/useI18n'
import type { SmartSuggestion } from '@/services/api.service'
import SmartSuggestionItem from './SmartSuggestionItem.vue'

defineProps<{
  visible: boolean
  suggestions: SmartSuggestion[]
  selectedPaths: Set<string>
  getFileIconClass: (path: string) => string
  getFileName: (path: string) => string
  getFilePath: (path: string) => string
  getSourceBadgeClass: (source: string) => string
  getSourceLabel: (source: string) => string
}>()

defineEmits<{
  (e: 'close'): void
  (e: 'toggle', path: string): void
  (e: 'add'): void
}>()

const { t } = useI18n()
</script>

<style scoped>
.smart-hud-header-content {
  display: flex;
  align-items: center;
  width: 100%;
}

.smart-hud-title-row {
  display: flex;
  align-items: center;
  gap: 0.625rem;
}

.smart-hud-sparkle-icon {
  width: 1.75rem;
  height: 1.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.15), rgba(249, 115, 22, 0.15));
  border-radius: 0.5rem;
  border: 1px solid rgba(251, 191, 36, 0.2);
}

.smart-hud-title-block {
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.smart-hud-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: #f3f4f6;
  line-height: 1.2;
}

.smart-hud-subtitle-inline {
  font-size: 0.6875rem;
  color: #6b7280;
  line-height: 1.2;
}

.smart-hud-body {
  overflow-y: auto;
  padding: 0.25rem;
}

.smart-hud-body::-webkit-scrollbar {
  width: 6px;
}

.smart-hud-body::-webkit-scrollbar-track {
  background: transparent;
}

.smart-hud-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
}

.smart-hud-body::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.12);
}

.smart-hud-empty {
  text-align: center;
  padding: 2rem 1rem;
  color: #4b5563;
  font-size: 0.8125rem;
}

.smart-hud-list {
  display: flex;
  flex-direction: column;
  padding: 0.25rem 0.5rem;
}

.smart-hud-footer-content {
  width: 100%;
}

.smart-hud-action {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.875rem 1.25rem;
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  border: none;
  border-radius: 0.75rem;
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.15),
    0 4px 16px rgba(139, 92, 246, 0.4);
  cursor: pointer;
  transition: all 150ms;
  overflow: hidden;
}

.smart-hud-action-icon {
  color: rgba(255, 255, 255, 0.9);
}

.smart-hud-action:hover:not(:disabled) {
  background: linear-gradient(135deg, #9d6ff8 0%, #7577f5 100%);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.2),
    0 8px 24px rgba(139, 92, 246, 0.5);
  transform: translateY(-1px);
}

.smart-hud-action:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 2px 8px rgba(139, 92, 246, 0.3);
}

.smart-hud-action:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  box-shadow: none;
}

.smart-hud-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(90deg, transparent 0%, rgba(255, 255, 255, 0.15) 50%, transparent 100%);
  transform: translateX(-100%) skewX(-15deg);
  pointer-events: none;
}

.smart-hud-action:hover:not(:disabled) .smart-hud-shimmer {
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%) skewX(-15deg); }
  100% { transform: translateX(200%) skewX(-15deg); }
}
</style>
