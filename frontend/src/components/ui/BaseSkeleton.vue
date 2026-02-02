<template>
  <div
    v-if="count > 1"
    class="base-skeleton-group"
    :style="{ gap: `${gap}px` }"
  >
    <div
      v-for="i in count"
      :key="i"
      :class="skeletonClasses"
      :style="skeletonStyles"
      role="status"
      aria-busy="true"
      :aria-label="ariaLabel"
    />
  </div>
  <div
    v-else
    :class="skeletonClasses"
    :style="skeletonStyles"
    role="status"
    aria-busy="true"
    :aria-label="ariaLabel"
  />
</template>

<script setup lang="ts">
/**
 * BaseSkeleton - Skeleton loading component with shimmer animation
 * 
 * @example
 * // Text skeleton
 * <BaseSkeleton variant="text" width="200px" />
 * 
 * @example
 * // Multiple lines
 * <BaseSkeleton variant="text" :count="3" />
 * 
 * @example
 * // Circle (avatar)
 * <BaseSkeleton variant="circle" width="40px" height="40px" />
 * 
 * @example
 * // Card skeleton
 * <BaseSkeleton variant="card" height="200px" />
 * 
 * @example
 * // Rectangle without animation
 * <BaseSkeleton variant="rect" width="100%" height="100px" :animated="false" />
 */
import { computed } from 'vue'

interface Props {
  variant?: 'text' | 'circle' | 'rect' | 'card'
  width?: string | number
  height?: string | number
  count?: number
  animated?: boolean
  gap?: number
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'text',
  width: '100%',
  height: undefined,
  count: 1,
  animated: true,
  gap: 8
})

const skeletonClasses = computed(() => [
  'base-skeleton',
  `base-skeleton--${props.variant}`,
  { 'base-skeleton--animated': props.animated }
])

const skeletonStyles = computed(() => {
  const styles: Record<string, string> = {}
  
  // Width
  if (props.width !== undefined) {
    styles.width = typeof props.width === 'number' ? `${props.width}px` : props.width
  }
  
  // Height
  if (props.height !== undefined) {
    styles.height = typeof props.height === 'number' ? `${props.height}px` : props.height
  } else {
    // Default heights based on variant
    if (props.variant === 'text') {
      styles.height = '1rem'
    } else if (props.variant === 'circle') {
      styles.height = styles.width || '40px'
    } else if (props.variant === 'card') {
      styles.height = '200px'
    }
  }
  
  return styles
})

const ariaLabel = computed(() => {
  return props.count > 1 
    ? `Loading ${props.count} items` 
    : 'Loading content'
})
</script>

<style scoped>
.base-skeleton {
  display: block;
  background: linear-gradient(
    90deg,
    rgba(55, 65, 81, 0.4) 0%,
    rgba(75, 85, 99, 0.4) 50%,
    rgba(55, 65, 81, 0.4) 100%
  );
  background-size: 200% 100%;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.base-skeleton--animated {
  animation: skeleton-shimmer 1.5s ease-in-out infinite;
}

.base-skeleton--text {
  border-radius: var(--radius-sm);
  height: 1rem;
}

.base-skeleton--circle {
  border-radius: 50%;
  aspect-ratio: 1;
}

.base-skeleton--rect {
  border-radius: var(--radius-md);
}

.base-skeleton--card {
  border-radius: var(--radius-lg);
  height: 200px;
}

.base-skeleton-group {
  display: flex;
  flex-direction: column;
}

@keyframes skeleton-shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* CSS Variables for customization */
.base-skeleton {
  --skeleton-bg-start: rgba(55, 65, 81, 0.4);
  --skeleton-bg-mid: rgba(75, 85, 99, 0.4);
  --skeleton-bg-end: rgba(55, 65, 81, 0.4);
  
  background: linear-gradient(
    90deg,
    var(--skeleton-bg-start) 0%,
    var(--skeleton-bg-mid) 50%,
    var(--skeleton-bg-end) 100%
  );
}
</style>
