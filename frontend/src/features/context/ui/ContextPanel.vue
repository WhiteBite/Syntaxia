<template>
  <div class="context-panel-root layout-fill layout-column layout-clip"
       data-tour="context-preview"
       @dragover.prevent="handleDragOver"
       @dragleave="handleDragLeave"
       @drop.prevent="handleDrop">
    
    <!-- Drop Zone Overlay -->
    <div v-if="isDragging" class="drop-zone-overlay">
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
      :show-search="showSearch"
      :file-count="contextStore.fileCount"
      :line-count="contextStore.lineCount"
      :token-count="contextStore.tokenCount"
      :output-format="settingsStore.settings.context.outputFormat"
      @toggle-search="showSearch = !showSearch"
      @show-stats="showStatsPopover = true"
      @format-change="handleFormatChange"
    />

    <!-- Search Bar -->
    <ContextPanelToolbar
      :visible="showSearch"
      :search-query="searchQuery"
      :results-count="searchResults.length"
      :current-index="currentSearchIndex"
      @update:search-query="searchQuery = $event"
      @search-next="searchNext"
      @search-prev="searchPrev"
      @close="showSearch = false"
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
      :highlighted-lines="highlightedLinesSet"
      :search-query="searchQuery"
      :chunk-boundaries="showChunkHUD ? chunkBoundaries : undefined"
      :output-format="settingsStore.settings.context.outputFormat"
    />

    <!-- UNIFIED CONTEXT HUD BAR -->
    <ContextPanelFooter
      :visible="contextStore.hasContext"
      :show-chunk-nav="showChunkHUD"
      :current-chunk="chunking.currentChunk.value"
      :total-chunks="chunking.totalChunks.value"
      :copy-success="copySuccess"
      :is-chunk-copied="chunking.isChunkCopied"
      @clear="contextStore.clearContext"
      @export="handleExport"
      @copy="showChunkHUD ? handleCopyCurrentChunk() : handleCopyText()"
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
import { computed, nextTick, ref, watch } from 'vue'
import { useChunking } from '../composables/useChunking'
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

// Copy state
const copySuccess = ref(false)

// Build status text for skeleton loader
const buildStatusText = computed(() => {
  const progress = contextStore.buildProgress
  if (progress < 20) return t('context.statusAnalyzing')
  if (progress < 50) return t('context.statusReading')
  if (progress < 80) return t('context.statusProcessing')
  if (progress < 95) return t('context.statusFormatting')
  return t('context.statusFinalizing')
})

// Chunking
const chunking = useChunking()

// Show HUD when chunking is enabled and total tokens exceed chunk limit
const showChunkHUD = computed(() => {
  if (!contextStore.hasContext) return false
  if (!settingsStore.settings.context.enableAutoSplit) return false
  const totalTokens = contextStore.tokenCount
  const maxPerChunk = settingsStore.settings.context.maxTokensPerChunk
  return totalTokens > maxPerChunk
})

// Auto-calculate chunks when context changes
watch(
  () => [
    contextStore.tokenCount, 
    contextStore.lineCount,
    contextStore.fileCount,
    settingsStore.settings.context.enableAutoSplit,
    settingsStore.settings.context.maxTokensPerChunk,
    settingsStore.settings.context.splitStrategy
  ],
  () => {
    if (!settingsStore.settings.context.enableAutoSplit) {
      chunking.setChunks([])
      return
    }
    
    const totalTokens = contextStore.tokenCount
    const totalLines = contextStore.lineCount
    const maxPerChunk = settingsStore.settings.context.maxTokensPerChunk
    
    if (totalTokens <= 0 || totalLines <= 0) {
      chunking.setChunks([])
      return
    }
    
    const numChunks = Math.ceil(totalTokens / maxPerChunk)
    
    if (numChunks <= 1) {
      chunking.setChunks([{
        index: 0,
        startLine: 0,
        endLine: totalLines - 1,
        tokenCount: totalTokens
      }])
      return
    }
    
    const tokensPerChunk = Math.ceil(totalTokens / numChunks)
    const linesPerChunk = Math.ceil(totalLines / numChunks)
    
    const chunks: { index: number; startLine: number; endLine: number; tokenCount: number }[] = []
    
    for (let i = 0; i < numChunks; i++) {
      const startLine = i * linesPerChunk
      const endLine = Math.min((i + 1) * linesPerChunk - 1, totalLines - 1)
      const chunkTokens = i === numChunks - 1 
        ? totalTokens - (tokensPerChunk * (numChunks - 1))
        : tokensPerChunk
      
      chunks.push({ index: i, startLine, endLine, tokenCount: chunkTokens })
    }
    
    chunking.setChunks(chunks)
  },
  { immediate: true }
)

