<template>
  <div class="base-dropdown">
    <!-- Trigger Slot -->
    <div
      ref="triggerRef"
      class="base-dropdown__trigger"
      @click="handleTriggerClick"
      :aria-expanded="modelValue"
      :aria-haspopup="true"
      :aria-disabled="disabled"
    >
      <slot name="trigger" :isOpen="modelValue" :toggle="toggle" />
    </div>

    <!-- Dropdown Content - Teleported to body -->
    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="modelValue"
          ref="dropdownRef"
          class="base-dropdown__content"
          :style="dropdownStyle"
          role="menu"
          :aria-label="ariaLabel"
        >
          <slot :close="close" />
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import type { DropdownPlacement } from '@/composables/ui/useDropdown'

interface Props {
  modelValue: boolean
  placement?: DropdownPlacement
  offset?: number
  closeOnClick?: boolean
  disabled?: boolean
  ariaLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  placement: 'bottom-start',
  offset: 8,
  closeOnClick: true,
  disabled: false,
  ariaLabel: 'Dropdown menu'
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const triggerRef = ref<HTMLElement | null>(null)
const dropdownRef = ref<HTMLElement | null>(null)
const position = ref({ top: '0px', left: '0px' })

/**
 * Calculate dropdown position based on trigger element and placement
 */
function calculatePosition(): void {
  if (!triggerRef.value || !dropdownRef.value) return

  const triggerRect = triggerRef.value.getBoundingClientRect()
  const dropdownRect = dropdownRef.value.getBoundingClientRect()
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight

  let top = 0
  let left = 0

  // Calculate base position based on placement
  switch (props.placement) {
    case 'bottom-start':
      top = triggerRect.bottom + props.offset
      left = triggerRect.left
      break
    case 'bottom-end':
      top = triggerRect.bottom + props.offset
      left = triggerRect.right - dropdownRect.width
      break
    case 'top-start':
      top = triggerRect.top - dropdownRect.height - props.offset
      left = triggerRect.left
      break
    case 'top-end':
      top = triggerRect.top - dropdownRect.height - props.offset
      left = triggerRect.right - dropdownRect.width
      break
    case 'left':
      top = triggerRect.top
      left = triggerRect.left - dropdownRect.width - props.offset
      break
    case 'right':
      top = triggerRect.top
      left = triggerRect.right + props.offset
      break
  }

  // Adjust if dropdown goes outside viewport (horizontal)
  if (left + dropdownRect.width > viewportWidth) {
    left = viewportWidth - dropdownRect.width - 8
  }
  if (left < 8) {
    left = 8
  }

  // Adjust if dropdown goes outside viewport (vertical)
  if (top + dropdownRect.height > viewportHeight) {
    // Flip to top if there's more space
    if (triggerRect.top > viewportHeight - triggerRect.bottom) {
      top = triggerRect.top - dropdownRect.height - props.offset
    } else {
      top = viewportHeight - dropdownRect.height - 8
    }
  }
  if (top < 8) {
    top = 8
  }

  position.value = {
    top: `${top}px`,
    left: `${left}px`
  }
}

/**
 * Close dropdown
 */
function close(): void {
  emit('update:modelValue', false)
}

/**
 * Toggle dropdown
 */
function toggle(): void {
  if (props.disabled) return
  emit('update:modelValue', !props.modelValue)
}

/**
 * Handle trigger click
 */
function handleTriggerClick(event: MouseEvent): void {
  event.stopPropagation()
  toggle()
}

/**
 * Handle click outside to close dropdown
 */
function handleClickOutside(event: MouseEvent): void {
  if (!props.modelValue) return

  const target = event.target as Node

  // Check if click is outside both trigger and dropdown
  const isOutsideTrigger = triggerRef.value && !triggerRef.value.contains(target)
  const isOutsideDropdown = dropdownRef.value && !dropdownRef.value.contains(target)

  if (isOutsideTrigger && isOutsideDropdown) {
    close()
  }
}

/**
 * Handle click inside dropdown
 */
function handleDropdownClick(): void {
  if (props.closeOnClick) {
    close()
  }
}

/**
 * Handle escape key to close dropdown
 */
function handleEscape(event: KeyboardEvent): void {
  if (props.modelValue && event.key === 'Escape') {
    close()
    // Return focus to trigger
    const focusableElement = triggerRef.value?.querySelector('button, [tabindex]') as HTMLElement
    focusableElement?.focus()
  }
}

/**
 * Recalculate position on window resize or scroll
 */
function handleResize(): void {
  if (props.modelValue) {
    calculatePosition()
  }
}

// Watch for modelValue changes to calculate position
watch(() => props.modelValue, async (isOpen) => {
  if (isOpen) {
    await nextTick()
    calculatePosition()
    
    // Setup event listeners when opened
    document.addEventListener('click', handleClickOutside)
    document.addEventListener('keydown', handleEscape)
    window.addEventListener('resize', handleResize)
    window.addEventListener('scroll', handleResize, true)
    
    if (props.closeOnClick && dropdownRef.value) {
      dropdownRef.value.addEventListener('click', handleDropdownClick)
    }
  } else {
    // Cleanup event listeners when closed
    document.removeEventListener('click', handleClickOutside)
    document.removeEventListener('keydown', handleEscape)
    window.removeEventListener('resize', handleResize)
    window.removeEventListener('scroll', handleResize, true)
    
    if (dropdownRef.value) {
      dropdownRef.value.removeEventListener('click', handleDropdownClick)
    }
  }
})

// Computed style for dropdown positioning
const dropdownStyle = computed(() => ({
  top: position.value.top,
  left: position.value.left
}))
</script>

<style scoped>
.base-dropdown {
  position: relative;
  display: inline-block;
}

.base-dropdown__trigger {
  display: inline-block;
  cursor: pointer;
}

.base-dropdown__trigger[aria-disabled="true"] {
  cursor: not-allowed;
  opacity: 0.5;
}

.base-dropdown__content {
  position: fixed;
  z-index: var(--z-dropdown);
  min-width: min(200px, 80vw);
  max-width: min(400px, 90vw);
  max-height: min(400px, 80vh);
  overflow-y: auto;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
}

/* Scrollbar styling */
.base-dropdown__content::-webkit-scrollbar {
  width: 6px;
}

.base-dropdown__content::-webkit-scrollbar-track {
  background: transparent;
}

.base-dropdown__content::-webkit-scrollbar-thumb {
  background: var(--bg-3);
  border-radius: var(--radius-full);
}

.base-dropdown__content::-webkit-scrollbar-thumb:hover {
  background: var(--border-strong);
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all var(--transition-fast);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
