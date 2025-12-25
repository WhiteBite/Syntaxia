<template>
  <div class="context-list-panel">
    <!-- Compact Toolbar -->
    <div v-if="contextStore.contextList.length > 0" class="context-toolbar">
      <!-- Search -->
      <div class="context-search">
        <svg class="context-search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input 
          v-model="list.searchQuery.value" 
          type="text" 
          :placeholder="t('context.search')" 
          class="context-search-input"
        />
        <!-- Favorites toggle inside search -->
        <button 
          @click="list.showFavoritesOnly.value = !list.showFavoritesOnly.value"
          class="context-filter-btn"
          :class="{ active: list.showFavoritesOnly.value }"
          :title="t('context.showFavorites')"
        >
          <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
          </svg>
        </button>
        <!-- Sort toggle -->
        <button 
          @click="cycleSortBy"
          class="context-filter-btn"
          :title="sortLabel"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
          </svg>
        </button>
      </div>
      
      <!-- Bulk Actions (when selected) -->
      <div v-if="list.selectedContexts.value.size > 0" class="context-bulk-actions">
        <span class="context-bulk-count">{{ list.selectedContexts.value.size }}</span>
        <button @click="list.copySelectedContext" class="context-bulk-btn" :title="t('context.copyToClipboard')">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        </button>
        <button @click="list.deleteSelected" class="context-bulk-btn context-bulk-btn--danger" :title="t('context.deleteSelected')">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="contextStore.isLoading" class="context-loading">
      <svg class="animate-spin h-6 w-6 text-purple-500" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
      </svg>
    </div>

    <!-- Empty State -->
    <div v-else-if="contextStore.contextList.length === 0" class="context-empty">
      <div class="context-empty-icon">
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
        </svg>
      </div>
      <p class="context-empty-title">{{ t('context.noSaved') }}</p>
      <p class="context-empty-hint">{{ t('context.buildToSee') }}</p>
    </div>

    <!-- Context List -->
    <div v-else class="context-list">
      <ContextListItem
        v-for="(context, index) in list.filteredContexts.value"
        :key="context.id"
        :context="context"
        :index="index"
        :is-selected="list.selectedContexts.value.has(context.id)"
        :is-active="contextStore.selectedListItem?.id === context.id"
        :is-editing="list.editingId.value === context.id"
        :editing-name="list.editingName.value"
        :drag-over-index="list.dragOverIndex.value"
        :search-query="list.searchQuery.value"
        :show-favorites-only="list.showFavoritesOnly.value"
        @select="selectContext"
        @toggle-select="list.toggleSelect"
        @toggle-favorite="list.toggleFavorite"
        @load="loadContext"
        @restore-selection="handleRestoreSelection"
        @start-rename="handleStartRename"
        @save-rename="list.saveRename"
        @cancel-rename="list.cancelRename"
        @update:editing-name="list.editingName.value = $event"
        @copy="list.copyContext"
        @duplicate="list.duplicateContext"
        @export="list.exportContext"
        @delete="list.confirmDelete"
        @drag-start="list.handleDragStart"
        @drag-over="list.handleDragOver"
        @drag-leave="list.handleDragLeave"
        @drop="list.handleDrop"
        @drag-end="list.handleDragEnd"
      />
    </div>

    <!-- Footer: Save Current Context -->
    <div v-if="canSaveContext" class="context-footer">
      <button @click="showSaveDialog = true" class="context-save-btn">
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        {{ t('context.saveCurrentContext') }}
      </button>
    </div>

    <!-- Save Dialog (Premium Design) -->
    <Teleport to="body">
      <div v-if="showSaveDialog" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="showSaveDialog = false"></div>
        <div class="save-modal">
          <!-- Gradient accent line -->
          <div class="save-modal-accent"></div>
          
          <h3 class="save-modal-title">{{ t('context.saveContext') }}</h3>
          
          <form @submit.prevent="saveContext" class="save-modal-content">
            <div class="save-modal-field">
              <label class="save-modal-label">{{ t('context.topic') }}</label>
              <input 
                ref="topicInputRef"
                v-model="saveTopic" 
                type="text" 
                class="save-modal-input" 
                :placeholder="t('context.topicPlaceholder')"
                @keyup.enter="saveTopic.trim() && saveContext()"
              />
            </div>
            <div class="save-modal-field">
              <label class="save-modal-label">{{ t('context.summary') }}</label>
              <textarea 
                v-model="saveSummary" 
                class="save-modal-textarea" 
                :placeholder="t('context.summaryPlaceholder')" 
              />
            </div>
          </form>
          
          <div class="save-modal-footer">
            <!-- File count badge -->
            <div class="save-modal-badge">
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span>{{ fileStore.selectedPaths.size }} {{ t('context.filesShort') }}</span>
            </div>
            
            <div class="save-modal-actions">
              <button type="button" @click="showSaveDialog = false" class="save-modal-cancel">
                {{ t('context.cancel') }}
              </button>
              <button 
                type="submit" 
                @click="saveContext" 
                class="save-modal-submit" 
                :disabled="!saveTopic.trim()"
              >
                {{ t('common.save') }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Delete Confirmation Modal -->
    <Teleport to="body">
      <div v-if="list.deleteModal.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="list.deleteModal.show = false"></div>
        <div class="relative confirm-modal">
          <div class="flex items-center gap-3 mb-4">
            <div class="w-10 h-10 rounded-lg bg-red-500/20 flex items-center justify-center">
              <svg class="w-5 h-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </div>
            <div>
              <h3 class="text-base font-semibold text-white">{{ t('context.confirmDelete') }}</h3>
              <p class="text-sm text-gray-400">{{ list.deleteModal.message }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2">
            <button @click="list.deleteModal.show = false" class="btn btn-ghost btn-sm">{{ t('context.cancel') }}</button>
            <button @click="list.executeDelete" class="btn btn-sm bg-red-500 hover:bg-red-600 text-white">{{ t('context.delete') }}</button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useFileStore } from '@/features/files/model/file.store'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useContextList } from '../composables/useContextList'
