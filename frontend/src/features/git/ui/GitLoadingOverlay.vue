<template>
  <div v-if="isLoading" class="loading-overlay absolute">
    <div class="text-center">
      <BaseSpinner size="lg" class="mx-auto mb-2" />
      <span class="text-sm text-gray-400">{{ message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { BaseSpinner } from '@/components/ui'

defineProps<{
  isLoading: boolean
  message: string
}>()
</script>
