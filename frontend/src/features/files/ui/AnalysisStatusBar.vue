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

    <!-- Dependent Files Button - only show if there are dependents -->
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

    <!-- Smart Suggestions HUD (Raycast Style) -->
    <Teleport to="body">
      <div v-if="showRelatedPopup" class="smart-hud-overlay" @click.self="showRelatedPopup = false">
        <div class="smart-hud">
          <!-- Premium Header -->
          <div class="smart-hud-header">
            <div class="smart-hud-title-row">
              <div class="smart-hud-sparkle-icon">✨</div>
              <div class="smart-hud-title-block">
                <span class="smart-hud-title">{{ t('context.aiRecommendations') }}</span>
                <span class="smart-hud-subtitle-inline">{{ t('context.foundFilesAnalysis', { count: suggestions.length }) }}</span>
              </div>
            </div>
            <button class="smart-hud-close" @click="showRelatedPopup = false">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M18 6L6 18M6 6l12 12"/>
              </svg>
            </button>
          </div>

          <!-- File List -->
          <div class="smart-hud-body">
            <div v-if="suggestions.length === 0" class="smart-hud-empty">
              {{ t('context.noRelatedFiles') }}
            </div>
            <div v-else class="smart-hud-list">
              <label 
                v-for="item in suggestions" 
                :key="item.path"
                class="smart-hud-item"
                :class="{ 'smart-hud-item-selected': selectedRelated.has(item.path) }"
              >
                <!-- Checkbox -->
                <div class="smart-hud-checkbox" :class="{ checked: selectedRelated.has(item.path) }">
                  <input 
                    type="checkbox" 
                    :checked="selectedRelated.has(item.path)"
                    @change="toggleRelated(item.path)"
                  />
                  <svg v-if="selectedRelated.has(item.path)" viewBox="0 0 12 12" fill="none">
                    <path d="M2 6L5 9L10 3" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </div>

                <!-- File Icon (SVG) -->
                <div class="smart-hud-file-icon" :class="getFileIconClass(item.path)">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8l-6-6z"/>
                    <path d="M14 2v6h6M10 12l-2 2 2 2M14 12l2 2-2 2"/>
                  </svg>
                </div>

                <!-- File Info -->
                <div class="smart-hud-file-info">
                  <span class="smart-hud-file-name">{{ getFileName(item.path) }}</span>
                  <span class="smart-hud-file-path">{{ getFilePath(item.path) }}</span>
                </div>

                <!-- Source Badge (tiny chip with icon) -->
                <span class="smart-hud-badge" :class="getSourceBadgeClass(item.source)">
                  <span class="smart-hud-badge-icon">🔗</span>
                  {{ getSourceLabel(item.source) }}
                </span>
              </label>
            </div>
          </div>

          <!-- Premium Footer with separated button -->
          <div v-if="suggestions.length > 0" class="smart-hud-footer">
            <button 
              class="smart-hud-action group"
              :disabled="selectedRelated.size === 0"
              @click="addSelectedRelated"
            >
              <svg class="smart-hud-action-icon" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <path d="M12 5v14M5 12h14"/>
              </svg>
              <span>{{ t('context.addFiles') }} ({{ selectedRelated.size }})</span>
              <!-- Shimmer effect -->
              <div class="smart-hud-shimmer"></div>
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Impact Analysis Popup -->
    <Teleport to="body">
      <div v-if="showImpactPopup" class="popup-overlay" @click.self="showImpactPopup = false">
        <div class="popup-content popup-impact">
          <div class="popup-header">
            <div>
              <h3 class="popup-title">{{ t('context.impactTitle') }}</h3>
              <p class="popup-subtitle">{{ t('context.impactSubtitle') }}</p>
            </div>
            <button class="popup-close" @click="showImpactPopup = false">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

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
                <div class="risk-bar" :class="getRiskBarClass(impactResult.riskLevel)" 
                  :style="{ width: `${Math.round(impactResult.aggregateRisk * 100)}%` }" />
              </div>
            </div>

            <!-- Affected Files -->
            <div v-if="impactResult?.affectedFiles.length" class="impact-section">
              <div class="impact-section-header">
                {{ t('context.affectedFiles') }} ({{ impactResult.affectedFiles.length }})
              </div>
              <div class="popup-list">
                <div v-for="file in impactResult.affectedFiles" :key="file.path" class="popup-item popup-item-readonly">
                  <span class="popup-item-icon">📄</span>
                  <span class="popup-item-path">{{ file.path }}</span>
                  <span class="popup-item-type" :class="file.type === 'direct' ? 'type-direct' : 'type-transitive'">
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
                <div v-for="test in impactResult.relatedTests" :key="test" class="popup-item popup-item-readonly">
                  <span class="popup-item-icon">🧪</span>
                  <span class="popup-item-path">{{ test }}</span>
                </div>
              </div>
            </div>

            <div v-if="!impactResult || impactResult.totalDependents === 0" class="popup-empty">
              {{ t('context.noImpactFiles') }}
            </div>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { apiService, type SmartSuggestion } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { computed, ref, watch } from 'vue'

