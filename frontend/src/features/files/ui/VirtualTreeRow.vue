<template>
  <div
    :class="[
      'tree-row group',
      props.isSelected ? 'tree-row-selected' : '',
      isDragging ? 'tree-row-dragging' : '',
      props.isFocused ? 'tree-row-focused' : '',
      !isRelevant ? 'tree-row-dimmed' : '',
      props.isDraggingSelection ? 'tree-row-drag-selecting' : '',
      item.relativePath ? 'tree-row-flat-search' : '',
      isHighlightedByDependency ? 'tree-row-dependency-highlight' : '',
      isOutgoingHighlight ? 'tree-row-dependency-outgoing' : '',
      isIncomingHighlight ? 'tree-row-dependency-incoming' : ''
    ]"
    :style="item.relativePath ? { paddingLeft: '8px' } : { paddingLeft: `${item.depth * 16 + 8}px` }"
    :tabindex="props.isFocused ? 0 : -1"
    @click="handleClick($event)"
    @contextmenu.prevent="handleContextMenu"
    :title="getSizeWarningTooltip()"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    draggable="true"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
  >
    <!-- Tree guide lines (hidden in flat search mode) -->
    <svg v-if="item.depth > 0 && !item.relativePath" class="tree-guides" :width="item.depth * 16 + 16" :height="rowHeight" style="shape-rendering: crispEdges; overflow: visible;">
      <!-- Vertical continuation lines for ancestors that have more siblings -->
      <template v-for="(hasMore, idx) in item.ancestorHasMoreSiblings" :key="'v-' + idx">
        <line v-if="hasMore" 
          :x1="8 + (idx + 1) * 16 + 0.5" y1="-4" 
          :x2="8 + (idx + 1) * 16 + 0.5" :y2="rowHeight + 4"
          :class="[
            'tree-guide-line',
            `tree-guide-${Math.min(idx + 1, 5)}`,
            { 'tree-guide-highlight': shouldHighlightGuide(idx) }
          ]" />
      </template>
      <!-- Current level: vertical line (full or half) + horizontal connector -->
      <line :x1="8 + item.depth * 16 + 0.5" y1="-4" 
            :x2="8 + item.depth * 16 + 0.5" :y2="item.isLast ? rowHeight / 2 : rowHeight + 4"
            :class="[
              'tree-guide-line',
              `tree-guide-${Math.min(item.depth, 5)}`,
              { 'tree-guide-highlight': shouldHighlightGuide(item.depth - 1) }
            ]" />
      <line :x1="8 + item.depth * 16" :y1="rowHeight / 2" 
            :x2="8 + item.depth * 16 + 10" :y2="rowHeight / 2"
            :class="[
              'tree-guide-line',
              `tree-guide-${Math.min(item.depth, 5)}`,
              { 'tree-guide-highlight': shouldHighlightGuide(item.depth - 1) }
            ]" />
    </svg>

    <!-- Expand/Collapse Icon -->
    <div
      v-if="item.node.isDir"
      :class="['tree-expand', item.node.isExpanded ? 'tree-expand-open' : '']"
      @click.stop="handleExpand"
    >
      <ChevronIcon />
    </div>
    <div v-else class="w-5"></div>

    <!-- Checkbox -->
    <div
      :class="[
        'tree-cb',
        props.checkboxState !== 'none' ? 'tree-cb-checked' : '',
        props.checkboxState === 'partial' ? 'tree-cb-partial' : '',
        isSelectionDisabled ? 'tree-cb-disabled' : '',
        props.isDraggingSelection ? 'tree-cb-dragging' : '',
        isCriticalSize && props.checkboxState === 'none' ? 'tree-cb-warning' : ''
      ]"
      @click.stop="handleToggleSelect($event)"
      @mousedown.stop="handleCheckboxMouseDown"
      :title="isSelectionDisabled ? t('files.binaryFile') : undefined"
    >
      <CheckIcon v-if="props.checkboxState === 'full'" />
      <div
        v-else-if="props.checkboxState === 'partial'"
        class="w-2 h-0.5 bg-white rounded-full"
      ></div>
      <AlertTriangleIcon 
        v-else-if="isCriticalSize && props.checkboxState === 'none'" 
        class="w-3 h-3 tree-warning-icon"
      />
    </div>

    <!-- Magic Wand Button (Select Related) -->
    <BaseButton
      v-if="!item.node.isDir"
      variant="ghost"
      size="xs"
      icon-only
      class="tree-wand"
      @click.stop="handleSelectRelated"
      :title="t('files.selectRelated')"
    >
      <template #icon>
        <WandIcon />
      </template>
    </BaseButton>

    <!-- Dependency Indicator -->
    <div
      v-if="hasDependencyLink"
      class="tree-dependency-container"
      @mouseenter="handleDependencyHover(true)"
      @mouseleave="handleDependencyHover(false)"
    >
      <button
        class="tree-dependency-indicator"
        :class="{ 'tree-dependency-indicator--hover': isDependencyHovered }"
        @click.stop="handleAddSingleDependency"
        :title="dependencyTooltip"
      >
        <LinkIcon class="w-3.5 h-3.5" :class="dependencyIconClass" />
        <span v-if="totalDependencyCount > 1" class="dependency-count">
          {{ totalDependencyCount }}
        </span>
      </button>
      
      <!-- Batch Add Button (shown on hover if multiple dependencies) -->
      <BaseButton
        v-if="totalDependencyCount > 1 && isDependencyHovered"
        variant="ghost"
        size="xs"
        class="tree-dependency-batch"
        @click.stop="handleAddAllDependencies"
        :title="t('files.addAllDependencies')"
      >
        <template #icon>
          <PlusIcon class="w-3 h-3" />
        </template>
      </BaseButton>
    </div>

    <!-- File/Folder Icon -->
    <div class="tree-icon" :class="{\n      'tree-icon-critical': isCriticalSize,\n      'tree-icon-heavy': isHeavySize\n    }\" :style=\"heatmapColor ? { color: heatmapColor } : {}\">\n      <FolderOpenIcon v-if=\"item.node.isDir && item.node.isExpanded\" />\n      <FolderIcon v-else-if=\"item.node.isDir\" />\n      <span v-else class=\"tree-file-icon\">{{ getFileIcon(item.node.name) }}</span>\n    </div>\n\n    <!-- Name with search highlighting -->\n    <span class=\"tree-name\">
      <template v-if="nameSegments.length > 1">
        <template v-for="(segment, idx) in nameSegments" :key="idx">
          <mark v-if="segment.isMatch" class="tree-highlight">{{ segment.text }}</mark>
          <span v-else>{{ segment.text }}</span>
        </template>
      </template>
      <template v-else>
        {{ item.displayName || item.node.name }}
      </template>
    </span>
    
    <!-- Dependency Indicator -->
    <DependencyIndicator
      v-if="fileStore.showDependencyIndicators && !item.node.isDir"
      :file-path="item.node.path"
      @show-dependencies="handleShowDependencies"
    />
    
    <!-- Search mode: show relative path -->
    <span v-if="item.relativePath" class="tree-search-path">
      {{ item.relativePath }}
    </span>

    <!-- Folder: file count + selected count + selected tokens weight -->
    <template v-if="item.node.isDir">
      <span 
        v-if="props.fileCount > 0" 
        class="tree-count"
        :title="t('files.fileCountTooltip').replace('{count}', String(props.fileCount))">
        {{ props.fileCount }}
      </span>
      <span 
        v-if="props.selectedFileCount > 0" 
        class="tree-count-selected"
        :title="t('files.selectedFileCountTooltip').replace('{count}', String(props.selectedFileCount))">
        +{{ props.selectedFileCount }}
      </span>
      <span 
        v-if="props.selectedTokens > 0"
        class="tree-weight"
        :class="`tree-weight--${weightLevel}`"
        :title="t('files.selectedTokensTooltip').replace('{count}', formatTokens(props.selectedTokens))"
      >
        {{ formatTokens(props.selectedTokens) }}
      </span>
    </template>

    <!-- File indicators -->
    <template v-else>
      <span 
        v-if="item.node.contentType === 'binary'" 
        class="tree-binary-badge"
        :title="t('files.binaryFile')"
      >
        BIN
      </span>
      <span 
        v-else-if="fileWeightLevel !== 'none'"
        class="tree-token-badge"
        :class="`tree-token-badge--${fileWeightLevel}`"
        :title="t('files.tokenCountTooltip').replace('{count}', formatTokens(fileTokens))"
      >
        {{ formatTokens(fileTokens) }}
      </span>
      <span v-else-if="item.node.size" class="tree-size">
        {{ formatSize(item.node.size) }}
      </span>
    </template>

    <!-- QuickLook button -->
    <BaseButton
      v-if="!item.node.isDir"
      variant="ghost"
      size="xs"
      icon-only
      class="tree-preview"
      @click.stop="handleQuickLook"
    >
      <template #icon>
        <EyeIcon />
      </template>
    </BaseButton>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { FlattenedNode } from '@/composables/useVirtualTree'
