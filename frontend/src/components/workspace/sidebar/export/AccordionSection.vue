<template>
  <div class="accordion">
    <button class="accordion-header" @click="isOpen = !isOpen">
      <ChevronRight class="accordion-chevron" :class="{ rotated: isOpen }" />
      <component :is="icon" class="accordion-icon" />
      <span class="accordion-label">{{ label }}</span>
    </button>
    <div v-if="isOpen" class="accordion-content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ChevronRight } from 'lucide-vue-next'
import { ref, type Component } from 'vue'

const props = withDefaults(defineProps<{
  label: string
  icon: Component
  defaultOpen?: boolean
}>(), {
  defaultOpen: true
})

const isOpen = ref(props.defaultOpen)
</script>

<style scoped>
.accordion {
  display: flex;
  flex-direction: column;
}

.accordion-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  width: 100%;
  padding: 0.75rem 0;
  background: transparent;
  border: none;
  color: #d1d5db;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  cursor: pointer;
  transition: color 0.15s;
}
.accordion-header:hover {
  color: white;
}

.accordion-chevron {
  width: 12px;
  height: 12px;
  color: #9ca3af;
  transition: transform 0.2s ease-out;
  flex-shrink: 0;
}
.accordion-chevron.rotated {
  transform: rotate(90deg);
}

.accordion-icon {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  color: #9ca3af;
}

.accordion-label {
  flex: 1;
  text-align: left;
}

.accordion-content {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0 0 0.5rem 0;
}
</style>
