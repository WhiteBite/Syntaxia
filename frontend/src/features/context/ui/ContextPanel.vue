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
    <div class="panel-header-unified">
      <div class="panel-header-unified-title">
        <div class="section-icon section-icon-indigo">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
        </div>
        <span>{{ t('context.preview') }}</span>
      </div>

      <div v-if="contextStore.hasContext" class="flex items-center gap-1.5 ml-2 overflow-x-auto">
        <!-- Search toggle -->
        <button @click="showSearch = !showSearch" class="icon-btn"
          :class="{ 'text-indigo-400 bg-indigo-500/10': showSearch }" :title="t('context.search')">
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </button>

        <!-- Stats (clickable badges - open popover) -->
        <button 
          @click="showStatsPopover = true"
          class="hidden xl:flex items-center gap-1.5 flex-shrink-0 stats-chips-btn"
          :title="t('stats.clickToExpand')"
        >
          <span class="chip-unified chip-unified-accent">{{ contextStore.fileCount }} {{ t('context.files') }}</span>
          <span class="chip-unified chip-unified-accent">{{ contextStore.lineCount }} {{ t('context.lines') }}</span>
          <span class="chip-unified chip-unified-accent">{{ contextStore.tokenCount }} {{ t('context.tokens') }}</span>
        </button>

        <!-- Format Selector Dropdown -->
        <FormatDropdown
          :model-value="settingsStore.settings.context.outputFormat"
          @update:model-value="handleFormatChange"
        />
      </div>
    </div>

    <!-- Search Bar -->
    <div v-if="showSearch" class="search-bar">
      <div class="relative">
        <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input v-model="searchQuery" type="text" :placeholder="t('context.search')" class="input pl-8 pr-20 text-sm"
          @keyup.enter="searchNext" @keyup.escape="showSearch = false" ref="searchInput" />
        <div v-if="searchQuery && searchResults.length > 0" class="search-nav">
          <span class="search-count">{{ currentSearchIndex + 1 }}/{{ searchResults.length }}</span>
          <button @click="searchPrev" class="search-nav-btn">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
            </svg>
          </button>
          <button @click="searchNext" class="search-nav-btn">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
            </svg>
          </button>
        </div>
        <div v-else-if="searchQuery && searchResults.length === 0" class="search-no-results">
          <span>{{ t('context.noResults') }}</span>
        </div>
      </div>
    </div>

    <!-- Scrollable Content Area -->
    <div class="context-content-area layout-fill layout-scroll" ref="contentContainer">
      <!-- In-place Skeleton Loader (AAA UX) -->
      <SkeletonLoader 
        v-if="contextStore.isBuilding"
        :progress="contextStore.buildProgress"
        :status-text="buildStatusText"
      />

      <div
        v-else-if="(contextStore.fileCount === 0 || contextStore.totalSize === 0 || contextStore.lineCount === 0) && contextStore.contextId"
        class="flex items-center justify-center h-full">
        <div class="text-center max-w-md mx-auto px-4">
          <div
            class="w-16 h-16 mx-auto mb-4 bg-amber-500/20 rounded-2xl flex items-center justify-center border border-amber-500/30">
            <svg class="w-8 h-8 text-amber-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z">
              </path>
            </svg>
          </div>
          <p class="text-lg font-semibold text-amber-400 mb-4">{{ t('error.emptyContext') }}</p>

          <div class="context-stats mb-4">
            <div class="stats-grid">
              <div class="stat-card">
                <div class="stat-card-value">{{ contextStore.fileCount }}</div>
                <div class="stat-label">{{ t('context.files') }}</div>
              </div>
              <div class="stat-card">
                <div class="stat-card-value">{{ contextStore.lineCount }}</div>
                <div class="stat-label">{{ t('context.lines') }}</div>
              </div>
            </div>
          </div>

          <div class="info-box text-left">
            <p class="text-sm font-medium text-gray-300 mb-2">{{ t('error.suggestions') }}:</p>
            <ul class="text-sm text-gray-400 space-y-1.5">
              <li>• {{ t('error.checkFiles') }}</li>
              <li>• {{ t('error.checkPaths') }}</li>
              <li>• {{ t('error.tryRefresh') }}</li>
            </ul>
          </div>
        </div>
      </div>

      <div v-else-if="contextStore.error" class="flex items-center justify-center h-full">
        <div class="text-center max-w-md">
          <div
            class="w-16 h-16 mx-auto mb-4 bg-red-500/20 rounded-2xl flex items-center justify-center border border-red-500/30">
            <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <p class="text-lg font-semibold text-red-400 mb-2">{{ contextStore.error }}</p>

          <div class="info-box text-left mt-4">
            <p class="text-sm font-medium text-gray-300 mb-2">{{ t('error.suggestions') }}:</p>
            <ul class="text-sm text-gray-400 space-y-1.5">
              <li>• {{ t('error.checkFiles') }}</li>
              <li>• {{ t('error.checkPaths') }}</li>
              <li>• {{ t('error.tryRefresh') }}</li>
            </ul>
          </div>
        </div>
      </div>

      <div v-else-if="!contextStore.hasContext" class="empty-state-enhanced h-full flex flex-col items-center justify-center">
        <div class="empty-state-icon-glow mb-4">
          <svg class="w-8 h-8 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z">
            </path>
          </svg>
        </div>
        <p class="text-base font-semibold text-white mb-3">{{ t('context.notBuilt') }}</p>
        
        <!-- Step-by-step instructions -->
        <div class="text-left max-w-xs space-y-2 mb-6">
          <div class="flex items-center gap-3 text-sm">
            <span class="flex-shrink-0 w-5 h-5 rounded-full bg-indigo-500/20 text-indigo-400 text-xs font-bold flex items-center justify-center">1</span>
            <span class="text-gray-400">{{ t('context.step1') }}</span>
          </div>
          <div class="flex items-center gap-3 text-sm">
            <span class="flex-shrink-0 w-5 h-5 rounded-full bg-indigo-500/20 text-indigo-400 text-xs font-bold flex items-center justify-center">2</span>
            <span class="text-gray-400">{{ t('context.step2') }}</span>
          </div>
        </div>

        <!-- Hint arrows -->
        <div class="flex items-center gap-6 text-gray-400">
          <div class="flex items-center gap-2">
            <svg class="w-4 h-4 animate-[pulse_3s_ease-in-out_infinite]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
            </svg>
            <span class="text-xs">{{ t('context.selectHint') }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs">{{ t('context.chatHint') }}</span>
            <svg class="w-4 h-4 animate-[pulse_3s_ease-in-out_infinite]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
            </svg>
          </div>
        </div>
      </div>

      <div v-else class="code-editor context-content min-h-full">
        <!-- Template Preview Block (Sticky HUD) -->
        <TemplatePreviewBlock />
        

        
        <!-- Virtual Code View for performance with large files -->
        <VirtualCodeView
          v-if="contextStore.currentChunk?.lines?.length"
          ref="virtualCodeRef"
          :lines="contextStore.currentChunk.lines"
          :highlighted-lines="highlightedLinesSet"
          :search-query="searchQuery"
          :chunk-boundaries="showChunkHUD ? chunkBoundaries : undefined"
          :output-format="settingsStore.settings.context.outputFormat"
        />
        <div v-else-if="contextStore.isLoading" class="text-center py-8">
          <svg class="animate-spin h-6 w-6 text-blue-500 mx-auto mb-2" xmlns="http://www.w3.org/2000/svg" fill="none"
            viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
          </svg>
          <p class="text-gray-400">{{ t('context.loading') }}</p>
        </div>
        <div v-else class="text-center py-8">
          <p v-if="contextStore.contextId" class="text-gray-400">{{ t('context.loading') }}</p>
          <p v-else class="text-gray-400">{{ t('context.notBuilt') }}</p>
        </div>
      </div>
      <ExportModal ref="exportModalRef" />
    </div>

    <!-- UNIFIED CONTEXT HUD BAR -->
    <Transition name="hud">
      <div v-if="contextStore.hasContext" class="unified-hud">
        <!-- Tools Section -->
        <div class="hud-section hud-tools">
          <Tooltip :text="t('context.clear')" position="top">
            <button @click="contextStore.clearContext" class="hud-tool-btn hud-tool-btn--danger">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </Tooltip>
          <Tooltip :text="t('context.export')" position="top">
            <button @click="handleExport" class="hud-tool-btn">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </button>
          </Tooltip>
        </div>

        <!-- Divider -->
        <div class="hud-divider" />

        <!-- Navigation Section (only if chunking) -->
        <template v-if="showChunkHUD">
          <div class="hud-section hud-nav">
            <button 
              class="hud-nav-btn" 
              :disabled="chunking.currentChunk.value <= 1"
              @click="goToPrevChunk"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <div class="hud-counter-wrap">
              <div class="hud-counter">
                <span class="hud-counter-current">{{ chunking.currentChunk.value }}</span>
                <span class="hud-counter-sep">/</span>
                <span class="hud-counter-total">{{ chunking.totalChunks.value }}</span>
              </div>
              <!-- Progress dots -->
              <div class="hud-progress-dots">
                <span 
                  v-for="i in chunking.totalChunks.value" 
                  :key="i"
                  class="hud-dot"
                  :class="{ 
                    active: i === chunking.currentChunk.value,
                    copied: chunking.isChunkCopied(i - 1)
                  }"
                />
              </div>
            </div>
            <button 
              class="hud-nav-btn" 
              :disabled="chunking.currentChunk.value >= chunking.totalChunks.value"
              @click="goToNextChunk"
            >
              <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>
          <div class="hud-divider" />
        </template>

        <!-- Main Action Button -->
        <button 
          class="hud-action-btn"
          :class="{ 'hud-action-btn--success': copySuccess }"
          @click="showChunkHUD ? handleCopyCurrentChunk() : handleCopyText()"
        >
          <svg v-if="!copySuccess" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <span>{{ copyButtonText }}</span>
        </button>
      </div>
    </Transition>

    <!-- Stats Popover -->
    <StatsPopover 
      :visible="showStatsPopover" 
      @close="showStatsPopover = false" 
    />
  </div>