import { TOKEN_THRESHOLDS } from '@/config/constants'
import { getFileIcon } from '@/utils/fileIcons'
import { highlightMatches } from '@/utils/searchHighlight'
import { useHoveredFile } from '../composables/useHoveredFile'
import { useFileStore, type FileNode } from '../model/file.store'
import { useSettingsStore } from '@/stores/settings.store'
import { CheckIcon, ChevronIcon, EyeIcon, FolderIcon, FolderOpenIcon, WandIcon } from '@/components/icons'
import { BaseButton } from '@/components/ui'
import { computed, ref } from 'vue'
import type { FuseResultMatch } from 'fuse.js'
import { Link as LinkIcon, AlertTriangle as AlertTriangleIcon, Plus as PlusIcon } from 'lucide-vue-next'
import DependencyIndicator from './DependencyIndicator.vue'

const { t } = useI18n()
const hoveredFile = useHoveredFile()
const settingsStore = useSettingsStore()

const DRAG_DATA_TYPE = 'application/x-file-paths'

// Row height for SVG calculations (must match CSS min-height)
const rowHeight = computed(() => 26 * settingsStore.settings.uiScale)

interface Props {
  item: FlattenedNode
  compactMode?: boolean
  isSelected?: boolean
  isFocused?: boolean
  checkboxState?: 'none' | 'partial' | 'full'
  fileCount?: number
  selectedFileCount?: number
  selectedTokens?: number
  allowSelectBinary?: boolean
  isDraggingSelection?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compactMode: false,
  isSelected: false,
  isFocused: false,
  checkboxState: 'none',
  fileCount: 0,
  selectedFileCount: 0,
  selectedTokens: 0,
  allowSelectBinary: false,
  isDraggingSelection: false
})