// Search state
const showSearch = ref(false)
const showStatsPopover = ref(false)
const searchQuery = ref('')
const searchResults = ref<number[]>([])
const currentSearchIndex = ref(0)

// Computed for VirtualCodeView
const highlightedLinesSet = computed(() => {
  if (searchResults.value.length === 0) return new Set<number>()
  const currentLine = searchResults.value[currentSearchIndex.value]
  return currentLine !== undefined ? new Set([currentLine]) : new Set<number>()
})

const chunkBoundaries = computed(() => {
  if (!showChunkHUD.value || chunking.chunks.value.length <= 1) return new Set<number>()
  return new Set(chunking.chunks.value.slice(1).map(c => c.startLine))
})

// Drag & drop state
const isDragging = ref(false)
let dragCounter = 0

// Handle format change and rebuild context
async function handleFormatChange(format: OutputFormat) {
  settingsStore.updateContextSettings({ outputFormat: format })
  if (contextStore.contextId) {
    await contextStore.rebuildContext()
  }
}

// Search functionality
watch(searchQuery, (query) => {
  if (!query || !contextStore.currentChunk?.lines) {
    searchResults.value = []
    currentSearchIndex.value = 0
    return
  }

  const results: number[] = []
  const lowerQuery = query.toLowerCase()

  contextStore.currentChunk.lines.forEach((line, index) => {
    if (line.toLowerCase().includes(lowerQuery)) {
      results.push(contextStore.currentChunk!.startLine + index)
    }
  })

  searchResults.value = results
  currentSearchIndex.value = 0

  if (results.length > 0) {
    scrollToLine(results[0])
  }
})

watch(showSearch, (show) => {
  if (!show) {
    searchQuery.value = ''
  }
})

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

function scrollToLine(lineNum: number) {
  nextTick(() => {
    contentRef.value?.scrollToLine(lineNum)
  })
}

function searchNext() {
  if (searchResults.value.length === 0) return
  currentSearchIndex.value = (currentSearchIndex.value + 1) % searchResults.value.length
  scrollToLine(searchResults.value[currentSearchIndex.value])
}

function searchPrev() {
  if (searchResults.value.length === 0) return
  currentSearchIndex.value = (currentSearchIndex.value - 1 + searchResults.value.length) % searchResults.value.length
  scrollToLine(searchResults.value[currentSearchIndex.value])
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

// Drag & drop handlers
function handleDragOver(e: DragEvent) {
  dragCounter++
  isDragging.value = true
  
  if (e.dataTransfer?.types.includes('Files') || e.dataTransfer?.types.includes('text/plain')) {
    e.dataTransfer.dropEffect = 'copy'
  }
}

function handleDragLeave() {
  dragCounter--
  if (dragCounter <= 0) {
    dragCounter = 0
    isDragging.value = false
  }
}

async function handleDrop(e: DragEvent) {
  isDragging.value = false
  dragCounter = 0
  
  if (!e.dataTransfer) return
  
  const textData = e.dataTransfer.getData('text/plain')
  if (textData) {
    const paths = textData.split('\n').filter(p => p.trim())
    if (paths.length > 0) {
      logger.debug('Dropped file paths:', paths.length)
      window.dispatchEvent(new CustomEvent('add-files-to-context', { detail: { paths } }))
    }
  }
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
