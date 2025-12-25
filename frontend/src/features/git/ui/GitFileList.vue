<template>
  <div class="space-y-2">
    <div class="flex items-center justify-between">
      <span class="text-sm text-gray-300">{{ title }}</span>
      <span class="text-xs text-gray-400">{{ filteredFiles.length }}/{{ files.length }}</span>
    </div>

    <!-- Search & Filter Bar -->
    <div class="space-y-2">
      <div class="relative">
        <svg class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input 
          v-model="searchQuery" 
          type="text" 
          :placeholder="t('git.searchFiles')"
          class="input pl-10 py-1.5 text-sm w-full" 
        />
        <button 
          v-if="searchQuery" 
          @click="searchQuery = ''"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-white"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- File Type Filters -->
      <div class="flex flex-wrap gap-1.5">
        <button 
          v-for="filter in fileTypeFilters" 
          :key="filter.id" 
          @click="toggleFilter(filter.id)"
          :class="['git-filter-chip', activeFilters.has(filter.id) ? 'git-filter-chip-active' : '']"
        >
          <span class="git-filter-icon">{{ filter.icon }}</span>
          <span>{{ filter.label }}</span>
          <span class="git-filter-count">{{ getFilterCount(filter.id) }}</span>
        </button>
      </div>
    </div>

    <!-- File Tree -->
    <div class="file-list-container">
      <SimpleFileTree 
        :files="filteredFiles" 
        :selected-paths="selectedFiles"
        @toggle-select="$emit('toggle-select', $event)" 
        @select-folder="$emit('select-folder', $event)"
        @preview-file="$emit('preview-file', $event)" 
      />
    </div>

    <!-- Selection Actions -->
    <div class="flex gap-2">
      <button @click="$emit('select-all')" class="action-btn action-btn-success btn-sm flex-1">
        {{ t('git.selectAll') }}
      </button>
      <button @click="$emit('clear-selection')" class="action-btn action-btn-danger btn-sm flex-1">
        {{ t('git.clearSelection') }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import SimpleFileTree from '@/components/SimpleFileTree.vue';
import { useI18n } from '@/composables/useI18n';
import { computed, ref } from 'vue';

const props = defineProps<{
  title: string
  files: string[]
  selectedFiles: Set<string>
}>()

defineEmits<{
  'toggle-select': [path: string]
  'select-folder': [files: string[]]
  'select-all': []
  'clear-selection': []
  'preview-file': [path: string]
}>()

const { t } = useI18n()

const searchQuery = ref('')
const activeFilters = ref<Set<string>>(new Set())

const fileTypeFilters = [
  { id: 'code', label: 'Code', icon: '💻', extensions: ['.ts', '.js', '.tsx', '.jsx', '.vue', '.go', '.py', '.java', '.rs', '.cpp', '.c', '.h', '.cs', '.rb', '.php', '.swift', '.kt'] },
  { id: 'styles', label: 'Styles', icon: '🎨', extensions: ['.css', '.scss', '.sass', '.less', '.styl'] },
  { id: 'config', label: 'Config', icon: '⚙️', extensions: ['.json', '.yaml', '.yml', '.toml', '.xml', '.ini', '.env', '.config'] },
  { id: 'docs', label: 'Docs', icon: '📄', extensions: ['.md', '.txt', '.rst', '.adoc'] },
]

const filteredFiles = computed(() => {
  let result = props.files

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(f => f.toLowerCase().includes(query))
  }

  if (activeFilters.value.size > 0) {
    const activeExtensions = new Set<string>()
    activeFilters.value.forEach(filterId => {
      const filter = fileTypeFilters.find(f => f.id === filterId)
      if (filter) filter.extensions.forEach(ext => activeExtensions.add(ext))
    })
    result = result.filter(f => {
      const ext = '.' + f.split('.').pop()?.toLowerCase()
      return activeExtensions.has(ext)
    })
  }

  return result
})

function getFilterCount(filterId: string): number {
  const filter = fileTypeFilters.find(f => f.id === filterId)
  if (!filter) return 0
  return props.files.filter(f => {
    const ext = '.' + f.split('.').pop()?.toLowerCase()
    return filter.extensions.includes(ext)
  }).length
}

function toggleFilter(filterId: string) {
  const newFilters = new Set(activeFilters.value)
  if (newFilters.has(filterId)) {
    newFilters.delete(filterId)
  } else {
    newFilters.add(filterId)
  }
  activeFilters.value = newFilters
}
</script>
