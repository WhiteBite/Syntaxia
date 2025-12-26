<template>
  <div class="context-panel-root layout-fill layout-column layout-clip"
       data-tour="context-preview"
       @dragover.prevent="dragDrop.handleDragOver"
       @dragleave="dragDrop.handleDragLeave"
       @drop.prevent="dragDrop.handleDrop">
    
    <!-- Drop Zone Overlay -->
    <div v-if="dragDrop.isDragging.value" class="drop-zone-overlay">
      <div class="drop-zone-content">
        <svg class="w-12 h-12 text-indigo-400 mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
            d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
        <p class="text-lg font-semibold text-white">{{ t('context.dropFiles') }}</p>
        <p class="text-sm text-gray-400">{{ t('context.dropFilesHint') }}</p>
      </div>
    </div>
    
    <!-- Header -->
    <ContextPanelHeader
      :has-context="contextStore.hasContext"
      :show-search="search.showSearch.value"
      :file-count="contextStore.fileCount"
      :line-count="contextStore.lineCount"
      :token-count="contextStore.tokenCount"
      :output-format="settingsStore.settings.context.outputFormat"
      @toggle-search="search.toggleSearch"
      @show-stats="showStatsPopover = true"
      @format-change="handleFormatChange"
    />

    <!-- Search Bar -->
    <ContextPanelToolbar
      :visible="search.showSearch.value"
      :search-query="search.searchQuery.value"
      :results-count="search.searchResults.value.length"
      :current-index="search.currentSearchIndex.value"
      @update:search-query="search.searchQuery.value = $event"
      @search-next="search.searchNext"
      @search-prev="search.searchPrev"
      @close="search.closeSearch"
    />

    <!-- Scrollable Content Area -->
    <ContextPanelContent
      ref="contentRef"
      :is-building="contextStore.isBuilding"
      :build-progress="contextStore.buildProgress"
      :status-text="buildStatusText"
      :file-count="contextStore.fileCount"
      :line-count="contextStore.lineCount"
      :total-size="contextStore.totalSize"
      :context-id="contextStore.contextId"
      :error="contextStore.error"
      :has-context="contextStore.hasContext"
      :is-loading="contextStore.isLoading"
      :lines="contextStore.currentChunk?.lines"
      :highlighted-lines="search.highlightedLinesSet.value"
      :search-query="search.searchQuery.value"
      :chunk-boundaries="chunking.showChunkHUD.value ? chunkBoundaries : undefined"
      :output-format="settingsStore.settings.context.outputFormat"
    />

    <!-- UNIFIED CONTEXT HUD BAR -->
    <ContextPanelFooter
      :visible="contextStore.hasContext"
      :show-chunk-nav="chunking.showChunkHUD.value"
      :current-chunk="chunking.currentChunk.value"
      :total-chunks="chunking.totalChunks.value"
      :copy-success="copySuccess"
      :is-chunk-copied="chunking.isChunkCopied"
      @clear="contextStore.clearContext"
      @export="handleExport"
      @copy="chunking.showChunkHUD.value ? handleCopyCurrentChunk() : handleCopyText()"
      @prev-chunk="goToPrevChunk"
      @next-chunk="goToNextChunk"
    />

    <!-- Stats Popover -->
    <StatsPopover 
      :visible="showStatsPopover" 
      @close="showStatsPopover = false" 
    />
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useTemplateStore, generateFileTree, detectLanguages } from '@/features/templates'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore, type OutputFormat } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, nextTick, ref, toRef } from 'vue'
import { useChunking } from '../composables/useChunking'
import { useContextSearch } from '../composables/useContextSearch'
import { useDragDrop } from '../composables/useDragDrop'
import { useContextStore } from '../model/context.store'

// Child components
import ContextPanelHeader from './ContextPanelHeader.vue'
import ContextPanelToolbar from './ContextPanelToolbar.vue'
import ContextPanelContent from './ContextPanelContent.vue'
import ContextPanelFooter from './ContextPanelFooter.vue'
import StatsPopover from './StatsPopover.vue'

const logger = useLogger('ContextPanel')
const contextStore = useContextStore()
const settingsStore = useSettingsStore()
const templateStore = useTemplateStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const { t } = useI18n()

// Refs
const contentRef = ref<InstanceType<typeof ContextPanelContent> | null>(null)
const copySuccess = ref(false)
const showStatsPopover = ref(false)