const isSelectionDisabled = computed(() => {
  if (props.item.node.isDir) return false
  if (props.item.node.contentType !== 'binary') return false
  return !props.allowSelectBinary
})

type WeightLevel = 'none' | 'medium' | 'heavy' | 'critical'

const weightLevel = computed((): WeightLevel => {
  const tokens = props.selectedTokens
  if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return 'critical'
  if (tokens >= TOKEN_THRESHOLDS.HEAVY) return 'heavy'
  if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return 'medium'
  return 'none'
})

const fileTokens = computed(() => {
  if (props.item.node.isDir || !props.item.node.size) return 0
  return Math.round(props.item.node.size / TOKEN_THRESHOLDS.BYTES_PER_TOKEN)
})

const heatmapColor = computed(() => {
  if (!settingsStore.settings.fileExplorer.showTokenBars) return null
  
  if (props.item.node.isDir) {
    if (props.selectedTokens >= TOKEN_THRESHOLDS.CRITICAL) return 'rgba(248, 113, 113, 0.4)'
    if (props.selectedTokens >= TOKEN_THRESHOLDS.HEAVY) return 'rgba(251, 146, 60, 0.4)'
    if (props.selectedTokens >= TOKEN_THRESHOLDS.MEDIUM) return 'rgba(252, 211, 77, 0.4)'
    return null
  }
  
  if (fileTokens.value >= TOKEN_THRESHOLDS.CRITICAL) return 'rgba(248, 113, 113, 0.4)'
  if (fileTokens.value >= TOKEN_THRESHOLDS.HEAVY) return 'rgba(251, 146, 60, 0.4)'
  if (fileTokens.value >= TOKEN_THRESHOLDS.MEDIUM) return 'rgba(252, 211, 77, 0.4)'
  return null
})

