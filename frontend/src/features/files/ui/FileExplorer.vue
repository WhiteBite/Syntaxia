<template>
  <div class="file-explorer">
    <!-- Header: High-density single row -->
    <div class="flex flex-col border-b border-white/5 bg-[#0f111a]/80 backdrop-blur-md">
      <div class="flex items-center justify-between h-9 px-3">
        <div class="flex items-center gap-2 overflow-hidden">
          <h2 class="text-[10px] font-bold text-gray-400 uppercase tracking-widest select-none truncate">
            {{ t('files.title') }}
          </h2>
          <div v-if="fileStore.selectedCount > 0" 
            class="px-1.5 py-0.5 bg-indigo-500/20 text-indigo-300 rounded text-[9px] font-black border border-indigo-500/20">
            {{ fileStore.selectedCount }}
          </div>
        </div>
        
        <div class="flex items-center gap-0.5">
          <!-- Recommendations Lamp -->
          <button 
            v-if="recommendationsCount > 0"
            class="flex items-center justify-center w-7 h-7 rounded-full bg-amber-500/10 text-amber-500 hover:bg-amber-500/20 transition-all group/lamp"
            @click="toggleAnalysisPopup"
            :title="t('context.recommendations')"
          >
            <div class="relative">
              <Sparkles class="w-3.5 h-3.5" />
              <span class="absolute -top-1.5 -right-1.5 flex h-3.5 w-3.5 items-center justify-center rounded-full bg-amber-600 text-[8px] font-black text-white border-2 border-[#0f111a] animate-pulse">
                {{ recommendationsCount }}
              </span>
            </div>
          </button>

          <div v-if="recommendationsCount > 0" class="w-px h-3 bg-white/10 mx-1"></div>

          <!-- Actions -->
          <BaseButton 
            variant="ghost" size="xs" icon-only 
            class="w-7 h-7 text-gray-500 hover:text-white"
            @click="fileStore.expandAll()" :title="t('files.expandAll')"
          >
            <ChevronDownSquare class="w-3.5 h-3.5" />
          </BaseButton>
          <BaseButton 
            variant="ghost" size="xs" icon-only 
            class="w-7 h-7 text-gray-500 hover:text-white"
            @click="fileStore.collapseAll()" :title="t('files.collapseAll')"
          >
            <ChevronUpSquare class="w-3.5 h-3.5" />
          </BaseButton>

          <div class="w-px h-3 bg-white/10 mx-1"></div>

          <ViewOptionsDropdown />
          <SystemFiltersDropdown @open-advanced="showAdvancedFilters = true" />

          <BaseButton 
            variant="ghost" size="xs" icon-only 
            class="w-7 h-7 text-gray-500 hover:text-white"
            @click="explorer.handleRefresh" :title="t('files.refresh')"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': fileStore.isLoading }" />
          </BaseButton>
          
          <SettingsPopover 
            @open-ignore-rules="ignoreRulesModalRef?.open()" 
            @settings-changed="explorer.handleSettingsChange"
          />
        </div>
      </div>

      <!-- Compact Search -->
      <div class="px-2 pb-2">
        <div class="relative group">
          <input 
            v-model="explorer.searchQuery.value"
            type="text"
            :placeholder="t('files.searchShort')"
            class="w-full h-7 bg-white/5 text-xs text-gray-300 placeholder-gray-600 rounded border border-transparent focus:border-indigo-500/30 focus:bg-white/10 focus:outline-none transition-all pl-7 pr-7"
            @input="explorer.handleSearch" 
            @keydown.escape="clearSearch"
          />
          <SearchIcon class="w-3 h-3 text-gray-500 absolute left-2.5 top-1/2 -translate-y-1/2 pointer-events-none group-focus-within:text-indigo-400" />
          <button v-if="explorer.searchQuery.value" @click="clearSearch" 
            class="absolute right-1.5 top-1/2 -translate-y-1/2 p-0.5 text-gray-500 hover:text-white">
            <X class="w-3 h-3" />
          </button>
        </div>
      </div>
    </div>

    <!-- File Tree -->
    <div class="file-explorer__tree" data-tour="file-tree">
      <div v-if="fileStore.isLoading" class="loading-state">
        <SkeletonFileTree :rows="10" />
      </div>

      <div v-else-if="fileStore.nodes.length === 0" class="h-full">
        <BaseEmptyState
          :title="t('files.noFiles')"
          size="md"
        >
          <template #icon>
            <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
            </svg>
          </template>
        </BaseEmptyState>
      </div>

      <template v-else>
        <!-- Search Results Header -->
        <div v-if="explorer.searchQuery.value && fileStore.searchResults.length > 0" class="text-xs text-gray-400 mb-2 px-2">
          {{ fileStore.searchResults.length }} {{ t('files.results') }}
        </div>

        <!-- No Search Results -->
        <div v-if="explorer.searchQuery.value && fileStore.searchResults.length === 0" class="h-full">
          <BaseEmptyState
            :title="t('files.noSearchResults')"
            :description="t('files.tryAdjustingFilters')"
            size="md"
          >
            <template #icon>
              <svg class="w-12 h-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </template>
            <template #action>
              <BaseButton variant="ghost" @click="explorer.searchQuery.value = ''">
                {{ t('files.clearSearch') }}
              </BaseButton>
            </template>
          </BaseEmptyState>
        </div>

        <!-- Virtualized File Tree -->
        <VirtualFileTree
          v-else-if="!explorer.searchQuery.value || fileStore.searchResults.length > 0"
          :nodes="explorer.searchQuery.value ? fileStore.searchResults : fileStore.filteredNodes"
          :compact-mode="settingsStore.settings.fileExplorer.compactNestedFolders"
          :allow-select-binary="settingsStore.settings.fileExplorer.allowSelectBinary"
          @toggle-select="explorer.handleToggleSelect"
          @toggle-expand="explorer.handleToggleExpand"
          @contextmenu="handleContextMenu"
          @quicklook="explorer.handleQuickLook"
        />
      </template>
    </div>

    <!-- Context Menu -->
    <FileContextMenu :node="contextMenu.targetNode.value" :position="contextMenu.position.value"
      :visible="contextMenu.isVisible.value" @action="handleContextMenuAction" @close="contextMenu.hide" />

    <!-- Modals -->
    <IgnoreRulesModal ref="ignoreRulesModalRef" />
    <AdvancedFiltersModal 
      :is-open="showAdvancedFilters" 
      :filters="quickFilters.typeFilters.value"
      :get-count="quickFilters.getFilterCount"
      @close="showAdvancedFilters = false"
      @reset="quickFilters.resetFilters"
      @update-extensions="quickFilters.updateFilterExtensions"
    />
    <QuickLookModal v-model="explorer.quickLookVisible.value" :file-path="explorer.quickLookPath.value" @add-to-context="explorer.handleAddToContext" />
    <DependencyVisualizerModal
      v-model="showDependencyModal"
      :file-path="selectedFileForDeps"
    />
    <AnalysisStatusBar v-if="false" :selected-files="[]" @add-files="handleAddSuggestedFiles" />

    <!-- Footer: Magic Control Bar -->
    <div class="file-explorer__footer">
      <CommandBar 
        data-tour="build-button"
        :selected-count="fileStore.selectedPaths.size"
        :is-building="contextStore.isBuilding"
        @build="$emit('build-context')"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useContextMenu } from '@/composables/useContextMenu'
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useContextStore } from '@/features/context'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { BaseButton, BaseEmptyState } from '@/components/ui'
import { RefreshCw, Search as SearchIcon, X, ChevronDownSquare, ChevronUpSquare, Sparkles } from 'lucide-vue-next'
import { defineAsyncComponent, onMounted, onUnmounted, ref, watch, computed } from 'vue'

