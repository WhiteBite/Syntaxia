<template>
  <span
    :class="[
      'base-badge',
      `base-badge--${variant}`,
      `base-badge--${size}`
    ]"
  >
    <span v-if="$slots.icon || icon" class="base-badge__icon">
      <slot name="icon">
        <component :is="icon" />
      </slot>
    </span>
    <span class="base-badge__text">
      <slot />
    </span>
  </span>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'info' | 'accent'
  size?: 'xs' | 'sm' | 'md'
  icon?: Component
}

withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'sm'
})
</script>

<style scoped>
.base-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  font-family: var(--font-sans);
  font-weight: var(--font-weight-semibold);
  border-radius: var(--radius-full);
  line-height: 1;
  white-space: nowrap;
}

/* Variants */
.base-badge--primary {
  background: var(--accent-purple-bg);
  color: white;
  border: 1px solid var(--accent-purple-border);
}

.base-badge--secondary {
  background: var(--bg-2);
  color: var(--text-secondary);
  border: 1px solid var(--border-default);
}

.base-badge--accent {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.base-badge--success {
  background: var(--color-success-soft);
  color: var(--color-success);
  border: 1px solid var(--color-success-border);
}

.base-badge--warning {
  background: var(--color-warning-soft);
  color: var(--color-warning);
  border: 1px solid var(--color-warning-border);
}

.base-badge--danger {
  background: var(--color-danger-soft);
  color: var(--color-danger);
  border: 1px solid var(--color-danger-border);
}

.base-badge--info {
  background: var(--color-info-soft);
  color: var(--color-info);
  border: 1px solid var(--color-info-border);
}

/* Sizes */
.base-badge--xs {
  padding: 1px 4px;
  font-size: calc(9px * var(--ui-scale));
}

.base-badge--sm {
  padding: 2px 8px;
  font-size: var(--font-size-xs);
}

.base-badge--md {
  padding: 4px 12px;
  font-size: var(--font-size-sm);
}

.base-badge__icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.base-badge__icon :deep(svg) {
  width: 1em;
  height: 1em;
}
</style>
