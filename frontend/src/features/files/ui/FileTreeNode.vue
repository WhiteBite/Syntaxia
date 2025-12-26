<template>
  <div class="select-none">
    <div :class="[
      'tree-row group',
      isSelected ? 'tree-row-selected' : '',
      hasSelectedChildren && !isSelected ? 'tree-row-has-selected' : '',
      { 'opacity-50': node.isIgnored },
      animating ? 'tree-stagger' : ''
    ]" :style="{ paddingLeft: `${depth * 16 + 8}px` }" @click="handleClick" @contextmenu.prevent="handleContextMenu"
      :title="node.isIgnored ? `${node.path} (ignored)` : node.path" tabindex="0" @keydown="handleKeydown"
      @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave"
      :data-path="node.path" :data-is-dir="node.isDir">
      <!-- Tree guide lines - continuous vertical + L/T connectors -->
      <svg v-if="depth > 0" class="tree-guides" :width="depth * 16 + 16" height="32" :style="{ left: '0' }">
        <!-- Vertical continuation lines for ancestors that have more siblings -->
        <template v-for="(hasMore, idx) in ancestorHasMoreSiblings" :key="'v-' + idx">
          <line v-if="hasMore" 
            :x1="idx * 16 + 16" y1="0" 
            :x2="idx * 16 + 16" y2="32"
            class="tree-guide-line" :class="`tree-guide-${Math.min(idx + 1, 5)}`" />
        </template>
        <!-- Current level: vertical part (L or └ shape) -->
        <line :x1="(depth - 1) * 16 + 16" y1="0" 
              :x2="(depth - 1) * 16 + 16" :y2="isLast ? 16 : 32"
              class="tree-guide-line" :class="`tree-guide-${Math.min(depth, 5)}`" />
        <!-- Current level: horizontal connector -->
        <line :x1="(depth - 1) * 16 + 16" y1="16" 
              :x2="depth * 16 + 12" y2="16"
              class="tree-guide-line" :class="`tree-guide-${Math.min(depth, 5)}`" />
      </svg>

      <!-- Expand/Collapse Icon -->
      <button v-if="node.isDir" 
        :class="['tree-expand', compactInfo.lastNode.isExpanded ? 'tree-expand-open' : '']"
        @click.stop="handleClick"
        :aria-label="compactInfo.lastNode.isExpanded ? t('files.collapseAll') : t('files.expandAll')"
        :aria-expanded="compactInfo.lastNode.isExpanded">
        <ChevronIcon />
      </button>
      <div v-else class="w-5" aria-hidden="true"></div>

      <!-- Checkbox -->
      <button :class="[
        'tree-cb',
        (isSelected && !node.isDir) || (node.isDir && selectionState !== 'none') ? 'tree-cb-checked' : '',
        node.isDir && selectionState === 'partial' ? 'tree-cb-partial' : '',
        ripple ? 'tree-cb-ripple' : ''
      ]" @click.stop="handleToggleSelect"
        :aria-label="isSelected ? t('contextMenu.deselectAll') : t('contextMenu.selectAll')"
        :aria-checked="isSelected"
        :title="partialSelectionTooltip || undefined"
        role="checkbox">
        <CheckIcon v-if="isSelected && !node.isDir" />
        <CheckIcon v-else-if="node.isDir && selectionState === 'full'" />
        <div v-else-if="node.isDir && selectionState === 'partial'" class="w-2 h-0.5 bg-white rounded-full" aria-hidden="true"></div>
      </button>

      <!-- File/Folder Icon -->
      <div class="tree-icon">
        <FolderOpenIcon v-if="node.isDir && compactInfo.lastNode.isExpanded" />
        <FolderIcon v-else-if="node.isDir" />
        <span v-else class="tree-file-icon">{{ getFileIcon(node.name) }}</span>
      </div>

      <!-- Ignored Icon -->
      <div v-if="node.isIgnored" class="relative group/ignored">
        <svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
        </svg>
        <div class="tree-tooltip hidden group-hover/ignored:block -top-8 left-0">
          Ignored by rules
        </div>
      </div>

      <!-- Name -->
      <span class="tree-name">{{ displayName }}</span>

      <!-- File count badge -->
      <span v-if="node.isDir && node.children" :class="['tree-count', countBounce ? 'tree-count-bounce' : '']"
        :title="`${fileCount} files`">
        {{ fileCount }}
      </span>

      <!-- File size -->
      <span v-else-if="node.size" class="tree-size">
        {{ formatSize(node.size) }}
      </span>

      <!-- QuickLook button (files only) -->
      <button v-if="!node.isDir" class="tree-preview" @click.stop="handleQuickLook" :title="t('files.quickLook')">
        <EyeIcon />
      </button>
    </div>

    <!-- Children (recursive) with animation -->
    <Transition name="tree-children">
      <div v-if="node.isDir && compactInfo.lastNode.isExpanded && displayChildren.length > 0" class="tree-children">
        <FileTreeNode v-for="(child, index) in displayChildren" :key="child.path" :node="child" :depth="depth + 1"
          :compact-mode="compactMode" :is-last="index === displayChildren.length - 1"
          :ancestor-has-more-siblings="childAncestorHasMoreSiblings" @toggle-select="$emit('toggle-select', $event)"
          @toggle-expand="$emit('toggle-expand', $event)"
          @contextmenu="(childNode, event) => $emit('contextmenu', childNode, event)"
          @quicklook="$emit('quicklook', $event)" />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { getFileIcon } from '@/utils/fileIcons'