interface AffectedFile {
  path: string
  type: 'direct' | 'transitive'
  dependents: number
}

interface ImpactResult {
  totalDependents: number
  aggregateRisk: number
  riskLevel: 'low' | 'medium' | 'high'
  affectedFiles: AffectedFile[]
  relatedTests: string[]
}

const props = defineProps<{
  selectedFiles: string[]
}>()

const emit = defineEmits<{
  (e: 'add-files', files: string[]): void
}>()

const { t } = useI18n()
const projectStore = useProjectStore()

// Constants
const FETCH_TIMEOUT_MS = 5000

// State
const suggestions = ref<SmartSuggestion[]>([])
const impactResult = ref<ImpactResult | null>(null)
const selectedRelated = ref<Set<string>>(new Set())
const isLoadingRelated = ref(false)
const showRelatedPopup = ref(false)
const showImpactPopup = ref(false)

// Computed
const projectPath = computed(() => projectStore.currentPath || '')
const hasSelectedFiles = computed(() => props.selectedFiles.length > 0)
const relatedCount = computed(() => suggestions.value.length)
const dependentCount = computed(() => impactResult.value?.totalDependents || 0)

// Show bar when: loading, has related files, or has dependents
const shouldShowBar = computed(() => 
  hasSelectedFiles.value && (isLoadingRelated.value || relatedCount.value > 0 || dependentCount.value > 0)
)

// Debounce timer
let fetchTimer: ReturnType<typeof setTimeout> | null = null

// Fetch with timeout helper
async function fetchWithTimeout<T>(promise: Promise<T>, timeoutMs: number): Promise<T | null> {
  const timeout = new Promise<null>((resolve) => setTimeout(() => resolve(null), timeoutMs))
  return Promise.race([promise, timeout])
}

// Fetch related files
async function fetchRelated() {
  if (!projectPath.value || props.selectedFiles.length === 0) {
    suggestions.value = []
    return
  }

  isLoadingRelated.value = true
  try {
    const result = await fetchWithTimeout(
      apiService.getSmartSuggestions(projectPath.value, props.selectedFiles),
      FETCH_TIMEOUT_MS
    )
    if (result) {
      suggestions.value = result.suggestions
      selectedRelated.value = new Set(result.suggestions.map(s => s.path))
    } else {
      suggestions.value = []
    }
  } catch {
    suggestions.value = []
  } finally {
    isLoadingRelated.value = false
  }
}

// Fetch impact analysis (background, no loading indicator)
async function fetchImpact() {
  if (!projectPath.value || props.selectedFiles.length === 0) {
    impactResult.value = null
    return
  }

  try {
    const result = await fetchWithTimeout(
      apiService.getImpactPreview(projectPath.value, props.selectedFiles),
      FETCH_TIMEOUT_MS
    )
    impactResult.value = result
  } catch {
    impactResult.value = null
  }
}

// Watch selected files
watch(() => props.selectedFiles, () => {
  if (fetchTimer) clearTimeout(fetchTimer)

  if (props.selectedFiles.length === 0) {
    suggestions.value = []
    impactResult.value = null
    return
  }

  // Debounce both fetches
  fetchTimer = setTimeout(() => {
    fetchRelated()
    fetchImpact()
  }, 500)
}, { immediate: true, deep: true })

