<template>
  <div
    :class="[
      'base-card',
      { 'base-card--interactive': interactive },
      { 'base-card--glass': glass }
    ]"
  >
    <div v-if="$slots.header" class="base-card__header">
      <slot name="header" />
    </div>
    <div class="base-card__body">
      <slot />
    </div>
    <div v-if="$slots.footer" class="base-card__footer">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  interactive?: boolean
  glass?: boolean
}

withDefaults(defineProps<Props>(), {
  interactive: false,
  glass: false
})
</script>

<style scoped>
.base-card {
  background: var(--bg-1);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xl);
  overflow: hidden;
  transition: all var(--transition-normal);
}

.base-card--glass {
  background: rgba(17, 24, 39, 0.7);
  backdrop-filter: blur(12px);
  border-color: rgba(255, 255, 255, 0.08);
}

.base-card--interactive {
  cursor: pointer;
}

.base-card--interactive:hover {
  background: var(--bg-2);
  border-color: var(--border-default);
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.base-card__header {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--border-subtle);
}

.base-card__body {
  padding: var(--space-4);
}

.base-card__footer {
  padding: var(--space-3) var(--space-4);
  border-top: 1px solid var(--border-subtle);
  background: rgba(0, 0, 0, 0.1);
}
</style>