const fileWeightLevel = computed((): WeightLevel => {
  const tokens = fileTokens.value
  if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return 'critical'
  if (tokens >= TOKEN_THRESHOLDS.HEAVY) return 'heavy'
  if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return 'medium'
  return 'none'
})

// Size guard: check if file exceeds size thresholds
const isCriticalSize = computed(() => {
  if (props.item.node.isDir) return false
  const tokens = fileTokens.value
  return tokens >= TOKEN_THRESHOLDS.CRITICAL
})

const isHeavySize = computed(() => {
  if (props.item.node.isDir) return false
  const tokens = fileTokens.value
  return tokens >= TOKEN_THRESHOLDS.HEAVY && tokens < TOKEN_THRESHOLDS.CRITICAL
})

const emit = defineEmits<{
  (e: 'toggle-select', payload: { path: string, shiftKey: boolean }): void
  (e: 'toggle-expand', path: string): void
  (e: 'contextmenu', node: FileNode, event: MouseEvent): void
  (e: 'quicklook', path: string): void
  (e: 'select-related', path: string): void
  (e: 'add-dependency', path: string): void
  (e: 'add-all-dependencies', path: string): void
  (e: 'checkbox-mousedown', payload: { path: string, isSelected: boolean }): void
  (e: 'row-mouseenter', path: string): void
  (e: 'dependency-hover', payload: { path: string, isHovering: boolean }): void
}>()

const fileStore = useFileStore()
const isDragging = ref(false)
const isDependencyHovered = ref(false)

// Dependency indicator logic
const hasDependencyLink = computed(() => {
  if (props.item.node.isDir) return false
  if (!fileStore.showDependencyIndicators) return false
  return fileStore.allFileDependencies.has(props.item.node.path)
})

const dependencyInfo = computed(() => {
  if (!hasDependencyLink.value) return null
  return fileStore.allFileDependencies.get(props.item.node.path)
})

const incomingCount = computed(() => dependencyInfo.value?.incoming.length || 0)
const outgoingCount = computed(() => dependencyInfo.value?.outgoing.length || 0)
const totalDependencyCount = computed(() => incomingCount.value + outgoingCount.value)

const dependencyIconClass = computed(() => {
  if (incomingCount.value > 0 && outgoingCount.value > 0) {
    return 'text-purple-400' // Bidirectional
  } else if (incomingCount.value > 0) {
    return 'text-blue-400' // Incoming
  } else if (outgoingCount.value > 0) {
    return 'text-green-400' // Outgoing
  }
  return ''
})

const dependencyTooltip = computed(() => {
  if (!hasDependencyLink.value || !dependencyInfo.value) return ''

  const parts: string[] = []
  
  if (incomingCount.value > 0) {
    const files = dependencyInfo.value.incoming.map(p => p.substring(p.lastIndexOf('/') + 1))
    if (files.length === 1) {
      parts.push(t('files.incomingDep', { file: files[0] }))
    } else {
      parts.push(t('files.incomingDependencies', { count: files.length }))
    }
  }
  
  if (outgoingCount.value > 0) {
    const files = dependencyInfo.value.outgoing.map(p => p.substring(p.lastIndexOf('/') + 1))
    if (files.length === 1) {
      parts.push(t('files.outgoingDep', { file: files[0] }))
    } else {
      parts.push(t('files.outgoingDependencies', { count: files.length }))
    }
  }

  if (totalDependencyCount.value > 1) {
    parts.push('\n' + t('files.addAllDependencies'))
  }

  return parts.join('\n')
})

