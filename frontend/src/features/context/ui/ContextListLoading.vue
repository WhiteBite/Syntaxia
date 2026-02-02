<template>
  <div class="context-loading">
    <BaseSpinner size="md" class="text-purple-500" />
  </div>
</template>

<script setup lang="ts">
import BaseSpinner from '@/components/ui/BaseSpinner.vue'
</script>

<style scoped>
.context-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}
</style>
