<template>
  <div class="workspace-container layout-grid-main-footer">
    <!-- Background decoration -->
    <div class="workspace-bg">
      <div class="workspace-glow workspace-glow-1"></div>
      <div class="workspace-glow workspace-glow-2"></div>
    </div>

    <!-- Main 3-column layout -->
    <div class="workspace-layout">
      <!-- Left Sidebar -->
      <div 
        ref="leftPanelRef"
        :style="{ width: `${leftWidth}px` }"
        class="workspace-panel workspace-panel-left"
      >
        <LeftSidebar @preview-file="handlePreviewFile" @build-context="handleBuildContext" />
      </div>

      <!-- Left Resize Handle -->
      <div 
        class="resize-handle"
        @mousedown="(e) => leftResize.onMouseDown(e)"
        :class="{ 'is-resizing': leftResize.isResizing.value }"
      >
        <div class="resize-handle-line"></div>
      </div>

      <!-- Center Panel -->
      <div class="workspace-center layout-fill layout-column layout-clip">
        <CenterWorkspace />
      </div>

      <!-- Right Sidebar Toggle -->
      <button
        @click="toggleRightSidebar"
        class="sidebar-toggle"
        :class="{ 'sidebar-toggle-open': showRightSidebar }"
        :style="showRightSidebar ? { right: `${rightWidth}px` } : {}"
        :title="t('sidebar.toggle')"
      >
        <svg class="w-4 h-4 transition-transform" :class="{ 'rotate-180': !showRightSidebar }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>

      <!-- Right Resize Handle -->
      <div 
        v-if="showRightSidebar"
        class="resize-handle"
        @mousedown="(e) => rightResize.onMouseDown(e)"
        :class="{ 'is-resizing': rightResize.isResizing.value }"
      >
        <div class="resize-handle-line"></div>
      </div>

      <!-- Right Sidebar -->
      <Transition name="slide-right">
        <div 
          v-if="showRightSidebar"
          ref="rightPanelRef"
          :style="{ width: `${rightWidth}px` }"
          class="workspace-panel workspace-panel-right"
        >
          <RightSidebar @open-export="handleOpenExport" />
        </div>
      </Transition>
    </div>

    <!-- Action Bar -->
    <ActionBar class="workspace-actionbar" @open-export="handleOpenExport" @reset-layout="resetPanelSizes" />

    <!-- Export Modal -->
    <ExportModal ref="exportModalRef" />
  </div>
</template>


<script setup lang="ts">
import ExportModal from '@/components/ExportModal.vue'
import { useI18n } from '@/composables/useI18n'
import { useLogger } from '@/composables/useLogger'
import { useResizablePanel } from '@/composables/useResizablePanel'
import { useContextStore } from '@/features/context'
import { useFileStore } from '@/features/files'
import { useTemplateStore, generateFileTree, detectLanguages } from '@/features/templates'
import { useProjectStore } from '@/stores/project.store'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { onMounted, onUnmounted, ref, watch } from 'vue'
import ActionBar from './ActionBar.vue'
import CenterWorkspace from './CenterWorkspace.vue'
import LeftSidebar from './LeftSidebar.vue'
import RightSidebar from './RightSidebar.vue'

const logger = useLogger('MainWorkspace')
const templateStore = useTemplateStore()
const projectStore = useProjectStore()
const { t } = useI18n()
const contextStore = useContextStore()
const fileStore = useFileStore()
const uiStore = useUIStore()
const settingsStore = useSettingsStore()
const exportModalRef = ref<InstanceType<typeof ExportModal> | null>(null)

// Right sidebar visibility
const showRightSidebar = ref(loadSidebarState())

function loadSidebarState(): boolean {
  try {
    const saved = localStorage.getItem('right-sidebar-visible')
    return saved !== 'false'
  } catch {
    return true
  }
}

function toggleRightSidebar() {
  showRightSidebar.value = !showRightSidebar.value
  try {
    localStorage.setItem('right-sidebar-visible', String(showRightSidebar.value))
  } catch {
    // Ignore localStorage errors
  }
}

function resetPanelSizes() {
  leftResize.resetToDefault()
  rightResize.resetToDefault()
  uiStore.addToast(t('workspace.layoutReset'), 'success')
}

// Panel resizing
const leftResize = useResizablePanel({
  minWidth: 280,
  maxWidth: 700,
  defaultWidth: 380,
  storageKey: 'workspace-left-width'
})
const leftWidth = leftResize.width

