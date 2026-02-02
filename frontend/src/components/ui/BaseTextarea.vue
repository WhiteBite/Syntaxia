<template>
  <div class="base-textarea-wrapper" :class="{ 'base-textarea-wrapper--disabled': disabled }">
    <div v-if="label || maxLength" class="base-textarea-header">
      <label v-if="label" class="base-textarea-label">{{ label }}</label>
      <span v-if="maxLength" class="base-textarea-counter" :class="{ 'base-textarea-counter--limit': isNearLimit }">
        {{ characterCount }}/{{ maxLength }}
      </span>
    </div>
    <div class="base-textarea-container" :class="{ 'base-textarea-container--focused': isFocused, 'base-textarea-container--error': error }">
      <textarea
        ref="textareaRef"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :maxlength="maxLength"
        :rows="computedRows"
        class="base-textarea"
        :class="{ 'base-textarea--no-resize': !resize || autoResize }"
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
      />
    </div>
    <span v-if="error" class="base-textarea-error">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'

interface Props {
  modelValue: string
  label?: string
  placeholder?: string
  rows?: number
  maxLength?: number
  error?: string
  disabled?: boolean
  resize?: boolean
  autoResize?: boolean
  maxRows?: number
}

const props = withDefaults(defineProps<Props>(), {
  rows: 3,
  disabled: false,
  resize: true,
  autoResize: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const isFocused = ref(false)
const textareaRef = ref<HTMLTextAreaElement | null>(null)
const computedRows = ref(props.rows)

const characterCount = computed(() => props.modelValue.length)

const isNearLimit = computed(() => {
  if (!props.maxLength) return false
  return characterCount.value >= props.maxLength * 0.9
})

function handleInput(event: Event) {
  const target = event.target as HTMLTextAreaElement
  emit('update:modelValue', target.value)
  
  if (props.autoResize) {
    adjustHeight()
  }
}

function handleFocus() {
  isFocused.value = true
}

function handleBlur() {
  isFocused.value = false
}

function adjustHeight() {
  if (!textareaRef.value || !props.autoResize) return
  
  // Reset height to auto to get the correct scrollHeight
  textareaRef.value.style.height = 'auto'
  
  const lineHeight = parseInt(getComputedStyle(textareaRef.value).lineHeight)
  const minHeight = lineHeight * props.rows
  const maxHeight = props.maxRows ? lineHeight * props.maxRows : Infinity
  
  const newHeight = Math.min(Math.max(textareaRef.value.scrollHeight, minHeight), maxHeight)
  
  textareaRef.value.style.height = `${newHeight}px`
}

// Watch for external value changes when autoResize is enabled
watch(() => props.modelValue, () => {
  if (props.autoResize) {
    nextTick(() => adjustHeight())
  }
})

// Initialize auto-resize on mount
watch(() => textareaRef.value, (textarea) => {
  if (textarea && props.autoResize) {
    nextTick(() => adjustHeight())
  }
}, { immediate: true })

defineExpose({
  textarea: textareaRef,
  focus: () => textareaRef.value?.focus()
})
</script>

<style scoped>
.base-textarea-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  width: 100%;
}

.base-textarea-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 var(--space-1);
}

.base-textarea-label {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--text-muted);
}

.base-textarea-counter {
  font-size: var(--font-size-xs);
  color: var(--text-subtle);
  transition: color var(--transition-fast);
}

.base-textarea-counter--limit {
  color: var(--color-warning);
  font-weight: var(--font-weight-medium);
}

.base-textarea-container {
  position: relative;
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.base-textarea-container:hover:not(.base-textarea-container--focused) {
  border-color: var(--border-strong);
  background: var(--bg-2);
}

.base-textarea-container--focused {
  border-color: var(--accent-indigo);
  background: var(--bg-2);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
  transform: translateY(-1px);
}

.base-textarea-container--error {
  border-color: var(--color-danger);
}

.base-textarea-container--error.base-textarea-container--focused {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.15);
}

.base-textarea {
  width: 100%;
  padding: var(--space-3);
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  font-family: var(--font-sans);
  line-height: 1.5;
  outline: none;
  resize: vertical;
  transition: height var(--transition-normal);
}

.base-textarea::placeholder {
  color: var(--text-subtle);
}

.base-textarea--no-resize {
  resize: none;
}

.base-textarea-error {
  font-size: var(--font-size-xs);
  color: var(--color-danger);
  padding-left: var(--space-1);
}

.base-textarea-wrapper--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.base-textarea-wrapper--disabled .base-textarea {
  cursor: not-allowed;
}
</style>