const isHoveredOrFocused = computed(() => {
  return props.isSelected || hoveredFile.isHovered(props.item.node.path)
})

const isRelevant = computed(() => {
  if (!fileStore.isZenMode) return true
  return !fileStore.shouldDimNode(props.item.node)
})

// Dependency highlighting
const hoveredFilePath = ref<string | null>(null)

const isHighlightedByDependency = computed(() => {
  if (!hoveredFilePath.value || !settingsStore.settings.fileExplorer.autoHighlightDependencies) {
    return false
  }
  
  const deps = fileStore.allFileDependencies.get(hoveredFilePath.value)
  if (!deps) return false
  
  return deps.incoming.includes(props.item.node.path) || deps.outgoing.includes(props.item.node.path)
})

const isOutgoingHighlight = computed(() => {
  if (!hoveredFilePath.value || !settingsStore.settings.fileExplorer.autoHighlightDependencies) {
    return false
  }
  
  const deps = fileStore.allFileDependencies.get(hoveredFilePath.value)
  if (!deps) return false
  
  return deps.outgoing.includes(props.item.node.path)
})

const isIncomingHighlight = computed(() => {
  if (!hoveredFilePath.value || !settingsStore.settings.fileExplorer.autoHighlightDependencies) {
    return false
  }
  
  const deps = fileStore.allFileDependencies.get(hoveredFilePath.value)
  if (!deps) return false
  
  return deps.incoming.includes(props.item.node.path)
})

// Search highlighting: split name into segments
const nameSegments = computed(() => {
  // Check if this is a search result with matches
  const searchResult = props.item as unknown as { matches?: ReadonlyArray<FuseResultMatch> }
  if (searchResult.matches && fileStore.searchQuery) {
    return highlightMatches(props.item.node.name, searchResult.matches)
  }
  return [{ text: props.item.node.name, isMatch: false }]
})

function shouldHighlightGuide(depth: number): boolean {
  if (!isHoveredOrFocused.value) return false
  // Highlight parent's guide line (depth - 1)
  return depth === props.item.depth - 1
}

function handleClick(event: MouseEvent) {
  if (props.item.node.isDir) {
    emit('toggle-expand', props.item.node.path)
  } else {
    emit('toggle-select', { path: props.item.node.path, shiftKey: event.shiftKey })
  }
}

function handleExpand() {
  emit('toggle-expand', props.item.node.path)
}

function handleToggleSelect(event: MouseEvent) {
  if (isSelectionDisabled.value) return
  emit('toggle-select', { path: props.item.node.path, shiftKey: event.shiftKey })
}

function handleContextMenu(event: MouseEvent) {
  emit('contextmenu', props.item.node, event)
}

function handleQuickLook() {
  emit('quicklook', props.item.node.path)
}

function handleMouseEnter() {
  hoveredFile.setHovered(props.item.node.path, props.item.node.isDir)
  
  // Track hovered file for dependency highlighting
  if (!props.item.node.isDir && settingsStore.settings.fileExplorer.autoHighlightDependencies) {
    hoveredFilePath.value = props.item.node.path
  }
  
  // Emit row-mouseenter for drag-to-select
  if (props.isDraggingSelection) {
    emit('row-mouseenter', props.item.node.path)
  }
}

function handleMouseLeave() {
  hoveredFile.clearHovered(props.item.node.path)
  
  // Clear hovered file for dependency highlighting
  if (!props.item.node.isDir) {
    hoveredFilePath.value = null
  }
}

function handleShowDependencies(direction: 'incoming' | 'outgoing') {
  // For now, just select the dependencies
  const deps = fileStore.allFileDependencies.get(props.item.node.path)
  if (!deps) return
  
  const pathsToSelect = direction === 'incoming' ? deps.incoming : deps.outgoing
  
  // Select all dependency files
  for (const path of pathsToSelect) {
    if (!fileStore.selectedPaths.has(path)) {
      fileStore.selectPath(path)
    }
  }
}

