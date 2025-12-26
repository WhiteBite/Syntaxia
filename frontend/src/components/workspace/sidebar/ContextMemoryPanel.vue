<template>
  <div class="context-detail-panel">
    <!-- Empty State: No context selected -->
    <div v-if="!selectedContext" class="detail-empty">
      <div class="detail-empty-icon">
        <svg class="w-10 h-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
            d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
        </svg>
      </div>
      <p class="detail-empty-title">{{ t('context.selectToPreview') }}</p>
      <p class="detail-empty-hint">{{ t('context.selectToPreviewHint') }}</p>
    </div>

    <!-- Context Details -->
    <template v-else>
      <!-- Header -->
      <div class="detail-header">
        <div class="detail-avatar" :class="`detail-avatar--${avatarColor}`">
          <span v-if="selectedContext.isFavorite">⭐</span>
          <span v-else>{{ initials }}</span>
        </div>
        <div class="detail-info">
          <div class="detail-title-row">
            <h2 class="detail-title">{{ contextName }}</h2>
            <span class="detail-badge">{{ selectedContext.fileCount }} {{ t('context.filesShort') }}</span>
          </div>
          <div class="detail-meta">
            <span>{{ formatSize(selectedContext.totalSize) }}</span>
            <span v-if="selectedContext.createdAt" class="detail-dot">•</span>
            <span v-if="selectedContext.createdAt">{{ formatDate(selectedContext.createdAt) }}</span>
          </div>
        </div>
      </div>

      <!-- Files List -->
      <div class="detail-content">
        <h3 class="detail-section-title">{{ t('context.contextFiles') }}</h3>
        
        <div v-if="!contextFiles.length && selectedContext.fileCount > 0" class="detail-no-files">
          <p>{{ t('context.filesListUnavailable') }}</p>
          <p class="detail-no-files-hint">{{ selectedContext.fileCount }} {{ t('context.filesInContext') }}</p>
        </div>
        
        <div v-else-if="!contextFiles.length" class="detail-no-files">
          {{ t('context.noFilesInfo') }}
        </div>
        
        <ul v-else class="detail-files">
          <li v-for="file in contextFiles" :key="file" class="detail-file">
            <div class="detail-file-icon-wrap">
              <svg class="detail-file-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div class="detail-file-info">
              <span class="detail-file-name">{{ getFileName(file) }}</span>
              <span class="detail-file-dir">{{ getFileDir(file) }}</span>
            </div>
          </li>
        </ul>
      </div>

      <!-- Footer Actions -->
      <div class="detail-footer">
        <button @click="loadAndBuildContext" :disabled="isLoading" class="detail-load-btn">
          <svg v-if="!isLoading" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          {{ isLoading ? t('context.building') : t('context.loadAndBuild') }}
        </button>
        <button @click="deleteContext" class="detail-delete-btn" :title="t('context.delete')">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useFileStore } from '@/features/files/model/file.store'
import { useContextStore } from '@/features/context/model/context.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, ref } from 'vue'
import { formatContextSize } from '@/features/context/lib/context-utils'

const emit = defineEmits<{
  (e: 'switch-to-preview'): void
}>()

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()
const settingsStore = useSettingsStore()
const uiStore = useUIStore()

const isLoading = ref(false)

// Get selected context from list (for right panel preview)
const selectedContext = computed(() => contextStore.selectedListItem)

// Get files from context
const contextFiles = computed(() => {
  return selectedContext.value?.files || []
})

// Context name without brackets
const contextName = computed(() => {
  const name = selectedContext.value?.name || t('context.untitled')
  // Remove [XX] suffix if present
  return name.replace(/\s*\[\d+\]\s*$/, '')
})