// Helpers
function getSourceLabel(source: string): string {
  switch (source) {
    case 'git': return t('context.sourceGitShort')
    case 'arch': return t('context.sourceArchShort')
    case 'semantic': return t('context.sourceSemanticShort')
    default: return ''
  }
}

function getSourceBadgeClass(source: string): string {
  switch (source) {
    case 'git': return 'badge-git'
    case 'arch': return 'badge-arch'
    case 'semantic': return 'badge-semantic'
    default: return 'badge-default'
  }
}

function getFileIconClass(path: string): string {
  const ext = path.split('.').pop()?.toLowerCase() || ''
  const classMap: Record<string, string> = {
    'vue': 'icon-vue',
    'ts': 'icon-ts',
    'tsx': 'icon-ts',
    'js': 'icon-js',
    'jsx': 'icon-js',
    'go': 'icon-go',
    'py': 'icon-py',
    'rs': 'icon-rs',
    'json': 'icon-json',
    'yaml': 'icon-json',
    'yml': 'icon-json',
    'css': 'icon-css',
    'scss': 'icon-css',
    'html': 'icon-html',
  }
  return classMap[ext] || 'icon-default'
}

function getFileName(path: string): string {
  return path.split('/').pop() || path
}

function getFilePath(path: string): string {
  const parts = path.split('/')
  if (parts.length <= 1) return ''
  return parts.slice(0, -1).join('/')
}

function getRiskClass(level: string): string {
  switch (level) {
    case 'low': return 'risk-low'
    case 'medium': return 'risk-medium'
    case 'high': return 'risk-high'
    default: return ''
  }
}

function getRiskBarClass(level: string): string {
  switch (level) {
    case 'low': return 'bg-green-500'
    case 'medium': return 'bg-amber-500'
    case 'high': return 'bg-red-500'
    default: return 'bg-gray-500'
  }
}

function getRiskLabel(level: string): string {
  switch (level) {
    case 'low': return t('context.riskLow')
    case 'medium': return t('context.riskMedium')
    case 'high': return t('context.riskHigh')
    default: return ''
  }
}

function toggleRelated(path: string) {
  if (selectedRelated.value.has(path)) {
    selectedRelated.value.delete(path)
  } else {
    selectedRelated.value.add(path)
  }
  selectedRelated.value = new Set(selectedRelated.value)
}

function addSelectedRelated() {
  if (selectedRelated.value.size > 0) {
    emit('add-files', Array.from(selectedRelated.value))
    // Remove added files from suggestions but keep the rest
    suggestions.value = suggestions.value.filter(s => !selectedRelated.value.has(s.path))
    selectedRelated.value.clear()
    showRelatedPopup.value = false
  }
}
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

.status-btn-empty {
  background: var(--bg-2);
  color: var(--text-muted);
  border-color: var(--border-default);
  cursor: default;
}

.status-btn-empty .status-icon {
  color: #22c55e;
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

/* Smart HUD - Raycast Style Command Palette */
.smart-hud-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 12vh;
  z-index: 100;
}

.smart-hud {
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.12);
  border-radius: 0.875rem;
  width: 100%;
  max-width: 40rem;
  max-height: 60vh;
  display: flex;
  flex-direction: column;
  box-shadow: 
    0 0 0 1px rgba(139, 92, 246, 0.15),
    0 25px 60px -10px rgba(139, 92, 246, 0.35),
    0 15px 40px -5px rgba(0, 0, 0, 0.6);
  overflow: hidden;
}

/* Premium Header */
.smart-hud-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.875rem 1rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  background: linear-gradient(180deg, rgba(139, 92, 246, 0.03) 0%, transparent 100%);
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

.smart-hud-close {
  position: absolute;
  right: 0.75rem;
  top: 50%;
  transform: translateY(-50%);
  padding: 0.25rem;
  border-radius: 0.25rem;
  color: #4b5563;
  transition: all 100ms;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smart-hud-close:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #9ca3af;
}

