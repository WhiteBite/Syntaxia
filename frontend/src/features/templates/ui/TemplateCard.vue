<template>
  <!-- Builder Card Component -->
  <div 
    class="tpl-card" 
    :class="{ collapsed: isCollapsed, disabled: !enabled }"
  >
    <div class="tpl-card-header" @click="enabled && toggleCollapse()">
      <component :is="icon" class="w-3.5 h-3.5" />
      <span class="tpl-card-title">{{ title }}</span>
      <span v-if="count !== undefined" class="tpl-card-count">{{ count }}</span>
      <button 
        v-if="toggleable" 
        @click.stop="$emit('toggle')" 
        class="tpl-card-toggle"
        :class="{ active: enabled }"
      >
        <Check v-if="enabled" class="w-3 h-3" />
      </button>
      <ChevronRight v-if="collapsible && enabled" class="w-3.5 h-3.5 tpl-card-chevron" :class="{ rotated: !isCollapsed }" />
    </div>
    
    <Transition name="expand">
      <div v-if="!isCollapsed && enabled" class="tpl-card-body">
        <slot />
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { Check, ChevronRight } from 'lucide-vue-next'
import type { Component } from 'vue'
import { ref } from 'vue'

defineProps<{
  title: string
  icon: Component
  count?: number
  enabled?: boolean
  toggleable?: boolean
  collapsible?: boolean
}>()

defineEmits<{
  (e: 'toggle'): void
}>()

const isCollapsed = ref(false)

function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
}
</script>

<style scoped>
.tpl-card {
  background: #0d1117;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: var(--radius-md);
  overflow: hidden;
  transition: all 0.2s ease-out;
}

.tpl-card:hover:not(.disabled) {
  border-color: rgba(255, 255, 255, 0.15);
}

.tpl-card:focus-within:not(.disabled) {
  border-color: var(--accent-indigo-border);
  box-shadow: 0 0 0 1px var(--accent-indigo-border);
}

.tpl-card.disabled {
  opacity: 0.4;
}

.tpl-card.collapsed .tpl-card-header {
  border-bottom: none;
}

.tpl-card-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 0.75rem;
  background: rgba(255, 255, 255, 0.02);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  user-select: none;
  transition: background 0.15s;
}

.tpl-card-header:hover {
  background: rgba(255, 255, 255, 0.04);
}

.tpl-card-header svg {
  color: var(--text-muted);
  flex-shrink: 0;
}

.tpl-card-title {
  flex: 1;
  font-size: 11px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.03em;
}

.tpl-card-count {
  font-size: 10px;
  font-family: var(--font-mono);
  color: var(--text-subtle);
}

.tpl-card-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 1.25rem;
  height: 1.25rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  color: transparent;
  cursor: pointer;
  transition: all 0.15s;
}

.tpl-card-toggle:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.tpl-card-toggle.active {
  background: var(--accent-indigo-bg);
  border-color: var(--accent-indigo-border);
  color: var(--accent-indigo);
}

.tpl-card-chevron {
  color: var(--text-subtle);
  transition: transform 0.2s ease-out;
}

.tpl-card-chevron.rotated {
  transform: rotate(90deg);
}

.tpl-card-body {
  padding: 0.75rem;
}

/* Expand animation */
.expand-enter-active,
.expand-leave-active {
  transition: all 0.2s ease-out;
  overflow: hidden;
}

.expand-enter-from,
.expand-leave-to {
  opacity: 0;
  max-height: 0;
  padding: 0 0.75rem;
}

.expand-enter-to,
.expand-leave-from {
  opacity: 1;
  max-height: 500px;
}
</style>