// Build status text for skeleton loader
const buildStatusText = computed(() => {
  const progress = contextStore.buildProgress
  if (progress < 20) return t('context.statusAnalyzing')
  if (progress < 50) return t('context.statusReading')
  if (progress < 80) return t('context.statusProcessing')
  if (progress < 95) return t('context.statusFormatting')
  return t('context.statusFinalizing')
})

// Composables
const dragDrop = useDragDrop()

const chunking = useChunking({
  tokenCount: toRef(() => contextStore.tokenCount),
  lineCount: toRef(() => contextStore.lineCount),
  fileCount: toRef(() => contextStore.fileCount),
  hasContext: toRef(() => contextStore.hasContext)
})

const search = useContextSearch({
  lines: toRef(() => contextStore.currentChunk?.lines),
  startLine: toRef(() => contextStore.currentChunk?.startLine ?? 0),
  scrollToLine: (lineNum: number) => contentRef.value?.scrollToLine(lineNum)
})

const chunkBoundaries = computed(() => {
  if (!chunking.showChunkHUD.value || chunking.chunks.value.length <= 1) return new Set<number>()
  return new Set(chunking.chunks.value.slice(1).map(c => c.startLine))
})

// Handle format change and rebuild context
async function handleFormatChange(format: OutputFormat) {
  settingsStore.updateContextSettings({ outputFormat: format })
  if (contextStore.contextId) {
    await contextStore.rebuildContext()
  }
}

// Chunk navigation with scroll
function goToPrevChunk() {
  chunking.prevChunk()
  scrollToCurrentChunk()
}

function goToNextChunk() {
  chunking.nextChunk()
  scrollToCurrentChunk()
}

function scrollToCurrentChunk() {
  nextTick(() => {
    const chunkInfo = chunking.currentChunkInfo.value
    if (!chunkInfo) return
    contentRef.value?.scrollToLine(chunkInfo.startLine)
  })
}

async function handleExport() {
  try {
    contentRef.value?.openExportModal()
  } catch (error) {
    logger.error('Failed to open export modal:', error)
  }
}

async function handleCopyText() {
  if (!contextStore.contextId) return
  
  try {
    const filesContent = await contextStore.getFullContextContent()
    
    let content: string
    if (settingsStore.settings.context.applyTemplateOnCopy && templateStore.activeTemplate) {
      const files = contextStore.summary?.files || []
      const templateContext = {
        fileTree: generateFileTree(files, projectStore.projectName),
        files: filesContent,
        task: templateStore.currentTask,
        userRules: templateStore.userRules,
        fileCount: contextStore.fileCount,
        tokenCount: contextStore.tokenCount,
        languages: detectLanguages(files),
        projectName: projectStore.projectName
      }
      content = templateStore.generatePrompt(templateContext)
    } else {
      content = filesContent
    }
    
    await navigator.clipboard.writeText(content)
    showCopySuccess()
    uiStore.addToast(t('toast.contextCopied'), 'success')
  } catch (error) {
    logger.error('Failed to copy context:', error)
    uiStore.addToast(t('toast.copyError'), 'error')
  }
}

async function handleCopyCurrentChunk() {
  if (!contextStore.contextId || !chunking.currentChunkInfo.value) return
  
  try {
    const chunkInfo = chunking.currentChunkInfo.value
    const lines = contextStore.currentChunk?.lines ?? []
    const chunkContent = lines.slice(chunkInfo.startLine, chunkInfo.endLine + 1).join('\n')
    
    await navigator.clipboard.writeText(chunkContent)
    showCopySuccess()
    chunking.markCopied(chunkInfo.index, true)
    uiStore.addToast(t('chunks.copied'), 'success')
  } catch (error) {
    logger.error('Failed to copy chunk:', error)
    uiStore.addToast(t('toast.copyError'), 'error')
  }
}

function showCopySuccess() {
  copySuccess.value = true
  setTimeout(() => { copySuccess.value = false }, 1500)
}
</script>

<style scoped>
.context-panel-root {
  width: 100%;
  background: var(--bg-1);
  position: relative;
}

.drop-zone-overlay {
  position: absolute;
  inset: 0;
  background: rgba(99, 102, 241, 0.1);
  backdrop-filter: blur(4px);
  border: 2px dashed rgba(99, 102, 241, 0.5);
  border-radius: 8px;
  z-index: 50;
  display: flex;
  align-items: center;
  justify-content: center;
}

.drop-zone-content {
  text-align: center;
}
</style>
