<template>
  <div class="base-popover-wrapper">
    <!-- Trigger Slot -->
    <div
      ref="triggerRef"
      class="base-popover-trigger"
      @click="handleTriggerClick"
      @mouseenter="handleTriggerMouseEnter"
      @mouseleave="handleTriggerMouseLeave"
    >
      <slot name="trigger" />
    </div>

    <!-- Popover Content via Teleport -->
    <Teleport to="body">
      <Transition name="popover-fade">
        <div
          v-if="isOpen"
          ref="popoverRef"
          class="base-popover"
          :class="[`base-popover--${placement}`, { 'base-popover--no-arrow': !arrow }]"
          :style="popoverStyle"
          @mouseenter="handlePopoverMouseEnter"
          @mouseleave="handlePopoverMouseLeave"
        >
          <!-- Arrow -->
          <div v-if="arrow" class="base-popover-arrow" :style="arrowStyle" />

          <!-- Content Slot -->
          <div class="base-popover-content">
            <slot />
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { usePopover, type PopoverPlacement, type PopoverTrigger } from '@/composables/ui/usePopover'
import { computed, watch } from 'vue'

interface Props {
  modelValue?: boolean
  trigger?: PopoverTrigger
  placement?: PopoverPlacement
  offset?: number
  delay?: number
  arrow?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: false,
  trigger: 'click',
  placement: 'top',
  offset: 8,
  delay: 0,
  arrow: true
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const {
  isOpen,
  open,
  close,
  toggle,
  position,
  arrowPosition
} = usePopover({
  placement: props.placement,
  offset: props.offset,
  trigger: props.trigger,
  delay: props.delay,
  arrow: props.arrow,
  onOpen: () => emit('update:modelValue', true),
  onClose: () => emit('update:modelValue', false)
})

// Sync with v-model
watch(
  () => props.modelValue,
  (value) => {
    if (value && !isOpen.value) {
      open()
    } else if (!value && isOpen.value) {
      close()
    }
  }
)

// Sync internal state with v-model
watch(isOpen, (value) => {
  if (value !== props.modelValue) {
    emit('update:modelValue', value)
  }
})

const popoverStyle = computed(() => ({
  top: position.value.top,
  left: position.value.left,
  transform: position.value.transform
}))

const arrowStyle = computed(() => ({
  top: arrowPosition.value.top,
  left: arrowPosition.value.left,
  bottom: arrowPosition.value.bottom,
  right: arrowPosition.value.right,
  transform: arrowPosition.value.transform
}))

function handleTriggerClick() {
  if (props.trigger === 'click') {
    toggle()
  }
}

function handleTriggerMouseEnter() {
  if (props.trigger === 'hover') {
    open()
  }
}

function handleTriggerMouseLeave() {
  if (props.trigger === 'hover') {
    close()
  }
}

function handlePopoverMouseEnter() {
  if (props.trigger === 'hover') {
    // Keep popover open when hovering over it
    open()
  }
}

function handlePopoverMouseLeave() {
  if (props.trigger === 'hover') {
    close()
  }
}
</script>

<style scoped>
.base-popover-wrapper {
  display: inline-block;
}

.base-popover-trigger {
  display: inline-block;
}

.base-popover {
  position: fixed;
  z-index: var(--z-dropdown);
  background: rgba(15, 23, 42, 0.98);
  backdrop-filter: blur(16px);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-lg);
  max-width: min(400px, 90vw);
}

.base-popover-content {
  position: relative;
  z-index: 1;
}

/* Arrow */
.base-popover-arrow {
  position: absolute;
  width: 12px;
  height: 12px;
  background: rgba(15, 23, 42, 0.98);
  border: 1px solid var(--border-default);
  z-index: 0;
}

/* Arrow positioning for top placements */
.base-popover--top .base-popover-arrow,
.base-popover--top-start .base-popover-arrow,
.base-popover--top-end .base-popover-arrow {
  border-top: none;
  border-left: none;
  transform: translateX(-50%) rotate(45deg);
  margin-top: -6px;
}

/* Arrow positioning for bottom placements */
.base-popover--bottom .base-popover-arrow,
.base-popover--bottom-start .base-popover-arrow,
.base-popover--bottom-end .base-popover-arrow {
  border-bottom: none;
  border-right: none;
  transform: translateX(-50%) rotate(45deg);
  margin-bottom: -6px;
}

/* Arrow positioning for left placement */
.base-popover--left .base-popover-arrow {
  border-left: none;
  border-bottom: none;
  transform: translateY(-50%) rotate(45deg);
  margin-left: -6px;
}

/* Arrow positioning for right placement */
.base-popover--right .base-popover-arrow {
  border-right: none;
  border-top: none;
  transform: translateY(-50%) rotate(45deg);
  margin-right: -6px;
}

/* Transitions */
.popover-fade-enter-active,
.popover-fade-leave-active {
  transition: opacity var(--transition-fast), transform var(--transition-fast);
}

.popover-fade-enter-from,
.popover-fade-leave-to {
  opacity: 0;
}

.base-popover--top.popover-fade-enter-from,
.base-popover--top-start.popover-fade-enter-from,
.base-popover--top-end.popover-fade-enter-from {
  transform: translateY(8px);
}

.base-popover--bottom.popover-fade-enter-from,
.base-popover--bottom-start.popover-fade-enter-from,
.base-popover--bottom-end.popover-fade-enter-from {
  transform: translateY(-8px);
}

.base-popover--left.popover-fade-enter-from {
  transform: translateX(8px);
}

.base-popover--right.popover-fade-enter-from {
  transform: translateX(-8px);
}

.base-popover--top.popover-fade-leave-to,
.base-popover--top-start.popover-fade-leave-to,
.base-popover--top-end.popover-fade-leave-to {
  transform: translateY(8px);
}

.base-popover--bottom.popover-fade-leave-to,
.base-popover--bottom-start.popover-fade-leave-to,
.base-popover--bottom-end.popover-fade-leave-to {
  transform: translateY(-8px);
}

.base-popover--left.popover-fade-leave-to {
  transform: translateX(8px);
}

.base-popover--right.popover-fade-leave-to {
  transform: translateX(-8px);
}
</style>