</template>

<script setup lang="ts">
import ExportModal from '@/components/ExportModal.vue'
import Tooltip from '@/components/ui/Tooltip.vue'
import FormatDropdown from '@/components/workspace/sidebar/export/FormatDropdown.vue'
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { TemplatePreviewBlock, useTemplateStore, generateFileTree, detectLanguages } from '@/features/templates'
import SkeletonLoader from './SkeletonLoader.vue'
import StatsPopover from './StatsPopover.vue'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore, type OutputFormat } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, nextTick, ref, watch } from 'vue'
import { useChunking } from '../composables/useChunking'
import { formatContextSize } from '../lib/context-utils'
import { useContextStore } from '../model/context.store'
import VirtualCodeView from './VirtualCodeView.vue'

const logger = useLogger('ContextPanel')
const contextStore = useContextStore()
const settingsStore = useSettingsStore()
const templateStore = useTemplateStore()
const projectStore = useProjectStore()
const uiStore = useUIStore()
const { t } = useI18n()
const exportModalRef = ref<InstanceType<typeof ExportModal> | null>(null)

// Copy state for unified HUD
const copySuccess = ref(false)
const copyButtonText = computed(() => {
  if (copySuccess.value) return '✓'
  if (showChunkHUD.value) return `#${chunking.currentChunk.value}`
  return t('context.copy')
})

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

