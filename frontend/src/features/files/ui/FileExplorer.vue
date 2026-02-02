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
              <BaseBadge variant="accent" class="cursor-help">
                {{ fileStore.selectedCount }}
              </BaseBadge>
              <!-- Tooltip with stats -->
              <div class="absolute right-0 top-full mt-2 hidden group-hover:block z-50">
                <BaseCard glass class="shadow-2xl">
                  <div class="text-xs space-y-1.5 whitespace-nowrap">
                    <div class="text-white font-semibold mb-2">Selection Stats</div>
                    <div class="text-gray-400">Files: <span class="text-white">{{ fileStore.selectedCount }}</span></div>
                    <div class="text-gray-400">Total: <span class="text-white">{{ explorer.totalFileCount.value }}</span> files</div>
                    <div class="text-gray-400">Progress: <span class="text-indigo-400">{{ explorer.selectionProgress.value }}%</span></div>
                    <div class="text-gray-400">Est. Size: <span class="text-white">{{ Math.round(fileStore.estimatedContextSize * 100) / 100 }}MB</span></div>
                    <div class="text-gray-400">Est. Tokens: <span class="text-emerald-400">~{{ Math.round(fileStore.estimatedTokenCount / 1000) }}K</span></div>
                  </div>
                </BaseCard>
              </div>
            </div>
            <span class="text-gray-400">{{ t('files.selected') }}</span>
            <!-- Clear selection button -->
            <BaseButton
              v-if="fileStore.selectedCount > 0"
              variant="ghost"
              size="xs"
              icon-only
              @click="fileStore.clearSelection"
              class="text-gray-500 hover:text-red-400"
              :title="t('files.clearSelection')"
            >
              <X class="w-3.5 h-3.5" />
            </BaseButton>
          </div>

          <!-- Settings Popover -->
          <SettingsPopover 
            @open-ignore-rules="ignoreRulesModalRef?.open()" 
            @settings-changed="explorer.handleSettingsChange"
          />
          
          <!-- Solo Expansion Mode Toggle -->
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            @click="fileStore.toggleSoloExpansionMode()"
            :class="{ 'text-indigo-400': fileStore.isSoloExpansionMode }"
            :title="t('files.soloExpansionMode')"
          >
            <List class="w-4 h-4" />
          </BaseButton>
          
          <!-- Selected Only Mode Toggle -->
          <BaseButton
            v-if="fileStore.selectedCount > 0"
            variant="ghost"
            size="sm"
            icon-only
            @click="fileStore.toggleSelectedOnlyMode()"
            :class="{ 'text-indigo-400': fileStore.isSelectedOnlyMode }"
            :title="t('files.selectedOnlyMode')"
          >
            <CheckSquare class="w-4 h-4" />
          </BaseButton>
          
          <!-- Zen Mode Toggle -->
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            @click="fileStore.toggleZenMode()"
            :class="{ 'text-indigo-400': fileStore.isZenMode }"
            :title="t('files.zenMode')"
          >
            <Eye class="w-4 h-4" />
          </BaseButton>

          <!-- Refresh Button -->
          <BaseButton
            variant="ghost"
            size="sm"
            icon-only
            @click="explorer.handleRefresh"
            :title="t('files.refresh')"
          >
            <BaseSpinner v-if="fileStore.isLoading" size="sm" />
            <RefreshCw v-else class="w-4 h-4" />
          </BaseButton>
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

    <!-- Folder Focus Mode Banner -->
    <div v-if="fileStore.focusedFolderPath" class="focus-mode-banner">
      <div class="focus-mode-content">
        <svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <circle cx="12" cy="12" r="3" stroke-width="2"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12c0-4.5 4-8 9-8s9 3.5 9 8-4 8-9 8-9-3.5-9-8z"/>
        </svg>
        <span class="focus-mode-label">{{ t('files.focusedOn', { folder: getFocusedFolderName() }) }}</span>
      </div>
      <BaseButton
        variant="ghost"
        size="sm"
        @click="fileStore.clearFolderFocus()"
        class="focus-mode-back-btn"
      >
        <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        {{ t('files.backToRoot') }}
      </BaseButton>
    </div>

    <!-- Selected Only Mode Banner -->
    <div v-if="fileStore.isSelectedOnlyMode" class="selected-only-banner">
      <div class="selected-only-content">
        <CheckSquare class="w-4 h-4 text-emerald-400" />
        <span class="selected-only-label">{{ t('files.selectedOnlyModeActive', { count: fileStore.selectedCount }) }}</span>
      </div>
      <BaseButton
        variant="ghost"
        size="sm"
        @click="fileStore.toggleSelectedOnlyMode()"
        class="selected-only-exit-btn"
      >
        <X class="w-4 h-4 mr-1.5" />
        {{ t('files.exitSelectedOnlyMode') }}
      </BaseButton>
    </div>

    <!-- Quick Filters -->
    <QuickFiltersBar />

    <!-- Presets Panel -->
    <PresetsPanel />

    <!-- Favorites Section -->
    <FavoritesPanel @select="handleFavoriteSelect" />

    <!-- Search -->
    <div class="file-explorer__search">
      <BaseInput 
        ref="searchInputRef"
        v-model="explorer.searchQuery.value" 
        :placeholder="t('files.searchShort')" 
        @input="explorer.handleSearch" 
        @keydown.escape="clearSearch"
      >
        <template #prefix>
          <SearchIcon class="w-4 h-4" />
        </template>
        <template #suffix v-if="explorer.searchQuery.value">
          <button 
            @click="clearSearch"
            class="p-1 rounded hover:bg-gray-700/50 text-gray-400 hover:text-white transition-colors"
            :title="t('files.clear')"
          >
            <X class="w-3.5 h-3.5" />
          </button>
        </template>
      </BaseInput>
    </div>


    <!-- File Tree - MAIN SCROLLABLE AREA -->
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
    <QuickLookModal v-model="explorer.quickLookVisible.value" :file-path="explorer.quickLookPath.value" @add-to-context="explorer.handleAddToContext" />
    <DependencyVisualizerModal
      v-model="showDependencyModal"
      :file-path="selectedFileForDeps"
    />

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
import { BaseBadge, BaseButton, BaseCard, BaseSpinner, BaseEmptyState } from '@/components/ui'
import { RefreshCw, Search as SearchIcon, X, Eye, CheckSquare, List } from 'lucide-vue-next'
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
import FavoritesPanel from './FavoritesPanel.vue'
import PresetsPanel from './PresetsPanel.vue'
import SkeletonFileTree from '@/components/SkeletonFileTree.vue'

