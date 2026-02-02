<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="base-modal-backdrop" @click="handleBackdropClick">
        <div 
          class="base-modal-container" 
          :class="sizeClass"
          @click.stop
          role="dialog"
          aria-modal="true"
          :aria-labelledby="title ? 'modal-title' : undefined"
        >
          <!-- Header -->
          <div v-if="$slots.header || title || showClose" class="base-modal-header">
            <slot name="header">
              <h2 v-if="title" id="modal-title" class="base-modal-title">{{ title }}</h2>
            </slot>
            <button
              v-if="showClose"
              @click="handleClose"
              class="base-modal-close"
              :aria-label="closeLabel"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <!-- Body -->
          <div class="base-modal-body">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="base-modal-footer">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, computed } from 'vue'

interface Props {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl' | 'full'
  closeOnBackdrop?: boolean
  closeOnEsc?: boolean
  showClose?: boolean
  closeLabel?: string
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closeOnBackdrop: true,
  closeOnEsc: true,
  showClose: true,
  closeLabel: 'Close'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'close'): void
}>()

const sizeClass = computed(() => `base-modal-size-${props.size}`)

function handleClose() {
  emit('update:modelValue', false)
  emit('close')
}

function handleBackdropClick() {
  if (props.closeOnBackdrop) {
    handleClose()
  }
}

function handleEscKey(event: KeyboardEvent) {
  if (event.key === 'Escape' && props.closeOnEsc && props.modelValue) {
    handleClose()
  }
}

onMounted(() => {
  if (props.closeOnEsc) {
    document.addEventListener('keydown', handleEscKey)
  }
})

onUnmounted(() => {
  if (props.closeOnEsc) {
    document.removeEventListener('keydown', handleEscKey)
  }
})
</script>

<style scoped>
.base-modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--z-modal-backdrop);
  backdrop-filter: blur(2px);
}

.base-modal-container {
  background: var(--bg-1);
  border-radius: var(--radius-2xl);
  box-shadow: var(--shadow-xl);
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border-subtle);
}

/* Size variants */
.base-modal-size-sm {
  max-width: 400px;
}

.base-modal-size-md {
  max-width: 500px;
}

.base-modal-size-lg {
  max-width: 700px;
}

.base-modal-size-xl {
  max-width: 900px;
}

.base-modal-size-full {
  max-width: 95vw;
  max-height: 95vh;
}

.base-modal-header {
  padding: var(--space-6);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-shrink: 0;
}

.base-modal-title {
  font-size: var(--font-size-xl);
  font-weight: var(--font-weight-semibold);
  color: var(--text-primary);
  margin: 0;
  flex: 1;
}

.base-modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: transparent;
  border: none;
  border-radius: var(--radius-md);
  color: var(--text-muted);
  cursor: pointer;
  transition: all var(--transition-fast);
  flex-shrink: 0;
}

.base-modal-close:hover {
  background: var(--bg-3);
  color: var(--text-primary);
}

.base-modal-body {
  padding: var(--space-6);
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

.base-modal-footer {
  padding: var(--space-6);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  flex-shrink: 0;
}

/* Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-normal);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .base-modal-container,
.modal-leave-active .base-modal-container {
  transition: transform var(--transition-normal) var(--ease-out);
}

.modal-enter-from .base-modal-container,
.modal-leave-to .base-modal-container {
  transform: scale(0.95);
}
</style>
