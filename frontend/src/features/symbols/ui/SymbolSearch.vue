<template>
  <div class="symbol-search">
    <!-- Search Input -->
    <div class="search-wrapper">
      <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input
        ref="inputRef"
        v-model="searchQuery"
        type="text"
        :placeholder="t('symbols.searchPlaceholder')"
        class="search-input"
        @input="handleInput"
        @keydown.escape="handleClear"
      />
      <button
        v-if="searchQuery"
        class="clear-btn"
        :title="t('symbols.clear')"
        @click="handleClear"
      >
        <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Kind Filters -->
    <div class="kind-filters">
      <button
        v-for="kind in availableKinds"
        :key="kind.value"
        class="kind-btn"
        :class="{ 'kind-btn--active': selectedKinds.includes(kind.value) }"
        :title="kind.label"
        @click="toggleKind(kind.value)"
      >
        <component :is="kind.icon" class="kind-icon" />
        <span v-if="kindStats[kind.value] > 0" class="kind-count">
          {{ kindStats[kind.value] }}
        </span>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { computed, h, ref, watch } from 'vue'
import type { SymbolKind } from '../types'

const props = defineProps<{
  modelValue: string
  selectedKinds: SymbolKind[]
  kindStats: Record<SymbolKind, number>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
  (e: 'search', query: string): void
  (e: 'toggle-kind', kind: SymbolKind): void
  (e: 'clear'): void
}>()

const { t } = useI18n()
const inputRef = ref<HTMLInputElement | null>(null)
const searchQuery = ref(props.modelValue)

// Debounce timer
let debounceTimer: ReturnType<typeof setTimeout> | null = null

// Icon components for each kind
const ClassIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10' })
])

const FunctionIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4' })
])

const MethodIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z' })
])

const InterfaceIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2z' })
])

const VariableIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 7v10c0 2 1 3 3 3h10c2 0 3-1 3-3V7c0-2-1-3-3-3H7c-2 0-3 1-3 3z' }),
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M8 12h8' })
])

const ConstantIcon = () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
  h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M7 20l4-16m2 16l4-16M6 9h14M4 15h14' })
])

const availableKinds = computed(() => [
  { value: 'class' as SymbolKind, label: t('symbols.kinds.class'), icon: ClassIcon },
  { value: 'function' as SymbolKind, label: t('symbols.kinds.function'), icon: FunctionIcon },
  { value: 'method' as SymbolKind, label: t('symbols.kinds.method'), icon: MethodIcon },
  { value: 'interface' as SymbolKind, label: t('symbols.kinds.interface'), icon: InterfaceIcon },
  { value: 'variable' as SymbolKind, label: t('symbols.kinds.variable'), icon: VariableIcon },
  { value: 'constant' as SymbolKind, label: t('symbols.kinds.constant'), icon: ConstantIcon },
])

watch(() => props.modelValue, (newValue) => {
  searchQuery.value = newValue
})

function handleInput() {
  emit('update:modelValue', searchQuery.value)
  
  // Debounce search
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  
  debounceTimer = setTimeout(() => {
    emit('search', searchQuery.value)
  }, 300)
}

function handleClear() {
  searchQuery.value = ''
  emit('update:modelValue', '')
  emit('clear')
  inputRef.value?.focus()
}

function toggleKind(kind: SymbolKind) {
  emit('toggle-kind', kind)
}

function focus() {
  inputRef.value?.focus()
}

defineExpose({ focus })
</script>

<style scoped>
.symbol-search {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px 12px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 10px;
  width: 16px;
  height: 16px;
  color: #6b7280;
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 8px 32px 8px 34px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  color: #e5e7eb;
  transition: all 0.15s ease-out;
}

.search-input::placeholder {
  color: #6b7280;
}

.search-input:focus {
  outline: none;
  border-color: rgba(139, 92, 246, 0.5);
  background: rgba(255, 255, 255, 0.06);
}

.clear-btn {
  position: absolute;
  right: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  background: none;
  border: none;
  border-radius: 4px;
  color: #6b7280;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.clear-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e5e7eb;
}

.kind-filters {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
}

.kind-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px;
  color: #9ca3af;
  font-size: 11px;
  cursor: pointer;
  transition: all 0.15s ease-out;
}

.kind-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e7eb;
}

.kind-btn--active {
  background: rgba(139, 92, 246, 0.15);
  border-color: rgba(139, 92, 246, 0.3);
  color: #a78bfa;
}

.kind-icon {
  width: 14px;
  height: 14px;
}

.kind-count {
  font-size: 10px;
  padding: 1px 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
}

.kind-btn--active .kind-count {
  background: rgba(139, 92, 246, 0.2);
}
</style>
