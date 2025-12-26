<template>
  <div class="context-content-area layout-fill layout-scroll" ref="contentContainer">
    <!-- In-place Skeleton Loader (AAA UX) -->
    <SkeletonLoader 
      v-if="isBuilding"
      :progress="buildProgress"
      :status-text="statusText"
    />

    <!-- Empty context with ID (warning state) -->
    <div
      v-else-if="(fileCount === 0 || totalSize === 0 || lineCount === 0) && contextId"
      class="flex items-center justify-center h-full"
    >
      <div class="text-center max-w-md mx-auto px-4">
        <div class="w-16 h-16 mx-auto mb-4 bg-amber-500/20 rounded-2xl flex items-center justify-center border border-amber-500/30">
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
              <div class="stat-card-value">{{ fileCount }}</div>
              <div class="stat-label">{{ t('context.files') }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-card-value">{{ lineCount }}</div>
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

    <!-- Error state -->
    <div v-else-if="error" class="flex items-center justify-center h-full">
      <div class="text-center max-w-md">
        <div class="w-16 h-16 mx-auto mb-4 bg-red-500/20 rounded-2xl flex items-center justify-center border border-red-500/30">
          <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          </svg>
        </div>
        <p class="text-lg font-semibold text-red-400 mb-2">{{ error }}</p>

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

    <!-- No context built state -->
    <div v-else-if="!hasContext" class="empty-state-enhanced h-full flex flex-col items-center justify-center">
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

    <!-- Main content with code view -->
    <div v-else class="code-editor context-content min-h-full">
      <!-- Template Preview Block (Sticky HUD) -->
      <TemplatePreviewBlock />

      <!-- Virtual Code View for performance with large files -->
      <VirtualCodeView
        v-if="lines?.length"
        ref="virtualCodeRef"
        :lines="lines"
        :highlighted-lines="highlightedLines"
        :search-query="searchQuery"
        :chunk-boundaries="chunkBoundaries"
        :output-format="outputFormat"
      />
      <div v-else-if="isLoading" class="text-center py-8">
        <svg class="animate-spin h-6 w-6 text-blue-500 mx-auto mb-2" xmlns="http://www.w3.org/2000/svg" fill="none"
          viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
        </svg>
        <p class="text-gray-400">{{ t('context.loading') }}</p>
      </div>
      <div v-else class="text-center py-8">
        <p v-if="contextId" class="text-gray-400">{{ t('context.loading') }}</p>
        <p v-else class="text-gray-400">{{ t('context.notBuilt') }}</p>
      </div>
    </div>
    
    <ExportModal ref="exportModalRef" />
  </div>
</template>

<script setup lang="ts">
import ExportModal from '@/components/ExportModal.vue'
import { useI18n } from '@/composables/useI18n'
import { TemplatePreviewBlock } from '@/features/templates'
import SkeletonLoader from './SkeletonLoader.vue'
import VirtualCodeView from './VirtualCodeView.vue'
import { ref } from 'vue'
import type { OutputFormat } from '@/stores/settings.store'

defineProps<{
  isBuilding: boolean
  buildProgress: number
  statusText: string
  fileCount: number
  lineCount: number
  totalSize: number
  contextId: string | null
  error: string | null
  hasContext: boolean
  isLoading: boolean
  lines: string[] | undefined
  highlightedLines: Set<number>
  searchQuery: string
  chunkBoundaries: Set<number> | undefined
  outputFormat: OutputFormat
}>()

const { t } = useI18n()
const virtualCodeRef = ref<InstanceType<typeof VirtualCodeView> | null>(null)
const exportModalRef = ref<InstanceType<typeof ExportModal> | null>(null)

// Expose methods for parent
function scrollToLine(lineNum: number) {
  virtualCodeRef.value?.scrollToLine(lineNum)
}

function openExportModal() {
  exportModalRef.value?.open()
}

defineExpose({
  scrollToLine,
  openExportModal
})
</script>

<style scoped>
.context-content-area {
  background: var(--bg-1);
}

.context-content {
  min-height: 100%;
  padding: 1rem;
}

.code-editor {
  @apply rounded-lg p-4 font-mono text-sm;
  background: var(--bg-1);
  color: var(--text-secondary);
}
</style>
