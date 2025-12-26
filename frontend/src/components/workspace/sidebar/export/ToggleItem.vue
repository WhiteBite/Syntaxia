<template>
  <label class="toggle-row">
    <span class="toggle-label">{{ label }}</span>
    <button
      type="button"
      class="toggle-switch"
      :class="{ active: modelValue }"
      @click="emit('update:modelValue', !modelValue)"
    >
      <span class="toggle-thumb" />
    </button>
  </label>
</template>

<script setup lang="ts">
defineProps<{
  modelValue: boolean
  label: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()
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
  font-size: 12px;
  color: #9ca3af;
  transition: color 0.15s ease-out;
}

.toggle-row:hover .toggle-label {
  color: #e5e7eb;
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
</style>
