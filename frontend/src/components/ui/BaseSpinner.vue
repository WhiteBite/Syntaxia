<template>
  <svg
    :class="[
      'base-spinner',
      `base-spinner--${size}`,
      'animate-spin'
    ]"
    :style="{ color }"
    viewBox="0 0 24 24"
    fill="none"
  >
    <circle
      class="opacity-25"
      cx="12"
      cy="12"
      r="10"
      stroke="currentColor"
      stroke-width="2"
    />
    <path
      class="opacity-75"
      fill="currentColor"
      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
    />
  </svg>
</template>

<script setup lang="ts">
interface Props {
  size?: 'xs' | 'sm' | 'md' | 'lg'
  color?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md',
  color: 'currentColor'
})
</script>

<style scoped>
.base-spinner {
  display: inline-block;
  flex-shrink: 0;
}

/* Sizes */
.base-spinner--xs {
  width: 12px;
  height: 12px;
}

.base-spinner--sm {
  width: 16px;
  height: 16px;
}

.base-spinner--md {
  width: 20px;
  height: 20px;
}

.base-spinner--lg {
  width: 24px;
  height: 24px;
}

/* Animation */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
