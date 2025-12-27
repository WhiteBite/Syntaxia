<template>
  <Transition name="slide-fade">
    <div v-if="visible && suggestions.length > 0" class="auto-suggest-panel">
      <!-- Header -->
      <div class="auto-suggest-header">
        <div class="auto-suggest-title">
          <svg class="w-4 h-4 text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
          </svg>
          <span>{{ t('chat.autoSuggest.title') }}</span>
          <span class="auto-suggest-count">({{ suggestions.length }})</span>
        </div>
        <button class="auto-suggest-hide" @click="$emit('hide')" :title="t('chat.autoSuggest.hide')">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- File List -->
      <div class="auto-suggest-list">
        <div
          v-for="item in suggestions"
          :key="item.path"
          class="auto-suggest-item"
          :class="{ 'is-selected': selectedPaths.has(item.path) }"
          @click="$emit('toggle', item.path)"
        >
          <!-- Checkbox -->
          <div class="auto-suggest-checkbox">
            <svg v-if="selectedPaths.has(item.path)" class="w-4 h-4 text-purple-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </div>

          <!-- File Info -->
          <div class="auto-suggest-file">
            <span class="auto-suggest-filename">{{ getFileName(item.path) }}</span>
            <span class="auto-suggest-filepath">{{ getFilePath(item.path) }}</span>
          </div>

          <!-- Relevance Score -->
          <div class="auto-suggest-relevance">
            <div class="auto-suggest-relevance-bar">
              <div 
                class="auto-suggest-relevance-fill" 
                :style="{ width: `${getRelevancePercent(item.confidence)}%` }"
              ></div>
            </div>
            <span class="auto-suggest-relevance-text">{{ getRelevancePercent(item.confidence) }}%</span>
          </div>

          <!-- Source Badge -->
          <span class="auto-suggest-badge" :class="getSourceBadgeClass(item.source)">
            {{ getSourceLabel(item.source) }}
          </span>
        </div>
      </div>

      <!-- Footer -->
      <div class="auto-suggest-footer">
        <button
          class="auto-suggest-add-btn"
          :disabled="selectedPaths.size === 0"
          @click="$emit('add')"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          {{ t('chat.autoSuggest.addSelected') }} ({{ selectedPaths.size }})
        </button>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { SmartSuggestion } from '@/services/api.service'

defineProps<{
  visible: boolean
  suggestions: SmartSuggestion[]
  selectedPaths: Set<string>
  getFileName: (path: string) => string
  getFilePath: (path: string) => string
  getSourceBadgeClass: (source: string) => string
  getSourceLabel: (source: string) => string
  getRelevancePercent: (confidence: number) => number
}>()

defineEmits<{
  (e: 'toggle', path: string): void
  (e: 'add'): void
  (e: 'hide'): void
}>()

const { t } = useI18n()
</script>

<style scoped>
.auto-suggest-panel {
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: 0.5rem;
  margin-bottom: 0.75rem;
  overflow: hidden;
}

.auto-suggest-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 0.75rem;
  background: var(--bg-1);
  border-bottom: 1px solid var(--border-subtle);
}

.auto-suggest-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-1);
}

.auto-suggest-count {
  color: var(--text-3);
  font-weight: 400;
}

.auto-suggest-hide {
  padding: 0.25rem;
  border-radius: 0.25rem;
  color: var(--text-3);
  transition: all 0.15s;
}

.auto-suggest-hide:hover {
  background: var(--bg-3);
  color: var(--text-1);
}

.auto-suggest-list {
  max-height: 12rem;
  overflow-y: auto;
}

.auto-suggest-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  cursor: pointer;
  transition: background 0.15s;
}

.auto-suggest-item:hover {
  background: var(--bg-3);
}

.auto-suggest-item.is-selected {
  background: rgba(139, 92, 246, 0.1);
}

.auto-suggest-checkbox {
  width: 1rem;
  height: 1rem;
  border: 1px solid var(--border-default);
  border-radius: 0.25rem;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.auto-suggest-item.is-selected .auto-suggest-checkbox {
  background: var(--purple-500);
  border-color: var(--purple-500);
}

.auto-suggest-file {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0.125rem;
}

.auto-suggest-filename {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--text-1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.auto-suggest-filepath {
  font-size: 0.6875rem;
  color: var(--text-3);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.auto-suggest-relevance {
  display: flex;
  align-items: center;
  gap: 0.375rem;
  flex-shrink: 0;
}

.auto-suggest-relevance-bar {
  width: 2.5rem;
  height: 0.25rem;
  background: var(--bg-3);
  border-radius: 0.125rem;
  overflow: hidden;
}

.auto-suggest-relevance-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--purple-500), var(--purple-400));
  border-radius: 0.125rem;
  transition: width 0.3s ease;
}

.auto-suggest-relevance-text {
  font-size: 0.6875rem;
  color: var(--text-3);
  min-width: 2rem;
  text-align: right;
}

.auto-suggest-badge {
  font-size: 0.625rem;
  padding: 0.125rem 0.375rem;
  border-radius: 0.25rem;
  font-weight: 500;
  text-transform: uppercase;
  flex-shrink: 0;
}

.badge-git {
  background: rgba(249, 115, 22, 0.15);
  color: #f97316;
}

.badge-arch {
  background: rgba(59, 130, 246, 0.15);
  color: #3b82f6;
}

.badge-semantic {
  background: rgba(16, 185, 129, 0.15);
  color: #10b981;
}

.badge-default {
  background: var(--bg-3);
  color: var(--text-3);
}

.auto-suggest-footer {
  padding: 0.5rem 0.75rem;
  border-top: 1px solid var(--border-subtle);
  background: var(--bg-1);
}

.auto-suggest-add-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  width: 100%;
  padding: 0.5rem 1rem;
  font-size: 0.8125rem;
  font-weight: 500;
  color: white;
  background: linear-gradient(135deg, var(--purple-500), var(--purple-600));
  border-radius: 0.375rem;
  transition: all 0.15s;
}

.auto-suggest-add-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, var(--purple-400), var(--purple-500));
  transform: translateY(-1px);
}

.auto-suggest-add-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Transition */
.slide-fade-enter-active {
  transition: all 0.2s ease-out;
}

.slide-fade-leave-active {
  transition: all 0.15s ease-in;
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-0.5rem);
}
</style>
