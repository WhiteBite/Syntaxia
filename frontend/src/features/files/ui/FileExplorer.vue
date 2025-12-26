<template>
  <div class="file-explorer">
    <!-- Header -->
    <div class="file-explorer__header">
      <div class="panel-header-unified">
        <div class="panel-header-unified-title">
          <div class="panel-header-unified-icon panel-header-unified-icon-indigo">
            <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"></path>
            </svg>
          </div>
          <h2>{{ t('files.title') }}</h2>
        </div>

        <div class="flex items-center gap-1">
          <!-- Selection Stats -->
          <div class="flex items-center gap-2 text-xs mr-2">
            <div class="relative group">
              <span class="chip-unified chip-unified-accent cursor-help">
                {{ fileStore.selectedCount }}
              </span>
              <!-- Tooltip with stats -->
              <div class="absolute right-0 top-full mt-2 hidden group-hover:block z-50 px-4 py-3 bg-gray-800/95 backdrop-blur-sm border border-gray-700/50 rounded-xl shadow-2xl whitespace-nowrap">
                <div class="text-xs space-y-1.5">
                  <div class="text-white font-semibold mb-2">Selection Stats</div>
                  <div class="text-gray-400">Files: <span class="text-white">{{ fileStore.selectedCount }}</span></div>
                  <div class="text-gray-400">Total: <span class="text-white">{{ explorer.totalFileCount.value }}</span> files</div>
                  <div class="text-gray-400">Progress: <span class="text-indigo-400">{{ explorer.selectionProgress.value }}%</span></div>
                  <div class="text-gray-400">Est. Size: <span class="text-white">{{ Math.round(fileStore.estimatedContextSize * 100) / 100 }}MB</span></div>
                  <div class="text-gray-400">Est. Tokens: <span class="text-emerald-400">~{{ Math.round(fileStore.estimatedTokenCount / 1000) }}K</span></div>
                </div>
              </div>
            </div>
            <span class="text-gray-400">{{ t('files.selected') }}</span>
            <!-- Clear selection button -->
            <button 
              v-if="fileStore.selectedCount > 0"
              @click="fileStore.clearSelection"
              class="p-1 rounded hover:bg-red-500/20 text-gray-500 hover:text-red-400 transition-colors"
              :title="t('files.clearSelection')"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Settings Popover -->
          <SettingsPopover 
            @open-ignore-rules="ignoreRulesModalRef?.open()" 
            @settings-changed="explorer.handleSettingsChange"
          />

          <!-- Refresh Button -->
          <button @click="explorer.handleRefresh"
            class="p-2 rounded-lg hover:bg-gray-700/50 text-gray-400 hover:text-white transition-colors"
            :title="t('files.refresh')" :aria-label="t('files.refresh')">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Breadcrumbs -->
      <div v-if="fileStore.breadcrumbs.length > 0" class="px-3 pb-3">
        <BreadcrumbsNav 
          :segments="fileStore.breadcrumbs" 
          :root-name="fileStore.projectName"
          @navigate="handleBreadcrumbNavigate"
          @open-in-explorer="handleOpenInExplorer"
        />
      </div>
    </div>

    <!-- Quick Filters -->
    <QuickFiltersBar />

    <!-- Search -->
    <div class="file-explorer__search">
      <div class="relative group">
        <input 
          ref="searchInputRef"
          v-model="explorer.searchQuery.value" 
          type="text" 
          :placeholder="t('files.searchShort')" 
          :aria-label="t('files.searchShort')"
          class="search-input pr-8"
          :class="{ 'search-input-active': explorer.searchQuery.value }"
          @input="explorer.handleSearch" 
          @keydown.escape="clearSearch"
        />
        <svg class="input-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <button 
          v-if="explorer.searchQuery.value"
          @click="clearSearch"
          class="search-clear-btn"
          :title="t('files.clear')"
        >
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- File Tree - MAIN SCROLLABLE AREA -->
    <div class="file-explorer__tree" data-tour="file-tree">
      <div v-if="fileStore.isLoading" class="flex items-center justify-center h-full">
        <svg class="loading-spinner" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
      </div>

      <div v-else-if="fileStore.nodes.length === 0" class="empty-state h-full">
        <p class="empty-state-text">{{ t('files.noFiles') }}</p>
      </div>

      <template v-else>
        <!-- Search Results Header -->
        <div v-if="explorer.searchQuery.value && fileStore.searchResults.length > 0" class="text-xs text-gray-400 mb-2 px-2">
          {{ fileStore.searchResults.length }} {{ t('files.results') }}
        </div>

        <!-- No Search Results -->
        <div v-if="explorer.searchQuery.value && fileStore.searchResults.length === 0" class="empty-state h-full">
          <div class="empty-state-content">
            <svg class="empty-state-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            <p class="empty-state-text">{{ t('files.noSearchResults') }}</p>
            <button class="empty-state-action" @click="explorer.searchQuery.value = ''">
              {{ t('files.clearSearch') }}
            </button>
          </div>
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
      :visible="contextMenu.isVisible.value" @action="explorer.handleContextMenuAction" @close="contextMenu.hide" />

    <!-- Modals -->
    <IgnoreRulesModal ref="ignoreRulesModalRef" />
    <QuickLookModal v-model="explorer.quickLookVisible.value" :file-path="explorer.quickLookPath.value" @add-to-context="explorer.handleAddToContext" />

    <!-- Analysis Status Bar -->
    <AnalysisStatusBar 
      :selected-files="Array.from(fileStore.selectedPaths)"
      @add-files="handleAddSuggestedFiles"
    />

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
import { defineAsyncComponent, onMounted, onUnmounted, ref, watch } from 'vue'
import { useFileExplorer } from '../composables/useFileExplorer'
import { provideHoveredFile } from '../composables/useHoveredFile'
import { useFileStore, type FileNode } from '../model/file.store'
import AnalysisStatusBar from './AnalysisStatusBar.vue'
import BreadcrumbsNav from './BreadcrumbsNav.vue'
import CommandBar from './CommandBar.vue'
import FileContextMenu from './FileContextMenu.vue'
import QuickFiltersBar from './QuickFiltersBar.vue'
import SettingsPopover from './SettingsPopover.vue'
import VirtualFileTree from './VirtualFileTree.vue'

