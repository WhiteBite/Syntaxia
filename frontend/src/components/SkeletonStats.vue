<template>
  <div class="stats-grid">
    <div v-for="i in 4" :key="i" class="stat-card">
      <div class="skeleton h-6 w-16 rounded mb-1"></div>
      <div class="skeleton h-3 w-12 rounded"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Skeleton stats component - uses .skeleton class from design-tokens.css
</script>
