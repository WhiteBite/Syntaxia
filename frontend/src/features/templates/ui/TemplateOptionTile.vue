<template>
  <!-- Grid Tile Toggle for Context Options -->
  <button 
    class="tpl-tile" 
    :class="{ active: modelValue }"
    @click="$emit('update:modelValue', !modelValue)"
  >
    <component :is="icon" class="w-4 h-4" />
    <span class="tpl-tile-label">{{ label }}</span>
    <span v-if="hint" class="tpl-tile-hint">{{ hint }}</span>
  </button>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  modelValue: boolean
  label: string
  icon: Component
  hint?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()
</script>

<style scoped>
.tpl-tile {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.375rem;
  padding: 0.75rem 0.5rem;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all 0.2s ease-out;
}

.tpl-tile:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.tpl-tile.active {
  background: var(--accent-indigo-bg);
  border-color: var(--accent-indigo-border);
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.2);
}

.tpl-tile.active svg {
  color: var(--accent-indigo);
}

.tpl-tile svg {
  color: var(--text-muted);
  transition: color 0.15s;
}

.tpl-tile-label {
  font-size: 11px;
  font-weight: 500;
  color: var(--text-primary);
}

.tpl-tile-hint {
  font-size: 9px;
  color: var(--text-subtle);
}
</style>