/* Body / List */
.smart-hud-body {
  flex: 1;
  min-height: 0;
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

/* List Item with separator */
.smart-hud-item {
  display: flex;
  align-items: center;
  gap: 0.625rem;
  padding: 0.625rem 0.75rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: background 100ms;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
}

.smart-hud-item:last-child {
  border-bottom: none;
}

.smart-hud-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.smart-hud-item-selected {
  background: rgba(139, 92, 246, 0.1);
}

.smart-hud-item-selected:hover {
  background: rgba(139, 92, 246, 0.15);
}

/* Checkbox */
.smart-hud-checkbox {
  position: relative;
  width: 1rem;
  height: 1rem;
  flex-shrink: 0;
  border: 1.5px solid rgba(255, 255, 255, 0.15);
  border-radius: 0.25rem;
  background: rgba(0, 0, 0, 0.3);
  transition: all 100ms;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smart-hud-checkbox.checked {
  background: linear-gradient(135deg, #8b5cf6, #6366f1);
  border-color: transparent;
}

.smart-hud-checkbox input {
  position: absolute;
  opacity: 0;
  width: 100%;
  height: 100%;
  cursor: pointer;
}

.smart-hud-checkbox svg {
  width: 0.625rem;
  height: 0.625rem;
  color: white;
}

/* File Icon (SVG based) */
.smart-hud-file-icon {
  width: 1.25rem;
  height: 1.25rem;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.smart-hud-file-icon svg {
  width: 1rem;
  height: 1rem;
}

/* Bright vivid file icons */
.smart-hud-file-icon.icon-vue { color: #6ee7b7; }
.smart-hud-file-icon.icon-ts { color: #7dd3fc; }
.smart-hud-file-icon.icon-js { color: #fde047; }
.smart-hud-file-icon.icon-go { color: #67e8f9; }
.smart-hud-file-icon.icon-py { color: #7dd3fc; }
.smart-hud-file-icon.icon-rs { color: #fdba74; }
.smart-hud-file-icon.icon-json { color: #fde047; }
.smart-hud-file-icon.icon-css { color: #c4b5fd; }
.smart-hud-file-icon.icon-html { color: #fca5a5; }
.smart-hud-file-icon.icon-default { color: #d1d5db; }

/* File Info */
.smart-hud-file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 0;
}

.smart-hud-file-name {
  font-size: 0.8125rem;
  font-weight: 500;
  color: #e5e7eb;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.3;
}

.smart-hud-file-path {
  font-size: 0.6875rem;
  color: #a1a1aa;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  line-height: 1.3;
}

/* Source Badges - Tiny info chips with icon */
.smart-hud-badge {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  font-size: 0.5625rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  border: none;
}

.smart-hud-badge-icon {
  font-size: 0.5rem;
  opacity: 0.8;
}

.badge-git {
  color: #93c5fd;
  background: rgba(59, 130, 246, 0.12);
}

.badge-arch {
  color: #86efac;
  background: rgba(34, 197, 94, 0.12);
}

.badge-semantic {
  color: #d8b4fe;
  background: rgba(168, 85, 247, 0.12);
}

.badge-default {
  color: #a1a1aa;
  background: rgba(107, 114, 128, 0.12);
}

/* Premium Footer */
.smart-hud-footer {
  padding: 0.75rem 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(28, 31, 46, 0.5);
  backdrop-filter: blur(8px);
}

/* Premium Action Button */
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

/* Shimmer effect */
.smart-hud-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.15) 50%,
    transparent 100%
  );
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

/* Legacy popup styles for Impact Analysis */
.popup-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
  padding: 1rem;
}

.popup-content {
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: 1rem;
  width: 100%;
  max-width: 32rem;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
}

.popup-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 1rem 1.25rem;
  border-bottom: 1px solid var(--border-default);
}

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

.popup-close {
  padding: 0.25rem;
  border-radius: 0.375rem;
  color: var(--text-muted);
  transition: all 150ms;
}

.popup-close:hover {
  background: var(--bg-2);
  color: var(--text-primary);
}

.popup-body {
  flex: 1;
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
