<template>
  <component
    :is="icon"
    :class="[
      'base-icon',
      `base-icon--${size}`
    ]"
    :style="color ? { color } : undefined"
  />
</template>

<script setup lang="ts">
import type { Component } from 'vue'

interface Props {
  icon: Component
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl'
  color?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md'
})
</script>

<style scoped>
.base-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* Sizes */
.base-icon--xs {
  width: 12px;
  height: 12px;
}

.base-icon--sm {
  width: 16px;
  height: 16px;
}

.base-icon--md {
  width: 20px;
  height: 20px;
}

.base-icon--lg {
  width: 24px;
  height: 24px;
}

.base-icon--xl {
  width: 32px;
  height: 32px;
}
</style>
