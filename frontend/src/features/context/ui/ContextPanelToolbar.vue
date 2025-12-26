<template>
  <div v-if="visible" class="search-bar">
    <div class="relative">
      <svg class="search-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
          d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
      <input 
        ref="searchInputRef"
        v-model="localQuery" 
        type="text" 
        :placeholder="t('context.search')" 
        class="input pl-8 pr-20 text-sm"
        @keyup.enter="$emit('search-next')" 
        @keyup.escape="$emit('close')"
      />
      <div v-if="localQuery && resultsCount > 0" class="search-nav">
        <span class="search-count">{{ currentIndex + 1 }}/{{ resultsCount }}</span>
        <button @click="$emit('search-prev')" class="search-nav-btn">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
          </svg>
        </button>
        <button @click="$emit('search-next')" class="search-nav-btn">
          <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
      </div>
      <div v-else-if="localQuery && resultsCount === 0" class="search-no-results">
        <span>{{ t('context.noResults') }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/composables/useI18n'
import { ref, watch, nextTick } from 'vue'

const props = defineProps<{
  visible: boolean
  searchQuery: string
  resultsCount: number
  currentIndex: number
}>()

const emit = defineEmits<{
  'update:searchQuery': [value: string]
  'search-next': []
  'search-prev': []
  'close': []
}>()

const { t } = useI18n()
const searchInputRef = ref<HTMLInputElement | null>(null)

// Local query with two-way binding
const localQuery = ref(props.searchQuery)

watch(() => props.searchQuery, (val) => {
  localQuery.value = val
})

watch(localQuery, (val) => {
  emit('update:searchQuery', val)
})

// Focus input when search opens
watch(() => props.visible, (show) => {
  if (show) {
    nextTick(() => searchInputRef.value?.focus())
  }
})

// Expose input ref for external focus
defineExpose({
  focus: () => searchInputRef.value?.focus()
})
</script>

<style scoped>
.search-bar {
  @apply p-2;
  background: var(--bg-1);
  border-bottom: 1px solid var(--border-default);
}

.search-icon {
  @apply absolute left-2.5 top-1/2 -translate-y-1/2 w-4 h-4;
  color: var(--text-muted);
}

.search-nav {
  @apply absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1;
}

.search-count {
  @apply text-xs;
  color: var(--text-muted);
}

.search-nav-btn {
  @apply p-1 rounded;
  color: var(--text-muted);
  transition: all 150ms ease-out;
}

.search-nav-btn:hover {
  background: var(--bg-2);
  color: var(--text-primary);
}

.search-no-results {
  @apply absolute right-2 top-1/2 -translate-y-1/2 text-xs;
  color: var(--text-subtle);
}
</style>
