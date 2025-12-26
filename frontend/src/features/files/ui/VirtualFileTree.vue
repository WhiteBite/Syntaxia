<template>
  <div class="virtual-tree-wrapper">
    <RecycleScroller
      v-if="flattenedNodes.length > 0"
      class="virtual-tree-scroller"
      :items="flattenedNodes"
      :item-size="26"
      key-field="id"
      v-slot="{ item }"
    >
      <VirtualTreeRow
        :item="item"
        :compact-mode="compactMode"
        :is-selected="isNodeSelected(item.node)"
        :checkbox-state="getCheckboxState(item.node)"
        :file-count="getFileCount(item.node)"
        :selected-tokens="getSelectedTokens(item.node)"
        :allow-select-binary="allowSelectBinary"
        @toggle-select="$emit('toggle-select', $event)"
        @toggle-expand="$emit('toggle-expand', $event)"
        @contextmenu="(node, event) => $emit('contextmenu', node, event)"
        @quicklook="$emit('quicklook', $event)"
      />
    </RecycleScroller>
    <div v-else class="empty-state">
      <p class="empty-state-text">{{ t('files.noFiles') }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { useVirtualTree } from '@/composables/useVirtualTree'
import { computed, toRef } from 'vue'
import { RecycleScroller } from 'vue-virtual-scroller'
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css'
import { useFileStore, type FileNode } from '../model/file.store'
import VirtualTreeRow from './VirtualTreeRow.vue'

const { t } = useI18n()
const fileStore = useFileStore()

interface Props {
  nodes: FileNode[]
  compactMode?: boolean
  allowSelectBinary?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compactMode: false,
  allowSelectBinary: false,
})

defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'toggle-expand', path: string): void
  (e: 'contextmenu', node: FileNode, event: MouseEvent): void
  (e: 'quicklook', path: string): void
}>()

const nodesRef = toRef(props, 'nodes')
const { flattenedVisibleNodes } = useVirtualTree({ nodes: nodesRef })

const flattenedNodes = computed(() => flattenedVisibleNodes.value)

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
