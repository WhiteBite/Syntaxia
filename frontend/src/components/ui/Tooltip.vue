<template>
  <div class="tooltip-wrapper" @mouseenter="show = true" @mouseleave="show = false">
    <slot />
    <Transition name="tooltip">
      <div v-if="show" class="tooltip" :class="[`tooltip--${position}`]">
        {{ text }}
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  text: string
  position?: 'top' | 'bottom' | 'left' | 'right'
}>()

const show = ref(false)
</script>

<style scoped>
.tooltip-wrapper {
  position: relative;
  display: inline-flex;
}

.tooltip {
  position: absolute;
  z-index: 100;
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

.tooltip--top {
  bottom: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
}

.tooltip--bottom {
  top: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
}

.tooltip--left {
  right: calc(100% + 8px);
  top: 50%;
  transform: translateY(-50%);
}

.tooltip--right {
  left: calc(100% + 8px);
  top: 50%;
  transform: translateY(-50%);
}

/* Transition */
.tooltip-enter-active,
.tooltip-leave-active {
  transition: all 0.15s ease-out;
}

.tooltip-enter-from,
.tooltip-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(4px);
}

.tooltip--bottom.tooltip-enter-from,
.tooltip--bottom.tooltip-leave-to {
  transform: translateX(-50%) translateY(-4px);
}

.tooltip--left.tooltip-enter-from,
.tooltip--left.tooltip-leave-to {
  transform: translateY(-50%) translateX(4px);
}

.tooltip--right.tooltip-enter-from,
.tooltip--right.tooltip-leave-to {
  transform: translateY(-50%) translateX(-4px);
}
</style>
