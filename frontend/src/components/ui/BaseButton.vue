<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'base-button',
      `base-button--${variant}`,
      `base-button--${size}`,
      { 'base-button--loading': loading },
      { 'base-button--icon-only': iconOnly }
    ]"
    @click="$emit('click', $event)"
  >
    <span v-if="loading" class="base-button__loader">
      <svg class="animate-spin h-4 w-4" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" fill="none"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
      </svg>
    </span>
    <span v-else-if="$slots.icon || icon" class="base-button__icon">
      <slot name="icon">
        <component :is="icon" />
      </slot>
    </span>
    <span v-if="!iconOnly" class="base-button__text">
      <slot />
    </span>
  </button>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success' | 'warning'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  icon?: Component
  iconOnly?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  iconOnly: false
})

defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()
</script>

<style scoped>
.base-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  font-family: var(--font-sans);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  cursor: pointer;
  border: 1px solid transparent;
  white-space: nowrap;
  user-select: none;
  position: relative;
  outline: none;
}

.base-button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
  filter: grayscale(0.5);
}

/* Variants */
.base-button--primary {
  background: var(--gradient-primary);
  color: white;
  border: none;
  box-shadow: var(--shadow-glow-accent);
}
.base-button--primary:hover:not(:disabled) {
  background: var(--gradient-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

.base-button--secondary {
  background: var(--bg-2);
  color: var(--text-primary);
  border-color: var(--border-default);
}
.base-button--secondary:hover:not(:disabled) {
  background: var(--bg-3);
  border-color: var(--border-strong);
}

.base-button--ghost {
  background: transparent;
  color: var(--text-muted);
}
.base-button--ghost:hover:not(:disabled) {
  background: var(--bg-2);
  color: var(--text-primary);
}

.base-button--danger {
  background: var(--color-danger-soft);
  color: var(--color-danger);
  border-color: var(--color-danger-border);
}
.base-button--danger:hover:not(:disabled) {
  background: rgba(248, 113, 113, 0.25);
  border-color: var(--color-danger);
}

/* Sizes */
.base-button--xs {
  padding: var(--space-1) var(--space-2);
  font-size: var(--font-size-xs);
  min-height: calc(24px * var(--ui-scale));
}

.base-button--sm {
  padding: var(--space-1) var(--space-3);
  font-size: var(--font-size-sm);
  min-height: calc(32px * var(--ui-scale));
}

.base-button--md {
  padding: var(--space-2) var(--space-4);
  font-size: var(--font-size-md);
  min-height: calc(40px * var(--ui-scale));
}

.base-button--lg {
  padding: var(--space-3) var(--space-6);
  font-size: var(--font-size-lg);
  min-height: calc(48px * var(--ui-scale));
}

.base-button--icon-only {
  padding: var(--space-2);
}

.base-button__icon {
  display: flex;
  align-items: center;
  justify-content: center;
}

.base-button__icon :deep(svg) {
  width: 1.1em;
  height: 1.1em;
}

.base-button:focus-visible {
  box-shadow: 0 0 0 2px var(--bg-0), 0 0 0 4px var(--accent-indigo);
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}
.animate-spin {
  animation: spin 1s linear infinite;
}
</style>
