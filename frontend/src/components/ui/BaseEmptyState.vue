<template>
  <div 
    class="base-empty-state" 
    :class="[
      `base-empty-state--${size}`,
      { 'base-empty-state--with-action': $slots.action }
    ]"
  >
    <!-- Icon -->
    <div class="base-empty-state__icon">
      <slot name="icon">
        <BaseIcon v-if="icon" :icon="icon" :size="iconSize" />
      </slot>
    </div>

    <!-- Title -->
    <h3 v-if="title" class="base-empty-state__title">
      {{ title }}
    </h3>

    <!-- Description -->
    <p v-if="description" class="base-empty-state__description">
      {{ description }}
    </p>

    <!-- Action Slot -->
    <div v-if="$slots.action" class="base-empty-state__action">
      <slot name="action" />
    </div>
  </div>
</template>

<script setup lang="ts">
/**
 * BaseEmptyState - Empty state component for displaying when no content is available
 * 
 * @example
 * // Basic usage with icon component
 * <BaseEmptyState
 *   :icon="FolderIcon"
 *   title="No files selected"
 *   description="Select files from the tree to build context"
 * />
 * 
 * @example
 * // With custom icon slot
 * <BaseEmptyState title="No results" description="Try adjusting your filters">
 *   <template #icon>
 *     <SearchIcon class="w-12 h-12" />
 *   </template>
 * </BaseEmptyState>
 * 
 * @example
 * // With action button
 * <BaseEmptyState
 *   :icon="RefreshIcon"
 *   title="No data"
 *   description="Click refresh to load data"
 *   size="lg"
 * >
 *   <template #action>
 *     <BaseButton variant="primary" @click="refresh">Refresh</BaseButton>
 *   </template>
 * </BaseEmptyState>
 * 
 * @example
 * // Small size for compact spaces
 * <BaseEmptyState
 *   :icon="InboxIcon"
 *   title="Empty"
 *   size="sm"
 * />
 */
import { computed, type Component } from 'vue'
import BaseIcon from './BaseIcon.vue'

interface Props {
  /** Icon component to display (from lucide-vue-next or heroicons) */
  icon?: Component
  /** Title text */
  title?: string
  /** Description text */
  description?: string
  /** Size variant */
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md'
})

// Map empty state size to icon size
const iconSize = computed(() => {
  const sizeMap: Record<string, 'md' | 'lg' | 'xl'> = {
    sm: 'md',
    md: 'lg',
    lg: 'xl'
  }
  return sizeMap[props.size]
})
</script>

<style scoped>
.base-empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: var(--space-8);
  flex: 1;
  min-height: 0;
}

/* Size variants */
.base-empty-state--sm {
  padding: var(--space-4);
  gap: var(--space-2);
}

.base-empty-state--md {
  padding: var(--space-8);
  gap: var(--space-3);
}

.base-empty-state--lg {
  padding: var(--space-11);
  gap: var(--space-4);
}

/* Icon container */
.base-empty-state__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-subtle);
  opacity: 0.6;
  transition: opacity var(--transition-normal);
}

.base-empty-state:hover .base-empty-state__icon {
  opacity: 0.8;
}

.base-empty-state--sm .base-empty-state__icon {
  width: 40px;
  height: 40px;
}

.base-empty-state--md .base-empty-state__icon {
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-2xl);
  padding: var(--space-3);
}

.base-empty-state--lg .base-empty-state__icon {
  width: 72px;
  height: 72px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: var(--radius-2xl);
  padding: var(--space-4);
}

/* Title */
.base-empty-state__title {
  margin: 0;
  font-weight: var(--font-weight-medium);
  color: var(--text-muted);
}

.base-empty-state--sm .base-empty-state__title {
  font-size: var(--font-size-sm);
  margin-top: var(--space-2);
}

.base-empty-state--md .base-empty-state__title {
  font-size: var(--font-size-md);
  margin-top: var(--space-4);
}

.base-empty-state--lg .base-empty-state__title {
  font-size: var(--font-size-lg);
  margin-top: var(--space-4);
}

/* Description */
.base-empty-state__description {
  margin: 0;
  color: var(--text-subtle);
  max-width: 400px;
}

.base-empty-state--sm .base-empty-state__description {
  font-size: var(--font-size-xs);
  margin-top: var(--space-1);
}

.base-empty-state--md .base-empty-state__description {
  font-size: var(--font-size-sm);
  margin-top: var(--space-1);
}

.base-empty-state--lg .base-empty-state__description {
  font-size: var(--font-size-md);
  margin-top: var(--space-2);
}

/* Action slot */
.base-empty-state__action {
  display: flex;
  gap: var(--space-3);
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
}

.base-empty-state--sm .base-empty-state__action {
  margin-top: var(--space-3);
}

.base-empty-state--md .base-empty-state__action {
  margin-top: var(--space-4);
}

.base-empty-state--lg .base-empty-state__action {
  margin-top: var(--space-6);
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .base-empty-state {
    padding: var(--space-6);
  }

  .base-empty-state--lg {
    padding: var(--space-8);
  }

  .base-empty-state__description {
    max-width: 300px;
  }
}
</style>
