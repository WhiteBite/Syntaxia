<template>
  <div class="skeleton-file-tree">
    <div
      v-for="i in rows"
      :key="i"
      class="skeleton-file-tree__row"
      :style="{ paddingLeft: getIndent(i) + 'px' }"
    >
      <BaseSkeleton variant="circle" width="16px" height="16px" />
      <BaseSkeleton variant="text" :width="getWidth(i) + 'px'" height="16px" />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * SkeletonFileTree - Specialized skeleton for file tree structure
 * Uses BaseSkeleton internally to create a tree-like loading state
 * 
 * @example
 * <SkeletonFileTree :rows="10" />
 */
import { BaseSkeleton } from '@/components/ui'

interface Props {
  rows?: number
}

const _props = withDefaults(defineProps<Props>(), {
  rows: 10
})

function getIndent(index: number): number {
  // Create a tree-like structure with varying indents
  const pattern = [0, 0, 16, 16, 32, 32, 16, 16, 0, 0]
  return pattern[(index - 1) % pattern.length] || 0
}

function getWidth(index: number): number {
  // Varying widths to simulate different file name lengths
  const widths = [120, 80, 100, 140, 90, 110, 70, 130, 95, 105]
  return widths[(index - 1) % widths.length] || 100
}
</script>

<style scoped>
.skeleton-file-tree {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  padding: var(--space-2);
}

.skeleton-file-tree__row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  transition: padding-left var(--transition-fast);
}
</style>

