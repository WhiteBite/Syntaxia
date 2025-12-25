<template>
  <div class="select-none p-2">
    <template v-for="node in visibleNodes" :key="node.path">
      <div :class="[
        'tree-row group',
        isSelected(node.path) ? 'tree-row-selected' : '',
        animatingNodes.has(node.path) ? 'tree-stagger' : ''
      ]" :style="{ paddingLeft: `${node.depth * 20 + 12}px` }" @click="handleClick(node)"
        @dblclick="handleDoubleClick(node)" tabindex="0">
        <!-- Tree guide lines -->
        <svg v-if="node.depth > 0" class="tree-guides" :width="node.depth * 20" height="32" style="left: 12px">
          <!-- Vertical lines for ancestors -->
          <template v-for="d in node.depth" :key="d">
            <line v-if="shouldDrawVerticalLine(node, d)" :x1="(d - 1) * 20 + 10" y1="0" :x2="(d - 1) * 20 + 10" y2="32"
              class="tree-guide-line" :class="`tree-guide-${Math.min(d, 5)}`" />
          </template>
          <!-- Connector path -->
          <path :d="getConnectorPath(node)" class="tree-guide-line" :class="`tree-guide-${Math.min(node.depth, 5)}`" />
        </svg>

        <!-- Expand/Collapse Icon -->
        <div v-if="node.isDir" :class="['tree-expand', isExpanded(node.path) ? 'tree-expand-open' : '']"
          @click.stop="handleClick(node)">
          <svg class="tree-expand-icon" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
              clip-rule="evenodd" />
          </svg>
        </div>
        <div v-else class="w-5"></div>

        <!-- Checkbox -->
        <div :class="[
          'tree-cb',
          isSelected(node.path) ? 'tree-cb-checked' : '',
          node.isDir && hasSomeSelected(node.path) && !hasAllSelected(node.path) ? 'tree-cb-partial' : '',
          rippleNode === node.path ? 'tree-cb-ripple' : ''
        ]" @click.stop="toggleSelect(node)">
          <svg v-if="isSelected(node.path) || (node.isDir && hasAllSelected(node.path))" class="tree-cb-icon"
            fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd"
              d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
              clip-rule="evenodd" />
          </svg>
          <div v-else-if="node.isDir && hasSomeSelected(node.path)" class="w-2 h-0.5 bg-white rounded-full"></div>
        </div>

        <!-- File/Folder Icon -->
        <div class="tree-icon">
          <svg v-if="node.isDir" class="tree-folder-icon" fill="currentColor" viewBox="0 0 20 20">
            <path d="M2 6a2 2 0 012-2h5l2 2h5a2 2 0 012 2v6a2 2 0 01-2 2H4a2 2 0 01-2-2V6z" />
          </svg>
          <span v-else class="tree-file-icon">{{ getFileIcon(node.name) }}</span>
        </div>

        <!-- Name -->
        <span class="tree-name">{{ node.name }}</span>

        <!-- Preview button -->
        <button v-if="!node.isDir" @click.stop="emit('preview-file', node.path)" class="tree-preview" title="Preview">
          <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        </button>

        <!-- Count -->
        <span v-if="node.isDir" :class="['tree-count', countAnimating === node.path ? 'tree-count-bounce' : '']">
          {{ getFilesInDir(node.path).length }}
        </span>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { getFileIcon } from '@/utils/fileIcons'
import { computed, ref, watch } from 'vue'

interface TreeNode {
  path: string
  name: string
  isDir: boolean
  depth: number
  isLast: boolean
  parentPath: string
}

const props = defineProps<{
  files: string[]
  selectedPaths: Set<string>
}>()

const emit = defineEmits<{
  (e: 'toggle-select', path: string): void
  (e: 'select-folder', paths: string[]): void
  (e: 'preview-file', path: string): void
}>()

const expandedPaths = ref<Set<string>>(new Set())
const animatingNodes = ref<Set<string>>(new Set())
const rippleNode = ref<string | null>(null)
const countAnimating = ref<string | null>(null)
const initializedFiles = ref<string | null>(null)

// Auto-expand root folders when files change
watch(() => props.files, (newFiles) => {
  if (!newFiles?.length) return

  // Create a signature to detect if files actually changed
  const signature = newFiles.slice(0, 5).join('|')
  if (signature === initializedFiles.value) return
  initializedFiles.value = signature

  // Find and expand root-level directories
  const rootDirs = new Set<string>()
  for (const file of newFiles) {
    const firstPart = file.split('/')[0]
    if (firstPart && file.includes('/')) {
      rootDirs.add(firstPart)
    }
  }

  // Add root dirs to expanded (preserve existing expanded state)
  rootDirs.forEach(dir => expandedPaths.value.add(dir))
  expandedPaths.value = new Set(expandedPaths.value)
}, { immediate: true })