import { useContextStore, type ContextSummary } from '../model/context.store'
import ContextListItem from './ContextListItem.vue'

const logger = useLogger('ContextList')

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const list = useContextList()

// Save dialog state
const showSaveDialog = ref(false)
const saveTopic = ref('')
const saveSummary = ref('')
const topicInputRef = ref<HTMLInputElement | null>(null)

// Auto-focus topic input when dialog opens
watch(showSaveDialog, (isOpen) => {
  if (isOpen) {
    nextTick(() => {
      topicInputRef.value?.focus()
    })
  }
})

const canSaveContext = computed(() => fileStore.selectedPaths.size > 0 && projectStore.currentPath)

const emit = defineEmits<{
  (e: 'switch-to-preview'): void
}>()

const sortLabel = computed(() => {
  const labels: Record<string, string> = {
    date: t('context.sortByDate'),
    name: t('context.sortByName'),
    size: t('context.sortBySize')
  }
  return labels[list.sortBy.value] || labels.date
})

function cycleSortBy() {
  const options: Array<'date' | 'name' | 'size'> = ['date', 'name', 'size']
  const current = options.indexOf(list.sortBy.value)
  list.sortBy.value = options[(current + 1) % options.length]
}

// Select context (single click) - shows details in right panel
function selectContext(contextId: string) {
  // Use selectListItem for instant preview without loading content
  contextStore.selectListItem(contextId)
}

// Load context (double click) - loads files into editor
async function loadContext(contextId: string) {
  try {
    contextStore.selectListItem(contextId)
    await contextStore.loadContextContent(contextId, 0, 0)
    emit('switch-to-preview')
  } catch (error) {
    logger.error('[ContextList] Failed to load context:', error)
  }
}

function handleRestoreSelection(context: ContextSummary) {
  if (!context.files || context.files.length === 0) {
    uiStore.addToast(t('context.noFilesToRestore'), 'warning')
    return
  }
  
  fileStore.clearSelection()
  fileStore.selectMultiple(context.files)
  uiStore.addToast(
    t('context.selectionRestored').replace('{count}', String(context.files.length)).replace('{name}', context.name || ''),
    'success'
  )
}

function handleStartRename(context: ContextSummary) {
  list.startRename(context)
}

async function saveContext() {
  if (!projectStore.currentPath || !saveTopic.value.trim()) return
  try {
    // Normalize paths before saving (convert backslashes to forward slashes)
    const normalizedFiles = Array.from(fileStore.selectedPaths).map(path => 
      path.replace(/\\/g, '/')
    )
    
    await apiService.saveContextMemory(
      projectStore.currentPath,
      saveTopic.value.trim(),
      saveSummary.value.trim(),
      normalizedFiles
    )
    uiStore.addToast(t('context.contextSaved'), 'success')
    showSaveDialog.value = false
    saveTopic.value = ''
    saveSummary.value = ''
    list.refresh()
  } catch {
    uiStore.addToast(t('context.saveError'), 'error')
  }
}

onMounted(() => {
  list.refresh()
  window.addEventListener('keydown', list.handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', list.handleKeydown)
})
</script>

<style scoped>
.context-list-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background: linear-gradient(180deg, #131620 0%, #0f111a 100%);
}

/* Toolbar */
.context-toolbar {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

/* Search */
.context-search {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 0 10px;
  height: 32px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  transition: all 0.15s ease-out;
}

.context-search:focus-within {
  border-color: #8b5cf6;
  background: rgba(139, 92, 246, 0.05);
}

.context-search-icon {
  width: 14px;
  height: 14px;
  color: #6b7280;
  flex-shrink: 0;
}

.context-search-input {
  flex: 1;
  background: none;
  border: none;
  outline: none;
  font-size: 12px;
  color: #e5e7eb;
}

.context-search-input::placeholder {
  color: #4b5563;
}

.context-filter-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-filter-btn:hover {
  color: #e5e7eb;
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.12);
}

.context-filter-btn.active {
  color: #facc15;
  background: rgba(250, 204, 21, 0.1);
  border-color: rgba(250, 204, 21, 0.2);
}

/* Bulk Actions */
.context-bulk-actions {
  display: flex;
  align-items: center;
  gap: 6px;
}

.context-bulk-count {
  font-size: 11px;
  font-weight: 600;
  color: #8b5cf6;
  padding: 2px 8px;
  background: rgba(139, 92, 246, 0.15);
  border-radius: 4px;
}

.context-bulk-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.context-bulk-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.context-bulk-btn--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

/* Loading */
.context-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}

