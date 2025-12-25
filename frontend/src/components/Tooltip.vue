<template>
  <div class="tooltip-wrapper" @mouseenter="handleMouseEnter" @mouseleave="handleMouseLeave">
    <slot />
    
    <Teleport to="body">
      <Transition name="tooltip">
        <div
          v-if="isVisible && !disabled"
          ref="tooltipRef"
          role="tooltip"
          :aria-hidden="!isVisible"
          class="tooltip fixed z-[1070] px-3 py-2 bg-gray-900 border border-gray-600 rounded-lg shadow-2xl pointer-events-none max-w-[320px]"
          :style="tooltipStyle"
        >
          <slot name="tooltip">
            <div class="text-xs text-white whitespace-normal">{{ content }}</div>
          </slot>
          
          <!-- Arrow -->
          <div
            class="tooltip-arrow absolute"
            :class="arrowClass"
          ></div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'

interface Props {
  content?: string
  placement?: 'top' | 'bottom' | 'left' | 'right'
  delay?: number
  disabled?: boolean
  maxWidth?: string
}

const props = withDefaults(defineProps<Props>(), {
  content: '',
  placement: 'top',
  delay: 300,
  disabled: false,
  maxWidth: '320px'
})

const isVisible = ref(false)
const position = ref({ x: 0, y: 0 })
const actualPlacement = ref(props.placement)
const tooltipRef = ref<HTMLElement>() // Used in template via ref="tooltipRef"
const wrapperRef = ref<HTMLElement>()

let showTimer: ReturnType<typeof setTimeout> | null = null

const tooltipStyle = computed(() => {
  const style: Record<string, string> = {
    maxWidth: props.maxWidth
  }

  switch (actualPlacement.value) {
    case 'top':
      style.left = `${position.value.x}px`
      style.bottom = `${window.innerHeight - position.value.y}px`
      style.transform = 'translateX(-50%)'
      break
    case 'bottom':
      style.left = `${position.value.x}px`
      style.top = `${position.value.y}px`
      style.transform = 'translateX(-50%)'
      break
    case 'left':
      style.right = `${window.innerWidth - position.value.x}px`
      style.top = `${position.value.y}px`
      style.transform = 'translateY(-50%)'
      break
    case 'right':
      style.left = `${position.value.x}px`
      style.top = `${position.value.y}px`
      style.transform = 'translateY(-50%)'
      break
  }

  return style
})

const arrowClass = computed(() => {
  const classes: Record<string, string> = {
    top: 'left-1/2 -translate-x-1/2 top-full border-l-[6px] border-r-[6px] border-t-[6px] border-transparent border-t-gray-600',
    bottom: 'left-1/2 -translate-x-1/2 bottom-full border-l-[6px] border-r-[6px] border-b-[6px] border-transparent border-b-gray-600',
    left: 'top-1/2 -translate-y-1/2 left-full border-t-[6px] border-b-[6px] border-l-[6px] border-transparent border-l-gray-600',
    right: 'top-1/2 -translate-y-1/2 right-full border-t-[6px] border-b-[6px] border-r-[6px] border-transparent border-r-gray-600'
  }
  return classes[actualPlacement.value] || classes.top
})

function calculatePosition(element: HTMLElement) {
  const rect = element.getBoundingClientRect()
  const offset = 8

  let x = 0
  let y = 0
  let placement = props.placement

  switch (props.placement) {
    case 'top':
      x = rect.left + rect.width / 2
      y = rect.top - offset
      if (y < 100) {
        placement = 'bottom'
        y = rect.bottom + offset
      }
      break
    case 'bottom':
      x = rect.left + rect.width / 2
      y = rect.bottom + offset
      if (y + 200 > window.innerHeight) {
        placement = 'top'
        y = rect.top - offset
      }
      break
    case 'left':
      x = rect.left - offset
      y = rect.top + rect.height / 2
      break
    case 'right':
      x = rect.right + offset
      y = rect.top + rect.height / 2
      break
  }

  actualPlacement.value = placement
  position.value = { x, y }
}

function handleMouseEnter(event: MouseEvent) {
  if (props.disabled) return

  const target = event.currentTarget as HTMLElement
  
  if (showTimer) {
    clearTimeout(showTimer)
  }

  showTimer = setTimeout(() => {
    calculatePosition(target)
    isVisible.value = true
    showTimer = null
  }, props.delay)
}

function handleMouseLeave() {
  if (showTimer) {
    clearTimeout(showTimer)
    showTimer = null
  }
  isVisible.value = false
}

onMounted(() => {
  wrapperRef.value = document.querySelector('.tooltip-wrapper') as HTMLElement
})

// Expose refs used in template
defineExpose({ tooltipRef })
</script>

<style scoped>
.tooltip-wrapper {
  display: inline-block;
}

.tooltip-enter-active {
  transition: opacity 0.2s ease-out, transform 0.2s ease-out;
}

.tooltip-leave-active {
  transition: opacity 0.15s ease-in;
}

.tooltip-enter-from {
  opacity: 0;
  transform: translateX(-50%) translateY(-4px);
}

.tooltip-leave-to {
  opacity: 0;
}

.tooltip-arrow {
  width: 0;
  height: 0;
}
</style>
