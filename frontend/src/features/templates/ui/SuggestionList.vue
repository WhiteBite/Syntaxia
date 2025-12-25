<template>
  <ul class="suggestion-list">
    <li
      v-for="(suggestion, index) in suggestions"
      :key="suggestion.key"
      class="suggestion-item"
      :class="{ selected: index === selectedIndex }"
      @click="selectSuggestion(suggestion)"
      @mouseenter="hoveredIndex = index"
    >
      <span class="suggestion-key">{{ suggestion.key }}</span>
      <div class="suggestion-info">
        <div class="suggestion-label">{{ suggestion.label }}</div>
        <div class="suggestion-description">{{ suggestion.description }}</div>
      </div>
    </li>
  </ul>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { TemplateVariable } from '../model/template.types'

defineProps<{
  suggestions: TemplateVariable[]
  selectedIndex: number
}>()

const emit = defineEmits<{
  (e: 'select', suggestion: TemplateVariable): void
}>()

const hoveredIndex = ref(-1)

function selectSuggestion(suggestion: TemplateVariable) {
  emit('select', suggestion)
}
</script>

<style scoped>
.suggestion-list {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  z-index: 10;
  background: var(--bg-2);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  max-height: 200px;
  overflow-y: auto;
  box-shadow: var(--shadow-md);
  list-style: none;
  margin: 0;
  padding: 0;
}

.suggestion-item {
  padding: 0.5rem;
  cursor: pointer;
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: flex-start;
  gap: 0.5rem;
}

.suggestion-item:hover,
.suggestion-item.selected {
  background: var(--accent-indigo-bg);
}

.suggestion-item:last-child {
  border-bottom: none;
}

.suggestion-key {
  font-family: monospace;
  background: var(--bg-3);
  padding: 0.125rem 0.25rem;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
  flex-shrink: 0;
  align-self: center;
}

.suggestion-info {
  flex: 1;
}

.suggestion-label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: var(--font-size-sm);
}

.suggestion-description {
  font-size: var(--font-size-xs);
  color: var(--text-muted);
  margin-top: 0.125rem;
}
</style>