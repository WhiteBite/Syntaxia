<template>
  <div 
    class="virtual-tree-wrapper" 
    :class="{ 'tree-row-drag-selecting': isDraggingSelection }"
    @keydown="handleKeyDown"
  >
    <div v-if="currentFolderPath && flattenedNodes.length > 20" class="tree-sticky-breadcrumb">
      <FolderIcon class="w-4 h-4" />
      <span>{{ currentFolderPath }}</span>
    </div>
    <RecycleScroller
      ref="scrollerRef"
      class="virtual-tree-scroller"
      :class="{ 'is-scrolling': isScrolling }"
      :items="flattenedNodes"
      :item-size="rowHeight"
      key-field="id"
      v-slot="{ item }"
      @scroll.passive="handleScroll"
    >
      <VirtualTreeRow
        :ref="(el) => setRowRef(item.id, el)"
        :item="item"
        :compact-mode="compactMode"
        :is-selected="isNodeSelected(item.node)"
        :is-focused="item.id === fileStore.focusedPath"
        :checkbox-state="getCheckboxState(item.node)"
        :file-count="getFileCount(item.node)"
        :selected-file-count="getSelectedFileCount(item.node)"
        :selected-tokens="getSelectedTokens(item.node)"
        :allow-select-binary="allowSelectBinary"
        :is-dragging-selection="isDraggingSelection"
        @toggle-select="handleToggleSelect"
        @toggle-expand="$emit('toggle-expand', $event)"
        @contextmenu="(node, event) => $emit('contextmenu', node, event)"
        @quicklook="$emit('quicklook', $event)"
        @select-related="handleSelectRelated"
        @add-dependency="handleAddDependency"
        @checkbox-mousedown="handleCheckboxMouseDown"
        @row-mouseenter="handleRowMouseEnter"
      />
    </RecycleScroller>
    <div v-if="flattenedNodes.length === 0" class="empty-state">
      <p class="empty-state-text">{{ t('files.noFiles') }}</p>
    </div>
  </div>

</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useVirtualTree } from '@/composables/useVirtualTree'
import { useSettingsStore } from '@/stores/settings.store'
import { useUIStore } from '@/stores/ui.store'
import { computed, toRef, ref } from 'vue'
import { FolderIcon } from '@heroicons/vue/24/outline'

import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { useFileStore, type FileNode } from '../model/file.store'
import VirtualTreeRow from './VirtualTreeRow.vue'

const { t } = useI18n()
const fileStore = useFileStore()
const settingsStore = useSettingsStore()
const uiStore = useUIStore()

const rowHeight = computed(() => 26 * settingsStore.settings.uiScale)


interface Props {
  nodes: FileNode[]
  compactMode?: boolean
  allowSelectBinary?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compactMode: false,
  allowSelectBinary: false,
})

const emit = defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'toggle-expand', path: string): void
  (e: 'contextmenu', node: FileNode, event: MouseEvent): void
  (e: 'quicklook', path: string): void
  (e: 'select-related', path: string): void
}>()

const nodesRef = toRef(props, 'nodes')
const { flattenedVisibleNodes } = useVirtualTree({ 
  nodes: nodesRef,
  isSelectedOnlyMode: computed(() => fileStore.isSelectedOnlyMode),
  selectedPaths: computed(() => fileStore.selectedPaths),
  rootPath: computed(() => fileStore.rootPath),
  focusedFolderPath: computed(() => fileStore.focusedFolderPath)
})

const flattenedNodes = computed(() => flattenedVisibleNodes.value)

// Drag-to-select state
const isDraggingSelection = ref(false)
const dragInitialAction = ref<'select' | 'deselect'>('select')
const draggedPaths = new Set<string>()

// Scroll performance optimization
const isScrolling = ref(false)
const currentFolderPath = ref('')
let scrollTimeout: number | null = null

// Keyboard navigation: refs to row elements
const rowRefs = new Map<string, InstanceType<typeof VirtualTreeRow>>()
const scrollerRef = ref<InstanceType<typeof RecycleScroller> | null>(null)

function setRowRef(id: string, el: unknown) {
  if (el) {
    rowRefs.set(id, el as InstanceType<typeof VirtualTreeRow>)
  } else {
    rowRefs.delete(id)
  }
}

