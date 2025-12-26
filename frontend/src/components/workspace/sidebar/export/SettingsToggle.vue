<template>
  <label class="toggle-label" :title="hint">
    <div 
      class="toggle-switch"
      :class="{ 'toggle-switch-active': modelValue }"
      @click.prevent="handleToggle"
    >
      <div 
        class="toggle-switch-thumb"
        :class="{ 'toggle-switch-thumb-active': modelValue }"
      />
    </div>
    <span class="toggle-text">{{ label }}</span>
    <!-- Savings badge inline -->
    <span v-if="savings" class="toggle-savings" :class="{ 'toggle-savings-active': modelValue }">
      {{ savings }}
    </span>
  </label>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const props = defineProps<{
  modelValue: boolean
  label: string
  hint?: string
  savings?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const isAnimating = ref(false)

function handleToggle() {
  if (isAnimating.value) return
  isAnimating.value = true
  emit('update:modelValue', !props.modelValue)
  setTimeout(() => {
    isAnimating.value = false
  }, 250)
}
</script>

<style scoped>
.toggle-savings {
  margin-left: auto;
  font-size: 0.625rem;
  color: var(--text-muted);
  padding: 0.125rem 0.375rem;
  background: var(--bg-3);
  border-radius: 0.25rem;
  transition: all 150ms;
}

.toggle-savings-active {
  background: rgba(34, 197, 94, 0.15);
  color: #22c55e;
}
</style>
