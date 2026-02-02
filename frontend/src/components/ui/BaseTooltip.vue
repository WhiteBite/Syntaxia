<template>
  <div 
    ref="triggerRef"
    class="tooltip-wrapper" 
    @mouseenter="show = true" 
    @mouseleave="show = false"
  >
    <slot />
    <Teleport to="body">
      <Transition name="tooltip">
        <div 
          v-if="show" 
          class="tooltip" 
          :class="[`tooltip--${position}`, { 'tooltip--multiline': multiline }]"
          :style="tooltipStyle"
        >
          {{ text }}
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = withDefaults(defineProps<{
  text: string
  position?: 'top' | 'bottom' | 'left' | 'right'
  multiline?: boolean
}>(), {
  position: 'top',
  multiline: false
})

const show = ref(false)
const triggerRef = ref<HTMLElement | null>(null)

const tooltipStyle = computed(() => {
  if (!triggerRef.value) return {}
  
  const rect = triggerRef.value.getBoundingClientRect()
  const style: Record<string, string> = {
    position: 'fixed',
    zIndex: '10000'
  }
  
  switch (props.position) {
    case 'top':
      style.top = `${rect.top - 8}px`
      style.left = `${rect.left + rect.width / 2}px`
      style.transform = 'translate(-50%, -100%)'
      break
    case 'bottom':
      style.top = `${rect.bottom + 8}px`
      style.left = `${rect.left + rect.width / 2}px`
      style.transform = 'translateX(-50%)'
      break
    case 'left':
      style.top = `${rect.top + rect.height / 2}px`
      style.left = `${rect.left - 8}px`
      style.transform = 'translate(-100%, -50%)'
      break
    case 'right':
      style.top = `${rect.top + rect.height / 2}px`
      style.left = `${rect.right + 8}px`
      style.transform = 'translateY(-50%)'
      break
  }
  
  return style
})
</script>

<style scoped>
.tooltip-wrapper {
  display: inline-flex;
}
</style>

<style>
/* Global styles for teleported tooltip */
.tooltip {
  width: max-content;
  padding: 6px 10px;
  background: #1c1f2e;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 6px;
  color: #e5e7eb;
  font-size: 11px;
  font-weight: 500;
  white-space: nowrap;
  pointer-events: none;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
}

.tooltip--multiline {
  white-space: normal;
  max-width: min(280px, 40vw);
  text-align: left;
  line-height: 1.4;
}

/* Transition */
.tooltip-enter-active,
.tooltip-leave-active {
  transition: all 0.15s ease-out;
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
}

.tooltip--top.tooltip-enter-from,
.tooltip--top.tooltip-leave-to {
  transform: translate(-50%, calc(-100% + 4px));
}

.tooltip--bottom.tooltip-enter-from,
.tooltip--bottom.tooltip-leave-to {
  transform: translate(-50%, -4px);
}

.tooltip--left.tooltip-enter-from,
.tooltip--left.tooltip-leave-to {
  transform: translate(calc(-100% + 4px), -50%);
}

.tooltip--right.tooltip-enter-from,
.tooltip--right.tooltip-leave-to {
  transform: translate(-4px, -50%);
}
</style>