import { useFileExplorer } from '../composables/useFileExplorer'
import { useQuickFilters } from '../composables/useQuickFilters'
import { provideHoveredFile } from '../composables/useHoveredFile'
import { useFileStore, type FileNode } from '../model/file.store'
import { useAnalysisStatus } from '../composables/useAnalysisStatus'
import AnalysisStatusBar from './AnalysisStatusBar.vue'
import CommandBar from './CommandBar.vue'
import FileContextMenu from './FileContextMenu.vue'
import SettingsPopover from './SettingsPopover.vue'
import VirtualFileTree from './VirtualFileTree.vue'
import ViewOptionsDropdown from './ViewOptionsDropdown.vue'
import SystemFiltersDropdown from './SystemFiltersDropdown.vue'
import SkeletonFileTree from '@/components/SkeletonFileTree.vue'

const QuickLookModal = defineAsyncComponent(() => import('@/components/QuickLookModal.vue'))
const IgnoreRulesModal = defineAsyncComponent(() => import('./IgnoreRulesModal.vue'))
const DependencyVisualizerModal = defineAsyncComponent(() => import('./DependencyVisualizerModal.vue'))
const AdvancedFiltersModal = defineAsyncComponent(() => import('./FilterSettingsModal.vue'))

provideHoveredFile()

