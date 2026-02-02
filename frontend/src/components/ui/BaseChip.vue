<template>
  <component
    :is="clickable ? 'button' : 'span'"
    :type="clickable ? 'button' : undefined"
    :role="clickable ? 'button' : undefined"
    :tabindex="clickable ? 0 : undefined"
    :class="[
      'base-chip',
      `base-chip--${variant}`,
      `base-chip--${size}`,
      { 'base-chip--clickable': clickable },
      { 'base-chip--removable': removable }
    ]"
    @click="handleClick"
    @keydown="handleKeydown"
  >
    <!-- Icon Slot -->
    <span v-if="$slots.icon || icon" class="base-chip__icon">
      <slot name="icon">
        <component :is="icon" v-if="icon" />
      </slot>
    </span>

    <!-- Text Content -->
    <span class="base-chip__text">
      <slot />
    </span>

    <!-- Remove Button -->
    <button
      v-if="removable"
      type="button"
      class="base-chip__remove"
      :aria-label="t('common.remove')"
      @click.stop="handleRemove"
      @keydown.stop="handleRemoveKeydown"
    >
      <X class="base-chip__remove-icon" />
    </button>
  </component>
</template>

<script setup lang="ts">
/**
 * BaseChip - Compact chip/tag component for labels, filters, and selections
 * 
 * @example
 * // Simple chip
 * <BaseChip variant="primary">JavaScript</BaseChip>
 * 
 * @example
 * // Removable chip
 * <BaseChip variant="success" removable @remove="handleRemove">
 *   Selected
 * </BaseChip>
 * 
 * @example
 * // Clickable chip with icon
 * <BaseChip clickable :icon="FilterIcon" @click="handleClick">
 *   Filter
 * </BaseChip>
 * 
 * @example
 * // Custom icon via slot
 * <BaseChip variant="warning">
 *   <template #icon>
 *     <span>⚠️</span>
 *   </template>
 *   Warning
 * </BaseChip>
 */
import { type Component } from 'vue'
import { X } from 'lucide-vue-next'
import { useI18n } from '@/composables/useI18n'

interface Props {
  variant?: 'default' | 'primary' | 'success' | 'warning' | 'danger'
  size?: 'xs' | 'sm' | 'md'
  removable?: boolean
  clickable?: boolean
  icon?: Component
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'default',
  size: 'sm',
  removable: false,
  clickable: false
})

const emit = defineEmits<{
  (e: 'remove'): void
  (e: 'click', event: MouseEvent): void
}>()

const { t } = useI18n()

function handleClick(event: MouseEvent) {
  if (props.clickable) {
    emit('click', event)
  }
}

function handleRemove() {
  emit('remove')
}

function handleKeydown(event: KeyboardEvent) {
  if (props.clickable && (event.key === 'Enter' || event.key === ' ')) {
    event.preventDefault()
    emit('click', event as unknown as MouseEvent)
  }
  if (props.removable && event.key === 'Delete') {
    event.preventDefault()
    emit('remove')
  }
}

function handleRemoveKeydown(event: KeyboardEvent) {
  if (event.key === 'Enter' || event.key === ' ') {
    event.preventDefault()
    emit('remove')
  }
}
</script>

<style scoped>
.base-chip {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-family: var(--font-sans);
  font-weight: var(--font-weight-medium);
  border-radius: var(--radius-full);
  line-height: 1;
  white-space: nowrap;
  border: 1px solid transparent;
  transition: all var(--transition-fast);
  user-select: none;
}

/* Variants */
.base-chip--default {
  background: var(--bg-2);
  color: var(--text-secondary);
  border-color: var(--border-default);
}

.base-chip--primary {
  background: var(--accent-indigo-bg);
  color: white;
  border-color: var(--accent-indigo-border);
}

.base-chip--success {
  background: var(--color-success-soft);
  color: var(--color-success);
  border-color: var(--color-success-border);
}

.base-chip--warning {
  background: var(--color-warning-soft);
  color: var(--color-warning);
  border-color: var(--color-warning-border);
}

.base-chip--danger {
  background: var(--color-danger-soft);
  color: var(--color-danger);
  border-color: var(--color-danger-border);
}

/* Sizes */
.base-chip--xs {
  padding: 1px 6px;
  font-size: calc(10px * var(--ui-scale));
  gap: 2px;
}

.base-chip--sm {
  padding: 2px 8px;
  font-size: var(--font-size-xs);
}

.base-chip--md {
  padding: 4px 12px;
  font-size: var(--font-size-sm);
  gap: var(--space-2);
}

/* Clickable state */
.base-chip--clickable {
  cursor: pointer;
}

.base-chip--clickable:hover {
  transform: scale(1.05);
  filter: brightness(1.1);
}

.base-chip--clickable:active {
  transform: scale(0.98);
}

.base-chip--clickable:focus-visible {
  outline: 2px solid var(--accent-indigo);
  outline-offset: 2px;
}

/* Removable state */
.base-chip--removable {
  padding-right: 2px;
}

.base-chip--removable.base-chip--xs {
  padding-right: 1px;
}

.base-chip--removable.base-chip--md {
  padding-right: 4px;
}

/* Icon */
.base-chip__icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.base-chip__icon :deep(svg) {
  width: 1em;
  height: 1em;
}

.base-chip--xs .base-chip__icon :deep(svg) {
  width: 10px;
  height: 10px;
}

.base-chip--sm .base-chip__icon :deep(svg) {
  width: 12px;
  height: 12px;
}

.base-chip--md .base-chip__icon :deep(svg) {
  width: 14px;
  height: 14px;
}

/* Text */
.base-chip__text {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Remove button */
.base-chip__remove {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px;
  border: none;
  background: transparent;
  color: currentColor;
  cursor: pointer;
  border-radius: var(--radius-sm);
  opacity: 0.6;
  transition: all var(--transition-fast);
  flex-shrink: 0;
}

.base-chip__remove:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.2);
}

.base-chip__remove:focus-visible {
  outline: 2px solid currentColor;
  outline-offset: 1px;
  opacity: 1;
}

.base-chip__remove-icon {
  width: 12px;
  height: 12px;
}

.base-chip--xs .base-chip__remove-icon {
  width: 10px;
  height: 10px;
}

.base-chip--md .base-chip__remove-icon {
  width: 14px;
  height: 14px;
}

/* Hover effects for variants */
.base-chip--clickable.base-chip--default:hover {
  background: var(--bg-3);
  border-color: var(--border-strong);
}

.base-chip--clickable.base-chip--primary:hover {
  background: rgba(99, 102, 241, 0.35);
  border-color: var(--accent-purple-border);
}

.base-chip--clickable.base-chip--success:hover {
  background: rgba(74, 222, 128, 0.25);
  border-color: var(--color-success);
}

.base-chip--clickable.base-chip--warning:hover {
  background: rgba(251, 191, 36, 0.25);
  border-color: var(--color-warning);
}

.base-chip--clickable.base-chip--danger:hover {
  background: rgba(248, 113, 113, 0.25);
  border-color: var(--color-danger);
}
</style>