function handleCheckboxMouseDown() {
  if (isSelectionDisabled.value) return
  
  // Only start drag-to-select for files, not directories
  if (!props.item.node.isDir) {
    // Emit checkbox-mousedown to start drag-to-select
    emit('checkbox-mousedown', {
      path: props.item.node.path,
      isSelected: props.checkboxState === 'full'
    })
  }
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return Math.round(bytes / 1024) + ' KB'
  return Math.round(bytes / (1024 * 1024)) + ' MB'
}

function formatTokens(tokens: number): string {
  if (tokens < 1000) return tokens + ''
  return Math.round(tokens / 1000) + 'k'
}

function getSizeWarningTooltip(): string {
  if (props.item.node.isDir) {
    return props.item.node.isIgnored ? `${props.item.node.path} (ignored)` : props.item.node.path
  }
  
  if (isCriticalSize.value) {
    const tokens = fileTokens.value
    const percent = Math.round((tokens / TOKEN_THRESHOLDS.MAX_CONTEXT) * 100)
    return t('files.criticalSizeWarning')
      .replace('{tokens}', formatTokens(tokens))
      .replace('{percent}', String(percent))
  }
  
  if (isHeavySize.value) {
    const tokens = fileTokens.value
    return t('files.heavySizeWarning')
      .replace('{tokens}', formatTokens(tokens))
  }
  
  return props.item.node.isIgnored ? `${props.item.node.path} (ignored)` : props.item.node.path
}

function handleDragStart(e: DragEvent) {
  if (!e.dataTransfer) return
  isDragging.value = true
  let pathsToDrag: string[]
  if (props.isSelected && fileStore.selectedCount > 1) {
    pathsToDrag = fileStore.selectedFilesList || []
  } else if (props.item.node.isDir) {
    pathsToDrag = fileStore.getAllFilesInNode(props.item.node)
  } else {
    pathsToDrag = [props.item.node.path]
  }
  e.dataTransfer.setData(DRAG_DATA_TYPE, JSON.stringify(pathsToDrag))
  e.dataTransfer.effectAllowed = 'copy'
}

function handleDragEnd() {
  isDragging.value = false
}

function handleSelectRelated() {
  emit('select-related', props.item.node.path)
}

function handleAddSingleDependency() {
  emit('add-dependency', props.item.node.path)
}

function handleAddAllDependencies() {
  emit('add-all-dependencies', props.item.node.path)
}

function handleDependencyHover(isHovering: boolean) {
  isDependencyHovered.value = isHovering
  
  if (settingsStore.settings.fileExplorer.autoHighlightDependencies) {
    emit('dependency-hover', {
      path: props.item.node.path,
      isHovering
    })
  }
}
</script>

<style scoped>
.tree-row {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    padding: 0;
    border-radius: 0;
    cursor: pointer;
    user-select: none;
    position: relative;
    color: var(--text-secondary);
    height: var(--tree-row-height);
    min-height: var(--tree-row-height);
    max-height: var(--tree-row-height);
    margin: 0;
    box-sizing: border-box;
}

.tree-row:hover {
    background: var(--tree-hover-bg);
    color: var(--text-primary);
}

.tree-row-selected {
    background: rgba(99, 102, 241, 0.18) !important;
    color: white !important;
}

.tree-row-selected:hover {
    background: rgba(99, 102, 241, 0.25) !important;
}

.tree-guides {
    position: absolute;
    top: 0;
    left: 0;
    pointer-events: none;
    overflow: visible;
    z-index: 0;
    height: var(--tree-row-height) !important;
    display: block;
}

.tree-guide-line {
    fill: none;
    stroke-width: 1px;
    stroke-linecap: square;
    stroke-linejoin: miter;
    opacity: 0.5;
    transition: opacity var(--transition-fast);
    stroke: var(--border-subtle);
}

.tree-row:hover .tree-guide-line {
    opacity: 0.8;
}

