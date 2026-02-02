<template>
  <Transition name="alert">
    <div
      v-if="!isDismissed"
      :class="[
        'base-alert',
        `base-alert--${variant}`
      ]"
      role="alert"
      :aria-live="variant === 'error' ? 'assertive' : 'polite'"
    >
      <!-- Icon -->
      <div class="base-alert__icon">
        <slot name="icon">
          <component :is="defaultIcon" :size="16" />
        </slot>
      </div>

      <!-- Content -->
      <div class="base-alert__content">
        <!-- Title -->
        <div v-if="$slots.title || title" class="base-alert__title">
          <slot name="title">{{ title }}</slot>
        </div>

        <!-- Message -->
        <div class="base-alert__message">
          <slot />
        </div>

        <!-- Actions -->
        <div v-if="$slots.actions" class="base-alert__actions">
          <slot name="actions" />
        </div>
      </div>

      <!-- Dismiss Button -->
      <button
        v-if="dismissible"
        class="base-alert__dismiss"
        @click="handleDismiss"
        :aria-label="dismissLabel"
        type="button"
      >
        <X :size="16" />
      </button>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Info, CheckCircle, AlertTriangle, XCircle, X } from 'lucide-vue-next'

interface Props {
  variant?: 'info' | 'success' | 'warning' | 'error'
  title?: string
  dismissible?: boolean
  dismissLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'info',
  dismissible: false,
  dismissLabel: 'Dismiss'
})

const emit = defineEmits<{
  (e: 'dismiss'): void
}>()

const isDismissed = ref(false)

const defaultIcon = computed(() => {
  const iconMap = {
    info: Info,
    success: CheckCircle,
    warning: AlertTriangle,
    error: XCircle
  }
  return iconMap[props.variant]
})

function handleDismiss() {
  isDismissed.value = true
  emit('dismiss')
}
</script>

<style scoped>
.base-alert {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
  border-radius: var(--radius-xl);
  border: 1px solid;
  font-size: var(--font-size-sm);
  line-height: 1.5;
  position: relative;
}

/* Variants */
.base-alert--info {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
  color: #93c5fd;
}

.base-alert--info .base-alert__icon {
  color: #3b82f6;
}

.base-alert--info .base-alert__title {
  color: #60a5fa;
}

.base-alert--success {
  background: var(--color-success-soft);
  border-color: var(--color-success-border);
  color: #86efac;
}

.base-alert--success .base-alert__icon {
  color: var(--color-success);
}

.base-alert--success .base-alert__title {
  color: var(--color-success);
}

.base-alert--warning {
  background: var(--color-warning-soft);
  border-color: var(--color-warning-border);
  color: #fcd34d;
}

.base-alert--warning .base-alert__icon {
  color: var(--color-warning);
}

.base-alert--warning .base-alert__title {
  color: var(--color-warning);
}

.base-alert--error {
  background: var(--color-danger-soft);
  border-color: var(--color-danger-border);
  color: #fca5a5;
}

.base-alert--error .base-alert__icon {
  color: var(--color-danger);
}

.base-alert--error .base-alert__title {
  color: var(--color-danger);
}

/* Icon */
.base-alert__icon {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  padding-top: 2px;
}

/* Content */
.base-alert__content {
  flex: 1;
  min-width: 0;
}

.base-alert__title {
  font-weight: var(--font-weight-semibold);
  margin-bottom: var(--space-1);
}

.base-alert__message {
  color: inherit;
  opacity: 0.9;
}

.base-alert__actions {
  margin-top: var(--space-3);
  display: flex;
  gap: var(--space-2);
}

/* Dismiss Button */
.base-alert__dismiss {
  flex-shrink: 0;
  display: flex;
  align-items: flex-start;
  padding: var(--space-1);
  margin: calc(var(--space-1) * -1);
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: currentColor;
  opacity: 0.6;
  cursor: pointer;
  transition: all var(--transition-fast);
}

.base-alert__dismiss:hover {
  opacity: 1;
  background: rgba(0, 0, 0, 0.1);
}

.base-alert__dismiss:focus-visible {
  outline: 2px solid currentColor;
  outline-offset: 2px;
  opacity: 1;
}

/* Transitions */
.alert-enter-active,
.alert-leave-active {
  transition: all var(--transition-normal);
}

.alert-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}

.alert-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