/* Empty State */
.context-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  flex: 1;
  padding: 32px;
  text-align: center;
}

.context-empty-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 16px;
  color: #4b5563;
  margin-bottom: 16px;
}

.context-empty-title {
  font-size: 14px;
  font-weight: 500;
  color: #9ca3af;
  margin-bottom: 4px;
}

.context-empty-hint {
  font-size: 12px;
  color: #6b7280;
}

/* List */
.context-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  scrollbar-width: thin;
  scrollbar-color: rgba(255, 255, 255, 0.1) transparent;
}

.context-list::-webkit-scrollbar {
  width: 6px;
}

.context-list::-webkit-scrollbar-track {
  background: transparent;
}

.context-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 20px;
}

.context-list::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}

/* Modal */
.confirm-modal {
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  max-width: 360px;
  width: 100%;
}

/* Footer */
.context-footer {
  padding: 12px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: linear-gradient(180deg, transparent 0%, rgba(139, 92, 246, 0.03) 100%);
}

.context-save-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 12px 16px;
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
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

.context-save-btn:hover {
  background: linear-gradient(135deg, #a78bfa 0%, #8b5cf6 100%);
  transform: translateY(-1px);
  box-shadow: 
    0 6px 20px rgba(139, 92, 246, 0.45),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

/* Save Modal - Premium Design */
.save-modal {
  position: relative;
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 28px;
  max-width: 440px;
  width: 100%;
  box-shadow: 
    0 24px 48px rgba(0, 0, 0, 0.5),
    0 0 80px rgba(139, 92, 246, 0.15);
  overflow: hidden;
}

.save-modal-accent {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #8b5cf6 0%, #6366f1 50%, #3b82f6 100%);
}

.save-modal-title {
  font-size: 18px;
  font-weight: 700;
  color: white;
  margin-bottom: 24px;
}

.save-modal-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.save-modal-field {
  display: flex;
  flex-direction: column;
}

.save-modal-label {
  display: block;
  font-size: 10px;
  font-weight: 700;
  color: #6b7280;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  margin-bottom: 8px;
}

.save-modal-input {
  width: 100%;
  padding: 14px 16px;
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  font-size: 14px;
  color: white;
  outline: none;
  transition: all 0.2s ease-out;
}

.save-modal-input::placeholder {
  color: #4b5563;
}

.save-modal-input:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.save-modal-textarea {
  width: 100%;
  height: 90px;
  padding: 14px 16px;
  background: #0f111a;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  font-size: 14px;
  color: white;
  outline: none;
  resize: none;
  transition: all 0.2s ease-out;
}

.save-modal-textarea::placeholder {
  color: #4b5563;
}

.save-modal-textarea:focus {
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.15);
}

.save-modal-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 28px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.save-modal-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  background: rgba(139, 92, 246, 0.1);
  border: 1px solid rgba(139, 92, 246, 0.2);
  border-radius: 8px;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #a78bfa;
  font-weight: 500;
}

.save-modal-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.save-modal-cancel {
  padding: 10px 18px;
  background: none;
  border: none;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 500;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.save-modal-cancel:hover {
  color: white;
  background: rgba(255, 255, 255, 0.05);
}

.save-modal-submit {
  padding: 12px 28px;
  background: linear-gradient(135deg, #9333ea 0%, #7c3aed 50%, #6366f1 100%);
  border: none;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 700;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease-out;
  box-shadow: 
    0 4px 16px rgba(139, 92, 246, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1) inset,
    0 1px 0 rgba(255, 255, 255, 0.15) inset;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.2);
}

.save-modal-submit:hover:not(:disabled) {
  background: linear-gradient(135deg, #a855f7 0%, #8b5cf6 50%, #818cf8 100%);
  transform: translateY(-2px);
  box-shadow: 
    0 8px 24px rgba(139, 92, 246, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.15) inset,
    0 1px 0 rgba(255, 255, 255, 0.2) inset;
}

.save-modal-submit:active:not(:disabled) {
  transform: translateY(0);
}

.save-modal-submit:disabled {
  opacity: 0.4;
  cursor: not-allowed;
  background: linear-gradient(135deg, #6b7280 0%, #4b5563 100%);
  box-shadow: none;
}
</style>
