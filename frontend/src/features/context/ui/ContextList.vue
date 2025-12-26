<template>
  <div class="context-list-panel">
    <!-- Compact Toolbar -->
    <ContextListSearch
      v-if="contextStore.contextList.length > 0"
      v-model:search-query="list.searchQuery.value"
      v-model:show-favorites-only="list.showFavoritesOnly.value"
      :sort-label="sortLabel"
      :selected-count="list.selectedContexts.value.size"
      @cycle-sort="cycleSortBy"
      @copy-selected="list.copySelectedContext"
      @delete-selected="list.deleteSelected"
    />

    <!-- Loading -->
    <ContextListLoading v-if="contextStore.isLoading" />

    <!-- Empty State -->
    <ContextListEmpty v-else-if="contextStore.contextList.length === 0" />

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
    <ContextListFooter
      :can-save="canSaveContext"
      @save="showSaveDialog = true"
    />

    <!-- Save Dialog -->
    <ContextSaveDialog
      :show="showSaveDialog"
      :file-count="fileStore.selectedPaths.size"
      @close="showSaveDialog = false"
      @save="saveContext"
    />

    <!-- Delete Confirmation Modal -->
    <ContextDeleteModal
      :show="list.deleteModal.show"
      :message="list.deleteModal.message"
      @close="list.deleteModal.show = false"
      @confirm="list.executeDelete"
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useFileStore } from '@/features/files/model/file.store'
import { apiService } from '@/services/api.service'
import { useProjectStore } from '@/stores/project.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useContextList } from '../composables/useContextList'
import { useContextStore, type ContextSummary } from '../model/context.store'
import ContextDeleteModal from './ContextDeleteModal.vue'
import ContextListEmpty from './ContextListEmpty.vue'
import ContextListFooter from './ContextListFooter.vue'
import ContextListItem from './ContextListItem.vue'
import ContextListLoading from './ContextListLoading.vue'
import ContextListSearch from './ContextListSearch.vue'
import ContextSaveDialog from './ContextSaveDialog.vue'

const logger = useLogger('ContextList')

const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const list = useContextList()

// Save dialog state
const showSaveDialog = ref(false)

const canSaveContext = computed(() => fileStore.selectedPaths.size > 0 && !!projectStore.currentPath)

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
  contextStore.selectListItem(contextId)
}

// Load context (double click) - loads files into editor
async function loadContext(contextId: string) {
  try {
    contextStore.selectListItem(contextId)
    await contextStore.loadContextContent(contextId, 0, 0)
    emit('switch-to-preview')
  } catch (error) {
    logger.error('Failed to load context:', error)
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

async function saveContext(topic: string, summary: string) {
  if (!projectStore.currentPath || !topic.trim()) return
  try {
    // Normalize paths before saving (convert backslashes to forward slashes)
    const normalizedFiles = Array.from(fileStore.selectedPaths).map(path => 
      path.replace(/\\/g, '/')
    )
    
    await apiService.saveContextMemory(
      projectStore.currentPath,
      topic.trim(),
      summary.trim(),
      normalizedFiles
    )
    uiStore.addToast(t('context.contextSaved'), 'success')
    showSaveDialog.value = false
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

/* List */
.context-list {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-height: 0;
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
</style>
