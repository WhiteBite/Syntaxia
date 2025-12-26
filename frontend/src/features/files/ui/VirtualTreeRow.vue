<template>
  <div
    :class="[
      'tree-row group',
      props.isSelected ? 'tree-row-selected' : '',
      { 'opacity-50': item.node.isIgnored }
    ]"
    :style="{ paddingLeft: `${item.depth * 16 + 8}px` }"
    @click="handleClick"
    @contextmenu.prevent="handleContextMenu"
    :title="item.node.isIgnored ? `${item.node.path} (ignored)` : item.node.path"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <!-- Tree guide lines -->
    <svg v-if="item.depth > 0" class="tree-guides" :width="item.depth * 16 + 16" :height="rowHeight" style="shape-rendering: crispEdges; overflow: visible;">
      <!-- Vertical continuation lines for ancestors that have more siblings -->
      <template v-for="(hasMore, idx) in item.ancestorHasMoreSiblings" :key="'v-' + idx">
        <line v-if="hasMore" 
          :x1="8 + (idx + 1) * 16 + 0.5" y1="-4" 
          :x2="8 + (idx + 1) * 16 + 0.5" :y2="rowHeight + 4"
          class="tree-guide-line" :class="`tree-guide-${Math.min(idx + 1, 5)}`" />
      </template>
      <!-- Current level: vertical line (full or half) + horizontal connector -->
      <line :x1="8 + item.depth * 16 + 0.5" y1="-4" 
            :x2="8 + item.depth * 16 + 0.5" :y2="item.isLast ? 13 : rowHeight + 4"
            class="tree-guide-line" :class="`tree-guide-${Math.min(item.depth, 5)}`" />
      <line :x1="8 + item.depth * 16" y1="13" 
            :x2="8 + item.depth * 16 + 10" y2="13"
            class="tree-guide-line" :class="`tree-guide-${Math.min(item.depth, 5)}`" />
    </svg>

    <!-- Expand/Collapse Icon -->
    <div
      v-if="item.node.isDir"
      :class="['tree-expand', item.node.isExpanded ? 'tree-expand-open' : '']"
      @click.stop="handleExpand"
    >
      <svg class="tree-expand-icon" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
          clip-rule="evenodd"
        />
      </svg>
    </div>
    <div v-else class="w-5"></div>

    <!-- Checkbox -->
    <div
      :class="[
        'tree-cb',
        props.checkboxState !== 'none' ? 'tree-cb-checked' : '',
        props.checkboxState === 'partial' ? 'tree-cb-partial' : '',
        isSelectionDisabled ? 'tree-cb-disabled' : ''
      ]"
      @click.stop="handleToggleSelect"
      :title="isSelectionDisabled ? t('files.binaryFile') : undefined"
    >
      <svg
        v-if="props.checkboxState === 'full'"
        class="tree-cb-icon"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path
          fill-rule="evenodd"
          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
          clip-rule="evenodd"
        />
      </svg>
      <div
        v-else-if="props.checkboxState === 'partial'"
        class="w-2 h-0.5 bg-white rounded-full"
      ></div>
    </div>

    <!-- File/Folder Icon -->
    <div class="tree-icon">
      <svg
        v-if="item.node.isDir"
        class="tree-folder-icon"
        fill="currentColor"
        viewBox="0 0 20 20"
      >
        <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
      </svg>
      <span v-else class="tree-file-icon">{{ getFileIcon(item.node.name) }}</span>
    </div>

    <!-- Name -->
    <span class="tree-name">{{ item.node.name }}</span>

    <!-- Folder: file count + selected tokens weight -->
    <template v-if="item.node.isDir">
      <span 
        v-if="props.fileCount > 0" 
        class="tree-count"
        :title="t('files.fileCountTooltip').replace('{count}', String(props.fileCount))">
        {{ props.fileCount }}
      </span>
      <!-- Bubble-up: show selected tokens inside folder -->
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
      <!-- Binary indicator -->
      <span 
        v-if="item.node.contentType === 'binary'" 
        class="tree-binary-badge"
        :title="t('files.binaryFile')"
      >
        BIN
      </span>
      <!-- Heavy file badge (only for heavy files) -->
      <span 
        v-else-if="fileWeightLevel !== 'none'"
        class="tree-token-badge"
        :class="`tree-token-badge--${fileWeightLevel}`"
        :title="t('files.tokenCountTooltip').replace('{count}', formatTokens(fileTokens))"
      >
        {{ formatTokens(fileTokens) }}
      </span>
      <!-- Normal file size (for light files) -->
      <span v-else-if="item.node.size" class="tree-size">
        {{ formatSize(item.node.size) }}
      </span>
    </template>

    <!-- QuickLook button -->
    <button
      v-if="!item.node.isDir"
      class="tree-preview"
      @click.stop="handleQuickLook"
    >
      <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
        />
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
        />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { FlattenedNode } from '@/composables/useVirtualTree'