const rightResize = useResizablePanel({
  minWidth: 320,
  maxWidth: 700,
  defaultWidth: 380,
  storageKey: 'workspace-right-width',
  invertDirection: true // Тянем влево = увеличиваем ширину
})
// Panel refs are used in template via ref="leftPanelRef" and ref="rightPanelRef"
const leftPanelRef = leftResize.panelRef
const rightPanelRef = rightResize.panelRef
const rightWidth = rightResize.width

function handlePreviewFile(_filePath: string) {
  // Preview file functionality - handled by QuickLook
}

// Global keyboard shortcut handlers
const handleGlobalBuildContext = () => handleBuildContext()
const handleGlobalOpenExport = () => handleOpenExport()
const handleGlobalCopyContext = async () => {
  if (contextStore.hasContext && contextStore.contextId) {
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
      uiStore.addToast(t('toast.contextCopied'), 'success')
    } catch (error) {
      logger.error('Failed to copy context:', error)
      uiStore.addToast(t('toast.copyError'), 'error')
    }
  }
}

onMounted(() => {
  window.addEventListener('global-build-context', handleGlobalBuildContext)
  window.addEventListener('global-open-export', handleGlobalOpenExport)
  window.addEventListener('global-copy-context', handleGlobalCopyContext)
})

onUnmounted(() => {
  window.removeEventListener('global-build-context', handleGlobalBuildContext)
  window.removeEventListener('global-open-export', handleGlobalOpenExport)
  window.removeEventListener('global-copy-context', handleGlobalCopyContext)
})

// Watch for format/settings changes and rebuild context automatically
watch(
  () => [
    settingsStore.settings.context.outputFormat,
    settingsStore.settings.context.stripComments
  ],
  async ([newFormat, newStripComments], [oldFormat, oldStripComments]) => {
    // Only rebuild if context exists and settings actually changed
    if (!contextStore.hasContext || contextStore.isBuilding) return
    if (newFormat === oldFormat && newStripComments === oldStripComments) return
    
    // Get the files from current context summary
    const currentFiles = contextStore.summary?.files
    if (!currentFiles || currentFiles.length === 0) {
      // Fallback to selected files if no files in summary
      if (fileStore.selectedPaths.size === 0) return
    }
    
    const filePaths = currentFiles && currentFiles.length > 0 
      ? currentFiles 
      : Array.from(fileStore.selectedPaths)
    
    // Settings changed, rebuilding context silently
    
    try {
      const options = {
        maxTokens: settingsStore.settings.context.maxTokens,
        stripComments: settingsStore.settings.context.stripComments,
        includeTests: settingsStore.settings.context.includeTests,
        splitStrategy: settingsStore.settings.context.splitStrategy,
        outputFormat: settingsStore.settings.context.outputFormat,
        // Content optimization options
        excludeTests: settingsStore.settings.context.excludeTests,
        collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
        stripLicense: settingsStore.settings.context.stripLicense,
        compactDataFiles: settingsStore.settings.context.compactDataFiles,
        trimWhitespace: settingsStore.settings.context.trimWhitespace
      }
      
      await contextStore.buildContext(filePaths, options)
      uiStore.addToast(t('context.rebuilt'), 'success')
    } catch (error) {
      logger.error('Failed to rebuild context:', error)
    }
  }
)

