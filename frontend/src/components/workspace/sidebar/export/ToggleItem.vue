<template>
  <label class="toggle-row">
    <span class="toggle-label">
      {{ label }}
      <span 
        v-if="hint" 
        ref="triggerRef"
        class="hint-trigger" 
        @mouseenter="showHint = true" 
        @mouseleave="showHint = false"
      >
        <svg class="hint-icon" viewBox="0 0 16 16" fill="currentColor">
          <path d="M8 15A7 7 0 1 1 8 1a7 7 0 0 1 0 14zm0 1A8 8 0 1 0 8 0a8 8 0 0 0 0 16z"/>
          <path d="m8.93 6.588-2.29.287-.082.38.45.083c.294.07.352.176.288.469l-.738 3.468c-.194.897.105 1.319.808 1.319.545 0 1.178-.252 1.465-.598l.088-.416c-.2.176-.492.246-.686.246-.275 0-.375-.193-.304-.533L8.93 6.588zM9 4.5a1 1 0 1 1-2 0 1 1 0 0 1 2 0z"/>
        </svg>
      </span>
    </span>
    <button
      type="button"
      class="toggle-switch"
      :class="{ active: modelValue }"
      @click="emit('update:modelValue', !modelValue)"
    >
      <span class="toggle-thumb" />
    </button>
    
    <Teleport to="body">
      <Transition name="hint">
        <div 
          v-if="showHint && hint" 
          class="hint-popup"
          :style="popupStyle"
        >
          {{ hint }}
        </div>
      </Transition>
    </Teleport>
  </label>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

defineProps<{
  modelValue: boolean
  label: string
  hint?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const showHint = ref(false)
const triggerRef = ref<HTMLElement | null>(null)

const popupStyle = computed(() => {
  if (!triggerRef.value) return {}
  const rect = triggerRef.value.getBoundingClientRect()
  return {
    top: `${rect.bottom + 8}px`,
    left: `${rect.left + rect.width / 2}px`,
    transform: 'translateX(-50%)'
  }
})
</script>

<style scoped>
.toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 32px;
  padding: 0 0.5rem;
  margin: 0 -0.5rem;
  cursor: pointer;
  border-radius: 4px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.03);
  transition: all 0.15s ease-out;
}

.toggle-row:last-child {
  border-bottom: none;
}

.toggle-row:hover {
  background: rgba(255, 255, 255, 0.04);
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #9ca3af;
  transition: color 0.15s ease-out;
}

.toggle-row:hover .toggle-label {
  color: #e5e7eb;
}

.hint-trigger {
  display: inline-flex;
  align-items: center;
}

.hint-icon {
  width: 12px;
  height: 12px;
  color: #6b7280;
  opacity: 0.5;
  cursor: help;
  transition: all 0.15s ease-out;
}

.hint-trigger:hover .hint-icon {
  color: #a78bfa;
  opacity: 1;
}

.toggle-switch {
  position: relative;
  width: 32px;
  height: 16px;
  background: rgba(255, 255, 255, 0.08);
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease-out;
  flex-shrink: 0;
}

.toggle-switch:hover {
  background: rgba(255, 255, 255, 0.12);
}

.toggle-switch.active {
  background: #8b5cf6;
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 12px;
  height: 12px;
  background: white;
  border-radius: 50%;
  transition: transform 0.15s ease-out;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}

.toggle-switch.active .toggle-thumb {
  transform: translateX(16px);
}

.hint-enter-active,
.hint-leave-active {
  transition: all 0.15s ease-out;
}

.hint-enter-from,
.hint-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-4px);
}
</style>

<style>
/* Global styles for teleported popup */
.hint-popup {
  position: fixed;
  z-index: 10000;
  width: max-content;
  max-width: 260px;
  padding: 8px 12px;
  background: #1e2235;
  border: 1px solid rgba(139, 92, 246, 0.3);
  border-radius: 8px;
  color: #e5e7eb;
  font-size: 11px;
  font-weight: 400;
  line-height: 1.5;
  text-align: left;
  white-space: normal;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5), 0 0 0 1px rgba(255, 255, 255, 0.05);
  pointer-events: none;
}
</style>