const fileStore = useFileStore()
const contextStore = useContextStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const explorer = useFileExplorer()
const quickFilters = useQuickFilters()
const contextMenu = useContextMenu()
const logger = useLogger('FileExplorer')

// Real analysis status integration for the "Lamp"
const selectedFilesRef = computed(() => Array.from(fileStore.selectedPaths))
const analysis = useAnalysisStatus({
  selectedFiles: selectedFilesRef,
  onAddFiles: (files) => fileStore.selectMultiple(files)
})

const recommendationsCount = computed(() => (analysis.relatedCount.value || 0) + (analysis.dependentCount.value || 0))
const toggleAnalysisPopup = () => {
    if (analysis.relatedCount.value > 0) analysis.showRelatedPopup.value = true
    else if (analysis.dependentCount.value > 0) analysis.showImpactPopup.value = true
}

const ignoreRulesModalRef = ref<InstanceType<typeof IgnoreRulesModal>>()
const searchInputRef = ref<HTMLInputElement | null>(null)
const showAdvancedFilters = ref(false)

// Dependency modal state
const showDependencyModal = ref(false)
const selectedFileForDeps = ref('')

function clearSearch() {
  explorer.clearSearch()
}

defineExpose({ searchInputRef })

defineEmits<{
  (e: 'preview-file', filePath: string): void
  (e: 'build-context'): void
}>()

async function handleContextMenuAction(payload: { type: string; node: FileNode }) {
  const result = await explorer.handleContextMenuAction(payload)
  
  // Handle special actions that need parent component interaction
  if (result && result.action === 'showDependencies') {
    selectedFileForDeps.value = result.path
    showDependencyModal.value = true
  }
}

function handleContextMenu(node: FileNode, event: MouseEvent) {
  contextMenu.show(node, event)
}

function handleAddSuggestedFiles(files: string[]) {
  const normalizedFiles = files
    .map(path => path.replace(/\\/g, '/'))
    .filter(path => !fileStore.selectedPaths.has(path))
  
  if (normalizedFiles.length > 0) {
    fileStore.selectMultiple(normalizedFiles)
  }
}

function handleUndoSelection() {
  if (fileStore.undoSelection()) {
    uiStore.addToast(t('files.undoSelection'), 'info')
  }
}

function handleRedoSelection() {
  if (fileStore.redoSelection()) {
    uiStore.addToast(t('files.redoSelection'), 'info')
  }
}

onMounted(async () => {
  explorer.initialize()
  
  window.addEventListener('global-undo-selection', handleUndoSelection)
  window.addEventListener('global-redo-selection', handleRedoSelection)
  
  try {
    if (!projectStore.currentPath) {
      uiStore.addToast('No project selected', 'warning')
      return
    }
    await fileStore.loadFileTree(projectStore.currentPath)
  } catch (error) {
    logger.error('Failed to load file tree:', error)
    uiStore.addToast('Failed to load project files. Please try again.', 'error')
  }
})

watch(() => projectStore.currentPath, async (newPath, oldPath) => {
  if (newPath && newPath !== oldPath) {
    fileStore.clearSelection()
    contextStore.clearContext()
    try {
      await fileStore.loadFileTree(newPath)
    } catch (error) {
      logger.error('Failed to load file tree after project change:', error)
      uiStore.addToast('Failed to load project files', 'error')
    }
  }
})

onUnmounted(() => {
  explorer.cleanup()
  window.removeEventListener('global-undo-selection', handleUndoSelection)
  window.removeEventListener('global-redo-selection', handleRedoSelection)
})
</script>

<style scoped>
.file-explorer {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
  background: transparent;
  overflow: hidden;
}

.file-explorer__tree {
  flex: 1 1 0;
  min-height: 0;
  overflow: hidden;
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
}

.file-explorer__tree .loading-state {
  flex: 1;
  overflow-y: auto;
}

.file-explorer__tree > :deep(.virtual-tree-wrapper) {
  flex: 1 1 0;
  min-height: 0;
}

.file-explorer__footer {
  flex-shrink: 0;
  padding: 0.5rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  background: #0f111a;
}
</style>
