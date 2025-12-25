<template>
  <button
    @click="$emit('remove')"
    :class="[
      'filter-chip',
      excluded ? 'filter-chip-excluded' : `filter-chip-${category}`,
    ]"
  >
    <span v-if="excluded" class="chip-icon">⊘</span>
    <span v-else-if="icon" class="chip-icon">{{ icon }}</span>
    <span :class="{ 'chip-label-excluded': excluded }">{{ label }}</span>
    <svg class="w-3 h-3 chip-close" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
    </svg>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  icon?: string
  category?: string
  excluded?: boolean
}>()

defineEmits<{
  remove: []
}>()
</script>

<style scoped>
.filter-chip {
  @apply flex items-center gap-1 px-2 py-1 text-xs font-medium rounded-full;
  transition: all 150ms ease-out;
}
.filter-chip:hover { transform: scale(1.05); }
.filter-chip:hover .chip-close { opacity: 1; }
.chip-close { opacity: 0.6; transition: opacity 150ms ease-out; }
.chip-icon { @apply text-sm leading-none; }
.chip-label-excluded { text-decoration: line-through; opacity: 0.8; }

.filter-chip-code { background: var(--accent-indigo-bg); border: 1px solid var(--accent-indigo-border); color: white; }
.filter-chip-test { background: var(--accent-emerald-bg); border: 1px solid var(--accent-emerald-border); color: white; }
.filter-chip-config { background: var(--accent-orange-bg); border: 1px solid var(--accent-orange-border); color: white; }
.filter-chip-docs { background: var(--accent-purple-bg); border: 1px solid var(--accent-purple-border); color: white; }
.filter-chip-styles { background: rgba(236, 72, 153, 0.25); border: 1px solid rgba(236, 72, 153, 0.5); color: white; }
.filter-chip-lang { background: linear-gradient(135deg, var(--accent-indigo-bg), var(--accent-purple-bg)); border: 1px solid var(--accent-purple-border); color: white; }
.filter-chip-smart { background: linear-gradient(135deg, var(--accent-emerald-bg), var(--accent-indigo-bg)); border: 1px solid var(--accent-emerald-border); color: white; }
.filter-chip-excluded { background: var(--color-danger-soft); border: 1px solid var(--color-danger-border); color: var(--color-danger); }
</style>