const QuickLookModal = defineAsyncComponent(() => import('@/components/QuickLookModal.vue'))
const IgnoreRulesModal = defineAsyncComponent(() => import('./IgnoreRulesModal.vue'))
const DependencyVisualizerModal = defineAsyncComponent(() => import('./DependencyVisualizerModal.vue'))

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

function handleFavoriteSelect(_path: string) {
  // File is already added to context in FavoritesPanel
  // This handler can be used for additional actions if needed
}

function getFocusedFolderName(): string {
  if (!fileStore.focusedFolderPath) return ''
  const parts = fileStore.focusedFolderPath.split(/[/\\]/)
  return parts[parts.length - 1] || fileStore.focusedFolderPath
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

.focus-mode-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.1) 0%, rgba(139, 92, 246, 0.1) 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-left: 3px solid rgb(99, 102, 241);
  border-radius: 0.375rem;
  margin: 0.5rem;
  gap: 1rem;
}

.focus-mode-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
  min-width: 0;
}

.focus-mode-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #e5e7eb;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.focus-mode-back-btn {
  flex-shrink: 0;
  color: rgb(99, 102, 241);
  font-weight: 500;
}

.focus-mode-back-btn:hover {
  color: rgb(129, 140, 248);
  background: rgba(99, 102, 241, 0.1);
}

.selected-only-banner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.75rem 1rem;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(5, 150, 105, 0.1) 100%);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-left: 3px solid rgb(16, 185, 129);
  border-radius: 0.375rem;
  margin: 0.5rem;
  gap: 1rem;
}

.selected-only-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  flex: 1;
  min-width: 0;
}

.selected-only-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #d1fae5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.selected-only-exit-btn {
  flex-shrink: 0;
  color: rgb(16, 185, 129);
  font-weight: 500;
}

.selected-only-exit-btn:hover {
  color: rgb(52, 211, 153);
  background: rgba(16, 185, 129, 0.1);
}
</style>
