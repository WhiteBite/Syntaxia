<template>
  <div class="space-y-1 p-2">
    <div v-for="i in 7" :key="i" class="flex items-center gap-2" :style="{ paddingLeft: getIndent(i) + 'px' }">
      <div class="skeleton w-4 h-4 rounded"></div>
      <div class="skeleton h-4 rounded" :style="{ width: getWidth(i) + 'px' }"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
function getIndent(index: number): number {
  const indents = [0, 0, 16, 16, 32, 16, 0]
  return indents[index - 1] || 0
}

function getWidth(index: number): number {
  const widths = [120, 80, 100, 140, 90, 110, 70]
  return widths[index - 1] || 100
}
</script>