import { computed, ref, watch } from 'vue'
import { useHoveredFile } from '../composables/useHoveredFile'
import { useFileStore, type FileNode } from '../model/file.store'

const { t } = useI18n()
const hoveredFile = useHoveredFile()

interface Props {
  node: FileNode
  depth?: number
  compactMode?: boolean
  isLast?: boolean
  ancestorHasMoreSiblings?: boolean[] // Array of booleans for each ancestor level
}

const props = withDefaults(defineProps<Props>(), {
  depth: 0,
  compactMode: false,
  isLast: false,
  ancestorHasMoreSiblings: () => []
})

const emit = defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'toggle-expand', path: string): void
  (e: 'contextmenu', node: FileNode, event: MouseEvent): void
  (e: 'quicklook', path: string): void
}>()

const fileStore = useFileStore()
const ripple = ref(false)
const countBounce = ref(false)
const animating = ref(false)

const isSelected = computed(() => fileStore.selectedPaths.has(props.node.path))

// Compact mode: merge single-child folder chains (e.g., src/main/java → "src/main/java")
const compactInfo = computed(() => {
  if (!props.compactMode || !props.node.isDir || !props.node.children) {
    return { name: props.node.name, skipChildren: null as FileNode[] | null, lastNode: props.node, chainPaths: [props.node.path] }
  }
  
  // Check if we should compact: single folder child only
  const children = props.node.children
  if (children.length === 1 && children[0].isDir) {
    const namePath: string[] = [props.node.name]
    const chainPaths: string[] = [props.node.path]
    let current = props.node
    
    // Traverse down while there's only one folder child
    while (current.children && current.children.length === 1 && current.children[0].isDir) {
      const child = current.children[0]
      namePath.push(child.name)
      chainPaths.push(child.path)
      current = child
    }
    
    // Only compact if we found multiple levels
    if (namePath.length > 1) {
      return { 
        name: namePath.join('/'), 
        skipChildren: current.children || null,
        lastNode: current,
        chainPaths
      }
    }
  }
  
  return { name: props.node.name, skipChildren: null, lastNode: props.node, chainPaths: [props.node.path] }
})

const displayName = computed(() => compactInfo.value.name)

// Children to display (may skip intermediate folders in compact mode)
const displayChildren = computed(() => {
  if (compactInfo.value.skipChildren) {
    return compactInfo.value.skipChildren
  }
  return props.node.children || []
})

