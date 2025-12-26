<template>
  <div v-if="isLoading" class="loading-overlay absolute">
    <div class="text-center">
      <svg class="loading-spinner mx-auto mb-2" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
      </svg>
      <span class="text-sm text-gray-400">{{ message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  isLoading: boolean
  message: string
}>()
</script>
