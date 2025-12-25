<template>
  <button
    @click="$emit('toggle', $event)"
    :class="[
      'filter-dropdown-item filter-item',
      { 'filter-dropdown-item-active': active },
      { 'filter-dropdown-item-excluded': excluded },
    ]"
    role="option"
    :aria-selected="active"
    tabindex="-1"
  >
    <div class="filter-item-left">
      <div v-if="iconType === 'category'" :class="['filter-item-icon', `filter-item-icon-${category}`]">
        <slot name="icon" />
      </div>
      <span v-else-if="icon" class="filter-item-emoji">{{ icon }}</span>
      <div v-if="subtitle" class="filter-item-info">
        <span :class="['filter-item-label', { 'filter-item-label-excluded': excluded }]">{{ label }}</span>
        <span class="filter-item-subtitle">{{ subtitle }}</span>
      </div>
      <template v-else>
        <span :class="['filter-item-label', { 'filter-item-label-excluded': excluded }]">{{ label }}</span>
        <span v-if="primary" class="filter-item-primary">★</span>
      </template>
    </div>
    <div class="filter-item-right">
      <span class="filter-item-count">{{ count }}</span>
      <div class="filter-item-bar">
        <div :class="['filter-item-bar-fill', barClass]" :style="{ width: percentage + '%' }"></div>
      </div>
    </div>
  </button>
</template>

<script setup lang="ts">
defineProps<{
  label: string
  icon?: string
  iconType?: 'emoji' | 'category'
  category?: string
  count: number
  percentage: number
  active?: boolean
  excluded?: boolean
  primary?: boolean
  subtitle?: string
  barClass?: string
}>()

defineEmits<{
  toggle: [event: MouseEvent]
}>()
</script>

<style scoped>
.filter-dropdown-item {
  @apply w-full flex items-center justify-between px-3 py-2 text-sm;
  color: var(--text-secondary);
  transition: all 150ms ease-out;
}
.filter-dropdown-item:hover,
.filter-dropdown-item:focus { background: var(--bg-2); color: var(--text-primary); outline: none; }
.filter-dropdown-item:focus-visible { box-shadow: inset 0 0 0 2px var(--accent-indigo); }
.filter-dropdown-item-active { background: var(--accent-indigo-bg); color: white; }
.filter-dropdown-item-active:hover { background: rgba(99, 102, 241, 0.35); }
.filter-dropdown-item-excluded { background: var(--color-danger-soft); color: var(--color-danger); }
.filter-dropdown-item-excluded:hover { background: rgba(248, 113, 113, 0.25); }

.filter-item-left { @apply flex items-center gap-2; }
.filter-item-icon { @apply w-6 h-6 rounded-md flex items-center justify-center; }
.filter-item-icon-code { background: var(--accent-indigo-bg); color: #a5b4fc; }
.filter-item-icon-test { background: var(--accent-emerald-bg); color: #6ee7b7; }
.filter-item-icon-config { background: var(--accent-orange-bg); color: #fdba74; }
.filter-item-icon-docs { background: var(--accent-purple-bg); color: #d8b4fe; }
.filter-item-icon-styles { background: rgba(236, 72, 153, 0.25); color: #f9a8d4; }
.filter-item-emoji { @apply text-base; }
.filter-item-label { @apply font-medium; }
.filter-item-label-excluded { text-decoration: line-through; opacity: 0.7; }
.filter-item-primary { @apply text-[10px] text-amber-400; }
.filter-item-info { @apply flex flex-col; }
.filter-item-subtitle { @apply text-[10px]; color: var(--text-muted); }

.filter-item-right { @apply flex items-center gap-2; }
.filter-item-count { @apply text-xs tabular-nums; color: var(--text-muted); }
.filter-item-bar { @apply w-12 h-1.5 rounded-full overflow-hidden; background: var(--bg-3); }
.filter-item-bar-fill {
  @apply h-full rounded-full;
  background: var(--accent-indigo);
  transition: width 300ms ease-out;
}
.filter-item-bar-fill.bar-lang { background: linear-gradient(90deg, var(--accent-indigo), var(--accent-purple)); }
.filter-item-bar-fill.bar-smart { background: linear-gradient(90deg, var(--accent-emerald), var(--accent-indigo)); }
</style>
