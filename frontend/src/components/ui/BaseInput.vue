<template>
  <div class="base-input-wrapper" :class="{ 'base-input-wrapper--disabled': disabled }">
    <label v-if="label" class="base-input-label">{{ label }}</label>
    <div 
      class="base-input-container" 
      :class="{ 
        'base-input-container--focused': isFocused,
        'base-input-container--ghost': variant === 'ghost'
      }"
    >
      <span v-if="$slots.prefix || prefixIcon" class="base-input__prefix">
        <slot name="prefix">
          <component :is="prefixIcon" />
        </slot>
      </span>
      <input
        ref="inputRef"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        class="base-input"
        @input="handleInput"
        @focus="isFocused = true"
        @blur="isFocused = false"
      />
      <span v-if="$slots.suffix || suffixIcon" class="base-input__suffix">
        <slot name="suffix">
          <component :is="suffixIcon" />
        </slot>
      </span>
    </div>
    <span v-if="error" class="base-input-error">{{ error }}</span>
  </div>
</template>

<script setup lang="ts">
import { ref, type Component } from 'vue'

interface Props {
  modelValue: string | number
  label?: string
  placeholder?: string
  type?: string
  disabled?: boolean
  error?: string
  prefixIcon?: Component
  suffixIcon?: Component
  variant?: 'default' | 'ghost'
}

const _props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  variant: 'default'
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const isFocused = ref(false)
const inputRef = ref<HTMLInputElement | null>(null)

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

defineExpose({
  input: inputRef,
  focus: () => inputRef.value?.focus()
})
</script>

<style scoped>
.base-input-wrapper {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  width: 100%;
}

.base-input-label {
  font-size: var(--font-size-xs);
  font-weight: var(--font-weight-medium);
  color: var(--text-muted);
  padding-left: var(--space-1);
}

.base-input-container {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--bg-1);
  border: 1px solid var(--border-default);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  min-height: calc(38px * var(--ui-scale));
}

.base-input-container:hover {
  border-color: var(--border-strong);
  background: var(--bg-2);
}

.base-input-container--ghost {
  background: transparent;
  border-color: transparent;
  padding-left: 0;
  padding-right: 0;
}

.base-input-container--ghost:hover {
  background: rgba(255, 255, 255, 0.03);
}

.base-input-container--ghost.base-input-container--focused {
  background: var(--bg-2);
  border-color: var(--accent-indigo);
  padding-left: var(--space-2);
  padding-right: var(--space-2);
}

.base-input-container--focused {
  border-color: var(--accent-indigo);
  background: var(--bg-2);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
  transform: translateY(-1px);
}

.base-input {
  flex: 1;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-size: var(--font-size-sm);
  font-family: var(--font-sans);
  outline: none;
  width: 100%;
}

.base-input::placeholder {
  color: var(--text-subtle);
}

.base-input__prefix,
.base-input__suffix {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-muted);
}

.base-input__prefix :deep(svg),
.base-input__suffix :deep(svg) {
  width: 1.1em;
  height: 1.1em;
}

.base-input-error {
  font-size: var(--font-size-xs);
  color: var(--color-danger);
  padding-left: var(--space-1);
}

.base-input-wrapper--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.base-input-wrapper--disabled .base-input {
  cursor: not-allowed;
}
</style>