async function handleBuildContext() {
  if (fileStore.selectedPaths.size === 0) {
    uiStore.addToast('Выберите файлы для построения контекста', 'warning')
    return
  }

  if (contextStore.isBuilding) return

  try {
    const filePaths = Array.from(fileStore.selectedPaths)
    const options = {
      maxTokens: settingsStore.settings.context.maxTokens,
      stripComments: settingsStore.settings.context.stripComments,
      includeTests: settingsStore.settings.context.includeTests,
      splitStrategy: settingsStore.settings.context.splitStrategy,
      outputFormat: settingsStore.settings.context.outputFormat,
      // Output options
      includeManifest: settingsStore.settings.context.includeManifest,
      includeLineNumbers: settingsStore.settings.context.includeLineNumbers,
      // Content optimization options
      excludeTests: settingsStore.settings.context.excludeTests,
      collapseEmptyLines: settingsStore.settings.context.collapseEmptyLines,
      stripLicense: settingsStore.settings.context.stripLicense,
      compactDataFiles: settingsStore.settings.context.compactDataFiles,
      trimWhitespace: settingsStore.settings.context.trimWhitespace
    }
    
    await contextStore.buildContext(filePaths, options)
    
    // Show success toast with copy action
    uiStore.addToast(t('toast.contextBuilt'), 'success', 5000, {
      label: t('context.copy'),
      icon: '📋',
      onClick: async () => {
        try {
          const content = await contextStore.getFullContextContent()
          await navigator.clipboard.writeText(content)
          uiStore.addToast(t('toast.contextCopied'), 'success')
        } catch {
          uiStore.addToast(t('toast.copyError'), 'error')
        }
      }
    })
  } catch (error) {
    logger.error('Failed to build context:', error)
    
    // Handle token limit exceeded error with detailed message
    if (error instanceof Error && error.message === 'TOKEN_LIMIT_EXCEEDED') {
      const storeError = contextStore.error
      if (storeError?.startsWith('TOKEN_LIMIT_EXCEEDED:')) {
        const parts = storeError.split(':')
        const actual = Number(parts[1])
        const limit = Number(parts[2])
        const actualK = Math.round(actual / 1000)
        const limitK = Math.round(limit / 1000)
        uiStore.addToast(
          t('error.tokenLimitExceeded', { actual: actualK, limit: limitK }),
          'error'
        )
      } else {
        uiStore.addToast(t('error.tokenLimitGeneric'), 'error')
      }
      return
    }
    
    const errorMsg = error instanceof Error ? error.message : 'Unknown error'
    uiStore.addToast(`${t('toast.contextError')}: ${errorMsg}`, 'error')
  }
}

function handleOpenExport() {
  if (!contextStore.hasContext) {
    uiStore.addToast('Сначала постройте контекст', 'warning')
    return
  }
  exportModalRef.value?.open()
}

// Expose refs used in template
defineExpose({ leftPanelRef, rightPanelRef })
</script>


<style scoped>
/* Workspace - использует layout-grid-main-footer из layout.css */
.workspace-container {
  position: relative;
  background: var(--bg-app);
}

/* Background */
.workspace-bg {
  @apply absolute inset-0 overflow-hidden pointer-events-none;
}

.workspace-glow {
  @apply absolute rounded-full blur-3xl opacity-30;
}

.workspace-glow-1 {
  @apply w-[600px] h-[600px] -top-64 -right-64;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.08) 0%, transparent 70%);
}

.workspace-glow-2 {
  @apply w-[500px] h-[500px] -bottom-48 -left-48;
  background: radial-gradient(circle, rgba(236, 72, 153, 0.06) 0%, transparent 70%);
}

/* Layout - 3-колоночный flex внутри grid-ячейки */
.workspace-layout {
  display: flex;
  overflow: hidden;
  position: relative;
  z-index: 10;
}

.workspace-panel {
  @apply backdrop-blur-sm overflow-hidden;
  background: var(--bg-panel-sidebar);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  align-self: stretch;
}

.workspace-panel-left {
  border-right: 1px solid var(--border-default);
}

.workspace-panel-right {
  border-left: 1px solid var(--border-default);
}

/* Центральная панель - использует layout-fill из layout.css */
.workspace-center {
  background: var(--bg-panel-center);
}

/* Resize Handle */
.resize-handle {
  @apply w-1 cursor-col-resize relative;
  @apply flex items-center justify-center;
  transition: all 150ms ease-out;
}

.resize-handle:hover,
.resize-handle.is-resizing {
  background: var(--accent-purple-bg);
}

.resize-handle-line {
  @apply w-0.5 h-16 rounded-full;
  background: var(--border-strong);
  transition: all 150ms ease-out;
}

.resize-handle:hover .resize-handle-line,
.resize-handle.is-resizing .resize-handle-line {
  @apply h-24;
  background: var(--accent-purple-border);
  box-shadow: 0 0 8px rgba(139, 92, 246, 0.4);
}

/* Sidebar Toggle */
.sidebar-toggle {
  @apply absolute top-1/2 -translate-y-1/2 z-20;
  @apply p-1.5 rounded-l-lg;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-right: 0;
  color: var(--text-muted);
  transition: color 150ms, background 150ms;
  right: 0;
}

.sidebar-toggle:hover {
  color: var(--text-primary);
  background: var(--bg-2);
}

.sidebar-toggle-open {
  border-radius: 0.5rem 0 0 0.5rem;
}

/* Action Bar */
.workspace-actionbar {
  @apply relative z-10;
  border-top: 1px solid var(--border-default);
  background: var(--bg-1);
  backdrop-filter: blur(12px);
}

/* Slide Animation */
.slide-right-enter-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-right-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 1, 1);
}

.slide-right-enter-from,
.slide-right-leave-to {
  opacity: 0;
  transform: translateX(100%);
}
</style>