// Generate initials - SAME LOGIC as ContextListItem
const initials = computed(() => {
  const name = selectedContext.value?.name || ''
  
  if (!name) {
    const fileCount = selectedContext.value?.fileCount || 0
    return fileCount > 0 ? String(fileCount) : '?'
  }
  
  if (name.includes('Пустой')) return 'ПК'
  
  const filesMatch = name.match(/^(\d+)\s*файл/)
  if (filesMatch) {
    return filesMatch[1].length > 2 ? filesMatch[1].slice(0, 2) : filesMatch[1]
  }
  
  const words = name.split(/[\s\-_]+/).filter(w => w.length > 0)
  if (words.length >= 2) {
    return (words[0][0] + words[1][0]).toUpperCase()
  }
  
  return name.slice(0, 2).toUpperCase()
})

// Generate color based on ID
const avatarColor = computed(() => {
  const colors = ['indigo', 'purple', 'pink', 'blue', 'cyan', 'teal', 'green', 'amber']
  const source = selectedContext.value?.id || selectedContext.value?.name || ''
  let hash = 0
  for (let i = 0; i < source.length; i++) {
    hash = source.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
})

// Extract filename from path
function getFileName(filePath: string): string {
  const parts = filePath.replace(/\\/g, '/').split('/')
  return parts[parts.length - 1] || filePath
}

// Extract directory from path (relative, without project root)
function getFileDir(filePath: string): string {
  const normalized = filePath.replace(/\\/g, '/')
  const parts = normalized.split('/')
  if (parts.length <= 1) return ''
  
  // Remove filename
  parts.pop()
  
  // Try to find common prefixes to remove
  const path = parts.join('/')
  
  // Remove drive letter and common prefixes
  const cleaned = path
    .replace(/^[A-Za-z]:/, '')
    .replace(/^\/?(Sources|Projects|repos|workspace|home|Users\/[^/]+)\/[^/]+\//i, '')
    .replace(/^\/?/, '')
  
  return cleaned || '/'
}

function formatSize(bytes: number): string {
  return formatContextSize(bytes)
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('ru-RU', { 
    day: 'numeric', 
    month: 'short',
    hour: '2-digit',
    minute: '2-digit'
  })
}

async function loadAndBuildContext() {
  if (!selectedContext.value?.files?.length) {
    uiStore.addToast(t('context.noFilesInfo'), 'warning')
    return
  }
  
  isLoading.value = true
  
  try {
    // 1. Normalize paths (convert backslashes to forward slashes for consistency)
    const normalizedFiles = selectedContext.value.files.map(file => 
      file.replace(/\\/g, '/')
    )
    
    // 2. Restore file selection using selectMultiple (doesn't require node existence)
    fileStore.clearSelection()
    fileStore.selectMultiple(normalizedFiles)
    
    // 3. Build context immediately
    const options = {
      outputFormat: settingsStore.settings.context.outputFormat,
      stripComments: settingsStore.settings.context.stripComments,
      excludeTests: settingsStore.settings.context.excludeTests,
      maxTokens: settingsStore.settings.context.maxTokens,
    }
    
    await contextStore.buildContext(normalizedFiles, options)
    
    // 4. Switch to preview
    emit('switch-to-preview')
    uiStore.addToast(t('context.contextRestored'), 'success')
  } catch (error) {
    uiStore.addToast(t('context.saveError'), 'error')
  } finally {
    isLoading.value = false
  }
}

async function deleteContext() {
  if (!selectedContext.value?.id) return
  
  try {
    await contextStore.deleteContext(selectedContext.value.id)
    contextStore.clearListSelection()
    uiStore.addToast(t('context.deleted'), 'success')
  } catch {
    uiStore.addToast(t('context.deleteError'), 'error')
  }
}
</script>

<style scoped>
.context-detail-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(180deg, #131620 0%, #0f111a 100%);
}

/* Empty State */
.detail-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 32px;
  text-align: center;
}

.detail-empty-icon {
  width: 72px;
  height: 72px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(139, 92, 246, 0.08);
  border: 1px solid rgba(139, 92, 246, 0.15);
  border-radius: 20px;
  color: #8b5cf6;
  margin-bottom: 20px;
}

.detail-empty-title {
  font-size: 15px;
  font-weight: 600;
  color: #e5e7eb;
  margin-bottom: 6px;
}

.detail-empty-hint {
  font-size: 13px;
  color: #6b7280;
  max-width: 200px;
}

/* Header */
.detail-header {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 20px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.detail-avatar {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 700;
  color: white;
  flex-shrink: 0;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

/* Avatar colors */
.detail-avatar--indigo { background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%); }
.detail-avatar--purple { background: linear-gradient(135deg, #a855f7 0%, #9333ea 100%); }
.detail-avatar--pink { background: linear-gradient(135deg, #ec4899 0%, #db2777 100%); }
.detail-avatar--blue { background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%); }
.detail-avatar--cyan { background: linear-gradient(135deg, #06b6d4 0%, #0891b2 100%); }
.detail-avatar--teal { background: linear-gradient(135deg, #14b8a6 0%, #0d9488 100%); }
.detail-avatar--green { background: linear-gradient(135deg, #22c55e 0%, #16a34a 100%); }
.detail-avatar--amber { background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%); }

.detail-info {
  flex: 1;
  min-width: 0;
}

.detail-title-row {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.detail-title {
  font-size: 16px;
  font-weight: 600;
  color: white;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.detail-badge {
  flex-shrink: 0;
  padding: 3px 8px;
  background: rgba(139, 92, 246, 0.15);
  border: 1px solid rgba(139, 92, 246, 0.25);
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  color: #a78bfa;
}

.detail-meta {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #6b7280;
}

.detail-dot {
  color: #4b5563;
}

/* Content */
.detail-content {
  flex: 1;
  overflow-y: auto;
  padding: 16px 20px;
  min-height: 0;
}

/* Custom scrollbar - thin and dark */
.detail-content::-webkit-scrollbar {
  width: 6px;
}

.detail-content::-webkit-scrollbar-track {
  background: transparent;
}

.detail-content::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 20px;
}

.detail-content::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.15);
}

/* Firefox scrollbar */
.detail-content {
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.08) transparent;
}

/* Force scrollbar styling for Windows */
.detail-files {
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.08) transparent;
}

.detail-section-title {
  font-size: 10px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 12px;
}

.detail-no-files {
  font-size: 13px;
  color: #6b7280;
  text-align: center;
  padding: 24px;
}

.detail-no-files-hint {
  margin-top: 8px;
  font-size: 12px;
  color: #8b5cf6;
  font-weight: 500;
}

.detail-files {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.detail-file {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 8px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  transition: background 0.15s ease-out;
}

.detail-file:last-child {
  border-bottom: none;
}

.detail-file:hover {
  background: rgba(255, 255, 255, 0.04);
}

.detail-file-icon-wrap {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 6px;
  flex-shrink: 0;
}

.detail-file-icon {
  width: 16px;
  height: 16px;
  color: #9ca3af;
}

.detail-file-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.detail-file-name {
  font-size: 13px;
  font-weight: 500;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.detail-file-dir {
  font-size: 10px;
  color: #4b5563;
  font-family: ui-monospace, monospace;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Footer */
.detail-footer {
  display: flex;
  gap: 10px;
  padding: 16px 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(0, 0, 0, 0.2);
}

/* Premium Button */
.detail-load-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 44px;
  padding: 0 20px;
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  color: white;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 
    0 4px 12px rgba(139, 92, 246, 0.35),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.detail-load-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #a78bfa 0%, #818cf8 100%);
  transform: translateY(-1px);
  box-shadow: 
    0 6px 20px rgba(139, 92, 246, 0.45),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

.detail-load-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.detail-load-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* Delete Button - same height, more prominent */
.detail-delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  background: rgba(239, 68, 68, 0.08);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 10px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.detail-delete-btn:hover {
  background: rgba(239, 68, 68, 0.15);
  border-color: rgba(239, 68, 68, 0.4);
  color: #f87171;
  transform: translateY(-1px);
}

.detail-delete-btn svg {
  width: 18px;
  height: 18px;
}
</style>
