<template>
  <button
    type="button"
    role="switch"
    :aria-checked="modelValue"
    :class="[
      'toggle-switch',
      modelValue ? 'toggle-switch--active' : '',
      disabled ? 'toggle-switch--disabled' : ''
    ]"
    :disabled="disabled"
    @click="toggle"
  >
    <span class="toggle-switch__track">
      <span class="toggle-switch__thumb" />
    </span>
  </button>
</template>

<script setup lang="ts">
interface Props {
  modelValue: boolean
  disabled?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  disabled: false
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'change', value: boolean): void
}>()

function toggle() {
  if (props.disabled) return
  const newValue = !props.modelValue
  emit('update:modelValue', newValue)
  emit('change', newValue)
}
</script>

<style scoped>
.toggle-switch {
  position: relative;
  flex-shrink: 0;
  width: 36px;
  height: 20px;
  padding: 0;
  border: none;
  background: transparent;
  cursor: pointer;
  outline: none;
}

.toggle-switch__track {
  display: block;
  width: 100%;
  height: 100%;
  border-radius: 10px;
  background: rgba(100, 116, 139, 0.4);
  border: 1px solid rgba(100, 116, 139, 0.3);
  transition: all 200ms ease-out;
}

.toggle-switch__thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #94a3b8;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
  transition: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
}

/* Active state */
.toggle-switch--active .toggle-switch__track {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-color: rgba(139, 92, 246, 0.5);
  box-shadow: 0 0 12px rgba(99, 102, 241, 0.3);
}

.toggle-switch--active .toggle-switch__thumb {
  left: 18px;
  background: #fff;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

/* Hover states */
.toggle-switch:hover:not(.toggle-switch--disabled) .toggle-switch__track {
  border-color: rgba(148, 163, 184, 0.5);
}

.toggle-switch--active:hover:not(.toggle-switch--disabled) .toggle-switch__track {
  box-shadow: 0 0 16px rgba(99, 102, 241, 0.4);
}

/* Disabled state */
.toggle-switch--disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

/* Focus state */
.toggle-switch:focus-visible .toggle-switch__track {
  outline: 2px solid rgba(99, 102, 241, 0.5);
  outline-offset: 2px;
}
</style>
