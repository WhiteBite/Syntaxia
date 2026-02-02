<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="[
      'base-button',
      `base-button--${variant}`,
      `base-button--${size}`,
      { 'base-button--loading': loading },
      { 'base-button--icon-only': iconOnly },
      { 'base-button--block': block },
      { 'base-button--outline': outline }
    ]"
    @click="$emit('click', $event)"
  >
    <!-- Icon Left or Loading Spinner -->
    <span 
      v-if="loading || (($slots.icon || icon) && iconPosition === 'left')" 
      class="base-button__icon base-button__icon--left"
    >
      <BaseSpinner v-if="loading" :size="spinnerSize" />
      <slot v-else name="icon">
        <BaseIcon v-if="icon" :icon="icon" :size="iconSize" />
      </slot>
    </span>

    <!-- Button Text -->
    <span v-if="!iconOnly" class="base-button__text">
      <slot />
    </span>

    <!-- Icon Right -->
    <span 
      v-if="!loading && ($slots.icon || icon) && iconPosition === 'right'" 
      class="base-button__icon base-button__icon--right"
    >
      <slot name="icon">
        <BaseIcon v-if="icon" :icon="icon" :size="iconSize" />
      </slot>
    </span>
  </button>
</template>

<script setup lang="ts">
/**
 * BaseButton - Enhanced button component with comprehensive features
 * 
 * @example
 * // Primary button with loading state
 * <BaseButton variant="primary" :loading="isLoading" :icon="PlayIcon">
 *   Run Tests
 * </BaseButton>
 * 
 * @example
 * // Outline button with icon on right
 * <BaseButton variant="outline" size="sm" :icon="SettingsIcon" icon-position="right">
 *   Settings
 * </BaseButton>
 * 
 * @example
 * // Link variant for cancel actions
 * <BaseButton variant="link" @click="cancel">
 *   Cancel
 * </BaseButton>
 * 
 * @example
 * // Full width success button
 * <BaseButton variant="success" block>
 *   Save Changes
 * </BaseButton>
 */
import { computed, type Component } from 'vue'
import BaseIcon from './BaseIcon.vue'
import BaseSpinner from './BaseSpinner.vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success' | 'warning' | 'outline' | 'link' | 'text'
  size?: 'xs' | 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  icon?: Component
  iconOnly?: boolean
  iconPosition?: 'left' | 'right'
  block?: boolean
  outline?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'secondary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  iconOnly: false,
  iconPosition: 'left',
  block: false,
  outline: false
})

defineEmits<{
  (e: 'click', event: MouseEvent): void
}>()

// Map button size to icon/spinner size
const iconSize = computed(() => {
  const sizeMap: Record<string, 'xs' | 'sm' | 'md' | 'lg'> = {
    xs: 'xs',
    sm: 'sm',
    md: 'sm',
    lg: 'md'
  }
  return sizeMap[props.size]
})

const spinnerSize = computed(() => iconSize.value)
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

.base-button--success {
  background: var(--color-success-soft);
  color: var(--color-success);
  border-color: var(--color-success-border);
}
.base-button--success:hover:not(:disabled) {
  background: rgba(74, 222, 128, 0.25);
  border-color: var(--color-success);
}

.base-button--warning {
  background: var(--color-warning-soft);
  color: var(--color-warning);
  border-color: var(--color-warning-border);
}
.base-button--warning:hover:not(:disabled) {
  background: rgba(251, 191, 36, 0.25);
  border-color: var(--color-warning);
}

.base-button--link {
  background: transparent;
  color: var(--accent-indigo);
  border: none;
  padding-left: 0;
  padding-right: 0;
  text-decoration: none;
}
.base-button--link:hover:not(:disabled) {
  text-decoration: underline;
  color: var(--accent-purple);
}

.base-button--text {
  background: transparent;
  color: var(--text-primary);
  border: none;
  padding-left: 0;
  padding-right: 0;
}
.base-button--text:hover:not(:disabled) {
  color: var(--text-primary);
  opacity: 0.8;
}

/* Outline variants */
.base-button--outline.base-button--primary {
  background: transparent;
  color: var(--accent-indigo);
  border-color: var(--accent-indigo);
  box-shadow: none;
}
.base-button--outline.base-button--primary:hover:not(:disabled) {
  background: var(--accent-indigo-bg);
  border-color: var(--accent-purple);
  color: var(--accent-purple);
  box-shadow: 0 0 20px rgba(139, 92, 246, 0.3);
}

.base-button--outline.base-button--secondary {
  background: transparent;
  color: var(--text-primary);
  border-color: var(--border-strong);
}
.base-button--outline.base-button--secondary:hover:not(:disabled) {
  background: var(--bg-2);
  border-color: var(--text-muted);
}

.base-button--outline.base-button--danger {
  background: transparent;
  color: var(--color-danger);
  border-color: var(--color-danger);
}
.base-button--outline.base-button--danger:hover:not(:disabled) {
  background: var(--color-danger-soft);
  border-color: var(--color-danger);
}

.base-button--outline.base-button--success {
  background: transparent;
  color: var(--color-success);
  border-color: var(--color-success);
}
.base-button--outline.base-button--success:hover:not(:disabled) {
  background: var(--color-success-soft);
  border-color: var(--color-success);
}

.base-button--outline.base-button--warning {
  background: transparent;
  color: var(--color-warning);
  border-color: var(--color-warning);
}
.base-button--outline.base-button--warning:hover:not(:disabled) {
  background: var(--color-warning-soft);
  border-color: var(--color-warning);
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

.base-button--icon-only.base-button--xs {
  padding: var(--space-1);
}

.base-button--icon-only.base-button--lg {
  padding: var(--space-3);
}

.base-button--block {
  width: 100%;
}

.base-button__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.base-button__text {
  flex: 1;
  min-width: 0;
}

.base-button:focus-visible {
  box-shadow: 0 0 0 2px var(--bg-0), 0 0 0 4px var(--accent-indigo);
}

/* Loading state - preserve button size */
.base-button--loading {
  position: relative;
}

.base-button--loading .base-button__text {
  visibility: visible;
  opacity: 0.7;
}

/* Ripple effect (optional) */
.base-button::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 0;
  height: 0;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%);
  transition: width 0.6s, height 0.6s;
  pointer-events: none;
}

.base-button:active:not(:disabled)::after {
  width: 200px;
  height: 200px;
  opacity: 0;
  transition: width 0s, height 0s, opacity 0.6s;
}
</style>
