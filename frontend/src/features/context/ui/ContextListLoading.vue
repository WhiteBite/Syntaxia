<template>
  <div class="context-loading">
    <svg class="animate-spin h-6 w-6 text-purple-500" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
    </svg>
  </div>
</template>

<style scoped>
.context-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
}
</style>