import { TOKEN_THRESHOLDS } from '@/config/constants'
import { getFileIcon } from '@/utils/fileIcons'
import { useHoveredFile } from '../composables/useHoveredFile'
import type { FileNode } from '../model/file.store'
import { computed } from 'vue'

const { t } = useI18n()

// Row height for SVG calculations (must match CSS min-height)
const rowHeight = 26

interface Props {
  item: FlattenedNode
  compactMode?: boolean
  isSelected?: boolean
  checkboxState?: 'none' | 'partial' | 'full'
  fileCount?: number
  selectedTokens?: number
  allowSelectBinary?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  compactMode: false,
  isSelected: false,
  checkboxState: 'none',
  fileCount: 0,
  selectedTokens: 0,
  allowSelectBinary: false
})

// Check if selection is disabled for this item
const isSelectionDisabled = computed(() => {
  if (props.item.node.isDir) return false
  if (props.item.node.contentType !== 'binary') return false
  return !props.allowSelectBinary
})

// Weight level for visual indicators
type WeightLevel = 'none' | 'medium' | 'heavy' | 'critical'

const weightLevel = computed((): WeightLevel => {
  const tokens = props.selectedTokens
  if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return 'critical'
  if (tokens >= TOKEN_THRESHOLDS.HEAVY) return 'heavy'
  if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return 'medium'
  return 'none'
})

// File's own token count (for files only)
const fileTokens = computed(() => {
  if (props.item.node.isDir || !props.item.node.size) return 0
  return Math.round(props.item.node.size / TOKEN_THRESHOLDS.BYTES_PER_TOKEN)
})

const fileWeightLevel = computed((): WeightLevel => {
  const tokens = fileTokens.value
  if (tokens >= TOKEN_THRESHOLDS.CRITICAL) return 'critical'
  if (tokens >= TOKEN_THRESHOLDS.HEAVY) return 'heavy'
  if (tokens >= TOKEN_THRESHOLDS.MEDIUM) return 'medium'
  return 'none'
})

const emit = defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'toggle-expand', path: string): void
  (e: 'contextmenu', node: FileNode, event: MouseEvent): void
  (e: 'quicklook', path: string): void
}>()

// Use singleton pattern for hovered file state (avoids provide/inject issues with virtual scrolling)
const { setHovered, clearHovered } = useHoveredFile()

function handleClick() {
  if (props.item.node.isDir) {
    emit('toggle-expand', props.item.node.path)
  } else {
    // Clicking on a file row toggles its selection
    emit('toggle-select', props.item.node.path)
  }
}

function handleExpand() {
  emit('toggle-expand', props.item.node.path)
}

function handleToggleSelect() {
  if (isSelectionDisabled.value) return
  emit('toggle-select', props.item.node.path)
}

function handleContextMenu(event: MouseEvent) {
  emit('contextmenu', props.item.node, event)
}

function handleQuickLook() {
  emit('quicklook', props.item.node.path)
}

function handleMouseEnter() {
  setHovered(props.item.node.path, props.item.node.isDir)
}

function handleMouseLeave() {
  clearHovered(props.item.node.path)
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
</script>