.tree-expand {
    flex-shrink: 0;
    width: var(--tree-expand-size);
    height: var(--tree-expand-size);
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: var(--radius-sm);
    transition: all var(--transition-fast);
    color: #94a3b8;
}

.tree-row:hover .tree-expand {
    color: #cbd5e1;
}

.tree-expand:hover {
    color: #f1f5f9;
    background: rgba(148, 163, 184, 0.15);
}

.tree-expand-icon {
    width: 1.125rem;
    height: 1.125rem;
    transition: transform var(--transition-fast);
}

.tree-expand-open .tree-expand-icon {
    transform: rotate(90deg);
}

.tree-cb {
    flex-shrink: 0;
    width: 16px;
    height: 16px;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1.5px solid #64748b;
    cursor: pointer;
    transition: all var(--transition-fast);
    background: transparent;
}

.tree-row:hover .tree-cb {
    border-color: #818cf8;
    background: rgba(99, 102, 241, 0.08);
}

.tree-cb:hover {
    border-color: #818cf8;
    background: rgba(99, 102, 241, 0.12);
}

.tree-cb-checked {
    background: #6366f1;
    border-color: #6366f1;
    box-shadow: 0 0 4px rgba(99, 102, 241, 0.4);
}

.tree-cb-partial {
    background: rgba(99, 102, 241, 0.35);
    border-color: #818cf8;
}

.tree-cb-icon {
    width: 0.75rem;
    height: 0.75rem;
    color: white;
}

.tree-cb-disabled {
    opacity: 0.35;
    cursor: not-allowed;
    border-color: #475569;
    background: rgba(71, 85, 105, 0.2);
}

.tree-cb-warning {
    border-color: #fb923c;
    background: rgba(251, 146, 60, 0.1);
}

.virtual-tree-scroller:not(.is-scrolling) .tree-cb-warning {
    animation: pulse-warning 2s ease-in-out infinite;
}

.tree-icon-critical {
}

.virtual-tree-scroller:not(.is-scrolling) .tree-icon-critical {
    animation: pulse-critical 2s ease-in-out infinite;
}

.tree-cb-warning:hover {
    border-color: #f97316;
    background: rgba(251, 146, 60, 0.2);
    animation: none;
}

.tree-warning-icon {
    color: #fb923c;
}

@keyframes pulse-warning {
    0%, 100% {
        box-shadow: 0 0 0 0 rgba(251, 146, 60, 0.4);
    }
    50% {
        box-shadow: 0 0 0 4px rgba(251, 146, 60, 0);
    }
}

.tree-icon {
    flex-shrink: 0;
    width: var(--tree-icon-size);
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition-fast);
}

.tree-icon-critical {
    animation: pulse-critical 2s ease-in-out infinite;
}

.tree-icon-heavy {
    filter: drop-shadow(0 0 2px rgba(251, 146, 60, 0.3));
}

@keyframes pulse-critical {
    0%, 100% {
        filter: drop-shadow(0 0 2px rgba(248, 113, 113, 0.5));
    }
    50% {
        filter: drop-shadow(0 0 4px rgba(248, 113, 113, 0.8));
    }
}

.tree-folder-icon {
    width: 18px;
    height: 18px;
    color: var(--tree-folder-color);
    transition: all var(--transition-fast);
}

.tree-name {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    font-size: var(--font-size-sm);
    font-weight: var(--font-weight-medium);
}

.tree-count {
    font-size: 10px;
    margin-left: auto;
    margin-right: var(--space-2);
    padding: 2px 6px;
    border-radius: var(--radius-md);
    color: #94a3b8;
    background: rgba(100, 116, 139, 0.25);
    font-weight: var(--font-weight-medium);
    min-width: 1.25rem;
    text-align: center;
    font-variant-numeric: tabular-nums;
}

.tree-size {
    font-size: 11px;
    margin-left: auto;
    margin-right: var(--space-2);
    color: #94a3b8;
    font-weight: var(--font-weight-medium);
    font-variant-numeric: tabular-nums;
}