// Current chunk line range for indicator (reserved for future HUD display)
const _currentChunkLineRange = computed(() => {
  const info = chunking.currentChunkInfo.value
  if (!info) return ''
  return `${info.startLine + 1}–${info.endLine + 1}`
})
void _currentChunkLineRange // suppress unused warning

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
    
    // Calculate number of chunks based on total tokens
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
    
    // Estimate chunk boundaries based on token distribution
    const tokensPerChunk = Math.ceil(totalTokens / numChunks)
    const linesPerChunk = Math.ceil(totalLines / numChunks)
    
    const chunks: { index: number; startLine: number; endLine: number; tokenCount: number }[] = []
    
    for (let i = 0; i < numChunks; i++) {
      const startLine = i * linesPerChunk
      const endLine = Math.min((i + 1) * linesPerChunk - 1, totalLines - 1)
      const chunkTokens = i === numChunks - 1 
        ? totalTokens - (tokensPerChunk * (numChunks - 1))
        : tokensPerChunk
      
      chunks.push({
        index: i,
        startLine,
        endLine,
        tokenCount: chunkTokens
      })
    }
    
    chunking.setChunks(chunks)
  },
  { immediate: true }
)

// Search state
const showSearch = ref(false)
const showStatsPopover = ref(false)
const searchQuery = ref('')
const searchResults = ref<number[]>([]) // Line numbers with matches
const currentSearchIndex = ref(0)
const searchInput = ref<HTMLInputElement | null>(null)
const lineRefs = ref<Map<number, HTMLElement>>(new Map())
const virtualCodeRef = ref<InstanceType<typeof VirtualCodeView> | null>(null)

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
  // Rebuild context with new format if we have one
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
  if (show) {
    nextTick(() => searchInput.value?.focus())
  } else {
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
    
    // Find the line element for the start of this chunk
    const lineEl = lineRefs.value.get(chunkInfo.startLine)
    if (lineEl) {
      lineEl.scrollIntoView({ behavior: 'smooth', block: 'start' })
    } else {
      // Fallback: scroll content container to top
      const container = document.querySelector('.context-content-area')
      if (container) {
        const scrollRatio = chunkInfo.startLine / (contextStore.lineCount || 1)
        container.scrollTop = container.scrollHeight * scrollRatio
      }
    }
  })
}