function isSelected(path: string): boolean {
  return props.selectedPaths?.has(path) ?? false
}

function isExpanded(path: string): boolean {
  return expandedPaths.value.has(path)
}

function getFilesInDir(dirPath: string): string[] {
  if (!props.files || !Array.isArray(props.files)) return []
  return props.files.filter(f => f.startsWith(dirPath + '/'))
}

function hasSomeSelected(dirPath: string): boolean {
  if (!props.selectedPaths) return false
  const filesInDir = getFilesInDir(dirPath)
  return filesInDir.some(f => props.selectedPaths.has(f))
}

function hasAllSelected(dirPath: string): boolean {
  if (!props.selectedPaths) return false
  const filesInDir = getFilesInDir(dirPath)
  return filesInDir.length > 0 && filesInDir.every(f => props.selectedPaths.has(f))
}

// Build tree with isLast info for proper lines
const visibleNodes = computed(() => {
  if (!props.files?.length) return []

  const nodes: TreeNode[] = []
  const addedDirs = new Set<string>()
  const sortedFiles = [...props.files].sort()

  // Track children per directory
  const dirChildren = new Map<string, string[]>()

  for (const file of sortedFiles) {
    const parts = file.split('/')
    let currentPath = ''

    for (let i = 0; i < parts.length; i++) {
      const parentPath = currentPath
      currentPath = currentPath ? `${currentPath}/${parts[i]}` : parts[i]

      if (!dirChildren.has(parentPath)) dirChildren.set(parentPath, [])
      const children = dirChildren.get(parentPath)!
      if (!children.includes(currentPath)) children.push(currentPath)
    }
  }

  for (const file of sortedFiles) {
    const parts = file.split('/')
    let currentPath = ''
    let isVisible = true

    for (let i = 0; i < parts.length - 1; i++) {
      const parentPath = currentPath
      currentPath = currentPath ? `${currentPath}/${parts[i]}` : parts[i]

      if (parentPath && !expandedPaths.value.has(parentPath)) {
        isVisible = false
        break
      }

      if (!addedDirs.has(currentPath)) {
        addedDirs.add(currentPath)
        const siblings = dirChildren.get(parentPath) || []
        nodes.push({
          path: currentPath,
          name: parts[i],
          isDir: true,
          depth: i,
          isLast: siblings.indexOf(currentPath) === siblings.length - 1,
          parentPath
        })
      }
    }

    if (isVisible) {
      const parentPath = parts.slice(0, -1).join('/')
      if (!parentPath || expandedPaths.value.has(parentPath)) {
        const siblings = dirChildren.get(parentPath) || []
        nodes.push({
          path: file,
          name: parts[parts.length - 1],
          isDir: false,
          depth: parts.length - 1,
          isLast: siblings.indexOf(file) === siblings.length - 1,
          parentPath
        })
      }
    }
  }

  return nodes
})

// Check if should draw vertical line at depth d for this node
function shouldDrawVerticalLine(node: TreeNode, d: number): boolean {
  if (d === node.depth) return !node.isLast

  // Check ancestor at depth d - find the ancestor node
  const parts = node.path.split('/')
  if (d > parts.length) return false

  const ancestorPath = parts.slice(0, d).join('/')
  const ancestorNode = visibleNodes.value.find(n => n.path === ancestorPath)

  if (!ancestorNode) return false
  return !ancestorNode.isLast
}

// Get SVG path for connector
function getConnectorPath(node: TreeNode): string {
  const x = (node.depth - 1) * 20 + 10
  const midY = 16
  if (node.isLast) {
    return `M ${x} 0 L ${x} ${midY} L ${x + 10} ${midY}`
  }
  return `M ${x} 0 L ${x} 32 M ${x} ${midY} L ${x + 10} ${midY}`
}

function handleClick(node: TreeNode) {
  if (node.isDir) {
    if (expandedPaths.value.has(node.path)) {
      expandedPaths.value.delete(node.path)
    } else {
      expandedPaths.value.add(node.path)
      // Trigger stagger animation
      setTimeout(() => {
        const children = visibleNodes.value.filter(n => n.parentPath === node.path)
        children.forEach(child => animatingNodes.value.add(child.path))
        setTimeout(() => children.forEach(child => animatingNodes.value.delete(child.path)), 300)
      }, 10)
    }
    expandedPaths.value = new Set(expandedPaths.value)
  }
}

function handleDoubleClick(node: TreeNode) {
  if (!node.isDir) {
    emit('preview-file', node.path)
  }
}

function toggleSelect(node: TreeNode) {
  rippleNode.value = node.path
  setTimeout(() => { rippleNode.value = null }, 400)

  if (node.isDir) {
    emit('select-folder', getFilesInDir(node.path))
    countAnimating.value = node.path
    setTimeout(() => { countAnimating.value = null }, 300)
  } else {
    emit('toggle-select', node.path)
  }
}
</script>