.tree-weight, .tree-token-badge {
    font-size: 10px;
    margin-right: var(--space-2);
    padding: 2px 6px;
    border-radius: var(--radius-md);
    font-weight: 600;
    font-variant-numeric: tabular-nums;
    font-family: ui-monospace, monospace;
}

.tree-weight--medium, .tree-token-badge--medium {
    color: #fcd34d;
    background: rgba(252, 211, 77, 0.15);
    border: 1px solid rgba(252, 211, 77, 0.25);
}

.tree-weight--heavy, .tree-token-badge--heavy {
    color: #fb923c;
    background: rgba(251, 146, 60, 0.15);
    border: 1px solid rgba(251, 146, 60, 0.25);
}

.tree-weight--critical, .tree-token-badge--critical {
    color: #f87171;
    background: rgba(248, 113, 113, 0.15);
    border: 1px solid rgba(248, 113, 113, 0.3);
}

.tree-preview {
    opacity: 0;
    visibility: hidden;
    margin-left: auto;
    margin-right: var(--space-1);
    transition: all var(--transition-fast);
}

.tree-row:hover .tree-preview {
    opacity: 1;
    visibility: visible;
}

.tree-wand {
    opacity: 0;
    visibility: hidden;
    transition: all var(--transition-fast);
}

.tree-row:hover .tree-wand {
    opacity: 1;
    visibility: visible;
}

.tree-dependency-container {
    display: flex;
    align-items: center;
    gap: 2px;
    margin-right: var(--space-1);
    opacity: 0.6;
    transition: opacity var(--transition-fast);
}

.tree-row:hover .tree-dependency-container {
    opacity: 1;
}

.tree-row:hover .tree-preview {
    opacity: 1;
}

.tree-dependency-container {
    display: flex;
    align-items: center;
    gap: 2px;
    margin-right: var(--space-1);
}

.tree-dependency-indicator {
    display: flex;
    align-items: center;
    gap: 2px;
    padding: 2px 4px;
    border-radius: var(--radius-sm);
    background: rgba(99, 102, 241, 0.1);
    border: 1px solid rgba(99, 102, 241, 0.2);
    color: var(--text-secondary);
    transition: all var(--transition-fast);
    cursor: pointer;
}

.tree-dependency-indicator:hover {
    background: rgba(99, 102, 241, 0.2);
    border-color: rgba(99, 102, 241, 0.4);
    transform: scale(1.05);
}

.tree-dependency-indicator--hover {
    animation: pulse-dependency 1.5s ease-in-out infinite;
}

@keyframes pulse-dependency {
    0%, 100% {
        box-shadow: 0 0 0 0 rgba(99, 102, 241, 0.4);
    }
    50% {
        box-shadow: 0 0 0 4px rgba(99, 102, 241, 0);
    }
}

.dependency-count {
    font-size: 10px;
    font-weight: 600;
    color: var(--text-primary);
    min-width: 12px;
    text-align: center;
}

.tree-dependency-batch {
}

.virtual-tree-scroller:not(.is-scrolling) .tree-dependency-batch {
    animation: slideIn 0.2s ease-out;
}

@keyframes slideIn {
    from {
        opacity: 0;
        transform: translateX(-4px);
    }
    to {
        opacity: 1;
        transform: translateX(0);
    }
}

.tree-wand {
    opacity: 0;
    transition: all var(--transition-fast);
}

.tree-row:hover .tree-wand {
    opacity: 1;
}

.tree-row-dragging {
    opacity: 0.6;
    cursor: grabbing;
}

.tree-row-highlighted {
    background: rgba(99, 102, 241, 0.12) !important;
}

.virtual-tree-scroller:not(.is-scrolling) .tree-row-highlighted {
    animation: highlight-fade 0.3s ease-out;
}

@keyframes highlight-fade {
    from {
        background: rgba(99, 102, 241, 0.25);
    }
    to {
        background: rgba(99, 102, 241, 0.12);
    }
}
</style>