// Unified selection info for folder - computed once, used by multiple dependents
const folderSelectionInfo = computed(() => {
  if (!props.node.isDir) {
    return { selected: 0, total: 0, state: isSelected.value ? 'full' as const : 'none' as const }
  }

  // Early exit if no files selected globally
  if (fileStore.selectedCount === 0) {
    const total = fileStore.getAllFilesInNode(props.node).length
    return { selected: 0, total, state: 'none' as const }
  }

  const allFiles = fileStore.getAllFilesInNode(props.node)
  const total = allFiles.length
  
  if (total === 0) {
    return { selected: 0, total: 0, state: isSelected.value ? 'full' as const : 'none' as const }
  }

  // Count selected
  let selected = 0
  const selectedSet = fileStore.selectedPaths
  for (const filePath of allFiles) {
    if (selectedSet.has(filePath)) selected++
  }

  const state = selected === 0 ? 'none' as const 
    : selected === total ? 'full' as const 
    : 'partial' as const

  return { selected, total, state }
})

const selectionState = computed(() => folderSelectionInfo.value.state)

// Check if this folder contains any selected files (for path highlighting)
const hasSelectedChildren = computed(() => {
  if (!props.node.isDir || fileStore.selectedCount === 0) return false
  return selectionState.value !== 'none'
})

const fileCount = computed(() => folderSelectionInfo.value.total)

const partialSelectionTooltip = computed(() => {
  if (selectionState.value !== 'partial') return ''
  const { selected, total } = folderSelectionInfo.value
  return t('tree.partialSelection', { selected, total })
})

// Build ancestor array for children
// In compact mode, we need to account for skipped levels
const childAncestorHasMoreSiblings = computed(() => {
  // Current node's contribution: does it have more siblings?
  const currentHasMore = !props.isLast
  
  // In compact mode, if we're showing a merged path like "java/com/example",
  // we only add ONE entry (for the visual row), not multiple
  return [...props.ancestorHasMoreSiblings, currentHasMore]
})

function handleClick() {
  if (props.node.isDir) {
    // In compact mode, toggle ALL nodes in the chain together
    const { chainPaths, lastNode } = compactInfo.value
    const isCurrentlyExpanded = lastNode.isExpanded
    
    if (chainPaths.length > 1) {
      // Compact mode: expand/collapse all nodes in chain
      // We emit a special event with all paths to handle atomically
      emit('toggle-expand', JSON.stringify({ 
        paths: chainPaths, 
        expand: !isCurrentlyExpanded 
      }))
    } else {
      // Single node - simple toggle
      emit('toggle-expand', props.node.path)
    }
  }
}

function handleToggleSelect() {
  ripple.value = true
  setTimeout(() => { ripple.value = false }, 400)

  if (props.node.isDir) {
    countBounce.value = true
    setTimeout(() => { countBounce.value = false }, 300)
  }

  emit('toggle-select', props.node.path)
}

function handleContextMenu(event: MouseEvent) {
  emit('contextmenu', props.node, event)
}

function handleQuickLook() {
  emit('quicklook', props.node.path)
}

function handleMouseEnter() {
  // Store hovered node path using provide/inject pattern
  hoveredFile.setHovered(props.node.path, props.node.isDir)
}

function handleMouseLeave() {
  // Clear only if this node was the hovered one
  hoveredFile.clearHovered(props.node.path)
}

function handleKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter') {
    event.preventDefault()
    handleToggleSelect()
  } else if (event.key === 'ArrowRight' && props.node.isDir && !compactInfo.value.lastNode.isExpanded) {
    handleClick()
  } else if (event.key === 'ArrowLeft' && props.node.isDir && compactInfo.value.lastNode.isExpanded) {
    handleClick()
  }
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return Math.round(bytes / 1024) + ' KB'
  return Math.round(bytes / (1024 * 1024)) + ' MB'
}

// Trigger animation when node becomes visible
watch(() => props.node.isExpanded, (expanded) => {
  if (expanded) {
    animating.value = true
    setTimeout(() => { animating.value = false }, 300)
  }
}, { immediate: false })
</script>