function handleScroll(event: Event) {
  isScrolling.value = true
  if (scrollTimeout) clearTimeout(scrollTimeout)
  scrollTimeout = setTimeout(() => {
    isScrolling.value = false
  }, 150) as unknown as number
  
  // Calculate current folder path
  const scroller = event.target as HTMLElement
  const scrollTop = scroller.scrollTop
  const firstVisibleIndex = Math.floor(scrollTop / rowHeight.value)
  const firstVisible = flattenedNodes.value[firstVisibleIndex]
  
  if (firstVisible) {
    const path = firstVisible.node.path
    const rootPath = fileStore.rootPath || ''
    
    // Get relative path from project root
    let relativePath = path
    if (rootPath && path.startsWith(rootPath)) {
      relativePath = path.slice(rootPath.length).replace(/^[/\\]/, '')
    }
    
    // Split path and get folder path (remove filename if it's a file)
    const parts = relativePath.split(/[/\\]/).filter(Boolean)
    if (!firstVisible.node.isDir && parts.length > 0) {
      parts.pop() // Remove filename
    }
    
    // Show last 3 levels maximum
    const displayParts = parts.slice(-3)
    currentFolderPath.value = displayParts.length > 0 ? displayParts.join(' › ') : fileStore.projectName || 'Root'
  }
}

function handleToggleSelect(payload: { path: string, shiftKey: boolean }) {
  const { path, shiftKey } = payload
  
  if (shiftKey && fileStore.lastSelectedPath && fileStore.lastSelectedPath !== path) {
    const lastIndex = flattenedNodes.value.findIndex(n => n.id === fileStore.lastSelectedPath)
    const currentIndex = flattenedNodes.value.findIndex(n => n.id === path)
    
    if (lastIndex !== -1 && currentIndex !== -1) {
      const start = Math.min(lastIndex, currentIndex)
      const end = Math.max(lastIndex, currentIndex)
      
      const pathsToSelect = flattenedNodes.value
        .slice(start, end + 1)
        .filter(n => !n.node.isDir)
        .map(n => n.id)
      
      fileStore.selectMultiple(pathsToSelect)
      fileStore.lastSelectedPath = path
      return
    }
  }
  
  fileStore.toggleSelect(path)
}

// Selection helpers - computed at parent level for better performance
function isNodeSelected(node: FileNode): boolean {
  if (!node.isDir) {
    return fileStore.selectedPaths.has(node.path)
  }
  // For directories, check if any files inside are selected
  const allFiles = fileStore.getAllFilesInNode(node)
  return allFiles.some(filePath => fileStore.selectedPaths.has(filePath))
}

function getCheckboxState(node: FileNode): 'none' | 'partial' | 'full' {
  if (!node.isDir) {
    return fileStore.selectedPaths.has(node.path) ? 'full' : 'none'
  }
  if (fileStore.selectedCount === 0) return 'none'
  
  const allFiles = fileStore.getAllFilesInNode(node)
  if (allFiles.length === 0) return 'none'
  
  let selectedCount = 0
  for (const filePath of allFiles) {
    if (fileStore.selectedPaths.has(filePath)) selectedCount++
  }
  
  if (selectedCount === 0) return 'none'
  if (selectedCount === allFiles.length) return 'full'
  return 'partial'
}

function getFileCount(node: FileNode): number {
  if (!node.isDir) return 0
  return fileStore.getAllFilesInNode(node).length
}

// Build a map of path -> size for quick lookup
const fileSizeMap = computed(() => {
  const map = new Map<string, number>()
  for (const item of flattenedVisibleNodes.value) {
    if (!item.node.isDir && item.node.size) {
      map.set(item.node.path, item.node.size)
    }
  }
  return map
})

// Get total tokens of selected files inside a folder (for bubble-up indicator)
function getSelectedTokens(node: FileNode): number {
  if (!node.isDir) {
    // For files, return their own token count if selected
    if (fileStore.selectedPaths.has(node.path) && node.size) {
      return Math.round(node.size / 4)
    }
    return 0
  }
  
  // For folders, sum tokens of all selected files inside
  const allFiles = fileStore.getAllFilesInNode(node)
  let totalTokens = 0
  for (const filePath of allFiles) {
    if (fileStore.selectedPaths.has(filePath)) {
      const size = fileSizeMap.value.get(filePath)
      if (size) {
        totalTokens += Math.round(size / 4)
      }
    }
  }
  return totalTokens
}

// Get count of selected files inside a folder (for bubble-up indicator)
function getSelectedFileCount(node: FileNode): number {
  if (!node.isDir) return 0
  return fileStore.getSelectedFileCountInNode(node)
}

