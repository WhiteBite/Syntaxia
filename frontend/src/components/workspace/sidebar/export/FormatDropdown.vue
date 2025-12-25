<template>
  <div class="format-dropdown" ref="dropdownRef">
    <!-- Trigger Button -->
    <button 
      @click="toggleDropdown"
      class="format-trigger"
      :class="{ 'format-trigger-open': isOpen }"
      :title="t('context.changeFormat')"
    >
      <span class="format-label">{{ currentFormat.label }}</span>
      <svg 
        class="format-chevron" 
        :class="{ 'format-chevron-open': isOpen }"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
      </svg>
    </button>

    <!-- Dropdown Menu - Teleported to body to avoid overflow issues -->
    <Teleport to="body">
      <Transition name="dropdown">
        <div 
          v-if="isOpen" 
          class="format-menu"
          :style="menuStyle"
        >
          <button
            v-for="format in formats"
            :key="format.id"
            @click="selectFormat(format.id)"
            class="format-option"
            :class="{ 'format-option-active': modelValue === format.id }"
            :title="format.tooltip"
          >
            <span class="format-option-label">{{ format.label }}</span>
            <svg 
              v-if="modelValue === format.id" 
              class="format-check" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import type { OutputFormat } from '@/stores/settings.store'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const { t } = useI18n()

const props = defineProps<{
  modelValue: OutputFormat
}>()

const emit = defineEmits<{
  'update:modelValue': [value: OutputFormat]
}>()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
const menuStyle = ref<{ top: string; left: string }>({ top: '0px', left: '0px' })

const formats: { id: OutputFormat; label: string; tooltip: string }[] = [
  { id: 'xml', label: 'XML', tooltip: 'Best for AI - structured, easy to parse' },
  { id: 'markdown', label: 'Markdown', tooltip: 'Good for documentation and readability' },
  { id: 'plain', label: 'Plain', tooltip: 'Simple format with separators' },
]

const currentFormat = computed(() => 
  formats.find(f => f.id === props.modelValue) || formats[0]
)

function selectFormat(id: OutputFormat) {
  emit('update:modelValue', id)
  isOpen.value = false
}

function updateMenuPosition() {
  if (!dropdownRef.value) return
  const rect = dropdownRef.value.getBoundingClientRect()
  menuStyle.value = {
    top: `${rect.bottom + 4}px`,
    left: `${rect.left + rect.width / 2}px`
  }
}

function toggleDropdown() {
  if (!isOpen.value) {
    updateMenuPosition()
  }
  isOpen.value = !isOpen.value
}

function handleClickOutside(e: MouseEvent) {
  if (dropdownRef.value && !dropdownRef.value.contains(e.target as Node)) {
    isOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.format-dropdown {
  position: relative;
}

.format-trigger {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  padding: 0.25rem 0.5rem;
  font-size: 0.75rem;
  font-weight: 600;
  border-radius: 0.5rem;
  background: var(--accent-indigo-bg);
  color: white;
  border: 1px solid var(--accent-indigo-border);
  transition: all 150ms ease-out;
  cursor: pointer;
}

.format-trigger:hover,
.format-trigger-open {
  background: rgba(99, 102, 241, 0.35);
  border-color: rgba(99, 102, 241, 0.6);
}

.format-label {
  min-width: 2rem;
  text-align: center;
}

.format-chevron {
  width: 0.75rem;
  height: 0.75rem;
  transition: transform 150ms ease-out;
}

.format-chevron-open {
  transform: rotate(180deg);
}
</style>

<!-- Global styles for teleported menu -->
<style>
.format-menu {
  position: fixed;
  transform: translateX(-50%);
  z-index: 9999;
  min-width: 120px;
  padding: 0.25rem 0;
  border-radius: 0.75rem;
  background: var(--bg-2);
  border: 1px solid var(--border-strong);
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.4), 0 8px 10px -6px rgba(0, 0, 0, 0.3);
}

.format-option {
  width: 100%;
  padding: 0.375rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 500;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
  color: var(--text-secondary);
  transition: all 100ms ease-out;
  cursor: pointer;
  background: transparent;
  border: none;
}

.format-option:hover {
  background: var(--bg-3);
  color: var(--text-primary);
}

.format-option-active {
  color: var(--accent-primary);
  background: var(--accent-indigo-bg);
}

.format-option-label {
  flex: 1;
  text-align: left;
}

.format-check {
  width: 0.875rem;
  height: 0.875rem;
  color: var(--accent-primary);
}

/* Dropdown animation */
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 150ms ease-out;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateX(-50%) translateY(-4px);
}
</style>