const QuickLookModal = defineAsyncComponent(() => import('@/components/QuickLookModal.vue'))
const IgnoreRulesModal = defineAsyncComponent(() => import('./IgnoreRulesModal.vue'))

provideHoveredFile()

const fileStore = useFileStore()
const contextStore = useContextStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const settingsStore = useSettingsStore()
const { t } = useI18n()

const explorer = useFileExplorer()
const contextMenu = useContextMenu()
const logger = useLogger('FileExplorer')

const ignoreRulesModalRef = ref<InstanceType<typeof IgnoreRulesModal>>()
const searchInputRef = ref<HTMLInputElement | null>(null)

function clearSearch() {
  explorer.clearSearch()
}

defineExpose({ searchInputRef })

defineEmits<{
  (e: 'preview-file', filePath: string): void
  (e: 'build-context'): void
}>()

function handleContextMenu(node: FileNode, event: MouseEvent) {
  contextMenu.show(node, event)
}

function handleBreadcrumbNavigate(path: string) {
  fileStore.expandPath(path)
}

function handleAddSuggestedFiles(files: string[]) {
  const normalizedFiles = files
    .map(path => path.replace(/\\/g, '/'))
    .filter(path => !fileStore.selectedPaths.has(path))
  
  if (normalizedFiles.length > 0) {
    fileStore.selectMultiple(normalizedFiles)
  }
}

async function handleOpenInExplorer(_path: string) {
  if (projectStore.currentPath) {
    try {
      const runtime = await import('#wailsjs/runtime/runtime')
      runtime.BrowserOpenURL('file://' + projectStore.currentPath)
    } catch (error) {
      logger.error('Failed to open in explorer:', error)
      uiStore.addToast('Failed to open in file explorer', 'error')
    }
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

.file-explorer__header {
  flex-shrink: 0;
  border-bottom: 1px solid var(--border-default);
}

.file-explorer__search {
  flex-shrink: 0;
  padding: 0.5rem;
  border-bottom: 1px solid var(--border-default);
}

.file-explorer__tree {
  flex: 1 1 0;
  min-height: 0;
  overflow: hidden;
  padding: 0.5rem;
  display: flex;
  flex-direction: column;
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
