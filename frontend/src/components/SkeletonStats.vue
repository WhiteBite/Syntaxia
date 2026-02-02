<template>
  <div class="skeleton-stats">
    <div
      v-for="i in cards"
      :key="i"
      class="skeleton-stats__card"
    >
      <BaseSkeleton variant="text" width="64px" height="24px" />
      <BaseSkeleton variant="text" width="48px" height="12px" />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * SkeletonStats - Specialized skeleton for statistics cards
 * Uses BaseSkeleton internally to create stat card loading states
 * 
 * @example
 * <SkeletonStats :cards="3" />
 */
import { BaseSkeleton } from '@/components/ui'

interface Props {
  cards?: number
}

withDefaults(defineProps<Props>(), {
  cards: 3
})
</script>

<style scoped>
.skeleton-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: var(--space-3);
}

.skeleton-stats__card {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  padding: var(--space-3);
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
}
</style>