function scrollToLine(lineNum: number) {
  nextTick(() => {
    // Use virtual scroller if available
    if (virtualCodeRef.value) {
      virtualCodeRef.value.scrollToLine(lineNum)
      return
    }
    // Fallback to old method
    const el = lineRefs.value.get(lineNum)
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'center' })
    }
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

// Debug logging to track state changes (dev only)
if (import.meta.env.DEV) {
  watch(() => [contextStore.contextId, contextStore.fileCount, contextStore.totalSize, contextStore.lineCount, contextStore.error, contextStore.isBuilding],
    ([contextId, fileCount, totalSize, lineCount, error, isBuilding]) => {
      logger.debug('State changed:', { contextId, fileCount, totalSize, lineCount, error, isBuilding })
    },
    { immediate: true }
  )
}

// Helper functions available for future use
void formatContextSize // suppress unused import warning

async function handleExport() {
  try {
    if (exportModalRef.value) {
      exportModalRef.value.open()
    }
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
    chunking.markCopied(chunkInfo.index, true) // Auto-advance
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
  
  // Check if dragging files
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
  
  // Handle file paths from internal drag (text/plain with file paths)
  const textData = e.dataTransfer.getData('text/plain')
  if (textData) {
    const paths = textData.split('\n').filter(p => p.trim())
    if (paths.length > 0) {
      logger.debug('Dropped file paths:', paths.length)
      // Emit event to add files to selection and build context
      window.dispatchEvent(new CustomEvent('add-files-to-context', { detail: { paths } }))
      return
    }
  }
  
  // Handle external file drops (from OS file manager)
  if (e.dataTransfer.files.length > 0) {
    const filePaths = Array.from(e.dataTransfer.files).map(f => f.name)
    logger.debug('Dropped external files:', filePaths.length)
    // Note: External files need special handling - for now just log
  }
}
</script>


<style scoped>
/* Stats Chips Button */
.stats-chips-btn {
  background: none;
  border: none;
  padding: 2px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.stats-chips-btn:hover {
  background: rgba(139, 92, 246, 0.1);
}

.stats-chips-btn:hover .chip-unified {
  border-color: rgba(139, 92, 246, 0.4);
}

/* Search Bar */
.search-bar {
  @apply p-2;
  background: var(--bg-1);
  border-bottom: 1px solid var(--border-default);
}

.search-icon {
  @apply absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4;
  color: var(--text-muted);
}

.search-nav {
  @apply absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1;
}

.search-count {
  @apply text-xs;
  color: var(--text-muted);
}

.search-nav-btn {
  @apply p-1 rounded;
  color: var(--text-muted);
  transition: all 150ms ease-out;
}

.search-nav-btn:hover {
  background: var(--bg-2);
  color: var(--text-primary);
}

.search-no-results {
  @apply absolute right-2 top-1/2 -translate-y-1/2 text-xs;
  color: var(--text-subtle);
}

/* Code Editor */
.code-editor {
  @apply rounded-lg p-4 font-mono text-sm;
  background: var(--bg-1);
  color: var(--text-secondary);
}

.code-line {
  @apply px-2 py-0.5 rounded;
  transition: background-color 150ms ease-out;
}

.code-line:hover:not(.code-line-highlight) {
  background: var(--bg-2);
}

.code-line-highlight {
  background: rgba(250, 204, 21, 0.15);
  border-left: 2px solid var(--color-warning);
}

.line-number {
  @apply mr-4 select-none;
  color: var(--text-subtle);
}

/* Root container - layout классы в template */
.context-panel-root {
  width: 100%;
  background: var(--bg-1);
  position: relative;
}

/* Scrollable Content Area - layout-fill layout-scroll в template */
.context-content-area {
  background: var(--bg-1);
}

/* Контент внутри области скролла */
.context-content {
  min-height: 100%; /* Гарантирует заполнение всей высоты контейнера */
  padding: 1rem;
}

/* Syntax Highlighting */
:deep(.syntax-tag) {
  color: #f472b6; /* pink-400 */
}

:deep(.syntax-attr) {
  color: #a78bfa; /* violet-400 */
}

:deep(.syntax-string) {
  color: #4ade80; /* green-400 */
}

:deep(.syntax-punct) {
  color: #94a3b8; /* slate-400 */
}

:deep(.syntax-heading) {
  color: #60a5fa; /* blue-400 */
  font-weight: 600;
}

:deep(.syntax-fence) {
  color: #f59e0b; /* amber-500 */
}

:deep(.syntax-key) {
  color: #a78bfa; /* violet-400 */
}

:deep(.syntax-value) {
  color: #f472b6; /* pink-400 */
}

:deep(.syntax-separator) {
  color: #60a5fa; /* blue-400 */
  font-weight: 500;
}

/* UNIFIED HUD BAR - Glassmorphism Island */
.unified-hud {
  position: absolute;
  bottom: 16px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: rgba(15, 17, 26, 0.85);
  backdrop-filter: blur(24px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  box-shadow: 
    0 16px 48px rgba(0, 0, 0, 0.6),
    0 4px 16px rgba(0, 0, 0, 0.4);
  z-index: 20;
}

.hud-section {
  display: flex;
  align-items: center;
  gap: 4px;
}

.hud-divider {
  width: 1px;
  height: 20px;
  background: rgba(255, 255, 255, 0.1);
  margin: 0 4px;
}

/* Tool buttons */
.hud-tool-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: 8px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.hud-tool-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.hud-tool-btn--danger:hover {
  background: rgba(239, 68, 68, 0.15);
  color: #f87171;
}

/* Navigation */
.hud-nav-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  background: rgba(255, 255, 255, 0.05);
  border: none;
  border-radius: 6px;
  color: #9ca3af;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.hud-nav-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.hud-nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.hud-counter-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 0 8px;
}

.hud-counter {
  display: flex;
  align-items: baseline;
  gap: 2px;
  font-family: ui-monospace, monospace;
}

.hud-progress-dots {
  display: flex;
  gap: 4px;
}

.hud-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  transition: all 0.15s ease-out;
}

.hud-dot.active {
  background: #8b5cf6;
  box-shadow: 0 0 8px rgba(139, 92, 246, 0.6);
}

.hud-dot.copied {
  background: #10b981;
}

.hud-counter-current {
  font-size: 14px;
  font-weight: 700;
  color: #e5e7eb;
}

.hud-counter-sep {
  font-size: 11px;
  color: #4b5563;
}

.hud-counter-total {
  font-size: 11px;
  color: #6b7280;
}

/* Main action button - Most prominent */
.hud-action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease-out;
  box-shadow: 
    0 4px 16px rgba(139, 92, 246, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

.hud-action-btn:hover {
  transform: translateY(-1px);
  box-shadow: 
    0 6px 24px rgba(139, 92, 246, 0.5),
    inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

.hud-action-btn--success {
  background: linear-gradient(135deg, #059669 0%, #10b981 100%);
  box-shadow: 
    0 4px 16px rgba(16, 185, 129, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.15);
}

/* HUD Transition */
.hud-enter-active,
.hud-leave-active {
  transition: all 0.25s ease-out;
}

.hud-enter-from,
.hud-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(16px);
}

.fab-icon--primary:hover {
  box-shadow: 0 4px 12px rgba(147, 51, 234, 0.6);
  transform: translateY(-2px);
}
</style>