// Keyboard navigation handler
function handleKeyDown(event: KeyboardEvent) {
  const nodes = flattenedNodes.value
  if (nodes.length === 0) return

  // Initialize focus to first node if not set
  if (!fileStore.focusedPath) {
    fileStore.setFocusedPath(nodes[0].id)
    return
  }

  const currentIndex = nodes.findIndex(n => n.id === fileStore.focusedPath)
  if (currentIndex === -1) {
    fileStore.setFocusedPath(nodes[0].id)
    return
  }

  const currentNode = nodes[currentIndex].node
  let handled = false

  switch (event.key) {
    case 'ArrowDown':
      // Move focus to next visible node
      if (currentIndex < nodes.length - 1) {
        fileStore.setFocusedPath(nodes[currentIndex + 1].id)
        scrollFocusedIntoView()
        handled = true
      }
      break

    case 'ArrowUp':
      // Move focus to previous visible node
      if (currentIndex > 0) {
        fileStore.setFocusedPath(nodes[currentIndex - 1].id)
        scrollFocusedIntoView()
        handled = true
      }
      break

    case 'ArrowRight':
      // Expand folder if collapsed
      if (currentNode.isDir && !currentNode.isExpanded) {
        emit('toggle-expand', currentNode.path)
        handled = true
      }
      break

    case 'ArrowLeft':
      // Collapse folder if expanded
      if (currentNode.isDir && currentNode.isExpanded) {
        emit('toggle-expand', currentNode.path)
        handled = true
      }
      break

    case ' ':
      // Toggle checkbox selection
      event.preventDefault() // Prevent page scroll
      handleToggleSelect({ path: currentNode.path, shiftKey: false })
      handled = true
      break

    case 'Enter':
      // Toggle expand/collapse for folders, select for files
      if (currentNode.isDir) {
        emit('toggle-expand', currentNode.path)
      } else {
        handleToggleSelect({ path: currentNode.path, shiftKey: false })
      }
      handled = true
      break
  }

  if (handled) {
    event.preventDefault()
  }
}

function handleSelectRelated(path: string) {
  const count = fileStore.selectRelated(path)
  if (count > 0) {
    uiStore.addToast(
      t('files.relatedFilesSelected').replace('{count}', String(count)),
      'success',
      2000
    )
  } else {
    uiStore.addToast(
      t('files.relatedFilesSelected').replace('{count}', '0'),
      'info',
      2000
    )
  }
}

function handleAddDependency(path: string) {
  // Select the dependency file
  fileStore.selectPath(path)
  
  // Show toast notification
  uiStore.addToast(
    t('files.dependencyAdded').replace('{count}', '1'),
    'success',
    2000
  )
}

// Drag-to-select handlers
function handleCheckboxMouseDown(payload: { path: string, isSelected: boolean }) {
  const node = fileStore.findNode(payload.path)
  if (!node) return
  
  // For directories, use normal toggle behavior
  if (node.isDir) {
    return
  }
  
  isDraggingSelection.value = true
  dragInitialAction.value = payload.isSelected ? 'deselect' : 'select'
  draggedPaths.clear()
  draggedPaths.add(payload.path)
  
  // Apply initial action immediately
  if (dragInitialAction.value === 'select') {
    fileStore.selectPath(payload.path)
  } else {
    fileStore.deselectPath(payload.path)
  }
  
  // Add global mouseup listener
  document.addEventListener('mouseup', handleDragEnd, { once: true })
}

function handleRowMouseEnter(path: string) {
  if (!isDraggingSelection.value) return
  if (draggedPaths.has(path)) return
  
  draggedPaths.add(path)
  
  // Apply the initial action to this path
  const node = fileStore.findNode(path)
  if (!node || node.isDir) return
  
  if (dragInitialAction.value === 'select') {
    fileStore.selectPath(path)
  } else {
    fileStore.deselectPath(path)
  }
}

function handleDragEnd() {
  isDraggingSelection.value = false
  draggedPaths.clear()
  
  // Save selection if auto-save is enabled
  if (settingsStore.settings.fileExplorer.autoSaveSelection) {
    fileStore.saveSelectionToStorage()
  }
}

function scrollFocusedIntoView() {
  // Use nextTick to ensure DOM is updated
  import('vue').then(({ nextTick }) => {
    nextTick(() => {
      if (!fileStore.focusedPath) return
      
      const focusedIndex = flattenedNodes.value.findIndex(n => n.id === fileStore.focusedPath)
      if (focusedIndex === -1) return

      // Scroll the RecycleScroller to the focused item
      if (scrollerRef.value) {
        const scroller = scrollerRef.value as { scrollToItem?: (index: number) => void }
        if (scroller.scrollToItem) {
          scroller.scrollToItem(focusedIndex)
        }
      }
    })
  })
}
</script>

<style scoped>
.virtual-tree-wrapper {
  flex: 1 1 0;
  min-height: 0;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.virtual-tree-scroller {
  flex: 1 1 0;
  min-height: 0;
  overflow-y: auto !important;
  overflow-x: hidden;
}

.virtual-tree-scroller :deep(.vue-recycle-scroller__item-view) {
  margin: 0 !important;
  padding: 0 !important;
  overflow: visible !important;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: var(--text-muted);
}
</style>
